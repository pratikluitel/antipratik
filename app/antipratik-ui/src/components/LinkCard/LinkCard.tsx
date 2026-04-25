'use client';

import Image from 'next/image';
import { useState } from 'react';
import type { LinkPost } from '../../lib/types';
import styles from './LinkCard.module.css';

interface Props {
  post: LinkPost;
}

function cardAccentClass(category: LinkPost['category']): string {
  switch (category) {
    case 'music':   return styles.cardMusic;
    case 'writing': return styles.cardEssay;
    case 'video':   return styles.cardVideo;
    default:        return styles.cardDefault;
  }
}

function tagClass(category: LinkPost['category']): string | undefined {
  switch (category) {
    case 'music':   return `${styles.tag} ${styles.tagMusic}`;
    case 'writing': return `${styles.tag} ${styles.tagEssay}`;
    case 'video':   return `${styles.tag} ${styles.tagVideo}`;
    default:        return undefined; // social or undefined → no tag
  }
}

function tagLabel(category: LinkPost['category']): string | null {
  switch (category) {
    case 'music':   return 'Music';
    case 'writing': return 'Essay';
    case 'video':   return 'Video';
    default:        return null;
  }
}

function LinkThumbnail({ post }: { post: LinkPost }) {
  const [loaded, setLoaded] = useState(false);
  return (
    <div className={styles.thumbnailWrapper}>
      {post.thumbnailTinyUrl && !loaded && (
        <div
          className={styles.thumbnailPlaceholder}
          style={{ backgroundImage: `url(${post.thumbnailTinyUrl})` }}
          aria-hidden="true"
        />
      )}
      <Image
        src={post.thumbnailSmallUrl ?? post.thumbnailUrl!}
        alt={post.title}
        width={52}
        height={52}
        onLoad={() => setLoaded(true)}
        className={`${styles.thumbnail} ${loaded ? styles.imageVisible : styles.imageHidden}`}
      />
    </div>
  );
}

export default function LinkCard({ post }: Props) {
  const tagCls = tagClass(post.category);
  const label = tagLabel(post.category);
  const date = new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(new Date(post.createdAt));

  const inner = (
    <>
      {post.thumbnailUrl && (
        <LinkThumbnail post={post} />
      )}

      <div className={styles.textBlock}>
        {tagCls && label && <span className={tagCls}>{label}</span>}
        <span className={styles.title}>{post.title}</span>
        <span className={styles.domain}>{post.domain}</span>
        {post.description && (
          <p className={styles.excerpt}>{post.description}</p>
        )}
        <time className={styles.date} dateTime={post.createdAt}>{date}</time>
      </div>

      <span className={styles.arrow} aria-hidden="true">{post.category === 'video' ? '▶' : '↗'}</span>
    </>
  );

  return (
    <article className={cardAccentClass(post.category)}>
      <a
        href={post.url}
        target="_blank"
        rel="noopener noreferrer"
        className={styles.link}
      >
        {inner}
      </a>
    </article>
  );
}
