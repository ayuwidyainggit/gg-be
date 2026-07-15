package service

import (
	"errors"
	"fmt"
	"master/adapter"
	"master/entity"
	"master/model"
	"master/pkg/config/env"
	"master/pkg/constant"
	"master/pkg/str"
	"master/pkg/structs"
	"master/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/exp/rand"
)

type SalesmanService interface {
	FindParentCustId(string) (entity.MCustomerResp, error)
	Detail(params entity.DetailSalesmanParams) (entity.SalesmanResponse, error)
	List(entity.SalesmanQueryFilter, string, string) (data []entity.SalesmanListResponse, total int, lastPage int, err error)
	LookupList(entity.SalesmanQueryFilter, string, string) (data []entity.SalesmanListResponse, total int, lastPage int, err error)
	Store(entity.CreateSalesmanBody) (entity.SalesmanResponse, error)
	Update(int64, entity.UpdateSalesmanRequest) error
	Delete(string, int64, int64) error
	LookupJobType() (data []entity.JobTypeLookupResponse, err error)
	LookupTaxOption() (data []entity.TaxOptionLookupResponse, err error)
	UpdateIsActive(custId string, empId int64, userId int64) (err error)
	UpdateDeActive(custId string, empId int64, userId int64) (err error)
	publishJobSalesman(job entity.PublishJob)
	GoSchedulerCustom(startDate string, url string) (err error)
}

func NewSalesmanService(salesmanRepository repository.SalesmanRepository, config env.ConfigEnv) *salesmanServiceImpl {
	return &salesmanServiceImpl{
		Config:             config,
		SalesmanRepository: salesmanRepository,
	}
}

var (
	GROUPTYPEPRODUCTLINE = 1
	GROUPTYPEBRAND       = 2
	GROUPTYPESUB         = 3
)

type salesmanServiceImpl struct {
	Config             env.ConfigEnv
	SalesmanRepository repository.SalesmanRepository
	// MProductRepository repository.MProductRepository
}

func (service *salesmanServiceImpl) FindParentCustId(custId string) (response entity.MCustomerResp, err error) {
	mCustomer, err := service.SalesmanRepository.FindOneParentCustId(custId)
	if err != nil {
		return response, err
	}

	if err = structs.Automapper(mCustomer, &response); err != nil {
		return response, err
	}

	return response, err
}

func (service *salesmanServiceImpl) Detail(params entity.DetailSalesmanParams) (response entity.SalesmanResponse, err error) {
	salesman, err := service.SalesmanRepository.FindOneByEmpIdAndCustId(params)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(salesman, &response)
	if err != nil {
		return response, err
	}

	details, err := service.SalesmanRepository.FindDetailByIdAndCustId(params)
	if err != nil {
		return response, err
	}
	for _, detail := range details {
		switch detail.GroupType {
		case 1:
			detailEntity := entity.SalesmanDetGroupProductLine{
				MSalesmanProductTypeID: detail.MSalesmanProductTypeID,
				RefID:                  detail.RefID,
				PlCode:                 detail.RefCode,
				PlName:                 detail.RefName,
				PlId:                   detail.PlId,
			}
			response.Details.ProductLine = append(response.Details.ProductLine, detailEntity)
		case 2:
			detailEntity := entity.SalesmanDetGroupBrand{
				MSalesmanProductTypeID: detail.MSalesmanProductTypeID,
				RefID:                  detail.RefID,
				BrandCode:              detail.RefCode,
				BrandName:              detail.RefName,
				PlId:                   detail.PlId,
			}
			response.Details.Brand = append(response.Details.Brand, detailEntity)
		case 3:
			detailEntity := entity.SalesmanDetGroupSubBrand{
				MSalesmanProductTypeID: detail.MSalesmanProductTypeID,
				RefID:                  detail.RefID,
				SBrand1Code:            detail.RefCode,
				SBrand1Name:            detail.RefName,
				PlId:                   detail.PlId,
			}
			response.Details.SubBrand = append(response.Details.SubBrand, detailEntity)
		}
	}
	tempArrProdLine := make([]entity.SalesmanDetGroupProductLine, 0)
	if len(response.Details.ProductLine) == 0 {
		response.Details.ProductLine = tempArrProdLine
	}
	tempArrBrand := make([]entity.SalesmanDetGroupBrand, 0)
	if len(response.Details.Brand) == 0 {
		response.Details.Brand = tempArrBrand
	}
	tempArrSubBrand := make([]entity.SalesmanDetGroupSubBrand, 0)
	if len(response.Details.SubBrand) == 0 {
		response.Details.SubBrand = tempArrSubBrand
	}

	if salesman.JobType != nil {
		jobTypeName := entity.GetJobTypeName(*response.JobType)
		response.JobTypeName = &jobTypeName
	}

	if salesman.TaxOption != nil {
		taxOptionName := entity.GetTaxOptionName(*response.TaxOption)
		response.TaxOptionName = &taxOptionName
	}

	if salesman.StartDate != nil {
		startDate := *salesman.StartDate
		response.StartDate = &startDate
	}

	if salesman.EndDate != nil {
		endDate := *salesman.EndDate
		response.EndDate = &endDate
	}

	return response, err
}

