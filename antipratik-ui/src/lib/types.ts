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
  audioUrl: string;
  duration: number; // seconds
  album?: string;
}

export interface PhotoPost extends BasePost {
  type: 'photo';
  images: Array<{
    url: string;
    alt: string;
    caption?: string;
  }>;
  location?: string;
}

export interface VideoPost extends BasePost {
  type: 'video';
  title: string;
  thumbnailUrl: string;
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
}

export type Post = EssayPost | ShortPost | MusicPost | PhotoPost | VideoPost | LinkPost;

// ─── EXTERNAL LINKS ──────────────────────────────────────────────────────────

export interface ExternalLink {
  id: string;
  title: string;
  url: string;
  domain: string;
  description: string;
  iconUrl?: string;
  featured: boolean;
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

// ─── FEED FILTER ─────────────────────────────────────────────────────────────

export interface FilterState {
  activeTypes: ContentType[];
  activeTags: string[];
}
