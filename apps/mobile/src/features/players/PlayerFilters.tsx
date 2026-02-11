import { View, ScrollView, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { Chip } from '@/shared/ui';
import { spacing } from '@/shared/theme';

interface PlayerFiltersProps {
  gender: string;
  onGenderChange: (v: string) => void;
}

const GENDERS = ['all', 'male', 'female'] as const;

export function PlayerFilters({ gender, onGenderChange }: PlayerFiltersProps) {
  const { t } = useTranslation();

  const genderLabels: Record<string, string> = {
    all: t('common.all'),
    male: t('auth.gender_male'),
    female: t('auth.gender_female'),
  };

  return (
    <View style={styles.container}>
      <ScrollView horizontal showsHorizontalScrollIndicator={false} contentContainerStyle={styles.row}>
        {GENDERS.map((g) => (
          <Chip key={g} label={genderLabels[g]} selected={gender === g} onPress={() => onGenderChange(g)} />
        ))}
      </ScrollView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { paddingVertical: spacing.sm },
  row: { paddingHorizontal: spacing.base },
});
