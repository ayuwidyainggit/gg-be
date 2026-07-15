1. # **Expense History** 

Create New endpoint   
URL : 

| {{url}}/mobile/v1/expense?limit=10 |
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
| start\_date |  |  | epoch  |
| end\_date |  |  | epoch |

### **Example Request** 

---

**Example Request Default**

| bash curl \-X GET "{{url}}/mobile/v1/expense?limit10" \\  \-H "Authorization: Bearer $TOKEN" \\  \-H "Content-Type: application/json" | jq '.' |
| :---- |

### **Response  :** *add filter data dengan data 3 bulan terakhir* 
---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| is\_clock\_out | int 1 |  | 1 : clock in 2 : clock outmobile.attendances.type  |
| expense\_id | Int | 8 | Id expense acf.expense.expense\_id |
| date | Date  |  | Tanggal create expense contoh : DD/MM/YYYY acf.expense.date |
| expense\_name | Varchar | 50 | Expense name   diambil dari acf.expense\_type.expense\_name |
| amount | Numeric | 11 | Nominal expense acf.expense.amount |
| reason | Varchar | 255 | reason  acf.expense.note |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data { "message": "", "data" : \[    "is\_clock\_out" : 1,    "expense\_data" : \[            {            "expense\_id" : 13,            "date": "01/12/2025",            "expense\_name": "Uang Kebersihan",            "amount": 12000,            "reason" : "",            },            {            "expense\_id" : 14,            "date": "01/12/2025",            "expense\_name": "Uang Kebersihan",            "amount": 12000,            "reason" : "",            }        \] "paging" : {  "total\_record": 84,  "page\_current": 1,  "page\_limit": 10,  "page\_total": 9 }, "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| Case : empty state / tidak terdapat data berdasarkan pencarian {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

2. # **Expense Detail   (Ready to dev)**

Create New endpoint   
URL : 

| {{url}}/mobile/v1/expense/{expense\_id} |
| :---- |

Method		: GET

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Path Variable** 

---

|  expense\_id | 32 |
| :---- | :---- |

### **Example Request** 

---

**Example Request Default**

| bash curl \-X GET "{{url}}/mobile/v1/expense/{expense\_id}" \\  \-H "Authorization: Bearer $TOKEN" \\  \-H "Content-Type: application/json" | jq '.' |
| :---- |

### **Response** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Object | \- |  \- |
| expense\_id  | Int | 8 | Id expense acf.expense.expense\_id |
| date | Date  |  | Tanggal create expense contoh : DD/MM/YYYY acf.expense.date |
| expense\_name | Varchar | 50 | Expense name   diambil dari acf.expense\_type.expense\_name |
| amount | Numeric | 11 | Nominal expense acf.expense.amount |
| reason | Varchar | 255 | reason  acf.expense.note |
| is\_clock\_out | int 1 |  | 1 : clock in 2 : clock outmobile.attendances.type  |
| visits | Array  |  |  |
| outlet\_id | int | 8 | outlet id  |
| outlet\_code | Varchar | 30 | outlet code |
| outlet\_name | Varchar | 150 | nama outlet  |
| outlet\_address1 | Varchar | 150 | alamat outlet |
| file\_name | varchar | 255 | Nama File  |
| file\_type | varchar | 50 | photo disimpan dalam format JPG |
| media\_category | ENUM('image','video') |  | Membedakan jenis konten: *image* atau *video* |
| file\_url | text |  | url file |
| file\_size | BIGINT |  | Untuk mengetahui ukuran file asli |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses  { "message": "", "data" : {  "expense\_id" : 13,  "date": "01/12/2025",  "expense\_name": "Uang Kebersihan",  "amount": 12000,  "reason" : "",  "is\_clock\_out" : 1,  "visits" : \[    {      "outlet\_id" : 12,      "outlet\_code" : "O021",      "outlet\_name" : "Outlet1",      "outlet\_address1" :"Jalan merpati"    },     {      "outlet\_id" : 13,      "outlet\_code" : "O022",      "outlet\_name" : "Outlet2",      "outlet\_address1" :"Jalan merpati"    }  \],   "files": \[  {     "file\_name": "arrival\_IMG\_20250119\_155922",     "file\_type": “JPG”,     "file\_key": "users/123/avatar/profile\_123.png",     "media\_category": "image",    "file\_url": "",     "file\_size": 204800  } \] } "paging" : {  "total\_record": 84,  "page\_current": 1,  "page\_limit": 10,  "page\_total": 9 }, "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

3. # **Expense Lookup  (dependency dg fitur lain → hardcode be)**

Create New endpoint   
URL : 

| {{url}}/mobile/v1/expense\_type?is\_active=1\&mode=lookup |
| :---- |

Method		: GET

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| is\_active | Integer | Yes | default : 1 |
| mode  | Varchar | Yes | lookup |

### **Example Request** 

---

**Example Request Default**

| bash curl \-X GET "{{url}}/mobile/v1/expense?limit10" \\  \-H "Authorization: Bearer $TOKEN" \\  \-H "Content-Type: application/json" | jq '.' |
| :---- |

### **Response  :  diambil dari expense\_type.source \= 2**  ---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- | \- |
| expense\_type\_id | INt | 8 | expense\_type.expense\_type\_id |
| expense\_type\_code | Varchar | 20 | expense\_type.expense\_type\_code |
| expense\_type\_name | Varchar | 50 | expense\_type.expense\_type\_name |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data  |
| :---- |
| Case : empty state / tidak terdapat data berdasarkan pencarian { "message": "", "data" : \[ {  "expense\_type\_id" : 13,  "expense\_type\_code": "E002",  "expense\_type\_name": "Uang Kebersihan", }, {  "expense\_type\_id" : 12,  "expense\_type\_code": "E002",  "expense\_type\_name": "Uang Keamanan", }, \] "paging" : {  "total\_record": 84,  "page\_current": 1,  "page\_limit": 10,  "page\_total": 9 }, "request\_id": "6915a5e8e3f53f84fe73517f" }   |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

4. # **Outlet Lookup  → (per tgl 1 Dec outlet di takeout)**

Create New endpoint   
URL : 

| {{url}}/mobile/v1/outlet?is\_active=1\&mode=lookup\&date=2025-12-31 |
| :---- |

Method		: GET

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| is\_active | Integer | Yes | default : 1 |
| mode  | Varchar | Yes | lookup |

### **Example Request** 

---

**Example Request Default**

| bash curl \-X GET "{{url}}/mobile/v1/expense?limit10" \\  \-H "Authorization: Bearer $TOKEN" \\  \-H "Content-Type: application/json" | jq '.' |
| :---- |

### **Response  :** 

| SELECT      o.outlet\_id,      o.outlet\_code,      o.outlet\_name FROM pjp.outlet\_visit\_list v JOIN mst.m\_outlet o      ON o.outlet\_id \= v.outlet\_id JOIN pjp.permanent\_journey\_plans p      ON p.id \= v.pjp\_id WHERE      v.is\_planned \= TRUE     AND v.date \= DATE '2025-11-26';  ![][image1] |
| :---- |

*pjp.permanent\_journey\_plans*

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| outlet\_id | int | 8 | m\_outlet.outlet\_id |
| outlet\_code | varchar | 30 | m\_outlet.outlet\_code |
| outlet\_name | varchar | 150 | m\_outlet.outlet\_name |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data  |
| :---- |
| Case : empty state / tidak terdapat data berdasarkan pencarian { "message": "", "data" : \[ { "outlet\_id" : 13, "outlet\_code": "E002", "outlet\_name": "Uang Kebersihan", }, { "outlet\_id" : 15, "outlet\_code": "E002", "outlet\_name": "Uang Kebersihan", }, \] "paging" : { "total\_record": 84, "page\_current": 1, "page\_limit": 10, "page\_total": 9 }, "request\_id": "6915a5e8e3f53f84fe73517f" }   |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

5. #  **Create Expense  (Ready to dev)**

| {{url}}/mobile/v1/expense |
| :---- |

### **Method : POST**

### **Headers**

---

| Accept | multipart/form-data |
| :---- | :---- |
| Authorization  | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

Upload file menggunakan Object Storage Huawei 

### **Body** 

| Field | Type | Required | Description |
| :---- | :---- | :---- | :---- |
| expense\_type\_id | integer | yes |  |
| outlet\_id | array(integer) | yes | per tgl 1 Des 2025 outlet di takeout |
| amount | integer | yes |  |
| note | varchar | No |  |
| file | array(file)  | yes | file yang akan di upload |
| folder | string | yes | folder tujuan bucket |

create data acf.expense 

#### 

| Nama Kolom | Value  |
| :---- | :---- |
| cust\_id\*\* | diisi cust\_id yang melakukan create expense |
| expense\_id\* | generate BE  |
| doc\_no | generate BE EYYYYMMDD-3digitRunnningNumber  E20261201222 |
| expense\_type\_id\*\* | req body expense\_type\_id |
| source | 2 |
| date | diisi tanggal create  |
| amount | dr req body amount |
| note | dr req body note |
| collector\_id\*\* | dari user\_id yang login   |
| created\_by\*\* | dari user\_id yang login   |
| created\_at | diisi waktu saat create  |
| updated\_by\*\* | NULL |
| updated\_at | NULL |
| deleted\_by\*\* | NULL |
| deleted\_at | NULL |
| is\_del | NULL |

## 

list error : 

- Data outlet not found  
- Expense type not found  
- File must not exceed 3 

6. #  **Delete Expense  (Ready to dev)**

| {{url}}/mobile/v1/expense/{expense\_id} |
| :---- |

### **Method : DELETE**

### **Headers**

---

| Accept | multipart/form-data |
| :---- | :---- |
| Authorization  | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

**Path Variable** 

---

|  expense\_id | 32 |
| :---- | :---- |

7. #  **Update Expense  (Ready to dev)**

| {{url}}/mobile/v1/expense/{expense\_id} |
| :---- |

### **Method : PATCH**

### **Headers**

---

| Accept | multipart/form-data |
| :---- | :---- |
| Authorization  | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

**Path Variable** 

---

|  expense\_id | 32 |
| :---- | :---- |

### **Body** 

| Field | Type | Required | Description |
| :---- | :---- | :---- | :---- |
| expense\_type\_id | integer | yes |  |
| outlet\_id | array(integer) | yes |  |
| amount | integer | yes |  |
| note | varchar | yes |  |
| files | array(file)   | yes | file yang akan di upload |
| delete\_file\_ids | multiple integer with separate coma | yes | folder tujuan bucket |

[image1]: <data:image/png;base64,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>