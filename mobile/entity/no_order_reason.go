package entity

type NoOrderReasonResp struct {
	TakingOrderId int64  `json:"taking_order_id"`
	Reason        string `json:"reason"`
	ImageUrl      string `json:"image_url"`
}
