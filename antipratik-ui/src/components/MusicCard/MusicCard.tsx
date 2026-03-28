'use client';

import Image from 'next/image';
import type { MusicPost } from '../../lib/types';
import { useMusicPlayer } from '../MusicProvider/MusicProvider';
import styles from './MusicCard.module.css';

interface Props {
  post: MusicPost;
}

function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${m}:${s.toString().padStart(2, '0')}`;
}

export default function MusicCard({ post }: Props) {
  const { play, isPlaying, activeTrack } = useMusicPlayer();
  const isThisTrackActive = activeTrack?.id === post.id;
  const isThisTrackPlaying = isThisTrackActive && isPlaying;

  const date = new Date(post.createdAt).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  });

  return (
    <article
      className={`${styles.card}${isThisTrackPlaying ? ` ${styles.playing}` : ''}`}
      onClick={() => play(post)}
    >
      <div className={styles.inner}>
        <div className={styles.topRow}>
          {post.albumArt ? (
            <Image
              src={post.albumArt}
              alt={post.title}
              width={52}
              height={52}
              className={styles.albumArt}
            />
          ) : (
            <div className={styles.albumArtPlaceholder} aria-hidden="true">
              ♪
            </div>
          )}
          <div className={styles.textColumn}>
            <span className={styles.tag}>Music</span>
            <h2 className={styles.title}>{post.title}</h2>
            {post.album && <span className={styles.album}>{post.album}</span>}
            <div className={styles.meta}>
              <time className={styles.date} dateTime={post.createdAt}>{date}</time>
              <span className={styles.metaSep}>·</span>
              <span className={styles.duration}>{formatDuration(post.duration)}</span>
            </div>
          </div>
          <button
            className={styles.playButton}
            onClick={(e) => {
              e.stopPropagation();
              play(post);
            }}
            aria-label={isThisTrackPlaying ? `Pause ${post.title}` : `Play ${post.title}`}
          >
            <span className={isThisTrackPlaying ? styles.pauseIcon : styles.playIcon} />
          </button>
        </div>
      </div>
    </article>
  );
}
