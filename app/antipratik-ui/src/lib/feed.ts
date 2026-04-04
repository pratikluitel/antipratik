/**
 * antipratik.com — Feed Logic Utilities
 * Pure functions for filtering, sorting, and clustering feed posts.
 * No React imports — usable in server and client contexts.
 */

import type { Post, FilterState, FilterAction, FeedItem, ContentType } from './types';
import { LAYOUT_MODE } from './types';

export const initialFilterState: FilterState = {
  activeTypes: [],
  activeTags: [],
  sortOrder: 'newest',
};

const ALL_TYPES: ContentType[] = ['essay', 'short', 'music', 'photo', 'video', 'link'];

export function filterReducer(state: FilterState, action: FilterAction): FilterState {
  switch (action.type) {
    case 'TOGGLE_TYPE': {
      const isActive = state.activeTypes.includes(action.contentType);
      let next: ContentType[];
      if (isActive) {
        next = state.activeTypes.filter((t) => t !== action.contentType);
      } else {
        next = [...state.activeTypes, action.contentType];
      }
      // If all types would be active, normalize to [] (same as "all")
      if (next.length === ALL_TYPES.length) next = [];
      return { ...state, activeTypes: next };
    }
    case 'TOGGLE_TAG': {
      const isActive = state.activeTags.includes(action.tag);
      const next = isActive
        ? state.activeTags.filter((t) => t !== action.tag)
        : [...state.activeTags, action.tag];
      return { ...state, activeTags: next };
    }
    case 'SET_SORT':
      return { ...state, sortOrder: action.order };
    case 'CLEAR_ALL':
      return { activeTypes: [], activeTags: [], sortOrder: 'newest' };
  }
}

export function applyFilters(posts: Post[], state: FilterState): Post[] {
  let result = posts;

  if (state.activeTypes.length > 0) {
    result = result.filter((p) => state.activeTypes.includes(p.type));
  }

  if (state.activeTags.length > 0) {
    result = result.filter((p) =>
      p.tags.some((tag) => state.activeTags.includes(tag))
    );
  }

  result = [...result].sort((a, b) => {
    const aTime = new Date(a.createdAt).getTime();
    const bTime = new Date(b.createdAt).getTime();
    return state.sortOrder === 'newest' ? bTime - aTime : aTime - bTime;
  });

  return result;
}

export function buildFeedClusters(posts: Post[]): FeedItem[] {
  const items: FeedItem[] = [];
  let currentMode: 'reading' | 'visual' | null = null;
  let lastVisualDate: string | null = null;

  for (const post of posts) {
    const layoutMode = LAYOUT_MODE[post.type];

    if (currentMode !== null && layoutMode !== currentMode) {
      items.push({ kind: 'divider', from: currentMode, to: layoutMode });
    }

    if (layoutMode === 'visual') {
      const date = post.createdAt.slice(0, 10); // YYYY-MM-DD
      if (date !== lastVisualDate) {
        items.push({ kind: 'date', date: post.createdAt });
        lastVisualDate = date;
      }
    }

    items.push({ kind: 'post', post });
    currentMode = layoutMode;
  }

  return items;
}
