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
	"scyllax-tms/model"
	"scyllax-tms/repository"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type PicklistService interface {
	Create(ctx context.Context, request entity.CreatePicklistRequest, customerId string) (entity.PicklistResponse, error)
	Update(ctx context.Context, request entity.UpdatePicklistRequest) (entity.PicklistResponse, error)
	Delete(ctx context.Context, request entity.DeletePicklistRequest) (entity.PicklistResponse, error)
	GetPicklist(ctx context.Context, request entity.GetPicklistRequest) (entity.PicklistResponse, error)
	GetAll(ctx context.Context, dataFilter entity.GeneralQueryFilter, picklistFilter entity.PicklistFilter, customerId string) ([]entity.PicklistResponse, entity.Meta, error)
	GetListInvoice(ctx context.Context, dataFilter entity.GeneralQueryFilter, header map[string]string, customerId string) ([]entity.CustomShipmentInvoice, entity.Meta, error)
}

type PicklistServiceImpl struct {
	picklistRepo repository.PicklistRepository
	validate     *validator.Validate
}

func NewPicklistServiceImpl(picklistRepo repository.PicklistRepository, validate *validator.Validate) PicklistService {
	return &PicklistServiceImpl{picklistRepo: picklistRepo, validate: validate}
}

func (s *PicklistServiceImpl) GeneratePicklistNo(ctx context.Context) (string, error) {
	// Get the current date
	now := time.Now()
	yy := now.Format("06") // Year in two digits
	mm := now.Format("01") // Month in two digits
	dd := now.Format("02") // Day in two digits

	// Get the current sequence number for the day
	seq, err := s.picklistRepo.GetNextSequence(ctx, yy, mm, dd)
	if err != nil {
		return "", err
	}

	// Format the picklist number
	picklistNo := fmt.Sprintf("PL%s%s%s%04d", yy, mm, dd, seq)
	return picklistNo, nil
}

