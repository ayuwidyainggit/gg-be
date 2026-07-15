package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"inventory/adapter"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg"
	"inventory/pkg/constant"
	"inventory/pkg/conversion"
	"inventory/pkg/errmsg"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
	"log"
	"math"
	"mime/multipart"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockOpnameService interface {
	Store(request entity.CreateStockOpname) (err error)
	StoreV2(request entity.CreateStockOpnameV2) (err error)
	Reports(params entity.ReportStockOpanmeParams) (response entity.StockOpnameReports, err error)
	Cancel(params entity.CancelStockOpanmeParams) (err error)
	List(dataFilter entity.StockOpnameQueryFilter) (data []entity.StockOpnameList, total int64, lastPage int, err error)
	ListV2(dataFilter entity.StockOpnameListV2QueryFilter) (data []entity.StockOpnameListV2Response, total int64, lastPage int, err error)
	DetailV2(params entity.StockOpnameDetailV2Params) (data entity.StockOpnameDetailV2Response, err error)
	UpdateStatusV2(params entity.UpdateStockOpnameStatusV2Params, request entity.UpdateStockOpnameStatusV2Request) error
	RevisedV2(params entity.RevisedStockOpnameV2Params, request entity.RevisedStockOpnameV2Request) error
	StartV2(params entity.StartStockOpnameV2Params) (entity.StartStockOpnameV2Response, error)
	SubmitV2(params entity.SubmitStockOpnameV2Params, request entity.SubmitStockOpnameV2Request, fileHeader *multipart.FileHeader) (entity.SubmitStockOpnameV2Response, error)
	ProductList(dataFilter entity.StockOpnameProductListQueryFilter) (data []entity.StockOpnameProductListResponse, total int64, lastPage int, err error)
	DownloadTemplate(params entity.StockOpnameTemplateDownloadParams) (response entity.StockOpnameTemplateDownloadResponse, err error)
	BulkUpload(params entity.BulkUploadStockOpnameV2Params, file []byte, filename string) (entity.BulkUploadStockOpnameV2Response, error)
	DownloadReport(params entity.StockOpnameDownloadParams) (response entity.StockOpnameDownloadResponse, err error)
}

func NewStockOpnameService(stockOpname repository.StockOpnameRepository, transaction repository.Dbtransaction, stockRepository repository.StockRepository, warehouseStockRepository repository.WarehouseStockRepository, obsAdapter adapter.ObsAdapter) *StockOpnameServiceImpl {
	return &StockOpnameServiceImpl{
		StockOpnameRepository:    stockOpname,
		Transaction:              transaction,
		StockRepository:          stockRepository,
		WarehouseStockRepository: warehouseStockRepository,
		ObsAdapter:               obsAdapter,
	}
}

type StockOpnameServiceImpl struct {
	StockOpnameRepository    repository.StockOpnameRepository
	Transaction              repository.Dbtransaction
	StockRepository          repository.StockRepository
	WarehouseStockRepository repository.WarehouseStockRepository
	ObsAdapter               adapter.ObsAdapter
}

