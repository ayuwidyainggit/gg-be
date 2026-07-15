package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sales/entity"
	"sales/model"
	"sales/pkg/str"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryOrderApprovalImpl struct {
		*gorm.DB
	}
)
type OrderApprovalRepository interface {
	FindNeedReview(filter entity.OrderApprovalQueryFilter) ([]model.OrderApprovalRead, int64, int, error)
	FindApproved(filter entity.OrderApprovalQueryFilter) ([]model.OrderApprovalRead, int64, int, error)
	FindRejected(filter entity.OrderApprovalQueryFilter) ([]model.OrderApprovalRead, int64, int, error)
	UpdateStatusDetail(c context.Context, OrderApprovalRequestID int64, empID int64, status int) error
	FindIfNeedsApproval(c context.Context, orderApprovalRequestID int64) (result model.OrderApprovalActiveRead, err error)
	UpdateOrderApprovalFinished(c context.Context, OrderApprovalRequestID int64) error
	UpdateStatusOrder(c context.Context, RoNo string, status int) error
	FindOrderApprovalByID(orderApprovalRequestID int64) (model.OrderApprovalRequestRead, error)
}

func NewOrderApprovalRepo(db *gorm.DB) *RepositoryOrderApprovalImpl {
	return &RepositoryOrderApprovalImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryOrderApprovalImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryOrderApprovalImpl) FindNeedReview(filter entity.OrderApprovalQueryFilter) ([]model.OrderApprovalRead, int64, int, error) {
	var results []model.OrderApprovalRead
	var total int64
	var limit int
	if filter.Limit == 0 {
		limit = 10
	} else {
		limit = filter.Limit
	}
	query := repository.Table("sls.order_approval_requests oar").
		Select(`
		oard.order_approval_request_id, ord.ro_no, ord.ro_date, o.outlet_id, o.outlet_code, o.outlet_name, 
		ord.validate_credit_limit_value, ord.validate_overdue_value, ord.validate_outstanding_value,
		ord.total, s.sales_name, o.credit_limit, o.sales_inv_limit,
		CASE WHEN o.credit_limit_type = 2 THEN 2 ELSE NULL END AS credit_limit_type,
		CASE 
			WHEN o.credit_limit_type = 1 THEN 'Unlimited'
			WHEN o.credit_limit_type = 2 THEN 'Limit By Total'
			WHEN o.credit_limit_type = 3 THEN 'Limit By Supplier'
			ELSE 'Unlimited'
		END AS credit_limit_type_name,
		CASE WHEN o.credit_limit_action = 1 THEN 1
			 WHEN o.credit_limit_action = 2 THEN 2
			 ELSE NULL
		END AS credit_limit_action,
		CASE 
			WHEN o.credit_limit_action = 1 THEN 'Warning'
			WHEN o.credit_limit_action = 2 THEN 'Restricted'
			ELSE ''
		END AS credit_limit_action_name,
		CASE WHEN o.sales_inv_limit_type = 2 THEN 2 ELSE NULL END AS sales_inv_limit_type,
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
			WHEN o.sales_inv_limit_action = 1 THEN 'Warning'
			WHEN o.sales_inv_limit_action = 2 THEN 'Restricted'
			ELSE ''
		END AS sales_inv_limit_action_name,
		CASE WHEN o.obs_type = 2 THEN 2 ELSE NULL END AS obs_type,
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
			WHEN o.obs_limit_action = 1 THEN 'Warning'
			WHEN o.obs_limit_action = 2 THEN 'Restricted'
			ELSE ''
		END AS obs_limit_action_name,
		oar.cust_id AS cust_id_origin
	`).
		Joins("JOIN sls.order_approval_requests_details oard ON oar.order_approval_request_id = oard.order_approval_request_id").
		Joins("LEFT JOIN sls.order ord ON oar.ro_no = ord.ro_no AND ord.cust_id = oar.cust_id").
		Joins("LEFT JOIN mst.m_outlet o ON ord.outlet_id = o.outlet_id AND o.cust_id = oar.cust_id").
		Joins("LEFT JOIN mst.m_salesman s ON ord.salesman_id = s.emp_id AND s.cust_id = oar.cust_id")

	queryCount := repository.Table("sls.order_approval_requests oar").Select("ord.ro_no").
		Joins("JOIN sls.order_approval_requests_details oard ON oar.order_approval_request_id = oard.order_approval_request_id").
		Joins("LEFT JOIN sls.order ord ON oar.ro_no = ord.ro_no AND ord.cust_id = oar.cust_id").
		Joins("LEFT JOIN mst.m_outlet o ON ord.outlet_id = o.outlet_id AND o.cust_id = oar.cust_id").
		Joins("LEFT JOIN mst.m_salesman s ON ord.salesman_id = s.emp_id AND s.cust_id = oar.cust_id")

	query.Where("oard.emp_id = ?", filter.EmpID).
		Where("oard.status IS NULL").
		Where(`
		NOT EXISTS (
			SELECT 1 FROM sls.order_approval_requests_details rj
			WHERE rj.order_approval_request_id = oar.order_approval_request_id
			AND rj.status = 9
			AND rj.level = (
				SELECT MAX(level) FROM sls.order_approval_requests_details
				WHERE order_approval_request_id = oar.order_approval_request_id
			)
		)
	`).
		Where(`
		NOT EXISTS (
			SELECT 1 FROM sls.order_approval_requests_details h
			WHERE h.order_approval_request_id = oar.order_approval_request_id
			AND h.level > oard.level
			AND h.status IS NULL
		)
	`)

	queryCount.Where("oard.emp_id = ?", filter.EmpID).
		Where("oard.status IS NULL").
		Where(`
	NOT EXISTS (
		SELECT 1 FROM sls.order_approval_requests_details rj
		WHERE rj.order_approval_request_id = oar.order_approval_request_id
		AND rj.status = 9
		AND rj.level = (
			SELECT MAX(level) FROM sls.order_approval_requests_details
			WHERE order_approval_request_id = oar.order_approval_request_id
		)
	)
`).
		Where(`
	NOT EXISTS (
		SELECT 1 FROM sls.order_approval_requests_details h
		WHERE h.order_approval_request_id = oar.order_approval_request_id
		AND h.level > oard.level
		AND h.status IS NULL
	)
`)

	// Optional filter: ro_date
	if filter.RoFrom != nil && filter.RoTo != nil {
		fromTime := str.UnixTimestampToUtcTime(*filter.RoFrom)
		toTime := str.UnixTimestampToUtcTime(*filter.RoTo)
		query.Where("ord.ro_date BETWEEN ? AND ?", fromTime, toTime)
		queryCount.Where("ord.ro_date BETWEEN ? AND ?", fromTime, toTime)
	}

	// Optional filter: salesman_id
	if len(filter.SalesmanId) > 0 {
		query.Where("s.emp_id IN ?", filter.SalesmanId)
		queryCount.Where("s.emp_id IN ?", filter.SalesmanId)

	}

	// Optional filter: outlet_id
	if len(filter.OutletID) > 0 {
		query.Where("o.outlet_id IN ?", filter.OutletID)
		queryCount.Where("o.outlet_id IN ?", filter.OutletID)

	}

	sortBy := ``
	if filter.Sort != "" {
		mSortBy := strings.Split(filter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("ord.ro_no DESC")
	}

	page := filter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * filter.Limit

	err := query.Limit(limit).Offset(offset).Find(&results).Error
	if err != nil {
		return results, total, 0, err
	}
	err = queryCount.Model(&results).Count(&total).Error
	if err != nil {
		return results, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))

	return results, total, lastPage, nil
}

