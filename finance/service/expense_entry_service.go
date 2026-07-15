package service

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"finance/repository"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrExpenseExists            = errors.New("expense already exists")
	ErrExpenseNotFound          = errors.New("record not found")
	ErrInvalidExpenseDateFilter = errors.New("invalid expense date filter")
)

type ExpenseEntryService interface {
	List(filter entity.ExpenseEntryQueryFilter) (data []entity.ExpenseEntryListResponse, total int64, lastPage int, err error)
	Store(custId string, request entity.CreateExpenseEntryBody, userId int64) (entity.CreateExpenseResponseData, error)
	Update(custId string, expenseId int64, request entity.UpdateExpenseEntryBody, userId int64) (entity.UpdateExpenseResponseData, error)
	Delete(custId string, expenseId int64, deletedBy int64) error
	Detail(custId string, expenseId int64) (entity.ExpenseDetailResponse, error)
}

type expenseEntryServiceImpl struct {
	Repo        repository.ExpenseEntryRepository
	Transaction repository.Dbtransaction
}

func NewExpenseEntryService(repo repository.ExpenseEntryRepository, tx repository.Dbtransaction) *expenseEntryServiceImpl {
	return &expenseEntryServiceImpl{Repo: repo, Transaction: tx}
}

func (s *expenseEntryServiceImpl) List(filter entity.ExpenseEntryQueryFilter) ([]entity.ExpenseEntryListResponse, int64, int, error) {
	ctx := context.Background()

	if filter.StartDate == "" || filter.EndDate == "" {
		endDate := time.Now().UTC()
		startDate := endDate.AddDate(0, -3, 0)
		if filter.StartDate == "" {
			filter.StartDate = startDate.Format("2006-01-02")
		}
		if filter.EndDate == "" {
			filter.EndDate = endDate.Format("2006-01-02")
		}
	}

	if filter.StartDate != "" {
		normalizedStartDate, err := normalizeExpenseDateEpoch(filter.StartDate)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("%w: start_date", ErrInvalidExpenseDateFilter)
		}
		filter.StartDate = normalizedStartDate
	}
	if filter.EndDate != "" {
		normalizedEndDate, err := normalizeExpenseDateEpoch(filter.EndDate)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("%w: end_date", ErrInvalidExpenseDateFilter)
		}
		filter.EndDate = normalizedEndDate
	}

	rows, total, lastPage, err := s.Repo.FindAll(ctx, filter)
	if err != nil {
		return nil, 0, 0, err
	}

	resp := make([]entity.ExpenseEntryListResponse, 0, len(rows))
	for _, r := range rows {
		e := entity.ExpenseEntryListResponse{}

		e.ExpenseID = r.ExpenseID
		if r.DocNo != nil {
			e.DocumentNo = *r.DocNo
		}
		if r.Date != nil {
			e.Date = r.Date.Format("02/01/2006")
		}
		e.ExpenseTypeID = r.ExpenseTypeID
		if r.ExpenseTypeCode != nil {
			e.ExpenseTypeCode = *r.ExpenseTypeCode
		}
		if r.ExpenseTypeName != nil {
			e.ExpenseTypeName = *r.ExpenseTypeName
		}
		e.CollectorID = r.CollectorID
		if r.CollectorName != nil {
			e.CollectorName = *r.CollectorName
		}
		if r.Balance != nil {
			e.Balance = *r.Balance
		}
		if r.Amount != nil {
			e.Amount = *r.Amount
		}
		if r.Note != nil {
			e.Reason = *r.Note
		}
		if r.Source != nil {
			e.IsClockOut = *r.Source
		}

		resp = append(resp, e)
	}

	return resp, total, lastPage, nil
}

func normalizeExpenseDateEpoch(value string) (string, error) {
	trimmedValue := strings.TrimSpace(value)
	if trimmedValue == "" {
		return "", fmt.Errorf("invalid epoch")
	}

	if parsedDate, err := time.Parse("2006-01-02", trimmedValue); err == nil {
		return parsedDate.Format("2006-01-02"), nil
	}

	epoch, err := strconv.ParseInt(trimmedValue, 10, 64)
	if err != nil {
		return "", err
	}

	if epoch <= 0 {
		return "", fmt.Errorf("invalid epoch")
	}

	parsedTime := time.Unix(epoch, 0)
	if epoch >= 1_000_000_000_000 {
		parsedTime = time.UnixMilli(epoch)
	}

	return parsedTime.UTC().Format("2006-01-02"), nil
}

