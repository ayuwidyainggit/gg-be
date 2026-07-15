package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"master/entity"
	"master/model"
	"master/repository"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrSurveyReportInvalidDateRange = errors.New("end_date must be after or equal to start_date")

type SurveyReportService interface {
	List(filter entity.SurveyReportQueryFilter) ([]entity.SurveyReportListResponse, int, int, error)
	Detail(surveyAnswerID int64, custID string) (entity.SurveyReportDetailResponse, error)
	Export(filter entity.SurveyReportQueryFilter, createdBy string) (entity.SurveyReportExportResponse, error)
}

type surveyReportServiceImpl struct {
	surveyReportRepo  repository.SurveyReportRepository
	reportListRepo    repository.ReportListRepository
	objectURLResolver objectURLResolver
}

type objectURLResolver interface {
	ResolvePublicFileURL(objectKey string) (string, bool)
}

const (
	surveyReportDefaultPage             = 1
	surveyReportDefaultLimit            = 5
	surveyReportDefaultSort             = "created_date:desc"
	surveyReportExportInProgressMessage = "Processing time may vary by file size. Please check Download History to access the file"
)

func NewSurveyReportService(surveyReportRepo repository.SurveyReportRepository, reportListRepo repository.ReportListRepository, objectURLResolver objectURLResolver) SurveyReportService {
	return &surveyReportServiceImpl{
		surveyReportRepo:  surveyReportRepo,
		reportListRepo:    reportListRepo,
		objectURLResolver: objectURLResolver,
	}
}

func (s *surveyReportServiceImpl) List(filter entity.SurveyReportQueryFilter) ([]entity.SurveyReportListResponse, int, int, error) {
	if err := normalizeSurveyReportFilter(&filter); err != nil {
		return nil, 0, 0, err
	}

	rows, total, lastPage, err := s.surveyReportRepo.FindList(filter, true)
	if err != nil {
		return nil, 0, 0, err
	}

	responses := make([]entity.SurveyReportListResponse, 0, len(rows))
	for _, row := range rows {
		responses = append(responses, entity.SurveyReportListResponse{
			SurveyAnswerID:     row.SurveyAnswerID,
			SurveyID:           row.SurveyID,
			SurveyTitle:        row.SurveyTitle,
			AnswerFrequency:    row.AnswerFrequency,
			ResponseType:       row.ResponseType,
			AnswerDate:         row.AnswerDate,
			CreatedDate:        row.CreatedDate,
			EffectiveDateStart: row.EffectiveDateStart,
			EffectiveDateEnd:   row.EffectiveDateEnd,
			AreaID:             row.AreaID,
			AreaCode:           row.AreaCode,
			AreaName:           row.AreaName,
			DistributorID:      row.DistributorID,
			DistributorCode:    row.DistributorCode,
			DistributorName:    row.DistributorName,
			OutletID:           row.OutletID,
			OutletCode:         row.OutletCode,
			OutletName:         row.OutletName,
			EmpID:              row.EmpID,
			EmpCode:            row.EmpCode,
			EmpName:            row.EmpName,
			SalesmanName:       row.SalesmanName,
			Status:             row.Status,
		})
	}

	return responses, total, lastPage, nil
}

