# Plan SX-2435 / SX-2451 — sales order import / order-type follow-up (sales only)

Readiness: `blocked`
Plan Quality Gate: `BLOCKED`
Mode: Maintenance Stability Mode
Task ID: `20260706-sx-2435-2451-sales-order-import`

## Goal

Siapkan artefak eksekusi yang aman dan minimal untuk implementasi SX-2435 dan SX-2451 di module `sales` saja, dengan fokus area sales order import / order-type follow-up, memakai reuse maksimum dari implementasi `order_type = "O"` yang sudah ada dan pola order enhancement yang sudah hidup.

## Non-goals

- Tidak mengubah module selain `sales`.
- Tidak mendesain endpoint/import contract baru dari nol tanpa requirement ticket yang terverifikasi.
- Tidak mengubah flow create order `SO`, `C`, nil, atau empty `order_type` yang sudah ditutup oleh SX-2154 / SX-2184, kecuali ticket SX-2435 / SX-2451 secara eksplisit menuntutnya.
- Tidak mengubah export Excel existing (`GET /sales/v1/download`) kecuali requirement final ticket membuktikan area itu memang bagian scope.
- Tidak mengubah formula stock, promo, invoice, approval, cancel, atau final-order di luar kebutuhan compile-safe / ticket-safe.
- Tidak menyentuh `.env`, credential, dump DB, atau module lain.

## Scope

- Target service: `sales`.
- Fokus discovery dan eksekusi lanjutan di area berikut yang sudah terverifikasi hidup:
  - create order / taking order path: `sales/controller/order_controller.go`, `sales/service/order_service.go`, `sales/service/validate_order_service.go`, `sales/service/order_type_helper.go`
  - order enhance/update path: `sales/controller/order_controller.go`, `sales/service/*update*`, `sales/entity/edit_order_enhance.go`
  - model/entity persistence: `sales/entity/order.go`, `sales/entity/order_detail.go`, `sales/model/order.go`, `sales/model/order_detail.go`
  - additive migration folder: `sales/migration/sls.order/**`
  - tests: `sales/controller/order_controller_test.go`, `sales/service/order_type_helper_test.go`, `sales/service/order_status_helper_test.go`, `sales/service/order_service_test.go`
- Plan ini mencakup eksekusi-ready worklist untuk implementasi **setelah** 1 requirement gap tertutup: definisi tepat SX-2435 dan SX-2451.

## Requirements

Berdasarkan repo evidence saat ini, requirement yang **sudah terverifikasi**:

- `POST /sales/v1/orders` sudah mendukung optional `order_type` dengan enum `O`, `C`, `SO` melalui `sales/entity/order.go`.
- `order_type = "O"` sudah punya helper dan test dedicated di `sales/service/order_type_helper.go` dan `sales/service/order_type_helper_test.go`.
- Controller create sudah punya branch khusus `IsTakingOrder` / `ShouldValidateStockOnCreate` di `sales/controller/order_controller.go`.
- Tidak ditemukan endpoint atau service sales order import file/excel pada codebase `sales` saat ini.
- Tidak ditemukan referensi `SX-2435` atau `SX-2451` di repo, docs, atau plan lama.

Requirement yang **belum terverifikasi dan wajib dikonfirmasi sebelum implementasi source**:

- Apakah SX-2435 / SX-2451 menambah endpoint import baru, atau hanya melengkapi existing order enhance / taking-order flow.
- Jika import baru: format input (JSON upload vs multipart file vs Excel template), source endpoint, dan target persistence fields.
- Jika follow-up order-type: exact delta terhadap perilaku SX-2154/SX-2184 yang sudah hidup.
- Apakah requirement menyentuh create order saja, update enhance tab purchase order, atau detail read/list/export.

## Acceptance Criteria

AC plan/handoff ini dianggap cukup bila:

- Executor tidak perlu menebak file awal, invariants, atau command validasi.
- Source of truth repo-local dan prior-plan reuse jelas.
- Scope gap ticket terdokumentasi eksplisit dan reversible.
- Begitu detail SX-2435 / SX-2451 diberikan, implementer bisa langsung lanjut dari task `T2` tanpa discovery ulang besar.

Acceptance criteria implementasi **sementara** yang aman (menunggu detail ticket):

