package entity

import (
	"encoding/json"
	"time"
)

// Helper type for flexible int array unmarshalling
type FlexibleIntArray []int

func (a *FlexibleIntArray) UnmarshalJSON(data []byte) error {
	// Try unmarshalling as array
	var arr []int
	if err := json.Unmarshal(data, &arr); err == nil {
		*a = arr
		return nil
	}

	// Try unmarshalling as single int
	var single int
	if err := json.Unmarshal(data, &single); err == nil {
		*a = []int{single}
		return nil
	}

	return nil // Or return error if strict
}

// ============ Response Structs ============

// SurveyListResponse for List endpoint
type SurveyListResponse struct {
	SurveyId          int        `json:"survey_id"`
	CreatedAt         *time.Time `json:"created_at"`
	AnswerFrequency   string     `json:"answer_frequency"`
	SurveyTitle       string     `json:"survey_title"`
	ResponseType      string     `json:"response_type"`
	EfectiveDateStart *time.Time `json:"efective_date_start"`
	EfectiveDateEnd   *time.Time `json:"efective_date_end"`
	Status            int        `json:"status"`
}

// SurveyDetailResponse for Detail endpoint
type SurveyDetailResponse struct {
	SurveyId           int                       `json:"survey_id"`
	CreatedAt          *time.Time                `json:"created_at"`
	AnswerFrequency    string                    `json:"answer_frequency"`
	SurveyTitle        string                    `json:"survey_title"`
	ResponseType       string                    `json:"response_type"`
	LevelTarget        string                    `json:"level_target"`
	EfectiveDateStart  *time.Time                `json:"efective_date_start"`
	EfectiveDateEnd    *time.Time                `json:"efective_date_end"`
	Status             int                       `json:"status"`
	DistributorId      FlexibleIntArray          `json:"distributor_id"`
	AreaId             FlexibleIntArray          `json:"area_id"`
	DistributorCode    string                    `json:"distributor_code"`
	DistributorName    string                    `json:"distributor_name"`
	BusinessUnits      []SurveyBusinessUnit      `json:"business_units"`
	TargetDistributor  []SurveyDistributorResponse `json:"target_distributor"`
	Outlet             []SurveyOutletResponse    `json:"outlet"`
	Salesman           []SurveySalesmanResponse  `json:"salesman"`
	TargetSurvey       *SurveyTargetResponse     `json:"target_survey"`
	Template           []SurveyTemplateNested    `json:"template"`
}

// SurveyBusinessUnit represents selected distributor-area mapping in detail response.
type SurveyBusinessUnit struct {
	DistributorId   int    `json:"distributor_id"`
	AreaId          int    `json:"area_id"`
	TargetCustId    string `json:"target_cust_id,omitempty"`
	TargetCustName  string `json:"target_cust_name,omitempty"`
	BusinessUnitName string `json:"business_unit_name,omitempty"`
	Name             string `json:"name,omitempty"`
	Type             string `json:"type,omitempty"`
}

// SurveyDistributorResponse represents the target_distributor entries exposed
// on the detail response. Mirrors mst.m_survey_distributor columns plus
// joined distributor metadata.
type SurveyDistributorResponse struct {
	MSurveyDistributorId int    `json:"m_survey_distributor_id"`
	DistributorId        int    `json:"distributor_id"`
	DistributorCode      string `json:"distributor_code"`
	DistributorName      string `json:"distributor_name"`
}

// SurveyTargetResponse for target_survey in detail
type SurveyTargetResponse struct {
	TargetType string                   `json:"target_type"`
	EmpId      *int                     `json:"emp_id"`
	SalesName  *string                  `json:"sales_name"`
	Area       []SurveyAreaResponse     `json:"area"`
	Outlet     []SurveyOutletResponse   `json:"outlet"`
	Salesman   []SurveySalesmanResponse `json:"salesman"`
}

// SurveyAreaResponse for area in target_survey
type SurveyAreaResponse struct {
	AreaId   int    `json:"area_id"`
	AreaName string `json:"area_name"`
}

// SurveyOutletResponse for outlet in target_survey
type SurveyOutletResponse struct {
	SurveyOutletId int    `json:"survey_outlet_id"`
	OutletId       int    `json:"outlet_id"`
	OutletCode     string `json:"outlet_code"`
	OutletName     string `json:"outlet_name"`
	OtClassId      *int   `json:"ot_class_id"`
	OtClassName    string `json:"ot_class_name"`
	OtGrpId        *int   `json:"ot_grp_id"`
	OtGrpName      string `json:"ot_grp_name"`
	OtTypeId       *int   `json:"ot_type_id"`
	OtTypeName     string `json:"ot_type_name"`
}

