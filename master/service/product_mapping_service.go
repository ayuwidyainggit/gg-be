package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/constant"
	"master/repository"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2/log"
	"github.com/xuri/excelize/v2"
)

const productMappingOrigin = "product_mapping"

type productMappingImportCandidate struct {
	rowNumber       int
	distributorCode string
	parentProCode   string
	proCode         string
	proName         string
	largestUOM      string
	middleUOM       string
	smallestUOM     string
	distCustID      string
	distributorID   int64
	parentProduct   model.Product
}

var (
	productMappingUOMPattern     = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	productMappingProCodePattern = regexp.MustCompile(`^[a-zA-Z0-9\W_]+$`)

	productMappingTemplateHeaders = []string{
		"distributor_code",
		"principal_product_code",
		"distributor_product_code",
		"distributor_product_name",
		"largest_unit",
		"middle_unit",
		"smallest_unit",
	}
)

type ProductMappingService interface {
	List(dataFilter entity.ProductMappingListQueryFilter) ([]entity.ProductMappingListItem, int, int, error)
	Detail(dataFilter entity.ProductMappingDetailQueryFilter) (entity.ProductMappingDetailResponse, int, int, error)
	Update(proID int64, req entity.ProductMappingUpdateRequest, principalCustID string, updatedBy int64) error
	Delete(proID int64, principalCustID string, deletedBy int64) error
	DownloadTemplate() (*bytes.Buffer, string, string, error)
	Import(req entity.ProductMappingImportRequest, principalCustID string, createdBy int64) (entity.ProductMappingImportResponse, error)
}

func NewProductMappingService(
	productMappingRepository repository.ProductMappingRepository,
	productRepository repository.ProductRepository,
	distributorRepository repository.DistributorRepository,
	tx repository.TransactionManager,
) ProductMappingService {
	return &productMappingServiceImpl{
		tx:                       tx,
		ProductMappingRepository: productMappingRepository,
		ProductRepository:        productRepository,
		DistributorRepository:    distributorRepository,
	}
}

type productMappingServiceImpl struct {
	tx                       repository.TransactionManager
	ProductMappingRepository repository.ProductMappingRepository
	ProductRepository        repository.ProductRepository
	DistributorRepository    repository.DistributorRepository
}

func (s *productMappingServiceImpl) List(dataFilter entity.ProductMappingListQueryFilter) ([]entity.ProductMappingListItem, int, int, error) {
	normalizeProductMappingPagination(&dataFilter.Page, &dataFilter.Limit)

	rows, total, lastPage, err := s.ProductMappingRepository.FindDistributorSummary(dataFilter)
	if err != nil {
		log.Error("ProductMappingService, List, err:", err.Error())
		return nil, 0, 0, err
	}

	items := make([]entity.ProductMappingListItem, 0, len(rows))
	for _, row := range rows {
		item := entity.ProductMappingListItem{
			DistributorID:   row.DistributorID,
			DistributorCode: row.DistributorCode,
			DistributorName: row.DistributorName,
			TotalProduct:    row.TotalProduct,
		}
		if row.CreatedBy != nil {
			item.CreatedBy = *row.CreatedBy
		}
		if row.CreatedByName != nil {
			item.CreatedByName = *row.CreatedByName
		}
		if row.CreatedAt != nil {
			item.CreatedAt = row.CreatedAt.Format(time.RFC3339)
		}
		if row.UpdatedBy != nil {
			item.UpdatedBy = *row.UpdatedBy
		}
		if row.UpdatedByName != nil {
			item.UpdatedByName = *row.UpdatedByName
		}
		if row.UpdatedAt != nil {
			item.UpdatedAt = row.UpdatedAt.Format(time.RFC3339)
		}
		items = append(items, item)
	}

	return items, total, lastPage, nil
}

