package service

import (
	"context"
	"fmt"
	"log"
	"scyllax-tms/entity"
	"scyllax-tms/exception"
	"scyllax-tms/helper"
	"scyllax-tms/model"
	"scyllax-tms/repository"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type ReportService interface {
	GetDriverReport(ctx context.Context, dataFilter entity.DriverReportQueryFilter) (response entity.DriverReportResponse)
	GetListShipmentNo(ctx context.Context) (response []entity.ShipmentNoDropdown)
	GetListReasons(ctx context.Context) (response []entity.ReasonDropdown)
	GetListOutlet(ctx context.Context) (response []entity.OutletDropdown)
	GetListDriver(ctx context.Context) (response []entity.DriverNameDropdown)
	GetListProductCode(ctx context.Context) (response []entity.ProductCodeDropdown)
	GetShipmentReportSummary(ctx context.Context, dataFilter entity.ShipmentReportQueryFilter) (response []entity.ShipmentReportSummary)
	GetShipmentReportDetail(ctx context.Context, dataFilter entity.ShipmentReportDetailQueryFilter) (response []entity.ShipmentReportDetail)
	GetShipmentReportReject(ctx context.Context, dataFilter entity.ShipmentReportRejectlQueryFilter) (response []entity.ShipmentReportReject)
}

type ReportServiceImpl struct {
	shipmentInvoicesRepo repository.ShipmentInvoicesRepo
	validate             *validator.Validate
}

func NewReportServiceImpl(shipmentInvoicesRepo repository.ShipmentInvoicesRepo, validate *validator.Validate) ReportService {
	return &ReportServiceImpl{
		shipmentInvoicesRepo: shipmentInvoicesRepo,
		validate:             validate,
	}
}

func (service *ReportServiceImpl) GetDriverReport(ctx context.Context, dataFilter entity.DriverReportQueryFilter) (response entity.DriverReportResponse) {
	shipment, finished, skipped, trip, progress, err := service.shipmentInvoicesRepo.CountReport(ctx, dataFilter)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	data, err := service.shipmentInvoicesRepo.GetReport(ctx, dataFilter)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}
	log.Println("Data: ", data)

	response = entity.DriverReportResponse{
		TotalShipment: shipment,
		TotalTrip:     trip,
		TotalFinished: finished,
		TotalSkipped:  skipped,
		Progress:      progress,
	}

	if len(data) > 0 {
		response.SkippedReasons = &data
	} else {
		response.SkippedReasons = nil
	}

	return response
}

func (service *ReportServiceImpl) GetListShipmentNo(ctx context.Context) (response []entity.ShipmentNoDropdown) {
	shipmentNumbers := service.shipmentInvoicesRepo.GetListShipmentNo(ctx)

	return shipmentNumbers
}

func (service *ReportServiceImpl) GetListReasons(ctx context.Context) (response []entity.ReasonDropdown) {
	reasons := service.shipmentInvoicesRepo.GetListReasons(ctx)

	return reasons
}

func (service *ReportServiceImpl) GetListOutlet(ctx context.Context) (response []entity.OutletDropdown) {
	outlets := service.shipmentInvoicesRepo.GetListOutlet(ctx)

	return outlets
}

func (service *ReportServiceImpl) GetListDriver(ctx context.Context) (response []entity.DriverNameDropdown) {
	drivers := service.shipmentInvoicesRepo.GetListDriver(ctx)

	return drivers
}

func (service *ReportServiceImpl) GetListProductCode(ctx context.Context) (response []entity.ProductCodeDropdown) {
	pcodes := service.shipmentInvoicesRepo.GetListProductCode(ctx)

	return pcodes
}

