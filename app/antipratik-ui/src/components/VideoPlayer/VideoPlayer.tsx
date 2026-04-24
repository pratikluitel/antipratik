'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import { createPortal } from 'react-dom';
import styles from './VideoPlayer.module.css';

type Quality = 'low' | 'med' | 'high';

function qualityUrl(base: string, quality: Quality): string {
  if (quality === 'med') return base;
  return `${base}?q=${quality}`;
}

// SVG icons (20×20 viewBox)
function IconPlay() {
  return (
    <svg width="16" height="16" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
      <polygon points="4,2 18,10 4,18" />
    </svg>
  );
}

function IconPause() {
  return (
    <svg width="16" height="16" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
      <rect x="3" y="2" width="5" height="16" rx="1" />
      <rect x="12" y="2" width="5" height="16" rx="1" />
    </svg>
  );
}

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

function IconEnterFullscreen() {
  return (
    <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" aria-hidden="true">
      <polyline points="2,8 2,2 8,2" />
      <polyline points="12,2 18,2 18,8" />
      <polyline points="18,12 18,18 12,18" />
      <polyline points="8,18 2,18 2,12" />
    </svg>
  );
}

function IconExitFullscreen() {
  return (
    <svg width="16" height="16" viewBox="0 0 20 20" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" aria-hidden="true">
      <polyline points="8,2 8,8 2,8" />
      <polyline points="18,8 12,8 12,2" />
      <polyline points="12,18 12,12 18,12" />
      <polyline points="2,12 8,12 8,18" />
    </svg>
  );
}

function VolumeIcon({ volume, muted }: { volume: number; muted: boolean }) {
  if (muted || volume === 0) return <IconVolumeMuted />;
  if (volume < 0.5) return <IconVolumeLow />;
  return <IconVolumeHigh />;
}

interface Props {
  videoUrl: string;
  title?: string;
  onClose: () => void;
}

export default function VideoPlayer({ videoUrl, title, onClose }: Props) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [isMuted, setIsMuted] = useState(false);
  const [volume, setVolume] = useState(1);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [quality, setQuality] = useState<Quality>('med');

  // Lock body scroll while modal is open
  useEffect(() => {
    const prev = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    return () => { document.body.style.overflow = prev; };
  }, []);

  // Close on Escape key (only when not in fullscreen — browser handles that natively)
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && !document.fullscreenElement) onClose();
    };
    document.addEventListener('keydown', handler);
    return () => document.removeEventListener('keydown', handler);
  }, [onClose]);

  // Track fullscreen state
  useEffect(() => {
    const handler = () => setIsFullscreen(!!document.fullscreenElement);
    document.addEventListener('fullscreenchange', handler);
    return () => document.removeEventListener('fullscreenchange', handler);
  }, []);

  // Wire video events
  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    const onTimeUpdate = () => setCurrentTime(video.currentTime);
    const onDurationChange = () => setDuration(video.duration);
    const onPlay = () => setIsPlaying(true);
    const onPause = () => setIsPlaying(false);
    const onVolumeChange = () => {
      setIsMuted(video.muted);
      setVolume(video.volume);
    };

    video.addEventListener('timeupdate', onTimeUpdate);
    video.addEventListener('durationchange', onDurationChange);
    video.addEventListener('play', onPlay);
    video.addEventListener('pause', onPause);
    video.addEventListener('volumechange', onVolumeChange);

    // Attempt autoplay; fall back to muted if blocked
    video.play().catch(() => {
      video.muted = true;
      setIsMuted(true);
      video.play().catch(() => {});
    });

    return () => {
      video.removeEventListener('timeupdate', onTimeUpdate);
      video.removeEventListener('durationchange', onDurationChange);
      video.removeEventListener('play', onPlay);
      video.removeEventListener('pause', onPause);
      video.removeEventListener('volumechange', onVolumeChange);
    };
  }, []);

  const togglePlay = useCallback(() => {
    const video = videoRef.current;
    if (!video) return;
    if (video.paused) { video.play().catch(() => {}); } else { video.pause(); }
  }, []);

  const handleSeek = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const video = videoRef.current;
    if (!video) return;
    video.currentTime = Number(e.target.value);
  }, []);

  const handleVolumeChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const video = videoRef.current;
    if (!video) return;
    const val = Number(e.target.value);
    video.volume = val;
    video.muted = val === 0;
  }, []);

  const toggleMute = useCallback(() => {
    const video = videoRef.current;
    if (!video) return;
    if (video.muted) {
      video.muted = false;
      if (video.volume === 0) video.volume = 0.5;
    } else {
      video.muted = true;
    }
  }, []);

  const handleFullscreen = useCallback(() => {
    const container = containerRef.current;
    if (!container) return;
    if (document.fullscreenElement) {
      document.exitFullscreen().catch(() => {});
    } else {
      container.requestFullscreen().catch(() => {});
    }
  }, []);

  const handleQualityChange = useCallback((next: Quality) => {
    const video = videoRef.current;
    if (!video) return;
    const savedTime = video.currentTime;
    const wasPlaying = !video.paused;
    setQuality(next);
    video.src = qualityUrl(videoUrl, next);
    video.load();
    video.addEventListener('canplay', () => {
      video.currentTime = savedTime;
      if (wasPlaying) video.play().catch(() => {});
    }, { once: true });
  }, [videoUrl]);

  const seekFill = duration > 0 ? (currentTime / duration) * 100 : 0;
  const volumeFill = isMuted ? 0 : volume * 100;

  const modal = (
    <div
      className={styles.overlay}
      onClick={onClose}
      role="dialog"
      aria-modal="true"
      aria-label={title ?? 'Video player'}
    >
      <div ref={containerRef} className={styles.container} onClick={(e) => e.stopPropagation()}>
        {title && <div className={styles.title}>{title}</div>}

        <video
          ref={videoRef}
          className={styles.video}
          src={qualityUrl(videoUrl, quality)}
          playsInline
        />

        <div className={styles.controls}>
          <button
            type="button"
            className={styles.controlBtn}
            onClick={togglePlay}
            aria-label={isPlaying ? 'Pause' : 'Play'}
          >
            {isPlaying ? <IconPause /> : <IconPlay />}
          </button>

          <div className={styles.seekWrapper}>
            <input
              type="range"
              className={styles.seek}
              min={0}
              max={duration || 0}
              step={0.1}
              value={currentTime}
              onChange={handleSeek}
              aria-label="Seek"
              style={{ '--seek-fill': `${seekFill}%` } as React.CSSProperties}
            />
          </div>

          <div className={styles.volumeGroup}>
            <button
              type="button"
              className={styles.controlBtn}
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

          <select
            className={styles.qualitySelect}
            value={quality}
            onChange={(e) => handleQualityChange(e.target.value as Quality)}
            aria-label="Quality"
          >
            <option value="low">Low</option>
            <option value="med">Med</option>
            <option value="high">High</option>
          </select>

          <button
            type="button"
            className={styles.controlBtn}
            onClick={handleFullscreen}
            aria-label={isFullscreen ? 'Exit fullscreen' : 'Fullscreen'}
          >
            {isFullscreen ? <IconExitFullscreen /> : <IconEnterFullscreen />}
          </button>
        </div>
      </div>
    </div>
  );

  return createPortal(modal, document.body);
}
