# Rencana Implementasi — SX-1958 Outlet PJP Tidak Sesuai karena Filter `is_active`

## Goal
Memperbaiki flow backend agar outlet yang dikonfigurasi di PJP tetap muncul pada flow Add New Order walaupun outlet tersebut `is_active = false`, tanpa merusak behavior default endpoint list outlet umum.

## Non-goals
- Tidak mengubah credential, token, atau data sensitif staging/Jira ke source code.
- Tidak mengubah business rule verifikasi outlet selain kebutuhan defect ini.
- Tidak mengubah seluruh contract list outlet menjadi selalu include inactive untuk semua caller.
- Tidak melakukan refactor besar pada arsitektur controller/service/repository di luar seam yang diperlukan.

## Scope
- Service `master`: contract filter endpoint `GET /v1/outlets`.
- Service `pjp`: caller service yang mengambil outlet berdasarkan salesman/PJP.
- Test coverage untuk mixed active/inactive outlet by explicit outlet IDs.
- Manual verification checklist terhadap PJP ID `257` dan salesman `EMP0025 - Piere Njangka`.

## Requirements
1. Endpoint outlet list harus bisa mengembalikan outlet inactive bila caller secara eksplisit meminta include inactive.
2. Default behavior existing tetap aman untuk screen/list lain yang membutuhkan active-only filtering.
3. Flow outlet by salesman/PJP harus mengekspresikan intent secara eksplisit agar tidak bergantung pada perilaku implisit query `outlet_id`.
4. Jika request mengandung `include_inactive=1`, backend harus mengabaikan filter `is_active` untuk request tersebut.
5. Expected outlet codes dari PJP harus dapat lolos filter detail master outlet:
   - `BMI260010`
   - `BMI260011`
   - `BMI260030`
   - `BMI260031`
   - `BMI260029`
   - `BMI260028`
   - `BMI260034`
   - `BMI260035`

## Acceptance Criteria
1. Flow **Sales → Add New Order → pilih salesman EMP0025 - Piere Njangka → pilih Outlet** menampilkan seluruh outlet yang dikonfigurasi di PJP terkait.
2. Outlet inactive yang ada di konfigurasi PJP tetap muncul.
3. Response endpoint list outlet untuk use case PJP mengandung seluruh expected outlet codes di atas.
4. Existing behavior default endpoint umum tetap tidak berubah ketika caller tidak mengirim `include_inactive=1`.
5. Test baru/updated mencakup skenario mixed active/inactive outlet IDs dan precedence `include_inactive` atas `is_active`.

## Existing Patterns/Reuse
- Route list outlet existing: `master/controller/outlet_controller.go` → `List`.
- Filter contract existing: `master/entity/outlet.go::OutletQueryFilter`.
- Filter SQL existing: `master/repository/outlet_repository.go::FindAllByCustId`.
- Caller PJP existing untuk outlet by salesman:
  - `pjp/service/third_party/get_outlet_by_sales_codes_service.go`
  - `pjp/service/third_party/get_outlet_picklist_by_sales_codes_service.go`
- Reuse helper parser existing untuk list query integer:
  - `parseIntSliceQuery`
  - `parseIntSliceQueryAllowZero`
- Reuse test location existing:
  - `master/controller/outlet_controller_test.go`
  - `master/controller/query_filter_parser_test.go`

## Constraints
- Harus mengikuti arsitektur existing **Controller → Service → Repository → DB**.
- Tidak boleh menghardcode data sensitif atau token untuk manual/staging verification.
- Kontrak umum `/v1/outlets` dipakai lintas service; perubahan default berisiko regression tinggi.
- Bukti lokal saat ini menunjukkan repository hanya menghormati `is_active` bila query dikirim; namun evidence Jira menunjukkan web flow Add New Order memang mengirim `is_active=1`, jadi solusi harus robust terhadap request itu.
- Instruksi repo mengharuskan command runtime/test menggunakan prefix `rtk` pada implementasi aktual.

## Risks
- **Regression risk** bila `outlet_id` otomatis mengabaikan `is_active` secara global.
- **Caller divergence risk** bila hanya master diubah tapi caller PJP tidak mengirim flag baru, sementara flow web aktual ternyata lewat caller lain yang tetap menyetel `is_active=1`.
- **Coverage risk** karena belum ada bukti lokal penuh untuk mapping web Add New Order ke endpoint aggregator tertentu.
- **Data risk** jika verification rule ternyata juga menyingkirkan sebagian outlet expected; perlu validasi staging/DB terhadap `verification_status` expected outlet.

