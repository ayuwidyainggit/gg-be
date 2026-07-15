package repository

import (
	"context"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/constant"
	"math"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	ErrNoRowsAffected   = "no rows affected"
	ErrNoFieldsToUpdate = "no fields to update"
)

type (
	RepositoryReplenishmentImpl struct {
		*gorm.DB
	}
)

type ReplenishmentRepository interface {
	FindAllByCustId(ctx context.Context, dataFilter entity.ReplenishmentQueryFilter, custId, parentCustId string, isPrincipal bool) ([]model.ReplenishmentOrderList, int64, int, error)
	Store(ctx context.Context, replenishment *model.ReplenishmentOrder) error
	CreateDetail(ctx context.Context, detail *model.ReplenishmentOrderDetail) error
	FindSupplierByID(ctx context.Context, supID int64, custID string) error
	FindWarehouseByID(ctx context.Context, whID int64, distributorID int64) error
	FindProductByID(ctx context.Context, proID int64, custID, parentCustID string) error
	FindByReplenishmentID(ctx context.Context, replenishmentID string, custID, parentCustID string, isPrincipal bool) (*model.ReplenishmentOrderRead, error)
	FindDetailByReplenishmentID(ctx context.Context, replenishmentID string, custID, parentCustID string, isPrincipal bool) ([]model.ReplenishmentOrderDetailRead, error)
	FindFinalByReplenishmentID(ctx context.Context, replenishmentNo string, custID, parentCustID string, isPrincipal bool) ([]model.ReplenishmentFinalRead, error)
	FindGoodReceiptByReplenishmentNo(ctx context.Context, replenishmentNo string, custID, parentCustID string, isPrincipal bool) ([]model.ReplenishmentGoodReceiptRead, error)
	FindProductList(ctx context.Context, dataFilter entity.ReplenishmentProductQueryFilter, custId, parentCustId string) ([]model.ReplenishmentProductList, int64, int, error)
	FindProductGrList(ctx context.Context, dataFilter entity.ProductGrListQueryFilter, custId, parentCustId string) (*model.ReplenishmentOrderRead, []model.ProductGrListDetail, int64, int, error)
	FindPoList(ctx context.Context, dataFilter entity.PoListQueryFilter, custId, parentCustId string) ([]model.PoList, int64, int, error)
	FindApprovalProducts(ctx context.Context, dataFilter entity.ReplenishmentApprovalProductQueryFilter, custId, parentCustId string, isPrincipal bool) ([]model.ReplenishmentApprovalProduct, int64, int, error)
	FindApprovalOrderList(ctx context.Context, dataFilter entity.ReplenishmentApprovalListQueryFilter, custId, parentCustId string, isPrincipal bool) ([]model.ReplenishmentOrderList, int64, int, error)
	FindByReplenishmentNo(ctx context.Context, replenishmentNo string, custID string) (*model.ReplenishmentOrder, error)
	GetReplenishmentOrderByID(ctx context.Context, replenishmentID int64, custID string) (*model.ReplenishmentOrder, error)
	FindDetailByReplenishmentIDForUpdate(ctx context.Context, replenishmentID int64, custID string) ([]model.ReplenishmentOrderDetail, error)
	UpdateApproval(ctx context.Context, replenishment *model.ReplenishmentOrder) error
	UpdateDetail(ctx context.Context, detail *model.ReplenishmentOrderDetail) error
	CreateDetailForApproval(ctx context.Context, detail *model.ReplenishmentOrderDetail) error
	SoftDeleteDetail(ctx context.Context, detailID int64, custID string, deletedBy int64) error
	CheckIsPrincipal(ctx context.Context, custID string) (bool, error)
	UpdateStatusByReplenishmentNo(ctx context.Context, replenishmentNo string, custID string, status int, updatedBy int64) error
	LookupProductByProCode(ctx context.Context, custID string, proCode string) (proID int64, unitID3 *string, err error)
	FindActiveReplenishmentDetail(ctx context.Context, custID string, replenishmentID, proID int64) (*model.ReplenishmentOrderDetail, error)
	UpdateReplenishmentDetailSAPFields(ctx context.Context, custID string, replenishmentDetailID int64, sapQty3, sapPurchPrice3 *float64, updatedBy int64) error
	FindCustIDByDistributorCode(ctx context.Context, distributorCode int) (string, error)
	GetDistributorIDByCustID(ctx context.Context, custID string) (int64, error)
	IsDistributorApprovalPIC(ctx context.Context, userID int64, supID int64, distributorID int64) (bool, error)
	GetDistributorApprovalRequirement(ctx context.Context, supID int64, distributorID int64) (bool, bool, error)
	GetInitialDistributorApproval(ctx context.Context, supID int64, distributorID int64) (int, int, int64, error)
	InsertReplenishmentOrderApproval(ctx context.Context, custID string, replenishmentID int64, level int, sequence int, pic int64) error
	IsUserApprovalPIC(ctx context.Context, userID int64) (bool, error)
	UpdateReplenishmentOrderApprovalStatus(ctx context.Context, custID string, replenishmentID int64, pic int64, status int, remarks *string) error
	IsReturnReasonDistributorExists(ctx context.Context, returnReasonID int64) (bool, error)
	FindSummarizeReplanishment(ctx context.Context, replanishmentIDs []int64, custID, parentCustID string, isPrincipal bool, userID int64) ([]model.SummarizeReplanishmentRow, error)
	FindSAPExportRows(ctx context.Context, distributorCode string, dateFrom, dateTo int64) ([]model.SAPReplExportRow, error)
}

