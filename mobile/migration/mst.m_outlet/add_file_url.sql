-- Add file_url field to mst.m_outlet
ALTER TABLE mst.m_outlet 
ADD COLUMN IF NOT EXISTS file_url TEXT;

COMMENT ON COLUMN mst.m_outlet.file_url IS 'URL file gambar outlet dari upload';
