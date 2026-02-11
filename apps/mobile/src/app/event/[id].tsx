import { View, Text, ScrollView, StyleSheet, RefreshControl } from 'react-native';
import { useLocalSearchParams } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { useEvent, useJoinEvent, useLeaveEvent } from '@/shared/api/hooks';
import { Avatar, Badge, Button, ScreenHeader, Skeleton, ErrorState } from '@/shared/ui';
import { colors, spacing, typography, radius } from '@/shared/theme';
import type { EventStatus } from '@/shared/api/events';

const STATUS_COLORS: Record<EventStatus, string> = {
  draft: colors.textMuted, open: colors.statusOpen, filling: colors.statusFilling,
  full: colors.statusFull, in_progress: colors.info, completed: colors.statusDone, cancelled: colors.statusCancel,
};

export default function EventDetailScreen() {
  const { t } = useTranslation();
  const { id } = useLocalSearchParams<{ id: string }>();
  const { data: event, isLoading, isError, refetch } = useEvent(id ?? '');
  const joinMutation = useJoinEvent();
  const leaveMutation = useLeaveEvent();

  if (isLoading) {
    return (
      <View style={styles.container}>
        <ScreenHeader title="" showBack />
        <View style={styles.skeleton}><Skeleton width="100%" height={200} radius={16} /></View>
      </View>
    );
  }

  if (isError || !event) {
    return (
      <View style={styles.container}>
        <ScreenHeader title="" showBack />
        <ErrorState onRetry={refetch} />
      </View>
    );
  }

  const statusColor = STATUS_COLORS[event.status] ?? colors.textMuted;
  const isFull = event.current_participants >= event.max_participants;
  const dateStr = new Date(event.start_time).toLocaleDateString('ru-RU', {
    weekday: 'long', day: 'numeric', month: 'long', hour: '2-digit', minute: '2-digit',
  });

  const getActionButton = () => {
    if (event.status === 'completed' || event.status === 'cancelled') return null;
    if (event.is_joined) {
      return (
        <Button variant="outline" title={t('events.joined') + ' ✓'} onPress={() => leaveMutation.mutate(event.id)} loading={leaveMutation.isPending} />
      );
    }
    if (isFull) {
      return <Button variant="outline" title={t('events.full')} onPress={() => {}} disabled />;
    }
    return (
      <Button
        variant="primary"
        title={event.price > 0 ? `${t('events.join')} (${event.price.toLocaleString()} ₸)` : t('events.join')}
        onPress={() => joinMutation.mutate(event.id)}
        loading={joinMutation.isPending}
      />
    );
  };

  return (
    <View style={styles.container}>
      <ScreenHeader title={t('events.title')} showBack />
      <ScrollView
        contentContainerStyle={styles.content}
        refreshControl={<RefreshControl refreshing={false} onRefresh={refetch} tintColor={colors.primary} />}
        showsVerticalScrollIndicator={false}
      >
        {/* Status */}
        <View style={styles.statusRow}>
          <View style={[styles.statusDot, { backgroundColor: statusColor }]} />
          <Badge variant="primary" text={t(`events.status_${event.status}`)} />
        </View>

        {/* Title */}
        <Text style={styles.title}>{event.title}</Text>

        {/* Info */}
        <View style={styles.infoSection}>
          <InfoRow emoji="📅" text={dateStr} />
          <InfoRow emoji="📍" text={event.location_name} />
          <InfoRow emoji="🎾" text={`${event.composition_type}, ${t('events.detail_sets', { count: event.sets_count })}`} />
          <InfoRow emoji="📊" text={t('events.detail_level', { min: event.level_min, max: event.level_max })} />
          {event.price > 0 ? <InfoRow emoji="💰" text={`${event.price.toLocaleString()} ₸`} /> : null}
          {event.creator ? <InfoRow emoji="👤" text={`${event.creator.first_name} ${event.creator.last_name}`} /> : null}
        </View>

        {/* Participants */}
        <Text style={styles.sectionTitle}>
          {t('events.detail_participants')} ({event.current_participants}/{event.max_participants})
        </Text>
        <View style={styles.participants}>
          {event.participants.map((p) => (
            <View key={p.id} style={styles.participant}>
              <Avatar uri={p.avatar_url} name={`${p.first_name} ${p.last_name}`} size="sm" />
              <Text style={styles.participantName}>{p.first_name}</Text>
            </View>
          ))}
          {event.current_participants < event.max_participants ? (
            <View style={styles.emptySlots}>
              <Text style={styles.emptySlotsText}>
                +{event.max_participants - event.current_participants} {t('events.spots_left', { count: event.max_participants - event.current_participants })}
              </Text>
            </View>
          ) : null}
        </View>

        {/* Description */}
        {event.description ? (
          <>
            <Text style={styles.sectionTitle}>{t('profile.about')}</Text>
            <Text style={styles.description}>{event.description}</Text>
          </>
        ) : null}

        {/* Action */}
        <View style={styles.actionSection}>
          {getActionButton()}
        </View>
      </ScrollView>
    </View>
  );
}

function InfoRow({ emoji, text }: { emoji: string; text: string }) {
  return (
    <View style={styles.infoRow}>
      <Text style={styles.infoEmoji}>{emoji}</Text>
      <Text style={styles.infoText}>{text}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  content: { padding: spacing.base, paddingBottom: spacing['3xl'] },
  skeleton: { padding: spacing.base },
  statusRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm, marginBottom: spacing.md },
  statusDot: { width: 10, height: 10, borderRadius: 5 },
  title: { ...typography.textStyles.h2, color: colors.text, marginBottom: spacing.lg },
  infoSection: { marginBottom: spacing.lg },
  infoRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm, marginBottom: spacing.sm },
  infoEmoji: { fontSize: 16, width: 24, textAlign: 'center' },
  infoText: { ...typography.textStyles.body, color: colors.textSecondary, flex: 1 },
  sectionTitle: { ...typography.textStyles.h4, color: colors.text, marginBottom: spacing.md, marginTop: spacing.lg },
  participants: { flexDirection: 'row', flexWrap: 'wrap', gap: spacing.md },
  participant: { alignItems: 'center', width: 56 },
  participantName: { ...typography.textStyles.caption, color: colors.text, marginTop: 2 },
  emptySlots: { justifyContent: 'center' },
  emptySlotsText: { ...typography.textStyles.caption, color: colors.textMuted },
  description: { ...typography.textStyles.body, color: colors.textSecondary },
  actionSection: { marginTop: spacing.xl },
});
