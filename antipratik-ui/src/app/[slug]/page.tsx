import { notFound } from 'next/navigation';
import type { Metadata } from 'next';
import type { EssayPost } from '../../lib/types';
import { getPost, getPosts } from '../../lib/api';
import ArticleClient from '../../components/ArticleClient/ArticleClient';

interface Props {
  params: Promise<{ slug: string }>;
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const post = await getPost(slug);
  if (!post) return { title: 'antipratik' };
  return {
    title: `${post.title} — antipratik`,
    description: post.excerpt,
  };
}

export async function generateStaticParams() {
  const posts = await getPosts();
  return posts
    .filter((p) => p.type === 'essay')
    .map((p) => ({ slug: (p as EssayPost).slug }));
}

export default async function ArticlePage({ params }: Props) {
  const { slug } = await params;
  const post = await getPost(slug);
  if (!post) notFound();
  return <ArticleClient post={post} />;
}
