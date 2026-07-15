---
trigger: always_on
---

# Sales Service Rules

Rules untuk **Sales Service** (Port 9004) - Fiber v2 + sqlx.

## Service Overview

Sales Service menangani:
- Sales Order (SO)
- Sales Invoice
- Sales Return
- Payment & Collection
- POS Transaction
- Discount & Promotions

---

## Project Structure

```
sales/
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
type OrderController struct {
    OrderService service.OrderService
    validator    *validation.Validate
}

func (controller *OrderController) Route(app *fiber.App) {
    route := app.Group("/v1/orders", middleware.JWTProtected())
    route.Get("", controller.List)
    route.Get("/:order_no", controller.Detail)
    route.Post("", controller.Create)
    route.Patch("/:order_no/status", controller.UpdateStatus)
    route.Delete("/:order_no", controller.Cancel)
}
```

---

## Service Pattern

```go
type OrderService interface {
    CreateOrder(request entity.CreateOrderRequest) (entity.OrderResponse, error)
    GetOrder(orderNo string) (entity.OrderDetailResponse, error)
    ListOrders(filter entity.OrderFilter) ([]entity.OrderResponse, int, int, error)
    UpdateStatus(orderNo string, status int) error
    CancelOrder(orderNo string, reason string) error
}

type orderServiceImpl struct {
    OrderRepository   repository.OrderRepository
    ProductRepository repository.ProductRepository
    StockRepository   repository.StockRepository
    TransactionRepo   repository.TransactionRepository
}
```

---

## Order Creation Pattern

```go
func (s *orderServiceImpl) CreateOrder(request entity.CreateOrderRequest) (entity.OrderResponse, error) {
    // 1. Generate order number
    orderNo := generateOrderNo("SO", request.CustId)
    
    // 2. Validate products and calculate totals
    var totalBruto, totalDisc, totalNetto float64
    for _, item := range request.Items {
        product, err := s.ProductRepository.FindById(item.ProductId)
        if err != nil {
            return entity.OrderResponse{}, errors.New("product not found: " + strconv.Itoa(int(item.ProductId)))
        }
        
        // Calculate item total
        itemTotal := product.SellPrice1 * item.Qty
        totalBruto += itemTotal
    }
    
    // 3. Apply discounts
    totalDisc = s.calculateDiscounts(request, totalBruto)
    totalNetto = totalBruto - totalDisc
    
    // 4. Create order header
    order := model.SalesOrder{
        OrderNo:    orderNo,
        OutletId:   request.OutletId,
        TotalBruto: totalBruto,
        TotalDisc:  totalDisc,
        TotalNetto: totalNetto,
        Status:     1, // Draft
        CreatedBy:  request.CreatedBy,
    }
    
    // 5. Create order items
    err := s.OrderRepository.CreateWithItems(order, request.Items)
    return s.mapToResponse(order), err
}
```

---

## POS Transaction Pattern

```go
type POSService interface {
    CreateTransaction(request entity.POSRequest) (entity.POSResponse, error)
    GetOpenShift(cashierId int64) (*entity.ShiftResponse, error)
    OpenShift(request entity.OpenShiftRequest) error
    CloseShift(request entity.CloseShiftRequest) error
}

func (s *posServiceImpl) CreateTransaction(request entity.POSRequest) (entity.POSResponse, error) {
    // 1. Verify shift is open
    shift, err := s.GetOpenShift(request.CashierId)
    if err != nil {
        return entity.POSResponse{}, errors.New("no open shift found")
    }
    
    // 2. Create transaction
    txNo := generateTxNo("TX", request.CustId)
    
    // 3. Process payment
    for _, payment := range request.Payments {
        // Handle different payment methods (cash, card, ewallet)
    }
    
    // 4. Deduct stock immediately for POS
    for _, item := range request.Items {
        err := s.StockRepository.DeductStock(item.ProductId, request.WarehouseId, item.Qty)
        if err != nil {
            return entity.POSResponse{}, err
        }
    }
    
    return entity.POSResponse{TransactionNo: txNo}, nil
}
```

---

## Payment Pattern

```go
type PaymentService interface {
    CreatePayment(request entity.PaymentRequest) (entity.PaymentResponse, error)
    GetOutstandingInvoices(outletId int64) ([]entity.InvoiceResponse, error)
    ApplyPayment(paymentId int64, invoices []entity.InvoicePayment) error
}

func (s *paymentServiceImpl) ApplyPayment(paymentId int64, invoices []entity.InvoicePayment) error {
    payment, _ := s.PaymentRepository.FindById(paymentId)
    remainingAmount := payment.Amount
    
    for _, inv := range invoices {
        if remainingAmount <= 0 {
            break
        }
        
        invoice, _ := s.InvoiceRepository.FindById(inv.InvoiceId)
        applyAmount := math.Min(remainingAmount, invoice.Outstanding)
        
        // Update invoice
        invoice.PaidAmount += applyAmount
        invoice.Outstanding -= applyAmount
        if invoice.Outstanding == 0 {
            invoice.Status = "PAID"
        }
        
        s.InvoiceRepository.Update(invoice)
        remainingAmount -= applyAmount
    }
    return nil
}
```

---

## Entity Conventions

```go
type CreateOrderRequest struct {
    OutletId    int64              `json:"outlet_id" validate:"required"`
    SalesmanId  int64              `json:"salesman_id" validate:"required"`
    Items       []OrderItemRequest `json:"items" validate:"required,dive"`
    DiscountIds []int64            `json:"discount_ids"`
    Notes       string             `json:"notes"`
    CustId      string             `json:"-"`
    CreatedBy   int64              `json:"-"`
}

type OrderItemRequest struct {
    ProductId int64   `json:"product_id" validate:"required"`
    Qty       float64 `json:"qty" validate:"required,gt=0"`
    UnitId    string  `json:"unit_id" validate:"required"`
    Price     float64 `json:"price"`
}
```

---

## Order Status Flow

```
1 = Draft
2 = Submitted  
3 = Approved
4 = Delivered
5 = Completed
6 = Cancelled
```

---

## Testing Requirements

- Unit tests for order calculations
- Test discount logic
- Test payment allocation
- Target coverage: >75%
