import { View, Text, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { colors, spacing, typography, radius } from '@/shared/theme';

interface StatsCardProps {
  totalGames: number;
  wins: number;
  winRate: number;
}

export function StatsCard({ totalGames, wins, winRate }: StatsCardProps) {
  const { t } = useTranslation();
  return (
    <View style={styles.container}>
      <StatItem emoji="🎾" value={String(totalGames)} label={t('profile.matches')} />
      <View style={styles.divider} />
      <StatItem emoji="🏆" value={String(wins)} label={t('profile.wins')} />
      <View style={styles.divider} />
      <StatItem emoji="📊" value={`${winRate}%`} label={t('profile.win_rate')} />
    </View>
  );
}

function StatItem({ emoji, value, label }: { emoji: string; value: string; label: string }) {
  return (
    <View style={styles.item}>
      <Text style={styles.emoji}>{emoji}</Text>
      <Text style={styles.value}>{value}</Text>
      <Text style={styles.label}>{label}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    backgroundColor: colors.card,
    borderRadius: radius.lg,
    borderWidth: 1,
    borderColor: colors.borderLight,
    paddingVertical: spacing.base,
    marginHorizontal: spacing.base,
  },
  item: { flex: 1, alignItems: 'center' },
  divider: { width: 1, backgroundColor: colors.borderLight },
  emoji: { fontSize: 20, marginBottom: spacing.xs },
  value: { ...typography.textStyles.h3, color: colors.text },
  label: { ...typography.textStyles.caption, color: colors.textMuted, marginTop: 2 },
});
