PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;

-- users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    current_token TEXT,
    token_expires_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

-- Recreate post_tags with ON DELETE CASCADE
CREATE TABLE post_tags_new (
    post_id TEXT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag     TEXT NOT NULL,
    PRIMARY KEY (post_id, tag)
);
INSERT INTO post_tags_new SELECT * FROM post_tags;
DROP TABLE post_tags;
ALTER TABLE post_tags_new RENAME TO post_tags;
CREATE INDEX IF NOT EXISTS idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX IF NOT EXISTS idx_post_tags_tag ON post_tags(tag);

-- essay_posts
CREATE TABLE essay_posts_new (
    post_id              TEXT PRIMARY KEY REFERENCES posts(id) ON DELETE CASCADE,
    title                TEXT NOT NULL,
    slug                 TEXT NOT NULL UNIQUE,
    excerpt              TEXT NOT NULL,
    body                 TEXT NOT NULL,
    reading_time_minutes INTEGER NOT NULL
);
INSERT INTO essay_posts_new SELECT * FROM essay_posts;
DROP TABLE essay_posts;
ALTER TABLE essay_posts_new RENAME TO essay_posts;

-- short_posts
CREATE TABLE short_posts_new (
    post_id TEXT PRIMARY KEY REFERENCES posts(id) ON DELETE CASCADE,
    body    TEXT NOT NULL
);
INSERT INTO short_posts_new SELECT * FROM short_posts;
DROP TABLE short_posts;
ALTER TABLE short_posts_new RENAME TO short_posts;

-- music_posts
CREATE TABLE music_posts_new (
    post_id   TEXT PRIMARY KEY REFERENCES posts(id) ON DELETE CASCADE,
    title     TEXT NOT NULL,
    album_art TEXT NOT NULL,
    audio_url TEXT NOT NULL,
    duration  INTEGER NOT NULL,
    album     TEXT
);
INSERT INTO music_posts_new SELECT * FROM music_posts;
DROP TABLE music_posts;
ALTER TABLE music_posts_new RENAME TO music_posts;

-- photo_posts
CREATE TABLE photo_posts_new (
    post_id  TEXT PRIMARY KEY REFERENCES posts(id) ON DELETE CASCADE,
    location TEXT
);
INSERT INTO photo_posts_new SELECT * FROM photo_posts;
DROP TABLE photo_posts;
ALTER TABLE photo_posts_new RENAME TO photo_posts;

-- photo_images
CREATE TABLE photo_images_new (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id    TEXT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    alt        TEXT NOT NULL,
    caption    TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0
);
INSERT INTO photo_images_new SELECT * FROM photo_images;
DROP TABLE photo_images;
ALTER TABLE photo_images_new RENAME TO photo_images;
CREATE INDEX IF NOT EXISTS idx_photo_images_post_id ON photo_images(post_id);

-- video_posts
CREATE TABLE video_posts_new (
    post_id       TEXT PRIMARY KEY REFERENCES posts(id) ON DELETE CASCADE,
    title         TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,
    video_url     TEXT NOT NULL,
    duration      INTEGER NOT NULL,
    playlist      TEXT
);
INSERT INTO video_posts_new SELECT * FROM video_posts;
DROP TABLE video_posts;
ALTER TABLE video_posts_new RENAME TO video_posts;

-- link_posts
CREATE TABLE link_posts_new (
    post_id       TEXT PRIMARY KEY REFERENCES posts(id) ON DELETE CASCADE,
    title         TEXT NOT NULL,
    url           TEXT NOT NULL,
    domain        TEXT NOT NULL,
    description   TEXT,
    thumbnail_url TEXT,
    category      TEXT
);
INSERT INTO link_posts_new SELECT * FROM link_posts;
DROP TABLE link_posts;
ALTER TABLE link_posts_new RENAME TO link_posts;

COMMIT;
PRAGMA foreign_keys=ON;
