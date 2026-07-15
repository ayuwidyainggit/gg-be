package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/constant"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type mOutletRepositoryImpl struct {
	*sqlx.DB
}

// Store implements OutletRepository.
func (*mOutletRepositoryImpl) Store(outlet model.MOutlet) (int, error) {
	panic("unimplemented")
}

type mOutletTransaction struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

type MOutletRepository interface {
	TrxBegin() (*mOutletTransaction, error)
	MOutletCheck(ctx context.Context, year int, cust_id string, status []string) (model.MOutletCode, error)
	FindOneByOutletCodeAndCustId(OutletName string, custId, parentCustId string) (model.MOutletRead, error)
	FindOneByIdOutlet(OutletId []int, custId, parentCustId string) ([]model.MOutletRead, error)
	FindByIdOutletFromList(ctx context.Context, OutletId []int) ([]model.MOutletReadFromList, error)
	FindAllByCustId(dataFilter entity.MOutletQueryFilter, custId, parentCustId string, empId int64) (consPro []model.MOutletRead, total int64, lastPage int, err error)
	FindAllByCustIdFromList(dataFilter entity.MOutletQueryFilter, custId string, empId int64, parentCustId string) (consPro []model.MOutletRead, total int64, lastPage int, err error)
	Delete(custId string, outletId int, deletedBy int64) error
	FindOneByOutletIdAndCustId(outletId int64, custId, parentCustId string) (model.OutletReads, error)
	GetDetailOutletContact(outletId int64, custId string) ([]model.MOutletContactReads, error)
	GetOutletCountsByNames(outletName []string) ([]model.MOutletFromListSimilar, error)
	FindMobileOutletList(dataFilter entity.MobileOutletListQueryFilter, custId string, empID int64) (outlets []model.MobileOutletList, total int64, lastPage int, err error)
	FindMobileOutletDetail(outletId int64, custId string) (model.MobileOutletDetail, error)
	FindMobileOutletContact(outletId int64, custId string) (*model.MobileOutletContactDetail, error)
	GetAllOutletLookupByCusttomerId(ctx context.Context, dataFilter entity.OutletLookupQueryFilter, custId string) ([]model.OutletLookupList, int64, int, error)
	FindTopStatus(ctx context.Context, custID string) (int, error)
	FindRegionByDistributorID(ctx context.Context, distributorID int) (model.Region, error)
	CheckOutletByCode(ctx context.Context, outletCode string, custId string) (bool, error)
	FindByDestinations(ctx context.Context, destinationType string, destinationIDs []int) ([]model.MOutletReadFromList, error)
}

func NewMOutletRepository(db *sqlx.DB) *mOutletRepositoryImpl {
	return &mOutletRepositoryImpl{db}
}

func NewMOutletTransaction(db *sqlx.DB) (trxObj *mOutletTransaction, err error) {
	trx := db.MustBegin()

	return &mOutletTransaction{tx: trx, db: db}, nil
}

func (repo *mOutletRepositoryImpl) TrxBegin() (*mOutletTransaction, error) {
	trxObj, err := NewMOutletTransaction(repo.DB)
	if err != nil {
		return nil, err
	}
	return trxObj, nil
}

func (repo *mOutletTransaction) TrxCommit() error {
	return repo.tx.Commit()
}

func (repo *mOutletTransaction) TrxRollback() error {
	return repo.tx.Rollback()
}

func (repository *mOutletRepositoryImpl) MOutletCheck(ctx context.Context, year int, cust_id string, status []string) (model.MOutletCode, error) {
	query :=
		`SELECT id, cust_id, serial_code, year_code, last_sequence_no, status, created_at, created_by, updated_at, updated_by
		FROM mst.m_outlet_code
		WHERE year_code = $1 AND cust_id = $2 AND status = ANY($3)`

	var outlet model.MOutletCode
	err := repository.DB.QueryRowContext(ctx, query, year, cust_id, pq.Array(status)).Scan(
		&outlet.Id,
		&outlet.CustId,
		&outlet.SerialCode,
		&outlet.YearCode,
		&outlet.LastSequenceNo,
		&outlet.Status,
		&outlet.CreatedAt,
		&outlet.CreatedBy,
		&outlet.UpdatedAt,
		&outlet.UpdatedBy,
	)
	if err != nil {
		log.Println("mOutletTransaction, MOutletCheck, err:", err.Error())
		return model.MOutletCode{}, err
	}

	return outlet, nil
}

func (repository *mOutletRepositoryImpl) FindOneByOutletCodeAndCustId(OutletName string, custId, parentCustId string) (model.MOutletRead, error) {
	outlet := model.MOutletRead{}
	query := `SELECT o.outlet_name FROM mst.m_outlet o
	LEFT JOIN sys.m_user u ON u.user_id = o.updated_by
	LEFT JOIN mst.m_disc_group dg ON dg.disc_grp_id = o.disc_grp_id AND dg.cust_id = '` + custId + `'
	LEFT JOIN mst.m_outlet_loc ol ON ol.ot_loc_id = o.ot_loc_id AND ol.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_outlet_group og ON og.ot_grp_id = o.ot_grp_id AND og.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_price_group pg ON pg.price_grp_id = o.price_grp_id AND pg.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_district dis ON dis.district_id = o.district_id AND dis.cust_id = '` + custId + `'
	LEFT JOIN mst.m_beat b ON b.beat_id = o.beat_id AND b.cust_id = '` + custId + `'
	LEFT JOIN mst.m_sub_beat sb ON sb.sbeat_id = o.sbeat_id AND sb.cust_id = '` + custId + `'
	LEFT JOIN mst.m_outlet_class oc ON oc.ot_class_id = o.ot_class_id AND oc.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_industry ind ON ind.industry_id = o.industry_id AND ind.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_market m ON m.market_id = o.market_id AND m.cust_id = '` + custId + `'
	LEFT JOIN mst.m_plu_group mpg ON mpg.plu_grp_id = o.plu_grp_id AND mpg.cust_id = '` + custId + `'
	LEFT JOIN mst.m_conv_group mcg ON mcg.conv_grp_id = o.conv_grp_id AND mcg.cust_id = '` + custId + `'
	LEFT JOIN mst.m_invoice_disc mid ON mid.inv_disc_id = o.disc_inv_id  AND mid.cust_id = '` + custId + `'
	WHERE o.outlet_name = $1 and o.is_del = FALSE AND o.cust_id = $2`
	err := repository.Get(&outlet, query, OutletName, custId)
	if err != nil {
		log.Println("outletRepository, FindOneByoutletCodeAndCustId, err:", err.Error())
		return outlet, err
	}

	return outlet, nil
}

