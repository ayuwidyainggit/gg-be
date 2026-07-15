# Extracted from Enhance_Create_Survey_BE.docx

BE Analysis Document
Fitur: Enhance Survey
ScyllaX Distribution Management System
Informasi Dokumen
Detail
Versi
1.0.0
Modul
Manage Survey > Survey List
Status
Draft
Tanggal
2026
Tujuan
Analisa teknis backend untuk penambahan level target Distributor
API Survey
Item
Detail
Method
GET
Endpoint
{{url}}/master/v1/survey?page=1&limit=10&sort:created_date:desc&status=1
Headers
Accept
application/json
Authorization
Bearer Token
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4
Query Params
Field
Type
Required
Description
q
No
template_title
page
Integer
Yes
default : 1
limit
Integer
Yes
default : 5
sort
String
Yes
default created_date:desc
response_frequency
Array String
No
Mandatory, Optional
answer_frequency
Array String
No
One Time
Multiple Times, One Day
Multiple Times, Different Day
status
Integer
No
1: mst.m_survey_template.is_active = true
0 : mst.m_survey_template.is_active = false
Null : show all status
Example Request
curl --location -g \
'{{url}}/master/v1/survey?page=1&limit=5&sort=created_at:desc&status=1' \
--header 'Accept: application/json' \
--header 'Authorization: Bearer {{token}}'
Response
Nama Atribut
Type
Length
Description
message
String
150
Response Message
data
Array
-
-
survey_id
int
8
mst.m_survey.survey_id
created_at
timestamptz
6
mst.m_survey.created_at
answer_frequency
enum
Multiple, One Time
mst.m_survey.answer_frequency
survey_title
varchar
150
mst.m_survey.survey_title
response_type
enum
Mandatory, Optional
mst.m_survey.response_type
efective_date_start
date
mst.m_survey.efective_date_start
efective_date_end
date
mst.m_survey.efective_date_end
status
int
8
mst.m_survey.status
paging
Object
total_record
Numeric
11
Total data seluruh halaman
page_current
Numeric
11
Halaman saat ini
page_limit
Numeric
11
Data yang ditampilkan per page
page_total
Numeric
11
Total page keseluruhan
request_id
String
150
Generate request id
Example Response :
Case : success
{
"request_id": "REQ-202512210001",
"message": "Success",
"data": [
{
"survey_id": 12345678,
"created_at": "2025-12-01T10:15:30Z",
"answer_frequency": "Multiple",
"survey_title": "Customer Satisfaction Survey",
"response_type": "Mandatory",
"efective_date_start": "2025-12-01",
"efective_date_end": "2025-12-31",
"status": 1
},
{
"survey_id": 12345679,
"created_at": "2025-12-05T08:00:00Z",
"answer_frequency": "One Time",
"survey_title": "Employee Feedback Survey",
"response_type": "Optional",
"efective_date_start": "2025-12-05",
"efective_date_end": "2026-01-05",
"status": 1
}
],
"paging": {
"total_record": 25,
"page_current": 1,
"page_limit": 10,
"page_total": 3
}
}
Case : empty state / tidak terdapat data berdasarkan pencarian
{
"message": "No Data",
"data": null,
"paging": {
"total_record": 0,
"page_current": 1,
"page_limit": 10,
"page_total": 0
},
"request_id": "6915a6dd2395083c685e8e16"
}
Case : error
{
"message": "record not found",
"request_id": "6915a6dd2395083c685e8e16"
}
API Survey Detail
Item
Detail
Method
GET
Endpoint
{{url}}/master/v1/survey/:survey_id
Headers
Accept
application/json
Authorization
Bearer Token
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4
Path Variable
Field
Type
Required
Note
survey_id
Integer
Yes
Example Request
Example Request Default
curl --location -g \
'{{url}}/master/v1/survey/21' \
--header 'Accept: application/json' \
--header 'Authorization: Bearer {{token}}'
Response
Nama Atribut
Type
Length
Description
message
String
150
Response Message
data
Array
-
-
survey_id
int
8
mst.m_survey.survey_template_id
created_at
timestamptz
6
mst.m_survey.created_at
answer_frequency
enum
Multiple, One Time
mst.m_survey.answer_frequency
survey_title
varchar
150
mst.m_survey.survey_title
response_type
enum
Mandatory, Optional
mst.m_survey.response_type
efective_date_start
date
mst.m_survey.efective_date_start
efective_date_end
date
mst.m_survey.efective_date_end
level_target
varchar
150
mst.m_survey.level_target
status
int
8
mst.m_survey.status
target_survey
Object
target_type
enum
(national,)
mst.m_survey.target_type
emp_id
int(8)
mst.m_survey.emp_id
sales_name
mst.m_salesman.sales_name
business_unit
Array
target_cust_id
varchar
10
mst.m_survey_area.target_cust_id
target_cust_name
varchar
150
smc.m_customer.cust_name based on target_cust_id
area_id
int
8
mst.m_survey_area.area_id
area_name
mst.m_area.area_name
distributor_id
int
8
mst.m_survey_area.distributor_id
distributor_code
mst.m_distributor.distributor_code
distributor_name
mst.m_distributor.distributor_name
template
Array
survey_template_id
mst.m_survey_detail.survey_template_id
template_code
varchar
10
mst.m_survey_template.template_code
template_title
varchar
150
mst.m_survey_template.template_title
question_template
Array
question_template_id
int
8
mst.m_question_template.question_template_id
question
varchar
225
mst.m_question_template.question
answer_type
enum
(Single, Multiple, Free Text)
mst.m_question_template.answer_type
options
Array
q_option_template_id
int
8
mst.m_q_option_template.q_option_template_id
option
varchar
225
mst.m_q_option_template.option
outlet
Array
survey_outlet_id
int
8
mst.m_survey_outlet.survey_outlet_id
outlet_id
int
8
mst.m_survey_outlet.outlet_id
outlet_code
mst.m_outlet.outlet_code
outlet_name
mst.m_outlet.outlet_name
ot_class_id
int
8
mst.m_outlet.ot_class_id
ot_class_name
mst.m_outlet_class.ot_class_name
ot_grp_id
int
8
mst.m_outlet.ot_grp_id
ot_grp_name
mst.m_outlet_group.ot_grp_name
ot_type_id
int
8
mst.m_outlet.ot_type_id
ot_type_name
mst.m_outlet_type.ot_type_name
salesman
Array
m_survey_salesman_id
mst.m_survey_salesman.m_survey_salesman_id
sales_id
int
8
mst.m_survey_salesman.sales_id
sales_team_id
mst.m_salesman.sales_team_id
sales_team_name
mst.m_sales_team.sales_team_name
sales_name
mst.m_salesman.sales_name
target_distributor
Array
m_survey_salesman_id
mst.m_survey_distributor.id
distributor_id
mst.m_survey_distributor.distributor_id
distributor_code
mst.m_distributor.distributor_code
distributor_name
mst.m_distributor.distributor_name
paging
Object
total_record
Numeric
11
Total data seluruh halaman
page_current
Numeric
11
Halaman saat ini
page_limit
Numeric
11
Data yang ditampilkan per page
page_total
Numeric
11
Total page keseluruhan
request_id
String
150
Generate request id
Example Response :
Case : Level Target = salesman
{
"message": "Success",
"data": [
{
"survey_id": 10000001,
"created_at": "2025-12-01T09:30:00Z",
"answer_frequency": "Multiple",
"survey_title": "Survey Kepuasan Pelanggan",
"response_type": "Mandatory",
"efective_date_start": "2025-12-01",
"efective_date_end": "2025-12-31",
"status": 1,
"business_units": [
{
"target_cust_id": NULL,
"target_cust_id": NULL,
"area_id": 70,
"area_name": "Jakarta",
"distributor_id": 44,
"distributor_code": "D002",
"distributor_name": "Distributor Jakarta 2"
}
],
"template": [
{
"survey_template_id": 20001,
"template_code": "TMP01",
"template_title": "Template Survey Umum",
"question_template": [
{
"question_template_id": 30001,
"question": "Bagaimana kualitas pelayanan kami?",
"answer_type": "Multiple",
"options": [
{
"q_option_template_id": 40001,
"option": "Sangat Baik"
},
{
"q_option_template_id": 40002,
"option": "Baik"
},
{
"q_option_template_id": 40003,
"option": "Cukup"
},
{
"q_option_template_id": 40004,
"option": "Kurang"
}
]
},
{
"question_template_id": 30002,
"question": "Apakah Anda puas dengan layanan kami?",
"answer_type": "Single",
"options": [
{
"q_option_template_id": 40005,
"option": "Ya"
},
{
"q_option_template_id": 40006,
"option": "Tidak"
}
]
},
{
"question_template_id": 30003,
"question": "Saran untuk peningkatan layanan",
"answer_type": "Free Text",
"options": []
}
]
}
],
"outlet": [],
"salesman": [
{
"m_survey_salesman_id": 676,
"sales_id": 87654321,
"sales_team_id": 890,
"sales_team_name": "Team Budi",
"sales_name": "Budi Santoso"
}
],
"distributor": []
}
],
"paging": {
"total_record": 12,
"page_current": 1,
"page_limit": 10,
"page_total": 2
},
"request_id": "REQ-202512210002"
}
Case : Level Target = outlet
{
"message": "Success",
"data": [
{
"survey_id": 10000001,
"created_at": "2025-12-01T09:30:00Z",
"answer_frequency": "Multiple",
"survey_title": "Survey Kepuasan Pelanggan",
"response_type": "Mandatory",
"efective_date_start": "2025-12-01",
"efective_date_end": "2025-12-31",
"status": 1,
"business_units": [
{
"target_cust_id": "C220001",
"target_cust_id": "PT Sejahtera",
"area_id": NULL,
"area_name": NULL,
"distributor_id": NULL,
"distributor_code": NULL",
"distributor_name": NULL
}
],
"template": [
{
"survey_template_id": 20001,
"template_code": "TMP01",
"template_title": "Template Survey Umum",
"question_template": [
{
"question_template_id": 30001,
"question": "Bagaimana kualitas pelayanan kami?",
"answer_type": "Multiple",
"options": [
{
"q_option_template_id": 40001,
"option": "Sangat Baik"
},
{
"q_option_template_id": 40002,
"option": "Baik"
},
{
"q_option_template_id": 40003,
"option": "Cukup"
},
{
"q_option_template_id": 40004,
"option": "Kurang"
}
]
},
{
"question_template_id": 30002,
"question": "Apakah Anda puas dengan layanan kami?",
"answer_type": "Single",
"options": [
{
"q_option_template_id": 40005,
"option": "Ya"
},
{
"q_option_template_id": 40006,
"option": "Tidak"
}
]
},
{
"question_template_id": 30003,
"question": "Saran untuk peningkatan layanan",
"answer_type": "Free Text",
"options": []
}
]
}
],
"outlet": [
{
"survey_outlet_id": 50001,
"outlet_id": 1001,
"outlet_code": "OUT001",
"outlet_name": "Outlet Jakarta 1",
"ot_class_id": 10,
"ot_class_name": "Retail",
"ot_grp_id": 20,
"ot_grp_name": "Modern Trade",
"ot_type_id": 30,
"ot_type_name": "Supermarket"
}
],
"salesman": [],
"distributor": []
}
],
"paging": {
"total_record": 12,
"page_current": 1,
"page_limit": 10,
"page_total": 2
},
"request_id": "REQ-202512210002"
}
Case : Level Target = Distributor
{
"message": "Success",
"data": [
{
"survey_id": 10000001,
"created_at": "2025-12-01T09:30:00Z",
"answer_frequency": "Multiple",
"survey_title": "Survey Kepuasan Pelanggan",
"response_type": "Mandatory",
"efective_date_start": "2025-12-01",
"efective_date_end": "2025-12-31",
"status": 1,
"business_units": [
{
"target_cust_id": NULL,
"target_cust_id": NULL,
"area_id": 70,
"area_name": "Jakarta",
"distributor_id": 44,
"distributor_code": "D002",
"distributor_name": "Distributor Jakarta 2"
}
],
"template": [
{
"survey_template_id": 20001,
"template_code": "TMP01",
"template_title": "Template Survey Umum",
"question_template": [
{
"question_template_id": 30001,
"question": "Bagaimana kualitas pelayanan kami?",
"answer_type": "Multiple",
"options": [
{
"q_option_template_id": 40001,
"option": "Sangat Baik"
},
{
"q_option_template_id": 40002,
"option": "Baik"
},
{
"q_option_template_id": 40003,
"option": "Cukup"
},
{
"q_option_template_id": 40004,
"option": "Kurang"
}
]
},
{
"question_template_id": 30002,
"question": "Apakah Anda puas dengan layanan kami?",
"answer_type": "Single",
"options": [
{
"q_option_template_id": 40005,
"option": "Ya"
},
{
"q_option_template_id": 40006,
"option": "Tidak"
}
]
},
{
"question_template_id": 30003,
"question": "Saran untuk peningkatan layanan",
"answer_type": "Free Text",
"options": []
}
]
}
],
"outlet": [],
"salesman": [],
"distributor": [
{
"id": 50001,
"distributor_id": 1001,
"distributor_code": "OUT001",
"distributor_name": "Outlet Jakarta 1"
}
],
}
],
"paging": {
"total_record": 12,
"page_current": 1,
"page_limit": 10,
"page_total": 2
},
"request_id": "REQ-202512210002"
}
API Survey Create
Item
Detail
Method
POST
Endpoint
{{url}}/master/v1/survey
Headers
Accept
application/json
Authorization
Bearer Token
Body
Field
Type
Required
Note
survey_title
varchar(150)
Yes
efective_date_start
date
Yes
efective_date_end
date
Yes
answer_frequency
enum
One Time
Multiple Times, One Day
Multiple Times, Different Day
Yes
level_target
enum
Distributor
Salesman
Outlet
Yes
Tambahan 3 Juli 2026
response_type
enum
Mandatory, Optional
Yes
area_id
Array(Int)
Yes
distributor_id
Array(Int)
No
Mandatory jika business unit yang dipilih adalah distributor
target_cust_id
varchar(10)
No
Tambahan 3 Juli 2026
Mandatory jika business unit yang dipilih adalah principal
NULL jika business unit yang dipilih adalah distributor
outlet_id
Array(Int)
No
Mandatory jika level target = Outlet
emp_id
Array(Int)
No
Mandatory jika level target = Salesman
target_distributor_id
Array(Int)
No
Tambahan 3 Juli 2026
Mandatory jika level target = Distributor
survey_template_id
Array(Int)
Example Request
Level Target =  Salesman  dan Business Unit = Distributor
curl 'https://best.scyllax.online/master/v1/survey' \
-H 'sec-ch-ua-platform: "macOS"' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyIiwiZGlzdF9wcmljZV9ncnBfaWQiOjAsImRpc3RyaWJ1dG9yX2lkIjowLCJlbWFpbCI6InByaW5jZXNzYUBnbWFpbC5jb20iLCJlbXBfaWQiOjM4MCwiZXhwaXJlcyI6MTc4MzE1MzAwMywiaXNfYWRtaW4iOmZhbHNlLCJsYW5nX2lkIjoiaWQiLCJtb2JpbGVfbm8iOiIwODEzMzMzMzMzMzMiLCJwYXJlbnRfY3VzdF9pZCI6IkMyNjAwMiIsInVzZXJfZnVsbG5hbWUiOiJQcmluY2Vzc2EgQWhzYW5pIFRhcXdpbSIsInVzZXJfaWQiOjE0MCwidXNlcl9uYW1lIjoiUHJpbmNlc3NhIEFoc2FuaSBUYXF3aW0iLCJ3aGF0c2FwcCI6IjA4MTMzMzMzMzMzMyJ9.uuXCEkb4cGxoMIzq8l7NVAKXmmq88-c796BnbVryNDc' \
-H 'Referer: https://staging.scyllax.online/' \
-H 'sec-ch-ua: "Google Chrome";v="149", "Chromium";v="149", "Not)A;Brand";v="24"' \
-H 'sec-ch-ua-mobile: ?0' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36' \
-H 'Accept: application/json, text/plain, */*' \
-H 'Content-Type: application/json' \
--data-raw '{
"survey_title": "Survey Salesman",
"efective_date_start": "2026-07-01",
"efective_date_end": "2026-07-31",
"answer_frequency": "One Time",
"level_target": "Salesman",
"response_type": "Optional",
"target_cust_id": NULL,
"distributor_id": [
102,
103
],
"area_id": [
96,
91,
88
],
"outlet_id": [],
"emp_id": [
476,
480,
479
],
"target_distributor_id":NULL,
"survey_template_id": 53
}'
{
"survey_title": "Survey Salesman",
"efective_date_start": "2026-07-01",
"efective_date_end": "2026-07-31",
"answer_frequency": "One Time",
"level_target": "Salesman",
"response_type": "Optional",
"target_cust_id": NULL,
"distributor_id": [
102,
103
],
"area_id": [
96,
91,
88
],
"outlet_id": [],
"emp_id": [
476,
480,
479
],
"target_distributor_id":[],
"survey_template_id": 53
}
Impact Database :
mst.m_survey
survey_id
cust_id
survey_title
answer_frequency
response_type
level_target
emp_id
efective_date_start
efective_date_end
status
is_del
created_at
created_by
updated_at
updated_by
deleted_at
deleted_by
1
CUST001
Survey Salesman
One Time
Optional
Salesman
476
2026-07-01
2026-07-31
1
false
2026-07-01 08:00:00+07
1001
NULL
NULL
NULL
NULL
mst.m_survey_area
survey_area_id
survey_id
distributor_id
area_id
targer_cust_id
is_del
1
1
102
88 (area_id dari distributor)
NULL
false
2
1
103
88
NULL
false
mst.m_survey_salesman
m_survey_salesman_id
cust_id
survey_id
salesman_id
is_del
1
CUST001
1
476
false
2
CUST001
1
480
false
3
CUST001
2
479
false
mst.m_survey_detail
survey_detail_id
survey_id
survey_template_id
is_del
1
1
53
false
Level Target =  Outlet  dan Business Unit = Distributor
curl 'https://best.scyllax.online/master/v1/survey' \
-H 'sec-ch-ua-platform: "macOS"' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyIiwiZGlzdF9wcmljZV9ncnBfaWQiOjAsImRpc3RyaWJ1dG9yX2lkIjowLCJlbWFpbCI6InByaW5jZXNzYUBnbWFpbC5jb20iLCJlbXBfaWQiOjM4MCwiZXhwaXJlcyI6MTc4MzE1MzAwMywiaXNfYWRtaW4iOmZhbHNlLCJsYW5nX2lkIjoiaWQiLCJtb2JpbGVfbm8iOiIwODEzMzMzMzMzMzMiLCJwYXJlbnRfY3VzdF9pZCI6IkMyNjAwMiIsInVzZXJfZnVsbG5hbWUiOiJQcmluY2Vzc2EgQWhzYW5pIFRhcXdpbSIsInVzZXJfaWQiOjE0MCwidXNlcl9uYW1lIjoiUHJpbmNlc3NhIEFoc2FuaSBUYXF3aW0iLCJ3aGF0c2FwcCI6IjA4MTMzMzMzMzMzMyJ9.uuXCEkb4cGxoMIzq8l7NVAKXmmq88-c796BnbVryNDc' \
-H 'Referer: https://staging.scyllax.online/' \
-H 'sec-ch-ua: "Google Chrome";v="149", "Chromium";v="149", "Not)A;Brand";v="24"' \
-H 'sec-ch-ua-mobile: ?0' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36' \
-H 'Accept: application/json, text/plain, */*' \
-H 'Content-Type: application/json' \
--data-raw '{
"survey_title": "Survey Outlet",
"efective_date_start": "2026-07-01",
"efective_date_end": "2026-07-30",
"answer_frequency": "Multiple Times, One Day",
"level_target": "Outlet",
"response_type": "Mandatory",
"target_type": "Specific",
"target_cust_id": NULL,
"distributor_id": [
102,
119,
103
],
"area_id": [
96,
91,
88
],
"outlet_id": [
3489,
3488
],
"emp_id": [],
"target_distributor_id":[],
"survey_template_id": 53
}'
{
"survey_title": "Survey Outlet",
"efective_date_start": "2026-07-01",
"efective_date_end": "2026-07-30",
"answer_frequency": "Multiple Times, One Day",
"level_target": "Outlet",
"response_type": "Mandatory",
"target_type": "Specific",
"target_cust_id": NULL,
"distributor_id": [
102,
119,
103
],
"area_id": [
96,
91,
88
],
"outlet_id": [
3489,
3488
],
"emp_id": [],
"target_distributor_id":[],
"survey_template_id": 53
}
Impact Database :
mst.m_survey
survey_id
cust_id
survey_title
answer_frequency
response_type
level_target
emp_id
efective_date_start
efective_date_end
status
is_del
created_at
created_by
updated_at
updated_by
deleted_at
deleted_by
1
CUST001
Survey Outlet
Multiple Times, One Day
Mandatory
Outlet
NULL
2026-07-01
2026-07-30
1
false
2026-07-01 08:00:00+07
1001
NULL
NULL
NULL
NULL
mst.m_survey_area
survey_area_id
survey_id
distributor_id
area_id
targer_cust_id
is_del
1
1
102
88 (area_id dari distributor)
NULL
false
2
1
119
88
NULL
false
3
1
103
88
NULL
false
mst.m_survey_outlet
m_survey_outlet_id
cust_id
survey_id
outlet_id
is_del
1
CUST001
1
3489
false
2
CUST001
1
3488
false
mst.m_survey_detail
survey_detail_id
survey_id
survey_template_id
is_del
1
1
53
false
Level Target =  Distributor  dan Business Unit = Distributor
curl 'https://best.scyllax.online/master/v1/survey' \
-H 'sec-ch-ua-platform: "macOS"' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyIiwiZGlzdF9wcmljZV9ncnBfaWQiOjAsImRpc3RyaWJ1dG9yX2lkIjowLCJlbWFpbCI6InByaW5jZXNzYUBnbWFpbC5jb20iLCJlbXBfaWQiOjM4MCwiZXhwaXJlcyI6MTc4MzE1MzAwMywiaXNfYWRtaW4iOmZhbHNlLCJsYW5nX2lkIjoiaWQiLCJtb2JpbGVfbm8iOiIwODEzMzMzMzMzMzMiLCJwYXJlbnRfY3VzdF9pZCI6IkMyNjAwMiIsInVzZXJfZnVsbG5hbWUiOiJQcmluY2Vzc2EgQWhzYW5pIFRhcXdpbSIsInVzZXJfaWQiOjE0MCwidXNlcl9uYW1lIjoiUHJpbmNlc3NhIEFoc2FuaSBUYXF3aW0iLCJ3aGF0c2FwcCI6IjA4MTMzMzMzMzMzMyJ9.uuXCEkb4cGxoMIzq8l7NVAKXmmq88-c796BnbVryNDc' \
-H 'Referer: https://staging.scyllax.online/' \
-H 'sec-ch-ua: "Google Chrome";v="149", "Chromium";v="149", "Not)A;Brand";v="24"' \
-H 'sec-ch-ua-mobile: ?0' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36' \
-H 'Accept: application/json, text/plain, */*' \
-H 'Content-Type: application/json' \
--data-raw '{
"survey_title": "Survey Outlet",
"efective_date_start": "2026-07-01",
"efective_date_end": "2026-07-30",
"answer_frequency": "Multiple Times, One Day",
"level_target": "Distributor",
"response_type": "Mandatory",
"target_type": "Specific",
"target_cust_id": NULL,
"distributor_id": [
102,
119,
103
],
"area_id": [
96,
91,
88
],
"outlet_id":[],
"target_distributor_id":[
3489,
3488
],
"emp_id": [],
"survey_template_id": 53
}
'
{
"survey_title": "Survey Outlet",
"efective_date_start": "2026-07-01",
"efective_date_end": "2026-07-30",
"answer_frequency": "Multiple Times, One Day",
"level_target": "Distributor",
"response_type": "Mandatory",
"target_type": "Specific",
"target_cust_id": NULL,
"distributor_id": [
102,
119,
103
],
"area_id": [
96,
91,
88
],
"outlet_id":[],
"target_distributor_id":[
3489,
3488
],
"emp_id": [],
"survey_template_id": 53
}
Impact Database :
mst.m_survey
survey_id
cust_id
survey_title
answer_frequency
response_type
level_target
emp_id
efective_date_start
efective_date_end
status
is_del
created_at
created_by
updated_at
updated_by
deleted_at
deleted_by
1
CUST001
Survey Outlet
Multiple Times, One Day
Mandatory
Distributor
NULL
2026-07-01
2026-07-30
1
false
2026-07-01 08:00:00+07
1001
NULL
NULL
NULL
NULL
mst.m_survey_area
survey_area_id
survey_id
distributor_id
area_id
targer_cust_id
is_del
1
1
102
88 (area_id dari distributor)
NULL
false
2
1
119
88
NULL
false
3
1
103
88
NULL
false
mst.m_survey_distributor
m_survey_distributor_id
cust_id
survey_id
distributor
is_del
1
CUST001
1
3489
false
2
CUST001
1
3488
false
mst.m_survey_detail
survey_detail_id
survey_id
survey_template_id
is_del
1
1
53
false
Level Target =  Salesman  dan Business Unit = Principal
curl 'https://best.scyllax.online/master/v1/survey' \
-H 'sec-ch-ua-platform: "macOS"' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyIiwiZGlzdF9wcmljZV9ncnBfaWQiOjAsImRpc3RyaWJ1dG9yX2lkIjowLCJlbWFpbCI6InByaW5jZXNzYUBnbWFpbC5jb20iLCJlbXBfaWQiOjM4MCwiZXhwaXJlcyI6MTc4MzE1MzAwMywiaXNfYWRtaW4iOmZhbHNlLCJsYW5nX2lkIjoiaWQiLCJtb2JpbGVfbm8iOiIwODEzMzMzMzMzMzMiLCJwYXJlbnRfY3VzdF9pZCI6IkMyNjAwMiIsInVzZXJfZnVsbG5hbWUiOiJQcmluY2Vzc2EgQWhzYW5pIFRhcXdpbSIsInVzZXJfaWQiOjE0MCwidXNlcl9uYW1lIjoiUHJpbmNlc3NhIEFoc2FuaSBUYXF3aW0iLCJ3aGF0c2FwcCI6IjA4MTMzMzMzMzMzMyJ9.uuXCEkb4cGxoMIzq8l7NVAKXmmq88-c796BnbVryNDc' \
-H 'Referer: https://staging.scyllax.online/' \
-H 'sec-ch-ua: "Google Chrome";v="149", "Chromium";v="149", "Not)A;Brand";v="24"' \
-H 'sec-ch-ua-mobile: ?0' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36' \
-H 'Accept: application/json, text/plain, */*' \
-H 'Content-Type: application/json' \
--data-raw '{
"survey_title": "Survey Salesman",
"efective_date_start": "2026-07-01",
"efective_date_end": "2026-07-31",
"answer_frequency": "One Time",
"level_target": "Salesman",
"response_type": "Optional",
"target_cust_id": "C22001",
"distributor_id": [],
"area_id": [
96,
91,
88
],
"outlet_id": [],
"emp_id": [
476,
480,
479
],
"target_distributor_id":NULL,
"survey_template_id": 53
}'
{
"survey_title": "Survey Salesman",
"efective_date_start": "2026-07-01",
"efective_date_end": "2026-07-31",
"answer_frequency": "One Time",
"level_target": "Salesman",
"response_type": "Optional",
"target_cust_id": "C22001",
"distributor_id": [],
"area_id": [
96,
91,
88
],
"outlet_id": [],
"emp_id": [
476,
480,
479
],
"target_distributor_id":[],
"survey_template_id": 53
}
Impact Database :
mst.m_survey
survey_id
cust_id
survey_title
answer_frequency
response_type
level_target
emp_id
efective_date_start
efective_date_end
status
is_del
created_at
created_by
updated_at
updated_by
deleted_at
deleted_by
1
CUST001
Survey Salesman
One Time
Optional
Salesman
476
2026-07-01
2026-07-31
1
false
2026-07-01 08:00:00+07
1001
NULL
NULL
NULL
NULL
mst.m_survey_area
survey_area_id
survey_id
distributor_id
area_id
targer_cust_id
is_del
1
1
NULL
NULL
C22001
false
mst.m_survey_salesman
m_survey_salesman_id
cust_id
survey_id
salesman_id
is_del
1
CUST001
1
476
false
2
CUST001
1
480
false
3
CUST001
2
479
false
mst.m_survey_detail
survey_detail_id
survey_id
survey_template_id
is_del
1
1
53
false
Level Target =  Outlet  dan Business Unit = Principal
curl 'https://best.scyllax.online/master/v1/survey' \
-H 'sec-ch-ua-platform: "macOS"' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyIiwiZGlzdF9wcmljZV9ncnBfaWQiOjAsImRpc3RyaWJ1dG9yX2lkIjowLCJlbWFpbCI6InByaW5jZXNzYUBnbWFpbC5jb20iLCJlbXBfaWQiOjM4MCwiZXhwaXJlcyI6MTc4MzE1MzAwMywiaXNfYWRtaW4iOmZhbHNlLCJsYW5nX2lkIjoiaWQiLCJtb2JpbGVfbm8iOiIwODEzMzMzMzMzMzMiLCJwYXJlbnRfY3VzdF9pZCI6IkMyNjAwMiIsInVzZXJfZnVsbG5hbWUiOiJQcmluY2Vzc2EgQWhzYW5pIFRhcXdpbSIsInVzZXJfaWQiOjE0MCwidXNlcl9uYW1lIjoiUHJpbmNlc3NhIEFoc2FuaSBUYXF3aW0iLCJ3aGF0c2FwcCI6IjA4MTMzMzMzMzMzMyJ9.uuXCEkb4cGxoMIzq8l7NVAKXmmq88-c796BnbVryNDc' \
-H 'Referer: https://staging.scyllax.online/' \
-H 'sec-ch-ua: "Google Chrome";v="149", "Chromium";v="149", "Not)A;Brand";v="24"' \
-H 'sec-ch-ua-mobile: ?0' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36' \
-H 'Accept: application/json, text/plain, */*' \
-H 'Content-Type: application/json' \
--data-raw '{
"survey_title": "Survey Outlet",
"efective_date_start": "2026-07-01",
"efective_date_end": "2026-07-30",
"answer_frequency": "Multiple Times, One Day",
"level_target": "Outlet",
"response_type": "Mandatory",
"target_type": "Specific",
"target_cust_id": "C22001",
"distributor_id": [],
"area_id": [
96,
91,
88
],
"outlet_id": [
3489,
3488
],
"emp_id": [],
"target_distributor_id":[],
"survey_template_id": 53
}'
{
"survey_title": "Survey Outlet",
"efective_date_start": "2026-07-01",
"efective_date_end": "2026-07-30",
"answer_frequency": "Multiple Times, One Day",
"level_target": "Outlet",
"response_type": "Mandatory",
"target_type": "Specific",
"target_cust_id": "C22001",
"distributor_id": [],
"area_id": [
96,
91,
88
],
"outlet_id": [
3489,
3488
],
"emp_id": [],
"target_distributor_id":[],
"survey_template_id": 53
}
Impact Database :
mst.m_survey
survey_id
cust_id
survey_title
answer_frequency
response_type
level_target
emp_id
efective_date_start
efective_date_end
status
is_del
created_at
created_by
updated_at
updated_by
deleted_at
deleted_by
1
CUST001
Survey Outlet
Multiple Times, One Day
Mandatory
Outlet
NULL
2026-07-01
2026-07-30
1
false
2026-07-01 08:00:00+07
1001
NULL
NULL
NULL
NULL
mst.m_survey_area
survey_area_id
survey_id
distributor_id
area_id
targer_cust_id
is_del
1
1
NULL
NULL
C22001
false
mst.m_survey_outlet
m_survey_outlet_id
cust_id
survey_id
outlet_id
is_del
1
CUST001
1
3489
false
2
CUST001
1
3488
false
mst.m_survey_detail
survey_detail_id
survey_id
survey_template_id
is_del
1
1
53
false
Level Target =  Distributor  dan Business Unit = Principal
curl 'https://best.scyllax.online/master/v1/survey' \
-H 'sec-ch-ua-platform: "macOS"' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyIiwiZGlzdF9wcmljZV9ncnBfaWQiOjAsImRpc3RyaWJ1dG9yX2lkIjowLCJlbWFpbCI6InByaW5jZXNzYUBnbWFpbC5jb20iLCJlbXBfaWQiOjM4MCwiZXhwaXJlcyI6MTc4MzE1MzAwMywiaXNfYWRtaW4iOmZhbHNlLCJsYW5nX2lkIjoiaWQiLCJtb2JpbGVfbm8iOiIwODEzMzMzMzMzMzMiLCJwYXJlbnRfY3VzdF9pZCI6IkMyNjAwMiIsInVzZXJfZnVsbG5hbWUiOiJQcmluY2Vzc2EgQWhzYW5pIFRhcXdpbSIsInVzZXJfaWQiOjE0MCwidXNlcl9uYW1lIjoiUHJpbmNlc3NhIEFoc2FuaSBUYXF3aW0iLCJ3aGF0c2FwcCI6IjA4MTMzMzMzMzMzMyJ9.uuXCEkb4cGxoMIzq8l7NVAKXmmq88-c796BnbVryNDc' \
-H 'Referer: https://staging.scyllax.online/' \
-H 'sec-ch-ua: "Google Chrome";v="149", "Chromium";v="149", "Not)A;Brand";v="24"' \
-H 'sec-ch-ua-mobile: ?0' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36' \
-H 'Accept: application/json, text/plain, */*' \
-H 'Content-Type: application/json' \
--data-raw '{
"survey_title": "Survey Outlet",
"efective_date_start": "2026-07-01",
"efective_date_end": "2026-07-30",
"answer_frequency": "Multiple Times, One Day",
"level_target": "Distributor",
"response_type": "Mandatory",
"target_type": "Specific",
"target_cust_id": "C22001",
"distributor_id": [],
"area_id": [
96,
91,
88
],
"outlet_id":[],
"target_distributor_id":[
3489,
3488
],
"emp_id": [],
"survey_template_id": 53
}
'
{
"survey_title": "Survey Outlet",
"efective_date_start": "2026-07-01",
"efective_date_end": "2026-07-30",
"answer_frequency": "Multiple Times, One Day",
"level_target": "Distributor",
"response_type": "Mandatory",
"target_type": "Specific",
"target_cust_id": "C22001",
"distributor_id": [],
"area_id": [
96,
91,
88
],
"outlet_id":[],
"target_distributor_id":[
3489,
3488
],
"emp_id": [],
"survey_template_id": 53
}
Impact Database :
mst.m_survey
survey_id
cust_id
survey_title
answer_frequency
response_type
level_target
emp_id
efective_date_start
efective_date_end
status
is_del
created_at
created_by
updated_at
updated_by
deleted_at
deleted_by
1
CUST001
Survey Outlet
Multiple Times, One Day
Mandatory
Distributor
NULL
2026-07-01
2026-07-30
1
false
2026-07-01 08:00:00+07
1001
NULL
NULL
NULL
NULL
mst.m_survey_area
survey_area_id
survey_id
distributor_id
area_id
targer_cust_id
is_del
1
1
NULL
NULL
C22001
false
mst.m_survey_distributor
m_survey_distributor_id
cust_id
survey_id
distributor
is_del
1
CUST001
1
3489
false
2
CUST001
1
3488
false
mst.m_survey_detail
survey_detail_id
survey_id
survey_template_id
is_del
1
1
53
false
Response
Case : sukses dan terdapat data
{
"message": "Survey has been successfully created",
"request_id": "6915a5e8e3f53f84fe73517f"
}
case : Survey Title tidak unik per periode aktif.
{
"message": "Survey Title already exists for the active period",
"request_id": "6915a5e8e3f53f84fe73517f"
}
Case : error
{
"message": "Failed to create survey data",
"request_id": "6915a5e8e3f53f84fe73517f"
}
Impact ke database :
mst.m_survey
mst.m_survey_area
mst.m_survey_detail
mst.m_survey_outlet
mst.m_survey_salesman
mst.m_survey_distributor
login principal = bisa input buat salesman punya principal itu sendiri, dan salesman  punya distributor di bawahnya
login distributor = hanya bisa input salesman punya distributor yang login
API Survey Edit
Create new endpoint :
Content-Type : application/json
Method: PUT
URL: {{url}}/master/v1/survey/:survey_id
Headers
Accept
application/json
Authorization
Bearer Token
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4
Body
Field
Type
Required
Note
survey_title
varchar(150)
Yes
efective_date_start
date
Yes
efective_date_end
date
Yes
answer_frequency
enum
One Time
Multiple Times, One Day
Multiple Times, Different Day
Yes
level_target
enum
Distributor
Salesman
Outlet
Yes
Tambahan 3 Juli 2026
response_type
enum
Mandatory, Optional
Yes
area_id
Array(Int)
Yes
distributor_id
Array(Int)
No
Mandatory jika business unit yang dipilih adalah distributor
target_cust_id
varchar(10)
No
Tambahan 3 Juli 2026
Mandatory jika business unit yang dipilih adalah principal
NULL jika business unit yang dipilih adalah distributor
outlet_id
Array(Int)
No
Mandatory jika level target = Outlet
emp_id
Array(Int)
No
Mandatory jika level target = Salesman
target_distributor_id
Array(Int)
No
Tambahan 3 Juli 2026
Mandatory jika level target = Distributor
survey_template_id
Array(Int)
Example Request
Example Request Default
curl -X PUT 'https://best.scyllax.online/master/v1/survey/123' \
-H 'sec-ch-ua-platform: "macOS"' \
-H 'Authorization: Bearer <YOUR_TOKEN>' \
-H 'Referer: https://staging.scyllax.online/' \
-H 'sec-ch-ua: "Google Chrome";v="149", "Chromium";v="149", "Not)A;Brand";v="24"' \
-H 'sec-ch-ua-mobile: ?0' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36' \
-H 'Accept: application/json, text/plain, */*' \
-H 'Content-Type: application/json' \
--data-raw '{
"survey_title": "Survey Salesman",
"efective_date_start": "2026-07-01",
"efective_date_end": "2026-07-31",
"answer_frequency": "One Time",
"level_target": "Salesman",
"response_type": "Optional",
"target_cust_id": null,
"distributor_id": [
102,
103
],
"area_id": [
96,
91,
88
],
"outlet_id": [],
"emp_id": [
476,
480,
479
],
"target_distributor_id": null,
"survey_template_id": 53
}'
Response :
Case : sukses dan terdapat data
{
"message": "Survey has been successfully updated",
"request_id": "6915a5e8e3f53f84fe73517f"
}
Case : error
{
"message": "Failed to update survey data",
"request_id": "6915a5e8e3f53f84fe73517f"
}
API Survey Deactive
Create new endpoint :
Content-Type : application/json
Method: PATCH
URL: {{url}}/master/v1/survey/:survey_id
Headers
Accept
application/json
Authorization
Bearer Token
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4
Body
Field
Type
Required
Note
is_active
bool
false
Example Request
Example Request Default
curl --location \
--request POST '{{url}}/master/v1/survey/:survey_id' \
--header 'Accept: application/json' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {{token}}' \
--data '{
"is_active": false
}'
Response
Response :
Case : sukses dan terdapat data
{
"message": "Survey successfully deactivated",
"request_id": "6915a5e8e3f53f84fe73517f"
}
Case : error
{
"message": "Survey not found",
"request_id": "6915a5e8e3f53f84fe73517f"
}
API  Outlet LIST
Create new endpoint :
Content-Type : application/json
Method: GET
URL: 103.28.219.73/v1/outlets?page=1&limit=10&sort=outlet_id:desc&is_active=1&identity_type=National ID&identity_no=312123456789098765
Link doc eksisting
Headers
Accept
application/json
Authorization
Bearer Token
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3Mjk3MzUyMDEsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.a0LeOmfzKXEhyPPc0UrDxZIyOJXYGxfcr49LaW7ksp4
Params
Field
Type
Required
Note
page
bool
false
limit
sort
q
is_active
tidak perlu
outlet_status
1,5,6,7,
outlet_id
verification_status
1
ot_grp_id
integer
No
ot_type_id
integer
No
identity_type
identity_no
distributor_id
contoh : 0, 67, 680 = principal 67, 68 = distributor
ot_class_id
integer
No
Filter berdasarkan mst.m_outlet.ot_class_id
Contoh Case:
User login : princessa@gmail.com | Admin123
Filter :
outlet berdasarkan PT Besi Makmur
ot_type = Type S , type A
class name = Class Outdoor
group = group R
curl :
curl 'https://best.scyllax.online/master/v1/outlets?page=1&limit=10&sort=outlet_id:desc&is_active=1&ot_type_id=99,97&ot_grp_id=111&ot_class_id=91&q=&page=1&limit=70' \
-H 'sec-ch-ua-platform: "macOS"' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyIiwiZGlzdF9wcmljZV9ncnBfaWQiOjAsImRpc3RyaWJ1dG9yX2lkIjowLCJlbWFpbCI6InByaW5jZXNzYUBnbWFpbC5jb20iLCJlbXBfaWQiOjM4MCwiZXhwaXJlcyI6MTc3NjE2MDM4NiwiaXNfYWRtaW4iOmZhbHNlLCJsYW5nX2lkIjoiaWQiLCJtb2JpbGVfbm8iOiIwODEzMzMzMzMzMzMiLCJwYXJlbnRfY3VzdF9pZCI6IkMyNjAwMiIsInVzZXJfZnVsbG5hbWUiOiJQcmluY2Vzc2EgQWhzYW5pIFRhcXdpbSIsInVzZXJfaWQiOjE0MCwidXNlcl9uYW1lIjoiUHJpbmNlc3NhIEFoc2FuaSBUYXF3aW0iLCJ3aGF0c2FwcCI6IjA4MTMzMzMzMzMzMyJ9.O6RyO5a6idD-AHxPaoGff1_z1y9l5K0yAwgYllUGpxQ' \
-H 'Referer: https://staging.scyllax.online/' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36' \
-H 'Accept: application/json, text/plain, */*' \
-H 'sec-ch-ua: "Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"' \
-H 'sec-ch-ua-mobile: ?0'
response eksisting (seharusnya ada 10 data dengan filter di atas)
query BE :
select mo.outlet_code , mo.outlet_name, mot.ot_type_name , moc.ot_class_name , mog.ot_grp_name  from mst.m_outlet mo
join mst.m_outlet_type mot
on mot.ot_type_id = mo.ot_type_id
join mst.m_outlet_class moc
on moc.ot_class_id =mo.ot_class_id
join mst.m_outlet_group mog
on mog.ot_grp_id =mo.ot_grp_id
join smc.m_customer mc
on mc.cust_id = mo.cust_id
where mc.distributor_id =102 and mo.ot_type_id in(99, 97) and mo.ot_grp_id in(111) and mo.ot_class_id in (91)
response yang seharusnya tampil :
Enhance FE  :
FE belum kirim distributor_id
FE belum kirim ot_class_id
API Salesman Lookup
Create new endpoint :
Content-Type : application/json
Method: GET
URL: {{url}}/v1/salesman?q=&page=1&limit=10&sort=sales_name:asc&sales_team_id=24,1
Doc eksisting
Headers
Accept
application/json
Authorization
Bearer Token
Body
Enhance : add param berikut
distributor_id → Optional (int[])
Data Test : Email : adminbm@gmail.com
Cust_id : C260020001
Distributor : PT. Besi Makmur  | distributor_id = 102
Sales Team yang dipilih:
- GT  | sales_team_id = 66
- MIX | sales_team_id= 65
Lookup Salesman yang tampil berdasarkan filter distributor dan sales team :
•⁠  ⁠Piere Njangka
•⁠  ⁠Jaka
Query :
user distributor :
select mc.cust_id , mc.distributor_id ,  ms.* from mst.m_salesman ms
join smc.m_customer mc
on mc.cust_id = ms.cust_id
where ms.sales_team_id in(66,65) and ms.cust_id = 'C260020001' and ms.is_active = true and ms.is_del = false
and mc.distributor_id =102
Ketika user principal  filter salesman berdasarkan User Principal :
select mc.cust_id , mc.distributor_id ,  ms.* from mst.m_salesman ms
join smc.m_customer mc
on mc.cust_id = ms.cust_id
where ms.sales_team_id in(78,77) and ms.cust_id like 'C26002%' and ms.is_active = true and ms.is_del = false
Ketika user principal  filter salesman berdasarkan User Principal + distributor :
select mc.cust_id , mc.distributor_id ,  ms.* from mst.m_salesman ms
join smc.m_customer mc
on mc.cust_id = ms.cust_id
where ms.sales_team_id in(78,77) and ms.cust_id like 'C26002%' and ms.is_active = true and ms.is_del = false
and mc.distributor_id =102
API Salesman Team Lookup
Create new endpoint :
Content-Type : application/json
Method: GET
URL: {{url}}/master/v1/sales-teams?mode=lookup&page=1&limit=70&distributor_id=1&distributor_id=45&q=&page=1&limit=70
Params
Field
Type
Required
Note
q
mode
No
page
Integer
Yes
limit
Integer
Yes
distributor_id
int[]
Yes
enhance tambahan : tambahkan distributor_id = 0 → artinya mencari sales_team berdasarkan user milik principal
secara business rules, sales_team dapat dipilih berdasarkan multiple dropdown dari business_unit, dimana resp business_unit adalah user principal itu sendiri dan distributor dibawahnya.
contoh : 0, 67, 680 = principal 67, 68 = distributor
Issue :
User Login :  princ@idetama.id
Business Unit yang di pilih :
Admin Principal 1  (ini user principal itu sendiri)
Distributor iDetama (distributor_id = 67 )
CURL :
curl 'https://best.scyllax.online/master/v1/sales-teams?mode=lookup&page=1&limit=70&distributor_id=1&distributor_id=67&q=&page=1&limit=70' \
-H 'Accept: application/json, text/plain, */*' \
-H 'Accept-Language: en-US,en;q=0.9' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxIiwiZGlzdF9wcmljZV9ncnBfaWQiOjAsImRpc3RyaWJ1dG9yX2lkIjowLCJlbWFpbCI6InByaW5jQGlkZXRhbWEuaWQiLCJlbXBfaWQiOjI3OCwiZXhwaXJlcyI6MTc3NjEzNTU1MSwiaXNfYWRtaW4iOnRydWUsImxhbmdfaWQiOiJlbiIsIm1vYmlsZV9ubyI6IjA4MTEzMjMyMzMyIiwicGFyZW50X2N1c3RfaWQiOiJDMjIwMDEiLCJ1c2VyX2Z1bGxuYW1lIjoiQWRtaW4gUHJpbmNpcGFsIDEiLCJ1c2VyX2lkIjoxLCJ1c2VyX25hbWUiOiJwcmluY0BpZGV0YW1hLmlkIiwid2hhdHNhcHAiOiIwODExMzIzMjMzMiJ9.j7sy-VIHtURjIvINFis-ljcw0SeicRUYOPRO7NEjNi8' \
-H 'Connection: keep-alive' \
-H 'Origin: https://staging.scyllax.online' \
-H 'Referer: https://staging.scyllax.online/' \
-H 'Sec-Fetch-Dest: empty' \
-H 'Sec-Fetch-Mode: cors' \
-H 'Sec-Fetch-Site: same-site' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36' \
-H 'sec-ch-ua: "Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"' \
-H 'sec-ch-ua-mobile: ?0' \
-H 'sec-ch-ua-platform: "macOS"'
Payload FE :
mode=lookup&
page=1&
limit=70&
distributor_id=1&
distributor_id=67&
q=&
page=1
&limit=70
Response :
{
"message": "",
"data": [
{
"sales_team_id": 60,
"sales_team_code": "900",
"sales_team_name": "Team HO"
},
{
"sales_team_id": 50,
"sales_team_code": "200",
"sales_team_name": "CANVAS"
},
{
"sales_team_id": 49,
"sales_team_code": "100",
"sales_team_name": "TAKING ORDER"
},
{
"sales_team_id": 48,
"sales_team_code": "000",
"sales_team_name": "SHOP SALES"
}
],
"paging": {
"total_record": 4,
"page_current": 1,
"page_limit": 70,
"page_total": 1
},
"request_id": "69dc5f9fc7fd667f49595760"
}
Fixing FE : untuk case business unit Admin Principal 1 , dan Distributor iDetama distributor seharusnya hanya di kirim 67 dan 0
Fixing BE : BE seharusnya filter berdasarkan user principal dan distributor idetama
Query :
SELECT
mst.cust_id,
md.distributor_id,
md.distributor_name,
mst.sales_team_id,
mst.sales_team_code,
mst.sales_team_name
FROM mst.m_sales_team mst
JOIN smc.m_customer mc
ON mst.cust_id = mc.cust_id
LEFT JOIN mst.m_distributor md
ON md.distributor_id = mc.distributor_id
WHERE mst.is_active = true
AND mst.is_del = false
AND (
md.distributor_id = :JWT_DISTRIBUTOR_ID  -- Parameter ID Distributor (Contoh: 67)
OR mst.cust_id = :PARENT_CUST_ID         -- Parameter Cust ID Principal (Contoh: 'C22001')
);
Issue
Issue SX 910
https://scyllax-pratesis.atlassian.net/browse/SX-910
mas yogi tolong survey_title dibuat unik API Create
Method: POST
URL: {{url}}/master/v1/survey
API Update
Method: PUT
URL: {{url}}/master/v1/survey/:survey_idsurvey_title tidak boleh sama jika effective date-nya overlap
contoh 1:
sekarang tgl 1 jan 2026
•⁠ ⁠S001 = saya buat data untuk tanggal 1 Feb 2026 - 28 Feb 2026 —> Survey A
•⁠ ⁠S002 = saya buat data untuk tanggal 1 Maret 2026 - 31 Maret 2026 —> Survey A
Success , karena beda efective date
contoh 2:
sekarang tgl 1 jan 2026
•⁠ ⁠S001 = saya buat data untuk tanggal 1 Feb 2026 - 28 Feb 2026 —> Survey A
•⁠ ⁠S002 = saya buat data untuk tanggal 1 Feb 2026 - 15 Feb 2026 —> Survey A ini
Failed, karena beda efective date
Note :
"Survey A" = "survey a" dianggap sama.
if effective_date_start > effective_date_end:
return error "effective_date_start must be <= effective_date_end"
exists = query overlap_check(survey_title, effective_date_start, effective_date_end, survey_id_if_update)
if exists:
return error "Survey title already exists in overlapping effective date range"
Issue Create Survey : (case salesman)
Issue BE :
Tabel :
mst.m_survey_area = kurang field distributor_id
kurang tabel mst.m_survey_salesman https://docs.google.com/document/d/1FRz0ym2cxYIlwCqjviwNLqOr4fyhZOZI2QCMyJ4wtmE/edit?tab=t.0
perbaiki bagian post
CURL :
curl 'https://best.scyllax.online/master/v1/survey' \
-H 'sec-ch-ua-platform: "macOS"' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzI2MDAyMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjowLCJkaXN0cmlidXRvcl9pZCI6MTAyLCJlbWFpbCI6ImFkbWluYm1AZ21haWwuY29tIiwiZW1wX2lkIjozODEsImV4cGlyZXMiOjE3NzU4MzA1NDYsImlzX2FkbWluIjpmYWxzZSwibGFuZ19pZCI6ImlkIiwibW9iaWxlX25vIjoiMDgxMzMzMzMzMzMzIiwicGFyZW50X2N1c3RfaWQiOiJDMjYwMDIiLCJ1c2VyX2Z1bGxuYW1lIjoiUGhpbGwgSm9uZXMiLCJ1c2VyX2lkIjoxNDEsInVzZXJfbmFtZSI6IlBoaWxsIEpvbmVzIiwid2hhdHNhcHAiOiIwODEzMzMzMzMzMzMifQ.WLGNdCqV5BoxDTL0iG9LU7HTbmp7EhmS-FuEsHi9Z-E' \
-H 'Referer: https://staging.scyllax.online/' \
-H 'sec-ch-ua: "Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"' \
-H 'sec-ch-ua-mobile: ?0' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36' \
-H 'Accept: application/json, text/plain, */*' \
-H 'Content-Type: application/json' \
--data-raw '{"survey_title":"survey salesman","efective_date_start":"2026-04-01","efective_date_end":"2026-04-18","answer_frequency":"One Time","response_type":"Mandatory","target_type":"Specific","distributor_id":[102],"area_id":[88],"outlet_id":[],"survey_template_id":40,"emp_id":[421]}'
Payload :
{
"survey_title": "survey salesman",
"efective_date_start": "2026-04-01",
"efective_date_end": "2026-04-18",
"answer_frequency": "One Time",
"response_type": "Mandatory",
"target_type": "Specific",
"distributor_id": [102],
"area_id": [88],
"outlet_id": [],
"survey_template_id": 40,
"emp_id": [421]
}
pengaruh ke tabel :
mst.m_survey (DONE)
mst.m_survey_area (TO DO )
insert ke field distributor_id
mst.m_survey_salesman  (TO DO )
cust_id            = berdasarkan cust_id user yang login
m_survey_salesman_id*  = generate BE
survey_id**      = berdasarkan survey_id dari tabel mst.m_survey
salesman_id**    = berdasrkan emp_id dari req body
mst.m_survey_outlet (DONE)
mst.m_survey_detail = (DONE)
