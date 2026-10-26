-- A focus row records work segments, not an independent time budget.
-- Preserve all IDs, timestamps, accumulated time, corrections and request tokens.
ALTER TABLE focus_sessions DROP COLUMN duration_minutes;
