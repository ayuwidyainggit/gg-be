package service

import (
	"database/sql"
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/constant"
	"master/pkg/structs"
	"master/repository"
	"time"

	"github.com/lib/pq"
)

type DistributorService interface {
	Store(request entity.CreateDistributorBody) (response entity.DistributorResponse, err error)
	List(dataFilter entity.DistributorQueryFilter, custId string) (data []entity.DistributorListRespone, total int, lastPage int, err error)
	LookupList(dataFilter entity.DistributorQueryFilter, custId string) (data []entity.DistributorLookupResponse, total int, lastPage int, err error)
	Detail(params entity.DetailDistributorParams) (response entity.DistributorResponse, err error)
	Update(distributorId int64, request entity.UpdateDistributorRequest) (err error)
	Delete(custId string, distributorId int64, userId int64) (err error)
	ListWithCustomer(dataFilter entity.DistributorQueryFilter, custId string) (data []entity.DistributorCustomerResp, total int, lastPage int, err error)
}

func NewDistributorService(repository repository.DistributorRepository) *distributorServiceImpl {
	return &distributorServiceImpl{
		DistributorRepository: repository,
	}
}

type distributorServiceImpl struct {
	DistributorRepository repository.DistributorRepository
}

const distributorCodeDuplicateErrMsg = "Distributor code already exists. Please use a different distributor code."

func mapDistributorDuplicateError(err error) error {
	if err == nil {
		return nil
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return errors.New(distributorCodeDuplicateErrMsg)
	}

	return err
}

func IsDistributorNotFoundError(err error) bool {
	return errors.Is(err, constant.ErrNoRowsAffected) || errors.Is(err, sql.ErrNoRows)
}

func (service *distributorServiceImpl) resolveUpdateCustID(distributorId int64, custId string) (string, error) {
	params := entity.DetailDistributorParams{
		CustId:        custId,
		ParentCustId:  custId,
		DistributorId: int(distributorId),
	}

	distributor, err := service.DistributorRepository.FindOneByDistributorIdAndCustId(params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", constant.ErrNoRowsAffected
		}

		return "", err
	}

	return distributor.CustId, nil
}

