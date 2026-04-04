/**
 * antipratik.com — Data Abstraction Layer
 *
 * Repository pattern: components only ever call these functions.
 * All data is fetched from the Go backend at NEXT_PUBLIC_API_URL.
 *
 * Never call fetch() directly in a component or page.
 */

import type { Post, EssayPost, ExternalLink, FilterState } from './types';

const API_URL = process.env.NEXT_PUBLIC_API_URL;
const IS_API_DISABLED = !API_URL;

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
  return response.json() as Promise<Post[]>;
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
