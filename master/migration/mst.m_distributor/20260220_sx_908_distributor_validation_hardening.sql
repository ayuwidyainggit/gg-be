ALTER TABLE mst.m_distributor
    ALTER COLUMN distributor_code SET NOT NULL,
    ALTER COLUMN distributor_name SET NOT NULL,
    ALTER COLUMN region_id SET NOT NULL,
    ALTER COLUMN area_id SET NOT NULL,
    ALTER COLUMN channel_id SET NOT NULL,
    ALTER COLUMN sub_distributor_group_id SET NOT NULL,
    ALTER COLUMN dist_price_grp_id SET NOT NULL,
    ALTER COLUMN address SET NOT NULL,
    ALTER COLUMN latitude SET NOT NULL,
    ALTER COLUMN longitude SET NOT NULL,
    ALTER COLUMN is_active SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_m_distributor_cust_code_active
ON mst.m_distributor (cust_id, distributor_code)
WHERE is_del = false;

ALTER TABLE mst.m_distributor_contact
    ALTER COLUMN email DROP NOT NULL;