func NewReplenishmentRepo(db *gorm.DB) *RepositoryReplenishmentImpl {
	return &RepositoryReplenishmentImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryReplenishmentImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repo *RepositoryReplenishmentImpl) buildListQueryBase(ctx context.Context, dataFilter entity.ReplenishmentQueryFilter, custId, parentCustId string, isPrincipal bool, forCount bool) *gorm.DB {
	_ = parentCustId
	_ = isPrincipal

	hasDistributorColumn := false
	_ = repo.model(ctx).Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_schema = 'inv'
			  AND table_name = 'replenishment_order'
			  AND column_name = 'distributor_id'
		)
	`).Scan(&hasDistributorColumn).Error

	distributorExpr := "mc.distributor_id"
	if hasDistributorColumn {
		distributorExpr = "COALESCE(replenishment_order.distributor_id, mc.distributor_id)"
	}

	var query *gorm.DB
	if forCount {
		query = repo.model(ctx).Select("replenishment_order.replenishment_id")
	} else {
		query = repo.model(ctx).
			Select(fmt.Sprintf("replenishment_order.*, us.user_name AS created_by_name, us2.user_name AS updated_by_name, sup.sup_code, sup.sup_name, wh.wh_code, wh.wh_name, rs.status_name, %s AS distributor_id, md.distributor_code, md.distributor_name, md.address", distributorExpr))
	}

	distributorJoin := fmt.Sprintf(
		"LEFT JOIN mst.m_distributor md ON ((%s > 0 AND md.distributor_id = %s) OR ((%s IS NULL OR %s = 0) AND md.cust_id = replenishment_order.cust_id))",
		distributorExpr, distributorExpr, distributorExpr, distributorExpr,
	)

	query = query.
		Table("inv.replenishment_order").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = replenishment_order.created_by").
		Joins("LEFT JOIN sys.m_user us2 ON us2.user_id = replenishment_order.updated_by").
		Joins("LEFT JOIN smc.m_customer mc ON mc.cust_id = replenishment_order.cust_id").
		Joins(distributorJoin).
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = replenishment_order.sup_id").
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = replenishment_order.wh_id AND wh.cust_id = replenishment_order.cust_id").
		Joins("LEFT JOIN inv.replenishment_status rs ON rs.status_code = replenishment_order.status").
		Where("replenishment_order.is_del = false")

	if dataFilter.IsPICUser {
		approvalExists := `EXISTS (
				SELECT 1
				FROM inv.replenishment_order_approval roa
				WHERE roa.cust_id = replenishment_order.cust_id
				  AND roa.replenishment_order_id = replenishment_order.replenishment_id
				  AND (roa.pic = ? OR roa.pic = ?)
			)`
		if dataFilter.EmpID > 0 {
			query = query.Where(
				`(`+approvalExists+` OR replenishment_order.created_by = ?)`,
				dataFilter.UserID, dataFilter.EmpID, dataFilter.UserID,
			)
		} else {
			query = query.Where(
				`EXISTS (
				SELECT 1
				FROM inv.replenishment_order_approval roa
				WHERE roa.cust_id = replenishment_order.cust_id
				  AND roa.replenishment_order_id = replenishment_order.replenishment_id
				  AND roa.pic = ?
			) OR replenishment_order.created_by = ?`,
				dataFilter.UserID, dataFilter.UserID,
			)
		}
	}

	if !dataFilter.IsPICUser {
		query = query.Where("replenishment_order.cust_id = ?", custId)
	} else if dataFilter.DistributorIDFromToken != nil && *dataFilter.DistributorIDFromToken > 0 {
		query = query.Where(fmt.Sprintf("%s = ?", distributorExpr), *dataFilter.DistributorIDFromToken)
	} else if isPrincipal {
		query = query.Where("(md.parent_cust_id = ? OR replenishment_order.cust_id = ?)", custId, custId)
	} else {
		query = query.Where("replenishment_order.cust_id = ?", custId)
	}

	if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
		startTime := time.Unix(*dataFilter.StartDate, 0).In(constant.AsiaJakartaLocation)
		endTime := time.Unix(*dataFilter.EndDate, 0).In(constant.AsiaJakartaLocation)
		endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, constant.AsiaJakartaLocation)
		query = query.Where("replenishment_order.date BETWEEN ? AND ?", startTime, endTime)
	}

	if dataFilter.StartDeliveryDate != nil && dataFilter.EndDeliveryDate != nil {
		startTime := time.Unix(*dataFilter.StartDeliveryDate, 0).In(constant.AsiaJakartaLocation)
		endTime := time.Unix(*dataFilter.EndDeliveryDate, 0).In(constant.AsiaJakartaLocation)
		endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, constant.AsiaJakartaLocation)
		query = query.Where("replenishment_order.delivery_date BETWEEN ? AND ?", startTime, endTime)
	}

	if len(dataFilter.SupIDParsed) > 0 {
		query = query.Where("replenishment_order.sup_id IN ?", dataFilter.SupIDParsed)
	}

	if len(dataFilter.StatusParsed) > 0 {
		query = query.Where("replenishment_order.status IN ?", dataFilter.StatusParsed)
	}

	if dataFilter.Distributor != nil {
		query = query.Where(fmt.Sprintf("%s = ?", distributorExpr), *dataFilter.Distributor)
	}

	if dataFilter.PoNo != "" {
		query = query.Where("replenishment_order.replenishment_no ILIKE ?", "%"+dataFilter.PoNo+"%")
	}

	if dataFilter.SoNo != "" {
		query = query.Where("replenishment_order.so_no ILIKE ?", "%"+dataFilter.SoNo+"%")
	}

	return query
}

func (repo *RepositoryReplenishmentImpl) FindAllByCustId(ctx context.Context, dataFilter entity.ReplenishmentQueryFilter, custId, parentCustId string, isPrincipal bool) ([]model.ReplenishmentOrderList, int64, int, error) {
	var replenishments []model.ReplenishmentOrderList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = constant.DEFAULT_PAGE_LIMIT
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repo.buildListQueryBase(ctx, dataFilter, custId, parentCustId, isPrincipal, true)
	query := repo.buildListQueryBase(ctx, dataFilter, custId, parentCustId, isPrincipal, false)

	// Sorting
	sortBy := ""
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				// Handle default sort field mapping
				sortField := colSort[0]
				if sortField == "created_date" {
					sortField = "created_at"
				}
				sortBy += fmt.Sprintf(`replenishment_order.%s %s, `, sortField, colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		if sortBy != "" {
			query.Order(sortBy)
		}
	}
	// Default sort
	if sortBy == "" {
		query.Order("replenishment_order.created_at DESC")
	}

	// Pagination
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	err := query.Limit(limit).Offset(offset).Find(&replenishments).Error
	if err != nil {
		return replenishments, total, 0, err
	}

	err = queryCount.Model(&model.ReplenishmentOrderList{}).Count(&total).Error
	if err != nil {
		return replenishments, total, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return replenishments, total, lastPage, nil
}

func (repo *RepositoryReplenishmentImpl) buildApprovalListQuery(ctx context.Context, dataFilter entity.ReplenishmentApprovalListQueryFilter, custId, parentCustId string, forCount bool) *gorm.DB {
	_ = parentCustId
	var query *gorm.DB
	if forCount {
		query = repo.model(ctx).Table("inv.replenishment_order AS ro").Select("count(DISTINCT ro.replenishment_id)")
	} else {
		query = repo.model(ctx).Table("inv.replenishment_order AS ro").
			Select(`ro.replenishment_id, ro.cust_id, ro.replenishment_no, ro.date, ro.delivery_date, ro.sup_id, ro.wh_id, ro.status, ro.created_by, ro.created_at, ro.updated_by, ro.updated_at, ro.is_del,
				us.user_name AS created_by_name, us2.user_name AS updated_by_name,
				sup.sup_code, sup.sup_name, mc.distributor_id, md.distributor_code, md.distributor_name, md.address, wh.wh_code, wh.wh_name, rs.status_name`)
	}

	query = query.
		Joins("INNER JOIN inv.replenishment_order_approval roa ON roa.cust_id = ro.cust_id AND roa.replenishment_order_id = ro.replenishment_id").
		Joins("LEFT JOIN sys.m_user us ON us.user_id = ro.created_by").
		Joins("LEFT JOIN sys.m_user us2 ON us2.user_id = ro.updated_by").
		Joins("LEFT JOIN smc.m_customer mc ON mc.cust_id = ro.cust_id").
		Joins("LEFT JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = ro.sup_id").
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = ro.wh_id AND wh.cust_id = ?", custId).
		Joins("LEFT JOIN inv.replenishment_status rs ON rs.status_code = ro.status").
		Where("ro.is_del = false")

	if dataFilter.EmpID > 0 {
		query = query.Where("(roa.pic = ? OR roa.pic = ?)", dataFilter.UserID, dataFilter.EmpID)
	} else {
		query = query.Where("roa.pic = ?", dataFilter.UserID)
	}

	if dataFilter.DistributorIDFromToken != nil && *dataFilter.DistributorIDFromToken > 0 {
		query = query.Where("mc.distributor_id = ?", *dataFilter.DistributorIDFromToken)
	} else {
		query = query.Where("ro.cust_id = ?", custId)
	}

	statusIn := dataFilter.StatusParsed
	if len(statusIn) == 0 {
		statusIn = []int{1, 2, 3}
	}
	query = query.Where("roa.status IN ?", statusIn)

	if dataFilter.Q != "" {
		query = query.Where("ro.replenishment_no ILIKE ?", "%"+dataFilter.Q+"%")
	}

	if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
		startTime := time.Unix(*dataFilter.StartDate, 0).In(constant.AsiaJakartaLocation)
		endTime := time.Unix(*dataFilter.EndDate, 0).In(constant.AsiaJakartaLocation)
		endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, constant.AsiaJakartaLocation)
		query = query.Where("ro.date BETWEEN ? AND ?", startTime, endTime)
	}

	if dataFilter.DeliveryDateStart != nil && dataFilter.DeliveryDateEnd != nil {
		startTime := time.Unix(*dataFilter.DeliveryDateStart, 0).In(constant.AsiaJakartaLocation)
		endTime := time.Unix(*dataFilter.DeliveryDateEnd, 0).In(constant.AsiaJakartaLocation)
		endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, constant.AsiaJakartaLocation)
		query = query.Where("ro.delivery_date BETWEEN ? AND ?", startTime, endTime)
	}

	return query
}

func (repo *RepositoryReplenishmentImpl) FindApprovalOrderList(ctx context.Context, dataFilter entity.ReplenishmentApprovalListQueryFilter, custId, parentCustId string, isPrincipal bool) ([]model.ReplenishmentOrderList, int64, int, error) {
	_ = isPrincipal
	var rows []model.ReplenishmentOrderList
	var total int64

	limit := dataFilter.Limit
	if limit == 0 {
		limit = constant.DEFAULT_PAGE_LIMIT
	}
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	queryCount := repo.buildApprovalListQuery(ctx, dataFilter, custId, parentCustId, true)
	if err := queryCount.Scan(&total).Error; err != nil {
		return rows, 0, 0, err
	}

	query := repo.buildApprovalListQuery(ctx, dataFilter, custId, parentCustId, false)
	err := query.Order("ro.created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error
	if err != nil {
		return rows, total, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return rows, total, lastPage, nil
}

func (repo *RepositoryReplenishmentImpl) Store(ctx context.Context, replenishment *model.ReplenishmentOrder) error {
	db := repo.model(ctx)
	err := replenishment.BeforeCreate(db)
	if err != nil {
		return fmt.Errorf("failed to generate replenishment_no: %w", err)
	}

	var nextID struct {
		ReplenishmentID int64 `gorm:"column:replenishment_id"`
	}

	// Try to get ID from sequence first, fallback to MAX+1 if sequence doesn't exist
	err = repo.model(ctx).Raw(`
		SELECT COALESCE(
			nextval('inv.replenishment_order_replenishment_id_seq'),
			(SELECT COALESCE(MAX(replenishment_id), 0) + 1 FROM inv.replenishment_order WHERE cust_id = $1)
		) AS replenishment_id
	`, replenishment.CustID).Scan(&nextID).Error

	if err != nil {
		// If sequence doesn't exist, use MAX+1 approach
		err = repo.model(ctx).Raw(`
			SELECT COALESCE(MAX(replenishment_id), 0) + 1 AS replenishment_id
			FROM inv.replenishment_order
			WHERE cust_id = $1
		`, replenishment.CustID).Scan(&nextID).Error
		if err != nil {
			return fmt.Errorf("failed to get next replenishment_id: %w", err)
		}
	}

	if nextID.ReplenishmentID == 0 {
		nextID.ReplenishmentID = 1
	}

	var hasDistributorColumn bool
	err = repo.model(ctx).Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_schema = 'inv'
			  AND table_name = 'replenishment_order'
			  AND column_name = 'distributor_id'
		)
	`).Scan(&hasDistributorColumn).Error
	if err != nil {
		return fmt.Errorf("failed to detect distributor_id column: %w", err)
	}

	const maxInsertAttempt = 5
	for attempt := 0; attempt < maxInsertAttempt; attempt++ {
		if attempt > 0 {
			if regenErr := replenishment.BeforeCreate(db); regenErr != nil {
				return fmt.Errorf("failed to regenerate replenishment_no: %w", regenErr)
			}
		}

		if hasDistributorColumn {
			err = db.Exec(`
				INSERT INTO inv.replenishment_order 
				(cust_id, replenishment_id, replenishment_no, date, distributor_id, sup_id, wh_id, delivery_type, replenishment_type, 
				 so_start_date, so_end_date, delivery_date, note, status, so_no, created_by, 
				 created_at, updated_by, updated_at, deleted_by, deleted_at, is_del)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
			`,
				replenishment.CustID,
				nextID.ReplenishmentID,
				replenishment.ReplenishmentNo,
				replenishment.Date,
				replenishment.DistributorID,
				replenishment.SupID,
				replenishment.WhID,
				replenishment.DeliveryType,
				replenishment.ReplenishmentType,
				replenishment.SoStartDate,
				replenishment.SoEndDate,
				replenishment.DeliveryDate,
				replenishment.Note,
				replenishment.Status,
				replenishment.SoNo,
				replenishment.CreatedBy,
				replenishment.CreatedAt,
				replenishment.UpdatedBy,
				replenishment.UpdatedAt,
				replenishment.DeletedBy,
				replenishment.DeletedAt,
				replenishment.IsDel,
			).Error
		} else {
			err = db.Exec(`
				INSERT INTO inv.replenishment_order 
				(cust_id, replenishment_id, replenishment_no, date, sup_id, wh_id, delivery_type, replenishment_type, 
				 so_start_date, so_end_date, delivery_date, note, status, so_no, created_by, 
				 created_at, updated_by, updated_at, deleted_by, deleted_at, is_del)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
			`,
				replenishment.CustID,
				nextID.ReplenishmentID,
				replenishment.ReplenishmentNo,
				replenishment.Date,
				replenishment.SupID,
				replenishment.WhID,
				replenishment.DeliveryType,
				replenishment.ReplenishmentType,
				replenishment.SoStartDate,
				replenishment.SoEndDate,
				replenishment.DeliveryDate,
				replenishment.Note,
				replenishment.Status,
				replenishment.SoNo,
				replenishment.CreatedBy,
				replenishment.CreatedAt,
				replenishment.UpdatedBy,
				replenishment.UpdatedAt,
				replenishment.DeletedBy,
				replenishment.DeletedAt,
				replenishment.IsDel,
			).Error
		}

		if err == nil {
			break
		}
		if !isDuplicateReplenishmentNoError(err) {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("failed to insert replenishment order: %w", err)
	}
	replenishment.ReplenishmentID = nextID.ReplenishmentID

	return nil
}

func isDuplicateReplenishmentNoError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "sqlstate 23505") &&
		(strings.Contains(errMsg, "ux_replenishment_order_no_per_cust") ||
			strings.Contains(errMsg, "replenishment_no"))
}

func (repo *RepositoryReplenishmentImpl) CreateDetail(ctx context.Context, detail *model.ReplenishmentOrderDetail) error {
	err := repo.model(ctx).Create(detail).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *RepositoryReplenishmentImpl) FindSupplierByID(ctx context.Context, supID int64, custID string) error {
	var count int64
	err := repo.model(ctx).
		Table("mst.m_supplier").
		Where("sup_id = ? AND is_del = false", supID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("supplier with id %d not found", supID)
	}
	return nil
}

