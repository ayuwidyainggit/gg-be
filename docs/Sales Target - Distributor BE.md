1. # **Sales Target Distributor \- List** 

**Flow bisnis:**

* User membuat **sales target per bulan** untuk **tahun 2025 (Jan–Des)**  
* Target hanya **aktif pada bulan berjalan**.  
* Ketika bulan berganti:  
  * Target bulan yang sudah lewat **otomatis menjadi non-aktif (deactive)**.  
* Proses harus **otomatis**, **konsisten**, dan **tidak bergantung user action**


 Create new endpoint :  
Type 		: application/json  
Method		: GET  
URL		: {{url}}/master/v1/sales-target-distributor?page=1\&limit=10\&sort=asc\&year=2025

### 

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| page | Integer | Yes | default : 1 |
| limit  | Integer | Yes | default : 20 |
| sort | String | Yes | default created\_date:desc  |
| year | Integer | No |  |
| status   | Array\<integer\> | No | jika null ditampilkan semua status 0 : Draft  1 : Aktif  2 : Nonaktif  Jika year \> tahun sekarang \= **Inactive**Jika year \<=bulan dan tahun sekarang, lanjut cek pada is\_active,      Jika is\_active \= true → **Active**     Jika is\_active \= false → **Inactive** |

### **Example Request** 

---

**Example Request Default**

| curl \-X GET "{{url}}/master/v1/sales-target-distributor?page=1\&limit=10\&sort=asc\&year=2025" \\   \-H "Authorization: Bearer $TOKEN" \\   \-H "Accept: application/json" |
| :---- |

