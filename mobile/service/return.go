package service

import (
	"context"
	"encoding/json"
	"fmt"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/config/env"
	"mobile/pkg/str"
	"mobile/pkg/structs"
	"mobile/repository"
	"sort"
	"time"
)

type ReturnService interface {
	ReturnReasons(entity.ReturnReasonsRequest) ([]string, error)
	// ReturnList(entity.ReturnReasonsRequest) ([]entity.ReturnListResponse, error)
	ReturnReasonLookupList(entity.GeneralQueryFilter) (data []entity.ReturnReasonsLookupResponse, total int64, lastPage int, err error)
	Store(request entity.CreateReturnBody) (err error)
	UpdateStatus(request entity.UpdateStatusReturnBody) (err error)
	UpdateQuantity(returnNo string, request entity.UpdateQuantityReturnBody) (err error)
}

func NewReturnService(config env.ConfigEnv, returnRepository repository.ReturnRepository, transaction repository.Dbtransaction) *returnServiceImpl {
	return &returnServiceImpl{
		Config:           config,
		ReturnRepository: returnRepository,
		Transaction:      transaction,
	}
}

type returnServiceImpl struct {
	Config           env.ConfigEnv
	ReturnRepository repository.ReturnRepository
	Transaction      repository.Dbtransaction
}

func (service *returnServiceImpl) ReturnReasons(request entity.ReturnReasonsRequest) (response []string, err error) {

	response = []string{"reason 1", "reason 2", "reason 3"}
	return response, err
}

// func (service *ReturnsServiceImpl) ReturnList(request entity.ReturnReasonsRequest) (response []entity.ReturnListResponse, err error) {

// 	for i := 1; i < 4; i++ {
// 		rowReturnList := entity.ReturnListResponse{
// 			OutletCode: "00" + strconv.Itoa(i),
// 			OutletName: "Toko makmur jaya",
// 			ReturnNo:   "100000100",
// 			ReturnDate: time.Now().Format("2006-01-02"),
// 			InvoiceNo:  "001",
// 			ProCode:    "001",
// 			ProName:    "Sabun Lifeboy",
// 			UnitId1:    "PCS",
// 			UnitId2:    "BOX",
// 			UnitId3:    "CTN",
// 			Qty:        5,
// 			Price:      500000,
// 		}
// 		response = append(response, rowReturnList)

// 	}

// 	return response, err
// }

