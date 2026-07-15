package model

import "time"

type SAPReplExportRow struct {
	ReplenishmentNo     string     `gorm:"column:replenishment_no"`
	ShipToPartyCode     *string    `gorm:"column:ship_to_party_code"`
	ReplenishmentType   string     `gorm:"column:replenishment_type"`
	IsAdditionFrom      bool       `gorm:"column:is_addition_from"`
	Status              int        `gorm:"column:status"`
	DistributionChannel *string    `gorm:"column:distribution_channel"`
	SalesOffice         *string    `gorm:"column:sales_office"`
	DeliveryDate        *time.Time `gorm:"column:delivery_date"`
	ShippingPoint       *string    `gorm:"column:shipping_point"`
	Plant               *string    `gorm:"column:plant"`
	CreatedAt           time.Time  `gorm:"column:created_at"`
	Material            *string    `gorm:"column:material"`
	OrderQty            float64    `gorm:"column:order_qty"`
	Uom                 *string    `gorm:"column:uom"`
	CustPoPrice         float64    `gorm:"column:cust_po_price"`
	Division            *string    `gorm:"column:division"`
}
