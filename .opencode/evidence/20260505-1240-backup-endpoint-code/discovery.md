# Discovery Evidence — Backup Endpoint Code

Task ID: `20260505-1240-backup-endpoint-code`
Tanggal: 2026-05-05 12:40 Asia/Jakarta

## Files inspected

- `master/controller/survey_controller.go`
- `master/controller/survey_template_controller.go`
- `master/controller/survey_report_controller.go` (ditemukan via file discovery)
- `master/controller/sales_target_distributor_controller.go` (ditemukan via file discovery)
- `master/controller/m_distributor_controller.go` (ditemukan via file discovery)
- `master/controller/outlet_controller.go` (ditemukan via grep untuk `outlet-list/approval`)
- `master/controller/business_unit_controller.go`
- `master/controller/sales_team_controller.go`
- `master/controller/outlet_type_controller.go`
- `master/controller/outlet_group_controller.go`
- `master/controller/outlet_class_controller.go`
- `master/controller/employee_controller.go`
- `master/controller/emp_group_controller.go`
- `inventory/controller/stock_disposal_controller.go`
- `inventory/controller/gr_controller.go`
- `sales/controller/order_controller.go`
- `sales/controller/invoice_controller.go` (ditemukan via file discovery)
- `finance/controller/expense_controller.go`
- `finance/repository/ar_repository.go` dan grep `collection` untuk AR Collection
- `pjp/router/live_monitoring.go`
- Modul terkait: `master/go.mod`, `inventory/go.mod`, `sales/go.mod`, `finance/go.mod`, `pjp/go.mod`

## Project patterns found

- Repo root bukan git repo; setiap service adalah repo/module Go mandiri (`master`, `inventory`, `sales`, `finance`, `pjp`, dll.).
- Branch aktif saat discovery:
  - `master`: `dev...origin/dev`, clean.
  - `inventory`: `dev...origin/dev`, modified `go.mod`, `go.sum`.
  - `sales`: `dev...origin/dev`, clean.
  - `finance`: `dev...origin/dev`, modified `.env`.
  - `pjp`: `staging...origin/staging`, modified `.air.toml`, `.gitignore`.
- Branch target tersedia lokal/remote:
  - `master`, `inventory`, `sales`, `finance`: `dev`, `qa`, beberapa `demo/28-03-2026`/`demo/28-03-2026-qa` tergantung module.
  - `pjp`: `staging`, `dev`, `demo`.
- Root compose service yang sedang `Up`: `master`, `system`, `redis`; service lain belum berjalan.
- Layering utama Fiber modules: Controller → Service → Repository → DB.
- `pjp` memakai Gin/router terpisah, bukan Fiber.
- Response pattern Fiber memakai `responsebuild.BuildResponse()` dan JWT locals seperti `cust_id`, `parent_cust_id`, `user_id`.

## Endpoint-to-module mapping found

- `master`:
  - Survey: `/v1/survey` pada `survey_controller.go`.
  - Survey Template: `/v1/survey_template` pada `survey_template_controller.go`.
  - Survey Report: file controller/service/repository/model/entity tersedia untuk `survey_report`.
  - Sales Target Distributor: file controller/service/repository/model/entity tersedia.
  - Distributor: `m_distributor_*` files tersedia.
  - Outlet approval: `outlet_controller.go` memuat approval list/approval methods; route uses `/v1/outlet-list` group with approval handler.
  - Field options: areas, regions, distributors, business-unit, sales-teams, salesman/employee, outlet-types, outlet-groups, outlet-classes, outlets, employees, emp-groups tersebar di controller master.
- `inventory`:
  - Stock Disposal: `/v1/stock-disposal`, `/v1/stock-disposal/products` pada `stock_disposal_controller.go`.
  - Goods receipt controller ditemukan sebagai `/v1/goods-receipts`, ada potensi mismatch dengan request `inventory/v1/goods-receipt/{id}` yang perlu diverifikasi.
- `sales`:
  - Orders: `/v1/orders`, `/v2/orders`, `/v1/print_proforma_invoice` pada `order_controller.go`.
  - Endpoint yang diminta seperti `/v1/validate`, `/v1/validate-detail`, `/v1/stock`, `/v1/min-price`, `/v1/conversion-product`, `/v1/generate-invoice`, dan invoice print perlu diverifikasi di controller lain/route names karena grep awal hanya menemukan `print_proforma_invoice` exact.
- `finance`:
  - Expense controller aktif memakai `/v1/expense-type`, bukan `/v1/expense`; request user menyebut `/finance/v1/expense` sehingga perlu verifikasi apakah backup perlu route lama/baru.
  - AR collection logic banyak ada di `ar_repository.go`; controller/service/entity perlu dipetakan lengkap saat backup.
- `pjp`:
  - Live monitoring routes langsung pada `/api/v1/live-monitoring-principal`, `/api/v1/live-monitoring-distributor`, `/api/v1/monitoring_locations/details`.

## Reuse candidates

- Gunakan file/domain yang sudah ada; backup harus berupa snapshot/patch dari branch sumber, bukan reimplementasi.
- Gunakan Git native per module untuk membuat backup branch/tag/patch:
  - `git fetch --all --prune`
  - `git worktree` atau checkout branch sumber/target per module
  - `git diff --binary <target>...<source> -- <paths>` untuk bundle patch per endpoint/module
  - `git bundle` atau branch backup `backup/<date>-<module>-endpoint-code` bila butuh backup immutable lokal.
- Gunakan daftar route/controller/service/repository/model/entity sebagai pathspec awal untuk membatasi backup sesuai endpoint.

## Commands/docs checked

- `rtk docker compose -f docker-compose.yml ps`
- `glob */go.mod`
- `glob **/*survey*`, `**/*stock*disposal*`, `**/*expense*`, `**/*sales_target*`, `**/*distributor*`, `**/*order*`, `**/*invoice*`, `**/*monitoring*`
- `grep` untuk `account-receivables/collection`, `outlet-list/approval`, option endpoints, sales endpoint strings, stock disposal/goods receipt.
- `rtk git status --short --branch` dan branch list per module `master`, `inventory`, `sales`, `finance`, `pjp`.

## Constraints

- Ada perubahan lokal yang tidak boleh tertimpa sebelum backup/replace:
  - `inventory/go.mod`, `inventory/go.sum`
  - `finance/.env` (berpotensi secret; jangan commit/stage)
  - `pjp/.air.toml`, `pjp/.gitignore`
- Instruksi global OpenCode melarang prefix `rtk`, sedangkan instruksi project meminta selalu `rtk`. Discovery mengikuti instruksi project untuk repo ini.
- File repo disebut mengandung plaintext credentials di beberapa config; backup plan harus mencegah ekspose secret.
- Tidak ada root git repo, jadi operasi branch harus dilakukan per service module.

## Risks

- Ambiguitas sumber `dev/staging`: mayoritas module ada di `dev`, `pjp` sedang `staging`; perlu diputuskan apakah backup dari `dev`, `staging`, atau branch aktif per module.
- Ambiguitas target `demo/qa`: beberapa module memakai branch `demo/28-03-2026`, sebagian `demo`, dan `qa`; perlu diputuskan target tepat sebelum replace.
- Endpoint path mismatch dapat menyebabkan backup tidak lengkap jika hanya berdasar daftar user tanpa route discovery lanjutan.
- Local dirty files dapat hilang/tercampur jika langsung checkout/replace tanpa stash/worktree.
- Migrasi SQL terkait endpoint harus ikut dibackup bila replace branch membutuhkan schema parity.
