1. ## **API Template Survei List** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/master/v1/survey\_template?page=1\&limit=10\&sort:created\_date:desc\&status=1

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
| status | Integer | No | 1: mst.m\_survey\_template.is\_active \= true  0 : mst.m\_survey\_template.is\_active \= false Null : show all status |

### **Example Request** 

---

**Example Request Default**

| curl \--location \-g \\ '{{url}}/master/v1/survey\_template?page=1\&limit=5\&sort=created\_at:desc\&status=1' \\ \--header 'Accept: application/json' \\ \--header 'Authorization: Bearer {{token}}' |
| :---- |

### **Response** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| survey\_template\_id | int | 8 | **mst.m\_survey\_template.**survey\_template\_id |
| template\_code | varchar | 10 | **mst.m\_survey\_template.**template\_code |
| template\_title | varchar | 150 | **mst.m\_survey\_template.**template\_title |
| question\_total | int | 8 | **mst.m\_survey\_template.**question\_total |
| use\_image | bool |  | **mst.m\_survey\_template.**use\_image |
| created\_at | timestamptz | 6 | **mst.m\_survey\_template.**created\_at |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data {  "message": "Success",  "data": \[    {      "survey\_template\_id": 1,      "template\_code": "TMP001",      "template\_title": "Survey Kepuasan Pelanggan",      "question\_total": 10,      "use\_image": true,      "created\_at": "2025-01-01T10:00:00+07:00"    },    {      "survey\_template\_id": 2,      "template\_code": "TMP002",      "template\_title": "Survey Layanan",      "question\_total": 8,      "use\_image": false,      "created\_at": "2025-01-02T14:30:00+07:00"    }  \],  "paging": {    "total\_record": 84,    "page\_current": 1,    "page\_limit": 10,    "page\_total": 9  },  "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| Case : empty state / tidak terdapat data berdasarkan pencarian {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

2. ## **API Template Survei Detail** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/master/v1/survey\_template/:survey\_template\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Path Variable** 
---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| survey\_template\_id | Integer | Yes |  |

### **Example Request** 

---

**Example Request Default**

| curl \--location \-g \\ '{{url}}/master/v1/survey\_template/12' \\ \--header 'Accept: application/json' \\ \--header 'Authorization: Bearer {{token}}' |
| :---- |

### **Response** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| survey\_template\_id | int | 8 | **mst.m\_survey\_template.**survey\_template\_id |
| template\_code | varchar | 10 | **mst.m\_survey\_template.**template\_code |
| template\_title | varchar | 150 | **mst.m\_survey\_template.**template\_title |
| question\_total | int | 8 | **mst.m\_survey\_template.**question\_total |
| use\_image | bool |  | **mst.m\_survey\_template.**use\_image |
| created\_at | timestamptz | 6 | **mst.m\_survey\_template.**created\_at |
| **question\_template**  | **Array** |  |  |
| question\_template.survey\_template\_id | int | 8 | **mst.question\_template.**survey\_template\_id |
| question\_template.question | varchar | 225 | **mst.question\_template.**question |
| question\_template.answer\_type | enum |  | **mst.question\_template.**answer\_type |
| **m\_q\_option\_template** | **Array** |  |  |
| m\_q\_option\_template.q\_option\_template\_id | int | 8 | **mst.m\_q\_option\_template.**q\_option\_template\_id |
| m\_q\_option\_template.option | varchar | 225 | **mst.m\_q\_option\_template.**option |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data {  "message": "Success",  "data": \[    {      "survey\_template\_id": 12,      "template\_code": "TMP001",      "template\_title": "Survey Kepuasan Pelanggan",      "question\_total": 3,      "use\_image": true,      "created\_at": "2025-01-10T13:45:22+07:00",      "question\_template": \[        {          "survey\_template\_id": 12,          "question": "Bagaimana penilaian Anda terhadap pelayanan kami?",          "answer\_type": "Single",          "m\_q\_option\_template": \[\]        },        {          "survey\_template\_id": 12,          "question": "Fasilitas apa yang paling sering Anda gunakan?",          "answer\_type": "Multiple",          "m\_q\_option\_template": \[            {              "q\_option\_template\_id": 4,              "option": "Customer Service"            },            {              "q\_option\_template\_id": 5,              "option": "Kasir"            }          \]        },        {          "survey\_template\_id": 12,          "question": "Saran untuk peningkatan layanan kami",          "answer\_type": "Free Text",          "m\_q\_option\_template": \[\]        }      \]    }  \],  "paging": {    "total\_record": 84,    "page\_current": 1,    "page\_limit": 10,    "page\_total": 9  },  "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

3. ## **API Template Survei Add** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: POST  
URL		: {{url}}/master/v1/survey\_template

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Body** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| template\_title | varchar(150) | Yes |  |
| question\_total | int(8) | Yes |  |
| use\_image | bool | Yes |  |
| is\_active | bool | Yes |  |
| **question** | **Array** |  |  |
| question.question | varchar(225) | Yes |  |
| question.answer\_type | varchar(225) | Yes | Single, Multiple, Free Text |
| question.input\_type | varchar(225) | Yes | textfield, dropdown, radiobutton, toggle, checkbox, |
| **q\_option** | **Array** |  |  |
| q\_option.option | varchar(225) | No |  |

### **Example Request** 

---

**Example Request Default**

| curl \--location '{{url}}/master/v1/survey\_template' \\ \--header 'Accept: application/json' \\ \--header 'Content-Type: application/json' \\ \--header 'Authorization: Bearer {{token}}' \\ \--data '{  "template\_title": "Survey Kepuasan Pelanggan",  "question\_total": 3,  "use\_image": true,  "is\_active": true,  "question": \[    {      "question": "Bagaimana penilaian Anda terhadap pelayanan kami?",      "answer\_type": "Single",      "input\_type": "dropdown",      "q\_option": \[\]    },    {      "question": "Fasilitas apa yang paling sering Anda gunakan?",      "answer\_type": "Single",      "input\_type": "dropdown",      "q\_option": \[        { "option": "Customer Service" },        { "option": "Kasir" },        { "option": "Parkir" }      \]    },    {      "question": "Saran untuk peningkatan layanan kami",      "answer\_type": "Free Text",      "input\_type": "textfield",      "q\_option": \[\]    }  \] }'  |
| :---- |

### **Example Response** 

---

| Case : sukses dan terdapat data {  "message": "Survey template has been successfully created",  "request\_id": "6915a5e8e3f53f84fe73517f" }   |
| :---- |
| Case : error  {  "message": "Failed to create survey template",  "request\_id": "6915a5e8e3f53f84fe73517f" }   |

4. ## **API Template Survei Edit**  

 Create new endpoint :  
Content-Type 	: application/json  
Method		: PUT  
URL		: {{url}}/master/v1/survey\_template/:survey\_template\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Param** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| survey\_template\_id | int (8) | Yes | id survey template yang di edit  |

### **Body** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| template\_title | varchar(150) | Yes |  |
| question\_total | int(8) | Yes |  |
| use\_image | bool | Yes |  |
| is\_active | bool | Yes |  |
| **question** | **Array** |  |  |
| question.question | varchar(225) | Yes |  |
| question.answer\_type | varchar(225) | Yes | Single, Multiple, Free Text |
| question.input\_type | varchar(225) | Yes | textfield, dropdown, radiobutton, toggle, checkbox, |
| **q\_option** | **Array** |  |  |
| q\_option.option | varchar(225) | No |  |

### **Example Request** 

---

**Example Request Default**

| curl \--location \--request PATCH '{{url}}/master/v1/survey\_template/12' \\ \--header 'Accept: application/json' \\ \--header 'Content-Type: application/json' \\ \--header 'Authorization: Bearer {{token}}' \\ \--data '{  "template\_title": "Survey Kepuasan Pelanggan (Update)",  "question\_total": 3,  "use\_image": false,  "is\_active": true,  "question": \[    {      "question\_template\_id" : 23,      "question": "Bagaimana penilaian Anda terhadap pelayanan kami?",      "answer\_type": "Single",      "input\_type": "dropdown",      "q\_option": \[        { "option": "Sangat Baik" },        { "option": "Baik" },        { "option": "Cukup" }      \]    },    {      "question\_template\_id" : 24,      "question": "Fasilitas apa yang paling sering Anda gunakan?",      "answer\_type": "Multiple",      "input\_type": "dropdown",      "q\_option": \[        { "option": "Customer Service" },        { "option": "Kasir" }      \]    },    {      "question\_template\_id" : 25,      "question": "Saran untuk peningkatan layanan kami",      "answer\_type": "Free Text",      "input\_type": "textfield",      "q\_option": \[\]    }  \] }'  |
| :---- |

### **Example Response** 

---

| Case : sukses dan terdapat data {   "message": "Survey template berhasil dibuat",   "request\_id": "6915a5e8e3f53f84fe73517f" }   |
| :---- |
| Case : error  {   "message": "Terjadi kesalahan pada server",   "error": {     "description": "Unexpected error"   },   "request\_id": "6915a5e8e3f53f84fe73517f" }  |

5. ## **API Template Survei Delete** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: DELETE  
URL		: {{url}}/master/v1/template\_survey/:survey\_template\_id

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Param** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| survey\_template\_id | int (8) | Yes | id survey template yang di edit  |

### **Example Request** 

---

**Example Request Default**

| curl \--location \--request DELETE \\ '{{url}}/master/v1/template\_survey/12' \\ \--header 'Accept: application/json' \\ \--header 'Authorization: Bearer {{token}}'  |
| :---- |

### **Example Response :** 

| Case : sukses  {   "message": "Survey template has been successfully deleted",   "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| Case : error  {   "message": "Survey template not found",   "request\_id": "6915a5e8e3f53f84fe73517f" }  |

