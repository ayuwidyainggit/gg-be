package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"master/pkg/str"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"

	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	IsExists(relationId int, custId, foreignKey string) (bool, error)
	IsKeyExists(relationId string, custId, foreignKey string) (bool, error)
	FindOne(params entity.DetailProductParams) (model.Product, error)
	FindOneByProductIdAndCustId(params entity.DetailProductParams) (model.Product, error)
	FindOneByProductIdAndDistributorId(params entity.DetailProductParams) (model.Product, error)
	FindOneByProductCodeAndCustId(productCode string, custId string) (model.Product, error)
	FindAll(dataFilter entity.ProductQueryFilter, custId string) ([]model.Product, int, int, error)
	FindAllByCustId(dataFilter entity.ProductQueryFilter, custId string) (product []model.Product, total int, lastPage int, err error)
	FindAllByCustIdLookup(dataFilter entity.ProductQueryFilter, custId string) (product []model.Product, total int, lastPage int, err error)
	FindAllByDistributorLookup(dataFilter entity.ProductQueryFilter, custId string) (product []model.Product, total int, lastPage int, err error)
	ReportList(filter entity.ProductReportQueryFilter) (data []entity.ProductReportResponse, total int, lastPage int, err error)
	FindAllByCustIdLookupDistPrice(dataFilter entity.ProductQueryFilter) (product []model.ProductDistPrice, total int, lastPage int, err error)
	FindAllByDistributorLookupDistPrice(dataFilter entity.ProductQueryFilter) (product []model.ProductDistPrice, total int, lastPage int, err error)
	FindAllByCustIdSearch(dataFilter entity.ProductQueryFilter, custId string) (product []model.Product, total int, lastPage int, err error)
	FindAllExport(dataFilter entity.ProductQueryFilter, custId string) ([]model.Product, int, error)
	FindByCodeProductCategory(custId, parentCustId, code, name string) (int64, error)
	FindByCodeProductLine(custId, parentCustId, code, name string) (int64, error)
	FindByCodeBrand(custId string, parentCustId string, plId int64, code, name string) (int64, error)
	FindByCodeSubBrand1(custId string, parentCustId string, brandId int64, code, name string) (int64, error)
	FindByCodeSubBrand2(custId, parentCustId, code, name string) (int64, error)
	FindByCodeFlavor(custId, parentCustId, code, name string) (int64, error)
	FindByCodePackSize(custId, parentCustId, code, name string) (int64, error)
	FindByCodePackType(custId, parentCustId, code, name string) (int64, error)
	FindByCodeSupplier(custId, parentCustId, code, name string) (int64, error)
	FindByCodePrincipal(custId, parentCustId, code, name string) (int64, error)
	FindByCodeCPro(custId, parentCustId, code, name string) (int64, error)
	FindByCodeUnit(custId, parentCustId, code, name string) (string, error)
	GetOrCreateProductCoretax(custId, parentCustId, code, name string, createdBy int64) (string, error)
	CreateProduct(data entity.ProcessedProductRow) error
	GetProductImportInstructions() ([]model.ImportInstruction, error)
	GetProductCategoryIdByCode(custId, code string) (int64, error)
	GetProductLineIdByCode(custId, code string) (int64, error)
	GetBrandIdByCode(custId, code string) (int64, error)
	GetSubBrand1IdByCode(custId, code string) (int64, error)
	GetSubBrand2IdByCode(custId, code string) (int64, error)
	GetFlavorIdByCode(custId, code string) (int64, error)
	GetPackSizeIdByCode(custId, code string) (int64, error)
	GetPackTypeIdByCode(custId, code string) (int64, error)
	GetSupplierIdByCode(custId, code string) (int64, error)
	GetPrincipalIdByCode(custId, code string) (int64, error)
	GetCProIdByCode(custId, code string) (int64, error)
	GetUnitIdByCode(custId, code string) (string, error)
	GetUnitProductCoretaxIdByCode(custId, parentCustId, code string) (string, error)
	CreateImportHistory(uploadType, fileName, custId string, uploadedBy int64, totalData int) (int64, error)
	UpdateImportHistory(historyId int64, success, failed int, statusReupload bool) error
	InsertProductTemp(historyId int64, statusInsert, CustId string, data entity.ImportProductRow) error
	CheckBrandExists(custId string, brandId int64) (bool, error)
	CheckBrandCodeDuplicate(custId string, brandId int64, code string) (bool, error)
	CheckProductCatExists(custId string, pcatId int64) (bool, error)
	CheckProductCatCodeDuplicate(custId string, pcatId int64, code string) (bool, error)
	CheckProductLineExists(custId string, plId int64) (bool, error)
	CheckProductLineCodeDuplicate(custId string, plId int64, code string) (bool, error)
	CheckSbrand1Exists(custId string, id int64) (bool, error)
	CheckSbrand1CodeDuplicate(custId string, id int64, code string) (bool, error)
	CheckSbrand2Exists(custId string, id int64) (bool, error)
	CheckSbrand2CodeDuplicate(custId string, id int64, code string) (bool, error)
	CheckFlavorExists(custId string, id int64) (bool, error)
	CheckFlavorCodeDuplicate(custId string, id int64, code string) (bool, error)
	CheckPackTypeExists(custId string, id int64) (bool, error)
	CheckPackTypeCodeDuplicate(custId string, id int64, code string) (bool, error)
	CheckPackSizeExists(custId string, id int64) (bool, error)
	CheckPackSizeCodeDuplicate(custId string, id int64, code string) (bool, error)
	CheckSupplierExists(custId string, id int64) (bool, error)
	CheckSupplierCodeDuplicate(custId string, id int64, code string) (bool, error)
	CheckPrincipalExists(custId string, id int64) (bool, error)
	CheckPrincipalCodeDuplicate(custId string, id int64, code string) (bool, error)
	CheckCProExists(custId string, id int64) (bool, error)
	CheckCProCodeDuplicate(custId string, id int64, code string) (bool, error)
	CheckCoretaxExists(custId string, code string) (bool, error)
	CheckUnitExists(custId string, code string) (bool, error)
	CheckProductExists(custId, parentCustId, proCode string) (bool, error)
	UpdateImportBrand(custId string, brandId int64, code, name string) error
	UpdateImportCategory(custId string, pcatId int64, code, name string) error
	UpdateImportProductLine(custId string, id int64, code, name string) error
	UpdateImportSubBrand1(custId string, sbrand1Id int64, code, name string) error
	UpdateImportSubBrand2(custId string, sbrand2Id int64, code, name string) error
	UpdateImportFlavor(custId string, flavorId int64, code, name string) error
	UpdateImportPackType(custId string, ptypeId int64, code, name string) error
	UpdateImportPackSize(custId string, psizeId int64, code, name string) error
	UpdateImportSupplier(custId string, supId int64, code, name string) error
	UpdateImportPrincipal(custId string, principalId int64, code, name string) error
	UpdateImportCPro(custId string, cproId int64, code, name string) error
	UpdateImportProductCoretax(custId string, code, name string) error
	UpdateImportUnit(custId string, id, name string) error
	FindProductIdByCode(custId, proCode string) (int, error)
	FindProductById(custId string, productId int) (*entity.ProcessedProductRow, error)
	UpdateImportProduct(p entity.ProcessedProductRow) error
	InsertProductUpdateTemp(temp entity.ImportProductUpdateTemp) error
	DeleteProductUpdateTempByCtid(ctid string) error
	DeleteProductTempByCtid(ctid string) error
	GetCtidsProductUpdateByHistory(historyId int64) ([]string, error)
	GetCtidsProductInsertByHistory(historyId int64) ([]string, error)
	CountProductUpdateTemp(historyId int64) (int, error)
	CountProductTemp(historyId int64) (int, error)
	GetImportProductTotalData(historyId int64) (int, error)
	Store(ctx context.Context, product model.Product) (int64, error)
	BulkStore(products Products) ([]int64, error)
	Update(productId int64, request entity.UpdateProductRequest) error
	Delete(custId string, productId int64, deletedBy int64) error
	DeleteWithContext(ctx context.Context, custId string, productId int64, deletedBy int64) error
	DeleteMultiple(custId string, productId []int64, deletedBy int64) error
	FindDistributor(custId string) ([]model.MCustomer, error)
	StoreDist(productsDist []model.ProductDistCreate) error
	FindAllPrincipal(dataFilter entity.ProductPrincipalQueryFilter, custId string) (product []model.Principal, total int, lastPage int, err error)
	FindAllCategory(dataFilter entity.ProductCategoryQueryFilter, custId string) (product []model.ProductCat, total int, lastPage int, err error)
	FindAllBrand(dataFilter entity.ProductBrandQueryFilter, custId string) (product []model.Brand, total int, lastPage int, err error)
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepositoryImpl{db}
}

const reportTypeExpr = `CASE
	WHEN LENGTH(mp.cust_id) = 6 THEN 'Own Products'
	WHEN COALESCE(mp.parent_pro_id, 0) = 0 THEN 'Own Products'
	WHEN mp.is_product_mapping = false THEN 'Product Assigned'
	ELSE 'Product Mapping'
END`

const reportNormalizeExpr = `mp.parent_pro_id <> 0 AND mp.is_product_mapping = true AND md.allow_upload_secondary_sales = true AND parent.pro_id IS NOT NULL`

func (r *productRepositoryImpl) ReportList(filter entity.ProductReportQueryFilter) (data []entity.ProductReportResponse, total int, lastPage int, err error) {
	page, limit := filter.Page, filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	from := ` FROM mst.m_product mp
		LEFT JOIN (
			SELECT cust_id, BOOL_OR(COALESCE(allow_upload_secondary_sales, false)) AS allow_upload_secondary_sales
			FROM mst.m_distributor
			GROUP BY cust_id
		) md ON md.cust_id = mp.cust_id
		LEFT JOIN mst.m_product parent ON parent.pro_id = mp.parent_pro_id
			AND parent.cust_id = LEFT(mp.cust_id, 6)
			AND parent.is_del = false AND parent.is_active = true`
	where := ` WHERE mp.is_del = false AND mp.is_active = true AND mp.cust_id IN (?)`
	args := []interface{}{filter.CustIDs}
	if filter.Query != "" {
		where += ` AND (mp.pro_name ILIKE ? OR mp.pro_code ILIKE ?)`
		wildcard := "%" + filter.Query + "%"
		args = append(args, wildcard, wildcard)
	}

	countQuery, countArgs, err := sqlx.In(`SELECT COUNT(*)`+from+where, args...)
	if err != nil {
		return nil, 0, 0, err
	}
	countQuery = r.Rebind(countQuery)
	if err = r.Get(&total, countQuery, countArgs...); err != nil {
		return nil, 0, 0, err
	}

	sortColumns := map[string]string{
		"pro_name": "pro_name",
		"pro_code": "pro_code",
		"type":     "type",
		"pro_id":   "pro_id",
	}
	sortColumn := sortColumns[filter.SortBy]
	if sortColumn == "" {
		sortColumn = "pro_name"
	}
	sortOrder := strings.ToUpper(filter.SortOrder)
	if sortOrder != "DESC" {
		sortOrder = "ASC"
	}

	normalized := reportNormalizeExpr
	selectQuery := `SELECT
		CASE WHEN ` + normalized + ` THEN parent.cust_id ELSE mp.cust_id END AS cust_id,
		CASE WHEN ` + normalized + ` THEN parent.pro_id ELSE mp.pro_id END AS pro_id,
		CASE WHEN ` + normalized + ` THEN parent.pro_code ELSE mp.pro_code END AS pro_code,
		CASE WHEN ` + normalized + ` THEN parent.pro_name ELSE mp.pro_name END AS pro_name,
		CASE WHEN ` + normalized + ` THEN mp.cust_id ELSE NULL END AS original_cust_id,
		CASE WHEN ` + normalized + ` THEN mp.pro_id ELSE NULL END AS original_pro_id,
		CASE WHEN ` + normalized + ` THEN mp.pro_code ELSE NULL END AS original_pro_code,
		CASE WHEN ` + normalized + ` THEN mp.parent_pro_id ELSE NULL END AS original_parent_pro_id,
		` + reportTypeExpr + ` AS type` + from + where + ` ORDER BY ` + sortColumn + ` ` + sortOrder + ` LIMIT ? OFFSET ?`
	dataArgs := append([]interface{}{}, args...)
	dataArgs = append(dataArgs, limit, offset)
	selectQuery, dataArgs, err = sqlx.In(selectQuery, dataArgs...)
	if err != nil {
		return nil, total, 0, err
	}
	selectQuery = r.Rebind(selectQuery)
	if err = r.Select(&data, selectQuery, dataArgs...); err != nil {
		return nil, total, 0, err
	}
	return data, total, sql_helper.CalculateLastPage(total, limit), nil
}

type productRepositoryImpl struct {
	*sqlx.DB
}

type Products []model.Product

func (repository *productRepositoryImpl) resolveDistPriceGroupID(parentCustID string, distributorID int64) (int, error) {
	if distributorID == 0 {
		return 0, nil
	}

	var distPriceGroupID sql.NullInt64
	query := `SELECT dist_price_grp_id
			FROM mst.m_distributor
			WHERE cust_id = $1
			AND distributor_id = $2`

	err := repository.Get(&distPriceGroupID, query, parentCustID, distributorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		log.Error("productRepository, resolveDistPriceGroupID, err:", err.Error())
		return 0, err
	}

	if !distPriceGroupID.Valid {
		return 0, nil
	}

	return int(distPriceGroupID.Int64), nil
}

func (repository *productRepositoryImpl) IsExists(relationId int, custId, foreignKey string) (bool, error) {
	isExists := model.IsExists{}
	query := `SELECT EXISTS(SELECT pro_id 
							FROM mst.m_product 
							WHERE is_del = false 
							AND ` + foreignKey + ` = $1 
							AND cust_id = $2);`
	err := repository.Get(&isExists, query, relationId, custId)
	if err != nil {
		log.Error("MProductRepository, IsExists, err:", err.Error())
		return false, err
	}

	return isExists.Exists, err
}

func (repository *productRepositoryImpl) IsKeyExists(relationId string, custId, foreignKey string) (bool, error) {
	isExists := model.IsExists{}
	query := `SELECT EXISTS(SELECT pro_id 
							FROM mst.m_product 
							WHERE is_del = false 
							AND ` + foreignKey + ` = $1 
							AND cust_id = $2);`
	err := repository.Get(&isExists, query, relationId, custId)
	if err != nil {
		log.Error("MProductRepository, IsExists, err:", err.Error())
		return false, err
	}

	return isExists.Exists, err
}

func (repository *productRepositoryImpl) FindOne(params entity.DetailProductParams) (model.Product, error) {
	product := model.Product{}
	query := `SELECT 
				p.cust_id, p.pro_id, p.pro_code, p.bar_code, p.pro_name,
				p.vat, p.vat_bg, p.vat_lg_purch, p.vat_lg_sell, p.cogs, p.pro_status,
				p.pcat_id, pc.pcat_code, pc.pcat_name,
				br.brand_id, br.brand_code, br.brand_name, 
				br.pl_id, pl.pl_code, pl.pl_name,
				p.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
				p.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
				p.flavor_id, fv.flavor_code, fv.flavor_name, 
				p.ptype_id, pt.ptype_code, pt.ptype_name,
				p.psize_id, ps.psize_code, ps.psize_name,
				p.sup_id, su.sup_code, su.sup_name,
				p.principal_id, pr.principal_code, pr.principal_name,
				p.c_pro_id, cp.c_pro_code, cp.c_pro_name,
				p.is_main_pro, p.sort_no, p.item_no, p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5, 
				COALESCE(un1.unit_name, p.unit_id1, '') AS unit_name1,
				COALESCE(un2.unit_name, p.unit_id2, '') AS unit_name2,
				COALESCE(un3.unit_name, p.unit_id3, '') AS unit_name3,
				NULLIF(unc1.unit_id_coretax, '') AS unit_id_coretax1, NULLIF(unc2.unit_id_coretax, '') AS unit_id_coretax2, NULLIF(unc3.unit_id_coretax, '') AS unit_id_coretax3,
				NULLIF(unc1.unit_name_coretax, '') AS unit_name_coretax1, NULLIF(unc2.unit_name_coretax, '') AS unit_name_coretax2, NULLIF(unc3.unit_name_coretax, '') AS unit_name_coretax3,
				p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5, 
				p.is_batch, p.is_exp_date, 
				p.weight,p.length, p.width, p.height, p.volume,
				p.purch_price1, p.sell_price1,
				p.purch_price2, p.sell_price2,
				p.purch_price3, p.sell_price3,
				p.purch_price4, p.sell_price4,
				p.purch_price5, p.sell_price5,
				p.weight1, p.length1, p.width1, p.height1, p.volume1, 
				p.weight2, p.length2, p.width2, p.height2, p.volume2, 
				p.weight3, p.length3, p.width3, p.height3, p.volume3, 
				p.weight4, p.length4, p.width4, p.height4, p.volume4, 
				p.weight5, p.length5, p.width5, p.height5, p.volume5,  
				p.parent_pro_id, 
				p.excise_rate, p.excise_tax, 
				p.is_active, p.created_by, p.created_at, 
				p.updated_by, p.updated_at, p.is_del, p.deleted_by, p.deleted_at, p.image_url,
				p.saf_stock_unit_id, p.saf_stock_qty, p.min_stock_unit_id, p.min_stock_qty,
				saf_unit.unit_name AS saf_stock_unit_name,
				min_unit.unit_name AS min_stock_unit_name,
				(SELECT uc.user_fullname FROM sys.m_user uc WHERE uc.user_id = p.created_by LIMIT 1) AS created_by_name,
				(SELECT u.user_fullname FROM sys.m_user u WHERE u.user_id = p.updated_by LIMIT 1) AS updated_by_name,
				p.distributor_id AS distributor_id,
				(SELECT d.distributor_name FROM mst.m_distributor d WHERE d.distributor_id = p.distributor_id LIMIT 1) AS distributor_name,
				parent.pro_code AS parent_pro_code,
				parent.pro_name AS parent_pro_name,
				p.pro_code_coretax,
				(SELECT pct2.pro_name_coretax FROM mst.m_product_coretax pct2 WHERE pct2.pro_code_coretax = p.pro_code_coretax LIMIT 1) AS pro_name_coretax
			FROM mst.m_product p
			LEFT JOIN mst.m_unit saf_unit ON saf_unit.unit_id = p.saf_stock_unit_id
			LEFT JOIN mst.m_unit min_unit ON min_unit.unit_id = p.min_stock_unit_id
			LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id 
			LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id
			LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id
			LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id
			LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id
			LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id 
			LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id
			LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id
			LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id 
			LEFT JOIN mst.m_unit un1 ON un1.unit_id = p.unit_id1
			LEFT JOIN mst.m_unit un2 ON un2.unit_id = p.unit_id2
			LEFT JOIN mst.m_unit un3 ON un3.unit_id = p.unit_id3
			LEFT JOIN mst.m_unit_coretax unc1 ON un1.unit_id_coretax = unc1.unit_id_coretax
			LEFT JOIN mst.m_unit_coretax unc2 ON un2.unit_id_coretax = unc2.unit_id_coretax
			LEFT JOIN mst.m_unit_coretax unc3 ON un3.unit_id_coretax = unc3.unit_id_coretax
			LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id
			LEFT JOIN LATERAL (
				SELECT pro_code, pro_name FROM mst.m_product par WHERE par.pro_id = p.parent_pro_id LIMIT 1
			) parent ON true
			LEFT JOIN LATERAL (
				SELECT c_pro_code, c_pro_name FROM mst.m_cons_product cp2 WHERE cp2.c_pro_id = p.c_pro_id LIMIT 1
			) cp ON true
			WHERE p.pro_id = $1`
	err := repository.Get(&product, query, params.ProductId)
	if err != nil {
		log.Error("productRepository, FindOne, err:", err.Error())
		return product, err
	}

	return product, nil
}

