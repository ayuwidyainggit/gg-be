# Dokumentasi API Endpoints Scylla Backend

Dokumentasi lengkap semua endpoint API yang tersedia di Scylla Backend, dikelompokkan berdasarkan service.

**Catatan**: Semua endpoint kecuali yang disebutkan secara eksplisit memerlukan JWT authentication melalui middleware `JWTProtected()`.

---

## 📦 Inventory Service

### Goods Receipt (GR)
- **GET** `/v1/goods-receipts` - Mendapatkan daftar goods receipt
- **GET** `/v1/goods-receipts/:gr_no` - Mendapatkan detail goods receipt berdasarkan nomor GR
- **GET** `/v1/goods-receipts/invoice/:invoice_no` - Mendapatkan detail goods receipt berdasarkan nomor invoice
- **GET** `/v1/goods-receipts/lookup-ap` - Lookup goods receipt untuk Account Payable
- **GET** `/v1/goods-receipts/suppliers/list` - Mendapatkan daftar supplier untuk GR
- **GET** `/v1/goods-receipts/warehouses/list` - Mendapatkan daftar warehouse untuk GR
- **GET** `/v1/goods-receipts/distributors/list` - Mendapatkan daftar distributor untuk GR
- **POST** `/v1/goods-receipts` - Membuat goods receipt baru

### Goods Receipt Branch (GR Branch)
- **GET** `/v1/goods-receipts-branch` - Mendapatkan daftar goods receipt branch
- **GET** `/v1/goods-receipts-branch/:gr_branch_no` - Mendapatkan detail goods receipt branch
- **GET** `/v1/goods-receipts-branch/invoice/:invoice_no` - Mendapatkan detail GR branch berdasarkan invoice
- **GET** `/v1/goods-receipts-branch/suppliers/list` - Daftar supplier untuk GR branch
- **GET** `/v1/goods-receipts-branch/distributors/list` - Daftar distributor untuk GR branch
- **GET** `/v1/goods-receipts-branch/warehouses/list` - Daftar warehouse untuk GR branch
- **GET** `/v1/goods-receipts-branch/order-bookings/list` - Daftar order booking untuk GR branch
- **GET** `/v1/goods-receipts-branch/order-bookings/:order_booking_id` - Detail order booking untuk GR branch
- **GET** `/v1/goods-receipts-branch/print/warehouses/list` - Daftar warehouse untuk print
- **POST** `/v1/goods-receipts-branch` - Membuat goods receipt branch baru
- **PATCH** `/v1/goods-receipts-branch/status` - Update status GR branch
- **PATCH** `/v1/goods-receipts-branch/print` - Print GR branch

### AR Branch
- **GET** `/v1/ar-branch` - Mendapatkan daftar AR branch
- **GET** `/v1/ar-branch/:gr_branch_no` - Mendapatkan detail AR branch
- **GET** `/v1/ar-branch/filter/distributors` - Filter distributor untuk AR branch
- **GET** `/v1/ar-branch/filter/suppliers` - Filter supplier untuk AR branch
- **POST** `/v1/ar-branch/payments/:gr_branch_no` - Membuat payment AR branch

### BPPR (Barang Pindah Pabrik ke Ritel)
- **GET** `/v1/bpprs` - Mendapatkan daftar BPPR
- **GET** `/v1/bpprs/:bppr_no` - Mendapatkan detail BPPR
- **POST** `/v1/bpprs` - Membuat BPPR baru
- **PATCH** `/v1/bpprs/:bppr_no` - Update BPPR
- **DELETE** `/v1/bpprs/:bppr_no` - Hapus BPPR

### Sample Issue (SMP ISS)
- **GET** `/v1/sample-issue` - Mendapatkan daftar sample issue
- **GET** `/v1/sample-issue/:smp_iss_no` - Mendapatkan detail sample issue
- **POST** `/v1/sample-issue` - Membuat sample issue baru
- **PATCH** `/v1/sample-issue/:smp_iss_no` - Update sample issue
- **DELETE** `/v1/sample-issue/:smp_iss_no` - Hapus sample issue

### Warehouse Transfer
- **GET** `/v1/stock-transfers` - Mendapatkan daftar stock transfer
- **GET** `/v1/stock-transfers/:wh_trf_no` - Mendapatkan detail stock transfer
- **GET** `/v1/stock-transfers/warehouse/list/:method` - Daftar warehouse berdasarkan method
- **POST** `/v1/stock-transfers` - Membuat stock transfer baru
- **PATCH** `/v1/stock-transfers/:wh_trf_no` - Update stock transfer
- **DELETE** `/v1/stock-transfers/:wh_trf_no` - Hapus stock transfer

### Warehouse Stock Adjustment
- **GET** `/v1/stock-adjustments` - Mendapatkan daftar stock adjustment
- **GET** `/v1/stock-adjustments/:adj_no` - Mendapatkan detail stock adjustment
- **GET** `/v1/stock-adjustments/warehouse/list` - Daftar warehouse untuk adjustment
- **POST** `/v1/stock-adjustments` - Membuat stock adjustment baru
- **PATCH** `/v1/stock-adjustments/:adj_no` - Update status stock adjustment

### Warehouse SO
- **GET** `/v1/warehouse-so` - Mendapatkan daftar warehouse SO
- **GET** `/v1/warehouse-so/:wh_so_no` - Mendapatkan detail warehouse SO
- **POST** `/v1/warehouse-so` - Membuat warehouse SO baru
- **PATCH** `/v1/warehouse-so/:wh_so_no` - Update warehouse SO
- **DELETE** `/v1/warehouse-so/:wh_so_no` - Hapus warehouse SO

### Van SO (Van Sales Order)
- **GET** `/v1/van-so` - Mendapatkan daftar van SO
- **GET** `/v1/van-so/:van_so_no` - Mendapatkan detail van SO
- **POST** `/v1/van-so` - Membuat van SO baru
- **PATCH** `/v1/van-so/:van_so_no` - Update van SO
- **DELETE** `/v1/van-so/:van_so_no` - Hapus van SO

### Van UL (Van Unloading)
- **GET** `/v1/van-ul` - Mendapatkan daftar van unloading
- **GET** `/v1/van-ul/:van_ul_no` - Mendapatkan detail van unloading
- **POST** `/v1/van-ul` - Membuat van unloading baru
- **PATCH** `/v1/van-ul/:van_ul_no` - Update van unloading
- **DELETE** `/v1/van-ul/:van_ul_no` - Hapus van unloading

### Van LO (Van Loading)
- **GET** `/v1/van-lo` - Mendapatkan daftar van loading
- **GET** `/v1/van-lo/:van_lo_no` - Mendapatkan detail van loading
- **POST** `/v1/van-lo` - Membuat van loading baru
- **PATCH** `/v1/van-lo/:van_lo_no` - Update van loading
- **DELETE** `/v1/van-lo/:van_lo_no` - Hapus van loading