func (repository *mOutletRepositoryImpl) FindOneByIdOutlet(OutletId []int, custId, parentCustId string) ([]model.MOutletRead, error) {

	outlets := []model.MOutletRead{}
	query := `SELECT o.outlet_id, o.outlet_name, o.outlet_code, o.address1, o.phone_no, o.wa_no, o.fax_no FROM mst.m_outlet o
	LEFT JOIN sys.m_user u ON u.user_id = o.updated_by
	LEFT JOIN mst.m_disc_group dg ON dg.disc_grp_id = o.disc_grp_id AND dg.cust_id = '` + custId + `'
	LEFT JOIN mst.m_outlet_loc ol ON ol.ot_loc_id = o.ot_loc_id AND ol.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_outlet_group og ON og.ot_grp_id = o.ot_grp_id AND og.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_price_group pg ON pg.price_grp_id = o.price_grp_id AND pg.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_district dis ON dis.district_id = o.district_id AND dis.cust_id = '` + custId + `'
	LEFT JOIN mst.m_beat b ON b.beat_id = o.beat_id AND b.cust_id = '` + custId + `'
	LEFT JOIN mst.m_sub_beat sb ON sb.sbeat_id = o.sbeat_id AND sb.cust_id = '` + custId + `'
	LEFT JOIN mst.m_outlet_class oc ON oc.ot_class_id = o.ot_class_id AND oc.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_industry ind ON ind.industry_id = o.industry_id AND ind.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_market m ON m.market_id = o.market_id AND m.cust_id = '` + custId + `'
	LEFT JOIN mst.m_plu_group mpg ON mpg.plu_grp_id = o.plu_grp_id AND mpg.cust_id = '` + custId + `'
	LEFT JOIN mst.m_conv_group mcg ON mcg.conv_grp_id = o.conv_grp_id AND mcg.cust_id = '` + custId + `'
	LEFT JOIN mst.m_invoice_disc mid ON mid.inv_disc_id = o.disc_inv_id  AND mid.cust_id = '` + custId + `'
	WHERE o.outlet_id IN (?) and o.is_del = FALSE AND o.cust_id = ?`
	// Convert []int to []interface{} for sqlx.In
	ids := make([]interface{}, len(OutletId))
	for i, v := range OutletId {
		ids[i] = v
	}

	// Rebuild query with sqlx.In to handle the IN clause
	fmt.Println("===>>>", ids)
	query, args, err := sqlx.In(query, ids, custId)
	if err != nil {
		log.Println("outletRepository, FindOneByIdOutlet, sqlx.In err:", err.Error())
		return outlets, err
	}

	// Rebind for the specific database (PostgreSQL uses $1, $2, etc.)
	query = repository.Rebind(query)

	err = repository.Select(&outlets, query, args...)
	if err != nil {
		log.Println("outletRepository, FindOneByIdOutlet, err:", err.Error())
		return outlets, err
	}

	return outlets, nil
}

func (repository *mOutletRepositoryImpl) FindByIdOutletFromList(ctx context.Context, OutletId []int) ([]model.MOutletReadFromList, error) {
	outletStatuses := []int{
		constant.NewOutletStatus,
		constant.DormantOutletStatus,
		constant.RegisteredOutletStatus,
		constant.ActiveOutletStatus,
	}

	query := `
				SELECT
						o.outlet_id,
						o.outlet_name,
						o.outlet_code,
						o.address1,
						o.cust_id,
						o.phone_no,
						o.wa_no,
						o.fax_no,
						o.latitude,
						o.longitude,
						o.cust_id
				FROM mst.m_outlet o
				WHERE
					o.outlet_id in (?)
					AND o.is_del = false
				  AND o.outlet_status in (?);`

	outlets := []model.MOutletReadFromList{}
	// Convert []int to []interface{} for sqlx.In
	ids := make([]interface{}, len(OutletId))
	for i, v := range OutletId {
		ids[i] = v
	}

	query, args, err := sqlx.In(query, ids, outletStatuses)
	if err != nil {
		log.Println("outletRepository, FindOneByIdOutlet, sqlx.In err:", err.Error())
		return outlets, err
	}

	// Rebind for the specific database (PostgreSQL uses $1, $2, etc.)
	query = repository.Rebind(query)

	err = repository.SelectContext(ctx, &outlets, query, args...)
	if err != nil {
		log.Println("outletRepository, FindOneByIdOutlet, err:", err.Error())
		return outlets, err
	}

	return outlets, nil
}

func (repository *mOutletTransaction) Store(outlet *model.MOutlet) error {
	query :=
		`INSERT INTO mst.m_outlet(
			cust_id, outlet_name, outlet_code, outlet_principal_code, source, outlet_status, address1, phone_no,
			wa_no, fax_no, building_own, latitude, longitude, file_url, is_active, created_by, created_at, updated_by, updated_at, is_del, deleted_by, deleted_at,  verification_status,disc_grp_id,ot_grp_id,price_grp_id,ot_class_id,ot_type_id, credit_limit_action, credit_limit_action_name)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
		) RETURNING outlet_id;`
	var lastInsertId int
	err := repository.tx.QueryRow(query,
		outlet.CustId, outlet.OutletName, outlet.OutletCode, outlet.OutletPrincipalCode, outlet.Source, outlet.OutletStatus, outlet.Address, outlet.PhoneNo,
		outlet.WaNo, outlet.FaxNo, outlet.BuldingOwn, outlet.Latitude, outlet.Longitude, outlet.FileUrl, outlet.IsActive, outlet.CreatedBy, outlet.CreatedAt, outlet.UpdatedBy,
		outlet.UpdatedAt, outlet.IsDel, outlet.DeletedBy, outlet.DeletedAt, outlet.VerificationStatus, outlet.NonDiscountGroupID, outlet.NonOutletGroupID, outlet.NonPriceGroupID, outlet.NonOutletClassID, outlet.NonOutletTypeID, outlet.CreditLimitAction, outlet.CreditLimitActionName,
	).Scan(&lastInsertId)
	if err != nil {
		log.Println("outletRepository, Store, err:", err.Error())
		return err
	}
	outlet.OutletId = lastInsertId
	return nil
}

