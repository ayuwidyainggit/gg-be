package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sales/entity"
	"sales/model"
	"sales/pkg/constant"
	"sales/pkg/str"
	"sales/pkg/structs"
	"sales/repository"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SoService interface {
	Store(request entity.CreateSoBody) (err error)
	Detail(SoNo string, custID string) (response entity.SoResponse, err error)
	List(dataFilter entity.SoQueryFilter) (data []entity.SoListResponse, total int64, lastPage int, err error)
	Delete(custId string, SoNo string, userId int64) (err error)
	Update(soNo string, request entity.UpdateSoBody) (err error)
	Download(filter entity.SoDownloadQueryFilter) (entity.ReportList, error)
}

func NewSoService(soRepository repository.SoRepository, reportRepository repository.ReportRepository, transaction repository.Dbtransaction) *soServiceImpl {
	return &soServiceImpl{
		SoRepository:     soRepository,
		ReportRepository: reportRepository,
		Transaction:      transaction,
	}
}

type soServiceImpl struct {
	SoRepository     repository.SoRepository
	ReportRepository repository.ReportRepository
	Transaction      repository.Dbtransaction
}

func (service *soServiceImpl) Store(request entity.CreateSoBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.SoDate != nil {
		SoDate, err := str.DateStrToRfc3339String(*request.SoDate)
		if err != nil {
			return err
		}
		request.SoDate = &SoDate
	}

	if request.SysDate != nil {
		SysDate, err := str.DateStrToRfc3339String(*request.SysDate)
		if err != nil {
			return err
		}
		request.SysDate = &SysDate
	}

	if request.InvoiceDate != nil {
		invoiceDate, err := str.DateStrToRfc3339String(*request.InvoiceDate)
		if err != nil {
			return err
		}
		request.InvoiceDate = &invoiceDate
	}

	if request.DeliveryDate != nil {
		deliveryDate, err := str.DateStrToRfc3339String(*request.DeliveryDate)
		if err != nil {
			return err
		}
		request.DeliveryDate = &deliveryDate
	}

	if request.DueDate != nil {
		dueDate, err := str.DateStrToRfc3339String(*request.DueDate)
		if err != nil {
			return err
		}
		request.DueDate = &dueDate
	}

	var SoModel model.So
	err = structs.Automapper(request, &SoModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.SoRepository.Store(txCtx, &SoModel)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details.Normal {
			var SoDetModel model.SoDet

			if Detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*Detail.ExpDate)
				if err != nil {
					return err
				}
				Detail.ExpDate = &expDate
			}

			err = structs.Automapper(Detail, &SoDetModel)
			if err != nil {
				return err
			}
			SoDetModel.CustID = request.CustID
			SoDetModel.SoNo = SoModel.SoNo
			SoDetModel.ItemType = 1
			err = service.SoRepository.StoreDetail(txCtx, &SoDetModel)
			if err != nil {
				return err
			}
		}
		for _, Detail := range request.Details.Promo {
			var SoDetModel model.SoDet

			if Detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*Detail.ExpDate)
				if err != nil {
					return err
				}
				Detail.ExpDate = &expDate
			}

			err = structs.Automapper(Detail, &SoDetModel)
			if err != nil {
				return err
			}
			SoDetModel.CustID = request.CustID
			SoDetModel.SoNo = SoModel.SoNo
			SoDetModel.ItemType = 2
			err = service.SoRepository.StoreDetail(txCtx, &SoDetModel)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (service *soServiceImpl) Detail(SoNo string, custID string) (response entity.SoResponse, err error) {
	so, err := service.SoRepository.FindByNo(SoNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(so, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.SoRepository.FindDetail(SoNo, custID)
	if err != nil {
		return response, err
	}
	for _, detail := range Details {
		var detailData entity.SoDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}
		if detail.ExpDate != nil {
			expDate := detail.ExpDate.Format("2006-01-02")
			detailData.ExpDate = &expDate
		}
		if detailData.ItemType == 1 {
			response.Details.Normal = append(response.Details.Normal, detailData)
		} else {
			response.Details.Promo = append(response.Details.Promo, detailData)
		}
	}
	if so.SoDate != nil {
		SoDate := so.SoDate.Format("2006-01-02")
		response.SoDate = &SoDate
	}
	if so.SysDate != nil {
		SysDate := so.SysDate.Format("2006-01-02")
		response.SysDate = &SysDate
	}
	if so.InvoiceDate != nil {
		InvDate := so.InvoiceDate.Format("2006-01-02")
		response.InvoiceDate = &InvDate

	}
	if so.DeliveryDate != nil {
		DelivDate := so.DeliveryDate.Format("2006-01-02")
		response.DeliveryDate = &DelivDate
	}

	payTypeName := response.GeneratePayTypeName()
	response.PayTypeName = payTypeName

	return response, nil
}
func (service *soServiceImpl) List(dataFilter entity.SoQueryFilter) (data []entity.SoListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.SoRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range whAdjs {
		var vResp entity.SoListResponse
		structs.Automapper(row, &vResp)
		if row.SoDate != nil {
			SoDate := row.SoDate.Format("2006-01-02")
			vResp.SoDate = &SoDate
		}

		if row.SysDate != nil {
			SysDate := row.SysDate.Format("2006-01-02")
			vResp.SysDate = &SysDate
		}

		if row.InvoiceDate != nil {
			InvDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = &InvDate

		}

		if row.DeliveryDate != nil {
			DelivDate := row.DeliveryDate.Format("2006-01-02")
			vResp.DeliveryDate = &DelivDate
		}

		payTypeName := vResp.GeneratePayTypeName()
		vResp.PayTypeName = payTypeName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *soServiceImpl) Delete(custId string, SoNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.SoRepository.Delete(txCtx, custId, SoNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (service *soServiceImpl) Update(soNo string, request entity.UpdateSoBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.SoDate != nil {
		SoDate, err := str.DateStrToRfc3339String(*request.SoDate)
		if err != nil {
			return err
		}
		request.SoDate = &SoDate
	}

	if request.SysDate != nil {
		SysDate, err := str.DateStrToRfc3339String(*request.SysDate)
		if err != nil {
			return err
		}
		request.SysDate = &SysDate
	}

	if request.InvoiceDate != nil {
		invoiceDate, err := str.DateStrToRfc3339String(*request.InvoiceDate)
		if err != nil {
			return err
		}
		request.InvoiceDate = &invoiceDate
	}

	if request.DeliveryDate != nil {
		deliveryDate, err := str.DateStrToRfc3339String(*request.DeliveryDate)
		if err != nil {
			return err
		}
		request.DeliveryDate = &deliveryDate
	}

	if request.DueDate != nil {
		dueDate, err := str.DateStrToRfc3339String(*request.DueDate)
		if err != nil {
			return err
		}
		request.DueDate = &dueDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.So
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.SoRepository.Update(txCtx, soNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details.Normal {
			if detail.SoDetID != nil {
				DetailIds = append(DetailIds, *detail.SoDetID)
			}
		}
		for _, detail := range request.Details.Promo {
			if detail.SoDetID != nil {
				DetailIds = append(DetailIds, *detail.SoDetID)
			}
		}
		if len(DetailIds) > 0 {
			err := service.SoRepository.DeleteDetailNotInIDs(txCtx, soNo, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details.Normal {
			// parse time format YYYY-mm-dd to Rfc3339

			var soDetModel model.SoDet
			if detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
				if err != nil {
					return err
				}
				detail.ExpDate = &expDate
			}

			err = structs.Automapper(detail, &soDetModel)
			if err != nil {
				return err
			}
			soDetModel.CustID = request.CustID
			soDetModel.SoNo = soNo
			soDetModel.ItemType = 1
			if detail.SoDetID == nil || *detail.SoDetID == 0 {
				soDetModel.SoDetID = nil
				err = service.SoRepository.StoreDetail(txCtx, &soDetModel)
				if err != nil {
					return err
				}
			} else {
				soDetModel.CustID = ""
				err = service.SoRepository.UpdateDetail(txCtx, &soDetModel)
				if err != nil {
					return err
				}

			}
		}

		for _, detail := range request.Details.Promo {
			// parse time format YYYY-mm-dd to Rfc3339

			var soDetModel model.SoDet
			if detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
				if err != nil {
					return err
				}
				detail.ExpDate = &expDate
			}

			err = structs.Automapper(detail, &soDetModel)
			if err != nil {
				return err
			}
			soDetModel.CustID = request.CustID
			soDetModel.SoNo = soNo
			soDetModel.ItemType = 2
			if detail.SoDetID == nil || *detail.SoDetID == 0 {
				soDetModel.SoDetID = nil
				err = service.SoRepository.StoreDetail(txCtx, &soDetModel)
				if err != nil {
					return err
				}
			} else {
				soDetModel.CustID = ""
				err = service.SoRepository.UpdateDetail(txCtx, &soDetModel)
				if err != nil {
					return err
				}

			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *soServiceImpl) Download(filter entity.SoDownloadQueryFilter) (entity.ReportList, error) {
	var data entity.ReportList

	// 1. Check if there's in-progress report
	inProgress, err := service.ReportRepository.CountDownloadSalesOrderInProgress(filter.CustId)
	if err != nil {
		log.Error("SoService, Download, CountDownloadSalesOrderInProgress:", err.Error())
		return data, err
	}
	if inProgress > 0 {
		return data, errors.New("Processing time may vary by file size. Please check Download History to access the file")
	}

	// 2. Generate report ID and name
	objectID := primitive.NewObjectID()
	objectIDString := objectID.Hex()

	loc, _ := time.LoadLocation("Asia/Jakarta")
	exportDate := time.Now().In(loc).Format("020106") // ddmmyy
	sequence := service.ReportRepository.CountDownloadSalesOrderByDate(filter.CustId, exportDate)

	reportName := entity.REPORT_NAME_DOWNLOAD_SALES_ORDER + "-" + exportDate + "-" + fmt.Sprintf("%03d", sequence)
	log.Info("SoService, Download, Generated report name:", reportName)

	// 3. Create report entry with status Processing
	reportList := model.ReportList{
		CustID:     filter.CustId,
		ReportID:   objectIDString,
		ReportName: reportName,
		StartDate:  str.UnixTimestampToUtcDate(filter.StartDate),
		EndDate:    str.UnixTimestampToUtcDate(filter.EndDate),
		FileStatus: entity.FILE_STATUS_PROCESSING,
		CreatedBy:  filter.ExportBy,
	}

	err = service.ReportRepository.StoreReportList(context.Background(), &reportList)
	if err != nil {
		log.Error("SoService, Download, StoreReportList:", err.Error())
		return data, err
	}

	// 4. Trigger async Excel generation
	filter.ReportID = objectIDString
	go service.generateDownloadSalesOrderExcel(filter)

	// 5. Return report metadata
	if err = structs.Automapper(reportList, &data); err != nil {
		log.Error("SoService, Download, Automapper:", err.Error())
		return data, err
	}
	data.StartDate = reportList.StartDate.Format(constant.YYYY_MM_DD)
	data.EndDate = reportList.EndDate.Format(constant.YYYY_MM_DD)

	return data, nil
}

// generateDownloadSalesOrderExcel generates Excel with 4 sheets and updates report with base64 content
func (service *soServiceImpl) generateDownloadSalesOrderExcel(filter entity.SoDownloadQueryFilter) {
	log.Info("SoService, generateDownloadSalesOrderExcel, Starting Excel generation for:", filter.ReportID)

	// Helper function to update report status on failure
	updateFailed := func(errMsg string) {
		log.Error("SoService, generateDownloadSalesOrderExcel,", errMsg)
		failedUpdate := model.ReportList{ReportID: filter.ReportID, FileStatus: entity.FILE_STATUS_FAILED}
		_ = service.ReportRepository.UpdateReportByReportID(context.Background(), filter.ReportID, &failedUpdate)
	}

	// Get data from repository
	dataPo, err := service.SoRepository.FindDownloadDataPo(filter)
	if err != nil {
		updateFailed("FindDownloadDataPo: " + err.Error())
		return
	}

	dataSo, err := service.SoRepository.FindDownloadDataSo(filter)
	if err != nil {
		updateFailed("FindDownloadDataSo: " + err.Error())
		return
	}

	dataFinal, err := service.SoRepository.FindDownloadDataFinal(filter)
	if err != nil {
		updateFailed("FindDownloadDataFinal: " + err.Error())
		return
	}

	dataQtySummary, err := service.SoRepository.FindDownloadQtySummary(filter)
	if err != nil {
		updateFailed("FindDownloadQtySummary: " + err.Error())
		return
	}

	// Prepare date range string for header
	startDateStr := str.UnixTimestampToUtcDate(filter.StartDate).Format("02 January 2006")
	endDateStr := str.UnixTimestampToUtcDate(filter.EndDate).Format("02 January 2006")
	dateRangeStr := startDateStr + " - " + endDateStr

	// Get salesman info for header
	salesmanInfo := ""
	if len(filter.SalesmanId) == 0 {
		salesmanInfo = "All Salesmen"
	} else if len(filter.SalesmanId) > 1 {
		salesmanInfo = "Multiple Salesmen"
	} else if len(dataPo) > 0 {
		salesmanCode := ""
		if dataPo[0].SalesmanCode != nil {
			salesmanCode = *dataPo[0].SalesmanCode
		}
		employeeName := ""
		if dataPo[0].EmployeeName != nil {
			employeeName = *dataPo[0].EmployeeName
		}
		if salesmanCode != "" && employeeName != "" {
			salesmanInfo = salesmanCode + " - " + employeeName
		} else if employeeName != "" {
			salesmanInfo = employeeName
		}
	}

	// Create Excel file with 4 sheets
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Error("SoService, generateDownloadSalesOrderExcel, Close:", err.Error())
		}
	}()

	// Sheet 1: Purchase Order
	service.createPurchaseOrderSheet(f, service.mapPoToEntity(filterDownloadDataPoWithPONumber(dataPo)), dateRangeStr, salesmanInfo)

	// Sheet 2: Sales Order
	service.createSalesOrderSheet(f, service.mapSoToEntity(dataSo), dateRangeStr, salesmanInfo)

	// Sheet 3: Final Order
	service.createFinalOrderSheet(f, service.mapFinalToEntity(dataFinal), dateRangeStr, salesmanInfo)

	// Sheet 4: QTY Summary
	service.createQtySummarySheet(f, service.mapQtySummaryToEntity(dataQtySummary), dateRangeStr, salesmanInfo)

	// Delete default Sheet1
	f.DeleteSheet("Sheet1")

	// Convert to buffer and encode as base64
	buf, err := f.WriteToBuffer()
	if err != nil {
		updateFailed("WriteToBuffer: " + err.Error())
		return
	}

	fileBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Update report with base64 content and status Ready
	updateData := model.ReportList{
		ReportID:   filter.ReportID,
		FileStatus: entity.FILE_STATUS_READY,
		FileBase64: fileBase64,
	}

	if err = service.ReportRepository.UpdateReportByReportID(context.Background(), filter.ReportID, &updateData); err != nil {
		log.Error("SoService, generateDownloadSalesOrderExcel, UpdateReportByReportID:", err.Error())
		return
	}

	log.Info("SoService, generateDownloadSalesOrderExcel, Excel generation completed for:", filter.ReportID)
}

// Helper functions to safely dereference pointer values for Excel output
func derefString(p *string) interface{} {
	if p == nil {
		return ""
	}
	return *p
}

func derefFloat64(p *float64) interface{} {
	if p == nil {
		return ""
	}
	return *p
}

func formatDownloadAmount(p *float64) string {
	if p == nil {
		return "0"
	}
	return formatDownloadAmountValue(*p)
}

func formatDownloadAmountValue(value float64) string {
	if value == 0 {
		return "0"
	}

	rounded := int64(value)
	if value > 0 && value != float64(rounded) {
		rounded = int64(value + 0.5)
	} else if value < 0 && value != float64(rounded) {
		rounded = int64(value - 0.5)
	}

	negative := rounded < 0
	if negative {
		rounded = -rounded
	}

	amount := strconv.FormatInt(rounded, 10)
	for i := len(amount) - 3; i > 0; i -= 3 {
		amount = amount[:i] + "." + amount[i:]
	}
	if negative {
		return "-" + amount
	}
	return amount
}

func resolveDownloadPONumber(poNo *string, orderNo *string) string {
	if poNo != nil && *poNo != "" {
		return *poNo
	}
	if orderNo != nil {
		return *orderNo
	}
	return ""
}

func hasValidDownloadPONumber(poNo *string, orderNo *string) bool {
	return (poNo != nil && strings.TrimSpace(*poNo) != "") || (orderNo != nil && strings.TrimSpace(*orderNo) != "")
}

func filterDownloadDataPoWithPONumber(data []model.SoDownloadPo) []model.SoDownloadPo {
	result := make([]model.SoDownloadPo, 0, len(data))
	for _, row := range data {
		if hasValidDownloadPONumber(row.PoNo, row.OrderNo) {
			result = append(result, row)
		}
	}

	return result
}

func derefFloat64Zero(p *float64) interface{} {
	if p == nil {
		return 0
	}
	return *p
}

func (service *soServiceImpl) createPurchaseOrderSheet(f *excelize.File, data []entity.SoDownloadPoRow, dateRange string, salesmanInfo string) {
	sheetName := "Purchase Order"
	f.NewSheet(sheetName)

	// Add header information rows
	f.SetCellValue(sheetName, "A1", "Order date")
	f.SetCellValue(sheetName, "B1", dateRange)
	if salesmanInfo != "" {
		f.SetCellValue(sheetName, "A2", "Salesman")
		f.SetCellValue(sheetName, "B2", salesmanInfo)
	}

	// Create style for header info
	styleID, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FCD5B4"}, Pattern: 1},
		Font: &excelize.Font{Bold: true},
	})
	f.SetCellStyle(sheetName, "A1", "A2", styleID)

	// Table headers start at row 3
	headerRow := 3
	headers := []interface{}{
		"Po No", "So No", "Order Date", "Invoice Date", "Invoice No", "Outlet Code", "Outlet Name",
		"Salesman Code", "Employee Name", "Supplier Code", "Supplier Name",
		"Product Code", "ProName", "Largest Unit", "Middle Unit", "Smallest Unit",
		"Largest Selling Price", "Middle Selling Price", "Smallest Selling Price",
		"Final Largest Selling Price", "Final Middle Selling Price", "Final Smallest Selling Price",
		"Largest QTY Order", "Middle QTY Order", "Smallest QTY Order",
		"GrossSales", "Promotion", "Discount", "Net Sales (ExcPPN)", "PPN", "Gross",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, headerRow)
		f.SetCellValue(sheetName, cell, h)
	}

	// Data rows start at row 4
	for rowIdx, row := range data {
		values := []interface{}{
			row.PoNo, row.SoNo, row.OrderDate, row.InvoiceDate, row.InvoiceNo, row.OutletCode, row.OutletName,
			derefString(row.SalesmanCode), row.EmployeeName, row.SupplierCode, row.SupplierName,
			row.ProductCode, row.ProductName, row.LargestUnit, row.MiddleUnit, row.SmallestUnit,
			formatDownloadAmount(row.LargestSellingPrice), formatDownloadAmount(row.MiddleSellingPrice), formatDownloadAmount(row.SmallestSellingPrice),
			formatDownloadAmount(row.FinalLargestSellingPrice), formatDownloadAmount(row.FinalMiddleSellingPrice), formatDownloadAmount(row.FinalSmallestSellingPrice),
			derefFloat64(row.LargestQtyOrder), derefFloat64(row.MiddleQtyOrder), derefFloat64(row.SmallestQtyOrder),
			formatDownloadAmount(row.GrossSales), formatDownloadAmount(row.Promotion), formatDownloadAmount(row.Discount), formatDownloadAmount(row.NetSales), formatDownloadAmount(row.Vat), formatDownloadAmount(row.Gross),
		}
		for colIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+headerRow+1)
			f.SetCellValue(sheetName, cell, v)
		}
	}

	f.SetColWidth(sheetName, "A", "AF", 15)
}

func (service *soServiceImpl) createSalesOrderSheet(f *excelize.File, data []entity.SoDownloadSoRow, dateRange string, salesmanInfo string) {
	sheetName := "Sales Order"
	f.NewSheet(sheetName)

	// Add header information rows
	f.SetCellValue(sheetName, "A1", "Order date")
	f.SetCellValue(sheetName, "B1", dateRange)
	if salesmanInfo != "" {
		f.SetCellValue(sheetName, "A2", "Salesman")
		f.SetCellValue(sheetName, "B2", salesmanInfo)
	}

	// Create style for header info
	styleID, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FCD5B4"}, Pattern: 1},
		Font: &excelize.Font{Bold: true},
	})
	f.SetCellStyle(sheetName, "A1", "A2", styleID)

	// Table headers start at row 3
	headerRow := 3
	headers := []interface{}{
		"Po No", "So No", "Order Date", "Invoice Date", "Invoice No", "Outlet Code", "Outlet Name",
		"Salesman Code", "Employee Name", "Supplier Code", "Supplier Name",
		"Product Code", "ProName", "Largest Unit", "Middle Unit", "Smallest Unit",
		"Largest Selling Price", "Middle Selling Price", "Smallest Selling Price",
		"Final Largest Selling Price", "Final Middle Selling Price", "Final Smallest Selling Price",
		"Largest QTY Order", "Middle QTY Order", "Smallest QTY Order",
		"GrossSales", "Promotion", "Discount", "Net Sales (ExcPPN)", "PPN", "Gross",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, headerRow)
		f.SetCellValue(sheetName, cell, h)
	}

	// Data rows start at row 4
	for rowIdx, row := range data {
		values := []interface{}{
			row.PoNo, row.SoNo, row.OrderDate, row.InvoiceDate, row.InvoiceNo, row.OutletCode, row.OutletName,
			derefString(row.SalesmanCode), row.EmployeeName, row.SupplierCode, row.SupplierName,
			row.ProductCode, row.ProductName, row.LargestUnit, row.MiddleUnit, row.SmallestUnit,
			formatDownloadAmount(row.LargestSellingPrice), formatDownloadAmount(row.MiddleSellingPrice), formatDownloadAmount(row.SmallestSellingPrice),
			formatDownloadAmount(row.FinalLargestSellingPrice), formatDownloadAmount(row.FinalMiddleSellingPrice), formatDownloadAmount(row.FinalSmallestSellingPrice),
			derefFloat64(row.LargestQtyOrder), derefFloat64(row.MiddleQtyOrder), derefFloat64(row.SmallestQtyOrder),
			formatDownloadAmount(row.GrossSales), formatDownloadAmount(row.Promotion), formatDownloadAmount(row.Discount), formatDownloadAmount(row.NetSales), formatDownloadAmount(row.Vat), formatDownloadAmount(row.Gross),
		}
		for colIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+headerRow+1)
			f.SetCellValue(sheetName, cell, v)
		}
	}

	f.SetColWidth(sheetName, "A", "AF", 15)
}

