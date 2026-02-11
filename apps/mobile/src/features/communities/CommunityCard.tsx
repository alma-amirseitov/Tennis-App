import { View, Text, Pressable, StyleSheet } from 'react-native';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { Avatar, Badge } from '@/shared/ui';
import { colors, spacing, typography, radius } from '@/shared/theme';
import type { CommunityItem, CommunityType } from '@/shared/api/communities';

const TYPE_LABELS: Record<CommunityType, string> = {
  club: 'communities.type_club',
  league: 'communities.type_league',
  organizer: 'communities.type_organizer',
  group: 'communities.type_group',
};

interface CommunityCardProps {
  community: CommunityItem;
}

export function CommunityCard({ community }: CommunityCardProps) {
  const { t } = useTranslation();
  const router = useRouter();
  return (
    <Pressable
      onPress={() => router.push(`/community/${community.id}`)}
      style={({ pressed }) => [styles.card, pressed && { transform: [{ scale: 0.98 }] }]}
    >
      <Avatar uri={community.logo_url} name={community.name} size="lg" />
      <View style={styles.info}>
        <View style={styles.nameRow}>
          <Text style={styles.name} numberOfLines={1}>{community.name}</Text>
          {community.is_verified ? <Text style={styles.verified}>✓</Text> : null}
        </View>
        <View style={styles.metaRow}>
          <Badge variant="info" text={t(TYPE_LABELS[community.community_type])} />
          <Text style={styles.members}>{t('communities.members', { count: community.member_count })}</Text>
        </View>
        {community.district ? (
          <Text style={styles.district}>📍 {community.district}</Text>
        ) : null}
      </View>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  card: {
    flexDirection: 'row',
    backgroundColor: colors.card,
    borderRadius: radius.lg,
    borderWidth: 1,
    borderColor: colors.borderLight,
    padding: spacing.base,
    marginHorizontal: spacing.base,
    marginBottom: spacing.md,
    gap: spacing.md,
  },
  info: { flex: 1 },
  nameRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.xs },
  name: { ...typography.textStyles.h4, color: colors.text, flex: 1 },
  verified: { ...typography.textStyles.body, color: colors.primary, fontWeight: typography.fontWeight.bold },
  metaRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm, marginTop: spacing.xs },
  members: { ...typography.textStyles.caption, color: colors.textMuted },
  district: { ...typography.textStyles.caption, color: colors.textMuted, marginTop: 2 },
});