func (service *distributorServiceImpl) Store(request entity.CreateDistributorBody) (response entity.DistributorResponse, err error) {
	_, err = service.DistributorRepository.FindOneByDistributorCodeAndCustId(request.DistributorCode, request.CustId)
	if err == nil {
		return response, errors.New(distributorCodeDuplicateErrMsg)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return response, err
	}

	timeNow := time.Now().In(time.UTC)
	var mDistributorData model.Distributor
	err = structs.Automapper(request, &mDistributorData)
	if err != nil {
		return response, err
	}

	// Map distributor_setup fields
	if request.DistributorSetup != nil {
		mDistributorData.AllowAddProduct = request.DistributorSetup.AllowAddProduct
		mDistributorData.AllowEditProduct = request.DistributorSetup.AllowEditProduct
		mDistributorData.AllowManagePricing = request.DistributorSetup.AllowManagePricing
		mDistributorData.AllowUploadSecondarySales = request.DistributorSetup.AllowUploadSecondarySales
	}

	mDistributorData.CreatedAt = &timeNow
	mDistributorData.CreatedBy = request.CreatedBy
	mDistributorData.UpdatedAt = &timeNow
	mDistributorData.UpdatedBy = request.CreatedBy

	trx, err := service.DistributorRepository.TrxBegin()
	if err != nil {
		return response, err
	}

	err = trx.InsertDistributor(&mDistributorData)
	if err != nil {
		trx.TrxRollback()
		return response, mapDistributorDuplicateError(err)
	}

	for _, detail := range request.Contacts {
		var MDistributorContactData model.DistributorContact
		err = structs.Automapper(detail, &MDistributorContactData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		MDistributorContactData.CustId = request.CustId

		err := trx.InsertDistributorContact(mDistributorData.DistributorId, &MDistributorContactData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
	}

	for _, detail := range request.Tax {
		var MDistributorTaxData model.DistributorTax
		err = structs.Automapper(detail, &MDistributorTaxData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		MDistributorTaxData.CustId = request.CustId

		err := trx.InsertDistributorTax(mDistributorData.DistributorId, &MDistributorTaxData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
	}

	trx.TrxCommit()

	return response, nil
}

func (service *distributorServiceImpl) List(dataFilter entity.DistributorQueryFilter, custId string) (data []entity.DistributorListRespone, total int, lastPage int, err error) {
	distributor, total, lastPage, err := service.DistributorRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}
	for _, row := range distributor {
		var vResp entity.DistributorListRespone
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *distributorServiceImpl) LookupList(dataFilter entity.DistributorQueryFilter, custId string) (data []entity.DistributorLookupResponse, total int, lastPage int, err error) {
	distributor, total, lastPage, err := service.DistributorRepository.FindAllLookupByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range distributor {
		var vResp entity.DistributorLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *distributorServiceImpl) Detail(params entity.DetailDistributorParams) (response entity.DistributorResponse, err error) {
	distributor, err := service.DistributorRepository.FindOneByDistributorIdAndCustId(params)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(distributor, &response)
	if err != nil {
		return response, err
	}

	// Map distributor_setup to response
	response.DistributorSetup = &entity.DistributorSetup{
		AllowAddProduct:           distributor.AllowAddProduct,
		AllowEditProduct:          distributor.AllowEditProduct,
		AllowManagePricing:        distributor.AllowManagePricing,
		AllowUploadSecondarySales: distributor.AllowUploadSecondarySales,
	}

	distributorContacts, err := service.DistributorRepository.FindAllDistributorContactByDistIdAndCustId(params)
	if err != nil {
		return response, err
	}

	detailContact := make([]entity.DistributorContact, 0)
	for _, distributorContact := range distributorContacts {
		var distributorContactResp entity.DistributorContact
		err = structs.Automapper(distributorContact, &distributorContactResp)
		if err != nil {
			return response, err
		}
		detailContact = append(detailContact, distributorContactResp)
	}
	response.DistributorContact = detailContact

	distributorTaxs, err := service.DistributorRepository.FindAllDistributorTaxByDistIdAndCustId(params)
	if err != nil {
		return response, err
	}

	detailTax := make([]entity.DistributorTax, 0)
	for _, distributorTax := range distributorTaxs {
		var distributorTaxResp entity.DistributorTax
		err = structs.Automapper(distributorTax, &distributorTaxResp)
		if err != nil {
			return response, err
		}
		detailTax = append(detailTax, distributorTaxResp)
	}
	response.DistributorTax = detailTax
	return response, nil
}

func (service *distributorServiceImpl) Update(distributorId int64, request entity.UpdateDistributorRequest) (err error) {
	targetCustID, err := service.resolveUpdateCustID(distributorId, request.CustId)
	if err != nil {
		return err
	}

	if request.DistributorCode != nil {
		distributor, err := service.DistributorRepository.FindOneByDistributorCodeAndCustId(*request.DistributorCode, targetCustID)
		if err == nil && distributor.DistributorId != distributorId {
			return errors.New(distributorCodeDuplicateErrMsg)
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	trx, err := service.DistributorRepository.TrxBegin()
	if err != nil {
		return err
	}

	// Handle distributor_setup fields for update
	if request.DistributorSetup != nil {
		allowAddProduct := request.DistributorSetup.AllowAddProduct
		allowEditProduct := request.DistributorSetup.AllowEditProduct
		allowManagePricing := request.DistributorSetup.AllowManagePricing
		allowUploadSecondarySales := request.DistributorSetup.AllowUploadSecondarySales
		request.AllowAddProduct = &allowAddProduct
		request.AllowEditProduct = &allowEditProduct
		request.AllowManagePricing = &allowManagePricing
		request.AllowUploadSecondarySales = &allowUploadSecondarySales
	}
	request.CustId = targetCustID

	err = trx.Update(distributorId, request)
	if err != nil {
		trx.TrxRollback()
		return mapDistributorDuplicateError(err)
	}
	distributorContactIDs := []int64{}
	distributorTaxIDs := []int64{}

	for _, detail := range request.Contacts {
		if detail.DistributorContactId != nil {
			distributorContactIDs = append(distributorContactIDs, *detail.DistributorContactId)
		}
	}

	for _, detail := range request.Tax {
		if detail.DistributorTaxId != nil {
			distributorTaxIDs = append(distributorTaxIDs, *detail.DistributorTaxId)
		}
	}

	if request.Contacts != nil && len(distributorContactIDs) == 0 {
		err := trx.DeleteAllDistributorContact(distributorId)
		if err != nil {
			trx.TrxRollback()
			return err
		}
	}
	if len(distributorContactIDs) > 0 {
		err := trx.DeleteDistributorContactNotIn(distributorId, distributorContactIDs)
		if err != nil {
			trx.TrxRollback()
			return err
		}
	}
	if request.Tax != nil && len(distributorTaxIDs) == 0 {
		err := trx.DeleteAllDistributorTax(distributorId)
		if err != nil {
			trx.TrxRollback()
			return err
		}
	}
	if len(distributorTaxIDs) > 0 {
		err := trx.DeleteDistributorTaxNotIn(distributorId, distributorTaxIDs)
		if err != nil {
			trx.TrxRollback()
			return err
		}
	}

	for _, detail := range request.Contacts {
		if detail.DistributorContactId == nil || *detail.DistributorContactId == 0 {
			detail.DistributorContactId = nil
			var MdistributorContact model.DistributorContact
			err = structs.Automapper(detail, &MdistributorContact)
			if err != nil {
				trx.TrxRollback()
				return err
			}
			MdistributorContact.CustId = request.CustId
			err := trx.InsertDistributorContact(distributorId, &MdistributorContact)
			if err != nil {
				trx.TrxRollback()
				return err
			}
		} else {
			err := trx.UpdateDistributorContact(distributorId, *detail.DistributorContactId, detail)
			if err != nil {
				trx.TrxRollback()
				return err
			}
		}
	}

	for _, detail := range request.Tax {
		if detail.DistributorTaxId == nil || *detail.DistributorTaxId == 0 {
			detail.DistributorTaxId = nil
			var MdistributorTax model.DistributorTax
			err = structs.Automapper(detail, &MdistributorTax)
			if err != nil {
				trx.TrxRollback()
				return err
			}
			MdistributorTax.CustId = request.CustId
			err := trx.InsertDistributorTax(distributorId, &MdistributorTax)
			if err != nil {
				trx.TrxRollback()
				return err
			}
		} else {
			err := trx.UpdateDistributorTax(distributorId, *detail.DistributorTaxId, detail)
			if err != nil {
				trx.TrxRollback()
				return err
			}
		}
	}

	err = trx.TrxCommit()
	if err != nil {
		return err
	}

	return nil
}

func (service *distributorServiceImpl) Delete(custId string, distributorId int64, userId int64) (err error) {
	err = service.DistributorRepository.Delete(custId, distributorId, userId)
	if err != nil {
		return err
	}

	return err
}

func (service *distributorServiceImpl) ListWithCustomer(dataFilter entity.DistributorQueryFilter, custId string) (data []entity.DistributorCustomerResp, total int, lastPage int, err error) {
	distributor, total, lastPage, err := service.DistributorRepository.FindAllByCustIdWithCustomer(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}
	for _, row := range distributor {
		var vResp entity.DistributorCustomerResp
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
