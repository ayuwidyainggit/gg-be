package repository

import (
	"encoding/json"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/constant"
	"master/pkg/sql_helper"
	"master/pkg/str"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

func NewDistributorRepository(db *sqlx.DB) *distributorRepositoryImpl {
	return &distributorRepositoryImpl{db}
}

type distributorRepositoryImpl struct {
	*sqlx.DB
}
type distributorTransaction struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

func buildDistributorScopeCondition(dataFilter entity.DistributorQueryFilter, custId string) string {
	if dataFilter.JwtDistributorId > 0 {
		return ` AND (sp.distributor_id = ` + strconv.FormatInt(dataFilter.JwtDistributorId, 10) + ` OR sp.cust_id LIKE '` + custId + `%' ) `
	}

	return ` AND sp.cust_id LIKE '` + custId + `%' `
}

func buildDistributorActiveCondition(isActive *int) string {
	if isActive != nil && *isActive == 2 {
		return ` AND sp.is_active = false `
	}

	return ` AND sp.is_active = true `
}

type DistributorRepository interface {
	TrxBegin() (*distributorTransaction, error)
	FindAllByCustId(dataFilter entity.DistributorQueryFilter, custId string) ([]model.DistributorList, int, int, error)
	FindAllLookupByCustId(dataFilter entity.DistributorQueryFilter, custId string) ([]model.DistributorList, int, int, error)
	FindAllAreaRegionByIDs(distIDs []int64, custID string) ([]model.DistributorAreaRegionDetail, error)
	FindAllAreaRegionByCodes(distCodes []string, custID string) ([]model.DistributorAreaRegionDetail, error)
	FindOneByDistributorIdAndCustId(params entity.DetailDistributorParams) (model.DistributorList, error)
	FindAllowManagePricingByDistributorID(distributorID int64) (bool, error)
	FindAllowUploadSecondarySalesByDistributorID(distributorID int64) (bool, error)
	FindDistributorIdByCustId(custId string) (int64, error)
	FindOneByDistributorCodeAndCustId(distributorId string, custId string) (model.DistributorList, error)
	FindOneByParentCustIdAndDistributorCode(custId, distributorCode string) (model.DistributorWithCustomer, error)
	FindAllDistributorContactByDistIdAndCustId(params entity.DetailDistributorParams) ([]model.DistributorContact, error)
	FindAllDistributorTaxByDistIdAndCustId(params entity.DetailDistributorParams) ([]model.DistributorTax, error)
	Delete(custId string, distributorId int64, deletedBy int64) error
	FindAllByCustIdWithCustomer(dataFilter entity.DistributorQueryFilter, custId string) ([]model.DistributorList, int, int, error)
}

func NewsDistributorRepository(db *sqlx.DB) *distributorRepositoryImpl {
	return &distributorRepositoryImpl{db}
}
func NewTransactionDistributor(db *sqlx.DB) (trxObj *distributorTransaction, err error) {
	trx := db.MustBegin()

	return &distributorTransaction{tx: trx, db: db}, nil
}

func (repo *distributorRepositoryImpl) TrxBegin() (*distributorTransaction, error) {
	trxObj, err := NewTransactionDistributor(repo.DB)
	if err != nil {
		return nil, err
	}
	return trxObj, nil
}

func (repo *distributorTransaction) TrxCommit() error {
	return repo.tx.Commit()
}

func (repo *distributorTransaction) TrxRollback() error {
	return repo.tx.Rollback()
}

func (repo *distributorTransaction) InsertDistributor(distributor *model.Distributor) error {
	query :=
		`INSERT INTO mst.m_distributor(
			cust_id, distributor_code, distributor_name, barcode, region_id,
			area_id, channel_id, sub_distributor_group_id, dist_price_grp_id, address, province_id,
			regency_id, sub_district_id, ward_id, zip_code, ot_loc_id,
			latitude, longitude, is_active, created_by, created_at, updated_by, updated_at,
			is_del, deleted_by, deleted_at, phone, fax_number,
			allow_add_product, allow_edit_product, allow_manage_pricing, allow_upload_secondary_sales
	)
	VALUES ( 
		$1, $2, $3, $4, $5, 
		$6, 
		$7, $8, $9, $10, $11,
		$12, $13, $14, $15, $16, 
		$17, $18, $19, $20, $21, 
		$22, $23, $24, $25, 
		$26, $27, $28,
		$29, $30, $31, $32
	) RETURNING distributor_id;`
	lastInsertId := distributor.DistributorId
	err := repo.tx.QueryRow(query,
		distributor.CustId, distributor.DistributorCode, distributor.DistributorName, distributor.Barcode, distributor.RegionId,
		distributor.AreaId, distributor.ChannelId, distributor.SubDistributorGroupId, distributor.DistPriceGrpId, distributor.Address,
		distributor.ProvinceId, distributor.RegencyId, distributor.SubDistrictId, distributor.WardId, distributor.ZipCode, distributor.OtLocId,
		distributor.Latitude, distributor.Longitude, distributor.IsActive, distributor.CreatedBy, distributor.CreatedAt, distributor.UpdatedBy, distributor.UpdatedAt,
		distributor.IsDel, distributor.DeletedBy, distributor.DeletedAt, distributor.Phone, distributor.FaxNumber,
		distributor.AllowAddProduct, distributor.AllowEditProduct, distributor.AllowManagePricing, distributor.AllowUploadSecondarySales).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	distributor.DistributorId = lastInsertId
	return nil
}

func (repo *distributorTransaction) InsertDistributorContact(disributorId int64, distributorContact *model.DistributorContact) error {
	query :=
		`INSERT INTO mst.m_distributor_contact(
		cust_id, distributor_id, contact_name,
		job_title, phone_no, is_wa_no, wa_no, email, identity_no, identity_type
	)
	VALUES ( 
		$1, $2, $3, $4, 
		$5, $6, $7, $8, $9, $10
	) RETURNING distributor_contact_id;`

	lastInsertId := distributorContact.DistributorContactId
	err := repo.tx.QueryRow(query, distributorContact.CustId, disributorId, distributorContact.ContactName,
		distributorContact.JobTitle, distributorContact.PhoneNo, distributorContact.IsWaNo, distributorContact.WaNo, distributorContact.Email, distributorContact.IdentityNo, distributorContact.IdentityType).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	distributorContact.DistributorContactId = lastInsertId
	return nil
}

func (repo *distributorTransaction) InsertDistributorTax(disributorId int64, distributorTax *model.DistributorTax) error {
	query :=
		`INSERT INTO mst.m_distributor_tax(
		cust_id, distributor_id, tax_identifier_no_type,
		tax_identifier_no, nitku, tax_name, tax_address
	)
	VALUES ( 
		$1, $2, $3, $4, 
		$5, $6, $7
	) RETURNING distributor_tax_id;`

	lastInsertId := distributorTax.DistributorTaxId
	err := repo.tx.QueryRow(query, distributorTax.CustId, disributorId, distributorTax.TaxIdentifierNoType,
		distributorTax.TaxIdentifierNo, distributorTax.Nitku, distributorTax.TaxName, distributorTax.TaxAddress).Scan(&lastInsertId)
	if err != nil {
		return err
	}
	distributorTax.DistributorTaxId = lastInsertId
	return nil
}

func (repository *distributorRepositoryImpl) FindAllByCustId(dataFilter entity.DistributorQueryFilter, custId string) ([]model.DistributorList, int, int, error) {
	mdistributor := []model.DistributorList{}

	selectCount := ` COUNT(DISTINCT sp.distributor_id) AS total `
	selectField := ` DISTINCT ON (sp.distributor_id) sp.cust_id, COALESCE(sp.parent_cust_id, '') AS parent_cust_id, sp.distributor_id, sp.distributor_code, distributor_name, sp.barcode,
					COALESCE(sp.region_id, 0) AS region_id, COALESCE(sp.area_id, 0) AS area_id,
					COALESCE(sp.channel_id, 0) AS channel_id, COALESCE(sp.sub_distributor_group_id, 0) AS sub_distributor_group_id,
					COALESCE(sp.dist_price_grp_id, 0) AS dist_price_grp_id, COALESCE(sp.address, '') AS address,
					COALESCE(sp.province_id, '') AS province_id, COALESCE(sp.regency_id, '') AS regency_id,
					COALESCE(sp.sub_district_id, '') AS sub_district_id, COALESCE(sp.ward_id, '') AS ward_id,
					COALESCE(sp.zip_code, '') AS zip_code, COALESCE(sp.ot_loc_id, 0) AS ot_loc_id, COALESCE(sp.latitude, '') AS latitude,
					COALESCE(ar.area_code, '') AS area_code, COALESCE(ar.area_name, '') AS area_name,
					COALESCE(pv.province_id, '') AS province_code, COALESCE(pv.province, '') AS province_name,
					COALESCE(rgc.regency_id, '') AS regency_code, COALESCE(rgc.regency, '') AS regency_name,
					COALESCE(sd.sub_district_id, '') AS sub_district_code, COALESCE(sd.sub_district, '') AS sub_district_name,
					COALESCE(wr.ward_id, '') AS ward_code, COALESCE(wr.ward, '') AS ward_name,
					COALESCE(sp.longitude, '') AS longitude, COALESCE(rgn.region_code, '') AS region_code, COALESCE(rgn.region_name, '') AS region_name, COALESCE(cnl.channel_code, '') AS channel_code, COALESCE(cnl.channel_name, '') AS channel_name, COALESCE(dpg.dist_price_grp_code, '') AS dist_price_grp_code, COALESCE(dpg.dist_price_grp_name, '') AS dist_price_grp_name, COALESCE(sp.is_active, false) AS is_active,
					COALESCE(u.user_fullname, '') AS updated_by_name `
	qWhere := ` WHERE sp.is_del = false `
	qWhere += buildDistributorScopeCondition(dataFilter, custId)
	qWhere += buildDistributorActiveCondition(dataFilter.IsActive)

	if dataFilter.Query != "" {
		qWhere += ` AND (CAST(sp.distributor_id AS TEXT) ILIKE '%` + dataFilter.Query + `%' 
								OR sp.distributor_code ILIKE '%` + dataFilter.Query + `%'
								OR sp.distributor_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if len(dataFilter.AreaID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.AreaID, ",")
		qWhere += ` AND sp.area_id IN (` + intArrStr + `) `
	}

	qFrom := ` 	FROM mst.m_distributor sp
				LEFT JOIN mst.m_dist_price_group dpg ON dpg.dist_price_grp_id = sp.dist_price_grp_id
				LEFT JOIN mst.m_channel cnl ON cnl.channel_id = sp.channel_id
				LEFT JOIN mst.m_region rgn ON rgn.region_id = sp.region_id
				
				LEFT JOIN mst.m_area ar ON ar.area_id = sp.area_id
				LEFT JOIN mst.m_province pv ON pv.province_id = sp.province_id
				LEFT JOIN mst.m_regency rgc ON rgc.regency_id = sp.regency_id
				LEFT JOIN mst.m_sub_district sd ON sd.sub_district_id = sp.sub_district_id
				LEFT JOIN mst.m_ward wr ON wr.ward_id = sp.ward_id

				LEFT JOIN sys.m_user u ON u.user_id = sp.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		return mdistributor, 0, 0, err
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
		querySelect += fmt.Sprintf(` ORDER BY sp.distributor_id DESC, %s`, sortBy)
	} else {
		querySelect += ` ORDER BY sp.distributor_id DESC `
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

	err = repository.Select(&mdistributor, querySelect)
	if err != nil {
		return mdistributor, total, lastPage, err
	}

	return mdistributor, total, lastPage, nil
}

func (repository *distributorRepositoryImpl) FindAllAreaRegionByIDs(distIDs []int64, custID string) ([]model.DistributorAreaRegionDetail, error) {
	distributors := []model.DistributorAreaRegionDetail{}
	if len(distIDs) == 0 {
		return distributors, nil
	}

	query, args, err := sqlx.In(`
		SELECT
			d.distributor_id,
			d.distributor_code,
			d.distributor_name,
			d.region_id,
			COALESCE(r.region_name, '') AS region_name,
			d.area_id,
			COALESCE(a.area_name, '') AS area_name
		FROM mst.m_distributor d
		LEFT JOIN mst.m_region r ON r.region_id = d.region_id
		LEFT JOIN mst.m_area a ON a.area_id = d.area_id AND a.cust_id = ?
		WHERE d.distributor_id IN (?)
		  AND d.parent_cust_id = ?
		  AND d.is_del = false
	`, custID, distIDs, custID)
	if err != nil {
		log.Error("distributorRepositoryImpl, FindAllAreaRegionByIDs, err:", err.Error())
		return nil, err
	}

	query = repository.Rebind(query)
	if err := repository.Select(&distributors, query, args...); err != nil {
		log.Error("distributorRepositoryImpl, FindAllAreaRegionByIDs, err:", err.Error())
		return distributors, err
	}

	return distributors, nil
}

func (repository *distributorRepositoryImpl) FindAllAreaRegionByCodes(distCodes []string, custID string) ([]model.DistributorAreaRegionDetail, error) {
	distributors := []model.DistributorAreaRegionDetail{}
	if len(distCodes) == 0 {
		return distributors, nil
	}

	query, args, err := sqlx.In(`
		SELECT
			d.distributor_id,
			d.distributor_code,
			d.distributor_name,
			d.region_id,
			COALESCE(r.region_name, '') AS region_name,
			d.area_id,
			COALESCE(a.area_name, '') AS area_name
		FROM mst.m_distributor d
		LEFT JOIN mst.m_region r ON r.region_id = d.region_id
		LEFT JOIN mst.m_area a ON a.area_id = d.area_id AND a.cust_id = ?
		WHERE d.distributor_code IN (?)
		  AND d.parent_cust_id = ?
		  AND d.is_del = false
	`, custID, distCodes, custID)
	if err != nil {
		return nil, err
	}

	query = repository.Rebind(query)
	if err := repository.Select(&distributors, query, args...); err != nil {
		log.Error("distributorRepositoryImpl, FindAllAreaRegionByCodes, err:", err.Error())
		return distributors, err
	}

	return distributors, nil
}

func (repository *distributorRepositoryImpl) FindAllLookupByCustId(dataFilter entity.DistributorQueryFilter, custId string) ([]model.DistributorList, int, int, error) {
	mdistributor := []model.DistributorList{}

	selectCount := ` COUNT(DISTINCT sp.distributor_id) AS total `
	selectField := ` sp.cust_id, COALESCE(sp.parent_cust_id, '') AS parent_cust_id, sp.distributor_id, sp.distributor_code, distributor_name, sp.barcode,
					COALESCE(sp.region_id, 0) AS region_id, COALESCE(sp.area_id, 0) AS area_id,
					COALESCE(sp.channel_id, 0) AS channel_id, COALESCE(sp.sub_distributor_group_id, 0) AS sub_distributor_group_id,
					COALESCE(sp.dist_price_grp_id, 0) AS dist_price_grp_id, COALESCE(sp.address, '') AS address,
					COALESCE(sp.province_id, '') AS province_id, COALESCE(sp.regency_id, '') AS regency_id,
					COALESCE(sp.sub_district_id, '') AS sub_district_id, COALESCE(sp.ward_id, '') AS ward_id,
					COALESCE(sp.zip_code, '') AS zip_code, COALESCE(sp.ot_loc_id, 0) AS ot_loc_id, COALESCE(sp.latitude, '') AS latitude, COALESCE(sp.is_active, false) AS is_active,
					COALESCE(sp.longitude, '') AS longitude, 
					COALESCE(u.user_fullname, '') AS updated_by_name `
	qWhere := ` WHERE sp.is_del = false `
	qWhere += buildDistributorScopeCondition(dataFilter, custId)
	qWhere += buildDistributorActiveCondition(dataFilter.IsActive)

	if len(dataFilter.AreaID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.AreaID, ",")
		qWhere += ` AND sp.area_id IN (` + intArrStr + `) `
	}

	qFrom := ` 	FROM mst.m_distributor sp
				LEFT JOIN sys.m_user u ON u.user_id = sp.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		return mdistributor, 0, 0, err
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
		sortBy := `distributor_id`
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

	err = repository.Select(&mdistributor, querySelect)
	if err != nil {
		return mdistributor, total, lastPage, err
	}

	return mdistributor, total, lastPage, nil
}

func (repo *distributorRepositoryImpl) FindOneByDistributorIdAndCustId(params entity.DetailDistributorParams) (model.DistributorList, error) {
	distributor := model.DistributorList{}
	args := []interface{}{params.DistributorId}
	scopeCondition := `AND mdist.cust_id LIKE $2`
	args = append(args, params.CustId+"%")
	if params.JwtDistributorId > 0 {
		scopeCondition = `AND (mdist.distributor_id = $2 OR mdist.cust_id LIKE $3)`
		args = []interface{}{params.DistributorId, params.JwtDistributorId, params.CustId + "%"}
	}

	query := `SELECT 
		mdist.cust_id,
		COALESCE(mdist.parent_cust_id, '') AS parent_cust_id,
		mdist.distributor_id,
		mdist.distributor_code,
		mdist.distributor_name,
		mdist.barcode,
		COALESCE(mdist.region_id, 0) AS region_id,
		COALESCE(mdist.area_id, 0) AS area_id,
		COALESCE(mdist.channel_id, 0) AS channel_id,
		COALESCE(mdist.sub_distributor_group_id, 0) AS sub_distributor_group_id,
		COALESCE(mdist.dist_price_grp_id, 0) AS dist_price_grp_id,
		COALESCE(mdist.address, '') AS address,
		COALESCE(mdist.province_id, '') AS province_id,
		COALESCE(mdist.regency_id, '') AS regency_id,
		COALESCE(mdist.sub_district_id, '') AS sub_district_id,
		COALESCE(mdist.ward_id, '') AS ward_id,
		COALESCE(mdist.zip_code, '') AS zip_code,
		COALESCE(mdist.ot_loc_id, 0) AS ot_loc_id,
		COALESCE(mdist.latitude, '') AS latitude,
		COALESCE(mdist.longitude, '') AS longitude,
		mdist.phone,
		mdist.fax_number,
		COALESCE(ar.area_code, '') AS area_code, COALESCE(ar.area_name, '') AS area_name,
		COALESCE(pv.province_id, '') AS province_code, COALESCE(pv.province, '') AS province_name,
		COALESCE(rgc.regency_id, '') AS regency_code, COALESCE(rgc.regency, '') AS regency_name,
		COALESCE(sd.sub_district_id, '') AS sub_district_code, COALESCE(sd.sub_district, '') AS sub_district_name,
		COALESCE(wr.ward_id, '') AS ward_code, COALESCE(wr.ward, '') AS ward_name,
		COALESCE(u.user_fullname, '') AS updated_by_name,
		COALESCE(mdist.is_active, false) AS is_active,
		mdist.is_del,
		mdist.created_by,
		mdist.created_at,
		mdist.updated_by,
		mdist.updated_at,
		mdist.deleted_by,
		mdist.deleted_at,
		COALESCE(mr.region_code, '') AS region_code, COALESCE(mr.region_name, '') AS region_name,
        COALESCE(mc.channel_code, '') AS channel_code, COALESCE(mc.channel_name, '') AS channel_name,
        COALESCE(msdg.sub_distributor_group_code, '') AS sub_distributor_group_code, COALESCE(msdg.sub_distributor_group_name, '') AS sub_distributor_group_name,
		COALESCE(mdpg.dist_price_grp_code, '') AS dist_price_grp_code, COALESCE(mdpg.dist_price_grp_name, '') AS dist_price_grp_name,
		COALESCE(mdist.allow_add_product, false) AS allow_add_product,
		COALESCE(mdist.allow_edit_product, false) AS allow_edit_product,
		COALESCE(mdist.allow_manage_pricing, false) AS allow_manage_pricing,
		COALESCE(mdist.allow_upload_secondary_sales, false) AS allow_upload_secondary_sales
	FROM mst.m_distributor mdist
	LEFT JOIN sys.m_user u ON u.user_id = mdist.updated_by 
	LEFT JOIN mst.m_area ar ON ar.area_id = mdist.area_id AND ar.cust_id = '` + params.ParentCustId + `'
	LEFT JOIN mst.m_province pv ON pv.province_id = mdist.province_id AND pv.cust_id = '` + params.ParentCustId + `'
	LEFT JOIN mst.m_regency rgc ON rgc.regency_id = mdist.regency_id AND rgc.cust_id = '` + params.ParentCustId + `'
	LEFT JOIN mst.m_sub_district sd ON sd.sub_district_id = mdist.sub_district_id AND sd.cust_id = '` + params.ParentCustId + `'
	LEFT JOIN mst.m_ward wr ON wr.ward_id = mdist.ward_id AND wr.cust_id = '` + params.ParentCustId + `'
	left join mst.m_region mr on mr.region_id = mdist.region_id and mr.cust_id = '` + params.ParentCustId + `'
	left join mst.m_channel mc on mc.channel_id = mdist.channel_id and mc.cust_id = '` + params.ParentCustId + `'
	left join mst.m_sub_distributor_group msdg on msdg.sub_distributor_group_id = mdist.sub_distributor_group_id and msdg.cust_id = '` + params.ParentCustId + `'
	left join mst.m_dist_price_group mdpg on mdpg.dist_price_grp_id = mdist.dist_price_grp_id and mdpg.cust_id = '` + params.ParentCustId + `'
	WHERE mdist.distributor_id = $1 
	` + scopeCondition + `
	AND mdist.is_del = false;`
	err := repo.Get(&distributor, query, args...)
	if err != nil {
		return distributor, err
	}

	return distributor, nil
}

func (repo *distributorRepositoryImpl) FindAllowManagePricingByDistributorID(distributorID int64) (bool, error) {
	var allowManagePricing bool
	query := `
		SELECT COALESCE(mdist.allow_manage_pricing, false) AS allow_manage_pricing
		FROM mst.m_distributor mdist
		WHERE mdist.distributor_id = $1
		AND mdist.is_del = false
		LIMIT 1;`
	err := repo.Get(&allowManagePricing, query, distributorID)
	if err != nil {
		return false, err
	}

	return allowManagePricing, nil
}

func (repo *distributorRepositoryImpl) FindAllowUploadSecondarySalesByDistributorID(distributorID int64) (bool, error) {
	var allowed bool
	query := `
		SELECT COALESCE(mdist.allow_upload_secondary_sales, false) AS allow_upload_secondary_sales
		FROM mst.m_distributor mdist
		WHERE mdist.distributor_id = $1
		AND mdist.is_del = false
		LIMIT 1;`
	err := repo.Get(&allowed, query, distributorID)
	if err != nil {
		return false, err
	}
	return allowed, nil
}

func (repo *distributorRepositoryImpl) FindDistributorIdByCustId(custId string) (int64, error) {
	var distributorID int64
	query := `
		SELECT COALESCE(mc.distributor_id, 0) AS distributor_id
		FROM smc.m_customer mc
		WHERE mc.cust_id = $1
		LIMIT 1;`
	err := repo.Get(&distributorID, query, custId)
	if err != nil {
		return 0, err
	}
	return distributorID, nil
}

func (repo *distributorRepositoryImpl) FindOneByDistributorCodeAndCustId(distributorCode string, custId string) (model.DistributorList, error) {
	distributor := model.DistributorList{}
	query := `SELECT sp.distributor_id, sp.distributor_code
			FROM mst.m_distributor sp
			WHERE sp.cust_id = $2
			AND sp.distributor_code = $1
			AND sp.is_del = false;`
	err := repo.Get(&distributor, query, distributorCode, custId)
	if err != nil {
		return distributor, err
	}

	return distributor, nil
}

func (repo *distributorRepositoryImpl) FindOneByParentCustIdAndDistributorCode(custId, distributorCode string) (model.DistributorWithCustomer, error) {
	detail := model.DistributorWithCustomer{}
	query := `
		SELECT
			md.parent_cust_id,
			md.distributor_id,
			md.distributor_code,
			md.distributor_name,
			mc.cust_id AS dist_cust_id,
			md.allow_add_product,
			md.allow_edit_product,
			md.allow_manage_pricing,
			md.allow_upload_secondary_sales
		FROM mst.m_distributor md
		INNER JOIN smc.m_customer mc ON mc.distributor_id = md.distributor_id
		WHERE
			md.parent_cust_id = $1
			AND md.distributor_code = $2
			AND md.is_del = false
		LIMIT 1;`
	err := repo.Get(&detail, query, custId, distributorCode)
	if err != nil {
		log.Error(err.Error())
		return detail, err
	}

	return detail, nil
}

func (repo *distributorRepositoryImpl) FindAllDistributorContactByDistIdAndCustId(params entity.DetailDistributorParams) ([]model.DistributorContact, error) {
	distributorContact := []model.DistributorContact{}
	query := `SELECT *
			  FROM mst.m_distributor_contact
			  WHERE distributor_id = $1
			  AND cust_id LIKE $2 `
	err := repo.Select(&distributorContact, query, params.DistributorId, params.CustId+"%")
	if err != nil {
		return distributorContact, err
	}

	return distributorContact, nil
}

func (repo *distributorRepositoryImpl) FindAllDistributorTaxByDistIdAndCustId(params entity.DetailDistributorParams) ([]model.DistributorTax, error) {
	distributorTax := []model.DistributorTax{}
	query := `SELECT *
			  FROM mst.m_distributor_tax
			  WHERE distributor_id = $1
			  AND cust_id LIKE $2 `
	err := repo.Select(&distributorTax, query, params.DistributorId, params.CustId+"%")
	if err != nil {
		return distributorTax, err
	}

	return distributorTax, nil
}

func (repository *distributorTransaction) Update(distributorID int64, request entity.UpdateDistributorRequest) error {
	var (
		r            model.DistributorUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	r.BarcodeProvided = request.BarcodeProvided
	r.ProvinceIdProvided = request.ProvinceIdProvided
	r.RegencyIdProvided = request.RegencyIdProvided
	r.SubDistrictIdProvided = request.SubDistrictIdProvided
	r.WardIdProvided = request.WardIdProvided
	r.ZipCodeProvided = request.ZipCodeProvided
	r.OtLocIdProvided = request.OtLocIdProvided
	r.PhoneProvided = request.PhoneProvided
	r.FaxNumberProvided = request.FaxNumberProvided
	sqlPatch := sql_helper.SQLPatches(r)
	hasField := func(field string) bool {
		for _, existingField := range sqlPatch.Fields {
			if existingField == field {
				return true
			}
		}
		return false
	}

	if request.BarcodeProvided {
		sqlPatch.Args["barcode"] = nil
		if request.Barcode != nil && *request.Barcode != "" {
			sqlPatch.Args["barcode"] = *request.Barcode
		}
		if !hasField("barcode = :barcode") {
			sqlPatch.Fields = append(sqlPatch.Fields, "barcode = :barcode")
		}
	}

	if request.ZipCodeProvided {
		sqlPatch.Args["zip_code"] = nil
		if request.ZipCode != nil && *request.ZipCode != "" {
			sqlPatch.Args["zip_code"] = *request.ZipCode
		}
		if !hasField("zip_code = :zip_code") {
			sqlPatch.Fields = append(sqlPatch.Fields, "zip_code = :zip_code")
		}
	}

	if request.ProvinceIdProvided {
		sqlPatch.Args["province_id"] = nil
		if request.ProvinceId != nil && *request.ProvinceId != "" {
			sqlPatch.Args["province_id"] = *request.ProvinceId
		}
		if !hasField("province_id = :province_id") {
			sqlPatch.Fields = append(sqlPatch.Fields, "province_id = :province_id")
		}
	}

	if request.RegencyIdProvided {
		sqlPatch.Args["regency_id"] = nil
		if request.RegencyId != nil && *request.RegencyId != "" {
			sqlPatch.Args["regency_id"] = *request.RegencyId
		}
		if !hasField("regency_id = :regency_id") {
			sqlPatch.Fields = append(sqlPatch.Fields, "regency_id = :regency_id")
		}
	}

	if request.SubDistrictIdProvided {
		sqlPatch.Args["sub_district_id"] = nil
		if request.SubDistrictId != nil && *request.SubDistrictId != "" {
			sqlPatch.Args["sub_district_id"] = *request.SubDistrictId
		}
		if !hasField("sub_district_id = :sub_district_id") {
			sqlPatch.Fields = append(sqlPatch.Fields, "sub_district_id = :sub_district_id")
		}
	}

	if request.WardIdProvided {
		sqlPatch.Args["ward_id"] = nil
		if request.WardId != nil && *request.WardId != "" {
			sqlPatch.Args["ward_id"] = *request.WardId
		}
		if !hasField("ward_id = :ward_id") {
			sqlPatch.Fields = append(sqlPatch.Fields, "ward_id = :ward_id")
		}
	}

	if request.OtLocIdProvided {
		sqlPatch.Args["ot_loc_id"] = nil
		if request.OtLocId != nil && *request.OtLocId != 0 {
			sqlPatch.Args["ot_loc_id"] = *request.OtLocId
		}
		if !hasField("ot_loc_id = :ot_loc_id") {
			sqlPatch.Fields = append(sqlPatch.Fields, "ot_loc_id = :ot_loc_id")
		}
	}

	if request.PhoneProvided {
		sqlPatch.Args["phone"] = nil
		if request.Phone != nil && *request.Phone != "" {
			sqlPatch.Args["phone"] = *request.Phone
		}
		if !hasField("phone = :phone") {
			sqlPatch.Fields = append(sqlPatch.Fields, "phone = :phone")
		}
	}

	if request.FaxNumberProvided {
		sqlPatch.Args["fax_number"] = nil
		if request.FaxNumber != nil && *request.FaxNumber != "" {
			sqlPatch.Args["fax_number"] = *request.FaxNumber
		}
		if !hasField("fax_number = :fax_number") {
			sqlPatch.Fields = append(sqlPatch.Fields, "fax_number = :fax_number")
		}
	}

	for i := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_distributor
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND distributor_id = :distributor_id;`

	sqlPatch.Args["distributor_id"] = distributorID
	sqlPatch.Args["cust_id"] = request.CustId
	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return constant.ErrNoRowsAffected
	}
	if nRows == 0 {
		return constant.ErrNoRowsAffected
	}

	return nil
}

