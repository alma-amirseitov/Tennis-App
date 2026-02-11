import { View, Text, Pressable, StyleSheet } from 'react-native';
import { useRouter } from 'expo-router';
import { Avatar, Badge } from '@/shared/ui';
import { colors, spacing, typography, radius } from '@/shared/theme';
import type { UserProfile } from '@/shared/api/users';

interface PlayerCardProps {
  player: UserProfile;
}

export function PlayerCard({ player }: PlayerCardProps) {
  const router = useRouter();
  return (
    <Pressable
      onPress={() => router.push(`/player/${player.id}`)}
      style={({ pressed }) => [styles.card, pressed && { transform: [{ scale: 0.98 }] }]}
    >
      <Avatar uri={player.avatar_url} name={`${player.first_name} ${player.last_name}`} size="lg" />
      <View style={styles.info}>
        <View style={styles.nameRow}>
          <Text style={styles.name} numberOfLines={1}>{player.first_name} {player.last_name}</Text>
          <Badge variant="primary" text={`NTRP ${player.ntrp_level}`} />
        </View>
        {player.district ? (
          <Text style={styles.location}>📍 {player.district}</Text>
        ) : null}
        <View style={styles.statsRow}>
          <Text style={styles.stat}>🎾 {player.total_games}</Text>
          <Text style={styles.stat}>🏆 {player.wins}</Text>
          <Text style={styles.stat}>📊 {player.win_rate}%</Text>
        </View>
      </View>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  card: {
    flexDirection: 'row',
    backgroundColor: colors.card,
    borderRadius: radius.lg,
    borderWidth: 1,
    borderColor: colors.borderLight,
    padding: spacing.base,
    marginHorizontal: spacing.base,
    marginBottom: spacing.md,
    gap: spacing.md,
  },
  info: { flex: 1 },
  nameRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', gap: spacing.sm },
  name: { ...typography.textStyles.h4, color: colors.text, flex: 1 },
  location: { ...typography.textStyles.caption, color: colors.textMuted, marginTop: 2 },
  statsRow: { flexDirection: 'row', gap: spacing.md, marginTop: spacing.sm },
  stat: { ...typography.textStyles.caption, color: colors.textSecondary },
});
