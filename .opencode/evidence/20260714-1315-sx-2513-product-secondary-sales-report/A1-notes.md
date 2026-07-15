A1 notes

changed paths:lines
- `master/entity/product.go:37-56` — added `ProductReportQueryFilter` and `ProductReportResponse`; `original_*` fields use pointers.
- `master/service/product_service.go:27-31,104-107` — added `ReportList` interface contract and temporary concrete implementation required for current `ProductService` construction to compile. A2 must replace with repository-backed implementation.
- `master/controller/product_controller.go:47,517-594` — registered `POST /v1/products/report` before `:pro_id`; added query parsing, cust_id validation, sort allowlist/order normalization, service call, and paging envelope.
- `master/controller/product_report_controller_test.go:1-185` — added five focused controller tests.

unverified facts
- Fiber/fasthttp `QueryArgs().PeekMulti("cust_id[]")` preserves repeated `cust_id[]` values and order; controller test verifies two ordered values.
- Exact requested command `rtk go test ./controller -run Product.Report -v` is not recognized by local `rtk` wrapper and reports `Go test: No tests found`; equivalent Go regex `rtk go test ./controller -run '^TestProductReport' -v` passes all 5 tests.
- No `PROJECT_STACK.md`, `PROJECT_COMMANDS.md`, `FRAMEWORK_PLAYBOOK.md`, or `PROJECT_DETECTED_TOOLS.md` exists in `.opencode/docs`; repo-local AGENTS and plan used instead.

plan deviations
- A1 adds temporary `productServiceImpl.ReportList` returning empty data so existing `NewProductService` satisfies expanded interface. Controller tests use fake service and verify call path. This is intentionally incomplete and must be replaced by A2 repository/service implementation; no repository or migration files touched.
- Limit values below 1 are not normalized/rejected in A1 because handoff explicitly binds only `page<1` normalization and A2 owns query behavior.
- JWT middleware is not exercised in unit tests; test route directly injects request ID, matching existing controller test pattern while validating controller scope does not read JWT locals.
