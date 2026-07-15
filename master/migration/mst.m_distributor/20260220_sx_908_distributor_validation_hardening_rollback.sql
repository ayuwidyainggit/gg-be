DROP INDEX IF EXISTS uq_m_distributor_cust_code_active;

ALTER TABLE mst.m_distributor
    ALTER COLUMN distributor_code DROP NOT NULL,
    ALTER COLUMN distributor_name DROP NOT NULL,
    ALTER COLUMN region_id DROP NOT NULL,
    ALTER COLUMN area_id DROP NOT NULL,
    ALTER COLUMN channel_id DROP NOT NULL,
    ALTER COLUMN sub_distributor_group_id DROP NOT NULL,
    ALTER COLUMN dist_price_grp_id DROP NOT NULL,
    ALTER COLUMN address DROP NOT NULL,
    ALTER COLUMN latitude DROP NOT NULL,
    ALTER COLUMN longitude DROP NOT NULL,
    ALTER COLUMN is_active DROP NOT NULL;

ALTER TABLE mst.m_distributor_contact
    ALTER COLUMN email SET NOT NULL;