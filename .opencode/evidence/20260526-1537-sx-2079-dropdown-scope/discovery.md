# Discovery SX-2079 Dropdown Scope

Task id: `20260526-1537-sx-2079-dropdown-scope`

## Files inspected

- `AGENTS.md`
- `.opencode/docs/index.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/docs/QUALITY.md`
- `docs/Monitoring Activity - BE.md`
- `master/controller/m_region_controller.go`
- `master/controller/m_area_controller.go`
- `master/controller/business_unit_controller.go`
- `master/service/m_region_service.go`
- `master/service/m_area_service.go`
- `master/service/business_unit_service.go`
- `master/repository/m_region_repository.go`
- `master/repository/m_area_repository.go`
- `master/repository/business_unit_repository.go`
- `master/entity/m_region.go`
- `master/entity/m_area.go`
- `master/entity/business_unit.go`
- `master/model/employee.go`
- `master/repository/employee_repository.go`
- `master/migration/mst.m_employee/20260520_add_employee_territory_mapping.sql`
- `master/pkg/middleware/jwt_middleware.go`
- `master/main.go`
- `master/go.mod`
- `master/repository/business_unit_repository_test.go`
- `master/service/business_unit_service_test.go`
- `master/controller/business_unit_controller_test.go`

## Project patterns found

- Target module: `master`, Go 1.20, Fiber service.
- Layer contract: Controller → Service → Repository → DB.
- JWT middleware sets `cust_id`, `parent_cust_id`, `user_id`, `user_name`, `employee_id`, `distributor_id` in Fiber locals.
- Principal path convention: `distributor_id == nil || distributor_id == 0` in `BusinessUnitService.GetBusinessUnit`.
- Region controller currently passes `cust_id` and `parent_cust_id`, but not `employee_id` or `distributor_id`.
- Area controller currently passes `cust_id` and `parent_cust_id`, but not `employee_id` or `distributor_id`.
- Business-unit controller parses multi-value `region_id`, `area_id`, `is_active` manually with `normalizeIntArrayQuery`.
- `mst.m_employee` already has `region_scope`, `area_scope`, `distributor_scope` in model and migration.
- Mapping tables already exist: `mst.m_employee_region_mapping`, `mst.m_employee_area_mapping`, `mst.m_employee_distributor_mapping`; key column is `emp_id`, not `employee_id`.
- Employee detail repository already queries mapping tables and uses `cust_id`, `emp_id`, `is_del = false`.
- Business-unit domain uses `mst.m_distributor` as master table and `model.BusinessUnitDistributor` response rows.

## Reuse candidates

- Reuse `normalizeEmployeeTerritoryScope` behavior as pattern, but add separate dropdown normalization to accept `specific`, `spesific`, `selected` as specific and default null/unknown to all.
- Reuse `normalizeIntArrayQuery` for array query parsing.
- Reuse `sqlx.In` and placeholder args style from `business_unit_repository.go` for safer dynamic queries.
- Reuse `EmployeeRepository.FindOneByEmployeeIdAndCustId` query shape or add narrow scope lookup method to avoid full employee detail load.
- Reuse business-unit query builder unit test style in `master/repository/business_unit_repository_test.go`.

## Constraints

- Do not change non-principal behavior unless requirement explicitly says.
- Preserve `q`, `page`, `limit`, `sort`, `is_active`, `region_id`, `area_id`.
- Principal filter must use current user employee scope from DB, not request payload.
- Defensive scope normalization required because docs mention `Specific`, `SPESIFIC`, `SELECTED`, `All`, `ALL`, `NULL`.
- Use `emp_id` for mapping joins based on migration and existing employee repository.
- Tenant filter must respect parent/current customer rules from architecture docs.
- Existing region/area repositories build SQL by string concatenation; plan should prefer safer args in new helpers/builders when practical.

## Risks

- Existing Area repository has likely alias bug in search branch: `a.area_code` / `a.area_name` while query alias is `ma`.
- Existing business-unit principal query filters via `EXISTS smc.m_customer parent_cust_id`, while issue says `mst.m_distributor.parent_cust_id`. Implementation must verify local data relationship before replacing or should preserve compatible parent filter.
- Existing region/area QueryParser may not parse `region_id[]` or `area_id[]` arrays as robustly as business-unit controller.
- Explicit `region_id`/`area_id` filter behavior can become bypass if scope is not still enforced. Default safe plan: enforce principal scope even when explicit filters are sent.
- Repository interface changes affect stubs/tests and `main.go` constructor wiring if employee scope repository is injected.

## Commands/docs checked

- Local reads/searches only.
- No external docs needed: behavior is project-specific SQL and Go code.
- No GitHub/web/browser needed: no upstream dependency or UI reference involved.

## Research gate result

- Local project discovery: needed and completed.
- Official docs/context7: not needed; no unknown library behavior.
- GitHub: not needed; no upstream behavior dependency.
- Brave/web: not needed; Jira/docs are local.
- Browser/screenshot: not needed; backend API task.
