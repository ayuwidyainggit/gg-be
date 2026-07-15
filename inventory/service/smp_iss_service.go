package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
)

type SmpIssService interface {
	Store(request entity.CreateSampleIssueBody) (err error)
	Detail(smpIssueNo string, custID, parentCustId string) (response entity.SampleIssueResponse, err error)
	List(dataFilter entity.GeneralQueryFilter, custId string, parentCustId string) (data []entity.SampleIssueListResponse, total int64, lastPage int, err error)
	Delete(custId string, smpIssNo string, userId int64) (err error)
	Update(smpIssNo string, request entity.UpdateSampleIssueBody) (err error)
}

type smpIssServiceImpl struct {
	SmpIssRepository  repository.SmpIssRepository
	WhStockRepository repository.WhStockRepository
	StockRepository   repository.StockRepository
	Transaction       repository.Dbtransaction
}

func NewSmpIssService(
	smpIssRepository repository.SmpIssRepository,
	whStockRepository repository.WhStockRepository,
	stockRepository repository.StockRepository,
	transaction repository.Dbtransaction) *smpIssServiceImpl {
	return &smpIssServiceImpl{
		SmpIssRepository:  smpIssRepository,
		WhStockRepository: whStockRepository,
		StockRepository:   stockRepository,
		Transaction:       transaction,
	}
}

