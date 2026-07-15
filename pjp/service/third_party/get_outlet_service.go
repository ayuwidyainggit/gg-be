package thirdparty

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"scyllax-pjp/config"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
)

func (service *thirdPartyService) GetOutlet(ctx context.Context, dataFilter model.DmsQueryFilter, custId string) ([]model.Outlet, response.Meta) {
	config, err := config.LoadConfig(".")
	helper.ErrorPanic(err)

	endpointURL := masterOutletEndpointURL(config.KongUrl, dataFilter)

	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
	helper.ErrorPanic(err)

	req.Header.Add("cust_id", custId)
	resp, err := http.DefaultClient.Do(req)
	helper.ErrorPanic(err)

	defer resp.Body.Close()

	var result model.OutletAPIResponse
	body, err := io.ReadAll(resp.Body)
	helper.ErrorPanic(err)

	if err := json.Unmarshal(body, &result); err != nil {
		helper.ErrorPanic(err)
	}

	pagination := &response.Meta{
		TotalData: result.Paging.TotalRecord,
		Page:      result.Paging.PageCurrent,
		Limit:     result.Paging.PageLimit,
		TotalPage: result.Paging.PageTotal,
	}

	return model.ConvertOutletAPIsToOutlets(result.Data), *pagination
}
