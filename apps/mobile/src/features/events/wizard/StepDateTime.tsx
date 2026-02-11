import { View, Text, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { Button, Input } from '@/shared/ui';
import { colors, spacing, typography } from '@/shared/theme';
import { useWizard } from './WizardContext';

export function StepDateTime() {
  const { t } = useTranslation();
  const { data, update, setStep } = useWizard();
  return (
    <View style={styles.container}>
      <Text style={styles.title}>{t('events.wizard_step_6')}</Text>
      <Input
        label={t('events.filter_date')}
        placeholder="2026-03-15"
        value={data.start_date}
        onChangeText={(v) => update({ start_date: v })}
      />
      <View style={styles.spacer} />
      <Input
        label={t('events.detail_time') + ' (start)'}
        placeholder="18:00"
        value={data.start_time}
        onChangeText={(v) => update({ start_time: v })}
      />
      <View style={styles.spacer} />
      <Input
        label={t('events.detail_time') + ' (end)'}
        placeholder="19:30"
        value={data.end_time}
        onChangeText={(v) => update({ end_time: v })}
      />
      <Button
        variant="primary"
        title={t('auth.quiz_next')}
        onPress={() => setStep(6)}
        disabled={!data.start_date || !data.start_time}
        style={styles.nextBtn}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { padding: spacing.xl },
  title: { ...typography.textStyles.h2, color: colors.text, marginBottom: spacing.xl },
  spacer: { height: spacing.lg },
  nextBtn: { marginTop: spacing.xl },
});
