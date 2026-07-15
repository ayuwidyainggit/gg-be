---
trigger: always_on
---

# Master Service Rules

Rules untuk **Master Service** (Port 9002) - Fiber v2 + sqlx.

## Service Overview

Master Service adalah service terbesar yang menangani semua master data:
- Products, Principals, Brands, Categories
- Outlets, Outlet Types, Outlet Groups
- Employees, Salesmen, Warehouses
- Pricing (Distributor Price, Special Price, TPR)
- Geographic (Province, Regency, District, Ward)
- Dan 80+ module master data lainnya

---

## Project Structure

```
master/
├── adapter/          # External adapters (OBS)
├── controller/       # HTTP handlers (80+ controllers)
├── entity/           # Request/Response DTOs (90+ entities)
├── model/            # Database models (90+ models)
├── pkg/              # Shared utilities
│   ├── config/       # Configuration & database
│   ├── middleware/   # Fiber middleware
│   ├── validation/   # Request validation
│   ├── structs/      # Automapper utility
│   └── server/       # Server setup
├── repository/       # Data access layer (80+ repos)
├── service/          # Business logic (80+ services)
└── main.go           # Entry point with DI setup
```

---

## Dependency Injection Pattern

```go
func main() {
    // 1. Load config
    envCfg := env.NewCfgEnv()
    validatorPkg := validation.NewValidator()
    postgreDB, _ := config.ConnToDb(envCfg)

    // 2. Setup Repository
    productRepository := repository.NewProductRepository(postgreDB)
    outletRepository := repository.NewOutletRepository(postgreDB)
    // ... more repositories

    // 3. Setup Service
    productService := service.NewProductService(productRepository)
    outletService := service.NewOutletService(outletRepository)
    // ... more services

    // 4. Setup Controller
    productController := controller.NewProductController(productService, validatorPkg)
    outletController := controller.NewOutletController(outletService, validatorPkg)
    // ... more controllers

    // 5. Setup Fiber & Routes
    app := fiber.New(fiberCfg)
    middleware.AppMiddleware(app, envCfg)
    productController.Route(app)
    outletController.Route(app)
    // ... more routes
}
```

---

## Controller Conventions

### Route Registration

```go
func (controller *ProductController) Route(app *fiber.App) {
    qParamId := ":pro_id"
    
    // File operations group
    productsFileRouteV1 := app.Group("/v1/products-file", middleware.JWTProtected())
    productsFileRouteV1.Get("/export", controller.Export)
    productsFileRouteV1.Get("/export-template", controller.ExportTemplate)
    productsFileRouteV1.Post("/import", controller.Import)

    // CRUD operations group
    productsRouteV1 := app.Group("/v1/products", middleware.JWTProtected())
    productsRouteV1.Get("/principals", controller.PrincipalList)
    productsRouteV1.Get("/categories", controller.CategoryList)
    productsRouteV1.Get("/"+qParamId, controller.Detail)
    productsRouteV1.Get("", controller.List)
    productsRouteV1.Post("", controller.Create)
    productsRouteV1.Post("/bulk", controller.Bulk)
    productsRouteV1.Patch("/"+qParamId, controller.Update)
    productsRouteV1.Delete("/"+qParamId, controller.Delete)
}
```

### List Handler with Pagination

```go
func (controller *ProductController) List(c *fiber.Ctx) error {
    var dataFilter entity.ProductQueryFilter
    responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

    if err := c.QueryParser(&dataFilter); err != nil {
        responsePayload.Setmsg(err.Error())
        return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
    }

    // Inject auth context
    dataFilter.CustId = c.Locals("cust_id").(string)
    dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)

    // Handle different modes
    switch dataFilter.Mode {
    case "search":
        data, total, lastPage, err = controller.ProductService.SearchList(dataFilter, dataFilter.CustId)
    case "lookup":
        data, total, lastPage, err = controller.ProductService.LookupList(dataFilter, dataFilter.CustId)
    default:
        data, total, lastPage, err = controller.ProductService.List(dataFilter, dataFilter.CustId)
    }

    responsePayload.Setdata(data)
    responsePayload.Setpaging(entity.Pagination{
        TotalRecord: total,
        PageCurrent: dataFilter.Page,
        PageLimit:   dataFilter.Limit,
        PageTotal:   lastPage,
    })
    return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}
```

