'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import type { MusicPost } from '../../lib/types';
import { formatTime } from '../../lib/utils';
import styles from './MusicPlayer.module.css';

function IconVolumeMuted() {
  return (
    <svg width="18" height="18" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
      <path d="M3 7h3l5-4v14l-5-4H3V7z" />
      <line x1="13" y1="7" x2="19" y2="13" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
      <line x1="19" y1="7" x2="13" y2="13" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
    </svg>
  );
}

function IconVolumeLow() {
  return (
    <svg width="18" height="18" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
      <path d="M3 7h3l5-4v14l-5-4H3V7z" />
      <path d="M14 7a4 4 0 0 1 0 6" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" fill="none" />
    </svg>
  );
}

function IconVolumeHigh() {
  return (
    <svg width="18" height="18" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
      <path d="M3 7h3l5-4v14l-5-4H3V7z" />
      <path d="M14 7a4 4 0 0 1 0 6" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" fill="none" />
      <path d="M16.5 4.5a7.5 7.5 0 0 1 0 11" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" fill="none" />
    </svg>
  );
}

function VolumeIcon({ volume, muted }: { volume: number; muted: boolean }) {
  if (muted || volume === 0) return <IconVolumeMuted />;
  if (volume < 0.5) return <IconVolumeLow />;
  return <IconVolumeHigh />;
}

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
  const [volume, setVolume] = useState(1);
  const [isMuted, setIsMuted] = useState(false);

  const audioRef = useRef<HTMLAudioElement | null>(null);
  const [albumArtError, setAlbumArtError] = useState(false);
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
    audio.addEventListener('volumechange', () => {
      setIsMuted(audio.muted);
      setVolume(audio.volume);
    });

    audioRef.current = audio;

    return () => {
      audio.pause();
      audio.removeAttribute('src');
    };
  }, []);

  // Step 2 — Load a new track when track.audioUrl changes (src only; play is Step 3's job)
  const [prevAudioUrl, setPrevAudioUrl] = useState(track.audioUrl);
  if (prevAudioUrl !== track.audioUrl) {
    setPrevAudioUrl(track.audioUrl);
    setCurrentTime(0);
    setDuration(track.duration); // DB value as placeholder until loadedmetadata fires
  }

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    if (!track.audioUrl) {
      audio.removeAttribute('src');
    } else {
      audio.src = track.audioUrl; // browser auto-loads on src change; no audio.load() needed
    }
  }, [track.audioUrl]);

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

  const toggleMute = useCallback(() => {
    const audio = audioRef.current;
    if (!audio) return;
    if (audio.muted) {
      audio.muted = false;
      if (audio.volume === 0) audio.volume = 0.5;
    } else {
      audio.muted = true;
    }
  }, []);

  const handleVolumeChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const audio = audioRef.current;
    if (!audio) return;
    const val = Number(e.target.value);
    audio.volume = val;
    audio.muted = val === 0;
  }, []);

  const volumeFill = isMuted ? 0 : volume * 100;

  return (
    <div className={`${styles.player}${isVisible && !isExiting ? ` ${styles.visible}` : ''}`}>

      {/* Player Bar — always shown */}
      <div className={styles.bar}>

        {/* Track info — clicking opens drawer */}
        <div className={styles.trackInfo}>
          <div className={styles.albumArt}>
            {track.albumArt
              // eslint-disable-next-line @next/next/no-img-element
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
          <div className={styles.volumeGroup}>
            <button
              type="button"
              className={styles.volumeBtn}
              onClick={toggleMute}
              aria-label={isMuted ? 'Unmute' : 'Mute'}
            >
              <VolumeIcon volume={volume} muted={isMuted} />
            </button>
            <div className={styles.volumeSliderWrapper}>
              <input
                type="range"
                className={styles.volumeSlider}
                min={0}
                max={1}
                step={0.02}
                value={isMuted ? 0 : volume}
                onChange={handleVolumeChange}
                aria-label="Volume"
                style={{ '--volume-fill': `${volumeFill}%` } as React.CSSProperties}
              />
            </div>
          </div>
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
