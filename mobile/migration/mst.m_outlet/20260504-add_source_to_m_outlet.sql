-- Add source column to mst.m_outlet table
-- This field indicates the origin of the outlet creation: 1 = Mobile, 0 = Web

ALTER TABLE mst.m_outlet 
ADD COLUMN IF NOT EXISTS source INTEGER;
