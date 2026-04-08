/**
 * antipratik.com — Data Abstraction Layer
 *
 * Repository pattern: components only ever call these functions.
 * All data is fetched from the Go backend at NEXT_PUBLIC_API_URL.
 *
 * Never call fetch() directly in a component or page.
 */

import type { Post, MusicPost, PhotoPost, VideoPost, LinkPost, EssayPost, ShortPost, ExternalLink, FilterState } from './types';

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

// ─── AUTH ─────────────────────────────────────────────────────────────────────

/**
 * POST /api/auth/login
 * Returns a JWT on success; throws on invalid credentials.
 */
export async function login(username: string, password: string): Promise<{ token: string }> {
  const response = await fetch(`${API_URL}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  });
  if (!response.ok) {
    const body = await response.json().catch(() => ({})) as { error?: string };
    throw new Error(body.error ?? `Login failed (${response.status})`);
  }
  return response.json() as Promise<{ token: string }>;
}

// ─── ADMIN WRITE HELPERS ──────────────────────────────────────────────────────

function authHeaders(token: string): HeadersInit {
  return { Authorization: `Bearer ${token}` };
}

async function throwOnError(response: Response, label: string): Promise<void> {
  if (!response.ok) {
    const body = await response.json().catch(() => ({})) as { error?: string };
    throw new Error(body.error ?? `API error: ${response.status} — ${label}`);
  }
}

// ─── ESSAY ────────────────────────────────────────────────────────────────────

export interface CreateEssayInput {
  title: string;
  slug: string;
  excerpt: string;
  body: string;
  tags?: string[];
}

export async function createEssay(data: CreateEssayInput, token: string): Promise<EssayPost> {
  const response = await fetch(`${API_URL}/api/posts/essay`, {
    method: 'POST',
    headers: { ...authHeaders(token), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  await throwOnError(response, 'createEssay');
  return response.json() as Promise<EssayPost>;
}

export interface UpdateEssayInput {
  title?: string;
  slug?: string;
  excerpt?: string;
  body?: string;
  tags?: string[];
}

export async function updateEssay(id: string, data: UpdateEssayInput, token: string): Promise<EssayPost> {
  const response = await fetch(`${API_URL}/api/posts/essay/${id}`, {
    method: 'PUT',
    headers: { ...authHeaders(token), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  await throwOnError(response, 'updateEssay');
  return response.json() as Promise<EssayPost>;
}

// ─── SHORT POST ───────────────────────────────────────────────────────────────

export interface CreateShortPostInput {
  body: string;
  tags?: string[];
}

export async function createShortPost(data: CreateShortPostInput, token: string): Promise<ShortPost> {
  const response = await fetch(`${API_URL}/api/posts/short`, {
    method: 'POST',
    headers: { ...authHeaders(token), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  await throwOnError(response, 'createShortPost');
  return response.json() as Promise<ShortPost>;
}

export interface UpdateShortPostInput {
  body?: string;
  tags?: string[];
}

export async function updateShortPost(id: string, data: UpdateShortPostInput, token: string): Promise<ShortPost> {
  const response = await fetch(`${API_URL}/api/posts/short/${id}`, {
    method: 'PUT',
    headers: { ...authHeaders(token), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  await throwOnError(response, 'updateShortPost');
  return response.json() as Promise<ShortPost>;
}

// ─── MUSIC ────────────────────────────────────────────────────────────────────

export async function createMusicPost(formData: FormData, token: string): Promise<MusicPost> {
  const response = await fetch(`${API_URL}/api/posts/music`, {
    method: 'POST',
    headers: authHeaders(token),
    body: formData,
  });
  await throwOnError(response, 'createMusicPost');
  return response.json() as Promise<MusicPost>;
}

export async function updateMusicPost(id: string, formData: FormData, token: string): Promise<MusicPost> {
  const response = await fetch(`${API_URL}/api/posts/music/${id}`, {
    method: 'PUT',
    headers: authHeaders(token),
    body: formData,
  });
  await throwOnError(response, 'updateMusicPost');
  return response.json() as Promise<MusicPost>;
}

// ─── PHOTO ────────────────────────────────────────────────────────────────────

export async function createPhotoPost(formData: FormData, token: string): Promise<PhotoPost> {
  const response = await fetch(`${API_URL}/api/posts/photo`, {
    method: 'POST',
    headers: authHeaders(token),
    body: formData,
  });
  await throwOnError(response, 'createPhotoPost');
  return response.json() as Promise<PhotoPost>;
}

export async function updatePhotoPost(id: string, formData: FormData, token: string): Promise<PhotoPost> {
  const response = await fetch(`${API_URL}/api/posts/photo/${id}`, {
    method: 'PUT',
    headers: authHeaders(token),
    body: formData,
  });
  await throwOnError(response, 'updatePhotoPost');
  return response.json() as Promise<PhotoPost>;
}

// ─── VIDEO ────────────────────────────────────────────────────────────────────

export async function createVideoPost(formData: FormData, token: string): Promise<VideoPost> {
  const response = await fetch(`${API_URL}/api/posts/video`, {
    method: 'POST',
    headers: authHeaders(token),
    body: formData,
  });
  await throwOnError(response, 'createVideoPost');
  return response.json() as Promise<VideoPost>;
}

export async function updateVideoPost(id: string, formData: FormData, token: string): Promise<VideoPost> {
  const response = await fetch(`${API_URL}/api/posts/video/${id}`, {
    method: 'PUT',
    headers: authHeaders(token),
    body: formData,
  });
  await throwOnError(response, 'updateVideoPost');
  return response.json() as Promise<VideoPost>;
}

// ─── LINK POST ────────────────────────────────────────────────────────────────

export async function createLinkPost(formData: FormData, token: string): Promise<LinkPost> {
  const response = await fetch(`${API_URL}/api/posts/link`, {
    method: 'POST',
    headers: authHeaders(token),
    body: formData,
  });
  await throwOnError(response, 'createLinkPost');
  return response.json() as Promise<LinkPost>;
}

export async function updateLinkPost(id: string, formData: FormData, token: string): Promise<LinkPost> {
  const response = await fetch(`${API_URL}/api/posts/link/${id}`, {
    method: 'PUT',
    headers: authHeaders(token),
    body: formData,
  });
  await throwOnError(response, 'updateLinkPost');
  return response.json() as Promise<LinkPost>;
}

// ─── EXTERNAL LINK ────────────────────────────────────────────────────────────

export interface CreateExternalLinkInput {
  title: string;
  url: string;
  description: string;
  category: 'music' | 'writing' | 'video' | 'social';
  featured?: boolean;
}

export async function createExternalLink(data: CreateExternalLinkInput, token: string): Promise<ExternalLink> {
  const response = await fetch(`${API_URL}/api/links`, {
    method: 'POST',
    headers: { ...authHeaders(token), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  await throwOnError(response, 'createExternalLink');
  return response.json() as Promise<ExternalLink>;
}

export interface UpdateExternalLinkInput {
  title?: string;
  url?: string;
  description?: string;
  category?: 'music' | 'writing' | 'video' | 'social';
  featured?: boolean;
}

export async function updateExternalLink(id: string, data: UpdateExternalLinkInput, token: string): Promise<ExternalLink> {
  const response = await fetch(`${API_URL}/api/links/${id}`, {
    method: 'PUT',
    headers: { ...authHeaders(token), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  await throwOnError(response, 'updateExternalLink');
  return response.json() as Promise<ExternalLink>;
}

// ─── DELETE ───────────────────────────────────────────────────────────────────

export async function deletePost(id: string, token: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/posts/${id}`, {
    method: 'DELETE',
    headers: authHeaders(token),
  });
  await throwOnError(response, 'deletePost');
}

export async function deleteExternalLink(id: string, token: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/links/${id}`, {
    method: 'DELETE',
    headers: authHeaders(token),
  });
  await throwOnError(response, 'deleteExternalLink');
}
