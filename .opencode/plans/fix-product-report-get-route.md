# Plan — fix-product-report-get-route

Task ID: `fix-product-report-get-route`
Mode: maintenance-stability (route shadowing regression)
Target service: `master`
Target endpoint: `GET /v1/products/report` (with existing `POST /v1/products/report` preserved)
Plan quality gate: `PASS_FOR_SLICE`
plan_status: PASS_FOR_SLICE
preflight_disposition: target-app

## Goal

`GET /v1/products/report?cust_id[]=C26002&cust_id[]=C260020001&page=1&limit=20` must reach the existing `Report` handler and return the master response envelope (`data` + `paging` + `request_id`). Currently the route shadowing causes Fiber to dispatch `GET /v1/products/report` to `GET /:pro_id`, where `ParamsParser` chokes on `report` as int64 and the request fails before reaching the report logic. Fix registers a literal `GET /report` route inside the existing JWT-protected group, ordered before the parameter route, and reuses the unchanged `Report` handler.

Slice is one route registration + one router-level regression test. No service/repository/DB change.

## Non-goals

- No handler, DTO, service, or repository change.
- No JWT middleware change, no auth weakening, no new middleware.
- No DDL/migration, no env, no other service touched.
- No new dependency.
- No expansion to other GET endpoints (`/principals`, `/categories`, `/brands` are already literal).
- No documentation beyond evidence already in this plan.

## Scope

Allowed file groups:

- `master/controller/product_controller.go` — add one `productsRouteV1.Get("/report", controller.Report)` line, positioned before `productsRouteV1.Get("/"+qParamId, controller.Detail)`.
- `master/controller/product_report_controller_test.go` — extend with router-level regression tests that mount `ctrl.Route(app)` and assert (a) `GET /v1/products/report` reaches `Report` (no int64 parse), (b) `GET /v1/products/:pro_id` with numeric id still reaches `Detail` or its equivalent dispatch path.

Out of scope:

- `master/service/`, `master/repository/`, `master/entity/`, `master/pkg/`, `master/main.go`.
- `go.mod`, `go.sum`, `docker-compose.yml`, `.env`, `master/Makefile`.
- `master/controller/*` other than the two files above.

## Requirements

1. Add `productsRouteV1.Get("/report", controller.Report)` inside the existing `Route(app)` method, registered before the `:pro_id` parameter route.
2. `POST /v1/products/report` registration and behavior remain unchanged.
3. Route group middleware remains `middleware.JWTProtected()`; no new middleware added.
4. `Report` handler unchanged; filter parsing, sort/page/sort_order validation, response envelope, paging shape preserved.
5. Add router-level regression test in `master/controller/product_report_controller_test.go` that:
   a. Calls `controller.Route(app)` on a fresh `fiber.New()` instance;
   b. Sends `GET /v1/products/report?cust_id[]=C26002&page=1&limit=20` and asserts HTTP 200 plus that `reportCalled` on the stub is true and `reportDetailCalled` is false;
   c. Sends `GET /v1/products/12345` and asserts dispatch reaches detail logic (or that it does not call `ReportList`); the existing `Detail` validator path is exercised as before.
6. Tests must not require a live DB or real JWT; they stub `ProductService` like the existing test file.
7. `cd master && rtk go test ./controller -run TestProductReport -v` and `cd master && rtk go test ./...` exit 0.
8. `cd master && rtk go build ./...` exits 0.
9. Diff is limited to the two files in Scope plus this plan and evidence notes.

## Acceptance Criteria

