package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"scyllax-tms/config"
	"scyllax-tms/entity"
	"scyllax-tms/exception"
	"scyllax-tms/helper"
	"scyllax-tms/model"
	"scyllax-tms/repository"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type ShipmentService interface {
	CreateManual(ctx context.Context, headers map[string]string, request entity.CreateShipmentRequest) (string, error)
	CreateAuto(ctx context.Context, headers map[string]string, request entity.CreateShipmentAutoRequest) ([]string, error)
	FindAll(ctx context.Context, dataFilter entity.ShipmentQueryFilter) (response []entity.ShipmentResponse)
	FindByShipmentNo(ctx context.Context, params entity.ShipmentParams) (response entity.ShipmentPreviewResponse)
	FindShipmentInvoiceByShipmentNo(ctx context.Context, params entity.ShipmentParams) (response entity.ShipmentPickList)
	SubmitShipment(ctx context.Context, request entity.SubmitShipmentRequest)
	Delete(ctx context.Context, headers map[string]string, params entity.ShipmentParams)
	DeleteBulk(ctx context.Context, request entity.DeleteShipmentRequest)
	LoginSendPick(ctx context.Context) (string, error)
	GenerateSendPick(ctx context.Context, request entity.CreateShipmentAutoRequest) (data []entity.SendPickResponse, err error)
	UpdateStatusOrder(ctx context.Context, header map[string]string, request entity.UpdateStatusOrder) (err error)
	MobileUpdateStatusOrder(ctx context.Context, header map[string]string, request entity.UpdateStatusOrder) (err error)
	UpdateStatusReturn(ctx context.Context, header map[string]string, request entity.UpdateStatusReturn) (err error)
	MobileUpdateStatusReturn(ctx context.Context, header map[string]string, request entity.UpdateStatusReturn) (err error)
}

type ShipmentServiceImpl struct {
	shipmentRepo         repository.ShipmentRepo
	shipmentInvoicesRepo repository.ShipmentInvoicesRepo
	validate             *validator.Validate
}

func NewShipmentServiceImpl(shipmentRepo repository.ShipmentRepo, shipmentInvoicesRepo repository.ShipmentInvoicesRepo, validate *validator.Validate) ShipmentService {
	return &ShipmentServiceImpl{
		shipmentRepo:         shipmentRepo,
		shipmentInvoicesRepo: shipmentInvoicesRepo,
		validate:             validate,
	}
}

