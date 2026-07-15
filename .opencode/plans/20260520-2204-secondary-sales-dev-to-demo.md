# Plan — Restore Secondary Sales Report endpoints dari `dev` ke branch demo berbasis `qa`

Task id: `20260520-2204-secondary-sales-dev-to-demo`
Tanggal: `2026-05-20 22:04 Asia/Jakarta`
Service target: `sales`
Source branch: `dev` (`8a8a0e6`, sudah termasuk merge `bugfix/secondary-sales-dev` menurut user dan `git log`)
Target base branch: `qa` (`e16c0a1`)
Target demo branch: `demo-20052026-2204`
Target worktree: `/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260520-2204/sales`
Primary source of truth: `.opencode/plans/20260520-2204-secondary-sales-dev-to-demo.md`
Evidence: `.opencode/evidence/20260520-2204-secondary-sales-dev-to-demo/discovery.md`

## Goal

Buat branch demo baru dari `qa`, lalu restore/copy kode endpoint Secondary Sales Report yang ada di `dev` ke branch demo tersebut, berdasarkan endpoint yang tercantum di `docs/Secondary Sales Report_BE.md`.

Endpoint target docs:

1. `GET /sales/v1/reports/secondary-sales/trend-sales?year=2026`
2. `GET /sales/v1/reports/secondary-sales/sum-date?month=5`
3. `GET /sales/v1/reports/secondary-sales/group?month=5&group_by=outlet`
4. `POST /sales/v1/reports/secondary-sales`

Endpoint pendamping yang ikut karena berbagi pipeline/helper Secondary Sales Report:

- `POST /sales/v1/extract/secondary-sales`

## Non-goals

- Tidak merge penuh `dev` ke `qa`/demo.
- Tidak cherry-pick commit lintas scope.
- Tidak push branch remote kecuali diminta setelah implementasi.
- Tidak buka MR/PR.
- Tidak edit source saat fase plan.
- Tidak mengubah secret, `.env`, compose, atau config deploy.
- Tidak membawa endpoint Activity Report Sales kecuali perubahan compile langsung tidak bisa dipisah.

## Scope

File scope yang berbeda antara `qa..dev` dan harus dicopy dari `dev` ke branch demo:

- `controller/report_controller.go`
- `controller/so_controller_test.go`
- `entity/report.go`
- `repository/report_repository.go`
- `repository/report_repository_test.go`
- `service/report_service.go`
- `service/report_service_test.go`

File yang dicek tetapi tidak berubah dalam scope `qa..dev`:

- `model/report.go` — tidak perlu copy.
- `pkg/config/env/env.go` — tidak perlu copy.
- `pkg/constant/constant.go` — tidak perlu copy.
- `main.go` — tidak perlu copy.

## Requirements

- Branch demo baru dibuat dari `origin/qa`.
- Nama branch mengikuti user: `demo-DDMMYYYY-HHII` → `demo-20052026-2204`.
- Mode sinkronisasi: copy file scope endpoint dari `dev`, bukan merge/cherry-pick.
- Kode di branch demo harus membawa fitur/fix Secondary Sales Report dari `dev`, termasuk:
  - Trend Sales filter `cust_id` body + `year` query.
  - Sum Date dan Group filter `year` + `cust_id`.
  - Export filter `cust_id` body.
  - Scope check `cust_id` untuk principal/distributor via `resolveSecondaryDashboardCustID` + `ExistsCustomerInParentScope`.
  - Export owner `report.list.cust_id = auth user`.
  - Export net sales formula tidak lagi bergantung ke `amount_final`.
  - Legacy `SecondarySales()` net sales formula juga tidak bergantung ke `amount_final`.
  - Regression tests terkait dari `dev` ikut terbawa.
- Validasi minimal dari worktree demo:
  - `rtk go test ./... -count=1`
  - targeted tests untuk report controller/service/repository.
- Tidak commit secret.

## Acceptance Criteria

1. Branch lokal `demo-20052026-2204` ada dan berbasis `origin/qa`.
2. Worktree baru ada di `/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260520-2204/sales`.
3. Diff `demo-20052026-2204` terhadap `origin/qa` hanya memuat file scope di atas, kecuali compile membutuhkan tambahan eksplisit.
4. Semua endpoint docs ada di route demo:
   - `POST /v1/reports/secondary-sales`
   - `GET /v1/reports/secondary-sales/sum-date`
   - `GET /v1/reports/secondary-sales/group`
   - `GET /v1/reports/secondary-sales/trend-sales`
5. Behavior `cust_id` sesuai dev terbaru:
   - empty/equal auth → auth cust.
   - principal bisa request child active.
   - distributor request sibling → 403.
   - invalid `cust_id` → 400.
