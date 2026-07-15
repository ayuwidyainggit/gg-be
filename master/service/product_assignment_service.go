package service

import (
	"bytes"
	"context"
	"fmt"
	"master/entity"
	"master/model"
	"master/repository"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/xuri/excelize/v2"
)

const (
	productAssignmentTypeAssignment       = "assignment"
	productAssignmentTypeRemoveAssignment = "remove_assignment"
)

type ProductAssignmentService interface {
	List(dataFilter entity.ProductAssignmentQueryFilter) ([]entity.ProductAssignmentResponse, int, int, error)
	DownloadTemplate() (*bytes.Buffer, string, string, error)
	Export(dataFilter entity.ProductAssignmentQueryFilter) (*bytes.Buffer, string, string, error)
	Import(req entity.ProductAssignmentImportRequest, custId string, createdBy int64) (entity.ProductAssignmentImportResponse, error)
	RemoveAssignment(req entity.ProductAssignmentImportRequest, custId string, createdBy int64) (entity.ProductAssignmentImportResponse, error)
}

func NewProductAssignmentService(productAssignmentRepository repository.ProductAssignmentRepository, productRepository repository.ProductRepository, distributorRepository repository.DistributorRepository, tx repository.TransactionManager) *productAssignmentServiceImpl {
	return &productAssignmentServiceImpl{
		tx:                          tx,
		ProductAssignmentRepository: productAssignmentRepository,
		ProductRepository:           productRepository,
		DistributorRepository:       distributorRepository,
	}
}

type productAssignmentServiceImpl struct {
	tx                          repository.TransactionManager
	ProductAssignmentRepository repository.ProductAssignmentRepository
	ProductRepository           repository.ProductRepository
	DistributorRepository       repository.DistributorRepository
}

func (s *productAssignmentServiceImpl) List(dataFilter entity.ProductAssignmentQueryFilter) ([]entity.ProductAssignmentResponse, int, int, error) {
	assignments, total, lastPage, err := s.ProductAssignmentRepository.FindAll(dataFilter)
	if err != nil {
		log.Error("ProductAssignmentService, List, err:", err.Error())
		return nil, 0, 0, err
	}

	responses := make([]entity.ProductAssignmentResponse, 0, len(assignments))
	for _, assignment := range assignments {
		response := entity.ProductAssignmentResponse{
			ID:              assignment.ID,
			CustID:          assignment.CustID,
			ProID:           assignment.ProID,
			ProCode:         assignment.ProCode,
			ProName:         assignment.ProName,
			DistributorID:   assignment.DistributorID,
			DistributorCode: assignment.DistributorCode,
			DistributorName: assignment.DistributorName,
			AssignmentType:  assignment.AssignmentType,
			CreatedBy:       assignment.CreatedBy,
			CreatedByName:   assignment.CreatedByName,
			CreatedAt:       assignment.CreatedAt.Format(time.RFC3339),
		}
		responses = append(responses, response)
	}

	return responses, total, lastPage, nil
}

func (s *productAssignmentServiceImpl) DownloadTemplate() (*bytes.Buffer, string, string, error) {
	f := excelize.NewFile()
	sheet := "Product_Assignment_Template"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)

	// Header
	headers := []string{"pro_code", "distributor_code"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	f.DeleteSheet("Sheet1")
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", "", err
	}

	return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "product_assignment_template.xlsx", nil
}

func (s *productAssignmentServiceImpl) Export(dataFilter entity.ProductAssignmentQueryFilter) (*bytes.Buffer, string, string, error) {
	assignments, err := s.ProductAssignmentRepository.FindAllHistory(dataFilter)
	if err != nil {
		log.Error("ProductAssignmentService, Export, FindAllHistory err:", err.Error())
		return nil, "", "", err
	}

	f := excelize.NewFile()
	sheet := "Product_Assignment"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)

	headers := []string{"Product Code", "Product Name", "Distributor Code", "Distributor", "Type", "Created by"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for r, assignment := range assignments {
		row := r + 2
		f.SetCellValue(sheet, "A"+fmt.Sprint(row), assignment.ProCode)
		f.SetCellValue(sheet, "B"+fmt.Sprint(row), assignment.ProName)
		f.SetCellValue(sheet, "C"+fmt.Sprint(row), assignment.DistributorCode)
		f.SetCellValue(sheet, "D"+fmt.Sprint(row), assignment.DistributorName)
		f.SetCellValue(sheet, "E"+fmt.Sprint(row), humanizeAssignmentType(assignment.AssignmentType))
		f.SetCellValue(sheet, "F"+fmt.Sprint(row), formatCreatedByLabel(assignment.CreatedByName, assignment.CreatedAt))
	}

	f.DeleteSheet("Sheet1")
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", "", err
	}

	return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "product_assignment.xlsx", nil
}

func (s *productAssignmentServiceImpl) Import(req entity.ProductAssignmentImportRequest, custId string, createdBy int64) (entity.ProductAssignmentImportResponse, error) {
	return s.processAssignmentFile(req, custId, createdBy, productAssignmentTypeAssignment)
}

func (s *productAssignmentServiceImpl) RemoveAssignment(req entity.ProductAssignmentImportRequest, custId string, createdBy int64) (entity.ProductAssignmentImportResponse, error) {
	return s.processAssignmentFile(req, custId, createdBy, productAssignmentTypeRemoveAssignment)
}

