import { useState, useMemo } from 'react';
import {
  View,
  Text,
  FlatList,
  Pressable,
  TextInput,
  StyleSheet,
  RefreshControl,
} from 'react-native';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { useChats } from '@/shared/api/hooks';
import { ScreenHeader, Avatar, Skeleton, ErrorState, EmptyState } from '@/shared/ui';
import { colors, spacing, typography, radius } from '@/shared/theme';
import type { ChatItem, ChatType } from '@/shared/api/chat';

type TabType = 'personal' | 'community' | 'event';

function formatTime(iso: string): string {
  const d = new Date(iso);
  const now = new Date();
  const diff = now.getTime() - d.getTime();
  if (diff < 86400000) return d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
  if (diff < 172800000) return 'Вчера';
  return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
}

function getChatName(chat: ChatItem): string {
  if (chat.chat_type === 'personal' && chat.other_user) {
    return chat.other_user.first_name;
  }
  if (chat.chat_type === 'community' && chat.community) {
    return chat.community.name;
  }
  return 'Chat';
}

function getChatAvatar(chat: ChatItem): string | null {
  if (chat.chat_type === 'personal' && chat.other_user) {
    return chat.other_user.avatar_url;
  }
  if (chat.chat_type === 'community' && chat.community) {
    return chat.community.logo_url;
  }
  return null;
}

interface ChatCardProps {
  chat: ChatItem;
}

function ChatCard({ chat }: ChatCardProps) {
  const router = useRouter();
  const name = getChatName(chat);
  const avatar = getChatAvatar(chat);
  const lastMsg = chat.last_message?.content ?? '';
  const time = chat.last_message?.created_at ? formatTime(chat.last_message.created_at) : '';

  return (
    <Pressable
      onPress={() => router.push(`/chat/${chat.id}`)}
      style={({ pressed }) => [styles.card, pressed && { opacity: 0.9 }]}
    >
      <Avatar uri={avatar} name={name} size="md" />
      <View style={styles.info}>
        <View style={styles.topRow}>
          <Text style={styles.name} numberOfLines={1}>{name}</Text>
          <Text style={styles.time}>{time}</Text>
        </View>
        <View style={styles.bottomRow}>
          <Text style={styles.preview} numberOfLines={1}>{lastMsg || '—'}</Text>
          {chat.unread_count > 0 && (
            <View style={styles.badge}>
              <Text style={styles.badgeText}>{chat.unread_count > 99 ? '99+' : chat.unread_count}</Text>
            </View>
          )}
        </View>
      </View>
    </Pressable>
  );
}