func (service *soServiceImpl) createFinalOrderSheet(f *excelize.File, data []entity.SoDownloadFinalRow, dateRange string, salesmanInfo string) {
	sheetName := "Final Order"
	f.NewSheet(sheetName)

	// Add header information rows
	f.SetCellValue(sheetName, "A1", "Order date")
	f.SetCellValue(sheetName, "B1", dateRange)
	if salesmanInfo != "" {
		f.SetCellValue(sheetName, "A2", "Salesman")
		f.SetCellValue(sheetName, "B2", salesmanInfo)
	}

	// Create style for header info
	styleID, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FCD5B4"}, Pattern: 1},
		Font: &excelize.Font{Bold: true},
	})
	f.SetCellStyle(sheetName, "A1", "A2", styleID)

	// Table headers start at row 3
	headerRow := 3
	headers := []interface{}{
		"Po No", "So No", "Order Date", "Invoice Date", "Invoice No", "Outlet Code", "Outlet Name",
		"Salesman Code", "Employee Name", "Supplier Code", "Supplier Name",
		"Product Code", "ProName", "Largest Unit", "Middle Unit", "Smallest Unit",
		"Largest Selling Price", "Middle Selling Price", "Smallest Selling Price",
		"Final Largest Selling Price", "Final Middle Selling Price", "Final Smallest Selling Price",
		"Largest QTY Order", "Middle QTY Order", "Smallest QTY Order",
		"GrossSales", "Promotion", "Discount", "Net Sales (ExcPPN)", "PPN", "Gross",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, headerRow)
		f.SetCellValue(sheetName, cell, h)
	}

	// Data rows start at row 4
	for rowIdx, row := range data {
		values := []interface{}{
			row.PoNo, row.SoNo, row.OrderDate, row.InvoiceDate, row.InvoiceNo, row.OutletCode, row.OutletName,
			derefString(row.SalesmanCode), row.EmployeeName, row.SupplierCode, row.SupplierName,
			row.ProductCode, row.ProductName, row.LargestUnit, row.MiddleUnit, row.SmallestUnit,
			formatDownloadAmount(row.LargestSellingPrice), formatDownloadAmount(row.MiddleSellingPrice), formatDownloadAmount(row.SmallestSellingPrice),
			formatDownloadAmount(row.FinalLargestSellingPrice), formatDownloadAmount(row.FinalMiddleSellingPrice), formatDownloadAmount(row.FinalSmallestSellingPrice),
			derefFloat64(row.LargestQtyOrder), derefFloat64(row.MiddleQtyOrder), derefFloat64(row.SmallestQtyOrder),
			formatDownloadAmount(row.GrossSales), formatDownloadAmount(row.Promotion), formatDownloadAmount(row.Discount), formatDownloadAmount(row.NetSales), formatDownloadAmount(row.Vat), formatDownloadAmount(row.Gross),
		}
		for colIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+headerRow+1)
			f.SetCellValue(sheetName, cell, v)
		}
	}

	f.SetColWidth(sheetName, "A", "AF", 15)
}

