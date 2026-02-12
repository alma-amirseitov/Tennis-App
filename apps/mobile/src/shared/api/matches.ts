import { api } from './client';

export interface MatchSetScore {
  p1: number;
  p2: number;
  tiebreak?: { p1: number; p2: number };
}

export interface MatchScore {
  sets: MatchSetScore[];
}

export interface MatchDetail {
  id: string;
  event_id: string;
  player1_id: string;
  player2_id: string;
  player1: { id: string; first_name: string; last_name: string; ntrp_level: number; avatar_url: string | null };
  player2: { id: string; first_name: string; last_name: string; ntrp_level: number; avatar_url: string | null };
  score: MatchScore | null;
  winner_id: string | null;
  result_status: 'pending' | 'confirmed' | 'disputed' | 'admin_confirmed';
  played_at: string;
}

export interface SubmitResultRequest {
  winner_id: string;
  score: MatchScore;
}

export interface SubmitResultResponse {
  match_id: string;
  result_status: 'pending';
  message: string;
}

export interface ConfirmResultRequest {
  action: 'confirm' | 'dispute';
  reason?: string;
}

export interface RatingChange {
  before: number;
  after: number;
  change: number;
}

export interface ConfirmResultResponse {
  result_status: 'confirmed' | 'disputed';
  rating_changes?: {
    player1: RatingChange;
    player2: RatingChange;
  };
}

function unwrapData<T>(response: { data: { data: T } }): T {
  return response.data.data;
}

export async function getMatch(id: string): Promise<MatchDetail> {
  const res = await api.get<{ data: MatchDetail }>(`/matches/${id}`);
  return unwrapData(res);
}

export async function submitMatchResult(
  matchId: string,
  data: SubmitResultRequest
): Promise<SubmitResultResponse> {
  const res = await api.post<{ data: SubmitResultResponse }>(
    `/matches/${matchId}/result`,
    data
  );
  return unwrapData(res);
}

export async function confirmMatchResult(
  matchId: string,
  data: ConfirmResultRequest
): Promise<ConfirmResultResponse> {
  const res = await api.post<{ data: ConfirmResultResponse }>(
    `/matches/${matchId}/confirm`,
    data
  );
  return unwrapData(res);
}
