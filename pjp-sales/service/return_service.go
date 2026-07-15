package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sales/entity"
	"sales/model"
	"sales/pkg/conversion"
	"sales/pkg/str"
	"sales/pkg/structs"
	"sales/repository"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type ReturnService interface {
	List(dataFilter entity.ReturnQueryFilter) (data []entity.ReturnListResponse, total int64, lastPage int, err error)
	ShipmentList(dataFilter entity.ReturnQueryFilter) (data []entity.ReturnShipmentListResponse, total int64, lastPage int, err error)
	SetCreateReturnRequest(request entity.CreateReturnBody) (responses []entity.CreateReturnRequestBody, err error)
	// Store(request entity.CreateReturnBody) (err error)
	Store(request entity.CreateReturnRequestBody) (err error)
	Detail(returnNo string, custID string, parentCustID string) (response entity.ReturnResponse, err error)
	// Delete(custId string, returnNo string, userId int64) (err error)
	Update(returnNo string, request entity.UpdateReturnBody) (err error)
	UpdateQuantity(returnNo string, request entity.UpdateQuantityReturnBody) (err error)
	Approve(returnNo string, request entity.ApproveReturnBody) (err error)
	Cancel(returnNo string, request entity.CancelReturnBody) (err error)
	UpdateStatus(request entity.UpdateStatusReturnBody) (err error)
	UpdateAssign(request entity.UpdateAssignReturnBody) (err error)
	Print(custId string, returnNo string, userId int64) (err error)

	SalesmanFilterLookupList(entity.GeneralQueryFilter) (data []entity.SalesmansLookupResponse, total int64, lastPage int, err error)
	EmployeeFilterLookupList(entity.SalesmanQueryFilter) (data []entity.EmployeeLookupResponse, total int64, lastPage int, err error)
	RoleFilterLookupList(entity.GeneralQueryFilter) (data []entity.EmpGroupLookupResponse, total int64, lastPage int, err error)
	OutletFilterLookupList(entity.GeneralQueryFilter) (data []entity.OutletsLookupResponse, total int64, lastPage int, err error)
	ReturnStatusesLookupList(entity.GeneralQueryFilter) (data []entity.ReturnStatusesLookupResponse, total int64, lastPage int, err error)
	SalesmanFilterLookupCreate(entity.GeneralQueryFilter) (data []entity.SalesmansLookupResponse, total int64, lastPage int, err error)
	OutletFilterLookupCreate(entity.OutletCreateReturnQueryFilter) (data []entity.OutletsLookupResponse, total int64, lastPage int, err error)
	ProductListCreate(dataFilter entity.ProductListQueryFilter) (data []entity.ProductListResponse, total int64, lastPage int, err error)
	ReturnReasonLookupList(entity.GeneralQueryFilter) (data []entity.ReturnReasonsLookupResponse, total int64, lastPage int, err error)
	WarehouseLookupList(entity.WarehouseQueryFilter) (data []entity.WarehousesLookupResponse, total int64, lastPage int, err error)
	ProductLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.ProductsLookupCreateResponse, total int64, lastPage int, err error)
	// ProductConditionLookupList(entity.GeneralQueryFilter) (data []entity.ProductConditionsLookupResponse, total int64, lastPage int, err error)
}

func NewReturnService(returnRepository repository.ReturnRepository, orderRepository repository.OrderRepository, promotionRepository repository.PromotionRepository, discountRepository repository.DiscountRepository, transaction repository.Dbtransaction) *returnServiceImpl {
	return &returnServiceImpl{
		ReturnRepository:    returnRepository,
		OrderRepository:     orderRepository,
		PromotionRepository: promotionRepository,
		DiscountRepository:  discountRepository,
		Transaction:         transaction,
	}
}

type returnServiceImpl struct {
	ReturnRepository    repository.ReturnRepository
	OrderRepository     repository.OrderRepository
	PromotionRepository repository.PromotionRepository
	DiscountRepository  repository.DiscountRepository
	Transaction         repository.Dbtransaction
}