func (service *ShipmentServiceImpl) CreateManual(ctx context.Context, headers map[string]string, request entity.CreateShipmentRequest) (string, error) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	dateTime, err := time.Parse("2006-01-02", request.DeliveryDate) // yyyy-mm-dd
	helper.ErrorPanic(err)

	// shipmentNo := helper.ManualGenerateShipmentNo(request.Vehicle.VehicleID)
	shipmentNo := helper.GenerateShipmentNo()

	shipmentType := "manual"

	dataset := model.Shipment{
		ShipmentNo:   shipmentNo,
		DriverID:     request.Vehicle.DriverID,
		DriverName:   request.Vehicle.DriverName,
		HelperID:     request.Vehicle.HelperID,
		HelperName:   request.Vehicle.HelperName,
		VehicleID:    request.Vehicle.VehicleID,
		VehicleNo:    request.Vehicle.VehicleNo,
		VehicleType:  request.Vehicle.VehicleType,
		VehicleName:  request.Vehicle.VehicleName,
		Length:       request.Vehicle.Length,
		Width:        request.Vehicle.Width,
		Height:       request.Vehicle.Height,
		Volume:       request.Vehicle.Volume,
		Weight:       request.Vehicle.Weight,
		DeliveryDate: dateTime,
		CustID:       request.Vehicle.CustID,
		ShipmentType: shipmentType,
	}

	error := service.shipmentRepo.Insert(ctx, dataset)
	if error != nil {
		return "", exception.NewBadRequestError(error.Error())
	}
	shipmentCurrent := dataset.ShipmentNo

	for _, value := range request.Shipment {
		dateTime, err = time.Parse("2006-01-02", request.DeliveryDate) // yyyy-mm-dd
		helper.ErrorPanic(err)

		var dueDate *time.Time

		// Check if DueDate is not empty or nil
		if value.DueDate != "" {
			parsedDate, err := time.Parse("2006-01-02", value.DueDate)
			// log.Printf("Due Date parsed: %s", parsedDate)
			if err != nil {
				log.Printf("Error parsing DueDate for OrderNo %s: %v", value.OrderNo, err)
				continue
			}

			dueDate = &parsedDate
		}

		var date *time.Time

		if value.Date != "" {
			parsedDate, err := time.Parse("2006-01-02", value.Date)
			// log.Printf("Due Date parsed: %s", parsedDate)
			if err != nil {
				log.Printf("Error parsing Date for OrderNo %s: %v", value.OrderNo, err)
				continue
			}

			date = &parsedDate
		}

		// check if unitId2 == unitId3 --> add qty2 to qty3 & set to 0 for qty2
		if value.UnitId2 == value.UnitId3 {
			value.Qty3 += value.Qty2
			value.Qty2 = 0
		}

		dataset := model.ShipmentInvoices{
			ShipmentNo:         &shipmentCurrent,
			OrderNo:            &value.OrderNo,
			InvoiceNo:          &value.InvoiceNo,
			OutletID:           value.OutletID,
			OutletCode:         value.OutletCode,
			OutletAddress:      value.OutletAddress,
			OutletStatus:       value.OutletStatus,
			OutletName:         value.OutletName,
			SalesmanID:         value.SalesmanID,
			SalesmanName:       value.SalesmanName,
			ProductID:          value.ProductID,
			ProductName:        value.ProductName,
			ProductStatus:      value.ProductStatus,
			ProductCode:        value.ProductCode,
			Sku:                value.Sku,
			Qty1:               value.Qty1,
			Qty2:               value.Qty2,
			Qty3:               value.Qty3,
			ConvUnit1:          value.ConvUnit1,
			ConvUnit2:          value.ConvUnit2,
			ConvUnit3:          value.ConvUnit3,
			UnitId1:            value.UnitId1,
			UnitId2:            value.UnitId2,
			UnitId3:            value.UnitId3,
			Volume:             value.Volume,
			Weight:             value.Weight,
			Status:             value.Status,
			CustID:             value.CustID,
			DeliveryDate:       dateTime,
			WarehouseLatitude:  fmt.Sprint(value.WarehouseLatitude),
			WarehouseLongitude: fmt.Sprint(value.WarehouseLongitude),
			OutletLatitude:     fmt.Sprint(value.OutletLatitude),
			OutletLongitude:    fmt.Sprint(value.OutletLongitude),
			TotalBruto:         value.TotalBruto,
			TotalVolumeDisc:    value.TotalVolumeDisc,
			TotalPromo:         value.TotalPromo,
			TotalPpn:           value.TotalPpn,
			DueDate:            dueDate,
			SellPrice1:         value.SellPrice1,
			SellPrice2:         value.SellPrice2,
			SellPrice3:         value.SellPrice3,
			TotalBelumBayar:    value.TotalNetto,
			PayTypeName:        value.PayTypeName,
			InvoiceDate:        date,
			TotalNetto:         value.TotalNetto,
			Vat:                value.Vat,
			VatValue:           value.VatValue,
			ItemCdnName:        value.ItemCdnName,
			OrderDetailID:      value.OrderDetailID,
		}

		error := service.shipmentInvoicesRepo.Insert(ctx, dataset)
		if error != nil {
			return "error when storing data", exception.NewBadRequestError(error.Error())
		}

		if strings.HasPrefix(value.OrderNo, "SO") {
			log.Printf("Processing OrderNo with SO prefix: %s", value.OrderNo)
			orderUpdate := entity.UpdateStatusOrder{
				Orders: []entity.OrderItem{
					{
						OrderNo: value.OrderNo,
						Status:  3,
					},
				},
			}
			log.Printf("Order Update: %+v", orderUpdate)
			err = service.UpdateStatusOrder(ctx, headers, orderUpdate)
			if err != nil {
				log.Printf("Error updating order status: %v", err)
				return "", err
			}
		} else if strings.HasPrefix(value.OrderNo, "SR") {
			log.Printf("Processing OrderNo with SR prefix: %s", value.OrderNo)
			returnUpdate := entity.UpdateStatusReturn{
				Returns: []entity.ReturnItem{
					{
						OrderNo: value.OrderNo,
						Status:  4,
					},
				},
			}
			log.Printf("Return Update: %+v", returnUpdate)
			err = service.UpdateStatusReturn(ctx, headers, returnUpdate)
			if err != nil {
				log.Printf("Error updating return status: %v", err)
				return "", err
			}
		} else {
			log.Printf("OrderNo does not match SO or SR prefixes: %s", value.OrderNo)
		}
	}

	return shipmentNo, nil
}

