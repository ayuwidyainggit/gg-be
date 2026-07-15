ALTER TABLE pjp.outlet_visit_list ADD is_extra_call bool NULL DEFAULT false;
ALTER TABLE pjp.route_outlet_history ADD is_extra_call bool NULL DEFAULT false;
ALTER TABLE pjp_principles.outlet_visit_list ADD is_extra_call bool NULL DEFAULT false;
ALTER TABLE pjp_principles.destinations_history  ADD is_extra_call bool NULL DEFAULT false;
