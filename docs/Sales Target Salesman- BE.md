1. # **Sales Target \- List** 

   

 Create new endpoint :  
Type 		: application/json  
Method		: GET  
URL		: {{url}}/master/v1/sales-target?page=1\&limit=10\&sort=asc\&year=2025

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

### **Example Request** 

---

**Example Request Default**

| curl \-X GET "{{url}}/master/v1/sales-target?page=1\&limit=10\&sort=asc\&year=2025" \\   \-H "Authorization: Bearer $TOKEN" \\   \-H "Accept: application/json" |
| :---- |

### **Response  :** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| sales\_target\_id | Int | 8 | **mst.m\_sales\_target.**sales\_target\_id |
| year | int  | 4 | **mst.m\_sales\_target.**year |
| month | Int  | 2 | **mst.m\_sales\_target.**month |
| allocated\_total | int | 11 | **mst.m\_sales\_target.**allocated\_total |
| monthly\_target | Int  | 11 | **mst.m\_sales\_target.**monthly\_target |
| remaining | int | 11 | **mst.m\_sales\_target.**remaining |
| status | int | 11 | **mst.m\_sales\_target.**status INPUT: current\_year, current\_month, year, month, status OUTPUT: result\_status IF current\_year \= year AND current\_month \= month THEN     IF status \= 1 THEN         result\_status \= "Active"     ELSE IF status \= 0 THEN         result\_status \= "Draft"     ELSE IF status \= 2 THEN         result\_status \= "Nonactive"     ELSE         result\_status \= "Unknown"     ENDIF ELSE     result\_status \= "Nonactive" ENDIF RETURN result\_status  |
| updated\_by | Int | 8 | sys.m\_user.user\_name jika data belum di edit : (updated\_at \= NULL) mst.m\_sales\_target.created\_by jika data sudah di edit : (updated\_at \= NOT NULL) mst.m\_sales\_target.updated\_by |
| updated\_at | Timestampz | 6 | jika data belum di edit : (updated\_at \= NULL) mst.m\_sales\_target.created\_at jika data sudah di edit : (updated\_at \= NOT NULL) mst.m\_sales\_target.updated\_at |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data {  "message": "",  "data": \[    {      "sales\_target\_id": 1,      "month": 2,      "year": 1,      "allocated\_total": 500000,      "monthly\_target": 300000,      "remaining": 200000,      "updated\_by": "Charles",      "updated\_at": "2025-01-01"    }  \],  "paging": {    "total\_record": 84,    "page\_current": 1,    "page\_limit": 10,    "page\_total": 9  },  "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| Case : empty state / tidak terdapat data berdasarkan pencarian {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

2. # **Sales Target \- Detail**  

 Create new endpoint :  
Type 		: application/json  
Method		: GET  
URL		: {{url}}/master/v1/sales-target/:sales\_target\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Path Variables:** 

---

|  sales\_target\_id | 32 |
| :---- | :---- |

### **Example Request** 

---

**Example Request Default**

| curl \-X GET "{{url}}/master/v1/sales-target/121" \\   \-H "Authorization: Bearer $TOKEN" \\   \-H "Accept: application/json" |
| :---- |

### **Response  :** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Object | \- |  \- |
| sales\_target\_id | Int | 8 | **mst.m\_sales\_target.**sales\_target\_id |
| year | int  | 4 | **mst.m\_sales\_target.**year |
| month | Int  | 2 | **mst.m\_sales\_target.**month |
| allocated\_total | int | 11 | **mst.m\_sales\_target.**allocated\_total |
| monthly\_target | Int  | 11 | **mst.m\_sales\_target.**monthly\_target |
| remaining | int | 11 | **mst.m\_sales\_target.**remaining |
| updated\_by | Int | 8 | sys.m\_user.user\_name jika data belum di edit : (updated\_at \= NULL) mst.m\_sales\_target.created\_by jika data sudah di edit : (updated\_at \= NOT NULL) mst.m\_sales\_target.updated\_by |
| updated\_at | Timestampz | 6 | jika data belum di edit : (updated\_at \= NULL) mst.m\_sales\_target.created\_at jika data sudah di edit : (updated\_at \= NOT NULL) mst.m\_sales\_target.updated\_at |
| details | Array | \- |  \- |
| sales\_allocated\_id | Integer | 8 | **mst.m\_sales\_allocated.**sales\_allocated\_id |
| sales\_target\_id | Integer | 8 | **mst.m\_sales\_target.**sales\_target\_id |
| salesman\_id | Integer | 8 | **mst.m\_sales\_allocated.**salesman\_id |
| sales\_name | Varchar | 150 | **mst.salesman.**sales\_name |
| opr\_type | Varchar | 150 | **mst.salesman.**opr\_type  jika ada 2 data , resp \= Canvas, Taking Order |
| distributor\_id | Integer | 8 | **mst.m\_sales\_allocated.**distributor\_id |
| distributor\_code | Varchar | 20 | **mst.m\_distributor.**distributor\_code |
| distributor\_name | Varchar | 150 | **mst.m\_distributor.**distributor\_name |
| sales\_team\_id |  |  | **mst.salesman.**sales\_team\_id  |
| sales\_team\_code |  |  | **mst.m\_sales\_team.**sales\_team\_code |
| sales\_team\_name |  |  | **mst.m\_sales\_team.**sales\_team\_name |
| allocated | Integer | 11 | **mst.m\_sales\_allocated.**allocated |
| is\_active | bool |  | **mst.m\_sales\_allocated.**is\_active |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data  { "message": "", "data":  {      "sales\_target\_id": 1,      "month": 2,      "year": 1,      "allocated\_total": 500000,      "monthly\_target": 300000,      "remaining": 200000,      "updated\_by": "Charles",      "updated\_at": "2025-01-01",      "details" : \[            {                "sales\_allocated\_id": 1,                "sales\_target\_id": 2,                "salesman\_id": 12,                "sales\_name": "Alexa",                "opr\_type": "Canvas, Taking Order",                "distributor\_id": 1,                "distributor\_code" :"D001",                "distributor\_name" :"John",                "channel\_id": "Charles",                "channel\_code" :"C01",                "channel\_name" : "MT",                "allocated": "2025-01-01",                "is\_active": true            },            {                "sales\_allocated\_id": 2,                "sales\_target\_id": 2,                "salesman\_id": 13,                "sales\_name": "Charles",                "distributor\_id": 1,                "distributor\_code" :"D001",                "distributor\_name" :"John",                "channel\_id": "Charles",                "channel\_code" :"C01",                "channel\_name" : "MT",                "allocated": "2025-01-01",                "is\_active": true            } \] }, "paging": {   "total\_record": 84,   "page\_current": 1,   "page\_limit": 10,   "page\_total": 9 }, "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

3. # **Sales Target \- Add**  

 Create new endpoint :  
Type 		: application/json  
Method		: POST  
URL		: {{url}}/master/v1/sales-target

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Body Variables:** 

---

| Field | Type | Required | Description |
| :---- | :---- | :---- | :---- |
| sales\_target\_distributor\_yearly\_id\*\* | int(8) | Yes | diisi dari dropdown year, ambil resp ***sales\_target\_distribut***or\_yearly\_id |
| sales\_target\_distributor\_monthly\_id\*\* | int(8) | Yes | diisi dari dropdown month, ambil resp ***details.sales\_target\_distributor\_monthly\_id*** |
| month | Int(2) | Yes | Diisi inputan month  |
| year  | int(4) | Yes | Diisi inputan year  |
| allocated\_total | int(11) | Yes | Diisi Total Allocated   |
| monthly\_target  | int(11) | Yes | Diisi Monthly Target  |
| remaining | int(11) | Yes | Diisi Remaining |
| status  | int(11) | Yes | dikirim 1  (Aktif)  |
| data  | Array  | Yes |  |
| salesman\_id | int(8) | Yes | diisi dari api [https://best.scyllax.online/master/v1/salesman?page=1\&limit=999](https://best.scyllax.online/master/v1/salesman?page=1&limit=999) ambil emp\_id |
| sales\_team\_id  | int(8) | Yes | diisi dari api [https://best.scyllax.online/master/v1/salesman?page=1\&limit=999](https://best.scyllax.online/master/v1/salesman?page=1&limit=999) ambil sales\_team\_id |
| allocated | int(11) | Yes | Dari inputan allocated |

### **Example Request** 

---

**Example Request Default**

| curl \-X POST "{{url}}/master/v1/sales-target" \\  \-H "Authorization: Bearer $TOKEN" \\  \-H "Accept: application/json" \\  \-H "Content-Type: application/json" \\  \-d '{    "month": 2,    "year": 2025,    "allocated\_total": 10000000,    "monthly\_target": 15000000,    "remaining": 5000000,    "status": 1,    "data": \[      {        "salesman\_id": 12,        "sales\_team\_id": 333,        "allocated": 5000000      },      {        "salesman\_id": 44,        "sales\_team\_id": 55,        "allocated": 5000000      }    \]  }' |
| :---- |

1) Case Sukses Menyimpan Data : 

{  
   "message": "Data saved successfully",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

2) Case Gagal Menyimpan Data :

{  
   "message": "Failed to save data, please try again",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

impact database : 

* insert ke tabel mst.m\_sales\_target

| Field | Data Source |
| :---- | :---- |
| cust\_id\*\* | diisi cust\_id user login |
| sales\_target\_id\* | generate BE |
| sales\_target\_distributor\_yearly\_id\*\* | req body : sales\_target\_distributor\_yearly\_id |
| sales\_target\_distributor\_monthly\_id\*\* | req body : sales\_target\_distributor\_monthly\_id |
| month | req body : month |
| year | req body : year |
| allocated\_total | req body : allocated\_total |
| monthly\_target | req body : monthly\_target |
| remaining | req body : remaining |
| status | 1   |
| created\_by | diisi user yang melakukan create |
| created\_at | diisi time saat execute  |
| updated\_by | NULL |
| updated\_at | NULL |
| deleted\_by | NULL |
| deleted\_at | NULL |
| is\_del | NULL |

* insert ke tabel mst.m\_sales\_allocated

| Field | Data Source |
| :---- | :---- |
| cust\_id\*\* | diisi cust\_id user login |
| sales\_allocated\_id\* | generate BE |
| sales\_target\_id\*\* | diisi sales\_target\_id  |
| salesman\_id\*\* | req body : salesman\_id\*\* |
| sales\_team\_id\*\* | req body : sales\_team\_id\*\* |
| allocated | req body : allocated |
| is\_active | TRUE |
| created\_by | diisi user yang melakukan create |
| created\_at | diisi time saat execute  |
| updated\_by | NULL |
| updated\_at | NULL |
| deleted\_by | NULL |
| deleted\_at | NULL |
| is\_del | NULL |

4. # **Sales Target \- Edit**  

 Create new endpoint :  
Type 		: application/json  
Method		: PATCH  
URL		: {{url}}/master/v1/sales-target/:sales\_target\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Path Variables:** 

---

|  sales\_target\_id | 32 |
| :---- | :---- |

### **Body Variables:** 

---

| Field | Type | Required | Description |
| :---- | :---- | :---- | :---- |
| sales\_target\_distributor\_yearly\_id\*\* | int(8) | No | diisi dari dropdown year, ambil resp ***sales\_target\_distribut***or\_yearly\_id |
| sales\_target\_distributor\_monthly\_id\*\* | int(8) | No | diisi dari dropdown month, ambil resp ***details.sales\_target\_distributor\_monthly\_id*** |
| month | Int(2) | No | Diisi inputan month  |
| year  | int(4) | No | Diisi inputan year  |
| allocated\_total | int(11) | No | Diisi Total Allocated   |
| monthly\_target  | int(11) | No | Diisi Monthly Target  |
| remaining | int(11) | No | Diisi Remaining |
| status  | int(11) | No | 0 : Draft  1 : Aktif  2 : Nonaktif  |
| data  | Array  | No |  |
| salesman\_id | int(8) | No | diisi dari api [https://best.scyllax.online/master/v1/salesman?page=1\&limit=999](https://best.scyllax.online/master/v1/salesman?page=1&limit=999) ambil emp\_id |
| sales\_team\_id  | int(8) | No | diisi dari api [https://best.scyllax.online/master/v1/salesman?page=1\&limit=999](https://best.scyllax.online/master/v1/salesman?page=1&limit=999) ambil sales\_team\_id |
| allocated | int(11) | No | Dari inputan allocated |

---

**Example Request** 

| curl \-X PATCH "{{url}}/master/v1/sales-target/121" \\  \-H "Authorization: Bearer $TOKEN" \\  \-H "Accept: application/json" \\  \-H "Content-Type: application/json" \\  \-d '{    "month": 2,    "year": 2025,    "allocated\_total": 10000000,    "monthly\_target": 15000000,    "remaining": 5000000,    "status": 1,    "data": \[      {        "salesman\_id": 12,        "sales\_team\_id": 333,        "allocated": 5000000      },      {        "salesman\_id": 44,        "sales\_team\_id": 55,        "allocated": 5000000      }    \]  }'  |
| :---- |

3) Case Sukses Menyimpan Data : 

{  
   "message": "Data saved successfully",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

4) Case Gagal Menyimpan Data (Selain duplicate data) : 

{  
   "message": "Failed to save data, please try again",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

