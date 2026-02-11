import { View, Text, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { Button, Input } from '@/shared/ui';
import { colors, spacing, typography } from '@/shared/theme';
import { useWizard } from './WizardContext';

export function StepLocation() {
  const { t } = useTranslation();
  const { data, update, setStep } = useWizard();
  return (
    <View style={styles.container}>
      <Text style={styles.title}>{t('events.wizard_step_5')}</Text>
      <Input
        label={t('events.detail_location')}
        placeholder={t('events.detail_location')}
        value={data.location_name}
        onChangeText={(v) => update({ location_name: v })}
      />
      <Button
        variant="primary"
        title={t('auth.quiz_next')}
        onPress={() => setStep(5)}
        disabled={!data.location_name}
        style={styles.nextBtn}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: { padding: spacing.xl },
  title: { ...typography.textStyles.h2, color: colors.text, marginBottom: spacing.xl },
  nextBtn: { marginTop: spacing.xl },
});