func (service *soServiceImpl) createQtySummarySheet(f *excelize.File, data []entity.SoDownloadQtySummaryRow, dateRange string, salesmanInfo string) {
	sheetName := "QTY Summary"
	f.NewSheet(sheetName)

	// Add header information rows
	f.SetCellValue(sheetName, "A1", "Order date")
	f.SetCellValue(sheetName, "B1", dateRange)
	if salesmanInfo != "" {
		f.SetCellValue(sheetName, "A2", "Salesman")
		f.SetCellValue(sheetName, "B2", salesmanInfo)
	}

	// Create style for header info
	styleID, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FCD5B4"}, Pattern: 1},
		Font: &excelize.Font{Bold: true},
	})
	f.SetCellStyle(sheetName, "A1", "A2", styleID)

	// Table headers start at row 3
	headerRow := 3
	headers := []interface{}{
		"Po No", "So No", "Order Date", "Invoice Date", "Invoice No", "Outlet Code", "Outlet Name",
		"Salesman Code", "Employee Name", "Supplier Code", "Supplier Name",
		"Product Code", "ProName", "Largest Unit", "Middle Unit", "Smallest Unit",
		"Largest QTY Purchase Order", "Middle QTY Purchase Order", "Smallest QTY Purchase Order",
		"Largest QTY Sales Order", "Middle QTY Sales Order", "Smallest QTY Sales Order",
		"Largest QTY Final Order", "Middle QTY Final Order", "Smallest QTY Final Order",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, headerRow)
		f.SetCellValue(sheetName, cell, h)
	}

	// Data rows start at row 4
	for rowIdx, row := range data {
		values := []interface{}{
			row.PoNo, row.SoNo, row.OrderDate, row.InvoiceDate, row.InvoiceNo, row.OutletCode, row.OutletName,
			derefString(row.SalesmanCode), row.EmployeeName, row.SupplierCode, row.SupplierName,
			row.ProductCode, row.ProductName, row.LargestUnit, row.MiddleUnit, row.SmallestUnit,
			derefFloat64Zero(row.LargestQtyPo), derefFloat64Zero(row.MiddleQtyPo), derefFloat64Zero(row.SmallestQtyPo),
			derefFloat64(row.LargestQtySo), derefFloat64(row.MiddleQtySo), derefFloat64(row.SmallestQtySo),
			derefFloat64(row.LargestQtyFinal), derefFloat64(row.MiddleQtyFinal), derefFloat64(row.SmallestQtyFinal),
		}
		for colIdx, v := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+headerRow+1)
			f.SetCellValue(sheetName, cell, v)
		}
	}

	f.SetColWidth(sheetName, "A", "Z", 15)
}

