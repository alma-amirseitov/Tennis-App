import { View, Text, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { Button, Select } from '@/shared/ui';
import { colors, spacing, typography } from '@/shared/theme';
import { useWizard } from './WizardContext';

const PARTICIPANTS = [2, 4, 8, 16, 32].map((n) => ({ value: String(n), label: String(n) }));
const SETS = [1, 3, 5].map((n) => ({ value: String(n), label: String(n) }));

export function StepDetails() {
  const { t } = useTranslation();
  const { data, update, setStep } = useWizard();
  return (
    <View style={styles.container}>
      <Text style={styles.title}>{t('events.wizard_step_4')}</Text>
      <Select
        label={t('events.detail_participants')}
        value={String(data.max_participants)}
        onSelect={(v) => update({ max_participants: parseInt(v, 10) })}
        options={PARTICIPANTS}
      />
      <View style={styles.spacer} />
      <Select
        label={t('events.detail_sets', { count: data.sets_count })}
        value={String(data.sets_count)}
        onSelect={(v) => update({ sets_count: parseInt(v, 10) })}
        options={SETS}
      />
      <Button variant="primary" title={t('auth.quiz_next')} onPress={() => setStep(4)} style={styles.nextBtn} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { padding: spacing.xl },
  title: { ...typography.textStyles.h2, color: colors.text, marginBottom: spacing.xl },
  spacer: { height: spacing.lg },
  nextBtn: { marginTop: spacing.xl },
});