func (repository *mOutletTransaction) StoreFromList(ctx context.Context, outlet *model.MOutletCreadFromList) error {

	query :=
		`INSERT INTO pjp.route_outlet_history(
			route_code, route_name, outlet_id, outlet_code, outlet_name, longitude, latitude, outlet_status, outlet_address, pjp_id, pjp_code, cust_id, status, created_at, updated_at, verified_date, 
			old_pjp_id, old_pjp_code, old_route_code, old_route_name, photo, avg_sales_week, index_day, start_week, is_in_current_year, week, year, date, is_additional, is_extra_call)
		VALUES ( 
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
		) RETURNING outlet_id;`
	var lastInsertID int
	err := repository.tx.QueryRowContext(ctx, query,
		outlet.RouteCode, outlet.RouteName, outlet.OutletId, outlet.OutletCode, outlet.OutletName, outlet.Longitude, outlet.Latitude, outlet.OutletStatus, outlet.OutletAddress, outlet.PjpId, outlet.PjpCode, outlet.CustId, outlet.Status, outlet.CreatedAt, outlet.UpdatedAt, outlet.VerifiedDate,
		outlet.OldPjpId, outlet.OldPjpCode, outlet.OldRouteCode, outlet.OldRouteName, outlet.Photo, outlet.AvgSalesWeek, outlet.IndexDay, outlet.StartWeek, outlet.IsInCurrentYear, outlet.Week, outlet.Year, outlet.Date, outlet.IsAdditional, outlet.IsExtraCall).Scan(&lastInsertID)
	if err != nil {
		log.Println("outletRepository, Store, err:", err.Error())
		return err
	}
	outlet.Id = lastInsertID
	return nil
}

func (repository *mOutletTransaction) StoreFromListPrinciple(ctx context.Context, outlet *model.MOutletCreadFromList) error {
	query := `
				INSERT INTO pjp_principles.destinations_history (	
						route_code, route_name, verified_date, "date", week, "year", index_day, start_week, is_in_current_year, is_additional,
						destination_id, destination_code, destination_status, destination_name, destination_address, destination_type, longitude, latitude,
						pjp_id, pjp_code, old_pjp_id, old_pjp_code, old_route_code, old_route_name, photo, avg_sales_week, cust_id, is_extra_call)
		VALUES( 
						$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
						$11, $12, $13, $14, $15, $16, $17, $18,
						$19, $20, $21, $22, $23, $24, $25, $26, $27, $28) RETURNING id;`

	var lastInsertID int
	err := repository.tx.QueryRowContext(ctx, query,
		outlet.RouteCode, outlet.RouteName, outlet.VerifiedDate, outlet.Date, outlet.Week, outlet.Year, outlet.IndexDay, outlet.StartWeek, outlet.IsInCurrentYear, outlet.IsAdditional,
		outlet.OutletId, outlet.OutletCode, outlet.OutletStatus, outlet.OutletName, outlet.OutletAddress, "outlet", outlet.Longitude, outlet.Latitude,
		outlet.PjpId, outlet.PjpCode, outlet.OldPjpId, outlet.OldPjpCode, outlet.OldRouteCode, outlet.OldRouteName, outlet.Photo, outlet.AvgSalesWeek, outlet.CustId, outlet.IsExtraCall).Scan(&lastInsertID)
	if err != nil {
		log.Println("outletRepository, Store, err:", err.Error())
		return err
	}

	outlet.Id = lastInsertID
	return nil
}

func (repository *mOutletTransaction) StoreFromListOutletVisitList(ctx context.Context, outlet *model.MOutletCreadFromListOutletVisitList) error {

	query :=
		`INSERT INTO pjp.outlet_visit_list(
			year, week, date, day, route_code, outlet_code, pjp_id, pjp_code, created_at, updated_at, outlet_id, is_planned, latitude, longitude, is_extra_call, start)
		VALUES ( 
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		) RETURNING outlet_id;`
	var lastInsertID int
	err := repository.tx.QueryRowContext(ctx, query,
		outlet.Year, outlet.Week, outlet.Date, outlet.Day, outlet.RouteCode, outlet.OutletCode, outlet.PjpId, outlet.PjpCode, outlet.CreatedAt, outlet.UpdatedAt, outlet.OutletId, outlet.IsPlanned, outlet.Latitude, outlet.Longitude, outlet.IsExtraCall, outlet.Start).Scan(&lastInsertID)
	if err != nil {
		log.Println("outletVisitRepository, Store, err:", err.Error())
		return err
	}
	outlet.Id = lastInsertID
	return nil
}

func (repository *mOutletTransaction) StoreFromListOutletVisitListPrinciple(ctx context.Context, outlet *model.MOutletCreadFromListOutletVisitList) error {
	query := `
				INSERT INTO pjp_principles.outlet_visit_list
						("year", week, "date", "day", route_code, outlet_code, pjp_id, pjp_code, outlet_id, is_planned, is_extra_call, start, destination_type) VALUES
						($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id;`

	var lastInsertId int
	err := repository.tx.QueryRowContext(ctx, query,
		outlet.Year, outlet.Week, outlet.Date, outlet.Day, outlet.RouteCode, outlet.OutletCode, outlet.PjpId, outlet.PjpCode, outlet.OutletId, outlet.IsPlanned, outlet.IsExtraCall, outlet.Start, outlet.DestinationType).Scan(&lastInsertId)
	if err != nil {
		log.Println("outletVisitRepository, Store, err:", err.Error())
		return err
	}
	outlet.Id = lastInsertId
	return nil
}

func (repository *mOutletTransaction) StoreDetailContact(outletContact *model.MOutletContact) error {
	query :=
		`INSERT INTO mst.m_outlet_contact(
			cust_id, outlet_id, contact_name, 
			job_title, phone_no, wa_no,
			email, identity_no)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6,
			$7, $8
		) RETURNING outlet_contact_id;`
	var lastInsertId int64
	err := repository.tx.QueryRow(query,
		outletContact.CustID, outletContact.OutletID, outletContact.ContactName,
		outletContact.JobTitle, outletContact.PhoneNo, outletContact.WaNo,
		outletContact.Email, outletContact.IdentityNo).Scan(&lastInsertId)
	if err != nil {
		log.Println("outletRepository, StoreDetailContact, err:", err.Error())
		return err
	}
	outletContact.OutletContactId = &lastInsertId
	return nil
}