6. Export formula untuk `NetSalesExcPPN` dan `NetSalesIncPPN` memakai qty*price formula, bukan `amount_final`.
7. `rtk go test ./... -count=1` lulus di worktree demo.
8. Commit lokal dibuat di branch demo setelah validasi lulus.
9. Jika ada conflict saat copy, implementor berhenti dan lapor path + conflict reason.

## Existing Patterns/Reuse

- Reuse plan sebelumnya: `.opencode/plans/20260518-2026-secondary-sales-dev-to-qa.md` untuk pola worktree dan scope branch demo.
- Reuse route existing di `ReportController.Route`.
- Reuse strict layer Controller → Service → Repository → DB.
- Reuse helper service `resolveSecondaryDashboardCustID` dari `dev`.
- Reuse repository SQL builder `buildSecondarySalesUnionQuery` dan helper `buildReportSecondarySalesReportOrderQuery` dari `dev`.
- Reuse test patterns dari `sales/service/report_service_test.go`, `sales/controller/so_controller_test.go`, dan `sales/repository/report_repository_test.go`.

## Constraints

- Shell command di repo ini harus pakai `rtk` untuk test/go command.
- Jangan push remote tanpa instruksi eksplisit.
- Jangan force-push.
- Jangan mengubah branch/worktree existing seperti `demo-18052026`.
- Branch `demo-20052026-2204` belum ada saat discovery.
- Source root `/Users/ujang/Projects/Geekgarden/scylla-be/sales` sedang di branch `dev`; jangan melakukan edit langsung di situ untuk demo implementation.
- Worktree target harus dibuat terpisah.

## Risks

- Copy file utuh dari `dev` bisa membawa perubahan test/helper yang bergantung pada file di luar scope. Mitigasi: jalankan compile/test penuh; jika gagal, tambahkan dependency compile langsung dengan alasan tertulis.
- Endpoint docs menyebut GET request body untuk `cust_id`; beberapa proxy/client bisa drop body GET. Behavior tetap aman karena fallback auth.
- Branch demo berbasis `qa` mungkin memiliki perbedaan config/runtime dari `dev`. Mitigasi: jangan copy config kecuali compile wajib.
- Worktree path bisa sudah ada. Mitigasi: cek path dulu; kalau ada, minta keputusan user, jangan hapus otomatis.

## Decisions/Assumptions

Keputusan user:

- Source code dari `dev` terbaru; `bugfix/secondary-sales-dev` sudah merge ke `dev`.
- Mode sync: copy file scope endpoint dari `dev` ke branch demo.
- Nama branch: `demo-DDMMYYYY-HHII`, berbasis `qa`.

Asumsi:

- `HHII` = jam-menit lokal Asia/Jakarta. Untuk sesi ini: `2204`, jadi `demo-20052026-2204`.
- Worktree path mengikuti pola repo: `/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260520-2204/sales`.
- `origin/dev` dan `origin/qa` setelah `git fetch` adalah source of truth. Lokal `dev` sudah sama dengan `origin/dev` pada `8a8a0e6`.

Open questions: tidak ada.

## TDD/Test Plan

TDD required: ya, untuk validasi restore behavior lintas branch.

Reason:

- Restore endpoint report lintas branch berisiko mematahkan tenant isolation, formula export, dan handler contracts.

Existing tests yang wajib ikut dan dijalankan di branch demo:

- Controller:
  - `TestSecondarySalesExportReturnsForbiddenForDistributorSiblingCust`
  - `TestSecondarySalesExportAuthCustNotOverwrittenByBody`
  - `TestSecondarySalesExportInvalidCustIDReturns400`
  - `TestSecondaryReportSalesTrendSalesReturnsForbiddenForDistributorSibling`
  - `TestSecondaryReportSalesTrendSalesPassesCustIDFromBody`
  - `TestSecondaryReportSalesTrendSalesInvalidCustIDReturns400`
- Service:
  - `TestPublishSecondarySalesReportUsesEffectiveCustButStoresAuthOwner`
  - `TestPublishSecondarySalesReportRejectsUnauthorizedDistributorSibling`
  - `TestSecondarySalesReportTrendSalesUsesChildCustWhenAllowed`
- Repository:
  - `TestBuildSecondarySalesUnionQueryNetSalesFromFormula`
  - `TestGetReportSecondarySalesReportOrderNetSalesFromFormula`
  - `TestSecondarySalesLegacyQueryNetSalesFromFormula`

First failing/regression test expectation:

- Setelah copy file scope dari `dev`, semua test di atas harus tersedia dan lulus.
- Jika test tidak ada di branch demo setelah copy, copy scope belum lengkap.

Commands dari worktree demo:

```bash
rtk go test ./controller -run 'TestSecondarySales|TestSecondaryReportSales' -count=1
rtk go test ./service -run 'TestPublishSecondarySalesReport|TestSecondarySalesReport' -count=1
rtk go test ./repository -run 'TestBuildSecondarySalesUnionQuery|TestGetReportSecondarySalesReportOrder|TestSecondarySalesLegacyQuery' -count=1
rtk go test ./... -count=1
```

## Implementation Steps

1. Fetch refs:

```bash
git fetch --all --prune
```

2. Verify target branch and worktree path absent:

```bash
git branch --list demo-20052026-2204
test ! -e /Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260520-2204/sales
```

3. Create parent dir if needed, then worktree from `origin/qa`:

```bash
mkdir -p /Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260520-2204
git worktree add -b demo-20052026-2204 /Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260520-2204/sales origin/qa
```

4. In worktree demo, copy scoped files from `dev`:

```bash
git checkout dev -- \
  controller/report_controller.go \
  controller/so_controller_test.go \
  entity/report.go \
  repository/report_repository.go \
  repository/report_repository_test.go \
  service/report_service.go \
  service/report_service_test.go
```

5. Inspect diff:

```bash
git status --short
git diff --stat origin/qa
git diff --name-only origin/qa
```

Expected diff files only:

```text
controller/report_controller.go
controller/so_controller_test.go
entity/report.go
repository/report_repository.go
repository/report_repository_test.go
service/report_service.go
service/report_service_test.go
```

6. Validate routes and endpoint symbols exist (read/grep only):

```bash
grep -n 'secondary-sales' controller/report_controller.go
grep -n 'SecondarySalesReportTrendSales\|SecondarySalesReportSumReportByMonth\|SecondarySalesReportGroupSales\|PublishSecondarySalesReport' service/report_service.go
```

7. Run tests:

```bash
rtk go test ./controller -run 'TestSecondarySales|TestSecondaryReportSales' -count=1
rtk go test ./service -run 'TestPublishSecondarySalesReport|TestSecondarySalesReport' -count=1
rtk go test ./repository -run 'TestBuildSecondarySalesUnionQuery|TestGetReportSecondarySalesReportOrder|TestSecondarySalesLegacyQuery' -count=1
rtk go test ./... -count=1
```

8. If tests pass, commit locally on branch demo:

```bash
git add controller/report_controller.go controller/so_controller_test.go entity/report.go repository/report_repository.go repository/report_repository_test.go service/report_service.go service/report_service_test.go
git commit -m "restore secondary sales report endpoints from dev"
```

9. Final status:

```bash
git status --short
git log --oneline -3
```

10. Do not push unless user explicitly asks.

## Expected Files to Change

In target worktree branch `demo-20052026-2204`:

- `controller/report_controller.go`
- `controller/so_controller_test.go`
- `entity/report.go`
- `repository/report_repository.go`
- `repository/report_repository_test.go`
- `service/report_service.go`
- `service/report_service_test.go`

No expected change:

- `model/report.go`
- `main.go`
- `pkg/config/env/env.go`
- `pkg/constant/constant.go`
- `.env`, compose, lockfiles.

## Agent/Tool Routing

- Execution: `@fixer` or direct orchestrator bounded git/worktree operation.
- Discovery: `@explorer` if file scope expands due compile failures.
- Final review: `@quality-gate` after tests pass, because this is branch-restore + tenant-sensitive endpoint.

## Execution-ready Worklist / Handoff Contract

`start_with`: `R01`

| id | depends_on | action | owner | validation | exit criteria | blocking | requires_user_decision |
|---|---|---|---|---|---|---|---|
| R01 | none | `git fetch --all --prune` di repo `sales` source | @fixer | `git rev-parse --short origin/dev origin/qa` | refs tersedia | ready | no |
| R02 | R01 | cek branch `demo-20052026-2204` dan worktree path belum ada | @fixer | `git branch --list demo-20052026-2204`; `test ! -e <path>` | tidak ada konflik nama/path | ready | no |
| R03 | R02 | buat worktree branch `demo-20052026-2204` dari `origin/qa` | @fixer | `git -C <worktree> branch --show-current`; `git rev-parse --short HEAD` | branch demo aktif, base = `origin/qa` | ready | no |
| R04 | R03 | copy 7 file scope endpoint dari `dev` ke worktree demo | @fixer | `git status --short`; `git diff --name-only origin/qa` | hanya file scope berubah | ready | no |
| R05 | R04 | cek endpoint docs dan route symbol di worktree demo | @fixer | grep routes + service methods | 4 endpoint docs + extract route tersedia | ready | no |
| R06 | R04 | jalankan targeted tests controller/service/repository | @fixer | commands di TDD plan | targeted tests hijau | ready | no |
| R07 | R06 | jalankan full suite `rtk go test ./... -count=1` | @fixer | full command | full suite hijau | ready | no |
| R08 | R07 | quality gate final | @quality-gate | review diff + test output | PASS/PASS_WITH_RISKS tanpa blocker | ready | no |
| R09 | R08 | commit lokal di branch demo | @fixer | `git log --oneline -1`; `git status --short` | commit ada, working tree clean | ready | no |
| R10 | R09 | laporkan branch, worktree path, commit hash, tests | @orchestrator | final summary | user dapat push/manual deploy | ready | no |

