package request

type SalesmanListQueryFilter struct {
	Page            string `query:"page"`
	Limit           string `query:"limit"`
	Query           string `query:"q"`
	DestinationCode string `query:"outlet_code"`
	DestinationID   int    `query:"outlet_id"`
	Mode            string `query:"mode"`
	Sort            string `query:"sort"`
	IsActive        string `query:"is_active"`
	SalesTeamID     string `query:"sales_team_id"`
}
