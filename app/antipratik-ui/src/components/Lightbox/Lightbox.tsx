'use client';

import { useState, useEffect } from 'react';
import { createPortal } from 'react-dom';
import type { PhotoPost } from '../../lib/types';
import styles from './Lightbox.module.css';

interface Props {
  images: PhotoPost['images'];
  startIndex: number;
  onClose: () => void;
}

export default function Lightbox({ images, startIndex, onClose }: Props) {
  const [currentIndex, setCurrentIndex] = useState(startIndex);
  const [loaded, setLoaded] = useState(false);

  const next = () => setCurrentIndex(i => Math.min(i + 1, images.length - 1));
  const prev = () => setCurrentIndex(i => Math.max(i - 1, 0));

  useEffect(() => {
    setLoaded(false);
  }, [currentIndex]);

  useEffect(() => {
    document.body.style.overflow = 'hidden';
    return () => { document.body.style.overflow = ''; };
  }, []);

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'ArrowRight') next();
      if (e.key === 'ArrowLeft') prev();
      if (e.key === 'Escape') onClose();
    };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentIndex]);

  return createPortal(
    <div className={styles.backdrop} onClick={onClose}>

      <button
        className={styles.closeBtn}
        onClick={onClose}
        aria-label="Close lightbox"
      >
        ✕
      </button>

      {images.length > 1 && (
        <div className={styles.counter}>
          {currentIndex + 1} / {images.length}
        </div>
      )}

      <div className={styles.imageContainer} onClick={e => e.stopPropagation()}>
        {images[currentIndex].thumbnailTinyUrl && !loaded && (
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={images[currentIndex].thumbnailTinyUrl}
            alt=""
            aria-hidden="true"
            className={`${styles.image} ${styles.imagePlaceholder}`}
          />
        )}
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src={images[currentIndex].url}
          alt={images[currentIndex].alt}
          className={`${styles.image} ${loaded ? styles.imageVisible : styles.imageHidden}`}
          onLoad={() => setLoaded(true)}
        />
        {images[currentIndex].caption && (
          <div className={styles.caption}>
            {images[currentIndex].caption}
          </div>
        )}
      </div>

      {images.length > 1 && (
        <>
          <button
            className={styles.navBtn}
            onClick={e => { e.stopPropagation(); prev(); }}
            disabled={currentIndex === 0}
            aria-label="Previous image"
          >←</button>
          <button
            className={`${styles.navBtn} ${styles.navBtnRight}`}
            onClick={e => { e.stopPropagation(); next(); }}
            disabled={currentIndex === images.length - 1}
            aria-label="Next image"
          >→</button>
        </>
      )}

      {images.length > 1 && (
        <div className={styles.dots}>
          {images.map((_, i) => (
            <button
              key={i}
              className={i === currentIndex ? styles.dotActive : styles.dot}
              onClick={e => { e.stopPropagation(); setCurrentIndex(i); }}
              aria-label={`Go to image ${i + 1}`}
            />
          ))}
        </div>
      )}

    </div>,
    document.body
  );
}