- Perubahan source harus tetap berada di `sales` only.
- Perubahan harus mempertahankan invariants SX-2184 yang sudah hidup.
- Jika ternyata ticket hanya follow-up taking order / purchase detail, implementasi harus reuse helper existing (`IsTakingOrder`, `takingOrderQtySource`, `applyTakingOrderDetailFields`) dan tidak membuat branch/order-type helper baru yang duplikatif.
- Jika ternyata ticket butuh import baru, implementasi harus mulai dari test/contract failure dulu dan memilih diff terkecil yang mengikuti layering Controller → Service → Repository → DB.

## Existing Patterns/Reuse

Repo evidence yang paling relevan:

- `.opencode/plans/20260604-1024-sx-2154-order-type.md`
  - menambah persistence `order_type` dan `original_qty_po*`.
- `.opencode/plans/20260608-1118-sx-2184-order-type-o-stock.md`
  - menetapkan invariants bypass stock + no inventory mutation hanya untuk `order_type = "O"`.
- `.opencode/plans/20260608-1347-sx-2131-2184-taking-order-audit.md`
  - audit repo-local bahwa flow `O` sudah hidup dan test pass.
- `sales/service/order_type_helper.go`
  - helper canonical `NormalizedOrderType`, `IsTakingOrder`, `ShouldValidateStockOnCreate`, `ShouldMutateInventoryOnCreate`, `BuildCreateOrderValidationBypassResponse`, `takingOrderQtySource`, `applyTakingOrderValidationSnapshot`, `applyTakingOrderDetailFields`.
- `sales/controller/order_controller.go`
  - create order sudah normalize empty string, build validation bypass, dan branch `IsTakingOrder` vs `ShouldValidateStockOnCreate`.
- `sales/controller/order_controller.go:isEnhancePatchRequest`
  - existing enhance payload keys untuk `purchase_order`, `purchase_details`, `sales_order`, `final_order` memberi reuse bila ticket ternyata tentang update/import-style patching, bukan endpoint baru.
- `sales/controller/order_controller_test.go`
  - sudah punya regression test matrix untuk `O`, `nil`, `empty`, `C`, `SO`.
- `sales/service/order_type_helper_test.go`
  - sudah punya high-signal taking-order persistence tests; paling cocok jadi baseline jika SX-2435 / SX-2451 adalah follow-up taking order.

## Constraints

- Repo rule: Controller → Service → Repository → DB.
- Write harus tetap di service transaction.
- Repo ini mewajibkan shell `rtk`-prefixed untuk workflow project-local.
- Validation harus dijalankan dari folder `sales`.
- Compose/runtime baseline mulai dari `rtk docker compose -f docker-compose.yml ps`.
- Tidak ada migration command documented khusus `sales`; additive SQL harus idempotent.
- Repo root bukan git repo; bila perlu branch/hash, cek dari `sales` module git context.

## Risks

- Risiko utama saat ini adalah **false certainty**: ticket SX-2435 / SX-2451 tidak ada di repo, jadi implementasi tanpa clarification berisiko salah endpoint atau salah acceptance.
- Karena belum ada existing order import endpoint/file flow, membangun import baru tanpa ticket detail akan menjadi architecture/product guess, bertentangan dengan maintenance-mode posture.
- Follow-up taking order berisiko regress `SO`/nil/empty/`C` bila implementer menambah branch baru di luar helper existing.
- Area `order_controller.go` dan `order_service.go` padat; perubahan tanpa targeted tests mudah merusak promo/stock/status flow.
- Jika ticket menyentuh update enhance purchase details, ada risiko overwrite `original_qty_po*` yang seharusnya immutable setelah create awal.

## Decisions/Assumptions

- Source strategy: repo-local evidence only. External docs/context7/web research tidak dipakai karena issue scope adalah internal ticket + repo behavior; tidak ada dependency/library decision baru yang version-sensitive.
- Asumsi kerja paling kuat: SX-2435 / SX-2451 berkaitan dengan sales order purchase-order / import follow-up karena user meminta cek prior plans “sales orders/export/import/order-type related”.
- Asumsi ini **belum cukup** untuk langsung coding; status plan tetap `BLOCKED` sampai ada 1 clarifying payload ticket.
- Reuse pertama yang wajib dipertimbangkan implementer adalah helper/order-type path existing, bukan menambah helper/enum/path baru.
- Jika detail ticket nantinya ternyata hanya meminta perbaikan field mapping pada update enhance purchase detail, endpoint baru import tidak perlu dibuat sama sekali. ponytail: ceiling = plan tetap netral terhadap bentuk solusi; upgrade path = konkretkan setelah ticket text diberikan.

