CREATE TABLE IF NOT EXISTS note_tags
(
    note_id TEXT NOT NULL REFERENCES notes (id) ON DELETE CASCADE,
    tag_id  TEXT NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (note_id, tag_id)

);