---
trigger: always_on
---

# TMS Service Rules

Rules untuk **TMS Service** (Port 9006) - Fiber v2 + GORM + Swagger.

## Service Overview

TMS (Transport Management System) Service menangani:
- Shipment management
- Delivery tracking
- Visit & Arrive management
- Unload & Reject handling
- Pickup management
- Picklist generation
- Driver & Vehicle assignment

---

## Project Structure

```
tms/
├── config/           # Configuration & database
├── controller/       # HTTP handlers
├── docs/             # Swagger documentation
├── entity/           # Request/Response DTOs
├── exception/        # Custom error handlers
├── helper/           # Utility functions
├── model/            # GORM models
├── repository/       # Data access layer
├── router/           # Route definitions
├── service/          # Business logic
├── utils/            # Additional utilities
└── main.go           # Entry point
```

---

## Key Differences from Other Services

| Aspect | TMS Service | Other Services |
|--------|-------------|----------------|
| ORM | GORM | sqlx |
| Error Handling | Custom Exception | Standard error |
| Router | Separate router package | Controller.Route() |
| Swagger | Swag annotations | Manual/None |

---

## Main Setup

```go
func main() {
    loadConfig, _ := config.LoadConfig(".")
    db := config.ConnectionDB(&loadConfig)
    validate := utils.InitializeValidator()

    // Swagger setup
    if loadConfig.Environment != "dev" {
        docs.SwaggerInfo.Host = loadConfig.SwaggerHost
        docs.SwaggerInfo.BasePath = loadConfig.SwaggerUrl
    }

    // Init layers
    shipmentRepo := repository.NewShipmentRepoImpl(db)
    shipmentService := service.NewShipmentServiceImpl(shipmentRepo, shipmentInvoicesRepo, validate)
    shipmentController := controller.NewShipmentController(shipmentService)

    routes := router.NewRouter(shipmentController, ...)

    app := fiber.New(fiber.Config{
        ErrorHandler: exception.ErrorHandler,
    })
    app.Use(recover.New())
    app.Use(requestid.New())
    app.Use(logger.New())
    app.Use(cors.New())
    app.Mount("/api/v1", routes)
    app.Get("/*", fiberSwagger.WrapHandler)
}
```

---

## Controller Pattern with Swagger

```go
// @Summary     Create shipment manually
// @Description Create a new shipment with manual vehicle assignment
// @Tags        Shipment
// @Accept      json
// @Produce     json
// @Param       request body entity.CreateShipmentRequest true "Shipment data"
// @Success     201 {object} entity.ShipmentResponse
// @Failure     400 {object} exception.ErrorResponse
// @Security    Bearer
// @Router      /shipments [post]
func (controller *ShipmentController) Create(ctx *fiber.Ctx) error {
    var request entity.CreateShipmentRequest
    
    if err := ctx.BodyParser(&request); err != nil {
        panic(exception.NewBadRequestError(err.Error()))
    }
    
    headers := map[string]string{
        "Authorization": ctx.Get("Authorization"),
        "X-Cust-Id":     ctx.Get("X-Cust-Id"),
    }
    
    result, err := controller.ShipmentService.CreateManual(ctx.Context(), headers, request)
    if err != nil {
        panic(exception.NewBadRequestError(err.Error()))
    }
    
    return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
        "code":    201,
        "status":  "CREATED",
        "data":    result,
    })
}
```

---

## Service Pattern with Context

```go
type ShipmentService interface {
    CreateManual(ctx context.Context, headers map[string]string, request entity.CreateShipmentRequest) (string, error)
    CreateAuto(ctx context.Context, headers map[string]string, request entity.CreateShipmentAutoRequest) ([]string, error)
    FindAll(ctx context.Context, dataFilter entity.ShipmentQueryFilter) []entity.ShipmentResponse
    Delete(ctx context.Context, headers map[string]string, params entity.ShipmentParams)
}

type ShipmentServiceImpl struct {
    shipmentRepo         repository.ShipmentRepo
    shipmentInvoicesRepo repository.ShipmentInvoicesRepo
    validate             *validator.Validate
}

func NewShipmentServiceImpl(shipmentRepo repository.ShipmentRepo, shipmentInvoicesRepo repository.ShipmentInvoicesRepo, validate *validator.Validate) ShipmentService {
    return &ShipmentServiceImpl{
        shipmentRepo:         shipmentRepo,
        shipmentInvoicesRepo: shipmentInvoicesRepo,
        validate:             validate,
    }
}
```

