1. **Homepage (Dashboard) Secondary Sales Report**

Saat ini, user Principal tidak dapat memilih data distributor berdasarkan wilayah tertentu pada halaman Secondary Sales Report. Hal ini menyebabkan proses monitoring dan analisa performa salesman menjadi kurang fleksibel dan sulit difokuskan pada area tertentu.

Enhancement ini akan menambahkan filter wilayah berupa:

* Region  
* Area  
* Distributor

2. **Modal Export Report**

Saat ini, user Principal tidak dapat memilih data distributor berdasarkan wilayah tertentu pada halaman Secondary Sales Report. Hal ini menyebabkan proses monitoring dan analisa performa salesman menjadi kurang fleksibel dan sulit difokuskan pada area tertentu.

Enhancement ini akan menambahkan filter wilayah berupa:

* Region  
* Area  
* Distributor


1. # **Trend Sales**

 Enhance endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/sales/v1/reports/secondary-sales/trend-sales?year=2026

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Example Request** 

---

| curl 'https://best.scyllax.online/sales/v1/reports/secondary-sales/trend-sales?year=2026' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjowLCJkaXN0cmlidXRvcl9pZCI6MTAyLCJlbWFpbCI6ImFkbWluYm1AZ21haWwuY29tIiwiZW1wX2lkIjozODEsImV4cGlyZXMiOjE3Nzg4OTk0OTUsImlzX2FkbWluIjpmYWxzZSwibGFuZ19pZCI6ImlkIiwibW9iaWxlX25vIjoiMDgxMzMzMzMzMzMzIiwicGFyZW50X2N1c3RfaWQiOiJDMjYwMDIiLCJ1c2VyX2Z1bGxuYW1lIjoiUGhpbGwgSm9uZXMiLCJ1c2VyX2lkIjoxNDEsInVzZXJfbmFtZSI6IlBoaWxsIEpvbmVzIiwid2hhdHNhcHAiOiIwODEzMzMzMzMzMzMifQ.k8XvFu9M-whaoVK6vg5DAWrHCZit-8XlKPawR\_AHReE' \\   \-H 'Connection: keep-alive' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Chromium";v="148", "Google Chrome";v="148", "Not/A)Brand";v="99"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"' |
| :---- |
| **Payload** year : 2026 **Tambahkan request body berikut :**  cust\_id (untuk filter berdasarkan business unit → report.fact\_orders.cust\_id)  Untuk user principal, karena mempunyai child cust\_id , maka cust\_id bisa berbeda-beda , bisa cust\_id milik user principal itu sendiri, atau cust\_id milik distributor Untuk user distributor, karena hanya mempunyai 1 cust\_id , maka cust\_id bisa berdasarkan user yang login  |

### **Response JSON**

---

| {     "message": "",     "data": \[         {             "month": 1,             "total\_gross\_sale": 0,             "total\_discount\_promo": 0,             "net\_sales": 0         },         {             "month": 2,             "total\_gross\_sale": 0,             "total\_discount\_promo": 0,             "net\_sales": 0         },         {             "month": 3,             "total\_gross\_sale": 0,             "total\_discount\_promo": 0,             "net\_sales": 0         },         {             "month": 4,             "total\_gross\_sale": 469500000,             "total\_discount\_promo": 2132000,             "net\_sales": 477368000         },         {             "month": 5,             "total\_gross\_sale": 0,             "total\_discount\_promo": 0,             "net\_sales": 0         },         {             "month": 6,             "total\_gross\_sale": 0,             "total\_discount\_promo": 0,             "net\_sales": 0         },         {             "month": 7,             "total\_gross\_sale": 0,             "total\_discount\_promo": 0,             "net\_sales": 0         },         {             "month": 8,             "total\_gross\_sale": 0,             "total\_discount\_promo": 0,             "net\_sales": 0         },         {             "month": 9,             "total\_gross\_sale": 0,             "total\_discount\_promo": 0,             "net\_sales": 0         },         {             "month": 10,             "total\_gross\_sale": 0,             "total\_discount\_promo": 0,             "net\_sales": 0         },         {             "month": 11,             "total\_gross\_sale": 0,             "total\_discount\_promo": 0,             "net\_sales": 0         },         {             "month": 12,             "total\_gross\_sale": 0,             "total\_discount\_promo": 0,             "net\_sales": 0         }     \],     "request\_id": "6f6571f9-378a-4364-b25d-705ec2382fea" }  |
| :---- |

