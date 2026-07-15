# Discovery — SX-1965 Survey Select All Save

## Ringkasan discovery

- Module target: `master`
- Endpoint target: `POST /master/v1/survey`
- Call stack utama:
  - `master/controller/survey_controller.go` → `SurveyController.Create()`
  - `master/service/survey_service.go` → `surveyServiceImpl.Store()`
  - `master/repository/survey_repository.go` / `master/repository/salesman_repository.go`

## File yang diinspeksi

- `master/controller/survey_controller.go`
- `master/service/survey_service.go`
- `master/service/survey_service_test.go`
- `master/repository/survey_repository.go`
- `master/repository/salesman_repository.go`
- `master/entity/survey.go`
- `.opencode/plans/20260504-0846-sx-1906-survey-principal-only.md`
- `.opencode/plans/20260504-2058-sx-1915-salesman-business-unit.md`
- `.opencode/docs/AGENT_ROUTING.md`
- `.opencode/docs/ARCHITECTURE.md`
- `.opencode/docs/QUALITY.md`
- `docker-compose.yml`
- `master/go.mod`

## Pola existing yang relevan

1. `surveyServiceImpl.Store()` dan `Update()` sudah memakai `txManager.WithinTransaction(...)`.
2. `normalizeBusinessUnitSelection()` sudah memperlakukan `0` sebagai sentinel non-distributor dan hanya meneruskan distributor positif.
3. `resolveSurveyCustIds()` saat ini:
   - jika tidak ada distributor positif → fallback ke `custId` request;
   - jika ada distributor positif → resolve child customer melalui `FindCustIdsByDistributorIds(parentCustId, distributorIds)`.
4. Validasi salesman saat ini dilakukan per `emp_id` lewat `resolveSalesmanCustIds()` dengan `salesmanRepo.FindOneByEmpIdAndCustId(...)` pada daftar `cust_id` hasil resolve.
5. Error handling saat ini masih generik untuk kegagalan DB insert:
   - `failed to create survey`
   - `failed to create survey areas`
   - `failed to create survey salesmen`
6. Controller hanya mengenali error sentinel umum seperti `ErrSurveySalesmanNotFound`, tetapi belum mendukung payload error terstruktur berisi daftar `invalid_emp_id` dan `invalid_salesman`.

## Temuan data master dari DB

Query dilakukan ke `scylla_citus_dev` pada tabel `mst.m_salesman`, `smc.m_customer`, dan `mst.m_sales_team`.

### Emp ID dari payload QA

| emp_id | sales_name | cust_id | parent_cust_id | distributor_id | sales_team_code | sales_team_name | is_active | is_del |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 415 | Jaka | C260020001 | C26002 | 102 | 01 | MIX | true | false |
| 421 | Piere Njangka | C260020001 | C26002 | 102 | 02 | GT | true | false |
| 435 | Erling Braut Caraka | C26002 | C26002 | null | 77 | Tim Yuhu | false | false |
| 450 | Bagus Prima | C26002 | C26002 | null | 77 | Tim Yuhu | true | false |
| 458 | Richard | C260020001 | C26002 | 102 | 02 | GT | true | false |
| 459 | Rizal | C260020001 | C26002 | 102 | 02 | GT | true | false |
| 466 | Subiwo | C260020001 | C26002 | 102 | 02 | GT | true | false |

### Salesman yang disorot FE / sales team 77

- `450` — `Bagus Prima`
  - `cust_id = C26002`
  - principal-owned
  - `sales_team_code = 77`
  - `is_active = true`
- `435` — `Erling Braut Caraka`
  - `cust_id = C26002`
  - principal-owned
  - `sales_team_code = 77`
  - `is_active = false`

## Analisis root-cause paling mungkin

Ada dua kandidat kuat yang harus ditangani bersama:

1. **Rule validasi salesman principal vs child distributor belum cukup eksplisit dan error-nya tidak actionable.**
   - Data payload QA mencampur salesman child distributor `C260020001` dan principal-owned `C26002`.
   - Secara code path, ini seharusnya lolos jika `distributor_id` mengandung `0` dan principal memang boleh memilih salesman principal + child distributor.
   - Namun bila ada satu salesman yang tidak lolos validasi, service hanya mengembalikan `ErrSurveySalesmanNotFound` generik tanpa daftar ID/nama invalid.

2. **Data master menunjukkan `emp_id 435` inactive.**
   - FE menandai `Erling Braut Caraka` sebagai suspect, dan DB menunjukkan `is_active = false`.
   - Bila bisnis menganggap salesman inactive tidak valid sebagai target survey, payload QA memang semestinya gagal, tetapi response saat ini tidak memberi tahu siapa yang invalid.
   - Bila bisnis menganggap selection “All” harus hanya mengirim data valid sesuai akses, maka akar masalah bisa berada di FE filtering; BE tetap perlu mengembalikan error yang jelas dan aman.

