# Plan SX-2085 — Monitoring Activity Collection Paid

Task ID: `20260528-1440-sx-2085-monitoring-collection-paid`
Jira: `SX-2085` — `[Defect] [BE] Collection doesnt show at view detail Monitoring Activity`
Module: `pjp` / Live Monitoring Detail
Endpoint: `GET /scylla-pjp/api/v1/monitoring_locations/details`
Source of truth: file ini.
Evidence kept: `.opencode/evidence/20260528-1440-sx-2085-monitoring-collection-paid/discovery.md`.
Open questions kept: `.opencode/draft/20260528-1440-sx-2085-monitoring-collection-paid/open-questions.md`.

## Goal

Perbaiki BE Monitoring Activity detail agar section `collection` menampilkan paid collection berdasarkan `selected_date` dan `emp_id` salesman yang dipilih. Amount harus berasal dari `acf.deposit_payment.payment_amount`, difilter oleh `acf.deposit.deposit_date = request.date`, `acf.deposit.emp_id = request.emp_id`, dan `acf.deposit.collection_no IS NOT NULL`.

## Non-goals

- Tidak mengubah route endpoint detail.
- Tidak mengubah auth/JWT behaviour.
- Tidak mengubah skema DB.
- Tidak memakai `sls.order.total` sebagai paid amount.
- Tidak menambah hardcoded date `2026-05-28` atau `emp_id=421`.
- Tidak mengubah kontrak response jadi object `collection.total_paid` tanpa koordinasi FE.
- Tidak mengubah list endpoint `live-monitoring-principal` / `live-monitoring-distributor` kecuali compile impact interface.

## Scope

- Tambah query collection paid di `pjp` live monitoring detail repository.
- Tambah row model untuk result collection per outlet.
- Isi response existing `collection: []CollectionData` dengan `outlet_id`, `outlet_code`, `outlet_name`, `collection_total`.
- Update `collection_summary` agar count/status mengikuti jumlah item collection.
- Tambah unit/regression tests untuk service detail.
- Jalankan DB/API evidence untuk sample QA `date=2026-05-28`, `emp_id=421` bila staging/token tersedia.

## Requirements

1. Endpoint detail tetap menerima `date` string `YYYY-MM-DD` dari query request.
2. Endpoint detail tetap menerima `emp_id` int dari query request.
3. Query collection wajib parameterized via GORM placeholders.
4. Query collection wajib filter `d.deposit_date = req.Date` atau `DATE(d.deposit_date) = req.Date` sesuai tipe kolom aktual; prefer exact date kalau kolom `date`.
5. Query collection wajib filter `d.emp_id = req.EmpID`.
6. Query collection wajib filter `d.collection_no IS NOT NULL`.
7. Amount wajib `SUM(COALESCE(acf.deposit_payment.payment_amount,0))`.
8. No data wajib return `collection: []` dan `collection_summary.count=0,status=none`, bukan error.
9. Multiple payment row per deposit wajib dijumlah.
10. Multiple detail row per deposit tidak boleh duplicate payment amount.
11. Tenant filter wajib memakai target salesman `cust_id` dari `GetSalesmanCustID` / `targetCustIDs`.
12. Existing sections `visit_information`, `sales`, `return`, `expense`, `shipment` tidak boleh berubah.

## Acceptance Criteria

- [ ] `GET /scylla-pjp/api/v1/monitoring_locations/details?emp_id=421&date=2026-05-28` mengembalikan `collection` berisi paid collection bila DB punya data.
- [ ] `collection[].collection_total` sama dengan sum `acf.deposit_payment.payment_amount` untuk date + emp_id + `collection_no IS NOT NULL`.
- [ ] `sls.order.total` tidak dipakai untuk paid collection amount.
- [ ] Jika tidak ada collection, response tetap sukses dengan `collection: []`.
- [ ] Multiple `deposit_payment` rows terjumlah.
- [ ] Multiple `deposit_detail` rows tidak menduplikasi amount.
- [ ] `collection_summary.count` sama dengan jumlah outlet collection.
- [ ] Regression detail fields lain tetap sama.
- [ ] Staging QA bisa melihat Collection section untuk salesman/date valid.

## Existing Patterns/Reuse

- Reuse endpoint detail existing: `GET /monitoring_locations/details`.
- Reuse request existing: `LiveMonitoringDetailRequest{EmpID, DistributorID, Date}`.
- Reuse response existing: `response.CollectionData` dengan JSON `collection_total`.
- Reuse service flow existing di `GetMonitoringDetail`: resolve `salesmanCustID`, build `targetCustIDs`, lalu load transactional sections.
- Reuse repository pattern `GetSales`, `GetReturns`, `GetShipments`.
- Reuse `buildVisitSummary(len(collection))`.
- Tidak ada util KiloCode/project lain yang sudah menghitung Monitoring Activity collection paid di `pjp`; perlu tambah repository method baru.

