# Runtime cURL + DB Validation — SX-2182

Task id: `20260608-1534-sx-2182-secondary-sales-multiselect`
Validated at: `2026-06-08T17:09+07:00` onward

## Runtime setup

Docker local services:

```bash
rtk docker compose -f docker-compose.yml up -d rabbitmq sales
rtk docker compose -f docker-compose.yml ps
```

Observed running services:

- `scylla-system` on `9001`
- `scylla-master` on `9002`
- `scylla-sales` on `9004`
- `scylla-rabbitmq` healthy
- `scylla-redis` healthy

Ping checks:

- `GET http://127.0.0.1:9002/ping` → `200`
- `GET http://127.0.0.1:9004/ping` → `200`

Local DB:

```bash
PGPASSWORD=postgres psql -h 127.0.0.1 -p 5432 -U postgres -d ggn_scyllax
```

DB check result: connected to `ggn_scyllax` as `postgres`.

Auth users:

- Principal login: `princ@idetama.id` / `admin` → HTTP 200, token saved outside repo under temp path.
- Distributor login: `dist@sda.idetama.id` / `admin` → HTTP 200, token saved outside repo under temp path.

Decoded claims used for validation:

- Principal:
  - `cust_id = C22001`
  - `parent_cust_id = C22001`
  - `employee_id = 278`
  - `distributor_id = 0`
- Distributor:
  - `cust_id = C220010001`
  - `parent_cust_id = C22001`
  - `distributor_id = 67`

## Master Business Unit cURL validation

Output files were saved under `/var/folders/2r/vlp3bhfn1cl5cpdhp23rkh440000gn/T/opencode/sx2182-curl/`.

### Single region/area

```bash
GET http://127.0.0.1:9002/v1/business-unit?is_active=1&region_id=67&area_id=82&q=&page=1&limit=99
```

Result: `200`

Response summary:

- `paging.total_record = 3`
- `distributor_data` length = `3`

DB validation:

```sql
SELECT COUNT(DISTINCT md.distributor_id)
FROM mst.m_distributor md
WHERE md.is_del=false
  AND md.parent_cust_id='C22001'
  AND md.region_id IN (67)
  AND md.area_id IN (82)
  AND md.is_active=true;
```

Result: `3`.

### Multi region/area comma

```bash
GET http://127.0.0.1:9002/v1/business-unit?is_active=1&region_id=67,68&area_id=82,70&q=&page=1&limit=99
```

Result: `200`

Response summary:

- `paging.total_record = 7`
- `distributor_data` length = `7`

DB validation result for equivalent filter: `7`.

### Multi region/area comma + spaces

```bash
GET http://127.0.0.1:9002/v1/business-unit?is_active=1&region_id=67,%2068&area_id=82,%2070&q=&page=1&limit=99
```

Result: `200`

Response summary:

- `paging.total_record = 7`
- `distributor_data` length = `7`

### Invalid numeric token

```bash
GET http://127.0.0.1:9002/v1/business-unit?is_active=1&region_id=67,abc&area_id=82&q=&page=1&limit=99
```

Result: `400`

Response summary:

```json
{
  "message": "invalid region_id value \"abc\"",
  "request_id": "..."
}
```

## Sales dashboard cURL validation

Dataset chosen from DB facts:

- `cust_id`: `C220010001`, `C220010002`
- `year`: `2026`
- `month`: `4`

DB fact sample showed both cust IDs have `report.fact_orders` data in April 2026.

### Trend single cust

```bash
GET http://127.0.0.1:9004/v1/reports/secondary-sales/trend-sales?year=2026&cust_id=C220010001
```

Result: `200`

Response summary:

- rows = `12`
- month 4:
  - `total_gross_sale = 9354000`
  - `total_discount_promo = 0`
  - `net_sales = 9354000`

### Trend multi cust

```bash
GET http://127.0.0.1:9004/v1/reports/secondary-sales/trend-sales?year=2026&cust_id=C220010001,C220010002
```

Result: `200`

Response summary:

- rows = `12`
- month 4:
  - `total_gross_sale = 23604000`
  - `total_discount_promo = 120000`
  - `net_sales = 23484000`

DB validation:

```sql
SELECT fo.cust_id, dt.month,
  SUM(fo.gross_sale),
  SUM(fo.discount + fo.special_discount),
  SUM(fo.net_sales_exclude_ppn)
FROM report.fact_orders fo
JOIN report.dim_dates dt ON dt.id=fo.date_id
WHERE dt.year=2026
  AND dt.month=4
  AND fo.cust_id IN ('C220010001','C220010002')
GROUP BY fo.cust_id, dt.month;
```

DB results:

- `C220010001`: gross `9354000`, discount `0`, net `9354000`
- `C220010002`: gross `14250000`, discount `120000`, net `14130000`
- Multi aggregate: gross `23604000`, discount `120000`, net `23484000`

