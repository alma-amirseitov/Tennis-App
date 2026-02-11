import { View, Text, ScrollView, StyleSheet, RefreshControl } from 'react-native';
import { useLocalSearchParams } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { usePlayer } from '@/shared/api/hooks';
import { ProfileHeader } from '@/features/profile/ProfileHeader';
import { StatsCard } from '@/features/profile/StatsCard';
import { Button, ScreenHeader, Skeleton, ErrorState } from '@/shared/ui';
import { colors, spacing, typography, radius } from '@/shared/theme';

export default function PlayerDetailScreen() {
  const { t } = useTranslation();
  const { id } = useLocalSearchParams<{ id: string }>();
  const { data: player, isLoading, isError, refetch } = usePlayer(id ?? '');

  if (isLoading) {
    return (
      <View style={styles.container}>
        <ScreenHeader title={t('players.title')} showBack />
        <View style={styles.skeletonContent}>
          <Skeleton width={80} height={80} radius={40} />
          <Skeleton width={180} height={22} radius={4} />
          <Skeleton width="100%" height={80} radius={16} />
        </View>
      </View>
    );
  }

  if (isError || !player) {
    return (
      <View style={styles.container}>
        <ScreenHeader title={t('players.title')} showBack />
        <ErrorState onRetry={refetch} />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <ScreenHeader title={`${player.first_name} ${player.last_name}`} showBack />
      <ScrollView
        contentContainerStyle={styles.content}
        refreshControl={<RefreshControl refreshing={false} onRefresh={refetch} tintColor={colors.primary} />}
        showsVerticalScrollIndicator={false}
      >
        <ProfileHeader user={player} />
        <StatsCard totalGames={player.total_games} wins={player.wins} winRate={player.win_rate} />

        {/* Action buttons */}
        <View style={styles.actions}>
          <Button variant="outline" title={t('players.write_message')} onPress={() => {}} style={styles.actionBtn} />
          <Button variant="outline" title={t('players.invite_to_game')} onPress={() => {}} style={styles.actionBtn} />
          <Button variant="primary" title={t('players.add_friend')} onPress={() => {}} style={styles.actionBtn} />
        </View>

        {/* Achievements placeholder */}
        <Text style={styles.sectionTitle}>{t('profile.badges')}</Text>
        <View style={styles.placeholder}>
          <Text style={styles.placeholderEmoji}>🏅 🎖 🏆 🔥</Text>
        </View>

        {/* Match history placeholder */}
        <Text style={styles.sectionTitle}>{t('profile.match_history')}</Text>
        <View style={styles.placeholder}>
          <Text style={styles.placeholderEmoji}>🎾</Text>
          <Text style={styles.placeholderText}>{t('common.loading')}</Text>
        </View>
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  content: { paddingBottom: spacing['3xl'] },
  skeletonContent: { alignItems: 'center', gap: spacing.lg, paddingTop: spacing['3xl'] },
  actions: { flexDirection: 'row', gap: spacing.sm, paddingHorizontal: spacing.base, marginTop: spacing.lg },
  actionBtn: { flex: 1 },
  sectionTitle: {
    ...typography.textStyles.h4,
    color: colors.text,
    paddingHorizontal: spacing.base,
    marginTop: spacing.xl,
    marginBottom: spacing.md,
  },
  placeholder: {
    alignItems: 'center',
    paddingVertical: spacing.xl,
    marginHorizontal: spacing.base,
    backgroundColor: colors.card,
    borderRadius: radius.lg,
    borderWidth: 1,
    borderColor: colors.borderLight,
  },
  placeholderEmoji: { fontSize: 24, marginBottom: spacing.xs },
  placeholderText: { ...typography.textStyles.caption, color: colors.textMuted },
});
