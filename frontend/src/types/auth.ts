export interface User {
  id: number;
  phone: string;
  email?: string;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  user_id: number;
  expires_at: string;
}

export interface LoginRequest {
  phone: string;
  password: string;
}

export interface RegisterRequest {
  phone: string;
  password: string;
  email?: string;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

export interface ApiError {
  code: string;
  message: string;
  details?: string;
}

export interface ApiResponse<T = any> {
  data?: T;
  error?: ApiError;
}
