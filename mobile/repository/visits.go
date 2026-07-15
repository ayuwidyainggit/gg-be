package repository

import (
	"context"
	"errors"
	"fmt"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/constant"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryVisitsImpl struct {
		*gorm.DB
	}
)

type VisitsRepository interface {
	Store(data *model.Visit) error
	FindAllByCustId(dataFilter entity.VisitQueryFilter) ([]model.Outlet, error)
	GetStartByEmployeeBetweenTime(empCode string, custID string, start, end time.Time) (visit *model.Visit, err error)
	GetLastStatusByOutletCodeBetweenTime(outletCode string, start, end time.Time) (visit *model.Visit, err error)
	FindAllSkipReasonByCustId(dataFilter entity.SkipReasonsQueryFilter) ([]model.SkipReason, error)
	GetLastVisitGroupOutlet(custID string, start, end time.Time) (visits []model.Visit, err error)
	StoreOutletVisitList(data *model.OutletVisitList) error
	FindOutletByCode(custID string, outletCode string) (model.Outlet, error)
	UpdateOutletLocation(outletID int64, latitude, longitude string) error
	GetVisitListByCustID(request entity.VisitsListRequest) ([]model.OutletVisitList, error)
}

func NewVisitsRepository(db *gorm.DB) *RepositoryVisitsImpl {
	return &RepositoryVisitsImpl{db}
}

func (repo *RepositoryVisitsImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
func (repository *RepositoryVisitsImpl) Store(data *model.Visit) error {
	err := repository.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryVisitsImpl) FindAllByCustId(dataFilter entity.VisitQueryFilter) ([]model.Outlet, error) {
	var outlets []model.Outlet

	query := repository.
		Select("m_outlet.outlet_id, m_outlet.outlet_code, m_outlet.outlet_name, m_outlet.barcode, m_outlet.email, m_outlet.wa_no, m_outlet.latitude, m_outlet.longitude, m_outlet.address1, m_outlet.address2, m_outlet.city, m_outlet.zip_code, m_outlet.phone_no, b.no_order_id").
		Joins("LEFT JOIN sls.no_order as b on b.outlet_id = m_outlet.outlet_id")
	// Joins("LEFT JOIN mst.m_industry ind ON ind.industry_id = m_outlet.industry_id AND ind.cust_id = ?", dataFilter.ParentCustId)
	// Joins("LEFT JOIN mst.m_supplier s ON s.sup_id = m_product.sup_id AND s.cust_id = ?", parentCustId).
	// Joins("LEFT JOIN mst.m_sub_brand1 sbr ON sbr.sbrand1_id = m_product.sbrand1_id AND sbr.cust_id = ?", parentCustId).
	// Joins("LEFT JOIN mst.m_brand br ON br.brand_id = sbr.brand_id AND br.cust_id = ?", parentCustId)

	query.Where("m_outlet.cust_id = ? AND m_outlet.is_active = true", dataFilter.CustId)

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
		query.Order("m_outlet.outlet_code ASC")
	}

	err := query.Limit(10).Find(&outlets).Error
	if err != nil {
		return outlets, err
	}

	return outlets, nil

}

func (repository *RepositoryVisitsImpl) GetStartByEmployeeBetweenTime(empCode string, custID string, start, end time.Time) (visit *model.Visit, err error) {
	query := repository.Select("*")
	if empCode != "" {
		query.Where("emp_code = ?", empCode)
	}
	if custID != "" {
		query.Where("cust_id = ?", custID)
	}

	query.
		Where("type = ? AND created_at between ? AND ?", entity.TYPE_START_ID, start, end).
		Order("created_at DESC")

	err = query.Take(&visit).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New(constant.STATUS_DB_NOT_FOUND)
	}
	return visit, err
}

