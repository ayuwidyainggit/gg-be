CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_roh_date_cust_pjp_outlet_extra_route
ON pjp.route_outlet_history (date, cust_id, pjp_id, outlet_id, is_extra_call, route_code);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_ovl_date_pjp_outlet_extra
ON pjp.outlet_visit_list (date, pjp_id, outlet_id, is_extra_call);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_visits_cust_emp_outlet_created_visit
ON mobile.visits (cust_id, emp_code, outlet_code, created_at DESC, visit_id DESC);
