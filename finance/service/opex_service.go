package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"finance/pkg/structs"
	"finance/repository"
)

type OpexService interface {
	Store(request entity.CreateOpexTrBody) (err error)
	Detail(opexTrNo string, custID string) (response entity.OpexTrResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.OpexTrListResponse, total int64, lastPage int, err error)
	Delete(custId string, opexTrNo string, userId int64) (err error)
	Update(opexTrNo string, request entity.UpdateOpexTrBody) (err error)
}

type opexServiceImpl struct {
	Repository  repository.OpexRepository
	Transaction repository.Dbtransaction
}

func NewOpexService(opexSoRepository repository.OpexRepository, transaction repository.Dbtransaction) *opexServiceImpl {
	return &opexServiceImpl{
		Repository:  opexSoRepository,
		Transaction: transaction,
	}
}
func (service *opexServiceImpl) Store(request entity.CreateOpexTrBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.OpexTrDate != nil {
		opexDate, err := str.DateStrToRfc3339String(*request.OpexTrDate)
		if err != nil {
			return err
		}
		request.OpexTrDate = &opexDate
	}

	var opexModel model.OpexTr
	err = structs.Automapper(request, &opexModel)
	if err != nil {
		return err
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.Store(txCtx, &opexModel)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details {
			var opexDetModel model.OpexTrDet

			err = structs.Automapper(Detail, &opexDetModel)
			if err != nil {
				return err
			}
			opexDetModel.CustID = request.CustID
			opexDetModel.OpexTrNo = opexModel.OpexTrNo
			err = service.Repository.StoreDetail(txCtx, &opexDetModel)
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
func (service *opexServiceImpl) Detail(opexTrNo string, custID string) (response entity.OpexTrResponse, err error) {
	opex, err := service.Repository.FindByNo(opexTrNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(opex, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.Repository.FindDetail(opexTrNo, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.OpexTrDetResponse
	for _, detail := range Details {
		var detailData entity.OpexTrDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}
		DetailsData = append(DetailsData, detailData)
	}
	if opex.OpexTrDate != nil {
		opexDate := opex.OpexTrDate.Format("2006-01-02")
		response.OpexTrDate = &opexDate
	}

	response.Details = DetailsData
	return response, nil
}
func (service *opexServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.OpexTrListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range whAdjs {
		var vResp entity.OpexTrListResponse
		structs.Automapper(row, &vResp)
		if row.OpexTrDate != nil {
			opexDate := row.OpexTrDate.Format("2006-01-02")
			vResp.OpexTrDate = &opexDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *opexServiceImpl) Delete(custId string, opexTrNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, opexTrNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (service *opexServiceImpl) Update(opexTrNo string, request entity.UpdateOpexTrBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.OpexTrDate != nil {
		opexDate, err := str.DateStrToRfc3339String(*request.OpexTrDate)
		if err != nil {
			return err
		}
		request.OpexTrDate = &opexDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.OpexTr
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Update(txCtx, opexTrNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details {
			if detail.OpexTrDetID != nil {
				DetailIds = append(DetailIds, *detail.OpexTrDetID)
			}
		}
		if len(DetailIds) > 0 {
			err := service.Repository.DeleteDetailNotInIDs(txCtx, opexTrNo, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details {
			// parse time format YYYY-mm-dd to Rfc3339

			var opexDetModel model.OpexTrDet

			err = structs.Automapper(detail, &opexDetModel)
			if err != nil {
				return err
			}
			opexDetModel.CustID = request.CustID
			opexDetModel.OpexTrNo = opexTrNo
			if detail.OpexTrDetID == nil || *detail.OpexTrDetID == 0 {
				opexDetModel.OpexTrDetID = nil
				err = service.Repository.StoreDetail(txCtx, &opexDetModel)
				if err != nil {
					return err
				}
			} else {
				opexDetModel.CustID = ""
				err = service.Repository.UpdateDetail(txCtx, &opexDetModel)
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
