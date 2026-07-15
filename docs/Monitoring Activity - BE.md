1. # **Region** 

 Create new endpoint   
 Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}master/v1/regions?q=\&page=1\&limit=10\&sort=region\_id:asc\&is\_active=1\&region\_id=1

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| q |  |  |  |
| page | Integer | Yes | default : 1 |
| limit  | Integer | Yes | default : 9999 |
| sort | String | Yes | default region\_id:asc  |
| is\_active | Integer | Yes | 1 |
| region\_id | Array\<Integer\> | Yes | 1 |

### 

**Example Request** 

| curl \--location \-g \\   '{{url}}/master/v1/regions?q=\&page=1\&limit=10\&sort=region\_id:asc\&is\_active=1\&region\_id\[\]=1\&region\_id\[\]=2\&region\_id\[\]=3' \\   \--header 'Accept: application/json' \\   \--header 'Authorization: Bearer {{token}}' |
| :---- |

### **Response** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| region\_id | Integer | 8 | **mst.m\_region.**region\_id |
| region\_code | Varchar | 10 | **mst.m\_region.**region\_code |
| region\_name | Varchar | 150 | **mst.m\_region.**region\_name |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data { "message": "", "data" : \[  {   "region\_id" : 12,   "region\_code": "R001",   "region\_name": "Jakarta"  } \], "paging" : {   "total\_record": 84,   "page\_current": 1,   "page\_limit": 10,   "page\_total": 9 }, "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| **Case : empty state / tidak terdapat data berdasarkan pencarian** {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| **Case : error ** {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" } |

### **Changes :** 

* ***26 Mei 2026*** 