### Van BS UL (Van Branch Stock Unloading)
- **GET** `/v1/van-bs-ul` - Mendapatkan daftar van BS unloading
- **GET** `/v1/van-bs-ul/:van_bs_ul_no` - Mendapatkan detail van BS unloading
- **POST** `/v1/van-bs-ul` - Membuat van BS unloading baru
- **PATCH** `/v1/van-bs-ul/:van_bs_ul_no` - Update van BS unloading
- **DELETE** `/v1/van-bs-ul/:van_bs_ul_no` - Hapus van BS unloading

### GDS (Goods Delivery Slip)
- **GET** `/v1/gds` - Mendapatkan daftar GDS
- **GET** `/v1/gds/:gds_no` - Mendapatkan detail GDS
- **POST** `/v1/gds` - Membuat GDS baru
- **PATCH** `/v1/gds/:gds_no` - Update GDS
- **DELETE** `/v1/gds/:gds_no` - Hapus GDS

### Item Stock Change
- **GET** `/v1/item-st-ch` - Mendapatkan daftar item stock change
- **GET** `/v1/item-st-ch/:isc_no` - Mendapatkan detail item stock change
- **POST** `/v1/item-st-ch` - Membuat item stock change baru
- **PATCH** `/v1/item-st-ch/:isc_no` - Update item stock change
- **DELETE** `/v1/item-st-ch/:isc_no` - Hapus item stock change

### Stock
- **GET** `/v1/stocks` - Mendapatkan daftar stock
- **GET** `/v1/stocks/opname-lookup` - Lookup stock untuk opname
- **GET** `/v1/stocks/report` - Laporan stock
- **POST** `/v1/stocks` - Membuat stock baru
- **POST** `/v1/stocks/bulk` - Membuat stock secara bulk

### Warehouse Stock
- **GET** `/v1/warehouse-stocks` - Mendapatkan daftar warehouse stock
- **GET** `/v1/warehouse-stocks/warehouses` - Daftar warehouse
- **GET** `/v1/warehouse-stocks/products` - Daftar produk
- **POST** `/v1/warehouse-stocks` - Upsert warehouse stock
- **POST** `/v1/warehouse-stocks/bulk` - Upsert warehouse stock secara bulk

### Stock Return
- **GET** `/v1/stock-returns` - Mendapatkan daftar stock return
- **GET** `/v1/stock-returns/:return_no` - Mendapatkan detail stock return
- **PATCH** `/v1/stock-returns/:return_no` - Update stock return
- **PATCH** `/v1/stock-returns` - Update batch stock return

### Supplier Return
- **GET** `/v1/supplier-returns` - Mendapatkan daftar supplier return
- **GET** `/v1/supplier-returns/:supplier_return_no` - Mendapatkan detail supplier return
- **GET** `/v1/supplier-returns/suppliers/list` - Daftar supplier untuk return
- **POST** `/v1/supplier-returns` - Membuat supplier return baru
- **PATCH** `/v1/supplier-returns/:supplier_return_no` - Update status supplier return

### Stock Opname
- **GET** `/v1/stock-opname` - Mendapatkan daftar stock opname
- **GET** `/v1/stock-opname/:doc_no` - Mendapatkan report stock opname
- **GET** `/v1/stock-opname/product-hierarchy` - Mendapatkan hierarki produk untuk opname
- **GET** `/v1/stock-opname/statuses` - Mendapatkan daftar status opname
- **POST** `/v1/stock-opname` - Membuat stock opname baru
- **PATCH** `/v1/stock-opname/:doc_no/cancel` - Cancel stock opname

### Order Booking
- **GET** `/v1/order-booking` - Mendapatkan daftar order booking
- **GET** `/v1/order-booking/:order_booking_id` - Mendapatkan detail order booking
- **GET** `/v1/order-booking/status` - Lookup status order booking
- **GET** `/v1/order-booking/approval` - Daftar order booking yang perlu approval
- **POST** `/v1/order-booking` - Membuat order booking baru
- **PUT** `/v1/order-booking/approve/:order_booking_id` - Approve order booking
- **PUT** `/v1/order-booking/reject/:order_booking_id` - Reject order booking
- **DELETE** `/v1/order-booking/:order_booking_id` - Hapus order booking

### Stock Disposal
- **GET** `/v1/stock-disposal` - Mendapatkan daftar stock disposal
- **GET** `/v1/stock-disposal/:stock_disposal_id` - Mendapatkan detail stock disposal
- **POST** `/v1/stock-disposal` - Membuat stock disposal baru

### Files
- **POST** `/v1/files/uploads` - Upload file

---

## 🏢 Master Service

### Product
- **GET** `/v1/products` - Mendapatkan daftar produk (dengan mode: search, lookup, lookup_dist_price)
- **GET** `/v1/products/:pro_id` - Mendapatkan detail produk
- **GET** `/v1/products/principals` - Daftar principal
- **GET** `/v1/products/categories` - Daftar kategori produk
- **GET** `/v1/products/brands` - Daftar brand
- **POST** `/v1/products` - Membuat produk baru
- **POST** `/v1/products/bulk` - Membuat produk secara bulk
- **PATCH** `/v1/products/:pro_id` - Update produk
- **DELETE** `/v1/products/:pro_id` - Hapus produk
- **DELETE** `/v1/products/` - Hapus multiple produk

### Product File Operations
- **GET** `/v1/products-file/export` - Export produk ke file
- **GET** `/v1/products-file/export-instructions` - Download instruksi export/import
- **GET** `/v1/products-file/export-template` - Download template export
- **GET** `/v1/products-file/export-template-update` - Download template update
- **POST** `/v1/products-file/import` - Import produk dari file
- **POST** `/v1/products-file/import-update` - Import update produk dari file

### Outlet
- **GET** `/v1/outlets` - Mendapatkan daftar outlet
- **GET** `/v1/outlets/:outlet_id` - Mendapatkan detail outlet
- **GET** `/v1/outlets/verification-status` - Daftar status verifikasi outlet
- **GET** `/v1/outlets/list-by-distributor` - Daftar outlet berdasarkan distributor
- **GET** `/v1/outlets/export` - Export outlet
- **GET** `/v1/outlets/export-template` - Download template export outlet
- **GET** `/v1/outlets/export-template-update` - Download template update outlet
- **POST** `/v1/outlets` - Membuat outlet baru
- **POST** `/v1/outlets/import` - Import outlet dari file
- **POST** `/v1/outlets/import-update` - Import update outlet dari file
- **POST** `/v1/outlets/approve` - Approve outlet
- **POST** `/v1/outlets/reject` - Reject outlet
- **PATCH** `/v1/outlets/:outlet_id` - Update outlet
- **DELETE** `/v1/outlets/:outlet_id` - Hapus outlet

