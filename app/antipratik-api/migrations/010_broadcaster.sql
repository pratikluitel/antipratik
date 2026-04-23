-- Migrate existing newsletter_subscribers rows into the new subscribers table.
-- Legacy rows get confirmed=FALSE and a freshly-generated random token.
CREATE TABLE IF NOT EXISTS subscribers (
    id              INTEGER  PRIMARY KEY AUTOINCREMENT,
    type            TEXT     NOT NULL,
    address         TEXT     NOT NULL UNIQUE,
    token           TEXT     NOT NULL UNIQUE,
    confirmed       BOOLEAN  NOT NULL DEFAULT FALSE,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    confirmed_at    DATETIME,
    unsubscribed_at DATETIME
);

INSERT OR IGNORE INTO subscribers (type, address, token, confirmed, created_at)
SELECT
    'email',
    lower(trim(email)),
    lower(hex(randomblob(32))),
    FALSE,
    created_at
FROM newsletter_subscribers;

DROP TABLE IF EXISTS newsletter_subscribers;

-- Broadcasts (email broadcasts and the singleton contact broadcast)
CREATE TABLE IF NOT EXISTS broadcasts (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    type       TEXT     NOT NULL,
    title      TEXT     NOT NULL,
    data       TEXT     NOT NULL DEFAULT '{}',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Rendered HTML bodies for email broadcasts
CREATE TABLE IF NOT EXISTS email_broadcasts (
    broadcast_id INTEGER PRIMARY KEY REFERENCES broadcasts(id) ON DELETE CASCADE,
    email_body   TEXT    NOT NULL
);

-- Per-subscriber send records; one row per (broadcast, subscriber) pair.
-- status: BUFFERED | SUCCESS | FAILED
CREATE TABLE IF NOT EXISTS broadcast_sends (
    id            INTEGER  PRIMARY KEY AUTOINCREMENT,
    broadcast_id  INTEGER  NOT NULL REFERENCES broadcasts(id) ON DELETE CASCADE,
    subscriber_id INTEGER  NOT NULL REFERENCES subscribers(id) ON DELETE CASCADE,
    status        TEXT     NOT NULL DEFAULT 'BUFFERED',
    message       TEXT,
    scheduled_at  DATETIME NOT NULL,
    sent_at       DATETIME,
    UNIQUE(broadcast_id, subscriber_id)
);

CREATE INDEX IF NOT EXISTS idx_broadcast_sends_broadcast ON broadcast_sends(broadcast_id);
CREATE INDEX IF NOT EXISTS idx_broadcast_sends_due
    ON broadcast_sends(broadcast_id, status, scheduled_at);

-- Inbound contact form messages saved for admin reference
CREATE TABLE IF NOT EXISTS contact_messages (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    name       TEXT     NOT NULL,
    email      TEXT     NOT NULL,
    message    TEXT     NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Singleton broadcast used by the contact form
INSERT OR IGNORE INTO broadcasts (id, type, title, data)
VALUES (1, 'contact', 'Contact Form Message', '{}');
