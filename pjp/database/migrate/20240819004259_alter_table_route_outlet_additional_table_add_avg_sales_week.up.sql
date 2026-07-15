ALTER TABLE pjp.route_outlet_additional
ADD COLUMN IF NOT EXISTS avg_sales_week NUMERIC(10,2) DEFAULT 0;