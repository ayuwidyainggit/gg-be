BEGIN;

ALTER TABLE IF EXISTS pjp.outlet_visit_list
    DROP COLUMN IF EXISTS folder,
    DROP COLUMN IF EXISTS is_update_location,
    DROP COLUMN IF EXISTS longitude,
    DROP COLUMN IF EXISTS latitude,
    DROP COLUMN IF EXISTS photo_path;

DROP INDEX IF EXISTS outlet_visit_list_outlet_idx;
DROP INDEX IF EXISTS outlet_visit_list_date_idx;

COMMIT;

