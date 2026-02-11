import { View, Text, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { Avatar, Badge, Button } from '@/shared/ui';
import { colors, spacing, typography } from '@/shared/theme';
import type { CommunityItem, CommunityType } from '@/shared/api/communities';

const TYPE_LABELS: Record<CommunityType, string> = {
  club: 'communities.type_club',
  league: 'communities.type_league',
  organizer: 'communities.type_organizer',
  group: 'communities.type_group',
};

interface CommunityHeaderProps {
  community: CommunityItem;
  onJoin?: () => void;
  onLeave?: () => void;
  loading?: boolean;
}

export function CommunityHeader({ community, onJoin, onLeave, loading }: CommunityHeaderProps) {
  const { t } = useTranslation();
  return (
    <View style={styles.container}>
      <Avatar uri={community.logo_url} name={community.name} size="xl" />
      <View style={styles.nameRow}>
        <Text style={styles.name}>{community.name}</Text>
        {community.is_verified ? <Text style={styles.verified}> ✓</Text> : null}
      </View>
      <View style={styles.metaRow}>
        <Badge variant="info" text={t(TYPE_LABELS[community.community_type])} />
        <Text style={styles.dot}>•</Text>
        <Text style={styles.meta}>{t('communities.members', { count: community.member_count })}</Text>
      </View>
      {community.district ? (
        <Text style={styles.location}>📍 {community.district}</Text>
      ) : null}
      <View style={styles.buttonRow}>
        {community.is_member ? (
          <Button variant="outline" title={t('communities.joined') + ' ✓'} onPress={onLeave ?? (() => {})} loading={loading} />
        ) : (
          <Button variant="primary" title={t('communities.join')} onPress={onJoin ?? (() => {})} loading={loading} />
        )}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { alignItems: 'center', paddingVertical: spacing.xl, paddingHorizontal: spacing.base },
  nameRow: { flexDirection: 'row', alignItems: 'center', marginTop: spacing.md },
  name: { ...typography.textStyles.h2, color: colors.text },
  verified: { ...typography.textStyles.h2, color: colors.primary },
  metaRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm, marginTop: spacing.sm },
  dot: { color: colors.textMuted },
  meta: { ...typography.textStyles.bodySm, color: colors.textSecondary },
  location: { ...typography.textStyles.bodySm, color: colors.textMuted, marginTop: spacing.xs },
  buttonRow: { marginTop: spacing.lg, width: '100%' },
});
