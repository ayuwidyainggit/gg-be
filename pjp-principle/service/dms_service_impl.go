package service

import (
	"scyllax-pjp/repository"
)

type DmsServiceImpl struct {
	outletRepo      repository.OutletRepository
	routeOutletRepo repository.RouteOutletRepository
}

func NewDmsServiceImpl(outletRepo repository.OutletRepository, routeOutletRepo repository.RouteOutletRepository) DmsService {
	return &DmsServiceImpl{
		outletRepo:      outletRepo,
		routeOutletRepo: routeOutletRepo,
	}
}

// func (service *DmsServiceImpl) GetSalesTeam(ctx context.Context, custId string) ([]model.SalesTeam, error) {
// 	config, err := config.LoadConfig(".")

// 	if err != nil {
// 		log.Fatal("could not load config", err)
// 	}

// 	endpointURL := config.KongUrl + "/v1/sales-teams?q&mode=lookup"
// 	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("cust_id", custId)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var response struct {
// 		Data []model.SalesTeam `json:"data"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 		return nil, err
// 	}

// 	return response.Data, nil
// }

// func (service *DmsServiceImpl) GetSalesman(ctx context.Context, dataFilter model.DmsQueryFilter, custId string) ([]model.NewSalesman, response.Meta, error) {
// 	config, err := config.LoadConfig(".")
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}

// 	// kong public add /master for dev remove for staging
// 	endpointURL := fmt.Sprintf("%s/v1/salesman?limit=%s&page=%s&sort=%s&is_active=%d&sales_team_id=%s",
// 		config.KongUrl, dataFilter.Limit, dataFilter.Page, dataFilter.Sort, 1, dataFilter.SalesTeamID)
// 	log.Printf("Endpoint URL: %s", endpointURL)

// 	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}
// 	req.Header.Set("cust_id", custId)

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, response.Meta{}, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}

// 	var result struct {
// 		Data   []model.NewSalesman `json:"data"`
// 		Paging struct {
// 			TotalRecord int `json:"total_record"`
// 			PageCurrent int `json:"page_current"`
// 			PageLimit   int `json:"page_limit"`
// 			PageTotal   int `json:"page_total"`
// 		} `json:"paging"`
// 	}
// 	if err := json.Unmarshal(body, &result); err != nil {
// 		log.Printf("JSON Unmarshal Error: %v", err)
// 		return nil, response.Meta{}, err
// 	}

// 	for i := range result.Data {
// 		s := &result.Data[i]

// 		var ops []string
// 		if s.IsActiveCanvas {
// 			ops = append(ops, "Canvas")
// 		}
// 		if s.IsTakingOrder {
// 			ops = append(ops, "Taking Order")
// 		}
// 		s.OperationTypeName = strings.Join(ops, ", ")
// 	}

// 	// Call assined salesman to pjp
// 	assignedSalesmanIDs := service.pjpRepo.FindAssignedSalesman(ctx)

// 	// Assign salesman id exist to Set
// 	assignedSalesmanSet := make(map[int]struct{})
// 	for _, id := range assignedSalesmanIDs {
// 		assignedSalesmanSet[id] = struct{}{}
// 	}

// 	// filter salesman where not exist and append as a result to show
// 	filteredSalesmen := make([]model.NewSalesman, 0)
// 	for _, salesman := range result.Data {
// 		if _, exists := assignedSalesmanSet[salesman.EmployeeID]; !exists {
// 			filteredSalesmen = append(filteredSalesmen, salesman)
// 		}
// 	}

// 	totalPages := 0
// 	if result.Paging.PageLimit > 0 {
// 		totalPages = (len(filteredSalesmen) + result.Paging.PageLimit - 1) / result.Paging.PageLimit
// 	}

// 	pagination := response.Meta{
// 		TotalData: len(filteredSalesmen),
// 		Page:      result.Paging.PageCurrent,
// 		Limit:     result.Paging.PageLimit,
// 		TotalPage: totalPages,
// 	}

// 	return filteredSalesmen, pagination, nil
// }

// func (service *DmsServiceImpl) GetSalesmanByID(ctx context.Context, empId int, headers map[string]string, custId string) (model.NewSalesman, error) {
// 	config, err := config.LoadConfig(".")
// 	if err != nil {
// 		return model.NewSalesman{}, err
// 	}

