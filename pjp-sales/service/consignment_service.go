package service

import (
	"context"
	"sales/entity"
	"sales/model"
	"sales/pkg/str"
	"sales/pkg/structs"
	"sales/repository"
)

type ConsignmentService interface {
	Store(request entity.CreateConsignmentBody) (err error)
	Detail(consNo string, custID string) (response entity.ConsignmentResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.ConsignmentListResponse, total int64, lastPage int, err error)
	Delete(custId string, consNo string, userId int64) (err error)
	Update(consNo string, request entity.UpdateConsignmentBody) (err error)
}

func NewConsignmentService(repository repository.ConsignmentRepository, transaction repository.Dbtransaction) *consignmentServiceImpl {
	return &consignmentServiceImpl{
		Repository:  repository,
		Transaction: transaction,
	}
}

type consignmentServiceImpl struct {
	Repository  repository.ConsignmentRepository
	Transaction repository.Dbtransaction
}

func (service *consignmentServiceImpl) Store(request entity.CreateConsignmentBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.ConsDate != nil {
		consDate, err := str.DateStrToRfc3339String(*request.ConsDate)
		if err != nil {
			return err
		}
		request.ConsDate = &consDate
	}

	var consModel model.Consignment
	err = structs.Automapper(request, &consModel)
	if err != nil {
		return err
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.Store(txCtx, &consModel)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details {
			var consDetModel model.ConsignmentDet

			if Detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*Detail.ExpDate)
				if err != nil {
					return err
				}
				Detail.ExpDate = &expDate
			}

			err = structs.Automapper(Detail, &consDetModel)
			if err != nil {
				return err
			}
			consDetModel.CustID = request.CustID
			consDetModel.ConsNo = consModel.ConsNo
			err = service.Repository.StoreDetail(txCtx, &consDetModel)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (service *consignmentServiceImpl) Detail(consNo string, custID string) (response entity.ConsignmentResponse, err error) {
	so, err := service.Repository.FindByNo(consNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(so, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.Repository.FindDetail(consNo, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.ConsignDetResponse
	for _, detail := range Details {
		var detailData entity.ConsignDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}
		if detail.ExpDate != nil {
			expDate := detail.ExpDate.Format("2006-01-02")
			detailData.ExpDate = &expDate
		}
		DetailsData = append(DetailsData, detailData)
	}
	if so.ConsDate != nil {
		consDate := so.ConsDate.Format("2006-01-02")
		response.ConsDate = &consDate
	}

	response.Details = DetailsData
	return response, nil
}
func (service *consignmentServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.ConsignmentListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range whAdjs {
		var vResp entity.ConsignmentListResponse
		structs.Automapper(row, &vResp)
		if row.ConsDate != nil {
			consDate := row.ConsDate.Format("2006-01-02")
			vResp.ConsDate = &consDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *consignmentServiceImpl) Delete(custId string, consNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, consNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (service *consignmentServiceImpl) Update(consNo string, request entity.UpdateConsignmentBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.ConsDate != nil {
		consDate, err := str.DateStrToRfc3339String(*request.ConsDate)
		if err != nil {
			return err
		}
		request.ConsDate = &consDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.Consignment
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Update(txCtx, consNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details {
			if detail.ConsDetID != nil {
				DetailIds = append(DetailIds, *detail.ConsDetID)
			}
		}
		if len(DetailIds) > 0 {
			err := service.Repository.DeleteDetailNotInIDs(txCtx, consNo, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details {
			// parse time format YYYY-mm-dd to Rfc3339

			var consDetModel model.ConsignmentDet
			if detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
				if err != nil {
					return err
				}
				detail.ExpDate = &expDate
			}

			err = structs.Automapper(detail, &consDetModel)
			if err != nil {
				return err
			}
			consDetModel.CustID = request.CustID
			consDetModel.ConsNo = consNo
			if detail.ConsDetID == nil || *detail.ConsDetID == 0 {
				consDetModel.ConsDetID = nil
				err = service.Repository.StoreDetail(txCtx, &consDetModel)
				if err != nil {
					return err
				}
			} else {
				consDetModel.CustID = ""
				err = service.Repository.UpdateDetail(txCtx, &consDetModel)
				if err != nil {
					return err
				}

			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
