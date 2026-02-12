import React, { useEffect, useRef, useCallback, useMemo, useState } from 'react';
import {
  View,
  Text,
  FlatList,
  TextInput,
  Pressable,
  StyleSheet,
  KeyboardAvoidingView,
  Platform,
  ActivityIndicator,
} from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { useChatMessagesInfinite, useSendMessage } from '@/shared/api/hooks';
import { useChats } from '@/shared/api/hooks';
import { chatStore } from '@/shared/stores/chat';
import { authStore } from '@/shared/stores/auth';
import { wsManager } from '@/shared/lib/websocket';
import { ScreenHeader } from '@/shared/ui';
import { colors, spacing, typography, radius } from '@/shared/theme';
import type { MessageItem } from '@/shared/api/chat';

function formatMessageTime(iso: string): string {
  return new Date(iso).toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
}

interface MessageBubbleProps {
  msg: MessageItem;
  isMine: boolean;
}

function MessageBubble({ msg, isMine }: MessageBubbleProps) {
  return (
    <View style={[styles.bubbleWrap, isMine ? styles.bubbleWrapMine : styles.bubbleWrapOther]}>
      <View style={[styles.bubble, isMine ? styles.bubbleMine : styles.bubbleOther]}>
        <Text style={[styles.bubbleText, isMine && styles.bubbleTextMine]}>{msg.content}</Text>
        <Text style={[styles.bubbleTime, isMine ? styles.bubbleTimeMine : styles.bubbleTimeOther]}>
          {formatMessageTime(msg.created_at)}
        </Text>
      </View>
    </View>
  );
}

