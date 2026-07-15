package service

import (
	"errors"
	"fmt"
	"master/adapter"
	"master/entity"
	"master/model"
	"master/pkg/config/env"
	"master/pkg/constant"
	"master/pkg/structs"
	"master/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DistPriceService interface {
	Detail(int64, string) (entity.DistPriceDetailResp, error)
	List(entity.DistPriceQueryFilter, string) (data []entity.DistPriceListResponse, total int, lastPage int, err error)
	LookupList(entity.DistPriceQueryFilter, string, string) (data []entity.DistPriceLookupResponse, total int, lastPage int, err error)
	LookupProductList(entity.DistPriceQueryFilter, string) (data []entity.DistPriceLookupProResp, total int, lastPage int, err error)
	Store(entity.CreateDistPriceBody) (entity.DistPriceResponse, error)
	Update(entity.TokenMetadata, int64, entity.UpdateDistPriceRequest) error
	Delete(entity.TokenMetadata, int64, int64) error
	PublishOrInactive(entity.PublishUnpublishDistPriceReq) error
}

func NewDistPriceService(
	config env.ConfigEnv,
	distPriceRepository repository.DistPriceRepository,
	transPriceRepository repository.MTransactionPriceRepository,
) *distPriceServiceImpl {
	return &distPriceServiceImpl{
		Config:                      config,
		DistPriceRepository:         distPriceRepository,
		MTransactionPriceRepository: transPriceRepository,
	}
}

type distPriceServiceImpl struct {
	Config                      env.ConfigEnv
	DistPriceRepository         repository.DistPriceRepository
	MTransactionPriceRepository repository.MTransactionPriceRepository
}