Matches API response.

### Sum-date single cust

```bash
GET http://127.0.0.1:9004/v1/reports/secondary-sales/sum-date?month=4&year=2026&cust_id=C220010001
```

Result: `200`

Response summary:

- `total_gross_sale = 6216544`
- `total_discount_promo = 0`
- `net_sales_exc_ppn = 6341344`
- `net_sales = 6657924`
- `qty = 217`
- `qty_return = 70`
- `return_rate = 32.25806451612903`

DB validation for equivalent order-minus-return formula matched those values.

### Sum-date multi cust

```bash
GET http://127.0.0.1:9004/v1/reports/secondary-sales/sum-date?month=4&year=2026&cust_id=C220010001,C220010002
```

Result: `200`

Response summary:

- `total_gross_sale = 20166544`
- `total_discount_promo = 120000`
- `net_sales_exc_ppn = 20171344`
- `net_sales = 21101724`
- `qty = 250`
- `qty_return = 270`
- `return_rate = 108`

DB validation for equivalent multi-cust order-minus-return formula:

- gross `20166544`
- discount `120000`
- net excl ppn `20171344`

Matches API response.

### Group outlet multi cust

```bash
GET http://127.0.0.1:9004/v1/reports/secondary-sales/group?month=4&year=2026&cust_id=C220010001,C220010002&group_by=outlet
```

Result: `200`

Response summary:

- rows = `49`
- first row:

```json
{
  "id": 244,
  "code": "B00003",
  "name": "Bangka III",
  "net_sales": 7950000
}
```

DB validation equivalent group query top rows:

- `244 | B00003 | Bangka III | 7950000.00`
- `1660 | IDEA260049 | Toko Pelan Tapi Sampai | 3750000.00`
- `1462 | 0019 | Babeh TK | 2762000.00`

Matches API first row and confirms `code` is populated.

### Group salesman multi cust

```bash
GET http://127.0.0.1:9004/v1/reports/secondary-sales/group?month=4&year=2026&cust_id=C220010001,C220010002&group_by=salesman
```

Result: `200`

Response summary:

- rows = `11`
- first row:

```json
{
  "id": 339,
  "code": "NS000001",
  "name": "New Salesman Test",
  "net_sales": 13830000
}
```

### Unauthorized distributor sibling

Distributor token (`cust_id=C220010001`) requesting sibling `C220010002`:

```bash
GET http://127.0.0.1:9004/v1/reports/secondary-sales/sum-date?month=4&year=2026&cust_id=C220010002
```

Result: `403`

Response summary:

```json
{
  "message": "cust_id is outside authorized scope",
  "request_id": "..."
}
```

## Export cURL validation

Date range used:

- `from = 1774976400` (`2026-04-01 00:00:00 Asia/Jakarta`)
- `to = 1777568399` (`2026-04-30 23:59:59 Asia/Jakarta`)

### Export legacy string cust_id

```bash
POST http://127.0.0.1:9004/v1/reports/secondary-sales
{
  "from": 1774976400,
  "to": 1777568399,
  "distributor_ids": [],
  "outlet_ids": [],
  "salesman_ids": [],
  "pro_ids": [],
  "cust_id": "C220010001"
}
```

Result: `200`

Response summary:

- `report_id = 6a26961d8161651c86508278`
- `report_name = SecondarySales-080626-001`
- initial response `file_status = 2`

DB `report.list` after async processing:

- `cust_id = C22001` (auth owner preserved)
- `file_status = 1`
- `file_url` populated
- date range `2026-04-01` to `2026-04-30`

### Export array multi cust_id

```bash
POST http://127.0.0.1:9004/v1/reports/secondary-sales
{
  "from": 1774976400,
  "to": 1777568399,
  "distributor_ids": [],
  "outlet_ids": [],
  "salesman_ids": [],
  "pro_ids": [],
  "cust_id": ["C220010001", "C220010002"]
}
```

Result: `200`

Response summary:

- `report_id = 6a26961e8161651c86508279`
- `report_name = SecondarySales-080626-002`
- initial response `file_status = 2`

DB `report.list` after async processing:

- `cust_id = C22001` (auth owner preserved)
- `file_status = 1`
- `file_url` populated
- date range `2026-04-01` to `2026-04-30`

DB export source row validation for multi-cust/date range:

- order rows:
  - `C220010001 = 216`
  - `C220010002 = 4`
- return rows:
  - `C220010001 = 4`

This confirms the export source has rows for both selected child cust IDs and the generated report job completed locally.

## Notes / limitations

- Runtime validation used existing local data in `ggn_scyllax`; no synthetic data was inserted for SX-2182.
- Token files were written only to temp path outside repo and not committed.
- Export file content was not downloaded/opened; DB `report.list` status and source-row SQL were validated instead.
