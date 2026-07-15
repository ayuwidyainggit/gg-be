# Plan SX-2253 — Available Stock Same `pro_id` Per Row

Task id: `20260619-0831-sx-2253-available-stock-row`
Jira: SX-2253 (blocked by SX-2250)
Repo target: `sales/` (ScyllaX Sales BE)

## Goal

Hitung `available_stock` per row detail Sales Order secara independen. Row dengan `pro_id` sama tidak boleh saling mempengaruhi nilai `available_stock`. Tetap aman untuk stock movement, on-customer stock, dan promo/reward semantics.

## Non-goals

- Tidak ubah business logic stock movement/on-customer stock.
- Tidak ubah `FindWarehouseStockByWhIdAndProIds` jadi row-keyed.
- Tidak ubah `item_type` semantics, promo flags, atau reward injection logic.
- Tidak hardcode SO number atau product ID tertentu.

## Scope

In scope:

- Field `qty1_stok`, `qty2_stok`, `qty3_stok` di response `OrderResponse.Details`, `DetailsFinal`, `PurchaseDetails`.
- Jalur `DetailV2` di `sales/service/order_service.go` yang menghitung display stock dari warehouse + row qty.
- Unit test di `sales/service/order_service_test.go` yang mengunci behavior row-level stock display untuk duplicate `pro_id`.
- Validasi via `rtk go test ./service -run 'DetailV2'`.

Out of scope:

- Tabel `sls.order_detail.qty*_stok` schema/column rename.
- Perilaku cancel/non-cancel stock display kecuali ditemukan bukti regression.
- Aggregation product-level seperti promo eligibility, total stock summary, atau `aggregatePromoByProduct*`.

## Requirements

1. Tiap `OrderDetResponse` menghitung `qty*_stok` dari `warehouseStockMap[pro_id] + current row qty` saja.
2. Row A dan row B dengan `pro_id` sama, qty berbeda, dan `item_type` berbeda (1 vs 2) wajib menghasilkan `qty*_stok` yang berbeda dan konsisten dengan qty masing-masing.
3. Tabel `sls.order_detail` tidak dimodifikasi, hanya logika mapping ke `OrderDetResponse` di service.
4. Behavior SX-2250 (reward product masuk on-customer stock) tidak diregressi.

## Acceptance Criteria

- `rtk go test ./service -run 'DetailV2'` lulus.
- Test baru row-level same pro_id tersedia dan lulus.
- `TestDetailV2_SameSKURewardDoesNotContaminateNormalRow` tetap lulus.
- `TestDetailV2_Cancelled_UsesWarehouseCurrentOnlyForDisplayedStock` dan `TestDetailV2_NonCancelled_KeepsExistingDisplayedStockBehavior` tetap lulus.
- Manual evidence (optional) menunjukkan row normal/reward same pro_id di staging punya `qty*_stok` independen.

## Existing Patterns / Reuse

- Reuse `computeDisplayedAvailableStockBreakdown` di `sales/service/order_stock_helper.go`. Helper ini sudah row-scoped.
- Reuse `applyStockBreakdownToPointers` untuk assign `Qty1Stok/Qty2Stok/Qty3Stok`.
- Reuse `mockOrderRepositoryDetailV2` di `sales/service/order_service_test.go` untuk test.
- Reuse test naming dan style `TestDetailV2_*` di file yang sama.

## Constraints

- Tenant/parent_cust_id rules dari `ARCHITECTURE.md` tidak dilanggar.
- Layer Controller → Service → Repository → DB tetap dijaga.
- Write operations tetap via service-layer transaction; field ini read-only mapping.
- `rtk` prefix dipakai untuk semua command Go di repo ini.
- Jangan commit kredensial test/staging.

## Risks

- Jika nilai `qty*_stok` di DB sudah tersimpan salah sebelum fix, response setelah fix bisa berbeda dari persisted. Executor harus verifikasi apakah ada call site yang menulis `qty*_stok` (`RefreshOrderDetailStock`) dan apakah perlu re-snap saat fix. Jika ya, fix harus mencakup re-snap row di `DetailV2` untuk row yang broken.
- Map `pro_id` aggregation lain (misal promo eligibility) sengaja tidak diubah.
- Jika FE menghitung ulang dari `qty*_stok`, perubahan nilai di BE akan otomatis ter-reflect.

## Decisions / Assumptions