### Outlet Dropdowns
- **GET** `/v1/dropdown-outlet-type/` - Daftar tipe outlet untuk dropdown
- **GET** `/v1/dropdown-outlet-group/` - Daftar grup outlet untuk dropdown

### Warehouse
- **GET** `/v1/warehouses` - Mendapatkan daftar warehouse
- **GET** `/v1/warehouses/:wh_id` - Mendapatkan detail warehouse
- **POST** `/v1/warehouses` - Membuat warehouse baru
- **PATCH** `/v1/warehouses/:wh_id` - Update warehouse
- **DELETE** `/v1/warehouses/:wh_id` - Hapus warehouse

### Supplier
- **GET** `/v1/suppliers` - Mendapatkan daftar supplier
- **GET** `/v1/suppliers/:sup_id` - Mendapatkan detail supplier
- **POST** `/v1/suppliers` - Membuat supplier baru
- **PATCH** `/v1/suppliers/:sup_id` - Update supplier
- **DELETE** `/v1/suppliers/:sup_id` - Hapus supplier

### Employee
- **GET** `/v1/employees` - Mendapatkan daftar employee
- **GET** `/v1/employees/:emp_id` - Mendapatkan detail employee
- **POST** `/v1/employees` - Membuat employee baru
- **POST** `/v1/employees/create-multiple` - Membuat multiple employee
- **PATCH** `/v1/employees/:emp_id` - Update employee
- **DELETE** `/v1/employees/:emp_id` - Hapus employee

### Salesman
- **GET** `/v1/salesmans` - Mendapatkan daftar salesman
- **GET** `/v1/salesmans/:salesman_id` - Mendapatkan detail salesman
- **POST** `/v1/salesmans` - Membuat salesman baru
- **PATCH** `/v1/salesmans/:salesman_id` - Update salesman
- **DELETE** `/v1/salesmans/:salesman_id` - Hapus salesman

### Vehicle
- **GET** `/v1/vehicles` - Mendapatkan daftar vehicle
- **GET** `/v1/vehicles/:vehicle_id` - Mendapatkan detail vehicle
- **POST** `/v1/vehicles` - Membuat vehicle baru
- **PATCH** `/v1/vehicles/:vehicle_id` - Update vehicle
- **DELETE** `/v1/vehicles/:vehicle_id` - Hapus vehicle

### Bank
- **GET** `/v1/banks` - Mendapatkan daftar bank
- **GET** `/v1/banks/:bank_id` - Mendapatkan detail bank
- **POST** `/v1/banks` - Membuat bank baru
- **PATCH** `/v1/banks/:bank_id` - Update bank
- **DELETE** `/v1/banks/:bank_id` - Hapus bank

### Distributor Price
- **GET** `/v1/distributor-prices` - Mendapatkan daftar distributor price
- **GET** `/v1/distributor-prices/:dist_price_id` - Mendapatkan detail distributor price
- **POST** `/v1/distributor-prices` - Membuat distributor price baru
- **POST** `/v1/distributor-prices/scheduler/publish-unpublish` - Publish/unpublish distributor price
- **PATCH** `/v1/distributor-prices/:dist_price_id` - Update distributor price
- **DELETE** `/v1/distributor-prices/:dist_price_id` - Hapus distributor price

### Master Data (Lainnya)
Semua endpoint berikut mengikuti pola CRUD standar (GET list, GET detail, POST create, PATCH update, DELETE):

- **Outlet Type** - `/v1/outlet-types`
- **Outlet Group** - `/v1/outlet-groups`
- **Outlet Class** - `/v1/outlet-classes`
- **Outlet Location** - `/v1/outlet-locations`
- **District** - `/v1/districts`
- **Market** - `/v1/markets`
- **Industry** - `/v1/industries`
- **Incentive Group** - `/v1/incentive-groups`
- **Sales Team** - `/v1/sales-teams`
- **Sales Type** - `/v1/sales-types`
- **Beat** - `/v1/beats`
- **Sub Beat** - `/v1/sub-beats`
- **Sub Brand 1** - `/v1/sub-brand1`
- **Sub Brand 2** - `/v1/sub-brand2`
- **Discount Group** - `/v1/discount-groups`
- **Price Group** - `/v1/price-groups`
- **PLU Group** - `/v1/plu-groups`
- **Conversion Group** - `/v1/conversion-groups`
- **Return Reason** - `/v1/return-reasons`
- **Return Reason Distributor** - `/v1/return-reason-distributors`
- **Pickup Reason** - `/v1/pickup-reasons`
- **CNDN** - `/v1/cndns`
- **Employee Group** - `/v1/employee-groups`
- **Official** - `/v1/officials`
- **Official Hierarchy** - `/v1/official-hierarchies`
- **Product Distribution** - `/v1/product-dists`
- **TOP (Terms of Payment)** - `/v1/tops`
- **Discount** - `/v1/discounts`
- **TPR (Temporary Price Reduction)** - `/v1/tprs`
- **TPR Limit** - `/v1/tpr-limits`
- **Discount Product** - `/v1/discount-products`
- **PLU Product** - `/v1/plu-products`
- **Conversion Group Detail** - `/v1/conversion-group-details`
- **MSP Price** - `/v1/msp-prices`
- **Remark Promo** - `/v1/remark-promos`
- **M Parking Fund** - `/v1/m-parking-funds`
- **M Periods** - `/v1/m-periods`
- **M Week** - `/v1/m-weeks`
- **M Working Day** - `/v1/m-working-days`
- **Dist Price Group** - `/v1/dist-price-groups`
- **Status** - `/v1/statuses`
- **Province** - `/v1/provinces`
- **Regency** - `/v1/regencies`
- **Division** - `/v1/divisions`
- **Sub District** - `/v1/sub-districts`
- **Ward** - `/v1/wards`
- **Vehicle Type** - `/v1/vehicle-types`
- **M Channel** - `/v1/m-channels`
- **Region** - `/v1/regions`
- **Area** - `/v1/areas`
- **Taking Order** - `/v1/taking-orders`
- **Sub Distributor Group** - `/v1/sub-distributor-groups`
- **Skip Reason** - `/v1/skip-reasons`
- **Reject Reason** - `/v1/reject-reasons`
- **Distributor** - `/v1/distributors`
- **Special Price Group** - `/v1/special-price-groups`
- **Price** - `/v1/prices`
- **Brand** - `/v1/brands`
- **Principal** - `/v1/principals`
- **Product Category** - `/v1/product-categories`
- **Product Line** - `/v1/product-lines`
- **Flavor** - `/v1/flavors`
- **Pack Type** - `/v1/pack-types`
- **Pack Size** - `/v1/pack-sizes`
- **Unit** - `/v1/units`
- **Unit CoreTax** - `/v1/unit-coretax`
- **Product CoreTax** - `/v1/product-coretax`
- **Cons Product** - `/v1/cons-products`
- **Invoice Disc** - `/v1/invoice-discs`
- **Missed Payment Reason** - `/v1/missed-payment-reasons`
- **Manage Minimum Price** - `/v1/manage-minimum-prices`
- **History** - `/v1/histories`

