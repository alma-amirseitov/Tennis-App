import { View, Text, ScrollView, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { useRouter } from 'expo-router';
import { Button, Card } from '@/shared/ui';
import { useCreateEvent } from '@/shared/api/hooks';
import { showToast } from '@/shared/lib/toast';
import { colors, spacing, typography } from '@/shared/theme';
import { useWizard } from './WizardContext';

export function StepReview() {
  const { t } = useTranslation();
  const router = useRouter();
  const { data } = useWizard();
  const createMutation = useCreateEvent();

  const handleCreate = async () => {
    if (!data.event_type || !data.composition_type) return;
    try {
      const startTime = `${data.start_date}T${data.start_time}:00`;
      const endTime = data.end_time ? `${data.start_date}T${data.end_time}:00` : undefined;
      await createMutation.mutateAsync({
        title: data.title,
        description: data.description,
        event_type: data.event_type as 'find_partner' | 'organized_game' | 'tournament' | 'training',
        composition_type: data.composition_type as 'singles' | 'doubles' | 'mixed' | 'team',
        level_min: data.level_min,
        level_max: data.level_max,
        max_participants: data.max_participants,
        sets_count: data.sets_count,
        start_time: startTime,
        end_time: endTime,
        location_name: data.location_name,
        price: data.price,
        community_id: data.community_id || undefined,
      });
      showToast(t('common.done'));
      router.back();
    } catch {
      showToast(t('errors.something_went_wrong'));
    }
  };

  return (
    <ScrollView contentContainerStyle={styles.container}>
      <Text style={styles.title}>{t('events.wizard_step_8')}</Text>
      <Card style={styles.reviewCard}>
        <ReviewRow label={t('events.filter_type')} value={data.event_type} />
        <ReviewRow label={t('events.wizard_step_2')} value={data.composition_type} />
        <ReviewRow label={t('events.filter_level')} value={`${data.level_min} — ${data.level_max}`} />
        <ReviewRow label={t('events.detail_participants')} value={String(data.max_participants)} />
        <ReviewRow label={t('events.detail_location')} value={data.location_name} />
        <ReviewRow label={t('events.detail_time')} value={`${data.start_date} ${data.start_time}${data.end_time ? ` — ${data.end_time}` : ''}`} />
        <ReviewRow label={t('events.create')} value={data.title} />
      </Card>
      <Button
        variant="primary"
        title={t('events.create')}
        onPress={handleCreate}
        loading={createMutation.isPending}
        style={styles.createBtn}
      />
    </ScrollView>
  );
}

function ReviewRow({ label, value }: { label: string; value: string }) {
  return (
    <View style={styles.row}>
      <Text style={styles.rowLabel}>{label}</Text>
      <Text style={styles.rowValue}>{value || '—'}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { padding: spacing.xl, paddingBottom: spacing['3xl'] },
  title: { ...typography.textStyles.h2, color: colors.text, marginBottom: spacing.xl },
  reviewCard: { marginBottom: spacing.xl },
  row: { flexDirection: 'row', justifyContent: 'space-between', paddingVertical: spacing.sm, borderBottomWidth: StyleSheet.hairlineWidth, borderBottomColor: colors.borderLight },
  rowLabel: { ...typography.textStyles.bodySm, color: colors.textMuted },
  rowValue: { ...typography.textStyles.bodySm, color: colors.text, fontWeight: typography.fontWeight.semibold, maxWidth: '60%', textAlign: 'right' },
  createBtn: { marginTop: spacing.lg },
});
