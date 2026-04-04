'use client';

import { createContext, useContext, useState } from 'react';

interface NavbarContextValue {
  articleTitle: string | null;
  setArticleTitle: (t: string | null) => void;
}

const NavbarContext = createContext<NavbarContextValue | null>(null);

export function NavbarProvider({ children }: { children: React.ReactNode }) {
  const [articleTitle, setArticleTitle] = useState<string | null>(null);
  return (
    <NavbarContext.Provider value={{ articleTitle, setArticleTitle }}>
      {children}
    </NavbarContext.Provider>
  );
}

export function useNavbarContext(): NavbarContextValue {
  const ctx = useContext(NavbarContext);
  if (!ctx) throw new Error('useNavbarContext must be used inside NavbarProvider');
  return ctx;
}
