import { useState, useEffect, useRef } from 'react';
import {
  View,
  Text,
  ScrollView,
  StyleSheet,
  Animated,
  Pressable,
  Alert,
  Platform,
} from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { useMatch, useConfirmMatchResult } from '@/shared/api/hooks';
import { ScreenHeader, Button, Avatar, Skeleton, ErrorState } from '@/shared/ui';
import { showToast } from '@/shared/lib/toast';
import { colors, spacing, typography, radius } from '@/shared/theme';
import { authStore } from '@/shared/stores/auth';

function formatSet(s: { p1: number; p2: number; tiebreak?: { p1: number; p2: number } }): string {
  let str = `${s.p1}:${s.p2}`;
  if (s.tiebreak) str += ` (${s.tiebreak.p1}:${s.tiebreak.p2})`;
  return str;
}

export default function ConfirmResultScreen() {
  const { t } = useTranslation();
  const router = useRouter();
  const { id } = useLocalSearchParams<{ id: string }>();
  const matchId = id ?? '';
  const myId = authStore((s) => s.user?.id) ?? '';

  const { data: match, isLoading, isError, refetch } = useMatch(matchId);
  const confirmMutation = useConfirmMatchResult();

  const [confirmed, setConfirmed] = useState<boolean | null>(null);
  const [ratingDelta, setRatingDelta] = useState<number | null>(null);
  const animValue = useRef(new Animated.Value(0)).current;

  const opponent = match
    ? match.player1_id === myId
      ? match.player2
      : match.player1
    : null;
  const opponentName = opponent
    ? `${opponent.first_name} ${opponent.last_name}`
    : '';
  const score = match?.score;
  const winnerId = match?.winner_id ?? null;
  const winnerName =
    match && winnerId
      ? winnerId === match.player1_id
        ? `${match.player1.first_name} ${match.player1.last_name}`
        : `${match.player2.first_name} ${match.player2.last_name}`
      : '';

  const handleConfirm = async () => {
    try {
      const res = await confirmMutation.mutateAsync({
        matchId,
        data: { action: 'confirm' },
      });
      setConfirmed(true);
      const rc = res.rating_changes;
      const myChange = rc
        ? myId === match?.player1_id
          ? rc.player1.change
          : rc.player2.change
        : 0;
      setRatingDelta(myChange);
      showToast(t('match.result_confirmed'));
    } catch {
      showToast(t('errors.something_went_wrong'));
    }
  };

  const handleDispute = () => {
    if (Platform.OS === 'web') {
      if (window.confirm(t('match.dispute_result') + '?')) {
        confirmMutation.mutate({
          matchId,
          data: { action: 'dispute', reason: 'Disputed by user' },
        });
        showToast(t('match.dispute_result'));
        router.back();
      }
    } else {
      Alert.alert(
        t('match.dispute_result'),
        t('match.dispute_result') + '?',
        [
          { text: t('common.cancel'), style: 'cancel' },
          {
            text: t('match.dispute_result'),
            style: 'destructive',
            onPress: () => {
              confirmMutation.mutate({
                matchId,
                data: { action: 'dispute', reason: 'Disputed by user' },
              });
              showToast(t('match.dispute_result'));
              router.back();
            },
          },
        ]
      );
    }
  };

  useEffect(() => {
    if (ratingDelta !== null) {
      animValue.setValue(0);
      Animated.spring(animValue, {
        toValue: 1,
        useNativeDriver: true,
        friction: 6,
        tension: 100,
      }).start();
    }
  }, [ratingDelta, animValue]);

  const scale = animValue.interpolate({ inputRange: [0, 1], outputRange: [0.5, 1] });
  const opacity = animValue.interpolate({ inputRange: [0, 1], outputRange: [0, 1] });

  if (isLoading) {
    return (
      <View style={styles.container}>
        <ScreenHeader title={t('match.confirm_result')} showBack />
        <View style={styles.skeleton}><Skeleton width="100%" height={200} radius={16} /></View>
      </View>
    );
  }

  if (isError || !match) {
    return (
      <View style={styles.container}>
        <ScreenHeader title={t('match.confirm_result')} showBack />
        <ErrorState onRetry={refetch} />
      </View>
    );
  }

  if (confirmed && ratingDelta !== null) {
    const isPositive = ratingDelta >= 0;
    return (
      <View style={styles.container}>
        <ScreenHeader title={t('match.result_confirmed')} showBack />
        <View style={styles.successContent}>
          <Text style={styles.successEmoji}>🎾</Text>
          <Text style={styles.successTitle}>{t('match.result_confirmed')}</Text>
          <Animated.View style={[styles.ratingChange, { transform: [{ scale }], opacity }]}>
            <Text style={[styles.ratingDelta, { color: isPositive ? colors.success : colors.danger }]}>
              {isPositive ? '+' : ''}{Math.round(ratingDelta)}
            </Text>
          </Animated.View>
          <Button
            variant="primary"
            title={t('common.done')}
            onPress={() => router.back()}
            style={styles.doneBtn}
          />
        </View>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <ScreenHeader title={t('match.confirm_result')} showBack />
      <ScrollView contentContainerStyle={styles.content}>
        <Text style={styles.subtitle}>
          {t('notifications.match_result_pending')}
        </Text>

        <View style={styles.opponentRow}>
          <Avatar uri={opponent?.avatar_url} name={opponentName} size="xl" />
          <Text style={styles.opponentName}>{opponentName}</Text>
        </View>

        {score && score.sets.length > 0 && (
          <>
            <Text style={styles.sectionTitle}>{t('match.score')}</Text>
            <View style={styles.scoreCard}>
              {score.sets.map((s, i) => (
                <Text key={i} style={styles.scoreLine}>
                  {t('match.set', { number: i + 1 })}: {formatSet(s)}
                </Text>
              ))}
            </View>
          </>
        )}

        {winnerName ? (
          <Text style={styles.winnerText}>
            {t('match.winner')}: {winnerName}
          </Text>
        ) : null}

        <View style={styles.actions}>
          <Button
            variant="primary"
            title={t('match.confirm_result') + ' ✓'}
            onPress={handleConfirm}
            loading={confirmMutation.isPending}
            style={styles.actionBtn}
          />
          <Button
            variant="outline"
            title={t('match.dispute_result') + ' ✕'}
            onPress={handleDispute}
            disabled={confirmMutation.isPending}
            style={[styles.actionBtn, styles.disputeBtn]}
          />
        </View>
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  content: { padding: spacing.xl, paddingBottom: spacing['3xl'] },
  skeleton: { padding: spacing.xl },
  subtitle: {
    ...typography.textStyles.body,
    color: colors.textSecondary,
    textAlign: 'center',
    marginBottom: spacing.xl,
  },
  opponentRow: { alignItems: 'center', marginBottom: spacing.xl },
  opponentName: { ...typography.textStyles.h3, color: colors.text, marginTop: spacing.md },
  sectionTitle: { ...typography.textStyles.h4, color: colors.text, marginBottom: spacing.md },
  scoreCard: {
    backgroundColor: colors.card,
    borderRadius: radius.lg,
    borderWidth: 1,
    borderColor: colors.borderLight,
    padding: spacing.base,
    marginBottom: spacing.xl,
  },
  scoreLine: { ...typography.textStyles.body, color: colors.text, marginBottom: spacing.sm },
  winnerText: {
    ...typography.textStyles.body,
    color: colors.text,
    textAlign: 'center',
    marginBottom: spacing.xl,
  },
  actions: { gap: spacing.md },
  actionBtn: {},
  disputeBtn: { borderColor: colors.danger },
  successContent: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: spacing.xl,
  },
  successEmoji: { fontSize: 64, marginBottom: spacing.lg },
  successTitle: { ...typography.textStyles.h2, color: colors.text, marginBottom: spacing.xl, textAlign: 'center' },
  ratingChange: { marginBottom: spacing.xl },
  ratingDelta: { fontSize: 48, fontWeight: typography.fontWeight.extrabold },
  doneBtn: { minWidth: 200 },
});
