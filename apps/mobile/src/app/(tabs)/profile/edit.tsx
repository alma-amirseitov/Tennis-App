import { useState } from 'react';
import { View, Text, ScrollView, StyleSheet, KeyboardAvoidingView, Platform, Pressable, Image } from 'react-native';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import * as ImagePicker from 'expo-image-picker';
import { useProfile, useUpdateProfile } from '@/shared/api/hooks';
import { uploadAvatar } from '@/shared/api/users';
import { Button, Input, Select, ScreenHeader, Avatar } from '@/shared/ui';
import { showToast } from '@/shared/lib/toast';
import { colors, spacing, typography, radius } from '@/shared/theme';

const DISTRICT_KEYS = [
  'district_esil', 'district_almaty', 'district_saryarka', 'district_baikonur', 'district_nurinsky',
] as const;

const schema = z.object({
  first_name: z.string().min(2).max(50),
  last_name: z.string().min(2).max(50),
  gender: z.enum(['male', 'female']),
  birth_year: z.number().min(1940).max(2012),
  district: z.string().min(1),
});

type FormData = z.infer<typeof schema>;

const BIRTH_YEARS = Array.from({ length: 2012 - 1940 + 1 }, (_, i) => 1940 + i).reverse();

export default function EditProfileScreen() {
  const { t } = useTranslation();
  const router = useRouter();
  const { data: user } = useProfile();
  const updateProfile = useUpdateProfile();
  const [avatarUri, setAvatarUri] = useState<string | null>(null);
  const [uploading, setUploading] = useState(false);

  const { control, handleSubmit, formState: { errors, isValid, isSubmitting } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      first_name: user?.first_name ?? '',
      last_name: user?.last_name ?? '',
      gender: user?.gender ?? 'male',
      birth_year: user?.birth_year ?? 1990,
      district: user?.district ?? '',
    },
    mode: 'onChange',
  });

  const districtOptions = DISTRICT_KEYS.map((key) => ({ value: key, label: t(`auth.${key}`) }));
  const yearOptions = BIRTH_YEARS.map((y) => ({ value: String(y), label: String(y) }));

  const pickImage = async () => {
    const result = await ImagePicker.launchImageLibraryAsync({
      mediaTypes: ['images'],
      allowsEditing: true,
      aspect: [1, 1],
      quality: 0.8,
    });
    if (!result.canceled && result.assets[0]) {
      setAvatarUri(result.assets[0].uri);
      setUploading(true);
      try {
        await uploadAvatar(result.assets[0].uri);
      } catch {
        showToast(t('errors.something_went_wrong'));
      } finally {
        setUploading(false);
      }
    }
  };

  const onSubmit = async (data: FormData) => {
    try {
      await updateProfile.mutateAsync({
        first_name: data.first_name,
        last_name: data.last_name,
        gender: data.gender,
        birth_year: data.birth_year,
        district: t(`auth.${data.district}`),
      });
      showToast(t('common.save'));
      router.back();
    } catch {
      showToast(t('errors.something_went_wrong'));
    }
  };

  return (
    <View style={styles.container}>
      <ScreenHeader title={t('profile.edit')} showBack />
      <KeyboardAvoidingView style={styles.flex} behavior={Platform.OS === 'ios' ? 'padding' : 'height'}>
        <ScrollView contentContainerStyle={styles.scrollContent} keyboardShouldPersistTaps="handled" showsVerticalScrollIndicator={false}>
          {/* Avatar */}
          <Pressable onPress={pickImage} style={styles.avatarSection}>
            {avatarUri ? (
              <Image source={{ uri: avatarUri }} style={styles.avatarImage} />
            ) : (
              <Avatar uri={user?.avatar_url} name={`${user?.first_name ?? ''} ${user?.last_name ?? ''}`} size="xl" />
            )}
            <Text style={styles.changePhoto}>{uploading ? t('common.loading') : '📷 ' + t('auth.continue')}</Text>
          </Pressable>

          <Controller control={control} name="first_name" render={({ field: { onChange, onBlur, value } }) => (
            <Input label={t('auth.first_name')} value={value} onChangeText={onChange} onBlur={onBlur} error={errors.first_name?.message} style={styles.field} />
          )} />

          <Controller control={control} name="last_name" render={({ field: { onChange, onBlur, value } }) => (
            <Input label={t('auth.last_name')} value={value} onChangeText={onChange} onBlur={onBlur} error={errors.last_name?.message} style={styles.field} />
          )} />

          <Controller control={control} name="gender" render={({ field: { onChange, value } }) => (
            <View style={styles.field}>
              <Text style={styles.label}>{t('auth.gender')}</Text>
              <View style={styles.genderRow}>
                <Button variant={value === 'male' ? 'primary' : 'outline'} title={t('auth.gender_male')} onPress={() => onChange('male')} style={styles.genderBtn} />
                <Button variant={value === 'female' ? 'primary' : 'outline'} title={t('auth.gender_female')} onPress={() => onChange('female')} style={styles.genderBtn} />
              </View>
            </View>
          )} />

          <View style={styles.field}>
            <Controller control={control} name="birth_year" render={({ field: { onChange, value } }) => (
              <Select label={t('auth.birth_year')} value={String(value)} onSelect={(v) => onChange(Number(v))} options={yearOptions} placeholder={t('auth.select_year')} />
            )} />
          </View>

          <View style={styles.field}>
            <Controller control={control} name="district" render={({ field: { onChange, value } }) => (
              <Select label={t('auth.district')} value={value} onSelect={onChange} options={districtOptions} placeholder={t('auth.select_district')} />
            )} />
          </View>

          <Button variant="primary" title={t('common.save')} onPress={handleSubmit(onSubmit)} disabled={!isValid} loading={isSubmitting} style={styles.submitBtn} />
        </ScrollView>
      </KeyboardAvoidingView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  flex: { flex: 1 },
  scrollContent: { padding: spacing.xl, paddingBottom: spacing['3xl'] },
  avatarSection: { alignItems: 'center', marginBottom: spacing.xl },
  avatarImage: { width: 80, height: 80, borderRadius: 40 },
  changePhoto: { ...typography.textStyles.bodySm, color: colors.primary, marginTop: spacing.sm },
  label: { ...typography.textStyles.bodySm, fontWeight: typography.fontWeight.semibold, color: colors.textSecondary, marginBottom: spacing.sm },
  field: { marginBottom: spacing.lg },
  genderRow: { flexDirection: 'row', gap: spacing.sm },
  genderBtn: { flex: 1 },
  submitBtn: { marginTop: spacing.lg },
});
