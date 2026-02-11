import { View, Text, StyleSheet } from 'react-native';
import { useTranslation } from 'react-i18next';

export default function QuizScreen() {
  const { t } = useTranslation();

  return (
    <View style={styles.container}>
      <Text style={styles.title}>{t('auth.quiz_title')}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#FAFBFC',
  },
  title: {
    fontSize: 18,
    fontWeight: '700',
    color: '#1A1D21',
  },
});
