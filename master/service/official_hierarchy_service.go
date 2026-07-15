package service

import (
	"errors"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"strconv"
	"time"

	"golang.org/x/exp/slices"
)

type OfficialHierarchyService interface {
	Detail(int, string) (entity.OfficialHierarchyResponse, error)
	List(entity.OfficialHierarchyQueryFilter, string) (data []entity.OfficialHierarchyListResponse, total int, lastPage int, err error)
	LookupList(entity.OfficialHierarchyQueryFilter, string) (data []entity.OfficialHierarchyLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateOfficialHierarchyBody) (entity.OfficialHierarchyResponse, error)
	Upsert(entity.CreateOfficialHierarchyBody) (entity.OfficialHierarchyResponse, error)
	Update(int, entity.UpdateOfficialHierarchyRequest) error
	Delete(string, int, int64) error
	BulkUpsert(entity.BulkUpsertOfficialHierarchyBody) (entity.OfficialHierarchyResponse, error)
}

func NewOfficialHierarchyService(officialHierarchyRepository repository.OfficialHierarchyRepository) *officialHierarchyServiceImpl {
	return &officialHierarchyServiceImpl{
		OfficialHierarchyRepository: officialHierarchyRepository,
	}
}

type officialHierarchyServiceImpl struct {
	OfficialHierarchyRepository repository.OfficialHierarchyRepository
}

func (service *officialHierarchyServiceImpl) Detail(officialId int, custId string) (response entity.OfficialHierarchyResponse, err error) {
	official, err := service.OfficialHierarchyRepository.FindOneByOfficialTypeAndCustId(officialId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(official, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *officialHierarchyServiceImpl) List(dataFilter entity.OfficialHierarchyQueryFilter, custId string) (data []entity.OfficialHierarchyListResponse, total int, lastPage int, err error) {
	officials, total, lastPage, err := service.OfficialHierarchyRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range officials {
		var vResp entity.OfficialHierarchyListResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *officialHierarchyServiceImpl) LookupList(dataFilter entity.OfficialHierarchyQueryFilter, custId string) (data []entity.OfficialHierarchyLookupResponse, total int, lastPage int, err error) {
	officials, total, lastPage, err := service.OfficialHierarchyRepository.FindAllByCustIdLookup(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range officials {
		var vResp entity.OfficialHierarchyLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *officialHierarchyServiceImpl) Store(request entity.CreateOfficialHierarchyBody) (response entity.OfficialHierarchyResponse, err error) {

	timeNow := time.Now().In(time.UTC)
	var officialData model.OfficialHierarchy

	err = structs.Automapper(request, &officialData)
	if err != nil {
		return response, err
	}
	officialData.CustId = request.CustId
	officialData.CreatedAt = &timeNow
	officialData.CreatedBy = &request.CreatedBy
	officialData.UpdatedAt = &timeNow
	officialData.UpdatedBy = &request.CreatedBy
	// officialData.IsActive = true

	officialType, err := service.OfficialHierarchyRepository.Store(officialData)
	if err != nil {
		return response, err
	}

	response.OfficialType = officialType

	return response, err
}

func (service *officialHierarchyServiceImpl) Upsert(request entity.CreateOfficialHierarchyBody) (response entity.OfficialHierarchyResponse, err error) {

	// log.Println("request:", structs.StructToJson(request))
	timeNow := time.Now().In(time.UTC)
	request.CreatedAt = timeNow
	request.UpdatedBy = request.CreatedBy
	request.UpdatedAt = timeNow

	if request.OfficialType == 1 {
		boolTrue := true
		request.IsActive = &boolTrue
	}

	err = service.OfficialHierarchyRepository.Upsert(request)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *officialHierarchyServiceImpl) Update(officialType int, request entity.UpdateOfficialHierarchyRequest) (err error) {

	err = service.OfficialHierarchyRepository.Update(officialType, request)
	if err != nil {
		return err
	}

	return err
}

func (service *officialHierarchyServiceImpl) Delete(custId string, officialId int, userId int64) (err error) {

	err = service.OfficialHierarchyRepository.Delete(custId, officialId, userId)
	if err != nil {
		return err
	}

	return err
}

func (service *officialHierarchyServiceImpl) BulkUpsert(request entity.BulkUpsertOfficialHierarchyBody) (response entity.OfficialHierarchyResponse, err error) {

	log.Println("BulkUpsert - SERVICE")

	trx, err := service.OfficialHierarchyRepository.TrxBegin()
	if err != nil {
		trx.TrxRollback()
		return response, err
	}

	offTypeTemp := make([]int, 0)
	hierarchyCodeTemp := make([]string, 0)

	timeNow := time.Now().In(time.UTC)
	for _, row := range request.UpsertOfficialHierarchy {
		if !slices.Contains(offTypeTemp, row.OfficialType) && !slices.Contains(hierarchyCodeTemp, row.HierarchyCode) {
			row.CreatedAt = timeNow
			row.UpdatedAt = timeNow
			row.CreatedBy = request.CreatedBy
			row.UpdatedBy = request.CreatedBy
			row.CustId = request.CustID
			if row.OfficialType == 1 && !*row.IsActive {
				trx.TrxRollback()
				return response, errors.New("for official_type: 1, value is_active must be true")
			}
			if err = trx.UpsertWithTrx(row); err != nil {
				trx.TrxRollback()
				return response, err
			}
			offTypeTemp = append(offTypeTemp, row.OfficialType)
			hierarchyCodeTemp = append(hierarchyCodeTemp, row.HierarchyCode)
		} else {
			trx.TrxRollback()
			errMsg := "official_type: " + strconv.Itoa(row.OfficialType) + " or hierarchy_code: " + row.HierarchyCode + " is duplicate"
			return response, errors.New(errMsg)
		}
	}
	trx.TrxCommit()
	return response, err
}
