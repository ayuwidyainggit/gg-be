# Verification — SX-2054 Answer Frequency

Task ID: `20260525-1738-sx-2054-answer-frequency`

## Changed files

- `master/entity/survey.go`
- `master/pkg/validation/validation.go`
- `master/pkg/constant/survey_answer_frequency.go`
- `master/migration/mst.survey/004_alter_answer_frequency_length.sql`
- `master/controller/survey_controller_test.go`
- `master/service/survey_service_test.go`

## Root cause confirmed

1. `master/migration/mst.survey/001_create_tables.sql` used `answer_frequency VARCHAR(20) NOT NULL`.
   - too short for:
     - `Multiple Times, One Day`
     - `Multiple Times, Different Day`
2. `master/entity/survey.go` used validator tag `oneof=Multiple 'One Time'` for create/update.
3. Read paths already passed raw DB value through service/repository layers.

## Implemented changes

### 1. Schema / persistence

Added migration:

- `master/migration/mst.survey/004_alter_answer_frequency_length.sql`

SQL:

```sql
ALTER TABLE mst.m_survey
ALTER COLUMN answer_frequency TYPE VARCHAR(50);
```

### 2. Centralized contract values

Added source of truth:

- `master/pkg/constant/survey_answer_frequency.go`

Supported new write values:

- `One Time`
- `Multiple Times, One Day`
- `Multiple Times, Different Day`

Legacy read-only compatibility value kept:

- `Multiple`

### 3. API validation

Updated request tags in:

- `master/entity/survey.go`

From:

- `validate:"required,oneof=Multiple 'One Time'"`

To:

- `validate:"required,answer_frequency"`

Registered custom validator and translation in:

- `master/pkg/validation/validation.go`

Behavior now:

- create/update accepts exact 3 new values
- create/update rejects legacy `Multiple`
- create/update rejects unknown typo/case mismatch

### 4. Read compatibility

No normalization added.

Read paths remain raw pass-through in:

- `master/repository/survey_repository.go`
- `master/service/survey_service.go`
- `master/repository/survey_report_repository.go`
- `mobile/repository/survey.go`

This preserves legacy row readability for `Multiple`.

## Tests run

From `master/`:

```bash
rtk go test ./controller -run 'TestSurveyController_.*AnswerFrequency|TestSurveyController_Create|TestSurveyController_Update' -v
```

Result:
- 12 passed in 1 package

```bash
rtk go test ./service -run 'TestSurveyService_.*AnswerFrequency|TestSurveyService_Store|TestSurveyService_Update|TestSurveyService_Detail|TestSurveyService_List' -v
```

Result:
- 30 passed in 1 package

```bash
rtk go test ./...
```

Result:
- 301 passed in 23 packages

## Test coverage added/updated

### Controller

- create accepts `Multiple Times, One Day`
- create accepts `Multiple Times, Different Day`
- create rejects legacy `Multiple`
- update accepts `Multiple Times, One Day`
- update accepts `Multiple Times, Different Day`
- update rejects legacy `Multiple`
- capture stub asserts exact value forwarded to service

### Service

- list keeps raw legacy `Multiple`
- detail keeps raw legacy `Multiple`

## Grep audit summary

Post-change grep in `master` confirmed:

- validator now uses custom tag `answer_frequency`
- constants/helper centralized in `master/pkg/constant/survey_answer_frequency.go`
- read paths still pass `answer_frequency` raw
- literal `Multiple` remains only in legacy-read tests and compatibility-related cases, not as valid write contract

## Runtime / database verification status

- `rtk docker compose -f docker-compose.yml ps` ran from repo root
- current local result: no compose services running locally during this verification
- Target DB validated directly with `psql` using provided connection
- Migration `master/migration/mst.survey/004_alter_answer_frequency_length.sql` executed successfully on target DB

Preflight before migration:

- column type: `character varying:20`
- existing values:
  - `Multiple` = 22
  - `One Time` = 51

Post-migration validation:

- column type: `character varying:50`
- existing values preserved:
  - `Multiple` = 22
  - `One Time` = 51

Transactional DB smoke checks executed and rolled back:

1. insert survey header with `answer_frequency = 'Multiple Times, One Day'`
   - inserted row length = 23
2. update same transactional row to `Multiple Times, Different Day`
   - updated row length = 29
3. update same transactional row to `One Time`
   - updated row length = 8
4. rollback confirmed temporary row removed (`count = 0`)
5. legacy row compatibility check:
   - transactional update from legacy `Multiple` to `Multiple Times, One Day` succeeded
   - rollback preserved original counts

API smoke create/update/read against running service was not executed in this run; DB-level migration and transactional validation were executed directly.

## Legacy compatibility note

Safe strategy implemented:

- new write only accepts new values
- legacy `Multiple` still readable on master report/list/detail and mobile read paths
- no arbitrary migration from `Multiple` to one of new values

## Known remaining risks

1. FE or other clients still writing `Multiple` will now get validation error by design.
2. Mobile duplicate-submission semantics are unchanged and still out-of-scope here:
   - `mobile/service/survey.go` still blocks repeat submission for all surveys via `CheckExistingSubmission(...)`
   - no rule difference for `One Day` vs `Different Day` implemented without explicit product requirement.
3. Direct API smoke against running service was not executed in this run; validation completed at unit-test and DB-transaction level.
