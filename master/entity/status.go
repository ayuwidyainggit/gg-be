package entity

type StatusQueryFilter struct {
	Sort     string `query:"sort"`
	Mode     string `query:"mode"`
	Page     int    `query:"page"`
	Limit    int    `query:"limit" validate:"required"`
	Query    string `query:"q"`
	StatusId string `query:"status_id"`
	LangId   string `query:"lang_id"`
}

type StatusResponse struct {
	StatusID    string `json:"status_id"`
	StatusName  string `json:"status_name"`
	StatusValue int    `json:"status_value"`
	LangID      string `json:"lang_id"`
}

type StatusListResponse struct {
	StatusID    string `json:"status_id"`
	StatusName  string `json:"status_name"`
	StatusValue int    `json:"status_value"`
	LangID      string `json:"lang_id"`
}

type StatusLookupResponse struct {
	StatusID    string `json:"status_id"`
	StatusName  string `json:"status_name"`
	StatusValue int    `json:"status_value"`
	LangID      string `json:"lang_id"`
}

type CreateStatusBody struct {
	StatusID    string `json:"status_id"`
	StatusName  string `json:"status_name"`
	StatusValue int    `json:"status_value"`
	LangID      string `json:"lang_id"`
}

type DetailStatusParams struct {
	StatusId    string `params:"status_id" validate:"required"`
	StatusValue int    `params:"status_value" validate:"required"`
}

type UpdateStatusParams struct {
	StatusId string `params:"status_id" validate:"required"`
}

type DeleteStatusParams struct {
	StatusId string `params:"status_id" validate:"required"`
}

type UpdateStatusRequest struct {
	CustId      string `json:"cust_id" validate:"required,max=10"`
	StatusName  string `json:"status_name"`
	StatusValue int    `json:"status_value"`
	LangID      string `json:"lang_id"`
}
