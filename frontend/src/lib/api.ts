export const API_BASE_URL = (import.meta.env.VITE_API_URL as string | undefined) ?? "http://localhost:8080";

export interface ApiError {
  error: string;
  fields?: Record<string, string>;
}
