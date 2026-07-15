package service

import (
	"context"
	"log"
	"scyllax-tms/entity"
	"scyllax-tms/exception"
	"scyllax-tms/helper"
	"scyllax-tms/model"
	"scyllax-tms/repository"
	"strings"

	"github.com/go-playground/validator/v10"
)

type VisitService interface {
	GetSummary(ctx context.Context, params entity.SummaryParams) *entity.SummaryResponse
	GetSummaryDailyByParams(ctx context.Context, params entity.SummaryDailyParams) (response entity.SummaryDailyResponse)
	GetDailyActivity(ctx context.Context, params entity.DailyActivityParams) (response []entity.DailyActivityResponse)
	Start(ctx context.Context, request entity.VisitRequest)
	End(ctx context.Context, request entity.VisitRequest)
	Unload(ctx context.Context, headers map[string]string, request entity.UnloadRequest)
	Leave(ctx context.Context, request entity.LeaveRequest)
	Arrive(ctx context.Context, request entity.ArriveRequest)
	Skip(ctx context.Context, request entity.SkipRequest)
}

type VisitServiceImpl struct {
	shipmentRepo         repository.ShipmentRepo
	shipmentInvoicesRepo repository.ShipmentInvoicesRepo
	shipmentService      ShipmentService
	validate             *validator.Validate
}

func NewVisitServiceImpl(shipmentRepo repository.ShipmentRepo, shipmentInvoicesRepo repository.ShipmentInvoicesRepo, shipmentService ShipmentService, validate *validator.Validate) VisitService {
	return &VisitServiceImpl{
		shipmentRepo:         shipmentRepo,
		shipmentInvoicesRepo: shipmentInvoicesRepo,
		shipmentService:      shipmentService,
		validate:             validate,
	}
}

