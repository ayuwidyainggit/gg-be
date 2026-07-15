package pjp

import (
	"scyllax-pjp/data/response"
	"testing"
	"time"
)

func TestBuildDestinationDetailsDataSplitsOutletsAndDistributors(t *testing.T) {
	date := time.Date(2026, 6, 29, 0, 0, 0, 0, time.UTC)
	rows := []response.DestinationDetailRow{
		{
			RouteCode: 3883, RouteName: "Route 1", Week: 84, Year: 2026, Date: date,
			DestinationID: 1722, DestinationCode: "BMI260003", DestinationType: "outlet",
			DestinationName: "Toko merah", Longitude: "1", Latitude: "2",
			DestinationStatus: "1", DestinationAddress: "Jalan Bangka 1",
		},
		{
			RouteCode: 3883, RouteName: "Route 1", Week: 84, Year: 2026, Date: date,
			DestinationID: 102, DestinationCode: "3100022", DestinationType: "distributor",
			DestinationName: "PT Besi Makmur", Longitude: "106.8456", Latitude: "-6.2088",
			DestinationStatus: "", DestinationAddress: "jalan besi makmur madura",
		},
	}

	data := buildDestinationDetailsData(rows)

	if data.RouteCode != 3883 || data.RouteName != "Route 1" {
		t.Fatalf("route = %d/%q, want 3883/Route 1", data.RouteCode, data.RouteName)
	}
	if len(data.Outlets) != 1 || data.Outlets[0].OutletID != 1722 {
		t.Fatalf("outlets = %+v, want outlet 1722", data.Outlets)
	}
	if len(data.Distributors) != 1 || data.Distributors[0].DistributorID != 102 {
		t.Fatalf("distributors = %+v, want distributor 102", data.Distributors)
	}
}
