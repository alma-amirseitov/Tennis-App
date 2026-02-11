import { useState, useEffect } from 'react';
import { View, Text, ScrollView, StyleSheet, Pressable } from 'react-native';
import { useRouter } from 'expo-router';
import { useTranslation } from 'react-i18next';
import { Button, Card, Skeleton } from '@/shared/ui';
import { getQuizQuestions, submitQuiz } from '@/shared/api/quiz';
import type { QuizQuestion, QuizOption } from '@/shared/api/quiz';
import { showToast } from '@/shared/lib/toast';
import { colors, spacing, typography, radius } from '@/shared/theme';
import axios from 'axios';

function getQuestionText(q: QuizQuestion, lang: string): string {
  if (lang === 'kk' && q.text_kz) return q.text_kz;
  if (lang === 'en' && q.text_en) return q.text_en;
  return q.text_ru;
}

function getOptionText(o: QuizOption, lang: string): string {
  if (lang === 'kk' && o.text_kz) return o.text_kz;
  if (lang === 'en' && o.text_en) return o.text_en;
  return o.text_ru;
}

export default function QuizScreen() {
  const { t, i18n } = useTranslation();
  const router = useRouter();
  const lang = i18n.language;

  const [questions, setQuestions] = useState<QuizQuestion[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [answers, setAnswers] = useState<Array<{ question_id: number; option_id: string }>>([]);
  const [submitting, setSubmitting] = useState(false);
  const [result, setResult] = useState<{
    level_label: string;
    ntrp_level: number;
    initial_rating: number;
  } | null>(null);

  useEffect(() => {
    getQuizQuestions()
      .then((res) => setQuestions(res.questions ?? []))
      .catch(() => {
        showToast(t('errors.something_went_wrong'));
        setQuestions([]);
      })
      .finally(() => setLoading(false));
  }, [t]);

  const currentQuestion = questions[currentIndex];
  const isLastQuestion = currentIndex === questions.length - 1;
  const selectedOptionId = answers.find(
    (a) => a.question_id === currentQuestion?.id
  )?.option_id;

  const handleSelectOption = (optionId: string) => {
    if (!currentQuestion) return;
    const newAnswers = answers.filter((a) => a.question_id !== currentQuestion.id);
    newAnswers.push({ question_id: currentQuestion.id, option_id: optionId });
    setAnswers(newAnswers);
  };

  const handleNext = async () => {
    if (!currentQuestion) return;
    const selected = answers.find((a) => a.question_id === currentQuestion.id);
    if (!selected) return;

    if (isLastQuestion) {
      setSubmitting(true);
      try {
        const res = await submitQuiz(answers);
        setResult({
          level_label: res.level_label,
          ntrp_level: res.ntrp_level,
          initial_rating: res.initial_rating,
        });
      } catch (err) {
        if (axios.isAxiosError(err)) {
          showToast(err.response?.data?.error?.message ?? t('errors.something_went_wrong'));
        } else {
          showToast(t('errors.something_went_wrong'));
        }
      } finally {
        setSubmitting(false);
      }
    } else {
      setCurrentIndex((i) => i + 1);
    }
  };

  const handleSkip = () => {
    router.replace('/(tabs)');
  };

  const handleStartPlaying = () => {
    router.replace('/(tabs)');
  };

  if (loading) {
    return (
      <View style={styles.container}>
        <View style={styles.skeletonTitle}>
          <Skeleton width="100%" height={24} />
        </View>
        <View style={styles.skeletonSubtitle}>
          <Skeleton width="80%" height={16} />
        </View>
        <View style={styles.skeletonCard}>
          <Skeleton width="100%" height={56} />
        </View>
        <View style={styles.skeletonCard}>
          <Skeleton width="100%" height={56} />
        </View>
      </View>
    );
  }

  if (questions.length === 0 && !result) {
    return (
      <View style={styles.container}>
        <Text style={styles.title}>{t('auth.quiz_title')}</Text>
        <Text style={styles.emptyText}>{t('common.no_results')}</Text>
        <Button
          variant="primary"
          title={t('auth.quiz_skip')}
          onPress={handleSkip}
          style={styles.skipBtn}
        />
      </View>
    );
  }

  if (result) {
    return (
      <View style={styles.container}>
        <Text style={styles.resultEmoji}>🎾</Text>
        <Text style={styles.resultTitle}>{t('auth.quiz_result_title')}</Text>
        <Card style={styles.resultCard}>
          <Text style={styles.resultLevel}>{result.level_label}</Text>
          <Text style={styles.resultNtrp}>
            {t('auth.quiz_result_subtitle', {
              level: result.ntrp_level,
              name: result.level_label,
            })}
          </Text>
        </Card>
        <Text style={styles.resultDesc}>{t('auth.quiz_result_desc')}</Text>
        <Button
          variant="primary"
          title={t('auth.quiz_start_playing')}
          onPress={handleStartPlaying}
          style={styles.startBtn}
        />
      </View>
    );
  }

  if (!currentQuestion) return null;

  const progress = ((currentIndex + 1) / questions.length) * 100;

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Pressable onPress={handleSkip}>
          <Text style={styles.skipLink}>{t('auth.quiz_skip')}</Text>
        </Pressable>
        <Text style={styles.progress}>
          {t('auth.quiz_progress', {
            current: currentIndex + 1,
            total: questions.length,
          })}
        </Text>
      </View>

      <View style={styles.progressBar}>
        <View style={[styles.progressFill, { width: `${progress}%` }]} />
      </View>

      <Text style={styles.title}>{t('auth.quiz_title')}</Text>
      <Text style={styles.subtitle}>{t('auth.quiz_subtitle')}</Text>

      <Text style={styles.questionText}>
        {getQuestionText(currentQuestion, lang)}
      </Text>

      <ScrollView style={styles.options} showsVerticalScrollIndicator={false}>
        {currentQuestion.options.map((opt) => (
          <Card
            key={opt.id}
            onPress={() => handleSelectOption(opt.id)}
            style={[
              styles.optionCard,
              selectedOptionId === opt.id && styles.optionCardSelected,
            ]}
          >
            <Text
              style={[
                styles.optionText,
                selectedOptionId === opt.id && styles.optionTextSelected,
              ]}
            >
              {getOptionText(opt, lang)}
            </Text>
          </Card>
        ))}
      </ScrollView>

      <Button
        variant="primary"
        title={t('auth.quiz_next')}
        onPress={handleNext}
        disabled={!selectedOptionId}
        loading={submitting}
        style={styles.nextBtn}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
    padding: spacing.xl,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing.lg,
  },
  skipLink: {
    ...typography.textStyles.body,
    color: colors.textMuted,
  },
  progress: {
    ...typography.textStyles.caption,
    color: colors.textSecondary,
  },
  progressBar: {
    height: 4,
    backgroundColor: colors.borderLight,
    borderRadius: 2,
    overflow: 'hidden',
    marginBottom: spacing.xl,
  },
  progressFill: {
    height: '100%',
    backgroundColor: colors.primary,
  },
  title: {
    ...typography.textStyles.h2,
    color: colors.text,
    marginBottom: spacing.sm,
  },
  subtitle: {
    ...typography.textStyles.bodySm,
    color: colors.textSecondary,
    marginBottom: spacing.xl,
  },
  questionText: {
    ...typography.textStyles.h3,
    color: colors.text,
    marginBottom: spacing.xl,
  },
  options: {
    flex: 1,
    marginBottom: spacing.lg,
  },
  optionCard: {
    marginBottom: spacing.sm,
    borderWidth: 2,
    borderColor: colors.border,
  },
  optionCardSelected: {
    borderColor: colors.primary,
    backgroundColor: colors.primaryLight,
  },
  optionText: {
    ...typography.textStyles.body,
    color: colors.text,
  },
  optionTextSelected: {
    color: colors.primary,
    fontWeight: typography.fontWeight.semibold,
  },
  nextBtn: {
    width: '100%',
  },
  skeletonTitle: {
    marginBottom: spacing.sm,
  },
  skeletonSubtitle: {
    marginBottom: spacing.xl,
  },
  skeletonCard: {
    marginBottom: spacing.sm,
  },
  emptyText: {
    ...typography.textStyles.body,
    color: colors.textMuted,
    marginBottom: spacing.xl,
  },
  skipBtn: {
    marginTop: spacing.lg,
  },
  resultEmoji: {
    fontSize: 48,
    textAlign: 'center',
    marginBottom: spacing.lg,
  },
  resultTitle: {
    ...typography.textStyles.h2,
    color: colors.text,
    textAlign: 'center',
    marginBottom: spacing.xl,
  },
  resultCard: {
    marginBottom: spacing.xl,
    alignItems: 'center',
  },
  resultLevel: {
    ...typography.textStyles.h2,
    color: colors.primary,
    marginBottom: spacing.sm,
  },
  resultNtrp: {
    ...typography.textStyles.body,
    color: colors.textSecondary,
  },
  resultDesc: {
    ...typography.textStyles.bodySm,
    color: colors.textMuted,
    textAlign: 'center',
    marginBottom: spacing.xl,
  },
  startBtn: {
    width: '100%',
  },
});
