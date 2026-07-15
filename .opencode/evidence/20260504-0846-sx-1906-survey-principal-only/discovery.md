# Discovery Evidence — SX-1906 Survey Principal-only BU

## Files inspected

- `AGENTS.md` — repo constraints: multi-module Go services, `master` service on port `9002`, Fiber `/ping`, strict controller → service → repository → DB flow, transaction requirement for writes, tenant filtering guidance.
- `master/controller/survey_controller.go` — `POST /v1/survey` is registered under JWT middleware; with root prefix this maps to `POST /master/v1/survey`. Controller unmarshals `entity.CreateSurveyBody`, injects `cust_id`, `parent_cust_id`, and `user_id`, validates, then calls `SurveyService.Store`.
- `master/entity/survey.go` — `CreateSurveyBody.DistributorId` and `EmpId` use `FlexibleIntArray`, so both scalar and array JSON are accepted. `distributor_id` currently has no validation tag that blocks `0`.
- `master/service/survey_service.go` — `Store` normalizes `distributor_id` through `normalizePositiveInts`, calls `FindSurveyAreasByDistributorIds`, builds survey area rows with `buildSurveyAreas`, resolves salesman cust scope with `resolveSurveyCustIds`, validates salesman, then writes survey, areas, salesmen, outlets, and details inside a transaction.
- `master/repository/survey_repository.go` — area rows insert into `mst.m_survey_area (survey_id, distributor_id, area_id)`. Area lookup only reads areas from `mst.m_distributor` for positive distributor IDs. Salesman rows insert into `mst.m_survey_salesman` with resolved `cust_id`.
- `master/service/survey_service_test.go` — existing tests already cover ignoring `0` when mixed with real distributor IDs, rollback on salesman failures, normal distributor area mapping, and detail/list behavior.
- `master/migration/mst.survey/001_create_tables.sql` — `mst.m_survey_area.distributor_id INT NOT NULL`, no FK to distributor declared in local migration.
- `master/migration/mst.survey/002_add_distributor_and_salesman.sql` — historical migration sets NULL `distributor_id` to `0`, indicating schema has tolerated sentinel/placeholder `0` in survey area rows.
- `docs/Create Survey_BE.md` — create survey documentation lists `distributor_id` body field; later outlet section documents `distributor_id` example where `0 = principal` and positive values are distributor IDs.

## Project patterns found

- Survey create/update keeps writes transactional in service layer via `txManager.WithinTransaction` and repository write methods.
- Existing service already normalizes non-positive distributor IDs out of distributor lookups.
- `buildSurveyAreas` currently returns `ErrSurveyAreaDistributorRequired` when `area_id` is present but normalized distributor IDs and distributor-derived areas are empty.
- `resolveSurveyCustIds` returns the request `custId` when no positive distributor IDs are available, which matches principal-owned survey context.
- Salesman validation uses `FindOneByEmpIdAndCustId` for each resolved cust ID and `parentCustId`.

## Reuse candidates

- Reuse `normalizePositiveInts` to keep `0` out of distributor FK/query paths.
- Extend `buildSurveyAreas` or add a small helper to represent principal-only area mappings without treating `0` as real distributor ID.
- Reuse existing `surveyRepositoryRedStub`, `transactionManagerStub`, and `salesmanRepositoryStub` tests for a red/green unit test around `distributor_id: [0]` and selected `area_id`/`emp_id`.

## Commands/docs checked

- `docker compose -f docker-compose.yml ps` showed `master`, `system`, and `redis` running. Per global OpenCode instruction, commands were run without `rtk`; repo `AGENTS.md` asks for `rtk`, but higher-priority global instruction says not to prefix OpenCode commands with `rtk` unless explicitly requested.
- Local static discovery only; no external library docs needed because behavior is local Go service logic and SQL schema.
- No GitHub/web/browser research needed; issue is backend-only and local evidence/docs are sufficient.

## Constraints

- Do not remove distributor validation/normalization globally.
- Do not insert/query distributor ID `0` as a real distributor; if `0` is stored, it must be treated as principal-only sentinel/placeholder because local schema has no FK and migration already backfilled `0`.
- Preserve existing behavior for positive distributor IDs, outlet-specific survey, normal distributor flow, and invalid salesman/distributor behavior.
- Tests should be run in `master/` module, not repo root, because each service has its own `go.mod`.

## Risks

- Domain ambiguity remains for exact persisted representation of selected areas in principal-only surveys. Current schema requires `distributor_id NOT NULL`; safest local-compatible representation is `distributor_id = 0` with selected `area_id`, but user acceptance says avoid treating `0` as real distributor FK except if schema supports sentinel. Local migration evidence supports sentinel usage.
- Existing test `TestSurveyService_Store_ShouldIgnoreZeroDistributorAndResolveAreas` expects mixed `[67,0,68]` to ignore `0`; implementation must not accidentally create extra sentinel rows for mixed payloads unless domain owner confirms mixed principal + distributor is allowed.
- Cannot reproduce staging failure without authorization token and staging data. Plan should require local test reproduction at service level and optional authenticated staging/local HTTP verification by implementer.