2. # **SUM Date**

 Enhance endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/sales/v1/reports/secondary-sales/sum-date?month=5

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Example Request** 

---

| curl 'https://best.scyllax.online/sales/v1/reports/secondary-sales/sum-date?month=5' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjowLCJkaXN0cmlidXRvcl9pZCI6MTAyLCJlbWFpbCI6ImFkbWluYm1AZ21haWwuY29tIiwiZW1wX2lkIjozODEsImV4cGlyZXMiOjE3Nzg4OTk0OTUsImlzX2FkbWluIjpmYWxzZSwibGFuZ19pZCI6ImlkIiwibW9iaWxlX25vIjoiMDgxMzMzMzMzMzMzIiwicGFyZW50X2N1c3RfaWQiOiJDMjYwMDIiLCJ1c2VyX2Z1bGxuYW1lIjoiUGhpbGwgSm9uZXMiLCJ1c2VyX2lkIjoxNDEsInVzZXJfbmFtZSI6IlBoaWxsIEpvbmVzIiwid2hhdHNhcHAiOiIwODEzMzMzMzMzMzMifQ.k8XvFu9M-whaoVK6vg5DAWrHCZit-8XlKPawR\_AHReE' \\   \-H 'Connection: keep-alive' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Chromium";v="148", "Google Chrome";v="148", "Not/A)Brand";v="99"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"'  |
| :---- |
| **Request Body**  month=5  **Tambahkan request body berikut :**  year \= integer  (report.dim\_dates) cust\_id (untuk filter berdasarkan business unit → report.fact\_orders.cust\_id)  Untuk user principal, karena mempunyai child cust\_id , maka cust\_id bisa berbeda-beda , bisa cust\_id milik user principal itu sendiri, atau cust\_id milik distributor Untuk user distributor, karena hanya mempunyai 1 cust\_id , maka cust\_id bisa berdasarkan user yang login   |

### **Response JSON**

---

| {     "message": "",     "data": {         "total\_gross\_sale": 0,         "total\_discount\_promo": 0,         "net\_sales": 0,         "total\_salesman": 0,         "total\_outlet": 0,         "total\_product": 0,         "qty": 0,         "qty\_return": 0,         "return\_rate": 0,         "net\_sales\_return": 0,         "last\_update": null     },     "request\_id": "6e6571f9-378a-4364-b25d-705ec2382fea" }  |
| :---- |

### **Query Eksisting:** 

---

| SALES SUMMARY SELECT SUM(report.fact\_orders.gross\_sale) AS *total\_gross\_sale*, SUM(report.fact\_orders.discount \+ report.fact\_orders.special\_discount) AS *total\_discount\_promo*, SUM(report.fact\_orders.net\_sales\_exclude\_ppn) AS *net\_sales*, COUNT(DISTINCT(salesman\_id)) AS *total\_salesman*, COUNT(DISTINCT(outlet\_id)) AS *total\_outlet*, COUNT(DISTINCT(pro\_id)) AS *total\_product*, SUM(qty) AS *qty*, MAX(extracted\_at) AS *last\_update* FROM report.fact\_orders JOIN report.dim\_dates *dt* ON report.fact\_orders.date\_id \= *dt*.id WHERE report.fact\_orders.cust\_id \= 'C2600200001' AND *dt*.month \= 5 LIMIT 1  |
| :---- |