func (repository *productRepositoryImpl) FindOneByProductIdAndCustId(params entity.DetailProductParams) (model.Product, error) {
	product := model.Product{}
	query := `SELECT 
				p.cust_id, p.pro_id, p.pro_code, p.bar_code, p.pro_name,
				p.vat, p.vat_bg, p.vat_lg_purch, p.vat_lg_sell, p.cogs, p.pro_status,
				p.pcat_id, pc.pcat_code, pc.pcat_name,
				br.brand_id, br.brand_code, br.brand_name, 
				br.pl_id, pl.pl_code, pl.pl_name,
				p.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
				p.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
				p.flavor_id, fv.flavor_code, fv.flavor_name, 
				p.ptype_id, pt.ptype_code, pt.ptype_name,
				p.psize_id, ps.psize_code, ps.psize_name,
				p.sup_id, su.sup_code, su.sup_name,
				p.principal_id, pr.principal_code, pr.principal_name,
				p.c_pro_id, cp.c_pro_code, cp.c_pro_name,
				p.is_main_pro, p.sort_no, p.item_no, p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5, 
				COALESCE(un1.unit_name, p.unit_id1, '') AS unit_name1,
				COALESCE(un2.unit_name, p.unit_id2, '') AS unit_name2,
				COALESCE(un3.unit_name, p.unit_id3, '') AS unit_name3,
				NULLIF(unc1.unit_id_coretax, '') AS unit_id_coretax1, NULLIF(unc2.unit_id_coretax, '') AS unit_id_coretax2, NULLIF(unc3.unit_id_coretax, '') AS unit_id_coretax3,
				NULLIF(unc1.unit_name_coretax, '') AS unit_name_coretax1, NULLIF(unc2.unit_name_coretax, '') AS unit_name_coretax2, NULLIF(unc3.unit_name_coretax, '') AS unit_name_coretax3,
				p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5, 
				p.is_batch, p.is_exp_date, 
				p.weight,p.length, p.width, p.height, p.volume,
				p.purch_price1, p.sell_price1,
				p.purch_price2, p.sell_price2,
				p.purch_price3, p.sell_price3,
				p.purch_price4, p.sell_price4,
				p.purch_price5, p.sell_price5,
				p.weight1, p.length1, p.width1, p.height1, p.volume1, 
				p.weight2, p.length2, p.width2, p.height2, p.volume2, 
				p.weight3, p.length3, p.width3, p.height3, p.volume3, 
				p.weight4, p.length4, p.width4, p.height4, p.volume4, 
				p.weight5, p.length5, p.width5, p.height5, p.volume5,  
				p.parent_pro_id, 
				p.excise_rate, p.excise_tax, 
				p.is_active, p.created_by, p.created_at, 
				p.updated_by, p.updated_at, p.is_del, p.deleted_by, p.deleted_at, p.image_url,
				p.saf_stock_unit_id, p.saf_stock_qty, p.min_stock_unit_id, p.min_stock_qty,
				saf_unit.unit_name AS saf_stock_unit_name,
				min_unit.unit_name AS min_stock_unit_name,
				uc.user_fullname AS created_by_name,
				u.user_fullname AS updated_by_name,
				dist.distributor_id, dist.distributor_name,
				parent.pro_code AS parent_pro_code, parent.pro_name AS parent_pro_name,p.pro_code_coretax,pct.pro_name_coretax
			FROM mst.m_product p
			LEFT JOIN sys.m_user u ON u.user_id = p.updated_by
			LEFT JOIN sys.m_user uc ON uc.user_id = p.created_by
			LEFT JOIN mst.m_distributor dist ON dist.distributor_id  = p.distributor_id AND dist.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit saf_unit ON saf_unit.unit_id = p.saf_stock_unit_id AND saf_unit.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit min_unit ON min_unit.unit_id = p.min_stock_unit_id AND min_unit.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + params.ParentCustID + `' 
			LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id AND pc.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id AND sb2.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id AND ps.cust_id = '` + params.ParentCustID + `' 
			LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id AND pt.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id AND fv.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id AND pr.cust_id = '` + params.ParentCustID + `' 
			LEFT JOIN mst.m_unit un1 ON un1.unit_id = p.unit_id1 AND un1.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit un2 ON un2.unit_id = p.unit_id2 AND un2.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit un3 ON un3.unit_id = p.unit_id3 AND un3.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit_coretax unc1 ON un1.unit_id_coretax = unc1.unit_id_coretax AND un1.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit_coretax unc2 ON un2.unit_id_coretax = unc2.unit_id_coretax AND un2.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit_coretax unc3 ON un3.unit_id_coretax = unc3.unit_id_coretax AND un3.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id
			LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id
			LEFT JOIN mst.m_product parent ON parent.pro_id = p.parent_pro_id
			LEFT JOIN mst.m_product_coretax pct ON pct.pro_code_coretax = p.pro_code_coretax
			WHERE p.pro_id = $1 
			AND p.cust_id = $2`
	err := repository.Get(&product, query, params.ProductId, params.ParentCustID)
	if err != nil {
		log.Error("productRepository, FindOneByProductTypeAndCustId, err:", err.Error())
		return product, err
	}

	return product, nil
}

func (repository *productRepositoryImpl) FindOneByProductIdAndDistributorId(params entity.DetailProductParams) (model.Product, error) {
	product := model.Product{}

	asiaJkt, _ := time.LoadLocation("Asia/Jakarta")
	timeNowJkt := time.Now().In(asiaJkt)
	currentDate := time.Date(timeNowJkt.Year(), timeNowJkt.Month(), timeNowJkt.Day(), 0, 0, 0, 0, timeNowJkt.Location())
	transDate := currentDate.Format("2006-01-02")

	query := `SELECT 
				p.cust_id, p.pro_id, p.pro_code, p.bar_code, p.pro_name,
				p.vat, p.vat_bg, p.vat_lg_purch, p.vat_lg_sell, p.cogs, p.pro_status,
				p.pcat_id, pc.pcat_code, pc.pcat_name,
				br.brand_id, br.brand_code, br.brand_name, 
				br.pl_id, pl.pl_code, pl.pl_name,
				p.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
				p.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
				p.flavor_id, fv.flavor_code, fv.flavor_name, 
				p.ptype_id, pt.ptype_code, pt.ptype_name,
				p.psize_id, ps.psize_code, ps.psize_name,
				p.sup_id, su.sup_code, su.sup_name,
				p.principal_id, pr.principal_code, pr.principal_name,
				p.c_pro_id, cp.c_pro_code, cp.c_pro_name,
				p.is_main_pro, p.sort_no, p.item_no, p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5, 
				COALESCE(un1.unit_name, p.unit_id1, '') AS unit_name1,
				COALESCE(un2.unit_name, p.unit_id2, '') AS unit_name2,
				COALESCE(un3.unit_name, p.unit_id3, '') AS unit_name3,
				NULLIF(unc1.unit_id_coretax, '') AS unit_id_coretax1, NULLIF(unc2.unit_id_coretax, '') AS unit_id_coretax2, NULLIF(unc3.unit_id_coretax, '') AS unit_id_coretax3,
				NULLIF(unc1.unit_name_coretax, '') AS unit_name_coretax1, NULLIF(unc2.unit_name_coretax, '') AS unit_name_coretax2, NULLIF(unc3.unit_name_coretax, '') AS unit_name_coretax3,
				p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5, 
				p.is_batch, p.is_exp_date, 
				p.weight,p.length, p.width, p.height, p.volume,
				
				CASE WHEN mtp_mg_pr.purch_price1 IS NULL THEN p.purch_price1 
					ELSE mtp_mg_pr.purch_price1 END AS purch_price1,
					
				CASE WHEN mtp_mg_pr.purch_price2 IS NULL THEN p.purch_price2
					ELSE mtp_mg_pr.purch_price2 END AS purch_price2,
					
				CASE WHEN mtp_mg_pr.purch_price3 IS NULL THEN p.purch_price3 
					ELSE mtp_mg_pr.purch_price3 END AS purch_price3,

				p.purch_price4, p.purch_price5,

				CASE WHEN (mtp.sell_price1=0 or mtp.sell_price1 is null) THEN 
				(CASE WHEN (mtp_mg_pr_sell.sell_price1=0 or mtp_mg_pr_sell.sell_price1 is null) THEN p.sell_price1 
					ELSE mtp_mg_pr_sell.sell_price1 END)
				ELSE mtp.sell_price1 END AS sell_price1,

				CASE WHEN (mtp.sell_price2=0 or mtp.sell_price2 is null) THEN 
				(CASE WHEN (mtp_mg_pr_sell.sell_price2=0 or mtp_mg_pr_sell.sell_price2 is null) THEN p.sell_price2 
					ELSE mtp_mg_pr_sell.sell_price2 END)
				ELSE mtp.sell_price2 END AS sell_price2,

				CASE WHEN (mtp.sell_price3=0 or mtp.sell_price3 is null) THEN 
				(CASE WHEN (mtp_mg_pr_sell.sell_price3=0 or mtp_mg_pr_sell.sell_price3 is null) THEN p.sell_price3 
					ELSE mtp_mg_pr_sell.sell_price3 END)
				ELSE mtp.sell_price3 END AS sell_price3,
				p.sell_price4, p.sell_price5, 

				p.weight1, p.length1, p.width1, p.height1, p.volume1, 
				p.weight2, p.length2, p.width2, p.height2, p.volume2, 
				p.weight3, p.length3, p.width3, p.height3, p.volume3, 
				p.weight4, p.length4, p.width4, p.height4, p.volume4, 
				p.weight5, p.length5, p.width5, p.height5, p.volume5,  
				p.parent_pro_id, 
				p.excise_rate, p.excise_tax, 
				p.is_active, p.created_by, p.created_at, 
				p.updated_by, p.updated_at, p.is_del, p.deleted_by, p.deleted_at, p.image_url,
				p.saf_stock_unit_id, p.saf_stock_qty, p.min_stock_unit_id, p.min_stock_qty,
				saf_unit.unit_name AS saf_stock_unit_name,
				min_unit.unit_name AS min_stock_unit_name,
				uc.user_fullname AS created_by_name,
				u.user_fullname AS updated_by_name,
				dist.distributor_id, dist.distributor_name,
				parent.pro_code AS parent_pro_code, parent.pro_name AS parent_pro_name,p.pro_code_coretax,pct.pro_name_coretax
			FROM mst.m_product p
			LEFT JOIN sys.m_user u ON u.user_id = p.updated_by
			LEFT JOIN sys.m_user uc ON uc.user_id = p.created_by
			LEFT JOIN mst.m_distributor dist ON dist.distributor_id  = p.distributor_id AND dist.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit saf_unit ON saf_unit.unit_id = p.saf_stock_unit_id AND saf_unit.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit min_unit ON min_unit.unit_id = p.min_stock_unit_id AND min_unit.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + params.ParentCustID + `' 
			LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id AND pc.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id AND sb2.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id AND ps.cust_id = '` + params.ParentCustID + `' 
			LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id AND pt.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id AND fv.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id AND pr.cust_id = '` + params.ParentCustID + `' 
			LEFT JOIN mst.m_unit un1 ON un1.unit_id = p.unit_id1 AND un1.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit un2 ON un2.unit_id = p.unit_id2 AND un2.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit un3 ON un3.unit_id = p.unit_id3 AND un3.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit_coretax unc1 ON un1.unit_id_coretax = unc1.unit_id_coretax AND un1.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit_coretax unc2 ON un2.unit_id_coretax = unc2.unit_id_coretax AND un2.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_unit_coretax unc3 ON un3.unit_id_coretax = unc3.unit_id_coretax AND un3.cust_id = '` + params.ParentCustID + `'
			LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id
			LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id
			LEFT JOIN mst.m_product parent ON parent.pro_id = p.parent_pro_id
			LEFT JOIN mst.m_product_coretax pct ON pct.pro_code_coretax = p.pro_code_coretax
			LEFT JOIN mst.m_distributor dist ON dist.cust_id = '` + params.ParentCustID + `' AND dist.distributor_id = ` + fmt.Sprintf("%d", params.DistributorID) + `
			LEFT JOIN mst.m_transaction_price mtp ON mtp.pro_id = p.pro_id 
					AND mtp.cust_id = '` + params.CustID + `' AND mtp.outlet_id = 0
					AND ('` + transDate + `' BETWEEN mtp.start_date AND mtp.end_date) 
			LEFT JOIN LATERAL (
				SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3
				FROM mst.m_transaction_price mtp_mg_pr
				WHERE mtp_mg_pr.cust_id = '` + params.ParentCustID + `' 
					AND mtp_mg_pr.pro_id = p.pro_id
					AND mtp_mg_pr.start_date <= '` + transDate + `' 
					AND (mtp_mg_pr.distributor_id = (CASE WHEN mtp_mg_pr.coverage='N' THEN 0 ELSE dist.distributor_id END)
								OR mtp_mg_pr.price_group_reff = dist.dist_price_grp_id)	
				ORDER BY mtp_mg_pr.start_date DESC LIMIT 1
			) mtp_mg_pr ON true
			LEFT JOIN LATERAL (
				SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3 
				FROM mst.m_transaction_price mtp_mg_pr_sell
				WHERE mtp_mg_pr_sell.cust_id = '` + params.ParentCustID + `' 
					AND mtp_mg_pr_sell.pro_id = p.pro_id
					and mtp_mg_pr_sell.source = 10
					AND mtp_mg_pr_sell.start_date <= '` + transDate + `' 
					AND (mtp_mg_pr_sell.distributor_id = (CASE WHEN mtp_mg_pr_sell.coverage = 'N' THEN 0 ELSE dist.distributor_id END)
								OR mtp_mg_pr_sell.price_group_reff = dist.dist_price_grp_id)	
				ORDER BY mtp_mg_pr_sell.start_date DESC LIMIT 1
			) mtp_mg_pr_sell ON true
			WHERE p.pro_id = $1 
			AND p.cust_id = $2`
	err := repository.Get(&product, query, params.ProductId, params.ParentCustID)
	if err != nil {
		log.Error("productRepository, FindOneByProductIdAndDistributorId, err:", err.Error())
		return product, err
	}

	return product, nil
}

func (repository *productRepositoryImpl) FindOneByProductCodeAndCustId(productCode string, custId string) (model.Product, error) {
	product := model.Product{}
	query := `SELECT 
				p.cust_id, p.pro_id, p.pro_code, p.bar_code, p.pro_name, 
				p.vat, p.vat_bg, p.vat_lg_purch, p.vat_lg_sell, p.cogs, p.pro_status,
				p.pcat_id, p.sbrand1_id, p.sbrand2_id, p.flavor_id, p.ptype_id, 
				p.psize_id, p.sup_id, p.principal_id, p.c_pro_id, p.is_main_pro, 
				p.sort_no, p.item_no, 
				p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5, 
				p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5, 
				p.is_batch, p.is_exp_date,  
				p.weight, p.length, p.width, p.height, p.volume,
				p.purch_price1, p.sell_price1,
				p.purch_price2, p.sell_price2,
				p.purch_price3, p.sell_price3,
				p.purch_price4, p.sell_price4,
				p.purch_price5, p.sell_price5,
				p.weight1, p.length1, p.width1, p.height1, p.volume1, 
				p.weight2, p.length2, p.width2, p.height2, p.volume2, 
				p.weight3, p.length3, p.width3, p.height3, p.volume3, 
				p.weight4, p.length4, p.width4, p.height4, p.volume4, 
				p.weight5, p.length5, p.width5, p.height5, p.volume5,  
				p.parent_pro_id, p.excise_rate, p.excise_tax, 
				p.is_active, p.created_by, p.created_at, 
				p.updated_by, p.updated_at, p.is_del, p.deleted_by, p.deleted_at, p.image_url,
				p.saf_stock_unit_id, p.saf_stock_qty, p.min_stock_unit_id, p.min_stock_qty,
				saf_unit.unit_name AS saf_stock_unit_name,
				min_unit.unit_name AS min_stock_unit_name,
				uc.user_fullname AS created_by_name,
				u.user_fullname AS updated_by_name,
				dist.distributor_id, dist.distributor_name,
				parent.pro_code AS parent_pro_code, parent.pro_name AS parent_pro_name 
			FROM mst.m_product p
			LEFT JOIN sys.m_user u ON u.user_id = p.updated_by
			LEFT JOIN sys.m_user uc ON uc.user_id = p.created_by
			LEFT JOIN mst.m_distributor dist ON dist.distributor_id  = p.distributor_id
			LEFT JOIN mst.m_unit saf_unit ON saf_unit.unit_id = p.saf_stock_unit_id
			LEFT JOIN mst.m_unit min_unit ON min_unit.unit_id = p.min_stock_unit_id
			LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id
			LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id
			LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id
			LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id
			LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id
			LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id
			LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id
			LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id
			LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id
			LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id
			LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id
			LEFT JOIN mst.m_product parent ON parent.pro_id = p.parent_pro_id
			WHERE p.pro_code = $1 
			AND p.cust_id = $2
			AND p.is_del = false`
	err := repository.Get(&product, query, productCode, custId)
	if err != nil {
		log.Error("productRepository, FindOneByProductCodeAndCustId, err:", err.Error())
		return product, err
	}

	return product, nil
}

func (repository *productRepositoryImpl) FindDistributor(custId string) ([]model.MCustomer, error) {
	customers := []model.MCustomer{}
	querySelect := `SELECT cust_id, cust_name from smc.m_customer where parent_cust_id ='` + custId + `' AND cust_id != parent_cust_id `
	err := repository.Select(&customers, querySelect)
	if err != nil {
		log.Error("productRepository, FindDistributor, err:", err.Error())
		return customers, err
	}

	return customers, nil
}

