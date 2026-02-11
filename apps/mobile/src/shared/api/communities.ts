import { api } from './client';
import type { PaginatedResponse } from './users';

/* ── Types ── */
export type CommunityType = 'club' | 'league' | 'organizer' | 'group';
export type CommunityAccess = 'open' | 'closed' | 'paid';
export type MemberRole = 'owner' | 'admin' | 'moderator' | 'member';

export interface CommunityItem {
  id: string;
  name: string;
  description: string;
  community_type: CommunityType;
  access_type: CommunityAccess;
  logo_url: string | null;
  cover_url: string | null;
  district: string;
  is_verified: boolean;
  member_count: number;
  active_events_count: number;
  is_member: boolean;
  my_role: MemberRole | null;
  created_at: string;
}

export interface CommunityMember {
  id: string;
  user_id: string;
  first_name: string;
  last_name: string;
  avatar_url: string | null;
  ntrp_level: number;
  rating: number;
  role: MemberRole;
  joined_at: string;
}

export interface CommunitySearchParams {
  query?: string;
  community_type?: CommunityType;
  my_communities?: boolean;
  page?: number;
  per_page?: number;
}

export interface CreateCommunityRequest {
  name: string;
  description: string;
  community_type: CommunityType;
  access_type: CommunityAccess;
  district?: string;
}

/* ── Helpers ── */
function unwrapData<T>(response: { data: { data: T } }): T {
  return response.data.data;
}

/* ── API ── */
export async function getCommunities(
  params: CommunitySearchParams
): Promise<PaginatedResponse<CommunityItem>> {
  const res = await api.get<{ data: PaginatedResponse<CommunityItem> }>('/communities', { params });
  return unwrapData(res);
}

export async function getCommunityById(id: string): Promise<CommunityItem> {
  const res = await api.get<{ data: CommunityItem }>(`/communities/${id}`);
  return unwrapData(res);
}

export async function getCommunityMembers(
  id: string,
  params?: { page?: number; per_page?: number; role?: MemberRole }
): Promise<PaginatedResponse<CommunityMember>> {
  const res = await api.get<{ data: PaginatedResponse<CommunityMember> }>(
    `/communities/${id}/members`,
    { params }
  );
  return unwrapData(res);
}

export async function createCommunity(data: CreateCommunityRequest): Promise<CommunityItem> {
  const res = await api.post<{ data: CommunityItem }>('/communities', data);
  return unwrapData(res);
}

export async function joinCommunity(id: string): Promise<void> {
  await api.post(`/communities/${id}/join`);
}

export async function leaveCommunity(id: string): Promise<void> {
  await api.delete(`/communities/${id}/join`);
}