func (repo *RepositoryReplenishmentImpl) FindWarehouseByID(ctx context.Context, whID int64, distributorID int64) error {
	if distributorID <= 0 {
		return fmt.Errorf("invalid distributor_id")
	}

	var count int64
	err := repo.model(ctx).
		Table("mst.m_warehouse mw").
		Joins("JOIN smc.m_customer mc ON mc.cust_id = mw.cust_id").
		Where("mw.wh_id = ? AND mc.distributor_id = ? AND mw.is_active = true AND mw.is_del = false", whID, distributorID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("warehouse with id %d not found", whID)
	}
	return nil
}

func (repo *RepositoryReplenishmentImpl) FindProductByID(ctx context.Context, proID int64, custID, parentCustID string) error {
	var count int64
	err := repo.model(ctx).
		Table("mst.m_product").
		Where("pro_id = ?", proID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("product with id %d not found", proID)
	}
	return nil
}

func (repo *RepositoryReplenishmentImpl) FindByReplenishmentID(ctx context.Context, replenishmentID string, custID, parentCustID string, isPrincipal bool) (*model.ReplenishmentOrderRead, error) {
	var replenishment model.ReplenishmentOrderRead

	replenishmentIDInt, err := strconv.ParseInt(replenishmentID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid replenishment_id: %w", err)
	}

	query := repo.model(ctx).
		Select(`replenishment_order.*, 
			sup.sup_code, 
			sup.sup_name, 
			wh.wh_code, 
			wh.wh_name, 
			rs.status_name,
			md.distributor_id,
			md.distributor_code,
			md.distributor_name,
			md.address,
			gr.delivery_fee`).
		Table("inv.replenishment_order").
		Joins("LEFT JOIN inv.replenishment_status rs ON rs.status_code = replenishment_order.status").
		Joins("LEFT JOIN smc.m_customer mc ON mc.cust_id = replenishment_order.cust_id").
		Joins("LEFT JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id").
		Where("replenishment_order.replenishment_id = ?", replenishmentIDInt)

	if isPrincipal {
		query = query.
			Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = replenishment_order.sup_id").
			Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = replenishment_order.wh_id AND wh.cust_id = replenishment_order.cust_id").
			Joins("LEFT JOIN inv.gr gr ON gr.po_no = replenishment_order.replenishment_no AND gr.cust_id = replenishment_order.cust_id")
	} else {
		query = query.
			Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = replenishment_order.sup_id").
			Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = replenishment_order.wh_id AND wh.cust_id = replenishment_order.cust_id").
			Joins("LEFT JOIN inv.gr gr ON gr.po_no = replenishment_order.replenishment_no AND gr.cust_id = replenishment_order.cust_id")
	}

	err = query.First(&replenishment).Error
	if err != nil {
		return nil, err
	}
	return &replenishment, nil
}

func (repo *RepositoryReplenishmentImpl) FindDetailByReplenishmentID(ctx context.Context, replenishmentID string, custID, parentCustID string, isPrincipal bool) ([]model.ReplenishmentOrderDetailRead, error) {
	var details []model.ReplenishmentOrderDetailRead

	replenishmentIDInt, err := strconv.ParseInt(replenishmentID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid replenishment_id: %w", err)
	}

	// Build subquery for in_transit_stock
	// Use rod.cust_id (from the replenishment order detail being viewed) instead of custID from JWT
	// This ensures we calculate in_transit_stock for the correct customer
	var inTransitStock1Subquery, inTransitStock2Subquery, inTransitStock3Subquery string
	if isPrincipal {
		inTransitStock1Subquery = `COALESCE((
			SELECT SUM(COALESCE(rod_in_transit.order_booking_qty1, 0))
			FROM inv.replenishment_order_detail rod_in_transit
			JOIN inv.replenishment_order ro_in_transit ON ro_in_transit.replenishment_id = rod_in_transit.replenishment_id AND ro_in_transit.cust_id = rod_in_transit.cust_id
			WHERE rod_in_transit.pro_id = rod.pro_id 
				AND rod_in_transit.cust_id = rod.cust_id
				AND ro_in_transit.status = 4
				AND ro_in_transit.is_del = false
				AND rod_in_transit.is_del = false
		), 0)`
		inTransitStock2Subquery = `COALESCE((
			SELECT SUM(COALESCE(rod_in_transit.order_booking_qty2, 0))
			FROM inv.replenishment_order_detail rod_in_transit
			JOIN inv.replenishment_order ro_in_transit ON ro_in_transit.replenishment_id = rod_in_transit.replenishment_id AND ro_in_transit.cust_id = rod_in_transit.cust_id
			WHERE rod_in_transit.pro_id = rod.pro_id 
				AND rod_in_transit.cust_id = rod.cust_id
				AND ro_in_transit.status = 4
				AND ro_in_transit.is_del = false
				AND rod_in_transit.is_del = false
		), 0)`
		inTransitStock3Subquery = `COALESCE((
			SELECT SUM(COALESCE(rod_in_transit.order_booking_qty3, 0))
			FROM inv.replenishment_order_detail rod_in_transit
			JOIN inv.replenishment_order ro_in_transit ON ro_in_transit.replenishment_id = rod_in_transit.replenishment_id AND ro_in_transit.cust_id = rod_in_transit.cust_id
			WHERE rod_in_transit.pro_id = rod.pro_id 
				AND rod_in_transit.cust_id = rod.cust_id
				AND ro_in_transit.status = 4
				AND ro_in_transit.is_del = false
				AND rod_in_transit.is_del = false
		), 0)`
	} else {
		// For non-principal: use rod.cust_id (same as the replenishment order detail being viewed)
		inTransitStock1Subquery = `COALESCE((
			SELECT SUM(COALESCE(rod_in_transit.order_booking_qty1, 0))
			FROM inv.replenishment_order_detail rod_in_transit
			JOIN inv.replenishment_order ro_in_transit ON ro_in_transit.replenishment_id = rod_in_transit.replenishment_id AND ro_in_transit.cust_id = rod_in_transit.cust_id
			WHERE rod_in_transit.pro_id = rod.pro_id 
				AND rod_in_transit.cust_id = rod.cust_id
				AND ro_in_transit.status = 4
				AND ro_in_transit.is_del = false
				AND rod_in_transit.is_del = false
		), 0)`
		inTransitStock2Subquery = `COALESCE((
			SELECT SUM(COALESCE(rod_in_transit.order_booking_qty2, 0))
			FROM inv.replenishment_order_detail rod_in_transit
			JOIN inv.replenishment_order ro_in_transit ON ro_in_transit.replenishment_id = rod_in_transit.replenishment_id AND ro_in_transit.cust_id = rod_in_transit.cust_id
			WHERE rod_in_transit.pro_id = rod.pro_id 
				AND rod_in_transit.cust_id = rod.cust_id
				AND ro_in_transit.status = 4
				AND ro_in_transit.is_del = false
				AND rod_in_transit.is_del = false
		), 0)`
		inTransitStock3Subquery = `COALESCE((
			SELECT SUM(COALESCE(rod_in_transit.order_booking_qty3, 0))
			FROM inv.replenishment_order_detail rod_in_transit
			JOIN inv.replenishment_order ro_in_transit ON ro_in_transit.replenishment_id = rod_in_transit.replenishment_id AND ro_in_transit.cust_id = rod_in_transit.cust_id
			WHERE rod_in_transit.pro_id = rod.pro_id 
				AND rod_in_transit.cust_id = rod.cust_id
				AND ro_in_transit.status = 4
				AND ro_in_transit.is_del = false
				AND rod_in_transit.is_del = false
		), 0)`
	}

	// Build subquery for warehouse stock total (SUM qty)
	// Use same logic as FindApprovalProducts
	var whsTotalQtySubquery string
	if isPrincipal {
		// For principal: no cust_id filter in warehouse_stock
		whsTotalQtySubquery = `COALESCE((
			SELECT SUM(COALESCE(whs_sub.qty, 0))
			FROM inv.warehouse_stock whs_sub
			WHERE whs_sub.pro_id = rod.pro_id 
				AND whs_sub.wh_id = ro.wh_id
		), 0)`
	} else {
		// For non-principal: filter by custID
		escapedCustID := strings.ReplaceAll(custID, "'", "''")
		whsTotalQtySubquery = fmt.Sprintf(`COALESCE((
			SELECT SUM(COALESCE(whs_sub.qty, 0))
			FROM inv.warehouse_stock whs_sub
			WHERE whs_sub.pro_id = rod.pro_id 
				AND whs_sub.cust_id = '%s'
				AND whs_sub.wh_id = ro.wh_id
		), 0)`, escapedCustID)
	}

	query := repo.model(ctx).
		Select(`rod.replenishment_detail_id,
			rod.replenishment_id,
			ro.replenishment_no,
			rod.pro_id,
			COALESCE(product.pro_code, '') AS pro_code,
			COALESCE(product.pro_name, '') AS pro_name,
			COALESCE(product.unit_id1, '') AS unit_id1,
			COALESCE(product.unit_id2, '') AS unit_id2,
			COALESCE(product.unit_id3, '') AS unit_id3,
			rod.order_booking_qty1,
			rod.order_booking_qty2,
			rod.order_booking_qty3,
			rod.purch_price1,
			rod.purch_price2,
			rod.purch_price3,
			rod.estimated_price,
			product.vat,
			rod.qty_order_allocation1,
			rod.qty_order_allocation2,
			rod.qty_order_allocation3,
			rod.qty_order_approval1,
			rod.qty_order_approval2,
			rod.qty_order_approval3,
			product.saf_stock_qty,
			product.min_stock_qty,
			CASE 
				WHEN `+whsTotalQtySubquery+` = 0 THEN 0
				ELSE 
					`+whsTotalQtySubquery+` -
					(FLOOR(`+whsTotalQtySubquery+` / NULLIF(product.conv_unit2 * product.conv_unit3, 0)) * product.conv_unit2 * product.conv_unit3) -
					(FLOOR(
						(
							`+whsTotalQtySubquery+` - 
							(FLOOR(`+whsTotalQtySubquery+` / NULLIF(product.conv_unit2 * product.conv_unit3, 0)) * product.conv_unit2 * product.conv_unit3)
						) / NULLIF(product.conv_unit2, 0)
					) * product.conv_unit2)
			END AS qty1,
			CASE 
				WHEN `+whsTotalQtySubquery+` = 0 OR product.conv_unit2 = 0 THEN 0
				ELSE FLOOR(
					(
						`+whsTotalQtySubquery+` - 
						(FLOOR(`+whsTotalQtySubquery+` / NULLIF(product.conv_unit2 * product.conv_unit3, 0)) * product.conv_unit2 * product.conv_unit3)
					) / NULLIF(product.conv_unit2, 0)
				)
			END AS qty2,
			CASE 
				WHEN `+whsTotalQtySubquery+` = 0 OR product.conv_unit2 = 0 OR product.conv_unit3 = 0 THEN 0
				ELSE FLOOR(`+whsTotalQtySubquery+` / (product.conv_unit2 * product.conv_unit3))
			END AS qty3,
			`+inTransitStock1Subquery+` AS in_transit_stock1,
			`+inTransitStock2Subquery+` AS in_transit_stock2,
			`+inTransitStock3Subquery+` AS in_transit_stock3`).
		Table("inv.replenishment_order_detail AS rod").
		Joins("LEFT JOIN inv.replenishment_order AS ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id").
		Joins(`LEFT JOIN LATERAL (
			SELECT pro_code, pro_name, vat, unit_id1, unit_id2, unit_id3, saf_stock_qty, min_stock_qty, conv_unit2, conv_unit3
			FROM mst.m_product 
			WHERE pro_id = rod.pro_id AND is_del = false 
			ORDER BY CASE WHEN cust_id = ? THEN 1 WHEN cust_id = ? THEN 2 ELSE 3 END 
			LIMIT 1
		) AS product ON true`, parentCustID, custID).
		Where("rod.replenishment_id = ?", replenishmentIDInt)

	err = query.Order("rod.replenishment_detail_id ASC").Find(&details).Error

	if err != nil {
		return nil, err
	}
	return details, nil
}