## Execution Source of Truth

Urutan precedence implementasi nanti:

1. Instruksi user terbaru + detail exact SX-2435 / SX-2451.
2. Safety/security repo (`AGENTS.md`, `.opencode/docs/ARCHITECTURE.md`, `.opencode/docs/QUALITY.md`).
3. Non-negotiable Implementation Invariants di plan ini.
4. Execution-ready Worklist / Handoff Contract.
5. Acceptance Criteria dan Done Criteria.
6. Existing Patterns/Reuse + prior plans SX-2154 / SX-2184 / SX-2131 audit.

Jika ada konflik, executor wajib mengikuti source dengan precedence lebih tinggi dan catat konflik di evidence.

## Non-negotiable Implementation Invariants

- Module boundary: `sales` only.
- `order_type = "O"` tetap satu-satunya trigger bypass stock validation / inventory mutation create, kecuali detail ticket baru eksplisit mengubah kontrak itu.
- `SO`, `C`, nil, dan empty `order_type` harus tetap flow existing/as-is sampai requirement baru terbukti.
- `original_qty_po*` tidak boleh dioverwrite sembarangan pada edit/update jika requirement tidak eksplisit.
- Reuse helper existing sebelum menambah helper/branch baru.
- Jangan membuat endpoint import/file-upload baru tanpa contract request/response + persistence target yang jelas dari ticket.
- Semua write tetap transaction-aware; repository tidak memegang business logic.
- Jangan copy secret/token/header auth ke source, test, evidence, atau plan.

## Do Not / Reject If

- Reject jika executor langsung membuat endpoint import/excel upload baru tanpa detail ticket.
- Reject jika perubahan mengubah behavior `SO`, `C`, nil, atau empty `order_type` tanpa test regression.
- Reject jika `original_qty_po*` diupdate dari flow edit/enhance tanpa requirement eksplisit.
- Reject jika implementasi menduplikasi helper order-type yang sudah ada di `sales/service/order_type_helper.go`.
- Reject jika fix menyebar ke module lain atau `.env`.
- Reject jika implementasi mengubah export Excel hanya karena kata “import/export” disebut user, tanpa bukti ticket memang menyasar area download.

## Diff Boundary

Allowed file groups untuk implementasi lanjutan setelah requirement ticket jelas:

- `sales/controller/order_controller.go`
- `sales/controller/order_controller_test.go`
- `sales/service/order_service.go`
- `sales/service/order_type_helper.go`
- `sales/service/order_type_helper_test.go`
- `sales/service/validate_order_service.go`
- `sales/service/order_status_helper.go`
- `sales/service/order_status_helper_test.go`
- `sales/entity/order.go`
- `sales/entity/order_detail.go`
- `sales/entity/edit_order_enhance.go`
- `sales/model/order.go`
- `sales/model/order_detail.go`
- `sales/repository/order_repository.go`
- `sales/migration/sls.order/**`
- `sales/migration/sls.order_detail/**`
- optional new test file under `sales/controller/` atau `sales/service/` jika memang paling kecil diff-nya
- `.opencode/evidence/20260706-sx-2435-2451-sales-order-import/**`

Out-of-boundary changes harus direvert atau dijustifikasi di evidence sebelum final quality gate.

## TDD/Test Plan

TDD required: ya.

Red-first plan setelah ticket detail tersedia:

1. Jika ticket adalah follow-up taking order / purchase detail mapping:
   - Tambah failing test di `sales/service/order_type_helper_test.go` atau `sales/controller/order_controller_test.go` yang mereproduksi gap exact field/path.
2. Jika ticket adalah update enhance purchase order/detail:
   - Tambah failing test di service/controller area update enhance; verifikasi field yang boleh berubah dan field immutable (`original_qty_po*`).
3. Jika ticket memang endpoint import baru:
   - Tambah failing controller/service contract test dulu; jangan mulai dari source route.

Baseline validation commands yang sudah valid di repo:

- `rtk docker compose -f docker-compose.yml ps`
- `rtk go test ./controller -run 'Test.*SX2184|Test.*Create.*OrderType'`
- `rtk go test ./service -run 'Test.*SX2184|Test.*OrderType|Test.*TakingOrder'`
- `rtk go test ./...`

