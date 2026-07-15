package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type HierarchyApprovalDet struct {
	HierarchyApprovalDetailID  *int64    `gorm:"column:hierarchy_approval_detail_id;default:nextval('sls.hierarchy_approval_detail_seq'::regclass);not null" json:"hierarchy_approval_detail_id"`
	HierarchyApprovalID        int64     `gorm:"column:hierarchy_approval_id" json:"hierarchy_approval_id"`
	Level                      int       `gorm:"column:level" json:"level"`
	HierarchyApprovalDetCustID string    `gorm:"column:hierarchy_approval_detail_cust_id" json:"hierarchy_approval_detail_cust_id"`
	IsActive                   bool      `gorm:"column:is_active" json:"is_active"`
	MaxOverLimit               *float64  `gorm:"column:max_over_limit" json:"max_over_limit"`
	CreatedBy                  *int64    `gorm:"created_by" json:"created_by"`
	CreatedAt                  time.Time `gorm:"created_at" json:"created_at"`
	UpdatedBy                  *int64    `gorm:"updated_by" json:"updated_by"`
	UpdatedAt                  time.Time `gorm:"updated_at" json:"updated_at"`
}

func (HierarchyApprovalDet) TableName() string {
	return "sls.hierarchy_approvals_details"
}

func (m *HierarchyApprovalDet) BeforeCreate(trx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

type HierarchyApprovalDetRead struct {
	HierarchyApprovalDetailID      int64     `gorm:"column:hierarchy_approval_detail_id;default:nextval('sls.hierarchy_approval_detail_seq'::regclass);not null" json:"hierarchy_approval_detail_id"`
	HierarchyApprovalID            int64     `gorm:"column:hierarchy_approval_id" json:"hierarchy_approval_id"`
	Level                          int       `gorm:"column:level" json:"level"`
	HierarchyApprovalDetCustID     string    `gorm:"column:hierarchy_approval_detail_cust_id" json:"hierarchy_approval_detail_cust_id"`
	HierarchyApprovalDetCustIDName string    `gorm:"column:setup_for_name" json:"setup_for_name"`
	IsActive                       bool      `gorm:"column:is_active" json:"is_active"`
	MaxOverLimit                   *float64  `gorm:"column:max_over_limit" json:"max_over_limit"`
	CreatedBy                      *int64    `gorm:"created_by" json:"created_by"`
	CreatedAt                      time.Time `gorm:"created_at" json:"created_at"`
	UpdatedBy                      *int64    `gorm:"updated_by" json:"updated_by"`
	UpdatedAt                      time.Time `gorm:"updated_at" json:"updated_at"`
}

func (HierarchyApprovalDetRead) TableName() string {
	return "sls.hierarchy_approvals_details"
}

type HierarchyApprovalUpdate struct {
	HierarchyApprovalDetCustID string    `gorm:"column:hierarchy_approval_detail_cust_id" json:"hierarchy_approval_detail_cust_id"`
	IsActive                   bool      `gorm:"column:is_active" json:"is_active"`
	MaxOverLimit               *float64  `gorm:"column:max_over_limit" json:"max_over_limit"`
	CreatedBy                  *int64    `gorm:"created_by" json:"created_by"`
	CreatedAt                  time.Time `gorm:"created_at" json:"created_at"`
	UpdatedBy                  *int64    `gorm:"updated_by" json:"updated_by"`
	UpdatedAt                  time.Time `gorm:"updated_at" json:"updated_at"`
}

func (HierarchyApprovalUpdate) TableName() string {
	return "sls.hierarchy_approvals_details"
}

func (m *HierarchyApprovalUpdate) BeforeUpdate(trx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now

	return nil
}

type HierarchyApprovalDetReadMap map[int]*HierarchyApprovalDetRead

func (m HierarchyApprovalDetReadMap) Set(level int, value HierarchyApprovalDetRead) {
	m[level] = &value
}
func (m HierarchyApprovalDetReadMap) GetByLevel(level int) (value *HierarchyApprovalDetRead, err error) {
	val, ok := m[level]

	// If the key exists
	if !ok {
		return value, errors.New(fmt.Sprintf("%v Not Found", level))
	}

	return val, nil
}
