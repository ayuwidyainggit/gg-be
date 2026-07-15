package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"scyllax-tms/config"
	"scyllax-tms/entity"
	"scyllax-tms/exception"
	"scyllax-tms/helper"
	"scyllax-tms/repository"
)

type DmsService interface {
	GetVehicleByDms(ctx context.Context, dataFilter entity.GeneralQueryFilter) ([]entity.VehicleResponse, entity.Meta, error)
	GetRejectReason(ctx context.Context, dataFilter entity.GeneralQueryFilter) ([]entity.ReasonRejectResponse, entity.Meta, error)
	GetListInvoice(ctx context.Context, dataFilter entity.GeneralQueryFilter, header map[string]string) ([]entity.CustomShipmentInvoice, entity.Meta, error)
	GetReturns(ctx context.Context, dataFilter entity.GeneralQueryFilter, headers map[string]string) ([]entity.ResponseReturn, entity.Meta, error)
}

type DmsServiceImpl struct {
	shipmentRepo repository.ShipmentInvoicesRepo
}

func NewDmsServiceImpl(shipmentRepo repository.ShipmentInvoicesRepo) DmsService {
	return &DmsServiceImpl{
		shipmentRepo: shipmentRepo,
	}
}

func (service *DmsServiceImpl) GetVehicleByDms(ctx context.Context, dataFilter entity.GeneralQueryFilter) ([]entity.VehicleResponse, entity.Meta, error) {
	config, err := config.LoadConfig(".")

	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	endpointURL := fmt.Sprintf("%s/v1/vehicles?limit=%d&page=%d&is_active=%d", config.KongUrl, dataFilter.Limit, dataFilter.Page, dataFilter.IsActive)
	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
	req.Header.Set("cust_id", "C220010001")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, entity.Meta{}, err
	}
	defer resp.Body.Close()

	var response struct {
		Data   []entity.VehicleResponse `json:"data"`
		Paging struct {
			TotalRecord int `json:"total_record"`
			PageCurrent int `json:"page_current"`
			PageLimit   int `json:"page_limit"`
			PageTotal   int `json:"page_total"`
		} `json:"paging"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, entity.Meta{}, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, entity.Meta{}, err
	}

	pagination := &entity.Meta{
		TotalData: response.Paging.TotalRecord,
		Page:      response.Paging.PageCurrent,
		Limit:     response.Paging.PageLimit,
		TotalPage: response.Paging.PageTotal,
	}

	return response.Data, *pagination, nil
}

func (service *DmsServiceImpl) GetRejectReason(ctx context.Context, dataFilter entity.GeneralQueryFilter) ([]entity.ReasonRejectResponse, entity.Meta, error) {
	config, err := config.LoadConfig(".")

	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	endpointURL := fmt.Sprintf("%s/v1/reject-reason?limit=%d&page=%d&is_active=%d", config.KongUrl, dataFilter.Limit, dataFilter.Page, dataFilter.IsActive)
	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
	req.Header.Set("cust_id", "C220010001")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, entity.Meta{}, err
	}
	defer resp.Body.Close()

	var response struct {
		Data   []entity.ReasonRejectResponse `json:"data"`
		Paging struct {
			TotalRecord int `json:"total_record"`
			PageCurrent int `json:"page_current"`
			PageLimit   int `json:"page_limit"`
			PageTotal   int `json:"page_total"`
		} `json:"paging"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, entity.Meta{}, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, entity.Meta{}, err
	}

	pagination := &entity.Meta{
		TotalData: response.Paging.TotalRecord,
		Page:      response.Paging.PageCurrent,
		Limit:     response.Paging.PageLimit,
		TotalPage: response.Paging.PageTotal,
	}

	return response.Data, *pagination, nil
}

