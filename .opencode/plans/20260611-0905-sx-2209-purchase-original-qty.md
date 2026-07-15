# SX-2209 — Purchase Details Tetap Tampil Saat Available Stock = 0

Task ID: `20260611-0905-sx-2209-purchase-original-qty`  
Readiness: `ready-for-implementation`  
Plan Quality Gate: `PASS`  
Mode: Maintenance Stability Mode  
Primary source of truth: file ini.

## Goal

Perbaiki response `GET /sales/v2/orders/{order_no}` supaya `data.purchase_details.normal[]` tetap menampilkan row Purchase Order ketika current `qty_po1/2/3` semuanya `0`, selama salah satu `original_qty_po1/2/3` non-zero.

Target sample Jira: `SO2606100013`, `order_detail_id = 7273`, `pro_id = 10813`, `original_qty_po3 = 3`, `qty_po1/2/3 = 0` harus muncul di `purchase_details.normal[]`.

## Non-goals

- Tidak mengubah create order `order_type = O`.
- Tidak mengubah stock validation, stock mutation, atau process order.
- Tidak mengganti nilai `qty_po*` dengan `original_qty_po*`.
- Tidak mengubah formula promo, discount, tax, gross, net sales.
- Tidak mengubah tab `details.normal[]` atau `details_final.normal[]` kecuali test menunjukkan dampak langsung dari fix.
- Tidak menyimpan token Jira/API/cURL di source, test snapshot, atau artifact.

## Scope

Masuk scope:

- `sales/controller/order_controller.go` hanya untuk memastikan route/flow; perubahan tidak diharapkan.
- `sales/service/order_service.go`, terutama filter pembentukan `response.PurchaseDetails.Normal` di `DetailV2`.
- `sales/service/order_service_test.go` untuk regression test SX-2209.
- Verifikasi model/entity yang sudah expose `original_qty_po*`.

Di luar scope kecuali discovery implementasi membuktikan perlu:

- Repository query `FindDetail`, karena saat ini sudah `Select("sls.order_detail.* ...")`.
- Migration baru, karena `sales/migration/sls.order/add_order_type_and_original_qty_po_fields.sql` sudah menambahkan `original_qty_po1/2/3`.

## Requirements

- Row masuk `purchase_details.normal[]` jika:

```text
any(qty_po1, qty_po2, qty_po3) non-zero
OR any(original_qty_po1, original_qty_po2, original_qty_po3) non-zero
```

- Null harus diperlakukan sebagai `0`.
- Current `qty_po*` tetap merepresentasikan qty yang bisa diproses.
- `original_qty_po*` tetap merepresentasikan original Taking Order/Purchase Order.
- Promo/reward injection tidak boleh panic saat row purchase punya current qty zero.
- Field JSON `original_qty_po1`, `original_qty_po2`, `original_qty_po3` harus tetap tersedia di `purchase_details.normal[]`; evidence menunjukkan field sudah ada di `entity.OrderDetResponse`.

## Acceptance Criteria

- `GET /sales/v2/orders/SO2606100013` mengembalikan `purchase_details.normal[]` yang berisi `order_detail_id = 7273`.
- Row dengan `qty_po1/2/3 = 0` tetap tampil jika salah satu `original_qty_po1/2/3` non-zero.
- Row existing dengan `qty_po* > 0` tetap tampil.
- Row dengan semua `qty_po*` dan semua `original_qty_po*` zero/null tetap tidak tampil di Purchase Order jika tidak aktif lewat fallback sales existing.
- Response null-safe untuk row lama dengan `original_qty_po* IS NULL`.
- `qty_po*` tidak diubah menjadi original qty.
- Test regression mencakup original-only purchase detail.
- Tidak ada secret/token yang masuk source atau test.

## Existing Patterns/Reuse

Evidence utama ada di `.opencode/evidence/20260611-0905-sx-2209-purchase-original-qty/discovery.md`.