### Files
- **POST** `/v1/files/uploads` - Upload file

---

## 💰 Finance Service

### Account Payable (AP)
- **GET** `/v1/ap` - Mendapatkan daftar account payable
- **GET** `/v1/ap/:ap_no` - Mendapatkan detail account payable
- **POST** `/v1/ap` - Membuat account payable baru
- **PATCH** `/v1/ap/:ap_no` - Update account payable
- **DELETE** `/v1/ap/:ap_no` - Hapus account payable

### AP Payment
- **GET** `/v1/account-payable-payments` - Mendapatkan daftar AP payment
- **GET** `/v1/account-payable-payments/:account_payable_payment_no` - Mendapatkan detail AP payment
- **POST** `/v1/account-payable-payments` - Membuat AP payment baru
- **PATCH** `/v1/account-payable-payments/:account_payable_payment_no` - Update AP payment
- **DELETE** `/v1/account-payable-payments/:account_payable_payment_no` - Hapus AP payment

### AP Payment Lookup
- **GET** `/v1/account-payable-payments-lookup/check-giro` - Lookup balance payment deposit (cheque giro)
- **GET** `/v1/account-payable-payments-lookup/bank-transfer` - Lookup balance payment deposit (bank transfer)
- **GET** `/v1/account-payable-payments-lookup/cndn` - Lookup balance payment deposit (CNDN)
- **GET** `/v1/account-payable-payments-lookup/return` - Lookup balance payment deposit (return)
- **GET** `/v1/account-payable-payments-lookup/invoice-no` - Lookup invoice number

### AP Pay
- **GET** `/v1/ap-pay` - Mendapatkan daftar AP pay
- **GET** `/v1/ap-pay/:ap_pay_no` - Mendapatkan detail AP pay
- **POST** `/v1/ap-pay` - Membuat AP pay baru
- **PATCH** `/v1/ap-pay/:ap_pay_no` - Update AP pay
- **DELETE** `/v1/ap-pay/:ap_pay_no` - Hapus AP pay

### AP List
- **GET** `/v1/ap-list` - Mendapatkan daftar AP list
- **GET** `/v1/ap-list/:ap_list_no` - Mendapatkan detail AP list
- **POST** `/v1/ap-list` - Membuat AP list baru
- **PATCH** `/v1/ap-list/:ap_list_no` - Update AP list
- **DELETE** `/v1/ap-list/:ap_list_no` - Hapus AP list

### AP CNDN
- **GET** `/v1/ap-cndn` - Mendapatkan daftar AP CNDN
- **GET** `/v1/ap-cndn/:ap_cndn_no` - Mendapatkan detail AP CNDN
- **POST** `/v1/ap-cndn` - Membuat AP CNDN baru
- **PATCH** `/v1/ap-cndn/:ap_cndn_no` - Update AP CNDN
- **DELETE** `/v1/ap-cndn/:ap_cndn_no` - Hapus AP CNDN

### AP Distributor Discount
- **GET** `/v1/ap-distributor-discount` - Mendapatkan daftar AP distributor discount
- **GET** `/v1/ap-distributor-discount/:ap_distributor_discount_no` - Mendapatkan detail AP distributor discount
- **POST** `/v1/ap-distributor-discount` - Membuat AP distributor discount baru
- **PATCH** `/v1/ap-distributor-discount/:ap_distributor_discount_no` - Update AP distributor discount
- **DELETE** `/v1/ap-distributor-discount/:ap_distributor_discount_no` - Hapus AP distributor discount

### AP Supplier Invoice Return
- **GET** `/v1/ap-supplier-invoice-return` - Mendapatkan daftar AP supplier invoice return
- **GET** `/v1/ap-supplier-invoice-return/:ap_supplier_invoice_return_no` - Mendapatkan detail AP supplier invoice return
- **POST** `/v1/ap-supplier-invoice-return` - Membuat AP supplier invoice return baru
- **PATCH** `/v1/ap-supplier-invoice-return/:ap_supplier_invoice_return_no` - Update AP supplier invoice return
- **DELETE** `/v1/ap-supplier-invoice-return/:ap_supplier_invoice_return_no` - Hapus AP supplier invoice return

### Account Receivable (AR)
- **GET** `/v1/account-receivables/:invoice_no` - Mendapatkan detail account receivable
- **GET** `/v1/account-receivables/collection` - Mendapatkan daftar collection
- **GET** `/v1/account-receivables/collection/:collection_no` - Mendapatkan detail collection
- **GET** `/v1/account-receivables/collection/job-titles` - Daftar job title untuk collection
- **GET** `/v1/account-receivables/collection/collectors` - Daftar collector
- **GET** `/v1/account-receivables/collection/invoices` - Daftar invoice untuk collection
- **GET** `/v1/account-receivables/collection/filter/outlet-groups` - Filter outlet group untuk collection
- **GET** `/v1/account-receivables/collection/filter/salesmans` - Filter salesman untuk collection
- **GET** `/v1/account-receivables/collection/filter/collectors` - Filter collector untuk collection
- **GET** `/v1/account-receivables/filter/outlets` - Filter outlet
- **GET** `/v1/account-receivables/filter/salesmans` - Filter salesman
- **POST** `/v1/account-receivables/collection` - Membuat collection baru
- **PATCH** `/v1/account-receivables/collection/:collection_no` - Update collection
- **PATCH** `/v1/account-receivables/collection/print/:collection_no` - Print collection
- **DELETE** `/v1/account-receivables/collection/:collection_no` - Hapus collection

### AR Pay
- **GET** `/v1/ar-pay` - Mendapatkan daftar AR pay
- **GET** `/v1/ar-pay/:ar_pay_no` - Mendapatkan detail AR pay
- **POST** `/v1/ar-pay` - Membuat AR pay baru
- **PATCH** `/v1/ar-pay/:ar_pay_no` - Update AR pay
- **DELETE** `/v1/ar-pay/:ar_pay_no` - Hapus AR pay

### AR Settlement
- **GET** `/v1/account-receivables/settlement` - Mendapatkan daftar AR settlement
- **GET** `/v1/account-receivables/settlement/:deposit_no` - Mendapatkan detail AR settlement
- **GET** `/v1/account-receivables/settlement/filter/collectors` - Filter collector untuk settlement
- **GET** `/v1/account-receivables/settlement/filter/deposit-statuses` - Filter deposit status
- **PATCH** `/v1/account-receivables/settlement/approve/:deposit_no` - Approve settlement
- **PATCH** `/v1/account-receivables/settlement/reject/:deposit_no` - Reject settlement

