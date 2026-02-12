import { useState, useMemo } from 'react';
import {
  View,
  Text,
  ScrollView,
  StyleSheet,
  TextInput,
  Pressable,
  Alert,
  Platform,
} from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { useMatch, useSubmitMatchResult } from '@/shared/api/hooks';
import { ScreenHeader, Button, Avatar, Skeleton, ErrorState } from '@/shared/ui';
import { showToast } from '@/shared/lib/toast';
import { colors, spacing, typography, radius } from '@/shared/theme';
import { authStore } from '@/shared/stores/auth';
import type { MatchSetScore } from '@/shared/api/matches';

function formatSet(score: MatchSetScore): string {
  let s = `${score.p1}:${score.p2}`;
  if (score.tiebreak) s += ` (${score.tiebreak.p1}:${score.tiebreak.p2})`;
  return s;
}

export default function SubmitResultScreen() {
  const { t } = useTranslation();
  const router = useRouter();
  const { id } = useLocalSearchParams<{ id: string }>();
  const matchId = id ?? '';
  const myId = authStore((s) => s.user?.id) ?? '';

  const { data: match, isLoading, isError, refetch } = useMatch(matchId);
  const submitMutation = useSubmitMatchResult();

  const [setsCount, setSetsCount] = useState<2 | 3>(2);
  const [scores, setScores] = useState<MatchSetScore[]>([
    { p1: 0, p2: 0 },
    { p1: 0, p2: 0 },
  ]);
  const [showPreview, setShowPreview] = useState(false);
  const [winnerOverride, setWinnerOverride] = useState<string | null>(null);

  const p1Id = match?.player1_id ?? '';
  const p2Id = match?.player2_id ?? '';
  const p1Name = match ? `${match.player1.first_name} ${match.player1.last_name[0]}.` : '';
  const p2Name = match ? `${match.player2.first_name} ${match.player2.last_name[0]}.` : '';

  const computedWinner = useMemo(() => {
    let p1Wins = 0;
    let p2Wins = 0;
    for (const s of scores) {
      if (s.p1 > s.p2) p1Wins++;
      else if (s.p2 > s.p1) p2Wins++;
    }
    if (p1Wins > p2Wins) return p1Id;
    if (p2Wins > p1Wins) return p2Id;
    return null;
  }, [scores, p1Id, p2Id]);

  const winnerId = winnerOverride ?? computedWinner;

  const updateScore = (setIndex: number, field: 'p1' | 'p2', value: number) => {
    setScores((prev) => {
      const next = [...prev];
      next[setIndex] = { ...next[setIndex], [field]: Math.max(0, Math.min(20, value)) };
      return next;
    });
  };

  const setTiebreak = (setIndex: number, p1: number, p2: number) => {
    setScores((prev) => {
      const next = [...prev];
      next[setIndex] = {
        ...next[setIndex],
        tiebreak: { p1: Math.max(0, Math.min(20, p1)), p2: Math.max(0, Math.min(20, p2)) },
      };
      return next;
    });
  };

  const changeSetsCount = (n: 2 | 3) => {
    if (n === 3 && scores.length < 3) {
      setScores((prev) => [...prev, { p1: 0, p2: 0 }]);
    } else if (n === 2 && scores.length > 2) {
      setScores((prev) => prev.slice(0, 2));
    }
    setSetsCount(n);
  };

  const needsTiebreak = (s: MatchSetScore) =>
    (s.p1 === 7 && s.p2 === 6) || (s.p1 === 6 && s.p2 === 7);
  const hasTiebreak = (s: MatchSetScore) => !!s.tiebreak;

  const handleSubmit = async () => {
    if (!winnerId || !matchId) return;
    if (scores.some((s) => s.p1 === 0 && s.p2 === 0)) {
      showToast(t('errors.something_went_wrong'));
      return;
    }

    try {
      await submitMutation.mutateAsync({
        matchId,
        data: {
          winner_id: winnerId,
          score: { sets: scores },
        },
      });
      showToast(t('match.result_submitted'));
      router.back();
    } catch {
      showToast(t('errors.something_went_wrong'));
    }
  };

  if (isLoading) {
    return (
      <View style={styles.container}>
        <ScreenHeader title={t('match.submit_result')} showBack />
        <View style={styles.skeleton}><Skeleton width="100%" height={200} radius={16} /></View>
      </View>
    );
  }

  if (isError || !match) {
    return (
      <View style={styles.container}>
        <ScreenHeader title={t('match.submit_result')} showBack />
        <ErrorState onRetry={refetch} />
      </View>
    );
  }

  if (showPreview) {
    return (
      <View style={styles.container}>
        <ScreenHeader title={t('match.submit_result')} showBack />
        <ScrollView contentContainerStyle={styles.content}>
          <Text style={styles.previewTitle}>{t('match.score')}</Text>
          <View style={styles.vsRow}>
            <Avatar uri={match.player1.avatar_url} name={p1Name} size="lg" />
            <Text style={styles.vs}>vs</Text>
            <Avatar uri={match.player2.avatar_url} name={p2Name} size="lg" />
          </View>
          <View style={styles.scorePreview}>
            {scores.map((s, i) => (
              <Text key={i} style={styles.scoreLine}>
                {t('match.set', { number: i + 1 })}: {formatSet(s)}
              </Text>
            ))}
          </View>
          <Text style={styles.winnerLabel}>{t('match.winner')}</Text>
          <View style={styles.winnerToggle}>
            <Pressable
              onPress={() => setWinnerOverride(p1Id)}
              style={[styles.winnerBtn, winnerId === p1Id && styles.winnerBtnActive]}
            >
              <Text style={[styles.winnerText, winnerId === p1Id && styles.winnerTextActive]}>{p1Name}</Text>
            </Pressable>
            <Pressable
              onPress={() => setWinnerOverride(p2Id)}
              style={[styles.winnerBtn, winnerId === p2Id && styles.winnerBtnActive]}
            >
              <Text style={[styles.winnerText, winnerId === p2Id && styles.winnerTextActive]}>{p2Name}</Text>
            </Pressable>
          </View>
          <View style={styles.actions}>
            <Button variant="outline" title={t('common.back')} onPress={() => setShowPreview(false)} style={styles.actionBtn} />
            <Button
              variant="primary"
              title={t('match.submit_result')}
              onPress={handleSubmit}
              loading={submitMutation.isPending}
              disabled={!winnerId}
              style={styles.actionBtn}
            />
          </View>
        </ScrollView>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <ScreenHeader title={t('match.submit_result')} showBack />
      <ScrollView contentContainerStyle={styles.content} keyboardShouldPersistTaps="handled">
        <View style={styles.vsRow}>
          <View style={styles.playerCol}>
            <Avatar uri={match.player1.avatar_url} name={p1Name} size="md" />
            <Text style={styles.playerName}>{p1Name}</Text>
            <Text style={styles.ntrp}>({match.player1.ntrp_level})</Text>
          </View>
          <Text style={styles.vs}>vs</Text>
          <View style={styles.playerCol}>
            <Avatar uri={match.player2.avatar_url} name={p2Name} size="md" />
            <Text style={styles.playerName}>{p2Name}</Text>
            <Text style={styles.ntrp}>({match.player2.ntrp_level})</Text>
          </View>
        </View>

        {/* Sets count */}
        <Text style={styles.sectionLabel}>{t('match.set', { number: 1 })} — {setsCount === 2 ? 'Best of 2' : 'Best of 3'}</Text>
        <View style={styles.setsToggle}>
          <Pressable onPress={() => changeSetsCount(2)} style={[styles.setBtn, setsCount === 2 && styles.setBtnActive]}>
            <Text style={[styles.setBtnText, setsCount === 2 && styles.setBtnTextActive]}>2</Text>
          </Pressable>
          <Pressable onPress={() => changeSetsCount(3)} style={[styles.setBtn, setsCount === 3 && styles.setBtnActive]}>
            <Text style={[styles.setBtnText, setsCount === 3 && styles.setBtnTextActive]}>3</Text>
          </Pressable>
        </View>

        {/* Set inputs */}
        {scores.slice(0, setsCount).map((s, i) => (
          <View key={i} style={styles.setRow}>
            <Text style={styles.setLabel}>{t('match.set', { number: i + 1 })}</Text>
            <View style={styles.scoreInputs}>
              <TextInput
                style={styles.scoreInput}
                value={String(s.p1)}
                onChangeText={(v) => updateScore(i, 'p1', parseInt(v, 10) || 0)}
                keyboardType="number-pad"
                maxLength={2}
              />
              <Text style={styles.colon}>:</Text>
              <TextInput
                style={styles.scoreInput}
                value={String(s.p2)}
                onChangeText={(v) => updateScore(i, 'p2', parseInt(v, 10) || 0)}
                keyboardType="number-pad"
                maxLength={2}
              />
            </View>
            {needsTiebreak(s) && (
              <View style={styles.tiebreakRow}>
                <Text style={styles.tiebreakLabel}>{t('match.tiebreak')}</Text>
                <TextInput
                  style={styles.tiebreakInput}
                  value={String(s.tiebreak?.p1 ?? 0)}
                  onChangeText={(v) => setTiebreak(i, parseInt(v, 10) || 0, s.tiebreak?.p2 ?? 0)}
                  keyboardType="number-pad"
                  maxLength={2}
                />
                <Text style={styles.colon}>:</Text>
                <TextInput
                  style={styles.tiebreakInput}
                  value={String(s.tiebreak?.p2 ?? 0)}
                  onChangeText={(v) => setTiebreak(i, s.tiebreak?.p1 ?? 0, parseInt(v, 10) || 0)}
                  keyboardType="number-pad"
                  maxLength={2}
                />
              </View>
            )}
          </View>
        ))}

        <Button
          variant="primary"
          title={t('common.confirm')}
          onPress={() => setShowPreview(true)}
          disabled={!computedWinner || scores.slice(0, setsCount).some((s) => s.p1 === 0 && s.p2 === 0)}
          style={styles.nextBtn}
        />
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.background },
  content: { padding: spacing.xl, paddingBottom: spacing['3xl'] },
  skeleton: { padding: spacing.xl },
  vsRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'space-around', marginBottom: spacing.xl },
  playerCol: { alignItems: 'center' },
  playerName: { ...typography.textStyles.body, color: colors.text, marginTop: spacing.sm },
  ntrp: { ...typography.textStyles.caption, color: colors.textMuted },
  vs: { ...typography.textStyles.h4, color: colors.textMuted },
  sectionLabel: { ...typography.textStyles.bodySm, color: colors.textSecondary, marginBottom: spacing.sm },
  setsToggle: { flexDirection: 'row', gap: spacing.sm, marginBottom: spacing.xl },
  setBtn: {
    flex: 1,
    paddingVertical: spacing.md,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    alignItems: 'center',
  },
  setBtnActive: { borderColor: colors.primary, backgroundColor: colors.primaryLight },
  setBtnText: { ...typography.textStyles.body, color: colors.text },
  setBtnTextActive: { color: colors.primary, fontWeight: typography.fontWeight.semibold },
  setRow: { marginBottom: spacing.lg },
  setLabel: { ...typography.textStyles.bodySm, color: colors.textSecondary, marginBottom: spacing.xs },
  scoreInputs: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm },
  scoreInput: {
    width: 56,
    height: 48,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.md,
    ...typography.textStyles.h3,
    color: colors.text,
    textAlign: 'center',
    padding: 0,
  },
  tiebreakRow: { flexDirection: 'row', alignItems: 'center', marginTop: spacing.sm, gap: spacing.sm },
  tiebreakLabel: { ...typography.textStyles.caption, color: colors.textMuted, width: 60 },
  tiebreakInput: {
    width: 48,
    height: 40,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.sm,
    ...typography.textStyles.body,
    color: colors.text,
    textAlign: 'center',
    padding: 0,
  },
  colon: { ...typography.textStyles.h3, color: colors.textMuted },
  nextBtn: { marginTop: spacing.xl },
  previewTitle: { ...typography.textStyles.h3, color: colors.text, marginBottom: spacing.lg, textAlign: 'center' },
  scorePreview: { marginBottom: spacing.xl },
  scoreLine: { ...typography.textStyles.body, color: colors.text, marginBottom: spacing.sm },
  winnerLabel: { ...typography.textStyles.bodySm, color: colors.textSecondary, marginBottom: spacing.sm },
  winnerToggle: { flexDirection: 'row', gap: spacing.sm, marginBottom: spacing.xl },
  winnerBtn: {
    flex: 1,
    paddingVertical: spacing.md,
    borderRadius: radius.md,
    borderWidth: 1,
    borderColor: colors.border,
    alignItems: 'center',
  },
  winnerBtnActive: { borderColor: colors.primary, backgroundColor: colors.primaryLight },
  winnerText: { ...typography.textStyles.body, color: colors.text },
  winnerTextActive: { color: colors.primary, fontWeight: typography.fontWeight.semibold },
  actions: { flexDirection: 'row', gap: spacing.sm },
  actionBtn: { flex: 1 },
});
