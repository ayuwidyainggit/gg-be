package service

import (
	"context"
	"inventory/adapter"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/config/env"
	"inventory/pkg/conversion"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
	"log"
	"net/http"
	"slices"
	"strings"
)

type WhTrfService interface {
	Store(request entity.CreateWhTrfBody) (err error)
	Detail(whTrfNo string, custID, parentCustID string, distributorID int64, token string) (response entity.WhTrfResponse, err error)
	List(dataFilter entity.WhQueryFilter, custId string) (data []entity.WhTrfListResponse, total int64, lastPage int, err error)
	Delete(custId string, whTrfNo string, userId int64) (err error)
	Update(whTrfNo string, request entity.UpdateWhTrfBody) (err error)
	ListWarehouse(dataFilter entity.StockTranferWarehouseQueryFilter, custId string) (data []entity.StockTranferWarehouse, total int64, lastPage int, err error)
}

type whTrfServiceImpl struct {
	WhTrfRepository repository.WhTrfRepository
	Transaction     repository.Dbtransaction
	StockRepository repository.StockRepository
	Config          env.ConfigEnv
}

func NewWhTrfService(WhTrfRepository repository.WhTrfRepository, transaction repository.Dbtransaction, stockRepository repository.StockRepository, config env.ConfigEnv) *whTrfServiceImpl {
	return &whTrfServiceImpl{
		WhTrfRepository: WhTrfRepository,
		Transaction:     transaction,
		StockRepository: stockRepository,
		Config:          config,
	}
}
func (service *whTrfServiceImpl) Store(request entity.CreateWhTrfBody) (err error) {
	c := context.Background()

	var whTrfModel model.WhTrf
	err = structs.Automapper(request, &whTrfModel)
	if err != nil {
		return err
	}
	stockTransferDate := str.GetJakartaDate()

	whTrfModel.WhTrfDate = &stockTransferDate
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.WhTrfRepository.Store(txCtx, &whTrfModel)
		if err != nil {
			return err
		}

		var productIDs []int64
		for _, detail := range request.Details {
			productIDs = append(productIDs, detail.ProID)
		}

		productsModel, err := service.WhTrfRepository.FindProductByListID(productIDs)
		if err != nil {
			return err
		}

		var productMap = model.MapProduct{}

		for _, productModel := range productsModel {
			productMap.SetProduct(productModel.ProductId, productModel)
		}

		for index, Detail := range request.Details {
			var stockUpdateEntities []*entity.StockUpdate

			productModel, err := productMap.GetByID(Detail.ProID)
			if err != nil {
				return err
			}

			QtyUnit := &conversion.QtyUnit{
				Qty1:      int(Detail.Qty1),
				Qty2:      int(Detail.Qty2),
				Qty3:      int(Detail.Qty3),
				ConvUnit2: int(productModel.ConvUnit2),
				ConvUnit3: int(productModel.ConvUnit3),
			}

			totalQty, err := QtyUnit.ToTotalQuantity()
			if err != nil {
				return err
			}

			var whTrfDetModel model.WhTrfDet
			err = structs.Automapper(Detail, &whTrfDetModel)
			if err != nil {
				return err
			}

			whTrfDetModel.Qty = totalQty
			whTrfDetModel.SeqNo = index + 1
			whTrfDetModel.CustID = request.CustID
			whTrfDetModel.WhTrfNo = whTrfModel.WhTrfNo
			err = service.WhTrfRepository.StoreDetail(txCtx, &whTrfDetModel)
			if err != nil {
				return err
			}

			for i := 0; i < 2; i++ {
				var qtyin, qtyout float64
				var whid int64
				if i == 0 {
					qtyout = float64(whTrfDetModel.Qty)
					whid = request.WhIDFrom
				} else {
					qtyin = float64(whTrfDetModel.Qty)
					whid = request.WhIDTo
				}

				stockUpdateEntity := entity.StockUpdate{
					CustID:    request.CustID,
					WhID:      whid,
					ProID:     int64(Detail.ProID),
					StockDate: *whTrfModel.WhTrfDate,
					TrCode:    whTrfModel.WhTrfNo[0:2],
					TrNo:      whTrfModel.WhTrfNo,
					QtyIn:     qtyin,
					QtyOut:    qtyout,
					UnitPrice: 0,
					RefDetId:  int64(*whTrfDetModel.WhTrfDetId),
				}

				stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)
			}

			err = service.StockRepository.StockUpdates(txCtx, stockUpdateEntities)
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

func (service *whTrfServiceImpl) Detail(whTrfNo string, custID, parentCustID string, distributorID int64, token string) (response entity.WhTrfResponse, err error) {
	smpIssue, err := service.WhTrfRepository.FindByNo(whTrfNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(smpIssue, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.WhTrfRepository.FindWhTrfdetail(whTrfNo, custID, parentCustID, distributorID)
	if err != nil {
		return response, err
	}

	var proIDs []int64
	for _, Detail := range Details {
		if slices.Contains(proIDs, Detail.ProID) {
			continue
		}
		proIDs = append(proIDs, Detail.ProID)
	}

	masterBase := strings.TrimSpace(service.Config.Get("MASTER_SERVICE_URL"))
	masterPath := strings.TrimSpace(service.Config.Get("MASTER_SERVICE_LIST_PRODUCT_PATH"))
	masterURL := masterBase + masterPath
	masterBaseLower := strings.ToLower(masterBase)
	canCallMaster := (strings.HasPrefix(masterBaseLower, "http://") || strings.HasPrefix(masterBaseLower, "https://")) &&
		masterURL != "" && len(proIDs) > 0

	var productDistPriceMap = entity.MapProductDistPrice{}

	if canCallMaster {
		var result entity.ProductListDistPriceResponse
		client := adapter.HttpClientInfo{
			Method: http.MethodGet,
			Url:    masterURL,
			Payload: map[string]interface{}{
				"mode":            "lookup_dist_price",
				"pro_id":          proIDs,
				"distributor_id":  distributorID,
				"include_deleted": true,
				"limit":           500,
				"is_active":       9,
			},
			Auth: token,
		}

		_, err = client.Dispatch(&result)
		if err != nil {
			log.Println("error:", err)
			return response, err
		}

		for _, productPrice := range result.Data {
			productDistPriceMap.SetProduct(productPrice.ProID, productPrice)
		}
	} else {
		// Local/dev: MASTER_SERVICE_* unset or invalid — use prices from m_product join (FindWhTrfdetail).
		for _, d := range Details {
			productDistPriceMap.SetProduct(d.ProID, entity.ProductListDistPriceDataResp{
				ProID:       d.ProID,
				PurchPrice1: d.PurchPrice1,
				PurchPrice2: d.PurchPrice2,
				PurchPrice3: d.PurchPrice3,
				SellPrice1:  d.SellPrice1,
				SellPrice2:  d.SellPrice2,
				SellPrice3:  d.SellPrice3,
			})
		}
	}

	var DetailsData []entity.WhTrfDetRespose
	for _, detail := range Details {
		qty := &conversion.Qty{
			Qty:       int(detail.Qty),
			ConvUnit2: int(detail.ConvUnit2),
			ConvUnit3: int(detail.ConvUnit3),
		}

		qtyConversion := qty.ConvToQtyConversion()

		var detailData entity.WhTrfDetRespose
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		productPrice, err := productDistPriceMap.GetByID(detailData.ProID)
		if err != nil {
			return response, err
		}

		detailData.PurchPrice1 = productPrice.PurchPrice1
		detailData.PurchPrice2 = productPrice.PurchPrice2
		detailData.PurchPrice3 = productPrice.PurchPrice3
		detailData.SellPrice1 = productPrice.SellPrice1
		detailData.SellPrice2 = productPrice.SellPrice2
		detailData.SellPrice3 = productPrice.SellPrice3

		detailData.Qty1 = float64(qtyConversion.Qty1)
		detailData.Qty2 = float64(qtyConversion.Qty2)
		detailData.Qty3 = float64(qtyConversion.Qty3)

		Subtotal := (detailData.SellPrice1 * detailData.Qty1) + (detailData.SellPrice2 * detailData.Qty2) + (detailData.SellPrice3 * detailData.Qty3)
		ppn := Subtotal * detail.Vat / 100
		pbn := Subtotal * detail.VatLgPurch / 100
		ppnDp := Subtotal * detail.VatBg / 100
		total := Subtotal + ppn + pbn + ppnDp

		response.SubTotal += Subtotal
		response.Total += total
		response.TotalVat += ppn
		response.TotalVatLgPurch += pbn
		response.TotalVatBg += ppnDp

		detailData.SubTotal = Subtotal
		detailData.Total = total
		DetailsData = append(DetailsData, detailData)
	}

	whTrfDate := smpIssue.WhTrfDate.Format("2006-01-02")
	response.WhTrfDate = &whTrfDate

	response.Details = DetailsData
	return response, nil
}
func (service *whTrfServiceImpl) List(dataFilter entity.WhQueryFilter, custId string) (data []entity.WhTrfListResponse, total int64, lastPage int, err error) {
	whTrfs, total, lastPage, err := service.WhTrfRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}
	if len(whTrfs) > 0 {
		for _, row := range whTrfs {
			var vResp entity.WhTrfListResponse
			structs.Automapper(row, &vResp)
			if row.WhTrfDate != nil {
				whTrfDate := row.WhTrfDate.Format("2006-01-02")
				vResp.WhTrfDate = &whTrfDate
			}
			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}
func (service *whTrfServiceImpl) Delete(custId string, whTrfNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.WhTrfRepository.Delete(txCtx, custId, whTrfNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (service *whTrfServiceImpl) Update(whTrfNo string, request entity.UpdateWhTrfBody) (err error) {
	c := context.Background()

	if request.WhTrfDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		if *request.WhTrfDate != "" {
			WhTrfDate, err := str.DateStrToRfc3339String(*request.WhTrfDate)
			if err != nil {
				return err
			}
			request.WhTrfDate = &WhTrfDate
		}

	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.WhTrf
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.WhTrfRepository.Update(txCtx, whTrfNo, request.CustID, Model)
		if err != nil {
			return err
		}
		DetailIds := []int{}

		for _, detail := range request.Details {
			if detail.WhTrfDetId != nil {
				DetailIds = append(DetailIds, *detail.WhTrfDetId)
			}
		}
		if len(DetailIds) > 0 {
			err := service.WhTrfRepository.DeleteDetailNotInIDs(txCtx, whTrfNo, request.CustID, DetailIds)
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
			detail.WhTrfNo = whTrfNo
			var whTrfDetModel model.WhTrfDet
			err = structs.Automapper(detail, &whTrfDetModel)
			if err != nil {
				return err
			}
			if detail.WhTrfDetId == nil || *detail.WhTrfDetId == 0 {
				whTrfDetModel.WhTrfDetId = nil
				err = service.WhTrfRepository.StoreDetail(txCtx, &whTrfDetModel)
				if err != nil {
					return err
				}
			} else {
				whTrfDetModel.CustID = ""

				err = service.WhTrfRepository.UpdateGrDetail(txCtx, &whTrfDetModel)
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

func (service *whTrfServiceImpl) ListWarehouse(dataFilter entity.StockTranferWarehouseQueryFilter, custId string) (data []entity.StockTranferWarehouse, total int64, lastPage int, err error) {
	warehouses, total, lastPage, err := service.WhTrfRepository.FindWarehouse(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range warehouses {
		var vResp entity.StockTranferWarehouse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}