func (service *returnServiceImpl) Store(request entity.CreateReturnBody) (err error) {
	c := context.Background()

	returnDetails := request.Details

	sort.Slice(returnDetails[:], func(i, j int) bool {
		if returnDetails[i].InvoiceNo != nil && returnDetails[j].InvoiceNo != nil {
			return *returnDetails[i].InvoiceNo < *returnDetails[j].InvoiceNo
		}

		if returnDetails[i].InvoiceNo != nil && returnDetails[j].InvoiceNo == nil {
			return false
		}

		return true
	})
	/*
		for i, Detail := range returnDetails {
			jsonF, _ := json.Marshal(Detail)

			// typecasting byte array to string
			fmt.Println(i, " : ", string(jsonF))
		}
	*/

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		// ReturnModel := model.Return{}
		var ReturnModel model.Return
		if err = structs.Automapper(request, &ReturnModel); err != nil {
			return err
		}

		// var ReturnModelList []model.Return
		var ReturnDetailModelList []model.ReturnDetail

		var invoiceNo string
		var nextInvoiceNo string
		for index, Detail := range returnDetails {
			data, err := service.ReturnRepository.FindInvoiceNoByProductId(Detail.ProductID)
			if err != nil {
				return err
			}
			if len(data) > 0 {
				Detail.InvoiceNo = &data[0].InvoiceNo
				parsedInvoiceDate, err := time.Parse(time.RFC3339, data[0].InvoiceDate)
				if err != nil {
					return err
				}

				formattedDate := parsedInvoiceDate.Format("2006-01-02")
				Detail.InvoiceDate = &formattedDate

			} else {
				Detail.InvoiceNo = nil
				Detail.InvoiceDate = nil
			}

			if Detail.InvoiceNo == nil {
				ReturnModel.DataStatus = 1
			} else {
				ReturnModel.DataStatus = 3
				if Detail.InvoiceDate != nil {
					invoiceDate, err := str.DateStrToRfc3339String(*Detail.InvoiceDate)
					if err != nil {
						return err
					}
					Detail.InvoiceDate = &invoiceDate
				}

				if index < len(returnDetails)-1 {
					resInvoiceNo, _ := json.Marshal(Detail.InvoiceNo)
					invoiceNo = string(resInvoiceNo)
					resNextInvoiceNo, _ := json.Marshal(returnDetails[index+1].InvoiceNo)
					nextInvoiceNo = string(resNextInvoiceNo)
				}
			}

			fmt.Println("&data[0].InvoiceNo >>", data[0].InvoiceNo)

			ReturnModel.ReturnDate = time.Now()
			ReturnModel.InvoiceNo = Detail.InvoiceNo
			ReturnModel.SalesmanID = *Detail.SalesmanID
			ReturnModel.OutletID = *Detail.OutletID
			ReturnModel.RefferenceNo = *Detail.RefferenceNo

			if Detail.InvoiceDate != nil {
				parsedInvoiceDate, err := time.Parse(time.RFC3339, *Detail.InvoiceDate)
				if err != nil {
					return err
				}
				ReturnModel.InvoiceDate = &parsedInvoiceDate
			}

			var ReturnDetailModel model.ReturnDetail
			if err = structs.Automapper(Detail, &ReturnDetailModel); err != nil {
				return err
			}

			ReturnDetailModel.CustID = ReturnModel.CustID
			ReturnDetailModel.SubTotal = (ReturnDetailModel.Qty1 * ReturnDetailModel.SellPrice1) + (ReturnDetailModel.Qty2 * ReturnDetailModel.SellPrice2) + (ReturnDetailModel.Qty3 * ReturnDetailModel.SellPrice3)
			ReturnDetailModel.VatValue = ReturnDetailModel.SubTotal * (ReturnDetailModel.Vat / 100.0)
			ReturnDetailModel.Total = ReturnDetailModel.SubTotal - ReturnDetailModel.VatValue
			ReturnModel.SubTotal += ReturnDetailModel.SubTotal
			ReturnModel.VatValue += ReturnDetailModel.VatValue

			ReturnDetailModel.OrderDetailID = &data[0].OrderDetailID

			ReturnDetailModelList = append(ReturnDetailModelList, ReturnDetailModel)

			if Detail.InvoiceNo == nil || index == len(returnDetails)-1 || (index < len(returnDetails)-1 && invoiceNo != nextInvoiceNo) {
				tmpWhId := Detail.WhID
				if data[0].InvoiceNo != "" {
					tmpWhId = &data[0].WhID
				} else {
					dataWhId, err := service.ReturnRepository.FindOneWhIdBySalesmanID(*Detail.SalesmanID)
					if err != nil {
						return err
					}
					tmpWhId = &dataWhId
				}

				ReturnModel.Total = ReturnModel.SubTotal - ReturnModel.VatValue
				if err := service.ReturnRepository.Store(txCtx, &ReturnModel); err != nil {
					return err
				}

				// ReturnModelList = append(ReturnModelList, ReturnModel)
				for _, returnDetail := range ReturnDetailModelList {
					returnDetail.ReturnNo = ReturnModel.ReturnNo
					returnDetail.WhId = *tmpWhId
					if err = service.ReturnRepository.StoreDetail(txCtx, &returnDetail); err != nil {
						return err
					}
				}

				// var ReturnModel model.Return
				ReturnModel = model.Return{}
				if err = structs.Automapper(request, &ReturnModel); err != nil {
					return err
				}

				ReturnDetailModelList = []model.ReturnDetail{}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (service *returnServiceImpl) UpdateStatus(request entity.UpdateStatusReturnBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		var Model model.Return
		err = structs.Automapper(request, &Model)
		if err != nil {
			return err
		}
		Model.CustID = ""

		for _, detail := range request.Returns {

			Model.ReturnNo = detail.ReturnNo
			Model.DataStatus = detail.DataStatus

			err = service.ReturnRepository.Update(txCtx, &Model)
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

func (service *returnServiceImpl) ReturnReasonLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.ReturnReasonsLookupResponse, total int64, lastPage int, err error) {
	Outlets, total, lastPage, err := service.ReturnRepository.FindAllMasterReturnReasonLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Outlets {
		var vResp entity.ReturnReasonsLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnServiceImpl) UpdateQuantity(returnNo string, request entity.UpdateQuantityReturnBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		DetailIds := []int64{}

		for _, detail := range request.Details {
			DetailIds = append(DetailIds, detail.ReturnDetailID)
		}
		if len(DetailIds) > 0 {
			err := service.ReturnRepository.DeleteDetailNotInIDs(txCtx, returnNo, DetailIds)
			if err != nil {
				return err
			}
		}

		var Model model.Return
		err = structs.Automapper(request, &Model)
		if err != nil {
			return err
		}

		for _, detail := range request.Details {

			var returnDetailModel model.ReturnQuantity
			err = structs.Automapper(detail, &returnDetailModel)
			if err != nil {
				return err
			}

			returnDetailModel.SubTotal = (returnDetailModel.Qty1 * returnDetailModel.SellPrice1) + (returnDetailModel.Qty2 * returnDetailModel.SellPrice2) + (returnDetailModel.Qty3 * returnDetailModel.SellPrice3)
			returnDetailModel.VatValue = returnDetailModel.Total * (returnDetailModel.Vat / 100.0)
			returnDetailModel.Total = returnDetailModel.SubTotal - returnDetailModel.VatValue
			Model.SubTotal += returnDetailModel.SubTotal
			Model.VatValue += returnDetailModel.VatValue

			returnDetailModel.CustID = ""
			err = service.ReturnRepository.UpdateQuantity(txCtx, &returnDetailModel)
			if err != nil {
				return err
			}
		}

		Model.CustID = ""
		Model.ReturnNo = returnNo
		Model.Total = Model.SubTotal - Model.VatValue
		err = service.ReturnRepository.Update(txCtx, &Model)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