// TODO Shipment Report Service
func (service *ReportServiceImpl) GetShipmentReportSummary(ctx context.Context, dataFilter entity.ShipmentReportQueryFilter) (response []entity.ShipmentReportSummary) {
	shipments, err := service.shipmentInvoicesRepo.GetShipmentReportSummary(ctx, dataFilter)
	// log.Printf("Date filter: %+v", dataFilter)
	helper.ErrorPanic(err)

	summaryMap := make(map[string]*entity.ShipmentReportSummary)

	for _, shipment := range shipments {
		plannedOutlets := make(map[int]bool)
		visitedOutlets := make(map[int]bool)
		skippedOutlets := make(map[int]bool)
		receivedOutlets := make(map[int]bool)
		rejectPartialOutlets := make(map[int]bool)
		rejectAllOutlets := make(map[int]bool)

		for _, invoice := range shipment.ShipmentInvoices {
			outletID := invoice.OutletID

			if _, exists := plannedOutlets[outletID]; !exists {
				plannedOutlets[outletID] = true
			}

			if invoice.ArriveAt != nil {
				visitedOutlets[outletID] = true
			}

			if invoice.SkipAt != nil {
				skippedOutlets[outletID] = true
			}

			if invoice.ProductStatus == "received" {
				receivedOutlets[outletID] = true
			}

			if invoice.ProductStatus == "Reject Partial" {
				rejectPartialOutlets[outletID] = true
			}

			if invoice.ProductStatus == "Reject" {
				rejectAllOutlets[outletID] = true
			}
		}

		if summary, exists := summaryMap[shipment.ShipmentNo]; exists {
			summary.Planned += int64(len(plannedOutlets))
			summary.Visited += int64(len(visitedOutlets))
			summary.Skipped += int64(len(skippedOutlets))
			summary.Received += int64(len(receivedOutlets))
			summary.RejectPartial += int64(len(rejectPartialOutlets))
			summary.RejectAll += int64(len(rejectAllOutlets))
		} else {
			summaryMap[shipment.ShipmentNo] = &entity.ShipmentReportSummary{
				DeliveryDate:  shipment.DeliveryDate,
				ShipmentNo:    shipment.ShipmentNo,
				DriverName:    shipment.DriverName,
				StartTime:     shipment.Start,
				EndTime:       shipment.Finish,
				Planned:       int64(len(plannedOutlets)),
				Visited:       int64(len(visitedOutlets)),
				Skipped:       int64(len(skippedOutlets)),
				Received:      int64(len(receivedOutlets)),
				RejectPartial: int64(len(rejectPartialOutlets)),
				RejectAll:     int64(len(rejectAllOutlets)),
			}
		}
	}

	reportSummaries := []entity.ShipmentReportSummary{}

	for _, summary := range summaryMap {
		reportSummaries = append(reportSummaries, *summary)
	}

	return reportSummaries
}

func (service *ReportServiceImpl) GetShipmentReportDetail(ctx context.Context, dataFilter entity.ShipmentReportDetailQueryFilter) (response []entity.ShipmentReportDetail) {
	shipments, err := service.shipmentInvoicesRepo.GetShipmentReportDetail(ctx, dataFilter)
	helper.ErrorPanic(err)

	var reportDetails []entity.ShipmentReportDetail
	groupedReports := make(map[string]*entity.ShipmentReportDetail)

	for _, shipment := range shipments {
		key := fmt.Sprintf("%s-%s-%s", shipment.DeliveryDate.Format("2006-01-02"), shipment.ShipmentNo, shipment.DriverName)

		if _, exists := groupedReports[key]; !exists {
			groupedReports[key] = &entity.ShipmentReportDetail{
				DeliveryDate:          shipment.DeliveryDate.Format("2006-01-02"),
				ShipmentNo:            shipment.ShipmentNo,
				DriverName:            shipment.DriverName,
				ShipmentReportDetails: []entity.ShipmentReportDetails{},
			}
		}

		orderNoSet := make(map[string]struct{})

		// Track the previous LeaveAt time for calculating DriveTime and ETA
		var previousLeaveAt *int64

		// Iterate through each invoice and set the statuses based on the `OrderNo`
		for _, invoice := range shipment.ShipmentInvoices {
			var documentNo string

			// if invoice.OrderNo != nil {
			// 	documentNo = *invoice.OrderNo
			// }
			if strings.HasPrefix(*invoice.OrderNo, "SO") {
				if invoice.InvoiceNo != nil {
					documentNo = *invoice.InvoiceNo
				}
			} else if strings.HasPrefix(*invoice.OrderNo, "SR") {
				if invoice.OrderNo != nil {
					documentNo = *invoice.OrderNo
				}
			}

			// Check if the OrderNo is unique for this shipment
			if _, alreadyExists := orderNoSet[*invoice.OrderNo]; alreadyExists {
				continue
			}
			orderNoSet[*invoice.OrderNo] = struct{}{}

			// Prepare time-related data
			var startTime, endTime, spent, driveTime, ETA int
			if shipment.Start != nil {
				startTime = int(*shipment.Start)
			}
			if shipment.Finish != nil {
				endTime = int(*shipment.Finish)
			}
			if invoice.ArriveAt != nil && invoice.LeaveAt != nil {
				spent = int((*invoice.LeaveAt - *invoice.ArriveAt) / 60)
			}
			if previousLeaveAt != nil && invoice.ArriveAt != nil {
				driveTime = int(*invoice.ArriveAt-*previousLeaveAt) / 60
			}
			if previousLeaveAt != nil {
				etaTime := time.Unix(*previousLeaveAt, 0).Add(time.Duration(driveTime) * time.Minute)
				ETA = int(etaTime.Unix())
			}

			// Calculate drive time and ETA if outlet_id is different
			// if previousOutletID != invoice.OutletID && previousLeaveAt != nil && invoice.ArriveAt != nil {
			// 	driveTime = int(*invoice.ArriveAt-*previousLeaveAt) / 60
			// 	etaTime := time.Unix(*previousLeaveAt, 0).Add(time.Duration(driveTime) * time.Minute)
			// 	ETA = int(etaTime.Unix())
			// }

			// previousOutletID = invoice.OutletID
			// previousLeaveAt = invoice.LeaveAt

			previousLeaveAt = invoice.LeaveAt

			var receivedStatus string
			switch {
			case strings.HasPrefix(*invoice.OrderNo, "SO"):
				receivedStatus = determineSOReceivedStatus(shipment.ShipmentInvoices)
			case strings.HasPrefix(*invoice.OrderNo, "SR"):
				if invoice.ProductStatus == "Pick Up" {
					receivedStatus = "Picked Up"
				} else if invoice.OutletStatus == "Skipped" || invoice.ProductStatus == "Skip" {
					receivedStatus = "Skipped"
				}
			default:
				receivedStatus = invoice.ProductStatus
			}

			detail := entity.ShipmentReportDetails{
				DocumentNo:     documentNo,
				OutletCode:     invoice.OutletCode,
				OutletName:     invoice.OutletName,
				VisitedStatus:  shipment.Status,
				StartTime:      startTime,
				ArriveAt:       invoice.ArriveAt,
				LeaveAt:        invoice.LeaveAt,
				EndTime:        endTime,
				DriveTime:      driveTime,
				UnloadAt:       invoice.UnloadAt,
				Spent:          spent,
				ETA:            ETA,
				ReceivedStatus: receivedStatus,
				Photo:          invoice.Photo,
				Signature:      invoice.Signature,
			}

			groupedReports[key].ShipmentReportDetails = append(groupedReports[key].ShipmentReportDetails, detail)
		}
	}

	for _, report := range groupedReports {
		reportDetails = append(reportDetails, *report)
	}

	return reportDetails
}