func (service *StockOpnameServiceImpl) Store(request entity.CreateStockOpname) (err error) {
	c := context.Background()

	scheduledAt, err := str.DateTimeStrToRfc3339StringInAsiaJkt(request.ScheduledAt)
	if err != nil {
		err := fmt.Errorf(errmsg.ERROR_DATE_FORMAT, "scheduled_at", request.ScheduledAt)
		return err
	}
	request.ScheduledAt = scheduledAt

	// validate scheduled_at must be greater than now
	tzJakarta, _ := time.LoadLocation("Asia/Jakarta")
	scheduledAtTime, _ := time.Parse(time.RFC3339, request.ScheduledAt)
	if scheduledAtTime.Before(time.Now().In(tzJakarta)) {
		err := fmt.Errorf(errmsg.ERROR_DATE_MUST_GREATER_THAN_NOW, "scheduled_at", request.ScheduledAt)
		return err
	}

	// validate warehouse id
	_, err = service.StockOpnameRepository.FindWarehouseByIDAndCustID(request.WhID, request.CustID)
	if err != nil {
		err := fmt.Errorf("wh_id: %d, %s", request.WhID, err.Error())
		return err
	}

	// validate employee id
	_, err = service.StockOpnameRepository.FindEmployeeByIDAndCustID(request.AssignToEmpID, request.CustID)
	if err != nil {
		err := fmt.Errorf("assign_to_emp_id: %d, %s", request.AssignToEmpID, err.Error())
		return err
	}

	request.DataStatus = 1 // make it default status 'Scheduled'
	var stockOpnameModel model.StockOpname
	err = structs.Automapper(request, &stockOpnameModel)
	if err != nil {
		return err
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.StockOpnameRepository.Store(txCtx, &stockOpnameModel)
		if err != nil {
			log.Println("Stock Opname Store, err:", err.Error())
			return err
		}

		docNo := stockOpnameModel.DocNo
		objectID := primitive.NewObjectID() // Generate a new ObjectID
		objectIDString := objectID.Hex()    // Convert ObjectID to string
		stockOpnameReport := model.StockOpnameReport{
			CustID:        request.CustID,
			DocNo:         docNo,
			StockReportID: objectIDString,
			Status:        1,
		}
		err = service.StockOpnameRepository.StoreReport(txCtx, &stockOpnameReport)
		if err != nil {
			log.Println("StoreReport, err:", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *StockOpnameServiceImpl) List(dataFilter entity.StockOpnameQueryFilter) (data []entity.StockOpnameList, total int64, lastPage int, err error) {
	stockOpname, total, lastPage, err := service.StockOpnameRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range stockOpname {
		var vResp entity.StockOpnameList
		structs.Automapper(row, &vResp)
		vResp.ScheduledAt = row.ScheduledAt.Format("2006-01-02 15:04:05")
		vResp.StatusDescription = vResp.GetStockOpnameStatusDesc()
		vResp.ProductHierarchyDesc = vResp.GetProductHierarchyDesc()
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *StockOpnameServiceImpl) ListV2(dataFilter entity.StockOpnameListV2QueryFilter) (data []entity.StockOpnameListV2Response, total int64, lastPage int, err error) {
	stockOpnameList, total, lastPage, err := service.StockOpnameRepository.FindAllByCustIdV2(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range stockOpnameList {
		vResp := entity.StockOpnameListV2Response{
			DocNo:         row.DocNo,
			CreatedDate:   row.CreatedDate,
			WhID:          row.WhID,
			WhCode:        row.WhCode,
			WhName:        row.WhName,
			CreatedBy:     row.CreatedBy,
			User:          row.UserName,
			ScheduledDate: row.ScheduledDate,
			EmpID:         row.EmpID,
			EmpName:       row.EmpName,
			Status:        row.Status,
		}
		vResp.StatusDesc = vResp.GetStatusDesc()
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *StockOpnameServiceImpl) Reports(params entity.ReportStockOpanmeParams) (response entity.StockOpnameReports, err error) {

	so, err := service.StockOpnameRepository.FindByNo(params)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(so, &response)
	if err != nil {
		return response, err
	}

	var responseList entity.StockOpnameList
	err = structs.Automapper(so, &responseList)
	if err != nil {
		return response, err
	}

	response.ScheduledAt = so.ScheduledAt.Format("2006-01-02 15:04:05")
	response.StatusDescription = responseList.GetStockOpnameStatusDesc()
	response.ProductHierarchyDesc = responseList.GetProductHierarchyDesc()

	soReports, err := service.StockOpnameRepository.FindAllStockOpnameReportByNo(params)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(soReports, &response.StockOpnameReports)
	if err != nil {
		return response, err
	}

	for i, row := range response.StockOpnameReports {
		response.StockOpnameReports[i].StatusDescription = row.GetOpnameReportStatusDesc()
	}

	return response, nil
}

func (service *StockOpnameServiceImpl) Cancel(params entity.CancelStockOpanmeParams) (err error) {
	c := context.Background()

	var findParams entity.ReportStockOpanmeParams
	err = structs.Automapper(params, &findParams)
	if err != nil {
		return err
	}
	so, err := service.StockOpnameRepository.FindByNo(findParams)
	if err != nil {
		return err
	}

	var vResp entity.StockOpnameList
	structs.Automapper(so, &vResp)

	if so.DataStatus >= 50 {
		return errors.New("already " + vResp.GetStockOpnameStatusDesc())
	}

	var updateParams entity.UpdateStockOpnameParams
	err = structs.Automapper(params, &updateParams)
	if err != nil {
		return err
	}

	stockOpname := model.StockOpname{
		UpdatedAt:  time.Now(),
		UpdatedBy:  &params.CancelBy,
		DataStatus: 50,
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.StockOpnameRepository.Update(txCtx, updateParams, stockOpname)
		if err != nil {
			log.Println("Stock Opname Update, err:", err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *StockOpnameServiceImpl) ProductList(dataFilter entity.StockOpnameProductListQueryFilter) (data []entity.StockOpnameProductListResponse, total int64, lastPage int, err error) {
	products, total, lastPage, err := service.StockOpnameRepository.ProductList(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, product := range products {
		var vResp entity.StockOpnameProductListResponse
		structs.Automapper(product, &vResp)

		// Calculate qty1, qty2, qty3 using conversion
		qty := &conversion.Qty{
			Qty:       int(product.Qty),
			ConvUnit2: int(product.ConvUnit2),
			ConvUnit3: int(product.ConvUnit3),
		}
		qtyConversion := qty.ConvToQtyConversion()
		vResp.Qty1 = float64(qtyConversion.Qty1)
		vResp.Qty2 = float64(qtyConversion.Qty2)
		vResp.Qty3 = float64(qtyConversion.Qty3)

		data = append(data, vResp)
	}

	return data, total, lastPage, nil
}

func (service *StockOpnameServiceImpl) StoreV2(request entity.CreateStockOpnameV2) (err error) {
	c := context.Background()

	if len(request.PrincipalID) == 0 && len(request.PLLane) == 0 && len(request.BrandID) == 0 && len(request.SBrand1ID) == 0 {
		return errors.New("at least one of principal_id, pl_lane, brand_id, or sbrand1_id must be provided")
	}

	scheduleDate, err := time.Parse("2006-01-02", request.ScheduleDate)
	if err != nil {
		return fmt.Errorf("invalid schedule_date format, expected YYYY-MM-DD: %w", err)
	}

	_, err = service.StockOpnameRepository.FindWarehouseByIDAndCustID(request.WhID, request.CustID)
	if err != nil {
		return fmt.Errorf("wh_id: %d, %s", request.WhID, err.Error())
	}

	proIDs := make([]int64, len(request.ProductList))
	for i, p := range request.ProductList {
		proIDs[i] = p.ProID
	}
	products, err := service.StockOpnameRepository.FindProductByIDs(c, proIDs, request.CustID)
	if err != nil {
		return fmt.Errorf("failed to get product prices: %w", err)
	}

	productMap := make(map[int64]model.Product)
	for _, p := range products {
		productMap[p.ProductId] = p
	}

	var missingProIDs []int64
	for _, productReq := range request.ProductList {
		if _, exists := productMap[productReq.ProID]; !exists {
			missingProIDs = append(missingProIDs, productReq.ProID)
		}
	}
	if len(missingProIDs) > 0 {
		return fmt.Errorf("products with ids %v not found for cust_id %s", missingProIDs, request.CustID)
	}

	includeZeroStockBool := request.IncludeZeroStock == 1
	now := time.Now()

	// Convert []int64 to pkg.Int64Array for PostgreSQL array support
	var principalIDArray pkg.Int64Array
	if len(request.PrincipalID) > 0 {
		principalIDArray = pkg.Int64Array(request.PrincipalID)
	}
	var plLaneArray pkg.Int64Array
	if len(request.PLLane) > 0 {
		plLaneArray = pkg.Int64Array(request.PLLane)
	}
	var brandIDArray pkg.Int64Array
	if len(request.BrandID) > 0 {
		brandIDArray = pkg.Int64Array(request.BrandID)
	}
	var sbrand1IDArray pkg.Int64Array
	if len(request.SBrand1ID) > 0 {
		sbrand1IDArray = pkg.Int64Array(request.SBrand1ID)
	}

	var stockOpnameModel model.StockOpnameV2
	var details []model.StockOpnameDetailV2
	var proIDsForStatusUpdate []int64

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		stockOpnameModel = model.StockOpnameV2{
			CustID:             request.CustID,
			DocNo:              "", // Will be auto-generated by BeforeCreate hook
			WhID:               request.WhID,
			StockType:          request.StockType,
			DataStatus:         1, // 1: Scheduled
			ProductHierarchy:   request.ProductHierarchy,
			IncludeZeroStock:   &includeZeroStockBool,
			PrincipalID:        principalIDArray,
			PLLane:             plLaneArray,
			BrandID:            brandIDArray,
			SBrand1ID:          sbrand1IDArray,
			DivisionID:         &request.DivisionID,
			InputBy:            request.InputBy,
			EmpID:              request.EmpID,
			CreatedBy:          &request.CreatedBy,
			CreatedAt:          &now,
			Notes:              nil,
			UpdatedBy:          nil,
			UpdatedAt:          nil,
			ScheduledAt:        &scheduleDate,
			AssignToEmpID:      nil,
			IsShowCurrentStock: nil,
			IsProcess:          false,
			IsRevised:          nil,
		}

		details = []model.StockOpnameDetailV2{}
		proIDsForStatusUpdate = []int64{}

		for _, productReq := range request.ProductList {
			product := productMap[productReq.ProID]

			convUnit1 := 0.0
			if productReq.ConvUnit1 != nil {
				convUnit1 = *productReq.ConvUnit1
			}
			convUnit2 := 0.0
			if productReq.ConvUnit2 != nil {
				convUnit2 = *productReq.ConvUnit2
			}
			convUnit3 := 0.0
			if productReq.ConvUnit3 != nil {
				convUnit3 = *productReq.ConvUnit3
			}

			qty1 := 0.0
			if productReq.Qty1 != nil {
				qty1 = *productReq.Qty1
			}
			qty2 := 0.0
			if productReq.Qty2 != nil {
				qty2 = *productReq.Qty2
			}
			qty3 := 0.0
			if productReq.Qty3 != nil {
				qty3 = *productReq.Qty3
			}

			detail := model.StockOpnameDetailV2{
				CustID:      request.CustID,
				DocNo:       "", // Will be set after StoreV2 generates doc_no
				ProID:       productReq.ProID,
				UnitID1:     productReq.UnitID1,
				UnitID2:     productReq.UnitID2,
				UnitID3:     productReq.UnitID3,
				ConvUnit1:   convUnit1,
				ConvUnit2:   convUnit2,
				ConvUnit3:   convUnit3,
				QtyStock1:   qty1,
				QtyStock2:   qty2,
				QtyStock3:   qty3,
				QtyOpname:   nil,
				PurchPrice1: product.PurchPrice1,
				PurchPrice2: product.PurchPrice2,
				PurchPrice3: product.PurchPrice3,
				CreatedBy:   &request.CreatedBy,
				CreatedAt:   &now,
			}
			details = append(details, detail)
			proIDsForStatusUpdate = append(proIDsForStatusUpdate, productReq.ProID)
		}

		err := service.StockOpnameRepository.StoreV2(txCtx, &stockOpnameModel)
		if err != nil {
			return fmt.Errorf("failed to store stock opname: %w", err)
		}

		for i := range details {
			details[i].DocNo = stockOpnameModel.DocNo
			err := service.StockOpnameRepository.StoreDetailV2(txCtx, &details[i])
			if err != nil {
				return fmt.Errorf("failed to store stock opname detail: %w", err)
			}
		}

		snapshot, err := service.StockOpnameRepository.GetProductStatusSnapshot(txCtx, proIDsForStatusUpdate, request.CustID)
		if err != nil {
			return fmt.Errorf("failed to get product status snapshot: %w", err)
		}
		if err := service.StockOpnameRepository.UpdateDetailsProductSnapshot(txCtx, stockOpnameModel.DocNo, request.CustID, snapshot); err != nil {
			return fmt.Errorf("failed to save product snapshot: %w", err)
		}

		err = service.StockOpnameRepository.UpdateProductStatus(txCtx, proIDsForStatusUpdate, request.CustID, 5)
		if err != nil {
			return fmt.Errorf("failed to update product status: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (service *StockOpnameServiceImpl) DetailV2(params entity.StockOpnameDetailV2Params) (data entity.StockOpnameDetailV2Response, err error) {
	header, err := service.StockOpnameRepository.FindDetailHeaderByDocNoV2(params)
	if err != nil {
		return data, err
	}

	products, err := service.StockOpnameRepository.FindDetailProductsByDocNoV2(params)
	if err != nil {
		return data, err
	}

	var divisionID int64 = 0
	if header.DivisionID != nil {
		divisionID = *header.DivisionID
	}

	var notes string = ""
	if header.Notes != nil {
		notes = *header.Notes
	}

	var subTotalSystem float64 = 0
	var subTotalPhysicalStock float64 = 0
	var difference float64 = 0

	var scheduledDate *string
	if header.ScheduledDate != nil {
		scheduledDate = header.ScheduledDate
	}

	data = entity.StockOpnameDetailV2Response{
		DocNo:         header.DocNo,
		CreatedDate:   header.CreatedDate,
		WhID:          header.WhID,
		WhCode:        header.WhCode,
		WhName:        header.WhName,
		StockType:     header.StockType,
		CreatedBy:     header.CreatedBy,
		UserName:      header.UserName,
		ScheduledDate: scheduledDate,
		Status:        header.Status,
		IsRevised:     header.IsRevised,
		IsProcess:     header.IsProcess,
		AssignTo: entity.StockOpnameDetailV2AssignTo{
			InputBy:      header.InputBy,
			DivisionID:   divisionID,
			DivisionName: header.DivisionName,
			EmpID:        header.EmpID,
			EmpName:      header.EmpName,
		},
		Notes: notes,
	}
	data.SetStatusDesc()

	for _, p := range products {
		qtyPhysical1 := p.QtyOpname1
		qtyPhysical2 := p.QtyOpname2
		qtyPhysical3 := p.QtyOpname3

		// qty1, qty2, qty3 berasal dari qty_stock1, qty_stock2, qty_stock3
		qty1 := p.QtyStock1
		qty2 := p.QtyStock2
		qty3 := p.QtyStock3

		differentStock1 := qtyPhysical1 - qty1
		differentStock2 := qtyPhysical2 - qty2
		differentStock3 := qtyPhysical3 - qty3

		differentPrice1 := differentStock1 * p.PurchPrice1
		differentPrice2 := differentStock2 * p.PurchPrice2
		differentPrice3 := differentStock3 * p.PurchPrice3

		// Hitung sub_total_system: (qty_stock1 * purch_price1) + (qty_stock2 * purch_price2) + (qty_stock3 * purch_price3)
		subTotalSystem += (qty1 * p.PurchPrice1) + (qty2 * p.PurchPrice2) + (qty3 * p.PurchPrice3)

		// Hitung sub_total_physical_stock: (qty_so1 * purch_price1) + (qty_so2 * purch_price2) + (qty_so3 * purch_price3)
		subTotalPhysicalStock += (qtyPhysical1 * p.PurchPrice1) + (qtyPhysical2 * p.PurchPrice2) + (qtyPhysical3 * p.PurchPrice3)

		// Hitung difference: ((qty_stock1-qty_so1) * purch_price1) + ((qty_stock2-qty_so2) * purch_price2) + ((qty_stock3-qty_so3) * purch_price3)
		difference += differentPrice1 + differentPrice2 + differentPrice3

		productItem := entity.StockOpnameDetailV2ProductItem{
			StockOpnameDetailID: p.StockOpnameDetailID,
			ProID:               p.ProID,
			ProCode:             p.ProCode,
			ProName:             p.ProName,
			UnitID1:             p.UnitID1,
			UnitID2:             p.UnitID2,
			UnitID3:             p.UnitID3,
			SellPrice1:          p.SellPrice1,
			SellPrice2:          p.SellPrice2,
			SellPrice3:          p.SellPrice3,
			UnitName1:           p.UnitName1,
			UnitName2:           p.UnitName2,
			UnitName3:           p.UnitName3,
			ConvUnit2:           p.ConvUnit2,
			ConvUnit3:           p.ConvUnit3,
			Qty1:                qty1,
			Qty2:                qty2,
			Qty3:                qty3,
			QtyPhysical1:        qtyPhysical1,
			QtyPhysical2:        qtyPhysical2,
			QtyPhysical3:        qtyPhysical3,
			DifferentStock1:     differentStock1,
			DifferentStock2:     differentStock2,
			DifferentStock3:     differentStock3,
			DifferentPrice1:     differentPrice1,
			DifferentPrice2:     differentPrice2,
			DifferentPrice3:     differentPrice3,
		}
		data.ProductList = append(data.ProductList, productItem)
	}

	// Set calculated totals
	data.SubTotalSystem = subTotalSystem
	data.SubTotalPhysicalStock = subTotalPhysicalStock
	data.Difference = difference

	return data, nil
}

func (service *StockOpnameServiceImpl) UpdateStatusV2(params entity.UpdateStockOpnameStatusV2Params, request entity.UpdateStockOpnameStatusV2Request) error {
	c := context.Background()

	stockOpname, err := service.StockOpnameRepository.FindStockOpnameForUpdate(c, params.DocNo, params.CustID)
	if err != nil {
		return err
	}

	if (request.IsProcess != nil && *request.IsProcess) || (request.IsAssigne != nil && *request.IsAssigne) {
		if stockOpname.CreatedBy == nil || *stockOpname.CreatedBy != params.UserID {
			return errors.New("you are not authorized to process or assign this stock opname. Only the creator can perform this action")
		}
	}

	var newStatus int
	var isProcess bool = stockOpname.IsProcess
	var logTitle string
	var logStatus int
	var shouldLog bool = false
	var shouldCompleteStock bool = false
	var shouldUpdateProductStatus bool = false

	if request.IsProcess != nil && *request.IsProcess {
		proIDs, err := service.StockOpnameRepository.GetProductIDsFromDetail(c, params.DocNo, params.CustID)
		if err != nil {
			return fmt.Errorf("failed to get product IDs: %w", err)
		}

		if len(proIDs) > 0 {
			// Check order status based on data_status
			hasStatusInRange, _, err := service.StockOpnameRepository.CheckProductsOrderStatus(c, proIDs, params.CustID)
			if err != nil {
				return fmt.Errorf("failed to check order status: %w", err)
			}

			if hasStatusInRange {
				roNos, err := service.StockOpnameRepository.GetBlockingOrderRos(c, proIDs, params.CustID)
				if err != nil {
					return fmt.Errorf("failed to get blocking ro_no: %w", err)
				}
				if len(roNos) > 0 {
					return fmt.Errorf("cannot process stock opname because there are still active orders for ro_no: %s", strings.Join(roNos, ", "))
				}
				return fmt.Errorf("cannot process stock opname because there are still active orders")
			}
			isProcess = true
		} else {
			isProcess = true
		}
		newStatus = request.OldStatus
		shouldLog = true
		logTitle = "Processed"
		logStatus = request.OldStatus
	} else if request.IsAssigne != nil && *request.IsAssigne {
		newStatus = entity.StockOpnameStatusAssign
		logStatus = entity.StockOpnameStatusAssign
		logTitle = "Assigne"
		shouldLog = true
	} else if request.IsCompleted != nil && *request.IsCompleted {
		newStatus = entity.StockOpnameStatusCompleted
		logStatus = entity.StockOpnameStatusCompleted
		logTitle = "Completed"
		shouldLog = true
		shouldCompleteStock = true
	} else if request.IsCancelled != nil && *request.IsCancelled {
		newStatus = entity.StockOpnameStatusRejected
		logStatus = entity.StockOpnameStatusRejected
		logTitle = "Cancelled"
		shouldLog = true
		shouldUpdateProductStatus = true
	} else {
		return errors.New("no action specified")
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		refID := primitive.NewObjectID().Hex()

		err := service.StockOpnameRepository.UpdateStockOpnameStatusV2(txCtx, params.DocNo, params.CustID, newStatus, isProcess, params.UserID)
		if err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}

		if shouldLog {
			now := time.Now()

			logEntry := &model.StockOpnameLog{
				Title:           logTitle,
				ExecutionTime:   now,
				OldStatus:       request.OldStatus,
				Status:          logStatus,
				RefID:           refID,
				TransactionCode: params.DocNo,
				RefTableName:    "inv.stock_opname",
				TriggeredBy:     "MANUAL",
				CreatedAt:       &now,
				CreatedBy:       &params.UserID,
				CustID:          params.CustID,
			}

			err = service.StockOpnameRepository.InsertStockOpnameLog(txCtx, logEntry)
			if err != nil {
				return fmt.Errorf("failed to insert log: %w", err)
			}
		}

		if shouldCompleteStock {
			err = service.handleStockCompletion(txCtx, params.DocNo, params.CustID, params.ParentCustID)
			if err != nil {
				return fmt.Errorf("failed to complete stock: %w", err)
			}
		}

		if shouldUpdateProductStatus {
			if err := service.StockOpnameRepository.RestoreProductStatusFromSnapshot(txCtx, params.DocNo, params.CustID); err != nil {
				return fmt.Errorf("failed to restore product status on cancel: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func resolveOpnamePhysicalQtyUOM(d model.StockOpnameDetailForCompleted) (q1, q2, q3 float64) {
	if d.RevisedDate != nil {
		if d.QtyRevised1 != nil {
			q1 = *d.QtyRevised1
		}
		if d.QtyRevised2 != nil {
			q2 = *d.QtyRevised2
		}
		if d.QtyRevised3 != nil {
			q3 = *d.QtyRevised3
		}
		return
	}
	if d.QtyRevised1 != nil {
		q1 = *d.QtyRevised1
	} else if d.QtySO1 != nil {
		q1 = *d.QtySO1
	}
	if d.QtyRevised2 != nil {
		q2 = *d.QtyRevised2
	} else if d.QtySO2 != nil {
		q2 = *d.QtySO2
	}
	if d.QtyRevised3 != nil {
		q3 = *d.QtyRevised3
	} else if d.QtySO3 != nil {
		q3 = *d.QtySO3
	}
	return
}

func (service *StockOpnameServiceImpl) handleStockCompletion(c context.Context, docNo, custID, parentCustID string) error {
	header, err := service.StockOpnameRepository.FindDetailHeaderByDocNoV2(entity.StockOpnameDetailV2Params{
		DocNo:        docNo,
		CustID:       custID,
		ParentCustID: custID,
	})
	if err != nil {
		return fmt.Errorf("failed to get stock opname header: %w", err)
	}

	details, err := service.StockOpnameRepository.GetStockOpnameDetailsForCompleted(c, docNo, custID, parentCustID)
	if err != nil {
		return fmt.Errorf("failed to get stock opname details: %w", err)
	}

	if len(details) == 0 {
		return nil // No details to process
	}

	// Get product IDs
	proIDs := make([]int64, len(details))
	for i, d := range details {
		proIDs[i] = d.ProID
	}

	priceMap, err := service.StockOpnameRepository.GetProductPrices(c, proIDs, custID, parentCustID)
	if err != nil {
		return fmt.Errorf("failed to get product prices: %w", err)
	}

	now := time.Now()
	var stocks []*model.Stock

	for _, detail := range details {
		qtySO1, qtySO2, qtySO3 := resolveOpnamePhysicalQtyUOM(detail)

		conv2 := detail.ConvUnit2
		if conv2 == 0 {
			conv2 = 1
		}
		conv3 := detail.ConvUnit3
		if conv3 == 0 {
			conv3 = 1
		}

		qtySOInSmall := qtySO1 + (qtySO2 * conv2) + (qtySO3 * conv2 * conv3)

		whStock, err := service.StockOpnameRepository.GetWarehouseStock(c, header.WhID, detail.ProID, custID)
		if err != nil {
			return fmt.Errorf("failed to get warehouse stock: %w", err)
		}

		totalDiff := qtySOInSmall - whStock.Qty

		productPrice, ok := priceMap[detail.ProID]
		if !ok {
			productPrice = model.ProductPrice{ProID: detail.ProID, SellPrice1: 0, Cogs: 0}
		}

		if totalDiff != 0 {
			stock := &model.Stock{
				CustID:      custID,
				StockDate:   now,
				TrCode:      "SO",
				TrNo:        docNo,
				WhID:        header.WhID,
				ProID:       detail.ProID,
				ItemCdn:     0,
				UnitPrice:   productPrice.SellPrice1,
				Cogs:        productPrice.Cogs,
				RefDetId:    detail.StockOpnameDetailID,
				QtyInOrder:  0,
				QtyOutOrder: 0,
			}
			if totalDiff > 0 {
				stock.QtyIn = totalDiff
				stock.QtyOut = 0
			} else {
				stock.QtyIn = 0
				stock.QtyOut = -totalDiff
			}
			stocks = append(stocks, stock)
		}

		if err := service.WarehouseStockRepository.UpdateQtyOnly(c, custID, header.WhID, detail.ProID, qtySOInSmall); err != nil {
			return fmt.Errorf("failed to update warehouse stock: %w", err)
		}
	}

	if len(stocks) > 0 {
		err = service.StockRepository.StoreBulk(c, stocks)
		if err != nil {
			return fmt.Errorf("failed to insert stocks: %w", err)
		}
	}

	if err := service.StockOpnameRepository.RestoreProductStatusFromSnapshot(c, docNo, custID); err != nil {
		return fmt.Errorf("failed to restore product status: %w", err)
	}

	return nil
}

func (service *StockOpnameServiceImpl) RevisedV2(params entity.RevisedStockOpnameV2Params, request entity.RevisedStockOpnameV2Request) error {
	c := context.Background()

	stockOpname, err := service.StockOpnameRepository.FindStockOpnameForUpdate(c, params.DocNo, params.CustID)
	if err != nil {
		return err
	}

	if stockOpname.CreatedBy == nil || *stockOpname.CreatedBy != params.UserID {
		return errors.New("you are not authorized to revised this stock opname. Only the creator can perform this action")
	}

	for _, item := range request.Data {
		_, err = service.StockOpnameRepository.FindStockOpnameDetailByProID(c, params.DocNo, params.CustID, item.ProID)
		if err != nil {
			return err
		}
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.StockOpnameRepository.UpdateStockOpnameIsRevised(txCtx, params.DocNo, params.CustID, true, params.UserID)
		if err != nil {
			return fmt.Errorf("failed to update is_revised: %w", err)
		}

		for _, item := range request.Data {
			err = service.StockOpnameRepository.UpdateStockOpnameDetailRevised(
				txCtx,
				params.DocNo,
				params.CustID,
				item.ProID,
				item.QtyRevised1,
				item.QtyRevised2,
				item.QtyRevised3,
				params.UserID,
			)
			if err != nil {
				return fmt.Errorf("failed to update revised quantities for pro_id %d: %w", item.ProID, err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (service *StockOpnameServiceImpl) DownloadTemplate(params entity.StockOpnameTemplateDownloadParams) (response entity.StockOpnameTemplateDownloadResponse, err error) {
	var header model.StockOpnameDetailV2Header
	var products []model.StockOpnameDetailV2Product

	if params.DocNo != "" {
		detailParams := entity.StockOpnameDetailV2Params{
			DocNo:        params.DocNo,
			CustID:       params.CustID,
			ParentCustID: params.ParentCustID,
		}

		header, err = service.StockOpnameRepository.FindDetailHeaderByDocNoV2(detailParams)
		if err != nil {
			return response, fmt.Errorf("failed to get stock opname header: %w", err)
		}

		products, err = service.StockOpnameRepository.FindDetailProductsByDocNoV2(detailParams)
		if err != nil {
			return response, fmt.Errorf("failed to get stock opname products: %w", err)
		}
	}

	file, err := service.generateStockOpnameTemplateExcel(header, products)
	if err != nil {
		return response, fmt.Errorf("failed to generate excel file: %w", err)
	}

	var buf bytes.Buffer
	err = file.Write(&buf)
	if err != nil {
		return response, fmt.Errorf("failed to write excel file: %w", err)
	}
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Generate report name: DownloadStockOpnameTemplate-DDMMYY-3digitRunningNumber
	now := time.Now()
	dateStr := now.Format("020106") // DDMMYY
	sequenceNumber, err := getNextSequenceNumberStockOpnameTemplate(dateStr)
	if err != nil {
		return response, fmt.Errorf("failed to get sequence number: %w", err)
	}
	reportName := fmt.Sprintf("DownloadStockOpnameTemplate-%s-%03d", dateStr, sequenceNumber)

	response.Name = reportName + ".xlsx"
	response.Version = "v1.2"
	response.ExpiresAt = time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	response.DownloadBase64 = base64Str

	return response, nil
}

// DownloadReport generates the final Stock Opname report in Excel format
// and returns it as base64 with a generated report name.
func (service *StockOpnameServiceImpl) DownloadReport(params entity.StockOpnameDownloadParams) (response entity.StockOpnameDownloadResponse, err error) {
	detailParams := entity.StockOpnameDetailV2Params{
		DocNo:        params.DocNo,
		CustID:       params.CustID,
		ParentCustID: params.ParentCustID,
	}

	header, err := service.StockOpnameRepository.FindDetailHeaderByDocNoV2(detailParams)
	if err != nil {
		return response, fmt.Errorf("failed to get stock opname header: %w", err)
	}

	products, err := service.StockOpnameRepository.FindDetailProductsByDocNoV2(detailParams)
	if err != nil {
		return response, fmt.Errorf("failed to get stock opname products: %w", err)
	}

	file, err := service.generateStockOpnameReportExcel(header, products)
	if err != nil {
		return response, fmt.Errorf("failed to generate excel file: %w", err)
	}

	var buf bytes.Buffer
	if err := file.Write(&buf); err != nil {
		return response, fmt.Errorf("failed to write excel file: %w", err)
	}

	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Generate report name: DownloadStockOpname-DDMMYY-3digitRunningNumber
	now := time.Now()
	dateStr := now.Format("020106") // DDMMYY
	sequenceNumber, err := getNextSequenceNumberStockOpnameTemplate(dateStr)
	if err != nil {
		return response, fmt.Errorf("failed to get sequence number: %w", err)
	}
	reportName := fmt.Sprintf("DownloadStockOpname-%s-%03d", dateStr, sequenceNumber)

	response.ReportName = reportName + ".xlsx"
	response.FileBase64 = base64Str

	return response, nil
}

func (service *StockOpnameServiceImpl) generateStockOpnameTemplateExcel(header model.StockOpnameDetailV2Header, products []model.StockOpnameDetailV2Product) (*excelize.File, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetName := "Stock Opname Detail"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Create styles
	boldStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	if err != nil {
		return nil, err
	}

	boldBorderStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D3D3D3"},
			Pattern: 1,
		},
	})
	if err != nil {
		return nil, err
	}

	textBorderStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	numberBorderStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	// Title
	f.SetCellValue(sheetName, "A1", "Stock Opname Detail")
	f.SetCellStyle(sheetName, "A1", "A1", boldStyle)

	// Format dates to DD/MM/YYYY
	createdDateFormatted := ""
	if header.CreatedDate != "" {
		if t, err := time.Parse(time.RFC3339, header.CreatedDate); err == nil {
			createdDateFormatted = t.Format("02/01/2006")
		} else if t, err := time.Parse("2006-01-02T15:04:05Z", header.CreatedDate); err == nil {
			createdDateFormatted = t.Format("02/01/2006")
		} else {
			createdDateFormatted = header.CreatedDate
		}
	}

	scheduledDateFormatted := ""
	if header.ScheduledDate != nil && *header.ScheduledDate != "" {
		if t, err := time.Parse("2006-01-02 15:04:05+00", *header.ScheduledDate); err == nil {
			scheduledDateFormatted = t.Format("02/01/2006")
		} else if t, err := time.Parse("2006-01-02", *header.ScheduledDate); err == nil {
			scheduledDateFormatted = t.Format("02/01/2006")
		} else {
			scheduledDateFormatted = *header.ScheduledDate
		}
	}

	statusDesc := ""
	if header.Status > 0 {
		statusDesc = entity.StockOpnameStatusDesc[header.Status]
	}

	warehouseStr := ""
	if header.WhCode != "" || header.WhName != "" {
		warehouseStr = fmt.Sprintf("%s - %s", header.WhCode, header.WhName)
	}

	// Header section (Rows 2-8)
	row := 2
	headerLabels := []string{"Document No:", "Created Date:", "Created By:", "Warehouse:", "Stock Type:", "Schedule Date:", "Status:"}
	headerValues := []string{header.DocNo, createdDateFormatted, header.UserName, warehouseStr, header.StockType, scheduledDateFormatted, statusDesc}

	for i, label := range headerLabels {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), label)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), headerValues[i])
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), textBorderStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), textBorderStyle)
		row++
	}

	// Assign to section (Rows 10-14)
	row = 10
	f.SetCellValue(sheetName, "A10", "Assign to:")
	f.SetCellStyle(sheetName, "A10", "A10", boldStyle)

	// Format division_id and emp_id
	divisionIDStr := ""
	if header.DivisionID != nil {
		divisionIDStr = fmt.Sprintf("%d", *header.DivisionID)
	}
	empIDStr := fmt.Sprintf("%d", header.EmpID)

	assignLabels := []string{"Input By:", "Division:", "Division ID:", "Employee:", "Emp ID:"}
	assignValues := []string{header.InputBy, header.DivisionName, divisionIDStr, header.EmpName, empIDStr}

	for i, label := range assignLabels {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row+1), label)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row+1), assignValues[i])
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row+1), fmt.Sprintf("A%d", row+1), textBorderStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row+1), fmt.Sprintf("B%d", row+1), textBorderStyle)
		row++
	}

	row = 16
	// Merge cells for UOM header (D-F)
	f.SetCellValue(sheetName, "D16", "UOM")
	f.MergeCell(sheetName, "D16", "F16")
	f.SetCellStyle(sheetName, "D16", "F16", boldBorderStyle)

	// Merge cells for Quantity header (G-I)
	f.SetCellValue(sheetName, "G16", "Quantity")
	f.MergeCell(sheetName, "G16", "I16")
	f.SetCellStyle(sheetName, "G16", "I16", boldBorderStyle)

	row = 17
	// Header: No, Product Code, Product Name, UOM (Largest,Middle,Smallest), Quantity (Largest,Middle,Smallest)
	headers := []string{"No", "Product Code", "Product Name", "Largest", "Middle", "Smallest", "Largest", "Middle", "Smallest"}
	cols := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}

	for i, header := range headers {
		cell := fmt.Sprintf("%s%d", cols[i], row)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, boldBorderStyle)
	}

	// Product list data (Row 18 onwards): Quantity = qty_physical (QtyOpname1=Smallest, QtyOpname2=Middle, QtyOpname3=Largest)
	row = 18
	for _, product := range products {
		// No: stock_opname_detail_id
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), product.StockOpnameDetailID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), product.ProCode)
		// Product Name
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), product.ProName)

		// UOM with conversion info (D, E, F)
		uom3Text := product.UnitName3
		if product.UnitName2 != "" && product.ConvUnit3 > 0 {
			uom3Text = fmt.Sprintf("%s\n1 %s = %.0f %s", product.UnitName3, product.UnitName3, product.ConvUnit3, product.UnitName2)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), uom3Text)

		uom2Text := product.UnitName2
		if product.UnitName1 != "" && product.ConvUnit2 > 0 {
			uom2Text = fmt.Sprintf("%s\n1 %s = %.0f %s", product.UnitName2, product.UnitName2, product.ConvUnit2, product.UnitName1)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), uom2Text)

		uom1Text := product.UnitName1
		if product.UnitName2 != "" && product.ConvUnit2 > 0 && product.ConvUnit3 > 0 {
			uom1Text = fmt.Sprintf("%s\n1 %s = %.0f * %.0f %s", product.UnitName1, product.UnitName2, product.ConvUnit2, product.ConvUnit3, product.UnitName1)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), uom1Text)

		// Quantity: Smallest=qty_physical1, Middle=qty_physical2, Largest=qty_physical3 (G=Largest, H=Middle, I=Smallest)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), product.QtyOpname3)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), product.QtyOpname2)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), product.QtyOpname1)

		// Apply styles
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), textBorderStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("F%d", row), textBorderStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("G%d", row), fmt.Sprintf("I%d", row), numberBorderStyle)
		row++
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 8)
	f.SetColWidth(sheetName, "B", "B", 15)
	f.SetColWidth(sheetName, "C", "C", 30)
	f.SetColWidth(sheetName, "D", "F", 15) // UOM columns
	f.SetColWidth(sheetName, "G", "I", 12) // Quantity columns

	return f, nil
}

