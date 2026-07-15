# Quick Start Guide

Panduan cepat untuk memulai development di Scylla Backend.

## 🚀 Quick Setup (5 Menit)

### 1. Prerequisites Check

```bash
# Check Go version
go version  # Should be 1.23.5+

# Check PostgreSQL
psql --version  # Should be 12+

# Check Docker (optional)
docker --version
```

### 2. Clone & Setup

```bash
# Clone repository
git clone <repository-url>
cd scylla-be

# Setup environment untuk service (contoh: cronjob)
cd cronjob
cp .env.example .env  # Jika ada, atau buat manual
```

### 3. Configure Environment

Edit file `.env` di direktori service:

```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=your_password
DB_NAME=scylla_db
DB_DEBUG=true
DB_MAX_CONNECTIONS=100
DB_MAX_IDLE_CONNECTIONS=10
DB_MAX_LIFETIME_CONNECTIONS=3600
```

### 4. Setup Database

```bash
# Create database
createdb scylla_db

# Atau menggunakan psql
psql -U postgres -c "CREATE DATABASE scylla_db;"
```

### 5. Install Dependencies

```bash
# Di direktori service
go mod download
go mod tidy
```

### 6. Run Service

```bash
# Run service
go run main.go

# Service akan running di http://localhost:8080
```

### 7. Test Service

```bash
# Health check
curl http://localhost:8080/ping

# Expected response: "It works"
```

## 📝 Contoh: Menambahkan Endpoint Baru

### 1. Buat Entity (DTO)

File: `entity/job_entity.go`
```go
type CreateJobBody struct {
    JobName string `json:"job_name" validate:"required,max=100"`
    JobDesc string `json:"job_desc" validate:"required,max=255"`
}
```

### 2. Buat Model

File: `model/job_model.go`
```go
type Job struct {
    JobID   uuid.UUID `gorm:"column:job_id;primaryKey;type:uuid;default:gen_random_uuid()"`
    JobName string    `gorm:"column:job_name"`
    JobDesc string    `gorm:"column:job_desc"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (Job) TableName() string {
    return "public.jobs"
}
```

### 3. Buat Repository

File: `repository/job_repository.go`
```go
type JobRepository interface {
    Store(c context.Context, data *model.Job) error
}

func (repo *RepositoryJobImpl) Store(c context.Context, data *model.Job) error {
    return repo.model(c).Create(data).Error
}
```

### 4. Buat Service

File: `service/job_service.go`
```go
func (service *jobServiceImpl) Store(request entity.CreateJobBody) error {
    c := context.Background()
    var jobModel model.Job
    
    structs.Automapper(request, &jobModel)
    
    return service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
        return service.JobRepository.Store(txCtx, &jobModel)
    })
}
```

### 5. Buat Controller

File: `controller/job_controller.go`
```go
func (controller *JobController) Create(c *fiber.Ctx) error {
    var request entity.CreateJobBody
    responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), "")
    
    if err := c.BodyParser(&request); err != nil {
        responsePayload.Setmsg(err.Error())
        return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
    }
    
    errs := controller.validator.ValidateStruct(request, "")
    if errs != nil {
        responsePayload.Setmsg(fiber.ErrBadRequest.Message)
        responsePayload.Seterrors(errs)
        return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
    }
    
    err := controller.JobService.Store(request)
    if err != nil {
        responsePayload.Setmsg(err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
    }
    
    return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *JobController) Route(app *fiber.App) {
    routeV1 := app.Group("/v1/jobs")
    routeV1.Post("", controller.Create)
}
```

### 6. Register Route di main.go

```go
jobController.Route(app)
```

### 7. Test Endpoint

```bash
curl -X POST http://localhost:8080/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "job_name": "Test Job",
    "job_desc": "Test Description"
  }'
```

## 🔧 Common Commands

### Development

```bash
# Run service
go run main.go

# Run dengan hot reload (jika menggunakan air)
air

# Build binary
go build -o bin/service-name main.go

# Run binary
./bin/service-name
```

### Testing

```bash
# Run tests
go test ./...

# Run tests dengan coverage
go test -cover ./...

# Run specific test
go test -v ./service -run TestJobService_Store
```

### Database

```bash
# Connect to database
psql -U postgres -d scylla_db

# Run SQL file
psql -U postgres -d scylla_db -f migration/001_create_table.sql
```

### Docker

```bash
# Build image
docker build -t scylla-service-name ./service-name

# Run container
docker run -p 8080:8080 --env-file ./service-name/.env scylla-service-name

# Run dengan docker-compose
docker-compose up
```

## 🐛 Troubleshooting

### Port Already in Use

```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>

# Atau change port di .env
SERVER_PORT=8081
```

### Database Connection Error

```bash
# Check PostgreSQL is running
pg_isready

# Check connection
psql -U postgres -d scylla_db -c "SELECT 1;"

# Verify environment variables
echo $DB_HOST
echo $DB_PORT
```

### Module Not Found

```bash
# Download dependencies
go mod download

# Tidy modules
go mod tidy

# Verify
go mod verify
```

## 📚 Next Steps

1. Baca [Development Guidelines](./DEVELOPMENT.md) untuk best practices
2. Lihat [API Structure](./API_STRUCTURE.md) untuk API documentation
3. Check [Database Documentation](./DATABASE.md) untuk database schema

## 💡 Tips

- Gunakan `DB_DEBUG=true` untuk melihat SQL queries
- Gunakan file `*_test.http` untuk manual testing
- Check logs untuk debugging
- Gunakan Postman/Insomnia untuk API testing

---

**Happy Coding!** 🎉

