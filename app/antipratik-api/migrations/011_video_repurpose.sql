-- Migrate existing video_posts (bare external URLs) to link_posts with category='video'.
-- Guards against re-running: INSERT OR IGNORE + the NOT LIKE '/files/%' check ensure
-- newly-uploaded video posts (which store relative /files/... URLs) are never migrated.

INSERT OR IGNORE INTO link_posts
  (post_id, title, url, domain, description, thumbnail_url,
   thumbnail_tiny_url, thumbnail_small_url, thumbnail_medium_url, thumbnail_large_url, category)
SELECT
  vp.post_id, vp.title, vp.video_url, '', NULL,
  vp.thumbnail_url, vp.thumbnail_tiny_url, vp.thumbnail_small_url,
  vp.thumbnail_medium_url, vp.thumbnail_large_url,
  'video'
FROM video_posts vp
JOIN posts p ON p.id = vp.post_id
WHERE p.type = 'video'
  AND vp.video_url NOT LIKE '/files/%';

UPDATE posts SET type = 'link'
WHERE id IN (
  SELECT lp.post_id
  FROM link_posts lp
  JOIN posts p ON p.id = lp.post_id
  WHERE lp.category = 'video' AND p.type = 'video'
);

DELETE FROM video_posts
WHERE post_id IN (SELECT post_id FROM link_posts WHERE category = 'video');

-- Recreate video_posts to: remove NOT NULL from duration, make thumbnail_url nullable,
-- drop playlist column, add description column.
PRAGMA foreign_keys = OFF;

CREATE TABLE video_posts_new (
    post_id               TEXT PRIMARY KEY REFERENCES posts(id) ON DELETE CASCADE,
    title                 TEXT NOT NULL,
    description           TEXT,
    thumbnail_url         TEXT,
    thumbnail_tiny_url    TEXT,
    thumbnail_small_url   TEXT,
    thumbnail_medium_url  TEXT,
    thumbnail_large_url   TEXT,
    video_url             TEXT NOT NULL
);

INSERT INTO video_posts_new
  (post_id, title, description, thumbnail_url, thumbnail_tiny_url,
   thumbnail_small_url, thumbnail_medium_url, thumbnail_large_url, video_url)
SELECT
  post_id, title, NULL, thumbnail_url, thumbnail_tiny_url,
  thumbnail_small_url, thumbnail_medium_url, thumbnail_large_url, video_url
FROM video_posts;

DROP TABLE video_posts;
ALTER TABLE video_posts_new RENAME TO video_posts;

PRAGMA foreign_keys = ON;
