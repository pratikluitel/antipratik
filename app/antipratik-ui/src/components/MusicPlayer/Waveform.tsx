'use client';

import { useRef } from 'react';

interface Props {
  progress: number;
  onScrub: (progress: number) => void;
}

const BAR_COUNT = 35;

const BAR_HEIGHTS: number[] = Array.from({ length: BAR_COUNT }, (_, i) => {
  const h = 20 + (Math.sin(i * 0.8) + Math.sin(i * 1.3)) * 14 + 8;
  return Math.max(8, Math.min(40, h));
});

export default function Waveform({ progress, onScrub }: Props) {
  const containerRef = useRef<HTMLDivElement>(null);
  const playedCount = Math.floor(progress * BAR_COUNT);

  function handleClick(e: React.MouseEvent<HTMLDivElement>) {
    if (!containerRef.current) return;
    const rect = containerRef.current.getBoundingClientRect();
    const fraction = (e.clientX - rect.left) / rect.width;
    onScrub(Math.max(0, Math.min(1, fraction)));
  }

  return (
    <div
      ref={containerRef}
      onClick={handleClick}
      style={{
        height: 'var(--waveform-height)',
        display: 'flex',
        alignItems: 'center',
        gap: 'var(--waveform-bar-gap)',
        cursor: 'pointer',
      }}
    >
      {BAR_HEIGHTS.map((h, i) => (
        <div
          key={i}
          style={{
            width: 'var(--waveform-bar-width)',
            height: `${h}px`,
            borderRadius: '2px',
            flexShrink: 0,
            background: i < playedCount
              ? 'var(--accent-music)'
              : 'var(--color-elevated-dark)',
          }}
        />
      ))}
    </div>
  );
}