## Constraints

- Repo multi-module Go; target module `pjp`.
- Layer wajib Controller → Service → Repository → DB.
- Repository hanya data access; business mapping tetap service.
- Tenant/scope jangan dilonggarkan; semua transactional query harus jaga `cust_id`.
- Shell workflow repo wajib pakai `rtk`.
- Planner tidak mengedit source; implementasi dilakukan oleh `@orchestrator`/`@fixer` setelah plan ini.
- Staging credentials/token tidak boleh disimpan di artifact.

## Risks

- Query referensi Jira berisiko duplicate amount karena join `deposit_payment` ke `deposit_detail` pada grain invoice/detail.
- Satu `deposit_no,cust_id` bisa terkait banyak outlet; perlu DB validation. Bila lintas outlet terjadi, business harus menentukan alokasi payment.
- Existing FE mungkin berharap `collection` array, bukan object. Rekomendasi: isi array existing dulu.
- `deposit_date` tipe aktual belum divalidasi; implementer wajib pilih comparison yang tidak kena timezone mismatch.
- Staging DB/token mungkin tidak tersedia; QA evidence bisa blocked walau code/test lokal selesai.

## Decisions/Assumptions

Decisions:

- Pertahankan contract existing `collection: []CollectionData`.
- Gunakan `collection_total` untuk paid amount per outlet.
- Jangan tambah top-level `collection.total_paid` di bugfix awal kecuali FE meminta.
- Gunakan query agregasi aman: aggregate payment per deposit dulu, lalu join ke distinct deposit-outlet.
- Gunakan `targetCustIDs` dari salesman cust_id untuk filter `acf.deposit` dan related transactional joins.

Assumptions / Open Questions:

- A1: `req.Date` sudah `YYYY-MM-DD` dan sesuai `acf.deposit.deposit_date`.
- A2: `req.EmpID` adalah `acf.deposit.emp_id` / salesman emp_id.
- A3: FE bisa render amount dari existing `collection[].collection_total`.
- Q1: FE butuh top-level total `collection.total_paid` atau cukup array existing? Default: array existing.
- Q2: Jika satu deposit lintas outlet, payment harus split atau blocked by business? Default: validate; bila terjadi, minta keputusan.
- Q3: Staging token/DB access tersedia untuk evidence? Default: implementer jalankan bila akses ada.

Question gate: tidak perlu tanya sebelum primary plan karena default non-breaking jelas. Open questions dicatat di draft dan tidak memblokir implementasi awal kecuali lintas-outlet grain terbukti terjadi.

## TDD/Test Plan

TDD wajib karena bug logic/query dan response behaviour.

Existing test patterns:

- `pjp/service/live_monitoring/get_detail_service_test.go` punya `detailRepoStub` dan tests untuk detail sections.
- Repository SQL tests belum tampak; fokus awal service tests + optional repository SQL/dry-run bila pattern tersedia.

Red step:

1. Tambah `model.CollectionRow` dan repository interface method `GetCollections`; test awal compile fail sampai stub/source disesuaikan.
2. Tambah `TestGetMonitoringDetail_IncludesCollectionPaid`:
   - stub `GetCollections` return satu row `OutletID=101`, `CollectionTotal=1500000`.
   - expect `len(result.Collection)==1`, `collection_total==1500000`, `collection_summary.count==1,status=completed`.
3. Tambah `TestGetMonitoringDetail_NoCollectionReturnsEmptyList`:
   - stub collection kosong.
   - expect `collection: []`, `collection_summary.none`.
4. Tambah `TestGetMonitoringDetail_CollectionUsesRequestDateAndEmpID`:
   - stub capture `date`, `empID`, `custIDs`.
   - expect captured `2026-05-28`, `421`, `[]string{salesmanCustID}`.
5. Jika repository test feasible, tambah SQL-level test/mock untuk memastikan query memakai `deposit_payment.payment_amount`, bukan `sls.order.total`.

Green step:

- Implement `GetCollections` repository method.
- Map rows ke `response.CollectionData` di service.
- Set `collectionSummary := buildVisitSummary(len(collection))`.
- Update stub tests.

Refactor step:

- Ekstrak transform helper bila mapping mulai panjang, misal `transformCollectionRows`.
- Keep naming selaras existing `GetSales`, `GetReturns`.
- Jangan refactor unrelated SX-2038 logic.

