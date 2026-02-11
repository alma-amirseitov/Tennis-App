import { useState } from 'react';
import { View, Text, ScrollView, Pressable, StyleSheet } from 'react-native';
import { useLocalSearchParams } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { useCommunity, useCommunityMembers, useJoinCommunity, useLeaveCommunity } from '@/shared/api/hooks';
import { CommunityHeader } from '@/features/communities/CommunityHeader';
import { MembersList } from '@/features/communities/MembersList';
import { ScreenHeader, ErrorState, Skeleton, EmptyState } from '@/shared/ui';
import { colors, spacing, typography, radius } from '@/shared/theme';

const TABS = ['events', 'members'] as const;
type TabKey = typeof TABS[number];

export default function CommunityDetailScreen() {
  const { t } = useTranslation();
  const { id } = useLocalSearchParams<{ id: string }>();
  const [activeTab, setActiveTab] = useState<TabKey>('events');

  const { data: community, isLoading, isError, refetch } = useCommunity(id ?? '');
  const { data: membersData } = useCommunityMembers(id ?? '', { per_page: 50 });
  const joinMutation = useJoinCommunity();
  const leaveMutation = useLeaveCommunity();

  const tabLabels: Record<TabKey, string> = {
    events: t('communities.tab_events'),
    members: t('communities.tab_members'),
  };

  if (isLoading) {
    return (
      <View style={styles.container}>
        <ScreenHeader title="" showBack />
        <View style={styles.skeleton}><Skeleton width={80} height={80} radius={40} /></View>
      </View>
    );
  }

  if (isError || !community) {
    return (
      <View style={styles.container}>
        <ScreenHeader title="" showBack />
        <ErrorState onRetry={refetch} />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <ScreenHeader title={community.name} showBack />
      <ScrollView showsVerticalScrollIndicator={false}>
        <CommunityHeader
          community={community}
          onJoin={() => joinMutation.mutate(community.id)}
          onLeave={() => leaveMutation.mutate(community.id)}
          loading={joinMutation.isPending || leaveMutation.isPending}
        />

        {/* Tab bar */}
        <View style={styles.tabBar}>
          {TABS.map((tab) => (
            <Pressable key={tab} onPress={() => setActiveTab(tab)} style={[styles.tab, activeTab === tab && styles.tabActive]}>
              <Text style={[styles.tabText, activeTab === tab && styles.tabTextActive]}>{tabLabels[tab]}</Text>
            </Pressable>
          ))}
        </View>

        {/* Tab content */}
        {activeTab === 'events' ? (
          <View style={styles.tabContent}>
            <EmptyState emoji="🎾" title={t('events.empty')} />
          </View>
        ) : (
          membersData ? (
            <MembersList members={membersData.data} />
          ) : (
            <View style={styles.tabContent}>
              <Skeleton width="100%" height={60} radius={8} />
            </View>
          )
        )}
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  skeleton: { alignItems: 'center', paddingTop: spacing['3xl'] },
  tabBar: {
    flexDirection: 'row',
    borderBottomWidth: 1,
    borderBottomColor: colors.borderLight,
    backgroundColor: colors.card,
  },
  tab: { flex: 1, paddingVertical: spacing.md, alignItems: 'center', borderBottomWidth: 2, borderBottomColor: 'transparent' },
  tabActive: { borderBottomColor: colors.primary },
  tabText: { ...typography.textStyles.bodySm, color: colors.textMuted },
  tabTextActive: { color: colors.primary, fontWeight: typography.fontWeight.semibold },
  tabContent: { minHeight: 200 },
});
