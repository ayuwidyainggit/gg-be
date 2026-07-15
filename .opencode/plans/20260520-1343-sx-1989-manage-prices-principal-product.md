# SX-1989 — Manage Prices not updated for principal-created products

Task ID: 20260520-1343-sx-1989-manage-prices-principal-product
Service: `master/` (Fiber)
Issue: Jira SX-1989 (Sprint TAP-18, staging FAILED)
Source of truth untuk implementasi BE.

## Goal

Pastikan harga yang disimpan via Manage Prices (principal scope) ter-refleksi di `GET /master/v1/products?mode=lookup_dist_price` untuk produk principal lama (`JY1-005`) maupun baru, tanpa regression ke produk distributor normal, dengan menyelaraskan kontrak data parent-child di `mst.m_product` dan menambahkan guardrail di write path.

## Non-goals

- Tidak mengubah skema atau makna kolom `mst.m_transaction_price`.
- Tidak menyentuh path lookup di `inventory`, `sales`, `mobile`, `pjp*` selain pembacaan tidak langsung.
- Tidak membuat child product otomatis untuk produk principal lama yang tidak punya child (masuk manual review queue).
- Tidak mengubah RBAC, authentication, atau API contract response.

## Scope

- `master/repository/m_price_repository.go`
- `master/service/m_price_service.go` (audit telemetry + guardrail saat publish)
- `master/repository/product_repository.go` (catatan: tidak ubah lookup logic, hanya validasi via test)
- `master/repository/product_dist_repository.go` (guardrail saat assign distributor product dari principal — minimal validasi `parent_pro_id`)
- `master/service/product_service.go` (pastikan `StoreDistProduct`/`BulkStoreDistProduct` mengisi `parent_pro_id` saat ada konteks principal product)
- `master/migration/mst.m_product/<tanggal>_backfill_parent_pro_id_sx_1989.sql` (forward + rollback) untuk legacy distributor child rows
- `master/service/m_price_service_test.go`, `master/repository/m_price_repository_test.go`, `master/repository/product_repository_test.go` (baru), `master/service/product_service_test.go` (baru jika belum ada untuk path ini)
- `.opencode/evidence/20260520-1343-sx-1989-manage-prices-principal-product/` untuk audit script + verifikasi staging

## Requirements

1. Konfirmasi root cause via audit SQL: distributor child rows untuk produk principal lama memiliki `parent_pro_id` NULL/0, sehingga lookup `pricing_lookup_pro_id = COALESCE(NULLIF(parent_pro_id,0), pro_id)` jatuh ke child `pro_id`, padahal `mst.m_transaction_price` ditulis dengan principal `pro_id`.
2. Backfill `mst.m_product.parent_pro_id` untuk distributor child rows yang dapat dipetakan **unambiguous** ke parent principal (`pro_code` sama dalam scope `parent_cust_id`, `is_del = false`, child memiliki `distributor_id IS NOT NULL`, parent `cust_id = parent_cust_id` dan `distributor_id IS NULL`).
3. Tambahkan guardrail write path agar setiap distributor child row baru yang asalnya dari principal selalu memiliki `parent_pro_id` valid.
4. Tambahkan audit telemetry di `MPriceService.PublishByRMQ`/`applyPublishedProductPrices` untuk kasus distributor target tanpa child terlinked atau dengan link rusak: `warn + audit + lanjut publish` (tidak fail).
5. Pertahankan canonical key `mst.m_transaction_price.pro_id = parent principal pro_id`. Jangan duplikasi ke child `pro_id`.
6. Verifikasi end-to-end di staging: produk JY1-005 setelah Manage Prices di-update menghasilkan `purch_price*` / `sell_price*` baru di response lookup.

## Acceptance Criteria

- [ ] `GET /master/v1/products?mode=lookup_dist_price&...&q=jersey%20persinga%20ngawi%20loh%20ya` mengembalikan harga terbaru untuk JY1-005 setelah Manage Prices di-update.
- [ ] Tidak ada perbedaan behavior antara produk principal lama (April 2026), principal baru (May 2027), dan produk distributor normal.
- [ ] Repository tests `(?s)`-regex baru memverifikasi:
  - lookup distributor menggunakan `pricing_lookup_pro_id = parent` saat `parent_pro_id` valid;
  - update Manage Prices menulis `m_transaction_price.pro_id = principal pro_id` dan tidak duplikasi ke child;
  - guardrail `StoreDist`/`BulkStoreDistProduct` menolak row tanpa `parent_pro_id` saat konteks principal.
