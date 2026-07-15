

1. ## **API View Invoice**

 Enhance endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/finance/v1/account-receivables/collection/invoices?limit=999999\&page=1\&inv\_date\_from=1767225600\&inv\_date\_to=1769903999\&due\_date\_from=1767225600\&due\_date\_to=1769903999\&salesman\_id=346

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| q |  | No | template\_title |
| page | Integer | Yes | default : 1 |
| limit  | Integer | Yes | 999999 |
| inv\_date\_from | epoch  | Yes | 1767225600 |
| inv\_date\_to | epoch | Yes | 1769903999 |
| due\_date\_from | epoch | Yes | 1767225600 |
| due\_date\_to | epoch | Yes | 1769903999 |
| salesman\_id | **Before** : int **After** : Array\<int\> | Yes | 346 |

### **Example Request :** 

| curl 'https://best.scyllax.online/finance/v1/account-receivables/collection/invoices?limit=999999\&page=1\&inv\_date\_from=1767225600\&inv\_date\_to=1769903999\&due\_date\_from=1767225600\&due\_date\_to=1769903999\&salesman\_id=346' \\   \-H 'sec-ch-ua-platform: "macOS"' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoic3lhaHJpemFAZ21haWwuY29tIiwiZW1wX2lkIjoyMTAsImV4cGlyZXMiOjE3Njk3NjEwNDMsImlzX2FkbWluIjpmYWxzZSwibGFuZ19pZCI6ImlkIiwibW9iaWxlX25vIjoiMDg1NzQ4MTIzMTIyMiIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IlN5YWhyaXphIFNhbGVzbWFuIiwidXNlcl9pZCI6MTAwLCJ1c2VyX25hbWUiOiJTeWFocml6YSIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4MiJ9.b8yHFqeRaHY9XhFFgiVRvYB-ZfROOxTUk\_2jq0YF3qk' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'sec-ch-ua: "Not(A:Brand";v="8", "Chromium";v="144", "Google Chrome";v="144"' \\   \-H 'sec-ch-ua-mobile: ?0' |
| :---- |

### **Response :** 

---

| {     "message": "",     "data": \[         {             "invoice\_no": "INV2601280001",             "invoice\_date": "2026-01-28",             "ro\_no": "SO2601270018",             "due\_date": "2026-01-28",             "outlet\_id": 420,             "outlet\_code": "LV000010",             "outlet\_name": "Lavina TK 10",             "salesman\_id": 346,             "salesman\_code": "DR0123",             "salesman\_name": "Dicky Ramdhan",             "invoice\_amount": 470400,             "remaining\_amount": 470400,             "paid\_amount": 0         },         {             "invoice\_no": "INV2601280002",             "invoice\_date": "2026-01-28",             "ro\_no": "SO2601270017",             "due\_date": "2026-01-28",             "outlet\_id": 419,             "outlet\_code": "LV000011",             "outlet\_name": "Lavina TK 11",             "salesman\_id": 346,             "salesman\_code": "DR0123",             "salesman\_name": "Dicky Ramdhan",             "invoice\_amount": 2256850,             "remaining\_amount": 2006850,             "paid\_amount": 250000         },         {             "invoice\_no": "INV2601280003",             "invoice\_date": "2026-01-28",             "ro\_no": "SO2601270016",             "due\_date": "2026-01-28",             "outlet\_id": 418,             "outlet\_code": "LV000013",             "outlet\_name": "Lavina TK 13",             "salesman\_id": 346,             "salesman\_code": "DR0123",             "salesman\_name": "Dicky Ramdhan",             "invoice\_amount": 92685,             "remaining\_amount": 92685,             "paid\_amount": 0         },         {             "invoice\_no": "INV2601280004",             "invoice\_date": "2026-01-28",             "ro\_no": "SO2601270011",             "due\_date": "2026-01-28",             "outlet\_id": 416,             "outlet\_code": "LV000009",             "outlet\_name": "Lavina TK 9",             "salesman\_id": 346,             "salesman\_code": "DR0123",             "salesman\_name": "Dicky Ramdhan",             "invoice\_amount": 97000,             "remaining\_amount": 97000,             "paid\_amount": 0         },         {             "invoice\_no": "INV2601280005",             "invoice\_date": "2026-01-28",             "ro\_no": "SO2601270010",             "due\_date": "2026-01-28",             "outlet\_id": 416,             "outlet\_code": "LV000009",             "outlet\_name": "Lavina TK 9",             "salesman\_id": 346,             "salesman\_code": "DR0123",             "salesman\_name": "Dicky Ramdhan",             "invoice\_amount": 97000,             "remaining\_amount": 97000,             "paid\_amount": 0         },         {             "invoice\_no": "INV2601190002",             "invoice\_date": "2026-01-19",             "ro\_no": "SO2601190003",             "due\_date": "2026-01-26",             "outlet\_id": 580,             "outlet\_code": "L000141",             "outlet\_name": "LOLLY 141",             "salesman\_id": 346,             "salesman\_code": "DR0123",             "salesman\_name": "Dicky Ramdhan",             "invoice\_amount": 651840,             "remaining\_amount": 651840,             "paid\_amount": 0         }     \],     "paging": {         "total\_record": 6,         "page\_current": 1,         "page\_limit": 999999,         "page\_total": 1     },     "request\_id": "697b71f7d8e779a4dbc62340" }  |
| :---- |

