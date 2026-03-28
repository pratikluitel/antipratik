import { getPosts } from '../lib/api';
import FeedPageClient from '../components/FeedPageClient/FeedPageClient';

export default async function Home() {
  const posts = await getPosts();
  return <FeedPageClient posts={posts} />;
}
