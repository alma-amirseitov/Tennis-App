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
import { getMatch, submitMatchResult, confirmMatchResult } from './matches';
import type { SubmitResultRequest, ConfirmResultRequest } from './matches';
import { getChats, getMessages, sendMessage, markChatRead, getUnreadCount } from './chat';

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

/* ── Matches ── */
export function useMatch(id: string) {
  return useQuery({ queryKey: ['match', id], queryFn: () => getMatch(id), enabled: !!id });
}

export function useSubmitMatchResult() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ matchId, data }: { matchId: string; data: SubmitResultRequest }) =>
      submitMatchResult(matchId, data),
    onSuccess: (_, { matchId }) => qc.invalidateQueries({ queryKey: ['match', matchId] }),
  });
}

export function useConfirmMatchResult() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ matchId, data }: { matchId: string; data: ConfirmResultRequest }) =>
      confirmMatchResult(matchId, data),
    onSuccess: (_, { matchId }) => {
      qc.invalidateQueries({ queryKey: ['match', matchId] });
      qc.invalidateQueries({ queryKey: ['profile'] });
    },
  });
}

/* ── Chat ── */
export function useChats() {
  return useQuery({ queryKey: ['chats'], queryFn: getChats, staleTime: 30_000 });
}

export function useChatMessages(chatId: string, params?: { before?: string; limit?: number }) {
  return useQuery({
    queryKey: ['chat-messages', chatId, params],
    queryFn: () => getMessages(chatId, params),
    enabled: !!chatId,
    staleTime: 10_000,
  });
}

export function useChatMessagesInfinite(chatId: string, limit = 50) {
  return useInfiniteQuery({
    queryKey: ['chat-messages-infinite', chatId],
    queryFn: ({ pageParam }: { pageParam?: string }) =>
      getMessages(chatId, { before: pageParam, limit }),
    getNextPageParam: (last) => {
      if (!last.has_more || last.data.length === 0) return undefined;
      return last.data[last.data.length - 1]?.id;
    },
    initialPageParam: undefined as string | undefined,
    enabled: !!chatId,
    staleTime: 10_000,
  });
}

export function useSendMessage() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ chatId, content, replyToId }: { chatId: string; content: string; replyToId?: string | null }) =>
      sendMessage(chatId, content, replyToId),
    onSuccess: (msg, { chatId }) => {
      qc.invalidateQueries({ queryKey: ['chat-messages', chatId] });
      qc.invalidateQueries({ queryKey: ['chat-messages-infinite', chatId] });
      qc.invalidateQueries({ queryKey: ['chats'] });
    },
  });
}

export function useMarkChatRead() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ chatId, lastReadAt }: { chatId: string; lastReadAt: string }) =>
      markChatRead(chatId, lastReadAt),
    onSuccess: (_, { chatId }) => {
      qc.invalidateQueries({ queryKey: ['chats'] });
    },
  });
}

export function useUnreadCount() {
  return useQuery({ queryKey: ['chats-unread'], queryFn: getUnreadCount, staleTime: 30_000 });
}