## Klarifikasi requirement terbaru dari user

- Set salesman yang **boleh diterima** pada create survey harus selaras dengan hasil endpoint list salesman untuk filter yang sama, yaitu pola request seperti:
  - `GET /master/v1/salesman`
  - `sales_team_id=82,81,80,78,77,66,65`
  - `distributor_id=0,102,103,119`
  - `limit=9999`
- Artinya, create survey tidak boleh lebih ketat daripada daftar salesman yang sudah ditampilkan dan bisa dipilih oleh FE untuk scope principal tersebut.
- Token yang dibagikan user diperlakukan sensitif dan **tidak** disalin ke artifact ini.

## Implikasi teknis dari klarifikasi terbaru

1. **Acuan validasi create harus kompatibel dengan scope endpoint `GET /master/v1/salesman`.**
   - Bila list endpoint mengembalikan principal-owned salesman dan child-distributor salesman untuk kombinasi `distributor_id=0,102,103,119`, maka create survey harus menerima himpunan salesman yang sama.
2. **Jika endpoint list masih mengembalikan salesman inactive**, maka ada mismatch kontrak antar endpoint.
   - Dalam kondisi itu implementer harus memutuskan salah satu dengan evidence:
     - menyamakan create dengan list; atau
     - mempertahankan penolakan inactive di create tetapi memperjelas error, lalu mencatat mismatch sebagai root cause lintas endpoint.
3. **Kemungkinan root cause utama bergeser menjadi mismatch antara source-of-truth FE picker dan validator create survey.**
   - Ini lebih spesifik daripada sekadar sentinel `0`.

## Constraint teknis yang perlu dijaga

- Layering tetap `Controller → Service → Repository → DB`.
- Write harus tetap transaction-safe di service layer.
- Sentinel `0` tidak boleh pernah dipakai sebagai distributor riil untuk lookup/insert distributor master.
- Implementasi tidak boleh merusak behavior principal-only dan mixed principal+distributor yang sudah punya regression test di `survey_service_test.go`.

## Reuse candidate

- `normalizeBusinessUnitSelection()` untuk normalisasi `distributor_id`.
- `resolveSurveyCustIds()` sebagai base resolver scope customer.
- `resolveSalesmanCustIds()` sebagai titik pusat validasi scope per salesman, tetapi perlu diperkaya agar bisa mengembalikan detail invalid dan mempertimbangkan status aktif.
- Test doubles di `master/service/survey_service_test.go` sudah sangat cocok untuk Red → Green → Refactor.

## Research gate

- Local project discovery: dipakai dan wajib.
- Official docs/context7: tidak diperlukan karena isu ini domain lokal dan tidak version-sensitive library.
- GitHub/upstream: tidak diperlukan.
- Brave/web search: tidak diperlukan.
- Browser/screenshot: tidak diperlukan untuk planning backend ini.

## Commands dan query yang diperiksa

- `which psql`
- Query data salesman untuk `emp_id IN (450,435,415,421,458,459,466)`
- Query salesman suspect `Bagus Prima`, `Erling Braut Caraka`, `sales_team_code = 77`
- Query `information_schema.columns` untuk `mst.m_survey`, `mst.m_survey_area`, `mst.m_survey_salesman`

## Evidence tambahan sesudah klarifikasi FE scope

- Dicoba hit endpoint production/staging-like salesman dengan scope FE yang diberikan user untuk membuktikan kandidat valid source-of-truth:
  - `GET /master/v1/salesman?page=1&sort=sales_name:asc&sales_team_id=82,81,80,78,77,66,65&distributor_id=0,102,103,119&limit=9999`
- Hasil request runtime yang dicoba dari environment ini: `Unauthorized`.
- Implikasi:
  - Plan implementasi tetap memakai scope builder lokal `buildSalesmanCustScopeCondition(...)` sebagai acuan kompatibilitas code-level.
  - Evidence runtime langsung terhadap endpoint salesman masih perlu dilengkapi saat token dev yang benar-benar valid tersedia di environment retest.

## Risiko

- Acceptance menyebut payload QA harus sukses jika semua salesman valid dalam scope principal; data DB saat ini menunjukkan satu salesman inactive, jadi implementer perlu konfirmasi perilaku bisnis active/inactive atau treat ini sebagai invalid yang harus dijelaskan.
- Perubahan error response pada endpoint create bisa menyentuh compatibility jika ada consumer lain yang mengandalkan message-only.
- Jika logging debug ditambah, harus bersifat sementara, terstruktur, dan tidak mencatat token sensitif.
