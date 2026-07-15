package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"master/pkg/structs"

	"github.com/jmoiron/sqlx"
)

type outletRepositoryImpl struct {
	*sqlx.DB
}

// Store implements OutletRepository.
func (*outletRepositoryImpl) Store(outlet model.Outlet) (int, error) {
	panic("unimplemented")
}

type outletTransaction struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

// OutletTx is the transaction handle for outlet operations. Used so callers (e.g. service) can hold the transaction and call Tx() for same-tx operations.
type OutletTx interface {
	Store(outlet *model.Outlet) error
	StoreDetail(outletSalesman *model.MOutletSalesman) error
	StoreDetailSalesman(outletSalesman *model.MOutletSalesman) error
	StoreDetailBank(outletBank *model.MOutletBank) error
	StoreDetailContact(outletContact *model.MOutletContact) error
	StoreDetailTax(outletTax *model.MOutletTax) error
	Update(outletId int, request entity.UpdateOutletRequest) error
	UpdateDetailOutletSalesman(outletId int, outletSalesId int64, request entity.OutletSalesman) error
	UpdateDetailOutletBank(outletId int, outletBankId int64, request entity.OutletBank) error
	UpdateDetailOutletContact(outletId int, outletContactId int64, request entity.OutletContact) error
	UpdateDetailOutletContactPartial(outletId int, outletContactId int64, patch model.MOutletContactUpdate) error
	DeleteDetailSalesNotIn(outletId int64, outletSalesIds []int64) error
	DeleteDetailSalesByOutletId(outletId int64) error
	DeleteDetailBankNotIn(outletId int64, outletBankIds []int64) error
	DeleteDetailBankByOutletId(outletId int64) error
	DeleteDetailContactNotIn(outletId int64, outletContactIds []int64) error
	DeleteDetailContactByOutletId(outletId int64) error
	UpdateDetailOutletTax(outletId int, outletTaxId int64, request entity.OutletTax) error
	UpdateDetailOutletTaxPartial(outletId int, outletTaxId int64, patch model.MOutletTaxUpdate) error
	DeleteDetailTaxNotIn(outletId int64, outletTaxIds []int64) error
	DeleteDetailTaxByOutletId(outletId int64) error
	Approve(request entity.ApproveOutletBody) error
	Reject(request entity.RejectOutletBody) error
	GetOutletCrByOutletCrIds(outletCrIds []int) ([]model.OutletCr, error)
	GetOutletCrDetByOutletCrIds(outletCrIds []int) ([]model.OutletCrDet, error)
	UpdateOutletCrStatus(outletCrIds []int, status int, updatedBy int64) error
	UpdateOutletLocationFromCr(outletId int64, latitude, longitude string) error
	TrxCommit() error
	TrxRollback() error
	Tx() *sqlx.Tx
}

type OutletRepository interface {
	// Export
	FindAllExport(dataFilter entity.OutletQueryFilter, custId string) ([]model.OutletExport, int, error)
	TrxBegin() (OutletTx, error)
	FindOneParentCustId(distCustId string) (model.MCustomer, error)
	FindOneByOutletIdAndCustId(outletId int64, custId, parentCustId string) (model.OutletRead, error)
	FindOneByOutletCodeAndCustId(outletCode string, custId, parentCustId string) (model.OutletRead, error)
	FindAllByCustId(dataFilter entity.OutletQueryFilter, custId, parentCustId string) (consPro []model.OutletRead, total int, lastPage int, err error)
	FindAllOutletTypeByCustId(dataFilter entity.OutletQueryFilter, custId, parentCustId string) (consPro []model.OutletRead, total int, lastPage int, err error)
	FindAllOutletGroupByCustId(dataFilter entity.OutletQueryFilter, custId, parentCustId string) (consPro []model.OutletRead, total int, lastPage int, err error)
	FindAllVerificationStatusByCustId(dataFilter entity.OutletQueryFilter, custId, parentCustId string) (consPro []model.OutletRead, total int, lastPage int, err error)
	// Store(outlet model.Outlet) (int, error)
	// Update(outletId int, request entity.UpdateOutletRequest) error
	// BulkInsertOutletTemp(outlets []model.OutletTemp) error
	CheckBankExists(custId string, bankId int64) (bool, error)
	CheckBankCodeDuplicate(custId string, bankId int64, code string) (bool, error)
	CheckOutletClassExists(custId string, otClassId int64) (bool, error)
	CheckOutletClassCodeDuplicate(custId string, otClassId int64, code string) (bool, error)
	CheckOutletGroupExists(custId string, otGrpId int64) (bool, error)
	CheckOutletGroupCodeDuplicate(custId string, otGrpId int64, code string) (bool, error)
	CheckOutletLocExists(custId string, otLocId int64) (bool, error)
	CheckOutletLocCodeDuplicate(custId string, otLocId int64, code string) (bool, error)
	CheckOutletTypeExists(custId string, otTypeId int64) (bool, error)
	CheckOutletTypeCodeDuplicate(custId string, otTypeId int64, code string) (bool, error)
	CheckDistrictExists(custId string, districtId int64) (bool, error)
	CheckDistrictCodeDuplicate(custId string, districtId int64, code string) (bool, error)
	CheckDiscGroupExists(custId string, discGrpId int64) (bool, error)
	CheckDiscGroupCodeDuplicate(custId string, discGrpId int64, code string) (bool, error)
	CheckMarketExists(custId string, marketId int64) (bool, error)
	CheckMarketCodeDuplicate(custId string, marketId int64, code string) (bool, error)
	CheckIndustryExists(custId string, industryId int64) (bool, error)
	CheckIndustryCodeDuplicate(custId string, industryId int64, code string) (bool, error)
	CheckPriceGroupExists(custId string, priceGrpId int64) (bool, error)
	CheckPriceGroupCodeDuplicate(custId string, priceGrpId int64, code string) (bool, error)
	FindBankIdByCode(custId, code string) (int64, error)
	FindOutletClassIdByCode(custId, code string) (int64, error)
	FindOutletGroupIdByCode(custId, code string) (int64, error)
	FindOutletLocIdByCode(custId, code string) (int64, error)
	FindOutletTypeIdByCode(custId, code string) (int64, error)
	FindDistrictIdByCode(custId, code string) (int64, error)
	FindDiscGroupIdByCode(custId, code string) (int64, error)
	FindMarketIdByCode(custId, code string) (int64, error)
	FindIndustryIdByCode(custId, code string) (int64, error)
	FindPriceGroupIdByCode(custId, code string) (int64, error)
	FindBankIdByName(custId, name string) (int64, error)
	FindOutletClassIdByName(custId, name string) (int64, error)
	FindOutletGroupIdByName(custId, name string) (int64, error)
	FindOutletLocIdByName(custId, name string) (int64, error)
	FindOutletTypeIdByName(custId, name string) (int64, error)
	FindDistrictIdByName(custId, name string) (int64, error)
	FindDiscGroupIdByName(custId, name string) (int64, error)
	FindMarketIdByName(custId, name string) (int64, error)
	FindMarketsByIDs(custId string, ids []int64) (map[int64]model.Market, error)
	FindIndustryIdByName(custId, name string) (int64, error)
	FindPriceGroupIdByName(custId, name string) (int64, error)
	CheckOutletExists(custId, outletId string) (bool, error)
	UpdateImportBank(custId, code, name string, bankId int64) error
	UpdateImportOutletClass(custId, code, name string, otClassId int64) error
	UpdateImportOutletGroup(custId, code, name string, otGrpId int64) error
	UpdateImportOutletLoc(custId, code, name string, otLocId int64) error
	UpdateImportOutletType(custId, code, name string, otTypeId int64) error
	UpdateImportDistrict(custId, code, name string, districtId int64) error
	UpdateImportDiscGroup(custId, code, name string, discGrpId int64) error
	UpdateImportMarket(custId, code, name string, marketId int64) error
	UpdateImportIndustry(custId, code, name string, industryId int64) error
	UpdateImportPriceGroup(custId string, priceGrpId int64, code, name string) (int64, error)
	UpdateOutletImport(data entity.ProcessedUpdateOutlet) error
	// Update import for 4 location tables
	CheckProvinceExists(custId, provinceId string) (bool, error)
	CheckRegencyExists(custId, regencyId string) (bool, error)
	CheckSubDistrictExists(custId, subDistrictId string) (bool, error)
	CheckWardExists(custId, wardId string) (bool, error)
	UpdateImportProvince(custId, provinceId, province string) error
	UpdateImportRegency(custId, regencyId, regency, provinceId string) error
	UpdateImportSubDistrict(custId, subDistrictId, subDistrict, provinceId, regencyId string) error
	UpdateImportWard(custId, wardId, ward, provinceId, regencyId, subDistrictId string) error
	// Reupload insert helpers (outlet_temp)
	GetCtidsOutletInsertByHistory(historyId int64) ([]string, error)
	DeleteOutletTempByCtid(ctid string) error
	CountOutletTemp(historyId int64) (int, error)
	InsertOutletUpdateTemp(temp entity.ImportOutletUpdateTemp) error
	CreateOutlet(data entity.ProcessedOutlet) (int64, error)
	CreateOutletBank(bank model.MOutletBank) (int64, error)
	CreateOutletContact(contact model.MOutletContact) (int64, error)
	CreateOutletTax(tax model.MOutletTax) (int64, error)
	GetWardIdByName(custId, name string) (string, error)
	GetSubDistrictIdByName(custId, name string) (string, error)
	GetRegencyIdByName(custId, name string) (string, error)
	GetProvinceIdByName(custId, name string) (string, error)
	GetOutletDataForTemplateUpdate(custId string, fields []string) (map[string][][]string, error)
	GetOrCreateOutletType(custId, otTypeCode, otTypeName string) (int64, error)
	GetOrCreateOutletGroup(custId, otGrpCode, otGrpName string) (int64, error)
	GetOrCreatePriceGroup(custId, priceGrpCode, priceGrpName string) (int64, error)
	GetOrCreateBank(custId, bankCode, bankName string) (int64, error)
	GetOrCreateOutletClass(custId, otClassCode, otClassName string) (int64, error)
	GetOrCreateDistrict(custId, districtCode, districtName string) (int64, error)
	GetOrCreateOutletLoc(custId, otLocCode, otLocName string) (int64, error)
	GetOrCreateDiscGroup(custId, discGrpCode, discGrpName string) (int64, error)
	GetOrCreateMarket(custId, marketCode, marketName string) (int64, error)
	GetOrCreateIndustry(custId, industryCode, industryName string) (int64, error)
	CheckOutletCodeExists(custId, outletCode string) (bool, error)
	// Location mapping upserts used during outlet import
	UpsertProvince(custId, provinceId, province string, userId int64) error
	UpsertRegency(custId, regencyId, regency, provinceId string, userId int64) error
	UpsertSubDistrict(custId, subDistrictId, subDistrict, provinceId, regencyId string, userId int64) error
	UpsertWard(custId, wardId, ward, provinceId, regencyId, subDistrictId string, userId int64) error
	CreateOutletTemp(historyId int64, status, custId string, data entity.OutletTemp) error
	CreateImportHistory(typeUpload, fileName, custId string, uploadedBy int64, totalData int) (int64, error)
	UpdateImportHistoryOutlet(historyId int64, success, failed int, statusReupload bool) error
	// GetOrCreateOutletUdf(custId string, outletId, udfId int64, udfValue string) (int64, error)
	Delete(custId string, outletId int, deletedBy int64) error
	FindDetailByOutletIdAndCustId(outleId int, custId string) ([]model.MOutletSalesman, error)
	GetDetailOutletSalesman(outletId int64, custId string, lang string) ([]model.MOutletSalesmanRead, error)
	GetDetailOutletbank(outletId int64, custId string) ([]model.MOutletBankRead, error)
	GetDetailOutletContact(outletId int64, custId string) ([]model.MOutletContactRead, error)
	GetDetailOutletTax(outletId int64, custId string) ([]model.MOutletTaxRead, error)
	FindAllByPriceGrpIDAndCustID(priceGrpID int, custID, parentCustID string) ([]model.Outlet, error)
	// Instructions for templates
	GetImportInstructions(instructionType string) ([]entity.ImportInstruction, error)

	// Temp maintenance for reupload
	DeleteOutletUpdateTempByCtid(ctid string) error
	CountOutletUpdateTemp(historyId int64) (int, error)
	GetImportTotalData(historyId int64) (int, error)
	// Resolve outlet_id by detail IDs for import-update
	FindOutletIdByOutletBankId(custId string, outletBankId int64) (int64, error)
	FindOutletIdByOutletContactId(custId string, outletContactId int64) (int64, error)
	FindOutletIdByOutletTaxId(custId string, outletTaxId int64) (int64, error)
	FindOneCustomerByDistributorID(distributorID int64) (model.MCustomer, error)
	FindCustIdsByDistributorIds(parentCustId string, distributorIds []int) ([]string, error)
	FindCustIdsByParentCustId(parentCustId string) ([]string, error)
	GetRegionAndDistributorCodeByCustId(custId, parentCustId string) (regionCode *string, distributorCode *string, err error)
	GetNextOutletPrincipalCodeSeq(prefix string) (seq int, err error)
	GetNextOutletPrincipalCodeSeqTx(tx *sqlx.Tx, prefix string) (seq int, err error)
	CreateOutletTx(tx *sqlx.Tx, data entity.ProcessedOutlet) (int64, error)
	BulkUpdateStatuses(ctx context.Context) (int64, error)
	BulkPromoteRegisteredWithTransToNoo(ctx context.Context) (int64, error)
	ExistsOutletContactIdentity(custId, parentCustId, identityType, identityNo string, excludeOutletId int64) (bool, error)
	UpdateOutletStatus(outletId int64, custId string, status int, updatedBy int64) error
	FindAllOutletCrByStatus(dataFilter entity.OutletListApprovalQueryFilter, custId string) ([]model.OutletCrList, int, int, error)
	UpdateOutletCrStatus(outletCrIds []int, status int, updatedBy int64) error
	GetOutletCrDetByOutletCrIds(outletCrIds []int) ([]model.OutletCrDet, error)
	GetOutletCrByOutletCrIds(outletCrIds []int) ([]model.OutletCr, error)
	UpdateOutletLocationFromCr(outletId int64, latitude, longitude string) error
}

func NewOutletRepository(db *sqlx.DB) *outletRepositoryImpl {
	return &outletRepositoryImpl{db}
}

func NewOutletTransaction(db *sqlx.DB) (trxObj *outletTransaction, err error) {
	trx := db.MustBegin()

	return &outletTransaction{tx: trx, db: db}, nil
}

func (repo *outletRepositoryImpl) TrxBegin() (OutletTx, error) {
	trxObj, err := NewOutletTransaction(repo.DB)
	if err != nil {
		return nil, err
	}
	return trxObj, nil
}

func (repo *outletTransaction) TrxCommit() error {
	return repo.tx.Commit()
}

func (repo *outletTransaction) TrxRollback() error {
	return repo.tx.Rollback()
}

// Tx returns the underlying *sqlx.Tx for use in same-transaction operations (e.g. outlet_code update after outlet insert).
func (repo *outletTransaction) Tx() *sqlx.Tx {
	return repo.tx
}