func determineSOReceivedStatus(invoices []model.ShipmentInvoices) string {
	var hasRejectAll, hasRejectPartial, outletSkipped, allReceived bool
	allReceived = true

	for _, invoice := range invoices {
		switch invoice.ProductStatus {
		case "Reject":
			hasRejectAll = true
			allReceived = false
		case "Reject Partial":
			hasRejectPartial = true
			allReceived = false
		case "Receive":
			continue
		default:
			outletSkipped = invoice.OutletStatus == "Skipped"
			allReceived = false
		}
	}

	if outletSkipped {
		return "Skipped"
	} else if hasRejectAll {
		return "Reject All"
	} else if hasRejectPartial {
		return "Reject Partial"
	} else if allReceived {
		return "Received All"
	}
	return "-"
}

func (service *ReportServiceImpl) GetShipmentReportReject(ctx context.Context, dataFilter entity.ShipmentReportRejectlQueryFilter) (response []entity.ShipmentReportReject) {
	shipments, err := service.shipmentInvoicesRepo.GetShipmentReportReject(ctx, dataFilter)
	helper.ErrorPanic(err)

	groupedReports := make(map[string]*entity.ShipmentReportReject)

	for _, shipment := range shipments {
		key := fmt.Sprintf("%s-%s-%s", shipment.ShipmentNo, shipment.DeliveryDate.Format("2006-01-02"), shipment.DriverName)

		if report, exists := groupedReports[key]; exists {
			for _, invoice := range shipment.ShipmentInvoices {
				if invoice.ProductStatus == "Reject" || invoice.ProductStatus == "Reject Partial" {
					detail := entity.ShipmentReportRejectDetails{
						InvoiceNo:   *invoice.InvoiceNo,
						OutletCode:  invoice.OutletCode,
						OutletName:  invoice.OutletName,
						ProductName: invoice.ProductName,
						ProductCode: invoice.ProductCode,
						QtyReject1:  *invoice.QtyReject1,
						QtyReject2:  *invoice.QtyReject2,
						QtyReject3:  *invoice.QtyReject3,
						ReasonName:  *invoice.ReasonName,
					}
					report.ShipmentReportRejectDetails = append(report.ShipmentReportRejectDetails, detail)
				}
			}
		} else {
			detailsList := []entity.ShipmentReportRejectDetails{}
			for _, invoice := range shipment.ShipmentInvoices {
				if invoice.ProductStatus == "Reject" || invoice.ProductStatus == "Reject Partial" {
					detail := entity.ShipmentReportRejectDetails{
						InvoiceNo:   *invoice.InvoiceNo,
						OutletCode:  invoice.OutletCode,
						OutletName:  invoice.OutletName,
						ProductName: invoice.ProductName,
						ProductCode: invoice.ProductCode,
						QtyReject1:  *invoice.QtyReject1,
						QtyReject2:  *invoice.QtyReject2,
						QtyReject3:  *invoice.QtyReject3,
						ReasonName:  *invoice.ReasonName,
					}
					detailsList = append(detailsList, detail)
				}
			}

			groupedReports[key] = &entity.ShipmentReportReject{
				DeliveryDate:                shipment.DeliveryDate.Format("2006-01-02"),
				ShipmentNo:                  shipment.ShipmentNo,
				DriverName:                  shipment.DriverName,
				ShipmentReportRejectDetails: detailsList,
			}
		}
	}

	for _, report := range groupedReports {
		response = append(response, *report)
	}

	return response
}
