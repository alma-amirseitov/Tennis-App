import { api } from './client';

export type ChatType = 'personal' | 'community' | 'event';

export interface ChatUser {
  id: string;
  first_name: string;
  avatar_url: string | null;
}

export interface ChatCommunity {
  id: string;
  name: string;
  logo_url: string | null;
}

export interface LastMessage {
  content: string;
  sender_id: string;
  created_at: string;
}

export interface ChatItem {
  id: string;
  chat_type: ChatType;
  other_user?: ChatUser;
  community?: ChatCommunity;
  last_message: LastMessage | null;
  unread_count: number;
  is_muted: boolean;
}

export interface MessageItem {
  id: string;
  sender: ChatUser;
  content: string;
  reply_to: string | null;
  created_at: string;
}

export interface MessagesResponse {
  data: MessageItem[];
  has_more: boolean;
}

function unwrapData<T>(response: { data: { data: T } }): T {
  return response.data.data;
}

export async function getChats(): Promise<ChatItem[]> {
  const res = await api.get<{ data: ChatItem[] }>('/chats');
  return unwrapData(res);
}

export async function createPersonalChat(userId: string): Promise<{ chat_id: string; is_new: boolean }> {
  const res = await api.post<{ data: { chat_id: string; is_new: boolean } }>('/chats/personal', {
    user_id: userId,
  });
  return unwrapData(res);
}

export async function getMessages(
  chatId: string,
  params?: { before?: string; limit?: number }
): Promise<MessagesResponse> {
  const res = await api.get<{ data: MessagesResponse }>(`/chats/${chatId}/messages`, { params });
  return unwrapData(res);
}

export async function sendMessage(
  chatId: string,
  content: string,
  replyToId?: string | null
): Promise<MessageItem> {
  const res = await api.post<{ data: MessageItem }>(`/chats/${chatId}/messages`, {
    content,
    reply_to_id: replyToId ?? null,
  });
  return unwrapData(res);
}

export async function markChatRead(chatId: string, lastReadAt: string): Promise<void> {
  await api.post(`/chats/${chatId}/read`, { last_read_at: lastReadAt });
}

export async function muteChat(chatId: string, isMuted: boolean): Promise<void> {
  await api.patch(`/chats/${chatId}/mute`, { is_muted: isMuted });
}

export async function getUnreadCount(): Promise<number> {
  const res = await api.get<{ data: { total_unread: number } }>('/chats/unread-count');
  return unwrapData(res).total_unread;
}
