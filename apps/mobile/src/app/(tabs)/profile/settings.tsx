import { View, Text, ScrollView, Pressable, StyleSheet, Alert, Platform } from 'react-native';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import i18n from '@/shared/i18n';
import { authStore } from '@/shared/stores/auth';
import { ScreenHeader } from '@/shared/ui';
import { colors, spacing, typography, radius } from '@/shared/theme';
import Constants from 'expo-constants';

type LangCode = 'ru' | 'kk' | 'en';
const LANGUAGES: { code: LangCode; label: string }[] = [
  { code: 'ru', label: 'Русский' },
  { code: 'kk', label: 'Қазақша' },
  { code: 'en', label: 'English' },
];

export default function SettingsScreen() {
  const { t } = useTranslation();
  const router = useRouter();
  const currentLang = (i18n.language || 'ru') as LangCode;

  const setLanguage = (lang: LangCode) => {
    i18n.changeLanguage(lang);
  };

  const handleLogout = () => {
    const doLogout = async () => {
      await authStore.getState().logout();
      router.replace('/(auth)');
    };

    if (Platform.OS === 'web') {
      if (window.confirm(t('profile.logout_confirm'))) {
        doLogout();
      }
    } else {
      Alert.alert(t('profile.logout'), t('profile.logout_confirm'), [
        { text: t('common.cancel'), style: 'cancel' },
        { text: t('profile.logout'), style: 'destructive', onPress: doLogout },
      ]);
    }
  };

  const version = Constants.expoConfig?.version ?? '1.0.0';

  return (
    <View style={styles.container}>
      <ScreenHeader title={t('profile.settings')} showBack />
      <ScrollView contentContainerStyle={styles.content} showsVerticalScrollIndicator={false}>
        {/* Language */}
        <Text style={styles.sectionTitle}>{t('profile.language')}</Text>
        <View style={styles.card}>
          {LANGUAGES.map((lang) => (
            <Pressable
              key={lang.code}
              onPress={() => setLanguage(lang.code)}
              style={[styles.langRow, currentLang === lang.code && styles.langRowActive]}
            >
              <Text style={[styles.langText, currentLang === lang.code && styles.langTextActive]}>
                {lang.label}
              </Text>
              {currentLang === lang.code ? <Text style={styles.check}>✓</Text> : null}
            </Pressable>
          ))}
        </View>

        {/* Notifications */}
        <Text style={styles.sectionTitle}>{t('profile.notifications_settings')}</Text>
        <View style={styles.card}>
          <SettingRow label={t('profile.notifications_settings')} value="→" />
        </View>

        {/* Privacy */}
        <Text style={styles.sectionTitle}>{t('profile.privacy')}</Text>
        <View style={styles.card}>
          <SettingRow label={t('profile.privacy')} value="→" />
        </View>

        {/* About */}
        <Text style={styles.sectionTitle}>{t('profile.about')}</Text>
        <View style={styles.card}>
          <SettingRow label={t('profile.about')} value={`v${version}`} />
        </View>

        {/* Logout */}
        <Pressable onPress={handleLogout} style={styles.logoutButton}>
          <Text style={styles.logoutText}>{t('profile.logout')}</Text>
        </Pressable>
      </ScrollView>
    </View>
  );
}

function SettingRow({ label, value }: { label: string; value: string }) {
  return (
    <View style={styles.settingRow}>
      <Text style={styles.settingLabel}>{label}</Text>
      <Text style={styles.settingValue}>{value}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  content: { padding: spacing.base, paddingBottom: spacing['3xl'] },
  sectionTitle: {
    ...typography.textStyles.caption,
    color: colors.textMuted,
    textTransform: 'uppercase',
    marginTop: spacing.xl,
    marginBottom: spacing.sm,
    paddingHorizontal: spacing.xs,
  },
  card: {
    backgroundColor: colors.card,
    borderRadius: radius.lg,
    borderWidth: 1,
    borderColor: colors.borderLight,
    overflow: 'hidden',
  },
  langRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 14,
    paddingHorizontal: spacing.base,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: colors.borderLight,
  },
  langRowActive: { backgroundColor: colors.primaryLight },
  langText: { ...typography.textStyles.body, color: colors.text },
  langTextActive: { color: colors.primary, fontWeight: typography.fontWeight.semibold },
  check: { ...typography.textStyles.body, color: colors.primary, fontWeight: typography.fontWeight.bold },
  settingRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 14,
    paddingHorizontal: spacing.base,
  },
  settingLabel: { ...typography.textStyles.body, color: colors.text },
  settingValue: { ...typography.textStyles.bodySm, color: colors.textMuted },
  logoutButton: {
    marginTop: spacing['2xl'],
    backgroundColor: colors.card,
    borderRadius: radius.lg,
    borderWidth: 1,
    borderColor: colors.danger,
    paddingVertical: 14,
    alignItems: 'center',
  },
  logoutText: { ...typography.textStyles.button, color: colors.danger },
});
