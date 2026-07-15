package service

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"inventory/adapter"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/constant"
	"inventory/pkg/conversion"
	"inventory/pkg/structs"
	"inventory/repository"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/xuri/excelize/v2"
)

type ReportsService interface {
	GetStockMovementReport(dataFilter entity.StockMovementReportQueryFilter) (*entity.StockMovementReportResponse, error)
	PreviewDownloadStockMovementReport(dataFilter entity.PreviewDownloadStockMovementReportQueryFilter) (res []*entity.PreviewDownloadStockMovementReportResponse, total int64, err error)
	DownloadStockMovementReport(dataFilter entity.DownloadStockMovementReportQueryFilter) (*entity.DownloadStockMovementReportResponse, error)
}

func NewReportsService(reportsRepository repository.ReportsRepository, obsAdapter *adapter.ObsAdapterImpl) *reportsServiceImpl {
	return &reportsServiceImpl{
		ReportsRepository: reportsRepository,
		ObsAdapter:        obsAdapter,
	}
}

type reportsServiceImpl struct {
	ReportsRepository repository.ReportsRepository
	ObsAdapter        *adapter.ObsAdapterImpl
}

func (service *reportsServiceImpl) GetStockMovementReport(dataFilter entity.StockMovementReportQueryFilter) (*entity.StockMovementReportResponse, error) {
	ctx := context.Background()

	now := time.Now()
	month := dataFilter.Month
	year := dataFilter.Year
	if month == 0 {
		month = int(now.Month())
	}
	if year == 0 {
		year = now.Year()
	}

	whTotalStock, err := service.ReportsRepository.GetStockMovementWarehouseTotalStock(ctx, dataFilter.CustID, month, year)
	if err != nil {
		return nil, err
	}

	var whTotalStockEntities []entity.StockMovementWarehouseTotalStock
	for _, wh := range whTotalStock {
		whTotalStockEntities = append(whTotalStockEntities, entity.StockMovementWarehouseTotalStock{
			WhID:          wh.WhID,
			WhCode:        wh.WhCode,
			WhName:        wh.WhName,
			OpeningStock:  wh.OpeningStock,
			ChangingStock: wh.ChangingStock,
			ClosingStock:  wh.ClosingStock,
		})
	}

	stockMovements, err := service.ReportsRepository.GetStockMovementTransactionTypes(ctx, dataFilter.CustID, month, year)
	if err != nil {
		return nil, err
	}

	var stockMovementEntities []entity.StockMovementTransactionType
	for _, sm := range stockMovements {
		if sm.TrCode == "CO" {
			continue // Skip "CO" transaction type as per spec
		}
		stockMovementEntities = append(stockMovementEntities, entity.StockMovementTransactionType{
			TrCode:  sm.TrCode,
			TrName:  sm.TrName,
			NoOfDoc: sm.NoOfDoc,
		})
	}

	topProductsIn, err := service.ReportsRepository.GetTopProductsByStockIn(ctx, dataFilter.CustID, month, year)
	if err != nil {
		return nil, err
	}

	topProductsInEntities := make([]entity.StockMovementTopProduct, 0)
	for _, product := range topProductsIn {
		convUnit2 := int(product.ConvUnit2)
		convUnit3 := int(product.ConvUnit3)
		if convUnit2 == 0 {
			convUnit2 = 1
		}
		if convUnit3 == 0 {
			convUnit3 = 1
		}

		qty := &conversion.Qty{
			Qty:       int(product.TotalQty),
			ConvUnit2: convUnit2,
			ConvUnit3: convUnit3,
		}
		qtyConversion := qty.ConvToQtyConversion()

		topProductsInEntities = append(topProductsInEntities, entity.StockMovementTopProduct{
			ProName:     product.ProName,
			QtyLargest:  int64(qtyConversion.Qty3), // Largest = Qty3
			QtyMedium:   int64(qtyConversion.Qty2), // Medium = Qty2
			QtySmallest: int64(qtyConversion.Qty1), // Smallest = Qty1
		})
	}

	topProductsOut, err := service.ReportsRepository.GetTopProductsByStockOut(ctx, dataFilter.CustID, month, year)
	if err != nil {
		return nil, err
	}

	topProductsOutEntities := make([]entity.StockMovementTopProduct, 0)
	for _, product := range topProductsOut {
		convUnit2 := int(product.ConvUnit2)
		convUnit3 := int(product.ConvUnit3)
		if convUnit2 == 0 {
			convUnit2 = 1
		}
		if convUnit3 == 0 {
			convUnit3 = 1
		}

		qty := &conversion.Qty{
			Qty:       int(product.TotalQty),
			ConvUnit2: convUnit2,
			ConvUnit3: convUnit3,
		}
		qtyConversion := qty.ConvToQtyConversion()

		topProductsOutEntities = append(topProductsOutEntities, entity.StockMovementTopProduct{
			ProName:     product.ProName,
			QtyLargest:  int64(qtyConversion.Qty3), // Largest = Qty3
			QtyMedium:   int64(qtyConversion.Qty2), // Medium = Qty2
			QtySmallest: int64(qtyConversion.Qty1), // Smallest = Qty1
		})
	}

	stockAkhir, err := service.ReportsRepository.GetNetStockChangesCurrentMonth(ctx, dataFilter.CustID, month, year)
	if err != nil {
		return nil, err
	}

	stockAwal, err := service.ReportsRepository.GetNetStockChangesPreviousMonth(ctx, dataFilter.CustID, month, year)
	if err != nil {
		return nil, err
	}

	growStock := stockAkhir - stockAwal

	response := &entity.StockMovementReportResponse{
		WhTotalStock: whTotalStockEntities,
		NetStockChanges: entity.StockMovementNetStockChanges{
			StockAwal:  stockAwal,
			StockAkhir: stockAkhir,
			GrowStock:  growStock,
		},
		StockMovement: stockMovementEntities,
		TopProductIn:  topProductsInEntities,
		TopProductOut: topProductsOutEntities,
	}

	return response, nil
}

