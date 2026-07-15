package service

import (
	"context"
	"sales/entity"
	"sales/model"
	"sales/pkg/str"
	"sales/pkg/structs"
	"sales/repository"
)

type TlsService interface {
	Store(request entity.CreateTlsBody) (err error)
	Detail(TlsId int, custID string) (response entity.TlsResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.TlsListResponse, total int64, lastPage int, err error)
	Update(TlsId int, request entity.UpdateTlsBody) (err error)
	Delete(custId string, TlsId int, userId int64) (err error)
}

func NewTlsService(tlsRepository repository.TlsRepository, transaction repository.Dbtransaction) *tlsServiceImpl {
	return &tlsServiceImpl{
		TlsRepository: tlsRepository,
		Transaction:   transaction,
	}
}

type tlsServiceImpl struct {
	TlsRepository repository.TlsRepository
	Transaction   repository.Dbtransaction
}

func (service *tlsServiceImpl) Store(request entity.CreateTlsBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.TlsDate != nil {
		tlsDate, err := str.DateStrToRfc3339String(*request.TlsDate)
		if err != nil {
			return err
		}
		request.TlsDate = &tlsDate
	}

	var tlsModel model.Tls
	err = structs.Automapper(request, &tlsModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.TlsRepository.Store(txCtx, &tlsModel)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details {

			var tlsDetModel model.TlsDet
			err = structs.Automapper(Detail, &tlsDetModel)
			if err != nil {
				return err
			}
			tlsDetModel.TlsId = tlsModel.TlsId
			tlsDetModel.CustId = request.CustId

			err = service.TlsRepository.StoreDetail(txCtx, &tlsDetModel)
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

func (service *tlsServiceImpl) Detail(TlsId int, custID string) (response entity.TlsResponse, err error) {
	tls, err := service.TlsRepository.FindByNo(TlsId, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(tls, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.TlsRepository.FindDetail(TlsId, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.TlsDetResponse
	for _, detail := range Details {
		var detailData entity.TlsDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}
		DetailsData = append(DetailsData, detailData)
	}
	if tls.TlsDate != nil {
		tlsDate := tls.TlsDate.Format("2006-01-02")
		response.TlsDate = &tlsDate
	}
	response.Details = DetailsData
	return response, nil
}

func (service *tlsServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.TlsListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.TlsRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range whAdjs {
		var vResp entity.TlsListResponse
		structs.Automapper(row, &vResp)
		if row.TlsDate != nil {
			tlsDate := row.TlsDate.Format("2006-01-02")
			vResp.TlsDate = &tlsDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *tlsServiceImpl) Update(TlsId int, request entity.UpdateTlsBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.TlsDate != nil {
		TlsDate, err := str.DateStrToRfc3339String(*request.TlsDate)
		if err != nil {
			return err
		}
		request.TlsDate = &TlsDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.Tls
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustId = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.TlsRepository.Update(txCtx, TlsId, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details {
			if detail.TlsDetId != nil {
				DetailIds = append(DetailIds, *detail.TlsDetId)
			}
		}
		if len(DetailIds) > 0 {
			err := service.TlsRepository.DeleteDetailNotInIDs(txCtx, TlsId, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details {
			// parse time format YYYY-mm-dd to Rfc3339

			var tlsDetModel model.TlsDet

			err = structs.Automapper(detail, &tlsDetModel)
			if err != nil {
				return err
			}
			tlsDetModel.CustId = request.CustId
			tlsDetModel.TlsId = int64(TlsId)
			if detail.TlsDetId == nil || *detail.TlsDetId == 0 {
				tlsDetModel.TlsDetId = nil
				err = service.TlsRepository.StoreDetail(txCtx, &tlsDetModel)
				if err != nil {
					return err
				}
			} else {
				tlsDetModel.CustId = ""
				err = service.TlsRepository.UpdateDetail(txCtx, &tlsDetModel)
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

func (service *tlsServiceImpl) Delete(custId string, TlsId int, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.TlsRepository.Delete(txCtx, custId, TlsId, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
