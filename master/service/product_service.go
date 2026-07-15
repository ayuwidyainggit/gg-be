package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"strconv"
	"time"

	// "io"
	"archive/zip"
	"regexp"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2/log"
	"github.com/xuri/excelize/v2"
)

type ProductService interface {
	Detail(entity.DetailProductParams) (entity.ProductDetailResponse, error)
	List(entity.ProductQueryFilter, string) (data []entity.ProductResponse, total int, lastPage int, err error)
	ReportList(entity.ProductReportQueryFilter) (data []entity.ProductReportResponse, total int, lastPage int, err error)
	Export(filter entity.ProductQueryFilter) (buffer *bytes.Buffer, contentType string, filename string, err error)
	ExportTemplate(format string) (buffer *bytes.Buffer, contentType string, filename string, err error)
	ExportTemplateUpdate(custId string, format string, fields []string) (*bytes.Buffer, string, string, error)
	ImportProductCSV4(req entity.ImportProductRequest) error
	ImportProductXLSX4(req entity.ImportProductRequest) error
	ImportUpdateXLSX(req entity.ImportProductRequest) error
	ImportUpdateCSV(req entity.ImportProductRequest) error
	ExportImportInstructions(format string) (*bytes.Buffer, string, string, error)
	createInstructionCSV(data []model.ImportInstruction) (*bytes.Buffer, error)
	createInstructionXLSX(data []model.ImportInstruction) (*bytes.Buffer, error)
	mapToProcessedProductRow(custId string, row map[string]string) (entity.ProcessedProductRow, error)
	ReuploadImportUpdateFile(custId string, historyId int64, req entity.ImportRequest) error
	ReuploadImportInsertFile(custId string, historyId int64, req entity.ImportRequest) error
	LookupList(entity.ProductQueryFilter, string) (data []entity.ProductLookupResponse, total int, lastPage int, err error)
	LookupDistPrice(entity.ProductQueryFilter) (data []entity.ProductLookupDistPrice, total int, lastPage int, err error)
	SearchList(entity.ProductQueryFilter, string) (data []entity.ProductSearchResponse, total int, lastPage int, err error)
	Store(entity.CreateProductBody) (entity.ProductResponse, error)
	BulkStore(entity.BulkProductBody) (entity.BulkProductResponse, error)
	Update(int64, entity.UpdateProductRequest) error
	Delete(string, int64, int64) error
	DeleteMultiple(string, []int64, int64) error
	PrincipalList(entity.ProductPrincipalQueryFilter, string) (data []entity.PrincipalLookupResponse, total int, lastPage int, err error)
	CategoryList(entity.ProductCategoryQueryFilter, string) (data []entity.ProductCategoryList, total int, lastPage int, err error)
	BrandList(entity.ProductBrandQueryFilter, string) (data []entity.ProductBrandList, total int, lastPage int, err error)
}

func NewProductService(productRepository repository.ProductRepository) *productServiceImpl {
	return &productServiceImpl{
		ProductRepository: productRepository,
	}
}

type productServiceImpl struct {
	ProductRepository repository.ProductRepository
}

func (service *productServiceImpl) Detail(params entity.DetailProductParams) (response entity.ProductDetailResponse, err error) {
	var product model.Product

	product, err = service.ProductRepository.FindOne(params)
	if err != nil {
		return response, err
	}

	payload, err := json.Marshal(product)
	if err != nil {
		return response, err
	}
	if err = json.Unmarshal(payload, &response); err != nil {
		return response, err
	}

	return response, err
}

func (service *productServiceImpl) List(dataFilter entity.ProductQueryFilter, custId string) (data []entity.ProductResponse, total int, lastPage int, err error) {
	products, total, lastPage, err := service.ProductRepository.FindAll(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range products {
		var vResp entity.ProductResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *productServiceImpl) ReportList(filter entity.ProductReportQueryFilter) (data []entity.ProductReportResponse, total int, lastPage int, err error) {
	return service.ProductRepository.ReportList(filter)
}

func (s *productServiceImpl) ExportImportInstructions(format string) (*bytes.Buffer, string, string, error) {
	instructions, err := s.ProductRepository.GetProductImportInstructions()
	if err != nil {
		return nil, "", "", err
	}

	switch format {
	case "csv":
		buf, err := s.createInstructionCSV(instructions)
		if err != nil {
			return nil, "", "", err
		}
		return buf, "text/csv", "import_instructions_product.csv", nil
	case "xls":
		buf, err := s.createInstructionXLSX(instructions)
		if err != nil {
			return nil, "", "", err
		}
		return buf, "application/vnd.ms-excel", "import_instructions_product.xls", nil
	default:
		buf, err := s.createInstructionXLSX(instructions)
		if err != nil {
			return nil, "", "", err
		}
		return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "import_instructions_product.xlsx", nil
	}
}

func (s *productServiceImpl) createInstructionCSV(data []model.ImportInstruction) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Comma = ';'

	// Header
	writer.Write([]string{"Kolom", "Mandatory", "Keterangan"})

	// Rows
	for _, d := range data {
		writer.Write([]string{
			d.Kolom,
			d.Mandatory,
			safeString(d.Keterangan),
		})
	}
	writer.Flush()
	return buf, nil
}

func (s *productServiceImpl) createInstructionXLSX(data []model.ImportInstruction) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheet := "Instructions"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)

	// Header
	headers := []string{"Kolom", "Mandatory", "Keterangan"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Rows
	for r, d := range data {
		f.SetCellValue(sheet, "A"+strconv.Itoa(r+2), d.Kolom)
		f.SetCellValue(sheet, "B"+strconv.Itoa(r+2), d.Mandatory)
		f.SetCellValue(sheet, "C"+strconv.Itoa(r+2), safeString(d.Keterangan))
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (service *productServiceImpl) Export(filter entity.ProductQueryFilter) (*bytes.Buffer, string, string, error) {
	var contentType, filename string
	var buffer *bytes.Buffer
	products, _, err := service.ProductRepository.FindAllExport(filter, filter.CustId)
	if err != nil {
		log.Info("ProductService, Export, FindAllExport err: %v", err)
		return nil, "", "", err
	}

	switch filter.Format {
	case "csv":
		buffer, err := service.createCSV(products)
		if err != nil {
			return nil, "", "", err
		}
		return buffer, "text/csv", "products.csv", nil
	case "xls":
		buffer, err = service.createXLS(products)
		contentType = "application/vnd.ms-excel"
		filename = "products.xls"
		return buffer, contentType, filename, nil
	default:
		buffer, err := service.createXLSX(products)
		if err != nil {
			return nil, "", "", err
		}
		contentType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		filename := "products.xlsx"
		return buffer, contentType, filename, nil
	}
}

func (service *productServiceImpl) createXLSX(products []model.Product) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Products"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		log.Errorf("createXLSX: Failed to create new sheet: %v", err)
		return nil, err
	}

	headers := []string{
		"pro_name", "pro_code", "bar_code",
		"pro_code_coretax", "pro_name_coretax",
		"pcat_code", "pcat_name",
		"pl_code", "pl_name",
		"brand_code", "brand_name",
		"sbrand1_code", "sbrand1_name",
		"sbrand2_code", "sbrand2_name",
		"flavor_code", "flavor_name",
		"ptype_code", "ptype_name",
		"psize_code", "psize_name",
		"principal_code", "principal_name",
		"sup_code", "sup_name",
		"c_pro_code", "c_pro_name",
		"parent_pro_id", "is_active",
		"is_batch", "is_exp_date", "pro_status",
		"unit_id3", "unit_name3", "unit_id2", "unit_name2", "unit_id1", "unit_name1",
		"conv_unit3", "conv_unit2",
		"purch_price3", "purch_price2", "purch_price1",
		"sell_price3", "sell_price2", "sell_price1",
		"weight3", "length3", "width3", "height3", "volume3",
		"weight2", "length2", "width2", "height2", "volume2",
		"weight1", "length1", "width1", "height1", "volume1",
		"saf_stock_qty", "saf_stock_unit_id",
		"min_stock_qty", "min_stock_unit_id",
		"vat", "vat_bg", "vat_lg_purch", "vat_lg_sell", "cogs",
		"excise_rate", "excise_tax",
	}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, MapHeaderToWeb(header))
	}

	for i, rowData := range products {
		var vResp entity.ProductExportResponse
		if err := structs.Automapper(rowData, &vResp); err != nil {
			log.Errorf("createXLSX: Automapper err: %v", err)
			return nil, err
		}
		currentRow := i + 2
		values := []interface{}{
			vResp.ProductName, vResp.ProductCode, strPtrProd(vResp.BarCode),
			strPtrProd(vResp.ProCodeCoreTax), strPtrProd(vResp.ProNameCoreTax),
			vResp.PCatCode, vResp.PCatName,
			vResp.ProductLineCode, vResp.ProductLineName,
			vResp.BrandCode, vResp.BrandName,
			vResp.Sbrand1Code, vResp.Sbrand1Name,
			vResp.Sbrand2Code, vResp.Sbrand2Name,
			vResp.FlavorCode, vResp.FlavorName,
			vResp.PTypeCode, vResp.PTypeName,
			vResp.PSizeCode, vResp.PSizeName,
			vResp.PrincipalCode, vResp.PrincipalName,
			vResp.SupCode, vResp.SupName,
			vResp.CProCode, vResp.CProName,
			vResp.ParentProId, vResp.IsActive,
			vResp.IsBatch, vResp.IsExpDate,
			mapProductStatusToString(intPtrProd(vResp.ProStatus)),
			vResp.UnitId3, vResp.UnitName3, vResp.UnitId2, vResp.UnitName2, vResp.UnitId1, vResp.UnitName1,
			vResp.ConvUnit3, vResp.ConvUnit2,
			vResp.PurchPrice3, vResp.PurchPrice2, vResp.PurchPrice1,
			vResp.SellPrice3, vResp.SellPrice2, vResp.SellPrice1,
			floatPtrProd(vResp.Weight3), floatPtrProd(vResp.Length3), floatPtrProd(vResp.Width3), floatPtrProd(vResp.Height3), floatPtrProd(vResp.Volume3),
			floatPtrProd(vResp.Weight2), floatPtrProd(vResp.Length2), floatPtrProd(vResp.Width2), floatPtrProd(vResp.Height2), floatPtrProd(vResp.Volume2),
			floatPtrProd(vResp.Weight1), floatPtrProd(vResp.Length1), floatPtrProd(vResp.Width1), floatPtrProd(vResp.Height1), floatPtrProd(vResp.Volume1),
			vResp.SafStockQty, strPtrProd(vResp.SafStockUnitId),
			vResp.MinStockQty, strPtrProd(vResp.MinStockUnitId),
			floatPtrProd(vResp.Vat), floatPtrProd(vResp.VatBg), floatPtrProd(vResp.VatLgPurch), floatPtrProd(vResp.VatLgSell), floatPtrProd(vResp.Cogs),
			vResp.ExciseRate, vResp.ExciseTax,
		}

		for col, val := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, currentRow)
			f.SetCellValue(sheetName, cell, val)
		}
	}

	f.SetActiveSheet(index)
	if err := f.DeleteSheet("Sheet1"); err != nil {
		log.Warnf("createXLSX: Could not delete default sheet: %v", err)
	}

	return f.WriteToBuffer()
}

func (service *productServiceImpl) createXLS(products []model.Product) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Products"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		log.Errorf("createXLS: Failed to create new sheet: %v", err)
		return nil, err
	}

	headers := []string{
		"pro_name", "pro_code", "bar_code",
		"pro_code_coretax", "pro_name_coretax",
		"pcat_code", "pcat_name",
		"pl_code", "pl_name",
		"brand_code", "brand_name",
		"sbrand1_code", "sbrand1_name",
		"sbrand2_code", "sbrand2_name",
		"flavor_code", "flavor_name",
		"ptype_code", "ptype_name",
		"psize_code", "psize_name",
		"principal_code", "principal_name",
		"sup_code", "sup_name",
		"c_pro_code", "c_pro_name",
		"parent_pro_id", "is_active",
		"is_batch", "is_exp_date", "pro_status",
		"unit_id3", "unit_name3", "unit_id2", "unit_name2", "unit_id1", "unit_name1",
		"conv_unit3", "conv_unit2",
		"purch_price3", "purch_price2", "purch_price1",
		"sell_price3", "sell_price2", "sell_price1",
		"weight3", "length3", "width3", "height3", "volume3",
		"weight2", "length2", "width2", "height2", "volume2",
		"weight1", "length1", "width1", "height1", "volume1",
		"saf_stock_qty", "saf_stock_unit_id",
		"min_stock_qty", "min_stock_unit_id",
		"vat", "vat_bg", "vat_lg_purch", "vat_lg_sell", "cogs",
		"excise_rate", "excise_tax",
	}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, MapHeaderToWeb(header))

		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
		})
		f.SetCellStyle(sheetName, cell, cell, style)
	}

	for i, rowData := range products {
		var vResp entity.ProductExportResponse
		if err := structs.Automapper(rowData, &vResp); err != nil {
			log.Errorf("createXLS: Automapper err: %v", err)
			return nil, err
		}
		currentRow := i + 2
		values := []interface{}{
			vResp.ProductName, vResp.ProductCode, strPtrProd(vResp.BarCode),
			strPtrProd(vResp.ProCodeCoreTax), strPtrProd(vResp.ProNameCoreTax),
			vResp.PCatCode, vResp.PCatName,
			vResp.ProductLineCode, vResp.ProductLineName,
			vResp.BrandCode, vResp.BrandName,
			vResp.Sbrand1Code, vResp.Sbrand1Name,
			vResp.Sbrand2Code, vResp.Sbrand2Name,
			vResp.FlavorCode, vResp.FlavorName,
			vResp.PTypeCode, vResp.PTypeName,
			vResp.PSizeCode, vResp.PSizeName,
			vResp.PrincipalCode, vResp.PrincipalName,
			vResp.SupCode, vResp.SupName,
			vResp.CProCode, vResp.CProName,
			vResp.ParentProId, vResp.IsActive,
			vResp.IsBatch, vResp.IsExpDate,
			mapProductStatusToString(intPtrProd(vResp.ProStatus)),
			vResp.UnitId3, vResp.UnitName3, vResp.UnitId2, vResp.UnitName2, vResp.UnitId1, vResp.UnitName1,
			vResp.ConvUnit3, vResp.ConvUnit2,
			vResp.PurchPrice3, vResp.PurchPrice2, vResp.PurchPrice1,
			vResp.SellPrice3, vResp.SellPrice2, vResp.SellPrice1,
			floatPtrProd(vResp.Weight3), floatPtrProd(vResp.Length3), floatPtrProd(vResp.Width3), floatPtrProd(vResp.Height3), floatPtrProd(vResp.Volume3),
			floatPtrProd(vResp.Weight2), floatPtrProd(vResp.Length2), floatPtrProd(vResp.Width2), floatPtrProd(vResp.Height2), floatPtrProd(vResp.Volume2),
			floatPtrProd(vResp.Weight1), floatPtrProd(vResp.Length1), floatPtrProd(vResp.Width1), floatPtrProd(vResp.Height1), floatPtrProd(vResp.Volume1),
			vResp.SafStockQty, strPtrProd(vResp.SafStockUnitId),
			vResp.MinStockQty, strPtrProd(vResp.MinStockUnitId),
			floatPtrProd(vResp.Vat), floatPtrProd(vResp.VatBg), floatPtrProd(vResp.VatLgPurch), floatPtrProd(vResp.VatLgSell), floatPtrProd(vResp.Cogs),
			vResp.ExciseRate, vResp.ExciseTax,
		}

		for col, val := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, currentRow)
			f.SetCellValue(sheetName, cell, val)
		}
	}

	f.SetActiveSheet(index)
	if err := f.DeleteSheet("Sheet1"); err != nil {
		log.Warnf("createXLS: Could not delete default sheet: %v", err)
	}

	return f.WriteToBuffer()
}

func (service *productServiceImpl) createCSV(products []model.Product) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)
	writer.Comma = ';'

	headers := []string{
		"pro_name", "pro_code", "bar_code",
		"pro_code_coretax", "pro_name_coretax",
		"pcat_code", "pcat_name",
		"pl_code", "pl_name",
		"brand_code", "brand_name",
		"sbrand1_code", "sbrand1_name",
		"sbrand2_code", "sbrand2_name",
		"flavor_code", "flavor_name",
		"ptype_code", "ptype_name",
		"psize_code", "psize_name",
		"principal_code", "principal_name",
		"sup_code", "sup_name",
		"c_pro_code", "c_pro_name",
		"parent_pro_id", "is_active",
		"is_batch", "is_exp_date", "pro_status",
		"unit_id3", "unit_name3", "unit_id2", "unit_name2", "unit_id1", "unit_name1",
		"conv_unit3", "conv_unit2",
		"purch_price3", "purch_price2", "purch_price1",
		"sell_price3", "sell_price2", "sell_price1",
		"weight3", "length3", "width3", "height3", "volume3",
		"weight2", "length2", "width2", "height2", "volume2",
		"weight1", "length1", "width1", "height1", "volume1",
		"saf_stock_qty", "saf_stock_unit_id",
		"min_stock_qty", "min_stock_unit_id",
		"vat", "vat_bg", "vat_lg_purch", "vat_lg_sell", "cogs",
		"excise_rate", "excise_tax",
	}

	webHeaders := make([]string, len(headers))
	for i, h := range headers {
		webHeaders[i] = MapHeaderToWeb(h) // ubah nama DB ke nama web
	}

	if err := writer.Write(webHeaders); err != nil {
		log.Errorf("createCSV: Failed to write header: %v", err)
		return nil, err
	}

	for _, rowData := range products {
		var vResp entity.ProductExportResponse
		if err := structs.Automapper(rowData, &vResp); err != nil {
			log.Errorf("createCSV: Automapper err: %v", err)
			return nil, err
		}

		record := []string{
			vResp.ProductName, vResp.ProductCode, strPtrProd(vResp.BarCode),
			strPtrProd(vResp.ProCodeCoreTax), strPtrProd(vResp.ProNameCoreTax),
			vResp.PCatCode, vResp.PCatName,
			vResp.ProductLineCode, vResp.ProductLineName,
			vResp.BrandCode, vResp.BrandName,
			vResp.Sbrand1Code, vResp.Sbrand1Name,
			vResp.Sbrand2Code, vResp.Sbrand2Name,
			vResp.FlavorCode, vResp.FlavorName,
			vResp.PTypeCode, vResp.PTypeName,
			vResp.PSizeCode, vResp.PSizeName,
			vResp.PrincipalCode, vResp.PrincipalName,
			vResp.SupCode, vResp.SupName,
			vResp.CProCode, vResp.CProName,
			strconv.Itoa(vResp.ParentProId), strconv.FormatBool(vResp.IsActive),
			strconv.FormatBool(vResp.IsBatch), strconv.FormatBool(vResp.IsExpDate),
			mapProductStatusToString(intPtrProd(vResp.ProStatus)),
			vResp.UnitId3, vResp.UnitName3, vResp.UnitId2, vResp.UnitName2, vResp.UnitId1, vResp.UnitName1,
			strconv.FormatFloat(float64(vResp.ConvUnit3), 'f', -1, 32), strconv.FormatFloat(float64(vResp.ConvUnit2), 'f', -1, 32),
			strconv.FormatFloat(vResp.PurchPrice3, 'f', -1, 64), strconv.FormatFloat(vResp.PurchPrice2, 'f', -1, 64), strconv.FormatFloat(vResp.PurchPrice1, 'f', -1, 64),
			strconv.FormatFloat(vResp.SellPrice3, 'f', -1, 64), strconv.FormatFloat(vResp.SellPrice2, 'f', -1, 64), strconv.FormatFloat(vResp.SellPrice1, 'f', -1, 64),
			floatPtrStrProd(vResp.Weight3), floatPtrStrProd(vResp.Length3), floatPtrStrProd(vResp.Width3), floatPtrStrProd(vResp.Height3), floatPtrStrProd(vResp.Volume3),
			floatPtrStrProd(vResp.Weight2), floatPtrStrProd(vResp.Length2), floatPtrStrProd(vResp.Width2), floatPtrStrProd(vResp.Height2), floatPtrStrProd(vResp.Volume2),
			floatPtrStrProd(vResp.Weight1), floatPtrStrProd(vResp.Length1), floatPtrStrProd(vResp.Width1), floatPtrStrProd(vResp.Height1), floatPtrStrProd(vResp.Volume1),
			strconv.FormatFloat(vResp.SafStockQty, 'f', -1, 64), strPtrProd(vResp.SafStockUnitId),
			strconv.FormatFloat(vResp.MinStockQty, 'f', -1, 64), strPtrProd(vResp.MinStockUnitId),
			floatPtrStrProd(vResp.Vat), floatPtrStrProd(vResp.VatBg), floatPtrStrProd(vResp.VatLgPurch), floatPtrStrProd(vResp.VatLgSell), floatPtrStrProd(vResp.Cogs),
			strconv.FormatFloat(vResp.ExciseRate, 'f', -1, 64), strconv.FormatFloat(vResp.ExciseTax, 'f', -1, 64),
		}

		if err := writer.Write(record); err != nil {
			log.Errorf("createCSV: Failed to write record: %v", err)
			return nil, err
		}
	}

	writer.Flush()

	return buffer, nil
}

