package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/repository"
	"time"
)

type InvoiceDiscService interface {
	Detail(int, string) (entity.InvoiceDiscDetailsResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.InvoiceDiscListResponse, total int, lastPage int, err error)
	Store(entity.CreateInvoiceDiscBody) (entity.InvoiceDiscResponse, error)
	Update(int, entity.UpdateInvoiceDiscRequest) error
	Delete(string, int, int64) error
}

func NewInvoiceDiscService(
	invoiceDiscRepository repository.InvoiceDiscRepository,
	invoiceDiscDetRepository repository.InvoiceDiscDetRepository,

) *InvoiceDiscServiceImpl {
	return &InvoiceDiscServiceImpl{
		InvoiceDiscRepository:    invoiceDiscRepository,
		InvoiceDiscDetRepository: invoiceDiscDetRepository,
	}
}

type InvoiceDiscServiceImpl struct {
	InvoiceDiscRepository    repository.InvoiceDiscRepository
	MProductRepository       repository.ProductRepository
	InvoiceDiscDetRepository repository.InvoiceDiscDetRepository
}

func (service *InvoiceDiscServiceImpl) Detail(invDiscId int, custId string) (response entity.InvoiceDiscDetailsResponse, err error) {
	invoiceDisc, err := service.InvoiceDiscRepository.FindOneByInvDiscIdAndCustId(invDiscId, custId)
	if err != nil {
		return response, err
	}

	invoiceDiscDet, err := service.InvoiceDiscDetRepository.FindOneByInvDiscIdAndCustId(invDiscId, custId)
	if err != nil {
		return response, err
	}

	response.InvDiscId = invoiceDisc.InvDiscId
	response.InvDiscCode = invoiceDisc.InvDiscCode
	response.InvDiscName = invoiceDisc.InvDiscName
	response.IsActive = invoiceDisc.IsActive
	response.UpdatedBy = invoiceDisc.UpdatedBy
	response.UpdatedAt = invoiceDisc.UpdatedAt
	response.Details = make([]entity.InvDiscDet, 0)

	for _, row := range invoiceDiscDet {
		rowResp := entity.InvDiscDet{
			RowNo:    row.RowNo,
			MinValue: row.MinValue,
			MaxValue: row.MaxValue,
			DiscPerc: row.DiscPerc,
		}
		response.Details = append(response.Details, rowResp)
	}

	return response, err
}

func (service *InvoiceDiscServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.InvoiceDiscListResponse, total int, lastPage int, err error) {
	InvoiceDiscs, total, lastPage, err := service.InvoiceDiscRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(InvoiceDiscs) > 0 {
		for _, row := range InvoiceDiscs {
			invoiceDisc := entity.InvoiceDiscListResponse{
				InvDiscId:   row.InvDiscId,
				InvDiscCode: row.InvDiscCode,
				InvDiscName: row.InvDiscName,
				IsActive:    row.IsActive,
				UpdatedBy:   row.UpdatedBy,
				UpdatedAt:   row.UpdatedAt,
			}
			if row.UpdatedByName != nil {
				invoiceDisc.UpdatedByName = *row.UpdatedByName
			}
			data = append(data, invoiceDisc)
		}
	}
	return data, total, lastPage, err
}

func (service *InvoiceDiscServiceImpl) Store(request entity.CreateInvoiceDiscBody) (response entity.InvoiceDiscResponse, err error) {

	// inv_disc_code & cust id validation, if err == nil, this means that code & cust id already exists
	invoiceDisc, err := service.InvoiceDiscRepository.FindOneByInvDiscCodeAndCustId(request.InvDiscCode, request.CustId)
	if err == nil {
		return response, errors.New("inv_disc_code: " + invoiceDisc.InvDiscCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	invDisc := model.InvoiceDisc{
		CustId:      request.CustId,
		InvDiscCode: request.InvDiscCode,
		InvDiscName: request.InvDiscName,
		IsActive:    request.IsActive,
		CreatedAt:   &timeNow,
		CreatedBy:   &request.CreatedBy,
		UpdatedAt:   &timeNow,
		UpdatedBy:   &request.CreatedBy,
	}

	invDiscId, err := service.InvoiceDiscRepository.Store(invDisc)
	if err != nil {
		return response, err
	}

	response.InvDiscId = invDiscId

	invDiscDetails := make([]model.InvoiceDiscDet, 0)
	for _, row := range request.Details {
		invDet := model.InvoiceDiscDet{
			CustId:    request.CustId,
			InvDiscId: invDiscId,
			RowNo:     row.RowNo,
			MinValue:  row.MinValue,
			MaxValue:  row.MaxValue,
			DiscPerc:  row.DiscPerc,
		}
		invDiscDetails = append(invDiscDetails, invDet)
	}
	err = service.InvoiceDiscDetRepository.BulkInsert(invDiscDetails)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *InvoiceDiscServiceImpl) Update(invDiscId int, request entity.UpdateInvoiceDiscRequest) (err error) {

	invoiceDisc, err := service.InvoiceDiscRepository.FindOneByInvDiscCodeAndCustId(request.InvDiscCode, request.CustId)
	if err == nil && invoiceDisc.InvDiscId != invDiscId {
		return errors.New("inv_disc_code: " + invoiceDisc.InvDiscCode + " is already exists")
	}

	detailsReq := request.Details
	request.Details = nil
	err = service.InvoiceDiscRepository.Update(invDiscId, request)
	if err != nil {
		return err
	}

	detailsDb := make([]model.InvoiceDiscDet, 0)
	if len(detailsReq) > 0 {
		err = service.InvoiceDiscDetRepository.DeleteByInvDiscId(invDiscId, request.CustId)
		if err != nil {
			// if no rows affected (no details), set err = nil
			err = nil
		}

		for _, row := range detailsReq {
			rowDb := model.InvoiceDiscDet{
				CustId:    request.CustId,
				InvDiscId: invDiscId,
				RowNo:     row.RowNo,
				MinValue:  row.MinValue,
				MaxValue:  row.MaxValue,
				DiscPerc:  row.DiscPerc,
			}
			detailsDb = append(detailsDb, rowDb)
		}

		err = service.InvoiceDiscDetRepository.BulkInsert(detailsDb)
		if err != nil {
			return err
		}
	}

	return err
}

func (service *InvoiceDiscServiceImpl) Delete(custId string, invDiscId int, userId int64) (err error) {

	err = service.InvoiceDiscRepository.Delete(custId, invDiscId, userId)
	if err != nil {
		return err
	}

	return err
}
