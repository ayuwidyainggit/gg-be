package model

// SurveyArea represents mst.m_survey_area table
type SurveyArea struct {
	SurveyAreaId  int    `db:"survey_area_id" json:"survey_area_id"`
	SurveyId      int    `db:"survey_id" json:"survey_id"`
	DistributorId int    `db:"distributor_id" json:"distributor_id"`
	AreaId        int    `db:"area_id" json:"area_id"`
	TargetCustId  string `db:"target_cust_id" json:"target_cust_id"`
	IsDel         bool   `db:"is_del" json:"is_del"`
	// Joined fields
	AreaName        *string `db:"area_name" json:"area_name,omitempty"`
	DistributorName *string `db:"distributor_name" json:"distributor_name,omitempty"`
	CustName        *string `db:"cust_name" json:"cust_name,omitempty"`
}
