'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { marked } from 'marked';
import type { EssayPost } from '../../lib/types';
import { useNavbarContext } from '../NavbarContext';
import styles from './ArticleClient.module.css';

interface Props {
  post: EssayPost;
}

function ChevronUpSVG() {
  return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
      <path
        d="M3 10.5L8 5.5L13 10.5"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}

// The body markdown is author-controlled content from a trusted backend.
// dangerouslySetInnerHTML is acceptable here.
const CATEGORY_LABEL = 'Essay';

export default function ArticleClient({ post }: Props) {
  const { setArticleTitle } = useNavbarContext();
  const [progress, setProgress] = useState(0);

  // Register article title in Navbar context
  useEffect(() => {
    setArticleTitle(post.title);
    return () => setArticleTitle(null);
  }, [post.title, setArticleTitle]);

  // Reading progress tracker
  useEffect(() => {
    function handleScroll() {
      const el = document.documentElement;
      const scrolled = el.scrollTop;
      const total = el.scrollHeight - el.clientHeight;
      if (total <= 0) return;
      setProgress(Math.min(1, Math.max(0, scrolled / total)));
    }
    window.addEventListener('scroll', handleScroll, { passive: true });
    handleScroll();
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const date = new Intl.DateTimeFormat('en-US', {
    month: 'long',
    day: 'numeric',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(new Date(post.createdAt));

  const htmlBody = marked.parse(post.body, { breaks: true }) as string;

  return (
    <>
      {/* Back button — desktop only, fixed left */}
      <Link href="/feed" className={styles.backBtn}>
        ← Feed
      </Link>

      {/* Reading progress + scroll-to-top — desktop only, fixed right */}
      <div className={styles.sideControls}>
        <button
          className={`${styles.scrollToTop}${progress > 0 ? ` ${styles.scrollToTopVisible}` : ''}`}
          onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
          aria-label="Scroll to top"
        >
          <ChevronUpSVG />
        </button>
        <div className={styles.progressTrack} aria-hidden="true">
          <div
            className={styles.progressFill}
            style={{ height: `${progress * 100}%` }}
          />
        </div>
      </div>

      <article className={styles.article}>
        <header className={styles.header}>
          <span className={styles.tag}>{CATEGORY_LABEL}</span>
          <h1 className={styles.title}>{post.title}</h1>
          <div className={styles.meta}>
            <span>{date}</span>
            <span className={styles.dot}>·</span>
            <span>{post.readingTimeMinutes} min read</span>
          </div>
        </header>

        <div
          className={styles.body}
          dangerouslySetInnerHTML={{ __html: htmlBody }}
        />
      </article>
    </>
  );
}