func (repository *productRepositoryImpl) FindAll(dataFilter entity.ProductQueryFilter, custId string) ([]model.Product, int, int, error) {

	products := []model.Product{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` p.cust_id, p.pro_id, p.pro_code, p.bar_code, p.pro_name, 
					p.vat, p.vat_bg, p.vat_lg_purch, p.vat_lg_sell, p.cogs, p.pro_status,
					p.pcat_id, pc.pcat_code, pc.pcat_name,
					br.brand_id, br.brand_code, br.brand_name, 
					br.pl_id, pl.pl_code, pl.pl_name,
					p.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
					p.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
					p.flavor_id, fv.flavor_code, fv.flavor_name, 
					p.ptype_id, pt.ptype_code, pt.ptype_name,
					p.psize_id, ps.psize_code, ps.psize_name,
					p.sup_id, su.sup_code, su.sup_name,
					p.principal_id, pr.principal_code, pr.principal_name,
					p.c_pro_id, cp.c_pro_code, cp.c_pro_name,
					p.is_main_pro,  p.sort_no, p.item_no, 
					p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5, 
					p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5, 
					p.is_batch, p.is_exp_date, 
					p.weight, p.length, p.width, p.height, p.volume,
					p.purch_price1, p.sell_price1,
					p.purch_price2, p.sell_price2,
					p.purch_price3, p.sell_price3,
					p.purch_price4, p.sell_price4,
					p.purch_price5, p.sell_price5,
					p.weight1, p.length1, p.width1, p.height1, p.volume1, 
					p.weight2, p.length2, p.width2, p.height2, p.volume2, 
					p.weight3, p.length3, p.width3, p.height3, p.volume3, 
					p.weight4, p.length4, p.width4, p.height4, p.volume4, 
					p.weight5, p.length5, p.width5, p.height5, p.volume5,   
					p.parent_pro_id,p.excise_rate, p.excise_tax, 
					p.is_active, p.created_by, p.created_at, 
					p.updated_by, p.updated_at, p.is_del, p.deleted_by, p.deleted_at, p.image_url,
					p.saf_stock_unit_id, p.saf_stock_qty, p.min_stock_unit_id, p.min_stock_qty,
					saf_unit.unit_name AS saf_stock_unit_name,
					min_unit.unit_name AS min_stock_unit_name,p.pro_code_coretax,pct.pro_name_coretax,
					uc.user_fullname AS created_by_name,
					u.user_fullname AS updated_by_name,
					dist.distributor_id, dist.distributor_name`
	qWhere := ` LEFT JOIN sys.m_user u ON u.user_id = p.updated_by
				LEFT JOIN sys.m_user uc ON uc.user_id = p.created_by
				LEFT JOIN mst.m_distributor dist ON dist.distributor_id  = p.distributor_id AND dist.parent_cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_unit saf_unit ON saf_unit.unit_id = p.saf_stock_unit_id AND saf_unit.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_unit min_unit ON min_unit.unit_id = p.min_stock_unit_id AND min_unit.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + dataFilter.ParentCustId + `' 
				LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id AND pc.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id AND sb2.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id AND ps.cust_id = '` + dataFilter.ParentCustId + `' 
				LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id AND pt.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id AND fv.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id AND pr.cust_id = '` + dataFilter.ParentCustId + `' 
				LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id AND su.cust_id = '` + dataFilter.ParentCustId + `' 
				LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id
				LEFT JOIN mst.m_product_coretax pct ON pct.pro_code_coretax = p.pro_code_coretax AND pct.cust_id = '` + dataFilter.ParentCustId + `'
				WHERE p.is_del = false `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SupCode != "" {
		qWhere += ` AND su.sup_code = '` + dataFilter.SupCode + `' `
	}

	if len(dataFilter.SupID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.SupID, ",")
		qWhere += ` AND p.sup_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.PcatId) > 0 {
		intArrStr := str.ArrayToString(dataFilter.PcatId, ",")
		qWhere += ` AND p.pcat_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.PrincipalID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.PrincipalID, ",")
		qWhere += ` AND p.principal_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.BrandID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.BrandID, ",")
		qWhere += ` AND br.brand_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.ProductLineId) > 0 {
		intArrStr := str.ArrayToString(dataFilter.ProductLineId, ",")
		qWhere += ` AND br.pl_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.SubBrand1Id) > 0 {
		intArrStr := str.ArrayToString(dataFilter.SubBrand1Id, ",")
		qWhere += ` AND p.sbrand1_id IN (` + intArrStr + `) `
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND p.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND p.is_active = false `
		}
	}

	// override distributor id from jwt if exist, to make sure distributor only get their own data
	if dataFilter.JwtDistributorId != 0 {
		dataFilter.DistributorIds = []int64{dataFilter.JwtDistributorId}
		dataFilter.AllDistributor = false
	} else {
		// if all distributor is true, get all distributor id based on parent cust id and filter in where clause
		if dataFilter.AllDistributor {
			qWhere += fmt.Sprintf(" AND (dist.distributor_id IS NULL OR dist.distributor_id IN (%s))", "SELECT distributor_id FROM mst.m_distributor WHERE cust_id = '"+dataFilter.ParentCustId+"'")
		}
	}

	// only apply distributor filter if not all distributor and distributor_ids is provided
	if !dataFilter.AllDistributor && len(dataFilter.DistributorIds) > 0 {
		idsStr := make([]string, len(dataFilter.DistributorIds))
		for i, id := range dataFilter.DistributorIds {
			idsStr[i] = strconv.FormatInt(id, 10)
		}
		intArrStr := strings.Join(idsStr, ", ")
		qWhere += fmt.Sprintf(" AND dist.distributor_id IN (%s)", intArrStr)
	}

	// if distributor filter is not provided, default to get only principal products
	if dataFilter.JwtDistributorId == 0 &&
		!dataFilter.AllDistributor &&
		len(dataFilter.DistributorIds) == 0 {
		qWhere += fmt.Sprintf(" AND p.distributor_id IS NULL AND p.parent_pro_id = 0 AND p.cust_id = '%s' ", dataFilter.ParentCustId)
	}

	queryFrom := ` FROM mst.m_product p `
	queryCount := `SELECT ` + selectCount + ` ` + queryFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere

	// log.Error("productRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Error("productRepository, count total, err:", err.Error())
		return products, 0, 0, err
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
		querySelect += fmt.Sprintf(` ORDER BY %s`, sortBy)
	} else {
		sortBy := `p.pro_id`
		querySelect += fmt.Sprintf(` ORDER BY %s DESC`, sortBy)
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

	// log.Error("productRepository, querySelect:", querySelect)
	err = repository.Select(&products, querySelect)
	if err != nil {
		log.Error("productRepository, FindAllByCustId, err:", err.Error())
		return products, total, lastPage, err
	}

	return products, total, lastPage, nil
}

func (repository *productRepositoryImpl) FindAllByCustId(dataFilter entity.ProductQueryFilter, custId string) ([]model.Product, int, int, error) {

	products := []model.Product{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` p.cust_id, p.pro_id, p.pro_code, p.bar_code, p.pro_name, 
					p.vat, p.vat_bg, p.vat_lg_purch, p.vat_lg_sell, p.cogs, p.pro_status,
					p.pcat_id, pc.pcat_code, pc.pcat_name,
					br.brand_id, br.brand_code, br.brand_name, 
					br.pl_id, pl.pl_code, pl.pl_name,
					p.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
					p.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
					p.flavor_id, fv.flavor_code, fv.flavor_name, 
					p.ptype_id, pt.ptype_code, pt.ptype_name,
					p.psize_id, ps.psize_code, ps.psize_name,
					p.sup_id, su.sup_code, su.sup_name,
					p.principal_id, pr.principal_code, pr.principal_name,
					p.c_pro_id, cp.c_pro_code, cp.c_pro_name,
					p.is_main_pro,  p.sort_no, p.item_no, 
					p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5, 
					p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5, 
					p.is_batch, p.is_exp_date, 
					p.weight, p.length, p.width, p.height, p.volume,
					p.purch_price1, p.sell_price1,
					p.purch_price2, p.sell_price2,
					p.purch_price3, p.sell_price3,
					p.purch_price4, p.sell_price4,
					p.purch_price5, p.sell_price5,
					p.weight1, p.length1, p.width1, p.height1, p.volume1, 
					p.weight2, p.length2, p.width2, p.height2, p.volume2, 
					p.weight3, p.length3, p.width3, p.height3, p.volume3, 
					p.weight4, p.length4, p.width4, p.height4, p.volume4, 
					p.weight5, p.length5, p.width5, p.height5, p.volume5,   
					p.parent_pro_id,p.excise_rate, p.excise_tax, 
					p.is_active, p.created_by, p.created_at, 
					p.updated_by, p.updated_at, p.is_del, p.deleted_by, p.deleted_at, p.image_url,
					p.saf_stock_unit_id, p.saf_stock_qty, p.min_stock_unit_id, p.min_stock_qty,
					saf_unit.unit_name AS saf_stock_unit_name,
					min_unit.unit_name AS min_stock_unit_name,p.pro_code_coretax,pct.pro_name_coretax,
					uc.user_fullname AS created_by_name,
					u.user_fullname AS updated_by_name`
	qWhere := ` LEFT JOIN sys.m_user u ON u.user_id = p.updated_by
				LEFT JOIN sys.m_user uc ON uc.user_id = p.created_by
				LEFT JOIN mst.m_unit saf_unit ON saf_unit.unit_id = p.saf_stock_unit_id AND saf_unit.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_unit min_unit ON min_unit.unit_id = p.min_stock_unit_id AND min_unit.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + dataFilter.ParentCustId + `' 
				LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id AND pc.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id AND sb2.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id AND ps.cust_id = '` + dataFilter.ParentCustId + `' 
				LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id AND pt.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id AND fv.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id AND pr.cust_id = '` + dataFilter.ParentCustId + `' 
				LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id AND su.cust_id = '` + dataFilter.ParentCustId + `' 
				LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id
				LEFT JOIN mst.m_product_coretax pct ON pct.pro_code_coretax = p.pro_code_coretax AND pct.cust_id = '` + dataFilter.ParentCustId + `'
				WHERE p.is_del = false 
				AND p.cust_id = '` + dataFilter.ParentCustId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SupCode != "" {
		qWhere += ` AND su.sup_code = '` + dataFilter.SupCode + `' `
	}

	if len(dataFilter.SupID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.SupID, ",")
		qWhere += ` AND p.sup_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.PcatId) > 0 {
		intArrStr := str.ArrayToString(dataFilter.PcatId, ",")
		qWhere += ` AND p.pcat_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.PrincipalID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.PrincipalID, ",")
		qWhere += ` AND p.principal_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.BrandID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.BrandID, ",")
		qWhere += ` AND br.brand_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.ProductLineId) > 0 {
		intArrStr := str.ArrayToString(dataFilter.ProductLineId, ",")
		qWhere += ` AND br.pl_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.SubBrand1Id) > 0 {
		intArrStr := str.ArrayToString(dataFilter.SubBrand1Id, ",")
		qWhere += ` AND p.sbrand1_id IN (` + intArrStr + `) `
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND p.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND p.is_active = false `
		}
	}

	queryFrom := ` FROM mst.m_product p `
	queryCount := `SELECT ` + selectCount + ` ` + queryFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere

	// log.Error("productRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Error("productRepository, count total, err:", err.Error())
		return products, 0, 0, err
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
		sortBy := `p.pro_id`
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

	// log.Error("productRepository, querySelect:", querySelect)
	err = repository.Select(&products, querySelect)
	if err != nil {
		log.Error("productRepository, FindAllByCustId, err:", err.Error())
		return products, total, lastPage, err
	}

	return products, total, lastPage, nil
}

func (repository *productRepositoryImpl) FindAllByCustIdLookup(dataFilter entity.ProductQueryFilter, custId string) ([]model.Product, int, int, error) {

	products := []model.Product{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` p.cust_id, p.pro_id, p.pro_code, p.bar_code, p.pro_name, 
					p.vat, p.vat_bg, p.vat_lg_purch, p.vat_lg_sell, p.cogs, p.pro_status,
					p.pcat_id, pc.pcat_code, pc.pcat_name,
					br.brand_id, br.brand_code, br.brand_name, 
					br.pl_id, pl.pl_code, pl.pl_name,
					p.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
					p.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
					p.flavor_id, fv.flavor_code, fv.flavor_name, 
					p.ptype_id, pt.ptype_code, pt.ptype_name,
					p.psize_id, ps.psize_code, ps.psize_name,
					p.sup_id, su.sup_code, su.sup_name,
					p.principal_id, pr.principal_code, pr.principal_name,
					p.c_pro_id, cp.c_pro_code, cp.c_pro_name,
					p.is_main_pro, p.sort_no, p.item_no, 
					p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5, 
					p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5, 
					p.is_batch, p.is_exp_date, 
					p.weight,  p.length, p.width, p.height, p.volume,
					p.purch_price1, p.sell_price1,
					p.purch_price2, p.sell_price2,
					p.purch_price3, p.sell_price3,
					p.purch_price4, p.sell_price4,
					p.purch_price5, p.sell_price5,
					p.weight1, p.length1, p.width1, p.height1, p.volume1, 
					p.weight2, p.length2, p.width2, p.height2, p.volume2, 
					p.weight3, p.length3, p.width3, p.height3, p.volume3, 
					p.weight4, p.length4, p.width4, p.height4, p.volume4, 
					p.weight5, p.length5, p.width5, p.height5, p.volume5,   
					p.parent_pro_id,
					p.excise_rate, p.excise_tax, p.is_active, p.created_by, p.created_at, 
					p.updated_by, p.updated_at, p.is_del, p.deleted_by, p.deleted_at, p.image_url,
					p.saf_stock_unit_id, p.saf_stock_qty, p.min_stock_unit_id, p.min_stock_qty,
					saf_unit.unit_name AS saf_stock_unit_name,
					min_unit.unit_name AS min_stock_unit_name,
					u.user_fullname AS updated_by_name `
	qWhere := ` LEFT JOIN sys.m_user u ON u.user_id = p.updated_by AND u.cust_id = '` + custId + `'
				LEFT JOIN mst.m_unit saf_unit ON saf_unit.unit_id = p.saf_stock_unit_id AND saf_unit.cust_id = '` + custId + `'
				LEFT JOIN mst.m_unit min_unit ON min_unit.unit_id = p.min_stock_unit_id AND min_unit.cust_id = '` + custId + `'
				LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + custId + `'
				LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + custId + `'
				LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + custId + `'
				LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id AND pc.cust_id = '` + custId + `'
				LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id AND sb2.cust_id = '` + custId + `'
				LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id AND ps.cust_id = '` + custId + `'
				LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id AND pt.cust_id = '` + custId + `'
				LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id AND fv.cust_id = '` + custId + `'
				LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id AND pr.cust_id = '` + custId + `'
				LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id AND su.cust_id = '` + custId + `'
				LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id AND cp.cust_id = '` + custId + `'
				WHERE p.is_del = false AND p.is_active = true 
				AND p.cust_id = '` + dataFilter.ParentCustId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SupCode != "" {
		qWhere += ` AND su.sup_code = '` + dataFilter.SupCode + `' `
	}

	tempPCatId := []interface{}{}
	for _, xx := range dataFilter.PcatId {
		tempPCatId = append(tempPCatId, ""+strconv.Itoa(xx)+"")
	}
	varPCatId := make([]string, len(tempPCatId))
	for i, v := range tempPCatId {
		varPCatId[i] = fmt.Sprint(v)
	}
	var PCatId = strings.Join(varPCatId, ",")

	if len(dataFilter.PcatId) > 0 {
		qWhere += ` AND p.pcat_id in (` + PCatId + `) `
	}

	queryFrom := ` FROM mst.m_product p `
	queryCount := `SELECT ` + selectCount + ` ` + queryFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere

	// log.Error("productRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Error("productRepository, count total, err:", err.Error())
		return products, 0, 0, err
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
		sortBy := `p.pro_id`
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

	// log.Error("productRepository, querySelect:", querySelect)
	err = repository.Select(&products, querySelect)
	if err != nil {
		log.Error("productRepository, FindAllByCustId, err:", err.Error())
		return products, total, lastPage, err
	}

	return products, total, lastPage, nil
}

