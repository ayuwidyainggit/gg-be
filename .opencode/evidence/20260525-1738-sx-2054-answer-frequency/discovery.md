# Discovery — SX-2054 Answer Frequency

Task ID: `20260525-1738-sx-2054-answer-frequency`

## Ringkasan

Discovery lokal menemukan dua blocker utama untuk SX-2054:

1. `master/migration/mst.survey/001_create_tables.sql` mendefinisikan `mst.m_survey.answer_frequency VARCHAR(20) NOT NULL`. Nilai baru `Multiple Times, One Day` dan `Multiple Times, Different Day` melebihi 20 karakter, sehingga insert/update akan gagal walau tidak ada `CHECK` constraint.
2. `master/entity/survey.go` masih memakai validator `validate:"required,oneof=Multiple 'One Time'"` pada `CreateSurveyBody.AnswerFrequency` dan `UpdateSurveyBody.AnswerFrequency`, sehingga dua nilai baru akan ditolak di API boundary.

## Files inspected

- `AGENTS.md`
- `.opencode/docs/index.md`
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `docker-compose.yml`
- `master/go.mod`
- `master/entity/survey.go`
- `master/model/survey.go`
- `master/controller/survey_controller.go`
- `master/controller/survey_controller_test.go`
- `master/service/survey_service.go`
- `master/service/survey_service_test.go`
- `master/repository/survey_repository.go`
- `master/repository/survey_report_repository.go`
- `master/model/survey_report.go`
- `master/entity/survey_report.go`
- `master/migration/mst.survey/001_create_tables.sql`
- `master/migration/mst.survey/002_add_distributor_and_salesman.sql`
- `master/migration/mst.survey/003_add_survey_salesman_guardrails.sql`
- `mobile/repository/survey.go`
- `mobile/model/survey.go`
- `mobile/service/survey.go`
- `mobile/migration/202603091235/ddl-survey.sql`

## Commands / tools checked

- `rtk docker compose -f docker-compose.yml ps`
  - Output: compose file parsed, no services running.
  - Warning: `.rtk/filters.toml` untrusted, filters not applied.
- Local content search for:
  - `answer_frequency`
  - `AnswerFrequency`
  - `One Time`
  - `Multiple Times`
  - `Multiple`
  - `CreateSurvey`
  - `UpdateSurvey`
  - `GetSurvey`
  - `SurveyDetail`
  - `mst.m_survey`
- Explorer read-only discovery.
- Oracle read-only migration/risk review.

## Project patterns found

- Repo is multi-module Go monorepo.
- Target service for create/update API is `master`.
- `master` uses Fiber, controller/service/repository layering.
- Survey route is mounted as `/v1/survey`; via gateway/docs it corresponds to `/master/v1/survey`.
- Write flow already uses service-layer transaction through `txManager.WithinTransaction(...)`.
- Repository already stores and reads `survey.AnswerFrequency` raw without mapping.
- List/detail/report/mobile read paths pass through stored `answer_frequency` string raw.
- Existing migration style for survey is raw SQL under `master/migration/mst.survey/` with numbered files.
- Existing tests use stubs in `master/service/survey_service_test.go` and controller stubs in `master/controller/survey_controller_test.go`.

## Key code points

### API request validation

- `master/entity/survey.go:127`
  - `CreateSurveyBody.AnswerFrequency string 'json:"answer_frequency" validate:"required,oneof=Multiple 'One Time'"'`
- `master/entity/survey.go:145`
  - `UpdateSurveyBody.AnswerFrequency string 'json:"answer_frequency" validate:"required,oneof=Multiple 'One Time'"'`

### Create/update persistence

- `master/service/survey_service.go:474`
  - create model uses `AnswerFrequency: request.AnswerFrequency`
- `master/service/survey_service.go:582`
  - update model uses `AnswerFrequency: request.AnswerFrequency`
- `master/repository/survey_repository.go:207-216`
  - `INSERT INTO mst.m_survey (... answer_frequency ...)`
