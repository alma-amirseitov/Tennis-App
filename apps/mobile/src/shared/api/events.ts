import { api } from './client';
import type { PaginatedResponse } from './users';

/* ── Types ── */
export type EventType = 'find_partner' | 'organized_game' | 'tournament' | 'training';
export type EventStatus = 'draft' | 'open' | 'filling' | 'full' | 'in_progress' | 'completed' | 'cancelled';
export type CompositionType = 'singles' | 'doubles' | 'mixed' | 'team';

export interface EventItem {
  id: string;
  community_id: string | null;
  creator_id: string;
  title: string;
  description: string;
  event_type: EventType;
  composition_type: CompositionType;
  level_min: number;
  level_max: number;
  max_participants: number;
  current_participants: number;
  sets_count: number;
  start_time: string;
  end_time: string | null;
  location_name: string;
  status: EventStatus;
  price: number;
  creator: { id: string; first_name: string; last_name: string; avatar_url: string | null; ntrp_level: number } | null;
  community: { id: string; name: string; logo_url: string | null; is_verified: boolean } | null;
  participants: EventParticipant[];
  is_joined: boolean;
  created_at: string;
}

export interface EventParticipant {
  id: string;
  user_id: string;
  first_name: string;
  last_name: string;
  avatar_url: string | null;
  ntrp_level: number;
}

export interface EventSearchParams {
  event_type?: EventType;
  status?: EventStatus;
  composition_type?: CompositionType;
  level_min?: number;
  level_max?: number;
  date_from?: string;
  date_to?: string;
  district?: string;
  community_id?: string;
  my_events?: boolean;
  page?: number;
  per_page?: number;
}

export interface CreateEventRequest {
  community_id?: string;
  title: string;
  description?: string;
  event_type: EventType;
  composition_type: CompositionType;
  level_min: number;
  level_max: number;
  max_participants: number;
  sets_count: number;
  start_time: string;
  end_time?: string;
  location_name: string;
  price?: number;
}

/* ── Helpers ── */
function unwrapData<T>(response: { data: { data: T } }): T {
  return response.data.data;
}

/* ── API ── */
export async function getEvents(params: EventSearchParams): Promise<PaginatedResponse<EventItem>> {
  const res = await api.get<{ data: PaginatedResponse<EventItem> }>('/events', { params });
  return unwrapData(res);
}

export async function getEventById(id: string): Promise<EventItem> {
  const res = await api.get<{ data: EventItem }>(`/events/${id}`);
  return unwrapData(res);
}

export async function createEvent(data: CreateEventRequest): Promise<EventItem> {
  const res = await api.post<{ data: EventItem }>('/events', data);
  return unwrapData(res);
}

export async function joinEvent(id: string): Promise<void> {
  await api.post(`/events/${id}/join`);
}

export async function leaveEvent(id: string): Promise<void> {
  await api.delete(`/events/${id}/join`);
}
