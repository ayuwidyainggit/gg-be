-- ============================================
-- Stock Disposal Dummy Data
-- Created: 2025-11-17
-- Description: Insert dummy data for Stock Disposal feature testing
-- ============================================

-- Insert transaction code "SD" if not exists
-- First check the structure of m_trans
INSERT INTO sys.m_trans (tr_code, tr_name, tr_desc)
SELECT 'SD', 'Stock Disposal', 'Stock Disposal Transaction'
WHERE NOT EXISTS (SELECT 1 FROM sys.m_trans WHERE tr_code = 'SD');

-- ============================================
-- Insert Stock Disposal Header Data
-- ============================================

-- Stock Disposal 1: Bad Stock disposal
INSERT INTO inv.stock_disposal (
    cust_id,
    sd_id,
    tr_code,
    disposal_date,
    sd_number,
    sup_id,
    wh_id,
    stock_type,
    gr_no,
    note,
    sub_total,
    vat_value,
    total,
    created_by,
    created_at,
    updated_at,
    is_del
) VALUES (
    'C24001',                    -- cust_id
    nextval('inv.stock_disposal_sd_id_seq'),  -- sd_id (auto)
    'SD',                        -- tr_code
    CURRENT_DATE,                -- disposal_date
    'SD251117001',               -- sd_number (manual, akan di-override oleh BeforeCreate jika ada)
    75,                          -- sup_id (PT. FUMAKILA)
    2,                           -- wh_id (Gudang BS - Bad Stock)
    'BS',                        -- stock_type
    NULL,                        -- gr_no (optional)
    'Stock disposal untuk barang rusak',  -- note
    1500000.00,                  -- sub_total
    165000.00,                   -- vat_value (11%)
    1665000.00,                  -- total
    27,                          -- created_by (test agus)
    NOW(),                       -- created_at
    NOW(),                       -- updated_at
    false                        -- is_del
) ON CONFLICT (cust_id, sd_id) DO NOTHING;

-- Get the sd_id that was just inserted (or use a fixed value for testing)
-- For dummy data, we'll use a fixed approach
DO $$
DECLARE
    v_sd_id_1 BIGINT;
    v_sd_id_2 BIGINT;
BEGIN
    -- Insert first stock disposal and get the ID
    INSERT INTO inv.stock_disposal (
        cust_id, tr_code, disposal_date, sd_number, sup_id, wh_id, stock_type, gr_no, note,
        sub_total, vat_value, total, created_by, created_at, updated_at, is_del
    ) VALUES (
        'C24001', 'SD', CURRENT_DATE, 'SD251117001', 75, 2, 'BS', NULL, 
        'Stock disposal untuk barang rusak', 1500000.00, 165000.00, 1665000.00, 
        27, NOW(), NOW(), false
    ) RETURNING sd_id INTO v_sd_id_1;

    -- Insert second stock disposal
    INSERT INTO inv.stock_disposal (
        cust_id, tr_code, disposal_date, sd_number, sup_id, wh_id, stock_type, gr_no, note,
        sub_total, vat_value, total, created_by, created_at, updated_at, is_del
    ) VALUES (
        'C24001', 'SD', CURRENT_DATE, 'SD251117002', 76, 3, 'E', NULL,
        'Stock disposal untuk barang expired', 2000000.00, 220000.00, 2220000.00,
        27, NOW(), NOW(), false
    ) RETURNING sd_id INTO v_sd_id_2;

    -- Insert detail for first stock disposal
    INSERT INTO inv.stock_disposal_detail (
        cust_id, sd_id, pro_id, file_name, file_type, media_category, file_base64, file_size,
        unit_id1, unit_id2, unit_id3, qty1, qty2, qty3,
        purch_price1, purch_price2, purch_price3,
        gross_price, vat, vat_value, sub_total,
        created_by, created_at, updated_at, is_del
    ) VALUES (
        'C24001', v_sd_id_1, 485, 'damaged_product_1.jpg', 'jpg', 'image',
        'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==',  -- dummy base64 (1x1 pixel)
        1024,  -- file_size
        'PCS', 'BOX', 'CART',  -- unit_id1, unit_id2, unit_id3
        5.0, 2.0, 1.0,  -- qty1, qty2, qty3
        100000.00, 500000.00, 1000000.00,  -- purch_price1, purch_price2, purch_price3
        1500000.00,  -- gross_price (qty1*price1 + qty2*price2 + qty3*price3)
        11.0,  -- vat (11%)
        165000.00,  -- vat_value
        1665000.00,  -- sub_total (gross_price + vat_value)
        27, NOW(), NOW(), false
    );

    -- Insert second detail for first stock disposal
    INSERT INTO inv.stock_disposal_detail (
        cust_id, sd_id, pro_id, file_name, file_type, media_category, file_base64, file_size,
        unit_id1, unit_id2, unit_id3, qty1, qty2, qty3,
        purch_price1, purch_price2, purch_price3,
        gross_price, vat, vat_value, sub_total,
        created_by, created_at, updated_at, is_del
    ) VALUES (
        'C24001', v_sd_id_1, 480, 'damaged_product_2.png', 'png', 'image',
        'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==',
        2048,
        'PCS', 'BOX', 'CART',
        3.0, 1.0, 0.0,
        120000.00, 600000.00, 0.00,
        960000.00,
        11.0,
        105600.00,
        1065600.00,
        27, NOW(), NOW(), false
    );

    -- Insert detail for second stock disposal (expired)
    INSERT INTO inv.stock_disposal_detail (
        cust_id, sd_id, pro_id, file_name, file_type, media_category, file_base64, file_size,
        unit_id1, unit_id2, unit_id3, qty1, qty2, qty3,
        purch_price1, purch_price2, purch_price3,
        gross_price, vat, vat_value, sub_total,
        created_by, created_at, updated_at, is_del
    ) VALUES (
        'C24001', v_sd_id_2, 477, 'expired_product_1.jpg', 'jpg', 'image',
        'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==',
        1536,
        'PCS', 'BOX', 'CART',
        10.0, 4.0, 2.0,
        150000.00, 750000.00, 1500000.00,
        2000000.00,
        11.0,
        220000.00,
        2220000.00,
        27, NOW(), NOW(), false
    );

END $$;

-- ============================================
-- Verify inserted data
-- ============================================
SELECT 'Stock Disposal Header Records:' as info;
SELECT sd_id, cust_id, sd_number, disposal_date, sup_id, wh_id, stock_type, note, total 
FROM inv.stock_disposal 
WHERE cust_id = 'C24001' 
ORDER BY sd_id DESC 
LIMIT 5;

SELECT 'Stock Disposal Detail Records:' as info;
SELECT sd_detail_id, sd_id, pro_id, file_name, media_category, qty1, qty2, qty3, sub_total
FROM inv.stock_disposal_detail 
WHERE cust_id = 'C24001'
ORDER BY sd_id DESC, sd_detail_id DESC
LIMIT 10;

