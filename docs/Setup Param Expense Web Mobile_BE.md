1. ## **Update Database**

   [Link](https://docs.google.com/document/d/1-xZGPB3ZmwVIewWcCm4BpuLfv5jcsQLfNI_txVO2Z6E/edit?tab=t.0#heading=h.mv1xf7683h6o)  
   **Nama Database:** expense\_type   
   **Schema** : acf  
   **Field yang ditambah :**   
* is\_active (Bool) → True (active) False (inactive)  
* source   
  1: web

  2 : mobile

2. ## **Expense List** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/acf/v1/expense?sort=created\_date:desc &\&limit=5\&page=1

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| q | Integer | No | query param by expense\_type\_code dan expense\_type\_name |
| source | Integer | Yes | 1:  web2:  mobile |
| page | Integer | Yes | default : 1 |
| limit  | Integer | Yes | default : 5  |
| sort | String | Yes | default created\_date:desc  |

### **Example Request** 

---

**Example Request Default**

| curl \--location \-g \\ '{{url}}/acf/v1/expense?source=1\&page=1\&limit=5\&sort=created\_date:desc' \\ \--header 'Accept: application/json' |
| :---- |

### **Response** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| expense\_type\_id | Integer | 8 | **acf.expense\_type.**expense\_type\_id |
| expense\_type\_code | Varchar | 20 | **acf.expense\_type.**expense\_type\_code |
| expense\_type\_name | Varchar | 50 | **acf.expense\_type.**expense\_type\_name |
| status | Integer | 3 | **acf.expense\_type.**is\_active True : ActiveFalse :  Inactive |
| update\_by | Integer | 8 | **sys.m\_user.**user\_name  **acf.expense\_type.**updated\_by **relasi** dengan tabel **sys.m\_user.user\_id** |
| update\_date | timestamptz | 6 | **acf.expense\_type.**updated\_at |

**Example Response :** 

| Case : sukses dan terdapat data {  "message": "",  "data": \[    {      "expense\_type\_id": 12,      "expense\_type\_code": "A001",      "expense\_type\_name": "Biaya Transportasi",      "status": True,      "update\_by": 1001,      "update\_date": "2025-12-01T10:00:00Z"    },    {      "expense\_type\_id": 13,      "expense\_type\_code": "A002",      "expense\_type\_name": "Biaya Konsumsi",      "status": False,      "update\_by": 1002,      "update\_date": "2025-12-01T11:00:00Z"    }  \],  "paging": {    "total\_record": 84,    "page\_current": 1,    "page\_limit": 10,    "page\_total": 9  },  "request\_id": "6915a5e8e3f53f84fe73517f" }   |
| :---- |
| Case : empty state / tidak terdapat data berdasarkan pencarian {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

3. # **Create Expense**  

 Create new endpoint :  
Content-Type 	: application/json  
Method		: POST  
URL		: {{url}}/acf/v1/expense

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Body** 

### ---

| Field | Type | Required | Description |
| :---- | :---- | :---- | :---- |
| expense\_type\_code | varchar(20) | Yes |  |
| expense\_type\_name | varchar(50) | Yes |  |
| source | int(3) | Yes | 1 : web → jika create lewat menu **setup parameter web** 2: mobile → jika create lewat menu **setup parameter mobile**  |

### **Response JSON :** 

### ---

1) Case Sukses Menyimpan Data : 

{  
   "message": "Data saved successfully",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

2) Case Sukses Menyimpan Data :   
   (expense\_type\_code dan expense\_type\_name tidak boleh sama dengan data sebelumnya) 

{  
   "message": "Data already exists",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

3) Case Gagal Menyimpan Data (Selain duplicate data) : 

{  
   "message": "Failed to save data, please try again",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

4. # **Update Expense**  

 Create new endpoint :  
Content-Type 	: application/json  
Method		: PATCH  
URL		: {{url}}/acf/v1/expense/:expense\_type\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Body** 

### ---

| Field | Type | Required | Description |
| :---- | :---- | :---- | :---- |
| expense\_type\_code | varchar(20) | Yes |  |
| expense\_type\_name | varchar(50) | Yes |  |

### **Response JSON :** 

### ---

4) Case Sukses Menyimpan Data : 

{  
   "message": "Data update successfully",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

5) Case Sukses Menyimpan Data :   
   (expense\_type\_code dan expense\_type\_name tidak boleh sama dengan data sebelumnya) 

{  
   "message": "Data already exists",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

6) Case Gagal Menyimpan Data (Selain duplicate data) : 

{  
   "message": "Failed to update data, please try again",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

5. # **Delete Expense List**  

 Create new endpoint :  
Content-Type 	: application/json  
Method		: Delete  
URL		: {{url}}/acf/v1/expense/:expense\_type\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4  |

**Path Variable** 

---

|  expense\_type\_id | 32 |
| :---- | :---- |

### **Response JSON :** 

### ---

* Case Sukses Menghapus Data : 

   "message": "Data delete successfully",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

* Case Gagal Menyimpan Data (Selain duplicate data) : 

{  
   "message": "Failed to delete data, please try again",  
   "request\_id": "6915a6dd2395083c685e8e16"  
}

