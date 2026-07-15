package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"finance/entity"
	"finance/model"
	"finance/repository"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/xuri/excelize/v2"
)

type PaymentDepositReportService interface {
	ListReport(dataFilter entity.PaymentDepositReportQueryFilter) (entity.PaymentDepositReportResponse, error)
	DownloadReport(dataFilter entity.PaymentDepositReportQueryFilter, createdBy string) (entity.ReportListResponse, error)
}

type paymentDepositReportServiceImpl struct {
	Repo        repository.PaymentDepositReportRepository
	Transaction repository.Dbtransaction
}

type paymentDepositExportMetadata struct {
	StartDate      time.Time
	EndDate        time.Time
	CollectorLabel string
}

func NewPaymentDepositReportService(repo repository.PaymentDepositReportRepository, transaction repository.Dbtransaction) PaymentDepositReportService {
	return &paymentDepositReportServiceImpl{
		Repo:        repo,
		Transaction: transaction,
	}
}

func (service *paymentDepositReportServiceImpl) ListReport(dataFilter entity.PaymentDepositReportQueryFilter) (entity.PaymentDepositReportResponse, error) {
	var response entity.PaymentDepositReportResponse

	// Get paginated data
	rows, total, _, err := service.Repo.FindAllPaymentDeposit(dataFilter, dataFilter.CustId)
	if err != nil {
		return response, err
	}

	// Calculate summary from ALL data (not just paginated)
	summaryRow, err := service.Repo.FindPaymentDepositSummary(dataFilter, dataFilter.CustId)
	if err != nil {
		return response, err
	}
	recapRows, err := service.Repo.FindPaymentDepositRecapRows(dataFilter, dataFilter.CustId, dataFilter.ParentCustId)
	if err != nil {
		return response, err
	}

	// Map rows to response items
	var items []entity.PaymentDepositReportItem
	for _, row := range rows {
		depositDate := ""
		if !row.DepositDate.IsZero() {
			depositDate = row.DepositDate.Format("2006-01-02")
		}
		item := entity.PaymentDepositReportItem{
			DepositDate:       depositDate,
			DepositType:       row.DepositType,
			DepositNo:         row.DepositNo,
			CollectorID:       row.CollectorID,
			CollectorCode:     row.CollectorCode,
			CollectorName:     row.CollectorName,
			CashAmount:        row.CashAmount,
			ChequeAmount:      row.ChequeAmount,
			TransferAmount:    row.TransferAmount,
			ReturnAmount:      row.ReturnAmount,
			CreditDebitAmount: row.CreditDebitAmount,
			ExpenseAmount:     row.ExpenseAmount,
			TotalPayment:      row.TotalPayment,
		}
		items = append(items, item)
	}

	// Map summary
	summary := entity.PaymentDepositReportSummary{
		TotalCash:        summaryRow.TotalCash,
		TotalCheque:      summaryRow.TotalCheque,
		TotalTransfer:    summaryRow.TotalTransfer,
		TotalReturn:      summaryRow.TotalReturn,
		TotalCreditDebit: summaryRow.TotalCreditDebit,
		TotalExpense:     summaryRow.TotalExpense,
		TotalAmount:      summaryRow.TotalCash + summaryRow.TotalCheque + summaryRow.TotalTransfer + summaryRow.TotalReturn + summaryRow.TotalCreditDebit,
	}
	summary.GrandTotal = summary.TotalAmount - summary.TotalExpense

	// Construct full response
	response.Items = items
	response.Summary = summary
	response.SummaryByDepositType = buildSummaryByDepositType(recapRows)
	response.Pagination = entity.PaymentDepositReportPagination{
		Page:      dataFilter.Page,
		Limit:     dataFilter.Limit,
		TotalData: total,
		TotalPage: int((total + int64(dataFilter.Limit) - 1) / int64(dataFilter.Limit)),
	}

	return response, nil
}

