# Verification — SX-2224 Secondary Sales Group (outlet-murni)

Task ID: `20260612-0916-sx-2224-secondary-sales-group-code`
Date: 2026-06-12
Executor: `@fixer`
Plan: `.opencode/plans/20260612-0916-sx-2224-secondary-sales-group-code.md`
Mode: Maintenance Stability Mode (regression-first, minimal diff)

---

## Summary

Implemented SX-2224 outlet-murni mapping for `group_by=outlet` in
`buildSecondarySalesReportGroupQuery`. The outlet branch previously used
`salesman_id` as the grouping key and `CONCAT_WS(' > ', emp_name, outlet_name)`
as the display name (SX-2172 behaviour). It now uses `outlet_id` as id and
`outlet_code`/`outlet_name` directly from `mst.m_outlet`, with no salesman or
employee joins.

Salesman, product_category, and product branches are unchanged.

---

## TDD Sequence

### Red

Test renamed from `TestSecondarySalesReportGroupOutletUsesSalesmanOutletDisplayMapping`
to `TestSecondarySalesReportGroupOutletUsesOutletDisplayMapping` with SX-2224
assertions. Before source change the test failed:

```
[FAIL] TestSecondarySalesReportGroupOutletUsesOutletDisplayMapping
  report_repository_test.go:670: expected outlet group SQL to contain "o.outlet_id AS id"
```

### Green

Source updated. All targeted and full-suite tests pass.

---

## Changed Files

| File | Change |
|------|--------|
| `sales/repository/report_repository.go` | `buildSecondarySalesReportGroupQuery` outlet branch only |
| `sales/repository/report_repository_test.go` | Renamed test + SX-2224 assertions + reject assertions |

No other files changed. Controller, service, model, entity, routes, migrations,
env, lockfiles, and `pjp-sales` untouched.

---

## SQL Behaviour — Before vs After

### order branch `outlet`

| Field | Before (SX-2172) | After (SX-2224) |
|-------|-----------------|-----------------|
| `id`   | `o.salesman_id` | `o.outlet_id` |
| `code` | `COALESCE(mo.outlet_code, '')` | `COALESCE(mo.outlet_code, '')` *(unchanged)* |
| `name` | `CONCAT_WS(' > ', NULLIF(e.emp_name,''), NULLIF(mo.outlet_name,''))` | `COALESCE(mo.outlet_name, '')` |
| joins  | `m_outlet` + `m_salesman` + `m_employee` | `m_outlet` only |

### return branch `outlet`

| Field | Before (SX-2172) | After (SX-2224) |
|-------|-----------------|-----------------|
| `id`   | `r.salesman_id` | `r.outlet_id` |
| `code` | `COALESCE(mo.outlet_code, '')` | `COALESCE(mo.outlet_code, '')` *(unchanged)* |
| `name` | `CONCAT_WS(' > ', NULLIF(e.emp_name,''), NULLIF(mo.outlet_name,''))` | `COALESCE(mo.outlet_name, '')` |
| joins  | `m_outlet` + `m_salesman` + `m_employee` | `m_outlet` only |

Preserved invariants:
- `cust_id IN ?` on both order and return branches
- Date range `invoice_date >= ? AND invoice_date < ?` from month/year params
- Return subtraction via `* -1 AS net_sales`
- `GROUP BY id, code, name ORDER BY net_sales DESC`
- Tenant row-level join: `mo.cust_id = o.cust_id` (order) / `mo.cust_id = rd.cust_id` (return)

---

## Validation Commands and Results

All commands run from `sales/` service directory with `rtk` prefix.

```
rtk go test ./repository -run 'TestSecondarySalesReportGroup'
```
Result: **8 passed** — includes new `TestSecondarySalesReportGroupOutletUsesOutletDisplayMapping`
plus `TestSecondarySalesReportGroupQueriesUseSourceTablesAndDateRange` (all 4 variants),
`TestSecondarySalesReportGroupProductCategoryUsesMasterCategoryFallback`,
`TestSecondarySalesReportGroupProductUsesMasterProductFallback`.

```
rtk go test ./service -run 'TestSecondarySalesReportGroupSales'
```
Result: **6 passed**

```
rtk go test ./controller -run 'TestSecondaryReportSalesGroup'
```
Result: **4 passed**

