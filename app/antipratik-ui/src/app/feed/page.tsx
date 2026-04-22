import type { Metadata } from 'next';
import { getPosts, getTags } from '../../lib/api';
import FeedPageClient from '../../components/FeedPageClient/FeedPageClient';

export const dynamic = 'force-dynamic';

export const metadata: Metadata = {
  title: 'Feed — antipratik',
};

interface Props {
  searchParams: Promise<{ tag?: string; photo?: string; track?: string }>;
}

export default async function FeedPage({ searchParams }: Props) {
  const { tag, photo, track } = await searchParams;
  const [posts, allTags] = await Promise.all([getPosts(), getTags()]);
  return (
    <FeedPageClient
      posts={posts}
      allTags={allTags}
      initialTag={tag}
      initialPhotoId={photo}
      initialTrackId={track}
    />
  );
}
