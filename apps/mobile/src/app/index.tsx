import { useEffect, useState } from 'react';
import { View, ActivityIndicator, StyleSheet } from 'react-native';
import { Redirect } from 'expo-router';
import { authStore } from '@/shared/stores/auth';
import { colors } from '@/shared/theme';

export default function Index() {
  const [ready, setReady] = useState(false);
  const isAuthenticated = authStore((s) => s.isAuthenticated);
  const isLoading = authStore((s) => s.isLoading);

  useEffect(() => {
    authStore.getState().loadFromKeychain().finally(() => setReady(true));
  }, []);

  if (!ready || isLoading) {
    return (
      <View style={styles.splash}>
        <ActivityIndicator size="large" color={colors.primary} />
      </View>
    );
  }

  if (isAuthenticated) {
    return <Redirect href="/(tabs)" />;
  }

  return <Redirect href="/(auth)" />;
}

const styles = StyleSheet.create({
  splash: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: colors.background,
  },
});
