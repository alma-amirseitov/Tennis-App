import { useState, useCallback, useMemo } from 'react';
import { View, TextInput, FlatList, StyleSheet, ActivityIndicator, RefreshControl } from 'react-native';
import { useTranslation } from 'react-i18next';
import { useSearchPlayers } from '@/shared/api/hooks';
import { PlayerCard } from '@/features/players/PlayerCard';
import { PlayerFilters } from '@/features/players/PlayerFilters';
import { EmptyState, ErrorState, Skeleton } from '@/shared/ui';
import { colors, spacing, typography, radius } from '@/shared/theme';
import type { UserProfile } from '@/shared/api/users';

export default function PlayersScreen() {
  const { t } = useTranslation();
  const [query, setQuery] = useState('');
  const [gender, setGender] = useState('all');

  const params = useMemo(() => ({
    query: query.length >= 2 ? query : undefined,
    gender: gender === 'all' ? undefined : gender,
  }), [query, gender]);

  const {
    data,
    isLoading,
    isError,
    refetch,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useSearchPlayers(params);

  const players = useMemo(
    () => data?.pages.flatMap((p) => p.data) ?? [],
    [data]
  );

  const renderItem = useCallback(
    ({ item }: { item: UserProfile }) => <PlayerCard player={item} />,
    []
  );

  const renderFooter = () => {
    if (isFetchingNextPage) {
      return <ActivityIndicator style={styles.footer} color={colors.primary} />;
    }
    return null;
  };

  return (
    <View style={styles.container}>
      {/* Search */}
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

      {/* Filters */}
      <PlayerFilters gender={gender} onGenderChange={setGender} />

      {/* List */}
      {isLoading ? (
        <View style={styles.skeletons}>
          {[1, 2, 3, 4].map((i) => (
            <View key={i} style={styles.skeletonCard}>
              <Skeleton width="100%" height={90} radius={16} />
            </View>
          ))}
        </View>
      ) : isError ? (
        <ErrorState onRetry={refetch} />
      ) : players.length === 0 ? (
        <EmptyState emoji="👥" title={t('players.empty')} />
      ) : (
        <FlatList
          data={players}
          keyExtractor={(item) => item.id}
          renderItem={renderItem}
          onEndReached={() => { if (hasNextPage) fetchNextPage(); }}
          onEndReachedThreshold={0.3}
          ListFooterComponent={renderFooter}
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
  skeletons: { paddingTop: spacing.md },
  skeletonCard: { paddingHorizontal: spacing.base, marginBottom: spacing.md },
  listContent: { paddingTop: spacing.sm, paddingBottom: spacing['3xl'] },
  footer: { paddingVertical: spacing.xl },
});
