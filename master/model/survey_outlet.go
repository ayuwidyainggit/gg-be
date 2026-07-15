package model

// SurveyOutlet represents mst.m_survey_outlet table
type SurveyOutlet struct {
	SurveyOutletId int  `db:"survey_outlet_id" json:"survey_outlet_id"`
	SurveyId       int  `db:"survey_id" json:"survey_id"`
	OutletId       int  `db:"outlet_id" json:"outlet_id"`
	IsDel          bool `db:"is_del" json:"is_del"`
	// Joined fields
	OutletCode  *string `db:"outlet_code" json:"outlet_code,omitempty"`
	OutletName  *string `db:"outlet_name" json:"outlet_name,omitempty"`
	OtClassId   *int    `db:"ot_class_id" json:"ot_class_id,omitempty"`
	OtClassName *string `db:"ot_class_name" json:"ot_class_name,omitempty"`
	OtGrpId     *int    `db:"ot_grp_id" json:"ot_grp_id,omitempty"`
	OtGrpName   *string `db:"ot_grp_name" json:"ot_grp_name,omitempty"`
	OtTypeId    *int    `db:"ot_type_id" json:"ot_type_id,omitempty"`
	OtTypeName  *string `db:"ot_type_name" json:"ot_type_name,omitempty"`
}