Jika ticket mengarah ke update enhance, tambahkan targeted run pada package test terkait setelah file final diketahui.

## Implementation Steps

1. Konfirmasi exact ticket text SX-2435 dan SX-2451 (cukup acceptance singkat / payload sample / endpoint target).
2. Dari `sales`, verifikasi current branch status dan bahwa baseline taking-order tests masih hijau.
3. Petakan ticket ke salah satu reuse path berikut:
   - Path A: create-order taking-order follow-up
   - Path B: update enhance purchase-order / purchase-details follow-up
   - Path C: endpoint import baru (hanya jika ticket benar-benar menyebut import contract baru)
4. Tulis failing test paling sempit pada file test existing yang paling dekat.
5. Implement minimal diff pada helper/controller/service/model sesuai path terpilih.
6. Tambah migration additive hanya jika field DB baru memang diperlukan dan belum ada.
7. Jalankan targeted tests, lalu full `rtk go test ./...` dari `sales`.
8. Jika perubahan menyentuh persistence behavior, lakukan smoke lokal atau minimal SQL verification bila runtime siap.
9. Simpan evidence command output + changed file list + residual risk.
10. Route ke `@quality-gate` untuk signoff final.

## Expected Files to Change

Paling mungkin berubah, tergantung detail ticket final:

- `sales/controller/order_controller.go`
- `sales/controller/order_controller_test.go`
- `sales/service/order_service.go`
- `sales/service/order_type_helper.go`
- `sales/service/order_type_helper_test.go`
- `sales/service/validate_order_service.go`
- `sales/entity/edit_order_enhance.go`
- `sales/repository/order_repository.go`
- optional additive migration under `sales/migration/sls.order/**` or `sales/migration/sls.order_detail/**`

File yang saat ini **tidak ada bukti perlu diubah**:

- `sales/service/so_service.go` (export excel)
- module selain `sales`
- swagger/openapi files (belum diverifikasi ada contract source aktif untuk area ini)

## Agent/Tool Routing

- `@orchestrator`: minta clarification ticket minimum lalu route eksekusi.
- `@fixer`: implement source/tests setelah T1 selesai.
- `@explorer`: optional jika butuh pemetaan test/update-enhance lebih dalam.
- `@quality-gate`: final signoff karena perubahan berpotensi menyentuh API/DB behavior.
- `@artifact-planner`: selesai setelah plan artifact ini.

## Executor Handoff Prompt

```text
Gunakan `.opencode/plans/20260706-sx-2435-2451-sales-order-import.md` sebagai source of truth awal. Scope: sales service only untuk SX-2435 dan SX-2451, dengan reuse maksimum dari flow order_type/taking-order yang sudah ada. Sebelum coding, tutup dulu requirement gap: dapatkan exact ticket text / acceptance singkat / payload sample untuk SX-2435 dan SX-2451, karena repo tidak memiliki referensi ticket maupun endpoint import existing. Must preserve: hanya order_type O yang punya bypass stock/inventory existing; SO/C/nil/empty tetap as-is; original_qty_po* tidak dioverwrite sembarangan; layering Controller → Service → Repository → DB; no secrets. Do not touch: module selain sales, .env, export excel path kecuali ticket membuktikan area itu scope. Validation minimum: dari sales jalankan baseline taking-order tests lalu targeted regression baru, dan akhiri dengan `rtk go test ./...`. Return evidence: clarified scope, changed files, tests, smoke/not-run reason, residual risks. Setelah implementasi selesai, route ke @quality-gate.
```

## Execution-ready Worklist / Handoff Contract

`start_with`: `T1`

