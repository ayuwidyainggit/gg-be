package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"finance/pkg/structs"
	"finance/repository"
	"time"
)

type ApService interface {
	Store(request entity.CreateApBody) (err error)
	Detail(arNo string, custID string) (response entity.ApResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.ApListResponse, total int64, lastPage int, err error)
	Delete(custId string, apNo string, userId int64) (err error)
	Update(apNo string, request entity.UpdateApBody) (err error)
}

type apServiceImpl struct {
	Repository  repository.ApRepository
	Transaction repository.Dbtransaction
}

func NewApService(repository repository.ApRepository, transaction repository.Dbtransaction) *apServiceImpl {
	return &apServiceImpl{
		Repository:  repository,
		Transaction: transaction,
	}
}

func (service *apServiceImpl) Store(request entity.CreateApBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.ApDate != nil {
		apDate, err := str.DateStrToRfc3339String(*request.ApDate)
		if err != nil {
			return err
		}
		request.ApDate = &apDate
	}

	// parse time format YYYY-mm-dd to Rfc3339
	if request.InvDate != nil {
		InvDate, err := str.DateStrToRfc3339String(*request.InvDate)
		if err != nil {
			return err
		}
		request.InvDate = &InvDate
	}

	// parse time format YYYY-mm-dd to Rfc3339
	if request.InvDueDate != nil {
		InvDueDate, err := str.DateStrToRfc3339String(*request.InvDueDate)
		if err != nil {
			return err
		}
		request.InvDueDate = &InvDueDate
	}

	// parse time format YYYY-mm-dd to Rfc3339
	if request.TaxInvDate != nil {
		TaxInvDate, err := str.DateStrToRfc3339String(*request.TaxInvDate)
		if err != nil {
			return err
		}
		request.TaxInvDate = &TaxInvDate
	}

	// parse time format YYYY-mm-dd to Rfc3339
	if request.TaxReturnDate != nil {
		TaxReturnDate, err := str.DateStrToRfc3339String(*request.TaxReturnDate)
		if err != nil {
			return err
		}
		request.TaxReturnDate = &TaxReturnDate
	}

	var Apmodel model.Ap
	err = structs.Automapper(request, &Apmodel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.Store(txCtx, &Apmodel)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details {
			var detModel model.ApDet
			err = structs.Automapper(Detail, &detModel)
			if err != nil {
				return err
			}
			detModel.CustID = request.CustID
			detModel.ApNo = Apmodel.ApNo
			err = service.Repository.StoreDetail(txCtx, &detModel)
			if err != nil {
				return err
			}
		}

		for _, MoneyPromoDetail := range request.MoneyPromoDetails {
			var apMoneyPromoModel model.ApMoneyPromo
			err = structs.Automapper(MoneyPromoDetail, &apMoneyPromoModel)
			if err != nil {
				return err
			}
			apMoneyPromoModel.CustID = request.CustID
			apMoneyPromoModel.ApNo = Apmodel.ApNo
			err = service.Repository.StoreMoneyPromoDetail(txCtx, &apMoneyPromoModel)
			if err != nil {
				return err
			}
		}

		for _, qtyPromoDetail := range request.QtyPromoDetails {
			var apQtyPromoModel model.ApQtyPromo
			err = structs.Automapper(qtyPromoDetail, &apQtyPromoModel)
			if err != nil {
				return err
			}
			apQtyPromoModel.CustID = request.CustID
			apQtyPromoModel.ApNo = Apmodel.ApNo
			err = service.Repository.StoreQtyPromoDetail(txCtx, &apQtyPromoModel)
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

func (service *apServiceImpl) Detail(apNo string, custID string) (response entity.ApResponse, err error) {
	ap, err := service.Repository.FindByNo(apNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ap, &response)
	if err != nil {
		return response, err
	}

	details, err := service.Repository.FindDetail(apNo, custID)
	if err != nil {
		return response, err
	}
	var detailsData []entity.ApDetResponse
	for _, detail := range details {
		var detailData entity.ApDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		detailsData = append(detailsData, detailData)
	}

	qtyPromoDetails, err := service.Repository.FindQtyPromoDetail(apNo, custID)
	if err != nil {
		return response, err
	}

	var qtyPromoData []entity.ApQtyPromoResponse
	for _, qtyPromoDetail := range qtyPromoDetails {
		var detailData entity.ApQtyPromoResponse
		err = structs.Automapper(qtyPromoDetail, &detailData)
		if err != nil {
			return response, err
		}
		qtyPromoData = append(qtyPromoData, detailData)
	}

	moneyPromoDetails, err := service.Repository.FindMoneyPromoDetail(apNo, custID)
	if err != nil {
		return response, err
	}

	var moneyPromoData []entity.ApMoneyPromoResponse
	for _, moneyPromoDetail := range moneyPromoDetails {
		var detailData entity.ApMoneyPromoResponse
		err = structs.Automapper(moneyPromoDetail, &detailData)
		if err != nil {
			return response, err
		}

		moneyPromoData = append(moneyPromoData, detailData)
	}

	apDate := ap.ApDate.Format("2006-01-02")
	response.ApDate = &apDate

	invDate := ap.InvDate.Format("2006-01-02")
	response.InvDate = &invDate

	invDueDate := ap.InvDueDate.Format("2006-01-02")
	response.InvDueDate = &invDueDate

	taxInvDate := ap.TaxInvDate.Format("2006-01-02")
	response.TaxInvDate = &taxInvDate

	taxReturnDate := ap.TaxReturnDate.Format("2006-01-02")
	response.TaxReturnDate = &taxReturnDate

	response.Details = detailsData
	response.QtyPromoDetails = qtyPromoData
	response.MoneyPromoDetails = moneyPromoData
	return response, nil
}

func (service *apServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.ApListResponse, total int64, lastPage int, err error) {
	aps, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range aps {
		var vResp entity.ApListResponse
		structs.Automapper(row, &vResp)
		if row.ApDate != nil {
			ApDate := row.ApDate.Format("2006-01-02")
			vResp.ApDate = &ApDate
		}
		if row.InvDate != nil {
			InvDate := row.InvDate.Format("2006-01-02")
			vResp.InvDate = &InvDate
		}
		if row.InvDueDate != nil {
			InvDueDate := row.InvDueDate.Format("2006-01-02")
			vResp.InvDueDate = &InvDueDate
		}
		if row.TaxInvDate != nil {
			TaxInvDate := row.TaxInvDate.Format("2006-01-02")
			vResp.TaxInvDate = &TaxInvDate
		}
		if row.TaxReturnDate != nil {
			TaxReturnDate := row.TaxReturnDate.Format("2006-01-02")
			vResp.TaxReturnDate = &TaxReturnDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *apServiceImpl) Delete(custId string, apNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, apNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (service *apServiceImpl) Update(apNo string, request entity.UpdateApBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.ApDate != nil {
		apDate, err := str.DateStrToRfc3339String(*request.ApDate)
		if err != nil {
			return err
		}
		request.ApDate = &apDate
	}

	// parse time format YYYY-mm-dd to Rfc3339
	if request.InvDate != nil {
		InvDate, err := str.DateStrToRfc3339String(*request.InvDate)
		if err != nil {
			return err
		}
		request.InvDate = &InvDate
	}

	// parse time format YYYY-mm-dd to Rfc3339
	if request.InvDueDate != nil {
		InvDueDate, err := str.DateStrToRfc3339String(*request.InvDueDate)
		if err != nil {
			return err
		}
		request.InvDueDate = &InvDueDate
	}

	// parse time format YYYY-mm-dd to Rfc3339
	if request.TaxInvDate != nil {
		TaxInvDate, err := str.DateStrToRfc3339String(*request.TaxInvDate)
		if err != nil {
			return err
		}
		request.TaxInvDate = &TaxInvDate
	}

	// parse time format YYYY-mm-dd to Rfc3339
	if request.TaxReturnDate != nil {
		TaxReturnDate, err := str.DateStrToRfc3339String(*request.TaxReturnDate)
		if err != nil {
			return err
		}
		request.TaxReturnDate = &TaxReturnDate
	}
	// End parse time format YYYY-mm-dd to Rfc339

	var Model model.Ap
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	if Model.IsPosted != nil {
		if *Model.IsPosted {
			now := time.Now()
			Model.PostedAt = &now
		}
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Update(txCtx, apNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}
		qtyPromoDetailIds := []int64{}
		moneyPromoDetailIds := []int64{}

		for _, detail := range request.Details {
			if detail.ApDetID != nil {
				DetailIds = append(DetailIds, *detail.ApDetID)
			}
		}
		for _, detail := range request.QtyPromoDetails {
			if detail.ApQtyPromoID != nil {
				qtyPromoDetailIds = append(qtyPromoDetailIds, *detail.ApQtyPromoID)
			}
		}
		for _, detail := range request.MoneyPromoDetails {
			if detail.ApMoneyPromoID != nil {
				moneyPromoDetailIds = append(moneyPromoDetailIds, *detail.ApMoneyPromoID)
			}
		}
		if len(DetailIds) > 0 {
			err := service.Repository.DeleteDetailNotInIDs(txCtx, apNo, DetailIds)
			if err != nil {
				return err
			}
		}
		if len(qtyPromoDetailIds) > 0 {
			err := service.Repository.DeleteQtyPromoDetailNotInIDs(txCtx, apNo, qtyPromoDetailIds)
			if err != nil {
				return err
			}
		}
		if len(moneyPromoDetailIds) > 0 {
			err := service.Repository.DeleteMoneyPromoDetailNotInIDs(txCtx, apNo, moneyPromoDetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details {
			var apDetModel model.ApDet

			err = structs.Automapper(detail, &apDetModel)
			if err != nil {
				return err
			}
			apDetModel.CustID = request.CustID
			apDetModel.ApNo = apNo
			if detail.ApDetID == nil || *detail.ApDetID == 0 {
				apDetModel.ApDetID = 0
				err = service.Repository.StoreDetail(txCtx, &apDetModel)
				if err != nil {
					return err
				}
			} else {
				apDetModel.CustID = ""
				err = service.Repository.UpdateDetail(txCtx, &apDetModel)
				if err != nil {
					return err
				}

			}
		}

		for _, detail := range request.QtyPromoDetails {
			var apQtyPromoDetModel model.ApQtyPromo

			err = structs.Automapper(detail, &apQtyPromoDetModel)
			if err != nil {
				return err
			}
			apQtyPromoDetModel.CustID = request.CustID
			apQtyPromoDetModel.ApNo = apNo
			if detail.ApQtyPromoID == nil || *detail.ApQtyPromoID == 0 {
				apQtyPromoDetModel.ApQtyPromoID = nil
				err = service.Repository.StoreQtyPromoDetail(txCtx, &apQtyPromoDetModel)
				if err != nil {
					return err
				}
			} else {
				apQtyPromoDetModel.CustID = ""
				err = service.Repository.UpdateQtyPromoDetail(txCtx, &apQtyPromoDetModel)
				if err != nil {
					return err
				}

			}
		}

		for _, detail := range request.MoneyPromoDetails {
			var apMoneyPromoDetModel model.ApMoneyPromo

			err = structs.Automapper(detail, &apMoneyPromoDetModel)
			if err != nil {
				return err
			}
			apMoneyPromoDetModel.CustID = request.CustID
			apMoneyPromoDetModel.ApNo = apNo
			if detail.ApMoneyPromoID == nil {
				err = service.Repository.StoreMoneyPromoDetail(txCtx, &apMoneyPromoDetModel)
				if err != nil {
					return err
				}
			} else {
				apMoneyPromoDetModel.CustID = ""
				err = service.Repository.UpdateMoneyPromoDetail(txCtx, &apMoneyPromoDetModel)
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