### AR CNDN
- **GET** `/v1/ar-cndn` - Mendapatkan daftar AR CNDN
- **GET** `/v1/ar-cndn/:ar_cndn_no` - Mendapatkan detail AR CNDN
- **POST** `/v1/ar-cndn` - Membuat AR CNDN baru
- **PATCH** `/v1/ar-cndn/:ar_cndn_no` - Update AR CNDN
- **DELETE** `/v1/ar-cndn/:ar_cndn_no` - Hapus AR CNDN

### Bank Transfer
- **GET** `/v1/bank-transfer` - Mendapatkan daftar bank transfer
- **GET** `/v1/bank-transfer/:bank_transfer_no` - Mendapatkan detail bank transfer
- **GET** `/v1/bank-transfer-filter` - Lookup bank untuk filter
- **POST** `/v1/bank-transfer` - Membuat bank transfer baru
- **PATCH** `/v1/bank-transfer/:bank_transfer_no` - Update bank transfer
- **DELETE** `/v1/bank-transfer/:bank_transfer_no` - Hapus bank transfer

### Cheque
- **GET** `/v1/cheque` - Mendapatkan daftar cheque
- **GET** `/v1/cheque/:chq_no` - Mendapatkan detail cheque
- **POST** `/v1/cheque` - Membuat cheque baru
- **PATCH** `/v1/cheque/:chq_no` - Update cheque
- **DELETE** `/v1/cheque/:chq_no` - Hapus cheque

### Cheque Giro
- **GET** `/v1/cheque-giro` - Mendapatkan daftar cheque giro
- **GET** `/v1/cheque-giro/:cheque_giro_no` - Mendapatkan detail cheque giro
- **GET** `/v1/cheque-giro-filter` - Lookup bank untuk cheque giro
- **POST** `/v1/cheque-giro` - Membuat cheque giro baru
- **PATCH** `/v1/cheque-giro/:cheque_giro_no` - Update cheque giro
- **DELETE** `/v1/cheque-giro/:cheque_giro_no` - Hapus cheque giro

### Cheque Giro Clearing
- **GET** `/v1/cheque-giro-clearing` - Mendapatkan daftar cheque giro clearing
- **GET** `/v1/cheque-giro-clearing/:cheque_giro_clearing_no` - Mendapatkan detail cheque giro clearing
- **POST** `/v1/cheque-giro-clearing` - Membuat cheque giro clearing baru
- **PATCH** `/v1/cheque-giro-clearing/:cheque_giro_clearing_no` - Update cheque giro clearing
- **DELETE** `/v1/cheque-giro-clearing/:cheque_giro_clearing_no` - Hapus cheque giro clearing

### Cash Transaction
- **GET** `/v1/cash-tr` - Mendapatkan daftar cash transaction
- **GET** `/v1/cash-tr/:cash_tr_no` - Mendapatkan detail cash transaction
- **POST** `/v1/cash-tr` - Membuat cash transaction baru
- **PATCH** `/v1/cash-tr/:cash_tr_no` - Update cash transaction
- **DELETE** `/v1/cash-tr/:cash_tr_no` - Hapus cash transaction

### Cash Bank Report
- **GET** `/v1/cash-bank-report` - Mendapatkan laporan cash bank
- **GET** `/v1/cash-bank-report/:cash_bank_report_id` - Mendapatkan detail laporan cash bank
- **POST** `/v1/cash-bank-report` - Membuat laporan cash bank baru
- **PATCH** `/v1/cash-bank-report/:cash_bank_report_id` - Update laporan cash bank
- **DELETE** `/v1/cash-bank-report/:cash_bank_report_id` - Hapus laporan cash bank

### Deposit
- **GET** `/v1/deposit` - Mendapatkan daftar deposit
- **GET** `/v1/deposit/:deposit_no` - Mendapatkan detail deposit
- **POST** `/v1/deposit/collection` - Membuat deposit collection
- **POST** `/v1/deposit/payment` - Membuat deposit payment
- **PATCH** `/v1/deposit/:deposit_no` - Update deposit
- **DELETE** `/v1/deposit/:deposit_no` - Hapus deposit

### Deposit Lookup
- **GET** `/v1/deposit-lookup-filter/` - Lookup index untuk deposit
- **GET** `/v1/invoice-list-collection/` - Daftar invoice untuk collection
- **GET** `/v1/deposit-payment/balance` - Balance payment deposit berdasarkan cust ID

### VAT Extract
- **GET** `/v1/vat-extract` - Mendapatkan daftar VAT extract
- **GET** `/v1/vat-extract/result` - Daftar hasil VAT extract
- **GET** `/v1/vat-extract/result/:vat_extract_id` - Detail hasil VAT extract
- **GET** `/v1/vat-extract/download-result/:vat_extract_id` - Download hasil VAT extract
- **POST** `/v1/vat-extract` - Extract VAT

### CoreTax VAT Extract
- **GET** `/v1/coretax-vat-extract` - Mendapatkan daftar CoreTax VAT extract
- **GET** `/v1/coretax-vat-extract/:coretax_vat_extract_id` - Mendapatkan detail CoreTax VAT extract
- **POST** `/v1/coretax-vat-extract` - Membuat CoreTax VAT extract baru
- **PATCH** `/v1/coretax-vat-extract/:coretax_vat_extract_id` - Update CoreTax VAT extract
- **DELETE** `/v1/coretax-vat-extract/:coretax_vat_extract_id` - Hapus CoreTax VAT extract

### Taxes
- **GET** `/v1/taxes/` - Mendapatkan laporan taxes
- **GET** `/v1/taxes/generate` - Daftar generate taxes
- **POST** `/v1/taxes/generate` - Generate taxes
- **POST** `/v1/taxes/bulk-delete` - Hapus taxes secara bulk
- **DELETE** `/v1/taxes/:taxes_id` - Hapus taxes

### OPEX Transaction
- **GET** `/v1/opex-tr` - Mendapatkan daftar OPEX transaction
- **GET** `/v1/opex-tr/:opex_tr_no` - Mendapatkan detail OPEX transaction
- **POST** `/v1/opex-tr` - Membuat OPEX transaction baru
- **PATCH** `/v1/opex-tr/:opex_tr_no` - Update OPEX transaction
- **DELETE** `/v1/opex-tr/:opex_tr_no` - Hapus OPEX transaction

### Memo JR
- **GET** `/v1/memo-jr` - Mendapatkan daftar memo JR
- **GET** `/v1/memo-jr/:mj_no` - Mendapatkan detail memo JR
- **POST** `/v1/memo-jr` - Membuat memo JR baru
- **PATCH** `/v1/memo-jr/:mj_no` - Update memo JR
- **DELETE** `/v1/memo-jr/:mj_no` - Hapus memo JR