1. `GET /v1/products/report?cust_id[]=C26002&page=1&limit=20` reaches `Report` handler (service called, 200 envelope, `data` + `paging` returned). Proven by router-level test using `controller.Route(app)`.
2. `GET /v1/products/report` no longer falls into `GET /v1/products/:pro_id`. Test asserts `reportDetailCalled == false` and `reportCalled == true`.
3. `GET /v1/products/12345` (numeric) still routes to detail logic; existing `Detail` flow (validator + service) is unchanged. Test asserts dispatch hits detail service and not `ReportList`.
4. `POST /v1/products/report` continues to pass all four existing tests (`TestProductReport_MissingCustID_Returns400`, `TestProductReport_BlankCustID_Returns400`, `TestProductReport_InvalidSortBy_Returns400`, `TestProductReport_InvalidSortOrder_Returns400`, `TestProductReport_ValidRequest_CallsService`).
5. JWT middleware is still wrapping `/v1/products/*` group; no new route escapes the group. Confirmed by reading `Route(app)` after the change.
6. `cd master && rtk go test ./controller -run TestProductReport -v` exit 0.
7. `cd master && rtk go test ./...` exit 0.
8. `cd master && rtk go build ./...` exit 0.
9. Diff is restricted to the two files named in Scope.

## Existing Patterns/Reuse

- `master/controller/product_controller.go:43-49` — route group; new line follows the same pattern as `productsRouteV1.Post("/report", controller.Report)`.
- `master/controller/product_controller.go:517-595` — handler. Reused as-is.
- `master/controller/product_report_controller_test.go:13-32` — stub and app-mount pattern. Extend with router-mount variant.
- `master/controller/supplier_controller_test.go:43-80` — precedent for setting `c.Locals` before invoking handler; reuse pattern for detail regression.
- `master/pkg/responsebuild/response.go:26` — `BuildResponse` envelope.
- `master/pkg/middleware/jwt_middleware.go:15-32` — `JWTProtected()`. Production code unchanged; test app must replicate `c.Locals` (not the middleware) to bypass JWT in unit tests, matching existing precedent.

## Source Anatomy

| Subsystem | Files | Note |
|---|---|---|
| Route registration | `master/controller/product_controller.go:34-55` | Add literal GET `/report` line; keep JWT group intact. |
| Report handler | `master/controller/product_controller.go:517-595` | Unchanged; serves both POST and GET. |
| Detail handler | `master/controller/product_controller.go:57-99` | Unchanged; still owns `:pro_id` path. |
| Stub test infra | `master/controller/product_report_controller_test.go:13-23` | Reused; add flag for `Detail` call count. |
| JWT middleware | `master/pkg/middleware/jwt_middleware.go:15-32` | Wrapping group preserved. |
| Service interface | `master/service/product_service.go:27-30` | `ReportList` reused by handler; no service edit. |
| Response/paging | `master/pkg/responsebuild/response.go:26`, `master/entity/api.go:16-22` | Unchanged. |

Route registration is controller-only and uses a static segment beside an existing `:pro_id` parameter segment. `Report` requires only `requestid` for normal report success/error response construction; request filter comes solely from query parameters. `Detail` requires `requestid`, `cust_id`, `parent_cust_id`, and `distributor_id`, then calls `ProductService.Detail`. Therefore test setup must provide the same locals expected after middleware but must continue through the actual registered JWT group. `jwtError` permits a missing bearer token when the request provides `Cust_id`; this is test-only transport setup, not production auth behavior. No datasource or query executes in the controller test because `ProductService` is a stub.

## Reference Map

- Routing fix: repo-backed (`master/controller/product_controller.go:43-49`).
- Handler reuse: repo-backed (`master/controller/product_controller.go:517-595`).
- Detail parser failure basis: repo-backed (`master/controller/product_controller.go:57-76`, `master/entity/product.go:1049-1054`).
- Test pattern: repo-backed (`master/controller/product_report_controller_test.go:13-32`, `master/controller/supplier_controller_test.go:43-80`).
- Auth preservation: docs-backed + repo-backed (`.opencode/docs/FRAMEWORK_PLAYBOOK.md:33` + `master/pkg/middleware/jwt_middleware.go:15-32,35-52`).
- Commands: docs-backed (`.opencode/docs/PROJECT_COMMANDS.md:18-22,40-44`).

## Constraints

- Maintain `Controller → Service → Repository → DB`. No direct repo call from controller.
- Service-layer writes remain transactional (no change here; slice is read-only path).
- Reuse `net/http/httptest`; do not introduce testify or new assertion library.
- Parameterized query values only. No new SQL — handler unchanged.
- Per-service Go test/build commands with `rtk` prefix.
- Do not weaken JWT, log bearer tokens, or expose protected routes outside the group.

