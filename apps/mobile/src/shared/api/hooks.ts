import { useQuery, useMutation, useInfiniteQuery, useQueryClient } from '@tanstack/react-query';
import { getMyProfile, updateMyProfile, searchPlayers, getUserById } from './users';
import type { UpdateProfileRequest, PlayerSearchParams } from './users';
import { getEvents, getEventById, createEvent, joinEvent, leaveEvent } from './events';
import type { EventSearchParams, CreateEventRequest } from './events';
import {
  getCommunities, getCommunityById, getCommunityMembers,
  createCommunity, joinCommunity, leaveCommunity,
} from './communities';
import type { CommunitySearchParams, CreateCommunityRequest, MemberRole } from './communities';

/* ── Profile ── */
export function useProfile() {
  return useQuery({ queryKey: ['profile'], queryFn: getMyProfile, staleTime: 60_000 });
}

export function useUpdateProfile() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: UpdateProfileRequest) => updateMyProfile(data),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['profile'] }); },
  });
}

/* ── Players ── */
export function useSearchPlayers(params: PlayerSearchParams) {
  return useInfiniteQuery({
    queryKey: ['players', params],
    queryFn: ({ pageParam = 1 }) => searchPlayers({ ...params, page: pageParam, per_page: 20 }),
    getNextPageParam: (last) => {
      const { page, total_pages } = last.pagination;
      return page < total_pages ? page + 1 : undefined;
    },
    initialPageParam: 1,
    staleTime: 30_000,
  });
}

export function usePlayer(id: string) {
  return useQuery({ queryKey: ['player', id], queryFn: () => getUserById(id), enabled: !!id });
}

/* ── Events ── */
export function useEvents(params: EventSearchParams) {
  return useInfiniteQuery({
    queryKey: ['events', params],
    queryFn: ({ pageParam = 1 }) => getEvents({ ...params, page: pageParam, per_page: 20 }),
    getNextPageParam: (last) => {
      const { page, total_pages } = last.pagination;
      return page < total_pages ? page + 1 : undefined;
    },
    initialPageParam: 1,
    staleTime: 30_000,
  });
}

export function useEvent(id: string) {
  return useQuery({ queryKey: ['event', id], queryFn: () => getEventById(id), enabled: !!id });
}

export function useCreateEvent() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateEventRequest) => createEvent(data),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['events'] }); },
  });
}

export function useJoinEvent() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => joinEvent(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['events'] });
      qc.invalidateQueries({ queryKey: ['event'] });
    },
  });
}

export function useLeaveEvent() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => leaveEvent(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['events'] });
      qc.invalidateQueries({ queryKey: ['event'] });
    },
  });
}

/* ── Communities ── */
export function useCommunities(params: CommunitySearchParams) {
  return useInfiniteQuery({
    queryKey: ['communities', params],
    queryFn: ({ pageParam = 1 }) => getCommunities({ ...params, page: pageParam, per_page: 20 }),
    getNextPageParam: (last) => {
      const { page, total_pages } = last.pagination;
      return page < total_pages ? page + 1 : undefined;
    },
    initialPageParam: 1,
    staleTime: 30_000,
  });
}

export function useCommunity(id: string) {
  return useQuery({ queryKey: ['community', id], queryFn: () => getCommunityById(id), enabled: !!id });
}

export function useCommunityMembers(id: string, params?: { page?: number; per_page?: number; role?: MemberRole }) {
  return useQuery({
    queryKey: ['community-members', id, params],
    queryFn: () => getCommunityMembers(id, params),
    enabled: !!id,
  });
}

export function useCreateCommunity() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateCommunityRequest) => createCommunity(data),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['communities'] }); },
  });
}

export function useJoinCommunity() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => joinCommunity(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['communities'] });
      qc.invalidateQueries({ queryKey: ['community'] });
    },
  });
}

export function useLeaveCommunity() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => leaveCommunity(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['communities'] });
      qc.invalidateQueries({ queryKey: ['community'] });
    },
  });
}
