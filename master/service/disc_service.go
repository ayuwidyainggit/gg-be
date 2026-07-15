package service

import (
	"errors"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/str"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type DiscService interface {
	Detail(int64, string) (entity.DiscResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.DiscListResponse, total int, lastPage int, err error)
	Store(entity.CreateDiscBody) (entity.DiscResponse, error)
	Update(int64, entity.UpdateDiscRequest) error
	Delete(string, int, int64) error
}

func NewDiscService(discRepository repository.DiscRepository) *discServiceImpl {
	return &discServiceImpl{
		DiscRepository: discRepository,
	}
}

type discServiceImpl struct {
	DiscRepository repository.DiscRepository
	// MProductRepository repository.MProductRepository
}

func (service *discServiceImpl) Detail(discId int64, custId string) (response entity.DiscResponse, err error) {
	disc, err := service.DiscRepository.FindOneByDiscIdAndCustId(discId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(disc, &response)
	if err != nil {
		return response, err
	}

	discDetails, err := service.DiscRepository.FindDetailByDiscIdAndCustId(discId, custId)
	if err != nil {
		return response, err
	}
	for _, discDetail := range discDetails {
		var discDetailResp entity.DiscDet
		err = structs.Automapper(discDetail, &discDetailResp)
		if err != nil {
			return response, err
		}
		response.Details = append(response.Details, discDetailResp)
	}

	return response, err
}

func (service *discServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.DiscListResponse, total int, lastPage int, err error) {
	discs, total, lastPage, err := service.DiscRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(discs) > 0 {
		for _, row := range discs {
			var vResp entity.DiscListResponse
			structs.Automapper(row, &vResp)
			if row.StartDate != nil {
				StartDate := row.StartDate.Format("2006-01-02")
				vResp.StartDate = &StartDate
			}
			if row.EndDate != nil {
				EndDate := row.EndDate.Format("2006-01-02")
				vResp.EndDate = &EndDate
			}

			data = append(data, vResp)
		}
	}

	// discsPrint, _ := json.Marshal(discs)
	// log.Println("### DiscService, List, discsPrint ###")
	// log.Println(string(discsPrint))
	// log.Println("### End Of discsPrint ###")

	return data, total, lastPage, err
}

func (service *discServiceImpl) Store(request entity.CreateDiscBody) (response entity.DiscResponse, err error) {

	// disc_code & cust id validation, if err == nil, this means that code & cust id already exists
	disc, err := service.DiscRepository.FindOneByDiscCodeAndCustId(request.DiscCode, request.CustId)
	if err == nil {
		return response, errors.New("disc_code: " + disc.DiscCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)

	startDate, err := str.DateStrToRfc3339String(request.StartDate)
	if err != nil {
		return response, err
	}
	request.StartDate = startDate

	endDate, err := str.DateStrToRfc3339String(request.EndDate)
	if err != nil {
		return response, err
	}
	request.EndDate = endDate

	var discData model.Disc
	structs.Automapper(request, &discData)

	discData.CreatedAt = &timeNow
	discData.CreatedBy = &request.CreatedBy
	discData.UpdatedBy = &request.CreatedBy
	discData.UpdatedAt = &timeNow

	trx, err := service.DiscRepository.TrxBegin()
	if err != nil {
		trx.TrxRollback()
		return response, err
	}
	err = trx.Store(&discData)
	if err != nil {
		trx.TrxRollback()
		return response, err
	}

	response.DiscId = discData.DiscId
	for _, detail := range request.Details {
		var discDetData model.DiscDet
		err = structs.Automapper(detail, &discDetData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		discDetData.CustID = request.CustId
		discDetData.DiscID = discData.DiscId

		err := trx.StoreDetail(&discDetData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
	}
	trx.TrxCommit()
	return response, err
}

func (service *discServiceImpl) Update(discId int64, request entity.UpdateDiscRequest) (err error) {
	discDetIDs := []int64{}

	// disc_code & cust id validation, if err == nil and params discId != disc.Id, this means that code & cust id already exists
	disc, err := service.DiscRepository.FindOneByDiscCodeAndCustId(request.DiscCode, request.CustId)
	if err == nil && disc.DiscId != discId {
		return errors.New("disc_code: " + disc.DiscCode + " is already exists")
	}

	trx, err := service.DiscRepository.TrxBegin()
	if err != nil {
		return err
	}
	err = trx.Update(discId, request)
	if err != nil {
		return err
	}

	for _, detail := range request.Details {
		if detail.DiscDetID != nil {
			discDetIDs = append(discDetIDs, *detail.DiscDetID)
		}
	}

	if len(discDetIDs) > 0 {
		err := trx.DeleteDetailNotIn(discId, discDetIDs)
		if err != nil {
			trx.TrxRollback()
			return err
		}
	}

	for _, detail := range request.Details {
		if detail.DiscDetID == nil || *detail.DiscDetID == 0 {
			detail.DiscDetID = nil
			var discDetData model.DiscDet
			err = structs.Automapper(detail, &discDetData)
			if err != nil {
				trx.TrxRollback()
				return err
			}
			discDetData.CustID = request.CustId
			discDetData.DiscID = discId

			err := trx.StoreDetail(&discDetData)
			if err != nil {
				trx.TrxRollback()
				return err
			}

		} else {
			err := trx.UpdateDetail(discId, *detail.DiscDetID, detail)
			if err != nil {
				log.Println("discService, UpdateDetail, err:", err.Error())
				trx.TrxRollback()
				return err
			}
		}
	}
	trx.TrxCommit()
	return err
}

func (service *discServiceImpl) Delete(custId string, discId int, userId int64) (err error) {

	// isExists, err := service.MProductRepository.IsKeyExists(discId, custId, "disc_id1")
	// if err != nil {
	// 	return err
	// }

	// if isExists {
	// 	return errors.New("disc_id is still being used")
	// }

	err = service.DiscRepository.Delete(custId, discId, userId)
	if err != nil {
		return err
	}

	return err
}