func (repo *RepositoryReplenishmentImpl) FindFinalByReplenishmentID(ctx context.Context, replenishmentNo string, custID, parentCustID string, isPrincipal bool) ([]model.ReplenishmentFinalRead, error) {
	var finalData []model.ReplenishmentFinalRead

	var replenishmentOrder model.ReplenishmentOrder
	err := repo.model(ctx).
		Where("replenishment_no = ?", replenishmentNo).
		First(&replenishmentOrder).Error

	if err != nil {
		return nil, err
	}

	replenishmentID := replenishmentOrder.ReplenishmentID

	query := repo.model(ctx).
		Select(`rod.replenishment_detail_id,
			rod.pro_id,
			product.pro_code,
			product.pro_name,
			product.unit_id1,
			product.unit_id2,
			product.unit_id3,
			rod.purch_price1 AS purch_price_delivery1,
			rod.purch_price2 AS purch_price_delivery2,
			rod.purch_price3 AS purch_price_delivery3,
			COALESCE(rod.qty_order_approval1, 0) AS final_order1,
			COALESCE(rod.qty_order_approval2, 0) AS final_order2,
			COALESCE(rod.qty_order_approval3, 0) AS final_order3,
			gr_det.unit_price1 AS gr_price1,
			gr_det.unit_price2 AS gr_price2,
			gr_det.unit_price3 AS gr_price3,
			COALESCE(gr_det.qty1, 0) AS stock_received1,
			COALESCE(gr_det.qty2, 0) AS stock_received2,
			COALESCE(gr_det.qty3, 0) AS stock_received3,
			product.vat`).
		Table("inv.replenishment_order_detail AS rod").
		Joins("LEFT JOIN inv.replenishment_order AS ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id").
		Where("rod.replenishment_id = ?", replenishmentID)

	if !isPrincipal {
		query = query.Where("rod.cust_id = ?", custID)
	}

	query = query.
		Joins(`LEFT JOIN LATERAL (
			SELECT pro_code, pro_name, unit_id1, unit_id2, unit_id3, vat
			FROM mst.m_product 
			WHERE pro_id = rod.pro_id AND is_del = false 
			ORDER BY CASE WHEN cust_id = ? THEN 1 WHEN cust_id = ? THEN 2 ELSE 3 END 
			LIMIT 1
		) AS product ON true`, parentCustID, custID)

	if isPrincipal {
		query = query.
			Joins("LEFT JOIN inv.gr AS gr ON gr.po_no = ro.replenishment_no AND gr.cust_id = ro.cust_id AND gr.with_reference = true").
			Joins("LEFT JOIN inv.gr_det AS gr_det ON gr_det.gr_no = gr.gr_no AND gr_det.pro_id = rod.pro_id AND gr_det.cust_id = ro.cust_id AND gr_det.item_type = 1")
	} else {
		query = query.
			Joins("LEFT JOIN inv.gr AS gr ON gr.po_no = ro.replenishment_no AND gr.cust_id = ? AND gr.with_reference = true", custID).
			Joins("LEFT JOIN inv.gr_det AS gr_det ON gr_det.gr_no = gr.gr_no AND gr_det.pro_id = rod.pro_id AND gr_det.cust_id = ? AND gr_det.item_type = 1", custID)
	}

	err = query.Order("rod.replenishment_detail_id ASC").Find(&finalData).Error

	if err != nil {
		if strings.Contains(err.Error(), "column") && strings.Contains(err.Error(), "does not exist") {
			retryQuery := repo.model(ctx).
				Select(`rod.replenishment_detail_id,
					rod.pro_id,
					product.pro_code,
					product.pro_name,
					product.unit_id1,
					product.unit_id2,
					product.unit_id3,
					rod.purch_price1 AS purch_price_delivery1,
					rod.purch_price2 AS purch_price_delivery2,
					rod.purch_price3 AS purch_price_delivery3,
					COALESCE(rod.qty_order_approval1, 0) AS final_order1,
					COALESCE(rod.qty_order_approval2, 0) AS final_order2,
					COALESCE(rod.qty_order_approval3, 0) AS final_order3,
					gr_det.unit_price1 AS gr_price1,
					gr_det.unit_price2 AS gr_price2,
					gr_det.unit_price3 AS gr_price3,
					COALESCE(gr_det.qty1, 0) AS stock_received1,
					COALESCE(gr_det.qty2, 0) AS stock_received2,
					COALESCE(gr_det.qty3, 0) AS stock_received3,
					product.vat`).
				Table("inv.replenishment_order_detail AS rod").
				Joins("LEFT JOIN inv.replenishment_order AS ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id").
				Where("rod.replenishment_id = ?", replenishmentID)

			if !isPrincipal {
				retryQuery = retryQuery.Where("rod.cust_id = ?", custID)
			}

			retryQuery = retryQuery.
				Joins(`LEFT JOIN LATERAL (
					SELECT pro_code, pro_name, unit_id1, unit_id2, unit_id3, vat
					FROM mst.m_product 
					WHERE pro_id = rod.pro_id AND is_del = false 
					ORDER BY CASE WHEN cust_id = ? THEN 1 WHEN cust_id = ? THEN 2 ELSE 3 END 
					LIMIT 1
				) AS product ON true`, parentCustID, custID)

			if isPrincipal {
				retryQuery = retryQuery.
					Joins("LEFT JOIN inv.gr AS gr ON gr.po_no = ro.replenishment_no AND gr.cust_id = ro.cust_id AND gr.with_reference = true").
					Joins("LEFT JOIN inv.gr_det AS gr_det ON gr_det.gr_no = gr.gr_no AND gr_det.pro_id = rod.pro_id AND gr_det.cust_id = ro.cust_id AND gr_det.item_type = 1")
			} else {
				retryQuery = retryQuery.
					Joins("LEFT JOIN inv.gr AS gr ON gr.po_no = ro.replenishment_no AND gr.cust_id = ? AND gr.with_reference = true", custID).
					Joins("LEFT JOIN inv.gr_det AS gr_det ON gr_det.gr_no = gr.gr_no AND gr_det.pro_id = rod.pro_id AND gr_det.cust_id = ? AND gr_det.item_type = 1", custID)
			}

			err = retryQuery.Order("rod.replenishment_detail_id ASC").Find(&finalData).Error
		}
	}

	if err != nil {
		return nil, err
	}
	return finalData, nil
}

func (repo *RepositoryReplenishmentImpl) FindGoodReceiptByReplenishmentNo(ctx context.Context, replenishmentNo string, custID, parentCustID string, isPrincipal bool) ([]model.ReplenishmentGoodReceiptRead, error) {
	var goodReceiptData []model.ReplenishmentGoodReceiptRead

	query := repo.model(ctx).
		Select(`gr.po_no,
			gr_det.pro_id,
			product.pro_code,
			product.pro_name,
			gr_det.unit_price1,
			gr_det.unit_price2,
			gr_det.unit_price3,
			gr_det.qty1 AS qty_received1,
			gr_det.qty2 AS qty_received2,
			gr_det.qty3 AS qty_received3,
			(gr_det.unit_price1 * gr_det.qty1 + gr_det.unit_price2 * gr_det.qty2 + gr_det.unit_price3 * gr_det.qty3) AS estimated_price,
			product.vat`).
		Table("inv.gr AS gr").
		Where("gr.po_no = ? AND gr.with_reference = true", replenishmentNo)

	if isPrincipal {
		query = query.
			Joins("JOIN inv.gr_det AS gr_det ON gr.gr_no = gr_det.gr_no AND gr_det.cust_id = gr.cust_id").
			Joins(`LEFT JOIN LATERAL (
				SELECT pro_code, pro_name, vat
				FROM mst.m_product 
				WHERE pro_id = gr_det.pro_id AND is_del = false 
				ORDER BY CASE WHEN cust_id = ? THEN 1 WHEN cust_id = ? THEN 2 ELSE 3 END 
				LIMIT 1
			) AS product ON true`, parentCustID, custID)
	} else {
		query = query.
			Joins("JOIN inv.gr_det AS gr_det ON gr.gr_no = gr_det.gr_no AND gr_det.cust_id = ?", custID).
			Joins(`LEFT JOIN LATERAL (
				SELECT pro_code, pro_name, vat
				FROM mst.m_product 
				WHERE pro_id = gr_det.pro_id AND is_del = false 
				ORDER BY CASE WHEN cust_id = ? THEN 1 WHEN cust_id = ? THEN 2 ELSE 3 END 
				LIMIT 1
			) AS product ON true`, parentCustID, custID).
			Where("gr.cust_id = ?", custID)
	}

	err := query.Order("gr_det.pro_id ASC").Find(&goodReceiptData).Error

	if err != nil {
		return nil, err
	}
	return goodReceiptData, nil
}

