---
trigger: always_on
---

# PJP Service Rules

Rules untuk **PJP Service** (Port 9009) - **Gin + GORM**.

## Service Overview

PJP (Permanent Journey Plan) Service menangani:
- PJP creation & management
- Route planning
- Route outlet assignment
- Daily route mapping
- Third-party integrations
- Visit tracking

---

## Project Structure

```
pjp/
├── config/               # Configuration & database
├── constant/             # Constants
├── controller/           # HTTP handlers
│   ├── pjp/              # PJP controllers
│   ├── pjp_auto/         # Auto PJP controllers
│   ├── pjp_enhance/      # Enhanced PJP controllers
│   ├── route/            # Route controllers
│   ├── visit/            # Visit controllers
│   └── third_party/      # External API controllers
├── data/                 # Request/Response structs
│   ├── request/          # Request DTOs
│   └── response/         # Response DTOs
├── database/             # Database connection
├── docs/                 # Swagger documentation
├── exception/            # Custom exceptions
├── helper/               # Utility functions
├── middleware/           # Gin middleware
├── model/                # GORM models
├── repository/           # Data access layer
├── router/               # Route definitions
├── service/              # Business logic
└── main.go               # Entry point
```

---

## Key Differences from Fiber Services

| Aspect | PJP (Gin) | Other Services (Fiber) |
|--------|-----------|------------------------|
| Framework | Gin | Fiber v2 |
| ORM | GORM | sqlx |
| Context | *gin.Context | *fiber.Ctx |
| Router | gin.Engine | fiber.App |
| Response | c.JSON() | c.JSON() |
| Params | c.Param() | c.Params() |
| Query | c.Query() | c.Query() |
| Body | c.ShouldBindJSON() | c.BodyParser() |

---

## Main Setup

```go
func main() {
    loadConfig, _ := config.LoadConfig(".")
    db := config.ConnectionDB(&loadConfig)
    validate := utils.InitializeValidator(db)

    // Swagger setup
    if loadConfig.Environment != "dev" {
        docs.SwaggerInfo.Host = loadConfig.SwaggerHost
        docs.SwaggerInfo.BasePath = "/scylla-pjp/api/v1"
    } else {
        docs.SwaggerInfo.Host = "localhost:9009"
        docs.SwaggerInfo.BasePath = "/api/v1"
    }

    // Init layers
    pjpRepository := pjpRepo.NewPjpRepository()
    pjpService := pjpServ.NewPjpService(pjpRepository, routeOutletRepository, routeRepository, validate, db)
    pjpController := pjpCtrl.NewPjpController(pjpService)

    routes := router.NewRouter(pjpController, ...)

    server := &http.Server{
        Addr:           ":" + loadConfig.ServerPort,
        Handler:        routes,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
    }
    server.ListenAndServe()
}
```

---

## Controller Pattern (Gin)

```go
type PjpController struct {
    pjpService pjp.PjpService
}

func NewPjpController(pjpService pjp.PjpService) *PjpController {
    return &PjpController{pjpService: pjpService}
}

// @Summary     Create PJP
// @Description Create a new PJP
// @Tags        PJP
// @Accept      json
// @Produce     json
// @Param       request body request.PjpRequest true "PJP data"
// @Success     201 {object} response.PjpResponse
// @Security    Bearer
// @Router      /pjp [post]
func (controller *PjpController) Create(ctx *gin.Context) {
    var req request.PjpRequest
    
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "code":    400,
            "status":  "BAD REQUEST",
            "message": err.Error(),
        })
        return
    }
    
    customerId := ctx.GetString("cust_id")
    controller.pjpService.Create(ctx, req, customerId)
    
    ctx.JSON(http.StatusCreated, gin.H{
        "code":    201,
        "status":  "CREATED",
        "message": "PJP created successfully",
    })
}


func (controller *PjpController) GetAll(ctx *gin.Context) {
    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
    
    filters := map[string]interface{}{
        "status":      ctx.Query("status"),
        "salesman_id": ctx.Query("salesman_id"),
    }
    
    customerId := ctx.GetString("cust_id")
    data, meta, err := controller.pjpService.GetAll(ctx, limit, page, filters, customerId)
    
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{
        "code":   200,
        "status": "OK",
        "data":   data,
        "meta":   meta,
    })
}
```

---

## Router Pattern (Gin)

