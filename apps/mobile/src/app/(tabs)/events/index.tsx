import { useState, useCallback, useMemo } from 'react';
import { View, Text, Pressable, FlatList, StyleSheet, ActivityIndicator, RefreshControl } from 'react-native';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { useEvents } from '@/shared/api/hooks';
import { EventCard } from '@/features/events/EventCard';
import { EventFilters } from '@/features/events/EventFilters';
import { EmptyState, ErrorState, Skeleton, FAB } from '@/shared/ui';
import { colors, spacing, typography } from '@/shared/theme';
import type { EventItem, EventType } from '@/shared/api/events';

const INNER_TABS = ['feed', 'calendar', 'my'] as const;
type InnerTab = typeof INNER_TABS[number];

export default function EventsScreen() {
  const { t } = useTranslation();
  const router = useRouter();
  const [innerTab, setInnerTab] = useState<InnerTab>('feed');
  const [typeFilter, setTypeFilter] = useState<EventType | 'all'>('all');

  const params = useMemo(() => ({
    event_type: typeFilter === 'all' ? undefined : typeFilter,
    my_events: innerTab === 'my' ? true : undefined,
  }), [typeFilter, innerTab]);

  const { data, isLoading, isError, refetch, fetchNextPage, hasNextPage, isFetchingNextPage } = useEvents(params);
  const events = useMemo(() => data?.pages.flatMap((p) => p.data) ?? [], [data]);

  const tabLabels: Record<InnerTab, string> = {
    feed: t('events.tab_feed'),
    calendar: t('events.tab_calendar'),
    my: t('events.tab_my'),
  };

  const renderItem = useCallback(({ item }: { item: EventItem }) => <EventCard event={item} />, []);

  return (
    <View style={styles.container}>
      {/* Inner tabs */}
      <View style={styles.tabBar}>
        {INNER_TABS.map((tab) => (
          <Pressable key={tab} onPress={() => setInnerTab(tab)} style={[styles.tab, innerTab === tab && styles.tabActive]}>
            <Text style={[styles.tabText, innerTab === tab && styles.tabTextActive]}>{tabLabels[tab]}</Text>
          </Pressable>
        ))}
      </View>

      {/* Filters (only for feed) */}
      {innerTab === 'feed' ? (
        <EventFilters selectedType={typeFilter} onTypeChange={setTypeFilter} />
      ) : null}

      {/* Calendar placeholder */}
      {innerTab === 'calendar' ? (
        <EmptyState emoji="📅" title={t('events.tab_calendar')} description={t('common.loading')} />
      ) : isLoading ? (
        <View style={styles.skeletons}>
          {[1, 2, 3].map((i) => (
            <View key={i} style={styles.skeletonCard}>
              <Skeleton width="100%" height={140} radius={16} />
            </View>
          ))}
        </View>
      ) : isError ? (
        <ErrorState onRetry={refetch} />
      ) : events.length === 0 ? (
        <EmptyState emoji="🎾" title={innerTab === 'my' ? t('events.empty_my') : t('events.empty')} />
      ) : (
        <FlatList
          data={events}
          keyExtractor={(item) => item.id}
          renderItem={renderItem}
          onEndReached={() => { if (hasNextPage) fetchNextPage(); }}
          onEndReachedThreshold={0.3}
          ListFooterComponent={isFetchingNextPage ? <ActivityIndicator style={styles.footer} color={colors.primary} /> : null}
          refreshControl={<RefreshControl refreshing={false} onRefresh={refetch} tintColor={colors.primary} />}
          showsVerticalScrollIndicator={false}
          contentContainerStyle={styles.listContent}
        />
      )}

      {/* FAB */}
      <FAB onPress={() => router.push('/event/create')} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  tabBar: {
    flexDirection: 'row',
    borderBottomWidth: 1,
    borderBottomColor: colors.borderLight,
    backgroundColor: colors.card,
  },
  tab: { flex: 1, paddingVertical: spacing.md, alignItems: 'center', borderBottomWidth: 2, borderBottomColor: 'transparent' },
  tabActive: { borderBottomColor: colors.primary },
  tabText: { ...typography.textStyles.bodySm, color: colors.textMuted },
  tabTextActive: { color: colors.primary, fontWeight: typography.fontWeight.semibold },
  skeletons: { paddingTop: spacing.md },
  skeletonCard: { paddingHorizontal: spacing.base, marginBottom: spacing.md },
  listContent: { paddingTop: spacing.sm, paddingBottom: spacing['3xl'] },
  footer: { paddingVertical: spacing.xl },
});