- [ ] Service tests untuk `MPriceService.PublishByRMQ` mencatat audit anomaly count saat distributor target tanpa child linked.
- [ ] Migration backfill memiliki rollback, idempoten, hanya menyentuh row `cust_id` distributor yang `distributor_id IS NOT NULL` dan `COALESCE(parent_pro_id, 0) = 0` dengan parent unambiguous.
- [ ] Validasi staging: audit query menunjukkan 0 row distributor product principal-origin dengan `parent_pro_id` rusak setelah backfill.
- [ ] Tidak ada regression di test `master/...`: `rtk go test ./...`.

## Existing Patterns / Reuse

- `MPriceProductSnapshot.ParentProID` (`master/model/m_price.go:135`) sudah membawa `parent_pro_id`.
- `UpdatePrincipalAssignedProductPrices` (`master/repository/m_price_repository.go:604`) sudah pakai `WHERE parent_pro_id = ? AND distributor_id IS NOT NULL` — pola yang sama dipakai untuk audit.
- `FindAffectedDistributorProductIDs`, `FindUpdatedDistributorProductIDs` (`master/repository/m_price_repository.go:528,554`) bisa diperluas untuk hitung anomali tanpa link.
- Pola test repository: `master/repository/m_price_repository_test.go` pakai `DATA-DOG/go-sqlmock` + regex `(?s)`.
- Pola test service: `master/service/m_price_service_test.go` pakai struct stub repo + `repository.MTransactionPriceRepository` stub.
- Pola migration: file SQL pasangan forward + rollback di `master/migration/mst.<table>/<YYYYMMDD>_<slug>.sql`. Contoh: `mst.m_distributor/20260407_add_parent_cust_id_nullable_fk.sql` + `_rollback.sql`.
- Tidak ada utility KiloCode yang menggantikan logika ini; ikut `Reuse > Extend > Create`.

## Constraints

- Layering Controller → Service → Repository wajib; logika audit di service, akses data via repository.
- `mst.m_transaction_price` punya konsumen lain (`m_sp_price_service`, `dist_price_service`, `inventory.wh_trf_service`); jangan mengubah skema atau makna `pro_id`.
- Repo banyak query string-concat. Pertahankan placeholder `$1..$N` untuk perubahan baru, hindari menambah string concat baru tanpa escaping.
- `pq.Array` / `sqlx.In` sudah dipakai; jika tambah filter array baru, pakai pola yang sama.
- Migration harus reversible dan tenant-safe; tidak boleh cross `cust_id` / `parent_cust_id` / `distributor_id`.
- Tidak boleh commit secrets, tidak menyentuh `.env`, dan tidak menjalankan migration langsung di staging tanpa persetujuan eksplisit.

## Risks

- Mismatch link parent-child saat backfill bisa menyebabkan harga principal bocor ke produk salah → mitigasi: hanya backfill saat mapping unambiguous (1:1 berdasarkan `pro_code` dalam parent tenant).
- Audit telemetry tambahan di publish dapat menambah load DB; mitigasi: query audit ringan, hanya jalan saat publish, log + counter di memory.
- Test sqlmock terhadap query CTE besar di `FindAllByDistributorLookupDistPrice` rapuh terhadap whitespace; mitigasi: gunakan regex `(?s)` longgar yang fokus ke clause kunci.
- Distributor 102 token expired dapat menggagalkan verifikasi staging; mitigasi: minta token baru pada tahap verifikasi.

## Decisions / Assumptions