func (repository *distributorTransaction) UpdateDistributorContact(distributorID, distributorContactID int64, request entity.DistributorContactUpdate) error {
	var (
		r            model.DistributorContactUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)

	sqlPatch := sql_helper.SQLPatches(r)

	if request.Email != nil {
		sqlPatch.Args["email"] = *request.Email
	}

	for i := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_distributor_contact
			  SET ` + sqlSetFields + `
			  WHERE  distributor_contact_id = :distributor_contact_id
			  AND distributor_id=:distributor_id;`

	sqlPatch.Args["distributor_id"] = distributorID
	sqlPatch.Args["distributor_contact_id"] = distributorContactID

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return constant.ErrNoRowsAffected
	}
	if nRows == 0 {
		return constant.ErrNoRowsAffected
	}

	return nil
}

func (repository *distributorTransaction) UpdateDistributorTax(distributorID, distributorTaxID int64, request entity.DistributorTaxUpdate) error {
	var (
		r            model.DistributorTaxUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)

	sqlPatch := sql_helper.SQLPatches(r)

	for i := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_distributor_tax
			  SET ` + sqlSetFields + `
			  WHERE  distributor_tax_id = :distributor_tax_id
			  AND distributor_id=:distributor_id;`

	sqlPatch.Args["distributor_id"] = distributorID
	sqlPatch.Args["distributor_tax_id"] = distributorTaxID

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return constant.ErrNoRowsAffected
	}
	if nRows == 0 {
		return constant.ErrNoRowsAffected
	}

	return nil
}