func (service *soServiceImpl) mapPoToEntity(data []model.SoDownloadPo) []entity.SoDownloadPoRow {
	var result []entity.SoDownloadPoRow

	for _, row := range data {
		orderNo := ""
		if row.OrderNo != nil {
			orderNo = *row.OrderNo
		}
		poNo := resolveDownloadPONumber(row.PoNo, row.OrderNo)
		orderDate := ""
		if row.RoDate != nil {
			orderDate = row.RoDate.Format("2006-01-02")
		}
		invoiceDate := ""
		if row.InvoiceDate != nil {
			invoiceDate = row.InvoiceDate.Format("2006-01-02")
		}
		invoiceNo := ""
		if row.InvoiceNo != nil {
			invoiceNo = *row.InvoiceNo
		}
		outletCode := ""
		if row.OutletCode != nil {
			outletCode = *row.OutletCode
		}
		outletName := ""
		if row.OutletName != nil {
			outletName = *row.OutletName
		}
		employeeName := ""
		if row.EmployeeName != nil {
			employeeName = *row.EmployeeName
		}
		supplierCode := ""
		if row.SupplierCode != nil {
			supplierCode = *row.SupplierCode
		}
		supplierName := ""
		if row.SupplierName != nil {
			supplierName = *row.SupplierName
		}
		unit3 := ""
		if row.UnitId3 != nil {
			unit3 = *row.UnitId3
		}
		unit2 := ""
		if row.UnitId2 != nil {
			unit2 = *row.UnitId2
		}
		unit1 := ""
		if row.UnitId1 != nil {
			unit1 = *row.UnitId1
		}

		// Calculate gross_sales
		var grossSales float64
		if row.QtyPo3 != nil && row.SellPricePo3 != nil {
			grossSales += *row.QtyPo3 * *row.SellPricePo3
		}
		if row.QtyPo2 != nil && row.SellPricePo2 != nil {
			grossSales += *row.QtyPo2 * *row.SellPricePo2
		}
		if row.QtyPo1 != nil && row.SellPricePo1 != nil {
			grossSales += *row.QtyPo1 * *row.SellPricePo1
		}

		// Calculate net_sales
		promotion := 0.0
		discount := 0.0
		if row.DiscValueFinal != nil {
			discount = *row.DiscValueFinal
		}
		netSales := grossSales - promotion - discount

		vat := 0.0
		if row.VatValueFinal != nil {
			vat = *row.VatValueFinal
		}
		gross := grossSales

		result = append(result, entity.SoDownloadPoRow{
			OrderNo:                   orderNo,
			PoNo:                      poNo,
			SoNo:                      row.SoNo,
			OrderDate:                 orderDate,
			InvoiceDate:               invoiceDate,
			InvoiceNo:                 invoiceNo,
			OutletCode:                outletCode,
			OutletName:                outletName,
			SalesmanCode:              row.SalesmanCode,
			EmployeeName:              employeeName,
			SupplierCode:              supplierCode,
			SupplierName:              supplierName,
			ProductCode:               row.ProductCode,
			ProductName:               row.ProductName,
			LargestUnit:               unit3,
			MiddleUnit:                unit2,
			SmallestUnit:              unit1,
			LargestSellingPrice:       row.SellPriceSystem3,
			MiddleSellingPrice:        row.SellPriceSystem2,
			SmallestSellingPrice:      row.SellPriceSystem1,
			FinalLargestSellingPrice:  row.SellPricePo3,
			FinalMiddleSellingPrice:   row.SellPricePo2,
			FinalSmallestSellingPrice: row.SellPricePo1,
			LargestQtyOrder:           row.QtyPo3,
			MiddleQtyOrder:            row.QtyPo2,
			SmallestQtyOrder:          row.QtyPo1,
			GrossSales:                &grossSales,
			Promotion:                 &promotion,
			Discount:                  &discount,
			NetSales:                  &netSales,
			Vat:                       &vat,
			Gross:                     &gross,
		})
	}

	return result
}

