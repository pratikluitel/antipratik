'use client';

import Image from 'next/image';
import { useState } from 'react';
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

function AlbumArtImage({ post }: { post: MusicPost }) {
  const [loaded, setLoaded] = useState(false);
  return (
    <div className={styles.albumArtPanel}>
      {post.albumArtTinyUrl && !loaded && (
        <div
          className={styles.albumArtPlaceholder_lqip}
          style={{ backgroundImage: `url(${post.albumArtTinyUrl})` }}
          aria-hidden="true"
        />
      )}
      <Image
        src={post.albumArt}
        alt={post.title}
        fill
        onLoad={() => setLoaded(true)}
        className={`${styles.albumArt} ${loaded ? styles.imageVisible : styles.imageHidden}`}
      />
    </div>
  );
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
      {post.albumArt ? (
        <AlbumArtImage post={post} />
      ) : (
        <div className={styles.albumArtPlaceholder} aria-hidden="true">
          ♪
        </div>
      )}
      <div className={styles.content}>
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
    </article>
  );
}
