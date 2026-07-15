1. # **Sales Summary** 

   URL : [http://103.28.219.73:5001/mobile/v1/sales/summary](http://103.28.219.73:5001/mobile/v1/sales/summary)  
   Method : GET   
   Header :   
   {  
    "content-type": "application/json",  
    "authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbGxvd19pbnB1dF9wcmljZSI6ZmFsc2UsImN1c3RfaWQiOiJDMjIwMDEwMDAxIiwiZGlzdF9wcmljZV9ncnBfaWQiOjEsImVtYWlsIjoiY2hhcmxlc0BnbWFpbC5jb20iLCJlbXBfY29kZSI6IjMwMTIyNDEiLCJlbXBfZ3JwX2lkIjo4OSwiZW1wX2lkIjoyMjgsImV4cGlyZXMiOjE3NjUyNTQ0MTAsImlzX2FjdGl2ZV9ndWRhbmdfY2FudmFzIjp0cnVlLCJpc19hY3RpdmVfZ3VkYW5nX3V0YW1hIjpmYWxzZSwibGFuZ19pZCI6ImlkIiwibW9iaWxlX25vIjoiMDgxMzMzMzMzMzMzIiwib3ByX3R5cGVfY2FudmFzIjoiQyIsIm9wcl90eXBlX29yZGVyX3Rha2luZyI6Ik8iLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInRheF9vcHRpb24iOiJFIiwidXNlcl9pZCI6MTIwLCJ1c2VyX3JvbGUiOiJzYWxlc21hbiIsIndoYXRzYXBwIjoiMDgxMzMzMzMzMzMzIn0.n7M8rRks6IDnu7B4D-ViM6XI-JUL0\_cfvuASQIM1iio"  
   }  
     
   Request Body :   
   N/A  
     
   Response :   
    "message": "",  
    "data": {  
      "current\_sales": 0,  
      "daily\_targe": 0  
    },  
    "request\_id": "69378795df3c349a68774722"  
   }  
     
   Enhance :   
* Current Sales

| Before Enhancement | After Enhancement |
| :---- | :---- |
| **COALESCE(SUM(total), 0\)** AS current\_sales, FROM sls.order WHERE salesman\_id \= $1 AND DATE(ro\_date) \= **CURRENT\_DATE** \-- Filter hanya penjualan hari ini AND is\_del \= false; \-- Pastikan order tidak dihapus | **current\_sales** sebelumnya hanya menghitung total order berdasarkan CURRENT\_DATE. Pada enhancement, perhitungan diubah menjadi **total order \- retur** dan memakai rentang tanggal dinamis: **start\_date** \= tanggal pertama pada bulan berjalan **end\_date** \= tanggal terbaru dari mobile.attendance.create\_at dengan type \= 1 (clock in)   Total Order :  SELECT **ro\_date, total::bigint, ro\_no, invoice\_no** FROM sls.order  WHERE cust\_id \= 'C220010001' AND salesman\_id \= 210  AND ro\_date BETWEEN '2025-12-01' AND '2025-12-30' AND is\_del \= false and invoice\_no is not NULL ORDER BY ro\_date; Return :  SELECT **return\_date, total::bigint, return\_no** FROM sls.return  WHERE invoice\_no in (SELECT invoice\_no FROM sls.order  WHERE cust\_id \= 'C220010001' AND salesman\_id \= 210  AND ro\_date BETWEEN '2025-12-01' AND '2025-12-24' AND is\_del \= false and invoice\_no is not NULL ORDER BY ro\_date) Total Current Sales \= Total Order \-  Return |

### 	**Enhance 13 Jan 2026**  **Response :** 

| {  "message": "",  "data": {    "current\_sales": 0,    "daily\_target": 0,    "monthly\_sales\_target": 0,  },  "request\_id": "69378795df3c349a68774722" }  |
| :---- |

Add response **monthly\_sales\_target**

data source : (ambil mst.allocated\_total, jika hasil ada lebih dr 1, SUM hasil tersebut) 

| select *mst*.sales\_target\_id , *mst*."month" , *mst*."year" , *mst*.allocated\_total ,   *msa*.salesman\_id , *msa*.allocated  from mst.m\_sales\_target *mst*  join mst.m\_sales\_allocated *msa*  	on *msa*.sales\_team\_id \= *mst*.sales\_target\_id  where *mst*.month=2 and *mst*.year=2025 and *msa*.salesman\_id \=2 and *msa*.is\_del \= false  and *msa*.is\_active \=true and *mst*.status \=1 and mst.cust\_id=’C220010001’  |
| :---- |

