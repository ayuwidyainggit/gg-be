# 

1. ## **Pseudocode** 

| START SET response.data \= data\_user IF user\_role \= "salesman distributor" THEN     IF salesman\_type \= "Taking order" THEN         IF plan \> 0 THEN             response.success \= true             response.message \= "Check-in Available"             response.description \= ""         ELSE             response.success \= false             response.message \= "Check-in Unavailable"             response.description \=                "Check-in cannot be completed because there is no scheduled route plan. Please contact your administrator."         ENDIF     ELSE IF salesman\_type \= "Taking order \+ canvas"           OR salesman\_type \= "Canvas" THEN         // Validasi PLAN & STOCK         IF plan \> 0 AND stock \> 0 THEN             response.success \= true             response.message \= "Check-in Available"             response.description \= ""         ELSE IF plan \= 0 AND stock \> 0 THEN             response.success \= false             response.message \= "Check-in Unavailable"             response.description \=                "Check-in cannot be completed because there is no scheduled route plan. Please contact your administrator."         ELSE IF plan \> 0 AND stock \= 0 THEN             response.success \= false             response.message \= "Check-in Unavailable"             response.description \=                "Check-in cannot be completed because the stock is unavailable. Please contact your administrator."         ELSE             response.success \= false             response.message \= "Check-in Unavailable"             response.description \=                "Check-in cannot be completed because there is no scheduled route plan and the stock is unavailable. Please contact your administrator."         ENDIF     ENDIF ELSE IF user\_role \= "principal" THEN     // Validasi PLAN saja     IF plan \> 0 THEN         response.success \= true         response.message \= "Check-in Available"         response.description \= ""     ELSE         response.success \= false         response.message \= "Check-in Unavailable"         response.description \=            "Check-in cannot be completed because there is no scheduled route plan. Please contact your administrator."     ENDIF ENDIF SET response.request\_id \= generateUUID() RETURN response END  |
| :---- |

## 

2. ## **Definisi user dan tipe salesman**

* user distributor  \= distributor\_id NOT NULL   
* user principal \= distributor\_id  NULL   
* User distributor taking order \= mst.m\_salesman.opr\_type \= O   
* User distributor taking canvas \= mst.m\_salesman\_canvas.is\_active \= true  &  mst.m\_salesman\_canvas.opr\_type \= C  
* User distributor taking canvas & taking\_order \=  
  mst.m\_salesman\_canvas.is\_active \= true  &  mst.m\_salesman\_canvas.opr\_type \= C &   
  mst.m\_salesman.opr\_type \= O 


3. ## **Query PLAN dan STOCK**

   1. ### **query plan \= is\_distributor FALSE**

   

| select COUNT(*rpp*.id) as *plan* from pjp\_principles.route\_pop\_permanent *rpp* join pjp\_principles.permanent\_journey\_plans *pjp* on *pjp*.id \=*rpp*.pjp\_id where  *pjp*.salesman\_id \= 345 AND *rpp*.date \= '2026-01-15' |
| :---- |

   

   2. ### **query plan \= is\_distributor TRUE** {#query-plan-=-is_distributor-true}

* Validasi 1 : ketika plan \= 0 

| select COUNT(*rpp*.id) as *plan* from pjp.route\_pop\_permanent *rpp* join pjp.permanent\_journey\_plans *pjp* on *pjp*.id \=*rpp*.pjp\_id where  *pjp*.salesman\_id \= 345 AND *rpp*.date \= '2026-01-15' |
| :---- |

## 

  3. ### **query stock  \= is\_distributor TRUE** {#query-stock-=-is_distributor-true}

* Validasi 2 : Ketika salesman canvas memiliki stock warehouse \= 0  

| select *msc*.wh\_id, COALESCE(*ws*.qty, 0) AS *qty*  from mst.m\_salesman *ms* left join mst.m\_salesman\_canvas *msc*  	on *msc*.emp\_id  \= *ms*.emp\_id left join inv.warehouse\_stock *ws*  	on *msc*.wh\_id \= *ws*.wh\_id where *ms*.emp\_id \= 213  |
| :---- |

## 

4. ## **API Validasi PJP**

 Create new endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: /mobile/v1/attendances/check

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer {token} |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| date | epoch | Yes | Unix timestamp (epoch time) example : 1719822823 |
| emp\_id  | Integer | Yes | diisi menggunakan emp\_id dari resp login |
| distributor\_id | Integer | No | User distributor \= NOT NULL  user principal \= NULL  |

**Example Request Default :** 

| curl \--location \-g \\   'https://be.scyllax.online/mobile/v1/attendances/check?date=1704067200\&emp\_id=12\&distributor\_id=99' \\   \--header 'Accept: application/json' \\   \--header 'Authorization: Bearer YOUR\_BEARER\_TOKEN' |
| :---- |

