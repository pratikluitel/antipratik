/**
 * antipratik.com — Data Abstraction Layer
 *
 * Repository pattern: components only ever call these functions.
 * When NEXT_PUBLIC_API_URL is set, functions fetch from the Go backend.
 * When not set, they return dummy data for local development.
 *
 * Never call fetch() directly in a component or page.
 * Never import from dummy-data/ directly in a component or page.
 */

import type { Post, EssayPost, ExternalLink, FilterState } from './types';

const API_URL = process.env.NEXT_PUBLIC_API_URL;

// ─── POSTS ───────────────────────────────────────────────────────────────────

/**
 * GET /api/posts
 * Returns all posts, newest first. Optionally filtered by content type and tags.
 *
 * @param filter - Optional filter state with activeTypes and activeTags
 */
export async function getPosts(filter?: FilterState): Promise<Post[]> {
  if (API_URL) {
    const params = new URLSearchParams();
    if (filter?.activeTypes?.length) {
      filter.activeTypes.forEach((t) => params.append('type', t));
    }
    if (filter?.activeTags?.length) {
      filter.activeTags.forEach((tag) => params.append('tag', tag));
    }
    const query = params.toString() ? `?${params.toString()}` : '';
    const response = await fetch(`${API_URL}/api/posts${query}`, {
      next: { revalidate: 60 },
    });
    if (!response.ok) {
      throw new Error(`API error: ${response.status} ${response.statusText} — getPosts`);
    }
    return response.json() as Promise<Post[]>;
  }

  const { posts } = await import('./dummy-data/posts');
  if (!filter || (!filter.activeTypes.length && !filter.activeTags.length)) {
    return posts;
  }
  return posts.filter((post) => {
    const typeMatch =
      !filter.activeTypes.length || filter.activeTypes.includes(post.type);
    const tagMatch =
      !filter.activeTags.length ||
      filter.activeTags.some((tag) => post.tags.includes(tag));
    return typeMatch && tagMatch;
  });
}

/**
 * GET /api/posts/:slug
 * Returns a single essay post by slug. Returns null if not found.
 *
 * @param slug - The URL slug of the essay
 */
export async function getPost(slug: string): Promise<EssayPost | null> {
  if (API_URL) {
    const response = await fetch(`${API_URL}/api/posts/${slug}`, {
      next: { revalidate: 300 },
    });
    if (response.status === 404) return null;
    if (!response.ok) {
      throw new Error(`API error: ${response.status} ${response.statusText} — getPost(${slug})`);
    }
    return response.json() as Promise<EssayPost>;
  }

  const { posts } = await import('./dummy-data/posts');
  const post = posts.find((p) => p.type === 'essay' && (p as EssayPost).slug === slug);
  return (post as EssayPost) ?? null;
}

// ─── LINKS ───────────────────────────────────────────────────────────────────

/**
 * GET /api/links
 * Returns all external links.
 */
export async function getLinks(): Promise<ExternalLink[]> {
  if (API_URL) {
    const response = await fetch(`${API_URL}/api/links`, {
      next: { revalidate: 3600 },
    });
    if (!response.ok) {
      throw new Error(`API error: ${response.status} ${response.statusText} — getLinks`);
    }
    return response.json() as Promise<ExternalLink[]>;
  }

  const { links } = await import('./dummy-data/links');
  return links;
}

/**
 * GET /api/links/featured
 * Returns up to 4 featured external links for the homepage snippet.
 */
export async function getFeaturedLinks(): Promise<ExternalLink[]> {
  if (API_URL) {
    const response = await fetch(`${API_URL}/api/links/featured`, {
      next: { revalidate: 3600 },
    });
    if (!response.ok) {
      throw new Error(`API error: ${response.status} ${response.statusText} — getFeaturedLinks`);
    }
    return response.json() as Promise<ExternalLink[]>;
  }

  const { links } = await import('./dummy-data/links');
  return links.filter((l) => l.featured).slice(0, 4);
}
