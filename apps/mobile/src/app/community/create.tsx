import { View, Text, ScrollView, StyleSheet, KeyboardAvoidingView, Platform } from 'react-native';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useCreateCommunity } from '@/shared/api/hooks';
import { Button, Input, ScreenHeader, Card } from '@/shared/ui';
import { showToast } from '@/shared/lib/toast';
import { colors, spacing, typography } from '@/shared/theme';
import type { CommunityType, CommunityAccess } from '@/shared/api/communities';

const TYPES: { key: CommunityType; emoji: string; labelKey: string }[] = [
  { key: 'club', emoji: '🎾', labelKey: 'communities.type_club' },
  { key: 'league', emoji: '🏆', labelKey: 'communities.type_league' },
  { key: 'organizer', emoji: '📋', labelKey: 'communities.type_organizer' },
  { key: 'group', emoji: '👥', labelKey: 'communities.type_group' },
];

const ACCESS_TYPES: { key: CommunityAccess; labelKey: string }[] = [
  { key: 'open', labelKey: 'communities.join' },
  { key: 'closed', labelKey: 'communities.pending' },
];

const schema = z.object({
  name: z.string().min(3).max(100),
  description: z.string().min(10).max(500),
  community_type: z.enum(['club', 'league', 'organizer', 'group']),
  access_type: z.enum(['open', 'closed']),
});

type FormData = z.infer<typeof schema>;

export default function CreateCommunityScreen() {
  const { t } = useTranslation();
  const router = useRouter();
  const createMutation = useCreateCommunity();

  const { control, handleSubmit, setValue, watch, formState: { isValid, isSubmitting } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { name: '', description: '', community_type: 'group', access_type: 'open' },
    mode: 'onChange',
  });

  const selectedType = watch('community_type');
  const selectedAccess = watch('access_type');

  const onSubmit = async (data: FormData) => {
    try {
      await createMutation.mutateAsync(data);
      showToast(t('common.done'));
      router.back();
    } catch {
      showToast(t('errors.something_went_wrong'));
    }
  };

  return (
    <View style={styles.container}>
      <ScreenHeader title={t('communities.create')} showBack />
      <KeyboardAvoidingView style={styles.flex} behavior={Platform.OS === 'ios' ? 'padding' : 'height'}>
        <ScrollView contentContainerStyle={styles.content} keyboardShouldPersistTaps="handled" showsVerticalScrollIndicator={false}>
          {/* Type selector */}
          <Text style={styles.sectionLabel}>{t('common.filter')}</Text>
          <View style={styles.typeGrid}>
            {TYPES.map((tp) => (
              <Card key={tp.key} onPress={() => setValue('community_type', tp.key, { shouldValidate: true })}
                style={[styles.typeCard, selectedType === tp.key && styles.typeCardSelected]}>
                <Text style={styles.typeEmoji}>{tp.emoji}</Text>
                <Text style={[styles.typeLabel, selectedType === tp.key && styles.typeLabelSelected]}>{t(tp.labelKey)}</Text>
              </Card>
            ))}
          </View>

          {/* Name */}
          <Controller control={control} name="name" render={({ field: { onChange, onBlur, value } }) => (
            <Input label={t('auth.first_name')} placeholder={t('communities.create')} value={value} onChangeText={onChange} onBlur={onBlur} style={styles.field} />
          )} />

          {/* Description */}
          <Controller control={control} name="description" render={({ field: { onChange, onBlur, value } }) => (
            <Input label={t('profile.about')} value={value} onChangeText={onChange} onBlur={onBlur} multiline numberOfLines={4} style={styles.field} />
          )} />

          {/* Access type */}
          <Text style={styles.sectionLabel}>{t('profile.privacy')}</Text>
          <View style={styles.accessRow}>
            {ACCESS_TYPES.map((ac) => (
              <Button key={ac.key} variant={selectedAccess === ac.key ? 'primary' : 'outline'} title={t(ac.labelKey)}
                onPress={() => setValue('access_type', ac.key, { shouldValidate: true })} style={styles.accessBtn} />
            ))}
          </View>

          <Button variant="primary" title={t('communities.create')} onPress={handleSubmit(onSubmit)} disabled={!isValid} loading={isSubmitting} style={styles.submit} />
        </ScrollView>
      </KeyboardAvoidingView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  flex: { flex: 1 },
  content: { padding: spacing.xl, paddingBottom: spacing['3xl'] },
  sectionLabel: { ...typography.textStyles.h4, color: colors.text, marginBottom: spacing.md, marginTop: spacing.lg },
  typeGrid: { flexDirection: 'row', flexWrap: 'wrap', gap: spacing.sm },
  typeCard: { width: '47%', alignItems: 'center', paddingVertical: spacing.lg },
  typeCardSelected: { borderColor: colors.primary, backgroundColor: colors.primaryLight },
  typeEmoji: { fontSize: 28, marginBottom: spacing.sm },
  typeLabel: { ...typography.textStyles.bodySm, color: colors.text },
  typeLabelSelected: { color: colors.primary, fontWeight: typography.fontWeight.semibold },
  field: { marginBottom: spacing.lg },
  accessRow: { flexDirection: 'row', gap: spacing.sm },
  accessBtn: { flex: 1 },
  submit: { marginTop: spacing.xl },
});
