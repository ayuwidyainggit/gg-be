# Plan SX-2184 — `order_type = O` skip stock validation dan inventory mutation

Readiness: `ready-for-implementation`
Plan Quality Gate: `PASS`
Mode: Maintenance Stability Mode
Task ID: `20260608-1118-sx-2184-order-type-o-stock`

## Goal

Perbaiki `POST /sales/v1/orders` supaya request dengan `order_type = "O"` diperlakukan sebagai Taking Order/Purchase Order: qty boleh melebihi `inv.warehouse_stock`, tidak menjalankan stock validation, tidak insert `inv.stock`, tidak update `inv.warehouse_stock`, dan tetap menyimpan data purchase-order sesuai docs user.

## Non-goals

- Tidak mengubah behavior `order_type = "SO"` dan payload lama tanpa `order_type`.
- Tidak mendefinisikan logic baru untuk `order_type = "C"`; pertahankan existing/as-is.
- Tidak mengubah promo consult, discount, VAT, reward product, approval, invoice, proforma, cancel, atau update/final-order flow kecuali compile/test menuntut perubahan field additive.
- Tidak menyalin token/auth header dari Google Docs, Jira, curl example, env, atau sumber mana pun.
- Tidak mengubah module selain `sales` kecuali ada compile break yang terbukti terkait.

## Scope

- Module target: `sales`.
- Endpoint target: `POST /sales/v1/orders` lewat `sales/controller/order_controller.go`.
- Area perubahan: DTO/entity, model, migration additive, create-order controller branch, create-order service mapping/stock branch, unit/integration tests.
- Area DB target: `sls.order.order_type`, `sls.order_detail.original_qty_po1/2/3`, serta existing `qty_po*`/PO fields.

## Requirements

- `CreateOrderBody` menerima optional `order_type` dengan valid values `O`, `C`, `SO`; missing/empty tetap backward compatible.
- Untuk `order_type = "O"`:
  - Controller tidak boleh memanggil warehouse stock validation yang membuat error insufficient stock.
  - Validation result untuk order header harus merepresentasikan stock validation tidak dilakukan: `validate_stok = false`, `validate_stok_message = NULL` atau equivalent nil-safe sesuai model.
  - Service tidak boleh memanggil `StockRepository.SalesStockUpdates` saat create order, sehingga tidak ada insert `inv.stock` dan tidak ada upsert/decrement `inv.warehouse_stock`.
  - `sls.order.order_type = "O"` tersimpan.
  - `sls.order.opr_type = "O"` bila docs/user mapping dan existing field tersedia.
  - `sls.order_detail.qty_po1/2/3` dan `original_qty_po1/2/3` tersimpan dari payload `details.normal[].qty1/2/3` atau `qty_po1/2/3` bila payload mengirim keduanya; gunakan helper eksplisit dan dokumentasikan precedence di test.
  - `sls.order_detail.qty_po` tetap hasil konversi ke satuan terkecil memakai `conversion.QtyUnit.ToTotalQuantity()`.
  - Sales-order qty fields untuk O (`qty`, `qty1/2/3`, `qty_final`, `qty*_final`) harus mengikuti docs: `NULL`/belum terisi sampai process order, kecuali executor menemukan constraint DB NOT NULL; bila constraint memaksa non-null, catat blocker dan pilih minimal compatibility dengan alasan.
- Untuk `order_type = "SO"`:
  - Tetap menjalankan `ValidateOrderService.ValidateOrder` dan existing stock mutation behavior.
  - Tidak boleh lolos qty > stock karena fix O.
- Untuk `order_type` nil/empty:
  - Treat as existing/as-is; jangan bypass stock.
- Untuk `order_type = "C"`:
  - Jangan tambah bypass baru kecuali existing code sudah punya pattern; default as-is.

## Acceptance Criteria

