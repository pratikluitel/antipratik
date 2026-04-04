'use client';

import { useEffect, useRef } from 'react';
import type { FilterState, FilterAction, ContentType } from '../../lib/types';
import styles from './FilterBar.module.css';

interface Props {
  state: FilterState;
  allTags: string[];
  dispatch: (action: FilterAction) => void;
}

const CONTENT_TYPES: { type: ContentType; label: string; pillClass: string }[] = [
  { type: 'essay',  label: 'Essays',  pillClass: styles.pillEssay },
  { type: 'short',  label: 'Short',   pillClass: styles.pillShort },
  { type: 'music',  label: 'Music',   pillClass: styles.pillMusic },
  { type: 'photo',  label: 'Photos',  pillClass: styles.pillPhoto },
  { type: 'video',  label: 'Videos',  pillClass: styles.pillVideo },
  { type: 'link',   label: 'Links',   pillClass: styles.pillLink  },
];

export default function FilterBar({ state, allTags: _allTags, dispatch }: Props) {
  const barRef = useRef<HTMLDivElement>(null);

  // Mirror data-theme on <html> to data-mode on this element (Rule 5)
  useEffect(() => {
    function sync() {
      const theme = document.documentElement.getAttribute('data-theme') ?? 'dark';
      barRef.current?.setAttribute('data-mode', theme);
    }
    sync();
    const observer = new MutationObserver(sync);
    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['data-theme'],
    });
    return () => observer.disconnect();
  }, []);

  return (
    <div ref={barRef} className={styles.bar} data-mode="dark">
      <div className={styles.pillsRow}>
        <button
          className={`${styles.pill} ${styles.pillAll}${state.activeTypes.length === 0 ? ` ${styles.selected}` : ''}`}
          onClick={() => dispatch({ type: 'CLEAR_ALL' })}
        >
          All
        </button>
        {CONTENT_TYPES.map(({ type, label, pillClass }) => (
          <button
            key={type}
            className={`${styles.pill} ${pillClass}${state.activeTypes.includes(type) ? ` ${styles.selected}` : ''}`}
            onClick={() => dispatch({ type: 'TOGGLE_TYPE', contentType: type })}
          >
            {label}
          </button>
        ))}
      </div>

      {state.activeTags.length > 0 && (
        <div className={styles.tagsRow}>
          {state.activeTags.map((tag) => (
            <span key={tag} className={styles.chip}>
              {tag}
              <button
                className={styles.chipRemove}
                onClick={() => dispatch({ type: 'TOGGLE_TAG', tag })}
                aria-label={`Remove ${tag} filter`}
              >
                ×
              </button>
            </span>
          ))}
          <button
            className={styles.clearTags}
            onClick={() => dispatch({ type: 'CLEAR_ALL' })}
          >
            Clear all
          </button>
        </div>
      )}
    </div>
  );
}