func (repository *mOutletTransaction) UpdateOutletCodeSequence(id string, lastSequenceNo string, updatedBy int64) error {
	query := `UPDATE mst.m_outlet_code SET last_sequence_no = $1, updated_by = $2, updated_at = NOW() WHERE id = $3`
	_, err := repository.tx.Exec(query, lastSequenceNo, updatedBy, id)
	if err != nil {
		log.Println("mOutletTransaction, UpdateOutletCodeSequence, err:", err.Error())
		return err
	}
	return nil
}

func (repository *mOutletRepositoryImpl) FindAllByCustId(dataFilter entity.MOutletQueryFilter, custId, parentCustId string, empId int64) ([]model.MOutletRead, int64, int, error) {

	outlets := []model.MOutletRead{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` o.outlet_name, o.outlet_id, o.address1 , o.phone_no, o.wa_no, o.fax_no, o.email, o.building_own,
	o.latitude, o.longitude, o.is_active, COALESCE(o.created_by, 0) AS salesman_id, o.verification_status, o.verified_at, o.verified_by,
u2.user_fullname as verified_by_name, o.outlet_code, COALESCE(o.outlet_status, 0) AS outlet_status
	`
	qWhere := ` WHERE o.is_del = false and o.verification_status = 1 and o.outlet_status <> 4 and o.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (o.outlet_code ILIKE '%` + dataFilter.Query + `%' 
					OR o.outlet_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SalesID != 0 {
		qWhere += `and o.created_by = ` + strconv.Itoa(dataFilter.SalesID) + ` `
	}

	if dataFilter.CreatedBy != 0 {
		qWhere += `and o.created_by = ` + strconv.Itoa(dataFilter.CreatedBy) + ` `
	}

	// if dataFilter.IsAdditional != nil {
	// 	if *dataFilter.IsAdditional {
	// 		qWhere += ` AND o.is_additional = true `
	// 	} else {
	// 		qWhere += ` AND o.is_additional = false `
	// 	}
	// }

	if len(dataFilter.VerificationStatus) > 0 {
		ids := make([]string, len(dataFilter.VerificationStatus))
		for i, id := range dataFilter.VerificationStatus {
			ids[i] = strconv.Itoa(id)
		}
		qWhere += ` and o.verification_status IN (` + strings.Join(ids, ",") + `) `
	}

	qFrom := ` FROM mst.m_outlet o
	LEFT JOIN sys.m_user u ON u.user_id = o.updated_by
	LEFT JOIN sys.m_user u2 ON u2.user_id = o.verified_by
	LEFT JOIN mst.m_disc_group dg ON dg.disc_grp_id = o.disc_grp_id AND dg.cust_id = '` + custId + `'
	LEFT JOIN mst.m_outlet_loc ol ON ol.ot_loc_id = o.ot_loc_id AND ol.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_outlet_group og ON og.ot_grp_id = o.ot_grp_id AND og.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_price_group pg ON pg.price_grp_id = o.price_grp_id AND pg.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_district dis ON dis.district_id = o.district_id AND dis.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_beat b ON b.beat_id = o.beat_id AND b.cust_id = '` + custId + `'
	LEFT JOIN mst.m_sub_beat sb ON sb.sbeat_id = o.sbeat_id AND sb.cust_id = '` + custId + `'
	LEFT JOIN mst.m_outlet_class oc ON oc.ot_class_id = o.ot_class_id AND oc.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_industry ind ON ind.industry_id = o.industry_id AND ind.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_market m ON m.market_id = o.market_id AND m.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_plu_group mpg ON mpg.plu_grp_id = o.plu_grp_id AND mpg.cust_id = '` + custId + `'
	LEFT JOIN mst.m_conv_group mcg ON mcg.conv_grp_id = o.conv_grp_id AND mcg.cust_id = '` + custId + `'
	LEFT JOIN mst.m_invoice_disc mid ON mid.inv_disc_id = o.disc_inv_id  AND mid.cust_id = '` + custId + `' 
	
	LEFT JOIN mst.m_ward owrd on owrd.ward_id = o.outlet_ward_id AND owrd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sub_district omsd on omsd.sub_district_id = owrd.sub_district_id AND omsd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_regency omr on omr.regency_id = omsd.regency_id AND omr.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_province omp on omp.province_id = omr.province_id AND omp.cust_id = '` + parentCustId + `'
	
	LEFT JOIN mst.m_ward dlvwrd on dlvwrd.ward_id = o.delv_ward_id AND dlvwrd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sub_district dlvmsd on dlvmsd.sub_district_id = dlvwrd.sub_district_id AND dlvmsd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_regency dlvmr on dlvmr.regency_id = dlvmsd.regency_id AND dlvmr.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_province dlvmp on dlvmp.province_id = dlvmr.province_id AND dlvmp.cust_id = '` + parentCustId + `'

	LEFT JOIN mst.m_ward invwrd on invwrd.ward_id = o.inv_ward_id AND invwrd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sub_district invmsd on invmsd.sub_district_id = invwrd.sub_district_id AND invmsd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_regency invmr on invmr.regency_id = invmsd.regency_id AND invmr.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_province invmp on invmp.province_id = invmr.province_id AND invmp.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_outlet_type ot on ot.ot_type_id = o.ot_type_id AND ot.cust_id = '` + parentCustId + `'
	`
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("outletRepository, queryCount:", queryCount)
	var total int64
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("outletRepository, count total, err:", err.Error())
		return outlets, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `outlet_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))

	querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))

	// log.Println("outletRepository, querySelect:", querySelect)
	err = repository.Select(&outlets, querySelect)
	if err != nil {
		log.Println("outletRepository, FindAllByCustId, err:", err.Error())
		return outlets, total, lastPage, err
	}

	return outlets, total, lastPage, nil
}