func (service *DmsServiceImpl) GetListInvoice(ctx context.Context, dataFilter entity.GeneralQueryFilter, headers map[string]string) ([]entity.CustomShipmentInvoice, entity.Meta, error) {
	config, err := config.LoadConfig(".")

	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	// fmt.Println(config.KongUrlSales)

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 9999
	}

	var endpointUrl string

	// Remove /sales when deploy to statging
	if len(dataFilter.OutletId) > 0 {
		endpointUrl = fmt.Sprintf("%s/v1/invoices?limit=%d&status=2&status=6&outlet_id=%s", config.KongUrlSales, dataFilter.Limit, dataFilter.OutletId)
	} else {
		endpointUrl = fmt.Sprintf("%s/v1/invoices?limit=%d&status=2&status=6", config.KongUrlSales, dataFilter.Limit)
	}
	// log.Printf("Request URL: %s", endpointUrl)
	req, err := http.NewRequestWithContext(ctx, "GET", endpointUrl, nil)
	if err != nil {
		return nil, entity.Meta{}, err
	}
	// req.Header.Set("cust_id", "C220010001")
	// req.Header.Set("Accept", "application/json")
	// req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoiQzIyMDAxMDAwMSIsImRpc3RfcHJpY2VfZ3JwX2lkIjoxLCJlbWFpbCI6ImRpc3RAc2RhLmlkZXRhbWEuaWQiLCJlbXBfaWQiOjAsImV4cGlyZXMiOjE3MjI4OTgwNjgsImlzX2FkbWluIjp0cnVlLCJsYW5nX2lkIjoiZW4iLCJtb2JpbGVfbm8iOiIwODU3NDgxMjMxMjgiLCJwYXJlbnRfY3VzdF9pZCI6IkMyMjAwMSIsInVzZXJfZnVsbG5hbWUiOiJEaXN0IElERSBTZGEiLCJ1c2VyX2lkIjoxMiwidXNlcl9uYW1lIjoiZGlzdEBzZGEuaWRldGFtYS5pZCIsIndoYXRzYXBwIjoiMDg1NzQ4MTIzMTI4In0.yZBTzRMzO8M_Qb0Z-FHkbO6IQ171jXv4FvIzfg0lmSg")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, entity.Meta{}, err
	}

	defer resp.Body.Close()

	// log.Printf("Response status: %s", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.Printf("Error reading response body: %v", err)
		return nil, entity.Meta{}, err
	}

	// log.Printf("Response body: %s", string(body))

	var response struct {
		Message string                           `json:"message"`
		Data    []entity.ResponseShipmentInvoive `json:"data"`
		Paging  struct {
			TotalRecord int `json:"total_record"`
			PageCurrent int `json:"page_current"`
			PageLimit   int `json:"page_limit"`
			PageTotal   int `json:"page_total"`
		} `json:"paging"`
		RequestID string `json:"request_id"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshall response: %v", err)
		return nil, entity.Meta{}, err
	}

	// log.Println("==============================================================================")
	// log.Printf("Parsed response: %+v", response)

	assignedInvoiceNo := service.shipmentRepo.FindAllByInvoiceNo(ctx)
	// log.Printf("Assigned invoice: %+v", assignedInvoiceNo)

	assignedInvoiceNoSet := make(map[string]struct{})
	for _, no := range assignedInvoiceNo {
		assignedInvoiceNoSet[no] = struct{}{}
	}

	filteredInvoice := make([]entity.CustomShipmentInvoice, 0)
	for _, invoice := range response.Data {
		log.Printf("Invocie number: %+v", invoice.InvoiceNo)
		if invoice.InvoiceNo == "" || invoice.InvoiceNo == "null" {
			// log.Printf("Invoice number is empty or null, including in filtered results.")
			var customInvoice entity.CustomShipmentInvoice
			helper.Automapper(invoice, &customInvoice)

			customInvoice.TotalPromo = invoice.PromoValue
			customInvoice.TotalBruto = invoice.SubTotal
			customInvoice.TotalPPN = invoice.VatValue
			customInvoice.TotalNetto = invoice.Total

			filteredInvoice = append(filteredInvoice, customInvoice)
		} else if _, exists := assignedInvoiceNoSet[invoice.InvoiceNo]; !exists {
			// Process as usual if invoice is not assigned
			// log.Printf("Invoice number is not assigned, adding to filtered results.")
			var customInvoice entity.CustomShipmentInvoice
			helper.Automapper(invoice, &customInvoice)

			customInvoice.TotalPromo = invoice.PromoValue
			customInvoice.TotalBruto = invoice.SubTotal
			customInvoice.TotalPPN = invoice.VatValue
			customInvoice.TotalNetto = invoice.Total

			filteredInvoice = append(filteredInvoice, customInvoice)
		}
	}

	// log.Printf("Filtered response: %+v", filteredInvoice)

	totalRecords := len(filteredInvoice)
	totalPages := 0
	if response.Paging.PageLimit > 0 {
		totalPages = (totalRecords + response.Paging.PageLimit - 1) / response.Paging.PageLimit
	}

	pagination := &entity.Meta{
		TotalData: totalRecords,
		Page:      response.Paging.PageCurrent,
		Limit:     response.Paging.PageLimit,
		TotalPage: totalPages,
	}

	return filteredInvoice, *pagination, nil

}

func (service *DmsServiceImpl) GetReturns(ctx context.Context, dataFilter entity.GeneralQueryFilter, headers map[string]string) ([]entity.ResponseReturn, entity.Meta, error) {
	config, err := config.LoadConfig(".")

	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 9999
	}

	// Remove /sales when deploy to staging or production
	var endpointURL string

	if len(dataFilter.OutletId) > 0 {
		endpointURL = fmt.Sprintf(
			"%s/v1/returns?limit=%d&sort=return_no:desc&status=3&mode=shipment&outlet_id=%s",
			config.KongUrlSales, dataFilter.Limit, dataFilter.OutletId,
		)
	} else {
		endpointURL = fmt.Sprintf(
			"%s/v1/returns?limit=%d&sort=return_no:desc&status=3&mode=shipment",
			config.KongUrlSales, dataFilter.Limit,
		)
	}
	
	// log.Printf("Request URL: %s", endpointURL)
	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
	if err != nil {
		return nil, entity.Meta{}, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, entity.Meta{}, err
	}
	defer resp.Body.Close()

	// log.Printf("Response status: %s", resp.Status)

	var response struct {
		Data   []entity.ResponseReturn `json:"data"`
		Paging struct {
			TotalRecord int `json:"total_record"`
			PageCurrent int `json:"page_current"`
			PageLimit   int `json:"page_limit"`
			PageTotal   int `json:"page_total"`
		} `json:"paging"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, entity.Meta{}, err
	}

	// log.Printf("Response body: %s", string(body))

	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshall response: %v", err)
		return nil, entity.Meta{}, err
	}

	// log.Printf("Parsed response: %+v", response)

	pagination := &entity.Meta{
		TotalData: response.Paging.TotalRecord,
		Page:      response.Paging.PageCurrent,
		Limit:     response.Paging.PageLimit,
		TotalPage: response.Paging.PageTotal,
	}

	return response.Data, *pagination, nil
}
