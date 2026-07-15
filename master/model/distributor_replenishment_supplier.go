package model

type DistributorReplenishmentSupplierRow struct {
	SupID   int    `db:"sup_id"`
	SupCode string `db:"sup_code"`
	SupName string `db:"sup_name"`
}