func (repository *productRepositoryImpl) FindAllByDistributorLookup(dataFilter entity.ProductQueryFilter, custId string) ([]model.Product, int, int, error) {

	products := []model.Product{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` p.cust_id, p.pro_id, p.pro_code, p.bar_code, p.pro_name, 
					p.vat, p.vat_bg, p.vat_lg_purch, p.vat_lg_sell, p.cogs, p.pro_status,
					p.pcat_id, pc.pcat_code, pc.pcat_name,
					br.brand_id, br.brand_code, br.brand_name, 
					br.pl_id, pl.pl_code, pl.pl_name,
					p.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
					p.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
					p.flavor_id, fv.flavor_code, fv.flavor_name, 
					p.ptype_id, pt.ptype_code, pt.ptype_name,
					p.psize_id, ps.psize_code, ps.psize_name,
					p.sup_id, su.sup_code, su.sup_name,
					p.principal_id, pr.principal_code, pr.principal_name,
					p.c_pro_id, cp.c_pro_code, cp.c_pro_name,
					p.is_main_pro, p.sort_no, p.item_no, 
					p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5, 
					p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5, 
					p.is_batch, p.is_exp_date, 
					p.weight,  p.length, p.width, p.height, p.volume,
					-- p.purch_price1, -- p.sell_price1,
					-- p.purch_price2, -- p.sell_price2,
					-- p.purch_price3, -- p.sell_price3,
					-- case when mtp.purch_price1 is null then 0 else mtp.purch_price1 end as purch_price1,
					-- case when mtp.sell_price1 is null then 0 else mtp.sell_price1 end as sell_price1,
					-- case when mtp.purch_price2 is null then 0 else mtp.purch_price2 end as purch_price2,
					-- case when mtp.sell_price2 is null then 0 else mtp.sell_price2 end as sell_price2,
					-- case when mtp.purch_price3 is null then 0 else mtp.purch_price3 end as purch_price3,
					-- case when mtp.sell_price3 is null then 0 else mtp.sell_price3 end as sell_price3,
					case when mtp.purch_price1 is null then p.purch_price1 else mtp.purch_price1 end as purch_price1,
					case when mtp.sell_price1 is null then p.sell_price1 else mtp.sell_price1 end as sell_price1,
					case when mtp.purch_price2 is null then p.purch_price2 else mtp.purch_price2 end as purch_price2,
					case when mtp.sell_price2 is null then p.sell_price2 else mtp.sell_price2 end as sell_price2,
					case when mtp.purch_price3 is null then p.purch_price3 else mtp.purch_price3 end as purch_price3,
					case when mtp.sell_price3 is null then p.sell_price3 else mtp.sell_price3 end as sell_price3,
					p.purch_price4, p.sell_price4,
					p.purch_price5, p.sell_price5,
					p.weight1, p.length1, p.width1, p.height1, p.volume1, 
					p.weight2, p.length2, p.width2, p.height2, p.volume2, 
					p.weight3, p.length3, p.width3, p.height3, p.volume3, 
					p.weight4, p.length4, p.width4, p.height4, p.volume4, 
					p.weight5, p.length5, p.width5, p.height5, p.volume5,   
					p.parent_pro_id,
					p.excise_rate, p.excise_tax, p.is_active, p.created_by, p.created_at, 
					p.updated_by, p.updated_at, p.is_del, p.deleted_by, p.deleted_at, p.image_url,
					p.saf_stock_unit_id, p.saf_stock_qty, p.min_stock_unit_id, p.min_stock_qty,
					saf_unit.unit_name AS saf_stock_unit_name,
					min_unit.unit_name AS min_stock_unit_name,
					u.user_fullname AS updated_by_name `
	qWhere := ` LEFT JOIN sys.m_user u ON u.user_id = p.updated_by AND u.cust_id = '` + custId + `'
				LEFT JOIN mst.m_unit saf_unit ON saf_unit.unit_id = p.saf_stock_unit_id AND saf_unit.cust_id = '` + custId + `'
				LEFT JOIN mst.m_unit min_unit ON min_unit.unit_id = p.min_stock_unit_id AND min_unit.cust_id = '` + custId + `'
				LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + custId + `'
				LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + custId + `'
				LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + custId + `'
				LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id AND pc.cust_id = '` + custId + `'
				LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id AND sb2.cust_id = '` + custId + `'
				LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id AND ps.cust_id = '` + custId + `'
				LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id AND pt.cust_id = '` + custId + `'
				LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id AND fv.cust_id = '` + custId + `'
				LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id AND pr.cust_id = '` + custId + `'
				LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id AND su.cust_id = '` + custId + `'
				LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id AND cp.cust_id = '` + custId + `'
				LEFT JOIN mst.m_transaction_price mtp on mtp.pro_id = p.pro_id and mtp.distributor_id = '` + fmt.Sprintf("%d", dataFilter.DistributorID) + `'
				WHERE p.is_del = false AND p.is_active = true 
				AND p.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SupCode != "" {
		qWhere += ` AND su.sup_code = '` + dataFilter.SupCode + `' `
	}

	tempPCatId := []interface{}{}
	for _, xx := range dataFilter.PcatId {
		tempPCatId = append(tempPCatId, ""+strconv.Itoa(xx)+"")
	}
	varPCatId := make([]string, len(tempPCatId))
	for i, v := range tempPCatId {
		varPCatId[i] = fmt.Sprint(v)
	}
	var PCatId = strings.Join(varPCatId, ",")

	if len(dataFilter.PcatId) > 0 {
		qWhere += ` AND p.pcat_id in (` + PCatId + `) `
	}

	queryFrom := ` FROM mst.m_product p `
	queryCount := `SELECT ` + selectCount + ` ` + queryFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere

	// log.Error("productRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Error("productRepository, count total, err:", err.Error())
		return products, 0, 0, err
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
		sortBy := `p.pro_id`
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

	// log.Error("productRepository, querySelect:", querySelect)
	err = repository.Select(&products, querySelect)
	if err != nil {
		log.Error("productRepository, FindAllByCustId, err:", err.Error())
		return products, total, lastPage, err
	}

	return products, total, lastPage, nil
}

func (repository *productRepositoryImpl) FindAllByCustIdLookupDistPrice(dataFilter entity.ProductQueryFilter) ([]model.ProductDistPrice, int, int, error) {

	var products []model.ProductDistPrice
	selectCount := ` COUNT(*) AS total `
	selectField := ` p.pro_id, p.pro_code, p.pro_name, 
					p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5, 
					p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5,
					p.purch_price1, p.purch_price2, p.purch_price3, p.purch_price4, p.purch_price5,
					p.sell_price1, p.sell_price2, p.sell_price3, p.sell_price4, p.sell_price5,p.vat,
					br.brand_id, br.brand_code, br.brand_name, 
					br.pl_id, pl.pl_code, pl.pl_name,
					p.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name `

	qWhere := `
				LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + dataFilter.ParentCustId + `' 
				LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id AND su.cust_id = '` + dataFilter.ParentCustId + `'
	WHERE p.is_del = false 
				AND p.cust_id = '` + dataFilter.ParentCustId + `'
				AND p.distributor_id IS NULL
				AND COALESCE(p.parent_pro_id, 0) = 0 `

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND p.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND p.is_active = false `
		}
		if *dataFilter.IsActive == 9 {
			qWhere += ` `
		}
	} else {
		qWhere += ` AND p.is_active = true `
	}

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SupCode != "" {
		qWhere += ` AND su.sup_code = '` + dataFilter.SupCode + `' `
	}

	tempPCatId := []interface{}{}
	for _, xx := range dataFilter.PcatId {
		tempPCatId = append(tempPCatId, ""+strconv.Itoa(xx)+"")
	}
	varPCatId := make([]string, len(tempPCatId))
	for i, v := range tempPCatId {
		varPCatId[i] = fmt.Sprint(v)
	}
	var PCatId = strings.Join(varPCatId, ",")

	if len(dataFilter.PcatId) > 0 {
		qWhere += ` AND p.pcat_id in (` + PCatId + `) `
	}

	if len(dataFilter.BrandID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.BrandID, ",")
		qWhere += ` AND br.brand_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.ProductLineId) > 0 {
		intArrStr := str.ArrayToString(dataFilter.ProductLineId, ",")
		qWhere += ` AND br.pl_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.SubBrand1Id) > 0 {
		intArrStr := str.ArrayToString(dataFilter.SubBrand1Id, ",")
		qWhere += ` AND p.sbrand1_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.ProID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.ProID, ",")
		qWhere += ` AND p.pro_id IN (` + intArrStr + `) `
	}

	queryFrom := ` FROM mst.m_product p `
	queryCount := `SELECT ` + selectCount + ` ` + queryFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere

	// log.Error("productRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Error("productRepository, count total, err:", err.Error())
		return products, 0, 0, err
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
		sortBy := `p.pro_id`
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

	// log.Error("productRepository, querySelect:", querySelect)
	err = repository.Select(&products, querySelect)
	if err != nil {
		log.Error("productRepository, FindAllByCustId, err:", err.Error())
		return products, total, lastPage, err
	}

	return products, total, lastPage, nil
}

func (repository *productRepositoryImpl) FindAllByDistributorLookupDistPrice(dataFilter entity.ProductQueryFilter) ([]model.ProductDistPrice, int, int, error) {
	var products []model.ProductDistPrice

	transDate := dataFilter.OrderDate
	if transDate == "" {
		asiaJkt, _ := time.LoadLocation("Asia/Jakarta")
		timeNowJkt := time.Now().In(asiaJkt)
		currentDate := time.Date(timeNowJkt.Year(), timeNowJkt.Month(), timeNowJkt.Day(), 0, 0, 0, 0, timeNowJkt.Location())
		transDate = currentDate.Format("2006-01-02")
	}

	distPriceGroupID, err := repository.resolveDistPriceGroupID(dataFilter.ParentCustId, dataFilter.DistributorID)
	if err != nil {
		return products, 0, 0, err
	}
	if distPriceGroupID == 0 {
		distPriceGroupID = dataFilter.DistPriceGroupId
	}

	baseSelect := ` p.pro_id, p.pro_code, p.pro_name,
					p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5,
					p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5,
					p.purch_price1, p.purch_price2, p.purch_price3, p.purch_price4, p.purch_price5,
					p.sell_price1, p.sell_price2, p.sell_price3, p.sell_price4, p.sell_price5,
					p.vat,
					br.brand_id, br.brand_code, br.brand_name,
					pl.pl_id, pl.pl_code, pl.pl_name,
					COALESCE(sb1.sbrand1_id, p.sbrand1_id) AS sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
					COALESCE(NULLIF(p.parent_pro_id, 0), p.pro_id) AS pricing_lookup_pro_id `

	baseWhere := `
	LEFT JOIN LATERAL (
		SELECT sb1.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name, sb1.brand_id
		FROM mst.m_sub_brand1 sb1
		WHERE sb1.sbrand1_id = p.sbrand1_id
			AND sb1.cust_id IN (p.cust_id, '` + dataFilter.ParentCustId + `')
		ORDER BY CASE WHEN sb1.cust_id = p.cust_id THEN 0 ELSE 1 END
		LIMIT 1
	) sb1 ON true
	LEFT JOIN LATERAL (
		SELECT br.brand_id, br.brand_code, br.brand_name, br.pl_id
		FROM mst.m_brand br
		WHERE br.brand_id = sb1.brand_id
			AND br.cust_id IN (p.cust_id, '` + dataFilter.ParentCustId + `')
		ORDER BY CASE WHEN br.cust_id = p.cust_id THEN 0 ELSE 1 END
		LIMIT 1
	) br ON true
	LEFT JOIN LATERAL (
		SELECT pl.pl_id, pl.pl_code, pl.pl_name
		FROM mst.m_product_line pl
		WHERE pl.pl_id = br.pl_id
			AND pl.cust_id IN (p.cust_id, '` + dataFilter.ParentCustId + `')
		ORDER BY CASE WHEN pl.cust_id = p.cust_id THEN 0 ELSE 1 END
		LIMIT 1
	) pl ON true
	LEFT JOIN LATERAL (
		SELECT su.sup_id, su.sup_code
		FROM mst.m_supplier su
		WHERE su.sup_id = p.sup_id
			AND su.cust_id IN (p.cust_id, '` + dataFilter.ParentCustId + `')
		ORDER BY CASE WHEN su.cust_id = p.cust_id THEN 0 ELSE 1 END
		LIMIT 1
	) su ON true
	WHERE p.cust_id = '` + dataFilter.CustId + `'
			AND p.distributor_id = ` + strconv.FormatInt(dataFilter.DistributorID, 10) + ` `

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			baseWhere += ` AND p.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			baseWhere += ` AND p.is_active = false `
		}
		if *dataFilter.IsActive == 9 {
			baseWhere += ` `
		}
	} else {
		baseWhere += ` AND p.is_active = true `
	}

	if !dataFilter.IncludeDeleted {
		baseWhere += ` AND p.is_del = false `
	}

	if dataFilter.Query != "" {
		baseWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SupCode != "" {
		baseWhere += ` AND su.sup_code = '` + dataFilter.SupCode + `' `
	}

	tempPCatId := []interface{}{}
	for _, xx := range dataFilter.PcatId {
		tempPCatId = append(tempPCatId, ""+strconv.Itoa(xx)+"")
	}
	varPCatId := make([]string, len(tempPCatId))
	for i, v := range tempPCatId {
		varPCatId[i] = fmt.Sprint(v)
	}
	var PCatId = strings.Join(varPCatId, ",")

	if len(dataFilter.PcatId) > 0 {
		baseWhere += ` AND p.pcat_id in (` + PCatId + `) `
	}

	if len(dataFilter.BrandID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.BrandID, ",")
		baseWhere += ` AND br.brand_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.ProductLineId) > 0 {
		intArrStr := str.ArrayToString(dataFilter.ProductLineId, ",")
		baseWhere += ` AND br.pl_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.SubBrand1Id) > 0 {
		intArrStr := str.ArrayToString(dataFilter.SubBrand1Id, ",")
		baseWhere += ` AND p.sbrand1_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.ProID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.ProID, ",")
		baseWhere += ` AND p.pro_id IN (` + intArrStr + `) `
	}

	queryFrom := ` FROM mst.m_product p `
	queryCount := `SELECT COUNT(*) AS total ` + queryFrom + ` ` + baseWhere

	// log.Error("productRepository, queryCount:", queryCount)
	var total int
	err = repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Error("productRepository, count total, err:", err.Error())
		return products, 0, 0, err
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
	} else {
		sortBy = `p.pro_id DESC`
	}

	finalSortBy := strings.NewReplacer(
		"p.", "priced_products.",
		"br.", "priced_products.",
		"pl.", "priced_products.",
		"sb1.", "priced_products.",
	).Replace(sortBy)

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))

	pagedBaseQuery := `SELECT ` + baseSelect + ` ` + queryFrom + ` ` + baseWhere + fmt.Sprintf(` ORDER BY %s LIMIT %s OFFSET %s`, sortBy, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))
	querySelect := `WITH paged_products AS MATERIALIZED (
		SELECT * FROM (` + pagedBaseQuery + `) base_page
	), priced_products AS MATERIALIZED (
		SELECT paged_products.*,
			CASE WHEN mtp_mg_pr.purch_price1 IS NULL THEN paged_products.purch_price1
				ELSE mtp_mg_pr.purch_price1 END AS effective_purch_price1,
			CASE WHEN mtp_mg_pr.purch_price2 IS NULL THEN paged_products.purch_price2
				ELSE mtp_mg_pr.purch_price2 END AS effective_purch_price2,
			CASE WHEN mtp_mg_pr.purch_price3 IS NULL THEN paged_products.purch_price3
				ELSE mtp_mg_pr.purch_price3 END AS effective_purch_price3,
			CASE WHEN (mtp.sell_price1 = 0 OR mtp.sell_price1 IS NULL) THEN
				(CASE WHEN (mtp_mg_pr_sell.sell_price1 = 0 OR mtp_mg_pr_sell.sell_price1 IS NULL) THEN paged_products.sell_price1
					ELSE mtp_mg_pr_sell.sell_price1 END)
				ELSE mtp.sell_price1 END AS effective_sell_price1,
			CASE WHEN (mtp.sell_price2 = 0 OR mtp.sell_price2 IS NULL) THEN
				(CASE WHEN (mtp_mg_pr_sell.sell_price2 = 0 OR mtp_mg_pr_sell.sell_price2 IS NULL) THEN paged_products.sell_price2
					ELSE mtp_mg_pr_sell.sell_price2 END)
				ELSE mtp.sell_price2 END AS effective_sell_price2,
			CASE WHEN (mtp.sell_price3 = 0 OR mtp.sell_price3 IS NULL) THEN
				(CASE WHEN (mtp_mg_pr_sell.sell_price3 = 0 OR mtp_mg_pr_sell.sell_price3 IS NULL) THEN paged_products.sell_price3
					ELSE mtp_mg_pr_sell.sell_price3 END)
				ELSE mtp.sell_price3 END AS effective_sell_price3
		FROM paged_products
		LEFT JOIN mst.m_transaction_price mtp ON mtp.pro_id = paged_products.pro_id
			AND mtp.cust_id = '` + dataFilter.CustId + `'
			AND mtp.outlet_id = ` + strconv.Itoa(dataFilter.OutletId) + `
			AND ('` + transDate + `' BETWEEN mtp.start_date AND mtp.end_date)
		LEFT JOIN LATERAL (
			-- child generic (scope_priority=0) wins over parent generic (scope_priority=1)
			-- when start_date is equal; otherwise latest start_date wins.
			SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3
			FROM (
				SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3,
					start_date, 0 AS scope_priority
				FROM mst.m_transaction_price
				WHERE cust_id = '` + dataFilter.CustId + `'
					AND pro_id = paged_products.pro_id
					AND COALESCE(outlet_id, 0) = 0
					AND start_date <= '` + transDate + `'
					AND (distributor_id = (CASE WHEN coverage = 'N' THEN 0 ELSE ` + strconv.FormatInt(dataFilter.DistributorID, 10) + ` END)
						OR price_group_reff = ` + strconv.Itoa(distPriceGroupID) + `)
				UNION ALL
				SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3,
					start_date, 1 AS scope_priority
				FROM mst.m_transaction_price
				WHERE cust_id = '` + dataFilter.ParentCustId + `'
					AND pro_id = paged_products.pricing_lookup_pro_id
					AND start_date <= '` + transDate + `'
					AND (distributor_id = (CASE WHEN coverage = 'N' THEN 0 ELSE ` + strconv.FormatInt(dataFilter.DistributorID, 10) + ` END)
						OR price_group_reff = ` + strconv.Itoa(distPriceGroupID) + `)
			) candidates
			ORDER BY start_date DESC, scope_priority ASC LIMIT 1
		) mtp_mg_pr ON true
		LEFT JOIN LATERAL (
			-- same UNION ALL logic but restricted to source=10 (sell price source)
			SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3
			FROM (
				SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3,
					start_date, 0 AS scope_priority
				FROM mst.m_transaction_price
				WHERE cust_id = '` + dataFilter.CustId + `'
					AND pro_id = paged_products.pro_id
					AND COALESCE(outlet_id, 0) = 0
					AND source = 10
					AND start_date <= '` + transDate + `'
					AND (distributor_id = (CASE WHEN coverage = 'N' THEN 0 ELSE ` + strconv.FormatInt(dataFilter.DistributorID, 10) + ` END)
						OR price_group_reff = ` + strconv.Itoa(distPriceGroupID) + `)
				UNION ALL
				SELECT purch_price1, purch_price2, purch_price3, sell_price1, sell_price2, sell_price3,
					start_date, 1 AS scope_priority
				FROM mst.m_transaction_price
				WHERE cust_id = '` + dataFilter.ParentCustId + `'
					AND pro_id = paged_products.pricing_lookup_pro_id
					AND source = 10
					AND start_date <= '` + transDate + `'
					AND (distributor_id = (CASE WHEN coverage = 'N' THEN 0 ELSE ` + strconv.FormatInt(dataFilter.DistributorID, 10) + ` END)
						OR price_group_reff = ` + strconv.Itoa(distPriceGroupID) + `)
			) candidates
			ORDER BY start_date DESC, scope_priority ASC LIMIT 1
		) mtp_mg_pr_sell ON true
	)
	SELECT priced_products.pro_id, priced_products.pro_code, priced_products.pro_name,
		priced_products.unit_id1, priced_products.unit_id2, priced_products.unit_id3, priced_products.unit_id4, priced_products.unit_id5,
		priced_products.conv_unit2, priced_products.conv_unit3, priced_products.conv_unit4, priced_products.conv_unit5,
		priced_products.effective_purch_price1 AS purch_price1,
		priced_products.effective_purch_price2 AS purch_price2,
		priced_products.effective_purch_price3 AS purch_price3,
		priced_products.purch_price4, priced_products.purch_price5,
		priced_products.effective_sell_price1 AS sell_price1,
		priced_products.effective_sell_price2 AS sell_price2,
		priced_products.effective_sell_price3 AS sell_price3,
		priced_products.sell_price4, priced_products.sell_price5,
		priced_products.vat,
		priced_products.brand_id, priced_products.brand_code, priced_products.brand_name,
		priced_products.pl_id, priced_products.pl_code, priced_products.pl_name,
		priced_products.sbrand1_id, priced_products.sbrand1_code, priced_products.sbrand1_name
	FROM priced_products
	ORDER BY ` + finalSortBy

	// log.Error("productRepository, querySelect:", querySelect)
	err = repository.Select(&products, querySelect)
	if err != nil {
		log.Error("productRepository, FindAllByCustId, err:", err.Error())
		return products, total, lastPage, err
	}

	return products, total, lastPage, nil
}

func (repository *productRepositoryImpl) FindAllByCustIdSearch(dataFilter entity.ProductQueryFilter, custId string) ([]model.Product, int, int, error) {

	products := []model.Product{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` p.cust_id, p.pro_id, p.pro_code, p.bar_code, p.pro_name, 
					p.unit_id1, 
					p.unit_id2, p.unit_id3, 
					p.unit_id4, p.unit_id5, 
					p.conv_unit2, p.conv_unit3, 
					p.conv_unit4, p.conv_unit5, 
					p.purch_price,
					p.purch_price1, p.purch_price2, p.purch_price3, p.purch_price4, p.purch_price5,  
					p.sell_price1, p.sell_price2, p.sell_price3, 
					p.sell_price4, p.sell_price5`
	qWhere := ` LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id
				WHERE p.is_del = false AND p.is_active = true 
				AND p.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SupCode != "" {
		qWhere += ` AND p.sup_code = '` + dataFilter.SupCode + `' `
	}

	tempPCatId := []interface{}{}
	for _, xx := range dataFilter.PcatId {
		tempPCatId = append(tempPCatId, ""+strconv.Itoa(xx)+"")
	}
	varPCatId := make([]string, len(tempPCatId))
	for i, v := range tempPCatId {
		varPCatId[i] = fmt.Sprint(v)
	}
	var PCatId = strings.Join(varPCatId, ",")

	if len(dataFilter.PcatId) > 0 {
		qWhere += ` AND p.pcat_id in (` + PCatId + `) `
	}

	queryFrom := ` FROM mst.m_product p `
	queryCount := `SELECT ` + selectCount + ` ` + queryFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere

	// log.Error("productRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Error("productRepository, FindAllByCustIdSearch, err:", err.Error())
		return products, 0, 0, err
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
		sortBy := `p.pro_id`
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

	// log.Error("productRepository, querySelect:", querySelect)
	err = repository.Select(&products, querySelect)
	if err != nil {
		log.Error("productRepository, FindAllByCustIdSearch, err:", err.Error())
		return products, total, lastPage, err
	}

	return products, total, lastPage, nil
}

