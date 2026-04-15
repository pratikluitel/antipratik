import type { ShortPost } from '../../lib/types';
import styles from './ShortPostCard.module.css';

interface Props {
  post: ShortPost;
  onTagClick?: (tag: string) => void;
}

export default function ShortPostCard({ post, onTagClick }: Props) {
  const date = new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    timeZone: 'UTC',
  }).format(new Date(post.createdAt));

  return (
    <article className={styles.card}>
      <div className={styles.inner}>
        <p className={styles.body}>{post.body}</p>
        <div className={styles.footer}>
          <time className={styles.date} dateTime={post.createdAt}>{date}</time>
          {post.tags.length > 0 && (
            <div className={styles.hashtags}>
              {post.tags.map((tag) => (
                <span key={tag} className={styles.hashtag} onClick={() => onTagClick?.(tag)}>#{tag}</span>
              ))}
            </div>
          )}
        </div>
      </div>
    </article>
  );
}
