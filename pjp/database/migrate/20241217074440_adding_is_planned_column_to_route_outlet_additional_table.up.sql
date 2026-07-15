ALTER TABLE pjp.route_outlet_additional
ADD COLUMN IF NOT EXISTS is_planned BOOLEAN DEFAULT false;