- `POST /sales/v1/orders` dengan `order_type = "O"`, wh stock 5, qty 10, berhasil create order.
- Response tidak mengandung error stock validation untuk `O`.
- Tidak ada row baru `inv.stock` akibat create order `O`.
- Tidak ada perubahan `inv.warehouse_stock` akibat create order `O`.
- `sls.order.order_type = 'O'` tersimpan.
- `sls.order.validate_stok = false` dan `validate_stok_message IS NULL` untuk `O`.
- `sls.order_detail.original_qty_po1/2/3` tersimpan sesuai payload original qty.
- `sls.order_detail.qty_po1/2/3` tersimpan sesuai payload qty.
- `sls.order_detail.qty_po` tersimpan sesuai konversi existing.
- `order_type = "SO"` tetap menjalankan validasi stock dan inventory behavior existing.
- Request lama tanpa `order_type` tetap berjalan seperti sebelum perubahan.
- Promo, discount, VAT, order number, reward product, dan status existing tidak regression.

## Existing Patterns/Reuse

- Route create order: `sales/controller/order_controller.go:36-42`, handler `Create` di `67-134`.
- Existing validation flow: `OrderController.Create` map `CreateOrderBody` ke `ValidateOrderBody` lalu `ValidateOrderService.ValidateOrder`.
- Stock validation source: `sales/service/validate_order_service.go:70-109` dan repository `GetWarehouseStockByProducts` di `sales/repository/validate_order_repository.go:252-259`.
- Existing order store transaction: `sales/service/order_service.go:288-490`.
- Existing inventory mutation path: `sales/service/order_service.go:377-390`, `450-463`, `467-472` menuju `sales/repository/stock_repository.go:469-475`.
- Existing PO detail columns: `sales/model/order_detail.go` sudah punya `QtyPo1/2/3`, `SellPricePo1/2/3`, `PromoPo*`, `DiscPo`, `VatValuePo`.
- Existing qty conversion: `conversion.QtyUnit.ToTotalQuantity()`.
- Existing test mocks and patterns: `sales/service/order_service_test.go`.
- Prior related plan: `.opencode/plans/20260604-1024-sx-2154-order-type.md`; reuse additive field plan, tetapi extend untuk SX-2184 stock validation/mutation bypass.

## Constraints

- Ikuti repo rule: Controller → Service → Repository → DB.
- Semua write tetap dalam service transaction.
- Repository writes harus tetap tx-aware.
- Shell command di repo ini harus pakai `rtk` sesuai `AGENTS.md`.
- Jangan commit atau menyalin secrets/token/env.
- Google Docs/Jira membutuhkan auth/JS atau terlalu besar via fetch; external docs tidak diverifikasi langsung. Requirement docs berasal dari prompt user dan harus dianggap reference input dari user.
- Migration command untuk `sales` tidak terdokumentasi; SQL harus idempotent.

## Risks

- Root cause utama ada sebelum service: controller memanggil `ValidateOrderService.ValidateOrder` tanpa branch `order_type`, sehingga fix hanya di service tidak cukup.
- Bila create `O` mengosongkan `qty`/`qty1/2/3` tetapi DB column tidak nullable, insert bisa gagal; executor wajib cek schema/migration/DB atau existing constraints.
- Promo/discount code saat create memakai `Qty1/2/3` dan harga untuk kalkulasi. Mengosongkan qty terlalu awal bisa merusak promo/discount. Simpan PO qty di model setelah semua calculation yang butuh request qty selesai, bukan mutasi request mentah sebelum promo/discount.
- `determineSalesOrderStatus` mungkin menganggap `Validate1Success=false` sebagai blocker. Untuk `O`, berikan validation response netral atau override stock flag secara eksplisit.
- Jika `validate_stok_message` model bertipe `string`, menyimpan SQL NULL perlu model pointer atau repository update explicit. Jangan klaim NULL bila DB menyimpan empty string; sesuaikan model/migration dan test.
- `SO` regression risk tinggi karena stock validation dan inventory mutation shared path.

## Decisions/Assumptions