func (repo *RepositoryReplenishmentImpl) FindProductList(ctx context.Context, dataFilter entity.ReplenishmentProductQueryFilter, custId, parentCustId string) ([]model.ReplenishmentProductList, int64, int, error) {
	var products []model.ReplenishmentProductList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = constant.DEFAULT_PAGE_LIMIT
	} else {
		limit = dataFilter.Limit
	}

	// Build base query for counting distinct products
	queryCount := repo.model(ctx).
		Table("sls.\"order\" o").
		Joins("JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = ?", custId).
		Where("o.cust_id = ?", custId).
		Where("o.data_status = 1").
		Where("o.wh_id = ?", dataFilter.WhID)

	// Build main query with aggregation - set Table first, then Select
	query := repo.model(ctx).
		Table("sls.\"order\" o").
		Select(`
			od.pro_id,
			COALESCE(p.pro_code, '') AS pro_code,
			COALESCE(p.pro_name, '') AS pro_name,
			od.purch_price3,
			od.purch_price2,
			od.purch_price1,
			COALESCE(p.unit_id3, '') AS unit_id3,
			COALESCE(p.unit_id2, '') AS unit_id2,
			COALESCE(p.unit_id1, '') AS unit_id1,
			COALESCE(p.vat, 0) AS vat,
			SUM(od.qty1) AS qty1,
			SUM(od.qty2) AS qty2,
			SUM(od.qty3) AS qty3,
			COALESCE((
				SELECT SUM(COALESCE(rod.qty_order_approval1, 0))
				FROM inv.replenishment_order_detail rod
				JOIN inv.replenishment_order ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id
				WHERE rod.pro_id = od.pro_id 
					AND rod.cust_id = ?
					AND ro.status = 4
					AND ro.is_del = false
					AND rod.is_del = false
			), 0) AS in_transit_stock1,
			COALESCE((
				SELECT SUM(COALESCE(rod.qty_order_approval2, 0))
				FROM inv.replenishment_order_detail rod
				JOIN inv.replenishment_order ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id
				WHERE rod.pro_id = od.pro_id 
					AND rod.cust_id = ?
					AND ro.status = 4
					AND ro.is_del = false
					AND rod.is_del = false
			), 0) AS in_transit_stock2,
			COALESCE((
				SELECT SUM(COALESCE(rod.qty_order_approval3, 0))
				FROM inv.replenishment_order_detail rod
				JOIN inv.replenishment_order ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id
				WHERE rod.pro_id = od.pro_id 
					AND rod.cust_id = ?
					AND ro.status = 4
					AND ro.is_del = false
					AND rod.is_del = false
			), 0) AS in_transit_stock3
		`, custId, custId, custId).
		Joins("JOIN sls.order_detail od ON od.ro_no = o.ro_no AND od.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_product p ON p.pro_id = od.pro_id AND p.cust_id = ?", custId).
		Where("o.cust_id = ?", custId).
		Where("o.data_status = 1").
		Where("o.wh_id = ?", dataFilter.WhID)

	// Date range filter with timezone Indonesia
	if dataFilter.StartDate != nil && dataFilter.EndDate != nil {
		startTimeUTC := time.Unix(*dataFilter.StartDate, 0).UTC()
		endTimeUTC := time.Unix(*dataFilter.EndDate, 0).UTC()
		startDateStr := startTimeUTC.Format("2006-01-02")
		endDateStr := endTimeUTC.Format("2006-01-02")
		query = query.Where("o.ro_date BETWEEN ? AND ?", startDateStr, endDateStr)
		queryCount = queryCount.Where("o.ro_date BETWEEN ? AND ?", startDateStr, endDateStr)
	}

	// Group by clause (matching requirement: group by pro_id and prices)
	query = query.Group(`
		od.pro_id,
		od.purch_price3,
		od.purch_price2,
		od.purch_price1,
		p.pro_code,
		p.pro_name,
		p.unit_id3,
		p.unit_id2,
		p.unit_id1,
		p.vat
	`)

	// Sorting
	sortBy := ""
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortField := colSort[0]
				if sortField == "created_date" {
					sortField = "od.pro_id"
				}
				sortBy += fmt.Sprintf(`%s %s, `, sortField, colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		if sortBy != "" {
			query.Order(sortBy)
		}
	}
	// Default sort by pro_id
	if sortBy == "" {
		query.Order("od.pro_id ASC")
	}

	// Count total distinct products
	var countResult struct {
		Count int64 `gorm:"column:count"`
	}
	err := queryCount.Select("COUNT(DISTINCT od.pro_id) as count").Scan(&countResult).Error
	if err != nil {
		return products, total, 0, err
	}
	total = countResult.Count

	// Pagination
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	err = query.Limit(limit).Offset(offset).Scan(&products).Error
	if err != nil {
		return products, total, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return products, total, lastPage, nil
}

func (repo *RepositoryReplenishmentImpl) FindByReplenishmentNo(ctx context.Context, replenishmentNo string, custID string) (*model.ReplenishmentOrder, error) {
	var replenishment model.ReplenishmentOrder
	query := repo.model(ctx).
		Where("replenishment_no = ? AND is_del = false", replenishmentNo)

	if custID != "" {
		query = query.Where("cust_id = ?", custID)
	}

	err := query.First(&replenishment).Error
	if err != nil {
		return nil, err
	}
	return &replenishment, nil
}

func (repo *RepositoryReplenishmentImpl) GetReplenishmentOrderByID(ctx context.Context, replenishmentID int64, custID string) (*model.ReplenishmentOrder, error) {
	var replenishment model.ReplenishmentOrder
	query := repo.model(ctx).
		Where("replenishment_id = ? AND is_del = false", replenishmentID)

	if custID != "" {
		query = query.Where("cust_id = ?", custID)
	}

	err := query.First(&replenishment).Error
	if err != nil {
		return nil, err
	}
	return &replenishment, nil
}

func (repo *RepositoryReplenishmentImpl) FindDetailByReplenishmentIDForUpdate(ctx context.Context, replenishmentID int64, custID string) ([]model.ReplenishmentOrderDetail, error) {
	var details []model.ReplenishmentOrderDetail
	query := repo.model(ctx).
		Where("replenishment_id = ? AND is_del = false", replenishmentID)

	// Only filter by cust_id if custID is provided (non-empty)
	// Principal users can access without cust_id filter
	if custID != "" {
		query = query.Where("cust_id = ?", custID)
	}

	err := query.Find(&details).Error
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (repo *RepositoryReplenishmentImpl) UpdateApproval(ctx context.Context, replenishment *model.ReplenishmentOrder) error {
	updates := map[string]interface{}{
		"status":      replenishment.Status,
		"is_approval": replenishment.IsApproval,
		"approve_by":  replenishment.ApproveBy,
		"approve_at":  replenishment.ApproveAt,
		"updated_at":  time.Now(),
	}
	if replenishment.UpdatedBy != nil {
		updates["updated_by"] = *replenishment.UpdatedBy
	}

	result := repo.model(ctx).
		Model(&model.ReplenishmentOrder{}).
		Where("replenishment_id = ? AND cust_id = ? AND is_del = false", replenishment.ReplenishmentID, replenishment.CustID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf(ErrNoRowsAffected)
	}
	return nil
}

func (repo *RepositoryReplenishmentImpl) UpdateDetail(ctx context.Context, detail *model.ReplenishmentOrderDetail) error {
	now := time.Now()
	detail.UpdatedAt = &now

	// Build update map dynamically to include all fields
	updates := map[string]interface{}{
		"updated_at": detail.UpdatedAt,
	}

	// Only include fields that should be updated (not nil)
	if detail.ReturnReasonID != nil {
		updates["return_reason_id"] = detail.ReturnReasonID
	}
	if detail.QtyOrderApproval1 != nil {
		updates["qty_order_approval1"] = detail.QtyOrderApproval1
	}
	if detail.QtyOrderApproval2 != nil {
		updates["qty_order_approval2"] = detail.QtyOrderApproval2
	}
	if detail.QtyOrderApproval3 != nil {
		updates["qty_order_approval3"] = detail.QtyOrderApproval3
	}
	if detail.QtyOrderAllocation1 != nil {
		updates["qty_order_allocation1"] = detail.QtyOrderAllocation1
	}
	if detail.QtyOrderAllocation2 != nil {
		updates["qty_order_allocation2"] = detail.QtyOrderAllocation2
	}
	if detail.QtyOrderAllocation3 != nil {
		updates["qty_order_allocation3"] = detail.QtyOrderAllocation3
	}

	// Always update price and estimated_price fields (they are non-pointer float64 in model)
	// These will be set from request if provided, otherwise keep existing value
	updates["purch_price1"] = detail.PurchPrice1
	updates["purch_price2"] = detail.PurchPrice2
	updates["purch_price3"] = detail.PurchPrice3
	updates["estimated_price"] = detail.EstimatedPrice

	if detail.UpdatedBy != nil {
		updates["updated_by"] = *detail.UpdatedBy
	}

	if len(updates) <= 1 {
		// Only updated_at means no actual data fields to update
		return fmt.Errorf(ErrNoFieldsToUpdate)
	}

	result := repo.model(ctx).
		Model(&model.ReplenishmentOrderDetail{}).
		Where("replenishment_detail_id = ? AND cust_id = ? AND is_del = false", detail.ReplenishmentDetailID, detail.CustID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf(ErrNoRowsAffected)
	}
	return nil
}

func (repo *RepositoryReplenishmentImpl) CreateDetailForApproval(ctx context.Context, detail *model.ReplenishmentOrderDetail) error {
	now := time.Now()
	detail.CreatedAt = now
	detail.UpdatedAt = &now
	if detail.CreatedBy != 0 {
		detail.UpdatedBy = &detail.CreatedBy
	}

	return repo.model(ctx).Create(detail).Error
}

func (repo *RepositoryReplenishmentImpl) SoftDeleteDetail(ctx context.Context, detailID int64, custID string, deletedBy int64) error {
	now := time.Now()
	query := repo.model(ctx).
		Model(&model.ReplenishmentOrderDetail{}).
		Where("replenishment_detail_id = ? AND is_del = false", detailID)

	if custID != "" {
		query = query.Where("cust_id = ?", custID)
	}

	result := query.Updates(map[string]interface{}{
		"is_del":     true,
		"deleted_by": deletedBy,
		"deleted_at": now,
	})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf(ErrNoRowsAffected)
	}
	return nil
}

