/**
 * antipratik.com — Data Abstraction Layer
 *
 * Repository pattern: components only ever call these functions.
 * All data is fetched from the Go backend at NEXT_PUBLIC_API_URL.
 *
 * Never call fetch() directly in a component or page.
 */

import type { Post, MusicPost, PhotoPost, VideoPost, LinkPost, EssayPost, ExternalLink, FilterState } from './types';

const API_URL = process.env.NEXT_PUBLIC_API_URL;
const IS_API_DISABLED = !API_URL;

// ─── URL PREFIXING ────────────────────────────────────────────────────────────
// The backend stores relative URLs (e.g. /files/abc.jpg, /thumbnails/abc-small.jpg).
// Prefix them with the API base URL so the browser can resolve them correctly.

function prefixUrl(url: string | undefined | null): string {
  if (!url) return '';
  if (url.startsWith('http') || url.startsWith('https') || url.startsWith('//')) return url;
  return url.startsWith('/') ? `${API_URL}${url}` : `${API_URL}/${url}`;
}

function prefixOptionalUrl(url: string | undefined | null): string | undefined {
  if (!url) return undefined;
  if (url.startsWith('http') || url.startsWith('https') || url.startsWith('//')) return url;
  return url.startsWith('/') ? `${API_URL}${url}` : `${API_URL}/${url}`;
}

function prefixPost(post: Post): Post {
  switch (post.type) {
    case 'music': {
      const p = post as MusicPost;
      return { ...p, albumArt: prefixUrl(p.albumArt), audioUrl: prefixUrl(p.audioUrl) };
    }
    case 'photo': {
      const p = post as PhotoPost;
      return {
        ...p,
        images: p.images.map((img) => ({
          ...img,
          url: prefixUrl(img.url),
          thumbnailSmallUrl: prefixOptionalUrl(img.thumbnailSmallUrl),
          thumbnailMediumUrl: prefixOptionalUrl(img.thumbnailMediumUrl),
          thumbnailLargeUrl: prefixOptionalUrl(img.thumbnailLargeUrl),
        })),
      };
    }
    case 'video': {
      const p = post as VideoPost;
      return { ...p, thumbnailUrl: prefixUrl(p.thumbnailUrl) };
    }
    case 'link': {
      const p = post as LinkPost;
      return { ...p, thumbnailUrl: prefixOptionalUrl(p.thumbnailUrl) };
    }
    default:
      return post;
  }
}

// During isolated UI builds without a backend, we still compile the app.
// Pages will render with empty collections, and actual data loading requires
// NEXT_PUBLIC_API_URL to be configured at runtime.

// ─── POSTS ───────────────────────────────────────────────────────────────────

/**
 * GET /api/posts
 * Returns all posts, newest first. Optionally filtered by content type and tags.
 *
 * @param filter - Optional filter state with activeTypes and activeTags
 */
export async function getPosts(filter?: FilterState): Promise<Post[]> {
  if (IS_API_DISABLED) {
    return [];
  }

  const params = new URLSearchParams();
  if (filter?.activeTypes?.length) {
    filter.activeTypes.forEach((t) => params.append('type', t));
  }
  if (filter?.activeTags?.length) {
    filter.activeTags.forEach((tag) => params.append('tag', tag));
  }
  const query = params.toString() ? `?${params.toString()}` : '';
  const response = await fetch(`${API_URL}/api/posts${query}`);
  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText} — getPosts`);
  }
  const posts: Post[] = await response.json();
  return posts.map(prefixPost);
}

/**
 * GET /api/posts/:slug
 * Returns a single essay post by slug. Returns null if not found.
 *
 * @param slug - The URL slug of the essay
 */
export async function getPost(slug: string): Promise<EssayPost | null> {
  if (IS_API_DISABLED) {
    return null;
  }

  const response = await fetch(`${API_URL}/api/posts/${slug}`);
  if (response.status === 404) return null;
  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText} — getPost(${slug})`);
  }
  return response.json() as Promise<EssayPost>;
}

// ─── LINKS ───────────────────────────────────────────────────────────────────

/**
 * GET /api/links
 * Returns all external links.
 */
export async function getLinks(): Promise<ExternalLink[]> {
  if (IS_API_DISABLED) {
    return [];
  }

  const response = await fetch(`${API_URL}/api/links`);
  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText} — getLinks`);
  }
  return response.json() as Promise<ExternalLink[]>;
}

/**
 * GET /api/links/featured
 * Returns up to 4 featured external links for the homepage snippet.
 */
export async function getFeaturedLinks(): Promise<ExternalLink[]> {
  if (IS_API_DISABLED) {
    return [];
  }

  const response = await fetch(`${API_URL}/api/links/featured`);
  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText} — getFeaturedLinks`);
  }
  return response.json() as Promise<ExternalLink[]>;
}