func (service *StockOpnameServiceImpl) generateStockOpnameReportExcel(header model.StockOpnameDetailV2Header, products []model.StockOpnameDetailV2Product) (*excelize.File, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetName := "Stock Opname Detail"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Styles
	boldStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	if err != nil {
		return nil, err
	}

	boldBorderStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D3D3D3"},
			Pattern: 1,
		},
	})
	if err != nil {
		return nil, err
	}

	textBorderStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	numberBorderStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	// Title
	f.SetCellValue(sheetName, "A1", "Stock Opname Detail")
	f.SetCellStyle(sheetName, "A1", "A1", boldStyle)

	// Date formatting
	createdDateFormatted := ""
	if header.CreatedDate != "" {
		if t, err := time.Parse(time.RFC3339, header.CreatedDate); err == nil {
			createdDateFormatted = t.Format("02/01/2006")
		} else if t, err := time.Parse("2006-01-02T15:04:05Z", header.CreatedDate); err == nil {
			createdDateFormatted = t.Format("02/01/2006")
		} else {
			createdDateFormatted = header.CreatedDate
		}
	}

	scheduledDateFormatted := ""
	if header.ScheduledDate != nil && *header.ScheduledDate != "" {
		if t, err := time.Parse("2006-01-02 15:04:05+00", *header.ScheduledDate); err == nil {
			scheduledDateFormatted = t.Format("02/01/2006")
		} else if t, err := time.Parse("2006-01-02", *header.ScheduledDate); err == nil {
			scheduledDateFormatted = t.Format("02/01/2006")
		} else {
			scheduledDateFormatted = *header.ScheduledDate
		}
	}

	statusDesc := ""
	if header.Status > 0 {
		statusDesc = entity.StockOpnameStatusDesc[header.Status]
	}

	warehouseStr := ""
	if header.WhCode != "" || header.WhName != "" {
		warehouseStr = fmt.Sprintf("%s - %s", header.WhCode, header.WhName)
	}

	// Header section (Rows 2-8)
	row := 2
	headerLabels := []string{"Document No:", "Created Date:", "Created By:", "Warehouse:", "Stock Type:", "Schedule Date:", "Status:"}
	headerValues := []string{header.DocNo, createdDateFormatted, header.UserName, warehouseStr, header.StockType, scheduledDateFormatted, statusDesc}

	for i, label := range headerLabels {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), label)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), headerValues[i])
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), textBorderStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), textBorderStyle)
		row++
	}

	// Assign to section
	row = 10
	f.SetCellValue(sheetName, "A10", "Assign to:")
	f.SetCellStyle(sheetName, "A10", "A10", boldStyle)

	divisionIDStr := ""
	if header.DivisionID != nil {
		divisionIDStr = fmt.Sprintf("%d", *header.DivisionID)
	}
	empIDStr := fmt.Sprintf("%d", header.EmpID)

	assignLabels := []string{"Input By:", "Division:", "Division ID:", "Employee:", "Emp ID:"}
	assignValues := []string{header.InputBy, header.DivisionName, divisionIDStr, header.EmpName, empIDStr}

	for i, label := range assignLabels {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row+1), label)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row+1), assignValues[i])
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row+1), fmt.Sprintf("A%d", row+1), textBorderStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row+1), fmt.Sprintf("B%d", row+1), textBorderStyle)
		row++
	}

	// Section headers
	row = 16
	f.SetCellValue(sheetName, "D16", "UOM")
	f.MergeCell(sheetName, "D16", "F16")
	f.SetCellStyle(sheetName, "D16", "F16", boldBorderStyle)

	f.SetCellValue(sheetName, "G16", "System Quantity")
	f.MergeCell(sheetName, "G16", "I16")
	f.SetCellStyle(sheetName, "G16", "I16", boldBorderStyle)

	f.SetCellValue(sheetName, "J16", "Physical Quantity")
	f.MergeCell(sheetName, "J16", "L16")
	f.SetCellStyle(sheetName, "J16", "L16", boldBorderStyle)

	f.SetCellValue(sheetName, "M16", "Difference")
	f.MergeCell(sheetName, "M16", "P16")
	f.SetCellStyle(sheetName, "M16", "P16", boldBorderStyle)

	// Column headers
	row = 17
	headers := []string{
		"No", "Product Code", "Product Name",
		"Largest", "Middle", "Smallest",
		"Largest", "Middle", "Smallest", // System
		"Largest", "Middle", "Smallest", // Physical
		"Largest", "Middle", "Smallest", "In Price", // Difference
	}
	cols := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P"}

	for i, h := range headers {
		cell := fmt.Sprintf("%s%d", cols[i], row)
		f.SetCellValue(sheetName, cell, h)
		f.SetCellStyle(sheetName, cell, cell, boldBorderStyle)
	}

	// Helper for physical qty (qty_revised if not nil and not zero, else qty_so)
	getPhysical := func(revised *float64, so float64) float64 {
		if revised != nil && *revised != 0 {
			return *revised
		}
		return so
	}

	sort.Slice(products, func(i, j int) bool {
		return products[i].ProCode < products[j].ProCode
	})

	// Data rows
	row = 18
	for idx, product := range products {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), idx+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), product.ProCode)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), product.ProName)

		// UOM section
		uom3Text := product.UnitName3
		if product.UnitName2 != "" && product.ConvUnit3 > 0 {
			uom3Text = fmt.Sprintf("%s\n1 %s = %.0f %s", product.UnitName3, product.UnitName3, product.ConvUnit3, product.UnitName2)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), uom3Text)

		uom2Text := product.UnitName2
		if product.UnitName1 != "" && product.ConvUnit2 > 0 {
			uom2Text = fmt.Sprintf("%s\n1 %s = %.0f %s", product.UnitName2, product.UnitName2, product.ConvUnit2, product.UnitName1)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), uom2Text)

		uom1Text := product.UnitName1
		if product.UnitName2 != "" && product.ConvUnit2 > 0 && product.ConvUnit3 > 0 {
			uom1Text = fmt.Sprintf("%s\n1 %s = %.0f * %.0f %s", product.UnitName1, product.UnitName2, product.ConvUnit2, product.ConvUnit3, product.UnitName1)
		}
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), uom1Text)

		// System quantity (G,H,I)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), product.QtyStock3)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), product.QtyStock2)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), product.QtyStock1)

		// Physical quantity (J,K,L) using revised or SO
		phys3 := getPhysical(product.QtyRevised3, product.QtyOpname3)
		phys2 := getPhysical(product.QtyRevised2, product.QtyOpname2)
		phys1 := getPhysical(product.QtyRevised1, product.QtyOpname1)

		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), phys3)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), phys2)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), phys1)

		// Difference: Physical Quantity - System Quantity for each unit
		diff3 := phys3 - product.QtyStock3 // Largest
		diff2 := phys2 - product.QtyStock2 // Middle
		diff1 := phys1 - product.QtyStock1 // Smallest

		f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), diff3)
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), diff2)
		f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), diff1)

		// In Price: (Diff Largest * SellPrice3) + (Diff Middle * SellPrice2) + (Diff Smallest * SellPrice1)
		inPrice := (diff3 * product.SellPrice3) + (diff2 * product.SellPrice2) + (diff1 * product.SellPrice1)
		f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), inPrice)

		// Styles
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("C%d", row), textBorderStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("F%d", row), textBorderStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("G%d", row), fmt.Sprintf("P%d", row), numberBorderStyle)

		row++
	}

	// Column widths
	f.SetColWidth(sheetName, "A", "A", 8)
	f.SetColWidth(sheetName, "B", "B", 15)
	f.SetColWidth(sheetName, "C", "C", 30)
	f.SetColWidth(sheetName, "D", "F", 15)
	f.SetColWidth(sheetName, "G", "P", 12)

	return f, nil
}

