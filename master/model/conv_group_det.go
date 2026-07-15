package model

import "time"

type ConvGroupDet struct {
	CustId        string     `json:"cust_id" db:"cust_id"`
	ConvGrpId     int        `json:"conv_grp_id" db:"conv_grp_id"`
	ProId         int        `json:"pro_id" db:"pro_id"`
	UnitId1       string     `json:"unit_id1" db:"unit_id1"`
	UnitId2       string     `json:"unit_id2" db:"unit_id2"`
	UnitId3       string     `json:"unit_id3" db:"unit_id3"`
	UnitId4       *string    `json:"unit_id4" db:"unit_id4"`
	UnitId5       *string    `json:"unit_id5" db:"unit_id5"`
	ConvUnit2     float64    `json:"conv_unit2" db:"conv_unit2"`
	ConvUnit3     float64    `json:"conv_unit3" db:"conv_unit3"`
	ConvUnit4     float64    `json:"conv_unit4" db:"conv_unit4"`
	ConvUnit5     float64    `json:"conv_unit5" db:"conv_unit5"`
	NewUnitId1    string     `json:"new_unit_id1" db:"new_unit_id1"`
	NewUnitId2    string     `json:"new_unit_id2" db:"new_unit_id2"`
	NewUnitId3    string     `json:"new_unit_id3" db:"new_unit_id3"`
	NewUnitId4    string     `json:"new_unit_id4" db:"new_unit_id4"`
	NewUnitId5    string     `json:"new_unit_id5" db:"new_unit_id5"`
	NewConvUnit2  float64    `json:"new_conv_unit2" db:"new_conv_unit2"`
	NewConvUnit3  float64    `json:"new_conv_unit3" db:"new_conv_unit3"`
	NewConvUnit4  float64    `json:"new_conv_unit4" db:"new_conv_unit4"`
	NewConvUnit5  float64    `json:"new_conv_unit5" db:"new_conv_unit5"`
	CreatedBy     *int64     `json:"created_by" db:"created_by,omitempty"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at,omitempty"`
	UpdatedBy     *int64     `json:"updated_by" db:"updated_by,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at,omitempty"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	IsDel         bool       `json:"is_del" db:"is_del"`
	DeletedBy     *int64     `json:"deleted_by" db:"deleted_by,omitempty"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at,omitempty"`
}

type ConvGroupDetUpdate struct {
	ConvGrpId    *int       `json:"conv_grp_id" sql:"conv_grp_id"`
	ProId        *int       `json:"pro_id" sql:"pro_id"`
	UnitIdS      *string    `json:"unit_id_s" sql:"unit_id_s"`
	UnitIdM      *string    `json:"unit_id_m" sql:"unit_id_m"`
	UnitIdL      *string    `json:"unit_id_l" sql:"unit_id_l"`
	ConvUnitM    *float64   `json:"conv_unit_m" sql:"conv_unit_m"`
	ConvUnitL    *float64   `json:"conv_unit_l" sql:"conv_unit_l"`
	NewUnitIdS   *string    `json:"new_unit_id_s" sql:"new_unit_id_s"`
	NewUnitIdM   *string    `json:"new_unit_id_m" sql:"new_unit_id_m"`
	NewUnitIdL   *string    `json:"new_unit_id_l" sql:"new_unit_id_l"`
	NewConvUnitM *float64   `json:"new_conv_unit_m" sql:"new_conv_unit_m"`
	NewConvUnitL *float64   `json:"new_conv_unit_l" sql:"new_conv_unit_l"`
	UpdatedAt    *time.Time `json:"updated_at" sql:"updated_at"`
	UpdatedBy    *int       `json:"updated_by" sql:"updated_by"`
}
type ConvGroupDetRead struct {
	CustId        string     `json:"cust_id" db:"cust_id"`
	ConvGrpId     int        `json:"conv_grp_id" db:"conv_grp_id"`
	ProId         int        `json:"pro_id" db:"pro_id"`
	ProCode       *string    `json:"pro_code" db:"pro_code"`
	ProName       *string    `json:"pro_name" db:"pro_name"`
	UnitId1       *string    `json:"unit_id1" db:"unit_id1"`
	UnitId2       *string    `json:"unit_id2" db:"unit_id2"`
	UnitId3       *string    `json:"unit_id3" db:"unit_id3"`
	UnitId4       *string    `json:"unit_id4" db:"unit_id4"`
	UnitId5       *string    `json:"unit_id5" db:"unit_id5"`
	ConvUnit2     *float64   `json:"conv_unit2" db:"conv_unit2"`
	ConvUnit3     *float64   `json:"conv_unit3" db:"conv_unit3"`
	ConvUnit4     *float64   `json:"conv_unit4" db:"conv_unit4"`
	ConvUnit5     *float64   `json:"conv_unit5" db:"conv_unit5"`
	NewUnitId1    *string    `json:"new_unit_id1" db:"new_unit_id1"`
	NewUnitId2    *string    `json:"new_unit_id2" db:"new_unit_id2"`
	NewUnitId3    *string    `json:"new_unit_id3" db:"new_unit_id3"`
	NewUnitId4    *string    `json:"new_unit_id4" db:"new_unit_id4"`
	NewUnitId5    *string    `json:"new_unit_id5" db:"new_unit_id5"`
	NewConvUnit2  *float64   `json:"new_conv_unit2" db:"new_conv_unit2"`
	NewConvUnit3  *float64   `json:"new_conv_unit3" db:"new_conv_unit3"`
	NewConvUnit4  *float64   `json:"new_conv_unit4" db:"new_conv_unit4"`
	NewConvUnit5  *float64   `json:"new_conv_unit5" db:"new_conv_unit5"`
	CreatedBy     *int64     `json:"created_by" db:"created_by,omitempty"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at,omitempty"`
	UpdatedBy     *int64     `json:"updated_by" db:"updated_by,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at,omitempty"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	IsDel         bool       `json:"is_del" db:"is_del"`
	DeletedBy     *int64     `json:"deleted_by" db:"deleted_by,omitempty"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at,omitempty"`
}
