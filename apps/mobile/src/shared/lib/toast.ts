import { Alert, Platform } from 'react-native';

export function showToast(message: string): void {
  if (Platform.OS === 'web') {
    // On web, use window.alert as fallback
    if (typeof window !== 'undefined' && window.alert) {
      window.alert(message);
    }
  } else {
    Alert.alert('', message);
  }
}
