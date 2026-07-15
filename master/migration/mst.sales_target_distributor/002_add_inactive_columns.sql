-- Add user_inactive and inactive_at columns to mst.m_sales_target_distributor_yearly
ALTER TABLE mst.m_sales_target_distributor_yearly
    ADD COLUMN IF NOT EXISTS user_inactive int8 NULL,
    ADD COLUMN IF NOT EXISTS inactive_at timestamptz(6) NULL;