- `order_type == "O"` adalah satu-satunya trigger bypass stock validation dan inventory mutation untuk issue ini.
- `nil`, empty string, dan `"SO"` menjalankan flow existing/as-is.
- `"C"` tidak dibypass kecuali existing implementation sudah demikian.
- Tipe DB `order_type`: gunakan style dari plan SX-2154, yaitu `VARCHAR(2) NULL` plus optional check `IN ('O','C','SO')`, kecuali executor menemukan migration style berbeda di branch terbaru.
- Tipe `original_qty_po1/2/3`: `FLOAT4 NULL`, mengikuti existing `qty_po1/2/3 FLOAT4`.
- Precedence source original qty untuk O: payload `qty1/2/3` adalah sumber utama sesuai prompt; jika executor juga menambahkan `qty_po1/2/3` ke DTO, helper boleh set `qty_po*` dari `qty1/2/3` untuk create O dan tidak perlu menuntut client mengirim `qty_po*`.
- Tidak ada question gate blocking; requirement cukup eksplisit. Jika DB actual tidak memungkinkan `qty` null, itu blocker implementasi yang harus dicatat saat executor memvalidasi schema.

## Execution Source of Truth

Urutan precedence saat implementasi:

1. Instruksi eksplisit user terbaru untuk SX-2184.
2. Safety/security repo: jangan copy token/secrets, jangan gunakan remote DB default, pakai `rtk`, jaga tenant `cust_id`.
3. Non-negotiable Implementation Invariants di plan ini.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria dan Done Criteria.
6. Implementation Steps.
7. Follow-up atau rekomendasi non-blocking.

Jika ada konflik, executor wajib mengikuti source yang lebih tinggi dan menulis konflik tersebut di evidence/verifikasi.

## Non-negotiable Implementation Invariants

- Bypass stock hanya untuk `strings.TrimSpace(order_type) == "O"`; jangan bypass untuk `SO`, nil, empty, atau `C`.
- Controller create harus tidak memanggil stock validation untuk `O`, karena root cause terjadi sebelum service store.
- Inventory mutation create harus tidak dipanggil untuk `O`; tidak cukup hanya mengabaikan error validation.
- `SO` dan no-`order_type` harus tetap memakai validation/mutation path lama.
- `original_qty_po*` hanya diisi pada create awal dari original request qty; process order/edit tidak boleh overwrite original qty kecuali requirement baru eksplisit.
- Jangan mutasi request quantity terlalu awal sehingga promo/discount/VAT berubah tanpa sengaja; branch field persistence di model layer setelah kalkulasi existing selesai.
- Semua DB writes tetap dalam existing transaction.
- Tidak boleh menambahkan log berisi payload credential/token/header auth.

## Do Not / Reject If

- Reject jika fix hanya mengubah `ValidateOrderService` tapi controller masih menjalankan stock query/error untuk `O` tanpa test.
- Reject jika `SO` qty > stock menjadi sukses.
- Reject jika payload tanpa `order_type` menjadi bypass stock.
- Reject jika create `O` masih memanggil `StockRepository.SalesStockUpdates` atau menghasilkan row `inv.stock`.
- Reject jika `inv.warehouse_stock` berubah untuk create `O`.
- Reject jika implementasi membandingkan `"PO"` sebagai nilai API utama tanpa mapping eksplisit; contract user adalah `"O"`.
- Reject jika field original qty diisi dari qty yang sudah disesuaikan stock, bukan original payload.
- Reject jika secret/token dari docs masuk ke kode, test, fixture, log, postman file, atau plan lanjutan.

## Diff Boundary

Allowed file groups:

- `sales/controller/order_controller.go`
- `sales/entity/order.go`
- `sales/entity/order_detail.go`
- `sales/entity/validate_order.go` bila `ValidateOrderBody` perlu order type awareness
- `sales/model/order.go`
- `sales/model/order_detail.go`
- `sales/service/order_service.go`
- `sales/service/validate_order_service.go` hanya jika helper/branch diletakkan di service validation, tetapi controller bypass tetap wajib
- `sales/repository/*` hanya jika persistence/read model membutuhkan explicit field handling
- `sales/migration/sls.order/**` atau `sales/migration/sls.order_detail/**` untuk additive migrations dan rollback optional
- `sales/service/order_service_test.go`, `sales/controller/*_test.go`, atau test baru di `sales/service/`/`sales/controller/`
- OpenAPI/swagger docs bila ditemukan dan digunakan repo
- Evidence path `.opencode/evidence/20260608-1118-sx-2184-order-type-o-stock/**`