func (service *reportsServiceImpl) PreviewDownloadStockMovementReport(dataFilter entity.PreviewDownloadStockMovementReportQueryFilter) (res []*entity.PreviewDownloadStockMovementReportResponse, total int64, err error) {
	ctx := context.Background()

	// Map query filter into model for repository
	var params model.StockLedgerRequest
	if err = structs.Automapper(dataFilter, &params); err != nil {
		return nil, 0, err
	}

	// Explicitly takeout TransactionType equals to "CO".
	filtered := make([]string, 0, len(params.TransactionTypes))
	for _, v := range params.TransactionTypes {
		if v != "CO" {
			filtered = append(filtered, v)
		}
	}
	params.TransactionTypes = filtered

	// Fetch stock ledger rows
	stockLedger, err := service.ReportsRepository.GetStockLedger(ctx, &params)
	if err != nil {
		log.Error("DownloadStockMovementReport GetStockLedger error: ", err)
		return nil, 0, err
	}

	if len(stockLedger) == 0 {
		return nil, 0, fmt.Errorf("no stock movement data found for given filter")
	}
	total = stockLedger[0].TotalRecord

	// // Map transaction type codes to names for better readability in the report
	// for i := range stockLedger {
	// 	if trxType, ok := constant.MapTransactionType[stockLedger[i].TransactionType]; ok {
	// 		stockLedger[i].TransactionType = trxType
	// 	}
	// }

	// Map model rows to entity response slice
	results := make([]*entity.PreviewDownloadStockMovementReportResponse, 0, len(stockLedger))
	for _, r := range stockLedger {
		remarks := ""
		if r.Updates >= 1 {
			remarks = "Tambah Barang"
		} else if r.Updates <= 1 {
			remarks = "Pengurangan Barang"
		}

		row := &entity.PreviewDownloadStockMovementReportResponse{
			DistributorID:   r.DistributorID,
			DistributorCode: r.DistributorCode,
			DistributorName: r.DistributorName,

			WhID:   r.WhID,
			WhCode: r.WhCode,
			WhName: r.Warehouse,
			Date:   r.Date,

			ProID:   r.ProID,
			ProCode: r.ProductCode,
			ProName: r.ProductName,

			OpeningStock1: r.OpeningStockLarge,
			OpeningStock2: r.OpeningStockMedium,
			OpeningStock3: r.OpeningStockSmall,

			ChangesStock1: r.UpdatesLarge,
			ChangesStock2: r.UpdatesMedium,
			ChangesStock3: r.UpdatesSmall,

			ClosingStock1: r.ClosingStockLarge,
			ClosingStock2: r.ClosingStockMedium,
			ClosingStock3: r.ClosingStockSmall,

			TransactionType: r.TransactionType,
			RefNo:           r.ReferenceNo,
			Remarks:         remarks,
		}
		results = append(results, row)
	}

	return results, total, nil
}