export default function ChatListScreen() {
  const { t } = useTranslation();
  const [tab, setTab] = useState<TabType>('personal');
  const [search, setSearch] = useState('');

  const { data: chats, isLoading, isError, refetch } = useChats();

  const filteredChats = useMemo(() => {
    if (!chats) return [];
    let list = chats.filter((c) => c.chat_type === tab);
    if (search.trim()) {
      const q = search.toLowerCase();
      list = list.filter((c) => {
        const name = getChatName(c);
        return name.toLowerCase().includes(q);
      });
    }
    list.sort((a, b) => {
      const aTime = a.last_message?.created_at ?? '';
      const bTime = b.last_message?.created_at ?? '';
      return bTime.localeCompare(aTime);
    });
    return list;
  }, [chats, tab, search]);

  const tabLabels: Record<TabType, string> = {
    personal: t('chat.tab_personal'),
    community: t('chat.tab_groups'),
    event: t('chat.tab_events'),
  };

  const counts = useMemo(() => {
    if (!chats) return { personal: 0, community: 0, event: 0 };
    return {
      personal: chats.filter((c) => c.chat_type === 'personal').length,
      community: chats.filter((c) => c.chat_type === 'community').length,
      event: chats.filter((c) => c.chat_type === 'event').length,
    };
  }, [chats]);

  const unreadByTab = useMemo(() => {
    if (!chats) return { personal: 0, community: 0, event: 0 };
    return {
      personal: chats.filter((c) => c.chat_type === 'personal').reduce((s, c) => s + c.unread_count, 0),
      community: chats.filter((c) => c.chat_type === 'community').reduce((s, c) => s + c.unread_count, 0),
      event: chats.filter((c) => c.chat_type === 'event').reduce((s, c) => s + c.unread_count, 0),
    };
  }, [chats]);

  if (isLoading) {
    return (
      <View style={styles.container}>
        <ScreenHeader title={t('chat.title')} showBack />
        <View style={styles.skeleton}>
          <Skeleton width="100%" height={72} radius={12} />
          <Skeleton width="100%" height={72} radius={12} />
          <Skeleton width="100%" height={72} radius={12} />
        </View>
      </View>
    );
  }

  if (isError) {
    return (
      <View style={styles.container}>
        <ScreenHeader title={t('chat.title')} showBack />
        <ErrorState onRetry={refetch} />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <ScreenHeader title={t('chat.title')} showBack />
      <TextInput
        style={styles.search}
        placeholder={t('chat.search_placeholder')}
        placeholderTextColor={colors.textMuted}
        value={search}
        onChangeText={setSearch}
      />
      <View style={styles.tabs}>
        {(['personal', 'community', 'event'] as TabType[]).map((t) => (
          <Pressable
            key={t}
            onPress={() => setTab(t)}
            style={[styles.tab, tab === t && styles.tabActive]}
          >
            <Text style={[styles.tabText, tab === t && styles.tabTextActive]}>{tabLabels[t]}</Text>
            {(counts[t] > 0 || unreadByTab[t] > 0) && (
              <View style={[styles.tabBadge, unreadByTab[t] > 0 && styles.tabBadgeUnread]}>
                <Text style={[styles.tabBadgeText, unreadByTab[t] > 0 && styles.tabBadgeTextUnread]}>
                  {unreadByTab[t] > 0 ? unreadByTab[t] : counts[t]}
                </Text>
              </View>
            )}
          </Pressable>
        ))}
      </View>
      {filteredChats.length === 0 ? (
        <EmptyState emoji="💬" title={t('chat.empty')} />
      ) : (
        <FlatList
          data={filteredChats}
          keyExtractor={(item) => item.id}
          renderItem={({ item }) => <ChatCard chat={item} />}
          refreshControl={<RefreshControl refreshing={false} onRefresh={refetch} tintColor={colors.primary} />}
          contentContainerStyle={styles.list}
          showsVerticalScrollIndicator={false}
        />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  search: {
    height: 44,
    marginHorizontal: spacing.base,
    marginVertical: spacing.sm,
    borderRadius: radius.md,
    backgroundColor: colors.card,
    borderWidth: 1,
    borderColor: colors.border,
    paddingHorizontal: spacing.base,
    ...typography.textStyles.body,
    color: colors.text,
  },
  tabs: { flexDirection: 'row', borderBottomWidth: 1, borderBottomColor: colors.borderLight, paddingHorizontal: spacing.base },
  tab: { flex: 1, paddingVertical: spacing.md, alignItems: 'center', flexDirection: 'row', justifyContent: 'center', gap: spacing.xs },
  tabActive: { borderBottomWidth: 2, borderBottomColor: colors.primary },
  tabText: { ...typography.textStyles.bodySm, color: colors.textMuted },
  tabTextActive: { color: colors.primary, fontWeight: typography.fontWeight.semibold },
  tabBadge: { backgroundColor: colors.border, borderRadius: 10, paddingHorizontal: 6, paddingVertical: 2 },
  tabBadgeUnread: { backgroundColor: colors.primary },
  tabBadgeText: { ...typography.textStyles.caption, color: colors.textSecondary, fontSize: 10 },
  tabBadgeTextUnread: { color: colors.white },
  skeleton: { padding: spacing.base, gap: spacing.md },
  list: { paddingBottom: spacing['3xl'] },
  card: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: spacing.base,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: colors.borderLight,
    gap: spacing.md,
  },
  info: { flex: 1 },
  topRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' },
  name: { ...typography.textStyles.body, fontWeight: typography.fontWeight.semibold, color: colors.text, flex: 1 },
  time: { ...typography.textStyles.caption, color: colors.textMuted },
  bottomRow: { flexDirection: 'row', alignItems: 'center', marginTop: 2 },
  preview: { ...typography.textStyles.bodySm, color: colors.textSecondary, flex: 1 },
  badge: { backgroundColor: colors.primary, borderRadius: 10, paddingHorizontal: 6, paddingVertical: 2 },
  badgeText: { ...typography.textStyles.caption, color: colors.white, fontSize: 10 },
});