// UpdateOutletCrStatus updates the status of multiple outlet change requests within a transaction.
// It sets the status, updated_by, updated_at, approval_by, and approval_at fields.
// Uses parameterized query with sqlx.In to prevent SQL injection.
func (repository *outletTransaction) UpdateOutletCrStatus(outletCrIds []int, status int, updatedBy int64) error {
	if len(outletCrIds) == 0 {
		return errors.New("outlet_cr_id is required")
	}

	query, args, err := sqlx.In(
		`UPDATE mst.outlet_cr 
		SET status = ?, updated_by = ?, updated_at = CURRENT_TIMESTAMP, approval_by = ?, approval_at = CURRENT_TIMESTAMP 
		WHERE outlet_cr_id IN (?)`,
		status, updatedBy, updatedBy, outletCrIds,
	)
	if err != nil {
		return err
	}

	query = repository.db.Rebind(query)
	_, err = repository.tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// GetOutletCrDetByOutletCrIds retrieves outlet change request details for latitude and longitude fields
// within a transaction. Returns details filtered by outlet_cr_id and field_name (latitude/longitude).
// Uses parameterized query with sqlx.In to prevent SQL injection.
func (repository *outletTransaction) GetOutletCrDetByOutletCrIds(outletCrIds []int) ([]model.OutletCrDet, error) {
	if len(outletCrIds) == 0 {
		return []model.OutletCrDet{}, nil
	}

	query, args, err := sqlx.In(
		`SELECT outlet_cr_det_id, outlet_cr_id, field_name, old_value, new_value
		FROM mst.outlet_cr_det
		WHERE outlet_cr_id IN (?)
		AND field_name IN (?)`,
		outletCrIds, []string{"latitude", "longitude"},
	)
	if err != nil {
		return nil, err
	}

	query = repository.db.Rebind(query)
	var details []model.OutletCrDet
	err = repository.tx.Select(&details, query, args...)
	if err != nil {
		return details, err
	}

	return details, nil
}

// GetOutletCrByOutletCrIds retrieves outlet change request records by their IDs within a transaction.
// Returns all fields including cust_id, status, approval information, and timestamps.
// Uses parameterized query with sqlx.In to prevent SQL injection.
func (repository *outletTransaction) GetOutletCrByOutletCrIds(outletCrIds []int) ([]model.OutletCr, error) {
	if len(outletCrIds) == 0 {
		return []model.OutletCr{}, nil
	}

	query, args, err := sqlx.In(
		`SELECT cust_id, outlet_cr_id, outlet_id, source, status, created_by, created_at, updated_by, updated_at, approval_by, approval_at
		FROM mst.outlet_cr
		WHERE outlet_cr_id IN (?)`,
		outletCrIds,
	)
	if err != nil {
		return nil, err
	}

	query = repository.db.Rebind(query)
	var outletCrs []model.OutletCr
	err = repository.tx.Select(&outletCrs, query, args...)
	if err != nil {
		return outletCrs, err
	}

	return outletCrs, nil
}

// UpdateOutletLocationFromCr updates the latitude and longitude of an outlet based on approved change request
// within a transaction. Only updates non-deleted outlets (is_del = false).
func (repository *outletTransaction) UpdateOutletLocationFromCr(outletId int64, latitude, longitude string) error {
	query := `UPDATE mst.m_outlet 
		SET latitude = $1, longitude = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE outlet_id = $3 AND is_del = false`

	_, err := repository.tx.Exec(query, latitude, longitude, outletId)
	if err != nil {
		return err
	}

	return nil
}

func (repository *outletRepositoryImpl) FindOneParentCustId(distCustId string) (model.MCustomer, error) {
	mCustomer := model.MCustomer{}
	query := `SELECT 
				cust_id, cust_name, parent_cust_id
			  FROM smc.m_customer
			  WHERE cust_id = $1`
	err := repository.Get(&mCustomer, query, distCustId)
	if err != nil {
		log.Println("salesTeamRepository, FindOneParentCustId, err:", err.Error())
		return mCustomer, err
	}

	return mCustomer, nil
}

func outletExportDistributorJoin(distributorIDs []int) string {
	if len(distributorIDs) == 0 {
		return ""
	}
	return `LEFT JOIN smc.m_customer mc ON mc.cust_id = o.cust_id`
}

func outletExportDistributorWhere(distributorIDs []int) string {
	if len(distributorIDs) == 0 {
		return ""
	}
	parts := make([]string, 0, len(distributorIDs))
	for _, id := range distributorIDs {
		parts = append(parts, fmt.Sprintf("%d", id))
	}
	return ` AND mc.distributor_id IN (` + strings.Join(parts, ",") + `)`
}

func (repository *outletRepositoryImpl) FindAllExport(dataFilter entity.OutletQueryFilter, custId string) ([]model.OutletExport, int, error) {
	outlets := []model.OutletExport{}

	parentCustId := custId
	if parent, errParent := repository.FindOneParentCustId(custId); errParent == nil && parent.ParentCustId != "" {
		parentCustId = parent.ParentCustId
	}

	distributorIDsForMc := dataFilter.DistributorID
	if len(dataFilter.ResolvedCustIdsForDistributor) > 0 {
		distributorIDsForMc = nil
	}

	selectCount := ` COUNT(DISTINCT o.outlet_id) AS total `

	selectField := `
		o.cust_id, o.outlet_id, o.outlet_code, o.outlet_name, o.outlet_principal_code, o.barcode, o.outlet_status, o.market_id,
		o.address1, o.address2, o.city, o.zip_code, o.phone_no, o.wa_no, o.fax_no, o.email,
		o.top, o.payment_type, o.is_contra_bon, o.agent_from, o.credit_limit_type, o.credit_limit,
		o.sales_inv_limit_type, o.sales_inv_limit, o.avg_sales_week, o.avg_sales_month,
		o.first_trans_date, o.last_trans_date, o.first_week_no, o.ot_start_date, o.ot_reg_date,
		o.building_own, o.dob, o.ar_status, o.ar_total, o.closed_date, o.is_emb_bail,
		o.owner_name, o.owner_addr1, o.owner_addr2, o.owner_city, o.owner_phone_no, o.owner_id_no,
		o.delv_addr1, o.delv_addr2, o.delv_city, o.inv_addr1, o.inv_addr2, o.inv_city,
		o.is_active, o.created_by, o.created_at, o.updated_by, o.updated_at, o.is_del, 
		o.latitude, o.longitude, o.image_url, o.file_url,
		o.is_obs, o.obs, o.is_wa_no, o.delv_zip_code, o.delv_is_same_addr,
		o.inv_zip_code, o.inv_is_same_addr,
		o.verification_status, o.verified_at, o.verified_by, o.tax_invoice_form, o.obs_type,
		o.credit_limit_action, o.sales_inv_limit_action, o.obs_limit_action,
		TO_CHAR(o.outlet_establishment_date, 'YYYY-MM-DD') AS outlet_establishment_date,
		o.delv_city2, o.delv_latitude, o.delv_longitude, o.delv_latitude2, o.delv_longitude2, o.delv_zip_code2,
		o.delv_ward_id, dw.sub_district_id AS delv_sub_district_id, dsd.regency_id AS delv_regency_id, drg.province_id AS delv_province_id,
		o.inv_ward_id, iw.sub_district_id AS inv_sub_district_id, isd.regency_id AS inv_regency_id, irg.province_id AS inv_province_id,

		-- master join
		oc.ot_class_code, oc.ot_class_name,
		og.ot_grp_code, og.ot_grp_name,
		ol.ot_loc_code, ol.ot_loc_name,
		ot.ot_type_code, ot.ot_type_name,
		dc.district_code, dc.district_name,
		dg.disc_grp_code, dg.disc_grp_name,
		mk.market_code, mk.market_name,
		ind.industry_code, ind.industry_name,
		pg.sp_price_grp_code AS price_grp_code, pg.sp_price_grp_name AS price_grp_name,
	CAST(bank_detail.bank_id AS TEXT) AS bank_id, bank_detail.bank_code, bank_detail.bank_name,
	w.ward AS outlet_ward,
	COALESCE(osd.sub_district, sd.sub_district) AS outlet_sub_district,
	COALESCE(org.regency, rg.regency) AS outlet_regency,
	COALESCE(opv.province, pv.province) AS outlet_province,
	dw.ward AS delv_ward,
	dsd.sub_district AS delv_sub_district,
	drg.regency AS delv_regency,
	dpv.province AS delv_province,
	iw.ward AS inv_ward,
	isd.sub_district AS inv_sub_district,
	irg.regency AS inv_regency,
	ipv.province AS inv_province,

		-- tax
		tax_detail.tax_invoice_id, tax_detail.tax_no, tax_detail.tax_name, tax_detail.tax_city, 
		tax_detail.tax_addr1, tax_detail.tax_addr2, tax_detail.tax_type, tax_detail.nitku, tax_detail.address_tax, 
		tax_detail.tax_identifier_type, tax_detail.tax_identifier_no,

	-- contact (only mapped fields)
	contact_detail.contact_name, contact_detail.job_title,
	contact_detail.identity_no, contact_detail.identity_type,
	contact_detail.phone_no AS contact_phone_no,
	contact_detail.wa_no AS contact_wa_no,
	contact_detail.email AS contact_email,
	contact_detail.is_wa_no AS contact_is_wa_no,

		-- bank
	bank_detail.account_no, bank_detail.account_name,

		u.user_fullname AS updated_by_name
	`

	qFrom := `
		FROM mst.m_outlet o
		` + outletExportDistributorJoin(distributorIDsForMc) + `
		LEFT JOIN mst.m_outlet_class oc ON oc.ot_class_id = o.ot_class_id AND oc.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_outlet_group og ON og.ot_grp_id = o.ot_grp_id AND og.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_outlet_loc   ol ON ol.ot_loc_id = o.ot_loc_id AND ol.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_outlet_type  ot ON ot.ot_type_id = o.ot_type_id AND ot.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_district     dc ON dc.district_id = o.district_id AND dc.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_disc_group   dg ON dg.disc_grp_id = o.disc_grp_id AND dg.cust_id = '` + custId + `'
		LEFT JOIN mst.m_market       mk ON mk.market_id = o.market_id AND mk.cust_id = '` + custId + `'
		LEFT JOIN mst.m_industry     ind ON ind.industry_id = o.industry_id AND ind.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_sp_price_group  pg ON pg.sp_price_grp_id = o.price_grp_id AND pg.cust_id = '` + custId + `'
	-- bank linked through outlet_bank -> bank
	LEFT JOIN (
		SELECT 
			mob.cust_id,
			mob.outlet_id,
			mob.bank_id,
			mob.account_no,
			mob.account_name,
			b.bank_code,
			b.bank_name,
			ROW_NUMBER() OVER (PARTITION BY mob.cust_id, mob.outlet_id ORDER BY mob.outlet_bank_id DESC) AS bank_rn
		FROM mst.m_outlet_bank mob
		LEFT JOIN mst.m_bank b ON b.bank_id = mob.bank_id AND b.cust_id = mob.cust_id
	) AS bank_detail ON bank_detail.cust_id = o.cust_id AND bank_detail.outlet_id = o.outlet_id AND bank_detail.bank_rn = 1
		LEFT JOIN mst.m_ward w          ON w.ward_id = o.outlet_ward_id AND w.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_sub_district sd ON sd.sub_district_id = w.sub_district_id AND sd.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_regency rg      ON rg.regency_id = sd.regency_id AND rg.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_province pv     ON pv.province_id = rg.province_id AND pv.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_province opv    ON opv.province_id = o.outlet_province_id AND opv.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_regency org     ON org.regency_id = o.outlet_regency_id AND org.cust_id = '` + parentCustId + `'
		LEFT JOIN mst.m_sub_district osd ON osd.sub_district_id = o.outlet_sub_district_id AND osd.cust_id = '` + parentCustId + `'

	-- detail tables
	LEFT JOIN (
		SELECT 
			c.cust_id,
			c.outlet_id,
			c.contact_name,
			c.job_title,
			c.identity_no,
			c.identity_type,
			c.phone_no,
			c.wa_no,
			c.email,
			c.is_wa_no,
			ROW_NUMBER() OVER (PARTITION BY c.cust_id, c.outlet_id ORDER BY c.outlet_contact_id DESC) AS contact_rn
		FROM mst.m_outlet_contact c
	) AS contact_detail ON contact_detail.cust_id = o.cust_id AND contact_detail.outlet_id = o.outlet_id AND contact_detail.contact_rn = 1
	LEFT JOIN (
		SELECT 
			t.cust_id,
			t.outlet_id,
			t.tax_invoice_id,
			t.tax_no,
			t.tax_name,
			t.tax_city,
			t.tax_addr1,
			t.tax_addr2,
			t.tax_type,
			t.nitku,
			t.address_tax,
			t.tax_identifier_type,
			t.tax_identifier_no,
			ROW_NUMBER() OVER (PARTITION BY t.cust_id, t.outlet_id ORDER BY t.outlet_tax_id DESC) AS tax_rn
		FROM mst.m_outlet_tax t
	) AS tax_detail ON tax_detail.cust_id = o.cust_id AND tax_detail.outlet_id = o.outlet_id AND tax_detail.tax_rn = 1

	LEFT JOIN sys.m_user u ON u.user_id = o.updated_by
	LEFT JOIN mst.m_ward dw          ON dw.ward_id = o.delv_ward_id AND dw.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sub_district dsd ON dsd.sub_district_id = dw.sub_district_id AND dsd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_regency drg      ON drg.regency_id = dsd.regency_id AND drg.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_province dpv     ON dpv.province_id = drg.province_id AND dpv.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_ward iw          ON iw.ward_id = o.inv_ward_id AND iw.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sub_district isd ON isd.sub_district_id = iw.sub_district_id AND isd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_regency irg      ON irg.regency_id = isd.regency_id AND irg.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_province ipv     ON ipv.province_id = irg.province_id AND ipv.cust_id = '` + parentCustId + `'
	`

	qWhere := ` WHERE o.is_del = false `
	if len(dataFilter.ResolvedCustIdsForDistributor) > 0 {
		quoted := make([]string, len(dataFilter.ResolvedCustIdsForDistributor))
		for i, cid := range dataFilter.ResolvedCustIdsForDistributor {
			quoted[i] = `'` + strings.ReplaceAll(cid, `'`, `''`) + `'`
		}
		qWhere += ` AND o.cust_id IN (` + strings.Join(quoted, ",") + `) `
	} else {
		qWhere += ` AND o.cust_id = '` + custId + `' `
	}

	if dataFilter.Query != "" {
		qWhere += ` AND (o.outlet_code ILIKE '%` + dataFilter.Query + `%' 
					OR o.outlet_name ILIKE '%` + dataFilter.Query + `%') `
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND o.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND o.is_active = false `
		}
	}

	if strings.TrimSpace(dataFilter.Status) != "" {
		statusStrings := strings.Split(dataFilter.Status, ",")
		statusValues := make([]int, 0, len(statusStrings))
		for _, s := range statusStrings {
			trimmed := strings.TrimSpace(s)
			if trimmed == "" {
				continue
			}
			if trimmed == "0" || strings.EqualFold(trimmed, "all") {
				statusValues = nil
				break
			}
			if val, err := strconv.Atoi(trimmed); err == nil {
				if val != 0 {
					statusValues = append(statusValues, val)
				}
			}
		}

		if len(statusValues) == 1 {
			qWhere += ` AND o.outlet_status = ` + strconv.Itoa(statusValues[0])
		} else if len(statusValues) > 1 {
			parts := make([]string, len(statusValues))
			for i, v := range statusValues {
				parts[i] = strconv.Itoa(v)
			}
			qWhere += ` AND o.outlet_status IN (` + strings.Join(parts, ",") + `)`
		}
	} else if dataFilter.OutletStatus != nil && *dataFilter.OutletStatus != 0 {
		// fallback lama (single outlet_status)
		qWhere += ` AND o.outlet_status = ` + fmt.Sprintf("%d", *dataFilter.OutletStatus)
	}
	qWhere += outletExportDistributorWhere(distributorIDsForMc)

	// count query
	queryCount := `SELECT ` + selectCount + qFrom + qWhere

	// select query
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// hitung total
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("outletRepository, FindAllExport, count err:", err.Error())
		return outlets, 0, err
	}

	// Sorting
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		sortBy := ""
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		if sortBy != "" {
			querySelect += fmt.Sprintf(` ORDER BY %s`, sortBy)
		}
	} else {
		querySelect += ` ORDER BY o.outlet_id DESC`
	}

	// Tidak ada LIMIT OFFSET (export semua)
	err = repository.Select(&outlets, querySelect)
	if err != nil {
		log.Println("outletRepository, FindAllExport, err:", err.Error())
		return outlets, total, err
	}

	return outlets, total, nil
}

// Resolve outlet_id by outlet_bank_id within a cust_id
func (r *outletRepositoryImpl) FindOutletIdByOutletBankId(custId string, outletBankId int64) (int64, error) {
	var outletId int64
	q := `SELECT outlet_id FROM mst.m_outlet_bank WHERE cust_id = $1 AND outlet_bank_id = $2 LIMIT 1`
	if err := r.Get(&outletId, q, custId, outletBankId); err != nil {
		return 0, err
	}
	return outletId, nil
}

// Resolve outlet_id by outlet_contact_id within a cust_id
func (r *outletRepositoryImpl) FindOutletIdByOutletContactId(custId string, outletContactId int64) (int64, error) {
	var outletId int64
	q := `SELECT outlet_id FROM mst.m_outlet_contact WHERE cust_id = $1 AND outlet_contact_id = $2 LIMIT 1`
	if err := r.Get(&outletId, q, custId, outletContactId); err != nil {
		return 0, err
	}
	return outletId, nil
}

// Resolve outlet_id by outlet_tax_id within a cust_id
func (r *outletRepositoryImpl) FindOutletIdByOutletTaxId(custId string, outletTaxId int64) (int64, error) {
	var outletId int64
	q := `SELECT outlet_id FROM mst.m_outlet_tax WHERE cust_id = $1 AND outlet_tax_id = $2 LIMIT 1`
	if err := r.Get(&outletId, q, custId, outletTaxId); err != nil {
		return 0, err
	}
	return outletId, nil
}

func (repository *outletRepositoryImpl) GetImportInstructions(instructionType string) ([]entity.ImportInstruction, error) {
	rows := []entity.ImportInstruction{}

	query := `
		SELECT 
			instruction_id, 
			instruction_type, 
			kolom, 
			mandatory, 
			keterangan, 
			step
		FROM import.import_instructions
		WHERE instruction_type = $1
		ORDER BY 
			CASE 
				WHEN step ILIKE 'Step 1%' THEN 1
				WHEN step ILIKE 'Step 2%' THEN 2
				WHEN step ILIKE 'Step 3%' THEN 3
				ELSE 4
			END, 
			instruction_id;
	`

	err := repository.Select(&rows, query, instructionType)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (repository *outletRepositoryImpl) FindOneByOutletIdAndCustId(outletId int64, custId, parentCustId string) (model.OutletRead, error) {
	log.Println("outletRepositoryImpl, outletId:", outletId)
	outlet := model.OutletRead{}
	query := `SELECT o.*,
					CASE 
						WHEN o.payment_type = 1 THEN 'Cash On Delivery'
						WHEN o.payment_type = 2 THEN 'Cash Before Delivery'
						WHEN o.payment_type = 3 THEN 'Credit'
						ELSE ''
					END AS payment_type_name,
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
	COALESCE(NULLIF(TRIM(o.outlet_sub_district_id), ''), omsd.sub_district_id) as outlet_sub_district_id,
	COALESCE(osdm.sub_district, omsd.sub_district) as outlet_sub_district,
	COALESCE(NULLIF(TRIM(o.outlet_regency_id), ''), omr.regency_id) as outlet_regency_id,
	COALESCE(oreg.regency, omr.regency) as outlet_regency,
	COALESCE(NULLIF(TRIM(o.outlet_province_id), ''), omp.province_id) as outlet_province_id,
	COALESCE(oprov.province, omp.province) as outlet_province,
	
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
	LEFT JOIN mst.m_disc_group dg ON dg.disc_grp_id = o.disc_grp_id AND dg.cust_id = o.cust_id
	LEFT JOIN mst.m_outlet_loc ol ON ol.ot_loc_id = o.ot_loc_id AND ol.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_outlet_type ot ON ot.ot_type_id = o.ot_type_id AND ot.cust_id  = '` + parentCustId + `'
	LEFT JOIN mst.m_outlet_group og ON og.ot_grp_id = o.ot_grp_id AND og.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sp_price_group pg ON pg.sp_price_grp_id = o.price_grp_id AND pg.cust_id = o.cust_id
	LEFT JOIN mst.m_district dis ON dis.district_id = o.district_id AND dis.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_beat b ON b.beat_id = o.beat_id AND b.cust_id = o.cust_id
	LEFT JOIN mst.m_sub_beat sb ON sb.sbeat_id = o.sbeat_id AND sb.cust_id = o.cust_id
	LEFT JOIN mst.m_outlet_class oc ON oc.ot_class_id = o.ot_class_id AND oc.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_industry ind ON ind.industry_id = o.industry_id AND ind.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_market m ON m.market_id = o.market_id AND m.cust_id = o.cust_id
	LEFT JOIN mst.m_plu_group mpg ON mpg.plu_grp_id = o.plu_grp_id AND mpg.cust_id = o.cust_id
	LEFT JOIN mst.m_conv_group mcg ON mcg.conv_grp_id = o.conv_grp_id AND mcg.cust_id = o.cust_id
	LEFT JOIN mst.m_invoice_disc mid ON mid.inv_disc_id = o.disc_inv_id  AND mid.cust_id = o.cust_id

	LEFT JOIN mst.m_ward owrd on owrd.ward_id = o.outlet_ward_id AND owrd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sub_district omsd on omsd.sub_district_id = owrd.sub_district_id AND omsd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_regency omr on omr.regency_id = omsd.regency_id AND omr.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_province omp on omp.province_id = omr.province_id AND omp.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_province oprov ON oprov.province_id = o.outlet_province_id AND oprov.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_regency oreg ON oreg.regency_id = o.outlet_regency_id AND oreg.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sub_district osdm ON osdm.sub_district_id = o.outlet_sub_district_id AND osdm.cust_id = '` + parentCustId + `'
	
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
			  AND (o.cust_id = $2 OR ($2 = $3 AND o.cust_id IN (SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = $3)))`
	err := repository.Get(&outlet, query, outletId, custId, parentCustId)
	if err != nil {
		log.Println("outletRepository, FindOneByOutletIdAndCustId, err:", err.Error())
		return outlet, err
	}

	return outlet, nil
}

func (repo *outletRepositoryImpl) FindDetailByOutletIdAndCustId(outletId int, custId string) ([]model.MOutletSalesman, error) {
	outletSales := []model.MOutletSalesman{}
	query := `SELECT 
					cust_id, outlet_id, sales_id, w1,
					w2, w3, w4, route_id, day_id,
					outlet_sales_id
				FROM mst.m_outlet_salesman 
				WHERE outlet_id = $1 
				AND cust_id = $2`
	err := repo.Select(&outletSales, query, outletId, custId)
	if err != nil {
		log.Println("outletRepository, FindOneBySBrand1IdAndCustId, err:", err.Error())
		return outletSales, err
	}

	return outletSales, nil
}

func (repository *outletRepositoryImpl) FindOneByOutletCodeAndCustId(outletCode string, custId, parentCustId string) (model.OutletRead, error) {
	outlet := model.OutletRead{}
	query := `SELECT o.*,
	u.user_fullname AS updated_by_name,dg.disc_grp_code,dg.disc_grp_name,ol.ot_loc_code,ol.ot_loc_name,
	og.ot_grp_code,og.ot_grp_name,pg.price_grp_code,pg.price_grp_name,dis.district_code,dis.district_name,
	b.beat_code,b.beat_name,sb.sbeat_code,sb.sbeat_name,oc.ot_class_code,oc.ot_class_name,ind.industry_code,ind.industry_name,
	m.market_code,m.market_name,mpg.plu_grp_code,mpg.plu_grp_name,mcg.conv_grp_code,mcg.conv_grp_name,mid.inv_disc_code as disc_inv_code,mid.inv_disc_name as disc_inv_name FROM mst.m_outlet o
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
	WHERE o.outlet_code = $1 and o.is_del = FALSE AND o.cust_id = $2`
	err := repository.Get(&outlet, query, outletCode, custId)
	if err != nil {
		log.Println("outletRepository, FindOneByoutletCodeAndCustId, err:", err.Error())
		return outlet, err
	}

	return outlet, nil
}

func appendIntInFilter(qWhere string, column string, values []int) string {
	if len(values) == 0 {
		return qWhere
	}

	ids := make([]string, 0, len(values))
	for _, id := range values {
		ids = append(ids, strconv.Itoa(id))
	}

	return qWhere + ` and ` + column + ` IN (` + strings.Join(ids, ",") + `) `
}

func buildOutletCustScopeWhere(custId string, parentCustId string, resolvedCustIds []string) string {
	if len(resolvedCustIds) > 0 {
		quoted := make([]string, len(resolvedCustIds))
		for i, cid := range resolvedCustIds {
			quoted[i] = `'` + strings.ReplaceAll(cid, `'`, `''`) + `'`
		}
		return ` WHERE o.is_del = false AND o.cust_id IN (` + strings.Join(quoted, ",") + `) `
	}

	trimmedCustID := strings.TrimSpace(custId)
	return ` WHERE o.is_del = false and o.cust_id = '` + trimmedCustID + `' `
}

func (repository *outletRepositoryImpl) FindAllByCustId(dataFilter entity.OutletQueryFilter, custId, parentCustId string) ([]model.OutletRead, int, int, error) {

	outlets := []model.OutletRead{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` DISTINCT ON (o.outlet_id) o.*,
	u.user_fullname AS updated_by_name,u2.user_fullname AS verified_by_name,u3.user_fullname AS created_by_name,dg.disc_grp_code,dg.disc_grp_name,ol.ot_loc_code,ol.ot_loc_name,
	og.ot_grp_code,og.ot_grp_name,pg.price_grp_code,pg.price_grp_name,dis.district_code,dis.district_name,
	b.beat_code,b.beat_name,sb.sbeat_code,sb.sbeat_name,oc.ot_class_code,oc.ot_class_name,ind.industry_code,ind.industry_name,
	m.market_code,m.market_name,mpg.plu_grp_code,mpg.plu_grp_name,mcg.conv_grp_code,mcg.conv_grp_name,mid.inv_disc_code as disc_inv_code,mid.inv_disc_name as disc_inv_name,
	co.identity_type, co.identity_no,

	owrd.ward as outlet_ward,
	COALESCE(NULLIF(TRIM(o.outlet_sub_district_id), ''), omsd.sub_district_id) as outlet_sub_district_id,
	COALESCE(osdm.sub_district, omsd.sub_district) as outlet_sub_district,
	COALESCE(NULLIF(TRIM(o.outlet_regency_id), ''), omr.regency_id) as outlet_regency_id,
	COALESCE(oreg.regency, omr.regency) as outlet_regency,
	COALESCE(NULLIF(TRIM(o.outlet_province_id), ''), omp.province_id) as outlet_province_id,
	COALESCE(oprov.province, omp.province) as outlet_province,
	
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
	invmp.province as inv_province,invmp.province_id as inv_province_id,
	ot.ot_type_code,ot.ot_type_name,
	mot.nitku, mot.address_tax, mot.tax_identifier_no, mot.tax_identifier_type, mot.tax_type, mot.tax_invoice_id,
	mb.bank_code, mb.bank_name, mob.bank_id, mob.account_no, mob.account_name,
	moc.contact_name, moc.job_title, moc.phone_no AS contact_phone_no,
	ocs.status_code AS outlet_status_code, ocs.status_description AS outlet_status_desc
	`
	qWhere := buildOutletCustScopeWhere(custId, parentCustId, dataFilter.ResolvedCustIdsForDistributor)
	if dataFilter.Query != "" {
		qWhere += ` AND (o.outlet_code ILIKE '%` + dataFilter.Query + `%' 
					OR o.outlet_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if shouldApplyOutletIsActiveFilter(dataFilter.IncludeInactive) && dataFilter.IsActive != nil {
		// fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND o.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND o.is_active = false `
		}
	}

	if len(dataFilter.OutletStatusIDs) > 0 {
		qWhere = appendIntInFilter(qWhere, "o.outlet_status", dataFilter.OutletStatusIDs)
	} else if dataFilter.OutletStatus != nil && *dataFilter.OutletStatus != 0 {
		qWhere += ` AND o.outlet_status = ` + fmt.Sprintf("%d", *dataFilter.OutletStatus)
	}

	qWhere = appendIntInFilter(qWhere, "o.outlet_id", dataFilter.OutletID)
	qWhere = appendIntInFilter(qWhere, "o.verification_status", dataFilter.VerificationStatus)
	qWhere = appendIntInFilter(qWhere, "o.ot_class_id", dataFilter.OtClassID)
	qWhere = appendIntInFilter(qWhere, "o.ot_grp_id", dataFilter.OtGrpID)
	qWhere = appendIntInFilter(qWhere, "o.ot_type_id", dataFilter.OtTypeID)

	if len(dataFilter.IdentityType) > 0 && len(dataFilter.IdentityNo) > 0 {
		// ids := make([]string, len(dataFilter.VerificationStatus))
		// for i, id := range dataFilter.VerificationStatus {
		// 	ids[i] = strconv.Itoa(id)
		// }
		qWhere += ` and (co.identity_type IN ('` + strings.Join(dataFilter.IdentityType, "','") + `') and co.identity_no IN ('` + strings.Join(dataFilter.IdentityNo, "','") + `')) `
	}

	qFrom := ` FROM mst.m_outlet o
	LEFT JOIN sys.m_user u ON u.user_id = o.updated_by
	LEFT JOIN sys.m_user u2 ON u2.user_id = o.verified_by
	LEFT JOIN sys.m_user u3 ON u3.user_id = o.created_by
	LEFT JOIN mst.m_disc_group dg ON dg.disc_grp_id = o.disc_grp_id AND dg.cust_id = o.cust_id
	LEFT JOIN mst.m_outlet_contact co ON co.outlet_id = o.outlet_id AND co.cust_id = o.cust_id
	LEFT JOIN mst.m_outlet_loc ol ON ol.ot_loc_id = o.ot_loc_id AND ol.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_outlet_group og ON og.ot_grp_id = o.ot_grp_id AND og.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_price_group pg ON pg.price_grp_id = o.price_grp_id AND pg.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_district dis ON dis.district_id = o.district_id AND dis.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_beat b ON b.beat_id = o.beat_id AND b.cust_id = o.cust_id
	LEFT JOIN mst.m_sub_beat sb ON sb.sbeat_id = o.sbeat_id AND sb.cust_id = o.cust_id
	LEFT JOIN mst.m_outlet_class oc ON oc.ot_class_id = o.ot_class_id AND oc.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_industry ind ON ind.industry_id = o.industry_id AND ind.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_market m ON m.market_id = o.market_id AND m.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_plu_group mpg ON mpg.plu_grp_id = o.plu_grp_id AND mpg.cust_id = o.cust_id
	LEFT JOIN mst.m_conv_group mcg ON mcg.conv_grp_id = o.conv_grp_id AND mcg.cust_id = o.cust_id
	LEFT JOIN mst.m_invoice_disc mid ON mid.inv_disc_id = o.disc_inv_id  AND mid.cust_id = o.cust_id 
	
	LEFT JOIN mst.m_ward owrd on owrd.ward_id = o.outlet_ward_id AND owrd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sub_district omsd on omsd.sub_district_id = owrd.sub_district_id AND omsd.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_regency omr on omr.regency_id = omsd.regency_id AND omr.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_province omp on omp.province_id = omr.province_id AND omp.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_province oprov ON oprov.province_id = o.outlet_province_id AND oprov.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_regency oreg ON oreg.regency_id = o.outlet_regency_id AND oreg.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_sub_district osdm ON osdm.sub_district_id = o.outlet_sub_district_id AND osdm.cust_id = '` + parentCustId + `'
	
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
	LEFT JOIN mst.m_outlet_type ot on ot.ot_type_id = o.ot_type_id AND ot.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_outlet_tax mot on mot.outlet_id = o.outlet_id  and mot.cust_id = '` + parentCustId + `'
	Left JOIN mst.m_outlet_bank mob on mob.outlet_id = o.outlet_id and mob.cust_id = o.cust_id
	LEFT JOIN mst.m_bank mb on mb.bank_id = mob.bank_id and mb.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_outlet_contact moc on moc.outlet_id = o.outlet_id and moc.cust_id = o.cust_id
	LEFT JOIN mst.m_outlet_config_status ocs ON ocs.status_code = o.outlet_status::varchar AND ocs.is_del = false
	`
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	log.Println("outletRepository, queryCount:", queryCount)
	var total int
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
		querySelect += fmt.Sprintf(`ORDER BY o.outlet_id DESC, %s`, sortBy)
	} else {
		sortBy := `o.outlet_id`
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

	log.Println("outletRepository, querySelect:", querySelect)
	err = repository.Select(&outlets, querySelect)
	if err != nil {
		log.Println("outletRepository, FindAllByCustId, err:", err.Error())
		return outlets, total, lastPage, err
	}

	return outlets, total, lastPage, nil
}