func (service *salesmanServiceImpl) List(dataFilter entity.SalesmanQueryFilter, custId, parentCustId string) (data []entity.SalesmanListResponse, total int, lastPage int, err error) {
	salesmans, total, lastPage, err := service.SalesmanRepository.FindAllByCustId(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range salesmans {
		var vResp entity.SalesmanListResponse
		structs.Automapper(row, &vResp)

		oprText := entity.ConvStringString(entity.OprType, vResp.OprType)
		oprTextCanvas := entity.ConvStringString(entity.OprType, vResp.OprTypeCanvas)
		vResp.OprTypeText = oprText
		vResp.OprTypeTextCanvas = &oprTextCanvas
		if *vResp.IsTakingOrder == true {
			vResp.WhNameView = vResp.WhName
		}
		if *vResp.IsActiveCanvas == true && *vResp.IsTakingOrder == true {
			vResp.WhNameView = vResp.WhNameView + "\n" + vResp.WhNameCanvas
		} else if *vResp.IsActiveCanvas == true {
			vResp.WhNameView = vResp.WhNameCanvas
		}

		detailCustId := custId
		if row.CustId != "" {
			detailCustId = row.CustId
		}
		details, err := service.SalesmanRepository.FindDetailById(row.EmpId, detailCustId)
		if err != nil {
			return nil, 0, 0, err
		}
		for _, detail := range details {
			switch detail.GroupType {
			case 1:
				detailEntity := entity.SalesmanDetGroupProductLine{
					MSalesmanProductTypeID: detail.MSalesmanProductTypeID,
					RefID:                  detail.RefID,
					PlCode:                 detail.RefCode,
					PlName:                 detail.RefName,
					PlId:                   detail.PlId,
				}

				vResp.Details.ProductLine = append(vResp.Details.ProductLine, detailEntity)
				// response.Details.ProductLine = append(response.Details.ProductLine, detailEntity)
			case 2:
				detailEntity := entity.SalesmanDetGroupBrand{
					MSalesmanProductTypeID: detail.MSalesmanProductTypeID,
					RefID:                  detail.RefID,
					BrandCode:              detail.RefCode,
					BrandName:              detail.RefName,
					PlId:                   detail.PlId,
				}
				vResp.Details.Brand = append(vResp.Details.Brand, detailEntity)
				// response.Details.Brand = append(response.Details.Brand, detailEntity)
			case 3:
				detailEntity := entity.SalesmanDetGroupSubBrand{
					MSalesmanProductTypeID: detail.MSalesmanProductTypeID,
					RefID:                  detail.RefID,
					SBrand1Code:            detail.RefCode,
					SBrand1Name:            detail.RefName,
					PlId:                   detail.PlId,
				}
				vResp.Details.SubBrand = append(vResp.Details.SubBrand, detailEntity)
				// response.Details.SubBrand = append(response.Details.SubBrand, detailEntity)
			}
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *salesmanServiceImpl) LookupList(dataFilter entity.SalesmanQueryFilter, custId, parentCustId string) (data []entity.SalesmanListResponse, total int, lastPage int, err error) {
	salesmans, total, lastPage, err := service.SalesmanRepository.FindAllByCustIdLookup(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range salesmans {
		var vResp entity.SalesmanListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *salesmanServiceImpl) Store(request entity.CreateSalesmanBody) (response entity.SalesmanResponse, err error) {
	service.SalesmanRepository.TrxBegin()

	paramsDetail := entity.DetailSalesmanParams{
		EmpId:        request.EmpId,
		CustId:       request.CustId,
		ParentCustId: request.ParentCustId,
	}

	timeNow := time.Now().In(time.UTC)
	if request.IsActiveCanvas {

		// service.SalesmanRepository.TrxBegin()

		var tmpWhCode string
		if len(request.WarehouseName) > 0 {

			letters := []rune("12345678910")
			randomStr := make([]rune, 3) // Generate a 3-character random string
			for i := range randomStr {
				randomStr[i] = letters[rand.Intn(len(letters))]
			}

			// Combine the Unix timestamp with the random alphabetic string
			result := string(randomStr)
			tmpEmpId := strconv.FormatInt(request.EmpId, 10) // get emp id
			tmpWhCode = result + tmpEmpId
		}

		warehouse, err := service.SalesmanRepository.FindOneByWarehouseCodeAndCustId(tmpWhCode, request.CustId)
		if err == nil {
			response.WhId = warehouse.WarehouseId
			// return response, errors.New("wh_code: " + warehouse.WarehouseCode + " is already exists")
		} else {
			warehouseData := model.Warehouse{}
			err = structs.Automapper(request, &warehouseData)
			if err != nil {
				return response, err
			}

			tmpStockType := "G"
			warehouseData.WarehouseCode = tmpWhCode
			warehouseData.IsActive = true
			warehouseData.CreatedAt = &timeNow
			warehouseData.CreatedBy = &request.CreatedBy
			warehouseData.UpdatedAt = &timeNow
			warehouseData.UpdatedBy = &request.CreatedBy
			warehouseData.CustId = request.CustId
			warehouseData.StockType = &tmpStockType

			defer func() {
				if p := recover(); p != nil {
					service.SalesmanRepository.TrxRollback()
					panic(p) // re-throw panic after Rollback
				}
			}()
			warehouseId, err := service.SalesmanRepository.StoreWarehouse(warehouseData)
			if err != nil {
				service.SalesmanRepository.TrxRollback()
				return response, err
			}

			response.WhId = warehouseId
		}

		var salesmanCanvasData model.SalesmanCanvas
		err = structs.Automapper(request, &salesmanCanvasData)
		if err != nil {
			return response, err
		}

		salesmanCanvasData.CreatedAt = &timeNow
		salesmanCanvasData.CreatedBy = &request.CreatedBy
		salesmanCanvasData.UpdatedAt = &timeNow
		salesmanCanvasData.UpdatedBy = &request.CreatedBy
		salesmanCanvasData.CustId = request.CustId
		salesmanCanvasData.WhId = response.WhId
		salesmanCanvasData.OprTypeCanvas = "C"

		defer func() {
			if p := recover(); p != nil {
				service.SalesmanRepository.TrxRollback()
				panic(p) // re-throw panic after Rollback
			}
		}()

		salesman, err := service.SalesmanRepository.FindOneSalesmanCanvasByEmpIdAndCustId(request.EmpId, request.CustId)

		if salesman.EmpId != 0 {
			// return response, errors.New("emp_id: " + fmt.Sprintf("%d", salesman.EmpId) + " is already exists")
		} else {
			empIdSalesmanCanvas, err := service.SalesmanRepository.StoreSalesmanCanvas(salesmanCanvasData)
			if err != nil {
				service.SalesmanRepository.TrxRollback()
				return response, err
			}
			response.EmpId = empIdSalesmanCanvas
		}

		// service.SalesmanRepository.TrxCommit()
	}

	// if request.IsTakingOrder {
	salesman, err := service.SalesmanRepository.FindOneByEmpIdAndCustId(paramsDetail)
	if err == nil {
		return response, fmt.Errorf("emp_id: %d (%s) is already exists", salesman.EmpId, salesman.EmpName)
	}

	// timeNow := time.Now().In(time.UTC)
	var salesmanData model.Salesman
	// log.Println("request.TransDate:", request.TransDate)
	if request.TransDate != nil {
		if *request.TransDate != "" {
			transdate, err := str.DateStrToRfc3339String(*request.TransDate)
			if err != nil {
				return response, err
			}
			request.TransDate = &transdate
		} else {
			request.TransDate = nil
		}
	}

	if request.StartDate != "" {
		startDate, err := str.DateStrToRfc3339String(request.StartDate)
		if err != nil {
			return response, err
		}
		request.StartDate = startDate
	} else {
		request.StartDate = ""
	}

	if request.JobType != "P" {
		if request.EndDate != nil {
			endDate, err := str.DateStrToRfc3339String(*request.EndDate)
			if err != nil {
				return response, err
			}

			// Validasi: EndDate harus lebih dari StartDate
			if request.StartDate != "" && endDate <= request.StartDate {
				return response, fmt.Errorf("EndDate must be greater than StartDate")
			}
			request.EndDate = &endDate
		}

	} else {
		request.EndDate = nil
	}

	err = structs.Automapper(request, &salesmanData)
	if err != nil {
		return response, err
	}

	// log.Println("salesman.TransDate:", salesman.TransDate)
	salesmanData.CreatedAt = &timeNow
	salesmanData.CreatedBy = &request.CreatedBy
	salesmanData.UpdatedAt = &timeNow
	salesmanData.UpdatedBy = &request.CreatedBy
	salesmanData.CustId = request.CustId
	salesmanData.OprType = "O"

	// service.SalesmanRepository.TrxBegin()

	defer func() {
		if p := recover(); p != nil {
			service.SalesmanRepository.TrxRollback()
			panic(p) // re-throw panic after Rollback
		}
	}()
	empId, err := service.SalesmanRepository.Store(salesmanData)
	if err != nil {
		service.SalesmanRepository.TrxRollback()
		return response, err
	}

	response.EmpId = empId
	for _, detail := range request.Details.ProductLine {
		modelDetail := model.SalesmanDetail{
			CustId:    &request.CustId,
			EmpId:     &empId,
			GroupType: &GROUPTYPEPRODUCTLINE,
			RefID:     &detail.RefID,
		}
		err := service.SalesmanRepository.StoreDetail(modelDetail)
		if err != nil {
			service.SalesmanRepository.TrxRollback()
			return response, err
		}
	}
	for _, detail := range request.Details.Brand {
		modelDetail := model.SalesmanDetail{
			CustId:    &request.CustId,
			EmpId:     &empId,
			GroupType: &GROUPTYPEBRAND,
			RefID:     &detail.RefID,
		}
		err := service.SalesmanRepository.StoreDetail(modelDetail)
		if err != nil {
			service.SalesmanRepository.TrxRollback()
			return response, err
		}
	}
	for _, detail := range request.Details.SubBrand {
		modelDetail := model.SalesmanDetail{
			CustId:    &request.CustId,
			EmpId:     &empId,
			GroupType: &GROUPTYPESUB,
			RefID:     &detail.RefID,
		}
		err := service.SalesmanRepository.StoreDetail(modelDetail)
		if err != nil {
			service.SalesmanRepository.TrxRollback()
			return response, err
		}
	}
	// service.SalesmanRepository.TrxCommit()
	// }
	service.SalesmanRepository.TrxCommit()

	startDate := ""
	endDate := ""
	if request.StartDate != "" {
		startDate = request.StartDate
	}

	if request.EndDate != nil {
		endDate = *request.EndDate
	}

	go goSchedulerSalesman(service, startDate, endDate, request.CustId, request.CreatedBy, request.EmpId)

	return response, err
}

func (service *salesmanServiceImpl) Update(salesmanId int64, request entity.UpdateSalesmanRequest) (err error) {
	timeNow := time.Now().In(time.UTC)

	service.SalesmanRepository.TrxBegin()

	defer func() {
		if p := recover(); p != nil {
			service.SalesmanRepository.TrxRollback()
			panic(p) // re-throw panic after Rollback
		}
	}()

	if *request.IsActiveCanvas {
		salesman, _ := service.SalesmanRepository.FindOneSalesmanCanvasByEmpIdAndCustId(salesmanId, request.CustId)
		// fmt.Println(">>>>", salesman.EmpId)

		if salesman.EmpId == 0 {
			var tmpWhCode string
			if len(request.WarehouseName) > 0 {

				letters := []rune("12345678910")
				randomStr := make([]rune, 3) // Generate a 3-character random string
				for i := range randomStr {
					randomStr[i] = letters[rand.Intn(len(letters))]
				}

				// Combine the Unix timestamp with the random alphabetic string
				result := string(randomStr)
				tmpEmpId := strconv.FormatInt(salesmanId, 10) // get salesman id
				tmpWhCode = result + tmpEmpId
			}

			warehouse, err := service.SalesmanRepository.FindOneByWarehouseCodeAndCustId(tmpWhCode, request.CustId)
			if err == nil {
				return errors.New("wh_code: " + warehouse.WarehouseCode + " is already exists")
			}

			warehouseData := model.Warehouse{}
			err = structs.Automapper(request, &warehouseData)
			if err != nil {
				return err
			}

			tmpStockType := "G"
			warehouseData.WarehouseCode = tmpWhCode
			warehouseData.IsActive = true
			warehouseData.CreatedAt = &timeNow
			warehouseData.UpdatedAt = &timeNow
			warehouseData.UpdatedBy = &request.UpdatedBy
			warehouseData.CustId = request.CustId
			warehouseData.StockType = &tmpStockType

			defer func() {
				if p := recover(); p != nil {
					service.SalesmanRepository.TrxRollback()
					panic(p) // re-throw panic after Rollback
				}
			}()
			warehouseId, err := service.SalesmanRepository.StoreWarehouse(warehouseData)
			if err != nil {
				service.SalesmanRepository.TrxRollback()
				return err
			}

			var salesmanCanvasData model.SalesmanCanvas
			err = structs.Automapper(request, &salesmanCanvasData)
			if err != nil {
				return err
			}

			salesmanCanvasData.EmpId = salesmanId
			salesmanCanvasData.CreatedAt = &timeNow
			salesmanCanvasData.UpdatedAt = &timeNow
			salesmanCanvasData.UpdatedBy = &request.UpdatedBy
			salesmanCanvasData.CustId = request.CustId
			salesmanCanvasData.WhId = warehouseId
			salesmanCanvasData.OprTypeCanvas = "C"

			defer func() {
				if p := recover(); p != nil {
					service.SalesmanRepository.TrxRollback()
					panic(p) // re-throw panic after Rollback
				}
			}()

			salesman, err := service.SalesmanRepository.FindOneSalesmanCanvasByEmpIdAndCustId(salesmanId, request.CustId)

			if salesman.EmpId != 0 {
				return errors.New("emp_id: " + fmt.Sprintf("%d", salesman.EmpId) + " is already exists")
			}

			empIdSalesmanCanvas, err := service.SalesmanRepository.StoreSalesmanCanvas(salesmanCanvasData)
			if err != nil {
				service.SalesmanRepository.TrxRollback()
				return err
			}
			fmt.Println(empIdSalesmanCanvas)
		} else {
			err = service.SalesmanRepository.UpdateCanvas(salesmanId, request)
			if err != nil {
				service.SalesmanRepository.TrxRollback()
				return err
			}
		}
	} else {
		err = service.SalesmanRepository.UpdateCanvas(salesmanId, request)
		if err != nil {
			service.SalesmanRepository.TrxRollback()
			return err
		}
	}

	if !request.IsTakingOrder {
		err = service.SalesmanRepository.UpdateIsTakingOrder(salesmanId, request.CustId)
		if err != nil {
			service.SalesmanRepository.TrxRollback()
			return err
		}
	}
	// else {
	// // parse time format YYYY-mm-dd to Rfc3339
	if request.StartDate != nil {
		startDate, err := str.DateStrToRfc3339String(*request.StartDate)
		if err != nil {
			return err
		}
		request.StartDate = &startDate
	}

	if request.JobType != "P" {

		if request.EndDate != nil {
			endDate, err := str.DateStrToRfc3339String(*request.EndDate)
			if err != nil {
				return err
			}

			// Validasi: EndDate harus lebih dari StartDate
			if request.StartDate != nil && endDate < *request.StartDate {
				return fmt.Errorf("EndDate must be greater than StartDate")
			}

			request.EndDate = &endDate
		}
	} else {
		request.EndDate = nil
	}
	// request.OprType = "O"
	err = service.SalesmanRepository.Update(salesmanId, request)
	if err != nil {
		service.SalesmanRepository.TrxRollback()
		return err
	}
	DetailIds := []int64{}
	for _, detail := range request.Details.ProductLine {
		if detail.MSalesmanProductTypeID != nil {
			DetailIds = append(DetailIds, *detail.MSalesmanProductTypeID)
		}
	}

	for _, detail := range request.Details.Brand {
		if detail.MSalesmanProductTypeID != nil {
			DetailIds = append(DetailIds, *detail.MSalesmanProductTypeID)
		}
	}

	for _, detail := range request.Details.SubBrand {
		if detail.MSalesmanProductTypeID != nil {
			DetailIds = append(DetailIds, *detail.MSalesmanProductTypeID)
		}
	}

	if err = service.SalesmanRepository.DeleteDetails(salesmanId, request.CustId); err != nil {
		service.SalesmanRepository.TrxRollback()
		return err
	}

	for _, detail := range request.Details.ProductLine {
		modelDetail := model.SalesmanDetail{
			CustId:    &request.CustId,
			EmpId:     &salesmanId,
			GroupType: &GROUPTYPEPRODUCTLINE,
			RefID:     detail.RefID,
		}
		err := service.SalesmanRepository.StoreDetail(modelDetail)
		if err != nil {
			service.SalesmanRepository.TrxRollback()
			return err
		}
	}

	for _, detail := range request.Details.Brand {
		modelDetail := model.SalesmanDetail{
			CustId:    &request.CustId,
			EmpId:     &salesmanId,
			GroupType: &GROUPTYPEBRAND,
			RefID:     detail.RefID,
		}
		err := service.SalesmanRepository.StoreDetail(modelDetail)
		if err != nil {
			service.SalesmanRepository.TrxRollback()
			return err
		}
	}

	for _, detail := range request.Details.SubBrand {
		modelDetail := model.SalesmanDetail{
			CustId:    &request.CustId,
			EmpId:     &salesmanId,
			GroupType: &GROUPTYPESUB,
			RefID:     detail.RefID,
		}
		err := service.SalesmanRepository.StoreDetail(modelDetail)
		if err != nil {
			service.SalesmanRepository.TrxRollback()
			return err
		}
	}
	// }

	service.SalesmanRepository.TrxCommit()

	startDate := ""
	endDate := ""
	if request.StartDate != nil {
		startDate = *request.StartDate
	}

	if request.EndDate != nil {
		endDate = *request.EndDate
	}

	go goSchedulerSalesman(service, startDate, endDate, request.CustId, request.UpdatedBy, salesmanId)

	return err
}

func goSchedulerSalesman(service *salesmanServiceImpl, startDate string, endDate string, CustId string, UserId int64, EmpId int64) {
	payloadRmq := entity.UpdateIsActiveRequest{
		CustId: CustId,
		UserId: UserId,
		EmpId:  EmpId,
	}

	if startDate != "" {

		layout := time.RFC3339
		parsedDate, err := time.Parse(layout, startDate)
		if err != nil {
			return
		}

		startDateStr := parsedDate.Format("2006-01-02") + "T00:00:00+07:00"

		pubJobStart := entity.PublishJob{
			JobName:   "Publish Scheduler Is Active",
			JobDesc:   fmt.Sprintf("scheduler is active, emp_id: %d", payloadRmq.EmpId),
			JobType:   constant.JOB_TYPE_ONE_TIME,
			Task:      constant.JOB_TASK_HTTP_REQ,
			RunAt:     startDateStr,
			Url:       service.Config.Get("MASTER_SERVICE_URL") + constant.SALESMAN_ISACTIVE,
			Payload:   structs.StructToJson(payloadRmq),
			CreatedBy: "system",
		}
		go service.publishJobSalesman(pubJobStart)
	}

	if endDate != "" {

		layout := time.RFC3339
		parsedDate, err := time.Parse(layout, endDate)
		if err != nil {
			return
		}

		// nextDay := parsedDate.AddDate(0, 0, 1)

		EndDateStr := parsedDate.Format("2006-01-02") + "T00:00:00+07:00"

		// fmt.Println("end date str", EndDateStr)

		pubJobStart := entity.PublishJob{
			JobName:   "Publish Scheduler De Active",
			JobDesc:   fmt.Sprintf("scheduler is active, emp_id: %d", payloadRmq.EmpId),
			JobType:   constant.JOB_TYPE_ONE_TIME,
			Task:      constant.JOB_TASK_HTTP_REQ,
			RunAt:     EndDateStr,
			Url:       service.Config.Get("MASTER_SERVICE_URL") + constant.SALESMAN_DEACTIVE,
			Payload:   structs.StructToJson(payloadRmq),
			CreatedBy: "system",
		}
		go service.publishJobSalesman(pubJobStart)
	}
}

func (service *salesmanServiceImpl) Delete(custId string, empId int64, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsKeyExists(salesmanId, custId, "salesman_id1")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("salesman_id is still being used")
	// }

	err = service.SalesmanRepository.Delete(custId, empId, userId)
	if err != nil {
		return err
	}

	return err
}

func (service *salesmanServiceImpl) LookupJobType() (data []entity.JobTypeLookupResponse, err error) {

	data = entity.JobType

	return data, err
}

func (service *salesmanServiceImpl) LookupTaxOption() (data []entity.TaxOptionLookupResponse, err error) {

	data = entity.TaxOption

	return data, err
}

func (service *salesmanServiceImpl) UpdateIsActive(custId string, empId int64, userId int64) (err error) {

	err = service.SalesmanRepository.UpdateIsActive(empId, custId)
	if err != nil {
		return err
	}

	return err
}

func (service *salesmanServiceImpl) UpdateDeActive(custId string, empId int64, userId int64) (err error) {

	err = service.SalesmanRepository.UpdateDeActive(empId, custId)
	if err != nil {
		return err
	}

	return err
}

func (service *salesmanServiceImpl) publishJobSalesman(job entity.PublishJob) {
	client := adapter.HttpClientInfo{
		Url:     service.Config.Get("CRONJOB_SERVICE_URL") + "/v1/jobs",
		Method:  "POST",
		Payload: job,
	}

	res, err := client.Dispatch()
	if err != nil {
		log.Errorf("Dispatch failed: %v", err)
		return
	}

	if res.StatusCode() >= 400 {
		log.Errorf("HTTP request failed with status %d: %s", res.StatusCode(), res.Body())
	}
}

func (service *salesmanServiceImpl) GoSchedulerCustom(startDate string, url string) (err error) {
	fmt.Println("startDate", startDate)
	service.SalesmanRepository.CheckDate()

	if startDate != "" {
		pubJobStart := entity.PublishJob{
			JobName:   "Publish Scheduler Is Active",
			JobDesc:   fmt.Sprintf("scheduler custom"),
			JobType:   constant.JOB_TYPE_ONE_TIME,
			Task:      constant.JOB_TASK_HTTP_REQ,
			RunAt:     startDate,
			Url:       service.Config.Get("MASTER_SERVICE_URL") + url,
			Payload:   structs.StructToJson(map[string]string{"start_date": startDate}),
			CreatedBy: "system",
		}
		go service.publishJobSalesman(pubJobStart)
	}

	return nil
}