func (repository *productRepositoryImpl) FindAllExport(dataFilter entity.ProductQueryFilter, custId string) ([]model.Product, int, error) {
	products := []model.Product{}

	selectCount := ` COUNT(*) AS total `
	selectField := ` p.cust_id, p.pro_id, p.pro_code, p.bar_code, p.pro_name, 
					p.vat, p.vat_bg, p.vat_lg_purch, p.vat_lg_sell, p.cogs, p.pro_status,
					p.pcat_id, pc.pcat_code, pc.pcat_name,
					br.brand_id, br.brand_code, br.brand_name, 
					br.pl_id, pl.pl_code, pl.pl_name,
					p.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
					p.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
					p.flavor_id, fv.flavor_code, fv.flavor_name, 
					p.ptype_id, pt.ptype_code, pt.ptype_name,
					p.psize_id, ps.psize_code, ps.psize_name,
					p.sup_id, su.sup_code, su.sup_name,
					p.principal_id, pr.principal_code, pr.principal_name,
					p.c_pro_id, cp.c_pro_code, cp.c_pro_name,
					p.is_main_pro,  p.sort_no, p.item_no, 
					p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5, 
					u1.unit_name AS unit_name1, u2.unit_name AS unit_name2, u3.unit_name AS unit_name3,
					p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5, 
					p.is_batch, p.is_exp_date, 
					p.weight, p.length, p.width, p.height, p.volume,
					p.purch_price1, p.sell_price1,
					p.purch_price2, p.sell_price2,
					p.purch_price3, p.sell_price3,
					p.purch_price4, p.sell_price4,
					p.purch_price5, p.sell_price5,
					p.weight1, p.length1, p.width1, p.height1, p.volume1, 
					p.weight2, p.length2, p.width2, p.height2, p.volume2, 
					p.weight3, p.length3, p.width3, p.height3, p.volume3, 
					p.weight4, p.length4, p.width4, p.height4, p.volume4, 
					p.weight5, p.length5, p.width5, p.height5, p.volume5,   
					p.parent_pro_id,p.excise_rate, p.excise_tax, 
					p.is_active, p.is_del, p.created_by, p.created_at, 
					p.updated_by, p.updated_at, p.is_del, p.deleted_by, p.deleted_at, p.image_url,
					p.saf_stock_unit_id, p.saf_stock_qty, p.min_stock_unit_id, p.min_stock_qty,
					saf_unit.unit_name AS saf_stock_unit_name,
					min_unit.unit_name AS min_stock_unit_name,p.pro_code_coretax,pct.pro_name_coretax,
					uc.user_fullname AS created_by_name,
					u.user_fullname AS updated_by_name,
					dist.distributor_id, dist.distributor_name`
	qWhere := ` LEFT JOIN sys.m_user u ON u.user_id = p.updated_by
				LEFT JOIN sys.m_user uc ON uc.user_id = p.created_by
				LEFT JOIN mst.m_distributor dist ON dist.distributor_id  = p.distributor_id AND dist.parent_cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_unit saf_unit ON saf_unit.unit_id = p.saf_stock_unit_id AND saf_unit.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_unit min_unit ON min_unit.unit_id = p.min_stock_unit_id AND min_unit.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + dataFilter.ParentCustId + `' 
				LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id AND pc.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id AND sb2.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id AND ps.cust_id = '` + dataFilter.ParentCustId + `' 
				LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id AND pt.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id AND fv.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id AND pr.cust_id = '` + dataFilter.ParentCustId + `' 
				LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id
				LEFT JOIN mst.m_unit u1 ON p.unit_id1 = u1.unit_id AND u1.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_unit u2 ON p.unit_id2 = u2.unit_id AND u2.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_unit u3 ON p.unit_id3 = u3.unit_id AND u3.cust_id = '` + dataFilter.ParentCustId + `'
				LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id
				LEFT JOIN mst.m_product_coretax pct ON pct.pro_code_coretax = p.pro_code_coretax AND pct.cust_id = '` + dataFilter.ParentCustId + `'
				WHERE p.is_del = false `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SupCode != "" {
		qWhere += ` AND su.sup_code = '` + dataFilter.SupCode + `' `
	}

	if len(dataFilter.SupID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.SupID, ",")
		qWhere += ` AND p.sup_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.PcatId) > 0 {
		intArrStr := str.ArrayToString(dataFilter.PcatId, ",")
		qWhere += ` AND p.pcat_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.PrincipalID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.PrincipalID, ",")
		qWhere += ` AND p.principal_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.BrandID) > 0 {
		intArrStr := str.ArrayToString(dataFilter.BrandID, ",")
		qWhere += ` AND br.brand_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.ProductLineId) > 0 {
		intArrStr := str.ArrayToString(dataFilter.ProductLineId, ",")
		qWhere += ` AND br.pl_id IN (` + intArrStr + `) `
	}

	if len(dataFilter.SubBrand1Id) > 0 {
		intArrStr := str.ArrayToString(dataFilter.SubBrand1Id, ",")
		qWhere += ` AND p.sbrand1_id IN (` + intArrStr + `) `
	}

	if dataFilter.Status != "" {
		switch strings.ToLower(dataFilter.Status) {
		case "active":
			qWhere += ` AND p.is_active = true `
		case "deactive":
			qWhere += ` AND p.is_active = false `
		default:
			if _, errConv := strconv.Atoi(dataFilter.Status); errConv == nil {
				qWhere += fmt.Sprintf(` AND p.pro_status = %s `, dataFilter.Status)
			}
		}
	}

	// override distributor id from jwt if exist, to make sure distributor only get their own data
	if dataFilter.JwtDistributorId != 0 {
		dataFilter.DistributorIds = []int64{dataFilter.JwtDistributorId}
		dataFilter.AllDistributor = false
	} else {
		// if all distributor is true, get all distributor id based on parent cust id and filter in where clause
		if dataFilter.AllDistributor {
			qWhere += fmt.Sprintf(" AND (dist.distributor_id IS NULL OR dist.distributor_id IN (%s))", "SELECT distributor_id FROM mst.m_distributor WHERE cust_id = '"+dataFilter.ParentCustId+"'")
		}
	}

	// only apply distributor filter if not all distributor and distributor_ids is provided
	if !dataFilter.AllDistributor && len(dataFilter.DistributorIds) > 0 {
		idsStr := make([]string, len(dataFilter.DistributorIds))
		for i, id := range dataFilter.DistributorIds {
			idsStr[i] = strconv.FormatInt(id, 10)
		}
		intArrStr := strings.Join(idsStr, ", ")
		qWhere += fmt.Sprintf(" AND dist.distributor_id IN (%s)", intArrStr)
	}

	// if distributor filter is not provided, default to get only principal products
	if dataFilter.JwtDistributorId == 0 &&
		!dataFilter.AllDistributor &&
		len(dataFilter.DistributorIds) == 0 {
		qWhere += fmt.Sprintf(" AND p.distributor_id IS NULL AND p.parent_pro_id = 0 AND p.cust_id = '%s' ", dataFilter.ParentCustId)
	}

	queryFrom := ` FROM mst.m_product p `
	queryCount := `SELECT ` + selectCount + ` ` + queryFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere

	// Hitung total rows
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Error("productRepository, FindAllExport, count err:", err.Error())
		return products, 0, err
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
		querySelect += ` ORDER BY p.pro_id DESC`
	}

	// Tidak ada LIMIT OFFSET
	err = repository.Select(&products, querySelect)
	if err != nil {
		log.Error("productRepository, FindAllExport, err:", err.Error())
		return products, total, err
	}

	return products, total, nil
}

func (r *productRepositoryImpl) FindByCodeProductCategory(custId, parentCustId, code, name string) (int64, error) {
	var id int64
	query := `SELECT pcat_id FROM mst.m_product_cat WHERE cust_id IN ($1, $2) AND pcat_code = $3 AND is_del = false`
	err := r.DB.Get(&id, query, custId, parentCustId, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("kode %s dan nama %s untuk Product Category tidak ditemukan, silakan setup parameter terlebih dahulu", code, name)
		}
		return 0, err
	}
	return id, nil
}
func (r *productRepositoryImpl) GetProductCategoryIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT pcat_id FROM mst.m_product_cat WHERE cust_id = $1 AND pcat_code = $2 AND is_del = false`
	err := r.DB.Get(&id, query, custId, code)
	if err == nil {
		return id, nil // Ditemukan, kembalikan ID.
	}
	return 0, err
}

func (r *productRepositoryImpl) FindByCodeProductLine(custId, parentCustId, code, name string) (int64, error) {
	var id int64
	query := `SELECT pl_id FROM mst.m_product_line WHERE cust_id IN ($1, $2) AND pl_code = $3 AND is_del = false`
	err := r.DB.Get(&id, query, custId, parentCustId, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("kode %s dan nama %s untuk Product Line tidak ditemukan, silakan setup parameter terlebih dahulu", code, name)
		}
		return 0, err
	}
	return id, nil
}

func (r *productRepositoryImpl) GetProductLineIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT pl_id FROM mst.m_product_line WHERE cust_id = $1 AND pl_code = $2 AND is_del = false`
	err := r.DB.Get(&id, query, custId, code)
	if err == nil {
		return id, nil // Ditemukan, kembalikan ID.
	}
	return 0, err
}

func (r *productRepositoryImpl) FindByCodeBrand(custId, parentCustId string, plId int64, code, name string) (int64, error) {
	var id int64
	query := `SELECT brand_id FROM mst.m_brand WHERE cust_id IN ($1, $2) AND brand_code = $3 AND is_del = false`
	err := r.DB.Get(&id, query, custId, parentCustId, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("kode %s dan nama %s untuk Brand tidak ditemukan, silakan setup parameter terlebih dahulu", code, name)
		}
		return 0, err
	}
	return id, nil
}

func (r *productRepositoryImpl) GetBrandIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT brand_id FROM mst.m_brand WHERE cust_id = $1 AND brand_code = $2 AND is_del = false`
	err := r.DB.Get(&id, query, custId, code)
	if err == nil {
		return id, nil // Ditemukan, kembalikan ID.
	}
	return 0, err
}

func (r *productRepositoryImpl) FindByCodeSubBrand1(custId, parentCustId string, brandId int64, code, name string) (int64, error) {
	var id int64
	query := `SELECT sbrand1_id FROM mst.m_sub_brand1 WHERE cust_id IN ($1, $2) AND sbrand1_code = $3 AND is_del = false`
	err := r.DB.Get(&id, query, custId, parentCustId, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("kode %s dan nama %s untuk Sub Brand tidak ditemukan, silakan setup parameter terlebih dahulu", code, name)
		}
		return 0, err
	}
	return id, nil
}

func (r *productRepositoryImpl) GetSubBrand1IdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT sbrand1_id FROM mst.m_sub_brand1 WHERE cust_id = $1 AND sbrand1_code = $2 AND is_del = false`
	err := r.DB.Get(&id, query, custId, code)
	if err == nil {
		return id, nil // Ditemukan, kembalikan ID.
	}
	return 0, err
}

func (r *productRepositoryImpl) FindByCodeSubBrand2(custId, parentCustId, code, name string) (int64, error) {
	var id int64
	query := `SELECT sbrand2_id FROM mst.m_sub_brand2 WHERE cust_id IN ($1, $2) AND sbrand2_code = $3 AND is_del = false`
	err := r.DB.Get(&id, query, custId, parentCustId, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("kode %s dan nama %s untuk Sub Brand 2 tidak ditemukan, silakan setup parameter terlebih dahulu", code, name)
		}
		return 0, err
	}
	return id, nil
}

func (r *productRepositoryImpl) GetSubBrand2IdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT sbrand2_id FROM mst.m_sub_brand2 WHERE cust_id = $1 AND sbrand2_code = $2 AND is_del = false`
	err := r.DB.Get(&id, query, custId, code)
	if err == nil {
		return id, nil // Ditemukan, kembalikan ID.
	}
	return 0, err
}

func (r *productRepositoryImpl) FindByCodeFlavor(custId, parentCustId, code, name string) (int64, error) {
	var id int64
	query := `SELECT flavor_id FROM mst.m_flavor WHERE cust_id IN ($1, $2) AND flavor_code = $3 AND is_del = false`
	err := r.DB.Get(&id, query, custId, parentCustId, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("kode %s dan nama %s untuk Flavor tidak ditemukan, silakan setup parameter terlebih dahulu", code, name)
		}
		return 0, err
	}
	return id, nil
}

func (r *productRepositoryImpl) GetFlavorIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT flavor_id FROM mst.m_flavor WHERE cust_id = $1 AND flavor_code = $2 AND is_del = false`
	err := r.DB.Get(&id, query, custId, code)
	if err == nil {
		return id, nil // Ditemukan, kembalikan ID.
	}
	return 0, err
}

func (r *productRepositoryImpl) FindByCodePackSize(custId, parentCustId, code, name string) (int64, error) {
	var id int64
	query := `SELECT psize_id FROM mst.m_pack_size WHERE cust_id IN ($1, $2) AND psize_code = $3 AND is_del = false`
	err := r.DB.Get(&id, query, custId, parentCustId, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("kode %s dan nama %s untuk Pack Size tidak ditemukan, silakan setup parameter terlebih dahulu", code, name)
		}
		return 0, err
	}
	return id, nil
}

func (r *productRepositoryImpl) GetPackSizeIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT psize_id FROM mst.m_pack_size WHERE cust_id = $1 AND psize_code = $2 AND is_del = false`
	err := r.DB.Get(&id, query, custId, code)
	if err == nil {
		return id, nil // Ditemukan, kembalikan ID.
	}
	return 0, err
}

func (r *productRepositoryImpl) FindByCodePackType(custId, parentCustId, code, name string) (int64, error) {
	var id int64
	query := `SELECT ptype_id FROM mst.m_pack_type WHERE cust_id IN ($1, $2) AND ptype_code = $3 AND is_del = false`
	err := r.DB.Get(&id, query, custId, parentCustId, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("kode %s dan nama %s untuk Pack Type tidak ditemukan, silakan setup parameter terlebih dahulu", code, name)
		}
		return 0, err
	}
	return id, nil
}

func (r *productRepositoryImpl) GetPackTypeIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT ptype_id FROM mst.m_pack_type WHERE cust_id = $1 AND ptype_code = $2 AND is_del = false`
	err := r.DB.Get(&id, query, custId, code)
	if err == nil {
		return id, nil // Ditemukan, kembalikan ID.
	}
	return 0, err
}

func (r *productRepositoryImpl) FindByCodeSupplier(custId, parentCustId, code, name string) (int64, error) {
	var id int64
	query := `SELECT sup_id FROM mst.m_supplier WHERE cust_id IN ($1, $2) AND sup_code = $3 AND is_del = false`
	err := r.DB.Get(&id, query, custId, parentCustId, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("kode %s dan nama %s untuk Supplier tidak ditemukan, silakan setup parameter terlebih dahulu", code, name)
		}
		return 0, err
	}
	return id, nil
}

func (r *productRepositoryImpl) GetSupplierIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT sup_id FROM mst.m_supplier WHERE cust_id = $1 AND sup_code = $2 AND is_del = false`
	err := r.DB.Get(&id, query, custId, code)
	if err == nil {
		return id, nil // Ditemukan, kembalikan ID.
	}
	return 0, err
}

func (r *productRepositoryImpl) FindByCodePrincipal(custId, parentCustId, code, name string) (int64, error) {
	var id int64
	query := `SELECT principal_id FROM mst.m_principal WHERE cust_id IN ($1, $2) AND principal_code = $3 AND is_del = false`
	err := r.DB.Get(&id, query, custId, parentCustId, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("kode %s dan nama %s untuk Principal tidak ditemukan, silakan setup parameter terlebih dahulu", code, name)
		}
		return 0, err
	}
	return id, nil
}