- Asumsi: defect terjadi karena ada aggregation atau write path yang salah sebelum `DetailV2` mengembalikan response, atau karena helper/fungsi internal menambahkan qty row lain. `DetailV2` saat ini sudah row-scoped pada `computeDisplayedAvailableStockBreakdown`; jika fix ternyata ada di mapping lain, executor akan trace dulu.
- Asumsi: `ItemType=1` adalah ordered product, `ItemType=2` adalah reward product (konsisten dengan MR !95 dan defect prompt).
- Asumsi: helper `computeDisplayedAvailableStockBreakdown` adalah canonical source-of-truth untuk display stock row-level; tidak diganti.
- Asumsi: `useWarehouseCurrentOnly` untuk status CANCELLED tetap dipakai.
- Open question: apakah ada path lain (mis. mobile sync, cancellation re-snap) yang nulis `qty*_stok` dan menyebabkan value stale. Executor harus trace saat implementasi jika ditemukan evidence dari `RefreshOrderDetailStock` callers.

## Execution Source of Truth

Precedence saat executor kerja:

1. Safety/permission/security rules (no secret commit, no test data leak).
2. Acceptance criteria di section ini.
3. TDD/Test Plan section.
4. Implementation Steps dan Execution-ready Worklist.
5. Non-negotiable Implementation Invariants.

Jika ada konflik, executor wajib pilih sumber yang lebih tinggi dan catat di evidence.

## Non-negotiable Implementation Invariants

- `qty*_stok` per row dihitung dari warehouse stock product + qty row itu sendiri saja. Tidak ada operasi `qtyByProduct[pro_id] += ...` yang memaksa row pakai total qty semua row dengan pro_id sama.
- Item row normal (`item_type=1`) dan row reward (`item_type=2`) tetap dipisah; tidak ada merge/collapse.
- `computeDisplayedAvailableStockBreakdown` signature dan semantic tidak diubah; hanya call site atau state input-nya yang dirapikan jika perlu.
- `FindWarehouseStockByWhIdAndProIds` tetap product-level lookup.
- Tidak menyentuh stock movement, on-customer stock, atau `RefreshOrderDetailStock` call sites kecuali evidence menunjukkan itu sumber bug.
- Tidak hardcode SO number atau product ID tertentu di fix.

## Do Not / Reject If

- Reject jika fix menambah `qtyByProduct` map keyed by `pro_id` untuk `qty*_stok` row-level.
- Reject jika fix meng-collapse `item_type=1` dan `item_type=2` row di response.
- Reject jika fix menghapus reward qty dari on-customer stock hanya demi display.
- Reject jika fix menulis hardcode `pro_id` atau `SO2606170005` di source.
- Reject jika fix membuat test dependensi pada token/credential staging.

## Diff Boundary

Boleh diubah:

- `sales/service/order_service.go` (mapping `qty*_stok` row-level di `DetailV2` dan helper terkait).
- `sales/service/order_service_test.go` (test baru row-level same pro_id).
- `sales/service/order_stock_helper.go` jika signature helper perlu di-extract (jarang).
- `sales/service/order_stock_helper_test.go` jika helper internal berubah.

Tidak boleh diubah:

- `sales/repository/order_repository.go` (kecuali ditemukan bug `RefreshOrderDetailStock` yang terjustifikasi; catat di evidence).
- `sales/entity/order_detail.go` (tidak ada perubahan schema/field).
- `sales/model/order_detail.go`.
- Semua service lain (`master`, `inventory`, `finance`, `tms`, dll).
- `docker-compose.yml`, `.env`, config, secrets, migrations di luar Sales.

Bukti wajib:

- Output `rtk go test ./service -run 'DetailV2'` (lampirkan ringkas).
- Output `rtk go test ./service -run 'Stock'`.
- Catatan line:offset perubahan di evidence.

## TDD / Test Plan

- TDD required: yes, karena ini bug di service logic yang punya test mock siap.
- Existing patterns: `TestDetailV2_*` di `order_service_test.go`, `mockOrderRepositoryDetailV2`.
- First failing test: `TestDetailV2_SameProIDNormalAndRewardRow_ComputesStockIndependently`. Asumsikan sebelum fix test gagal, lalu implementasi, lalu hijau.

Test baru yang ditambahkan:

1. `TestDetailV2_SameProIDNormalAndRewardRow_ComputesStockIndependently`
   - `wh_stock = 100` untuk pro_id=X.
   - Row A: pro_id=X, item_type=1, qty1=10.
   - Row B: pro_id=X, item_type=2, qty1=5.
   - Expect: Row A `qty*_stok` mencerminkan `wh_stock + 10`, Row B mencerminkan `wh_stock + 5`. Assert masing-masing `Qty1Stok`/`Qty2Stok`/`Qty3Stok` sesuai canonical mapping.
2. `TestDetailV2_SameProIDTwoNormalRows_ComputesStockIndependently`
   - Row A dan Row B pro_id sama, item_type=1, qty berbeda.
   - Expect: `qty*_stok` independen.
3. `TestDetailV2_DifferentProID_StockByProduct`
   - Row A pro_id=123 wh=100, Row B pro_id=456 wh=50.
   - Expect: row A stock=wh+qty, row B stock=wh+qty masing-masing.
4. `TestDetailV2_SameProIDNormalAndReward_PromoFlagsUnchanged`
   - Assertion tambahan: row normal `IsProductPromotion=false`, row reward `IsProductPromotion=true`, `item_type` masing-masing 1 dan 2.

Commands:

```bash
rtk go test ./service -run 'DetailV2' -count=1
rtk go test ./service -run 'Stock' -count=1
```

## Implementation Steps

1. Trace `qty*_stok` write path: cari semua call site `RefreshOrderDetailStock`. Tentukan apakah ada path yang menulis value row-keyed benar atau salah. Catat hasil trace di evidence.
2. Tulis test `TestDetailV2_SameProIDNormalAndRewardRow_ComputesStockIndependently` di `order_service_test.go`. Jalankan, harus gagal (red).
3. Tulis test `TestDetailV2_SameProIDTwoNormalRows_ComputesStockIndependently` dan `TestDetailV2_DifferentProID_StockByProduct`. Jalankan, harus gagal/red untuk skenario same pro_id duplicate.
4. Perbaiki logika `DetailV2`/`computeDisplayedAvailableStockBreakdown` sehingga `qty*_stok` row-level independen. Hindari map aggregation by `pro_id` untuk field row display.
5. Jalankan semua test `DetailV2_*` dan `TestComputeDisplayedAvailableStockBreakdown*`. Hijau.
6. Jalankan `rtk go test ./...` di `sales/`. Hijau.
7. Verifikasi tidak ada perubahan di luar diff boundary dengan `git diff --stat` (atau setara) sebelum commit.
8. Update evidence dengan output test, line:offset perubahan, dan trace `RefreshOrderDetailStock`.

## Expected Files to Change

- `sales/service/order_service.go`
- `sales/service/order_service_test.go`
- Optional: `sales/service/order_stock_helper.go` (hanya jika ada helper tambahan diekstrak)

## Agent / Tool Routing

- `@fixer` implementasi bounded.
- `@designer` tidak relevan (no UI change).
- `@oracle` review perubahan sebelum quality gate.
- `@quality-gate` final signoff.

## Executor Handoff Prompt

Copy-paste ke executor (mis. `@orchestrator`/`@fixer`):

```
Scope: Fix Jira SX-2253 di service Sales agar available_stock (qty*_stok) per row detail dihitung independen, tidak dijumlah berdasarkan pro_id.

must_preserve:
- item_type semantics (1=normal, 2=reward)
- promo flags dan reward injection
- stock movement / on-customer stock semantics (SX-2250)
- computeDisplayedAvailableStockBreakdown signature & semantic
- FindWarehouseStockByWhIdAndProIds product-level behavior

do_not_touch:
- sales/repository/order_repository.go (kecuali ada bukti RefreshOrderDetailStock bug)
- sales/entity/order_detail.go
- sales/model/order_detail.go
- service lain di luar sales/
- docker-compose.yml, .env, secrets

validation:
- rtk go test ./service -run 'DetailV2' -count=1 (lulus)
- rtk go test ./service -run 'Stock' -count=1 (lulus)
- rtk go test ./service -count=1 (lulus, no regression)

return/evidence:
- ringkasan output test
- list file yang berubah dan line:offset perubahan
- trace singkat call site RefreshOrderDetailStock (jika disentuh)
- tidak ada kredensial staging/test di commit/log/evidence
```

## Execution-ready Worklist / Handoff Contract

