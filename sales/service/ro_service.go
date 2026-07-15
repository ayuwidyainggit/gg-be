package service

import (
	"context"
	"sales/entity"
	"sales/model"
	"sales/pkg/str"
	"sales/pkg/structs"
	"sales/repository"
)

type RoService interface {
	Store(request entity.CreateRoBody) (err error)
	Detail(RoNo string, custID string) (response entity.RoResponse, err error)
	List(dataFilter entity.RoQueryFilter) (data []entity.RoListResponse, total int64, lastPage int, err error)
	Update(roNo string, request entity.UpdateRoBody) (err error)
	Delete(custId string, RoNo string, userId int64) (err error)
}

func NewRoService(roRepository repository.RoRepository, transaction repository.Dbtransaction) *roServiceImpl {
	return &roServiceImpl{
		RoRepository: roRepository,
		Transaction:  transaction,
	}
}

type roServiceImpl struct {
	RoRepository repository.RoRepository
	Transaction  repository.Dbtransaction
}

var detCreatemapperToModel = func(det entity.CreateRoDetBody, request entity.CreateRoBody, roNo string) (model.RoDet, error) {
	var gdsDetModel model.RoDet

	// parse time format YYYY-mm-dd to Rfc3339
	if det.ExpDate != nil {
		expDate, err := str.DateStrToRfc3339String(*det.ExpDate)
		if err != nil {
			return gdsDetModel, err
		}
		det.ExpDate = &expDate
	}
	gdsDetModel.CustId = request.CustId
	gdsDetModel.RoNo = roNo
	err := structs.Automapper(det, &gdsDetModel)
	if err != nil {
		return gdsDetModel, err
	}
	return gdsDetModel, nil
}

func (service *roServiceImpl) Store(request entity.CreateRoBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.RoDate != nil {
		roDate, err := str.DateStrToRfc3339String(*request.RoDate)
		if err != nil {
			return err
		}
		request.RoDate = &roDate
	}

	if request.ValDate != nil {
		valDate, err := str.DateStrToRfc3339String(*request.ValDate)
		if err != nil {
			return err
		}
		request.ValDate = &valDate
	}

	if request.DueDate != nil {
		DueDate, err := str.DateStrToRfc3339String(*request.DueDate)
		if err != nil {
			return err
		}
		request.DueDate = &DueDate
	}

	if request.DeliveryDate != nil {
		deliveryDate, err := str.DateStrToRfc3339String(*request.DeliveryDate)
		if err != nil {
			return err
		}
		request.DeliveryDate = &deliveryDate
	}

	var roModel model.Ro
	err = structs.Automapper(request, &roModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.RoRepository.Store(txCtx, &roModel)
		if err != nil {
			return err
		}

		for _, detail := range request.Details.Normal {
			var gdsDetModel model.RoDet

			// parse time format YYYY-mm-dd to Rfc3339
			if detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
				if err != nil {
					return err
				}
				detail.ExpDate = &expDate
			}
			gdsDetModel.CustId = request.CustId
			gdsDetModel.RoNo = roModel.RoNo
			gdsDetModel.ItemType = 1
			err := structs.Automapper(detail, &gdsDetModel)
			if err != nil {
				return err
			}
			err = service.RoRepository.StoreDetail(txCtx, &gdsDetModel)
			if err != nil {
				return err
			}

		}
		for _, Detail := range request.Details.Promo {
			var gdsDetModel model.RoDet

			// parse time format YYYY-mm-dd to Rfc3339
			if Detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*Detail.ExpDate)
				if err != nil {
					return err
				}
				Detail.ExpDate = &expDate
			}
			gdsDetModel.CustId = request.CustId
			gdsDetModel.RoNo = roModel.RoNo
			gdsDetModel.ItemType = 2
			err := structs.Automapper(Detail, &gdsDetModel)
			if err != nil {
				return err
			}
			err = service.RoRepository.StoreDetail(txCtx, &gdsDetModel)
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

func (service *roServiceImpl) Detail(RoNo string, custID string) (response entity.RoResponse, err error) {
	ro, err := service.RoRepository.FindByNo(RoNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ro, &response)
	if err != nil {
		return response, err
	}

	details, err := service.RoRepository.FindDetail(RoNo, custID)
	if err != nil {
		return response, err
	}
	for _, detail := range details {
		var detailData entity.RoDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}
		if detail.ExpDate != nil {
			expDate := detail.ExpDate.Format("2006-01-02")
			detailData.ExpDate = &expDate
		}
		if detailData.ItemType == 1 {
			response.Details.Normal = append(response.Details.Normal, detailData)
		} else {
			response.Details.Promo = append(response.Details.Promo, detailData)
		}
	}
	if ro.RoDate != nil {
		roDate := ro.RoDate.Format("2006-01-02")
		response.RoDate = &roDate
	}
	if ro.ValDate != nil {
		ValDate := ro.ValDate.Format("2006-01-02")
		response.ValDate = &ValDate
	}
	if ro.DeliveryDate != nil {
		DelivDate := ro.DeliveryDate.Format("2006-01-02")
		response.DeliveryDate = &DelivDate
	}

	statusName := response.GenerateDataStatusName()
	response.DataStatusName = statusName

	payTypeName := response.GeneratePayTypeName()
	response.PayTypeName = payTypeName

	return response, nil
}

