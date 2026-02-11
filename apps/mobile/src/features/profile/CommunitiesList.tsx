import { View, Text, ScrollView, Pressable, StyleSheet } from 'react-native';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { Avatar } from '@/shared/ui';
import { colors, spacing, typography, radius } from '@/shared/theme';

interface MiniCommunity {
  id: string;
  name: string;
  logo_url: string | null;
}

interface CommunitiesListProps {
  communities: MiniCommunity[];
}

export function CommunitiesList({ communities }: CommunitiesListProps) {
  const { t } = useTranslation();
  const router = useRouter();

  if (communities.length === 0) {
    return (
      <View style={styles.emptyContainer}>
        <Text style={styles.emptyText}>{t('profile.communities')}: 0</Text>
      </View>
    );
  }

  return (
    <ScrollView horizontal showsHorizontalScrollIndicator={false} contentContainerStyle={styles.scroll}>
      {communities.map((c) => (
        <Pressable key={c.id} onPress={() => router.push(`/community/${c.id}`)} style={styles.item}>
          <Avatar uri={c.logo_url} name={c.name} size="md" />
          <Text style={styles.name} numberOfLines={1}>{c.name}</Text>
        </Pressable>
      ))}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  scroll: { paddingHorizontal: spacing.base, gap: spacing.base },
  item: { alignItems: 'center', width: 72 },
  name: { ...typography.textStyles.caption, color: colors.text, marginTop: spacing.xs, textAlign: 'center' },
  emptyContainer: { paddingHorizontal: spacing.base, paddingVertical: spacing.sm },
  emptyText: { ...typography.textStyles.bodySm, color: colors.textMuted },
});
