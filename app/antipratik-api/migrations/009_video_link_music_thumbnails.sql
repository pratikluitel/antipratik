ALTER TABLE video_posts ADD COLUMN thumbnail_small_url  TEXT;
ALTER TABLE video_posts ADD COLUMN thumbnail_medium_url TEXT;
ALTER TABLE video_posts ADD COLUMN thumbnail_large_url  TEXT;

ALTER TABLE link_posts ADD COLUMN thumbnail_small_url  TEXT;
ALTER TABLE link_posts ADD COLUMN thumbnail_medium_url TEXT;
ALTER TABLE link_posts ADD COLUMN thumbnail_large_url  TEXT;

ALTER TABLE music_posts ADD COLUMN album_art_small  TEXT;
ALTER TABLE music_posts ADD COLUMN album_art_medium TEXT;
ALTER TABLE music_posts ADD COLUMN album_art_large  TEXT;
