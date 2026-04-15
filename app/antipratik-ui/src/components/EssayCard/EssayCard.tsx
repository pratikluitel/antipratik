import Link from 'next/link';
import type { EssayPost } from '../../lib/types';
import styles from './EssayCard.module.css';

interface Props {
  post: EssayPost;
}

export default function EssayCard({ post }: Props) {
  const date = new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(new Date(post.createdAt));

  return (
    <article className={styles.card}>
      <Link href={`/${post.slug}`} className={styles.inner}>
        <span className={styles.tag}>Essay</span>
        <h2 className={styles.title}>{post.title}</h2>
        <p className={styles.excerpt}>{post.excerpt}</p>
        <div className={styles.meta}>
          <span>{date}</span>
          <span className={styles.separator}>·</span>
          <span>{post.readingTimeMinutes} min read</span>
        </div>
      </Link>
    </article>
  );
}
