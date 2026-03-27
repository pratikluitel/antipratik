'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useTheme } from '../ThemeProvider';
import styles from './Navbar.module.css';

interface NavbarProps {
  articleTitle?: string;
}

export default function Navbar({ articleTitle }: NavbarProps) {
  const [scrolled, setScrolled] = useState(false);
  const pathname = usePathname();
  const { theme, toggle } = useTheme();

  useEffect(() => {
    if (typeof window === 'undefined') return;

    const handleScroll = () => {
      setScrolled(window.scrollY > 20);
    };

    window.addEventListener('scroll', handleScroll);
    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, []);

  const showArticleTitle = articleTitle && scrolled && window.scrollY > 80;

  return (
    <nav className={styles.navbar} data-scrolled={scrolled ? 'true' : 'false'}>
      <div className={styles.inner}>
        <a href="/" className={styles.logo}>
          antipratik
        </a>

        {articleTitle && (
          <span
            className={styles.articleTitle}
            data-show={showArticleTitle ? 'true' : 'false'}
          >
            {articleTitle}
          </span>
        )}

        <div className={styles.navLinks}>
          <Link
            href="/feed"
            className={pathname === '/feed' ? styles.navLink + ' ' + styles.active : styles.navLink}
          >
            Feed
          </Link>
          <Link
            href="/links"
            className={pathname === '/links' ? styles.navLink + ' ' + styles.active : styles.navLink}
          >
            Links
          </Link>
          <Link
            href="/about"
            className={pathname === '/about' ? styles.navLink + ' ' + styles.active : styles.navLink}
          >
            About
          </Link>
        </div>

        <div className={styles.controls}>
          <button
            className={styles.themeToggle}
            onClick={toggle}
            aria-label={`Switch to ${theme === 'dark' ? 'light' : 'dark'} mode`}
            aria-pressed={theme === 'light'}
          >
            <span className={styles.toggleTrack}>
              <span className={styles.toggleThumb} />
            </span>
          </button>
        </div>
      </div>
    </nav>
  );
}