2. # **Summary Daily**

   URL : [http://103.28.219.73:5001/mobile/v1/activities/summary/daily](http://103.28.219.73:5001/mobile/v1/activities/summary/daily)   
   Method : GET   
   Header :   
   {  
    "content-type": "application/json",  
    "authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhbGxvd19pbnB1dF9wcmljZSI6ZmFsc2UsImN1c3RfaWQiOiJDMjIwMDEwMDAxIiwiZGlzdF9wcmljZV9ncnBfaWQiOjEsImVtYWlsIjoiY2hhcmxlc0BnbWFpbC5jb20iLCJlbXBfY29kZSI6IjMwMTIyNDEiLCJlbXBfZ3JwX2lkIjo4OSwiZW1wX2lkIjoyMjgsImV4cGlyZXMiOjE3NjUyNTQ0MTAsImlzX2FjdGl2ZV9ndWRhbmdfY2FudmFzIjp0cnVlLCJpc19hY3RpdmVfZ3VkYW5nX3V0YW1hIjpmYWxzZSwibGFuZ19pZCI6ImlkIiwibW9iaWxlX25vIjoiMDgxMzMzMzMzMzMzIiwib3ByX3R5cGVfY2FudmFzIjoiQyIsIm9wcl90eXBlX29yZGVyX3Rha2luZyI6Ik8iLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInRheF9vcHRpb24iOiJFIiwidXNlcl9pZCI6MTIwLCJ1c2VyX3JvbGUiOiJzYWxlc21hbiIsIndoYXRzYXBwIjoiMDgxMzMzMzMzMzMzIn0.n7M8rRks6IDnu7B4D-ViM6XI-JUL0\_cfvuASQIM1iio"  
   }  
     
   Response :   
   {  
    "message": "",  
    "data": {  
      "last\_update": "",  
      "plan": 127,  
      "visit": 41,  
      "effective\_call": 1291,  
      "start\_time": "",  
      "end\_time": "",  
      "drive\_time": "",  
      "est\_time": ""  
    },  
    "request\_id": "69378795df3c349a6877471f"  
   }  
     
   Enhance :   
* Current Sales  
- Simulasi : [link](https://docs.google.com/spreadsheets/d/10MNiRK2VpYJbpje8n0koXz_QRjZdVsSg2rsFUD0HDyg/edit?usp=sharing)  
- Contoh di hari ini ada total 10 outlet yang harus dikunjungi sama sales.   
- Dari 10, yg sudah dikunjungi ada 3\.   
- Dan dari 3 yang sudah dikunjungi ada 2 yang belanja. Maka datanya akan:  
  **Plan** 10  
  **Visit**  3  
  **Effective Call** 2

| Field | Before Enhancement | After Enhancement |
| :---- | :---- | :---- |
| Plan | *Total kunjungan (Outlet) di hari itu* |  |
|  | SELECT COUNT(ovl.\*) FROM pjp.outlet\_visit\_list AS ovl JOIN pjp.permanent\_journey\_plans AS perpjp ON perpjp.pjp\_code \= ovl.pjp\_code WHERE perpjp.salesman\_id \= $1 AND ovl.arrive\_at IS NULL AND ovl.skip\_at IS NULL AND ovl.leave\_at IS NULL ) AS plan, | SELECT COUNT(ovl.\*)     FROM pjp.outlet\_visit\_list AS ovl     JOIN pjp.permanent\_journey\_plans AS perpjp  ON perpjp.pjp\_code \= ovl.pjp\_code     WHERE perpjp.salesman\_id \= $1 AND ovl.arrive\_at IS NULL AND ovl.date \= "2025-09-29" ) AS plan,   **ovl.date \= current date **  |
| Visit |  *yg berhasil dikunjungi* |  |
|  | SELECT COUNT(ovl.\*) FROM pjp.outlet\_visit\_list AS ovl JOIN pjp.permanent\_journey\_plans AS perpjp ON perpjp.pjp\_code \= ovl.pjp\_code WHERE perpjp.salesman\_id \= $1 AND ovl.arrive\_at IS NOT NULL AND ovl.skip\_at IS NULL AND ovl.leave\_at IS NULL ) AS visit,  | SELECT COUNT(ovl.\*)     FROM pjp.outlet\_visit\_list AS ovl     JOIN pjp.permanent\_journey\_plans AS perpjp  ON perpjp.pjp\_code \= ovl.pjp\_code     WHERE perpjp.salesman\_id \= $1 AND ovl.arrive\_at IS NOT NULL AND ovl.date \= "2025-09-29" ) AS visit,  **ovl.date \= current date**  |
| Effective Call | *Total outlet yang dikunjungi \+ order* |  |
|  | SELECT COUNT(orders.\*) FROM sls.order AS orders JOIN pjp.outlet\_visit\_list AS ovl ON ovl.outlet\_id \= orders.outlet\_id WHERE orders.salesman\_id \= $1 \-- Biasanya ditambahkan filter tanggal hari ini agar sesuai konteks "Daily" \-- AND DATE(orders.ro\_date) \= CURRENT\_DATE ) AS effective\_call;  | SELECT COUNT(orders.\*) FROM sls.order AS orders JOIN pjp.outlet\_visit\_list AS ovl ON ovl.outlet\_id \= orders.outlet\_id WHERE orders.salesman\_id \= $1 AND orders.ro\_date \= "2025-09-29" ) AS effective\_call;   **orders.ro\_date \= current date**  |
| start\_time  | NULL | mobile.attendance where type \= 1 and created\_date \= current\_date  |
| end\_time  | NULL  | mobile.attendance where type \= 2 and created\_date \= current\_date  |


  

  