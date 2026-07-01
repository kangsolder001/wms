import { create } from 'zustand';
import type { UserResponse } from '../api/auth';

interface AuthState {
  token: string | null;
  user: UserResponse | null;
  isAuthenticated: boolean;
  setAuth: (token: string, user: UserResponse) => void;
  logout: () => void;
}

function loadUser(): UserResponse | null {
  const raw = localStorage.getItem('user');
  if (!raw) return null;
  try { return JSON.parse(raw); } catch { return null; }
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('token'),
  user: loadUser(),
  isAuthenticated: !!localStorage.getItem('token'),
  setAuth: (token, user) => {
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    set({ token, user, isAuthenticated: true });
  },
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    set({ token: null, user: null, isAuthenticated: false });
  },
}));
