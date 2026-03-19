-- =========================
-- EVENTS
-- =========================
CREATE TABLE events (
    event_id INTEGER PRIMARY KEY AUTOINCREMENT,
    delivery_id TEXT PRIMARY KEY,
    event_name TEXT NOT NULL,
    target_url TEXT NOT NULL,
    body_json TEXT NOT NULL,
    secret TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now'))
);