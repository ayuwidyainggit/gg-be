# Verification — Restore Sales Orders Dev to SX-2209 QA

Task ID: `20260617-1845-restore-sales-orders-dev-to-sx-2209-qa`

## Branch status

- Repo: `/Users/ujang/Projects/Geekgarden/scylla-be/sales`
- Current branch: `bugfix/SX-2209-qa`
- Base branch: `qa`
- `qa` SHA: `20e03da17da51179df91754b6ef7a03c73541993`
- HEAD SHA: `20e03da17da51179df91754b6ef7a03c73541993`
- Working tree status after patch: modified, uncommitted files only in scoped endpoint/test/model files.
- No commit made. No push made.

## Source used

Primary source commit:

- `f784cf8 fix(order): show purchase rows with original qty`

Source evidence used:

- `rtk git show --stat --oneline f784cf8 -- service/order_service.go service/order_service_test.go`
- local `dev` file content for `service/order_service.go`
- local `dev` file content for `service/order_service_test.go`

No full `dev` merge.
No cherry-pick commit.
Manual endpoint-scoped adaptation only.

## Why files changed

Minimal source files from `f784cf8`:

- `service/order_service.go`
- `service/order_service_test.go`

Narrow dependency files required for compile/runtime mapping of restored SX-2209 behavior:

- `model/order_detail.go`
  - added `original_qty_po1/2/3` to `OrderDetail` and `OrderDetailRead`
  - required because restored service logic and test access `detail.OriginalQtyPo*`
- `entity/order_detail.go`
  - added `original_qty_po1/2/3` to `OrderDetResponse`
  - required because `DetailV2` response must expose original purchase qty values used by regression test
- `model/order.go`
  - added `order_type` to `OrderList` and `Order`
  - required because existing `DetailV2` tests in this branch construct `model.OrderList{OrderType: ...}`
  - `Order` kept aligned with `OrderList` for same persisted column shape

Not changed:

- `controller/order_controller.go`
- `repository/order_repository.go`
- invoice/report/Open API files
- route path / param shape
- auth middleware wiring
- transaction boundary for create-order writes
- tx-aware repository write pattern

## Behavior restored

`GET /sales/v2/orders/:ro_no`

- purchase rows now built from persisted purchase fields, not copied from sales rows
- purchase row include rule now accepts rows with `original_qty_po* > 0` even when current `qty_po* == 0`
- preserves promo exclusion for `item_type = 2`
- preserves sales/final detail filtering behavior
- preserves `middleware.JWTProtected()` route usage because controller route file unchanged

`POST /sales/v1/orders`

- no code change applied
- compile/tests did not prove dependency from `f784cf8`
- existing create flow preserved: Controller → Service → Repository → DB, service transaction boundary, tx-aware repository writes, tenant inputs `cust_id` / `parent_cust_id`

## Commands run summary

Repo/runtime checks:

```bash
rtk docker compose -f docker-compose.yml ps
rtk git status --short --branch
rtk git branch --list
rtk git remote -v
rtk git fetch origin dev qa
rtk git branch -r --list "origin/bugfix/SX-2209-dev"
rtk git show --stat --oneline f784cf8
rtk git show --function-context f784cf8 -- service/order_service.go service/order_service_test.go
rtk git diff --function-context qa..dev -- controller/order_controller.go service/order_service.go repository/order_repository.go entity/order.go model/order.go model/order_detail.go
```

Branch setup / baseline:

```bash
rtk git checkout qa
rtk git pull --ff-only origin qa
rtk git checkout -b bugfix/SX-2209-qa qa
rtk git merge-base --is-ancestor qa HEAD
rtk go test ./service -run TestDetailV2
rtk go test ./...
```

Implementation validation:

```bash
rtk go test ./service -run TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero
rtk go test ./service -run TestDetailV2
rtk go test ./service
rtk go test ./...
rtk git diff --name-status
rtk git diff --stat
```

## Validation results

Baseline before patch:

- `rtk go test ./service -run TestDetailV2` → pass
- `rtk go test ./...` → pass

After patch:

- `rtk go test ./service -run TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero` → pass
- `rtk go test ./service -run TestDetailV2` → pass
- `rtk go test ./service` → pass
- `rtk go test ./...` → pass

No baseline blocker.

## Final changed files

```text
entity/order_detail.go
model/order.go
model/order_detail.go
service/order_service.go
service/order_service_test.go
```

## Scope guard check

Confirmed:

- branch base remains `qa`
- no full `dev` merge
- route param not renamed; remains `:ro_no`
- no unrelated invoice/report/Open API files touched
- no secrets or `.env` touched
- service create transaction boundary unchanged
- repository tx-aware write pattern unchanged
- tenant behavior path for create unchanged

## Risks / remaining notes

- `POST /sales/v1/orders` was not widened because source commit `f784cf8` only restored SX-2209 read-path behavior and tests/compile did not prove create-path dependency.
- `HEAD` still equals `qa` because changes remain uncommitted. Review should use working tree diff, not commit diff.
- Quality gate review still pending if repo workflow requires final signoff before commit/push/MR.
