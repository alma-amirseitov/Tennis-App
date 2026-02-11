import { View, Text, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { Button } from './Button';
import { colors, spacing, typography } from '@/shared/theme';

interface ErrorStateProps {
  message?: string;
  onRetry?: () => void;
}

export function ErrorState({ message, onRetry }: ErrorStateProps) {
  const { t } = useTranslation();
  return (
    <View style={styles.container}>
      <Text style={styles.emoji}>⚠️</Text>
      <Text style={styles.title}>{t('errors.something_went_wrong')}</Text>
      {message ? <Text style={styles.description}>{message}</Text> : null}
      {onRetry ? (
        <Button variant="primary" title={t('common.retry')} onPress={onRetry} style={styles.button} />
      ) : null}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, justifyContent: 'center', alignItems: 'center', padding: spacing.xl },
  emoji: { fontSize: 48, marginBottom: spacing.base },
  title: { ...typography.textStyles.h3, color: colors.text, textAlign: 'center', marginBottom: spacing.xs },
  description: { ...typography.textStyles.bodySm, color: colors.textMuted, textAlign: 'center', marginBottom: spacing.base },
  button: { marginTop: spacing.base },
});
