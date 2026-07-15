package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"master/entity"
	"master/model"
	"master/repository"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/xuri/excelize/v2"
)

type ProductRipeningService interface {
	List(filter entity.ProductRipeningQueryFilter, picUserID int64) ([]entity.ProductRipeningListItem, int, int, error)
	Detail(params entity.ProductRipeningDetailParams, custID, parentCustID string, picUserID int64) (entity.ProductRipeningDetailResponse, error)
	Update(params entity.ProductRipeningDetailParams, req entity.ProductRipeningUpdateRequest, custID, parentCustID string, picUserID int64) error
	Export(filter entity.ProductRipeningQueryFilter, picUserID int64) (*bytes.Buffer, string, string, error)
	DownloadTemplate(format, custID, parentCustID string, picUserID int64) (*bytes.Buffer, string, string, error)
	Import(req entity.ProductRipeningImportRequest, custID, parentCustID string, picUserID int64) (entity.ProductRipeningImportResponse, error)
}

type productRipeningService struct {
	repo        repository.ProductRipeningRepository
	productRepo repository.ProductRepository
}

func NewProductRipeningService(repo repository.ProductRipeningRepository, productRepo repository.ProductRepository) ProductRipeningService {
	return &productRipeningService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *productRipeningService) List(filter entity.ProductRipeningQueryFilter, picUserID int64) ([]entity.ProductRipeningListItem, int, int, error) {
	// if err := s.ensurePICAccess(filter.CustId, filter.ParentCustId, picUserID); err != nil {
	// 	return nil, 0, 0, err
	// }

	rows, total, lastPage, err := s.repo.ListPlans(filter, picUserID, time.Now())
	if err != nil {
		return nil, 0, 0, err
	}

	items := make([]entity.ProductRipeningListItem, 0, len(rows))
	today := startOfDay(time.Now())
	for _, row := range rows {
		isActive := isRipeningPlanActive(row.WeekEnd, today)
		item := entity.ProductRipeningListItem{
			ID:              row.ID,
			DistributorID:   row.DistributorID,
			DistributorCode: row.DistributorCode,
			DistributorName: row.DistributorName,
			PerYear:         row.PerYear,
			PerID:           row.PerID,
			WeekID:          row.WeekID,
			WeekStart:       first10(row.WeekStart),
			WeekEnd:         first10(row.WeekEnd),
			WeekLabel:       fmt.Sprintf("Week %d (%s - %s)", row.WeekID, first10(row.WeekStart), first10(row.WeekEnd)),
			IsActive:        isActive,
			CanEdit:         isActive,
			TotalProduct:    row.TotalProduct,
			CreatedBy:       row.CreatedBy,
			CreatedByName:   row.CreatedByName,
			CreatedAt:       row.CreatedAt.Format(time.RFC3339),
			UpdatedBy:       row.UpdatedBy,
			UpdatedByName:   row.UpdatedByName,
		}
		if row.UpdatedAt != nil {
			updatedAt := row.UpdatedAt.Format(time.RFC3339)
			item.UpdatedAt = &updatedAt
		}
		items = append(items, item)
	}
	return items, total, lastPage, nil
}

func (s *productRipeningService) Detail(params entity.ProductRipeningDetailParams, custID, parentCustID string, picUserID int64) (entity.ProductRipeningDetailResponse, error) {
	rows, err := s.repo.FindRowsByPlan(custID, parentCustID, picUserID, params.DistributorID, params.PerYear, params.WeekID)
	if err != nil {
		return entity.ProductRipeningDetailResponse{}, fmt.Errorf("product ripening not found")
	}
	if len(rows) == 0 {
		return entity.ProductRipeningDetailResponse{}, fmt.Errorf("product ripening not found")
	}

	today := startOfDay(time.Now())
	return mapProductRipeningDetail(rows, today), nil
}

