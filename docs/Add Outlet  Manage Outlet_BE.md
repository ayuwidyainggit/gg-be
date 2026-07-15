1. # **Upload File** 

   API eksisting : {{kong\_url}}/mobile/v1/files/uploads 

2. # **Enhance Create New Outlet**

* Add field file\_url type text di mst.m\_outlet  
* Enhance API: simpan URL hasil upload file ke mst.m\_outlet.file\_url.

Enhance  
Content-Type 	: application/json  
Method		: POST  
URL		: {{url}}[/mobile/v1/m-outlets](http://103.28.219.73:5001/mobile/v1/m-outlets)

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

**Req Body :**  Tambahkan file\_url  
---

| {   "outlet\_name": "Outlet1",   "created\_by": 120,   "address": "Jalan Nangka III Kecamatan Depok Maguwoharjo Kabupaten Sleman 55281",   "phone\_no": "098977778888",   "building\_own": 1,   "latitude": "-7.7671403",   "longitude": "110.4243715",   "file\_url": "diisi url resp upload file ",   "details": {     "contact": \[       {         "contact\_name": "jojo",         "job\_title": "Outlet Manager",         "phone\_no": "0878655666",         "wa\_no": "0878655666",         "email": "jojo@gmail.com",         "identity\_no": null       }     \]   } }  |
| :---- |

* Saat save, input url pada mst.m\_outlet.file\_url

3. # **New Outlet List** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/mobile/v1/outley-list?page=1\&limit=10\&is\_active=0

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| q |  |  | search by outlet\_name dan outlet\_code |
| page | Integer | Yes | default : 1 |
| limit  | Integer | Yes | default : 9999 |
| sort | String | Yes | default outlet\_code:asc  |
| verification\_status | Array\<int\> | Yes | **1, 3** ![][image1] |

### **Example Request** 

---

**Example Request Default**

| curl \--location '{{url}}/mobile/v1/outlet-list?q=\&page=1\&limit=10\&sort=outlet\_code:asc\&verification\_status=1,3' \\   \--header 'Accept: application/json' \\   \--header 'Authorization: Bearer {{token}}' |
| :---- |

### **Response** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| outlet\_id | int | 8 | **mst.m\_outlet.**outlet\_id |
| outlet\_code | varchar | 30 | **mst.m\_outlet.**outlet\_code |
| outlet\_name | varchar | 150 | **mst.m\_outlet.**outlet\_name |
| address1 | varchar | 150 | **mst.m\_outlet.**address1 |
| longitude | varchar | 50 | **mst.m\_outlet.**longitude |
| latitude | varchar | 50 | **mst.m\_outlet.**latitude |
| verification\_status | int |  | **mst.m\_outlet.**verification\_status |

**Example Response :** 

| Case : sukses dan terdapat data {  "message": "success",  "data": {    "distributor\_id": "123",    "distributor\_code": "DST001",    "distributor\_name": "PT Distributor Maju",    "address": "Jl. Contoh No. 123, Jakarta",    "latitude": \-6.2,    "longitude": 106.816666  },  "paging": {    "total\_record": 100,    "page\_current": 1,    "page\_limit": 10,    "page\_total": 10  },  "request\_id": "req-abc-123456" }  |
| :---- |
| Case : empty state / tidak terdapat data berdasarkan pencarian {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" } |

4. # **Outlet Detail** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/mobile/v1/outley-list/:outlet\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Path Variable** 
---

|  outlet\_id | 32 |
| :---- | :---- |

### **Example Request** 

---

| curl \--location '{{url}}/mobile/v1/outlet-list/123' \\   \--header 'Accept: application/json' \\   \--header 'Authorization: Bearer {{token}}' |
| :---- |

### **Response** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| outlet\_id | int | 8 | **mst.m\_outlet.**outlet\_id |
| outlet\_code | varchar | 30 | **mst.m\_outlet.**outlet\_code |
| outlet\_name | varchar | 150 | **mst.m\_outlet.**outlet\_name |
| address1 | varchar | 150 | **mst.m\_outlet.**address1 |
| phone\_no | varchar | 20 | **mst.m\_outlet.**phone\_no |
| building\_own | int | 2 | **mst.m\_outlet.**building\_own |
| file\_url | text |  | **mst.m\_outlet.**file\_url |
| longitude | varchar | 50 | **mst.m\_outlet.**longitude |
| latitude | varchar | 50 | **mst.m\_outlet.**latitude |
| **other\_contact** | **Object** |  |  |
| contact\_name | varchar | 150 | **mst.m\_outlet\_contact.**contact\_name |
| job\_title | varchar | 100 | **mst.m\_outlet\_contact.**job\_title |
| phone\_no | varchar | 20 | **mst.m\_outlet\_contact.**phone\_no |
| wa\_no | varchar | 20 | **mst.m\_outlet\_contact.**wa\_no |
| email | varchar | 100 | **mst.m\_outlet\_contact.**email |

**Example Response :** 

| Case : sukses dan terdapat data {   "message": "success",   "data": \[     {       "outlet\_id": 1,       "outlet\_code": "OTL001",       "outlet\_name": "Outlet Jaya Abadi",       "address1": "Jl. Merdeka No.10, Jakarta",       "phone\_no": "02188990011",       "building\_own": 1,       "file\_url": "https://example.com/file/outlet1.jpg",       "longitude": "106.827153",       "latitude": "-6.175392",       "other\_contact": {         "contact\_name": "Budi Santoso",         "job\_title": "Owner",         "phone\_no": "081234567890",         "wa\_no": "081234567890",         "email": "budi@example.com"       }     }   \] }  |
| :---- |
| {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

5. # **Remove Outlet List** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: DELETE  
URL		: {{url}}/mobile/v1/outley-list/:outlet\_id

### **Headers**

---

| Accept | multipart/form-data |
| :---- | :---- |
| Authorization  | Bearer Token |

**Path Variable** 

---

|  expense\_id | Array\<integer\> |
| :---- | :---- |

### **Example Request** 

---

| curl \--location \--request DELETE '{{url}}/mobile/v1/outlet-list?outlet\_id=1\&outlet\_id=2\&outlet\_id=3\&outlet\_id=4' \\   \--header 'Accept: application/json' \\   \--header 'Authorization: Bearer {{token}}'  |
| :---- |

update di bagian field (mst.m\_outlet) berikut : 

- is\_del  
- deleted\_by  
- deleted\_at

6. # **Edit Outlet List** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: PATCH  
URL		: /mobile/v1/m-outlets/:outlet\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

**Req Body :**  Kirim data yang hanya di edit saja 

| curl \--location \--request PATCH '{{url}}/mobile/v1/m-outlets/1' \\   \--header 'Accept: application/json' \\   \--header 'Authorization: Bearer {{token}}' \\   \--header 'Content-Type: application/json' \\   \--data '{     "outlet\_name": "Outlet1",     "created\_by": 120,     "address": "Jalan Nangka III Kecamatan Depok Maguwoharjo Kabupaten Sleman 55281",     "phone\_no": "0989777888",     "building\_own": 1,     "latitude": "-7.7671403",     "longitude": "110.4243715",     "file\_url": "diisi url resp upload file",     "details": {       "contact": \[         {           "contact\_name": "jojo",           "job\_title": "Outlet Manager",           "phone\_no": "0878655666",           "wa\_no": "0878655666",           "email": "jojo@gmail.com",           "identity\_no": null         }       \]     }   }'   |
| :---- |

**Example Response :** 

| Case : sukses dan terdapat data {    "message": "Data successfully updated",    "request\_id": "6915a6dd2395083c685e8e16" }  |
| :---- |
| case tdk ada perubahan data   |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" } |
