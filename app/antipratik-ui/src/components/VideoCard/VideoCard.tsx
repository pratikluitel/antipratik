'use client';

import Image from 'next/image';
import { useState } from 'react';
import type { VideoPost } from '../../lib/types';
import styles from './VideoCard.module.css';

interface Props {
  post: VideoPost;
}

function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${m}:${s.toString().padStart(2, '0')}`;
}

export default function VideoCard({ post }: Props) {
  const [loaded, setLoaded] = useState(false);
  const date = new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(new Date(post.createdAt));

  return (
    <article className={styles.card}>
      <a
        href={post.videoUrl}
        target="_blank"
        rel="noopener noreferrer"
        className={styles.link}
      >
        <div className={styles.thumbnail}>
          {post.thumbnailUrl && post.thumbnailTinyUrl && (
            <div
              className={styles.imagePlaceholder}
              style={{ backgroundImage: `url(${post.thumbnailTinyUrl})` }}
              aria-hidden="true"
            />
          )}
          {post.thumbnailUrl && (
            <Image
              src={post.thumbnailLargeUrl ?? post.thumbnailUrl}
              alt={post.title}
              fill
              sizes="(max-width: 860px) 100vw, 860px"
              onLoad={() => setLoaded(true)}
              className={loaded ? styles.imageVisible : styles.imageHidden}
            />
          )}
          <div className={styles.scrim} />
          <div className={styles.playButton} aria-hidden="true">
            <span className={styles.playIcon} />
          </div>
          <time className={styles.dateOverlay} dateTime={post.createdAt}>
            {date}
          </time>
        </div>

        <div className={styles.body}>
          <span className={styles.tag}>Video</span>
          <h2 className={styles.title}>{post.title}</h2>
          <div className={styles.footerRow}>
            <span>{post.playlist ?? ''}</span>
            <span>{formatDuration(post.duration)}</span>
          </div>
        </div>
      </a>
    </article>
  );
}
