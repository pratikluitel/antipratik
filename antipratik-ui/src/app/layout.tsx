import type { Metadata } from 'next';
import { ThemeProvider } from '../components/ThemeProvider';
import Navbar from '../components/Navbar';
import '../styles/tokens.css';
import '../styles/globals.css';

export const metadata: Metadata = {
  title: 'antipratik',
  description: 'Developer, musician, writer, and photographer based in Kathmandu, Nepal.',
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
          <Navbar />
          {children}
        </ThemeProvider>
      </body>
    </html>
  );
}