Reuse yang harus dipakai:

- `getValueOrDefault` di `sales/service/order_service.go` untuk null-safe pointer float.
- `mockOrderRepositoryDetailV2` di `sales/service/order_service_test.go`.
- Test existing `TestDetailV2_PurchaseDetailsUsesPurchaseActiveRowsForOrderTypeO` sebagai pola regression test.
- Model/entity existing:
  - `model.OrderDetailRead.OriginalQtyPo1/2/3`
  - `entity.OrderDetResponse.OriginalQtyPo1/2/3`
- Route existing:
  - `OrderController.Route()` → `GET /v2/orders/:ro_no` → `DetailV2`.

Bug candidate repo-backed:

```go
func activeQtyForTab(detail model.OrderDetailRead, tab promoSnapshotTab) float64 {
    switch tab {
    case promoSnapshotTabPurchase:
        return getValueOrDefault(detail.QtyPo1, 0) + getValueOrDefault(detail.QtyPo2, 0) + getValueOrDefault(detail.QtyPo3, 0)
    }
}

func isActiveDetailForTab(detail model.OrderDetailRead, tab promoSnapshotTab) bool {
    if detail.ItemType == 2 { return false }
    return activeQtyForTab(detail, tab) > 0
}
```

`DetailV2` saat ini append purchase normal jika `isActiveDetailForTab(detail, promoSnapshotTabPurchase) || isActiveDetailForTab(detail, promoSnapshotTabSalesOrder)`.

## Constraints

- Ikuti aturan repo: jalankan command dengan prefix `rtk`.
- Validasi di module `sales`, bukan root.
- Preserve layer Controller → Service → Repository → DB.
- Query harus tetap tenant-safe dengan `cust_id`; tidak perlu ubah query bila tidak perlu.
- Tool write saat planning melaporkan LSP errors existing pada file lain (`report_repository.go`, `order_controller.go`, `order_service.go`, `order_status_helper.go`, `order_controller_test.go`). Executor harus membedakan pre-existing diagnostics dari perubahan SX-2209 dan mencatatnya jika validasi global gagal.

## Risks

- Mengubah `activeQtyForTab` langsung dapat memengaruhi promo snapshot, promo consult, recompute promo, dan helper lain. Risiko ini lebih besar dari kebutuhan display response.
- Row original-only dengan current qty zero dapat masuk promo consult payload; kalkulasi harus tetap berbasis `qty_po* = 0`, bukan original qty.
- Jika executor menambahkan fallback ke `Qty1/2/3`, tab Purchase Order bisa memunculkan Sales Order biasa secara tidak sengaja.
- Jika test global sudah gagal karena error existing, klaim selesai harus berdasarkan targeted test plus bukti pre-existing failure.

## Decisions/Assumptions

- Keputusan: fix utama sebaiknya service-layer display predicate khusus untuk `purchase_details`, bukan repository SQL, karena repository sudah fetch semua detail.
- Keputusan: jangan ubah `activeQtyForTab` kecuali executor menemukan semua pengguna helper memang menginginkan original qty purchase. Berdasarkan discovery, helper dipakai luas sehingga lebih aman tidak disentuh.
- Keputusan: operator `> 0` mengikuti convention existing `isActiveDetailForTab`; jika domain perlu retur/negative qty, executor boleh pakai `!= 0` hanya setelah menemukan convention repo yang mendukung.
- Asumsi slice-safe: field `original_qty_po*` sudah ada di database target karena migration SX-2184 tersedia.
- Question gate: tidak ada pertanyaan blocking; prompt user sudah memberikan rule bisnis, sample, dan non-goals.

## Execution Source of Truth

Urutan precedence executor:

1. Instruksi eksplisit terbaru dari user.
2. Aturan safety/security repo: jangan commit/copy secret, token, `.env`; gunakan `rtk`; validasi di module `sales`.
3. Non-negotiable Implementation Invariants di plan ini.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria dan Done Criteria.
6. Implementation Steps.
7. Rekomendasi tambahan dari executor selama tetap dalam Diff Boundary.

