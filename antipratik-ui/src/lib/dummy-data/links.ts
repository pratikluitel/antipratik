/**
 * Dummy external links for local development.
 * 8 links across music, writing, video, and social categories.
 *
 * Do not import this file directly in components or pages.
 * Access it via src/lib/api.ts → getLinks() / getFeaturedLinks().
 */

import type { ExternalLink } from '../types';

export const links: ExternalLink[] = [
  // ── Music ────────────────────────────────────────────────────────────────
  {
    id: 'link-ext-001',
    title: 'antipratik on SoundCloud',
    url: 'https://soundcloud.com/antipratik',
    domain: 'soundcloud.com',
    description: 'Ambient and electronic tracks. Field recordings from Kathmandu and the Himalayas.',
    featured: true,
    category: 'music',
  },
  {
    id: 'link-ext-002',
    title: 'antipratik on Bandcamp',
    url: 'https://antipratik.bandcamp.com',
    domain: 'bandcamp.com',
    description: 'Full albums and EPs. Pay what you want or nothing — the music is meant to be heard.',
    featured: true,
    category: 'music',
  },
  // ── Writing ──────────────────────────────────────────────────────────────
  {
    id: 'link-ext-003',
    title: 'Essays on Substack',
    url: 'https://antipratik.substack.com',
    domain: 'substack.com',
    description: 'Longer essays on music, code, and living at altitude. Published when ready, not on a schedule.',
    featured: true,
    category: 'writing',
  },
  {
    id: 'link-ext-004',
    title: 'Writing on Medium',
    url: 'https://medium.com/@antipratik',
    domain: 'medium.com',
    description: 'Technical writing on distributed systems, developer tooling, and the intersection of craft in music and code.',
    featured: false,
    category: 'writing',
  },
  // ── Video ────────────────────────────────────────────────────────────────
  {
    id: 'link-ext-005',
    title: 'YouTube — Studio Sessions',
    url: 'https://youtube.com/@antipratik',
    domain: 'youtube.com',
    description: 'Behind-the-scenes studio sessions, gear walkthroughs, and long-form process videos.',
    featured: true,
    category: 'video',
  },
  {
    id: 'link-ext-006',
    title: 'Vimeo — Short Films',
    url: 'https://vimeo.com/antipratik',
    domain: 'vimeo.com',
    description: 'Short films and visual essays shot in Nepal. Higher quality than YouTube for the cinematic work.',
    featured: false,
    category: 'video',
  },
  // ── Social ───────────────────────────────────────────────────────────────
  {
    id: 'link-ext-007',
    title: '@antipratik on X',
    url: 'https://x.com/antipratik',
    domain: 'x.com',
    description: 'Sporadic thoughts on music, code, and Kathmandu. The short-form version of everything else.',
    featured: false,
    category: 'social',
  },
  {
    id: 'link-ext-008',
    title: 'GitHub',
    url: 'https://github.com/antipratik',
    domain: 'github.com',
    description: 'Open source code. Tools, utilities, and the occasional library. Most of it is small and useful.',
    featured: false,
    category: 'social',
  },
];
