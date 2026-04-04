'use client';

import { createContext, useContext, useState, useRef, useEffect } from 'react';
import type { MusicPost } from '../../lib/types';

interface MusicContextValue {
  activeTrack: MusicPost | null;
  isPlaying: boolean;
  isExiting: boolean;
  play: (track: MusicPost) => void;
  pause: () => void;
  stop: () => void;
}

const MusicContext = createContext<MusicContextValue | null>(null);

export function useMusicPlayer(): MusicContextValue {
  const ctx = useContext(MusicContext);
  if (!ctx) {
    throw new Error('useMusicPlayer must be used within a MusicProvider');
  }
  return ctx;
}

export default function MusicProvider({ children }: { children: React.ReactNode }) {
  const [activeTrack, setActiveTrack] = useState<MusicPost | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [isExiting, setIsExiting] = useState(false);
  const exitTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    return () => {
      if (exitTimerRef.current) clearTimeout(exitTimerRef.current);
    };
  }, []);

  function play(track: MusicPost) {
    if (activeTrack?.id === track.id) {
      setIsPlaying((prev) => !prev);
    } else {
      setActiveTrack(track);
      setIsPlaying(true);
    }
  }

  function pause() {
    setIsPlaying(false);
  }

  function stop() {
    setIsExiting(true);
    setIsPlaying(false);
    if (exitTimerRef.current) clearTimeout(exitTimerRef.current);
    exitTimerRef.current = setTimeout(() => {
      setActiveTrack(null);
      setIsExiting(false);
    }, 400);
  }

  return (
    <MusicContext.Provider value={{ activeTrack, isPlaying, isExiting, play, pause, stop }}>
      {children}
    </MusicContext.Provider>
  );
}
