CREATE OR REPLACE VIEW current_sessions AS
SELECT
    *
FROM
    sessions
WHERE
    expires_after > now();

CREATE INDEX IF NOT EXISTS sessions_expiresafter_idx ON sessions USING btree (expires_after);

CREATE OR REPLACE VIEW items_with_price AS
WITH CTE AS (
    SELECT
        item_id,
        MAX(time) AS time,
        price,
        available
    FROM
        item_prices
    GROUP BY
        item_id,
        price,
        available
)
SELECT
    i.*,
    CTE.time,
    CTE.price,
    CTE.available
FROM
    items i
    INNER JOIN CTE ON i.ID = CTE.item_id;

