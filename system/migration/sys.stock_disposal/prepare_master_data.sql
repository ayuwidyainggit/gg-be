-- ============================================
-- Stock Disposal - Prepare Master Data (Combined)
-- Created: 2025-11-17
-- Description: Combined script untuk prepare semua master data Stock Disposal
-- ============================================

-- ============================================
-- 1. Prepare sys.m_trans_no (SD Number Format)
-- ============================================
-- Insert format SD number untuk semua customer yang ada
-- Format: SD[YY][MM][DD][3-digit sequential]

INSERT INTO sys.m_trans_no (
    cust_id, 
    tr_code, 
    prefix, 
    trans_format, 
    trans_format_type, 
    no_length, 
    trans_desc, 
    midfix, 
    suffix, 
    seq_name, 
    sch
)
SELECT 
    cust_id, 
    'SD',           -- tr_code
    'SD',           -- prefix
    'YY',           -- trans_format (untuk year format, actual implementation pakai YYMMDD di code)
    1,              -- trans_format_type
    3,              -- no_length (untuk 3-digit sequential: 001, 002, dst)
    'Stock Disposal Number',  -- trans_desc
    '',             -- midfix (empty)
    '',             -- suffix (empty)
    '',             -- seq_name (empty, karena sequence di-handle di BeforeCreate hook)
    ''              -- sch (empty)
FROM smc.m_customer
WHERE NOT EXISTS (
    SELECT 1 
    FROM sys.m_trans_no 
    WHERE sys.m_trans_no.cust_id = smc.m_customer.cust_id 
    AND sys.m_trans_no.tr_code = 'SD'
)
ON CONFLICT (cust_id, tr_code) DO NOTHING;

-- ============================================
-- 2. Prepare sys.m_menu (Stock Disposal Menu)
-- ============================================
-- Insert menu Stock Disposal sebagai child dari Inventory menu
-- menu_id: 1912 (next sequential setelah 1911 - Approval Goods Receipt Branch)

INSERT INTO sys.m_menu (
    menu_id,
    menu_name,
    parent_id,
    tr_code,
    level,
    sort_index,
    package_id,
    form_pos,
    form_class,
    icon_index,
    menu_action,
    menu_type,
    icon_web,
    is_header,
    url_web,
    params,
    shortcut,
    breadcrumbs,
    tr_code2,
    "targetType"
)
VALUES (
    '1912',                 -- menu_id (next sequential setelah 1911)
    'Stock Disposal',       -- menu_name
    '19',                   -- parent_id (Inventory menu)
    'SD',                   -- tr_code
    2,                      -- level (child dari Inventory yang level 1)
    12,                     -- sort_index (next setelah max sort_index 11)
    NULL,                   -- package_id
    0,                      -- form_pos
    NULL,                   -- form_class
    0,                      -- icon_index
    2,                      -- menu_action (default)
    0,                      -- menu_type (default)
    NULL,                   -- icon_web
    false,                  -- is_header
    NULL,                   -- url_web
    NULL,                   -- params
    NULL,                   -- shortcut
    NULL,                   -- breadcrumbs
    NULL,                   -- tr_code2
    'iframe-tab'            -- targetType (default, menggunakan quotes karena case-sensitive)
)
ON CONFLICT (menu_id) DO NOTHING;

-- ============================================
-- 3. Prepare sys.m_trans (Transaction Code SD)
-- ============================================
-- Insert transaction code "SD" untuk Stock Disposal

INSERT INTO sys.m_trans (
    tr_code,
    tr_name,
    tr_desc,
    tr_group_id,
    menu_id,
    seq_no,
    tr_caption,
    menu_id_parent
)
VALUES (
    'SD',                           -- tr_code (PRIMARY KEY)
    'Stock Disposal',               -- tr_name
    'Stock Disposal Transaction',   -- tr_desc
    'SD',                           -- tr_group_id
    '1912',                         -- menu_id (menu_id dari step 2)
    1,                              -- seq_no (default)
    'Stock Disposal',               -- tr_caption
    '19'                            -- menu_id_parent (parent menu Inventory)
)
ON CONFLICT (tr_code) DO NOTHING;

-- ============================================
-- Verification Queries
-- ============================================

-- Check m_trans_no
SELECT 
    'm_trans_no' as table_name,
    COUNT(*) as record_count
FROM sys.m_trans_no 
WHERE tr_code = 'SD';

-- Check m_menu
SELECT 
    'm_menu' as table_name,
    menu_id, 
    menu_name, 
    parent_id, 
    tr_code, 
    level, 
    sort_index 
FROM sys.m_menu 
WHERE menu_id = '1912' OR tr_code = 'SD';

-- Check m_trans
SELECT 
    'm_trans' as table_name,
    tr_code, 
    tr_name, 
    tr_desc, 
    tr_group_id, 
    menu_id, 
    seq_no 
FROM sys.m_trans 
WHERE tr_code = 'SD';