func (s *productMappingServiceImpl) Detail(dataFilter entity.ProductMappingDetailQueryFilter) (entity.ProductMappingDetailResponse, int, int, error) {
	normalizeProductMappingPagination(&dataFilter.Page, &dataFilter.Limit)

	dist, err := s.DistributorRepository.FindOneByDistributorIdAndCustId(entity.DetailDistributorParams{
		CustId:        dataFilter.CustId,
		DistributorId: int(dataFilter.DistributorId),
	})
	if err != nil {
		return entity.ProductMappingDetailResponse{}, 0, 0, fmt.Errorf("distributor not found")
	}

	rows, total, lastPage, err := s.ProductMappingRepository.FindDetailByDistributor(dataFilter)
	if err != nil {
		log.Error("ProductMappingService, Detail, err:", err.Error())
		return entity.ProductMappingDetailResponse{}, 0, 0, err
	}

	items := make([]entity.ProductMappingDetailItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.ProductMappingDetailItem{
			ProID:         row.ProID,
			ParentProID:   row.ParentProID,
			ParentProCode: row.ParentProCode,
			ParentProName: row.ParentProName,
			ProCode:       row.ProCode,
			ProName:       row.ProName,
			LargestUOM:    row.LargestUOM,
			MiddleUOM:     row.MiddleUOM,
			SmallestUOM:   row.SmallestUOM,
		})
	}

	response := entity.ProductMappingDetailResponse{
		DistributorCode: dist.DistributorCode,
		DistributorName: dist.DistributorName,
		TotalProduct:    total,
		Items:           items,
	}

	return response, total, lastPage, nil
}

func (s *productMappingServiceImpl) Update(proID int64, req entity.ProductMappingUpdateRequest, principalCustID string, updatedBy int64) error {
	mapping, err := s.ProductMappingRepository.FindOneByProIDAndPrincipal(proID, principalCustID)
	if err != nil {
		return fmt.Errorf("product mapping not found")
	}

	proCode := strings.TrimSpace(req.ProCode)
	proName := strings.TrimSpace(req.ProName)
	unitID1 := strings.TrimSpace(req.UnitID1)
	unitID2 := strings.TrimSpace(req.UnitID2)
	unitID3 := strings.TrimSpace(req.UnitID3)

	if err := validateProductMappingUOM(unitID1, "largest"); err != nil {
		return err
	}
	if unitID2 != "" {
		if err := validateProductMappingUOM(unitID2, "middle"); err != nil {
			return err
		}
	}
	if unitID3 != "" {
		if err := validateProductMappingUOM(unitID3, "smallest"); err != nil {
			return err
		}
	}
	if err := validateDistinctProductMappingUOM(unitID1, unitID2, unitID3); err != nil {
		return err
	}

	dbUnitID1, dbUnitID2, dbUnitID3 := mapProductMappingUOMToProductColumns(unitID1, unitID2, unitID3)

	exists, err := s.ProductMappingRepository.ExistsProName(mapping.CustID, proName, proID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("product name already exist")
	}

	if proCode != "" {
		exists, err = s.ProductMappingRepository.ExistsProCode(mapping.CustID, proCode, proID)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("product code already exist")
		}
	}

	return s.ProductMappingRepository.UpdateMapping(context.Background(), proID, proCode, proName, dbUnitID1, dbUnitID2, dbUnitID3, updatedBy)
}

func (s *productMappingServiceImpl) Delete(proID int64, principalCustID string, deletedBy int64) error {
	if _, err := s.ProductMappingRepository.FindOneByProIDAndPrincipal(proID, principalCustID); err != nil {
		return fmt.Errorf("product mapping not found")
	}
	return s.ProductMappingRepository.SoftDelete(context.Background(), proID, deletedBy)
}

func (s *productMappingServiceImpl) DownloadTemplate() (*bytes.Buffer, string, string, error) {
	buf, err := createProductMappingTemplateXLSX()
	if err != nil {
		return nil, "", "", err
	}
	return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "product_mapping_template.xlsx", nil
}

func createProductMappingTemplateXLSX() (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Template"
	index, _ := f.NewSheet(sheetName)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
	})
	widths := make([]int, len(productMappingTemplateHeaders))

	for i, header := range productMappingTemplateHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheetName, cell, header)
		_ = f.SetCellStyle(sheetName, cell, cell, headerStyle)
		widths[i] = utf8.RuneCountInString(header)
	}

	insSheet := "Instruction"
	_, _ = f.NewSheet(insSheet)
	insHeaders := []string{"Field", "Instruction"}
	for i, h := range insHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(insSheet, cell, h)
		_ = f.SetCellStyle(insSheet, cell, cell, headerStyle)
	}
	for rowIdx, row := range productMappingTemplateInstructions() {
		_ = f.SetCellValue(insSheet, "A"+strconv.Itoa(rowIdx+2), row[0])
		_ = f.SetCellValue(insSheet, "B"+strconv.Itoa(rowIdx+2), row[1])
	}
	_ = f.SetColWidth(insSheet, "A", "A", 28)
	_ = f.SetColWidth(insSheet, "B", "B", 80)

	f.SetActiveSheet(index)
	_ = f.DeleteSheet("Sheet1")

	for i := range widths {
		if colName, err := excelize.ColumnNumberToName(i + 1); err == nil {
			width := float64(widths[i] + 4)
			if width < 18 {
				width = 18
			}
			_ = f.SetColWidth(sheetName, colName, colName, width)
		}
	}

	return f.WriteToBuffer()
}

