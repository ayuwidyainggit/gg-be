---
trigger: always_on
---

# Mobile Service Rules

Rules untuk **Mobile Service** (Port 9008) - Fiber v2 + sqlx.

## Service Overview

Mobile Service adalah backend khusus untuk mobile apps:
- Mobile authentication
- Sync data for offline mode
- Mobile-specific endpoints
- Push notification handling
- Photo/file uploads from devices

---

## Project Structure

```
mobile/
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
type SyncController struct {
    SyncService service.SyncService
    validator   *validation.Validate
}

func (controller *SyncController) Route(app *fiber.App) {
    route := app.Group("/v1/mobile", middleware.JWTProtected())
    route.Get("/sync/outlets", controller.SyncOutlets)
    route.Get("/sync/products", controller.SyncProducts)
    route.Get("/sync/prices", controller.SyncPrices)
    route.Post("/sync/orders", controller.SyncOrders)
    route.Post("/sync/visits", controller.SyncVisits)
}
```

---

## Mobile-Optimized Response

```go
// Minimal response for mobile bandwidth
type MobileOutletResponse struct {
    OutletId   int64   `json:"o_id"`
    OutletCode string  `json:"o_cd"`
    OutletName string  `json:"o_nm"`
    Address    string  `json:"addr"`
    Lat        float64 `json:"lat"`
    Lng        float64 `json:"lng"`
    LastSync   int64   `json:"ls"` // Unix timestamp
}

func (s *syncServiceImpl) GetOutletsForSync(filter entity.SyncFilter) ([]MobileOutletResponse, error) {
    outlets, _ := s.OutletRepository.FindModifiedSince(filter.LastSync)
    
    var result []MobileOutletResponse
    for _, o := range outlets {
        result = append(result, MobileOutletResponse{
            OutletId:   o.OutletId,
            OutletCode: o.OutletCode,
            OutletName: o.OutletName,
            Address:    o.Address,
            Lat:        o.Latitude,
            Lng:        o.Longitude,
            LastSync:   o.UpdatedAt.Unix(),
        })
    }
    return result, nil
}
```

---

## Delta Sync Pattern

```go
type SyncFilter struct {
    LastSync    int64  `query:"last_sync"` // Unix timestamp
    SalesmanId  int64  `query:"salesman_id"`
    WarehouseId int64  `query:"warehouse_id"`
    CustId      string `query:"-"`
}

func (r *outletRepositoryImpl) FindModifiedSince(lastSync int64) ([]model.Outlet, error) {
    var outlets []model.Outlet
    syncTime := time.Unix(lastSync, 0)
    
    query := `SELECT * FROM mst.m_outlet 
              WHERE updated_at > $1 AND is_del = false`
    err := r.Select(&outlets, query, syncTime)
    return outlets, err
}
```

---

## Offline Order Sync

```go
type SyncOrderRequest struct {
    Orders []MobileOrder `json:"orders"`
}

type MobileOrder struct {
    LocalId     string            `json:"local_id"` // UUID from device
    OutletId    int64             `json:"outlet_id"`
    Items       []MobileOrderItem `json:"items"`
    CreatedAtLocal int64          `json:"created_at"` // Device timestamp
}

func (s *syncServiceImpl) SyncOrders(request entity.SyncOrderRequest) ([]entity.SyncResult, error) {
    var results []entity.SyncResult
    
    for _, order := range request.Orders {
        // Check if already synced (idempotency)
        existing, _ := s.OrderRepository.FindByLocalId(order.LocalId)
        if existing != nil {
            results = append(results, entity.SyncResult{
                LocalId:  order.LocalId,
                ServerId: existing.OrderNo,
                Status:   "ALREADY_SYNCED",
            })
            continue
        }
        
        // Create new order
        orderNo, err := s.OrderService.CreateFromMobile(order)
        status := "SUCCESS"
        if err != nil {
            status = "FAILED"
        }
        
        results = append(results, entity.SyncResult{
            LocalId:  order.LocalId,
            ServerId: orderNo,
            Status:   status,
            Error:    err,
        })
    }
    
    return results, nil
}
```

---

## Photo Upload Pattern

```go
type PhotoController struct {
    PhotoService service.PhotoService
}

func (controller *PhotoController) Route(app *fiber.App) {
    route := app.Group("/v1/mobile/photos", middleware.JWTProtected())
    route.Post("/outlet/:outlet_id", controller.UploadOutletPhoto)
    route.Post("/visit/:visit_id", controller.UploadVisitPhoto)
}

func (controller *PhotoController) UploadOutletPhoto(c *fiber.Ctx) error {
    outletId, _ := strconv.ParseInt(c.Params("outlet_id"), 10, 64)
    
    file, err := c.FormFile("photo")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "no file uploaded"})
    }
    
    // Validate file size (max 5MB for mobile)
    if file.Size > 5*1024*1024 {
        return c.Status(400).JSON(fiber.Map{"error": "file too large"})
    }
    
    // Upload to OBS
    url, err := controller.PhotoService.UploadPhoto(file, "outlets", outletId)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    
    return c.JSON(fiber.Map{"url": url})
}
```

---

## GPS Location Tracking

```go
type VisitRequest struct {
    OutletId    int64   `json:"outlet_id" validate:"required"`
    Latitude    float64 `json:"latitude" validate:"required"`
    Longitude   float64 `json:"longitude" validate:"required"`
    Accuracy    float64 `json:"accuracy"`
    VisitType   string  `json:"visit_type"` // CHECK_IN, CHECK_OUT
    DeviceTime  int64   `json:"device_time"`
}

func (s *visitServiceImpl) RecordVisit(request entity.VisitRequest) error {
    // Validate distance from outlet
    outlet, _ := s.OutletRepository.FindById(request.OutletId)
    distance := haversineDistance(
        request.Latitude, request.Longitude,
        outlet.Latitude, outlet.Longitude,
    )
    
    if distance > 100 { // 100 meters threshold
        return errors.New("too far from outlet location")
    }
    
    visit := model.Visit{
        OutletId:   request.OutletId,
        Latitude:   request.Latitude,
        Longitude:  request.Longitude,
        VisitType:  request.VisitType,
        Distance:   distance,
        DeviceTime: time.Unix(request.DeviceTime, 0),
    }
    
    return s.VisitRepository.Create(visit)
}
```

---

## Entity Conventions

```go
// Use short field names for mobile bandwidth
type MobileProductResponse struct {
    Id   int64   `json:"id"`
    Cd   string  `json:"cd"`   // code
    Nm   string  `json:"nm"`   // name
    P1   float64 `json:"p1"`   // price1
    P2   float64 `json:"p2"`   // price2
    P3   float64 `json:"p3"`   // price3
    U1   string  `json:"u1"`   // unit1
    U2   string  `json:"u2"`   // unit2
    U3   string  `json:"u3"`   // unit3
    Img  string  `json:"img"`  // image_url
}
```

---

## Testing Requirements

- Unit tests for sync logic
- Test offline order handling
- Test GPS validation
- Target coverage: >75%
