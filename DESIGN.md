# Scylla Backend — API Surface Design System

> Category: multi-module Go backend services for distribution/field-sales operations.
> Audience: backend engineers and integrators consuming HTTP APIs from `system`, `master`, `inventory`, `sales`, `finance`, `tms`, `mobile`, `pjp`, `cronjob`.
> Source: derived from existing service conventions (`master/pkg/responsebuild/response.go`, controller patterns, JWT middleware, `entity.Pagination`, `pkg/sql_helper/sql_patch.go`) and project-wide rules in `.opencode/docs/ARCHITECTURE.md`, `QUALITY.md`, `SECURITY.md`.

This file replaces a visual design system. Repo contains no UI/frontend; all surfaces are HTTP APIs. "Design" here means the **shape, color, and depth of API behavior**, not pixels.

## Visual Theme & Atmosphere

The API surface must feel:

- **Predictable**: same envelope, paging, error, and authentication shape across every service. No creative deviations per controller.
- **Tight and operational**: minimal payload bloat, no decorative fields, every response field has a downstream consumer.
- **Tenant-aware but not tenant-noisy**: tenant (`cust_id`) appears in requests, errors, and audit context, never as a marketing field.
- **Idiomatic to Go HTTP ecosystem**: uses Fiber or Gin native constructs, framework-provided middleware, and Go stdlib pagination math (`CalculateLastPage`).
- **Audit-friendly**: every response carries a `request_id` for cross-service tracing; no log lines omit it.

Reject the opposite atmosphere: a "RESTful" surface that mixes snake_case and camelCase, returns wrapped responses for some endpoints and raw objects for others, or fabricates fictional fields for symmetry.

## Color Palette & Roles

This is a token map for response keys and HTTP behavior, not visual colors. Each role is mandatory and non-negotiable unless the active plan states otherwise.

| Token | Role | Default value | Notes |
|---|---|---|---|
| `MESSAGE` | Human-readable top-level message | string | Localized via `texttranslator` (EN/ID). |
| `DATA` | Payload | object/array/null | Omitted from response when nil (`omitempty`). |
| `ERRORS` | Error map | object/null | Field-level or single error shape; never mixed with `DATA`. |
| `PAGING` | Pagination object | object/null | Required for list endpoints; omitted for single-resource responses. |
| `REQUEST_ID` | Cross-service correlation | string | Always emitted; never empty. |
| `STATUS` | HTTP status | 200/400/401/403/404/409/422/500 | Use the most specific applicable status. |

Binding (in `ApiPayload`): `Message`, `Data`, `Errors`, `Paging`, `RequestId` follow Fiber's `json` tags with `omitempty`; `request_id` is always present. Source: `master/pkg/responsebuild/response.go:18-24`.

## Typography Rules

This is a response-content and naming convention map, not typography. The repo speaks in two natural languages and two naming dialects. Pick consistently.

- Field names: `snake_case` for JSON keys, `CamelCase` for Go struct fields. Do not mix in the same struct.
- Error keys: prefer `code` (machine-readable enum) + `message` (human-readable, localized). Avoid mixing verbose English sentences into machine fields.
- Identifier prefixes: `mp_` for `m_product`-derived fields when reused in joins, `parent_` for the principal-parent join, `original_` for distributor-side identity preserved alongside normalized identity. Do not invent alternative prefixes.
- Enums (e.g. `type` in product report): `Own Products`, `Product Assigned`, `Product Mapping`. Match docs exactly; do not auto-format or lowercase.
- i18n: message strings flow through `texttranslator` for `id`. Worker code passes the language code, never the translated string. Source: `master/pkg/texttranslator/...`.

## Component Stylings

Each "component" below is a mandatory building block for new endpoints. Reuse instead of inventing.

- **Endpoint pattern**: HTTP verb + path; controllers register within the service's existing route group; preserve service's existing prefix (compose-level `master`, `sales`, `tms`, etc.). Do not register internal prefixes inside the service.
- **Authentication wrapper**: existing service `middleware.JWTProtected()` (Fiber) or `utils.DecodeJWT` (Gin). The protected route group is the only place this is added. Do not bypass or recreate it.
- **Request parser**: Fiber's `c.Locals("jwt")` and `c.Queries()` for query parameters; Gin's `c.QueryArray` and `c.PostForm`. Do not introduce a third parser.
- **Response builder**: `responsebuild.BuildResponse(requestID, lang)` in Fiber services. Gin services must keep the same `message/data/errors/paging/request_id` shape locally.
- **Error translator**: `texttranslator` for human messages; `entity.ApiResponse` for non-200 paths. Keep the JSON shape consistent across success and failure.
- **Paging helper**: `sql_helper.CalculateLastPage` (or equivalent) for total pages. Do not roll a new formula.
- **List query style**: `master` uses sqlx/pgx pattern with bound parameters; GORM services use existing repository interfaces. Match the target service.
- **Validation**: `go-playground/validator` with struct tags; missing field maps to 400 with `errors` populated.

## Layout Principles

- **Endpoint grouping**: keep endpoints under existing service path prefixes; new endpoint joins the controller family it belongs to (e.g. `product` family under `/v1/products`).
- **Layered composition**: Controller → Service → Repository → DB. Controllers never call repositories; repositories never hold business logic; only services orchestrate transactions.
- **Stateless requests**: every endpoint reads required context from request + JWT; do not rely on hidden process state. The `cust_id` scope is request-driven and explicit; do not silently substitute JWT locals.
- **Request body vs query parameters**: if the active plan calls for a query-parameter endpoint (e.g. `POST .../report?cust_id[]=...`), follow the spec. Do not silently switch to a JSON body contract.
- **Paging layout**: `paging.total_record`, `paging.page_current`, `paging.page_limit`, `paging.page_total` — preserve the established order; do not rename fields.
- **Idempotency**: write operations are scoped to a transaction inside the service. Retries from clients should rely on idempotency keys only when explicitly designed.