Jika ada konflik, ikuti sumber yang lebih tinggi dan catat konflik di evidence final.

## Non-negotiable Implementation Invariants

- `qty_po*` tidak boleh diisi dari `original_qty_po*` untuk response atau DB.
- `original_qty_po*` hanya boleh dipakai untuk keputusan display row Purchase Order.
- Kalkulasi promo, VAT, amount, final, dan stock harus tetap berbasis current qty existing.
- Scope perubahan utama adalah `DetailV2` response `PurchaseDetails.Normal`.
- Jangan mengubah mutation/process order/stock untuk menyelesaikan bug display.
- Jangan mengubah `details.normal[]` atau `details_final.normal[]` semantics.
- Jangan menambahkan token/cURL Authorization di test atau artifact.
- Jika validation global gagal karena error existing, wajib catat exact failure dan targeted test result.

## Do Not / Reject If

Tolak atau revert implementasi jika:

- Mengubah `qty_po3` dari `0` menjadi `3` pada sample hanya agar FE tampil.
- Menggunakan `original_qty_po*` untuk stock mutation, process order, promo calculation, VAT, gross/net sales.
- Membuat filter purchase details menjadi semua product item tanpa PO/original PO.
- Menghapus fallback existing `isActiveDetailForTab(detail, promoSnapshotTabSalesOrder)` tanpa alasan dan regression evidence.
- Menambah migration duplicate untuk `original_qty_po*` tanpa bukti DB schema missing.
- Mengubah file di service lain tanpa bukti kebutuhan langsung.
- Menyimpan credential/token dari Jira atau curl.

## Diff Boundary

Allowed source/test changes:

- `sales/service/order_service.go`
- `sales/service/order_service_test.go`

Allowed only if evidence shows necessary:

- `sales/entity/order_detail.go` jika ternyata response build tidak expose `original_qty_po*` pada branch aktual executor.
- `sales/model/order_detail.go` jika branch aktual belum punya field `OriginalQtyPo*`.
- `sales/repository/order_repository.go` jika executor menemukan query branch aktual tidak memilih `sls.order_detail.*` atau tidak mengambil `original_qty_po*`.

Evidence/output paths:

- `.opencode/evidence/20260611-0905-sx-2209-purchase-original-qty/`
- `.opencode/plans/20260611-0905-sx-2209-purchase-original-qty.md`

Any out-of-boundary change must be reverted or justified in verification evidence before final quality gate.

## TDD/Test Plan

TDD required: yes. Ini production behavior bug dan regression risk ada di response filter.

Existing test pattern:

- `sales/service/order_service_test.go`
- `TestDetailV2_PurchaseDetailsUsesPurchaseActiveRowsForOrderTypeO`
- `mockOrderRepositoryDetailV2`

Red step:

- Tambah test baru sebelum fix, misalnya:

```go
func TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero(t *testing.T)
```

Scenario:

- `OrderType = "O"`
- Detail normal item `ItemType = 1`
- `OrderDetailID = 7273`
- `ProId = 10813`
- `QtyPo1/2/3 = 0`
- `OriginalQtyPo1 = 0`, `OriginalQtyPo2 = 0`, `OriginalQtyPo3 = 3`
- Sales qty dan final qty nil/zero agar row tidak lolos lewat fallback tab lain.
- Expected `len(response.PurchaseDetails.Normal) == 1`, row ID/pro ID match, `QtyPo3 == 0`, `OriginalQtyPo3 == 3`, `Qty3` tetap nil atau zero sesuai behavior existing untuk current purchase qty zero.

Green step:

- Tambah helper predicate display purchase yang null-safe, misalnya:

