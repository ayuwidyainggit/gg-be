ALTER TABLE pjp.outlet_visit_list
ADD COLUMN distance_meter INTEGER,
ADD COLUMN allowed_radius INTEGER DEFAULT 100,
ADD COLUMN location_status SMALLINT CHECK (location_status IN (0,1));
