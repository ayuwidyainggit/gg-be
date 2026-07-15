-- Migration: add parent customer relation to distributor
-- Strategy: nullable first to avoid failing on existing legacy rows

ALTER TABLE mst.m_distributor
    ADD COLUMN IF NOT EXISTS parent_cust_id VARCHAR(10);

WITH candidates AS (
    SELECT
        d.distributor_id,
        CASE
            WHEN LENGTH(d.cust_id) > 6 THEN LEFT(d.cust_id, 6)
            ELSE d.cust_id
        END AS parent_candidate
    FROM mst.m_distributor d
)
UPDATE mst.m_distributor d
SET parent_cust_id = c.parent_candidate
FROM candidates c
JOIN smc.m_customer mc ON mc.cust_id = c.parent_candidate
WHERE d.distributor_id = c.distributor_id
  AND d.parent_cust_id IS NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_m_distributor_parent_cust_id'
    ) THEN
        ALTER TABLE mst.m_distributor
            ADD CONSTRAINT fk_m_distributor_parent_cust_id
            FOREIGN KEY (parent_cust_id)
            REFERENCES smc.m_customer(cust_id)
            ON UPDATE CASCADE
            ON DELETE RESTRICT;
    END IF;
END
$$;

CREATE INDEX IF NOT EXISTS idx_m_distributor_parent_cust_id
    ON mst.m_distributor(parent_cust_id);