func (service *roServiceImpl) List(dataFilter entity.RoQueryFilter) (data []entity.RoListResponse, total int64, lastPage int, err error) {
	realOrders, total, lastPage, err := service.RoRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range realOrders {
		var vResp entity.RoListResponse
		structs.Automapper(row, &vResp)
		if row.RoDate != nil {
			roDate := row.RoDate.Format("2006-01-02")
			vResp.RoDate = &roDate
		}
		if row.ValDate != nil {
			ValDate := row.ValDate.Format("2006-01-02")
			vResp.ValDate = &ValDate
		}
		if row.DeliveryDate != nil {
			DelivDate := row.DeliveryDate.Format("2006-01-02")
			vResp.DeliveryDate = &DelivDate
		}

		statusName := vResp.GenerateDataStatusName()
		vResp.DataStatusName = statusName

		payTypeName := vResp.GeneratePayTypeName()
		vResp.PayTypeName = payTypeName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *roServiceImpl) Update(roNo string, request entity.UpdateRoBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.RoDate != nil {
		RoDate, err := str.DateStrToRfc3339String(*request.RoDate)
		if err != nil {
			return err
		}
		request.RoDate = &RoDate
	}

	if request.ValDate != nil {
		ValDate, err := str.DateStrToRfc3339String(*request.ValDate)
		if err != nil {
			return err
		}
		request.ValDate = &ValDate
	}

	if request.DueDate != nil {
		DueDate, err := str.DateStrToRfc3339String(*request.DueDate)
		if err != nil {
			return err
		}
		request.DueDate = &DueDate
	}

	if request.DeliveryDate != nil {
		deliveryDate, err := str.DateStrToRfc3339String(*request.DeliveryDate)
		if err != nil {
			return err
		}
		request.DeliveryDate = &deliveryDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.Ro
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.RoRepository.Update(txCtx, roNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details.Normal {
			if detail.RoDetId != nil {
				DetailIds = append(DetailIds, *detail.RoDetId)
			}
		}
		for _, detail := range request.Details.Promo {
			if detail.RoDetId != nil {
				DetailIds = append(DetailIds, *detail.RoDetId)
			}
		}
		if len(DetailIds) > 0 {
			err := service.RoRepository.DeleteDetailNotInIDs(txCtx, roNo, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details.Normal {
			// parse time format YYYY-mm-dd to Rfc3339

			var roDetModel model.RoDet
			if detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
				if err != nil {
					return err
				}
				detail.ExpDate = &expDate
			}

			err = structs.Automapper(detail, &roDetModel)
			if err != nil {
				return err
			}
			roDetModel.CustId = request.CustId
			roDetModel.RoNo = roNo
			if detail.RoDetId == nil || *detail.RoDetId == 0 {
				roDetModel.RoDetID = nil
				err = service.RoRepository.StoreDetail(txCtx, &roDetModel)
				if err != nil {
					return err
				}
			} else {
				roDetModel.CustId = ""
				err = service.RoRepository.UpdateDetail(txCtx, &roDetModel)
				if err != nil {
					return err
				}

			}
		}
		for _, detail := range request.Details.Promo {
			// parse time format YYYY-mm-dd to Rfc3339

			var roDetModel model.RoDet
			if detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
				if err != nil {
					return err
				}
				detail.ExpDate = &expDate
			}

			err = structs.Automapper(detail, &roDetModel)
			if err != nil {
				return err
			}
			roDetModel.CustId = request.CustId
			roDetModel.RoNo = roNo
			if detail.RoDetId == nil || *detail.RoDetId == 0 {
				roDetModel.RoDetID = nil
				err = service.RoRepository.StoreDetail(txCtx, &roDetModel)
				if err != nil {
					return err
				}
			} else {
				roDetModel.CustId = ""
				err = service.RoRepository.UpdateDetail(txCtx, &roDetModel)
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

func (service *roServiceImpl) Delete(custId string, RoNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.RoRepository.Delete(txCtx, custId, RoNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