func strPtrProd(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func floatPtrProd(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func intPtrProd(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func int64PtrProd(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

func timePtrProd(v *time.Time) string {
	if v == nil {
		return ""
	}
	return v.Format("2006-01-02 15:04:05")
}

func floatPtrStrProd(v *float64) string {
	if v == nil {
		return ""
	}
	return strconv.FormatFloat(*v, 'f', -1, 64)
}

func boolStrProd(v bool) string {
	return strconv.FormatBool(v)
}

func int64PtrStrProd(v *int64) string {
	if v == nil {
		return ""
	}
	return strconv.FormatInt(*v, 10)
}

func timePtrStrProd(v *time.Time) string {
	if v == nil {
		return ""
	}
	return v.Format("2006-01-02 15:04:05")
}

func Int64PtrProd2(v int64) *int64 {
	return &v
}

func TimePtrProd2(v time.Time) *time.Time {
	return &v
}

var HeaderMapping = map[string]string{
	"pro_id":   "Product ID",
	"pro_code": "Product Code",
	"bar_code": "Barcode",
	"pro_name": "Product Name",

	// Category
	"pcat_id":   "Product Category ID",
	"pcat_code": "Product Category Code",
	"pcat_name": "Product Category",

	// Product Line
	"pl_id":   "Product Line ID",
	"pl_code": "Product Line Code",
	"pl_name": "Product Line",

	// Brand
	"brand_id":   "Brand ID",
	"brand_code": "Brand Code",
	"brand_name": "Brand",

	// Sub Brand 1
	"sbrand1_id":   "Sub Brand ID",
	"sbrand1_code": "Sub Brand Code",
	"sbrand1_name": "Sub Brand",

	// Sub Brand 2
	"sbrand2_id":   "Sub Brand 2 ID",
	"sbrand2_code": "Sub Brand 2 Code",
	"sbrand2_name": "Sub Brand 2",

	// Flavor
	"flavor_id":   "Flavor ID",
	"flavor_code": "Flavor Code",
	"flavor_name": "Flavor",

	// Pack Type
	"ptype_id":   "Pack Type ID",
	"ptype_code": "Pack Type Code",
	"ptype_name": "Pack Type",

	// Pack Size
	"psize_id":   "Pack Size ID",
	"psize_code": "Pack Size Code",
	"psize_name": "Pack Size",

	// Supplier
	"sup_id":   "Supplier ID",
	"sup_code": "Supplier Code",
	"sup_name": "Supplier",

	// Principal
	"principal_id":   "Principal ID",
	"principal_code": "Principal Code",
	"principal_name": "Principal",

	// Consumer Product
	"c_pro_id":   "Consumer Product ID",
	"c_pro_code": "Consumer Product Code",
	"c_pro_name": "Consumer Product",

	// Flags
	"is_main_pro": "Main Product",
	"sort_no":     "Sort",
	"item_no":     "Item",

	// Units
	"unit_id1": "Smallest Unit Code",
	"unit_id2": "Middle Unit Code",
	"unit_id3": "Largest Unit Code",
	"unit_id4": "Extra Unit Code 1",
	"unit_id5": "Extra Unit Code 2",

	"unit_name1": "Smallest Unit",
	"unit_name2": "Middle Unit",
	"unit_name3": "Largest Unit",

	// unit
	"unit_id":   "Unit Code",
	"unit_name": "Unit Name",

	"conv_unit2": "Conversion Middle Unit",
	"conv_unit3": "Conversion Largest Unit",
	"conv_unit4": "Extra Unit Conversion 1",
	"conv_unit5": "Extra Unit Conversion 2",

	// Flags
	"is_batch":    "Batch",
	"is_exp_date": "Expired Date",

	// Dimensions
	"length": "Length (cm)",
	"width":  "Width (cm)",
	"height": "Height (cm)",
	"weight": "Weight (gram)",
	"volume": "Volume (cm3)",

	// Parent & New
	"parent_pro_id": "Parent Product",
	"is_new_pro":    "New Product",

	// Purchase Prices
	"purch_price1": "Purchase Price Smallest Unit",
	"purch_price2": "Purchase Price Middle Unit",
	"purch_price3": "Purchase Price Largest Unit",
	"purch_price4": "Extra Unit 1 Purchase Price",
	"purch_price5": "Extra Unit 2 Purchase Price",

	// Sell Prices
	"sell_price1": "Sell Price Smallest Unit",
	"sell_price2": "Sell Price Middle Unit",
	"sell_price3": "Sell Price Largest Unit",
	"sell_price4": "Extra Unit 1 Sell Price",
	"sell_price5": "Extra Unit 2 Sell Price",

	// Dimension per unit
	"length1": "Dimension Smallest Unit Length (cm)",
	"length2": "Dimension Middle Unit Length (cm)",
	"length3": "Dimension Largest Unit Length (cm)",
	"length4": "Dimension Extra Unit 1 Length (cm)",
	"length5": "Dimension Extra Unit 2 Length (cm)",

	"width1": "Dimension Smallest Unit Width (cm)",
	"width2": "Dimension Middle Unit Width (cm)",
	"width3": "Dimension Largest Unit Width (cm)",
	"width4": "Dimension Extra Unit 1 Width (cm)",
	"width5": "Dimension Extra Unit 2 Width (cm)",

	"height1": "Dimension Smallest Unit Height (cm)",
	"height2": "Dimension Middle Unit Height (cm)",
	"height3": "Dimension Largest Unit Height (cm)",
	"height4": "Dimension Extra Unit 1 Height (cm)",
	"height5": "Dimension Extra Unit 2 Height (cm)",

	"weight1": "Dimension Smallest Unit Weight (gram)",
	"weight2": "Dimension Middle Unit Weight (gram)",
	"weight3": "Dimension Largest Unit Weight (gram)",
	"weight4": "Dimension Extra Unit 1 Weight (gram)",
	"weight5": "Dimension Extra Unit 2 Weight (gram)",

	"volume1": "Dimension Smallest Unit Volume (cm3)",
	"volume2": "Dimension Middle Unit Volume (cm3)",
	"volume3": "Dimension Largest Unit Volume (cm3)",
	"volume4": "Dimension Extra Unit 1 Volume (cm3)",
	"volume5": "Dimension Extra Unit 2 Volume (cm3)",

	// Stock
	"saf_stock_qty":       "Safety Stock",
	"saf_stock_unit_id":   "Safety Stock Unit Code",
	"saf_stock_unit_name": "Safety Stock Unit Name",
	"min_stock_qty":       "Min Stock",
	"min_stock_unit_id":   "Min Stock Unit Code",
	"min_stock_unit_name": "Min Stock Unit Name",

	// Tax & Status
	"excise_rate": "Tarif Cukai",
	"excise_tax":  "Tarif Pajak (%)",
	"is_active":   "Status Active",
	"is_del":      "Status Deleted",

	// Misc
	"image_url": "Image Link",

	"vat":              "PPN (%)",
	"vat_bg":           "PPN DP (%)",
	"vat_lg_purch":     "PPN BM Dist (%)",
	"vat_lg_sell":      "PPN BM Toko (%)",
	"cogs":             "HPP",
	"pro_status":       "Product Status",
	"pro_code_coretax": "Referral Code (Coretax)",
	"pro_name_coretax": "Description",

	// Status temp
	"status_insert": "Status Data",
	"created_at":    "Created Date",
	"history_id":    "History ID",
	"cust_id":       "Customer ID",
	"error_message": "Error Message",
}

var headerProductMaxLength = map[string]int{
	"pro_name":   100,
	"pro_code":   20,
	"bar_code":   50,
	"brand_code": 5, "brand_name": 100,
	"flavor_code": 5, "flavor_name": 100,
	"sbrand1_code": 5, "sbrand1_name": 100,
	"sbrand2_code": 5, "sbrand2_name": 100,
	"pl_code": 5, "pl_name": 100,
	"pcat_code": 5, "pcat_name": 100,
	"ptype_code": 5, "ptype_name": 50,
	"psize_code": 5, "psize_name": 100,
	"principal_code": 8, "principal_name": 150,
	"sup_code": 20, "sup_name": 150,
	"c_pro_code": 5, "c_pro_name": 150,
	"unit_id1": 5, "unit_id2": 5, "unit_id3": 5,
	"unit_name1": 50, "unit_name2": 50, "unit_name3": 50,
	"saf_stock_unit_id": 10, "min_stock_unit_id": 10,
	"pro_code_coretax": 10,
	"weight1":          10,
	"weight2":          10,
	"weight3":          10,
	"weight4":          10,
	"weight5":          10,
	"vat":              10,
	"vat_bg":           10,
	"vat_lg_purch":     10,
	"vat_lg_sell":      10,
	"excise_rate":      20,
	"excise_tax":       10,
	"saf_stock_qty":    10,
	"min_stock_qty":    10,
	"purch_price1":     10,
	"purch_price2":     10,
	"purch_price3":     10,
	"purch_price4":     10,
	"purch_price5":     10,
	"sell_price1":      10,
	"sell_price2":      10,
	"sell_price3":      10,
	"sell_price4":      10,
	"sell_price5":      10,
	"cogs":             10,
}

func MapHeaderToWeb(header string) string {
	if val, ok := HeaderMapping[header]; ok {
		return val
	}
	return header
}

func MapHeaderFromWeb(header string) string {
	// Hilangkan tambahan seperti (maksimal 100 karakter)
	re := regexp.MustCompile(`\s*\(maksimal\s*\d+\s*karakter\)`)
	header = re.ReplaceAllString(header, "")

	// Hilangkan tanda * di belakang kalau ada
	header = strings.TrimSuffix(header, "*")

	// Hilangkan spasi berlebih di depan/belakang
	header = strings.TrimSpace(header)

	// Lakukan mapping balik ke nama kolom DB
	for k, v := range HeaderMapping {
		if v == header {
			return k
		}
	}

	return header
}

var mandatoryColumns = map[string]bool{
	"pro_name":     true,
	"pro_code":     true,
	"bar_code":     true,
	"pcat_code":    true,
	"pl_code":      true,
	"brand_code":   true,
	"unit_id1":     true,
	"unit_name1":   true,
	"purch_price1": true,
	"sell_price1":  true,
	"is_active":    true,
}

var headerProductStepColor = map[string]struct {
	Step  string
	Color string
}{
	"pro_code":          {"step1_basic", "#C6EFCE"},
	"bar_code":          {"step1_basic", "#C6EFCE"},
	"pro_name":          {"step1_basic", "#C6EFCE"},
	"pcat_code":         {"step1_basic", "#C6EFCE"},
	"pcat_name":         {"step1_basic", "#C6EFCE"},
	"pl_code":           {"step1_basic", "#C6EFCE"},
	"pl_name":           {"step1_basic", "#C6EFCE"},
	"brand_code":        {"step1_basic", "#C6EFCE"},
	"brand_name":        {"step1_basic", "#C6EFCE"},
	"sbrand1_code":      {"step1_basic", "#C6EFCE"},
	"sbrand1_name":      {"step1_basic", "#C6EFCE"},
	"sbrand2_code":      {"step1_basic", "#C6EFCE"},
	"sbrand2_name":      {"step1_basic", "#C6EFCE"},
	"flavor_code":       {"step1_basic", "#C6EFCE"},
	"flavor_name":       {"step1_basic", "#C6EFCE"},
	"ptype_code":        {"step1_basic", "#C6EFCE"},
	"ptype_name":        {"step1_basic", "#C6EFCE"},
	"psize_code":        {"step1_basic", "#C6EFCE"},
	"psize_name":        {"step1_basic", "#C6EFCE"},
	"sup_code":          {"step1_basic", "#C6EFCE"},
	"sup_name":          {"step1_basic", "#C6EFCE"},
	"principal_code":    {"step1_basic", "#C6EFCE"},
	"principal_name":    {"step1_basic", "#C6EFCE"},
	"c_pro_code":        {"step1_basic", "#C6EFCE"},
	"c_pro_name":        {"step1_basic", "#C6EFCE"},
	"unit_id1":          {"step2_unit", "#FFEB9C"},
	"unit_id2":          {"step2_unit", "#FFEB9C"},
	"unit_id3":          {"step2_unit", "#FFEB9C"},
	"unit_name1":        {"step2_unit", "#FFEB9C"},
	"unit_name2":        {"step2_unit", "#FFEB9C"},
	"unit_name3":        {"step2_unit", "#FFEB9C"},
	"conv_unit2":        {"step2_unit", "#FFEB9C"},
	"conv_unit3":        {"step2_unit", "#FFEB9C"},
	"is_batch":          {"step1_basic", "#C6EFCE"},
	"is_exp_date":       {"step1_basic", "#C6EFCE"},
	"parent_pro_id":     {"step1_basic", "#C6EFCE"},
	"pro_code_coretax":  {"step1_basic", "#C6EFCE"},
	"length1":           {"step2_unit", "#FFEB9C"},
	"width1":            {"step2_unit", "#FFEB9C"},
	"height1":           {"step2_unit", "#FFEB9C"},
	"weight1":           {"step2_unit", "#FFEB9C"},
	"length2":           {"step2_unit", "#FFEB9C"},
	"width2":            {"step2_unit", "#FFEB9C"},
	"height2":           {"step2_unit", "#FFEB9C"},
	"weight2":           {"step2_unit", "#FFEB9C"},
	"length3":           {"step2_unit", "#FFEB9C"},
	"width3":            {"step2_unit", "#FFEB9C"},
	"height3":           {"step2_unit", "#FFEB9C"},
	"weight3":           {"step2_unit", "#FFEB9C"},
	"saf_stock_qty":     {"step2_unit", "#FFEB9C"},
	"saf_stock_unit_id": {"step2_unit", "#FFEB9C"},
	"min_stock_qty":     {"step2_unit", "#FFEB9C"},
	"min_stock_unit_id": {"step2_unit", "#FFEB9C"},
	"is_active":         {"step1_basic", "#C6EFCE"},
	"pro_status":        {"step1_basic", "#C6EFCE"},
	"purch_price1":      {"step2_unit", "#FFEB9C"},
	"purch_price2":      {"step2_unit", "#FFEB9C"},
	"purch_price3":      {"step2_unit", "#FFEB9C"},
	"sell_price1":       {"step2_unit", "#FFEB9C"},
	"sell_price2":       {"step2_unit", "#FFEB9C"},
	"sell_price3":       {"step2_unit", "#FFEB9C"},
	"excise_rate":       {"step2_unit", "#FFEB9C"},
	"excise_tax":        {"step2_unit", "#FFEB9C"},
	"vat":               {"step2_unit", "#FFEB9C"},
	"vat_bg":            {"step2_unit", "#FFEB9C"},
	"vat_lg_purch":      {"step2_unit", "#FFEB9C"},
	"vat_lg_sell":       {"step2_unit", "#FFEB9C"},
	"cogs":              {"step2_unit", "#FFEB9C"},
}

var mandatoryProductColumns = map[string]bool{
	// Step 1
	"Product Code":            true,
	"Product Name":            true,
	"Referral Code (Coretax)": true,
	"Product Category Code":   true,
	"Product Category":        true,
	"Product Line Code":       true,
	"Product Line":            true,
	"Brand Code":              true,
	"Brand":                   true,
	"Sub Brand Code":          true,
	"Sub Brand":               true,
	"Sub Brand 2 Code":        true,
	"Sub Brand 2":             true,
	"Flavor Code":             true,
	"Flavor":                  true,
	"Pack Type Code":          true,
	"Pack Type":               true,
	"Pack Size Code":          true,
	"Pack Size":               true,
	"Principal Code":          true,
	"Principal":               true,
	"Supplier Code":           true,
	"Supplier":                true,
	"Consumer Product Code":   true,
	"Consumer Product":        true,
	"Product Status":          true,
	"Expired Date":            true,

	// Step 2
	"Smallest Unit Code":      true,
	"Smallest Unit":           true,
	"Middle Unit Code":        true,
	"Middle Unit":             true,
	"Largest Unit Code":       true,
	"Largest Unit":            true,
	"Conversion Middle Unit":  true,
	"Conversion Largest Unit": true,

	"Purchase Price Smallest Unit": true,
	"Purchase Price Middle Unit":   true,
	"Purchase Price Largest Unit":  true,

	"Sell Price Smallest Unit": true,
	"Sell Price Middle Unit":   true,
	"Sell Price Largest Unit":  true,

	"Dimension Smallest Unit Weight (gram)": true,
	"Dimension Smallest Unit Length (cm)":   true,
	"Dimension Smallest Unit Width (cm)":    true,
	"Dimension Smallest Unit Height (cm)":   true,

	"Dimension Middle Unit Weight (gram)": true,
	"Dimension Middle Unit Length (cm)":   true,
	"Dimension Middle Unit Width (cm)":    true,
	"Dimension Middle Unit Height (cm)":   true,

	"Dimension Largest Unit Weight (gram)": true,
	"Dimension Largest Unit Length (cm)":   true,
	"Dimension Largest Unit Width (cm)":    true,
	"Dimension Largest Unit Height (cm)":   true,

	"Safety Stock":           true,
	"Safety Stock Unit Code": true,
	"Min Stock":              true,
	"Min Stock Unit Code":    true,
}

func (service *productServiceImpl) ExportTemplate(format string) (*bytes.Buffer, string, string, error) {
	var buffer *bytes.Buffer
	var contentType, filename string
	var err error

	instructions, err := service.ProductRepository.GetProductImportInstructions()
	if err != nil {
		return nil, "", "", err
	}

	// Logika untuk memilih format file
	switch format {
	case "csv":
		buffer, err = service.createTemplateCSV(instructions)
		contentType = "application/zip"
		filename = "product_template.zip"
	case "xls":
		buffer, err = service.createTemplateXLS(instructions)
		contentType = "application/vnd.ms-excel"
		filename = "product_template.xls"
	default: // Default ke xlsx jika format tidak diset atau tidak valid
		buffer, err = service.createTemplateXLSX(instructions)
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		filename = "product_template.xlsx"
	}

	if err != nil {
		return nil, "", "", err
	}

	return buffer, contentType, filename, nil
}

func (service *productServiceImpl) createTemplateXLSX(data []model.ImportInstruction) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Product Template"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	// Definisikan semua header kolom sesuai requirement
	headers := []string{
		"pro_name", "pro_code", "bar_code", "pro_code_coretax",
		"pcat_code", "pcat_name",
		"pl_code", "pl_name",
		"brand_code", "brand_name",
		"sbrand1_code", "sbrand1_name",
		"sbrand2_code", "sbrand2_name",
		"flavor_code", "flavor_name",
		"ptype_code", "ptype_name",
		"psize_code", "psize_name",
		"principal_code", "principal_name",
		"sup_code", "sup_name",
		"c_pro_code", "c_pro_name",
		"parent_pro_id",
		"is_active", "is_batch",
		"is_exp_date", "pro_status",
		"unit_id3", "unit_name3", "unit_id2", "unit_name2", "unit_id1", "unit_name1",
		"conv_unit3", "conv_unit2",
		"purch_price3", "purch_price2", "purch_price1",
		"sell_price3", "sell_price2", "sell_price1",
		"weight3", "length3", "width3", "height3",
		"weight2", "length2", "width2", "height2",
		"weight1", "length1", "width1", "height1",
		"saf_stock_qty", "saf_stock_unit_id",
		"min_stock_qty", "min_stock_unit_id",
		"vat", "vat_bg", "vat_lg_purch", "vat_lg_sell", "cogs",
		"excise_rate", "excise_tax",
	}

	for i, header := range headers {
		displayName := MapHeaderToWeb(header)
		if mandatoryProductColumns[displayName] {
			displayName = displayName + "*"
		}
		if max, ok := headerProductMaxLength[header]; ok {
			displayName += fmt.Sprintf("(maksimal %d karakter)", max)
		}
		cell := fmt.Sprintf("%c1", 'A'+i)
		if i >= 26 {
			// Handle columns beyond Z (AA, AB, etc.)
			firstLetter := 'A' + (i / 26) - 1
			secondLetter := 'A' + (i % 26)
			cell = fmt.Sprintf("%c%c1", firstLetter, secondLetter)
		}
		f.SetCellValue(sheetName, cell, displayName)

		color := "#CCCCCC"
		if info, ok := headerProductStepColor[header]; ok {
			color = info.Color
		}
		// Set style untuk header
		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
		})
		f.SetCellStyle(sheetName, cell, cell, style)

		col, _ := excelize.ColumnNumberToName(i + 1)
		width := float64(len(displayName)) * 1.1 // faktor skala biar pas
		_ = f.SetColWidth(sheetName, col, col, width)
	}

	instSheet := "Instructions"
	_, err = f.NewSheet(instSheet)
	if err != nil {
		return nil, err
	}

	// Header
	instHeaders := []string{"Kolom", "Mandatory", "Keterangan", "Step", "Warna"}
	for i, h := range instHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(instSheet, cell, h)

		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
		})
		f.SetCellStyle(instSheet, cell, cell, style)
	}

	// Data baris
	maxLen := make([]int, len(instHeaders))
	for r, d := range data {
		values := []string{
			d.Kolom,
			d.Mandatory,
			safeString(d.Keterangan),
			d.Step,
			d.Color,
		}
		for i, v := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, r+2)
			f.SetCellValue(instSheet, cell, v)

			if l := len(v); l > maxLen[i] {
				maxLen[i] = l
			}
		}
	}

	for i, h := range instHeaders {
		col, _ := excelize.ColumnNumberToName(i + 1)

		// Ambil panjang maksimum antara header & data
		headerLen := len(h)
		maxLen[i] = maxProductLen(headerLen, maxLen[i])

		// Skala sedikit agar tidak terlalu pas (lebih enak dibaca)
		width := float64(maxLen[i]) * 1.1

		_ = f.SetColWidth(instSheet, col, col, width)
	}

	f.SetActiveSheet(index)
	if err := f.DeleteSheet("Sheet1"); err != nil {
		log.Warnf("createXLS: Could not delete default sheet: %v", err)
	}

	return f.WriteToBuffer()
}

