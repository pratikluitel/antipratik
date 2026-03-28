'use client';

import type { Post, PhotoPost } from '../../lib/types';
import EssayCard from '../EssayCard/EssayCard';
import ShortPostCard from '../ShortPostCard/ShortPostCard';
import MusicCard from '../MusicCard/MusicCard';
import PhotoCard from '../PhotoCard/PhotoCard';
import VideoCard from '../VideoCard/VideoCard';
import LinkCard from '../LinkCard/LinkCard';

interface Props {
  post: Post;
  onPhotoOpen: (images: PhotoPost['images'], startIndex: number) => void;
  onTagClick?: (tag: string) => void;
}

export default function PostCard({ post, onPhotoOpen, onTagClick }: Props) {
  switch (post.type) {
    case 'essay':
      return <EssayCard post={post} />;
    case 'short':
      return <ShortPostCard post={post} onTagClick={onTagClick} />;
    case 'music':
      return <MusicCard post={post} />;
    case 'photo':
      return <PhotoCard post={post} onOpen={onPhotoOpen} />;
    case 'video':
      return <VideoCard post={post} />;
    case 'link':
      return <LinkCard post={post} />;
    default:
      throw new Error(`Unhandled post type: ${(post as { type: string }).type}`);
  }
}