func (service *soServiceImpl) mapSoToEntity(data []model.SoDownloadSo) []entity.SoDownloadSoRow {
	var result []entity.SoDownloadSoRow

	for _, row := range data {
		orderNo := ""
		if row.OrderNo != nil {
			orderNo = *row.OrderNo
		}
		poNo := resolveDownloadPONumber(row.PoNo, row.OrderNo)
		orderDate := ""
		if row.RoDate != nil {
			orderDate = row.RoDate.Format("2006-01-02")
		}
		invoiceDate := ""
		if row.InvoiceDate != nil {
			invoiceDate = row.InvoiceDate.Format("2006-01-02")
		}
		invoiceNo := ""
		if row.InvoiceNo != nil {
			invoiceNo = *row.InvoiceNo
		}
		outletCode := ""
		if row.OutletCode != nil {
			outletCode = *row.OutletCode
		}
		outletName := ""
		if row.OutletName != nil {
			outletName = *row.OutletName
		}
		employeeName := ""
		if row.EmployeeName != nil {
			employeeName = *row.EmployeeName
		}
		supplierCode := ""
		if row.SupplierCode != nil {
			supplierCode = *row.SupplierCode
		}
		supplierName := ""
		if row.SupplierName != nil {
			supplierName = *row.SupplierName
		}
		unit3 := ""
		if row.UnitId3 != nil {
			unit3 = *row.UnitId3
		}
		unit2 := ""
		if row.UnitId2 != nil {
			unit2 = *row.UnitId2
		}
		unit1 := ""
		if row.UnitId1 != nil {
			unit1 = *row.UnitId1
		}

		// Calculate gross_sales
		var grossSales float64
		if row.Qty3 != nil && row.SellPrice3 != nil {
			grossSales += *row.Qty3 * *row.SellPrice3
		}
		if row.Qty2 != nil && row.SellPrice2 != nil {
			grossSales += *row.Qty2 * *row.SellPrice2
		}
		if row.Qty1 != nil && row.SellPrice1 != nil {
			grossSales += *row.Qty1 * *row.SellPrice1
		}

		// Calculate net_sales
		promotion := 0.0
		discount := 0.0
		if row.DiscValueFinal != nil {
			discount = *row.DiscValueFinal
		}
		netSales := grossSales - promotion - discount

		vat := 0.0
		if row.VatValueFinal != nil {
			vat = *row.VatValueFinal
		}
		gross := grossSales

		result = append(result, entity.SoDownloadSoRow{
			OrderNo:                   orderNo,
			PoNo:                      poNo,
			SoNo:                      row.SoNo,
			OrderDate:                 orderDate,
			InvoiceDate:               invoiceDate,
			InvoiceNo:                 invoiceNo,
			OutletCode:                outletCode,
			OutletName:                outletName,
			SalesmanCode:              row.SalesmanCode,
			EmployeeName:              employeeName,
			SupplierCode:              supplierCode,
			SupplierName:              supplierName,
			ProductCode:               row.ProductCode,
			ProductName:               row.ProductName,
			LargestUnit:               unit3,
			MiddleUnit:                unit2,
			SmallestUnit:              unit1,
			LargestSellingPrice:       row.SellPriceSystem3,
			MiddleSellingPrice:        row.SellPriceSystem2,
			SmallestSellingPrice:      row.SellPriceSystem1,
			FinalLargestSellingPrice:  row.SellPrice3,
			FinalMiddleSellingPrice:   row.SellPrice2,
			FinalSmallestSellingPrice: row.SellPrice1,
			LargestQtyOrder:           row.Qty3,
			MiddleQtyOrder:            row.Qty2,
			SmallestQtyOrder:          row.Qty1,
			GrossSales:                &grossSales,
			Promotion:                 &promotion,
			Discount:                  &discount,
			NetSales:                  &netSales,
			Vat:                       &vat,
			Gross:                     &gross,
		})
	}

	return result
}