| Task | Action | depends_on | owner/lane | validation | exit criteria | status | requires_user_decision | must_preserve | do_not_touch | evidence_update | exit_verification |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T1 | Dapatkan detail exact SX-2435 dan SX-2451 (endpoint/payload/AC singkat) dari user atau source internal yang diberikan user | none | `@orchestrator` | clarifying note tersimpan di evidence/summary | Scope ticket tidak ambigu lagi | blocked | yes | jangan isi gap dengan asumsi | source code | update discovery/evidence summary | exact scope tertulis |
| T2 | Verifikasi baseline current sales order-type path masih hijau | T1 | `@fixer` | `rtk go test ./controller -run 'Test.*SX2184|Test.*Create.*OrderType'` dan `rtk go test ./service -run 'Test.*SX2184|Test.*OrderType|Test.*TakingOrder'` | Baseline existing behavior terbukti | ready | no | invariants SX-2184 tetap hidup | unrelated packages bila tidak perlu | test output | baseline pass atau failure detail |
| T3 | Petakan ticket ke path A/B/C dan pilih file test paling sempit | T2 | `@fixer` | repo read + short mapping note | Executor tahu file target dan tidak replanning besar | ready | no | prefer reuse existing tests/helpers | module lain | mapping note | chosen path justified |
| T4 | Tambah failing regression/contract test sesuai gap ticket | T3 | `@fixer` | targeted `rtk go test` pada package terkait | Test gagal dengan alasan yang sesuai gap | ready | no | non-scope behavior tidak berubah | broad refactor | failing test output | red state tercatat |
| T5 | Implement minimal diff pada controller/service/helper/model/repository | T4 | `@fixer` | targeted tests | Gap ticket tertutup tanpa regress known invariants | ready | no | O-only bypass semantics, transaction/layering | module lain, export path unless in-scope | changed files list | targeted tests pass |
| T6 | Tambah migration additive jika dan hanya jika field DB baru diperlukan | T5 | `@fixer` | SQL review, compile/tests | Migration idempotent dan scoped | ready | no | additive only | destructive schema edits | migration notes | reviewed/applied locally |
| T7 | Full validation dari sales | T5/T6 | `@fixer` | `rtk go test ./...` | Full suite pass atau unrelated failures terdokumentasi | ready | no | no secrets in logs | repo-wide destructive commands | full test output | pass/failure notes |
| T8 | Optional smoke/API/DB verification jika persistence behavior berubah dan runtime siap | T7 | `@fixer` | local compose + local-only request/SQL | DB/API proof captured atau not-run reason jelas | ready | no | local DB only; no token persistence | remote DB defaults | smoke notes | SQL/API summary |
| T9 | Final conformance/risk review | T7/T8 | `@quality-gate` | diff + tests + evidence review | PASS atau blocker eksplisit | ready | no | scope and invariants preserved | source edits by reviewer | quality gate note | signoff status |

### Handoff Payload — T1 clarification gate

```yaml
handoff:
  task_id: 20260706-sx-2435-2451-sales-order-import
  plan_id: 20260706-sx-2435-2451-sales-order-import
  caller: orchestrator
  callee: orchestrator
  scope: Dapatkan detail exact SX-2435 dan SX-2451 sebelum source implementation dimulai.
  claim_level: scoped
  claim_scope: Boleh klaim requirement clarified atau still-blocked-with-question. Tidak boleh klaim implementation started.
  source_basis: .opencode/plans/20260706-sx-2435-2451-sales-order-import.md,.opencode/evidence/20260706-sx-2435-2451-sales-order-import/discovery.md
  must_preserve: Jangan isi gap ticket dengan asumsi repo. Klarifikasi minimal cukup untuk endpoint target/payload/AC singkat.
  do_not_touch: source code,module selain sales
  validation: Simpan summary klarifikasi di evidence atau handoff summary
  exit_criteria: Exact scope SX-2435/SX-2451 tidak ambigu lagi atau blocker ke user tercatat eksplisit.
  evidence_required: .opencode/evidence/20260706-sx-2435-2451-sales-order-import/discovery.md
  depends_on: none
  context_bundle: .opencode/plans/20260706-sx-2435-2451-sales-order-import.md,.opencode/evidence/20260706-sx-2435-2451-sales-order-import/discovery.md
```

### Handoff Payload — T2 baseline verification

```yaml
handoff:
  task_id: 20260706-sx-2435-2451-sales-order-import
  plan_id: 20260706-sx-2435-2451-sales-order-import
  caller: orchestrator
  callee: fixer
  scope: Verifikasi baseline flow order_type/taking-order existing di sales sebelum perubahan SX-2435/SX-2451.
  claim_level: scoped
  claim_scope: Boleh klaim baseline existing verified atau failed-with-evidence. Tidak boleh klaim ticket implemented.
  source_basis: .opencode/plans/20260604-1024-sx-2154-order-type.md,.opencode/plans/20260608-1118-sx-2184-order-type-o-stock.md,.opencode/plans/20260608-1347-sx-2131-2184-taking-order-audit.md,sales/controller/order_controller.go,sales/controller/order_controller_test.go,sales/service/order_type_helper.go,sales/service/order_type_helper_test.go
  must_preserve: O-only bypass semantics existing tidak berubah. Tidak ada source edit di luar kebutuhan validasi discovery. Validation berjalan dari folder sales.
  do_not_touch: module selain sales,.env,export excel path
  validation: rtk go test ./controller -run 'Test.*SX2184|Test.*Create.*OrderType',rtk go test ./service -run 'Test.*SX2184|Test.*OrderType|Test.*TakingOrder'
  exit_criteria: Hasil baseline pass/fail terdokumentasi jelas untuk menjadi titik mulai implementasi.
  evidence_required: .opencode/evidence/20260706-sx-2435-2451-sales-order-import/baseline-tests.md
  depends_on: T1
  context_bundle: sales/controller/order_controller.go,sales/controller/order_controller_test.go,sales/service/order_type_helper.go,sales/service/order_type_helper_test.go
```

