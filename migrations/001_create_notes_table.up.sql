CREATE TABLE IF NOT EXISTS notes
(
    id      TEXT PRIMARY KEY, -- UUID
    text    TEXT NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);