Blocked policy:

- Jika branch/path target sudah ada → stop dan minta keputusan user.
- Jika tests gagal karena dependency file di luar scope → catat file, alasan, dan minta persetujuan untuk perluas scope.
- Jika conflict terjadi saat checkout file dari `dev` → stop dan lapor conflict path.

## Validation Commands

Dari target worktree:

```bash
rtk go test ./controller -run 'TestSecondarySales|TestSecondaryReportSales' -count=1
rtk go test ./service -run 'TestPublishSecondarySalesReport|TestSecondarySalesReport' -count=1
rtk go test ./repository -run 'TestBuildSecondarySalesUnionQuery|TestGetReportSecondarySalesReportOrder|TestSecondarySalesLegacyQuery' -count=1
rtk go test ./... -count=1
```

Optional read-only DB checks setelah service demo deploy:

```bash
# Trend expected child data
PGPASSWORD='<redacted>' psql -h 103.28.219.73 -p 25431 -U postgres -d scylla_citus_dev -At -c "select m.month, coalesce(sum(fo.gross_sale),0)::bigint, coalesce(sum(fo.discount + fo.special_discount),0)::bigint, coalesce(sum(fo.net_sales_exclude_ppn),0)::bigint from (select generate_series(1,12) as month) m left join report.dim_dates dt on dt.month=m.month and dt.year=2026 left join report.fact_orders fo on fo.date_id=dt.id and fo.cust_id='C260020001' group by m.month order by m.month"

# Export formula sample for INV2605200015
PGPASSWORD='<redacted>' psql -h 103.28.219.73 -p 25431 -U postgres -d scylla_citus_dev -At -c "select ((coalesce(od.qty1_final,0)*coalesce(od.sell_price1,0)) + (coalesce(od.qty2_final,0)*coalesce(od.sell_price2,0)) + (coalesce(od.qty3_final,0)*coalesce(od.sell_price3,0))) - coalesce(od.promo_value_final,0) - coalesce(od.disc_value_final,0) as net_exc, (((coalesce(od.qty1_final,0)*coalesce(od.sell_price1,0)) + (coalesce(od.qty2_final,0)*coalesce(od.sell_price2,0)) + (coalesce(od.qty3_final,0)*coalesce(od.sell_price3,0))) - coalesce(od.promo_value_final,0) - coalesce(od.disc_value_final,0)) + coalesce(od.vat_value_final,0) as net_inc from sls.order_detail od join sls.\"order\" o on o.ro_no=od.ro_no and o.cust_id=od.cust_id where o.invoice_no='INV2605200015' and o.cust_id='C260020001'"
```

## Evidence Requirements

Keep under `.opencode/evidence/20260520-2204-secondary-sales-dev-to-demo/`:

- `discovery.md` (created).
- `implementation-log.md` (to create during execution): branch/worktree creation, copied files, commit hash.
- `validation.md` (to create during execution): targeted tests + full test output.
- Optional `db-validation.md` if read-only DB checks run.

## Done Criteria

- Branch `demo-20052026-2204` exists locally and is based on `origin/qa`.
- Worktree path exists and points to branch demo.
- 7 scoped files copied from `dev`.
- Diff against `origin/qa` limited to scope.
- Tests pass.
- Quality gate completed.
- Local commit exists on demo branch.
- No push performed unless separately requested.

## Final Planning Summary

- Artifacts created:
  - `.opencode/plans/20260520-2204-secondary-sales-dev-to-demo.md`
  - `.opencode/evidence/20260520-2204-secondary-sales-dev-to-demo/discovery.md`
- Questions asked and answered:
  - Source branch: `dev` terbaru, sudah termasuk merge bugfix.
  - Sync mode: copy file scope endpoint.
  - Branch naming: single branch baru dari `qa`, pattern `demo-DDMMYYYY-HHII`.
- Key decisions:
  - Target branch `demo-20052026-2204`.
  - Target worktree `/Users/ujang/Projects/Geekgarden/scylla-be-worktrees-20260520-2204/sales`.
  - Copy only 7 changed endpoint/test files from `dev`.
- Readiness: ready for implementation from `R01`.
- Cleanup: no draft artifacts created. Evidence kept because it contains branch/file scope and endpoint mapping.