func (service *soServiceImpl) mapFinalToEntity(data []model.SoDownloadFinal) []entity.SoDownloadFinalRow {
	var result []entity.SoDownloadFinalRow

	for _, row := range data {
		orderNo := ""
		if row.OrderNo != nil {
			orderNo = *row.OrderNo
		}
		poNo := resolveDownloadPONumber(row.PoNo, row.OrderNo)
		orderDate := ""
		if row.RoDate != nil {
			orderDate = row.RoDate.Format("2006-01-02")
		}
		invoiceDate := ""
		if row.InvoiceDate != nil {
			invoiceDate = row.InvoiceDate.Format("2006-01-02")
		}
		invoiceNo := ""
		if row.InvoiceNo != nil {
			invoiceNo = *row.InvoiceNo
		}
		outletCode := ""
		if row.OutletCode != nil {
			outletCode = *row.OutletCode
		}
		outletName := ""
		if row.OutletName != nil {
			outletName = *row.OutletName
		}
		employeeName := ""
		if row.EmployeeName != nil {
			employeeName = *row.EmployeeName
		}
		supplierCode := ""
		if row.SupplierCode != nil {
			supplierCode = *row.SupplierCode
		}
		supplierName := ""
		if row.SupplierName != nil {
			supplierName = *row.SupplierName
		}
		unit3 := ""
		if row.UnitId3 != nil {
			unit3 = *row.UnitId3
		}
		unit2 := ""
		if row.UnitId2 != nil {
			unit2 = *row.UnitId2
		}
		unit1 := ""
		if row.UnitId1 != nil {
			unit1 = *row.UnitId1
		}

		// Calculate gross_sales
		var grossSales float64
		if row.Qty3Final != nil && row.SellPriceFinal3 != nil {
			grossSales += *row.Qty3Final * *row.SellPriceFinal3
		}
		if row.Qty2Final != nil && row.SellPriceFinal2 != nil {
			grossSales += *row.Qty2Final * *row.SellPriceFinal2
		}
		if row.Qty1Final != nil && row.SellPriceFinal1 != nil {
			grossSales += *row.Qty1Final * *row.SellPriceFinal1
		}

		// Calculate net_sales
		promotion := 0.0
		discount := 0.0
		if row.DiscValueFinal != nil {
			discount = *row.DiscValueFinal
		}
		netSales := grossSales - promotion - discount

		vat := 0.0
		if row.VatValueFinal != nil {
			vat = *row.VatValueFinal
		}
		gross := grossSales

		result = append(result, entity.SoDownloadFinalRow{
			OrderNo:                   orderNo,
			PoNo:                      poNo,
			SoNo:                      row.SoNo,
			OrderDate:                 orderDate,
			InvoiceDate:               invoiceDate,
			InvoiceNo:                 invoiceNo,
			OutletCode:                outletCode,
			OutletName:                outletName,
			SalesmanCode:              row.SalesmanCode,
			EmployeeName:              employeeName,
			SupplierCode:              supplierCode,
			SupplierName:              supplierName,
			ProductCode:               row.ProductCode,
			ProductName:               row.ProductName,
			LargestUnit:               unit3,
			MiddleUnit:                unit2,
			SmallestUnit:              unit1,
			LargestSellingPrice:       row.SellPriceSystem3,
			MiddleSellingPrice:        row.SellPriceSystem2,
			SmallestSellingPrice:      row.SellPriceSystem1,
			FinalLargestSellingPrice:  row.SellPriceFinal3,
			FinalMiddleSellingPrice:   row.SellPriceFinal2,
			FinalSmallestSellingPrice: row.SellPriceFinal1,
			LargestQtyOrder:           row.Qty3Final,
			MiddleQtyOrder:            row.Qty2Final,
			SmallestQtyOrder:          row.Qty1Final,
			GrossSales:                &grossSales,
			Promotion:                 &promotion,
			Discount:                  &discount,
			NetSales:                  &netSales,
			Vat:                       &vat,
			Gross:                     &gross,
		})
	}

	return result
}

