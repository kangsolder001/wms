import api from './client';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: string;
  meta?: {
    page: number;
    limit: number;
    total: number;
  };
}

export interface LoginData {
  token: string;
  user: UserResponse;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  full_name: string;
  role: string;
}

export interface UserResponse {
  id: string;
  username: string;
  email: string;
  full_name: string;
  role: string;
  is_active: boolean;
}

export const authApi = {
  login: (data: LoginRequest) => api.post<ApiResponse<LoginData>>('/auth/login', data),
  register: (data: RegisterRequest) => api.post<ApiResponse<UserResponse>>('/auth/register', data),
  getProfile: () => api.get<ApiResponse<UserResponse>>('/auth/me'),
  listUsers: () => api.get<ApiResponse<UserResponse[]>>('/auth/users'),
  updateUser: (id: string, data: Partial<RegisterRequest & { is_active: boolean }>) =>
    api.put<ApiResponse<UserResponse>>(`/auth/users/${id}`, data),
  deleteUser: (id: string) => api.delete<ApiResponse<{ message: string }>>(`/auth/users/${id}`),
};