Task ids berurutan; tiap task atomic, bisa selesai terpisah. Default owner `@fixer` kecuali ditandai.

### T1 — Trace write path `qty*_stok`

- action: cari semua call site `RefreshOrderDetailStock` dan penulisan `qty1_stok/qty2_stok/qty3_stok`. Catat line:offset.
- depends_on: none
- owner: `@fixer`
- validation: `rg -n 'RefreshOrderDetailStock|qty1_stok|qty2_stok|qty3_stok' sales/`
- exit criteria: daftar call site + line:offset di evidence.
- blocking: ready
- requires_user_decision: no
- must_preserve: write path tidak diubah tanpa justifikasi.
- do_not_touch: `sales/repository/order_repository.go` sampai ada justifikasi.
- evidence_update: tambah ringkasan trace di `.opencode/evidence/<task-id>/discovery.md`.
- exit_verification: `rg` output + ringkasan.
- start_with: T1

### T2 — Tulis test red same pro_id normal+reward

- action: tambah `TestDetailV2_SameProIDNormalAndRewardRow_ComputesStockIndependently` di `order_service_test.go`.
- depends_on: T1
- owner: `@fixer`
- validation: `rtk go test ./service -run 'TestDetailV2_SameProIDNormalAndRewardRow_ComputesStockIndependently' -count=1` (red/fail).
- exit criteria: test ada dan gagal karena bug.
- blocking: ready
- requires_user_decision: no
- must_preserve: gaya test existing.
- do_not_touch: file lain.
- evidence_update: catat output test red.
- exit_verification: test gagal dengan message konsisten.
- start_with: T1

### T3 — Tulis test red two normal rows same pro_id

- action: tambah `TestDetailV2_SameProIDTwoNormalRows_ComputesStockIndependently` di `order_service_test.go`.
- depends_on: T1
- owner: `@fixer`
- validation: `rtk go test ./service -run 'TestDetailV2_SameProIDTwoNormalRows_ComputesStockIndependently' -count=1` (red/fail).
- exit criteria: test ada dan gagal.
- blocking: ready
- requires_user_decision: no
- must_preserve: gaya test existing.
- do_not_touch: file lain.
- evidence_update: catat output test red.
- exit_verification: test gagal konsisten.
- start_with: T1

### T4 — Tulis test green different pro_id (regression)

- action: tambah `TestDetailV2_DifferentProID_StockByProduct` (skenario pro_id berbeda) di `order_service_test.go`. Test ini harus langsung hijau (sanity).
- depends_on: T1
- owner: `@fixer`
- validation: `rtk go test ./service -run 'TestDetailV2_DifferentProID_StockByProduct' -count=1` (green).
- exit criteria: test lulus.
- blocking: ready
- requires_user_decision: no
- must_preserve: gaya test existing.
- do_not_touch: file lain.
- evidence_update: catat output.
- exit_verification: test lulus.
- start_with: T1

### T5 — Implement fix `DetailV2` row-level stock display

- action: perbaiki logika mapping `qty*_stok` di `DetailV2` (sales, purchase, final) dan/atau helper terkait sehingga row-level independen. Jangan pakai map aggregation by pro_id untuk field row display.
- depends_on: T2, T3
- owner: `@fixer`
- validation: `rtk go test ./service -run 'DetailV2' -count=1` (semua hijau).
- exit criteria: T2 dan T3 hijau; test existing `TestDetailV2_SameSKURewardDoesNotContaminateNormalRow`, `TestDetailV2_Cancelled_UsesWarehouseCurrentOnlyForDisplayedStock`, `TestDetailV2_NonCancelled_KeepsExistingDisplayedStockBehavior` tetap hijau.
- blocking: ready
- requires_user_decision: no
- must_preserve: invariant section.
- do_not_touch: list di Non-negotiable.
- evidence_update: catat file:line yang berubah.
- exit_verification: semua `DetailV2_*` hijau.
- start_with: T2

### T6 — Run full service tests

- action: jalankan `rtk go test ./service -count=1`.
- depends_on: T5
- owner: `@fixer`
- validation: exit code 0.
- exit criteria: tidak ada regression di test service lain.
- blocking: ready
- requires_user_decision: no
- must_preserve: invariant.
- do_not_touch: file di luar diff boundary.
- evidence_update: ringkasan output.
- exit_verification: `rtk go test ./service -count=1` exit 0.
- start_with: T5

