package entity

type AreaQueryFilter struct {
	CustID string `query:"cust_id" validate:"required"`
}
