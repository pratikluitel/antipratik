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

  const date = new Date(post.createdAt).toLocaleDateString('en-US', {
    month: 'long',
    day: 'numeric',
    year: 'numeric',
  });

  const htmlBody = marked.parse(post.body) as string;

  return (
    <>
      {/* Reading progress — desktop only, fixed right */}
      <div className={styles.progressTrack} aria-hidden="true">
        <div
          className={styles.progressFill}
          style={{ height: `${progress * 100}%` }}
        />
      </div>

      {/* Back button — desktop only, fixed left */}
      <Link href="/feed" className={styles.backBtn}>
        ← Feed
      </Link>

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
