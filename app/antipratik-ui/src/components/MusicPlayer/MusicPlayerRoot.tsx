'use client';

import { useMusicPlayer } from '../MusicProvider/MusicProvider';
import MusicPlayer from './MusicPlayer';

export default function MusicPlayerRoot() {
  const { activeTrack, isPlaying, isExiting, play, pause, stop } = useMusicPlayer();

  return (activeTrack !== null || isExiting) ? (
    <MusicPlayer
      track={activeTrack!}
      isPlaying={isPlaying}
      isExiting={isExiting}
      onPlay={play}
      onPause={pause}
      onStop={stop}
    />
  ) : <div style={{ display: 'none' }} aria-hidden="true" />;
}
