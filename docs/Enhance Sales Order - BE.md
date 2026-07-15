1. # **Taking Order** 

URL: [http://103.28.219.73:5001/mobile/v1/orders](http://103.28.219.73:5001/mobile/v1/orders)

| \=== Request Body \=== {   "ro\_date": "2026-01-22",   "opr\_type": "O", // untuk FE dihandle mas nouval   "due\_date": "2026-01-22",   "salesman\_id": 228,   "wh\_id": 291,   "outlet\_id": 234,   "delivery\_date": "2026-01-22",   "order\_no": "PO26012230122410001",   "pay\_type": 0,   "sub\_total": 1389099,   "disc": 0,   "disc\_value": 0,   "promo\_value": 0,   "cash\_disc\_value": 0.0,   "tot\_disc1": 0,   "tot\_disc2": 0,   "vat": 85842,   "vat\_value": 171685,   "total": 1560784,   "data\_status": 1,   "created\_by": 228,   "data\_source": 2,   "details": {     "normal": \[       {         "pro\_id": 472,         "item\_type": 1,         "disc\_value": 0,         "vat": 11,         "vat\_value": 111246,         "vat\_bg": 0,         "vat\_bg\_value": 0,         "qty1\_stok": 0,         "qty2\_stok": 0,         "qty3\_stok": 6,         "qty1": 0,         "qty1\_final": 0,         "sell\_price1": 40128,         "unit\_id1": "KLG",         "conv\_unit1": 0,         "qty2": 2,         "qty2\_final": 2,         "sell\_price2": 505667,         "unit\_id2": "KRT",         "conv\_unit2": 12,         "qty3": 0,         "sell\_price3": 505667,         "unit\_id3": "KRT",         "conv\_unit3": 12,         "amount": 1122580       },       {         "pro\_id": 475,         "item\_type": 1,         "disc\_value": 0,         "vat": 11,         "vat\_value": 60439,         "vat\_bg": 0,         "vat\_bg\_value": 0,         "qty1\_stok": 0,         "qty2\_stok": 0,         "qty3\_stok": 15,         "qty1": 0,         "qty1\_final": 0,         "sell\_price1": 3053,         "unit\_id1": "PCS",         "conv\_unit1": 0,         "qty2": 3,         "qty2\_final": 3,         "sell\_price2": 183150,         "unit\_id2": "KRT",         "conv\_unit2": 60,         "qty3": 0,         "sell\_price3": 183150,         "unit\_id3": "KRT",         "conv\_unit3": 60,         "amount": 609889       }     \]   }  |
| :---- |

Enhance Request :

1) tambahkan opr\_type \= char (1) untuk membedakan data order  Canvas / Taking Order  
- O \= Taking Order  
- C \= Canvas


Enhance Response : 

1) sls.oder (tdk ada)  
2) sls.oder\_detail

| Field | Data Source |
| :---- | :---- |
| sell\_price\_system1 | sell\_price1 |
| sell\_price\_system2 | sell\_price2 |
| sell\_price\_system3 | sell\_price3 |
| qty\_po1 | qty1 |
| qty\_po2 | qty2 |
| qty\_po3 | qty3 |
| sell\_price\_po1 | sell\_price1 |
| sell\_price\_po2 | sell\_price2 |
| sell\_price\_po3 | sell\_price3 |
| disc\_po | disc |
| disc\_value\_po | disc\_value |
| vat\_po | vat |
| vat\_value\_po | vat\_value |

   

2. # **Canvas**

URL: [http://103.28.219.73:5001/mobile/v1/orders](http://103.28.219.73:5001/mobile/v1/orders)

| \=== Request Body \=== {   "ro\_date": "2026-01-22",   "opr\_type": "C", // untuk FE dihandle mas firman   "due\_date": "2026-01-22",   "salesman\_id": 228,   "wh\_id": 248,   "outlet\_id": 234,   "delivery\_date": "2026-01-22",   "order\_no": "PO26012230122410002",   "pay\_type": 0,   "sub\_total": 3014256,   "disc": 0,   "disc\_value": 0,   "promo\_value": 0,   "cash\_disc\_value": 0.0,   "tot\_disc1": 0,   "tot\_disc2": 0,   "vat": 186274,   "vat\_value": 372548,   "total": 3386804,   "data\_status": 1,   "created\_by": 228,   "data\_source": 2,   "details": {     "normal": \[       {         "pro\_id": 482,         "item\_type": 1,         "disc\_value": 0,         "vat": 11,         "vat\_value": 20900,         "vat\_bg": 0,         "vat\_bg\_value": 0,         "qty1\_stok": 0,         "qty2\_stok": 0,         "qty3\_stok": 49,         "qty1": 1,         "qty1\_final": 1,         "sell\_price1": 4634,         "unit\_id1": "PCS",         "conv\_unit1": 0,         "qty2": 2,         "qty2\_final": 2,         "sell\_price2": 92685,         "unit\_id2": "KRT",         "conv\_unit2": 20,         "qty3": 0,         "sell\_price3": 92685,         "unit\_id3": "KRT",         "conv\_unit3": 20,         "amount": 210904       },       {         "pro\_id": 484,         "item\_type": 1,         "disc\_value": 0,         "vat": 11,         "vat\_value": 351648,         "vat\_bg": 0,         "vat\_bg\_value": 0,         "qty1\_stok": 0,         "qty2\_stok": 0,         "qty3\_stok": 14,         "qty1": 0,         "qty1\_final": 0,         "sell\_price1": 3330,         "unit\_id1": "PCS",         "conv\_unit1": 0,         "qty2": 0,         "qty2\_final": 0,         "sell\_price2": 79920,         "unit\_id2": "PCK",         "conv\_unit2": 24,         "qty3": 2,         "qty3\_final": 2,         "sell\_price3": 1598400,         "unit\_id3": "KRT",         "conv\_unit3": 20,         "amount": 3548448       }     \]   } } \=== Response \=== {   "message": "Berhasil Dibuat",   "request\_id": "6971fd75396bbc0ee9a9689c" }  |
| :---- |

Enhance : 

3) sls.oder

| Field | Data Source |
| :---- | :---- |
| invoice\_no | khusus untuk opr\_type \= C BE generate invoice  ***(tolong handle supaya saat hit bersamaan invoice\_no tetep unik)***   INVYYMMDD-4 digit running numberINV2501150010 |
| invoice\_date | diisi tanggal hit api  |
| validate\_stok | true |
| validate\_stok\_message | Sufficient Stock Note : qty order gaboleh lebih besar dibanding stock gudang (gudang yang dipilih) qty1 qty2 qty3 tidak boleh melebihi qty1\_stok qty2\_stok qty3\_stok  |

4) sls.oder\_detail

| Field | Data Source |
| :---- | :---- |
| sell\_price\_system1 | sell\_price1 |
| sell\_price\_system2 | sell\_price2 |
| sell\_price\_system3 | sell\_price3 |
| qty\_po1 | qty1 |
| qty\_po2 | qty2 |
| qty\_po3 | qty3 |
| sell\_price\_po1 | sell\_price1 |
| sell\_price\_po2 | sell\_price2 |
| sell\_price\_po3 | sell\_price3 |
| disc\_po | disc |
| dis\_value\_po | disc\_value |
| vat\_po | vat |
| vat\_value\_po | vat\_value |

