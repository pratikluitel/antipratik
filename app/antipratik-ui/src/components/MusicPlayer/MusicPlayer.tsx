'use client';

import { useState, useEffect, useRef } from 'react';
import type { MusicPost } from '../../lib/types';
import { formatTime } from '../../lib/utils';
import Waveform from './Waveform';
import styles from './MusicPlayer.module.css';

interface Props {
  track: MusicPost;
  isPlaying: boolean;
  isExiting: boolean;
  onPlay: (track: MusicPost) => void;
  onPause: () => void;
  onStop: () => void;
}

export default function MusicPlayer({ track, isPlaying, isExiting, onPlay, onPause, onStop }: Props) {
  const [isVisible, setIsVisible] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(track.duration);

  const audioRef = useRef<HTMLAudioElement | null>(null);
  // Stable ref for onStop to avoid stale closure in the 'ended' event listener
  const onStopRef = useRef(onStop);
  useEffect(() => { onStopRef.current = onStop; }, [onStop]);

  const progress = duration > 0 ? currentTime / duration : 0;

  // Entry animation: insert into DOM first, then apply visible class after one frame
  useEffect(() => {
    const id = setTimeout(() => setIsVisible(true), 16);
    return () => clearTimeout(id);
  }, []);

  // Step 1 — Create the audio engine on mount
  useEffect(() => {
    const audio = new Audio();

    audio.addEventListener('timeupdate', () => {
      setCurrentTime(Math.floor(audio.currentTime));
    });
    audio.addEventListener('loadedmetadata', () => {
      if (isFinite(audio.duration)) setDuration(Math.floor(audio.duration));
    });
    audio.addEventListener('ended', () => {
      onStopRef.current();
    });
    audio.addEventListener('error', () => {
      console.warn('Audio error:', audio.error);
    });

    audioRef.current = audio;

    return () => {
      audio.pause();
      audio.removeAttribute('src');
    };
  }, []);

  // Step 2 — Load a new track when track.audioUrl changes (src only; play is Step 3's job)
  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    if (!track.audioUrl) {
      audio.removeAttribute('src');
    } else {
      audio.src = track.audioUrl; // browser auto-loads on src change; no audio.load() needed
    }
    setCurrentTime(0);
    setDuration(track.duration); // reset to dummy until loadedmetadata fires
  }, [track.audioUrl, track.duration]);

  // Step 3 — Sync isPlaying prop → audio element
  // Watches track.audioUrl too so it re-fires on track switch even when isPlaying stays true
  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    if (isPlaying && track.audioUrl) {
      audio.play().catch(err => console.warn('Play failed:', err));
    } else {
      // Always pause: covers explicit pause, track-with-no-audio, and switching tracks
      audio.pause();
    }
  }, [isPlaying, track.audioUrl]);

  // Step 4 — Real seeking
  function seek(fraction: number) {
    const audio = audioRef.current;
    if (!audio || !isFinite(audio.duration)) return;
    const newTime = fraction * audio.duration;
    audio.currentTime = newTime;
    setCurrentTime(Math.floor(newTime));
  }

  function handleProgressClick(e: React.MouseEvent<HTMLDivElement>) {
    const rect = e.currentTarget.getBoundingClientRect();
    seek((e.clientX - rect.left) / rect.width);
  }

  return (
    <div className={`${styles.player}${isVisible && !isExiting ? ` ${styles.visible}` : ''}`}>

      {/* Player Bar — always shown */}
      <div className={styles.bar}>

        {/* Track info — clicking opens drawer */}
        <div className={styles.trackInfo}>
          <div className={styles.albumArt}>
            {track.albumArt
              ? <img src={track.albumArt} alt={track.title} />
              : <span>♪</span>}
          </div>
          <div>
            <div className={styles.trackTitle}>{track.title}</div>
            <div className={styles.trackAlbum}>{track.album}</div>
          </div>
        </div>

        {/* Controls — stopPropagation so they don't open drawer */}
        <div className={styles.controls} onClick={(e) => e.stopPropagation()}>
          <button
            className={styles.playBtn}
            onClick={() => isPlaying ? onPause() : onPlay(track)}
            aria-label={isPlaying ? 'Pause' : 'Play'}
          >
            <span className={isPlaying ? styles.pauseIcon : styles.playIcon} />
          </button>
        </div>

        {/* Progress bar — desktop only, hidden mobile */}
        <div className={styles.progressArea} onClick={(e) => e.stopPropagation()}>
          <span className={styles.timeLabel}>{formatTime(currentTime)}</span>
          <div className={styles.progressTrack} onClick={handleProgressClick}>
            <div className={styles.progressFill} style={{ width: `${progress * 100}%` }} />
          </div>
          <span className={styles.timeLabel}>{formatTime(duration)}</span>
        </div>

      </div>

    </div>
  );
}