func (s *productRipeningService) Update(params entity.ProductRipeningDetailParams, req entity.ProductRipeningUpdateRequest, custID, parentCustID string, picUserID int64) error {
	if len(req.Rows) == 0 {
		return fmt.Errorf("product ripening rows are required")
	}

	currentRows, err := s.repo.FindRowsByPlan(custID, parentCustID, picUserID, params.DistributorID, params.PerYear, params.WeekID)
	if err != nil {
		return fmt.Errorf("product ripening not found")
	}
	if len(currentRows) == 0 {
		return fmt.Errorf("product ripening not found")
	}

	preserveBefore, canEdit, err := resolveProductRipeningEditWindow(currentRows[0].WeekStart, currentRows[0].WeekEnd, time.Now())
	if err != nil {
		return err
	}
	if !canEdit {
		return fmt.Errorf("product ripening can no longer be edited")
	}

	existingByID := make(map[int64]model.ProductRipening, len(currentRows))
	finalRows := make([]model.ProductRipening, 0, len(currentRows))
	for _, row := range currentRows {
		existingByID[row.ID] = row
		finalRows = append(finalRows, row)
	}

	indexByID := make(map[int64]int, len(finalRows))
	for idx, row := range finalRows {
		indexByID[row.ID] = idx
	}

	seen := make(map[int64]struct{}, len(req.Rows))
	for _, item := range req.Rows {
		if _, exists := seen[item.ID]; exists {
			return fmt.Errorf("duplicate product ripening row %d", item.ID)
		}
		seen[item.ID] = struct{}{}

		existing, ok := existingByID[item.ID]
		if !ok {
			return fmt.Errorf("product ripening row %d not found in selected plan", item.ID)
		}

		finalRows[indexByID[item.ID]] = model.ProductRipening{
			ID:            existing.ID,
			CustID:        existing.CustID,
			DistributorID: existing.DistributorID,
			ProID:         existing.ProID,
			PerYear:       existing.PerYear,
			PerID:         existing.PerID,
			WeekID:        existing.WeekID,
			SundayQty:     item.SundayQty,
			MondayQty:     item.MondayQty,
			TuesdayQty:    item.TuesdayQty,
			WednesdayQty:  item.WednesdayQty,
			ThursdayQty:   item.ThursdayQty,
			FridayQty:     item.FridayQty,
			SaturdayQty:   item.SaturdayQty,
			CreatedBy:     existing.CreatedBy,
			CreatedAt:     existing.CreatedAt,
		}
	}

	return s.repo.ReplacePlanRows(
		custID,
		currentRows[0].DistributorID,
		currentRows[0].PerYear,
		currentRows[0].PerID,
		currentRows[0].WeekID,
		finalRows,
		preserveBefore,
		picUserID,
	)
}

func (s *productRipeningService) Export(filter entity.ProductRipeningQueryFilter, picUserID int64) (*bytes.Buffer, string, string, error) {
	// if err := s.ensurePICAccess(filter.CustId, filter.ParentCustId, picUserID); err != nil {
	// 	return nil, "", "", err
	// }

	rows, err := s.repo.ExportRows(filter, picUserID)
	if err != nil {
		return nil, "", "", err
	}
	switch strings.ToLower(strings.TrimSpace(filter.Format)) {
	case "csv":
		buf, err := createProductRipeningExportCSV(rows)
		if err != nil {
			return nil, "", "", err
		}
		return buf, "text/csv", "product_ripening.csv", nil
	case "xls":
		buf, err := createProductRipeningExportWorkbook(rows)
		if err != nil {
			return nil, "", "", err
		}
		return buf, "application/vnd.ms-excel", "product_ripening.xls", nil
	default:
		buf, err := createProductRipeningExportWorkbook(rows)
		if err != nil {
			return nil, "", "", err
		}
		return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "product_ripening.xlsx", nil
	}
}

