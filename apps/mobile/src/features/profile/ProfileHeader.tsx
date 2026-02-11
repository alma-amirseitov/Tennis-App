import { View, Text, StyleSheet } from 'react-native';
import { Avatar, Badge } from '@/shared/ui';
import { colors, spacing, typography } from '@/shared/theme';
import type { UserProfile } from '@/shared/api/users';

interface ProfileHeaderProps {
  user: UserProfile;
}

export function ProfileHeader({ user }: ProfileHeaderProps) {
  const ratingSign = user.rating_delta >= 0 ? '▲' : '▼';
  const ratingColor = user.rating_delta >= 0 ? colors.success : colors.danger;
  return (
    <View style={styles.container}>
      <Avatar uri={user.avatar_url} name={`${user.first_name} ${user.last_name}`} size="xl" />
      <Text style={styles.name}>{user.first_name} {user.last_name}</Text>
      <View style={styles.levelRow}>
        <Badge variant="primary" text={`NTRP ${user.ntrp_level}`} />
        <Text style={styles.levelLabel}>{user.level_label}</Text>
      </View>
      <View style={styles.ratingRow}>
        <Text style={styles.ratingValue}>{user.rating.toLocaleString()}</Text>
        <Text style={[styles.ratingDelta, { color: ratingColor }]}>
          {ratingSign} {Math.abs(user.rating_delta)}
        </Text>
      </View>
      {user.district ? (
        <Text style={styles.location}>📍 {user.city}, {user.district}</Text>
      ) : null}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { alignItems: 'center', paddingVertical: spacing.xl },
  name: { ...typography.textStyles.h2, color: colors.text, marginTop: spacing.md },
  levelRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm, marginTop: spacing.sm },
  levelLabel: { ...typography.textStyles.bodySm, color: colors.textSecondary },
  ratingRow: { flexDirection: 'row', alignItems: 'baseline', gap: spacing.sm, marginTop: spacing.sm },
  ratingValue: { ...typography.textStyles.h3, color: colors.text },
  ratingDelta: { ...typography.textStyles.bodySm, fontWeight: typography.fontWeight.semibold },
  location: { ...typography.textStyles.bodySm, color: colors.textMuted, marginTop: spacing.xs },
});
