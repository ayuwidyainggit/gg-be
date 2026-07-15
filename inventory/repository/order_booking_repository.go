package repository

import (
	"context"
	"errors"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryOrderBookingImpl struct {
		*gorm.DB
	}
)

type OrderBookingRepository interface {
	FindAllByCustId(dataFilter entity.OrderQueryFilter) ([]model.OrderBookingList, int64, int, error)
	FindByNo(OrderBookingId int, custId string, parentCustId string) (realOrder model.OrderBookingList, err error)
	FindDetail(OrderBookingId int, custId string) (details []model.OrderBookingDetailRead, err error)
	Store(c context.Context, data *model.OrderBooking) error
	StoreDetail(c context.Context, data *model.OrderBookingDetail) error
	DeleteDetailNotInIDs(c context.Context, OrderBookingId int, IDs []int64) error
	Delete(c context.Context, custId string, OrderBookingId int, deletedBy int64) error
	UpdateApproval(c context.Context, OrderBookingId int, data model.OrderBookingDetailStatus) error
	UpdateCompleted(c context.Context, PoNo string, data model.OrderBookingDetailStatus) error
	UpdateOrder(c context.Context, OrderBookingId int, data model.OrderBookingDetailStatusApproval) error
	UpdateOrderDetail(c context.Context, OrderBookingDetailId int, data model.OrderBookingDetailApproval) error
	FindDetailByNotInDetailIDs(detailIDs []int64, OrderBookingId int, custId string) (details []model.OrderBookingDetailRead, err error)
	FindDetailFinal(grBranchNo string, custId string, orderBookingId int) (details []model.OrderBookingDetailFinalRead, err error)
	CountAllByCustId(custId string) (int, error)
}

func NewOrderBookingRepo(db *gorm.DB) *RepositoryOrderBookingImpl {
	return &RepositoryOrderBookingImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryOrderBookingImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryOrderBookingImpl) FindAllByCustId(dataFilter entity.OrderQueryFilter) ([]model.OrderBookingList, int64, int, error) {
	var ro []model.OrderBookingList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("order_booking_id")
	query := repository.Select(
		`inv.order_booking.*, 
			md.distributor_name AS distributor_name,
			md.address AS distributor_address,
			up.user_fullname AS updated_by_name, 
			uc.user_fullname AS created_by_name,
			sup.sup_name,sup.sup_code,sup.credit_limit`).
		Joins("left join sys.m_user up on up.user_id = inv.order_booking.updated_by").
		Joins("left join sys.m_user uc on uc.user_id = inv.order_booking.created_by").
		Joins("left join smc.m_customer mc on mc.cust_id = inv.order_booking.cust_id").
		Joins("left join mst.m_distributor md on md.distributor_id = mc.distributor_id").
		Joins("left join mst.m_supplier sup on sup.sup_id = inv.order_booking.sup_id AND sup.cust_id = ?", dataFilter.ParentCustId)

	if dataFilter.CustId != dataFilter.ParentCustId {
		queryCount.Where("inv.order_booking.cust_id=?", dataFilter.CustId)
		query.Where("inv.order_booking.cust_id=?", dataFilter.CustId)
	} else {
		queryCount.Where("inv.order_booking.cust_id LIKE ?", dataFilter.ParentCustId+"%")
		query.Where("inv.order_booking.cust_id LIKE ?", dataFilter.ParentCustId+"%")
	}

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("inv.order_booking.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("inv.order_booking.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		search := "%" + dataFilter.Query + "%"
		queryCount.Where(
			"(inv.order_booking.po_no ILIKE ? OR inv.order_booking.so_po ILIKE ?)",
			search, search,
		)
		query.Where(
			"(inv.order_booking.po_no ILIKE ? OR inv.order_booking.so_po ILIKE ?)",
			search, search,
		)
	}

	if len(dataFilter.Status) > 0 {
		queryCount.Where("inv.order_booking.status_order_booking in ?", dataFilter.Status)
		query.Where("inv.order_booking.status_order_booking in ?", dataFilter.Status)
	}

	// sortBy := ``
	// if dataFilter.Sort != "" {
	// 	mSortBy := strings.Split(dataFilter.Sort, ",")
	// 	for _, row := range mSortBy {
	// 		colSort := strings.Split(row, ":")
	// 		if len(colSort) > 1 {
	// 			sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
	// 		}
	// 	}
	// 	sortBy = strings.TrimSuffix(sortBy, ", ")
	// 	query.Order(sortBy)
	// } else {
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
		query.Order("order_booking_id DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ro).Error
	if err != nil {
		return ro, total, 0, err
	}
	err = queryCount.Model(&ro).Count(&total).Error
	if err != nil {
		return ro, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ro, total, lastPage, nil
}

func (repository *RepositoryOrderBookingImpl) Store(c context.Context, data *model.OrderBooking) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderBookingImpl) StoreDetail(c context.Context, data *model.OrderBookingDetail) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderBookingImpl) FindByNo(OrderBookingId int, custId string, parentCustId string) (realOrder model.OrderBookingList, err error) {
	query := repository.
		Select(`inv.order_booking.*, 
		up.user_fullname AS updated_by_name, 
		uc.user_fullname AS created_by_name,
		mc.distributor_id,
		md.distributor_name AS distributor_name,
		md.distributor_code AS distributor_code,
		md.address AS distributor_address,
		grb.gr_branch_no,grb.sub_total as sub_total_final,grb.vat_value as vat_value_final,grb.total as total_final,grb.delivery_fee as delivery_fee_final,
		sup.sup_name,sup.sup_code,sup.credit_limit`).
		Joins("left join sys.m_user up on up.user_id = inv.order_booking.updated_by").
		Joins("left join sys.m_user uc on uc.user_id = inv.order_booking.created_by").
		Joins("left join smc.m_customer mc on mc.cust_id = inv.order_booking.cust_id").
		Joins("left join mst.m_distributor md on md.distributor_id = mc.distributor_id").
		Joins("left join inv.gr_branch grb on inv.order_booking.po_no = grb.po_no AND grb.cust_id = ?", custId).
		Joins("left join mst.m_supplier sup on sup.sup_id = inv.order_booking.sup_id AND sup.cust_id = ?", parentCustId).
		Where("inv.order_booking.order_booking_id = ?", OrderBookingId)

	if custId != parentCustId {
		query = query.Where("inv.order_booking.cust_id = ?", custId)
	} else {
		query = query.Where("inv.order_booking.cust_id LIKE ?", parentCustId+"%")
	}

	err = query.Take(&realOrder).Error
	return realOrder, err
}

