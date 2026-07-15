package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"scyllax-pjp/config"
	"scyllax-pjp/data/entity"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/exception"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"scyllax-pjp/repository"
	"scyllax-pjp/repository/pjp"

	"github.com/go-playground/validator/v10"
)

type VisitService interface {
	FinishVisit(ctx context.Context, request request.FinishVisitRequest)
	ArriveVisit(ctx context.Context, request request.ArriveVisitRequest)
	SkipVisit(ctx context.Context, request request.SkipVisitRequest)
	ResumeVisit(ctx context.Context, request request.ResumeVisitRequest)
	LeaveVisit(ctx context.Context, request request.LeaveVisitRequest)
	UnloadVisit(ctx context.Context, request request.OnholdVisitRequest)
	OutletVisit(ctx context.Context, request request.OutletVisitRequest)
	GetAllOutletBySalesCode(ctx context.Context, salesCode, custId, date, routeCode string) (responses []response.DestinationResponse)
	SummaryVisit(ctx context.Context, dataFilter entity.SummaryQueryFilter) (res response.SummaryResponse)
	TravelList(ctx context.Context, outletVisitId int) (response response.TodoListResponse)
	SummaryVisitStatus(ctx context.Context, dataFilter entity.SummaryQueryFilter) (res response.VisitStatusResponse)
	GetVisitOutletList(ctx context.Context, dataFilter entity.SummaryQueryFilter) (res []response.OutletVisitListResponse)
	GetSalesmanReport(ctx context.Context, dataFilter entity.SalesmanReportQueryFilter) entity.SalesmanReport
}

type VisitServiceImpl struct {
	RouteOutletRepo     repository.RouteOutletRepository
	OutletVisitListRepo repository.OutletVisitRepo
	PjpRepo             pjp.PjpRepository
	Validate            *validator.Validate
}

func NewVisitServiceImpl(RouteOutletRepo repository.RouteOutletRepository, OutletVisitListRepo repository.OutletVisitRepo, PjpRepo pjp.PjpRepository, validate *validator.Validate) VisitService {
	return &VisitServiceImpl{
		RouteOutletRepo:     RouteOutletRepo,
		OutletVisitListRepo: OutletVisitListRepo,
		PjpRepo:             PjpRepo,
		Validate:            validate,
	}
}

func (service *VisitServiceImpl) FinishVisit(ctx context.Context, request request.FinishVisitRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	service.OutletVisitListRepo.UpdateOutletVisitListMultiRowColumnAt(ctx, request.SalesmanCode, request.CustID, "finish", request.CurrentTime, request.Date)
}

func (service *VisitServiceImpl) ArriveVisit(ctx context.Context, request request.ArriveVisitRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	service.OutletVisitListRepo.UpdateOutletVisitListColumnAt(ctx, request.SalesmanCode, request.CustID, "arrive_at", request.CurrentTime, request.Date, request.Id)
}

func (service *VisitServiceImpl) SkipVisit(ctx context.Context, request request.SkipVisitRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	service.OutletVisitListRepo.UpdateOutletVisitListSkipColumnAt(ctx, request.SalesmanCode, request.CustID, "skip_at", request.CurrentTime, request.Date, request.Id, request.SkipReason)
}

func (service *VisitServiceImpl) ResumeVisit(ctx context.Context, request request.ResumeVisitRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	service.OutletVisitListRepo.UpdateOutletVisitListColumnAt(ctx, request.SalesmanCode, request.CustID, "resume_at", request.CurrentTime, request.Date, request.Id)
	// add on_hold to be null as resume is updated
	service.OutletVisitListRepo.UpdateOutletVisitListColumnAt(ctx, request.SalesmanCode, request.CustID, "on_hold", nil, request.Date, request.Id)
}

func (service *VisitServiceImpl) LeaveVisit(ctx context.Context, request request.LeaveVisitRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	service.OutletVisitListRepo.UpdateOutletVisitListColumnAt(ctx, request.SalesmanCode, request.CustID, "leave_at", request.CurrentTime, request.Date, request.Id)
}

func (service *VisitServiceImpl) UnloadVisit(ctx context.Context, request request.OnholdVisitRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	service.OutletVisitListRepo.UpdateOutletVisitListColumnAt(ctx, request.SalesmanCode, request.CustID, "on_hold", request.CurrentTime, request.Date, request.Id)
	// add resume to be null as on_hold is updated
	service.OutletVisitListRepo.UpdateOutletVisitListColumnAt(ctx, request.SalesmanCode, request.CustID, "resume_at", nil, request.Date, request.Id)

}

func (service *VisitServiceImpl) OutletVisit(ctx context.Context, request request.OutletVisitRequest) {
	err := service.Validate.Struct(request)
	helper.ErrorPanic(err)

	service.OutletVisitListRepo.CreateOutletVisit(ctx, request.SalesmanCode, request.CustID, request.Date)
}