func (service *paymentDepositReportServiceImpl) DownloadReport(dataFilter entity.PaymentDepositReportQueryFilter, createdBy string) (entity.ReportListResponse, error) {
	// 1. Generate Report ID & Name
	reportID := xid.New().String()
	now := time.Now()
	startDate, err := time.Parse("2006-01-02", dataFilter.StartDate)
	if err != nil {
		return entity.ReportListResponse{}, err
	}
	endDate, err := time.Parse("2006-01-02", dataFilter.EndDate)
	if err != nil {
		return entity.ReportListResponse{}, err
	}

	// Running number logic
	runningNum, err := service.Repo.GetReportRunningNumber(dataFilter.CustId, now)
	if err != nil {
		return entity.ReportListResponse{}, err
	}
	reportName := fmt.Sprintf("%s-%s-%03d", entity.PaymentDepositReportDownloadPrefix, now.Format("020106"), runningNum+1)

	// 2. Prepare ReportList entry

	reportEntry := model.ReportList{
		CustID:     dataFilter.CustId,
		ReportID:   reportID,
		ReportName: reportName,
		StartDate:  startDate,
		EndDate:    endDate,
		FileStatus: entity.PaymentDepositReportStatusProcessing,
		CreatedBy:  &createdBy,
		CreatedAt:  now,
	}

	// 3. Insert to DB within transaction
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		return service.Repo.InsertReportList(txCtx, reportEntry)
	})
	if err != nil {
		return entity.ReportListResponse{}, err
	}

	// 4. Async Generation
	exportMeta := buildPaymentDepositExportMetadata(dataFilter, startDate, endDate)

	go func(meta paymentDepositExportMetadata) {
		// Fetch ALL data (no pagination limit)
		rows, fetchErr := service.Repo.FindAllPaymentDepositDownload(dataFilter, dataFilter.CustId, dataFilter.ParentCustId)
		if fetchErr != nil {
			log.Printf("payment deposit download fetch error: %v", fetchErr)
			return
		}

		meta.CollectorLabel = resolvePaymentDepositCollectorLabel(dataFilter, rows)
		fileBase64, genErr := service.generateExcel(rows, meta)
		if genErr != nil {
			log.Printf("payment deposit download generate error: %v", genErr)
			return
		}

		bgCtx := context.Background()
		if err := service.Transaction.WithinTransaction(bgCtx, func(txCtx context.Context) error {
			return service.Repo.UpdateReportList(txCtx, reportID, entity.PaymentDepositReportStatusReady, fileBase64)
		}); err != nil {
			log.Printf("payment deposit download update error: %v", err)
			return
		}
	}(exportMeta)

	// 5. Return immediate response
	return entity.ReportListResponse{
		ReportID:       reportID,
		ReportName:     reportName,
		StartDate:      startDate.Format("2006-01-02"),
		EndDate:        endDate.Format("2006-01-02"),
		FileStatus:     entity.PaymentDepositReportStatusProcessing,
		FileStatusName: entity.PaymentDepositReportStatusNameProcessing,
		CreatedBy:      createdBy,
		CreatedAt:      now.Format(time.RFC3339),
	}, nil
}

