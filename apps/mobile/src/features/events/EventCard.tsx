import { View, Text, Pressable, StyleSheet } from 'react-native';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { Badge } from '@/shared/ui';
import { colors, spacing, typography, radius } from '@/shared/theme';
import type { EventItem, EventStatus, EventType } from '@/shared/api/events';

const STATUS_COLORS: Record<EventStatus, string> = {
  draft: colors.textMuted,
  open: colors.statusOpen,
  filling: colors.statusFilling,
  full: colors.statusFull,
  in_progress: colors.info,
  completed: colors.statusDone,
  cancelled: colors.statusCancel,
};

const TYPE_EMOJI: Record<EventType, string> = {
  find_partner: '🔍',
  organized_game: '🎾',
  tournament: '🏆',
  training: '🏋️',
};

interface EventCardProps {
  event: EventItem;
}

export function EventCard({ event }: EventCardProps) {
  const { t } = useTranslation();
  const router = useRouter();

  const statusColor = STATUS_COLORS[event.status] ?? colors.textMuted;
  const statusKey = `events.status_${event.status}` as const;
  const typeEmoji = TYPE_EMOJI[event.event_type] ?? '🎾';

  const dateStr = new Date(event.start_time).toLocaleDateString('ru-RU', {
    day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit',
  });

  return (
    <Pressable
      onPress={() => router.push(`/event/${event.id}`)}
      style={({ pressed }) => [styles.card, pressed && { transform: [{ scale: 0.98 }] }]}
    >
      {/* Status bar */}
      <View style={styles.topRow}>
        <View style={[styles.statusDot, { backgroundColor: statusColor }]} />
        <Text style={[styles.statusText, { color: statusColor }]}>{t(statusKey)}</Text>
        <Text style={styles.composition}>{event.composition_type}</Text>
      </View>

      {/* Title */}
      <Text style={styles.title}>{typeEmoji} {event.title}</Text>

      {/* Info rows */}
      {event.creator ? (
        <Text style={styles.info}>👤 {event.creator.first_name} ({event.creator.ntrp_level})</Text>
      ) : null}
      {event.location_name ? <Text style={styles.info}>📍 {event.location_name}</Text> : null}
      <Text style={styles.info}>📅 {dateStr}</Text>
      <Text style={styles.info}>👥 {t('events.spots', { current: event.current_participants, max: event.max_participants })}</Text>

      {/* Level badge */}
      <View style={styles.bottomRow}>
        <Badge variant="primary" text={t('events.detail_level', { min: event.level_min, max: event.level_max })} />
        {event.price > 0 ? <Text style={styles.price}>💰 {event.price.toLocaleString()} ₸</Text> : null}
      </View>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.lg,
    borderWidth: 1,
    borderColor: colors.borderLight,
    padding: spacing.base,
    marginHorizontal: spacing.base,
    marginBottom: spacing.md,
  },
  topRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm, marginBottom: spacing.sm },
  statusDot: { width: 8, height: 8, borderRadius: 4 },
  statusText: { ...typography.textStyles.caption, fontWeight: typography.fontWeight.semibold, flex: 1 },
  composition: { ...typography.textStyles.caption, color: colors.textMuted },
  title: { ...typography.textStyles.h4, color: colors.text, marginBottom: spacing.sm },
  info: { ...typography.textStyles.bodySm, color: colors.textSecondary, marginBottom: 2 },
  bottomRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', marginTop: spacing.md },
  price: { ...typography.textStyles.bodySm, color: colors.accent, fontWeight: typography.fontWeight.semibold },
});