## Decisions/Assumptions
- **Keputusan utama:** gunakan **Option A** — tambah flag eksplisit `include_inactive=1` pada endpoint master, dan ubah caller PJP yang relevan untuk mengirim flag itu.
- **Precedence rule:** jika `include_inactive=1`, backend **mengabaikan** `is_active` meskipun caller masih mengirim `is_active=1`.
- **Asumsi:** flow Add New Order dapat dikoordinasikan agar caller/backend aggregator mengirim flag baru.
- **Asumsi:** `verification_status=1` tetap valid untuk use case ini kecuali data staging membuktikan sebaliknya.
- **Open question residual:** route web Add New Order yang persis mem-build request dengan `is_active=1` belum terbukti dari repo lokal ini; perlu manual verification setelah implementasi untuk memastikan caller yang dipakai production/staging sudah ter-cover.

## TDD/Test Plan
### TDD Required
Ya. Ini defect production pada behavior API/filtering dan wajib mengikuti Red → Green → Refactor.

### Reason
Perubahan contract filter mudah menimbulkan regression diam-diam pada screen lain. Test harus mengunci precedence `include_inactive` serta default behavior existing.

### Existing Test Patterns
- Parser/helper tests di `master/controller/query_filter_parser_test.go`
- Controller helper tests di `master/controller/outlet_controller_test.go`
- Bila dibutuhkan stub service/controller minimal, pola lightweight test file yang sama dapat diperluas.

### First Failing / Regression Test (Red)
Tambahkan test yang mengunci parsing/contract baru pada `master/controller/outlet_controller_test.go` atau file test sejenis:
- query `include_inactive=1&is_active=1&outlet_id=10,11&verification_status=1`
- expected: filter yang diteruskan downstream menandai include inactive aktif dan `is_active` effectively ignored.

Tambahkan test repository/service logic:
- Given `OutletQueryFilter{OutletID:[activeID,inactiveID], IsActive=1, IncludeInactive=1}`
- Expected SQL/filter builder **tidak** menambahkan clause `o.is_active = true`.

### Green Step
- Tambah field contract baru `IncludeInactive` pada `OutletQueryFilter`.
- Update controller parsing agar field ini terbaca dari query.
- Update repository filter application: apply `is_active` **hanya bila** `IncludeInactive` false/nil.
- Update caller PJP by salesman & picklist agar mengirim `include_inactive=1` saat memanggil master outlet by `outlet_id`.

### Refactor Step
- Jika pola precedence ini berpotensi dipakai ulang, ekstrak helper kecil untuk menentukan apakah active filter harus diaplikasikan.
- Jaga refactor tetap kecil; jangan ubah struktur besar query builder.

### Edge Cases
- `include_inactive=1` + `is_active=1` → include inactive wins.
- `include_inactive=1` + `is_active=2` → include inactive wins.
- `include_inactive` absent + `is_active=1` → behavior existing tetap active-only.
- `include_inactive` absent + `outlet_id=<mixed ids>` → behavior existing tetap sama.
- `include_inactive=0` atau invalid → fallback ke behavior existing.

### Commands
- `rtk go test ./...` pada module yang diubah (`master`, `pjp`)
- Minimum targeted:
  - `rtk go test ./controller/...`
  - `rtk go test ./service/...`
  - `rtk go test ./repository/...`

## Implementation Steps
1. **Tambah contract filter baru di master entity**
   - Tambahkan field `IncludeInactive *int` atau tipe bool-compatible pada `entity.OutletQueryFilter` dengan tag query `include_inactive`.
   - Gunakan pola yang konsisten dengan `IsActive` agar parsing Fiber tetap sederhana.

2. **Update controller list outlet**
   - Pastikan `c.QueryParser(&dataFilter)` membaca `include_inactive`.
   - Tambahkan guard/normalization ringan bila diperlukan (mis. hanya `1` yang dianggap true).
   - Tidak perlu menghapus `is_active` dari request; cukup biarkan repository memutuskan precedence.

3. **Update repository filter precedence**
   - Di `master/repository/outlet_repository.go::FindAllByCustId`, ubah logic:
     - apply `o.is_active = true/false` hanya jika `IncludeInactive` tidak aktif.
   - Jangan ubah filter `verification_status`, `outlet_id`, pagination, atau scope tenant/distributor yang sudah ada.

4. **Update PJP caller services**
   - `pjp/service/third_party/get_outlet_by_sales_codes_service.go`
   - `pjp/service/third_party/get_outlet_picklist_by_sales_codes_service.go`
   - Tambahkan `include_inactive=1` ke query saat memanggil `.../v1/outlets?outlet_id=...`.
   - Pertahankan `limit=9999` existing behavior kecuali ditemukan constraint baru.

