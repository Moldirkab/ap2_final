import api from "./client";

export interface AuthResponse {
  token: string;
}

export interface ProfileResponse {
  user_id: string;
  role: string;
}

export const register = (email: string, password: string) =>
  api.post<AuthResponse>("/auth/register", { email, password });

export const login = (email: string, password: string) =>
  api.post<AuthResponse>("/auth/login", { email, password });

export const getProfile = () => api.get<ProfileResponse>("/auth/profile");
