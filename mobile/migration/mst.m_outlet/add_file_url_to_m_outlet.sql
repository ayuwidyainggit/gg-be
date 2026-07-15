-- Add file_url column to mst.m_outlet table
-- This field will store the URL of uploaded outlet photos/images

ALTER TABLE mst.m_outlet 
ADD COLUMN IF NOT EXISTS file_url TEXT NULL;

