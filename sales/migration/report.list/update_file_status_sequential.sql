-- Migration: Update file_status values to sequential (1,2,3,4)
-- Created: 2024-12-20
-- Purpose: Refactor file_status from non-sequential (1,5,6,10) to sequential (1,2,3,4)
-- Mapping: 5→2 (Processing), 6→3 (Failed), 10→4 (Expired), 1→1 (Ready, no change)

UPDATE report.list 
SET file_status = CASE 
	WHEN file_status = 5 THEN 2  -- Processing
	WHEN file_status = 6 THEN 3  -- Failed
	WHEN file_status = 10 THEN 4 -- Expired
	ELSE file_status              -- Ready (1) stays the same
END
WHERE file_status IN (5, 6, 10);

COMMENT ON COLUMN report.list.file_status IS 'File status: 1=Ready, 2=Processing, 3=Failed, 4=Expired';
