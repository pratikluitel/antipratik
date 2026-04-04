CREATE TABLE IF NOT EXISTS posts (
    id         TEXT PRIMARY KEY,
    type       TEXT NOT NULL CHECK(type IN ('essay','short','music','photo','video','link')),
    created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_posts_type       ON posts(type);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);

CREATE TABLE IF NOT EXISTS post_tags (
    post_id TEXT NOT NULL REFERENCES posts(id),
    tag     TEXT NOT NULL,
    PRIMARY KEY (post_id, tag)
);
CREATE INDEX IF NOT EXISTS idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX IF NOT EXISTS idx_post_tags_tag     ON post_tags(tag);

CREATE TABLE IF NOT EXISTS essay_posts (
    post_id              TEXT PRIMARY KEY REFERENCES posts(id),
    title                TEXT NOT NULL,
    slug                 TEXT NOT NULL UNIQUE,
    excerpt              TEXT NOT NULL,
    body                 TEXT NOT NULL,
    reading_time_minutes INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS short_posts (
    post_id TEXT PRIMARY KEY REFERENCES posts(id),
    body    TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS music_posts (
    post_id   TEXT PRIMARY KEY REFERENCES posts(id),
    title     TEXT NOT NULL,
    album_art TEXT NOT NULL,
    audio_url TEXT NOT NULL,
    duration  INTEGER NOT NULL,
    album     TEXT
);

CREATE TABLE IF NOT EXISTS photo_posts (
    post_id  TEXT PRIMARY KEY REFERENCES posts(id),
    location TEXT
);

CREATE TABLE IF NOT EXISTS photo_images (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id    TEXT NOT NULL REFERENCES posts(id),
    url        TEXT NOT NULL,
    alt        TEXT NOT NULL,
    caption    TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_photo_images_post_id ON photo_images(post_id);

CREATE TABLE IF NOT EXISTS video_posts (
    post_id       TEXT PRIMARY KEY REFERENCES posts(id),
    title         TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,
    video_url     TEXT NOT NULL,
    duration      INTEGER NOT NULL,
    playlist      TEXT
);

CREATE TABLE IF NOT EXISTS link_posts (
    post_id       TEXT PRIMARY KEY REFERENCES posts(id),
    title         TEXT NOT NULL,
    url           TEXT NOT NULL,
    domain        TEXT NOT NULL,
    description   TEXT,
    thumbnail_url TEXT,
    category      TEXT
);

CREATE TABLE IF NOT EXISTS links (
    id          TEXT PRIMARY KEY,
    title       TEXT NOT NULL,
    url         TEXT NOT NULL,
    domain      TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    featured    INTEGER NOT NULL DEFAULT 0 CHECK(featured IN (0,1)),
    category    TEXT NOT NULL CHECK(category IN ('music','writing','video','social'))
);
CREATE INDEX IF NOT EXISTS idx_links_featured ON links(featured);
