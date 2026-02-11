import { api } from './client';

export interface SendOTPResponse {
  session_id: string;
  expires_in: number;
  retry_after?: number;
}

export interface VerifyOTPResponseExisting {
  is_new: false;
  access_token: string;
  refresh_token: string;
  user: { id: string; first_name: string; is_profile_complete: boolean };
}

export interface VerifyOTPResponseNew {
  is_new: true;
  temp_token: string;
  user_id: string;
}

export type VerifyOTPResponse = VerifyOTPResponseExisting | VerifyOTPResponseNew;

export interface RefreshResponse {
  access_token: string;
  refresh_token: string;
}

export interface ApiErrorResponse {
  error: {
    code: string;
    message?: string;
    details?: Array<{ field: string; message: string }>;
  };
}

function unwrapData<T>(response: { data: { data: T } }): T {
  return response.data.data;
}

export async function sendOTP(phone: string): Promise<SendOTPResponse> {
  const response = await api.post<{ data: SendOTPResponse }>('/auth/otp/send', {
    phone,
  });
  return unwrapData(response);
}

export async function verifyOTP(
  sessionId: string,
  code: string
): Promise<VerifyOTPResponse> {
  const response = await api.post<{ data: VerifyOTPResponse }>(
    '/auth/otp/verify',
    { session_id: sessionId, code }
  );
  return unwrapData(response);
}

export async function refreshToken(
  refreshTokenValue: string
): Promise<RefreshResponse> {
  const response = await api.post<{ data: RefreshResponse }>('/auth/refresh', {
    refresh_token: refreshTokenValue,
  });
  return unwrapData(response);
}

export interface ProfileSetupRequest {
  first_name: string;
  last_name: string;
  gender: 'male' | 'female';
  birth_year: number;
  city: string;
  district: string;
  language: string;
}

export interface ProfileSetupResponse {
  access_token: string;
  refresh_token: string;
  user: unknown;
}

export async function profileSetup(
  data: ProfileSetupRequest
): Promise<ProfileSetupResponse> {
  const response = await api.post<{ data: ProfileSetupResponse }>(
    '/auth/profile/setup',
    data
  );
  return unwrapData(response);
}