func (repo *RepositoryReplenishmentImpl) FindProductGrList(ctx context.Context, dataFilter entity.ProductGrListQueryFilter, custId, parentCustId string) (*model.ReplenishmentOrderRead, []model.ProductGrListDetail, int64, int, error) {
	var header model.ReplenishmentOrderRead
	var details []model.ProductGrListDetail
	var total int64

	limit := dataFilter.Limit
	if limit == 0 {
		limit = constant.DEFAULT_PAGE_LIMIT
	}
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	// Get header by replenishment_no
	headerQuery := repo.model(ctx).
		Table("inv.replenishment_order AS ro").
		Select(`ro.cust_id,
			ro.replenishment_id,
			ro.replenishment_no,
			ro.date,
			ro.sup_id,
			sup.sup_code,
			sup.sup_name,
			ro.wh_id,
			wh.wh_code,
			wh.wh_name,
			ro.delivery_type,
			ro.replenishment_type,
			ro.so_start_date,
			ro.so_end_date,
			ro.delivery_date,
			ro.note,
			ro.status,
			rs.status_name,
			ro.so_no,
			ro.created_by,
			ro.created_at,
			ro.updated_by,
			ro.updated_at,
			ro.deleted_by,
			ro.deleted_at,
			ro.is_del`).
		Joins("LEFT JOIN mst.m_supplier AS sup ON sup.sup_id = ro.sup_id AND sup.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_warehouse AS wh ON wh.wh_id = ro.wh_id AND wh.cust_id = ?", custId).
		Joins("LEFT JOIN inv.replenishment_status AS rs ON rs.status_code = ro.status").
		Where("ro.cust_id = ?", custId).
		Where("ro.replenishment_no = ?", dataFilter.ReplenishmentNo).
		Where("ro.is_del = false").
		Limit(1)

	if err := headerQuery.Scan(&header).Error; err != nil {
		return nil, nil, 0, 0, err
	}

	// Get details with search filter
	detailQuery := repo.model(ctx).
		Table("inv.replenishment_order_detail AS rod").
		Select(`rod.replenishment_detail_id,
			rod.pro_id,
			COALESCE(product.pro_code, '') AS pro_code,
			COALESCE(product.pro_name, '') AS pro_name,
			product.vat,
			rod.qty_order_approval1,
			rod.qty_order_approval2,
			rod.qty_order_approval3,
			COALESCE(product.unit_id1, '') AS unit_id1,
			COALESCE(product.unit_id2, '') AS unit_id2,
			COALESCE(product.unit_id3, '') AS unit_id3,
			rod.purch_price1,
			rod.purch_price2,
			rod.purch_price3`).
		Joins("LEFT JOIN inv.replenishment_order AS ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id").
		Joins(`LEFT JOIN LATERAL (
			SELECT pro_code, pro_name, vat, unit_id1, unit_id2, unit_id3 
			FROM mst.m_product 
			WHERE pro_id = rod.pro_id AND is_del = false 
			ORDER BY CASE WHEN cust_id = ? THEN 1 WHEN cust_id = ? THEN 2 ELSE 3 END 
			LIMIT 1
		) AS product ON true`, parentCustId, header.CustID).
		Where("rod.replenishment_id = ? AND rod.cust_id = ? AND rod.is_del = false", header.ReplenishmentID, header.CustID)

	// Search filter
	if dataFilter.Q != "" {
		like := "%" + dataFilter.Q + "%"
		detailQuery = detailQuery.Where("(product.pro_code ILIKE ? OR product.pro_name ILIKE ?)", like, like)
	}

	// Count total
	if err := detailQuery.Count(&total).Error; err != nil {
		return nil, nil, 0, 0, err
	}

	// Execute query with pagination
	err := detailQuery.Order("rod.replenishment_detail_id ASC").Limit(limit).Offset(offset).Scan(&details).Error
	if err != nil {
		return nil, nil, 0, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return &header, details, total, lastPage, nil
}

func (repo *RepositoryReplenishmentImpl) FindPoList(ctx context.Context, dataFilter entity.PoListQueryFilter, custId, parentCustId string) ([]model.PoList, int64, int, error) {
	var poList []model.PoList
	var total int64

	limit := dataFilter.Limit
	if limit == 0 {
		limit = constant.DEFAULT_PAGE_LIMIT
	}
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	query := repo.model(ctx).
		Table("inv.replenishment_order ro").
		Select(`
			ro.replenishment_no,
			ro.replenishment_type,
			ro.wh_id,
			wh.wh_code,
			wh.wh_name,
			ro.sup_id,
			sup.sup_code,
			sup.sup_name`).
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = ro.wh_id AND wh.cust_id = ?", custId).
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = ro.sup_id AND sup.cust_id = ?", custId).
		Where("ro.cust_id = ?", custId).
		Where("ro.is_del = false").
		Where("ro.status = ?", 4)

	// Count query
	countQuery := repo.model(ctx).
		Table("inv.replenishment_order ro").
		Where("ro.cust_id = ?", custId).
		Where("ro.is_del = false").
		Where("ro.status = ?", 4)

	// Count total
	if err := countQuery.Count(&total).Error; err != nil {
		return poList, 0, 0, err
	}

	// Execute query with pagination
	err := query.Order("ro.replenishment_no DESC").Limit(limit).Offset(offset).Scan(&poList).Error
	if err != nil {
		return poList, 0, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return poList, total, lastPage, nil
}

func (repo *RepositoryReplenishmentImpl) FindApprovalProducts(ctx context.Context, dataFilter entity.ReplenishmentApprovalProductQueryFilter, custId, parentCustId string, isPrincipal bool) ([]model.ReplenishmentApprovalProduct, int64, int, error) {
	var products []model.ReplenishmentApprovalProduct
	var total int64

	limit := dataFilter.Limit
	if limit == 0 {
		limit = constant.DEFAULT_PAGE_LIMIT
	}
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	warehouseFilter := ""
	if dataFilter.WhID > 0 {
		warehouseFilter = fmt.Sprintf(" AND whs_sub.wh_id = %d", dataFilter.WhID)
	}
	whsTotalQtySubquery := fmt.Sprintf(`COALESCE((
		SELECT SUM(COALESCE(whs_sub.qty, 0))
		FROM inv.warehouse_stock whs_sub
		WHERE whs_sub.pro_id = p.pro_id
		  AND whs_sub.cust_id = p.cust_id
		  %s
	), 0)`, warehouseFilter)

	qty1Expression, qty2Expression, qty3Expression := BuildQtyCalculationExpressions(whsTotalQtySubquery, "p")

	inTransitStock1Subquery := `COALESCE((
		SELECT SUM(COALESCE(rod.order_booking_qty1, 0))
		FROM inv.replenishment_order_detail rod
		JOIN inv.replenishment_order ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id
		WHERE rod.pro_id = p.pro_id
		  AND rod.cust_id = p.cust_id
		  AND ro.status = 4
		  AND rod.is_del = false
		  AND ro.is_del = false
	), 0)`
	inTransitStock2Subquery := `COALESCE((
		SELECT SUM(COALESCE(rod.order_booking_qty2, 0))
		FROM inv.replenishment_order_detail rod
		JOIN inv.replenishment_order ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id
		WHERE rod.pro_id = p.pro_id
		  AND rod.cust_id = p.cust_id
		  AND ro.status = 4
		  AND rod.is_del = false
		  AND ro.is_del = false
	), 0)`
	inTransitStock3Subquery := `COALESCE((
		SELECT SUM(COALESCE(rod.order_booking_qty3, 0))
		FROM inv.replenishment_order_detail rod
		JOIN inv.replenishment_order ro ON ro.replenishment_id = rod.replenishment_id AND ro.cust_id = rod.cust_id
		WHERE rod.pro_id = p.pro_id
		  AND rod.cust_id = p.cust_id
		  AND ro.status = 4
		  AND rod.is_del = false
		  AND ro.is_del = false
	), 0)`
	ripeningDistributorExpr := "mc.distributor_id"
	if dataFilter.DistributorID != nil && *dataFilter.DistributorID > 0 {
		ripeningDistributorExpr = fmt.Sprintf("%d", *dataFilter.DistributorID)
	}
	ripeningSubquery := `COALESCE((
		SELECT MAX(
			CASE EXTRACT(DOW FROM ((CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'))::date)
				WHEN 1 THEN pr.monday_qty
				WHEN 2 THEN pr.tuesday_qty
				WHEN 3 THEN pr.wednesday_qty
				WHEN 4 THEN pr.thursday_qty
				WHEN 5 THEN pr.friday_qty
				WHEN 6 THEN pr.saturday_qty
				WHEN 0 THEN pr.sunday_qty
			END
		)
		FROM mst.product_ripening pr
		WHERE pr.pro_id = p.pro_id
		  AND pr.created_by = ` + fmt.Sprintf("%d", dataFilter.UserID) + `
		  AND ` + ripeningDistributorExpr + ` IS NOT NULL
		  AND ` + ripeningDistributorExpr + ` > 0
		  AND pr.distributor_id = ` + ripeningDistributorExpr + `
	), 0)`

	query := repo.model(ctx).
		Table("mst.m_product p").
		Select(`
			p.pro_id,
			COALESCE(p.pro_code, '') AS pro_code,
			COALESCE(p.pro_name, '') AS pro_name,
			` + ripeningSubquery + ` AS ripening,
			` + inTransitStock1Subquery + ` AS in_transit_stock1,
			` + inTransitStock2Subquery + ` AS in_transit_stock2,
			` + inTransitStock3Subquery + ` AS in_transit_stock3,
			p.saf_stock_qty,
			p.min_stock_qty,
			p.vat,
			p.conv_unit2,
			p.conv_unit3,
			p.unit_id1,
			p.unit_id2,
			p.unit_id3,
			p.purch_price1,
			p.purch_price2,
			p.purch_price3,
			` + whsTotalQtySubquery + ` AS total_qty,
			` + qty3Expression + `,
			` + qty2Expression + `,
			` + qty1Expression)

	query = query.
		Joins("JOIN smc.m_customer mc ON mc.cust_id = p.cust_id").
		Where("p.is_del = false").
		Where("p.is_active = true")
	if isPrincipal {
		query = query.
			Joins("JOIN smc.m_customer prod_cust ON prod_cust.cust_id = p.cust_id").
			Where(`prod_cust.distributor_id IN (
				SELECT DISTINCT child.distributor_id
				FROM smc.m_customer child
				WHERE child.parent_cust_id = ?
				  AND child.distributor_id IS NOT NULL
				  AND child.distributor_id > 0
			)`, custId)
	} else {
		query = query.Where("p.cust_id = ?", custId)
	}

	if dataFilter.Q != "" {
		like := "%" + dataFilter.Q + "%"
		query = query.Where("(p.pro_code ILIKE ? OR p.pro_name ILIKE ?)", like, like)
	}

	if dataFilter.SupID != nil {
		query = query.Where("p.sup_id = ?", *dataFilter.SupID)
	}

	// Apply zero_stock filter (using subquery instead of HAVING since we removed GROUP BY)
	if dataFilter.ZeroStock != nil {
		if *dataFilter.ZeroStock {
			query = query.Where(whsTotalQtySubquery + " = 0")
		} else {
			query = query.Where(whsTotalQtySubquery + " > 0")
		}
	}
	countQuery := repo.model(ctx).
		Table("mst.m_product p")

	countQuery = countQuery.Where("p.is_del = false").
		Where("p.is_active = true")
	if isPrincipal {
		countQuery = countQuery.
			Joins("JOIN smc.m_customer prod_cust ON prod_cust.cust_id = p.cust_id").
			Where(`prod_cust.distributor_id IN (
				SELECT DISTINCT child.distributor_id
				FROM smc.m_customer child
				WHERE child.parent_cust_id = ?
				  AND child.distributor_id IS NOT NULL
				  AND child.distributor_id > 0
			)`, custId)
	} else {
		countQuery = countQuery.Where("p.cust_id = ?", custId)
	}

	if dataFilter.Q != "" {
		like := "%" + dataFilter.Q + "%"
		countQuery = countQuery.Where("(p.pro_code ILIKE ? OR p.pro_name ILIKE ?)", like, like)
	}
	if dataFilter.SupID != nil {
		countQuery = countQuery.Where("p.sup_id = ?", *dataFilter.SupID)
	}

	if dataFilter.ZeroStock != nil {
		if *dataFilter.ZeroStock {
			countQuery = countQuery.Where(whsTotalQtySubquery + " = 0")
		} else {
			countQuery = countQuery.Where(whsTotalQtySubquery + " > 0")
		}
	}

	if err := countQuery.Select("COUNT(DISTINCT p.pro_id)").Scan(&total).Error; err != nil {
		return products, 0, 0, err
	}

	err := query.Order("p.pro_name ASC").Limit(limit).Offset(offset).Scan(&products).Error
	if err != nil {
		return products, 0, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return products, total, lastPage, nil
}

// CheckIsPrincipal checks if cust_id exists in mst.m_principal table
func (repo *RepositoryReplenishmentImpl) CheckIsPrincipal(ctx context.Context, custID string) (bool, error) {
	var count int64
	err := repo.model(ctx).
		Table("mst.m_principal").
		Where("cust_id = ? AND is_del = false", custID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *RepositoryReplenishmentImpl) UpdateStatusByReplenishmentNo(ctx context.Context, replenishmentNo string, custID string, status int, updatedBy int64) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
		"updated_by": updatedBy,
	}

	result := repo.model(ctx).
		Model(&model.ReplenishmentOrder{}).
		Where("replenishment_no = ? AND cust_id = ? AND is_del = false", replenishmentNo, custID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf(ErrNoRowsAffected)
	}
	return nil
}

func (repo *RepositoryReplenishmentImpl) LookupProductByProCode(ctx context.Context, custID string, proCode string) (proID int64, unitID3 *string, err error) {
	proCode = strings.TrimSpace(proCode)
	custID = strings.TrimSpace(custID)
	if proCode == "" || custID == "" {
		return 0, nil, fmt.Errorf("invalid product or customer")
	}

	type row struct {
		ProID   int64   `gorm:"column:pro_id"`
		UnitID3 *string `gorm:"column:unit_id3"`
	}
	var picked row
	err = repo.model(ctx).Raw(`
SELECT p.pro_id, p.unit_id3
FROM mst.m_product p
WHERE p.is_del = false
  AND UPPER(TRIM(p.pro_code)) = UPPER(TRIM(?))
  AND (
	UPPER(TRIM(p.cust_id)) = UPPER(TRIM(?))
	OR UPPER(TRIM(p.cust_id)) = (
		SELECT UPPER(TRIM(COALESCE(NULLIF(TRIM(mc.parent_cust_id), ''), mc.cust_id)))
		FROM smc.m_customer mc
		WHERE UPPER(TRIM(mc.cust_id)) = UPPER(TRIM(?))
		LIMIT 1
	)
  )
ORDER BY CASE WHEN UPPER(TRIM(p.cust_id)) = UPPER(TRIM(?)) THEN 0 ELSE 1 END
LIMIT 1
`, proCode, custID, custID, custID).Scan(&picked).Error

	if err != nil {
		return 0, nil, err
	}
	if picked.ProID <= 0 {
		return 0, nil, gorm.ErrRecordNotFound
	}
	return picked.ProID, picked.UnitID3, nil
}

func (repo *RepositoryReplenishmentImpl) FindActiveReplenishmentDetail(ctx context.Context, custID string, replenishmentID, proID int64) (*model.ReplenishmentOrderDetail, error) {
	var d model.ReplenishmentOrderDetail
	err := repo.model(ctx).
		Where("cust_id = ? AND replenishment_id = ? AND pro_id = ? AND is_del = ?", custID, replenishmentID, proID, false).
		First(&d).Error
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (repo *RepositoryReplenishmentImpl) UpdateReplenishmentDetailSAPFields(ctx context.Context, custID string, replenishmentDetailID int64, sapQty3, sapPurchPrice3 *float64, updatedBy int64) error {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
		"updated_by": updatedBy,
	}
	if sapQty3 != nil {
		updates["sap_qty3"] = *sapQty3
	}
	if sapPurchPrice3 != nil {
		updates["sap_purch_price3"] = *sapPurchPrice3
	}
	if len(updates) == 2 {
		return nil
	}
	result := repo.model(ctx).
		Model(&model.ReplenishmentOrderDetail{}).
		Where("cust_id = ? AND replenishment_detail_id = ? AND is_del = ?", custID, replenishmentDetailID, false).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf(ErrNoRowsAffected)
	}
	return nil
}

func (repo *RepositoryReplenishmentImpl) GetDistributorIDByCustID(ctx context.Context, custID string) (int64, error) {
	var distributorID int64
	err := repo.model(ctx).
		Table("smc.m_customer").
		Select("COALESCE(distributor_id, 0)").
		Where("cust_id = ?", custID).
		Limit(1).
		Scan(&distributorID).Error
	return distributorID, err
}

func (repo *RepositoryReplenishmentImpl) FindCustIDByDistributorCode(ctx context.Context, distributorCode int) (string, error) {
	if distributorCode <= 0 {
		return "", gorm.ErrRecordNotFound
	}

	codeStr := strconv.Itoa(distributorCode)
	var custIDs []string
	err := repo.model(ctx).
		Table("mst.m_distributor").
		Where("is_del = ?", false).
		Where("TRIM(distributor_code) = ?", codeStr).
		Where("COALESCE(NULLIF(TRIM(cust_id), ''), '') <> ''").
		Pluck("cust_id", &custIDs).Error
	if err != nil {
		return "", err
	}
	if len(custIDs) == 0 {
		return "", gorm.ErrRecordNotFound
	}
	if len(custIDs) > 1 {
		return "", fmt.Errorf("distributor_code matches multiple distributors")
	}
	return strings.TrimSpace(custIDs[0]), nil
}

func (repo *RepositoryReplenishmentImpl) IsDistributorApprovalPIC(ctx context.Context, userID int64, supID int64, distributorID int64) (bool, error) {
	var count int64
	err := repo.model(ctx).
		Table("mst.distributor_replenishment_approval dra").
		Joins("JOIN mst.distributor_replenishment_setup drs ON drs.id = dra.dist_replenishment_setup_id").
		Where("drs.sup_id = ? AND drs.distributor_id = ? AND dra.pic = ? AND drs.is_del = false AND dra.is_del = false", supID, distributorID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *RepositoryReplenishmentImpl) GetDistributorApprovalRequirement(ctx context.Context, supID int64, distributorID int64) (bool, bool, error) {
	type setupResult struct {
		IsApprovalRequired *bool `gorm:"column:is_approval_required"`
	}
	var result setupResult
	err := repo.model(ctx).
		Table("mst.distributor_replenishment_setup").
		Select("is_approval_required").
		Where("sup_id = ? AND distributor_id = ? AND is_del = false", supID, distributorID).
		Order("id DESC").
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return false, false, err
	}
	if result.IsApprovalRequired == nil {
		return false, false, nil
	}
	return true, *result.IsApprovalRequired, nil
}

func (repo *RepositoryReplenishmentImpl) GetInitialDistributorApproval(ctx context.Context, supID int64, distributorID int64) (int, int, int64, error) {
	type approvalRow struct {
		Level    int   `gorm:"column:level"`
		Sequence int   `gorm:"column:sequence"`
		Pic      int64 `gorm:"column:pic"`
	}
	var row approvalRow
	err := repo.model(ctx).
		Table("mst.distributor_replenishment_approval dra").
		Select("dra.level, dra.sequence, dra.pic").
		Joins("JOIN mst.distributor_replenishment_setup drs ON drs.id = dra.dist_replenishment_setup_id").
		Where("drs.sup_id = ? AND drs.distributor_id = ? AND drs.is_del = false AND dra.is_del = false", supID, distributorID).
		Order("dra.sequence ASC, dra.level ASC").
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return 0, 0, 0, err
	}
	if row.Pic == 0 {
		return 0, 0, 0, fmt.Errorf("no approval PIC configured")
	}
	return row.Level, row.Sequence, row.Pic, nil
}

