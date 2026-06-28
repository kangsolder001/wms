import { create } from 'zustand';

interface ThemeState {
  isDark: boolean;
  toggle: () => void;
}

function applyTheme(dark: boolean) {
  document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light');
}

const initialDark = localStorage.getItem('theme') === 'dark';
applyTheme(initialDark);

export const useThemeStore = create<ThemeState>((set) => ({
  isDark: initialDark,
  toggle: () =>
    set((state) => {
      const next = !state.isDark;
      localStorage.setItem('theme', next ? 'dark' : 'light');
      applyTheme(next);
      return { isDark: next };
    }),
}));