func (s *PicklistServiceImpl) Create(ctx context.Context, request entity.CreatePicklistRequest, customerId string) (entity.PicklistResponse, error) {
	// Start a new transaction
	tx, err := s.picklistRepo.BeginTransaction(ctx)
	if err != nil {
		return entity.PicklistResponse{}, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	err = s.validate.Struct(request)
	if err != nil {
		return entity.PicklistResponse{}, err
	}

	picklistNo, err := s.GeneratePicklistNo(ctx)
	fmt.Println("data", picklistNo)
	if err != nil {
		return entity.PicklistResponse{}, err
	}

	picklist := model.Picklist{
		PicklistNo: picklistNo,
		CustId:     customerId,
		Driver:     request.Driver,
		Helper:     request.Helper,
		Vehicle:    request.Vehicle,
		UpdatedBy:  request.UpdatedBy,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = s.picklistRepo.InsertPicklist(ctx, tx, picklist)
	if err != nil {
		return entity.PicklistResponse{}, err
	}

	for _, order := range request.Orders {
		var dueDate *time.Time
		if order.DueDate != "" {
			parsedDueDate, err := helper.ParseDate(order.DueDate)
			if err != nil {
				return entity.PicklistResponse{}, err
			}
			dueDate = &parsedDueDate
		}

		// Parse InvoiceDate string to time.Time
		var invoiceDate *time.Time
		if order.InvoiceDate != "" {
			parsedDate, err := helper.ParseDate(order.InvoiceDate)
			if err != nil {
				return entity.PicklistResponse{}, err
			}
			invoiceDate = &parsedDate
		}

		orderPicklist := model.OrderPicklist{
			OrderNo:     order.OrderNo,
			CustId:      customerId,
			PicklistNo:  picklistNo,
			InvoiceNo:   order.InvoiceNo,
			OutletName:  order.OutletName,
			Salesman:    order.Salesman,
			InvoiceDate: invoiceDate, // Now using parsed time.Time
			DueDate:     dueDate,
			TotalPrice:  order.TotalPrice,
			Ppn:         order.Ppn,
			Discount:    order.Discount,
			TotalUnpaid: order.TotalUnpaid,
			TotalPromo:  order.TotalPromo,
			PaymentType: order.PaymentType,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		orderId, err := s.picklistRepo.InsertOrderPicklist(ctx, tx, orderPicklist)
		if err != nil {
			return entity.PicklistResponse{}, err
		}

		for _, orderProduct := range order.Products {
			orderProductModel := model.OrderProduct{
				OrderID:       uint(orderId),
				CustId:        customerId,
				ProductName:   orderProduct.ProductName,
				ProductCode:   orderProduct.ProductCode,
				ProductId:     orderProduct.ProductId,
				Quantity1:     orderProduct.Quantity1,
				Quantity2:     orderProduct.Quantity2,
				Quantity3:     orderProduct.Quantity3,
				QuantityUnit1: orderProduct.QuantityUnit1,
				QuantityUnit2: orderProduct.QuantityUnit2,
				QuantityUnit3: orderProduct.QuantityUnit3,
				// Volume:        orderProduct.Volume,
				// Weight:        orderProduct.Weight,
				Unit1Price: orderProduct.Unit1Price,
				Unit2Price: orderProduct.Unit2Price,
				Unit3Price: orderProduct.Unit3Price,
				Ppn:        orderProduct.Ppn,
				Volume1:    orderProduct.Volume1,
				Volume2:    orderProduct.Volume2,
				Volume3:    orderProduct.Volume3,
				Weight1:    orderProduct.Weight1,
				Weight2:    orderProduct.Weight2,
				Weight3:    orderProduct.Weight3,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			err = s.picklistRepo.InsertOrderProduct(ctx, tx, orderProductModel)
			if err != nil {
				return entity.PicklistResponse{}, err
			}
		}
	}

	return entity.PicklistResponse{
		PicklistNo: picklist.PicklistNo,
		Driver:     picklist.Driver,
		Helper:     picklist.Helper,
		Vehicle:    picklist.Vehicle,
		UpdatedBy:  picklist.UpdatedBy,
		CreatedAt:  picklist.CreatedAt,
		UpdatedAt:  picklist.UpdatedAt,
	}, nil
}

func (s *PicklistServiceImpl) Update(ctx context.Context, request entity.UpdatePicklistRequest) (entity.PicklistResponse, error) {
	err := s.validate.Struct(request)
	if err != nil {
		return entity.PicklistResponse{}, err
	}

	picklist, err := s.picklistRepo.FindByID(ctx, request.PicklistNo)
	if err != nil {
		return entity.PicklistResponse{}, err
	}

	picklist.Driver = request.Driver
	picklist.Helper = request.Helper
	picklist.Vehicle = request.Vehicle
	picklist.UpdatedBy = request.UpdatedBy
	picklist.UpdatedAt = time.Now()

	err = s.picklistRepo.UpdatePicklist(ctx, picklist)
	if err != nil {
		return entity.PicklistResponse{}, err
	}

	return entity.PicklistResponse{
		PicklistNo: picklist.PicklistNo,
		Driver:     picklist.Driver,
		Helper:     picklist.Helper,
		Vehicle:    picklist.Vehicle,
		UpdatedBy:  picklist.UpdatedBy,
		CreatedAt:  picklist.CreatedAt,
		UpdatedAt:  picklist.UpdatedAt,
	}, nil
}

func (s *PicklistServiceImpl) Delete(ctx context.Context, request entity.DeletePicklistRequest) (entity.PicklistResponse, error) {
	err := s.validate.Struct(request)
	if err != nil {
		return entity.PicklistResponse{}, err
	}

	picklist, err := s.picklistRepo.FindByID(ctx, request.PicklistNo)
	if err != nil {
		return entity.PicklistResponse{}, err
	}

	err = s.picklistRepo.DeletePicklist(ctx, request.PicklistNo)
	if err != nil {
		return entity.PicklistResponse{}, err
	}

	return entity.PicklistResponse{
		PicklistNo: picklist.PicklistNo,
		Driver:     picklist.Driver,
		Helper:     picklist.Helper,
		Vehicle:    picklist.Vehicle,
		UpdatedBy:  picklist.UpdatedBy,
		CreatedAt:  picklist.CreatedAt,
		UpdatedAt:  picklist.UpdatedAt,
	}, nil
}

func (s *PicklistServiceImpl) GetPicklist(ctx context.Context, request entity.GetPicklistRequest) (entity.PicklistResponse, error) {
	err := s.validate.Struct(request)
	if err != nil {
		return entity.PicklistResponse{}, err
	}

	picklist, err := s.picklistRepo.FindByID(ctx, request.PicklistNo)
	if err != nil {
		return entity.PicklistResponse{}, err
	}

	// Fetch orders and products associated with the picklist
	orders, err := s.picklistRepo.GetOrdersByPicklistNo(ctx, request.PicklistNo)
	if err != nil {
		return entity.PicklistResponse{}, err
	}
	fmt.Println(orders)

	var orderResponses []entity.OrderResponse
	for _, order := range orders {
		products, err := s.picklistRepo.GetProductsByOrderID(ctx, order.ID)
		if err != nil {
			return entity.PicklistResponse{}, err
		}

		var productResponses []entity.OrderProductResponse
		for _, product := range products {
			productResponses = append(productResponses, entity.OrderProductResponse{
				ProductName:   product.ProductName,
				ProductCode:   product.ProductCode,
				ProductId:     product.ProductId,
				Quantity1:     product.Quantity1,
				Quantity2:     product.Quantity2,
				Quantity3:     product.Quantity3,
				QuantityUnit1: product.QuantityUnit1,
				QuantityUnit2: product.QuantityUnit2,
				QuantityUnit3: product.QuantityUnit3,
				Volume1:       product.Volume1,
				Volume2:       product.Volume2,
				Volume3:       product.Volume3,
				Weight1:       product.Weight1,
				Weight2:       product.Weight2,
				Weight3:       product.Weight3,
				Unit1Price:    product.Unit1Price,
				Unit2Price:    product.Unit2Price,
				Unit3Price:    product.Unit3Price,
				Ppn:           product.Ppn,
			})
		}

		orderResponses = append(orderResponses, entity.OrderResponse{
			OrderNo:     order.OrderNo,
			InvoiceNo:   order.InvoiceNo,
			OutletName:  order.OutletName,
			Salesman:    order.Salesman,
			InvoiceDate: order.InvoiceDate,
			DueDate:     order.DueDate,
			TotalPrice:  order.TotalPrice,
			Ppn:         order.Ppn,
			Discount:    order.Discount,
			TotalUnpaid: order.TotalUnpaid,
			TotalPromo:  order.TotalPromo,
			PaymentType: order.PaymentType,
			Products:    productResponses,
		})
	}

	return entity.PicklistResponse{
		PicklistNo: picklist.PicklistNo,
		Driver:     picklist.Driver,
		Helper:     picklist.Helper,
		Vehicle:    picklist.Vehicle,
		UpdatedBy:  picklist.UpdatedBy,
		CreatedAt:  picklist.CreatedAt,
		UpdatedAt:  picklist.UpdatedAt,
		Orders:     orderResponses,
	}, nil
}

func (s *PicklistServiceImpl) GetAll(ctx context.Context, dataFilter entity.GeneralQueryFilter, picklistFilter entity.PicklistFilter, customerId string) ([]entity.PicklistResponse, entity.Meta, error) {
	// First, get total count of all records without pagination
	totalData, err := s.picklistRepo.CountAll(ctx, picklistFilter, customerId)
	if err != nil {
		return nil, entity.Meta{}, err
	}

	// Calculate total pages
	totalPage := (totalData + int64(dataFilter.Limit) - 1) / int64(dataFilter.Limit)

	// Get paginated picklists
	picklists, err := s.picklistRepo.FindAll(ctx, dataFilter, picklistFilter, customerId)
	if err != nil {
		return nil, entity.Meta{}, err
	}

	var responses []entity.PicklistResponse
	for _, picklist := range picklists {
		orders, err := s.picklistRepo.GetOrdersByPicklistNo(ctx, picklist.PicklistNo)
		if err != nil {
			return nil, entity.Meta{}, err
		}

		var orderResponses []entity.OrderResponse
		for _, order := range orders {
			products, err := s.picklistRepo.GetProductsByOrderID(ctx, order.ID)
			if err != nil {
				return nil, entity.Meta{}, err
			}

			var productResponses []entity.OrderProductResponse
			for _, product := range products {
				productResponses = append(productResponses, entity.OrderProductResponse{
					ProductName:   product.ProductName,
					ProductCode:   product.ProductCode,
					ProductId:     product.ProductId,
					Quantity1:     product.Quantity1,
					Quantity2:     product.Quantity2,
					Quantity3:     product.Quantity3,
					QuantityUnit1: product.QuantityUnit1,
					QuantityUnit2: product.QuantityUnit2,
					QuantityUnit3: product.QuantityUnit3,
					Unit1Price:    product.Unit1Price,
					Unit2Price:    product.Unit2Price,
					Unit3Price:    product.Unit3Price,
					Ppn:           product.Ppn,
				})
			}

			orderResponses = append(orderResponses, entity.OrderResponse{
				OrderNo:     order.OrderNo,
				InvoiceNo:   order.InvoiceNo,
				OutletName:  order.OutletName,
				Salesman:    order.Salesman,
				InvoiceDate: order.InvoiceDate,
				DueDate:     order.DueDate,
				TotalPrice:  order.TotalPrice,
				Ppn:         order.Ppn,
				Discount:    order.Discount,
				TotalUnpaid: order.TotalUnpaid,
				TotalPromo:  order.TotalPromo,
				PaymentType: order.PaymentType,
				Products:    productResponses,
			})
		}

		responses = append(responses, entity.PicklistResponse{
			PicklistNo: picklist.PicklistNo,
			Driver:     picklist.Driver,
			Helper:     picklist.Helper,
			Vehicle:    picklist.Vehicle,
			UpdatedBy:  picklist.UpdatedBy,
			CreatedAt:  picklist.CreatedAt,
			UpdatedAt:  picklist.UpdatedAt,
			Orders:     orderResponses,
		})
	}

	meta := entity.Meta{
		Page:      dataFilter.Page,
		Limit:     dataFilter.Limit,
		TotalData: int(totalData),
		TotalPage: int(totalPage),
	}

	return responses, meta, nil
}

func (service *PicklistServiceImpl) GetListInvoice(ctx context.Context, dataFilter entity.GeneralQueryFilter, headers map[string]string, customerId string) ([]entity.CustomShipmentInvoice, entity.Meta, error) {
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 999
	}

	endpointUrl := buildInvoiceURL(config.KongUrlSales, dataFilter)
	log.Printf("[GetListInvoice] endpointUrl: %s", endpointUrl)

	req, err := http.NewRequestWithContext(ctx, "GET", endpointUrl, nil)
	if err != nil {
		log.Printf("[GetListInvoice] error build request: %v", err)
		return nil, entity.Meta{}, err
	}

	setRequestHeaders(req, headers)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[GetListInvoice] error request ke %s: %v", endpointUrl, err)
		return nil, entity.Meta{}, err
	}
	defer resp.Body.Close()

	log.Printf("[GetListInvoice] status response: %d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[GetListInvoice] gagal baca body: %v", err)
		return nil, entity.Meta{}, err
	}

	log.Printf("[GetListInvoice] raw response body: %s", string(body))

	var response struct {
		Message   string                           `json:"message"`
		Data      []entity.ResponseShipmentInvoive `json:"data"`
		Paging    pagingMeta                       `json:"paging"`
		RequestID string                           `json:"request_id"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("[GetListInvoice] error unmarshall response: %v", err)
		return nil, entity.Meta{}, err
	}

	assignedInvoiceNoSet := buildAssignedInvoiceSet(service.picklistRepo.FindAllByInvoiceNo(ctx, customerId))
	log.Printf("[GetListInvoice] assignedInvoiceNoSet: %v", assignedInvoiceNoSet)

	filteredInvoice := filterUnassignedInvoices(response.Data, assignedInvoiceNoSet)

	pagination := calculatePagination(response.Paging, len(filteredInvoice))

	return filteredInvoice, *pagination, nil
}

// ==== Helper Struct ====

type pagingMeta struct {
	TotalRecord int `json:"total_record"`
	PageCurrent int `json:"page_current"`
	PageLimit   int `json:"page_limit"`
	PageTotal   int `json:"page_total"`
}

// ==== Helper Functions ====

func buildInvoiceURL(base string, filter entity.GeneralQueryFilter) string {
	var query strings.Builder

	query.WriteString(fmt.Sprintf("%s/v1/invoices?limit=%d", base, filter.Limit))

	if filter.OutletId != "" {
		query.WriteString(fmt.Sprintf("&outlet_id=%s", filter.OutletId))
	}

	if filter.EmpId != "" {
		query.WriteString(fmt.Sprintf("&salesman_id=%s", filter.EmpId))
	}

	switch filter.DocumentType {
	case "invoice":
		query.WriteString("&status=6")
	case "sales_order":
		query.WriteString("&status=2")
	default:
		query.WriteString("&status=2&status=6")
	}

	return query.String()
}

func setRequestHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func buildAssignedInvoiceSet(invoiceNos []string) map[string]struct{} {
	set := make(map[string]struct{}, len(invoiceNos))
	for _, no := range invoiceNos {
		set[no] = struct{}{}
	}
	return set
}

func filterUnassignedInvoices(invoices []entity.ResponseShipmentInvoive, assigned map[string]struct{}) []entity.CustomShipmentInvoice {
	var result []entity.CustomShipmentInvoice

	for _, inv := range invoices {
		if _, exists := assigned[inv.OrderNo]; exists {
			continue
		}
		var custom entity.CustomShipmentInvoice
		helper.Automapper(inv, &custom)

		custom.TotalPromo = inv.PromoValue + inv.PromoBgValue
		custom.TotalBruto = inv.SubTotal + inv.PromoBgValue
		custom.TotalPPN = inv.VatValue
		custom.TotalNetto = inv.Total

		custom.Details = make([]entity.ResponseShipmentInvoiveDetails, len(inv.Details))
		copy(custom.Details, inv.Details)

		result = append(result, custom)
	}
	return result
}

func calculatePagination(p pagingMeta, totalRecords int) *entity.Meta {
	totalPages := 0
	if p.PageLimit > 0 {
		totalPages = (totalRecords + p.PageLimit - 1) / p.PageLimit
	}
	return &entity.Meta{
		TotalData: totalRecords,
		Page:      p.PageCurrent,
		Limit:     p.PageLimit,
		TotalPage: totalPages,
	}
}
