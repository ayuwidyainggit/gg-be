1. ### **mst.outlet\_cr**

**Deskripsi:** Tabel ini adalah tabel master yang menyimpan semua perubahan mengenai outlet yang dibedakan dari 2 sumber (source : mobile / web) serta approval dari setiap perubahan 

**Kolom (*Columns / Fields*):**

| Nama Kolom | Tipe Data | *Nullable* | Nilai *Default* | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| cust\_id | int(8) |  | \- | id customer |
| outlet\_cr\_id | int(8) |  | \- | Primary Key  |
| outlet\_id | int(8) |  | nextval('mst.outlet\_id\_seq'::regclass) | FK : mst.m\_outlet |
| source | int(8) |  |  | 1: web2: mobile |
| created\_by | int(8) |  |  |  |
| created\_at | timestamp(6) |  |  |  |
| status | int(3) |  |  | 1: pending2: approve3: reject |
| approval\_by | int(8) |  |  | diisi user yang  |
| aproveal\_at | timestamp |  | \- | Status outlet. |

2. ### **mst.outlet\_cr\_det**

**Deskripsi:** Tabel ini adalah tabel untuk menyimpan field yang berubah. 1x perubahan untuk merubah lebih dari 1 field 

**Kolom (*Columns / Fields*):**

| Nama Kolom | Tipe Data | *Nullable* | Nilai *Default* | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| outlet\_cr\_det\_id | varchar(10) |  | \- | Primary Key  |
| outlet\_cr\_id | int(8) |  | nextval('mst.outlet\_id\_seq'::regclass) | FK.mst.outlet\_cr |
| field\_name | varchar(30) |  | \- | diisi field name dari mst.m\_outlet yang dilakukan perubahan  |
| new\_value | varchar(225) |  | \- | diisi value perubahan  |
| old\_value | varchar(225) |  | \- | diisi value sebelum perubahan  |

