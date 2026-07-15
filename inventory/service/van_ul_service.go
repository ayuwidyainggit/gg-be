package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
)

type VanUlService interface {
	Store(request entity.CreateVanUlBody) (err error)
	Detail(vanUlNo string, custID string) (response entity.VanUlResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.VanUlListResponse, total int64, lastPage int, err error)
	Delete(custId string, vanBsUlNo string, userId int64) (err error)
	Update(vanUlNo string, request entity.VanUlUpdateBody) (err error)
}

type vanUlServiceImpl struct {
	Repository  repository.VanUlRepository
	Transaction repository.Dbtransaction
}

func NewVanUlService(Repository repository.VanUlRepository, transaction repository.Dbtransaction) *vanUlServiceImpl {
	return &vanUlServiceImpl{
		Repository:  Repository,
		Transaction: transaction,
	}
}
func (service *vanUlServiceImpl) Store(request entity.CreateVanUlBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.VanUlDate != nil {
		vanUlDate, err := str.DateStrToRfc3339String(*request.VanUlDate)
		if err != nil {
			return err
		}
		request.VanUlDate = &vanUlDate
	}

	var Model model.VanUl
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.Store(txCtx, &Model)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details.Normal {
			// parse time format YYYY-mm-dd to Rfc3339
			if Detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*Detail.ExpDate)
				if err != nil {
					return err
				}
				Detail.ExpDate = &expDate
			}
			var vanUlDetModel model.VanUlDet

			err = structs.Automapper(Detail, &vanUlDetModel)
			if err != nil {
				return err
			}
			vanUlDetModel.CustID = request.CustID
			vanUlDetModel.VanUlNo = Model.VanUlNo
			vanUlDetModel.ItemType = 1
			err = service.Repository.StoreDetail(txCtx, &vanUlDetModel)
			if err != nil {
				return err
			}

		}
		for _, Detail := range request.Details.Promo {
			// parse time format YYYY-mm-dd to Rfc3339
			if Detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*Detail.ExpDate)
				if err != nil {
					return err
				}
				Detail.ExpDate = &expDate
			}
			var vanUlDetModel model.VanUlDet

			err = structs.Automapper(Detail, &vanUlDetModel)
			if err != nil {
				return err
			}
			vanUlDetModel.CustID = request.CustID
			vanUlDetModel.VanUlNo = Model.VanUlNo
			vanUlDetModel.ItemType = 2
			err = service.Repository.StoreDetail(txCtx, &vanUlDetModel)
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
func (service *vanUlServiceImpl) Detail(vanUlNo string, custID string) (response entity.VanUlResponse, err error) {
	vanUl, err := service.Repository.FindByNo(vanUlNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(vanUl, &response)
	if err != nil {
		return response, err
	}

	vanUlDetails, err := service.Repository.FindDetail(vanUlNo, custID)
	if err != nil {
		return response, err
	}
	for _, Detail := range vanUlDetails {
		var vanUlDetailData entity.VanUlDetReadResponse
		err = structs.Automapper(Detail, &vanUlDetailData)
		if err != nil {
			return response, err
		}

		if vanUlDetailData.ExpDate != nil {
			expDate := Detail.ExpDate.Format("2006-01-02")
			vanUlDetailData.ExpDate = &expDate
		}

		if vanUlDetailData.ItemType == 1 {
			response.Details.Normal = append(response.Details.Normal, vanUlDetailData)
		} else {
			response.Details.Promo = append(response.Details.Promo, vanUlDetailData)

		}
	}
	vanUlData := vanUl.VanUlDate.Format("2006-01-02")

	response.VanUlDate = vanUlData
	return response, nil
}

func (service *vanUlServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.VanUlListResponse, total int64, lastPage int, err error) {
	grs, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range grs {
		var vResp entity.VanUlListResponse
		structs.Automapper(row, &vResp)
		grData := row.VanUlDate.Format("2006-01-02")
		vResp.VanUlDate = grData

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *vanUlServiceImpl) Update(vanUlNo string, request entity.VanUlUpdateBody) (err error) {
	c := context.Background()

	if request.VanUlDate != nil {
		if *request.VanUlDate != "" {
			Date, err := str.DateStrToRfc3339String(*request.VanUlDate)
			if err != nil {
				return err
			}
			request.VanUlDate = &Date
		}
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var vanUlModel model.VanUl
	err = structs.Automapper(request, &vanUlModel)
	if err != nil {
		return err
	}
	vanUlModel.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Update(txCtx, vanUlNo, vanUlModel)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details.Normal {
			if detail.VanUlDetID != nil {
				DetailIds = append(DetailIds, *detail.VanUlDetID)
			}
		}
		for _, detail := range request.Details.Promo {
			if detail.VanUlDetID != nil {
				DetailIds = append(DetailIds, *detail.VanUlDetID)
			}
		}
		if len(DetailIds) > 0 {
			// log.Println("grDetailIds:", structs.StructToJson(grDetailIds))
			err := service.Repository.DeleteDetailNotInIDs(txCtx, vanUlNo, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details.Normal {
			// sequence := detail.SeqNo
			if detail.VanUlDetID == nil || *detail.VanUlDetID == 0 {
				detail.VanUlDetID = nil
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
				var DetailModel model.VanUlDet
				err = structs.Automapper(detail, &DetailModel)
				if err != nil {
					return err
				}
				// grDetailModel.SeqNo = sequence
				DetailModel.CustID = request.CustID
				DetailModel.VanUlNo = vanUlNo
				DetailModel.ItemType = 1
				DetailModel.CustID = ""
				err = service.Repository.StoreDetail(txCtx, &DetailModel)
				if err != nil {
					return err
				}
			} else {
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
				var DetailModel model.VanUlDet
				err = structs.Automapper(detail, &DetailModel)
				if err != nil {
					return err
				}
				DetailModel.ItemType = 1
				DetailModel.CustID = ""
				err = service.Repository.UpdateDetail(txCtx, &DetailModel)
				if err != nil {
					return err
				}
			}
		}

		for _, detail := range request.Details.Promo {
			// sequence := detail.SeqNo
			if detail.VanUlDetID == nil || *detail.VanUlDetID == 0 {
				detail.VanUlDetID = nil
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
				var DetailModel model.VanUlDet
				err = structs.Automapper(detail, &DetailModel)
				if err != nil {
					return err
				}
				// grDetailModel.SeqNo = sequence
				DetailModel.CustID = request.CustID
				DetailModel.VanUlNo = vanUlNo
				DetailModel.ItemType = 2
				err = service.Repository.StoreDetail(txCtx, &DetailModel)
				if err != nil {
					return err
				}
			} else {
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
				var DetailModel model.VanUlDet
				err = structs.Automapper(detail, &DetailModel)
				if err != nil {
					return err
				}
				DetailModel.ItemType = 2
				DetailModel.CustID = ""
				err = service.Repository.UpdateDetail(txCtx, &DetailModel)
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
	// grDetails, err := service.GrRepository.FindGrdetailDetails(grNo, request.CustID)
	// if err != nil {
	// 	return err
	// }

	return nil
}
func (service *vanUlServiceImpl) Delete(custId string, vanUlNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, vanUlNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