// createTemplateXLS adalah fungsi privat untuk membuat template Excel format lama (XLS)
func (service *productServiceImpl) createTemplateXLS(data []model.ImportInstruction) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Product Template"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	// Definisikan semua header kolom sesuai requirement (sama dengan XLSX)
	headers := []string{
		"pro_name", "pro_code", "bar_code", "pro_code_coretax",
		"pcat_code", "pcat_name",
		"pl_code", "pl_name",
		"brand_code", "brand_name",
		"sbrand1_code", "sbrand1_name",
		"sbrand2_code", "sbrand2_name",
		"flavor_code", "flavor_name",
		"ptype_code", "ptype_name",
		"psize_code", "psize_name",
		"principal_code", "principal_name",
		"sup_code", "sup_name",
		"c_pro_code", "c_pro_name",
		"parent_pro_id",
		"is_active", "is_batch",
		"is_exp_date", "pro_status",
		"unit_id3", "unit_name3", "unit_id2", "unit_name2", "unit_id1", "unit_name1",
		"conv_unit3", "conv_unit2",
		"purch_price3", "purch_price2", "purch_price1",
		"sell_price3", "sell_price2", "sell_price1",
		"weight3", "length3", "width3", "height3",
		"weight2", "length2", "width2", "height2",
		"weight1", "length1", "width1", "height1",
		"saf_stock_qty", "saf_stock_unit_id",
		"min_stock_qty", "min_stock_unit_id",
		"vat", "vat_bg", "vat_lg_purch", "vat_lg_sell", "cogs",
		"excise_rate", "excise_tax",
	}

	for i, header := range headers {
		displayName := MapHeaderToWeb(header)
		if mandatoryProductColumns[displayName] {
			displayName = displayName + "*"
		}
		if max, ok := headerProductMaxLength[header]; ok {
			displayName += fmt.Sprintf("(maksimal %d karakter)", max)
		}
		cell := fmt.Sprintf("%c1", 'A'+i)
		if i >= 26 {
			firstLetter := 'A' + (i / 26) - 1
			secondLetter := 'A' + (i % 26)
			cell = fmt.Sprintf("%c%c1", firstLetter, secondLetter)
		}
		f.SetCellValue(sheetName, cell, displayName)

		color := "#CCCCCC"
		if info, ok := headerProductStepColor[header]; ok {
			color = info.Color
		}
		// Set style untuk header
		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
		})
		f.SetCellStyle(sheetName, cell, cell, style)

		col, _ := excelize.ColumnNumberToName(i + 1)
		width := float64(len(displayName)) * 1.1 // faktor skala biar pas
		// if width < 10 {
		// 	width = 10 // minimal lebar
		// }
		_ = f.SetColWidth(sheetName, col, col, width)
	}

	instSheet := "Instructions"
	_, err = f.NewSheet(instSheet)
	if err != nil {
		return nil, err
	}

	// Header
	instHeaders := []string{"Kolom", "Mandatory", "Keterangan", "Step", "Warna"}
	for i, h := range instHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(instSheet, cell, h)

		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
		})
		f.SetCellStyle(instSheet, cell, cell, style)
	}

	// Data baris
	maxLen := make([]int, len(instHeaders))
	for r, d := range data {
		values := []string{
			d.Kolom,
			d.Mandatory,
			safeString(d.Keterangan),
			d.Step,
			d.Color,
		}
		for i, v := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, r+2)
			f.SetCellValue(instSheet, cell, v)

			if l := len(v); l > maxLen[i] {
				maxLen[i] = l
			}
		}
	}

	for i, h := range instHeaders {
		col, _ := excelize.ColumnNumberToName(i + 1)

		// Ambil panjang maksimum antara header & data
		headerLen := len(h)
		maxLen[i] = maxProductLen(headerLen, maxLen[i])

		// Skala sedikit agar tidak terlalu pas (lebih enak dibaca)
		width := float64(maxLen[i]) * 1.1

		_ = f.SetColWidth(instSheet, col, col, width)
	}

	f.SetActiveSheet(index)
	if err := f.DeleteSheet("Sheet1"); err != nil {
		log.Warnf("createXLS: Could not delete default sheet: %v", err)
	}

	return f.WriteToBuffer()
}

// createTemplateCSV adalah fungsi privat untuk membuat template CSV
func (service *productServiceImpl) createTemplateCSV(instructions []model.ImportInstruction) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	fw, err := zw.Create("product_template.csv")
	if err != nil {
		return nil, err
	}
	cw := csv.NewWriter(fw)
	cw.Comma = ';'

	// Tulis header
	headers := []string{
		"pro_name", "pro_code", "bar_code", "pro_code_coretax",
		"pcat_code", "pcat_name",
		"pl_code", "pl_name",
		"brand_code", "brand_name",
		"sbrand1_code", "sbrand1_name",
		"sbrand2_code", "sbrand2_name",
		"flavor_code", "flavor_name",
		"ptype_code", "ptype_name",
		"psize_code", "psize_name",
		"principal_code", "principal_name",
		"sup_code", "sup_name",
		"c_pro_code", "c_pro_name",
		"parent_pro_id",
		"is_active", "is_batch",
		"is_exp_date", "pro_status",
		"unit_id3", "unit_name3", "unit_id2", "unit_name2", "unit_id1", "unit_name1",
		"conv_unit3", "conv_unit2",
		"purch_price3", "purch_price2", "purch_price1",
		"sell_price3", "sell_price2", "sell_price1",
		"weight3", "length3", "width3", "height3",
		"weight2", "length2", "width2", "height2",
		"weight1", "length1", "width1", "height1",
		"saf_stock_qty", "saf_stock_unit_id",
		"min_stock_qty", "min_stock_unit_id",
		"vat", "vat_bg", "vat_lg_purch", "vat_lg_sell", "cogs",
		"excise_rate", "excise_tax",
	}

	webHeaders := make([]string, len(headers))
	for i, h := range headers {
		displayName := MapHeaderToWeb(h)
		if mandatoryProductColumns[displayName] {
			displayName = displayName + "*"
		}
		if max, ok := headerProductMaxLength[h]; ok {
			displayName += fmt.Sprintf("(maksimal %d karakter)", max)
		}
		webHeaders[i] = displayName // ubah nama DB ke nama web
	}

	if err := cw.Write(webHeaders); err != nil {
		return nil, err
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		return nil, err
	}

	// 2) instructions.csv
	iw, err := zw.Create("instructions.csv")
	if err != nil {
		return nil, err
	}
	icw := csv.NewWriter(iw)
	icw.Comma = ';'
	_ = icw.Write([]string{"kolom", "mandatory", "keterangan", "step", "warna"})
	for _, it := range instructions {
		mand := it.Mandatory
		ket := ""
		if it.Keterangan != nil {
			ket = *it.Keterangan
		}
		_ = icw.Write([]string{it.Kolom, mand, ket, it.Step, it.Color})
	}
	icw.Flush()
	if err := icw.Error(); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}

func maxProductLen(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var FieldGroups = map[string]string{
	// Brand
	// "brand_id":   "brand",
	"brand_code": "brand",
	"brand_name": "brand",

	// Sub Brand 1
	"sbrand1_code": "sbrand1",
	"sbrand1_name": "sbrand1",

	// Sub Brand 2
	"sbrand2_code": "sbrand2",
	"sbrand2_name": "sbrand2",

	// Category
	"pcat_code": "category",
	"pcat_name": "category",

	// Product Line
	"pl_code": "productline",
	"pl_name": "productline",

	// Flavor
	"flavor_code": "flavor",
	"flavor_name": "flavor",

	// Pack Type
	"ptype_code": "packtype",
	"ptype_name": "packtype",

	// Pack Size
	"psize_code": "size",
	"psize_name": "size",

	// Supplier
	"sup_code": "supplier",
	"sup_name": "supplier",

	// Principal
	"principal_code": "principal",
	"principal_name": "principal",

	// Cons Product
	"c_pro_code": "conspro",
	"c_pro_name": "conspro",

	// Master Product (semua field product)
	"pro_code":   "product",
	"pro_name":   "product",
	"bar_code":   "product",
	"cogs":       "product",
	"pro_status": "product",
	"is_active":  "product",

	// Tambahan product fields
	// Units
	"unit_id1": "product", "unit_id2": "product", "unit_id3": "product",
	"unit_name1": "product", "unit_name2": "product", "unit_name3": "product",
	"conv_unit2": "product", "conv_unit3": "product",

	// unit
	"unit_id": "unit", "unit_name": "unit",

	// Flags
	"is_batch":    "product",
	"is_exp_date": "product",

	// Parent & New
	"parent_pro_id": "product",
	// "is_new_pro":    "product",

	// Purchase Prices
	"purch_price1": "product", "purch_price2": "product",
	"purch_price3": "product",

	// Sell Prices
	"sell_price1": "product", "sell_price2": "product",
	"sell_price3": "product",

	// Dimension per unit
	"length1": "product", "length2": "product", "length3": "product",
	"width1": "product", "width2": "product", "width3": "product",
	"height1": "product", "height2": "product", "height3": "product",
	"weight1": "product", "weight2": "product", "weight3": "product",

	// Stock
	"saf_stock_qty":     "product",
	"saf_stock_unit_id": "product",
	"min_stock_qty":     "product",
	"min_stock_unit_id": "product",

	// Tax & Status
	"excise_rate": "product",
	"excise_tax":  "product",

	"vat":          "product",
	"vat_bg":       "product",
	"vat_lg_purch": "product",
	"vat_lg_sell":  "product",

	// coretax
	"pro_code_coretax": "product",
	// "pro_name_coretax": "coretax",
}

func (service *productServiceImpl) ExportTemplateUpdate(custId string, format string, fields []string) (*bytes.Buffer, string, string, error) {
	// Tentukan headers dari fields
	headers := []string{}
	added := make(map[string]bool)

	ensure := func(f string) {
		if !added[f] {
			headers = append(headers, f)
			added[f] = true
		}
	}

	// 2. Konversi field dari UI (web name) → internal name
	internalFields := []string{}
	for _, f := range fields {
		internalFields = append(internalFields, MapHeaderFromWeb(f))
	}
	log.Info(internalFields)

	// 3. Terapkan rules FieldGroups
	for _, f := range internalFields {
		group := FieldGroups[f]

		switch group {
		case "brand", "sbrand1", "sbrand2", "category", "principal",
			"productline", "flavor", "packtype", "size", "supplier", "conspro":

			// kalau pilih code/name saja → tambahkan id
			if strings.HasSuffix(f, "_code") || strings.HasSuffix(f, "_name") {
				base := strings.TrimSuffix(f, "_code")
				base = strings.TrimSuffix(base, "_name")
				// ensure(base + "_id")
				ensure(base + "_code")
				ensure(base + "_name")
			}
			// kalau pilih id saja → tambahkan code dan name
			if strings.HasSuffix(f, "_id") {
				base := strings.TrimSuffix(f, "_id")
				ensure(f)
				ensure(base + "_code")
				ensure(base + "_name")
				continue
			}
			ensure("pro_code")
			ensure("pro_name")
			ensure(f)

		case "product":
			// product field → selalu tambahkan pro_id, pro_code, pro_name
			ensure("pro_code")
			ensure("pro_name")
			ensure(f)

			if f == "unit_id1" || f == "unit_id2" || f == "unit_id3" || f == "unit_name1" || f == "unit_name2" || f == "unit_name3" {
				// ambil suffix angka setelah "unit_id"/"unit_name"
				num := strings.TrimPrefix(strings.TrimPrefix(f, "unit_id"), "unit_name")
				if num != "" {
					ensure("unit_id" + num)
					ensure("unit_name" + num)
					ensure(f)
				}
			}

		case "unit":
			if strings.HasSuffix(f, "_name") {
				base := strings.TrimSuffix(f, "_name")
				ensure(f)
				ensure(base + "_id")
			}
			if strings.HasSuffix(f, "_id") {
				base := strings.TrimSuffix(f, "_id")
				ensure(f)
				ensure(base + "_name")
				continue
			}
			ensure(f)

		default:
			// fallback → langsung tambahkan
			// ensure(f)
		}
	}

	hasProduct := false
	hasBrand := false
	hasProductLine := false
	hasUnitX := false
	hasUnit := false
	for _, f := range headers {
		switch FieldGroups[f] {
		case "product":
			hasProduct = true
		case "brand":
			hasBrand = true
		case "productline":
			hasProductLine = true
		case "unit":
			hasUnit = true
		}
	}
	for _, f := range headers {
		if f == "unit_id1" || f == "unit_id2" || f == "unit_id3" {
			hasUnitX = true
			break
		}
	}

	if hasUnitX || (hasProduct && hasUnit) {
		// buang semua entry "unit" dari headers
		filtered := []string{}
		for _, f := range headers {
			if f == "unit_id" || f == "unit_name" {
				continue
			}
			filtered = append(filtered, f)
		}
		headers = filtered
	}

	if hasProduct && hasBrand {
		// ensure("sbrand1_id")
		ensure("sbrand1_code")
		ensure("sbrand1_name")
	}

	// Rule 2: product + productline → tambahkan brand + subbrand1
	if hasProduct && hasProductLine {
		// brand
		// ensure("brand_id")
		ensure("brand_code")
		ensure("brand_name")

		// subbrand1
		// ensure("sbrand1_id")
		ensure("sbrand1_code")
		ensure("sbrand1_name")
	}

	if len(headers) == 0 {
		return nil, "", "", fmt.Errorf("tidak ada header valid dari fields %v", fields)
	}
	log.Info(headers)

	// 4. Ambil data dari repo pakai internal headers
	// data, err := service.ProductRepository.GetDataForTemplateUpdate(custId, headers)
	// if err != nil {
	// 	return nil, "", "", err
	// }
	instructions, err := service.ProductRepository.GetProductImportInstructions()
	if err != nil {
		return nil, "", "", err
	}

	// 6. Generate file sesuai format
	var buffer *bytes.Buffer
	var contentType, filename string

	switch format {
	case "csv":
		buffer, err = service.createTemplateUpdateCSV(headers, instructions)
		contentType = "application/zip"
		filename = "product_template_update.zip"
	case "xlsx":
		buffer, err = service.createTemplateUpdateXLSX(headers, instructions)
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		filename = "product_template_update.xlsx"
	case "xls":
		buffer, err = service.createTemplateUpdateXLSX(headers, instructions)
		contentType = "application/vnd.ms-excel"
		filename = "product_template_update.xls"
	default:
		return nil, "", "", fmt.Errorf("format %s tidak didukung", format)
	}

	if err != nil {
		return nil, "", "", err
	}
	return buffer, contentType, filename, nil
}

func reorderHeaders(headers []string) []string {
	priority := []string{"pro_name", "pro_code", "bar_code", "pro_code_coretax",
		"pcat_code", "pcat_name",
		"pl_code", "pl_name",
		"brand_code", "brand_name",
		"sbrand1_code", "sbrand1_name",
		"sbrand2_code", "sbrand2_name",
		"flavor_code", "flavor_name",
		"ptype_code", "ptype_name",
		"psize_code", "psize_name",
		"principal_code", "principal_name",
		"sup_code", "sup_name",
		"c_pro_code", "c_pro_name",
		"parent_pro_id",
		"is_active", "is_batch",
		"is_exp_date", "pro_status",
		"unit_id3", "unit_name3", "unit_id2", "unit_name2", "unit_id1", "unit_name1",
		"conv_unit3", "conv_unit2",
		"purch_price3", "purch_price2", "purch_price1",
		"sell_price3", "sell_price2", "sell_price1",
		"weight3", "length3", "width3", "height3",
		"weight2", "length2", "width2", "height2",
		"weight1", "length1", "width1", "height1",
		"saf_stock_qty", "saf_stock_unit_id",
		"min_stock_qty", "min_stock_unit_id",
		"vat", "vat_bg", "vat_lg_purch", "vat_lg_sell", "cogs",
		"excise_rate", "excise_tax"} // yang harus muncul duluan
	reordered := []string{}
	seen := make(map[string]bool)

	// tambahkan yang prioritas
	for _, p := range priority {
		for _, h := range headers {
			if h == p {
				reordered = append(reordered, h)
				seen[h] = true
			}
		}
	}

	// tambahkan sisanya sesuai urutan asli
	for _, h := range headers {
		if !seen[h] {
			reordered = append(reordered, h)
		}
	}

	return reordered
}

// func (service *productServiceImpl) createTemplateUpdateXLSX(headers []string, data map[string][][]string, instruction []model.ImportInstruction) (*bytes.Buffer, error) {
func (service *productServiceImpl) createTemplateUpdateXLSX(headers []string, data []model.ImportInstruction) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Product Template Update"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	// 1. Tulis headers di row pertama
	headers = reorderHeaders(headers)
	for i, header := range headers {
		displayName := MapHeaderToWeb(header)

		// ✅ Tambahkan info maksimal karakter (tanpa tanda *)
		if max, ok := headerProductMaxLength[header]; ok {
			displayName += fmt.Sprintf("(maksimal %d karakter)", max)
		}

		// Tulis header ke cell
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, displayName)

		color := "#CCCCCC"
		if info, ok := headerProductStepColor[header]; ok {
			color = info.Color
		}
		// Style untuk header
		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
		})
		f.SetCellStyle(sheetName, cell, cell, style)

		col, _ := excelize.ColumnNumberToName(i + 1)
		width := float64(len(displayName)) * 1.1 // faktor skala biar pas
		_ = f.SetColWidth(sheetName, col, col, width)
	}

	instSheet := "Instructions"
	_, err = f.NewSheet(instSheet)
	if err != nil {
		return nil, err
	}

	// Header
	instHeaders := []string{"Kolom", "Mandatory", "Keterangan", "Step", "Warna"}
	for i, h := range instHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(instSheet, cell, h)

		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
		})
		f.SetCellStyle(instSheet, cell, cell, style)
	}

	// Data baris
	maxLen := make([]int, len(instHeaders))
	for r, d := range data {
		values := []string{
			d.Kolom,
			d.Mandatory,
			safeString(d.Keterangan),
			d.Step,
			d.Color,
		}
		for i, v := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, r+2)
			f.SetCellValue(instSheet, cell, v)

			if l := len(v); l > maxLen[i] {
				maxLen[i] = l
			}
		}
	}

	for i, h := range instHeaders {
		col, _ := excelize.ColumnNumberToName(i + 1)

		// Ambil panjang maksimum antara header & data
		headerLen := len(h)
		maxLen[i] = maxProductLen(headerLen, maxLen[i])

		// Skala sedikit agar tidak terlalu pas (lebih enak dibaca)
		width := float64(maxLen[i]) * 1.1

		_ = f.SetColWidth(instSheet, col, col, width)
	}

	f.SetActiveSheet(index)
	if err := f.DeleteSheet("Sheet1"); err != nil {
		log.Warnf("createXLS: Could not delete default sheet: %v", err)
	}

	return f.WriteToBuffer()
}

