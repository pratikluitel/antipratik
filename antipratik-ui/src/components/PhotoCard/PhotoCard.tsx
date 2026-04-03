'use client';

import Image from 'next/image';
import type { PhotoPost } from '../../lib/types';
import styles from './PhotoCard.module.css';

interface Props {
  post: PhotoPost;
  onOpen: (images: PhotoPost['images'], startIndex: number) => void;
}

const MultiImageIcon = () => (
  <svg 
    width="20" 
    height="20" 
    viewBox="0 0 24 24" 
    fill="none" 
    stroke="currentColor" 
    strokeWidth="2" 
    strokeLinecap="round" 
    strokeLinejoin="round"
  >
    <rect width="14" height="14" x="8" y="8" rx="2" ry="2"/>
    <path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/>
  </svg>
);

export default function PhotoCard({ post, onOpen }: Props) {
  const { images } = post;
  const count = images.length;
  const mainImage = images[0];

  const date = new Date(post.createdAt).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  });

  return (
    <article className={styles.card}>
      <div 
        className={styles.imageWrapper}
        onClick={() => onOpen(images, 0)}
      >
        <Image
          src={mainImage.url}
          alt={mainImage.alt}
          fill
          sizes="(max-width: 680px) 100vw, 680px"
        />
        {count > 1 && (
          <div className={styles.galleryIndicator}>
            <MultiImageIcon />
          </div>
        )}
      </div>

      <div className={styles.body}>
        <span className={styles.tag}>Photo</span>
        <div className={styles.meta}>
          {post.location && <span>· {post.location}</span>}
          <time dateTime={post.createdAt}>{date}</time>
        </div>
      </div>
    </article>
  );
}
