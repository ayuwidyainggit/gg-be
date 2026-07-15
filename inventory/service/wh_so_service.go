package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
)

type WhSoService interface {
	Store(request entity.CreateWhSoBody) (err error)
	Detail(whSoNo string, custID string) (response entity.WhSoResponse, err error)
	List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.WhSoListResponse, total int64, lastPage int, err error)
	Delete(custId string, whSoNo string, userId int64) (err error)
	Update(whSoNo string, request entity.UpdateWhSoBody) (err error)
}

type whSoServiceImpl struct {
	WhSoRepository repository.WhSoRepository
	Transaction    repository.Dbtransaction
}

func NewWhSoService(WhSoRepository repository.WhSoRepository, transaction repository.Dbtransaction) *whSoServiceImpl {
	return &whSoServiceImpl{
		WhSoRepository: WhSoRepository,
		Transaction:    transaction,
	}
}
func (service *whSoServiceImpl) Store(request entity.CreateWhSoBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.WhSoDate != nil {
		whSoDate, err := str.DateStrToRfc3339String(*request.WhSoDate)
		if err != nil {
			return err
		}
		request.WhSoDate = &whSoDate
	}

	var whSoModel model.WhSo
	err = structs.Automapper(request, &whSoModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.WhSoRepository.Store(txCtx, &whSoModel)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details {
			var whAdjDetModel model.WhSoDet
			whAdjDetModel.CustID = request.CustID
			whAdjDetModel.WhSoNo = whSoModel.WhSoNo
			err = structs.Automapper(Detail, &whAdjDetModel)
			if err != nil {
				return err
			}

			err = service.WhSoRepository.StoreDetail(txCtx, &whAdjDetModel)
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
func (service *whSoServiceImpl) Detail(whSoNo string, custID string) (response entity.WhSoResponse, err error) {
	whAdj, err := service.WhSoRepository.FindByNo(whSoNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(whAdj, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.WhSoRepository.FindDetail(whSoNo, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.WhSoDetResponse
	for _, detail := range Details {
		var detailData entity.WhSoDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		DetailsData = append(DetailsData, detailData)
	}

	whSoDate := whAdj.WhSoDate.Format("2006-01-02")
	response.WhSoDate = &whSoDate

	response.Details = DetailsData
	return response, nil
}
func (service *whSoServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.WhSoListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.WhSoRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}
	if len(whAdjs) > 0 {
		for _, row := range whAdjs {
			var vResp entity.WhSoListResponse
			structs.Automapper(row, &vResp)
			if row.WhSoDate != nil {
				WhSoDate := row.WhSoDate.Format("2006-01-02")
				vResp.WhSoDate = &WhSoDate
			}
			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}
func (service *whSoServiceImpl) Delete(custId string, whSoNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.WhSoRepository.Delete(txCtx, custId, whSoNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *whSoServiceImpl) Update(whSoNo string, request entity.UpdateWhSoBody) (err error) {
	c := context.Background()

	if request.WhSoDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		if *request.WhSoDate != "" {
			WhTrfDate, err := str.DateStrToRfc3339String(*request.WhSoDate)
			if err != nil {
				return err
			}
			request.WhSoDate = &WhTrfDate
		}

	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.WhSo
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.WhSoRepository.Update(txCtx, whSoNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details {
			if detail.WhSoDetId != nil {
				DetailIds = append(DetailIds, *detail.WhSoDetId)
			}
		}
		if len(DetailIds) > 0 {
			err := service.WhSoRepository.DeleteDetailNotInIDs(txCtx, whSoNo, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details {
			// parse time format YYYY-mm-dd to Rfc3339

			var whSoDetModel model.WhSoDet

			err = structs.Automapper(detail, &whSoDetModel)
			if err != nil {
				return err
			}
			whSoDetModel.CustID = request.CustID
			whSoDetModel.WhSoNo = whSoNo
			if detail.WhSoDetId == nil || *detail.WhSoDetId == 0 {
				detail.WhSoDetId = nil
				err = service.WhSoRepository.StoreDetail(txCtx, &whSoDetModel)
				if err != nil {
					return err
				}
			} else {
				whSoDetModel.CustID = ""
				err = service.WhSoRepository.UpdateGrDetail(txCtx, &whSoDetModel)
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