- Strategi data lama: backfill `parent_pro_id` + guardrail + audit telemetry (hasil question gate user).
- Produk tanpa child distributor: skip dulu, masuk manual review queue (hasil question gate user).
- Behavior anomali saat publish: warn + audit + lanjut publish (hasil question gate user).
- Akses staging tersedia saat implementasi (hasil question gate user); test unit/regression dijalankan dulu sebelum verifikasi end-to-end.
- Asumsi: produk principal di `mst.m_product` punya `parent_cust_id == cust_id` dan `distributor_id IS NULL`. Distributor child punya `distributor_id IS NOT NULL` dan `cust_id` = distributor tenant.
- Asumsi: `pro_code` adalah kunci unik untuk pemetaan parent ↔ child dalam scope `parent_cust_id`. Jika tidak unik → row tersebut masuk manual review.
- Audit telemetry minimal: log structured + counter via existing `log.Info/Warn`. Tidak menambah dependency observability baru.

## TDD / Test Plan

TDD wajib karena perubahan menyentuh logika harga dan multi-tenant. Plan Red → Green → Refactor.

**Step Red (failing tests dulu):**

1. `repository/product_repository_test.go::TestFindAllByDistributorLookupDistPrice_HealthyParentLink_ReadsParentTransactionPrice`
   - Setup `mst.m_product` parent (`distributor_id IS NULL`, `pro_id = 100`) dan child (`pro_id = 200`, `parent_pro_id = 100`, `distributor_id = 102`) + `mst.m_transaction_price` dengan `pro_id = 100`.
   - Expect lookup mengembalikan `effective_purch_price1 = parent transaction price`.
   - Sebelum fix: lulus (sanity).
2. `repository/product_repository_test.go::TestFindAllByDistributorLookupDistPrice_LegacyBrokenParent_DoesNotReadParentTransactionPrice`
   - Setup `parent_pro_id = 0` di child untuk produk legacy.
   - Expect harga tetap fallback ke `paged_products.purch_price*` (base) — sebelum backfill.
   - Sebelum fix: lulus; setelah backfill (post-fix DB) skenario ini tidak relevan; test ini menjaga regression behavior code path tetap deterministik.
3. `repository/m_price_repository_test.go::TestUpdatePrincipalAssignedProductPrices_DoesNotTouchProductsWithoutParentLink`
   - Verifikasi WHERE clause hanya sentuh `parent_pro_id = ? AND distributor_id IS NOT NULL`.
4. `repository/m_price_repository_test.go::TestFindBrokenDistributorChildLinks_ReturnsAnomalies`
   - **Red**: method baru `FindBrokenDistributorChildLinks(parentProID, parentCustID, distributorIDs)` belum ada → test fail kompilasi.
   - Setelah ditambah, return distributor ID dengan child rusak (`parent_pro_id = 0` atau row tidak ada).
5. `service/m_price_service_test.go::TestPublishByRMQ_PrincipalScope_LogsAnomaliesButContinues`
   - Stub repository return 1 distributor sehat + 1 distributor anomaly.
   - Expect: `applyPublishedProductPrices` tetap dipanggil, `syncTransactionPrices` insert/update transaction price, anomaly logged dengan counter > 0.
6. `service/m_price_service_test.go::TestPublishByRMQ_DoesNotDuplicateTransactionPriceToChildProID`
   - Verifikasi `setupTransactionPriceData/Update` tetap pakai `price.ProID` (= principal pro_id) bukan child.
7. `service/product_service_test.go::TestStoreDistProduct_PrincipalContext_RequiresParentProID`
   - **Red**: guard belum ada. Test memastikan `ProductDistCreate` di-build dengan `parent_pro_id = principal pro_id` jika `productID` memiliki `distributor_id IS NULL` (principal-origin).

**Step Green:**

