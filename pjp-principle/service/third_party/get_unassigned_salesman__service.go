package thirdparty

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"scyllax-pjp/config"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
)

func (service *thirdPartyService) GetUnassignedSalesman(ctx context.Context, headers map[string]string, dataFilter request.SalesmanListQueryFilter, custId string) ([]model.NewSalesman, response.Meta) {
	config, err := config.LoadConfig(".")
	helper.ErrorPanic(err)

	endpointURL := masterSalesmanEndpointURL(config.KongUrl, dataFilter)
	log.Printf("Endpoint URL: %s", endpointURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL, nil)
	helper.ErrorPanic(err)

	req.Header.Set("cust_id", custId)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	helper.ErrorPanic(err)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-OK status: %v", resp.StatusCode)
		helper.ErrorPanic(fmt.Errorf("Non-OK status: %v", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	helper.ErrorPanic(err)

	var result struct {
		Data   []model.NewSalesman `json:"data"`
		Paging struct {
			TotalRecord int `json:"total_record"`
			PageCurrent int `json:"page_current"`
			PageLimit   int `json:"page_limit"`
			PageTotal   int `json:"page_total"`
		} `json:"paging"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("JSON Unmarshal Error: %v", err)
		helper.ErrorPanic(err)
	}

	setOperationTypeName(result.Data)

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}

	defer helper.CommitOrRollback(tx)

	pjps := service.pjpRepo.GetPjpWithRoute(ctx, tx, "", custId)
	salesmanIDSet := make(map[int]struct{})
	for _, pjp := range pjps {
		salesmanIDSet[pjp.SalesManID] = struct{}{}
	}

	// Convert ke slice jika perlu
	salesmanIDs := make([]int, 0, len(salesmanIDSet))
	for id := range salesmanIDSet {
		salesmanIDs = append(salesmanIDs, id)
	}

	filteredSalesmen := filterUnassignedSalesmen(result.Data, salesmanIDs)
	pagination := generatePagination(result.Paging.PageCurrent, result.Paging.PageLimit, len(filteredSalesmen))

	return filteredSalesmen, pagination
}