// 	// kong public add /master for dev remove for staging
// 	endpointURL := fmt.Sprintf("%s/v1/salesman/%d",
// 		config.KongUrl, empId)
// 	log.Printf("Endpoint URL: %s", endpointURL)

// 	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
// 	if err != nil {
// 		return model.NewSalesman{}, err
// 	}
// 	req.Header.Set("cust_id", custId)
// 	for key, value := range headers {
// 		req.Header.Set(key, value)
// 	}

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return model.NewSalesman{}, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		log.Printf("HTTP request failed with status code %d", resp.StatusCode)
// 		return model.NewSalesman{}, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return model.NewSalesman{}, err
// 	}

// 	var result struct {
// 		Message   string            `json:"message"`
// 		Data      model.NewSalesman `json:"data"`
// 		RequestID string            `json:"request_id"`
// 	}
// 	if err := json.Unmarshal(body, &result); err != nil {
// 		log.Printf("JSON Unmarshal Error: %v", err)
// 		return model.NewSalesman{}, err
// 	}

// 	return result.Data, nil
// }

// func (service *DmsServiceImpl) GetWarehouse(ctx context.Context, custId string) ([]model.Warehouse, error) {
// 	config, err := config.LoadConfig(".")

// 	if err != nil {
// 		log.Fatal("could not load config", err)
// 	}

// 	endpointURL := fmt.Sprintf("%s/v1/warehouses?mode=lookup&sort=wh_name:asc", config.KongUrl)
// 	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("cust_id", custId)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var response struct {
// 		Data []model.Warehouse `json:"data"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 		return nil, err
// 	}

// 	return response.Data, nil
// }

// func (service *DmsServiceImpl) GetOutlet(ctx context.Context, dataFilter model.DmsQueryFilter, custId string) ([]model.Outlet, response.Meta, error) {
// 	config, err := config.LoadConfig(".")

// 	if err != nil {
// 		log.Fatal("could not load config", err)
// 	}

// 	endpointURL := fmt.Sprintf("%s/v1/outlets?limit=%s&page=%s&sort=%s&is_active=%d",
// 		config.KongUrl, dataFilter.Limit, dataFilter.Page, dataFilter.Sort, 1)

// 	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}
// 	req.Header.Add("cust_id", custId)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}
// 	defer resp.Body.Close()

// 	var result struct {
// 		Data   []model.Outlet `json:"data"`
// 		Paging struct {
// 			TotalRecord int `json:"total_record"`
// 			PageCurrent int `json:"page_current"`
// 			PageLimit   int `json:"page_limit"`
// 			PageTotal   int `json:"page_total"`
// 		} `json:"paging"`
// 	}
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}

// 	if err := json.Unmarshal(body, &result); err != nil {
// 		return nil, response.Meta{}, err
// 	}

// 	pagination := &response.Meta{
// 		TotalData: result.Paging.TotalRecord,
// 		Page:      result.Paging.PageCurrent,
// 		Limit:     result.Paging.PageLimit,
// 		TotalPage: result.Paging.PageTotal,
// 	}

// 	return result.Data, *pagination, nil
// }

// func (service *DmsServiceImpl) GetListOutlet(ctx context.Context, dataFilter model.OutletQueryFilter) ([]model.OutletList, response.Meta, error) {
// 	config, err := config.LoadConfig(".")

// 	if err != nil {
// 		log.Fatal("could not load config ", err)
// 	}

// 	endpointUrl := fmt.Sprintf("%s/master/v1/outlets?limit=%s&page=%s&sort=%s&is_active=%s&outlet_type_name=%s&outlet_group_name=%s",
// 		config.KongUrl, dataFilter.Limit, dataFilter.Page, dataFilter.Sort, dataFilter.IsActive, dataFilter.OutletTypeName, dataFilter.OutletGroupName)

// 	req, err := http.NewRequestWithContext(ctx, "GET", endpointUrl, nil)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}
// 	req.Header.Set("cust_id", "C220010001")
// 	req.Header.Set("Accept", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, response.Meta{}, fmt.Errorf("request failed: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, response.Meta{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// 	}