```go
func NewRouter(
    pjpController *pjpCtrl.PjpController,
    routeController *route.RouteController,
    visitController *visit.VisitController,
) *gin.Engine {
    router := gin.Default()
    
    // Middleware
    router.Use(middleware.CORSMiddleware())
    router.Use(middleware.AuthMiddleware())
    
    // Swagger
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // API v1
    v1 := router.Group("/api/v1")
    {
        pjp := v1.Group("/pjp")
        {
            pjp.GET("", pjpController.GetAll)
            pjp.GET("/:id", pjpController.GetById)
            pjp.POST("", pjpController.Create)
            pjp.PUT("/:id", pjpController.Update)
            pjp.DELETE("/:id", pjpController.Delete)
        }
        
        routes := v1.Group("/routes")
        {
            routes.GET("", routeController.GetAll)
            routes.PATCH("/:id/status", routeController.UpdateStatus)
        }
    }
    
    return router
}
```

---

## Service Pattern

```go
type PjpService interface {
    Create(ctx context.Context, request request.PjpRequest, customerId string)
    GetAll(ctx context.Context, limit, page int, filters map[string]interface{}, customerId string) ([]response.PjpResponse, response.Meta, error)
    GetById(ctx context.Context, pjpId int, customerId string) response.PjpResponse
    Update(ctx context.Context, request request.PjpRequest, customerId string)
    Delete(ctx context.Context, pjpId int, customerId string)
}

type pjpService struct {
    pjpRepository         pjp.PjpRepository
    routeOutletRepository routeoutlet.RouteOutletRepository
    routeRepository       route.RouteRepository
    validate              *validator.Validate
    db                    *gorm.DB
}

func NewPjpService(pjpRepo pjp.PjpRepository, ..., db *gorm.DB) PjpService {
    return &pjpService{
        pjpRepository:         pjpRepo,
        routeOutletRepository: routeOutletRepository,
        routeRepository:       routeRepository,
        validate:              validate,
        db:                    db,
    }
}
```

---

## Repository Pattern (GORM)

```go
type PjpRepository interface {
    FindAll(db *gorm.DB, filters map[string]interface{}) ([]model.Pjp, error)
    FindById(db *gorm.DB, id int) (model.Pjp, error)
    Create(db *gorm.DB, pjp model.Pjp) error
    Update(db *gorm.DB, pjp model.Pjp) error
    Delete(db *gorm.DB, id int) error
}

type pjpRepositoryImpl struct{}

func NewPjpRepository() PjpRepository {
    return &pjpRepositoryImpl{}
}

func (r *pjpRepositoryImpl) FindAll(db *gorm.DB, filters map[string]interface{}) ([]model.Pjp, error) {
    var pjps []model.Pjp
    query := db.Preload("Routes")
    
    if status, ok := filters["status"].(string); ok && status != "" {
        query = query.Where("status = ?", status)
    }
    
    err := query.Find(&pjps).Error
    return pjps, err
}
```

---

## Automapper Pattern

```go
func toPjpResponse(value model.Pjp) response.PjpResponse {
    statusBool, _ := strconv.ParseBool(value.Status)

    res := response.PjpResponse{
        PjpCode: helper.FormatPjpCode(value.PjpCode),
        Status:  statusBool,
    }
    helper.Automapper(value, &res)
    
    return res
}

func mapRequestToModel(req request.PjpRequest, customerId string) model.Pjp {
    return model.Pjp{
        ID:            req.ID,
        PjpCode:       req.PjpCode,
        SalesManID:    req.SalesManID,
        Status:        "false",
        PjpMode:       "manual",
        CustID:        customerId,
    }
}
```

---

## Request/Response DTOs

```go
// data/request/pjp_request.go
type PjpRequest struct {
    ID            int    `json:"id"`
    PjpCode       int    `json:"pjp_code"`
    SalesManID    int    `json:"salesman_id" validate:"required"`
    SalesmanName  string `json:"salesman_name"`
    SalesmanCode  string `json:"salesman_code"`
    WarehouseID   int    `json:"warehouse_id"`
    WarehouseName string `json:"warehouse_name"`
    OperationType string `json:"operation_type"`
}

// data/response/pjp_response.go
type PjpResponse struct {
    ID           int            `json:"id"`
    PjpCode      string         `json:"pjp_code"`
    Status       bool           `json:"status"`
    SalesmanName string         `json:"salesman_name"`
    Routes       []RoutesEntity `json:"routes,omitempty"`
}
```

---

## Testing Requirements

- Unit tests with mocked GORM
- Test Gin context handling
- Target coverage: >75%
