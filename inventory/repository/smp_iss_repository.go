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
	RepositorySmpIssImpl struct {
		*gorm.DB
	}
)
type SmpIssRepository interface {
	Store(c context.Context, data *model.SampleIssue) error
	StoreDetail(c context.Context, data *model.SampleIssueDet) error
	FindByNo(smpIssueNo string, custId, parentCustId string) (smpIssue model.SampleIssueList, err error)
	FindSmpIssuedetail(smpIssueNo string, custId string) (Details []model.SampleIssueDetRead, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string, parentCustId string) ([]model.SampleIssueList, int64, int, error)
	Delete(c context.Context, custId string, smpIsNo string, deletedBy int64) error
	Update(c context.Context, smpIsNo string, data model.SampleIssue) error
	DeleteDetailNotInIDs(c context.Context, smpIsNo string, IDs []int) error
	UpdateGrDetail(c context.Context, Details *model.SampleIssueDet) error
}

func NewSmpIssRepo(db *gorm.DB) *RepositorySmpIssImpl {
	return &RepositorySmpIssImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositorySmpIssImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
func (repository *RepositorySmpIssImpl) Store(c context.Context, data *model.SampleIssue) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositorySmpIssImpl) StoreDetail(c context.Context, data *model.SampleIssueDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositorySmpIssImpl) FindByNo(smpIssueNo string, custId, parentCustId string) (smpIssue model.SampleIssueList, err error) {
	err = repository.
		Select("sample_issue.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, cndn.cndn_code, cndn.cndn_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = sample_issue.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = sample_issue.closed_by").
		Joins("left join mst.m_warehouse wh on wh.wh_id = sample_issue.wh_id AND wh.cust_id = ?", custId).
		Joins("left join mst.m_cndn cndn on cndn.cndn_id = sample_issue.cndn_id AND cndn.cust_id = ?", parentCustId).
		Where("sample_issue.smp_iss_no = ? AND sample_issue.cust_id=?", smpIssueNo, custId).
		Take(&smpIssue).Error
	return smpIssue, err
}

func (repository *RepositorySmpIssImpl) FindSmpIssuedetail(smpIssueNo string, custId string) (Details []model.SampleIssueDetRead, err error) {
	err = repository.Select("inv.sample_issue_det.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = inv.sample_issue_det.pro_id").
		Where("smp_iss_no = ? AND inv.sample_issue_det.cust_id=?", smpIssueNo, custId).Order("seq_no ASC").
		Find(&Details).Error
	return Details, err
}

func (repository *RepositorySmpIssImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string, parentCustId string) ([]model.SampleIssueList, int64, int, error) {
	var smpIss []model.SampleIssueList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("smp_iss_no")
	query := repository.Select("sample_issue.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, cndn.cndn_code, cndn.cndn_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = sample_issue.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = sample_issue.closed_by").
		Joins("left join mst.m_warehouse wh on wh.wh_id = sample_issue.wh_id AND wh.cust_id = ?", custId).
		Joins("left join mst.m_cndn cndn on cndn.cndn_id = sample_issue.cndn_id AND cndn.cust_id = ?", parentCustId)

	queryCount.Where("sample_issue.cust_id=?", custId)
	query.Where("sample_issue.cust_id=?", custId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sample_issue.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sample_issue.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("sample_issue.smp_iss_no=?", dataFilter.Query)
		queryCount.Where("sample_issue.smp_iss_no=?", dataFilter.Query)
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
		query.Order("smp_iss_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&smpIss).Error
	if err != nil {
		return smpIss, total, 0, err
	}
	err = queryCount.Model(&smpIss).Count(&total).Error
	if err != nil {
		return smpIss, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return smpIss, total, lastPage, nil
}

func (repository *RepositorySmpIssImpl) Delete(c context.Context, custId string, smpIsNo string, deletedBy int64) error {
	var data model.SampleIssue
	result := repository.model(c).Model(&data).Where("smp_iss_no=? AND cust_id = ? AND is_del= ? ", smpIsNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositorySmpIssImpl) Update(c context.Context, smpIsNo string, data model.SampleIssue) error {
	result := repository.model(c).Model(&data).Where("smp_iss_no=?", smpIsNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositorySmpIssImpl) DeleteDetailNotInIDs(c context.Context, smpIsNo string, IDs []int) error {
	var Details model.SampleIssueDet
	err := repository.model(c).Where("smp_iss_no=? AND smp_iss_det_id not in (?) ", smpIsNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositorySmpIssImpl) UpdateGrDetail(c context.Context, Details *model.SampleIssueDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
