import { View, Text, StyleSheet } from 'react-native';
import { Button } from './Button';
import { colors, spacing, typography } from '@/shared/theme';

interface EmptyStateProps {
  emoji?: string;
  title: string;
  description?: string;
  actionTitle?: string;
  onAction?: () => void;
}

export function EmptyState({ emoji = '🔍', title, description, actionTitle, onAction }: EmptyStateProps) {
  return (
    <View style={styles.container}>
      <Text style={styles.emoji}>{emoji}</Text>
      <Text style={styles.title}>{title}</Text>
      {description ? <Text style={styles.description}>{description}</Text> : null}
      {actionTitle && onAction ? (
        <Button variant="primary" title={actionTitle} onPress={onAction} style={styles.button} />
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
