import { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  StyleSheet,
  ViewStyle,
  TextInputProps,
} from 'react-native';
import { colors, spacing, radius, typography } from '@/shared/theme';

export interface InputProps extends Omit<TextInputProps, 'style'> {
  label?: string;
  placeholder?: string;
  value: string;
  onChangeText: (text: string) => void;
  error?: string;
  secureTextEntry?: boolean;
  style?: ViewStyle;
}

export function Input({
  label,
  placeholder,
  value,
  onChangeText,
  error,
  secureTextEntry,
  editable = true,
  style,
  ...rest
}: InputProps) {
  const [focused, setFocused] = useState(false);

  const borderColor = error
    ? colors.danger
    : focused
      ? colors.primary
      : colors.border;

  const borderWidth = focused || error ? 2 : 1;

  return (
    <View style={[styles.container, style]}>
      {label ? (
        <Text style={styles.label}>{label}</Text>
      ) : null}
      <TextInput
        value={value}
        onChangeText={onChangeText}
        onFocus={() => setFocused(true)}
        onBlur={() => setFocused(false)}
        placeholder={placeholder}
        placeholderTextColor={colors.textMuted}
        secureTextEntry={secureTextEntry}
        editable={editable}
        style={[
          styles.input,
          { borderColor, borderWidth },
          !editable && styles.inputDisabled,
        ]}
        {...rest}
      />
      {error ? <Text style={styles.errorText}>{error}</Text> : null}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    width: '100%',
  },
  label: {
    ...typography.textStyles.bodySm,
    fontWeight: typography.fontWeight.semibold,
    color: colors.textSecondary,
    marginBottom: spacing.sm,
  },
  input: {
    height: 52,
    borderRadius: radius.md,
    backgroundColor: colors.background,
    paddingHorizontal: spacing.base,
    paddingVertical: 14,
    ...typography.textStyles.body,
    color: colors.text,
  },
  inputDisabled: {
    opacity: 0.6,
  },
  errorText: {
    ...typography.textStyles.caption,
    color: colors.danger,
    marginTop: spacing.sm,
  },
});
