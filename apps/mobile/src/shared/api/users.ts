import { api } from './client';

/* ── Types ── */
export interface UserProfile {
  id: string;
  phone: string;
  first_name: string;
  last_name: string;
  gender: 'male' | 'female';
  birth_year: number;
  city: string;
  district: string;
  language: string;
  avatar_url: string | null;
  ntrp_level: number;
  level_label: string;
  rating: number;
  rating_delta: number;
  total_games: number;
  wins: number;
  losses: number;
  win_rate: number;
  streak: number;
  is_profile_complete: boolean;
  last_active_at: string | null;
  created_at: string;
}

export interface UpdateProfileRequest {
  first_name?: string;
  last_name?: string;
  gender?: 'male' | 'female';
  birth_year?: number;
  district?: string;
  language?: string;
}

export interface PlayerSearchParams {
  query?: string;
  level_min?: number;
  level_max?: number;
  district?: string;
  gender?: string;
  sort_by?: 'rating' | 'activity' | 'name' | 'games';
  page?: number;
  per_page?: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
}

/* ── Helpers ── */
function unwrapData<T>(response: { data: { data: T } }): T {
  return response.data.data;
}

/* ── API ── */
export async function getMyProfile(): Promise<UserProfile> {
  const res = await api.get<{ data: UserProfile }>('/users/me');
  return unwrapData(res);
}

export async function updateMyProfile(data: UpdateProfileRequest): Promise<UserProfile> {
  const res = await api.patch<{ data: UserProfile }>('/users/me', data);
  return unwrapData(res);
}

export async function uploadAvatar(uri: string): Promise<{ avatar_url: string }> {
  const formData = new FormData();
  const name = uri.split('/').pop() ?? 'avatar.jpg';
  const type = name.endsWith('.png') ? 'image/png' : 'image/jpeg';
  formData.append('avatar', { uri, name, type } as unknown as Blob);
  const res = await api.post<{ data: { avatar_url: string } }>('/users/me/avatar', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
  return unwrapData(res);
}

export async function searchPlayers(
  params: PlayerSearchParams
): Promise<PaginatedResponse<UserProfile>> {
  const res = await api.get<{ data: PaginatedResponse<UserProfile> }>('/users/search', { params });
  return unwrapData(res);
}

export async function getUserById(id: string): Promise<UserProfile> {
  const res = await api.get<{ data: UserProfile }>(`/users/${id}`);
  return unwrapData(res);
}