```
rtk go test ./...
```
Result: **245 passed in 22 packages** — zero failures, zero unrelated failures.

---

## Runtime Smoke

Run on 2026-06-12 using local Docker Compose sales service and the user-provided bearer token.

Environment:
- Service: `scylla-sales` on `http://127.0.0.1:9004`
- Compose command: `rtk docker compose -f docker-compose.yml up -d redis rabbitmq sales`
- DB target from compose: `host.docker.internal:5432`, DB `ggn_scyllax`
- Direct DB cross-check: `psql -h localhost -p 5432 -U postgres -d ggn_scyllax`

API smoke results for `month=6&year=2026&cust_id=C260020001`:

### `group_by=outlet`
HTTP 200. Top rows:
```json
[
  {"id":1841,"code":"BMI260029","name":"TK Mawar Melati","net_sales":5019250000},
  {"id":1840,"code":"BMI260028","name":"TK Hijau Mekar","net_sales":181505260},
  {"id":1943,"code":"BMI260040","name":"Outlet NOO","net_sales":58800000},
  {"id":1831,"code":"BMI260024","name":"Toko Berkah jaya","net_sales":52335000},
  {"id":1918,"code":"BMI260037","name":"Botol Hijau","net_sales":45500000}
]
```
DB cross-check with equivalent `sls.order` + `sls.return` query returned matching top rows:
```text
1841|BMI260029|TK Mawar Melati|5019250000
1840|BMI260028|TK Hijau Mekar|181505260
1943|BMI260040|Outlet NOO|58800000
1831|BMI260024|Toko Berkah jaya|52335000
1918|BMI260037|Botol Hijau|45500000
```

### `group_by=salesman`
HTTP 200. Top rows:
```json
[
  {"id":479,"code":"BM300","name":"Herman Lee","net_sales":5241596260},
  {"id":466,"code":"2026","name":"Subiwo","net_sales":62000000},
  {"id":478,"code":"BM200","name":"Arifin","net_sales":58800000},
  {"id":421,"code":"EMP0025","name":"Piere Njangka","net_sales":37494000},
  {"id":415,"code":"EMP0021","name":"Jaka","net_sales":27500000}
]
```

### `group_by=product_category`
HTTP 200. Rows:
```json
[
  {"id":74,"code":"01","name":"Mainan","net_sales":5335929000},
  {"id":77,"code":"02","name":"Jersey","net_sales":114561260}
]
```

### `group_by=product`
HTTP 200. Top rows:
```json
[
  {"id":10745,"code":"AF-003","name":"Action Figure Messi","net_sales":5013835000},
  {"id":10751,"code":"AF-005","name":"Action Figure Ronaldo","net_sales":180000000},
  {"id":10743,"code":"AF-001","name":"Action Figure Maradona","net_sales":75250000},
  {"id":10733,"code":"JY1-002","name":"Jersey Manchester City FC","net_sales":50000000},
  {"id":8436,"code":"LPI-002","name":"Jersey Medan Chief","net_sales":39200000}
]
```

DB master cross-check for product/category:
```text
10733|JY1-002|Jersey Manchester City FC|C260020001|77
77|02|Jersey
```

Runtime smoke conclusion: PASS. All four group variants return `code`, product/product_category names are populated from master data, and outlet now uses outlet id/code/name as required by SX-2224.

---

## Source Strategy

Repo-local code and tests only. External docs and web skipped — no external library
behaviour material to this SQL mapping change.

---

## Invariant Conformance Check

| Invariant | Status |
|-----------|--------|
| outlet uses `outlet_id` / `outlet_code` / `outlet_name` | PASS |
| salesman branch unchanged | PASS |
| product_category branch unchanged | PASS |
| product branch unchanged | PASS |
| `cust_id IN ?` preserved on both CTEs | PASS |
| date range from month/year preserved | PASS |
| return subtraction (`* -1`) preserved | PASS |
| JSON aliases `id/code/name/net_sales` unchanged | PASS |
| `LEFT JOIN` on m_outlet kept (not inner join) | PASS |
| no env / migration / package / lockfile changes | PASS |
| pjp-sales untouched | PASS |
| Controller → Service → Repository → DB contract intact | PASS |

---

## Quality Gate

Routed to `@quality-gate` for final signoff per plan T5.