```go
func hasPurchaseDisplayQty(detail model.OrderDetailRead) bool {
    return getValueOrDefault(detail.QtyPo1, 0) > 0 ||
        getValueOrDefault(detail.QtyPo2, 0) > 0 ||
        getValueOrDefault(detail.QtyPo3, 0) > 0 ||
        getValueOrDefault(detail.OriginalQtyPo1, 0) > 0 ||
        getValueOrDefault(detail.OriginalQtyPo2, 0) > 0 ||
        getValueOrDefault(detail.OriginalQtyPo3, 0) > 0
}
```

- Gunakan helper di append `PurchaseDetails.Normal`:

```go
if hasPurchaseDisplayQty(detail) || isActiveDetailForTab(detail, promoSnapshotTabSalesOrder) {
    response.PurchaseDetails.Normal = append(response.PurchaseDetails.Normal, detailData)
}
```

- Preserve mapping `Qty1/2/3` dari `QtyPo1/2/3` hanya saat `purchaseQtyPoTotal > 0` kecuali test existing meminta lain.

Refactor step:

- Jika helper juga perlu `ItemType` guard, buat explicit:

```go
func shouldIncludePurchaseDetailRow(detail model.OrderDetailRead) bool {
    if detail.ItemType == 2 { return false }
    return hasPurchaseDisplayQty(detail) || isActiveDetailForTab(detail, promoSnapshotTabSalesOrder)
}
```

- Hindari duplicate boolean expression di loop.

Edge cases:

- `qty_po* > 0`, original nil: tetap tampil.
- semua current/original nil atau zero dan sales fallback inactive: tidak tampil.
- original nil pada row lama: tidak panic.
- promo repository nil dan non-nil path tidak panic.

Commands:

```bash
rtk go test ./service -run 'TestDetailV2_PurchaseDetails(UsesPurchaseActiveRowsForOrderTypeO|IncludesOriginalQtyWhenCurrentQtyZero)'
rtk go test ./service -run 'TestDetailV2'
rtk go test ./...
```

Jika `rtk go test ./...` gagal karena pre-existing compile/LSP issues, catat exact failure dan jalankan targeted package/test yang relevan setelah memastikan perubahan SX-2209 lulus.

## Implementation Steps

1. Dari repo root, cek runtime/status sesuai harness jika perlu:

```bash
rtk docker compose -f docker-compose.yml ps
```

2. Masuk validasi module `sales` via `workdir=/Users/ujang/Projects/Geekgarden/scylla-be/sales` untuk command test.
3. Tambah regression test original-only purchase detail di `sales/service/order_service_test.go`.
4. Jalankan targeted test dan pastikan gagal karena row belum tampil.
5. Tambah helper predicate di `sales/service/order_service.go` dekat `activeQtyForTab` / `isActiveDetailForTab`.
6. Ubah kondisi append `response.PurchaseDetails.Normal` di `DetailV2` agar memakai helper display purchase.
7. Pastikan `QtyPo*` dan `OriginalQtyPo*` tetap ikut Automapper dan tidak dimodifikasi.
8. Jalankan targeted tests, lalu broader `DetailV2`, lalu `rtk go test ./...` jika feasible.
9. Jika punya local DB/token aman, smoke manual:

```bash
curl -s "${BASE_URL}/sales/v2/orders/SO2606100013" \
  -H "Authorization: Bearer ${TOKEN}" \
  | jq '.data.purchase_details.normal[] | select(.order_detail_id == 7273)'
```

10. Catat hasil validasi dan limitations untuk `@quality-gate`.

## Expected Files to Change

- `sales/service/order_service.go`
- `sales/service/order_service_test.go`

Tidak diharapkan berubah:

- migration files
- DB schema docs
- stock mutation helpers
- controller route
- repository query

## Agent/Tool Routing