// 	var result struct {
// 		Data   []model.OutletList `json:"data"`
// 		Paging struct {
// 			TotalRecord int `json:"total_record"`
// 			PageCurrent int `json:"page_current"`
// 			PageLimit   int `json:"page_limit"`
// 			PageTotal   int `json:"page_total"`
// 		} `json:"paging"`
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, response.Meta{}, fmt.Errorf("failed to read response body: %w", err)
// 	}

// 	if err := json.Unmarshal(body, &result); err != nil {
// 		return nil, response.Meta{}, fmt.Errorf("failed to unmarshal response body: %w", err)
// 	}

// 	pagination := &response.Meta{
// 		TotalData: result.Paging.TotalRecord,
// 		Page:      result.Paging.PageCurrent,
// 		Limit:     result.Paging.PageLimit,
// 		TotalPage: result.Paging.PageTotal,
// 	}

// 	return result.Data, *pagination, nil

// }

// func (service *DmsServiceImpl) GetListOutlet(ctx context.Context, dataFilter model.OutletQueryFilter, custId string) ([]model.OutletList, response.Meta, error) {
// 	config, err := config.LoadConfig(".")
// 	// log.Printf(config.KongUrl)

// 	if err != nil {
// 		log.Fatal("could not load config ", err)
// 	}

// 	// Add default limit to 15
// 	if dataFilter.Limit == "" {
// 		dataFilter.Limit = "15"
// 	}

// 	endpointUrl := fmt.Sprintf("%s/v1/outlets?limit=%s&page=%s&sort=%s&is_active=%s",
// 		config.KongUrl, dataFilter.Limit, dataFilter.Page, dataFilter.Sort, dataFilter.IsActive)

// 	req, err := http.NewRequestWithContext(ctx, "GET", endpointUrl, nil)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}
// 	req.Header.Add("cust_id", custId)
// 	req.Header.Set("Accept", "application/json")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, response.Meta{}, fmt.Errorf("request failed: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, response.Meta{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// 	}

// 	var result struct {
// 		Data   []model.OutletList `json:"data"`
// 		Paging struct {
// 			TotalRecord int `json:"total_record"`
// 			PageCurrent int `json:"page_current"`
// 			PageLimit   int `json:"page_limit"`
// 			PageTotal   int `json:"page_total"`
// 		} `json:"paging"`
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, response.Meta{}, fmt.Errorf("failed to read response body: %w", err)
// 	}

// 	if err := json.Unmarshal(body, &result); err != nil {
// 		return nil, response.Meta{}, fmt.Errorf("failed to unmarshal response body: %w", err)
// 	}

// 	// apply filter
// 	filteredData := filterOutletByTypeName(result.Data, dataFilter)

// 	pagination := &response.Meta{
// 		TotalData: len(filteredData),
// 		Page:      result.Paging.PageCurrent,
// 		Limit:     result.Paging.PageLimit,
// 		TotalPage: (len(filteredData) + result.Paging.PageLimit - 1) / result.Paging.PageLimit,
// 	}

// 	return filteredData, *pagination, nil
// }

// // new filter for outlet type name and outley group name
// func filterOutletByTypeName(data []model.OutletList, filter model.OutletQueryFilter) []model.OutletList {
// 	var filteredData []model.OutletList

// 	for _, outlet := range data {
// 		if (filter.OutletTypeName == "" || outlet.OtTypeName == filter.OutletTypeName) &&
// 			(filter.OutletGroupName == "" || outlet.OtGrpName == filter.OutletGroupName) {
// 			filteredData = append(filteredData, outlet)
// 		}
// 	}

// 	return filteredData
// }

// func (service *DmsServiceImpl) GetOutletNotAssign(ctx context.Context, dataFilter model.DmsQueryFilter) []model.Outlet {
// 	result := service.outletRepo.GetOutlet(ctx, dataFilter)

// 	var data []model.Outlet
// 	for _, row := range result {
// 		var res model.Outlet
// 		helper.Automapper(row, &res)
// 		data = append(data, res)
// 	}

// 	return data
// }