func (service *ShipmentServiceImpl) CreateAuto(ctx context.Context, headers map[string]string, request entity.CreateShipmentAutoRequest) ([]string, error) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	// Begin a transaction through the repository
	tx, err := service.shipmentRepo.BeginTx(ctx)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	dataSendPick, err := service.GenerateSendPick(ctx, request)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	dateTime, err := time.Parse("2006-01-02", request.DeliveryDate) // yyyy-mm-dd
	if err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	var mapper entity.MapperSendPick
	shipmentMapping := make(map[string][]string)

	for _, response := range dataSendPick {
		parts := strings.Split(response.ShipmentNo, "_")
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid shipment_no format: %s", response.ShipmentNo)
		}

		// vehicleIdStr := parts[1]
		// lastNumber := parts[2]

		vehicleId, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid vehicle_id format: %s", parts[1])
		}

		// Mapping generate shipment number
		// shipmentNo := helper.MappingAutoShipmentNo(vehicleIdStr, lastNumber)
		shipmentNo := helper.GenerateShipmentNo()

		// Store in mapping
		shipmentMapping[shipmentNo] = response.Orders
		mapper.ShipmentNumbers = append(mapper.ShipmentNumbers, shipmentNo)
		mapper.VehicleIDs = append(mapper.VehicleIDs, vehicleId)
	}

	if len(mapper.ShipmentNumbers) == 0 {
		return nil, fmt.Errorf("no shipmentNos available from dataSendPick")
	}

	countVehicle := make(map[int]int)
	for _, id := range mapper.VehicleIDs {
		countVehicle[id]++
	}

	var extendsVehicle []entity.VehicleBody
	for _, data := range request.Vehicle {
		count := countVehicle[data.VehicleID]
		for i := 0; i < count; i++ {
			extendsVehicle = append(extendsVehicle, data)
		}
	}

	for i, value := range extendsVehicle {
		var shipmentNo string
		if i < len(mapper.ShipmentNumbers) {
			shipmentNo = mapper.ShipmentNumbers[i]
		}

		shipmentType := "auto"

		dataset := model.Shipment{
			ShipmentNo:   shipmentNo,
			DriverID:     value.DriverID,
			DriverName:   value.DriverName,
			HelperID:     value.HelperID,
			HelperName:   value.HelperName,
			VehicleID:    value.VehicleID,
			VehicleNo:    value.VehicleNo,
			VehicleType:  value.VehicleType,
			VehicleName:  value.VehicleName,
			Length:       value.Length,
			Width:        value.Width,
			Height:       value.Height,
			Volume:       value.Volume,
			Weight:       value.Weight,
			DeliveryDate: dateTime,
			CustID:       value.CustID,
			ShipmentType: shipmentType,
		}

		//fmt.Println("shipment", dataset)
		err := service.shipmentRepo.InsertWithTx(tx, dataset)
		if err != nil {
			return nil, exception.NewBadRequestError(err.Error())
		}

		if i < len(mapper.ShipmentNumbers) {
			mapper.ShipmentNumbersMapper = append(mapper.ShipmentNumbersMapper, shipmentNo)
		}
	}

	for _, value := range request.Shipment {
		dateTime, err := time.Parse("2006-01-02", request.DeliveryDate) // yyyy-mm-dd
		if err != nil {
			panic(exception.NewBadRequestError(err.Error()))
		}

		var shipmentNo string
		for no, orders := range shipmentMapping {
			for _, order := range orders {
				if order == value.OrderNo {
					shipmentNo = no
					break
				}
			}
			if shipmentNo != "" {
				break
			}
		}

		if shipmentNo == "" {
			continue
		}

		var dueDate *time.Time

		// Check if DueDate is not empty or nil
		if value.DueDate != "" {
			parsedDate, err := time.Parse("2006-01-02", value.DueDate)
			// log.Printf("Due Date parsed: %s", parsedDate)
			if err != nil {
				log.Printf("Error parsing DueDate for OrderNo %s: %v", value.OrderNo, err)
				continue
			}

			dueDate = &parsedDate
		}

		var date *time.Time

		if value.Date != "" {
			parsedDate, err := time.Parse("2006-01-02", value.Date)
			// log.Printf("Due Date parsed: %s", parsedDate)
			if err != nil {
				log.Printf("Error parsing Date for OrderNo %s: %v", value.OrderNo, err)
				continue
			}

			date = &parsedDate
		}

		// check if unitId2 == unitId3 --> add qty2 to qty3 & set to 0 for qty2
		if value.UnitId2 == value.UnitId3 {
			value.Qty3 += value.Qty2
			value.Qty2 = 0
		}

		dataset := model.ShipmentInvoices{
			ShipmentNo:         &shipmentNo,
			OrderNo:            &value.OrderNo,
			InvoiceNo:          &value.InvoiceNo, // Todo
			OutletID:           value.OutletID,
			OutletCode:         value.OutletCode,
			OutletAddress:      value.OutletAddress,
			OutletStatus:       value.OutletStatus,
			OutletName:         value.OutletName,
			SalesmanID:         value.SalesmanID,
			SalesmanName:       value.SalesmanName,
			ProductID:          value.ProductID,
			ProductName:        value.ProductName,
			ProductStatus:      value.ProductStatus,
			ProductCode:        value.ProductCode,
			Sku:                value.Sku,
			Qty1:               value.Qty1,
			Qty2:               value.Qty2,
			Qty3:               value.Qty3,
			ConvUnit1:          value.ConvUnit1,
			ConvUnit2:          value.ConvUnit2,
			ConvUnit3:          value.ConvUnit3,
			UnitId1:            value.UnitId1,
			UnitId2:            value.UnitId2,
			UnitId3:            value.UnitId3,
			Volume:             value.Volume,
			Weight:             value.Weight,
			Status:             value.Status,
			CustID:             value.CustID,
			DeliveryDate:       dateTime,
			WarehouseLatitude:  fmt.Sprint(value.WarehouseLatitude),
			WarehouseLongitude: fmt.Sprint(value.WarehouseLongitude),
			OutletLatitude:     fmt.Sprint(value.OutletLatitude),
			OutletLongitude:    fmt.Sprint(value.OutletLongitude),
			TotalBruto:         value.TotalBruto,
			TotalVolumeDisc:    value.TotalVolumeDisc,
			TotalPromo:         value.TotalPromo,
			TotalPpn:           value.TotalPpn,
			DueDate:            dueDate,
			SellPrice1:         value.SellPrice1,
			SellPrice2:         value.SellPrice2,
			SellPrice3:         value.SellPrice3,
			TotalBelumBayar:    value.TotalBruto,
			PayTypeName:        value.PayTypeName,
			InvoiceDate:        date,
			TotalNetto:         value.TotalNetto,
			Vat:                value.Vat,
			VatValue:           value.VatValue,
			ItemCdnName:        value.ItemCdnName,
			OrderDetailID:      value.OrderDetailID,
		}

		//fmt.Println("shipment invoices", dataset)
		err = service.shipmentInvoicesRepo.InsertWithTx(tx, dataset) // Use transaction
		if err != nil {
			return nil, exception.NewBadRequestError(err.Error())
		}

		if strings.HasPrefix(value.OrderNo, "SO") {
			log.Printf("Processing OrderNo with SO prefix: %s", value.OrderNo)
			orderUpdate := entity.UpdateStatusOrder{
				Orders: []entity.OrderItem{
					{
						OrderNo: value.OrderNo,
						Status:  3,
					},
				},
			}
			log.Printf("Order Update: %+v", orderUpdate)
			err = service.UpdateStatusOrder(ctx, headers, orderUpdate)
			if err != nil {
				log.Printf("Error updating order status: %v", err)
				return nil, err
			}
		} else if strings.HasPrefix(value.OrderNo, "SR") {
			log.Printf("Processing OrderNo with SR prefix: %s", value.OrderNo)
			returnUpdate := entity.UpdateStatusReturn{
				Returns: []entity.ReturnItem{
					{
						OrderNo: value.OrderNo,
						Status:  4,
					},
				},
			}
			log.Printf("Return Update: %+v", returnUpdate)
			err = service.UpdateStatusReturn(ctx, headers, returnUpdate)
			if err != nil {
				log.Printf("Error updating return status: %v", err)
				return nil, err
			}
		} else {
			log.Printf("OrderNo does not match SO or SR prefixes: %s", value.OrderNo)
		}
	}

	return mapper.ShipmentNumbersMapper, nil
}

