package entity

type SAPReplExportQuery struct {
	DistributorCode string `query:"distributor_code" validate:"required"`
	DateFrom        int64  `query:"date_from" validate:"required,gte=1000000000"`
	DateTo          int64  `query:"date_to" validate:"required,lte=9999999999,gtefield=DateFrom"`
}

type SAPReplExportItem struct {
	CustReference       string                `json:"cust_reference"`
	ShipToPartyCode     string                `json:"ship_to_party_code,omitempty"`
	ReplenishmentType   string                `json:"replenishment_type"`
	IsAdditionFrom      bool                  `json:"is_addition_from"`
	Status              int                   `json:"status"`
	OrderType           string                `json:"order_type"`
	SalesOrganization   int                   `json:"sales_organization"`
	DistributionChannel string                `json:"distribution_channel"`
	SalesOffice         string                `json:"sales_office"`
	DeliveryDate        string                `json:"delivery_date"`
	ShippingPoint       string                `json:"shipping_point"`
	Plant               string                `json:"plant"`
	CreatedAt           string                `json:"created_at"`
	Details             []SAPReplExportDetail `json:"details"`
}

type SAPReplExportDetail struct {
	Material    string  `json:"material"`
	Division    string  `json:"division"`
	OrderQty    float64 `json:"order_qty"`
	Uom         string  `json:"uom"`
	CustPoPrice float64 `json:"cust_po_price"`
}