### Master Data Finance
Semua endpoint berikut mengikuti pola CRUD standar:

- **M AP Disc** - `/v1/m-ap-disc`
- **M Cheque Reject** - `/v1/cheque-reject`
- **M COA** - `/v1/mcoa`
- **M COA Type** - `/v1/mcoa-type`
- **M OPEX** - `/v1/opex`
- **M Taxes** - `/v1/m-taxes`
- **CNDN** - `/v1/cndns`

### Files
- **POST** `/v1/files/uploads` - Upload file

---

## 🛒 Sales Service

### Order
- **GET** `/v1/orders` - Mendapatkan daftar order
- **GET** `/v1/orders/:ro_no` - Mendapatkan detail order
- **GET** `/v1/orders/discount` - Mendapatkan detail discount untuk order
- **GET** `/v1/orders/minimum-price/:pro_id/product` - Mendapatkan minimum price produk
- **POST** `/v1/orders` - Membuat order baru
- **POST** `/v1/orders/conversion` - Konversi unit produk
- **PATCH** `/v1/orders/:ro_no` - Update order
- **PATCH** `/v1/orders/final/:ro_no` - Update final order
- **PATCH** `/v1/orders/status` - Update status order secara bulk
- **DELETE** `/v1/orders/:ro_no` - Hapus order

### Outlet Lookup (dari Order Controller)
- **GET** `/v1/outlets` - Lookup salesman untuk outlet

### RO (Request Order)
- **GET** `/v1/ro` - Mendapatkan daftar RO
- **GET** `/v1/ro/:ro_no` - Mendapatkan detail RO
- **POST** `/v1/ro` - Membuat RO baru
- **PATCH** `/v1/ro/:ro_no` - Update RO
- **DELETE** `/v1/ro/:ro_no` - Hapus RO

### SO (Sales Order)
- **GET** `/v1/so` - Mendapatkan daftar SO
- **GET** `/v1/so/:so_no` - Mendapatkan detail SO
- **POST** `/v1/so` - Membuat SO baru
- **PATCH** `/v1/so/:so_no` - Update SO
- **DELETE** `/v1/so/:so_no` - Hapus SO

### Invoice
- **GET** `/v1/invoices` - Mendapatkan daftar invoice
- **GET** `/v1/invoices/details` - Mendapatkan detail invoice
- **GET** `/v1/invoices/:ro_no` - Mendapatkan invoice berdasarkan RO
- **POST** `/v1/invoices/` - Update invoice
- **PATCH** `/v1/invoices/print/:invoice_no` - Print invoice

### Return
- **GET** `/v1/returns` - Mendapatkan daftar return
- **GET** `/v1/returns/:return_no` - Mendapatkan detail return
- **GET** `/v1/returns/filter/outlets` - Filter outlet untuk return
- **GET** `/v1/returns/filter/salesmans` - Filter salesman untuk return
- **GET** `/v1/returns/filter/employees` - Filter employee untuk return
- **GET** `/v1/returns/filter/roles` - Filter role untuk return
- **GET** `/v1/returns/filter/return-statuses` - Filter return status
- **GET** `/v1/returns/create/filter/outlets` - Filter outlet untuk create return
- **GET** `/v1/returns/create/filter/salesmans` - Filter salesman untuk create return
- **GET** `/v1/returns/create/products` - Daftar produk untuk create return
- **GET** `/v1/returns/master/warehouses` - Daftar warehouse master
- **GET** `/v1/returns/master/return-reasons` - Daftar return reason master
- **GET** `/v1/returns/master/products` - Daftar produk master
- **POST** `/v1/returns` - Membuat return baru
- **PATCH** `/v1/returns/:return_no` - Update return
- **DELETE** `/v1/returns/:return_no` - Hapus return

### Promotion
- **GET** `/v1/promotions` - Mendapatkan daftar promotion
- **GET** `/v1/promotions/:promo_id` - Mendapatkan detail promotion
- **GET** `/v1/promotions/statuses` - Daftar status promotion
- **POST** `/v1/promotions` - Membuat promotion baru
- **POST** `/v1/promotions/consult` - Consult promotion
- **POST** `/v1/promotions/bulk-update-status` - Update status promotion secara bulk
- **PATCH** `/v1/promotions/:promo_id` - Update promotion
- **DELETE** `/v1/promotions/:promo_id` - Hapus promotion

### Promotion V2
- **POST** `/v2/promotions` - Membuat promotion baru (v2)
- **PATCH** `/v2/promotions/:promo_id` - Update promotion (v2)

### Promo Template
- **GET** `/v1/promo-templates` - Mendapatkan daftar promo template
- **GET** `/v1/promo-templates/:promo_template_id` - Mendapatkan detail promo template
- **GET** `/v1/promo-templates/statuses` - Daftar status promo template
- **POST** `/v1/promo-templates` - Membuat promo template baru
- **PATCH** `/v1/promo-templates/:promo_template_id` - Update promo template
- **DELETE** `/v1/promo-templates/:promo_template_id` - Hapus promo template

### Discount
- **GET** `/v1/discounts` - Mendapatkan daftar discount
- **GET** `/v1/discounts/:discount_id` - Mendapatkan detail discount
- **GET** `/v1/discounts/statuses` - Daftar status discount
- **GET** `/v1/discounts/publish/statuses` - Daftar publish status discount
- **POST** `/v1/discounts` - Membuat discount baru
- **POST** `/v1/discounts/publish` - Publish discount
- **POST** `/v1/discounts/consult` - Consult discount
- **PATCH** `/v1/discounts/:discount_id` - Update discount
- **DELETE** `/v1/discounts/:discount_id` - Hapus discount

### Validate Order
- **POST** `/v1/validate-order/` - Validasi order
- **POST** `/v1/validate-order/detail` - Validasi detail order

### Consignment
- **GET** `/v1/consignment` - Mendapatkan daftar consignment
- **GET** `/v1/consignment/:cons_no` - Mendapatkan detail consignment
- **POST** `/v1/consignment` - Membuat consignment baru
- **PATCH** `/v1/consignment/:cons_no` - Update consignment
- **DELETE** `/v1/consignment/:cons_no` - Hapus consignment

### TLS (Temporary Loan Stock)
- **GET** `/v1/tls` - Mendapatkan daftar TLS
- **GET** `/v1/tls/:tls_id` - Mendapatkan detail TLS
- **POST** `/v1/tls` - Membuat TLS baru
- **PATCH** `/v1/tls/:tls_id` - Update TLS
- **DELETE** `/v1/tls/:tls_id` - Hapus TLS