func (service *ShipmentServiceImpl) FindAll(ctx context.Context, dataFilter entity.ShipmentQueryFilter) (response []entity.ShipmentResponse) {
	result := service.shipmentRepo.FindAll(ctx, dataFilter)

	for _, row := range result {
		var res entity.ShipmentResponse
		helper.Automapper(row, &res)
		response = append(response, res)
	}
	return response
}

func (service *ShipmentServiceImpl) FindByShipmentNo(ctx context.Context, params entity.ShipmentParams) (response entity.ShipmentPreviewResponse) {
	data, err := service.shipmentRepo.FindByShipmentNo(ctx, params.ShipmentNo)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error())) 
	}

	var entries []entity.ShipmentEntryResponse
	outletMap := make(map[int][]entity.ShipmentItemResponse)

	type SoGroup struct {
		TotalBruto      float64
		TotalPromo      float64
		TotalVolumeDisc float64
		TotalPPN        float64
		TotalNetto      float64
	}
	soGroup := make(map[string]*SoGroup)
	processedOrders := make(map[string]bool)

	var uniqueSalesman string
	isSameSalesman := true

	updateGroupTotals := func(shipment model.ShipmentInvoices, group *SoGroup) {
		if shipment.TotalBruto != nil {
			group.TotalBruto += *shipment.TotalBruto
		}
		if shipment.TotalPromo != nil {
			group.TotalPromo += *shipment.TotalPromo
		}
		if shipment.TotalVolumeDisc != nil {
			group.TotalVolumeDisc += *shipment.TotalVolumeDisc
		}
		if shipment.TotalPpn != nil {
			group.TotalPPN += *shipment.TotalPpn
		}
		// Calculate Netto for this group: Bruto + PPN - Promo - Volume Discount
		group.TotalNetto = group.TotalBruto + group.TotalPPN - group.TotalPromo - group.TotalVolumeDisc
	}

	for i, shipment := range data.ShipmentInvoices {
		var outlet entity.ShipmentItemResponse
		helper.Automapper(shipment, &outlet)

		outlet.InvoiceNo = ""
		if shipment.InvoiceNo != nil {
			outlet.InvoiceNo = *shipment.InvoiceNo
		}

		outlet.ShipmentStatus = shipment.Status
		outletMap[shipment.OutletID] = append(outletMap[shipment.OutletID], outlet)

		soNumber := ""
		if shipment.OrderNo != nil {
			soNumber = *shipment.OrderNo
		}

		if processedOrders[soNumber] {
			log.Printf("Skipping already processed OrderNo: %s", soNumber)
			continue 
		}

		processedOrders[soNumber] = true
		log.Printf("Processing OrderNo: %s", soNumber)

		if _, exists := soGroup[soNumber]; !exists {
			soGroup[soNumber] = &SoGroup{}
		}

		updateGroupTotals(shipment, soGroup[soNumber])

		if i == 0 {
			uniqueSalesman = shipment.SalesmanName
		} else if shipment.SalesmanName != uniqueSalesman {
			isSameSalesman = false
		}
	}

	var totalBruto, totalPromo, totalVolumeDisc, totalPPN, totalNetto float64
	for soNumber, group := range soGroup {
		log.Printf("OrderNo: %s | Bruto: %f, Promo: %f, VolumeDisc: %f, PPN: %f, Netto: %f", soNumber, group.TotalBruto, group.TotalPromo, group.TotalVolumeDisc, group.TotalPPN, group.TotalNetto)
		totalBruto += group.TotalBruto
		totalPromo += group.TotalPromo
		totalVolumeDisc += group.TotalVolumeDisc
		totalPPN += group.TotalPPN
		totalNetto += group.TotalNetto
	}

	response.TotalBruto = totalBruto
	response.TotalPromo = totalPromo
	response.TotalVolumeDisc = totalVolumeDisc
	response.TotalPpn = totalPPN
	response.TotalNetto = totalNetto

	for outletID, outletItems := range outletMap {
		entry := entity.ShipmentEntryResponse{
			OutletID: outletID,
			Outlets:  outletItems,
		}
		entries = append(entries, entry)
	}

	if strings.HasPrefix(data.ShipmentNo, "DO") {
		response.PickListNo = strings.Replace(data.ShipmentNo, "DO", "PL", 1)
	} else {
		response.PickListNo = data.ShipmentNo
	}

	response.PickListDate = data.CreatedAt

	// If the salesman is unique, use the unique name; otherwise, set to "Semua"
	if isSameSalesman {
		response.SalesmanName = uniqueSalesman
	} else {
		response.SalesmanName = "Semua"
	}

	helper.Automapper(data, &response)
	response.Shipment = entries

	return response
}