### **Response  :** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| sales\_target\_distributor\_yearly\_id | Int | 8 | **mst.m\_sales\_target\_distributor\_yearly.**sales\_target\_distributor\_yearly\_id |
| distributor\_id | int | 4 | **mst.m\_sales\_target\_distributor\_yearly.**distributor\_id |
| distributor\_code | varchar | 20 | **mst.m\_distributor.**distributor\_code |
| distributor\_name | varchar | 150 | **mst.m\_distributor.**distributor\_name |
| year | int  | 4 | **mst.m\_sales\_target\_distributor\_yearly.**year |
| yearly\_target | Int  | 11 | **mst.m\_sales\_target\_distributor\_yearly.**yearly\_target |
| updated\_by | Int | 8 | sys.m\_user.user\_name jika data belum di edit : (updated\_at \= NULL) **mst.m\_sales\_target.**created\_by jika data sudah di edit : (updated\_at \= NOT NULL) **mst.m\_sales\_target.**updated\_by |
| updated\_at | Timestampz | 6 | jika data belum di edit : (updated\_at \= NULL) **mst.m\_sales\_target.**created\_at jika data sudah di edit : (updated\_at \= NOT NULL) **mst.m\_sales\_target.**updated\_at |
| status | varchar | 50 | **mst.m\_sales\_target\_distributor\_yearly.**status |
| user\_inactive | int | 8 | **mst.m\_sales\_target\_distributor\_yearly.**user\_inactive |
| inactive\_at | timestamptz | 6 | **mst.m\_sales\_target\_distributor\_yearly.**inactive\_at |
| is\_active | bool |  | mst.m\_sales\_target\_distributor\_yearly.is\_active **delete response** |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data {   "message": "",   "data": \[     {       "sales\_target\_distributor\_yearly\_id": 1,       "distributor\_id": 10,       "distributor\_code": "DST-001",       "distributor\_name": "Distributor Utama",       "year": 2025,       "yearly\_target": 300000,       "updated\_by": 1001,       "updated\_at": "2025-01-01T00:00:00Z",      "status": "Active",       "is\_active": true     }   \],   "paging": {     "total\_record": 84,     "page\_current": 1,     "page\_limit": 10,     "page\_total": 9   },   "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| Case : empty state / tidak terdapat data berdasarkan pencarian {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

2. # **Sales Target Distributor \- Detail**  

 Create new endpoint :  
Type 		: application/json  
Method		: GET  
URL		: {{url}}/master/v1/sales-target-distributor/:sales\_target\_distributor\_yearly\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Path Variables:** 

---

|  sales\_target\_distributor\_yearly\_id | 32 |
| :---- | :---- |

### **Example Request** 

---

**Example Request Default**

| curl \-X GET "{{url}}/master/v1/sales-target-distributor/121" \\   \-H "Authorization: Bearer $TOKEN" \\   \-H "Accept: application/json" |
| :---- |

#### **query allocation :**  {#query-allocation-:}

| SELECT      y.year,     y.cust\_id,     m.sales\_target\_distributor\_monthly\_id,     m.month,     m.monthly\_target,     t.allocated\_total,     a.salesman\_id,     a.allocated FROM mst.m\_sales\_target\_distributor\_yearly y 	JOIN mst.m\_sales\_target\_distributor\_monthly m ON y.sales\_target\_distributor\_yearly\_id \= m.sales\_target\_distributor\_yearly\_id 	left JOIN mst.m\_sales\_target t ON m.sales\_target\_distributor\_monthly\_id \= t.sales\_target\_distributor\_monthly\_id 	left JOIN mst.m\_sales\_allocated a ON t.sales\_target\_id \= a.sales\_target\_id WHERE  	y."year" \= 2026 and 	y.cust\_id \= 'C22001'  Notes :  jika allocated\_total atau salesman\_id \= null , artinya sales target belum di alokasikan  is\_allocated \= FALSE allocation\_total \= allocated\_total  jika allocated\_total atau salesman\_id \= not null , artinya sales target sudah di alokasikan  is\_allocated \= TRUE allocation\_total \= allocated\_total   ![][image1]  |
| :---- |

### **Response  :** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Object | \- |  \- |
| sales\_target\_distributor\_yearly\_id | Int | 8 | **mst.m\_sales\_target\_distributor\_yearly.**sales\_target\_distributor\_yearly\_id |
| year | int  | 4 | **mst.m\_sales\_target\_distributor\_yearly.**year |
| yearly\_target | Int  | 11 | **mst.m\_sales\_target\_distributor\_yearly.**yearly\_target |
| updated\_by | Int | 8 | **sys.m\_user**.user\_name jika data belum di edit : (updated\_at \= NULL) **mst.m\_sales\_target.**created\_by jika data sudah di edit : (updated\_at \= NOT NULL) **mst.m\_sales\_target.**updated\_by |
| updated\_at | Timestampz | 6 | jika data belum di edit : (updated\_at \= NULL) **mst.m\_sales\_target.**created\_at jika data sudah di edit : (updated\_at \= NOT NULL) **mst.m\_sales\_target.**updated\_at |
| area\_id | int  | 8 | **mst.m\_area.**area\_id |
| area\_code | varchar | 10 | **mst.m\_area.**area\_code |
| area\_name | varchar | 150 | **mst.m\_area.**area\_name |
| region\_id | int  | 8 | **mst.m\_sales\_target\_distributor\_yearly.**region\_id |
| region\_code | varchar | 10 | **mst.m\_region.**region\_code |
| region\_name | varchar | 150 | **mst.m\_region.**region\_name |
| distributor\_id | int  | 8 | **mst.m\_sales\_target\_distributor\_yearly.**distributor\_id |
| distributor\_code | varchar | 20 | **m\_distributor.**distributor\_code |
| distributor\_name | varchar | 150 | **m\_distributor.**distributor\_name |
| status | varchar | 50 | **mst.m\_sales\_target\_distributor\_yearly.**status |
| is\_allocated | bool |  | **cek penjelasan [disini](#query-allocation-:)**   |
| allocation\_total | numeric | 11 |  |
| **details** | **Array** | **\-** |  \- |
| sales\_target\_distributor\_monthly\_id | int  | 8 | **mst.m\_sales\_target\_distributor\_monthly.**sales\_target\_distributor\_monthly\_id |
| month  | int | 2 | **mst.m\_sales\_target\_distributor\_monthly.**month  |
| monthly\_target | int | 11 | **mst.m\_sales\_target\_distributor\_monthly.**monthly\_target |
| is\_active | bool |  | **mst.m\_sales\_target\_distributor\_monthly.**is\_active |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data  {   "message": "",   "data": {     "sales\_target\_distributor\_yearly\_id": 1,     "year": 2025,     "yearly\_target": 500000,     "updated\_by": 1001,     "updated\_at": "2025-01-01T00:00:00Z",     "area\_id": 10,     "area\_code": "AR-01",     "area\_name": "Area Jakarta",     "region\_id": 5,     "region\_code": "RG-01",     "region\_name": "Region Barat",     "distributor\_id": 1,     "distributor\_code": "D001",     "distributor\_name": "John",     "status": "Active" ,     "is\_allocated": TRUE ,     "allocated\_total": 12000000 ,     "details": \[       {         "sales\_target\_distributor\_monthly\_id": 1,         "month": 2,         "monthly\_target": 300000,         "is\_active": true       }     \]   },   "paging": {     "total\_record": 84,     "page\_current": 1,     "page\_limit": 10,     "page\_total": 9   },   "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

3. # **Sales Target Distributor \- Add**  

 Create new endpoint :  
Type 		: application/json  
Method		: POST  
URL		: {{url}}/master/v1/sales-target-distributor

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Body Variables:** 

---

| Field | Type | Required | Description |
| :---- | :---- | :---- | :---- |
| year  | int(4) | Yes | Diisi inputan year  |
| area\_id | int(8) | Yes | Diisi dari dropdown area |
| region\_id  | int(8) | Yes | Diisi dari dropdown region |
| distributor\_id | int(8) | Yes | Diisi dari dropdown distributor |
| yearly\_target | int(11)  | Yes | Diisi dari hasil summary monthly target |
| status | int(8) | Yes | 0 : Draft  1 : Aktif  2 : Nonaktif  (dikirim 1 jika sudah di submit, 0 jika masih draft )  |
| data  | Array  | Yes |  |
| month | int(2) | Yes | Diisi bulan  1: Jan 2: Feb  …. dst |
| monthly\_target | int(11) | Yes | diisi amount sesuai bulan yang diinput |

### **Example Request** 

---

**Example Request Default**

| curl \-X POST "{{url}}/master/v1/sales-target-distributor" \\  \-H "Authorization: Bearer $TOKEN" \\  \-H "Accept: application/json" \\  \-H "Content-Type: application/json" \\  \-d ' {  "year": 2025,  "area\_id": 12,  "region\_id": 3,  "distributor\_id": 45,  "yearly\_target": 120000000,  "status": 2,  "data": \[    {      "month": 1,      "monthly\_target": 10000000    },    {      "month": 2,      "monthly\_target": 9500000    },    {      "month": 3,      "monthly\_target": 10500000    }  \] }' |
| :---- |

1) Case Sukses Menyimpan Data : 

{  
   "message": "Data saved successfully",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

2) Case Gagal Menyimpan Data (Selain duplicate data) : 

{  
   "message": "Failed to save data, please try again",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

# **4\. Sales Target Distributor \- Edit**  

 Create new endpoint :  
Type 		: application/json  
Method		: PATCH  
URL		: {{url}}/master/v1/sales-target-distributor/:sales\_target\_distributor\_yearly\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Path Variables:** 

---

|  sales\_target\_distributor\_yearly\_id | 32 |
| :---- | :---- |

### **Body Variables: (kirim data yang diubah saja )** 

---

| Field | Type | Required | Description |
| :---- | :---- | :---- | :---- |
| year  | int(4) | No | Diisi inputan year  |
| area\_id | int(8) | No | Diisi dari dropdown area |
| region\_id  | int(8) | No | Diisi dari dropdown region |
| distributor\_id | int(8) | No | Diisi dari dropdown distributor |
| yearly\_target | int(11)  | No | Diisi dari hasil summary monthly target |
| status | int(8) | No | 0 : Draft  1 : Aktif  2 : Nonaktif  |
| is\_active | bool | No | True : Active False : Inactive |
| data  | Array  | No |  |
| month | int(2) | No | Diisi bulan  1: Jan 2: Feb  …. dst |
| monthly\_target | int(11) | No | diisi amount sesuai bulan yang diinput |

tambahkan logic jika FE kirim status \= 2 → BE update ke user\_inactive dan inactive\_at  
---

**Example Request** 

| curl \-X PATCH "{{url}}/master/v1/sales-target/121" \\  \-H "Authorization: Bearer $TOKEN" \\  \-H "Accept: application/json" \\  \-H "Content-Type: application/json" \\  \-d '{  "year": 2025,  "area\_id": 12,  "region\_id": 3,  "distributor\_id": 45,  "yearly\_target": 120000000,  "status": 2,  "data": \[    {      "month": 1,      "monthly\_target": 10000000    },    {      "month": 2,      "monthly\_target": 9500000    },    {      "month": 3,      "monthly\_target": 10500000    }  \] }' |
| :---- |

1) Case Sukses Menyimpan Data : 

{  
   "message": "Data saved successfully",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

2) Case Gagal Menyimpan Data (Selain duplicate data) : 

{  
   "message": "Failed to save data, please try again",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}
