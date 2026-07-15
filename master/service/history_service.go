package service

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/xuri/excelize/v2"
	"master/entity"
	"master/repository"
)

type HistoryService interface {
	List(uploadType, search *string, page, limit int, custId string) ([]entity.ImportHistoryRow, int, int, error)
	Download(custId string, historyId int64, format string) (*bytes.Buffer, string, string, error)
	GetDetail(historyId int64) (entity.ImportHistoryRow, error)
	ReuploadFile(custId string, historyId int64, format string, file entity.ImportRequest) error
}

type historyServiceImpl struct {
	repo            repository.HistoryRepository
	outletSvc       OutletService
	productService  ProductService
	employeeService EmployeeService
}

func NewHistoryService(repo repository.HistoryRepository, outletSvc OutletService, productService ProductService, employeeService EmployeeService) HistoryService {
	return &historyServiceImpl{repo: repo, outletSvc: outletSvc, productService: productService, employeeService: employeeService}
}

func (s *historyServiceImpl) List(uploadType, search *string, page, limit int, custId string) ([]entity.ImportHistoryRow, int, int, error) {
	return s.repo.ListImportHistory(uploadType, search, page, limit, custId)
}

func (s *historyServiceImpl) Download(custId string, historyId int64, format string) (*bytes.Buffer, string, string, error) {
	hist, err := s.repo.GetImportHistoryById(historyId)
	if err != nil {
		return nil, "", "", err
	}
	if hist.UploadType == nil || *hist.UploadType == "" {
		return nil, "", "", errors.New("upload_type is empty for this history")
	}

	uploadTypeSlug := strings.ToLower(strings.TrimSpace(*hist.UploadType))
	if uploadTypeSlug == "" {
		uploadTypeSlug = "unknown"
	} else {
		uploadTypeSlug = strings.NewReplacer(" ", "_", "/", "-", "\\", "-", ".", "_").Replace(uploadTypeSlug)
	}

	timestamp := time.Now().Format("2006-01-02_1504")
	failedCSVFilename := fmt.Sprintf("history_%d_%s_failed_%s.csv", historyId, uploadTypeSlug, timestamp)
	failedZIPFilename := fmt.Sprintf("history_%d_%s_failed_%s.zip", historyId, uploadTypeSlug, timestamp)
	failedExcelFilename := fmt.Sprintf("history_%d_%s_failed_%s.xlsx", historyId, uploadTypeSlug, timestamp)

	table, includeCtid, err := tempTableByUploadType(*hist.UploadType)
	if err != nil {
		return nil, "", "", err
	}

	cols, rows, err := s.repo.FetchTempRows(table, historyId, custId, includeCtid)
	if err != nil {
		return nil, "", "", err
	}

	// For outlet reupload insert only, convert raw column names to display headers
	isOutletNew := isOutletNewUploadType(*hist.UploadType)
	isOutlet := strings.ToLower(*hist.UploadType) == "outlet"
	isOutletInsert := isOutlet || isOutletNew
	isOutletUpdate := strings.ToLower(*hist.UploadType) == "outlet-update"
	displayCols := cols
	if isOutletInsert {
		displayCols = make([]string, len(cols))
		for i := range cols {
			if isOutletNew {
				displayCols[i] = importNewDisplayHeader(cols[i])
			} else {
				displayCols[i] = toDisplayHeader(cols[i])
			}
		}
	}
	isProduct := strings.ToLower(*hist.UploadType) == "product"
	if isProduct {
		displayCols = make([]string, len(cols))
		for i := range cols {
			displayCols[i] = MapHeaderToWeb(cols[i])
		}
	}
	isEmployee := strings.ToLower(*hist.UploadType) == "employee"
	if isEmployee {
		displayCols = make([]string, len(cols))
		for i := range cols {
			displayCols[i] = employeeDisplayHeader(cols[i])
		}
	}
	if isOutletUpdate {
		displayCols = make([]string, len(cols))
		for i := range cols {
			displayCols[i] = toDisplayHeader(cols[i])
		}
	}
	isProductUpdate := strings.ToLower(*hist.UploadType) == "product-update"
	if strings.ToLower(*hist.UploadType) == "product-update" {
		displayCols = make([]string, len(cols))
		for i := range cols {
			displayCols[i] = MapHeaderToWeb(cols[i])
		}
	}
	isEmployeeUpdate := strings.ToLower(*hist.UploadType) == "employee-update"
	if strings.ToLower(*hist.UploadType) == "employee-update" {
		displayCols = make([]string, len(cols))
		for i := range cols {
			displayCols[i] = employeeDisplayHeader(cols[i])
		}
	}

	// Hide cust_id, ctid, history_id for product/product-update downloads (both history and reupload)
	if strings.ToLower(*hist.UploadType) == "product" || strings.ToLower(*hist.UploadType) == "product-update" {
		// determine indices to keep based on original column names
		hidden := func(name string) bool {
			n := strings.ToLower(strings.TrimSpace(name))
			return n == "cust_id" || n == "ctid" || n == "history_id"
		}
		keepIdx := make([]int, 0, len(cols))
		for i, c := range cols {
			if !hidden(c) {
				keepIdx = append(keepIdx, i)
			}
		}
		// filter cols and displayCols
		newCols := make([]string, 0, len(keepIdx))
		newDisplay := make([]string, 0, len(keepIdx))
		for _, i := range keepIdx {
			newCols = append(newCols, cols[i])
			newDisplay = append(newDisplay, displayCols[i])
		}
		cols = newCols
		displayCols = newDisplay
		// filter rows
		newRows := make([][]interface{}, 0, len(rows))
		for _, r := range rows {
			// guard against unexpected row length
			rec := make([]interface{}, 0, len(keepIdx))
			for _, i := range keepIdx {
				if i >= 0 && i < len(r) {
					rec = append(rec, r[i])
				} else {
					rec = append(rec, "")
				}
			}
			newRows = append(newRows, rec)
		}
		rows = newRows
	}

	// Hide internal IDs for outlet/outlet-new/outlet-update downloads (both history and reupload)
	if isOutletInsert || strings.ToLower(*hist.UploadType) == "outlet-update" {
		// list of internal / technical columns to hide from user templates
		internalIDs := map[string]struct{}{
			"cust_id":           {},
			"ctid":              {},
			"history_id":        {},
			"outlet_id":         {},
			"outlet_contact_id": {},
			"bank_id":           {},
			"outlet_tax_id":     {},
			"outlet_bank_id":    {},
		}
		hidden := func(name string) bool { _, ok := internalIDs[strings.ToLower(strings.TrimSpace(name))]; return ok }
		keepIdx := make([]int, 0, len(cols))
		for i, c := range cols {
			if !hidden(c) {
				keepIdx = append(keepIdx, i)
			}
		}
		if len(keepIdx) != len(cols) { // only rebuild if something filtered
			newCols := make([]string, 0, len(keepIdx))
			newDisplay := make([]string, 0, len(keepIdx))
			for _, i := range keepIdx {
				newCols = append(newCols, cols[i])
				newDisplay = append(newDisplay, displayCols[i])
			}
			cols = newCols
			displayCols = newDisplay
			newRows := make([][]interface{}, 0, len(rows))
			for _, r := range rows {
				rec := make([]interface{}, 0, len(keepIdx))
				for _, i := range keepIdx {
					if i >= 0 && i < len(r) {
						rec = append(rec, r[i])
					} else {
						rec = append(rec, "")
					}
				}
				newRows = append(newRows, rec)
			}
			rows = newRows
		}
	}

	// Hide internal IDs for employee / employee-update (history + reupload)
	if strings.ToLower(*hist.UploadType) == "employee" || strings.ToLower(*hist.UploadType) == "employee-update" {
		// Enumerate known internal/technical columns. (If you need to hide more, add here.)
		internalEmp := map[string]struct{}{
			"cust_id": {}, "ctid": {}, "history_id": {},
			"emp_id": {}, "emp_grp_id": {}, "emp_type_id": {},
			"division_id": {}, "province_id": {}, "city_id": {},
			"sub_district_id": {}, "ward_id": {},
		}
		// Also hide any column that strictly ends with "_id" but keep human-readable display versions already mapped.
		isInternal := func(name string) bool {
			n := strings.ToLower(strings.TrimSpace(name))
			if _, ok := internalEmp[n]; ok {
				return true
			}
			if strings.HasSuffix(n, "_id") {
				return true
			}
			return false
		}
		keepIdx := make([]int, 0, len(cols))
		for i, c := range cols {
			if !isInternal(c) {
				keepIdx = append(keepIdx, i)
			}
		}
		if len(keepIdx) != len(cols) {
			newCols := make([]string, 0, len(keepIdx))
			newDisplay := make([]string, 0, len(keepIdx))
			for _, i := range keepIdx {
				newCols = append(newCols, cols[i])
				newDisplay = append(newDisplay, displayCols[i])
			}
			cols = newCols
			displayCols = newDisplay
			newRows := make([][]interface{}, 0, len(rows))
			for _, r := range rows {
				rec := make([]interface{}, 0, len(keepIdx))
				for _, i := range keepIdx {
					if i >= 0 && i < len(r) {
						rec = append(rec, r[i])
					} else {
						rec = append(rec, "")
					}
				}
				newRows = append(newRows, rec)
			}
			rows = newRows
		}
	}

	if strings.ToLower(format) == "csv" {
		// For outlet insert reupload: return a ZIP containing failed rows CSV + instructions.csv
		if isOutletInsert {
			zipBuf := new(bytes.Buffer)
			zw := zip.NewWriter(zipBuf)

			// 1) failed rows csv
			failedName := failedCSVFilename
			fw, err := zw.Create(failedName)
			if err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			cw := csv.NewWriter(fw)
			cw.Comma = ';'
			if err := cw.Write(displayCols); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			for _, r := range rows {
				rec := make([]string, len(cols))
				for i, v := range r {
					rec[i] = toString(v)
				}
				if err := cw.Write(rec); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			}
			cw.Flush()
			if err := cw.Error(); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}

			// 2) instructions.csv (kolom, mandatory, keterangan)
			if isOutlet {
				if instructions, err := s.repo.GetImportInstructions("outlet"); err == nil && len(instructions) > 0 {
					iw, err := zw.Create("instructions.csv")
					if err != nil {
						_ = zw.Close()
						return nil, "", "", err
					}
					icw := csv.NewWriter(iw)
					icw.Comma = ';'
					_ = icw.Write([]string{"kolom", "mandatory", "keterangan", "step", "color"})
					for _, it := range instructions {
						mand := "Kolom Tidak Wajib Diisi"
						if it.Mandatory {
							mand = "Kolom Wajib Diisi"
						}
						stepVal := ""
						if it.Step != nil {
							stepVal = *it.Step
						}
						stepLabel, colorLabel := instructionStepAndColor(it.Kolom, stepVal)
						_ = icw.Write([]string{it.Kolom, mand, it.Keterangan, stepLabel, colorLabel})
					}
					icw.Flush()
					if err := icw.Error(); err != nil {
						_ = zw.Close()
						return nil, "", "", err
					}
				}
			}

			if err := zw.Close(); err != nil {
				return nil, "", "", err
			}
			return zipBuf, "application/zip", failedZIPFilename, nil
		}

		// For outlet-update reupload: same behavior (ZIP: failed.csv + instructions.csv)
		if isOutletUpdate {
			zipBuf := new(bytes.Buffer)
			zw := zip.NewWriter(zipBuf)

			// 1) failed rows csv
			failedName := failedCSVFilename
			fw, err := zw.Create(failedName)
			if err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			cw := csv.NewWriter(fw)
			cw.Comma = ';'
			if err := cw.Write(displayCols); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			for _, r := range rows {
				rec := make([]string, len(cols))
				for i, v := range r {
					rec[i] = toString(v)
				}
				if err := cw.Write(rec); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			}
			cw.Flush()
			if err := cw.Error(); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}

			// 2) instructions.csv, try outlet_update first, fallback to outlet
			if instructions, err := s.repo.GetImportInstructions("outlet_update"); err == nil && len(instructions) > 0 {
				iw, err := zw.Create("instructions.csv")
				if err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
				icw := csv.NewWriter(iw)
				icw.Comma = ';'
				_ = icw.Write([]string{"kolom", "mandatory", "keterangan", "step", "color"})
				for _, it := range instructions {
					mand := "Kolom Tidak Wajib Diisi"
					if it.Mandatory {
						mand = "Kolom Wajib Diisi"
					}
					stepVal := ""
					if it.Step != nil {
						stepVal = *it.Step
					}
					stepLabel, colorLabel := instructionStepAndColor(it.Kolom, stepVal)
					_ = icw.Write([]string{it.Kolom, mand, it.Keterangan, stepLabel, colorLabel})
				}
				icw.Flush()
				if err := icw.Error(); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			} else if instructions, err := s.repo.GetImportInstructions("outlet"); err == nil && len(instructions) > 0 {
				iw, err := zw.Create("instructions.csv")
				if err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
				icw := csv.NewWriter(iw)
				icw.Comma = ';'
				_ = icw.Write([]string{"kolom", "mandatory", "keterangan", "step", "color"})
				for _, it := range instructions {
					mand := "Kolom Tidak Wajib Diisi"
					if it.Mandatory {
						mand = "Kolom Wajib Diisi"
					}
					stepVal := ""
					if it.Step != nil {
						stepVal = *it.Step
					}
					stepLabel, colorLabel := instructionStepAndColor(it.Kolom, stepVal)
					_ = icw.Write([]string{it.Kolom, mand, it.Keterangan, stepLabel, colorLabel})
				}
				icw.Flush()
				if err := icw.Error(); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			}

			if err := zw.Close(); err != nil {
				return nil, "", "", err
			}
			return zipBuf, "application/zip", failedZIPFilename, nil
		}

		// For outlet-update reupload: return a ZIP containing failed rows CSV + instructions.csv
		if strings.ToLower(*hist.UploadType) == "outlet-update" {
			zipBuf := new(bytes.Buffer)
			zw := zip.NewWriter(zipBuf)

			failedName := failedCSVFilename
			fw, err := zw.Create(failedName)
			if err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			cw := csv.NewWriter(fw)
			cw.Comma = ';'
			if err := cw.Write(displayCols); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			for _, r := range rows {
				rec := make([]string, len(cols))
				for i, v := range r {
					rec[i] = toString(v)
				}
				if err := cw.Write(rec); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			}
			cw.Flush()
			if err := cw.Error(); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}

			if instructions, err := s.repo.GetImportInstructions("outlet_update"); err == nil && len(instructions) > 0 {
				iw, err := zw.Create("instructions.csv")
				if err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
				icw := csv.NewWriter(iw)
				icw.Comma = ';'
				_ = icw.Write([]string{"kolom", "mandatory", "keterangan"})
				for _, it := range instructions {
					mand := "Kolom Tidak Wajib Diisi"
					if it.Mandatory {
						mand = "Kolom Wajib Diisi"
					}
					_ = icw.Write([]string{it.Kolom, mand, it.Keterangan})
				}
				icw.Flush()
				if err := icw.Error(); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			}

			if err := zw.Close(); err != nil {
				return nil, "", "", err
			}
			return zipBuf, "application/zip", failedZIPFilename, nil
		}

		if isProduct {
			zipBuf := new(bytes.Buffer)
			zw := zip.NewWriter(zipBuf)

			// 1) failed rows csv
			failedName := failedCSVFilename
			fw, err := zw.Create(failedName)
			if err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			cw := csv.NewWriter(fw)
			cw.Comma = ';'
			if err := cw.Write(displayCols); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			for _, r := range rows {
				rec := make([]string, len(cols))
				for i, v := range r {
					rec[i] = toString(v)
				}
				if err := cw.Write(rec); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			}
			cw.Flush()
			if err := cw.Error(); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}

			// 2) instructions.csv (kolom, mandatory, keterangan)
			if instructions, err := s.repo.GetImportInstructions("product"); err == nil && len(instructions) > 0 {
				iw, err := zw.Create("instructions.csv")
				if err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
				icw := csv.NewWriter(iw)
				icw.Comma = ';'
				_ = icw.Write([]string{"kolom", "mandatory", "keterangan", "step", "color"})
				for _, it := range instructions {
					mand := "Kolom Tidak Wajib Diisi"
					if it.Mandatory {
						mand = "Kolom Wajib Diisi"
					}
					stepVal := ""
					if it.Step != nil {
						stepVal = *it.Step
					}
					stepLabel, colorLabel := instructionStepAndColor(it.Kolom, stepVal)
					_ = icw.Write([]string{it.Kolom, mand, it.Keterangan, stepLabel, colorLabel})
				}
				icw.Flush()
				if err := icw.Error(); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			}

			if err := zw.Close(); err != nil {
				return nil, "", "", err
			}
			return zipBuf, "application/zip", failedZIPFilename, nil
		}

		if isProductUpdate {
			zipBuf := new(bytes.Buffer)
			zw := zip.NewWriter(zipBuf)

			// 1) failed rows csv
			failedName := failedCSVFilename
			fw, err := zw.Create(failedName)
			if err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			cw := csv.NewWriter(fw)
			cw.Comma = ';'
			if err := cw.Write(displayCols); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			for _, r := range rows {
				rec := make([]string, len(cols))
				for i, v := range r {
					rec[i] = toString(v)
				}
				if err := cw.Write(rec); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			}
			cw.Flush()
			if err := cw.Error(); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}

			// 2) instructions.csv (kolom, mandatory, keterangan)
			if instructions, err := s.repo.GetImportInstructions("product"); err == nil && len(instructions) > 0 {
				iw, err := zw.Create("instructions.csv")
				if err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
				icw := csv.NewWriter(iw)
				icw.Comma = ';'
				_ = icw.Write([]string{"kolom", "mandatory", "keterangan", "step", "color"})
				for _, it := range instructions {
					mand := "Kolom Tidak Wajib Diisi"
					if it.Mandatory {
						mand = "Kolom Wajib Diisi"
					}
					stepVal := ""
					if it.Step != nil {
						stepVal = *it.Step
					}
					stepLabel, colorLabel := instructionStepAndColor(it.Kolom, stepVal)
					_ = icw.Write([]string{it.Kolom, mand, it.Keterangan, stepLabel, colorLabel})
				}
				icw.Flush()
				if err := icw.Error(); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			}

			if err := zw.Close(); err != nil {
				return nil, "", "", err
			}
			return zipBuf, "application/zip", failedZIPFilename, nil
		}

		if isEmployee {
			zipBuf := new(bytes.Buffer)
			zw := zip.NewWriter(zipBuf)

			failedName := failedCSVFilename
			fw, err := zw.Create(failedName)
			if err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			cw := csv.NewWriter(fw)
			cw.Comma = ';'
			if err := cw.Write(displayCols); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			for _, r := range rows {
				rec := make([]string, len(cols))
				for i, v := range r {
					rec[i] = toString(v)
				}
				if err := cw.Write(rec); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			}
			cw.Flush()
			if err := cw.Error(); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}

			if instructions, err := s.repo.GetImportInstructions("employee"); err == nil && len(instructions) > 0 {
				iw, err := zw.Create("instructions.csv")
				if err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
				icw := csv.NewWriter(iw)
				icw.Comma = ';'
				_ = icw.Write([]string{"kolom", "mandatory", "keterangan", "step", "color"})
				for _, it := range instructions {
					mand := "Kolom Tidak Wajib Diisi"
					if it.Mandatory {
						mand = "Kolom Wajib Diisi"
					}
					stepVal := ""
					if it.Step != nil {
						stepVal = *it.Step
					}
					stepLabel, colorLabel := instructionStepAndColor(it.Kolom, stepVal)
					_ = icw.Write([]string{it.Kolom, mand, it.Keterangan, stepLabel, colorLabel})
				}
				icw.Flush()
				if err := icw.Error(); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			}

			if err := zw.Close(); err != nil {
				return nil, "", "", err
			}
			return zipBuf, "application/zip", failedZIPFilename, nil
		}

		if isEmployeeUpdate {
			zipBuf := new(bytes.Buffer)
			zw := zip.NewWriter(zipBuf)

			failedName := failedCSVFilename
			fw, err := zw.Create(failedName)
			if err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			cw := csv.NewWriter(fw)
			cw.Comma = ';'
			if err := cw.Write(displayCols); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}
			for _, r := range rows {
				rec := make([]string, len(cols))
				for i, v := range r {
					rec[i] = toString(v)
				}
				if err := cw.Write(rec); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			}
			cw.Flush()
			if err := cw.Error(); err != nil {
				_ = zw.Close()
				return nil, "", "", err
			}

			if instructions, err := s.repo.GetImportInstructions("employee"); err == nil && len(instructions) > 0 {
				iw, err := zw.Create("instructions.csv")
				if err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
				icw := csv.NewWriter(iw)
				icw.Comma = ';'
				_ = icw.Write([]string{"kolom", "mandatory", "keterangan", "step", "color"})
				for _, it := range instructions {
					mand := "Kolom Tidak Wajib Diisi"
					if it.Mandatory {
						mand = "Kolom Wajib Diisi"
					}
					stepVal := ""
					if it.Step != nil {
						stepVal = *it.Step
					}
					stepLabel, colorLabel := instructionStepAndColor(it.Kolom, stepVal)
					_ = icw.Write([]string{it.Kolom, mand, it.Keterangan, stepLabel, colorLabel})
				}
				icw.Flush()
				if err := icw.Error(); err != nil {
					_ = zw.Close()
					return nil, "", "", err
				}
			}

			if err := zw.Close(); err != nil {
				return nil, "", "", err
			}
			return zipBuf, "application/zip", failedZIPFilename, nil
		}

		// Default: single CSV for other upload types
		buf := new(bytes.Buffer)
		w := csv.NewWriter(buf)
		w.Comma = ';'
		if err := w.Write(displayCols); err != nil {
			return nil, "", "", err
		}
		for _, r := range rows {
			rec := make([]string, len(cols))
			for i, v := range r {
				rec[i] = toString(v)
			}
			if err := w.Write(rec); err != nil {
				return nil, "", "", err
			}
		}
		w.Flush()
		if err := w.Error(); err != nil {
			return nil, "", "", err
		}
		return buf, "text/csv", failedCSVFilename, nil
	}

	f := excelize.NewFile()
	sheet := "Failed"
	idx, _ := f.NewSheet(sheet)
	for i, h := range displayCols {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	for rIdx, row := range rows {
		for cIdx, v := range row {
			cell, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+2)
			f.SetCellValue(sheet, cell, toString(v))
		}
	}
	// Optional: add Instructions sheet for reupload downloads
	if isOutlet {
		if instructions, err := s.repo.GetImportInstructions("outlet"); err == nil && len(instructions) > 0 {
			addInstructionSheet(f, instructions)
		}
	}
	if isOutletUpdate {
		if instructions, err := s.repo.GetImportInstructions("outlet_update"); err == nil && len(instructions) > 0 {
			addInstructionSheet(f, instructions)
		} else if instructions, err := s.repo.GetImportInstructions("outlet"); err == nil && len(instructions) > 0 {
			addInstructionSheet(f, instructions)
		}
	}
	if strings.ToLower(*hist.UploadType) == "outlet-update" {
		if instructions, err := s.repo.GetImportInstructions("outlet_update"); err == nil && len(instructions) > 0 {
			addInstructionSheet(f, instructions)
		}
	}
	if strings.ToLower(*hist.UploadType) == "product" {
		if instructions, err := s.repo.GetImportInstructions("product"); err == nil && len(instructions) > 0 {
			addInstructionSheet(f, instructions)
		}
	}
	if strings.ToLower(*hist.UploadType) == "product-update" {
		if instructions, err := s.repo.GetImportInstructions("product"); err == nil && len(instructions) > 0 {
			addInstructionSheet(f, instructions)
		}
	}
	if strings.ToLower(*hist.UploadType) == "employee" {
		if instructions, err := s.repo.GetImportInstructions("employee"); err == nil && len(instructions) > 0 {
			addInstructionSheet(f, instructions)
		}
	}
	if strings.ToLower(*hist.UploadType) == "employee-update" {
		if instructions, err := s.repo.GetImportInstructions("employee"); err == nil && len(instructions) > 0 {
			addInstructionSheet(f, instructions)
		}
	}
	f.SetActiveSheet(idx)
	_ = f.DeleteSheet("Sheet1")
	buf, err := f.WriteToBuffer()
	return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", failedExcelFilename, err
}