func productMappingTemplateInstructions() [][]string {
	return [][]string{
		{"Distributor Code", "Mandatory. Must exist in distributor master under principal account."},
		{"Principal Product Code", "Mandatory. Must exist in principal product master and cannot be mapped more than once per distributor."},
		{"Distributor Product Code", "Optional. Alphanumeric/special characters, maximum 10 characters, unique per distributor (case-insensitive trim)."},
		{"Distributor Product Name", "Mandatory. Maximum 20 characters, unique per distributor (case-insensitive trim)."},
		{"Largest Unit", "Mandatory. Alphanumeric, maximum 5 characters."},
		{"Middle Unit", "Optional. Alphanumeric, maximum 5 characters."},
		{"Smallest Unit", "Optional. Alphanumeric, maximum 5 characters."},
		{"UOM Uniqueness", "Largest/Middle/Smallest unit values must be unique within one row (case-insensitive)."},
		{"Duplicate Row Rule", "Combination of principal_product_code + distributor_product_code + distributor_product_name must be unique in one import file."},
	}
}

func (s *productMappingServiceImpl) Import(req entity.ProductMappingImportRequest, principalCustID string, createdBy int64) (entity.ProductMappingImportResponse, error) {
	fileURL := strings.TrimSpace(req.URL)
	if fileURL == "" {
		fileURL = strings.TrimSpace(req.FileURL)
	}

	response := entity.ProductMappingImportResponse{
		URL:         fileURL,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}

	rows, err := downloadProductMappingRows(fileURL)
	if err != nil {
		return response, err
	}

	candidates := make([]productMappingImportCandidate, 0)
	failedReasons := make([]string, 0)
	total := 0
	seenCombinationInFile := make(map[string]int)

	headerMap := map[string]int{}
	if len(rows) > 0 {
		for idx, header := range rows[0] {
			headerMap[normalizeProductMappingHeader(header)] = idx
		}
	}

	for rowIndex, row := range rows {
		if rowIndex == 0 {
			continue
		}
		if isProductMappingRowEmpty(row) {
			continue
		}
		total++

		candidate, rowErrs := s.validateImportRow(rowIndex+1, row, headerMap, principalCustID)
		if len(rowErrs) > 0 {
			failedReasons = append(failedReasons, rowErrs...)
			continue
		}

		combinationKey := buildProductMappingImportCombinationKey(candidate.parentProCode, candidate.proCode, candidate.proName)
		if firstRow, exists := seenCombinationInFile[combinationKey]; exists {
			failedReasons = append(failedReasons, fmt.Sprintf("row %d: duplicate combination of parent_pro_code, pro_code, and pro_name in import file (already exists on row %d)", rowIndex+1, firstRow))
			continue
		}
		seenCombinationInFile[combinationKey] = rowIndex + 1

		candidates = append(candidates, candidate)
	}

	if len(failedReasons) > 0 {
		failedRows := total - len(candidates)
		response.TotalRow = total
		response.FailedRow = failedRows
		response.SuccessRow = len(candidates)
		response.FailedReasons = failedReasons
		if req.Validate {
			return response, fmt.Errorf("import has validation errors")
		}
		return response, nil
	}

	if !req.Validate {
		response.TotalRow = total
		response.SuccessRow = total
		response.FailedRow = 0
		return response, nil
	}

	ctx := context.Background()
	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		for _, candidate := range candidates {
			productReplica := candidate.parentProduct
			productReplica.ProductId = 0
			productReplica.ParentProId = int(candidate.parentProduct.ProductId)
			productReplica.CustId = candidate.distCustID
			productReplica.DistributorID = &candidate.distributorID
			productReplica.Level = 1
			productReplica.Origin = productMappingOrigin
			productReplica.AssignerUserID = &createdBy
			productReplica.IsProductMapping = true
			productReplica.ProductCode = candidate.proCode
			productReplica.ProductName = candidate.proName
			productReplica.UnitId1, productReplica.UnitId2, productReplica.UnitId3 = mapProductMappingUOMToProductColumns(
				candidate.largestUOM,
				candidate.middleUOM,
				candidate.smallestUOM,
			)
			productReplica.IsActive = true
			productReplica.IsDel = false
			now := time.Now()
			productReplica.CreatedAt = &now
			productReplica.UpdatedAt = &now
			productReplica.CreatedBy = &createdBy
			productReplica.UpdatedBy = &createdBy

			if _, err := s.ProductRepository.Store(txCtx, productReplica); err != nil {
				return fmt.Errorf("row %d: failed to save product mapping: %v", candidate.rowNumber, err)
			}
		}
		return nil
	})

	response.TotalRow = total
	response.FailedReasons = failedReasons
	if err != nil {
		response.SuccessRow = 0
		response.FailedRow = total
		return response, err
	}

	response.SuccessRow = total
	response.FailedRow = 0
	return response, nil
}

