package entity

import "time"

type ConvGroupDetResponse struct {
	ConvGrpId     int        `json:"conv_grp_id"`
	ProId         int        `json:"pro_id"`
	ProCode       string     `json:"pro_code" db:"pro_code"`
	ProName       string     `json:"pro_name" db:"pro_name"`
	UnitId1       string     `json:"unit_id1"`
	UnitId2       string     `json:"unit_id2"`
	UnitId3       string     `json:"unit_id3"`
	UnitId4       string     `json:"unit_id4"`
	UnitId5       string     `json:"unit_id5"`
	ConvUnit2     float64    `json:"conv_unit2"`
	ConvUnit3     float64    `json:"conv_unit3"`
	ConvUnit4     float64    `json:"conv_unit4"`
	ConvUnit5     float64    `json:"conv_unit5"`
	NewUnitId1    string     `json:"new_unit_id1"`
	NewUnitId2    string     `json:"new_unit_id2"`
	NewUnitId3    string     `json:"new_unit_id3"`
	NewUnitId4    string     `json:"new_unit_id4"`
	NewUnitId5    string     `json:"new_unit_id5"`
	NewConvUnit2  float64    `json:"new_conv_unit2"`
	NewConvUnit3  float64    `json:"new_conv_unit3"`
	NewConvUnit4  float64    `json:"new_conv_unit4"`
	NewConvUnit5  float64    `json:"new_conv_unit5"`
	UpdatedBy     *int       `json:"updated_by"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type CreateConvGroupDetBody struct {
	CustId       string  `json:"cust_id" validate:"required,max=10"`
	ConvGrpId    int     `json:"conv_grp_id"`
	ProId        int     `json:"pro_id"`
	UnitId1      string  `json:"unit_id1"`
	UnitId2      string  `json:"unit_id2"`
	UnitId3      string  `json:"unit_id3"`
	UnitId4      string  `json:"unit_id4"`
	UnitId5      string  `json:"unit_id5"`
	ConvUnit2    float64 `json:"conv_unit2"`
	ConvUnit3    float64 `json:"conv_unit3"`
	ConvUnit4    float64 `json:"conv_unit4"`
	ConvUnit5    float64 `json:"conv_unit5"`
	NewUnitId1   string  `json:"new_unit_id1"`
	NewUnitId2   string  `json:"new_unit_id2"`
	NewUnitId3   string  `json:"new_unit_id3"`
	NewUnitId4   string  `json:"new_unit_id4"`
	NewUnitId5   string  `json:"new_unit_id5"`
	NewConvUnit2 float64 `json:"new_conv_unit2"`
	NewConvUnit3 float64 `json:"new_conv_unit3"`
	NewConvUnit4 float64 `json:"new_conv_unit4"`
	NewConvUnit5 float64 `json:"new_conv_unit5"`
	CreatedBy    int64   `json:"created_by"`
	UpdatedBy    int64   `json:"updated_by"`
}

type DetailConvGroupDetParams struct {
	ConvGrpId int `params:"conv_grp_id" validate:"required"`
	ProId     int `params:"pro_id" validate:"required"`
}

type UpdateConvGroupDetParams struct {
	ConvGrpId int `params:"conv_grp_id" validate:"required"`
	ProId     int `params:"pro_id" validate:"required"`
}

type DeleteConvGroupDetParams struct {
	ConvGrpId int `params:"conv_grp_id" validate:"required"`
	ProId     int `params:"pro_id" validate:"required"`
}

type UpdateConvGroupDetRequest struct {
	CustId       string  `json:"cust_id" validate:"required,max=10"`
	ConvGrpId    int     `json:"conv_grp_id"`
	ProId        int     `json:"pro_id"`
	UnitIdS      string  `json:"unit_id_s"`
	UnitIdM      string  `json:"unit_id_m"`
	UnitIdL      string  `json:"unit_id_l"`
	ConvUnitM    float64 `json:"conv_unit_m"`
	ConvUnitL    float64 `json:"conv_unit_l"`
	NewUnitIdS   string  `json:"new_unit_id_s"`
	NewUnitIdM   string  `json:"new_unit_id_m"`
	NewUnitIdL   string  `json:"new_unit_id_l"`
	NewConvUnitM float64 `json:"new_conv_unit_m"`
	NewConvUnitL float64 `json:"new_conv_unit_l"`
	UpdatedBy    int64   `json:"updated_by" validate:"required"`
}
