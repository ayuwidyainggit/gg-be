-- SX-1989 ROLLBACK: Reset parent_pro_id to 0 for rows backfilled by
-- 20260520_backfill_parent_pro_id_sx_1989.sql
--
-- WARNING: This rollback is only safe if no further changes have been made
-- to the affected rows after the backfill. If rows have been updated since,
-- restore from a pre-backfill database snapshot instead.
--
-- Affected rows (confirmed on staging 2026-05-20):
--   pro_id=10776 cust_id=C260020001 pro_code=AF-007 parent_pro_id=10753
--   pro_id=10777 cust_id=C260020001 pro_code=AF-008 parent_pro_id=10754
--   pro_id=10778 cust_id=C260020001 pro_code=AF-009 parent_pro_id=10755
--
-- For a full rollback of all tenants, restore from backup.

UPDATE mst.m_product
   SET parent_pro_id = 0,
       updated_at    = CURRENT_TIMESTAMP
 WHERE pro_id IN (10776, 10777, 10778)
   AND cust_id = 'C260020001';
