package model

import (
	"time"

	"gorm.io/gorm"
)

type MCoaType struct {
	CoaTypeID   int64          `gorm:"column:coa_type_id;primaryKey" json:"coa_type_id"`
	CoaTypeName *string        `gorm:"column:coa_type_name" json:"coa_type_name"`
	CoaGroup    *string        `gorm:"column:coa_group" json:"coa_group"`
	DefBlc      *string        `gorm:"column:def_blc" json:"def_blc"`
	SortIndex   *int64         `gorm:"column:sort_index" json:"sort_index"`
	CoaKind     *int64         `gorm:"column:coa_kind" json:"coa_kind"`
	IsActive    bool           `gorm:"column:is_active" json:"is_active"`
	CreatedBy   *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy   *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt   *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	IsDel       bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy   *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (m *MCoaType) BeforeCreate(trx *gorm.DB) (err error) {
	now := time.Now()

	// intTmpsStr := now.UnixNano() / int64(time.Millisecond)
	// m.ChqTrNo = strconv.Itoa(int(intTmpsStr))
	m.CreatedAt = now
	m.UpdatedBy = m.CreatedBy
	return nil
}

func (MCoaType) TableName() string {
	return "acf.m_coa_type"
}

type MCoaTypeList struct {
	CoaTypeID     int64          `gorm:"column:coa_type_id;primaryKey" json:"coa_type_id"`
	CoaTypeName   *string        `gorm:"column:coa_type_name" json:"coa_type_name"`
	CoaGroup      *string        `gorm:"column:coa_group" json:"coa_group"`
	DefBlc        *string        `gorm:"column:def_blc" json:"def_blc"`
	SortIndex     *int64         `gorm:"column:sort_index" json:"sort_index"`
	CoaKind       *int64         `gorm:"column:coa_kind" json:"coa_kind"`
	IsActive      *bool          `gorm:"column:is_active" json:"is_active"`
	CreatedBy     *int64         `gorm:"column:created_by" json:"created_by"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedBy     *int64         `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt     *time.Time     `gorm:"column:updated_at" json:"updated_at"`
	UpdatedByName *string        `gorm:"column:updated_by_name" json:"updated_by_name"`
	IsDel         bool           `gorm:"column:is_del" json:"is_del"`
	DeletedBy     *int64         `gorm:"column:deleted_by" json:"deleted_by"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MCoaTypeList) TableName() string {
	return "acf.m_coa_type"
}