- Tambahkan `MPriceRepository.FindBrokenDistributorChildLinks` (read-only) yang me-`SELECT distributor_id FROM mst.m_product WHERE cust_id IN (...)` dengan `parent_pro_id IS DISTINCT FROM <parentProID>` dan `distributor_id IN (...)`. Jangan tulis apa pun.
- Update `MPriceService.applyPublishedProductPrices` (principal scope) untuk juga menghitung anomali via method baru, lalu `log.Warn` + counter (mis. `mPricePublishAnomalyCounter`).
- Update `productServiceImpl.StoreDistProduct` (`master/service/product_service.go:4257`) dan `BulkStoreDistProduct` (4350) agar mengisi `ParentProId` di `model.ProductDistCreate` saat product principal asal punya `distributor_id IS NULL`. Tambahkan field `ParentProId` ke `model.ProductDistCreate` + kolom di INSERT `m_product_dist` (`product_dist_repository.go`/`product_repository.go::StoreDist`).
- Backfill SQL forward:
  ```sql
  -- master/migration/mst.m_product/20260520_backfill_parent_pro_id_sx_1989.sql
  WITH parent_lookup AS (
    SELECT parent.pro_code, parent.pro_id AS parent_pro_id, parent.cust_id AS parent_cust_id
    FROM mst.m_product parent
    WHERE parent.distributor_id IS NULL
      AND parent.is_del = false
  ), child_to_fix AS (
    SELECT child.pro_id, child.cust_id, child.pro_code, parent_lookup.parent_pro_id
    FROM mst.m_product child
    JOIN mst.m_distributor d
      ON d.cust_id = child.cust_id
     AND d.distributor_id = child.distributor_id
    JOIN parent_lookup
      ON parent_lookup.parent_cust_id = d.parent_cust_id
     AND parent_lookup.pro_code = child.pro_code
    WHERE child.distributor_id IS NOT NULL
      AND COALESCE(child.parent_pro_id, 0) = 0
      AND child.is_del = false
    GROUP BY child.pro_id, child.cust_id, child.pro_code, parent_lookup.parent_pro_id
    HAVING COUNT(DISTINCT parent_lookup.parent_pro_id) = 1
  )
  UPDATE mst.m_product child
     SET parent_pro_id = child_to_fix.parent_pro_id,
         updated_at = CURRENT_TIMESTAMP
    FROM child_to_fix
   WHERE child.pro_id = child_to_fix.pro_id
     AND child.cust_id = child_to_fix.cust_id;
  ```
- Backfill SQL rollback (kembalikan ke NULL hanya untuk row yang baru diisi via batch ini → simpan list di tabel `mst.m_product_parent_backfill_audit` opsional; jika tidak, rollback per row tidak deterministik, jadi rollback berbentuk dokumentasi `restore from backup`).
- Tidak menambah duplikasi `m_transaction_price`; canonical tetap principal `pro_id`.

**Step Refactor:**

- Ekstrak konstruksi WHERE audit ke helper jika dipakai berulang.
- Pastikan kode anomaly tidak menambah waktu publish > N ms (target: 1 query tambahan per publish).

**Edge cases:**

- Produk principal yang sengaja tidak punya child di distributor target: jangan tag sebagai anomali jika distributor target tidak dalam `price.DistributorIDs`.
- Produk dengan `pro_code` duplikat antar parent: backfill skip (HAVING = 1).
- Coverage `N` (national) tetap pakai `distributor_id = 0` di `m_transaction_price`; tidak terdampak fix ini.

**Commands:**

```bash
rtk go mod download && rtk go mod tidy   # workdir master
rtk go test ./...                        # workdir master
rtk go test ./service -run TestPublishByRMQ
rtk go test ./repository -run TestFindAllByDistributorLookupDistPrice
```

## Implementation Steps

1. Tambahkan migration backfill + rollback (file SQL pair) di `master/migration/mst.m_product/`.
2. Tambah field `ParentProId` di `model.ProductDistCreate`; perbarui INSERT `m_product_dist` di `product_repository.go::StoreDist` + assignment di `productServiceImpl.StoreDistProduct` / `BulkStoreDistProduct` agar mengisi parent saat product principal-origin.
3. Tambah method baru `MPriceRepository.FindBrokenDistributorChildLinks` (read-only) + test sqlmock.
4. Update `MPriceService.applyPublishedProductPrices` untuk principal scope: panggil method baru, `log.Warn` dengan struktur `{ price_id, parent_pro_id, distributor_id, reason }`, increment counter (sederhana: variable `int` di service + log).
5. Tulis test repository + service mengikuti TDD plan; pastikan stub `MPriceRepository` di-extend untuk method baru.
6. Verifikasi local: `rtk go test ./...` di `master/`.
7. Setelah lulus test, jalankan migration di staging (memerlukan persetujuan), lalu re-publish Manage Prices terdampak (atau panggil `Publish` ulang) untuk produk JY1-005.
8. Verifikasi staging via curl evidence di Jira (lihat `Validation Commands`).
9. Update `evidence/visual-comparison` tidak relevan (BE only); cukup catat hasil curl + audit SQL ke `evidence/`.