func (repo *RepositoryReplenishmentImpl) InsertReplenishmentOrderApproval(ctx context.Context, custID string, replenishmentID int64, level int, sequence int, pic int64) error {
	now := time.Now()

	var tableName string
	err := repo.model(ctx).Raw(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'inv'
		  AND table_name IN ('replenishment_order_approval', 'replanishment_order_approval')
		ORDER BY CASE
			WHEN table_name = 'replenishment_order_approval' THEN 1
			ELSE 2
		END
		LIMIT 1
	`).Scan(&tableName).Error
	if err != nil {
		return fmt.Errorf("failed to resolve approval queue table: %w", err)
	}
	if strings.TrimSpace(tableName) == "" {
		return fmt.Errorf("approval queue table not found in schema inv")
	}

	var orderIDColumn string
	err = repo.model(ctx).Raw(`
		SELECT column_name
		FROM information_schema.columns
		WHERE table_schema = 'inv'
		  AND table_name = ?
		  AND column_name IN ('replenishment_order_id', 'replanishment_order_id')
		ORDER BY CASE
			WHEN column_name = 'replenishment_order_id' THEN 1
			ELSE 2
		END
		LIMIT 1
	`, tableName).Scan(&orderIDColumn).Error
	if err != nil {
		return fmt.Errorf("failed to resolve approval queue order id column: %w", err)
	}
	if strings.TrimSpace(orderIDColumn) == "" {
		return fmt.Errorf("approval queue order id column not found on inv.%s", tableName)
	}

	type idMeta struct {
		IsNullable    string `gorm:"column:is_nullable"`
		ColumnDefault string `gorm:"column:column_default"`
		IsIdentity    string `gorm:"column:is_identity"`
	}
	var meta idMeta
	_ = repo.model(ctx).Raw(`
		SELECT
			COALESCE(is_nullable, 'YES') AS is_nullable,
			COALESCE(column_default, '') AS column_default,
			COALESCE(is_identity, 'NO') AS is_identity
		FROM information_schema.columns
		WHERE table_schema = 'inv'
		  AND table_name = ?
		  AND column_name = 'id'
		LIMIT 1
	`, tableName).Scan(&meta).Error

	needsExplicitID := strings.EqualFold(meta.IsNullable, "NO") &&
		strings.TrimSpace(meta.ColumnDefault) == "" &&
		!strings.EqualFold(meta.IsIdentity, "YES")

	var sql string
	var args []interface{}
	if needsExplicitID {
		sql = fmt.Sprintf(`
			INSERT INTO inv.%s
			(id, cust_id, %s, level, sequence, pic, status, created_at, remarks)
			VALUES (
				(
					SELECT COALESCE(MAX(id), 0) + 1
					FROM inv.%s
					WHERE cust_id = ?
				),
				?, ?, ?, ?, ?, 1, ?, NULL
			)
		`, tableName, orderIDColumn, tableName)
		args = []interface{}{custID, custID, replenishmentID, level, sequence, pic, now}
	} else {
		sql = fmt.Sprintf(`
			INSERT INTO inv.%s
			(cust_id, %s, level, sequence, pic, status, created_at, remarks)
			VALUES (?, ?, ?, ?, ?, 1, ?, NULL)
		`, tableName, orderIDColumn)
		args = []interface{}{custID, replenishmentID, level, sequence, pic, now}
	}

	if err := repo.model(ctx).Exec(sql, args...).Error; err != nil {
		return fmt.Errorf("failed to insert replenishment approval queue on inv.%s: %w", tableName, err)
	}
	return nil
}

func (repo *RepositoryReplenishmentImpl) IsUserApprovalPIC(ctx context.Context, userID int64) (bool, error) {
	var count int64
	err := repo.model(ctx).
		Table("mst.distributor_replenishment_approval").
		Where("pic = ? AND is_del = false", userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *RepositoryReplenishmentImpl) UpdateReplenishmentOrderApprovalStatus(ctx context.Context, custID string, replenishmentID int64, pic int64, status int, remarks *string) error {
	var tableName string
	err := repo.model(ctx).Raw(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'inv'
		  AND table_name IN ('replenishment_order_approval', 'replanishment_order_approval')
		ORDER BY CASE
			WHEN table_name = 'replenishment_order_approval' THEN 1
			ELSE 2
		END
		LIMIT 1
	`).Scan(&tableName).Error
	if err != nil {
		return fmt.Errorf("failed to resolve approval queue table: %w", err)
	}
	if strings.TrimSpace(tableName) == "" {
		return fmt.Errorf("approval queue table not found in schema inv")
	}

	var orderIDColumn string
	err = repo.model(ctx).Raw(`
		SELECT column_name
		FROM information_schema.columns
		WHERE table_schema = 'inv'
		  AND table_name = ?
		  AND column_name IN ('replenishment_order_id', 'replanishment_order_id')
		ORDER BY CASE
			WHEN column_name = 'replenishment_order_id' THEN 1
			ELSE 2
		END
		LIMIT 1
	`, tableName).Scan(&orderIDColumn).Error
	if err != nil {
		return fmt.Errorf("failed to resolve approval queue order id column: %w", err)
	}
	if strings.TrimSpace(orderIDColumn) == "" {
		return fmt.Errorf("approval queue order id column not found on inv.%s", tableName)
	}

	var remarksVal interface{}
	if remarks != nil && strings.TrimSpace(*remarks) != "" {
		remarksVal = strings.TrimSpace(*remarks)
	} else {
		remarksVal = nil
	}

	updateSQL := fmt.Sprintf(`
		UPDATE inv.%s
		SET status = ?, remarks = ?
		WHERE cust_id = ? AND %s = ? AND pic = ?
	`, tableName, orderIDColumn)
	result := repo.model(ctx).Exec(updateSQL, status, remarksVal, custID, replenishmentID, pic)
	if result.Error != nil {
		return fmt.Errorf("failed to update approval queue status: %w", result.Error)
	}

	// Fallback update if PIC specific row was not found.
	if result.RowsAffected == 0 {
		fallbackSQL := fmt.Sprintf(`
			UPDATE inv.%s
			SET status = ?, remarks = ?
			WHERE cust_id = ? AND %s = ?
		`, tableName, orderIDColumn)
		result = repo.model(ctx).Exec(fallbackSQL, status, remarksVal, custID, replenishmentID)
		if result.Error != nil {
			return fmt.Errorf("failed to update approval queue status (fallback): %w", result.Error)
		}
	}

	return nil
}