// func (service *DmsServiceImpl) GetOutletBySalesman(ctx context.Context, dataFilter model.OutletBySalesman, headers map[string]string, custId string) ([]model.OutletNew, response.Meta, error) {
// 	config, err := config.LoadConfig(".")

// 	if err != nil {
// 		log.Fatal("could not load config", err)
// 	}

// 	// find pjpIds by salesman
// 	pjpIds := service.pjpRepo.FindPjpIdBySalesmanCode(ctx, dataFilter.SalesmanCode)
// 	log.Printf("Retreived pjpID: %d", pjpIds)

// 	// find DestinationID by pjpIds
// 	DestinationIDs := service.routeOutletRepo.FindAllDestinationIDByPjpIdToday(ctx, pjpIds)
// 	if len(DestinationIDs) == 0 {
// 		emptyPagination := response.Meta{
// 			TotalData: 0,
// 			Page:      0,
// 			Limit:     0,
// 			TotalPage: 0,
// 		}
// 		return []model.OutletNew{}, emptyPagination, nil
// 	}
// 	log.Printf("Retreived DestinationIDs: %d", DestinationIDs)

// 	var DestinationIDStrings []string
// 	for _, DestinationID := range DestinationIDs {
// 		DestinationIDStrings = append(DestinationIDStrings, strconv.Itoa(DestinationID))
// 	}
// 	DestinationIDsFormatted := strings.Join(DestinationIDStrings, ",")

// 	// remove /master when deploy to staging or production
// 	endpointURL := fmt.Sprintf("%s/v1/outlets?&outlet_id=%s&limit=99", config.KongUrl, DestinationIDsFormatted)
// 	log.Printf("Request URL: %s", endpointURL)

// 	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}
// 	req.Header.Add("cust_id", custId)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}
// 	defer resp.Body.Close()

// 	var result struct {
// 		Data   []model.OutletNew `json:"data"`
// 		Paging struct {
// 			TotalRecord int `json:"total_record"`
// 			PageCurrent int `json:"page_current"`
// 			PageLimit   int `json:"page_limit"`
// 			PageTotal   int `json:"page_total"`
// 		} `json:"paging"`
// 	}
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}

// 	if err := json.Unmarshal(body, &result); err != nil {
// 		return nil, response.Meta{}, err
// 	}

// 	pagination := &response.Meta{
// 		TotalData: result.Paging.TotalRecord,
// 		Page:      result.Paging.PageCurrent,
// 		Limit:     result.Paging.PageLimit,
// 		TotalPage: result.Paging.PageTotal,
// 	}

// 	return result.Data, *pagination, nil
// }

// func (service *DmsServiceImpl) GetOutletBySalesmanId(ctx context.Context, dataFilter model.OutletBySalesmanId, headers map[string]string, custId string) ([]model.OutletNew, response.Meta, error) {
// 	config, err := config.LoadConfig(".")
// 	if err != nil {
// 		log.Fatal("could not load config", err)
// 	}

// 	// Find pjpId by salesman
// 	pjpId := service.pjpRepo.FindPjpIdBySalesmanId(ctx, dataFilter.SalesmanId)

// 	// Convert pjpId to []int
// 	pjpIds := []int{pjpId}

// 	// Determine DestinationIDs based on dataFilter.Search
// 	var DestinationIDs []int
// 	if dataFilter.Search != "" {
// 		log.Printf("Search filter applied: %s", dataFilter.Search)
// 		DestinationIDs = service.routeOutletRepo.SearchDestinationIDByPjpId(ctx, pjpIds, dataFilter.Search)
// 	} else {
// 		log.Printf("No search filter applied")
// 		DestinationIDs = service.routeOutletRepo.FindAllDestinationIDByPjpId(ctx, pjpIds)
// 	}

// 	// Handle empty DestinationIDs
// 	if len(DestinationIDs) == 0 {
// 		emptyPagination := response.Meta{
// 			TotalData: 0,
// 			Page:      0,
// 			Limit:     0,
// 			TotalPage: 0,
// 		}
// 		return []model.OutletNew{}, emptyPagination, nil
// 	}
// 	log.Printf("Retrieved DestinationIDs: %d", DestinationIDs)