func (service *ShipmentServiceImpl) FindShipmentInvoiceByShipmentNo(ctx context.Context, params entity.ShipmentParams) (response entity.ShipmentPickList) {
	data, err := service.shipmentInvoicesRepo.FindByShipmentNo(ctx, params.ShipmentNo)
	log.Printf("Retreived: %+v", data)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	var invoiceList []entity.ShipmentInvoiceList
	var totalInvoice float64
	// var totalBelumBayar float64
	var totalNetto float64

	for _, shipmentInvoice := range data {
		var invoice entity.ShipmentInvoiceList

		helper.Automapper(shipmentInvoice, &invoice)

		if shipmentInvoice.TotalBruto != nil && shipmentInvoice.TotalPpn != nil {
			invoice.NilaiInvoice = *shipmentInvoice.TotalBruto + *shipmentInvoice.TotalPpn
		} else if shipmentInvoice.TotalBruto != nil {
			invoice.NilaiInvoice = *shipmentInvoice.TotalBruto
		} else {
			invoice.NilaiInvoice = 0
		}

		// if shipmentInvoice.TotalBruto != nil {
		// 	// invoice.TotalBelumBayar = *shipmentInvoice.TotalBruto
		// 	totalInvoice = *shipmentInvoice.TotalBruto
		// 	totalBelumBayar = *shipmentInvoice.TotalBruto
		// }

		if shipmentInvoice.TotalNetto != nil {
			totalNetto += *shipmentInvoice.TotalNetto
		}

		totalInvoice += invoice.NilaiInvoice
		// totalBelumBayar += invoice.NilaiInvoice

		invoiceList = append(invoiceList, invoice)
		// invoice.NilaiInvoice = *shipmentInvoice.TotalNetto + *shipmentInvoice.TotalBruto
		// log.Printf("Data mapping: %+v", invoiceList)

		// totalInvoice += *shipmentInvoice.TotalBruto
		// totalBelumBayar += *shipmentInvoice.TotalBruto
		// totalNetto += *shipmentInvoice.TotalNetto
	}

	response.TotalInvoice = totalNetto
	response.TotalBelumBayar = totalNetto
	response.TotalNetto = totalNetto
	response.SipmentInvoice = invoiceList

	return response
}

