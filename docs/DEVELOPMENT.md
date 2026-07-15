# Development Guidelines

Panduan lengkap untuk development di project Scylla Backend.

## 📋 Table of Contents

- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Code Review](#code-review)

## 🚀 Getting Started

### Prerequisites

Pastikan sudah install:
- Go 1.23.5+
- PostgreSQL 12+
- Redis 6+ (untuk service yang memerlukan)
- Git
- IDE dengan Go support

### Initial Setup

1. **Clone Repository**
   ```bash
   git clone <repository-url>
   cd scylla-be
   ```

2. **Setup Environment**
   ```bash
   # Copy env template
   cp service-name/.env.example service-name/.env
   
   # Edit .env sesuai kebutuhan
   nano service-name/.env
   ```

3. **Install Dependencies**
   ```bash
   cd service-name
   go mod download
   go mod tidy
   ```

4. **Run Service**
   ```bash
   go run main.go
   ```

## 📁 Project Structure

### Service Structure

Setiap service mengikuti struktur berikut:

```
service-name/
├── main.go              # Entry point
├── controller/          # HTTP layer
├── service/            # Business logic
├── repository/         # Data access
├── entity/             # DTOs
├── model/              # Database models
├── adapter/            # External integrations
└── pkg/                # Utilities
```

### Layer Responsibilities

#### Controller Layer
- Handle HTTP requests/responses
- Request validation
- Route definition
- Response formatting

#### Service Layer
- Business logic
- Transaction management
- Data transformation
- Business rules validation

#### Repository Layer
- Database operations
- Query building
- Data mapping
- Transaction context handling

#### Entity Layer
- Request DTOs
- Response DTOs
- Validation tags
- API contracts

#### Model Layer
- Database schema
- GORM tags
- Table relationships
- Database constraints

## 📝 Coding Standards

### Naming Conventions

#### Files
- Use `snake_case` untuk file names
- Contoh: `job_controller.go`, `user_service.go`

#### Packages
- Use `lowercase` untuk package names
- Contoh: `controller`, `service`, `repository`

#### Types
- Use `PascalCase` untuk exported types
- Contoh: `JobController`, `JobService`, `JobRepository`

#### Functions
- Use `PascalCase` untuk exported functions
- Use `camelCase` untuk unexported functions
- Contoh: `NewJobController()`, `Store()`, `findOneByID()`

#### Variables
- Use `camelCase` untuk variables
- Contoh: `jobService`, `validatorPkg`, `postgreDB`

### Code Organization

#### Import Order
```go
import (
    // Standard library
    "context"
    "fmt"
    "time"
    
    // Third-party packages
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
    
    // Local packages
    "service-name/entity"
    "service-name/service"
)
```

#### Function Structure
```go
// 1. Function signature
func (service *jobServiceImpl) Store(request entity.CreateJobBody) error {
    // 2. Context setup
    c := context.Background()
    
    // 3. Validation
    if err := validate(request); err != nil {
        return err
    }
    
    // 4. Business logic
    // ...
    
    // 5. Transaction
    return service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
        // 6. Repository call
        return service.JobRepository.Store(txCtx, &jobModel)
    })
}
```

### Error Handling

#### Always Return Errors
```go
// ✅ Good
func (repo *Repository) Find(id int) (*Model, error) {
    var model Model
    err := repo.DB.Where("id = ?", id).First(&model).Error
    if err != nil {
        return nil, err
    }
    return &model, nil
}

// ❌ Bad
func (repo *Repository) Find(id int) *Model {
    var model Model
    repo.DB.Where("id = ?", id).First(&model)
    return &model
}
```

#### Log Errors Appropriately
```go
// ✅ Good
if err != nil {
    log.Errorf("Error storing job: %+v", err)
    return err
}

// ❌ Bad
if err != nil {
    panic(err)
}
```

### Transaction Management

#### Always Use Transaction for Write Operations
```go
// ✅ Good
err := service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
    if err := service.JobRepository.Store(txCtx, &jobModel); err != nil {
        return err
    }
    // Additional operations...
    return nil
})

// ❌ Bad
err := service.JobRepository.Store(c, &jobModel)
```

#### Repository Context Handling
```go
func (repo *RepositoryJobImpl) model(ctx context.Context) *gorm.DB {
    tx := extractTx(ctx)
    if tx != nil {
        return tx.WithContext(ctx)
    }
    return repo.WithContext(ctx)
}
```

### Validation

#### Use Struct Tags for Validation
```go
type CreateJobBody struct {
    JobName string `json:"job_name" validate:"required,max=100"`
    JobDesc string `json:"job_desc" validate:"required,max=255"`
    JobType string `json:"job_type" validate:"required,oneof=duration cron daily"`
}
```

#### Validate in Controller
```go
errs := controller.validator.ValidateStruct(request, headerAcceptLang)
if errs != nil {
    responsePayload.Setmsg(fiber.ErrBadRequest.Message)
    responsePayload.Seterrors(errs)
    return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
}
```

### Response Formatting

#### Use Response Builder
```go
responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
responsePayload.Setdata(data)
responsePayload.Setpaging(entity.Pagination{
    TotalRecord: total,
    PageCurrent: page,
    PageLimit:   limit,
    PageTotal:   lastPage,
})
return c.JSON(responsePayload.GetRespPayload())
```

## 🔄 Development Workflow

### Branch Strategy

1. **Main Branch**: Production-ready code
2. **Feature Branches**: `feature/feature-name`
3. **Bugfix Branches**: `bugfix/bug-name`
4. **Hotfix Branches**: `hotfix/issue-name`

### Commit Messages

Format commit message:
```
type(scope): subject

body (optional)

footer (optional)
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

Contoh:
```
feat(cronjob): add job scheduler functionality

- Add gocron integration
- Add job CRUD endpoints
- Add job execution logic
```

### Creating New Feature

1. **Create Feature Branch**
   ```bash
   git checkout -b feature/job-scheduler
   ```

2. **Develop Feature**
   - Follow coding standards
   - Write clean code
   - Add comments where needed

3. **Test Locally**
   ```bash
   go run main.go
   # Test endpoints
   ```

4. **Commit Changes**
   ```bash
   git add .
   git commit -m "feat(cronjob): add job scheduler"
   ```

5. **Push and Create MR**
   ```bash
   git push origin feature/job-scheduler
   # Create merge request
   ```

### Adding New Service

1. **Create Service Structure**
   ```bash
   mkdir -p new-service/{controller,service,repository,entity,model,adapter,pkg}
   ```

2. **Initialize Go Module**
   ```bash
   cd new-service
   go mod init new-service
   ```

3. **Copy Base Structure**
   - Copy `main.go` template
   - Copy `pkg/` utilities
   - Setup environment config

4. **Implement Layers**
   - Create models
   - Create entities
   - Create repositories
   - Create services
   - Create controllers

5. **Register Routes**
   ```go
   // In main.go
   controller.Route(app)
   ```

## 🧪 Testing

### Unit Testing

```go
func TestJobService_Store(t *testing.T) {
    // Setup
    mockRepo := &MockJobRepository{}
    service := NewJobService(mockRepo, mockTransaction)
    
    // Test
    err := service.Store(entity.CreateJobBody{
        JobName: "Test Job",
        // ...
    })
    
    // Assert
    assert.NoError(t, err)
    assert.True(t, mockRepo.StoreCalled)
}
```

### Integration Testing

```go
func TestJobController_Create(t *testing.T) {
    // Setup test server
    app := fiber.New()
    controller.Route(app)
    
    // Make request
    req := httptest.NewRequest("POST", "/v1/jobs", bytes.NewBuffer(jsonData))
    resp, _ := app.Test(req)
    
    // Assert
    assert.Equal(t, 201, resp.StatusCode)
}
```

### Manual Testing

Gunakan file `*_test.http` untuk manual testing:

```http
### Create Job
POST http://localhost:8080/v1/jobs
Content-Type: application/json

{
  "job_name": "Test Job",
  "job_desc": "Test Description",
  "job_type": "daily",
  "task": "http_request",
  "url": "https://example.com",
  "created_by": "admin"
}
```

## 👀 Code Review

### Review Checklist

- [ ] Code follows naming conventions
- [ ] Error handling is proper
- [ ] Transactions are used for write operations
- [ ] Validation is implemented
- [ ] Response format is consistent
- [ ] Logging is appropriate
- [ ] No hardcoded values
- [ ] Comments are clear and necessary
- [ ] No unused imports/variables
- [ ] Code is readable and maintainable

### Common Issues to Avoid

1. **Panic in Service/Repository**
   ```go
   // ❌ Bad
   if err != nil {
       panic(err)
   }
   
   // ✅ Good
   if err != nil {
       log.Errorf("Error: %+v", err)
       return err
   }
   ```

2. **Missing Transaction**
   ```go
   // ❌ Bad
   service.Repo1.Store(data1)
   service.Repo2.Store(data2)
   
   // ✅ Good
   service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
       service.Repo1.Store(txCtx, data1)
       service.Repo2.Store(txCtx, data2)
       return nil
   })
   ```

3. **Hardcoded Values**
   ```go
   // ❌ Bad
   limit := 10
   
   // ✅ Good
   limit := dataFilter.Limit
   if limit == 0 {
       limit = 10 // default
   }
   ```

## 📚 Resources

- [Go Best Practices](https://golang.org/doc/effective_go)
- [Fiber Documentation](https://docs.gofiber.io)
- [GORM Documentation](https://gorm.io/docs)
- [Project Wiki](./WIKI.md)

---

**Last Updated**: 2024

