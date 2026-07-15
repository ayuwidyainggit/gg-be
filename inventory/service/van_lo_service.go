package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
)

type VanLoService interface {
	Store(request entity.CreateVanLoBody) (err error)
	Detail(vanLoNo string, custID string) (response entity.VanLoResponse, err error)
	List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.VanLoListResponse, total int64, lastPage int, err error)
	Delete(custId string, vanLoNo string, userId int64) (err error)
	Update(vanLoNo string, request entity.VanLoUpdateBody) (err error)
}

type vanLoServiceImpl struct {
	Repository  repository.VanLoRepository
	Transaction repository.Dbtransaction
}

func NewVanLoService(Repository repository.VanLoRepository, transaction repository.Dbtransaction) *vanLoServiceImpl {
	return &vanLoServiceImpl{
		Repository:  Repository,
		Transaction: transaction,
	}
}
func (service *vanLoServiceImpl) Store(request entity.CreateVanLoBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.VanLoDate != nil {
		vanLoDate, err := str.DateStrToRfc3339String(*request.VanLoDate)
		if err != nil {
			return err
		}
		request.VanLoDate = &vanLoDate
	}

	var Model model.VanLo
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
			var vanLoDetModel model.VanLoDet

			err = structs.Automapper(Detail, &vanLoDetModel)
			if err != nil {
				return err
			}
			vanLoDetModel.CustID = request.CustID
			vanLoDetModel.VanLoNo = Model.VanLoNo
			vanLoDetModel.ItemType = 1
			err = service.Repository.StoreDetail(txCtx, &vanLoDetModel)
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
			var vanLoDetModel model.VanLoDet

			err = structs.Automapper(Detail, &vanLoDetModel)
			if err != nil {
				return err
			}
			vanLoDetModel.CustID = request.CustID
			vanLoDetModel.VanLoNo = Model.VanLoNo
			vanLoDetModel.ItemType = 2
			err = service.Repository.StoreDetail(txCtx, &vanLoDetModel)
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
func (service *vanLoServiceImpl) Detail(vanLoNo string, custID string) (response entity.VanLoResponse, err error) {
	vanLo, err := service.Repository.FindByNo(vanLoNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(vanLo, &response)
	if err != nil {
		return response, err
	}

	vanLoDetails, err := service.Repository.FindDetail(vanLoNo, custID)
	if err != nil {
		return response, err
	}
	for _, Detail := range vanLoDetails {
		var vanLoDetailData entity.VanLoDetReadResponse
		err = structs.Automapper(Detail, &vanLoDetailData)
		if err != nil {
			return response, err
		}
		if Detail.ExpDate != nil {
			expDate := Detail.ExpDate.Format("2006-01-02")
			vanLoDetailData.ExpDate = &expDate
		}

		if vanLoDetailData.ItemType == 1 {
			response.Details.Normal = append(response.Details.Normal, vanLoDetailData)
		} else {
			response.Details.Promo = append(response.Details.Promo, vanLoDetailData)
		}
	}
	vanLoDate := vanLo.VanLoDate.Format("2006-01-02")

	response.VanLoDate = vanLoDate
	return response, nil
}
func (service *vanLoServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.VanLoListResponse, total int64, lastPage int, err error) {
	grs, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}
	if len(grs) > 0 {
		for _, row := range grs {
			var vResp entity.VanLoListResponse
			structs.Automapper(row, &vResp)
			if row.VanLoDate != nil {
				vanLoDate := row.VanLoDate.Format("2006-01-02")
				vResp.VanLoDate = &vanLoDate
			}

			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}
func (service *vanLoServiceImpl) Update(vanLoNo string, request entity.VanLoUpdateBody) (err error) {
	c := context.Background()

	if request.VanLoDate != nil {
		if *request.VanLoDate != "" {
			Date, err := str.DateStrToRfc3339String(*request.VanLoDate)
			if err != nil {
				return err
			}
			request.VanLoDate = &Date
		}
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var vanLoModel model.VanLo
	err = structs.Automapper(request, &vanLoModel)
	if err != nil {
		return err
	}
	vanLoModel.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Update(txCtx, vanLoNo, vanLoModel)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details.Normal {
			if detail.VanLoDetID != nil {
				DetailIds = append(DetailIds, *detail.VanLoDetID)
			}
		}
		for _, detail := range request.Details.Promo {
			if detail.VanLoDetID != nil {
				DetailIds = append(DetailIds, *detail.VanLoDetID)
			}
		}
		if len(DetailIds) > 0 {
			// log.Println("grDetailIds:", structs.StructToJson(grDetailIds))
			err := service.Repository.DeleteDetailNotInIDs(txCtx, vanLoNo, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details.Normal {
			// sequence := detail.SeqNo
			if detail.VanLoDetID == nil || *detail.VanLoDetID == 0 {
				detail.VanLoDetID = nil
				if detail.ExpDate != nil {
					if *detail.ExpDate != "" {
						expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
						if err != nil {
							return err
						}
						detail.ExpDate = &expDate
					}
				}

				var DetailModel model.VanLoDet
				err = structs.Automapper(detail, &DetailModel)
				if err != nil {
					return err
				}

				// grDetailModel.SeqNo = sequence
				DetailModel.CustID = request.CustID
				DetailModel.VanLoNo = vanLoNo
				DetailModel.ItemType = 1
				DetailModel.CustID = ""
				err = service.Repository.StoreDetail(txCtx, &DetailModel)
				if err != nil {
					return err
				}
			} else {
				if detail.ExpDate != nil {
					if *detail.ExpDate != "" {
						expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
						if err != nil {
							return err
						}
						detail.ExpDate = &expDate
					}
				}
				var DetailModel model.VanLoDet
				err = structs.Automapper(detail, &DetailModel)
				if err != nil {
					return err
				}
				DetailModel.CustID = ""
				DetailModel.ItemType = 1
				err = service.Repository.UpdateDetail(txCtx, &DetailModel)
				if err != nil {
					return err
				}
			}
		}

		for _, detail := range request.Details.Promo {
			// sequence := detail.SeqNo
			if detail.VanLoDetID == nil || *detail.VanLoDetID == 0 {
				detail.VanLoDetID = nil
				if detail.ExpDate != nil {
					if *detail.ExpDate != "" {
						expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
						if err != nil {
							return err
						}
						detail.ExpDate = &expDate
					}
				}

				var DetailModel model.VanLoDet
				err = structs.Automapper(detail, &DetailModel)
				if err != nil {
					return err
				}

				// grDetailModel.SeqNo = sequence
				DetailModel.CustID = request.CustID
				DetailModel.VanLoNo = vanLoNo
				DetailModel.ItemType = 2
				err = service.Repository.StoreDetail(txCtx, &DetailModel)
				if err != nil {
					return err
				}
			} else {
				if detail.ExpDate != nil {
					if *detail.ExpDate != "" {
						expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
						if err != nil {
							return err
						}
						detail.ExpDate = &expDate
					}
				}
				var DetailModel model.VanLoDet
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
func (service *vanLoServiceImpl) Delete(custId string, vanLoNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, vanLoNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
