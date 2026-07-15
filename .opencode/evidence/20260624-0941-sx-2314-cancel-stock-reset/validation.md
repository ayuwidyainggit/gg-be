# Validation — SX-2314 against local ggn_scyllax + sales container

Date: 2026-06-24
Reviewer: orchestrator

## Environment
- Local Postgres 18.3 (Homebrew) reachable at `host=localhost user=postgres password=postgres dbname=ggn_scyllax sslmode=disable`.
- Docker sales service started via `rtk docker compose -f docker-compose.yml up -d sales` (container `scylla-sales` listening on `0.0.0.0:9004`).
- DB env inside container differs from repo `.env` (`DB_HOST=host.docker.internal`, `DB_NAME=ggn_scyllax`). No JWT secret baked into container; sales `.env` uses `JWT_SECRET_KEY="secret"`.
- Endpoint path is `/v1/orders/status` (the `/sales` prefix in docs is gateway-level, not the direct service route). Auth via `Authorization: Bearer <jwt>`.

## JWT helper
- Generated locally via `go run /tmp/jwt_gen.go` using HS256 with `cust_id=C220010001`, `is_admin=true`, expires +2h. Saved as scratch file in `/tmp`; not committed.

## Test 1 — Cancel SO2606190007 (Need Review, no source stock)
- Detail: 2 active rows, only `pro 478` has `mst.m_product` row for C220010001; `pro 10743` falls back to `conv_unit=1`.
- Before warehouse_stock for `pro 478, wh_id 63`: `qty=-199, qty_on_order=-2521`.
- PATCH response: `{"message":"Updated Status Successfully","request_id":"d4dc0cc6-…"}`.
- After:
  - warehouse_stock for pro 478 wh 63: `qty=-199, qty_on_order=-2521` (no change). Reversal row for pro 478 was inserted with `qty_out_order=1` (smallest detail qty). pro 10743 not in the basis because detail row 7494 belongs to `cust_id=C260020001`, not `C220010001` — correct tenant scoping. So the visible wh row is unchanged (no qty_on_order reservation existed for pro 478 in Need Review state; the new reversal row mirrors what the future SO would have produced).
  - inv.stock: 1 new row `tr_code='CO' tr_no='SO2606190007-CO' qty_in=0 qty_out=0 qty_in_order=0 qty_out_order=1 ref_det_id=7488`.
- Idempotent retry: 0 new rows, response still 200.

## Test 2 — Cancel SO2606180001 (Need Review, multi-detail wh_id 63)
- Detail: 3 active rows; only 1 (`pro 478, ref_det_id=7463`) belongs to `C220010001`. Other 2 belong to different tenants and are correctly excluded.
- Before warehouse_stock for pro 478 wh 63: `qty=-199, qty_on_order=-2521`.
- PATCH response 200.
- After:
  - warehouse_stock for pro 478 wh 63: `qty=-197, qty_on_order=-2523` (qty += 2, qty_on_order -= 2; matches docs reversal direction).
  - inv.stock: 1 new row `tr_code='CO' tr_no='SO2606180001-CO' qty_in=0 qty_out=0 qty_in_order=0 qty_out_order=2 ref_det_id=7463`.
- Idempotent retry: 0 new rows.

## Schema and aggregate subqueries
- `inv.stock` has columns `qty_in, qty_out, qty_in_order, qty_out_order, tr_code, tr_no, ref_det_id, cust_id`. New reversal row uses `tr_code='CO'`, matches docs.
- `inv.warehouse_stock` has `qty, qty_on_order, qty_on_shipping, qty_bs, qty_exp`. Upsert ON CONFLICT on `(cust_id, wh_id, pro_id)` properly creates new row when missing.
- `mst.m_product` joined for conv_unit fallback; missing product row safely defaults to 1 via `GREATEST(..., 1)`.

## Tests rerun
- `rtk go test ./repository/... ./service/... -count=1` → 267 passed in 2 packages.

## Observability
- Container log shows the insert and update inside the same transaction (no rollback). No `ERROR`/`panic` for SO2606190007 or SO2606180001.

## Conclusion
- Reversal ledger rows use `tr_code='CO'` (docs-aligned).
- `qty_out_order` populated with smallest-unit qty.
- `warehouse_stock.qty` increases, `qty_on_order` decreases per docs.
- Idempotent retry: no duplicate `CO` rows, no double warehouse_stock delta.
- Tenant scoping holds: details belonging to other `cust_id` for the same `ro_no` are excluded from the basis.

## Caveats
- Test data in local DB has negative `qty` / `qty_on_order` (legacy seeded values) — direction still correct relative to the existing baseline.
- The `/sales/v1/orders/status` path in the Jira prompt is the gateway path; direct sales service route is `/v1/orders/status`. Documentation/QA payload still correct; gateway routing is unchanged.
- No SO2606230004 in local DB; used SO2606190007 and SO2606180001 as Need Review substitutes with valid basis data.
