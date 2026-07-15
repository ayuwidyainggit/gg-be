1. ### **mst.m\_survey\_template**

**Deskripsi:** Tabel ini adalah tabel master yang menyimpan semua data survey  
**Kolom (*Columns / Fields*):**

| Nama Kolom | Tipe Data | *Nullable* | Nilai *Default* | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| cust\_id | int(8) |  | \- | id customer |
| survey\_template\_id\* | int(8) |  | nextval('mst.survey\_template\_id\_seq'::regclass) | Primary Key   |
| template\_code | varchar(10) |  |  | Kode template |
| template\_title | varchar(150) |  |  | Judul Template |
| question\_total | int(8) |  |  | Jumlah question per template |
| use\_image | bool |  |  | true : wajib input photo  false : tidak wajib input photo |
| is\_active | bool |  | true | Status yang menandakan apakah data produk ini aktif. |
| created\_by | int(8) |  | \- | Pengguna yang menambahkan data ini. |
| created\_at | timestamptz(6) |  | \- | Waktu ketika data ini dibuat. |
| updated\_by | int(8) |  | \- | Pengguna yang melakukan modifikasi pada data. |
| updated\_at | timestamptz(6) |  | \- | Waktu ketika modifikasi pada data dilakukan. |
| is\_del | bool |  | false | Status untuk menandai penghapusan data. |
| deleted\_by | int(8) |  | \- | Pengguna yang melakukan penghapusan data. |
| deleted\_at | timestamptz(6) |  | \- | Waktu ketika penghapusan pada data dilakukan. |

2. ### **mst.question\_template**

**Deskripsi:** Tabel ini adalah tabel master yang menyimpan data question dari masing-masing data survey  
**Kolom (*Columns / Fields*):**

| Nama Kolom | Tipe Data | *Nullable* | Nilai *Default* | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| cust\_id | int(8) |  | \- | id customer |
| question\_template\_id\* | int(8) |  | nextval('mst.question\_template\_id\_seq'::regclass) | Primary Key  |
| survey\_template\_id\*\* | int(8) |  | \- | Foreign Key  |
| question | varchar(225) |  |  | question list |
| input\_type | enum |  |  | textfield, dropdown, radiobutton, toggle, checkbox, |
| answer\_type | enum |  |  | Single, Multiple, Free Text |
| seq | int 4 |  |  | urutan  |

3. ### **mst.m\_q\_option\_template**

**Deskripsi:** Tabel ini adalah tabel master yang menyimpan data question dari masing-masing data survey  
**Kolom (*Columns / Fields*):**

| Nama Kolom | Tipe Data | *Nullable* | Nilai *Default* | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| cust\_id | int(8) |  | \- | id customer |
| q\_option\_template\_id\* | int(8) |  | nextval('mst.q\_option\_template\_id\_seq'::regclass) | Primary Key  |
| question\_template\_id\*\* | int(8) |  | \- | Foreign Key  |
| option | varchar(225) |  |  | option list  |

