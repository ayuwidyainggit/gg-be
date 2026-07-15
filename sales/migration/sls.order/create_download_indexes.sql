-- Migration: Create indexes for Sales Order Download endpoint optimization
-- Purpose: Improve query performance for /v1/download endpoint
-- Date: 2025-12-02
-- 
-- These indexes are designed to optimize:
-- 1. Date range filtering (ro_date)
-- 2. Customer filtering (cust_id)
-- 3. Salesman filtering (salesman_id)
-- 4. Product filtering (pro_id)
-- 5. Join operations between sls.order and sls.order_detail

-- Index for ro_date filtering (most common filter)
-- Used in WHERE clauses: sls.order.ro_date >= ? AND sls.order.ro_date <= ?
CREATE INDEX IF NOT EXISTS idx_order_ro_date 
ON sls.order(ro_date);

-- Composite index for order table filtering
-- Optimizes queries with cust_id + ro_date + salesman_id filters
-- Used in WHERE clauses: cust_id = ? AND ro_date BETWEEN ? AND ? AND salesman_id = ?
CREATE INDEX IF NOT EXISTS idx_order_cust_date_salesman 
ON sls.order(cust_id, ro_date, salesman_id);

-- Index for order_detail table filtering
-- Optimizes queries filtering by cust_id and item_type (always item_type = 1 for download)
-- Used in WHERE clauses: cust_id = ? AND item_type = 1
CREATE INDEX IF NOT EXISTS idx_order_detail_cust_item 
ON sls.order_detail(cust_id, item_type);

-- Index for join optimization between order and order_detail
-- Optimizes JOIN: sls.order.ro_no = sls.order_detail.ro_no AND sls.order.cust_id = sls.order_detail.cust_id
CREATE INDEX IF NOT EXISTS idx_order_ro_cust 
ON sls.order(ro_no, cust_id);

CREATE INDEX IF NOT EXISTS idx_order_detail_ro_cust 
ON sls.order_detail(ro_no, cust_id);

-- Partial index for product filtering
-- Only indexes rows where item_type = 1 (normal items, not promo)
-- Used in WHERE clauses: pro_id IN (?) AND item_type = 1
CREATE INDEX IF NOT EXISTS idx_order_detail_pro_id 
ON sls.order_detail(pro_id) 
WHERE item_type = 1;

-- Comments for documentation
COMMENT ON INDEX idx_order_ro_date IS 'Index for date range filtering on order table';
COMMENT ON INDEX idx_order_cust_date_salesman IS 'Composite index for customer, date, and salesman filtering';
COMMENT ON INDEX idx_order_detail_cust_item IS 'Index for customer and item type filtering on order_detail';
COMMENT ON INDEX idx_order_ro_cust IS 'Index for join optimization on order table (ro_no + cust_id)';
COMMENT ON INDEX idx_order_detail_ro_cust IS 'Index for join optimization on order_detail table (ro_no + cust_id)';
COMMENT ON INDEX idx_order_detail_pro_id IS 'Partial index for product filtering (only item_type = 1)';
