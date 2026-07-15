package entity

type RegionQueryFilter struct {
	CustID string `query:"cust_id" validate:"required"`
}