Khusus untuk user principal, perhatikan pada mst.m\_employee   
pada mst.m\_employee terdapat changes :   
 \- \`region\_scope\`  
 \- \`area\_scope\`  
 \- \`distributor\_scope\`

**Rule untuk Dropdown Region Berdasarkan region\_scope**

Sistem harus mengecek nilai field **region\_scope** pada tabel mst.m\_employee.

***1\. region\_scope \= 'Specific'***

* Sistem harus mengambil data region dari tabel mst.**m\_employee\_region\_mapping** berdasarkan user yang sedang login.  
* Region yang ditampilkan pada dropdown hanya region yang terdaftar pada mapping user tersebut.  
  ***2\. region\_scope \= 'All' or NULL***   
* User tidak memiliki data mapping pada tabel mst.**m\_employee\_region\_mapping**.  
* Region yang ditampilkan pada dropdown adalah seluruh region berdasarkan cust\_id dari user yang sedang login.

 

2. # **Business Unit** 

 Create new endpoint   
 Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}master/v1/business-unit?q=\&page=1\&limit=9999\&sort=area\_id:asc\&is\_active=1\&region\_id=1

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| q |  |  |  |
| page | Integer | Yes | default : 1 |
| limit  | Integer | Yes | default : 9999 |
| user\_name | varchar | No | dari username user yg login  princ@idetama.id |
| distributor\_id | int | No | cek di local storage  jika NULL or 0 → user principal  Jika Not NULL → user distributor |
| cust\_id | varchar | No |  ***biasanya langsung include dr token (untuk filter berdasarkan principal)***  |
| region\_id | int | No | **before** : int  **Enhance** :  Array\<int\> |
| area\_id | int | No | **before** : int  **Enhance** :  Array\<int\> |
| is\_active | Array\<int\> | Yes | 0 : inactive 1 : active null : show all status // disini FE kirim 1  |

### 

* Contoh CURL untuk user distributor : 


| curl 'https://best.scyllax.online/master/v1/business-unit?is\_active=1\&distributor\_id=67\&page=1\&limit=70'' \\   \-H 'sec-ch-ua-platform: "macOS"' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsImVtcF9pZCI6MCwiZXhwaXJlcyI6MTc3MTU1NzU3MywiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4NTc0ODEyMzEyOCIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IkRpc3QgSURFIFNkYSIsInVzZXJfaWQiOjEyLCJ1c2VyX25hbWUiOiJkaXN0QHNkYS5pZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODU3NDgxMjMxMjgifQ.jOznFNEW7P1xGO80X755W3XDFiVrAPMcmauiHY44b9E' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'sec-ch-ua: "Not:A-Brand";v="99", "Google Chrome";v="145", "Chromium";v="145"' \\   \-H 'sec-ch-ua-mobile: ?0' yang perlu di kirim FE :  distributor\_id \= Local Storage distributor\_id  is\_active \= 1 page  limit  |
| :---- |


* Contoh CURL untuk user principal 


| curl 'https://best.scyllax.online/master/v1/business-unit?is\_active=1\&q=\&page=1\&limit=70' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxIiwiZGlzdF9wcmljZV9ncnBfaWQiOjAsImRpc3RyaWJ1dG9yX2lkIjowLCJlbWFpbCI6InByaW5jQGlkZXRhbWEuaWQiLCJlbXBfaWQiOjI3OCwiZXhwaXJlcyI6MTc3MTU1NzcwMCwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4MTEzMjMyMzMyIiwicGFyZW50X2N1c3RfaWQiOiJDMjIwMDEiLCJ1c2VyX2Z1bGxuYW1lIjoiQWRtaW4gUHJpbmNpcGFsIDEiLCJ1c2VyX2lkIjoxLCJ1c2VyX25hbWUiOiJwcmluY0BpZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODExMzIzMjMzMiJ9.qrNJ6QKRCA2\_bHRxDEWugtO-sHw6MRcdyH1EmQOd7qE' \\   \-H 'Connection: keep-alive' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Not:A-Brand";v="99", "Google Chrome";v="145", "Chromium";v="145"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"' yang perlu di kirim FE :  distributor\_id \= Local Storage distributor\_id pasti mengembalikan null , sehingga tdk perlu dikirim is\_active \= 1 page  limit  |
| :---- |


### **Response** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Object | \- |  \- |
| user\_id | Integer | 8 | **sys.m\_user.**user\_id user principal yg login |
| user\_fullname | varchar | 150 | **sys.m\_user.**user\_fullname user principal yg login |
| distributor\_id | int | 8 | **smc.m\_customer.**distributor\_id // jika NULL → user principal     jika NOT NULL → user distributor |
| distributor\_data | Array |  |  |
| distributor\_id | Integer | 8 | **mst.m\_distributor.**distributor\_id |
| distributor\_code | Varchar | 20 | **mst.m\_distributor.**distributor\_code |
| distributor\_name | varchar | 150 | **mst.m\_distributor.**distributor\_name |
| area\_id | int | 8 | **mst.m\_distributor.**area\_id |
| area\_code | varchar | 10 | **mst.m\_area.**area\_code |
| area\_name | varchar | 150 | **mst.m\_area.**area\_name |
| region\_id | int | 8 | **mst.m\_distributor.**region\_id |
| region\_code | varchar | 10 | **mst.m\_region.**region\_code |
| region\_name | varchar | 150 | **mst.m\_region.**region\_name |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

#### **Example Response User Principal :**  

**SELECT** \*  
**FROM** mst.m\_distributor *md*  
**WHERE** *md*.cust\_id **LIKE** 'C26002%' **and** *md*.is\_active \= **true** **and** *md*.is\_del **is** **false** **and** *md*.is\_active \= **true** **and** *md*.is\_del  \= **false**

atau 

**SELECT** \*  
**FROM** mst.m\_distributor *md*  
**where** *md*.parent\_cust\_id \='C26002' **and** *md*.is\_active \= **true** **and** *md*.is\_del  \= **false**

| Case : sukses dan terdapat data  {  "message": "Success",  "data": {    "cust\_id": "C26002", //sys.m\_user.cust\_id    "user\_id": 101,    "user\_fullname": "Widya Ayu",    "distributor\_id": “”,    "distributor\_data": \[      {        "cust\_id": "C260020001", //mst.m\_distributor.cust\_id        "distributor\_id": 23,        "distributor\_code": "DST001",        "distributor\_name": "PT Sumber Makmur",        "area\_id": 12,        "area\_code": "A001",        "area\_name": "Jakarta",        "region\_id": 1,        "region\_code": "R001",        "region\_name": "Jabodetabek"      },      {        "cust\_id": "C260020002",        "distributor\_id": 24,        "distributor\_code": "DST002",        "distributor\_name": "PT Sejahtera Abadi",        "area\_id": 13,        "area\_code": "A002",        "area\_name": "Bandung",        "region\_id": 2,        "region\_code": "R002",        "region\_name": "Jawa Barat"      }    \]  },  "paging": {    "total\_record": 2,    "page\_current": 1,    "page\_limit": 10,    "page\_total": 1  },  "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |

#### **Example Response User Distributor :** 

**select** \* **from** mst.m\_distributor *md*  **where** *md*.distributor\_id \='68'

| Case : sukses dan terdapat data { "message": "Success", "data": {   "user\_id": 101,   "user\_fullname": "Widya Ayu",   "user\_principal": "Principal A Company",   "distributor\_id": 21,   "distributor\_code": "DST001",   "distributor\_name": "PT Sumber Makmur",   "area\_id": 12,   "area\_code": "A001",   "area\_name": "Jakarta",   "region\_id": 1,   "region\_code": "R001",   "region\_name": "Jabodetabek" }, "paging": {   "total\_record": 2,   "page\_current": 1,   "page\_limit": 10,   "page\_total": 1, }, "request\_id": "6915a5e8e3f53f84fe73517f" } |
| :---- |
| **Case : empty state / tidak terdapat data berdasarkan pencarian** {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| **Case : error ** {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }   |

### **Changes :** 

* 26 Mei 2026 

Khusus untuk user principal, perhatikan pada mst.m\_employee   
pada mst.m\_employee terdapat changes :   
 \- \`region\_scope\`  
 \- \`area\_scope\`  
 \- \`distributor\_scope\`

**Rule untuk Dropdown Business Unit** 

Jika pada payload region\_id dan area\_id \= null, maka sistem harus mengecek nilai field distributor\_scope pada mst.m\_employee.

***1\. distributor\_scope \= 'Specific'***

* Sistem harus mengambil data distributor dari tabel mst.**m\_employee\_distributor\_mapping** berdasarkan user yang sedang login.  
* Distributor yang ditampilkan pada dropdown adalah data yang terkait dengan distributor mapping milik user tersebut.

  ***2\. distributor\_scope \= 'All'  or NULL dan region\_scope \= ‘All’ or NULL dan area\_scope \= ‘All’ or NULL*** 

* User tidak memiliki data pada tabel mst.**m\_employee\_distributor\_mapping**.  
* Distributor yang ditampilkan pada dropdown adalah seluruh data berdasarkan data ***parent\_cust\_id*** dari user yang sedang login. 

  ***3\. distributor\_scope \= 'All' or NULL  namun region\_scope \= ‘SPESIFIC’ dan area\_scope \= ‘SPESIFIC’***

* User tidak memiliki data pada tabel mst.**m\_employee\_distributor\_mapping**.  
* Data yang ditampilkan pada dropdown adalah seluruh data distributor berdasarkan area\_id dan region\_id pada master employee.

  ***4\. distributor\_scope \= 'All'  or NULL namun region\_scope \= ‘ALL’ or NULL  dan area\_scope \= ‘SPESIFIC’***

* User tidak memiliki data pada tabel mst.**m\_employee\_distributor\_mapping**.  
* Data yang ditampilkan pada dropdown adalah seluruh data distributor berdasarkan area\_id pada master employee.

  ***5\. distributor\_scope \= 'All' or NULL namun region\_scope \= ‘SPESIFIC’ dan area\_scope \= ‘ALL’ or NULL***

* User tidak memiliki data pada tabel mst.**m\_employee\_distributor\_mapping**.  
* Data yang ditampilkan pada dropdown adalah seluruh data distributor berdasarkan region\_id pada master employee.


3. # **Area**

**eksisting**   
 Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}master/v1/areas?q=\&page=1\&limit=9999\&sort=area\_id:asc\&is\_active=1\&region\_id=1

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| q |  |  |  |
| page | Integer | Yes | default : 1 |
| limit  | Integer | Yes | default : 9999 |
| sort | String | Yes | default region\_id:asc  |
| is\_active | Integer | Yes | 1 |
| area\_id | Array\<Integer\> | Yes | 1 |

### 

**Example Request** 

| curl \--location \-g \\   '{{url}}/master/v1/area?q=\&page=1\&limit=10\&sort=area\_id:asc\&is\_active=1\&area\_id\[\]=1\&area\_id\[\]=2\&area\_id\[\]=3' \\   \--header 'Accept: application/json' \\   \--header 'Authorization: Bearer {{token}}'  |
| :---- |

### **Response** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| area\_id | Integer | 8 | **mst.m\_area.**area\_id |
| area\_code | Varchar | 10 | **mst.m\_area.**area\_code |
| area\_name | Varchar | 150 | **mst.m\_area.**area\_name |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data { "message": "", "data" : \[  {   "area\_id" : 12,   "area\_code": "A001",   "area\_name": "Jakarta"  } \], "paging" : {   "total\_record": 84,   "page\_current": 1,   "page\_limit": 10,   "page\_total": 9 }, "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| **Case : empty state / tidak terdapat data berdasarkan pencarian** {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| **Case : error ** {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" } |

### **Changes :** 

* 26 Mei 2026 

Khusus untuk user principal, perhatikan pada mst.m\_employee   
pada mst.m\_employee terdapat changes :   
 \- \`region\_scope\`  
 \- \`area\_scope\`  
 \- \`distributor\_scope\`

**Rule untuk Dropdown Area**

Jika pada payload region\_id \= null, maka sistem harus mengecek nilai field **area\_scope** pada mst.m\_employee.

***1\. area\_scope \= 'Specific'*** 

* Sistem harus mengambil data area dari tabel mst.m\_employee\_area\_mapping berdasarkan user yang sedang login.  
* Area yang ditampilkan pada dropdown adalah Area yang terkait dengan area mapping milik user tersebut.

  ***2\. area\_scope \= 'All' or NULL dan region\_scope \= ‘All’ or NULL***

* User tidak memiliki data pada tabel mst.**m\_employee\_area\_mapping**.  
* Area yang ditampilkan pada dropdown adalah seluruh area berdasarkan cust\_id dari user yang sedang login.

  ***3\. area\_scope \= 'All' or ‘NULL’ namun region\_scope ‘SELECTED’***

* User tidak memiliki data pada tabel mst.**m\_employee\_area\_mapping**.  
* Area yang ditampilkan pada dropdown adalah seluruh area berdasarkan cust\_id dan region\_id pada master employee (***m\_employee\_region\_mapping***)


4. # **Employee**

 Create new endpoint   
 Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/master/v1/employee-pjp?q=\&page=1\&limit=9999\&sort=area\_id:asc\&is\_active\[\]=1\&region\_id\[\]=1\&cust\_id=10\&distributor\_id=5

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| q |  |  |  |
| page | Integer | Yes | default : 1 |
| limit  | Integer | Yes | default : 9999 |
| sort | String | Yes | default region\_id:asc  |
| status | Array\<int\> | Yes | 0 : inactive 1 : active null : show all status // disini FE kirim 1  |
| cust\_id | Integer | No | Jika user principal memilih user principal itu sendiri → ditampilkan employee dari cust\_id user tersebut |
| distributor\_id | Integer | No | Jika user principal memilih dropdown distributor → ditampilkan employee dari distributor tersebut |

### 

**Example Request :**

| curl \--location \-g \\   '{{url}}/master/v1/employee?q=\&page=1\&limit=999\&sort=area\_id:asc\&is\_active\[\]=1\&region\_id\[\]=1\&cust\_id=10\&distributor\_id=5' \\   \--header 'Accept: application/json' \\   \--header 'Authorization: Bearer {{token}}'  |
| :---- |

### **Query:** 

### ---

* jika filter berdasarkan distributor : (**cust\_id \> 6 digit)**

| select *me*.\*, *mr*.role\_name  from mst.m\_employee *me*   join sys.m\_user *mu*  	on *mu*.emp\_id \= *me*.emp\_id and *me*.cust\_id \= *mu*.cust\_id join sys.user\_roles *ur*  		on *ur*.user\_id  \= *mu*.user\_id  and *ur*.cust\_id \= *mu*.cust\_id  	join sys.m\_role *mr*  		on *ur*.role\_id  \= *mr*.role\_id and *mr*.cust\_id  \= *ur*.cust\_id  where  *me*.cust\_id in ('C260020002', 'C260020001') and *me*.is\_active \= true and *me*.is\_del \= false and *mr*.role\_name \='salesman'  |
| :---- |


* jika filter berdasarkan user principal itu sendiri ((**cust\_id \= 6 digit)**

| select *me*.\*, *mr*.role\_name  from mst.m\_employee *me*   join sys.m\_user *mu*  	on *mu*.emp\_id \= *me*.emp\_id and *me*.cust\_id \= *mu*.cust\_id join sys.user\_roles *ur*  		on *ur*.user\_id  \= *mu*.user\_id  and *ur*.cust\_id \= *mu*.cust\_id  	join sys.m\_role *mr*  		on *ur*.role\_id  \= *mr*.role\_id and *mr*.cust\_id  \= *ur*.cust\_id  where  *me*.cust\_id in ('C26002') and *me*.is\_active \= true and *me*.is\_del \= false and *mr*.role\_name \='salesman'   |
| :---- |


  

* jika filter berdasarkan user principal itu sendiri (cust\_id) dan user distributor 

| select *me*.\*, *mr*.role\_name  from mst.m\_employee *me*   join sys.m\_user *mu*  	on *mu*.emp\_id \= *me*.emp\_id and *me*.cust\_id \= *mu*.cust\_id join sys.user\_roles *ur*  		on *ur*.user\_id  \= *mu*.user\_id  and *ur*.cust\_id \= *mu*.cust\_id  	join sys.m\_role *mr*  		on *ur*.role\_id  \= *mr*.role\_id and *mr*.cust\_id  \= *ur*.cust\_id  where  *me*.cust\_id in ('C26002','C260020002', 'C260020001') and *me*.is\_active \= true and *me*.is\_del \= false and *mr*.role\_name \='salesman'  |
| :---- |


### 

### 

### **Response:** 

### ---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| emp\_id | Integer | 8 | **mst.m\_employee.**area\_id |
| emp\_code | Varchar | 10 | **mst.m\_employee.**area\_code |
| emp\_name | Varchar | 150 | **mst.m\_employee.**area\_name |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data { "message": "", "data" : \[  {   "emp\_id" : 12,   "emp\_code": "A001",   "emp\_name": "Jakarta"  } \], "paging" : {   "total\_record": 84,   "page\_current": 1,   "page\_limit": 10,   "page\_total": 9 }, "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| **Case : empty state / tidak terdapat data berdasarkan pencarian** {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| **Case : error ** {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }   |

5. # **Location Monitoring (Principal)**

 Create new endpoint   
    
 Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}scylla-pjp/v1/live-monitoring-principal

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Param** 
---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| date | epoch  | Yes | pjp\_principles.route\_pop\_permanent.date |
| emp\_id | Array\<int\> | No | pjp\_principles.permanent\_journey\_plans.salesman\_id |
| status | Array\<varchar\> | Yes | untuk saat ini FE kirim ‘Approved’ BE filter pada : **pjp\_principles.**permanent\_journey\_plans.approval\_status |

**Example Request :** 

| curl \--location \-g \\   '{{url}}scylla-pjp/v1/live-monitoring-principal?date=1738195200\&emp\_id\[\]=221\&emp\_id\[\]=279\&emp\_id\[\]=264\&status\[\]=Approved' \\   \--header 'Accept: application/json' \\   \--header 'Authorization: Bearer {{token}}' |
| :---- |

### **Response** 

**Query :** 

| select  	*pjp*.id, 	*pjp*.pjp\_code , 	*r*.route\_code , 	*r*.route\_name ,  	*pjp*.approval\_status ,  	*pjp*.salesman\_name,  	*dh*.id as *destination\_id*, 	*dh*.destination\_code , 	*dh*.destination\_type , 	*dh*.destination\_id , 	*mo*.outlet\_name  ,     *v*.longitude as *arrive\_longitude*,      *v*.latitude as *arrive\_latitude*, 	*ovl*.outlet\_id as *ovl\_outlet\_id*, 	*mo*.outlet\_code , *mo*.outlet\_name , *mo*.address1 , 	*md*.distributor\_code , *md*.distributor\_name ,     *dh*.id , *dh*.longitude , *dh*.latitude ,      *ovl*."start" ,     *ovl*.finish ,     *ovl*.skip\_at  from pjp\_principles.destinations\_history *dh*  join pjp\_principles.permanent\_journey\_plans *pjp*  	on *pjp*.id \= *dh*.pjp\_id  and *dh*.cust\_id \= *pjp*.cust\_id  left join mst.m\_outlet *mo* on *mo*.outlet\_id \= *dh*.destination\_id  left join pjp\_principles.routes *r*  	on *r*.route\_code \= *dh*.route\_code  left join mst.m\_distributor *md*  	on *md*.distributor\_id \= *dh*.destination\_id  left join pjp\_principles.outlet\_visit\_list *ovl*  	on *ovl*.outlet\_id \= *dh*.destination\_id  	and  *ovl*.pjp\_id \= *dh*.pjp\_id  	and ovl."date"::date \= '2026-05-21' 	and *ovl*.is\_extra\_call \= true left join mobile.visits *v*  		on *v*.emp\_code \= *pjp*.salesman\_code  	and *v*.outlet\_code \= *mo*.outlet\_code  	and v.created\_at::date \= '2026-05-21' where *pjp*.salesman\_id \=482 and *pjp*.approval\_status in('Approved', 'Need Review') and dh."date"::date \='2026-05-21' and *dh*.is\_extra\_call \= true  |
| :---- |

### ---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| emp\_id | int | 8 | **mst.m\_employee.**emp\_id |
| emp\_code | varchar | 10 | **mst.m\_employee.**emp\_code |
| emp\_name | varchar | 150 | **mst.m\_employee.**emp\_name |
| distributor\_id | int | 8 | **mst.m\_distributor.**distributor\_id |
| area\_id | varchar | 20 | **mst.m\_distributor.**area\_id |
| region\_id | varchar | 150 | **mst.m\_distributor.**region\_id |
| attendance\_id | int | 8 | **mobile.attendance.**attendance\_id (filter attendance\_typ=1 (clock in)  |
| attendance\_longitude  | varchar  | 50 | **mobile.attendance.**attendance\_longitude attendance\_typ=1 (clock in)  |
| attendance\_latitude | varchar  | 50 | **mobile.attendance.**attendance\_latitude attendance\_typ=1 (clock in)  |
| attendance\_at | timestamp  | 6 | **mobile.attendance.**created\_at attendance\_typ=1 (clock in)  |
| clock\_out | int | 8 | **mobile.attendance.**attendance\_id (filter attendance\_typ=2(clock in)  |
| clock\_out\_longitude  | varchar  | 50 | **mobile.attendance.**attendance\_longitude  (filter attendance\_typ=2(clock in)  |
| clock\_out\_latitude | varchar  | 50 | **mobile.attendance.**attendance\_latitude  (filter attendance\_typ=2(clock in)  |
| clock\_out\_at | timestamp  | 6 | **mobile.attendance.**created\_at  (filter attendance\_typ=2(clock in)  |
| **pjp\_data** | **Array** |  |  |
| pjp\_code | int | 8 | **pjp\_principles.permanent\_journey\_plans.**pjp\_code |
| approval\_status | varchar | 32 | **pjp\_principles.permanent\_journey\_plans.**approval\_status |
| route\_data | **Array** |  |  |
| route\_code | int | 8 | **pjp\_principles.routes.**route\_code |
| route\_name | varchar | 125 | **pjp\_principles.routes.**route\_name |
| destination\_data | **Array** |  |  |
| destination\_id | int | 8 | **pjp\_principles.destinations.**destination\_id |
| destination\_code | varchar | 125 | **pjp\_principles.destinations.**destination\_code |
| destination\_type | varchar | 125 | **pjp\_principles.destinations.**destination\_type |
| destination\_name | varchar | 125 | **pjp\_principles.destinations.**destination\_name |
| longitude | varchar | 125 | **pjp\_principles.destinations.**longitude |
| latitude | varchar | 125 | **pjp\_principles.destinations.**latitude |
| arrive\_at | varchar | 125 | **pjp\_principles.outlet\_visit\_list ovl.**arrive\_at |
| leave\_at | varchar | 125 | **pjp\_principles.outlet\_visit\_list ovl.**leave\_at |
| arrive\_longitude | varchar | 125 | **mobile.visits.**longitude |
| arrive\_latitude | varchar | 125 | **mobile.visits.**latitude |
| start | int | 8 | **pjp\_principles.outlet\_visit\_list ovl.**start |
| finish | int | 8 | **pjp\_principles.outlet\_visit\_list ovl.**finish |
| skip\_at | int | 8 | **pjp\_principles.outlet\_visit\_list ovl.**skip\_at |

**Example Response :** 

| Case : sukses dan terdapat data {  "message": "Success",  "data": \[    {      "emp\_id": 1001,      "emp\_code": "EMP001",      "emp\_name": "Budi Santoso",      "distributor\_id": 23,      "area\_id": 22,      "region\_id": 24,      "pjp\_data": \[        {          "pjp\_code": 1,          "approval\_status": "Approved",          "route\_data": \[            {              "route\_code": "R001",              "route\_name": "Route 001",              "destination\_data": \[                {                  "destination\_id": 1,                  "destination\_code": "D12121",                  "destination\_type" : "Outlet",                  "destination\_name": "destination\_name",                  "longitude": \-6.234567,                  "latitude": 106.123456,                  "arrive\_at": null,                  "start": null,                  "finish": null,                  "skip\_at": null                },                 {                  "destination\_id": 1,                  "destination\_code": "D12121",                  "destination\_type" : "Distributor",                  "destination\_name": "destination\_name",                  "longitude": \-6.234567,                  "latitude": 106.123456,                  "arrive\_at": null,                  "start": null,                  "finish": null,                  "skip\_at": null                }              \]            }          \]        }      \]    }  \],  "paging": {    "total\_record": 84,    "page\_current": 1,    "page\_limit": 10,    "page\_total": 9  },  "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| **Case : empty state / tidak terdapat data berdasarkan pencarian** {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| **Case : error** {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

6. # **Location Monitoring (Distributor)**

 Create new endpoint   
 Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}scylla-pjp/v1/live-monitoring-distributor

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Param** 
---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| date | epoch  | Yes | pjp.route\_pop\_permanent.date |
| emp\_id | Array\<int\> | No | pjp.permanent\_journey\_plans.salesman\_id |
| status | Array\<varchar\> | Yes | untuk saat ini FE kirim ‘Approved’ & ‘Need Review’ pjp.permanent\_journey\_plans.approval\_status |

**Example Request :** 

| curl \--location \-g \\   '{{url}}scylla-pjp/v1/live-monitoring-distributor?date=1738195200\&emp\_id\[\]=358\&status\[\]=Approved' \\   \--header 'Accept: application/json' \\   \--header 'Authorization: Bearer {{token}}' |
| :---- |

### **Response :**  {#response-:}

Query PJP : 

| select 	*md*.distributor\_id , 	*pjp*.approval\_status , 	*md*.area\_id , 	*md*.region\_id , 	*pjp*.id , 	*pjp*.approval\_status , 	*pjp*.salesman\_id , 	*pjp*.salesman\_code , 	*pjp*.salesman\_name , 	*r*.route\_code , 	*r*.route\_name , 	*mo*.outlet\_id , 	*mo*.outlet\_name , 	*mo*.outlet\_code , 	*mo*.address1 , 	*roh*.longitude , 	*roh*.latitude ,  	*v*.longitude as *arrive\_longitude*, 	*v*.latitude as *arrive\_latitude*, 	*ovl*.on\_hold,  	*ovl*.id  from pjp.route\_outlet\_history *roh* join pjp.permanent\_journey\_plans *pjp* 	on *roh*.pjp\_id \= *pjp*.id join mst.m\_outlet *mo* 	on *roh*.outlet\_id \= *mo*.outlet\_id join mst.m\_salesman *ms* 	on *ms*.emp\_id \= *pjp*.salesman\_id join smc.m\_customer *mc* 	on *ms*.cust\_id \= *mc*.cust\_id join mst.m\_distributor *md*  	on *md*.distributor\_id \=*mc*.distributor\_id join pjp.routes *r* 	on *roh*.route\_code  \= *r*.route\_code left join mobile.visits *v* 	on *v*.emp\_code \= *pjp*.salesman\_code 	AND *v*.outlet\_code \= *mo*.outlet\_code 	and v.created\_at::date \= '2026-03-06' left join pjp.outlet\_visit\_list *ovl* 	on *ovl*.outlet\_id \= *roh*.outlet\_id 	and *ovl*.pjp\_id \= *roh*.pjp\_id 	and ovl."date"::date \= '2026-03-06' 	and *ovl*.is\_extra\_call \= false where *pjp*.salesman\_id IN ( '360')  and *roh*."date" \='2026-03-06' and *roh*.is\_extra\_call \= false and *pjp*.approval\_status in('Approved', 'Need Review')  |
| :---- |
| **PJP Extra Call   select 	*md*.distributor\_id , 	*pjp*.approval\_status , 	*md*.area\_id , 	*md*.region\_id , 	*pjp*.id , 	*pjp*.approval\_status , 	*pjp*.salesman\_id , 	*pjp*.salesman\_code , 	*pjp*.salesman\_name , 	*r*.route\_code , 	*r*.route\_name , 	*mo*.outlet\_id , 	*mo*.outlet\_name , 	*mo*.outlet\_code , 	*mo*.address1 , 	*roh*.longitude , 	*roh*.latitude ,  	*v*.longitude as *arrive\_longitude*, 	*v*.latitude as *arrive\_latitude*, 	*ovl*.on\_hold,  	*ovl*.id  from pjp.route\_outlet\_history *roh* join pjp.permanent\_journey\_plans *pjp* 	on *roh*.pjp\_id \= *pjp*.id join mst.m\_outlet *mo* 	on *roh*.outlet\_id \= *mo*.outlet\_id join mst.m\_salesman *ms* 	on *ms*.emp\_id \= *pjp*.salesman\_id join smc.m\_customer *mc* 	on *ms*.cust\_id \= *mc*.cust\_id join mst.m\_distributor *md*  	on *md*.distributor\_id \=*mc*.distributor\_id join pjp.routes *r* 	on *roh*.route\_code  \= *r*.route\_code left join mobile.visits *v* 	on *v*.emp\_code \= *pjp*.salesman\_code 	AND *v*.outlet\_code \= *mo*.outlet\_code 	and v.created\_at::date \= '2026-03-12' left join pjp.outlet\_visit\_list *ovl* 	on *ovl*.outlet\_id \= *roh*.outlet\_id 	and *ovl*.pjp\_id \= *roh*.pjp\_id 	and ovl."date"::date \= '2026-03-12' 	and *ovl*.is\_extra\_call \= true where *pjp*.salesman\_id IN ( '360')  and *roh*."date" \='2026-03-12' and *roh*.is\_extra\_call \= true and *pjp*.approval\_status in('Approved', 'Need Review')** |

### 

### ---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| emp\_id | int | 8 | **mst.m\_employee.**emp\_id |
| emp\_code | varchar | 10 | **mst.m\_employee.**emp\_code |
| emp\_name | varchar | 150 | **mst.m\_employee.**emp\_name |
| distributor\_id | int | 8 | **mst.m\_distributor.**distributor\_id |
| area\_id | varchar | 20 | **mst.m\_distributor.**area\_id |
| region\_id | varchar | 150 | **mst.m\_distributor.**region\_id |
| attendance\_id | int | 8 | **mobile.attendance.**attendance\_id (filter attendance\_typ=1 (clock in)  |
| attendance\_longitude  | varchar  | 50 | **mobile.attendance.**attendance\_longitude attendance\_typ=1 (clock in)  |
| attendance\_latitude | varchar  | 50 | **mobile.attendance.**attendance\_latitude attendance\_typ=1 (clock in)  |
| attendance\_at | timestamp  | 6 | **mobile.attendance.**created\_at attendance\_typ=1 (clock in)  |
| clock\_out | int | 8 | **mobile.attendance.**attendance\_id (filter attendance\_typ=2(clock in)  |
| clock\_out\_longitude  | varchar  | 50 | **mobile.attendance.**attendance\_longitude  (filter attendance\_typ=2(clock in)  |
| clock\_out\_latitude | varchar  | 50 | **mobile.attendance.**attendance\_latitude  (filter attendance\_typ=2(clock in)  |
| clock\_out\_at | timestamp  | 6 | **mobile.attendance.**created\_at  (filter attendance\_typ=2(clock in)  |
| **pjp\_data** | **Array** |  |  |
| pjp\_id | int | 8 | **pjp.permanent\_journey\_plans.**pjp\_id |
| approval\_status | varchar | 32 | **pjp.permanent\_journey\_plans.**approval\_status |
| route\_data | **Array** |  |  |
| route\_code | int | 8 | **pjp.routes.**route\_code |
| route\_name | varchar | 125 | **pjp.routes.**route\_name |
| destination\_data | **Array** |  |  |
| destination\_id | int | 8 | **pjp.route\_outlet.**outlet\_id |
| destination\_code | varchar | 125 | **pjp.route\_outlet.**outlet\_code |
| destination\_type | varchar | 125 | **outlet** |
| destination\_name | varchar | 125 | **pjp.route\_outlet.**outlet\_name |
| longitude | varchar | 125 | **pjp.route\_outlet.**longitude |
| latitude | varchar | 125 | **pjp.route\_outlet.**latitude |
| arrive\_at | int | 8 | **pjp.outlet\_visit\_list ovl.**arrive\_at |
| leave\_at | int | 8 | **pjp.outlet\_visit\_list ovl.**leave\_at |
| arrive\_longitude | varchar | 125 | **mobile.visits.**longitude |
| arrive\_latitude | varchar | 125 | **mobile.visits.**latitude |
| start | int | 8 | **pjp.outlet\_visit\_list ovl.**start |
| finish | int | 8 | **pjp.outlet\_visit\_list ovl.**finish |
| skip\_at | int | 8 | **pjp.outlet\_visit\_list ovl.**skip\_at |

**Example Response :** 

| Case : sukses dan terdapat data {  "message": "Success",  "data": \[    {      "emp\_id": 1001,      "emp\_code": "EMP001",      "emp\_name": "Budi Santoso",      "distributor\_id": 23,      "area\_id": 22,      "region\_id": 24,      "attendance\_id": ,      "attendance\_longitude": ,      "attendance\_latitude":,      "attendance\_at":,      "pjp\_data": \[        {          "pjp\_code": 1,          "approval\_status": "Approved",          "route\_data": \[            {              "route\_code": "R001",              "route\_name": "Route 001",              "destination\_data": \[                {                  "destination\_id": 1,                  "destination\_code": "D12121",                  "destination\_type" : "Outlet",                  "destination\_name": "destination\_name",                  "longitude": \-6.234567,                  "latitude": 106.123456,                  "arrive\_at": null,                  "leave\_at": null,                  "arrive\_longitude": null,                  "arrive\_latitude": null,                  "start": null,                  "finish": null,                  "skip\_at": null                },                 {                  "destination\_id": 1,                  "destination\_code": "D12121",                  "destination\_type" : "Distributor",                  "destination\_name": "destination\_name",                  "longitude": \-6.234567,                  "latitude": 106.123456,                  "arrive\_at": null,                  "leave\_at": null,                  "arrive\_longitude": null,                  "arrive\_latitude": null,                  "start": null,                  "finish": null,                  "skip\_at": null                }              \]            }          \],           "extra\_call\_data": \[        {          "pjp\_code": 1,          "approval\_status": "Approved",          "route\_data": \[            {              "route\_code": "R001",              "route\_name": "Route 001",              "destination\_data": \[                {                  "destination\_id": 1,                  "destination\_code": "D12121",                  "destination\_type" : "Outlet",                  "destination\_name": "destination\_name",                  "longitude": \-6.234567,                  "latitude": 106.123456,                  "arrive\_at": null,                  "leave\_at": null,                  "arrive\_longitude": null,                  "arrive\_latitude": null,                  "start": null,                  "finish": null,                  "skip\_at": null                },                 {                  "destination\_id": 1,                  "destination\_code": "D12121",                  "destination\_type" : "Distributor",                  "destination\_name": "destination\_name",                  "longitude": \-6.234567,                  "latitude": 106.123456,                  "arrive\_at": null,                  "leave\_at": null,                  "arrive\_longitude": null,                  "arrive\_latitude": null,                  "start": null,                  "finish": null,                  "skip\_at": null                }              \]            }          \],  "paging": {    "total\_record": 84,    "page\_current": 1,    "page\_limit": 10,    "page\_total": 9  },  "request\_id": "6915a5e8e3f53f84fe73517f"  |
| :---- |
| **Case : empty state / tidak terdapat data berdasarkan pencarian** {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| **Case : error ** {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }   |

attendance :   
clock in  
![][image1]

clock out  
![][image2]

7. # **Location Monitoring \- Detail** 

 Create new endpoint   
  Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}scylla-pjp/v1/monitoring\_locations/details

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Param** 
---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| emp\_id | int(8) | Yes | *pjp*.salesman\_id |
| distributor\_id | int(8) | No  | If NULL → user principal  If Not NULL → user distributor ***tidak untuk filter data*** |
| date | date | Yes | pjp\_principles.route\_pop\_permanent.date |

**Example Request :** 

| curl \--location \-g \\   '{{url}}/v1/monitoring\_locations/details?emp\_id=339\&distributor\_id=68\&date=2025-12-15' \\   \--header 'Accept: application/json' \\   \--header 'Authorization: Bearer {{token}}'  |
| :---- |

### 

### **Query Visit Information**

![][image3]  
**query principles :** 

| select  	*pjp*.id, 	*ms*.sales\_name,  	SUM( 		case 			when *dh*.is\_extra\_call \= false  			then 1 else 0 		end 		) as *Planned*, 	SUM(         CASE              WHEN *ovl*.arrive\_at IS NOT NULL               AND *ovl*.leave\_at IS NULL              THEN 1 ELSE 0          END     ) AS *on\_going*,     SUM(         CASE              WHEN *dh*.is\_extra\_call \= true              THEN 1 ELSE 0          END     ) AS *extra\_call*,      SUM(         CASE              WHEN *ovl*.arrive\_at IS NOT NULL               AND *ovl*.leave\_at IS NOT NULL              THEN 1 ELSE 0          END     ) AS *visited*,      SUM(         CASE              WHEN *ovl*.skip\_at IS NOT NULL              THEN 1 ELSE 0          END     ) AS *skipped* from pjp\_principles.destinations\_history *dh*  join pjp\_principles.permanent\_journey\_plans *pjp*  	on *pjp*.id \= *dh*.pjp\_id  and *dh*.cust\_id \= *pjp*.cust\_id  join mst.m\_outlet *mo* on *mo*.outlet\_id \= *dh*.destination\_id  left join pjp\_principles.routes *r*  	on *r*.route\_code \= *dh*.route\_code  left join mst.m\_distributor *md*  	on *md*.distributor\_id \= *dh*.destination\_id  left join pjp\_principles.outlet\_visit\_list *ovl*  	on *ovl*.outlet\_id \= *dh*.destination\_id  	and  *ovl*.pjp\_id \= *dh*.pjp\_id  	and ovl."date"::date \= '2026-05-22' left join mobile.visits *v*  		on *v*.emp\_code \= *pjp*.salesman\_code  	and *v*.outlet\_code \= *mo*.outlet\_code  	and v.created\_at::date \= '2026-05-22' left join mst.m\_salesman *ms*  	on *ms*.emp\_id \= *pjp*.salesman\_id  where *pjp*.salesman\_id \=482 and *pjp*.approval\_status in('Approved', 'Need Review') and dh."date"::date \='2026-05-22'  group by *pjp*.id, *ms*.sales\_name  |
| :---- |

**query distributor :**

| select  	*me*.emp\_id ,  	*me*.emp\_code ,  	*me*.emp\_name ,  	count(*ro*.outlet\_id ) *PLAN*, 	COUNT(CASE WHEN *ovl*."start" IS NOT NULL THEN 1 END) AS *on\_going*, 	COUNT(CASE WHEN *ovl*."finish" IS NOT NULL THEN 1 END) AS *visited*,     COUNT(CASE WHEN *ovl*.skip\_at IS NOT NULL THEN 1 END) AS *total\_skip* from pjp.permanent\_journey\_plans *pjp*  join pjp.route\_pop\_permanent *rpp*  	on *rpp*.pjp\_id \= *pjp*.id  join pjp.routes *r*  	on *r*.route\_code \=*rpp*.route\_code  join pjp.route\_outlet *ro*  	on *ro*.route\_code \= *r*.route\_code  join mst.m\_salesman *ms*  	on *pjp*.salesman\_id \= *ms*.emp\_id  join mst.m\_employee *me*  	on *me*.emp\_id \= *ms*.emp\_id  join smc.m\_customer *mc*  	on *mc*.cust\_id \= *ms*.cust\_id  join mst.m\_distributor *md*  	on *md*.distributor\_id \= *mc*.distributor\_id  left join pjp.outlet\_visit\_list *ovl*  	on *ovl*.route\_code  \= *ro*.route\_code  where  *pjp*.salesman\_id  \= 358 and *rpp*."date" \='2026-01-30' and *pjp*.approval\_status \='Approved' group  by *me*.emp\_id, *me*.emp\_code, *me*.emp\_name  |
| :---- |

 

### **Query Sales :** 

![][image4]

### **![][image5]**

| SELECT      o.outlet\_name,     o.outlet\_code,     COUNT(ord.ro\_no) AS *order\_count*,     SUM(ord.sub\_total) AS *gross\_sales*,     SUM(ord.disc\_value) AS *total\_discount*,     SUM(ord.vat\_value) AS *total\_vat*,     SUM(ord.total) AS *net\_sales* FROM      sls.order ord INNER JOIN      mst.m\_outlet o ON ord.outlet\_id \= o.outlet\_id      AND ord.cust\_id \= o.cust\_id WHERE      ord.ro\_date \= '2026-01-30'     AND ord.salesman\_id=228     AND ord.is\_del \= false GROUP BY      o.outlet\_name,     o.outlet\_code ORDER BY      net\_sales DESC;  |
| :---- |

### **Query Return :** 

![][image6]  
![][image7]

| SELECT      *o*.outlet\_name,     *o*.outlet\_code,     SUM(*r*.total) AS *total\_return* FROM      sls.return *r* INNER JOIN      mst.m\_outlet *o* ON *r*.outlet\_id \= *o*.outlet\_id      AND *r*.cust\_id \= *o*.cust\_id WHERE      *r*.emp\_id \= '346'     AND *r*.return\_date \= '2026-01-28'     AND *r*.is\_del \= false GROUP BY      *o*.outlet\_name,     *o*.outlet\_code ORDER BY      *total\_return* DESC;  |
| :---- |

### **Collection** 

| Enhance 28 May 2026  Menampilkan tagihan yang sudah di bayarkan pada tanggal yang sudah di pilih  with *deposit\_data* as ( select  	*d*.emp\_id , 	*o*.outlet\_id,  	*d*.deposit\_no ,  	*d*.deposit\_date ,  	*d*.deposit\_status ,  	*me*.emp\_name ,  	*d*.collection\_no,  	*dd*.invoice\_no,  	*o*.total,  	*dp*.pay\_type , 	*dp*.payment\_amount  from acf.deposit *d* join mst.m\_employee *me*      on *me*.emp\_id \= *d*.emp\_id join acf.deposit\_detail *dd*  	on *dd*.deposit\_no \= *d*.deposit\_no and *dd*.cust\_id \= *d*.cust\_id  join sls."order" *o*  	on *o*.invoice\_no \= *dd*.invoice\_no and *o*.cust\_id \= *dd*.cust\_id  left join acf.deposit\_payment *dp*  	on *dp*.deposit\_no \= *d*.deposit\_no and *d*.cust\_id \= *dp*.cust\_id  where *d*.deposit\_date \= '2026-05-28' and *d*.collection\_no  is not null  and *d*.emp\_id  \= 421 ) select *mo*.outlet\_code ,  *mo*.outlet\_name , sum(*dp*.payment\_amount) FROM *deposit\_data* *dp* join mst.m\_outlet *mo*  on *mo*.outlet\_id \= *dp*.outlet\_id  group by *mo*.outlet\_code ,  *mo*.outlet\_name   |
| :---- |

### **Expense :** 

![][image8]

| select *e*.created\_by , *e*.date, *et*.expense\_type\_name , *e*.note , *e*.amount  from acf.expense *e* join acf.expense\_type *et*  	on *e*.expense\_type\_id \=*et*.expense\_type\_id where *e*.collector\_id \=28 and *e*.date \= '2026-01-23'  |
| :---- |

**Shipment** : 

### **![][image9]**

| select   *si*.shipment\_no ,   *si*.status ,   *o*.outlet\_name,   *o*.outlet\_code,   SUM(*si*.total\_netto) from tms.shipment\_invoices *si* INNER JOIN    mst.m\_outlet *o* on *si*.outlet\_id \=*o*.outlet\_id    AND *si*.cust\_id \= *o*.cust\_id where     *si*.salesman\_id  \= '210'    AND *si*.delivery\_date  \= '2025-02-25' GROUP BY    *o*.outlet\_name,    *o*.outlet\_code,    *si*.status ,    *si*.shipment\_no |
| :---- |

### **Survey** 

| Enhance 29 May 2026  Menampilkan survey yang sudah di submit oleh salesman   select count(*sa*.survey\_answer\_id ) as *submissions*,  *ms*.survey\_title , *mo*.outlet\_code , *mo*.outlet\_name    from mst.survey\_answer *sa*  join mst.m\_survey *ms*  	on *ms*.survey\_id \= *sa*.survey\_id  join mst.m\_outlet *mo*  	on *mo*.outlet\_id  \= *sa*.outlet\_id  where sa.answer\_date::date \= '2026-05-28' and *sa*.emp\_id \= 210 and *sa*.status \='Submitted' group by *ms*.survey\_title , *mo*.outlet\_code , *mo*.outlet\_name  ![][image10]  |
| :---- |

### **Response** 

### ---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| visit\_information | Object |  |  |
| activity\_date | date | \- | dari request body  |
| company\_name | varchar | 150 | **if req distributor\_id \= NULL  sys.m\_user.**user\_fullname **if req distributor\_id \= NOT NULL**  **mst.m\_distributor.**distributor\_name |
| company\_code | varchar | 20 | **if req distributor\_id \= NULL  \-** **if req distributor\_id \= NOT NULL**  **mst.m\_distributor.**distributor\_code |
| level  | varchar | 20 | **if req distributor\_id \= NULL  then “Principal” if req distributor\_id \= NOT NULL  then “Distributor”** |
| emp\_id | int | 8 | **mst.m\_employee.**emp\_id |
| emp\_code | varchar | 10 | **mst.m\_employee.**emp\_code |
| emp\_name | varchar | 150 | **mst.m\_employee.**emp\_name |
| activity\_time | timestamp | 6 | **mobile.attendances.**attendance\_id where type 1 (checkin)  |
| planned | int | 8 | **count destination.destination\_id** |
| on\_going | int | 8 | count outlet\_visit\_list where start is not null  |
| extra\_call |  |  |  |
| visited | int | 8 | count outlet\_visit\_list where finish is not null  |
| skipped | int | 8 | count outlet\_visit\_list where skip\_at  is not null  |
| **sales** | **Array**  |  |  |
| outlet\_id | int | 8 | sls.order.outlet\_id |
| outlet\_code | varchar | 30 | mst.m\_outlet.outlet\_code |
| outlet\_name | varchar | 150 | mst.m\_outlet.outlet\_name |
| sales\_order | numeric | 20,4 | sls.order.total |
| **return** | **Array**  |  |  |
| outlet\_id | int | 8 | sls.return.outlet\_id |
| outlet\_code | varchar | 30 | mst.m\_outlet.outlet\_code |
| outlet\_name | varchar | 150 | mst.m\_outlet.outlet\_name |
| return\_total | numeric | 20,4 | sls.return.total |
| **collection** | **Array**  |  |  |
| outlet\_id | int | 8 | NULL (menunggu dev sby) |
| outlet\_code | varchar | 30 | NULL (menunggu dev sby) |
| outlet\_name | varchar | 150 | NULL (menunggu dev sby) |
| collection\_total | numeric | 20,4 | NULL (menunggu dev sby) |
| **expense** | **Array**  |  |  |
| expense\_type\_id | int | 8 | acf.expense.expense\_type\_id |
| expense\_type | varchar | 100 | acf.expense\_type.expense\_type\_name |
| note | varchar | 100 | acf.expense.note |
| expense\_total | numeric | 20,4 | acf.expense.amount |
| **shipment**  | **Array**  |  |  |
| shipment\_no | varchar | 125 | tms.shipmen\_invoices.shipment\_no |
| **shipment\_data** | **Array** |  |  |
| outlet\_id | int | 8 | tms.shipmen\_invoices.outlet\_id |
| outlet\_name | varchar | 30 | mst.m\_outlet.outlet\_code |
| outlet\_code | varchar | 150 | mst.m\_outlet.outlet\_name |
| total\_netto | numeric | 20,4 | tms.shipmen\_invoices.total\_netto |
| **survey\_data** | **Array** |  |  |
| submission | int | 8 | count survey answer  |
| survey\_title | varchar | 150 | mst.m\_survey.survey\_title |
| outlet\_code | varchar | 150 | mst.m\_outlet.outlet\_name |
| outlet\_name | varchar | 30 | mst.m\_outlet.outlet\_code |

**Example Response :** 

| {  "message": "Success",  "data": \[    {      "visit\_information": {        "activity\_date": "2026-01-12",        "company\_name": "PT Sukses Makmur",        "company\_code": "SM001",        "level": "Principal",        "emp\_id": 12,        "emp\_code": "EMP012",        "emp\_name": "Andi Pratama",        "activity\_time": "2026-01-12 09:15:00",        "planned": 10,        "on\_going": 2,        "extra\_call": 1,        "visited": 6,        "skipped": 1      },      "sales": \[        {          "outlet\_id": 101,          "outlet\_code": "OTL101",          "outlet\_name": "Toko Jaya",          "sales\_order": 2500000        },        {          "outlet\_id": 102,          "outlet\_code": "OTL102",          "outlet\_name": "Toko Makmur",          "sales\_order": 1750000        }      \],      "return": \[        {          "outlet\_id": 101,          "outlet\_code": "OTL101",          "outlet\_name": "Toko Jaya",          "return\_total": 300000        }      \],      "collection": \[\],      "expense": \[        {          "expense\_type\_id": 1,          "expense\_type": "Transport",          "note": "Bensin",          "expense\_total": 100000        },        {          "expense\_type\_id": 2,          "expense\_type": "Makan",          "note": "Lunch",          "expense\_total": 75000        }      \],      "shipment": \[        {          "shipment\_no": "SHP20260112001",          "shipment\_data": \[            {              "outlet\_id": 101,              "outlet\_name": "Toko Jaya",              "outlet\_code": "OTL101",              "total\_netto": 1250000            },            {              "outlet\_id": 102,              "outlet\_name": "Toko Makmur",              "outlet\_code": "OTL102",              "total\_netto": 900000            }          \]        },        {          "shipment\_no": "SHP20260112002",          "shipment\_data": \[            {              "outlet\_id": 101,              "outlet\_name": "Toko Jaya",              "outlet\_code": "OTL101",              "total\_netto": 1250000            },            {              "outlet\_id": 102,              "outlet\_name": "Toko Makmur",              "outlet\_code": "OTL102",              "total\_netto": 900000            }          \]        }      \],  "survey\_data": \[ { "submission": 8, "survey\_title": "Customer Satisfaction Survey", "outlet\_code": "OTL101", "outlet\_name": "Toko Jaya" }, { "submission": 5, "survey\_title": "Product Availability Survey", "outlet\_code": "OTL102", "outlet\_name": "Toko Makmur" } \]    }  \] }  |
| :---- |

# **ISSUE**

# **List Monitoring** 

Method		: GET  
URL		: {{url}}scylla-pjp/v1/live-monitoring-distributor

| curl 'https://best.scyllax.online/scylla-pjp/api/v1/live-monitoring-distributor?date=1773835200\&status%5B%5D=Approved\&status%5B%5D=Need+Review\&emp\_id=360' \\   \-H 'Accept: application/json, text/plain, /' \\   \-H 'Accept-Language: en-US,en;q=0.9,id-ID;q=0.8,id;q=0.7' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsImVtcF9pZCI6MCwiZXhwaXJlcyI6MTc3Mzg5MjYwMSwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4NTc0ODEyMzEyOCIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IkRpc3QgSURFIFNkYSIsInVzZXJfaWQiOjEyLCJ1c2VyX25hbWUiOiJkaXN0QHNkYS5pZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODU3NDgxMjMxMjgifQ.p-o0cVGgOAQZUBGCF0JE17-jKKdUTw1SKidR0QNHc2I' \\   \-H 'Connection: keep-alive' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "Windows"'  |
| :---- |

Payload :   
![][image11]

| date | 1773835200 |
| :---- | :---- |

Response : 

| {    "data": \[        {            "emp\_id": 360,            "emp\_code": "2025120204",            "emp\_name": "Yogie Setya",            "distributor\_id": 67,            "area\_id": 82,            "region\_id": 67,            "attendance\_longitude": 106.6371410311893,            "attendance\_latitude": \-6.1748177746599735,            "attendance\_at": 1773820570,            "pjp\_data": \[                {                    "pjp\_id": 212,                    "approval\_status": "Approved",                    "route\_data": \[                        {                            "route\_code": "4869",                            "route\_name": "Route 3",                            "destination\_data": \[                                {                                    "destination\_id": 1462,                                    "destination\_code": "0019",                                    "destination\_type": "Outlet",                                    "destination\_name": "Babeh TK",                                    "longitude": 106.81473652161236,                                    "latitude": \-6.251901277564605,                                    "arrive\_at": null,                                    "arrive\_longitude": 0,                                    "arrive\_latitude": 0,                                    "start": 1773820773377,                                    "finish": null,                                    "skip\_at": null                                },                                {                                    "destination\_id": 1525,                                    "destination\_code": "09845",                                    "destination\_type": "Outlet",                                    "destination\_name": "Toko Abadi Membumi",                                    "longitude": 106.82795438915491,                                    "latitude": \-6.2523470625577815,                                    "arrive\_at": 1773824277382,                                    "arrive\_longitude": 106.818474,                                    "arrive\_latitude": \-6.360242,                                    "start": 1773820773377,                                    "finish": null,                                    "skip\_at": null                                }                            \]                        }                    \],                    "extra\_call\_data": \[                        {                            "route\_code": "4869",                            "route\_name": "Route 3",                            "destination\_data": \[                                {                                    "destination\_id": 910,                                    "destination\_code": "B000006",                                    "destination\_type": "Distributor",                                    "destination\_name": "BOEDY 6",                                    "longitude": 106.814415,                                    "latitude": \-6.253218,                                    "arrive\_at": 1773821111842,                                    "arrive\_longitude": 106.818514,                                    "arrive\_latitude": \-6.360262,                                    "start": 1773820773377,                                    "finish": null,                                    "skip\_at": null                                },                                {                                    "destination\_id": 912,                                    "destination\_code": "B000008",                                    "destination\_type": "Distributor",                                    "destination\_name": "BOEDY 8",                                    "longitude": 106.817895,                                    "latitude": \-6.251889,                                    "arrive\_at": 1773820784983,                                    "arrive\_longitude": 106.818745,                                    "arrive\_latitude": \-6.251872,                                    "start": 1773820773377,                                    "finish": null,                                    "skip\_at": null                                }                            \]                        }                    \]                }            \]        }    \],    "message": "Success",    "paging": {        "total\_record": 1,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 1    },    "request\_id": "7ed60f96-e82d-49d9-bef5-589806938742" }  |
| :---- |

## Issue : 

![][image12]

Boedy 8 : 

* Arrive : Wednesday, March 18, 2026 at 2:59:44.983 PM [GMT+07:00](https://www.epochconverter.com/timezones?q=1773820784983)  
* Leave **:** Wednesday, March 18, 2026 at 3:00:41.188 PM [GMT+07:00](https://www.epochconverter.com/timezones?q=1773820841188)

Enhance :   
Tambahkan **leave\_at** pada [response](#response-:) 

Query yang di gunakan saat cek data: 

| select 	*md*.distributor\_id , 	*pjp*.approval\_status , 	*md*.area\_id , 	*md*.region\_id , 	*pjp*.id , 	*pjp*.approval\_status , 	*pjp*.salesman\_id , 	*pjp*.salesman\_code , 	*pjp*.salesman\_name , 	*r*.route\_code , 	*r*.route\_name , 	*mo*.outlet\_id , 	*mo*.outlet\_name , 	*mo*.outlet\_code , 	*mo*.address1 , 	*roh*.longitude , 	*roh*.latitude , 	*v*.longitude as *arrive\_longitude*, 	*v*.latitude as *arrive\_latitude*, 	*ovl*.arrive\_at , 	*ovl*.leave\_at , 	*ovl*.on\_hold, 	*ovl*.id from pjp.route\_outlet\_history *roh* join pjp.permanent\_journey\_plans *pjp* 	on *roh*.pjp\_id \= *pjp*.id join mst.m\_outlet *mo* 	on *roh*.outlet\_id \= *mo*.outlet\_id join mst.m\_salesman *ms* 	on *ms*.emp\_id \= *pjp*.salesman\_id join smc.m\_customer *mc* 	on *ms*.cust\_id \= *mc*.cust\_id join mst.m\_distributor *md* 	on *md*.distributor\_id \=*mc*.distributor\_id join pjp.routes *r* 	on *roh*.route\_code  \= *r*.route\_code left join mobile.visits *v* 	on *v*.emp\_code \= *pjp*.salesman\_code 	AND *v*.outlet\_code \= *mo*.outlet\_code 	and v.created\_at::date \= '2026-03-18' left join pjp.outlet\_visit\_list *ovl* 	on *ovl*.outlet\_id \= *roh*.outlet\_id 	and *ovl*.pjp\_id \= *roh*.pjp\_id 	and ovl."date"::date \= '2026-03-18' 	and *ovl*.is\_extra\_call \= true where *pjp*.salesman\_id IN ( '360')  and *roh*."date" \='2026-03-18' and *roh*.is\_extra\_call \= true and *pjp*.approval\_status in('Approved', 'Need Review')  |
| :---- |

# **Detail Monitoring** 

Method		: GET  
URL		: {{url}}scylla-pjp/v1/monitoring\_locations/details

| curl 'https://best.scyllax.online/scylla-pjp/api/v1/monitoring\_locations/details?emp\_id=360\&date=2026-03-18\&distributor\_id=67' \\   \-H 'Accept: application/json, text/plain, /' \\   \-H 'Accept-Language: en-US,en;q=0.9,id-ID;q=0.8,id;q=0.7' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsImVtcF9pZCI6MCwiZXhwaXJlcyI6MTc3Mzg5MjYwMSwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4NTc0ODEyMzEyOCIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IkRpc3QgSURFIFNkYSIsInVzZXJfaWQiOjEyLCJ1c2VyX25hbWUiOiJkaXN0QHNkYS5pZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODU3NDgxMjMxMjgifQ.p-o0cVGgOAQZUBGCF0JE17-jKKdUTw1SKidR0QNHc2I' \\   \-H 'Connection: keep-alive' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "Windows"'  |
| :---- |

Payload : 

| ![][image13] |
| :---- |

Response : 

| {    "data": \[        {            "visit\_information": {                "activity\_date": "2026-03-18",                "company\_name": "Distributor iDetama",                "company\_code": "3434",                "level": "Distributor",                "emp\_id": 360,                "emp\_code": "2025120204",                "emp\_name": "Yogie Setya",                "activity\_time": "2026-03-18T07:56:09.52461Z",                "planned": 5,                "on\_going": 5,                "extra\_call": 0,                "visited": 0,                "skipped": 0            },            "sales": \[                {                    "outlet\_id": 912,                    "outlet\_code": "B000008",                    "outlet\_name": "BOEDY 8",                    "sales\_order": 2220000                },                {                    "outlet\_id": 910,                    "outlet\_code": "B000006",                    "outlet\_name": "BOEDY 6",                    "sales\_order": 444000                }            \],            "return": \[\],            "collection": \[\],            "expense": \[\],            "shipment": \[\]        }    \],    "message": "Success",    "request\_id": "fb61a23b-e13d-49d8-a768-8cfb9a05ddeb" }  |
| :---- |

## Issue  : 

| Seharusnya :  •⁠  ⁠Planned : 2 •⁠  ⁠Extra Call : 2 •⁠  ⁠On going : 1 •⁠  ⁠Visited : 2 •⁠  ⁠Skipped : 0 User : yogie.set@gmail.com  |
| :---- |
| ![][image14] |

Perbaikan Query : 

| SELECT    *pjp*.salesman\_id,    *pjp*.salesman\_code,    *pjp*.salesman\_name,    \-- PLAN (hanya non extra call)    COUNT(DISTINCT CASE        WHEN *ovl\_plan*.outlet\_id IS NOT NULL THEN *ovl\_plan*.outlet\_id    END) AS *planned*,    \-- EXTRA CALL (hanya extra call)    COUNT(DISTINCT CASE        WHEN *ovl\_extra*.outlet\_id IS NOT NULL THEN *ovl\_extra*.outlet\_id    END) AS *extra\_call*,    \-- ON GOING (tanpa lihat is\_extra\_call)    COUNT(CASE        WHEN *ovl\_all*.arrive\_at IS NOT NULL and *ovl\_all*.leave\_at IS NULL THEN 1    END) AS *on\_going*,    \-- VISITED (tanpa lihat is\_extra\_call)    COUNT(CASE        WHEN *ovl\_all*.arrive\_at  IS NOT NULL and *ovl\_all*.leave\_at IS NOT NULL THEN 1    END) AS *visited*,    \-- SKIPPED (tanpa lihat is\_extra\_call)    COUNT(CASE        WHEN *ovl\_all*.skip\_at IS NOT NULL and *ovl\_all*.leave\_at IS NULL THEN 1    END) AS *skipped* FROM pjp.route\_outlet\_history *roh* JOIN pjp.permanent\_journey\_plans *pjp*    ON *roh*.pjp\_id \= *pjp*.id JOIN mst.m\_outlet *mo*    ON *roh*.outlet\_id \= *mo*.outlet\_id JOIN mst.m\_salesman *ms*    ON *ms*.emp\_id \= *pjp*.salesman\_id JOIN smc.m\_customer *mc*    ON *ms*.cust\_id \= *mc*.cust\_id JOIN mst.m\_distributor *md*    ON *md*.distributor\_id \= *mc*.distributor\_id JOIN pjp.routes *r*    ON *roh*.route\_code \= *r*.route\_code \-- PLAN (non extra call) LEFT JOIN pjp.outlet\_visit\_list *ovl\_plan*    ON *ovl\_plan*.outlet\_id \= *roh*.outlet\_id    AND *ovl\_plan*.pjp\_id \= *roh*.pjp\_id    AND ovl\_plan.date::date \= '2026-03-18'    AND *ovl\_plan*.is\_extra\_call \= false \-- EXTRA CALL LEFT JOIN pjp.outlet\_visit\_list *ovl\_extra*    ON *ovl\_extra*.outlet\_id \= *roh*.outlet\_id    AND *ovl\_extra*.pjp\_id \= *roh*.pjp\_id    AND ovl\_extra.date::date \= '2026-03-18'    AND *ovl\_extra*.is\_extra\_call \= true \-- ALL VISIT STATUS (tanpa filter extra call) LEFT JOIN pjp.outlet\_visit\_list *ovl\_all*    ON *ovl\_all*.outlet\_id \= *roh*.outlet\_id    AND *ovl\_all*.pjp\_id \= *roh*.pjp\_id    AND ovl\_all.date::date \= '2026-03-18' WHERE    *pjp*.salesman\_id IN ('360')    AND *roh*.date \= '2026-03-18'    AND *pjp*.approval\_status IN ('Approved', 'Need Review') GROUP BY    *pjp*.salesman\_id,    *pjp*.salesman\_code,    *pjp*.salesman\_name;  ![][image15] |
| :---- |

data test tanggal 25 Maret 2026   
[syahriza@gmail.com](mailto:syahriza@gmail.com) (approve by mas fahmi)   
![][image16]  
