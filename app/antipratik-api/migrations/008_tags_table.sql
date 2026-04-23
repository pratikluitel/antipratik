-- Step 1: Create normalized tags table
CREATE TABLE IF NOT EXISTS tags (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

-- Step 2: Populate tags from existing post_tags strings
INSERT OR IGNORE INTO tags (name) SELECT DISTINCT tag FROM post_tags;

-- Step 3: Save existing associations into a temp table before dropping post_tags
CREATE TEMP TABLE post_tags_backup AS
SELECT post_id, tag FROM post_tags;

-- Step 4: Drop old post_tags table entirely
DROP TABLE post_tags;

-- Step 5: Recreate post_tags with normalized tag_id FK
CREATE TABLE post_tags (
    post_id TEXT    NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag_id  INTEGER NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);
CREATE INDEX idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX idx_post_tags_tag_id  ON post_tags(tag_id);

-- Step 6: Re-insert associations using the new tag IDs
INSERT INTO post_tags (post_id, tag_id)
SELECT b.post_id, t.id
FROM post_tags_backup b
JOIN tags t ON t.name = b.tag;

DROP TABLE post_tags_backup;

-- Step 7: Trigger to clean up orphan tags when a post is deleted
CREATE TRIGGER IF NOT EXISTS delete_orphan_tags
AFTER DELETE ON post_tags
BEGIN
    DELETE FROM tags
    WHERE id = OLD.tag_id
      AND NOT EXISTS (SELECT 1 FROM post_tags WHERE tag_id = OLD.tag_id);
END;
