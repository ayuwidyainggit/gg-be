-- ============================================
-- Rollback Migration
-- File: inventory/migration/inv.replenishment/rollback_replenishment_order_tables.sql
-- Description: Rollback script to drop replenishment order tables
-- ============================================

-- Drop tables (in reverse order due to foreign keys)
DROP TABLE IF EXISTS inv.replenishment_order_detail CASCADE;
DROP TABLE IF EXISTS inv.replenishment_order CASCADE;