func (s *surveyReportServiceImpl) Detail(surveyAnswerID int64, custID string) (entity.SurveyReportDetailResponse, error) {
	var response entity.SurveyReportDetailResponse

	row, err := s.surveyReportRepo.FindDetail(surveyAnswerID, custID)
	if err != nil {
		return response, err
	}

	questionRows, err := s.surveyReportRepo.FindQuestionDetails(surveyAnswerID, row.AnswerCustID)
	if err != nil {
		return response, err
	}

	questionTemplateIDs := make([]int64, 0, len(questionRows))
	detailIDs := make([]int64, 0, len(questionRows))
	for _, questionRow := range questionRows {
		questionTemplateIDs = append(questionTemplateIDs, questionRow.QuestionTemplateID)
		detailIDs = append(detailIDs, questionRow.SurveyAnswerDetailID)
	}

	questionOptions, err := s.surveyReportRepo.FindQuestionOptions(questionTemplateIDs)
	if err != nil {
		return response, err
	}

	selectedOptions, err := s.surveyReportRepo.FindSelectedOptions(detailIDs, row.AnswerCustID)
	if err != nil {
		return response, err
	}

	answerFiles, err := s.surveyReportRepo.FindAnswerFiles(detailIDs, row.AnswerCustID)
	if err != nil {
		return response, err
	}

	optionsMap := make(map[int64][]entity.SurveyReportQuestionOption)
	for _, optionRow := range questionOptions {
		optionsMap[optionRow.QuestionTemplateID] = append(optionsMap[optionRow.QuestionTemplateID], entity.SurveyReportQuestionOption{
			QOptionTemplateID: optionRow.QOptionTemplateID,
			Option:            optionRow.Option,
		})
	}

	selectedMap := make(map[int64][]entity.SurveyReportSelectedOption)
	for _, selectedRow := range selectedOptions {
		selectedMap[selectedRow.SurveyAnswerDetailID] = append(selectedMap[selectedRow.SurveyAnswerDetailID], entity.SurveyReportSelectedOption{
			SurveyAnswerOptionID: selectedRow.SurveyAnswerOptionID,
			QOptionTemplateID:    selectedRow.QOptionTemplateID,
			OptionLabel:          selectedRow.OptionLabel,
		})
	}

	filesMap := make(map[int64][]entity.SurveyReportAnswerFile)
	for _, fileRow := range answerFiles {
		filesMap[fileRow.SurveyAnswerDetailID] = append(filesMap[fileRow.SurveyAnswerDetailID], entity.SurveyReportAnswerFile{
			SurveyAnswerFilesID: fileRow.SurveyAnswerFilesID,
			FileName:            fileRow.FileName,
			FileKey:             fileRow.FileKey,
			MediaCategory:       fileRow.MediaCategory,
			FileSize:            fileRow.FileSize,
		})
	}

	questions := make([]entity.SurveyReportQuestion, 0, len(questionRows))
	for _, questionRow := range questionRows {
		question := entity.SurveyReportQuestion{
			SurveyAnswerDetailID: questionRow.SurveyAnswerDetailID,
			QuestionTemplateID:   questionRow.QuestionTemplateID,
			SurveyTemplateID:     questionRow.SurveyTemplateID,
			Question:             questionRow.Question,
			InputType:            questionRow.InputType,
			AnswerType:           questionRow.AnswerType,
			Answer:               surveyReportResolveAnswer(questionRow, selectedMap[questionRow.SurveyAnswerDetailID]),
			Seq:                  questionRow.Seq,
			IsAnswered:           questionRow.IsAnswered,
			FreeTextAnswer:       questionRow.FreeTextAnswer,
			PhotoPath:            questionRow.PhotoPath,
			Options:              optionsMap[questionRow.QuestionTemplateID],
			SelectedOptions:      selectedMap[questionRow.SurveyAnswerDetailID],
			Files:                filesMap[questionRow.SurveyAnswerDetailID],
		}
		if question.Options == nil {
			question.Options = []entity.SurveyReportQuestionOption{}
		}
		if question.SelectedOptions == nil {
			question.SelectedOptions = []entity.SurveyReportSelectedOption{}
		}
		if question.Files == nil {
			question.Files = []entity.SurveyReportAnswerFile{}
		}
		questions = append(questions, question)
	}

	response = entity.SurveyReportDetailResponse{
		SurveyAnswerID:     row.SurveyAnswerID,
		SurveyID:           row.SurveyID,
		SurveyTitle:        row.SurveyTitle,
		AnswerFrequency:    row.AnswerFrequency,
		ResponseType:       row.ResponseType,
		AnswerDate:         row.AnswerDate,
		CreatedDate:        row.CreatedDate,
		EffectiveDateStart: row.EffectiveDateStart,
		EffectiveDateEnd:   row.EffectiveDateEnd,
		Area: entity.SurveyReportAreaResponse{
			AreaID:   row.AreaID,
			AreaCode: row.AreaCode,
			AreaName: row.AreaName,
		},
		Distributor: entity.SurveyReportDistributorResponse{
			DistributorID:   row.DistributorID,
			DistributorCode: row.DistributorCode,
			DistributorName: row.DistributorName,
		},
		Outlet: entity.SurveyReportOutletResponse{
			OutletID:   row.OutletID,
			OutletCode: row.OutletCode,
			OutletName: row.OutletName,
		},
		Salesman: entity.SurveyReportSalesmanResponse{
			EmpID:        row.EmpID,
			EmpCode:      row.EmpCode,
			EmpName:      row.EmpName,
			SalesmanName: row.SalesmanName,
		},
		Status:  row.Status,
		Details: questions,
	}
	if response.Details == nil {
		response.Details = []entity.SurveyReportQuestion{}
	}

	return response, nil
}