- `master/repository/survey_repository.go:220-227`
  - `UPDATE mst.m_survey SET ... answer_frequency = $2 ...`

### Read paths

- `master/repository/survey_repository.go:122-128`
  - list selects `answer_frequency` raw.
- `master/repository/survey_repository.go:138-155`
  - detail selects `s.answer_frequency` raw.
- `master/service/survey_service.go:222`
  - list response maps `sv.AnswerFrequency` raw.
- `master/service/survey_service.go:245`
  - detail response maps `survey.AnswerFrequency` raw.
- `master/repository/survey_report_repository.go:58` and `178`
  - report list/detail selects `s.answer_frequency` raw.
- `mobile/repository/survey.go:62` and `123`
  - mobile list/detail selects `ms.answer_frequency` raw.
- `mobile/model/survey.go`
  - mobile response structs expose `AnswerFrequency string`.

### Schema

- `master/migration/mst.survey/001_create_tables.sql:10`
  - `answer_frequency VARCHAR(20) NOT NULL`
- No `CHECK` constraint found for `mst.m_survey.answer_frequency` in survey migrations.
- No Postgres enum found for `answer_frequency` in local migrations.

## Reuse candidates

- Reuse existing raw pass-through repository/service mapping; do not add read normalization.
- Reuse transaction flow in `Store()` and `Update()`.
- Reuse validation boundary in controller, but move value list away from fragile inline validator tag if custom validation/helper pattern is acceptable.
- Reuse survey migration folder `master/migration/mst.survey/` for schema widen migration.
- Reuse controller/service test patterns already present.

## Constraints

- New write values must be exact:
  - `One Time`
  - `Multiple Times, One Day`
  - `Multiple Times, Different Day`
- Legacy `Multiple` must not be mass-migrated without product decision.
- Read path should remain tolerant of legacy `Multiple`.
- Do not commit Authorization/cookie/token/credential from docs/browser/cURL evidence.
- Follow `Controller → Service → Repository → DB`.
- Validate in target service directory `master/` using `rtk` commands.

## Risks

1. `VARCHAR(20)` will truncate/fail new values. Must widen before or with app release.
2. Validator currently rejects required new values.
3. Inline `oneof` tags are brittle for comma/spaced values and duplicate create/update rules.
4. Adding strict DB `CHECK` immediately can break unknown legacy values if data contains values outside known set.
5. Mobile submit flow currently blocks duplicate submissions for all surveys via `CheckExistingSubmission(...)`; no code evidence found that it differentiates `One Day` vs `Different Day`.
6. `take_survey` in mobile list/detail is based on existing answer, not new frequency semantics.
7. `FindAllByCustId` accepts raw `sort` in `ORDER BY`; out of SX-2054 but noted as existing risk.

## Evidence-backed recommendation

- Schema: add new migration to widen `mst.m_survey.answer_frequency` to at least `VARCHAR(50)`; do not mass-update `Multiple`.
- Validation: make create/update accept only three new write values; reject legacy `Multiple` for new write unless product explicitly requires temporary write compatibility.
- Read compatibility: keep list/detail/report/mobile returning raw stored value, including legacy `Multiple`.
- DB hardening: optional transitional `CHECK` should only be added after running distinct-value audit and should include legacy `Multiple`; safest immediate migration is widen-only.
- Downstream logic: document mobile duplicate-submission gap; do not invent business behavior for `One Day` vs `Different Day` without product rule.

## Useful preflight SQL

```sql
SELECT answer_frequency, COUNT(*)
FROM mst.m_survey
GROUP BY 1
ORDER BY 1;
```

## Suggested validation commands

From repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

From `master/`:

```bash
rtk go mod download && rtk go mod tidy
rtk go test ./controller -run 'TestSurveyController_.*AnswerFrequency|TestSurveyController_Create|TestSurveyController_Update' -v
rtk go test ./service -run 'TestSurveyService_.*AnswerFrequency|TestSurveyService_Store|TestSurveyService_Update|TestSurveyService_Detail' -v
rtk go test ./repository -run 'TestSurveyReport' -v
rtk go test ./...
```
