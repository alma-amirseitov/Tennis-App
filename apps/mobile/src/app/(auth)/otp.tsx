import { useState, useRef, useEffect } from 'react';
import {
  View,
  Text,
  TextInput,
  StyleSheet,
  Animated,
  Pressable,
  Platform,
} from 'react-native';
import { useRouter, useLocalSearchParams } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { verifyOTP, sendOTP } from '@/shared/api/auth';
import type { ApiErrorResponse } from '@/shared/api/auth';
import { authStore } from '@/shared/stores/auth';
import { showToast } from '@/shared/lib/toast';
import { colors, spacing, typography, radius } from '@/shared/theme';
import axios from 'axios';

const CELL_COUNT = 4;
const RESEND_SECONDS = 60;

function maskPhone(phone: string): string {
  if (phone.length < 10) return phone;
  const digits = phone.replace(/\D/g, '');
  return `+7 ${digits.slice(0, 3)} ***-**-${digits.slice(-2)}`;
}

export default function OtpScreen() {
  const { t } = useTranslation();
  const router = useRouter();
  const params = useLocalSearchParams<{ session_id: string; phone: string }>();
  const [sessionId, setSessionId] = useState(params.session_id ?? '');
  const phone = params.phone ?? '';

  useEffect(() => {
    if (params.session_id) setSessionId(params.session_id);
  }, [params.session_id]);

  const [code, setCode] = useState(['', '', '', '']);
  const [loading, setLoading] = useState(false);
  const submittedRef = useRef(false);
  const loadingRef = useRef(false);
  const [resendSeconds, setResendSeconds] = useState(RESEND_SECONDS);
  const [canResend, setCanResend] = useState(false);
  const resendTimerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const shakeAnim = useRef(new Animated.Value(0)).current;
  const inputRefs = useRef<(TextInput | null)[]>([]);

  useEffect(() => {
    inputRefs.current[0]?.focus();
  }, []);

  useEffect(() => {
    resendTimerRef.current = setInterval(() => {
      setResendSeconds((s) => {
        if (s <= 1) {
          if (resendTimerRef.current) {
            clearInterval(resendTimerRef.current);
            resendTimerRef.current = null;
          }
          setCanResend(true);
          return RESEND_SECONDS;
        }
        return s - 1;
      });
    }, 1000);
    return () => {
      if (resendTimerRef.current) clearInterval(resendTimerRef.current);
    };
  }, []);

  const triggerShake = () => {
    shakeAnim.setValue(0);
    Animated.sequence([
      Animated.timing(shakeAnim, { toValue: 1, duration: 50, useNativeDriver: true }),
      Animated.timing(shakeAnim, { toValue: 2, duration: 50, useNativeDriver: true }),
      Animated.timing(shakeAnim, { toValue: 3, duration: 50, useNativeDriver: true }),
      Animated.timing(shakeAnim, { toValue: 0, duration: 50, useNativeDriver: true }),
    ]).start();
  };

  const handleChange = (index: number, value: string) => {
    if (submittedRef.current || loadingRef.current) return;

    if (value.length > 1) {
      const digits = value.replace(/\D/g, '').slice(0, CELL_COUNT).split('');
      const newCode = [...code];
      digits.forEach((d, i) => {
        if (i < CELL_COUNT) newCode[i] = d;
      });
      setCode(newCode);
      const nextEmpty = newCode.findIndex((c) => !c);
      const focusIndex = nextEmpty === -1 ? CELL_COUNT - 1 : nextEmpty;
      inputRefs.current[focusIndex]?.focus();
      if (newCode.every((c) => c)) submitCode(newCode.join(''));
      return;
    }

    const newCode = [...code];
    newCode[index] = value.replace(/\D/g, '').slice(-1);
    setCode(newCode);

    if (value && index < CELL_COUNT - 1) {
      inputRefs.current[index + 1]?.focus();
    }
    if (newCode.every((c) => c)) submitCode(newCode.join(''));
  };

  const handleKeyPress = (index: number, key: string) => {
    if (key === 'Backspace' && !code[index] && index > 0) {
      inputRefs.current[index - 1]?.focus();
    }
  };

  const submitCode = async (fullCode: string) => {
    if (!sessionId || loadingRef.current || submittedRef.current) return;
    submittedRef.current = true;
    loadingRef.current = true;
    setLoading(true);
    try {
      const data = await verifyOTP(sessionId, fullCode);

      if (data.is_new) {
        const userId = data.user_id ?? (data as unknown as { user?: { id?: string } }).user?.id ?? '';
        await authStore.getState().setTempToken(data.temp_token, userId);
        router.replace('/(auth)/profile-setup');
      } else {
        await authStore.getState().login(
          { access_token: data.access_token, refresh_token: data.refresh_token },
          data.user
        );
        router.replace('/(tabs)');
      }
    } catch (err) {
      submittedRef.current = false;
      loadingRef.current = false;
      setLoading(false);
      setCode(['', '', '', '']);
      inputRefs.current[0]?.focus();
      triggerShake();

      if (axios.isAxiosError(err)) {
        const errData = err.response?.data as ApiErrorResponse | undefined;
        const errCode = errData?.error?.code;

        if (errCode === 'OTP_INVALID_CODE') {
          showToast(t('auth.wrong_code'));
        } else if (errCode === 'OTP_SESSION_EXPIRED') {
          showToast(t('errors.OTP_SESSION_EXPIRED'));
        } else if (errCode === 'OTP_MAX_ATTEMPTS') {
          showToast(t('errors.OTP_MAX_ATTEMPTS'));
        } else if (err.message === 'Network Error' || !err.response) {
          showToast(t('errors.no_internet'));
        } else {
          showToast(errData?.error?.message ?? t('errors.something_went_wrong'));
        }
      } else {
        showToast(t('errors.something_went_wrong'));
      }
    }
  };

  const handleResend = async () => {
    if (!canResend || !phone) return;
    setCanResend(false);
    setResendSeconds(RESEND_SECONDS);
    submittedRef.current = false;
    try {
      const { session_id } = await sendOTP(phone);
      setSessionId(session_id);
      router.setParams({ session_id, phone });
      resendTimerRef.current = setInterval(() => {
        setResendSeconds((s) => {
          if (s <= 1) {
            if (resendTimerRef.current) {
              clearInterval(resendTimerRef.current);
              resendTimerRef.current = null;
            }
            setCanResend(true);
            return RESEND_SECONDS;
          }
          return s - 1;
        });
      }, 1000);
    } catch {
      setCanResend(true);
      showToast(t('errors.SMS_RATE_LIMITED'));
    }
  };

  const handleBack = () => {
    router.back();
  };

  const shakeTranslate = shakeAnim.interpolate({
    inputRange: [0, 1, 2, 3],
    outputRange: [0, 10, -10, 0],
  });

  if (!sessionId) {
    router.replace('/(auth)');
    return null;
  }

  return (
    <View style={styles.container}>
      <Pressable onPress={handleBack} style={styles.backButton}>
        <Text style={styles.backText}>{t('common.back')}</Text>
      </Pressable>

      <Text style={styles.title}>{t('auth.otp_title')}</Text>
      <Text style={styles.subtitle}>
        {t('auth.otp_subtitle', { phone: maskPhone(phone) })}
      </Text>

      <Animated.View
        style={[styles.cells, { transform: [{ translateX: shakeTranslate }] }]}
      >
        {code.map((digit, index) => (
          <TextInput
            key={index}
            ref={(el) => { inputRefs.current[index] = el; }}
            style={[styles.cell, loading && styles.cellDisabled]}
            value={digit}
            onChangeText={(v) => handleChange(index, v)}
            onKeyPress={({ nativeEvent }) => handleKeyPress(index, nativeEvent.key)}
            keyboardType="number-pad"
            textContentType={Platform.OS === 'ios' ? 'oneTimeCode' : undefined}
            autoComplete={Platform.OS === 'android' ? 'sms-otp' : undefined}
            maxLength={4}
            selectTextOnFocus
            editable={!loading}
          />
        ))}
      </Animated.View>

      <View style={styles.resend}>
        {canResend ? (
          <Pressable onPress={handleResend} hitSlop={12}>
            <Text style={styles.resendLink}>{t('auth.resend_code')}</Text>
          </Pressable>
        ) : (
          <Text style={styles.resendText}>
            {t('auth.resend_in', { seconds: resendSeconds })}
          </Text>
        )}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
    paddingHorizontal: spacing.xl,
    paddingTop: spacing.xl,
  },
  backButton: {
    alignSelf: 'flex-start',
    paddingVertical: spacing.sm,
    marginBottom: spacing.lg,
  },
  backText: {
    ...typography.textStyles.body,
    color: colors.primary,
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
  cells: {
    flexDirection: 'row',
    gap: spacing.sm,
    marginBottom: spacing.xl,
  },
  cell: {
    flex: 1,
    height: 56,
    borderRadius: radius.md,
    borderWidth: 2,
    borderColor: colors.border,
    backgroundColor: colors.card,
    ...typography.textStyles.h3,
    color: colors.text,
    textAlign: 'center',
    padding: 0,
  },
  cellDisabled: {
    opacity: 0.6,
  },
  resend: {
    alignItems: 'center',
  },
  resendLink: {
    ...typography.textStyles.body,
    color: colors.primary,
  },
  resendText: {
    ...typography.textStyles.bodySm,
    color: colors.textMuted,
  },
});
