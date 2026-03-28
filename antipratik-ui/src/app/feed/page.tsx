import type { Metadata } from 'next';
import { getPosts } from '../../lib/api';
import FeedPageClient from '../../components/FeedPageClient/FeedPageClient';

export const metadata: Metadata = {
  title: 'Feed — antipratik',
};

export default async function FeedPage() {
  const posts = await getPosts();
  return <FeedPageClient posts={posts} />;
}
