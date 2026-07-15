1. ## **Outlet List Approval** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/master/v1/outlet-list?page=1\&limit=10\&sort:created\_date:desc\&status=1

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
| sort | String | Yes | default created\_date:desc  |
| status | Array(int)  | No  | 1 : Need Review  2 : Approved 	 3 : Rejected  |

### **Example Request** 

---

**Example Request Default**

| curl \--location \-g \\ '{{url}}/master/v1/outlet-list?page=1\&limit=5\&sort=created\_date:desc\&status=1' \\ \--header 'Accept: application/json' |
| :---- |

### **Response** 

---

* Jika FE tidak mengirimkan param status \= show all status 

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| message | String | 150 | Response Message |
| data | Array | \- |  \- |
| outlet\_cr\_id | Integer | 8 | **mst.outlet\_cr.**outlet\_cr\_id |
| outlet\_id | Integer | 8 | **mst.outlet\_cr.** |
| oultet\_code | varchar | 30 | **mst.m\_outlet**.outlet\_code |
| outlet\_name | varchar | 50 | **mst.m\_outlet**.outlet\_name |
| current\_long | varchar | 50 | **mst.outlet\_cr\_det.**current\_long |
| current\_lat | varchar | 50 | **mst.outlet\_cr\_det.**current\_lat |
| new\_long | varchar | 50 | **mst.outlet\_cr\_det.**new\_long |
| new\_lat | varchar | 50 | **mst.outlet\_cr\_det.**new\_lat |
| source | Integer | 8 | **mst.outlet\_cr.**source |
| status | Integer | 3 | **mst.outlet\_cr.**status |
| status\_desc | varchar | 50 | Need Review	: 1  Approved 	: 2 Rejected 	: 3 |
| request\_by | varchar | 50 | **sys.m\_user.**user\_namerelasi : **mst.outlet\_cr.**creted\_at \= sys.m\_user.user\_id |
| request\_date | timestamp | 6 | **mst.outlet\_cr.**creted\_at |
| paging  | Object  |  |  |
| total\_record | Numeric | 11 | Total data seluruh halaman |
| page\_current | Numeric | 11 | Halaman saat ini  |
| page\_limit | Numeric | 11 | Data yang ditampilkan per page  |
| page\_total | Numeric | 11 | Total page keseluruhan |
| request\_id | String | 150 | Generate request id  |

**Example Response :** 

| Case : sukses dan terdapat data {  "message": "",  "data": \[    {      "outlet\_cr\_id": 12,      "outlet\_id": 1001,      "outlet\_code": "OUT-001",      "outlet\_name": "Outlet Contoh",      "current\_long": "106.812345",      "current\_lat": "-6.201234",      "new\_long": "106.812999",      "new\_lat": "-6.201999",      "source": 1,      "status": 1,      "status\_desc": "Need Review",      "request\_by": "Admin User",      "request\_date": "2025-12-01T10:15:00"    },    {      "outlet\_cr\_id": 13,      "outlet\_id": 1002,      "outlet\_code": "OUT-002",      "outlet\_name": "Outlet Dua",      "current\_long": "107.002200",      "current\_lat": "-6.120330",      "new\_long": "107.002450",      "new\_lat": "-6.120500",      "source": 2,      "status": 2,      "status\_desc": "Approved",      "request\_by": "Supervisor",      "request\_date": "2025-12-02T14:22:00"    }  \],  "paging": {    "total\_record": 84,    "page\_current": 1,    "page\_limit": 10,    "page\_total": 9  },  "request\_id": "6915a5e8e3f53f84fe73517f" }   |
| :---- |
| Case : empty state / tidak terdapat data berdasarkan pencarian {    "message": "No Data",    "data": null,    "paging": {        "total\_record": 0,        "page\_current": 1,        "page\_limit": 10,        "page\_total": 0    },    "request\_id": "6915a6dd2395083c685e8e16" }  |
| Case : error  {    "message": "record not found",    "request\_id": "6915a6dd2395083c685e8e16" }  |

2. ## **Approval Outlet List** 

 Create new endpoint :  
Content-Type 	: application/json  
Method		: PATCH  
URL		: {{url}}/master/v1/outlet-list/approval

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Body** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| outlet\_cr\_id | Array(int) | Yes |  |
| status | Integer | Yes | 2 : Approved 	 3 : Rejected  |

### **Example Request** 

---

**Example Request Default**

| curl \--location \-g \\ '{{url}}/master/v1/outlet-list/approval' \\ \--header 'Accept: application/json' |
| :---- |

**Example Response :** 

| Case : sukses approval  {  "message": "Approval successfully",  "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |
| Case : gagal approval  {  "message": "Approval failed",  "request\_id": "6915a5e8e3f53f84fe73517f" }  |

* Be update tabel mst.outlet\_cr.status   
* Jika status yang dikirim **FE \= 2** , update di mst.m\_outlet    
  Nilai dilihat dari tabel mst.outlet\_cr\_det field\_name \= latitude / longitude , ambil value dari new\_value, dan berdasarkan outlet\_id yang dikirim FE   
- latitude   
- longitude   