var (
	stockOpnameTemplateSequenceMutex sync.Mutex
	stockOpnameTemplateSequenceFile  = ".download_stock_opname_template_sequence.json"
)

type StockOpnameTemplateSequenceStorage struct {
	Sequences map[string]int `json:"sequences"` // map[DDMMYY]sequenceNumber
}

func getNextSequenceNumberStockOpnameTemplate(dateStr string) (int, error) {
	stockOpnameTemplateSequenceMutex.Lock()
	defer stockOpnameTemplateSequenceMutex.Unlock()

	storage, err := readStockOpnameTemplateSequenceStorage()
	if err != nil {
		return 0, err
	}

	if storage.Sequences == nil {
		storage.Sequences = make(map[string]int)
	}

	currentSeq := storage.Sequences[dateStr]

	nextSeq := currentSeq + 1
	storage.Sequences[dateStr] = nextSeq

	err = writeStockOpnameTemplateSequenceStorage(storage)
	if err != nil {
		return 0, err
	}

	return nextSeq, nil
}

func readStockOpnameTemplateSequenceStorage() (*StockOpnameTemplateSequenceStorage, error) {
	storage := &StockOpnameTemplateSequenceStorage{
		Sequences: make(map[string]int),
	}

	if _, err := os.Stat(stockOpnameTemplateSequenceFile); os.IsNotExist(err) {
		return storage, nil
	}

	data, err := os.ReadFile(stockOpnameTemplateSequenceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read sequence file: %w", err)
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, storage)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sequence file: %w", err)
		}
	}

	return storage, nil
}

