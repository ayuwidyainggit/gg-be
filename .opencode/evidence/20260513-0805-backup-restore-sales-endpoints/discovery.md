## Discovery

### Ringkasan tugas
- User meminta backup code terkait endpoint sales berikut:
  - `PATCH https://best.scyllax.online/sales/v1/orders/enhance/SO2604290004`
  - `POST https://best.scyllax.online/sales/v1/validate-order`
- User mengonfirmasi restore berarti **copy file/logic dari branch `dev` ke branch worktree `demo-13052026` yang dibuat dari `qa`**.
- User mengonfirmasi endpoint enhance yang dimaksud adalah **PATCH saja**.

### Files inspected
- `AGENTS.md`
- `sales/controller/order_controller.go`
- `sales/controller/validate_order_controller.go`
- `sales/service/order_service.go`
- `sales/service/validate_order_service.go`
- `sales/entity/edit_order_enhance.go`
- `sales/controller/so_controller_test.go`
- `sales/go.mod`
- `sales/README.md`

### Project patterns found
- Root repo **bukan git repo**; module `sales/` adalah git repo terpisah.
- Command wajib mengikuti instruksi lokal dengan prefix `rtk` saat menjalankan shell command.
- Service `sales` aktif di docker compose root dan status container `scylla-sales` adalah `Up`.
- Routing controller Fiber ditemukan di:
  - `sales/controller/order_controller.go`
  - `sales/controller/validate_order_controller.go`
- Route yang relevan:
  - `PATCH /v1/orders/enhance/:ro_no` → `OrderController.UpdateEnhance`
  - `POST /v1/validate-order/` → `ValidateOrderController.ValidateOrder`
- Alur arsitektur yang terlihat konsisten dengan aturan repo:
  - Controller → Service
  - `OrderController.UpdateEnhance` → `OrderService.UpdateEnhance`
  - `ValidateOrderController.ValidateOrder` → `ValidateOrderService.ValidateOrder`

### Reuse candidates
- Implementasi existing di branch `dev` harus diprioritaskan sebagai source restore, bukan reimplementasi manual.
- Candidate file yang paling mungkin perlu dibandingkan/copy dari `dev` ke branch target:
  - `sales/controller/order_controller.go`
  - `sales/controller/validate_order_controller.go`
  - `sales/service/order_service.go`
  - `sales/service/validate_order_service.go`
  - `sales/entity/edit_order_enhance.go`
- Test pattern yang bisa direuse untuk quick verification controller/helper ada di `sales/controller/so_controller_test.go`.

### Commands/docs checked
- `rtk docker compose -f docker-compose.yml ps`
- `rtk git status` (di `sales/`)
- `rtk git branch --all && rtk git worktree list` (di `sales/`)
- Pembacaan file controller/service/entity lokal melalui tool read.

### Branch/worktree facts
- Branch aktif pada repo `sales/`: `dev`
- Branch target source yang tersedia: `qa`
- Existing branch/worktree terkait demo yang terlihat:
  - branch `demo-05052026`
  - worktree `~/Projects/Geekgarden/scylla-be-restore-worktrees-20260505/sales [demo-05052026]`
- Ada detached worktree lain yang tampak operasional, sehingga perlu hati-hati agar tidak bentrok saat membuat worktree baru.

### Constraints
- Planner ini tidak boleh mengubah source implementation; hanya menulis artefak `.opencode/`.
- Root folder bukan git repo sehingga seluruh operasi branch/worktree harus dilakukan dari `sales/`.
- Nama target sudah diputuskan user: `demo-13052026`.
- Restore yang dimaksud adalah copy logic/file dari `dev` ke branch demo berbasis `qa`, bukan cherry-pick kecuali nanti dipilih sebagai teknik implementasi internal.

### Risks
- Endpoint live user menuliskan `POST` untuk enhance, tetapi code lokal hanya menunjukkan `PATCH`; user sudah mengklarifikasi `PATCH`, jadi implementasi harus mengikuti klarifikasi ini.
- Logic endpoint dapat tersebar di beberapa file service/entity, sehingga restore berbasis file tunggal berisiko tidak lengkap.
- Ada potensi branch `demo-13052026` atau worktree path target sudah ada; implementasi perlu cek eksistensi sebelum create.
- Karena root bukan git repo, menjalankan perintah git dari root akan gagal.

### Research gate decision
- Local project discovery: **required dan sudah dilakukan**.
- Official docs/context7: **tidak diperlukan**, karena tugas ini bergantung pada fakta repo dan git workflow lokal, bukan behavior library yang ambigu.
- GitHub/GitLab upstream API research: **tidak diperlukan** untuk planning awal, karena branch dan file target sudah dapat diinspeksi lokal.
- Brave/web search: **tidak diperlukan**.
- Browser/screenshot evidence: **tidak diperlukan** karena ini bukan task visual/reference UI.