### Handoff Payload — T4/T5 implementation lane

```yaml
handoff:
  task_id: 20260706-sx-2435-2451-sales-order-import
  plan_id: 20260706-sx-2435-2451-sales-order-import
  caller: orchestrator
  callee: fixer
  scope: Implement minimal source diff for clarified SX-2435/SX-2451 in sales after failing test exists.
  claim_level: partial
  claim_scope: Boleh klaim source + tests for exact clarified ticket behavior in sales only. Tidak boleh klaim release-ready tanpa quality gate.
  source_basis: .opencode/plans/20260706-sx-2435-2451-sales-order-import.md,.opencode/plans/20260604-1024-sx-2154-order-type.md,.opencode/plans/20260608-1118-sx-2184-order-type-o-stock.md,.opencode/plans/20260608-1347-sx-2131-2184-taking-order-audit.md,sales/controller/order_controller.go,sales/service/order_type_helper.go,sales/service/order_service.go,sales/entity/order.go,sales/entity/edit_order_enhance.go,sales/model/order.go,sales/model/order_detail.go
  must_preserve: sales only,controller-service-repository-db layering,O-only bypass semantics unless clarified ticket explicitly changes it,original_qty_po* tidak dioverwrite sembarangan,all writes transaction-aware
  do_not_touch: module selain sales,.env / secrets / tokens,sales export excel path kecuali clarified ticket membuktikan area itu scope
  validation: targeted rtk go test pada package/tes baru yang dibuat,rtk go test ./...
  exit_criteria: failing regression test berubah hijau dan full sales suite selesai dengan hasil terdokumentasi
  evidence_required: .opencode/evidence/20260706-sx-2435-2451-sales-order-import/implementation.md,.opencode/evidence/20260706-sx-2435-2451-sales-order-import/test-results.md
  depends_on: T1,T2,T3,T4
  context_bundle: sales/controller/order_controller.go,sales/controller/order_controller_test.go,sales/service/order_type_helper.go,sales/service/order_type_helper_test.go,sales/service/order_service.go,sales/entity/order.go,sales/entity/edit_order_enhance.go
```

## Subagent Context Bundle

### Verified by planner

- `confirmed_repo`: `.opencode/docs/ARCHITECTURE.md:7-10` menetapkan layering Controller → Service → Repository → DB, write di service transaction, tenant rules penting.
- `confirmed_repo`: `.opencode/docs/QUALITY.md:3-5, 25` validasi untuk `sales` dijalankan dari folder `sales` dengan `rtk go test ./...` dan compose check dari root.
- `confirmed_repo`: `.opencode/plans/20260604-1024-sx-2154-order-type.md` adalah prior plan persistence `order_type` + `original_qty_po*`.
- `confirmed_repo`: `.opencode/plans/20260608-1118-sx-2184-order-type-o-stock.md` adalah prior plan bypass stock/inventory untuk `O`.
- `confirmed_repo`: `.opencode/plans/20260608-1347-sx-2131-2184-taking-order-audit.md` menyatakan implementasi local repo untuk `O` sudah hidup dan pernah pass targeted/full tests.
- `confirmed_repo`: `sales/controller/order_controller.go` sudah punya branch `IsTakingOrder` dan `ShouldValidateStockOnCreate`.
- `confirmed_repo`: `sales/service/order_type_helper.go` adalah helper canonical untuk taking-order field/validation snapshot.
- `confirmed_repo`: `sales/controller/order_controller_test.go` sudah punya matrix regression `O`, nil, empty, `C`, `SO`.
- `confirmed_repo`: tidak ada endpoint/service order import berbasis file/excel di codebase `sales` saat ini.
- `confirmed_repo`: tidak ada referensi `SX-2435` atau `SX-2451` di repo/plans/docs saat ini.