func (repository *mOutletRepositoryImpl) FindAllByCustIdFromList(dataFilter entity.MOutletQueryFilter, custId string, empId int64, parentCustId string) ([]model.MOutletRead, int64, int, error) {

	outlets := []model.MOutletRead{}
	selectCount := ` COUNT(*) AS total `
	selectField := `c.outlet_name, c.outlet_id, c.Address1, b.outlet_code, COALESCE(c.verification_status, '0') as verification_status, 
	COALESCE(c.outlet_status, '0') as outlet_status, COALESCE(c.created_by, '0') as salesman_id, c.phone_no, c.wa_no, c.fax_no, c.email, c.building_own, c.created_at`

	qWhere := ` WHERE a.salesman_id = '` + strconv.FormatInt(empId, 10) + `' and  b.is_planned = TRUE and b.date >= CURRENT_DATE `

	if dataFilter.Query != "" {
		qWhere += ` AND (b.outlet_code ILIKE '%` + dataFilter.Query + `%' 
					OR b.outlet_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SalesID != 0 {
		qWhere += `and a.salesman_id = ` + strconv.Itoa(dataFilter.SalesID) + ` `
	}

	if dataFilter.CreatedBy != 0 {
		qWhere += `and a.salesman_id = ` + strconv.Itoa(dataFilter.CreatedBy) + ` `
	}

	if len(dataFilter.VerificationStatus) > 0 {
		ids := make([]string, len(dataFilter.VerificationStatus))
		for i, id := range dataFilter.VerificationStatus {
			ids[i] = strconv.Itoa(id)
		}
		qWhere += ` and c.verification_status IN (` + strings.Join(ids, ",") + `) `
	}

	qFrom := ` FROM pjp.permanent_journey_plans a LEFT JOIN pjp.outlet_visit_list b on b.pjp_id = a.id LEFT JOIN mst.m_outlet c on c.outlet_id = b.outlet_id `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int64
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		return outlets, 0, 0, err
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `b.created_at`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))

	querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))

	err = repository.Select(&outlets, querySelect)
	if err != nil {
		return outlets, total, lastPage, err
	}

	return outlets, total, lastPage, nil
}

func (repository *mOutletRepositoryImpl) FindMobileOutletDetail(outletId int64, custId string) (model.MobileOutletDetail, error) {
	outlet := model.MobileOutletDetail{}
	query := `SELECT 
		o.outlet_id,
		COALESCE(o.outlet_code, '') AS outlet_code,
		COALESCE(o.outlet_name, '') AS outlet_name,
		COALESCE(o.address1, '') AS address1,
		COALESCE(o.phone_no, '') AS phone_no,
		COALESCE(o.building_own, 0) AS building_own,
		COALESCE(o.file_url, '') AS file_url,
		COALESCE(o.longitude, '') AS longitude,
		COALESCE(o.latitude, '') AS latitude
	FROM mst.m_outlet o
	WHERE o.outlet_id = $1 
		AND o.cust_id = $2 
		AND o.is_del = false`

	err := repository.Get(&outlet, query, outletId, custId)
	if err != nil {
		return outlet, err
	}
	return outlet, nil
}

func (repository *mOutletRepositoryImpl) FindMobileOutletContact(outletId int64, custId string) (*model.MobileOutletContactDetail, error) {
	contact := model.MobileOutletContactDetail{}
	query := `SELECT 
		COALESCE(contact_name, '') AS contact_name,
		COALESCE(job_title, '') AS job_title,
		COALESCE(phone_no, '') AS phone_no,
		COALESCE(wa_no, '') AS wa_no,
		COALESCE(email, '') AS email
	FROM mst.m_outlet_contact
	WHERE outlet_id = $1 AND cust_id = $2
	LIMIT 1`

	err := repository.Get(&contact, query, outletId, custId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &contact, nil
}

func (repository *mOutletRepositoryImpl) Delete(custId string, outletId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_outlet
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND outlet_id = :outlet_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"outlet_id":  outletId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("OutletRepository, Delete, err:", err.Error())
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (repository *mOutletRepositoryImpl) GetAllOutletLookupByCusttomerId(ctx context.Context, dataFilter entity.OutletLookupQueryFilter, custId string) ([]model.OutletLookupList, int64, int, error) {
	outlets := []model.OutletLookupList{}

	selectCount := ` COUNT(*) AS total `
	selectField := ` mo.outlet_id, COALESCE(mo.outlet_code, '') AS outlet_code, COALESCE(mo.outlet_name, '') AS outlet_name, COALESCE(mo.address1, '') AS address, mo.latitude, mo.longitude, mo.outlet_status, md.distributor_id, COALESCE(md.region_id, 0) AS region_id, COALESCE(md.area_id,0) AS area_id `
	qFrom := ` FROM mst.m_outlet mo LEFT JOIN mst.m_distributor md ON mo.cust_id = md.cust_id `
	qWhere := fmt.Sprintf(` WHERE mo.cust_id LIKE '%s%%' AND mo.is_del = false AND mo.verification_status = 1 `, custId)

	// if dataFilter.IsActive != nil {
	// 	if *dataFilter.IsActive == 1 {
	// 		qWhere += ` AND mo.is_active = true `
	// 	} else if *dataFilter.IsActive == 2 {
	// 		qWhere += ` AND mo.is_active = false `
	// 	}
	// } else {
	// 	qWhere += ` AND mo.is_active = true `
	// }

	if dataFilter.Search != "" {
		qWhere += ` AND (mo.outlet_code ILIKE '%` + dataFilter.Search + `%' OR mo.outlet_name ILIKE '%` + dataFilter.Search + `%') `
	}

	// Collect hierarchy conditions; when OutletPrinciple is true,
	// they are grouped into one parenthesised OR with mo.cust_id.
	// ponytail: simple slice build — fine for <100 filters, O(n) per slice.
	var principleConds []string

	if len(dataFilter.RegionID) > 0 {
		regionsID := make([]string, len(dataFilter.RegionID))
		for i, v := range dataFilter.RegionID {
			regionsID[i] = strconv.Itoa(v)
		}
		cond := `md.region_id IN (` + strings.Join(regionsID, ",") + `)`
		if dataFilter.OutletPrinciple {
			principleConds = append(principleConds, cond)
		} else {
			qWhere += ` AND ` + cond + ` `
		}
	}

	if len(dataFilter.AreaID) > 0 {
		areasID := make([]string, len(dataFilter.AreaID))
		for i, v := range dataFilter.AreaID {
			areasID[i] = strconv.Itoa(v)
		}
		cond := `md.area_id IN (` + strings.Join(areasID, ",") + `)`
		if dataFilter.OutletPrinciple {
			principleConds = append(principleConds, cond)
		} else {
			qWhere += ` AND ` + cond + ` `
		}
	}

	if len(dataFilter.DistributorID) > 0 {
		distributorsID := make([]string, len(dataFilter.DistributorID))
		for i, v := range dataFilter.DistributorID {
			distributorsID[i] = strconv.Itoa(v)
		}
		cond := `md.distributor_id IN (` + strings.Join(distributorsID, ",") + `)`
		if dataFilter.OutletPrinciple {
			principleConds = append(principleConds, cond)
		} else {
			qWhere += ` AND (` + cond + `) `
		}
	}

	if dataFilter.OutletPrinciple {
		if len(principleConds) > 0 {
			qWhere += ` AND (` + strings.Join(principleConds, " AND ") + ` OR mo.cust_id = '` + custId + `') `
		} else {
			qWhere += ` AND (mo.cust_id = '` + custId + `') `
		}
	}

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	queryCount += fmt.Sprintf(` AND mo.outlet_status != 4`)
	querySelect += fmt.Sprintf(` AND mo.outlet_status != 4`)

	var total int64
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("outletRepository, GetAllOutletLookupByCusttomerId, count total, err:", err.Error())
		return outlets, 0, 0, err
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`mo.%s %s, `, colSort[0], strings.ToUpper(colSort[1]))
			}
		}
		if sortBy != "" {
			sortBy = strings.TrimSuffix(sortBy, ", ")
			querySelect += fmt.Sprintf(` ORDER BY %s`, sortBy)
		}
	} else {
		querySelect += ` ORDER BY mo.outlet_id DESC`
	}

	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 10
	}

	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(total) / float64(dataFilter.Limit)))

	querySelect += fmt.Sprintf(` LIMIT %d OFFSET %d`, dataFilter.Limit, offset)

	err = repository.Select(&outlets, querySelect)
	if err != nil {
		log.Println("outletRepository, GetAllOutletLookupByCusttomerId, select, err:", err.Error())
		return outlets, total, lastPage, err
	}

	return outlets, total, lastPage, nil
}

func (repository *mOutletRepositoryImpl) FindOneByOutletIdAndCustId(outletId int64, custId, parentCustId string) (model.OutletReads, error) {
	log.Println("outletRepositoryImpl, outletId:", outletId)
	outlet := model.OutletReads{}

	// if custId is principle, must find distributor custId first.
	if len(custId) == 6 {
		var distributorCustId string

		qwFind := `SELECT cust_id FROM mst.m_outlet WHERE outlet_id = $1 AND is_del = false`
		errFind := repository.QueryRow(qwFind, outletId).Scan(&distributorCustId)
		if errFind != nil {
			log.Println("mOutletRepositoryImpl, FindOneByOutletIdAndCustId, err:", errFind.Error())
			return outlet, errFind
		}
		if distributorCustId != "" {
			custId = distributorCustId
		}
	}
	query := `SELECT o.*,
					CASE 
						WHEN o.payment_type = 1 THEN 'Cash On Delivery'
						WHEN o.payment_type = 2 THEN 'Cash Before Delivery'
						WHEN o.payment_type = 3 THEN 'Credit'
						ELSE ''
					END AS payment_type_name,
					CASE 
						WHEN o.credit_limit_type = 2 THEN 2
						ELSE NULL
					END AS credit_limit_type,
					CASE 
						WHEN o.credit_limit_type = 1 THEN 'Unlimited'
						WHEN o.credit_limit_type = 2 THEN 'Limit By Total'
						WHEN o.credit_limit_type = 3 THEN 'Limit By Supplier'
						ELSE 'Unlimited'
					END AS credit_limit_type_name,
					CASE 
						WHEN o.credit_limit_action = 1 THEN 1
						WHEN o.credit_limit_action = 2 THEN 2
						ELSE NULL
					END AS credit_limit_action,
					CASE 
						WHEN credit_limit_action = 1 THEN 'Warning'
						WHEN credit_limit_action = 2 THEN 'Restricted'
						ELSE ''
					END AS credit_limit_action_name,
					CASE 
						WHEN o.sales_inv_limit_type = 2 THEN 2
						ELSE NULL
					END AS sales_inv_limit_type,
					CASE 
						WHEN o.sales_inv_limit_type = 1 THEN 'Unlimited'
						WHEN o.sales_inv_limit_type = 2 THEN 'Limited By Invoice'
						ELSE 'Unlimited'
					END AS sales_inv_limit_type_name,
					CASE 
						WHEN o.sales_inv_limit_action = 1 THEN 1
						WHEN o.sales_inv_limit_action = 2 THEN 2
						ELSE NULL
					END AS sales_inv_limit_action,
					CASE 
						WHEN sales_inv_limit_action = 1 THEN 'Warning'
						WHEN sales_inv_limit_action = 2 THEN 'Restricted'
						ELSE ''
					END AS sales_inv_limit_action_name,
					CASE 
						WHEN o.obs_type = 2 THEN 2
						ELSE NULL
					END AS obs_type,
					CASE 
						WHEN o.obs_type = 1 THEN 'Unlimited'
						WHEN o.obs_type = 2 THEN 'Limited By Invoice'
						ELSE 'Unlimited'
					END AS obs_type_name,
					CASE 
						WHEN o.obs_limit_action = 1 THEN 1
						WHEN o.obs_limit_action = 2 THEN 2
						ELSE NULL
					END AS obs_limit_action,
					CASE 
						WHEN obs_limit_action = 1 THEN 'Warning'
						WHEN obs_limit_action = 2 THEN 'Restricted'
						ELSE ''
					END AS obs_limit_action_name,
	u.user_fullname AS updated_by_name,u2.user_fullname AS verified_by_name,dg.disc_grp_code,dg.disc_grp_name,ol.ot_loc_code,ol.ot_loc_name,
	og.ot_grp_code,og.ot_grp_name,pg.sp_price_grp_code as price_grp_code,pg.sp_price_grp_name as price_grp_name,dis.district_code,dis.district_name,
	b.beat_code,b.beat_name,sb.sbeat_code,sb.sbeat_name,oc.ot_class_code,oc.ot_class_name,ind.industry_code,ind.industry_name,
	m.market_code,m.market_name,mpg.plu_grp_code,mpg.plu_grp_name,mcg.conv_grp_code,mcg.conv_grp_name,mid.inv_disc_code as disc_inv_code,mid.inv_disc_name as disc_inv_name, ot.ot_type_name,
	
	owrd.ward as outlet_ward,
	omsd.sub_district_id as outlet_sub_district_id,omsd.sub_district as outlet_sub_district,
	omr.regency as outlet_regency,omr.regency_id as outlet_regency_id,
	omp.province as outlet_province,omp.province_id as outlet_province_id,
	
	dlvwrd.ward as delv_ward,
	dlvmsd.sub_district as delv_sub_district,dlvmsd.sub_district_id as delv_sub_district_id,
	dlvmr.regency as delv_regency,dlvmr.regency_id as delv_regency_id,
	dlvmp.province as delv_province,dlvmp.province_id as delv_province_id,

	dlvwrd2.ward as delv_ward2,
	dlvmsd2.sub_district as delv_sub_district2,dlvmsd2.sub_district_id as delv_sub_district_id2,
	dlvmr2.regency as delv_regency2,dlvmr2.regency_id as delv_regency_id2,
	dlvmp2.province as delv_province2,dlvmp2.province_id as delv_province_id2,
	
	invwrd.ward as inv_ward,
	invmsd.sub_district as inv_sub_district,invmsd.sub_district_id as inv_sub_district_id,
	invmr.regency as inv_regency,invmr.regency_id as inv_regency_id,
	invmp.province as inv_province,invmp.province_id as inv_province_id
	
	FROM mst.m_outlet o
	LEFT JOIN sys.m_user u ON u.user_id = o.updated_by
	LEFT JOIN sys.m_user u2 ON u2.user_id = o.verified_by
	LEFT JOIN mst.m_disc_group dg ON dg.disc_grp_id = o.disc_grp_id AND dg.cust_id = '` + custId + `'
	LEFT JOIN mst.m_outlet_loc ol ON ol.ot_loc_id = o.ot_loc_id AND ol.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_outlet_type ot ON ot.ot_type_id = o.ot_type_id AND ot.cust_id  = '` + parentCustId + `'
	LEFT JOIN mst.m_outlet_group og ON og.ot_grp_id = o.ot_grp_id AND og.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sp_price_group pg ON pg.sp_price_grp_id = o.price_grp_id AND pg.cust_id = '` + custId + `'
	LEFT JOIN mst.m_district dis ON dis.district_id = o.district_id AND dis.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_beat b ON b.beat_id = o.beat_id AND b.cust_id = '` + custId + `'
	LEFT JOIN mst.m_sub_beat sb ON sb.sbeat_id = o.sbeat_id AND sb.cust_id = '` + custId + `'
	LEFT JOIN mst.m_outlet_class oc ON oc.ot_class_id = o.ot_class_id AND oc.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_industry ind ON ind.industry_id = o.industry_id AND ind.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_market m ON m.market_id = o.market_id AND m.cust_id = '` + custId + `'
	LEFT JOIN mst.m_plu_group mpg ON mpg.plu_grp_id = o.plu_grp_id AND mpg.cust_id = '` + custId + `'
	LEFT JOIN mst.m_conv_group mcg ON mcg.conv_grp_id = o.conv_grp_id AND mcg.cust_id = '` + custId + `'
	LEFT JOIN mst.m_invoice_disc mid ON mid.inv_disc_id = o.disc_inv_id  AND mid.cust_id = '` + custId + `'

	LEFT JOIN mst.m_ward owrd on owrd.ward_id = o.outlet_ward_id AND owrd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sub_district omsd on omsd.sub_district_id = owrd.sub_district_id AND omsd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_regency omr on omr.regency_id = omsd.regency_id AND omr.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_province omp on omp.province_id = omr.province_id AND omp.cust_id = '` + parentCustId + `'
	
	LEFT JOIN mst.m_ward dlvwrd on dlvwrd.ward_id = o.delv_ward_id AND dlvwrd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sub_district dlvmsd on dlvmsd.sub_district_id = dlvwrd.sub_district_id AND dlvmsd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_regency dlvmr on dlvmr.regency_id = dlvmsd.regency_id AND dlvmr.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_province dlvmp on dlvmp.province_id = dlvmr.province_id AND dlvmp.cust_id = '` + parentCustId + `'

	LEFT JOIN mst.m_ward dlvwrd2 on dlvwrd2.ward_id = o.delv_ward_id2 AND dlvwrd2.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sub_district dlvmsd2 on dlvmsd2.sub_district_id = dlvwrd2.sub_district_id AND dlvmsd2.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_regency dlvmr2 on dlvmr2.regency_id = dlvmsd2.regency_id AND dlvmr2.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_province dlvmp2 on dlvmp2.province_id = dlvmr2.province_id AND dlvmp2.cust_id = '` + parentCustId + `'

	LEFT JOIN mst.m_ward invwrd on invwrd.ward_id = o.inv_ward_id AND invwrd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sub_district invmsd on invmsd.sub_district_id = invwrd.sub_district_id AND invmsd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_regency invmr on invmr.regency_id = invmsd.regency_id AND invmr.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_province invmp on invmp.province_id = invmr.province_id AND invmp.cust_id = '` + parentCustId + `'

			  WHERE o.outlet_id = $1 
			  AND o.cust_id = $2`
	err := repository.Get(&outlet, query, outletId, custId)
	if err != nil {
		log.Println("outletRepository, FindOneByOutletIdAndCustId, err:", err.Error())
		return outlet, err
	}

	return outlet, nil
}

func (repository *mOutletRepositoryImpl) GetDetailOutletContact(outletId int64, custId string) ([]model.MOutletContactReads, error) {
	outletContact := []model.MOutletContactReads{}
	query := `SELECT cust_id, outlet_id, contact_name, 
	job_title, phone_no, wa_no, is_wa_no, identity_type, identity_no, fax_number, outlet_establishment_date,
	email, outlet_contact_id from mst.m_outlet_contact
			  WHERE outlet_id = $1 
			  AND cust_id = $2`
	err := repository.Select(&outletContact, query, outletId, custId)
	if err != nil {
		log.Println("outletRepository, FindOneByOutletIdAndCustId, err:", err.Error())
		return outletContact, err
	}
	return outletContact, nil
}

func (repository *mOutletRepositoryImpl) GetOutletCountsByNames(outletNames []string) ([]model.MOutletFromListSimilar, error) {
	log.Println("outletRepositoryImpl, outletNames:", outletNames)
	outlet := []model.MOutletFromListSimilar{}

	// Convert []string to []any for sqlx.In
	args := make([]any, len(outletNames))
	for i, v := range outletNames {
		args[i] = v
	}

	query := `SELECT count(outlet_name) as count, outlet_name 
				FROM mst.m_outlet WHERE outlet_name IN (?)
				GROUP BY outlet_name ORDER BY count(DISTINCT outlet_name) DESC`

	// Use sqlx.In to handle the IN clause
	query, args, err := sqlx.In(query, args)
	if err != nil {
		log.Println("outletRepository, GetOutletCountsByNames, sqlx.In err:", err.Error())
		return outlet, err
	}

	// Rebind for PostgreSQL
	query = repository.Rebind(query)

	err = repository.Select(&outlet, query, args...)
	if err != nil {
		log.Println("outletRepository, GetOutletCountsByNames, err:", err.Error())
		return outlet, err
	}

	return outlet, nil
}

func (repository *mOutletRepositoryImpl) FindMobileOutletList(dataFilter entity.MobileOutletListQueryFilter, custId string, empID int64) ([]model.MobileOutletList, int64, int, error) {
	outlets := []model.MobileOutletList{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` o.outlet_id,
		COALESCE(o.outlet_code, '') AS outlet_code,
		COALESCE(o.outlet_name, '') AS outlet_name,
		COALESCE(o.outlet_status, 0) AS outlet_status,
		COALESCE(o.address1, '') AS address1,
		COALESCE(o.latitude, '') AS latitude,
		COALESCE(o.longitude, '') AS longitude,
		COALESCE(o.avg_sales_week, 0) AS avg_sales_week
	`
	qWhere := ` WHERE o.is_del = false AND o.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (o.outlet_code ILIKE '%` + dataFilter.Query + `%'
					OR o.outlet_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if len(dataFilter.OutletStatus) > 0 {
		ids := make([]string, len(dataFilter.OutletStatus))
		for i, id := range dataFilter.OutletStatus {
			ids[i] = strconv.Itoa(id)
		}
		qWhere += ` AND o.outlet_status IN (` + strings.Join(ids, ",") + `) `
	}

	if empID > 0 {
		qWhere += ` AND o.created_by = ` + strconv.FormatInt(empID, 10)
	}

	qWhere += ` AND o.verification_status= 2`

	qFrom := ` FROM mst.m_outlet o `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int64
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("outletRepository, FindMobileOutletList, count total, err:", err.Error())
		return outlets, 0, 0, err
	}

	sortBy := `outlet_code ASC` // default sort
	if dataFilter.Sort != "" {
		sortParts := strings.Split(dataFilter.Sort, ":")
		if len(sortParts) == 2 {
			sortField := sortParts[0]
			sortOrder := strings.ToUpper(sortParts[1])
			if sortOrder == "ASC" || sortOrder == "DESC" {
				sortBy = fmt.Sprintf(`%s %s`, sortField, sortOrder)
			}
		}
	}
	querySelect += fmt.Sprintf(` ORDER BY %s`, sortBy)

	if dataFilter.Limit <= 0 || dataFilter.Limit > 999 {
		dataFilter.Limit = 5
	}
	page := dataFilter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))
	if lastPage <= 0 {
		lastPage = 1
	}

	querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))

	err = repository.Select(&outlets, querySelect)
	if err != nil {
		log.Println("outletRepository, FindMobileOutletList, err:", err.Error())
		return outlets, 0, 0, err
	}

	return outlets, total, lastPage, nil
}

