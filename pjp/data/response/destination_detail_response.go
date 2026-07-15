package response

import "time"

type DestinationDetailsResponse struct {
	Message   string                 `json:"message"`
	Data      DestinationDetailsData `json:"data"`
	Paging    LiveMonitoringPaging   `json:"paging"`
	RequestID string                 `json:"request_id"`
}

type DestinationDetailsData struct {
	RouteCode    int                            `json:"route_code"`
	RouteName    string                         `json:"route_name"`
	Week         int                            `json:"week"`
	Year         int                            `json:"year"`
	Date         time.Time                      `json:"date"`
	Outlets      []DestinationDetailOutlet      `json:"outlets"`
	Distributors []DestinationDetailDistributor `json:"distributors"`
}

type DestinationDetailOutlet struct {
	OutletID      int    `json:"outlet_id"`
	OutletCode    string `json:"outlet_code"`
	OutletName    string `json:"outlet_name"`
	Longitude     string `json:"longitude"`
	Latitude      string `json:"latitude"`
	OutletStatus  string `json:"outlet_status"`
	OutletAddress string `json:"outlet_address"`
}

type DestinationDetailDistributor struct {
	DistributorID      int    `json:"distributor_id"`
	DistributorCode    string `json:"distributor_code"`
	DistributorName    string `json:"distributor_name"`
	DistributorStatus  string `json:"distributor_status"`
	DistributorAddress string `json:"distributor_address"`
	Longitude          string `json:"longitude"`
	Latitude           string `json:"latitude"`
}

type DestinationDetailRow struct {
	RouteCode          int       `gorm:"column:route_code"`
	RouteName          string    `gorm:"column:route_name"`
	Week               int       `gorm:"column:week"`
	Year               int       `gorm:"column:year"`
	Date               time.Time `gorm:"column:date"`
	DestinationID      int       `gorm:"column:destination_id"`
	DestinationCode    string    `gorm:"column:destination_code"`
	DestinationType    string    `gorm:"column:destination_type"`
	DestinationName    string    `gorm:"column:destination_name"`
	Longitude          string    `gorm:"column:longitude"`
	Latitude           string    `gorm:"column:latitude"`
	DestinationStatus  string    `gorm:"column:destination_status"`
	DestinationAddress string    `gorm:"column:destination_address"`
}