### Files already read

- `.opencode/docs/index.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `.opencode/plans/20260604-1024-sx-2154-order-type.md`
- `.opencode/plans/20260608-1118-sx-2184-order-type-o-stock.md`
- `.opencode/plans/20260608-1347-sx-2131-2184-taking-order-audit.md`
- `sales/controller/order_controller.go`
- `sales/controller/order_controller_test.go` (grep evidence)
- `sales/service/order_type_helper.go`
- `sales/service/order_type_helper_test.go`
- `docs/Sales Order Enhancement_BE.md`

### Open assumptions

- `assumption`: SX-2435 / SX-2451 masih berada di domain taking-order / purchase-order / import follow-up.
- `assumption`: user belum memberikan exact ticket AC/payload karena berharap planner menemukan konteks dari repo; repo ternyata tidak menyimpan referensi ticket tersebut.
- `assumption`: perubahan paling mungkin ada di create/update order flow, bukan export, tetapi ini belum bisa dinaikkan menjadi fakta.

### Source of truth order

1. User clarification untuk SX-2435 / SX-2451.
2. Plan ini.
3. Prior plans SX-2154 / SX-2184 / audit SX-2131-2184.
4. Current repo source/tests di `sales`.

## Validation Commands

Dari repo root:

- `rtk docker compose -f docker-compose.yml ps`

Dari `sales`:

- `rtk go test ./controller -run 'Test.*SX2184|Test.*Create.*OrderType'`
- `rtk go test ./service -run 'Test.*SX2184|Test.*OrderType|Test.*TakingOrder'`
- `rtk go test ./...`

Jika ticket final menyentuh update enhance path, tambahkan command targeted package yang sesuai setelah file test final dipilih.

## Evidence Requirements

- `.opencode/evidence/20260706-sx-2435-2451-sales-order-import/index.json`
- `.opencode/evidence/20260706-sx-2435-2451-sales-order-import/discovery.md`
- Saat eksekusi: `baseline-tests.md`, `implementation.md`, `test-results.md`, optional `smoke.md`
- Final report harus eksplisit menyebut bahwa planner ini blocked oleh missing ticket text, bukan oleh hambatan teknis repo.
- Mechanical handoff validation wajib dijalankan setelah plan ditulis:

  - `python3 ~/.config/opencode/scripts/subagent-handoff-check.py --plan .opencode/plans/20260706-sx-2435-2451-sales-order-import.md`

## Done Criteria

Plan ini dianggap done bila:

- Artefak plan durable tersimpan di `.opencode/plans/20260706-sx-2435-2451-sales-order-import.md`.
- Exact candidate files, invariants, validation commands, evidence path, done criteria, dan worklist sudah ada.
- Reuse path terhadap SX-2154 / SX-2184 dan helper current repo sudah dijelaskan.
- Status gate jujur: `BLOCKED` karena missing exact ticket scope, bukan dibuat-buat menjadi `PASS`.

Implementasi nanti dianggap done bila:

- exact SX-2435 / SX-2451 behavior jelas,
- targeted regression test hijau,
- `rtk go test ./...` dari `sales` selesai,
- evidence lengkap,
- quality gate final selesai.

## Final Planning Summary

- Artefak dibuat:
  - `.opencode/plans/20260706-sx-2435-2451-sales-order-import.md`
  - `.opencode/evidence/20260706-sx-2435-2451-sales-order-import/`
- Keputusan utama:
  - tidak mengarang contract import baru karena repo tidak memiliki endpoint/import flow maupun referensi SX-2435 / SX-2451,
  - gunakan prior plans SX-2154 / SX-2184 + helper existing sebagai reuse backbone,
  - blok hanya pada 1 clarification input, bukan pada discovery teknis.
- Asumsi terbuka:
  - domain ticket memang terkait taking-order / purchase-order follow-up,
  - belum ada bukti ticket menyasar export.
- Readiness: `blocked`.
- Plan Quality Gate: `BLOCKED`.
- Active-lane reset note: eksekusi source code berikutnya harus dilakukan di lane implementasi (`@orchestrator` → `@fixer` / `@quality-gate`); batas read-only planner ini tidak terbawa ke lane berikutnya.