func (service *productServiceImpl) createTemplateUpdateCSV(headers []string, data []model.ImportInstruction) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	// === 1) product_template_update.csv ===
	tw, err := zw.Create("product_template_update.csv")
	if err != nil {
		return nil, err
	}
	cw := csv.NewWriter(tw)
	cw.Comma = ';'

	headers = reorderHeaders(headers)
	webHeaders := make([]string, len(headers))
	for i, h := range headers {
		displayName := MapHeaderToWeb(h)

		// ✅ tambahkan maksimal karakter (tanpa tanda *)
		if max, ok := headerProductMaxLength[h]; ok {
			displayName += fmt.Sprintf("(maksimal %d karakter)", max)
		}

		webHeaders[i] = displayName
	}

	if err := cw.Write(webHeaders); err != nil {
		return nil, err
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		return nil, err
	}

	// === 2) instructions.csv ===
	iw, err := zw.Create("instructions.csv")
	if err != nil {
		return nil, err
	}
	icw := csv.NewWriter(iw)
	icw.Comma = ';'

	_ = icw.Write([]string{"Kolom", "Mandatory", "Keterangan", "Step", "Warna"})
	for _, it := range data {
		mand := it.Mandatory
		ket := ""
		if it.Keterangan != nil {
			ket = *it.Keterangan
		}
		_ = icw.Write([]string{it.Kolom, mand, ket, it.Step, it.Color})
	}
	icw.Flush()
	if err := icw.Error(); err != nil {
		return nil, err
	}

	// Tutup ZIP writer
	if err := zw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *productServiceImpl) ImportProductCSV4(req entity.ImportProductRequest) error {
	reader := csv.NewReader(req.File)
	reader.Comma = ';' // Gunakan ini jika pemisah adalah titik koma

	// 1. Baca Header
	rawHeader, err := reader.Read()
	if err != nil {
		return errors.New("Gagal membaca header dari file CSV. Pastikan file tidak kosong dan formatnya benar.")
	}

	// 2. Validasi Header
	if err := s.validateHeader(rawHeader); err != nil {
		return err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return errors.New("Gagal membaca data dari file CSV. Periksa kembali format pemisah dan struktur kolom.")
	}
	headers := make([]string, len(rawHeader))
	for i, h := range rawHeader {
		headers[i] = MapHeaderFromWeb(strings.TrimSpace(h))
	}
	log.Info(headers)

	// --- FR-9: Simpan history awal dengan status "processing" ---
	historyId, err := s.ProductRepository.CreateImportHistory("product", req.FileName, req.CustId, req.CreatedBy, len(records))
	if err != nil {
		return err
	}

	errChan := make(chan error, len(records))
	var wg sync.WaitGroup

	for i, row := range records {
		wg.Add(1)
		go func(i int, r []string) {
			defer wg.Done()
			_, err := s.mapRowToStruct(headers, r)
			if err != nil {
				importData := s.mapRowToStructTemp(headers, r, err) // meski error, tetap ambil data mentah
				_ = s.ProductRepository.InsertProductTemp(historyId, "failed", req.CustId, importData)
				errChan <- fmt.Errorf("Kode Produk %s: %w", importData.ProCode, err) // +2 karena header + index slice
			}
		}(i, row)
	}

	wg.Wait()
	close(errChan)

	var errs []string
	for e := range errChan {
		errs = append(errs, e.Error())
	}

	if len(errs) > 0 {
		failedCount := len(errs)
		successCount := len(records) - failedCount

		// Update history (status masih processing, tapi update progress awal)
		_ = s.ProductRepository.UpdateImportHistory(historyId, successCount, failedCount, false)

		return fmt.Errorf("Validasi data gagal karena: %s", strings.Join(errs, "\n"))
	}

	// --- FR-8: Jalankan async ---
	go func() {
		err := s.processRows(headers, req, records, req.FileName, historyId)
		if err != nil {
			// update status history jadi "failed"
			_ = s.ProductRepository.UpdateImportHistory(historyId, 0, len(records), true)
			log.Error("Import async failed: %v", err)
			return
		}
		log.Info("Import async success, historyId=%d", historyId)
	}()

	return nil
}

// --- FUNGSI UTAMA UNTUK XLSX ---
func (s *productServiceImpl) ImportProductXLSX4(req entity.ImportProductRequest) error {
	f, err := excelize.OpenReader(req.File)
	if err != nil {
		return errors.New("Gagal membaca file. Pastikan file tidak kosong dan formatnya benar.")
	}
	defer f.Close()

	var rows [][]string
	found := false

	for _, sheet := range f.GetSheetList() {
		r, err := f.GetRows(sheet)
		if err != nil || len(r) < 2 {
			continue
		}

		if s.validateHeader(r[0]) == nil {
			rows = r
			found = true
			break
		}
	}

	if !found {
		return errors.New("Tidak ditemukan sheet dengan header yang valid di dalam file Excel.")
	}

	if len(rows) < 2 {
		return errors.New("File Excel harus memiliki header dan minimal satu baris data.")
	}

	// 1. Validasi Header
	if err := s.validateHeader(rows[0]); err != nil {
		return err
	}

	rawHeaders := rows[0]
	headers := make([]string, len(rawHeaders))
	for i, h := range rawHeaders {
		headers[i] = MapHeaderFromWeb(strings.TrimSpace(h))
	}
	log.Info(headers)

	historyId, err := s.ProductRepository.CreateImportHistory("product", req.FileName, req.CustId, req.CreatedBy, len(rows)-1)
	if err != nil {
		return err
	}

	errChan := make(chan error, len(rows)-1)
	var wg sync.WaitGroup

	for i, row := range rows[1:] {
		wg.Add(1)
		go func(i int, r []string) {
			defer wg.Done()
			_, err := s.mapRowToStruct(headers, r)
			if err != nil {
				importData := s.mapRowToStructTemp(headers, r, err) // meski error, tetap ambil data mentah
				importData.DistributorId = req.DistributorId
				if req.DistributorId != 0 {
					importData.Level = 1
				}
				_ = s.ProductRepository.InsertProductTemp(historyId, "failed", req.CustId, importData)
				errChan <- fmt.Errorf("Kode Produk %s: %w", importData.ProCode, err) // +2 karena header + index slice
			}
		}(i, row)
	}

	wg.Wait()
	close(errChan)

	var errs []string
	for e := range errChan {
		errs = append(errs, e.Error())
	}
	if len(errs) > 0 {
		failedCount := len(errs)
		successCount := (len(rows) - 1) - failedCount

		// Update history (status masih processing, tapi update progress awal)
		_ = s.ProductRepository.UpdateImportHistory(historyId, successCount, failedCount, false)
		return fmt.Errorf("Validasi data gagal karena: %s", strings.Join(errs, "\n"))
	}

	// --- FR-8: Jalankan async ---
	go func() {
		err := s.processRows(headers, req, rows[1:], req.FileName, historyId)
		if err != nil {
			// update status history jadi "failed"
			_ = s.ProductRepository.UpdateImportHistory(historyId, 0, len(rows)-1, true)
			log.Error("Import XLSX async failed: %v", err)
			return
		}
		log.Info("Import XLSX async success, historyId=%d", historyId)
	}()

	// --- FR-7: Return cepat ke UI (≤ 5 detik) ---
	return nil
}

