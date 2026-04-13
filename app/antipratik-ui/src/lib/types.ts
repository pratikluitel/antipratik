/**
 * antipratik.com — TypeScript Type System
 * Single source of truth for all data interfaces.
 * Designed to match what the Go REST backend will return.
 */

export type ContentType = 'essay' | 'short' | 'music' | 'photo' | 'video' | 'link';

export type Theme = 'dark' | 'light';

// ─── BASE ────────────────────────────────────────────────────────────────────

interface BasePost {
  id: string;
  type: ContentType;
  createdAt: string; // ISO 8601
  tags: string[];    // For ShortPost, tags serve as hashtags
}

// ─── CONTENT TYPES ───────────────────────────────────────────────────────────

export interface EssayPost extends BasePost {
  type: 'essay';
  title: string;
  slug: string;
  excerpt: string;
  body: string; // markdown
  readingTimeMinutes: number;
}

export interface ShortPost extends BasePost {
  type: 'short';
  body: string;
  // tags from BasePost serve as hashtags — no separate field needed
}

export interface MusicPost extends BasePost {
  type: 'music';
  title: string;
  albumArt: string; // URL
  albumArtTinyUrl?: string; // 20px wide — LQIP blur placeholder
  audioUrl: string;
  duration: number; // seconds
  album?: string;
}

export interface PhotoImage {
  url: string;                // original — served via GET /files/{fileId}
  alt: string;
  caption?: string;
  thumbnailTinyUrl?: string;   // 20px wide — used as LQIP blur placeholder
  thumbnailSmallUrl?: string;  // 300px wide — served via GET /thumbnails/{id}-small.ext
  thumbnailMediumUrl?: string; // 600px wide
  thumbnailLargeUrl?: string;  // 1200px wide
}

export interface PhotoPost extends BasePost {
  type: 'photo';
  images: PhotoImage[];
  location?: string;
}

export interface VideoPost extends BasePost {
  type: 'video';
  title: string;
  thumbnailUrl: string;
  thumbnailTinyUrl?: string; // 20px wide — LQIP blur placeholder
  videoUrl: string;
  duration: number; // seconds
  playlist?: string;
}

export interface LinkPost extends BasePost {
  type: 'link';
  title: string;
  url: string;
  domain: string;
  description?: string;
  thumbnailUrl?: string;
  thumbnailTinyUrl?: string; // 20px wide — LQIP blur placeholder
  category?: 'music' | 'writing' | 'video' | 'social';
}

export type Post = EssayPost | ShortPost | MusicPost | PhotoPost | VideoPost | LinkPost;

// ─── EXTERNAL LINKS ──────────────────────────────────────────────────────────

export interface ExternalLink {
  id: string;
  title: string;
  url: string;
  domain: string;
  description: string;
  featured: boolean;
  category: 'music' | 'writing' | 'video' | 'social';
}

// ─── MUSIC PLAYER ────────────────────────────────────────────────────────────

export interface Track {
  id: string;
  title: string;
  albumArt: string;
  audioUrl: string;
  duration: number; // seconds
  album?: string;
}

// ─── FEED LAYOUT ─────────────────────────────────────────────────────────────

export type LayoutMode = 'reading' | 'visual';

// Maps each ContentType to its layout mode
// reading: essay, short
// visual: music, photo, video, link
export const LAYOUT_MODE: Record<ContentType, LayoutMode> = {
  essay: 'reading',
  short: 'reading',
  music: 'visual',
  photo: 'visual',
  video: 'visual',
  link: 'visual',
};

export type FeedItem =
  | { kind: 'post'; post: Post }
  | { kind: 'divider'; from: LayoutMode; to: LayoutMode }
  | { kind: 'date'; date: string };

// ─── FEED FILTER ─────────────────────────────────────────────────────────────

export interface FilterState {
  activeTypes: ContentType[];
  activeTags: string[];
  sortOrder: 'newest' | 'oldest';
}

export type FilterAction =
  | { type: 'TOGGLE_TYPE'; contentType: ContentType }
  | { type: 'TOGGLE_TAG'; tag: string }
  | { type: 'SET_SORT'; order: FilterState['sortOrder'] }
  | { type: 'CLEAR_ALL' };