func writeStockOpnameTemplateSequenceStorage(storage *StockOpnameTemplateSequenceStorage) error {
	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sequence storage: %w", err)
	}

	err = os.WriteFile(stockOpnameTemplateSequenceFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write sequence file: %w", err)
	}

	return nil
}
func (service *StockOpnameServiceImpl) StartV2(params entity.StartStockOpnameV2Params) (entity.StartStockOpnameV2Response, error) {
	c := context.Background()
	var response entity.StartStockOpnameV2Response

	stockOpname, err := service.StockOpnameRepository.FindStockOpnameForStart(c, params.DocNo, params.CustID)
	if err != nil {
		return response, err
	}

	if stockOpname.DataStatus != entity.StockOpnameStatusAssign {
		return response, constant.ErrStockOpnameCannotBeStarted
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.StockOpnameRepository.UpdateStockOpnameStart(txCtx, params.DocNo, params.CustID, params.UserID)
		if err != nil {
			return fmt.Errorf("failed to start stock opname: %w", err)
		}

		return nil
	})

	if err != nil {
		return response, err
	}

	now := time.Now()
	response = entity.StartStockOpnameV2Response{
		DocNo:         stockOpname.DocNo,
		WhID:          stockOpname.WhID,
		WhCode:        stockOpname.WhCode,
		WhName:        stockOpname.WhName,
		Status:        entity.StockOpnameStatusOnGoing,
		StatusDesc:    entity.StockOpnameStatusDesc[entity.StockOpnameStatusOnGoing],
		StartedAt:     now.Format("2006-01-02 15:04:05"),
		StartedBy:     params.UserID,
		StartedByName: "",
	}

	return response, nil
}

