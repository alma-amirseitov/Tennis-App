import { View, Text, FlatList, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { Avatar, Badge } from '@/shared/ui';
import { colors, spacing, typography } from '@/shared/theme';
import type { CommunityMember } from '@/shared/api/communities';

const ROLE_COLORS: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'muted'> = {
  owner: 'warning',
  admin: 'danger' as 'primary',
  moderator: 'info',
  member: 'muted',
};

interface MembersListProps {
  members: CommunityMember[];
}

export function MembersList({ members }: MembersListProps) {
  const { t } = useTranslation();
  return (
    <FlatList
      data={members}
      keyExtractor={(item) => item.id}
      renderItem={({ item }) => (
        <View style={styles.row}>
          <Avatar uri={item.avatar_url} name={`${item.first_name} ${item.last_name}`} size="md" />
          <View style={styles.info}>
            <Text style={styles.name}>{item.first_name} {item.last_name}</Text>
            <Text style={styles.rating}>NTRP {item.ntrp_level} • {item.rating}</Text>
          </View>
          <Badge variant={ROLE_COLORS[item.role] ?? 'muted'} text={item.role} />
        </View>
      )}
      ItemSeparatorComponent={() => <View style={styles.separator} />}
      contentContainerStyle={styles.list}
    />
  );
}

const styles = StyleSheet.create({
  list: { paddingVertical: spacing.sm },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.base,
    gap: spacing.md,
  },
  info: { flex: 1 },
  name: { ...typography.textStyles.body, color: colors.text, fontWeight: typography.fontWeight.semibold },
  rating: { ...typography.textStyles.caption, color: colors.textMuted, marginTop: 1 },
  separator: { height: StyleSheet.hairlineWidth, backgroundColor: colors.borderLight, marginHorizontal: spacing.base },
});
