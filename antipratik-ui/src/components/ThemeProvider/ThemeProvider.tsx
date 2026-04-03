'use client';

import React, { createContext, useContext, useEffect, useState } from 'react';
import type { Theme } from '../../lib/types';

interface ThemeContextType {
  theme: Theme;
  toggle: () => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setTheme] = useState<Theme>('dark');
  const [mounted, setMounted] = useState(false);

  // On mount, read localStorage and set document attribute
  useEffect(() => {
    try {
      const stored = localStorage.getItem('ap-theme') as Theme | null;
      const initial = stored || 'dark';
      setTheme(initial);
      document.documentElement.dataset.theme = initial;
    } catch {
      // localStorage unavailable
      document.documentElement.dataset.theme = 'dark';
    }
    setMounted(true);
  }, []);

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

  // Don't render children until mounted to avoid hydration mismatch
  if (!mounted) {
    return null;
  }

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