func (service *reportsServiceImpl) DownloadStockMovementReport(dataFilter entity.DownloadStockMovementReportQueryFilter) (res *entity.DownloadStockMovementReportResponse, err error) {
	ctx := context.Background()

	// Map query filter into model for repository
	var params model.StockLedgerRequest
	if err = structs.Automapper(dataFilter, &params); err != nil {
		return nil, err
	}

	// Explicitly takeout TransactionType equals to "CO".
	filtered := make([]string, 0, len(params.TransactionTypes))
	for _, v := range params.TransactionTypes {
		if v != "CO" {
			filtered = append(filtered, v)
		}
	}
	params.TransactionTypes = filtered

	// Fetch stock ledger rows
	stockLedger, err := service.ReportsRepository.GetStockLedger(ctx, &params)
	if err != nil {
		log.Error("DownloadStockMovementReport GetStockLedger error: ", err)
		return nil, err
	}

	if len(stockLedger) == 0 {
		return nil, fmt.Errorf("no stock movement data found for given filter")
	}

	// Map transaction type codes to names for better readability in the report
	for i := range stockLedger {
		if trxType, ok := constant.MapTransactionType[stockLedger[i].TransactionType]; ok {
			stockLedger[i].TransactionType = trxType
		}
	}

	// Generate Excel file based on stock ledger rows
	excelFile, err := generateStockMovementExcel(stockLedger)
	if err != nil {
		log.Error("DownloadStockMovementReport generateStockMovementExcel error: ", err)
		return nil, err
	}
	defer func() {
		if cerr := excelFile.Close(); cerr != nil {
			log.Error("DownloadStockMovementReport excel close error: ", cerr)
		}
	}()

	var buf bytes.Buffer
	if err = excelFile.Write(&buf); err != nil {
		log.Error("DownloadStockMovementReport excel write error: ", err)
		return nil, err
	}

	// // Ensure tmp directory exists
	// tmpDir := "./tmp"
	// if err := os.MkdirAll(tmpDir, 0755); err != nil {
	// 	log.Error("failed to create tmp dir: ", err)
	// 	return nil, err
	// }

	// // Build filename for preview
	// previewFilename := fmt.Sprintf("%s.xlsx", "reportName")
	// previewPath := filepath.Join(tmpDir, previewFilename)

	// // Save file for preview
	// if err := os.WriteFile(previewPath, buf.Bytes(), 0644); err != nil {
	// 	log.Error("failed to write preview excel file: ", err)
	// 	return nil, err
	// }

	fileBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Build report metadata
	now := time.Now()

	startDate := now
	endDate := now

	if dataFilter.StartDate != "" {
		if t, parseErr := time.Parse("2006-01-02", dataFilter.StartDate); parseErr == nil {
			startDate = t
		} else {
			return nil, parseErr
		}
	}

	if dataFilter.EndDate != "" {
		if t, parseErr := time.Parse("2006-01-02", dataFilter.EndDate); parseErr == nil {
			endDate = t
		} else {
			return nil, parseErr
		}
	}

	// Generate report name: DownloadStockMovement-DDMMYY-3digitRunningNumber
	dateStr := now.Format("020106") // DDMMYY
	sequenceNumber, err := getNextSequenceNumber(dateStr)
	if err != nil {
		log.Error("DownloadStockMovementReport getNextSequenceNumber error: ", err)
		return nil, fmt.Errorf("failed to get sequence number: %w", err)
	}
	reportName := fmt.Sprintf("DownloadStockMovement-%s-%03d", dateStr, sequenceNumber)

	// upload the buffer to OBS if configured
	var fileURL string
	if service.ObsAdapter != nil {
		uploadReq := &model.UploadBytes{
			Folder:      "reports",
			FileName:    fmt.Sprintf("%s.xlsx", reportName),
			Data:        buf.Bytes(),
			ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		}
		url, upErr := service.ObsAdapter.UploadBytes(uploadReq)
		if upErr != nil {
			log.Error("DownloadStockMovementReport upload to OBS error: ", upErr)
			// continue even if upload fails; we'll still return base64
		} else {
			fileURL = url
		}
	}

	// Generate report ID (24 hex characters similar to ObjectID)
	reportID := generateReportID()

	// First insert to report.list with status 0 and empty file content
	reportRow := &model.ReportList{
		CustID:     dataFilter.CustID,
		ReportID:   reportID,
		ReportName: reportName,
		StartDate:  startDate,
		EndDate:    endDate,
		FileStatus: 0, // 0 = processing
		FileURL:    fileURL,
		FileBase64: fileBase64,
		CreatedBy:  dataFilter.UserFullName,
		CreatedAt:  now,
	}

	if err = service.ReportsRepository.CreateReportList(ctx, reportRow); err != nil {
		return nil, err
	}

	// Update report.list with final file content & completed status (1 as per spec)
	const fileStatusCompleted = 1
	if err = service.ReportsRepository.UpdateReportListFile(ctx, reportID, fileStatusCompleted, fileBase64, fileURL); err != nil {
		return nil, err
	}

	reportRow.FileStatus = fileStatusCompleted
	reportRow.FileBase64 = fileBase64

	// Build response payload
	res = &entity.DownloadStockMovementReportResponse{
		ReportID:   reportRow.ReportID,
		ReportName: reportRow.ReportName,
		StartDate:  reportRow.StartDate.Format("2006-01-02"),
		EndDate:    reportRow.EndDate.Format("2006-01-02"),
		FileStatus: reportRow.FileStatus,
		// FileStatusName: "Completed",
		FileURL:   reportRow.FileURL,
		CreatedBy: reportRow.CreatedBy,
		CreatedAt: reportRow.CreatedAt,
	}

	return res, nil
}

