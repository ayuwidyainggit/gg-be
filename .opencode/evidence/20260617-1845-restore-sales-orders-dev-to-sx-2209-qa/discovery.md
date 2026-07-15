# Discovery Evidence — Restore Sales Order Endpoints

Task ID: `20260617-1845-restore-sales-orders-dev-to-sx-2209-qa`

## Scope requested

Plan lengkap untuk ambil code endpoint:

- `POST /sales/v1/orders`
- `GET /sales/v2/orders/{order_no}`

Source branch: `dev`.
Target branch: `bugfix/SX-2209-qa`, dibuat dari `qa`.

## Local repo findings

- Workspace root: `/Users/ujang/Projects/Geekgarden/scylla-be`.
- Root folder is not a Git repo.
- Active Git repo for this task: `/Users/ujang/Projects/Geekgarden/scylla-be/sales`.
- Current branch before planning: `qa...origin/qa`.
- Relevant existing branches:
  - local `dev`
  - local `qa`
  - remote `origin/dev`
  - remote `origin/qa`
  - local `bugfix/SX-2209-dev`
  - remote `origin/bugfix/SX-2209-dev`
- Target branch `bugfix/SX-2209-qa` did not exist locally/remotely during discovery.
- Working tree was clean during discovery.

## Repo docs inspected

- `AGENTS.md`
- `.opencode/docs/index.md`
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`

Missing optional docs:

- `.opencode/docs/PROJECT_STACK.md`
- `.opencode/docs/PROJECT_COMMANDS.md`
- `.opencode/docs/FRAMEWORK_PLAYBOOK.md`

## Existing endpoint ownership

Endpoint route source:

- `sales/controller/order_controller.go`
  - `POST /v1/orders` → `OrderController.Create`
  - `GET /v2/orders/:ro_no` → `OrderController.DetailV2`
  - Both routes protected by `middleware.JWTProtected()`.

Service source:

- `sales/service/order_service.go`
  - `OrderService.Store(...)`
  - `OrderService.DetailV2(...)`

Repository source:

- `sales/repository/order_repository.go`
  - order writes use `Store`, `StoreDetail`, `StoreReward`
  - detail reads use `FindByNo`, `FindDetail`, `FindReward`, stock/product helpers
  - tx-aware helper exists: `RepositoryOrderImpl.model(ctx)` extracts transaction context.

Entity/model source likely involved:

- `sales/entity/order.go`
- `sales/model/order.go`
- `sales/model/order_detail.go`

Tests found:

- `sales/service/order_service_test.go`
  - has `mockOrderRepositoryDetailV2`
  - has DetailV2 tests around purchase details and SX-2209 behavior.
- `sales/controller/so_controller_test.go`
  - controller test pattern uses Fiber + `httptest` + local middleware setting `Locals`.

## Source branch evidence

`git diff qa..dev -- controller/order_controller.go service/order_service.go repository/order_repository.go entity/order.go model pkg main.go go.mod go.sum` shows broad dev changes:

- `controller/order_controller.go`
- `entity/order.go`
- `go.mod`
- `main.go`
- `model/invoice.go`
- `model/invoice_detail.go`
- `model/open_api_config.go`
- `model/order.go`
- `model/order_detail.go`
- `model/promotionV2.go`
- `model/report.go`
- `pkg/constant/constant.go`
- `pkg/middleware/open_api_middleware.go`
- `pkg/validation/validation.go`
- `repository/order_repository.go`
- `service/order_service.go`

Dev contains unrelated Open API/config/invoice/report changes. Endpoint-scope restore must avoid blind merge/cherry-pick of whole `dev` unless user explicitly accepts unrelated changes.

## SX-2209 source evidence

Commit found on `dev`/`bugfix/SX-2209-dev`:

- `f784cf8 fix(order): show purchase rows with original qty`

Changed files:

- `service/order_service.go`
- `service/order_service_test.go`

Core behavior from commit:

- Adds `hasPurchaseDisplayQty(detail model.OrderDetailRead) bool`.
- Adds `shouldIncludePurchaseDetailRow(detail model.OrderDetailRead) bool`.
- Changes purchase detail inclusion in `DetailV2` from active purchase/sales qty only to include rows with `original_qty_po*` even when current purchase qty is zero.
- Adds regression test `TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero`.

This is directly relevant to `GET /sales/v2/orders/{order_no}` and issue key `SX-2209`.

## Constraints

- Preserve Controller → Service → Repository → DB.
- Write operations for `POST /sales/v1/orders` must stay inside service-layer transaction.
- Repository writes must keep tx-context extraction.
- Preserve tenant filters: `cust_id`, `parent_cust_id`, `custId` rules.
- Do not copy secrets or `.env`.
- Use `rtk` prefix for shell workflows in this repo.
- Validate inside `sales/`, not root.

## Research source strategy

Used:

- Repo-local docs.
- Local Git branch/diff/log evidence.
- Local source files and tests.

Skipped:

- Official docs/context7: not needed; behavior is repo-local Go/Fiber/GORM code and branch restore.
- GitHub/web: not needed; source branch and target repo are local Git.
- Browser/screenshot: not UI task.

## Main risk

`dev` has many unrelated changes beyond two requested endpoints. Safe plan must use a patch/cherry-pick strategy constrained to endpoint behavior, not merge all `dev` into `qa`.
