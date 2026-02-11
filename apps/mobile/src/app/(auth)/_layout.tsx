import { Stack } from 'expo-router';

export default function AuthLayout() {
  return (
    <Stack
      screenOptions={{
        headerShown: false,
        animation: 'slide_from_right',
        gestureEnabled: true,
      }}
    >
      <Stack.Screen name="index" />
      <Stack.Screen name="otp" />
      <Stack.Screen name="profile-setup" options={{ gestureEnabled: false }} />
      <Stack.Screen name="quiz" options={{ gestureEnabled: false }} />
    </Stack>
  );
}
