package thirdparty

import (
	"context"
	"fmt"
	"net/url"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/model"
	"scyllax-pjp/repository/pjp"
	routeoutlet "scyllax-pjp/repository/route_outlet"
	routeoutlethistory "scyllax-pjp/repository/route_outlet_history"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type ThirdPartyService interface {
	GetAssignedSalesman(ctx context.Context, custId string) response.ListSalesmanAPIResponse
	GetUnassignedSalesman(ctx context.Context, dataFilter request.SalesmanListQueryFilter, custId string) ([]model.NewSalesman, response.Meta)
	GetOutlet(ctx context.Context, dataFilter model.DmsQueryFilter, custId string) ([]model.Outlet, response.Meta)
	GetSalesmanByID(ctx context.Context, empId int, headers map[string]string, custId string) model.NewSalesman
	GetOutletBySalesCodes(ctx context.Context, dataFilter model.OutletBySalesman, headers map[string]string, custId string) ([]model.OutletNew, response.Meta, error)
	GetOutletPicklistBySalesCodes(ctx context.Context, dataFilter model.OutletBySalesman, headers map[string]string, custId string) ([]model.OutletNew, response.Meta, error)
}
type thirdPartyService struct {
	pjpRepo                pjp.PjpRepository
	routeOutletRepo        routeoutlet.RouteOutletRepository
	routeOutletHistoryRepo routeoutlethistory.RouteOutletHistoryRepository
	db                     *gorm.DB
}

func NewThirdPartyService(pjpRepo pjp.PjpRepository, routeOutlet routeoutlet.RouteOutletRepository, routeOutletHistory routeoutlethistory.RouteOutletHistoryRepository, db *gorm.DB) ThirdPartyService {
	return &thirdPartyService{
		pjpRepo:                pjpRepo,
		routeOutletRepo:        routeOutlet,
		routeOutletHistoryRepo: routeOutletHistory,
		db:                     db,
	}
}

func masterSalesmanEndpointURL(baseURL string, filter request.SalesmanListQueryFilter) string {
	return fmt.Sprintf("%s/v1/salesman?limit=%s&page=%s&sort=%s&is_active=%d&sales_team_id=%s",
		baseURL, filter.Limit, filter.Page, filter.Sort, 1, filter.SalesTeamID)
}

func masterOutletEndpointURL(baseURL string, filter model.DmsQueryFilter) string {
	query := url.Values{}
	query.Set("verification_status", "1")

	if filter.Limit != "" {
		query.Set("limit", filter.Limit)
	}
	if filter.Page != "" {
		query.Set("page", filter.Page)
	}
	if filter.Sort != "" {
		query.Set("sort", filter.Sort)
	}
	if filter.Query != "" {
		query.Set("q", filter.Query)
	}
	if filter.OutletCode != "" {
		query.Set("outlet_code", filter.OutletCode)
	}
	if filter.OutletID > 0 {
		query.Set("outlet_id", strconv.Itoa(filter.OutletID))
	}
	if filter.Mode != "" {
		query.Set("mode", filter.Mode)
	}
	if filter.IsActive != "" {
		query.Set("is_active", filter.IsActive)
	}
	if filter.SalesTeamID != "" {
		query.Set("sales_team_id", filter.SalesTeamID)
	}

	return fmt.Sprintf("%s/v1/outlets?%s", baseURL, query.Encode())
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

func masterOutletByIDsEndpointURL(baseURL, outletIDs string, limit int, includeInactive int) string {
	query := fmt.Sprintf("%s/v1/outlets?outlet_id=%s&limit=%s", baseURL, outletIDs, strconv.Itoa(limit))
	if includeInactive == 1 {
		query += "&include_inactive=1"
	}
	return query
}
