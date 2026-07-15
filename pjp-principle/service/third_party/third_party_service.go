package thirdparty

import (
	"context"
	"fmt"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/model"
	distributordms "scyllax-pjp/repository/distributor_dms"
	outletdms "scyllax-pjp/repository/outlet_dms"
	"scyllax-pjp/repository/pjp"
	"strings"

	"gorm.io/gorm"
)

type ThirdPartyService interface {
	GetAssignedSalesman(ctx context.Context, headers map[string]string, custId string) response.ListSalesmanAPIResponse
	GetUnassignedSalesman(ctx context.Context, headers map[string]string, dataFilter request.SalesmanListQueryFilter, custId string) ([]model.NewSalesman, response.Meta)
	GetOutlet(ctx context.Context, dataFilter model.OutletQueryFilter, custId string) ([]model.OutletDms, response.Meta)
	GetDistributor(ctx context.Context, dataFilter model.DistributorQueryFilter, custId string) ([]model.DistributorDms, response.Meta)
	GetSalesmanByID(ctx context.Context, empId int, headers map[string]string, custId string) model.NewSalesman
}
type thirdPartyService struct {
	pjpRepo         pjp.PjpRepository
	outletRepo      outletdms.OutletDmsRepository
	distributorRepo distributordms.DistributorDmsRepository
	db              *gorm.DB
}

func NewThirdPartyService(pjpRepo pjp.PjpRepository, outletRepo outletdms.OutletDmsRepository, distributorRepo distributordms.DistributorDmsRepository, db *gorm.DB) ThirdPartyService {
	return &thirdPartyService{
		pjpRepo:         pjpRepo,
		outletRepo:      outletRepo,
		distributorRepo: distributorRepo,
		db:              db,
	}
}

func masterSalesmanEndpointURL(baseURL string, filter request.SalesmanListQueryFilter) string {
	return fmt.Sprintf("%s/v1/salesman?limit=%s&page=%s&sort=%s&is_active=%d&sales_team_id=%s",
		baseURL, filter.Limit, filter.Page, filter.Sort, 1, filter.SalesTeamID)
}

func setOperationTypeName(salesmen []model.NewSalesman) {
	for i := range salesmen {
		s := &salesmen[i]
		var ops []string
		if s.IsActiveCanvas {
			ops = append(ops, "Canvas")
		}
		if s.IsTakingOrder {
			ops = append(ops, "Taking Order")
		}
		s.OperationTypeName = strings.Join(ops, ", ")
	}
}

func filterUnassignedSalesmen(all []model.NewSalesman, assignedIDs []int) []model.NewSalesman {
	assignedSet := make(map[int]struct{}, len(assignedIDs))
	for _, id := range assignedIDs {
		assignedSet[id] = struct{}{}
	}

	var filtered []model.NewSalesman
	for _, s := range all {
		if _, exists := assignedSet[s.EmployeeID]; !exists {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func generatePagination(pageCurrent, pageLimit int, totalFiltered int) response.Meta {
	totalPages := 0
	if pageLimit > 0 {
		totalPages = (totalFiltered + pageLimit - 1) / pageLimit
	}
	return response.Meta{
		TotalData: totalFiltered,
		Page:      pageCurrent,
		Limit:     pageLimit,
		TotalPage: totalPages,
	}
}
