CREATE TABLE IF NOT EXISTS newsletter_subscribers (
    email      TEXT NOT NULL PRIMARY KEY COLLATE NOCASE,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
);