Edge cases:

- No data → empty slice.
- `deposit_payment.payment_amount` NULL → 0.
- Multiple payment rows → sum per deposit.
- Multiple deposit rows same outlet → sum per outlet.
- Multiple detail rows same deposit → no duplicate payment.
- Date mismatch timezone → prefer date column exact compare or explicit `DATE()` only bila kolom timestamp.

Commands:

```bash
rtk go test ./service/live_monitoring -run 'TestGetMonitoringDetail_.*Collection'
rtk go test ./repository/live_monitoring ./service/live_monitoring
rtk go test ./...
```

## Implementation Steps

1. Run DB grain validation for `date='2026-05-28'`, `emp_id=421`:

```sql
SELECT
    d.deposit_no,
    d.cust_id,
    COUNT(DISTINCT dd.invoice_no) AS invoice_count,
    COUNT(dp.*) AS payment_join_count,
    SUM(dp.payment_amount) AS joined_payment_sum
FROM acf.deposit d
JOIN acf.deposit_detail dd
    ON dd.deposit_no = d.deposit_no
   AND dd.cust_id = d.cust_id
LEFT JOIN acf.deposit_payment dp
    ON dp.deposit_no = d.deposit_no
   AND dp.cust_id = d.cust_id
WHERE d.deposit_date = '2026-05-28'
  AND d.collection_no IS NOT NULL
  AND d.emp_id = 421
GROUP BY d.deposit_no, d.cust_id;
```

2. Add model row in `pjp/model/live_monitoring.go`:

```go
type CollectionRow struct {
    OutletID         int     `gorm:"column:outlet_id"`
    OutletCode       string  `gorm:"column:outlet_code"`
    OutletName       string  `gorm:"column:outlet_name"`
    CollectionTotal  float64 `gorm:"column:collection_total"`
}
```

3. Add interface method in `pjp/repository/live_monitoring/live_monitoring_repository.go`:

```go
GetCollections(ctx context.Context, tx *gorm.DB, custIDs []string, date string, empID int) ([]model.CollectionRow, error)
```

4. Implement repository method in `pjp/repository/live_monitoring/get_detail_repository.go` after `GetReturns` or near other detail queries. Use safe aggregation:

```sql
WITH payment_per_deposit AS (
    SELECT
        deposit_no,
        cust_id,
        SUM(COALESCE(payment_amount, 0)) AS payment_amount
    FROM acf.deposit_payment
    GROUP BY deposit_no, cust_id
), deposit_outlet AS (
    SELECT DISTINCT
        d.emp_id,
        o.outlet_id,
        d.deposit_no,
        d.cust_id,
        d.deposit_date,
        d.collection_no
    FROM acf.deposit d
    JOIN acf.deposit_detail dd
        ON dd.deposit_no = d.deposit_no
       AND dd.cust_id = d.cust_id
    JOIN sls."order" o
        ON o.invoice_no = dd.invoice_no
       AND o.cust_id = dd.cust_id
    WHERE d.cust_id IN ?
      AND d.deposit_date = ?
      AND d.collection_no IS NOT NULL
      AND d.emp_id = ?
)
SELECT
    mo.outlet_id,
    mo.outlet_code,
    mo.outlet_name,
    SUM(COALESCE(ppd.payment_amount, 0)) AS collection_total
FROM deposit_outlet do2
LEFT JOIN payment_per_deposit ppd
    ON ppd.deposit_no = do2.deposit_no
   AND ppd.cust_id = do2.cust_id
JOIN mst.m_outlet mo
    ON mo.outlet_id = do2.outlet_id
   AND mo.cust_id = do2.cust_id
GROUP BY mo.outlet_id, mo.outlet_code, mo.outlet_name
ORDER BY mo.outlet_code;
```

GORM option: use `tx.Raw(query, custIDs, date, empID).Scan(&results)` for CTE clarity, or `Table` with subquery if repo avoids raw. Raw is acceptable if parameterized.

5. Update `GetMonitoringDetail` service:
   - call `collectionRows, err := s.repository.GetCollections(ctx, s.db, targetCustIDs, req.Date, req.EmpID)` after returns or before expenses.
   - map to `[]response.CollectionData` using pointer fields.
   - replace hardcoded empty collection.
   - keep `collectionSummary := buildVisitSummary(len(collection))`.

6. Update `detailRepoStub` in tests:
   - add `collections []model.CollectionRow`, `receivedCollectionCustIDs []string`, `receivedCollectionDate string`, `receivedCollectionEmpID int`.
   - implement `GetCollections`.