func (s *productRipeningService) DownloadTemplate(format, custID, parentCustID string, picUserID int64) (*bytes.Buffer, string, string, error) {
	// if err := s.ensurePICAccess(custID, parentCustID, picUserID); err != nil {
	// 	return nil, "", "", err
	// }

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "csv":
		buf, err := createProductRipeningTemplateCSV()
		if err != nil {
			return nil, "", "", err
		}
		return buf, "text/csv", "product_ripening_template.csv", nil
	case "xls":
		buf, err := createProductRipeningTemplateWorkbook()
		if err != nil {
			return nil, "", "", err
		}
		return buf, "application/vnd.ms-excel", "product_ripening_template.xls", nil
	default:
		buf, err := createProductRipeningTemplateWorkbook()
		if err != nil {
			return nil, "", "", err
		}
		return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "product_ripening_template.xlsx", nil
	}
}

func (s *productRipeningService) Import(req entity.ProductRipeningImportRequest, custID, parentCustID string, picUserID int64) (entity.ProductRipeningImportResponse, error) {
	// if err := s.ensurePICAccess(custID, parentCustID, picUserID); err != nil {
	// 	return entity.ProductRipeningImportResponse{}, err
	// }

	resp := entity.ProductRipeningImportResponse{
		FileURL:     req.FileURL,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}

	rows, fileName, err := downloadProductRipeningImportRows(req.FileURL)
	if err != nil {
		return resp, err
	}
	resp.FileName = fileName

	headerIndex, err := detectProductRipeningHeader(rows)
	if err != nil {
		return resp, err
	}
	dataRows := rows[headerIndex+1:]
	if len(dataRows) == 0 {
		return resp, fmt.Errorf("template does not contain data rows")
	}

	type groupedPlan struct {
		distributor model.ProductRipeningAssignedDistributor
		week        model.ProductRipeningWeek
		rows        []model.ProductRipening
	}

	plans := map[string]*groupedPlan{}
	failed := make([]string, 0)
	for idx, row := range dataRows {
		actualRow := idx + headerIndex + 2
		if isProductRipeningImportRowEmpty(row) {
			continue
		}
		resp.TotalRow++

		distributorCode := productRipeningCell(row, 0)
		productCode := productRipeningCell(row, 1)
		perYearText := productRipeningCell(row, 2)
		weekIDText := productRipeningCell(row, 3)
		if distributorCode == "" || productCode == "" || perYearText == "" || weekIDText == "" {
			failed = append(failed, fmt.Sprintf("row %d: distributor code, product code, year, and week are required", actualRow))
			continue
		}

		perYear, err := parseProductRipeningWholeNumber(perYearText, "year", false)
		if err != nil {
			failed = append(failed, fmt.Sprintf("row %d: %v", actualRow, err))
			continue
		}

		weekID, err := parseProductRipeningWholeNumber(weekIDText, "week", false)
		if err != nil {
			failed = append(failed, fmt.Sprintf("row %d: %v", actualRow, err))
			continue
		}

		qtys, err := parseProductRipeningQtys(row, 4)
		if err != nil {
			failed = append(failed, fmt.Sprintf("row %d: %v", actualRow, err))
			continue
		}

		distributor, err := s.repo.FindDistributorByCode(parentCustID, distributorCode)
		if err != nil {
			failed = append(failed, fmt.Sprintf("row %d: distributor %s not found", actualRow, distributorCode))
			continue
		}

		week, err := s.repo.FindWeekByYearAndWeekID(parentCustID, distributor.DistributorCustID, perYear, weekID, time.Now())
		if err != nil {
			failed = append(failed, fmt.Sprintf("row %d: year %d week %d not found or already passed for distributor %s", actualRow, perYear, weekID, distributorCode))
			continue
		}

		product, err := s.productRepo.FindOneByProductCodeAndCustId(productCode, distributor.DistributorCustID)
		if err != nil || product.ProductId == 0 {
			failed = append(failed, fmt.Sprintf("row %d: product %s not found for distributor %s", actualRow, productCode, distributorCode))
			continue
		}

		groupKey := fmt.Sprintf("%d:%d:%d:%d", distributor.DistributorID, week.PerYear, week.PerID, week.WeekID)
		if _, ok := plans[groupKey]; !ok {
			plans[groupKey] = &groupedPlan{
				distributor: distributor,
				week:        week,
				rows:        []model.ProductRipening{},
			}
		}
		for _, existing := range plans[groupKey].rows {
			if existing.ProID == int64(product.ProductId) {
				failed = append(failed, fmt.Sprintf("row %d: duplicate product %s for distributor %s week %d", actualRow, productCode, distributorCode, week.WeekID))
				err = fmt.Errorf("duplicate")
				break
			}
		}
		if err != nil && err.Error() == "duplicate" {
			continue
		}

		plans[groupKey].rows = append(plans[groupKey].rows, model.ProductRipening{
			CustID:        custID,
			DistributorID: distributor.DistributorID,
			ProID:         int64(product.ProductId),
			PerYear:       week.PerYear,
			PerID:         week.PerID,
			WeekID:        week.WeekID,
			SundayQty:     qtys[0],
			MondayQty:     qtys[1],
			TuesdayQty:    qtys[2],
			WednesdayQty:  qtys[3],
			ThursdayQty:   qtys[4],
			FridayQty:     qtys[5],
			SaturdayQty:   qtys[6],
		})
	}

	resp.FailedReasons = failed
	resp.FailedRow = len(failed)
	if len(failed) > 0 {
		return resp, fmt.Errorf("import validation failed")
	}

	success := 0
	planKeys := make([]string, 0, len(plans))
	for key := range plans {
		planKeys = append(planKeys, key)
	}
	sort.Strings(planKeys)

	for _, key := range planKeys {
		plan := plans[key]
		weekEnd, err := time.Parse("2006-01-02", first10(plan.week.WeekEnd))
		if err != nil {
			resp.FailedReasons = append(resp.FailedReasons, fmt.Sprintf("week %d: invalid week end date", plan.week.WeekID))
			continue
		}
		today := startOfDay(time.Now())
		if weekEnd.Before(today) {
			resp.FailedReasons = append(resp.FailedReasons, fmt.Sprintf("distributor %s week %d already passed", plan.distributor.DistributorCode, plan.week.WeekID))
			continue
		}

		weekStart, _ := time.Parse("2006-01-02", first10(plan.week.WeekStart))
		var preserveBefore time.Time
		if !today.Before(weekStart) {
			preserveBefore = today
		}

		if err := s.repo.ReplacePlanRows(custID, plan.distributor.DistributorID, plan.week.PerYear, plan.week.PerID, plan.week.WeekID, plan.rows, preserveBefore, picUserID); err != nil {
			resp.FailedReasons = append(resp.FailedReasons, fmt.Sprintf("distributor %s week %d: %v", plan.distributor.DistributorCode, plan.week.WeekID, err))
			continue
		}
		success += len(plan.rows)
	}

	resp.SuccessRow = success
	resp.FailedRow = len(resp.FailedReasons)
	if resp.FailedRow > 0 {
		return resp, fmt.Errorf("import failed")
	}
	return resp, nil
}