func (service *VisitServiceImpl) GetAllOutletBySalesCode(ctx context.Context, salesCode, custId, date, routeCode string) (responses []response.DestinationResponse) {
	result := service.RouteOutletRepo.GetAllOutletBySalesCode(ctx, salesCode, custId, date, routeCode)

	for _, row := range result {
		var res response.DestinationResponse
		helper.Automapper(row, &res)
		res.Status = row.Status
		responses = append(responses, res)
	}

	return responses
}

func (service *VisitServiceImpl) SummaryVisit(ctx context.Context, dataFilter entity.SummaryQueryFilter) (res response.SummaryResponse) {
	result, err := service.OutletVisitListRepo.GetSummary(ctx, dataFilter)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	if result.Start != nil {
		res.StartTime = *result.Start
	}

	if result.Finish != nil {
		res.EndTime = *result.Finish
	}

	return res
}

func (service *VisitServiceImpl) SummaryVisitStatus(ctx context.Context, dataFilter entity.SummaryQueryFilter) (res response.VisitStatusResponse) {
	result, err := service.OutletVisitListRepo.GetSummaryStatus(ctx, dataFilter)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	res = response.VisitStatusResponse{
		Planned:    result.Planned,
		Finished:   result.Finished,
		Skipped:    result.Skipped,
		OnProgress: result.OnProgress,
		OnHold:     result.OnHold,
	}

	return res
}

func (service *VisitServiceImpl) TravelList(ctx context.Context, outleVisitId int) (response response.TodoListResponse) {
	// selectColumns := []string{"arrive_at, unload_at, leave_at"}
	// data, err := service.OutletVisitListRepo.FindByOneColumn(ctx, selectColumns, "outlet_id", DestinationID)
	// if err != nil {
	// 	panic(exception.NewNotFoundError(err.Error()))
	// }

	// helper.Automapper(data, &response)
	data, err := service.OutletVisitListRepo.FindById(ctx, outleVisitId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return data

}

func (service *VisitServiceImpl) GetVisitOutletList(ctx context.Context, dataFilter entity.SummaryQueryFilter) (response []response.OutletVisitListResponse) {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	data, err := service.OutletVisitListRepo.FindAll(ctx, dataFilter)
	if err != nil {
		helper.ErrorPanic(err)
	}

	// remove /master when deploy to staging or production
	endpointURL := fmt.Sprintf("%s/v1/outlets?limit=999999", config.KongUrl)
	// log.Printf("Request URL: %s", endpointURL)

	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
	if err != nil {
		helper.ErrorPanic(err)
	}
	req.Header.Add("cust_id", "C220010001")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		helper.ErrorPanic(err)
	}
	defer resp.Body.Close()

	var result struct {
		Data   []model.OutletNew `json:"data"`
		Paging struct {
			TotalRecord int `json:"total_record"`
			PageCurrent int `json:"page_current"`
			PageLimit   int `json:"page_limit"`
			PageTotal   int `json:"page_total"`
		} `json:"paging"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helper.ErrorPanic(err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		helper.ErrorPanic(err)
	}

	for i := range data {
		found := false
		for _, outlet := range result.Data {
			fmt.Println("okok", data)
			if data[i].DestinationID == outlet.DestinationID {
				found = true
				data[i].DueDate = outlet.DueDate
				data[i].DestinationName = outlet.DestinationName
				data[i].DestinationAddress = outlet.Address1
			}
		}

		if !found {
			log.Printf("DestinationID not found: %d", data[i].DestinationID)
		}
	}

	// log.Printf("Result: %+v", data)
	helper.Automapper(data, &response)

	return response
}

func (service *VisitServiceImpl) GetSalesmanReport(ctx context.Context, dataFilter entity.SalesmanReportQueryFilter) entity.SalesmanReport {
	var report entity.SalesmanReport

	// Get visit data
	visits, err := service.OutletVisitListRepo.GetVisitsByDateAndSalesman(ctx, dataFilter)
	if err != nil {
		helper.ErrorPanic(err)
	}

	report.TotalOutlets = len(visits)
	skipReasonCount := make(map[string]int)

	for _, visit := range visits {
		var isPlanned bool
		existsInAdditional, err := service.OutletVisitListRepo.CheckRouteOutletAdditional(ctx, visit.RouteCode, visit.DestinationID)
		if err != nil {
			helper.ErrorPanic(err)
		}
		isPlanned = !existsInAdditional

		if visit.LeaveAt != nil {
			report.TotalVisit++
			if isPlanned {
				report.Visit.Planned++
			} else {
				report.Visit.NotPlanned++
			}
		} else {
			report.TotalNotVisit++
			if isPlanned {
				report.NotVisit.Planned++
			} else {
				report.NotVisit.NotPlanned++
			}
		}

		if visit.SkipAt != nil && visit.SkipReason != nil && *visit.SkipReason != "" {
			skipReasonCount[*visit.SkipReason]++
		}
	}

	totalNotVisited := report.TotalNotVisit
	for reason, count := range skipReasonCount {
		percentage := int(math.Round((float64(count) / float64(totalNotVisited)) * 100))
		report.NotVisitReasons = append(report.NotVisitReasons, entity.NotVisitReason{
			Reason:     reason,
			Count:      count,
			Percentage: percentage,
		})
	}

	return report
}
