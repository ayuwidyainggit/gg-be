# Execution Plan: fix-global-week-sequence

Status: `PASS_FOR_SLICE`

Deterministic implementation and validation slice is defined. Remote repair remains gated on exact row mapping, backup artifacts, and live preflight because supplied facts omit exact customer, period, and calendar IDs. No mutations executed.

## Objective

Keep declared primary key `(cust_id, per_year, per_id, week_id)`. Repair calendar 39 and future calendar generation so `m_week.per_id` and `m_week.week_id` use one global sequence for each applicable customer/year scope, while `calendar_week_no` remains calendar-local starting at 1.

## Source of truth

- User-approved option 1 and supplied verified remote facts:
  - Active calendar 38: January-June, 25 weeks.
  - Active calendar 39: June-December, 27 weeks.
  - `m_week` duplicate conflicts exist for `per_id/week_id` 1..25.
  - `m_work_day` currently has no duplicates.
- `master/migration/mst.working_day_calendar/20260529_create_working_day_calendar.up.sql`:
  - Adds `working_day_calendar_id` and `calendar_week_no`.
  - Documents `week_id` as compatible global week number.
  - Documents `calendar_week_no` as title-local sequence.
  - Adds calendar/week indexes and foreign keys.
- `master/pkg/generator/working_day_calendar.go`:
  - `FirstWeekID` controls generated `WeekID`.
  - `calendarWeekNo` starts at 1 for every generated calendar.
- `master/service/working_day_calendar_service.go`:
  - `Create` currently passes `FirstWeekID: 1` on every calendar.
  - `ImportHolidays` also regenerates with `FirstWeekID: 1`; this must not recreate conflicting IDs.
  - `materializeWorkingDayCalendarRows` currently writes `PerId` and `WeekId` from generated `WeekID`.
- `master/pkg/generator/mweek.go`:
  - Existing period generator increments `WeekID` across periods; use as compatibility evidence, not as remote repair authority.
- `master/repository/m_week_repository.go`:
  - Existing transaction seam exists, but calendar detail writes use `WorkingDayCalendarRepository`; implementation must preserve transaction ownership and tx-context rules.

## Invariants

1. Do not alter or drop declared PK `(cust_id, per_year, per_id, week_id)`.
2. For each affected customer/year, `(per_id, week_id)` is globally sequenced across calendars; no duplicate PK tuples.
3. `calendar_week_no` remains local: each calendar has values `1..number_of_weeks`.
4. Calendar 39 keeps 27 weeks and date continuity; calendar 38 keeps 25 weeks and original dates.
5. Every affected `m_week` row keeps its calendar association, dates, active/closed flags, audit fields, and other non-key data.
6. Every `m_work_day` row keeps date, work/holiday state, holiday source/note, calendar association, and foreign-key relationship; update references only when required by key repair.
7. `m_work_day` has no duplicate rows before repair; do not create duplicates during repair.
8. All remote repair writes run in one transaction with explicit lock/preflight and rollback on any mismatch.
9. Future create/import paths derive next global week ID from authoritative DB state, not constant `1`.
10. No remote mutation, migration execution, or production deployment occurs during planning.

## Diff boundary

In scope:

- `master/pkg/generator/working_day_calendar.go` and related generator tests.
- `master/service/working_day_calendar_service.go` and related service/repository tests.
- Minimal repository interface/query additions needed to allocate the next global sequence safely inside the calendar write transaction.
- One reviewed remote repair script/runbook or migration artifact, only after exact row mapping and backup plan are approved.
- Validation SQL and evidence capture under `.opencode/evidence/` during execution.

Out of scope:

- PK redesign, constraint removal, or broad legacy week renumbering outside calendars 38/39.
- Changes to `m_work_day` business semantics.
- Changes to calendar-local numbering or UI labels beyond corrected global IDs.
- Remote execution in this planning task.
- Unrelated period, route, sales, finance, or product-ripening logic.

## Worklist

### M0 — Preflight and exact mapping

Owner: `@explorer` / DB operator

Dependencies: none.

Actions:

- Confirm current branch/worktree and target service module.
- Capture schema facts from remote, including exact PK/index definitions, FK dependencies, triggers, and all rows for calendars 38/39.
- Resolve exact `cust_id`, `per_year`, `per_id`, `week_id`, calendar IDs, row counts, date ranges, and downstream references.
- Produce deterministic mapping: old calendar 39 week IDs `1..27` to global IDs after calendar 38. Do not assume offset until query proves calendar 38 sequence and scope.
- Check for references to affected week keys in dependent tables before mutation.

Validation/evidence:

- Preflight SQL output with timestamps and DB identity.
- Schema/constraint dump.
- Affected-row CSV/JSON snapshot.
- Mapping table with old key, new key, date, calendar ID, and dependent-reference counts.

Gate: stop if mapping is not one-to-one, calendar date/order assumptions fail, closed/legacy rows overlap, or external references cannot be updated atomically.

### M1 — Backup and rollback package

Owner: DB operator

Dependencies: M0 exact mapping.

Actions:

- Take logical backup or targeted immutable exports for affected `m_week`, `m_work_day`, calendar rows, and all dependent references.
- Record backup checksum, DB identity, transaction-independent timestamp, retention location, and restore command.
- Prepare rollback SQL from captured old key values. Rollback must restore both `m_week` keys and any updated references in one transaction.
- Dry-run rollback against disposable restore/database when available.

Validation/evidence:

- Backup artifact path and checksum.
- Restore/dry-run log.
- Rollback script review record.

Gate: no remote write without verified backup and tested rollback path.

