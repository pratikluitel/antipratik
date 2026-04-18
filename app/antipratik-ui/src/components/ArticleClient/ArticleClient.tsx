'use client';

import { useEffect, useMemo, useState } from 'react';
import Link from 'next/link';
import { marked } from 'marked';
import type { EssayPost } from '../../lib/types';
import { useNavbarContext } from '../NavbarContext';
import styles from './ArticleClient.module.css';

interface Props {
  post: EssayPost;
}

interface Heading {
  id: string;
  text: string;
  level: number;
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

function slugify(text: string): string {
  return text.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '');
}

export default function ArticleClient({ post }: Props) {
  const { setArticleTitle } = useNavbarContext();
  const [progress, setProgress] = useState(0);
  const [tocOpen, setTocOpen] = useState(true);

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

  // Parse markdown, collecting headings and adding IDs
  const { htmlBody, headings } = useMemo(() => {
    const collected: Heading[] = [];
    const renderer = new marked.Renderer();
    renderer.heading = ({ text, depth }: { text: string; depth: number }) => {
      if (depth === 2 || depth === 3) {
        const id = slugify(text);
        collected.push({ id, text, level: depth });
        return `<h${depth} id="${id}">${text}</h${depth}>`;
      }
      return `<h${depth}>${text}</h${depth}>`;
    };
    const html = marked.parse(post.body, { breaks: true, renderer }) as string;
    return { htmlBody: html, headings: collected };
  }, [post.body]);

  return (
    <>
      {/* Back button — desktop only, fixed left */}
      <Link href="/feed" className={styles.backBtn}>
        ← Feed
      </Link>

      {/* Ruler-style progress + scroll-to-top — desktop only, fixed right */}
      <div className={styles.sideControls}>
        <button
          className={`${styles.scrollToTop}${progress > 0 ? ` ${styles.scrollToTopVisible}` : ''}`}
          onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
          aria-label="Scroll to top"
        >
          <ChevronUpSVG />
        </button>
        <div className={styles.progressRuler} aria-hidden="true">
          {Array.from({ length: 11 }, (_, i) => {
            const isRead = i / 10 <= progress;
            const isLong = i === 5;
            const cls = [
              styles.tick,
              isLong ? styles.tickLong : '',
              isRead ? styles.tickRead : '',
            ].filter(Boolean).join(' ');
            return (
              <span key={i} className={cls} style={{ top: `${i * 10}%` }} />
            );
          })}
        </div>
      </div>

      {/* Table of Contents — desktop only, fixed right panel */}
      {headings.length > 0 && (
        <nav
          className={`${styles.toc}${tocOpen ? ` ${styles.tocOpen}` : ''}`}
          aria-label="Table of contents"
        >
          <button
            className={styles.tocToggle}
            onClick={() => setTocOpen((o) => !o)}
            aria-label={tocOpen ? 'Collapse contents' : 'Expand contents'}
          >
            {tocOpen ? 'Contents' : '≡'}
          </button>
          {tocOpen && (
            <ol className={styles.tocList}>
              {headings.map((h) => (
                <li
                  key={h.id}
                  className={h.level === 3 ? styles.tocItemSub : styles.tocItem}
                >
                  <a href={`#${h.id}`} className={styles.tocLink}>
                    {h.text}
                  </a>
                </li>
              ))}
            </ol>
          )}
        </nav>
      )}

      <article className={styles.article}>
        <header className={styles.header}>
          <span className={styles.tag}>{CATEGORY_LABEL}</span>
          <h1 className={styles.title}>{post.title}</h1>
          <div className={styles.meta}>
            <span>{date}</span>
            <span className={styles.dot}>·</span>
            <span>{post.readingTimeMinutes} min read</span>
          </div>
          {post.tags.length > 0 && (
            <div className={styles.tags}>
              {post.tags.map((tag) => (
                <Link
                  key={tag}
                  href={`/feed?tag=${encodeURIComponent(tag)}`}
                  className={styles.tagLink}
                >
                  #{tag}
                </Link>
              ))}
            </div>
          )}
        </header>

        <div
          className={styles.body}
          dangerouslySetInnerHTML={{ __html: htmlBody }}
        />
      </article>
    </>
  );
}
