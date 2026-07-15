-- ============================================
-- Rollback Migration
-- File: inventory/migration/inv.replenishment/rollback_replenishment_status_table.sql
-- Description: Rollback script to drop replenishment_status table
-- ============================================

-- Drop table (CASCADE to remove all dependencies)
DROP TABLE IF EXISTS inv.replenishment_status CASCADE;