- `@orchestrator`: mulai eksekusi dari plan ini, koordinasi evidence dan handoff.
- `@fixer`: implementasi bounded source/test di module `sales`.
- `@quality-gate`: final signoff karena perubahan material pada backend response behavior.
- `@explorer`: optional hanya jika executor menemukan branch code berbeda dari discovery ini.
- `@librarian`: tidak diperlukan kecuali muncul dependency/library behavior baru.

## Executor Handoff Prompt

Copyable prompt untuk `@orchestrator` / implementation lane:

```text
Implementasikan SX-2209 berdasarkan .opencode/plans/20260611-0905-sx-2209-purchase-original-qty.md sebagai source of truth.

Scope: service sales, endpoint GET /sales/v2/orders/{order_no}, response data.purchase_details.normal[]. Row Purchase Order harus tampil jika qty_po1/2/3 non-zero ATAU original_qty_po1/2/3 non-zero. Current qty_po* tetap 0 jika memang 0; jangan overwrite dengan original qty.

Must preserve: no stock mutation/process order changes, no promo/VAT/amount formula changes, no details/details_final semantic changes, no secrets/tokens in code/tests, use rtk commands, validate inside sales module.

Do not touch outside allowed diff boundary unless evidence proves necessary: sales/service/order_service.go and sales/service/order_service_test.go. Any out-of-boundary change must be justified or reverted.

Start with failing regression test for original_qty_po3 > 0 with qty_po1/2/3 = 0. Then implement minimal predicate for purchase_details display. Validate with targeted go tests and broader go test if feasible. Return changed files, test results, any pre-existing failures, and manual smoke status if performed.
```

## Execution-ready Worklist / Handoff Contract

`start_with`: `SX2209-1`

| Task | Action | depends_on | owner/lane | validation | exit criteria | status | requires_user_decision | must_preserve | do_not_touch | evidence_update | exit_verification |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| SX2209-1 | Baca plan dan discovery, konfirmasi branch code masih sesuai evidence | none | `@orchestrator`/`@explorer` | Read plan + grep `activeQtyForTab`, `DetailV2` purchase loop | Flow dan bug candidate terkonfirmasi atau discrepancy dicatat | ready | no | artifact source of truth | source edits | Catat file/line aktual | Notes sebelum edit |
| SX2209-2 | Tambah regression test original-only purchase detail | SX2209-1 | `@fixer` | `rtk go test ./service -run TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero` | Test gagal dengan alasan row tidak muncul | ready | no | test harus menjaga `qty_po*` tetap zero | production logic sebelum Red | Simpan output Red | Failing test due expected bug |
| SX2209-3 | Implement minimal display predicate untuk purchase details | SX2209-2 | `@fixer` | targeted test dari SX2209-2 | Test baru lulus | ready | no | jangan ubah stock/promo formula; jangan overwrite qty_po | repository/migration kecuali perlu | Catat diff helper dan loop | Targeted test Green |
| SX2209-4 | Regression test existing purchase qty dan negative zero/null case | SX2209-3 | `@fixer` | `rtk go test ./service -run 'TestDetailV2_PurchaseDetails'` | Existing test tetap lulus; negative case tercover atau alasan reuse jelas | ready | no | tab sales/final tidak berubah | controller/DB | Catat output | Test package targeted lulus |
| SX2209-5 | Broader validation dan optional manual smoke | SX2209-4 | `@fixer`/`@orchestrator` | `rtk go test ./service -run 'TestDetailV2'`; `rtk go test ./...`; optional curl with env token | Hasil validasi jelas; pre-existing failures dipisahkan | ready | no | no secrets in logs/code | `.env`, token, unrelated services | Catat commands dan failures | Evidence siap quality gate |
| SX2209-6 | Final review/signoff | SX2209-5 | `@quality-gate` | Review diff + evidence | Pass atau remediation list | ready | no | all invariants | implementation outside boundary | Signoff evidence | Final gate complete |

## Validation Commands

Dari `workdir=/Users/ujang/Projects/Geekgarden/scylla-be/sales`:

```bash
rtk go test ./service -run 'TestDetailV2_PurchaseDetailsIncludesOriginalQtyWhenCurrentQtyZero'
rtk go test ./service -run 'TestDetailV2_PurchaseDetails'
rtk go test ./service -run 'TestDetailV2'
rtk go test ./...
```

Optional DB/API verification dengan token dari env lokal saja:

```bash
curl -s "${BASE_URL}/sales/v2/orders/SO2606100013" \
  -H "Authorization: Bearer ${TOKEN}" \
  | jq '.data.purchase_details.normal[] | select(.order_detail_id == 7273)'
```

Manual SQL optional:

```sql
SELECT od.order_detail_id, od.pro_id,
       od.original_qty_po1, od.original_qty_po2, od.original_qty_po3,
       od.qty_po1, od.qty_po2, od.qty_po3
FROM sls.order_detail od
WHERE od.order_no = 'SO2606100013'
   OR od.ro_no = 'SO2606100013'
ORDER BY od.order_detail_id;
```

## Evidence Requirements

Executor harus menyimpan atau melaporkan:

- Red test output untuk regression test baru.
- Green test output setelah fix.
- Output targeted `TestDetailV2` atau alasan tidak bisa menjalankan.
- Output `rtk go test ./...` atau exact failure dan apakah failure pre-existing.
- Jika manual smoke dilakukan, payload sanitized untuk row `7273`; jangan simpan token.
- Jika ada perubahan di luar expected files, justification dan diff evidence.

Source strategy yang dipakai plan:

- Repo-local evidence dan prompt Jira user dipakai.
- Atlassian/web dilewati karena prompt sudah memuat detail issue dan akses Jira kemungkinan credential-gated.
- Official/library docs dilewati karena tidak ada API library baru atau behavior version-sensitive.
- Runtime/browser dilewati pada planning; runtime verification direncanakan untuk executor.

## Done Criteria

- Test baru membuktikan original-only purchase row tampil.
- Existing `TestDetailV2_PurchaseDetailsUsesPurchaseActiveRowsForOrderTypeO` tetap lulus.
- Negative zero/null case tidak memunculkan empty purchase row kecuali fallback sales existing berlaku.
- `qty_po*` tetap current qty dan `original_qty_po*` tetap original qty di response.
- Tidak ada perubahan stock mutation/process order.
- Evidence validasi lengkap dan siap `@quality-gate`.

## Final Planning Summary

Artifacts dibuat:

- Primary plan: `.opencode/plans/20260611-0905-sx-2209-purchase-original-qty.md`
- Evidence kept: `.opencode/evidence/20260611-0905-sx-2209-purchase-original-qty/discovery.md`
- Evidence manifest kept: `.opencode/evidence/20260611-0905-sx-2209-purchase-original-qty/index.json`

Evidence kept karena operasional untuk executor: berisi file/line discovery, bug candidate, constraints, dan test pattern.

Key decisions:

- Fix di service-layer display predicate untuk `purchase_details.normal[]`.
- Jangan ubah repository/migration karena field sudah tersedia dan query memilih `sls.order_detail.*`.
- Jangan ubah `activeQtyForTab` global kecuali executor menemukan bukti kuat, karena helper dipakai luas untuk promo snapshot/consult.

Assumptions:

- Migration `original_qty_po*` sudah diterapkan di environment target.
- Operator `> 0` cukup mengikuti convention existing; jika domain negative qty valid, executor harus menyesuaikan dengan evidence repo.

Open questions:

- Tidak ada yang blocking. Jika manual smoke ke `best.scyllax.online` diperlukan, executor/user harus menyediakan `BASE_URL` dan `TOKEN` via env lokal, tidak di artifact.

Cleanup:

- Tidak ada draft dibuat.
- Evidence tidak dihapus karena tetap berguna untuk implementasi dan quality gate.

Readiness final: `ready-for-implementation`.
