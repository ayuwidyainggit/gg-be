

1. ### **mst.m\_sales\_target**

   **Nama Database:** m\_sales\_target 

   **Schema** : mst

   **Tanggal Pembuatan: (**27 Desember 2025\)

   **Note for developer** : **create new table** 

   **Description :** digunakan untuk menyimpan data sales target  salesman

#### **Kolom (*Columns / Fields*):**

| Nama Kolom | Tipe Data | Nullable | Nilai Default | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| cust\_id\*\* | varchar(10) | No  | \- | ID customer Ref : ***smc.m\_customer***   |
| sales\_target\_id\* | int(8) | No | \- | ***Primary Key*** |
| sales\_target\_distributor\_yearly\_id\*\* | int(8) | No | \- | ***Foreign Key*** Ref :  mst.m\_sales\_target\_distributor\_yearly |
| sales\_target\_distributor\_monthly\_id\*\* | int(8) | No | \- | ***Foreign Key*** Ref :  mst.m\_sales\_target\_distributor\_monthly |
| month | int(2) | No | \- | target bulan  |
| year | int(4) | No | \- | target tahun |
| allocated\_total | int(11) | No |  |  |
| monthly\_target | int(11) | No | \- | jumlah target per bulan dan tahun yang telah di tentukan  |
| remaining | int(11) | No |  |  |
| status | int(11) | No |  | 0 : Draft  1 : Aktif  2 : Nonaktif  |
| created\_by | int(4) | No | \- | ID pengguna yang membuat data. Ref : ***sys.m\_user*** |
| created\_at | timestamptz(6) | No | \- | Waktu data dibuat. |
| updated\_by | int(8) | Yes | \- | ID pengguna yang terakhir mengubah. Ref : ***sys.m\_user*** |
| updated\_at | timestamptz(6) | Yes | \- | Waktu data terakhir diubah. |
| deleted\_by | int(4) | Yes | \- | ID pengguna yang menghapus. Ref : ***sys.m\_user*** |
| deleted\_at | timestamptz(6) | Yes | \- | Waktu data dihapus. |
| is\_del | bool | Yes | \- | Status soft delete. |

## 

2. ### **mst.m\_sales\_allocated**

   **Nama Database:** m\_sales\_allocated

   **Schema** : mst

   **Tanggal Pembuatan: (**27 Desember 2025\)

   **Note for developer** : **create new table** 

   **Description :** digunakan untuk menyimpan data sales target  salesman

#### **Kolom (*Columns / Fields*):**

| Nama Kolom | Tipe Data | Nullable | Nilai Default | Deskripsi |
| :---- | :---- | :---- | :---- | :---- |
| cust\_id\*\* | varchar(10) | No  | \- | ID customer Ref : ***smc.m\_customer***   |
| sales\_allocated\_id\* | int(8) | No | \- | ***Primary Key*** |
| sales\_target\_id\*\* | int(8) | No | \- | ***Foreign Key** mst.m\_sales\_target* |
| salesman\_id\*\* | int(8) | No | \- | ***Foreign Key** mst.m\_salesman (emp\_id)* |
| sales\_team\_id\*\* | int(8) | Yes | \- | ***Foreign Key** mst.m\_sales\_team* |
| allocated | int(11) | No | \- | target per salesman |
| is\_active | bool | Yes |  | True : Active False : Inactive |
| created\_by | int(4) | No | \- | ID pengguna yang membuat data. Ref : ***sys.m\_user*** |
| created\_at | timestamptz(6) | No | \- | Waktu data dibuat. |
| updated\_by | int(8) | Yes | \- | ID pengguna yang terakhir mengubah. Ref : ***sys.m\_user*** |
| updated\_at | timestamptz(6) | Yes | \- | Waktu data terakhir diubah. |
| deleted\_by | int(4) | Yes | \- | ID pengguna yang menghapus. Ref : ***sys.m\_user*** |
| deleted\_at | timestamptz(6) | Yes | \- | Waktu data dihapus. |
| is\_del | bool | Yes | \- | Status soft delete. |

## 

