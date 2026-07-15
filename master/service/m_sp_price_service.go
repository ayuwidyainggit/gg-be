package service

import (
	"errors"
	"fmt"
	"master/adapter"
	"master/entity"
	"master/model"
	"master/pkg/config/env"
	"master/pkg/constant"
	"master/pkg/errmsg"
	"master/pkg/str"
	"master/pkg/structs"
	"master/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpService interface {
	Preview(request entity.CreateMSpPriceBody) (response entity.PreviewMSpPriceResp, err error)
	Store(request entity.CreateMSpPriceBody) (response entity.MSpPriceResponse, err error)
	Detail(params entity.MSpPriceParams) (response entity.MSpPriceWithDetailResp, err error)
	List(dataFilter entity.MSpPriceQueryFilter, custId string) (data []entity.MSpPriceResponse, total int, lastPage int, err error)
	Delete(custId string, sppriceId string) (err error)
	Update(sppriceId string, request entity.UpdateMSpPriceBody) (err error)
	Cancel(params entity.MSpPriceCancelParams) (err error)
	PublishOrInactive(entity.PublishUnpublishSPriceReq) error
}

func NewSpPriceService(
	config env.ConfigEnv,
	spRepository repository.SpPriceRepository,
	specialPriceGroupRepository repository.SpecialPriceGroupRepository,
	outletRepository repository.OutletRepository,
	outletTypeRepository repository.OutletTypeRepository,
	outletGroupRepository repository.OutletGroupRepository,
	transPriceRepository repository.MTransactionPriceRepository,
) *spServiceImpl {
	return &spServiceImpl{
		Config:                      config,
		SpRepository:                spRepository,
		SpecialPriceGroupRepository: specialPriceGroupRepository,
		OutletRepository:            outletRepository,
		OutletTypeRepository:        outletTypeRepository,
		OutletGroupRepository:       outletGroupRepository,
		MTransactionPriceRepository: transPriceRepository,
	}
}

type spServiceImpl struct {
	Config                      env.ConfigEnv
	SpRepository                repository.SpPriceRepository
	SpecialPriceGroupRepository repository.SpecialPriceGroupRepository
	OutletRepository            repository.OutletRepository
	OutletTypeRepository        repository.OutletTypeRepository
	OutletGroupRepository       repository.OutletGroupRepository
	MTransactionPriceRepository repository.MTransactionPriceRepository
}

const (
	SALESTEAM   = 1
	OUTLETTYPE  = 5
	OUTLETGROUP = 10
	OUTLET      = 15
	OUTPUT      = 20
)

type GenerateOutputParams struct {
	Outlets          []model.Outlet
	OutletTypesReq   []entity.MSpPriceDet
	OutletGroupsReq  []entity.MSpPriceDet
	OutletsReq       []entity.MSpPriceDet
	MasterSellPrice1 float64
	MasterSellPrice2 float64
	MasterSellPrice3 float64
	NewSellPrice1    float64
	NewSellPrice2    float64
	NewSellPrice3    float64
	CustID           string
	ParentCustID     string
}

func (service *spServiceImpl) GetOutputByPriceGroupIdOutletTypeOutletGroupOutlet(params GenerateOutputParams) (output []entity.OutputSprice, err error) {

	for _, row := range params.Outlets {
		// default price based on master product
		output = append(output, entity.OutputSprice{
			OutletID:         row.OutletId,
			OutletCode:       *row.OutletCode,
			OutletName:       *row.OutletName,
			MasterSellPrice1: params.MasterSellPrice1,
			MasterSellPrice2: params.MasterSellPrice2,
			MasterSellPrice3: params.MasterSellPrice3,
			NewSellPrice1:    params.NewSellPrice1,
			NewSellPrice2:    params.NewSellPrice2,
			NewSellPrice3:    params.NewSellPrice3,
		})

		// prioritized based on outlet type
		for _, rowOtType := range params.OutletTypesReq {
			refID := int64(*rowOtType.RefID)
			otTypeID := int64(*row.OtTypeId)

			// validate ref_id first
			_, err := service.OutletTypeRepository.FindOneByOutletTypeIdAndCustId(refID, params.ParentCustID)
			if err != nil {
				outletTypeIDStr := strconv.Itoa(int(refID))
				return output, errors.New("Outlet Type ID: " + outletTypeIDStr + " not found")
			}

			if refID == otTypeID {
				for i, rowResp := range output {
					if rowResp.OutletID == row.OutletId {
						// Remove the element at index i
						output = append(output[:i], output[i+1:]...)
						break // Exit the inner loop after deletion
					}
				}
				output = append(output, entity.OutputSprice{
					OutletID:         row.OutletId,
					OutletCode:       *row.OutletCode,
					OutletName:       *row.OutletName,
					MasterSellPrice1: params.MasterSellPrice1,
					MasterSellPrice2: params.MasterSellPrice2,
					MasterSellPrice3: params.MasterSellPrice3,
					NewSellPrice1:    *rowOtType.NewSellPrice1,
					NewSellPrice2:    *rowOtType.NewSellPrice2,
					NewSellPrice3:    *rowOtType.NewSellPrice3,
				})
			}
		}

		// prioritized based on outlet group
		for _, rowOtGroup := range params.OutletGroupsReq {
			refID := int64(*rowOtGroup.RefID)
			otGroupID := int64(*row.OtGrpId)

			// validate ref_id first
			_, err := service.OutletGroupRepository.FindOneByOutletGroupIdAndCustId(refID, params.ParentCustID)
			if err != nil {
				outletGroupIDStr := strconv.Itoa(int(refID))
				return output, errors.New("Outlet Group ID: " + outletGroupIDStr + " not found")
			}

			if refID == otGroupID {
				for i, rowResp := range output {
					if rowResp.OutletID == row.OutletId {
						// Remove the element at index i
						output = append(output[:i], output[i+1:]...)
						break // Exit the inner loop after deletion
					}
				}
				output = append(output, entity.OutputSprice{
					OutletID:         row.OutletId,
					OutletCode:       *row.OutletCode,
					OutletName:       *row.OutletName,
					MasterSellPrice1: params.MasterSellPrice1,
					MasterSellPrice2: params.MasterSellPrice2,
					MasterSellPrice3: params.MasterSellPrice3,
					NewSellPrice1:    *rowOtGroup.NewSellPrice1,
					NewSellPrice2:    *rowOtGroup.NewSellPrice2,
					NewSellPrice3:    *rowOtGroup.NewSellPrice3,
				})
			}
		}

		// prioritized based on specific outlet
		for _, rowOutlet := range params.OutletsReq {
			refID := int64(*rowOutlet.RefID)
			outletID := int64(row.OutletId)

			// validate ref_id first
			_, err := service.OutletRepository.FindOneByOutletIdAndCustId(refID, params.CustID, params.ParentCustID)
			if err != nil {
				outletIDStr := strconv.Itoa(int(refID))
				return output, errors.New("Outlet ID: " + outletIDStr + " not found")
			}

			if refID == outletID {
				for i, rowResp := range output {
					if rowResp.OutletID == row.OutletId {
						// Remove the element at index i
						output = append(output[:i], output[i+1:]...)
						break // Exit the inner loop after deletion
					}
				}
				output = append(output, entity.OutputSprice{
					OutletID:         row.OutletId,
					OutletCode:       *row.OutletCode,
					OutletName:       *row.OutletName,
					MasterSellPrice1: params.MasterSellPrice1,
					MasterSellPrice2: params.MasterSellPrice2,
					MasterSellPrice3: params.MasterSellPrice3,
					NewSellPrice1:    *rowOutlet.NewSellPrice1,
					NewSellPrice2:    *rowOutlet.NewSellPrice2,
					NewSellPrice3:    *rowOutlet.NewSellPrice3,
				})
			}
		}
	}
	return output, err
}

