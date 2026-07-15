package controller

import (
	"cronjob/adapter"
	"cronjob/entity"
	"cronjob/pkg/config/env"
	"cronjob/pkg/constant"
	"cronjob/pkg/responsebuild"
	"cronjob/pkg/structs"
	"cronjob/pkg/validation"
	"cronjob/service"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

var (
	scheduler gocron.Scheduler
	jobMutex  sync.Mutex
	jobMap    = make(map[uuid.UUID]gocron.Job) // Track job objects by job ID
)

type JobController struct {
	Config     env.ConfigEnv
	JobService service.JobService
	validator  *validation.Validate
}

func NewJobController(
	config env.ConfigEnv,
	JobService service.JobService,
	validator *validation.Validate,
) *JobController {
	return &JobController{
		Config:     config,
		JobService: JobService,
		validator:  validator,
	}
}

func (controller *JobController) Route(app *fiber.App) {
	qParamId := ":job_id"
	routeV1 := app.Group("/v1/jobs")
	routeV1.Post("", controller.Create)
	routeV1.Get("/"+qParamId, controller.Detail)
	routeV1.Delete("/"+qParamId, controller.Delete)
	routeV1.Get("", controller.List)

	controller.Cronjob(app)
}

func (controller *JobController) Create(c *fiber.Ctx) error {
	var request entity.CreateJobBody
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.BodyParser(&request); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(request, headerAcceptLang)
	if errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	err := controller.JobService.Store(request)
	if err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *JobController) Detail(c *fiber.Ctx) error {
	var params entity.DetailJobBodyParam
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.ParamsParser(&params); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// custId := c.Locals("cust_id").(string)
	data, err := controller.JobService.Detail(params.JobID)
	if err != nil {
		log.Error(err.Error())
		statusCode := fiber.StatusBadRequest
		errMsg := err.Error()
		if err.Error() == "sql: no rows in result set" {
			statusCode = fiber.StatusNotFound
			errMsg = "Not found"
		}

		responsePayload.Setmsg(errMsg)
		return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *JobController) Delete(c *fiber.Ctx) error {
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	var params entity.DeleteJobBodyParams
	if err := c.ParamsParser(&params); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	errs := controller.validator.ValidateStruct(params, headerAcceptLang)
	if errs != nil {
		log.Error(errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// custId := c.Locals("cust_id").(string)
	// userId := c.Locals("user_id").(int64)
	// log.Println("VehicleController, Delete, CustId:", custId)

	err := controller.JobService.Delete(params.JobID)
	if err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setmsg("Deleted Successfully")
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *JobController) List(c *fiber.Ctx) error {
	var dataFilter entity.GeneralQueryFilter
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)
	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	// custId := c.Locals("cust_id").(string)
	data, total, lastPage, err := controller.JobService.List(dataFilter)
	if err != nil {
		log.Error(err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
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

/**
 * Cronjob is a function that sets up a scheduler for executing jobs periodically.
 *
 * Parameters:
 *   - app: A pointer to a fiber.App instance.
 *
 * Returns:
 *   - error: An error if the scheduler creation fails.
 */
func (controller *JobController) Cronjob(app *fiber.App) error {
	// Create scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Errorf("Failed to create scheduler: %+v", err)
		return err
	}

	// Load and schedule jobs dynamically
	controller.loadAndScheduleJobs(scheduler)

	intervalReloadJobs, err := strconv.Atoi(controller.Config.Get("INTERVAL_RELOAD_JOBS_IN_SECOND"))
	if err != nil {
		intervalReloadJobs = 30
	}
	// Periodically reload jobs
	ticker := time.NewTicker(time.Duration(intervalReloadJobs) * time.Second) // Refresh every x seconds
	go func() {
		for range ticker.C {
			controller.loadAndScheduleJobs(scheduler)
		}
	}()

	scheduler.Start()
	log.Info("Scheduler started.")
	return nil
}

/**
 * loadAndScheduleJobs loads and schedules jobs using the provided scheduler.
 *
 * Parameters:
 *     - scheduler: The gocron.Scheduler used to schedule the jobs.
 *
 * Returns:
 *     None
 */
func (controller *JobController) loadAndScheduleJobs(scheduler gocron.Scheduler) {
	jobMutex.Lock()
	defer jobMutex.Unlock()

	var dataFilter entity.GeneralQueryFilter
	dataFilter.Page = 1
	dataFilter.Limit = 100
	dataFilter.Active = 1
	dataFilter.Sort = "run_at:ASC,day_of_week_or_month:ASC,time_of_day:ASC"

	jobs, _, _, err := controller.JobService.List(dataFilter)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Create a set of active job IDs
	activeJobIDs := make(map[uuid.UUID]struct{})
	for _, job := range jobs {
		activeJobIDs[job.ID] = struct{}{}
	}

	// Remove inactive jobs
	for jobID, jobRef := range jobMap {
		if jobRef == nil {
			continue
		}
		printJobDataSingleLine(jobRef)
		nextRunTimeNil := time.Time{}
		nextRunTime, err := jobRef.NextRun()
		if err != nil {
			log.Error("nextRunTime:", nextRunTime)
			continue
		}
		if _, exists := activeJobIDs[jobID]; !exists || nextRunTime == nextRunTimeNil {
			scheduler.RemoveJob(jobRef.ID())
			delete(jobMap, jobID)
			log.Infof("Removed inactive job ID: %s", jobID)

			// Update inactive job in DB
			if err := controller.JobService.Inactive(jobID); err != nil {
				log.Errorf("Error inactive job %s: %+v", jobID, err)
			}
		}
	}

	// Schedule active jobs
	for _, job := range jobs {
		// log.Info("jobs >>>", structs.StructToJson(jobs))
		if _, exists := jobMap[job.ID]; exists {
			continue
		}

		var err error
		var jobRef gocron.Job

		switch job.JobType {
		case "one_time":
			if len(job.RunAt) == 0 {
				log.Warnf("One-time job ID %v has empty RunAt", job.ID)
				continue
			}
			log.Infof("job.RunAt: %s", job.RunAt)
			runAtTime, parseErr := time.Parse(time.RFC3339, job.RunAt)
			if parseErr != nil {
				log.Errorf("Error parsing one-time job time: %v", parseErr)
				continue
			}
			log.Infof("runAtTime: %v", runAtTime)
			if runAtTime.Before(time.Now().UTC()) {
				// Update inactive job in DB
				if err := controller.JobService.Inactive(job.ID); err != nil {
					log.Errorf("Error inactive job %s: %+v", job.ID, err)
				}
				continue
			}
			jobRef, err = scheduler.NewJob(
				gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(runAtTime)),
				gocron.NewTask(executeJob, job),
			)
		case "cron":
			jobRef, err = scheduler.NewJob(
				gocron.CronJob(
					job.CronExpression,
					false,
				),
				gocron.NewTask(executeJob, job),
			)
		case "daily":
			if len(job.TimeOfDay) == 0 {
				log.Warnf("Daily job ID %v has empty TimeOfDay", job.ID)
				continue
			}

			hours, err := strconv.Atoi(job.TimeOfDay[:2])
			if err != nil {
				log.Errorf("Error parsing daily job time hours: %v", err)
				continue
			}
			minutes, err := strconv.Atoi(job.TimeOfDay[3:5])
			if err != nil {
				log.Errorf("Error parsing daily job time minutes: %v", err)
				continue
			}
			seconds, err := strconv.Atoi(job.TimeOfDay[3:5])
			if err != nil {
				log.Errorf("Error parsing daily job time seconds: %v", err)
				continue
			}
			timeHours := uint(hours)
			timeMinutes := uint(minutes)
			timeSeconds := uint(seconds)
			jobRef, err = scheduler.NewJob(
				gocron.DailyJob(
					1, // Default to every day
					gocron.NewAtTimes(
						gocron.NewAtTime(timeHours, timeMinutes, timeSeconds),
					),
				),
				gocron.NewTask(executeJob, job),
			)
			if err != nil {
				break
			}
		case "weekly":
			if job.DayOfWeekOrMonth == 0 || len(job.TimeOfDay) == 0 {
				log.Warnf("Weekly job ID %v has empty DayOfWeek or TimeOfDay", job.ID)
				continue
			}

			hours, err := strconv.Atoi(job.TimeOfDay[:2])
			if err != nil {
				log.Errorf("Error parsing daily job time hours: %v", err)
				continue
			}
			minutes, err := strconv.Atoi(job.TimeOfDay[3:5])
			if err != nil {
				log.Errorf("Error parsing daily job time minutes: %v", err)
				continue
			}
			seconds, err := strconv.Atoi(job.TimeOfDay[3:5])
			if err != nil {
				log.Errorf("Error parsing daily job time seconds: %v", err)
				continue
			}
			timeHours := uint(hours)
			timeMinutes := uint(minutes)
			timeSeconds := uint(seconds)

			jobRef, err = scheduler.NewJob(
				gocron.WeeklyJob(
					1, // Default to every week
					gocron.NewWeekdays(time.Weekday(job.DayOfWeekOrMonth)),
					gocron.NewAtTimes(
						gocron.NewAtTime(timeHours, timeMinutes, timeSeconds),
					),
				),
				gocron.NewTask(executeJob, job),
			)

		default:
			log.Warnf("Job type %s not supported", job.JobType)
			continue
		}

		if err != nil {
			log.Errorf("Error scheduling job ID %v: %+v", job.ID, err)
		} else {
			jobMap[job.ID] = jobRef
			log.Debugf("Scheduled job ID %v: %s", job.ID, job.Task)
		}
	}
}

/**
 * executeJob executes a job based on the provided JobList.
 *
 * Parameters:
 *   - job: The JobList containing the job details.
 */
func executeJob(job entity.JobList) {
	log.Infof("Executing job ID %v: %s", job.ID, job.Task)
	switch job.Task {
	case "http_request":
		// Run HTTP request asynchronously using a goroutine
		go httpRequest(job)
	default:
		log.Warnf("Task '%s' not supported", job.Task)
	}

}

/**
 * httpRequest is a function that performs an HTTP request.
 *
 * Parameters:
 * - job: an entity.JobList object representing the job details.
 *
 * Returns: None.
 *
 * Description: This function sends an HTTP request to the specified URL using the POST method. It expects a non-empty URL in the job object. If the URL is empty, an error message is logged and the function returns. The function also expects a valid JSON payload in the job object, which is converted to a struct using the structs.JsonToStruct function. The HTTP request is dispatched using the adapter.HttpClientInfo object, and the response is checked for any errors. If the response status code is less than 200, an error message is logged. Finally, the response body is unmarshaled into a map[string]interface{} object.
 */
func httpRequest(job entity.JobList) {
	log.Infof("httpRequest %v: %s", job.ID, job.Task)
	if job.Url == "" {
		log.Error("HTTP job has empty URL")
		return
	}

	resp := map[string]interface{}{}
	pi, err := structs.JsonToStruct(job.Payload)
	if err != nil {
		log.Error(err)
	}
	client := adapter.HttpClientInfo{
		Url:     job.Url,
		Method:  "POST",
		Payload: pi,
	}

	res, err := client.Dispatch()
	if err != nil {
		log.Error("Dispatch, error: ", err)
	}

	if res.StatusCode() < 200 {
		log.Error("res.StatusCode:", res.Body())
	}

	err = json.Unmarshal(res.Body(), &resp)
	if err != nil {
		log.Error("Error unmarshal response:", err)
	}
}

/**
 * printJobDataSingleLine prints the job data in a single line.
 *
 * Parameters:
 *   - job: The gocron.Job object to print the data for.
 */
func printJobDataSingleLine(job gocron.Job) {
	var parts []string

	jobNextRun, _ := job.NextRun()
	jobLastRun, _ := job.LastRun()
	parts = append(parts, fmt.Sprintf("ID: %s", job.ID()))
	parts = append(parts, fmt.Sprintf("NextRun: %v", jobNextRun))
	parts = append(parts, fmt.Sprintf("LastRun: %v", jobLastRun))
	fmt.Println(strings.Join(parts, " | "))
}