func (s *productRipeningService) ensurePICAccess(custID, parentCustID string, picUserID int64) error {
	hasAccess, err := s.repo.HasAssignedDistributor(custID, parentCustID, picUserID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return fmt.Errorf("product ripening is only available for users assigned as distributor PIC")
	}
	return nil
}

func mapProductRipeningDetail(rows []model.ProductRipening, today time.Time) entity.ProductRipeningDetailResponse {
	head := rows[0]
	isActive := isRipeningPlanActive(head.WeekEnd, today)
	response := entity.ProductRipeningDetailResponse{
		DistributorID:   head.DistributorID,
		DistributorCode: head.DistributorCode,
		DistributorName: head.DistributorName,
		PerYear:         head.PerYear,
		PerID:           head.PerID,
		WeekID:          head.WeekID,
		WeekStart:       first10(head.WeekStart),
		WeekEnd:         first10(head.WeekEnd),
		WeekLabel:       fmt.Sprintf("Week %d (%s - %s)", head.WeekID, first10(head.WeekStart), first10(head.WeekEnd)),
		IsActive:        isActive,
		CanEdit:         isActive,
		EditableDays:    buildProductRipeningEditableDays(head.WeekStart, head.WeekEnd, today),
		Rows:            make([]entity.ProductRipeningDetailRow, 0, len(rows)),
	}
	for _, row := range rows {
		item := entity.ProductRipeningDetailRow{
			ID:            row.ID,
			ProID:         row.ProID,
			ProductCode:   row.ProductCode,
			ProductName:   row.ProductName,
			SundayQty:     row.SundayQty,
			MondayQty:     row.MondayQty,
			TuesdayQty:    row.TuesdayQty,
			WednesdayQty:  row.WednesdayQty,
			ThursdayQty:   row.ThursdayQty,
			FridayQty:     row.FridayQty,
			SaturdayQty:   row.SaturdayQty,
			CreatedBy:     row.CreatedBy,
			CreatedByName: row.CreatedByName,
			CreatedAt:     row.CreatedAt.Format(time.RFC3339),
			UpdatedBy:     row.UpdatedBy,
			UpdatedByName: row.UpdatedByName,
		}
		if row.UpdatedAt != nil {
			updatedAt := row.UpdatedAt.Format(time.RFC3339)
			item.UpdatedAt = &updatedAt
		}
		response.Rows = append(response.Rows, item)
	}
	return response
}