func (service *soServiceImpl) mapQtySummaryToEntity(data []model.SoDownloadQtySummary) []entity.SoDownloadQtySummaryRow {
	var result []entity.SoDownloadQtySummaryRow

	for _, row := range data {
		orderNo := ""
		if row.OrderNo != nil {
			orderNo = *row.OrderNo
		}
		poNo := resolveDownloadPONumber(row.PoNo, row.OrderNo)
		orderDate := ""
		if row.RoDate != nil {
			orderDate = row.RoDate.Format("2006-01-02")
		}
		invoiceDate := ""
		if row.InvoiceDate != nil {
			invoiceDate = row.InvoiceDate.Format("2006-01-02")
		}
		invoiceNo := ""
		if row.InvoiceNo != nil {
			invoiceNo = *row.InvoiceNo
		}
		outletCode := ""
		if row.OutletCode != nil {
			outletCode = *row.OutletCode
		}
		outletName := ""
		if row.OutletName != nil {
			outletName = *row.OutletName
		}
		employeeName := ""
		if row.EmployeeName != nil {
			employeeName = *row.EmployeeName
		}
		supplierCode := ""
		if row.SupplierCode != nil {
			supplierCode = *row.SupplierCode
		}
		supplierName := ""
		if row.SupplierName != nil {
			supplierName = *row.SupplierName
		}
		unit3 := ""
		if row.UnitId3 != nil {
			unit3 = *row.UnitId3
		}
		unit2 := ""
		if row.UnitId2 != nil {
			unit2 = *row.UnitId2
		}
		unit1 := ""
		if row.UnitId1 != nil {
			unit1 = *row.UnitId1
		}

		result = append(result, entity.SoDownloadQtySummaryRow{
			OrderNo:          orderNo,
			PoNo:             poNo,
			SoNo:             row.SoNo,
			OrderDate:        orderDate,
			InvoiceDate:      invoiceDate,
			InvoiceNo:        invoiceNo,
			OutletCode:       outletCode,
			OutletName:       outletName,
			SalesmanCode:     row.SalesmanCode,
			EmployeeName:     employeeName,
			SupplierCode:     supplierCode,
			SupplierName:     supplierName,
			ProductCode:      row.ProductCode,
			ProductName:      row.ProductName,
			LargestUnit:      unit3,
			MiddleUnit:       unit2,
			SmallestUnit:     unit1,
			LargestQtyPo:     row.QtyPo3,
			MiddleQtyPo:      row.QtyPo2,
			SmallestQtyPo:    row.QtyPo1,
			LargestQtySo:     row.Qty3,
			MiddleQtySo:      row.Qty2,
			SmallestQtySo:    row.Qty1,
			LargestQtyFinal:  row.Qty3Final,
			MiddleQtyFinal:   row.Qty2Final,
			SmallestQtyFinal: row.Qty1Final,
		})
	}

	return result
}