## Risks

1. Router-level test must install required `c.Locals` (`requestid`, plus `cust_id`/`parent_cust_id`/`distributor_id` for `Detail`) before calling `controller.Route(app)` and send the existing `Cust_id` header. `JWTProtected().jwtError` allows missing JWT only when that header exists. Test must not bypass JWT in production code. Mitigation: test-only pre-route middleware sets locals and request sends `Cust_id`; production `Route` unchanged.
2. `Detail` may panic on missing locals if test hits detail path. Mitigation: in router-level detail test, pre-route middleware sets all required locals, request sends `Cust_id`, and stub overrides `Detail`.
3. Test must verify "GET report no longer hits detail" with a flag (e.g., `detailCalled`) added to the existing `productServiceStub`; the change in the stub is test-only and limited to this file.
4. `lib/pq` int64 parse failure for literal "report" string is observed behavior; rely on route ordering rather than per-handler int64 hardening, because the handler already returns 422 cleanly when a numeric id parses. Source of int64 parse is `Detail` handler, not the fix target.
5. Live curl with supplied JWT: marked `not-ready` unless executor proves token+runtime; evidence path stores curl attempt log or `not-ready` note.

## Decisions/Assumptions

- Decision: register literal `GET /report` using the existing handler, no new handler. Rationale: handler already reads query and validates filter; behavior is method-agnostic for this endpoint.
- Decision: router-level regression test mounts `controller.Route(app)`. Rationale: only a router-level test can prove dispatch order; handler-direct test cannot.
- Decision: add `detailCalled` flag on `productServiceStub`. Rationale: minimal change, test-only, used only inside this test file.
- Assumption: `Route(app)` ordering inside one method call is preserved by Fiber registration order. Fiber v2 path trie resolves static segments before parameter segments when both are registered; this is the framework contract used in the rest of the repo. Status: `assumption`, low risk given prior plans in this repo rely on the same ordering.
- Assumption: existing JWT secret and runtime state are not part of slice; tests bypass JWT via test middleware, not by changing production route code.
- Open question: none blocking. JWT token and runtime for supplied curl are executor-side; plan records `not-ready` if they are not produced.

## Execution Source of Truth

Precedence for this slice:

1. Handoff payload from `@orchestrator` (preserve POST, no auth weakening, fix GET shadowing).
2. Non-negotiable Implementation Invariants below.
3. Execution-ready Worklist / Handoff Contract.
4. Acceptance Criteria + Done Criteria.
5. Implementation Steps in this plan.
6. Follow-up notes.

Any conflict must be recorded in `.opencode/evidence/fix-product-report-get-route/` before final claim.

## Non-negotiable Implementation Invariants

- Handler `Report` body must NOT be edited. The fix is at the route table only.
- The JWT middleware on `/v1/products` group must remain `middleware.JWTProtected()`. Do not split or weaken.
- `POST /v1/products/report` registration stays. Adding a new `Get` must not remove the existing `Post`.
- Router-level regression test is mandatory; handler-direct tests are not enough to prove the bug is fixed.
- `Detail` behavior remains intact for numeric `:pro_id`. No regression on existing detail flow.
- `go.mod`, `go.sum`, service registration, env, JWT secret: untouched.

## Do Not / Reject If

- Do NOT add a new handler (`GetReport`, etc.). Reuse `Report`.
- Do NOT move `GET /:pro_id` registration above `GET /report` — that re-introduces the bug.
- Do NOT introduce a new package or split `ProductController`.
- Do NOT change `entity.DetailProductParams` or its validator.
- Do NOT touch JWT, env, or any other service.
- Do NOT claim done if router-level test is missing or if `cd master && rtk go test ./...` does not pass.
- Reject if the diff touches anything outside the two files named in Scope.

## Diff Boundary

