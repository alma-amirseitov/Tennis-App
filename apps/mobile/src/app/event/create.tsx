import { View, Text, Pressable, StyleSheet, ScrollView } from 'react-native';
import { useTranslation } from 'react-i18next';
import { ScreenHeader } from '@/shared/ui';
import { colors, spacing, typography } from '@/shared/theme';
import { WizardProvider, useWizard } from '@/features/events/wizard/WizardContext';
import { StepType } from '@/features/events/wizard/StepType';
import { StepComposition } from '@/features/events/wizard/StepComposition';
import { StepLevel } from '@/features/events/wizard/StepLevel';
import { StepDetails } from '@/features/events/wizard/StepDetails';
import { StepLocation } from '@/features/events/wizard/StepLocation';
import { StepDateTime } from '@/features/events/wizard/StepDateTime';
import { StepDescription } from '@/features/events/wizard/StepDescription';
import { StepReview } from '@/features/events/wizard/StepReview';

function WizardContent() {
  const { t } = useTranslation();
  const { step, setStep, totalSteps } = useWizard();
  const progress = ((step + 1) / totalSteps) * 100;

  const STEPS = [StepType, StepComposition, StepLevel, StepDetails, StepLocation, StepDateTime, StepDescription, StepReview];
  const StepComponent = STEPS[step] ?? StepType;

  return (
    <View style={styles.container}>
      <ScreenHeader
        title={t('events.create')}
        showBack
        right={
          <Text style={styles.stepLabel}>{step + 1}/{totalSteps}</Text>
        }
      />

      {/* Progress bar */}
      <View style={styles.progressBar}>
        <View style={[styles.progressFill, { width: `${progress}%` }]} />
      </View>

      {/* Back navigation within wizard */}
      {step > 0 ? (
        <Pressable onPress={() => setStep(step - 1)} style={styles.wizardBack}>
          <Text style={styles.wizardBackText}>← {t('common.back')}</Text>
        </Pressable>
      ) : null}

      <ScrollView showsVerticalScrollIndicator={false} keyboardShouldPersistTaps="handled">
        <StepComponent />
      </ScrollView>
    </View>
  );
}

export default function CreateEventScreen() {
  return (
    <WizardProvider>
      <WizardContent />
    </WizardProvider>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  stepLabel: { ...typography.textStyles.caption, color: colors.textMuted },
  progressBar: { height: 3, backgroundColor: colors.borderLight },
  progressFill: { height: '100%', backgroundColor: colors.primary },
  wizardBack: { paddingHorizontal: spacing.xl, paddingVertical: spacing.sm },
  wizardBackText: { ...typography.textStyles.body, color: colors.primary },
});