func (service *returnServiceImpl) List(dataFilter entity.ReturnQueryFilter) (data []entity.ReturnListResponse, total int64, lastPage int, err error) {
	rtns, total, lastPage, err := service.ReturnRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range rtns {
		var vResp entity.ReturnListResponse
		structs.Automapper(row, &vResp)

		if row.ReturnDate != nil {
			ReturnDate := row.ReturnDate.Format("2006-01-02")
			vResp.ReturnDate = ReturnDate
		}

		if row.InvoiceDate != nil {
			InvDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = &InvDate
		}

		returnStatusName := vResp.GenerateReturnStatusName()
		vResp.DataStatusName = &returnStatusName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnServiceImpl) ShipmentList(dataFilter entity.ReturnQueryFilter) (data []entity.ReturnShipmentListResponse, total int64, lastPage int, err error) {
	rtns, total, lastPage, err := service.ReturnRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range rtns {
		var vResp entity.ReturnShipmentListResponse
		structs.Automapper(row, &vResp)

		if row.ReturnDate != nil {
			ReturnDate := row.ReturnDate.Format("2006-01-02")
			vResp.ReturnDate = ReturnDate
		}

		if row.InvoiceDate != nil {
			InvDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = &InvDate
		}

		returnStatusName := vResp.GenerateReturnStatusName()
		vResp.DataStatusName = &returnStatusName

		Details, err := service.ReturnRepository.FindReturnDetail(row.ReturnNo, dataFilter.CustId, dataFilter.ParentCustId)
		if err != nil {
			return data, total, lastPage, err
		}

		var DetailsData []entity.ReturnDetailResponse
		for _, detail := range Details {
			var detailData entity.ReturnDetailResponse
			err = structs.Automapper(detail, &detailData)
			if err != nil {
				return data, total, lastPage, err
			}

			if row.InvoiceNo != nil {
				returnedProducts, err := service.ReturnRepository.CountReturnedProductQty(*row.InvoiceNo, detail.ProductID, dataFilter.CustId)
				if err != nil {
					return data, total, lastPage, err
				}

				detailData.RemainingQty1 = detail.InvoiceQty1 - returnedProducts.ReturnedQty1
				detailData.RemainingQty2 = detail.InvoiceQty2 - returnedProducts.ReturnedQty2
				detailData.RemainingQty3 = detail.InvoiceQty3 - returnedProducts.ReturnedQty3
			}

			itemConditionName := detailData.GenerateItemConditionName()
			detailData.ItemCndName = &itemConditionName

			// vResp.TotalVolume += detailData.Volume
			// vResp.TotalWeight += detailData.Weight

			vResp.TotalVolume += (detail.Volume1 * detail.Qty1) + (detail.Volume2 * detail.Qty2) + (detail.Volume3 * detail.Qty3)
			vResp.TotalWeight += (detail.Weight * detail.Qty1) + (detail.Weight2 * detail.Qty2) + (detail.Weight3 * detail.Qty3)

			DetailsData = append(DetailsData, detailData)
		}
		vResp.Details = DetailsData

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnServiceImpl) SetCreateReturnRequest(request entity.CreateReturnBody) (responses []entity.CreateReturnRequestBody, err error) {
	var response entity.CreateReturnRequestBody
	if err := structs.Automapper(request, &response); err != nil {
		return responses, err
	}
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

	var invoiceNo string
	var nextInvoiceNo string
	var detailResponses []entity.CreateReturnDetailRequestBody
	for index, Detail := range returnDetails {

		if Detail.InvoiceNo == nil {
			response.DataStatus = 1
		} else {
			response.DataStatus = 3

			if index < len(returnDetails)-1 {
				resInvoiceNo, _ := json.Marshal(Detail.InvoiceNo)
				invoiceNo = string(resInvoiceNo)
				resNextInvoiceNo, _ := json.Marshal(returnDetails[index+1].InvoiceNo)
				nextInvoiceNo = string(resNextInvoiceNo)
			}
		}

		if Detail.InvoiceDate != nil {
			invoiceDate, err := str.DateStrToRfc3339String(*Detail.InvoiceDate)
			if err != nil {
				return responses, err
			}
			Detail.InvoiceDate = &invoiceDate
		}

		// response.ReturnDate = time.Now().Format("2006-10-20")
		response.InvoiceNo = Detail.InvoiceNo
		response.InvoiceDate = Detail.InvoiceDate
		response.SalesmanID = *Detail.SalesmanID
		response.OutletID = *Detail.OutletID
		response.WhId = Detail.WhID

		var detailResponse entity.CreateReturnDetailRequestBody
		if err = structs.Automapper(Detail, &detailResponse); err != nil {
			return responses, err
		}
		detailResponse.CustID = response.CustID

		detailResponses = append(detailResponses, detailResponse)

		if Detail.InvoiceNo == nil || index == len(returnDetails)-1 || (index < len(returnDetails)-1 && invoiceNo != nextInvoiceNo) {

			response.Details = detailResponses
			responses = append(responses, response)

			response = entity.CreateReturnRequestBody{}
			if err = structs.Automapper(request, &response); err != nil {
				return responses, err
			}

			detailResponses = []entity.CreateReturnDetailRequestBody{}
		}
	}

	if err != nil {
		return responses, err
	}
	return responses, nil
}

func (service *returnServiceImpl) Store(request entity.CreateReturnRequestBody) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to create return: %v", r)
		}
	}()

	c := context.Background()
	log.Info("Masuk sebelum Initial Commit")
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		var returnModel model.Return
		if err = structs.Automapper(request, &returnModel); err != nil {
			return err
		}
		// returnModel.ReturnDate = time.Now()
		log.Info("Masuk setelah Automapper")

		if request.InvoiceDate != nil {
			parsedInvoiceDate, err := time.Parse(time.RFC3339, *request.InvoiceDate)
			if err != nil {
				return err
			}
			returnModel.InvoiceDate = &parsedInvoiceDate
		}

		returnModel.ReturnDate = time.Now()
		// parsedReturnDate, err := time.Parse(time.RFC3339, returnModel.ReturnDate)
		// if err != nil {
		// 	log.Info("Error di sini")
		// 	return err
		// }
		// returnModel.ReturnDate = parsedReturnDate

		var order model.OrderList
		var orderDetails []model.OrderDetailRead
		var orderRewards []model.FullPromoRewardRead
		// var qtyOrderDetail float64 = 0
		orderDetailNormalMaps := make(map[int64]*model.OrderDetailRead)
		orderDetailByIDMaps := make(map[int64]*model.OrderDetailRead)
		orderDetailPromoMaps := make(map[string]map[int64]*model.OrderDetailRead)
		orderDetailPromoProductIDMaps := make(map[string][]int64)
		qtyOrderDetail := 0.0
		if returnModel.InvoiceNo != nil {
			// AMBIL DATA ORDER
			order, err = service.OrderRepository.FindByInvoiceNo(*returnModel.InvoiceNo, returnModel.CustID)
			if err != nil {
				return err
			}
			// AMBIL DATA ORDER DETAIL
			orderDetails, err = service.OrderRepository.FindDetail(order.RoNo, order.CustID)
			if err != nil {
				return err
			}
			// MAPPING DATA ORDER DETAIL UNTUK NORMAL DAN PROMO
			for i := range orderDetails {
				detail := &orderDetails[i]
				if detail.OrderDetailID != nil {
					orderDetailByIDMaps[int64(*detail.OrderDetailID)] = detail
				}
				if detail.ItemType == 1 {
					// ORDER DETAIL DENGAN ITEM TYPE "NORMAL"
					orderDetailNormalMaps[int64(detail.ProId)] = detail
					qtyOrderDetail += *detail.Qty
				} else {
					// ORDER DETAIL DENGAN ITEM TYPE "PROMO"
					if _, isExist := orderDetailPromoMaps[*detail.PromoID]; !isExist {
						orderDetailPromoMaps[*detail.PromoID] = make(map[int64]*model.OrderDetailRead)
					}
					orderDetailPromoMaps[*detail.PromoID][int64(detail.ProId)] = detail
					orderDetailPromoProductIDMaps[*detail.PromoID] = append(orderDetailPromoProductIDMaps[*detail.PromoID], int64(detail.ProId))
				}
			}
		}
		// INPUT DATA PADA MODEL RETURN DAN RETURN DETAIL
		qtyReturnDetail := 0.0
		var returnDetailModelList []model.ReturnDetail
		for _, detailRequest := range request.Details {
			var returnDetailModel model.ReturnDetail
			if err = structs.Automapper(detailRequest, &returnDetailModel); err != nil {
				return err
			}
			promoValue := float64(0)
			promoBgValue := float64(0)
			discValue := float64(0)
			if returnDetailModel.ConvUnit2 == nil || returnDetailModel.ConvUnit3 == nil {
				return fmt.Errorf("conversion unit is required for product_id %d", detailRequest.ProductID)
			}
			totalQtyReturnDetail := returnDetailModel.Qty1 + (returnDetailModel.Qty2 * *returnDetailModel.ConvUnit2) + (returnDetailModel.Qty3 * *returnDetailModel.ConvUnit3 * *returnDetailModel.ConvUnit2)
			qtyReturnDetail += totalQtyReturnDetail
			if returnModel.InvoiceNo != nil {
				orderDetail := orderDetailNormalMaps[int64(detailRequest.ProductID)]
				if detailRequest.OrderDetailID != nil {
					if detailByID, ok := orderDetailByIDMaps[*detailRequest.OrderDetailID]; ok && detailByID.ItemType == 1 {
						orderDetail = detailByID
					}
				}

				if orderDetail == nil {
					return fmt.Errorf("order detail not found for product_id %d", detailRequest.ProductID)
				}

				qty1 := getValueOrDefault(orderDetail.Qty1, 0)
				qty2 := getValueOrDefault(orderDetail.Qty2, 0)
				qty3 := getValueOrDefault(orderDetail.Qty3, 0)
				conv2 := orderDetail.ConvUnit2
				if conv2 == nil {
					conv2 = orderDetail.MpConvUnit2
				}
				conv3 := orderDetail.ConvUnit3
				if conv3 == nil {
					conv3 = orderDetail.MpConvUnit3
				}
				if conv2 == nil || *conv2 == 0 {
					defaultConv2 := int(getValueOrDefault(returnDetailModel.ConvUnit2, 1))
					conv2 = &defaultConv2
				}
				if conv3 == nil || *conv3 == 0 {
					defaultConv3 := int(getValueOrDefault(returnDetailModel.ConvUnit3, 1))
					conv3 = &defaultConv3
				}

				totalQtyOrderDetail := qty1 + (qty2 * float64(*conv2)) + (qty3 * float64(*conv3) * float64(*conv2))
				if totalQtyOrderDetail <= 0 && orderDetail.Qty != nil && *orderDetail.Qty > 0 {
					// fallback untuk data lama yang tidak menyimpan qty1/qty2/qty3 dengan lengkap
					totalQtyOrderDetail = *orderDetail.Qty
				}
				if totalQtyOrderDetail <= 0 {
					// fallback terakhir agar request tidak gagal karena data historis tidak konsisten
					totalQtyOrderDetail = totalQtyReturnDetail
				}

				if orderDetail.ItemType == 1 {
					promoValue = math.Round(getValueOrDefault(orderDetail.PromoValue, 0) * (totalQtyReturnDetail / totalQtyOrderDetail))
					discValue = math.Round(getValueOrDefault(orderDetail.DiscValue, 0) * (totalQtyReturnDetail / totalQtyOrderDetail))
				}
				qtySisa := totalQtyOrderDetail - totalQtyReturnDetail
				orderDetail.Qty = &qtySisa

				qtySisaConversion := &conversion.Qty{
					Qty:       int(getValueOrDefault(&qtySisa, 0)),
					ConvUnit2: int(*conv2),
					ConvUnit3: int(*conv3),
				}

				qtyConversion := qtySisaConversion.ConvToQtyConversion()

				orderDetailQty := float64(qtySisaConversion.Qty)
				orderDetailQty1 := float64(qtyConversion.Qty1)
				orderDetailQty2 := float64(qtyConversion.Qty2)
				orderDetailQty3 := float64(qtyConversion.Qty3)
				orderDetail.Qty = &orderDetailQty
				orderDetail.Qty1 = &orderDetailQty1
				orderDetail.Qty2 = &orderDetailQty2
				orderDetail.Qty3 = &orderDetailQty3
			}
			returnDetailModel.PromoValue = &promoValue
			returnDetailModel.DiscValue = &discValue
			returnDetailModel.CustID = returnModel.CustID
			returnDetailModel.ItemType = int64(1)
			if orderDetailNormalMap, ok := orderDetailNormalMaps[int64(detailRequest.ProductID)]; ok {
				returnDetailModel.ItemType = int64(orderDetailNormalMap.ItemType)
			}
			returnDetailModel.Qty = returnDetailModel.Qty1 + (returnDetailModel.Qty2 * *returnDetailModel.ConvUnit2) + (returnDetailModel.Qty3 * *returnDetailModel.ConvUnit2 * *returnDetailModel.ConvUnit3)
			returnDetailModel.SubTotal = (returnDetailModel.Qty1 * returnDetailModel.SellPrice1) + (returnDetailModel.Qty2 * returnDetailModel.SellPrice2) + (returnDetailModel.Qty3 * returnDetailModel.SellPrice3)
			returnDetailModel.VatValue = math.Round((returnDetailModel.SubTotal - *returnDetailModel.PromoValue - *returnDetailModel.DiscValue) * (returnDetailModel.Vat / 100.0))
			returnDetailModel.Total = returnDetailModel.SubTotal - *returnDetailModel.PromoValue - *returnDetailModel.DiscValue + returnDetailModel.VatValue
			returnModel.SubTotal += returnDetailModel.SubTotal
			returnModel.VatValue += returnDetailModel.VatValue
			returnModel.DiscValue += discValue
			returnModel.PromoValue += promoValue
			returnModel.PromoBgValue += promoBgValue

			returnDetailModelList = append(returnDetailModelList, returnDetailModel)
		}

		// CEK APAKAH FULL RETURN.
		// JIKA IYA, MAKA FULL RETURN BARANG PROMO
		if qtyOrderDetail == qtyReturnDetail {
			for promoID := range orderDetailPromoMaps {
				for _, orderDetailPromoMap := range orderDetailPromoMaps[promoID] {
					var returnDetailModel model.ReturnDetail
					if err = structs.Automapper(orderDetailPromoMap, &returnDetailModel); err != nil {
						return err
					}

					returnDetailModel.CustID = returnModel.CustID
					returnDetailModel.ProductID = orderDetailPromoMap.ProId
					returnDetailModel.WhId = *request.WhId
					returnDetailModel.ItemCnd = 1
					returnDetailModel.ReturnReasonID = 0
					returnDetailModel.SubTotal = (returnDetailModel.Qty1 * returnDetailModel.SellPrice1) + (returnDetailModel.Qty2 * returnDetailModel.SellPrice2) + (returnDetailModel.Qty3 * returnDetailModel.SellPrice3)
					// returnDetailModel.VatValue = math.Round((returnDetailModel.SubTotal - *returnDetailModel.PromoValue - *returnDetailModel.DiscValue) * (returnDetailModel.Vat / 100.0))
					returnDetailModel.Total = *orderDetailPromoMap.Amount
					returnModel.SubTotal += returnDetailModel.SubTotal
					// returnModel.VatValue += returnDetailModel.VatValue
					// returnModel.DiscValue += discValue
					// returnModel.PromoValue += promoValue
					returnModel.PromoBgValue += *orderDetailPromoMap.PromoValue

					returnDetailModelList = append(returnDetailModelList, returnDetailModel)
				}
			}
		} else {
			// AMBIL DATA PROMO BARANG YANG TELAH DIDAPATKAN
			if returnModel.InvoiceNo != nil {
				orderRewards, err = service.OrderRepository.FindFullPromoRewards(*returnModel.InvoiceNo, returnModel.CustID)
				if err != nil {
					return err
				}
			}
			for _, orderReward := range orderRewards {
				// AMBIL DATA PROMO ADDITIONAL CRITERIA DENGAN ATTRIBUTE PRO
				promoAdditionalCriterias, err := service.PromotionRepository.FindPromoAdditionalCriteriasWithProductAttributeByPromoID(orderReward.PromoID, returnModel.CustID)
				if err != nil {
					return err
				}

				slabRules := 0.0
				var orderDetailWithPAC []*model.OrderDetailRead
				for _, promoAdditionalCriteria := range promoAdditionalCriterias {
					slabRule := 0.0

					// CEK APAKAH REFERENCE_ID (PRODUCT_ID) PAC ADA DI DETAIL ORDER NORMAL
					if _, isExist := orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID]; isExist {
						orderDetailWithPAC = append(orderDetailWithPAC, orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID])

						switch orderReward.SlabRuleUom {
						case 2:
							slabRule = (*orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].Qty3 * float64(*orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].ConvUnit3)) + *orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].Qty2
						case 1:
							slabRule = (*orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].Qty3 * float64(*orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].ConvUnit3) * float64(*orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].ConvUnit2)) + (*orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].Qty2 * float64(*orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].ConvUnit2)) + *orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].Qty1
						default:
							slabRule = *orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].Qty3
						}

						// UNTUK PAC YANG MANDATORY, CEK APAKAH MIN BUY VALUE NYA SUDAH TERPENUHI
						if promoAdditionalCriteria.IsMandatory {
							// PENDEFINISIAN BUY VALUE
							buyValue := slabRule
							if promoAdditionalCriteria.MinBuyType == 2 {
								buyValue = (*orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].Qty3 * float64(*orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].SellPrice3)) + (*orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].Qty2 * float64(*orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].SellPrice2)) + (*orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].Qty1 * float64(*orderDetailNormalMaps[promoAdditionalCriteria.ReferenceID].SellPrice1))
							}

							// CEK BUY VALUE LEBIH KECIL DARI MIN BUY VALUE. JIKA IYA, FULL RETURN DENGAN PROMO YBS
							if promoAdditionalCriteria.MinBuyValue > buyValue {
								for _, orderDetailPromoMap := range orderDetailPromoMaps[promoAdditionalCriteria.PromoID] {
									var returnDetailModel model.ReturnDetail
									if err = structs.Automapper(orderDetailPromoMap, &returnDetailModel); err != nil {
										return err
									}

									// returnDetailModel.PromoValue = &promoValue
									// returnDetailModel.DiscValue = &discValue

									returnDetailModel.CustID = returnModel.CustID
									returnDetailModel.ProductID = orderDetailPromoMap.ProId
									returnDetailModel.WhId = *request.WhId
									returnDetailModel.ItemCnd = 1
									returnDetailModel.ReturnReasonID = 0
									returnDetailModel.SubTotal = (returnDetailModel.Qty1 * returnDetailModel.SellPrice1) + (returnDetailModel.Qty2 * returnDetailModel.SellPrice2) + (returnDetailModel.Qty3 * returnDetailModel.SellPrice3)
									// returnDetailModel.VatValue = math.Round((returnDetailModel.SubTotal - *returnDetailModel.PromoValue - *returnDetailModel.DiscValue) * (returnDetailModel.Vat / 100.0))
									returnDetailModel.Total = *orderDetailPromoMap.Amount
									returnModel.SubTotal += returnDetailModel.SubTotal
									// returnModel.VatValue += returnDetailModel.VatValue
									// returnModel.DiscValue += discValue
									// returnModel.PromoValue += promoValue
									returnModel.PromoBgValue += *orderDetailPromoMap.PromoValue

									returnDetailModelList = append(returnDetailModelList, returnDetailModel)
								}

								break
							}
						}
					}

					slabRules += slabRule
				}

				if slabRules == 0.0 {
					continue
				}

				updatedPromoCriteria, err := service.PromotionRepository.FindPromoCriteriaBySlabRule(orderReward, returnModel.CustID, slabRules)
				// JIKA DATA PROMO CRITERIA DENGAN SLAB RULE TERBARU TIDAK DITEMUKAN, MAKA FULL RETURN BARANG PROMO
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						for _, orderDetailPromoMap := range orderDetailPromoMaps[orderReward.PromoID] {
							var returnDetailModel model.ReturnDetail
							if err = structs.Automapper(orderDetailPromoMap, &returnDetailModel); err != nil {
								return err
							}

							// returnDetailModel.PromoValue = &promoValue
							// returnDetailModel.DiscValue = &discValue

							returnDetailModel.CustID = returnModel.CustID
							returnDetailModel.ProductID = orderDetailPromoMap.ProId
							returnDetailModel.WhId = *request.WhId
							returnDetailModel.ItemCnd = 1
							returnDetailModel.ReturnReasonID = 0
							returnDetailModel.SubTotal = (returnDetailModel.Qty1 * returnDetailModel.SellPrice1) + (returnDetailModel.Qty2 * returnDetailModel.SellPrice2) + (returnDetailModel.Qty3 * returnDetailModel.SellPrice3)
							// returnDetailModel.VatValue = math.Round((returnDetailModel.SubTotal - *returnDetailModel.PromoValue - *returnDetailModel.DiscValue) * (returnDetailModel.Vat / 100.0))
							returnDetailModel.Total = *orderDetailPromoMap.Amount
							returnModel.SubTotal += returnDetailModel.SubTotal
							// returnModel.VatValue += returnDetailModel.VatValue
							// returnModel.DiscValue += discValue
							// returnModel.PromoValue += promoValue
							returnModel.PromoBgValue += *orderDetailPromoMap.PromoValue

							returnDetailModelList = append(returnDetailModelList, returnDetailModel)
						}
					} else {
						return err
					}
				}

				// if *newPromoCriteria.SlabID == int64(*orderReward.SlabId) {
				// 	continue
				// }
				// PENGHITUNGAN REWARD KELIPATAN
				multipliedValue := int64(1)
				if orderReward.IsMultiplied {
					// log.Info("validatedPromoCriterias["+promoCriteria.PromoID+"].SlabRule : ", validatedPromoCriterias[promoCriteria.PromoID].SlabRule)
					multipliedValue = int64(slabRules / updatedPromoCriteria.SlabRuleTo)
					// log.Info("MULTIPLIED_VALUE : ", multipliedValue)
				}
				updatedQtyReward := updatedPromoCriteria.SlabReward * float64(multipliedValue)

				// prevQtyReward := qtyOrderDetailPromoMaps[updatedPromoCriteria.PromoID]

				// MENGHITUNG QTY PADA ORDER DETAIL PROMO DISESUAIKAN ORDER REWARD TYPE
				qtyOrderDetailPromo := 0.0
				for _, orderDetailPromoMap := range orderDetailPromoMaps[orderReward.PromoID] {
					qtyPromo := 0.0

					switch updatedPromoCriteria.SlabRewardType {
					case entity.PromoRewardTypeFixedValue:
						qtyPromo = (*orderDetailPromoMap.Qty3 * float64(*orderDetailPromoMap.ConvUnit3)) + *orderDetailPromoMap.Qty2
					case entity.PromoRewardTypeQuantity:
						qtyPromo = (*orderDetailPromoMap.Qty3 * float64(*orderDetailPromoMap.ConvUnit3) * float64(*orderDetailPromoMap.ConvUnit2)) + (*orderDetailPromoMap.Qty2 * float64(*orderDetailPromoMap.ConvUnit2)) + *orderDetailPromoMap.Qty1
					default:
						qtyPromo = *orderDetailPromoMap.Qty3
					}

					orderDetailPromoMap.Qty = &qtyPromo
					qtyOrderDetailPromo += qtyPromo
				}

				if updatedQtyReward >= qtyOrderDetailPromo {
					continue
				}

				sisaQtyReward := qtyOrderDetailPromo - updatedQtyReward
				for _, orderDetailPromoMap := range orderDetailPromoMaps[orderReward.PromoID] {
					var returnDetailModel model.ReturnDetail
					if err = structs.Automapper(orderDetailPromoMap, &returnDetailModel); err != nil {
						return err
					}

					sisaQtyReward -= *orderDetailPromoMap.Qty
					returnDetailModel.Qty = (*orderDetailPromoMap.Qty3 * float64(*orderDetailPromoMap.ConvUnit3) * float64(*orderDetailPromoMap.ConvUnit2)) + (*orderDetailPromoMap.Qty2 * float64(*orderDetailPromoMap.ConvUnit2)) + *orderDetailPromoMap.Qty1
					if sisaQtyReward < 0 {
						qtyReturn := *orderDetailPromoMap.Qty + sisaQtyReward
						switch updatedPromoCriteria.SlabRewardType {
						case entity.PromoRewardTypeFixedValue:
							qtyReturn = qtyReturn * float64(*orderDetailPromoMap.ConvUnit2)
						case entity.PromoRewardTypePercent:
							qtyReturn = qtyReturn * float64(*orderDetailPromoMap.ConvUnit3) * float64(*orderDetailPromoMap.ConvUnit2)
						}

						qtyReturnConversion := &conversion.Qty{
							Qty:       int(getValueOrDefault(&qtyReturn, 0)),
							ConvUnit2: int(*orderDetailPromoMap.ConvUnit2),
							ConvUnit3: int(*orderDetailPromoMap.ConvUnit3),
						}

						qtyConversion := qtyReturnConversion.ConvToQtyConversion()

						returnDetailModel.Qty = float64(qtyReturnConversion.Qty)
						returnDetailModel.Qty1 = float64(qtyConversion.Qty1)
						returnDetailModel.Qty2 = float64(qtyConversion.Qty2)
						returnDetailModel.Qty3 = float64(qtyConversion.Qty3)
					}
					promoValue := math.Round(float64(*orderDetailPromoMap.PromoValue) * (returnDetailModel.Qty / ((*orderDetailPromoMap.Qty3 * float64(*orderDetailPromoMap.ConvUnit3) * float64(*orderDetailPromoMap.ConvUnit2)) + (*orderDetailPromoMap.Qty2 * float64(*orderDetailPromoMap.ConvUnit2)) + *orderDetailPromoMap.Qty1)))
					returnDetailModel.PromoValue = &promoValue
					// returnDetailModel.DiscValue = &discValue

					returnDetailModel.CustID = returnModel.CustID
					returnDetailModel.ProductID = orderDetailPromoMap.ProId
					returnDetailModel.WhId = *request.WhId
					returnDetailModel.ItemCnd = 1
					returnDetailModel.ReturnReasonID = 0
					returnDetailModel.SubTotal = (returnDetailModel.Qty1 * returnDetailModel.SellPrice1) + (returnDetailModel.Qty2 * returnDetailModel.SellPrice2) + (returnDetailModel.Qty3 * returnDetailModel.SellPrice3)
					// returnDetailModel.VatValue = math.Round((returnDetailModel.SubTotal - *returnDetailModel.PromoValue - *returnDetailModel.DiscValue) * (returnDetailModel.Vat / 100.0))
					// returnDetailModel.Total = *orderDetailPromoMap.Amount
					returnDetailModel.Total = returnDetailModel.SubTotal - *returnDetailModel.PromoValue
					returnModel.SubTotal += returnDetailModel.SubTotal
					// returnModel.VatValue += returnDetailModel.VatValue
					// returnModel.DiscValue += discValue
					// returnModel.PromoValue += promoValue
					returnModel.PromoBgValue += promoValue

					returnDetailModelList = append(returnDetailModelList, returnDetailModel)

					if sisaQtyReward <= 0 {
						break
					}
				}
			}
		}

		returnModel.Total = returnModel.SubTotal - returnModel.PromoValue - returnModel.PromoBgValue - returnModel.DiscValue + returnModel.VatValue
		if err := service.ReturnRepository.Store(txCtx, &returnModel); err != nil {
			return err
		}

		for _, returnDetail := range returnDetailModelList {
			returnDetail.ReturnNo = returnModel.ReturnNo
			if err = service.ReturnRepository.StoreDetail(txCtx, &returnDetail); err != nil {
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

func (service *returnServiceImpl) StoreOld(request entity.CreateReturnBody) (err error) {
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
			log.Info(i, " : ", string(jsonF))
		}
	*/
	// var createdReturnList []entity.CreatedReturnBody
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
		/*
			orderDetailList := make(map[int64]*model.OrderDetailRead)
			var orderDetailIdList []int64
			var consultPromotionBody entity.ConsultPromotionBody
		*/
		for index, Detail := range returnDetails {

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

			ReturnModel.ReturnDate = time.Now()
			ReturnModel.InvoiceNo = Detail.InvoiceNo
			ReturnModel.SalesmanID = *Detail.SalesmanID
			ReturnModel.OutletID = *Detail.OutletID

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

			promoValue := float64(0)
			// promoBgValue := float64(0)
			discValue := float64(0)
			if Detail.InvoiceNo != nil {

				orderDetail, err := service.OrderRepository.FindOrderDetailByDetailID(*Detail.OrderDetailID, ReturnModel.CustID)
				if err != nil {
					return err
				}

				totalQtyOrderDetail := *orderDetail.Qty1 + (*orderDetail.Qty2 * float64(*orderDetail.ConvUnit2)) + (*orderDetail.Qty3 * float64(*orderDetail.ConvUnit3) * float64(*orderDetail.ConvUnit2))
				totalQtyReturnDetail := ReturnDetailModel.Qty1 + (ReturnDetailModel.Qty2 * *ReturnDetailModel.ConvUnit2) + (ReturnDetailModel.Qty3 * *ReturnDetailModel.ConvUnit3 * *ReturnDetailModel.ConvUnit2)

				if orderDetail.ItemType == 1 {
					promoValue = math.Round(float64(*orderDetail.PromoValue) * (totalQtyReturnDetail / totalQtyOrderDetail))
				}
				// else {
				// 	promoBgValue = math.Round(float64(*orderDetail.PromoValue) * (totalQtyReturnDetail / totalQtyOrderDetail))
				// }
				discValue = math.Round(float64(*orderDetail.DiscValue) * (totalQtyReturnDetail / totalQtyOrderDetail))
				/*
					orderDetailList[int64(orderDetail.ProId)] = &orderDetail
					orderDetailIdList = append(orderDetailIdList, int64(*orderDetail.OrderDetailID))

					consultPromotionBody.Details = append(consultPromotionBody.Details, entity.ConsultPromotionSubBody{
						ProID      :
						Qty1       :
						Qty2       :
						Qty3       :
						ConvUnit2  :
						ConvUnit3  :
						SubTotal   :
						SellPrice1 :
						SellPrice2 :
						SellPrice3 :
						SellPrice4 :

						SellPrice5 :
					})
				*/
			}
			ReturnDetailModel.PromoValue = &promoValue
			// ReturnDetailModel.PromoBgValue = &promoBgValue
			ReturnDetailModel.DiscValue = &discValue

			ReturnDetailModel.CustID = ReturnModel.CustID
			ReturnDetailModel.SubTotal = (ReturnDetailModel.Qty1 * ReturnDetailModel.SellPrice1) + (ReturnDetailModel.Qty2 * ReturnDetailModel.SellPrice2) + (ReturnDetailModel.Qty3 * ReturnDetailModel.SellPrice3)
			ReturnDetailModel.VatValue = math.Round((ReturnDetailModel.SubTotal - *ReturnDetailModel.PromoValue - *ReturnDetailModel.DiscValue) * (ReturnDetailModel.Vat / 100.0))
			ReturnDetailModel.Total = ReturnDetailModel.SubTotal - *ReturnDetailModel.PromoValue - *ReturnDetailModel.DiscValue + ReturnDetailModel.VatValue
			ReturnModel.SubTotal += ReturnDetailModel.SubTotal
			ReturnModel.VatValue += ReturnDetailModel.VatValue
			ReturnModel.DiscValue += *ReturnDetailModel.DiscValue
			ReturnModel.PromoValue += *ReturnDetailModel.PromoValue

			ReturnDetailModelList = append(ReturnDetailModelList, ReturnDetailModel)

			if Detail.InvoiceNo == nil || index == len(returnDetails)-1 || (index < len(returnDetails)-1 && invoiceNo != nextInvoiceNo) {

				/*
					var ro model.OrderList
					var orderRewards []model.FullPromoRewardRead
				*/
				// for _, returnDetail := range ReturnDetailModelList {
				// 	returnDetail.ReturnNo = ReturnModel.ReturnNo
				// 	if err = service.ReturnRepository.StoreDetail(txCtx, &returnDetail); err != nil {
				// 		return err
				// 	}
				// }

				/*
					if ReturnModel.InvoiceNo != nil {
						ro, err = service.OrderRepository.FindByInvoiceNo(*ReturnModel.InvoiceNo, ReturnModel.CustID)
						if err != nil {
							return err
						}



						if orderRewards, err = service.OrderRepository.FindFullPromoRewards(ro.RoNo, ro.CustID); err == nil {
							for index, reward := range orderRewards {
								if reward.SlabRewardType == 2 || reward.SlabRewardType == 3 {
									continue
								}


							}
						}
					}
				*/

				ReturnModel.Total = ReturnModel.SubTotal - ReturnModel.PromoValue - ReturnModel.DiscValue + ReturnModel.VatValue
				if err := service.ReturnRepository.Store(txCtx, &ReturnModel); err != nil {
					return err
				}
				/*
					var CreatedReturn entity.CreatedReturnBody
					if ReturnModel.InvoiceNo != nil {
						if err = structs.Automapper(ReturnModel, &CreatedReturn); err != nil {
							return err
						}

						CreatedReturn.ReturnDate = ReturnModel.ReturnDate.Format("2006-01-02")

						if ReturnModel.InvoiceDate != nil {
							invoiceDate := ReturnModel.InvoiceDate.Format("2006-01-02")
							CreatedReturn.InvoiceDate = &invoiceDate
						}
					}

					ReturnModelList = append(ReturnModelList, ReturnModel)
				*/
				for _, returnDetail := range ReturnDetailModelList {
					returnDetail.ReturnNo = ReturnModel.ReturnNo
					if err = service.ReturnRepository.StoreDetail(txCtx, &returnDetail); err != nil {
						return err
					}
					/*
						var CreatedReturnDetail entity.CreatedReturnDetailBody
						if ReturnModel.InvoiceNo != nil {
							if err = structs.Automapper(returnDetail, &CreatedReturnDetail); err != nil {
								return err
							}

							CreatedReturn.Details = append(CreatedReturn.Details, CreatedReturnDetail)
						}
					*/
				}
				/*
					if ReturnModel.InvoiceNo != nil {
						createdReturnList = append(createdReturnList, CreatedReturn)
					}

					var ReturnModel model.Return
				*/
				ReturnModel = model.Return{}
				if err = structs.Automapper(request, &ReturnModel); err != nil {
					return err
				}

				ReturnDetailModelList = []model.ReturnDetail{}
				/*
					orderDetailList = make(map[int64]*model.OrderDetailRead)
					orderDetailIdList = []int64{}
				*/
			}
		}
		/*
			log.Info("===================================RETURN=======================================")
			for i, ReturnResult := range ReturnModelList {
				jsonF, _ := json.Marshal(ReturnResult)
				log.Info(i, " : ", string(jsonF))
			}
			log.Info("===================================RETURN DETAIL=======================================")
			for i, ReturnDetailResult := range ReturnDetailModelList {
				jsonF, _ := json.Marshal(ReturnDetailResult)
				log.Info(i, " : ", string(jsonF))
			}
		*/
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

/*
	func ConsultPromotionReturn(service *returnServiceImpl, request entity.ConsultPromotionBody) (responses []entity.ConsultPromotionResponse, err error) {
		for index, detail := range request.Details {
			var detailConversion entity.CreateConversionBody
			detailConversion.CustId = request.CustID
			detailConversion.ProductId = detail.ProID
			detailConversion.Qty1 = int64(detail.Qty1)
			detailConversion.Qty2 = int64(detail.Qty2)
			detailConversion.Qty3 = int64(detail.Qty3)
			consultPromotionBodyDetail, _ := service.Conversion(detailConversion, request.CustID, request.ParentCustID)

			request.Details[index].Qty1 = float64(consultPromotionBodyDetail.Qty1)
			request.Details[index].Qty2 = float64(consultPromotionBodyDetail.Qty2)
			request.Details[index].Qty3 = float64(consultPromotionBodyDetail.Qty3)
		}

		outlet, err := service.PromotionRepository.FindOutletByID(int64(request.OutletId), request.CustID, request.ParentCustID)
		if err != nil {
			return responses, err
		}

		salesman, err := service.PromotionRepository.FindSalesmanByID(int64(request.SalesmanId), request.CustID, request.ParentCustID)
		if err != nil {
			return responses, err
		}

		if request.WhId == 0 {
			request.WhId = salesman.WhId
		}
		log.Info("request : ", request)
		attributePromoValidationCriteriaList := make(map[string]int64)
		attributePromoValidationCriteriaList["OCL"] = int64(outlet.OtClassId)
		attributePromoValidationCriteriaList["OTG"] = int64(outlet.OtGrpId)
		attributePromoValidationCriteriaList["OTY"] = int64(outlet.OtTypeId)
		attributePromoValidationCriteriaList["STE"] = int64(salesman.SalesTeamId)

		promoAdditionalCriteriaGroups := make(map[string]map[string][]model.PromoAdditionalCriteriaByActivePromo)
		if promoAdditionalCriterias, err := service.PromotionRepository.FindAllPromoAdditionalCriteriasByActivePromo(request); err == nil {
			attributeListWithPro := [5]string{"OCL", "OTG", "OTY", "STE", "PRO"}
			attributePromoValidationList := make(map[string]map[string]bool)
			finalAttributePromoValidationList := make(map[string]map[string]bool)
			log.Info("PENGECEKAN PROMO ADDITIONAL CRITERIA")
			for index := range promoAdditionalCriterias {
				if index == 0 || promoAdditionalCriterias[index-1].PromoID != promoAdditionalCriterias[index].PromoID {
					log.Info("PENGECEKAN PROMO ", promoAdditionalCriterias[index].PromoID)
					log.Info("INISIALISASI FLAGGING PENGECEKAN PROMO " + promoAdditionalCriterias[index].PromoID + " DENGAN KRITERIA AWAL (OUTLET CLASS, OUTLET GROUP, OUTLET TYPE, SALES TEAM)")
					for _, attribute := range attributeListWithPro {
						promoAdditionalCriteriaGroups[promoAdditionalCriterias[index].PromoID] = make(map[string][]model.PromoAdditionalCriteriaByActivePromo)

						attributePromoValidationList[promoAdditionalCriterias[index].PromoID] = make(map[string]bool)
						finalAttributePromoValidationList[promoAdditionalCriterias[index].PromoID] = make(map[string]bool)
						attributePromoValidationList[promoAdditionalCriterias[index].PromoID][attribute] = false
						finalAttributePromoValidationList[promoAdditionalCriterias[index].PromoID][attribute] = false
						// log.Info("INISIALISASI attributePromoValidationList["+promoAdditionalCriterias[index].PromoID+"]["+attribute+"] : ", attributePromoValidationList[promoAdditionalCriterias[index].PromoID][attribute])
					}
				}

				// log.Info("GROUPING " + promoAdditionalCriterias[index].Attribute + " ATTRIBUTE PADA PROMO " + promoAdditionalCriterias[index].PromoID)
				promoAdditionalCriteriaGroups[promoAdditionalCriterias[index].PromoID][promoAdditionalCriterias[index].Attribute] = append(promoAdditionalCriteriaGroups[promoAdditionalCriterias[index].PromoID][promoAdditionalCriterias[index].Attribute], promoAdditionalCriterias[index])

				if promoAdditionalCriterias[index].Attribute == "PRO" {
					continue
				}

				// log.Info("PENGECEKAN PROMO ADDITIONAL CRITERIA DENGAN KRITERIA AWAL (OUTLET CLASS, OUTLET GROUP, OUTLET TYPE, SALES TEAM)")
				if !finalAttributePromoValidationList[promoAdditionalCriterias[index].PromoID][promoAdditionalCriterias[index].Attribute] {
					if promoAdditionalCriterias[index].ReferenceID == attributePromoValidationCriteriaList[promoAdditionalCriterias[index].Attribute] {
						if promoAdditionalCriterias[index].Condition == "I" {
							attributePromoValidationList[promoAdditionalCriterias[index].PromoID][promoAdditionalCriterias[index].Attribute] = true
						}
						finalAttributePromoValidationList[promoAdditionalCriterias[index].PromoID][promoAdditionalCriterias[index].Attribute] = true
					}
				}

				if index == len(promoAdditionalCriterias)-1 || promoAdditionalCriterias[index].Attribute != promoAdditionalCriterias[index+1].Attribute {
					if !finalAttributePromoValidationList[promoAdditionalCriterias[index].PromoID][promoAdditionalCriterias[index].Attribute] {
						attributePromoValidationList[promoAdditionalCriterias[index].PromoID][promoAdditionalCriterias[index].Attribute] = false
						if promoAdditionalCriterias[index].Condition == "E" {
							attributePromoValidationList[promoAdditionalCriterias[index].PromoID][promoAdditionalCriterias[index].Attribute] = true
						}
						finalAttributePromoValidationList[promoAdditionalCriterias[index].PromoID][promoAdditionalCriterias[index].Attribute] = true
					}
				}
			}

			// log.Info("SETELAH MAPPING attributePromoValidationList :", len(attributePromoValidationList))
			// for promoID := range attributePromoValidationList {
			// 	for attribute := range attributePromoValidationList[promoID] {
			// 		log.Info("attributePromoValidationList["+promoID+"]["+attribute+"] : ", attributePromoValidationList[promoID][attribute])
			// 	}
			// }

			// ELIMINASI PROMO YANG TIDAK SESUAI KETENTUAN AWAL
			log.Info("ELIMINASI PROMO YANG TIDAK SESUAI KETENTUAN AWAL")
			validatedPromoAdditionalCriteriaGroups := make(map[string]map[string][]model.PromoAdditionalCriteriaByActivePromo)
			for promoID := range attributePromoValidationList {
				isPromoValid := true
				for attribute := range attributePromoValidationList[promoID] {
					if !attributePromoValidationList[promoID][attribute] && attribute != "PRO" {
						isPromoValid = false
						break
					}
				}

				if isPromoValid {
					validatedPromoAdditionalCriteriaGroups[promoID] = make(map[string][]model.PromoAdditionalCriteriaByActivePromo)
					validatedPromoAdditionalCriteriaGroups[promoID] = promoAdditionalCriteriaGroups[promoID]
				}
			}

			log.Info("SETELAH ELIMINASI validatedPromoAdditionalCriteriaGroups :", len(validatedPromoAdditionalCriteriaGroups))
			// for promoID := range validatedPromoAdditionalCriteriaGroups {
			// 	log.Info("validatedPromoAdditionalCriteriaGroups["+promoID+"] :", len(validatedPromoAdditionalCriteriaGroups[promoID]))
			// 	for attribute := range validatedPromoAdditionalCriteriaGroups[promoID] {
			// 		log.Info("validatedPromoAdditionalCriteriaGroups["+promoID+"]["+attribute+"] : ", validatedPromoAdditionalCriteriaGroups[promoID][attribute])
			// 	}
			// }

			log.Info("PENGECEKAN PROMO ADDITIONAL CRITERIA DENGAN PRODUK YANG DIBELI")
			validatedPromoAdditionalCriteriaByProductGroups := make(map[string]map[int64]*entity.ConsultPromotionSubBody)
			subTotalValidatedPromoAdditionalCriteriaByProductGroups := make(map[string]int64)
			var validatedPromoList []string
			for promoID := range validatedPromoAdditionalCriteriaGroups {
				validatedPromoAdditionalCriteriaByProductGroups[promoID] = make(map[int64]*entity.ConsultPromotionSubBody)
				isPromoAdditionalCriteriaValid := false
				subTotalValidatedPromoAdditionalCriteriaByProductGroups[promoID] = 0
				// log.Info("LEN PAC "+promoID+" : ", len(validatedPromoAdditionalCriteriaGroups[promoID]["PRO"]))
				for _, validatedPromoAdditionalCriteria := range validatedPromoAdditionalCriteriaGroups[promoID]["PRO"] {
					// log.Info("PENGECEKAN PROMO ADDITIONAL CRITERIA DENGAN PRODUK YANG DIBELI")
					isMandatoryPromoAdditionalCriteriaNotValid := false
					for index, req := range request.Details {
						if validatedPromoAdditionalCriteria.ReferenceID != int64(req.ProID) {
							if !isPromoAdditionalCriteriaValid && index == len(request.Details)-1 && validatedPromoAdditionalCriteria.IsMandatory {
								isPromoAdditionalCriteriaValid = false
								isMandatoryPromoAdditionalCriteriaNotValid = true
							}
							continue
						}

						buyValue := float64(req.SubTotal)
						if validatedPromoAdditionalCriteria.MinBuyType == 1 {
							switch validatedPromoAdditionalCriteria.MinBuyUom {
							case 2:
								buyValue = (req.Qty3 * req.ConvUnit3) + req.Qty2
							case 1:
								buyValue = (req.Qty3 * req.ConvUnit3 * req.ConvUnit2) + (req.Qty2 * req.ConvUnit2) + req.Qty1
							default:
								buyValue = req.Qty3
							}
						}
						// log.Info("BUY VALUE : ", buyValue)

						if validatedPromoAdditionalCriteria.MinBuyValue > buyValue {
							if validatedPromoAdditionalCriteria.IsMandatory {
								isPromoAdditionalCriteriaValid = false
								isMandatoryPromoAdditionalCriteriaNotValid = true
								break
							}
							continue
						}

						isPromoAdditionalCriteriaValid = true
						validatedPromoAdditionalCriteriaByProductGroups[promoID][req.ProID] = &request.Details[index]
						log.Info("validatedPromoAdditionalCriteriaByProductGroups["+promoID+"]["+strconv.FormatInt(req.ProID, 10)+"] :", int64(req.SubTotal))
						subTotalValidatedPromoAdditionalCriteriaByProductGroups[promoID] += int64(req.SubTotal)
					}
					// log.Info("subTotalValidatedPromoAdditionalCriteriaByProductGroups["+promoID+"] :", subTotalValidatedPromoAdditionalCriteriaByProductGroups[promoID])

					if isMandatoryPromoAdditionalCriteriaNotValid {
						break
					}
				}

				if !isPromoAdditionalCriteriaValid {
					delete(validatedPromoAdditionalCriteriaByProductGroups, promoID)
					delete(subTotalValidatedPromoAdditionalCriteriaByProductGroups, promoID)
					continue
				}

				validatedPromoList = append(validatedPromoList, promoID)
			}

			log.Info("SETELAH PENGECEKAN PROMO ADDITIONAL CRITERIA validatedPromoAdditionalCriteriaByProductGroups : ", len(validatedPromoAdditionalCriteriaByProductGroups))
			// for promoID := range validatedPromoAdditionalCriteriaByProductGroups {
			// 	log.Info("validatedPromoAdditionalCriteriaByProductGroups["+promoID+"] : ", len(validatedPromoAdditionalCriteriaByProductGroups[promoID]))
			// 	for proID := range validatedPromoAdditionalCriteriaByProductGroups[promoID] {
			// 		log.Info("validatedPromoAdditionalCriteriaByProductGroups["+promoID+"]["+string(proID)+"] : ", validatedPromoAdditionalCriteriaByProductGroups[promoID][proID])
			// 	}
			// }

			log.Info("PENGECEKAN PROMO CRITERIA SESUAI DENGAN PROMO ADDITIONAL CRITERIA")
			validatedPromoCriterias := make(map[string]model.ConsultPromoCriteria)
			var validatedSlabIDs []string
			if promoCriterias, err := service.PromotionRepository.FindAllPromoCriteriasByPromoIDs(validatedPromoList); err == nil {
				isPromoCriteriaValid := false
				for index := range promoCriterias {
					if index == 0 || promoCriterias[index-1].PromoID != promoCriterias[index].PromoID {
						isPromoCriteriaValid = false
					}

					if isPromoCriteriaValid {
						continue
					}

					promoCriterias[index].SlabRule = 0
					slabRuleValue := float64(subTotalValidatedPromoAdditionalCriteriaByProductGroups[promoCriterias[index].PromoID])
					if promoCriterias[index].SlabRuleType == 1 {
						slabRuleValue = 0
						for _, req := range validatedPromoAdditionalCriteriaByProductGroups[promoCriterias[index].PromoID] {
							buyValue := 0.0

							// log.Info("SlabRuleUom : ", promoCriterias[index].SlabRuleUom)
							// log.Info("Qty1 : ", req.Qty1)
							// log.Info("Qty2 : ", req.Qty2)
							// log.Info("Qty3 : ", req.Qty3)
							switch promoCriterias[index].SlabRuleUom {
							case 2:
								buyValue = (req.Qty3 * req.ConvUnit3) + req.Qty2
							case 1:
								buyValue = (req.Qty3 * req.ConvUnit3 * req.ConvUnit2) + (req.Qty2 * req.ConvUnit2) + req.Qty1
							default:
								buyValue = req.Qty3
							}

							slabRuleValue += buyValue
						}

						promoCriterias[index].SlabRule = int64(slabRuleValue)
						// log.Info("SlabRule : ", promoCriterias[index].SlabRule)
					}
					// log.Info("slabRuleValue SLAB ("+strconv.FormatInt(*promoCriterias[index].SlabID, 10)+") :", slabRuleValue)

					// log.Info("promoCriterias["+strconv.Itoa(index)+"].IsMultiplied : ", promoCriterias[index].IsMultiplied)
					// log.Info("promoCriterias["+strconv.Itoa(index)+"].SlabRuleFrom : ", promoCriterias[index].SlabRuleFrom)
					// log.Info("promoCriterias["+strconv.Itoa(index)+"].SlabRuleTo : ", promoCriterias[index].SlabRuleTo)
					// log.Info("slabRuleValue : ", slabRuleValue)
					if (promoCriterias[index].IsMultiplied && slabRuleValue >= promoCriterias[index].SlabRuleTo) ||
						(!promoCriterias[index].IsMultiplied && slabRuleValue >= promoCriterias[index].SlabRuleFrom && slabRuleValue <= promoCriterias[index].SlabRuleTo) {
						validatedPromoCriterias[promoCriterias[index].PromoID] = promoCriterias[index]
						validatedSlabIDs = append(validatedSlabIDs, strconv.FormatInt(*promoCriterias[index].SlabID, 10))
						isPromoCriteriaValid = true
					}
				}
			}

			if sortedPromoCriterias, err := service.PromotionRepository.FindPromoCriteriasBySlabIDs(validatedSlabIDs); err == nil {
				log.Info("LENGTH VALIDATED PROMO CRITERIA : ", len(validatedPromoCriterias))
				log.Info("PENGHITUNGAN REWARD PRICE / REWARD PRODUCT")
				stockOfProductRewards := make(map[int64]float64)
				for _, promoCriteria := range sortedPromoCriterias {
					log.Info("validatedPromoCriterias["+promoCriteria.PromoID+"] :", validatedPromoCriterias[promoCriteria.PromoID])

					var response entity.ConsultPromotionResponse
					response.PromotionID = promoCriteria.PromoID
					response.PromotionDesc = validatedPromoCriterias[promoCriteria.PromoID].PromoDesc
				response.SlabId = *validatedPromoCriterias[promoCriteria.PromoID].SlabID
				response.SlabDesc = validatedPromoCriterias[promoCriteria.PromoID].SlabDesc

				switch validatedPromoCriterias[promoCriteria.PromoID].SlabRewardType {
				case entity.PromoRewardTypePercent:
					response.SlabReward = math.Round((float64(subTotalValidatedPromoAdditionalCriteriaByProductGroups[promoCriteria.PromoID]) * validatedPromoCriterias[promoCriteria.PromoID].SlabReward) / 100.0)
				case entity.PromoRewardTypeFixedValue:
					response.SlabReward = validatedPromoCriterias[promoCriteria.PromoID].SlabReward
				default:
					response.SlabReward = 0
				}

				slabReward := response.SlabReward
					for proID := range validatedPromoAdditionalCriteriaByProductGroups[promoCriteria.PromoID] {
						response.Products = append(response.Products, proID)

						if response.SlabReward == 0 {
							continue
						}

						log.Info("PENGHITUNGAN REWARD PRICE")
						// log.Info("validatedPromoAdditionalCriteriaByProductGroups["+promoCriteria.PromoID+"]["+strconv.Itoa(proID)+"].SubTotal :", validatedPromoAdditionalCriteriaByProductGroups[promoCriteria.PromoID][proID].SubTotal)
						rewardPrice := entity.ConsultPromotionRewardPriceResponse{}
						rewardPrice.ProID = proID
						rewardPrice.SubTotal = float64(validatedPromoAdditionalCriteriaByProductGroups[promoCriteria.PromoID][proID].SubTotal)

						reward := math.Round((rewardPrice.SubTotal * response.SlabReward) / float64(subTotalValidatedPromoAdditionalCriteriaByProductGroups[promoCriteria.PromoID]))
						slabReward -= reward
						if slabReward <= 0 {
							reward += slabReward
						}
						// log.Info("slabReward "+promoCriteria.PromoID+" :", slabReward)
						rewardPrice.Reward = reward
						rewardPrice.Total = rewardPrice.SubTotal - rewardPrice.Reward

						response.RewardPrice = append(response.RewardPrice, rewardPrice)
					}

					// log.Info("SlabRule", validatedPromoCriterias[promoCriteria.PromoID].SlabRule)
					// log.Info("SlabRuleTo", validatedPromoCriterias[promoCriteria.PromoID].SlabRuleTo)
					log.Info("PENGHITUNGAN REWARD PRODUCT")
					if response.SlabReward == 0 {
						// var convertedTotalQtyReward float64
						rewards, _ := service.PromotionRepository.GetAllRewardProductFromStock(request, validatedPromoCriterias[promoCriteria.PromoID])

						totalQtyStock := float64(0)
						for index, reward := range rewards {
							if _, exists := stockOfProductRewards[reward.ProID]; !exists {
								stockOfProductRewards[reward.ProID] = reward.QtyStock
							}

							convertedQtyStock := float64(0)
							switch validatedPromoCriterias[promoCriteria.PromoID].SlabRewardUom {
							case 3:
								convertedQtyStock = (stockOfProductRewards[reward.ProID] / validatedPromoAdditionalCriteriaByProductGroups[promoCriteria.PromoID][reward.ProID].ConvUnit2) / validatedPromoAdditionalCriteriaByProductGroups[promoCriteria.PromoID][reward.ProID].ConvUnit3
							case 2:
								convertedQtyStock = stockOfProductRewards[reward.ProID] / validatedPromoAdditionalCriteriaByProductGroups[promoCriteria.PromoID][reward.ProID].ConvUnit2
							default:
								convertedQtyStock = stockOfProductRewards[reward.ProID]
							}

							rewards[index].QtyStock = convertedQtyStock
							totalQtyStock += rewards[index].QtyStock
						}

						multipliedValue := int64(1)
						if validatedPromoCriterias[promoCriteria.PromoID].IsMultiplied {
							// log.Info("validatedPromoCriterias["+promoCriteria.PromoID+"].SlabRule : ", validatedPromoCriterias[promoCriteria.PromoID].SlabRule)
							multipliedValue = validatedPromoCriterias[promoCriteria.PromoID].SlabRule / int64(validatedPromoCriterias[promoCriteria.PromoID].SlabRuleTo)
							// log.Info("MULTIPLIED_VALUE : ", multipliedValue)
						}
						totalQtyReward := validatedPromoCriterias[promoCriteria.PromoID].SlabReward * float64(multipliedValue)
						log.Info("TOTAL QTY REWARD : ", totalQtyReward)
						log.Info("TOTAL REWARDS : ", len(rewards))

						response.SlabReward = totalQtyReward
						if totalQtyReward > totalQtyStock {
							if len(request.PromoIDs) > 0 {
								responses = append(responses, response)
							}

							continue
						}

						for _, reward := range rewards {
							log.Info("REWARD : ", reward)
							var rewardProduct entity.ConsultPromotionRewardProductResponse
							log.Info("REWARD QTY STOCK : ", reward.QtyStock)
							qtyReward := totalQtyReward
							if totalQtyReward >= reward.QtyStock {
								qtyReward = reward.QtyStock
							}
							totalQtyReward -= qtyReward

							qty := int64(0)
							convertedQtyReward := qtyReward
							var rewardProductConversion entity.CreateConversionBody
							rewardProductConversion.CustId = request.CustID
							rewardProductConversion.ProductId = reward.ProID

							switch validatedPromoCriterias[promoCriteria.PromoID].SlabRewardUom {
							case 1:
								rewardProductConversion.Qty1 = int64(qtyReward)
								rewardProductConversion.Qty2 = qty
								rewardProductConversion.Qty3 = qty

							case 2:
								rewardProductConversion.Qty1 = qty
								rewardProductConversion.Qty2 = int64(qtyReward)
								rewardProductConversion.Qty3 = qty
								convertedQtyReward = qtyReward * validatedPromoAdditionalCriteriaByProductGroups[promoCriteria.PromoID][reward.ProID].ConvUnit2
							default:
								rewardProductConversion.Qty1 = qty
								rewardProductConversion.Qty2 = qty
								rewardProductConversion.Qty3 = int64(qtyReward)
								convertedQtyReward = qtyReward * validatedPromoAdditionalCriteriaByProductGroups[promoCriteria.PromoID][reward.ProID].ConvUnit2 * validatedPromoAdditionalCriteriaByProductGroups[promoCriteria.PromoID][reward.ProID].ConvUnit3
							}
							rewardProductConversionResut, _ := service.Conversion(rewardProductConversion, request.CustID, request.ParentCustID)

							rewardProduct.ProID = reward.ProID
							rewardProduct.Qty1 = float64(rewardProductConversionResut.Qty1)
							rewardProduct.Qty2 = float64(rewardProductConversionResut.Qty2)
							rewardProduct.Qty3 = float64(rewardProductConversionResut.Qty3)
							// rewardProduct.UnitId = reward.UnitId
							// rewardProduct.Uom = validatedPromoCriterias[promoCriteria.PromoID].SlabRewardUom
							stockOfProductRewards[reward.ProID] -= convertedQtyReward

							response.RewardProduct = append(response.RewardProduct, rewardProduct)

							if totalQtyReward <= 0 {
								break
							}
						}
					}

					if len(response.RewardProduct) > 0 || len(response.RewardPrice) > 0 {
						responses = append(responses, response)
					}
				}
			}

		}

		return responses, nil
	}
*/
func (service *returnServiceImpl) Detail(returnNo string, custID string, parentCustID string) (response entity.ReturnResponse, err error) {
	rtn, err := service.ReturnRepository.FindOneByReturnNo(returnNo, custID, parentCustID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(rtn, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.ReturnRepository.FindReturnDetail(returnNo, custID, parentCustID)
	if err != nil {
		return response, err
	}

	var DetailsData []entity.ReturnDetailResponse
	for _, detail := range Details {
		var detailData entity.ReturnDetailResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		// if response.InvoiceNo != nil {
		// 	returnedProducts, err := service.ReturnRepository.CountReturnedProductQty(*response.InvoiceNo, detail.ProductID, custID)
		// 	if err != nil {
		// 		return response, err
		// 	}

		// 	detailData.RemainingQty1 = detail.InvoiceQty1 - returnedProducts.ReturnedQty1
		// 	detailData.RemainingQty2 = detail.InvoiceQty2 - returnedProducts.ReturnedQty2
		// 	detailData.RemainingQty3 = detail.InvoiceQty3 - returnedProducts.ReturnedQty3
		// }

		remainingQty1 := *detail.RemainingQty1
		remainingQty2 := *detail.RemainingQty2
		remainingQty3 := *detail.RemainingQty3
		convUnit2 := *detail.ConvUnit2
		convUnit3 := *detail.ConvUnit3
		for remainingQty1 < 0 {
			remainingQty2--
			remainingQty1 += convUnit2
		}

		for remainingQty2 < 0 {
			remainingQty3--
			remainingQty2 += convUnit3
		}

		detailData.RemainingQty1 = remainingQty1
		detailData.RemainingQty2 = remainingQty2
		detailData.RemainingQty3 = remainingQty3

		itemConditionName := detailData.GenerateItemConditionName()
		detailData.ItemCndName = &itemConditionName

		DetailsData = append(DetailsData, detailData)
	}

	if rtn.ReturnDate != nil {
		returnDate := rtn.ReturnDate.Format("2006-01-02")
		response.ReturnDate = &returnDate
	}

	if rtn.InvoiceDate != nil {
		invoiceDate := rtn.InvoiceDate.Format("2006-01-02")
		response.InvoiceDate = &invoiceDate
	}

	returnStatusName := response.GenerateReturnStatusName()
	response.DataStatusName = &returnStatusName

	response.Details = DetailsData
	return response, nil
}

/*
	func (service *returnServiceImpl) Delete(custId string, returnNo string, userId int64) (err error) {
		c := context.Background()
		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
			err = service.ReturnRepository.Delete(txCtx, custId, returnNo, userId)
			if err != nil {
				return err
			}
			return nil
		})

		return err
	}
*/

func (service *returnServiceImpl) Update(returnNo string, request entity.UpdateReturnBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		DetailIds := []int64{}

		for _, detail := range request.Details {
			DetailIds = append(DetailIds, detail.ReturnDetailID)
		}
		if len(DetailIds) > 0 {
			if err := service.ReturnRepository.DeleteDetailNotInIDs(txCtx, returnNo, DetailIds, request.CustID); err != nil {
				return err
			}
		}

		var returnModel model.Return
		if err = structs.Automapper(request, &returnModel); err != nil {
			return err
		}

		var order model.OrderList
		var orderDetails []model.OrderDetailRead
		// var qtyOrderDetail float64 = 0
		orderDetailNormalMaps := make(map[int64]*model.OrderDetailRead)
		orderDetailPromoMaps := make(map[string]map[int64]*model.OrderDetailRead)
		orderDetailPromoProductIDMaps := make(map[string][]int64)
		qtyOrderDetail := 0.0
		if returnModel.InvoiceNo != nil {
			order, err = service.OrderRepository.FindByInvoiceNo(*returnModel.InvoiceNo, returnModel.CustID)
			if err != nil {
				return err
			}

			orderDetails, err = service.OrderRepository.FindDetail(order.RoNo, order.CustID)
			if err != nil {
				return err
			}

			for i := range orderDetails {
				detail := &orderDetails[i]
				if detail.ItemType == 1 {
					orderDetailNormalMaps[int64(detail.ProId)] = detail
					qtyOrderDetail += *detail.Qty
				} else {
					if _, isExist := orderDetailPromoMaps[*detail.PromoID]; !isExist {
						orderDetailPromoMaps[*detail.PromoID] = make(map[int64]*model.OrderDetailRead)
					}
					orderDetailPromoMaps[*detail.PromoID][int64(detail.ProId)] = detail
					orderDetailPromoProductIDMaps[*detail.PromoID] = append(orderDetailPromoProductIDMaps[*detail.PromoID], int64(detail.ProId))
				}
			}
			// if orderRewards, err := service.OrderRepository.FindFullPromoRewards(*returnModel.InvoiceNo, returnModel.CustID); err == nil {
			// 	for index, reward := range orderRewards {

			// 	}
			// }
		}

		var whId int64 = 0
		qtyReturnDetail := 0.0
		// var returnDetailModelList []model.ReturnDetail
		for _, detailRequest := range request.Details {
			var returnDetailModel model.ReturnDetail
			if err = structs.Automapper(detailRequest, &returnDetailModel); err != nil {
				return err
			}

			if whId == 0 {
				whId = returnDetailModel.WhId
			}

			promoValue := float64(0)
			promoBgValue := float64(0)
			discValue := float64(0)
			totalQtyReturnDetail := returnDetailModel.Qty1 + (returnDetailModel.Qty2 * *returnDetailModel.ConvUnit2) + (returnDetailModel.Qty3 * *returnDetailModel.ConvUnit3 * *returnDetailModel.ConvUnit2)
			qtyReturnDetail += totalQtyReturnDetail
			if returnDetailModel.OrderDetailID != nil {
				totalQtyOrderDetail := *orderDetailNormalMaps[int64(detailRequest.ProductID)].Qty1 + (*orderDetailNormalMaps[int64(detailRequest.ProductID)].Qty2 * float64(*orderDetailNormalMaps[int64(detailRequest.ProductID)].ConvUnit2)) + (*orderDetailNormalMaps[int64(detailRequest.ProductID)].Qty3 * float64(*orderDetailNormalMaps[int64(detailRequest.ProductID)].ConvUnit3) * float64(*orderDetailNormalMaps[int64(detailRequest.ProductID)].ConvUnit2))

				if orderDetailNormalMaps[int64(detailRequest.ProductID)].ItemType == 1 {
					promoValue = math.Round(float64(*orderDetailNormalMaps[int64(detailRequest.ProductID)].PromoValue) * (totalQtyReturnDetail / totalQtyOrderDetail))
					discValue = math.Round(float64(*orderDetailNormalMaps[int64(detailRequest.ProductID)].DiscValue) * (totalQtyReturnDetail / totalQtyOrderDetail))
				}

				qtySisa := totalQtyOrderDetail - totalQtyReturnDetail
				orderDetailNormalMaps[int64(detailRequest.ProductID)].Qty = &qtySisa

				qtySisaConversion := &conversion.Qty{
					Qty:       int(getValueOrDefault(&qtySisa, 0)),
					ConvUnit2: int(*orderDetailNormalMaps[int64(detailRequest.ProductID)].ConvUnit2),
					ConvUnit3: int(*orderDetailNormalMaps[int64(detailRequest.ProductID)].ConvUnit3),
				}

				qtyConversion := qtySisaConversion.ConvToQtyConversion()

				orderDetailQty1 := float64(qtyConversion.Qty1)
				orderDetailQty2 := float64(qtyConversion.Qty2)
				orderDetailQty3 := float64(qtyConversion.Qty3)
				orderDetailNormalMaps[int64(detailRequest.ProductID)].Qty1 = &orderDetailQty1
				orderDetailNormalMaps[int64(detailRequest.ProductID)].Qty2 = &orderDetailQty2
				orderDetailNormalMaps[int64(detailRequest.ProductID)].Qty3 = &orderDetailQty3
			}
			returnDetailModel.PromoValue = &promoValue
			returnDetailModel.DiscValue = &discValue

			returnDetailModel.SubTotal = (returnDetailModel.Qty1 * returnDetailModel.SellPrice1) + (returnDetailModel.Qty2 * returnDetailModel.SellPrice2) + (returnDetailModel.Qty3 * returnDetailModel.SellPrice3)
			returnDetailModel.VatValue = math.Round((returnDetailModel.SubTotal - *returnDetailModel.PromoValue - *returnDetailModel.DiscValue) * (returnDetailModel.Vat / 100.0))
			returnDetailModel.Total = returnDetailModel.SubTotal - *returnDetailModel.PromoValue - *returnDetailModel.DiscValue + returnDetailModel.VatValue
			returnModel.SubTotal += returnDetailModel.SubTotal
			returnModel.VatValue += returnDetailModel.VatValue
			returnModel.DiscValue += discValue
			returnModel.PromoValue += promoValue
			returnModel.PromoBgValue += promoBgValue

			// returnDetailModel.SubTotal = (returnDetailModel.Qty1 * returnDetailModel.SellPrice1) + (returnDetailModel.Qty2 * returnDetailModel.SellPrice2) + (returnDetailModel.Qty3 * returnDetailModel.SellPrice3)
			// returnDetailModel.VatValue = returnDetailModel.SubTotal * (returnDetailModel.Vat / 100.0)
			// returnDetailModel.Total = returnDetailModel.SubTotal + returnDetailModel.VatValue
			// Model.SubTotal += returnDetailModel.SubTotal
			// Model.VatValue += returnDetailModel.VatValue

			returnDetailModel.CustID = ""
			if err = service.ReturnRepository.UpdateDetail(txCtx, &returnDetailModel); err != nil {
				return err
			}
		}

		var promoReturnDetailModelList []model.ReturnDetail
		if qtyOrderDetail == qtyReturnDetail {
			for promoID := range orderDetailPromoMaps {
				for _, orderDetailPromoMap := range orderDetailPromoMaps[promoID] {
					var returnDetailModel model.ReturnDetail
					if err = structs.Automapper(orderDetailPromoMap, &returnDetailModel); err != nil {
						return err
					}

					// returnDetailModel.PromoValue = &promoValue
					// returnDetailModel.DiscValue = &discValue

					returnDetailModel.CustID = returnModel.CustID
					returnDetailModel.ProductID = orderDetailPromoMap.ProId
					returnDetailModel.WhId = whId
					returnDetailModel.ItemCnd = 1
					returnDetailModel.ReturnReasonID = 0
					returnDetailModel.SubTotal = (returnDetailModel.Qty1 * returnDetailModel.SellPrice1) + (returnDetailModel.Qty2 * returnDetailModel.SellPrice2) + (returnDetailModel.Qty3 * returnDetailModel.SellPrice3)
					// returnDetailModel.VatValue = math.Round((returnDetailModel.SubTotal - *returnDetailModel.PromoValue - *returnDetailModel.DiscValue) * (returnDetailModel.Vat / 100.0))
					returnDetailModel.Total = *orderDetailPromoMap.Amount
					returnModel.SubTotal += returnDetailModel.SubTotal
					// returnModel.VatValue += returnDetailModel.VatValue
					// returnModel.DiscValue += discValue
					// returnModel.PromoValue += promoValue
					returnModel.PromoBgValue += *orderDetailPromoMap.PromoValue

					promoReturnDetailModelList = append(promoReturnDetailModelList, returnDetailModel)
				}
			}
		}

		returnModel.CustID = request.CustID
		returnModel.ReturnNo = returnNo
		returnModel.Total = returnModel.SubTotal - returnModel.PromoValue - returnModel.PromoBgValue - returnModel.DiscValue + returnModel.VatValue

		if err = service.ReturnRepository.Update(txCtx, &returnModel); err != nil {
			return err
		}

		for _, returnDetail := range promoReturnDetailModelList {
			returnDetail.ReturnNo = returnModel.ReturnNo
			if err = service.ReturnRepository.StoreDetail(txCtx, &returnDetail); err != nil {
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

func (service *returnServiceImpl) UpdateQuantity(returnNo string, request entity.UpdateQuantityReturnBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		DetailIds := []int64{}

		for _, detail := range request.Details {
			DetailIds = append(DetailIds, detail.ReturnDetailID)
		}
		if len(DetailIds) > 0 {
			err := service.ReturnRepository.DeleteDetailNotInIDs(txCtx, returnNo, DetailIds, request.CustID)
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
			returnDetailModel.VatValue = returnDetailModel.SubTotal * (returnDetailModel.Vat / 100.0)
			returnDetailModel.Total = returnDetailModel.SubTotal + returnDetailModel.VatValue
			Model.SubTotal += returnDetailModel.SubTotal
			Model.VatValue += returnDetailModel.VatValue

			returnDetailModel.CustID = ""
			err = service.ReturnRepository.UpdateQuantity(txCtx, &returnDetailModel)
			if err != nil {
				return err
			}
		}

		Model.CustID = request.CustID
		Model.ReturnNo = returnNo
		Model.Total = Model.SubTotal + Model.VatValue
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

func (service *returnServiceImpl) Approve(returnNo string, request entity.ApproveReturnBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		DetailIds := []int64{}

		for _, detail := range request.Details {
			DetailIds = append(DetailIds, detail.ReturnDetailID)
		}
		if len(DetailIds) > 0 {
			err := service.ReturnRepository.DeleteDetailNotInIDs(txCtx, returnNo, DetailIds, request.CustID)
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

			var returnDetailModel model.ReturnDetail
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
			err = service.ReturnRepository.UpdateDetail(txCtx, &returnDetailModel)
			if err != nil {
				return err
			}
		}

		Model.CustID = request.CustID
		Model.ReturnNo = returnNo
		Model.Total = Model.SubTotal - Model.VatValue

		var reviewedBy = request.UpdatedBy
		var reviewedAt = time.Now()

		Model.DataStatus = 3
		Model.IsReviewed = true
		Model.ReviewedBy = &reviewedBy
		Model.ReviewedAt = &reviewedAt
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

func (service *returnServiceImpl) Cancel(returnNo string, request entity.CancelReturnBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		var Model model.Return
		err = structs.Automapper(request, &Model)
		if err != nil {
			return err
		}

		Model.CustID = request.CustID
		Model.ReturnNo = returnNo

		var reviewedBy = request.UpdatedBy
		var reviewedAt = time.Now()

		Model.DataStatus = 9
		Model.IsReviewed = true
		Model.ReviewedBy = &reviewedBy
		Model.ReviewedAt = &reviewedAt
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

func (service *returnServiceImpl) UpdateStatus(request entity.UpdateStatusReturnBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		var Model model.Return
		err = structs.Automapper(request, &Model)
		if err != nil {
			return err
		}
		Model.CustID = request.CustID

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

func (service *returnServiceImpl) UpdateAssign(request entity.UpdateAssignReturnBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		var Model model.Return
		err = structs.Automapper(request, &Model)
		if err != nil {
			return err
		}
		Model.CustID = request.CustID

		for _, detail := range request.Returns {

			Model.ReturnNo = detail.ReturnNo
			Model.EmpId = detail.EmpId
			Model.DataStatus = 7

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

func (service *returnServiceImpl) SalesmanFilterLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.SalesmansLookupResponse, total int64, lastPage int, err error) {
	Salesmans, total, lastPage, err := service.ReturnRepository.FindAllSalesmanFilterByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Salesmans {
		var vResp entity.SalesmansLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnServiceImpl) EmployeeFilterLookupList(dataFilter entity.SalesmanQueryFilter) (data []entity.EmployeeLookupResponse, total int64, lastPage int, err error) {
	empGroups, total, lastPage, err := service.ReturnRepository.FindAllEmployeeByEmpGrpIdFilterByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range empGroups {
		var vResp entity.EmployeeLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err

}
func (service *returnServiceImpl) RoleFilterLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.EmpGroupLookupResponse, total int64, lastPage int, err error) {
	empGroups, total, lastPage, err := service.ReturnRepository.FindAllRolesFilterByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range empGroups {
		var vResp entity.EmpGroupLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err

}

func (service *returnServiceImpl) OutletFilterLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.OutletsLookupResponse, total int64, lastPage int, err error) {
	Outlets, total, lastPage, err := service.ReturnRepository.FindAllOutletFilterByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Outlets {
		var vResp entity.OutletsLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnServiceImpl) ReturnStatusesLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.ReturnStatusesLookupResponse, total int64, lastPage int, err error) {
	DepositStatuses, total, lastPage, err := service.ReturnRepository.FindAllReturnStatusesLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, depositStatus := range DepositStatuses {
		var vResp entity.ReturnStatusesLookupResponse
		structs.Automapper(depositStatus, &vResp)

		depositStatusName := vResp.GenerateDataReturnStatusName()
		vResp.ReturnStatusName = &depositStatusName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnServiceImpl) SalesmanFilterLookupCreate(dataFilter entity.GeneralQueryFilter) (data []entity.SalesmansLookupResponse, total int64, lastPage int, err error) {
	Salesmans, total, lastPage, err := service.ReturnRepository.FindAllSalesmanFilterByCustIdLookupModeCreate(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Salesmans {
		var vResp entity.SalesmansLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnServiceImpl) OutletFilterLookupCreate(dataFilter entity.OutletCreateReturnQueryFilter) (data []entity.OutletsLookupResponse, total int64, lastPage int, err error) {
	Outlets, total, lastPage, err := service.ReturnRepository.FindAllOutletFilterByCustIdLookupModeCreate(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Outlets {
		var vResp entity.OutletsLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnServiceImpl) ProductListCreate(dataFilter entity.ProductListQueryFilter) (data []entity.ProductListResponse, total int64, lastPage int, err error) {
	products, total, lastPage, err := service.ReturnRepository.FindAllProductByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range products {
		var vResp entity.ProductListResponse
		structs.Automapper(row, &vResp)
		if row.InvoiceDate != nil {
			invoiceDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = &invoiceDate
		}

		remainingQty1 := row.RemainingQty1
		remainingQty2 := row.RemainingQty2
		remainingQty3 := row.RemainingQty3
		convUnit2 := *row.ConvUnit2
		convUnit3 := *row.ConvUnit3
		for remainingQty1 < 0 {
			remainingQty2--
			remainingQty1 += convUnit2
		}

		for remainingQty2 < 0 {
			remainingQty3--
			remainingQty2 += convUnit3
		}

		vResp.RemainingQty1 = remainingQty1
		vResp.RemainingQty2 = remainingQty2
		vResp.RemainingQty3 = remainingQty3

		vResp.SubTotal1 = vResp.RemainingQty1 * vResp.SellPrice1
		vResp.SubTotal2 = vResp.RemainingQty2 * vResp.SellPrice2
		vResp.SubTotal3 = vResp.RemainingQty3 * vResp.SellPrice3
		vResp.Total = vResp.SubTotal1 + vResp.SubTotal2 + vResp.SubTotal3

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnServiceImpl) ProductListCreateOld(dataFilter entity.ProductListQueryFilter) (data []entity.ProductListResponse, total int64, lastPage int, err error) {
	products, total, lastPage, err := service.ReturnRepository.FindAllProductByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range products {
		var vResp entity.ProductListResponse
		structs.Automapper(row, &vResp)
		if row.InvoiceDate != nil {
			invoiceDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = &invoiceDate
		}

		returnedProducts, err := service.ReturnRepository.CountReturnedProductQty(*vResp.InvoiceNo, row.ProductId, dataFilter.CustId)
		if err != nil {
			return data, total, lastPage, err
		}

		vResp.RemainingQty1 = row.InvoiceQty1 - returnedProducts.ReturnedQty1
		vResp.RemainingQty2 = row.InvoiceQty2 - returnedProducts.ReturnedQty2
		vResp.RemainingQty3 = row.InvoiceQty3 - returnedProducts.ReturnedQty3

		vResp.SubTotal1 = vResp.RemainingQty1 * vResp.SellPrice1
		vResp.SubTotal2 = vResp.RemainingQty2 * vResp.SellPrice2
		vResp.SubTotal3 = vResp.RemainingQty3 * vResp.SellPrice3
		vResp.Total = vResp.SubTotal1 + vResp.SubTotal2 + vResp.SubTotal3

		data = append(data, vResp)
	}

	return data, total, lastPage, err
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

func (service *returnServiceImpl) WarehouseLookupList(dataFilter entity.WarehouseQueryFilter) (data []entity.WarehousesLookupResponse, total int64, lastPage int, err error) {
	if dataFilter.ItemCnd == 1 {
		sales, err := service.ReturnRepository.FindSalesmanById(dataFilter.SalesmanId, dataFilter.CustId, dataFilter.ParentCustId)
		if err != nil {
			return data, 0, 0, err
		}

		for _, row := range sales {
			dataFilter.WhId = append(dataFilter.WhId, row.WhId)
		}
	} else {
		if dataFilter.ItemCnd == 2 {
			dataFilter.StockType = "BS"
		}
		if dataFilter.ItemCnd == 3 {
			dataFilter.StockType = "E"
		}
	}

	Warehouses, total, lastPage, err := service.ReturnRepository.FindAllMasterWarehouseLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Warehouses {
		var vResp entity.WarehousesLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnServiceImpl) ProductLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.ProductsLookupCreateResponse, total int64, lastPage int, err error) {
	products, total, lastPage, err := service.ReturnRepository.FindAllMasterProductByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range products {
		var vResp entity.ProductsLookupCreateResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *returnServiceImpl) Print(custId string, returnNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ReturnRepository.Print(txCtx, custId, returnNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

/*
func (service *returnServiceImpl) ProductConditionLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.ProductConditionsLookupResponse, total int64, lastPage int, err error) {
	// productConditions, total, lastPage, err := service.ReturnRepository.FindAllReturnStatusesLookupMode(dataFilter)
	// if err != nil {
	// 	return data, total, lastPage, err
	// }
	var proCon entity.ProductConditionsLookupResponse
	proCons := proCon.GetProductConditionList
	// for id, dataProductCondition := range  {
	// 	proCon.ProductConditionId = int(id)
	// 	proCon.ProductConditionName = dataProductCondition

	// 	data = append(data, proCon)
	// }
	log.Info("Procons : ", proCons)

	return data, total, lastPage, err
}
*/
