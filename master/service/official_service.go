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

type OfficialService interface {
	Detail(int, string) (entity.OfficialResponse, error)
	List(entity.OfficialQueryFilter, string) (data []entity.OfficialListResponse, total int, lastPage int, err error)
	LookupList(entity.OfficialQueryFilter, string) (data []entity.OfficialLookupResponse, total int, lastPage int, err error)
	StoreHierarchy(entity.CreateOfficialBodyHierarchy) (entity.OfficialResponse, error)
	OfficialHierarchy(entity.OfficialQueryFilter, string) (data []entity.OfficialHierarchyResp, err error)
}

func NewOfficialService(officialRepository repository.OfficialRepository) *officialServiceImpl {
	return &officialServiceImpl{
		OfficialRepository: officialRepository,
	}
}

type officialServiceImpl struct {
	OfficialRepository repository.OfficialRepository
}

func (service *officialServiceImpl) Detail(officialId int, custId string) (response entity.OfficialResponse, err error) {
	official, err := service.OfficialRepository.FindOneByOfficialIdAndCustId(officialId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(official, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *officialServiceImpl) List(dataFilter entity.OfficialQueryFilter, custId string) (data []entity.OfficialListResponse, total int, lastPage int, err error) {
	officials, total, lastPage, err := service.OfficialRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range officials {
		var vResp entity.OfficialListResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *officialServiceImpl) LookupList(dataFilter entity.OfficialQueryFilter, custId string) (data []entity.OfficialLookupResponse, total int, lastPage int, err error) {
	officials, total, lastPage, err := service.OfficialRepository.FindAllByCustIdLookup(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range officials {
		var vResp entity.OfficialLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func duplicateEmpIdValidation(officials []*entity.CreateOfficialDataHierarchy) (empId int, err error) {
	empIds := make([]int, 0)
	for _, row := range officials {
		if slices.Contains(empIds, row.EmpId) {
			errMsg := "duplicate employe id: " + strconv.Itoa(row.EmpId)
			return row.EmpId, errors.New(errMsg)
		}
		empIds = append(empIds, row.EmpId)
	}
	// log.Println("empIds:", structs.StructToJson(empIds))
	return 0, nil
}

func (service *officialServiceImpl) StoreHierarchy(request entity.CreateOfficialBodyHierarchy) (response entity.OfficialResponse, err error) {

	timeNow := time.Now().In(time.UTC)
	var officialData model.Official

	var convertedOfficials []*entity.CreateOfficialDataHierarchy
	for _, official := range request.Officials {
		convertOfficialsRekursif([]entity.CreateOfficialDataHierarchy{official}, &convertedOfficials, 0)
	}
	// log.Println("convertedOfficials:", structs.StructToJson(convertedOfficials))

	dupEmpId, err := duplicateEmpIdValidation(convertedOfficials)
	if err != nil {
		errMsg := err.Error()
		dupEmp, errFindOneEmp := service.OfficialRepository.FindOneEmployeeByEmpIdAndCustId(dupEmpId, request.CustId)
		if errFindOneEmp == nil {
			errMsg = errMsg + ", " + dupEmp.EmpCode + " | " + dupEmp.EmpName
		}
		return response, errors.New(errMsg)
	}

	trx, err := service.OfficialRepository.TrxBegin()
	if err != nil {
		return response, err
	}

	err = trx.DeleteAllByCustIdWithTrx(request.CustId)
	if err != nil {
		log.Println("DeleteAllByCustIdWithTrx, error:", err.Error())
		trx.TrxRollback()
		return response, err
	}

	for _, official := range convertedOfficials {
		err = structs.Automapper(official, &officialData)
		if err != nil {
			return response, err
		}

		officialData.CustId = request.CustId
		officialData.CreatedAt = &timeNow
		officialData.CreatedBy = &request.CreatedBy

		officialId, err := trx.StoreWithTrx(officialData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}

		response.OfficialId = officialId
	}

	trx.TrxCommit()

	return response, err
}

func convertOfficialsRekursif(officials []entity.CreateOfficialDataHierarchy, convertedOfficials *[]*entity.CreateOfficialDataHierarchy, supervisorID int) {
	for _, official := range officials {
		convertedOfficial := &entity.CreateOfficialDataHierarchy{
			OfficialType: official.OfficialType,
			EmpId:        official.EmpId,
			SupervisorId: supervisorID,
		}
		*convertedOfficials = append(*convertedOfficials, convertedOfficial)

		if len(official.Children) > 0 {
			convertOfficialsRekursif(official.Children, convertedOfficials, official.EmpId)
		}
	}
}


func (service *officialServiceImpl) OfficialHierarchy(dataFilter entity.OfficialQueryFilter, custId string) (data []entity.OfficialHierarchyResp, err error) {
	officials, err := service.OfficialRepository.HierarchyByCustId(dataFilter, custId)
	if err != nil {
		return data, err
	}

	officialsEntity := []entity.OfficialHierarchyResp{}
	officialHierarchyMap := entity.NewOfficialHierarchyMap()
	for _, row := range officials {
		var vResp entity.OfficialHierarchyResp
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, err
		}
		officialsEntity = append(officialsEntity, vResp)
		officialHierarchyMap.Db[vResp.SupervisorId] = append(officialHierarchyMap.Db[vResp.SupervisorId], vResp)
	}

	for i := range officialsEntity {
		if officialsEntity[i].SupervisorId == 0 {
			umParent := entity.OfficialHierarchyResp{}
			structs.Automapper(officialsEntity[i], &umParent)
			parent := umParent // no parents
			officialHierarchyMap.SetChildrenRecursively(&parent)
			officialHierarchyMap.Append(parent)
		}
	}
	data = officialHierarchyMap.Resp

	return data, err
}