func (s *historyServiceImpl) GetDetail(historyId int64) (entity.ImportHistoryRow, error) {
	row, err := s.repo.GetImportHistoryById(historyId)
	if err != nil {
		return row, err
	}
	if row.UploadType == nil || strings.TrimSpace(*row.UploadType) == "" {
		return row, nil
	}

	table, _, err := tempTableByUploadType(*row.UploadType)
	if err != nil {
		return row, err
	}

	var errorMessages []string
	if table != "" {
		if messages, err := s.repo.GetDistinctTempValues(table, "error_message", historyId); err != nil {
			log.Warnf("failed to fetch error messages for history %d (%s): %v", historyId, *row.UploadType, err)
		} else if len(messages) > 0 {
			cleanMessages := make([]string, 0, len(messages))
			for _, msg := range messages {
				msg = strings.TrimSpace(msg)
				if msg != "" {
					cleanMessages = append(cleanMessages, msg)
				}
			}
			if len(cleanMessages) > 0 {
				row.ErrorMessages = cleanMessages
				errorMessages = cleanMessages
			}
		}

		if statuses, err := s.repo.GetDistinctTempValues(table, "status_insert", historyId); err != nil {
			log.Warnf("failed to fetch error statuses for history %d (%s): %v", historyId, *row.UploadType, err)
		} else if len(statuses) > 0 {
			humanStatuses := make([]string, 0, len(statuses))
			for _, st := range statuses {
				if friendly := humanizeStatus(st); friendly != "" {
					humanStatuses = append(humanStatuses, friendly)
				}
			}
			if len(humanStatuses) > 0 {
				row.ErrorStatuses = humanStatuses
			}
		}
	}

	if len(errorMessages) > 0 || len(row.ErrorStatuses) > 0 {
		header := ""
		if len(row.ErrorStatuses) > 0 {
			header = strings.Join(row.ErrorStatuses, ", ")
		}
		if len(errorMessages) > 0 {
			body := strings.Join(errorMessages, "\n")
			if header == "" {
				header = "validation failed"
			}
			logText := fmt.Sprintf("%s:\n%s", header, body)
			row.LogError = &logText
		} else if header != "" {
			logText := header
			row.LogError = &logText
		}
	}

	return row, nil
}

