# Execution Evidence SX-2131 / SX-2184 Taking Order

Task ID: `20260608-1347-sx-2131-2184-taking-order-audit`
Date: 2026-06-08 14:13 Asia/Jakarta
Mode: Maintenance Stability Mode

## Plan source

- `.opencode/plans/20260608-1347-sx-2131-2184-taking-order-audit.md`

## Branch / status

Command from `sales`:

```bash
git status --short && git branch --show-current && git rev-parse --short HEAD
```

Result:

```text
dev
53d68b3
```

Interpretation:

- Working tree source was clean; `git status --short` produced no file rows.
- Current service branch: `dev`.
- Current commit: `53d68b3`.

## Automated validation

From `sales`:

```bash
rtk go test ./service -run 'Test.*SX2184|Test.*OrderType|Test.*TakingOrder'
rtk go test ./controller -run 'Test.*SX2184|Test.*Create.*OrderType'
```

Result:

```text
Go test: 5 passed in 1 packages
Go test: 5 passed in 1 packages
```

Full suite:

```bash
rtk go test ./...
```

Result:

```text
Go test: 219 passed in 22 packages
```

## Runtime smoke setup

Command from repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

Initial result summary:

- `master`, `redis`, and `system` were already running.
- `sales` and `rabbitmq` were not running at first.

Started runtime dependencies for smoke:

```bash
rtk docker compose -f docker-compose.yml up -d rabbitmq sales
```

Result summary:

- `scylla-rabbitmq` started and became healthy.
- `scylla-sales` started.

Ping check:

```text
HTTP 200 It works
```

Migration apply check:

```bash
PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d ggn_scyllax -v ON_ERROR_STOP=1 -f sales/migration/sls.order/add_order_type_and_original_qty_po_fields.sql
```

Result summary:

- Migration completed successfully.
- All target columns already existed and were skipped by `IF NOT EXISTS` notices.

## Fixture used

Local DB fixture:

- `cust_id = C220010001`
- `parent_cust_id = C22001`
- `salesman_id = 210`
- `outlet_id = 1390`
- `wh_id = 241`
- `pro_id = 472`
- pre-smoke `inv.warehouse_stock.qty = 0`
- pre-smoke `inv.warehouse_stock.qty_on_order = 108`

Pre-smoke stock query result:

```text
C220010001	241	472	0	108
```

## API smoke

Request:

```text
POST http://localhost:9004/v1/orders
```

Security note:

- A local JWT was generated in-memory from local `.env` for this smoke only.
- Token value was not printed, written to repo, copied into evidence, or reused in final output.

Payload summary:

- `order_type = "O"`
- `ro_date = 2026-06-08`
- `wh_id = 241`
- `pro_id = 472`
- `qty1 = 10`, `qty2 = 0`, `qty3 = 0`
- warehouse stock fixture was `0`
- marker note: `sx2184-smoke-1780902798`

Result:

```text
HTTP_STATUS=201
RESPONSE={"message":"Created Successfully","data":{"ro_no":"SO2606080002"},"request_id":"0f53e449-dbfb-4121-a9d2-7ce0f4a131cf"}
```

## DB verification after API smoke

Header/order query result:

```text
SO2606080002	O	O	f	t	2	1	sx2184-smoke-1780902798
```

Column meaning:

- `ro_no = SO2606080002`
- `order_type = O`
- `opr_type = O`
- `validate_stok = false`
- `validate_stok_message IS NULL = true`
- `data_status = 2`
- `data_source = 1`

Detail query result:

```text
10	0	0	10	0	0	10	t	t	t	0	0
```

Column meaning:

- `original_qty_po1/2/3 = 10/0/0`
- `qty_po1/2/3 = 10/0/0`
- `qty_po = 10`
- `qty1/2/3 IS NULL = true/true/true`
- `qty = 0`
- `qty_final = 0`

Inventory stock mutation query:

```text
SELECT count(*) FROM inv.stock WHERE cust_id='C220010001' AND tr_no='SO2606080002';
```

Result:

```text
0
```

Warehouse stock query after create:

```text
0	108
```

Interpretation:

- `inv.warehouse_stock.qty` stayed `0`.
- `inv.warehouse_stock.qty_on_order` stayed `108`.
- No `inv.stock` row was created for this Taking Order create.

## Runtime cleanup

Stopped services started for this smoke:

```bash
rtk docker compose -f docker-compose.yml stop sales rabbitmq
```

Result summary:

- `scylla-sales` stopped.
- `scylla-rabbitmq` stopped.

## Diff boundary / changed files

No source files were changed in this execution pass.

Only evidence artifact added/updated under allowed plan evidence path:

- `.opencode/evidence/20260608-1347-sx-2131-2184-taking-order-audit/execution.md`
- `.opencode/evidence/20260608-1347-sx-2131-2184-taking-order-audit/index.json`

## Residual risks

- Runtime smoke was performed for `order_type = "O"` only.
- `SO/C/nil/empty` runtime smoke was not executed; regression coverage for those paths remains test-level by design because the plan explicitly avoids changing existing non-`O` behavior and avoids expanding stock formula scope.
- The smoke created local DB order `SO2606080002` in `ggn_scyllax`; this is local smoke data, not a remote/staging/dev DB change.