func (s *productServiceImpl) validateHeader(header []string) error {
	requiredHeaders := map[string]bool{
		"pro_code":         false,
		"pro_name":         false,
		"pcat_code":        false,
		"pcat_name":        false,
		"brand_code":       false,
		"brand_name":       false,
		"sbrand1_code":     false,
		"sbrand1_name":     false,
		"sbrand2_code":     false,
		"sbrand2_name":     false,
		"flavor_code":      false,
		"flavor_name":      false,
		"ptype_code":       false,
		"ptype_name":       false,
		"psize_code":       false,
		"psize_name":       false,
		"sup_code":         false,
		"sup_name":         false,
		"principal_code":   false,
		"principal_name":   false,
		"c_pro_code":       false,
		"c_pro_name":       false,
		"pro_code_coretax": false,
	}

	for _, h := range header {
		internalName := strings.ToLower(MapHeaderFromWeb(h))
		if _, ok := requiredHeaders[internalName]; ok {
			requiredHeaders[internalName] = true
		}
	}

	missing := []string{}
	for h, found := range requiredHeaders {
		if !found {
			missing = append(missing, MapHeaderToWeb(h))
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("Beberapa kolom wajib tidak ditemukan dalam file: - Kolom yang hilang: %v Pastikan nama kolom sesuai dengan template import yang disediakan.", strings.Join(missing, " - "))
	}
	return nil
}

func (s *productServiceImpl) processRows(headers []string, req entity.ImportProductRequest, rows [][]string, filename string, historyId int64) error {
	total := len(rows)
	success := 0
	// start := time.Now()
	targetCustId := req.CustId
	// if req.CustId != req.ParentCustId && req.ParentCustId != "" {
	// 	targetCustId = req.ParentCustId
	// }

	for i, row := range rows {
		// Konversi baris mentah menjadi struct ImportProductRow
		importData, err := s.mapRowToStruct(headers, row)
		if err != nil {
			log.Info("error in row %d: %w", i+2, err)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed", req.CustId, importData)
			continue
		}

		exist, _ := s.ProductRepository.CheckProductExists(targetCustId, req.ParentCustId, importData.ProCode)
		if exist {
			errMsg := fmt.Sprintf("Kode Produk %s sudah ada. Silakan ganti kode terlebih dahulu.", importData.ProCode)
			importData.ErrorMessage = errMsg

			// Simpan ke temp table agar bisa dilihat user
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-duplicate", req.CustId, importData)

			log.Warnf("Baris %d dilewati: %s", i+2, errMsg)
			continue
		}

		// 1. Dapatkan atau buat Product Line
		pcatId, err := s.ProductRepository.FindByCodeProductCategory(targetCustId, req.ParentCustId, importData.PcatCode, importData.PcatName)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: %v", i+2, err)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-product-category", req.CustId, importData)
			continue
		}

		plId, err := s.ProductRepository.FindByCodeProductLine(targetCustId, req.ParentCustId, importData.PlCode, importData.PlName)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: %v", i+2, err)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-product-line", req.CustId, importData)
			continue
		}

		brandId, err := s.ProductRepository.FindByCodeBrand(targetCustId, req.ParentCustId, plId, importData.BrandCode, importData.BrandName)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: %v", i+2, err)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-brand", req.CustId, importData)
			continue
		}

		sbrand1Id, err := s.ProductRepository.FindByCodeSubBrand1(targetCustId, req.ParentCustId, brandId, importData.Sbrand1Code, importData.Sbrand1Name)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: %v", i+2, err)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-sub-brand", req.CustId, importData)
			continue
		}

		sbrand2Id, err := s.ProductRepository.FindByCodeSubBrand2(targetCustId, req.ParentCustId, importData.Sbrand2Code, importData.Sbrand2Name)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: %v", i+2, err)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-sub-brand2", req.CustId, importData)
			continue
		}

		flavorId, err := s.ProductRepository.FindByCodeFlavor(targetCustId, req.ParentCustId, importData.FlavorCode, importData.FlavorName)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: %v", i+2, err)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-flavor", req.CustId, importData)
			continue
		}

		psizeId, err := s.ProductRepository.FindByCodePackSize(targetCustId, req.ParentCustId, importData.PSizeCode, importData.PSizeName)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: %v", i+2, err)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-pack-size", req.CustId, importData)
			continue
		}

		ptypeId, err := s.ProductRepository.FindByCodePackType(targetCustId, req.ParentCustId, importData.PTypeCode, importData.PTypeName)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: %v", i+2, err)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-pack-type", req.CustId, importData)
			continue
		}

		supId, err := s.ProductRepository.FindByCodeSupplier(targetCustId, req.ParentCustId, importData.SupCode, importData.SupName)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: %v", i+2, err)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-supplier", req.CustId, importData)
			continue
		}

		principalId, err := s.ProductRepository.FindByCodePrincipal(targetCustId, req.ParentCustId, importData.PrincipalCode, importData.PrincipalName)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: %v", i+2, err)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-principal", req.CustId, importData)
			continue
		}

		cproId, err := s.ProductRepository.FindByCodeCPro(targetCustId, req.ParentCustId, importData.CProCode, importData.CProName)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: %v", i+2, err)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-consumer-product", req.CustId, importData)
			continue
		}

		unitId1, err := s.ProductRepository.FindByCodeUnit(targetCustId, req.ParentCustId, importData.UnitId1, importData.UnitName1)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: kode %s dan nama %s untuk Smallest Unit tidak ditemukan, silakan setup parameter terlebih dahulu", i+2, importData.UnitId1, importData.UnitName1)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-smallest-unit", req.CustId, importData)
			continue
		}

		unitId2, err := s.ProductRepository.FindByCodeUnit(targetCustId, req.ParentCustId, importData.UnitId2, importData.UnitName2)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: kode %s dan nama %s untuk Middle Unit tidak ditemukan, silakan setup parameter terlebih dahulu", i+2, importData.UnitId2, importData.UnitName2)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-middle-unit", req.CustId, importData)
			continue
		}

		unitId3, err := s.ProductRepository.FindByCodeUnit(targetCustId, req.ParentCustId, importData.UnitId3, importData.UnitName3)
		if err != nil {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: kode %s dan nama %s untuk Largest Unit tidak ditemukan, silakan setup parameter terlebih dahulu", i+2, importData.UnitId3, importData.UnitName3)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-largest-unit", req.CustId, importData)
			continue
		}

		proCodeCoretax, err := s.ProductRepository.GetUnitProductCoretaxIdByCode(targetCustId, req.ParentCustId, importData.ProCodeCoretax)
		if err != nil || proCodeCoretax == "" {
			importData.ErrorMessage = fmt.Sprintf("Baris %d: Referral Code (Coretax) %s tidak ditemukan", i+2, importData.ProCodeCoretax)
			_ = s.ProductRepository.InsertProductTemp(historyId, "failed-coretax", req.CustId, importData)
			continue
		}

		// ... Lanjutkan proses ini untuk semua entitas lain (sbrand2, flavor, pcat, dll.) ...

		isMainPro := false
		sortNo, _ := strconv.Atoi(importData.SortNo)
		itemNo, _ := strconv.Atoi(importData.ItemNo)
		convUnit2, _ := strconv.ParseFloat(importData.ConvUnit2, 64)
		convUnit3, _ := strconv.ParseFloat(importData.ConvUnit3, 64)
		convUnit4, _ := strconv.ParseFloat(importData.ConvUnit4, 64)
		convUnit5, _ := strconv.ParseFloat(importData.ConvUnit5, 64)
		isBatch, _ := strconv.ParseBool(importData.IsBatch)
		isExpDate, _ := strconv.ParseBool(importData.IsExpDate)
		length, _ := strconv.ParseFloat(importData.Length, 64)
		width, _ := strconv.ParseFloat(importData.Width, 64)
		height, _ := strconv.ParseFloat(importData.Height, 64)
		weight, _ := strconv.ParseFloat(importData.Weight, 64)
		volume := length * width * height
		isNewPro, _ := strconv.ParseBool(importData.IsNewPro)
		purchPrice1, _ := strconv.ParseFloat(importData.PurchPrice1, 64)
		purchPrice2, _ := strconv.ParseFloat(importData.PurchPrice2, 64)
		purchPrice3, _ := strconv.ParseFloat(importData.PurchPrice3, 64)
		purchPrice4, _ := strconv.ParseFloat(importData.PurchPrice4, 64)
		purchPrice5, _ := strconv.ParseFloat(importData.PurchPrice5, 64)
		sellPrice1, _ := strconv.ParseFloat(importData.SellPrice1, 64)
		sellPrice2, _ := strconv.ParseFloat(importData.SellPrice2, 64)
		sellPrice3, _ := strconv.ParseFloat(importData.SellPrice3, 64)
		sellPrice4, _ := strconv.ParseFloat(importData.SellPrice4, 64)
		sellPrice5, _ := strconv.ParseFloat(importData.SellPrice5, 64)
		length1, _ := strconv.ParseFloat(importData.Length1, 64)
		length2, _ := strconv.ParseFloat(importData.Length2, 64)
		length3, _ := strconv.ParseFloat(importData.Length3, 64)
		length4, _ := strconv.ParseFloat(importData.Length4, 64)
		length5, _ := strconv.ParseFloat(importData.Length5, 64)
		width1, _ := strconv.ParseFloat(importData.Width1, 64)
		width2, _ := strconv.ParseFloat(importData.Width2, 64)
		width3, _ := strconv.ParseFloat(importData.Width3, 64)
		width4, _ := strconv.ParseFloat(importData.Width4, 64)
		width5, _ := strconv.ParseFloat(importData.Width5, 64)
		height1, _ := strconv.ParseFloat(importData.Height1, 64)
		height2, _ := strconv.ParseFloat(importData.Height2, 64)
		height3, _ := strconv.ParseFloat(importData.Height3, 64)
		height4, _ := strconv.ParseFloat(importData.Height4, 64)
		height5, _ := strconv.ParseFloat(importData.Height5, 64)
		weight1, _ := strconv.ParseFloat(importData.Weight1, 64)
		weight2, _ := strconv.ParseFloat(importData.Weight2, 64)
		weight3, _ := strconv.ParseFloat(importData.Weight3, 64)
		weight4, _ := strconv.ParseFloat(importData.Weight4, 64)
		weight5, _ := strconv.ParseFloat(importData.Weight5, 64)
		volume1 := length1 * width1 * height1
		volume2 := length2 * width2 * height2
		volume3 := length3 * width3 * height3
		volume4 := length4 * width4 * height4
		volume5 := length5 * width5 * height5
		safStockQty, _ := strconv.ParseFloat(importData.SafStockQty, 64)
		parentProId, _ := toInt64(importData.ParentProId)
		// safStockUnitId, _ := strconv.Atoi(importData.SafStockUnitId)
		minStockQty, _ := strconv.ParseFloat(importData.MinStockQty, 64)
		// minStockUnitId, _ := strconv.Atoi(importData.MinStockUnitId)
		exciseRate, _ := strconv.ParseFloat(importData.ExciseRate, 64)
		exciseTax, _ := strconv.ParseFloat(importData.ExciseTax, 64)
		isActive, _ := strconv.ParseBool(importData.IsActive)
		isDel, _ := strconv.ParseBool(importData.IsDel)
		vat, _ := strconv.ParseFloat(importData.Vat, 64)
		vatBg, _ := strconv.ParseFloat(importData.VatBg, 64)
		vatLgPurch, _ := strconv.ParseFloat(importData.VatLgPurch, 64)
		vatLgSell, _ := strconv.ParseFloat(importData.VatLgSell, 64)
		cogs, _ := strconv.ParseFloat(importData.Cogs, 64)

		importData.DistributorId = req.DistributorId
		if req.DistributorId != 0 {
			importData.Level = 1
		}

		// Setelah semua ID didapatkan, siapkan data untuk tabel utama
		processedData := entity.ProcessedProductRow{
			CustId:         targetCustId,
			ProCode:        importData.ProCode,
			ProName:        importData.ProName,
			BarCode:        importData.BarCode,
			PcatId:         pcatId,
			BrandId:        brandId,
			Sbrand1Id:      sbrand1Id,
			Sbrand2Id:      sbrand2Id,
			FlavorId:       flavorId,
			PSizeId:        psizeId,
			PTypeId:        ptypeId,
			SupId:          supId,
			PrincipalId:    principalId,
			CProId:         cproId,
			IsMainPro:      isMainPro,
			SortNo:         sortNo,
			ItemNo:         itemNo,
			UnitId1:        unitId1,
			UnitId2:        unitId2,
			UnitId3:        unitId3,
			UnitId4:        importData.UnitId4,
			UnitId5:        importData.UnitId5,
			ConvUnit2:      convUnit2,
			ConvUnit3:      convUnit3,
			ConvUnit4:      convUnit4,
			ConvUnit5:      convUnit5,
			IsBatch:        isBatch,
			IsExpDate:      isExpDate,
			Length:         length,
			Width:          width,
			Height:         height,
			Weight:         weight,
			Volume:         volume,
			ParentProId:    parentProId,
			IsNewPro:       isNewPro,
			PurchPrice1:    purchPrice1,
			PurchPrice2:    purchPrice2,
			PurchPrice3:    purchPrice3,
			PurchPrice4:    purchPrice4,
			PurchPrice5:    purchPrice5,
			SellPrice1:     sellPrice1,
			SellPrice2:     sellPrice2,
			SellPrice3:     sellPrice3,
			SellPrice4:     sellPrice4,
			SellPrice5:     sellPrice5,
			Length1:        length1,
			Length2:        length2,
			Length3:        length3,
			Length4:        length4,
			Length5:        length5,
			Width1:         width1,
			Width2:         width2,
			Width3:         width3,
			Width4:         width4,
			Width5:         width5,
			Height1:        height1,
			Height2:        height2,
			Height3:        height3,
			Height4:        height4,
			Height5:        height5,
			Weight1:        weight1,
			Weight2:        weight2,
			Weight3:        weight3,
			Weight4:        weight4,
			Weight5:        weight5,
			Volume1:        volume1,
			Volume2:        volume2,
			Volume3:        volume3,
			Volume4:        volume4,
			Volume5:        volume5,
			SafStockQty:    safStockQty,
			SafStockUnitId: importData.SafStockUnitId,
			MinStockQty:    minStockQty,
			MinStockUnitId: importData.MinStockUnitId,
			ExciseRate:     exciseRate,
			ExciseTax:      exciseTax,
			IsActive:       isActive,
			IsDel:          isDel,
			ImageUrl:       importData.ImageUrl,
			Vat:            vat,
			VatBg:          vatBg,
			VatLgPurch:     vatLgPurch,
			VatLgSell:      vatLgSell,
			Cogs:           cogs,
			ProStatus:      mapProductStatusToInt(importData.ProStatus),
			ProCodeCoretax: proCodeCoretax,
			DistributorId:  nullableProductDistributorID(importData.DistributorId),
			Level:          importData.Level,
			Origin:         "import",
			CreatedBy:      Int64PtrProd2(req.CreatedBy),
			CreatedAt:      TimePtrProd2(time.Now().In(time.UTC)),
			UpdatedBy:      Int64PtrProd2(req.CreatedBy),
			UpdatedAt:      TimePtrProd2(time.Now().In(time.UTC)),
		}
		log.Info("Processed row %d: %+v", i+2, processedData)

		// Terakhir, simpan produk utama
		if err := s.ProductRepository.CreateProduct(processedData); err != nil {
			return fmt.Errorf("failed to save product from row %d: %w", i+2, err)
		}

		success++

	}

	failed := total - success
	if err := s.ProductRepository.UpdateImportHistory(historyId, success, failed, false); err != nil {
		return err
	}
	// Jika semua berhasil, commit transaksi.
	// return tx.Commit()
	return nil
}

func (s *productServiceImpl) mapRowToStruct(headers, row []string) (entity.ImportProductRow, error) {
	// Buat map header -> value
	values := map[string]string{}
	for i, h := range headers {
		if i < len(row) {
			values[h] = strings.TrimSpace(row[i])
		} else {
			values[h] = ""
		}
	}

	data := entity.ImportProductRow{
		ProCode:        values["pro_code"],
		BarCode:        values["bar_code"],
		ProName:        values["pro_name"],
		PcatId:         values["pcat_id"],
		PcatCode:       values["pcat_code"],
		PcatName:       values["pcat_name"],
		PlId:           values["pl_id"],
		PlCode:         values["pl_code"],
		PlName:         values["pl_name"],
		BrandId:        values["brand_id"],
		BrandCode:      values["brand_code"],
		BrandName:      values["brand_name"],
		Sbrand1Id:      values["sbrand1_id"],
		Sbrand1Code:    values["sbrand1_code"],
		Sbrand1Name:    values["sbrand1_name"],
		Sbrand2Id:      values["sbrand2_id"],
		Sbrand2Code:    values["sbrand2_code"],
		Sbrand2Name:    values["sbrand2_name"],
		FlavorId:       values["flavor_id"],
		FlavorCode:     values["flavor_code"],
		FlavorName:     values["flavor_name"],
		PTypeId:        values["ptype_id"],
		PTypeCode:      values["ptype_code"],
		PTypeName:      values["ptype_name"],
		PSizeId:        values["psize_id"],
		PSizeCode:      values["psize_code"],
		PSizeName:      values["psize_name"],
		SupId:          values["sup_id"],
		SupCode:        values["sup_code"],
		SupName:        values["sup_name"],
		PrincipalId:    values["principal_id"],
		PrincipalCode:  values["principal_code"],
		PrincipalName:  values["principal_name"],
		CProId:         values["c_pro_id"],
		CProCode:       values["c_pro_code"],
		CProName:       values["c_pro_name"],
		IsMainPro:      values["is_main_pro"],
		SortNo:         "0",
		ItemNo:         "0",
		UnitId1:        values["unit_id1"],
		UnitId2:        values["unit_id2"],
		UnitId3:        values["unit_id3"],
		UnitId4:        values["unit_id4"],
		UnitId5:        values["unit_id5"],
		UnitName1:      values["unit_name1"],
		UnitName2:      values["unit_name2"],
		UnitName3:      values["unit_name3"],
		ConvUnit2:      values["conv_unit2"],
		ConvUnit3:      values["conv_unit3"],
		ConvUnit4:      values["conv_unit4"],
		ConvUnit5:      values["conv_unit5"],
		IsBatch:        values["is_batch"],
		IsExpDate:      values["is_exp_date"],
		Length:         values["length"],
		Width:          values["width"],
		Height:         values["height"],
		Weight:         values["weight"],
		Volume:         values["volume"],
		ParentProId:    values["parent_pro_id"],
		IsNewPro:       values["is_new_pro"],
		PurchPrice1:    values["purch_price1"],
		PurchPrice2:    values["purch_price2"],
		PurchPrice3:    values["purch_price3"],
		PurchPrice4:    values["purch_price4"],
		PurchPrice5:    values["purch_price5"],
		SellPrice1:     values["sell_price1"],
		SellPrice2:     values["sell_price2"],
		SellPrice3:     values["sell_price3"],
		SellPrice4:     values["sell_price4"],
		SellPrice5:     values["sell_price5"],
		Length1:        values["length1"],
		Length2:        values["length2"],
		Length3:        values["length3"],
		Length4:        values["length4"],
		Length5:        values["length5"],
		Width1:         values["width1"],
		Width2:         values["width2"],
		Width3:         values["width3"],
		Width4:         values["width4"],
		Width5:         values["width5"],
		Height1:        values["height1"],
		Height2:        values["height2"],
		Height3:        values["height3"],
		Height4:        values["height4"],
		Height5:        values["height5"],
		Weight1:        values["weight1"],
		Weight2:        values["weight2"],
		Weight3:        values["weight3"],
		Weight4:        values["weight4"],
		Weight5:        values["weight5"],
		Volume1:        values["volume1"],
		Volume2:        values["volume2"],
		Volume3:        values["volume3"],
		Volume4:        values["volume4"],
		Volume5:        values["volume5"],
		SafStockQty:    values["saf_stock_qty"],
		SafStockUnitId: values["saf_stock_unit_id"],
		MinStockQty:    values["min_stock_qty"],
		MinStockUnitId: values["min_stock_unit_id"],
		ExciseRate:     values["excise_rate"],
		ExciseTax:      values["excise_tax"],
		IsActive:       values["is_active"],
		IsDel:          values["is_del"],
		ImageUrl:       values["image_url"],
		Vat:            values["vat"],
		VatBg:          values["vat_bg"],
		VatLgPurch:     values["vat_lg_purch"],
		VatLgSell:      values["vat_lg_sell"],
		Cogs:           values["cogs"],
		ProStatus:      values["pro_status"],
		ProCodeCoretax: values["pro_code_coretax"],
		ProNameCoretax: values["pro_name_coretax"],
	}

	// Validasi
	// if values["pro_code"] == "" || values["pro_name"] == "" {
	// 	return data, errors.New("pro_code and pro_name cannot be empty")
	// }
	// step 1
	for _, f := range []string{
		"pro_code", "pro_name", "pro_code_coretax",
		"pcat_code", "pcat_name",
		"pl_code", "pl_name",
		"brand_code", "brand_name",
		"sbrand1_code", "sbrand1_name",
		"sbrand2_code", "sbrand2_name",
		"flavor_code", "flavor_name",
		"ptype_code", "ptype_name",
		"psize_code", "psize_name",
		"principal_code", "principal_name",
		"sup_code", "sup_name",
		"c_pro_code", "c_pro_name",
		"pro_status", "is_exp_date",
	} {
		if values[f] == "" {
			return data, fmt.Errorf("%s tidak boleh kosong", MapHeaderToWeb(f))
		}
	}
	// step 2
	for _, f := range []string{
		"unit_id1", "unit_name1",
		"unit_id2", "unit_name2",
		"unit_id3", "unit_name3",
		"conv_unit2", "conv_unit3",
		"purch_price1", "purch_price2", "purch_price3",
		"sell_price1", "sell_price2", "sell_price3",
		"weight1", "length1", "width1", "height1",
		"weight2", "length2", "width2", "height2",
		"weight3", "length3", "width3", "height3",
		"saf_stock_qty", "saf_stock_unit_id",
		"min_stock_qty", "min_stock_unit_id",
	} {
		if values[f] == "" {
			return data, fmt.Errorf("%s tidak boleh kosong", MapHeaderToWeb(f))
		}
	}

	for field, maxLen := range lengthProductRules {
		if len(values[field]) > maxLen {
			return data, fmt.Errorf("%s melebihi panjang maksimal %d karakter (panjang saat ini %d)", MapHeaderToWeb(field), maxLen, len(values[field]))
		}
	}

	for field, rule := range numericProductRules {
		if err := s.validateNumeric(field, values[field], rule[0], rule[1]); err != nil {
			return data, err
		}
	}

	if err := s.validateConvUnitInteger("conv_unit2", values["conv_unit2"]); err != nil {
		return data, err
	}
	if err := s.validateConvUnitInteger("conv_unit3", values["conv_unit3"]); err != nil {
		return data, err
	}

	for _, f := range numericDimensionProductRules {
		if values[f] != "" {
			if _, err := strconv.ParseFloat(values[f], 64); err != nil {
				return data, fmt.Errorf("%s harus berupa angka", MapHeaderToWeb(f))
			}
		}
	}

	return data, nil
}

func (s *productServiceImpl) mapRowToStructTemp(headers, row []string, errMsg error) entity.ImportProductRow {
	// Buat map header -> value
	values := map[string]string{}
	for i, h := range headers {
		if i < len(row) {
			values[h] = strings.TrimSpace(row[i])
		} else {
			values[h] = ""
		}
	}

	data := entity.ImportProductRow{
		ProCode:        values["pro_code"],
		BarCode:        values["bar_code"],
		ProName:        values["pro_name"],
		PcatId:         values["pcat_id"],
		PcatCode:       values["pcat_code"],
		PcatName:       values["pcat_name"],
		PlId:           values["pl_id"],
		PlCode:         values["pl_code"],
		PlName:         values["pl_name"],
		BrandId:        values["brand_id"],
		BrandCode:      values["brand_code"],
		BrandName:      values["brand_name"],
		Sbrand1Id:      values["sbrand1_id"],
		Sbrand1Code:    values["sbrand1_code"],
		Sbrand1Name:    values["sbrand1_name"],
		Sbrand2Id:      values["sbrand2_id"],
		Sbrand2Code:    values["sbrand2_code"],
		Sbrand2Name:    values["sbrand2_name"],
		FlavorId:       values["flavor_id"],
		FlavorCode:     values["flavor_code"],
		FlavorName:     values["flavor_name"],
		PTypeId:        values["ptype_id"],
		PTypeCode:      values["ptype_code"],
		PTypeName:      values["ptype_name"],
		PSizeId:        values["psize_id"],
		PSizeCode:      values["psize_code"],
		PSizeName:      values["psize_name"],
		SupId:          values["sup_id"],
		SupCode:        values["sup_code"],
		SupName:        values["sup_name"],
		PrincipalId:    values["principal_id"],
		PrincipalCode:  values["principal_code"],
		PrincipalName:  values["principal_name"],
		CProId:         values["c_pro_id"],
		CProCode:       values["c_pro_code"],
		CProName:       values["c_pro_name"],
		IsMainPro:      values["is_main_pro"],
		SortNo:         "0",
		ItemNo:         "0",
		UnitId1:        values["unit_id1"],
		UnitId2:        values["unit_id2"],
		UnitId3:        values["unit_id3"],
		UnitId4:        values["unit_id4"],
		UnitId5:        values["unit_id5"],
		UnitName1:      values["unit_name1"],
		UnitName2:      values["unit_name2"],
		UnitName3:      values["unit_name3"],
		ConvUnit2:      values["conv_unit2"],
		ConvUnit3:      values["conv_unit3"],
		ConvUnit4:      values["conv_unit4"],
		ConvUnit5:      values["conv_unit5"],
		IsBatch:        values["is_batch"],
		IsExpDate:      values["is_exp_date"],
		Length:         values["length"],
		Width:          values["width"],
		Height:         values["height"],
		Weight:         values["weight"],
		Volume:         values["volume"],
		ParentProId:    values["parent_pro_id"],
		IsNewPro:       values["is_new_pro"],
		PurchPrice1:    values["purch_price1"],
		PurchPrice2:    values["purch_price2"],
		PurchPrice3:    values["purch_price3"],
		PurchPrice4:    values["purch_price4"],
		PurchPrice5:    values["purch_price5"],
		SellPrice1:     values["sell_price1"],
		SellPrice2:     values["sell_price2"],
		SellPrice3:     values["sell_price3"],
		SellPrice4:     values["sell_price4"],
		SellPrice5:     values["sell_price5"],
		Length1:        values["length1"],
		Length2:        values["length2"],
		Length3:        values["length3"],
		Length4:        values["length4"],
		Length5:        values["length5"],
		Width1:         values["width1"],
		Width2:         values["width2"],
		Width3:         values["width3"],
		Width4:         values["width4"],
		Width5:         values["width5"],
		Height1:        values["height1"],
		Height2:        values["height2"],
		Height3:        values["height3"],
		Height4:        values["height4"],
		Height5:        values["height5"],
		Weight1:        values["weight1"],
		Weight2:        values["weight2"],
		Weight3:        values["weight3"],
		Weight4:        values["weight4"],
		Weight5:        values["weight5"],
		Volume1:        values["volume1"],
		Volume2:        values["volume2"],
		Volume3:        values["volume3"],
		Volume4:        values["volume4"],
		Volume5:        values["volume5"],
		SafStockQty:    values["saf_stock_qty"],
		SafStockUnitId: values["saf_stock_unit_id"],
		MinStockQty:    values["min_stock_qty"],
		MinStockUnitId: values["min_stock_unit_id"],
		ExciseRate:     values["excise_rate"],
		ExciseTax:      values["excise_tax"],
		IsActive:       values["is_active"],
		IsDel:          values["is_del"],
		ImageUrl:       values["image_url"],
		Vat:            values["vat"],
		VatBg:          values["vat_bg"],
		VatLgPurch:     values["vat_lg_purch"],
		VatLgSell:      values["vat_lg_sell"],
		Cogs:           values["cogs"],
		ProStatus:      values["pro_status"],
		ProCodeCoretax: values["pro_code_coretax"],
		ProNameCoretax: values["pro_name_coretax"],
		ErrorMessage:   errMsg.Error(),
	}

	return data
}

