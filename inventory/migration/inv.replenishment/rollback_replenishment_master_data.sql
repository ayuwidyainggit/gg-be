-- ============================================
-- Rollback Migration
-- File: inventory/migration/inv.replenishment/rollback_replenishment_master_data.sql
-- Description: Rollback script to drop replenishment master data tables
-- ============================================

-- Drop tables (CASCADE to remove all dependencies)
-- Note: Drop order matters - drop tables that reference these first
DROP TABLE IF EXISTS inv.delivery_type CASCADE;
DROP TABLE IF EXISTS inv.replenishment_type CASCADE;

