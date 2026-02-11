import { ScrollView, View, Text, Pressable, StyleSheet, RefreshControl } from 'react-native';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { useProfile } from '@/shared/api/hooks';
import { ProfileHeader } from '@/features/profile/ProfileHeader';
import { StatsCard } from '@/features/profile/StatsCard';
import { CommunitiesList } from '@/features/profile/CommunitiesList';
import { Button, Skeleton, ErrorState } from '@/shared/ui';
import { colors, spacing, typography, radius } from '@/shared/theme';

export default function ProfileScreen() {
  const { t } = useTranslation();
  const router = useRouter();
  const { data: user, isLoading, isError, refetch } = useProfile();

  if (isLoading) {
    return (
      <View style={styles.container}>
        <View style={styles.skeletonHeader}>
          <Skeleton width={80} height={80} radius={40} />
          <Skeleton width={180} height={22} radius={4} />
          <Skeleton width={120} height={16} radius={4} />
        </View>
        <View style={styles.skeletonStats}>
          <Skeleton width="100%" height={80} radius={16} />
        </View>
      </View>
    );
  }

  if (isError || !user) {
    return (
      <View style={styles.container}>
        <ErrorState onRetry={refetch} />
      </View>
    );
  }

  return (
    <ScrollView
      style={styles.container}
      contentContainerStyle={styles.content}
      refreshControl={<RefreshControl refreshing={false} onRefresh={refetch} tintColor={colors.primary} />}
      showsVerticalScrollIndicator={false}
    >
      {/* Header with settings icon */}
      <View style={styles.topBar}>
        <Text style={styles.screenTitle}>{t('profile.title')}</Text>
        <Pressable onPress={() => router.push('/(tabs)/profile/settings')} hitSlop={12}>
          <Text style={styles.settingsIcon}>⚙️</Text>
        </Pressable>
      </View>

      <ProfileHeader user={user} />

      {/* Quick Stats */}
      <StatsCard totalGames={user.total_games} wins={user.wins} winRate={user.win_rate} />

      {/* Edit button */}
      <View style={styles.editSection}>
        <Button
          variant="outline"
          title={t('profile.edit')}
          onPress={() => router.push('/(tabs)/profile/edit')}
        />
      </View>

      {/* Communities */}
      <SectionTitle title={t('profile.communities')} />
      <CommunitiesList communities={[]} />

      {/* Badges placeholder */}
      <SectionTitle title={t('profile.badges')} />
      <View style={styles.placeholderSection}>
        <Text style={styles.placeholderText}>🏅 🎖 🏆 🔥 ⭐ 👑</Text>
        <Text style={styles.placeholderDesc}>{t('common.loading')}</Text>
      </View>

      {/* Friends placeholder */}
      <SectionTitle title={t('profile.friends')} />
      <View style={styles.placeholderSection}>
        <Text style={styles.placeholderText}>👥</Text>
        <Text style={styles.placeholderDesc}>{t('common.loading')}</Text>
      </View>

      {/* Match history placeholder */}
      <SectionTitle title={t('profile.match_history')} />
      <View style={styles.placeholderSection}>
        <Text style={styles.placeholderText}>🎾</Text>
        <Text style={styles.placeholderDesc}>{t('common.loading')}</Text>
      </View>
    </ScrollView>
  );
}

function SectionTitle({ title }: { title: string }) {
  return <Text style={styles.sectionTitle}>{title}</Text>;
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  content: { paddingBottom: spacing['3xl'] },
  topBar: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: spacing.base,
    paddingVertical: spacing.md,
  },
  screenTitle: { ...typography.textStyles.h2, color: colors.text },
  settingsIcon: { fontSize: 22 },
  editSection: { paddingHorizontal: spacing.base, marginTop: spacing.lg },
  sectionTitle: {
    ...typography.textStyles.h4,
    color: colors.text,
    paddingHorizontal: spacing.base,
    marginTop: spacing.xl,
    marginBottom: spacing.md,
  },
  placeholderSection: {
    alignItems: 'center',
    paddingVertical: spacing.xl,
    marginHorizontal: spacing.base,
    backgroundColor: colors.card,
    borderRadius: radius.lg,
    borderWidth: 1,
    borderColor: colors.borderLight,
  },
  placeholderText: { fontSize: 24, marginBottom: spacing.xs },
  placeholderDesc: { ...typography.textStyles.caption, color: colors.textMuted },
  skeletonHeader: { alignItems: 'center', gap: spacing.md, paddingTop: spacing['3xl'] },
  skeletonStats: { padding: spacing.base, marginTop: spacing.xl },
});
