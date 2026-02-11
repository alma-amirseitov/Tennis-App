import { Pressable, Text, ActivityIndicator, ViewStyle, TextStyle } from 'react-native';
import { colors, spacing, radius, typography } from '@/shared/theme';

export type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'small';

export interface ButtonProps {
  variant?: ButtonVariant;
  title: string;
  onPress: () => void;
  disabled?: boolean;
  loading?: boolean;
  icon?: React.ReactNode;
  style?: ViewStyle;
}

const HEIGHT_DEFAULT = 52;
const HEIGHT_SMALL = 36;

export function Button({
  variant = 'primary',
  title,
  onPress,
  disabled = false,
  loading = false,
  icon,
  style,
}: ButtonProps) {
  const isSmall = variant === 'small';
  const isDisabled = disabled || loading;

  const getBackgroundColor = (): string => {
    if (isDisabled && (variant === 'primary' || variant === 'small')) return '#9BCDB5';
    if (variant === 'primary' || variant === 'small') return colors.primary;
    if (variant === 'secondary') return colors.primaryLight;
    return 'transparent';
  };

  const getTextColor = (): string => {
    if (variant === 'primary' || variant === 'small') return colors.white;
    if (variant === 'secondary') return colors.primary;
    if (variant === 'outline') return isDisabled ? colors.textMuted : colors.text;
    return colors.text;
  };

  const getBorderColor = (): string => {
    if (variant === 'secondary') return colors.primary;
    if (variant === 'outline') return isDisabled ? colors.borderLight : colors.border;
    return 'transparent';
  };

  const containerStyle: ViewStyle = {
    height: isSmall ? HEIGHT_SMALL : HEIGHT_DEFAULT,
    backgroundColor: getBackgroundColor(),
    borderRadius: isSmall ? radius.sm : 14,
    borderWidth: variant === 'outline' || variant === 'secondary' ? 2 : 0,
    borderColor: getBorderColor(),
    paddingHorizontal: isSmall ? spacing.base : spacing.xl,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: spacing.sm,
  };

  const textStyle: TextStyle = {
    ...(isSmall ? typography.textStyles.buttonSm : typography.textStyles.button),
    color: getTextColor(),
  };

  return (
    <Pressable
      onPress={onPress}
      disabled={isDisabled}
      accessibilityRole="button"
      accessibilityState={{ disabled: isDisabled }}
      style={({ pressed }) => [
        containerStyle,
        pressed && !isDisabled && { transform: [{ scale: 0.97 }] },
        style,
      ]}
    >
      {loading ? (
        <ActivityIndicator
          size="small"
          color={variant === 'primary' || variant === 'small' ? colors.white : colors.primary}
        />
      ) : (
        <>
          {icon}
          <Text style={textStyle}>{title}</Text>
        </>
      )}
    </Pressable>
  );
}
