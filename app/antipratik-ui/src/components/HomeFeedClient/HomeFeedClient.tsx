'use client';

import { useState, useCallback } from 'react';
import type { Post, PhotoPost, VideoPost } from '../../lib/types';
import PostCard from '../PostCard/PostCard';
import Lightbox from '../Lightbox/Lightbox';
import VideoPlayer from '../VideoPlayer/VideoPlayer';
import styles from './HomeFeedClient.module.css';

interface Props {
  posts: Post[];
}

export default function HomeFeedClient({ posts }: Props) {
  const [lightboxImages, setLightboxImages] = useState<PhotoPost['images'] | null>(null);
  const [lightboxIndex, setLightboxIndex] = useState(0);
  const [activeVideo, setActiveVideo] = useState<{ url: string; title: string } | null>(null);

  const openVideoPlayer = useCallback((post: VideoPost) => setActiveVideo({ url: post.videoUrl, title: post.title }), []);

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
          onVideoPlay={openVideoPlayer}
        />
      ))}
      {lightboxImages !== null && (
        <Lightbox
          images={lightboxImages}
          startIndex={lightboxIndex}
          onClose={() => setLightboxImages(null)}
        />
      )}
      {activeVideo && (
        <VideoPlayer
          videoUrl={activeVideo.url}
          title={activeVideo.title}
          onClose={() => setActiveVideo(null)}
        />
      )}
    </div>
  );
}
