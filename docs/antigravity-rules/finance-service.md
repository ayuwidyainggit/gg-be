---
trigger: always_on
---

# Finance Service Rules

Rules untuk **Finance Service** (Port 9005) - Fiber v2 + sqlx.

## Service Overview

Finance Service menangani:
- Accounts Payable (AP)
- Accounts Receivable (AR)
- Journal Entries
- Bank Reconciliation
- Payment Schedule
- Aging Reports
- Chart of Accounts

---

## Project Structure

```
finance/
├── adapter/          # External adapters
├── controller/       # HTTP handlers
├── entity/           # Request/Response DTOs
├── model/            # Database models
├── pkg/              # Shared utilities
├── repository/       # Data access layer
├── service/          # Business logic
├── migration/        # Database migrations
└── main.go           # Entry point
```

---

## Controller Pattern

```go
type JournalController struct {
    JournalService service.JournalService
    validator      *validation.Validate
}

func (controller *JournalController) Route(app *fiber.App) {
    route := app.Group("/v1/journals", middleware.JWTProtected())
    route.Get("", controller.List)
    route.Get("/:journal_id", controller.Detail)
    route.Post("", controller.Create)
    route.Post("/:journal_id/reverse", controller.Reverse)
    route.Delete("/:journal_id", controller.Delete)
}
```

---

## Service Pattern

```go
type JournalService interface {
    CreateJournal(request entity.CreateJournalRequest) (entity.JournalResponse, error)
    GetJournal(journalId int64) (entity.JournalDetailResponse, error)
    ListJournals(filter entity.JournalFilter) ([]entity.JournalResponse, int, int, error)
    ReverseJournal(journalId int64, reason string) error
    DeleteJournal(journalId int64) error
}

type journalServiceImpl struct {
    JournalRepository repository.JournalRepository
    AccountRepository repository.AccountRepository
    TransactionRepo   repository.TransactionRepository
}
```

---

## Double Entry Accounting

```go
func (s *journalServiceImpl) CreateJournal(request entity.CreateJournalRequest) (entity.JournalResponse, error) {
    // 1. Validate debit = credit
    var totalDebit, totalCredit float64
    for _, entry := range request.Entries {
        totalDebit += entry.Debit
        totalCredit += entry.Credit
    }
    
    if totalDebit != totalCredit {
        return entity.JournalResponse{}, errors.New("journal not balanced: debit != credit")
    }
    
    // 2. Generate journal number
    journalNo := generateJournalNo(request.CustId, request.TransDate)
    
    // 3. Create journal header
    journal := model.Journal{
        JournalNo:   journalNo,
        TransDate:   request.TransDate,
        Description: request.Description,
        TotalAmount: totalDebit,
        Status:      "POSTED",
        CreatedBy:   request.CreatedBy,
    }
    
    // 4. Create journal entries
    err := s.JournalRepository.CreateWithEntries(journal, request.Entries)
    
    // 5. Update account balances
    for _, entry := range request.Entries {
        s.AccountRepository.UpdateBalance(entry.AccountId, entry.Debit, entry.Credit)
    }
    
    return s.mapToResponse(journal), err
}
```

---

## Journal Reversal Pattern

```go
func (s *journalServiceImpl) ReverseJournal(journalId int64, reason string) error {
    // 1. Get original journal
    original, err := s.JournalRepository.FindById(journalId)
    if err != nil {
        return err
    }
    
    if original.Status == "REVERSED" {
        return errors.New("journal already reversed")
    }
    
    // 2. Get original entries
    entries, _ := s.JournalRepository.GetEntries(journalId)
    
    // 3. Create reversal entries (swap debit/credit)
    var reversalEntries []model.JournalEntry
    for _, entry := range entries {
        reversalEntries = append(reversalEntries, model.JournalEntry{
            AccountId: entry.AccountId,
            Debit:     entry.Credit,  // Swap
            Credit:    entry.Debit,   // Swap
        })
    }
    
    // 4. Create reversal journal
    reversal := model.Journal{
        JournalNo:       generateJournalNo(original.CustId, time.Now()),
        TransDate:       time.Now(),
        Description:     "Reversal of " + original.JournalNo + ": " + reason,
        RefJournalId:    &journalId,
        Status:          "POSTED",
    }
    
    err = s.JournalRepository.CreateWithEntries(reversal, reversalEntries)
    
    // 5. Update original status
    original.Status = "REVERSED"
    s.JournalRepository.Update(original)
    
    // 6. Update account balances
    for _, entry := range reversalEntries {
        s.AccountRepository.UpdateBalance(entry.AccountId, entry.Debit, entry.Credit)
    }
    
    return err
}
```

---

## Aging Report Pattern

```go
type AgingService interface {
    GetAPAging(filter entity.AgingFilter) ([]entity.AgingResponse, error)
    GetARAging(filter entity.AgingFilter) ([]entity.AgingResponse, error)
}

func (s *agingServiceImpl) GetAPAging(filter entity.AgingFilter) ([]entity.AgingResponse, error) {
    // Group invoices by aging buckets: Current, 1-30, 31-60, 61-90, >90
    invoices, _ := s.InvoiceRepository.FindUnpaidAP(filter)
    
    agingMap := make(map[int64]*entity.AgingResponse) // keyed by supplier_id
    
    for _, inv := range invoices {
        daysOverdue := int(time.Since(inv.DueDate).Hours() / 24)
        
        aging, exists := agingMap[inv.SupplierId]
        if !exists {
            aging = &entity.AgingResponse{SupplierId: inv.SupplierId}
            agingMap[inv.SupplierId] = aging
        }
        
        switch {
        case daysOverdue <= 0:
            aging.Current += inv.Outstanding
        case daysOverdue <= 30:
            aging.Days1To30 += inv.Outstanding
        case daysOverdue <= 60:
            aging.Days31To60 += inv.Outstanding
        case daysOverdue <= 90:
            aging.Days61To90 += inv.Outstanding
        default:
            aging.DaysOver90 += inv.Outstanding
        }
        aging.Total += inv.Outstanding
    }
    
    return maps.Values(agingMap), nil
}
```

---

## Entity Conventions

```go
type CreateJournalRequest struct {
    TransDate   time.Time           `json:"trans_date" validate:"required"`
    Description string              `json:"description" validate:"required"`
    Entries     []JournalEntryInput `json:"entries" validate:"required,min=2,dive"`
    CustId      string              `json:"-"`
    CreatedBy   int64               `json:"-"`
}

type JournalEntryInput struct {
    AccountId int64   `json:"account_id" validate:"required"`
    Debit     float64 `json:"debit" validate:"gte=0"`
    Credit    float64 `json:"credit" validate:"gte=0"`
}

type AgingFilter struct {
    AsOfDate   time.Time `query:"as_of_date"`
    SupplierId int64     `query:"supplier_id"`
    CustomerId int64     `query:"customer_id"`
    CustId     string    `query:"-"`
}
```

---

## Testing Requirements

- Unit tests for double-entry validation
- Test journal reversal
- Test aging calculations
- Target coverage: >75%
