'use client';

import Image from 'next/image';
import type { PhotoPost } from '../../lib/types';
import styles from './PhotoCard.module.css';

interface Props {
  post: PhotoPost;
  onOpen: (images: PhotoPost['images'], startIndex: number) => void;
}

export default function PhotoCard({ post, onOpen }: Props) {
  const { images } = post;
  const count = images.length;

  const date = new Date(post.createdAt).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  });

  const gridClass =
    count === 1 ? styles.single : count === 2 ? styles.double : styles.triple;
  const visibleImages = images.slice(0, 3);
  const extraCount = count > 3 ? count - 3 : 0;

  return (
    <article className={styles.card}>
      <div className={`${styles.imageGrid} ${gridClass}`}>
        {count <= 2
          ? visibleImages.map((img, i) => (
              <div
                key={img.url}
                className={styles.imageWrapper}
                onClick={() => onOpen(images, i)}
              >
                <Image
                  src={img.url}
                  alt={img.alt}
                  fill
                  sizes="(max-width: 680px) 50vw, 340px"
                />
              </div>
            ))
          : (
            <>
              <div
                className={styles.imageWrapperMain}
                onClick={() => onOpen(images, 0)}
              >
                <Image
                  src={visibleImages[0].url}
                  alt={visibleImages[0].alt}
                  fill
                  sizes="(max-width: 680px) 60vw, 408px"
                />
              </div>
              <div className={styles.imageStack}>
                {visibleImages.slice(1).map((img, i) => (
                  <div
                    key={img.url}
                    className={styles.imageWrapper}
                    onClick={() => onOpen(images, i + 1)}
                  >
                    <Image
                      src={img.url}
                      alt={img.alt}
                      fill
                      sizes="(max-width: 680px) 40vw, 272px"
                    />
                    {i === 1 && extraCount > 0 && (
                      <div
                        className={styles.countOverlay}
                        onClick={(e) => {
                          e.stopPropagation();
                          onOpen(images, i + 1);
                        }}
                      >
                        +{extraCount} more
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </>
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
