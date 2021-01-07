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
        *
    FROM (
        SELECT
            item_id,
            time,
            price,
            available,
            row_number() OVER (PARTITION BY item_id ORDER BY time DESC) AS rn
    FROM
        item_prices) AS t
    WHERE
        t.rn = 1
)
SELECT
    i.*,
    CTE.time,
    CTE.price,
    CTE.available
FROM
    items i
    INNER JOIN CTE ON i.ID = CTE.item_id;

