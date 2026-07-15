DROP INDEX IF EXISTS idx_m_distributor_parent_cust_id;

ALTER TABLE mst.m_distributor
    DROP CONSTRAINT IF EXISTS fk_m_distributor_parent_cust_id;

ALTER TABLE mst.m_distributor
    DROP COLUMN IF EXISTS parent_cust_id;