func (s *productServiceImpl) validateNumeric(fieldName, value string, precision, scale int) error {
	// value kosong
	if value == "" {
		return nil
	}

	// cek numeric valid
	if _, err := strconv.ParseFloat(value, 64); err != nil {
		return fmt.Errorf("%s harus berupa angka", MapHeaderToWeb(fieldName))
	}

	// cek digit & desimal
	parts := strings.Split(value, ".")
	totalDigits := len(strings.ReplaceAll(value, ".", ""))
	fracDigits := 0
	if len(parts) == 2 {
		fracDigits = len(parts[1])
	}

	if totalDigits > precision {
		return fmt.Errorf("%s melebihi jumlah digit maksimum (%d digit)", MapHeaderToWeb(fieldName), precision)
	}
	if fracDigits > scale {
		return fmt.Errorf("%s memiliki angka desimal terlalu banyak (maksimal %d angka di belakang koma)", MapHeaderToWeb(fieldName), scale)
	}
	return nil
}

var lengthProductRules = map[string]int{
	"pro_code": 30, "pro_name": 150,
	"bar_code":   50,
	"brand_code": 5, "brand_name": 100,
	"flavor_code": 5, "flavor_name": 100,
	"sbrand1_code": 5, "sbrand1_name": 100,
	"sbrand2_code": 5, "sbrand2_name": 100,
	"pl_code": 5, "pl_name": 100,
	"pcat_code": 5, "pcat_name": 100,
	"ptype_code": 5, "ptype_name": 50,
	"psize_code": 5, "psize_name": 100,
	"principal_code": 8, "principal_name": 150,
	"sup_code": 20, "sup_name": 150,
	"c_pro_code": 5, "c_pro_name": 150,
	"unit_id1": 5, "unit_id2": 5, "unit_id3": 5,
	// "unit_id4": 5, "unit_id5": 5,
	"saf_stock_unit_id": 10, "min_stock_unit_id": 10,
	"pro_code_coretax": 10,
}

var numericProductRules = map[string][2]int{
	"weight1":       {10, 2},
	"weight2":       {10, 2},
	"weight3":       {10, 2},
	"weight4":       {10, 2},
	"weight5":       {10, 2},
	"vat":           {10, 2},
	"vat_bg":        {10, 2},
	"vat_lg_purch":  {10, 2},
	"vat_lg_sell":   {10, 2},
	"excise_rate":   {20, 4},
	"excise_tax":    {10, 2},
	"saf_stock_qty": {10, 2},
	"min_stock_qty": {10, 2},
	"purch_price1":  {10, 0},
	"purch_price2":  {10, 0},
	"purch_price3":  {10, 0},
	"purch_price4":  {10, 0},
	"purch_price5":  {10, 0},
	"sell_price1":   {10, 0},
	"sell_price2":   {10, 0},
	"sell_price3":   {10, 0},
	"sell_price4":   {10, 0},
	"sell_price5":   {10, 0},
	"cogs":          {10, 2},
}

var numericDimensionProductRules = []string{
	"length1", "length2", "length3",
	"width1", "width2", "width3",
	"height1", "height2", "height3",
}

func (s *productServiceImpl) validateConvUnitInteger(fieldName, value string) error {
	const (
		defaultPrecision = 18
		noDecimalScale   = 0
	)

	err := s.validateNumeric(fieldName, value, defaultPrecision, noDecimalScale)

	if err != nil {
		if strings.Contains(err.Error(), "angka desimal terlalu banyak") {
			return fmt.Errorf("%s tidak boleh desimal", MapHeaderToWeb(fieldName))
		}

		return err
	}

	return nil
}

// var numericNonZeroProductRules = []string{
// 	"length1", "width1", "height1",
// 	"length2", "width2", "height2",
// 	"length3", "width3", "height3",
// }

func (s *productServiceImpl) validateLength(fieldName, value string, maxLen int) error {
	if value == "" {
		return nil
	}
	if len(value) > maxLen {
		return fmt.Errorf("%s melebihi jumlah digit maksimum (%d digit)", MapHeaderToWeb(fieldName), maxLen)
	}
	return nil
}

func (s *productServiceImpl) ImportUpdateXLSX(req entity.ImportProductRequest) error {
	f, err := excelize.OpenReader(req.File)
	if err != nil {
		return errors.New("Gagal membuka file Excel. Pastikan format file .xlsx valid dan tidak rusak.")
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return errors.New("Tidak ada sheet yang ditemukan di dalam file Excel.")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return errors.New("Gagal membaca data dari sheet Excel.")
	}

	if len(rows) < 2 {
		return errors.New("File harus memiliki header dan minimal satu baris data.")
	}

	rawHeaders := rows[0]
	headers := make([]string, len(rawHeaders))
	for i, h := range rawHeaders {
		headers[i] = MapHeaderFromWeb(strings.TrimSpace(h))
	}

	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.ToLower(h)] = i
	}

	historyId, err := s.ProductRepository.CreateImportHistory("product-update", req.FileName, req.CustId, req.CreatedBy, len(rows)-1)
	if err != nil {
		return fmt.Errorf("Gagal membuat riwayat impor data: %v", err)
	}

	var (
		mu   sync.Mutex
		errs []string
	)
	maxErr := 50 // tampilkan maksimal 50 error agar FE tidak overload
	var wg sync.WaitGroup

	for i, row := range rows[1:] {
		wg.Add(1)
		go func(i int, r []string) {
			defer wg.Done()
			localErrs := []string{}

			var proCode string
			if idx, ok := headerMap["pro_code"]; ok && idx < len(r) {
				proCode = strings.TrimSpace(r[idx])
			} else {
				proCode = fmt.Sprintf("Baris %d", i+2) // fallback kalau kolom tidak ada
			}

			for j, val := range r {
				field := headers[j]

				// validasi panjang string
				if maxLen, ok := lengthProductRules[field]; ok {
					if err := s.validateLength(field, val, maxLen); err != nil {
						localErrs = append(localErrs, fmt.Sprintf("Kode Produk %s: %v", proCode, err))
					}
				}

				if field == "conv_unit2" || field == "conv_unit3" {
					if err := s.validateConvUnitInteger(field, val); err != nil {
						localErrs = append(localErrs, fmt.Sprintf("Kode Produk %s: %v", proCode, err))
					}
				}

				// validasi numeric
				if rule, ok := numericProductRules[field]; ok {
					if err := s.validateNumeric(field, val, rule[0], rule[1]); err != nil {
						localErrs = append(localErrs, fmt.Sprintf("Kode Produk %s: %v", proCode, err))
					}
				}

				// validasi dimension numeric
				// for _, f := range numericDimensionProductRules {
				// 	if field == f && val != "" {
				// 		if _, err := strconv.ParseFloat(val, 64); err != nil {
				// 			localErrs = append(localErrs, fmt.Sprintf("row %d: %s must be numeric", i+2, MapHeaderToWeb(field)))
				// 		}
				// 	}
				// }
			}

			if len(localErrs) > 0 {
				errMsg := strings.Join(localErrs, "; ")
				temp := s.mapRowToUpdateTemp(historyId, req.CustId, r, headerMap, errMsg)
				temp.StatusInsert = "failed"
				if e := s.ProductRepository.InsertProductUpdateTemp(temp); e != nil {
					log.Errorf("Gagal menyimpan data sementara (baris %d): %v", i+2, e)
				}

				mu.Lock()
				if len(errs) < maxErr {
					limit := maxErr - len(errs)
					if len(localErrs) > limit {
						localErrs = localErrs[:limit]
					}
					errs = append(errs, localErrs...)
				}
				mu.Unlock()
			}
		}(i, row)
	}

	wg.Wait()

	if len(errs) > 0 {
		failedCount := len(errs)
		successCount := (len(rows) - 1) - failedCount

		_ = s.ProductRepository.UpdateImportHistory(historyId, successCount, failedCount, false)
		return fmt.Errorf("Ditemukan kesalahan pada data impor (%d produk gagal):\n%s", failedCount, strings.Join(errs, "\n"))
	}

	go func() {
		err := s.processUpdateRows(req, rows, historyId)
		if err != nil {
			// update status history ke failed kalau error
			_ = s.ProductRepository.UpdateImportHistory(historyId, 0, len(rows)-1, true)
			log.Errorf("Import XLSX async failed (historyId=%d): %v", historyId, err)
			return
		}
		log.Infof("Import XLSX async success, historyId=%d", historyId)
	}()

	// return cepat ke UI
	return nil
}

func (s *productServiceImpl) ImportUpdateCSV(req entity.ImportProductRequest) error {
	reader := csv.NewReader(req.File)
	reader.TrimLeadingSpace = true
	reader.Comma = ';'

	// Baca semua baris
	rows, err := reader.ReadAll()
	if err != nil {
		return errors.New("gagal membaca file CSV. Pastikan file tidak rusak dan formatnya benar.")
	}

	if len(rows) < 2 {
		return errors.New("file CSV harus memiliki header dan minimal satu baris data.")
	}

	rawHeaders := rows[0]
	headers := make([]string, len(rawHeaders))
	for i, h := range rawHeaders {
		headers[i] = MapHeaderFromWeb(strings.TrimSpace(h))
	}

	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.ToLower(h)] = i
	}

	historyId, err := s.ProductRepository.CreateImportHistory("product-update", req.FileName, req.CustId, req.CreatedBy, len(rows)-1)
	if err != nil {
		return fmt.Errorf("gagal membuat riwayat impor data: %v", err)
	}

	var (
		mu   sync.Mutex
		errs []string
	)
	maxErr := 50 // tampilkan maksimal 50 error agar FE tidak overload
	var wg sync.WaitGroup

	for i, row := range rows[1:] {
		wg.Add(1)
		go func(i int, r []string) {
			defer wg.Done()
			localErrs := []string{}

			var proCode string
			if idx, ok := headerMap["pro_code"]; ok && idx < len(r) {
				proCode = strings.TrimSpace(r[idx])
			} else {
				proCode = fmt.Sprintf("Baris %d", i+2) // fallback kalau kolom tidak ada
			}

			for j, val := range r {
				field := headers[j]

				// validasi panjang string
				if maxLen, ok := lengthProductRules[field]; ok {
					if err := s.validateLength(field, val, maxLen); err != nil {
						localErrs = append(localErrs, fmt.Sprintf("Kode Produk %s: %v", proCode, err))
					}
				}

				if field == "conv_unit2" || field == "conv_unit3" {
					if err := s.validateConvUnitInteger(field, val); err != nil {
						localErrs = append(localErrs, fmt.Sprintf("Kode Produk %s: %v", proCode, err))
					}
				}

				// validasi numeric
				if rule, ok := numericProductRules[field]; ok {
					if err := s.validateNumeric(field, val, rule[0], rule[1]); err != nil {
						localErrs = append(localErrs, fmt.Sprintf("Kode Produk %s: %v", proCode, err))
					}
				}

				// validasi dimension numeric
				// for _, f := range numericDimensionProductRules {
				// 	if field == f && val != "" {
				// 		if _, err := strconv.ParseFloat(val, 64); err != nil {
				// 			localErrs = append(localErrs, fmt.Sprintf("row %d: %s must be numeric", i+2, MapHeaderToWeb(field)))
				// 		}
				// 	}
				// }
			}

			if len(localErrs) > 0 {
				errMsg := strings.Join(localErrs, "; ")
				temp := s.mapRowToUpdateTemp(historyId, req.CustId, r, headerMap, errMsg)
				temp.StatusInsert = "failed"
				if e := s.ProductRepository.InsertProductUpdateTemp(temp); e != nil {
					log.Errorf("Gagal menyimpan data sementara (baris %d): %v", i+2, e)
				}

				mu.Lock()
				if len(errs) < maxErr {
					limit := maxErr - len(errs)
					if len(localErrs) > limit {
						localErrs = localErrs[:limit]
					}
					errs = append(errs, localErrs...)
				}
				mu.Unlock()
			}
		}(i, row)
	}

	wg.Wait()

	if len(errs) > 0 {
		failedCount := len(errs)
		successCount := (len(rows) - 1) - failedCount

		_ = s.ProductRepository.UpdateImportHistory(historyId, successCount, failedCount, false)
		return fmt.Errorf("Ditemukan kesalahan pada data impor (%d produk gagal):\n%s", failedCount, strings.Join(errs, "\n"))
	}

	go func() {
		err := s.processUpdateRows(req, rows, historyId)
		if err != nil {
			// update status history ke failed kalau error
			_ = s.ProductRepository.UpdateImportHistory(historyId, 0, len(rows)-1, true)
			log.Errorf("Import CSV async failed (historyId=%d): %v", historyId, err)
			return
		}
		log.Infof("Import CSV async success, historyId=%d", historyId)
	}()

	// return cepat ke UI
	return nil
}

func toInt64(s string) (int64, error) {
	if s == "" {
		return 0, nil // atau tergantung kebutuhan: mau dianggap kosong / error
	}
	return strconv.ParseInt(s, 10, 64)
}

func toInt(s string) (int, error) {
	i, err := strconv.Atoi(strings.TrimSpace(s))
	return i, err
}

func toFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

func toBool(s string) (bool, error) {
	return strconv.ParseBool(strings.TrimSpace(s))
}

// mapping grup -> fungsi update (untuk foreign tables)
// di dalam service
func (s *productServiceImpl) getUpdateFuncs() map[string]func(custId string, row map[string]string) error {
	return map[string]func(custId string, row map[string]string) error{

		"brand": func(custId string, row map[string]string) error {
			brandID, err := s.ProductRepository.GetBrandIdByCode(custId, row["brand_code"])
			if err != nil {
				return fmt.Errorf("Brand Code '%s' tidak ditemukan", row["brand_code"])
			}
			return s.ProductRepository.UpdateImportBrand(
				custId,
				brandID,
				row["brand_code"],
				row["brand_name"],
			)
		},

		"sub-brand1": func(custId string, row map[string]string) error {
			sbrand1ID, err := s.ProductRepository.GetSubBrand1IdByCode(custId, row["sbrand1_code"])
			if err != nil {
				return fmt.Errorf("Sub Brand Code '%s' tidak ditemukan", row["sbrand1_code"])
			}
			return s.ProductRepository.UpdateImportSubBrand1(
				custId,
				sbrand1ID,
				row["sbrand1_code"],
				row["sbrand1_name"],
			)
		},

		"sub-brand2": func(custId string, row map[string]string) error {
			sbrand2ID, err := s.ProductRepository.GetSubBrand2IdByCode(custId, row["sbrand2_code"])
			if err != nil {
				return fmt.Errorf("Sub Brand 2 Code '%s' tidak ditemukan", row["sbrand2_code"])
			}
			return s.ProductRepository.UpdateImportSubBrand2(
				custId,
				sbrand2ID,
				row["sbrand2_code"],
				row["sbrand2_name"],
			)
		},

		"category": func(custId string, row map[string]string) error {
			pcatID, err := s.ProductRepository.GetProductCategoryIdByCode(custId, row["pcat_code"])
			if err != nil {
				return fmt.Errorf("Product Category Code '%s' tidak ditemukan", row["pcat_code"])
			}
			return s.ProductRepository.UpdateImportCategory(
				custId,
				pcatID,
				row["pcat_code"],
				row["pcat_name"],
			)
		},

		"product-line": func(custId string, row map[string]string) error {
			plID, err := s.ProductRepository.GetProductLineIdByCode(custId, row["pl_code"])
			if err != nil {
				return fmt.Errorf("Product Line Code '%s' tidak ditemukan", row["pl_code"])
			}
			return s.ProductRepository.UpdateImportProductLine(
				custId,
				plID,
				row["pl_code"],
				row["pl_name"],
			)
		},

		"flavor": func(custId string, row map[string]string) error {
			flavorID, err := s.ProductRepository.GetFlavorIdByCode(custId, row["flavor_code"])
			if err != nil {
				return fmt.Errorf("Flavor Code '%s' tidak ditemukan", row["flavor_code"])
			}
			return s.ProductRepository.UpdateImportFlavor(
				custId,
				flavorID,
				row["flavor_code"],
				row["flavor_name"],
			)
		},

		"pack-type": func(custId string, row map[string]string) error {
			ptypeID, err := s.ProductRepository.GetPackTypeIdByCode(custId, row["ptype_code"])
			if err != nil {
				return fmt.Errorf("Pack Type Code '%s' tidak ditemukan", row["ptype_code"])
			}
			return s.ProductRepository.UpdateImportPackType(
				custId,
				ptypeID,
				row["ptype_code"],
				row["ptype_name"],
			)
		},

		"pack-size": func(custId string, row map[string]string) error {
			psizeID, err := s.ProductRepository.GetPackSizeIdByCode(custId, row["psize_code"])
			if err != nil {
				return fmt.Errorf("Pack Size Code '%s' tidak ditemukan", row["psize_code"])
			}
			return s.ProductRepository.UpdateImportPackSize(
				custId,
				psizeID,
				row["psize_code"],
				row["psize_name"],
			)
		},

		"supplier": func(custId string, row map[string]string) error {
			supID, err := s.ProductRepository.GetSupplierIdByCode(custId, row["sup_code"])
			if err != nil {
				return fmt.Errorf("Supplier Code '%s' tidak ditemukan", row["sup_code"])
			}
			return s.ProductRepository.UpdateImportSupplier(
				custId,
				supID,
				row["sup_code"],
				row["sup_name"],
			)
		},

		"principal": func(custId string, row map[string]string) error {
			principalID, err := s.ProductRepository.GetPrincipalIdByCode(custId, row["principal_code"])
			if err != nil {
				return fmt.Errorf("Principal Code '%s' tidak ditemukan", row["principal_code"])
			}
			return s.ProductRepository.UpdateImportPrincipal(
				custId,
				principalID,
				row["principal_code"],
				row["principal_name"],
			)
		},

		"consumer-product": func(custId string, row map[string]string) error {
			cproID, err := s.ProductRepository.GetCProIdByCode(custId, row["c_pro_code"])
			if err != nil {
				return fmt.Errorf("Consumer Product Code '%s' tidak ditemukan", row["c_pro_code"])
			}
			return s.ProductRepository.UpdateImportCPro(
				custId,
				cproID,
				row["c_pro_code"],
				row["c_pro_name"],
			)
		},

		"unit": func(custId string, row map[string]string) error {
			unitID, err := s.ProductRepository.GetUnitIdByCode(custId, row["unit_id"])
			if err != nil {
				return fmt.Errorf("Unit Code '%s' tidak ditemukan", row["unit_id"])
			}
			return s.ProductRepository.UpdateImportUnit(
				custId,
				unitID,
				row["unit_name"],
			)
		},
	}
}

