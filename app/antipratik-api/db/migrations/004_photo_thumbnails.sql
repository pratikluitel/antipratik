-- Add thumbnail URL columns to photo_images.
-- Nullable for backward compatibility with existing rows.
ALTER TABLE photo_images ADD COLUMN thumbnail_small_url  TEXT;
ALTER TABLE photo_images ADD COLUMN thumbnail_medium_url TEXT;
ALTER TABLE photo_images ADD COLUMN thumbnail_large_url  TEXT;
