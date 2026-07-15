package model

import "database/sql"

// DistributorReplenishmentSetup row for list query (joined)
type DistributorReplenishmentSetup struct {
	ID                 int            `db:"id"`
	SupID              int            `db:"sup_id"`
	SupCode            string         `db:"sup_code"`
	SupName            string         `db:"sup_name"`
	DistributorID      int            `db:"distributor_id"`
	DistributorCode    string         `db:"distributor_code"`
	DistributorName    string         `db:"distributor_name"`
	DistributorType    string         `db:"distributor_type"`
	WhLimitAction      string         `db:"wh_limit_action"`
	WhCapacity         sql.NullInt64  `db:"wh_capacity"`
	WhVolume           sql.NullInt64  `db:"wh_volume"`
	CreditLimitAction  int            `db:"credit_limit_action"`
	PlafonCredit       sql.NullInt64  `db:"plafon_credit"`
	LeadTimeDays       int            `db:"lead_time_days"`
	IsApprovalRequired bool           `db:"is_approval_required"`
	CreatedBy          sql.NullInt64  `db:"created_by"`
	CreatedByName      sql.NullString `db:"created_by_name"`
	CreatedAt          sql.NullTime   `db:"created_at"`
	UpdatedBy          sql.NullInt64  `db:"updated_by"`
	UpdatedByName      sql.NullString `db:"updated_by_name"`
	UpdatedAt          sql.NullTime   `db:"updated_at"`
}

type DistributorReplenishmentDistributorRow struct {
	ID              int    `db:"id"`
	DistributorID   int    `db:"distributor_id"`
	DistributorCode string `db:"distributor_code"`
	DistributorName string `db:"distributor_name"`
}
