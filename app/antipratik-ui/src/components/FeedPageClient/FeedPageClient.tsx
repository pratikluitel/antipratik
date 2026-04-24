'use client';

import { useReducer, useMemo, useState, useEffect, useCallback } from 'react';
import type { Post, PhotoPost, MusicPost, LinkPost, VideoPost } from '../../lib/types';
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
import VideoPlayer from '../VideoPlayer/VideoPlayer';
import { useMusicPlayer } from '../MusicProvider/MusicProvider';
import styles from './FeedPageClient.module.css';

interface Props {
  posts: Post[];
  allTags: string[];
  initialTag?: string;
  initialPhotoId?: string;
  initialTrackId?: string;
  initialVideoId?: string;
}

export default function FeedPageClient({ posts, allTags, initialTag, initialPhotoId, initialTrackId, initialVideoId }: Props) {
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
  const [activeVideo, setActiveVideo] = useState<{ url: string; title: string } | null>(null);
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

  useEffect(() => {
    if (!initialVideoId) return;
    const lp = posts.find((p): p is LinkPost => p.type === 'link' && (p as LinkPost).category === 'video' && p.id === initialVideoId);
    if (lp) Promise.resolve().then(() => setActiveVideo({ url: lp.url, title: lp.title }));
    const vp = posts.find((p): p is VideoPost => p.type === 'video' && p.id === initialVideoId);
    if (vp) Promise.resolve().then(() => setActiveVideo({ url: vp.videoUrl, title: vp.title }));
  }, [initialVideoId, posts]);

  const openVideoPlayer = useCallback((post: VideoPost) => setActiveVideo({ url: post.videoUrl, title: post.title }), []);

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
              onVideoPlay={openVideoPlayer}
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
