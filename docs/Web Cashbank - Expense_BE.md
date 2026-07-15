1. # **Expense List**   

Create New endpoint   
URL : 

| {{url}}/finance/v1/expense |
| :---- |

Method		: GET

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| page | Integer | Yes | default : 1 |
| limit  | Integer | Yes | default : 20 |
| sort | String | Yes | default created\_date:desc  |
| start\_date |  | No | epoch (filter minimal tanggal) |
| end\_date |  | No | epoch (filter maksimal tanggal) |
| min\_balance | decimal  | Yes | filter balance \> nilai |
| collector\_id | Array(INT)  | NO  | filter berdasarkan acf.deposit.collector\_id  |

*Jika start\_date dan end\_date tidak dikirim, data difilter otomatis ke 3 bulan terakhir. dan BE filter berdasarkan acf.expense.created\_by berdasrkan user\_id yang login* 

### **Example Request** 

---

**Example Request Default (3 bulan terakhir)** 

| `curl -X GET "{{url}}/acf/v1/expense?page=1&limit=10&sort=created_date:desc" \   -H "Authorization: Bearer $TOKEN" \   -H "Accept: application/json" | jq '.'` |
| :---- |

**Dengan filter tanggal:**

| `curl -X GET "{{url}}/acf/v1/expense?page=1&limit=10&sort=created_date:desc&start_date=1730419200&end_date=1738262399" \   -H "Authorization: Bearer $TOKEN" \   -H "Accept: application/json" | jq '.'` |
| :---- |

* jika filter berdasarkan distributor :   
  **select** \* **from** acf.expense *e* **where** *e*.cust\_id \= ( **select** *mc*.cust\_id  **from** smc.m\_customer *mc*  **where**  *mc*.distributor\_id \= 67 )  
* jika filter berdasarkan user principal itu sendiri (cust\_id)  
  **select \* from acf.expense *e* where *e*.cust\_id \= 'C22001'**

### **Response  :** *add filter data dengan data 3 bulan terakhir* 
---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| expense\_id | Int | 8 | Id expense acf.expense.expense\_id |
| doc\_no | varchar | 50 | acf.expense.doc\_no |
| date | Date  |  | Tanggal create expense contoh : DD/MM/YYYY acf.expense.date |
| expense\_type\_id | Int | 8 | acf.expense.expense\_type\_id |
| expense\_type\_code | varchar | 20 | acf.expense\_type.expense\_type\_code |
| expense\_type\_name | Varchar | 50 | acf.expense\_type.expense\_name |
| collector\_id | int | 8 | acf.expense.collector\_id |
| collector\_name | Varchar | 150 | sys.m\_user.user\_name |
| balance | Numeric | 11 | Sisa acf.expense.balance |
| amount | Numeric | 11 | Nominal expense acf.expense.amount |
| reason | Varchar | 255 | reason  acf.expense.note |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data `{   "message": "",   "data": [     {       "expense_id": 13,       "document_no": "EXP/2025/12/0001",       "date": "01/12/2025",       "expense_type_id": 1,       "expense_type_code": "CLEAN",       "expense_type_name": "Uang Kebersihan",       "collector_id": 123,       "collector_name": "Budi",       "balance": 12000,       "amount": 12000,       "reason": "",       "is_clock_out": 1     },     {       "expense_id": 14,       "document_no": "EXP/2025/12/0002",       "date": "01/12/2025",       "expense_type_id": 1,       "expense_type_code": "CLEAN",       "expense_type_name": "Uang Kebersihan",       "collector_id": 123,       "collector_name": "Budi",       "balance": 12000,       "amount": 12000,       "reason": "",       "is_clock_out": 1     }   ],   "paging": {     "total_record": 84,     "page_current": 1,     "page_limit": 10,     "page_total": 9   },   "request_id": "6915a5e8e3f53f84fe73517f" }`  |
| :---- |
| **Case : empty state / tidak terdapat data berdasarkan pencarian**  `{   "message": "No Data",   "data": null,   "paging": {     "total_record": 0,     "page_current": 1,     "page_limit": 10,     "page_total": 0   },   "request_id": "6915a6dd2395083c685e8e16" }`  |
| **Case : error ** `{   "message": "record not found",   "request_id": "6915a6dd2395083c685e8e16" }` |

