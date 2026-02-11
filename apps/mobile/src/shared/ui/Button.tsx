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

  const getBackgroundColor = (): string => {
    if (variant === 'primary' || variant === 'small') return colors.primary;
    if (variant === 'secondary') return colors.primaryLight;
    return 'transparent';
  };

  const getTextColor = (): string => {
    if (variant === 'primary' || variant === 'small') return colors.white;
    if (variant === 'secondary') return colors.primary;
    if (variant === 'outline') return colors.text;
    return colors.text;
  };

  const getBorderColor = (): string => {
    if (variant === 'secondary') return colors.primary;
    if (variant === 'outline') return colors.border;
    return 'transparent';
  };

  const containerStyle: ViewStyle = {
    height: isSmall ? HEIGHT_SMALL : HEIGHT_DEFAULT,
    backgroundColor: variant === 'outline' ? 'transparent' : getBackgroundColor(),
    borderRadius: isSmall ? radius.sm : 14,
    borderWidth: variant === 'outline' || variant === 'secondary' ? 2 : 0,
    borderColor: getBorderColor(),
    paddingHorizontal: isSmall ? spacing.base : spacing.xl,
    opacity: disabled ? 0.5 : 1,
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
      disabled={disabled || loading}
      style={({ pressed }) => [
        containerStyle,
        pressed && !disabled && !loading && { transform: [{ scale: 0.97 }] },
        style,
      ]}
    >
      {loading ? (
        <ActivityIndicator
          size="small"
          color={variant === 'primary' ? colors.white : colors.primary}
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
