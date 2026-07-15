package model

type User struct {
	UserId   int64  `db:"user_id"`
	Username string `db:"user_name"`
	Userpass string `db:"user_pass"`
	Fullname string `db:"user_fullname"`
	IsAdmin  bool   `db:"is_admin"`
	Email    string `db:"email"`
	LangId   string `db:"lang_id"`
	MobileNo string `db:"mobile_no"`
	Whatsapp string `db:"whatsapp"`
	Status   int    `db:"user_status"`
	CustId   string `db:"cust_id"`
	EmpId    int64  `db:"emp_id"`
}
