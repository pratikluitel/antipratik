/**
 * antipratik.com — Data Abstraction Layer
 *
 * Repository pattern: components only ever call these functions.
 * All data is fetched from the Go backend at NEXT_PUBLIC_API_URL.
 *
 * Never call fetch() directly in a component or page.
 */

import type { Post, MusicPost, PhotoPost, PhotoImage, VideoPost, LinkPost, EssayPost, ShortPost, ExternalLink, FilterState, SubscriberSummary, BroadcastSummary, BroadcastPreview, CreateBroadcastInput, UpdateBroadcastInput, BroadcastSendDetail } from './types';

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? '';
// Server-side internal URL (not exposed to browser). Set SERVER_API_URL in the
// container environment (e.g. http://api:8080) so SSR can reach the API over
// the Docker internal network while the browser uses relative URLs via nginx.
const SERVER_API_URL = process.env.SERVER_API_URL ?? '';

function getFetchBase(): string {
  if (typeof window === 'undefined' && SERVER_API_URL) {
    return SERVER_API_URL;
  }
  return API_URL;
}

// Disable API calls only when there is genuinely no URL available to fetch from.
const IS_API_DISABLED =
  typeof window === 'undefined'
    ? !SERVER_API_URL && !API_URL.startsWith('http')
    : false;

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
      return { ...p, albumArt: prefixUrl(p.albumArt), albumArtTinyUrl: prefixOptionalUrl(p.albumArtTinyUrl), albumArtSmallUrl: prefixOptionalUrl(p.albumArtSmallUrl), albumArtMediumUrl: prefixOptionalUrl(p.albumArtMediumUrl), albumArtLargeUrl: prefixOptionalUrl(p.albumArtLargeUrl), audioUrl: prefixUrl(p.audioUrl) };
    }
    case 'photo': {
      const p = post as PhotoPost;
      return {
        ...p,
        images: p.images.map((img) => ({
          ...img,
          url: prefixUrl(img.url),
          thumbnailTinyUrl: prefixOptionalUrl(img.thumbnailTinyUrl),
          thumbnailSmallUrl: prefixOptionalUrl(img.thumbnailSmallUrl),
          thumbnailMediumUrl: prefixOptionalUrl(img.thumbnailMediumUrl),
          thumbnailLargeUrl: prefixOptionalUrl(img.thumbnailLargeUrl),
        })),
      };
    }
    case 'video': {
      const p = post as VideoPost;
      return { ...p, thumbnailUrl: prefixUrl(p.thumbnailUrl), thumbnailTinyUrl: prefixOptionalUrl(p.thumbnailTinyUrl), thumbnailSmallUrl: prefixOptionalUrl(p.thumbnailSmallUrl), thumbnailMediumUrl: prefixOptionalUrl(p.thumbnailMediumUrl), thumbnailLargeUrl: prefixOptionalUrl(p.thumbnailLargeUrl) };
    }
    case 'link': {
      const p = post as LinkPost;
      return { ...p, thumbnailUrl: prefixOptionalUrl(p.thumbnailUrl), thumbnailTinyUrl: prefixOptionalUrl(p.thumbnailTinyUrl), thumbnailSmallUrl: prefixOptionalUrl(p.thumbnailSmallUrl), thumbnailMediumUrl: prefixOptionalUrl(p.thumbnailMediumUrl), thumbnailLargeUrl: prefixOptionalUrl(p.thumbnailLargeUrl) };
    }
    default:
      return post;
  }
}

// During isolated UI builds without a backend, we still compile the app.
// Pages will render with empty collections, and actual data loading requires
// NEXT_PUBLIC_API_URL to be configured at runtime.

// ─── TAGS ────────────────────────────────────────────────────────────────────

/**
 * GET /api/tags
 * Returns all tag names sorted alphabetically.
 */
export async function getTags(): Promise<string[]> {
  if (IS_API_DISABLED) {
    return [];
  }
  const response = await fetch(`${getFetchBase()}/api/tags`, { cache: 'no-store' });
  if (!response.ok) {
    return [];
  }
  return response.json();
}

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
  const response = await fetch(`${getFetchBase()}/api/posts${query}`);
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

  const response = await fetch(`${getFetchBase()}/api/posts/${slug}`);
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

  const response = await fetch(`${getFetchBase()}/api/links`);
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

  const response = await fetch(`${getFetchBase()}/api/links/featured`);
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

