import { View, Text, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { Button, Input } from '@/shared/ui';
import { colors, spacing, typography } from '@/shared/theme';
import { useWizard } from './WizardContext';

export function StepDescription() {
  const { t } = useTranslation();
  const { data, update, setStep } = useWizard();
  return (
    <View style={styles.container}>
      <Text style={styles.title}>{t('events.wizard_step_7')}</Text>
      <Input
        label={t('events.create')}
        placeholder={t('events.create')}
        value={data.title}
        onChangeText={(v) => update({ title: v })}
      />
      <View style={styles.spacer} />
      <Input
        label={t('profile.about')}
        placeholder=""
        value={data.description}
        onChangeText={(v) => update({ description: v })}
        multiline
        numberOfLines={4}
      />
      <Button
        variant="primary"
        title={t('auth.quiz_next')}
        onPress={() => setStep(7)}
        disabled={!data.title}
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