func (repo *RepositoryReplenishmentImpl) IsReturnReasonDistributorExists(ctx context.Context, returnReasonID int64) (bool, error) {
	var count int64
	err := repo.model(ctx).
		Table("mst.m_return_reason_distributor").
		Where("return_reason_id = ? AND (is_del = false OR is_del IS NULL)", returnReasonID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *RepositoryReplenishmentImpl) FindSummarizeReplanishment(ctx context.Context, replanishmentIDs []int64, custID, parentCustID string, isPrincipal bool, userID int64) ([]model.SummarizeReplanishmentRow, error) {
	rows := make([]model.SummarizeReplanishmentRow, 0)
	if len(replanishmentIDs) == 0 {
		return rows, nil
	}

	distributorExpr := "COALESCE(ro.distributor_id, mc.distributor_id, 0)"
	ripeningDistExpr := "COALESCE(NULLIF(ro.distributor_id, 0), NULLIF(mc.distributor_id, 0))"
	whsTotalQtySubquery := `COALESCE((
		SELECT SUM(COALESCE(whs_sub.qty, 0))
		FROM inv.warehouse_stock whs_sub
		WHERE whs_sub.pro_id = rod.pro_id
		  AND whs_sub.wh_id = ro.wh_id
		  AND whs_sub.cust_id = rod.cust_id
	), 0)`
	qty1Expr, qty2Expr, qty3Expr := BuildQtyCalculationExpressions(whsTotalQtySubquery, "product", "wh_stock_small", "wh_stock_medium", "wh_stock_large")

	query := repo.model(ctx).
		Table("inv.replenishment_order ro").
		Select(`
			ro.replenishment_id AS replanishment_id,
			ro.replenishment_no AS replanishment_no,
			`+distributorExpr+` AS disributor_id,
			md.distributor_code,
			md.distributor_name,
			ro.sup_id,
			sup.sup_code,
			sup.sup_name,
			ro.wh_id,
			wh.wh_code,
			wh.wh_name,
			rod.replenishment_detail_id AS replanishment_detail_id,
			rod.pro_id,
			COALESCE(product.pro_code, '') AS pro_code,
			COALESCE(product.pro_name, '') AS pro_name,
			`+qty3Expr+`,
			`+qty2Expr+`,
			`+qty1Expr+`,
			0::numeric AS optimum_qty,
			COALESCE(rp.ripening_qty, 0) AS ripening,
			rod.return_reason_id AS return_reason_id,
			COALESCE(product.unit_id1, '') AS unit_id1,
			COALESCE(product.unit_id2, '') AS unit_id2,
			COALESCE(product.unit_id3, '') AS unit_id3,
			rod.purch_price1,
			rod.purch_price2,
			rod.purch_price3,
			rod.order_booking_qty1 AS qty_ro1,
			rod.order_booking_qty2 AS qty_ro2,
			rod.order_booking_qty3 AS qty_ro3,
			(rod.order_booking_qty1 * rod.purch_price1) +
			(rod.order_booking_qty2 * rod.purch_price2) +
			(rod.order_booking_qty3 * rod.purch_price3) AS estimated_price`).
		Joins("JOIN inv.replenishment_order_detail rod ON rod.replenishment_id = ro.replenishment_id AND rod.cust_id = ro.cust_id AND rod.is_del = false").
		Joins("LEFT JOIN smc.m_customer mc ON mc.cust_id = ro.cust_id").
		Joins("LEFT JOIN mst.m_distributor md ON md.distributor_id = "+distributorExpr).
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = ro.sup_id").
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = ro.wh_id AND wh.cust_id = ro.cust_id").
		Joins(`LEFT JOIN LATERAL (
			SELECT pro_code, pro_name, unit_id1, unit_id2, unit_id3, conv_unit2, conv_unit3
			FROM mst.m_product
			WHERE pro_id = rod.pro_id AND is_del = false
			ORDER BY CASE WHEN cust_id = ? THEN 1 WHEN cust_id = ? THEN 2 ELSE 3 END
			LIMIT 1
		) AS product ON true`, parentCustID, custID).
		Joins(`LEFT JOIN LATERAL (
			SELECT COALESCE(MAX(
				CASE EXTRACT(DOW FROM ((CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'))::date)
					WHEN 1 THEN pr.monday_qty
					WHEN 2 THEN pr.tuesday_qty
					WHEN 3 THEN pr.wednesday_qty
					WHEN 4 THEN pr.thursday_qty
					WHEN 5 THEN pr.friday_qty
					WHEN 6 THEN pr.saturday_qty
					WHEN 0 THEN pr.sunday_qty
				END
			), 0) AS ripening_qty
			FROM mst.product_ripening pr
			WHERE pr.pro_id = rod.pro_id
			  AND pr.created_by = ?
			  AND `+ripeningDistExpr+` IS NOT NULL
			  AND `+ripeningDistExpr+` > 0
			  AND pr.distributor_id = `+ripeningDistExpr+`
		) rp ON true`, userID).
		Where("ro.replenishment_id IN ?", replanishmentIDs).
		Where("ro.is_del = false")

	if isPrincipal {
		query = query.Where("(md.parent_cust_id = ? OR ro.cust_id = ?)", custID, custID)
	} else {
		query = query.Where("ro.cust_id = ?", custID)
	}

	err := query.Order("ro.replenishment_id ASC, rod.replenishment_detail_id ASC").Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (repo *RepositoryReplenishmentImpl) FindSAPExportRows(ctx context.Context, distributorCode string, dateFrom, dateTo int64) ([]model.SAPReplExportRow, error) {
	distributorCode = strings.TrimSpace(distributorCode)
	if distributorCode == "" {
		return []model.SAPReplExportRow{}, nil
	}

	startTime := time.Unix(dateFrom, 0).In(constant.AsiaJakartaLocation)
	endTime := time.Unix(dateTo, 0).In(constant.AsiaJakartaLocation)

	distributorExpr := "COALESCE(ro.distributor_id, mc.distributor_id, 0)"
	distributorJoin := fmt.Sprintf(`JOIN mst.m_distributor md ON (
		(%[1]s > 0 AND md.distributor_id = %[1]s)
		OR (%[1]s = 0 AND md.cust_id = ro.cust_id)
	) AND md.is_del = false`, distributorExpr)

	rows := make([]model.SAPReplExportRow, 0)
	err := repo.model(ctx).
		Table("inv.replenishment_order AS ro").
		Select(`
			ro.replenishment_no,
			TRIM(md.distributor_code) AS ship_to_party_code,
			ro.replenishment_type,
			COALESCE(ro.is_addition_from, true) AS is_addition_from,
			ro.status,
			COALESCE(ch.channel_code, '') AS distribution_channel,
			COALESCE(ar.area_code, '') AS sales_office,
			ro.delivery_date,
			COALESCE(md.barcode, '') AS shipping_point,
			COALESCE(grp.sub_distributor_group_code, '') AS plant,
			ro.created_at,
			COALESCE(product.pro_code, '') AS material,
			COALESCE(rod.qty_order_approval3, rod.order_booking_qty3, 0) AS order_qty,
			COALESCE(product.unit_id3, '') AS uom,
			rod.purch_price3 AS cust_po_price,
			COALESCE(product.pl_code, '') AS division`).
		Joins("JOIN smc.m_customer mc ON mc.cust_id = ro.cust_id").
		Joins(distributorJoin).
		Joins("JOIN inv.replenishment_order_detail rod ON rod.replenishment_id = ro.replenishment_id AND rod.cust_id = ro.cust_id AND rod.is_del = false").
		Joins(`LEFT JOIN mst.m_channel ch ON ch.channel_id = md.channel_id`).
		Joins(`LEFT JOIN mst.m_area ar ON ar.area_id = md.area_id`).
		Joins(`LEFT JOIN mst.m_sub_distributor_group grp ON grp.sub_distributor_group_id = md.sub_distributor_group_id AND grp.is_del = false`).
		Joins(`LEFT JOIN LATERAL (
			SELECT p.pro_code, p.unit_id3, pl.pl_code
			FROM mst.m_product p
			LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id
			LEFT JOIN mst.m_brand b ON b.brand_id = sb1.brand_id
			LEFT JOIN mst.m_product_line pl ON pl.pl_id = b.pl_id
			WHERE p.pro_id = rod.pro_id AND p.is_del = false
			ORDER BY CASE
				WHEN p.cust_id = COALESCE(NULLIF(TRIM(mc.parent_cust_id), ''), ro.cust_id) THEN 1
				WHEN p.cust_id = ro.cust_id THEN 2
				ELSE 3
			END
			LIMIT 1
		) AS product ON true`).
		Where("ro.is_del = ?", false).
		Where("TRIM(md.distributor_code) = ?", distributorCode).
		Where("ro.created_at >= ? AND ro.created_at <= ?", startTime, endTime).
		Order("ro.replenishment_no ASC, rod.replenishment_detail_id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}
