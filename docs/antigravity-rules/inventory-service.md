---
trigger: always_on
---

# Inventory Service Rules

Rules untuk **Inventory Service** (Port 9003) - Fiber v2 + sqlx.

## Service Overview

Inventory Service menangani:
- Stock management
- Goods Receipt (GR)
- Stock Opname
- Stock Transfer
- Stock Adjustment
- Batch & Expiry tracking

---

## Project Structure

```
inventory/
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
type StockController struct {
    StockService service.StockService
    validator    *validation.Validate
}

func NewStockController(stockService service.StockService, validator *validation.Validate) *StockController {
    return &StockController{
        StockService: stockService,
        validator:    validator,
    }
}

func (controller *StockController) Route(app *fiber.App) {
    route := app.Group("/v1/stocks", middleware.JWTProtected())
    route.Get("", controller.List)
    route.Get("/:id", controller.Detail)
    route.Post("/adjustment", controller.Adjustment)
    route.Post("/opname", controller.Opname)
    route.Post("/transfer", controller.Transfer)
}
```

---

## Service Pattern

```go
type StockService interface {
    GetStock(params entity.StockParams) (entity.StockResponse, error)
    ListStock(filter entity.StockFilter) ([]entity.StockResponse, int, int, error)
    AdjustStock(request entity.AdjustmentRequest) error
    TransferStock(request entity.TransferRequest) error
    StockOpname(request entity.OpnameRequest) error
}

type stockServiceImpl struct {
    StockRepository    repository.StockRepository
    TransactionRepo    repository.TransactionRepository
}
```

---

## Repository Pattern

```go
type StockRepository interface {
    FindByProductAndWarehouse(productId, warehouseId int64) (model.Stock, error)
    UpdateStock(stock model.Stock) error
    CreateStockMovement(movement model.StockMovement) error
}

func (r *stockRepositoryImpl) UpdateStock(stock model.Stock) error {
    query := `UPDATE inv.m_stock SET 
              qty = $1, updated_at = NOW(), updated_by = $2
              WHERE stock_id = $3`
    _, err := r.Exec(query, stock.Qty, stock.UpdatedBy, stock.StockId)
    return err
}
```

---

## Stock Movement Pattern

```go
// Always record stock movement for audit trail
func (s *stockServiceImpl) AdjustStock(request entity.AdjustmentRequest) error {
    // 1. Get current stock
    stock, err := s.StockRepository.FindByProductAndWarehouse(request.ProductId, request.WarehouseId)
    
    // 2. Calculate new quantity
    oldQty := stock.Qty
    newQty := oldQty + request.AdjustmentQty
    
    // 3. Update stock
    stock.Qty = newQty
    err = s.StockRepository.UpdateStock(stock)
    
    // 4. Create movement record
    movement := model.StockMovement{
        StockId:     stock.StockId,
        MovementType: "ADJUSTMENT",
        QtyBefore:   oldQty,
        QtyAfter:    newQty,
        QtyChange:   request.AdjustmentQty,
        Reason:      request.Reason,
        CreatedBy:   request.CreatedBy,
    }
    return s.StockRepository.CreateStockMovement(movement)
}
```

---

## Batch & Expiry Handling

```go
type StockBatch struct {
    StockId     int64      `db:"stock_id"`
    ProductId   int64      `db:"pro_id"`
    WarehouseId int64      `db:"wh_id"`
    BatchNo     string     `db:"batch_no"`
    ExpDate     *time.Time `db:"exp_date"`
    Qty         float64    `db:"qty"`
}

// FIFO/FEFO logic for stock deduction
func (r *stockRepositoryImpl) DeductStockFEFO(productId, warehouseId int64, qty float64) error {
    // Get batches ordered by expiry date (earliest first)
    batches, _ := r.FindBatchesByProductWarehouse(productId, warehouseId)
    
    remainingQty := qty
    for _, batch := range batches {
        if remainingQty <= 0 {
            break
        }
        deductQty := math.Min(batch.Qty, remainingQty)
        // Deduct from this batch
        remainingQty -= deductQty
    }
    return nil
}
```

---

## Entity Conventions

```go
type StockFilter struct {
    WarehouseId int64  `query:"warehouse_id"`
    ProductId   int64  `query:"product_id"`
    BatchNo     string `query:"batch_no"`
    Query       string `query:"q"`
    Page        int    `query:"page"`
    Limit       int    `query:"limit"`
    CustId      string `query:"-"`
}

type AdjustmentRequest struct {
    ProductId     int64   `json:"product_id" validate:"required"`
    WarehouseId   int64   `json:"warehouse_id" validate:"required"`
    AdjustmentQty float64 `json:"adjustment_qty" validate:"required"`
    Reason        string  `json:"reason" validate:"required"`
    CustId        string  `json:"-"`
    CreatedBy     int64   `json:"-"`
}
```

---

## Error Handling

- Check stock availability before deduction
- Return clear error messages for insufficient stock
- Log all stock movements

---

## Testing Requirements

- Unit tests for stock calculations
- Test FIFO/FEFO logic
- Target coverage: >75%
