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
	ImageUrl     string  `json:"image_url,omitempty"`
	SkinName     string  `json:"skin_name"`
}

type UpdateUserImageBody struct {
	CustId   string `json:"cust_id"`
	ImageUrl string `json:"image_url"`
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

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email,min=1,max=100"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
	DeviceId        string `json:"device_id" validate:"required"`
	MacAddress      string `json:"mac_address" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email,min=1,max=100"`
	// DeviceId   string `json:"device_id" validate:"required"`
	// MacAddress string `json:"mac_address" validate:"required"`
}

type ForgotPasswordValidateRequest struct {
	RequestID string `json:"request_id" validate:"required"`
	OtpCode   string `json:"otp_code" validate:"required,len=4"`
}

type ForgotPasswordResendRequest struct {
	RequestID string `json:"request_id" validate:"required"`
}

type ResetPasswordRequest struct {
	ResetToken      string `json:"reset_token" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,min=1,max=100"`
	Password string `json:"password" validate:"required"`
	// DeviceId   string `json:"device_id" validate:"required"`
	// MacAddress string `json:"mac_address" validate:"required"`
	FcmToken string `json:"fcm_token" `
}

type LoginResponse struct {
	UserRole             string `json:"user_role"`
	Email                string `json:"email"`
	LangId               string `json:"lang_id"`
	MobileNo             string `json:"mobile_no"`
	Whatsapp             string `json:"whatsapp"`
	CustId               string `json:"cust_id"`
	AccessToken          string `json:"access_token"`
	EmpId                int64  `json:"emp_id"`
	EmpGrpId             int64  `json:"emp_grp_id"`
	EmpCode              string `json:"emp_code"`
	EmpName              string `json:"emp_name"`
	UserId               int64  `json:"user_id"`
	OprTypeOrderTaking   string `json:"opr_type_order_taking"`
	OprTypeCanvas        string `json:"opr_type_canvas"`
	AllowInputPrice      bool   `json:"allow_input_price"`
	TaxOption            string `json:"tax_option"`
	IsActiveGudangUtama  bool   `json:"is_active_gudang_utama"`
	IsActiveGudangCanvas bool   `json:"is_active_gudang_canvas"`
	DistributorId        *int64 `json:"distributor_id"` // NULL if distributor not found
	DistributorCode      string `json:"distributor_code"`
	DistributorName      string `json:"distributor_name"`
	DistributorAddress   string `json:"distributor_address"`
}

type UserData struct {
	UserId               int64   `json:"user_id"`
	UserRole             string  `json:"user_role"`
	Email                string  `json:"email"`
	LangId               string  `json:"lang_id"`
	MobileNo             string  `json:"mobile_no"`
	Whatsapp             string  `json:"whatsapp"`
	CustId               string  `json:"cust_id"`
	EmpId                *int64  `json:"emp_id"`
	EmpGrpId             *int64  `json:"emp_grp_id"`
	EmpCode              *string `json:"emp_code"`
	ParentCustId         string  `json:"parent_cust_id"`
	DistPriceGrpId       int     `json:"dist_price_grp_id"`
	OprTypeOrderTaking   string  `json:"opr_type_order_taking"`
	OprTypeCanvas        string  `json:"opr_type_canvas"`
	AllowInputPrice      bool    `json:"allow_input_price"`
	TaxOption            string  `json:"tax_option"`
	IsActiveGudangCanvas bool    `json:"is_active_gudang_canvas"`
	IsActiveGudangUtama  bool    `json:"is_active_gudang_utama"`
	DistributorID        int     `json:"distributor_id"`
	Username             string  `json:"user_name"`
	UserFullname         string  `json:"user_fullname"`
	IsAdmin              bool    `json:"is_admin"`
}

type TokenMetadata struct {
	UserId               int64
	UserRole             string
	Email                string
	CustId               string
	ParentCustId         string
	DistPriceGrpId       int
	EmpId                int64
	EmpGrpId             int64
	EmpCode              string
	LangId               string
	MobileNo             string
	Whatsapp             string
	Expires              int64
	OprTypeOrderTaking   string
	OprTypeCanvas        string
	AllowInputPrice      bool
	TaxOption            string `json:"tax_option"`
	IsActiveGudangCanvas bool   `json:"is_active_gudang_canvas"`
	IsActiveGudangUtama  bool   `json:"is_active_gudang_utama"`
	DistributorID        int64
	DistributorCode      string
}

type EmailData struct {
	Email   string `json:"email"`
	OtpCode string `json:"otp_code"`
}
type UpdatePasswordRequest struct {
	Email           string `json:"email" validate:"required,email,min=1,max=100"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}
type ChangePasswordRequest struct {
	Email           string `json:"email" validate:"required,email"`
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}
type ProfileResponse struct {
	UserRole        string `json:"user_role"`
	UserImg         string `json:"user_img"`
	SalesCode       string `json:"sales_code"`
	SalesName       string `json:"sales_name"`
	CustID          string `json:"cust_id"`
	Custname        string `json:"cust_name"`
	Email           string `json:"email"`
	SalesTeamName   string `json:"sales_team_name"`
	LangID          string `json:"lang_id"`
	MobileNo        string `json:"mobile_no"`
	Whatsapp        string `json:"whatsapp"`
	IsValidRoute    bool   `json:"is_valid_route"`
	MaxRadius       int64  `json:"max_radius"`
	DistributorCode string `json:"distributor_code"`
	DistributorName string `json:"distributor_name"`
	IsCaptureOutlet bool   `json:"is_captured_outlet"`
	Duration        int64  `json:"duration"`
	Distance        int64  `json:"distance"`
}

type SendLocationRequest struct {
	CustID    string
	EmpID     int64
	Longitude string `json:"longitude" validate:"required"`
	Latitude  string `json:"latitude" validate:"required"`
}