func resolveProductRipeningEditWindow(weekStartText, weekEndText string, now time.Time) (time.Time, bool, error) {
	weekStart, err := time.Parse("2006-01-02", first10(weekStartText))
	if err != nil {
		return time.Time{}, false, fmt.Errorf("invalid product ripening week start")
	}
	weekEnd, err := time.Parse("2006-01-02", first10(weekEndText))
	if err != nil {
		return time.Time{}, false, fmt.Errorf("invalid product ripening week end")
	}

	today := startOfDay(now)
	if weekEnd.Before(today) {
		return time.Time{}, false, fmt.Errorf("product ripening can no longer be edited")
	}
	if weekStart.After(today) {
		return time.Time{}, true, nil
	}
	return today, true, nil
}

func buildProductRipeningEditableDays(weekStartText, weekEndText string, now time.Time) entity.ProductRipeningEditableDays {
	editable := entity.ProductRipeningEditableDays{}
	weekStart, errStart := time.Parse("2006-01-02", first10(weekStartText))
	weekEnd, errEnd := time.Parse("2006-01-02", first10(weekEndText))
	if errStart != nil || errEnd != nil {
		return editable
	}

	today := startOfDay(now)
	if weekEnd.Before(today) {
		return editable
	}
	if weekStart.After(today) {
		return entity.ProductRipeningEditableDays{Sunday: true, Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true, Saturday: true}
	}

	flags := []bool{false, false, false, false, false, false, false}
	for i := ripeningWeekdayIndex(today); i < len(flags); i++ {
		flags[i] = true
	}
	editable.Sunday = flags[0]
	editable.Monday = flags[1]
	editable.Tuesday = flags[2]
	editable.Wednesday = flags[3]
	editable.Thursday = flags[4]
	editable.Friday = flags[5]
	editable.Saturday = flags[6]
	return editable
}

func isRipeningPlanActive(weekEndText string, today time.Time) bool {
	weekEnd, err := time.Parse("2006-01-02", first10(weekEndText))
	if err != nil {
		return false
	}
	return !weekEnd.Before(today)
}

func productRipeningExportHeaders() []string {
	return []string{"Plan ID", "Distributor Code", "Distributor Name", "Per Year", "Period ID", "Week", "Week Start", "Week End", "Total Product", "Created By", "Created By Name", "Created At", "Updated By", "Updated By Name", "Updated At"}
}

func createProductRipeningExportWorkbook(rows []model.ProductRipeningPlanListRow) (*bytes.Buffer, error) {
	loc := productRipeningExportLocation()
	f := excelize.NewFile()
	sheet := "Product Ripening"
	idx, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}
	style, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}, Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1}})
	headers := productRipeningExportHeaders()
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, style)
	}
	for r, row := range rows {
		values := []interface{}{row.ID, row.DistributorCode, row.DistributorName, row.PerYear, row.PerID, row.WeekID, first10(row.WeekStart), first10(row.WeekEnd), row.TotalProduct, row.CreatedBy, row.CreatedByName, formatProductRipeningExportTime(row.CreatedAt, loc), derefInt64(row.UpdatedBy), derefString(row.UpdatedByName), formatNullableProductRipeningExportTime(row.UpdatedAt, loc)}
		for c, value := range values {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			f.SetCellValue(sheet, cell, value)
		}
	}
	f.SetActiveSheet(idx)
	_ = f.DeleteSheet("Sheet1")
	return f.WriteToBuffer()
}