3. # **Secondary Sales Group** 

 Enhance endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: {{url}}/sales/v1/reports/secondary-sales/group?month=5\&group\_by=outlet

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **Example Request** 

---

| curl 'https://best.scyllax.online/sales/v1/reports/secondary-sales/group?month=5\&group\_by=outlet' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjowLCJkaXN0cmlidXRvcl9pZCI6MTAyLCJlbWFpbCI6ImFkbWluYm1AZ21haWwuY29tIiwiZW1wX2lkIjozODEsImV4cGlyZXMiOjE3Nzg4OTk0OTUsImlzX2FkbWluIjpmYWxzZSwibGFuZ19pZCI6ImlkIiwibW9iaWxlX25vIjoiMDgxMzMzMzMzMzMzIiwicGFyZW50X2N1c3RfaWQiOiJDMjYwMDIiLCJ1c2VyX2Z1bGxuYW1lIjoiUGhpbGwgSm9uZXMiLCJ1c2VyX2lkIjoxNDEsInVzZXJfbmFtZSI6IlBoaWxsIEpvbmVzIiwid2hhdHNhcHAiOiIwODEzMzMzMzMzMzMifQ.k8XvFu9M-whaoVK6vg5DAWrHCZit-8XlKPawR\_AHReE' \\   \-H 'Connection: keep-alive' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Chromium";v="148", "Google Chrome";v="148", "Not/A)Brand";v="99"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"'  |
| :---- |
| **Payload** month=5\&group\_by=outlet **Tambahkan request body berikut :**  year \= integer  cust\_id (untuk filter berdasarkan business unit → report.fact\_orders.cust\_id)  Untuk user principal, karena mempunyai child cust\_id , maka cust\_id bisa berbeda-beda , bisa cust\_id milik user principal itu sendiri, atau cust\_id milik distributor Untuk user distributor, karena hanya mempunyai 1 cust\_id , maka cust\_id bisa berdasarkan user yang login  |

### **Response JSON**

---

| {"message":"","data":\[\],"request\_id":"776571f9-378a-4364-b25d-705ec2382fea"} |
| :---- |

4. # **Export Secondary Sales** 

 Enhance endpoint :  
Content-Type 	: application/json  
Method		: POST  
URL		: {{url}}/sales/v1/reports/secondary-sales

### **CURL** 

| curl 'https://best.scyllax.online/sales/v1/reports/secondary-sales' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsImVtcF9pZCI6MCwiZXhwaXJlcyI6MTc3OTE1NzQzNCwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4NTc0ODEyMzEyOCIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IkRpc3QgSURFIFNkYSIsInVzZXJfaWQiOjEyLCJ1c2VyX25hbWUiOiJkaXN0QHNkYS5pZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODU3NDgxMjMxMjgifQ.gJy9fiCXWTXM1joP1ifm7v7kXoXaHsPfdI\_6cT4NnWQ' \\   \-H 'Connection: keep-alive' \\   \-H 'Content-Type: application/json' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Chromium";v="148", "Google Chrome";v="148", "Not/A)Brand";v="99"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"' \\   \--data-raw '{"from":1777568400,"to":1779123599,"distributor\_ids":\[\],"outlet\_ids":\[1903,1900,1899\],"salesman\_ids":\[\],"pro\_ids":\[10756,10741,10740\]}'  |
| :---- |

### **Response** 

| {     "message": "",     "data": {         "report\_id": "6a0a805cb486af7bf375fe65",         "report\_name": "SecondarySales-180526-002",         "start\_date": "2026-05-01",         "end\_date": "2026-05-18",         "file\_status": 2,         "file\_status\_name": "",         "file\_url": "",         "created\_by": "Dist IDE Sda",         "created\_at": "2026-05-18T02:58:36.985201953Z"     },     "paging": {         "total\_record": 1,         "page\_current": 0,         "page\_limit": 0,         "page\_total": 1     },     "request\_id": "806d15b5-85a3-4133-a94b-59e3d1e88406" }  |
| :---- |

