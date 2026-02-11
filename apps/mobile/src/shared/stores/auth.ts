import { Platform } from 'react-native';
import { create } from 'zustand';

const AUTH_KEY = 'auth';

// Platform-safe storage abstraction
const storage = {
  async setItem(key: string, value: string): Promise<void> {
    if (Platform.OS === 'web') {
      localStorage.setItem(key, value);
    } else {
      const SecureStore = await import('expo-secure-store');
      await SecureStore.setItemAsync(key, value);
    }
  },
  async getItem(key: string): Promise<string | null> {
    if (Platform.OS === 'web') {
      return localStorage.getItem(key);
    }
    const SecureStore = await import('expo-secure-store');
    return SecureStore.getItemAsync(key);
  },
  async removeItem(key: string): Promise<void> {
    if (Platform.OS === 'web') {
      localStorage.removeItem(key);
    } else {
      const SecureStore = await import('expo-secure-store');
      await SecureStore.deleteItemAsync(key);
    }
  },
};

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
}

export interface AuthUser {
  id: string;
  first_name: string;
  is_profile_complete?: boolean;
}

interface AuthState {
  isAuthenticated: boolean;
  isLoading: boolean;
  accessToken: string | null;
  refreshToken: string | null;
  user: AuthUser | null;

  login: (tokens: AuthTokens, user?: AuthUser) => Promise<void>;
  setTempToken: (tempToken: string, userId: string) => Promise<void>;
  logout: () => Promise<void>;
  loadFromKeychain: () => Promise<void>;
}

export const authStore = create<AuthState>((set) => ({
  isAuthenticated: false,
  isLoading: true,
  accessToken: null,
  refreshToken: null,
  user: null,

  login: async (tokens, user) => {
    await storage.setItem(
      AUTH_KEY,
      JSON.stringify({
        accessToken: tokens.access_token,
        refreshToken: tokens.refresh_token,
        user: user ?? null,
      })
    );
    set({
      isAuthenticated: true,
      accessToken: tokens.access_token,
      refreshToken: tokens.refresh_token,
      user: user ?? null,
    });
  },

  setTempToken: async (tempToken, userId) => {
    await storage.setItem(
      AUTH_KEY,
      JSON.stringify({
        tempToken,
        userId,
        user: null,
      })
    );
    set({
      isAuthenticated: false,
      accessToken: tempToken,
      refreshToken: null,
      user: { id: userId, first_name: '', is_profile_complete: false },
    });
  },

  logout: async () => {
    await storage.removeItem(AUTH_KEY);
    set({
      isAuthenticated: false,
      accessToken: null,
      refreshToken: null,
      user: null,
    });
  },

  loadFromKeychain: async () => {
    set({ isLoading: true });
    try {
      const stored = await storage.getItem(AUTH_KEY);
      if (stored) {
        const data = JSON.parse(stored) as {
          accessToken?: string;
          refreshToken?: string;
          tempToken?: string;
          userId?: string;
          user?: AuthUser | null;
        };
        if (data.accessToken && data.refreshToken) {
          set({
            isAuthenticated: true,
            accessToken: data.accessToken,
            refreshToken: data.refreshToken,
            user: data.user ?? null,
            isLoading: false,
          });
          return;
        }
        if (data.tempToken && data.userId) {
          set({
            isAuthenticated: false,
            accessToken: data.tempToken,
            refreshToken: null,
            user: data.user ?? { id: data.userId, first_name: '', is_profile_complete: false },
            isLoading: false,
          });
          return;
        }
      }
      set({ isLoading: false });
    } catch {
      set({ isLoading: false });
    }
  },
}));
