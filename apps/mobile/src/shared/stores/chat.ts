import { create } from 'zustand';
import type { MessageItem } from '@/shared/api/chat';

interface ChatState {
  messagesByChatId: Record<string, MessageItem[]>;
  unreadCounts: Record<string, number>;
  typingByChatId: Record<string, string>; // chatId -> "FirstName" who is typing
  totalUnread: number;

  addMessage: (chatId: string, message: MessageItem) => void;
  setMessages: (chatId: string, messages: MessageItem[]) => void;
  prependMessages: (chatId: string, messages: MessageItem[]) => void;
  setUnreadCount: (chatId: string, count: number) => void;
  setTotalUnread: (count: number) => void;
  setTyping: (chatId: string, userName: string | null) => void;
  clearChat: (chatId: string) => void;
}

export const chatStore = create<ChatState>((set) => ({
  messagesByChatId: {},
  unreadCounts: {},
  typingByChatId: {},
  totalUnread: 0,

  addMessage: (chatId, message) =>
    set((s) => ({
      messagesByChatId: {
        ...s.messagesByChatId,
        [chatId]: [...(s.messagesByChatId[chatId] ?? []), message],
      },
    })),

  setMessages: (chatId, messages) =>
    set((s) => ({
      messagesByChatId: { ...s.messagesByChatId, [chatId]: messages },
    })),

  prependMessages: (chatId, messages) =>
    set((s) => {
      const existing = s.messagesByChatId[chatId] ?? [];
      const merged = [...messages, ...existing];
      const seen = new Set<string>();
      const unique = merged.filter((m) => {
        if (seen.has(m.id)) return false;
        seen.add(m.id);
        return true;
      });
      return {
        messagesByChatId: { ...s.messagesByChatId, [chatId]: unique },
      };
    }),

  setUnreadCount: (chatId, count) =>
    set((s) => ({
      unreadCounts: { ...s.unreadCounts, [chatId]: count },
    })),

  setTotalUnread: (count) => set({ totalUnread: count }),

  setTyping: (chatId, userName) =>
    set((s) => ({
      typingByChatId: userName
        ? { ...s.typingByChatId, [chatId]: userName }
        : (() => {
            const next = { ...s.typingByChatId };
            delete next[chatId];
            return next;
          })(),
    })),

  clearChat: (chatId) =>
    set((s) => {
      const next = { ...s.messagesByChatId };
      delete next[chatId];
      return { messagesByChatId: next };
    }),
}));