func (service *smpIssServiceImpl) Store(request entity.CreateSampleIssueBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	smpIssDate, err := str.DateStrToRfc3339String(request.SmpIssDate)
	if err != nil {
		return err
	}
	request.SmpIssDate = smpIssDate
	var smpIssModel model.SampleIssue
	err = structs.Automapper(request, &smpIssModel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.SmpIssRepository.Store(txCtx, &smpIssModel)
		if err != nil {
			return err
		}

		// smpIssue, err := service.SmpIssRepository.FindByNo(smpIssNo, request.CustID, request.ParentCustID)
		// if err != nil {
		// 	return err
		// }

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
			Detail.SmpIssNo = smpIssModel.SmpIssNo
			var smpIssDetModel model.SampleIssueDet
			err = structs.Automapper(Detail, &smpIssDetModel)
			if err != nil {
				return err
			}

			err = service.SmpIssRepository.StoreDetail(txCtx, &smpIssDetModel)
			if err != nil {
				return err
			}
			/*
				whStockQuery := entity.WhStockQuery{
					CustID: smpIssModel.CustID,
					WhId:   *smpIssModel.WhID,
					ProId:  int64(smpIssDetModel.ProID),
				}

				whStock, err := service.WhStockRepository.FindByWhIdAndProId(whStockQuery)
				if err != nil {
					return err
				}

				newWhStock := model.WhStock{}
				oldQty := *whStock.Qty
				detailQty := *smpIssDetModel.Qty
				newQty := oldQty - detailQty
				newWhStock.Qty = &newQty
				if *newWhStock.Qty < 0 {
					return errors.New("quantity stock must be greater than 0")
				}
				err = service.WhStockRepository.UpdateWhStockByWhIdAndProId(txCtx, whStock.CustID, *whStock.WhID, *whStock.ProID, newWhStock)
				if err != nil {
					return err
				}

				var stock model.Stock
				err = structs.Automapper(whStock, &stock)
				if err != nil {
					return err
				}

				tempCogs := float64(0)
				itemCnd := int64(1)
				stock.Cogs = &tempCogs
				stock.ItemCnd = &itemCnd
				stock.Qty = smpIssDetModel.Qty
				stock.TrCode = request.TrCode
				stock.TrNo = smpIssModel.SmpIssNo
				stock.StockDate = smpIssModel.SmpIssDate
				stock.WhIDFrom = smpIssModel.WhID
				stock.UnitPrice = smpIssDetModel.UnitPrice1
				smpIssDetId := int64(*smpIssDetModel.SmpIssDetId)
				stock.RefDetId = &smpIssDetId
				err = service.StockRepository.StoreStock(txCtx, &stock)
				if err != nil {
					return err
				}
			*/
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (service *smpIssServiceImpl) Detail(smpIssueNo string, custID, parentCustId string) (response entity.SampleIssueResponse, err error) {
	smpIssue, err := service.SmpIssRepository.FindByNo(smpIssueNo, custID, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(smpIssue, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.SmpIssRepository.FindSmpIssuedetail(smpIssueNo, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.SampleIssueDetResp
	for _, detail := range Details {
		var detailData entity.SampleIssueDetResp
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

	smpIssueDate := smpIssue.SmpIssDate.Format("2006-01-02")
	response.SmpIssDate = smpIssueDate

	response.Details = DetailsData
	return response, nil
}

func (service *smpIssServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string, parentCustId string) (data []entity.SampleIssueListResponse, total int64, lastPage int, err error) {
	smpIssues, total, lastPage, err := service.SmpIssRepository.FindAllByCustId(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range smpIssues {
		var vResp entity.SampleIssueListResponse
		structs.Automapper(row, &vResp)
		smpIssueDate := row.SmpIssDate.Format("2006-01-02")
		vResp.SmpIssDate = smpIssueDate

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *smpIssServiceImpl) Delete(custId string, smpIssNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.SmpIssRepository.Delete(txCtx, custId, smpIssNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *smpIssServiceImpl) Update(smpIssNo string, request entity.UpdateSampleIssueBody) (err error) {
	c := context.Background()

	if request.SmpIssDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		if *request.SmpIssDate != "" {
			deliveryDate, err := str.DateStrToRfc3339String(*request.SmpIssDate)
			if err != nil {
				return err
			}
			request.SmpIssDate = &deliveryDate
		}

	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.SampleIssue
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	// custID := Model.CustID
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.SmpIssRepository.Update(txCtx, smpIssNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int{}

		for _, detail := range request.Details {
			if detail.SmpIssDetId != nil {
				DetailIds = append(DetailIds, *detail.SmpIssDetId)
			}
		}
		if len(DetailIds) > 0 {
			err := service.SmpIssRepository.DeleteDetailNotInIDs(txCtx, smpIssNo, DetailIds)
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
			detail.SmpIssNo = smpIssNo
			var smpIssDetModel model.SampleIssueDet
			err = structs.Automapper(detail, &smpIssDetModel)
			if err != nil {
				return err
			}

			if detail.SmpIssDetId == nil || *detail.SmpIssDetId == 0 {
				smpIssDetModel.SmpIssDetId = nil
				err = service.SmpIssRepository.StoreDetail(txCtx, &smpIssDetModel)
				if err != nil {
					return err
				}
			} else {
				smpIssDetModel.CustID = ""
				err = service.SmpIssRepository.UpdateGrDetail(txCtx, &smpIssDetModel)
				if err != nil {
					return err
				}

			}
			/*
				whStockQuery := entity.WhStockQuery{
					CustID: custID,
					WhId:   *Model.WhID,
					ProId:  int64(smpIssDetModel.ProID),
				}

				whStock, err := service.WhStockRepository.FindByWhIdAndProId(whStockQuery)
				if err != nil {
					return err
				}

				newWhStock := model.WhStock{}
				oldQty := *whStock.Qty
				detailQty := *smpIssDetModel.Qty
				newQty := oldQty - detailQty
				newWhStock.Qty = &newQty
				if *newWhStock.Qty < 0 {
					return errors.New("quantity stock must be greater than 0")
				}

				err = service.WhStockRepository.UpdateWhStockByWhIdAndProId(txCtx, whStock.CustID, *whStock.WhID, *whStock.ProID, newWhStock)
				if err != nil {
					return err
				}

				var stock model.Stock
				err = structs.Automapper(whStock, &stock)
				if err != nil {
					return err
				}

				tempCogs := float64(0)
				itemCnd := int64(1)
				stock.Cogs = &tempCogs
				stock.ItemCnd = &itemCnd
				stock.Qty = smpIssDetModel.Qty
				stock.TrCode = *Model.TrCode
				stock.TrNo = smpIssNo
				stock.StockDate = Model.SmpIssDate
				stock.WhIDFrom = Model.WhID
				stock.UnitPrice = smpIssDetModel.UnitPrice1
				smpIssDetId := int64(*smpIssDetModel.SmpIssDetId)
				stock.RefDetId = &smpIssDetId
				err = service.StockRepository.StoreStock(txCtx, &stock)
				if err != nil {
					return err
				}
			*/
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
