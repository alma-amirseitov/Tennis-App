import { api } from './client';

export interface QuizOption {
  id: string;
  text_ru: string;
  text_kz?: string;
  text_en?: string;
  weight: number;
}

export interface QuizQuestion {
  id: number;
  text_ru: string;
  text_kz?: string;
  text_en?: string;
  options: QuizOption[];
}

export interface QuizQuestionsResponse {
  questions: QuizQuestion[];
}

export interface QuizAnswer {
  question_id: number;
  option_id: string;
}

export interface QuizSubmitResponse {
  ntrp_level: number;
  level_label: string;
  initial_rating: number;
}

function unwrapData<T>(response: { data: { data: T } }): T {
  return response.data.data;
}

export async function getQuizQuestions(): Promise<QuizQuestionsResponse> {
  const response = await api.get<{ data: QuizQuestionsResponse }>(
    '/quiz/questions'
  );
  return unwrapData(response);
}

export async function submitQuiz(
  answers: QuizAnswer[]
): Promise<QuizSubmitResponse> {
  const response = await api.post<{ data: QuizSubmitResponse }>(
    '/quiz/submit',
    { answers }
  );
  return unwrapData(response);
}