func shouldApplyOutletIsActiveFilter(includeInactive *int) bool {
	return includeInactive == nil || *includeInactive != 1
}

func (repository *outletRepositoryImpl) FindAllOutletTypeByCustId(dataFilter entity.OutletQueryFilter, custId, parentCustId string) ([]model.OutletRead, int, int, error) {

	outlets := []model.OutletRead{}
	selectCount := ` COUNT(*) AS total `
	selectField := `a.ot_type_id, b.ot_type_code, b.ot_type_name  `
	qWhere := `where a.deleted_at is null and a.cust_id = '` + custId + `'`

	if dataFilter.Query != "" {
		qWhere += ` AND (b.ot_type_code ILIKE '%` + dataFilter.Query + `%' 
					OR b.ot_type_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if len(dataFilter.OtGrpID) > 0 {
		ids := make([]string, len(dataFilter.OtGrpID))
		for i, id := range dataFilter.OtGrpID {
			ids[i] = strconv.Itoa(id)
		}
		qWhere += ` and a.ot_grp_id IN (` + strings.Join(ids, ",") + `) `
	}

	if len(dataFilter.OutletID) > 0 {
		ids := make([]string, len(dataFilter.OutletID))
		for i, id := range dataFilter.OutletID {
			ids[i] = strconv.Itoa(id)
		}
		qWhere += ` and a.outlet_id IN (` + strings.Join(ids, ",") + `) `
	}
	qWhere += ` GROUP BY a.ot_type_id, b.ot_type_code, b.ot_type_name `

	qFrom := ` from mst.m_outlet a
			left join mst.m_outlet_type b on b.ot_type_id=a.ot_type_id
			left join mst.m_outlet_group c on c.ot_grp_id=a.ot_grp_id 
	`
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("outletRepository, queryCount:", queryCount)
	var total int
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
		querySelect += fmt.Sprintf(`ORDER BY a.%s`, sortBy)
	} else {
		sortBy := `ot_type_id`
		querySelect += fmt.Sprintf(`ORDER BY a.%s DESC`, sortBy)
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

func (repository *outletRepositoryImpl) FindAllOutletGroupByCustId(dataFilter entity.OutletQueryFilter, custId, parentCustId string) ([]model.OutletRead, int, int, error) {

	outlets := []model.OutletRead{}
	selectCount := ` COUNT(*) AS total `
	selectField := `a.ot_grp_id, c.ot_grp_code , c.ot_grp_name `
	qWhere := `where a.deleted_at is null and a.cust_id = '` + custId + `'`

	if dataFilter.Query != "" {
		qWhere += ` AND (c.ot_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR c.ot_grp_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if len(dataFilter.OtTypeID) > 0 {
		ids := make([]string, len(dataFilter.OtTypeID))
		for i, id := range dataFilter.OtTypeID {
			ids[i] = strconv.Itoa(id)
		}
		qWhere += ` and a.ot_type_id IN (` + strings.Join(ids, ",") + `) `
	}

	if len(dataFilter.OutletID) > 0 {
		ids := make([]string, len(dataFilter.OutletID))
		for i, id := range dataFilter.OutletID {
			ids[i] = strconv.Itoa(id)
		}
		qWhere += ` and a.outlet_id IN (` + strings.Join(ids, ",") + `) `
	}
	qWhere += ` GROUP BY a.ot_grp_id, c.ot_grp_code , c.ot_grp_name `

	qFrom := ` from mst.m_outlet a
			left join mst.m_outlet_type b on b.ot_type_id=a.ot_type_id
			left join mst.m_outlet_group c on c.ot_grp_id=a.ot_grp_id 
	`
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("outletRepository, queryCount:", queryCount)
	var total int
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
		querySelect += fmt.Sprintf(`ORDER BY a.%s`, sortBy)
	} else {
		sortBy := `ot_type_id`
		querySelect += fmt.Sprintf(`ORDER BY a.%s DESC`, sortBy)
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

func (repository *outletTransaction) Store(outlet *model.Outlet) error {
	query :=
		`INSERT INTO mst.m_outlet(
			cust_id, outlet_code, barcode, outlet_name, outlet_status, address1, address2, city, zip_code, phone_no, 
			wa_no, fax_no, email, disc_grp_id, ot_loc_id, ot_grp_id, price_grp_id, district_id, beat_id, sbeat_id, 
			ot_class_id, industry_id, market_id, top, payment_type, is_contra_bon, plu_grp_id, conv_grp_id, disc_inv_id, agent_from, 
			credit_limit_type, credit_limit_type_name, credit_limit, sales_inv_limit_type, sales_inv_limit_type_name, sales_inv_limit, avg_sales_week, avg_sales_month, first_trans_date, last_trans_date, first_week_no, ot_start_date, 
			ot_reg_date, building_own, dob, ar_status, ar_total, closed_date, is_emb_bail, is_pkp_outlet, tax_name, tax_addr1, tax_addr2, 
			tax_city, tax_no, tax_invoice_form, obs_type, owner_name, owner_addr1, owner_addr2, owner_city, owner_phone_no, owner_id_no, delv_addr1, 
			delv_addr2, delv_city, inv_addr1, inv_addr2, inv_city, is_active, created_by, created_at, updated_by, updated_at, is_del, deleted_by, deleted_at, image_url, latitude, longitude, ot_type_id, is_obs, obs,
			outlet_province_id, outlet_regency_id, outlet_sub_district_id, outlet_ward_id,is_wa_no,delv_ward_id,delv_zip_code,delv_is_same_addr,inv_ward_id,inv_zip_code,inv_is_same_addr,verification_status,outlet_establishment_date, credit_limit_action, credit_limit_action_name, sales_inv_limit_action, sales_inv_limit_action_name, obs_limit_action,
			delv_latitude, delv_longitude, delv_city2, delv_latitude2, delv_longitude2,delv_ward_id2,delv_zip_code2, outlet_principal_code)
		VALUES ( 
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
			$41, $42, $43, $44, $45, $46, $47, $48, $49, $50,
			$51, $52, $53, $54, $55, $56, $57, $58, $59, $60, 
			$61, $62, $63, $64, $65, $66, $67, $68, $69, $70,
			$71, $72, $73, $74, $75, $76, $77, $78, $79, $80,
			$81, $82, $83, $84, $85, $86, $87, $88, $89, $90,
			$91, $92, $93, $94, $95, $96, $97, $98, $99, $100,
			$101, $102, $103, $104, $105, $106, $107, $108, $109
		) RETURNING outlet_id;`
	var lastInsertId int64
	err := repository.tx.QueryRow(query,
		outlet.CustId, outlet.OutletCode, outlet.Barcode, outlet.OutletName, outlet.OutletStatus, outlet.Address1, outlet.Address2, outlet.City, outlet.ZipCode, outlet.PhoneNo,
		outlet.WaNo, outlet.FaxNo, outlet.Email, outlet.DiscGrpId, outlet.OtLocId, outlet.OtGrpId, outlet.PriceGrpId, outlet.DistrictId, outlet.BeatId, outlet.SbeatId,
		outlet.OtClassId, outlet.IndustryId, outlet.MarketId, outlet.Top, outlet.PaymentType, outlet.IsContraBon, outlet.PluGrpId, outlet.ConvGrpId, outlet.DiscInvId, outlet.AgentFrom,
		outlet.CreditLimitType, outlet.CreditLimitTypeName, outlet.CreditLimit, outlet.SalesInvLimitType, outlet.SalesInvLimitTypeName, outlet.SalesInvLimit, outlet.AvgSalesWeek, outlet.AvgSalesMonth, outlet.FirstTransDate, outlet.LastTransDate, outlet.FirstWeekNo, outlet.OtStartDate,
		outlet.OtRegDate, outlet.BuldingOwn, outlet.Dob, outlet.ArStatus, outlet.ArTotal, outlet.CloseDate, outlet.IsEmbBail, outlet.IsPkpOutlet, outlet.TaxName, outlet.TaxAddr1, outlet.TaxAddr2,
		outlet.TaxCity, outlet.TaxNo, outlet.TaxInvoiceForm, outlet.ObsType, outlet.OwnerName, outlet.OwnerAdd1, outlet.OwnerAddr2, outlet.OwnerCity, outlet.OwnerPhoneNo, outlet.OwnerIdNo, outlet.DelvAdd1,
		outlet.DelvAddr2, outlet.DelvCity, outlet.InvAddr1, outlet.InvAddr2, outlet.InvCity, outlet.IsActive, outlet.CreatedBy, outlet.CreatedAt, outlet.UpdatedBy,
		outlet.UpdatedAt, outlet.IsDel, outlet.DeletedBy, outlet.DeletedAt, outlet.ImageUrl, outlet.Latitude, outlet.Longitude, outlet.OtTypeId, outlet.IsObs, outlet.Obs,
		outlet.OutletProvinceId, outlet.OutletRegencyId, outlet.OutletSubDistrictId, outlet.OutletWardId, outlet.IsWaNo, outlet.DelvWardId, outlet.DelvZipCode, outlet.DelvIsSameAddress, outlet.InvWardId, outlet.InvZipCode, outlet.InvIsSameAddress, outlet.VerificationStatus, outlet.OutletEstablishmentDate, outlet.CreditLimitAction, outlet.CreditLimitActionName, outlet.SalesInvLimitAction, outlet.SalesInvLimitActionName, outlet.ObsLimitAction,
		outlet.DelvLatitude, outlet.DelvLongitude, outlet.DelvCity2, outlet.DelvLatitude2, outlet.DelvLongitude2, outlet.DelvWardId2, outlet.DelvZipCode2, outlet.OutletPrincipalCode).Scan(&lastInsertId)
	if err != nil {
		log.Println("outletRepository, Store, err:", err.Error())
		return err
	}
	outlet.OutletId = lastInsertId
	return nil
}

func (repository *outletRepositoryImpl) GetDetailOutletSalesman(outletId int64, custId string, lang string) ([]model.MOutletSalesmanRead, error) {
	outletSalesman := []model.MOutletSalesmanRead{}
	query := `SELECT 
			mos.cust_id, mos.outlet_id, mos.sales_id,sl.sales_name, mos.w1,
			mos.w2, mos.w3, mos.w4, mos.route_id, mos.day_id,d.day_name,mos.outlet_sales_id 
		FROM mst.m_outlet_salesman mos
		LEFT JOIN mst.m_salesman sl on mos.sales_id = sl.emp_id AND sl.cust_id = '` + custId + `'
		LEFT JOIN sys.m_day d on mos.day_id = d.day_id AND d.lang_id=$1
		WHERE mos.outlet_id = $2 
		AND mos.cust_id = $3`
	err := repository.Select(&outletSalesman, query, lang, outletId, custId)
	if err != nil {
		log.Println("outletRepository, FindOneByOutletIdAndCustId, err:", err.Error())
		return outletSalesman, err
	}
	return outletSalesman, nil
}

func (repository *outletRepositoryImpl) GetDetailOutletbank(outletId int64, custId string) ([]model.MOutletBankRead, error) {
	outletBank := []model.MOutletBankRead{}
	query := `SELECT mob.cust_id, mob.outlet_id, mob.bank_id, b.bank_code,
	          CASE WHEN mob.bank_id = 0 THEN 'Non Bank' ELSE b.bank_name END AS bank_name,
	          mob.account_no, mob.account_name,mob.outlet_bank_id from mst.m_outlet_bank mob
			  left join mst.m_bank b on mob.bank_id = b.bank_id
			  WHERE mob.outlet_id = $1 
			  AND mob.cust_id = $2`
	err := repository.Select(&outletBank, query, outletId, custId)
	if err != nil {
		log.Println("outletRepository, FindOneByOutletIdAndCustId, err:", err.Error())
		return outletBank, err
	}
	return outletBank, nil
}

func (repository *outletRepositoryImpl) GetDetailOutletContact(outletId int64, custId string) ([]model.MOutletContactRead, error) {
	outletContact := []model.MOutletContactRead{}
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

func (repository *outletRepositoryImpl) GetDetailOutletTax(outletId int64, custId string) ([]model.MOutletTaxRead, error) {
	outletTax := []model.MOutletTaxRead{}
	query := `SELECT cust_id, outlet_id, tax_invoice_id, 
	is_emb_bail, tax_no, tax_name,
	tax_city, tax_addr1, tax_addr2, outlet_tax_id, tax_type, nitku, address_tax, tax_identifier_type, tax_identifier_no 
	from mst.m_outlet_tax
			  WHERE outlet_id = $1 
			  AND cust_id = $2`
	err := repository.Select(&outletTax, query, outletId, custId)
	if err != nil {
		log.Println("outletRepository, FindOneByOutletIdAndCustId, err:", err.Error())
		return outletTax, err
	}
	return outletTax, nil
}

func (repository *outletTransaction) StoreDetail(outletSalesman *model.MOutletSalesman) error {
	query :=
		`INSERT INTO mst.m_outlet_salesman(
			cust_id, outlet_id, sales_id, w1,
			w2, w3, w4, route_id, day_id)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6,
			$7, $8, $9
		) RETURNING 2;`
	var lastInsertId int64
	err := repository.tx.QueryRow(query,
		outletSalesman.CustID, outletSalesman.OutletID, outletSalesman.SalesID, outletSalesman.W1,
		outletSalesman.W2, outletSalesman.W3, outletSalesman.W4, outletSalesman.RouteID, outletSalesman.DayID).Scan(&lastInsertId)
	if err != nil {
		log.Println("outletRepository, StoreDetail, err:", err.Error())
		return err
	}
	outletSalesman.OutletSalesId = &lastInsertId
	return nil
}

func (repository *outletTransaction) StoreDetailSalesman(outletSalesman *model.MOutletSalesman) error {
	query :=
		`INSERT INTO mst.m_outlet_salesman(
			cust_id, outlet_id, sales_id, w1,
			w2, w3, w4, route_id, day_id)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6,
			$7, $8, $9
		) RETURNING outlet_sales_id;`
	var lastInsertId int64
	err := repository.tx.QueryRow(query,
		outletSalesman.CustID, outletSalesman.OutletID, outletSalesman.SalesID, outletSalesman.W1,
		outletSalesman.W2, outletSalesman.W3, outletSalesman.W4, outletSalesman.RouteID, outletSalesman.DayID).Scan(&lastInsertId)
	if err != nil {
		log.Println("outletRepository, StoreDetailSalesman, err:", err.Error())
		return err
	}
	outletSalesman.OutletSalesId = &lastInsertId
	return nil
}

func (repository *outletTransaction) StoreDetailBank(outletBank *model.MOutletBank) error {
	query :=
		`INSERT INTO mst.m_outlet_bank(
			cust_id, outlet_id, bank_id, 
			account_no, account_name)
		VALUES ( 
			$1, $2, $3, 
			$4, $5
		) RETURNING outlet_bank_id;`
	var lastInsertId int64
	err := repository.tx.QueryRow(query,
		outletBank.CustID, outletBank.OutletID, outletBank.BankId,
		outletBank.AccountNo, outletBank.AccountName).Scan(&lastInsertId)
	if err != nil {
		log.Println("outletRepository, StoreDetailSalesman, err:", err.Error())
		return err
	}
	outletBank.OutletBankId = &lastInsertId
	return nil
}

func (repository *outletTransaction) StoreDetailContact(outletContact *model.MOutletContact) error {
	query :=
		`INSERT INTO mst.m_outlet_contact(
			cust_id, outlet_id, contact_name, 
			job_title, phone_no, wa_no,
			email, identity_no, is_wa_no, identity_type, fax_number)
		VALUES ( 
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING outlet_contact_id;`
	var lastInsertId int64
	err := repository.tx.QueryRow(query,
		outletContact.CustID, outletContact.OutletID, outletContact.ContactName,
		outletContact.JobTitle, outletContact.PhoneNo, outletContact.WaNo,
		outletContact.Email, outletContact.IdentityNo, outletContact.IsWaNo, outletContact.IdentityType, outletContact.FaxNumber).Scan(&lastInsertId)
	if err != nil {
		log.Println("outletRepository, StoreDetailContact, err:", err.Error())
		return err
	}
	outletContact.OutletContactId = &lastInsertId
	return nil
}

func (repository *outletTransaction) StoreDetailTax(outletTax *model.MOutletTax) error {
	query :=
		`INSERT INTO mst.m_outlet_tax(
			cust_id, outlet_id, tax_invoice_id, 
			is_emb_bail, tax_no, tax_name,
			tax_addr1, tax_addr2, tax_city, tax_type, nitku, address_tax, tax_identifier_type, tax_identifier_no)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13, $14
		) RETURNING outlet_tax_id;`
	var lastInsertId int64
	err := repository.tx.QueryRow(query,
		outletTax.CustID, outletTax.OutletID, outletTax.TaxInvoiceId,
		outletTax.IsEmbBail, outletTax.TaxNo, outletTax.TaxName,
		outletTax.TaxAddr1, outletTax.TaxAddr2, outletTax.TaxCity, outletTax.TaxType, outletTax.Nitku,
		outletTax.AdddressTax, outletTax.TaxIdentifierType, outletTax.TaxIdentifierNo).Scan(&lastInsertId)
	if err != nil {
		log.Println("outletRepository, StoreDetailTax, err:", err.Error())
		return err
	}
	outletTax.OutletTaxId = &lastInsertId
	return nil
}

func (repository *outletTransaction) Update(outletId int, request entity.UpdateOutletRequest) error {
	var (
		r            model.OutletUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)
	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("outletRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_outlet SET ` + sqlSetFields + `, 
				updated_at = CURRENT_TIMESTAMP 
				WHERE is_del = false AND cust_id = :cust_id AND outlet_id = :outlet_id_old;`

	// log.Println("outletRepository, Update, query:", query)

	sqlPatch.Args["outlet_id_old"] = outletId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("outletRepository, Update, err:", err.Error())
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

func (repository *outletTransaction) UpdateDetailOutletSalesman(outletId int, outletSalesId int64, request entity.OutletSalesman) error {
	var (
		r            model.MOutletSalesmanUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("tprRepository, Update, Fields & Args: %s\n", data)
	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")
	fmt.Println(sqlSetFields)
	query := `UPDATE mst.m_outlet_salesman
			  SET ` + sqlSetFields + `
			  WHERE  outlet_sales_id = :outlet_sales_id
			  AND outlet_id=:outlet_id;`

	// log.Println("tprRepository, Update, query:", query)
	sqlPatch.Args["outlet_sales_id"] = outletSalesId
	sqlPatch.Args["outlet_id"] = outletId

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("outletSalesmanRepository, Update, err:", err.Error())
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

func (repository *outletTransaction) UpdateDetailOutletBank(outletId int, outletBankId int64, request entity.OutletBank) error {
	var (
		r            model.MOutletBankUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("tprRepository, Update, Fields & Args: %s\n", data)
	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")
	fmt.Println(sqlSetFields)
	query := `UPDATE mst.m_outlet_bank
			  SET ` + sqlSetFields + `
			  WHERE  outlet_bank_id = :outlet_bank_id
			  AND outlet_id=:outlet_id;`

	// log.Println("tprRepository, Update, query:", query)
	sqlPatch.Args["outlet_bank_id"] = outletBankId
	sqlPatch.Args["outlet_id"] = outletId

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("outletBankRepository, Update, err:", err.Error())
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

func (repository *outletTransaction) UpdateDetailOutletContact(outletId int, outletContactId int64, request entity.OutletContact) error {
	var (
		r            model.MOutletContactUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("tprRepository, Update, Fields & Args: %s\n", data)
	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")
	// fmt.Println("contact ===>", sqlSetFields)
	query := `UPDATE mst.m_outlet_contact
			  SET ` + sqlSetFields + `
			  WHERE  outlet_contact_id = :outlet_contact_id
			  AND outlet_id=:outlet_id;`

	// log.Println("tprRepository, Update, query:", query)
	sqlPatch.Args["outlet_contact_id"] = outletContactId
	sqlPatch.Args["outlet_id"] = outletId

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("outletContactRepository, Update, err:", err.Error())
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

func (repository *outletTransaction) UpdateDetailOutletTax(outletId int, outletTaxId int64, request entity.OutletTax) error {
	var (
		r            model.MOutletTaxUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("tprRepository, Update, Fields & Args: %s\n", data)
	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")
	fmt.Println("OUTLET TAX ===>", sqlSetFields)
	query := `UPDATE mst.m_outlet_tax
			  SET ` + sqlSetFields + `
			  WHERE  outlet_tax_id = :outlet_tax_id
			  AND outlet_id=:outlet_id;`

	// log.Println("tprRepository, Update, query:", query)
	sqlPatch.Args["outlet_tax_id"] = outletTaxId
	sqlPatch.Args["outlet_id"] = outletId

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("outletTaxRepository, Update, err:", err.Error())
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

func (repository *outletTransaction) UpdateDetailOutletContactPartial(outletId int, outletContactId int64, patch model.MOutletContactUpdate) error {
	sqlPatch := sql_helper.SQLPatches(patch)
	if len(sqlPatch.Fields) == 0 {
		return nil
	}
	var sqlSetFields string
	for i := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")
	query := `UPDATE mst.m_outlet_contact
			  SET ` + sqlSetFields + `
			  WHERE  outlet_contact_id = :outlet_contact_id
			  AND outlet_id=:outlet_id;`
	sqlPatch.Args["outlet_contact_id"] = outletContactId
	sqlPatch.Args["outlet_id"] = outletId
	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("outletContactRepository, UpdatePartial, err:", err.Error())
		return err
	}
	nRows, err := result.RowsAffected()
	if err != nil {
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *outletTransaction) UpdateDetailOutletTaxPartial(outletId int, outletTaxId int64, patch model.MOutletTaxUpdate) error {
	sqlPatch := sql_helper.SQLPatches(patch)
	if len(sqlPatch.Fields) == 0 {
		return nil
	}
	var sqlSetFields string
	for i := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")
	query := `UPDATE mst.m_outlet_tax
			  SET ` + sqlSetFields + `
			  WHERE  outlet_tax_id = :outlet_tax_id
			  AND outlet_id=:outlet_id;`
	sqlPatch.Args["outlet_tax_id"] = outletTaxId
	sqlPatch.Args["outlet_id"] = outletId
	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("outletTaxRepository, UpdatePartial, err:", err.Error())
		return err
	}
	nRows, err := result.RowsAffected()
	if err != nil {
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *outletTransaction) DeleteDetailSalesNotIn(outleId int64, outletSalesId []int64) error {
	query, args, err := sqlx.In("DELETE FROM mst.m_outlet_salesman WHERE outlet_id = ? AND outlet_sales_id not in (?);", outleId, outletSalesId)
	if err != nil {
		return err
	}
	query = repository.db.Rebind(query)
	_, err = repository.tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repository *outletTransaction) DeleteDetailSalesByOutletId(outleId int64) error {
	query, args, err := sqlx.In("DELETE FROM mst.m_outlet_salesman WHERE outlet_id = ?;", outleId)
	if err != nil {
		return err
	}
	query = repository.db.Rebind(query)
	_, err = repository.tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repository *outletTransaction) DeleteDetailBankNotIn(outleId int64, outletBankId []int64) error {
	query, args, err := sqlx.In("DELETE FROM mst.m_outlet_bank WHERE outlet_id = ? AND outlet_bank_id not in (?);", outleId, outletBankId)
	if err != nil {
		return err
	}
	query = repository.db.Rebind(query)
	_, err = repository.tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repository *outletTransaction) DeleteDetailBankByOutletId(outleId int64) error {
	query, args, err := sqlx.In("DELETE FROM mst.m_outlet_bank WHERE outlet_id = ?;", outleId)
	if err != nil {
		return err
	}
	query = repository.db.Rebind(query)
	_, err = repository.tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repository *outletTransaction) DeleteDetailContactNotIn(outleId int64, outletContactId []int64) error {
	query, args, err := sqlx.In("DELETE FROM mst.m_outlet_contact WHERE outlet_id = ? AND outlet_contact_id not in (?);", outleId, outletContactId)
	if err != nil {
		return err
	}
	query = repository.db.Rebind(query)
	_, err = repository.tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repository *outletTransaction) DeleteDetailContactByOutletId(outleId int64) error {
	query, args, err := sqlx.In("DELETE FROM mst.m_outlet_contact WHERE outlet_id = ?;", outleId)
	if err != nil {
		return err
	}
	query = repository.db.Rebind(query)
	_, err = repository.tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repository *outletTransaction) DeleteDetailTaxNotIn(outleId int64, outletTaxId []int64) error {
	query, args, err := sqlx.In("DELETE FROM mst.m_outlet_tax WHERE outlet_id = ? AND outlet_tax_id not in (?);", outleId, outletTaxId)
	if err != nil {
		return err
	}
	query = repository.db.Rebind(query)
	_, err = repository.tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repository *outletTransaction) DeleteDetailTaxByOutletId(outleId int64) error {
	query, args, err := sqlx.In("DELETE FROM mst.m_outlet_tax WHERE outlet_id = ?;", outleId)
	if err != nil {
		return err
	}
	query = repository.db.Rebind(query)
	_, err = repository.tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repository *outletRepositoryImpl) Delete(custId string, outletId int, deletedBy int64) error {
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

func (repository *outletTransaction) Approve(request entity.ApproveOutletBody) error {
	var (
		r            model.OutletApprove
		sqlSetFields string
		nRows        int64
	)

	err := structs.Automapper(request, &r)
	if err != nil {
		return err
	}

	sqlPatch := sql_helper.SQLPatches(r)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")
	log.Println(sqlSetFields)

	query := `UPDATE mst.m_outlet SET ` + sqlSetFields + `, 
				verified_at = CURRENT_TIMESTAMP 
				WHERE is_del = false AND cust_id = :cust_id AND outlet_id in (` + strings.Trim(strings.Join(strings.Split(fmt.Sprint(request.Outlets), " "), ","), "[]") + `);`

	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("outletRepository, Approve, err:", err.Error())
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

func (repository *outletTransaction) Reject(request entity.RejectOutletBody) error {
	var (
		r            model.OutletReject
		sqlSetFields string
		nRows        int64
	)

	err := structs.Automapper(request, &r)
	if err != nil {
		return err
	}

	sqlPatch := sql_helper.SQLPatches(r)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")
	log.Println(sqlSetFields)

	query := `UPDATE mst.m_outlet SET ` + sqlSetFields + `, 
				verified_at = CURRENT_TIMESTAMP 
				WHERE is_del = false AND cust_id = :cust_id AND outlet_id in (` + strings.Trim(strings.Join(strings.Split(fmt.Sprint(request.Outlets), " "), ","), "[]") + `);`

	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("outletRepository, Reject, err:", err.Error())
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

func (repo *outletRepositoryImpl) FindAllByPriceGrpIDAndCustID(priceGrpID int, custID, parentCustID string) ([]model.Outlet, error) {
	outlets := []model.Outlet{}
	query := `SELECT o.cust_id, o.outlet_id, o.outlet_code, o.outlet_name, o.ot_type_id, mt.ot_type_name, o.ot_grp_id, mg.ot_grp_name
	FROM mst.m_outlet o 
	LEFT JOIN mst.m_outlet_type mt ON mt.ot_type_id = o.ot_type_id AND mt.cust_id = '` + custID + `'
	LEFT JOIN mst.m_outlet_group mg ON mg.ot_grp_id = o.ot_grp_id AND mg.cust_id = '` + custID + `'
	WHERE o.is_del = false AND o.is_active = true AND o.cust_id = $1 AND o.price_grp_id = $2
	ORDER BY o.outlet_id ASC`
	err := repo.Select(&outlets, query, custID, priceGrpID)
	if err != nil {
		log.Println("outletRepository, FindAllByPriceGrpIDAndCustID, err:", err.Error())
		return outlets, err
	}

	return outlets, nil
}

func (repository *outletRepositoryImpl) FindAllVerificationStatusByCustId(dataFilter entity.OutletQueryFilter, custId, parentCustId string) ([]model.OutletRead, int, int, error) {

	outlets := []model.OutletRead{}
	selectCount := ` COUNT(DISTINCT a.verification_status) AS total `
	selectField := `DISTINCT a.verification_status `
	qWhere := `where a.deleted_at is null and a.cust_id = '` + custId + `'`
	// if dataFilter.Query != "" {
	// 	qWhere += ` AND (b.ot_type_code ILIKE '%` + dataFilter.Query + `%'
	// 				OR b.ot_type_name ILIKE '%` + dataFilter.Query + `%' )`
	// }
	// if len(dataFilter.IdentityType) > 0 && len(dataFilter.IdentityNo) > 0 {
	// ids := make([]string, len(dataFilter.IdentityType))
	// for i, id := range dataFilter.IdentityType {
	// 	ids[i] = strconv.Itoa(id)
	// }
	// qWhere += ` and a.ot_grp_id IN (` + strings.Join(ids, ",") + `) `
	// }

	// if len(dataFilter.OutletID) > 0 {
	// 	ids := make([]string, len(dataFilter.OutletID))
	// 	for i, id := range dataFilter.OutletID {
	// 		ids[i] = strconv.Itoa(id)
	// 	}
	// 	qWhere += ` and a.outlet_id IN (` + strings.Join(ids, ",") + `) `
	// }
	// qWhere += ` GROUP BY a.verification_status `
	qFrom := ` from mst.m_outlet a `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere
	// log.Println("outletRepository, queryCount:", queryCount)
	var total int
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
		querySelect += fmt.Sprintf(`ORDER BY a.%s`, sortBy)
	} else {
		sortBy := `verification_status`
		querySelect += fmt.Sprintf(`ORDER BY a.%s ASC`, sortBy)
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

// ========================================================
// 1. Bulk Insert ke outlet_temp
// ========================================================
func (r *outletRepositoryImpl) BulkInsertOutletTemp(outlets []model.OutletTemp) error {
	query := `
		INSERT INTO upload.outlet_temp (
			history_id, cust_id, outlet_code, outlet_name,
			address1, address2, city,
			phone_no, wa_no, email,
			disc_grp_code, ot_class_code, ot_grp_code, ot_type_code,
			district_code, market_code, industry_code,
			status_insert
		) VALUES (
			:history_id, :cust_id, :outlet_code, :outlet_name,
			:address1, :address2, :city,
			:phone_no, :wa_no, :email,
			:disc_grp_code, :ot_class_code, :ot_grp_code, :ot_type_code,
			:district_code, :market_code, :industry_code,
			:status_insert
		) RETURNING cust_id;`

	var insertedIDs []string
	for _, row := range outlets {
		// Default value
		if row.StatusInsert == "" {
			row.StatusInsert = "true"
		}
		if row.HistoryId == "" {
			row.HistoryId = "0"
		}

		res, err := r.DB.NamedQuery(query, row)
		if err != nil {
			log.Println("BulkInsertOutletTemp error:", err)
			return err
		}
		defer res.Close()

		for res.Next() {
			var custId string
			if err := res.Scan(&custId); err != nil {
				return err
			}
			insertedIDs = append(insertedIDs, custId)
		}
	}

	return nil
}

// 		}
// 	}

// 	return nil
// }

// ========================================================
// 2. GetOrCreate master reference
// ========================================================
func (r *outletRepositoryImpl) GetOrCreateBank(custId, bankCode, bankName string) (int64, error) {
	var id int64
	query := `SELECT bank_id FROM mst.m_bank WHERE cust_id = $1 AND bank_code = $2 AND bank_name = $3`
	err := r.DB.Get(&id, query, custId, bankCode, bankName)
	if err == nil {
		return id, nil
	}

	query = `INSERT INTO mst.m_bank (cust_id, bank_code, bank_name) 
             VALUES ($1, $2, $3) RETURNING bank_id`
	err = r.DB.QueryRow(query, custId, bankCode, bankName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) GetOrCreateOutletClass(custId, otClassCode, otClassName string) (int64, error) {
	var id int64
	query := `SELECT ot_class_id FROM mst.m_outlet_class WHERE cust_id = $1 AND ot_class_code = $2`
	err := r.DB.Get(&id, query, custId, otClassCode)
	if err == nil {
		return id, nil
	}

	query = `INSERT INTO mst.m_outlet_class (cust_id, ot_class_code, ot_class_name, is_active, created_at) 
             VALUES ($1, $2, $3, true, NOW()) RETURNING ot_class_id`
	err = r.DB.QueryRow(query, custId, otClassCode, otClassName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) GetOrCreateDistrict(custId, districtCode, districtName string) (int64, error) {
	var id int64
	query := `SELECT district_id FROM mst.m_district WHERE cust_id = $1 AND district_code = $2 AND district_name = $3`
	err := r.DB.Get(&id, query, custId, districtCode, districtName)
	if err == nil {
		return id, nil
	}

	query = `INSERT INTO mst.m_district (cust_id, district_code, district_name) 
             VALUES ($1, $2, $3) RETURNING district_id`
	err = r.DB.QueryRow(query, custId, districtCode, districtName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) GetOrCreateOutletGroup(custId, otGrpCode, otGrpName string) (int64, error) {
	var id int64
	query := `SELECT ot_grp_id FROM mst.m_outlet_group WHERE cust_id = $1 AND ot_grp_code = $2`
	err := r.DB.Get(&id, query, custId, otGrpCode)
	if err == nil {
		return id, nil
	}

	query = `INSERT INTO mst.m_outlet_group (cust_id, ot_grp_code, ot_grp_name, is_active, created_at) 
             VALUES ($1, $2, $3, true, NOW()) RETURNING ot_grp_id`
	err = r.DB.QueryRow(query, custId, otGrpCode, otGrpName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) GetOrCreatePriceGroup(custId, priceGrpCode, priceGrpName string) (int64, error) {
	var id int64
	query := `SELECT sp_price_grp_id FROM mst.m_sp_price_group WHERE cust_id = $1 AND sp_price_grp_code = $2`
	err := r.DB.Get(&id, query, custId, priceGrpCode)
	if err == nil {
		return id, nil
	}

	query = `INSERT INTO mst.m_sp_price_group (cust_id, sp_price_grp_code, sp_price_grp_name, is_active, created_at) 
             VALUES ($1, $2, $3, true, NOW()) RETURNING sp_price_grp_id`
	err = r.DB.QueryRow(query, custId, priceGrpCode, priceGrpName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) GetOrCreateOutletLoc(custId, otLocCode, otLocName string) (int64, error) {
	var id int64
	query := `SELECT ot_loc_id FROM mst.m_outlet_loc WHERE cust_id = $1 AND ot_loc_code = $2`
	err := r.DB.Get(&id, query, custId, otLocCode)
	if err == nil {
		return id, nil
	}

	query = `INSERT INTO mst.m_outlet_loc (cust_id, ot_loc_code, ot_loc_name, is_active, created_at) 
             VALUES ($1, $2, $3, true, NOW()) RETURNING ot_loc_id`
	err = r.DB.QueryRow(query, custId, otLocCode, otLocName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) GetOrCreateDiscGroup(custId, discGrpCode, discGrpName string) (int64, error) {
	var id int64
	query := `SELECT disc_grp_id FROM mst.m_disc_group WHERE cust_id = $1 AND disc_grp_code = $2 AND disc_grp_name = $3`
	err := r.DB.Get(&id, query, custId, discGrpCode, discGrpName)
	if err == nil {
		return id, nil
	}

	query = `INSERT INTO mst.m_disc_group (cust_id, disc_grp_code, disc_grp_name) 
             VALUES ($1, $2, $3) RETURNING disc_grp_id`
	err = r.DB.QueryRow(query, custId, discGrpCode, discGrpName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) GetOrCreateMarket(custId, marketCode, marketName string) (int64, error) {
	var id int64
	query := `SELECT market_id FROM mst.m_market WHERE cust_id = $1 AND market_code = $2 AND market_name = $3`
	err := r.DB.Get(&id, query, custId, marketCode, marketName)
	if err == nil {
		return id, nil
	}

	query = `INSERT INTO mst.m_market (cust_id, market_code, market_name) 
             VALUES ($1, $2, $3) RETURNING market_id`
	err = r.DB.QueryRow(query, custId, marketCode, marketName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) GetOrCreateIndustry(custId, industryCode, industryName string) (int64, error) {
	var id int64
	query := `SELECT industry_id FROM mst.m_industry WHERE cust_id = $1 AND industry_code = $2 AND industry_name = $3`
	err := r.DB.Get(&id, query, custId, industryCode, industryName)
	if err == nil {
		return id, nil
	}

	query = `INSERT INTO mst.m_industry (cust_id, industry_code, industry_name) 
             VALUES ($1, $2, $3) RETURNING industry_id`
	err = r.DB.QueryRow(query, custId, industryCode, industryName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) CheckOutletCodeExists(custId, outletCode string) (bool, error) {
	code := strings.TrimSpace(outletCode)
	if code == "" {
		return false, nil
	}
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_outlet WHERE cust_id = $1 AND LOWER(outlet_code) = LOWER($2) AND is_del = FALSE`
	err := r.DB.Get(&exists, query, custId, code)
	return exists, err
}

// --- Location mapping upserts (province/regency/sub-district/ward) ---

func (r *outletRepositoryImpl) UpsertProvince(custId, provinceId, province string, userId int64) error {
	if provinceId == "" {
		return nil
	}
	name := strings.TrimSpace(province)
	if name == "" {
		return nil
	}

	// Check existence to avoid duplicate inserts
	var exists bool
	if err := r.DB.Get(&exists, `SELECT EXISTS(SELECT 1 FROM mst.m_province WHERE cust_id = $1 AND province_id = $2)`, custId, provinceId); err != nil {
		return err
	}

	if exists {
		// Update only when something changes or was inactive
		_, err := r.DB.Exec(`
            UPDATE mst.m_province
               SET province = $3,
                   is_active = TRUE,
                   updated_by = $4,
                   updated_at = NOW()
             WHERE cust_id = $1 AND province_id = $2
               AND (province IS DISTINCT FROM $3 OR is_active = FALSE)
        `, custId, provinceId, name, userId)
		return err
	}

	_, err := r.DB.Exec(`
        INSERT INTO mst.m_province (cust_id, province_id, province, is_active, created_by, created_at, updated_by, updated_at)
        VALUES ($1, $2, $3, TRUE, $4, NOW(), $4, NOW())
    `, custId, provinceId, name, userId)
	return err
}

func (r *outletRepositoryImpl) UpsertRegency(custId, regencyId, regency, provinceId string, userId int64) error {
	if regencyId == "" {
		return nil
	}
	name := strings.TrimSpace(regency)
	if name == "" {
		return nil
	}

	var exists bool
	if err := r.DB.Get(&exists, `SELECT EXISTS(SELECT 1 FROM mst.m_regency WHERE cust_id = $1 AND regency_id = $2)`, custId, regencyId); err != nil {
		return err
	}

	if exists {
		_, err := r.DB.Exec(`
            UPDATE mst.m_regency
               SET regency = $3,
                   province_id = NULLIF($4, ''),
                   is_active = TRUE,
                   updated_by = $5,
                   updated_at = NOW()
             WHERE cust_id = $1 AND regency_id = $2
               AND (regency IS DISTINCT FROM $3 OR province_id IS DISTINCT FROM NULLIF($4, '') OR is_active = FALSE)
        `, custId, regencyId, name, provinceId, userId)
		return err
	}

	_, err := r.DB.Exec(`
        INSERT INTO mst.m_regency (cust_id, regency_id, regency, province_id, is_active, created_by, created_at, updated_by, updated_at)
        VALUES ($1, $2, $3, NULLIF($4, ''), TRUE, $5, NOW(), $5, NOW())
    `, custId, regencyId, name, provinceId, userId)
	return err
}

func (r *outletRepositoryImpl) UpsertSubDistrict(custId, subDistrictId, subDistrict, provinceId, regencyId string, userId int64) error {
	if subDistrictId == "" {
		return nil
	}
	name := strings.TrimSpace(subDistrict)
	if name == "" {
		return nil
	}

	var exists bool
	if err := r.DB.Get(&exists, `SELECT EXISTS(SELECT 1 FROM mst.m_sub_district WHERE cust_id = $1 AND sub_district_id = $2)`, custId, subDistrictId); err != nil {
		return err
	}

	if exists {
		_, err := r.DB.Exec(`
            UPDATE mst.m_sub_district
               SET sub_district  = $3,
                   province_id   = NULLIF($4, ''),
                   regency_id    = NULLIF($5, ''),
                   is_active     = TRUE,
                   updated_by    = $6,
                   updated_at    = NOW()
             WHERE cust_id = $1 AND sub_district_id = $2
               AND (sub_district IS DISTINCT FROM $3 OR province_id IS DISTINCT FROM NULLIF($4, '') OR regency_id IS DISTINCT FROM NULLIF($5, '') OR is_active = FALSE)
        `, custId, subDistrictId, name, provinceId, regencyId, userId)
		return err
	}

	_, err := r.DB.Exec(`
        INSERT INTO mst.m_sub_district (cust_id, sub_district_id, sub_district, province_id, regency_id, is_active, created_by, created_at, updated_by, updated_at)
        VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), TRUE, $6, NOW(), $6, NOW())
    `, custId, subDistrictId, name, provinceId, regencyId, userId)
	return err
}

func (r *outletRepositoryImpl) UpsertWard(custId, wardId, ward, provinceId, regencyId, subDistrictId string, userId int64) error {
	if wardId == "" {
		return nil
	}
	name := strings.TrimSpace(ward)
	if name == "" {
		return nil
	}

	var exists bool
	if err := r.DB.Get(&exists, `SELECT EXISTS(SELECT 1 FROM mst.m_ward WHERE cust_id = $1 AND ward_id = $2)`, custId, wardId); err != nil {
		return err
	}

	if exists {
		_, err := r.DB.Exec(`
            UPDATE mst.m_ward
               SET ward            = $3,
                   province_id     = NULLIF($4, ''),
                   regency_id      = NULLIF($5, ''),
                   sub_district_id = NULLIF($6, ''),
                   is_active       = TRUE,
                   updated_by      = $7,
                   updated_at      = NOW()
             WHERE cust_id = $1 AND ward_id = $2
               AND (ward IS DISTINCT FROM $3 OR province_id IS DISTINCT FROM NULLIF($4, '') OR regency_id IS DISTINCT FROM NULLIF($5, '') OR sub_district_id IS DISTINCT FROM NULLIF($6, '') OR is_active = FALSE)
        `, custId, wardId, name, provinceId, regencyId, subDistrictId, userId)
		return err
	}

	_, err := r.DB.Exec(`
        INSERT INTO mst.m_ward (cust_id, ward_id, ward, province_id, regency_id, sub_district_id, is_active, created_by, created_at, updated_by, updated_at)
        VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), NULLIF($6, ''), TRUE, $7, NOW(), $7, NOW())
    `, custId, wardId, name, provinceId, regencyId, subDistrictId, userId)
	return err
}

func (r *outletRepositoryImpl) CreateOutletTemp(historyId int64, status, custId string, data entity.OutletTemp) error {
	data.HistoryId = historyId
	data.CustId = custId
	data.StatusInsert = status
	query := `INSERT INTO import.outlet_temp (
		history_id, outlet_code, outlet_name, outlet_status,
		address1, zip_code, phone_no, fax_no, email, 
		disc_grp_name, ot_loc_name, ot_grp_name, price_grp_name,
		district_name, ot_class_name, industry_name, market_name,
		top, is_contra_bon, agent_from,
		credit_limit_type_name, credit_limit, credit_limit_action_name,
		sales_inv_limit_type_name, sales_inv_limit, sales_inv_limit_action_name,
		obs_type_name, obs, obs_limit_action_name,
		building_own, ar_status_name,
		tax_name, tax_identifier_type, tax_identifier_no, address_tax,
		ot_type_name, outlet_establishment_date,
		delv_addr1, delv_province, delv_regency, delv_sub_district,
		delv_ward, delv_longitude, delv_latitude, delv_zip_code, delv_is_same_addr,
		inv_addr1, inv_province, inv_regency, inv_sub_district,
		inv_ward, inv_zip_code, inv_is_same_addr,
		bank_name, account_no, account_name,
		contact_name, job_title, contact_phone_no, contact_wa_no,
		contact_email, contact_is_wa_no, identity_type, identity_no,
		longitude, latitude, nitku,
		tax_invoice_form_name, payment_type_name,
		status_insert, error_message, updated_at, deleted_at, cust_id
	) VALUES (
		:history_id, :outlet_code, :outlet_name, :outlet_status,
		:address1, :zip_code, :phone_no, :fax_no, :email, 
		:disc_grp_name, :ot_loc_name, :ot_grp_name, :price_grp_name,
		:district_name, :ot_class_name, :industry_name, :market_name,
		:top, :is_contra_bon, :agent_from,
		:credit_limit_type_name, :credit_limit, :credit_limit_action_name,
		:sales_inv_limit_type_name, :sales_inv_limit, :sales_inv_limit_action_name,
		:obs_type_name, :obs, :obs_limit_action_name,
		:building_own, :ar_status_name,
		:tax_name, :tax_identifier_type, :tax_identifier_no, :address_tax,
		:ot_type_name, :outlet_establishment_date,
		:delv_addr1, :delv_province, :delv_regency, :delv_sub_district,
		:delv_ward, :delv_longitude, :delv_latitude, :delv_zip_code, :delv_is_same_addr,
		:inv_addr1, :inv_province, :inv_regency, :inv_sub_district,
		:inv_ward, :inv_zip_code, :inv_is_same_addr,
		:bank_name, :account_no, :account_name,
		:contact_name, :job_title, :contact_phone_no, :contact_wa_no,
		:contact_email, :contact_is_wa_no, :identity_type, :identity_no,
		:longitude, :latitude, :nitku,
		:tax_invoice_form_name, :payment_type_name,
		:status_insert, :error_message, :updated_at, :deleted_at, :cust_id
    )`

	_, err := r.DB.NamedExec(query, data)
	return err
}

func (r *outletRepositoryImpl) GetOrCreateOutletType(custId, otTypeCode, otTypeName string) (int64, error) {
	var id int64
	query := `SELECT ot_type_id FROM mst.m_outlet_type WHERE cust_id = $1 AND ot_type_code = $2`
	err := r.DB.Get(&id, query, custId, otTypeCode)
	if err == nil {
		return id, nil
	}

	query = `INSERT INTO mst.m_outlet_type (cust_id, ot_type_code, ot_type_name, is_active, created_at) 
             VALUES ($1, $2, $3, true, NOW()) RETURNING ot_type_id`
	err = r.DB.QueryRow(query, custId, otTypeCode, otTypeName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// func (r *outletRepositoryImpl) GetOrCreateOutletUdf(custId string, outletId, udfId int64, udfValue string) (int64, error) {
//     var id int64
//     query := `SELECT udf_id FROM mst.m_outlet_udf WHERE cust_id = $1 AND outlet_id = $2 AND udf_id = $3`
//     err := r.DB.Get(&id, query, custId, outletId, udfId)
//     if err == nil {
//         return id, nil
//     }

//     query = `INSERT INTO mst.m_outlet_udf (cust_id, outlet_id, udf_id, udf_value)
//              VALUES ($1, $2, $3, $4) RETURNING udf_id`
//     err = r.DB.QueryRow(query, custId, outletId, udfId, udfValue).Scan(&id)
//     if err != nil {
//         return 0, err
//     }
//     return id, nil
// }

const createOutletQuery = `INSERT INTO mst.m_outlet (
                cust_id, outlet_code, outlet_name, outlet_status, address1, address2, city, zip_code, 
                phone_no, wa_no, fax_no, email, disc_grp_id, ot_loc_id, ot_grp_id, price_grp_id, 
                district_id, beat_id, sbeat_id, ot_class_id, industry_id, market_id, top, payment_type, 
                is_contra_bon, plu_grp_id, conv_grp_id, disc_inv_id, agent_from, credit_limit_type, credit_limit_type_name,
                credit_limit, sales_inv_limit_type, sales_inv_limit_type_name, sales_inv_limit, avg_sales_week, avg_sales_month, 
                first_trans_date, last_trans_date, first_week_no, ot_start_date, ot_reg_date, 
                building_own, dob, ar_status, ar_total, closed_date, is_emb_bail, is_pkp_outlet, tax_name,
                tax_addr1, tax_addr2, tax_city, tax_no, owner_name, owner_addr1, owner_addr2,
                owner_city, owner_phone_no, owner_id_no, delv_addr1, delv_addr2, delv_city, 
                inv_addr1, inv_addr2, inv_city, is_active, created_by, created_at, updated_by, updated_at, latitude, longitude,
                image_url, ot_type_id, is_obs, obs, outlet_province_id, outlet_regency_id, outlet_sub_district_id, outlet_ward_id, is_wa_no, delv_ward_id,
                delv_zip_code, delv_is_same_addr, inv_ward_id, inv_zip_code, inv_is_same_addr, 
                verification_status, verified_at, verified_by, tax_invoice_form, obs_type, 
                credit_limit_action, credit_limit_action_name, sales_inv_limit_action, sales_inv_limit_action_name, obs_limit_action, outlet_establishment_date, 
                delv_city2, delv_latitude, delv_longitude, delv_latitude2, delv_longitude2, 
                delv_ward_id2, delv_zip_code2, outlet_principal_code
             ) VALUES (
                :cust_id, :outlet_code, :outlet_name, :outlet_status, :address1, :address2, :city, :zip_code, 
                :phone_no, :wa_no, :fax_no, :email, :disc_grp_id, :ot_loc_id, :ot_grp_id, :price_grp_id, 
                :district_id, :beat_id, :sbeat_id, :ot_class_id, :industry_id, :market_id, :top, :payment_type, 
                :is_contra_bon, :plu_grp_id, :conv_grp_id, :disc_inv_id, :agent_from, :credit_limit_type, :credit_limit_type_name,
                :credit_limit, :sales_inv_limit_type, :sales_inv_limit_type_name, :sales_inv_limit, :avg_sales_week, :avg_sales_month, 
                :first_trans_date, :last_trans_date, :first_week_no, :ot_start_date, :ot_reg_date, 
                :building_own, :dob, :ar_status, :ar_total, :closed_date, :is_emb_bail, :is_pkp_outlet, :tax_name,
                :tax_addr1, :tax_addr2, :tax_city, :tax_no, :owner_name, :owner_addr1, :owner_addr2,
                :owner_city, :owner_phone_no, :owner_id_no, :delv_addr1, :delv_addr2, :delv_city, 
                :inv_addr1, :inv_addr2, :inv_city, :is_active, :created_by, :created_at, :updated_by, :updated_at, :latitude, :longitude,
                :image_url, :ot_type_id, :is_obs, :obs, :outlet_province_id, :outlet_regency_id, :outlet_sub_district_id, :outlet_ward_id, :is_wa_no, :delv_ward_id,
                :delv_zip_code, :delv_is_same_addr, :inv_ward_id, :inv_zip_code, :inv_is_same_addr, 
                :verification_status, :verified_at, :verified_by, :tax_invoice_form, :obs_type, 
                :credit_limit_action, :credit_limit_action_name, :sales_inv_limit_action, :sales_inv_limit_action_name, :obs_limit_action, :outlet_establishment_date, 
                :delv_city2, :delv_latitude, :delv_longitude, :delv_latitude2, :delv_longitude2, 
                :delv_ward_id2, :delv_zip_code2, :outlet_principal_code
             ) RETURNING outlet_id`

type outletNamedQueryer interface {
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

func createOutletWithNamedQuery(exec outletNamedQueryer, data entity.ProcessedOutlet) (int64, error) {
	rows, err := exec.NamedQuery(createOutletQuery, data)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		var outletID int64
		if err := rows.Scan(&outletID); err != nil {
			return 0, err
		}
		return outletID, nil
	}

	return 0, errors.New("failed to retrieve inserted outlet_id")
}

func (r *outletRepositoryImpl) CreateOutlet(data entity.ProcessedOutlet) (int64, error) {
	return createOutletWithNamedQuery(r.DB, data)
}

func (r *outletRepositoryImpl) CreateOutletTx(tx *sqlx.Tx, data entity.ProcessedOutlet) (int64, error) {
	return createOutletWithNamedQuery(tx, data)
}

func (r *outletRepositoryImpl) CreateOutletBank(bank model.MOutletBank) (int64, error) {
	query := `
        INSERT INTO mst.m_outlet_bank (
            cust_id, outlet_id, bank_id, account_no, account_name
        ) VALUES ($1, $2, $3, $4, $5)
        RETURNING outlet_bank_id`

	var id int64
	if err := r.DB.QueryRow(query, bank.CustID, bank.OutletID, bank.BankId, bank.AccountNo, bank.AccountName).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) CreateOutletContact(contact model.MOutletContact) (int64, error) {
	query := `
        INSERT INTO mst.m_outlet_contact (
            cust_id, outlet_id, contact_name, job_title, phone_no, wa_no,
            email, identity_no, is_wa_no, identity_type, fax_number
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING outlet_contact_id`

	var id int64
	if err := r.DB.QueryRow(query,
		contact.CustID, contact.OutletID, contact.ContactName, contact.JobTitle,
		contact.PhoneNo, contact.WaNo, contact.Email, contact.IdentityNo,
		contact.IsWaNo, contact.IdentityType, contact.FaxNumber,
	).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) CreateOutletTax(tax model.MOutletTax) (int64, error) {
	query := `
        INSERT INTO mst.m_outlet_tax (
            cust_id, outlet_id, tax_invoice_id, is_emb_bail, tax_no, tax_name,
            tax_addr1, tax_addr2, tax_city, tax_type, nitku, address_tax,
            tax_identifier_type, tax_identifier_no
        ) VALUES (
            $1, $2, $3, $4, $5, $6,
            $7, $8, $9, $10, $11, $12,
            $13, $14
        ) RETURNING outlet_tax_id`

	var id int64
	if err := r.DB.QueryRow(query,
		tax.CustID, tax.OutletID, tax.TaxInvoiceId, tax.IsEmbBail, tax.TaxNo, tax.TaxName,
		tax.TaxAddr1, tax.TaxAddr2, tax.TaxCity, tax.TaxType, tax.Nitku, tax.AdddressTax,
		tax.TaxIdentifierType, tax.TaxIdentifierNo,
	).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) GetProvinceIdByName(custId, name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", nil
	}

	var provinceId string
	err := r.DB.Get(&provinceId, `
		SELECT province_id
		  FROM mst.m_province
		 WHERE cust_id = $1
		   AND LOWER(province) = LOWER($2)
		   AND is_active = TRUE
		 LIMIT 1
	`, custId, name)

	if err != nil {
		// kalau tidak ketemu atau ada error, kembalikan kosong
		return "", nil
	}
	return provinceId, err
}

func (r *outletRepositoryImpl) GetRegencyIdByName(custId, name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", nil
	}

	var regencyId string
	err := r.DB.Get(&regencyId, `
		SELECT regency_id
		  FROM mst.m_regency
		 WHERE cust_id = $1
		   AND LOWER(regency) = LOWER($2)
		   AND is_active = TRUE
		 LIMIT 1
	`, custId, name)

	if err != nil {
		// kalau tidak ketemu atau ada error, kembalikan kosong
		return "", nil
	}
	return regencyId, err
}

func (r *outletRepositoryImpl) GetSubDistrictIdByName(custId, name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", nil
	}

	var subDistrictId string
	err := r.DB.Get(&subDistrictId, `
		SELECT sub_district_id
		  FROM mst.m_sub_district
		 WHERE cust_id = $1
		   AND LOWER(sub_district) = LOWER($2)
		   AND is_active = TRUE
		 LIMIT 1
	`, custId, name)

	if err != nil {
		// kalau tidak ketemu atau ada error, kembalikan kosong
		return "", nil
	}
	return subDistrictId, err
}

func (r *outletRepositoryImpl) GetWardIdByName(custId, name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", nil
	}

	var wardId string
	err := r.DB.Get(&wardId, `
		SELECT ward_id
		  FROM mst.m_ward
		 WHERE cust_id = $1
		   AND LOWER(ward) = LOWER($2)
		   AND is_active = TRUE
		 LIMIT 1
	`, custId, name)

	if err != nil {
		// kalau tidak ketemu atau ada error, kembalikan kosong
		return "", nil
	}
	return wardId, err
}

func (r *outletRepositoryImpl) CreateImportHistory(typeUpload, fileName, custId string, uploadedBy int64, totalData int) (int64, error) {
	var id int64
	query := `
		INSERT INTO import.import_history (file_name, uploaded_by, total_data, upload_type, cust_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING history_id
    `
	err := r.DB.QueryRow(query, fileName, uploadedBy, totalData, typeUpload, custId).Scan(&id)
	return id, err
}

func (r *outletRepositoryImpl) UpdateImportHistoryOutlet(historyId int64, success, failed int, statusReupload bool) error {
	query := `
        UPDATE import.import_history
        SET successful_data = $1,
            failed_data = $2,
            status_reupload = $3
        WHERE history_id = $4
    `
	_, err := r.DB.Exec(query, success, failed, statusReupload, historyId)
	return err
}

// Temp maintenance for reupload
func (r *outletRepositoryImpl) DeleteOutletUpdateTempByCtid(ctid string) error {
	_, err := r.DB.Exec(`DELETE FROM import.outlet_update_temp WHERE ctid = $1`, ctid)
	return err
}

func (r *outletRepositoryImpl) CountOutletUpdateTemp(historyId int64) (int, error) {
	var n int
	err := r.DB.Get(&n, `SELECT COUNT(1) FROM import.outlet_update_temp WHERE history_id = $1`, historyId)
	return n, err
}

func (r *outletRepositoryImpl) GetImportTotalData(historyId int64) (int, error) {
	var n int
	err := r.DB.Get(&n, `SELECT total_data FROM import.import_history WHERE history_id = $1`, historyId)
	return n, err
}

func (r *outletRepositoryImpl) InsertOutletUpdateTemp(temp entity.ImportOutletUpdateTemp) error {
	query := `
        INSERT INTO import.outlet_update_temp (
			history_id, cust_id,
			outlet_loc_name,
			outlet_type_name,
			outlet_grp_name,
			district_name,
			ot_class_name,
			disc_grp_name,
			market_name,
			industry_name,
			price_grp_name,
			bank_name,

			outlet_code, outlet_name, address1,
			outlet_province, outlet_regency, outlet_sub_district, outlet_ward,
			zip_code, longitude, latitude, building_own, outlet_establishment_date,

			phone_no, fax_no, barcode,

			contact_name, job_title, identity_type, identity_no,
			contact_phone_no, contact_is_wa_no, contact_wa_no, contact_email, outlet_contact_id,

			tax_invoice_form_name, tax_identifier_type, tax_identifier_no, nitku,
			tax_name, address_tax, outlet_tax_id,

			is_contra_bon, agent_from,

			delv_addr1, delv_province, delv_regency, delv_sub_district,
			delv_ward, delv_longitude, delv_latitude, delv_zip_code, delv_is_same_addr,

			inv_addr1, inv_province, inv_regency, inv_sub_district,
			inv_ward, inv_zip_code, inv_is_same_addr,

			payment_type_name, ar_status_name,

			account_no, account_name, outlet_bank_id,

			top, credit_limit_type_name, credit_limit, credit_limit_action_name,
			sales_inv_limit_type_name, sales_inv_limit, sales_inv_limit_action_name,
			obs_type_name, obs, obs_limit_action_name,

			outlet_id,
			status_insert, error_message, created_at
		) VALUES (
			:history_id, :cust_id,
			:outlet_loc_name,
			:outlet_type_name,
			:outlet_grp_name,
			:district_name,
			:ot_class_name,
			:disc_grp_name,
			:market_name,
			:industry_name,
			:price_grp_name,
			:bank_name,

			:outlet_code, :outlet_name, :address1,
			:outlet_province, :outlet_regency, :outlet_sub_district, :outlet_ward,
			:zip_code, :longitude, :latitude, :building_own, :outlet_establishment_date,

			:phone_no, :fax_no, :barcode,

			:contact_name, :job_title, :identity_type, :identity_no,
			:contact_phone_no, :contact_is_wa_no, :contact_wa_no, :contact_email, :outlet_contact_id,

			:tax_invoice_form_name, :tax_identifier_type, :tax_identifier_no, :nitku,
			:tax_name, :address_tax, :outlet_tax_id,

			:is_contra_bon, :agent_from,

			:delv_addr1, :delv_province, :delv_regency, :delv_sub_district,
			:delv_ward, :delv_longitude, :delv_latitude, :delv_zip_code, :delv_is_same_addr,

			:inv_addr1, :inv_province, :inv_regency, :inv_sub_district,
			:inv_ward, :inv_zip_code, :inv_is_same_addr,

			:payment_type_name, :ar_status_name,

			:account_no, :account_name, :outlet_bank_id,

			:top, :credit_limit_type_name, :credit_limit, :credit_limit_action_name,
			:sales_inv_limit_type_name, :sales_inv_limit, :sales_inv_limit_action_name,
			:obs_type_name, :obs, :obs_limit_action_name,

			:outlet_id,
			:status_insert, :error_message, NOW()
		)
    `
	_, err := r.DB.NamedExec(query, temp)
	return err
}

// Reupload insert helpers for import.outlet_temp
func (r *outletRepositoryImpl) GetCtidsOutletInsertByHistory(historyId int64) ([]string, error) {
	query := `
        SELECT ctid::text
        FROM import.outlet_temp
        WHERE history_id = $1
        ORDER BY ctid
    `
	rows, err := r.DB.Query(query, historyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ctids []string
	for rows.Next() {
		var ctid string
		if err := rows.Scan(&ctid); err != nil {
			return nil, err
		}
		ctids = append(ctids, ctid)
	}
	return ctids, nil
}

func (r *outletRepositoryImpl) DeleteOutletTempByCtid(ctid string) error {
	_, err := r.DB.Exec(`DELETE FROM import.outlet_temp WHERE ctid = $1`, ctid)
	return err
}

func (r *outletRepositoryImpl) CountOutletTemp(historyId int64) (int, error) {
	var n int
	err := r.DB.Get(&n, `SELECT COUNT(1) FROM import.outlet_temp WHERE history_id = $1`, historyId)
	return n, err
}

// ========================================================
// 4. Data untuk Template Update
// ========================================================
func (r *outletRepositoryImpl) GetOutletDataForTemplateUpdate(custId string, fields []string) (map[string][][]string, error) {
	result := make(map[string][][]string)

	// Whitelist kolom outlet + relasi (termasuk kode/nama master yang sering diminta)
	outletSelectOrder := []string{
		"o.cust_id", "o.outlet_id", "o.outlet_code", "o.barcode", "o.outlet_name", "o.outlet_status",
		"o.address1", "o.address2", "o.city", "o.zip_code", "o.phone_no", "o.wa_no", "o.fax_no", "o.email",

		// Disc group
		"o.disc_grp_id", "dg.disc_grp_code", "dg.disc_grp_name",
		// Other masters (id + code + name)
		"o.ot_loc_id", "ol.ot_loc_code", "ol.ot_loc_name",
		"o.ot_grp_id", "og.ot_grp_code", "og.ot_grp_name",
		"o.price_grp_id", "pg.price_grp_code", "pg.price_grp_name",
		"o.district_id", "d.district_code", "d.district_name",
		"o.ot_class_id", "cls.ot_class_code", "cls.ot_class_name",
		"o.industry_id", "ind.industry_code", "ind.industry_name",
		"o.market_id", "m.market_code", "m.market_name",
		"o.top", "o.payment_type", "o.is_contra_bon",

		"o.plu_grp_id", "o.conv_grp_id", "o.disc_inv_id", "o.agent_from",
		"o.credit_limit_type", "o.credit_limit", "o.sales_inv_limit_type", "o.sales_inv_limit",

		"o.avg_sales_week", "o.avg_sales_month", "o.first_trans_date", "o.last_trans_date", "o.first_week_no",
		"o.ot_start_date", "o.ot_reg_date", "o.building_own", "o.dob", "o.ar_status", "o.ar_total", "o.closed_date", "o.is_emb_bail",

		// Tax & Owner come from outlet master
		"o.tax_name", "o.tax_addr1", "o.tax_addr2", "o.tax_city", "o.tax_no",
		"o.owner_name", "o.owner_addr1", "o.owner_addr2", "o.owner_city", "o.owner_phone_no", "o.owner_id_no",

		"o.delv_addr1", "o.delv_addr2", "o.delv_city",
		"o.inv_addr1", "o.inv_addr2", "o.inv_city",

		"o.is_active", "o.created_by", "o.created_at", "o.updated_by", "o.updated_at", "o.is_del", "o.deleted_by", "o.deleted_at",
		"o.latitude", "o.longitude", "o.image_url",

		"o.ot_type_id", "o.is_obs", "o.obs",
		"o.outlet_ward_id", "w.ward AS outlet_ward", "o.is_wa_no",

		"COALESCE(NULLIF(TRIM(o.outlet_sub_district_id), ''), sd.sub_district_id) AS outlet_sub_district_id",
		"COALESCE(osdm.sub_district, sd.sub_district) AS outlet_sub_district",
		"COALESCE(NULLIF(TRIM(o.outlet_regency_id), ''), r.regency_id) AS outlet_regency_id",
		"COALESCE(oreg.regency, r.regency) AS outlet_regency",
		"COALESCE(NULLIF(TRIM(o.outlet_province_id), ''), p.province_id) AS outlet_province_id",
		"COALESCE(oprov.province, p.province) AS outlet_province",

		// Delivery region
		"o.delv_ward_id", "w2.ward AS delv_ward",
		"sd2.sub_district_id AS delv_sub_district_id", "sd2.sub_district AS delv_sub_district",
		"r2.regency_id AS delv_regency_id", "r2.regency AS delv_regency",
		"p2.province_id AS delv_province_id", "p2.province AS delv_province",
		"o.delv_zip_code", "o.delv_is_same_addr",

		// Invoice region
		"o.inv_ward_id", "w3.ward AS inv_ward",
		"sd3.sub_district_id AS inv_sub_district_id", "sd3.sub_district AS inv_sub_district",
		"r3.regency_id AS inv_regency_id", "r3.regency AS inv_regency",
		"p3.province_id AS inv_province_id", "p3.province AS inv_province",
		"o.inv_zip_code", "o.inv_is_same_addr",

		"o.verification_status", "o.verified_at", "o.verified_by", "o.tax_invoice_form", "o.obs_type",
		"o.credit_limit_action", "o.sales_inv_limit_action", "o.obs_limit_action",
		"o.outlet_establishment_date",

		"o.delv_city2", "o.delv_latitude", "o.delv_longitude", "o.delv_latitude2", "o.delv_longitude2", "o.delv_ward_id2", "o.delv_zip_code2",
	}

	// lookup validasi kolom
	// Map kolom ke kunci header tanpa prefix alias tabel dan mendukung alias SQL (AS ...)
	aliasKey := func(col string) string {
		s := strings.TrimSpace(col)
		// handle AS alias (case-insensitive)
		lower := strings.ToLower(s)
		if idx := strings.Index(lower, " as "); idx != -1 {
			return strings.TrimSpace(s[idx+4:])
		}
		// otherwise take last segment after '.'
		parts := strings.Split(s, ".")
		return strings.TrimSpace(parts[len(parts)-1])
	}
	outletCols := map[string]struct{}{}
	for _, c := range outletSelectOrder {
		outletCols[aliasKey(c)] = struct{}{}
	}

	includeAllOutlet := false
	requestedOutletCols := make([]string, 0, len(outletSelectOrder))
	seenOutlet := map[string]struct{}{}

	// loop field
	for _, f := range fields {
		token := strings.ToLower(strings.TrimSpace(f))
		if token == "outlet" {
			includeAllOutlet = true
		}
		if _, ok := outletCols[token]; ok {
			if _, seen := seenOutlet[token]; !seen {
				seenOutlet[token] = struct{}{}
				requestedOutletCols = append(requestedOutletCols, token)
			}
			continue
		}

		// Switch case untuk master/reference tables
		switch token {
		case "bank":
			rows, err := r.fetchData(
				`SELECT bank_id, bank_code, bank_name 
				 FROM mst.m_bank 
				 WHERE cust_id = $1 ORDER BY bank_id`, custId)
			if err != nil {
				return nil, err
			}
			result["bank_id"] = getColumn(rows, 0)
			result["bank_code"] = getColumn(rows, 1)
			result["bank_name"] = getColumn(rows, 2)

		case "outlet_bank":
			rows, err := r.fetchData(
				`SELECT ob.outlet_bank_id, ob.bank_id, b.bank_code, b.bank_name, ob.account_no, ob.account_name
				 FROM mst.m_outlet_bank ob
				 LEFT JOIN mst.m_bank b ON b.cust_id = ob.cust_id AND b.bank_id = ob.bank_id
				 WHERE ob.cust_id = $1
				 ORDER BY ob.outlet_id, ob.outlet_bank_id`, custId)
			if err != nil {
				return nil, err
			}
			result["outlet_bank_id"] = getColumn(rows, 0)
			result["bank_id"] = getColumn(rows, 1)
			result["bank_code"] = getColumn(rows, 2)
			result["bank_name"] = getColumn(rows, 3)
			result["account_no"] = getColumn(rows, 4)
			result["account_name"] = getColumn(rows, 5)

		case "outlet_contact":
			rows, err := r.fetchData(
				`SELECT outlet_contact_id, contact_name, job_title, phone_no AS contact_phone_no, wa_no AS contact_wa_no,
						email AS contact_email, identity_no, is_wa_no AS contact_is_wa_no, identity_type, fax_number
				 FROM mst.m_outlet_contact
				 WHERE cust_id = $1
				 ORDER BY outlet_id, outlet_contact_id`, custId)
			if err != nil {
				return nil, err
			}
			result["outlet_contact_id"] = getColumn(rows, 0)
			result["contact_name"] = getColumn(rows, 1)
			result["job_title"] = getColumn(rows, 2)
			result["contact_phone_no"] = getColumn(rows, 3)
			result["contact_wa_no"] = getColumn(rows, 4)
			result["contact_email"] = getColumn(rows, 5)
			result["identity_no"] = getColumn(rows, 6)
			result["contact_is_wa_no"] = getColumn(rows, 7)
			result["identity_type"] = getColumn(rows, 8)
			result["fax_number"] = getColumn(rows, 9)

		case "outlet_tax":
			rows, err := r.fetchData(
				`SELECT outlet_tax_id, tax_invoice_id, tax_type, nitku, address_tax, tax_identifier_type, tax_identifier_no
				 FROM mst.m_outlet_tax
				 WHERE cust_id = $1
				 ORDER BY outlet_id, outlet_tax_id`, custId)
			if err != nil {
				return nil, err
			}
			result["outlet_tax_id"] = getColumn(rows, 0)
			result["tax_invoice_id"] = getColumn(rows, 1)
			result["tax_type"] = getColumn(rows, 2)
			result["nitku"] = getColumn(rows, 3)
			result["address_tax"] = getColumn(rows, 4)
			result["tax_identifier_type"] = getColumn(rows, 5)
			result["tax_identifier_no"] = getColumn(rows, 6)

		case "discount":
			rows, err := r.fetchData(
				`SELECT disc_grp_id, disc_grp_code, disc_grp_name 
				 FROM mst.m_disc_group 
				 WHERE cust_id = $1 ORDER BY disc_grp_id`, custId)
			if err != nil {
				return nil, err
			}
			result["disc_grp_id"] = getColumn(rows, 0)
			result["disc_grp_code"] = getColumn(rows, 1)
			result["disc_grp_name"] = getColumn(rows, 2)

		case "price":
			rows, err := r.fetchData(
				`SELECT price_grp_id, price_grp_code, price_grp_name 
				 FROM mst.m_price_group 
				 WHERE cust_id = $1 ORDER BY price_grp_id`, custId)
			if err != nil {
				return nil, err
			}
			result["price_grp_id"] = getColumn(rows, 0)
			result["price_grp_code"] = getColumn(rows, 1)
			result["price_grp_name"] = getColumn(rows, 2)

		case "class":
			rows, err := r.fetchData(
				`SELECT ot_class_id, ot_class_code, ot_class_name 
				 FROM mst.m_outlet_class 
				 WHERE cust_id = $1 ORDER BY ot_class_id`, custId)
			if err != nil {
				return nil, err
			}
			result["ot_class_id"] = getColumn(rows, 0)
			result["ot_class_code"] = getColumn(rows, 1)
			result["ot_class_name"] = getColumn(rows, 2)

		case "group":
			rows, err := r.fetchData(
				`SELECT ot_grp_id, ot_grp_code, ot_grp_name 
				 FROM mst.m_outlet_group 
				 WHERE cust_id = $1 ORDER BY ot_grp_id`, custId)
			if err != nil {
				return nil, err
			}
			result["ot_grp_id"] = getColumn(rows, 0)
			result["ot_grp_code"] = getColumn(rows, 1)
			result["ot_grp_name"] = getColumn(rows, 2)

		case "type":
			rows, err := r.fetchData(
				`SELECT ot_type_id, ot_type_code, ot_type_name 
				 FROM mst.m_outlet_type 
				 WHERE cust_id = $1 ORDER BY ot_type_id`, custId)
			if err != nil {
				return nil, err
			}
			result["ot_type_id"] = getColumn(rows, 0)
			result["ot_type_code"] = getColumn(rows, 1)
			result["ot_type_name"] = getColumn(rows, 2)

		case "district":
			rows, err := r.fetchData(
				`SELECT district_id, district_code, district_name 
				 FROM mst.m_district 
				 WHERE cust_id = $1 ORDER BY district_id`, custId)
			if err != nil {
				return nil, err
			}
			result["district_id"] = getColumn(rows, 0)
			result["district_code"] = getColumn(rows, 1)
			result["district_name"] = getColumn(rows, 2)

		case "market":
			rows, err := r.fetchData(
				`SELECT market_id, market_code, market_name 
				 FROM mst.m_market 
				 WHERE cust_id = $1 ORDER BY market_id`, custId)
			if err != nil {
				return nil, err
			}
			result["market_id"] = getColumn(rows, 0)
			result["market_code"] = getColumn(rows, 1)
			result["market_name"] = getColumn(rows, 2)

		case "industry":
			rows, err := r.fetchData(
				`SELECT industry_id, industry_code, industry_name 
				 FROM mst.m_industry 
				 WHERE cust_id = $1 ORDER BY industry_id`, custId)
			if err != nil {
				return nil, err
			}
			result["industry_id"] = getColumn(rows, 0)
			result["industry_code"] = getColumn(rows, 1)
			result["industry_name"] = getColumn(rows, 2)

		case "province":
			rows, err := r.fetchData(
				`SELECT province_id, province
				 FROM mst.m_province
				 WHERE cust_id = $1 ORDER BY province_id`, custId)
			if err != nil {
				return nil, err
			}
			result["province_id"] = getColumn(rows, 0)
			result["province"] = getColumn(rows, 1)

		case "regency":
			rows, err := r.fetchData(
				`SELECT regency_id, regency, province_id
				 FROM mst.m_regency
				 WHERE cust_id = $1 ORDER BY regency_id`, custId)
			if err != nil {
				return nil, err
			}
			result["regency_id"] = getColumn(rows, 0)
			result["regency"] = getColumn(rows, 1)
			result["province_id"] = getColumn(rows, 2)

		case "sub_district":
			rows, err := r.fetchData(
				`SELECT sub_district_id, sub_district, province_id, regency_id
				 FROM mst.m_sub_district
				 WHERE cust_id = $1 ORDER BY sub_district_id`, custId)
			if err != nil {
				return nil, err
			}
			result["sub_district_id"] = getColumn(rows, 0)
			result["sub_district"] = getColumn(rows, 1)
			result["province_id"] = getColumn(rows, 2)
			result["regency_id"] = getColumn(rows, 3)

		case "ward":
			rows, err := r.fetchData(
				`SELECT ward_id, ward, province_id, regency_id, sub_district_id
				 FROM mst.m_ward
				 WHERE cust_id = $1 ORDER BY ward_id`, custId)
			if err != nil {
				return nil, err
			}
			result["ward_id"] = getColumn(rows, 0)
			result["ward"] = getColumn(rows, 1)
			result["province_id"] = getColumn(rows, 2)
			result["regency_id"] = getColumn(rows, 3)
			result["sub_district_id"] = getColumn(rows, 4)

		default:
			// field tidak dikenal → skip
		}
	}

	// Ambil data outlet (join semua reference langsung)
	if includeAllOutlet || len(requestedOutletCols) > 0 {
		cols := outletSelectOrder
		if !includeAllOutlet && len(requestedOutletCols) > 0 {
			rq := map[string]struct{}{}
			for _, k := range requestedOutletCols {
				rq[k] = struct{}{}
			}
			filtered := make([]string, 0, len(requestedOutletCols))
			for _, k := range outletSelectOrder {
				if _, ok := rq[aliasKey(k)]; ok {
					filtered = append(filtered, k)
				}
			}
			cols = filtered
		}

		if len(cols) > 0 {
			q := `SELECT ` + strings.Join(cols, ",") + `
		    FROM mst.m_outlet o
		    LEFT JOIN mst.m_disc_group dg 
			    ON dg.disc_grp_id = o.disc_grp_id AND dg.cust_id = o.cust_id
		    LEFT JOIN mst.m_outlet_loc ol 
			    ON ol.ot_loc_id = o.ot_loc_id AND ol.cust_id = o.cust_id
		    LEFT JOIN mst.m_price_group pg 
			    ON pg.price_grp_id = o.price_grp_id AND pg.cust_id = o.cust_id
		    LEFT JOIN mst.m_outlet_class cls 
			    ON cls.ot_class_id = o.ot_class_id AND cls.cust_id = o.cust_id
		    LEFT JOIN mst.m_outlet_group og 
			    ON og.ot_grp_id = o.ot_grp_id AND og.cust_id = o.cust_id
		    LEFT JOIN mst.m_outlet_type ot 
			    ON ot.ot_type_id = o.ot_type_id AND ot.cust_id = o.cust_id
		    LEFT JOIN mst.m_district d 
			    ON d.district_id = o.district_id AND d.cust_id = o.cust_id
		    LEFT JOIN mst.m_market m 
			    ON m.market_id = o.market_id AND m.cust_id = o.cust_id
		    LEFT JOIN mst.m_industry ind 
			    ON ind.industry_id = o.industry_id AND ind.cust_id = o.cust_id
			LEFT JOIN mst.m_ward w 
			    ON w.ward_id = o.outlet_ward_id AND w.cust_id = o.cust_id
			LEFT JOIN mst.m_sub_district sd 
			    ON sd.sub_district_id = w.sub_district_id AND sd.cust_id = w.cust_id
			LEFT JOIN mst.m_regency r 
			    ON r.regency_id = sd.regency_id AND r.cust_id = sd.cust_id
			LEFT JOIN mst.m_province p 
			    ON p.province_id = r.province_id AND p.cust_id = r.cust_id
			LEFT JOIN mst.m_province oprov
			    ON oprov.province_id = o.outlet_province_id AND oprov.cust_id = o.cust_id
			LEFT JOIN mst.m_regency oreg
			    ON oreg.regency_id = o.outlet_regency_id AND oreg.cust_id = o.cust_id
			LEFT JOIN mst.m_sub_district osdm
			    ON osdm.sub_district_id = o.outlet_sub_district_id AND osdm.cust_id = o.cust_id

			-- Delivery hierarchy
			LEFT JOIN mst.m_ward w2 
			    ON w2.ward_id = o.delv_ward_id AND w2.cust_id = o.cust_id
			LEFT JOIN mst.m_sub_district sd2 
			    ON sd2.sub_district_id = w2.sub_district_id AND sd2.cust_id = w2.cust_id
			LEFT JOIN mst.m_regency r2 
			    ON r2.regency_id = sd2.regency_id AND r2.cust_id = sd2.cust_id
			LEFT JOIN mst.m_province p2 
			    ON p2.province_id = r2.province_id AND p2.cust_id = r2.cust_id

			-- Invoice hierarchy
			LEFT JOIN mst.m_ward w3 
			    ON w3.ward_id = o.inv_ward_id AND w3.cust_id = o.cust_id
			LEFT JOIN mst.m_sub_district sd3 
			    ON sd3.sub_district_id = w3.sub_district_id AND sd3.cust_id = w3.cust_id
			LEFT JOIN mst.m_regency r3 
			    ON r3.regency_id = sd3.regency_id AND r3.cust_id = sd3.cust_id
			LEFT JOIN mst.m_province p3 
			    ON p3.province_id = r3.province_id AND p3.cust_id = r3.cust_id
		    WHERE o.cust_id = $1
		    ORDER BY o.outlet_id`

			rows, err := r.fetchData(q, custId)
			if err != nil {
				return nil, err
			}
			for i, col := range cols {
				key := aliasKey(col)
				result[key] = getColumn(rows, i)
			}
		}
	}

	return result, nil
}

// ========================================================
// helper fetchData & getColumn
// ========================================================
func (r *outletRepositoryImpl) fetchData(query string, args ...interface{}) ([][]string, error) {
	rows, err := r.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := [][]string{}
	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range cols {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := []string{}
		for _, val := range values {
			if val != nil {
				row = append(row, fmt.Sprintf("%v", val))
			} else {
				row = append(row, "")
			}
		}
		results = append(results, row)
	}
	return results, nil
}

func getColumn(rows [][]string, idx int) [][]string {
	col := [][]string{}
	for _, row := range rows {
		if idx < len(row) {
			col = append(col, []string{row[idx]})
		} else {
			col = append(col, []string{""})
		}
	}
	return col
}

// 4 LOCATION: existence checks
func (r *outletRepositoryImpl) CheckProvinceExists(custId, provinceId string) (bool, error) {
	var exists bool
	err := r.DB.Get(&exists, `SELECT COUNT(1) > 0 FROM mst.m_province WHERE cust_id = $1 AND province_id = $2`, custId, provinceId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckRegencyExists(custId, regencyId string) (bool, error) {
	var exists bool
	err := r.DB.Get(&exists, `SELECT COUNT(1) > 0 FROM mst.m_regency WHERE cust_id = $1 AND regency_id = $2`, custId, regencyId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckSubDistrictExists(custId, subDistrictId string) (bool, error) {
	var exists bool
	err := r.DB.Get(&exists, `SELECT COUNT(1) > 0 FROM mst.m_sub_district WHERE cust_id = $1 AND sub_district_id = $2`, custId, subDistrictId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckWardExists(custId, wardId string) (bool, error) {
	var exists bool
	err := r.DB.Get(&exists, `SELECT COUNT(1) > 0 FROM mst.m_ward WHERE cust_id = $1 AND ward_id = $2`, custId, wardId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckOutletExists(custId, outletId string) (bool, error) {
	var exists bool
	err := r.DB.Get(&exists, `SELECT COUNT(1) > 0 FROM mst.m_outlet WHERE cust_id = $1 AND outlet_id = $2`, custId, outletId)
	return exists, err
}

// 4 LOCATION: update methods
func (r *outletRepositoryImpl) UpdateImportProvince(custId, provinceId, province string) error {
	query := `UPDATE mst.m_province
              SET province = NULLIF($3, ''), updated_at = NOW()
              WHERE cust_id = $1 AND province_id = $2`
	res, err := r.DB.Exec(query, custId, provinceId, province)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no province updated (cust_id=%s, province_id=%s)", custId, provinceId)
	}
	return nil
}

func (r *outletRepositoryImpl) UpdateImportRegency(custId, regencyId, regency, provinceId string) error {
	query := `UPDATE mst.m_regency
              SET regency = NULLIF($3, ''), province_id = NULLIF($4, ''), updated_at = NOW()
              WHERE cust_id = $1 AND regency_id = $2`
	res, err := r.DB.Exec(query, custId, regencyId, regency, provinceId)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no regency updated (cust_id=%s, regency_id=%s)", custId, regencyId)
	}
	return nil
}

func (r *outletRepositoryImpl) UpdateImportSubDistrict(custId, subDistrictId, subDistrict, provinceId, regencyId string) error {
	query := `UPDATE mst.m_sub_district
              SET sub_district = NULLIF($3, ''), province_id = NULLIF($4, ''), regency_id = NULLIF($5, ''), updated_at = NOW()
              WHERE cust_id = $1 AND sub_district_id = $2`
	res, err := r.DB.Exec(query, custId, subDistrictId, subDistrict, provinceId, regencyId)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no sub_district updated (cust_id=%s, sub_district_id=%s)", custId, subDistrictId)
	}
	return nil
}

func (r *outletRepositoryImpl) UpdateImportWard(custId, wardId, ward, provinceId, regencyId, subDistrictId string) error {
	query := `UPDATE mst.m_ward
              SET ward = NULLIF($3, ''), province_id = NULLIF($4, ''), regency_id = NULLIF($5, ''), sub_district_id = NULLIF($6, ''), updated_at = NOW()
              WHERE cust_id = $1 AND ward_id = $2`
	res, err := r.DB.Exec(query, custId, wardId, ward, provinceId, regencyId, subDistrictId)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no ward updated (cust_id=%s, ward_id=%s)", custId, wardId)
	}
	return nil
}

func (r *outletRepositoryImpl) UpdateOutletImport(data entity.ProcessedUpdateOutlet) error {
	// Query diperbarui dengan semua kolom yang relevan dari definisi tabel
	query := `
		UPDATE mst.m_outlet
		SET
		outlet_code    = COALESCE(NULLIF($1, ''), outlet_code),
		barcode        = COALESCE(NULLIF($2, ''), barcode),
		outlet_name    = COALESCE(NULLIF($3, ''), outlet_name),
		outlet_status  = COALESCE(NULLIF($4, 0), outlet_status),
		address1       = COALESCE(NULLIF($5, ''), address1),
		address2       = COALESCE(NULLIF($6, ''), address2),
		city           = COALESCE(NULLIF($7, ''), city),
		zip_code       = COALESCE(NULLIF($8, ''), zip_code),
		phone_no       = COALESCE(NULLIF($9, ''), phone_no),
		wa_no          = COALESCE(NULLIF($10, ''), wa_no),
		fax_no         = COALESCE(NULLIF($11, ''), fax_no),
		email          = COALESCE(NULLIF($12, ''), email),
		disc_grp_id    = COALESCE(NULLIF($13, 0), disc_grp_id),
		ot_loc_id      = COALESCE(NULLIF($14, 0), ot_loc_id),
		ot_grp_id      = COALESCE(NULLIF($15, 0), ot_grp_id),
		price_grp_id   = COALESCE(NULLIF($16, 0), price_grp_id),
		district_id    = COALESCE(NULLIF($17, 0), district_id),
		beat_id        = COALESCE(NULLIF($18, 0), beat_id),
		sbeat_id       = COALESCE(NULLIF($19, 0), sbeat_id),
		ot_class_id    = COALESCE(NULLIF($20, 0), ot_class_id),
		industry_id    = COALESCE(NULLIF($21, 0), industry_id),
		market_id      = COALESCE(NULLIF($22, 0), market_id),
		top            = COALESCE(NULLIF($23, 0), top),
		payment_type   = COALESCE(NULLIF($24, 0), payment_type),
		is_contra_bon  = COALESCE(NULLIF($25, true), is_contra_bon),
		plu_grp_id     = COALESCE(NULLIF($26, 0), plu_grp_id),
		conv_grp_id    = COALESCE(NULLIF($27, 0), conv_grp_id),
		disc_inv_id    = COALESCE(NULLIF($28, 0), disc_inv_id),
		agent_from     = COALESCE(NULLIF($29, ''), agent_from),
		credit_limit_type   = COALESCE(NULLIF($30, 0), credit_limit_type),
		credit_limit        = COALESCE(NULLIF($31, 0), credit_limit),
		sales_inv_limit_type= COALESCE(NULLIF($32, 0), sales_inv_limit_type),
		sales_inv_limit     = COALESCE(NULLIF($33, 0), sales_inv_limit),
		avg_sales_week      = COALESCE(NULLIF($34, 0), avg_sales_week),
		avg_sales_month     = COALESCE(NULLIF($35, 0), avg_sales_month),
		first_week_no       = COALESCE(NULLIF($38, 0), first_week_no),
		building_own        = COALESCE(NULLIF($41, 0), building_own),
		ar_status = COALESCE(NULLIF($43, 0), ar_status), 
		ar_total = COALESCE(NULLIF($44, 0), ar_total),
		is_emb_bail         = COALESCE(NULLIF($46, false), is_emb_bail),
		tax_name            = COALESCE(NULLIF($47, ''), tax_name),
		tax_addr1           = COALESCE(NULLIF($48, ''), tax_addr1),
		tax_addr2           = COALESCE(NULLIF($49, ''), tax_addr2),
		tax_city            = COALESCE(NULLIF($50, ''), tax_city),
		tax_no              = COALESCE(NULLIF($51, ''), tax_no),
		owner_name          = COALESCE(NULLIF($52, ''), owner_name),
		owner_addr1         = COALESCE(NULLIF($53, ''), owner_addr1),
		owner_addr2         = COALESCE(NULLIF($54, ''), owner_addr2),
		owner_city          = COALESCE(NULLIF($55, ''), owner_city),
		owner_phone_no      = COALESCE(NULLIF($56, ''), owner_phone_no),
		owner_id_no         = COALESCE(NULLIF($57, ''), owner_id_no),
		delv_addr1          = COALESCE(NULLIF($58, ''), delv_addr1),
		delv_addr2          = COALESCE(NULLIF($59, ''), delv_addr2),
		delv_city           = COALESCE(NULLIF($60, ''), delv_city),
		inv_addr1           = COALESCE(NULLIF($61, ''), inv_addr1),
		inv_addr2           = COALESCE(NULLIF($62, ''), inv_addr2),
		inv_city            = COALESCE(NULLIF($63, ''), inv_city),
		is_active           = COALESCE($64, is_active),
		latitude            = COALESCE(NULLIF($65, ''), latitude),
		longitude           = COALESCE(NULLIF($66, ''), longitude),
		image_url           = COALESCE(NULLIF($67, ''), image_url),
		ot_type_id          = COALESCE(NULLIF($68, 0), ot_type_id),
		is_obs              = COALESCE(NULLIF($69, false), is_obs),
		obs                 = COALESCE(NULLIF($70, 0), obs),
		outlet_ward_id      = COALESCE(NULLIF($71, ''), outlet_ward_id),
		is_wa_no            = COALESCE(NULLIF($72, false), is_wa_no),
		delv_ward_id        = COALESCE(NULLIF($73, ''), delv_ward_id),
		delv_zip_code       = COALESCE(NULLIF($74, ''), delv_zip_code),
		delv_is_same_addr   = COALESCE(NULLIF($75, false), delv_is_same_addr),
		inv_ward_id         = COALESCE(NULLIF($76, ''), inv_ward_id),
		inv_zip_code        = COALESCE(NULLIF($77, ''), inv_zip_code),
		inv_is_same_addr    = COALESCE(NULLIF($78, false), inv_is_same_addr),
		verification_status = COALESCE(NULLIF($79, 0), verification_status),
		verified_by         = COALESCE(NULLIF($81, 0), verified_by),
		tax_invoice_form    = COALESCE(NULLIF($82, 0), tax_invoice_form),
		obs_type            = COALESCE(NULLIF($83, 0), obs_type),
		credit_limit_action = COALESCE(NULLIF($84, 0), credit_limit_action),
		sales_inv_limit_action = COALESCE(NULLIF($85, 0), sales_inv_limit_action),
		obs_limit_action    = COALESCE(NULLIF($86, 0), obs_limit_action),
		first_trans_date = COALESCE($36, first_trans_date),
		last_trans_date  = COALESCE($37, last_trans_date),
		ot_start_date  = COALESCE($39, ot_start_date),
		ot_reg_date  = COALESCE($40, ot_reg_date),
		dob              = COALESCE($42, dob),
		closed_date      = COALESCE($45, closed_date),
		outlet_establishment_date = COALESCE($87, outlet_establishment_date),
		verified_at      = COALESCE($80, verified_at),
		delv_city2          = COALESCE(NULLIF($88, ''), delv_city2),
		delv_latitude       = COALESCE(NULLIF($89, ''), delv_latitude),
		delv_longitude      = COALESCE(NULLIF($90, ''), delv_longitude),
		delv_latitude2      = COALESCE(NULLIF($91, ''), delv_latitude2),
		delv_longitude2     = COALESCE(NULLIF($92, ''), delv_longitude2),
		delv_ward_id2       = COALESCE(NULLIF($93, ''), delv_ward_id2),
		delv_zip_code2      = COALESCE(NULLIF($94, ''), delv_zip_code2),
		is_pkp_outlet       = COALESCE($95, is_pkp_outlet),
			updated_by 				 = $96,
			updated_at 				 = now()
		WHERE cust_id = $97 AND outlet_id = $98;
	`

	nullTime := func(t time.Time) interface{} {
		if t.IsZero() {
			return nil
		}
		return t
	}

	firstTransDate := nullTime(data.FirstTransDate)
	lastTransDate := nullTime(data.LastTransDate)
	otStartDate := nullTime(data.OtStartDate)
	otRegDate := nullTime(data.OtRegDate)
	dob := nullTime(data.Dob)
	closedDate := nullTime(data.ClosedDate)
	outletEstablishmentDate := nullTime(data.OutletEstablishmentDate)
	verifiedAt := nullTime(data.VerifiedAt)
	outletWardIDStr := ""
	if data.OutletWardId != "0" {
		// outletWardIDStr = strconv.FormatInt(data.OutletWardId, 10)
	}
	isActiveParam := sql.NullBool{}
	if data.IsActive != nil {
		isActiveParam.Bool = *data.IsActive
		isActiveParam.Valid = true
	}
	isPkpOutletParam := sql.NullBool{}
	if data.IsPkpOutlet != nil {
		isPkpOutletParam.Bool = *data.IsPkpOutlet
		isPkpOutletParam.Valid = true
	}

	// Daftar parameter diperpanjang sesuai dengan jumlah kolom di query
	_, err := r.DB.Exec(query,
		data.OutletCode,
		data.Barcode,
		data.OutletName,
		data.OutletStatus,
		data.Address1,
		data.Address2,
		data.City,
		data.ZipCode,
		data.PhoneNo,
		data.WaNo,
		data.FaxNo,
		data.Email,
		data.DiscGrpId,
		data.OtLocId,
		data.OtGrpId,
		data.PriceGrpId,
		data.DistrictId,
		data.BeatId,
		data.SbeatId,
		data.OtClassId,
		data.IndustryId,
		data.MarketId,
		data.Top,
		data.PaymentType,
		data.IsContraBon,
		data.PluGrpId,
		data.ConvGrpId,
		data.DiscInvId,
		data.AgentFrom,
		data.CreditLimitType,
		data.CreditLimit,
		data.SalesInvLimitType,
		data.SalesInvLimit,
		data.AvgSalesWeek,
		data.AvgSalesMonth,
		firstTransDate,
		lastTransDate,
		data.FirstWeekNo,
		otStartDate,
		otRegDate,
		data.BuildingOwn,
		dob,
		data.ArStatus,
		data.ArTotal,
		closedDate,
		data.IsEmbBail,
		data.TaxName,
		data.TaxAddr1,
		data.TaxAddr2,
		data.TaxCity,
		data.TaxNo,
		data.OwnerName,
		data.OwnerAddr1,
		data.OwnerAddr2,
		data.OwnerCity,
		data.OwnerPhoneNo,
		data.OwnerIdNo,
		data.DelvAddr1,
		data.DelvAddr2,
		data.DelvCity,
		data.InvAddr1,
		data.InvAddr2,
		data.InvCity,
		isActiveParam,
		// -- Parameter tambahan dimulai di sini --
		data.Latitude,
		data.Longitude,
		data.ImageUrl,
		data.OtTypeId,
		data.IsObs,
		data.Obs,
		outletWardIDStr,
		data.IsWaNo,
		data.DelvWardId,
		data.DelvZipCode,
		data.DelvIsSameAddr,
		data.InvWardId,
		data.InvZipCode,
		data.InvIsSameAddr,
		data.VerificationStatus,
		verifiedAt,
		data.VerifiedBy,
		data.TaxInvoiceForm,
		data.ObsType,
		data.CreditLimitAction,
		data.SalesInvLimitAction,
		data.ObsLimitAction,
		outletEstablishmentDate,
		data.DelvCity2,
		data.DelvLatitude,
		data.DelvLongitude,
		data.DelvLatitude2,
		data.DelvLongitude2,
		data.DelvWardId2,
		data.DelvZipCode2,
		isPkpOutletParam,
		// -- Parameter meta dan WHERE clause --
		data.UpdatedBy,
		data.CustId,
		data.OutletId,
	)

	return err
}

// BANK
func (r *outletRepositoryImpl) CheckBankExists(custId string, bankId int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_bank 
			  WHERE cust_id = $1 AND bank_id = $2`
	err := r.DB.Get(&exists, query, custId, bankId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckBankCodeDuplicate(custId string, bankId int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_bank 
			  WHERE cust_id = $1 AND bank_code = $2 AND bank_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, bankId)
	return exists, err
}

// OUTLET CLASS
func (r *outletRepositoryImpl) CheckOutletClassExists(custId string, otClassId int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_outlet_class 
			  WHERE cust_id = $1 AND ot_class_id = $2`
	err := r.DB.Get(&exists, query, custId, otClassId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckOutletClassCodeDuplicate(custId string, otClassId int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_outlet_class 
			  WHERE cust_id = $1 AND ot_class_code = $2 AND ot_class_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, otClassId)
	return exists, err
}

// OUTLET GROUP
func (r *outletRepositoryImpl) CheckOutletGroupExists(custId string, otGrpId int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_outlet_group 
			  WHERE cust_id = $1 AND ot_grp_id = $2`
	err := r.DB.Get(&exists, query, custId, otGrpId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckOutletGroupCodeDuplicate(custId string, otGrpId int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_outlet_group 
			  WHERE cust_id = $1 AND ot_grp_code = $2 AND ot_grp_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, otGrpId)
	return exists, err
}

// OUTLET LOC
func (r *outletRepositoryImpl) CheckOutletLocExists(custId string, otLocId int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_outlet_loc 
			  WHERE cust_id = $1 AND ot_loc_id = $2`
	err := r.DB.Get(&exists, query, custId, otLocId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckOutletLocCodeDuplicate(custId string, otLocId int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_outlet_loc 
			  WHERE cust_id = $1 AND ot_loc_code = $2 AND ot_loc_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, otLocId)
	return exists, err
}

// OUTLET TYPE
func (r *outletRepositoryImpl) CheckOutletTypeExists(custId string, otTypeId int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_outlet_type 
			  WHERE cust_id = $1 AND ot_type_id = $2`
	err := r.DB.Get(&exists, query, custId, otTypeId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckOutletTypeCodeDuplicate(custId string, otTypeId int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_outlet_type 
			  WHERE cust_id = $1 AND ot_type_code = $2 AND ot_type_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, otTypeId)
	return exists, err
}

// DISTRICT
func (r *outletRepositoryImpl) CheckDistrictExists(custId string, districtId int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_district 
			  WHERE cust_id = $1 AND district_id = $2`
	err := r.DB.Get(&exists, query, custId, districtId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckDistrictCodeDuplicate(custId string, districtId int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_district 
			  WHERE cust_id = $1 AND district_code = $2 AND district_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, districtId)
	return exists, err
}

// DISC GROUP
func (r *outletRepositoryImpl) CheckDiscGroupExists(custId string, discGrpId int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_disc_group 
			  WHERE cust_id = $1 AND disc_grp_id = $2`
	err := r.DB.Get(&exists, query, custId, discGrpId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckDiscGroupCodeDuplicate(custId string, discGrpId int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_disc_group 
			  WHERE cust_id = $1 AND disc_grp_code = $2 AND disc_grp_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, discGrpId)
	return exists, err
}

// MARKET
func (r *outletRepositoryImpl) CheckMarketExists(custId string, marketId int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_market 
			  WHERE cust_id = $1 AND market_id = $2`
	err := r.DB.Get(&exists, query, custId, marketId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckMarketCodeDuplicate(custId string, marketId int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_market 
			  WHERE cust_id = $1 AND market_code = $2 AND market_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, marketId)
	return exists, err
}

// INDUSTRY
func (r *outletRepositoryImpl) CheckIndustryExists(custId string, industryId int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_industry 
			  WHERE cust_id = $1 AND industry_id = $2`
	err := r.DB.Get(&exists, query, custId, industryId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckIndustryCodeDuplicate(custId string, industryId int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_industry 
			  WHERE cust_id = $1 AND industry_code = $2 AND industry_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, industryId)
	return exists, err
}

// PRICE GROUP
func (r *outletRepositoryImpl) CheckPriceGroupExists(custId string, priceGrpId int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
			  FROM mst.m_sp_price_group 
			  WHERE cust_id = $1 AND sp_price_grp_id = $2`
	err := r.DB.Get(&exists, query, custId, priceGrpId)
	return exists, err
}

func (r *outletRepositoryImpl) CheckPriceGroupCodeDuplicate(custId string, priceGrpId int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 
              FROM mst.m_sp_price_group 
              WHERE cust_id = $1 AND sp_price_grp_code = $2 AND sp_price_grp_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, priceGrpId)
	return exists, err
}

// FindBankIdByCode returns bank_id for given custId and bank_code
func (r *outletRepositoryImpl) FindBankIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT bank_id FROM mst.m_bank WHERE cust_id = $1 AND bank_code = $2 AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, code)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindOutletClassIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT ot_class_id FROM mst.m_outlet_class WHERE cust_id = $1 AND ot_class_code = $2 AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, code)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindOutletGroupIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT ot_grp_id FROM mst.m_outlet_group WHERE cust_id = $1 AND ot_grp_code = $2 AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, code)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindOutletLocIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT ot_loc_id FROM mst.m_outlet_loc WHERE cust_id = $1 AND ot_loc_code = $2 AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, code)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindOutletTypeIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT ot_type_id FROM mst.m_outlet_type WHERE cust_id = $1 AND ot_type_code = $2 AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, code)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindDistrictIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT district_id FROM mst.m_district WHERE cust_id = $1 AND district_code = $2 AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, code)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindDiscGroupIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT disc_grp_id FROM mst.m_disc_group WHERE cust_id = $1 AND disc_grp_code = $2 AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, code)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindMarketIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT market_id FROM mst.m_market WHERE cust_id = $1 AND market_code = $2 AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, code)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindIndustryIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT industry_id FROM mst.m_industry WHERE cust_id = $1 AND industry_code = $2 AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, code)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindPriceGroupIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT sp_price_grp_id FROM mst.m_sp_price_group WHERE cust_id = $1 AND sp_price_grp_code = $2 AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, code)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindBankIdByName(custId, name string) (int64, error) {
	var id int64
	query := `SELECT bank_id FROM mst.m_bank WHERE cust_id = $1 AND LOWER(bank_name) = LOWER($2) AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindOutletClassIdByName(custId, name string) (int64, error) {
	var id int64
	query := `SELECT ot_class_id FROM mst.m_outlet_class WHERE cust_id = $1 AND LOWER(ot_class_name) = LOWER($2) AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindOutletGroupIdByName(custId, name string) (int64, error) {
	var id int64
	query := `SELECT ot_grp_id FROM mst.m_outlet_group WHERE cust_id = $1 AND LOWER(ot_grp_name) = LOWER($2) AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindOutletLocIdByName(custId, name string) (int64, error) {
	var id int64
	query := `SELECT ot_loc_id FROM mst.m_outlet_loc WHERE cust_id = $1 AND LOWER(ot_loc_name) = LOWER($2) AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindOutletTypeIdByName(custId, name string) (int64, error) {
	var id int64
	query := `SELECT ot_type_id FROM mst.m_outlet_type WHERE cust_id = $1 AND LOWER(ot_type_name) = LOWER($2) AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindDistrictIdByName(custId, name string) (int64, error) {
	var id int64
	query := `SELECT district_id FROM mst.m_district WHERE cust_id = $1 AND LOWER(district_name) = LOWER($2) AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindDiscGroupIdByName(custId, name string) (int64, error) {
	var id int64
	query := `SELECT disc_grp_id FROM mst.m_disc_group WHERE cust_id = $1 AND LOWER(disc_grp_name) = LOWER($2) AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindMarketIdByName(custId, name string) (int64, error) {
	var id int64
	query := `SELECT market_id FROM mst.m_market WHERE cust_id = $1 AND LOWER(market_name) = LOWER($2) AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindMarketsByIDs(custId string, ids []int64) (map[int64]model.Market, error) {
	result := make(map[int64]model.Market)
	if len(ids) == 0 {
		return result, nil
	}

	query, args, err := sqlx.In(`SELECT market_id, market_code, market_name FROM mst.m_market WHERE cust_id = ? AND market_id IN (?)`, custId, ids)
	if err != nil {
		return result, err
	}
	query = r.DB.Rebind(query)

	var markets []model.Market
	if err := r.DB.Select(&markets, query, args...); err != nil {
		return result, err
	}

	for _, mkt := range markets {
		result[int64(mkt.MarketId)] = mkt
	}

	return result, nil
}

func (r *outletRepositoryImpl) FindIndustryIdByName(custId, name string) (int64, error) {
	var id int64
	query := `SELECT industry_id FROM mst.m_industry WHERE cust_id = $1 AND LOWER(industry_name) = LOWER($2) AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *outletRepositoryImpl) FindPriceGroupIdByName(custId, name string) (int64, error) {
	var id int64
	query := `SELECT sp_price_grp_id FROM mst.m_sp_price_group WHERE cust_id = $1 AND LOWER(sp_price_grp_name) = LOWER($2) AND is_active = true LIMIT 1`
	err := r.DB.Get(&id, query, custId, name)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// BANK
func (r *outletRepositoryImpl) UpdateImportBank(custId, code, name string, bankId int64) error {
	query := `
        UPDATE mst.m_bank
        SET bank_code = $1, bank_name = $2
        WHERE cust_id = $3 AND bank_id = $4
    `
	res, err := r.DB.Exec(query, code, name, custId, bankId)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no bank record updated (cust_id=%s, bank_id=%d)", custId, bankId)
	}
	return nil
}

// OUTLET CLASS
func (r *outletRepositoryImpl) UpdateImportOutletClass(custId, code, name string, otClassId int64) error {
	query := `
        UPDATE mst.m_outlet_class
        SET ot_class_code = $1, ot_class_name = $2
        WHERE cust_id = $3 AND ot_class_id = $4
    `
	res, err := r.DB.Exec(query, code, name, custId, otClassId)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no outlet class updated (cust_id=%s, ot_class_id=%d)", custId, otClassId)
	}
	return nil
}

// OUTLET GROUP
func (r *outletRepositoryImpl) UpdateImportOutletGroup(custId, code, name string, otGrpId int64) error {
	query := `
        UPDATE mst.m_outlet_group
        SET ot_grp_code = $1, ot_grp_name = $2
        WHERE cust_id = $3 AND ot_grp_id = $4
    `
	res, err := r.DB.Exec(query, code, name, custId, otGrpId)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no outlet group updated (cust_id=%s, ot_grp_id=%d)", custId, otGrpId)
	}
	return nil
}

// OUTLET LOC
func (r *outletRepositoryImpl) UpdateImportOutletLoc(custId, code, name string, otLocId int64) error {
	query := `
        UPDATE mst.m_outlet_loc
        SET ot_loc_code = $1, ot_loc_name = $2
        WHERE cust_id = $3 AND ot_loc_id = $4
    `
	res, err := r.DB.Exec(query, code, name, custId, otLocId)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no outlet loc updated (cust_id=%s, ot_loc_id=%d)", custId, otLocId)
	}
	return nil
}

// OUTLET TYPE
func (r *outletRepositoryImpl) UpdateImportOutletType(custId, code, name string, otTypeId int64) error {
	query := `
        UPDATE mst.m_outlet_type
        SET ot_type_code = $1, ot_type_name = $2
        WHERE cust_id = $3 AND ot_type_id = $4
    `
	res, err := r.DB.Exec(query, code, name, custId, otTypeId)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no outlet type updated (cust_id=%s, ot_type_id=%d)", custId, otTypeId)
	}
	return nil
}

// DISTRICT
func (r *outletRepositoryImpl) UpdateImportDistrict(custId, code, name string, districtId int64) error {
	query := `
        UPDATE mst.m_district
        SET district_code = $1, district_name = $2
        WHERE cust_id = $3 AND district_id = $4
    `
	res, err := r.DB.Exec(query, code, name, custId, districtId)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no district updated (cust_id=%s, district_id=%d)", custId, districtId)
	}
	return nil
}

// DISC GROUP
func (r *outletRepositoryImpl) UpdateImportDiscGroup(custId, code, name string, discGrpId int64) error {
	query := `
        UPDATE mst.m_disc_group
        SET disc_grp_code = $1, disc_grp_name = $2
        WHERE cust_id = $3 AND disc_grp_id = $4
    `
	res, err := r.DB.Exec(query, code, name, custId, discGrpId)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no disc group updated (cust_id=%s, disc_grp_id=%d)", custId, discGrpId)
	}
	return nil
}

// MARKET
func (r *outletRepositoryImpl) UpdateImportMarket(custId, code, name string, marketId int64) error {
	query := `
        UPDATE mst.m_market
        SET market_code = $1, market_name = $2
        WHERE cust_id = $3 AND market_id = $4
    `
	res, err := r.DB.Exec(query, code, name, custId, marketId)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no market updated (cust_id=%s, market_id=%d)", custId, marketId)
	}
	return nil
}

// INDUSTRY
func (r *outletRepositoryImpl) UpdateImportIndustry(custId, code, name string, industryId int64) error {
	query := `
        UPDATE mst.m_industry
        SET industry_code = $1, industry_name = $2
        WHERE cust_id = $3 AND industry_id = $4
    `
	res, err := r.DB.Exec(query, code, name, custId, industryId)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("no industry updated (cust_id=%s, industry_id=%d)", custId, industryId)
	}
	return nil
}

// Price group

func (r *outletRepositoryImpl) UpdateImportPriceGroup(custId string, priceGrpId int64, code, name string) (int64, error) {
	var id int64
	query := `
		UPDATE mst.m_sp_price_group
		SET sp_price_grp_code = $1, sp_price_grp_name = $2
		WHERE cust_id = $3 AND sp_price_grp_id = $4
		RETURNING sp_price_grp_id
	`
	err := r.DB.QueryRow(query, code, name, custId, priceGrpId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repository *outletRepositoryImpl) GetRegionAndDistributorCodeByCustId(custId, parentCustId string) (regionCode *string, distributorCode *string, err error) {
	var rCode, dCode sql.NullString
	query := `WITH c AS (
			SELECT cust_id, COALESCE(NULLIF($2, ''), cust_id) AS parent_cust_id, distributor_id
			FROM smc.m_customer
			WHERE cust_id = $1
		)
		SELECT mr.region_code, md.distributor_code
		FROM c
		LEFT JOIN LATERAL (
			SELECT d.distributor_code, d.region_id, d.cust_id
			FROM mst.m_distributor d
			WHERE d.distributor_id = c.distributor_id
			  AND d.cust_id IN (c.cust_id, c.parent_cust_id)
			ORDER BY CASE WHEN d.cust_id = c.cust_id THEN 0 ELSE 1 END
			LIMIT 1
		) md ON true
		LEFT JOIN LATERAL (
			SELECT r.region_code, r.cust_id
			FROM mst.m_region r
			WHERE r.region_id = md.region_id
			  AND r.cust_id IN (c.cust_id, c.parent_cust_id)
			ORDER BY CASE WHEN r.cust_id = c.cust_id THEN 0 ELSE 1 END
			LIMIT 1
		) mr ON true`
	err = repository.QueryRow(query, custId, parentCustId).Scan(&rCode, &dCode)
	if err != nil {
		return nil, nil, err
	}
	if rCode.Valid {
		rc := rCode.String
		regionCode = &rc
	}
	if dCode.Valid {
		dc := dCode.String
		distributorCode = &dc
	}
	return regionCode, distributorCode, nil
}

// GetNextOutletPrincipalCodeSeq returns the next 4-digit sequence (1-9999) for the given prefix.
func (repository *outletRepositoryImpl) GetNextOutletPrincipalCodeSeq(prefix string) (seq int, err error) {
	tx, err := repository.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	seq, err = repository.GetNextOutletPrincipalCodeSeqTx(tx, prefix)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return seq, nil
}

func (repository *outletRepositoryImpl) GetNextOutletPrincipalCodeSeqTx(tx *sqlx.Tx, prefix string) (seq int, err error) {
	query := `
		INSERT INTO mst.m_outlet_principal_code_seq(prefix, last_sequence_no, created_at, updated_at)
		VALUES ($1, 1, NOW(), NOW())
		ON CONFLICT (prefix)
		DO UPDATE SET
			last_sequence_no = mst.m_outlet_principal_code_seq.last_sequence_no + 1,
			updated_at = NOW()
		RETURNING last_sequence_no`
	if err := tx.QueryRow(query, prefix).Scan(&seq); err != nil {
		return 0, err
	}
	if seq > 9999 {
		return 0, errors.New("outlet_principal_code sequence exceeded 9999 for today")
	}
	return seq, nil
}

func (repository *outletRepositoryImpl) FindOneCustomerByDistributorID(distributorID int64) (model.MCustomer, error) {
	mCustomer := model.MCustomer{}
	query := `SELECT 
				cust_id, cust_name, parent_cust_id
			  FROM smc.m_customer
			  WHERE distributor_id = $1`
	err := repository.Get(&mCustomer, query, distributorID)
	if err != nil {
		log.Println("outletRepositoryImpl, FindOneCustomerByDistributorID, err:", err.Error())
		return mCustomer, err
	}

	return mCustomer, nil
}

func (repository *outletRepositoryImpl) FindCustIdsByDistributorIds(parentCustId string, distributorIds []int) ([]string, error) {
	if len(distributorIds) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In(`SELECT cust_id FROM smc.m_customer WHERE distributor_id IN (?) AND parent_cust_id = ?`, distributorIds, parentCustId)
	if err != nil {
		return nil, err
	}
	query = repository.Rebind(query)
	var rows []struct {
		CustId string `db:"cust_id"`
	}
	err = repository.Select(&rows, query, args...)
	if err != nil {
		log.Println("outletRepositoryImpl, FindCustIdsByDistributorIds, err:", err.Error())
		return nil, err
	}
	custIds := make([]string, 0, len(rows))
	for _, r := range rows {
		custIds = append(custIds, r.CustId)
	}
	return custIds, nil
}

func (repository *outletRepositoryImpl) FindCustIdsByParentCustId(parentCustId string) ([]string, error) {
	trimmedParentCustID := strings.TrimSpace(parentCustId)
	if trimmedParentCustID == "" {
		return nil, nil
	}

	query := `SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = $1`
	var rows []struct {
		CustId string `db:"cust_id"`
	}
	err := repository.Select(&rows, query, trimmedParentCustID)
	if err != nil {
		log.Println("outletRepositoryImpl, FindCustIdsByParentCustId, err:", err.Error())
		return nil, err
	}

	custIds := make([]string, 0, len(rows))
	for _, row := range rows {
		custIds = append(custIds, row.CustId)
	}

	return custIds, nil
}

func (repository *outletRepositoryImpl) ExistsOutletContactIdentity(custId, parentCustId, identityType, identityNo string, excludeOutletId int64) (bool, error) {
	if strings.TrimSpace(identityType) == "" || strings.TrimSpace(identityNo) == "" {
		return false, nil
	}
	query := `SELECT 1 FROM mst.m_outlet_contact WHERE cust_id = $1 AND TRIM(identity_type) = $2 AND TRIM(identity_no) = $3`
	args := []interface{}{custId, strings.TrimSpace(identityType), strings.TrimSpace(identityNo)}
	if excludeOutletId > 0 {
		query += ` AND outlet_id != $4`
		args = append(args, excludeOutletId)
	}
	query += ` LIMIT 1`
	var one int
	err := repository.QueryRow(query, args...).Scan(&one)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (repository *outletRepositoryImpl) BulkUpdateStatuses(ctx context.Context) (int64, error) {
	const q = `
WITH base AS (
	SELECT
		o.outlet_id,
		o.cust_id,
		o.outlet_status,
		COALESCE(o.pre_dormant_status, 6) AS pre_status,
		CASE WHEN o.last_trans_date IS NULL THEN NULL ELSE (CURRENT_DATE - o.last_trans_date)::int END AS dslt_trans,
		(CURRENT_DATE - COALESCE(o.last_trans_date, o.created_at::date))::int AS dslt_overall
	FROM mst.m_outlet o
	WHERE o.is_del = false
		AND o.outlet_status IN (1, 5, 6, 7)
		AND (o.last_trans_date IS NULL OR o.outlet_status = 5)
),
cfg AS (
	SELECT
		b.*,
		COALESCE(NULLIF(cfg5.counting_period, 0), 10) AS cp_to_5,
		COALESCE(cfg5.validate_trx, true) AS vt_to_5,
		COALESCE(NULLIF(cfg1.counting_period, 0), 1) AS cp_to_1,
		COALESCE(cfg1.validate_trx, true) AS vt_to_1,
		COALESCE(NULLIF(cfg7.counting_period, 0), 5) AS cp_to_7,
		COALESCE(cfg7.validate_trx, true) AS vt_to_7,
		COALESCE(NULLIF(cfgpre.counting_period, 0), 1) AS cp_to_pre,
		COALESCE(cfgpre.validate_trx, true) AS vt_to_pre
	FROM base b
	LEFT JOIN LATERAL (
		SELECT d.counting_period, d.validate_trx
		FROM mst.m_outlet_config_det d
		WHERE d.is_del = false
			AND d.status = 5
			AND substring(b.cust_id, 1, length(d.cust_id)) = d.cust_id
		ORDER BY length(d.cust_id) DESC, d.outlet_config_det_id DESC
		LIMIT 1
	) cfg5 ON true
	LEFT JOIN LATERAL (
		SELECT d.counting_period, d.validate_trx
		FROM mst.m_outlet_config_det d
		WHERE d.is_del = false
			AND d.status = 1
			AND substring(b.cust_id, 1, length(d.cust_id)) = d.cust_id
		ORDER BY length(d.cust_id) DESC, d.outlet_config_det_id DESC
		LIMIT 1
	) cfg1 ON true
	LEFT JOIN LATERAL (
		SELECT d.counting_period, d.validate_trx
		FROM mst.m_outlet_config_det d
		WHERE d.is_del = false
			AND d.status = 7
			AND substring(b.cust_id, 1, length(d.cust_id)) = d.cust_id
		ORDER BY length(d.cust_id) DESC, d.outlet_config_det_id DESC
		LIMIT 1
	) cfg7 ON true
	LEFT JOIN LATERAL (
		SELECT d.counting_period, d.validate_trx
		FROM mst.m_outlet_config_det d
		WHERE d.is_del = false
			AND d.status = b.pre_status
			AND substring(b.cust_id, 1, length(d.cust_id)) = d.cust_id
		ORDER BY length(d.cust_id) DESC, d.outlet_config_det_id DESC
		LIMIT 1
	) cfgpre ON true
),
decision AS (
	SELECT
		c.*,
		CASE
			WHEN c.outlet_status = 5
				AND c.dslt_trans IS NOT NULL
				AND c.dslt_trans <= c.cp_to_pre
				AND c.vt_to_pre
				THEN c.pre_status
			WHEN c.outlet_status = 6
				AND c.dslt_trans IS NOT NULL
				AND c.dslt_trans <= c.cp_to_1
				AND c.vt_to_1
				THEN 1
			WHEN c.outlet_status = 6
				AND c.dslt_overall > c.cp_to_5
				THEN 5
			WHEN c.outlet_status = 1
				AND c.dslt_trans IS NOT NULL
				AND c.dslt_trans <= c.cp_to_7
				AND c.vt_to_7
				THEN 7
			WHEN c.outlet_status = 1
				AND c.dslt_overall > c.cp_to_5
				THEN 5
			ELSE c.outlet_status
		END AS new_status
	FROM cfg c
),
updated AS (
	UPDATE mst.m_outlet o
	SET
		outlet_status = d.new_status,
		pre_dormant_status = CASE
			WHEN d.new_status = 5 THEN o.outlet_status::int2
			WHEN o.outlet_status = 5 AND d.new_status <> 5 THEN NULL
			ELSE o.pre_dormant_status
		END,
		updated_at = NOW()
	FROM decision d
	WHERE o.outlet_id = d.outlet_id
		AND d.new_status <> d.outlet_status
	RETURNING o.outlet_id
)
SELECT COUNT(*)::bigint AS rows_affected FROM updated;
`
	var rowsAffected int64
	if err := repository.QueryRowxContext(ctx, q).Scan(&rowsAffected); err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (repository *outletRepositoryImpl) BulkPromoteRegisteredWithTransToNoo(ctx context.Context) (int64, error) {
	const q = `
UPDATE mst.m_outlet o
SET
	outlet_status = 1,
	updated_at = NOW()
WHERE o.is_del = false
	AND o.outlet_status = 6
	AND o.last_trans_date IS NOT NULL
`
	result, err := repository.ExecContext(ctx, q)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (repository *outletRepositoryImpl) UpdateOutletStatus(outletId int64, custId string, status int, updatedBy int64) error {
	query := `UPDATE mst.m_outlet SET outlet_status = $1, updated_by = $2, updated_at = now()`
	args := []interface{}{status, updatedBy}
	if status == 4 {
		query += `, closed_date = now()`
	}
	query += ` WHERE outlet_id = $3 AND cust_id = $4`
	args = append(args, outletId, custId)
	result, err := repository.Exec(query, args...)
	if err != nil {
		log.Println("outletRepository, UpdateOutletStatus, err:", err.Error())
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("outlet not found or access denied")
	}
	return nil
}

// FindAllOutletCrByStatus retrieves a paginated list of outlet change requests filtered by status and customer ID.
// It joins with m_outlet to get outlet details and outlet_cr_det to get latitude/longitude changes.
// Returns the list of outlet change requests, total count, last page number, and any error encountered.
func (repository *outletRepositoryImpl) FindAllOutletCrByStatus(dataFilter entity.OutletListApprovalQueryFilter, custId string) ([]model.OutletCrList, int, int, error) {
	outlets := []model.OutletCrList{}

	selectCount := ` COUNT(DISTINCT ocr.outlet_cr_id) AS total `

	selectField := `
		ocr.outlet_cr_id,
		ocr.outlet_id,
		mo.outlet_code,
		mo.outlet_name,
		ocd_long.old_value AS current_long,
		ocd_lat.old_value AS current_lat,
		ocd_long.new_value AS new_long,
		ocd_lat.new_value AS new_lat,
		ocr.source,
		ocr.status,
		CASE ocr.status 
			WHEN 1 THEN 'Need Review'
			WHEN 2 THEN 'Approved'
			WHEN 3 THEN 'Rejected'
		END AS status_desc,
		mu.user_name AS request_by,
		ocr.created_at AS request_date
	`

	qFrom := `
		FROM mst.outlet_cr ocr
		LEFT JOIN mst.m_outlet mo ON mo.outlet_id = ocr.outlet_id AND mo.cust_id = ocr.cust_id AND mo.is_del = false
		LEFT JOIN sys.m_user mu ON mu.user_id = ocr.created_by
		LEFT JOIN mst.outlet_cr_det ocd_long ON ocd_long.outlet_cr_id = ocr.outlet_cr_id AND ocd_long.field_name = 'longitude'
		LEFT JOIN mst.outlet_cr_det ocd_lat ON ocd_lat.outlet_cr_id = ocr.outlet_cr_id AND ocd_lat.field_name = 'latitude'
	`

	qWhere := ` WHERE ocr.cust_id = $1`
	args := []interface{}{custId}
	argIndex := 2

	if dataFilter.Status != "" {
		statusStrings := strings.Split(dataFilter.Status, ",")
		statusValues := []int{}
		for _, s := range statusStrings {
			trimmed := strings.TrimSpace(s)
			if val, err := strconv.Atoi(trimmed); err == nil {
				statusValues = append(statusValues, val)
			}
		}
		if len(statusValues) > 0 {
			// Build IN clause with dynamic placeholders
			placeholders := make([]string, len(statusValues))
			for i, v := range statusValues {
				placeholders[i] = "$" + strconv.Itoa(argIndex)
				args = append(args, v)
				argIndex++
			}
			qWhere += ` AND ocr.status IN (` + strings.Join(placeholders, ", ") + `)`
		}
	}

	if dataFilter.Q != "" {
		qWhere += ` AND (mo.outlet_code ILIKE $` + strconv.Itoa(argIndex) + ` OR mo.outlet_name ILIKE $` + strconv.Itoa(argIndex) + `)`
		searchPattern := "%" + dataFilter.Q + "%"
		args = append(args, searchPattern)
		argIndex++
	}

	qOrderBy := ` ORDER BY ocr.created_at DESC`
	if dataFilter.Sort != "" {
		sortParts := strings.Split(dataFilter.Sort, ":")
		if len(sortParts) == 2 {
			field := sortParts[0]
			direction := strings.ToUpper(sortParts[1])
			if direction == "ASC" || direction == "DESC" {
				if field == "created_date" {
					qOrderBy = ` ORDER BY ocr.created_at ` + direction
				}
			}
		}
	}

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	var total int
	err := repository.Get(&total, queryCount, args...)
	if err != nil {
		return outlets, total, 0, err
	}

	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	limit := dataFilter.Limit
	if limit < 1 {
		limit = 5
	}

	offset := (page - 1) * limit
	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	querySelect := `SELECT ` + selectField + qFrom + qWhere + qOrderBy + ` LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	err = repository.Select(&outlets, querySelect, args...)
	if err != nil {
		return outlets, total, lastPage, err
	}

	return outlets, total, lastPage, nil
}

// UpdateOutletCrStatus updates the status of multiple outlet change requests.
// It sets the status, updated_by, updated_at, approval_by, and approval_at fields.
// Uses parameterized query with sqlx.In to prevent SQL injection.
func (repository *outletRepositoryImpl) UpdateOutletCrStatus(outletCrIds []int, status int, updatedBy int64) error {
	if len(outletCrIds) == 0 {
		return errors.New("outlet_cr_id is required")
	}

	query, args, err := sqlx.In(
		`UPDATE mst.outlet_cr 
		SET status = ?, updated_by = ?, updated_at = CURRENT_TIMESTAMP, approval_by = ?, approval_at = CURRENT_TIMESTAMP 
		WHERE outlet_cr_id IN (?)`,
		status, updatedBy, updatedBy, outletCrIds,
	)
	if err != nil {
		return err
	}

	query = repository.Rebind(query)
	_, err = repository.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// GetOutletCrDetByOutletCrIds retrieves outlet change request details for latitude and longitude fields.
// Returns details filtered by outlet_cr_id and field_name (latitude/longitude).
// Uses parameterized query with sqlx.In to prevent SQL injection.
func (repository *outletRepositoryImpl) GetOutletCrDetByOutletCrIds(outletCrIds []int) ([]model.OutletCrDet, error) {
	if len(outletCrIds) == 0 {
		return []model.OutletCrDet{}, nil
	}

	query, args, err := sqlx.In(
		`SELECT outlet_cr_det_id, outlet_cr_id, field_name, old_value, new_value
		FROM mst.outlet_cr_det
		WHERE outlet_cr_id IN (?)
		AND field_name IN (?)`,
		outletCrIds, []string{"latitude", "longitude"},
	)
	if err != nil {
		return nil, err
	}

	query = repository.Rebind(query)
	var details []model.OutletCrDet
	err = repository.Select(&details, query, args...)
	if err != nil {
		return details, err
	}

	return details, nil
}

// GetOutletCrByOutletCrIds retrieves outlet change request records by their IDs.
// Returns all fields including cust_id, status, approval information, and timestamps.
// Uses parameterized query with sqlx.In to prevent SQL injection.
func (repository *outletRepositoryImpl) GetOutletCrByOutletCrIds(outletCrIds []int) ([]model.OutletCr, error) {
	if len(outletCrIds) == 0 {
		return []model.OutletCr{}, nil
	}

	query, args, err := sqlx.In(
		`SELECT cust_id, outlet_cr_id, outlet_id, source, status, created_by, created_at, updated_by, updated_at, approval_by, approval_at
		FROM mst.outlet_cr
		WHERE outlet_cr_id IN (?)`,
		outletCrIds,
	)
	if err != nil {
		log.Println("outletRepository, GetOutletCrByOutletCrIds, sqlx.In err:", err.Error())
		return nil, err
	}

	query = repository.Rebind(query)
	var outletCrs []model.OutletCr
	err = repository.Select(&outletCrs, query, args...)
	if err != nil {
		return outletCrs, err
	}

	return outletCrs, nil
}

// UpdateOutletLocationFromCr updates the latitude and longitude of an outlet based on approved change request.
// Only updates non-deleted outlets (is_del = false).
func (repository *outletRepositoryImpl) UpdateOutletLocationFromCr(outletId int64, latitude, longitude string) error {
	query := `UPDATE mst.m_outlet 
		SET latitude = $1, longitude = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE outlet_id = $3 AND is_del = false`

	_, err := repository.Exec(query, latitude, longitude, outletId)
	if err != nil {
		return err
	}

	return nil
}