func (repository *RepositoryOrderBookingImpl) FindDetail(OrderBookingId int, custId string) (details []model.OrderBookingDetailRead, err error) {
	err = repository.Select("inv.order_booking_detail.*, p.pro_code, p.pro_name, p.conv_unit2 as mconv_unit2,p.conv_unit3 as mconv_unit3,p.sup_id,sp.sup_name").
		Joins("left join mst.m_product p on p.pro_id = inv.order_booking_detail.pro_id").
		Joins("left join mst.m_supplier sp on p.sup_id = sp.sup_id").
		Where("order_booking_id = ?", OrderBookingId).
		Find(&details).Error
	return details, err
}

func (repository *RepositoryOrderBookingImpl) FindDetailFinal(grBranchNo string, custId string, orderBookingId int) (details []model.OrderBookingDetailFinalRead, err error) {

	subQuery := repository.
		Select("pro_id").
		Table("inv.order_booking_detail").
		Where("order_booking_id = ?", orderBookingId) // contoh filter di tabel lain

	err = repository.Select("inv.gr_branch_det.*, p.pro_code, p.pro_name, p.conv_unit2 as mconv_unit2,p.conv_unit3 as mconv_unit3").
		Joins("left join mst.m_product p on p.pro_id = inv.gr_branch_det.pro_id").
		Where("gr_branch_no = ?", grBranchNo).
		Where("item_type = ?", 1).
		Where("inv.gr_branch_det.cust_id = ?", custId).
		Where("inv.gr_branch_det.pro_id IN (?)", subQuery).
		Find(&details).Error
	return details, err
}

func (repository *RepositoryOrderBookingImpl) DeleteDetailNotInIDs(c context.Context, OrderBookingId int, IDs []int64) error {
	var Details model.OrderBookingDetail

	if len(IDs) == 0 {
		IDs = append(IDs, 0)
	}

	err := repository.model(c).Where("order_booking_id=? AND order_detail_id not in (?) AND item_type = 1", OrderBookingId, IDs).Delete(&Details).Error
	return err
}

func (repository *RepositoryOrderBookingImpl) Delete(c context.Context, custId string, OrderBookingId int, deletedBy int64) error {
	var data model.OrderBooking
	result := repository.model(c).Model(&data).Where("order_booking_id=? AND cust_id = ? AND is_del= ? ", OrderBookingId, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryOrderBookingImpl) UpdateApproval(c context.Context, OrderBookingId int, data model.OrderBookingDetailStatus) error {
	result := repository.model(c).Model(&data).Where("order_booking_id=?", OrderBookingId).Updates(data)

	if result.Error != nil {
		return result.Error
	}

	// if result.OrderwsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryOrderBookingImpl) UpdateCompleted(c context.Context, PoNo string, data model.OrderBookingDetailStatus) error {
	result := repository.model(c).Model(&data).Where("po_no=? AND cust_id=? AND status_order_booking=2", PoNo, data.CustID).Updates(data)

	if result.Error != nil {
		return result.Error
	}

	// if result.OrderwsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryOrderBookingImpl) UpdateOrder(c context.Context, OrderBookingId int, data model.OrderBookingDetailStatusApproval) error {
	result := repository.model(c).Model(&data).Where("order_booking_id=?", OrderBookingId).Updates(data)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		fmt.Println("no rows affected")
	}
	return nil
}

func (repository *RepositoryOrderBookingImpl) UpdateOrderDetail(c context.Context, OrderBookingDetailId int, data model.OrderBookingDetailApproval) error {
	result := repository.model(c).Model(&data).Where("order_booking_detail_id=?", OrderBookingDetailId).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		fmt.Println("no rows affected")
	}
	return nil
}

func (repository *RepositoryOrderBookingImpl) FindDetailByNotInDetailIDs(detailIDs []int64, OrderBookingId int, custId string) (details []model.OrderBookingDetailRead, err error) {

	if len(detailIDs) == 0 {
		detailIDs = append(detailIDs, 0)
	}

	err = repository.Select("inv.order_booking_detail.*").
		Where("order_booking_detail_id not in ? AND order_booking_id = ? AND inv.order_booking_detail.cust_id = ?", detailIDs, OrderBookingId, custId).
		Find(&details).Error

	return details, err
}

func (repository *RepositoryOrderBookingImpl) CountAllByCustId(custId string) (int, error) {
	var ob []model.OrderBookingList
	var total int64

	queryCount := repository.Select("order_booking_id")

	queryCount.Where("inv.order_booking.cust_id = ?", custId)
	queryCount.Where("inv.order_booking.created_at::date = CURRENT_DATE") // Menambahkan kondisi tanggal sekarang

	err := queryCount.Model(&ob).Count(&total).Error
	if err != nil {
		return 0, err
	}

	return int(total), nil
}
