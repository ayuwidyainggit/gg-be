package thirdparty

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"scyllax-pjp/config"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
)

func (service *thirdPartyService) GetSalesmanByID(ctx context.Context, empId int, headers map[string]string, custId string) model.NewSalesman {
	config, err := config.LoadConfig(".")
	helper.ErrorPanic(err)

	// kong public add /master for dev remove for staging
	endpointURL := fmt.Sprintf("%s/v1/salesman/%d",
		config.KongUrl, empId)
	log.Printf("Endpoint URL: %s", endpointURL)

	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
	helper.ErrorPanic(err)

	req.Header.Set("cust_id", custId)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	helper.ErrorPanic(err)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTP request failed with status code %d", resp.StatusCode)
		helper.ErrorPanic(fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	helper.ErrorPanic(err)

	var result struct {
		Message   string            `json:"message"`
		Data      model.NewSalesman `json:"data"`
		RequestID string            `json:"request_id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("JSON Unmarshal Error: %v", err)
		helper.ErrorPanic(err)
	}

	return result.Data
}