func (s *productMappingServiceImpl) validateImportRow(rowNumber int, row []string, headerMap map[string]int, principalCustID string) (productMappingImportCandidate, []string) {
	var errs []string
	getValue := func(keys ...string) string {
		for _, key := range keys {
			if idx, ok := headerMap[key]; ok && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
		}
		return ""
	}

	distributorCode := getValue("distributor_code")
	parentProCode := getValue("parent_pro_code", "principal_product_code")
	proCode := getValue("pro_code", "distributor_product_code")
	proName := getValue("pro_name", "distributor_product_name")
	largestUOM := getValue("largest_uom", "largest_unit")
	middleUOM := getValue("middle_uom", "middle_unit")
	smallestUOM := getValue("smallest_uom", "smallest_unit")

	if distributorCode == "" {
		errs = append(errs, fmt.Sprintf("row %d: distributor_code is mandatory", rowNumber))
	}
	if parentProCode == "" {
		errs = append(errs, fmt.Sprintf("row %d: parent_pro_code is mandatory", rowNumber))
	}
	if proName == "" {
		errs = append(errs, fmt.Sprintf("row %d: pro_name is mandatory", rowNumber))
	}
	if largestUOM == "" {
		errs = append(errs, fmt.Sprintf("row %d: largest_unit is mandatory", rowNumber))
	}

	if len(errs) > 0 {
		return productMappingImportCandidate{}, errs
	}

	if utf8.RuneCountInString(proCode) > 10 {
		errs = append(errs, fmt.Sprintf("row %d: product code exceeds maximum length", rowNumber))
	}
	if proCode != "" && !productMappingProCodePattern.MatchString(proCode) {
		errs = append(errs, fmt.Sprintf("row %d: invalid product code format", rowNumber))
	}
	if utf8.RuneCountInString(proName) > 20 {
		errs = append(errs, fmt.Sprintf("row %d: product name exceeds maximum length", rowNumber))
	}
	if err := validateProductMappingUOM(largestUOM, "largest"); err != nil {
		errs = append(errs, fmt.Sprintf("row %d: %s", rowNumber, err.Error()))
	}
	if middleUOM != "" {
		if err := validateProductMappingUOM(middleUOM, "middle"); err != nil {
			errs = append(errs, fmt.Sprintf("row %d: %s", rowNumber, err.Error()))
		}
	}
	if smallestUOM != "" {
		if err := validateProductMappingUOM(smallestUOM, "smallest"); err != nil {
			errs = append(errs, fmt.Sprintf("row %d: %s", rowNumber, err.Error()))
		}
	}
	if err := validateDistinctProductMappingUOM(largestUOM, middleUOM, smallestUOM); err != nil {
		errs = append(errs, fmt.Sprintf("row %d: %s", rowNumber, err.Error()))
	}

	dist, err := s.DistributorRepository.FindOneByParentCustIdAndDistributorCode(principalCustID, distributorCode)
	if err != nil || dist.DistributorId == 0 {
		errs = append(errs, fmt.Sprintf("row %d: distributor %s not found", rowNumber, distributorCode))
	}

	parentProduct, err := s.ProductRepository.FindOneByProductCodeAndCustId(parentProCode, principalCustID)
	if err != nil || parentProduct.ProductId == 0 {
		errs = append(errs, fmt.Sprintf("row %d: parent product %s not found", rowNumber, parentProCode))
	}

	if dist.DistributorCustID != "" {
		existsParentMapping, err := s.ProductMappingRepository.ExistsParentMappingByDistributor(dist.DistributorId, int64(parentProduct.ProductId))
		if err != nil {
			errs = append(errs, fmt.Sprintf("row %d: failed to validate parent_pro_code mapping", rowNumber))
		} else if existsParentMapping {
			errs = append(errs, fmt.Sprintf("row %d: Parent Product Code Already Exist", rowNumber))
		}

		if proCode != "" {
			exists, err := s.ProductMappingRepository.ExistsProCode(dist.DistributorCustID, proCode, 0)
			if err != nil {
				errs = append(errs, fmt.Sprintf("row %d: failed to validate product code", rowNumber))
			} else if exists {
				errs = append(errs, fmt.Sprintf("row %d: Product Code Already Exist", rowNumber))
			}
		}
		exists, err := s.ProductMappingRepository.ExistsProName(dist.DistributorCustID, proName, 0)
		if err != nil {
			errs = append(errs, fmt.Sprintf("row %d: failed to validate product name", rowNumber))
		} else if exists {
			errs = append(errs, fmt.Sprintf("row %d: Product Name Already Exist", rowNumber))
		}
	}

	if len(errs) > 0 {
		return productMappingImportCandidate{}, errs
	}

	return productMappingImportCandidate{
		rowNumber:       rowNumber,
		distributorCode: distributorCode,
		parentProCode:   parentProCode,
		proCode:         proCode,
		proName:         proName,
		largestUOM:      largestUOM,
		middleUOM:       middleUOM,
		smallestUOM:     smallestUOM,
		distCustID:      dist.DistributorCustID,
		distributorID:   dist.DistributorId,
		parentProduct:   parentProduct,
	}, nil
}

