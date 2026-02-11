import { useState } from 'react';
import { View, Text, TextInput, StyleSheet, KeyboardAvoidingView, Platform, ScrollView } from 'react-native';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { Button } from '@/shared/ui/Button';
import { sendOTP } from '@/shared/api/auth';
import type { ApiErrorResponse } from '@/shared/api/auth';
import { showToast } from '@/shared/lib/toast';
import { colors, spacing, typography, radius } from '@/shared/theme';
import axios from 'axios';

const PHONE_PREFIX = '+7';
const DIGITS_LENGTH = 10;

function formatDisplay(digits: string): string {
  if (digits.length === 0) return '';
  if (digits.length <= 3) return `(${digits}`;
  if (digits.length <= 6) return `(${digits.slice(0, 3)}) ${digits.slice(3)}`;
  if (digits.length <= 8) return `(${digits.slice(0, 3)}) ${digits.slice(3, 6)}-${digits.slice(6)}`;
  return `(${digits.slice(0, 3)}) ${digits.slice(3, 6)}-${digits.slice(6, 8)}-${digits.slice(8)}`;
}

export default function PhoneInputScreen() {
  const { t } = useTranslation();
  const router = useRouter();
  const [rawDigits, setRawDigits] = useState('');
  const [loading, setLoading] = useState(false);

  const isValid = rawDigits.length === DIGITS_LENGTH;

  const handleChange = (text: string) => {
    // Extract only digits from whatever the user types
    const digits = text.replace(/\D/g, '').slice(0, DIGITS_LENGTH);
    setRawDigits(digits);
  };

  const handleSubmit = async () => {
    if (!isValid || loading) return;
    setLoading(true);
    try {
      const fullPhone = `${PHONE_PREFIX}${rawDigits}`;
      const { session_id } = await sendOTP(fullPhone);
      router.push({
        pathname: '/(auth)/otp',
        params: { session_id, phone: fullPhone },
      });
    } catch (err) {
      if (axios.isAxiosError(err)) {
        const data = err.response?.data as ApiErrorResponse | undefined;
        const code = data?.error?.code;
        const message = data?.error?.message;

        if (code === 'RATE_LIMITED' || code === 'SMS_RATE_LIMITED') {
          const retryAfter = err.response?.headers?.['retry-after'];
          const msg = retryAfter
            ? t('errors.SMS_RATE_LIMITED') + ` (${retryAfter}s)`
            : t('errors.SMS_RATE_LIMITED');
          showToast(msg);
        } else if (err.message === 'Network Error' || !err.response) {
          showToast(t('errors.no_internet'));
        } else {
          showToast(message ?? t('errors.something_went_wrong'));
        }
      } else {
        showToast(t('errors.something_went_wrong'));
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
    >
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        keyboardShouldPersistTaps="handled"
        bounces={false}
      >
        <View style={styles.content}>
          <Text style={styles.title}>{t('auth.phone_title')}</Text>
          <Text style={styles.subtitle}>{t('auth.phone_subtitle')}</Text>

          <View style={styles.inputWrapper}>
            <Text style={styles.prefix}>{PHONE_PREFIX} </Text>
            <TextInput
              style={styles.input}
              value={formatDisplay(rawDigits)}
              onChangeText={handleChange}
              placeholder="(___) ___-__-__"
              placeholderTextColor={colors.textMuted}
              keyboardType="number-pad"
              maxLength={16}
              editable={!loading}
              autoFocus
              autoComplete="tel"
            />
          </View>

          <Button
            variant="primary"
            title={t('auth.get_code')}
            onPress={handleSubmit}
            disabled={!isValid}
            loading={loading}
            style={styles.button}
          />
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  scrollContent: {
    flexGrow: 1,
    justifyContent: 'center',
    paddingHorizontal: spacing.xl,
  },
  content: {
    width: '100%',
  },
  title: {
    ...typography.textStyles.h2,
    color: colors.text,
    marginBottom: spacing.sm,
  },
  subtitle: {
    ...typography.textStyles.bodySm,
    color: colors.textSecondary,
    marginBottom: spacing.xl,
  },
  inputWrapper: {
    flexDirection: 'row',
    alignItems: 'center',
    height: 52,
    borderRadius: radius.md,
    backgroundColor: colors.card,
    borderWidth: 1,
    borderColor: colors.border,
    paddingHorizontal: spacing.base,
    marginBottom: spacing.xl,
  },
  prefix: {
    ...typography.textStyles.body,
    color: colors.text,
  },
  input: {
    flex: 1,
    ...typography.textStyles.body,
    color: colors.text,
    padding: 0,
  },
  button: {
    width: '100%',
  },
});
