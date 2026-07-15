package model

// SurveySalesman represents mst.m_survey_salesman table.
type SurveySalesman struct {
	MSurveySalesmanId int    `db:"m_survey_salesman_id" json:"m_survey_salesman_id"`
	CustId            string `db:"cust_id" json:"cust_id"`
	SurveyId          int    `db:"survey_id" json:"survey_id"`
	SalesmanId        int    `db:"salesman_id" json:"salesman_id"`
	IsDel             bool   `db:"is_del" json:"is_del"`
	// Joined fields
	SalesTeamId   *int    `db:"sales_team_id" json:"sales_team_id,omitempty"`
	SalesTeamName *string `db:"sales_team_name" json:"sales_team_name,omitempty"`
	SalesName     *string `db:"sales_name" json:"sales_name,omitempty"`
}