func downloadProductMappingRows(fileURL string) ([][]string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status %d", resp.StatusCode)
	}

	f, err := excelize.OpenReader(resp.Body)
	if err != nil {
		return nil, err
	}

	sheetName := "Template"
	if idx, err := f.GetSheetIndex(sheetName); err != nil || idx < 0 {
		sheets := f.GetSheetList()
		if len(sheets) == 0 {
			return nil, fmt.Errorf("excel file has no sheets")
		}
		sheetName = sheets[0]
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func normalizeProductMappingHeader(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, " ", "_")
	return value
}

func isProductMappingRowEmpty(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func normalizeProductMappingPagination(page, limit *int) {
	if *page <= 0 {
		*page = 1
	}
	if *limit <= 0 {
		*limit = 10
	}
}

func validateProductMappingUOM(value, label string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s uom is mandatory", label)
	}
	if utf8.RuneCountInString(value) > 5 {
		return fmt.Errorf("%s uom exceeds maximum length", label)
	}
	if !productMappingUOMPattern.MatchString(value) {
		return fmt.Errorf("%s uom must be alphanumeric", label)
	}
	return nil
}

func validateDistinctProductMappingUOM(largest, middle, smallest string) error {
	if middle != "" && isSameProductMappingUOM(largest, middle) {
		return errors.New(constant.ProductMappingDuplicateUOMErrorMsg)
	}
	if smallest != "" && isSameProductMappingUOM(largest, smallest) {
		return errors.New(constant.ProductMappingDuplicateUOMErrorMsg)
	}
	if middle != "" && smallest != "" && isSameProductMappingUOM(middle, smallest) {
		return errors.New(constant.ProductMappingDuplicateUOMErrorMsg)
	}
	return nil
}

func isSameProductMappingUOM(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func mapProductMappingUOMToProductColumns(largest, middle, smallest string) (unitID1, unitID2, unitID3 string) {
	return strings.TrimSpace(smallest), strings.TrimSpace(middle), strings.TrimSpace(largest)
}

func buildProductMappingImportCombinationKey(parentProCode, proCode, proName string) string {
	return strings.ToLower(strings.TrimSpace(parentProCode)) + "|" +
		strings.ToLower(strings.TrimSpace(proCode)) + "|" +
		strings.ToLower(strings.TrimSpace(proName))
}
