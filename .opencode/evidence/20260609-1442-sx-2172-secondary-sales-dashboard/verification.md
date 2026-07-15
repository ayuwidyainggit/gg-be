# Verification Evidence — SX-2172 Secondary Sales Dashboard

Task ID: `20260609-1442-sx-2172-secondary-sales-dashboard`

## Source changes
- `sales/repository/report_repository.go`
- `sales/repository/report_repository_test.go`

## Test validation
Commands run from `sales`:

```bash
rtk go test ./repository -run 'TestSecondarySalesReportGroup'
```
Result: pass — `Go test: 8 passed in 1 packages`

```bash
rtk go test ./service -run 'TestSecondarySalesReportGroupSales'
```
Result: pass — `Go test: 6 passed in 1 packages`

```bash
rtk go test ./...
```
Result: pass — `Go test: 240 passed in 22 packages`

## Docker/runtime status
Command from repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```
Result: compose services running; `scylla-sales` exposed on `0.0.0.0:9004->9004/tcp`, `scylla-system` exposed on `0.0.0.0:9001->9001/tcp`.

Health checks:
- `GET http://localhost:9004/ping` returned HTTP `200`.
- `GET http://localhost:9001/ping` returned HTTP `200`.

## Local DB validation: `ggn_scyllax`
Connection used: local Postgres on `localhost:5432`, DB `ggn_scyllax`, user `postgres`.

Dataset check for `cust_id='C260020001'`, month `4`, year `2026`:
- `report.fact_orders`: `37` rows, net sales `477368000.0000`.
- Rows with `report.dim_products.category_id` empty/0/null for that slice: `37` rows.

Master-category query result sample:
- `id=77`, `code=02`, `name=Jersey`, `net_sales=373392000.0000`
- `id=74`, `code=01`, `name=Mainan`, `net_sales=103976000.0000`

Product query result sample:
- `id=10733`, `code=JY1-002`, `name=Jersey Manchester City FC`, `net_sales=97475000.0000`
- `id=10743`, `code=AF-001`, `name=Action Figure Maradona`, `net_sales=75000000.0000`

Outlet query result sample:
- `id=421`, `code=BMI260011`, `name=Piere Njangka > Toko tosca`, `net_sales=180500000.0000`

## Docker endpoint validation
Login tokens were obtained via local `system` service and used only in memory; credentials/tokens were not written to files.

Principal user scope from local login:
- `cust_id=C26002`
- `parent_cust_id=C26002`
- `distributor_id=0`

Distributor user scope from local login:
- `cust_id=C260020001`
- `parent_cust_id=C26002`
- `distributor_id=102`

Principal endpoint checks against `http://localhost:9004/v1/reports/secondary-sales/group?month=4&year=2026&cust_id=C260020001`:
- `group_by=outlet`: HTTP OK JSON, count `8`, first `{id:421, code:"BMI260011", name:"Piere Njangka > Toko tosca", net_sales:180500000}`.
- `group_by=salesman`: HTTP OK JSON, count `1`, first `{id:421, code:"EMP0025", name:"Piere Njangka", net_sales:477368000}`.
- `group_by=product_category`: HTTP OK JSON, count `2`, first `{id:77, code:"02", name:"Jersey", net_sales:373392000}`.
- `group_by=product`: HTTP OK JSON, count `13`, first `{id:10733, code:"JY1-002", name:"Jersey Manchester City FC", net_sales:97475000}`.

Distributor endpoint checks against `http://localhost:9004/v1/reports/secondary-sales/group?month=4&year=2026`:
- `group_by=outlet`: HTTP OK JSON, count `8`, first `{id:421, code:"BMI260011", name:"Piere Njangka > Toko tosca", net_sales:180500000}`.
- `group_by=salesman`: HTTP OK JSON, count `1`, first `{id:421, code:"EMP0025", name:"Piere Njangka", net_sales:477368000}`.
- `group_by=product_category`: HTTP OK JSON, count `2`, first `{id:77, code:"02", name:"Jersey", net_sales:373392000}`.
- `group_by=product`: HTTP OK JSON, count `13`, first `{id:10733, code:"JY1-002", name:"Jersey Manchester City FC", net_sales:97475000}`.

## Plan compliance
- Product/category now use master-priority joins with fallback to report dims.
- Product name appears in `group_by=product` endpoint result.
- Category appears despite all checked fact-order rows having report dim category missing/0/null.
- Outlet mapping follows user-selected mapping: `id=salesman_id`, `code=outlet_code`, `name=salesman_name > outlet_name`.
- Salesman mapping remains `id=salesman_id`, `code=salesman_code`, `name=salesman_name`.
- Response aliases remain `id`, `code`, `name`, `net_sales`.
- Return subtraction and `dt."year"` filters are covered by tests.

## Notes
- Repo root is not a Git repository in this workspace, so `git status` could not be used for final diff boundary. Changed file set is based on implementation task output and file inspection.
- `rtk` warned that project filters are untrusted; filters were not applied. Tests still passed.
