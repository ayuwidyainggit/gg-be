package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	CustId     string         `gorm:"column:cust_id" json:"cust_id"`
	UserId     *int64         `gorm:"column:user_id;primaryKey" json:"user_id"`
	Username   *string        `gorm:"column:user_name" json:"user_name"`
	Userpass   *string        `gorm:"column:user_pass" json:"user_pass"`
	Fullname   *string        `gorm:"column:user_fullname" json:"user_fullname"`
	IsAdmin    *bool          `gorm:"column:is_admin" json:"is_admin"`
	Email      *string        `gorm:"column:email" json:"email"`
	LangId     *string        `gorm:"column:lang_id" json:"lang_id"`
	MobileNo   *string        `gorm:"column:mobile_no" json:"mobile_no"`
	Whatsapp   *string        `gorm:"column:whatsapp" json:"whatsapp"`
	UserStatus *int           `gorm:"column:user_status" json:"user_status"`
	EmpStatus  *int64         `gorm:"column:emp_status" json:"emp_status"`
	EmpId      *int64         `gorm:"column:emp_id" json:"emp_id"`
	CreatedBy  *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel      bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	ImageUrl   *string        `gorm:"column:image_url" json:"image_url"`
	SkinName   *string        `gorm:"column:skin_name" json:"skin_name"`
	FcmToken   *string        `gorm:"column:fcm_token" json:"fcm_token"`
}

func (m *User) BeforeCreate(trx *gorm.DB) (err error) {
	password := []byte(*m.Userpass)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	tempHashPass := string(hashedPassword)
	m.Userpass = &tempHashPass
	m.CreatedAt = time.Now()

	return nil
}

func (m *User) BeforeUpdate(trx *gorm.DB) (err error) {
	if m.Userpass != nil {
		tempHashPass, err := HashPasswordString(*m.Userpass)
		if err != nil {
			return err
		}
		m.Userpass = &tempHashPass
		now := time.Now()
		m.UpdatedAt = &now
	}
	return
}

func HashPasswordString(pass string) (string, error) {
	password := []byte(pass)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (User) TableName() string {
	return "sys.m_user"
}

type UserRead struct {
	CustId     string         `gorm:"column:cust_id" json:"cust_id"`
	UserId     *int64         `gorm:"column:user_id;primaryKey" json:"user_id"`
	Username   *string        `gorm:"column:user_name" json:"user_name"`
	Userpass   *string        `gorm:"column:user_pass" json:"user_pass"`
	Fullname   *string        `gorm:"column:user_fullname" json:"user_fullname"`
	IsAdmin    *bool          `gorm:"column:is_admin" json:"is_admin"`
	Email      *string        `gorm:"column:email" json:"email"`
	LangId     *string        `gorm:"column:lang_id" json:"lang_id"`
	MobileNo   *string        `gorm:"column:mobile_no" json:"mobile_no"`
	Whatsapp   *string        `gorm:"column:whatsapp" json:"whatsapp"`
	UserStatus *int           `gorm:"column:user_status" json:"user_status"`
	EmpStatus  *int64         `gorm:"column:emp_status" json:"emp_status"`
	EmpId      *int64         `gorm:"column:emp_id" json:"emp_id"`
	CreatedBy  *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy  *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt  *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel      bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy  *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	ImageUrl   *string        `gorm:"column:image_url" json:"image_url"`
	SkinName   *string        `gorm:"column:skin_name" json:"skin_name"`
	RoleName   *string        `gorm:"column:role_name" json:"role_name"`
}

func (UserRead) TableName() string {
	return "sys.m_user"
}