func (service *ShipmentServiceImpl) SubmitShipment(ctx context.Context, request entity.SubmitShipmentRequest) {

	for i, outletName := range request.OutletName {
		routeID := request.RouteID[i]

		data, err := service.shipmentInvoicesRepo.FindByTwoColumns(ctx, "shipment_no", "outlet_name", request.ShipmentNo, outletName)
		if err != nil {
			panic(exception.NewNotFoundError(err.Error()))
		}

		for _, value := range data {
			value.RouteID = &routeID
			err := service.shipmentInvoicesRepo.Update(ctx, value)
			if err != nil {
				panic(exception.NewNotFoundError(err.Error()))
			}
		}
	}
}

func (service *ShipmentServiceImpl) Delete(ctx context.Context, headers map[string]string, params entity.ShipmentParams) {
	log.Printf("Shipment no: %s", params.ShipmentNo)
	orderNos := service.shipmentInvoicesRepo.GetAllOrderNoByShipmentNo(ctx, params.ShipmentNo)
	log.Printf("Order nos retrieved: %v", orderNos)

	for _, orderNo := range orderNos {
		if strings.HasPrefix(orderNo, "SO") {
			log.Printf("Processing OrderNo with SO prefix: %s", orderNo)
			orderUpdate := entity.UpdateStatusOrder{
				Orders: []entity.OrderItem{
					{
						OrderNo: orderNo,
						Status:  6,
					},
				},
			}
			log.Printf("Order Update: %+v", orderUpdate)
			err := service.UpdateStatusOrder(ctx, headers, orderUpdate)
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
						Status:  3,
					},
				},
			}
			log.Printf("Return Update: %+v", returnUpdate)
			err := service.UpdateStatusReturn(ctx, headers, returnUpdate)
			if err != nil {
				log.Printf("Error updating return status: %v", err)
				return
			}
		} else if strings.HasPrefix(orderNo, "") {
			log.Printf("Processing OrderNo with SO prefix: nil")
			returnUpdate := entity.UpdateStatusReturn{
				Returns: []entity.ReturnItem{
					{
						OrderNo: orderNo,
						Status:  2,
					},
				},
			}
			log.Printf("Order Update: %+v", returnUpdate)
			err := service.UpdateStatusReturn(ctx, headers, returnUpdate)
			if err != nil {
				log.Printf("Error updating order status: %v", err)
				return
			}
		} else {
			log.Printf("OrderNo does not match SO or SR prefixes: %s", orderNo)
		}
	}

	err := service.shipmentRepo.DeleteByQuery(ctx, "shipment_no", params.ShipmentNo)

	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
}

