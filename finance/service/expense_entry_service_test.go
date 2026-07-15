package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/repository"
	"strconv"
	"testing"
	"time"
)

type expenseEntryRepositoryListMock struct {
	repository.ExpenseEntryRepository
	capturedFilter entity.ExpenseEntryQueryFilter
	findAllFn      func(ctx context.Context, filter entity.ExpenseEntryQueryFilter) ([]model.ExpenseList, int64, int, error)
}

func (m *expenseEntryRepositoryListMock) FindAll(ctx context.Context, filter entity.ExpenseEntryQueryFilter) ([]model.ExpenseList, int64, int, error) {
	m.capturedFilter = filter
	if m.findAllFn != nil {
		return m.findAllFn(ctx, filter)
	}
	return []model.ExpenseList{}, 0, 0, nil
}

func TestExpenseEntryService_List_NormalizesUnixSecondsDateAndMapsResponse(t *testing.T) {
	startEpoch := int64(1775779200)
	endEpoch := int64(1775865599)
	startEpochStr := strconv.FormatInt(startEpoch, 10)
	endEpochStr := strconv.FormatInt(endEpoch, 10)
	expenseDate := time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC)
	docNo := "EXP-001"
	expenseTypeCode := "TRAVEL"
	expenseTypeName := "Travel"
	collectorName := "collector-a"
	balance := 12500.0
	amount := 18000.0
	note := "Taxi"
	source := 1
	collectorID := int64(99)

	repoMock := &expenseEntryRepositoryListMock{
		findAllFn: func(ctx context.Context, filter entity.ExpenseEntryQueryFilter) ([]model.ExpenseList, int64, int, error) {
			return []model.ExpenseList{
				{
					ExpenseID:       10,
					DocNo:           &docNo,
					Date:            &expenseDate,
					ExpenseTypeID:   7,
					ExpenseTypeCode: &expenseTypeCode,
					ExpenseTypeName: &expenseTypeName,
					CollectorID:     &collectorID,
					CollectorName:   &collectorName,
					Balance:         &balance,
					Amount:          &amount,
					Note:            &note,
					Source:          &source,
				},
			}, 1, 1, nil
		},
	}

	svc := &expenseEntryServiceImpl{Repo: repoMock}

	result, total, lastPage, err := svc.List(entity.ExpenseEntryQueryFilter{
		CustID:    "C001",
		UserID:    77,
		StartDate: startEpochStr,
		EndDate:   endEpochStr,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repoMock.capturedFilter.StartDate != time.Unix(startEpoch, 0).UTC().Format("2006-01-02") {
		t.Fatalf("expected normalized start_date %s, got %s", time.Unix(startEpoch, 0).UTC().Format("2006-01-02"), repoMock.capturedFilter.StartDate)
	}
	if repoMock.capturedFilter.EndDate != time.Unix(endEpoch, 0).UTC().Format("2006-01-02") {
		t.Fatalf("expected normalized end_date %s, got %s", time.Unix(endEpoch, 0).UTC().Format("2006-01-02"), repoMock.capturedFilter.EndDate)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
	if lastPage != 1 {
		t.Fatalf("expected lastPage 1, got %d", lastPage)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result))
	}
	if result[0].DocumentNo != docNo {
		t.Fatalf("expected document_no %s, got %s", docNo, result[0].DocumentNo)
	}
	if result[0].Date != "10/04/2026" {
		t.Fatalf("expected formatted date 10/04/2026, got %s", result[0].Date)
	}
	if result[0].CollectorName != collectorName {
		t.Fatalf("expected collector_name %s, got %s", collectorName, result[0].CollectorName)
	}
}

func TestExpenseEntryService_List_NormalizesUnixMillisecondsDate(t *testing.T) {
	startEpoch := time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC).UnixMilli()
	endEpoch := time.Date(2026, 4, 10, 23, 59, 59, 0, time.UTC).UnixMilli()

	repoMock := &expenseEntryRepositoryListMock{}
	svc := &expenseEntryServiceImpl{Repo: repoMock}

	_, _, _, err := svc.List(entity.ExpenseEntryQueryFilter{
		CustID:    "C001",
		UserID:    77,
		StartDate: strconv.FormatInt(startEpoch, 10),
		EndDate:   strconv.FormatInt(endEpoch, 10),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repoMock.capturedFilter.StartDate != "2026-04-10" {
		t.Fatalf("expected normalized start_date 2026-04-10, got %s", repoMock.capturedFilter.StartDate)
	}
	if repoMock.capturedFilter.EndDate != "2026-04-10" {
		t.Fatalf("expected normalized end_date 2026-04-10, got %s", repoMock.capturedFilter.EndDate)
	}
}

func TestExpenseEntryService_List_AppliesDefaultDateRangeWhenMissing(t *testing.T) {
	repoMock := &expenseEntryRepositoryListMock{}
	svc := &expenseEntryServiceImpl{Repo: repoMock}

	_, _, _, err := svc.List(entity.ExpenseEntryQueryFilter{
		CustID: "C001",
		UserID: 77,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repoMock.capturedFilter.StartDate == "" {
		t.Fatal("expected default start_date to be set")
	}
	if repoMock.capturedFilter.EndDate == "" {
		t.Fatal("expected default end_date to be set")
	}

	startDate, err := time.Parse("2006-01-02", repoMock.capturedFilter.StartDate)
	if err != nil {
		t.Fatalf("expected parsable start_date, got error: %v", err)
	}
	endDate, err := time.Parse("2006-01-02", repoMock.capturedFilter.EndDate)
	if err != nil {
		t.Fatalf("expected parsable end_date, got error: %v", err)
	}

	if !endDate.After(startDate) {
		t.Fatalf("expected end_date after start_date, got start=%s end=%s", repoMock.capturedFilter.StartDate, repoMock.capturedFilter.EndDate)
	}

	days := int(endDate.Sub(startDate).Hours() / 24)
	if days < 80 || days > 100 {
		t.Fatalf("expected roughly last 3 months window, got %d days", days)
	}
}

func TestExpenseEntryService_List_ReturnsErrorForInvalidEpochDate(t *testing.T) {
	repoMock := &expenseEntryRepositoryListMock{}
	svc := &expenseEntryServiceImpl{Repo: repoMock}

	_, _, _, err := svc.List(entity.ExpenseEntryQueryFilter{
		CustID:    "C001",
		UserID:    77,
		StartDate: "invalid-epoch",
		EndDate:   "1736899200000",
	})
	if err == nil {
		t.Fatal("expected invalid epoch date error")
	}
}
