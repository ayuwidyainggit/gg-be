1. # **Tabel mst.m\_distributor**

add field : parent\_cust\_id → varchar 10 , mandatory , relasi dengan smc.m\_customer.cust\_id

2. # **Distributor List** 

Enhance endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: master/v1/distributors

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **CURL** 

---

| curl 'https://best.scyllax.online/master/v1/distributors?page=1\&limit=10' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxIiwiZGlzdF9wcmljZV9ncnBfaWQiOjAsImRpc3RyaWJ1dG9yX2lkIjowLCJlbWFpbCI6InByaW5jQGlkZXRhbWEuaWQiLCJlbXBfaWQiOjI3OCwiZXhwaXJlcyI6MTc3NTE4NjU3NSwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4MTEzMjMyMzMyIiwicGFyZW50X2N1c3RfaWQiOiJDMjIwMDEiLCJ1c2VyX2Z1bGxuYW1lIjoiQWRtaW4gUHJpbmNpcGFsIDEiLCJ1c2VyX2lkIjoxLCJ1c2VyX25hbWUiOiJwcmluY0BpZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODExMzIzMjMzMiJ9.4xV\_MdkU9i\_oHKssTWjxpPJhWrTvuAg0HBNu74JRNdg' \\   \-H 'Connection: keep-alive' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"'  |
| :---- |

### **Perbaikan Query :** 

### ---

| SELECT \* FROM mst.m\_distributor *md* WHERE *md*.cust\_id LIKE (     SELECT *mu*.cust\_id || '%'     FROM sys.m\_user *mu*     WHERE *mu*.email \= 'princessa@gmail.com' ) AND *md*.is\_active \= true AND *md*.is\_del IS false; |
| :---- |

### 

### **1.1. User Login Sebagai Principal**

* Email : '[princessa@gmail.com](mailto:princessa@gmail.com)'  
* Menampilkan distributor bawahan dari user principal yang login   
  ![][image1]

**1.2. User Login Sebagai Distributor**

* Email : '[adminbm@gmail.com](mailto:adminbm@gmail.com)'  
* Menampilkan distributor milik dirinya sendiri   
  ![][image2]

3. # **Distributor Detail**  

Enhance endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: /master/v1/distributors/103

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token |

### **CURL** 

---

| curl 'https://best.scyllax.online/master/v1/distributors/103' \\   \-H 'sec-ch-ua-platform: "macOS"' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyIiwiZGlzdF9wcmljZV9ncnBfaWQiOjAsImRpc3RyaWJ1dG9yX2lkIjowLCJlbWFpbCI6InByaW5jZXNzYUBnbWFpbC5jb20iLCJlbXBfaWQiOjM4MCwiZXhwaXJlcyI6MTc3NTY2MDE4MiwiaXNfYWRtaW4iOmZhbHNlLCJsYW5nX2lkIjoiaWQiLCJtb2JpbGVfbm8iOiIwODEzMzMzMzMzMzMiLCJwYXJlbnRfY3VzdF9pZCI6IkMyNjAwMiIsInVzZXJfZnVsbG5hbWUiOiJQcmluY2Vzc2EgQWhzYW5pIFRhcXdpbSIsInVzZXJfaWQiOjE0MCwidXNlcl9uYW1lIjoiUHJpbmNlc3NhIEFoc2FuaSBUYXF3aW0iLCJ3aGF0c2FwcCI6IjA4MTMzMzMzMzMzMyJ9.NiVBrBTBEIae0PS7zJV8JcypuMHe8WE-QolDSgTbBjM' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'sec-ch-ua: "Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"' \\   \-H 'sec-ch-ua-mobile: ?0' |
| :---- |

### **Cek Issue** 

---

| {     "message": "nama tujuan parent\_cust\_id tidak ada di \*model.DistributorList",     "request\_id": "69d51b48942ab494e883a93c" } |
| :---- |
