package model

import (
	"time"

	"gorm.io/gorm"
)

type HierarchyApproval struct {
	HierarchyApprovalID   *int64          `gorm:"column:hierarchy_approval_id;default:nextval('sls.hierarchy_approval_seq'::regclass);not null" json:"hierarchy_approval_id"`
	SetupFor              string          `gorm:"setup_for" json:"setup_for"`
	HierarchyApprovalType int             `gorm:"hierarchy_approval_type" json:"hierarchy_approval_type"`
	CreatedBy             *int64          `gorm:"created_by" json:"created_by"`
	CreatedAt             time.Time       `gorm:"created_at" json:"created_at"`
	UpdatedBy             *int64          `gorm:"updated_by" json:"updated_by"`
	UpdatedAt             time.Time       `gorm:"updated_at" json:"updated_at"`
	DeletedBy             *int64          `gorm:"deleted_by" json:"deleted_by"`
	DeletedAt             *gorm.DeletedAt `gorm:"deleted_at" json:"deleted_at"`
}

func (HierarchyApproval) TableName() string {
	return "sls.hierarchy_approvals"
}

func (m *HierarchyApproval) BeforeCreate(trx *gorm.DB) (err error) {
	// intTmpsStr := time.Now().UnixNano() / int64(time.Millisecond)
	// m.RoNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.UpdatedBy = m.CreatedBy
	return nil
}

func (m *HierarchyApproval) BeforeUpdate(trx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now

	return nil
}

type HierarchyApprovalList struct {
	HierarchyApprovalID   int64           `gorm:"column:hierarchy_approval_id" json:"hierarchy_approval_id"`
	SetupFor              string          `gorm:"column:setup_for" json:"setup_for"`
	SetupForName          string          `gorm:"column:setup_for_name" json:"setup_for_name"`
	CompanyName           string          `gorm:"column:company_name" json:"company_name"`
	CompanyCode           string          `gorm:"column:company_code" json:"company_code"`
	HierarchyApprovalType int             `gorm:"column:hierarchy_approval_type" json:"hierarchy_approval_type"`
	CreatedBy             *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt             time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedBy             *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedByname         string          `gorm:"column:updated_by_name" json:"updated_by_name"`
	UpdatedAt             time.Time       `gorm:"column:updated_at" json:"updated_at"`
	DeletedBy             *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt             *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (HierarchyApprovalList) TableName() string {
	return "sls.hierarchy_approvals"
}

type HierarchyApprovalRead struct {
	HierarchyApprovalID   int64           `gorm:"column:hierarchy_approval_id" json:"hierarchy_approval_id"`
	SetupFor              string          `gorm:"column:setup_for" json:"setup_for"`
	SetupForName          string          `gorm:"column:setup_for_name" json:"setup_for_name"`
	CompanyName           string          `gorm:"column:company_name" json:"company_name"`
	CompanyCode           string          `gorm:"column:company_code" json:"company_code"`
	HierarchyApprovalType int             `gorm:"column:hierarchy_approval_type" json:"hierarchy_approval_type"`
	CreatedBy             *int64          `gorm:"column:created_by" json:"created_by"`
	CreatedAt             time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedBy             *int64          `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt             time.Time       `gorm:"column:updated_at" json:"updated_at"`
	DeletedBy             *int64          `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt             *gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (HierarchyApprovalRead) TableName() string {
	return "sls.hierarchy_approvals"
}