## Depth & Elevation

How "elevated" or "deep" an endpoint is in the system's responsibilities.

- **Flat endpoints** (depth 0): health, version, public config. No tenant context, no auth. Future use only if explicitly approved.
- **Standard endpoints** (depth 1): business reads and writes under JWT, single service, single DB. The default for all new routes. Wrapped in `responsebuild`, audited via `request_id`.
- **Aggregated endpoints** (depth 2): reports, search, join-heavy queries spanning `m_product` + `m_distributor` + parent principal rows. Use composite parameter binding; protect with parameter allowlist; document the join in plan evidence.
- **Cross-service endpoints** (depth 3): call out to other services (`tms` referencing `master`/`sales`/`mobile` per compose). Wrap with retry/time-out policy; never trust downstream response shape changes without coordinated plan.

Do not add depth 3+ endpoints without a plan-level architecture decision; route to `@architect` instead.

## Do's and Don'ts

**Do**

- Use parameterized SQL placeholders for every user-controlled value (`q`, `cust_id[]`, `limit`, `offset`, sort).
- Allowlist `sort_by` to a closed Go map; default and unknown values fall back to the safe default.
- Reuse `entity.Pagination` and `sql_helper.CalculateLastPage` instead of inventing new shapes.
- Emit `original_*` as JSON `null` (not empty string, not zero) when no normalized mapping is applied.
- Keep `request_id` non-empty on every response.
- Read `cust_id` scope from the validated request only; never overwrite from `c.Locals("cust_id")` unless the active plan says so.
- Cover new endpoints with controller tests (using `httptest`) and repository tests (using `sqlmock` for `master`, GORM/DB tests elsewhere).

**Don't**

- Don't log bearer tokens, JWT secrets, DSN strings, or full env values in evidence.
- Don't commit `.env`, `*.pem`, `*.key`, or backup dump files.
- Don't add migrations to services that don't have an existing migration practice.
- Don't switch between sqlx and GORM in the same module; respect the target module's data layer.
- Don't use raw Fiber patterns inside Gin services, or vice versa.
- Don't wrap errors in opaque 200 responses; surface HTTP 4xx/5xx with `errors` populated.
- Don't claim runtime/Staging evidence when env/DB/token are not configured; mark `not-ready` and stop.

## Responsive Behavior

Backend analog: **resilience to input variability** across client types (web, mobile, integration scripts, third-party).

- **Input format tolerance**: `cust_id` may arrive as repeated query keys (`cust_id[]=X&cust_id[]=Y`) or comma-separated depending on caller; parser must normalize and reject empty entries. Decision and test go in the plan, not in mid-execution.
- **Pagination input**: `page<1` normalizes to `1`; `limit` has an explicit maximum (small ceiling, e.g. 100) consistent with target service convention. Document the ceiling in the plan.
- **Auth fault tolerance**: missing or malformed JWT returns 400 with `Missing or malformed JWT` only when a valid `cust_id` header is present; otherwise 401. Source: `master/pkg/middleware/jwt_middleware.go:35-53`.
- **Query boundaries**: large `q` strings or very wide `IN` lists must be capped; use `limit`/offset as the primary throttle.
- **Reduced-motion (data refresh)**: prefer `Cache-Control: no-cache` only on endpoints that must always be live; do not add caching headers without plan approval.

## Agent Prompt Guide

- **Read first**: target service `controller/<feature>_controller.go`, `service/<feature>_service.go`, `repository/<feature>_repository.go`, target entity file, `pkg/responsebuild/response.go` (or service-local equivalent), `entity/api.go` for `Pagination`, `pkg/sql_helper/sql_patch.go` for paging math, `pkg/middleware/jwt_middleware.go` for auth boundaries.
- **Apply this design system** by: (1) choosing the right depth (Section "Depth & Elevation"), (2) following the layered composition, (3) reusing the seven components above without re-implementing, (4) running `cd <service> && rtk go test ./...` and `rtk go build ./...` before claiming done.
- **Reuse vs create**: reuse the response builder, error translator, paging helper, JWT middleware, and validator. Create only narrow DTO/entity files when existing types cannot express nullability or new fields without contaminating other endpoints.
- **When to ask questions**: when (a) the request is multi-tenant but the plan is unclear about the source of `cust_id`, (b) the user asks for new auth/scope semantics, (c) migration/DDL is implied for a module without a migration practice, (d) a new external dependency is proposed.
- **When to run visual/render validation**: not applicable (no UI). Instead, run controller `httptest` suite, repository `sqlmock`/`go test ./repository`, and an env-gated curl smoke when DB + token are configured. Do not claim "verified" if curl was skipped; mark `not-ready`.
- **Reporting deviations**: if any field/key/path/casing must diverge from this guide, document the deviation in the active plan's `## Decisions/Assumptions` and link evidence; do not silently drift.
- **Downstream ownership chain for substantial API work**: `@orchestrator` routes/integrates; `@designer` is not the default lane for backend work but may be consulted when the surface includes a UI/UX design that touches API shape; `@fixer` (or `@backend`) implements bounded changes with tests; `@artifact-planner` writes the durable plan and may consult `@designer` as read-only advisory for cross-channel concerns; `@oracle` reviews architecture tradeoffs; `@quality-gate` is final signoff for non-trivial/prompt/config/security-sensitive changes.
