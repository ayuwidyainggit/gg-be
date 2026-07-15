BEGIN;

CREATE SCHEMA IF NOT EXISTS pjp;

ALTER TABLE IF EXISTS pjp.outlet_visit_list
    ADD COLUMN IF NOT EXISTS photo_path VARCHAR(500),
    ADD COLUMN IF NOT EXISTS latitude VARCHAR(50),
    ADD COLUMN IF NOT EXISTS longitude VARCHAR(50),
    ADD COLUMN IF NOT EXISTS is_update_location BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS folder VARCHAR(255);

CREATE INDEX IF NOT EXISTS outlet_visit_list_outlet_idx ON pjp.outlet_visit_list (outlet_id);
CREATE INDEX IF NOT EXISTS outlet_visit_list_date_idx ON pjp.outlet_visit_list (date);

COMMIT;

