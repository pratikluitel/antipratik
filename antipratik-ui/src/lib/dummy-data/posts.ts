/**
 * Dummy post data for local development.
 * Content is Nepal/mountain/music contextual — not Lorem Ipsum.
 * Sorted newest first (descending createdAt).
 *
 * Do not import this file directly in components or pages.
 * Access it via src/lib/api.ts → getPosts() / getPost().
 */

import type { Post, EssayPost, ShortPost, MusicPost, PhotoPost, VideoPost, LinkPost } from '../types';

const essays: EssayPost[] = [
  {
    id: 'essay-001',
    type: 'essay',
    createdAt: '2026-03-15T08:30:00Z',
    tags: ['philosophy', 'music', 'impermanence', 'buddhism'],
    title: 'On Impermanence and Code',
    slug: 'on-impermanence-and-code',
    excerpt:
      'A sound exists only in the moment it is heard. A function exists only as long as the runtime allows. These are not different things.',
    readingTimeMinutes: 7,
    body: `
A sound exists only in the moment it is heard. Unlike a painting, which persists on its canvas across centuries, a piece of music is always already disappearing — each note a small act of vanishing. This is not a flaw in the medium. It is the medium.

I've been thinking about this since spending a week at a gompa above Langtang, watching monks build a sand mandala. Three days of careful work, every grain placed with intention. On the fourth morning, they swept it into a bowl and poured it into the river.

The mandala was not destroyed. It was completed.

## What Altitude Does to Thinking

There's something about being at 4,000 metres that clarifies the mind. Maybe it's the thinner air, the slower pace of everything, the absence of notifications. Maybe it's proximity to something larger than yourself.

I was refactoring a payment service when I left for the trek. When I returned, I saw it differently. The code I'd been precious about — the clever abstractions, the elegant edge-case handling — none of it would survive the next major version of the library it depended on. I'd been building a sand mandala.

This isn't nihilism. The monks weren't nihilists. The point is that impermanence is not a problem to be solved. It's a condition to be worked *with*.

## The Repository Pattern as Buddhist Practice

There's a pattern in software called the repository pattern. You write a thin layer that stands between your business logic and your data source. When the data source changes — when you move from SQLite to Postgres, or from REST to GraphQL — you change the repository, not the logic.

The business logic doesn't grieve the lost SQLite connection. It never knew about SQLite. It only knew about the interface.

I think about this when I write functions that will outlast the infrastructure they currently run on. The function is the mandala. The infrastructure is the sand.

Write good interfaces. Hold the implementation loosely.

## On Music Production in a City That Doesn't Sleep

Kathmandu is loud in a specific way. It's not the uniform drone of a highway. It's layered — a dog starting a chain reaction across three neighbourhoods, a temple bell, a truck navigating a lane too small for it, rain on corrugated iron.

I record this. Not as samples, exactly — more as texture. The city is a collaborator.

When I sit down to mix a track at 2am, the sounds outside become part of the sonic context, even when I'm wearing headphones. The mix that sounds right in that context is the mix I keep.

The present moment of music is all there is. The past note is memory. The future note is anticipation. Only the note you're hearing now exists.

I try to make music that knows this about itself.`,
  },
  {
    id: 'essay-002',
    type: 'essay',
    createdAt: '2026-02-20T10:00:00Z',
    tags: ['code', 'craft', 'tools', 'unix'],
    title: 'The Terminal as a Musical Instrument',
    slug: 'the-terminal-as-a-musical-instrument',
    excerpt:
      'A musician who truly knows their instrument stops thinking about the instrument. The same is true for a developer who has made the terminal their own.',
    readingTimeMinutes: 5,
    body: `
There's a guitarist I know who has played the same guitar for twelve years. He doesn't think about the instrument anymore. He thinks about the music, and the instrument disappears.

I've been using the same terminal configuration for four years. I barely notice it. This is the goal.

## What "Making It Your Own" Actually Means

There's a phase in learning an instrument — and a development environment — where you're constantly aware of the tool. You're thinking about key bindings, about where things live, about the friction between your intention and the outcome.

When this phase ends, something important happens. The tool becomes transparent. You see through it to the problem.

This isn't about memorising shortcuts. It's about building muscle memory through genuine use. The shortcuts you use every day become unconscious. The ones you don't need don't matter.

## The Prompt as Interface Design

Your shell prompt is the most frequently rendered piece of UI in your workflow. I spent more time designing mine than I've spent on any UI component in production code.

It shows: current directory (abbreviated), git branch and status, and the exit code of the last command if it was non-zero. That's it. Everything else is noise.

Less information, more quickly understood, is always better than more information.

## On Portability

I keep my dotfiles in a repository. When I set up a new machine, everything I need is one command away. This is not just convenience — it's a discipline.

Every configuration decision I've made is documented in code. I can read the git log and understand why I made each choice. The setup is reproducible and, more importantly, understandable.

Write tools you can explain.`,
  },
  {
    id: 'essay-003',
    type: 'essay',
    createdAt: '2026-01-10T07:00:00Z',
    tags: ['nepal', 'monsoon', 'writing', 'place'],
    title: 'What the Monsoon Teaches About Latency',
    slug: 'what-the-monsoon-teaches-about-latency',
    excerpt:
      'For four months each year, Kathmandu operates at a different clock speed. There are things you can only learn by waiting.',
    readingTimeMinutes: 6,
    body: `# What the Monsoon Teaches About Latency

For four months each year — roughly June through September — the valley fills with cloud and rain. The light changes. The pace changes. You stop planning more than a day ahead because the weather will change your plans anyway.

I've lived through five monsoons here. Each one has taught me something about working with conditions you can't control.

## Latency Is Not Lag

In distributed systems, latency is the time between a request and its response. We optimise to reduce it. We build caches. We precompute. We get frustrated when it exceeds expectations.

But there's a kind of latency that can't be optimised away, and shouldn't be. The latency between planting a thought and seeing what grows from it. The latency between writing a piece of music and understanding whether it's good.

The monsoon enforces this latency. You can't rush the river.

## On Working Slowly

I write better code in the monsoon. I'm not sure why. Maybe the grey light is easier on the eyes. Maybe the rain is good ambient sound. Maybe the general slowing of the city's metabolism affects mine.

I've stopped fighting it. I've started scheduling the exploratory work — the architectural thinking, the research — for the monsoon months. The heads-down execution work for the dry season.

There's wisdom in matching your pace to the environment's pace.

## The Specific Smell of Rain on Hot Stone

There's a word for it: petrichor. The smell of rain on dry earth. In Kathmandu it smells different from anywhere else I've been, because the stone and soil are different.

I've tried to put this smell into music twice. Both times I failed, but the failures were interesting. They led me somewhere else.

Sometimes the unreachable thing is the most useful thing.`,
  },
];

