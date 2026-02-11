import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import i18n from '@/shared/i18n';
import {
  View,
  Text,
  ScrollView,
  StyleSheet,
  KeyboardAvoidingView,
  Platform,
} from 'react-native';
import { Button, Input, Select } from '@/shared/ui';
import { profileSetup } from '@/shared/api/auth';
import { authStore } from '@/shared/stores/auth';
import { showToast } from '@/shared/lib/toast';
import { colors, spacing, typography } from '@/shared/theme';
import axios from 'axios';
import type { ApiErrorResponse } from '@/shared/api/auth';

const DISTRICT_KEYS = [
  'district_esil',
  'district_almaty',
  'district_saryarka',
  'district_baikonur',
  'district_nurinsky',
] as const;

const BIRTH_YEARS = Array.from(
  { length: 2012 - 1940 + 1 },
  (_, i) => 1940 + i
).reverse();

const schema = z.object({
    first_name: z.string().min(2).max(50),
    last_name: z.string().min(2).max(50),
    gender: z.enum(['male', 'female']),
    birth_year: z.number().min(1940).max(2012),
    district: z.string().min(1),
  });

type FormData = z.infer<typeof schema>;

export default function ProfileSetupScreen() {
  const { t } = useTranslation();
  const router = useRouter();

  const {
    control,
    handleSubmit,
    formState: { errors, isValid, isSubmitting },
  } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      first_name: '',
      last_name: '',
      gender: 'male',
      birth_year: 1990,
      district: '',
    },
    mode: 'onChange',
  });

  const districtOptions = DISTRICT_KEYS.map((key) => ({
    value: key,
    label: t(`auth.${key}`),
  }));

  const yearOptions = BIRTH_YEARS.map((y) => ({
    value: String(y),
    label: String(y),
  }));

  const onSubmit = async (data: FormData) => {
    try {
      const districtLabel = t(`auth.${data.district}`);
      const lang = i18n.language === 'kk' ? 'kk' : i18n.language === 'en' ? 'en' : 'ru';
      const result = await profileSetup({
        first_name: data.first_name,
        last_name: data.last_name,
        gender: data.gender,
        birth_year: data.birth_year,
        city: t('auth.city_astana'),
        district: districtLabel,
        language: lang,
      });
      await authStore.getState().login(
        {
          access_token: result.access_token,
          refresh_token: result.refresh_token,
        },
        result.user as { id: string; first_name: string }
      );
      router.replace('/(auth)/quiz');
    } catch (err) {
      if (axios.isAxiosError(err)) {
        const data = err.response?.data as ApiErrorResponse | undefined;
        showToast(data?.error?.message ?? t('errors.something_went_wrong'));
      } else {
        showToast(t('errors.something_went_wrong'));
      }
    }
  };

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      keyboardVerticalOffset={Platform.OS === 'ios' ? 0 : 20}
    >
      <ScrollView
        style={styles.scroll}
        contentContainerStyle={styles.scrollContent}
        keyboardShouldPersistTaps="handled"
        showsVerticalScrollIndicator={false}
      >
        <Text style={styles.title}>{t('auth.profile_title')}</Text>

        <Controller
          control={control}
          name="first_name"
          render={({ field: { onChange, onBlur, value } }) => (
            <Input
              label={t('auth.first_name')}
              value={value}
              onChangeText={onChange}
              onBlur={onBlur}
              error={errors.first_name?.message}
              style={styles.field}
            />
          )}
        />

        <Controller
          control={control}
          name="last_name"
          render={({ field: { onChange, onBlur, value } }) => (
            <Input
              label={t('auth.last_name')}
              value={value}
              onChangeText={onChange}
              onBlur={onBlur}
              error={errors.last_name?.message}
              style={styles.field}
            />
          )}
        />

        <Controller
          control={control}
          name="gender"
          render={({ field: { onChange, value } }) => (
            <View style={styles.field}>
              <Text style={styles.label}>{t('auth.gender')}</Text>
              <View style={styles.genderRow}>
                <Button
                  variant={value === 'male' ? 'primary' : 'outline'}
                  title={t('auth.gender_male')}
                  onPress={() => onChange('male')}
                  style={styles.genderBtn}
                />
                <Button
                  variant={value === 'female' ? 'primary' : 'outline'}
                  title={t('auth.gender_female')}
                  onPress={() => onChange('female')}
                  style={styles.genderBtn}
                />
              </View>
            </View>
          )}
        />

        <View style={styles.field}>
          <Controller
            control={control}
            name="birth_year"
            render={({ field: { onChange, value } }) => (
              <Select
                label={t('auth.birth_year')}
                value={String(value)}
                onSelect={(v) => onChange(Number(v))}
                options={yearOptions}
                placeholder={t('auth.select_year')}
                error={errors.birth_year?.message}
              />
            )}
          />
        </View>

        <View style={styles.field}>
          <Controller
            control={control}
            name="district"
            render={({ field: { onChange, value } }) => (
              <Select
                label={t('auth.district')}
                value={value}
                onSelect={onChange}
                options={districtOptions}
                placeholder={t('auth.select_district')}
                error={errors.district?.message}
              />
            )}
          />
        </View>

        <Button
          variant="primary"
          title={t('auth.continue')}
          onPress={handleSubmit(onSubmit)}
          disabled={!isValid}
          loading={isSubmitting}
          style={styles.submitBtn}
        />
      </ScrollView>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  scroll: {
    flex: 1,
  },
  scrollContent: {
    padding: spacing.xl,
    paddingBottom: spacing['3xl'],
  },
  title: {
    ...typography.textStyles.h2,
    color: colors.text,
    marginBottom: spacing.xl,
  },
  label: {
    ...typography.textStyles.bodySm,
    fontWeight: typography.fontWeight.semibold,
    color: colors.textSecondary,
    marginBottom: spacing.sm,
  },
  field: {
    marginBottom: spacing.lg,
  },
  genderRow: {
    flexDirection: 'row',
    gap: spacing.sm,
  },
  genderBtn: {
    flex: 1,
  },
  submitBtn: {
    marginTop: spacing.lg,
  },
});