### **Response** 

---

| Nama Atribut | Type  | Length | Description |
| :---- | :---- | :---- | :---- |
| success | bool |  | if stock \= 0 OR plan \= 0 → false  |
| message | varchar | 255 |  \- |
| description | varchar | 255 | Id Stock Disposal  |
| data | Numeric | 11 | Total page keseluruhan |
| warehouse\_stock  | int | 8 | berdasarkan query [stock](#query-stock-=-is_distributor-true) |
| plan | int | 8 | berdasarkan query [plan](#query-plan-=-is_distributor-true) |

### **Response Principal** 

#### sucess

| {   "success": true,   "message": "Check-in Available",   “description: : “”   "data": {     "plan": 10   } "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |

#### error

| {   "success": false,   "message": "Check-in Unavailable",   “description: : “Check-in cannot be completed as there is no scheduled route plan. Please reach out to your administrator for further assistance.”   "data": {     "plan": 0   } "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |

### **Response Distributor**  

#### **sucess → (TO & C) or C**

| {   "success": true,   "message": "Check-in Available",   "description": "",   "data": {     "emp\_id": 213,     "emp\_code": "8989",     "emp\_name": "UCING \- ASM",     "opr\_type": "O",     "opr\_type\_canvas": "C",     "wh\_id": 67,     "wh\_code": "001",     "wh\_name\_canvas": "Gudang Warehouse HO",     "stock": 5,     "plan": 10   },   "request\_id": "6915a5e8e3f53f84fe73517f" }   |
| :---- |

#### **sucess → TO** 

| {   "success": true,   "message": "Check-in Available",   "description": "",   "data": {     "emp\_id": 213,     "emp\_code": "8989",     "emp\_name": "UCING \- ASM",     "opr\_type": "O",     "opr\_type\_canvas": "",     "wh\_id": ,     "wh\_code": "",     "wh\_name\_canvas": "",     "stock": “”,     "plan": 10   },   "request\_id": "6915a5e8e3f53f84fe73517f" }   |
| :---- |

#### **Error → (TO ) → plan 0** 

| {   "success": false,   "message": "Check-in Unavailable",   "description": "Check-in cannot be completed because there is no scheduled route plan. Please contact your administrator.",   "data": {     "emp\_id": 213,     "emp\_code": "8989",     "emp\_name": "UCING \- ASM",     "opr\_type": "O",     "opr\_type\_canvas": "C",     "wh\_id": 67,     "wh\_code": "001",     "wh\_name\_canvas": "Gudang Warehouse HO",     "stock": 10,     "plan": 0   },   "request\_id": "6915a5e8e3f53f84fe73517f" }   |
| :---- |

#### **Error (TO & C) or C →  plan \= 0 and stock \=0** 

| {   "success": false,   "message": "Check-in Unavailable",   "description": "Check-in cannot be completed because there is no scheduled route plan and the stock is unavailable. Please contact your administrator.",   "data": {     "emp\_id": 213,     "emp\_code": "8989",     "emp\_name": "UCING \- ASM",     "opr\_type": "O",     "opr\_type\_canvas": "C",     "wh\_id": 67,     "wh\_code": "001",     "wh\_name\_canvas": "Gudang Warehouse HO",     "stock": 0,     "plan": 0   },   "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |

#### **Error  (TO & C) or C→ plan \= 0** 

| {   "success": false,   "message": "Check-in Unavailable",   "description": "Check-in cannot be completed because there is no scheduled route plan. Please contact your administrator.",   "data": {     "emp\_id": 213,     "emp\_code": "8989",     "emp\_name": "UCING \- ASM",     "opr\_type": "O",     "opr\_type\_canvas": "C",     "wh\_id": 67,     "wh\_code": "001",     "wh\_name\_canvas": "Gudang Warehouse HO",     "stock": 10,     "plan": 0   },   "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |

#### **Error  (TO & C) or C → stock \= 0** 

| {   "success": false,   "message": "Check-in Unavailable",   "description": "Check-in cannot be completed because the stock is unavailable. Please contact your administrator.",   "data": {     "emp\_id": 213,     "emp\_code": "8989",     "emp\_name": "UCING \- ASM",     "opr\_type": "O",     "opr\_type\_canvas": "C",     "wh\_id": 67,     "wh\_code": "001",     "wh\_name\_canvas": "Gudang Warehouse HO",     "stock": 10,     "plan": 0   },   "request\_id": "6915a5e8e3f53f84fe73517f" }  |
| :---- |

**validasi : (distributor)** 

* Salesman to \= validasi pjp   
* Taking order \+ canvas \= validasi pjp stock → lanjut canvas dan TO   
* Canvas \= validasi pjp \+ stock 

**validasi : (principal) \= pjp**  