package request

// TODO Request for process auto mapping by distance and revenue
type PjpProcesDistanceWithRevenue struct {
	Day     []string         `validate:"required" json:"Day"`
	Route   [][]Outlets      `validate:"required" json:"Route"`
	Sales   []PjpAutoRequest `validate:"required" json:"Sales"`
	Revenue []int64          `validate:"required" json:"Revenue"`
}

type CreatePjpAuto struct {
	Day   []string         `validate:"required" json:"Day" example:"[\"Monday\",\"Monday\",\"Sunday\",\"Sunday\"]"`
	Route [][]OutletReq    `validate:"required" json:"Route"`
	Sales []PjpAutoRequest `validate:"required" json:"Sales"`
}

// type CreatePjpAuto struct {
// 	Day   []string         `validate:"required" json:"Day" example:"[\"Monday\",\"Monday\",\"Sunday\",\"Sunday\"]"`
// 	Route [][]OutletReq    `validate:"required" json:"Route" example:"[[{\"latitude\":\"29.355400\",\"longitude\":\"-98.445857\",\"outlet_address\":\"Semarang\",\"outlet_code\":\"OUT10\",\"outlet_id\":10,\"outlet_name\":\"Outlet 10\",\"outlet_status\":1},{\"latitude\":\"29.351701\",\"longitude\":\"-98.514740\",\"outlet_address\":\"Magelang\",\"outlet_code\":\"OUT76\",\"outlet_id\":76,\"outlet_name\":\"Outlet 76\",\"outlet_status\":1}],[{\"latitude\":\"29.417867\",\"longitude\":\"-98.680534\",\"outlet_address\":\"Demak\",\"outlet_code\":\"OUT6\",\"outlet_id\":6,\"outlet_name\":\"Outlet 6\",\"outlet_status\":1}],[{\"latitude\":\"29.486833\",\"longitude\":\"-98.508355\",\"outlet_address\":\"Semarang\",\"outlet_code\":\"OUT99\",\"outlet_id\":99,\"outlet_name\":\"Outlet 99\",\"outlet_status\":1},{\"latitude\":\"29.468601\",\"longitude\":\"-98.524849\",\"outlet_address\":\"Jogja\",\"outlet_code\":\"OUT89\",\"outlet_id\":89,\"outlet_name\":\"Outlet 89\",\"outlet_status\":1},{\"latitude\":\"29.435115\",\"longitude\":\"-98.593962\",\"outlet_address\":\"Solo\",\"outlet_code\":\"OUT5\",\"outlet_id\":5,\"outlet_name\":\"Outlet 5\",\"outlet_status\":1},{\"latitude\":\"29.394394\",\"longitude\":\"-98.530070\",\"outlet_address\":\"Sleman\",\"outlet_code\":\"OUT2\",\"outlet_id\":2,\"outlet_name\":\"Outlet 2\",\"outlet_status\":1}],[{\"latitude\":\"29.459497\",\"longitude\":\"-98.434057\",\"outlet_address\":\"Bantul\",\"outlet_code\":\"OUT87\",\"outlet_id\":87,\"outlet_name\":\"Outlet 87\",\"outlet_status\":1},{\"latitude\":\"29.417361\",\"longitude\":\"-98.437544\",\"outlet_address\":\"Temanggung\",\"outlet_code\":\"OUT94\",\"outlet_id\":94,\"outlet_name\":\"Outlet 94\",\"outlet_status\":1}]]]"`
// 	Sales []PjpAutoRequest `validate:"required" json:"Sales" example:"[{\"id\":18,\"pjp_code\":1000,\"salesman_name\":\"DEN BEI\"},{\"id\":19,\"pjp_code\":4000,\"salesman_name\":\"JACK DANIEL\"},{\"id\":18,\"pjp_code\":1000,\"salesman_name\":\"DEN BEI\"},{\"id\":19,\"pjp_code\":4000,\"salesman_name\":\"JACK DANIEL\"}]"`
// }

type PjpAutoRequest struct {
	ID           int    `validate:"required" json:"id" example:"18"`
	PjpCode      string `validate:"required" json:"pjp_code" example:"1000"`
	SalesmanName string `validate:"required" json:"salesman_name" example:"DEN BEI"`
}

type Pjp struct {
	ID            int    `json:"id"`
	PjpCode       int    `json:"pjp_code"`
	OperationType string `json:"operation_type"`
	TeamSalesMan  string `json:"team_salesman"`
	SalesManID    int    `json:"salesman_id"`
	SalesmanName  string `json:"salesman_name"`
	SalesmanCode  string `json:"salesman_code"`
	PjpMode       string `json:"pjp_mode"`
	Status        bool   `json:"status"`
}

type PjpProces struct {
	ID      int `json:"id"`
	PjpCode int `json:"pjp_code"`
	// SalesManID   int    `json:"salesman_id"`
	SalesmanName string `json:"salesman_name"`
	// SalesmanCode string `json:"salesman_code"`
}

type OutletReq struct {
	Latitude           string `json:"latitude" example:"29.355400"`
	Longitude          string `json:"longitude" example:"-98.445857"`
	DestinationAddress string `json:"outlet_address" example:"Semarang"`
	DestinationCode    string `json:"outlet_code" example:"OUT10"`
	DestinationID      int    `json:"outlet_id" example:"10"`
	DestinationName    string `json:"outlet_name" example:"Outlet 10"`
	DestinationStatus  int    `json:"outlet_status" example:"1"`
}

type CreatePjpAutoExampleReq struct {
	Day []string `json:"Day" example:"[\"Monday\",\"Monday\",\"Sunday\",\"Sunday\"]"`
	// Day   []string         `json:"Day" example: ["1000"]`
	Route [][]Outlets      `json:"Route" swaggertype:"array,object" example:"[[{\"latitude\":\"29.355400\",\"longitude\":\"-98.445857\",\"outlet_address\":\"Semarang\",\"outlet_code\":\"OUT10\",\"outlet_id\":10,\"outlet_name\":\"Outlet 10\",\"outlet_status\":1}],[{\"latitude\":\"29.417867\",\"longitude\":\"-98.680534\",\"outlet_address\":\"Demak\",\"outlet_code\":\"OUT6\",\"outlet_id\":6,\"outlet_name\":\"Outlet 6\",\"outlet_status\":1}]]"`
	Sales []PjpAutoRequest `json:"Sales" example:"[{\"id\":18,\"pjp_code\":1000,\"salesman_name\":\"DEN BEI\"},{\"id\":19,\"pjp_code\":4000,\"salesman_name\":\"JACK DANIEL\"}]"`
}