func (service *distPriceServiceImpl) Detail(distPriceId int64, custId string) (response entity.DistPriceDetailResp, err error) {
	distPrice, err := service.DistPriceRepository.FindOneByDistPriceIdAndCustId(distPriceId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(distPrice, &response)
	if err != nil {
		return response, err
	}

	if distPrice.StartDate != nil {
		startDate := distPrice.StartDate.Format("2006-01-02")
		response.StartDate = startDate
	}
	if distPrice.EndDate != nil {
		endDate := distPrice.EndDate.Format("2006-01-02")
		response.EndDate = endDate
	}

	statusResp := entity.DistPriceListResponse{
		Status: response.Status,
	}
	response.StatusDesc = statusResp.GetDistPriceStatusDesc()

	return response, err
}

func (service *distPriceServiceImpl) List(dataFilter entity.DistPriceQueryFilter, custId string) (data []entity.DistPriceListResponse, total int, lastPage int, err error) {

	distPrices, total, lastPage, err := service.DistPriceRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range distPrices {
		var vResp entity.DistPriceListResponse
		structs.Automapper(row, &vResp)
		if row.StartDate != nil {
			startDate := row.StartDate.Format("2006-01-02")
			vResp.StartDate = startDate
		}
		if row.EndDate != nil {
			endDate := row.EndDate.Format("2006-01-02")
			vResp.EndDate = endDate
		}
		statusResp := entity.DistPriceListResponse{
			Status: row.Status,
		}
		vResp.StatusDesc = statusResp.GetDistPriceStatusDesc()
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *distPriceServiceImpl) LookupList(dataFilter entity.DistPriceQueryFilter, parentCustId, custId string) (data []entity.DistPriceLookupResponse, total int, lastPage int, err error) {

	distPrices, total, lastPage, err := service.DistPriceRepository.FindAllByCustIdLookupMode(dataFilter, parentCustId, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range distPrices {
		var vResp entity.DistPriceLookupResponse
		structs.Automapper(row, &vResp)
		if row.StartDate != nil {
			startDate := row.StartDate.Format("2006-01-02")
			vResp.StartDate = startDate
		}
		if row.EndDate != nil {
			endDate := row.EndDate.Format("2006-01-02")
			vResp.EndDate = endDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *distPriceServiceImpl) LookupProductList(dataFilter entity.DistPriceQueryFilter, custId string) (
	data []entity.DistPriceLookupProResp, total int, lastPage int, err error) {

	distPrices, total, lastPage, err := service.DistPriceRepository.FindAllByCustIdLookupProduct(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range distPrices {
		var vResp entity.DistPriceLookupProResp
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *distPriceServiceImpl) Store(request entity.CreateDistPriceBody) (response entity.DistPriceResponse, err error) {
	timeNow := time.Now().In(time.UTC)

	asiaJkt, _ := time.LoadLocation("Asia/Jakarta")
	startDateStr := request.StartDate
	startDate, err := time.ParseInLocation(constant.DATE_LAYOUT_YYYY_MM_DD, startDateStr, asiaJkt)
	if err != nil {
		return response, err
	}

	// endDate, err := time.Parse("2006-01-02", request.EndDate)
	// if err != nil {
	// 	return response, err
	// }

	endDateOld := ""
	if request.DistPricePriceIdOld != nil {
		if *request.DistPricePriceIdOld > 0 {
			endDateOld = startDate.AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
		}
	}

	// validate product
	countProduct := service.DistPriceRepository.CountByProIDAndCustID(request.ProId, request.CustId)
	log.Info("countProduct:", countProduct)
	if countProduct < 1 {
		return response, fmt.Errorf("pro_id %d not found", request.ProId)
	}

	var distPriceData model.DistPrice
	err = structs.Automapper(request, &distPriceData)
	if err != nil {
		return response, err
	}

	// distPriceData.StartDate = &startDate
	// distPriceData.EndDate = &endDate
	distPriceData.CreatedAt = &timeNow
	distPriceData.CreatedBy = &request.CreatedBy
	distPriceData.UpdatedAt = &timeNow
	distPriceData.UpdatedBy = &request.CreatedBy

	service.DistPriceRepository.TrxBegin()

	defer func() {
		if p := recover(); p != nil {
			service.DistPriceRepository.TrxRollback()
			panic(p) // re-throw panic after Rollback
		}
	}()

	distPriceId, err := service.DistPriceRepository.Store(distPriceData)
	if err != nil {
		service.DistPriceRepository.TrxRollback()
		return response, err
	}

	response.DistPricePriceId = distPriceId

	if endDateOld != "" {
		requestOld := entity.UpdateDistPriceRequest{
			CustId:    distPriceData.CustId,
			UpdatedBy: request.CreatedBy,
			EndDate:   &endDateOld,
		}
		err = service.DistPriceRepository.Update(*distPriceData.DistPriceIdOld, requestOld)
		if err != nil {
			log.Info("no old data needs to be updated:", err.Error())
		}
	}

	service.DistPriceRepository.TrxCommit()

	payloadRmq := entity.PublishUnpublishDistPriceReq{
		CustID:      request.CustId,
		DistPriceID: distPriceId,
		Status:      10,
	}
	pubJobStart := entity.PublishJob{
		JobName:   "Publish Dist Price",
		JobDesc:   fmt.Sprintf("publish dist price, dist_price_id: %d", distPriceId),
		JobType:   constant.JOB_TYPE_ONE_TIME,
		Task:      constant.JOB_TASK_HTTP_REQ,
		RunAt:     startDate.Format(time.RFC3339),
		Url:       service.Config.Get("MASTER_SERVICE_URL") + constant.DIST_PRICE_PUBLISH_UNPUBLISH,
		Payload:   structs.StructToJson(payloadRmq),
		CreatedBy: "system",
	}
	go service.publishJob(pubJobStart)

	return response, err
}

func (service *distPriceServiceImpl) publishJob(job entity.PublishJob) {
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

func (service *distPriceServiceImpl) Update(claims entity.TokenMetadata, distPriceId int64, request entity.UpdateDistPriceRequest) (err error) {
	if claims.CustId != claims.ParentCustId {
		return errors.New("this endpoint only for principal")
	}

	service.DistPriceRepository.TrxBegin()
	err = service.DistPriceRepository.Update(distPriceId, request)
	if err != nil {
		service.DistPriceRepository.TrxRollback()
		return err
	}
	service.DistPriceRepository.TrxCommit()
	return err
}

func (service *distPriceServiceImpl) Delete(claims entity.TokenMetadata, distPriceId int64, userId int64) (err error) {

	log.Info("distPriceServiceImpl, Delete, claims:", structs.StructToJson(claims))
	if claims.CustId != claims.ParentCustId {
		return errors.New("this endpoint only for principal")
	}

	distPrice, err := service.DistPriceRepository.FindOneByDistPriceIdAndCustId(distPriceId, claims.CustId)
	if err != nil {
		return err
	}

	log.Info("distPriceServiceImpl, Delete, distPrice:", structs.StructToJson(distPrice))

	service.DistPriceRepository.TrxBegin()

	if distPrice.DistPriceIdOld != nil && *distPrice.DistPriceIdOld != 0 {
		// update dist price old, set end date = null
		updateDistPriceOld := entity.UpdateDistPriceRequest{
			CustId:    claims.CustId,
			UpdatedBy: userId,
		}

		err = service.DistPriceRepository.UpdateEndDateNullByDistPriceId(*distPrice.DistPriceIdOld, updateDistPriceOld)
		if err != nil {
			service.DistPriceRepository.TrxRollback()
			log.Info("error:", err.Error())
			return err
		}

	}

	err = service.DistPriceRepository.Delete(claims.CustId, distPriceId, claims.UserId)
	if err != nil {
		service.DistPriceRepository.TrxRollback()
		return err
	}

	service.DistPriceRepository.TrxCommit()
	return err
}

func (service *distPriceServiceImpl) PublishOrInactive(request entity.PublishUnpublishDistPriceReq) (err error) {
	log.Info("distPriceServiceImpl, PublishOrInactive, request -> ", structs.StructToJson(request))

	distPrice, err := service.DistPriceRepository.FindOneByDistPriceIdAndCustId(request.DistPriceID, request.CustID)
	if err != nil {
		log.Error("distPriceServiceImpl, PublishOrInactive, error:", err.Error())
		return fmt.Errorf("dist_price_id %d not found", request.DistPriceID)
	}

	// scheduler event publish
	if request.Status == 10 {
		// setup insert transaction price data
		transPrice := setDistPriceToTransPriceData(distPrice, 0)
		if err = service.MTransactionPriceRepository.StoreIfNotExists(&transPrice); err != nil {
			return err
		}
	}

	if err = service.DistPriceRepository.UpdateStatusByRMQ(request); err != nil {
		log.Info("distPriceServiceImpl, UpdateStatusByRMQ, error:", err.Error())
	}

	return nil
}

func setDistPriceToTransPriceData(price model.DistPriceDetail, distributorID int64) (transPrice model.MTransactionPrice) {
	objectID := primitive.NewObjectID() // Generate a new ObjectID
	objectIDString := objectID.Hex()    // Convert ObjectID to string
	strDistPriceID := strconv.Itoa(int(price.DistPriceId))
	transPrice = model.MTransactionPrice{
		CustID:             price.CustId,
		TransactionPriceID: objectIDString,
		ProID:              price.ProId,
		PurchPrice1:        price.PurchPrice1,
		PurchPrice2:        price.PurchPrice2,
		PurchPrice3:        price.PurchPrice3,
		SellPrice1:         0,
		SellPrice2:         0,
		SellPrice3:         0,
		Source:             1, // distributor price = 1
		CreatedBy:          "scheduler",
		CreatedAt:          time.Now().UTC(),
		StartDate:          price.StartDate,
		EndDate:            nil,
		Coverage:           "D",
		DistributorID:      distributorID,
		PriceGroupReff:     price.DistPriceGroupId,
		ReferenceID:        strDistPriceID,
		OutletID:           0,
	}

	return transPrice
}