// SurveySalesmanResponse for salesman in survey detail.
type SurveySalesmanResponse struct {
	MSurveySalesmanId int    `json:"m_survey_salesman_id"`
	SalesId           int    `json:"sales_id"`
	SalesTeamId       *int   `json:"sales_team_id"`
	SalesTeamName     string `json:"sales_team_name"`
	SalesName         string `json:"sales_name"`
}

// SurveyTemplateNested for template in detail (reuses existing QuestionTemplateResponse)
type SurveyTemplateNested struct {
	SurveyTemplateId int                        `json:"survey_template_id"`
	TemplateCode     string                     `json:"template_code"`
	TemplateTitle    string                     `json:"template_title"`
	QuestionTemplate []QuestionTemplateResponse `json:"question_template"`
}

// ============ Request Structs ============

// CreateSurveyBody for POST request
type CreateSurveyBody struct {
	SurveyTitle        string           `json:"survey_title" validate:"required,max=150"`
	EfectiveDateStart  string           `json:"efective_date_start" validate:"required"`
	EfectiveDateEnd    string           `json:"efective_date_end" validate:"required"`
	AnswerFrequency    string           `json:"answer_frequency" validate:"required,answer_frequency"`
	ResponseType       string           `json:"response_type" validate:"required,oneof=Mandatory Optional"`
	LevelTarget        string           `json:"level_target" validate:"omitempty,level_target"`
	TargetType         string           `json:"target_type"`
	TargetCustId       string           `json:"target_cust_id"`
	AreaId             []int            `json:"area_id"`
	DistributorId      FlexibleIntArray `json:"distributor_id"`
	TargetDistributorId FlexibleIntArray `json:"target_distributor_id"`
	OutletId           []int            `json:"outlet_id"`
	SurveyTemplateId   FlexibleIntArray `json:"survey_template_id" validate:"required,min=1"`
	EmpId              FlexibleIntArray `json:"emp_id"`
	CustId             string           `json:"-"`
	ParentCustId       string           `json:"-"`
	CreatedBy          int64            `json:"-"`
}

// UpdateSurveyBody for PUT request
type UpdateSurveyBody struct {
	SurveyTitle        string           `json:"survey_title" validate:"required,max=150"`
	EfectiveDateStart  string           `json:"efective_date_start" validate:"required"`
	EfectiveDateEnd    string           `json:"efective_date_end" validate:"required"`
	AnswerFrequency    string           `json:"answer_frequency" validate:"required,answer_frequency"`
	ResponseType       string           `json:"response_type" validate:"required,oneof=Mandatory Optional"`
	LevelTarget        string           `json:"level_target" validate:"omitempty,level_target"`
	TargetType         string           `json:"target_type"`
	TargetCustId       string           `json:"target_cust_id"`
	AreaId             []int            `json:"area_id"`
	DistributorId      FlexibleIntArray `json:"distributor_id"`
	TargetDistributorId FlexibleIntArray `json:"target_distributor_id"`
	OutletId           []int            `json:"outlet_id"`
	SurveyTemplateId   FlexibleIntArray `json:"survey_template_id" validate:"required,min=1"`
	EmpId              FlexibleIntArray `json:"emp_id"`
	CustId             string           `json:"-"`
	ParentCustId       string           `json:"-"`
	UpdatedBy          int64            `json:"-"`
}

// DeactivateSurveyBody for PATCH request
type DeactivateSurveyBody struct {
	IsActive  bool   `json:"is_active"`
	CustId    string `json:"-"`
	UpdatedBy int64  `json:"-"`
}

// ============ Params & Filter Structs ============

// SurveyParams for path parameter
type SurveyParams struct {
	SurveyId int `params:"survey_id" validate:"required"`
}

// SurveyQueryFilter for List endpoint
type SurveyQueryFilter struct {
	Page              int     `query:"page"`
	Limit             int     `query:"limit" validate:"required"`
	Query             string  `query:"q"`
	Sort              string  `query:"sort"`
	Status            *int    `query:"status"`
	AnswerFrequency   *string `query:"answer_frequency"`
	ResponseFrequency *string `query:"response_frequency"` // Maps to response_type column in DB
	ResponseType      *string `query:"response_type"`      // Alternative param name
}