export default function ChatDetailScreen() {
  const { id: chatId } = useLocalSearchParams<{ id: string }>();
  const router = useRouter();
  const { t } = useTranslation();

  const currentUserId = authStore((s) => s.user?.id) ?? '';
  const accessToken = authStore((s) => s.accessToken);
  const messages = chatStore((s) => s.messagesByChatId[chatId ?? ''] ?? []);
  const typingUser = chatStore((s) => s.typingByChatId[chatId ?? ''] ?? null);
  const setMessages = chatStore((s) => s.setMessages);
  const addMessage = chatStore((s) => s.addMessage);
  const prependMessages = chatStore((s) => s.prependMessages);
  const setTyping = chatStore((s) => s.setTyping);

  const [input, setInput] = useState('');
  const flatListRef = useRef<FlatList>(null);
  const typingTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastTypingSentRef = useRef(0);

  const { data: chats } = useChats();
  const chat = chats?.find((c) => c.id === chatId);

  const {
    data: infiniteData,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useChatMessagesInfinite(chatId ?? '', 50);

  const sendMessageMutation = useSendMessage();

  const displayMessages = messages;

  useEffect(() => {
    const pages = infiniteData?.pages;
    if (!pages?.length || !chatId) return;
    const first = pages[0].data;
    if (first.length > 0) setMessages(chatId, first);
  }, [infiniteData?.pages?.[0], chatId, setMessages]);

  const prevPagesLengthRef = useRef(0);
  useEffect(() => {
    const pages = infiniteData?.pages;
    if (!pages || pages.length <= 1 || pages.length <= prevPagesLengthRef.current) {
      prevPagesLengthRef.current = pages?.length ?? 0;
      return;
    }
    prevPagesLengthRef.current = pages.length;
    const lastPage = pages[pages.length - 1];
    if (lastPage.data.length > 0) prependMessages(chatId ?? '', lastPage.data);
  }, [infiniteData?.pages, chatId, prependMessages]);

  useEffect(() => {
    if (!chatId || !accessToken) return;

    if (!wsManager.isConnected()) {
      wsManager.connect(accessToken).catch(() => {});
    }

    const unsubMessage = wsManager.on('message', (data: unknown) => {
      const d = data as { chat_id?: string; id?: string; sender?: MessageItem['sender']; content?: string; reply_to?: string | null; created_at?: string };
      if (d.chat_id === chatId && d.id && d.sender && d.content != null && d.created_at) {
        addMessage(chatId, {
          id: d.id,
          sender: d.sender,
          content: d.content,
          reply_to: d.reply_to ?? null,
          created_at: d.created_at,
        });
      }
    });

    const unsubTyping = wsManager.on('typing', (data: unknown) => {
      const d = data as { chat_id?: string; first_name?: string };
      if (d.chat_id === chatId) {
        setTyping(chatId, d.first_name ?? null);
        if (typingTimeoutRef.current) clearTimeout(typingTimeoutRef.current);
        typingTimeoutRef.current = setTimeout(() => setTyping(chatId, null), 3000);
      }
    });

    wsManager.sendRead(chatId);

    return () => {
      unsubMessage();
      unsubTyping();
      if (typingTimeoutRef.current) clearTimeout(typingTimeoutRef.current);
    };
  }, [chatId, accessToken, addMessage, setTyping]);

  const sendTypingThrottled = useCallback(() => {
    if (!chatId || !wsManager.isConnected()) return;
    const now = Date.now();
    if (now - lastTypingSentRef.current < 2000) return;
    lastTypingSentRef.current = now;
    wsManager.sendTyping(chatId);
  }, [chatId]);

  const handleInputChange = useCallback(
    (text: string) => {
      setInput(text);
      if (text.trim()) sendTypingThrottled();
    },
    [sendTypingThrottled]
  );

  const handleSend = useCallback(() => {
    const text = input.trim();
    if (!text || !chatId) return;

    setInput('');
    if (wsManager.isConnected()) {
      wsManager.sendMessage(chatId, text);
    } else {
      sendMessageMutation.mutate(
        { chatId, content: text },
        {
          onSuccess: (msg) => addMessage(chatId, msg),
        }
      );
    }
  }, [input, chatId, sendMessageMutation, addMessage]);

  const handleLoadMore = useCallback(() => {
    if (hasNextPage && !isFetchingNextPage) fetchNextPage();
  }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  const chatName = useMemo(() => {
    if (!chat) return '';
    if (chat.chat_type === 'personal' && chat.other_user) return chat.other_user.first_name;
    if (chat.chat_type === 'community' && chat.community) return chat.community.name;
    return t('chat.title');
  }, [chat, t]);

  if (!chatId) return null;

  if (isLoading) {
    return (
      <View style={styles.container}>
        <ScreenHeader title="..." showBack />
        <View style={styles.loading}>
          <ActivityIndicator size="large" color={colors.primary} />
        </View>
      </View>
    );
  }

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : undefined}
      keyboardVerticalOffset={Platform.OS === 'ios' ? 90 : 0}
    >
      <ScreenHeader title={chatName} showBack />
      <FlatList
        ref={flatListRef}
        data={displayMessages}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <MessageBubble msg={item} isMine={item.sender.id === currentUserId} />
        )}
        inverted
        contentContainerStyle={styles.list}
        onEndReached={handleLoadMore}
        onEndReachedThreshold={0.3}
        ListFooterComponent={
          isFetchingNextPage ? (
            <View style={styles.loadMore}>
              <ActivityIndicator size="small" color={colors.primary} />
            </View>
          ) : null
        }
        ListHeaderComponent={
          typingUser ? (
            <View style={styles.typingWrap}>
              <Text style={styles.typingText}>{typingUser} {t('chat.typing')}</Text>
            </View>
          ) : null
        }
      />
      <View style={styles.inputBar}>
        <TextInput
          style={styles.input}
          placeholder={t('chat.type_message')}
          placeholderTextColor={colors.textMuted}
          value={input}
          onChangeText={handleInputChange}
          multiline
          maxLength={2000}
          onSubmitEditing={handleSend}
        />
        <Pressable
          onPress={handleSend}
          style={({ pressed }) => [styles.sendBtn, pressed && styles.sendBtnPressed]}
          disabled={!input.trim()}
        >
          <Text style={[styles.sendText, !input.trim() && styles.sendTextDisabled]}>➤</Text>
        </Pressable>
      </View>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  loading: { flex: 1, justifyContent: 'center', alignItems: 'center' },
  list: { paddingHorizontal: spacing.base, paddingBottom: spacing.md, paddingTop: spacing.sm },
  loadMore: { paddingVertical: spacing.md, alignItems: 'center' },
  bubbleWrap: { marginVertical: 2, flexDirection: 'row' },
  bubbleWrapMine: { justifyContent: 'flex-end' },
  bubbleWrapOther: { justifyContent: 'flex-start' },
  bubble: {
    maxWidth: '80%',
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    borderRadius: radius.lg,
  },
  bubbleMine: { backgroundColor: colors.primary },
  bubbleOther: { backgroundColor: colors.card, borderWidth: 1, borderColor: colors.border },
  bubbleText: { ...typography.textStyles.body, color: colors.text },
  bubbleTextMine: { color: colors.white },
  bubbleTime: { ...typography.textStyles.caption, marginTop: 2, fontSize: 10 },
  bubbleTimeMine: { color: 'rgba(255,255,255,0.8)' },
  bubbleTimeOther: { color: colors.textMuted },
  typingWrap: { paddingVertical: spacing.sm, paddingHorizontal: spacing.md },
  typingText: { ...typography.textStyles.caption, color: colors.textMuted, fontStyle: 'italic' },
  inputBar: {
    flexDirection: 'row',
    alignItems: 'flex-end',
    padding: spacing.sm,
    backgroundColor: colors.card,
    borderTopWidth: 1,
    borderTopColor: colors.borderLight,
    gap: spacing.sm,
  },
  input: {
    flex: 1,
    minHeight: 40,
    maxHeight: 100,
    borderRadius: radius.lg,
    backgroundColor: colors.background,
    borderWidth: 1,
    borderColor: colors.border,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    ...typography.textStyles.body,
    color: colors.text,
  },
  sendBtn: {
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: colors.primary,
    justifyContent: 'center',
    alignItems: 'center',
  },
  sendBtnPressed: { opacity: 0.8 },
  sendText: { color: colors.white, fontSize: 18 },
  sendTextDisabled: { opacity: 0.5 },
});
