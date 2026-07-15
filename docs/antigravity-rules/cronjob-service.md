---
trigger: always_on
---

# Cronjob Service Rules

Rules untuk **Cronjob Service** (Port 9100) - Fiber v2 + sqlx.

## Service Overview

Cronjob Service menangani scheduled tasks:
- Scheduled data synchronization
- Automated report generation
- Data cleanup & archival
- Email/notification scheduling
- Periodic calculations

---

## Project Structure

```
cronjob/
├── adapter/          # External adapters
├── controller/       # HTTP handlers (for manual triggers)
├── entity/           # Request/Response DTOs
├── model/            # Database models
├── pkg/              # Shared utilities
├── repository/       # Data access layer
├── service/          # Business logic
└── main.go           # Entry point with scheduler
```

---

## Main Setup with Scheduler

```go
func main() {
    envCfg := env.NewCfgEnv()
    postgreDB := config.PostgreSQLConnection(envCfg)
    
    // Setup Repository
    jobRepository := repository.NewJobRepository(postgreDB)
    
    // Setup Service
    jobService := service.NewJobService(jobRepository)
    
    // Setup Controller (for manual triggers)
    jobController := controller.NewJobController(jobService)
    
    // Setup Fiber
    app := fiber.New(fiberCfg)
    middleware.AppMiddleware(app)
    jobController.Route(app)
    
    // Start scheduler in goroutine
    go startScheduler(jobService)
    
    // Start server
    server.FiberServerWithGracefulShutdown(app)
}

func startScheduler(jobService service.JobService) {
    // Daily jobs at midnight
    gocron.Every(1).Day().At("00:00").Do(jobService.RunDailyCleanup)
    
    // Hourly sync
    gocron.Every(1).Hour().Do(jobService.SyncData)
    
    // Weekly report
    gocron.Every(1).Week().Monday().At("06:00").Do(jobService.GenerateWeeklyReport)
    
    gocron.StartBlocking()
}
```

---

## Controller Pattern (Manual Triggers)

```go
type JobController struct {
    JobService service.JobService
    validator  *validation.Validate
}

func (controller *JobController) Route(app *fiber.App) {
    route := app.Group("/v1/jobs", middleware.JWTProtected())
    route.Post("/run/:job_name", controller.RunJob)
    route.Get("/status/:job_id", controller.GetJobStatus)
    route.Get("/history", controller.GetJobHistory)
}

func (controller *JobController) RunJob(c *fiber.Ctx) error {
    jobName := c.Params("job_name")
    responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), "")
    
    // Run job asynchronously
    go func() {
        switch jobName {
        case "daily_cleanup":
            controller.JobService.RunDailyCleanup()
        case "sync_data":
            controller.JobService.SyncData()
        case "weekly_report":
            controller.JobService.GenerateWeeklyReport()
        default:
            log.Error("Unknown job:", jobName)
        }
    }()
    
    responsePayload.Setmsg("Job scheduled: " + jobName)
    return c.Status(fiber.StatusAccepted).JSON(responsePayload.GetRespPayload())
}
```

---

## Service Pattern

```go
type JobService interface {
    RunDailyCleanup() error
    SyncData() error
    GenerateWeeklyReport() error
    GetJobStatus(jobId int64) (entity.JobStatus, error)
    GetJobHistory(filter entity.JobFilter) ([]entity.JobHistory, error)
}

type jobServiceImpl struct {
    JobRepository    repository.JobRepository
    ReportRepository repository.ReportRepository
}
```

---

## Job Execution Pattern

