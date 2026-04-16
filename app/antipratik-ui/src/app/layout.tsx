import type { Metadata } from 'next';
import { ThemeProvider } from '../components/ThemeProvider';
import { NavbarProvider } from '../components/NavbarContext';
import { MusicProvider } from '../components/MusicProvider';
import { MusicPlayerRoot } from '../components/MusicPlayer';
import Navbar from '../components/Navbar';
import '../styles/tokens.css';
import '../styles/globals.css';

export const metadata: Metadata = {
  title: 'antipratik',
  description: 'Developer, musician, and blogger based in Kathmandu, Nepal.',
  openGraph: {
    title: 'antipratik',
    description: 'Developer, music tinkerer, and blogger based in Kathmandu, Nepal.',
    url: 'https://antipratik.com',
    siteName: 'antipratik',
    locale: 'en_US',
    type: 'website',
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        {/* Google Fonts — DM Serif Display + DM Sans */}
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        {/* eslint-disable-next-line @next/next/no-page-custom-font -- App Router uses layout.tsx, not pages/_document.js */}
        <link
          href="https://fonts.googleapis.com/css2?family=DM+Sans:ital,wght@0,300;0,400;0,500;1,400&family=DM+Serif+Display:ital@0;1&display=swap"
          rel="stylesheet"
        />
        {/* Inline script: set data-theme before first paint to prevent flash */}
        <script
          dangerouslySetInnerHTML={{
            __html: `(function(){try{var t=localStorage.getItem('ap-theme')||'dark';document.documentElement.setAttribute('data-theme',t);}catch(e){document.documentElement.setAttribute('data-theme','dark');}})();`,
          }}
        />
      </head>
      <body>
        <ThemeProvider>
          <NavbarProvider>
            <MusicProvider>
              <Navbar />
              {children}
              <MusicPlayerRoot />
            </MusicProvider>
          </NavbarProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
