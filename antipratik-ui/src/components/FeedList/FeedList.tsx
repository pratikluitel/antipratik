'use client';

import type { Post, PhotoPost } from '../../lib/types';
import PostCard from '../PostCard/PostCard';
import styles from './FeedList.module.css';

interface Props {
  posts: Post[];
}

export default function FeedList({ posts }: Props) {
  function handlePhotoOpen(_images: PhotoPost['images'], _startIndex: number) {
    // TODO: wire up lightbox
  }

  return (
    <div className={styles.feed}>
      {posts.map((post) => (
        <PostCard
          key={post.id}
          post={post}
          onPhotoOpen={handlePhotoOpen}
        />
      ))}
    </div>
  );
}
