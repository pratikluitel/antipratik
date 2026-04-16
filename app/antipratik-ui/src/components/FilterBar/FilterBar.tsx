'use client';

import { useLayoutEffect, useRef, useState, useEffect } from 'react';
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

const SCROLL_STEP = 80;

export default function FilterBar({ state, allTags, dispatch }: Props) {
  // 'dark' is the SSR-safe default — useLayoutEffect syncs to data-theme before
  // the browser paints, so light-mode users never see a flash. Driving data-mode
  // from React state (rather than direct setAttribute) means React re-renders
  // can never reset the attribute back to the hardcoded JSX value.
  const [mode, setMode] = useState<'dark' | 'light'>('dark');
  const dropdownRef = useRef<HTMLDivElement>(null);
  const dropdownElRef = useRef<HTMLDivElement>(null);
  const listRef = useRef<HTMLDivElement>(null);
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [isClosing, setIsClosing] = useState(false);
  const [dropdownPos, setDropdownPos] = useState<{ top: number; left: number } | null>(null);
  const [canScrollUp, setCanScrollUp] = useState(false);
  const [canScrollDown, setCanScrollDown] = useState(false);

  // Sync mode with data-theme on <html> (Rule 5). Using setMode keeps
  // data-mode React-controlled so re-renders never reset it to "dark".
  useLayoutEffect(() => {
    function sync() {
      setMode((document.documentElement.getAttribute('data-theme') as 'dark' | 'light') ?? 'dark');
    }
    sync();
    const observer = new MutationObserver(sync);
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] });
    return () => observer.disconnect();
  }, []);

  const ANIMATION_DURATION = 200;

  function closeDropdown() {
    setIsClosing(true);
    setTimeout(() => {
      setDropdownOpen(false);
      setIsClosing(false);
    }, ANIMATION_DURATION);
  }

  // Close dropdown when clicking outside
  useEffect(() => {
    if (!dropdownOpen) return;
    function handleClick(e: MouseEvent) {
      if (
        dropdownRef.current && !dropdownRef.current.contains(e.target as Node) &&
        dropdownElRef.current && !dropdownElRef.current.contains(e.target as Node)
      ) {
        closeDropdown();
      }
    }
    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, [dropdownOpen]);

  // Calculate fixed position for the dropdown — re-runs on open and on window resize
  // so the dropdown tracks the button even if the viewport changes while open.
  useEffect(() => {
    if (!dropdownOpen) return;

    function recalc() {
      const btn = dropdownRef.current?.querySelector('button');
      if (!btn) return;
      const rect = btn.getBoundingClientRect();
      const DROPDOWN_WIDTH = 160;
      const MARGIN = 8;
      let left = rect.left;
      if (left + DROPDOWN_WIDTH > window.innerWidth - MARGIN) {
        left = window.innerWidth - DROPDOWN_WIDTH - MARGIN;
      }
      setDropdownPos({ top: rect.bottom + 6, left: Math.max(MARGIN, left) });
    }

    const raf = requestAnimationFrame(recalc);
    window.addEventListener('resize', recalc);
    window.addEventListener('scroll', recalc, { passive: true });

    return () => {
      cancelAnimationFrame(raf);
      window.removeEventListener('resize', recalc);
      window.removeEventListener('scroll', recalc);
    };
  }, [dropdownOpen]);

  // Update scroll arrow visibility whenever the list scrolls or opens
  function syncScrollArrows() {
    const el = listRef.current;
    if (!el) return;
    setCanScrollUp(el.scrollTop > 0);
    setCanScrollDown(el.scrollTop + el.clientHeight < el.scrollHeight - 1);
  }

  useEffect(() => {
    if (dropdownOpen) {
      requestAnimationFrame(syncScrollArrows);
    }
  }, [dropdownOpen, allTags]);

  function scrollList(direction: 'up' | 'down') {
    const el = listRef.current;
    if (!el) return;
    el.scrollBy({ top: direction === 'up' ? -SCROLL_STEP : SCROLL_STEP, behavior: 'smooth' });
  }

  return (
    <div className={styles.bar} data-mode={mode}>
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

        {allTags.length > 0 && (
          <div ref={dropdownRef} className={styles.tagDropdownWrap}>
            <button
              className={`${styles.pill} ${styles.pillAll} ${styles.tagDropdownBtn}${state.activeTags.length > 0 ? ` ${styles.tagDropdownBtnActive}` : ''}`}
              onClick={() => dropdownOpen ? closeDropdown() : setDropdownOpen(true)}
              aria-expanded={dropdownOpen}
              aria-haspopup="listbox"
            >
              Tags{state.activeTags.length > 0 ? ` (${state.activeTags.length})` : ''}
              <span className={styles.tagDropdownArrow} aria-hidden="true">{(dropdownOpen || isClosing) ? '▲' : '▼'}</span>
            </button>
            {(dropdownOpen || isClosing) && dropdownPos && (
              <div
                ref={dropdownElRef}
                className={`${styles.tagDropdown}${isClosing ? ` ${styles.tagDropdownClosing}` : ''}`}
                style={{ top: dropdownPos.top, left: dropdownPos.left }}
              >
                {/* Mobile-only scroll up arrow */}
                {canScrollUp && (
                  <button
                    className={`${styles.mobileScrollArrow} ${styles.mobileScrollArrowUp}`}
                    onClick={() => scrollList('up')}
                    aria-label="Scroll tags up"
                  >
                    ▲
                  </button>
                )}
                <div
                  ref={listRef}
                  role="listbox"
                  aria-multiselectable="true"
                  className={styles.tagDropdownList}
                  onScroll={syncScrollArrows}
                >
                  {allTags.map((tag) => {
                    const selected = state.activeTags.includes(tag);
                    return (
                      <button
                        key={tag}
                        role="option"
                        aria-selected={selected}
                        className={`${styles.tagOption}${selected ? ` ${styles.tagOptionSelected}` : ''}`}
                        onClick={() => dispatch({ type: 'TOGGLE_TAG', tag })}
                      >
                        {tag}
                      </button>
                    );
                  })}
                </div>
                {/* Scroll down arrow — shown when list has more content below */}
                {canScrollDown && (
                  <button
                    className={`${styles.mobileScrollArrow} ${styles.mobileScrollArrowDown}`}
                    onClick={() => scrollList('down')}
                    aria-label="Scroll tags down"
                  >
                    ▼
                  </button>
                )}
                {/* Clear tags — always pinned at bottom, outside scroll area */}
                {state.activeTags.length > 0 && (
                  <button
                    className={styles.clearTagsDropdown}
                    onClick={() => dispatch({ type: 'CLEAR_TAGS' })}
                  >
                    Clear tags
                  </button>
                )}
              </div>
            )}
          </div>
        )}
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
        </div>
      )}
    </div>
  );
}