func (s *expenseEntryServiceImpl) Store(custId string, request entity.CreateExpenseEntryBody, userId int64) (entity.CreateExpenseResponseData, error) {
	ctx := context.Background()
	var result entity.CreateExpenseResponseData

	now := time.Now()
	m := model.Expense{
		CustID:        custId,
		ExpenseTypeID: request.ExpenseTypeID,
		Date:          &now,
		Amount:        &request.Amount,
		Note:          request.Note,
		CreatedBy:     userId,
		IsDel:         false,
		Balance:       &request.Amount,
	}
	if request.CollectorID != nil {
		cid := int64(*request.CollectorID)
		m.CollectorID = &cid
	}

	var docNo string
	if err := s.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		count, err := s.Repo.CountExpensesInCurrentMonth(txCtx)
		if err != nil {
			return err
		}
		year, month, day := now.Year(), now.Month(), now.Day()
		docNo = fmt.Sprintf("E%d%02d%02d%03d", year, month, day, count+1)
		m.DocNo = &docNo
		if err := s.Repo.Store(txCtx, &m); err != nil {
			return err
		}
		result.ExpenseID = m.ExpenseID
		result.DocumentNo = docNo

		if len(request.FileURL) > 0 {
			files := buildExpenseFilesFromURLs(custId, m.ExpenseID, request.FileURL, now)
			if err := s.Repo.StoreExpenseFiles(txCtx, custId, m.ExpenseID, files); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return result, errors.New("failed to save data, please try again")
	}
	return result, nil
}

// buildExpenseFilesFromURLs builds ExpenseFile rows from URLs.
// file_key (acf.media_category_type) is set to "image" or "video" from file extension; file_size is 0 when not provided.
func buildExpenseFilesFromURLs(custId string, expenseId int64, urls []string, createdAt time.Time) []model.ExpenseFile {
	const maxFileNameLen = 255
	const maxFileURLLen = 500

	files := make([]model.ExpenseFile, 0, len(urls))
	for i, rawURL := range urls {
		rawURL = strings.TrimSpace(rawURL)
		if rawURL == "" {
			continue
		}
		fileURL := rawURL
		if len(fileURL) > maxFileURLLen {
			fileURL = fileURL[:maxFileURLLen]
		}
		fileName := fileNameFromURL(rawURL, i)
		if len(fileName) > maxFileNameLen {
			fileName = fileName[:maxFileNameLen]
		}
		if fileName == "" {
			fileName = fmt.Sprintf("file_%d", i+1)
		}
		fileKey := mediaCategoryFromExtension(fileName)
		files = append(files, model.ExpenseFile{
			CustID:        custId,
			ExpenseID:     expenseId,
			FileName:      fileName,
			FileURL:       fileURL,
			FileKey:       fileKey,
			MediaCategory: nil,
			FileSize:      0,
			CreatedAt:     createdAt,
		})
	}
	return files
}

// mediaCategoryFromExtension returns "image" or "video" for acf.media_category_type based on file extension; defaults to "image".
func mediaCategoryFromExtension(fileName string) string {
	ext := strings.ToLower(strings.TrimPrefix(path.Ext(fileName), "."))
	switch ext {
	case "mp4", "webm", "mov", "avi", "mkv", "m4v", "flv", "wmv", "3gp":
		return "video"
	default:
		// jpg, jpeg, png, gif, webp, bmp, svg, etc. and unknown -> image
		return "image"
	}
}

func fileNameFromURL(rawURL string, index int) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Sprintf("file_%d", index+1)
	}
	name := path.Base(u.Path)
	name = strings.TrimSpace(name)
	if name == "" || name == "." {
		return fmt.Sprintf("file_%d", index+1)
	}
	return name
}

func ptrFloat64(v float64) *float64 { return &v }