- Allowed writes: `master/controller/product_controller.go`, `master/controller/product_report_controller_test.go`.
- Allowed artifacts: `.opencode/plans/fix-product-report-get-route.md`, `.opencode/evidence/fix-product-report-get-route/*`, this evidence directory.
- Any out-of-boundary change must be reverted before final quality gate.

## TDD / Test Plan

- Red first: add router-level test asserting `GET /v1/products/report` reaches `Report` and does not call `Detail`. Run → it fails because the route is missing. Run also asserts `GET /v1/products/12345` reaches `Detail` (it does) → this part passes.
- Green: register `productsRouteV1.Get("/report", controller.Report)` before `:pro_id`.
- Refactor: nothing to refactor; the change is one line and one test flag.
- Edge cases covered: missing `cust_id[]` (existing test), invalid `sort_by` (existing test), invalid `sort_order` (existing test), shadowing regression (new test), numeric detail dispatch (new test).
- Validation command: `cd master && rtk go test ./controller -run TestProductReport -v`. Expected: PASS.
- Service-wide validation: `cd master && rtk go test ./...` then `cd master && rtk go build ./...`. Expected: exit 0.
- Live curl: optional; if no JWT/runtime, record `not-ready` note at `.opencode/evidence/fix-product-report-get-route/runtime-curl.md`.

## Implementation Steps

1. Open `master/controller/product_controller.go`. Locate `Route(app *fiber.App)` (line 34).
2. Add `productsRouteV1.Get("/report", controller.Report)` immediately after the existing `productsRouteV1.Post("/report", controller.Report)` and before `productsRouteV1.Get("/"+qParamId, controller.Detail)`. Match existing indentation (tabs).
3. Do not change any other line in this file.
4. Open `master/controller/product_report_controller_test.go`.
5. Extend `productServiceStub` with `detailCalled bool` and a `Detail` method that flips the flag and returns a zero `entity.ProductDetailResponse, nil`. (Stub already embeds `service.ProductService`; override `Detail` so direct calls are intercepted.)
6. Add test `TestProductReportRoute_GET_ReachesReport`:
   a. `app := fiber.New()`.
   b. Pre-route middleware sets `requestid`/`cust_id`/`parent_cust_id`/`distributor_id`/`user_id` locals.
   c. `ctrl := &ProductController{ProductService: svc}; ctrl.Route(app)`.
   d. `req := httptest.NewRequest("GET", "/v1/products/report?cust_id[]=C26002&page=1&limit=20", nil); req.Header.Set("Cust_id", "C26002")` (header is required by `jwtError` to allow missing bearer token).
   e. `res, _ := app.Test(req)`.
   f. Assert `res.StatusCode == fiber.StatusOK`.
   g. Assert `svc.reportCalled == true` and `svc.detailCalled == false`.
7. Add test `TestProductReportRoute_GET_Detail_ReachesDetail`:
   a. Same setup as step 6.
   b. `req := httptest.NewRequest("GET", "/v1/products/12345", nil); req.Header.Set("Cust_id", "C26002")`.
   c. `res, _ := app.Test(req)`.
   d. Assert `svc.detailCalled == true` and `svc.reportCalled == false`.
8. (Optional) `TestProductReportRoute_GET_Report_MissingCustID`: send GET without `cust_id[]`, assert 400. Keeps GET parity with existing POST missing-cust test.
9. Run `cd master && rtk go test ./controller -run TestProductReport -v`. Expect PASS.
10. Run `cd master && rtk go test ./...`. Expect PASS.
11. Run `cd master && rtk go build ./...`. Expect exit 0.
12. If executor has valid JWT and running service, run the supplied GET curl and save body (token redacted) to `.opencode/evidence/fix-product-report-get-route/runtime-curl.md`. Otherwise mark `not-ready`.

## Expected Files to Change

- `master/controller/product_controller.go` (one line added in `Route`).
- `master/controller/product_report_controller_test.go` (stub flag + 2-3 new tests).

## Agent / Tool Routing

