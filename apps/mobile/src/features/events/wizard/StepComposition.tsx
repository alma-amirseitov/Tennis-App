import { View, Text, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { Card, Button } from '@/shared/ui';
import { colors, spacing, typography } from '@/shared/theme';
import { useWizard } from './WizardContext';
import type { CompositionType } from '@/shared/api/events';

const COMPOSITIONS: { key: CompositionType; labelKey: string }[] = [
  { key: 'singles', labelKey: 'events.composition_singles' },
  { key: 'doubles', labelKey: 'events.composition_doubles' },
  { key: 'mixed', labelKey: 'events.composition_mixed' },
  { key: 'team', labelKey: 'events.composition_team' },
];

export function StepComposition() {
  const { t } = useTranslation();
  const { data, update, setStep } = useWizard();
  return (
    <View style={styles.container}>
      <Text style={styles.title}>{t('events.wizard_step_2')}</Text>
      {COMPOSITIONS.map((c) => (
        <Card
          key={c.key}
          onPress={() => update({ composition_type: c.key })}
          style={[styles.card, data.composition_type === c.key && styles.cardSelected]}
        >
          <Text style={[styles.label, data.composition_type === c.key && styles.labelSelected]}>
            {t(c.labelKey)}
          </Text>
        </Card>
      ))}
      <Button
        variant="primary"
        title={t('auth.quiz_next')}
        onPress={() => setStep(2)}
        disabled={!data.composition_type}
        style={styles.nextBtn}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { padding: spacing.xl },
  title: { ...typography.textStyles.h2, color: colors.text, marginBottom: spacing.xl },
  card: { marginBottom: spacing.md, borderWidth: 2, borderColor: colors.border },
  cardSelected: { borderColor: colors.primary, backgroundColor: colors.primaryLight },
  label: { ...typography.textStyles.body, color: colors.text, textAlign: 'center' },
  labelSelected: { color: colors.primary, fontWeight: typography.fontWeight.semibold },
  nextBtn: { marginTop: spacing.lg },
});
