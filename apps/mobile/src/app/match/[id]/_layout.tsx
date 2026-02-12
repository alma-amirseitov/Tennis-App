import { Stack } from 'expo-router';

export default function MatchLayout() {
  return (
    <Stack screenOptions={{ headerShown: false }}>
      <Stack.Screen name="submit-result" />
      <Stack.Screen name="confirm-result" />
    </Stack>
  );
}