func (service *ShipmentServiceImpl) DeleteBulk(ctx context.Context, request entity.DeleteShipmentRequest) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	err = service.shipmentRepo.DeleteBulk(ctx, request.ShipmentNo)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}
}

func (service *ShipmentServiceImpl) LoginSendPick(ctx context.Context) (string, error) {
	config, err := config.LoadConfig(".")

	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	data := entity.LoginSendPick{
		Email:    "superadmin@pilarmedia.com",
		Password: "password",
	}

	requestBody, err := json.Marshal(data)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	endpointURL := fmt.Sprintf("%s/auth/login", config.SendPickUrl)
	req, err := http.NewRequestWithContext(ctx, "POST", endpointURL, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	accessToken, ok := response["data"].(map[string]interface{})["access_token"].(string)
	if !ok {
		errors.New("unexpected response format")
	}

	return accessToken, nil
}

func (service *ShipmentServiceImpl) GenerateSendPick(ctx context.Context, request entity.CreateShipmentAutoRequest) (data []entity.SendPickResponse, err error) {
	tokenJwt, err := service.LoginSendPick(ctx)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	orders := []map[string]interface{}{}
	orderMap := make(map[string]bool)

	for _, shipment := range request.Shipment {
		if !orderMap[shipment.OrderNo] {
			order := map[string]interface{}{
				"order_id":            shipment.OrderNo, // change to use order no (SO number)
				"start_time":          "07:00",
				"end_time":            "08:00",
				"demand":              shipment.Volume,
				"commodities":         []string{"general"},
				"stuffing_latitude":   shipment.WarehouseLatitude,
				"stuffing_longitude":  shipment.WarehouseLongitude,
				"stripping_latitude":  shipment.OutletLatitude,
				"stripping_longitude": shipment.OutletLongitude,
			}
			orders = append(orders, order)
			orderMap[shipment.OrderNo] = true
		}
	}

	vehicles := []map[string]interface{}{}
	for _, vehicle := range request.Vehicle {
		veh := map[string]interface{}{
			"vehicle_id":   vehicle.VehicleID,
			"vehicle_type": vehicle.VehicleType,
			"latitude":     vehicle.WarehouseLatitude,     // warehouse latitude
			"longitude":    vehicle.WarehouseLongitude,    // warehouse longitude
			"capacity":     []int{int(vehicle.Volume), 0}, // [volume, hardcoded 0]
			"commodities":  []string{"general"},
		}
		vehicles = append(vehicles, veh)
	}

	payload := map[string]interface{}{
		"stuffing_duration":  0, // hardcoded
		"stripping_duration": 0, // hardcoded
		"shipment_date":      request.DeliveryDate,
		"orders":             orders,
		"vehicles":           vehicles,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	//fmt.Println(string(requestBody))

	endpointURL := fmt.Sprintf("%s/generate/sync", config.SendPickUrl)
	req, err := http.NewRequestWithContext(ctx, "POST", endpointURL, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenJwt))

	type result struct {
		data []entity.SendPickResponse
		err  error
	}

	resultChan := make(chan result)
	go func() {
		defer close(resultChan)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			resultChan <- result{nil, exception.NewBadRequestError(err.Error())}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			resultChan <- result{nil, exception.NewUnauthorizedError("token invalid")}
			return
		}

		if resp.StatusCode == http.StatusInternalServerError {
			resultChan <- result{nil, exception.NewInternalServerError("internal server error")}
			return
		}

		if resp.StatusCode == http.StatusUnprocessableEntity {
			var errorResp map[string][]string
			err = json.NewDecoder(resp.Body).Decode(&errorResp)
			if err != nil {
				resultChan <- result{nil, exception.NewInternalServerError(err.Error())}
				return
			}

			var errorMessages []string
			for _, errors := range errorResp {
				for _, errMsg := range errors {
					errorMessages = append(errorMessages, errMsg)
				}
			}

			resultChan <- result{nil, exception.NewInternalServerError(strings.Join(errorMessages, ", "))}
			return
		}

		var response entity.GenerateResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			resultChan <- result{nil, exception.NewBadRequestError(err.Error())}
			return
		}

		if response.Data.Report.TotalItemUnassigned > 0 && response.Data.Report.TotalItemAssigned == 0 {
			var descriptions []string
			for _, unassigned := range response.Data.Unassigned {
				for _, reason := range unassigned.Reason {
					descriptions = append(descriptions, reason.Description)
				}
			}
			resultChan <- result{nil, exception.NewInternalServerError(strings.Join(descriptions, ", "))}
			return
		}

		shipmentMap := make(map[string][]string)
		for _, res := range response.Data.Result {
			for _, item := range res.ItemDetails {
				shipmentMap[res.Vehicle.ShipmentNo] = append(shipmentMap[res.Vehicle.ShipmentNo], item.OrderID)
			}
		}

		for shipmentNo, orders := range shipmentMap {
			data = append(data, entity.SendPickResponse{
				ShipmentNo: shipmentNo,
				Orders:     orders,
			})
		}
		resultChan <- result{data, nil}
	}()

	select {
	case res := <-resultChan:
		if res.err != nil {
			return nil, res.err
		}
		return res.data, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (service *ShipmentServiceImpl) UpdateStatusOrder(ctx context.Context, headers map[string]string, request entity.UpdateStatusOrder) (err error) {

	config, err := config.LoadConfig(".")

	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	// Remove /sales for staging or production
	endpointUrl := fmt.Sprintf("%s/v1/orders/status", config.KongUrlSales)
	log.Printf("Request URL: %s", endpointUrl)

	// Marshal the request body to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// Create a new PATCH request with the context and request body
	req, err := http.NewRequestWithContext(ctx, "PATCH", endpointUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	// Set content type and add headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Initialize an HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update status: received %d status code", resp.StatusCode)
	}

	log.Println("Status update request successful")
	return nil
}
func (service *ShipmentServiceImpl) MobileUpdateStatusOrder(ctx context.Context, headers map[string]string, request entity.UpdateStatusOrder) (err error) {

	config, err := config.LoadConfig(".")

	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	// Remove /sales for staging or production
	endpointUrl := fmt.Sprintf("%s/v1/orders/status", config.KongUrlMobile)
	log.Printf("Request URL: %s", endpointUrl)

	// Marshal the request body to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// Create a new PATCH request with the context and request body
	req, err := http.NewRequestWithContext(ctx, "PATCH", endpointUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	// Set content type and add headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Initialize an HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update status: received %d status code", resp.StatusCode)
	}

	log.Println("Status update request successful")
	return nil
}

func (service *ShipmentServiceImpl) UpdateStatusReturn(ctx context.Context, headers map[string]string, request entity.UpdateStatusReturn) (err error) {
	config, err := config.LoadConfig(".")

	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	// Remove /sales for staging or production
	endpointUrl := fmt.Sprintf("%s/v1/returns/status", config.KongUrlSales)
	log.Printf("Request URL: %s", endpointUrl)

	// Marshal the request body to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// Create a new POST request with the context and request body
	req, err := http.NewRequestWithContext(ctx, "POST", endpointUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	// Set content type and add headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Initialize an HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update status: received %d status code", resp.StatusCode)
	}

	log.Println("Status update request successful")
	return nil
}
func (service *ShipmentServiceImpl) MobileUpdateStatusReturn(ctx context.Context, headers map[string]string, request entity.UpdateStatusReturn) (err error) {
	config, err := config.LoadConfig(".")

	if err != nil {
		panic(exception.NewInternalServerError(err.Error()))
	}

	// Remove /sales for staging or production
	endpointUrl := fmt.Sprintf("%s/v1/returns/status", config.KongUrlMobile)
	log.Printf("Request URL: %s", endpointUrl)

	// Marshal the request body to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// Create a new POST request with the context and request body
	req, err := http.NewRequestWithContext(ctx, "POST", endpointUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	// Set content type and add headers
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Initialize an HTTP client and make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update status: received %d status code", resp.StatusCode)
	}

	log.Println("Status update request successful")
	return nil
}
