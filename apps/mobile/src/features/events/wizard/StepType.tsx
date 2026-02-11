import { View, Text, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { Card } from '@/shared/ui';
import { colors, spacing, typography } from '@/shared/theme';
import { useWizard } from './WizardContext';
import type { EventType } from '@/shared/api/events';

const TYPES: { key: EventType; emoji: string; labelKey: string; descKey: string }[] = [
  { key: 'find_partner', emoji: '🔍', labelKey: 'events.type_find_partner', descKey: 'events.type_find_partner' },
  { key: 'organized_game', emoji: '🎾', labelKey: 'events.type_organized_game', descKey: 'events.type_organized_game' },
  { key: 'tournament', emoji: '🏆', labelKey: 'events.type_tournament', descKey: 'events.type_tournament' },
  { key: 'training', emoji: '🏋️', labelKey: 'events.type_training', descKey: 'events.type_training' },
];

export function StepType() {
  const { t } = useTranslation();
  const { data, update, setStep } = useWizard();
  return (
    <View style={styles.container}>
      <Text style={styles.title}>{t('events.wizard_step_1')}</Text>
      {TYPES.map((tp) => (
        <Card
          key={tp.key}
          onPress={() => { update({ event_type: tp.key }); setStep(1); }}
          style={[styles.card, data.event_type === tp.key && styles.cardSelected]}
        >
          <Text style={styles.emoji}>{tp.emoji}</Text>
          <Text style={styles.label}>{t(tp.labelKey)}</Text>
        </Card>
      ))}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { padding: spacing.xl },
  title: { ...typography.textStyles.h2, color: colors.text, marginBottom: spacing.xl },
  card: { flexDirection: 'row', alignItems: 'center', gap: spacing.md, marginBottom: spacing.md, borderWidth: 2, borderColor: colors.border },
  cardSelected: { borderColor: colors.primary, backgroundColor: colors.primaryLight },
  emoji: { fontSize: 28 },
  label: { ...typography.textStyles.h4, color: colors.text },
});
