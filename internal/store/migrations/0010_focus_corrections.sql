-- Corrections are deliberate ledger entries, independent of generic task/plan undo.
CREATE TABLE focus_corrections (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 session_id INTEGER NOT NULL REFERENCES focus_sessions(id),
 before_ms INTEGER NOT NULL CHECK (before_ms >= 0),
 after_ms INTEGER NOT NULL CHECK (after_ms >= 0),
 reason TEXT NOT NULL CHECK (length(trim(reason)) > 0),
 actor TEXT NOT NULL CHECK (actor IN ('user','agent')),
 request_id TEXT NOT NULL UNIQUE,
 created_at TEXT NOT NULL
);