func (s *historyServiceImpl) ReuploadFile(custId string, historyId int64, format string, file entity.ImportRequest) error {
	log.Info("DISINI++++++++++++++++++++++++++++++++++++++++")
	hist, err := s.repo.GetImportHistoryById(historyId)
	if err != nil {
		return err
	}
	if hist.UploadType == nil || *hist.UploadType == "" {
		return errors.New("upload_type is empty")
	}
	switch strings.ToLower(*hist.UploadType) {
	case "outlet-update":
		return s.outletSvc.ReuploadImportUpdateFile(custId, historyId, file)
	case "product-update":
		return s.productService.ReuploadImportUpdateFile(custId, historyId, file)
	case "employee-update":
		return s.employeeService.ReuploadImportUpdateFile(custId, historyId, file)
	case "outlet":
		return s.outletSvc.ReuploadImportInsertFile(custId, historyId, file)
	case "outlet-new":
		file.IsImportNew = true
		return s.outletSvc.ReuploadImportInsertFile(custId, historyId, file)
	case "product":
		return s.productService.ReuploadImportInsertFile(custId, historyId, file)
	case "employee":
		return s.employeeService.ReuploadImportInsertFile(custId, historyId, file)
	default:
		return fmt.Errorf("reupload not implemented for upload_type %s", *hist.UploadType)
	}
}

