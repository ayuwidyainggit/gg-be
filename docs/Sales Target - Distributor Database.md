1. ### **mst.m\_sales\_target\_distributor\_yearly**

   **Nama Database:** m\_sales\_target\_distributor\_yearly

   **Schema** : mst

   **Tanggal Pembuatan: (**27 Desember 2025\)

   **Note for developer** : **create new table** 

   **Description :** digunakan untuk menyimpan data sales target  

#### **Kolom (*Columns / Fields*):**

| Nama Kolom | Tipe Data | Nullable | Nilai Default | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| cust\_id\*\* | varchar(10) | No  | \- | ID customer Ref : ***smc.m\_customer***   |
| sales\_target\_distributor\_yearly\_id\* | int(8) | No | \- | ***Primary Key*** |
| area\_id\*\* | int(8) | No  | \- | ***Foreign Key*** mst.m\_area |
| region\_id\*\* | int(8) | No  | \- | ***Foreign Key*** mst.m\_region |
| distributor\_id\*\* | int(8) | No | \- | ***Foreign Key*** mst.m\_distributor |
| year | int(4) | No | \- | target tahun |
| yearly\_target | int(11) | No | \- | jumlah target per tahun yang telah di tentukan  |
| status | int(8) | No | \- | 0 : Draft  1 : Aktif  2 : Nonaktif  **perubahan status** |
| is\_active | bool | Yes | \- | True : Active False : Inactive **delete field** |
| user\_inactive | int(8) | Yes | \- | ID pengguna yang melakukan inactive data Ref : ***sys.m\_user*** |
| inactive\_at | timestamptz(6) | No | \- | Waktu data diinactivekan. |
| created\_by | int(4) | No | \- | ID pengguna yang membuat data. Ref : ***sys.m\_user*** |
| created\_at | timestamptz(6) | No | \- | Waktu data dibuat. |
| updated\_by | int(8) | Yes | \- | ID pengguna yang terakhir mengubah. Ref : ***sys.m\_user*** |
| updated\_at | timestamptz(6) | Yes | \- | Waktu data terakhir diubah. |
| deleted\_by | int(4) | Yes | \- | ID pengguna yang menghapus. Ref : ***sys.m\_user*** |
| deleted\_at | timestamptz(6) | Yes | \- | Waktu data dihapus. |
| is\_del | bool | Yes | \- | Status soft delete. |

## 

2. ### **mst.m\_sales\_target\_distributor\_monthly**

   **Nama Database:** m\_sales\_target\_distributor\_monthly

   **Schema** : mst

   **Tanggal Pembuatan: (**27 Desember 2025\)

   **Note for developer** : **create new table** 

   **Description :** digunakan untuk menyimpan data sales target  

#### **Kolom (*Columns / Fields*):**

| Nama Kolom | Tipe Data | Nullable | Nilai Default | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| cust\_id\*\* | varchar(10) | No  | \- | ID customer Ref : ***smc.m\_customer***   |
| sales\_target\_distributor\_monthly\_id\* | int(8) | No | \- | ***Primary Key*** |
| sales\_target\_distributor\_yearly\_id\*\* | int(8) | No | \- | ***Foreign Key*** mst.m\_sales\_target\_distributor\_monthly |
| month | int(2) | No | \- | bulan  |
| monthly\_target | int(11) | No | \- | target perbulan |
| is\_active | bool | Yes |  |  |
| created\_by | int(4) | No | \- | ID pengguna yang membuat data. Ref : ***sys.m\_user*** |
| created\_at | timestamptz(6) | No | \- | Waktu data dibuat. |
| updated\_by | int(8) | Yes | \- | ID pengguna yang terakhir mengubah. Ref : ***sys.m\_user*** |
| updated\_at | timestamptz(6) | Yes | \- | Waktu data terakhir diubah. |
| deleted\_by | int(4) | Yes | \- | ID pengguna yang menghapus. Ref : ***sys.m\_user*** |
| deleted\_at | timestamptz(6) | Yes | \- | Waktu data dihapus. |
| is\_del | bool | Yes | \- | Status soft delete. |

## 

