package entity

type CreateGamificationParticipantsBody struct {
	CustID         string `json:"cust_id"`
	GamificationNo string `json:"gamification_no"`
	SalesmanID     int64  `json:"salesman_id"`
}

type CreateGamificationProductsBody struct {
	CustID         string `json:"cust_id"`
	GamificationNo string `json:"gamification_no"`
	ProID          int64  `json:"pro_id"`
}

type CreateGamificationOutletsBody struct {
	CustID         string `json:"cust_id"`
	GamificationNo string `json:"gamification_no"`
	OutletID       int64  `json:"outlet_id"`
}

type CreateGamificationAnnouncementsBody struct {
	CustID           string `json:"cust_id"`
	GamificationNo   string `json:"gamification_no"`
	AnnouncementDate string `json:"announcement_date"`
}

type GamificationParticipantsResponse struct {
	GamificationParticipantID int64  `json:"gamification_participant_id"`
	SalesmanID                int64  `json:"salesman_id"`
	SalesmanCode              string `json:"salesman_code"`
	SalesmanName              string `json:"salesman_name"`
}

type GamificationProductsResponse struct {
	GamificationProductID int64  `json:"gamification_product_id"`
	ProID                 int64  `json:"pro_id"`
	ProCode               string `json:"pro_code"`
	ProName               string `json:"pro_name"`
}

type GamificationOutletsResponse struct {
	GamificationOutletID int64  `json:"gamification_outlet_id"`
	OutletID             int64  `json:"outlet_id"`
	OutletCode           string `json:"outlet_code"`
	OutletName           string `json:"outlet_name"`
}

type GamificationAnnouncementsResponse struct {
	GamificationAnnouncementID int64  `json:"gamification_announcement_id"`
	AnnouncementDate           string `json:"announcement_date"`
}

type GamificationRankingsResponse struct {
	Ranking      int64   `json:"ranking"`
	SalesmanID   int64   `json:"salesman_id"`
	SalesmanCode string  `json:"salesman_code"`
	SalesmanName string  `json:"salesman_name"`
	Total        float64 `json:"total"`
}
