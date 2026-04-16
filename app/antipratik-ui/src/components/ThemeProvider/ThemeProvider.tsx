'use client';

import React, { createContext, useContext, useState } from 'react';
import type { Theme } from '../../lib/types';

interface ThemeContextType {
  theme: Theme;
  toggle: () => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  // Read the data-theme attribute already set by the inline script in layout.tsx,
  // so we never need to call setTheme inside an effect.
  const [theme, setTheme] = useState<Theme>(() => {
    if (typeof document === 'undefined') return 'dark';
    return (document.documentElement.dataset.theme as Theme) || 'dark';
  });

  const toggle = () => {
    setTheme((prev) => {
      const next = prev === 'dark' ? 'light' : 'dark';
      document.documentElement.dataset.theme = next;
      try {
        localStorage.setItem('ap-theme', next);
      } catch {
        // localStorage unavailable
      }
      return next;
    });
  };

  // We always render children to keep the DOM tree stable for Next.js hydration.
  // The RootLayout inline script handles the theme attribute immediately to prevent flash.
  return (
    <ThemeContext.Provider value={{ theme, toggle }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme(): ThemeContextType {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be called inside <ThemeProvider>');
  }
  return context;
}