func (s *expenseEntryServiceImpl) Update(custId string, expenseId int64, request entity.UpdateExpenseEntryBody, userId int64) (entity.UpdateExpenseResponseData, error) {
	ctx := context.Background()
	result := entity.UpdateExpenseResponseData{ExpenseID: expenseId}

	existing, err := s.Repo.FindById(ctx, custId, expenseId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return result, ErrExpenseNotFound
	}
	if err != nil {
		return result, errors.New("failed to update data, please try again")
	}

	updateData := map[string]interface{}{
		"updated_by": userId,
		"updated_at": time.Now(),
	}
	if request.Amount != nil {
		updateData["amount"] = *request.Amount
		result.Amount = *request.Amount
	} else if existing.Amount != nil {
		result.Amount = *existing.Amount
	}
	if request.Note != nil {
		updateData["note"] = *request.Note
		result.Note = request.Note
	} else {
		result.Note = existing.Note
	}

	if err := s.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Update(txCtx, custId, expenseId, updateData)
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, ErrExpenseNotFound
		}
		return result, errors.New("failed to update data, please try again")
	}

	return result, nil
}

func (s *expenseEntryServiceImpl) Delete(custId string, expenseId int64, deletedBy int64) error {
	ctx := context.Background()

	if err := s.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.Repo.Delete(txCtx, custId, expenseId, deletedBy)
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrExpenseNotFound
		}
		return errors.New("failed to delete data, please try again")
	}
	return nil
}

func (s *expenseEntryServiceImpl) Detail(custId string, expenseId int64) (entity.ExpenseDetailResponse, error) {
	ctx := context.Background()
	var resp entity.ExpenseDetailResponse
	row, err := s.Repo.FindDetailById(ctx, custId, expenseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, ErrExpenseNotFound
		}
		return resp, err
	}

	resp.ExpenseID = row.ExpenseID
	resp.ExpenseTypeID = row.ExpenseTypeID
	if row.Date != nil {
		resp.Date = row.Date.Format("2006-01-02")
	}
	if row.DocNo != nil {
		resp.DocNo = *row.DocNo
	}
	if row.ExpenseTypeCode != nil {
		resp.ExpenseTypeCode = *row.ExpenseTypeCode
	}
	if row.ExpenseTypeName != nil {
		resp.ExpenseTypeName = *row.ExpenseTypeName
	}
	resp.CollectorID = row.CollectorID
	if row.CollectorName != nil {
		resp.CollectorName = *row.CollectorName
	}
	if row.Amount != nil {
		resp.Amount = *row.Amount
	}
	if row.RemainingAmount != nil {
		resp.RemainingAmount = *row.RemainingAmount
	}
	if row.Note != nil {
		resp.Note = *row.Note
	}

	deposits, err := s.Repo.FindDepositExpensesByExpenseId(ctx, custId, expenseId)
	if err != nil {
		return resp, err
	}
	resp.Deposits = make([]entity.ExpenseDetailDeposit, 0, len(deposits))
	for _, de := range deposits {
		usedAmount, _ := de.PaymentAmount.Float64()
		updateDate := de.CreatedAt
		if de.UpdatedAt != nil {
			updateDate = *de.UpdatedAt
		}
		resp.Deposits = append(resp.Deposits, entity.ExpenseDetailDeposit{
			DepositExpenseID: de.DepositExpenseID,
			DepositID:        0, // not in schema; use 0 when absent
			UsedAmount:       usedAmount,
			DepositNo:        de.DepositNo,
			UpdateDate:       updateDate,
		})
	}

	files, err := s.Repo.FindExpenseFilesByExpenseId(ctx, custId, expenseId)
	if err != nil {
		return resp, err
	}
	resp.Files = make([]entity.ExpenseDetailFile, 0, len(files))
	for _, f := range files {
		ef := entity.ExpenseDetailFile{
			FileName: f.FileName,
			FileURL:  f.FileURL,
			FileSize: f.FileSize,
		}
		if f.MediaCategory != nil {
			ef.MediaCategory = *f.MediaCategory
		}
		if f.FileKey != "" {
			ef.FileType = f.FileKey
		}
		resp.Files = append(resp.Files, ef)
	}

	return resp, nil
}
