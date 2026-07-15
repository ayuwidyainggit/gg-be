package repository

import (
	"fmt"
	"math"
	"mobile/entity"
	"mobile/model"
	"strings"

	"gorm.io/gorm"
)

type (
	RepositoryMProductDistImpl struct {
		*gorm.DB
	}
)
type MProductDistRepository interface {
	FindAllByCustId(dataFilter entity.ProductsQueryFilter, custId, parentCustId string) ([]model.MProductDist, int64, int, error)
}

func NewMProductDistRepository(db *gorm.DB) *RepositoryMProductDistImpl {
	return &RepositoryMProductDistImpl{db}
}

// func (repo *RepositoryMProductDistImpl) model(ctx context.Context) *gorm.DB {
// 	tx := extractTx(ctx)
// 	if tx != nil {
// 		return tx.WithContext(ctx)
// 	}
// 	return repo.WithContext(ctx)
// }

func (repository *RepositoryMProductDistImpl) FindAllByCustId(dataFilter entity.ProductsQueryFilter, custId, parentCustId string) ([]model.MProductDist, int64, int, error) {
	var products []model.MProductDist
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("pro_id")
	queryCount.Where("m_product_dist.cust_id=?", custId)

	query := repository.
		Select(`
			m_product_dist.cust_id, m_product_dist.pro_id, m_product_dist.is_active, 
			p.pro_code, p.pro_name, p.sup_id, 
			sbr.brand_id, p.sbrand1_id, 
			p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5,
			p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5, 
			s.sup_code, s.sup_name,
			br.brand_code, br.brand_name,
			sbr.sbrand1_code, sbr.sbrand1_name, p.sell_price1, p.sell_price2, p.sell_price3, p.sell_price4, p.sell_price5`).
		Joins("LEFT JOIN mst.m_product p ON p.pro_id = m_product_dist.pro_id AND p.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_unit un1 ON un1.unit_id = p.unit_id1 AND p.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_supplier s ON s.sup_id = p.sup_id AND s.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_sub_brand1 sbr ON sbr.sbrand1_id = p.sbrand1_id AND sbr.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_brand br ON br.brand_id = sbr.brand_id AND br.cust_id = ?", parentCustId)

	query.Where("m_product_dist.cust_id = ? AND m_product_dist.is_active = true", custId)

	if dataFilter.SupId != nil {
		if *dataFilter.SupId > 0 {
			query.Where("p.sup_id = ?", dataFilter.SupId)
			queryCount.Where("p.sup_id = ?", dataFilter.SupId)
		}
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
		query.Order("p.pro_code ASC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return products, total, 0, err
	}

	err = queryCount.Model(&products).Count(&total).Error
	if err != nil {
		return products, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return products, total, lastPage, nil

}