func (service *VisitServiceImpl) GetSummary(ctx context.Context, params entity.SummaryParams) *entity.SummaryResponse {
	shipment, trip, finished, inProgress, err := service.shipmentRepo.CountSummary(ctx, params.DriverId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	recordStart, err := service.shipmentRepo.FindByStartTimes(ctx, []string{"driver_id", "cust_id", "delivery_date"}, []any{params.DriverId, params.CustId, "CURRENT_DATE"})
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	recordFinish, err := service.shipmentRepo.FindByEndTimes(ctx, []string{"driver_id", "cust_id", "delivery_date"}, []any{params.DriverId, params.CustId, "CURRENT_DATE"})
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	response := &entity.SummaryResponse{
		Trip:       trip,
		Finished:   finished,
		InProgress: inProgress,
		Shipment:   shipment,
	}

	var minStart *int64
	var maxFinish *int64

	for _, val := range recordStart {
		if val.Start != nil {
			if minStart == nil || *val.Start < *minStart {
				minStart = val.Start
			}
		}
	}

	for _, val := range recordFinish {
		if val.Finish != nil {
			if maxFinish == nil || *val.Finish > *maxFinish {
				maxFinish = val.Finish
			}
		}
	}

	if minStart != nil {
		response.StartTime = *minStart
	}

	if maxFinish != nil {
		response.EndTime = *maxFinish
	}

	return response
}

func (service *VisitServiceImpl) GetSummaryDailyByParams(ctx context.Context, params entity.SummaryDailyParams) (response entity.SummaryDailyResponse) {
	data, err := service.shipmentRepo.FindByColumns(ctx, []string{"start, finish"}, []string{"shipment_no", "cust_id"}, []any{params.ShipmentNo, params.CustId})

	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	if data.Start != nil {
		response.StartTime = *data.Start
	}

	if data.Finish != nil {
		response.EndTime = *data.Finish
	}

	return response
}

func (service *VisitServiceImpl) GetDailyActivity(ctx context.Context, params entity.DailyActivityParams) (response []entity.DailyActivityResponse) {
	const (
		StatusInProgress = "In Progress"
		StatusPlanned    = "Planned"
	)

	data := service.shipmentRepo.FindByDriverID(ctx, params.DriverId)
	foundActive := false

	for _, value := range data {
		var res entity.DailyActivityResponse
		if !foundActive && (value.Status == StatusInProgress || value.Status == StatusPlanned) {
			res.IsActive = true
			foundActive = true
		} else {
			res.IsActive = false
		}
		helper.Automapper(value, &res)
		response = append(response, res)
	}
	return response
}

func (service *VisitServiceImpl) Start(ctx context.Context, request entity.VisitRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.Shipment{
		Start:      &request.CurrentTime,
		Status:     "In Progress",
		ShipmentNo: request.ShipmentNo,
	}

	err = service.shipmentRepo.UpdateByQuery(ctx, "driver_id", request.DriverID, dataset)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
}

func (service *VisitServiceImpl) End(ctx context.Context, request entity.VisitRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.Shipment{
		Finish:     &request.CurrentTime,
		Status:     "Finished",
		ShipmentNo: request.ShipmentNo,
	}

	err = service.shipmentRepo.UpdateByQuery(ctx, "driver_id", request.DriverID, dataset)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
}

func (service *VisitServiceImpl) Unload(ctx context.Context, headers map[string]string, request entity.UnloadRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.ShipmentInvoices{
		UnloadAt:      &request.CurrentTime,
		OutletStatus:  "Receive All",
		ProductStatus: "Receive",
		ShipmentNo:    &request.ShipmentNo,
	}

	err = service.shipmentInvoicesRepo.UpdateByColumnAtUnload(ctx, request.OutletID, dataset)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	// Fetch all order numbers for the given shipment number
	orderNos := service.shipmentInvoicesRepo.FindAllOrderNoByShipmentNo(ctx, request.ShipmentNo, request.OutletID)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	for _, orderNo := range orderNos {
		if strings.HasPrefix(orderNo, "SO") {
			log.Printf("Processing OrderNo with SO prefix: %s", orderNo)
			orderUpdate := entity.UpdateStatusOrder{
				Orders: []entity.OrderItem{
					{
						OrderNo: orderNo,
						Status:  4,
					},
				},
			}
			log.Printf("Order Update: %+v", orderUpdate)
			err = service.shipmentService.UpdateStatusOrder(ctx, headers, orderUpdate)
			if err != nil {
				log.Printf("Error updating order status: %v", err)
				return
			}
		} else if strings.HasPrefix(orderNo, "SR") {
			log.Printf("Processing OrderNo with SR prefix: %s", orderNo)
			returnUpdate := entity.UpdateStatusReturn{
				Returns: []entity.ReturnItem{
					{
						OrderNo: orderNo,
						Status:  4,
					},
				},
			}
			log.Printf("Return Update: %+v", returnUpdate)
			err = service.shipmentService.UpdateStatusReturn(ctx, headers, returnUpdate)
			if err != nil {
				log.Printf("Error updating return status: %v", err)
				return
			}
		} else {
			log.Printf("OrderNo does not match SO or SR prefixes: %s", orderNo)
		}
	}
}

func (service *VisitServiceImpl) Leave(ctx context.Context, request entity.LeaveRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.ShipmentInvoices{
		LeaveAt:      &request.CurrentTime,
		ShipmentNo:   &request.ShipmentNo,
		OutletStatus: "Finished",
	}

	err = service.shipmentInvoicesRepo.UpdateByColumnAt(ctx, request.OutletID, dataset)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
}

func (service *VisitServiceImpl) Arrive(ctx context.Context, request entity.ArriveRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.ShipmentInvoices{
		ArriveAt:     &request.CurrentTime,
		OutletStatus: "On Progress",
		ShipmentNo:   &request.ShipmentNo,
	}

	err = service.shipmentInvoicesRepo.UpdateByColumnAt(ctx, request.OutletID, dataset)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
}

func (service *VisitServiceImpl) Skip(ctx context.Context, request entity.SkipRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dataset := model.ShipmentInvoices{
		SkipAt:       &request.CurrentTime,
		OutletStatus: "Skipped",
		ShipmentNo:   &request.ShipmentNo,
		SkipReason:   request.SkipReason,
		InOutlet:     *request.InOutlet,
	}

	err = service.shipmentInvoicesRepo.UpdateSkip(ctx, request.OutletID, dataset)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
}
