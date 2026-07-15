package entity

type MobileDistributorListQueryFilter struct {
	Query    string `query:"q"`
	Page     int    `query:"page" validate:"required"`
	Limit    int    `query:"limit" validate:"required"`
	Sort     string `query:"sort" validate:"required"`
	IsActive *int   `query:"is_active"`
	RegionID []int  `query:"region_id"`
	AreaID   []int  `query:"area_id"`
}

type MobileDistributorListResponse struct {
	DistributorId   int    `json:"distributor_id"`
	DistributorCode string `json:"distributor_code"`
	DistributorName string `json:"distributor_name"`
	Address         string `json:"address"`
	Latitude        string `json:"latitude"`
	Longitude       string `json:"longitude"`
	RegionID        int    `json:"region_id"`
	AreaID          int    `json:"area_id"`
}
