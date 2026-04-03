'use client';

import { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useTheme } from '../ThemeProvider';
import { useNavbarContext } from '../NavbarContext';
import styles from './Navbar.module.css';

function SunSVG() {
  return (
    <svg width="18" height="18" viewBox="0 0 18 18" fill="none" aria-hidden="true" className={styles.themeIcon}>
      <circle cx="9" cy="9" r="3" stroke="currentColor" strokeWidth="1.5" />
      <line x1="9"    y1="1.5"  x2="9"    y2="3.5"  stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <line x1="9"    y1="14.5" x2="9"    y2="16.5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <line x1="1.5"  y1="9"   x2="3.5"  y2="9"    stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <line x1="14.5" y1="9"   x2="16.5" y2="9"    stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <line x1="3.69" y1="3.69" x2="5.11" y2="5.11" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <line x1="12.89" y1="12.89" x2="14.31" y2="14.31" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <line x1="14.31" y1="3.69" x2="12.89" y2="5.11" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <line x1="5.11" y1="12.89" x2="3.69" y2="14.31" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
    </svg>
  );
}

function MoonSVG() {
  return (
    <svg width="18" height="18" viewBox="0 0 18 18" fill="none" aria-hidden="true" className={styles.themeIcon}>
      <path
        d="M15.75 9.6 A6.75 6.75 0 1 1 8.4 2.25 A5.25 5.25 0 0 0 15.75 9.6 Z"
        fill="currentColor"
      />
    </svg>
  );
}

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
        <Link href="/" className={styles.logo}>
          antipratik
        </Link>

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
            <div className={styles.toggleTrack}>
              <div className={styles.toggleThumb}>
                {theme === 'dark' ? <SunSVG /> : <MoonSVG />}
              </div>
            </div>
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
              <div className={styles.toggleTrack}>
                <div className={styles.toggleThumb}>
                  {theme === 'dark' ? <SunSVG /> : <MoonSVG />}
                </div>
              </div>
            </button>
          </div>
        </div>
      )}
    </nav>
  );
}