func (r *productRepositoryImpl) GetPrincipalIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT principal_id FROM mst.m_principal WHERE cust_id = $1 AND principal_code = $2 AND is_del = false`
	err := r.DB.Get(&id, query, custId, code)
	if err == nil {
		return id, nil // Ditemukan, kembalikan ID.
	}
	return 0, err
}

func (r *productRepositoryImpl) FindByCodeCPro(custId, parentCustId, code, name string) (int64, error) {
	var id int64
	query := `SELECT c_pro_id FROM mst.m_cons_product WHERE cust_id IN ($1, $2) AND c_pro_code = $3 AND is_del = false`
	err := r.DB.Get(&id, query, custId, parentCustId, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("kode %s dan nama %s untuk Consumer Product tidak ditemukan, silakan setup parameter terlebih dahulu", code, name)
		}
		return 0, err
	}
	return id, nil
}

func (r *productRepositoryImpl) GetCProIdByCode(custId, code string) (int64, error) {
	var id int64
	query := `SELECT c_pro_id FROM mst.m_cons_product WHERE cust_id = $1 AND c_pro_code = $2 AND is_del = false`
	err := r.DB.Get(&id, query, custId, code)
	if err == nil {
		return id, nil // Ditemukan, kembalikan ID.
	}
	return 0, err
}

func (r *productRepositoryImpl) FindByCodeUnit(custId, parentCustId, code, name string) (string, error) {
	var id string
	query := `SELECT unit_id FROM mst.m_unit WHERE cust_id IN ($1, $2) AND unit_id = $3 AND is_del = false`
	err := r.DB.Get(&id, query, custId, parentCustId, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("kode %s dan nama %s untuk Unit tidak ditemukan, silakan setup parameter terlebih dahulu", code, name)
		}
		return "", err
	}
	return id, nil
}

func (r *productRepositoryImpl) GetUnitIdByCode(custId, code string) (string, error) {
	var id string
	query := `SELECT unit_id FROM mst.m_unit WHERE cust_id = $1 AND unit_id = $2 AND is_del = false`
	err := r.DB.Get(&id, query, custId, code)
	if err == nil {
		return id, nil // Ditemukan, kembalikan ID.
	}
	return "", err
}

func (r *productRepositoryImpl) GetOrCreateProductCoretax(custId, parentCustId, code, name string, createdBy int64) (string, error) {
	var codeResult string
	query := `SELECT pro_code_coretax FROM mst.m_product_coretax WHERE cust_id IN ($1, $2) AND pro_code_coretax = $3 AND is_del = false`
	err := r.DB.Get(&codeResult, query, custId, parentCustId, code)
	if err == nil {
		return codeResult, nil
	}

	query = `INSERT INTO mst.m_product_coretax (
			cust_id, pro_code_coretax, pro_name_coretax, is_active,
			created_by, created_at, updated_by, updated_at
		) VALUES ($1, $2, $3, true, $4, NOW(), $4, NOW()) RETURNING pro_code_coretax`
	err = r.DB.QueryRow(query, custId, code, name, createdBy).Scan(&codeResult)
	if err != nil {
		return "", err
	}
	// codeResult = code
	return codeResult, nil
}

func (r *productRepositoryImpl) GetUnitProductCoretaxIdByCode(custId, parentCustId, code string) (string, error) {
	var id string
	query := `SELECT pro_code_coretax FROM mst.m_product_coretax WHERE cust_id IN ($1, $2) AND pro_code_coretax = $3 AND is_del = false`
	err := r.DB.Get(&id, query, custId, parentCustId, code)
	if err == nil {
		return id, nil // Ditemukan, kembalikan ID.
	}
	return "", err
}

func (r *productRepositoryImpl) CreateProduct(data entity.ProcessedProductRow) error {
	query := `INSERT INTO mst.m_product (
			cust_id, pro_code, pro_name, bar_code,
			pcat_id, sbrand1_id, sbrand2_id,
			flavor_id, psize_id, ptype_id, sup_id, principal_id, c_pro_id,
			is_main_pro, sort_no, item_no,
			unit_id1, unit_id2, unit_id3,
			conv_unit2, conv_unit3,
			is_batch, is_exp_date,
			length, width, height, weight, volume,
			parent_pro_id, is_new_pro,
			purch_price1, purch_price2, purch_price3, purch_price4, purch_price5,
			sell_price1, sell_price2, sell_price3, sell_price4, sell_price5,
			length1, length2, length3, length4, length5,
			width1, width2, width3, width4, width5,
			height1, height2, height3, height4, height5,
			weight1, weight2, weight3, weight4, weight5,
			volume1, volume2, volume3, volume4, volume5,
			saf_stock_qty, saf_stock_unit_id,
			min_stock_qty, min_stock_unit_id,
			excise_rate, excise_tax,
			is_active, is_del,
			image_url, vat, vat_bg, vat_lg_purch, vat_lg_sell, cogs,
			pro_status, pro_code_coretax, 
			distributor_id, level, origin, assigner_user_id,
			created_by, created_at, updated_by, updated_at
		)
		VALUES (
			:cust_id, :pro_code, :pro_name, :bar_code,
			:pcat_id, :sbrand1_id, :sbrand2_id,
			:flavor_id, :psize_id, :ptype_id, :sup_id, :principal_id, :c_pro_id,
			:is_main_pro, :sort_no, :item_no,
			:unit_id1, :unit_id2, :unit_id3,
			:conv_unit2, :conv_unit3,
			:is_batch, :is_exp_date,
			:length, :width, :height, :weight, :volume,
			:parent_pro_id, :is_new_pro,
			:purch_price1, :purch_price2, :purch_price3, :purch_price4, :purch_price5,
			:sell_price1, :sell_price2, :sell_price3, :sell_price4, :sell_price5,
			:length1, :length2, :length3, :length4, :length5,
			:width1, :width2, :width3, :width4, :width5,
			:height1, :height2, :height3, :height4, :height5,
			:weight1, :weight2, :weight3, :weight4, :weight5,
			:volume1, :volume2, :volume3, :volume4, :volume5,
			:saf_stock_qty, :saf_stock_unit_id,
			:min_stock_qty, :min_stock_unit_id,
			:excise_rate, :excise_tax,
			:is_active, :is_del,
			:image_url, :vat, :vat_bg, :vat_lg_purch, :vat_lg_sell, :cogs,
			:pro_status, :pro_code_coretax, 
			:distributor_id, :level, :origin, :assigner_user_id,
			:created_by, :created_at, :updated_by, :updated_at
		)`

	_, err := r.DB.NamedExec(query, data)
	return err
	// return nil
}

func (r *productRepositoryImpl) GetProductImportInstructions() ([]model.ImportInstruction, error) {
	query := `
		SELECT instruction_id, instruction_type, kolom,
			CASE
				WHEN mandatory is true
				THEN 'Kolom Wajib Diisi'
				ELSE 'Kolom Tidak Wajib Diisi'
			END AS mandatory,
			keterangan, step, 
			CASE
				WHEN step = 'Step 1' THEN 'Hijau Muda'
				WHEN step = 'Step 2' THEN 'Kuning Muda'
				WHEN step = 'Step 3' THEN 'Merah'
				ELSE ''
			END as color
		FROM import.import_instructions
		WHERE instruction_type = 'product'
		ORDER BY 
			CASE 
				WHEN step ILIKE 'Step 1%' THEN 1
				WHEN step ILIKE 'Step 2%' THEN 2
				WHEN step ILIKE 'Step 3%' THEN 3
				ELSE 4
			END, instruction_id;
	`

	var instructions []model.ImportInstruction
	if err := r.DB.Select(&instructions, query); err != nil {
		return nil, err
	}
	return instructions, nil
}

func (r *productRepositoryImpl) CreateImportHistory(uploadType, fileName, custId string, uploadedBy int64, totalData int) (int64, error) {
	var id int64
	query := `
        INSERT INTO import.import_history (file_name, uploaded_by, total_data, upload_type, cust_id)
        VALUES ($1, $2, $3, $4, $5) RETURNING history_id
    `
	err := r.DB.QueryRow(query, fileName, uploadedBy, totalData, uploadType, custId).Scan(&id)
	return id, err
}

func (r *productRepositoryImpl) UpdateImportHistory(historyId int64, success, failed int, statusReupload bool) error {
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

func (r *productRepositoryImpl) FindProductById(custId string, productId int) (*entity.ProcessedProductRow, error) {
	query := `
        SELECT 
            pro_id, cust_id, pro_code, bar_code, pro_name,
			length, width, height, volume,
            length1, width1, height1, volume1,
            length2, width2, height2, volume2,
            length3, width3, height3, volume3,
            length4, width4, height4, volume4,
            length5, width5, height5, volume5,
			is_active, is_del
        FROM mst.m_product
        WHERE cust_id = $1 AND pro_id = $2 AND is_del = false
        LIMIT 1;
    `
	var p entity.ProcessedProductRow
	err := r.DB.Get(&p, query, custId, productId)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepositoryImpl) InsertProductTemp(historyId int64, statusInsert, CustId string, data entity.ImportProductRow) error {
	data.HistoryId = historyId
	data.CustId = CustId
	data.StatusInsert = statusInsert
	query := `
        INSERT INTO import.product_temp (
            history_id, cust_id,
            pro_code, bar_code, pro_name,
            pcat_code, pcat_name,
            pl_code, pl_name,
            brand_code, brand_name,
            sbrand1_code, sbrand1_name,
            sbrand2_code, sbrand2_name,
            flavor_code, flavor_name,
            ptype_code, ptype_name,
            psize_code, psize_name,
            sup_code, sup_name,
            principal_code, principal_name,
            c_pro_code, c_pro_name,
            unit_id1, unit_id2, unit_id3,
            unit_name1, unit_name2, unit_name3,
            conv_unit2, conv_unit3,
            is_batch, is_exp_date,
            parent_pro_id,
            purch_price1, purch_price2, purch_price3,
            sell_price1, sell_price2, sell_price3,
            length1, length2, length3,
            width1, width2, width3,
            height1, height2, height3,
            weight1, weight2, weight3,
            saf_stock_qty, saf_stock_unit_id,
            min_stock_qty, min_stock_unit_id,
            excise_rate, excise_tax,
            is_active,
            vat, vat_bg, vat_lg_purch, vat_lg_sell,
            cogs, pro_status,
            pro_code_coretax,
			distributor_id, level,
            status_insert, error_message
        )
        VALUES (
            :history_id, :cust_id,
            :pro_code, :bar_code, :pro_name,
            :pcat_code, :pcat_name,
            :pl_code, :pl_name,
            :brand_code, :brand_name,
            :sbrand1_code, :sbrand1_name,
            :sbrand2_code, :sbrand2_name,
            :flavor_code, :flavor_name,
            :ptype_code, :ptype_name,
            :psize_code, :psize_name,
            :sup_code, :sup_name,
            :principal_code, :principal_name,
            :c_pro_code, :c_pro_name,
            :unit_id1, :unit_id2, :unit_id3,
            :unit_name1, :unit_name2, :unit_name3,
            :conv_unit2, :conv_unit3,
            :is_batch, :is_exp_date,
            :parent_pro_id,
            :purch_price1, :purch_price2, :purch_price3,
            :sell_price1, :sell_price2, :sell_price3,
            :length1, :length2, :length3,
            :width1, :width2, :width3,
            :height1, :height2, :height3,
            :weight1, :weight2, :weight3,
            :saf_stock_qty, :saf_stock_unit_id,
            :min_stock_qty, :min_stock_unit_id,
            :excise_rate, :excise_tax,
            :is_active,
            :vat, :vat_bg, :vat_lg_purch, :vat_lg_sell,
            :cogs, :pro_status,
            :pro_code_coretax,
			:distributor_id, :level,
            :status_insert, :error_message
        )`

	_, err := r.DB.NamedExec(query, data)
	return err
}

func (r *productRepositoryImpl) CheckProductExists(custId, parentCustId, proCode string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_product WHERE cust_id IN ($1, $2) AND pro_code = $3 AND deleted_at IS NULL`
	err := r.DB.Get(&exists, query, custId, parentCustId, proCode)
	return exists, err
}

func (r *productRepositoryImpl) CheckBrandExists(custId string, brandId int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_brand WHERE cust_id = $1 AND brand_id = $2`
	err := r.DB.Get(&exists, query, custId, brandId)
	return exists, err
}

func (r *productRepositoryImpl) CheckBrandCodeDuplicate(custId string, brandId int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_brand WHERE cust_id = $1 AND brand_code = $2 AND brand_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, brandId)
	return exists, err
}

func (r *productRepositoryImpl) CheckProductCatExists(custId string, pcatId int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_product_cat WHERE cust_id = $1 AND pcat_id = $2`
	err := r.DB.Get(&exists, query, custId, pcatId)
	return exists, err
}

func (r *productRepositoryImpl) CheckProductCatCodeDuplicate(custId string, pcatId int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_product_cat WHERE cust_id = $1 AND pcat_code = $2 AND pcat_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, pcatId)
	return exists, err
}

