# Extracted from Create_Survey_Database.docx

mst.m_survey
Deskripsi: Tabel ini adalah tabel master yang menyimpan semua data surveyKolom (Columns / Fields):
Changes :
25 Mei 2026 = Perubahan enum answer_frequency Before : One Time, MultipleAfter :
One Time
Multiple Times, One Day
Multiple Times, Different Day
Nama Kolom
Tipe Data
Nullable
Nilai Default
Deskripsi
cust_id
int(8)
-
id customer
survey_id*
int(8)
nextval('mst.survey_id_seq'::regclass)
Primary Key
survey_title
varchar(150)
answer_frequency
enum
One Time
Multiple Times, One Day
Multiple Times, Different Day
response_type
enum
Mandatory, Optional
target_type
enum
National
Spesific
Area
level_target
enum
Outlet
Distributor
Salesman
efective_date_start
date
efective_date_end
date
emp_id**
INT8[]
Foreign Key
mst.m_salesman
is_active
bool
true
Status yang menandakan apakah data produk ini aktif.
true : active
false : nonactive
created_by
int(8)
-
Pengguna yang menambahkan data ini.
created_at
timestamptz(6)
-
Waktu ketika data ini dibuat.
updated_by
int(8)
-
Pengguna yang melakukan modifikasi pada data.
updated_at
timestamptz(6)
-
Waktu ketika modifikasi pada data dilakukan.
is_del
bool
false
Status untuk menandai penghapusan data.
deleted_by
int(8)
-
Pengguna yang melakukan penghapusan data.
deleted_at
timestamptz(6)
-
Waktu ketika penghapusan pada data dilakukan.
mst.m_survey_area
Deskripsi: Tabel ini adalah tabel master yang menyimpan semua data survey
Enhance :
3 Juli 2026 : Tambahan target_cust_id
Kolom (Columns / Fields):
Nama Kolom
Tipe Data
Nullable
Nilai Default
Deskripsi
cust_id
int(8)
-
id customer
m_survey_area_id*
int(8)
nextval('mst.m_survey_detail_id_seq'::regclass)
Primary Key
survey_id**
int(8)
-
Foreign Key
mst.m_survey
area_id
int(8)
-
Foreign Key
mst.m_area (ini berdasarkan data mst.m_distributor)
distributor_id
int(8)
dari business_unit multiple dropdown
target_cust_id
varchar(10)
cust_id dari business unit yang dipilih
impact di mobile adalah, cust_id tersebut memiliki task untuk mengisi survey
mst.m_survey_outlet
Deskripsi: Tabel ini adalah tabel master yang menyimpan semua data survey untuk outlet yang dikunjungiKolom (Columns / Fields):
Nama Kolom
Tipe Data
Nullable
Nilai Default
Deskripsi
cust_id
int(8)
-
id customer
m_survey_outlet_id*
int(8)
nextval('mst.m_survey_detail_id_seq'::regclass)
Primary Key
survey_id**
int(8)
-
Foreign Key
mst.m_survey
outlet_id
int(8)
-
Foreign Key
mst.m_outlet
mst.m_survey_salesman
Deskripsi: Tabel ini adalah tabel master yang menyimpan semua data survey untuk salesman (distributor dan outlet yang dikunjungi)
Kolom (Columns / Fields):
Nama Kolom
Tipe Data
Nullable
Nilai Default
Deskripsi
cust_id
int(8)
-
id customer
m_survey_salesman_id*
int(8)
nextval('mst.m_survey_detail_id_seq'::regclass)
Primary Key
survey_id**
int(8)
-
Foreign Key
mst.m_survey
salesman_id**
int(8)
-
Foreign Key
mst.m_salesman key : emp_id
mst.m_survey_distributor
Deskripsi: Tabel ini adalah tabel master yang menyimpan semua data survey untuk distributor yang dikunjungi
Enhance :
3 Juli 2026 (create new table)
Kolom (Columns / Fields):
Nama Kolom
Tipe Data
Nullable
Nilai Default
Deskripsi
cust_id
int(8)
-
id customer
id*
int(8)
nextval('mst.m_survey_detail_id_seq'::regclass)
Primary Key
survey_id**
int(8)
-
Foreign Key
mst.m_survey
distributor_id**
int(8)
-
Foreign Key
mst.m_distributor key : emp_id
mst.m_survey_detail
Deskripsi: Tabel ini adalah tabel master yang menyimpan semua data surveyKolom (Columns / Fields):
Nama Kolom
Tipe Data
Nullable
Nilai Default
Deskripsi
cust_id
int(8)
-
id customer
m_survey_detail_id*
int(8)
nextval('mst.m_survey_detail_id_seq'::regclass)
Primary Key
survey_id**
int(8)
-
Foreign Key
mst.m_survey
survey_template_id
int(8)
-
Foreign Key
mst.m_survey_template
