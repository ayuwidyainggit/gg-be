# Discovery — fix-product-report-get-route

Task ID: `fix-product-report-get-route`
Mode: maintenance-stability
Planner scope: plan/evidence only. Source untouched.

## Source strategy

- Used: repo-local controller, entity, service interface, controller tests, architecture/stack/commands/playbook/tool/MCP docs.
- Skipped: Context7/official Fiber docs. Reason: no version-sensitive API choice; minimal method registration follows existing Fiber route API already used in target file.
- Skipped: GitHub/web/browser/runtime curl. Reason: routing regression reproducible with Fiber `httptest`; live curl needs valid JWT and service runtime, so implementation evidence must label it `not-ready` if unavailable.
- Skipped: compose execution during planning. Reason: user requested plan only; no source/runtime change needed to identify static route shadowing.

## Files inspected

| Path | Finding | Verification |
|---|---|---|
| `master/controller/product_controller.go:34-54` | Protected `/v1/products` group registers `POST /report`, then `GET /:pro_id`; no literal `GET /report` exists. | `confirmed_repo` |
| `master/controller/product_controller.go:57-76` | `Detail` parses `:pro_id` into `entity.DetailProductParams.ProductId int64`; nonnumeric `report` fails parser. | `confirmed_repo` |
| `master/controller/product_controller.go:517-595` | `Report` parses repeated `cust_id[]`, applies query defaults/validation, calls `ProductService.ReportList`, and returns existing response/paging envelope. | `confirmed_repo` |
| `master/controller/product_report_controller_test.go:25-185` | Existing report tests mount only `POST /v1/products/report`; no router-conflict test covers `GET /report`. | `confirmed_repo` |
| `master/entity/product.go:37-44,1049-1054` | Report filter supports `cust_id[]`, `q`, paging, sorting; detail param `pro_id` is `int64`. | `confirmed_repo` |
| `master/service/product_service.go:27-31` | Controller delegates report work through `ProductService.ReportList`; no controller-to-repository bypass needed. | `confirmed_repo` |
| `master/pkg/middleware/jwt_middleware.go:15-32` | `JWTProtected()` configures protected routes; endpoint must remain inside existing group. | `confirmed_repo` |
| `.opencode/docs/PROJECT_STACK.md:25-46` | `master` is Go 1.20/Fiber; controller-service-repository and response envelope conventions apply; per-service test/build commands required. | `confirmed_docs` |
| `.opencode/docs/PROJECT_COMMANDS.md:1-44` | Commands run in target service and require `rtk` prefix. | `confirmed_docs` |
| `.opencode/docs/FRAMEWORK_PLAYBOOK.md:23-33,71-77` | Route belongs in controller group; preserve JWT group; controller tests reuse `httptest`; manual edit is valid for existing hand-authored controller. | `confirmed_docs` |

## Confirmed vs Assumed Audit

| Claim | Level | Basis |
|---|---|---|
| `GET /v1/products/report` currently has no literal registration. | `confirmed_repo` | `product_controller.go:43-49` |
| GET path may match `GET /:pro_id` and parser attempts `report` as int64. | `confirmed_repo` | registration plus `DetailProductParams.ProductId int64` and `Detail` parser |
| `POST /v1/products/report` must remain compatible. | `user_confirmed` | task handoff |
| Same `Report` handler can serve GET without changing filter/service/repository behavior. | `confirmed_repo` | handler reads query only, existing POST tests send no body |
| Existing JWT policy must stay intact. | `confirmed_repo` + `user_confirmed` | route group/middleware and handoff |
| Supplied curl has valid JWT/runtime config. | `unverified` | token/runtime not supplied; runtime smoke conditional |

## Reuse candidates

- `ProductController.Route`: add one literal route registration; no new handler, DTO, service, repository, dependency, or middleware.
- `productServiceStub` and `TestProductReport_ValidRequest_CallsService`: extend existing test file and stub.
- `net/http/httptest` + `fiber.New()`: existing controller test pattern.

## Risks

1. Test that mounts `ctrl.Report` directly cannot prove dispatch order. Regression test must call `ctrl.Route(app)` so literal route and parameter route compete.
2. `ctrl.Route(app)` installs JWT middleware. Test needs a local app-level test middleware before registration to provide required locals and bypass JWT only inside unit-test app setup; production route code must not bypass middleware. Alternative: test registration in isolated Fiber app with a valid signed test JWT if existing JWT config permits it. Executor must choose existing local pattern, preserve production auth.
3. `Detail` may need a stub implementation if router regression test verifies numeric detail dispatch. Embedded `service.ProductService` in current stub can provide promoted methods; explicit detail assertion may need a dedicated spy stub or a handler-level marker.

## Artifact validator fallback

`python3 scripts/plan-execution-readiness.py` exists under `~/.config/opencode/scripts/`, not target repo `scripts/`. Plan uses absolute config-root script path for planner validation.