func (service *StockOpnameServiceImpl) SubmitV2(params entity.SubmitStockOpnameV2Params, request entity.SubmitStockOpnameV2Request, fileHeader *multipart.FileHeader) (entity.SubmitStockOpnameV2Response, error) {
	c := context.Background()
	var response entity.SubmitStockOpnameV2Response

	_, err := service.StockOpnameRepository.FindStockOpnameForUpdate(c, params.DocNo, params.CustID)
	if err != nil {
		return response, fmt.Errorf("stock opname not found: %w", err)
	}

	detailIDs := make([]int64, len(request.Details))
	for i, detail := range request.Details {
		detailIDs[i] = detail.StockOpnameDetID
	}

	validIDs, err := service.StockOpnameRepository.ValidateStockOpnameDetailIDs(c, params.DocNo, params.CustID, detailIDs)
	if err != nil {
		return response, fmt.Errorf("failed to validate detail IDs: %w", err)
	}

	validIDMap := make(map[int64]bool)
	for _, id := range validIDs {
		validIDMap[id] = true
	}

	var invalidIDs []int64
	var validDetails []model.StockOpnameDetailQtySO
	for _, detail := range request.Details {
		if !validIDMap[detail.StockOpnameDetID] {
			invalidIDs = append(invalidIDs, detail.StockOpnameDetID)
		} else {
			validDetails = append(validDetails, model.StockOpnameDetailQtySO{
				StockOpnameDetID: detail.StockOpnameDetID,
				QtySO1:           *detail.QtySO1,
				QtySO2:           *detail.QtySO2,
				QtySO3:           *detail.QtySO3,
			})
		}
	}

	if len(invalidIDs) > 0 {
		return response, fmt.Errorf("stock_opname_det_id %v not found in document %s", invalidIDs, params.DocNo)
	}

	// Upload file to OBS if provided
	var filePath string
	if fileHeader != nil {
		const maxUploadedByLen = 50
		const maxFilePathLen = 50

		file, err := fileHeader.Open()
		if err != nil {
			return response, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		uploadModel := &model.Upload{
			Folder: "stock_opname",
			File:   fileHeader,
		}

		filePath, err = service.ObsAdapter.UploadFile(uploadModel)
		if err != nil {
			return response, fmt.Errorf("failed to upload file to OBS: %w", err)
		}

		const maxBulkUploadFilePathLen = 50
		if len(filePath) > maxBulkUploadFilePathLen {
			filePath = filePath[len(filePath)-maxBulkUploadFilePathLen:]
		}
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.StockOpnameRepository.UpdateStockOpnameDetailsQtySO(txCtx, params.DocNo, params.CustID, validDetails)
		if err != nil {
			return fmt.Errorf("failed to update stock opname details: %w", err)
		}

		err = service.StockOpnameRepository.UpdateStockOpnameStatusToSubmit(txCtx, params.DocNo, params.CustID, params.UserID)
		if err != nil {
			return fmt.Errorf("failed to update stock opname status: %w", err)
		}

		if filePath != "" {
			status := 5
			bulkUpload := model.StockOpnameBulkUpload{
				DocNo:      params.DocNo,
				FilePath:   filePath,
				UploadedBy: fmt.Sprintf("%d", params.UserID), // Assuming UserID is int64, convert to string
				UploadedAt: time.Now().UTC().Unix(),
				Status:     &status, // Submitted status
			}
			err = service.StockOpnameRepository.InsertBulkUpload(txCtx, &bulkUpload)
			if err != nil {
				return fmt.Errorf("failed to insert stock opname bulk upload: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return response, err
	}

	response = entity.SubmitStockOpnameV2Response{
		DocNo:          params.DocNo,
		UpdatedCount:   len(validDetails),
		UpdatedDetails: make([]entity.SubmitStockOpnameV2Detail, len(validDetails)),
	}

	validDetailMap := make(map[int64]model.StockOpnameDetailQtySO)
	for _, d := range validDetails {
		validDetailMap[d.StockOpnameDetID] = d
	}

	idx := 0
	for _, reqDetail := range request.Details {
		if !validIDMap[reqDetail.StockOpnameDetID] {
			continue
		}
		d := validDetailMap[reqDetail.StockOpnameDetID]
		q1, q2, q3 := d.QtySO1, d.QtySO2, d.QtySO3
		response.UpdatedDetails[idx] = entity.SubmitStockOpnameV2Detail{
			StockOpnameDetID: reqDetail.StockOpnameDetID,
			ProID:            reqDetail.ProID,
			QtySO1:           &q1,
			QtySO2:           &q2,
			QtySO3:           &q3,
		}
		idx++
	}

	return response, nil
}

const (
	maxBulkUploadFileSizeMB = 100
)

type bulkUploadExcelRow struct {
	DetailID     int64
	ProID        int64
	QtyPhysical1 float64
	QtyPhysical2 float64
	QtyPhysical3 float64
	Valid        bool
	ErrMsg       string
}

func (service *StockOpnameServiceImpl) BulkUpload(params entity.BulkUploadStockOpnameV2Params, file []byte, filename string) (entity.BulkUploadStockOpnameV2Response, error) {
	c := context.Background()
	var response entity.BulkUploadStockOpnameV2Response
	response.DocNo = params.DocNo
	response.ProcessedAt = time.Now().UTC().Format(time.RFC3339)

	// File type: allow .xlsx (requirement says .xlsx in Body)
	lower := strings.ToLower(filename)
	if !strings.HasSuffix(lower, ".xlsx") {
		return response, fmt.Errorf("invalid file format: only .xlsx is allowed")
	}
	// Max size 100MB
	if len(file) > maxBulkUploadFileSizeMB*1024*1024 {
		return response, fmt.Errorf("uploaded file exceeds maximum size limit (max %d MB)", maxBulkUploadFileSizeMB)
	}

	// Stock opname must exist and status = 4 (On Going)
	soForUpdate, err := service.StockOpnameRepository.FindStockOpnameForUpdate(c, params.DocNo, params.CustID)
	if err != nil {
		return response, fmt.Errorf("stock opname not found: %w", err)
	}
	if soForUpdate.DataStatus != entity.StockOpnameStatusOnGoing {
		statusDesc := entity.StockOpnameStatusDesc[soForUpdate.DataStatus]
		if statusDesc == "" {
			statusDesc = "Unknown"
		}
		return response, fmt.Errorf("stock opname cannot be updated in current status (status: %s)", statusDesc)
	}

	// Get current details for validation and response
	detailParams := entity.StockOpnameDetailV2Params{
		DocNo:        params.DocNo,
		CustID:       params.CustID,
		ParentCustID: params.ParentCustID,
	}
	header, err := service.StockOpnameRepository.FindDetailHeaderByDocNoV2(detailParams)
	if err != nil {
		return response, fmt.Errorf("failed to get stock opname header: %w", err)
	}
	products, err := service.StockOpnameRepository.FindDetailProductsByDocNoV2(detailParams)
	if err != nil {
		return response, fmt.Errorf("failed to get stock opname details: %w", err)
	}
	detailMap := make(map[int64]model.StockOpnameDetailV2Product)
	proIDMap := make(map[int64]model.StockOpnameDetailV2Product)    // Fallback map by ProID
	proCodeMap := make(map[string]model.StockOpnameDetailV2Product) // Lookup by Product Code (template uses ProCode in col B)
	for _, p := range products {
		detailMap[p.StockOpnameDetailID] = p
		proIDMap[p.ProID] = p
		if p.ProCode != "" {
			proCodeMap[strings.TrimSpace(p.ProCode)] = p
		}
	}

	// Parse Excel: sheet "Stock Opname Detail". A=No, B=Product Code, C=Product Name, D-F=UOM, G=Largest(qty3), H=Middle(qty2), I=Smallest(qty1)
	f, err := excelize.OpenReader(bytes.NewReader(file))
	if err != nil {
		return response, fmt.Errorf("failed to read excel file: %w", err)
	}
	defer f.Close()
	sheetName := "Stock Opname Detail"
	rows, err := f.GetRows(sheetName)
	if err != nil {
		// Try fallback to "Sheet1"
		sheetName = "Sheet1"
		rows, err = f.GetRows(sheetName)
		if err != nil {
			// If still fails, try getting the first sheet name
			if len(f.GetSheetMap()) > 0 {
				sheetName = f.GetSheetName(0)
				rows, err = f.GetRows(sheetName)
			}
		}
	}

	if err != nil {
		return response, fmt.Errorf("sheet %q not found or invalid, and no fallback sheet available", "Stock Opname Detail")
	}

	var parsedRows []bulkUploadExcelRow
	startRow := 18

	// Dynamic header detection
	// 1. Find the row containing "No" in the first column
	headerIndex := -1
	for i, row := range rows {
		if len(row) > 0 && strings.ToLower(strings.TrimSpace(row[0])) == "no" {
			headerIndex = i
			break
		}
	}

	if headerIndex != -1 {
		// 2. Find the first valid data row after the header
		foundDataRow := false
		for i := headerIndex + 1; i < len(rows); i++ {
			row := rows[i]
			if len(row) < 2 {
				continue
			}
			cell1 := strings.TrimSpace(row[1])
			if cell1 == "" {
				continue
			}
			if _, hasProCode := proCodeMap[cell1]; hasProCode {
				startRow = i + 1 // i is 0-based index, startRow is 1-based
				foundDataRow = true
				break
			}
		}

		// If no valid data row found, but header was found, maybe the file is empty of data?
		// Default to headerIndex + 2 (assuming 1 header row) so we can at least process (and likely fail) the next row
		if !foundDataRow {
			startRow = headerIndex + 2
		}
	} else {
		// Header "No" not found.
		// If rows < 18, we can't rely on default startRow=18.
		if len(rows) < 18 {
			response.Status = "FAILED"
			response.StatusDescription = "Failed"
			response.TotalRow = 0
			response.SuccessRow = 0
			response.FailedRow = 0
			response.SOData = buildBulkUploadSOData(header, products)
			return response, nil
		}
		// If rows >= 18, we stick to default startRow = 18 (assuming standard template with header at 17)
	}

	hasInvalidID := false
	for i := startRow - 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 2 {
			continue
		}
		cell0 := strings.TrimSpace(row[0])
		cell1 := strings.TrimSpace(row[1])
		if cell0 == "" || cell1 == "" {
			continue
		}
		p, ok := proCodeMap[cell1]
		if !ok {
			hasInvalidID = true
			parsedRows = append(parsedRows, bulkUploadExcelRow{Valid: false, ErrMsg: "Product code not found in template"})
			continue
		}
		detailID := p.StockOpnameDetailID
		proID := p.ProID
		// Template: G=Largest(6), H=Middle(7), I=Smallest(8). DB: QtyPhysical1=Smallest, QtyPhysical2=Middle, QtyPhysical3=Largest
		qty1, e1 := parseFloatCell(row, 8) // I = Smallest
		qty2, e2 := parseFloatCell(row, 7) // H = Middle
		qty3, e3 := parseFloatCell(row, 6) // G = Largest
		if e1 != nil || e2 != nil || e3 != nil {
			parsedRows = append(parsedRows, bulkUploadExcelRow{DetailID: detailID, ProID: proID, Valid: false, ErrMsg: "Invalid quantity value"})
			continue
		}
		// Qty harus nominal (bilangan bulat): decimal di-floor, contoh 1.5 -> 1, 2.5 -> 2
		qty1 = math.Floor(qty1)
		qty2 = math.Floor(qty2)
		qty3 = math.Floor(qty3)
		if qty1 < 0 || qty2 < 0 || qty3 < 0 {
			parsedRows = append(parsedRows, bulkUploadExcelRow{DetailID: detailID, ProID: proID, QtyPhysical1: qty1, QtyPhysical2: qty2, QtyPhysical3: qty3, Valid: false, ErrMsg: "Quantity cannot be negative"})
			continue
		}
		parsedRows = append(parsedRows, bulkUploadExcelRow{
			DetailID:     detailID,
			ProID:        proID,
			QtyPhysical1: qty1,
			QtyPhysical2: qty2,
			QtyPhysical3: qty3,
			Valid:        true,
		})
	}

	response.TotalRow = len(parsedRows)
	successCount := 0
	for _, r := range parsedRows {
		if r.Valid {
			successCount++
		}
	}
	response.SuccessRow = successCount
	response.FailedRow = response.TotalRow - successCount
	response.RowErrors = buildBulkUploadRowErrors(parsedRows, startRow)

	if hasInvalidID {
		response.Status = "FAILED"
		response.StatusDescription = "Failed"
		response.SOData = buildBulkUploadSODataWithRowResults(header, products, detailMap, parsedRows)
		return response, nil
	}

	if response.TotalRow == 0 {
		// DEBUG: Return detailed debug info
		debugMsg := fmt.Sprintf("Debug: rows=%d, headerIndex=%d, startRow=%d", len(rows), headerIndex, startRow)
		if headerIndex != -1 && len(rows) > headerIndex {
			debugMsg += fmt.Sprintf(", HeaderRow=%v", rows[headerIndex])
		}
		if startRow-1 < len(rows) {
			debugMsg += fmt.Sprintf(", FirstDataRow=%v", rows[startRow-1])
		}

		response.Status = "FAILED"
		response.StatusDescription = "Failed"
		response.SOData = buildBulkUploadSOData(header, products)
		// Return error with debug info instead of just returning response
		return response, fmt.Errorf("no valid rows found. %s", debugMsg)
	}

	// Transaction: insert bulk_upload, bulk_upload_items
	statusOnGoing := 4
	validRow := successCount
	invalidRow := response.FailedRow
	uploadedBy := strconv.FormatInt(params.UserID, 10)
	bulkUpload := &model.StockOpnameBulkUpload{
		DocNo:      params.DocNo,
		FilePath:   "",
		Status:     &statusOnGoing,
		TotalRow:   &response.TotalRow,
		ValidRow:   &validRow,
		InvalidRow: &invalidRow,
		UploadedBy: uploadedBy,
		UploadedAt: time.Now().UTC().Unix(),
	}
	var bulkItems []model.StockOpnameBulkUploadItem

	for _, r := range parsedRows {
		if !r.Valid {
			continue
		}
		bulkItems = append(bulkItems, model.StockOpnameBulkUploadItem{
			UploadID:    0,
			ProductID:   r.ProID,
			QtySO1:      r.QtyPhysical1,
			QtySO2:      r.QtyPhysical2,
			QtySO3:      r.QtyPhysical3,
			QtyRevised1: 0,
			QtyRevised2: 0,
			QtyRevised3: 0,
			UnitID1:     0,
			UnitID2:     0,
			UnitID3:     0,
		})
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		if err := service.StockOpnameRepository.InsertBulkUpload(txCtx, bulkUpload); err != nil {
			return err
		}
		for i := range bulkItems {
			bulkItems[i].UploadID = bulkUpload.UploadID
		}
		if err := service.StockOpnameRepository.InsertBulkUploadItems(txCtx, bulkItems); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return response, fmt.Errorf("failed to process bulk upload: %w", err)
	}

	if response.FailedRow == 0 {
		response.Status = "FULL_SUCCESS"
		response.StatusDescription = "Success"
	} else {
		response.Status = "PARTIAL_SUCCESS"
		response.StatusDescription = "Partially Success"
	}
	response.RowErrors = buildBulkUploadRowErrors(parsedRows, startRow)
	response.SOData = buildBulkUploadSODataWithRowResults(header, products, detailMap, parsedRows)
	return response, nil
}

func parseFloatCell(row []string, colIndex int) (float64, error) {
	if colIndex >= len(row) {
		return 0, nil
	}
	s := strings.TrimSpace(row[colIndex])
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}

func buildBulkUploadRowErrors(parsedRows []bulkUploadExcelRow, startRow int) []entity.BulkUploadRowError {
	var out []entity.BulkUploadRowError
	for i, r := range parsedRows {
		if r.Valid || r.ErrMsg == "" {
			continue
		}
		out = append(out, entity.BulkUploadRowError{
			RowIndex:   startRow + i,
			DetailID:   r.DetailID,
			ProID:      r.ProID,
			ErrMessage: r.ErrMsg,
		})
	}
	return out
}

func buildBulkUploadSOData(header model.StockOpnameDetailV2Header, products []model.StockOpnameDetailV2Product) []entity.BulkUploadStockOpnameV2SOData {
	productList := make([]entity.BulkUploadStockOpnameV2ProductItem, 0, len(products))
	for _, p := range products {
		qty1, qty2, qty3 := p.QtyStock1, p.QtyStock2, p.QtyStock3
		qtyP1, qtyP2, qtyP3 := p.QtyOpname1, p.QtyOpname2, p.QtyOpname3
		productList = append(productList, entity.BulkUploadStockOpnameV2ProductItem{
			StockOpnameDetailID: p.StockOpnameDetailID,
			ProID:               p.ProID,
			ProCode:             p.ProCode,
			ProName:             p.ProName,
			UnitID1:             p.UnitID1,
			UnitID2:             p.UnitID2,
			UnitID3:             p.UnitID3,
			UnitName1:           p.UnitName1,
			UnitName2:           p.UnitName2,
			UnitName3:           p.UnitName3,
			ConvUnit1:           1,
			ConvUnit2:           p.ConvUnit2,
			ConvUnit3:           p.ConvUnit3,
			Qty1:                qty1,
			Qty2:                qty2,
			Qty3:                qty3,
			QtyPhysical1:        qtyP1,
			QtyPhysical2:        qtyP2,
			QtyPhysical3:        qtyP3,
			DifferentStock1:     qty1 - qtyP1,
			DifferentStock2:     qty2 - qtyP2,
			DifferentStock3:     qty3 - qtyP3,
			DifferentPrice1:     (qty1 - qtyP1) * p.PurchPrice1,
			DifferentPrice2:     (qty2 - qtyP2) * p.PurchPrice2,
			DifferentPrice3:     (qty3 - qtyP3) * p.PurchPrice3,
			Status:              "success",
			ErrorMessage:        nil,
		})
	}
	return []entity.BulkUploadStockOpnameV2SOData{buildBulkUploadSODataHeader(header, productList)}
}

func buildBulkUploadSODataWithRowResults(header model.StockOpnameDetailV2Header, products []model.StockOpnameDetailV2Product, detailMap map[int64]model.StockOpnameDetailV2Product, parsedRows []bulkUploadExcelRow) []entity.BulkUploadStockOpnameV2SOData {
	rowByDetailID := make(map[int64]bulkUploadExcelRow)
	for _, r := range parsedRows {
		if r.DetailID > 0 {
			rowByDetailID[r.DetailID] = r
		}
	}
	productList := make([]entity.BulkUploadStockOpnameV2ProductItem, 0, len(products))
	for _, p := range products {
		r, hasRow := rowByDetailID[p.StockOpnameDetailID]
		qty1, qty2, qty3 := p.QtyStock1, p.QtyStock2, p.QtyStock3
		var qtyP1, qtyP2, qtyP3 float64
		var status string = "success"
		var errMsg *string
		if hasRow {
			qtyP1, qtyP2, qtyP3 = r.QtyPhysical1, r.QtyPhysical2, r.QtyPhysical3
			if !r.Valid {
				status = "failed"
				errMsg = &r.ErrMsg
			}
		} else {
			qtyP1, qtyP2, qtyP3 = p.QtyOpname1, p.QtyOpname2, p.QtyOpname3
		}
		productList = append(productList, entity.BulkUploadStockOpnameV2ProductItem{
			StockOpnameDetailID: p.StockOpnameDetailID,
			ProID:               p.ProID,
			ProCode:             p.ProCode,
			ProName:             p.ProName,
			UnitID1:             p.UnitID1,
			UnitID2:             p.UnitID2,
			UnitID3:             p.UnitID3,
			UnitName1:           p.UnitName1,
			UnitName2:           p.UnitName2,
			UnitName3:           p.UnitName3,
			ConvUnit1:           1,
			ConvUnit2:           p.ConvUnit2,
			ConvUnit3:           p.ConvUnit3,
			Qty1:                qty1,
			Qty2:                qty2,
			Qty3:                qty3,
			QtyPhysical1:        qtyP1,
			QtyPhysical2:        qtyP2,
			QtyPhysical3:        qtyP3,
			DifferentStock1:     qty1 - qtyP1,
			DifferentStock2:     qty2 - qtyP2,
			DifferentStock3:     qty3 - qtyP3,
			DifferentPrice1:     (qty1 - qtyP1) * p.PurchPrice1,
			DifferentPrice2:     (qty2 - qtyP2) * p.PurchPrice2,
			DifferentPrice3:     (qty3 - qtyP3) * p.PurchPrice3,
			Status:              status,
			ErrorMessage:        errMsg,
		})
	}
	return []entity.BulkUploadStockOpnameV2SOData{buildBulkUploadSODataHeader(header, productList)}
}

func buildBulkUploadSODataHeader(header model.StockOpnameDetailV2Header, productList []entity.BulkUploadStockOpnameV2ProductItem) entity.BulkUploadStockOpnameV2SOData {
	scheduledDate := ""
	if header.ScheduledDate != nil {
		scheduledDate = *header.ScheduledDate
	}
	divisionID := int64(0)
	if header.DivisionID != nil {
		divisionID = *header.DivisionID
	}
	return entity.BulkUploadStockOpnameV2SOData{
		DocNo:         header.DocNo,
		CreatedDate:   header.CreatedDate,
		WhID:          header.WhID,
		WhCode:        header.WhCode,
		WhName:        header.WhName,
		StockType:     header.StockType,
		CreatedBy:     header.CreatedBy,
		UserName:      header.UserName,
		ScheduledDate: scheduledDate,
		Status:        header.Status,
		StatusDesc:    entity.StockOpnameStatusDesc[header.Status],
		AssignTo: entity.StockOpnameDetailV2AssignTo{
			InputBy:      header.InputBy,
			DivisionID:   divisionID,
			DivisionName: header.DivisionName,
			EmpID:        header.EmpID,
			EmpName:      header.EmpName,
		},
		ProductList: productList,
	}
}
