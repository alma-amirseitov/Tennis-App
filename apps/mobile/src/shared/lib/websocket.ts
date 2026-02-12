import Constants from 'expo-constants';
import { authStore } from '@/shared/stores/auth';

type WSMessageHandler = (data: unknown) => void;

const API_BASE =
  (Constants.expoConfig?.extra as { apiUrl?: string } | undefined)?.apiUrl ??
  (Constants.manifest?.extra as { apiUrl?: string } | undefined) ??
  'http://localhost:8080';

function getWsUrl(): string {
  const base = API_BASE.replace(/^http/, 'ws');
  return `${base}/ws`;
}

const BACKOFF_DELAYS = [1000, 2000, 4000, 8000, 16000, 30000];
const HEARTBEAT_INTERVAL = 30000;

class WebSocketManager {
  private ws: WebSocket | null = null;
  private url = '';
  private token = '';
  private reconnectAttempt = 0;
  private heartbeatTimer: ReturnType<typeof setInterval> | null = null;
  private handlers: Map<string, Set<WSMessageHandler>> = new Map();
  private connectResolve: (() => void) | null = null;

  connect(accessToken: string): Promise<void> {
    this.token = accessToken;
    this.url = `${getWsUrl()}?token=${encodeURIComponent(accessToken)}`;

    return new Promise((resolve, reject) => {
      try {
        this.ws = new WebSocket(this.url);
        this.connectResolve = resolve;

        this.ws.onopen = () => {
          this.reconnectAttempt = 0;
          this.startHeartbeat();
          this.connectResolve?.();
          this.connectResolve = null;
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const msg = JSON.parse(event.data as string) as { type: string; data?: unknown };
            this.dispatch(msg.type, msg.data ?? msg);
            if (msg.type === 'pong') {
              // heartbeat ack
            }
          } catch {
            // ignore parse errors
          }
        };

        this.ws.onerror = () => {
          reject(new Error('WebSocket error'));
        };

        this.ws.onclose = () => {
          this.stopHeartbeat();
          if (this.connectResolve) {
            this.connectResolve = null;
            reject(new Error('WebSocket closed'));
          }
          this.scheduleReconnect();
        };
      } catch (err) {
        reject(err);
      }
    });
  }

  private scheduleReconnect() {
    const delay = BACKOFF_DELAYS[Math.min(this.reconnectAttempt, BACKOFF_DELAYS.length - 1)];
    this.reconnectAttempt++;

    setTimeout(() => {
      const token = authStore.getState().accessToken;
      if (token) {
        this.connect(token).catch(() => {});
      }
    }, delay);
  }

  private startHeartbeat() {
    this.stopHeartbeat();
    this.heartbeatTimer = setInterval(() => {
      this.send({ type: 'ping' });
    }, HEARTBEAT_INTERVAL);
  }

  private stopHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer);
      this.heartbeatTimer = null;
    }
  }

  private dispatch(type: string, data: unknown) {
    const handlers = this.handlers.get(type);
    if (handlers) {
      handlers.forEach((h) => h(data));
    }
    this.handlers.get('*')?.forEach((h) => h({ type, data }));
  }

  on(type: string, handler: WSMessageHandler): () => void {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, new Set());
    }
    this.handlers.get(type)!.add(handler);
    return () => this.handlers.get(type)?.delete(handler);
  }

  send(payload: object) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(payload));
    }
  }

  sendMessage(chatId: string, content: string, replyToId?: string | null, clientId?: string) {
    this.send({
      type: 'message',
      chat_id: chatId,
      content,
      reply_to: replyToId ?? null,
      client_id: clientId ?? `temp-${Date.now()}-${Math.random().toString(36).slice(2, 11)}`,
    });
  }

  sendTyping(chatId: string) {
    this.send({ type: 'typing', chat_id: chatId });
  }

  sendRead(chatId: string) {
    this.send({ type: 'read', chat_id: chatId });
  }

  disconnect() {
    this.stopHeartbeat();
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.handlers.clear();
    this.reconnectAttempt = 0;
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}

export const wsManager = new WebSocketManager();
