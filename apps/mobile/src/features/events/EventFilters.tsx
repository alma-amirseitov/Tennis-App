import { ScrollView, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';
import { Chip } from '@/shared/ui';
import { spacing } from '@/shared/theme';
import type { EventType } from '@/shared/api/events';

const TYPES: (EventType | 'all')[] = ['all', 'find_partner', 'organized_game', 'tournament', 'training'];

interface EventFiltersProps {
  selectedType: EventType | 'all';
  onTypeChange: (t: EventType | 'all') => void;
}

export function EventFilters({ selectedType, onTypeChange }: EventFiltersProps) {
  const { t } = useTranslation();

  const labels: Record<string, string> = {
    all: t('common.all'),
    find_partner: t('events.type_find_partner'),
    organized_game: t('events.type_organized_game'),
    tournament: t('events.type_tournament'),
    training: t('events.type_training'),
  };

  return (
    <ScrollView horizontal showsHorizontalScrollIndicator={false} contentContainerStyle={styles.row}>
      {TYPES.map((type) => (
        <Chip key={type} label={labels[type]} selected={selectedType === type} onPress={() => onTypeChange(type)} />
      ))}
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  row: { paddingHorizontal: spacing.base, paddingVertical: spacing.md },
});