2. # **Expense Detail** 

Create New endpoint   
URL : 

| {{url}}/finance/v1/expense/{expense\_id} |
| :---- |

Method		: GET

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Path Variable** 

---

|  expense\_id | 32 |
| :---- | :---- |

### **Example Request** 

---

**Example Request Default**

| bash `curl -X GET "{{url}}/acf/v1/expense/{expense_id}" \   -H "Authorization: Bearer $TOKEN" \   -H "Accept: application/json" | jq '.'` |
| :---- |

### **Response** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Object | \- |  \- |
| expense\_id  | Int | 8 | Id expense acf.expense.expense\_id |
| date | Date  |  | Tanggal create expense contoh : DD/MM/YYYY acf.expense.date |
| doc\_no | varchar | 50 | acf.expense.doc\_no |
| expense\_type\_id | Int | 8 | acf.expense.expense\_type\_id |
| expense\_type\_code | varchar | 20 | acf.expense\_type.expense\_type\_code |
| expense\_type\_name | Varchar | 50 | acf.expense\_type.expense\_name |
| collector\_id | int | 8 | acf.expense.collector\_id |
| collector\_name | Varchar | 150 | sys.m\_user.user\_name |
| amount | Numeric | 11 | Nominal expense acf.expense.amount |
| note | Varchar | 255 | reason  acf.expense.note |
| remainig\_amount  | numeric | 20,4 | acf.expense.amount \- (SUM(acf.deposit\_expense.payment\_amount) |
| **files**  | **Array**  |  |  |
| file\_name | varchar | 255 | Nama File  |
| file\_type | varchar | 50 | photo disimpan dalam format JPG |
| media\_category | ENUM('image','video') |  | Membedakan jenis konten: *image* atau *video* |
| file\_url | text |  | url file |
| file\_size | BIGINT |  | Untuk mengetahui ukuran file asli |
| request\_id | String | 150 | Generate request id  |
| **deposits**  | **Array**  |  |  |
| deposit\_expense\_id | int | 8 | acf.deposit\_expense.deposit\_expense\_id |
| deposit\_id | int | 8 | acf.deposit\_expense.deposit\_id |
| used\_amount | numeric | 20,4 | acf.deposit\_expense.payment\_amount |
| deposit\_no | varchar | 30 | acf.deposit.deposit\_no |
| update\_date | timestamp | 6 | acf.deposit\_expense.create\_date or update\_date |

**Example Response : (sukses)**

| {   "message": "",   "data": {     "expense\_id": 13,     "date": "2025-12-01",     "doc\_no": "EXP/2025/12/0001",     "expense\_type\_id": 1,     "expense\_type\_code": "CLEAN",     "expense\_type\_name": "Uang Kebersihan",     "collector\_id": 123,     "collector\_name": "Budi",     "amount": 12000,     "note": "",     "files": \[       {         "file\_name": "arrival\_IMG\_20250119\_155922",         "file\_type": "JPG",         "media\_category": "image",         "file\_url": "https://example.com/users/123/avatar/profile\_123.png",         "file\_size": 204800       }     \],     "deposits": \[       {         "deposit\_expense\_id": 1,         "deposit\_id": 10,         "used\_amount": 5000,         "deposit\_no": "DEP/2025/01/0001",         "update\_date": "2026-01-23T10:30:00Z"       }     \]   },   "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |

**Example Response : (error)**

| `{   "message": "record not found",   "request_id": "6915a6dd2395083c685e8e16" }` |
| :---- |

3. #  **Create Expense** 

| {{url}}/finance/v1/expense |
| :---- |

### **Method : POST**

### **Headers :**

---

| Accept | application/json |
| :---- | :---- |
| Authorization  | Bearer Token |

Upload file menggunakan Object Storage Huawei 

### **Body** 

| Field | Type | Required | Description |
| :---- | :---- | :---- | :---- |
| expense\_type\_id | integer | yes |  |
| amount | integer | yes |  |
| collector\_id | integer | yes |  |
| note | varchar | No |  |
| file\_url | array(text)  | yes | file yang akan di upload (max 3\)  |

### **Example Request** 

| curl \--location '{{url}}/acf/v1/expense' \\ \--header 'Accept: application/json' \\ \--header 'Authorization: Bearer YOUR\_TOKEN' \\ \--header 'Content-Type: application/json' \\ \--data '{   "expense\_type\_id": 1,   "amount": 50000,   "collector\_id": 10,   "note": "Biaya parkir",   "file\_url": \[     "https://obs.huawei.com/expense/file1.jpg",     "https://obs.huawei.com/expense/file2.jpg"   \] }'  |
| :---- |

### **Response :** 

| {   "message": "Expense created successfully",   "data": {     "expense\_id": 15,     "document\_no": "E20261255111"   },   "request\_id": "7015a5e8e3f53f84fe73517f" }  |
| :---- |

1\) Data outlet not found

| {   "message": "Data outlet not found",   "errors": \[     {       "field": "outlet\_id",       "description": "outlet\_id tidak ditemukan atau tidak aktif"     }   \],   "request\_id": "7015a6dd2395083c685e8e16" } |
| :---- |

2\) Expense type not found