2. ## **Create Collection**

 Enhance endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/finance/v1/account-receivables/collection

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Example Request :** 

| {   "collection\_date": "2026-01-01",   "emp\_id": 243,   "ot\_grp\_id": 0,   "total\_amount": 474524,   "remaining\_amount": 125668,   "invoice\_date\_from": "2026-01-01T00:00:00Z",   "invoice\_date\_to": "2026-01-31T00:00:00Z",   "due\_date\_from": "2026-01-01T00:00:00Z",   "due\_date\_to": "2026-01-31T00:00:00Z",   "notes": "",   "details": \[     {       "invoice\_no": "INV2601120001",       "salesman\_id": 234, //tambahkan req salesman\_id       "invoice\_amount": 474524,       "remaining\_amount": 125668     }   \] }  |
| :---- |

### **Response :** 

---

| {      "message":"Berhasil Dibuat",      "request\_id":"697b75d9d8e779a4dbc6234f" } |
| :---- |

3. ## **API Collection Detail** 

 No Enhance :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/finance/v1/account-receivables/collection/CL2601290001

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Example Request :** 

| curl 'https://best.scyllax.online/finance/v1/account-receivables/collection/CL2601290001' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoic3lhaHJpemFAZ21haWwuY29tIiwiZW1wX2lkIjoyMTAsImV4cGlyZXMiOjE3Njk3NjEwNDMsImlzX2FkbWluIjpmYWxzZSwibGFuZ19pZCI6ImlkIiwibW9iaWxlX25vIjoiMDg1NzQ4MTIzMTIyMiIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IlN5YWhyaXphIFNhbGVzbWFuIiwidXNlcl9pZCI6MTAwLCJ1c2VyX25hbWUiOiJTeWFocml6YSIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4MiJ9.b8yHFqeRaHY9XhFFgiVRvYB-ZfROOxTUk\_2jq0YF3qk' \\   \-H 'Connection: keep-alive' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Not(A:Brand";v="8", "Chromium";v="144", "Google Chrome";v="144"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"'  |
| :---- |

### **Response :** 

**enhance** : 

- invoice\_amount  	\= Ekspektasi Business → nilai invoice **acf.collection\_det.invoice\_amount**

- invoice\_payment  	\= Ekspektasi Business → payment deposit by collection \+ payment deposit by invoice   
  **acf.collection\_det.paid\_by\_invoice \+ acf.collection\_det.paid\_amount**

- paid\_amount		\= Ekspektasi Business → payment deposit by collection  
   **acf.collection\_det.paid\_amount** 

- remaining\_amount	\= Ekspektasi Business → total amount \- (pembayaran by invoice \+ pembayaran by collection)   
  **acf.collection.invoice\_amount \- (acf.collection\_det.paid\_by\_invoice \+ acf.collection\_det.paid\_amount)**

### ---

| {     "message": "",     "data": {         "cust\_id": "C220010001",         "collection\_no": "CL2601290001",         "collection\_date": "2026-01-29",         "emp\_id": 260,         "emp\_code": "EMP1001",         "emp\_name": "ANTONNY",         "ot\_grp\_id": null,         "ot\_grp\_code": null,         "ot\_grp\_name": null,         "notes": "",         "total\_amount": 5829769,         "remaining\_amount": 5829769,         "invoice\_date\_from": "2026-01-01",         "invoice\_date\_to": "2026-01-31",         "due\_date\_from": "2026-01-01",         "due\_date\_to": "2026-01-31",         "created\_by": 12,         "created\_by\_name": "Dist IDE Sda",         "created\_at": "2026-01-29T04:47:40.123409Z",         "updated\_by": 12,         "updated\_by\_name": "Dist IDE Sda",         "updated\_at": "2026-01-29T04:48:31.740898Z",         "deleted\_by": null,         "deleted\_by\_name": null,         "deleted\_at": null,         "printed\_by": 12,         "printed\_by\_name": "Dist IDE Sda",         "printed\_at": "2026-01-29T04:48:31.740838Z",         "details": \[             {                 "collection\_det\_id": 420,                 "collection\_no": "CL2601290001",                 "invoice\_no": "INV2601130037",                 "sales\_order": "SO2601120013",                 "invoice\_date": "2026-01-13",                 "due\_date": "2026-01-13",                 "salesman\_id": 261,                 "salesman\_name": "Nyeck Nyobe",                 "salesman\_code": "BDG001",                 "outlet\_id": 190,                 "outlet\_code": "250114",                 "outlet\_name": "Fahmi, Tk",                 "invoice\_amount": 1032300, //collection\_det                 "invoice\_payment": 4797469,                 "remaining\_amount": 1032300,                 "paid\_amount": 0,                   "total\_invoice\_amount": 0,                 "created\_by": 12,                 "created\_by\_name": "Dist IDE Sda",                 "created\_at": "2026-01-29T04:47:40.243551Z"             },             {                 "collection\_det\_id": 419,                 "collection\_no": "CL2601290001",                 "invoice\_no": "INV2601130038",                 "sales\_order": "SO2512230004",                 "invoice\_date": "2026-01-13",                 "due\_date": "2026-01-13",                 "salesman\_id": 261,                 "salesman\_name": "Nyeck Nyobe",                 "salesman\_code": "BDG001",                 "outlet\_id": 258,                 "outlet\_code": "tes123",                 "outlet\_name": "tes approval order",                 "invoice\_amount": 4797469,                 "invoice\_payment": 4797469,                 "remaining\_amount": 4797469,                 "paid\_amount": 0,                 "total\_invoice\_amount": 0,                 "created\_by": 12,                 "created\_by\_name": "Dist IDE Sda",                 "created\_at": "2026-01-29T04:47:40.125503Z"             }         \]     },     "request\_id": "697b785fd8e779a4dbc6235b" }  |
| :---- |

4. ## **Edit Collection**

 Create new endpoint :  
Content-Type 	: application/json  
Method		: PUT  
URL		: {{url}}/finance/v1/account-receivables/collection/:collection\_no  
FE kirim all in 

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Example Request :** 

| {   "collection\_date": "2026-01-01",   "emp\_id": 243,   "ot\_grp\_id": 0,   "total\_amount": 474524,   "remaining\_amount": 125668,   "invoice\_date\_from": "2026-01-01T00:00:00Z",   "invoice\_date\_to": "2026-01-31T00:00:00Z",   "due\_date\_from": "2026-01-01T00:00:00Z",   "due\_date\_to": "2026-01-31T00:00:00Z",   "notes": "",   "details": \[     {       "invoice\_no": "INV2601120001",       "salesman\_id": 234, //tambahkan req salesman\_id       "invoice\_amount": 474524,       "remaining\_amount": 125668     }   \] }  |
| :---- |

### **Response :** 

---

| {      "message":"Berhasil Dibuat",      "request\_id":"697b75d9d8e779a4dbc6234f" } |
| :---- |