func (repository *RepositoryVisitsImpl) GetLastStatusByOutletCodeBetweenTime(outletCode string, start, end time.Time) (visit *model.Visit, err error) {
	query := repository.Select("*")
	if outletCode != "" {
		query = query.Where("outlet_code = ? AND created_at between ? AND ?", outletCode, start, end)
	} else {
		query = query.Where("created_at between ? AND ?", start, end)
	}

	query = query.Order("created_at DESC")

	err = query.Take(&visit).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New(constant.STATUS_DB_NOT_FOUND)
	}
	return visit, err
}

func (repository *RepositoryVisitsImpl) FindAllSkipReasonByCustId(dataFilter entity.SkipReasonsQueryFilter) ([]model.SkipReason, error) {
	var skipReasons []model.SkipReason

	query := repository.Where("m_skip_reason.cust_id = ? AND m_skip_reason.is_active = true", dataFilter.CustId)

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
		query.Order("m_skip_reason.skip_reason_name ASC")
	}

	err := query.Find(&skipReasons).Error
	if err != nil {
		return skipReasons, err
	}

	return skipReasons, nil

}
func (repository *RepositoryVisitsImpl) GetLastVisitGroupOutlet(custID string, start, end time.Time) (visits []model.Visit, err error) {
	err = repository.Select("mobile.visits.*").
		Joins("INNER JOIN (select outlet_code, max(created_at) as max_date from mobile.visits WHERE cust_id=? AND  created_at between ? AND ? AND outlet_code is not null GROUP BY outlet_code) vi2 on mobile.visits.outlet_code = vi2.outlet_code AND mobile.visits.created_at = vi2.max_date ", custID, start, end).
		Find(&visits).Error
	if err != nil {
		return visits, err
	}

	return visits, nil
}

func (repository *RepositoryVisitsImpl) StoreOutletVisitList(data *model.OutletVisitList) error {
	err := repository.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryVisitsImpl) FindOutletByCode(custID string, outletCode string) (model.Outlet, error) {
	var outlet model.Outlet
	err := repository.
		Where("cust_id = ? AND outlet_code = ?", custID, outletCode).
		Take(&outlet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return outlet, errors.New(constant.STATUS_DB_NOT_FOUND)
	}
	return outlet, err
}

func (repository *RepositoryVisitsImpl) UpdateOutletLocation(outletID int64, latitude, longitude string) error {
	updateData := map[string]interface{}{
		"latitude":   latitude,
		"longitude":  longitude,
		"updated_at": time.Now().UTC(),
	}
	err := repository.Table("mst.m_outlet").
		Where("outlet_id = ?", outletID).
		Updates(updateData).Error
	return err
}

func (repository *RepositoryVisitsImpl) GetVisitListByCustID(request entity.VisitsListRequest) ([]model.OutletVisitList, error) {
	var visits []model.OutletVisitList

	query := repository.
		Table("pjp.outlet_visit_list as ovl").
		Select(`
			ovl.id,
			ovl.outlet_id,
			ovl.outlet_code,
			m_outlet.outlet_name,
			ovl.date,
			ovl.arrive_at,
			ovl.latitude,
			ovl.longitude,
			ovl.file_url as photo_path,
			ovl.folder,
			ovl.is_update_location,
			ovl.created_at,
			ovl.updated_at
		`).
		Joins("INNER JOIN mst.m_outlet ON m_outlet.outlet_id = ovl.outlet_id AND m_outlet.cust_id = ?", request.CustID)
		// Joins("INNER JOIN mobile.visits ON mobile.visits.outlet_code = ovl.outlet_code and mobile.visits.emp_code= ?", request.EmpCode)
		// visits.latitude,
		// visits.longitude,
		// visits.file_url as photo_path,
		// comment cause support skip visit

	// Add outlet_id filter if provided
	if request.OutletID != nil {
		query = query.Where("ovl.outlet_id = ?", &request.OutletID)
	}

	if request.PJPID > 0 {
		query = query.Where("ovl.pjp_id = ?", request.PJPID)
	}

	err := query.
		Order("ovl.created_at DESC").
		First(&visits).Error

	if err != nil {
		return visits, err
	}

	return visits, nil
}
