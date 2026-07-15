## 

# **API Documentation**

**Scylla X Pratesis**

# 

| Prepared forScylla X Pratesis Prepared by GeekGarden Software House | Version	: 1.0Date		: 14 Nov 2025 |
| :---- | :---- |

# 

# **Document Version History**

| Version | Date | Description | Author |
| :---- | :---- | :---- | :---- |
| 1.0 | 14 Nov 2025 | Initialization | Ayu Widya Inggit |
|  |  |  |  |

# **Document Approval**

| Name | Role | Organization | Signature | Date |
| :---- | :---- | :---- | :---- | :---- |
|  |  |  |  |  |
|  |  |  |  |  |
|  |  |  |  |  |

# 

1. ## **Appendix** 

Status   
![][image1]

| No  | Status  | Description |
| :---- | :---- | :---- |
|  | 1 | Need Review |
|  | 2 | Processed |
|  | 3 | On Delivery |
|  | 4 | Received |
|  | 5 | Partial Received |
|  | 6 | Invoicing |
|  | 7 | Completed |
|  | 9 | Cancelled |

[simulasi](https://docs.google.com/spreadsheets/d/1DAhmaxMPIdp6pYuaPRi9Gz93ABhiUulLzksRpUuqiuU/edit?usp=sharing) 

2. ## **Sales Order List** 

URL Eksisting : [link](https://docs.google.com/document/d/1uLiXs-ansJHSWPDeTLEp5yuyFGtjzjjV/edit#heading=h.lutg6whvnds4)

| https://be.scyllax.online/sales/v1/orders?page=1\&limit=20 |
| :---- |

Method : GET   
**Enhance Response :** 

| Field  | Type | Desc |
| :---- | :---- | :---- |
| outlet\_address1 | varchar 150 | alamat outlet  |
| source | enum(mobile, web) | Add new field at sls.order Digunakan untuk logic tab di detail order |

**Payload Default :**   
![][image2]

**Payload with filter start date and end date :**   
![][image3]  
**Payload with filter start date, end date, salesman, outlet, status:**   
![][image4]

**Response Sukses → case mobile status : Need Review:** 

| {     "message": "",     "data": \[         {             "ro\_no": "SO2509230002",            "source" : "mobile",             "ro\_date": "2025-09-23",             "val\_date": null,             "salesman\_id": 299,             "sales\_name": "Jhon Takbor",             "wh\_id": 254,             "wh\_code": "001",             "wh\_name": "gudang baik",             "outlet\_id": 1445,             "outlet\_code": "A0001",             "outlet\_name": "TK Agus",            "outlet\_address1": "Jl Proklamasi",             "delivery\_date": "2025-09-23",             "order\_no": null,             "po\_no": null,             "vehicle\_no": null,             "pay\_type": 1,             "pay\_type\_name": "Cash On Delivery",             "reff\_no": null,             "mobile\_id": 1,             "sub\_total": 408000,             "disc": 0,             "disc\_value": 0,             "promo\_value": 0,             "cash\_disc\_value": 0,             "tot\_disc1": 0,             "tot\_disc2": 0,             "vat": 11,             "vat\_value": 44880,             "total": 452880,             "data\_status": 6,             "data\_status\_name": "Invoicing",             "updated\_at": "2025-11-20T04:19:54.052613Z",             "updated\_by\_name": "Admin DT Bangka 01",             "due\_date": "2025-11-20T00:00:00Z",             "tr\_code": null,             "is\_closed": false,             "notes": "",             "invoice\_no": "INV2511200001",             "invoice\_date": "2025-11-20"         },         {             "ro\_no": "SO2509030002",             "source" : "mobile",             "ro\_date": "2025-09-03",             "val\_date": null,             "salesman\_id": 299,             "sales\_name": "Jhon Takbor",             "wh\_id": 254,             "wh\_code": "001",             "wh\_name": "gudang baik",             "outlet\_id": 1445,             "outlet\_code": "A0001",             "outlet\_name": "TK Agus",             "outlet\_address1": "Jl Proklamasi",             "delivery\_date": "2025-09-03",             "order\_no": null,             "po\_no": null,             "vehicle\_no": null,             "pay\_type": 1,             "pay\_type\_name": "Cash On Delivery",             "reff\_no": null,             "mobile\_id": 1,             "sub\_total": 584256,             "disc": 0,             "disc\_value": 0,             "promo\_value": 0,             "cash\_disc\_value": 0,             "tot\_disc1": 0,             "tot\_disc2": 0,             "vat": 11,             "vat\_value": 64268,             "total": 648524,             "data\_status": 6,             "data\_status\_name": "Invoicing",             "updated\_at": "2025-09-03T01:50:23.448612Z",             "updated\_by\_name": "Subardan Distributor",             "due\_date": "2025-09-18T00:00:00Z",             "tr\_code": null,             "is\_closed": false,             "notes": "",             "invoice\_no": "INV2509030002",             "invoice\_date": "2025-09-03"         }     \],     "paging": {         "total\_record": 21,         "page\_current": 1,         "page\_limit": 20,         "page\_total": 2     },     "request\_id": "8e8e5c7f-6ae0-46ef-af05-ac496d67812c" }  |
| :---- |

**Response Negatif : (apabila no response)** 

| {     "message": "",     "data": null,     "paging": {         "total\_record": 0,         "page\_current": 1,         "page\_limit": 20,         "page\_total": 0     },     "request\_id": "a08e5c7f-6ae0-46ef-af05-ac496d67812c" }  |
| :---- |

3. ## **Sales Order Detail** 

**URL** : 

| https://be.scyllax.online/sales/v2/orders/SO2509230002 |
| :---- |

**Method** : GET 

**URL Promo :** 

| /sales/v2/promotions/consult |
| :---- |

Contoh req dan response : click [here](https://docs.google.com/spreadsheets/d/1UGsfcV0-Lwhi9rv6cNOdkXZwCOF3mo5Ei-BNhbowPYE/edit?usp=sharing)  
**Enhance** : tambahkan field pada response dengan text bold di bawah ini : 

| Field  | source |
| :---- | :---- |
| **opr\_type** | sls.order.opr\_type |
| **source** | order.data\_source  1: web 2: mobile  |
| **is\_performa\_invs** | // null karena masih menunggu enhance performa invoice  |
| **purchase\_details** |  **lihat tabel purchase\_detail di bawah ini :**  |

**purchase\_details**

| purchase\_details | source  |
| :---- | :---- |
| **"order\_detail\_id": 2657,**                | sls.order\_detail.order\_detail\_id |
| **"seq\_no": 0,**              | sls.order\_detail.seq\_no |
| **"pro\_id": 792,** | sls.order\_detail.pro\_id |
| **"pro\_code": "RJ0001",**  | mst.m\_product.pro\_code |
| **"pro\_name": "REJ SHP KOREAN LAVENDER 350MLX24",**  | mst.m\_product.pro\_name |
| **"order\_status": "",** | **sls.order.data\_status** |
| **"item\_type": 1,** | sls.order\_detail.item\_type |
| **"qty": 24,** | sls.order\_detail.qty |
| **"qty\_final": 24,** | sls.order\_detail.qty\_final |
| **"qty\_po": 0,** | sls.order\_detail.qty\_po |
| **"qty\_po1": 0,** | sls.order\_detail.qty\_po1 |
| **"qty\_po2": 0,** | sls.order\_detail.qty\_po2 |
| **"qty\_po3": 1,** | sls.order\_detail.qty\_po3 |
| **"qty1": 0,** | sls.order\_detail.qty1 |
| **"qty2": 0,** | sls.order\_detail.qty2 |
| **"qty3": 1,** | sls.order\_detail.qty3 |
| **"qty4": null,** | sls.order\_detail.qty4 |
| **"qty5": null,** | sls.order\_detail.qty5 |
| **"qty1\_final": 0,** | sls.order\_detail.qty1\_final |
| **"qty2\_final": 0,** | sls.order\_detail.qty2\_final |
| **"qty3\_final": 1,** | sls.order\_detail.qty3\_final |
| **"qty4\_final": null,** | sls.order\_detail.qty4\_final |
| **"qty5\_final": null,** | sls.order\_detail.qty5\_final |
| **"qty1\_stok": 0,** | **Before :**  sls.order\_detail.qty1\_stok sls.order\_detail.qty2\_stok sls.order\_detail.qty3\_stok **After :**  Konversi **inv.warehouse\_stock.qty**  menjadi S, M, L ***(code eksisting)*** masing” di tambah qty1(L) , qty2(M) , qty3 (S) |
| **"qty2\_stok": 0,** |  |
| **"qty3\_stok": 10,** |  |
| "promo1": 0,  | promo1 |
| "promo2": 0,  | promo2 |
| "promo3": 0,  | promo3 |
| "promo4": 0,  | promo4 |
| "promo5": 0,  | promo5 |
| **"purch\_price1": 16000,** | sls.order\_detail.purch\_price1 |
| **"purch\_price2": 96000,** | sls.order\_detail.purch\_price2 |
| **"purch\_price3": 384000,** | sls.order\_detail.purch\_price3 |
| **"purch\_price4": 0,** | sls.order\_detail.purch\_price4 |
| **"purch\_price5": 0,** | sls.order\_detail.purch\_price5 |
| **"sell\_price1": 17000,** | sls.order\_detail.sell\_price\_po1 |
| **"sell\_price2": 102000,** | sls.order\_detail.sell\_price\_po1 |
| **"sell\_price3": 408000,** | sls.order\_detail.sell\_price\_po1 |
| **"sell\_price4": 0,** | sls.order\_detail.sell\_price4 |
| **"sell\_price5": 0** | sls.order\_detail.sell\_price5 |
| **"sell\_price\_system1": 17000,** | sls.order\_detail.sell\_price\_system1 |
| **"sell\_price\_system2": 102000,** | sls.order\_detail.sell\_price\_system2 |
| **"sell\_price\_system3": 408000,** | sls.order\_detail.sell\_price\_system3 |
| **"sell\_price\_system4": 0,** | sls.order\_detail.sell\_price\_system4 |
| **"sell\_price\_system5": 0,** | sls.order\_detail.sell\_price\_system5 |
| **"amount": 452880,** | sls.order\_detail.amount |
| **"amount\_final": 452880,** | sls.order\_detail.amount\_final |
| **"promo\_value": 0,** | sls.order\_detail.promo\_value |
| **"promo\_value\_final": 0,** | sls.order\_detail.promo\_value\_final |
| **"disc\_value": 0,** | sls.order\_detail.disc\_value |
| **"disc\_value\_final": 0,** | sls.order\_detail.disc\_value\_final |
| **"batch\_no": null,** | sls.order\_detail.batch\_no |
| **"exp\_date": null,** | sls.order\_detail.exp\_date |
| **"vat": 11,** | sls.order\_detail.vat |
| **"vat\_bg": 0,** | sls.order\_detail.bg |
| **"vat\_lg\_sell": 0,** | sls.order\_detail.vat.lg.sell |
| **"vat\_value": 44880,** | sls.order\_detail.vat\_value |
| **"vat\_value\_final": 44880,** | sls.order\_detail.vat\_value\_final |
| **"vat\_bg\_value": 0,** | sls.order\_detail.vat\_bg\_value |
| **"vat\_lg\_value": null,** | **tdk ada** |
| **"vat\_lg\_sell\_value": 0,** | sls.order\_detail.vat\_lg\_sell\_value |
| **"unit\_id1": "BTL",** | sls.order\_detail.unit\_ud1 |
| **"unit\_id2": "BOX",** | sls.order\_detail.unit\_ud2 |
| **"unit\_id3": "KRT",** | sls.order\_detail.unit\_ud3 |
| **"unit\_id4": "",** | sls.order\_detail.unit\_ud4 |
| **"unit\_id5": "",** | sls.order\_detail.unit\_ud5 |
| **"conv\_unit2": 6,** | sls.order\_detail.conv\_unit2 |
| **"conv\_unit3": 4,** | sls.order\_detail.conv\_unit3 |
| **"conv\_unit4": 0,** | sls.order\_detail.conv\_unit4 |
| **"conv\_unit5": 0,** | sls.order\_detail.conv\_unit5 |
| **"notes": "",** | sls.order\_detail.notes |
| **"discount\_id": "",** | sls.order\_detail.discount\_id |
| **"promoa\_id": null** | **tdk ada**  |

**Update Eksisting  (details\_final) \--\> add response**

| purchase\_details | source  |
| :---- | :---- |
| **sell\_price\_final1** | sls.order\_detail.sell\_price\_final1 |
| **sell\_price\_final2** | sls.order\_detail.sell\_price\_final2 |
| **sell\_price\_final3** | sls.order\_detail.sell\_price\_final3 |

**Response** : 

| {     "message": "",     "data": {         "ro\_no": "SO2602210001",         "source": "web",         "is\_proforma\_inv": false,         "order\_no": null,         "ro\_date": "2026-02-21",         "val\_date": null,         "salesman\_id": 357,         "salesman\_code": "2025120201",         "sales\_name": "Ilham Abdullah",         "wh\_id": 63,         "wh\_code": "008",         "wh\_name": "Gudang Utama",         "outlet\_id": 404,         "outlet\_code": "P00025Ar",         "outlet\_name": "Princes 9",         "outlet\_address1": "PERUM CANDRALOKA BLOK AA.12 NO.9",         "outlet\_address2": "",         "inv\_addr1": "PERUM CANDRALOKA BLOK AA.12 NO.4",         "inv\_addr2": "",         "delivery\_date": "2026-02-21",         "po\_no": null,         "vehicle\_no": null,         "pay\_type": 3,         "pay\_type\_name": "Credit",         "reff\_no": null,         "mobile\_id": 1,         "sub\_total": 310000,         "sub\_total\_final": 310000,         "disc": 0,         "disc\_value": 0,         "disc\_value\_final": 0,         "promo\_value": 0,         "promo\_value\_final": 0,         "promo\_bg\_value": 0,         "promo\_bg\_value\_final": 0,         "cash\_disc\_value": 0,         "tot\_disc1": 0,         "tot\_disc2": 0,         "vat": 11,         "vat\_value": 37200,         "vat\_value\_final": 37200,         "total": 347200,         "total\_final": 347200,         "data\_status": 2,         "data\_status\_name": "Processed",         "data\_source": 1,         "updated\_at": "2026-02-21T08:03:03.471215Z",         "updated\_by\_name": "Dist IDE Sda",         "due\_date": null,         "details": {             "normal": \[                 {                     "order\_detail\_id": 5215,                     "seq\_no": 0,                     "pro\_id": 688,                     "pro\_code": "FF-0101-0040",                     "pro\_name": "Cavendish Hand",                     "order\_status": "",                     "item\_type": 1,                     "qty": 12,                     "qty\_final": 12,                     "qty\_po": 12,                     "qty\_po1": null,                     "qty\_po2": null,                     "qty\_po3": null,                     "qty1": 0,                     "qty2": 0,                     "qty3": 1,                     "qty4": null,                     "qty5": null,                     "qty1\_final": 0,                     "qty2\_final": 0,                     "qty3\_final": 1,                     "qty4\_final": null,                     "qty5\_final": null,                     "qty1\_stok": 0,                     "qty2\_stok": 0,                     "qty3\_stok": 59,                     "purch\_price1": 10000,                     "purch\_price2": 200000,                     "purch\_price3": 300000,                     "purch\_price4": 0,                     "purch\_price5": 0,                     "sell\_price1": 110000,                     "sell\_price2": 210000,                     "sell\_price3": 310000,                     "sell\_price4": 0,                     "sell\_price5": 0,                     "sell\_price\_po1": null,                     "sell\_price\_po2": null,                     "sell\_price\_po3": null,                     "sell\_price\_final1": 110000,                     "sell\_price\_final2": 210000,                     "sell\_price\_final3": 310000,                     "sell\_price\_system1": 110000,                     "sell\_price\_system2": 110000,                     "sell\_price\_system3": 110000,                     "sell\_price\_system4": null,                     "sell\_price\_system5": null,                     "amount": 347200,                     "amount\_final": 347200,                     "promo\_value": 0,                     "promo\_value\_final": 0,                     "disc\_value": 0,                     "disc\_value\_final": 0,                     "batch\_no": null,                     "exp\_date": null,                     "vat": 12,                     "vat\_bg": 0,                     "vat\_lg\_sell": 0,                     "vat\_value": 37200,                     "vat\_value\_final": 37200,                     "vat\_bg\_value": 0,                     "vat\_lg\_value": null,                     "vat\_lg\_sell\_value": 0,                     "disc\_po": 0,                     "vat\_value\_po": 37200,                     "unit\_id1": "PNCH",                     "unit\_id2": "CRT",                     "unit\_id3": "CRT",                     "unit\_id4": "",                     "unit\_id5": "",                     "conv\_unit2": 12,                     "conv\_unit3": 1,                     "conv\_unit4": 0,                     "conv\_unit5": 0,                     "notes": "",                     "discount\_id": "",                     "promoa\_id": null                 }             \],             "promo": null         },         "details\_final": {             "normal": \[                 {                     "order\_detail\_id": 5215,                     "seq\_no": 0,                     "pro\_id": 688,                     "pro\_code": "FF-0101-0040",                     "pro\_name": "Cavendish Hand",                     "order\_status": "",                     "item\_type": 1,                     "qty": 12,                     "qty\_final": 12,                     "qty\_po": 12,                     "qty\_po1": null,                     "qty\_po2": null,                     "qty\_po3": null,                     "qty1": 0,                     "qty2": 0,                     "qty3": 1,                     "qty4": null,                     "qty5": null,                     "qty1\_final": 0,                     "qty2\_final": 0,                     "qty3\_final": 1,                     "qty4\_final": null,                     "qty5\_final": null,                     "qty1\_stok": 0,                     "qty2\_stok": 0,                     "qty3\_stok": 59,                     "purch\_price1": 10000,                     "purch\_price2": 200000,                     "purch\_price3": 300000,                     "purch\_price4": 0,                     "purch\_price5": 0,                     "sell\_price1": 110000,                     "sell\_price2": 210000,                     "sell\_price3": 310000,                     "sell\_price4": 0,                     "sell\_price5": 0,                     "sell\_price\_po1": null,                     "sell\_price\_po2": null,                     "sell\_price\_po3": null,                     "sell\_price\_final1": 110000,                     "sell\_price\_final2": 210000,                     "sell\_price\_final3": 310000,                     "sell\_price\_system1": 110000,                     "sell\_price\_system2": 110000,                     "sell\_price\_system3": 110000,                     "sell\_price\_system4": null,                     "sell\_price\_system5": null,                     "amount": 347200,                     "amount\_final": 347200,                     "promo\_value": 0,                     "promo\_value\_final": 0,                     "disc\_value": 0,                     "disc\_value\_final": 0,                     "batch\_no": null,                     "exp\_date": null,                     "vat": 12,                     "vat\_bg": 0,                     "vat\_lg\_sell": 0,                     "vat\_value": 37200,                     "vat\_value\_final": 37200,                     "vat\_bg\_value": 0,                     "vat\_lg\_value": null,                     "vat\_lg\_sell\_value": 0,                     "disc\_po": 0,                     "vat\_value\_po": 37200,                     "unit\_id1": "PNCH",                     "unit\_id2": "CRT",                     "unit\_id3": "CRT",                     "unit\_id4": "",                     "unit\_id5": "",                     "conv\_unit2": 12,                     "conv\_unit3": 1,                     "conv\_unit4": 0,                     "conv\_unit5": 0,                     "notes": "",                     "discount\_id": "",                     "promoa\_id": null                 }             \],             "promo": null         },         "purchase\_details": {             "normal": \[                 {                     "order\_detail\_id": 5215,                     "seq\_no": 0,                     "pro\_id": 688,                     "pro\_code": "FF-0101-0040",                     "pro\_name": "Cavendish Hand",                     "order\_status": "2",                     "item\_type": 1,                     "qty": 12,                     "qty\_final": 12,                     "qty\_po": 12,                     "qty\_po1": null,                     "qty\_po2": null,                     "qty\_po3": null,                     "qty1": 0,                     "qty2": 0,                     "qty3": 1,                     "qty4": null,                     "qty5": null,                     "qty1\_final": 0,                     "qty2\_final": 0,                     "qty3\_final": 1,                     "qty4\_final": null,                     "qty5\_final": null,                     "qty1\_stok": 0,                     "qty2\_stok": 0,                     "qty3\_stok": 59,                     "purch\_price1": 10000,                     "purch\_price2": 200000,                     "purch\_price3": 300000,                     "purch\_price4": 0,                     "purch\_price5": 0,                     "sell\_price1": 110000,                     "sell\_price2": 210000,                     "sell\_price3": 310000,                     "sell\_price4": 0,                     "sell\_price5": 0,                     "sell\_price\_po1": 110000,                     "sell\_price\_po2": 210000,                     "sell\_price\_po3": 310000,                     "sell\_price\_final1": 110000,                     "sell\_price\_final2": 210000,                     "sell\_price\_final3": 310000,                     "sell\_price\_system1": 110000,                     "sell\_price\_system2": 110000,                     "sell\_price\_system3": 110000,                     "sell\_price\_system4": null,                     "sell\_price\_system5": null,                     "amount": 347200,                     "amount\_final": 347200,                     "promo\_value": 0,                     "promo\_value\_final": 0,                     "disc\_value": 0,                     "disc\_value\_final": 0,                     "batch\_no": null,                     "exp\_date": null,                     "vat": 12,                     "vat\_bg": 0,                     "vat\_lg\_sell": 0,                     "vat\_value": 0,                     "vat\_value\_final": 0,                     "vat\_bg\_value": 0,                     "vat\_lg\_value": null,                     "vat\_lg\_sell\_value": 0,                     "disc\_po": 0,                     "vat\_value\_po": 37200,                     "unit\_id1": "PNCH",                     "unit\_id2": "CRT",                     "unit\_id3": "CRT",                     "unit\_id4": "",                     "unit\_id5": "",                     "conv\_unit2": 12,                     "conv\_unit3": 1,                     "conv\_unit4": 0,                     "conv\_unit5": 0,                     "notes": "",                     "discount\_id": "",                     "promoa\_id": null                 }             \],             "promo": \[\]         },         "remarks": null,         "tr\_code": null,         "is\_closed": false,         "notes": "",         "invoice\_no": null,         "invoice\_date": null,         "is\_printed": false,         "printed\_by": null,         "printed\_by\_name": null,         "printed\_at": null,         "validate\_stok": true,         "validate\_stok\_message": "Sufficient Stock",         "validate\_credit\_limit": true,         "validate\_credit\_limit\_message": "Within Limit",         "validate\_overdue": true,         "validate\_overdue\_message": "Allowed (Unlimited)",         "validate\_outstanding": true,         "validate\_outstanding\_message": "Allowed (Unlimited)",         "validate\_summary": true,         "credit\_limit\_type": 2,         "credit\_limit\_action": 1,         "credit\_limit\_action\_name": "Warning",         "sales\_inv\_limit\_type": null,         "sales\_inv\_limit\_action": null,         "sales\_inv\_limit\_action\_name": "",         "obs\_type": null,         "obs\_limit\_action": null,         "obs\_limit\_action\_name": "",         "order\_approval\_request\_id": null,         "order\_approval\_request\_emp\_approval\_status": null     },     "request\_id": "dea1bf04-3920-426a-90ff-f42040486208" }  |
| :---- |

## 

* ### Sales Order Detail 

4. ## **Sales Order Detail include Promo**

Dokumen Promo : [https://docs.google.com/spreadsheets/d/1UGsfcV0-Lwhi9rv6cNOdkXZwCOF3mo5Ei-BNhbowPYE/edit?gid=0\#gid=0](https://docs.google.com/spreadsheets/d/1UGsfcV0-Lwhi9rv6cNOdkXZwCOF3mo5Ei-BNhbowPYE/edit?gid=0#gid=0)   
23 Feb 2025   
(Update value PROMO)  
→ Melanjutkan API detail sales order  
Sebelumnya pada masing-masing tab dan masing-masing produk belum ada response promo   
Enhance  : Tambahkan response pada details.normal, details\_final.normal, purchase\_details.normal pada masing-masing tab dan masing-masing produk.

### **Sequence**

![][image5]

| sequenceDiagram     autonumber     participant FE as Frontend     participant BE as Sales Order Service     participant DB as PostgreSQL     participant PROMO as Promotion Service     FE-\>\>BE: GET /sales/v2/orders/{ro\_no}     activate BE     BE-\>\>DB: Fetch order \+ details.normal \+ details\_final \+ purchase\_details     activate DB     DB--\>\>BE: Order aggregate data     deactivate DB     BE-\>\>BE: Generate hash/signature per tab     BE-\>\>BE: Compare normal vs final vs purchase     alt All Tabs Identical         BE-\>\>PROMO: POST /promotions/consult (single payload)         activate PROMO         PROMO--\>\>BE: Promo result         deactivate PROMO         BE-\>\>BE: Inject same promo to all tabs     else Tabs Different         par Hit Promo Normal             BE-\>\>PROMO: POST consult (details.normal)             PROMO--\>\>BE: Promo result (normal)         and Hit Promo Final             BE-\>\>PROMO: POST consult (details\_final.normal)             PROMO--\>\>BE: Promo result (final)         and Hit Promo Purchase             BE-\>\>PROMO: POST consult (purchase\_details.normal)             PROMO--\>\>BE: Promo result (purchase)         end         BE-\>\>BE: Inject promo per respective tab     end     BE--\>\>FE: Return order response (with promo per tab)     deactivate BE  |
| :---- |

### **Pseudocode**

| FUNCTION getSalesOrder(ro\_no):     \# \==========================================     \# 1\. Fetch Order Aggregate     \# \==========================================     order \= DB.fetchOrderWithAllDetails(ro\_no)     normalTab    \= order.details.normal     finalTab     \= order.details\_final.normal     purchaseTab  \= order.purchase\_details.normal     \# \==========================================     \# 2\. Generate Signature / Hash Per Tab     \# \==========================================     hashNormal   \= generateSignature(normalTab)     hashFinal    \= generateSignature(finalTab)     hashPurchase \= generateSignature(purchaseTab)     \# \==========================================     \# 3\. Compare Tabs     \# \==========================================     IF hashNormal \== hashFinal AND hashNormal \== hashPurchase THEN         \# \======================================         \# CASE A: All Tabs Identical         \# \======================================         payload \= buildPromoPayload(normalTab)         promoResult \= PromotionService.consult(payload)         \# Inject same promo result to all tabs         injectPromo(order.details.normal, promoResult)         injectPromo(order.details\_final.normal, promoResult)         injectPromo(order.purchase\_details.normal, promoResult)     ELSE         \# \======================================         \# CASE B: Tabs Different (Parallel Call)         \# \======================================         futureNormal \= ASYNC PromotionService.consult(             buildPromoPayload(normalTab)         )         futureFinal \= ASYNC PromotionService.consult(             buildPromoPayload(finalTab)         )         futurePurchase \= ASYNC PromotionService.consult(             buildPromoPayload(purchaseTab)         )         \# Wait all async jobs         WAIT ALL futureNormal, futureFinal, futurePurchase         promoNormal   \= futureNormal.result         promoFinal    \= futureFinal.result         promoPurchase \= futurePurchase.result         \# Inject promo per tab         injectPromo(order.details.normal, promoNormal)         injectPromo(order.details\_final.normal, promoFinal)         injectPromo(order.purchase\_details.normal, promoPurchase)     END IF     \# \==========================================     \# 4\. Return Response     \# \==========================================     RETURN order END FUNCTION  |
| :---- |

| purchase\_details | source  |
| :---- | :---- |
| "promo1": 0,  | promo1 |
| "promo2": 0,  | promo2 |
| "promo3": 0,  | promo3 |
| "promo4": 0,  | promo4 |
| "promo5": 0,  | promo5 |

url : **/promotions/consult**

### **Contoh Request JSON** 

| {   "order\_date": "2025-12-17",   "outlet\_id": 1404,   "salesman\_id": 359,   "wh\_id": 63,   "details": \[      {       "pro\_id": 709,       "qty1": 10,       "qty2": 0,       "qty3": 0,       "gross\_value": 2000000     },     {       "pro\_id": 710,       "qty1": 12,       "qty2": 0,       "qty3": 0,       "gross\_value": 3600000     },     {       "pro\_id": 711,       "qty1": 10,       "qty2": 0,       "qty3": 0,       "gross\_value": 800000     }       \] }  |
| :---- |

### **Request (Source Data)** 

| purchase\_details | source  (dari hasil Fetch order \+ details.normal \+ details\_final \+ purchase\_details)  atau sama dengan response sales order detail |
| :---- | :---- |
| order\_date | ro\_date |
| outlet\_id | outlet\_id |
| salesman\_id | salesman\_id |
| wh\_id | wh\_id |
| details.pro\_id | **TAB PURCHASE ORDER=**       details.normal\[\].pro\_id **TAB SALES ORDER=**       purchase\_details.normal\[\].pro\_id **TAB FINAL ORDER=**       details\_final.normal\[\].pro\_id |
| details.qty1 | **TAB PURCHASE ORDER=**       details.normal\[\].qty1 **TAB SALES ORDER=**       purchase\_details.normal\[\].qty1 **TAB FINAL ORDER=**       details\_final.normal\[\].qty1 |
| details.qty2 | **TAB PURCHASE ORDER=**       details.normal\[\].qty2 **TAB SALES ORDER=**       purchase\_details.normal\[\].qty2 **TAB FINAL ORDER=**       details\_final.normal\[\].qty2 |
| details.qty3 | **TAB PURCHASE ORDER=**       details.normal\[\].qty3 **TAB SALES ORDER=**       purchase\_details.normal\[\].qty3 **TAB FINAL ORDER=**       details\_final.normal\[\].qty3 |
| details.gross\_value | (qty1 \* sell\_price1) \+(qty2 \* sell\_price2) \+ (qty3 \* sell\_price3) |

### **Contoh Response JSON :** 

| {     "message": "Consulted V2 Successfully",     "data": \[         {             "promo\_id": "slab01",             "promo\_desc": "Syarat Value & Reward Percentage (Value)",             "slab\_id": "695ba0009f8c1325e300ae2f",             "slab\_desc": "slab 3 Value",             "slab\_reward": 3,             "slab\_reward\_uom": "",             "slab\_reward\_type": "percentage",             "slab\_per\_scope": "",             "total\_gross\_value": 5600000,             "products\_eligible": \[                 709,                 710             \],             "reward\_value": null,             "reward\_percentage": \[                 {                     "pro\_id": 709,                     "gross\_value": 2000000,                     "promo1": 60000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1940000                 },                 {                     "pro\_id": 710,                     "gross\_value": 3600000,                     "promo1": 108000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 3492000                 }             \],             "reward\_product": null         }     \],     "request\_id": "fd9a72d5-b1af-4c75-83dd-a08a579d9877" }  |
| :---- |

### **Response JSON :** 

| FUNCTION aggregatePromo(promoResponse): productMap \= MAP\<pro\_id, PromoAccumulator\> FOR each promo IN promoResponse.data: promoCode \= promo.promo\_id IF promo.reward\_percentage IS NOT NULL: FOR each item IN promo.reward\_percentage: proId \= item.pro\_id IF productMap\[proId\] NOT EXISTS: productMap\[proId\] \= { promo1: 0, promo2: 0, promo3: 0, promo4: 0, promo5: 0, remarks: \[\] } productMap\[proId\].promo1 \+= item.promo1 productMap\[proId\].promo2 \+= item.promo2 productMap\[proId\].promo3 \+= item.promo3 productMap\[proId\].promo4 \+= item.promo4 productMap\[proId\].promo5 \+= item.promo5 productMap\[proId\].remarks.ADD(promoCode) RETURN productMap END FUNCTION |
| :---- |

### **Case : Reward Percentage** 

jika reward\_value \= NULL   
**reward\_percentage** \= NOT NULL   
reward\_product \= NULL 

| Response Detail Order | Source Data |
| :---- | :---- |
| details.normal\[\].pro\_id | data promo diambil dari **reward\_percentage**  |
| purchase\_details.normal\[\].pro\_id |  |
| details\_final.normal\[\].pro\_id |  |
| details.normal\[\].promo1 | SUM(data\[\].reward\_percentage\[\].promo1) GROUP BY pro\_id |
| purchase\_details.normal\[\].promo1 |  |
| details\_final.normal\[\].promo1 |  |
| details.normal\[\].promo2 | SUM(data\[\].reward\_percentage\[\].promo2) GROUP BY pro\_id |
| purchase\_details.normal\[\].promo2 |  |
| details\_final.normal\[\].promo2 |  |
| details.normal\[\].promo3 | SUM(data\[\].reward\_percentage\[\].promo3) GROUP BY pro\_id |
| purchase\_details.normal\[\].promo3 |  |
| details\_final.normal\[\].promo3 |  |
| details.normal\[\].promo4 | SUM(data\[\].reward\_percentage\[\].promo4) GROUP BY pro\_id |
|  |  |
|  |  |
| purchase\_details.normal\[\].promo4 |  |
| details\_final.normal\[\].promo4 |  |
| details.normal\[\].promo5 | SUM(data\[\].reward\_percentage\[\].promo5) GROUP BY pro\_id |
| purchase\_details.normal\[\].promo5 |  |
| details\_final.normal\[\].promo5 |  |
| details.normal\[\].total\_promo | promo1 \+ promo2 \+ promo3 \+ promo4 \+ promo5 |
| purchase\_details.normal\[\].total\_promo |  |
| details\_final.normal\[\].total\_promo |  |
| purchase\_details.normal\[\].remarks  | FUNCTION buildRemarksPerProduct(promoResponse, orderProducts):     productRemarks \= MAP\<pro\_id, ARRAY\<promo\_id\>\>     FOR each product IN orderProducts:         proId \= product.pro\_id         remarks \= \[\]         FOR each promo IN promoResponse.data:             IF proId IN promo.products\_eligible:                 remarks.ADD(promo.promo\_id)         productRemarks\[proId\] \= remarks     RETURN productRemarks END FUNCTION ***contoh** saat request pro\_id \= 709, dan 710  ternyata pro\_id \= 709 ada di products\_eligible di slab01, dan slab02 sedangkan pro\_id \= 710 ada di products\_eligible di slab02  jadi remarks untuk pro\_id \= 709 \--\> slab01, dan slab02  dan remarks untuk pro\_id \= 710 \--\> slab02*  |
| purchase\_details.normal\[\].remarks |  |
| details\_final.normal\[\].remarks |  |
| details.final\_remarks | data\[\].promo\_id (ini adalah promo\_id yang di dapat untuk semua product\_id) |
| purchase\_details.final\_remarks |  |
| details\_final.final\_remarks |  |

#### **Response yang sudah di sesuaikan :** 

| {     "message": "",     "data": {         "ro\_no": "SO2602210001",         "source": "web",         "is\_proforma\_inv": false,         "order\_no": null,         "ro\_date": "2026-02-21",         "val\_date": null,         "salesman\_id": 357,         "salesman\_code": "2025120201",         "sales\_name": "Ilham Abdullah",         "wh\_id": 63,         "wh\_code": "008",         "wh\_name": "Gudang Utama",         "outlet\_id": 404,         "outlet\_code": "P00025Ar",         "outlet\_name": "Princes 9",         "outlet\_address1": "PERUM CANDRALOKA BLOK AA.12 NO.9",         "outlet\_address2": "",         "inv\_addr1": "PERUM CANDRALOKA BLOK AA.12 NO.4",         "inv\_addr2": "",         "delivery\_date": "2026-02-21",         "po\_no": null,         "vehicle\_no": null,         "pay\_type": 3,         "pay\_type\_name": "Credit",         "reff\_no": null,         "mobile\_id": 1,         "sub\_total": 310000,         "sub\_total\_final": 310000,         "disc": 0,         "disc\_value": 0,         "disc\_value\_final": 0,         "promo\_value": 0,         "promo\_value\_final": 0,         "promo\_bg\_value": 0,         "promo\_bg\_value\_final": 0,         "cash\_disc\_value": 0,         "tot\_disc1": 0,         "tot\_disc2": 0,         "vat": 11,         "vat\_value": 37200,         "vat\_value\_final": 37200,         "total": 347200,         "total\_final": 347200,         "data\_status": 2,         "data\_status\_name": "Processed",         "data\_source": 1,         "updated\_at": "2026-02-21T08:03:03.471215Z",         "updated\_by\_name": "Dist IDE Sda",         "due\_date": null,         "details": {             "normal": \[                 {                     "order\_detail\_id": 5215,                     "seq\_no": 0,                     "pro\_id": 688,                     "pro\_code": "FF-0101-0040",                     "pro\_name": "Cavendish Hand",                     "order\_status": "",                     "item\_type": 1,                     "qty": 12,                     "qty\_final": 12,                     "qty\_po": 12,                     "qty\_po1": null,                     "qty\_po2": null,                     "qty\_po3": null,                     "qty1": 0,                     "qty2": 0,                     "qty3": 1,                     "qty4": null,                     "qty5": null,                     "qty1\_final": 0,                     "qty2\_final": 0,                     "qty3\_final": 1,                     "qty4\_final": null,                     "qty5\_final": null,                     "qty1\_stok": 0,                     "qty2\_stok": 0,                     "qty3\_stok": 59,                     "purch\_price1": 10000,                     "purch\_price2": 200000,                     "purch\_price3": 300000,                     "purch\_price4": 0,                     "purch\_price5": 0,                     "sell\_price1": 110000,                     "sell\_price2": 210000,                     "sell\_price3": 310000,                     "sell\_price4": 0,                     "sell\_price5": 0,                     "sell\_price\_po1": null,                     "sell\_price\_po2": null,                     "sell\_price\_po3": null,                     "sell\_price\_final1": 110000,                     "sell\_price\_final2": 210000,                     "sell\_price\_final3": 310000,                     "sell\_price\_system1": 110000,                     "sell\_price\_system2": 110000,                     "sell\_price\_system3": 110000,                     "sell\_price\_system4": null,                     "sell\_price\_system5": null,                     "amount": 347200,                     "amount\_final": 347200,                     "promo\_value": 0,                     "promo\_value\_final": 0,                     "disc\_value": 0,                     "disc\_value\_final": 0,                     "batch\_no": null,                     "exp\_date": null,                     "vat": 12,                     "vat\_bg": 0,                     "vat\_lg\_sell": 0,                     "vat\_value": 37200,                     "vat\_value\_final": 37200,                     "vat\_bg\_value": 0,                     "vat\_lg\_value": null,                     "vat\_lg\_sell\_value": 0,                     "disc\_po": 0,                     "vat\_value\_po": 37200,                     "unit\_id1": "PNCH",                     "unit\_id2": "CRT",                     "unit\_id3": "CRT",                     "unit\_id4": "",                     "unit\_id5": "",                     "conv\_unit2": 12,                     "conv\_unit3": 1,                     "conv\_unit4": 0,                     "conv\_unit5": 0,                     "notes": "",                     "discount\_id": "",                     "promo1": 120000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "promo\_total": 120000,                     "remarks": \["slab01", "slab02"\]                 }             \],             "final\_remarks" : \["slab01", "slab02"\],             "reward\_products" : \[\],             "promo": null         },         "details\_final": {             "normal": \[                 {                     "order\_detail\_id": 5215,                     "seq\_no": 0,                     "pro\_id": 688,                     "pro\_code": "FF-0101-0040",                     "pro\_name": "Cavendish Hand",                     "order\_status": "",                     "item\_type": 1,                     "qty": 12,                     "qty\_final": 12,                     "qty\_po": 12,                     "qty\_po1": null,                     "qty\_po2": null,                     "qty\_po3": null,                     "qty1": 0,                     "qty2": 0,                     "qty3": 1,                     "qty4": null,                     "qty5": null,                     "qty1\_final": 0,                     "qty2\_final": 0,                     "qty3\_final": 1,                     "qty4\_final": null,                     "qty5\_final": null,                     "qty1\_stok": 0,                     "qty2\_stok": 0,                     "qty3\_stok": 59,                     "purch\_price1": 10000,                     "purch\_price2": 200000,                     "purch\_price3": 300000,                     "purch\_price4": 0,                     "purch\_price5": 0,                     "sell\_price1": 110000,                     "sell\_price2": 210000,                     "sell\_price3": 310000,                     "sell\_price4": 0,                     "sell\_price5": 0,                     "sell\_price\_po1": null,                     "sell\_price\_po2": null,                     "sell\_price\_po3": null,                     "sell\_price\_final1": 110000,                     "sell\_price\_final2": 210000,                     "sell\_price\_final3": 310000,                     "sell\_price\_system1": 110000,                     "sell\_price\_system2": 110000,                     "sell\_price\_system3": 110000,                     "sell\_price\_system4": null,                     "sell\_price\_system5": null,                     "amount": 347200,                     "amount\_final": 347200,                     "promo\_value": 0,                     "promo\_value\_final": 0,                     "disc\_value": 0,                     "disc\_value\_final": 0,                     "batch\_no": null,                     "exp\_date": null,                     "vat": 12,                     "vat\_bg": 0,                     "vat\_lg\_sell": 0,                     "vat\_value": 37200,                     "vat\_value\_final": 37200,                     "vat\_bg\_value": 0,                     "vat\_lg\_value": null,                     "vat\_lg\_sell\_value": 0,                     "disc\_po": 0,                     "vat\_value\_po": 37200,                     "unit\_id1": "PNCH",                     "unit\_id2": "CRT",                     "unit\_id3": "CRT",                     "unit\_id4": "",                     "unit\_id5": "",                     "conv\_unit2": 12,                     "conv\_unit3": 1,                     "conv\_unit4": 0,                     "conv\_unit5": 0,                     "notes": "",                     "discount\_id": "",                     "promo1": 120000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "promo\_total": 120000,                     "remarks": \["slab01", "slab02"\]                 }             \],             "final\_remarks" : \["slab01", "slab02"\],             "reward\_products" : \[\],             "promo": null         },         "purchase\_details": {           "normal": \[             {               "order\_detail\_id": 5215,               "seq\_no": 0,               "pro\_id": 688,               "pro\_code": "FF-0101-0040",               "pro\_name": "Cavendish Hand",               "order\_status": "2",               "item\_type": 1,               "qty": 12,               "qty\_final": 12,               "qty\_po": 12,               "qty\_po1": null,               "qty\_po2": null,               "qty\_po3": null,               "qty1": 0,               "qty2": 0,               "qty3": 1,               "qty4": null,               "qty5": null,               "qty1\_final": 0,               "qty2\_final": 0,               "qty3\_final": 1,               "qty4\_final": null,               "qty5\_final": null,               "qty1\_stok": 0,               "qty2\_stok": 0,               "qty3\_stok": 59,               "purch\_price1": 10000,               "purch\_price2": 200000,               "purch\_price3": 300000,               "purch\_price4": 0,               "purch\_price5": 0,               "sell\_price1": 110000,               "sell\_price2": 210000,               "sell\_price3": 310000,               "sell\_price4": 0,               "sell\_price5": 0,               "sell\_price\_po1": 110000,               "sell\_price\_po2": 210000,               "sell\_price\_po3": 310000,               "sell\_price\_final1": 110000,               "sell\_price\_final2": 210000,               "sell\_price\_final3": 310000,               "sell\_price\_system1": 110000,               "sell\_price\_system2": 110000,               "sell\_price\_system3": 110000,               "sell\_price\_system4": null,               "sell\_price\_system5": null,               "amount": 347200,               "amount\_final": 347200,               "promo\_value": 0,               "promo\_value\_final": 0,               "disc\_value": 0,               "disc\_value\_final": 0,               "batch\_no": null,               "exp\_date": null,               "vat": 12,               "vat\_bg": 0,               "vat\_lg\_sell": 0,               "vat\_value": 0,               "vat\_value\_final": 0,               "vat\_bg\_value": 0,               "vat\_lg\_value": null,               "vat\_lg\_sell\_value": 0,               "disc\_po": 0,               "vat\_value\_po": 37200,               "unit\_id1": "PNCH",               "unit\_id2": "CRT",               "unit\_id3": "CRT",               "unit\_id4": "",               "unit\_id5": "",               "conv\_unit2": 12,               "conv\_unit3": 1,               "conv\_unit4": 0,               "conv\_unit5": 0,               "notes": "",               "discount\_id": "",               "promo1": 120000,               "promo2": 0,               "promo3": 0,               "promo4": 0,               "promo5": 0,               "promo\_total": 120000,               "remarks": \["slab01", "slab02"\]             }           \],           "final\_remarks": \["slab01", "slab02"\],           "reward\_products": \[\],           "promo": \[\]         },         "remarks": null,         "tr\_code": null,         "is\_closed": false,         "notes": "",         "invoice\_no": null,         "invoice\_date": null,         "is\_printed": false,         "printed\_by": null,         "printed\_by\_name": null,         "printed\_at": null,         "validate\_stok": true,         "validate\_stok\_message": "Sufficient Stock",         "validate\_credit\_limit": true,         "validate\_credit\_limit\_message": "Within Limit",         "validate\_overdue": true,         "validate\_overdue\_message": "Allowed (Unlimited)",         "validate\_outstanding": true,         "validate\_outstanding\_message": "Allowed (Unlimited)",         "validate\_summary": true,         "credit\_limit\_type": 2,         "credit\_limit\_action": 1,         "credit\_limit\_action\_name": "Warning",         "sales\_inv\_limit\_type": null,         "sales\_inv\_limit\_action": null,         "sales\_inv\_limit\_action\_name": "",         "obs\_type": null,         "obs\_limit\_action": null,         "obs\_limit\_action\_name": "",         "order\_approval\_request\_id": null,         "order\_approval\_request\_emp\_approval\_status": null     },     "request\_id": "dea1bf04-3920-426a-90ff-f42040486208" }  |
| :---- |

### **Case : Reward Value** 

jika **reward\_value** \= NOT NULL   
reward\_percentage \= NULL   
reward\_product \= NULL 

| {     "message": "Consulted V2 Successfully",     "data": \[         {             "promo\_id": "slab02",             "promo\_desc": "Syarat Value & Reward Value Perorder (Quantity)",             "slab\_id": "695deb11e3dc19398700bd68",             "slab\_desc": "slab 3 Quantity Smallest",             "slab\_reward": 20000,             "slab\_rule\_type": "value",             "slab\_rule\_uom": "",             "slab\_reward\_uom": "",             "slab\_reward\_type": "fixed\_value",             "slab\_per\_scope": "per\_order",             "total\_gross\_value": 3600000,             "products\_eligible": \[                 710,                 673             \],             "reward\_value": \[                 {                     "pro\_id": 710,                     "gross\_value": 1600000,                     "promo1": 10000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1590000                 },                 {                     "pro\_id": 673,                     "gross\_value": 2000000,                     "promo1": 10000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1990000                 }             \],             "reward\_percentage": null,             "reward\_product": null         }     \],     "request\_id": "e1a30a69-4acc-44fb-81c6-c0b2500d8f9e" }  |
| :---- |

| Response Detail Order | Source Data |
| :---- | :---- |
| details.normal\[\].pro\_id | data promo diambil dari **reward\_value**  |
| purchase\_details.normal\[\].pro\_id |  |
| details\_final.normal\[\].pro\_id |  |
| details.normal\[\].promo1 | SUM(data\[\].reward\_value\[\].promo1) GROUP BY pro\_id |
| purchase\_details.normal\[\].promo1 |  |
| details\_final.normal\[\].promo1 |  |
| details.normal\[\].promo2 | SUM(data\[\].reward\_value\[\].promo2) GROUP BY pro\_id |
| purchase\_details.normal\[\].promo2 |  |
| details\_final.normal\[\].promo2 |  |
| details.normal\[\].promo3 | SUM(data\[\].reward\_value\[\].promo3) GROUP BY pro\_id |
| purchase\_details.normal\[\].promo3 |  |
| details\_final.normal\[\].promo3 |  |
| details.normal\[\].promo4 | SUM(data\[\].reward\_value\[\].promo4) GROUP BY pro\_id |
|  |  |
|  |  |
| purchase\_details.normal\[\].promo4 |  |
| details\_final.normal\[\].promo4 |  |
| details.normal\[\].promo5 | SUM(data\[\].reward\_value\[\].promo5) GROUP BY pro\_id |
| purchase\_details.normal\[\].promo5 |  |
| details\_final.normal\[\].promo5 |  |
| details.normal\[\].total\_promo | promo1 \+ promo2 \+ promo3 \+ promo4 \+ promo5 |
| purchase\_details.normal\[\].total\_promo |  |
| details\_final.normal\[\].total\_promo |  |
| purchase\_details.normal\[\].remarks  | FUNCTION buildRemarksPerProduct(promoResponse, orderProducts):     productRemarks \= MAP\<pro\_id, ARRAY\<promo\_id\>\>     FOR each product IN orderProducts:         proId \= product.pro\_id         remarks \= \[\]         FOR each promo IN promoResponse.data:             IF proId IN promo.products\_eligible:                 remarks.ADD(promo.promo\_id)         productRemarks\[proId\] \= remarks     RETURN productRemarks END FUNCTION ***contoh** saat request pro\_id \= 709, dan 710  ternyata pro\_id \= 709 ada di products\_eligible di slab01, dan slab02 sedangkan pro\_id \= 710 ada di products\_eligible di slab02  jadi remarks untuk pro\_id \= 709 \--\> slab01, dan slab02  dan remarks untuk pro\_id \= 710 \--\> slab02*  |
| purchase\_details.normal\[\].remarks |  |
| details\_final.normal\[\].remarks |  |
| details.final\_remarks | data\[\].promo\_id (ini adalah promo\_id yang di dapat untuk semua product\_id) |
| purchase\_details.final\_remarks |  |
| details\_final.final\_remarks |  |

Contoh response promo pada detail : 

| {   "pro\_id": 709,   "promo1": 120000,   "promo2": 0,   "promo3": 0,   "promo4": 0,   "promo5": 0,   "promo\_total": 120000,   "remarks": \["slab01", "slab02"\] }  |
| :---- |

#### **Response yang sudah di sesuaikan :** 

| {     "message": "",     "data": {         "ro\_no": "SO2602210001",         "source": "web",         "is\_proforma\_inv": false,         "order\_no": null,         "ro\_date": "2026-02-21",         "val\_date": null,         "salesman\_id": 357,         "salesman\_code": "2025120201",         "sales\_name": "Ilham Abdullah",         "wh\_id": 63,         "wh\_code": "008",         "wh\_name": "Gudang Utama",         "outlet\_id": 404,         "outlet\_code": "P00025Ar",         "outlet\_name": "Princes 9",         "outlet\_address1": "PERUM CANDRALOKA BLOK AA.12 NO.9",         "outlet\_address2": "",         "inv\_addr1": "PERUM CANDRALOKA BLOK AA.12 NO.4",         "inv\_addr2": "",         "delivery\_date": "2026-02-21",         "po\_no": null,         "vehicle\_no": null,         "pay\_type": 3,         "pay\_type\_name": "Credit",         "reff\_no": null,         "mobile\_id": 1,         "sub\_total": 310000,         "sub\_total\_final": 310000,         "disc": 0,         "disc\_value": 0,         "disc\_value\_final": 0,         "promo\_value": 0,         "promo\_value\_final": 0,         "promo\_bg\_value": 0,         "promo\_bg\_value\_final": 0,         "cash\_disc\_value": 0,         "tot\_disc1": 0,         "tot\_disc2": 0,         "vat": 11,         "vat\_value": 37200,         "vat\_value\_final": 37200,         "total": 347200,         "total\_final": 347200,         "data\_status": 2,         "data\_status\_name": "Processed",         "data\_source": 1,         "updated\_at": "2026-02-21T08:03:03.471215Z",         "updated\_by\_name": "Dist IDE Sda",         "due\_date": null,         "details": {             "normal": \[                 {                     "order\_detail\_id": 5215,                     "seq\_no": 0,                     "pro\_id": 688,                     "pro\_code": "FF-0101-0040",                     "pro\_name": "Cavendish Hand",                     "order\_status": "",                     "item\_type": 1,                     "qty": 12,                     "qty\_final": 12,                     "qty\_po": 12,                     "qty\_po1": null,                     "qty\_po2": null,                     "qty\_po3": null,                     "qty1": 0,                     "qty2": 0,                     "qty3": 1,                     "qty4": null,                     "qty5": null,                     "qty1\_final": 0,                     "qty2\_final": 0,                     "qty3\_final": 1,                     "qty4\_final": null,                     "qty5\_final": null,                     "qty1\_stok": 0,                     "qty2\_stok": 0,                     "qty3\_stok": 59,                     "purch\_price1": 10000,                     "purch\_price2": 200000,                     "purch\_price3": 300000,                     "purch\_price4": 0,                     "purch\_price5": 0,                     "sell\_price1": 110000,                     "sell\_price2": 210000,                     "sell\_price3": 310000,                     "sell\_price4": 0,                     "sell\_price5": 0,                     "sell\_price\_po1": null,                     "sell\_price\_po2": null,                     "sell\_price\_po3": null,                     "sell\_price\_final1": 110000,                     "sell\_price\_final2": 210000,                     "sell\_price\_final3": 310000,                     "sell\_price\_system1": 110000,                     "sell\_price\_system2": 110000,                     "sell\_price\_system3": 110000,                     "sell\_price\_system4": null,                     "sell\_price\_system5": null,                     "amount": 347200,                     "amount\_final": 347200,                     "promo\_value": 0,                     "promo\_value\_final": 0,                     "disc\_value": 0,                     "disc\_value\_final": 0,                     "batch\_no": null,                     "exp\_date": null,                     "vat": 12,                     "vat\_bg": 0,                     "vat\_lg\_sell": 0,                     "vat\_value": 37200,                     "vat\_value\_final": 37200,                     "vat\_bg\_value": 0,                     "vat\_lg\_value": null,                     "vat\_lg\_sell\_value": 0,                     "disc\_po": 0,                     "vat\_value\_po": 37200,                     "unit\_id1": "PNCH",                     "unit\_id2": "CRT",                     "unit\_id3": "CRT",                     "unit\_id4": "",                     "unit\_id5": "",                     "conv\_unit2": 12,                     "conv\_unit3": 1,                     "conv\_unit4": 0,                     "conv\_unit5": 0,                     "notes": "",                     "discount\_id": "",                     "promo1": 120000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "promo\_total": 120000,                     "remarks": \["slab01", "slab02"\]                 }             \],             "final\_remarks" : \["slab01", "slab02"\],             "reward\_products" : \[\],             "promo": null         },         "details\_final": {             "normal": \[                 {                     "order\_detail\_id": 5215,                     "seq\_no": 0,                     "pro\_id": 688,                     "pro\_code": "FF-0101-0040",                     "pro\_name": "Cavendish Hand",                     "order\_status": "",                     "item\_type": 1,                     "qty": 12,                     "qty\_final": 12,                     "qty\_po": 12,                     "qty\_po1": null,                     "qty\_po2": null,                     "qty\_po3": null,                     "qty1": 0,                     "qty2": 0,                     "qty3": 1,                     "qty4": null,                     "qty5": null,                     "qty1\_final": 0,                     "qty2\_final": 0,                     "qty3\_final": 1,                     "qty4\_final": null,                     "qty5\_final": null,                     "qty1\_stok": 0,                     "qty2\_stok": 0,                     "qty3\_stok": 59,                     "purch\_price1": 10000,                     "purch\_price2": 200000,                     "purch\_price3": 300000,                     "purch\_price4": 0,                     "purch\_price5": 0,                     "sell\_price1": 110000,                     "sell\_price2": 210000,                     "sell\_price3": 310000,                     "sell\_price4": 0,                     "sell\_price5": 0,                     "sell\_price\_po1": null,                     "sell\_price\_po2": null,                     "sell\_price\_po3": null,                     "sell\_price\_final1": 110000,                     "sell\_price\_final2": 210000,                     "sell\_price\_final3": 310000,                     "sell\_price\_system1": 110000,                     "sell\_price\_system2": 110000,                     "sell\_price\_system3": 110000,                     "sell\_price\_system4": null,                     "sell\_price\_system5": null,                     "amount": 347200,                     "amount\_final": 347200,                     "promo\_value": 0,                     "promo\_value\_final": 0,                     "disc\_value": 0,                     "disc\_value\_final": 0,                     "batch\_no": null,                     "exp\_date": null,                     "vat": 12,                     "vat\_bg": 0,                     "vat\_lg\_sell": 0,                     "vat\_value": 37200,                     "vat\_value\_final": 37200,                     "vat\_bg\_value": 0,                     "vat\_lg\_value": null,                     "vat\_lg\_sell\_value": 0,                     "disc\_po": 0,                     "vat\_value\_po": 37200,                     "unit\_id1": "PNCH",                     "unit\_id2": "CRT",                     "unit\_id3": "CRT",                     "unit\_id4": "",                     "unit\_id5": "",                     "conv\_unit2": 12,                     "conv\_unit3": 1,                     "conv\_unit4": 0,                     "conv\_unit5": 0,                     "notes": "",                     "discount\_id": "",                     "promo1": 120000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "promo\_total": 120000,                     "remarks": \["slab01", "slab02"\]                 }             \],             "final\_remarks" : \["slab01", "slab02"\],             "reward\_products" : \[\],             "promo": null         },         "purchase\_details": {           "normal": \[             {               "order\_detail\_id": 5215,               "seq\_no": 0,               "pro\_id": 688,               "pro\_code": "FF-0101-0040",               "pro\_name": "Cavendish Hand",               "order\_status": "2",               "item\_type": 1,               "qty": 12,               "qty\_final": 12,               "qty\_po": 12,               "qty\_po1": null,               "qty\_po2": null,               "qty\_po3": null,               "qty1": 0,               "qty2": 0,               "qty3": 1,               "qty4": null,               "qty5": null,               "qty1\_final": 0,               "qty2\_final": 0,               "qty3\_final": 1,               "qty4\_final": null,               "qty5\_final": null,               "qty1\_stok": 0,               "qty2\_stok": 0,               "qty3\_stok": 59,               "purch\_price1": 10000,               "purch\_price2": 200000,               "purch\_price3": 300000,               "purch\_price4": 0,               "purch\_price5": 0,               "sell\_price1": 110000,               "sell\_price2": 210000,               "sell\_price3": 310000,               "sell\_price4": 0,               "sell\_price5": 0,               "sell\_price\_po1": 110000,               "sell\_price\_po2": 210000,               "sell\_price\_po3": 310000,               "sell\_price\_final1": 110000,               "sell\_price\_final2": 210000,               "sell\_price\_final3": 310000,               "sell\_price\_system1": 110000,               "sell\_price\_system2": 110000,               "sell\_price\_system3": 110000,               "sell\_price\_system4": null,               "sell\_price\_system5": null,               "amount": 347200,               "amount\_final": 347200,               "promo\_value": 0,               "promo\_value\_final": 0,               "disc\_value": 0,               "disc\_value\_final": 0,               "batch\_no": null,               "exp\_date": null,               "vat": 12,               "vat\_bg": 0,               "vat\_lg\_sell": 0,               "vat\_value": 0,               "vat\_value\_final": 0,               "vat\_bg\_value": 0,               "vat\_lg\_value": null,               "vat\_lg\_sell\_value": 0,               "disc\_po": 0,               "vat\_value\_po": 37200,               "unit\_id1": "PNCH",               "unit\_id2": "CRT",               "unit\_id3": "CRT",               "unit\_id4": "",               "unit\_id5": "",               "conv\_unit2": 12,               "conv\_unit3": 1,               "conv\_unit4": 0,               "conv\_unit5": 0,               "notes": "",               "discount\_id": "",               "promo1": 120000,               "promo2": 0,               "promo3": 0,               "promo4": 0,               "promo5": 0,               "promo\_total": 120000,               "remarks": \["slab01", "slab02"\]             }           \],           "final\_remarks": \["slab01", "slab02"\],           "reward\_products": \[\],           "promo": \[\]         },         "remarks": null,         "tr\_code": null,         "is\_closed": false,         "notes": "",         "invoice\_no": null,         "invoice\_date": null,         "is\_printed": false,         "printed\_by": null,         "printed\_by\_name": null,         "printed\_at": null,         "validate\_stok": true,         "validate\_stok\_message": "Sufficient Stock",         "validate\_credit\_limit": true,         "validate\_credit\_limit\_message": "Within Limit",         "validate\_overdue": true,         "validate\_overdue\_message": "Allowed (Unlimited)",         "validate\_outstanding": true,         "validate\_outstanding\_message": "Allowed (Unlimited)",         "validate\_summary": true,         "credit\_limit\_type": 2,         "credit\_limit\_action": 1,         "credit\_limit\_action\_name": "Warning",         "sales\_inv\_limit\_type": null,         "sales\_inv\_limit\_action": null,         "sales\_inv\_limit\_action\_name": "",         "obs\_type": null,         "obs\_limit\_action": null,         "obs\_limit\_action\_name": "",         "order\_approval\_request\_id": null,         "order\_approval\_request\_emp\_approval\_status": null     },     "request\_id": "dea1bf04-3920-426a-90ff-f42040486208" }  |
| :---- |

### **Case : Reward Product** 

jika reward\_value \= NULL   
reward\_percentage \= NULL   
**reward\_product** \= NOT NULL 

| {     "message": "Consulted V2 Successfully",     "data": \[         {             "promo\_id": "slab04",             "promo\_desc": "Syarat Value & Reward Produk (Quantity)",             "slab\_id": "694caff8cc48b19f80002ceb",             "slab\_desc": "slab 1 Quantity Smallest",             "slab\_reward": 1,             "slab\_rule\_type": "value",             "slab\_rule\_uom": "",             "slab\_reward\_uom": "smallest",             "slab\_reward\_type": "product",             "slab\_per\_scope": "",             "total\_gross\_value": 1600000,             "products\_eligible": \[                 711,                 710             \],             "reward\_value": null,             "reward\_percentage": null,             "reward\_product": \[                 {                     "pro\_id": 491,                     "qty1": 1,                     "qty2": 0,                     "qty3": 0,                     "gross\_value": 10000,                     "promo1": 10000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0                 }             \]         }     \],     "request\_id": "1da40a69-4acc-44fb-81c6-c0b2500d8f9e" }  |
| :---- |

| Response Detail Order | Source Data |
| :---- | :---- |
| details.normal\[\].pro\_id | data promo diambil dari **reward\_product**  |
| purchase\_details.normal\[\].pro\_id |  |
| details\_final.normal\[\].pro\_id |  |
| details.normal\[\].promo1 | SUM(data\[\].reward\_product\[\].promo1) GROUP BY pro\_id |
| purchase\_details.normal\[\].promo1 |  |
| details\_final.normal\[\].promo1 |  |
| details.normal\[\].promo2 | SUM(data\[\].reward\_product\[\].promo2) GROUP BY pro\_id |
| purchase\_details.normal\[\].promo2 |  |
| details\_final.normal\[\].promo2 |  |
| details.normal\[\].promo3 | SUM(data\[\].reward\_product\[\].promo3) GROUP BY pro\_id |
| purchase\_details.normal\[\].promo3 |  |
| details\_final.normal\[\].promo3 |  |
| details.normal\[\].promo4 | SUM(data\[\].reward\_product\[\].promo4) GROUP BY pro\_id |
|  |  |
|  |  |
| purchase\_details.normal\[\].promo4 |  |
| details\_final.normal\[\].promo4 |  |
| details.normal\[\].promo5 | SUM(data\[\].reward\_product\[\].promo5) GROUP BY pro\_id |
| purchase\_details.normal\[\].promo5 |  |
| details\_final.normal\[\].promo5 |  |
| details.normal\[\].total\_promo | promo1 \+ promo2 \+ promo3 \+ promo4 \+ promo5 |
| purchase\_details.normal\[\].total\_promo |  |
| details\_final.normal\[\].total\_promo |  |
| purchase\_details.normal\[\].remarks  | FUNCTION buildRemarksPerProduct(promoResponse, orderProducts):     productRemarks \= MAP\<pro\_id, ARRAY\<promo\_id\>\>     FOR each product IN orderProducts:         proId \= product.pro\_id         remarks \= \[\]         FOR each promo IN promoResponse.data:             IF proId IN promo.products\_eligible:                 remarks.ADD(promo.promo\_id)         productRemarks\[proId\] \= remarks     RETURN productRemarks END FUNCTION ***contoh** saat request pro\_id \= 709, dan 710  ternyata pro\_id \= 709 ada di products\_eligible di slab01, dan slab02 sedangkan pro\_id \= 710 ada di products\_eligible di slab02  jadi remarks untuk pro\_id \= 709 \--\> slab01, dan slab02  dan remarks untuk pro\_id \= 710 \--\> slab02*  |
| purchase\_details.normal\[\].remarks |  |
| details\_final.normal\[\].remarks |  |
| details.final\_remarks | data\[\].promo\_id (ini adalah promo\_id yang di dapat untuk semua product\_id) |
| purchase\_details.final\_remarks |  |
| details\_final.final\_remarks |  |

khusus untuk reward product, tambahkan Array to Object reward\_product pada ***details, purchase\_details, details\_final***

|  "reward\_products" : \[                 {                     "pro\_id": 709,                     "pro\_code" : "FF-0101-0040",                     "pro\_name" : "Cavendish Hand",                     "qty1": 0,                     "qty2": 0,                     "qty3": 1,                     "sell\_price1": 110000,                     "sell\_price2": 210000,                     "sell\_price3": 310000                 }             \],  |
| :---- |

| Response Detail Order | Source Data |
| :---- | :---- |
| **details.reward\_products\[\]** |  |
| pro\_id | reward\_product\[\].pro\_id |
| pro\_code | ambil dari mst.m\_product.pro\_code |
| pro\_name | ambil dari mst.m\_product.pro\_name |
| qty1 | reward\_product\[\].qty1 |
| qty2 | reward\_product\[\].qty2 |
| qty3 | reward\_product\[\].qty3 |
| sell\_price1 | ambil dari mst.m\_product.sell\_price1 |
| sell\_price2 | ambil dari mst.m\_product.sell\_price2 |
| sell\_price3 | ambil dari mst.m\_product.sell\_price3 |
| **purchase\_details.reward\_products\[\]** |  |
| pro\_id | reward\_product\[\].pro\_id |
| pro\_code | ambil dari mst.m\_product.pro\_code |
| pro\_name | ambil dari mst.m\_product.pro\_name |
| qty1 | reward\_product\[\].qty1 |
| qty2 | reward\_product\[\].qty2 |
| qty3 | reward\_product\[\].qty3 |
| sell\_price1 | ambil dari mst.m\_product.sell\_price1 |
| sell\_price2 | ambil dari mst.m\_product.sell\_price2 |
| sell\_price3 | ambil dari mst.m\_product.sell\_price3 |
| **details\_final.reward\_products\[\]** |  |
| pro\_id | reward\_product\[\].pro\_id |
| pro\_code | ambil dari mst.m\_product.pro\_code |
| pro\_name | ambil dari mst.m\_product.pro\_name |
| qty1 | reward\_product\[\].qty1 |
| qty2 | reward\_product\[\].qty2 |
| qty3 | reward\_product\[\].qty3 |
| sell\_price1 | ambil dari mst.m\_product.sell\_price1 |
| sell\_price2 | ambil dari mst.m\_product.sell\_price2 |
| sell\_price3 | ambil dari mst.m\_product.sell\_price3 |

#### **Response yang sudah di sesuaikan :** 

| {     "message": "",     "data": {         "ro\_no": "SO2602210001",         "source": "web",         "is\_proforma\_inv": false,         "order\_no": null,         "ro\_date": "2026-02-21",         "val\_date": null,         "salesman\_id": 357,         "salesman\_code": "2025120201",         "sales\_name": "Ilham Abdullah",         "wh\_id": 63,         "wh\_code": "008",         "wh\_name": "Gudang Utama",         "outlet\_id": 404,         "outlet\_code": "P00025Ar",         "outlet\_name": "Princes 9",         "outlet\_address1": "PERUM CANDRALOKA BLOK AA.12 NO.9",         "outlet\_address2": "",         "inv\_addr1": "PERUM CANDRALOKA BLOK AA.12 NO.4",         "inv\_addr2": "",         "delivery\_date": "2026-02-21",         "po\_no": null,         "vehicle\_no": null,         "pay\_type": 3,         "pay\_type\_name": "Credit",         "reff\_no": null,         "mobile\_id": 1,         "sub\_total": 310000,         "sub\_total\_final": 310000,         "disc": 0,         "disc\_value": 0,         "disc\_value\_final": 0,         "promo\_value": 0,         "promo\_value\_final": 0,         "promo\_bg\_value": 0,         "promo\_bg\_value\_final": 0,         "cash\_disc\_value": 0,         "tot\_disc1": 0,         "tot\_disc2": 0,         "vat": 11,         "vat\_value": 37200,         "vat\_value\_final": 37200,         "total": 347200,         "total\_final": 347200,         "data\_status": 2,         "data\_status\_name": "Processed",         "data\_source": 1,         "updated\_at": "2026-02-21T08:03:03.471215Z",         "updated\_by\_name": "Dist IDE Sda",         "due\_date": null,         "details": {             "normal": \[                 {                     "order\_detail\_id": 5215,                     "seq\_no": 0,                     "pro\_id": 688,                     "pro\_code": "FF-0101-0040",                     "pro\_name": "Cavendish Hand",                     "order\_status": "",                     "item\_type": 1,                     "qty": 12,                     "qty\_final": 12,                     "qty\_po": 12,                     "qty\_po1": null,                     "qty\_po2": null,                     "qty\_po3": null,                     "qty1": 0,                     "qty2": 0,                     "qty3": 1,                     "qty4": null,                     "qty5": null,                     "qty1\_final": 0,                     "qty2\_final": 0,                     "qty3\_final": 1,                     "qty4\_final": null,                     "qty5\_final": null,                     "qty1\_stok": 0,                     "qty2\_stok": 0,                     "qty3\_stok": 59,                     "purch\_price1": 10000,                     "purch\_price2": 200000,                     "purch\_price3": 300000,                     "purch\_price4": 0,                     "purch\_price5": 0,                     "sell\_price1": 110000,                     "sell\_price2": 210000,                     "sell\_price3": 310000,                     "sell\_price4": 0,                     "sell\_price5": 0,                     "sell\_price\_po1": null,                     "sell\_price\_po2": null,                     "sell\_price\_po3": null,                     "sell\_price\_final1": 110000,                     "sell\_price\_final2": 210000,                     "sell\_price\_final3": 310000,                     "sell\_price\_system1": 110000,                     "sell\_price\_system2": 110000,                     "sell\_price\_system3": 110000,                     "sell\_price\_system4": null,                     "sell\_price\_system5": null,                     "amount": 347200,                     "amount\_final": 347200,                     "promo\_value": 0,                     "promo\_value\_final": 0,                     "disc\_value": 0,                     "disc\_value\_final": 0,                     "batch\_no": null,                     "exp\_date": null,                     "vat": 12,                     "vat\_bg": 0,                     "vat\_lg\_sell": 0,                     "vat\_value": 37200,                     "vat\_value\_final": 37200,                     "vat\_bg\_value": 0,                     "vat\_lg\_value": null,                     "vat\_lg\_sell\_value": 0,                     "disc\_po": 0,                     "vat\_value\_po": 37200,                     "unit\_id1": "PNCH",                     "unit\_id2": "CRT",                     "unit\_id3": "CRT",                     "unit\_id4": "",                     "unit\_id5": "",                     "conv\_unit2": 12,                     "conv\_unit3": 1,                     "conv\_unit4": 0,                     "conv\_unit5": 0,                     "notes": "",                     "discount\_id": "",                     "promo1": 120000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "promo\_total": 120000,                     "remarks": \["slab01", "slab02"\]                 }             \],             "final\_remarks" : \["slab01", "slab02"\],             "reward\_products" : \[                  {                    "pro\_id": 709,                    "gross\_value": 2000000,                    "promo1": 60000,                    "promo2": 0,                    "promo3": 0,                    "promo4": 0,                    "promo5": 0,                    "net\_value": 1940000                  }             \],             "promo": null         },         "details\_final": {             "normal": \[                 {                     "order\_detail\_id": 5215,                     "seq\_no": 0,                     "pro\_id": 688,                     "pro\_code": "FF-0101-0040",                     "pro\_name": "Cavendish Hand",                     "order\_status": "",                     "item\_type": 1,                     "qty": 12,                     "qty\_final": 12,                     "qty\_po": 12,                     "qty\_po1": null,                     "qty\_po2": null,                     "qty\_po3": null,                     "qty1": 0,                     "qty2": 0,                     "qty3": 1,                     "qty4": null,                     "qty5": null,                     "qty1\_final": 0,                     "qty2\_final": 0,                     "qty3\_final": 1,                     "qty4\_final": null,                     "qty5\_final": null,                     "qty1\_stok": 0,                     "qty2\_stok": 0,                     "qty3\_stok": 59,                     "purch\_price1": 10000,                     "purch\_price2": 200000,                     "purch\_price3": 300000,                     "purch\_price4": 0,                     "purch\_price5": 0,                     "sell\_price1": 110000,                     "sell\_price2": 210000,                     "sell\_price3": 310000,                     "sell\_price4": 0,                     "sell\_price5": 0,                     "sell\_price\_po1": null,                     "sell\_price\_po2": null,                     "sell\_price\_po3": null,                     "sell\_price\_final1": 110000,                     "sell\_price\_final2": 210000,                     "sell\_price\_final3": 310000,                     "sell\_price\_system1": 110000,                     "sell\_price\_system2": 110000,                     "sell\_price\_system3": 110000,                     "sell\_price\_system4": null,                     "sell\_price\_system5": null,                     "amount": 347200,                     "amount\_final": 347200,                     "promo\_value": 0,                     "promo\_value\_final": 0,                     "disc\_value": 0,                     "disc\_value\_final": 0,                     "batch\_no": null,                     "exp\_date": null,                     "vat": 12,                     "vat\_bg": 0,                     "vat\_lg\_sell": 0,                     "vat\_value": 37200,                     "vat\_value\_final": 37200,                     "vat\_bg\_value": 0,                     "vat\_lg\_value": null,                     "vat\_lg\_sell\_value": 0,                     "disc\_po": 0,                     "vat\_value\_po": 37200,                     "unit\_id1": "PNCH",                     "unit\_id2": "CRT",                     "unit\_id3": "CRT",                     "unit\_id4": "",                     "unit\_id5": "",                     "conv\_unit2": 12,                     "conv\_unit3": 1,                     "conv\_unit4": 0,                     "conv\_unit5": 0,                     "notes": "",                     "discount\_id": "",                     "promo1": 120000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "promo\_total": 120000,                     "remarks": \["slab01", "slab02"\]                 }             \],             "final\_remarks" : \["slab01", "slab02"\],             "reward\_products" : \[                  {                    "pro\_id": 709,                    "gross\_value": 2000000,                    "promo1": 60000,                    "promo2": 0,                    "promo3": 0,                    "promo4": 0,                    "promo5": 0,                    "net\_value": 1940000                  }             \],             "promo": null         },         "purchase\_details": {           "normal": \[             {               "order\_detail\_id": 5215,               "seq\_no": 0,               "pro\_id": 688,               "pro\_code": "FF-0101-0040",               "pro\_name": "Cavendish Hand",               "order\_status": "2",               "item\_type": 1,               "qty": 12,               "qty\_final": 12,               "qty\_po": 12,               "qty\_po1": null,               "qty\_po2": null,               "qty\_po3": null,               "qty1": 0,               "qty2": 0,               "qty3": 1,               "qty4": null,               "qty5": null,               "qty1\_final": 0,               "qty2\_final": 0,               "qty3\_final": 1,               "qty4\_final": null,               "qty5\_final": null,               "qty1\_stok": 0,               "qty2\_stok": 0,               "qty3\_stok": 59,               "purch\_price1": 10000,               "purch\_price2": 200000,               "purch\_price3": 300000,               "purch\_price4": 0,               "purch\_price5": 0,               "sell\_price1": 110000,               "sell\_price2": 210000,               "sell\_price3": 310000,               "sell\_price4": 0,               "sell\_price5": 0,               "sell\_price\_po1": 110000,               "sell\_price\_po2": 210000,               "sell\_price\_po3": 310000,               "sell\_price\_final1": 110000,               "sell\_price\_final2": 210000,               "sell\_price\_final3": 310000,               "sell\_price\_system1": 110000,               "sell\_price\_system2": 110000,               "sell\_price\_system3": 110000,               "sell\_price\_system4": null,               "sell\_price\_system5": null,               "amount": 347200,               "amount\_final": 347200,               "promo\_value": 0,               "promo\_value\_final": 0,               "disc\_value": 0,               "disc\_value\_final": 0,               "batch\_no": null,               "exp\_date": null,               "vat": 12,               "vat\_bg": 0,               "vat\_lg\_sell": 0,               "vat\_value": 0,               "vat\_value\_final": 0,               "vat\_bg\_value": 0,               "vat\_lg\_value": null,               "vat\_lg\_sell\_value": 0,               "disc\_po": 0,               "vat\_value\_po": 37200,               "unit\_id1": "PNCH",               "unit\_id2": "CRT",               "unit\_id3": "CRT",               "unit\_id4": "",               "unit\_id5": "",               "conv\_unit2": 12,               "conv\_unit3": 1,               "conv\_unit4": 0,               "conv\_unit5": 0,               "notes": "",               "discount\_id": "",               "promo1": 120000,               "promo2": 0,               "promo3": 0,               "promo4": 0,               "promo5": 0,               "promo\_total": 120000,               "remarks": \["slab01", "slab02"\]             }           \],           "final\_remarks": \["slab01", "slab02"\],           "reward\_products": \[              {                "pro\_id": 709,                "gross\_value": 2000000,                "promo1": 60000,                "promo2": 0,                "promo3": 0,                "promo4": 0,                "promo5": 0,                "net\_value": 1940000              }           \],           "promo": \[\]         },         "remarks": null,         "tr\_code": null,         "is\_closed": false,         "notes": "",         "invoice\_no": null,         "invoice\_date": null,         "is\_printed": false,         "printed\_by": null,         "printed\_by\_name": null,         "printed\_at": null,         "validate\_stok": true,         "validate\_stok\_message": "Sufficient Stock",         "validate\_credit\_limit": true,         "validate\_credit\_limit\_message": "Within Limit",         "validate\_overdue": true,         "validate\_overdue\_message": "Allowed (Unlimited)",         "validate\_outstanding": true,         "validate\_outstanding\_message": "Allowed (Unlimited)",         "validate\_summary": true,         "credit\_limit\_type": 2,         "credit\_limit\_action": 1,         "credit\_limit\_action\_name": "Warning",         "sales\_inv\_limit\_type": null,         "sales\_inv\_limit\_action": null,         "sales\_inv\_limit\_action\_name": "",         "obs\_type": null,         "obs\_limit\_action": null,         "obs\_limit\_action\_name": "",         "order\_approval\_request\_id": null,         "order\_approval\_request\_emp\_approval\_status": null     },     "request\_id": "dea1bf04-3920-426a-90ff-f42040486208" } |
| :---- |

5. ## **Validate Order** 

URL : 

| https://be.scyllax.online/sales/v1/validate-order |
| :---- |

Method : POST   
Enhance : Yes 

#### **Enhance Payload :  (FE dan BE)** 

| Before :  {"product":\[{"pro\_id":792,"qty1":0,"qty2":0,"qty3":1}\],"outlet\_id":1445,"wh\_id":254,"total":452880} After :  {   "product": \[     {       "pro\_id": 792,       "qty1": 0,       "qty2": 0,       "qty3": 10,       "qty\_change1": 0,       "qty\_change2": 0,       "qty\_change3": 5,     }   \],   "outlet\_id": 1445,   "wh\_id": 254,   "total": 452880,  } |
| :---- |

Response : 

| {     "message": "",     "data": {         "validate1\_success": true,         "validate1\_message": "Sufficient Stock",         "validate2\_success": true,         "validate2\_message": "Allowed (Unlimited)",         "validate2\_value": 0,         "validate3\_success": true,         "validate3\_message": "Allowed (Unlimited)",         "validate3\_value": 0,         "validate4\_success": true,         "validate4\_message": "Allowed (Unlimited)",         "validate4\_value": 0,         "validate\_summary\_success": true     },     "request\_id": "038f5c7f-6ae0-46ef-af05-ac496d67812c" } |
| :---- |

#### **Enhance function di BE :** 

| perubahan di warehouse stock  ![][image6] |
| :---- |
| **Perubahan Rumus di Backend :**  (wh stock \+ (oncust \- qty changes) **compare** dengan qty order |
| **case qty changes bertambah ** *![][image7]*  |
| **case qty changes berkurang  ![][image8]**  |

6. ## **Salesman Lookup ( Enhance)**

URL : 

| https://be.scyllax.online/scylla-pjp/api/v1/list-salesman?q=\&page=1\&limit=70 |
| :---- |

Method : GET   
![][image9]

**Enhance** :   
add response : **allow\_input\_price →** mst.salesman.allow\_input\_price

**Response :** 

| {     "message": "",     "data": \[         {             "emp\_id": 299,             "emp\_code": "SLD001",             "emp\_name": "Jhon Takbor",             "allow\_input\_price": TRUE,             "sales\_team\_id": 62,             "sales\_team\_code": "15",             "sales\_team\_name": "united",             "sales\_name": "Jhon Takbor",             "opr\_type": "O",             "opr\_type\_canvas": "",             "opr\_type\_name": "Taking Order",             "opr\_type\_name\_canvas": "Unknown",             "is\_bonus\_rep": false,             "trans\_date": null,             "wh\_id": 254,             "wh\_code": "001",             "wh\_name": "gudang baik",             "wh\_name\_canvas": "",             "wh\_canvas\_id": null,             "wh\_name\_view": "gudang baik",             "inc\_grp\_id": 0,             "official\_id": 0,             "official\_type": 0,             "hierarchy\_code": "",             "sale\_system": "N",             "sm\_is\_transfer": false,             "sm\_valid\_route": false,             "sm\_geoloc\_valid": false,             "sm\_radius": 0,             "sm\_password": "",             "is\_active": true,             "updated\_at": "2025-09-02T09:13:31.881682Z",             "updated\_by\_name": "Admin DT Bangka 01",             "image\_url": "",             "details": {                 "product\_line": \[                     {                         "pl\_id": 75,                         "m\_salesman\_product\_type\_id": 2577,                         "ref\_id": 75,                         "pl\_code": "PL00",                         "pl\_name": "Non Product Jual"                     },                     {                         "pl\_id": 76,                         "m\_salesman\_product\_type\_id": 2578,                         "ref\_id": 76,                         "pl\_code": "PL01",                         "pl\_name": "Home Care"                     },                     {                         "pl\_id": 77,                         "m\_salesman\_product\_type\_id": 2579,                         "ref\_id": 77,                         "pl\_code": "PL02",                         "pl\_name": "Hair Care"                     },                     {                         "pl\_id": 78,                         "m\_salesman\_product\_type\_id": 2580,                         "ref\_id": 78,                         "pl\_code": "PL03",                         "pl\_name": "Oral Care"                     },                     {                         "pl\_id": 79,                         "m\_salesman\_product\_type\_id": 2581,                         "ref\_id": 79,                         "pl\_code": "PL99",                         "pl\_name": "Other Care"                     },                     {                         "pl\_id": 80,                         "m\_salesman\_product\_type\_id": 2582,                         "ref\_id": 80,                         "pl\_code": "PL04",                         "pl\_name": "Pestisida"                     }                 \],                 "brand": \[                     {                         "pl\_id": 75,                         "m\_salesman\_product\_type\_id": 2583,                         "ref\_id": 168,                         "brand\_code": "B00",                         "brand\_name": "Non Brand"                     },                     {                         "pl\_id": 76,                         "m\_salesman\_product\_type\_id": 2584,                         "ref\_id": 169,                         "brand\_code": "B01",                         "brand\_name": "VAPE"                     },                     {                         "pl\_id": 77,                         "m\_salesman\_product\_type\_id": 2585,                         "ref\_id": 170,                         "brand\_code": "B02",                         "brand\_name": "Sampho"                     },                     {                         "pl\_id": 80,                         "m\_salesman\_product\_type\_id": 2586,                         "ref\_id": 171,                         "brand\_code": "B03",                         "brand\_name": "AEROSOL SPRAY   "                     }                 \],                 "sub\_brand": \[                     {                         "pl\_id": 75,                         "m\_salesman\_product\_type\_id": 2587,                         "ref\_id": 186,                         "sbrand1\_code": "SB00",                         "sbrand1\_name": "Non Sub Brgand"                     },                     {                         "pl\_id": 77,                         "m\_salesman\_product\_type\_id": 2588,                         "ref\_id": 187,                         "sbrand1\_code": "SB01",                         "sbrand1\_name": "ReJoice"                     },                     {                         "pl\_id": 80,                         "m\_salesman\_product\_type\_id": 2589,                         "ref\_id": 188,                         "sbrand1\_code": "S02",                         "sbrand1\_name": "PEstisida"                     }                 \]             },             "sm\_is\_barcode": false,             "sm\_is\_photo\_profile": false,             "is\_active\_canvas": false,             "is\_taking\_order": true         }     \],     "paging": {         "total\_record": 1,         "page\_current": 1,         "page\_limit": 9999,         "page\_total": 1     },     "request\_id": "6927b4ef250e1766daa12c9d" }  |
| :---- |

7. ## **Outlet Lookup (No Enhance)**

URL : 

| https://be.scyllax.online/master/v1/outlets?is\_active=1\&verification\_status=1\&outlet\_id=\&q=\&page=1\&limit=99999 |
| :---- |

Method : GET   
![][image10]

**Response :** 

| {  "message": "",  "data": \[    {      "outlet\_id": 246,      "outlet\_code": "B0005",      "outlet\_name": "Beni Simanman",      "barcode": "",      "outlet\_status": 0,      "address1": "Jl.Banga Raya no 20 C Kenag Jakarta Selatan",      "address2": "",      "city": "",      "zip\_code": "098760",      "phone\_no": "087868016020",      "wa\_no": "",      "fax\_no": "",      "email": "beni@gmail.com",      "disc\_grp\_id": 111,      "disc\_grp\_code": "00",      "disc\_grp\_name": "Non Group Discount",      "ot\_loc\_id": 88,      "ot\_loc\_code": "0012",      "ot\_loc\_name": "PASAR",      "ot\_grp\_id": 71,      "ot\_grp\_code": "000",      "ot\_grp\_name": "Non Group",      "price\_grp\_id": 41,      "price\_grp\_code": "",      "price\_grp\_name": "",      "district\_id": 66,      "district\_code": "D66",      "district\_name": "Gunung Putri",      "beat\_id": 0,      "beat\_code": "",      "beat\_name": "",      "sbeat\_id": 0,      "sbeat\_code": "",      "sbeat\_name": "",      "ot\_class\_id": 56,      "ot\_class\_code": "000",      "ot\_class\_name": "NON Classifi",      "industry\_id": 38,      "industry\_code": "453",      "industry\_name": "Non Industry",      "market\_id": 50,      "market\_code": "",      "market\_name": "",      "top": 0,      "due\_date": "2025-11-27",      "payment\_type": 1,      "is\_contra\_bon": false,      "plu\_grp\_id": 0,      "plu\_grp\_code": "",      "plu\_grp\_name": "",      "conv\_grp\_id": 0,      "conv\_grp\_code": "",      "conv\_grp\_name": "",      "disc\_inv\_id": 0,      "disc\_inv\_code": "",      "disc\_inv\_name": "",      "agent\_from": "",      "credit\_limit\_type": 2,      "credit\_limit": 500000,      "sales\_inv\_limit\_type": 0,      "sales\_inv\_limit": 0,      "avg\_sales\_week": 0,      "avg\_sales\_month": 0,      "first\_trans\_date": "",      "last\_trans\_date": "",      "first\_week\_no": 0,      "ot\_start\_date": "",      "ot\_reg\_date": "",      "building\_own": 1,      "dob": "",      "ar\_status": 1,      "ar\_total": 0,      "closed\_date": "",      "is\_emb\_bail": false,      "tax\_name": "",      "tax\_addr1": "",      "tax\_addr2": "",      "tax\_city": "",      "tax\_no": "",      "tax\_invoice\_form": 1,      "tax\_invoice\_form\_name": "Standart",      "owner\_name": "",      "owner\_addr1": "",      "owner\_addr2": "",      "owner\_city": "",      "owner\_phone\_no": "",      "owner\_id\_no": "",      "delv\_addr1": "Jl.Banga Raya no 20 C Kenag Jakarta Selatan",      "delv\_city": "3737",      "delv\_latitude": "5",      "delv\_longitude": "5",      "delv\_addr2": "",      "delv\_city2": "",      "delv\_latitude2": "",      "delv\_longitude2": "",      "inv\_addr1": "Jl.Banga Raya no 20 C Kenag Jakarta Selatan",      "inv\_addr2": "",      "inv\_city": "3737",      "is\_active": true,      "updated\_by": 105,      "updated\_at": "2025-06-04T06:17:34.312664Z",      "updated\_by\_name": "Admin DT Bangka 01",      "latitude": "5",      "longitude": "5",      "ot\_type\_id": 62,      "ot\_type\_code": "000",      "ot\_type\_name": "Retail",      "is\_obs": false,      "obs": 0,      "outlet\_ward\_id": "123456",      "outlet\_ward": "Jaya Mekar",      "outlet\_sub\_district\_id": "363636",      "outlet\_sub\_district": "Cileungsi",      "outlet\_regency\_id": "3737",      "outlet\_regency": "Bogor",      "outlet\_province\_id": "36",      "outlet\_province": "Jawa Barat",      "is\_wa\_no": false,      "delv\_ward\_id": "123456",      "delv\_ward": "Jaya Mekar",      "delv\_sub\_district\_id": "363636",      "delv\_sub\_district": "Cileungsi",      "delv\_regency\_id": "3737",      "delv\_regency": "Bogor",      "delv\_province\_id": "36",      "delv\_province": "Jawa Barat",      "delv\_zip\_code": "098760",      "delv\_is\_same\_addr": true,      "delv\_ward\_id2": null,      "delv\_ward2": null,      "delv\_sub\_district\_id2": null,      "delv\_sub\_district2": null,      "delv\_regency\_id2": null,      "delv\_regency2": null,      "delv\_province\_id2": null,      "delv\_province2": null,      "delv\_zip\_code2": null,      "inv\_ward\_id": "123456",      "inv\_ward": "Jaya Mekar",      "inv\_sub\_district\_id": "363636",      "inv\_sub\_district": "Cileungsi",      "inv\_regency\_id": "3737",      "inv\_regency": "Bogor",      "inv\_province\_id": "36",      "inv\_province": "Jawa Barat",      "inv\_zip\_code": "098760",      "inv\_is\_same\_addr": true,      "verification\_status": 1,      "verification\_status\_name": "Approved",      "obs\_type": "",      "credit\_limit\_action": 2,      "sales\_inv\_limit\_action": 0,      "obs\_limit\_action": 0,      "outlet\_establishment\_date": "2025-06-04T00:00:00Z",      "identity\_type": "National ID",      "identity\_no": "09876789098760"    }  \],  "paging": {    "total\_record": 5,    "page\_current": 1,    "page\_limit": 99999,    "page\_total": 1  },  "request\_id": "692859193d65fe6458b7f268" }  |
| :---- |

8. ## **Product Lookup Enhancement**

URL : 

| https://best.scyllax.online/master/v1/products?mode=lookup\_dist\_price\&dist\_price\_group\_id=2\&outlet\_id=246\&order\_date=2025-11-27\&q=\&page=1\&limit=70 |
| :---- |

**Method : GET**   
**![][image11]**

**Response :** 

| {    "message": "",    "data": \[        {            "pro\_id": 711,            "pro\_code": "DD-FM02-0003",            "pro\_name": "Hometown Choco 450ml SP1",            "unit\_id1": "PCS",            "unit\_id2": "BOX",            "unit\_id3": "BOX",            "unit\_id4": "",            "unit\_id5": "",            "conv\_unit2": 24,            "conv\_unit3": 1,            "conv\_unit4": 0,            "conv\_unit5": 0,            "purch\_price1": 1,            "purch\_price2": 1,            "purch\_price3": 1,            "purch\_price4": 0,            "purch\_price5": 0,            "sell\_price1": 1,            "sell\_price2": 1,            "sell\_price3": 1,            "sell\_price4": 0,            "sell\_price5": 0,            "vat": 0,            "pl\_id": 124,            "pl\_code": "DD",            "pl\_name": "Dairy and Daily",            "brand\_id": 126,            "brand\_code": "HT",            "sbrand1\_id": 145,            "sbrand1\_code": "BAA",            "sbrand1\_name": "Liquid Milk"        },        {            "pro\_id": 655,            "pro\_code": "SP00006",            "pro\_name": "Sunpride Lype Salt  \\u0026 Garlic 55gr",            "unit\_id1": "PNCH",            "unit\_id2": "CRT",            "unit\_id3": "CRT",            "unit\_id4": "",            "unit\_id5": "",            "conv\_unit2": 12,            "conv\_unit3": 1,            "conv\_unit4": 0,            "conv\_unit5": 0,            "purch\_price1": 10000,            "purch\_price2": 110000,            "purch\_price3": 900000,            "purch\_price4": 0,            "purch\_price5": 0,            "sell\_price1": 11000,            "sell\_price2": 120000,            "sell\_price3": 950000,            "sell\_price4": 0,            "sell\_price5": 0,            "vat": 11,            "pl\_id": 120,            "pl\_code": "A001",            "pl\_name": "Sunpride",            "brand\_id": 124,            "brand\_code": "AP01",            "sbrand1\_id": 143,            "sbrand1\_code": "AP001",            "sbrand1\_name": "Sunpride"        }    \],    "paging": {        "total\_record": 168,        "page\_current": 1,        "page\_limit": 70,        "page\_total": 3    },    "request\_id": "69286b1f3d65fe6458b7f275" }  |
| :---- |

### **Enhance :** 

| Perbaiki filter pada query :   select  *mp*.cust\_id , *mp*.pro\_code ,  *mp*.distributor\_id , *mp*.pro\_name ,  *mp*.pro\_status ,  *mp*.is\_active ,  *mp*.is\_del   from mst.m\_product *mp*  where *mp*.cust\_id \='C22001'  and *mp*.distributor\_id \=68  //jika user principal mp.distributor is null and *mp*.is\_active is true  and mp.is\_del is false  order by *mp*.pro\_code asc  |
| :---- |
| **Data Test User Distributor:**  Email : [admdistributor01.bangka@pratesis.com](mailto:admdistributor01.bangka@pratesis.com) dapat di samakan dengan fitur master product :  Klik Master \> Product \> Product  ![][image12] |
| **Data Test User Principal:**  email : [princ@idetama.id](mailto:princ@idetama.id)  dapat di samakan dengan fitur master product :  Klik Master \> Product \> Product  Filter berdasarkan status \= Active dan Created By : Admin Principal 1 ![][image13] |

9. ## **Process Order**  

API ini digunakan ketika user membuat sales order melalui web.

**URL :** 

| https://best.scyllax.online/sales/v1/orders/ |
| :---- |

**Method** : POST

**Payload** : 

| {  "ro\_no": "SO2507100002",  "source": null,  "is\_performa\_inv": null,  "order\_no": null,  "ro\_date": "2025-07-10",  "val\_date": null,  "salesman\_id": 22,  "salesman\_code": "22222",  "sales\_name": "Salesman Bangka 01",  "wh\_id": 1,  "wh\_code": "008",  "wh\_name": "Gudang Utama",  "outlet\_id": 246,  "outlet\_code": "B0005",  "outlet\_name": "Beni Simanman",  "outlet\_address1": "Jl.Banga Raya no 20 C Kenag Jakarta Selatan",  "outlet\_address2": "",  "inv\_addr1": "Jl.Banga Raya no 20 C Kenag Jakarta Selatan",  "inv\_addr2": "",  "delivery\_date": "2025-07-10",  "po\_no": null,  "vehicle\_no": null,  "pay\_type": 1,  "pay\_type\_name": "Cash On Delivery",  "reff\_no": null,  "mobile\_id": 1,  "sub\_total": 950400,  "sub\_total\_final": 950400,  "disc": 0,  "disc\_value": 0,  "disc\_value\_final": 0,  "promo\_value": 0,  "promo\_value\_final": 0,  "promo\_bg\_value": 0,  "promo\_bg\_value\_final": 0,  "cash\_disc\_value": 0,  "tot\_disc1": 0,  "tot\_disc2": 0,  "vat": 11,  "vat\_value": 104544,  "vat\_value\_final": 104544,  "total": 1054944,  "total\_final": 1054944,  "data\_status": 1,  "data\_status\_name": "Need Review",  "data\_source": null,  "updated\_at": "2025-07-10T09:34:09.404405Z",  "updated\_by\_name": "Admin DT Bangka 03",  "due\_date": null,  "details": {    "normal": \[      {        "order\_detail\_id": 2901,        "seq\_no": 0,        "pro\_id": 473,        "pro\_code": "01002",        "pro\_name": "VAPE ULTRA MAT 45 LVNDR 36 PCS",        "order\_status": "",        "item\_type": 1,        "qty": 108,        "qty\_final": 108,        "qty\_po": 36,        "qty1": 0,        "qty2": 0,        "qty3": 3,        "qty4": null,        "qty5": null,        "qty1\_final": 0,        "qty2\_final": 0,        "qty3\_final": 3,        "qty4\_final": null,        "qty5\_final": null,        "qty1\_stok": 0,        "qty2\_stok": 0,        "qty3\_stok": 976,        "purch\_price1": 7117,        "purch\_price2": 256216,        "purch\_price3": 256216,        "purch\_price4": 0,        "purch\_price5": 0,        "sell\_price1": 8800,        "sell\_price2": 316800,        "sell\_price3": 316800,        "sell\_price4": 0,        "sell\_price5": 0,        "sell\_price\_system1": null,        "sell\_price\_system2": null,        "sell\_price\_system3": null,        "sell\_price\_system4": null,        "sell\_price\_system5": null,        "amount": 1054944,        "amount\_final": 1054944,        "promo\_value": 0,        "promo\_value\_final": 0,        "disc\_value": 0,        "disc\_value\_final": 0,        "batch\_no": null,        "exp\_date": null,        "vat": 11,        "vat\_bg": 0,        "vat\_lg\_sell": 0,        "vat\_value": 104544,        "vat\_value\_final": 104544,        "vat\_bg\_value": 0,        "vat\_lg\_value": null,        "vat\_lg\_sell\_value": 0,        "unit\_id1": "PCS",        "unit\_id2": "KRT",        "unit\_id3": "KRT",        "unit\_id4": "",        "unit\_id5": "",        "conv\_unit2": 36,        "conv\_unit3": 1,        "conv\_unit4": 0,        "conv\_unit5": 0,        "notes": "",        "discount\_id": "",        "promoa\_id": null,        "is\_promo": false      }    \]  },  "details\_final": {    "normal": \[      {        "order\_detail\_id": 2901,        "seq\_no": 0,        "pro\_id": 473,        "pro\_code": "01002",        "pro\_name": "VAPE ULTRA MAT 45 LVNDR 36 PCS",        "order\_status": "",        "item\_type": 1,        "qty": 108,        "qty\_final": 108,        "qty\_po": 36,        "qty1": 0,        "qty2": 0,        "qty3": 3,        "qty4": null,        "qty5": null,        "qty1\_final": 0,        "qty2\_final": 0,        "qty3\_final": 3,        "qty4\_final": null,        "qty5\_final": null,        "qty1\_stok": 0,        "qty2\_stok": 0,        "qty3\_stok": 976,        "purch\_price1": 7117,        "purch\_price2": 256216,        "purch\_price3": 256216,        "purch\_price4": 0,        "purch\_price5": 0,        "sell\_price1": 8800,        "sell\_price2": 316800,        "sell\_price3": 316800,        "sell\_price4": 0,        "sell\_price5": 0,        "sell\_price\_system1": null,        "sell\_price\_system2": null,        "sell\_price\_system3": null,        "sell\_price\_system4": null,        "sell\_price\_system5": null,        "amount": 1054944,        "amount\_final": 1054944,        "promo\_value": 0,        "promo\_value\_final": 0,        "disc\_value": 0,        "disc\_value\_final": 0,        "batch\_no": null,        "exp\_date": null,        "vat": 11,        "vat\_bg": 0,        "vat\_lg\_sell": 0,        "vat\_value": 104544,        "vat\_value\_final": 104544,        "vat\_bg\_value": 0,        "vat\_lg\_value": null,        "vat\_lg\_sell\_value": 0,        "unit\_id1": "PCS",        "unit\_id2": "KRT",        "unit\_id3": "KRT",        "unit\_id4": "",        "unit\_id5": "",        "conv\_unit2": 36,        "conv\_unit3": 1,        "conv\_unit4": 0,        "conv\_unit5": 0,        "notes": "",        "discount\_id": "",        "promoa\_id": null      }    \],    "promo": null  },  "purchase\_details": { "normal": null, "promo": null },  "remarks": null,  "tr\_code": null,  "is\_closed": false,  "notes": "coba ",  "invoice\_no": null,  "invoice\_date": null,  "is\_printed": false,  "printed\_by": null,  "printed\_by\_name": null,  "printed\_at": null,  "validate\_stok": true,  "validate\_stok\_message": "Sufficient Stock",  "validate\_credit\_limit": false,  "validate\_credit\_limit\_message": "Over Limit (11.598.096)",  "validate\_overdue": true,  "validate\_overdue\_message": "Allowed (Unlimited)",  "validate\_outstanding": true,  "validate\_outstanding\_message": "Allowed (Unlimited)",  "validate\_summary": false,  "credit\_limit\_type": 2,  "credit\_limit\_action": 2,  "credit\_limit\_action\_name": "Restricted",  "sales\_inv\_limit\_type": null,  "sales\_inv\_limit\_action": null,  "sales\_inv\_limit\_action\_name": "",  "obs\_type": null,  "obs\_limit\_action": null,  "obs\_limit\_action\_name": "",  "order\_approval\_request\_id": null,  "order\_approval\_request\_emp\_approval\_status": null }  |
| :---- |

**Response:** 

| {     "message": "Berhasil Diperbarui",     "request\_id": "5e506370-e334-40ac-8675-a8e40de2ba46" }  |
| :---- |

Enhane BE : 

* add logic , if source \= 2  (mobile)(expect : tab purchase order)  
- qty1 , qty2, qty3 → update ke field **qty\_po1, qty\_po2, qty\_po3**  
- sell\_price1, sell\_price2, sell\_price3 → **update ke  field sell\_price\_po1, sell\_price\_po2, sell\_price\_po3**  
-   
* add logic , if source \= 1 (web)   
- qty1 , qty2, qty3   
  → update ke field **qty1, qty2, qty3**(tab sales order)  
  → update ke field **qty1\_final, qty2\_final, qty3\_final**(tab final order)  
- sell\_price1, sell\_price2, sell\_price3   
  → update ke  field **sell\_price1, sell\_price2, sell\_price3**(tab sales order)

  → update ke  field **sell\_price\_final1, sell\_pricefinal2, sell\_pricefinal3**(tab sales order)


- disc\_value **(tab sales order)**  
  disc\_value\_final **(tab final order)**  
- vat\_value **(tab sales order)**

  **vat\_value\_final** (tab final order)

10. ## **Orders Conversion**

API ini digunakan untuk melakukan edit data produk .   
URL : 

| https://best.scyllax.online/sales/v1/orders/conversion |
| :---- |

Enhance endpoint : 

* Edit Final Order **( Enhance)**  
  * Tambahkan request berikut : apabila FE edit pada tab Final Order, req yang dikirim adalah qty1\_final, qty2\_final, dan qty6\_final (Tanpa mengubah value pada qty1, qty2, dan qty3) 

| {    "pro\_id":529,    "qty1\_final":0,    "qty2\_final":0,    "qty3\_final":6 } catatan : qty1, qty2, qty3 tidak di set  |
| :---- |

Response : (conversion mengikuti eksisting, tidak ada perubahan perhitungan namun sesuaikan pada field yang di edit, nilai qty1, qty2, dan qty3 tidak berubah) 

| {     "message": "",     "data": {         "qty1\_final": 0,         "qty2\_final": 0,         "qty3\_final": 18,         "total\_qty": 180     },     "request\_id": "c45a070a-642f-46b8-8718-373ddaa96712" }  |
| :---- |

* Edit Sales Order : **No Enhance**

https://best.scyllax.online/sales/v1/orders/conversion  
Request :

|  {    pro\_id: 635,     qty1: 0,     qty2: 0,     qty3: 34 } |
| :---- |

Response : 

| {     "message": "",     "data": {         "qty1": 0,         "qty2": 0,         "qty3": 18,         "total\_qty": 180     },     "request\_id": "c45a070a-642f-46b8-8718-373ddaa96712" }  |
| :---- |

11. ## **Edit Order ( Enhance) → enhance**

API ini digunakan ketika user klik process pada tab Purchase Order, Sales Order atau Final Order  
![][image14]

* purchase order   
- update di purchase order, sales order , dan final   
* sales order   
- update di sales order & final   
* final order   
- update di final 

### URL : 

| https://be.scyllax.online/sales/v1/orders/SO2509230001  |
| :---- |

**Method** : PATCH

Pada case ini (edit qty, tanpa ubah data, dan add sku) jika sls.order.data\_status sebelumnya \= Need Review (1) maka update sls.order.data\_status menjadi Processed (2)

#### **Case 1 \- Edit (Tab Purchase Order)**

##### **NO EDIT PRODUK** 

Frontend : 

* **Before** : FE kirim semua payload  
* **After**     :   
  * jika pada detail opr\_type \= C FE tidak perlu mengimkan data apapun , karena tdk ada data yang di edit dari sisi FE (sehingga BE tidak ada update stock dan on cust order) 

  * jika pada detail opr\_type \= O FE perlu mengimkan data order , karena tdk ada data yang di edit dari sisi FE (sehingga BE ada update stock dan on cust order) 

  Handle BE jika FE **TIDAK KIRIM** data ke BE 

* BE update 1 data saja yaitu :  data\_status  \= 2  
* **Before :**   
  BE insert ke inv.stock 

  **After :**   
  BE tidak perlu insert ke inv.stock karena tidak ada perubahan qty   
* **Before  :**   
  BE update ke inv.warehouse\_stock 

  **After :**   
  BE tidak perlu insert ke inv.warehouse\_stock karena tidak ada perubahan qty 

##### **EDIT PRODUK** 

Untuk data yang berada pada purchase\_order, **ketika ada edit produk,**maka be update **sls.order\_detail** berdasarkan order\_detail\_id pada field berikut :

| Request | Update  |
| :---- | :---- |
| purchase\_details.qty\_po1 | qty\_po1 |
| purchase\_details.qty\_po2 | qty\_po2 |
| purchase\_details.qty\_po3 | qty\_po3 |
| purchase\_details.sell\_price\_po1 | sell\_price\_po1 |
| purchase\_details.sell\_price\_po2 | sell\_price\_po2 |
| purchase\_details.sell\_price\_po3 | sell\_price\_po3 |
|  | **Discount :** Karena ada perubahan gross total, terdapat perubahan discount  Perhitungan eksisting BE : **ConsultDiscountBeforeStore** Update : *(dg nilai yg sama)***disc\_value** (untuk tab sales order) **disc\_value\_final** (untuk tab final order)  |
|  | **VAT Value :** Karena ada perubahan gross total, terdapat perubahan ppn  Perhitungan :  **((qty1 \* sell\_price1) \+ (qty2\* sell\_price2) \+ (qty3 \* sell\_price3)-promo\_value\_final \-disc\_value\_final) \* vat%**Update : *(dg nilai yg sama)* **vat\_value** (untuk tab sales order)  **vat\_value\_final** (untuk tab final order)   |

**Update di inv.stock**   
  (untuk saat ini insert double atau 4 data , seharusnya hanya 2 data aja SOXXX . dan SOXXX-CO  
**Update Warehouse Stock & On Customer Order**  
**Tabel : inv.warehouse\_stock ws**  

* qty\_update	\=  conversi satuan terkecil dari selisih (qty1, qty2, qty3) sebelum di edit \-  (qty1, qty2, qty3) request FE   
* Update ws.qty  \= ws.qty \- **qty\_update**  
* Update ws.qty\_on\_order  \= ws.qty \+ **qty\_update**  
* **Contoh 1 :** sebelum di edit, user input QTY 10 , user melakukan edit QTY 15   
  qty\_update \= \+5 (bertambah 5\)   
  maka ws.qty harus berkurang 5 dan ws.qty\_on\_order bertambah 5   
* **Contoh 2 :** sebelum di edit, user input QTY 10 , user melakukan edit QTY 5   
  qty\_update \= \-5 (berkurang 5\)   
  maka ws.qty harus bertambah 5 dan ws.qty\_on\_order bertambah 5 

Contoh Request : 

| { "purchase\_details": \[ { "order\_detail\_id": 6036, "qty\_po1": 5, "qty\_po2": 0, "qty\_po3": 0, "sell\_price\_po1": 2750, "sell\_price\_po2": 165000, "sell\_price\_po3": 165000, "is\_produc } \] } |
| :---- |

##### **ADD PRODUK** 

Untuk data yang berada pada purchase\_order, **ketika ada add produk,** maka be add **sls.order\_detail** dengan data sebagai berikut : 

| Request | Data Source  |
| :---- | :---- |
| pro\_id | pro\_id |
| sell\_price\_system1 | request FE : sell\_price\_system1 |
| sell\_price\_system2 | request FE : sell\_price\_system2 |
| sell\_price\_system3 | request FE : sell\_price\_system3 |
| purchase\_details.qty\_po1 | request FE : qty\_po1 |
| purchase\_details.qty\_po2 | request FE : qty\_po2 |
| purchase\_details.qty\_po3 | request FE : qty\_po3 |
| unit\_id1 | unit\_id1 |
| unit\_id2 | unit\_id2 |
| unit\_id3 | unit\_id3 |
| qty1\_stock | Fe Ambil dari initial stock sebelum di process order [https://prnt.sc/\_9vg0OcWdZRW](https://prnt.sc/_9vg0OcWdZRW) atau saat hit  API [https://best.scyllax.online/inventory/v1/stocks/report?page=1\&limit=10\&q=\&wh\_id=63\&pro\_id=676\&outlet\_id=1387\&order\_date=2026-01-13\&include\_zero\_stock=true\&active\_product\_only=true\&sort=pro\_code:asc](https://best.scyllax.online/inventory/v1/stocks/report?page=1&limit=10&q=&wh_id=63&pro_id=676&outlet_id=1387&order_date=2026-01-13&include_zero_stock=true&active_product_only=true&sort=pro_code:asc)kirim dr response \= qty1, qty2, qty3  |
| qty2\_stock |  |
| qty3\_stock |  |
| purchase\_details.sell\_price\_po1 | request FE : sell\_price\_po1 |
| purchase\_details.sell\_price\_po2 | request FE : sell\_price\_po2 |
| purchase\_details.sell\_price\_po3 | request FE : sell\_price\_po3 |
| purchase\_details.qty1 | request FE : qty\_po1 |
| purchase\_details.qty2 | request FE : qty\_po2 |
| purchase\_details.qty3 | request FE : qty\_po3 |
| purchase\_details.sell\_price1 | request FE : sell\_price\_po1 |
| purchase\_details.sell\_price2 | request FE : sell\_price\_po2 |
| purchase\_details.sell\_price3 | request FE : sell\_price\_po3 |
| purchase\_details.qty1\_final | request FE : qty\_po1 |
| purchase\_details.qty2\_final | request FE : qty\_po2 |
| purchase\_details.qty3\_final | request FE : qty\_po3 |
| purchase\_details.sell\_price\_final1 | request FE : sell\_price\_po1 |
| purchase\_details.sell\_price\_final2 | request FE : sell\_price\_po2 |
| purchase\_details.sell\_price\_final3 | request FE : sell\_price\_po3 |

* Contoh Request update qty\_po1, qty\_po2, qty\_po3, sell\_price\_po1, sell\_price\_po2, sell\_price\_po3 : 

{  
 "purchase\_details" :  
 \[  
   {  
     "order\_detail\_id" : 323,  
     "qty\_po1" : 21,  
     "qty\_po2" : 22,  
     "qty\_po3" : 0,  
     "sell\_price\_po1" : 21000,  
     "sell\_price\_po2" : 23000,  
     "sell\_price\_po3" : 21000  
   }  
 \]  
}

* Contoh Request update qty\_po1 dan sell\_price\_po1 : 

{  
 "purchase\_details" :  
 \[  
   {  
     "order\_detail\_id" : 323,  
     "qty\_po1" : 21,  
     "sell\_price\_po1" : 21000  
   }  
 \]  
}

* Contoh Request update qty\_po1 dan sell\_price\_po1, dan delete order\_detail 222 

{  
"purchase\_details" :  
\[  
  {  
    "order\_detail\_id" : 323,  
    "qty\_po1" : 21,  
    "sell\_price\_po1" : 21000  
  },  
 {  
    "order\_detail\_id" : 222,  
    "qty\_po1" : 0,  
    "qty\_po2" : 0,  
    "qty\_po3" : 0,  
  }  
\],

* Contoh Request update qty\_po1 dan sell\_price\_po1, delete order\_detail 222, dan **add produk 123** 

{  
 "purchase\_details": \[  
   {  
     "order\_detail\_id": 323,  
     "qty\_po1": 21,  
     "sell\_price\_po1": 21000  
   },  
   {  
     "order\_detail\_id": 222,  
     "qty\_po1": 0,  
     "qty\_po2": 0,  
     "qty\_po3": 0  
   }  
 \],  
 "add\_purchase\_details": \[  
   {  
     "pro\_id" :121,  
     "qty\_po1": 21,  
     "qty\_po2": 21,  
     "qty\_po3": 21,  
     "sell\_price\_system1": 21000,  
     "sell\_price\_system2": 21000,  
     "sell\_price\_system3": 21000,  
     "sell\_price\_po1": 25000,  
     "sell\_price\_po2": 26000,  
     "sell\_price\_po3": 27000,  
     "unit\_id1": "PCS",  
     "unit\_id2": "PCS",  
     "unit\_id3": "PCS",   
     "is\_product\_promotion\_so": false  
   },  
   {  
     "pro\_id": 122,  
     "qty\_po1": 21,  
     "qty\_po2": 0,  
     "qty\_po3": 0,  
     "sell\_price\_system1": 21000,  
     "sell\_price\_system2": 21000,  
     "sell\_price\_system3": 21000,  
     "sell\_price\_po1": 21000,  
     "sell\_price\_po2": 21000,  
     "sell\_price\_po3": 21000,  
     "unit\_id1": "PCS",  
     "unit\_id2": "PCS",  
     "unit\_id3": "PCS",  
 "is\_product\_promotion\_so": false  
   }  
 \]  
}

**Update di inv.stock**  
**Update Warehouse Stock & On Customer Order**  
**Tabel : inv.warehouse\_stock ws**  

* qty\_update	\=  conversi satuan terkecil dari selisih (qty1, qty2, qty3) sebelum di edit \-  (qty1, qty2, qty3) request FE   
* Update ws.qty  \= ws.qty \- **qty\_update**  
* Update ws.qty\_on\_order  \= ws.qty \+ **qty\_update**  
* **Contoh 1 :** sebelum di edit, user input QTY 10 , user melakukan edit QTY 15   
  qty\_update \= \+5 (bertambah 5\)   
  maka ws.qty harus berkurang 5 dan ws.qty\_on\_order bertambah 5   
* **Contoh 2 :** sebelum di edit, user input QTY 10 , user melakukan edit QTY 5   
  qty\_update \= \-5 (berkurang 5\)   
  maka ws.qty harus bertambah 5 dan ws.qty\_on\_order bertambah 5 

#### 

#### 

#### **Case 2 \- (Tab Sales Order)**

→ jika pada detail, is\_performa\_inv \= false / null, ketika user edit pada sales order, maka final order juga akan berubah mengikuti sales order.   
Untuk data yang berada pada sales\_order, maka be update **sls.order\_detail** berdasarkan order\_detail\_id pada field berikut : 

##### **EDIT PRODUK  & DELETE**  **Tabel : sls.order\_detail**

| Request | Update  |
| :---- | :---- |
| details.qty1 | qty1, qty1\_final |
| details.qty2 | qty2, qty2\_final |
| details.qty3 | qty3, qty3\_final |
| details.sell\_price1 | sell\_price1, sell\_price\_final1 |
| details.sell\_price2 | sell\_price2, sell\_price\_final2 |
| details.sell\_price3 | sell\_price3, sell\_price\_final3 |
|  | **Discount :** Karena ada perubahan gross total, terdapat perubahan discount  Perhitungan eksisting BE : **ConsultDiscountBeforeStore** Update : *(dg nilai yg sama)***disc\_value** (untuk tab sales order) **disc\_value\_final** (untuk tab final order)  |
|  | **VAT Value :** Karena ada perubahan gross total, terdapat perubahan ppn  Perhitungan :  **((qty1 \* sell\_price1) \+ (qty2\* sell\_price2) \+ (qty3 \* sell\_price3)-promo\_value\_final \-disc\_value\_final) \* vat%**Update : *(dg nilai yg sama)* **vat\_value** (untuk tab sales order)  **vat\_value\_final** (untuk tab final order)   |

**Update di inv.stock**  
**Update Warehouse Stock & On Customer Order**  
**Tabel : inv.warehouse\_stock ws**  

* qty\_update	\=  conversi satuan terkecil dari selisih (qty1, qty2, qty3) sebelum di edit \-  (qty1, qty2, qty3) request FE   
* Update ws.qty  \= ws.qty \- **qty\_update**  
* Update ws.qty\_on\_order  \= ws.qty \+ **qty\_update**  
* **Contoh 1 :** sebelum di edit, user input QTY 10 , user melakukan edit QTY 15   
  qty\_update \= \+5 (bertambah 5\)   
  maka ws.qty harus berkurang 5 dan ws.qty\_on\_order bertambah 5   
* **Contoh 2 :** sebelum di edit, user input QTY 10 , user melakukan edit QTY 5   
  qty\_update \= \-5 (berkurang 5\)   
  maka ws.qty harus bertambah 5 dan ws.qty\_on\_order bertambah 5 

##### **ADD SKU** 

Untuk data yang berada pada sales\_order, **ketika ada add produk,** maka be add **sls.order\_detail** dengan data sebagai berikut : 

| Request | Add  |
| :---- | :---- |
| pro\_id | pro\_id |
| sell\_price\_system1 | sell\_price\_system1 |
| sell\_price\_system2 | sell\_price\_system2 |
| sell\_price\_system3 | sell\_price\_system3 |
| qty1 | qty1 , qty1\_final |
| qty2 | qty2, qty2\_final |
| qty3 | qty3, qty3\_final |
| unit\_id1 | unit\_id1 |
| unit\_id2 | unit\_id2 |
| unit\_id3 | unit\_id3 |
| qty1\_stock | Fe Ambil dari initial stock sebelum di process order [https://prnt.sc/\_9vg0OcWdZRW](https://prnt.sc/_9vg0OcWdZRW) atau saat hit  API [https://best.scyllax.online/inventory/v1/stocks/report?page=1\&limit=10\&q=\&wh\_id=63\&pro\_id=676\&outlet\_id=1387\&order\_date=2026-01-13\&include\_zero\_stock=true\&active\_product\_only=true\&sort=pro\_code:asc](https://best.scyllax.online/inventory/v1/stocks/report?page=1&limit=10&q=&wh_id=63&pro_id=676&outlet_id=1387&order_date=2026-01-13&include_zero_stock=true&active_product_only=true&sort=pro_code:asc)kirim dr response \= qty1, qty2, qty3  |
| qty2\_stock |  |
| qty3\_stock |  |
| unit\_id1 | unit\_id1 |
| unit\_id2 | unit\_id2 |
| unit\_id3 | unit\_id3 |
| qty1\_stock | qty1\_stock |
| qty2\_stock | qty2\_stock |
| qty3\_stock | qty3\_stock |
| sell\_price1 | sell\_price1, sell\_price\_final1 |
| sell\_price2 | sell\_price2, sell\_price\_final2 |
| sell\_price3 | sell\_price3, sell\_price\_final3 |
|  | **Discount** Perhitungan eksisting BE : **ConsultDiscountBeforeStore** Update : *(dg nilai yg sama)***disc\_value** (untuk tab sales order) **disc\_value\_final** (untuk tab final order)  |
|  | **VAT\_VALUE** Perhitungan BE:  **((qty1 \* sell\_price1) \+ (qty2\* sell\_price2) \+ (qty3 \* sell\_price3)-promo\_value\_final \-disc\_value\_final) \* vat%**Update : *(dg nilai yg sama)* **vat\_value** (untuk tab sales order)  **vat\_value\_final** (untuk tab final order)   |
|  | **Initial Stock** *(disamakan dg eksisting seperti API /sales/v1/orders/)* qty1\_stock qty2\_stock  qty3\_stock  |

**Update di inv.stock**  
**Update Warehouse Stock & On Customer Order**  
**Tabel : inv.warehouse\_stock ws**  

* qty\_update	\=  conversi satuan terkecil dari selisih (qty1, qty2, qty3) sebelum di edit \-  (qty1, qty2, qty3) request FE   
* Update ws.qty  \= ws.qty \- **qty\_update**  
* Update ws.qty\_on\_order  \= ws.qty \+ **qty\_update**  
* **Contoh 1 :** sebelum di edit, user input QTY 10 , user melakukan edit QTY 15   
  qty\_update \= \+5 (bertambah 5\)   
  maka ws.qty harus berkurang 5 dan ws.qty\_on\_order bertambah 5   
* **Contoh 2 :** sebelum di edit, user input QTY 10 , user melakukan edit QTY 5   
  qty\_update \= \-5 (berkurang 5\)   
  maka ws.qty harus bertambah 5 dan ws.qty\_on\_order bertambah 5 

* Contoh Request update qty1, qty2, qty3, sell\_price1, sell\_price2, sell\_price3 : 

{  
 "sales\_order" : \[  
   {  
     "order\_detail\_id" : 323,  
     "qty1" : 21,  
     "qty2" : 22,  
     "qty3" : 0,  
     "sell\_price1" : 21000,  
     "sell\_price2" : 23000,  
     "sell\_price3" : 21000  
   }  
 \]  
}

* Contoh Request update qty1, sell\_price1: 

{  
 "sales\_order" : \[  
   {  
     "order\_detail\_id" : 323,  
     "qty1" : 21,  
     "sell\_price1" : 21000  
   }  
 \]  
}

* Contoh Request update qty1, sell\_price1, dan delete order\_detail 222 dan 223: 

{  
 "sales\_order" : \[  
   {  
     "order\_detail\_id" : 323,  
     "qty1" : 21,  
     "sell\_price1" : 21000  
   }  
 \]  
}

* Contoh Request update oder\_detail\_id \= 323, dan penambahan produk pro\_id= 121 dan 122

{  
"sales\_order" : \[  
  {  
    "order\_detail\_id" : 323,  
    "qty1" : 21,  
    "qty2" : 22,  
    "qty3" : 0,  
    "sell\_price1" : 21000,  
    "sell\_price2" : 23000,  
    "sell\_price3" : 21000  
  }  
\],  
"add\_sales\_order": \[  
   {  
     "pro\_id" :121,  
     "qty1": 21,  
     "qty2": 21,  
     "qty3": 21,  
     "sell\_price\_system1": 21000,  
     "sell\_price\_system2": 21000,  
     "sell\_price\_system3": 21000,  
     "sell\_price1": 25000,  
     "sell\_price2": 26000,  
     "sell\_price3": 27000  
   },  
   {  
     "pro\_id": 122,  
     "qty1": 21,  
     "qty2": 21,  
     "qty3": 21,  
     "sell\_price\_system1": 21000,  
     "sell\_price\_system2": 21000,  
     "sell\_price\_system3": 21000,  
     "sell\_price1": 25000,  
     "sell\_price2": 26000,  
     "sell\_price3": 27000  
   }  
 \]  
}

#### 

#### **Case 3 \- (Tab Final Order)**

##### **EDIT PRODUK  & DELETE**

Untuk data yang berada pada final\_order, maka be update **sls.order\_detail** berdasarkan order\_detail\_id pada field berikut : 

| Request | Update Tabel  |
| :---- | :---- |
| final\_order.qty1\_final | qty1\_final |
| final\_order.qty2\_final | qty2\_final |
| final\_order.qty3\_final | qty3\_final |
| final\_order.sell\_price\_final1 | sell\_price\_final1 |
| final\_order.sell\_price\_final2 | sell\_price\_final2 |
| final\_order.sell\_price\_final3 | sell\_price\_final3 |

**Update di inv.stock**  
**Update Warehouse Stock & On Customer Order**  
**Tabel : inv.warehouse\_stock ws**  

* qty\_update	\=  conversi satuan terkecil dari selisih (qty1, qty2, qty3) sebelum di edit \-  (qty1, qty2, qty3) request FE   
* Update ws.qty  \= ws.qty \- **qty\_update**  
* Update ws.qty\_on\_order  \= ws.qty \+ **qty\_update**  
* **Contoh 1 :** sebelum di edit, user input QTY 10 , user melakukan edit QTY 15   
  qty\_update \= \+5 (bertambah 5\)   
  maka ws.qty harus berkurang 5 dan ws.qty\_on\_order bertambah 5   
* **Contoh 2 :** sebelum di edit, user input QTY 10 , user melakukan edit QTY 5   
  qty\_update \= \-5 (berkurang 5\)   
  maka ws.qty harus bertambah 5 dan ws.qty\_on\_order bertambah 5 

##### **ADD SKU** 

Untuk data yang berada pada final\_order, **ketika ada add produk,** maka be add **sls.order\_detail** dengan data sebagai berikut : 

| Request | Data Source  |
| :---- | :---- |
| sell\_price\_system1 | request FE : sell\_price\_system1 |
| sell\_price\_system2 | request FE : sell\_price\_system2 |
| sell\_price\_system3 | request FE : sell\_price\_system3 |
| unit\_id1 | unit\_id1 |
| unit\_id2 | unit\_id2 |
| unit\_id3 | unit\_id3 |
| qty1\_stock | qty1\_stock |
| qty2\_stock | qty2\_stock |
| qty3\_stock | qty3\_stock |
| purchase\_details.qty1\_final | request FE : qty1\_final |
| purchase\_details.qty2\_final | request FE : qty2\_final |
| purchase\_details.qty3\_final | request FE : qty3\_final |
| purchase\_details.sell\_price\_final1 | request FE : sell\_price\_final1 |
| purchase\_details.sell\_price\_final2 | request FE : sell\_price\_final2 |
| purchase\_details.sell\_price\_final3 | request FE : sell\_price\_final3 |
|  | **Discount** Perhitungan eksisting BE : **ConsultDiscountBeforeStore** Update : *(dg nilai yg sama)***disc\_value\_final** (untuk tab final order)  |
|  | **VAT\_VALUE** Perhitungan BE:  **((qty1 \* sell\_price1) \+ (qty2\* sell\_price2) \+ (qty3 \* sell\_price3)-promo\_value\_final \-disc\_value\_final) \* vat%**Update : *(dg nilai yg sama)* **vat\_value\_final** (untuk tab final order)  |
|  | **Initial Stock** *(disamakan dg eksisting seperti API /sales/v1/orders/)* qty1\_stock qty2\_stock  qty3\_stock  |

**Update di inv.stock**  
**Update Warehouse Stock & On Customer Order**  
**Tabel : inv.warehouse\_stock ws**  

* qty\_update	\=  conversi satuan terkecil dari selisih (qty1, qty2, qty3) sebelum di edit \-  (qty1, qty2, qty3) request FE   
* Update ws.qty  \= ws.qty \- **qty\_update**  
* Update ws.qty\_on\_order  \= ws.qty \+ **qty\_update**  
* **Contoh 1 :** sebelum di edit, user input QTY 10 , user melakukan edit QTY 15   
  qty\_update \= \+5 (bertambah 5\)   
  maka ws.qty harus berkurang 5 dan ws.qty\_on\_order bertambah 5   
* **Contoh 2 :** sebelum di edit, user input QTY 10 , user melakukan edit QTY 5   
  qty\_update \= \-5 (berkurang 5\)   
  maka ws.qty harus bertambah 5 dan ws.qty\_on\_order bertambah 5   
    
    
* Contoh Request update qty1\_final, qty2\_final, qty3\_final, sell\_price\_final1, sell\_price\_final2, sell\_price\_final3 : 

{  
 {  
 "ro\_no": "SO2512150001",  
 "final\_order" : \[  
   {  
     "order\_detail\_id" : 323,  
     "qty1\_final" : 21,  
     "qty2\_final" : 22,  
     "qty3\_final" : 0,  
     "sell\_price\_final1" : 21000,  
     "sell\_price\_final2" : 23000,  
     "sell\_price\_final3" : 21000,  
   }  
 \]  
}  
 

* Contoh Request update qty1\_final, sell\_price\_final1 : {

{  
 "ro\_no": "SO2512150001",  
 "final\_order" : \[  
   {  
     "order\_detail\_id" : 323,  
     "qty1\_final" : 21,  
     "sell\_price\_final1" : 21000  
   }  
 \]  
}

* Contoh Request update order\_detail\_id \= 323 dan add\_final\_order pro\_id \= 121 dan 122

{  
"ro\_no": "SO2512150001",  
"final\_order" : \[  
  {  
    "order\_detail\_id" : 323,  
    "qty1\_final" : 21,  
    "qty2\_final" : 22,  
    "qty3\_final" : 0,  
    "sell\_price\_final1" : 21000,  
    "sell\_price\_final2" : 23000,  
    "sell\_price\_final3" : 21000,  
  }  
\],  
 "add\_final\_order": \[  
   {  
     "pro\_id" :121,  
     "qty1\_final": 21,  
     "qty2\_final": 21,  
     "qty3\_final": 21,  
     "sell\_price\_system1": 21000,  
     "sell\_price\_system2": 21000,  
     "sell\_price\_system3": 21000,  
     "sell\_price\_final1": 25000,  
     "sell\_price\_final2": 26000,  
     "sell\_price\_final3": 27000  
   },  
   {  
     "pro\_id": 122,  
     "qty1\_final": 21,  
     "qty2\_final": 21,  
     "qty3\_final": 21,  
     "sell\_price\_system1": 21000,  
     "sell\_price\_system2": 21000,  
     "sell\_price\_system3": 21000,  
     "sell\_price\_final1": 25000,  
     "sell\_price\_final2": 26000,  
     "sell\_price\_final3": 27000  
   }  
 \]  
}

### Response 

{  
 "message": "Berhasil Diperbarui",  
 "request\_id": "999cb52d-adf7-47c4-b10b-304a3b6e553a"  
}

12. ## **Download Excel Sales Order**

![][image15]

 Create new endpoint :  
Content-Type 	: application/json  
Method		: GET  
URL		: https://be.scyllax.online/sales/v1/download

### **Headers**

---

| Accept | application/json |
| :---- | :---- |
| Authorization | Bearer Token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4 |

### **Params** 

---

| Field | Type | Required | Note |
| :---- | :---- | :---- | :---- |
| start\_date | date | Yes (max 31 days) | epoch |
| end\_date  | date | Yes  (max 31 days) | epoch |
| salesman\_id | Integer | No | Untuk melakukan filter berdasarkan id supplier  |

### **Example Request** 

---

**Example Request Default**

| curl \--location \-g \\ 'https://be.scyllax.online/sales/v1/download?start\_date=1701388800\&end\_date=1704067199\&salesman\_id=12' \\ \--header 'Accept: application/json' |
| :---- |

### **Response :** 

* Sukses generate 

**{**

   **"message": "",**

   **"data": {**

       **"report\_id": "69424720904ba4ec4bd832d4",**

       **"report\_name": "DownloadSalesOrder-171225-001",**

       **"start\_date": "2025-12-17",**

       **"end\_date": "2025-12-17",**

       **"file\_status": 5,**

       **"file\_status\_name": "",**

       **"file\_url": "",**

       **"created\_by": "Admin DT Bangka 01",**

       **"created\_at": "2025-12-17T06:01:04.405513856Z"**

   **},**

   **"paging": {**

       **"total\_record": 1,**

       **"page\_current": 0,**

       **"page\_limit": 0,**

       **"page\_total": 1**

   **},**

   **"request\_id": "9e79811e-4214-4567-94fe-84d8f37ef172"**

**}**

* Gagal generate 

cek pada report.list jika terdapat report\_name dengan prefix DownloadSalesOrder status \= inprogress → maka BE mengembalikan resp gagal   
{  
   "message": "Processing time may vary by file size. Please check Download History to access the file",  
   "data": null ,  
   "request\_id": "9e79811e-4214-4567-94fe-84d8f37ef172"  
}

* **Format Generate Name File :** DownloadSalesOrder-DDMMYY-3digitRunningNumber  
* **Proses no 3 : (insert data → report.list**

| Field | Data Source |
| :---- | :---- |
| cust\_id | Diisi cust id user yang login |
| report\_id |  |
| report\_name | Diisi dengan format : DownloadSalesOrder-DDMMYY-3digitRunningNumber |
| start\_date | Diisi Tanggal sekarang  |
| end\_date | start\_date \+ 30  |
| file\_status | 0 jika proses generate selesai , update value 1  |
| file\_url  | NULL  |
| file\_base64 | add new field  diisi file base64 hasil generate be , di proses ini masih kosong karena masih proses generate  jika proses generate selesai , update value base64 |

Format Excel : [https://docs.google.com/spreadsheets/d/17LCTqne6rXv4x00bpFtEBUIrzReM3rcLq6rwTwHDxXI/edit?gid=236671606\#gid=236671606](https://docs.google.com/spreadsheets/d/17LCTqne6rXv4x00bpFtEBUIrzReM3rcLq6rwTwHDxXI/edit?gid=236671606#gid=236671606) 

* #### **Sheet Purchase Order** 

| Field | Data Source |
| :---- | :---- |
| Po No  | **sls.order.**po\_no |
| So No  | **sls.order.**ro\_no |
| Order Date  | **sls.order**.ro\_date |
| Outlet Code  | **mst.outlet.m\_outlet.**outlet\_code |
| Outlet Name  | **mst.m\_outlet.**outlet\_name |
| Salesman Code  | **mst.m\_salesman.**emp.id |
| Employee Name | **mst.m\_salesman.**sales\_name |
| Supplier Code  | **mst.m\_supplier.**sup\_code relasi dengan tabel produk  |
| Supplier Name  | **mst.m\_supplier.**sup\_name relasi dengan tabel produk  |
| Product Code  | **mst.m\_product**.pro\_code |
| ProName | **mst.m\_product**.pro\_name |
| Largest Unit  | **sls.order\_detail.**unit\_id3 / conv\_unit2 |
| Middle Unit  | **sls.order\_detail.**unit\_id2 / conv\_unit3 |
| Smalles Unit  | **sls.order\_detail.**unit\_id1 /1 |
| Largest Selling Price | **sls.order\_detail.**sell\_price\_system3 |
| Middle Selling Price | **sls.order\_detail.**sell\_price\_system2 |
| Smallest Selling Price | **sls.order\_detail.**sell\_price\_system1 |
| Final Largest Selling Price | **sls.order\_detail.**sell\_price\_po3 |
| Final Middle Selling Price | **sls.order\_detail.**sell\_price\_po2 |
| FInal Smallest Selling Price | **sls.order\_detail.**sell\_price\_po1 |
| Largest QTY Order | **sls.order\_detail.**qty\_po3 |
| Middle QTY Order | **sls.order\_detail.**qty\_po2 |
| Smallest QTY Order | **sls.order\_detail.**qty\_po1 |
| GrossSales | (sell\_price\_po1\* qty\_po1) \+ (sell\_price\_po2\* qty\_po2) \+ (sell\_price\_po3\* qty\_po3) |
| Promotion | **sls.order\_detail.**vat\_value\_final |
| Discount | **sls.order\_detail.**disc\_value\_final |
| Net Sales (ExcPPN) | gross\_sales \-promotion-discount |
| PPN  | net\_sales \*vat |
| Gross | Net Sales  \+ vat |

* #### **Sheet Sales Order**

| Field | Data Source |
| :---- | :---- |
| Po No  | **sls.order.**po\_no |
| So No  | **sls.order.**ro\_no |
| Order Date  | **sls.order**.ro\_date |
| Outlet Code  | **mst.outlet.m\_outlet.**outlet\_code |
| Outlet Name  | **mst.m\_outlet.**outlet\_name |
| Salesman Code  | **mst.m\_salesman.**emp.id |
| Employee Name | **mst.m\_salesman.**sales\_name |
| Supplier Code  | **mst.m\_supplier.**sup\_code relasi dengan tabel produk  |
| Supplier Name  | **mst.m\_supplier.**sup\_name relasi dengan tabel produk  |
| Product Code  | **mst.m\_product**.pro\_code |
| ProName | **mst.m\_product**.pro\_name |
| Largest Unit  | **sls.order\_detail.**unit\_id3 / conv\_unit2 |
| Middle Unit  | **sls.order\_detail.**unit\_id2 / conv\_unit3 |
| Smalles Unit  | **sls.order\_detail.**unit\_id1 /1 |
| Largest Selling Price | **sls.order\_detail.**sell\_price\_system1 |
| Middle Selling Price | **sls.order\_detail.**sell\_price\_system2 |
| Smallest Selling Price | **sls.order\_detail.**sell\_price\_system3 |
| Final Largest Selling Price | **sls.order\_detail.**sell\_price3 |
| Final Middle Selling Price | **sls.order\_detail.**sell\_price2 |
| FInal Smallest Selling Price | **sls.order\_detail.**sell\_price1 |
| Largest QTY Order | **sls.order\_detail.**qty3 |
| Middle QTY Order | **sls.order\_detail.**qty2 |
| Smallest QTY Order | **sls.order\_detail.**qty1 |
| GrossSales | (sell\_price1\* qty1) \+ (sell\_price2\* qty2) \+ (sell\_price3\* qty3) |
| Promotion | **sls.order\_detail.**vat\_value\_final |
| Discount | **sls.order\_detail.**disc\_value\_final |
| Net Sales (ExcPPN) | gross\_sales \-promotion-discount |
| PPN  | net\_sales \*vat |
| Gross | Net Sales  \+ vat |

* #### **Sheet Final Order**

| Field | Data Source |
| :---- | :---- |
| Po No  | **sls.order.**po\_no |
| So No  | **sls.order.**ro\_no |
| Order Date  | **sls.order**.ro\_date |
| Outlet Code  | **mst.outlet.m\_outlet.**outlet\_code |
| Outlet Name  | **mst.m\_outlet.**outlet\_name |
| Salesman Code  | **mst.m\_salesman.**emp.id |
| Employee Name | **mst.m\_salesman.**sales\_name |
| Supplier Code  | **mst.m\_supplier.**sup\_code relasi dengan tabel produk  |
| Supplier Name  | **mst.m\_supplier.**sup\_name relasi dengan tabel produk  |
| Product Code  | **mst.m\_product**.pro\_code |
| ProName | **mst.m\_product**.pro\_name |
| Largest Unit  | **sls.order\_detail.**unit\_id3 / conv\_unit2 |
| Middle Unit  | **sls.order\_detail.**unit\_id2 / conv\_unit3 |
| Smalles Unit  | **sls.order\_detail.**unit\_id1 /1 |
| Largest Selling Price | **sls.order\_detail.**sell\_price\_system1 |
| Middle Selling Price | **sls.order\_detail.**sell\_price\_system2 |
| Smallest Selling Price | **sls.order\_detail.**sell\_price\_system3 |
| Final Largest Selling Price | **sls.order\_detail.**sell\_price\_final3 |
| Final Middle Selling Price | **sls.order\_detail.**sell\_price\_final2 |
| FInal Smallest Selling Price | **sls.order\_detail.**sell\_price\_final1 |
| Largest QTY Order | **sls.order\_detail.**qty3\_final |
| Middle QTY Order | **sls.order\_detail.**qty2\_final |
| Smallest QTY Order | **sls.order\_detail.**qty1\_final |
| GrossSales | (sell\_price\_final1\* qty1\_final) \+ (sell\_price\_final2\* qty2\_final) \+ (sell\_price\_final3\* qty3\_final) |
| Promotion | **sls.order\_detail.**vat\_value\_final |
| Discount | **sls.order\_detail.**disc\_value\_final |
| Net Sales (ExcPPN) | gross\_sales \-promotion-discount |
| PPN  | net\_sales \*vat |
| Gross | Net Sales  \+ vat |

* #### **Sheet QTY Summary** 

| Field | Data Source |
| :---- | :---- |
| Po No  | **sls.order.**po\_no |
| So No  | **sls.order.**ro\_no |
| Order Date  | **sls.order**.ro\_date |
| Outlet Code  | **mst.outlet.m\_outlet.**outlet\_code |
| Outlet Name  | **mst.m\_outlet.**outlet\_name |
| Salesman Code  | **mst.m\_salesman.**emp.id |
| Employee Name | **mst.m\_salesman.**sales\_name |
| Supplier Code  | **mst.m\_supplier.**sup\_code relasi dengan tabel produk  |
| Supplier Name  | **mst.m\_supplier.**sup\_name relasi dengan tabel produk  |
| Product Code  | **mst.m\_product**.pro\_code |
| ProName | **mst.m\_product**.pro\_name |
| Largest Unit  | **sls.order\_detail.**unit\_id3 |
| Middle Unit  | **sls.order\_detail.**unit\_id2 |
| Smalles Unit  | **sls.order\_detail.**unit\_id1 |
| Largest QTY Purchase Order | **sls.order\_detail.**qty\_po3 |
| Middle QTY Purchase Order | **sls.order\_detail.**qty\_po2 |
| Smallest QTY Purchase Order | **sls.order\_detail.**qty\_po1 |
| Largest QTY Sales  Order | **sls.order\_detail.**qty3 |
| Middle QTY Sales  Order | **sls.order\_detail.**qty2 |
| Smallest QTY Sales  Order | **sls.order\_detail.**qty1 |
| Largest QTY Final Order | **sls.order\_detail.**qty3\_final |
| Middle QTY Final Order | **sls.order\_detail.**qty2\_final |
| Middle QTY Final Order | **sls.order\_detail.**qty1\_final |

* Pada proses no 8 , enhance API 

https://best.scyllax.online/sales/v1/reports?page=1\&limit=10  
**add response : file\_base64** 

| Field | Data Source |
| :---- | :---- |
| cust\_id | Diisi cust id user yang login |
| report\_id |  |
| report\_name | Diisi dengan format : DownloadSalesOrder-DDMMYY-3digitRunningNumber |
| start\_date | Diisi Tanggal sekarang  |
| end\_date | start\_date \+ 30  |
| file\_status | 1 |
| file\_url  | NULL  |
| file\_base64 | diisi file base64 hasil generate be  |

## **Cancel Order** {#cancel-order}

**https://best.scyllax.online/sales/v1/orders/status**

| curl 'https://best.scyllax.online/sales/v1/orders/status' \\   \-X 'PATCH' \\   \-H 'sec-ch-ua-platform: "macOS"' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsImVtcF9pZCI6MCwiZXhwaXJlcyI6MTc3MDY4NDc4NSwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4NTc0ODEyMzEyOCIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IkRpc3QgSURFIFNkYSIsInVzZXJfaWQiOjEyLCJ1c2VyX25hbWUiOiJkaXN0QHNkYS5pZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODU3NDgxMjMxMjgifQ.A1fCYdI1juSyKF7nWNLnKl1wqTzGwb-mPLOnKiWs6TE' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'sec-ch-ua: "Not(A:Brand";v="8", "Chromium";v="144", "Google Chrome";v="144"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36' \\   \-H 'Accept: application/json, text/plain, /' \\   \-H 'Content-Type: application/json' \\   \--data-raw '{"orders":\[{"ro\_no":"SO2602060007","data\_status":9}\]}'  |
| :---- |

| select (sum(*s*.qty\_in ) \- sum(*s*.qty\_out )) as *qty*, (sum(*s*.qty\_in\_order) \- sum(*s*.qty\_out\_order )) as *qty\_in\_order*  from inv.stock *s* where *s*.pro\_id \=474 and *s*.wh\_id \=63 and *s*.cust\_id \='C220010001'  ![][image16] ![][image17] saldo awal :  qty \= 30388 conv\_unit2 \=24 conv\_unit1 \= 1  Large     \=  1266 Middle  \= 0 Small     \= 4 Sales Order (Processing) select *s*.qty\_in , *s*.qty\_out , *s*.qty\_in\_order , *s*.qty\_out\_order  from inv.stock *s* where *s*.tr\_no \='SO2602090002'  L \= 100  M \= 0  S \= 0  ![][image18]    handle by BE saat cancel order  curl 'https://best.scyllax.online/sales/v1/orders/status' \\ Be cek pada no SO yang di cancel pada case diatas , qty\_out (satuan terkecil) \= 2.400  Maka saat cancel order Insert ke inv.stock  cust\_id stock\_id stock\_date tr\_code tr\_no wh\_id pro\_id item\_cdn qty\_in qty\_out unit\_price cogs ref\_det\_id created\_at qty\_in\_order qty\_out\_order C220010001 19572 2026-02-09 SO SO2602090002 63 474 1 2400 0 6432.0000 0.0000 4949 1770647314 0 0 C220010001 19572 2026-02-09 CO SO2602090002-CO 63 474 1 0 0 6432.0000 0.0000 4949 1770647314 0 2400  update ke inv.warehouse\_stock berdasarkan cust\_id, wh\_id, dan pro\_id qty \= bertambah sebesar qty\_out sales order (untuk case ini 2400\) qty\_on\_order \= berkurang sebesar qty\_out sales order (untuk case ini 2400\)    |
| :---- |

## **PROMO** 

1. ### **Detail Sales Order** 

Dokumen Promo : [https://docs.google.com/spreadsheets/d/1UGsfcV0-Lwhi9rv6cNOdkXZwCOF3mo5Ei-BNhbowPYE/edit?gid=0\#gid=0](https://docs.google.com/spreadsheets/d/1UGsfcV0-Lwhi9rv6cNOdkXZwCOF3mo5Ei-BNhbowPYE/edit?gid=0#gid=0) 

| Response Detail Order | Source Data |
| :---- | :---- |
| **details.normal\[\]** |  |
| promo\_so1 | sls.order\_detail.promo\_so1 |
| promo\_so2 | sls.order\_detail.promo\_so2 |
| promo\_so3 | sls.order\_detail.promo\_so3 |
| promo\_so4 | sls.order\_detail.promo\_so4 |
| promo\_so5 | sls.order\_detail.promo\_so5 |
| promo\_remarks\_so | sls.order\_detail.promo\_remarks\_so |
| is\_product\_promotion\_so | sls.order\_detail.is\_product\_promotion\_so |
| **details.promo\_remarks\_so\[\]** | sls.order.promo\_remarks\_so |
| **details\_final.normal\[\]** |  |
| promo\_final1 | sls.order\_detail.promo\_final1 |
| promo\_final2 | sls.order\_detail.promo\_final2 |
| promo\_final3 | sls.order\_detail.promo\_final3 |
| promo\_final4 | sls.order\_detail.promo\_final4 |
| promo\_final5 | sls.order\_detail.promo\_final5 |
| promo\_remarks\_final | sls.order\_detail.promo\_remarks\_final |
| is\_product\_promotion\_final | sls.order\_detail.is\_product\_promotion\_final |
| **details.promo\_remarks\_final\[\]** | sls.order.promo\_remarks\_final |
| **purchase\_details\_\[\]** |  |
| promo\_po1 | sls.order\_detail.promo\_po1 |
| promo\_po2 | sls.order\_detail.promo\_po2 |
| promo\_po3 | sls.order\_detail.promo\_po3 |
| promo\_po4 | sls.order\_detail.promo\_po4 |
| promo\_po5 | sls.order\_detail.promo\_po5 |
| promo\_remarks\_po | sls.order\_detail.promo\_remarks\_po |
| is\_product\_promotion\_po | sls.order\_detail.is\_product\_promotion\_po |
| **details.promo\_remarks\_po\[\]** | sls.order.promo\_remarks\_po |

#### **Contoh Response :** 

| contoh : promo reward value / percentage   {     "message": "",     "data": {       "ro\_no": "SO2603070001",       "source": "web",       "is\_proforma\_inv": false,       "order\_no": null,       "ro\_date": "2026-03-07",       "val\_date": null,       "salesman\_id": 347,       "salesman\_code": "AD003",       "sales\_name": "Ady Nugroho",       "wh\_id": 301,       "wh\_code": "842347",       "wh\_name": "Canvas \- Ady Nugroho",       "outlet\_id": 415,       "outlet\_code": "LV000008",       "outlet\_name": "Lavina TK 8",       "outlet\_address1": "poris",       "outlet\_address2": "",       "inv\_addr1": "Cibitung Blok Z No.8",       "inv\_addr2": "",       "delivery\_date": "2026-03-07",       "po\_no": null,       "vehicle\_no": null,       "pay\_type": 1,       "pay\_type\_name": "Cash On Delivery",       "reff\_no": null,       "mobile\_id": 1,       "sub\_total": 400000,       "sub\_total\_final": 400000,       "disc": 0,       "disc\_value": 20000,       "disc\_value\_final": 20000,       "promo\_value": 0,       "promo\_value\_final": 0,       "promo\_bg\_value": 0,       "promo\_bg\_value\_final": 0,       "cash\_disc\_value": 0,       "tot\_disc1": 0,       "tot\_disc2": 0,       "vat": 0,       "vat\_value": 0,       "vat\_value\_final": 0,       "total": 380000,       "total\_final": 380000,       "data\_status": 2,       "data\_status\_name": "Processed",       "data\_source": 1,       "updated\_at": "2026-03-07T10:37:24.590384Z",       "updated\_by\_name": "Dist IDE Sda",       "due\_date": null,       "details": {         "normal": \[           {             "order\_detail\_id": 5642,             "seq\_no": 0,             "pro\_id": 714,             "pro\_code": "AF-0001",             "pro\_name": "Action Figure Superman",             "order\_status": "",             "item\_type": 1,             "promo1": 0,             "promo2": 0,             "promo3": 0,             "promo4": 0,             "promo5": 0,             "promo\_total": 0,             "remarks": \[\],             "qty": 2,             "qty\_final": 2,             "qty\_po": 10,             "qty\_po1": null,             "qty\_po2": null,             "qty\_po3": null,             "qty1": 2,             "qty2": 0,             "qty3": 0,             "qty4": null,             "qty5": null,             "qty1\_final": 2,             "qty2\_final": 0,             "qty3\_final": 0,             "qty4\_final": null,             "qty5\_final": null,             "qty1\_stok": 2,             "qty2\_stok": 0,             "qty3\_stok": 4,             "purch\_price1": 150000,             "purch\_price2": 1500000,             "purch\_price3": 1500000,             "purch\_price4": 0,             "purch\_price5": 0,             "sell\_price1": 200000,             "sell\_price2": 500000,             "sell\_price3": 1000000,             "sell\_price4": 0,             "sell\_price5": 0,             "sell\_price\_po1": null,             "sell\_price\_po2": null,             "sell\_price\_po3": null,             "sell\_price\_final1": 200000,             "sell\_price\_final2": 500000,             "sell\_price\_final3": 1000000,             "sell\_price\_system1": 200000,             "sell\_price\_system2": 200000,             "sell\_price\_system3": 200000,             "sell\_price\_system4": null,             "sell\_price\_system5": null,             "amount": 380000,             "amount\_final": 380000,             "promo\_value": 0,             "promo\_value\_final": 0,             "disc\_value": 20000,             "disc\_value\_final": 20000,             "batch\_no": null,             "exp\_date": null,             "vat": 0,             "vat\_bg": 0,             "vat\_lg\_sell": 0,             "vat\_value": 0,             "vat\_value\_final": 0,             "vat\_bg\_value": 0,             "vat\_lg\_value": null,             "vat\_lg\_sell\_value": 0,             "disc\_po": 100000,             "vat\_value\_po": 0,             "unit\_id1": "PCS",             "unit\_id2": "CRT",             "unit\_id3": "CRT",             "unit\_id4": "",             "unit\_id5": "",             "conv\_unit2": 10,             "conv\_unit3": 1,             "conv\_unit4": 0,             "conv\_unit5": 0,             "notes": "",             "discount\_id": "DISC0001",             "promo\_so1": 20000,             "promo\_so2": 0,             "promo\_so3": 0,             "promo\_so4": 0,             "promo\_so5": 0,             "promo\_remarks\_so": \["slab01", "slab02"\],             "is\_product\_promotion\_so": false           }         \],         "promo": null,         "promo\_remarks\_so": \[         {             "promo\_id": "slab01",             "promo\_desc": "Syarat Value & Reward Percentage (Value)",             "slab\_id": "695ba0009f8c1325e300ae2f",             "slab\_desc": "slab 3 Value",             "slab\_reward": 3,             "slab\_reward\_uom": "",             "slab\_reward\_type": "percentage",             "slab\_per\_scope": "",             "total\_gross\_value": 5600000,             "products\_eligible": \[                 709,                 710             \],             "reward\_value": null,             "reward\_percentage": \[                 {                     "pro\_id": 709,                     "gross\_value": 2000000,                     "promo1": 60000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1940000                 },                 {                     "pro\_id": 710,                     "gross\_value": 3600000,                     "promo1": 108000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 3492000                 }             \],             "reward\_product": null         }     \],          "reward\_products": \[\]       },       "details\_final": {         "normal": \[           {             "order\_detail\_id": 5642,             "seq\_no": 0,             "pro\_id": 714,             "pro\_code": "AF-0001",             "pro\_name": "Action Figure Superman",             "order\_status": "",             "item\_type": 1,             "promo1": 0,             "promo2": 0,             "promo3": 0,             "promo4": 0,             "promo5": 0,             "promo\_total": 0,             "remarks": \[\],             "qty": 2,             "qty\_final": 2,             "qty\_po": 10,             "qty\_po1": null,             "qty\_po2": null,             "qty\_po3": null,             "qty1": 2,             "qty2": 0,             "qty3": 0,             "qty4": null,             "qty5": null,             "qty1\_final": 2,             "qty2\_final": 0,             "qty3\_final": 0,             "qty4\_final": null,             "qty5\_final": null,             "qty1\_stok": 2,             "qty2\_stok": 0,             "qty3\_stok": 4,             "purch\_price1": 150000,             "purch\_price2": 1500000,             "purch\_price3": 1500000,             "purch\_price4": 0,             "purch\_price5": 0,             "sell\_price1": 200000,             "sell\_price2": 500000,             "sell\_price3": 1000000,             "sell\_price4": 0,             "sell\_price5": 0,             "sell\_price\_po1": null,             "sell\_price\_po2": null,             "sell\_price\_po3": null,             "sell\_price\_final1": 200000,             "sell\_price\_final2": 500000,             "sell\_price\_final3": 1000000,             "sell\_price\_system1": 200000,             "sell\_price\_system2": 200000,             "sell\_price\_system3": 200000,             "sell\_price\_system4": null,             "sell\_price\_system5": null,             "amount": 380000,             "amount\_final": 380000,             "promo\_value": 0,             "promo\_value\_final": 0,             "disc\_value": 20000,             "disc\_value\_final": 20000,             "batch\_no": null,             "exp\_date": null,             "vat": 0,             "vat\_bg": 0,             "vat\_lg\_sell": 0,             "vat\_value": 0,             "vat\_value\_final": 0,             "vat\_bg\_value": 0,             "vat\_lg\_value": null,             "vat\_lg\_sell\_value": 0,             "disc\_po": 100000,             "vat\_value\_po": 0,             "unit\_id1": "PCS",             "unit\_id2": "CRT",             "unit\_id3": "CRT",             "unit\_id4": "",             "unit\_id5": "",             "conv\_unit2": 10,             "conv\_unit3": 1,             "conv\_unit4": 0,             "conv\_unit5": 0,             "notes": "",             "discount\_id": "DISC0001",             "promo\_final1": 20000,             "promo\_final2": 0,             "promo\_final3": 0,             "promo\_final4": 0,             "promo\_final5": 0,             "promo\_remarks\_final": \["slab01", "slab02"\],             "is\_product\_promotion\_final": false           }         \],         "promo": null,         "promo\_remarks\_final":\[         {             "promo\_id": "slab01",             "promo\_desc": "Syarat Value & Reward Percentage (Value)",             "slab\_id": "695ba0009f8c1325e300ae2f",             "slab\_desc": "slab 3 Value",             "slab\_reward": 3,             "slab\_reward\_uom": "",             "slab\_reward\_type": "percentage",             "slab\_per\_scope": "",             "total\_gross\_value": 5600000,             "products\_eligible": \[                 709,                 710             \],             "reward\_value": null,             "reward\_percentage": \[                 {                     "pro\_id": 709,                     "gross\_value": 2000000,                     "promo1": 60000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1940000                 },                 {                     "pro\_id": 710,                     "gross\_value": 3600000,                     "promo1": 108000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 3492000                 }             \],             "reward\_product": null         }     \], ,         "reward\_products": \[\]       },       "purchase\_details": {         "normal": \[           {             "order\_detail\_id": 5642,             "seq\_no": 0,             "pro\_id": 714,             "pro\_code": "AF-0001",             "pro\_name": "Action Figure Superman",             "order\_status": "2",             "item\_type": 1,             "promo\_total": 0,             "remarks": \[\],             "qty": 2,             "qty\_final": 2,             "qty\_po": 10,             "qty1": 2,             "qty2": 0,             "qty3": 0,             "qty4": null,             "qty5": null,             "qty1\_final": 2,             "qty2\_final": 0,             "qty3\_final": 0,             "qty4\_final": null,             "qty5\_final": null,             "qty1\_stok": 2,             "qty2\_stok": 0,             "qty3\_stok": 4,             "purch\_price1": 150000,             "purch\_price2": 1500000,             "purch\_price3": 1500000,             "purch\_price4": 0,             "purch\_price5": 0,             "sell\_price1": 200000,             "sell\_price2": 500000,             "sell\_price3": 1000000,             "sell\_price4": 0,             "sell\_price5": 0,             "sell\_price\_po1": 200000,             "sell\_price\_po2": 500000,             "sell\_price\_po3": 1000000,             "sell\_price\_final1": 200000,             "sell\_price\_final2": 500000,             "sell\_price\_final3": 1000000,             "sell\_price\_system1": 200000,             "sell\_price\_system2": 200000,             "sell\_price\_system3": 200000,             "sell\_price\_system4": null,             "sell\_price\_system5": null,             "amount": 380000,             "amount\_final": 380000,             "promo\_value": 0,             "promo\_value\_final": 0,             "disc\_value": 20000,             "disc\_value\_final": 20000,             "batch\_no": null,             "exp\_date": null,             "vat": 0,             "vat\_bg": 0,             "vat\_lg\_sell": 0,             "vat\_value": 0,             "vat\_value\_final": 0,             "vat\_bg\_value": 0,             "vat\_lg\_value": null,             "vat\_lg\_sell\_value": 0,             "disc\_po": 100000,             "vat\_value\_po": 0,             "unit\_id1": "PCS",             "unit\_id2": "CRT",             "unit\_id3": "CRT",             "unit\_id4": "",             "unit\_id5": "",             "conv\_unit2": 10,             "conv\_unit3": 1,             "conv\_unit4": 0,             "conv\_unit5": 0,             "notes": "",             "discount\_id": "DISC0001",             "promo\_po1": 20000,             "promo\_po2": 0,             "promo\_po3": 0,             "promo\_po4": 0,             "promo\_po5": 0,             "promo\_remarks\_po": \["slab01", "slab02"\],             "is\_product\_promotion\_po": false           }         \],         "promo": \[\],         "promo\_remarks\_po": \["slab01", "slab02"\],         "reward\_products": \[\]       },       "tr\_code": null,       "is\_closed": false,       "notes": "",       "invoice\_no": null,       "invoice\_date": null,       "is\_printed": false,       "printed\_by": null,       "printed\_by\_name": null,       "printed\_at": null,       "validate\_stok": true,       "validate\_stok\_message": "Sufficient Stock",       "validate\_credit\_limit": false,       "validate\_credit\_limit\_message": "Over Limit (370.687)",       "validate\_overdue": true,       "validate\_overdue\_message": "Allowed (Unlimited)",       "validate\_outstanding": true,       "validate\_outstanding\_message": "Allowed (Unlimited)",       "validate\_summary": false,       "credit\_limit\_type": 2,       "credit\_limit\_action": 1,       "credit\_limit\_action\_name": "Warning",       "sales\_inv\_limit\_type": null,       "sales\_inv\_limit\_action": 1,       "sales\_inv\_limit\_action\_name": "Warning",       "obs\_type": null,       "obs\_limit\_action": 1,       "obs\_limit\_action\_name": "Warning",       "order\_approval\_request\_id": null,       "order\_approval\_request\_emp\_approval\_status": null     },     "request\_id": "5d7d7afd-4638-4627-9e79-c42fe75d9978" }  |
| :---- |
| {     "message": "",     "data": {       "ro\_no": "SO2603070001",       "source": "web",       "is\_proforma\_inv": false,       "order\_no": null,       "ro\_date": "2026-03-07",       "val\_date": null,       "salesman\_id": 347,       "salesman\_code": "AD003",       "sales\_name": "Ady Nugroho",       "wh\_id": 301,       "wh\_code": "842347",       "wh\_name": "Canvas \- Ady Nugroho",       "outlet\_id": 415,       "outlet\_code": "LV000008",       "outlet\_name": "Lavina TK 8",       "outlet\_address1": "poris",       "outlet\_address2": "",       "inv\_addr1": "Cibitung Blok Z No.8",       "inv\_addr2": "",       "delivery\_date": "2026-03-07",       "po\_no": null,       "vehicle\_no": null,       "pay\_type": 1,       "pay\_type\_name": "Cash On Delivery",       "reff\_no": null,       "mobile\_id": 1,       "sub\_total": 400000,       "sub\_total\_final": 400000,       "disc": 0,       "disc\_value": 20000,       "disc\_value\_final": 20000,       "promo\_value": 0,       "promo\_value\_final": 0,       "promo\_bg\_value": 0,       "promo\_bg\_value\_final": 0,       "cash\_disc\_value": 0,       "tot\_disc1": 0,       "tot\_disc2": 0,       "vat": 0,       "vat\_value": 0,       "vat\_value\_final": 0,       "total": 380000,       "total\_final": 380000,       "data\_status": 2,       "data\_status\_name": "Processed",       "data\_source": 1,       "updated\_at": "2026-03-07T10:37:24.590384Z",       "updated\_by\_name": "Dist IDE Sda",       "due\_date": null,       "details": {         "normal": \[           {             "order\_detail\_id": 5642,             "seq\_no": 0,             "pro\_id": 714,             "pro\_code": "AF-0001",             "pro\_name": "Action Figure Superman",             "order\_status": "",             "item\_type": 1,             "promo1": 0,             "promo2": 0,             "promo3": 0,             "promo4": 0,             "promo5": 0,             "promo\_total": 0,             "remarks": \[\],             "qty": 2,             "qty\_final": 2,             "qty\_po": 10,             "qty\_po1": null,             "qty\_po2": null,             "qty\_po3": null,             "qty1": 2,             "qty2": 0,             "qty3": 0,             "qty4": null,             "qty5": null,             "qty1\_final": 2,             "qty2\_final": 0,             "qty3\_final": 0,             "qty4\_final": null,             "qty5\_final": null,             "qty1\_stok": 2,             "qty2\_stok": 0,             "qty3\_stok": 4,             "purch\_price1": 150000,             "purch\_price2": 1500000,             "purch\_price3": 1500000,             "purch\_price4": 0,             "purch\_price5": 0,             "sell\_price1": 200000,             "sell\_price2": 500000,             "sell\_price3": 1000000,             "sell\_price4": 0,             "sell\_price5": 0,             "sell\_price\_po1": null,             "sell\_price\_po2": null,             "sell\_price\_po3": null,             "sell\_price\_final1": 200000,             "sell\_price\_final2": 500000,             "sell\_price\_final3": 1000000,             "sell\_price\_system1": 200000,             "sell\_price\_system2": 200000,             "sell\_price\_system3": 200000,             "sell\_price\_system4": null,             "sell\_price\_system5": null,             "amount": 380000,             "amount\_final": 380000,             "promo\_value": 0,             "promo\_value\_final": 0,             "disc\_value": 20000,             "disc\_value\_final": 20000,             "batch\_no": null,             "exp\_date": null,             "vat": 0,             "vat\_bg": 0,             "vat\_lg\_sell": 0,             "vat\_value": 0,             "vat\_value\_final": 0,             "vat\_bg\_value": 0,             "vat\_lg\_value": null,             "vat\_lg\_sell\_value": 0,             "disc\_po": 100000,             "vat\_value\_po": 0,             "unit\_id1": "PCS",             "unit\_id2": "CRT",             "unit\_id3": "CRT",             "unit\_id4": "",             "unit\_id5": "",             "conv\_unit2": 10,             "conv\_unit3": 1,             "conv\_unit4": 0,             "conv\_unit5": 0,             "notes": "",             "discount\_id": "DISC0001",             "promo\_so1": 20000,             "promo\_so2": 0,             "promo\_so3": 0,             "promo\_so4": 0,             "promo\_so5": 0,             "promo\_remarks\_so": \["slab01", "slab02"\],             "is\_product\_promotion\_so": false           },           {             "order\_detail\_id": 5645,             "seq\_no": 0,             "pro\_id": 716,             "pro\_code": "AF-0001",             "pro\_name": "Action Figure Spongebob",             "order\_status": "",             "item\_type": 1,             "promo1": 0,             "promo2": 0,             "promo3": 0,             "promo4": 0,             "promo5": 0,             "promo\_total": 0,             "remarks": \[\],             "qty": 2,             "qty\_final": 2,             "qty\_po": 10,             "qty\_po1": null,             "qty\_po2": null,             "qty\_po3": null,             "qty1": 2,             "qty2": 0,             "qty3": 0,             "qty4": null,             "qty5": null,             "qty1\_final": 2,             "qty2\_final": 0,             "qty3\_final": 0,             "qty4\_final": null,             "qty5\_final": null,             "qty1\_stok": 2,             "qty2\_stok": 0,             "qty3\_stok": 4,             "purch\_price1": 150000,             "purch\_price2": 1500000,             "purch\_price3": 1500000,             "purch\_price4": 0,             "purch\_price5": 0,             "sell\_price1": 200000,             "sell\_price2": 500000,             "sell\_price3": 1000000,             "sell\_price4": 0,             "sell\_price5": 0,             "sell\_price\_po1": null,             "sell\_price\_po2": null,             "sell\_price\_po3": null,             "sell\_price\_final1": 200000,             "sell\_price\_final2": 500000,             "sell\_price\_final3": 1000000,             "sell\_price\_system1": 200000,             "sell\_price\_system2": 200000,             "sell\_price\_system3": 200000,             "sell\_price\_system4": null,             "sell\_price\_system5": null,             "amount": 380000,             "amount\_final": 380000,             "promo\_value": 0,             "promo\_value\_final": 0,             "disc\_value": 20000,             "disc\_value\_final": 20000,             "batch\_no": null,             "exp\_date": null,             "vat": 0,             "vat\_bg": 0,             "vat\_lg\_sell": 0,             "vat\_value": 0,             "vat\_value\_final": 0,             "vat\_bg\_value": 0,             "vat\_lg\_value": null,             "vat\_lg\_sell\_value": 0,             "disc\_po": 100000,             "vat\_value\_po": 0,             "unit\_id1": "PCS",             "unit\_id2": "CRT",             "unit\_id3": "CRT",             "unit\_id4": "",             "unit\_id5": "",             "conv\_unit2": 10,             "conv\_unit3": 1,             "conv\_unit4": 0,             "conv\_unit5": 0,             "notes": "",             "discount\_id": "DISC0001",             "promo\_so1": 20000,             "promo\_so2": 0,             "promo\_so3": 0,             "promo\_so4": 0,             "promo\_so5": 0,             "promo\_remarks\_so": \["slab01", "slab02"\],             "is\_product\_promotion\_so": true           }         \],         "promo": null,         "promo\_remarks\_so": \[         {             "promo\_id": "slab01",             "promo\_desc": "Syarat Value & Reward Percentage (Value)",             "slab\_id": "695ba0009f8c1325e300ae2f",             "slab\_desc": "slab 3 Value",             "slab\_reward": 3,             "slab\_reward\_uom": "",             "slab\_reward\_type": "percentage",             "slab\_per\_scope": "",             "total\_gross\_value": 5600000,             "products\_eligible": \[                 709,                 710             \],             "reward\_value": null,             "reward\_percentage": \[                 {                     "pro\_id": 709,                     "gross\_value": 2000000,                     "promo1": 60000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1940000                 },                 {                     "pro\_id": 710,                     "gross\_value": 3600000,                     "promo1": 108000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 3492000                 }             \],             "reward\_product": null         }     \],         "reward\_products": \[\]       },       "details\_final": {         "normal": \[           {             "order\_detail\_id": 5642,             "seq\_no": 0,             "pro\_id": 714,             "pro\_code": "AF-0001",             "pro\_name": "Action Figure Superman",             "order\_status": "",             "item\_type": 1,             "promo1": 0,             "promo2": 0,             "promo3": 0,             "promo4": 0,             "promo5": 0,             "promo\_total": 0,             "remarks": \[\],             "qty": 2,             "qty\_final": 2,             "qty\_po": 10,             "qty\_po1": null,             "qty\_po2": null,             "qty\_po3": null,             "qty1": 2,             "qty2": 0,             "qty3": 0,             "qty4": null,             "qty5": null,             "qty1\_final": 2,             "qty2\_final": 0,             "qty3\_final": 0,             "qty4\_final": null,             "qty5\_final": null,             "qty1\_stok": 2,             "qty2\_stok": 0,             "qty3\_stok": 4,             "purch\_price1": 150000,             "purch\_price2": 1500000,             "purch\_price3": 1500000,             "purch\_price4": 0,             "purch\_price5": 0,             "sell\_price1": 200000,             "sell\_price2": 500000,             "sell\_price3": 1000000,             "sell\_price4": 0,             "sell\_price5": 0,             "sell\_price\_po1": null,             "sell\_price\_po2": null,             "sell\_price\_po3": null,             "sell\_price\_final1": 200000,             "sell\_price\_final2": 500000,             "sell\_price\_final3": 1000000,             "sell\_price\_system1": 200000,             "sell\_price\_system2": 200000,             "sell\_price\_system3": 200000,             "sell\_price\_system4": null,             "sell\_price\_system5": null,             "amount": 380000,             "amount\_final": 380000,             "promo\_value": 0,             "promo\_value\_final": 0,             "disc\_value": 20000,             "disc\_value\_final": 20000,             "batch\_no": null,             "exp\_date": null,             "vat": 0,             "vat\_bg": 0,             "vat\_lg\_sell": 0,             "vat\_value": 0,             "vat\_value\_final": 0,             "vat\_bg\_value": 0,             "vat\_lg\_value": null,             "vat\_lg\_sell\_value": 0,             "disc\_po": 100000,             "vat\_value\_po": 0,             "unit\_id1": "PCS",             "unit\_id2": "CRT",             "unit\_id3": "CRT",             "unit\_id4": "",             "unit\_id5": "",             "conv\_unit2": 10,             "conv\_unit3": 1,             "conv\_unit4": 0,             "conv\_unit5": 0,             "notes": "",             "discount\_id": "DISC0001",             "promo\_final1": 20000,             "promo\_final2": 0,             "promo\_final3": 0,             "promo\_final4": 0,             "promo\_final5": 0,             "promo\_remarks\_final": \["slab01", "slab02"\],             "is\_product\_promotion\_final": false           },            {             "order\_detail\_id": 5645,             "seq\_no": 0,             "pro\_id": 714,             "pro\_code": "AF-0001",             "pro\_name": "Action Figure Spongebob",             "order\_status": "",             "item\_type": 1,             "promo1": 0,             "promo2": 0,             "promo3": 0,             "promo4": 0,             "promo5": 0,             "promo\_total": 0,             "remarks": \[\],             "qty": 2,             "qty\_final": 2,             "qty\_po": 10,             "qty\_po1": null,             "qty\_po2": null,             "qty\_po3": null,             "qty1": 2,             "qty2": 0,             "qty3": 0,             "qty4": null,             "qty5": null,             "qty1\_final": 2,             "qty2\_final": 0,             "qty3\_final": 0,             "qty4\_final": null,             "qty5\_final": null,             "qty1\_stok": 2,             "qty2\_stok": 0,             "qty3\_stok": 4,             "purch\_price1": 150000,             "purch\_price2": 1500000,             "purch\_price3": 1500000,             "purch\_price4": 0,             "purch\_price5": 0,             "sell\_price1": 200000,             "sell\_price2": 500000,             "sell\_price3": 1000000,             "sell\_price4": 0,             "sell\_price5": 0,             "sell\_price\_po1": null,             "sell\_price\_po2": null,             "sell\_price\_po3": null,             "sell\_price\_final1": 200000,             "sell\_price\_final2": 500000,             "sell\_price\_final3": 1000000,             "sell\_price\_system1": 200000,             "sell\_price\_system2": 200000,             "sell\_price\_system3": 200000,             "sell\_price\_system4": null,             "sell\_price\_system5": null,             "amount": 380000,             "amount\_final": 380000,             "promo\_value": 0,             "promo\_value\_final": 0,             "disc\_value": 20000,             "disc\_value\_final": 20000,             "batch\_no": null,             "exp\_date": null,             "vat": 0,             "vat\_bg": 0,             "vat\_lg\_sell": 0,             "vat\_value": 0,             "vat\_value\_final": 0,             "vat\_bg\_value": 0,             "vat\_lg\_value": null,             "vat\_lg\_sell\_value": 0,             "disc\_po": 100000,             "vat\_value\_po": 0,             "unit\_id1": "PCS",             "unit\_id2": "CRT",             "unit\_id3": "CRT",             "unit\_id4": "",             "unit\_id5": "",             "conv\_unit2": 10,             "conv\_unit3": 1,             "conv\_unit4": 0,             "conv\_unit5": 0,             "notes": "",             "discount\_id": "DISC0001",             "promo\_final1": 20000,             "promo\_final2": 0,             "promo\_final3": 0,             "promo\_final4": 0,             "promo\_final5": 0,             "promo\_remarks\_final": \["slab01", "slab02"\],             "is\_product\_promotion\_final": true           }         \],         "promo": null,         "promo\_remarks\_final":\[         {             "promo\_id": "slab01",             "promo\_desc": "Syarat Value & Reward Percentage (Value)",             "slab\_id": "695ba0009f8c1325e300ae2f",             "slab\_desc": "slab 3 Value",             "slab\_reward": 3,             "slab\_reward\_uom": "",             "slab\_reward\_type": "percentage",             "slab\_per\_scope": "",             "total\_gross\_value": 5600000,             "products\_eligible": \[                 709,                 710             \],             "reward\_value": null,             "reward\_percentage": \[                 {                     "pro\_id": 709,                     "gross\_value": 2000000,                     "promo1": 60000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1940000                 },                 {                     "pro\_id": 710,                     "gross\_value": 3600000,                     "promo1": 108000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 3492000                 }             \],             "reward\_product": null         }     \], ,         "reward\_products": \[\]       },       "purchase\_details": {         "normal": \[           {             "order\_detail\_id": 5642,             "seq\_no": 0,             "pro\_id": 714,             "pro\_code": "AF-0001",             "pro\_name": "Action Figure Superman",             "order\_status": "2",             "item\_type": 1,             "promo\_total": 0,             "remarks": \[\],             "qty": 2,             "qty\_final": 2,             "qty\_po": 10,             "qty1": 2,             "qty2": 0,             "qty3": 0,             "qty4": null,             "qty5": null,             "qty1\_final": 2,             "qty2\_final": 0,             "qty3\_final": 0,             "qty4\_final": null,             "qty5\_final": null,             "qty1\_stok": 2,             "qty2\_stok": 0,             "qty3\_stok": 4,             "purch\_price1": 150000,             "purch\_price2": 1500000,             "purch\_price3": 1500000,             "purch\_price4": 0,             "purch\_price5": 0,             "sell\_price1": 200000,             "sell\_price2": 500000,             "sell\_price3": 1000000,             "sell\_price4": 0,             "sell\_price5": 0,             "sell\_price\_po1": 200000,             "sell\_price\_po2": 500000,             "sell\_price\_po3": 1000000,             "sell\_price\_final1": 200000,             "sell\_price\_final2": 500000,             "sell\_price\_final3": 1000000,             "sell\_price\_system1": 200000,             "sell\_price\_system2": 200000,             "sell\_price\_system3": 200000,             "sell\_price\_system4": null,             "sell\_price\_system5": null,             "amount": 380000,             "amount\_final": 380000,             "promo\_value": 0,             "promo\_value\_final": 0,             "disc\_value": 20000,             "disc\_value\_final": 20000,             "batch\_no": null,             "exp\_date": null,             "vat": 0,             "vat\_bg": 0,             "vat\_lg\_sell": 0,             "vat\_value": 0,             "vat\_value\_final": 0,             "vat\_bg\_value": 0,             "vat\_lg\_value": null,             "vat\_lg\_sell\_value": 0,             "disc\_po": 100000,             "vat\_value\_po": 0,             "unit\_id1": "PCS",             "unit\_id2": "CRT",             "unit\_id3": "CRT",             "unit\_id4": "",             "unit\_id5": "",             "conv\_unit2": 10,             "conv\_unit3": 1,             "conv\_unit4": 0,             "conv\_unit5": 0,             "notes": "",             "discount\_id": "DISC0001",             "promo\_po1": 20000,             "promo\_po2": 0,             "promo\_po3": 0,             "promo\_po4": 0,             "promo\_po5": 0,             "promo\_remarks\_po": \["slab01", "slab02"\],             "is\_product\_promotion\_po": false           }         \],         "promo": \[\],         "promo\_remarks\_po": \["slab01", "slab02"\],         "reward\_products": \[\]       },       "tr\_code": null,       "is\_closed": false,       "notes": "",       "invoice\_no": null,       "invoice\_date": null,       "is\_printed": false,       "printed\_by": null,       "printed\_by\_name": null,       "printed\_at": null,       "validate\_stok": true,       "validate\_stok\_message": "Sufficient Stock",       "validate\_credit\_limit": false,       "validate\_credit\_limit\_message": "Over Limit (370.687)",       "validate\_overdue": true,       "validate\_overdue\_message": "Allowed (Unlimited)",       "validate\_outstanding": true,       "validate\_outstanding\_message": "Allowed (Unlimited)",       "validate\_summary": false,       "credit\_limit\_type": 2,       "credit\_limit\_action": 1,       "credit\_limit\_action\_name": "Warning",       "sales\_inv\_limit\_type": null,       "sales\_inv\_limit\_action": 1,       "sales\_inv\_limit\_action\_name": "Warning",       "obs\_type": null,       "obs\_limit\_action": 1,       "obs\_limit\_action\_name": "Warning",       "order\_approval\_request\_id": null,       "order\_approval\_request\_emp\_approval\_status": null     },     "request\_id": "5d7d7afd-4638-4627-9e79-c42fe75d9978" }  |

2. ### **Create Sales Order** 

| ![][image19] |
| :---- |
| **sequenceDiagram**    actor User    participant FE as Frontend    participant SO as Sales Order API    participant Promo as Promotion API    participant DB as Database     User**\-\>\>**FE**:** Input Sales Order (items, qty, price)     FE**\-\>\>**SO**:** POST /sales/v1/orders    Note right of FE**:** items, qty, price, outlet\_id, salesman\_id     SO**\-\>\>**Promo**:** GET /sales/v2/promotions/consult    Note right of SO**:** Send item list, qty, outlet, date     Promo**\--\>\>**SO**:** Return eligible promotions    Note left of Promo**:** discount / bonus / reward     SO**\-\>\>**SO**:** Calculate promotion result     SO**\-\>\>**DB**:** Insert sls.order    Note right of DB**:** header order\\n(order\_no, outlet\_id, total\_amount, promo\_summary)     SO**\-\>\>**DB**:** Insert sls.order\_detail (loop items)       DB**\--\>\>**SO**:** Success     SO**\--\>\>**FE**:** Order Created Response    FE**\--\>\>**User**:** Show Order Success  |

* #### FE hit API  : [https://best.scyllax.online/sales/v1/orders](https://best.scyllax.online/sales/v1/orders) 

  CURL : 


| curl 'https://best.scyllax.online/sales/v1/orders' \\   \-H 'sec-ch-ua-platform: "macOS"' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsImVtcF9pZCI6MCwiZXhwaXJlcyI6MTc3MjYxNTU3NiwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4NTc0ODEyMzEyOCIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IkRpc3QgSURFIFNkYSIsInVzZXJfaWQiOjEyLCJ1c2VyX25hbWUiOiJkaXN0QHNkYS5pZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODU3NDgxMjMxMjgifQ.\_bkN46a5FsGJJWHtEV4HjAYf1HkPBKy9RGPjmmlBKb0' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'sec-ch-ua: "Not:A-Brand";v="99", "Google Chrome";v="145", "Chromium";v="145"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Content-Type: application/json' \\   \--data-raw '{"mobile\_id":1,"data\_status":2,"data\_source":1,"disc":0,"disc\_value":0,"promo\_value":0,"cash\_disc\_value":0,"tot\_disc1":0,"tot\_disc2":0,"vat":11,"vat\_value":16500,"total":166500,"salesman\_id":357,"outlet\_id":199,"ro\_date":"2026-03-04","due\_date":null,"delivery\_date":"2026-03-04","pay\_type":1,"sub\_total":150000,"wh\_id":63,"notes":"","details":{"normal":\[{"vat\_bg":0,"vat\_lg\_sell":0,"vat\_bg\_value":0,"vat\_lg\_sell\_value":0,"pro\_id":5896,"pro\_code":"YP5218","pro\_name":"YUPI ROLETTO BOX 12X24X6GR","unit\_id1":"KLG","unit\_id2":"KRT","unit\_id3":"KRT","unit\_id4":"","unit\_id5":"","conv\_unit2":50,"conv\_unit3":1,"conv\_unit4":0,"conv\_unit5":0,"purch\_price1":70000,"purch\_price2":50000,"purch\_price3":30000,"purch\_price4":0,"purch\_price5":0,"sell\_price1":23000,"sell\_price2":800000,"sell\_price3":75000,"sell\_price4":0,"sell\_price5":0,"vat":11,"pl\_id":61,"pl\_code":"001","pl\_name":"Fumakila","brand\_id":104,"brand\_code":"2","sbrand1\_id":124,"sbrand1\_code":"002","sbrand1\_name":"Vape Non Coil","vat\_value":16500,"amount":166500,"notes":"","item\_type":1,"disc\_value":0,"qty1\_stok":0,"qty2\_stok":0,"qty3\_stok":70,"qty1":0,"qty2":0,"qty3":2}\],"promo":null}}'  |
| :---- |


  ##### Payload 


| {   "mobile\_id": 1,   "data\_status": 2,   "data\_source": 1,   "disc": 0,   "disc\_value": 0,   "promo\_value": 0,   "cash\_disc\_value": 0,   "tot\_disc1": 0,   "tot\_disc2": 0,   "vat": 11,   "vat\_value": 16500,   "total": 166500,   "salesman\_id": 357,   "outlet\_id": 225,   "ro\_date": "2026-03-04",   "due\_date": null,   "delivery\_date": "2026-03-04",   "pay\_type": 1,   "sub\_total": 150000,   "wh\_id": 63,   "notes": "",   "details": {     "normal": \[       {         "vat\_bg": 0,         "vat\_lg\_sell": 0,         "vat\_bg\_value": 0,         "vat\_lg\_sell\_value": 0,         "pro\_id": 5896,         "pro\_code": "YP5218",         "pro\_name": "YUPI ROLETTO BOX 12X24X6GR",         "unit\_id1": "KLG",         "unit\_id2": "KRT",         "unit\_id3": "KRT",         "unit\_id4": "",         "unit\_id5": "",         "conv\_unit2": 50,         "conv\_unit3": 1,         "conv\_unit4": 0,         "conv\_unit5": 0,         "purch\_price1": 70000,         "purch\_price2": 50000,         "purch\_price3": 30000,         "purch\_price4": 0,         "purch\_price5": 0,         "sell\_price1": 23000,         "sell\_price2": 800000,         "sell\_price3": 75000,         "sell\_price4": 0,         "sell\_price5": 0,         "vat": 11,         "pl\_id": 61,         "pl\_code": "001",         "pl\_name": "Fumakila",         "brand\_id": 104,         "brand\_code": "2",         "sbrand1\_id": 124,         "sbrand1\_code": "002",         "sbrand1\_name": "Vape Non Coil",         "vat\_value": 16500,         "amount": 166500,         "notes": "",         "item\_type": 1,         "disc\_value": 0,         "qty1\_stok": 0,         "qty2\_stok": 0,         "qty3\_stok": 66,         "qty1": 0,         "qty2": 0,         "qty3": 2       }     \],     "promo": null   } }  |
| :---- |


* BE hit promo consult  
  **URL Promo :** 

| /sales/v2/promotions/consult |
| :---- |

  Contoh req dan response : click [here](https://docs.google.com/spreadsheets/d/1UGsfcV0-Lwhi9rv6cNOdkXZwCOF3mo5Ei-BNhbowPYE/edit?usp=sharing)


  

| {   "order\_date": "2025-12-17",   "outlet\_id": 1404,   "salesman\_id": 359,   "wh\_id": 63,   "details": \[      {       "pro\_id": 709,       "qty1": 10,       "qty2": 0,       "qty3": 0,       "gross\_value": 2000000     },     {       "pro\_id": 710,       "qty1": 12,       "qty2": 0,       "qty3": 0,       "gross\_value": 3600000     },     {       "pro\_id": 711,       "qty1": 10,       "qty2": 0,       "qty3": 0,       "gross\_value": 800000     }       \] }  |
| :---- |


  

**Penjelasan :**

| Field | Data Source |
| :---- | :---- |
| order\_date | payload ro\_date |
| outlet\_id | outlet\_id |
| salesman\_id | salesman\_id |
| wh\_id | wh\_id |
| details\[\].pro\_id | details.normal\[\].pro\_id |
| details\[\].qty1 | details.normal\[\].qty1 |
| details\[\].qty2 | details.normal\[\].qty2 |
| details\[\].qty3 | details.normal\[\].qty3 |
| details\[\].gross\_value | (qty1 \* sell\_price1) \+(qty2 \* sell\_price2) \+ (qty3 \* sell\_price3) |

##### **Response Reward Percentage:** {#response-reward-percentage:}

jika reward\_value \= NULL   
**reward\_percentage** \= NOT NULL   
reward\_product \= NULL 

| {     "message": "Consulted V2 Successfully",     "data": \[         {             "promo\_id": "slab01",             "promo\_desc": "Syarat Value & Reward Percentage (Value)",             "slab\_id": "695ba0009f8c1325e300ae2f",             "slab\_desc": "slab 3 Value",             "slab\_reward": 3,             "slab\_reward\_uom": "",             "slab\_reward\_type": "percentage",             "slab\_per\_scope": "",             "total\_gross\_value": 5600000,             "products\_eligible": \[                 709,                 710             \],             "reward\_value": null,             "reward\_percentage": \[                 {                     "pro\_id": 709,                     "gross\_value": 2000000,                     "promo1": 60000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1940000                 },                 {                     "pro\_id": 710,                     "gross\_value": 3600000,                     "promo1": 108000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 3492000                 }             \],             "reward\_product": null         },   {             "promo\_id": "slab02",             "promo\_desc": "Syarat Value & Reward Percentage (Value)",             "slab\_id": "695ba0009f8c1325e300ae2f",             "slab\_desc": "slab 3 Value",             "slab\_reward": 3,             "slab\_reward\_uom": "",             "slab\_reward\_type": "percentage",             "slab\_per\_scope": "",             "total\_gross\_value": 5600000,             "products\_eligible": \[                 709,                 710             \],             "reward\_value": null,             "reward\_percentage": \[                 {                     "pro\_id": 709,                     "gross\_value": 2000000,                     "promo1": 60000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1940000                 },             \],             "reward\_product": null         }     \],     "request\_id": "fd9a72d5-b1af-4c75-83dd-a08a579d9877" }  |
| :---- |

######  *Update ke sls.order :*

| field | value  |
| :---- | :---- |
| promo\_remarks\_so | \[         {             "promo\_id": "slab01",             "promo\_desc": "Syarat Value & Reward Percentage (Value)",             "slab\_id": "695ba0009f8c1325e300ae2f",             "slab\_desc": "slab 3 Value",             "slab\_reward": 3,             "slab\_reward\_uom": "",             "slab\_reward\_type": "percentage",             "slab\_per\_scope": "",             "total\_gross\_value": 5600000,             "products\_eligible": \[                 709,                 710             \],             "reward\_value": null,             "reward\_percentage": \[                 {                     "pro\_id": 709,                     "gross\_value": 2000000,                     "promo1": 60000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1940000                 },                 {                     "pro\_id": 710,                     "gross\_value": 3600000,                     "promo1": 108000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 3492000                 }             \],             "reward\_product": null         }     \],  |
| promo\_remarks\_final |  |

######  *Update ke sls.order\_detail :* 

| field | value  |
| :---- | :---- |
| promo\_po1 | NULL (ini akan terisi apabila order via mobile)  |
| promo\_po2 |  |
| promo\_po3 |  |
| promo\_po4 |  |
| promo\_po5 |  |
| promo\_so1 | SUM(data\[\].reward\_percentage\[\].promo1) GROUP BY pro\_id |
| promo\_so2 | SUM(data\[\].reward\_percentage\[\].promo2) GROUP BY pro\_id |
| promo\_so3 | SUM(data\[\].reward\_percentage\[\].promo3) GROUP BY pro\_id |
| promo\_so4 | SUM(data\[\].reward\_percentage\[\].promo4) GROUP BY pro\_id |
| promo\_so5 | SUM(data\[\].reward\_percentage\[\].promo5) GROUP BY pro\_id |
| promo\_final1 | SUM(data\[\].reward\_percentage\[\].promo1) GROUP BY pro\_id |
| promo\_final2 | SUM(data\[\].reward\_percentage\[\].promo2) GROUP BY pro\_id |
| promo\_final3 | SUM(data\[\].reward\_percentage\[\].promo3) GROUP BY pro\_id |
| promo\_final4 | SUM(data\[\].reward\_percentage\[\].promo2) GROUP BY pro\_id |
| promo\_final5 | SUM(data\[\].reward\_percentage\[\].promo3) GROUP BY pro\_id |
| is\_product\_promotion\_so | FALSE |
| is\_product\_promotion\_po | NULL (ini akan terisi apabila order via mobile)   |
| is\_product\_promotion\_final | FALSE |
| promo\_remarks\_so | NULL (ini akan terisi apabila order via mobile)   |
| promo\_remarks\_po | data\[\].pro\_id  berdasarkan pro\_id pada contoh diatas , ekspektasinya adalah produk 709 \= {“"slab01", "slab02"} produk 710 \= {“"slab01"} |
| promo\_remarks\_final |  |

##### **Response Reward Value:** {#response-reward-value:}

jika **reward\_value** \= NOT NULL   
reward\_percentage \= NULL   
reward\_product \= NULL  

| {     "message": "Consulted V2 Successfully",     "data": \[         {             "promo\_id": "slab01",             "promo\_desc": "Syarat Value & Reward Value Perorder (Quantity)",             "slab\_id": "695deb11e3dc19398700bd68",             "slab\_desc": "slab 3 Quantity Smallest",             "slab\_reward": 20000,             "slab\_rule\_type": "value",             "slab\_rule\_uom": "",             "slab\_reward\_uom": "",             "slab\_reward\_type": "fixed\_value",             "slab\_per\_scope": "per\_order",             "total\_gross\_value": 3600000,             "products\_eligible": \[                 710,                 673             \],             "reward\_value": \[                 {                     "pro\_id": 710,                     "gross\_value": 1600000,                     "promo1": 10000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1590000                 },                 {                     "pro\_id": 673,                     "gross\_value": 2000000,                     "promo1": 10000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1990000                 }             \],             "reward\_percentage": null,             "reward\_product": null         },  {             "promo\_id": "slab02",             "promo\_desc": "Syarat Value & Reward Value Perorder (Quantity)",             "slab\_id": "695deb11e3dc19398700bd68",             "slab\_desc": "slab 3 Quantity Smallest",             "slab\_reward": 20000,             "slab\_rule\_type": "value",             "slab\_rule\_uom": "",             "slab\_reward\_uom": "",             "slab\_reward\_type": "fixed\_value",             "slab\_per\_scope": "per\_order",             "total\_gross\_value": 3600000,             "products\_eligible": \[                 710,                 673             \],             "reward\_value": \[                 {                     "pro\_id": 710,                     "gross\_value": 1600000,                     "promo1": 10000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1590000                 }             \],             "reward\_percentage": null,             "reward\_product": null         }     \],     "request\_id": "e1a30a69-4acc-44fb-81c6-c0b2500d8f9e" }  |
| :---- |

###### *Update ke sls.order :*

| field | value  |
| :---- | :---- |
| promo\_remarks\_so | data\[\].pro\_id  pada contoh diatas , ekspektasinya adalah {“"slab01", "slab02"} |
| promo\_remarks\_final |  |

######  *Update ke sls.order\_detail :* 

| field | value  |
| :---- | :---- |
| promo\_po1 | NULL (ini akan terisi apabila order via mobile)  |
| promo\_po2 |  |
| promo\_po3 |  |
| promo\_po4 |  |
| promo\_po5 |  |
| promo\_so1 | SUM(data\[\].reward\_value\[\].promo1) GROUP BY pro\_id |
| promo\_so2 | SUM(data\[\].reward\_value\[\].promo2) GROUP BY pro\_id |
| promo\_so3 | SUM(data\[\].reward\_value\[\].promo3) GROUP BY pro\_id |
| promo\_so4 | SUM(data\[\].reward\_value\[\].promo4) GROUP BY pro\_id |
| promo\_so5 | SUM(data\[\].reward\_value\[\].promo5) GROUP BY pro\_id |
| promo\_final1 | SUM(data\[\].reward\_value\[\].promo1) GROUP BY pro\_id |
| promo\_final2 | SUM(data\[\].reward\_value\[\].promo2) GROUP BY pro\_id |
| promo\_final3 | SUM(data\[\].reward\_value\[\].promo3) GROUP BY pro\_id |
| promo\_final4 | SUM(data\[\].reward\_value\[\].promo4) GROUP BY pro\_id |
| promo\_final5 | SUM(data\[\].reward\_value\[\].promo5) GROUP BY pro\_id |
| is\_produc\_promotion\_so | FALSE |
| is\_product\_promotion\_po | NULL (ini akan terisi apabila order via mobile)   |
| is\_product\_promotion\_final | FALSE |
| promo\_remarks\_so | NULL (ini akan terisi apabila order via mobile)   |
| promo\_remarks\_po | data\[\].pro\_id  berdasarkan pro\_id pada contoh diatas , ekspektasinya adalah produk 710 \= {“"slab01", "slab02"} produk 673 \= {“"slab01"} |
| promo\_remarks\_final |  |

##### **Response Reward Product:** {#response-reward-product:}

jika reward\_value \= NULL   
reward\_percentage \= NULL   
**reward\_product** \= NOT NULL 

| {     "message": "Consulted V2 Successfully",     "data": \[         {             "promo\_id": "slab04",             "promo\_desc": "Syarat Value & Reward Produk (Quantity)",             "slab\_id": "694caff8cc48b19f80002ceb",             "slab\_desc": "slab 1 Quantity Smallest",             "slab\_reward": 1,             "slab\_rule\_type": "value",             "slab\_rule\_uom": "",             "slab\_reward\_uom": "smallest",             "slab\_reward\_type": "product",             "slab\_per\_scope": "",             "total\_gross\_value": 1600000,             "products\_eligible": \[                 711,                 710             \],             "reward\_value": null,             "reward\_percentage": null,             "reward\_product": \[                 {                     "pro\_id": 491,                     "qty1": 1,                     "qty2": 0,                     "qty3": 0,                     "gross\_value": 10000,                     "promo1": 10000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0                 }             \]         },  {             "promo\_id": "slab05",             "promo\_desc": "Syarat Value & Reward Produk (Quantity)",             "slab\_id": "694caff8cc48b19f80002ceb",             "slab\_desc": "slab 1 Quantity Smallest",             "slab\_reward": 1,             "slab\_rule\_type": "value",             "slab\_rule\_uom": "",             "slab\_reward\_uom": "smallest",             "slab\_reward\_type": "product",             "slab\_per\_scope": "",             "total\_gross\_value": 1600000,             "products\_eligible": \[                 711,                 710             \],             "reward\_value": null,             "reward\_percentage": null,             "reward\_product": \[                 {                     "pro\_id": 491,                     "qty1": 1,                     "qty2": 0,                     "qty3": 0,                     "gross\_value": 10000,                     "promo1": 10000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0                 }             \]         }     \],     "request\_id": "1da40a69-4acc-44fb-81c6-c0b2500d8f9e" }  |
| :---- |

###### *Update ke sls.order :*

| field | value  |
| :---- | :---- |
| promo\_remarks\_so | data\[\].pro\_id  pada contoh diatas , ekspektasinya adalah {“"slab01", "slab02"} |
| promo\_remarks\_final |  |

######  *Update ke sls.order\_detail :* 

| field | value  |
| :---- | :---- |
| pro\_id | data.reward\_product\[\].pro\_id |
| qty\_po1 | NULL (ini akan terisi apabila order via mobile)   |
| qty\_po2 |  |
| qty\_po3 |  |
| promo\_po1 | NULL(ini akan terisi apabila order via mobile)   |
| promo\_po2 |  |
| promo\_po3 |  |
| promo\_po4 |  |
| promo\_po5 |  |
| promo\_remarks\_po | NULL(ini akan terisi apabila order via mobile)   |
| promo\_so1 | SUM(data\[\].reward\_product\[\].promo1) GROUP BY pro\_id |
| promo\_so2 | SUM(data\[\].reward\_product\[\].promo2) GROUP BY pro\_id |
| promo\_so3 | SUM(data\[\].reward\_product\[\].promo3) GROUP BY pro\_id |
| promo\_so4 | SUM(data\[\].reward\_product\[\].promo4) GROUP BY pro\_id |
| promo\_so5 | SUM(data\[\].reward\_product\[\].promo5) GROUP BY pro\_id |
| promo\_remarks\_so | data\[\].pro\_id  pada contoh diatas , ekspektasinya adalah {“"slab04", "slab05"} notes : hanya untuk promo yang berlaku di produk tersebut |
| qty1 | data\[\].reward\_product\[\].qty1 |
| qty2 | data\[\].reward\_product\[\].qty2 |
| qty3 | data\[\].reward\_product\[\].qty3 |
| promo\_final1 | SUM(data\[\].reward\_product\[\].promo1) GROUP BY pro\_id |
| promo\_final2 | SUM(data\[\].reward\_product\[\].promo2) GROUP BY pro\_id |
| promo\_final3 | SUM(data\[\].reward\_product\[\].promo3) GROUP BY pro\_id |
| promo\_final4 | SUM(data\[\].reward\_product\[\].promo4) GROUP BY pro\_id |
| promo\_final5 | SUM(data\[\].reward\_product\[\].promo5) GROUP BY pro\_id |
| promo\_remarks\_final | data\[\].pro\_id  pada contoh diatas , ekspektasinya adalah {“"slab04", "slab05"} notes : hanya untuk promo yang berlaku di produk tersebut |
| qty1\_final | data\[\].reward\_product\[\].qty1 |
| qty2\_final | data\[\].reward\_product\[\].qty2 |
| qty3\_final | data\[\].reward\_product\[\].qty3 |
| is\_produc\_promotion\_so | TRUE |
| is\_product\_promotion\_po | NULL (ini akan terisi apabila order via mobile)   |
| is\_product\_promotion\_final | TRUE |

###### 

######  *Insert ke inv.stock :* 

| cust\_id | stock\_id | stock\_date | tr\_code | tr\_no | wh\_id | pro\_id | item\_cdn | qty\_in | qty\_out | unit\_price | cogs | ref\_det\_id | created\_at | qty\_in\_order | qty\_out\_order |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| dari user yang login | generate dr BE | tanggal hari ini  | SO | SO2602090002 | dari request promo | pro\_id dari reward\_product  | 1 | 0 | konversi small dari response promo reward\_product (qty1, qty2 , qty3) | sell\_price satuan terkecil | dari mst.m\_product.cogs | dari id sls.order | 1770647314 | 0 | 0 |
| dari user yang login | generate dr BE | tanggal hari ini  | CO | SO2602090002-CO | dari request promo | pro\_id dari reward\_product  | 1 | 0 | 0 | sell\_price satuan terkecil | dari mst.m\_product.cogs  | dari id sls.order | 1770647314 | 0 | konversi small dari response promo reward\_product (qty1, qty2 , qty3) |

######  *Update ke inv.warehouse\_stock :* 

| cust\_id | wh\_id | pro\_id | qty | qty\_on\_order | qty\_on\_shipping | qty\_bs | qty\_exp | updated\_at |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| dari user yang login | dari request promo | pro\_id dari reward\_product  | qty sebelumnya \- (konversi small dari response promo reward\_product (qty1, qty2 , qty3)) | qty sebelumnya \+ (konversi small dari response promo reward\_product (qty1, qty2 , qty3)) | tidak ada update value | tidak ada update value | tidak ada update value | update  |

##### **Fixing rumus discount :** 

| Before :  Gross itemgross\_item \= (qty1\*sell\_price1) \+ (qty2\*sell\_price2) \+ (qty3\*sell\_price3)  Diskon item nominaldisc\_value\_item \= gross\_item \* (disc\_item / 100\)  |
| :---- |
| **After :  Gross item**gross\_item \= (qty1\*sell\_price1) \+ (qty2\*sell\_price2) \+ (qty3\*sell\_price3)  **Promo**promo \= (promo1 \+ promo2 \+ promo3+ promo4 \+ promo5) **Diskon item nominal**disc\_value\_item \= (gross\_item \- promo)  \* (disc\_item / 100\)  |

3. ### **Edit Sales Order  \- Tab Sales Order** 

     
   **URL** 		: [https://best.scyllax.online/sales/v1/orders/enhance/SO2603070001](https://best.scyllax.online/sales/v1/orders/enhance/SO2603070001)  
   **Method** 	: PATCH   
   **Payload** 	: 

| {   "sales\_order": \[     {       "order\_detail\_id": 5642,       "qty1": 0,       "qty2": 0,       "qty3": 3,       "sell\_price1": 200000,       "sell\_price2": 500000,       "sell\_price3": 1000000,        "is\_product\_promotion\_so" : true / false     }   \] }  |
| :---- |

   

   **Enhance :** 

| ![][image20] |
| :---- |
| **sequenceDiagram**  actor User  participant FE as Frontend  participant SO as Sales Order API  participant Promo as Promotion API  participant DB as Database   User**\-\>\>**FE**:** Edit Sales Order (items, qty, price)  FE**\-\>\>**SO**:** PATCH /sales/v1/orders/enhance/SO2603070001  Note right of FE**:** items, qty, price, outlet\_id, salesman\_id  SO**\-\>\>**DB**:** Get existing order\_detail  Note right of SO**:** Retrieve previous promo bonus qty  SO**\-\>\>**Promo**:** GET /sales/v2/promotions/consult  Note right of SO**:** Send item list, qty, outlet, date  Promo**\--\>\>**SO**:** Return eligible promotions  Note left of Promo**:** discount / bonus / reward  SO**\-\>\>**SO**:** Calculate promotion result  SO**\-\>\>**SO**:** Compare old promo vs new promo  **alt** Promo bonus increased      SO**\-\>\>**DB**:** Reduce stock for additional bonus      Note right of DB**:** Update inv.stock Update inv.warehouse\_stock  **else** Promo bonus decreased      SO**\-\>\>**DB**:** Return stock from removed bonus      Note right of DB**:** Update inv.stock Update inv.warehouse\_stock  **else** Promo unchanged      SO**\-\>\>**SO**:** No stock adjustment  **end**   SO**\-\>\>**DB**:** Update sls.order  Note right of DB**:** header order\\n(order\_no, outlet\_id, total\_amount, promo\_summary)  SO**\-\>\>**DB**:** Update sls.order\_detail (loop items)  DB**\--\>\>**SO**:** Success  SO**\--\>\>**FE**:** Order Updated Response  FE**\--\>\>**User**:** Show Order Success   |

   

   

###### 

* ###### *Update promo product ke sls.order & sls.order\_detail* 

- ###### Reward Percentage \= sama dengan [create order](#response-reward-percentage:)

- Reward Value \= sama dengan [create order](#response-reward-value:)  
- Reward Product \= sama dengan [create order](#response-reward-product:)

* ###### *Khusus reward product, update promo product ke inv.stock :* 

| cust\_id | stock\_id | stock\_date | tr\_code | tr\_no | wh\_id | pro\_id | item\_cdn | qty\_in | qty\_out | unit\_price | cogs | ref\_det\_id | created\_at | qty\_in\_order | qty\_out\_order |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| dari user yang login | generate dr BE | tanggal hari ini  | SO | SO2602090002 | dari request promo | pro\_id dari reward\_product  | 1 | 0 | konversi small dari **perbandingan qty sebelumnya dg**  **response promo** reward\_product (qty1, qty2 , qty3) | sell\_price satuan terkecil | dari mst.m\_product.cogs | dari id sls.order | 1770647314 | 0 | 0 |
| dari user yang login | generate dr BE | tanggal hari ini  | CO | SO2602090002-CO | dari request promo | pro\_id dari reward\_product  | 1 | 0 | 0 | sell\_price satuan terkecil | dari mst.m\_product.cogs  | dari id sls.order | 1770647314 | 0 | konversi small dari **perbandingan qty sebelumnya dg**  **response promo** reward\_product (qty1, qty2 , qty3) |


* ###### *Update promo product ke inv.warehouse\_stock :* 

| cust\_id | wh\_id | pro\_id | qty | qty\_on\_order | qty\_on\_shipping | qty\_bs | qty\_exp | updated\_at |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| dari user yang login | dari request promo | pro\_id dari reward\_product  | qty sebelumnya \- (konversi small dari **perbandingan qty sebelumnya dg**  **response promo** reward\_product (qty1, qty2 , qty3)) | qty sebelumnya \+ (konversi small dari **perbandingan qty sebelumnya dg**  **response promo** reward\_product (qty1, qty2 , qty3)) | tidak ada update value | tidak ada update value | tidak ada update value | update  |


  

##### **Fixing rumus discount :** 

| Before :  Gross itemgross\_item \= (qty1\*sell\_price1) \+ (qty2\*sell\_price2) \+ (qty3\*sell\_price3)  Diskon item nominaldisc\_value\_item \= gross\_item \* (disc\_item / 100\)  |
| :---- |
| **After :  Gross item**gross\_item \= (qty1\*sell\_price1) \+ (qty2\*sell\_price2) \+ (qty3\*sell\_price3)  **Promo**promo \= (promo1 \+ promo2 \+ promo3+ promo4 \+ promo5) **Diskon item nominal**disc\_value\_item \= (gross\_item \- promo)  \* (disc\_item / 100\)  |

4. ### **Edit Sales Order  \- Tab Final Order** 

     
   **URL** 		: [https://best.scyllax.online/sales/v1/orders/enhance/SO2603070001](https://best.scyllax.online/sales/v1/orders/enhance/SO2603070001)  
   **Method** 	: PATCH   
   **Payload** 	: 

| {   "final\_order": \[     {       "order\_detail\_id": 5609,       "qty1\_final": 1,       "qty2\_final": 0,       "qty3\_final": 2,       "sell\_price\_final1": 600000,       "sell\_price\_final2": 6000000,       "sell\_price\_final3": 6000000,        "is\_produc\_promotion\_final" : true / false     }   \] }  |
| :---- |

   

   **Enhance :** 

| ![][image20] |
| :---- |
| **sequenceDiagram**  actor User  participant FE as Frontend  participant SO as Sales Order API  participant Promo as Promotion API  participant DB as Database   User**\-\>\>**FE**:** Edit Sales Order (items, qty, price)  FE**\-\>\>**SO**:** PATCH /sales/v1/orders/enhance/SO2603070001  Note right of FE**:** items, qty, price, outlet\_id, salesman\_id  SO**\-\>\>**DB**:** Get existing order\_detail  Note right of SO**:** Retrieve previous promo bonus qty  SO**\-\>\>**Promo**:** GET /sales/v2/promotions/consult  Note right of SO**:** Send item list, qty, outlet, date  Promo**\--\>\>**SO**:** Return eligible promotions  Note left of Promo**:** discount / bonus / reward  SO**\-\>\>**SO**:** Calculate promotion result  SO**\-\>\>**SO**:** Compare old promo vs new promo  **alt** Promo bonus increased      SO**\-\>\>**DB**:** Reduce stock for additional bonus      Note right of DB**:** Update inv.stock Update inv.warehouse\_stock  **else** Promo bonus decreased      SO**\-\>\>**DB**:** Return stock from removed bonus      Note right of DB**:** Update inv.stock Update inv.warehouse\_stock  **else** Promo unchanged      SO**\-\>\>**SO**:** No stock adjustment  **end**   SO**\-\>\>**DB**:** Update sls.order  Note right of DB**:** header order\\n(order\_no, outlet\_id, total\_amount, promo\_summary)  SO**\-\>\>**DB**:** Update sls.order\_detail (loop items)  DB**\--\>\>**SO**:** Success  SO**\--\>\>**FE**:** Order Updated Response  FE**\--\>\>**User**:** Show Order Success   |

   

   

###### 

* ###### *Update promo product ke sls.order\_detail* 

###### *Update ke sls.order :*

| field | value  |
| :---- | :---- |
| promo\_remarks\_final | \[         {             "promo\_id": "slab01",             "promo\_desc": "Syarat Value & Reward Percentage (Value)",             "slab\_id": "695ba0009f8c1325e300ae2f",             "slab\_desc": "slab 3 Value",             "slab\_reward": 3,             "slab\_reward\_uom": "",             "slab\_reward\_type": "percentage",             "slab\_per\_scope": "",             "total\_gross\_value": 5600000,             "products\_eligible": \[                 709,                 710             \],             "reward\_value": null,             "reward\_percentage": \[                 {                     "pro\_id": 709,                     "gross\_value": 2000000,                     "promo1": 60000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 1940000                 },                 {                     "pro\_id": 710,                     "gross\_value": 3600000,                     "promo1": 108000,                     "promo2": 0,                     "promo3": 0,                     "promo4": 0,                     "promo5": 0,                     "net\_value": 3492000                 }             \],             "reward\_product": null         }     \],  |

######  *Update ke sls.order\_detail :* 

| field | value  |
| :---- | :---- |
| promo\_final1 | SUM(data\[\].reward\_percentage\[\].promo1) GROUP BY pro\_id |
| promo\_final2 | SUM(data\[\].reward\_percentage\[\].promo2) GROUP BY pro\_id |
| promo\_final3 | SUM(data\[\].reward\_percentage\[\].promo3) GROUP BY pro\_id |
| promo\_final4 | SUM(data\[\].reward\_percentage\[\].promo2) GROUP BY pro\_id |
| promo\_final5 | SUM(data\[\].reward\_percentage\[\].promo3) GROUP BY pro\_id |
| is\_product\_promotion\_final | FALSE |
| promo\_remarks\_final | data\[\].pro\_id  berdasarkan pro\_id pada contoh diatas , ekspektasinya adalah produk 709 \= {“"slab01", "slab02"} produk 710 \= {“"slab01"} |

* ###### *Update promo product ke inv.stock :* 

| cust\_id | stock\_id | stock\_date | tr\_code | tr\_no | wh\_id | pro\_id | item\_cdn | qty\_in | qty\_out | unit\_price | cogs | ref\_det\_id | created\_at | qty\_in\_order | qty\_out\_order |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| dari user yang login | generate dr BE | tanggal hari ini  | SO | SO2602090002 | dari request promo | pro\_id dari reward\_product  | 1 | 0 | konversi small dari **perbandingan qty sebelumnya dg**  **response promo** reward\_product (qty1, qty2 , qty3) | sell\_price satuan terkecil | dari mst.m\_product.cogs | dari id sls.order | 1770647314 | 0 | 0 |
| dari user yang login | generate dr BE | tanggal hari ini  | CO | SO2602090002-CO | dari request promo | pro\_id dari reward\_product  | 1 | 0 | 0 | sell\_price satuan terkecil | dari mst.m\_product.cogs  | dari id sls.order | 1770647314 | 0 | konversi small dari **perbandingan qty sebelumnya dg**  **response promo** reward\_product (qty1, qty2 , qty3) |


* ###### *Update promo product ke inv.warehouse\_stock :* 

| cust\_id | wh\_id | pro\_id | qty | qty\_on\_order | qty\_on\_shipping | qty\_bs | qty\_exp | updated\_at |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| dari user yang login | dari request promo | pro\_id dari reward\_product  | qty sebelumnya \- (konversi small dari **perbandingan qty sebelumnya dg**  **response promo** reward\_product (qty1, qty2 , qty3)) | qty sebelumnya \+ (konversi small dari **perbandingan qty sebelumnya dg**  **response promo** reward\_product (qty1, qty2 , qty3)) | tidak ada update value | tidak ada update value | tidak ada update value | update  |


  

##### **Fixing rumus discount :** 

| Before :  Gross itemgross\_item \= (qty1\*sell\_price1) \+ (qty2\*sell\_price2) \+ (qty3\*sell\_price3)  Diskon item nominaldisc\_value\_item \= gross\_item \* (disc\_item / 100\)  |
| :---- |
| **After :  Gross item**gross\_item \= (qty1\*sell\_price1) \+ (qty2\*sell\_price2) \+ (qty3\*sell\_price3)  **Promo**promo \= (promo1 \+ promo2 \+ promo3+ promo4 \+ promo5) **Diskon item nominal**disc\_value\_item \= (gross\_item \- promo)  \* (disc\_item / 100\)  |

#  **ISSUE** 

1. ### **SX-521, SX-781**

   [https://scyllax-pratesis.atlassian.net/browse/SX-521](https://scyllax-pratesis.atlassian.net/browse/SX-521)   
     
   Handle by BE :   
* **inv.order** 

  add field 

#### 

| Nama Kolom | Tipe Data | *Nullable* | Nilai *Default* | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| address1 | varchar(150) |  | \- | adress dari mst.m\_outlet |


* [https://best.scyllax.online/sales/v1/orders/](https://best.scyllax.online/sales/v1/orders/) 

  POST 

  ketika di hit, BE save **outlet\_address1** (dari payload)  → disimpan pada field address1


* [https://be.scyllax.online/sales/v1/orders?page=1\&limit=20](https://be.scyllax.online/sales/v1/orders?page=1&limit=20)   
  GET   
  ubah resp **outlet\_address1**  diambil dari inv.order.address1  
    
* [https://best.scyllax.online/sales/v2/orders/SO2601220008](https://best.scyllax.online/sales/v2/orders/SO2601220008)   
  GET   
  ubah resp **outlet\_address1**  diambil dari inv.order.address1  
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
  


2. ### **SX-574**

[https://scyllax-pratesis.atlassian.net/browse/SX-574](https://scyllax-pratesis.atlassian.net/browse/SX-574) 

No SO 		: SO2601120001  
Product 		: Lampu Depan Astrea LEDID: BS00040005

![][image21]

**Add New Order**  
![][image22]

1) #### Stock Awal 

| Process | Warehouse Stock  |  |  | On Cust Order |  |  |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- |
|  | **Largest**  | **Middle** | **Small** | **Largest**  | **Middle** | **Small** |
| Saldo Awal  | 970 | 0 | 0 | 30 | 0 | 0 |

2) Add Product 

Warehouse Stock diambil dari warehoustock \= 970 → **EXPECTED**  
Input QTY Order \= 2 0 0   
![][image23]![][image24]

3) Process Order 

| Process | Warehouse Stock  |  |  | On Cust Order |  |  |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- |
|  | **Largest**  | **Middle** | **Small** | **Largest**  | **Middle** | **Small** |
| Saldo Awal  | 970 | 0 | 0 | 30 | 0 | 0 |
| Proses no 3 (add kemudian process order)  | 968 | 0 | 0 | 32 | 0 | 0 |

**![][image25]**

4) Edit Order 

Ekspektasi : 

- Available stock \= wh\_stock 968 \+ qty\_order 2= 970  
- Resp BE : [https://best.scyllax.online/inventory/v1/stocks/report?page=1\&limit=10\&q=\&wh\_id=1\&pro\_id=536\&outlet\_id=246\&order\_date=2026-01-12\&include\_zero\_stock=true\&active\_product\_only=true\&sort=pro\_code:asc](https://best.scyllax.online/inventory/v1/stocks/report?page=1&limit=10&q=&wh_id=1&pro_id=536&outlet_id=246&order_date=2026-01-12&include_zero_stock=true&active_product_only=true&sort=pro_code:asc)   
  ![][image26]

![][image27]

5) Edit Order 

![][image28]  
Before qty order \= 2   0  0   
After qty order    \= 10  0  0   
Process Order 

| Process | Warehouse Stock  |  |  | On Cust Order |  |  |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- |
|  | **Largest**  | **Middle** | **Small** | **Largest**  | **Middle** | **Small** |
| Saldo Awal  | 970 | 0 | 0 | 30 | 0 | 0 |
| Proses no 3 (add  2 kemudian process order)  | 968 | 0 | 0 | 32 | 0 | 0 |
| Proses no 4 Edit 10  | 960 | 0 | 0 | 40 | 0 | 0 |

![][image29]

Wh Stock After Process  
![][image30]

3. ### **Issue \- Lookup Product**

* **API** : [https://best.scyllax.online/inventory/v1/stocks/report?page=1\&limit=10\&q=\&wh\_id=63\&pro\_id=711\&outlet\_id=865\&order\_date=2026-01-19\&include\_zero\_stock=true\&active\_product\_only=true\&sort=pro\_code:asc](https://best.scyllax.online/inventory/v1/stocks/report?page=1&limit=10&q=&wh_id=63&pro_id=711&outlet_id=865&order_date=2026-01-19&include_zero_stock=true&active_product_only=true&sort=pro_code:asc)   
* **Method** :  GET   
* **Payload** : page=1\&limit=10\&q=\&wh\_id=63\&pro\_id=711\&outlet\_id=865\&order\_date=2026-01-19\&include\_zero\_stock=true\&active\_product\_only=true\&sort=pro\_code:asc  
* **Enhance :**   
- qty1  
- qty2  
- qty3

Konversi **inv.warehouse\_stock.qty**  menjadi S, M, L ***(code eksisting)***  
masing” di tambah sls.order\_detail.qty1(L) , sls.order\_detail.qty2(M) , sls.order\_detail.qty3 (S)

* **Response** : 

{  
    "message": "",  
    "data": \[  
        {  
            "pro\_id": 711,  
            "pro\_code": "DD-FM02-0003",  
            "pro\_name": "Hometown Choco 450ml SP1",  
            "unit\_id1": "PCS",  
            "unit\_id2": "BOX",  
            "unit\_id3": "BOX",  
            "purch\_price1": 0,  
            "purch\_price2": 0,  
            "purch\_price3": 0,  
            "sell\_price1": 0,  
            "sell\_price2": 0,  
            "sell\_price3": 0,  
            "total\_qty": 91,  
            "qty1": 19,  
            "qty2": 1,  
            "qty3": 1,  
            "total\_qty\_order": 735,  
            "qty\_order1": 15,  
            "qty\_order2": 0,  
            "qty\_order3": 15,  
            "total\_qty\_inc\_on\_order": 1470,  
            "qty\_inc\_on\_order1": 10,  
            "qty\_inc\_on\_order2": 0,  
            "qty\_inc\_on\_order3": 17,  
            "conv\_unit2": 24,  
            "conv\_unit3": 2,  
            "is\_active": true,  
            "vat": 0,  
            "vat\_lg\_purch": 0,  
            "vat\_lg\_sell": 0  
        }  
    \],  
    "paging": {  
        "total\_record": 1,  
        "page\_current": 1,  
        "page\_limit": 10,  
        "page\_total": 1  
    },  
    "filter": {  
        "cust\_id": "C220010001",  
        "parent\_cust\_id": "C22001",  
        "date": "",  
        "page": 1,  
        "limit": 10,  
        "q": "",  
        "sort": "pro\_code:asc",  
        "pro\_id": \[  
            711  
        \],  
        "wh\_id": \[  
            63  
        \],  
        "sup\_id": null,  
        "show\_price": "",  
        "include\_zero\_stock": "true",  
        "active\_product\_only": "true",  
        "brand\_id": null,  
        "pcat\_id": null,  
        "pl\_id": null  
    },  
    "request\_id": "696e4f3d1ab02641cac1ebdc"  
}

4. ### **SX-980 | Issue \- Edit Order (Final Order Tab)** 

   [https://scyllax-pratesis.atlassian.net/browse/SX-980](https://scyllax-pratesis.atlassian.net/browse/SX-980)   
     
   Stock awal   
   ![][image31]  
     
   **wh\_id : 301**  
   **conv\_unit2 \= 10**   
   **conv\_unit3 \= 1** 

| Note  | pro\_id | wh\_stock  |  |  | conversi S | on\_cust\_order |  |  | conversi S |
| :---- | :---- | ----- | :---- | :---- | :---- | ----- | :---- | :---- | :---- |
|  |  | **L** | **M** | **S** |  | **L** | **M** | **S** |  |
| 1.Saldo Awal  | 748 | 1 | 0 | 3 | 13 | 1 | 0 | 0 | 10 |
| 2\. process order (final order)  |  | 0 | 0 | 2 | 2 | 2 | 0 | 1 | 21 |
|  |  |  |  |  |  |  |  |  |  |

* Sales order sebelum di edit   
  ![][image32]  
    
* sales order setelah di edit   
  2 0 1   
  process order 

| curl 'https://best.scyllax.online/sales/v1/orders/enhance/SO2603060001' \\   \-X 'PATCH' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsImVtcF9pZCI6MCwiZXhwaXJlcyI6MTc3Mjg1MDY4NiwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4NTc0ODEyMzEyOCIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IkRpc3QgSURFIFNkYSIsInVzZXJfaWQiOjEyLCJ1c2VyX25hbWUiOiJkaXN0QHNkYS5pZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODU3NDgxMjMxMjgifQ.G2m6WzCuusczt04Goik1AxiaPDnFofxlbLx7-5ND5gw' \\   \-H 'Connection: keep-alive' \\   \-H 'Content-Type: application/json' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Not:A-Brand";v="99", "Google Chrome";v="145", "Chromium";v="145"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"' \\   \--data-raw '{"final\_order":\[{"order\_detail\_id":5609,"qty1\_final":1,"qty2\_final":0,"qty3\_final":2,"sell\_price\_final1":600000,"sell\_price\_final2":6000000,"sell\_price\_final3":6000000}\]}'  |
| :---- |
| **payload :**  {   "final\_order": \[     {       "order\_detail\_id": 5609,       "qty1\_final": 1,       "qty2\_final": 0,       "qty3\_final": 2,       "sell\_price\_final1": 600000,       "sell\_price\_final2": 6000000,       "sell\_price\_final3": 6000000     }   \] }  |


  

  Issue : 

  ![][image33]

- Perhitungan perubahan qty  
  **data sebelumny**a \= sls.order\_detail.qty1\_final (0), qty2\_final (0), qty3\_final (1)  
  **payload** \= qty1\_final (1), qty2\_final (0), qty3\_final (2)  
  **perubahan \=** Small 1 , Medium 0 , Large 1 

  **hitung konversi \=** menggunakan function eksisting   
  untuk data ini seharusnya 11 

  data di mst.m\_product \= 

- **conv\_unit2 \= 10**   
- **conv\_unit3 \= 1** 

- saat insert ke inv.stcok   
  perhatikan pada baris ke 3 dan ke 4 

  value masih kosong sehingga di warehouse stock data tidak berubah 

  seharusnya   
  apabila data awal 1 0 0 \= qty\_out 10 dan qty\_out\_order 10 

  saat process final menjadi 2 0 1  
  **artinya \=** perubahan stock adalah 1 0 1 \= sehingga value **qty\_out** dan **qty\_out\_order** yang seharusnya adalah 

  11   
  (11 ini nilainya dari konversi di mst.m\_product)


- saat update ke inv.warehouse\_stock    
  data sebelumnya : ![][image34]  
  enhance :  
  data awal \= 13   
  perubahan qty\_order \= 11  
  **field qty** \=  data awal \- perubahan qty\_order  
  13 \- 11  \= 2   
  **qty\_on\_order** \= data awal \+ perubahan qty\_order  
  10 \+ 11 \= 21 




5. ### **SX-1241 | Cancel Order not updating stock for On Cust Order (Tab Sales Order)** 

   [https://scyllax-pratesis.atlassian.net/browse/SX-1241](https://scyllax-pratesis.atlassian.net/browse/SX-1241)   
   Reference doc : click [here](#cancel-order)  
     
* Data sebelum di cancel 

  ![][image35]


  Product Code 		\= 02006

  120 KRT | 0 KRT | 16 PCS

  ![][image36]


  SO2603090005

* CURL : 

| curl 'https://best.scyllax.online/sales/v1/orders/status' \\   \-X 'PATCH' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoiZmlybWFudG9rMTk4OUBnbWFpbC5jb20iLCJlbXBfaWQiOjYyLCJleHBpcmVzIjoxNzczMTE1MDkzLCJpc19hZG1pbiI6ZmFsc2UsImxhbmdfaWQiOiJpZCIsIm1vYmlsZV9ubyI6IjA4NTc0ODEyMzEyOCIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IkRpc3QgSURFIFNkYSIsInVzZXJfaWQiOjEzLCJ1c2VyX25hbWUiOiJmaXJtYW50b2sxOTg5QGdtYWlsLmNvbSIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.irahl1SnGaKZs6fnvDmTtsnxjUPf4dJKFLMPLOOTNY4' \\   \-H 'Connection: keep-alive' \\   \-H 'Content-Type: application/json' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Not:A-Brand";v="99", "Google Chrome";v="145", "Chromium";v="145"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"' \\   \--data-raw '{"orders":\[{"ro\_no":"SO2603090005","data\_status":9}\]}'  |
| :---- |


  

* Update ke inv.stock dan inv.warehouse\_stock 

untuk qty , perhatikan untuk masing-masing kebutuhan tab 

- purchase\_order : qty\_po1, qty\_po2, qty\_po3 **(prioritas 3\)**   
- sales\_order : qty1, qty2, qty3 **(prioritas 2\)**   
- final\_order : qty1\_final, qty2\_final, qty3\_final  **(prioritas 1\)**

* cek apabila ada di **sales\_order** dan tidak ada di **purchase\_order** & **final\_order,** maka update qty stock sesuai data sales\_order 

* cek apabila ada data di **sales\_order** dan **purchase\_order** , sedangkan final\_order null, maka update qty stock sesuai data **sales\_order** 

* cek apabila ada data di **purchase\_order** sedangkan **sales\_order** & **final\_order** null, maka update qty stock sesuai data **purchase\_order**  
* cek apabila ada data di purchase\_order ,  sales\_order, final\_order,   
  maka update qty stock sesuai data **final\_order**

kesimpulan , data di update berdasarkan urutan prioritas. 

| cust\_id | stock\_id | stock\_date | tr\_code | tr\_no | wh\_id | pro\_id | item\_cdn | qty\_in | qty\_out | unit\_price | cogs | ref\_det\_id | created\_at | qty\_in\_order | qty\_out\_order |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| dari user yang login | generate dr BE | tanggal hari ini  | SO | ro\_no (payload) | wh\_id sesuai ro\_no (cek di sls.order) | pro\_id dar ro\_no ,cek di sls.order\_detaili   | 1 | ambil dari konversi small \= qty1\_final, qty2\_final, qty3\_final atau qty1, qty2, qty3 atau qty\_po1, qty\_po2, qty\_po3 | 0 | sell\_price satuan terkecil | dari mst.m\_product.cogs | dari id sls.order | 1770647314 | 0 | 0 |
| dari user yang login | generate dr BE | tanggal hari ini  | CO | SO2602090002-CO | dari request promo | pro\_id dari reward\_product  | 1 | 0 | 0 | sell\_price satuan terkecil | dari mst.m\_product.cogs  | dari id sls.order | 1770647314 | ambil dari konversi small \=  qty1\_final, qty2\_final, qty3\_final atau qty1, qty2, qty3 atau qty\_po1, qty\_po2, qty\_po3 |  |

* ###### *Update promo product ke inv.warehouse\_stock :* 

| cust\_id | wh\_id | pro\_id | qty | qty\_on\_order | qty\_on\_shipping | qty\_bs | qty\_exp | updated\_at |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| dari user yang login | dari request promo | pro\_id dari reward\_product  | qty sebelumnya \+ (konversi small  qty1\_final, qty2\_final, qty3\_final atau qty1, qty2, qty3 atau qty\_po1, qty\_po2, qty\_po3)  | qty sebelumnya \- (konversi small  qty1\_final, qty2\_final, qty3\_final atau qty1, qty2, qty3 atau qty\_po1, qty\_po2, qty\_po3) | tidak ada update value | tidak ada update value | tidak ada update value | update  |


  


  


  


  


  


  


  


  


  


  


  


  


6. ### **SX-755 | Add product then delete existing product then process order will appears Qty Order, Initial Stock and PPN with 0 value**

   [https://scyllax-pratesis.atlassian.net/browse/SX-755](https://scyllax-pratesis.atlassian.net/browse/SX-755)  
   

   #### **Stock AF-0003:** 

* ##### **Stock Awal :** 

- wh\_stock : 43 0 2   
- on\_cust\_order : 29 0 0    
  ![][image37]


* ##### **Step 2 (setelah [proses](#step-2-process-order-\(201-\)) order)** 

- wh\_stock : 41 0 1  
- on\_cust\_order : 31 0 1

![][image38]

* ##### **Step 3 ([setelah edit  order dan process order)](#step-3-edit-order-:)** 

![][image39]

#### **Stock AF-0038:** 

* ##### **Stock Awal :** 

- wh\_stock : 67 0 0   
- on cust order : 31 0 9 

![][image40]

* ##### [**Skenario step 3  :**](#skenario-step-3-:)  {#skenario-step-3-:}

![][image41]

#### **Transaksi Sales Order**

* ##### **Step 1 : Create new order**

  ![][image42]


* ##### **Step 2 process order (201 )**  {#step-2-process-order-(201-)}

* ##### **Step 3 Edit Order :**  {#step-3-edit-order-:}

- Hapus action figure spiderman   
- Tambah produk baru   
  ![][image43]

![][image44]

- CURL : 

| curl 'https://best.scyllax.online/sales/v1/orders/enhance/SO2603130012' \\   \-X 'PATCH' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsImVtcF9pZCI6MCwiZXhwaXJlcyI6MTc3MzQ2NzcwOSwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4NTc0ODEyMzEyOCIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IkRpc3QgSURFIFNkYSIsInVzZXJfaWQiOjEyLCJ1c2VyX25hbWUiOiJkaXN0QHNkYS5pZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODU3NDgxMjMxMjgifQ.i6m6j5HbaxjbU8XClGqPRuaaAiEOLY7YaJyVC74\_d4w' \\   \-H 'Connection: keep-alive' \\   \-H 'Content-Type: application/json' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Not:A-Brand";v="99", "Google Chrome";v="145", "Chromium";v="145"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"' \\   \--data-raw '{"sales\_order":\[{"order\_detail\_id":5797,"qty1":0,"qty2":0,"qty3":0,"sell\_price1":200000,"sell\_price2":2000000,"sell\_price3":2000000}\],"add\_sales\_order":\[{"pro\_id":751,"qty1":0,"qty2":0,"qty3":2,"sell\_price\_system1":200000,"sell\_price\_system2":2000000,"sell\_price\_system3":2000000,"sell\_price1":200000,"sell\_price2":2000000,"sell\_price3":2000000,"unit\_id1":"PCS","unit\_id2":"CRT","unit\_id3":"CRT"}\]}'  |
| :---- |


  


  


  


  


  


  


  


  


7. ### **SX-1353 | Conversion not working well**

     
   [https://scyllax-pratesis.atlassian.net/browse/SX-1353](https://scyllax-pratesis.atlassian.net/browse/SX-1353)   
     
   

1) #### **Stock Awal** 

   AF-0006 

   ![][image45]

   

   ![][image46]

   conv\_unit2 \= 10 

   conv\_unit3 \= 1 

| Case | Warehouse Stock  |  |  | On Cust Order |  |  |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- |
|  | **Large**  | **Middle** | **Smallest** | **Large**  | **Middle** | **Smallest** |
| Saldo Awal  | 30 | 0 | 0 | 0 | 0 | 0 |
|  |  |  |  |  |  |  |
|  |  |  |  |  |  |  |

   

   

2) #### **Create Order**

   ![][image47]

   

   order : small 3 

   

   **impact ke stock** 

   stock small \= 300 

   qty order small \= 3 

   sisa stock small \= 297 (No Issue ) 

   ![][image48]

   

   

   available stock (No Issue) 

   ![][image49]

   

   

   

   

   

   

   

   

   

   

8. ### **SX-1402 | Download Sales Order**

   yang ini [https://scyllax-pratesis.atlassian.net/browse/SX-1402](https://scyllax-pratesis.atlassian.net/browse/SX-1402)   
   

| SELECT    *calc*.\*,    (        *calc*.*gross\_price*        \- COALESCE(*calc*.*promotion\_value*,0)        \- COALESCE(*calc*.*discount\_value*,0)    ) AS *net\_sales* FROM (    SELECT        *base*.\*,        (            COALESCE(*base*.*final\_largest\_selling\_price*,0) \* COALESCE(*base*.*largest\_qty\_order*,0) \+            COALESCE(*base*.*final\_middle\_selling\_price*,0) \* COALESCE(*base*.*middle\_qty\_order*,0) \+            COALESCE(*base*.*final\_small\_selling\_price*,0) \* COALESCE(*base*.*smallest\_qty\_order*,0)        ) AS *gross\_price*    FROM (        SELECT            *o*.order\_no ,            *o*.ro\_no ,            *o*.ro\_date ,            *o*.invoice\_date ,            *o*.invoice\_no ,            *mo*.outlet\_code ,            *mo*.outlet\_name ,            *me*.emp\_code ,            *me*.emp\_name ,            *ms2*.sup\_code ,            *ms2*.sup\_name ,            *mp*.pro\_code ,            *mp*.pro\_name ,            CONCAT(*mp*.unit\_id3,'/',*mp*.conv\_unit2) AS *largest\_unit*,            CONCAT(*mp*.unit\_id2,'/',*mp*.conv\_unit3) AS *middle\_unit*,            CONCAT(*mp*.unit\_id1,'/','1') AS *small\_unit*,            COALESCE(*mp*.sell\_price3,0) AS *largest\_system\_price*,            COALESCE(*mp*.sell\_price2,0) AS *middle\_system\_price*,            COALESCE(*mp*.sell\_price1,0) AS *small\_system\_price*,            COALESCE(*od*.sell\_price\_final3, *od*.sell\_price\_po3, *od*.sell\_price3,0) AS *final\_largest\_selling\_price*,            COALESCE(*od*.sell\_price\_final2, *od*.sell\_price\_po2, *od*.sell\_price2,0) AS *final\_middle\_selling\_price*,            COALESCE(*od*.sell\_price\_final1, *od*.sell\_price\_po1, *od*.sell\_price1,0) AS *final\_small\_selling\_price*,            COALESCE(*od*.qty3\_final, *od*.qty\_po3, *od*.qty3,0) AS *largest\_qty\_order*,            COALESCE(*od*.qty2\_final, *od*.qty\_po2, *od*.qty2,0) AS *middle\_qty\_order*,            COALESCE(*od*.qty1\_final, *od*.qty\_po1, *od*.qty1,0) AS *smallest\_qty\_order*,            CASE                WHEN                    COALESCE(*od*.promo\_final1,0)+COALESCE(*od*.promo\_final2,0)+                    COALESCE(*od*.promo\_final3,0)+COALESCE(*od*.promo\_final4,0)+                    COALESCE(*od*.promo\_final5,0) \> 0                THEN                    COALESCE(*od*.promo\_final1,0)+COALESCE(*od*.promo\_final2,0)+                    COALESCE(*od*.promo\_final3,0)+COALESCE(*od*.promo\_final4,0)+                    COALESCE(*od*.promo\_final5,0)                WHEN                    COALESCE(*od*.promo\_so1,0)+COALESCE(*od*.promo\_so2,0)+                    COALESCE(*od*.promo\_so3,0)+COALESCE(*od*.promo\_so4,0)+                    COALESCE(*od*.promo\_so5,0) \> 0                THEN                    COALESCE(*od*.promo\_so1,0)+COALESCE(*od*.promo\_so2,0)+                    COALESCE(*od*.promo\_so3,0)+COALESCE(*od*.promo\_so4,0)+                    COALESCE(*od*.promo\_so5,0)                ELSE                    COALESCE(*od*.promo\_po1,0)+COALESCE(*od*.promo\_po2,0)+                    COALESCE(*od*.promo\_po3,0)+COALESCE(*od*.promo\_po4,0)+                    COALESCE(*od*.promo\_po5,0)            END AS *promotion\_value*,            COALESCE(*od*.disc\_value\_final, *od*.disc\_value, *od*.disc\_value\_po,0) AS *discount\_value*,            COALESCE(*od*.vat\_value\_final, *od*.vat\_value, *od*.vat\_value\_po,0) AS *vat\_value*        FROM sls."order" *o*        JOIN mst.m\_outlet *mo*            ON *mo*.outlet\_id \= *o*.outlet\_id        LEFT JOIN (            SELECT emp\_id, emp\_code, emp\_name            FROM mst.m\_employee        ) *me*            ON *me*.emp\_id \= *o*.salesman\_id        JOIN sls.order\_detail *od*            ON *od*.ro\_no \= *o*.ro\_no        JOIN mst.m\_product *mp*            ON *mp*.pro\_id \= *od*.pro\_id        JOIN mst.m\_supplier *ms2*            ON *ms2*.sup\_id \= *mp*.sup\_id        WHERE *o*.ro\_date BETWEEN '2026-03-01' AND '2026-03-14'        AND *o*.salesman\_id IN (            62,204,206,209,210,211,226,228,229,232,234,235,236,240,242,244,            246,252,253,260,261,276,281,289,340,345,346,347,348,355,357,358,            359,360,362,370,382        )    ) *base* ) *calc*;  |
| :---- |

   

   

   

9. ### **SX-1241 | Cancel Sales Order**

[https://scyllax-pratesis.atlassian.net/browse/SX-1241](https://scyllax-pratesis.atlassian.net/browse/SX-1241) 

| curl 'https://best.scyllax.online/sales/v2/orders/SO2603130016' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsImVtcF9pZCI6MCwiZXhwaXJlcyI6MTc3MzU0MTA3OCwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4NTc0ODEyMzEyOCIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IkRpc3QgSURFIFNkYSIsInVzZXJfaWQiOjEyLCJ1c2VyX25hbWUiOiJkaXN0QHNkYS5pZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODU3NDgxMjMxMjgifQ.YYiGiPamPYqe7jeJs2NPxfc6Dvl82TXqhxR4G040sbA' \\   \-H 'Connection: keep-alive' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Not:A-Brand";v="99", "Google Chrome";v="145", "Chromium";v="145"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "Windows"'  |
| :---- |

**Issue** : setelah di cek, qty yang berpengaruh ke inv.stock dan inv.warehouse\_stock belum sesuai (dalam artian mengambil dari qty , sedangkan pada case ini sudah di tahap final order) artinya BE menggunakan qty1\_final, qty2\_final dan qty3\_final  untuk mengebalikan stock : 

- stock bertambah 1(L) 0(M) 2(S) → conversi 12   
- qty\_order berkurang 1(L) 0(M) 2(S) → conversi 12 

Query untuk mengambil qty order terakhir : 

| SELECT    *od*.pro\_id,    COALESCE(*od*.qty3\_final, *od*.qty3, *od*.qty\_po3) AS *qty3\_result*,    COALESCE(*od*.qty2\_final, *od*.qty2, *od*.qty\_po2) AS *qty2\_result*,    COALESCE(*od*.qty1\_final, *od*.qty1, *od*.qty\_po1) AS *qty1\_result* FROM sls."order" *o* JOIN sls.order\_detail *od*    ON *o*.ro\_no \= *od*.ro\_no WHERE *o*.ro\_no \= 'SO2603130014';  |
| :---- |

query diatas qty order diambil dari urutan prioritas : 

- qty\_final (prioritas 1\)  
- qty (prioritas 2 )   
- qty\_po (prioritas 3\)

* Insert ke inv.stock 


| cust\_id | stock\_id | stock\_date | tr\_code | tr\_no | wh\_id | pro\_id | item\_cdn | qty\_in | qty\_out | unit\_price | cogs | ref\_det\_id | created\_at | qty\_in\_order | qty\_out\_order |
| :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- | :---- |
| **C220010001** | 19572 | 2026-02-09 | SO | SO2603130014 | 304 | 739 | 1 | 12 | 0 | 6432.0000 | 0.0000 | 4949 | 1770647314 | 0 | 0 |
| **C220010001** | 19572 | 2026-02-09 | CO | SO2603130014-CO | 304 | 739 | 1 | 0 | 0 | 6432.0000 | 0.0000 | 4949 | 1770647314 | 0 | 12 |


* update ke inv.warehouse\_stock berdasarkan cust\_id, wh\_id, dan pro\_id  
- qty \= bertambah sebesar qty\_out final order (untuk case ini 12\)  
- qty\_on\_order \= berkurang sebesar qty\_out final order (untuk case ini 12\)




10. ### **SX-?? | Available Stock (Pop Up Edit Product)**

API 		: /inventory/v1/stocks/report  
Method 	: GET   
CURL 		: 

| curl 'https://best.scyllax.online/inventory/v1/stocks/report?page=1\&limit=10\&q=\&wh\_id=291\&pro\_id=715\&outlet\_id=1660\&order\_date=2026-03-24\&include\_zero\_stock=true\&active\_product\_only=true\&sort=pro\_code:asc' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsImVtcF9pZCI6MCwiZXhwaXJlcyI6MTc3NDQxMjU4MiwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4NTc0ODEyMzEyOCIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IkRpc3QgSURFIFNkYSIsInVzZXJfaWQiOjEyLCJ1c2VyX25hbWUiOiJkaXN0QHNkYS5pZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODU3NDgxMjMxMjgifQ.2lqm2xljCF6oyHbq8vQlV1nQxdMaqjxvid0zqPEH-DE' \\   \-H 'Connection: keep-alive' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"'  |
| :---- |

RESPONSE : 

| {     "message": "",     "data": \[         {             "pro\_id": 715,             "pro\_code": "AF-0002",             "pro\_name": "Action Figure Batman",             "unit\_id1": "PCS",             "unit\_id2": "CRT",             "unit\_id3": "CRT",             "purch\_price1": 0,             "purch\_price2": 0,             "purch\_price3": 0,             "sell\_price1": 0,             "sell\_price2": 0,             "sell\_price3": 0,             "total\_qty": 34,             "qty1": 8,             "qty2": 0,             "qty3": 6,             "total\_qty\_order": 12,             "qty\_order1": 2,             "qty\_order2": 0,             "qty\_order3": 1,             "total\_qty\_inc\_on\_order": 24,             "qty\_inc\_on\_order1": 6,             "qty\_inc\_on\_order2": 0,             "qty\_inc\_on\_order3": 4,             "conv\_unit2": 10,             "conv\_unit3": 1,             "is\_active": true,             "vat": 0,             "vat\_lg\_purch": 0,             "vat\_lg\_sell": 0         }     \],     "paging": {         "total\_record": 1,         "page\_current": 1,         "page\_limit": 10,         "page\_total": 1     },     "filter": {         "cust\_id": "C220010001",         "parent\_cust\_id": "C22001",         "date": "",         "order\_date": "2026-03-24",         "page": 1,         "limit": 10,         "q": "",         "sort": "pro\_code:asc",         "pro\_id": \[             715         \],         "wh\_id": \[             291         \],         "sup\_id": null,         "show\_price": "",         "include\_zero\_stock": "true",         "active\_product\_only": "true",         "brand\_id": null,         "pcat\_id": null,         "pl\_id": null,         "outlet\_id": 1660     },     "request\_id": "69c254dc7d31348d88ea8efa" }  |
| :---- |

#### 

#### 

#### **Enhance BE :** 

response qty1, qty2, qty3 → diambil dari stock warehouse 

Rumus warehouse stock :   
(dapat dilihat pada fitur eksisting **Inventory \> Distributor Stock**   
**![][image50]**

| curl 'https://best.scyllax.online/inventory/v1/stocks/report?limit=10\&page=1\&date=2026-03-24\&wh\_id=291\&pro\_id=715\&show\_conversion=true\&show\_price=true\&include\_zero\_stock=true\&active\_product\_only=true\&sort=pro\_code%3Aasc' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJkaXN0cmlidXRvcl9pZCI6NjcsImVtYWlsIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsImVtcF9pZCI6MCwiZXhwaXJlcyI6MTc3NDQxMjU4MiwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4NTc0ODEyMzEyOCIsInBhcmVudF9jdXN0X2lkIjoiQzIyMDAxIiwidXNlcl9mdWxsbmFtZSI6IkRpc3QgSURFIFNkYSIsInVzZXJfaWQiOjEyLCJ1c2VyX25hbWUiOiJkaXN0QHNkYS5pZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODU3NDgxMjMxMjgifQ.2lqm2xljCF6oyHbq8vQlV1nQxdMaqjxvid0zqPEH-DE' \\   \-H 'Connection: keep-alive' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"'  |
| :---- |
| **select** **SUM**(*s*.qty\_in) \- **SUM**(*s*.qty\_out) **as** *stock* **from** inv.stock *s*  **where** *s*.wh\_id \=291 **and** *s*.pro\_id \=715 result : 34 konversi menjadi L M S → 3 0 4 |

#### **Data Test :**

ro\_no \= SO2603240004  
pro\_id \= 715  
wh\_id \= 291 

 

#### **Handle by FE :**

FE ambil data dari : 

- Small 		\= qty1 \+ qty\_order1  
- Medium 	\= qty2 \+ qty\_order2  
- Large		\= qty3 \+ qty\_order3

### 

11. ###  **SX-1878 Issue Detail Order** 

[https://scyllax-pratesis.atlassian.net/browse/SX-1878](https://scyllax-pratesis.atlassian.net/browse/SX-1878)  

#### **Curl :** 

| curl 'https://best.scyllax.online/sales/v2/orders/SO2605060004' \\   \-H 'Accept: application/json, text/plain, \*/\*' \\   \-H 'Accept-Language: en-US,en;q=0.9' \\   \-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjowLCJkaXN0cmlidXRvcl9pZCI6MTAyLCJlbWFpbCI6ImFkbWluYm1AZ21haWwuY29tIiwiZW1wX2lkIjozODEsImV4cGlyZXMiOjE3NzgxMjU0MDUsImlzX2FkbWluIjpmYWxzZSwibGFuZ19pZCI6ImlkIiwibW9iaWxlX25vIjoiMDgxMzMzMzMzMzMzIiwicGFyZW50X2N1c3RfaWQiOiJDMjYwMDIiLCJ1c2VyX2Z1bGxuYW1lIjoiUGhpbGwgSm9uZXMiLCJ1c2VyX2lkIjoxNDEsInVzZXJfbmFtZSI6IlBoaWxsIEpvbmVzIiwid2hhdHNhcHAiOiIwODEzMzMzMzMzMzMifQ.7ocYl10aq5pFHtiidGH3dx\_uLwT6a5hCGCySs43xJjk' \\   \-H 'Connection: keep-alive' \\   \-H 'Origin: https://staging.scyllax.online' \\   \-H 'Referer: https://staging.scyllax.online/' \\   \-H 'Sec-Fetch-Dest: empty' \\   \-H 'Sec-Fetch-Mode: cors' \\   \-H 'Sec-Fetch-Site: same-site' \\   \-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10\_15\_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36' \\   \-H 'sec-ch-ua: "Google Chrome";v="147", "Not.A/Brand";v="8", "Chromium";v="147"' \\   \-H 'sec-ch-ua-mobile: ?0' \\   \-H 'sec-ch-ua-platform: "macOS"'  |
| :---- |

#### **Response :** 

| Tab Purchase Order purchase\_details.normal\[\].qty1\_stok \= 5 purchase\_details.normal\[\].qty2\_stok \= 0 purchase\_details.normal\[\].qty3\_stok \= 2 Tab Sales Order details.normal\[\].qty1\_stok \= 5 details.normal\[\].qty2\_stok \= 0 details.normal\[\].qty3\_stok \= 2 Tab Final Order details\_final.normal\[\].qty1\_stok \= 8 details\_final.normal\[\].qty2\_stok \= 0 details\_final.normal\[\].qty3\_stok \= 1  |
| :---- |

#### **Ekspektasi  :** 

| Tab Purchase Order purchase\_details.normal\[\].qty1\_stok \= 5 purchase\_details.normal\[\].qty2\_stok \= 0 purchase\_details.normal\[\].qty3\_stok \= 3 Tab Sales Order details.normal\[\].qty1\_stok \= 5 details.normal\[\].qty2\_stok \= 0 details.normal\[\].qty3\_stok \= 3 Tab Final Order details\_final.normal\[\].qty1\_stok \= 5 details\_final.normal\[\].qty2\_stok \= 0 details\_final.normal\[\].qty3\_stok \= 3 *nilai 8 0 1 ini secara konversi ke small , ini sudah benar  yaitu 13  namun secara konversi (large, medium, small) seharusnya adalah 2 0 3*   |
| :---- |

#### **Hasil Analisa Eksisting  :** 

| Hasil analisa dari repo sales terbaru: \- Endpoint GET /sales/v2/orders/:ro\_no masuk ke OrderController.DetailV2, lalu OrderService.DetailV2. \- Data detail order dibaca dari DB:   \- Header: sls.order   \- Detail: sls.order\_detail   \- Stock warehouse terbaru: model.WarehouseStock lewat FindWarehouseStockByWhIdAndProIds() Jawaban “response ini diambil dari rumus apa?” Untuk qty\*\_stok, rumusnya di service/order\_service.go. Tab Sales order / details.normal\[\].qty\*\_stok Flow: \- Ambil warehouseStockMap\[pro\_id\] \= stock warehouse terbaru dalam small qty. \- Convert ke qty1/qty2/qty3. \- Kalau order tidak cancelled, ditambah qty order tab Sales Order. Rumusnya: stockConverted \= ConvToQtyConversion(warehouseStockQty) qty1\_stok \= stockConverted.Qty1 \+ details.qty1 qty2\_stok \= stockConverted.Qty2 \+ details.qty2 qty3\_stok \= stockConverted.Qty3 \+ details.qty3 Lokasi: \- order\_service.go:2832-2857 Tab Purchase order / purchase\_details.normal\[\].qty\*\_stok purchase\_details awalnya cuma copy dari details. response.PurchaseDetails.Normal \= copy(response.Details.Normal) Jadi qty\*\_stok Purchase Order mengikuti hasil Sales Order, bukan dihitung sendiri dari PO qty. Lokasi: \- order\_service.go:2917-2921 Tab Final order / details\_final.normal\[\].qty\*\_stok Flow hampir sama, tapi yang ditambahkan adalah qty final hasil konversi dari qty\_final. Rumusnya: stockConvertedFinal \= ConvToQtyConversion(warehouseStockQty) qty1\_stok \= stockConvertedFinal.Qty1 \+ details\_final.qty1 qty2\_stok \= stockConvertedFinal.Qty2 \+ details\_final.qty2 qty3\_stok \= stockConvertedFinal.Qty3 \+ details\_final.qty3 Lokasi: \- order\_service.go:3042-3067 Kenapa hasilnya beda: PO/SO \= 5 0 2, Final \= 8 0 1 Karena kode memang menghitung basis qty yang berbeda: \- details / Sales Order pakai qty \- details\_final / Final Order pakai qty\_final \- purchase\_details hanya copy dari details Jadi ekspektasi bahwa semua tab mengikuti final order belum sesuai implementasi saat ini. Issue 2: 8 0 1 total small \= 13 tapi harusnya jadi 2 0 3 Benar, dari kode conversion sekarang orientasinya: qty3 \= total / (conv\_unit2 \* conv\_unit3) qty2 \= sisa / conv\_unit2 qty1 \= sisa  Lokasi: \- pkg/conversion/quantity.go:15-23 Artinya output struct saat ini: \- qty1 \= small \- qty2 \= medium \- qty3 \= large Kalau total small \= 13 dan konversi misalnya conv\_unit2=5, conv\_unit3=2, maka kode akan menghasilkan: \- large \= 1 \- medium \= 0 \- small \= 3 \- tetapi disimpan ke response sebagai qty1=3, qty2=0, qty3=1 Kalau UI/domain menganggap urutan field adalah (large, medium, small), maka harus dibalik mapping-nya menjadi: \- qty1 \= large \- qty2 \= medium \- qty3 \= small Jadi root cause Issue 2: hasil konversi benar secara total small, tapi mapping field qty1/qty2/qty3 di response tidak sesuai ekspektasi urutan large-medium-small.  |
| :---- |
| **Rumus Conversi :**  qty3 \= total / (conv\_unit2 \* conv\_unit3) qty2 \= sisa / conv\_unit2 qty1 \= sisa  |

#### **New Analysis :** 

| Kenapa hasilnya beda: PO/SO \= 5 0 2, Final \= 8 0 1 Karena kode memang menghitung basis qty yang berbeda: \- details / Sales Order pakai qty \- details\_final / Final Order pakai qty\_final \- purchase\_details hanya copy dari details Jadi ekspektasi bahwa semua tab mengikuti final order belum sesuai implementasi saat ini. Issue 2: 8 0 1 total small \= 13 tapi harusnya jadi 2 0 3 Benar, dari kode conversion sekarang orientasinya: qty3 \= total / (conv\_unit2 \* conv\_unit3) qty2 \= sisa / conv\_unit2 qty1 \= sisa  Flow Sebelumnya :  \- Ambil warehouseStockMap\[pro\_id\] \= stock warehouse terbaru dalam small qty. \- Convert ke qty1/qty2/qty3. \- Kalau order tidak cancelled, ditambah qty order tab Sales Order. conv\_unit2 \= 5  conv\_unit3 \= 1 wh\_stock          \=  4     | 0 0 4 qty\_order         \=   9    | 1 0 4 available stock \=  13   | 1 0 8                                         2 0 3  Flow yang harus diperbaiki :  \- WAREHOUSE\_STOCK \= Ambil warehouseStockMap\[pro\_id\] \= stock warehouse terbaru dalam small qty. \- QTY\_ORDER \= Kalau order tidak cancelled, Ambil QTY\_ORDER= convert ke small qty  \- Available Stock : WAREHOUSE\_STOCK \+ QTY\_ORDER  \- Convert Available Stock ke satuan Large , Medium, Small (masih menggunakan kode conversion yang eksisting)  QTY\_ORDER  \= untuk qty order , perhatikan untuk masing-masing kebutuhan tab  purchase\_order : qty\_po1, qty\_po2, qty\_po3 (prioritas 3\)  sales\_order : qty1, qty2, qty3 (prioritas 2\)  final\_order : qty1\_final, qty2\_final, qty3\_final  (prioritas 1\) jadi, jika qty\*\_final sudah terisi , maka qty\_order ambil dr qty\*\_final  Data Skenario :  ro\_no   \= SO2605060004 cust\_id \= C260020001 pro\_id \= 10778 Purchase Order :  qty\_po1  \= 2  (Large)  qty\_po2  \= 0 (Medium)  qty\_po3  \=  1   (Small)  Sales Order Tab :  qty1 \= 1 (Large) qty2 \= 0  (Medium) qty3 \= 2  (Small)  Final Order Tab : qty1\_final \=  4  (Large) qty2\_final \=  0  (Medium) qty3\_final \=  1  (Small)  warehouse\_stock pada wh\_id \=  350 , dan pro\_id \= 10778 adalah 4 (4 ini adalah satuan small) conversion pro\_id 10778 conv\_unit2 \= 5 conv\_unit3 \= 3  Implementasi Rumus Available Stock :  WAREHOUSE\_STOCK    \= 4 (Small)   QTY\_ORDER                   \= karena qty\*final sudah terisi, maka ambil pada qty\*\_final (semua mengikuti qty order terbaru) 1 0 4 \= 9 Small  Available Stock \= 13 Small → convert large : 2 Medium : 0 Small : 3  Jadi response pada detail order :  (ini di pakai untuk available stock tab Purchase Order) purchase\_details\[\].normal.qty1\_stock (Large)     \= 2  purchase\_details\[\].normal.qty2\_stock (Medium) \= 0  purchase\_details\[\].normal.qty3\_stock (Small)     \= 3 (ini di pakai untuk available stock tab Sales Order) details\[\].normal.qty1\_stock (Large)     \= 2  details\[\].normal.qty2\_stock (Medium) \= 0  details\[\].normal.qty3\_stock (Small)     \= 3 (ini di pakai untuk available stock tab Final Order) details\_final\[\].normal.qty1\_stock (Large)     \= 2  details\_final\[\].normal.qty2\_stock (Medium) \= 0  details\_final\[\].normal.qty3\_stock (Small)     \= 3 |
| :---- |