5. **Audit caller terkait lain**
   - Review `pjp/service/visit_service.go` dan caller lain yang memanggil `v1/outlets?outlet_id=...` untuk menilai apakah flow tersebut juga seharusnya include inactive atau tetap default.
   - Jika use case tidak jelas, jangan ubah semuanya sekaligus; catat sebagai follow-up jika perlu.

6. **Tambahkan test coverage**
   - Contract/parser test untuk `include_inactive`.
   - Logic test untuk precedence `include_inactive` over `is_active`.
   - Regression test untuk default behavior tanpa flag.

7. **Manual verification staging**
   - Jalankan flow PJP/Add New Order dengan auth valid dari environment lokal.
   - Pastikan request/caller baru benar-benar mengirim `include_inactive=1`.
   - Verifikasi seluruh expected outlet codes muncul.

## Expected Files to Change
- `master/entity/outlet.go`
- `master/controller/outlet_controller.go`
- `master/repository/outlet_repository.go`
- `master/controller/outlet_controller_test.go`
- mungkin file test tambahan di `master/repository/` atau `master/service/`
- `pjp/service/third_party/get_outlet_by_sales_codes_service.go`
- `pjp/service/third_party/get_outlet_picklist_by_sales_codes_service.go`
- mungkin file test tambahan di `pjp/service/third_party/` bila test harness tersedia

## Agent/Tool Routing
- **Explorer/read tools:** discovery code path dan reuse candidate.
- **Artifact planner:** menulis plan ini sebagai source of truth.
- **Fixer/implementer berikutnya:** melakukan perubahan code + tests sesuai plan.
- **Quality gate setelah implementasi:** review regression, evidence, dan security posture.

## Validation Commands
> Jalankan dari module directory terkait sesuai struktur monorepo.

### Master
- `rtk go test ./...`
- `rtk go test ./controller/...`
- `rtk go test ./repository/...`

### PJP
- `rtk go test ./...`
- `rtk go test ./service/...`

### Manual / API Verification
- Verifikasi request outlet by salesman/PJP sekarang mengandung `include_inactive=1`.
- Verifikasi response untuk outlet IDs configured PJP tetap mengandung outlet inactive.

## Evidence Requirements
1. Screenshot/log request builder atau API trace yang menunjukkan `include_inactive=1` terkirim pada flow PJP/Add New Order.
2. Output test yang membuktikan:
   - default behavior tetap aman
   - include inactive behavior aktif saat flag dikirim
3. Manual validation note yang memetakan expected outlet codes dengan response aktual.
4. Jika tersedia akses DB/staging aman, catat validasi untuk tiap outlet minimum:
   - `outlet_id`
   - `outlet_code`
   - `outlet_name`
   - `is_active`
   - `verification_status`
5. Karena `.opencode/docs/AGENT_ROUTING.md` dan `.opencode/docs/SKILLS.md` tidak ada di repo, keputusan plan ini didasarkan pada `AGENTS.md` dan executable code evidence.

## Done Criteria
- Code change di master dan PJP selesai sesuai scope.
- Test baru/updated pass.
- Manual verification memastikan expected outlet codes tampil.
- Tidak ada hardcoded secret/token/data sensitif di code/test/log.
- Evidence terkumpul cukup untuk quality gate.

## Final Planning Summary
- **Question gate:** sudah dijawab. Arah kontrak dipilih untuk menampilkan outlet regardless of `is_active`, dengan asumsi caller dapat diubah dan validation mencakup local + staging steps.
- **Keputusan kunci:** pakai flag eksplisit `include_inactive=1` pada endpoint master dan ubah caller PJP terkait agar mengirim flag tersebut.
- **Alasan keputusan:** ini paling aman terhadap regression dibanding mengubah perilaku global saat `outlet_id` ada.
- **Asumsi:** web flow dapat/akan melewati caller yang bisa diubah atau minimal backend master akan robust bila flag terkirim bersama `is_active=1`.
- **Open question tersisa:** route caller persis dari web Add New Order yang menambah `is_active=1` belum terpetakan penuh dari repo lokal; perlu diverifikasi saat implementasi/manual staging.
- **Artifact dibuat:**
  - primary plan: `.opencode/plans/20260512-0908-sx-1958-pjp-outlet-inactive.md`
  - evidence kept: `.opencode/evidence/20260512-0908-sx-1958-pjp-outlet-inactive/discovery.md`
- **Cleanup:** draft artifact tidak dibuat karena belum diperlukan. Discovery evidence dipertahankan karena masih operasional untuk implementer dan quality gate.
- **Readiness:** siap untuk implementasi bounded change pada service `master` dan `pjp` dengan TDD/regression coverage.
