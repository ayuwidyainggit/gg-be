package thirdparty

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"scyllax-pjp/config"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
)

func (service *thirdPartyService) GetAssignedSalesman(ctx context.Context, custId string) response.ListSalesmanAPIResponse {
	var result response.ListSalesmanAPIResponse
	var pageLimit int = 9999
	var page int = 1

	config, err := config.LoadConfig(".")
	helper.ErrorPanic(err)

	endpointURL := fmt.Sprintf("%s/v1/salesman?limit=%d&page=%d&is_active=%d",
		config.KongUrl, pageLimit, page, 1)
	log.Printf("Requesting salesman list from URL: %s", endpointURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL, nil)
	helper.ErrorPanic(err)

	req.Header.Set("cust_id", custId)

	resp, err := http.DefaultClient.Do(req)
	helper.ErrorPanic(err)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("received non-OK HTTP status: %d", resp.StatusCode)
		helper.ErrorPanic(fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	helper.ErrorPanic(err)

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("failed to unmarshal response: %v", err)
		helper.ErrorPanic(err)
	}

	// Ambil PJP aktif
	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}

	defer helper.CommitOrRollback(tx)

	pjps := service.pjpRepo.GetPjpWithRoute(ctx, tx, "", custId)

	if len(pjps) == 0 {
		// Kosongkan data dan paging
		result.Data = []response.ListSalesmanResponse{}
		result.Paging.TotalRecord = 0
		result.Paging.PageTotal = 0
		result.Paging.PageCurrent = page
		result.Paging.PageLimit = pageLimit
		return result
	}

	// Buat map salesmanId ke outlet_ids
	salesmanOutletMap := make(map[int][]int)
	for _, pjp := range pjps {
		if pjp.OutletID != 0 {
			salesmanOutletMap[pjp.SalesManID] = append(salesmanOutletMap[pjp.SalesManID], pjp.OutletID)
		}
	}

	// Filter salesman dan tambahkan outlet_ids
	filtered := make([]response.ListSalesmanResponse, 0)
	for _, s := range result.Data {
		if outletIDs, exists := salesmanOutletMap[s.EmpID]; exists {
			s.OutletIDs = outletIDs
			filtered = append(filtered, s)
		}
	}

	// Update result
	result.Data = filtered
	result.Paging.TotalRecord = len(filtered)
	result.Paging.PageCurrent = page
	result.Paging.PageLimit = pageLimit
	result.Paging.PageTotal = 1

	return result
}