// generateReportID returns a pseudo ObjectID-like string (24 hex chars).
func generateReportID() string {
	b := make([]byte, 12) // 24 hex characters
	if _, err := crand.Read(b); err != nil {
		// Fallback to timestamp-based ID if random fails
		ts := time.Now().UnixNano()
		return fmt.Sprintf("%024x", ts)
	}
	return hex.EncodeToString(b)
}

// generateStockMovementExcel builds the Stock Movement Excel file layout from ledger rows.
func generateStockMovementExcel(rows []*model.StockLedgerRow) (*excelize.File, error) {
	f := excelize.NewFile()

	sheetName := "Stock Movement"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")
	f.SetDefaultFont("Poppins")

	// Styles
	standardFont := &excelize.Font{
		Family: "Poppins",
		Bold:   false,
		Color:  "#353535",
		Size:   11,
	}

	groupHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font: standardFont,
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			Indent:     1,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#F2F6FE"},
			Pattern: 1,
		},
		// Border: []excelize.Border{
		// 	{Type: "top", Style: 0},
		// 	{Type: "right", Style: 0},
		// 	{Type: "bottom", Style: 0},
		// 	{Type: "left", Style: 0},
		// },
	})

	subHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font: standardFont,
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#F2F6FE"},
			Pattern: 1,
		},
		// Border: []excelize.Border{
		// 	{Type: "top", Style: 0},
		// 	{Type: "right", Style: 0},
		// 	{Type: "bottom", Style: 0},
		// 	{Type: "left", Style: 0},
		// },
	})

	// altRowStyle, _ := f.NewStyle(&excelize.Style{
	// 	Font: standardFont,
	// 	Fill: excelize.Fill{
	// 		Type:    "pattern",
	// 		Color:   []string{"#FFFFFF", "#F2F6FE"},
	// 		Pattern: 5,
	// 	},
	// 	Border: []excelize.Border{
	// 		{Type: "top", Style: 0},
	// 		{Type: "right", Style: 1},
	// 		{Type: "bottom", Style: 1},
	// 		{Type: "left", Style: 0},
	// 	},
	// })

	// headerStyle, err := f.NewStyle(&excelize.Style{
	// 	Font:      &excelize.Font{Bold: true},
	// 	Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	// 	Border: []excelize.Border{
	// 		{Type: "left", Color: "000000", Style: 1},
	// 		{Type: "top", Color: "000000", Style: 1},
	// 		{Type: "bottom", Color: "000000", Style: 1},
	// 		{Type: "right", Color: "000000", Style: 1},
	// 	},
	// })
	// if err != nil {
	// 	return nil, err
	// }

	textCellStyle, err := f.NewStyle(&excelize.Style{
		Font:      standardFont,
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
		// Fill: excelize.Fill{
		// 	Type:    "pattern",
		// 	Color:   []string{"#FFFFFF", "#F2F6FE"},
		// 	Pattern: 5,
		// },
		// Border: []excelize.Border{
		// 	{Type: "top", Style: 0},
		// 	{Type: "right", Style: 0},
		// 	{Type: "bottom", Style: 0},
		// 	{Type: "left", Style: 0},
		// },
	})
	if err != nil {
		return nil, err
	}

	numberCellStyle, _ := f.NewStyle(&excelize.Style{
		Font: standardFont,
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		// Fill: excelize.Fill{
		// 	Type:    "pattern",
		// 	Color:   []string{"#FFFFFF", "#F2F6FE"},
		// 	Pattern: 5,
		// },
		// Border: []excelize.Border{
		// 	{Type: "top", Style: 0},
		// 	{Type: "right", Style: 0},
		// 	{Type: "bottom", Style: 0},
		// 	{Type: "left", Style: 0},
		// },
		NumFmt: 3, // #,##0
	})

	dateCellStyle, err := f.NewStyle(&excelize.Style{
		Font:      standardFont,
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		// Fill: excelize.Fill{
		// 	Type:    "pattern",
		// 	Color:   []string{"#F2F6FE"},
		// 	Pattern: 5,
		// },
		// Border: []excelize.Border{
		// 	{Type: "top", Style: 0},
		// 	{Type: "right", Style: 0},
		// 	{Type: "bottom", Style: 0},
		// 	{Type: "left", Style: 0},
		// },
		NumFmt: 14, // dd-mm-yy
	})
	if err != nil {
		return nil, err
	}

	// Static headers
	staticHeaders := map[string]string{
		"A": "Distributor Code",
		"B": "Distributor Name",
		"C": "Warehouse",
		"D": "Date",
		"E": "Product Code",
		"F": "Product Name",
		"P": "Transaction Type",
		"Q": "Reference No",
	}

	for col, text := range staticHeaders {
		cell := col + "1"
		f.SetCellValue(sheetName, cell, text)
		f.SetCellStyle(sheetName, cell, cell, groupHeaderStyle)
		f.MergeCell(sheetName, col+"1", col+"2")
	}

	// Group headers
	f.SetCellValue(sheetName, "G1", "Opening Stock")
	f.MergeCell(sheetName, "G1", "I1")
	f.SetCellStyle(sheetName, "G1", "I1", groupHeaderStyle)

	f.SetCellValue(sheetName, "J1", "Change (+ / -)")
	f.MergeCell(sheetName, "J1", "L1")
	f.SetCellStyle(sheetName, "J1", "L1", groupHeaderStyle)

	f.SetCellValue(sheetName, "M1", "Closing Stock")
	f.MergeCell(sheetName, "M1", "O1")
	f.SetCellStyle(sheetName, "M1", "O1", groupHeaderStyle)

	subHeaders := map[string]string{
		"G": "Largest", "H": "Medium", "I": "Smallest",
		"J": "Largest", "K": "Medium", "L": "Smallest",
		"M": "Largest", "N": "Medium", "O": "Smallest",
	}

	for col, text := range subHeaders {
		cell := col + "2"
		f.SetCellValue(sheetName, cell, text)
		f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)
	}

	// Data rows
	for idx, row := range rows {
		r := idx + 3

		// if r%2 == 0 {
		// 	f.SetRowStyle(sheetName, r, r, altRowStyle)
		// }

		// Distributor & warehouse info
		if err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", r), row.DistributorCode); err != nil {
			return nil, err
		}
		if err := f.SetCellStyle(sheetName, fmt.Sprintf("A%d", r), fmt.Sprintf("A%d", r), textCellStyle); err != nil {
			return nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", r), row.DistributorName); err != nil {
			return nil, err
		}
		if err := f.SetCellStyle(sheetName, fmt.Sprintf("B%d", r), fmt.Sprintf("B%d", r), textCellStyle); err != nil {
			return nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("C%d", r), row.Warehouse); err != nil {
			return nil, err
		}
		if err := f.SetCellStyle(sheetName, fmt.Sprintf("C%d", r), fmt.Sprintf("C%d", r), textCellStyle); err != nil {
			return nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("D%d", r), row.Date); err != nil {
			return nil, err
		}
		if err := f.SetCellStyle(sheetName, fmt.Sprintf("D%d", r), fmt.Sprintf("D%d", r), dateCellStyle); err != nil {
			return nil, err
		}

		// Product info
		if err := f.SetCellValue(sheetName, fmt.Sprintf("E%d", r), row.ProductCode); err != nil {
			return nil, err
		}
		if err := f.SetCellStyle(sheetName, fmt.Sprintf("E%d", r), fmt.Sprintf("E%d", r), textCellStyle); err != nil {
			return nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("F%d", r), row.ProductName); err != nil {
			return nil, err
		}
		if err := f.SetCellStyle(sheetName, fmt.Sprintf("F%d", r), fmt.Sprintf("F%d", r), textCellStyle); err != nil {
			return nil, err
		}

		// Opening stock
		openCols := []struct {
			col string
			val int64
		}{
			{"G", row.OpeningStockLarge},
			{"H", row.OpeningStockMedium},
			{"I", row.OpeningStockSmall},
		}
		for _, c := range openCols {
			if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", c.col, r), c.val); err != nil {
				return nil, err
			}
			if err := f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", c.col, r), fmt.Sprintf("%s%d", c.col, r), numberCellStyle); err != nil {
				return nil, err
			}
		}

		// Change qty
		changeCols := []struct {
			col string
			val int64
		}{
			{"J", row.UpdatesLarge},
			{"K", row.UpdatesMedium},
			{"L", row.UpdatesSmall},
		}
		for _, c := range changeCols {
			if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", c.col, r), c.val); err != nil {
				return nil, err
			}
			if err := f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", c.col, r), fmt.Sprintf("%s%d", c.col, r), numberCellStyle); err != nil {
				return nil, err
			}
		}

		// Closing stock
		closeCols := []struct {
			col string
			val int64
		}{
			{"M", row.ClosingStockLarge},
			{"N", row.ClosingStockMedium},
			{"O", row.ClosingStockSmall},
		}
		for _, c := range closeCols {
			if err := f.SetCellValue(sheetName, fmt.Sprintf("%s%d", c.col, r), c.val); err != nil {
				return nil, err
			}
			if err := f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", c.col, r), fmt.Sprintf("%s%d", c.col, r), numberCellStyle); err != nil {
				return nil, err
			}
		}

		// Transaction type & reference no
		if err := f.SetCellValue(sheetName, fmt.Sprintf("P%d", r), row.TransactionType); err != nil {
			return nil, err
		}
		if err := f.SetCellStyle(sheetName, fmt.Sprintf("P%d", r), fmt.Sprintf("P%d", r), textCellStyle); err != nil {
			return nil, err
		}

		if err := f.SetCellValue(sheetName, fmt.Sprintf("Q%d", r), row.ReferenceNo); err != nil {
			return nil, err
		}
		if err := f.SetCellStyle(sheetName, fmt.Sprintf("Q%d", r), fmt.Sprintf("Q%d", r), textCellStyle); err != nil {
			return nil, err
		}
	}

	// Auto-fit basic column widths for readability
	colsToWidth := map[string]float64{
		"A": 25,
		"B": 25,
		"C": 20,
		"D": 20,
		"E": 25,
		"F": 40,
		"P": 25,
		"Q": 25,
	}
	for col, width := range colsToWidth {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return nil, err
		}
	}

	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      2,
		TopLeftCell: "A3",
		ActivePane:  "bottomLeft",
	})

	f.AutoFilter(sheetName, "A2", nil)

	for _, col := range []string{"G", "H", "I", "J", "K", "L", "M", "N", "O"} {
		_ = f.SetColWidth(sheetName, col, col, 11)
	}

	return f, nil
}
