package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
)

type itemStChServiceImpl struct {
	ItemStChRepository repository.ItemStChRepository
	Transaction        repository.Dbtransaction
}
type ItemStChService interface {
	Store(request entity.CreateItemStChBody) (err error)
	Detail(isc_no string, custID string) (response entity.ItemStChResponse, err error)
	List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.ItemStChListResponse, total int64, lastPage int, err error)
	Delete(custId string, isc_no string, userId int64) (err error)
	Update(isc_no string, request entity.UpdateItemStChBody) (err error)
}

func NewItemStChService(ItemStChRepository repository.ItemStChRepository, transaction repository.Dbtransaction) *itemStChServiceImpl {
	return &itemStChServiceImpl{
		ItemStChRepository: ItemStChRepository,
		Transaction:        transaction,
	}
}
func (service *itemStChServiceImpl) Store(request entity.CreateItemStChBody) (err error) {
	c := context.Background()

	if request.IscDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		IscDate, err := str.DateStrToRfc3339String(*request.IscDate)
		if err != nil {
			return err
		}
		request.IscDate = &IscDate
	}

	var itemSchModel model.ItemStCh
	err = structs.Automapper(request, &itemSchModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.ItemStChRepository.Store(txCtx, &itemSchModel)
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
			Detail.CustID = request.CustID
			Detail.IscNo = itemSchModel.IscNo
			var itemStChDetModel model.ItemStChDet
			err = structs.Automapper(Detail, &itemStChDetModel)
			if err != nil {
				return err
			}

			err = service.ItemStChRepository.StoreDetail(txCtx, &itemStChDetModel)
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
func (service *itemStChServiceImpl) Detail(isc_no string, custID string) (response entity.ItemStChResponse, err error) {
	itemStCh, err := service.ItemStChRepository.FindByNo(isc_no, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(itemStCh, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.ItemStChRepository.FindSmpIssuedetail(isc_no, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.ItemStChDetResponse
	for _, detail := range Details {
		var detailData entity.ItemStChDetResponse
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

	iscDate := itemStCh.IscDate.Format("2006-01-02")
	response.IscDate = &iscDate

	response.Details = DetailsData
	return response, nil
}

func (service *itemStChServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.ItemStChListResponse, total int64, lastPage int, err error) {
	// log.Println("service itemStChServiceImpl, List")
	itemStChanges, total, lastPage, err := service.ItemStChRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	// log.Println("itemStChanges:", structs.StructToJson(itemStChanges))
	for _, row := range itemStChanges {
		var vResp entity.ItemStChListResponse
		structs.Automapper(row, &vResp)
		if row.IscDate != nil {
			iscDate := row.IscDate.Format("2006-01-02")
			vResp.IscDate = &iscDate
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *itemStChServiceImpl) Delete(custId string, isc_no string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ItemStChRepository.Delete(txCtx, custId, isc_no, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *itemStChServiceImpl) Update(isc_no string, request entity.UpdateItemStChBody) (err error) {
	c := context.Background()

	if request.IscDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		if *request.IscDate != "" {
			deliveryDate, err := str.DateStrToRfc3339String(*request.IscDate)
			if err != nil {
				return err
			}
			request.IscDate = &deliveryDate
		}

	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.ItemStCh
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ItemStChRepository.Update(txCtx, isc_no, Model)
		if err != nil {
			return err
		}
		DetailIds := []int{}

		for _, detail := range request.Details {
			if detail.IscDetID != nil {
				DetailIds = append(DetailIds, *detail.IscDetID)
			}
		}
		if len(DetailIds) > 0 {
			err := service.ItemStChRepository.DeleteDetailNotInIDs(txCtx, isc_no, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details {
			sequence := detail.SeqNo
			// parse time format YYYY-mm-dd to Rfc3339
			if detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
				if err != nil {
					return err
				}
				detail.ExpDate = &expDate
			}
			detail.SeqNo = sequence
			detail.CustID = request.CustID
			detail.IscNo = isc_no

			var iscDetModel model.ItemStChDet
			err = structs.Automapper(detail, &iscDetModel)
			if err != nil {
				return err
			}
			if detail.IscDetID == nil || *detail.IscDetID == 0 {
				iscDetModel.IscDetId = nil
				err = service.ItemStChRepository.StoreDetail(txCtx, &iscDetModel)
				if err != nil {
					return err
				}
			} else {
				iscDetModel.CustID = ""
				err = service.ItemStChRepository.UpdateGrDetail(txCtx, &iscDetModel)
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