###  **Template  \=** click [here](https://docs.google.com/spreadsheets/d/1WkYPe81DgL6Y3SNsUjdnlw-DXYXNJ7e9/edit?gid=813953497#gid=813953497)

Query : 

| \--SECONDARY SALES TRX TYPE \= ORDER----- WITH *order\_data* AS (     SELECT          *o*.invoice\_no AS *document\_no*,         *o*.invoice\_date AS *document\_date*,         *o*.outlet\_id,         *o*.salesman\_id,         *o*.cust\_id,         *od*.pro\_id,         *od*.unit\_id1,         *od*.unit\_id2,         *od*.unit\_id3,         *od*.conv\_unit2,         *od*.conv\_unit3,         COALESCE(*od*.qty1\_final, 0\) AS *qty1*,         COALESCE(*od*.qty2\_final, 0\) AS *qty2*,         COALESCE(*od*.qty3\_final, 0\) AS *qty3*,         COALESCE(*od*.sell\_price\_final1, 0\) AS *price1*,         COALESCE(*od*.sell\_price\_final2, 0\) AS *price2*,         COALESCE(*od*.sell\_price\_final3, 0\) AS *price3*,         COALESCE(*od*.disc\_value\_final, 0\) AS *special\_discount*,         COALESCE(             *od*.promo\_final1 \+             *od*.promo\_final2 \+             *od*.promo\_final3 \+             *od*.promo\_final4 \+             *od*.promo\_final5,             0         ) AS *discount*,         COALESCE(*od*.vat\_value\_final, 0\) AS *ppn*,         'ORDER' AS *trx\_type*     FROM sls."order" *o*     JOIN sls.order\_detail *od*          ON *od*.ro\_no \= *o*.ro\_no        AND *od*.cust\_id \= *o*.cust\_id     WHERE *o*.cust\_id \= :cust\_id       AND *o*.data\_status IN (6,7)       AND (             :date\_from IS NULL             OR :date\_to IS NULL             OR o.invoice\_date BETWEEN :date\_from AND :date\_to       )       AND (             :outlet\_ids IS NULL             OR o.outlet\_id \= ANY(:outlet\_ids)       )       AND (             :salesman\_ids IS NULL             OR o.salesman\_id \= ANY(:salesman\_ids)       )       AND (             :product\_ids IS NULL             OR od.pro\_id \= ANY(:product\_ids)       ) ), *return\_data* AS (     SELECT         *r*.return\_no AS *document\_no*,         *r*.return\_date AS *document\_date*,         *r*.outlet\_id,         *r*.salesman\_id,         *r*.cust\_id,         *rd*.product\_id AS *pro\_id*,         *rd*.unit\_id1,         *rd*.unit\_id2,         *rd*.unit\_id3,         *rd*.conv\_unit2,         *rd*.conv\_unit3,         COALESCE(*rd*.qty1, 0\) AS *qty1*,         COALESCE(*rd*.qty2, 0\) AS *qty2*,         COALESCE(*rd*.qty3, 0\) AS *qty3*,         COALESCE(*rd*.sell\_price1, 0\) AS *price1*,         COALESCE(*rd*.sell\_price2, 0\) AS *price2*,         COALESCE(*rd*.sell\_price3, 0\) AS *price3*,         COALESCE(*rd*.disc\_value, 0\) AS *special\_discount*,         0 AS *discount*,         COALESCE(*rd*.vat\_value, 0\) AS *ppn*,         'RETURN' AS *trx\_type*     FROM sls.return\_det *rd*     JOIN sls."return" *r*         ON *r*.return\_no \= *rd*.return\_no        AND *r*.cust\_id \= *rd*.cust\_id     JOIN sls."order" *o*         ON *o*.invoice\_no \= *r*.invoice\_no        AND *o*.cust\_id \= *r*.cust\_id     WHERE *rd*.cust\_id \= :cust\_id       AND *o*.data\_status IN (6,7)       AND (             :date\_from IS NULL             OR :date\_to IS NULL             OR o.invoice\_date BETWEEN :date\_from AND :date\_to       )       AND (             :outlet\_ids IS NULL             OR r.outlet\_id \= ANY(:outlet\_ids)       )       AND (             :salesman\_ids IS NULL             OR r.salesman\_id \= ANY(:salesman\_ids)       )       AND (             :product\_ids IS NULL             OR rd.product\_id \= ANY(:product\_ids)       ) ), *trx* AS (     SELECT \* FROM *order\_data*     UNION ALL     SELECT \* FROM *return\_data* ) SELECT     *md*.distributor\_code,     *md*.distributor\_name,     *t*.*trx\_type*,     *t*.*document\_no*,     *t*.*document\_date*,     *mo*.outlet\_code,     *mo*.outlet\_principal\_code,     *mo*.outlet\_name,     *me*.emp\_code,     *ms*.sales\_name AS *emp\_name*,     *ms2*.sup\_code,     *ms2*.sup\_name,     *mp*.pro\_code,     *mp*.pro\_name,     *t*.*price3* AS *price\_unit3*,     *t*.*price2* AS *price\_unit2*,     *t*.*price1* AS *price\_unit1*,     *t*.unit\_id3,     *t*.unit\_id2,     *t*.unit\_id1,     *t*.conv\_unit2,     *t*.conv\_unit3,     *t*.*qty3* AS *qty3\_final*,     *t*.*qty2* AS *qty2\_final*,     *t*.*qty1* AS *qty1\_final*,     (         (*t*.*qty1* \* *t*.*price1*) \+         (*t*.*qty2* \* *t*.*price2*) \+         (*t*.*qty3* \* *t*.*price3*)     ) AS *gross\_sales*,     *t*.*special\_discount*,     *t*.*discount*,     (         (             (*t*.*qty1* \* *t*.*price1*) \+             (*t*.*qty2* \* *t*.*price2*) \+             (*t*.*qty3* \* *t*.*price3*)         )         \- *t*.*special\_discount*         \- *t*.*discount*     ) AS *net\_sales\_exc\_ppn*,     *t*.*ppn*,     (         (             (                 (*t*.*qty1* \* *t*.*price1*) \+                 (*t*.*qty2* \* *t*.*price2*) \+                 (*t*.*qty3* \* *t*.*price3*)             )             \- *t*.*special\_discount*             \- *t*.*discount*         )         \+ *t*.*ppn*     ) AS *net\_sales\_inc\_ppn* FROM *trx* *t* LEFT JOIN smc.m\_customer *mc*     ON *mc*.cust\_id \= *t*.cust\_id LEFT JOIN mst.m\_distributor *md*     ON *md*.distributor\_id \= *mc*.distributor\_id    AND *md*.cust\_id \= *mc*.cust\_id LEFT JOIN mst.m\_outlet *mo*     ON *mo*.outlet\_id \= *t*.outlet\_id    AND *mo*.cust\_id \= *t*.cust\_id LEFT JOIN mst.m\_salesman *ms*     ON *ms*.emp\_id \= *t*.salesman\_id    AND *ms*.cust\_id \= *t*.cust\_id    AND *ms*.is\_del \= FALSE LEFT JOIN mst.m\_employee *me*     ON *me*.emp\_id \= *ms*.emp\_id    AND *me*.cust\_id \= *ms*.cust\_id LEFT JOIN mst.m\_product *mp*     ON *mp*.pro\_id \= *t*.pro\_id    AND *mp*.cust\_id \= *t*.cust\_id JOIN mst.m\_supplier *ms2*     ON *ms2*.sup\_id \= *mp*.sup\_id ORDER BY *t*.*document\_date*, *t*.*document\_no*;  |
| :---- |