func mapProductStatusToInt(code string) int {
	switch code {
	case "Active":
		return 1
	case "In Active":
		return 2
	case "Dead Stock":
		return 3
	case "Dry Up":
		return 4
	default:
		return 0
	}
}

func mapProductStatusToString(code int) string {
	switch code {
	case 1:
		return "Active"
	case 2:
		return "In Active"
	case 3:
		return "Dead Stock"
	case 4:
		return "Dry Up"
	default:
		return "Undefined"
	}
}

func (s *productServiceImpl) mapToProcessedProductRow(custId string, row map[string]string) (entity.ProcessedProductRow, error) {
	var p entity.ProcessedProductRow
	var err error

	p.CustId = custId
	// p.ProId, _ = toInt(row["pro_id"])
	p.ProCode = row["pro_code"]
	p.BarCode = row["bar_code"]
	p.ProName = row["pro_name"]

	if row["pcat_code"] != "" {
		p.PcatId, err = s.ProductRepository.GetProductCategoryIdByCode(p.CustId, row["pcat_code"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Kategori Produk dengan kode %s", row["pcat_code"])
		}
	}

	if row["brand_code"] != "" {
		p.BrandId, err = s.ProductRepository.GetBrandIdByCode(p.CustId, row["brand_code"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Brand dengan kode %s", row["brand_code"])
		}
	}

	if row["sbrand1_code"] != "" {
		p.Sbrand1Id, err = s.ProductRepository.GetSubBrand1IdByCode(p.CustId, row["sbrand1_code"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Sub Brand 1 dengan kode %s", row["sbrand1_code"])
		}
	}

	if row["sbrand2_code"] != "" {
		p.Sbrand2Id, err = s.ProductRepository.GetSubBrand2IdByCode(p.CustId, row["sbrand2_code"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Sub Brand 2 dengan kode %s", row["sbrand2_code"])
		}
	}

	if row["flavor_code"] != "" {
		p.FlavorId, err = s.ProductRepository.GetFlavorIdByCode(p.CustId, row["flavor_code"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Flavor dengan kode %s", row["flavor_code"])
		}
	}

	if row["ptype_code"] != "" {
		p.PTypeId, err = s.ProductRepository.GetPackTypeIdByCode(p.CustId, row["ptype_code"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Jenis Kemasan (Pack Type) dengan kode %s", row["ptype_code"])
		}
	}

	if row["psize_code"] != "" {
		p.PSizeId, err = s.ProductRepository.GetPackSizeIdByCode(p.CustId, row["psize_code"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Ukuran Kemasan (Pack Size) dengan kode %s", row["psize_code"])
		}
	}

	if row["sup_code"] != "" {
		p.SupId, err = s.ProductRepository.GetSupplierIdByCode(p.CustId, row["sup_code"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Supplier dengan kode %s", row["sup_code"])
		}
	}

	if row["principal_code"] != "" {
		p.PrincipalId, err = s.ProductRepository.GetPrincipalIdByCode(p.CustId, row["principal_code"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Principal dengan kode %s", row["principal_code"])
		}
	}

	if row["c_pro_code"] != "" {
		p.CProId, err = s.ProductRepository.GetCProIdByCode(p.CustId, row["c_pro_code"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Consumer Product dengan kode %s", row["c_pro_code"])
		}
	}

	// --- Unit: berupa string ---
	if row["unit_id1"] != "" {
		p.UnitId1, err = s.ProductRepository.GetUnitIdByCode(p.CustId, row["unit_id1"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Unit Terkecil dengan kode %s", row["unit_id1"])
		}
	}
	if row["unit_id2"] != "" {
		p.UnitId2, err = s.ProductRepository.GetUnitIdByCode(p.CustId, row["unit_id2"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Unit Menengah dengan kode %s", row["unit_id2"])
		}
	}
	if row["unit_id3"] != "" {
		p.UnitId3, err = s.ProductRepository.GetUnitIdByCode(p.CustId, row["unit_id3"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Unit Terbesar dengan kode %s", row["unit_id3"])
		}
	}

	if row["pro_code_coretax"] != "" {
		p.ProCodeCoretax, err = s.ProductRepository.GetUnitProductCoretaxIdByCode(p.CustId, p.CustId, row["pro_code_coretax"])
		if err != nil {
			return p, fmt.Errorf("Gagal mendapatkan Referral Code (Coretax) dengan kode %s", row["pro_code_coretax"])
		}
	}

	p.IsMainPro = false
	p.SortNo = 0
	p.ItemNo = 0

	p.UnitId4 = row["unit_id4"]
	p.UnitId5 = row["unit_id5"]

	p.ConvUnit2, _ = toFloat(row["conv_unit2"])
	p.ConvUnit3, _ = toFloat(row["conv_unit3"])
	p.ConvUnit4, _ = toFloat(row["conv_unit4"])
	p.ConvUnit5, _ = toFloat(row["conv_unit5"])

	// Physical
	p.Weight, _ = toFloat(row["weight"])
	p.IsBatch, _ = toBool(row["is_batch"])
	p.IsExpDate, _ = toBool(row["is_exp_date"])
	p.Length, _ = toFloat(row["length"])
	p.Width, _ = toFloat(row["width"])
	p.Height, _ = toFloat(row["height"])
	p.Volume, _ = toFloat(row["volume"])

	// Parent / new
	p.ParentProId, _ = toInt64(row["parent_pro_id"])
	p.IsNewPro, _ = toBool(row["is_new_pro"])

	// Prices
	p.PurchPrice1, _ = toFloat(row["purch_price1"])
	p.PurchPrice2, _ = toFloat(row["purch_price2"])
	p.PurchPrice3, _ = toFloat(row["purch_price3"])
	p.PurchPrice4, _ = toFloat(row["purch_price4"])
	p.PurchPrice5, _ = toFloat(row["purch_price5"])

	p.SellPrice1, _ = toFloat(row["sell_price1"])
	p.SellPrice2, _ = toFloat(row["sell_price2"])
	p.SellPrice3, _ = toFloat(row["sell_price3"])
	p.SellPrice4, _ = toFloat(row["sell_price4"])
	p.SellPrice5, _ = toFloat(row["sell_price5"])

	// Per-unit dimensions
	p.Length1, _ = toFloat(row["length1"])
	p.Length2, _ = toFloat(row["length2"])
	p.Length3, _ = toFloat(row["length3"])
	p.Length4, _ = toFloat(row["length4"])
	p.Length5, _ = toFloat(row["length5"])

	p.Width1, _ = toFloat(row["width1"])
	p.Width2, _ = toFloat(row["width2"])
	p.Width3, _ = toFloat(row["width3"])
	p.Width4, _ = toFloat(row["width4"])
	p.Width5, _ = toFloat(row["width5"])

	p.Height1, _ = toFloat(row["height1"])
	p.Height2, _ = toFloat(row["height2"])
	p.Height3, _ = toFloat(row["height3"])
	p.Height4, _ = toFloat(row["height4"])
	p.Height5, _ = toFloat(row["height5"])

	p.Weight1, _ = toFloat(row["weight1"])
	p.Weight2, _ = toFloat(row["weight2"])
	p.Weight3, _ = toFloat(row["weight3"])
	p.Weight4, _ = toFloat(row["weight4"])
	p.Weight5, _ = toFloat(row["weight5"])

	p.Volume1, _ = toFloat(row["volume1"])
	p.Volume2, _ = toFloat(row["volume2"])
	p.Volume3, _ = toFloat(row["volume3"])
	p.Volume4, _ = toFloat(row["volume4"])
	p.Volume5, _ = toFloat(row["volume5"])

	// Stock
	p.SafStockQty, _ = toFloat(row["saf_stock_qty"])
	p.SafStockUnitId, _ = row["saf_stock_unit_id"]
	p.MinStockQty, _ = toFloat(row["min_stock_qty"])
	p.MinStockUnitId, _ = row["min_stock_unit_id"]

	// Tax
	p.Vat, _ = toFloat(row["vat"])
	p.VatBg, _ = toFloat(row["vat_bg"])
	p.VatLgPurch, _ = toFloat(row["vat_lg_purch"])
	p.VatLgSell, _ = toFloat(row["vat_lg_sell"])
	p.ExciseRate, _ = toFloat(row["excise_rate"])
	p.ExciseTax, _ = toFloat(row["excise_tax"])

	// Flags
	p.IsActive, _ = toBool(row["is_active"])
	p.IsDel, _ = toBool(row["is_del"])

	// Other
	p.ImageUrl = row["image_url"]
	p.Cogs, _ = toFloat(row["cogs"])
	// p.ProCodeCoretax = row["pro_code_coretax"]

	// Map status
	p.ProStatus = mapProductStatusToInt(row["pro_status"])

	return p, err
}

func toBoolWithDefaultProduct(val string, existing bool) bool {
	val = strings.TrimSpace(strings.ToLower(val))
	if val == "" {
		// Kolom kosong → kembalikan nilai lama dari DB
		return existing
	}
	return val == "true" || val == "1"
}

// mapping product khusus
func (s *productServiceImpl) updateProduct(custId string, row map[string]string) error {
	p, err := s.mapToProcessedProductRow(custId, row)
	if err != nil {
		return fmt.Errorf("Gagal memproses data produk: %v", err)
	}

	productId, err := s.ProductRepository.FindProductIdByCode(custId, p.ProCode)
	if err != nil {
		return fmt.Errorf("Gagal mencari produk dengan kode %s: %v", p.ProCode, err)
	}
	if productId == 0 {
		return fmt.Errorf("Produk dengan kode %s tidak ditemukan", p.ProCode)
	}

	// set product_id ke struct
	p.ProId = productId

	existing, err := s.ProductRepository.FindProductById(custId, productId)
	if err != nil {
		return fmt.Errorf("Gagal mengambil data produk yang sudah ada: %v", err)
	}

	p.IsActive = toBoolWithDefaultProduct(row["is_active"], existing.IsActive)
	p.IsDel = toBoolWithDefaultProduct(row["is_del"], existing.IsDel)

	p.Length = s.coalesceFloat(row["length"], existing.Length)
	p.Width = s.coalesceFloat(row["width"], existing.Width)
	p.Height = s.coalesceFloat(row["height"], existing.Height)

	p.Volume = p.Length * p.Width * p.Height

	p.Length1 = s.coalesceFloat(row["length1"], existing.Length1)
	p.Width1 = s.coalesceFloat(row["width1"], existing.Width1)
	p.Height1 = s.coalesceFloat(row["height1"], existing.Height1)

	p.Volume1 = p.Length1 * p.Width1 * p.Height1

	p.Length2 = s.coalesceFloat(row["length2"], existing.Length2)
	p.Width2 = s.coalesceFloat(row["width2"], existing.Width2)
	p.Height2 = s.coalesceFloat(row["height2"], existing.Height2)

	p.Volume2 = p.Length2 * p.Width2 * p.Height2

	p.Length3 = s.coalesceFloat(row["length3"], existing.Length3)
	p.Width3 = s.coalesceFloat(row["width3"], existing.Width3)
	p.Height3 = s.coalesceFloat(row["height3"], existing.Height3)

	p.Volume3 = p.Length3 * p.Width3 * p.Height3

	// pastikan custId diisi
	p.CustId = custId
	return s.ProductRepository.UpdateImportProduct(p)
}

func (s *productServiceImpl) coalesceFloat(val string, fallback float64) float64 {
	if val == "" {
		return fallback
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return fallback
	}
	return f
}

func (s *productServiceImpl) processUpdateRows(req entity.ImportProductRequest, rows [][]string, historyId int64) error {
	total := len(rows) - 1
	success := 0
	// start := time.Now()

	targetCustId := req.CustId
	// if req.ParentCustId != "" && req.ParentCustId != req.CustId {
	// 	targetCustId = req.ParentCustId
	// }

	// baca header
	headers := rows[0]
	headerMap := map[string]int{}
	for i, h := range headers {
		internal := MapHeaderFromWeb(strings.TrimSpace(h))
		if internal != "" {
			headerMap[internal] = i
		}
	}
	log.Info("Header Map:", headerMap)
	dataRows := rows[1:]

	// loop baris data
	for i, row := range dataRows {
		get := func(key string) string {
			if idx, ok := headerMap[key]; ok && idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		rowMap := map[string]string{}
		for k := range headerMap {
			rowMap[k] = get(k)
		}

		proCode := get("pro_code")
		if proCode == "" {
			proCode = fmt.Sprintf("Baris %d", i+2) // fallback kalau kolom kosong
		}

		var rowErrors []string
		var failedType string

		// --- update foreign tables ---
		updatedGroups := map[string]bool{}
		updateFuncs := s.getUpdateFuncs()
		for field, group := range FieldGroups {
			if group == "product" {
				continue
			}
			if rowMap[field] != "" && !updatedGroups[group] {
				if fn, ok := updateFuncs[group]; ok {
					if err := fn(targetCustId, rowMap); err != nil {
						rowErrors = append(rowErrors,
							fmt.Sprintf("Kode Produk %s: Gagal memperbarui data %s — %v",
								proCode, strings.Title(group), err))
						if failedType == "" {
							failedType = group // catat jenis error pertama
						}
					}
				}
				updatedGroups[group] = true
			}
		}

		// --- update product jika ada ---
		if len(rowErrors) == 0 && rowMap["pro_code"] != "" {
			if err := s.updateProduct(targetCustId, rowMap); err != nil {
				rowErrors = append(rowErrors,
					fmt.Sprintf("Kode Produk %s: Gagal memperbarui data produk — %v",
						proCode, err))
				if failedType == "" {
					failedType = "product"
				}
			}
		}

		if len(rowErrors) > 0 {
			// kalau error, simpan ke product_update_temp
			errMsg := strings.Join(rowErrors, "; ")
			status := "failed"
			if failedType != "" {
				status = fmt.Sprintf("failed-%s", failedType)
			}
			temp := s.mapRowToUpdateTemp(historyId, req.CustId, row, headerMap, errMsg)
			temp.StatusInsert = status
			if err := s.ProductRepository.InsertProductUpdateTemp(temp); err != nil {
				log.Errorf("Gagal menyimpan ke table temporary (Kode Produk %s): %v", proCode, err)
			}
			continue // skip row ini
		}
		success++
	}

	failed := total - success
	if err := s.ProductRepository.UpdateImportHistory(historyId, success, failed, false); err != nil {
		return err
	}

	return nil
}

func (s *productServiceImpl) mapRowToUpdateTemp(historyId int64, custId string, row []string, headerMap map[string]int, errMsg string) entity.ImportProductUpdateTemp {
	get := func(key string) string {
		if idx, ok := headerMap[strings.ToLower(strings.TrimSpace(key))]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	return entity.ImportProductUpdateTemp{
		HistoryId: historyId,
		CustId:    custId,

		// Brand & Category
		BrandId:   get("brand_id"),
		BrandCode: get("brand_code"),
		BrandName: get("brand_name"),
		PcatId:    get("pcat_id"),
		PcatCode:  get("pcat_code"),
		PcatName:  get("pcat_name"),
		PlId:      get("pl_id"),
		PlCode:    get("pl_code"),
		PlName:    get("pl_name"),

		// SubBrand
		Sbrand1Id:   get("sbrand1_id"),
		Sbrand1Code: get("sbrand1_code"),
		Sbrand1Name: get("sbrand1_name"),
		Sbrand2Id:   get("sbrand2_id"),
		Sbrand2Code: get("sbrand2_code"),
		Sbrand2Name: get("sbrand2_name"),

		// Flavor
		FlavorId:   get("flavor_id"),
		FlavorCode: get("flavor_code"),
		FlavorName: get("flavor_name"),

		// Unit
		UnitId:    get("unit_id"),
		UnitName:  get("unit_name"),
		UnitId1:   get("unit_id1"),
		UnitName1: get("unit_name1"),
		UnitId2:   get("unit_id2"),
		UnitName2: get("unit_name2"),
		UnitId3:   get("unit_id3"),
		UnitName3: get("unit_name3"),
		UnitId4:   get("unit_id4"),
		UnitId5:   get("unit_id5"),
		ConvUnit2: get("conv_unit2"),
		ConvUnit3: get("conv_unit3"),
		ConvUnit4: get("conv_unit4"),
		ConvUnit5: get("conv_unit5"),

		// Product Type, Size, Supplier, Principal
		PTypeId:       get("ptype_id"),
		PTypeCode:     get("ptype_code"),
		PTypeName:     get("ptype_name"),
		PSizeId:       get("psize_id"),
		PSizeCode:     get("psize_code"),
		PSizeName:     get("psize_name"),
		SupId:         get("sup_id"),
		SupCode:       get("sup_code"),
		SupName:       get("sup_name"),
		PrincipalId:   get("principal_id"),
		PrincipalCode: get("principal_code"),
		PrincipalName: get("principal_name"),

		// CoreTax & Child Product
		CProId:         get("c_pro_id"),
		CProCode:       get("c_pro_code"),
		CProName:       get("c_pro_name"),
		ProCodeCoretax: get("pro_code_coretax"),
		ProNameCoretax: get("pro_name_coretax"),

		// Product Info
		ProId:     get("pro_id"),
		ProCode:   get("pro_code"),
		ProName:   get("pro_name"),
		Barcode:   get("bar_code"),
		Cogs:      get("cogs"),
		ProStatus: get("pro_status"),
		IsActive:  get("is_active"),
		IsMainPro: "false",
		SortNo:    "0",
		ItemNo:    "0",

		// Dimension & Weight
		Length: get("length"),
		Width:  get("width"),
		Height: get("height"),
		Weight: get("weight"),
		Volume: get("volume"),

		Length1: get("length1"),
		Length2: get("length2"),
		Length3: get("length3"),
		Length4: get("length4"),
		Length5: get("length5"),
		Width1:  get("width1"),
		Width2:  get("width2"),
		Width3:  get("width3"),
		Width4:  get("width4"),
		Width5:  get("width5"),
		Height1: get("height1"),
		Height2: get("height2"),
		Height3: get("height3"),
		Height4: get("height4"),
		Height5: get("height5"),
		Weight1: get("weight1"),
		Weight2: get("weight2"),
		Weight3: get("weight3"),
		Weight4: get("weight4"),
		Weight5: get("weight5"),
		Volume1: get("volume1"),
		Volume2: get("volume2"),
		Volume3: get("volume3"),
		Volume4: get("volume4"),
		Volume5: get("volume5"),

		// Stock
		SafStockQty:      get("saf_stock_qty"),
		SafStockUnitId:   get("saf_stock_unit_id"),
		SafStockUnitName: get("saf_stock_unit_name"),
		MinStockQty:      get("min_stock_qty"),
		MinStockUnitId:   get("min_stock_unit_id"),
		MinStockUnitName: get("min_stock_unit_name"),

		// Price
		PurchPrice1: get("purch_price1"),
		PurchPrice2: get("purch_price2"),
		PurchPrice3: get("purch_price3"),
		PurchPrice4: get("purch_price4"),
		PurchPrice5: get("purch_price5"),
		SellPrice1:  get("sell_price1"),
		SellPrice2:  get("sell_price2"),
		SellPrice3:  get("sell_price3"),
		SellPrice4:  get("sell_price4"),
		SellPrice5:  get("sell_price5"),

		// Flags & Tax
		IsBatch:     get("is_batch"),
		IsExpDate:   get("is_exp_date"),
		ParentProId: get("parent_pro_id"),
		IsNewPro:    get("is_new_pro"),
		ExciseRate:  get("excise_rate"),
		ExciseTax:   get("excise_tax"),
		IsDel:       get("is_del"),
		Vat:         get("vat"),
		VatBg:       get("vat_bg"),
		VatLgPurch:  get("vat_lg_purch"),
		VatLgSell:   get("vat_lg_sell"),

		// Metadata
		ImageURL:     get("image_url"),
		StatusInsert: "failed",
		ErrorMessage: errMsg,
	}
}

func (s *productServiceImpl) ReuploadImportUpdateFile(custId string, historyId int64, req entity.ImportRequest) error {
	var rows [][]string
	var err error

	// --- Baca file ---
	switch strings.ToLower(req.Format) {
	case "csv":
		reader := csv.NewReader(req.File)
		reader.TrimLeadingSpace = true
		reader.Comma = ';'
		rows, err = reader.ReadAll()
	default:
		f, e := excelize.OpenReader(req.File)
		if e != nil {
			return errors.New("Gagal membuka file Excel. Pastikan format file .xlsx valid dan tidak rusak.")
		}
		defer f.Close()
		rows, err = f.GetRows(f.GetSheetName(0))
	}

	convertedReq := entity.ImportProductRequest{
		CustId:       req.CustId,
		ParentCustId: req.ParentCustId,
		File:         req.File,
		FileName:     req.Filename,
		CreatedBy:    req.UserId,
	}

	if err != nil {
		return errors.New("Gagal membaca data dari sheet Excel.")
	}
	if len(rows) < 2 {
		return errors.New("File harus memiliki header dan minimal satu baris data.")
	}

	// mapping header
	rawHeaders := rows[0]
	headers := make([]string, len(rawHeaders))
	for i, h := range rawHeaders {
		headers[i] = MapHeaderFromWeb(strings.TrimSpace(h))
	}

	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.ToLower(h)] = i
	}

	// --- Hitung total dari history lama ---
	total, err := s.ProductRepository.GetImportProductTotalData(historyId)
	if err != nil {
		return fmt.Errorf("Gagal membuat riwayat impor data: %v", err)
	}

	// --- Buat history baru ---
	newHistoryId, err := s.ProductRepository.CreateImportHistory(
		"product-update",
		req.Filename,
		req.CustId,
		req.UserId,
		total,
	)
	if err != nil {
		return err
	}

	// --- Validasi tiap baris ---
	var (
		mu   sync.Mutex
		errs []string
	)
	maxErr := 50
	var wg sync.WaitGroup

	for i, row := range rows[1:] {
		wg.Add(1)
		go func(i int, r []string) {
			defer wg.Done()
			localErrs := []string{}

			var proCode string
			if idx, ok := headerMap["pro_code"]; ok && idx < len(r) {
				proCode = strings.TrimSpace(r[idx])
				if proCode == "" {
					proCode = fmt.Sprintf("(baris %d)", i+2)
				}
			} else {
				proCode = fmt.Sprintf("(baris %d)", i+2)
			}

			for j, val := range r {
				field := headers[j]

				// Validasi panjang string
				if maxLen, ok := lengthProductRules[field]; ok {
					if err := s.validateLength(field, val, maxLen); err != nil {
						localErrs = append(localErrs, fmt.Sprintf("Kode Produk %s: %v", proCode, err))
					}
				}

				if field == "conv_unit2" || field == "conv_unit3" {
					if err := s.validateConvUnitInteger(field, val); err != nil {
						localErrs = append(localErrs, fmt.Sprintf("Kode Produk %s: %v", proCode, err))
					}
				}

				// Validasi numeric
				if rule, ok := numericProductRules[field]; ok {
					if err := s.validateNumeric(field, val, rule[0], rule[1]); err != nil {
						localErrs = append(localErrs, fmt.Sprintf("Kode Produk %s: %v", proCode, err))
					}
				}
			}

			if len(localErrs) > 0 {
				errMsg := strings.Join(localErrs, "; ")
				temp := s.mapRowToUpdateTemp(newHistoryId, req.CustId, r, headerMap, errMsg)
				temp.StatusInsert = "failed"
				if e := s.ProductRepository.InsertProductUpdateTemp(temp); e != nil {
					log.Errorf("Gagal menyimpan data sementara (Kode Produk %s): %v", proCode, e)
				}

				mu.Lock()
				if len(errs) < maxErr {
					limit := maxErr - len(errs)
					if len(localErrs) > limit {
						localErrs = localErrs[:limit]
					}
					errs = append(errs, localErrs...)
				}
				mu.Unlock()
			}
		}(i, row)
	}

	wg.Wait()

	// --- Jika ada error validasi ---
	if len(errs) > 0 {
		failedCount := len(errs)
		successCount := (len(rows) - 1) - failedCount

		_ = s.ProductRepository.UpdateImportHistory(newHistoryId, successCount, failedCount, false)
		return fmt.Errorf("Ditemukan kesalahan pada data impor (%d produk gagal):\n%s", failedCount, strings.Join(errs, "\n"))
	}

	// --- Jalankan async proses row ---
	go func() {
		err := s.processUpdateRows(convertedReq, rows, newHistoryId)
		if err != nil {
			// kalau error, set semua gagal
			_ = s.ProductRepository.UpdateImportHistory(newHistoryId, 0, len(rows)-1, true)
			log.Errorf("Reupload async failed (historyId=%d): %v", newHistoryId, err)
			return
		}
		log.Infof("Reupload async success, historyId=%d", newHistoryId)
	}()

	// return cepat ke UI
	return nil
}

func (s *productServiceImpl) ReuploadImportInsertFile(custId string, historyId int64, req entity.ImportRequest) error {
	var rows [][]string
	var err error

	// --- Baca file sesuai format ---
	switch strings.ToLower(req.Format) {
	case "csv":
		reader := csv.NewReader(req.File)
		reader.TrimLeadingSpace = true
		reader.Comma = ';'
		rows, err = reader.ReadAll()
	default:
		f, e := excelize.OpenReader(req.File)
		if e != nil {
			return errors.New("Gagal membaca file. Pastikan file tidak kosong dan formatnya benar.")
		}
		defer f.Close()
		rows, err = f.GetRows(f.GetSheetName(0))
	}
	if err != nil {
		return errors.New("Gagal membaca file. Pastikan file tidak kosong dan formatnya benar.")
	}
	if len(rows) < 2 {
		return errors.New("File Excel harus memiliki header dan minimal satu baris data.")
	}

	// --- Buat header mapping (sama seperti ImportProductXLSX4) ---
	rawHeaders := rows[0]
	headers := make([]string, len(rawHeaders))
	for i, h := range rawHeaders {
		headers[i] = MapHeaderFromWeb(strings.TrimSpace(h))
	}
	log.Infof("Reupload headers: %+v", headers)

	// Hitung total dari history lama
	total, err := s.ProductRepository.GetImportProductTotalData(historyId)
	if err != nil {
		return err
	}

	// Buat history baru untuk reupload
	newHistoryId, err := s.ProductRepository.CreateImportHistory("product", req.Filename, req.CustId, req.UserId, total)
	if err != nil {
		return err
	}

	// validasi data
	errChan := make(chan error, len(rows)-1)
	var wg sync.WaitGroup

	for i, row := range rows[1:] {
		wg.Add(1)
		go func(i int, r []string) {
			defer wg.Done()
			_, err := s.mapRowToStruct(headers, r)
			if err != nil {
				importData := s.mapRowToStructTemp(headers, r, err) // meski error, tetap ambil data mentah
				_ = s.ProductRepository.InsertProductTemp(historyId, "failed", req.CustId, importData)
				errChan <- fmt.Errorf("Kode Produk %s: %w", importData.ProCode, err)
			}
		}(i, row)
	}

	wg.Wait()
	close(errChan)

	// --- Kumpulkan error validasi ---
	var errs []string
	for e := range errChan {
		errs = append(errs, e.Error())
	}
	if len(errs) > 0 {
		failedCount := len(errs)
		successCount := (len(rows) - 1) - failedCount

		// Update history (status masih processing, tapi update progress awal)
		_ = s.ProductRepository.UpdateImportHistory(newHistoryId, successCount, failedCount, false)
		return fmt.Errorf("Validasi data gagal karena: %s", strings.Join(errs, "\n"))
	}

	// Jalankan ulang proses pakai pipeline processRows
	go func() {
		err := s.processRows(headers, entity.ImportProductRequest{
			CustId:    req.CustId,
			FileName:  req.Filename,
			CreatedBy: req.UserId,
		}, rows[1:], req.Filename, newHistoryId)

		if err != nil {
			// kalau gagal, update history sebagai failed
			_ = s.ProductRepository.UpdateImportHistory(newHistoryId, 0, total, true)
			log.Errorf("Reupload insert failed: %v", err)
			return
		}

		log.Infof("Reupload insert success, newHistoryId=%d", newHistoryId)
	}()

	return nil
}

func (service *productServiceImpl) LookupList(dataFilter entity.ProductQueryFilter, custId string) (data []entity.ProductLookupResponse, total int, lastPage int, err error) {
	var products []model.Product

	if dataFilter.DistributorID > 0 {
		products, total, lastPage, err = service.ProductRepository.FindAllByDistributorLookup(dataFilter, custId)
	} else {
		products, total, lastPage, err = service.ProductRepository.FindAllByCustIdLookup(dataFilter, custId)
	}

	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range products {
		var vResp entity.ProductLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *productServiceImpl) LookupDistPrice(dataFilter entity.ProductQueryFilter) (data []entity.ProductLookupDistPrice, total int, lastPage int, err error) {
	log.Debug(structs.StructToJson(dataFilter))
	var products []model.ProductDistPrice

	if dataFilter.DistributorID > 0 {
		products, total, lastPage, err = service.ProductRepository.FindAllByDistributorLookupDistPrice(dataFilter)
	} else {
		products, total, lastPage, err = service.ProductRepository.FindAllByCustIdLookupDistPrice(dataFilter)
	}

	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range products {
		var vResp entity.ProductLookupDistPrice
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *productServiceImpl) SearchList(dataFilter entity.ProductQueryFilter, custId string) (data []entity.ProductSearchResponse, total int, lastPage int, err error) {

	products, total, lastPage, err := service.ProductRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range products {
		var vResp entity.ProductSearchResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *productServiceImpl) Store(request entity.CreateProductBody) (response entity.ProductResponse, err error) {

	product, err := service.ProductRepository.FindOneByProductCodeAndCustId(request.ProductCode, request.CustId)
	if err == nil {
		return response, errors.New("product_code: " + product.ProductCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	productData := model.Product{}
	err = structs.Automapper(request, &productData)
	if err != nil {
		return response, err
	}

	productData.Level = 1
	productData.Origin = "create"
	productData.DistributorID = nullableProductDistributorID(request.DistributorId)
	if productData.DistributorID == nil {
		productData.Level = 0

	}
	productData.CustId = request.CustId
	productData.ProductCode = request.ProductCode
	productData.ProductName = request.ProductName
	productData.IsActive = request.IsActive
	productData.CreatedAt = &timeNow
	productData.CreatedBy = &request.CreatedBy
	productData.UpdatedAt = &timeNow
	productData.UpdatedBy = &request.CreatedBy

	productId, err := service.ProductRepository.Store(context.Background(), productData)
	if err != nil {
		return response, err
	}

	go service.StoreDistProduct(request.CustId, request.CreatedBy, productId)

	response.ProductId = productId

	return response, err
}

func nullableProductDistributorID(distributorID int64) *int64 {
	if distributorID <= 0 {
		return nil
	}
	return &distributorID
}

func (service *productServiceImpl) StoreDistProduct(DistID string, userID, productID int64) (err error) {
	distributors, err := service.ProductRepository.FindDistributor(DistID)
	if err != nil {
		return err
	}
	var distProducts []model.ProductDistCreate
	for _, distributor := range distributors {
		distProduct := model.ProductDistCreate{
			CustId:         distributor.CustId,
			ProId:          &productID,
			IsActive:       false,
			IsAlloc:        false,
			MinStock:       0,
			MinStockStr:    "00000.000.000",
			SafetyStock:    0,
			SafetyStockStr: "00000.000.000",
			PoFormula:      1,
			IsNewPro:       true,
			Vat:            11.0,
			VatBg:          0.00,
			VatLgPurch:     0.00,
			VatLgSell:      0.00,
			Cogs:           0.0000,
			SMweek1:        1,
			SMweek2:        1,
			UpdatedAt:      time.Now(),
			UpdatedBy:      userID,
			ParentProId:    &productID,
		}
		distProducts = append(distProducts, distProduct)
	}
	if len(distProducts) > 0 {
		err = service.ProductRepository.StoreDist(distProducts)
		if err != nil {
			log.Info("StoreDistProduct, StoreDist, err:", err.Error())
			return err
		}
	}

	return nil
}
func (service *productServiceImpl) BulkStore(request entity.BulkProductBody) (response entity.BulkProductResponse, err error) {
	var (
		products []model.Product
		custId   string
		userId   int64
	)

	timeNow := time.Now().In(time.UTC)
	for index, request := range request.Products {
		if product, err := service.ProductRepository.FindOneByProductCodeAndCustId(request.ProductCode, request.CustId); err == nil {
			return response, errors.New("Product in line " + fmt.Sprint(index+2) + " with product code " + product.ProductCode + " is already exists")
		}

		productData := model.Product{}
		if err = structs.Automapper(request, &productData); err != nil {
			return response, err
		}

		if custId == "" {
			custId = request.CustId
		}

		if userId == 0 {
			userId = request.CreatedBy
		}

		productData.CustId = request.CustId
		productData.ProductCode = request.ProductCode
		productData.ProductName = request.ProductName
		productData.IsActive = request.IsActive
		productData.CreatedAt = &timeNow
		productData.CreatedBy = &request.CreatedBy
		productData.UpdatedAt = &timeNow
		productData.UpdatedBy = &request.CreatedBy

		products = append(products, productData)
	}

	productIds, err := service.ProductRepository.BulkStore(products)
	if err != nil {
		return response, err
	}

	go service.BulkStoreDistProduct(custId, userId, productIds)

	for _, productId := range productIds {
		productResponse := entity.ProductResponse{}
		productResponse.ProductId = productId
		response.Products = append(response.Products, productResponse)
	}

	return response, err
}
func (service *productServiceImpl) BulkStoreDistProduct(DistID string, userID int64, productIDs []int64) (err error) {
	distributors, err := service.ProductRepository.FindDistributor(DistID)
	log.Info("Distributors :  ", distributors)
	if err != nil {
		return err
	}
	for _, productID := range productIDs {
		pid := productID // local copy to avoid loop-variable capture on Go <1.22
		var distProducts []model.ProductDistCreate
		for _, distributor := range distributors {
			distProduct := model.ProductDistCreate{
				CustId:         distributor.CustId,
				ProId:          &pid,
				IsActive:       false,
				IsAlloc:        false,
				MinStock:       0,
				MinStockStr:    "00000.000.000",
				SafetyStock:    0,
				SafetyStockStr: "00000.000.000",
				PoFormula:      1,
				IsNewPro:       true,
				Vat:            11.0,
				VatBg:          0.00,
				VatLgPurch:     0.00,
				VatLgSell:      0.00,
				Cogs:           0.0000,
				SMweek1:        1,
				SMweek2:        1,
				UpdatedAt:      time.Now(),
				UpdatedBy:      userID,
				ParentProId:    &pid,
			}
			distProducts = append(distProducts, distProduct)
		}
		if len(distProducts) > 0 {
			err = service.ProductRepository.StoreDist(distProducts)
			if err != nil {
				log.Info("BulkStoreDistProduct, BulkStoreDist, err:", err.Error())
				return err
			}
		}
	}

	return nil
}
func (service *productServiceImpl) Update(productId int64, request entity.UpdateProductRequest) (err error) {

	// product_type & cust id validation, if err == nil and params productId != product.Id, this means that code & cust id already exists
	product, err := service.ProductRepository.FindOneByProductCodeAndCustId(request.ProductCode, request.CustId)
	if err == nil && product.ProductId != productId {
		return errors.New("product_code: " + product.ProductCode + " is already exists")
	}

	// requestPrint, _ := json.Marshal(request)
	// log.Info("### ProductService, Update, requestPrint ###")
	// log.Info(string(requestPrint))
	// log.Info("### End Of requestPrint ###")

	err = service.ProductRepository.Update(productId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *productServiceImpl) Delete(custId string, productId int64, userId int64) (err error) {

	err = service.ProductRepository.Delete(custId, productId, userId)
	if err != nil {
		return err
	}

	return err
}

func (service *productServiceImpl) DeleteMultiple(custId string, productId []int64, userId int64) (err error) {

	err = service.ProductRepository.DeleteMultiple(custId, productId, userId)
	if err != nil {
		return err
	}

	return err
}

func (service *productServiceImpl) PrincipalList(dataFilter entity.ProductPrincipalQueryFilter, custId string) (data []entity.PrincipalLookupResponse, total int, lastPage int, err error) {
	principals, total, lastPage, err := service.ProductRepository.FindAllPrincipal(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	data = make([]entity.PrincipalLookupResponse, 0)
	for _, row := range principals {
		var vResp entity.PrincipalLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *productServiceImpl) CategoryList(dataFilter entity.ProductCategoryQueryFilter, custId string) (data []entity.ProductCategoryList, total int, lastPage int, err error) {
	productCategories, total, lastPage, err := service.ProductRepository.FindAllCategory(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	log.Info("productCategories:", structs.StructToJson(productCategories))

	data = make([]entity.ProductCategoryList, 0)
	for _, row := range productCategories {
		var vResp entity.ProductCategoryList
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *productServiceImpl) BrandList(dataFilter entity.ProductBrandQueryFilter, custId string) (data []entity.ProductBrandList, total int, lastPage int, err error) {
	brands, total, lastPage, err := service.ProductRepository.FindAllBrand(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	log.Info("brands:", structs.StructToJson(brands))

	data = make([]entity.ProductBrandList, 0)
	for _, row := range brands {
		var vResp entity.ProductBrandList
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