func (service *spServiceImpl) Preview(request entity.CreateMSpPriceBody) (response entity.PreviewMSpPriceResp, err error) {
	//init the loc
	asiaJkt, _ := time.LoadLocation("Asia/Jakarta")
	startDateStr := request.StartDate
	startDate, errParseDate := time.ParseInLocation(constant.DATE_LAYOUT_YYYY_MM_DD, startDateStr, asiaJkt)
	if errParseDate != nil {
		log.Error("Error parsing startDate:", errParseDate.Error())
		return response, errParseDate
	}

	endDateStr := request.EndDate
	endDate, errParseDate := time.ParseInLocation(constant.DATE_LAYOUT_YYYY_MM_DD, endDateStr, asiaJkt)
	if errParseDate != nil {
		log.Error("Error parsing endDate:", errParseDate.Error())
		return response, errParseDate
	}

	timeNowJkt := time.Now().In(asiaJkt)
	currentDate := time.Date(timeNowJkt.Year(), timeNowJkt.Month(), timeNowJkt.Day(), 0, 0, 0, 0, timeNowJkt.Location())
	tomorrow := currentDate.AddDate(0, 0, 1)

	if startDate.Before(tomorrow) { // rules: current date + 1, validate
		return response, errors.New(`start_date ` + errmsg.ERROR_DATE_MUST_GT_NOW)
	}
	if endDate.Before(tomorrow) { // rules: current date + 1, validate
		return response, errors.New(`end_date ` + errmsg.ERROR_DATE_MUST_GT_NOW)
	}

	// get special price group
	specialPriceGroup, err := service.SpecialPriceGroupRepository.FindOneBySpecialPriceGroupIdAndCustId(int(request.PriceGrpID), request.CustID)
	if err != nil {
		return response, errors.New("special price group id not found")
	}

	// get product
	product, err := service.SpRepository.FindOneProductByProID(request.ProID, request.ParentCustID)
	if err != nil {
		return response, errors.New("pro_id not found")
	}

	// get outlets by price_grp_id
	outlets, err := service.OutletRepository.FindAllByPriceGrpIDAndCustID(int(specialPriceGroup.SpecialPriceGroupId), request.CustID, request.ParentCustID)
	if err != nil {
		return response, errors.New("find outlet by price group not found")
	}

	// find details by ref_id
	outletTypesIDs := []int64{}
	for _, row := range request.Details.OutletType {
		outletTypesIDs = append(outletTypesIDs, *row.RefID)
	}

	outletGroupsIDs := []int64{}
	for _, row := range request.Details.OutletGroup {
		outletGroupsIDs = append(outletGroupsIDs, *row.RefID)
	}

	outletsIDs := []int64{}
	for _, row := range request.Details.Outlet {
		outletsIDs = append(outletsIDs, *row.RefID)
	}

	err = structs.Automapper(product, &response)
	if err != nil {
		return response, err
	}

	response.StartDate = request.StartDate
	response.EndDate = request.EndDate
	response.PriceGrpID = int64(specialPriceGroup.SpecialPriceGroupId)
	response.PriceGrpCode = specialPriceGroup.SpecialPriceGroupCode
	response.PriceGrpName = specialPriceGroup.SpecialPriceGroupName
	response.PurchPrice1 = product.PurchPrice1
	response.PurchPrice2 = product.PurchPrice2
	response.PurchPrice3 = product.PurchPrice3
	response.MasterSellPrice1 = product.SellPrice1
	response.MasterSellPrice2 = product.SellPrice2
	response.MasterSellPrice3 = product.SellPrice3
	response.NewSellPrice1 = request.NewSellPrice1
	response.NewSellPrice2 = request.NewSellPrice2
	response.NewSellPrice3 = request.NewSellPrice3
	response.Details = request.Details
	response.Status = 1
	statusResp := entity.MSpPriceResponse{
		Status: response.Status,
	}
	response.StatusDesc = statusResp.GetSpPriceStatusDesc()

	generateOutputParams := GenerateOutputParams{
		Outlets:          outlets,
		OutletTypesReq:   request.Details.OutletType,
		OutletGroupsReq:  request.Details.OutletGroup,
		OutletsReq:       request.Details.Outlet,
		MasterSellPrice1: product.SellPrice1,
		MasterSellPrice2: product.SellPrice2,
		MasterSellPrice3: product.SellPrice3,
		NewSellPrice1:    request.NewSellPrice1,
		NewSellPrice2:    request.NewSellPrice2,
		NewSellPrice3:    request.NewSellPrice3,
		CustID:           request.CustID,
		ParentCustID:     request.ParentCustID,
	}
	response.OutputSprice, err = service.GetOutputByPriceGroupIdOutletTypeOutletGroupOutlet(generateOutputParams)
	if err != nil {
		log.Error("generate output:", err.Error())
		return response, err
	}

	log.Info("response, Preview:", structs.StructToJson(response))
	return response, nil
}

