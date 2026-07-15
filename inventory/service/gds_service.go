package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
)

type GdsService interface {
	Store(request entity.CreateGdsBody) (err error)
	Detail(gdsNo, custID, langId, parentCustId string) (response entity.GdsResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.GdsListResponse, total int64, lastPage int, err error)
	Delete(custId string, gdsNo string, userId int64) (err error)
	Update(gdsNo string, request entity.UpdateGdsBody) (err error)
}

type gdsServiceImpl struct {
	GdsRepository repository.GdsRepository
	Transaction   repository.Dbtransaction
}

func NewGdsService(gdsRepository repository.GdsRepository, transaction repository.Dbtransaction) *gdsServiceImpl {
	return &gdsServiceImpl{
		GdsRepository: gdsRepository,
		Transaction:   transaction,
	}
}
func (service *gdsServiceImpl) Store(request entity.CreateGdsBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	gdsDate, err := str.DateStrToRfc3339String(request.GdsDate)
	if err != nil {
		return err
	}
	request.GdsDate = gdsDate

	var gdsModel model.Gds
	err = structs.Automapper(request, &gdsModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.GdsRepository.Store(txCtx, &gdsModel)
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
			var gdsDetModel model.GdsDet
			gdsDetModel.CustID = request.CustID
			gdsDetModel.GdsNo = gdsModel.GdsNo
			err = structs.Automapper(Detail, &gdsDetModel)
			if err != nil {
				return err
			}

			err = service.GdsRepository.StoreDetail(txCtx, &gdsDetModel)
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

func (service *gdsServiceImpl) Detail(gdsNo, custID, langId, parentCustId string) (response entity.GdsResponse, err error) {
	gds, err := service.GdsRepository.FindByNo(gdsNo, custID, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(gds, &response)
	if err != nil {
		return response, err
	}

	details, err := service.GdsRepository.FindDetail(gdsNo, custID, langId)
	if err != nil {
		return response, err
	}
	var detailsData []entity.GdsDetResponse
	for _, detail := range details {
		var detailData entity.GdsDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		if detailData.ExpDate != nil {
			expDate := detail.ExpDate.Format("2006-01-02")
			detailData.ExpDate = &expDate
		}

		detailsData = append(detailsData, detailData)
	}

	gdsDate := gds.GdsDate.Format("2006-01-02")
	response.GdsDate = &gdsDate

	response.Details = detailsData
	return response, nil
}
func (service *gdsServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.GdsListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.GdsRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range whAdjs {
		var vResp entity.GdsListResponse
		structs.Automapper(row, &vResp)
		if row.GdsDate != nil {
			GdsDate := row.GdsDate.Format("2006-01-02")
			vResp.GdsDate = &GdsDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *gdsServiceImpl) Delete(custId string, gdsNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.GdsRepository.Delete(txCtx, custId, gdsNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *gdsServiceImpl) Update(gdsNo string, request entity.UpdateGdsBody) (err error) {
	c := context.Background()

	if request.GdsDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		if *request.GdsDate != "" {
			gdsDate, err := str.DateStrToRfc3339String(*request.GdsDate)
			if err != nil {
				return err
			}
			request.GdsDate = &gdsDate
		}

	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.Gds
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.GdsRepository.Update(txCtx, gdsNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details {
			if detail.GdsDetId != nil {
				DetailIds = append(DetailIds, *detail.GdsDetId)
			}
		}
		if len(DetailIds) > 0 {
			err := service.GdsRepository.DeleteDetailNotInIDs(txCtx, gdsNo, DetailIds)
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
			var gdsDetModel model.GdsDet

			err = structs.Automapper(detail, &gdsDetModel)
			if err != nil {
				return err
			}
			gdsDetModel.CustID = request.CustID
			gdsDetModel.GdsNo = gdsNo
			if detail.GdsDetId == nil || *detail.GdsDetId == 0 {
				detail.GdsDetId = nil
				err = service.GdsRepository.StoreDetail(txCtx, &gdsDetModel)
				if err != nil {
					return err
				}
			} else {
				err = service.GdsRepository.UpdateDetail(txCtx, &gdsDetModel)
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