7. Add collection service tests listed in TDD Red step.

8. Run targeted tests, then full module tests.

9. Run staging SQL manual and API before/after if access exists. Capture only sanitized request/response snippets, no token.

10. Send to `@quality-gate` for final signoff.

## Expected Files to Change

- `pjp/model/live_monitoring.go` — add `CollectionRow`.
- `pjp/repository/live_monitoring/live_monitoring_repository.go` — add `GetCollections` method.
- `pjp/repository/live_monitoring/get_detail_repository.go` — add collection paid query.
- `pjp/service/live_monitoring/get_detail_service.go` — fetch/map collection and summary.
- `pjp/service/live_monitoring/get_detail_service_test.go` — update stub and add collection regression tests.
- Optional: `pjp/repository/live_monitoring/get_detail_repository_test.go` — if SQL-level test pattern added without new heavy dependency.

## Agent/Tool Routing

- `@orchestrator`: start execution from this plan, route implementation, integrate evidence.
- `@fixer`: implement code, tests, Red → Green → Refactor.
- `@oracle`: review only if DB grain shows cross-outlet deposit allocation ambiguity.
- `@quality-gate`: final review for tenant scope, query safety, response compatibility, QA evidence.
- `@librarian`: not needed; no external/library behaviour.
- `@architect`: not needed unless response contract changes to object/top-level totals.

## Execution-ready Worklist / Handoff Contract

start_with: `T1`

```yaml
- id: T1
  action: Jalankan DB validation untuk sample SX-2085 date 2026-05-28 emp_id 421, termasuk duplicate amount risk query dan manual expected result.
  depends_on: none
  owner: @fixer
  validation: hasil SQL ringkas tersimpan sebagai evidence implementasi; jika DB inaccessible, catat blocker dan lanjut local TDD.
  exit: grain data terkonfirmasi atau blocker DB terdokumentasi.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T2
  action: Tambah failing collection tests di service detail stub untuk paid collection, no data, dan passthrough date/emp_id/custIDs.
  depends_on: T1
  owner: @fixer
  validation: rtk go test ./service/live_monitoring -run 'TestGetMonitoringDetail_.*Collection' gagal karena method/logic belum ada.
  exit: Red tests mengekspresikan SX-2085.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T3
  action: Tambah model CollectionRow dan repository interface GetCollections.
  depends_on: T2
  owner: @fixer
  validation: compile bergerak ke missing implementation/stub yang jelas.
  exit: kontrak repo siap untuk implementation.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T4
  action: Implement GetCollections dengan CTE parameterized dan aggregation safe per deposit sebelum join outlet.
  depends_on: T3
  owner: @fixer
  validation: rtk go test ./repository/live_monitoring ./service/live_monitoring
  exit: repository compile, no raw string interpolation, no use sls.order.total as collection amount.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T5
  action: Update GetMonitoringDetail service untuk fetch collection, map response.CollectionData, dan update collection_summary.
  depends_on: T4
  owner: @fixer
  validation: rtk go test ./service/live_monitoring -run 'TestGetMonitoringDetail_.*Collection'
  exit: collection tests hijau.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T6
  action: Refactor kecil bila perlu, pastikan existing detail fields dan tests SX-2038 tetap hijau.
  depends_on: T5
  owner: @fixer
  validation: rtk go test ./service/live_monitoring -run TestGetMonitoringDetail
  exit: no regression detail behaviour.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T7
  action: Jalankan full test module pjp.
  depends_on: T6
  owner: @fixer
  validation: rtk go test ./...
  exit: semua test hijau atau unrelated failure terdokumentasi.
  status: ready
  blocker: none
  requires_user_decision: no

- id: T8
  action: Smoke API staging/local untuk request detail emp_id=421 date=2026-05-28; bandingkan response collection dengan SQL manual.
  depends_on: T7
  owner: @fixer
  validation: response snippet menunjukkan collection amount sama dengan SQL; token tidak disimpan.
  exit: QA evidence siap atau blocker akses tertulis.
  status: ready
  blocker: butuh token staging/DB access valid
  requires_user_decision: no

- id: T9
  action: Jika DB validation membuktikan satu deposit lintas outlet dan amount allocation ambigu, stop dan minta keputusan business.
  depends_on: T1
  owner: @oracle
  validation: contoh deposit lintas outlet disertakan dengan candidate allocation options.
  exit: keputusan business tertulis sebelum merge.
  status: blocked
  blocker: hanya aktif bila lintas-outlet grain terbukti.
  requires_user_decision: yes

- id: T10
  action: Review final dengan @quality-gate untuk query safety, tenant scope, API contract, tests, dan QA evidence.
  depends_on: T8
  owner: @quality-gate
  validation: quality-gate signoff atau actionable blockers.
  exit: siap PR/deploy.
  status: ready
  blocker: none
  requires_user_decision: no
```