func (service *paymentDepositReportServiceImpl) generateExcel(rows []model.PaymentDepositReportDownloadRow, meta paymentDepositExportMetadata) (string, error) {
	f := excelize.NewFile()
	sheetName := "Payment Deposit Report"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1") // Remove default sheet

	// Styles
	styleBold, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	styleHeaderBold, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	styleBorder, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	styleCurrencyBorder, _ := f.NewStyle(&excelize.Style{
		NumFmt: 3,
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	styleCurrency, _ := f.NewStyle(&excelize.Style{
		NumFmt: 3,
	})

	// Row 1: Title
	f.SetCellValue(sheetName, "A1", "Payment Deposit Report")
	f.SetCellStyle(sheetName, "A1", "A1", styleBold)

	// Row 2: Date Range
	dateRange := fmt.Sprintf("%s - %s", meta.StartDate.Format("02-01-2006"), meta.EndDate.Format("02-01-2006"))
	f.SetCellValue(sheetName, "A2", "Deposit Date")
	f.SetCellValue(sheetName, "B2", dateRange)

	// Row 3: Collector Info
	collectorInfo := strings.TrimSpace(meta.CollectorLabel)
	if collectorInfo == "" {
		collectorInfo = "All"
	}
	f.SetCellValue(sheetName, "A3", "Collector")
	f.SetCellValue(sheetName, "B3", collectorInfo)

	// Row 5: Table Header
	headers := []string{
		"Deposit Date", "Deposit Type", "Deposit No", "Collector", "Document Date", "Code", "Business Name", "Document No",
		"Cash", "Cheque / Giro", "Transfer", "Return", "Credit / Debit", "Discount", "Payment Balance", "Expense", "Expense Name",
	}

	startRow := 5
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, startRow)
		f.SetCellValue(sheetName, cell, h)
		f.SetCellStyle(sheetName, cell, cell, styleHeaderBold)
	}

	// Data Rows
	currentRow := startRow + 1
	for _, row := range rows {
		collector := ""
		if row.Collector != nil {
			collector = *row.Collector
		}
		documentDate := ""
		if row.DocumentDate != nil {
			documentDate = row.DocumentDate.Format("2006-01-02")
		}
		code := ""
		if row.Code != nil {
			code = *row.Code
		}
		businessName := ""
		if row.BusinessName != nil {
			businessName = *row.BusinessName
		}
		documentNo := ""
		if row.DocumentNo != nil {
			documentNo = *row.DocumentNo
		}
		expenseName := ""
		if row.ExpenseName != nil {
			expenseName = *row.ExpenseName
		}

		values := []interface{}{row.DepositDate.Format("2006-01-02"), row.DepositType, row.DepositNo, collector, documentDate, code, businessName, documentNo, row.Cash, row.ChequeGiro, row.Transfer, row.ReturnAmount, row.CreditDebit, row.Discount, row.PaymentBalance, row.Expense, expenseName}
		for i, value := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, currentRow)
			f.SetCellValue(sheetName, cell, value)
		}

		// Apply border and currency styles to data cells
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("H%d", currentRow), styleBorder)
		f.SetCellStyle(sheetName, fmt.Sprintf("I%d", currentRow), fmt.Sprintf("Q%d", currentRow), styleCurrencyBorder)
		currentRow++
	}

	recapStartRow := currentRow + 2
	recapByType := buildRecapFromDownloadRows(rows)
	ar := recapByType["Account Receivable"]
	ap := recapByType["Account Payable"]

	f.SetCellValue(sheetName, fmt.Sprintf("B%d", recapStartRow), "Account Receivable")
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", recapStartRow), "Account Payable")
	f.SetCellStyle(sheetName, fmt.Sprintf("B%d", recapStartRow), fmt.Sprintf("B%d", recapStartRow), styleBold)
	f.SetCellStyle(sheetName, fmt.Sprintf("E%d", recapStartRow), fmt.Sprintf("E%d", recapStartRow), styleBold)

	arLabels := []string{"Total Cash", "Total Cheque / Giro", "Total Transfer", "Total Return", "Total Credit / Debit", "Total Discount", "Total Payment Balance", "Total Expense"}
	arValues := []float64{ar.Cash, ar.ChequeGiro, ar.Transfer, ar.ReturnAmount, ar.CreditDebit, ar.Discount, ar.PaymentBalance, ar.Expense}
	for i, label := range arLabels {
		rowNo := recapStartRow + 1 + i
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNo), label)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNo), arValues[i])
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowNo), fmt.Sprintf("A%d", rowNo), styleBold)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", rowNo), fmt.Sprintf("B%d", rowNo), styleCurrency)
	}

	apLabels := []string{"Total Cash", "Total Cheque / Giro", "Total Transfer", "Total Return", "Total Credit / Debit", "Total Discount", "Total Payment Balance"}
	apValues := []float64{ap.Cash, ap.ChequeGiro, ap.Transfer, ap.ReturnAmount, ap.CreditDebit, ap.Discount, ap.PaymentBalance}
	for i, label := range apLabels {
		rowNo := recapStartRow + 1 + i
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNo), label)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNo), apValues[i])
		f.SetCellStyle(sheetName, fmt.Sprintf("D%d", rowNo), fmt.Sprintf("D%d", rowNo), styleBold)
		f.SetCellStyle(sheetName, fmt.Sprintf("E%d", rowNo), fmt.Sprintf("E%d", rowNo), styleCurrency)
	}

	// Buffer to Base64
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return encoded, nil
}