---

## Repository Pattern with GORM

```go
type ShipmentRepo interface {
    Insert(ctx context.Context, shipment model.Shipment) error
    InsertWithTx(tx *gorm.DB, shipment model.Shipment) error
    FindAll(ctx context.Context, filter entity.ShipmentQueryFilter) []model.Shipment
    FindByShipmentNo(ctx context.Context, shipmentNo string) (model.Shipment, error)
    DeleteByQuery(ctx context.Context, column, value string) error
    BeginTx(ctx context.Context) (*gorm.DB, error)
}

type ShipmentRepoImpl struct {
    db *gorm.DB
}

func (r *ShipmentRepoImpl) Insert(ctx context.Context, shipment model.Shipment) error {
    return r.db.WithContext(ctx).Create(&shipment).Error
}

func (r *ShipmentRepoImpl) FindAll(ctx context.Context, filter entity.ShipmentQueryFilter) []model.Shipment {
    var shipments []model.Shipment
    query := r.db.WithContext(ctx).Preload("ShipmentInvoices")
    
    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }
    
    query.Order("created_at DESC").Find(&shipments)
    return shipments
}
```

---

## Transaction Pattern

```go
func (s *ShipmentServiceImpl) CreateAuto(ctx context.Context, headers map[string]string, request entity.CreateShipmentAutoRequest) ([]string, error) {
    tx, err := s.shipmentRepo.BeginTx(ctx)
    if err != nil {
        panic(exception.NewInternalServerError(err.Error()))
    }

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        } else if err != nil {
            tx.Rollback()
        } else {
            err = tx.Commit().Error
        }
    }()

    // Use InsertWithTx for all operations
    err = s.shipmentRepo.InsertWithTx(tx, shipment)
    err = s.shipmentInvoicesRepo.InsertWithTx(tx, invoice)
    
    return shipmentNumbers, nil
}
```

---

## Exception Handling

```go
// exception/error_handler.go
func ErrorHandler(ctx *fiber.Ctx, err error) error {
    if e, ok := err.(*BadRequestError); ok {
        return ctx.Status(400).JSON(ErrorResponse{
            Code:    400,
            Status:  "BAD REQUEST",
            Errors:  e.Error,
            TraceId: ctx.Locals("requestid").(string),
        })
    }
    
    if e, ok := err.(*NotFoundError); ok {
        return ctx.Status(404).JSON(ErrorResponse{
            Code:    404,
            Status:  "NOT FOUND",
            Errors:  e.Error,
            TraceId: ctx.Locals("requestid").(string),
        })
    }
    
    return ctx.Status(500).JSON(ErrorResponse{
        Code:   500,
        Status: "INTERNAL SERVER ERROR",
    })
}

// Usage in service
panic(exception.NewBadRequestError("invalid data"))
panic(exception.NewNotFoundError("shipment not found"))
```

---

## GORM Model

```go
type Shipment struct {
    ID           int64              `gorm:"primaryKey;autoIncrement" json:"id"`
    ShipmentNo   string             `gorm:"uniqueIndex" json:"shipment_no"`
    DriverID     int64              `json:"driver_id"`
    DriverName   string             `json:"driver_name"`
    VehicleID    int64              `json:"vehicle_id"`
    VehicleNo    string             `json:"vehicle_no"`
    DeliveryDate time.Time          `json:"delivery_date"`
    Status       string             `json:"status"`
    CustID       string             `json:"cust_id"`
    CreatedAt    time.Time          `json:"created_at"`
    UpdatedAt    time.Time          `json:"updated_at"`
    Invoices     []ShipmentInvoices `gorm:"foreignKey:ShipmentNo;references:ShipmentNo"`
}
```

---

## Testing Requirements

- Unit tests with mocked GORM
- Test transaction rollback scenarios
- Target coverage: >75%
