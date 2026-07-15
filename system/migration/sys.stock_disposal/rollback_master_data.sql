-- ============================================
-- Stock Disposal - Rollback Master Data
-- Created: 2025-11-17
-- Description: Rollback script untuk menghapus data preparation Stock Disposal
-- WARNING: Hati-hati dengan DELETE, pastikan tidak ada data yang sudah digunakan
-- ============================================

-- ============================================
-- 1. Delete sys.m_trans (Transaction Code SD)
-- ============================================
-- Hapus transaction code "SD"
-- Note: Pastikan tidak ada data stock_disposal yang sudah menggunakan tr_code='SD'

DELETE FROM sys.m_trans 
WHERE tr_code = 'SD';

-- ============================================
-- 2. Delete sys.m_menu (Stock Disposal Menu)
-- ============================================
-- Hapus menu Stock Disposal
-- Note: Pastikan tidak ada reference dari sys.m_trans yang menggunakan menu_id='1912'

DELETE FROM sys.m_menu 
WHERE menu_id = '1912';

-- ============================================
-- 3. Delete sys.m_trans_no (SD Number Format)
-- ============================================
-- Hapus format SD number untuk semua customer
-- Note: Pastikan tidak ada data stock_disposal yang sudah menggunakan format ini

DELETE FROM sys.m_trans_no 
WHERE tr_code = 'SD';

-- ============================================
-- Verification Queries
-- ============================================

-- Check m_trans_no (should return 0 rows)
SELECT 
    'm_trans_no' as table_name,
    COUNT(*) as record_count
FROM sys.m_trans_no 
WHERE tr_code = 'SD';

-- Check m_menu (should return 0 rows)
SELECT 
    'm_menu' as table_name,
    COUNT(*) as record_count
FROM sys.m_menu 
WHERE menu_id = '1912' OR tr_code = 'SD';

-- Check m_trans (should return 0 rows)
SELECT 
    'm_trans' as table_name,
    COUNT(*) as record_count
FROM sys.m_trans 
WHERE tr_code = 'SD';

