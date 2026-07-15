package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
)

type VanBsUlService interface {
	Store(request entity.CreateVanBsUlBody) (err error)
	Detail(vanBsUlNo string, custID string) (response entity.VanBsUlResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.VanBsUlListResponse, total int64, lastPage int, err error)
	Delete(custId string, vanBsUlNo string, userId int64) (err error)
	Update(vanBsUlNo string, request entity.UpdateVanBsUlBody) (err error)
}

type vanBsUlServiceImpl struct {
	Repository  repository.VanBsUlRepository
	Transaction repository.Dbtransaction
}

func NewVanBsUlService(Repository repository.VanBsUlRepository, transaction repository.Dbtransaction) *vanBsUlServiceImpl {
	return &vanBsUlServiceImpl{
		Repository:  Repository,
		Transaction: transaction,
	}
}

func (service *vanBsUlServiceImpl) Store(request entity.CreateVanBsUlBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.VanBsUlDate != nil {
		vanBsUlDate, err := str.DateStrToRfc3339String(*request.VanBsUlDate)
		if err != nil {
			return err
		}
		request.VanBsUlDate = &vanBsUlDate
	}

	var Model model.VanBsUl
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.Store(txCtx, &Model)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details {
			// parse time format YYYY-mm-dd to Rfc3339
			if Detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*Detail.ExpDate)
				if err != nil {
					return err
				}
				Detail.ExpDate = &expDate
			}
			var vanBsUlDetModel model.VanBsUlDet
			vanBsUlDetModel.CustID = request.CustID
			vanBsUlDetModel.VanBsUlNo = Model.VanBsUlNo
			err = structs.Automapper(Detail, &vanBsUlDetModel)
			if err != nil {
				return err
			}

			err = service.Repository.StoreDetail(txCtx, &vanBsUlDetModel)
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
func (service *vanBsUlServiceImpl) Detail(vanBsUlNo string, custID string) (response entity.VanBsUlResponse, err error) {
	gds, err := service.Repository.FindByNo(vanBsUlNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(gds, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.Repository.FindDetail(vanBsUlNo, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.VanBsUlDetResponse
	for _, detail := range Details {
		var detailData entity.VanBsUlDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		if detailData.ExpDate != nil {
			expDate := detail.ExpDate.Format("2006-01-02")
			detailData.ExpDate = &expDate
		}

		DetailsData = append(DetailsData, detailData)
	}

	vanBsUlDate := gds.VanBsUlDate.Format("2006-01-02")
	response.VanBsUlDate = &vanBsUlDate

	response.Details = DetailsData
	return response, nil
}
func (service *vanBsUlServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.VanBsUlListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range whAdjs {
		var vResp entity.VanBsUlListResponse
		structs.Automapper(row, &vResp)
		if row.VanBsUlDate != nil {
			vanBsUlDate := row.VanBsUlDate.Format("2006-01-02")
			vResp.VanBsUlDate = &vanBsUlDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *vanBsUlServiceImpl) Delete(custId string, vanBsUlNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, vanBsUlNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (service *vanBsUlServiceImpl) Update(vanBsUlNo string, request entity.UpdateVanBsUlBody) (err error) {
	c := context.Background()

	if request.VanBsUlDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		if *request.VanBsUlDate != "" {
			vanBsUlDate, err := str.DateStrToRfc3339String(*request.VanBsUlDate)
			if err != nil {
				return err
			}
			request.VanBsUlDate = &vanBsUlDate
		}

	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.VanBsUl
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Update(txCtx, vanBsUlNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details {
			if detail.VanBsUlDetID != nil {
				DetailIds = append(DetailIds, *detail.VanBsUlDetID)
			}
		}
		if len(DetailIds) > 0 {
			err := service.Repository.DeleteDetailNotInIDs(txCtx, vanBsUlNo, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details {
			// parse time format YYYY-mm-dd to Rfc3339
			if detail.ExpDate != nil {
				// parse time format YYYY-mm-dd to Rfc3339
				if *detail.ExpDate != "" {
					expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
					if err != nil {
						return err
					}
					detail.ExpDate = &expDate
				}

			}
			var vanBsUlDetModel model.VanBsUlDet

			err = structs.Automapper(detail, &vanBsUlDetModel)
			if err != nil {
				return err
			}
			vanBsUlDetModel.CustID = request.CustID
			vanBsUlDetModel.VanBsUlNo = vanBsUlNo
			if detail.VanBsUlDetID == nil || *detail.VanBsUlDetID == 0 {
				vanBsUlDetModel.VanBsUlDetID = nil
				err = service.Repository.StoreDetail(txCtx, &vanBsUlDetModel)
				if err != nil {
					return err
				}
			} else {
				vanBsUlDetModel.CustID = ""
				err = service.Repository.UpdateDetail(txCtx, &vanBsUlDetModel)
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