### M2 — Generator and service patch

Owner: `@fixer` with `@backend` review

Dependencies: M0 confirms sequence scope and transaction boundary.

Actions:

- Replace `FirstWeekID: 1` in calendar create path with a transaction-safe next-global-week allocation for the owning customer/year/calendar sequence.
- Ensure generated rows retain `calendar_week_no = 1..N` while `WeekID` starts after current maximum in the same authoritative scope.
- Apply same allocation rule when `ImportHolidays` regenerates work days; it must reuse existing calendar week IDs, not allocate a second sequence or reset to 1.
- Keep allocation and detail writes in one transaction; lock the sequence source/rows to prevent concurrent calendars receiving overlapping IDs.
- Preserve controller → service → repository → DB layering. Repository writes must honor tx-context extraction.
- Avoid new sequence table unless existing schema cannot provide safe locking; prefer existing calendar/m_week state and minimal query surface.

Validation/evidence:

- Unit tests for first calendar, second calendar, 25+27 calendar sequence, repeated holiday import, concurrent allocation behavior, and invalid input.
- Service/repository transaction tests proving rollback on week/day insert failure.
- `rtk go test ./...` from `master` module, plus focused generator/service tests.
- Diff summary limited to generator, service, repository, and tests.

Gate: no patch merge if new calendar can still emit `week_id=1` in an occupied scope or import can alter existing week IDs.

### M3 — Remote transactional repair

Owner: DB operator; approval: `@orchestrator` and `@quality-gate`

Dependencies: M0, M1, M2 merged and tested.

Actions:

- Enter maintenance/write-quiescence window for calendar creation/import and affected week writes.
- Begin transaction with appropriate isolation and row locks.
- Re-run M0 preflight inside transaction; abort on row-count, checksum, date, ownership, or dependency mismatch.
- Update calendar 39 `m_week` key columns from verified old IDs to globally sequenced IDs. Use collision-safe two-phase temporary IDs or equivalent ordered update because PK is retained and final IDs may overlap existing keys during transition.
- Update `m_work_day` key/reference columns only according to verified schema and mapping. Since current facts show no duplicates, preserve rows and update only necessary key columns/references.
- Update dependent tables discovered in M0 atomically. Do not use broad unscoped updates.
- Re-check PK uniqueness, calendar-local sequence, date continuity, FK integrity, and dependent-reference completeness before commit.
- Commit only after all checks pass; otherwise rollback.

Validation/evidence:

- Transaction log with preflight counts, update counts, post-check results, commit/rollback result.
- Before/after snapshots and checksums.
- No-op rerun result proving idempotent guard rejects already-repaired state or produces zero changes.

### M4 — Post-repair validation and rollout

Owner: `@quality-gate` / release owner

Dependencies: M3 committed repair and deployed M2 patch.

Actions:

- Query calendars 38/39 and all affected dependent rows.
- Exercise create calendar after calendar 39 and verify next global week ID.
- Exercise holiday import on calendar 39 and verify week IDs remain unchanged, days remain one row per date, and holiday state remains correct.
- Verify reads by `(per_year, per_id, week_id, cust_id)`, generated-only filters, calendar views, and downstream consumers.
- Monitor duplicate-key errors, calendar creation/import failures, and week/day count anomalies.
- Keep backup and rollback package until post-release observation window closes.

Release gates:

- Focused and full `master` tests pass.
- Remote post-checks pass with zero duplicate PK tuples and zero orphan references.
- No unexpected row-count delta.
- Quality/security/data signoff complete.
- Rollback command and owner documented before release.

## Critical path and parallel work

Critical path: M0 exact mapping → M1 backup/rollback → M2 patch/tests → M3 transactional repair → M4 post-repair gate.

M2 test design can begin after M0 confirms scope. M1 backup tooling and M2 code review can proceed in parallel, but M3 cannot start until both finish. M4 monitoring preparation can proceed during M3.

## Risks and mitigations

- High: supplied facts omit exact scope and row mapping. Mitigation: M0 live preflight; abort on mismatch.
- High: retained PK causes transient collisions during renumbering. Mitigation: temporary collision-free IDs or safe ordered update inside one transaction; prove with dry run.
- High: hidden downstream references to week keys. Mitigation: dependency inventory before write; update atomically or stop.
- High: concurrent calendar creation can allocate same IDs. Mitigation: transaction lock/serialization around max-ID allocation; concurrency test.
- Medium: import regeneration may reset IDs. Mitigation: import path reuses persisted calendar week mapping; regression test.
- Medium: rollback may be incomplete if dependent tables are missed. Mitigation: targeted backup plus checksum and restore dry run.
- Medium: legacy rows may share customer/year scope. Mitigation: include legacy/generated classification in preflight; never infer from calendar dates alone.
- Medium: remote schema differs from repository migration. Mitigation: live schema dump is authoritative; stop on drift.

## Handoff contracts

- `@explorer` returns exact mapping, schema, dependency inventory, and preflight evidence.
- DB operator returns backup checksum, restore test, and reviewed rollback SQL.
- `@fixer` returns bounded patch, focused tests, full `master` test result, and diff boundary.
- `@quality-gate` returns data-integrity, security, and release signoff after M3/M4 evidence.
- Release owner returns observation-window metrics and rollback decision.

## Determinism decision

`PASS_FOR_SLICE`: implementation scope, invariants, test matrix, transaction strategy, and release gates are deterministic. Full remote execution is not deterministic yet because exact remote row mapping and schema/dependency evidence were not included in current inputs. Planning task performed no mutations.