func (repository *RepositoryOrderApprovalImpl) FindApproved(filter entity.OrderApprovalQueryFilter) ([]model.OrderApprovalRead, int64, int, error) {
	var results []model.OrderApprovalRead
	var total int64
	var limit int

	if filter.Limit == 0 {
		limit = 10
	} else {
		limit = filter.Limit
	}

	query := repository.Table("sls.order_approval_requests oar").
		Select(`
		oard.order_approval_request_id, ord.ro_no, ord.ro_date, o.outlet_id, o.outlet_code, o.outlet_name,
		ord.validate_credit_limit_value, ord.validate_overdue_value, ord.validate_outstanding_value,
		ord.total, s.sales_name, o.credit_limit, o.sales_inv_limit,
		CASE WHEN o.credit_limit_type = 2 THEN 2 ELSE NULL END AS credit_limit_type,
		CASE 
			WHEN o.credit_limit_type = 1 THEN 'Unlimited'
			WHEN o.credit_limit_type = 2 THEN 'Limit By Total'
			WHEN o.credit_limit_type = 3 THEN 'Limit By Supplier'
			ELSE 'Unlimited'
		END AS credit_limit_type_name,
		CASE WHEN o.credit_limit_action = 1 THEN 1
			 WHEN o.credit_limit_action = 2 THEN 2
			 ELSE NULL
		END AS credit_limit_action,
		CASE 
			WHEN o.credit_limit_action = 1 THEN 'Warning'
			WHEN o.credit_limit_action = 2 THEN 'Restricted'
			ELSE ''
		END AS credit_limit_action_name,
		CASE WHEN o.sales_inv_limit_type = 2 THEN 2 ELSE NULL END AS sales_inv_limit_type,
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
			WHEN o.sales_inv_limit_action = 1 THEN 'Warning'
			WHEN o.sales_inv_limit_action = 2 THEN 'Restricted'
			ELSE ''
		END AS sales_inv_limit_action_name,
		CASE WHEN o.obs_type = 2 THEN 2 ELSE NULL END AS obs_type,
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
			WHEN o.obs_limit_action = 1 THEN 'Warning'
			WHEN o.obs_limit_action = 2 THEN 'Restricted'
			ELSE ''
		END AS obs_limit_action_name,
		oar.cust_id AS cust_id_origin
	`).
		Joins("JOIN sls.order_approval_requests_details oard ON oar.order_approval_request_id = oard.order_approval_request_id").
		Joins("LEFT JOIN sls.order ord ON oar.ro_no = ord.ro_no AND ord.cust_id = oar.cust_id").
		Joins("LEFT JOIN mst.m_outlet o ON ord.outlet_id = o.outlet_id AND o.cust_id = oar.cust_id").
		Joins("LEFT JOIN mst.m_salesman s ON ord.salesman_id = s.emp_id AND s.cust_id = oar.cust_id").
		Where("oard.emp_id = ?", filter.EmpID).
		Where("oard.status = ?", 1)

	// Filter tanggal RO
	if filter.RoFrom != nil && filter.RoTo != nil {
		fromTime := str.UnixTimestampToUtcTime(*filter.RoFrom)
		toTime := str.UnixTimestampToUtcTime(*filter.RoTo)
		query.Where("ord.ro_date BETWEEN ? AND ?", fromTime, toTime)
	}

	// Filter salesman
	if len(filter.SalesmanId) > 0 {
		query.Where("s.emp_id IN ?", filter.SalesmanId)
	}

	// Filter outlet
	if len(filter.OutletID) > 0 {
		query.Where("o.outlet_id IN ?", filter.OutletID)
	}

	sortBy := ``
	if filter.Sort != "" {
		mSortBy := strings.Split(filter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("ord.ro_no DESC")
	}

	// Pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	err := query.Limit(limit).Offset(offset).Find(&results).Error
	if err != nil {
		return results, total, 0, err
	}

	// Count total
	queryCount := repository.Table("sls.order_approval_requests oar").
		Select("COUNT(DISTINCT ord.ro_no)").
		Joins("JOIN sls.order_approval_requests_details oard ON oar.order_approval_request_id = oard.order_approval_request_id").
		Joins("LEFT JOIN sls.order ord ON oar.ro_no = ord.ro_no AND ord.cust_id = oar.cust_id").
		Joins("LEFT JOIN mst.m_outlet o ON ord.outlet_id = o.outlet_id AND o.cust_id = oar.cust_id").
		Joins("LEFT JOIN mst.m_salesman s ON ord.salesman_id = s.emp_id AND s.cust_id = oar.cust_id").
		Where("oard.emp_id = ?", filter.EmpID).
		Where("oard.status = ?", 1)

	if filter.RoFrom != nil && filter.RoTo != nil {
		fromTime := str.UnixTimestampToUtcTime(*filter.RoFrom)
		toTime := str.UnixTimestampToUtcTime(*filter.RoTo)
		queryCount.Where("ord.ro_date BETWEEN ? AND ?", fromTime, toTime)
	}

	if len(filter.SalesmanId) > 0 {
		queryCount.Where("s.emp_id IN ?", filter.SalesmanId)
	}

	if len(filter.OutletID) > 0 {
		queryCount.Where("o.outlet_id IN ?", filter.OutletID)
	}

	err = queryCount.Count(&total).Error
	if err != nil {
		return results, total, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return results, total, lastPage, nil
}

func (repository *RepositoryOrderApprovalImpl) FindRejected(filter entity.OrderApprovalQueryFilter) ([]model.OrderApprovalRead, int64, int, error) {
	var results []model.OrderApprovalRead
	var total int64
	var limit int

	if filter.Limit == 0 {
		limit = 10
	} else {
		limit = filter.Limit
	}

	query := repository.Table("sls.order_approval_requests oar").
		Select(`
		oard.order_approval_request_id, ord.ro_no, ord.ro_date, o.outlet_id, o.outlet_code, o.outlet_name,
		ord.validate_credit_limit_value, ord.validate_overdue_value, ord.validate_outstanding_value,
		ord.total, s.sales_name, o.credit_limit, o.sales_inv_limit,
		CASE WHEN o.credit_limit_type = 2 THEN 2 ELSE NULL END AS credit_limit_type,
		CASE 
			WHEN o.credit_limit_type = 1 THEN 'Unlimited'
			WHEN o.credit_limit_type = 2 THEN 'Limit By Total'
			WHEN o.credit_limit_type = 3 THEN 'Limit By Supplier'
			ELSE 'Unlimited'
		END AS credit_limit_type_name,
		CASE WHEN o.credit_limit_action = 1 THEN 1
			 WHEN o.credit_limit_action = 2 THEN 2
			 ELSE NULL
		END AS credit_limit_action,
		CASE 
			WHEN o.credit_limit_action = 1 THEN 'Warning'
			WHEN o.credit_limit_action = 2 THEN 'Restricted'
			ELSE ''
		END AS credit_limit_action_name,
		CASE WHEN o.sales_inv_limit_type = 2 THEN 2 ELSE NULL END AS sales_inv_limit_type,
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
			WHEN o.sales_inv_limit_action = 1 THEN 'Warning'
			WHEN o.sales_inv_limit_action = 2 THEN 'Restricted'
			ELSE ''
		END AS sales_inv_limit_action_name,
		CASE WHEN o.obs_type = 2 THEN 2 ELSE NULL END AS obs_type,
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
			WHEN o.obs_limit_action = 1 THEN 'Warning'
			WHEN o.obs_limit_action = 2 THEN 'Restricted'
			ELSE ''
		END AS obs_limit_action_name,
		oar.cust_id AS cust_id_origin
	`).
		Joins("JOIN sls.order_approval_requests_details oard ON oar.order_approval_request_id = oard.order_approval_request_id").
		Joins("LEFT JOIN sls.order ord ON oar.ro_no = ord.ro_no AND ord.cust_id = oar.cust_id").
		Joins("LEFT JOIN mst.m_outlet o ON ord.outlet_id = o.outlet_id AND o.cust_id = oar.cust_id").
		Joins("LEFT JOIN mst.m_salesman s ON ord.salesman_id = s.emp_id AND s.cust_id = oar.cust_id").
		Where("oard.emp_id = ?", filter.EmpID).
		Where("oard.status = ?", 9)

	// Filter tanggal RO
	if filter.RoFrom != nil && filter.RoTo != nil {
		fromTime := str.UnixTimestampToUtcTime(*filter.RoFrom)
		toTime := str.UnixTimestampToUtcTime(*filter.RoTo)
		query.Where("ord.ro_date BETWEEN ? AND ?", fromTime, toTime)
	}

	// Filter salesman
	if len(filter.SalesmanId) > 0 {
		query.Where("s.emp_id IN ?", filter.SalesmanId)
	}

	// Filter outlet
	if len(filter.OutletID) > 0 {
		query.Where("o.outlet_id IN ?", filter.OutletID)
	}

	sortBy := ``
	if filter.Sort != "" {
		mSortBy := strings.Split(filter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("ord.ro_no DESC")
	}

	// Pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	err := query.Limit(limit).Offset(offset).Find(&results).Error
	if err != nil {
		return results, total, 0, err
	}

	// Count total
	queryCount := repository.Table("sls.order_approval_requests oar").
		Select("COUNT(DISTINCT ord.ro_no)").
		Joins("JOIN sls.order_approval_requests_details oard ON oar.order_approval_request_id = oard.order_approval_request_id").
		Joins("LEFT JOIN sls.order ord ON oar.ro_no = ord.ro_no AND ord.cust_id = oar.cust_id").
		Joins("LEFT JOIN mst.m_outlet o ON ord.outlet_id = o.outlet_id AND o.cust_id = oar.cust_id").
		Joins("LEFT JOIN mst.m_salesman s ON ord.salesman_id = s.emp_id AND s.cust_id = oar.cust_id").
		Where("oard.emp_id = ?", filter.EmpID).
		Where("oard.status = ?", 9)

	if filter.RoFrom != nil && filter.RoTo != nil {
		fromTime := str.UnixTimestampToUtcTime(*filter.RoFrom)
		toTime := str.UnixTimestampToUtcTime(*filter.RoTo)
		queryCount.Where("ord.ro_date BETWEEN ? AND ?", fromTime, toTime)
	}

	if len(filter.SalesmanId) > 0 {
		queryCount.Where("s.emp_id IN ?", filter.SalesmanId)
	}

	if len(filter.OutletID) > 0 {
		queryCount.Where("o.outlet_id IN ?", filter.OutletID)
	}

	err = queryCount.Count(&total).Error
	if err != nil {
		return results, total, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return results, total, lastPage, nil
}

func (repository *RepositoryOrderApprovalImpl) UpdateStatusDetail(c context.Context, OrderApprovalRequestID int64, empID int64, status int) error {
	result := repository.model(c).Table("sls.order_approval_requests_details oard").Where("order_approval_request_id = ? AND emp_id =? ", OrderApprovalRequestID, empID).
		Updates(map[string]interface{}{"status": status, "act_date": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryOrderApprovalImpl) UpdateOrderApprovalFinished(c context.Context, OrderApprovalRequestID int64) error {
	result := repository.model(c).Table("sls.order_approval_requests oar").Where("order_approval_request_id = ? ", OrderApprovalRequestID).
		Updates(map[string]interface{}{"finished_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryOrderApprovalImpl) FindIfNeedsApproval(c context.Context, orderApprovalRequestID int64) (result model.OrderApprovalActiveRead, err error) {
	// Subquery: ambil level paling atas (paling kecil) yang masih NULL
	subQueryMinLevel := repository.model(c).
		Table("sls.order_approval_requests_details").
		Select("MIN(level)").
		Where("order_approval_request_id = ?", orderApprovalRequestID).
		Where("status IS NULL")

	err = repository.model(c).
		Table("sls.order_approval_requests AS oar").
		Where("oar.order_approval_request_id = ?", orderApprovalRequestID).

		// ❌ Tidak boleh ada status = 9 (rejected)
		Where("NOT EXISTS (?)",
			repository.model(c).
				Table("sls.order_approval_requests_details AS d").
				Select("1").
				Where("d.order_approval_request_id = oar.order_approval_request_id").
				Where("d.status = 9"),
		).

		// ✅ Masih ada level yang butuh approval (status NULL di level aktif)
		Where("EXISTS (?)",
			repository.model(c).
				Table("sls.order_approval_requests_details AS d1").
				Select("1").
				Where("d1.order_approval_request_id = oar.order_approval_request_id").
				Where("d1.level = (?)", subQueryMinLevel).
				Where("d1.status IS NULL"),
		).

		// ❌ Belum ada yang approve (status=1) di level aktif
		Where("NOT EXISTS (?)",
			repository.model(c).
				Table("sls.order_approval_requests_details AS d3").
				Select("1").
				Where("d3.order_approval_request_id = oar.order_approval_request_id").
				Where("d3.level = (?)", subQueryMinLevel).
				Where("d3.status = 1"),
		).
		First(&result).Error

	return
}

func (repository *RepositoryOrderApprovalImpl) UpdateStatusOrder(c context.Context, RoNo string, status int) error {
	var data model.Order
	result := repository.model(c).Model(&data).Where("ro_no=?", RoNo).
		Updates(map[string]interface{}{"data_status": status})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (repository *RepositoryOrderApprovalImpl) FindOrderApprovalByID(orderApprovalRequestID int64) (model.OrderApprovalRequestRead, error) {
	orderApprovalRequest := model.OrderApprovalRequestRead{}
	err := repository.Select("*").Where("order_approval_request_id = ?", orderApprovalRequestID).Take(&orderApprovalRequest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return orderApprovalRequest, errors.New(fmt.Sprintf("order approval request id : %v not found", orderApprovalRequestID))
		}

		return orderApprovalRequest, err
	}
	return orderApprovalRequest, nil
}