### T7 — Diff boundary check

- action: `git diff --stat` (atau setara) pastikan perubahan hanya di file yang diizinkan.
- depends_on: T6
- owner: `@fixer`
- validation: tidak ada perubahan di luar `sales/service/*.go` (kecuali test).
- exit criteria: list file berubah hanya dalam diff boundary.
- blocking: ready
- requires_user_decision: no
- must_preserve: invariant.
- do_not_touch: file di luar boundary.
- evidence_update: tempel output diff stat.
- exit_verification: hanya file dalam scope.
- start_with: T6

### T8 — Cleanup evidence dan final summary

- action: hapus `.opencode/draft/20260619-0831-sx-2253-available-stock-row/` kecuali ada open question tersisa. Tulis `Final Planning Summary` di plan ini.
- depends_on: T7
- owner: `@fixer`
- validation: file artifact utama + evidence ringkas.
- exit criteria: ringkasan final tersedia untuk `@quality-gate`.
- blocking: ready
- requires_user_decision: no
- must_preserve: plan ini.
- do_not_touch: source/test.
- evidence_update: update final summary.
- exit_verification: file plan final + evidence lengkap.
- start_with: T7

First action untuk orchestrator: mulai dari T1.

## Validation Commands

```bash
rtk go test ./service -run 'DetailV2' -count=1
rtk go test ./service -run 'Stock' -count=1
rtk go test ./service -count=1
rg -n 'RefreshOrderDetailStock|qty1_stok|qty2_stok|qty3_stok' sales/
```

## Evidence Requirements

- File `.opencode/evidence/20260619-0831-sx-2253-available-stock-row/discovery.md` (sudah ditulis).
- Output ringkas test red/green per task (T2–T6).
- Trace `RefreshOrderDetailStock` (T1).
- Diff stat (T7).
- Manual staging validation optional, hanya jika executor punya env/token lokal; tidak disimpan di repo.

## Done Criteria

- T1–T8 selesai.
- Acceptance Criteria section hijau semua.
- Tidak ada file berubah di luar diff boundary.
- Tidak ada secret/kredensial di commit/log/evidence.

## Final Planning Summary

- Artifacts consulted: `.opencode/docs/ARCHITECTURE.md`, `.opencode/docs/QUALITY.md`, source Sales BE.
- Artifacts created: plan ini, `.opencode/evidence/20260619-0831-sx-2253-available-stock-row/discovery.md`.
- Key decision: fix di `sales/service/order_service.go` row-level, reuse `computeDisplayedAvailableStockBreakdown`, tambah test TDD.
- Assumptions: helper row-level sudah benar, bug kemungkinan di upstream write atau di salah satu call site mapping; executor trace dulu.
- Open questions: apakah `RefreshOrderDetailStock` path menyimpan `qty*_stok` salah sebelum `DetailV2`; executor jawab di T1.
- Readiness: `ready-for-implementation` (executor handoff di worklist).
- Cleanup: draft dir dibuat kosong; akan dihapus di T8 kecuali ada open question.

## Final Execution Summary (post-T8)

- T1-T8 executed. Status: PASS_WITH_RISKS.
- New tests added in `sales/service/order_service_test.go`:
  - `TestDetailV2_SameProIDNormalAndRewardRow_ComputesStockIndependently` (L1326-1511)
  - `TestDetailV2_SameProIDTwoNormalRows_ComputesStockIndependently` (L1513-1606)
  - `TestDetailV2_DifferentProID_StockByProduct` (L1608-1668)
- `DetailV2` mapping in `sales/service/order_service.go:2898-3072` already row-keyed via per-row `computeDisplayedAvailableStockBreakdown(whStockQty, current row Qty1/Qty2/Qty3, ...)`. No code change applied.
- All `TestDetailV2_*` and `TestComputeDisplayedAvailableStockBreakdown*` tests PASS (16/16). Full `rtk go test ./service -count=1` 195/195 PASS.
- Diff boundary: only `sales/service/order_service_test.go` and `.opencode/evidence/.../discovery.md`. No `sales/repository/*`, no `sales/entity/*`, no `sales/model/*` changes.
- Open question remains: where is the user-reported aggregated-display bug actually reproducible? Not in `DetailV2` mapping. Route to `@oracle` for plan-premise review before any further code edit.