func (service *spServiceImpl) Store(request entity.CreateMSpPriceBody) (response entity.MSpPriceResponse, err error) {
	//init the loc
	asiaJkt, _ := time.LoadLocation("Asia/Jakarta")
	startDateStr := request.StartDate
	startDate, errParseDate := time.ParseInLocation(constant.DATE_LAYOUT_YYYY_MM_DD, startDateStr, asiaJkt)
	if errParseDate != nil {
		log.Error("Error parsing startDate:", errParseDate.Error())
		return response, errParseDate
	}
	request.StartDate, err = str.DateStrToRfc3339String(request.StartDate)
	if err != nil {
		return response, err
	}

	endDateStr := request.EndDate
	endDate, errParseDate := time.ParseInLocation(constant.DATE_LAYOUT_YYYY_MM_DD, endDateStr, asiaJkt)
	if errParseDate != nil {
		log.Error("Error parsing endDate:", errParseDate.Error())
		return response, errParseDate
	}
	request.EndDate, err = str.DateStrToRfc3339String(request.EndDate)
	if err != nil {
		return response, err
	}

	timeNowJkt := time.Now().In(asiaJkt)
	currentDate := time.Date(timeNowJkt.Year(), timeNowJkt.Month(), timeNowJkt.Day(), 0, 0, 0, 0, timeNowJkt.Location())
	tomorrow := currentDate.AddDate(0, 0, 1)

	if startDate.Before(tomorrow) { // rules: current date + 1, validate
		return response, errors.New(`start_date ` + errmsg.ERROR_DATE_MUST_GT_NOW)
	}
	if endDate.Before(tomorrow) { // rules: current date + 1, validate
		return response, errors.New(`end_date ` + errmsg.ERROR_DATE_MUST_GT_NOW)
	}

	// get special price group
	specialPriceGroup, err := service.SpecialPriceGroupRepository.FindOneBySpecialPriceGroupIdAndCustId(int(request.PriceGrpID), request.CustID)
	if err != nil {
		log.Error("special price group id not found")
		return response, errors.New("special price group id not found")
	}

	// get product
	product, err := service.SpRepository.FindOneProductByProID(request.ProID, request.ParentCustID)
	if err != nil {
		log.Error("pro_id not found")
		return response, errors.New("pro_id not found")
	}
	request.UnitId1 = product.UnitId1
	request.UnitId2 = product.UnitId2
	request.UnitId3 = product.UnitId3
	request.ConvUnit2 = product.ConvUnit2
	request.ConvUnit3 = product.ConvUnit3
	request.SellPrice1 = product.SellPrice1
	request.SellPrice2 = product.SellPrice2
	request.SellPrice3 = product.SellPrice3
	request.Status = 1

	// get outlets by price_grp_id
	outlets, err := service.OutletRepository.FindAllByPriceGrpIDAndCustID(int(specialPriceGroup.SpecialPriceGroupId), request.CustID, request.ParentCustID)
	if err != nil {
		log.Error("find outlet by price group not found")
		return response, errors.New("find outlet by price group not found")
	}

	// find details by ref_id
	outletTypesIDs := []int64{}
	for _, row := range request.Details.OutletType {
		outletTypesIDs = append(outletTypesIDs, *row.RefID)
	}

	outletGroupsIDs := []int64{}
	for _, row := range request.Details.OutletGroup {
		outletGroupsIDs = append(outletGroupsIDs, *row.RefID)
	}

	outletsIDs := []int64{}
	for _, row := range request.Details.Outlet {
		outletsIDs = append(outletsIDs, *row.RefID)
	}

	err = structs.Automapper(product, &response)
	if err != nil {
		return response, err
	}

	generateOutputParams := GenerateOutputParams{
		Outlets:          outlets,
		OutletTypesReq:   request.Details.OutletType,
		OutletGroupsReq:  request.Details.OutletGroup,
		OutletsReq:       request.Details.Outlet,
		MasterSellPrice1: product.SellPrice1,
		MasterSellPrice2: product.SellPrice2,
		MasterSellPrice3: product.SellPrice3,
		NewSellPrice1:    request.NewSellPrice1,
		NewSellPrice2:    request.NewSellPrice2,
		NewSellPrice3:    request.NewSellPrice3,
		CustID:           request.CustID,
		ParentCustID:     request.ParentCustID,
	}
	output, err := service.GetOutputByPriceGroupIdOutletTypeOutletGroupOutlet(generateOutputParams)
	if err != nil {
		log.Error("generate output:", err.Error())
		return response, err
	}
	// log.Info("output:", structs.StructToJson(output))

	var mSpriceData model.MSpPrice
	err = structs.Automapper(request, &mSpriceData)
	if err != nil {
		log.Error("error automapper mSpriceData")
		return response, err
	}
	timeNow := time.Now()
	mSpriceData.CreatedAt = &timeNow
	mSpriceData.CreatedBy = request.CreatedBy
	mSpriceData.UpdatedAt = &timeNow
	mSpriceData.UpdatedBy = request.CreatedBy
	objectID := primitive.NewObjectID() // Generate a new ObjectID
	objectIDString := objectID.Hex()    // Convert ObjectID to string
	mSpriceData.SpPriceID = objectIDString
	// log.Info("mSpriceData:", structs.StructToJson(mSpriceData))

	trx, err := service.SpRepository.TrxBegin()
	if err != nil {
		return response, err
	}

	err = trx.InsertPrice(&mSpriceData)
	if err != nil {
		trx.TrxRollback()
		return response, err
	}

	lastNewSellPrice1 := request.NewSellPrice1
	lastNewSellPrice2 := request.NewSellPrice2
	lastNewSellPrice3 := request.NewSellPrice3

	// insert sp_price_det ref_type -> outlet type
	for _, detail := range request.Details.OutletType {
		var MspriceDetData model.MSpPriceDet
		err = structs.Automapper(detail, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		objectDetID := primitive.NewObjectID() // Generate a new ObjectID
		objectDetIDString := objectDetID.Hex() // Convert ObjectID to string
		MspriceDetData.SpPriceDetID = objectDetIDString
		MspriceDetData.SpPriceID = objectIDString
		MspriceDetData.CustID = request.CustID
		MspriceDetData.CreatedAt = &timeNow
		MspriceDetData.CreatedBy = request.CreatedBy
		MspriceDetData.UpdatedAt = &timeNow
		MspriceDetData.UpdatedBy = request.CreatedBy
		MspriceDetData.RefType = OUTLETTYPE
		MspriceDetData.SellPrice1 = lastNewSellPrice1
		MspriceDetData.SellPrice2 = lastNewSellPrice2
		MspriceDetData.SellPrice3 = lastNewSellPrice3

		err := trx.InsertPriceDetail(mSpriceData.SpPriceID, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}

		lastNewSellPrice1 = MspriceDetData.NewSellPrice1
		lastNewSellPrice2 = MspriceDetData.NewSellPrice2
		lastNewSellPrice3 = MspriceDetData.NewSellPrice3
	}

	// insert sp_price_det ref_type -> outlet group
	for _, detail := range request.Details.OutletGroup {
		var MspriceDetData model.MSpPriceDet
		err = structs.Automapper(detail, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		objectDetID := primitive.NewObjectID() // Generate a new ObjectID
		objectDetIDString := objectDetID.Hex() // Convert ObjectID to string
		MspriceDetData.SpPriceDetID = objectDetIDString
		MspriceDetData.SpPriceID = objectIDString
		MspriceDetData.CustID = request.CustID
		MspriceDetData.CreatedAt = &timeNow
		MspriceDetData.CreatedBy = request.CreatedBy
		MspriceDetData.UpdatedAt = &timeNow
		MspriceDetData.UpdatedBy = request.CreatedBy
		MspriceDetData.RefType = OUTLETGROUP
		MspriceDetData.SellPrice1 = lastNewSellPrice1
		MspriceDetData.SellPrice2 = lastNewSellPrice2
		MspriceDetData.SellPrice3 = lastNewSellPrice3

		err := trx.InsertPriceDetail(mSpriceData.SpPriceID, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}

		lastNewSellPrice1 = MspriceDetData.NewSellPrice1
		lastNewSellPrice2 = MspriceDetData.NewSellPrice2
		lastNewSellPrice3 = MspriceDetData.NewSellPrice3
	}

	// insert sp_price_det ref_type -> outlet
	for _, detail := range request.Details.Outlet {
		var MspriceDetData model.MSpPriceDet
		err = structs.Automapper(detail, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		objectDetID := primitive.NewObjectID() // Generate a new ObjectID
		objectDetIDString := objectDetID.Hex() // Convert ObjectID to string
		MspriceDetData.SpPriceDetID = objectDetIDString
		MspriceDetData.SpPriceID = objectIDString
		MspriceDetData.CustID = request.CustID
		MspriceDetData.CreatedAt = &timeNow
		MspriceDetData.CreatedBy = request.CreatedBy
		MspriceDetData.UpdatedAt = &timeNow
		MspriceDetData.UpdatedBy = request.CreatedBy
		MspriceDetData.RefType = OUTLET
		MspriceDetData.SellPrice1 = lastNewSellPrice1
		MspriceDetData.SellPrice2 = lastNewSellPrice2
		MspriceDetData.SellPrice3 = lastNewSellPrice3

		err := trx.InsertPriceDetail(mSpriceData.SpPriceID, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}

		lastNewSellPrice1 = MspriceDetData.NewSellPrice1
		lastNewSellPrice2 = MspriceDetData.NewSellPrice2
		lastNewSellPrice3 = MspriceDetData.NewSellPrice3
	}

	// insert sp_price_det ref_type -> output
	for _, detail := range output {
		var MspriceDetData model.MSpPriceDet
		err = structs.Automapper(detail, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		objectDetID := primitive.NewObjectID() // Generate a new ObjectID
		objectDetIDString := objectDetID.Hex() // Convert ObjectID to string
		MspriceDetData.SpPriceDetID = objectDetIDString
		MspriceDetData.SpPriceID = objectIDString
		MspriceDetData.CustID = request.CustID
		MspriceDetData.CreatedAt = &timeNow
		MspriceDetData.CreatedBy = request.CreatedBy
		MspriceDetData.UpdatedAt = &timeNow
		MspriceDetData.UpdatedBy = request.CreatedBy
		MspriceDetData.RefType = OUTPUT
		MspriceDetData.SellPrice1 = lastNewSellPrice1
		MspriceDetData.SellPrice2 = lastNewSellPrice2
		MspriceDetData.SellPrice3 = lastNewSellPrice3
		MspriceDetData.RefID = detail.OutletID

		err := trx.InsertPriceDetail(mSpriceData.SpPriceID, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
	}
	trx.TrxCommit()

	startDateStr = startDateStr + "T00:00:00+07:00"
	start := time.Now().In(asiaJkt)              // Record the start time
	delta := startDate.Sub(start).Milliseconds() // Calculate the delta time in milliseconds
	if delta < 0 {
		return response, errors.New("start is not valid")
	}
	if request.ExpirationMs > 0 {
		delta = int64(request.ExpirationMs)
		testStartDateTime := start.Add(time.Duration(delta) * time.Millisecond)
		// log.Info("testDateTime:", testStartDateTime)
		startDateStr = testStartDateTime.Format(time.RFC3339)
	}

	deltaStr := strconv.Itoa(int(delta))
	log.Info("SP Service, start_date deltaStr -> ", deltaStr)

	payloadRmq := entity.PublishUnpublishSPriceReq{
		CustID:       request.CustID,
		ParentCustID: request.ParentCustID,
		SpPriceID:    mSpriceData.SpPriceID,
		Status:       10,
		UpdatedBy:    "system",
	}

	pubJobStart := entity.PublishJob{
		JobName:   "Publish price",
		JobDesc:   fmt.Sprintf("publish price, pro_id: %d", mSpriceData.ProID),
		JobType:   constant.JOB_TYPE_ONE_TIME,
		Task:      constant.JOB_TASK_HTTP_REQ,
		RunAt:     startDateStr,
		Url:       service.Config.Get("MASTER_SERVICE_URL") + constant.SP_PRICE_PUBLISH_UNPUBLISH,
		Payload:   structs.StructToJson(payloadRmq),
		CreatedBy: "system",
	}
	go service.publishJob(pubJobStart)

	endDelta := endDate.Sub(start).Milliseconds() // Calculate the delta time in milliseconds
	if endDelta < 0 {
		return response, errors.New("end date is not valid")
	}
	endDateStr = endDateStr + "T23:59:59+07:00"
	if request.EndExpirationMs > 0 {
		endDelta = int64(request.EndExpirationMs)
		testEndDateTime := start.Add(time.Duration(endDelta) * time.Millisecond)
		// log.Info("testEndDateTime:", testEndDateTime)
		endDateStr = testEndDateTime.Format(time.RFC3339)
	}

	endDeltaStr := strconv.Itoa(int(endDelta))
	log.Info("SP Service, end_date endDeltaStr -> ", endDeltaStr)

	payloadRmq.Status = 7
	pubJobEnd := entity.PublishJob{
		JobName:   "Unpublish price",
		JobDesc:   fmt.Sprintf("unpublish price, pro_id: %d", mSpriceData.ProID),
		JobType:   constant.JOB_TYPE_ONE_TIME,
		Task:      constant.JOB_TASK_HTTP_REQ,
		RunAt:     endDateStr,
		Url:       service.Config.Get("MASTER_SERVICE_URL") + constant.SP_PRICE_PUBLISH_UNPUBLISH,
		Payload:   structs.StructToJson(payloadRmq),
		CreatedBy: "system",
	}
	go service.publishJob(pubJobEnd)

	return response, nil
}

func (service *spServiceImpl) Detail(params entity.MSpPriceParams) (response entity.MSpPriceWithDetailResp, err error) {
	spPrice, err := service.SpRepository.FindOneBySpPriceIdAndCustID(params)
	if err != nil {
		return response, err
	}
	err = structs.Automapper(spPrice, &response)
	if err != nil {
		return response, err
	}
	startDate, endDate := "", ""
	if spPrice.StartDate != nil {
		startDate = spPrice.StartDate.Format("2006-01-02")
		response.StartDate = &startDate
	}
	if spPrice.EndDate != nil {
		endDate = spPrice.EndDate.Format("2006-01-02")
		response.EndDate = &endDate
	}

	statusResp := entity.MSpPriceResponse{
		Status: response.Status,
	}
	response.StatusDesc = statusResp.GetSpPriceStatusDesc()

	spPriceDetails, err := service.SpRepository.FindDetailSpPriceIdAndCustId(params)
	if err != nil {
		return response, err
	}
	response.Details.SalesTeam = make([]entity.MSpPriceDetResp, 0)
	response.Details.OutletType = make([]entity.MSpPriceDetResp, 0)
	response.Details.OutletGroup = make([]entity.MSpPriceDetResp, 0)
	response.Details.Outlet = make([]entity.MSpPriceDetResp, 0)
	for _, spPriceDetail := range spPriceDetails {
		var spPriceDetailResp entity.MSpPriceDetResp
		err = structs.Automapper(spPriceDetail, &spPriceDetailResp)
		if err != nil {
			return response, err
		}
		if spPriceDetailResp.RefType == SALESTEAM {
			response.Details.SalesTeam = append(response.Details.SalesTeam, spPriceDetailResp)
		}
		if spPriceDetailResp.RefType == OUTLETTYPE {
			response.Details.OutletType = append(response.Details.OutletType, spPriceDetailResp)
		}
		if spPriceDetailResp.RefType == OUTLETGROUP {
			response.Details.OutletGroup = append(response.Details.OutletGroup, spPriceDetailResp)
		}
		if spPriceDetailResp.RefType == OUTLET {
			response.Details.Outlet = append(response.Details.Outlet, spPriceDetailResp)
		}
		if spPriceDetailResp.RefType == OUTPUT {
			var spPriceOutputResp entity.OutputSprice
			err = structs.Automapper(spPriceDetailResp, &spPriceOutputResp)
			if err != nil {
				return response, err
			}
			spPriceOutputResp.MasterSellPrice1 = *spPriceDetail.SellPrice1
			spPriceOutputResp.MasterSellPrice2 = *spPriceDetail.SellPrice2
			spPriceOutputResp.MasterSellPrice3 = *spPriceDetail.SellPrice3
			spPriceOutputResp.OutletID = *spPriceDetail.RefID
			spPriceOutputResp.OutletCode = *spPriceDetail.RefCode
			spPriceOutputResp.OutletName = *spPriceDetail.RefName
			response.OutputSprice = append(response.OutputSprice, spPriceOutputResp)
		}
	}

	return response, nil
}

func (service *spServiceImpl) List(dataFilter entity.MSpPriceQueryFilter, custId string) (data []entity.MSpPriceResponse, total int, lastPage int, err error) {
	spPrices, total, lastPage, err := service.SpRepository.FindAllByCustID(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}
	for _, row := range spPrices {
		var vResp entity.MSpPriceResponse
		structs.Automapper(row, &vResp)
		startDate, endDate := "", ""
		if row.StartDate != nil {
			startDate = row.StartDate.Format("2006-01-02")
			vResp.StartDate = &startDate
		}
		if row.EndDate != nil {
			endDate = row.EndDate.Format("2006-01-02")
			vResp.EndDate = &endDate
		}
		statusResp := entity.MSpPriceResponse{
			Status: row.Status,
		}
		vResp.StatusDesc = statusResp.GetSpPriceStatusDesc()
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *spServiceImpl) Delete(custId string, sppriceId string) (err error) {

	trx, err := service.SpRepository.TrxBegin()
	if err != nil {
		return err
	}

	err = trx.Delete(custId, sppriceId)
	if err != nil {
		trx.TrxRollback()
		return err
	}
	trx.TrxCommit()

	return err
}

func (service *spServiceImpl) Update(sppriceId string, request entity.UpdateMSpPriceBody) (err error) {
	params := entity.MSpPriceParams{
		CustID:       request.CustID,
		ParentCustID: request.ParentCustID,
		SpPriceID:    sppriceId,
	}
	spPrice, err := service.SpRepository.FindOneBySpPriceIdAndCustID(params)
	if err != nil {
		return err
	}
	if spPrice.Status != 1 {
		return errors.New("status not allowed to update")
	}

	//init the loc
	asiaJkt, _ := time.LoadLocation("Asia/Jakarta")
	startDateStr := request.StartDate
	startDate, errParseDate := time.ParseInLocation(constant.DATE_LAYOUT_YYYY_MM_DD, startDateStr, asiaJkt)
	if errParseDate != nil {
		log.Error("Error parsing startDate:", errParseDate.Error())
		return errParseDate
	}
	request.StartDate, err = str.DateStrToRfc3339String(request.StartDate)
	if err != nil {
		return err
	}

	endDateStr := request.EndDate
	endDate, errParseDate := time.ParseInLocation(constant.DATE_LAYOUT_YYYY_MM_DD, endDateStr, asiaJkt)
	if errParseDate != nil {
		log.Error("Error parsing endDate:", errParseDate.Error())
		return errParseDate
	}
	request.EndDate, err = str.DateStrToRfc3339String(request.EndDate)
	if err != nil {
		return err
	}

	timeNowJkt := time.Now().In(asiaJkt)
	currentDate := time.Date(timeNowJkt.Year(), timeNowJkt.Month(), timeNowJkt.Day(), 0, 0, 0, 0, timeNowJkt.Location())
	tomorrow := currentDate.AddDate(0, 0, 1)

	if startDate.Before(tomorrow) { // rules: current date + 1, validate
		return errors.New(`start_date ` + errmsg.ERROR_DATE_MUST_GT_NOW)
	}
	if endDate.Before(tomorrow) { // rules: current date + 1, validate
		return errors.New(`end_date ` + errmsg.ERROR_DATE_MUST_GT_NOW)
	}

	// get special price group
	specialPriceGroup, err := service.SpecialPriceGroupRepository.FindOneBySpecialPriceGroupIdAndCustId(int(request.PriceGrpID), request.CustID)
	if err != nil {
		log.Error("special price group id not found")
		return errors.New("special price group id not found")
	}

	// get product
	product, err := service.SpRepository.FindOneProductByProID(request.ProID, request.ParentCustID)
	if err != nil {
		log.Error("pro_id not found")
		return errors.New("pro_id not found")
	}
	request.UnitId1 = product.UnitId1
	request.UnitId2 = product.UnitId2
	request.UnitId3 = product.UnitId3
	request.ConvUnit2 = product.ConvUnit2
	request.ConvUnit3 = product.ConvUnit3
	request.SellPrice1 = product.SellPrice1
	request.SellPrice2 = product.SellPrice2
	request.SellPrice3 = product.SellPrice3
	request.Status = 1

	// get outlets by price_grp_id
	outlets, err := service.OutletRepository.FindAllByPriceGrpIDAndCustID(int(specialPriceGroup.SpecialPriceGroupId), request.CustID, request.ParentCustID)
	if err != nil {
		log.Error("find outlet by price group not found")
		return errors.New("find outlet by price group not found")
	}

	// find details by ref_id
	outletTypesIDs := []int64{}
	for _, row := range request.Details.OutletType {
		outletTypesIDs = append(outletTypesIDs, *row.RefID)
	}

	outletGroupsIDs := []int64{}
	for _, row := range request.Details.OutletGroup {
		outletGroupsIDs = append(outletGroupsIDs, *row.RefID)
	}

	outletsIDs := []int64{}
	for _, row := range request.Details.Outlet {
		outletsIDs = append(outletsIDs, *row.RefID)
	}

	generateOutputParams := GenerateOutputParams{
		Outlets:          outlets,
		OutletTypesReq:   request.Details.OutletType,
		OutletGroupsReq:  request.Details.OutletGroup,
		OutletsReq:       request.Details.Outlet,
		MasterSellPrice1: product.SellPrice1,
		MasterSellPrice2: product.SellPrice2,
		MasterSellPrice3: product.SellPrice3,
		NewSellPrice1:    request.NewSellPrice1,
		NewSellPrice2:    request.NewSellPrice2,
		NewSellPrice3:    request.NewSellPrice3,
		CustID:           request.CustID,
		ParentCustID:     request.ParentCustID,
	}
	output, err := service.GetOutputByPriceGroupIdOutletTypeOutletGroupOutlet(generateOutputParams)
	if err != nil {
		log.Error("generate output:", err.Error())
		return err
	}
	log.Info("output:", structs.StructToJson(output))

	timeNow := time.Now().UTC()
	request.UpdatedAt = &timeNow
	trx, err := service.SpRepository.TrxBegin()
	if err != nil {
		return err
	}

	// delete old data, all details
	err = trx.DeleteDetails(params.CustID, sppriceId)
	if err != nil {
		trx.TrxRollback()
		return err
	}

	log.Info("update, request:", structs.StructToJson(request))
	err = trx.Update(sppriceId, request)
	if err != nil {
		trx.TrxRollback()
		return err
	}

	lastNewSellPrice1 := request.NewSellPrice1
	lastNewSellPrice2 := request.NewSellPrice2
	lastNewSellPrice3 := request.NewSellPrice3

	// insert sp_price_det ref_type -> outlet type
	for _, detail := range request.Details.OutletType {
		var MspriceDetData model.MSpPriceDet
		err = structs.Automapper(detail, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return err
		}
		objectDetID := primitive.NewObjectID() // Generate a new ObjectID
		objectDetIDString := objectDetID.Hex() // Convert ObjectID to string
		MspriceDetData.SpPriceDetID = objectDetIDString
		MspriceDetData.SpPriceID = sppriceId
		MspriceDetData.CustID = request.CustID
		MspriceDetData.CreatedAt = &timeNow
		MspriceDetData.CreatedBy = request.UpdatedBy
		MspriceDetData.UpdatedAt = &timeNow
		MspriceDetData.UpdatedBy = request.UpdatedBy
		MspriceDetData.RefType = OUTLETTYPE
		MspriceDetData.SellPrice1 = lastNewSellPrice1
		MspriceDetData.SellPrice2 = lastNewSellPrice2
		MspriceDetData.SellPrice3 = lastNewSellPrice3

		err := trx.InsertPriceDetail(sppriceId, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return err
		}

		lastNewSellPrice1 = MspriceDetData.NewSellPrice1
		lastNewSellPrice2 = MspriceDetData.NewSellPrice2
		lastNewSellPrice3 = MspriceDetData.NewSellPrice3
	}

	// insert sp_price_det ref_type -> outlet group
	for _, detail := range request.Details.OutletGroup {
		var MspriceDetData model.MSpPriceDet
		err = structs.Automapper(detail, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return err
		}
		objectDetID := primitive.NewObjectID() // Generate a new ObjectID
		objectDetIDString := objectDetID.Hex() // Convert ObjectID to string
		MspriceDetData.SpPriceDetID = objectDetIDString
		MspriceDetData.SpPriceID = sppriceId
		MspriceDetData.CustID = request.CustID
		MspriceDetData.CreatedAt = &timeNow
		MspriceDetData.CreatedBy = request.UpdatedBy
		MspriceDetData.UpdatedAt = &timeNow
		MspriceDetData.UpdatedBy = request.UpdatedBy
		MspriceDetData.RefType = OUTLETGROUP
		MspriceDetData.SellPrice1 = lastNewSellPrice1
		MspriceDetData.SellPrice2 = lastNewSellPrice2
		MspriceDetData.SellPrice3 = lastNewSellPrice3

		err := trx.InsertPriceDetail(sppriceId, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return err
		}

		lastNewSellPrice1 = MspriceDetData.NewSellPrice1
		lastNewSellPrice2 = MspriceDetData.NewSellPrice2
		lastNewSellPrice3 = MspriceDetData.NewSellPrice3
	}

	// insert sp_price_det ref_type -> outlet
	for _, detail := range request.Details.Outlet {
		var MspriceDetData model.MSpPriceDet
		err = structs.Automapper(detail, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return err
		}
		objectDetID := primitive.NewObjectID() // Generate a new ObjectID
		objectDetIDString := objectDetID.Hex() // Convert ObjectID to string
		MspriceDetData.SpPriceDetID = objectDetIDString
		MspriceDetData.SpPriceID = sppriceId
		MspriceDetData.CustID = request.CustID
		MspriceDetData.CreatedAt = &timeNow
		MspriceDetData.CreatedBy = request.UpdatedBy
		MspriceDetData.UpdatedAt = &timeNow
		MspriceDetData.UpdatedBy = request.UpdatedBy
		MspriceDetData.RefType = OUTLET
		MspriceDetData.SellPrice1 = lastNewSellPrice1
		MspriceDetData.SellPrice2 = lastNewSellPrice2
		MspriceDetData.SellPrice3 = lastNewSellPrice3

		err := trx.InsertPriceDetail(sppriceId, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return err
		}

		lastNewSellPrice1 = MspriceDetData.NewSellPrice1
		lastNewSellPrice2 = MspriceDetData.NewSellPrice2
		lastNewSellPrice3 = MspriceDetData.NewSellPrice3
	}

	// insert sp_price_det ref_type -> output
	for _, detail := range output {
		var MspriceDetData model.MSpPriceDet
		err = structs.Automapper(detail, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return err
		}
		objectDetID := primitive.NewObjectID() // Generate a new ObjectID
		objectDetIDString := objectDetID.Hex() // Convert ObjectID to string
		MspriceDetData.SpPriceDetID = objectDetIDString
		MspriceDetData.SpPriceID = sppriceId
		MspriceDetData.CustID = request.CustID
		MspriceDetData.CreatedAt = &timeNow
		MspriceDetData.CreatedBy = request.UpdatedBy
		MspriceDetData.UpdatedAt = &timeNow
		MspriceDetData.UpdatedBy = request.UpdatedBy
		MspriceDetData.RefType = OUTPUT
		MspriceDetData.SellPrice1 = lastNewSellPrice1
		MspriceDetData.SellPrice2 = lastNewSellPrice2
		MspriceDetData.SellPrice3 = lastNewSellPrice3
		MspriceDetData.RefID = detail.OutletID

		err := trx.InsertPriceDetail(sppriceId, &MspriceDetData)
		if err != nil {
			trx.TrxRollback()
			return err
		}
	}

	trx.TrxCommit()
	return nil
}

func (service *spServiceImpl) Cancel(params entity.MSpPriceCancelParams) (err error) {
	log.Info("masuk cancel")
	mParams := entity.MSpPriceParams{
		CustID:       params.CustID,
		ParentCustID: params.ParentCustID,
		SpPriceID:    params.SpPriceID,
	}
	log.Infof("mParams: %+v", mParams)
	spPrice, err := service.SpRepository.FindOneBySpPriceIdAndCustID(mParams)
	if err != nil {
		log.Error("error:", err.Error())
		return err
	}
	if spPrice.Status != 1 {
		log.Error("status not allowed to cancel")
		return errors.New("status not allowed to cancel")
	}

	request := entity.UpdateMSpPriceBody{
		ParentCustID: params.ParentCustID,
		CustID:       params.CustID,
		SpPriceID:    params.SpPriceID,
		Status:       5,
		UpdatedBy:    params.UpdatedBy,
	}

	log.Info("cancel, request:", structs.StructToJson(request))

	trx, err := service.SpRepository.TrxBegin()
	if err != nil {
		return err
	}

	err = trx.Update(params.SpPriceID, request)
	if err != nil {
		trx.TrxRollback()
		return err
	}
	trx.TrxCommit()

	return nil
}

func SetupTransactionPriceData(price model.MSpPriceDetPublish, distributorID int64) (transPrice model.MTransactionPrice) {
	objectID := primitive.NewObjectID() // Generate a new ObjectID
	objectIDString := objectID.Hex()    // Convert ObjectID to string
	transPrice = model.MTransactionPrice{
		CustID:             price.CustID,
		TransactionPriceID: objectIDString,
		ProID:              price.ProID,
		PurchPrice1:        0,
		PurchPrice2:        0,
		PurchPrice3:        0,
		SellPrice1:         *price.NewSellPrice1,
		SellPrice2:         *price.NewSellPrice2,
		SellPrice3:         *price.NewSellPrice3,
		Source:             5, // selling price = 5
		CreatedBy:          price.CreatedBy,
		CreatedAt:          time.Now().UTC(),
		StartDate:          price.StartDate,
		EndDate:            price.EndDate,
		Coverage:           "N",
		DistributorID:      distributorID,
		ReferenceID:        price.SpPriceID,
		OutletID:           *price.RefID,
	}

	return transPrice
}

func (service *spServiceImpl) PublishOrInactive(request entity.PublishUnpublishSPriceReq) (err error) {
	log.Info("spServiceImpl, request -> ", structs.StructToJson(request))
	mParams := entity.MSpPriceParams{
		CustID:       request.CustID,
		ParentCustID: request.ParentCustID,
		SpPriceID:    request.SpPriceID,
	}

	spPrice, err := service.SpRepository.FindOneBySpPriceIdAndCustID(mParams)
	if err != nil {
		log.Error("error:", err.Error())
		return err
	}

	if request.Status == 10 && spPrice.Status != 1 {
		log.Error("status not allowed to publish")
		return errors.New("status not allowed to publish")
	}

	spPriceDetails, err := service.SpRepository.FindDetailSpPriceIdAndCustIdPublish(mParams)
	if err != nil {
		return err
	}

	// scheduler event publish
	if request.Status == 10 {
		// setup insert transaction price data
		for _, row := range spPriceDetails {
			if row.RefType == 20 {
				transPrice := SetupTransactionPriceData(row, 0)
				log.Info("transPrice:", transPrice)
				if err = service.MTransactionPriceRepository.Store(&transPrice); err != nil {
					return err
				}
			}
		}
	}

	if err = service.SpRepository.UpdateStatusByRMQ(request); err != nil {
		return err
	}

	return nil
}

func (service *spServiceImpl) publishJob(job entity.PublishJob) {
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
