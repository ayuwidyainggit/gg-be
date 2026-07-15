![][image1]

[Link](https://dbdiagram.io/d/mobile-expense-6925c8a37d9416ddff0ec96c) ERD 

1. ### **acf.expense\_type**

   **Nama Database:** expense\_type 

   **Schema** : acf

   **Tanggal Pembuatan: (**26 November 2025\)

   **Note for developer** : **create new table** 

   **Description :** digunakan untuk menyimpan data expense 

#### **Kolom (*Columns / Fields*):**

| Nama Kolom | Tipe Data | Nullable | Nilai Default | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| cust\_id\*\* | varchar(10) | No  | \- | ID customer Ref : ***smc.m\_customer***   |
| expense\_type\_id\*\* | int(8) | No | \- | ***PK*** ID untuk setiap entri. |
| expense\_type\_code | varchar(20) | No | \- | code expense type  |
| expense\_type\_name | varchar(50) | No | \- | expense name  |
| note | varchar(100) | Yes | \- | Catatan tambahan. |
| source | int | No |  | 1: web 2 : mobile |
| is\_active | bool | Yes |  | True : Active False : Inactive |
| created\_by | int(4) | No | \- | ID pengguna yang membuat data. Ref : ***sys.m\_user*** |
| created\_at | timestamptz(6) | No | \- | Waktu data dibuat. |
| updated\_by | int(8) | Yes | \- | ID pengguna yang terakhir mengubah. Ref : ***sys.m\_user*** |
| updated\_at | timestamptz(6) | Yes | \- | Waktu data terakhir diubah. |
| deleted\_by | int(4) | Yes | \- | ID pengguna yang menghapus. Ref : ***sys.m\_user*** |
| deleted\_at | timestamptz(6) | Yes | \- | Waktu data dihapus. |
| is\_del | bool | Yes | \- | Status soft delete. |

## 

**Contoh Data :** 

1. E001 \- Parking Fee  
2. E002 \- Gasoline  
3. E003 \- Meal Allowance 

   

   2. ### **acf.expense**

* **Nama Database:** expense 

* **Schema** : acf

* **Tanggal Pembuatan: (**26 November 2025\)

* **Note for developer** : **create new table** 

* **Description :** digunakan untuk menyimpan data expense 

#### **Kolom (*Columns / Fields*):**

| Nama Kolom | Tipe Data | Nullable | Nilai Default | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| cust\_id\*\* | varchar(10) | No  | \- | ID customer Ref : ***smc.m\_customer***   |
| expense\_id\* | int(8) | No  | \- | ***PK*** ID untuk setiap entri. |
| expense\_type\_id\*\* | int(8) | No | \- | expense type id Ref : ***acf.expense\_type*** |
| doc\_no | varchar(50) | No |  | EYYYYMMDD-3digitRunnningNumber  E20261201222 |
| source | int | No |  | 1: web 2 / null: mobile |
| date | date | No  | \- | Tanggal entry data  |
| amount | int(11) | No | \- | SUM (amount dari acf.expense\_det) |
| balance | int(11) | No | \- | sisa saldo  |
| created\_by\*\* | int(4) | No | \- | ID pengguna yang membuat data. Ref : ***sys.m\_user*** |
| note | varchar | 100 | \- |  |
| created\_at | timestamptz(6) | No | \- | Waktu data dibuat. |
| updated\_by\*\* | int(8) | Yes | \- | ID pengguna yang terakhir mengubah. Ref : ***sys.m\_user*** |
| updated\_at | timestamptz(6) | Yes | \- | Waktu data terakhir diubah. |
| deleted\_by\*\* | int(4) | Yes | \- | ID pengguna yang menghapus. Ref : ***sys.m\_user*** |
| deleted\_at | timestamptz(6) | Yes | \- | Waktu data dihapus. |
| is\_del | bool | Yes | \- | Status soft delete. |

## 

3. ### **acf.expense\_det (tdk dipakai)**

* **Nama Database:** expense 

* **Schema** : acf

* **Tanggal Pembuatan: (**26 November 2025\)

* **Note for developer** : **create new table** 

* **Description :** digunakan untuk menyimpan data outlet dari expense yang di input   
  pada figma, 1 expense terdiri dari lebih dari 1 outlet  

#### **Kolom (*Columns / Fields*):**

| Nama Kolom | Tipe Data | Nullable | Nilai Default | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| cust\_id\*\* | varchar(10) | No  | \- | ID customer Ref : ***smc.m\_customer***   |
| expense\_det\_id | int(8) | No  |  | ***PK*** ID untuk setiap entri. |
| expense\_id\* | int(8) | No  | \- | Ref : ***acf.expense*** |
| collector\_id | int(8) | No  | \- | Ref : ***sys.m\_user*** |
| expense\_type\_id\*\* | int(8) | No | \- | expense type id Ref : ***acf.expense\_type*** |
| amount | int(11) | No | \- | SUM (amount dari acf.expense\_det) |
| notes | varchar(100) | No | \- | Catatan tambahan. |

4. ### **acf.expense\_file**

* **Nama Database:** expense\_file 

* **Schema** : acf

* **Tanggal Pembuatan: (**26 November 2025\)

* **Note for developer** : **create new table** 

* **Description :** digunakan untuk menyimpan data file dari expense yang di input   
  pada figma, 1 expense terdiri maksimal 3 file  (disediakan tabel baru apabila  ada enhance untuk file lebih dari 3\)  

#### **Kolom (*Columns / Fields*):**

| Nama Kolom | Tipe Data | Nullable | Nilai Default | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| cust\_id\*\* | varchar(10) | No  | \- | ID customer Ref : ***smc.m\_customer***   |
| expense\_file\_id | int(8) | No  |  | ***PK*** ID untuk setiap entri. |
| expense\_det\_id | int(8) | No  |  | *Ref :*  acf.expense\_det.expense\_det\_id |
| file\_name | varchar(255) | No | \- | Nama File  |
| file\_url | varchar(50) | No | JPG | photo disimpan dalam format JPG |
| file\_key | ENUM('image','video') | No | \- | Membedakan jenis konten: *image* atau *video* |
| media\_category | text | No | \- | url file |
| file\_size | BIGINT | No | \- | Untuk mengetahui ukuran file asli |