func (repository *mOutletRepositoryImpl) FindTopStatus(ctx context.Context, custID string) (int, error) {
	query := `
		SELECT
			moc.status
		FROM
			mst.m_outlet_config moc
		INNER JOIN mst.m_outlet_config_det mocd ON
			mocd.outlet_config_id = moc.outlet_config_id
			AND mocd.is_del = false
		INNER JOIN mst.m_outlet_config_status mocs ON 
		  mocs.outlet_config_status_id = mocd.status
		WHERE
			moc.is_del = false
		  AND moc.cust_id = $1
		  AND mocd.cust_id = $1
		ORDER BY mocs.sort_order ASC 
		LIMIT 1;`

	var status int
	err := repository.QueryRowContext(ctx, query, custID).Scan(&status)
	if err != nil {
		return status, err
	}

	return status, nil
}

func (repository *mOutletRepositoryImpl) FindRegionByDistributorID(ctx context.Context, distributorID int) (model.Region, error) {
	var region model.Region
	query := `
		SELECT m_region.*
		FROM mst.m_distributor
		INNER JOIN mst.m_region ON m_region.region_id = m_distributor.region_id
		WHERE m_distributor.distributor_id = $1`

	err := repository.GetContext(ctx, &region, query, distributorID)
	if err != nil {
		log.Println("mOutletRepositoryImpl, FindRegionByDistributorID, err:", err.Error())
		return region, err
	}

	return region, nil
}