### Gamification
- **GET** `/v1/gamifications` - Mendapatkan daftar gamification
- **GET** `/v1/gamifications/:gamification_id` - Mendapatkan detail gamification
- **GET** `/v1/gamifications/filter/gamification-statuses` - Filter status gamification
- **GET** `/v1/gamifications/master/job-titles` - Daftar job title master
- **GET** `/v1/gamifications/master/sales-teams` - Daftar sales team master
- **GET** `/v1/gamifications/master/operation-types` - Daftar operation type master
- **GET** `/v1/gamifications/master/warehouses` - Daftar warehouse master
- **GET** `/v1/gamifications/master/salesmans` - Daftar salesman master
- **GET** `/v1/gamifications/master/product-categories` - Daftar product category master
- **GET** `/v1/gamifications/master/product-lines` - Daftar product line master
- **GET** `/v1/gamifications/master/brands` - Daftar brand master
- **GET** `/v1/gamifications/master/sub-brands1` - Daftar sub brand 1 master
- **GET** `/v1/gamifications/master/sub-brands2` - Daftar sub brand 2 master
- **POST** `/v1/gamifications` - Membuat gamification baru
- **PATCH** `/v1/gamifications/:gamification_id` - Update gamification
- **DELETE** `/v1/gamifications/:gamification_id` - Hapus gamification

### Order Approval
- **GET** `/v1/order-approval` - Mendapatkan daftar order approval
- **PATCH** `/v1/order-approval/:order_approval_request_id` - Update order approval

### Hierarchy Approval
- **GET** `/v1/hierarchy-approval/companies` - Daftar companies untuk hierarchy approval
- **GET** `/v1/hierarchy-approval` - Mendapatkan daftar hierarchy approval
- **GET** `/v1/hierarchy-approval/:hierarchy_approval_id` - Mendapatkan detail hierarchy approval
- **GET** `/v1/hierarchy-approval/request/:request_approval_id` - Mendapatkan request approval
- **POST** `/v1/hierarchy-approval` - Membuat hierarchy approval baru
- **PATCH** `/v1/hierarchy-approval/:hierarchy_approval_id` - Update hierarchy approval
- **DELETE** `/v1/hierarchy-approval/:hierarchy_approval_id` - Hapus hierarchy approval

### Report
- **GET** `/v1/reports` - Mendapatkan daftar report
- **POST** `/v1/reports/secondary-sales` - Generate secondary sales report
- **GET** `/v1/reports/secondary-sales/sum-date` - Summary secondary sales per bulan
- **GET** `/v1/reports/secondary-sales/group` - Group secondary sales report
- **GET** `/v1/reports/secondary-sales/trend-sales` - Trend sales report
- **POST** `/v1/reports/activity-report-sales` - Generate activity report sales
- **GET** `/v1/reports/activity-report-sales` - Daftar activity report sales
- **GET** `/v1/reports/activity-report-sales/sum-date` - Summary activity report per bulan
- **GET** `/v1/reports/activity-report-sales/group` - Group activity report sales

### Files
- **POST** `/v1/files/uploads` - Upload file

---

## 📱 Mobile Service

### User
- **POST** `v1/users/login` - Login user (tanpa JWT)
- **POST** `v1/users/register` - Register user (tanpa JWT)
- **POST** `v1/users/forgot-password` - Lupa password (tanpa JWT)
- **POST** `v1/users/forgot-password-validate` - Validasi forgot password (tanpa JWT)
- **PATCH** `v1/users/password` - Update password (tanpa JWT)
- **PATCH** `v1/users/change-password` - Change password (tanpa JWT)
- **GET** `/v1/users/profile` - Mendapatkan profile user
- **PATCH** `/v1/users/:user_id` - Update user

### Order
- **GET** `/v1/orders` - Mendapatkan daftar order
- **GET** `/v1/orders/:ro_no` - Mendapatkan detail order
- **POST** `/v1/orders` - Membuat order baru
- **POST** `/v1/orders/no-order` - Membuat no order
- **GET** `/v1/orders/no-order` - Daftar no order
- **GET** `/v1/orders/sales-report` - Summary order per salesman
- **POST** `/v1/orders/conversion` - Konversi unit produk
- **PATCH** `/v1/orders/:ro_no` - Update order
- **PATCH** `/v1/orders/final/:ro_no` - Update final order
- **PATCH** `/v1/orders/status` - Update status order secara bulk
- **DELETE** `/v1/orders/:ro_no` - Hapus order

### Order Canvas
- **GET** `/v1/orders-canvas` - Mendapatkan daftar order canvas
- **GET** `/v1/orders-canvas/:ro_no` - Mendapatkan detail order canvas
- **POST** `/v1/orders-canvas` - Membuat order canvas baru
- **POST** `/v1/orders-canvas/no-order` - Membuat no order canvas
- **GET** `/v1/orders-canvas/no-order` - Daftar no order canvas
- **GET** `/v1/orders-canvas/sales-report` - Summary order canvas per salesman
- **POST** `/v1/orders-canvas/conversion` - Konversi unit produk
- **PATCH** `/v1/orders-canvas/:ro_no` - Update order canvas
- **PATCH** `/v1/orders-canvas/final/:ro_no` - Update final order canvas
- **PATCH** `/v1/orders-canvas/status` - Update status order canvas secara bulk
- **DELETE** `/v1/orders-canvas/:ro_no` - Hapus order canvas

### Order History
- **GET** `/v1/orders-history` - Mendapatkan daftar order history
- **GET** `/v1/orders-history/:ro_no` - Mendapatkan detail order history

### Outlet Lookup (dari Order Controller)
- **GET** `/v1/outlets` - Lookup salesman untuk outlet

### Product
- **GET** `/v1/products` - Mendapatkan daftar produk
- **GET** `/v1/products/:pro_id` - Mendapatkan detail produk

### Stock
- **GET** `/v1/stocks/gudang-utama` - Daftar stock gudang utama
- **GET** `/v1/stocks/gudang-canvas` - Daftar stock gudang canvas

### Return
- **GET** `/v1/returns` - Mendapatkan daftar return
- **GET** `/v1/returns/master/return-reasons` - Daftar return reason master
- **POST** `/v1/returns` - Membuat return baru
- **POST** `/v1/returns/status` - Update status return
- **PATCH** `/v1/returns/quantity/:return_no` - Update quantity return

### Return Reasons
- **GET** `/v1/return-reasons/` - Daftar return reasons

### Promotion
- **GET** `/v1/promotions` - Mendapatkan daftar promotion
- **POST** `/v1/promotions/consult` - Consult promotion

### Discount
- **GET** `/v1/discounts` - Mendapatkan daftar discount
- **GET** `/v1/discounts/:discount_id` - Mendapatkan detail discount
- **GET** `/v1/discounts/statuses` - Daftar status discount
- **GET** `/v1/discounts/publish/statuses` - Daftar publish status discount
- **POST** `/v1/discounts` - Membuat discount baru
- **POST** `/v1/discounts/publish` - Publish discount
- **POST** `/v1/discounts/consult` - Consult discount
- **PATCH** `/v1/discounts/:discount_id` - Update discount
- **DELETE** `/v1/discounts/:discount_id` - Hapus discount

