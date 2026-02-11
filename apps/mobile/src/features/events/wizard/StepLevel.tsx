import { View, Text, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { Button, Select } from '@/shared/ui';
import { colors, spacing, typography } from '@/shared/theme';
import { useWizard } from './WizardContext';

const LEVELS = ['1.0', '1.5', '2.0', '2.5', '3.0', '3.5', '4.0', '4.5', '5.0', '5.5', '6.0'];
const LEVEL_OPTIONS = LEVELS.map((l) => ({ value: l, label: `NTRP ${l}` }));

export function StepLevel() {
  const { t } = useTranslation();
  const { data, update, setStep } = useWizard();
  return (
    <View style={styles.container}>
      <Text style={styles.title}>{t('events.wizard_step_3')}</Text>
      <Select
        label={t('events.filter_level') + ' (min)'}
        value={String(data.level_min)}
        onSelect={(v) => update({ level_min: parseFloat(v) })}
        options={LEVEL_OPTIONS}
      />
      <View style={styles.spacer} />
      <Select
        label={t('events.filter_level') + ' (max)'}
        value={String(data.level_max)}
        onSelect={(v) => update({ level_max: parseFloat(v) })}
        options={LEVEL_OPTIONS}
      />
      <Button variant="primary" title={t('auth.quiz_next')} onPress={() => setStep(3)} style={styles.nextBtn} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { padding: spacing.xl },
  title: { ...typography.textStyles.h2, color: colors.text, marginBottom: spacing.xl },
  spacer: { height: spacing.lg },
  nextBtn: { marginTop: spacing.xl },
});