Out-of-boundary changes harus direvert atau dijustifikasi di verification evidence sebelum final quality gate. Jangan mengubah `.env`, dump DB, postman env berisi credential, module lain, atau docs umum tanpa kebutuhan terbukti.

## TDD/Test Plan

TDD required: ya, ini perubahan behavior API, validation, persistence, dan inventory side effect.

Existing test patterns:

- `sales/service/order_service_test.go` sudah berisi mock repository/service dan beberapa tests untuk stock mutation skip pada mobile no-change.
- Controller tests ada di `sales/controller/so_controller_test.go`, tetapi create order controller test mungkin perlu dibuat baru jika belum ada.

First failing/regression tests:

1. Controller/unit test: `Create` dengan `order_type = "O"` tidak memanggil `ValidateOrderService.ValidateOrder` atau memakai bypass validation response; tetap memanggil `OrderService.Store` dengan validation response yang membuat stock flag false/null.
2. Controller/unit test: `Create` dengan `order_type = "SO"` tetap memanggil `ValidateOrderService.ValidateOrder` dan meneruskan error insufficient stock seperti existing.
3. Service/unit test: `Store` dengan `order_type = "O"` dan status processed/taking-order tidak memanggil `StockRepository.SalesStockUpdates`.
4. Service/unit test: `Store` dengan `order_type = "SO"` pada status processed tetap memanggil `SalesStockUpdates` sesuai existing.
5. Mapper/unit test: detail create `O` menyimpan `QtyPo1/2/3`, `OriginalQtyPo1/2/3`, dan `QtyPo` dari konversi qty original.
6. Backward compatibility test: nil/empty `order_type` tetap memakai validation existing.

Green step:

- Tambahkan helper kecil dan branch minimal sampai semua test di atas pass.
- Jalankan targeted test terlebih dahulu, lalu `rtk go test ./...` dari folder `sales`.

Refactor step:

- Extract helper yang mudah dibaca, misalnya:
  - `normalizedOrderType(orderType *string) string`
  - `isTakingOrder(orderType *string) bool`
  - `shouldValidateStockOnCreate(orderType *string) bool`
  - `shouldMutateInventoryOnCreate(orderType *string) bool`
  - `applyTakingOrderDetailFields(detail entity.CreateOrderDetBody, target *model.OrderDetail, totalQty float64)`
- Hindari duplicate branch normal/promo bila bisa tanpa memperluas diff.

Edge cases:

- `order_type` nil, `""`, lowercase/space input, invalid enum, `"C"`, product stock row missing, nil qty pointers, zero qty, promo reward product, status `NEED_REVIEW` vs `PROCESSED`.

Commands:

```bash
rtk docker compose -f docker-compose.yml ps
```

```bash
rtk go test ./service -run 'Test.*SX2184|Test.*OrderType|Test.*TakingOrder'
```

```bash
rtk go test ./controller -run 'Test.*SX2184|Test.*Create.*OrderType'
```

```bash
rtk go test ./...
```

## Implementation Steps

1. Sync with plan SX-2154 status: verify whether current branch already has `order_type` and `original_qty_po*` changes. If absent, implement additive fields/migrations from SX-2154 plus SX-2184 behavior in one bounded diff.
2. Add `OrderType *string` to `CreateOrderBody`, `OrderResponse`/`OrderListResponse` only if response/read paths need it, and `model.Order`/`OrderList` with `gorm:"column:order_type"`.
3. Add `OriginalQtyPo1/2/3 *float64` to `model.OrderDetail` and read/response structs if needed for verification or response.
4. Add migration:
   - `sls.order.order_type VARCHAR(2) NULL` with check constraint if compatible.
   - `sls.order_detail.original_qty_po1/2/3 FLOAT4 NULL` with comments.
   - Optional rollback dropping constraint/columns.
