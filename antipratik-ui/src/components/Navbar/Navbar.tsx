'use client';

import { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useTheme } from '../ThemeProvider';
import { useNavbarContext } from '../NavbarContext';
import styles from './Navbar.module.css';

export default function Navbar() {
  const [scrolled, setScrolled] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);
  const navRef = useRef<HTMLElement>(null);
  const pathname = usePathname();
  const { theme, toggle } = useTheme();
  const { articleTitle } = useNavbarContext();

  useEffect(() => {
    if (typeof window === 'undefined') return;

    const handleScroll = () => {
      const isScrolled = window.scrollY > 20;
      setScrolled(isScrolled);
      document.documentElement.setAttribute('data-scrolled', isScrolled ? 'true' : 'false');
    };

    handleScroll(); // set initial state
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  // Close menu when clicking outside
  useEffect(() => {
    if (!menuOpen) return;
    function handleOutsideClick(e: MouseEvent) {
      if (navRef.current && !navRef.current.contains(e.target as Node)) {
        setMenuOpen(false);
      }
    }
    document.addEventListener('mousedown', handleOutsideClick);
    return () => document.removeEventListener('mousedown', handleOutsideClick);
  }, [menuOpen]);

  // Close menu on route change
  useEffect(() => {
    setMenuOpen(false);
  }, [pathname]);

  const showArticleTitle =
    articleTitle && scrolled && typeof window !== 'undefined' && window.scrollY > 80;

  function navLinkClass(href: string) {
    return pathname === href
      ? `${styles.navLink} ${styles.active}`
      : styles.navLink;
  }

  return (
    <nav
      ref={navRef}
      className={styles.navbar}
      data-scrolled={scrolled ? 'true' : 'false'}
    >
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
          <Link href="/feed" className={navLinkClass('/feed')}>Feed</Link>
          <Link href="/links" className={navLinkClass('/links')}>Links</Link>
          <Link href="/about" className={navLinkClass('/about')}>About</Link>
        </div>

        <div className={styles.controls}>
          {/* Desktop theme toggle */}
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

          {/* Hamburger — mobile only */}
          <button
            className={`${styles.hamburger}${menuOpen ? ` ${styles.hamburgerOpen}` : ''}`}
            onClick={() => setMenuOpen((o) => !o)}
            aria-label={menuOpen ? 'Close menu' : 'Open menu'}
            aria-expanded={menuOpen}
          >
            <span className={styles.bar1} />
            <span className={styles.bar2} />
            <span className={styles.bar3} />
          </button>
        </div>
      </div>

      {/* Mobile dropdown */}
      {menuOpen && (
        <div className={styles.mobileMenu}>
          <Link href="/feed" className={navLinkClass('/feed')} onClick={() => setMenuOpen(false)}>
            Feed
          </Link>
          <Link href="/links" className={navLinkClass('/links')} onClick={() => setMenuOpen(false)}>
            Links
          </Link>
          <Link href="/about" className={navLinkClass('/about')} onClick={() => setMenuOpen(false)}>
            About
          </Link>
          <div className={styles.mobileMenuDivider} />
          <div className={styles.mobileMenuToggleRow}>
            <span className={styles.mobileMenuToggleLabel}>
              {theme === 'dark' ? 'Dark mode' : 'Light mode'}
            </span>
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
      )}
    </nav>
  );
}