func (r *productRepositoryImpl) CheckProductLineExists(custId string, plId int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_product_line WHERE cust_id = $1 AND pl_id = $2`
	err := r.DB.Get(&exists, query, custId, plId)
	return exists, err
}

func (r *productRepositoryImpl) CheckProductLineCodeDuplicate(custId string, plId int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_product_line WHERE cust_id = $1 AND pl_code = $2 AND pl_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, plId)
	return exists, err
}

func (r *productRepositoryImpl) CheckSbrand1Exists(custId string, id int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_sub_brand1 WHERE cust_id = $1 AND sbrand1_id = $2`
	err := r.DB.Get(&exists, query, custId, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckSbrand1CodeDuplicate(custId string, id int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_sub_brand1 WHERE cust_id = $1 AND sbrand1_code = $2 AND sbrand1_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckSbrand2Exists(custId string, id int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_sub_brand2 WHERE cust_id = $1 AND sbrand2_id = $2`
	err := r.DB.Get(&exists, query, custId, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckSbrand2CodeDuplicate(custId string, id int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_sub_brand2 WHERE cust_id = $1 AND sbrand2_code = $2 AND sbrand2_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckFlavorExists(custId string, id int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_flavor WHERE cust_id = $1 AND flavor_id = $2`
	err := r.DB.Get(&exists, query, custId, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckFlavorCodeDuplicate(custId string, id int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_flavor WHERE cust_id = $1 AND flavor_code = $2 AND flavor_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckPackTypeExists(custId string, id int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_pack_type WHERE cust_id = $1 AND ptype_id = $2`
	err := r.DB.Get(&exists, query, custId, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckPackTypeCodeDuplicate(custId string, id int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_pack_type WHERE cust_id = $1 AND ptype_code = $2 AND ptype_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckPackSizeExists(custId string, id int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_pack_size WHERE cust_id = $1 AND psize_id = $2`
	err := r.DB.Get(&exists, query, custId, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckPackSizeCodeDuplicate(custId string, id int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_pack_size WHERE cust_id = $1 AND psize_code = $2 AND psize_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckSupplierExists(custId string, id int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_supplier WHERE cust_id = $1 AND sup_id = $2`
	err := r.DB.Get(&exists, query, custId, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckSupplierCodeDuplicate(custId string, id int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_supplier WHERE cust_id = $1 AND sup_code = $2 AND sup_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckPrincipalExists(custId string, id int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_principal WHERE cust_id = $1 AND principal_id = $2`
	err := r.DB.Get(&exists, query, custId, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckPrincipalCodeDuplicate(custId string, id int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_principal WHERE cust_id = $1 AND principal_code = $2 AND principal_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckCProExists(custId string, id int64) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_cons_product WHERE cust_id = $1 AND c_pro_id = $2`
	err := r.DB.Get(&exists, query, custId, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckCProCodeDuplicate(custId string, id int64, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_cons_product WHERE cust_id = $1 AND c_pro_code = $2 AND c_pro_id <> $3`
	err := r.DB.Get(&exists, query, custId, code, id)
	return exists, err
}

func (r *productRepositoryImpl) CheckCoretaxExists(custId string, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_product_coretax WHERE cust_id = $1 AND pro_code_coretax = $2`
	err := r.DB.Get(&exists, query, custId, code)
	return exists, err
}

func (r *productRepositoryImpl) CheckUnitExists(custId string, code string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(1) > 0 FROM mst.m_unit WHERE cust_id = $1 AND unit_id = $2`
	err := r.DB.Get(&exists, query, custId, code)
	return exists, err
}

func (r *productRepositoryImpl) UpdateImportBrand(custId string, id int64, code, name string) error {
	query := `
		UPDATE mst.m_brand
		SET brand_code = COALESCE(NULLIF($1, ''), brand_code),
		    brand_name = COALESCE(NULLIF($2, ''), brand_name)
		WHERE cust_id = $3 AND brand_id = $4
	`
	_, err := r.DB.Exec(query, code, name, custId, id)
	return err
}

func (r *productRepositoryImpl) UpdateImportCategory(custId string, id int64, code, name string) error {
	query := `
		UPDATE mst.m_product_cat
		SET pcat_code = COALESCE(NULLIF($1, ''), pcat_code),
		    pcat_name = COALESCE(NULLIF($2, ''), pcat_name)
		WHERE cust_id = $3 AND pcat_id = $4
	`
	_, err := r.DB.Exec(query, code, name, custId, id)
	return err
}

func (r *productRepositoryImpl) UpdateImportProductLine(custId string, id int64, code, name string) error {
	query := `
		UPDATE mst.m_product_line
		SET pl_code = COALESCE(NULLIF($1, ''), pl_code),
		    pl_name = COALESCE(NULLIF($2, ''), pl_name)
		WHERE cust_id = $3 AND pl_id = $4
	`
	_, err := r.DB.Exec(query, code, name, custId, id)
	return err
}

func (r *productRepositoryImpl) UpdateImportSubBrand1(custId string, id int64, code, name string) error {
	query := `
		UPDATE mst.m_sub_brand1
		SET sbrand1_code = COALESCE(NULLIF($1, ''), sbrand1_code),
		    sbrand1_name = COALESCE(NULLIF($2, ''), sbrand1_name)
		WHERE cust_id = $3 AND sbrand1_id = $4
	`
	_, err := r.DB.Exec(query, code, name, custId, id)
	return err
}

func (r *productRepositoryImpl) UpdateImportSubBrand2(custId string, sbrand2Id int64, code, name string) error {
	query := `
		UPDATE mst.m_sub_brand2
		SET sbrand2_code = COALESCE(NULLIF($1, ''), sbrand2_code),
		    sbrand2_name = COALESCE(NULLIF($2, ''), sbrand2_name)
		WHERE sbrand2_id = $3 AND cust_id = $4
	`
	_, err := r.DB.Exec(query, code, name, sbrand2Id, custId)
	return err
}

func (r *productRepositoryImpl) UpdateImportFlavor(custId string, flavorId int64, code, name string) error {
	query := `
		UPDATE mst.m_flavor
		SET flavor_code = COALESCE(NULLIF($1, ''), flavor_code),
		    flavor_name = COALESCE(NULLIF($2, ''), flavor_name)
		WHERE flavor_id = $3 AND cust_id = $4
	`
	_, err := r.DB.Exec(query, code, name, flavorId, custId)
	return err
}

func (r *productRepositoryImpl) UpdateImportPackType(custId string, ptypeId int64, code, name string) error {
	query := `
		UPDATE mst.m_pack_type
		SET ptype_code = COALESCE(NULLIF($1, ''), ptype_code),
		    ptype_name = COALESCE(NULLIF($2, ''), ptype_name)
		WHERE ptype_id = $3 AND cust_id = $4
	`
	_, err := r.DB.Exec(query, code, name, ptypeId, custId)
	return err
}

func (r *productRepositoryImpl) UpdateImportPackSize(custId string, psizeId int64, code, name string) error {
	query := `
		UPDATE mst.m_pack_size
		SET psize_code = COALESCE(NULLIF($1, ''), psize_code),
		    psize_name = COALESCE(NULLIF($2, ''), psize_name)
		WHERE psize_id = $3 AND cust_id = $4
	`
	_, err := r.DB.Exec(query, code, name, psizeId, custId)
	return err
}

func (r *productRepositoryImpl) UpdateImportSupplier(custId string, supId int64, code, name string) error {
	query := `
		UPDATE mst.m_supplier
		SET sup_code = COALESCE(NULLIF($1, ''), sup_code),
		    sup_name = COALESCE(NULLIF($2, ''), sup_name)
		WHERE sup_id = $3 AND cust_id = $4
	`
	_, err := r.DB.Exec(query, code, name, supId, custId)
	return err
}

func (r *productRepositoryImpl) UpdateImportPrincipal(custId string, principalId int64, code, name string) error {
	query := `
		UPDATE mst.m_principal
		SET principal_code = COALESCE(NULLIF($1, ''), principal_code),
		    principal_name = COALESCE(NULLIF($2, ''), principal_name)
		WHERE principal_id = $3 AND cust_id = $4
	`
	_, err := r.DB.Exec(query, code, name, principalId, custId)
	return err
}

func (r *productRepositoryImpl) UpdateImportCPro(custId string, cproId int64, code, name string) error {
	query := `
		UPDATE mst.m_cons_product
		SET c_pro_code = COALESCE(NULLIF($1, ''), c_pro_code),
		    c_pro_name = COALESCE(NULLIF($2, ''), c_pro_name)
		WHERE c_pro_id = $3 AND cust_id = $4
	`
	_, err := r.DB.Exec(query, code, name, cproId, custId)
	return err
}

func (r *productRepositoryImpl) UpdateImportProductCoretax(custId string, code, name string) error {
	query := `
		UPDATE mst.m_product_coretax
		SET pro_name_coretax = COALESCE(NULLIF($1, ''), pro_name_coretax)
		WHERE pro_code_coretax = $2 AND cust_id = $3
	`
	_, err := r.DB.Exec(query, name, code, custId)
	return err
}

func (r *productRepositoryImpl) UpdateImportUnit(custId string, id, name string) error {
	query := `
		UPDATE mst.m_unit
		SET unit_name = COALESCE(NULLIF($1, ''), unit_name)
		WHERE unit_id = $2 AND cust_id = $3
	`
	_, err := r.DB.Exec(query, name, id, custId)
	return err
}

func (r *productRepositoryImpl) FindProductIdByCode(custId, proCode string) (int, error) {
	var id int
	query := `
        SELECT pro_id 
        FROM mst.m_product 
        WHERE cust_id = $1 
          AND pro_code = $2 
          AND is_del = FALSE
        LIMIT 1
    `
	err := r.DB.QueryRow(query, custId, proCode).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *productRepositoryImpl) UpdateImportProduct(p entity.ProcessedProductRow) error {
	query := `
    UPDATE mst.m_product SET
        pro_code        = COALESCE(NULLIF(:pro_code, ''), pro_code),
		bar_code        = COALESCE(NULLIF(:bar_code, ''), bar_code),
		pro_name        = COALESCE(NULLIF(:pro_name, ''), pro_name),

		pcat_id         = COALESCE(NULLIF(:pcat_id, 0), pcat_id),
		sbrand1_id      = COALESCE(NULLIF(:sbrand1_id, 0), sbrand1_id),
		sbrand2_id      = COALESCE(NULLIF(:sbrand2_id, 0), sbrand2_id),
		flavor_id       = COALESCE(NULLIF(:flavor_id, 0), flavor_id),
		ptype_id        = COALESCE(NULLIF(:ptype_id, 0), ptype_id),
		psize_id        = COALESCE(NULLIF(:psize_id, 0), psize_id),
		sup_id          = COALESCE(NULLIF(:sup_id, 0), sup_id),
		principal_id    = COALESCE(NULLIF(:principal_id, 0), principal_id),
		c_pro_id        = COALESCE(NULLIF(:c_pro_id, 0), c_pro_id),

		is_main_pro     = COALESCE(:is_main_pro, is_main_pro),
		sort_no         = COALESCE(NULLIF(:sort_no, 0), sort_no),
		item_no         = COALESCE(NULLIF(:item_no, 0), item_no),

		unit_id1        = COALESCE(NULLIF(:unit_id1, ''), unit_id1),
		unit_id2        = COALESCE(NULLIF(:unit_id2, ''), unit_id2),
		unit_id3        = COALESCE(NULLIF(:unit_id3, ''), unit_id3),

		conv_unit2      = COALESCE(NULLIF(:conv_unit2, 0), conv_unit2),
		conv_unit3      = COALESCE(NULLIF(:conv_unit3, 0), conv_unit3),
		weight          = COALESCE(NULLIF(:weight, 0), weight),

		is_batch        = COALESCE(:is_batch, is_batch),
		is_exp_date     = COALESCE(:is_exp_date, is_exp_date),

		length          = COALESCE(NULLIF(:length, 0), length),
		width           = COALESCE(NULLIF(:width, 0), width),
		height          = COALESCE(NULLIF(:height, 0), height),
		volume          = COALESCE(NULLIF(:volume, 0), volume),

		parent_pro_id   = COALESCE(NULLIF(:parent_pro_id, 0), parent_pro_id),
		is_new_pro      = COALESCE(:is_new_pro, is_new_pro),

		vat             = COALESCE(NULLIF(:vat, 0), vat),
		vat_bg          = COALESCE(NULLIF(:vat_bg, 0), vat_bg),
		vat_lg_purch    = COALESCE(NULLIF(:vat_lg_purch, 0), vat_lg_purch),
		vat_lg_sell     = COALESCE(NULLIF(:vat_lg_sell, 0), vat_lg_sell),
		excise_rate     = COALESCE(NULLIF(:excise_rate, 0), excise_rate),
		excise_tax      = COALESCE(NULLIF(:excise_tax, 0), excise_tax),

		is_active       = COALESCE(:is_active, is_active),
		updated_by      = COALESCE(:updated_by, updated_by),
		updated_at      = NOW(),
		is_del          = COALESCE(:is_del, is_del),

		image_url       = COALESCE(NULLIF(:image_url, ''), image_url),
		unit_id4        = COALESCE(NULLIF(:unit_id4, ''), unit_id4),
		unit_id5        = COALESCE(NULLIF(:unit_id5, ''), unit_id5),

		conv_unit4      = COALESCE(NULLIF(:conv_unit4, 0), conv_unit4),
		conv_unit5      = COALESCE(NULLIF(:conv_unit5, 0), conv_unit5),
		pro_status      = COALESCE(NULLIF(:pro_status, 0), pro_status),

		purch_price1    = COALESCE(NULLIF(:purch_price1, 0), purch_price1),
		purch_price2    = COALESCE(NULLIF(:purch_price2, 0), purch_price2),
		purch_price3    = COALESCE(NULLIF(:purch_price3, 0), purch_price3),
		purch_price4    = COALESCE(NULLIF(:purch_price4, 0), purch_price4),
		purch_price5    = COALESCE(NULLIF(:purch_price5, 0), purch_price5),

		sell_price1     = COALESCE(NULLIF(:sell_price1, 0), sell_price1),
		sell_price2     = COALESCE(NULLIF(:sell_price2, 0), sell_price2),
		sell_price3     = COALESCE(NULLIF(:sell_price3, 0), sell_price3),
		sell_price4     = COALESCE(NULLIF(:sell_price4, 0), sell_price4),
		sell_price5     = COALESCE(NULLIF(:sell_price5, 0), sell_price5),

		cogs            = COALESCE(NULLIF(:cogs, 0), cogs),

		weight1         = COALESCE(NULLIF(:weight1, 0), weight1),
		weight2         = COALESCE(NULLIF(:weight2, 0), weight2),
		weight3         = COALESCE(NULLIF(:weight3, 0), weight3),
		weight4         = COALESCE(NULLIF(:weight4, 0), weight4),
		weight5         = COALESCE(NULLIF(:weight5, 0), weight5),

		length1         = COALESCE(NULLIF(:length1, 0), length1),
		length2         = COALESCE(NULLIF(:length2, 0), length2),
		length3         = COALESCE(NULLIF(:length3, 0), length3),
		length4         = COALESCE(NULLIF(:length4, 0), length4),
		length5         = COALESCE(NULLIF(:length5, 0), length5),

		width1          = COALESCE(NULLIF(:width1, 0), width1),
		width2          = COALESCE(NULLIF(:width2, 0), width2),
		width3          = COALESCE(NULLIF(:width3, 0), width3),
		width4          = COALESCE(NULLIF(:width4, 0), width4),
		width5          = COALESCE(NULLIF(:width5, 0), width5),

		height1         = COALESCE(NULLIF(:height1, 0), height1),
		height2         = COALESCE(NULLIF(:height2, 0), height2),
		height3         = COALESCE(NULLIF(:height3, 0), height3),
		height4         = COALESCE(NULLIF(:height4, 0), height4),
		height5         = COALESCE(NULLIF(:height5, 0), height5),

		volume1         = COALESCE(NULLIF(:volume1, 0), volume1),
		volume2         = COALESCE(NULLIF(:volume2, 0), volume2),
		volume3         = COALESCE(NULLIF(:volume3, 0), volume3),
		volume4         = COALESCE(NULLIF(:volume4, 0), volume4),
		volume5         = COALESCE(NULLIF(:volume5, 0), volume5),

		saf_stock_unit_id = COALESCE(NULLIF(:saf_stock_unit_id, ''), saf_stock_unit_id),
		saf_stock_qty   = COALESCE(NULLIF(:saf_stock_qty, 0), saf_stock_qty),

		min_stock_unit_id = COALESCE(NULLIF(:min_stock_unit_id, ''), min_stock_unit_id),
		min_stock_qty   = COALESCE(NULLIF(:min_stock_qty, 0), min_stock_qty),

		pro_code_coretax = COALESCE(NULLIF(:pro_code_coretax, ''), pro_code_coretax)
    WHERE cust_id = :cust_id AND pro_id = :pro_id;
    `

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Error("productRepository, UpdateImportProduct, begin tx err:", err.Error())
		return err
	}

	_, err = tx.NamedExec(query, p)
	if err != nil {
		_ = tx.Rollback()
		log.Error("productRepository, UpdateImportProduct, err:", err.Error())
		return err
	}

	syncAssignedProductsQuery := `UPDATE mst.m_product child
		SET pro_code = parent.pro_code,
			bar_code = parent.bar_code,
			pro_name = parent.pro_name,
			pcat_id = parent.pcat_id,
			sbrand1_id = parent.sbrand1_id,
			sbrand2_id = parent.sbrand2_id,
			flavor_id = parent.flavor_id,
			ptype_id = parent.ptype_id,
			psize_id = parent.psize_id,
			sup_id = parent.sup_id,
			principal_id = parent.principal_id,
			c_pro_id = parent.c_pro_id,
			is_main_pro = parent.is_main_pro,
			sort_no = parent.sort_no,
			item_no = parent.item_no,
			unit_id1 = parent.unit_id1,
			unit_id2 = parent.unit_id2,
			unit_id3 = parent.unit_id3,
			unit_id4 = parent.unit_id4,
			unit_id5 = parent.unit_id5,
			conv_unit2 = parent.conv_unit2,
			conv_unit3 = parent.conv_unit3,
			conv_unit4 = parent.conv_unit4,
			conv_unit5 = parent.conv_unit5,
			weight = parent.weight,
			is_batch = parent.is_batch,
			is_exp_date = parent.is_exp_date,
			length = parent.length,
			width = parent.width,
			height = parent.height,
			volume = parent.volume,
			purch_price1 = parent.purch_price1,
			sell_price1 = parent.sell_price1,
			purch_price2 = parent.purch_price2,
			sell_price2 = parent.sell_price2,
			purch_price3 = parent.purch_price3,
			sell_price3 = parent.sell_price3,
			purch_price4 = parent.purch_price4,
			sell_price4 = parent.sell_price4,
			purch_price5 = parent.purch_price5,
			sell_price5 = parent.sell_price5,
			weight1 = parent.weight1,
			length1 = parent.length1,
			width1 = parent.width1,
			height1 = parent.height1,
			volume1 = parent.volume1,
			weight2 = parent.weight2,
			length2 = parent.length2,
			width2 = parent.width2,
			height2 = parent.height2,
			volume2 = parent.volume2,
			weight3 = parent.weight3,
			length3 = parent.length3,
			width3 = parent.width3,
			height3 = parent.height3,
			volume3 = parent.volume3,
			weight4 = parent.weight4,
			length4 = parent.length4,
			width4 = parent.width4,
			height4 = parent.height4,
			volume4 = parent.volume4,
			weight5 = parent.weight5,
			length5 = parent.length5,
			width5 = parent.width5,
			height5 = parent.height5,
			volume5 = parent.volume5,
			saf_stock_qty = parent.saf_stock_qty,
			saf_stock_unit_id = parent.saf_stock_unit_id,
			min_stock_qty = parent.min_stock_qty,
			min_stock_unit_id = parent.min_stock_unit_id,
			excise_rate = parent.excise_rate,
			excise_tax = parent.excise_tax,
			is_active = parent.is_active,
			image_url = parent.image_url,
			vat = parent.vat,
			vat_bg = parent.vat_bg,
			vat_lg_purch = parent.vat_lg_purch,
			vat_lg_sell = parent.vat_lg_sell,
			cogs = parent.cogs,
			pro_status = parent.pro_status,
			pro_code_coretax = parent.pro_code_coretax,
			updated_by = parent.updated_by,
			updated_at = CURRENT_TIMESTAMP
		FROM mst.m_product parent
		WHERE parent.pro_id = $1
			AND parent.cust_id = $2
			AND parent.is_del = false
			AND child.parent_pro_id = parent.pro_id
			AND child.is_del = false`

	_, err = tx.Exec(syncAssignedProductsQuery, p.ProId, p.CustId)
	if err != nil {
		_ = tx.Rollback()
		log.Error("productRepository, UpdateImportProduct, sync assigned products err:", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error("productRepository, UpdateImportProduct, commit err:", err.Error())
		return err
	}

	return nil
}

func (r *productRepositoryImpl) InsertProductUpdateTemp(temp entity.ImportProductUpdateTemp) error {
	query := `
		INSERT INTO import.product_update_temp (
			history_id, cust_id,
			brand_code, brand_name,
			pcat_code, pcat_name,
			pl_code, pl_name,
			sbrand1_code, sbrand1_name,
			sbrand2_code, sbrand2_name,
			flavor_code, flavor_name,
			unit_id, unit_name,
			unit_id1, unit_name1,
			unit_id2, unit_name2,
			unit_id3, unit_name3,
			conv_unit2, conv_unit3,
			ptype_code, ptype_name,
			psize_code, psize_name,
			sup_code, sup_name,
			principal_code, principal_name,
			c_pro_code, c_pro_name,
			pro_code_coretax,
			pro_code, pro_name, bar_code,
			cogs, pro_status, is_active,
			length1, length2, length3,
			width1, width2, width3,
			height1, height2, height3,
			weight1, weight2, weight3,
			saf_stock_qty, saf_stock_unit_id,
			min_stock_qty, min_stock_unit_id,
			purch_price1, purch_price2, purch_price3,
			sell_price1, sell_price2, sell_price3,
			is_batch, is_exp_date,
			parent_pro_id,
			excise_rate, excise_tax,
			vat, vat_bg, vat_lg_purch, vat_lg_sell,
			status_insert, error_message, created_at
		) VALUES (
			:history_id, :cust_id,
			:brand_code, :brand_name,
			:pcat_code, :pcat_name,
			:pl_code, :pl_name,
			:sbrand1_code, :sbrand1_name,
			:sbrand2_code, :sbrand2_name,
			:flavor_code, :flavor_name,
			:unit_id, :unit_name,
			:unit_id1, :unit_name1,
			:unit_id2, :unit_name2,
			:unit_id3, :unit_name3,
			:conv_unit2, :conv_unit3,
			:ptype_code, :ptype_name,
			:psize_code, :psize_name,
			:sup_code, :sup_name,
			:principal_code, :principal_name,
			:c_pro_code, :c_pro_name,
			:pro_code_coretax,
			:pro_code, :pro_name, :bar_code,
			:cogs, :pro_status, :is_active,
			:length1, :length2, :length3,
			:width1, :width2, :width3,
			:height1, :height2, :height3,
			:weight1, :weight2, :weight3,
			:saf_stock_qty, :saf_stock_unit_id,
			:min_stock_qty, :min_stock_unit_id,
			:purch_price1, :purch_price2, :purch_price3,
			:sell_price1, :sell_price2, :sell_price3,
			:is_batch, :is_exp_date,
			:parent_pro_id,
			:excise_rate, :excise_tax,
			:vat, :vat_bg, :vat_lg_purch, :vat_lg_sell,
			:status_insert, :error_message, NOW()
		)
	`

	_, err := r.DB.NamedExec(query, temp)
	return err
}

func (r *productRepositoryImpl) GetCtidsProductUpdateByHistory(historyId int64) ([]string, error) {
	query := `
        SELECT ctid::text 
        FROM import.product_update_temp 
        WHERE history_id = $1 
        ORDER BY created_at ASC
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

func (r *productRepositoryImpl) GetCtidsProductInsertByHistory(historyId int64) ([]string, error) {
	query := `
        SELECT ctid::text 
        FROM import.product_temp 
        WHERE history_id = $1 
        ORDER BY created_at ASC
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

func (r *productRepositoryImpl) DeleteProductUpdateTempByCtid(ctid string) error {
	query := `DELETE FROM import.product_update_temp WHERE ctid = $1`
	_, err := r.DB.Exec(query, ctid)
	return err
}

func (r *productRepositoryImpl) DeleteProductTempByCtid(ctid string) error {
	query := `DELETE FROM import.product_temp WHERE ctid = $1`
	_, err := r.DB.Exec(query, ctid)
	return err
}

func (r *productRepositoryImpl) CountProductUpdateTemp(historyId int64) (int, error) {
	var n int
	err := r.DB.Get(&n, `SELECT COUNT(1) FROM import.product_update_temp WHERE history_id = $1`, historyId)
	return n, err
}

func (r *productRepositoryImpl) GetImportProductTotalData(historyId int64) (int, error) {
	var n int
	err := r.DB.Get(&n, `SELECT failed_data FROM import.import_history WHERE history_id = $1`, historyId)
	return n, err
}

func (r *productRepositoryImpl) CountProductTemp(historyId int64) (int, error) {
	var n int
	err := r.DB.Get(&n, `SELECT COUNT(1) FROM import.product_temp WHERE history_id = $1`, historyId)
	return n, err
}

func (repository *productRepositoryImpl) Store(ctx context.Context, product model.Product) (int64, error) {
	if product.DistributorID != nil && *product.DistributorID <= 0 {
		product.DistributorID = nil
	}

	// Try to get transaction from context
	tx := GetTxFromContext(ctx)
	if tx == nil {
		// Fallback to db if no transaction (for non-transactional calls)
		tx = (*sqlx.Tx)(nil)
	}

	query :=
		`INSERT INTO mst.m_product(
			cust_id, pro_code, bar_code, 
			pro_name, pcat_id, 
			sbrand1_id, sbrand2_id, flavor_id, 
			ptype_id, psize_id, sup_id, principal_id, 
			c_pro_id, is_main_pro, sort_no, item_no, 
			unit_id1, unit_id2, unit_id3, unit_id4, unit_id5, 
			conv_unit2, conv_unit3, conv_unit4, conv_unit5, 
			weight, is_batch, is_exp_date, 
			length, width, height, volume, parent_pro_id,
			purch_price1, sell_price1,
			purch_price2, sell_price2,
			purch_price3, sell_price3,
			purch_price4, sell_price4,
			purch_price5, sell_price5,
			weight1, length1, width1, height1, volume1,
			weight2, length2, width2, height2, volume2,
			weight3, length3, width3, height3, volume3,
			weight4, length4, width4, height4, volume4,
			weight5, length5, width5, height5, volume5,
			saf_stock_qty, saf_stock_unit_id, min_stock_qty, min_stock_unit_id,
			excise_rate, excise_tax, is_active, image_url,
			vat, vat_bg, vat_lg_purch, vat_lg_sell, pro_status, cogs,pro_code_coretax,
			created_by, created_at, 
			updated_by, updated_at, 
			is_del, deleted_by, deleted_at,
			distributor_id, level, origin, assigner_user_id, is_product_mapping)
		VALUES ( 
			:cust_id, :pro_code, :bar_code, 
			:pro_name, :pcat_id, 
			:sbrand1_id, :sbrand2_id, :flavor_id, 
			:ptype_id, :psize_id, :sup_id, :principal_id, 
			:c_pro_id, :is_main_pro, :sort_no, :item_no, 
			:unit_id1, :unit_id2, :unit_id3, :unit_id4, :unit_id5, 
			:conv_unit2, :conv_unit3, :conv_unit4, :conv_unit5, 
			:weight, 
			:is_batch, :is_exp_date, 
			:length, :width, :height, :volume, 
			:parent_pro_id,
			:purch_price1, :sell_price1,
			:purch_price2, :sell_price2,
			:purch_price3, :sell_price3,
			:purch_price4, :sell_price4,
			:purch_price5, :sell_price5,
			:weight1, :length1, :width1, :height1, :volume1,
			:weight2, :length2, :width2, :height2, :volume2,
			:weight3, :length3, :width3, :height3, :volume3,
			:weight4, :length4, :width4, :height4, :volume4,
			:weight5, :length5, :width5, :height5, :volume5,
			:saf_stock_qty, :saf_stock_unit_id, :min_stock_qty, :min_stock_unit_id,
			:excise_rate, :excise_tax, :is_active, :image_url,
			:vat, :vat_bg, :vat_lg_purch, :vat_lg_sell, :pro_status, :cogs, :pro_code_coretax,
			:created_by, :created_at, 
			:updated_by, :updated_at, 
			:is_del, :deleted_by, :deleted_at,
			:distributor_id, :level, :origin, :assigner_user_id, :is_product_mapping
		) RETURNING pro_id;`
	// lastInsertId := product.ProductId

	var lastInsertID int64
	var err error
	var rows *sqlx.Rows
	if tx != nil {
		// Execute within transaction
		rows, err = tx.NamedQuery(query, product) // .Scan(&lastInsertId)
	} else {
		// Execute without transaction
		rows, err = repository.NamedQuery(query, product) // .Scan(&lastInsertId)
	}

	if err != nil {
		log.Error("productRepository, Store, err:", err.Error())
		return product.ProductId, err
	}

	// IMPORTANT: Always defer close to avoid connection corruption
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&lastInsertID)
	}
	return lastInsertID, nil
}

func (repository *productRepositoryImpl) BulkStore(products Products) ([]int64, error) {
	query :=
		`INSERT INTO mst.m_product(
			cust_id, pro_code, bar_code, 
			pro_name, pcat_id, 
			sbrand1_id, sbrand2_id, flavor_id, 
			ptype_id, psize_id, sup_id, principal_id, 
			c_pro_id, is_main_pro, sort_no, item_no, 
			unit_id1, unit_id2, unit_id3, unit_id4, unit_id5, 
			conv_unit2, conv_unit3, conv_unit4, conv_unit5, 
			weight, is_batch, is_exp_date, 
			length, width, height, volume, parent_pro_id,
			purch_price1, sell_price1,
			purch_price2, sell_price2,
			purch_price3, sell_price3,
			purch_price4, sell_price4,
			purch_price5, sell_price5,
			weight1, length1, width1, height1, volume1,
			weight2, length2, width2, height2, volume2,
			weight3, length3, width3, height3, volume3,
			weight4, length4, width4, height4, volume4,
			weight5, length5, width5, height5, volume5, 
			saf_stock_qty, saf_stock_unit_id, min_stock_qty, min_stock_unit_id,
			excise_rate, excise_tax, is_active, 
			vat, vat_bg, vat_lg_purch, vat_lg_sell, pro_status, cogs,
			created_by, created_at, 
			updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			:cust_id, :pro_code, :bar_code, 
			:pro_name, :pcat_id, 
			:sbrand1_id, :sbrand2_id, :flavor_id, 
			:ptype_id, :psize_id, :sup_id, :principal_id, 
			:c_pro_id, :is_main_pro, :sort_no, :item_no, 
			:unit_id1, :unit_id2, :unit_id3, :unit_id4, :unit_id5, 
			:conv_unit2, :conv_unit3, :conv_unit4, :conv_unit5, 
			:weight, 
			:is_batch, :is_exp_date, 
			:length, :width, :height, :volume, 
			:parent_pro_id,
			:purch_price1, :sell_price1,
			:purch_price2, :sell_price2,
			:purch_price3, :sell_price3,
			:purch_price4, :sell_price4,
			:purch_price5, :sell_price5,
			:weight1, :length1, :width1, :height1, :volume1,
			:weight2, :length2, :width2, :height2, :volume2,
			:weight3, :length3, :width3, :height3, :volume3,
			:weight4, :length4, :width4, :height4, :volume4,
			:weight5, :length5, :width5, :height5, :volume5, 
			:saf_stock_qty, :saf_stock_unit_id, :min_stock_qty, :min_stock_unit_id,
			:excise_rate, :excise_tax, :is_active, 
			:vat, :vat_bg, :vat_lg_purch, :vat_lg_sell, :pro_status, :cogs,
			:created_by, :created_at, 
			:updated_by, :updated_at, 
			:is_del, :deleted_by, :deleted_at
		) RETURNING pro_id;`
	// lastInsertId := product.ProductId
	// var insertedIDs []int64
	insertedIDs := []int64{}
	rows, err := repository.NamedQuery(query, products) // .Scan(&lastInsertId)
	if err != nil {
		log.Error("productRepository, Bulk Store, err:", err.Error())
		return products.getProductIDs(), err
	}

	for rows.Next() {
		productId := 0
		err := rows.Scan(&productId)
		if err != nil {
			log.Fatal(err)
		}
		insertedIDs = append(insertedIDs, int64(productId))
	}

	return insertedIDs, nil
}

func (products Products) getProductIDs() []int64 {
	var list []int64
	for _, products := range products {
		list = append(list, products.ProductId)
	}
	return list
}

func (repository *productRepositoryImpl) Update(productId int64, request entity.UpdateProductRequest) error {
	var (
		r            model.ProductUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("productRepository, Update, Fields & Ar. gs: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_product
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND pro_id = :pro_id_old;`

	// log.Error("productRepository, Update, query:", query)

	sqlPatch.Args["pro_id_old"] = productId
	sqlPatch.Args["cust_id"] = request.CustId

	tx, err := repository.Beginx()
	if err != nil {
		log.Error("productRepository, Update, begin tx err:", err.Error())
		return err
	}

	result, err := tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		_ = tx.Rollback()
		log.Error("productRepository, Update, err:", err.Error())
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		_ = tx.Rollback()
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		_ = tx.Rollback()
		return errors.New("no rows affected")
	}

	syncAssignedProductsQuery := `UPDATE mst.m_product child
		SET pro_code = parent.pro_code,
			bar_code = parent.bar_code,
			pro_name = parent.pro_name,
			pcat_id = parent.pcat_id,
			sbrand1_id = parent.sbrand1_id,
			sbrand2_id = parent.sbrand2_id,
			flavor_id = parent.flavor_id,
			ptype_id = parent.ptype_id,
			psize_id = parent.psize_id,
			sup_id = parent.sup_id,
			principal_id = parent.principal_id,
			c_pro_id = parent.c_pro_id,
			is_main_pro = parent.is_main_pro,
			sort_no = parent.sort_no,
			item_no = parent.item_no,
			unit_id1 = parent.unit_id1,
			unit_id2 = parent.unit_id2,
			unit_id3 = parent.unit_id3,
			unit_id4 = parent.unit_id4,
			unit_id5 = parent.unit_id5,
			conv_unit2 = parent.conv_unit2,
			conv_unit3 = parent.conv_unit3,
			conv_unit4 = parent.conv_unit4,
			conv_unit5 = parent.conv_unit5,
			weight = parent.weight,
			is_batch = parent.is_batch,
			is_exp_date = parent.is_exp_date,
			length = parent.length,
			width = parent.width,
			height = parent.height,
			volume = parent.volume,
			purch_price1 = parent.purch_price1,
			sell_price1 = parent.sell_price1,
			purch_price2 = parent.purch_price2,
			sell_price2 = parent.sell_price2,
			purch_price3 = parent.purch_price3,
			sell_price3 = parent.sell_price3,
			purch_price4 = parent.purch_price4,
			sell_price4 = parent.sell_price4,
			purch_price5 = parent.purch_price5,
			sell_price5 = parent.sell_price5,
			weight1 = parent.weight1,
			length1 = parent.length1,
			width1 = parent.width1,
			height1 = parent.height1,
			volume1 = parent.volume1,
			weight2 = parent.weight2,
			length2 = parent.length2,
			width2 = parent.width2,
			height2 = parent.height2,
			volume2 = parent.volume2,
			weight3 = parent.weight3,
			length3 = parent.length3,
			width3 = parent.width3,
			height3 = parent.height3,
			volume3 = parent.volume3,
			weight4 = parent.weight4,
			length4 = parent.length4,
			width4 = parent.width4,
			height4 = parent.height4,
			volume4 = parent.volume4,
			weight5 = parent.weight5,
			length5 = parent.length5,
			width5 = parent.width5,
			height5 = parent.height5,
			volume5 = parent.volume5,
			saf_stock_qty = parent.saf_stock_qty,
			saf_stock_unit_id = parent.saf_stock_unit_id,
			min_stock_qty = parent.min_stock_qty,
			min_stock_unit_id = parent.min_stock_unit_id,
			excise_rate = parent.excise_rate,
			excise_tax = parent.excise_tax,
			is_active = parent.is_active,
			image_url = parent.image_url,
			vat = parent.vat,
			vat_bg = parent.vat_bg,
			vat_lg_purch = parent.vat_lg_purch,
			vat_lg_sell = parent.vat_lg_sell,
			cogs = parent.cogs,
			pro_status = parent.pro_status,
			pro_code_coretax = parent.pro_code_coretax,
			updated_by = parent.updated_by,
			updated_at = CURRENT_TIMESTAMP
		FROM mst.m_product parent
		WHERE parent.pro_id = :parent_pro_id
			AND parent.cust_id = :parent_cust_id
			AND parent.is_del = false
			AND child.parent_pro_id = parent.pro_id
			AND child.is_del = false`

	_, err = tx.NamedExec(syncAssignedProductsQuery, map[string]interface{}{
		"parent_pro_id":  productId,
		"parent_cust_id": request.CustId,
	})
	if err != nil {
		_ = tx.Rollback()
		log.Error("productRepository, Update, sync assigned products err:", err.Error())
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error("productRepository, Update, commit err:", err.Error())
		return err
	}

	return nil
}

func (repository *productRepositoryImpl) Delete(custId string, productId int64, deletedBy int64) error {
	return repository.DeleteWithContext(context.Background(), custId, productId, deletedBy)
}

func (repository *productRepositoryImpl) DeleteWithContext(ctx context.Context, custId string, productId int64, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_product
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND pro_id = :pro_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"pro_id":     productId,
		"deleted_by": deletedBy,
	}

	tx := GetTxFromContext(ctx)

	var (
		result sql.Result
		err    error
	)
	if tx != nil {
		result, err = tx.NamedExec(query, wMap)
	} else {
		result, err = repository.NamedExec(query, wMap)
	}
	if err != nil {
		log.Error("ProductRepository, DeleteWithContext, err:", err.Error())
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

func (repository *productRepositoryImpl) DeleteMultiple(custId string, productId []int64, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_product
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND pro_id IN (:pro_id);`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"pro_id":     strings.Trim(strings.Replace(fmt.Sprint(productId), " ", ", ", -1), "[]"),
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Error("ProductRepository, DeleteMultiple, err:", err.Error())
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

func (repository *productRepositoryImpl) StoreDist(productsDist []model.ProductDistCreate) error {
	query :=
		`INSERT INTO mst.m_product_dist(
			cust_id,pro_id,is_active,is_alloc,min_stock,
            min_stock_str,safety_stock,safety_stock_str,po_formula,is_new_pro,vat,vat_bg,
            vat_lg_purch,vat_lg_sell,s_mweek1,s_mweek2,cogs,updated_by,updated_at,parent_pro_id)
		VALUES ( 
			:cust_id, :pro_id,:is_active,:is_alloc,:min_stock,:min_stock_str,:safety_stock,
		    :safety_stock_str,:po_formula,:is_new_pro,:vat,:vat_bg,:vat_lg_purch,:vat_lg_sell,
		    :s_mweek1,:s_mweek2,:cogs,:updated_by,:updated_at,:parent_pro_id
		);`
	// lastInsertId := product.ProductId
	_, err := repository.NamedExec(query, productsDist) // .Scan(&lastInsertId)
	if err != nil {
		log.Error("productDistRepository, Store, err:", err.Error())
		return err
	}
	return nil
}

func (repository *productRepositoryImpl) FindAllPrincipal(dataFilter entity.ProductPrincipalQueryFilter, custId string) ([]model.Principal, int, int, error) {

	principals := []model.Principal{}
	selectField := ` pro.principal_id, pri.principal_code, pri.principal_name `
	qWhere := ` LEFT JOIN mst.m_principal pri on pri.principal_id = pro.principal_id and pri.cust_id = pro.cust_id
				WHERE pro.cust_id = '` + dataFilter.ParentCustId + `' AND pro.is_active = true AND pro.deleted_at IS NULL `

	if dataFilter.Query != "" {
		qWhere += ` AND (pri.principal_code ILIKE '%` + dataFilter.Query + `%' 
					OR pri.principal_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	queryGroup := ` GROUP BY pro.principal_id, pri.principal_code, pri.principal_name `
	queryFrom := ` FROM mst.m_product pro `
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere + ` ` + queryGroup

	// log.Error("productRepository, queryCount:", queryCount)
	var total int

	sortBy := ` pro.principal_id`
	querySelect += fmt.Sprintf(`ORDER BY %s ASC`, sortBy)

	// log.Error("productRepository, querySelect:", querySelect)
	err := repository.Select(&principals, querySelect)
	if err != nil {
		log.Error("productRepository, FindAllPrincipal, err:", err.Error())
		return principals, total, 1, err
	}

	total = len(principals)

	return principals, total, 1, nil
}

func (repository *productRepositoryImpl) FindAllCategory(dataFilter entity.ProductCategoryQueryFilter, custId string) ([]model.ProductCat, int, int, error) {

	productCategories := []model.ProductCat{}
	selectField := ` pro.pcat_id, pcat.pcat_code, pcat.pcat_name  `
	qWhere := ` LEFT JOIN mst.m_product_cat pcat ON pcat.pcat_id = pro.pcat_id AND pcat.cust_id = pro.cust_id 
				WHERE pro.cust_id = '` + dataFilter.ParentCustId + `' AND pro.is_active = true AND pro.deleted_at IS NULL `

	if dataFilter.Query != "" {
		qWhere += ` AND (pcat.pcat_code ILIKE '%` + dataFilter.Query + `%' 
					OR pcat.pcat_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	queryGroup := ` GROUP BY pro.pcat_id, pcat.pcat_code, pcat.pcat_name  `
	queryFrom := ` FROM mst.m_product pro `
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere + ` ` + queryGroup

	// log.Error("productRepository, queryCount:", queryCount)
	var total int

	sortBy := ` pro.pcat_id`
	querySelect += fmt.Sprintf(`ORDER BY %s ASC`, sortBy)

	// log.Error("productRepository, querySelect:", querySelect)
	err := repository.Select(&productCategories, querySelect)
	if err != nil {
		log.Error("productRepository, FindAllCategory, err:", err.Error())
		return productCategories, total, 1, err
	}

	total = len(productCategories)

	return productCategories, total, 1, nil
}

func (repository *productRepositoryImpl) FindAllBrand(dataFilter entity.ProductBrandQueryFilter, custId string) ([]model.Brand, int, int, error) {

	brands := []model.Brand{}
	selectField := ` br.brand_id, br.brand_code, br.brand_name  `
	qWhere := ` LEFT JOIN mst.m_sub_brand1 sbr ON sbr.sbrand1_id = pro.sbrand1_id AND sbr.cust_id = pro.cust_id
				LEFT JOIN mst.m_brand br ON br.brand_id = sbr.brand_id AND sbr.cust_id = sbr.cust_id 
				WHERE pro.cust_id = '` + dataFilter.ParentCustId + `' AND pro.is_active = true AND pro.deleted_at IS NULL `

	if dataFilter.Query != "" {
		qWhere += ` AND (br.brand_code ILIKE '%` + dataFilter.Query + `%' 
					OR br.brand_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	queryGroup := ` GROUP BY br.brand_id, br.brand_code, br.brand_name  `
	queryFrom := ` FROM mst.m_product pro `
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere + ` ` + queryGroup

	// log.Error("productRepository, queryCount:", queryCount)
	var total int

	sortBy := ` br.brand_id`
	querySelect += fmt.Sprintf(`ORDER BY %s ASC`, sortBy)

	// log.Error("productRepository, querySelect:", querySelect)
	err := repository.Select(&brands, querySelect)
	if err != nil {
		log.Error("productRepository, FindAllBrand, err:", err.Error())
		return brands, total, 1, err
	}

	total = len(brands)

	return brands, total, 1, nil
}
