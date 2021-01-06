CREATE OR REPLACE VIEW current_sessions AS
SELECT
    *
FROM
    sessions
WHERE
    expires_after > now();

CREATE INDEX IF NOT EXISTS sessions_expiresafter_idx ON sessions USING btree (expires_after);