---

## Service Conventions

### Interface Definition

```go
type ProductService interface {
    Detail(entity.DetailProductParams) (entity.ProductDetailResponse, error)
    List(entity.ProductQueryFilter, string) ([]entity.ProductResponse, int, int, error)
    Export(filter entity.ProductQueryFilter) (*bytes.Buffer, string, string, error)
    Store(entity.CreateProductBody) (entity.ProductResponse, error)
    Update(int64, entity.UpdateProductRequest) error
    Delete(string, int64, int64) error
}
```

### Response Mapping with Automapper

```go
func (service *productServiceImpl) List(dataFilter entity.ProductQueryFilter, custId string) ([]entity.ProductResponse, int, int, error) {
    products, total, lastPage, err := service.ProductRepository.FindAllByCustId(dataFilter, custId)
    if err != nil {
        return nil, 0, 0, err
    }

    var data []entity.ProductResponse
    for _, row := range products {
        var vResp entity.ProductResponse
        err = structs.Automapper(row, &vResp)
        if err != nil {
            return nil, 0, 0, err
        }
        data = append(data, vResp)
    }
    return data, total, lastPage, nil
}
```

---

## Repository Query Patterns

### Dynamic Query Builder

```go
func (r *productRepositoryImpl) FindAllByCustId(dataFilter entity.ProductQueryFilter, custId string) ([]model.Product, int, int, error) {
    selectCount := ` COUNT(*) AS total `
    selectField := ` p.cust_id, p.pro_id, p.pro_code, p.pro_name, ... `
    
    qWhere := ` LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id
                WHERE p.is_del = false AND p.cust_id = '` + custId + `' `

    // Dynamic filters
    if dataFilter.Query != "" {
        qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
                    OR p.pro_name ILIKE '%` + dataFilter.Query + `%')`
    }

    if len(dataFilter.PrincipalID) > 0 {
        intArrStr := str.ArrayToString(dataFilter.PrincipalID, ",")
        qWhere += ` AND p.principal_id IN (` + intArrStr + `) `
    }

    // Sorting
    if dataFilter.Sort != "" {
        mSortBy := strings.Split(dataFilter.Sort, ",")
        for _, row := range mSortBy {
            colSort := strings.Split(row, ":")
            sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
        }
    }

    // Pagination
    offset := (page - 1) * dataFilter.Limit
    lastPage := int(math.Ceil(float64(total) / float64(dataFilter.Limit)))
    querySelect += fmt.Sprintf(` LIMIT %d OFFSET %d`, dataFilter.Limit, offset)

    err = r.Select(&products, querySelect)
    return products, total, lastPage, err
}
```

---

## Import/Export Pattern

### Export to Excel/CSV

```go
func (service *productServiceImpl) Export(filter entity.ProductQueryFilter) (*bytes.Buffer, string, string, error) {
    products, _, err := service.ProductRepository.FindAllExport(filter, filter.CustId)
    
    switch filter.Format {
    case "csv":
        return service.createCSV(products), "text/csv", "products.csv", nil
    case "xlsx":
        return service.createXLSX(products), "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "products.xlsx", nil
    }
}
```

---

## Entity Conventions

### Query Filter

```go
type ProductQueryFilter struct {
    Query         string  `query:"q"`
    Page          int     `query:"page"`
    Limit         int     `query:"limit"`
    Sort          string  `query:"sort"`
    Format        string  `query:"format"`
    Mode          string  `query:"mode"`
    CustId        string  `query:"-"`
    ParentCustId  string  `query:"-"`
    PrincipalID   []int   `query:"principal_id"`
    BrandID       []int   `query:"brand_id"`
    IsActive      *int    `query:"is_active"`
}
```

---

## Testing Requirements

- Unit tests with mocked interfaces
- Target coverage: >75%
- Test import/export functionality
