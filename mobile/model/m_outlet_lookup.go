package model

type OutletLookupList struct {
	OutletID      int64   `db:"outlet_id" json:"outlet_id"`
	OutletCode    string  `db:"outlet_code" json:"outlet_code"`
	OutletName    string  `db:"outlet_name" json:"outlet_name"`
	Address       string  `db:"address" json:"address"`
	Latitude      *string `db:"latitude" json:"latitude"`
	Longitude     *string `db:"longitude" json:"longitude"`
	DistributorID *int    `db:"distributor_id" json:"distributor_id"`
	RegionID      *int    `db:"region_id" json:"region_id"`
	AreaID        *int    `db:"area_id" json:"area_id"`
	OutletStatus  *int    `db:"outlet_status" json:"outlet_status"`
}
