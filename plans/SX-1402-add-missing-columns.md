# SX-1402 Phase 2: Add Missing Columns to Download Sales Order

## Context

Setelah fix awal SX-1402 (parsing `salesman_id[]` dan `IN clause`), QA melaporkan bahwa query existing sudah benar tetapi kurang beberapa data:

1. **PO No** → seharusnya diambil dari `sls.order.order_no` (bukan `sls.order.po_no`)
2. **Invoice No** → dari `sls.order.invoice_no`
3. **Invoice Date** → dari `sls.order.invoice_date`

Kolom Invoice No dan Invoice Date ditempatkan setelah Order Date di Excel.

## Scope Perubahan

- **Query tetap sama** (4 query terpisah untuk 4 sheet)
- **Filter tetap sama** (cust_id, parent_cust_id, item_type=1, salesman_id IN)
- **Hanya menambahkan 3 kolom** yang belum di-SELECT dari tabel `sls.order`

## File Yang Perlu Diubah

### 1. `sales/repository/so_repository.go`
Pada semua 4 method `FindDownloadData*`:
- Tambah `sls.order.order_no` di SELECT
- Tambah `sls.order.invoice_no` di SELECT
- Tambah `sls.order.invoice_date` di SELECT

### 2. `sales/model/so_download.go`
Pada semua 4 struct (`SoDownloadPo`, `SoDownloadSo`, `SoDownloadFinal`, `SoDownloadQtySummary`):
- Tambah field `OrderNo *string` (gorm column: order_no)
- Tambah field `InvoiceNo *string` (gorm column: invoice_no)
- Tambah field `InvoiceDate *time.Time` (gorm column: invoice_date)

### 3. `sales/entity/so_download.go`
Pada semua 4 entity row struct:
- Tambah field `OrderNo string` (json: order_no)
- Tambah field `InvoiceNo string` (json: invoice_no)
- Tambah field `InvoiceDate string` (json: invoice_date)

### 4. `sales/service/so_service.go`
Pada semua 4 mapper functions:
- Map `OrderNo` dari model ke entity
- Map `InvoiceNo` dari model ke entity
- Map `InvoiceDate` dari model ke entity (format: 2006-01-02)

Pada semua 4 sheet creator functions:
- Tambahkan header "Order No" sebelum "Po No"
- Tambahkan header "Invoice Date" dan "Invoice No" setelah "Order Date"
- Sesuaikan data row output

### 5. `sales/service/so_service_test.go`
- Update test data model dengan field baru
- Update test assertions untuk memastikan mapping benar

## Urutan Kolom Excel (Setelah Perubahan)

Untuk sheet Purchase Order, Sales Order, Final Order:
```
Order No | Po No | So No | Order Date | Invoice Date | Invoice No | Outlet Code | Outlet Name | ...sisanya sama
```

Untuk sheet QTY Summary:
```
Order No | Po No | So No | Order Date | Invoice Date | Invoice No | Outlet Code | Outlet Name | ...sisanya sama
```

## Catatan
- `sls.order.order_no` sudah pasti ada di tabel karena di-reference di SQL referensi user
- `sls.order.invoice_no` dan `sls.order.invoice_date` juga ada di tabel (dipakai di service lain)
- Tidak ada perubahan di controller atau routing