func (repository *mOutletRepositoryImpl) CheckOutletByCode(ctx context.Context, outletCode string, custId string) (bool, error) {
	var count int
	query := `SELECT COUNT(outlet_id) FROM mst.m_outlet WHERE outlet_code = $1 AND cust_id = $2 AND is_del = false`

	err := repository.QueryRowContext(ctx, query, outletCode, custId).Scan(&count)
	if err != nil {
		log.Println("mOutletRepositoryImpl, CheckOutletByCode, err:", err.Error())
		return false, err
	}

	return count > 0, nil
}

func (repository *mOutletRepositoryImpl) FindByDestinations(ctx context.Context, destinationType string, destinationIDs []int) ([]model.MOutletReadFromList, error) {
	var query string

	if destinationType == "outlet" {
		query = `
				SELECT
						o.outlet_id,
						o.outlet_name,
						o.outlet_code,
						o.address1,
						o.cust_id,
						o.phone_no,
						o.wa_no,
						o.fax_no,
						o.latitude,
						o.longitude
				FROM mst.m_outlet o
				WHERE
					o.outlet_id in (?)
					AND o.is_del = false
				  	AND o.verification_status = 1
				  	AND o.outlet_status != 4;`
	} else {
		query = `
				SELECT
						md.distributor_id as outlet_id,
						md.distributor_name as outlet_name,
						md.distributor_code as outlet_code,
						md.address as address1,
						md.cust_id,
						md.phone as phone_no,
						md.fax_number as fax_no,
						md.latitude,
						md.longitude
				FROM mst.m_distributor md
				WHERE
					md.distributor_id in (?)
					AND md.is_active = true
				  	AND md.is_del = false;`
	}

	outlets := []model.MOutletReadFromList{}
	ids := make([]interface{}, len(destinationIDs))
	for i, v := range destinationIDs {
		ids[i] = v
	}

	query, args, err := sqlx.In(query, ids)
	if err != nil {
		log.Println("outletRepository, FindByDestinations, sqlx.In err:", err.Error())
		return nil, err
	}
	query = repository.Rebind(query)

	err = repository.SelectContext(ctx, &outlets, query, args...)
	if err != nil {
		log.Println("outletRepository, FindByDestinations, err:", err.Error())
		return nil, err
	}

	return outlets, nil
}