- Implementation: `@fixer` (small bounded edit + tests in two files of one service).
- Review gate: `@quality-gate` (final signoff: regression test green, scope clean, auth preserved).
- No need for `@orchestrator` to delegate further; the work is single-track.
- No need for `@designer`/`@architect`/`@librarian` (no UI, no architecture boundary, no version-sensitive API).

## Executor Handoff Prompt

```text
Plan: .opencode/plans/fix-product-report-get-route.md
Task: fix-product-report-get-route
Scope: master service only
Allowed files:
  - master/controller/product_controller.go  (one-line route add)
  - master/controller/product_report_controller_test.go  (router-level tests)
Do not touch: master/service, master/repository, master/entity, master/main.go, master/pkg, master/go.mod/go.sum, master/.env, any other service, JWT middleware.
Must preserve:
  - POST /v1/products/report still works.
  - middleware.JWTProtected() still wraps /v1/products group.
  - Detail handler, entity.DetailProductParams, validator unchanged.
  - responsebuild envelope and paging shape unchanged.
Validation:
  - cd master && rtk go test ./controller -run TestProductReport -v
  - cd master && rtk go test ./...
  - cd master && rtk go build ./...
Evidence:
  - Save test logs under .opencode/evidence/fix-product-report-get-route/
  - If live curl feasible, save redacted body to .opencode/evidence/fix-product-report-get-route/runtime-curl.md; otherwise mark not-ready.
Worker contract: implement only; report back. Do not modify other files. Do not change JWT. Do not claim done if any test fails or diff exceeds the two files above.

## Execution-ready Worklist / Handoff Contract

```text
handoff:
  task_id: fix-product-report-get-route
  plan_id: fix-product-report-get-route
  caller: orchestrator
  callee: @fixer
  scope: add one literal GET /v1/products/report route and a router-level regression test
  claim_level: scoped
  claim_scope:
    may_claim:
      - "GET /v1/products/report reaches Report handler"
      - "POST /v1/products/report still works"
      - "router-level test proves both routes"
    may_not_claim:
      - "service, repository, or DB changes (none expected)"
      - "JWT or env changes (none expected)"
      - "any other service touched"
  source_basis:
    - master/controller/product_controller.go:34-55
    - master/controller/product_controller.go:57-99
    - master/controller/product_controller.go:517-595
    - master/controller/product_report_controller_test.go:13-32
    - master/controller/supplier_controller_test.go:43-80
    - .opencode/docs/PROJECT_STACK.md
    - .opencode/docs/PROJECT_COMMANDS.md
    - .opencode/docs/FRAMEWORK_PLAYBOOK.md
  must_preserve:
    - "POST /v1/products/report compatibility"
    - "middleware.JWTProtected() on /v1/products group"
    - "Controller -> Service -> Repository layering"
    - "Existing response/paging envelope"
  do_not_touch:
    - "master/service/**"
    - "master/repository/**"
    - "master/entity/**"
    - "master/main.go"
    - "master/pkg/**"
    - "master/go.mod"
    - "master/go.sum"
    - "master/.env"
    - "any other service or middleware"
  validation:
    - "cd master && rtk go test ./controller -run TestProductReport -v"
    - "cd master && rtk go test ./..."
    - "cd master && rtk go build ./..."
  exit_criteria:
    - "GET /v1/products/report reaches Report handler (test green)"
    - "GET /v1/products/:pro_id numeric still dispatches to Detail (test green)"
    - "POST tests still pass"
    - "diff limited to two files"
  evidence_required:
    - ".opencode/evidence/fix-product-report-get-route/test-controller.log"
    - ".opencode/evidence/fix-product-report-get-route/test-all.log"
    - ".opencode/evidence/fix-product-report-get-route/build.log"
    - ".opencode/evidence/fix-product-report-get-route/runtime-curl.md (optional, not-ready acceptable)"
  depends_on: []
  context_bundle:
    verified_by_planner:
      - "Route registration order: GET /:pro_id currently shadows GET /report. (confirmed_repo, master/controller/product_controller.go:43-49)"
      - "Report handler parses cust_id[] from query and is method-agnostic. (confirmed_repo, master/controller/product_controller.go:517-595)"
      - "Detail handler tries int64 parse on :pro_id and fails for 'report'. (confirmed_repo, master/controller/product_controller.go:65-77)"
      - "Existing report test uses handler-direct mount, not router mount. (confirmed_repo, master/controller/product_report_controller_test.go:25-32)"
      - "JWT middleware still required on group; test must set c.Locals without weakening production auth. (confirmed_repo, master/pkg/middleware/jwt_middleware.go:15-32)"
      - "Per-service go test/build commands with rtk prefix. (confirmed_docs, .opencode/docs/PROJECT_COMMANDS.md:18-22,40-44)"
    files_already_read:
      - "master/controller/product_controller.go"
      - "master/controller/product_report_controller_test.go"
      - "master/controller/supplier_controller_test.go"
      - "master/entity/product.go"
      - "master/service/product_service.go"
      - "master/pkg/middleware/jwt_middleware.go"
      - ".opencode/docs/PROJECT_STACK.md"
      - ".opencode/docs/PROJECT_COMMANDS.md"
      - ".opencode/docs/FRAMEWORK_PLAYBOOK.md"
      - ".opencode/docs/PROJECT_DETECTED_TOOLS.md"
      - ".opencode/docs/MCP.md"
    open_assumptions:
      - "Fiber v2 path trie resolves static segments before parameter segments. (assumption, low risk)"
      - "No JWT token supplied to executor; runtime curl marked not-ready if env not available."
    source_of_truth_order:
      - "task handoff payload"
      - "this plan"
      - "source files in master/controller/"