const shorts: ShortPost[] = [
  {
    id: 'short-001',
    type: 'short',
    createdAt: '2026-03-20T14:22:00Z',
    tags: ['kathmandu', 'fog', 'morning'],
    body: 'just spent an hour watching fog roll over Shivapuri. some mornings the city disappears completely. those are the good mornings.',
  },
  {
    id: 'short-002',
    type: 'short',
    createdAt: '2026-03-12T02:15:00Z',
    tags: ['code', 'night', 'rain'],
    body: 'refactoring at 2am while rain hits the tin roofs. there is no better context for this kind of work. the city is quiet and the problem is the only thing.',
  },
  {
    id: 'short-003',
    type: 'short',
    createdAt: '2026-02-28T09:45:00Z',
    tags: ['boudhanath', 'observation', 'ritual'],
    body: 'watched an old man do his morning kora at Boudha. same direction, same pace, same beads. 108 times. wondered how many of the world\'s problems would dissolve if we each had one ritual we trusted that completely.',
  },
];

const musicPosts: MusicPost[] = [
  {
    id: 'music-001',
    type: 'music',
    createdAt: '2026-03-01T12:00:00Z',
    tags: ['ambient', 'kathmandu', 'field-recording'],
    title: 'Threshold (Nagarkot Dawn)',
    albumArt: '/images/music/threshold.jpg',
    audioUrl: '',
    duration: 144,
    album: 'Altitude Studies',
  },
  {
    id: 'music-002',
    type: 'music',
    createdAt: '2026-01-18T16:00:00Z',
    tags: ['electronic', 'folk-fusion', 'sarangi'],
    title: 'Trisuli',
    albumArt: '/images/music/trisuli.jpg',
    audioUrl: '/audio/trisuli.mp3',
    duration: 144,
    album: 'River Sounds',
  },
];