func isOutletNewUploadType(uploadType string) bool {
	return strings.ToLower(strings.TrimSpace(uploadType)) == "outlet-new"
}

func tempTableByUploadType(uploadType string) (string, bool, error) {
	switch strings.ToLower(strings.TrimSpace(uploadType)) {
	case "outlet":
		return "import.outlet_temp", false, nil
	case "outlet-new":
		return "import.outlet_temp", false, nil
	case "product":
		return "import.product_temp", false, nil
	case "employee":
		return "import.employee_temp", false, nil
	case "outlet-update":
		return "import.outlet_update_temp", false, nil
	case "product-update":
		return "import.product_update_temp", false, nil
	case "employee-update":
		return "import.employee_update_temp", false, nil
	default:
		return "", false, fmt.Errorf("unknown upload_type: %s", uploadType)
	}
}

func humanizeStatus(status string) string {
	status = strings.TrimSpace(status)
	if status == "" {
		return ""
	}
	status = strings.ReplaceAll(status, "_", " ")
	status = strings.ReplaceAll(status, "-", " ")
	return strings.Join(strings.Fields(strings.ToLower(status)), " ")
}

func toString(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case []byte:
		return string(t)
	default:
		return fmt.Sprint(v)
	}
}

func addInstructionSheet(f *excelize.File, instructions []entity.ImportInstruction) {
	if len(instructions) == 0 {
		return
	}
	sheetName := "Instructions"
	if _, err := f.GetSheetIndex(sheetName); err == nil {
		for i := 2; i < 10; i++ {
			candidate := fmt.Sprintf("Instructions %d", i)
			if _, err := f.GetSheetIndex(candidate); err != nil {
				sheetName = candidate
				break
			}
		}
	}
	if _, err := f.NewSheet(sheetName); err != nil {
		return
	}
	_ = f.SetCellValue(sheetName, "A1", "kolom")
	_ = f.SetCellValue(sheetName, "B1", "mandatory")
	_ = f.SetCellValue(sheetName, "C1", "keterangan")
	_ = f.SetCellValue(sheetName, "D1", "step")
	_ = f.SetCellValue(sheetName, "E1", "color")

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#CCCCCC"}, // abu-abu muda
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	_ = f.SetCellStyle(sheetName, "A1", "E1", headerStyle)

	for i, ins := range instructions {
		row := i + 2
		mand := "Kolom Tidak Wajib Diisi"
		if ins.Mandatory {
			mand = "Kolom Wajib Diisi"
		}
		stepVal := ""
		if ins.Step != nil {
			stepVal = *ins.Step
		}
		stepLabel, colorLabel := instructionStepAndColor(ins.Kolom, stepVal)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), ins.Kolom)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), mand)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), ins.Keterangan)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), stepLabel)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), colorLabel)
	}
}
