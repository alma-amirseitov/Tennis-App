import { useState, useCallback, useMemo } from 'react';
import { View, TextInput, FlatList, ScrollView, StyleSheet, ActivityIndicator, RefreshControl } from 'react-native';
import { useTranslation } from 'react-i18next';
import { useCommunities } from '@/shared/api/hooks';
import { CommunityCard } from '@/features/communities/CommunityCard';
import { EmptyState, ErrorState, Skeleton, Chip } from '@/shared/ui';
import { colors, spacing, typography, radius } from '@/shared/theme';
import type { CommunityItem, CommunityType } from '@/shared/api/communities';

const TYPE_FILTERS: (CommunityType | 'all')[] = ['all', 'club', 'league', 'organizer', 'group'];

export default function CommunitiesScreen() {
  const { t } = useTranslation();
  const [query, setQuery] = useState('');
  const [typeFilter, setTypeFilter] = useState<CommunityType | 'all'>('all');

  const typeLabels: Record<string, string> = {
    all: t('common.all'),
    club: t('communities.type_club'),
    league: t('communities.type_league'),
    organizer: t('communities.type_organizer'),
    group: t('communities.type_group'),
  };

  const params = useMemo(() => ({
    query: query.length >= 2 ? query : undefined,
    community_type: typeFilter === 'all' ? undefined : typeFilter,
  }), [query, typeFilter]);

  const { data, isLoading, isError, refetch, fetchNextPage, hasNextPage, isFetchingNextPage } = useCommunities(params);

  const communities = useMemo(() => data?.pages.flatMap((p) => p.data) ?? [], [data]);

  const renderItem = useCallback(({ item }: { item: CommunityItem }) => <CommunityCard community={item} />, []);

  return (
    <View style={styles.container}>
      <View style={styles.searchContainer}>
        <TextInput
          style={styles.searchInput}
          placeholder={t('common.search_placeholder')}
          placeholderTextColor={colors.textMuted}
          value={query}
          onChangeText={setQuery}
          autoCorrect={false}
        />
      </View>

      {/* Type filters */}
      <ScrollView horizontal showsHorizontalScrollIndicator={false} contentContainerStyle={styles.filters}>
        {TYPE_FILTERS.map((type) => (
          <Chip key={type} label={typeLabels[type]} selected={typeFilter === type} onPress={() => setTypeFilter(type)} />
        ))}
      </ScrollView>

      {isLoading ? (
        <View style={styles.skeletons}>
          {[1, 2, 3].map((i) => (
            <View key={i} style={styles.skeletonCard}>
              <Skeleton width="100%" height={90} radius={16} />
            </View>
          ))}
        </View>
      ) : isError ? (
        <ErrorState onRetry={refetch} />
      ) : communities.length === 0 ? (
        <EmptyState emoji="🏛" title={t('communities.empty')} />
      ) : (
        <FlatList
          data={communities}
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
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  searchContainer: { paddingHorizontal: spacing.base, paddingTop: spacing.md },
  searchInput: {
    height: 44,
    borderRadius: radius.md,
    backgroundColor: colors.card,
    borderWidth: 1,
    borderColor: colors.border,
    paddingHorizontal: spacing.base,
    ...typography.textStyles.body,
    color: colors.text,
  },
  filters: { paddingHorizontal: spacing.base, paddingVertical: spacing.md },
  skeletons: { paddingTop: spacing.md },
  skeletonCard: { paddingHorizontal: spacing.base, marginBottom: spacing.md },
  listContent: { paddingTop: spacing.sm, paddingBottom: spacing['3xl'] },
  footer: { paddingVertical: spacing.xl },
});