func createProductRipeningExportCSV(rows []model.ProductRipeningPlanListRow) (*bytes.Buffer, error) {
	loc := productRipeningExportLocation()
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	if err := writer.Write(productRipeningExportHeaders()); err != nil {
		return nil, err
	}
	for _, row := range rows {
		record := []string{strconv.FormatInt(row.ID, 10), row.DistributorCode, row.DistributorName, strconv.Itoa(row.PerYear), strconv.Itoa(row.PerID), strconv.Itoa(row.WeekID), first10(row.WeekStart), first10(row.WeekEnd), strconv.Itoa(row.TotalProduct), strconv.FormatInt(row.CreatedBy, 10), row.CreatedByName, formatProductRipeningExportTime(row.CreatedAt, loc), formatNullableInt64(row.UpdatedBy), derefString(row.UpdatedByName), formatNullableProductRipeningExportTime(row.UpdatedAt, loc)}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return buf, writer.Error()
}

func createProductRipeningTemplateWorkbook() (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheet := "Template"
	idx, err := f.NewSheet(sheet)
	if err != nil {
		return nil, err
	}
	headerStyle, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}, Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1}, Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"}})
	dayStyle, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}, Fill: excelize.Fill{Type: "pattern", Color: []string{"D9EAF7"}, Pattern: 1}, Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"}})

	headers := []string{"Distributor Code", "Product Code", "Year", "Week"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.MergeCell(sheet, cell, fmt.Sprintf("%s2", string(rune('A'+i))))
		f.SetCellStyle(sheet, cell, fmt.Sprintf("%s2", string(rune('A'+i))), headerStyle)
	}
	f.SetCellValue(sheet, "E1", "Ripening QTY (in Largest Unit)")
	f.MergeCell(sheet, "E1", "K1")
	f.SetCellStyle(sheet, "E1", "K1", headerStyle)
	days := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	for i, day := range days {
		cell, _ := excelize.CoordinatesToCellName(i+5, 2)
		f.SetCellValue(sheet, cell, day)
		f.SetCellStyle(sheet, cell, cell, dayStyle)
	}

	instructionSheet := "Instructions"
	if _, err := f.NewSheet(instructionSheet); err != nil {
		return nil, err
	}
	instructions := [][]string{
		{"Column", "Instruction"},
		{"Distributor Code", "Required. Distributor must be assigned to the current PIC in distributor replenishment setup."},
		{"Product Code", "Required. Product must exist in the distributor product master (mst.m_product)."},
		{"Year", "Required. Use the distributor week master year from mst.m_week."},
		{"Week", "Required. Use week_id from mst.m_week. The selected year + week must still be active/not passed for the distributor."},
		{"Sunday-Saturday", "Required. Integer only. Zero is allowed. Blank value is not allowed."},
		{"Re-upload", "Upload the full weekly plan for the distributor. Existing future values will be replaced. Past dates in the current week are preserved."},
	}
	for r, row := range instructions {
		for c, value := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+1)
			f.SetCellValue(instructionSheet, cell, value)
			if r == 0 {
				f.SetCellStyle(instructionSheet, cell, cell, headerStyle)
			}
		}
	}
	f.SetActiveSheet(idx)
	_ = f.DeleteSheet("Sheet1")
	return f.WriteToBuffer()
}

func createProductRipeningTemplateCSV() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	if err := writer.Write([]string{"Distributor Code", "Product Code", "Year", "Week", "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}); err != nil {
		return nil, err
	}
	writer.Flush()
	return buf, writer.Error()
}

