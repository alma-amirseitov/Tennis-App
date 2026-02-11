import { Text, StyleSheet } from 'react-native';
import { colors, spacing, radius, typography } from '@/shared/theme';

export type BadgeVariant = 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'muted';

const VARIANT_COLORS: Record<BadgeVariant, string> = {
  primary: colors.primary,
  success: colors.success,
  warning: colors.warning,
  danger: colors.danger,
  info: colors.info,
  muted: colors.textMuted,
};

function hexToRgba(hex: string, alpha: number): string {
  const r = parseInt(hex.slice(1, 3), 16);
  const g = parseInt(hex.slice(3, 5), 16);
  const b = parseInt(hex.slice(5, 7), 16);
  return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}

export interface BadgeProps {
  variant?: BadgeVariant;
  text: string;
}

export function Badge({ variant = 'primary', text }: BadgeProps) {
  const color = VARIANT_COLORS[variant];
  const bgColor = hexToRgba(color, 0.15);

  return (
    <Text
      style={[
        styles.badge,
        { backgroundColor: bgColor, color },
      ]}
    >
      {text}
    </Text>
  );
}

const styles = StyleSheet.create({
  badge: {
    ...typography.textStyles.badge,
    paddingVertical: 3,
    paddingHorizontal: 10,
    borderRadius: radius.pill,
    alignSelf: 'flex-start',
  },
});
