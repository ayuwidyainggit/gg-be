package entity

type CreateUserBody struct {
	CustId       string `json:"cust_id"`
	UserName     string `json:"user_name"`
	UserPass     string `json:"user_pass"`
	UserFullName string `json:"user_fullname"`
	IsAdmin      bool   `json:"is_admin"`
	Email        string `json:"email" validate:"required,email"`
	LangId       string `json:"lang_id" validate:"required,oneof='id' 'en'"`
	MobileNo     string `json:"mobile_no"`
	Whatsapp     string `json:"whatsapp"`
	UserStatus   int    `json:"user_status"`
	EmpId        int64  `json:"emp_id"`
	EmpStatus    int64  `json:"emp_status"`
	CreatedBy    int64  `json:"created_by"`
	ImageUrl     string `json:"image_url"`
	SkinName     string `json:"skin_name"`
}
type UpdateUserBodyParam struct {
	UserID int64 `params:"user_id" validate:"required"`
}
type DetailUserBodyParam struct {
	UserID int64 `params:"user_id" validate:"required"`
}
type DeleteUserBodyParams struct {
	UserID int64 `params:"user_id" validate:"required"`
}
type UpdateUserBody struct {
	CustId       string  `json:"cust_id"`
	UserName     *string `json:"user_name"`
	UserPass     *string `json:"user_pass"`
	UserFullName *string `json:"user_fullname"`
	IsAdmin      *bool   `json:"is_admin"`
	Email        *string `json:"email" validate:"required,email"`
	LangId       *string `json:"lang_id" validate:"required,oneof='id' 'en'"`
	MobileNo     *string `json:"mobile_no"`
	Whatsapp     *string `json:"whatsapp"`
	UserStatus   *int    `json:"user_status"`
	EmpId        *int64  `json:"emp_id"`
	EmpStatus    *int64  `json:"emp_status"`
	UpdatedBy    int64   `json:"updated_by"`
	ImageUrl     string  `json:"image_url"`
	SkinName     string  `json:"skin_name"`
}

type UserResponse struct {
	UserId       int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	UserFullName string `json:"user_fullname"`
	IsAdmin      bool   `json:"is_admin"`
	Email        string `json:"email"`
	LangId       string `json:"lang_id"`
	MobileNo     string `json:"mobile_no"`
	Whatsapp     string `json:"whatsapp"`
	UserStatus   int    `json:"user_status"`
	EmpId        int64  `json:"emp_id"`
	EmpStatus    int64  `json:"emp_status"`
	CreatedBy    int64  `json:"created_by"`
	UpdatedBy    int64  `json:"updated_by"`
}

type UserListResponse struct {
	UserId       int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	UserFullName string `json:"user_fullname"`
	IsAdmin      bool   `json:"is_admin"`
	Email        string `json:"email"`
	LangId       string `json:"lang_id"`
	MobileNo     string `json:"mobile_no"`
	Whatsapp     string `json:"whatsapp"`
	UserStatus   int    `json:"user_status"`
	EmpId        int64  `json:"emp_id"`
	EmpStatus    int64  `json:"emp_status"`
	CreatedBy    int64  `json:"created_by"`
	UpdatedBy    int64  `json:"updated_by"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email,min=1,max=100"`
}

type ForgotPasswordValidateRequest struct {
	Email   string `json:"email" validate:"required,email,min=1,max=100"`
	OtpCode string `json:"otp_code" validate:"required,len=4"`
}
type UpdatePasswordRequest struct {
	Email           string `json:"email" validate:"required,email,min=1,max=100"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}
type Token struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type LoginResponse struct {
	UserId        *int64 `json:"user_id"`
	UserName      string `json:"user_name"`
	UserFullname  string `json:"user_fullname"`
	IsAdmin       bool   `json:"is_admin"`
	Email         string `json:"email"`
	LangId        string `json:"lang_id"`
	MobileNo      string `json:"mobile_no"`
	Whatsapp      string `json:"whatsapp"`
	CustId        string `json:"cust_id"`
	EmpId         int64  `json:"emp_id"`
	ParentCustId  string `json:"parent_cust_id"`
	Token         *Token `json:"token"`
	ImageUrl      string `json:"image_url"`
	SkinName      string `json:"skin_name"`
	DistributorID int64  `json:"distributor_id"`
}

type UserData struct {
	CustId         *string `json:"cust_id"`
	UserId         *int64  `json:"user_id"`
	Username       *string `json:"user_name"`
	Userpass       *string `json:"user_pass"`
	Fullname       *string `json:"user_fullname"`
	IsAdmin        *bool   `json:"is_admin"`
	Email          *string `json:"email"`
	LangId         *string `json:"lang_id"`
	MobileNo       *string `json:"mobile_no"`
	Whatsapp       *string `json:"whatsapp"`
	UserStatus     *int    `json:"user_status"`
	EmpStatus      *int64  `json:"emp_status"`
	EmpId          *int64  `json:"emp_id"`
	ParentCustId   string  `json:"parent_cust_id"`
	DistPriceGrpId int     `json:"dist_price_grp_id"`
	ImageUrl       string  `json:"image_url"`
	SkinName       string  `json:"skin_name"`
	DistributorID  int64   `json:"distributor_id"`
}

type TokenMetadata struct {
	UserId         int64
	UserName       string
	UserFullName   string
	Email          string
	IsAdmin        bool
	CustId         string
	ParentCustId   string
	DistPriceGrpId int
	EmpId          int64
	LangId         string
	MobileNo       string
	Whatsapp       string
	Expires        int64
}
type MenuBodyParam struct {
	MenuParam string `params:"menu_param" validate:"required"`
}

const (
	MENUWEB     = "web"
	MENUDESKTOP = "desktop"
)

type EmailData struct {
	Email   string `json:"email"`
	OtpCode string `json:"otp_code"`
}