## Expected Files to Change

- `master/migration/mst.m_product/20260520_backfill_parent_pro_id_sx_1989.sql` (new)
- `master/migration/mst.m_product/20260520_backfill_parent_pro_id_sx_1989_rollback.sql` (new, dokumentasi restore)
- `master/model/product_dist.go` (extend `ProductDistCreate`)
- `master/repository/product_repository.go` (`StoreDist` SQL: tambah kolom `parent_pro_id`)
- `master/service/product_service.go` (`StoreDistProduct`, `BulkStoreDistProduct` isi `ParentProId`)
- `master/repository/m_price_repository.go` (method `FindBrokenDistributorChildLinks`)
- `master/service/m_price_service.go` (audit telemetry di `applyPublishedProductPrices`)
- `master/service/m_price_service_test.go` (extend stub + test baru)
- `master/repository/m_price_repository_test.go` (test method baru + invariants existing)
- `master/repository/product_repository_test.go` (new) — coverage `FindAllByDistributorLookupDistPrice`
- `master/service/product_service_test.go` (new atau extend) — coverage `StoreDistProduct` parent linking

## Agent / Tool Routing

- Implementasi: `@fixer` di service `master`. Tidak boleh menyentuh service lain di luar scope.
- Validasi DB-sensitive query + tenant safety: minta review `@oracle`/`@architect`.
- Final signoff sebelum merge: `@quality-gate` (security/tenant/data integrity).
- Verifikasi staging end-to-end (curl + audit SQL): `@quality-gate` atau dev BE setelah credential disediakan.

## Execution-ready Worklist / Handoff Contract

| id | depends_on | owner | task | validation | exit | status | requires_user_decision |
|----|------------|-------|------|------------|------|--------|------------------------|
| T1 | none | @fixer | Tambah field `ParentProId` di `model.ProductDistCreate` + perbarui INSERT `StoreDist` (`product_repository.go`) | `rtk go build ./...` workdir master | model & repo kompilasi ulang dengan field baru | ready | no |
| T2 | T1 | @fixer | Pastikan `productServiceImpl.StoreDistProduct` & `BulkStoreDistProduct` mengisi `ParentProId` jika product principal-origin (`distributor_id IS NULL`) | `rtk go test ./service -run StoreDistProduct` | test guardrail principal-origin lulus | ready | no |
| T3 | none | @fixer | Tambah `MPriceRepository.FindBrokenDistributorChildLinks` + test sqlmock | `rtk go test ./repository -run FindBrokenDistributorChildLinks` | method ada, test merah → hijau | ready | no |
| T4 | T3 | @fixer | Wire audit telemetry di `MPriceService.applyPublishedProductPrices` (principal scope) + extend stub repo + test `TestPublishByRMQ_PrincipalScope_LogsAnomaliesButContinues` & `TestPublishByRMQ_DoesNotDuplicateTransactionPriceToChildProID` | `rtk go test ./service -run TestPublishByRMQ` | publish lanjut + anomali dilog | ready | no |
| T5 | none | @fixer | Tambah test repository `FindAllByDistributorLookupDistPrice` healthy + legacy-broken (regex `(?s)`) | `rtk go test ./repository -run FindAllByDistributorLookupDistPrice` | dua test baru lulus | ready | no |
| T6 | T1..T5 | @fixer | Jalankan `rtk go mod tidy` + `rtk go test ./...` workdir master | exit 0 | semua test hijau | ready | no |
| T7 | T6 | @oracle | Review tenant safety query backfill + audit telemetry | review tertulis | approval | ready | yes |
| T8 | T7 | @fixer | Tambah migration SQL forward + rollback `20260520_backfill_parent_pro_id_sx_1989*.sql` | `psql --dry-run` jika tersedia, atau review SQL | file siap dijalankan | ready | no |
| T9 | T8 | @quality-gate | Approve eksekusi migration di staging + jadwalkan window | approval tertulis | go/no-go | blocked | yes — perlu credential & window staging |
| T10 | T9 | @fixer | Eksekusi migration backfill di staging via tools resmi (bukan langsung dari OpenCode) | audit query 0 row rusak | data sehat | blocked | yes |
| T11 | T10 | @fixer | Re-publish Manage Prices JY1-005 + verifikasi curl Jira evidence | response harga sesuai update | TC-1 hijau | blocked | yes |
| T12 | T11 | @quality-gate | Final signoff: regression TC-2..TC-4 + tenant audit + close SX-1989 | report tertulis | done | blocked | yes |

