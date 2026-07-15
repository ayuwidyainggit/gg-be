package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
)

type VanSoService interface {
	Store(request entity.CreateVanSoBody) (err error)
	Detail(vanSoNo string, custID string) (response entity.VanSoResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.VanSoListResponse, total int64, lastPage int, err error)
	Delete(custId string, vanSoNo string, userId int64) (err error)
	Update(vanSoNo string, request entity.UpdateVanSoBody) (err error)
}

type vanSoServiceImpl struct {
	VanSoRepository repository.VanSoRepository
	Transaction     repository.Dbtransaction
}

func NewVanSoService(vanSoRepository repository.VanSoRepository, transaction repository.Dbtransaction) *vanSoServiceImpl {
	return &vanSoServiceImpl{
		VanSoRepository: vanSoRepository,
		Transaction:     transaction,
	}
}
func (service *vanSoServiceImpl) Store(request entity.CreateVanSoBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.VanSoDate != nil {
		vanSoDate, err := str.DateStrToRfc3339String(*request.VanSoDate)
		if err != nil {
			return err
		}
		request.VanSoDate = &vanSoDate
	}

	var vanSoModel model.VanSo
	err = structs.Automapper(request, &vanSoModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.VanSoRepository.Store(txCtx, &vanSoModel)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details {
			var vanSoDetModel model.VanSoDet

			err = structs.Automapper(Detail, &vanSoDetModel)
			if err != nil {
				return err
			}
			vanSoDetModel.CustID = request.CustID
			vanSoDetModel.VanSoNo = vanSoModel.VanSoNo
			err = service.VanSoRepository.StoreDetail(txCtx, &vanSoDetModel)
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
func (service *vanSoServiceImpl) Detail(vanSoNo string, custID string) (response entity.VanSoResponse, err error) {
	vanSo, err := service.VanSoRepository.FindByNo(vanSoNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(vanSo, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.VanSoRepository.FindDetail(vanSoNo, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.VanSoDetResponse
	for _, detail := range Details {
		var detailData entity.VanSoDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		DetailsData = append(DetailsData, detailData)
	}

	vanSoDate := vanSo.VanSoDate.Format("2006-01-02")
	response.VanSoDate = &vanSoDate

	response.Details = DetailsData
	return response, nil
}
func (service *vanSoServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.VanSoListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.VanSoRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range whAdjs {
		var vResp entity.VanSoListResponse
		structs.Automapper(row, &vResp)
		if row.VanSoDate != nil {
			vanSoDate := row.VanSoDate.Format("2006-01-02")
			vResp.VanSoDate = &vanSoDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *vanSoServiceImpl) Delete(custId string, vanSoNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.VanSoRepository.Delete(txCtx, custId, vanSoNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (service *vanSoServiceImpl) Update(vanSoNo string, request entity.UpdateVanSoBody) (err error) {
	c := context.Background()

	if request.VanSoDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		if *request.VanSoDate != "" {
			vanSoDate, err := str.DateStrToRfc3339String(*request.VanSoDate)
			if err != nil {
				return err
			}
			request.VanSoDate = &vanSoDate
		}

	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.VanSo
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.VanSoRepository.Update(txCtx, vanSoNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details {
			if detail.VanSoDetID != nil {
				DetailIds = append(DetailIds, *detail.VanSoDetID)
			}
		}
		if len(DetailIds) > 0 {
			err := service.VanSoRepository.DeleteDetailNotInIDs(txCtx, vanSoNo, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details {
			// parse time format YYYY-mm-dd to Rfc3339

			var vanSoDetModel model.VanSoDet

			err = structs.Automapper(detail, &vanSoDetModel)
			if err != nil {
				return err
			}
			vanSoDetModel.CustID = request.CustID
			vanSoDetModel.VanSoNo = vanSoNo
			if detail.VanSoDetID == nil || *detail.VanSoDetID == 0 {
				vanSoDetModel.VanSoDetID = nil
				err = service.VanSoRepository.StoreDetail(txCtx, &vanSoDetModel)
				if err != nil {
					return err
				}
			} else {
				vanSoDetModel.CustID = ""
				err = service.VanSoRepository.UpdateDetail(txCtx, &vanSoDetModel)
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
