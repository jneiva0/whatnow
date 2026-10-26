ALTER TABLE tasks ADD COLUMN estimate_minutes INTEGER CHECK (estimate_minutes IS NULL OR estimate_minutes BETWEEN 1 AND 1440);
CREATE TABLE steps (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 task_id INTEGER NOT NULL REFERENCES tasks(id),
 description TEXT NOT NULL CHECK (length(trim(description)) > 0),
 estimate_minutes INTEGER CHECK (estimate_minutes IS NULL OR estimate_minutes BETWEEN 1 AND 1440),
 position INTEGER NOT NULL,
 done INTEGER NOT NULL DEFAULT 0 CHECK (done IN (0,1)),
 deleted_at TEXT
);
CREATE INDEX steps_task_order ON steps(task_id, position, id);
CREATE TABLE focus_sessions (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 task_id INTEGER NOT NULL REFERENCES tasks(id),
 step_id INTEGER REFERENCES steps(id),
 duration_minutes INTEGER NOT NULL CHECK (duration_minutes BETWEEN 1 AND 1440),
 state TEXT NOT NULL CHECK (state IN ('running','paused','ended')),
 segment_started_at TEXT,
 accumulated_ms INTEGER NOT NULL DEFAULT 0 CHECK (accumulated_ms >= 0),
 created_at TEXT NOT NULL,
 ended_at TEXT,
 CHECK ((state = 'running') = (segment_started_at IS NOT NULL))
);
CREATE UNIQUE INDEX focus_one_running ON focus_sessions(state) WHERE state = 'running';
CREATE TABLE focus_selection (id INTEGER PRIMARY KEY CHECK (id = 1), session_id INTEGER REFERENCES focus_sessions(id));
CREATE TABLE focus_requests (request_id TEXT PRIMARY KEY, session_id INTEGER NOT NULL REFERENCES focus_sessions(id));
