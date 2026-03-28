'use client';

import { useReducer, useMemo, useState } from 'react';
import type { Post, PhotoPost } from '../../lib/types';
import {
  filterReducer,
  initialFilterState,
  applyFilters,
  buildFeedClusters,
} from '../../lib/feed';
import FilterBar from '../FilterBar/FilterBar';
import ClusterDivider from '../ClusterDivider/ClusterDivider';
import PostCard from '../PostCard/PostCard';
import styles from './FeedPageClient.module.css';

interface Props {
  posts: Post[];
}

export default function FeedPageClient({ posts }: Props) {
  const [state, dispatch] = useReducer(filterReducer, initialFilterState);

  const allTags = useMemo(() => {
    const tagSet = new Set<string>();
    for (const post of posts) {
      for (const tag of post.tags) {
        tagSet.add(tag);
      }
    }
    return Array.from(tagSet).sort();
  }, [posts]);

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

  // Suppress unused variable warnings until lightbox is wired
  void lightboxImages;
  void lightboxIndex;

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
    </div>
  );
}