function prefixPhotoImage(img: PhotoImage): PhotoImage {
  return {
    ...img,
    url: prefixUrl(img.url),
    thumbnailTinyUrl: prefixOptionalUrl(img.thumbnailTinyUrl),
    thumbnailSmallUrl: prefixOptionalUrl(img.thumbnailSmallUrl),
    thumbnailMediumUrl: prefixOptionalUrl(img.thumbnailMediumUrl),
    thumbnailLargeUrl: prefixOptionalUrl(img.thumbnailLargeUrl),
  };
}

export async function addPhotoImage(postID: string, formData: FormData, token: string): Promise<PhotoImage> {
  const response = await fetch(`${API_URL}/api/posts/${postID}/images`, {
    method: 'POST',
    headers: authHeaders(token),
    body: formData,
  });
  await throwOnError(response, 'addPhotoImage');
  const img = await response.json() as PhotoImage;
  return prefixPhotoImage(img);
}

export async function updatePhotoImage(
  postID: string,
  imageID: number,
  input: { caption?: string; alt?: string },
  token: string,
): Promise<PhotoImage> {
  const response = await fetch(`${API_URL}/api/posts/${postID}/images/${imageID}`, {
    method: 'PUT',
    headers: { ...authHeaders(token), 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  });
  await throwOnError(response, 'updatePhotoImage');
  const img = await response.json() as PhotoImage;
  return prefixPhotoImage(img);
}

export async function deletePhotoImage(postID: string, imageID: number, token: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/posts/${postID}/images/${imageID}`, {
    method: 'DELETE',
    headers: authHeaders(token),
  });
  await throwOnError(response, 'deletePhotoImage');
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

export async function subscribeNewsletter(email: string): Promise<void> {
  const res = await fetch(`${getFetchBase()}/api/subscribe`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ type: 'email', address: email }),
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error((body as { error?: string }).error ?? 'subscription failed');
  }
}

export async function confirmSubscription(token: string): Promise<void> {
  const res = await fetch(`${getFetchBase()}/api/confirm?token=${encodeURIComponent(token)}`);
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error((body as { error?: string }).error ?? 'confirmation failed');
  }
}

export async function unsubscribeNewsletter(token: string): Promise<void> {
  const res = await fetch(`${getFetchBase()}/api/unsubscribe?token=${encodeURIComponent(token)}`);
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error((body as { error?: string }).error ?? 'unsubscribe failed');
  }
}

// ─── BROADCASTER ADMIN ────────────────────────────────────────────────────────

export async function getSubscribers(type: string, token: string): Promise<SubscriberSummary[]> {
  const response = await fetch(`${API_URL}/api/subscribers?type=${encodeURIComponent(type)}`, {
    headers: authHeaders(token),
  });
  await throwOnError(response, 'getSubscribers');
  return response.json() as Promise<SubscriberSummary[]>;
}

export async function deleteSubscriber(address: string, token: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/subscribers/${encodeURIComponent(address)}`, {
    method: 'DELETE',
    headers: authHeaders(token),
  });
  await throwOnError(response, 'deleteSubscriber');
}

export async function getBroadcasts(type: string, token: string): Promise<BroadcastSummary[]> {
  const response = await fetch(`${API_URL}/api/broadcasts?type=${encodeURIComponent(type)}`, {
    headers: authHeaders(token),
  });
  await throwOnError(response, 'getBroadcasts');
  return response.json() as Promise<BroadcastSummary[]>;
}

export async function createBroadcast(data: CreateBroadcastInput, token: string): Promise<BroadcastPreview> {
  const response = await fetch(`${API_URL}/api/broadcasts`, {
    method: 'POST',
    headers: { ...authHeaders(token), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  await throwOnError(response, 'createBroadcast');
  return response.json() as Promise<BroadcastPreview>;
}

export async function updateBroadcast(id: number, data: UpdateBroadcastInput, token: string): Promise<BroadcastPreview> {
  const response = await fetch(`${API_URL}/api/broadcasts/${id}`, {
    method: 'PUT',
    headers: { ...authHeaders(token), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  await throwOnError(response, 'updateBroadcast');
  return response.json() as Promise<BroadcastPreview>;
}

export async function getBroadcastSendDetails(id: number, token: string): Promise<BroadcastSendDetail[]> {
  const response = await fetch(`${API_URL}/api/broadcasts/${id}/sends`, {
    headers: authHeaders(token),
  });
  await throwOnError(response, 'getBroadcastSendDetails');
  return response.json() as Promise<BroadcastSendDetail[]>;
}

export async function dispatchBroadcast(id: number, token: string): Promise<{ buffered_count: number }> {
  const response = await fetch(`${API_URL}/api/broadcasts/${id}/dispatch`, {
    method: 'POST',
    headers: authHeaders(token),
  });
  await throwOnError(response, 'dispatchBroadcast');
  return response.json() as Promise<{ buffered_count: number }>;
}