func (repository *distributorTransaction) DeleteDistributorContactNotIn(distributorId int64, distributorContactId []int64) error {
	query, args, err := sqlx.In("DELETE FROM mst.m_distributor_contact	WHERE distributor_id = ? AND distributor_contact_id not in (?);", distributorId, distributorContactId)
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

func (repository *distributorTransaction) DeleteAllDistributorContact(distributorId int64) error {
	_, err := repository.tx.Exec("DELETE FROM mst.m_distributor_contact WHERE distributor_id = $1;", distributorId)
	if err != nil {
		return err
	}

	return nil
}

func (repository *distributorTransaction) DeleteDistributorTaxNotIn(distributorId int64, distributorTaxId []int64) error {
	query, args, err := sqlx.In("DELETE FROM mst.m_distributor_tax	WHERE distributor_id = ? AND distributor_tax_id not in (?);", distributorId, distributorTaxId)
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

func (repository *distributorTransaction) DeleteAllDistributorTax(distributorId int64) error {
	_, err := repository.tx.Exec("DELETE FROM mst.m_distributor_tax WHERE distributor_id = $1;", distributorId)
	if err != nil {
		return err
	}

	return nil
}

func (repository *distributorRepositoryImpl) Delete(custId string, distributorId int64, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_distributor
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND distributor_id = :distributor_id;`

	wMap := map[string]interface{}{
		"cust_id":        custId,
		"distributor_id": distributorId,
		"deleted_by":     deletedBy,
	}
	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		return err
	}
	if nRows, err = result.RowsAffected(); err != nil {
		return constant.ErrNoRowsAffected
	}
	if nRows == 0 {
		return constant.ErrNoRowsAffected
	}

	return nil
}

func (repository *distributorRepositoryImpl) FindAllByCustIdWithCustomer(dataFilter entity.DistributorQueryFilter, custId string) ([]model.DistributorList, int, int, error) {
	mdistributor := []model.DistributorList{}

	selectCount := ` COUNT(DISTINCT sp.distributor_id) AS total `
	selectField := ` DISTINCT ON (sp.distributor_id) sp.cust_id, COALESCE(sp.parent_cust_id, '') AS parent_cust_id, sp.distributor_id, sp.distributor_code, distributor_name, sp.barcode,
					COALESCE(sp.region_id, 0) AS region_id, COALESCE(sp.area_id, 0) AS area_id,
					COALESCE(sp.channel_id, 0) AS channel_id, COALESCE(sp.sub_distributor_group_id, 0) AS sub_distributor_group_id,
					COALESCE(sp.dist_price_grp_id, 0) AS dist_price_grp_id, COALESCE(sp.address, '') AS address,
					COALESCE(sp.province_id, '') AS province_id, COALESCE(sp.regency_id, '') AS regency_id,
					COALESCE(sp.sub_district_id, '') AS sub_district_id, COALESCE(sp.ward_id, '') AS ward_id,
					COALESCE(sp.zip_code, '') AS zip_code, COALESCE(sp.ot_loc_id, 0) AS ot_loc_id, COALESCE(sp.latitude, '') AS latitude,
					COALESCE(ar.area_code, '') AS area_code, COALESCE(ar.area_name, '') AS area_name,
					COALESCE(pv.province_id, '') AS province_code, COALESCE(pv.province, '') AS province_name,
					COALESCE(rgc.regency_id, '') AS regency_code, COALESCE(rgc.regency, '') AS regency_name,
					COALESCE(sd.sub_district_id, '') AS sub_district_code, COALESCE(sd.sub_district, '') AS sub_district_name,
					COALESCE(wr.ward_id, '') AS ward_code, COALESCE(wr.ward, '') AS ward_name,
					COALESCE(sp.longitude, '') AS longitude, COALESCE(rgn.region_code, '') AS region_code, COALESCE(rgn.region_name, '') AS region_name, COALESCE(cnl.channel_code, '') AS channel_code, COALESCE(cnl.channel_name, '') AS channel_name, COALESCE(dpg.dist_price_grp_code, '') AS dist_price_grp_code, COALESCE(dpg.dist_price_grp_name, '') AS dist_price_grp_name, COALESCE(sp.is_active, false) AS is_active,
					COALESCE(u.user_fullname, '') AS updated_by_name, COALESCE(mc.cust_id, '') AS customer_id `

	qWhere := ` WHERE sp.is_del = false 
				AND sp.cust_id LIKE '%` + custId + `%` + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (CAST(sp.distributor_id AS TEXT) ILIKE '%` + dataFilter.Query + `%' 
								OR sp.distributor_code ILIKE '%` + dataFilter.Query + `%'
								OR sp.distributor_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if len(dataFilter.AreaID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.AreaID, ",")
		qWhere += ` AND sp.area_id IN (` + intArrStr + `) `
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND sp.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND sp.is_active = false `
		}
	}
	qFrom := ` 	FROM mst.m_distributor sp
				LEFT JOIN mst.m_dist_price_group dpg ON dpg.dist_price_grp_id = sp.dist_price_grp_id
				LEFT JOIN mst.m_channel cnl ON cnl.channel_id = sp.channel_id
				LEFT JOIN mst.m_region rgn ON rgn.region_id = sp.region_id
				
				LEFT JOIN mst.m_area ar ON ar.area_id = sp.area_id
				LEFT JOIN mst.m_province pv ON pv.province_id = sp.province_id
				LEFT JOIN mst.m_regency rgc ON rgc.regency_id = sp.regency_id
				LEFT JOIN mst.m_sub_district sd ON sd.sub_district_id = sp.sub_district_id
				LEFT JOIN mst.m_ward wr ON wr.ward_id = sp.ward_id
				LEFT JOIN smc.m_customer mc ON mc.distributor_id = sp.distributor_id

				LEFT JOIN sys.m_user u ON u.user_id = sp.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		return mdistributor, 0, 0, err
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
		sortBy := `distributor_id`
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

	err = repository.Select(&mdistributor, querySelect)
	if err != nil {
		return mdistributor, total, lastPage, err
	}

	return mdistributor, total, lastPage, nil
}
