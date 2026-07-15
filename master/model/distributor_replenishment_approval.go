package model

import "database/sql"

type DistributorReplenishmentApproval struct {
	ID                       int            `db:"id"`
	DistReplenishmentSetupID int            `db:"dist_replenishment_setup_id"`
	Level                    int            `db:"level"`
	Sequence                 int            `db:"sequence"`
	BusinessUnit             int            `db:"business_unit"`
	BusinessUnitName         sql.NullString `db:"business_unit_name"`
	Pic                      int            `db:"pic"`
	PicName                  sql.NullString `db:"pic_name"`
	IsActive                 bool           `db:"is_active"`
}
