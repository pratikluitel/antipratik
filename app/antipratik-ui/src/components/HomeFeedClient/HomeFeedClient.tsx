'use client';

import { useState } from 'react';
import type { Post, PhotoPost } from '../../lib/types';
import PostCard from '../PostCard/PostCard';
import Lightbox from '../Lightbox/Lightbox';
import styles from './HomeFeedClient.module.css';

interface Props {
  posts: Post[];
}

export default function HomeFeedClient({ posts }: Props) {
  const [lightboxImages, setLightboxImages] = useState<PhotoPost['images'] | null>(null);
  const [lightboxIndex, setLightboxIndex] = useState(0);

  return (
    <div className={styles.feed}>
      {posts.map((post) => (
        <PostCard
          key={post.id}
          post={post}
          onPhotoOpen={(imgs, idx) => {
            setLightboxImages(imgs);
            setLightboxIndex(idx);
          }}
        />
      ))}
      {lightboxImages !== null && (
        <Lightbox
          images={lightboxImages}
          startIndex={lightboxIndex}
          onClose={() => setLightboxImages(null)}
        />
      )}
    </div>
  );
}
