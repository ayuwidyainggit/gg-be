package thirdparty

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"scyllax-pjp/config"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"strconv"
	"strings"
)

func (service *thirdPartyService) GetOutletBySalesCodes(ctx context.Context, dataFilter model.OutletBySalesman, headers map[string]string, custId string) ([]model.OutletNew, response.Meta, error) {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("could not load config", err)
	}

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	pjps := service.pjpRepo.GetPjpsByEmpCodes(ctx, tx, dataFilter.SalesmanCode, custId)

	pjpIdSet := make(map[int]struct{}, len(pjps))
	for _, pjp := range pjps {
		pjpIdSet[pjp.ID] = struct{}{}
	}

	pjpIds := make([]int, 0, len(pjpIdSet))
	for id := range pjpIdSet {
		pjpIds = append(pjpIds, id)
	}

	outlets := service.routeOutletHistoryRepo.FindByPjpIdToday(ctx, tx, pjpIds, custId)
	if len(outlets) == 0 {
		emptyPagination := response.Meta{
			TotalData: 0,
			Page:      0,
			Limit:     0,
			TotalPage: 0,
		}
		return []model.OutletNew{}, emptyPagination, nil
	}

	idSet := make(map[int]struct{})
	outletIdStrings := make([]string, 0, len(outlets))
	for _, outlet := range outlets {
		if _, exists := idSet[outlet.OutletID]; !exists {
			idSet[outlet.OutletID] = struct{}{}
			outletIdStrings = append(outletIdStrings, strconv.Itoa(outlet.OutletID))
		}
	}

	outletIdsFormatted := strings.Join(outletIdStrings, ",")
	endpointURL := masterOutletByIDsEndpointURL(config.KongUrl, outletIdsFormatted, 9999, 1)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, response.Meta{}, err
	}
	req.Header.Add("cust_id", custId)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, response.Meta{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, response.Meta{}, err
	}

	var result model.OutletAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, response.Meta{}, err
	}

	pagination := response.Meta{
		TotalData: result.Paging.TotalRecord,
		Page:      result.Paging.PageCurrent,
		Limit:     result.Paging.PageLimit,
		TotalPage: result.Paging.PageTotal,
	}

	return model.ConvertOutletAPIsToOutletNews(result.Data), pagination, nil
}