5. Implement helper `isTakingOrder` in a suitable package/file. Use exact API value `"O"`; trim space. Prefer not to uppercase silently unless validation permits; if uppercase normalization is added, test it.
6. Update `OrderController.Create`:
   - Parse/validate request as existing.
   - If `isTakingOrder(request.OrderType)`, skip `ValidateOrderService.ValidateOrder` stock/credit flow and pass a validation response tailored for O to store, or call a non-stock validation path only if business requires credit/overdue still checked. Prompt says `validate_stok=false`; it does not require bypass credit. Minimal safe implementation: bypass stock validation but preserve other validations only if easily separable. If not separable, document and use neutral validation only for O per docs.
   - If not O, keep exact existing mapping and `ValidateOrder` call.
7. Update `OrderService.Store`:
   - Ensure `orderModel.OrderType` persists.
   - For O, set `orderModel.OprType = "O"` if nil/empty; preserve explicit valid value if existing mapping already does it.
   - For O, ensure stock validation fields persist as `false`/`NULL`. If `ValidateStokMessage` needs nullable model, adjust model pointer and verify all reads compile.
   - For normal details O, calculate `totalQty` from request `Qty1/2/3` as existing, set `QtyPo`, `QtyPo1/2/3`, `OriginalQtyPo1/2/3`, `SellPricePo1/2/3`, `DiscPo`, `VatValuePo`, promo PO fields as existing snapshot supports.
   - For O, leave sales-order qty/final fields nil if DB permits. If code/DB requires non-null, make the smallest compatibility choice and record it.
   - For non-O, do not change existing qty mapping.
8. Gate inventory mutation:
   - In normal and promo create loops, only append `SalesOrderStockUpdate` when `shouldMutateInventoryOnCreate(request.OrderType)` and existing `isProcessedDataStatus(orderModel.DataStatus)`.
   - Keep final `if len(salesOrderStockUpdateEntities) > 0` unchanged.
9. Add tests from TDD section.
10. Run validations from `sales` folder. If DB/API smoke available, start compose locally and run manual API + SQL checks without copying tokens.
11. Capture evidence: changed files, test output, migration notes, manual DB/API proof or not-run reason.
12. Route to `@quality-gate` for final review because this touches DB/API/inventory behavior.

## Expected Files to Change

- `sales/controller/order_controller.go`
- `sales/entity/order.go`
- `sales/entity/order_detail.go`
- `sales/model/order.go`
- `sales/model/order_detail.go`
- `sales/service/order_service.go`
- `sales/service/order_service_test.go`
- Controller test file if needed, for example `sales/controller/order_controller_test.go`
- `sales/migration/sls.order/add_order_type_and_original_qty_po_fields.sql` or equivalent split migration
- Optional rollback migration
- Optional swagger/openapi docs if present after search

## Agent/Tool Routing

- `@orchestrator`: execute this plan, route bounded work, integrate evidence.
- `@fixer`: implement code/tests/migration in `sales` only.
- `@explorer`: only if executor needs deeper schema/test fixture discovery.
- `@quality-gate`: final signoff for API/DB/inventory/security regression.
- `@librarian`: not needed unless authenticated docs become accessible and conflict with prompt.

## Executor Handoff Prompt

Copyable handoff:

```text
Implement SX-2184 in scylla-be using `.opencode/plans/20260608-1118-sx-2184-order-type-o-stock.md` as source of truth. Scope is module `sales`, endpoint `POST /sales/v1/orders`. Must preserve: `order_type = "O"` bypasses stock validation and create inventory mutation only; `SO`, nil/empty, and `C` stay existing/as-is; no secrets/token copied; writes stay in service transactions; tenant/cust_id rules preserved. Do not touch `.env`, postman credentials, unrelated modules, promo/discount logic beyond required compile-safe mapping, or process-order original qty semantics. Start with failing tests for controller validation bypass and service inventory skip. Add/verify additive migration for `sls.order.order_type` and `sls.order_detail.original_qty_po1/2/3`. Validate with targeted `rtk go test` and full `rtk go test ./...` from `sales`. Return changed files, root cause, fix summary, tests, migration notes, API/DB smoke evidence or not-run reason, and risks for SO/promo/inventory.
```