func (s *surveyReportServiceImpl) Export(filter entity.SurveyReportQueryFilter, createdBy string) (entity.SurveyReportExportResponse, error) {
	var response entity.SurveyReportExportResponse

	if err := normalizeSurveyReportFilter(&filter); err != nil {
		return response, err
	}

	const reportPrefix = "DownloadSurveyReport"
	inProgress, err := s.reportListRepo.CountInProgress(reportPrefix, filter.CustID, entity.FILE_STATUS_PROCESSING)
	if err != nil {
		return response, err
	}
	if inProgress > 0 {
		return response, errors.New(surveyReportExportInProgressMessage)
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.UTC
	}
	nowUTC := time.Now().UTC()
	now := nowUTC.In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)
	sequence, err := s.reportListRepo.CountByPrefixAndDate(reportPrefix, filter.CustID, startOfDay, endOfDay)
	if err != nil {
		return response, err
	}

	reportID := primitive.NewObjectID().Hex()
	reportName := fmt.Sprintf("%s-%s-%03d", reportPrefix, now.Format("020106"), sequence+1)
	createdAt := nowUTC
	createdByValue := createdBy
	storedReport := model.ReportList{
		ReportID:   reportID,
		CustID:     filter.CustID,
		ReportName: reportName,
		StartDate:  filter.StartDate,
		EndDate:    filter.EndDate,
		FileStatus: entity.FILE_STATUS_PROCESSING,
		CreatedBy:  &createdByValue,
		CreatedAt:  &createdAt,
	}
	if err := s.reportListRepo.Store(storedReport); err != nil {
		return response, err
	}

	rows, err := s.surveyReportRepo.FindExportRows(filter)
	if err != nil {
		_ = s.reportListRepo.UpdateFileResult(reportID, entity.FILE_STATUS_FAILED, nil, nil)
		return response, err
	}

	fileBase64, err := generateSurveyReportExcel(rows, s.objectURLResolver)
	if err != nil {
		_ = s.reportListRepo.UpdateFileResult(reportID, entity.FILE_STATUS_FAILED, nil, nil)
		return response, err
	}

	if err := s.reportListRepo.UpdateFileResult(reportID, entity.FILE_STATUS_READY, &fileBase64, nil); err != nil {
		return response, err
	}

	response = entity.SurveyReportExportResponse{
		ReportID:   reportID,
		ReportName: reportName,
		FileStatus: entity.FILE_STATUS_READY,
		FileBase64: fileBase64,
		CreatedBy:  createdBy,
	}
	response.FileStatusName = response.GetFileStatusName()

	return response, nil
}

func normalizeSurveyReportFilter(filter *entity.SurveyReportQueryFilter) error {
	if filter.Page < 1 {
		filter.Page = surveyReportDefaultPage
	}
	if filter.Limit < 1 {
		filter.Limit = surveyReportDefaultLimit
	}
	if filter.Limit > 9999 {
		filter.Limit = 9999
	}
	if filter.Sort == "" {
		filter.Sort = surveyReportDefaultSort
	}

	if filter.StartDate != nil && filter.EndDate != nil && filter.EndDate.Before(*filter.StartDate) {
		return ErrSurveyReportInvalidDateRange
	}
	return nil
}

func generateSurveyReportExcel(rows []model.SurveyReportExportRow, resolver objectURLResolver) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Survey Report"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{
		"Survey Date",
		"Survey Title",
		"Area Code",
		"Area Name",
		"Distributor Code",
		"Distributor Name",
		"Employee Code",
		"Employee Name",
		"Outlet Code",
		"Outlet Name",
		"Question",
		"Answer",
		"Attachment 1",
		"Attachment 2",
		"Attachment 3",
	}

	headStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#D9EAD3"}, Pattern: 1},
	})

	for idx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(idx+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headStyle)
	}

	for rowIdx, row := range rows {
		currentRow := rowIdx + 2
		values := []any{
			surveyReportFormatNullableDate(row.SurveyDate),
			row.SurveyTitle,
			row.AreaCode,
			row.AreaName,
			row.DistributorCode,
			row.DistributorName,
			row.EmpCode,
			row.EmpName,
			row.OutletCode,
			row.OutletName,
			row.Question,
			row.Answer,
			row.Attachment1,
			row.Attachment2,
			row.Attachment3,
		}

		for colIdx, value := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, currentRow)
			if colIdx >= 12 {
				if attachment, ok := value.(string); ok {
					if urlValue, display, valid := resolveSurveyReportAttachment(resolver, attachment); valid {
						f.SetCellValue(sheetName, cell, display)
						if err := f.SetCellHyperLink(sheetName, cell, urlValue, "External"); err != nil {
							return "", err
						}
						continue
					}
				}
			}
			f.SetCellValue(sheetName, cell, value)
		}
	}

	for idx := range headers {
		col, _ := excelize.ColumnNumberToName(idx + 1)
		f.SetColWidth(sheetName, col, col, 20)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func resolveSurveyReportAttachment(resolver objectURLResolver, raw string) (string, string, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || resolver == nil {
		return "", "", false
	}

	resolved, ok := resolver.ResolvePublicFileURL(trimmed)
	if !ok || resolved == "" {
		return "", "", false
	}

	display := trimmed
	if parsed, err := url.Parse(trimmed); err == nil && parsed.Scheme != "" && parsed.Host != "" {
		display = path.Base(parsed.Path)
	}
	if display == "/" || display == "." || display == "" {
		display = trimmed
	}
	return resolved, display, true
}

func surveyReportFormatNullableDate(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format("02/01/2006")
}

func surveyReportResolveAnswer(questionRow model.SurveyReportQuestionRow, selectedOptions []entity.SurveyReportSelectedOption) string {
	if questionRow.FreeTextAnswer != nil && *questionRow.FreeTextAnswer != "" {
		return *questionRow.FreeTextAnswer
	}

	labels := make([]string, 0, len(selectedOptions))
	for _, selectedOption := range selectedOptions {
		if selectedOption.OptionLabel != "" {
			labels = append(labels, selectedOption.OptionLabel)
		}
	}

	return strings.Join(labels, ", ")
}
