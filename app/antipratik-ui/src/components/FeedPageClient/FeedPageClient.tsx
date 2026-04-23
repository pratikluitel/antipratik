'use client';

import { useReducer, useMemo, useState, useEffect } from 'react';
import type { Post, PhotoPost, MusicPost } from '../../lib/types';
import {
  filterReducer,
  initialFilterState,
  applyFilters,
  buildFeedClusters,
} from '../../lib/feed';
import FilterBar from '../FilterBar/FilterBar';
import ClusterDivider from '../ClusterDivider/ClusterDivider';
import PostCard from '../PostCard/PostCard';
import Lightbox from '../Lightbox/Lightbox';
import { useMusicPlayer } from '../MusicProvider/MusicProvider';
import styles from './FeedPageClient.module.css';

interface Props {
  posts: Post[];
  allTags: string[];
  initialTag?: string;
  initialPhotoId?: string;
  initialTrackId?: string;
}

export default function FeedPageClient({ posts, allTags, initialTag, initialPhotoId, initialTrackId }: Props) {
  const [state, dispatch] = useReducer(
    filterReducer,
    initialTag
      ? { ...initialFilterState, activeTags: [initialTag] }
      : initialFilterState,
  );

  const filteredPosts = useMemo(
    () => applyFilters(posts, state),
    [posts, state]
  );

  const feedItems = useMemo(
    () => buildFeedClusters(filteredPosts),
    [filteredPosts]
  );

  const [lightboxImages, setLightboxImages] = useState<PhotoPost['images'] | null>(null);
  const [lightboxIndex, setLightboxIndex] = useState(0);
  const { play } = useMusicPlayer();

  useEffect(() => {
    if (!initialPhotoId) return;
    const photoPost = posts.find((p): p is PhotoPost => p.type === 'photo' && p.id === initialPhotoId);
    if (photoPost) {
      const images = photoPost.images;
      Promise.resolve().then(() => {
        setLightboxImages(images);
        setLightboxIndex(0);
      });
    }
  }, [initialPhotoId, posts]);

  useEffect(() => {
    if (!initialTrackId) return;
    const musicPost = posts.find((p): p is MusicPost => p.type === 'music' && p.id === initialTrackId);
    if (musicPost) {
      const track = musicPost;
      Promise.resolve().then(() => play(track));
    }
  }, [initialTrackId, posts, play]);

  return (
    <div className={styles.page}>
      <FilterBar state={state} allTags={allTags} dispatch={dispatch} />
      <div className={styles.feed}>
        {feedItems.map((item, i) => {
          if (item.kind === 'divider') {
            return <ClusterDivider key={`divider-${i}`} from={item.from} to={item.to} />;
          }
          if (item.kind === 'date') {
            return null;
          }
          return (
            <PostCard
              key={item.post.id}
              post={item.post}
              onPhotoOpen={(imgs, idx) => {
                setLightboxImages(imgs);
                setLightboxIndex(idx);
              }}
              onTagClick={(tag) => dispatch({ type: 'TOGGLE_TAG', tag })}
            />
          );
        })}
      </div>
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
