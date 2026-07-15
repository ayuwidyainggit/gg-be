

1. ## **API Survey** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/master/v1/survey?page=1\&limit=10\&sort:created\_date:desc\&status=1

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| q |  | No | template\_title |
| page | Integer | Yes | default : 1 |
| limit  | Integer | Yes | default : 5  |
| sort | String | Yes | default created\_date:desc  |
| response\_frequency | Array String | No | Mandatory, Optional |
| answer\_frequency | Array String | No | Multiple, One Time |
| status | Integer | No | 1: mst.m\_survey\_template.is\_active \= true  0 : mst.m\_survey\_template.is\_active \= false Null : show all status |

### **Example Request** 

---

**Example Request Default**

| curl \--location \-g \\ '{{url}}/master/v1/survey?page=1\&limit=5\&sort=created\_at:desc\&status=1' \\ \--header 'Accept: application/json' \\ \--header 'Authorization: Bearer {{token}}' |
| :---- |

### **Response** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| survey\_id | int | 8 | **mst.m\_survey.**survey\_id |
| created\_at | timestamptz | 6 | **mst.m\_survey.**created\_at |
| answer\_frequency | enum | Multiple, One Time | **mst.m\_survey.**answer\_frequency |
| survey\_title | varchar | 150 | **mst.m\_survey.**survey\_title |
| response\_type | enum | Mandatory, Optional | **mst.m\_survey.**response\_type |
| efective\_date\_start | date |  | **mst.m\_survey.**efective\_date\_start |
| efective\_date\_end | date |  | **mst.m\_survey.**efective\_date\_end |
| status | int | 8 | **mst.m\_survey.**status |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data {  "request\_id": "REQ-202512210001",  "message": "Success",  "data": \[    {      "survey\_id": 12345678,      "created\_at": "2025-12-01T10:15:30Z",      "answer\_frequency": "Multiple",      "survey\_title": "Customer Satisfaction Survey",      "response\_type": "Mandatory",      "efective\_date\_start": "2025-12-01",      "efective\_date\_end": "2025-12-31",      "status": 1    },    {      "survey\_id": 12345679,      "created\_at": "2025-12-05T08:00:00Z",      "answer\_frequency": "One Time",      "survey\_title": "Employee Feedback Survey",      "response\_type": "Optional",      "efective\_date\_start": "2025-12-05",      "efective\_date\_end": "2026-01-05",      "status": 1    }  \],  "paging": {    "total\_record": 25,    "page\_current": 1,    "page\_limit": 10,    "page\_total": 3  } } |
| :---- |
| Case : empty state / tidak terdapat data berdasarkan pencarian {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

2. ## **API Survey Detail**

 Create new endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/master/v1/survey/:survey\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Path Variable** 
---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| survey\_id | Integer | Yes |  |

### **Example Request** 

---

**Example Request Default**

| curl \--location \-g \\ '{{url}}/master/v1/survey/21' \\ \--header 'Accept: application/json' \\ \--header 'Authorization: Bearer {{token}}' |
| :---- |

### **Response** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| survey\_id | int | 8 | **mst.m\_survey.**survey\_template\_id |
| created\_at | timestamptz | 6 | **mst.m\_survey.**created\_at |
| answer\_frequency | enum | Multiple, One Time | **mst.m\_survey.**answer\_frequency |
| survey\_title | varchar | 150 | **mst.m\_survey.**survey\_title |
| response\_type | enum | Mandatory, Optional | **mst.m\_survey.**response\_type |
| efective\_date\_start | date |  | **mst.m\_survey.**efective\_date\_start |
| efective\_date\_end | date |  | **mst.m\_survey.**efective\_date\_end |
| status | int | 8 | **mst.m\_survey.**status |
| **target\_survey** | Object |  |  |
| target\_type | enum (national,) |  | **mst.m\_survey.**target\_type |
| emp\_id | int(8) |  | **mst.m\_survey.**emp\_id |
| sales\_name |  |  | **mst.m\_salesman.**sales\_name |
| **area** | **Array** |  |  |
| area\_id | int | 8 | **mst.m\_survey\_area.**area\_id |
| area\_name |  |  | **mst.m\_area.**area\_name |
| distributor\_id | int | 8 | **mst.m\_survey\_area.**distributor\_id |
| distributor\_code |  |  | **mst.m\_distributor.**distributor\_code |
| distributor\_name |  |  | **mst.m\_distributor.**distributor\_name |
| **template** | Array |  |  |
| survey\_template\_id |  |  | **mst.m\_survey\_detail.**survey\_template\_id |
| template\_code | varchar | 10 | **mst.m\_survey\_template.**template\_code |
| template\_title | varchar | 150 | **mst.m\_survey\_template.**template\_title |
| **question\_template** | Array |  |  |
| question\_template\_id | int | 8 | **mst.m\_question\_template.**question\_template\_id |
| question | varchar | 225 | **mst.m\_question\_template.**question |
| answer\_type | enum (Single, Multiple, Free Text) |  | **mst.m\_question\_template.**answer\_type |
| **options** | Array |  |  |
| q\_option\_template\_id | int | 8 | **mst.m\_q\_option\_template.**q\_option\_template\_id |
| option | varchar | 225 | **mst.m\_q\_option\_template.**option |
| **outlet** | **Array** |  |  |
| survey\_outlet\_id | int | 8 | **mst.m\_survey\_outlet.**survey\_outlet\_id |
| outlet\_id | int | 8 | **mst.m\_survey\_outlet.**outlet\_id |
| outlet\_code |  |  | **mst.m\_outlet.**outlet\_code |
| outlet\_name |  |  | **mst.m\_outlet.**outlet\_name |
| ot\_class\_id | int | 8 | **mst.m\_outlet.**ot\_class\_id |
| ot\_class\_name |  |  | **mst.m\_outlet\_class.**ot\_class\_name |
| ot\_grp\_id | int | 8 | **mst.m\_outlet.**ot\_grp\_id |
| ot\_grp\_name |  |  | **mst.m\_outlet\_group.**ot\_grp\_name |
| ot\_type\_id | int | 8 | **mst.m\_outlet.**ot\_type\_id |
| ot\_type\_name |  |  | **mst.m\_outlet\_type.**ot\_type\_name |
| **salesman** | **Array** |  |  |
| m\_survey\_salesman\_id |  |  | **mst.m\_survey\_salesman.**m\_survey\_salesman\_id |
| sales\_id | int | 8 | **mst.m\_survey\_salesman.**sales\_id |
| sales\_team\_id |  |  | **mst.m\_salesman.**sales\_team\_id |
| sales\_team\_name |  |  | **mst.m\_sales\_team.**sales\_team\_name |
| sales\_name |  |  | **mst.m\_salesman.**sales\_name |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

#### **Case : Level Target \= salesman** 

| {   "request\_id": "REQ-202512210002",   "message": "Success",   "data": \[     {       "survey\_id": 10000001,       "created\_at": "2025-12-01T09:30:00Z",       "answer\_frequency": "Multiple",       "survey\_title": "Survey Kepuasan Pelanggan",       "response\_type": "Mandatory",       "efective\_date\_start": "2025-12-01",       "efective\_date\_end": "2025-12-31",       "status": 1,       "target\_survey": {         "target\_type": "national",         "emp\_id": 87654321,         "sales\_name": "Budi Santoso",         "area": \[           {             "area\_id": 101,             "area\_name": "Jakarta"           },           {             "area\_id": 102,             "area\_name": "Bandung"           }         \],       },       "template": \[         {           "survey\_template\_id": 20001,           "template\_code": "TMP01",           "template\_title": "Template Survey Umum",           "question\_template": \[             {               "question\_template\_id": 30001,               "question": "Bagaimana kualitas pelayanan kami?",               "answer\_type": "Multiple",               "options": \[                 {                   "q\_option\_template\_id": 40001,                   "option": "Sangat Baik"                 },                 {                   "q\_option\_template\_id": 40002,                   "option": "Baik"                 },                 {                   "q\_option\_template\_id": 40003,                   "option": "Cukup"                 },                 {                   "q\_option\_template\_id": 40004,                   "option": "Kurang"                 }               \]             },             {               "question\_template\_id": 30001,               "question": "Bagaimana kualitas pelayanan kami?",               "answer\_type": "Single",               "options": \[                 {                   "q\_option\_template\_id": 40001,                   "option": "Yes / No"                 }               \]             },             {               "question\_template\_id": 30002,               "question": "Saran untuk peningkatan layanan",               "answer\_type": "Free Text",               "options": \[\]             }           \]         }       \],       "outlet": \[\],       "salesman":\[        {          "m\_survey\_salesman\_id": 676,          "sales\_id": 87654321,          "sales\_team\_id": 890,          "sales\_team\_name": "Team Budi",          "sales\_name": "Budi Santoso"        }       \]     }   \],   "paging": {     "total\_record": 12,     "page\_current": 1,     "page\_limit": 10,     "page\_total": 2   } }  |
| :---- |

#### **Case : Level Target \= outlet** 

| {   "request\_id": "REQ-202512210002",   "message": "Success",   "data": \[     {       "survey\_id": 10000001,       "created\_at": "2025-12-01T09:30:00Z",       "answer\_frequency": "Multiple",       "survey\_title": "Survey Kepuasan Pelanggan",       "response\_type": "Mandatory",       "efective\_date\_start": "2025-12-01",       "efective\_date\_end": "2025-12-31",       "status": 1,       "target\_survey": {         "target\_type": "national",         "emp\_id": 87654321,         "sales\_name": "Budi Santoso",         "area": \[           {             "area\_id": 101,             "area\_name": "Jakarta"           },           {             "area\_id": 102,             "area\_name": "Bandung"           }         \]       },       "template": \[         {           "survey\_template\_id": 20001,           "template\_code": "TMP01",           "template\_title": "Template Survey Umum",           "question\_template": \[             {               "question\_template\_id": 30001,               "question": "Bagaimana kualitas pelayanan kami?",               "answer\_type": "Multiple",               "options": \[                 {                   "q\_option\_template\_id": 40001,                   "option": "Sangat Baik"                 },                 {                   "q\_option\_template\_id": 40002,                   "option": "Baik"                 },                 {                   "q\_option\_template\_id": 40003,                   "option": "Cukup"                 },                 {                   "q\_option\_template\_id": 40004,                   "option": "Kurang"                 }               \]             },             {               "question\_template\_id": 30001,               "question": "Bagaimana kualitas pelayanan kami?",               "answer\_type": "Single",               "options": \[                 {                   "q\_option\_template\_id": 40001,                   "option": "Yes / No"                 }               \]             },             {               "question\_template\_id": 30002,               "question": "Saran untuk peningkatan layanan",               "answer\_type": "Free Text",               "options": \[\]             }           \]         }       \],       "outlet": \[         {           "survey\_outlet\_id": 10001,           "outlet\_id": 1001,           "outlet\_code": "S1001",           "outlet\_name": "Outlet A",           "ot\_class\_id": 56,           "ot\_class\_name": "class A",           "ot\_grp\_id": 1,           "ot\_grp\_name" : "Group A",           "ot\_type\_id": 23,           "ot\_type\_name": "Type A"         }       \],       "salesman": \[\]     }   \],   "paging": {     "total\_record": 12,     "page\_current": 1,     "page\_limit": 10,     "page\_total": 2   } }  |
| :---- |

3. ## **API Survey Create**  {#api-survey-create}

 Create new endpoint :  
Content-Type 	: application/json  
Method		: POST  
URL		: {{url}}/master/v1/survey

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Body** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| survey\_title | varchar(150) |  |  |
| efective\_date\_start | date |  |  |
| efective\_date\_end | date |  |  |
| answer\_frequency | enum Multiple, One Time |  |  |
| response\_type | enum Mandatory, Optional |  |  |
| target\_type | enum (Nasional,) |  | **takeout** |
| area\_id | Array(Int) |  |  |
| distributor\_id | Array(Int) |  | **add body** |
| outlet\_id | Array(Int) |  |  |
| survey\_template\_id | Array(Int) |  | before  : Array(Int) **after : Int(8)** |
| emp\_id | Int(8) |  | before : Int(8) **after : Array(Int)** |

### **Example Request** 

---

**Example Request Default**

| curl \--location \\ \--request POST '{{url}}/master/v1/survey' \\ \--header 'Accept: application/json' \\ \--header 'Content-Type: application/json' \\ \--header 'Authorization: Bearer {{token}}' \\ \--data '{  "survey\_title": "Survey Kepuasan Pelanggan Nasional",  "efective\_date\_start": "2025-12-01",  "efective\_date\_end": "2025-12-31",  "answer\_frequency": "Multiple",  "response\_type": "Mandatory",  "target\_type": "Nasional",  "area\_id": \[101, 102, 103\],  "outlet\_id": \[1001, 1002, 1003\],  "survey\_template\_id": \[20001, 20002\],  "emp\_id": \[20001, 20002\], }'  |
| :---- |

### **Response** 

---

| Case : sukses dan terdapat data {  "message": "Survey has been successfully created",  "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| case : Survey Title tidak unik per periode aktif. {  "message": "Survey Title already exists for the active period",  "request\_id": "6915a5e8e3f53f84fe73517f" } |
| Case : error  {  "message": "Failed to create survey data",  "request\_id": "6915a5e8e3f53f84fe73517f" }  |

Impact ke database : 

* #### mst.m\_survey

* mst.m\_survey\_area  
* mst.m\_survey\_detail   
* mst.m\_survey\_outlet 

4. ## **API Survey Edit**  {#api-survey-edit}

 Create new endpoint :  
Content-Type 	: application/json  
Method		: PUT  
URL		: {{url}}/master/v1/survey/:survey\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Body** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| survey\_title | varchar(150) |  |  |
| efective\_date\_start | date |  |  |
| efective\_date\_end | date |  |  |
| answer\_frequency | enum Multiple, One Time |  |  |
| response\_type | enum Mandatory, Optional |  |  |
| target\_type | enum (Nasional,) |  | **takeout** |
| area\_id | Array(Int) |  |  |
| distributor\_id | Array(Int) |  | **add response**  |
| outlet\_id | Array(Int) |  |  |
| survey\_template\_id | Array(Int) |  | before  : Array(Int) **after : Int(8)** |
| emp\_id | Int(8) |  | before : Int(8) **after : Array(Int)** |

### **Example Request** 

---

**Example Request Default**

| curl \--location \\ \--request POST '{{url}}/master/v1/survey/:survey\_id' \\ \--header 'Accept: application/json' \\ \--header 'Content-Type: application/json' \\ \--header 'Authorization: Bearer {{token}}' \\ \--data '{  "survey\_title": "Survey Kepuasan Pelanggan Nasional",  "efective\_date\_start": "2025-12-01",  "efective\_date\_end": "2025-12-31",  "answer\_frequency": "Multiple",  "response\_type": "Mandatory",  "target\_type": "Nasional",  "area\_id": \[101, 102, 103\],  "outlet\_id": \[1001, 1002, 1003\],  "survey\_template\_id": \[20001, 20002\],  "emp\_id": 87654321 }'  |
| :---- |

### **Response :**

| Case : sukses dan terdapat data {  "message": "Survey has been successfully updated",  "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| Case : error  {  "message": "Failed to update survey data",  "request\_id": "6915a5e8e3f53f84fe73517f" }  |

5. ## **API Survey Deactive** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: PATCH  
URL		: {{url}}/master/v1/survey/:survey\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Body** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| is\_active | bool |  | false |

### **Example Request** 

---

**Example Request Default**

| curl \--location \\ \--request POST '{{url}}/master/v1/survey/:survey\_id' \\ \--header 'Accept: application/json' \\ \--header 'Content-Type: application/json' \\ \--header 'Authorization: Bearer {{token}}' \\ \--data '{  "is\_active": false }'  |
| :---- |

### **Response** 

### **Response :**

| Case : sukses dan terdapat data {  "message": "Survey successfully deactivated",  "request\_id": "6915a5e8e3f53f84fe73517f" } |
| :---- |
| Case : error  {  "message": "Survey not found",  "request\_id": "6915a5e8e3f53f84fe73517f" }  |

6. ## **API  Outlet LIST** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: 103.28.219.73/v1/outlets?page=1\&limit=10\&sort=outlet\_id:desc\&is\_active=1\&identity\_type=National ID\&identity\_no=312123456789098765  
[Link doc eksisting](https://docs.google.com/document/d/1cfldOtEZ0G41VFAFTEYXQh0Yd6J2CtaP/edit#heading=h.cv02mhec31t8)

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| page | bool |  | false |
| limit |  |  |  |
| sort |  |  |  |
| q |  |  |  |
| is\_active |  |  | tidak perlu |
| outlet\_status |  |  | 1,5,6,7, |
| outlet\_id |  |  |  |
| verification\_status |  |  | 1 |
| ot\_grp\_id | **integer** | No |  |
| ot\_type\_id | **integer** | No |  |
| identity\_type |  |  |  |
|  identity\_no |  |  |  |
| distributor\_id |  |  | contoh : 0, 67, 680 \= principal 67, 68 \= distributor |
| **ot\_class\_id** | **integer** | **No**  | **Filter berdasarkan mst.m\_outlet.ot\_class\_id** |

### **Contoh Case:** 

### ---

User login : [princessa@gmail.com](mailto:princessa@gmail.com) | Admin123

Filter : 

outlet berdasarkan PT Besi Makmur

ot\_type \= Type S , type A 

class name \= Class Outdoor 

group \= group R

![][image1]

#### **curl :** 

| curl 'https://best.scyllax.online/master/v1/outlets?page=1\&limit=10\&sort=outlet\_id:desc\&is\_active=1\&ot\_type\_id=99,97\&ot\_grp\_id=111\&ot\_class\_id=91\&q=\&page=1\&limit=70' \\   \-H 'sec-ch-ua-platform: "macOS"' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyIiwiZGlzdF9wcmljZV9ncnBfaWQiOjAsImRpc3RyaWJ1dG9yX2lkIjowLCJlbWFpbCI6InByaW5jZXNzYUBnbWFpbC5jb20iLCJlbXBfaWQiOjM4MCwiZXhwaXJlcyI6MTc3NjE2MDM4NiwiaXNfYWRtaW4iOmZhbHNlLCJsYW5nX2lkIjoiaWQiLCJtb2JpbGVfbm8iOiIwODEzMzMzMzMzMzMiLCJwYXJlbnRfY3VzdF9pZCI6IkMyNjAwMiIsInVzZXJfZnVsbG5hbWUiOiJQcmluY2Vzc2EgQWhzYW5pIFRhcXdpbSIsInVzZXJfaWQiOjE0MCwidXNlcl9uYW1lIjoiUHJpbmNlc3NhIEFoc2FuaSBUYXF3aW0iLCJ3aGF0c2FwcCI6IjA4MTMzMzMzMzMzMyJ9.O6RyO5a6idD-AHxPaoGff1\_z1y9l5K0yAwgYllUGpxQ' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'sec-ch-ua: "Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"' \\   \-H 'sec-ch-ua-mobile: ?0'  |
| :---- |

![][image2]

#### **response eksisting (seharusnya ada 10 data dengan filter di atas)** 

![][image3]

#### **query BE :** 

---

| select *mo*.outlet\_code , *mo*.outlet\_name, *mot*.ot\_type\_name , *moc*.ot\_class\_name , *mog*.ot\_grp\_name  from mst.m\_outlet *mo*  join mst.m\_outlet\_type *mot*  	on *mot*.ot\_type\_id \= *mo*.ot\_type\_id  join mst.m\_outlet\_class *moc*  	on *moc*.ot\_class\_id \=*mo*.ot\_class\_id  join mst.m\_outlet\_group *mog*  	on *mog*.ot\_grp\_id \=*mo*.ot\_grp\_id  join smc.m\_customer *mc*  	on *mc*.cust\_id \= *mo*.cust\_id  where *mc*.distributor\_id \=102 and *mo*.ot\_type\_id in(99, 97\) and *mo*.ot\_grp\_id in(111) and *mo*.ot\_class\_id in (91)   |
| :---- |

response yang seharusnya tampil :   
![][image4]

#### **Enhance FE  :**

* FE belum kirim distributor\_id   
* FE belum kirim ot\_class\_id 

7. ## **API Salesman Lookup** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/v1/salesman?q=\&page=1\&limit=10\&sort=sales\_name:asc\&sales\_team\_id=24,1

[Doc eksisting](https://docs.google.com/document/d/1cfldOtEZ0G41VFAFTEYXQh0Yd6J2CtaP/edit#heading=h.a4bkhddc4o8z) 

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Body** 

---

Enhance : add param berikut 

distributor\_id → Optional (int\[\]) 

**Data Test :**   
**Email :** adminbm@gmail.com 


Cust\_id : C260020001

Distributor : PT. Besi Makmur  | distributor\_id \= 102 

Sales Team yang dipilih: 

  \- GT  | sales\_team\_id \= 66

  \- MIX | sales\_team\_id= 65

**Lookup Salesman yang tampil berdasarkan filter distributor dan sales team :** 

**•⁠  ⁠Piere Njangka**

**•⁠  ⁠Jaka**

**Query :** 

| user distributor :  select *mc*.cust\_id , *mc*.distributor\_id ,  *ms*.\* from mst.m\_salesman *ms* join smc.m\_customer *mc* 	on *mc*.cust\_id \= *ms*.cust\_id where *ms*.sales\_team\_id in(66,65) and *ms*.cust\_id \= 'C260020001' and *ms*.is\_active \= true and *ms*.is\_del \= false and *mc*.distributor\_id \=102  |
| :---- |
| **Ketika user principal  filter salesman berdasarkan User Principal :  select *mc*.cust\_id , *mc*.distributor\_id ,  *ms*.\* from mst.m\_salesman *ms* join smc.m\_customer *mc* 	on *mc*.cust\_id \= *ms*.cust\_id where *ms*.sales\_team\_id in(78,77) and *ms*.cust\_id like 'C26002%' and *ms*.is\_active \= true and *ms*.is\_del \= false  Ketika user principal  filter salesman berdasarkan User Principal \+ distributor :  select *mc*.cust\_id , *mc*.distributor\_id ,  *ms*.\* from mst.m\_salesman *ms* join smc.m\_customer *mc* 	on *mc*.cust\_id \= *ms*.cust\_id where *ms*.sales\_team\_id in(78,77) and *ms*.cust\_id like 'C26002%' and *ms*.is\_active \= true and *ms*.is\_del \= false and *mc*.distributor\_id \=102**  |

8. ## **API Salesman Team Lookup** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/master/v1/sales-teams?mode=lookup\&page=1\&limit=70\&distributor\_id=1\&distributor\_id=45\&q=\&page=1\&limit=70

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| q |  |  |  |
| mode |  | No |  |
| page | Integer | Yes |  |
| limit  | Integer | Yes |  |
| distributor\_id | int\[\] | Yes | enhance tambahan : tambahkan distributor\_id \= 0 → artinya mencari sales\_team berdasarkan user milik principal  secara business rules, sales\_team dapat dipilih berdasarkan multiple dropdown dari business\_unit, dimana resp business\_unit adalah user principal itu sendiri dan distributor dibawahnya.  contoh : 0, 67, 680 \= principal 67, 68 \= distributor |

### **Issue :** 

* User Login 		:  [princ@idetama.id](mailto:princ@idetama.id)  
* Business Unit yang di pilih :   
  * Admin Principal 1  ***(ini user principal itu sendiri)***  
  * Distributor iDetama (distributor\_id \= 67 )

  ![][image5]

* CURL : 

| curl 'https://best.scyllax.online/master/v1/sales-teams?mode=lookup\&page=1\&limit=70\&distributor\_id=1\&distributor\_id=67\&q=\&page=1\&limit=70' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxIiwiZGlzdF9wcmljZV9ncnBfaWQiOjAsImRpc3RyaWJ1dG9yX2lkIjowLCJlbWFpbCI6InByaW5jQGlkZXRhbWEuaWQiLCJlbXBfaWQiOjI3OCwiZXhwaXJlcyI6MTc3NjEzNTU1MSwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4MTEzMjMyMzMyIiwicGFyZW50X2N1c3RfaWQiOiJDMjIwMDEiLCJ1c2VyX2Z1bGxuYW1lIjoiQWRtaW4gUHJpbmNpcGFsIDEiLCJ1c2VyX2lkIjoxLCJ1c2VyX25hbWUiOiJwcmluY0BpZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODExMzIzMjMzMiJ9.j7sy-VIHtURjIvINFis-ljcw0SeicRUYOPRO7NEjNi8' \\   \-H 'Connection: keep-alive' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"'  |
| :---- |

* **Payload FE :** 

| mode=lookup& page=1& limit=70& distributor\_id=1& distributor\_id=67& q=& page=1 \&limit=70 |
| :---- |

* **Response :** 

| {     "message": "",     "data": \[         {             "sales\_team\_id": 60,             "sales\_team\_code": "900",             "sales\_team\_name": "Team HO"         },         {             "sales\_team\_id": 50,             "sales\_team\_code": "200",             "sales\_team\_name": "CANVAS"         },         {             "sales\_team\_id": 49,             "sales\_team\_code": "100",             "sales\_team\_name": "TAKING ORDER"         },         {             "sales\_team\_id": 48,             "sales\_team\_code": "000",             "sales\_team\_name": "SHOP SALES"         }     \],     "paging": {         "total\_record": 4,         "page\_current": 1,         "page\_limit": 70,         "page\_total": 1     },     "request\_id": "69dc5f9fc7fd667f49595760" }  |
| :---- |

* Fixing FE :   
  untuk case business unit **Admin Principal 1 , dan Distributor iDetama** distributor seharusnya hanya di kirim 67 dan 0   
* Fixing BE : BE seharusnya filter berdasarkan user principal dan distributor idetama   
* Query : 

| SELECT      *mst*.cust\_id,      *md*.distributor\_id,      *md*.distributor\_name,      *mst*.sales\_team\_id,      *mst*.sales\_team\_code,      *mst*.sales\_team\_name FROM mst.m\_sales\_team *mst*  JOIN smc.m\_customer *mc*      ON *mst*.cust\_id \= *mc*.cust\_id  LEFT JOIN mst.m\_distributor *md*      ON *md*.distributor\_id \= *mc*.distributor\_id  WHERE *mst*.is\_active \= true    AND *mst*.is\_del \= false   AND (       *md*.distributor\_id \= :JWT\_DISTRIBUTOR\_ID  \-- Parameter ID Distributor (Contoh: 67\)       OR mst.cust\_id \= :PARENT\_CUST\_ID         \-- Parameter Cust ID Principal (Contoh: 'C22001')   );  |
| :---- |

# 

# **Issue** 

1. ## **Issue SX 910**

 [https://scyllax-pratesis.atlassian.net/browse/SX-910](https://scyllax-pratesis.atlassian.net/browse/SX-910) 

mas yogi tolong survey\_title dibuat unik   
[**API Create**](#api-survey-create)  
  Method		: POST  
URL			: {{url}}/master/v1/survey

[**API Update**](#api-survey-edit)  
  Method		: PUT  
URL			: {{url}}/master/v1/survey/:survey\_id

survey\_title tidak boleh sama jika effective date-nya overlap  
contoh 1:   
sekarang tgl 1 jan 2026   
•⁠ ⁠S001 \= saya buat data untuk tanggal 1 Feb 2026 \- 28 Feb 2026 —\> Survey A   
•⁠ ⁠S002 \= saya buat data untuk tanggal 1 Maret 2026 \- 31 Maret 2026 —\> Survey A   
Success , karena beda efective date 

contoh 2:   
sekarang tgl 1 jan 2026   
•⁠ ⁠S001 \= saya buat data untuk tanggal 1 Feb 2026 \- 28 Feb 2026 —\> Survey A   
•⁠ ⁠S002 \= saya buat data untuk tanggal 1 Feb 2026 \- 15 Feb 2026 —\> Survey A ini   
Failed, karena beda efective date 

Note :   
"Survey A" \= "survey a"  
 **dianggap sama.**

| if effective\_date\_start \> effective\_date\_end:     return error "effective\_date\_start must be \<= effective\_date\_end" exists \= query overlap\_check(survey\_title, effective\_date\_start, effective\_date\_end, survey\_id\_if\_update) if exists:     return error "Survey title already exists in overlapping effective date range"  |
| :---- |

## **Issue Create Survey : (case salesman)**

Issue BE : 

1) Tabel :   
- mst.m\_survey\_area \= kurang field distributor\_id   
- kurang tabel mst.m\_survey\_salesman [https://docs.google.com/document/d/1FRz0ym2cxYIlwCqjviwNLqOr4fyhZOZI2QCMyJ4wtmE/edit?tab=t.0](https://docs.google.com/document/d/1FRz0ym2cxYIlwCqjviwNLqOr4fyhZOZI2QCMyJ4wtmE/edit?tab=t.0) 

2) perbaiki bagian post 

**CURL :** 

| curl 'https://best.scyllax.online/master/v1/survey' \\   \-H 'sec-ch-ua-platform: "macOS"' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjowLCJkaXN0cmlidXRvcl9pZCI6MTAyLCJlbWFpbCI6ImFkbWluYm1AZ21haWwuY29tIiwiZW1wX2lkIjozODEsImV4cGlyZXMiOjE3NzU4MzA1NDYsImlzX2FkbWluIjpmYWxzZSwibGFuZ19pZCI6ImlkIiwibW9iaWxlX25vIjoiMDgxMzMzMzMzMzMzIiwicGFyZW50X2N1c3RfaWQiOiJDMjYwMDIiLCJ1c2VyX2Z1bGxuYW1lIjoiUGhpbGwgSm9uZXMiLCJ1c2VyX2lkIjoxNDEsInVzZXJfbmFtZSI6IlBoaWxsIEpvbmVzIiwid2hhdHNhcHAiOiIwODEzMzMzMzMzMzMifQ.WLGNdCqV5BoxDTL0iG9LU7HTbmp7EhmS-FuEsHi9Z-E' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'sec-ch-ua: "Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Content-Type: application/json' \\   \--data-raw '{"survey\_title":"survey salesman","efective\_date\_start":"2026-04-01","efective\_date\_end":"2026-04-18","answer\_frequency":"One Time","response\_type":"Mandatory","target\_type":"Specific","distributor\_id":\[102\],"area\_id":\[88\],"outlet\_id":\[\],"survey\_template\_id":40,"emp\_id":\[421\]}'  |
| :---- |

**Payload :** 

| {   "survey\_title": "survey salesman",   "efective\_date\_start": "2026-04-01",   "efective\_date\_end": "2026-04-18",   "answer\_frequency": "One Time",   "response\_type": "Mandatory",   "target\_type": "Specific",   "distributor\_id": \[102\],   "area\_id": \[88\],   "outlet\_id": \[\],   "survey\_template\_id": 40,   "emp\_id": \[421\] }  |
| :---- |

**pengaruh ke tabel :** 

* mst.m\_survey (DONE)  
* mst.m\_survey\_area (TO DO ) 

 insert ke field distributor\_id

* mst.m\_survey\_salesman  (TO DO ) 

   cust\_id            \= berdasarkan cust\_id user yang login   
   m\_survey\_salesman\_id\*  \= generate BE   
   survey\_id\*\*      \= berdasarkan survey\_id dari tabel mst.m\_survey   
   salesman\_id\*\*    \= berdasrkan emp\_id dari req body 

* mst.m\_survey\_outlet (DONE)  
* mst.m\_survey\_detail \= (DONE) 