### Validate Order
- **POST** `/v1/validate-order/` - Validasi order
- **POST** `/v1/validate-order/detail` - Validasi detail order

### Invoice
- **GET** `/v1/invoices/payment/:invoice_no` - Mendapatkan payment invoice
- **POST** `/v1/invoices/payment` - Membuat payment invoice

### Collection
- **GET** `/v1/collections` - Mendapatkan daftar collection
- **GET** `/v1/collections/deposit/:deposit_no` - Mendapatkan detail collection berdasarkan deposit
- **GET** `/v1/collections/invoice/:invoice_no` - Mendapatkan detail collection berdasarkan invoice
- **GET** `/v1/collections/missed-payment-reasons` - Daftar missed payment reason
- **POST** `/v1/collections` - Membuat collection baru
- **POST** `/v1/collections/no-payment` - Membuat collection tanpa payment

### Sales
- **GET** `/v1/sales/summary` - Summary sales

### Visits
- **GET** `/v1/visits/` - Mendapatkan daftar visits
- **GET** `/v1/visits/summary` - Summary visits
- **GET** `/v1/visits/list` - List visits
- **GET** `/v1/visits/skip/reasons` - Daftar skip reason
- **POST** `/v1/visits/start` - Start visit
- **POST** `/v1/visits/skip` - Skip visit
- **POST** `/v1/visits/Arrive` - Arrive visit
- **POST** `/v1/visits/Hold` - Hold visit
- **POST** `/v1/visits/Resume` - Resume visit
- **POST** `/v1/visits/Leave` - Leave visit
- **POST** `/v1/visits/End` - End visit

### Activities
- **GET** `/v1/activities/summary/daily` - Summary activities harian

### Announcements
- **GET** `/v1/announcements/` - Mendapatkan daftar announcements

### Events
- **GET** `/v1/events/` - Mendapatkan daftar events

### Leaderboards
- **GET** `/v1/leaderboards/` - Mendapatkan daftar leaderboards

### Attendances
- **GET** `/v1/attendances` - Mendapatkan daftar attendance
- **POST** `/v1/attendances` - Membuat attendance baru

### Employee
- **GET** `/v1/employees` - Mendapatkan daftar employee
- **GET** `/v1/employees/:emp_id` - Mendapatkan detail employee
- **POST** `/v1/employees` - Membuat employee baru
- **POST** `/v1/employees/create-multiple` - Membuat multiple employee
- **PATCH** `/v1/employees/:emp_id` - Update employee
- **DELETE** `/v1/employees/:emp_id` - Hapus employee

### Employee Group
- **GET** `/v1/employee-groups` - Mendapatkan daftar employee group
- **GET** `/v1/employee-groups/:emp_group_id` - Mendapatkan detail employee group
- **POST** `/v1/employee-groups` - Membuat employee group baru
- **PATCH** `/v1/employee-groups/:emp_group_id` - Update employee group
- **DELETE** `/v1/employee-groups/:emp_group_id` - Hapus employee group

### M Outlet
- **GET** `/v1/m-outlets` - Mendapatkan daftar m-outlet
- **GET** `/v1/m-outlets/:outlet_id` - Mendapatkan detail m-outlet
- **POST** `/v1/m-outlets` - Membuat m-outlet baru

### Taking Order
- **GET** `/v1/no-order-reasons` - Daftar no order reason

### Pickup Reason
- **GET** `/v1/pickup-reasons` - Daftar pickup reason

### Files
- **POST** `/v1/files/uploads` - Upload file

---

## ⚙️ System Service

### User
- **POST** `v1/users/login` - Login user (tanpa JWT)
- **POST** `v1/users/forgot-password` - Lupa password (tanpa JWT)
- **POST** `v1/users/forgot-password-validate` - Validasi forgot password (tanpa JWT)
- **PATCH** `v1/users/password` - Update password (tanpa JWT)
- **GET** `/v1/users-menu/menus/:menu_param` - Mendapatkan semua menu (tanpa JWT)
- **GET** `/v1/users` - Mendapatkan daftar user
- **GET** `/v1/users/:user_id` - Mendapatkan detail user
- **GET** `/v1/users/menus/:menu_param` - Mendapatkan menu user (web/desktop)
- **GET** `/v1/cust` - Mendapatkan detail customer
- **POST** `/v1/users` - Membuat user baru
- **PATCH** `/v1/users/:user_id` - Update user
- **DELETE** `/v1/users/:user_id` - Hapus user

### Config
- **GET** `/v1/config` - Mendapatkan daftar config
- **GET** `/v1/config/:config_id` - Mendapatkan detail config
- **GET** `/v1/config/details/list` - Daftar detail config
- **POST** `/v1/config` - Membuat config baru
- **PATCH** `/v1/config/:config_id` - Update config
- **DELETE** `/v1/config/:config_id` - Hapus config

### M Day
- **GET** `/v1/m-days` - Mendapatkan daftar m-day
- **GET** `/v1/m-days/:day_id` - Mendapatkan detail m-day

### Notification
- **POST** `v1/notifications/whatsapp-cicd` - Kirim notifikasi WhatsApp untuk CI/CD (tanpa JWT)

### Files
- **POST** `/v1/files/uploads` - Upload file

---

## ⏰ Cronjob Service

### Job
- **GET** `/v1/jobs` - Mendapatkan daftar job
- **GET** `/v1/jobs/:job_id` - Mendapatkan detail job
- **POST** `/v1/jobs` - Membuat job baru
- **DELETE** `/v1/jobs/:job_id` - Hapus job

---

## 📝 Catatan Penting

1. **Authentication**: Sebagian besar endpoint memerlukan JWT token yang dikirim melalui header `Authorization: Bearer {token}`. Endpoint yang tidak memerlukan JWT telah ditandai dengan "(tanpa JWT)".

2. **Pagination**: Endpoint yang mengembalikan list biasanya mendukung pagination melalui query parameters:
   - `page`: Nomor halaman (default: 1)
   - `limit`: Jumlah data per halaman (default: 10)
   - `sort`: Sorting (format: `field:ASC` atau `field:DESC`)
   - `active`: Filter active status (0=all, 1=active, 2=inactive)
   - `query`: Search query (optional)

3. **Response Format**: Semua response mengikuti format standar:
   ```json
   {
     "message": "Success",
     "data": {...},
     "errors": null,
     "paging": {...},
     "request_id": "..."
   }
   ```

4. **Error Handling**: Error response mengikuti format yang sama dengan field `errors` berisi array error detail.

5. **Content Type**: Semua request dan response menggunakan `Content-Type: application/json`.

---

**Last Updated**: 2024

