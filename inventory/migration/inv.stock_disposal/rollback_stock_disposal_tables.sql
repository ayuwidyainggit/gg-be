-- ============================================
-- Rollback Migration
-- File: inventory/migration/inv.stock_disposal/rollback_stock_disposal_tables.sql
-- ============================================

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS inv.stock_disposal_detail CASCADE;
DROP TABLE IF EXISTS inv.stock_disposal CASCADE;

-- Drop ENUM type
DROP TYPE IF EXISTS inv.media_category_type CASCADE;

-- Drop sequences (if created)
-- DROP SEQUENCE IF EXISTS inv.stock_disposal_detail_sd_detail_id_seq CASCADE;
-- DROP SEQUENCE IF EXISTS inv.stock_disposal_sd_id_seq CASCADE;

