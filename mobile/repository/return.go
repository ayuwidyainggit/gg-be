package repository

import (
	"context"
	"fmt"
	"mobile/entity"
	"mobile/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryReturnImpl struct {
		*gorm.DB
	}
)
type ReturnRepository interface {
	Store(c context.Context, data *model.Return) error
	StoreDetail(c context.Context, data *model.ReturnDetail) error
	Update(c context.Context, data *model.Return) error
	FindAllMasterReturnReasonLookupMode(dataFilter entity.GeneralQueryFilter) (outlets []model.ReturnReasonLookup, total int64, lastPage int, err error)
	DeleteDetailNotInIDs(c context.Context, returnNo string, IDs []int64) error
	UpdateQuantity(c context.Context, Details *model.ReturnQuantity) error
	FindDistinctReturnsByProductIdAndDate(productId int64) ([]model.ReturnInfo, error)
	FindInvoiceNoByProductId(productId int64) ([]model.ReturnInfo, error)
	FindOrderDetailIdByInvoiceNo(invoiceNo int64) ([]model.ReturnDetail, error)
	FindOneWhIdBySalesmanID(salesmanID int64) (int64, error)
}

func NewReturnRepo(db *gorm.DB) *RepositoryReturnImpl {
	return &RepositoryReturnImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryReturnImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryReturnImpl) Store(c context.Context, data *model.Return) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryReturnImpl) StoreDetail(c context.Context, data *model.ReturnDetail) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryReturnImpl) Update(c context.Context, data *model.Return) error {
	result := repository.model(c).Model(&data).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryReturnImpl) FindAllMasterReturnReasonLookupMode(dataFilter entity.GeneralQueryFilter) ([]model.ReturnReasonLookup, int64, int, error) {

	var return_reasons []model.ReturnReasonLookup
	var total int64

	query := repository.Select("mst.m_return_reason.return_reason_id, mst.m_return_reason.return_reason_name")

	query.Where("mst.m_return_reason.cust_id=?", dataFilter.ParentCustId)

	query.Where("mst.m_return_reason.is_active=?", true)

	query.Where("mst.m_return_reason.is_del=?", false)

	if dataFilter.Query != "" {
		query.Where("mst.m_return_reason.return_reason_name LIKE ?", "%"+dataFilter.Query+"%")
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
		query.Order(sortBy)
	} else {
		query.Order("return_reason_id DESC")
	}

	err := query.Find(&return_reasons).Error
	if err != nil {
		return return_reasons, total, 0, err
	}

	total = int64(len(return_reasons))
	lastPage := 1
	return return_reasons, total, lastPage, nil
}

func (repository *RepositoryReturnImpl) DeleteDetailNotInIDs(c context.Context, returnNo string, IDs []int64) error {
	var Details model.ReturnDetail
	err := repository.model(c).Where("return_no=? AND return_detail_id not in (?) ", returnNo, IDs).Delete(&Details).Error
	return err
}

func (repository *RepositoryReturnImpl) UpdateQuantity(c context.Context, Details *model.ReturnQuantity) error {
	// result := repository.model(c).Updates(&Details)
	result := repository.model(c).Select("Qty1", "Qty2", "Qty3").Updates(&Details)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryReturnImpl) FindDistinctReturnsByProductIdAndDate(productId int64) ([]model.ReturnInfo, error) {
	var returns []model.ReturnInfo

	// Calculate the date one month ago
	OneMonthAgo := time.Now().AddDate(0, -1, 0).Format("2006-01-02")

	// Define the query to select distinct return_no, invoice_no, invoice_date, and product_id
	err := repository.Select("DISTINCT r.return_no, r.invoice_no, r.invoice_date, rd.product_id").
		Table("(SELECT * FROM sls.return WHERE invoice_date >= ?) AS r", OneMonthAgo).
		Joins("JOIN sls.return_det AS rd ON r.return_no = rd.return_no").
		Where("rd.product_id = ? ", productId).
		Order("r.invoice_date DESC").
		Find(&returns).Error

	if err != nil {
		return nil, err
	}

	return returns, nil
}

func (repository *RepositoryReturnImpl) FindInvoiceNoByProductId(productId int64) ([]model.ReturnInfo, error) {
	var returns []model.ReturnInfo

	// Calculate the date one month ago
	OneMonthAgo := time.Now().AddDate(0, -1, 0).Format("2006-01-02")

	// Define the query to select distinct return_no, invoice_no, invoice_date, and product_id
	err := repository.Select("DISTINCT o.invoice_no, o.invoice_date, od.pro_id, od.order_detail_id, o.wh_id").
		Table("(SELECT * FROM sls.order WHERE invoice_date >= ?) AS o", OneMonthAgo).
		Joins("JOIN sls.order_detail AS od ON o.ro_no = od.ro_no").
		Where("od.pro_id = ?", productId).
		Order("o.invoice_date DESC").
		Find(&returns).Error

	if err != nil {
		return nil, err
	}

	return returns, nil
}

func (repository *RepositoryReturnImpl) FindOneWhIdBySalesmanID(salesmanID int64) (int64, error) {
	var returns model.SalesmanWhId

	// Calculate the date one month ago

	// Define the query to select distinct return_no, invoice_no, invoice_date, and product_id
	err := repository.Select("wh_id").
		Table("mst.m_salesman").
		Where(" emp_id  = ?", salesmanID).
		Order("o.invoice_date DESC").
		Find(&returns).Error

	if err != nil {
		return 0, err
	}

	return returns.WhID, nil
}

func (repository *RepositoryReturnImpl) FindOrderDetailIdByInvoiceNo(invoiceNo int64) ([]model.ReturnDetail, error) {
	var returns []model.ReturnDetail

	// Calculate the date one month ago
	OneMonthAgo := time.Now().AddDate(0, -1, 0).Format("2006-01-02")

	// Define the query to select distinct return_no, invoice_no, invoice_date, and product_id
	err := repository.Select("DISTINCT r.return_no, r.invoice_no, r.invoice_date, rd.product_id").
		Table("(SELECT * FROM sls.return WHERE invoice_date >= ?) AS r", OneMonthAgo).
		Joins("JOIN sls.return_det AS rd ON r.return_no = rd.return_no").
		Where("rd.product_id = ? ", invoiceNo).
		Order("r.invoice_date DESC").
		Find(&returns).Error

	if err != nil {
		return nil, err
	}

	return returns, nil
}