## Execution-ready Worklist / Handoff Contract

`start_with`: `T1`

| Task | Action | depends_on | owner/lane | validation | exit criteria | status | requires_user_decision | must_preserve | do_not_touch | evidence_update | exit_verification |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T1 | Inspect actual branch for existing SX-2154 field changes and DB column constraints | none | `@explorer`/`@fixer` | grep/read target files; optional DB schema query if local DB running | Executor knows whether fields/migration already exist and whether qty columns allow NULL | ready | no | no source edits during discovery | secrets/env | update SX-2184 evidence notes | file/line/schema notes captured |
| T2 | Add failing tests for `O` controller bypass and `SO` validation as-is | T1 | `@fixer` | targeted controller test command | Tests fail for missing branch or field | ready | no | SO path still calls validate | implementation logic beyond tests | test output in evidence | failing tests demonstrate bug |
| T3 | Add failing tests for service `O` no `SalesStockUpdates` and detail PO/original qty mapping | T1 | `@fixer` | targeted service test command | Tests fail for missing mapping/branch | ready | no | nil/no order_type remains as-is | unrelated service behavior | test output in evidence | failing tests demonstrate expected behavior |
| T4 | Add/verify additive migrations for `order_type` and `original_qty_po*` | T2,T3 | `@fixer` | SQL review; optional local apply | Idempotent SQL exists; no destructive data change | ready | no | nullable/backward compatible | existing data, remote DB | migration path/status | migration reviewed or applied locally |
| T5 | Add DTO/model fields and helper functions | T4 | `@fixer` | compile via targeted tests | Code compiles far enough for branch implementation | ready | no | optional order_type, exact O trigger | validation behavior for SO | changed files list | no compile errors from fields |
| T6 | Implement controller create branch for O stock validation bypass | T5 | `@fixer` | controller tests | O bypasses stock validation; SO/nil still validate | ready | no | credit/other behavior not broadened without evidence | Update/Final endpoints | test output | targeted controller tests pass |
| T7 | Implement service store O mapping and inventory mutation gate | T6 | `@fixer` | service tests | O persists PO/original qty and does not call `SalesStockUpdates`; SO unchanged | ready | no | transaction boundaries and promo/discount calculations | unrelated promo/invoice/update flow | test output | targeted service tests pass |
| T8 | Run full sales validation | T7 | `@fixer` | `rtk go test ./...` from `sales` | All tests pass or unrelated failures documented with evidence | ready | no | no secrets in logs | no repo-wide destructive commands | command output | full validation result captured |
| T9 | Optional local API + DB smoke | T8 | `@fixer` | compose up if needed; API request with local token only; SQL checks | O qty > stock succeeds; no inventory side effect; SO regression checked if feasible | ready | no | use local/dev credentials only, never docs token | remote DB defaults | API/SQL result or not-run reason | smoke evidence captured |
| T10 | Final quality gate review | T8 | `@quality-gate` | diff review, tests, migration, evidence | PASS or actionable blockers | ready | no | acceptance criteria and diff boundary | implementation edits by reviewer | quality gate notes | signoff recorded |

## Validation Commands

From repo root:

```bash
rtk docker compose -f docker-compose.yml ps
```

From `sales` module:

```bash
rtk go test ./controller -run 'Test.*SX2184|Test.*Create.*OrderType'
```

```bash
rtk go test ./service -run 'Test.*SX2184|Test.*OrderType|Test.*TakingOrder'
```

```bash
rtk go test ./...
```

Optional DB verification after local migration and API create:

```sql
SELECT order_type, validate_stok, validate_stok_message, opr_type
FROM sls."order"
WHERE ro_no = '<created_ro_no>' AND cust_id = '<cust_id>';
```

```sql
SELECT original_qty_po1, original_qty_po2, original_qty_po3,
       qty_po1, qty_po2, qty_po3, qty_po,
       qty1, qty2, qty3, qty, qty_final
FROM sls.order_detail
WHERE ro_no = '<created_ro_no>' AND cust_id = '<cust_id>';
```

```sql
SELECT * FROM inv.stock WHERE tr_no = '<created_ro_no>' OR ref_no = '<created_ro_no>' OR ro_no = '<created_ro_no>';
```

```sql
SELECT qty, qty_on_order FROM inv.warehouse_stock
WHERE cust_id = '<cust_id>' AND wh_id = '<wh_id>' AND pro_id = '<pro_id>';
```

Sesuaikan nama kolom ref inventory dengan schema actual sebelum menjalankan query.

## Evidence Requirements

Source strategy:

- Local project discovery: digunakan dan wajib untuk implementasi.
- Official docs/context7: diskip karena tidak ada library/API eksternal version-sensitive; Google Docs user dicoba via webfetch tetapi tidak bisa diekstrak reliably tanpa auth/JS atau terlalu besar.
- GitHub/upstream: diskip karena bugfix bergantung pada repo lokal, bukan upstream.
- Web search: diskip karena requirement berasal dari Jira/docs internal yang diringkas user.
- Browser/screenshot: tidak relevan untuk backend API task.

Evidence yang harus dikumpulkan executor:

- File/line root cause final.
- Diff summary semua file berubah.
- Test output targeted dan full `rtk go test ./...` dari `sales`.
- Migration apply evidence atau not-run reason.
- API/DB smoke evidence bila runtime tersedia.
- Bukti tidak ada token/credential baru dalam diff.
- Quality gate review result.

## Done Criteria

- API `POST /sales/v1/orders` menerima optional `order_type`.
- `order_type = "O"` bisa create order walau qty > wh stock.
- `order_type = "O"` tidak menjalankan stock validation yang menghasilkan insufficient stock.
- `order_type = "O"` tidak insert `inv.stock` dan tidak update/decrement `inv.warehouse_stock` pada create order.
- `sls.order.order_type`, `validate_stok=false`, `validate_stok_message=NULL`, dan `opr_type=O` tersimpan sesuai docs bila DB mendukung nullable message.
- `sls.order_detail.original_qty_po1/2/3`, `qty_po1/2/3`, dan `qty_po` tersimpan sesuai original payload/konversi.
- `SO` qty > stock tetap gagal/as-is; `SO` qty cukup tetap sukses/as-is.
- Payload tanpa `order_type` tetap backward compatible.
- Tests dan migration siap; manual API/DB validation terdokumentasi bila tersedia.
- Tidak ada secret/token/auth header baru di diff.

## Final Planning Summary

- Artifacts dibuat:
  - `.opencode/plans/20260608-1118-sx-2184-order-type-o-stock.md`
  - `.opencode/evidence/20260608-1118-sx-2184-order-type-o-stock/discovery.md`
  - `.opencode/evidence/20260608-1118-sx-2184-order-type-o-stock/index.json`
- Evidence discovery disimpan karena implementer butuh file/line root cause dan branch points; tidak dibersihkan sebagai stale.
- Draft tidak dibuat, jadi tidak ada draft cleanup.
- Key decision: root cause utama ada di controller pre-store validation; fix harus mencakup controller validation bypass untuk O dan service inventory mutation gate untuk O.
- Assumptions: docs user adalah sumber requirement utama; Google Docs/Jira tidak dapat diakses authenticated dari tool; `C` tetap existing/as-is; nil/empty order_type tetap existing/as-is.
- Open questions: tidak ada yang blocking. Potential blocker runtime: jika schema actual tidak mengizinkan nullable qty fields untuk O, executor harus catat dan pilih minimal safe compatibility.
- Readiness: `ready-for-implementation`.