`start_with`: T1.

## Validation Commands

```bash
# workdir: master
rtk go mod download && rtk go mod tidy
rtk go test ./...
rtk go test ./service -run TestPublishByRMQ
rtk go test ./repository -run "FindAllByDistributorLookupDistPrice|FindBrokenDistributorChildLinks|UpdatePrincipalAssignedProductPrices"

# staging audit (after backfill)
SELECT COUNT(*) FROM mst.m_product
WHERE distributor_id IS NOT NULL
  AND COALESCE(parent_pro_id, 0) = 0
  AND is_del = false;
-- expected: 0 (atau hanya row yang masuk manual review)

# Reproduksi Jira (token valid)
curl 'https://best.scyllax.online/master/v1/products?mode=lookup_dist_price&dist_price_group_id=2&brand_id=135,138,139&pl_id=131,132,138&sbrand1_id=150,151,155,154&outlet_id=1843&order_date=2026-05-16&q=jersey%20persinga%20ngawi%20loh%20ya&page=1&limit=99' \
  -H 'Authorization: Bearer <staging-token>'
```

## Evidence Requirements

- Hasil `rtk go test ./...` (workdir master) sebelum & sesudah perubahan.
- Audit SQL pre-backfill: jumlah row `distributor_id IS NOT NULL AND COALESCE(parent_pro_id,0) = 0` per `cust_id`.
- Audit SQL post-backfill: 0 row tersisa (atau list manual review).
- Response curl JY1-005 sebelum & sesudah re-publish, simpan ke `.opencode/evidence/20260520-1343-sx-1989-manage-prices-principal-product/curl-staging.md`.
- Log audit telemetry sample dari publish principal price (anomali count > 0 saat ada distributor target tanpa link).

## Done Criteria

- Semua acceptance criteria hijau.
- `@quality-gate` final signoff.
- Migration backfill telah dijalankan di staging dan diverifikasi.
- Dokumentasi evidence lengkap di folder `.opencode/evidence/...`.
- Tidak ada perubahan ke service di luar `master/` selain test yang relevan.

## Final Planning Summary

- Artifacts dibuat:
  - `.opencode/plans/20260520-1343-sx-1989-manage-prices-principal-product.md` (primary, source of truth).
  - `.opencode/evidence/20260520-1343-sx-1989-manage-prices-principal-product/discovery.md` (dipertahankan: rujukan path file & root cause untuk `@fixer`).
  - `.opencode/draft/20260520-1343-sx-1989-manage-prices-principal-product/open-questions.md` (akan dihapus setelah konsolidasi karena semua pertanyaan sudah dijawab di Decisions/Assumptions).
- Keputusan kunci:
  - Backfill `parent_pro_id` untuk row distributor child legacy unambiguous.
  - Skip auto-create child untuk produk principal lama tanpa child (manual review).
  - Audit telemetry warn-and-continue saat publish menemukan anomali.
  - Pertahankan canonical `m_transaction_price.pro_id = principal pro_id`.
- Asumsi: `pro_code` unik dalam scope `parent_cust_id`; staging access tersedia saat implementasi.
- Open questions: tidak ada yang masih material; semua sudah dijawab via question gate.
- Readiness: implementasi siap dimulai dari T1; T9–T12 menunggu credential & window staging dari user.
- Cleanup: draft `open-questions.md` akan dihapus setelah primary plan ini di-merge sebagai source of truth.