## Validation Commands

Dari repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

Dari `pjp/`:

```bash
rtk go test ./service/live_monitoring -run 'TestGetMonitoringDetail_.*Collection'
rtk go test ./service/live_monitoring -run TestGetMonitoringDetail
rtk go test ./repository/live_monitoring ./service/live_monitoring
rtk go test ./...
```

Manual API smoke, token jangan disimpan:

```bash
curl -G 'https://best.scyllax.online/scylla-pjp/api/v1/monitoring_locations/details' \
  --data-urlencode 'emp_id=421' \
  --data-urlencode 'date=2026-05-28' \
  -H 'Authorization: Bearer <token>' \
  -H 'Accept: application/json, text/plain, */*'
```

Manual SQL expected:

```sql
WITH payment_per_deposit AS (
    SELECT deposit_no, cust_id, SUM(COALESCE(payment_amount, 0)) AS payment_amount
    FROM acf.deposit_payment
    GROUP BY deposit_no, cust_id
), deposit_outlet AS (
    SELECT DISTINCT d.emp_id, o.outlet_id, d.deposit_no, d.cust_id, d.deposit_date, d.collection_no
    FROM acf.deposit d
    JOIN acf.deposit_detail dd ON dd.deposit_no = d.deposit_no AND dd.cust_id = d.cust_id
    JOIN sls."order" o ON o.invoice_no = dd.invoice_no AND o.cust_id = dd.cust_id
    WHERE d.deposit_date = '2026-05-28'
      AND d.collection_no IS NOT NULL
      AND d.emp_id = 421
)
SELECT mo.outlet_code, mo.outlet_name, SUM(COALESCE(ppd.payment_amount, 0)) AS collection_total
FROM deposit_outlet do2
LEFT JOIN payment_per_deposit ppd ON ppd.deposit_no = do2.deposit_no AND ppd.cust_id = do2.cust_id
JOIN mst.m_outlet mo ON mo.outlet_id = do2.outlet_id AND mo.cust_id = do2.cust_id
GROUP BY mo.outlet_code, mo.outlet_name
ORDER BY mo.outlet_code;
```

## Evidence Requirements

- DB validation result untuk `2026-05-28`, `emp_id=421`.
- Duplicate amount risk query result.
- Before response showing `collection: []` or missing amount.
- After response snippet showing `collection[].collection_total`.
- Test output targeted and `rtk go test ./...`.
- Diff note proving `payment_amount` used, not `order.total`.
- If DB/token unavailable, record exact blocker.
- Research gate:
  - Local discovery: required, done.
  - Official docs/context7: not needed, no library/API behaviour new.
  - GitHub: not needed, bug internal repo.
  - Brave/web: not needed, no external current facts.
  - Browser/screenshot: not needed, BE-only; QA may validate UI manually after deploy.

## Done Criteria

- Worklist `T1`–`T8` complete or access blockers documented.
- `collection` array populated from `acf.deposit_payment.payment_amount`.
- No duplicate payment due to `deposit_detail` join.
- No hardcoded QA sample values in code.
- Unit/regression tests added and green.
- Full `pjp` tests green or unrelated failures documented.
- Staging response matches SQL manual result.
- `@quality-gate` signoff.

## Final Planning Summary

- Primary plan path: `.opencode/plans/20260528-1440-sx-2085-monitoring-collection-paid.md` — source of truth implementation.
- Evidence created/kept: `.opencode/evidence/20260528-1440-sx-2085-monitoring-collection-paid/discovery.md` and `index.json`; kept because implementation needs inspected files, existing contract, root cause, and query grain risk.
- Draft kept: `.opencode/draft/20260528-1440-sx-2085-monitoring-collection-paid/open-questions.md`; kept because FE top-level total and cross-outlet allocation remain conditional questions.
- Key decisions: fill existing `collection` array; use `collection_total`; use safe payment-per-deposit CTE; keep response non-breaking.
- Assumptions: `req.Date` maps to `deposit_date`; `req.EmpID` maps to `deposit.emp_id`; FE can sum array if total needed.
- Questions asked: not asked in chat; non-blocking defaults chosen. Blocking only if cross-outlet payment allocation ambiguity appears in DB.
- Cleanup: no draft/evidence deleted because both remain operationally useful for implementer/QA.
- Readiness: ready for `@orchestrator` start at `T1`; `T8` needs staging token/DB access.