| `{   "message": "Expense type not found",   "errors": [     {       "field": "expense_type_id",       "description": "expense_type_id tidak ditemukan"     }   ],   "request_id": "7015a6dd2395083c685e8e17" }` |
| :---- |

3\) File must not exceed 3

| `{   "message": "Validation error",   "errors": [     {       "field": "file_url",       "description": "Jumlah file tidak boleh lebih dari 3"     }   ],   "request_id": "7015a6dd2395083c685e8e18" }`  |
| :---- |

### **Impact ke Database   :**

- generate doc\_no \= EYYYYMMDD-3digitRunnningNumber   
  contoh : E20261201222  
- Insert ke acf.expense

#### 

| Nama Kolom | Value |
| :---- | :---- |
| cust\_id\*\* | diisi cust\_id user login |
| expense\_id\* | generate BE |
| doc\_no | generate doc\_no |
| expense\_type\_id\*\* | expense\_type\_id |
| source | 1 |
| date | diisi tgl sekarang |
| amount | amount |
| note | note |
| collector\_id\*\* | collector\_id |
| created\_by\*\* | diisi user yang input |
| created\_at | diisi waktu input data  |
| updated\_by\*\* | NULL |
| updated\_at | NULL |
| deleted\_by\*\* | NULL |
| deleted\_at | NULL |
| is\_del | False |

## 


- Insert detail file ke tabel file  acf.expense\_file

#### 

| Nama Kolom | Tipe Data |
| :---- | :---- |
| cust\_id\*\* | diisi cust\_id user login |
| expense\_file\_id | generate BE |
| expense\_id\* | expense\_id |
| file\_name | diisi file name |
| file\_url | diisi file\_url  |
| file\_key | image |
| media\_category | png / jpg/jpeg |
| file\_size | ukuran file |

4. #  **Update Expense**  

| {{url}}/finance/v1/expense/{expense\_id} |
| :---- |

### **Method : PATCH**

### **Headers**

---

| Accept | application /json  |
| :---- | :---- |
| Authorization  | Bearer Token |

**Path Variable** 

---

|  expense\_id | 32 |
| :---- | :---- |

### **Body** 

| Field | Type | Required | Description |
| :---- | :---- | :---- | :---- |
| amount | integer | yes |  |
| note | varchar | yes |  |

### **Example Request :**

| bash `curl -X PATCH "{{url}}/acf/v1/expense/{expense_id}" \   -H "Authorization: Bearer $TOKEN" \   -H "Accept: application/json" \   -H "Content-Type: application/json" \   -d '{     "amount": 15000,     "note": "Revisi nominal expense"   }'`  |
| :---- |

### 

### **Response Success:**

| `{   "message": "Expense updated successfully",   "data": {     "expense_id": 13,     "amount": 15000,     "note": "Revisi nominal expense"   },   "request_id": "7015a6dd2395083c685e8e16" }`  |
| :---- |

1\) Data expense not found

| `{   "message": "Data expense_id not found",   "errors": [     {       "field": "expense_id",       "description": "expense_id tidak ditemukan atau tidak aktif"     }   ],   "request_id": "7015a7aa2395083c685e8e19" }`  |
| :---- |

### **Impact ke Database**   

- Update ke acf.expense

