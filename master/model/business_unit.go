package model

// BusinessUnitDistributor for query result
type BusinessUnitDistributor struct {
	CustId          string  `db:"cust_id"`
	DistributorId   int     `db:"distributor_id"`
	DistributorCode string  `db:"distributor_code"`
	DistributorName string  `db:"distributor_name"`
	AreaId          *int    `db:"area_id"`
	AreaCode        *string `db:"area_code"`
	AreaName        *string `db:"area_name"`
	RegionId        *int    `db:"region_id"`
	RegionCode      *string `db:"region_code"`
	RegionName      *string `db:"region_name"`
}

// UserInfo for sys.m_user query
type UserInfo struct {
	UserId       int    `db:"user_id"`
	UserFullname string `db:"user_fullname"`
}