func buildPaymentDepositExportMetadata(dataFilter entity.PaymentDepositReportQueryFilter, startDate, endDate time.Time) paymentDepositExportMetadata {
	return paymentDepositExportMetadata{
		StartDate:      startDate,
		EndDate:        endDate,
		CollectorLabel: "All",
	}
}

func resolvePaymentDepositCollectorLabel(dataFilter entity.PaymentDepositReportQueryFilter, rows []model.PaymentDepositReportDownloadRow) string {
	if len(dataFilter.EmpID) != 1 {
		return "All"
	}

	collectorName := ""
	for _, row := range rows {
		if row.Collector == nil {
			continue
		}
		name := strings.TrimSpace(*row.Collector)
		if name == "" {
			continue
		}
		if collectorName == "" {
			collectorName = name
			continue
		}
		if collectorName != name {
			return "All"
		}
	}

	if collectorName == "" {
		return "All"
	}
	return collectorName
}

func buildRecapFromDownloadRows(rows []model.PaymentDepositReportDownloadRow) map[string]model.PaymentDepositReportRecapRow {
	result := map[string]model.PaymentDepositReportRecapRow{
		"Account Receivable": {DepositType: "Account Receivable"},
		"Account Payable":    {DepositType: "Account Payable"},
	}
	for _, row := range rows {
		recap, ok := result[row.DepositType]
		if !ok {
			continue
		}
		recap.Cash += row.Cash
		recap.ChequeGiro += row.ChequeGiro
		recap.Transfer += row.Transfer
		recap.ReturnAmount += row.ReturnAmount
		recap.CreditDebit += row.CreditDebit
		recap.Discount += row.Discount
		recap.PaymentBalance += row.PaymentBalance
		recap.Expense += row.Expense
		result[row.DepositType] = recap
	}
	result["Account Payable"] = model.PaymentDepositReportRecapRow{
		DepositType:    "Account Payable",
		Cash:           result["Account Payable"].Cash,
		ChequeGiro:     result["Account Payable"].ChequeGiro,
		Transfer:       result["Account Payable"].Transfer,
		ReturnAmount:   result["Account Payable"].ReturnAmount,
		CreditDebit:    result["Account Payable"].CreditDebit,
		Discount:       result["Account Payable"].Discount,
		PaymentBalance: result["Account Payable"].PaymentBalance,
		Expense:        0,
	}
	return result
}

func buildSummaryByDepositType(rows []model.PaymentDepositReportRecapRow) []entity.PaymentDepositReportSummaryByDepositTypeItem {
	ordered := []string{"Account Receivable", "Account Payable"}
	byType := map[string]model.PaymentDepositReportRecapRow{}
	for _, row := range rows {
		byType[row.DepositType] = row
	}
	result := make([]entity.PaymentDepositReportSummaryByDepositTypeItem, 0, len(ordered))
	for _, label := range ordered {
		row, ok := byType[label]
		if !ok {
			row = model.PaymentDepositReportRecapRow{DepositType: label}
		}
		if label == "Account Payable" {
			row.Expense = 0
		}
		result = append(result, entity.PaymentDepositReportSummaryByDepositTypeItem{
			DepositTypeLabel:    label,
			SummaryCash:         row.Cash,
			ChequeGiro:          row.ChequeGiro,
			Transfer:            row.Transfer,
			ReturnAmount:        row.ReturnAmount,
			CreditDebit:         row.CreditDebit,
			Discount:            row.Discount,
			TotalPaymentBalance: row.PaymentBalance,
			TotalExpense:        row.Expense,
		})
	}
	return result
}