const photoPosts: PhotoPost[] = [
  {
    id: 'photo-001',
    type: 'photo',
    createdAt: '2026-02-14T06:30:00Z',
    tags: ['langtang', 'mountains', 'winter', 'trek'],
    images: [
      {
        url: '/images/photos/langtang-ridge-dawn.jpg',
        alt: 'Langtang ridge at dawn, thin clouds below the peaks',
        caption: 'The valley floor was still dark when the peaks caught the first light.',
      },
      {
        url: '/images/photos/langtang-village-snow.jpg',
        alt: 'Stone houses in Langtang village with fresh snow',
        caption: 'Langtang village. A week after the snowfall.',
      },
      {
        url: '/images/photos/langtang-prayer-flags.jpg',
        alt: 'Prayer flags strung between two stone pillars against a clear blue sky',
      },
    ],
    location: 'Langtang Valley, 3,500m',
  },
  {
    id: 'photo-002',
    type: 'photo',
    createdAt: '2026-01-28T17:00:00Z',
    tags: ['kathmandu', 'patan', 'durbar-square', 'evening'],
    images: [
      {
        url: '/images/photos/patan-durbar-evening.jpg',
        alt: 'Patan Durbar Square at dusk, warm light on stone temples',
        caption: 'Patan Durbar, January. The square empties after sunset.',
      }
    ],
    location: 'Patan, Lalitpur',
  },
];

const videoPosts: VideoPost[] = [
  {
    id: 'video-001',
    type: 'video',
    createdAt: '2026-02-05T11:00:00Z',
    tags: ['music', 'studio', 'process', 'ambient'],
    title: 'Making "Threshold" — a studio session',
    thumbnailUrl: '/images/videos/threshold-session-thumb.jpg',
    videoUrl: 'https://vimeo.com/placeholder/threshold-session',
    duration: 1140,
    playlist: 'Studio Sessions',
  },
];

const linkPosts: LinkPost[] = [
  {
    id: 'link-001',
    type: 'link',
    createdAt: '2026-03-10T09:00:00Z',
    tags: ['music', 'technology', 'generative'],
    title: 'Why generative music keeps failing to be interesting',
    url: 'https://example.com/generative-music-failure',
    domain: 'example.com',
    description:
      'A long-form piece on why most generative music sounds like waiting room ambience, and what the exceptions do differently.',
    category: 'music',
  },
];

export const posts: Post[] = [
  shorts[0],    // 2026-03-20 — short
  essays[0],    // 2026-03-15 — essay
  musicPosts[0], // 2026-03-01 — music
  linkPosts[0],  // 2026-03-10 — link
  shorts[1],    // 2026-03-12 — short
  essays[1],    // 2026-02-20 — essay
  photoPosts[0], // 2026-02-14 — photo
  videoPosts[0], // 2026-02-05 — video
  shorts[2],    // 2026-02-28 — short
  musicPosts[1], // 2026-01-18 — music
  photoPosts[1], // 2026-01-28 — photo
  essays[2],    // 2026-01-10 — essay
].sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