// 	// Convert DestinationIDs to comma-separated string
// 	var DestinationIDStrings []string
// 	for _, DestinationID := range DestinationIDs {
// 		DestinationIDStrings = append(DestinationIDStrings, strconv.Itoa(DestinationID))
// 	}
// 	DestinationIDsFormatted := strings.Join(DestinationIDStrings, ",")

// 	// Remove /master when deploy to staging or production
// 	endpointURL := fmt.Sprintf("%s/v1/outlets?&outlet_id=%s", config.KongUrl, DestinationIDsFormatted)
// 	log.Printf("Request URL: %s", endpointURL)

// 	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}
// 	req.Header.Add("cust_id", custId)
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}
// 	defer resp.Body.Close()

// 	var result struct {
// 		Data   []model.OutletNew `json:"data"`
// 		Paging struct {
// 			TotalRecord int `json:"total_record"`
// 			PageCurrent int `json:"page_current"`
// 			PageLimit   int `json:"page_limit"`
// 			PageTotal   int `json:"page_total"`
// 		} `json:"paging"`
// 	}
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, response.Meta{}, err
// 	}

// 	if err := json.Unmarshal(body, &result); err != nil {
// 		return nil, response.Meta{}, err
// 	}

// 	pagination := &response.Meta{
// 		TotalData: result.Paging.TotalRecord,
// 		Page:      result.Paging.PageCurrent,
// 		Limit:     result.Paging.PageLimit,
// 		TotalPage: result.Paging.PageTotal,
// 	}

// 	return result.Data, *pagination, nil
// }

// func (service *DmsServiceImpl) GetListSalesman(ctx context.Context, custId string) (response.ListSalesmanAPIResponse, error) {
// 	var result response.ListSalesmanAPIResponse
// 	var pageLimit int = 9999
// 	var page int = 1

// 	config, err := config.LoadConfig(".")
// 	if err != nil {
// 		log.Printf("failed to load config: %v", err)
// 		return result, err
// 	}

// 	endpointURL := fmt.Sprintf("%s/v1/salesman?limit=%d&page=%d&is_active=%d",
// 		config.KongUrl, pageLimit, page, 1)
// 	log.Printf("Requesting salesman list from URL: %s", endpointURL)

// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL, nil)
// 	if err != nil {
// 		log.Printf("failed to create HTTP request: %v", err)
// 		return result, err
// 	}
// 	req.Header.Set("cust_id", custId)

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Printf("failed to perform HTTP request: %v", err)
// 		return result, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		log.Printf("received non-OK HTTP status: %d", resp.StatusCode)
// 		return result, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Printf("failed to read response body: %v", err)
// 		return result, err
// 	}

// 	if err := json.Unmarshal(body, &result); err != nil {
// 		log.Printf("failed to unmarshal response: %v", err)
// 		return result, err
// 	}

// 	// Ambil PJP aktif
// 	pjpList, err := service.pjpRepo.GetActivePjp(ctx, custId)
// 	if err != nil {
// 		log.Printf("failed to get active PJP: %v", err)
// 		return result, err
// 	}

// 	if len(pjpList) == 0 {
// 		// Kosongkan data dan paging
// 		result.Data = []response.ListSalesmanResponse{}
// 		result.Paging.TotalRecord = 0
// 		result.Paging.PageTotal = 0
// 		result.Paging.PageCurrent = page
// 		result.Paging.PageLimit = pageLimit
// 		return result, nil
// 	}

// 	// Buat map salesmanId untuk filter cepat
// 	salesmanMap := make(map[int]struct{})
// 	for _, pjp := range pjpList {
// 		salesmanMap[pjp.SalesManID] = struct{}{}
// 	}

// 	// Filter salesman
// 	filtered := make([]response.ListSalesmanResponse, 0)
// 	for _, s := range result.Data {
// 		if _, exists := salesmanMap[s.EmpID]; exists {
// 			filtered = append(filtered, s)
// 		}
// 	}

// 	// Update result
// 	result.Data = filtered
// 	result.Paging.TotalRecord = len(filtered)
// 	result.Paging.PageCurrent = page
// 	result.Paging.PageLimit = pageLimit
// 	result.Paging.PageTotal = 1

// 	return result, nil
// }