```

## Progress Tracking

- tracker_path: `.opencode/state/fix-product-report-get-route/progress.json`
- init_command: `python3 ~/.config/opencode/scripts/task-progress.py fix-product-report-get-route --init --plan .opencode/plans/fix-product-report-get-route.md`
- summary_command: `python3 ~/.config/opencode/scripts/task-progress.py fix-product-report-get-route --summary`
- checklist_command: `python3 ~/.config/opencode/scripts/task-progress.py fix-product-report-get-route --checklist`
- update_rules: update `pending -> in_progress` immediately before starting; update to `completed`/`blocked`/`cancelled` immediately after; update whenever evidence is written; update at every cross-lane handoff.
- task_map:

1. **A1** | `@fixer` | Add `productsRouteV1.Get("/report", controller.Report)` in `Route(app)`, before `:pro_id` line. | evidence: `.opencode/evidence/fix-product-report-get-route/A1-route-registration.log`
2. **A2** | `@fixer` | Extend `productServiceStub` with `Detail` override + `detailCalled` flag; add `TestProductReportRoute_GET_ReachesReport`. | evidence: `.opencode/evidence/fix-product-report-get-route/A2-test-get-report.log`
3. **A3** | `@fixer` | Add `TestProductReportRoute_GET_Detail_ReachesDetail` and optional `TestProductReportRoute_GET_Report_MissingCustID`. | evidence: `.opencode/evidence/fix-product-report-get-route/A3-test-get-detail.log`
4. **A4** | `@fixer` | Run `cd master && rtk go test ./controller -run TestProductReport -v`, `rtk go test ./...`, `rtk go build ./...`; save logs. | evidence: `.opencode/evidence/fix-product-report-get-route/test-controller.log`, `test-all.log`, `build.log`
5. **A5** | `@fixer` | If runtime available, run supplied GET curl, save redacted body to `runtime-curl.md`; else mark `not-ready`. | evidence: `.opencode/evidence/fix-product-report-get-route/runtime-curl.md`

| id | owner | evidence path | update command |
|---|---|---|---|
| A1 | @fixer | `.opencode/evidence/fix-product-report-get-route/test-controller.log` | `python3 ~/.config/opencode/scripts/task-progress.py fix-product-report-get-route --update A1 --status completed --owner @fixer --evidence .opencode/evidence/fix-product-report-get-route/test-controller.log` |
| A2 | @fixer | `.opencode/evidence/fix-product-report-get-route/test-all.log` | `python3 ~/.config/opencode/scripts/task-progress.py fix-product-report-get-route --update A2 --status completed --owner @fixer --evidence .opencode/evidence/fix-product-report-get-route/test-all.log` |
| A3 | @fixer | `.opencode/evidence/fix-product-report-get-route/build.log` | `python3 ~/.config/opencode/scripts/task-progress.py fix-product-report-get-route --update A3 --status completed --owner @fixer --evidence .opencode/evidence/fix-product-report-get-route/build.log` |
| A4 | @fixer | `.opencode/evidence/fix-product-report-get-route/test-controller.log`, `test-all.log`, `build.log` | `python3 ~/.config/opencode/scripts/task-progress.py fix-product-report-get-route --update A4 --status completed --owner @fixer --evidence .opencode/evidence/fix-product-report-get-route/build.log` |
| A5 | @fixer | `.opencode/evidence/fix-product-report-get-route/runtime-curl.md` (optional) | `python3 ~/.config/opencode/scripts/task-progress.py fix-product-report-get-route --update A5 --status completed --owner @fixer --evidence .opencode/evidence/fix-product-report-get-route/runtime-curl.md` |

## Validation Commands

- `cd master && rtk go test ./controller -run TestProductReport -v` — expect PASS with all `TestProductReport*` tests green.
- `cd master && rtk go test ./...` — expect exit 0.
- `cd master && rtk go build ./...` — expect exit 0.
- Optional runtime: `curl -H "Authorization: Bearer <token>" "http://localhost:9002/v1/products/report?cust_id[]=C26002&page=1&limit=20"` — expect 200 with master envelope (only if executor has valid JWT + running service).

## Evidence Requirements

- `.opencode/evidence/fix-product-report-get-route/discovery.md` (already written).
- `.opencode/evidence/fix-product-report-get-route/test-controller.log` — `rtk go test -run TestProductReport -v` output.
- `.opencode/evidence/fix-product-report-get-route/test-all.log` — `rtk go test ./...` output.
- `.opencode/evidence/fix-product-report-get-route/build.log` — `rtk go build ./...` output.
- `.opencode/evidence/fix-product-report-get-route/runtime-curl.md` — live curl evidence or `not-ready` note.

## Done Criteria

- All Acceptance Criteria checked.
- All Validation Commands exit 0.
- Diff limited to `master/controller/product_controller.go` and `master/controller/product_report_controller_test.go`.
- No JWT, env, or other service changes.
- Evidence logs saved at the paths above.
- @quality-gate signoff received.

## Final Planning Summary

- Artifacts consulted: `master/controller/product_controller.go`, `master/controller/product_report_controller_test.go`, `master/controller/supplier_controller_test.go`, `master/entity/product.go`, `master/service/product_service.go`, `master/pkg/middleware/jwt_middleware.go`, `.opencode/docs/PROJECT_STACK.md`, `.opencode/docs/PROJECT_COMMANDS.md`, `.opencode/docs/FRAMEWORK_PLAYBOOK.md`, `.opencode/docs/PROJECT_DETECTED_TOOLS.md`, `.opencode/docs/MCP.md`, prior plans `20260714-1315-sx-2513-product-secondary-sales-report.md` and `20260715-sx-2516.md` (context only).
- Artifacts created: `.opencode/plans/fix-product-report-get-route.md`, `.opencode/evidence/fix-product-report-get-route/discovery.md`.
- Key decisions: literal GET route registration before parameter route, no handler edit, router-level regression test, test-local JWT bypass via Locals only.
- Assumptions: Fiber v2 static-before-parameter resolution; supplied curl has runtime/token available to executor.
- Open questions: none blocking. Live curl evidence is optional and labeled `not-ready` if env not present.
- Readiness: `PASS_FOR_SLICE`.
- Cleanup: no stale drafts to remove; this plan is the single primary plan.

---

Active-lane reset note: planner read-only restrictions do not persist into the next lane. The next active lane (`@orchestrator` or `@fixer`/`@backend` implementation lane) must refresh its own permissions and context before editing source files.