func downloadProductRipeningImportRows(fileURL string) ([][]string, string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download file: status %d", resp.StatusCode)
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	fileName := fileURL
	if parts := strings.Split(fileURL, "/"); len(parts) > 0 {
		fileName = parts[len(parts)-1]
	}
	if strings.HasSuffix(strings.ToLower(fileName), ".csv") {
		reader := csv.NewReader(bytes.NewReader(content))
		rows, err := reader.ReadAll()
		return rows, fileName, err
	}
	f, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = f.Close() }()
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, "", fmt.Errorf("excel file has no sheet")
	}
	rows, err := f.GetRows(sheetName)
	return rows, fileName, err
}

func detectProductRipeningHeader(rows [][]string) (int, error) {
	if len(rows) == 0 {
		return 0, fmt.Errorf("template does not contain header")
	}
	if matchesFlatRipeningHeader(rows[0]) {
		return 0, nil
	}
	if len(rows) > 1 && matchesMergedRipeningHeader(rows[0], rows[1]) {
		return 1, nil
	}
	return 0, fmt.Errorf("invalid template header")
}

func matchesFlatRipeningHeader(row []string) bool {
	expected := []string{"Distributor Code", "Product Code", "Year", "Week", "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	for i, want := range expected {
		if normalizeProductRipeningHeader(productRipeningCell(row, i)) != normalizeProductRipeningHeader(want) {
			return false
		}
	}
	return true
}

func matchesMergedRipeningHeader(row1, row2 []string) bool {
	if productRipeningCell(row1, 0) != "Distributor Code" || productRipeningCell(row1, 1) != "Product Code" || productRipeningCell(row1, 2) != "Year" || productRipeningCell(row1, 3) != "Week" {
		return false
	}
	expectedDays := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	for i, want := range expectedDays {
		if normalizeProductRipeningHeader(productRipeningCell(row2, i+4)) != normalizeProductRipeningHeader(want) {
			return false
		}
	}
	return true
}

func parseProductRipeningQtys(row []string, startIndex int) ([]int, error) {
	out := make([]int, 7)
	for i := 0; i < 7; i++ {
		n, err := parseProductRipeningWholeNumber(productRipeningCell(row, i+startIndex), strings.ToLower([]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}[i])+" quantity", true)
		if err != nil {
			return nil, err
		}
		out[i] = n
	}
	return out, nil
}

func parseProductRipeningWholeNumber(value, field string, allowZero bool) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, fmt.Errorf("%s is required", field)
	}
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", field)
	}
	if number != float64(int(number)) {
		return 0, fmt.Errorf("%s must be an integer", field)
	}
	result := int(number)
	if allowZero {
		if result < 0 {
			return 0, fmt.Errorf("%s must be greater than or equal to 0", field)
		}
		return result, nil
	}
	if result < 1 {
		return 0, fmt.Errorf("%s must be greater than or equal to 1", field)
	}
	return result, nil
}

func normalizeProductRipeningHeader(value string) string {
	return strings.TrimSpace(strings.ToUpper(value))
}

func productRipeningCell(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func isProductRipeningImportRowEmpty(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func first10(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func derefInt64(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}

func formatNullableInt64(value *int64) string {
	if value == nil {
		return ""
	}
	return strconv.FormatInt(*value, 10)
}

func productRipeningExportLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.FixedZone("WIB", 7*60*60)
	}
	return loc
}

func formatProductRipeningExportTime(value time.Time, loc *time.Location) string {
	return value.In(loc).Format(time.RFC3339)
}

func formatNullableProductRipeningExportTime(value *time.Time, loc *time.Location) string {
	if value == nil {
		return ""
	}
	return formatProductRipeningExportTime(*value, loc)
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func ripeningWeekdayIndex(t time.Time) int {
	switch t.Weekday() {
	case time.Sunday:
		return 0
	case time.Monday:
		return 1
	case time.Tuesday:
		return 2
	case time.Wednesday:
		return 3
	case time.Thursday:
		return 4
	case time.Friday:
		return 5
	default:
		return 6
	}
}

var _ ProductRipeningService = (*productRipeningService)(nil)
var _ = log.Info
