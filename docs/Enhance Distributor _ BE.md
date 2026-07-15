1. # **Distributor List**

No Enhance endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: [https://best.scyllax.online/master/v1/distributors?page=1\&limit=10](https://best.scyllax.online/master/v1/distributors?page=1&limit=10) 

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| q |  |  |  |
| page | Integer | Yes | default : 1 |
| limit  | Integer | Yes | default : 5  |

### **Example Request** 

| curl \--location \-g \\ '{{url}}/master/v1/outlet\_list?page=1\&limit=5\&sort=created\_date:desc\&status=1' \\ \--header 'Accept: application/json' |
| :---- |

### **Response Example Response :** 

| {     "message": "",     "data": \[         {             "cust\_id": "C22001",             "distributor\_id": 68,             "distributor\_code": "232323",             "distributor\_name": "Pratesis Bangka",             "barcode": null,             "region\_id": 67,             "area\_id": 82,             "channel\_id": 44,             "sub\_distributor\_group\_id": 36,             "sub\_distributor\_group\_code": "",             "sub\_distributor\_group\_name": "",             "dist\_price\_grp\_id": 38,             "address": "jl kedondong",             "province\_id": "36",             "regency\_id": "3737",             "sub\_district\_id": "363636",             "ward\_id": "3366363",             "zip\_code": "62134",             "ot\_loc\_id": 0,             "latitude": "3434343",             "longitude": "34343",             "phone": "",             "fax\_number": "",             "contact\_name": "",             "job\_title": "",             "phone\_no": "",             "wa\_no": "",             "email": "",             "dist\_price\_grp\_code": "003",             "dist\_price\_grp\_name": "JAWA",             "region\_code": "2000",             "region\_name": "Central",             "channel\_code": "000",             "channel\_name": "Non Chanel",             "area\_code": "08",             "area\_name": "Central Jawa",             "province\_code": "36",             "province\_name": "JAWA BARAT",             "regency\_code": "3737",             "regency\_name": "Bogor",             "sub\_district\_code": "363636",             "sub\_district\_name": "Cileungsi",             "ward\_code": "3366363",             "ward\_name": "Cileungsi",             "is\_active": true,             "is\_del": false,             "createdby": null,             "created\_at": null,             "updatedby": null,             "updated\_at": null,             "updated\_by\_name": "Admin Principal 1",             "deleted\_by": null,             "deleted\_at": null,             "customer\_id": ""         },         {             "cust\_id": "C22001",             "distributor\_id": 67,             "distributor\_code": "3434",             "distributor\_name": "Distributor iDetama",             "barcode": null,             "region\_id": 68,             "area\_id": 82,             "channel\_id": 44,             "sub\_distributor\_group\_id": 36,             "sub\_distributor\_group\_code": "",             "sub\_distributor\_group\_name": "",             "dist\_price\_grp\_id": 39,             "address": "jl kedondong",             "province\_id": "36",             "regency\_id": "3737",             "sub\_district\_id": "363636",             "ward\_id": "123456",             "zip\_code": "62125",             "ot\_loc\_id": 0,             "latitude": "3434343",             "longitude": "34343",             "phone": "",             "fax\_number": "",             "contact\_name": "",             "job\_title": "",             "phone\_no": "",             "wa\_no": "",             "email": "",             "dist\_price\_grp\_code": "000",             "dist\_price\_grp\_name": "NON Group Discount",             "region\_code": "3000",             "region\_name": "West",             "channel\_code": "000",             "channel\_name": "Non Chanel",             "area\_code": "08",             "area\_name": "Central Jawa",             "province\_code": "36",             "province\_name": "JAWA BARAT",             "regency\_code": "3737",             "regency\_name": "Bogor",             "sub\_district\_code": "363636",             "sub\_district\_name": "Cileungsi",             "ward\_code": "123456",             "ward\_name": "Jaya Mekar",             "is\_active": true,             "is\_del": false,             "createdby": null,             "created\_at": null,             "updatedby": null,             "updated\_at": null,             "updated\_by\_name": "Admin Principal 1",             "deleted\_by": null,             "deleted\_at": null,             "customer\_id": ""         }     \],     "paging": {         "total\_record": 2,         "page\_current": 1,         "page\_limit": 10,         "page\_total": 1     },     "request\_id": "696747a251f841f20353ce57" }  |
| :---- |

### 

2. # **Distributor Detail** 

Enhance endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: [https://best.scyllax.online/master/v1/distributors/68](https://best.scyllax.online/master/v1/distributors/68) 

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| q |  |  |  |
| page | Integer | Yes | default : 1 |
| limit  | Integer | Yes | default : 5  |

### **Example Request** 

| curl \--location \-g \\ '{{url}}/master/v1/outlet\_list?page=1\&limit=5\&sort=created\_date:desc\&status=1' \\ \--header 'Accept: application/json' |
| :---- |

### **Add Response :**

| Field | Type | Note |
| :---- | :---- | :---- |
| **distributor\_setup** |  |  |
| allow\_add\_product | Bool | **mst.m\_distributor.**allow\_add\_product |
| allow\_edit\_product | Bool | **mst.m\_distributor.**allow\_edit\_product |
| allow\_manage\_pricing | Bool | **mst.m\_distributor.**allow\_manage\_pricing |
| allow\_upload\_secondary\_sales | Bool | **mst.m\_distributor.**allow\_upload\_secondary\_sales |
|  |  |  |

### **Response Example Response :** 

| {     "message": "",     "data": {         "distributor\_id": 68,         "distributor\_code": "232323",         "distributor\_name": "Pratesis Bangka",         "barcode": "",         "region\_id": 67,         "area\_id": 82,         "channel\_id": 44,         "sub\_distributor\_group\_id": 36,         "sub\_distributor\_group\_code": "000",         "sub\_distributor\_group\_name": "NON",         "dist\_price\_grp\_id": 38,         "address": "jl kedondong",         "province\_id": "36",         "regency\_id": "3737",         "sub\_district\_id": "363636",         "ward\_id": "3366363",         "zip\_code": "62134",         "ot\_loc\_id": 0,         "latitude": "3434343",         "longitude": "34343",         "phone": "342342343423",         "fax\_number": "",         "dist\_price\_grp\_code": "003",         "dist\_price\_grp\_name": "JAWA",         "region\_code": "2000",         "region\_name": "Central",         "channel\_code": "000",         "channel\_name": "Non Chanel",         "area\_code": "08",         "area\_name": "Central Jawa",         "province\_code": "36",         "province\_name": "JAWA BARAT",         "regency\_code": "3737",         "regency\_name": "Bogor",         "sub\_district\_code": "363636",         "sub\_district\_name": "Cileungsi",         "ward\_code": "3366363",         "ward\_name": "Cileungsi",         "is\_active": true,         "updated\_by": 1,         "updated\_by\_name": "Admin Principal 1",         "updated\_at": "2025-10-16T02:57:25.099069Z",         "contacts": \[             {                 "distributor\_contact\_id": 91,                 "distributor\_id": 68,                 "contact\_name": "Nasywa Syafinka Widyamara",                 "job\_title": "hr",                 "phone\_no": "08584545454455",                 "is\_wa\_no": true,                 "wa\_no": "08584545454455",                 "email": "syafinkaforcolle@gmail.com",                 "identity\_no": "45454545",                 "identity\_type": "National ID"             }         \],         "tax": \[             {                 "distributor\_tax\_id": 27,                 "distributor\_id": 68,                 "tax\_identifier\_no\_type": "National ID",                 "tax\_identifier\_no": "0000000000000000",                 "nitku": "000000",                 "tax\_name": "fg",                 "tax\_address": "jl"             }         \],         "distributor\_setup" : {             "allow\_add\_product": true,             "allow\_edit\_product": true,             "allow\_manage\_pricing": false,             "allow\_upload\_secondary\_sales": true         }     },     "request\_id": "6967484051f841f20353ce94" }  |
| :---- |

3. # **Distributor Create** 

Enhance endpoint :  
Content-Type 	: application/json  
Method		: POST  
URL		: https://best.scyllax.online/master/v1/distributors 

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Body** 
---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| distributor\_code | String | YES |  |
| distributor\_name | String | YES |  |
| barcode | String | NO |  |
| region\_id | Integer | YES |  |
| area\_id | Integer | YES |  |
| channel\_id | Integer | YES |  |
| sub\_distributor\_group\_id | Integer | YES |  |
| dist\_price\_grp\_id | Integer | YES |  |
| address | String | YES |  |
| province\_id | Integer | NO |  |
| regency\_id | String | NO |  |
| sub\_district\_id | String | NO |  |
| ward\_id | String | NO |  |
| zip\_code | String | NO |  |
| ot\_loc\_id | Integer | YES |  |
| latitude | String | YES |  |
| longitude | String | YES |  |
| phone | String | NO |  |
| fax\_number | String | NO |  |
| is\_active | Boolean | YES |  |
| **contact** |  |  |  |
| contact\_name | String | YES |  |
| phone\_no | String | YES |  |
| job\_title | String | YES |  |
| is\_wa\_no | Boolean | YES |  |
| wa\_no | String | YES |  |
| email | String | NO |  |
| identity\_no | String | YES |  |
| identity\_type | String | YES |  |
| **tax** |  |  |  |
| tax\_identifier\_no\_type | String | YES |  |
| tax\_identifier\_no | String | YES |  |
| nitku | String | YES |  |
| tax\_name | String | YES |  |
| tax\_address | String | YES |  |
| **distributor\_setup** |  |  |  |
| allow\_add\_product | Bool | YES |  |
| allow\_edit\_product | Bool | YES |  |
| allow\_manage\_pricing | Bool | YES |  |
| allow\_upload\_secondary\_sales | Bool | YES |  |

### **Example Request** 

| {   "distributor\_code": "D00232",   "distributor\_name": "Distributor Geek 1",   "barcode": "12322121223232",   "region\_id": 68,   "area\_id": 85,   "channel\_id": 44,   "sub\_distributor\_group\_id": 36,   "dist\_price\_grp\_id": 44,   "address": "Jalan Kaliurang ",   "province\_id": 0,   "regency\_id": "0101",   "sub\_district\_id": "353535",   "ward\_id": "Not set",   "zip\_code": "78766",   "ot\_loc\_id": 0,   "latitude": "106.827153",   "longitude": "-6.175392",   "phone": "08786767878",   "fax\_number": "08786767878",   "is\_active": true,   "contacts": \[     {       "contact\_name": "Widya ",       "phone\_no": "0878655878",       "job\_title": "CEO",       "is\_wa\_no": true,       "wa\_no": "0878655878",       "email": "widya@gmail.com",       "identity\_no": "213208098000009",       "identity\_type": "National ID"     }   \],   "tax": \[     {       "tax\_identifier\_no\_type": "TIN",       "tax\_identifier\_no": "233222121",       "nitku": "212133",       "tax\_name": "tax A",       "tax\_address": "Jalan kemana "     }   \],   "distributor\_setup" : {      "allow\_add\_product": true,     "allow\_edit\_product": true,     "allow\_manage\_pricing": false,     "allow\_upload\_secondary\_sales": true   } }  |
| :---- |

### **Response :** 

| *resp menyesuaikan eksisting (tidak ada perubahan)* |
| :---- |

### **Impact ke db  :** 

***saat insert ke tabel mst.m\_distributor, tambahkan data :*** 

- allow\_add\_product  
- allow\_edit\_product  
- allow\_manage\_pricing  
- allow\_upload\_secondary\_sales

4. # **Distributor Update** 

Enhance endpoint :  
Content-Type 	: application/json  
Method		: PATCH  
URL		: https://best.scyllax.online/master/v1/distributors 

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Body** 
---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| distributor\_code | String | YES |  |
| distributor\_name | String | YES |  |
| barcode | String | NO |  |
| region\_id | Integer | YES |  |
| area\_id | Integer | YES |  |
| channel\_id | Integer | YES |  |
| sub\_distributor\_group\_id | Integer | YES |  |
| dist\_price\_grp\_id | Integer | YES |  |
| address | String | YES |  |
| province\_id | Integer | NO |  |
| regency\_id | String | NO |  |
| sub\_district\_id | String | NO |  |
| ward\_id | String | NO |  |
| zip\_code | String | NO |  |
| ot\_loc\_id | Integer | YES |  |
| latitude | String | YES |  |
| longitude | String | YES |  |
| phone | String | NO |  |
| fax\_number | String | NO |  |
| is\_active | Boolean | YES |  |
| **contact** |  |  |  |
| contact\_name | String | YES |  |
| phone\_no | String | YES |  |
| job\_title | String | YES |  |
| is\_wa\_no | Boolean | YES |  |
| wa\_no | String | YES |  |
| email | String | NO |  |
| identity\_no | String | YES |  |
| identity\_type | String | YES |  |
| **tax** |  |  |  |
| tax\_identifier\_no\_type | String | YES |  |
| tax\_identifier\_no | String | YES |  |
| nitku | String | YES |  |
| tax\_name | String | YES |  |
| tax\_address | String | YES |  |
| **distributor\_setup** |  |  |  |
| allow\_add\_product | Bool | YES |  |
| allow\_edit\_product | Bool | YES |  |
| allow\_manage\_pricing | Bool | YES |  |
| allow\_upload\_secondary\_sales | Bool | YES |  |

| {   "distributor\_code": "232323",   "distributor\_name": "Pratesis Bangka 333",   "barcode": "",   "region\_id": 67,   "area\_id": 82,   "channel\_id": 44,   "sub\_distributor\_group\_id": 36,   "dist\_price\_grp\_id": 38,   "address": "jl kedondong",   "province\_id": "36",   "regency\_id": "3737",   "sub\_district\_id": "363636",   "ward\_id": "3366363",   "zip\_code": "62134",   "ot\_loc\_id": 0,   "latitude": "3434343",   "longitude": "34343",   "phone": "342342343423",   "fax\_number": "",   "is\_active": true,   "contacts": \[     {       "distributor\_contact\_id": 91,       "distributor\_id": 68,       "contact\_name": "Nasywa Syafinka Widyamara",       "job\_title": "hr",       "phone\_no": "08584545454455",       "is\_wa\_no": true,       "wa\_no": "08584545454455",       "email": "syafinkaforcolle@gmail.com",       "identity\_no": "45454545",       "identity\_type": "National ID"     }   \],   "tax": \[     {       "distributor\_tax\_id": 27,       "distributor\_id": 68,       "tax\_identifier\_no\_type": "National ID",       "tax\_identifier\_no": "0000000000000000",       "nitku": "000000",       "tax\_name": "fg",       "tax\_address": "jl"     }   \],   "distributor\_setup" : {      "allow\_add\_product": true,     "allow\_edit\_product": true,     "allow\_manage\_pricing": false,     "allow\_upload\_secondary\_sales": true   } }  |
| :---- |