func (s *productAssignmentServiceImpl) processAssignmentFile(req entity.ProductAssignmentImportRequest, custId string, createdBy int64, assignmentType string) (entity.ProductAssignmentImportResponse, error) {
	response := entity.ProductAssignmentImportResponse{
		FileURL:     req.FileURL,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}

	rows, err := downloadAssignmentRows(req.FileURL)
	if err != nil {
		return response, err
	}

	var failedReasons []string
	var successCount int
	var total int
	hasErrors := false

	ctx := context.Background()
	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		for rowIndex, row := range rows {
			if rowIndex == 0 {
				continue
			}
			total++

			proCode := ""
			if len(row) > 0 {
				proCode = strings.TrimSpace(row[0])
			}
			distributorCode := ""
			if len(row) > 1 {
				distributorCode = strings.TrimSpace(row[1])
			}

			if proCode == "" || distributorCode == "" {
				failedReasons = append(failedReasons, fmt.Sprintf("row %d: missing pro_code or distributor_code", rowIndex+1))
				hasErrors = true
				continue
			}

			principalProduct, err := s.ProductRepository.FindOneByProductCodeAndCustId(proCode, custId)
			if err != nil {
				failedReasons = append(failedReasons, fmt.Sprintf("row %d: product %s not found", rowIndex+1, proCode))
				hasErrors = true
				continue
			}

			dist, err := s.DistributorRepository.FindOneByParentCustIdAndDistributorCode(custId, distributorCode)
			if err != nil {
				failedReasons = append(failedReasons, fmt.Sprintf("row %d: distributor %s not found", rowIndex+1, distributorCode))
				hasErrors = true
				continue
			}

			now := time.Now()

			switch assignmentType {
			case productAssignmentTypeAssignment:
				existingDistProduct, _ := s.ProductRepository.FindOneByProductCodeAndCustId(proCode, dist.DistributorCustID)
				if existingDistProduct.ProductId != 0 {
					failedReasons = append(failedReasons, fmt.Sprintf("row %d: product %s already assigned to distributor %s", rowIndex+1, proCode, distributorCode))
					hasErrors = true
					continue
				}

				productReplica := principalProduct
				productReplica.ParentProId = int(principalProduct.ProductId)
				productReplica.CustId = dist.DistributorCustID
				productReplica.DistributorID = &dist.DistributorId
				productReplica.Level = 1
				productReplica.Origin = productAssignmentTypeAssignment
				productReplica.AssignerUserID = &createdBy
				productReplica.CreatedAt = &now
				productReplica.UpdatedAt = &now

				_, err = s.ProductRepository.Store(txCtx, productReplica)
				if err != nil {
					failedReasons = append(failedReasons, fmt.Sprintf("row %d: failed to replicate product %s: %v", rowIndex+1, proCode, err))
					hasErrors = true
					continue
				}
			case productAssignmentTypeRemoveAssignment:
				existingDistProduct, _ := s.ProductRepository.FindOneByProductCodeAndCustId(proCode, dist.DistributorCustID)
				if existingDistProduct.ProductId == 0 {
					failedReasons = append(failedReasons, fmt.Sprintf("row %d: product %s is not assigned to distributor %s", rowIndex+1, proCode, distributorCode))
					hasErrors = true
					continue
				}

				err = s.ProductRepository.DeleteWithContext(txCtx, dist.DistributorCustID, existingDistProduct.ProductId, createdBy)
				if err != nil {
					failedReasons = append(failedReasons, fmt.Sprintf("row %d: failed to remove assignment for product %s: %v", rowIndex+1, proCode, err))
					hasErrors = true
					continue
				}
			}

			assignmentLog := model.ProductAssignmentMutation{
				CustID:         custId,
				ActionDate:     now,
				ProID:          principalProduct.ProductId,
				DistributorID:  dist.DistributorId,
				AssignmentType: assignmentType,
				CreatedBy:      createdBy,
				CreatedAt:      now,
			}
			_, err = s.ProductAssignmentRepository.Store(txCtx, assignmentLog)
			if err != nil {
				failedReasons = append(failedReasons, fmt.Sprintf("row %d: failed to record assignment history for product %s: %v", rowIndex+1, proCode, err))
				hasErrors = true
				continue
			}

			successCount++
		}

		if hasErrors {
			return fmt.Errorf("import has %d failed rows, all changes rolled back", total-successCount)
		}

		return nil
	})

	response.TotalRow = total
	response.SuccessRow = successCount
	response.FailedRow = total - successCount
	response.FailedReasons = failedReasons

	if err != nil {
		log.Error(err)
		return response, err
	}

	return response, nil
}

func downloadAssignmentRows(fileURL string) ([][]string, error) {
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

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("excel file has no sheets")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func humanizeAssignmentType(value string) string {
	switch strings.ToLower(value) {
	case productAssignmentTypeRemoveAssignment:
		return "Remove Assignment"
	default:
		return "Product Assignment"
	}
}

func formatCreatedByLabel(name string, createdAt time.Time) string {
	displayName := strings.TrimSpace(name)
	if displayName == "" {
		displayName = "-"
	}
	return fmt.Sprintf("%s | %s", displayName, createdAt.Format("15:04, 2 Jan 2006"))
}