```go
func (s *jobServiceImpl) RunDailyCleanup() error {
    jobId := s.startJob("daily_cleanup")
    
    defer func() {
        if r := recover(); r != nil {
            s.failJob(jobId, fmt.Sprintf("panic: %v", r))
        }
    }()
    
    // Record start
    log.Info("Starting daily cleanup job")
    
    // Execute cleanup tasks
    err := s.cleanupOldLogs(30) // 30 days
    if err != nil {
        s.failJob(jobId, err.Error())
        return err
    }
    
    err = s.cleanupTempData()
    if err != nil {
        s.failJob(jobId, err.Error())
        return err
    }
    
    // Record success
    s.completeJob(jobId)
    log.Info("Daily cleanup completed successfully")
    return nil
}

func (s *jobServiceImpl) startJob(jobName string) int64 {
    job := model.JobExecution{
        JobName:   jobName,
        Status:    "RUNNING",
        StartedAt: time.Now(),
    }
    id, _ := s.JobRepository.CreateExecution(job)
    return id
}

func (s *jobServiceImpl) completeJob(jobId int64) {
    s.JobRepository.UpdateExecution(jobId, model.JobExecution{
        Status:      "COMPLETED",
        CompletedAt: timePtr(time.Now()),
    })
}

func (s *jobServiceImpl) failJob(jobId int64, errorMsg string) {
    s.JobRepository.UpdateExecution(jobId, model.JobExecution{
        Status:      "FAILED",
        ErrorMsg:    &errorMsg,
        CompletedAt: timePtr(time.Now()),
    })
}
```

---

## Batch Processing Pattern

```go
func (s *jobServiceImpl) SyncData() error {
    jobId := s.startJob("sync_data")
    
    // Process in batches to avoid memory issues
    batchSize := 1000
    offset := 0
    totalProcessed := 0
    
    for {
        records, err := s.DataRepository.FindBatch(batchSize, offset)
        if err != nil {
            s.failJob(jobId, err.Error())
            return err
        }
        
        if len(records) == 0 {
            break
        }
        
        // Process batch
        for _, record := range records {
            err := s.processRecord(record)
            if err != nil {
                log.Warn("Failed to process record:", record.ID, err)
                continue // Skip failed records, continue with others
            }
            totalProcessed++
        }
        
        offset += batchSize
        
        // Update progress
        s.updateJobProgress(jobId, totalProcessed)
    }
    
    s.completeJob(jobId)
    log.Infof("Sync completed: %d records processed", totalProcessed)
    return nil
}
```

---

## Report Generation Pattern

```go
func (s *jobServiceImpl) GenerateWeeklyReport() error {
    jobId := s.startJob("weekly_report")
    
    // Calculate date range
    now := time.Now()
    endDate := now.AddDate(0, 0, -int(now.Weekday()))
    startDate := endDate.AddDate(0, 0, -7)
    
    // Generate report data
    salesData, _ := s.ReportRepository.GetSalesSummary(startDate, endDate)
    inventoryData, _ := s.ReportRepository.GetInventorySummary(endDate)
    
    // Create report
    report := model.Report{
        ReportType: "WEEKLY_SUMMARY",
        StartDate:  startDate,
        EndDate:    endDate,
        Data:       jsonEncode(map[string]interface{}{
            "sales":     salesData,
            "inventory": inventoryData,
        }),
        GeneratedAt: time.Now(),
    }
    
    err := s.ReportRepository.Create(report)
    if err != nil {
        s.failJob(jobId, err.Error())
        return err
    }
    
    // Send notification
    s.sendReportNotification(report)
    
    s.completeJob(jobId)
    return nil
}
```

---

## Job History Model

```go
type JobExecution struct {
    ID          int64      `db:"job_id"`
    JobName     string     `db:"job_name"`
    Status      string     `db:"status"` // RUNNING, COMPLETED, FAILED
    StartedAt   time.Time  `db:"started_at"`
    CompletedAt *time.Time `db:"completed_at"`
    ErrorMsg    *string    `db:"error_msg"`
    Progress    int        `db:"progress"`
    CreatedAt   time.Time  `db:"created_at"`
}
```

---

## Error Handling

- Always wrap job execution in recover()
- Log all job start/complete/fail events
- Track job history for debugging
- Continue processing on individual record failures

---

## Testing Requirements

- Unit tests for job logic
- Test batch processing
- Test error recovery
- Target coverage: >75%
