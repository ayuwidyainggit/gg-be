package service

import (
	"context"
	"errors"
	"math"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/constant"
	"mobile/pkg/str"
	"mobile/pkg/structs"
	"mobile/repository"

	"github.com/gofiber/fiber/v2/log"
)

type DiscountService interface {
	Store(request entity.CreateDiscountBody) (err error)
	Detail(params entity.DetailDiscountParams) (response entity.Discount, err error)
	DetailGrp(DiscGrpId string) (data []entity.DetailDiscountGrp, err error)
	List(dataFilter entity.DiscountQueryFilter) (data []entity.Discount, total int64, lastPage int, err error)
	Update(discountID string, request entity.UpdateDiscountBody) (err error)
	Delete(params entity.DetailDiscountParams, deletedBy string) (err error)
	PublishDiscount(equest entity.PublishDiscountBody) (err error)
	ConsultDiscount(equest entity.ConsultDiscountBody) (responses []entity.ConsultDiscountResponse, err error)
	// ConsultDiscountNew(equest entity.ConsultDiscountBody) (responses []entity.ConsultDiscountResponseNew, err error)
}

func NewDiscountService(discountRepository repository.DiscountRepository, transaction repository.Dbtransaction) *discountServiceImpl {
	return &discountServiceImpl{
		DiscountRepository: discountRepository,
		Transaction:        transaction,
	}
}

type discountServiceImpl struct {
	DiscountRepository repository.DiscountRepository
	Transaction        repository.Dbtransaction
}

func (service *discountServiceImpl) Store(request entity.CreateDiscountBody) (err error) {
	c := context.Background()

	effectiveFrom, err := str.DateStrToRfc3339String(request.EffectiveFrom)
	if err != nil {
		return err
	}
	request.EffectiveFrom = effectiveFrom

	effectiveTo, err := str.DateStrToRfc3339String(request.EffectiveTo)
	if err != nil {
		return err
	}
	request.EffectiveTo = effectiveTo

	var discountModel model.Discount
	err = structs.Automapper(request, &discountModel)
	if err != nil {
		return err
	}
	discountModel.DiscountStatusID = 1 // make it default status 'Inactive'
	discountModel.PublishStatusID = 1  // make it default status 'New'

	// var discountStatusLogModel model.DiscountStatusLog
	// err = structs.Automapper(request, &discountStatusLogModel)
	// if err != nil {
	// 	return err
	// }

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.DiscountRepository.Store(txCtx, &discountModel)
		if err != nil {
			log.Error("err:", err.Error())
			return err
		}

		for _, row := range request.DiscountPrincipals {
			var discountPrincipalModel model.DiscountPrincipal
			err := structs.Automapper(row, &discountPrincipalModel)
			if err != nil {
				return err
			}
			discountPrincipalModel.CustID = request.CustID
			discountPrincipalModel.DiscountID = request.DiscountID
			// log.Info("discountCriteriaModel:", structs.StructToJson(discountCriteriaModel))
			err = service.DiscountRepository.StoreDiscountPrincipal(txCtx, &discountPrincipalModel)
			if err != nil {
				return err
			}

		}

		for _, row := range request.DiscountGroups {
			var discountGroupModel model.DiscountGroup
			err := structs.Automapper(row, &discountGroupModel)
			if err != nil {
				return err
			}
			discountGroupModel.CustID = request.CustID
			discountGroupModel.DiscountID = request.DiscountID
			// log.Info("discountCriteriaModel:", structs.StructToJson(discountCriteriaModel))
			err = service.DiscountRepository.StoreDiscountGroup(txCtx, &discountGroupModel)
			if err != nil {
				return err
			}

		}

		for _, row := range request.DiscountCriterias {
			var discountCriteriaModel model.DiscountCriteria
			err := structs.Automapper(row, &discountCriteriaModel)
			if err != nil {
				return err
			}
			discountCriteriaModel.CustID = request.CustID
			discountCriteriaModel.DiscountID = request.DiscountID
			// log.Info("discountCriteriaModel:", structs.StructToJson(discountCriteriaModel))
			err = service.DiscountRepository.StoreDiscountCriteria(txCtx, &discountCriteriaModel)
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

func (service *discountServiceImpl) Detail(params entity.DetailDiscountParams) (response entity.Discount, err error) {
	discount, err := service.DiscountRepository.FindByDiscountID(params)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(discount, &response)
	if err != nil {
		return response, err
	}

	response.CustID = ""
	response.EffectiveFrom = discount.EffectiveFrom.Format(constant.YYYY_MM_DD)
	response.EffectiveTo = discount.EffectiveTo.Format(constant.YYYY_MM_DD)

	response.DiscountStatusDesc = response.GetDiscountStatusDesc()
	response.PublishStatusDesc = response.GetPublishStatusDesc()

	err = service.GetDiscountPrincipals(params, &response)
	if err != nil {
		return response, err
	}

	err = service.GetDiscountGroups(params, &response)
	if err != nil {
		return response, err
	}

	err = service.GetDiscountCriterias(params, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (service *discountServiceImpl) List(dataFilter entity.DiscountQueryFilter) (data []entity.Discount, total int64, lastPage int, err error) {
	discounts, total, lastPage, err := service.DiscountRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range discounts {
		var vResp entity.Discount
		structs.Automapper(row, &vResp)

		vResp.CustID = ""
		vResp.EffectiveFrom = row.EffectiveFrom.Format(constant.YYYY_MM_DD)
		vResp.EffectiveTo = row.EffectiveTo.Format(constant.YYYY_MM_DD)
		vResp.DiscountStatusDesc = vResp.GetDiscountStatusDesc()
		vResp.PublishStatusDesc = vResp.GetPublishStatusDesc()

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *discountServiceImpl) Update(discountID string, request entity.UpdateDiscountBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	effectiveFrom, err := str.DateStrToRfc3339String(request.EffectiveFrom)
	if err != nil {
		return err
	}
	request.EffectiveFrom = effectiveFrom

	effectiveTo, err := str.DateStrToRfc3339String(request.EffectiveTo)
	if err != nil {
		return err
	}
	request.EffectiveTo = effectiveTo
	// End parse time format YYYY-mm-dd to Rfc339

	var discountModel model.Discount
	err = structs.Automapper(request, &discountModel)
	if err != nil {
		return err
	}

	discountModel.CustID = ""
	discountModel.PublishStatusID = 0
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.DiscountRepository.Update(txCtx, discountID, discountModel)
		if err != nil {
			return err
		}

		err := service.DiscountRepository.DeleteDiscountPrincipals(txCtx, request.CustID, discountID)
		if err != nil {
			return err
		}

		for _, row := range request.DiscountPrincipals {

			var discountPrincipalModel model.DiscountPrincipal
			err = structs.Automapper(row, &discountPrincipalModel)
			if err != nil {
				return err
			}
			discountPrincipalModel.CustID = request.CustID
			discountPrincipalModel.DiscountID = discountID

			err = service.DiscountRepository.StoreDiscountPrincipal(txCtx, &discountPrincipalModel)
			if err != nil {
				return err
			}

		}

		err = service.DiscountRepository.DeleteDiscountGroups(txCtx, request.CustID, discountID)
		if err != nil {
			return err
		}

		for _, row := range request.DiscountGroups {

			var discountGroupModel model.DiscountGroup
			err = structs.Automapper(row, &discountGroupModel)
			if err != nil {
				return err
			}
			discountGroupModel.CustID = request.CustID
			discountGroupModel.DiscountID = discountID

			err = service.DiscountRepository.StoreDiscountGroup(txCtx, &discountGroupModel)
			if err != nil {
				return err
			}

		}

		err = service.DiscountRepository.DeleteDiscountCriterias(txCtx, request.CustID, discountID)
		if err != nil {
			return err
		}

		for _, row := range request.DiscountCriterias {

			var discountCritModel model.DiscountCriteria
			err = structs.Automapper(row, &discountCritModel)
			if err != nil {
				return err
			}
			discountCritModel.CustID = request.CustID
			discountCritModel.DiscountID = discountID

			err = service.DiscountRepository.StoreDiscountCriteria(txCtx, &discountCritModel)
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

func (service *discountServiceImpl) Delete(params entity.DetailDiscountParams, deletedBy string) (err error) {
	c := context.Background()

	discount, err := service.DiscountRepository.FindByDiscountID(params)
	if err != nil {
		return err
	}

	if discount.DiscountStatusID != 1 {
		return errors.New("the discount is not allow to be delete")
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.DiscountRepository.Delete(txCtx, params.CustID, params.DiscountID)
		if err != nil {
			return err
		}

		err = service.DiscountRepository.DeleteDiscountPrincipals(txCtx, params.CustID, params.DiscountID)
		if err != nil {
			log.Error("DeleteDiscountPrincipals, error:", err.Error())
		}

		err = service.DiscountRepository.DeleteDiscountGroups(txCtx, params.CustID, params.DiscountID)
		if err != nil {
			log.Error("DeleteDiscountGroups, error:", err.Error())
		}

		err = service.DiscountRepository.DeleteDiscountCriterias(txCtx, params.CustID, params.DiscountID)
		if err != nil {
			log.Error("DeleteDiscountCriterias, error:", err.Error())
		}

		return nil
	})

	return err
}

func (service *discountServiceImpl) GetDiscountPrincipals(params entity.DetailDiscountParams, discountResponse *entity.Discount) (err error) {
	discountPrincipals, err := service.DiscountRepository.FindAllDiscountPrincipalsByDiscountID(params)
	if err != nil {
		return err
	}

	for _, row := range discountPrincipals {
		var discountPrincipal entity.DiscountPrincipal
		err = structs.Automapper(row, &discountPrincipal)
		if err != nil {
			return err
		}
		discountPrincipal.CustID = ""
		discountPrincipal.DiscountID = ""
		discountResponse.DiscountPrincipals = append(discountResponse.DiscountPrincipals, discountPrincipal)
	}

	return
}

func (service *discountServiceImpl) GetDiscountGroups(params entity.DetailDiscountParams, discountResponse *entity.Discount) (err error) {
	discountGroups, err := service.DiscountRepository.FindAllDiscountGroupsByDiscountID(params)
	if err != nil {
		return err
	}

	for _, row := range discountGroups {
		var discountGroup entity.DiscountGroup
		err = structs.Automapper(row, &discountGroup)
		if err != nil {
			return err
		}
		discountGroup.CustID = ""
		discountGroup.DiscountID = ""
		discountResponse.DiscountGroups = append(discountResponse.DiscountGroups, discountGroup)
	}

	return
}

func (service *discountServiceImpl) GetDiscountCriterias(params entity.DetailDiscountParams, discountResponse *entity.Discount) (err error) {
	discountCriterias, err := service.DiscountRepository.FindAllDiscountCriteriasByDiscountID(params)
	if err != nil {
		return err
	}

	for _, row := range discountCriterias {
		var discountCriteria entity.DiscountCriteria
		err = structs.Automapper(row, &discountCriteria)
		if err != nil {
			return err
		}
		discountCriteria.CustID = ""
		discountCriteria.DiscountID = ""

		discountCriteria.SlabRewardTypeName = constant.GetQtyAmountPercentDisplayName(discountCriteria.SlabRewardType)

		discountResponse.DiscountCriterias = append(discountResponse.DiscountCriterias, discountCriteria)
	}

	return
}

func (service *discountServiceImpl) PublishDiscount(request entity.PublishDiscountBody) (err error) {
	c := context.Background()

	discounts, err := service.DiscountRepository.FindAllByCustIdAndDiscountID(request)
	if err != nil {
		return err
	}

	if len(discounts) < len(request.DiscountID) {
		return errors.New("some of discount id not found")
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.DiscountRepository.PublishDiscount(txCtx, request)
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

func (service *discountServiceImpl) ConsultDiscount(request entity.ConsultDiscountBody) (responses []entity.ConsultDiscountResponse, err error) {
	outlet, err := service.DiscountRepository.FindOutletByID(request.OutletId, request.CustID, request.ParentCustID)
	if err != nil {
		return responses, err
	}

	err = structs.Automapper(request.Details, &responses)
	if err != nil {
		return responses, err
	}
	log.Info("Panjang response : ", len(responses))
	SubTotalPrincipals := map[string]int{}
	TotalProductPerDiscount := map[string]int{}
	for index := range responses {
		product, err := service.DiscountRepository.FindProductByID(responses[index].ProID)
		if err != nil {
			return responses, err
		}

		if discount, err := service.DiscountRepository.FindDiscountByProductAndOutlet(product, outlet, request); err == nil {
			responses[index].DiscountID = discount.DiscountId
			if subTotal, isExist := SubTotalPrincipals[responses[index].DiscountID]; isExist {
				SubTotalPrincipals[responses[index].DiscountID] = subTotal + responses[index].SubTotal
				TotalProductPerDiscount[responses[index].DiscountID] += 1
			} else {
				SubTotalPrincipals[responses[index].DiscountID] = responses[index].SubTotal
				TotalProductPerDiscount[responses[index].DiscountID] = 1
			}

			var param entity.DetailDiscountParams
			param.CustID = request.CustID
			param.ParentCustId = request.ParentCustID
			param.DiscountID = responses[index].DiscountID

			if discountPrincipals, err := service.DiscountRepository.FindAllDiscountPrincipalsByDiscountID(param); err == nil && len(discountPrincipals) > 0 {
				for _, discountPrincipal := range discountPrincipals {
					responses[index].PrincipalID = append(responses[index].PrincipalID, int(discountPrincipal.PrincipalID))
				}
			}
		}
	}
	// log.Info("SubTotalPrincipals Awal : ", SubTotalPrincipals)
	// log.Info("TotalProductPerDiscount Awal : ", TotalProductPerDiscount)
	// log.Info("Responses Awal : ", responses)

	slabRewards := map[string]int{}
	decreaseSlabRewards := map[string]int{}
	discountCriterias := map[string]model.DiscountCriteria{}
	for discountID, SubTotalPrincipal := range SubTotalPrincipals {
		if discountCriteria, err := service.DiscountRepository.FindDiscountCriteriaBySubTotal(discountID, SubTotalPrincipal); err == nil {
			slabReward := discountCriteria.SlabReward
			if discountCriteria.SlabRewardType == 2 {
				slabReward = math.Round((float64(slabReward) * float64(SubTotalPrincipal)) / 100)
			}
			slabRewards[discountID] = int(slabReward)
			decreaseSlabRewards[discountID] = slabRewards[discountID]
			discountCriterias[discountID] = discountCriteria
		}
	}
	// log.Info("slabRewards Awal : ", slabRewards)
	// log.Info("decreaseSlabRewards Awal : ", decreaseSlabRewards)
	// log.Info("discountCriterias Awal : ", discountCriterias)
	// log.Info("SubTotalPrincipals Awal : ", SubTotalPrincipals)
	// log.Info("TotalProductPerDiscount Awal : ", TotalProductPerDiscount)

	for index := range responses {
		if responses[index].DiscountID == "" {
			responses[index].SubTotalPrincipal = responses[index].SubTotal
		} else {
			responses[index].SubTotalPrincipal = SubTotalPrincipals[responses[index].DiscountID]
			responses[index].SlabDesc = discountCriterias[responses[index].DiscountID].SlabDesc
			responses[index].SlabReward = slabRewards[responses[index].DiscountID]

			rewardProduct := int(math.Round((float64(slabRewards[responses[index].DiscountID]) * float64(responses[index].SubTotal)) / float64(SubTotalPrincipals[responses[index].DiscountID])))
			TotalProductPerDiscount[responses[index].DiscountID]--
			if TotalProductPerDiscount[responses[index].DiscountID] <= 0 {
				rewardProduct = decreaseSlabRewards[responses[index].DiscountID]
			}
			decreaseSlabRewards[responses[index].DiscountID] -= rewardProduct
			responses[index].RewardProduct = rewardProduct
		}

		// log.Info("TotalProductPerDiscount[", responses[index].DiscountID, "]  : ", TotalProductPerDiscount[responses[index].DiscountID])
		// log.Info("decreaseSlabRewards[", responses[index].DiscountID, "] : ", decreaseSlabRewards[responses[index].DiscountID])
		// log.Info("Response ", index, " : ", responses[index])
	}

	return responses, nil
}

/*
func (service *discountServiceImpl) ConsultDiscountNew(request entity.ConsultDiscountBody) (responses []entity.ConsultDiscountResponse, err error) {
	outlet, err := service.DiscountRepository.FindOutletByID(request.OutletId, request.CustID, request.ParentCustID)
	if err != nil {
		return responses, err
	}

	err = structs.Automapper(request.Details, &responses)
	if err != nil {
		return responses, err
	}
	log.Info("Panjang response : ", len(responses))
	SubTotalPrincipals := map[string]int{}
	TotalProductPerDiscount := map[string]int{}
	for index := range responses {
		product, err := service.DiscountRepository.FindProductByID(responses[index].ProID)
		if err != nil {
			return responses, err
		}

		if discount, err := service.DiscountRepository.FindDiscountByProductAndOutlet(product, outlet); err == nil {
			responses[index].DiscountID = discount.DiscountId
			if subTotal, isExist := SubTotalPrincipals[responses[index].DiscountID]; isExist {
				SubTotalPrincipals[responses[index].DiscountID] = subTotal + responses[index].SubTotal
				TotalProductPerDiscount[responses[index].DiscountID] += 1
			} else {
				SubTotalPrincipals[responses[index].DiscountID] = responses[index].SubTotal
				TotalProductPerDiscount[responses[index].DiscountID] = 1
			}

			var param entity.DetailDiscountParams
			param.CustID = request.CustID
			param.ParentCustId = request.ParentCustID
			param.DiscountID = responses[index].DiscountID

			if discountPrincipals, err := service.DiscountRepository.FindAllDiscountPrincipalsByDiscountID(param); err == nil && len(discountPrincipals) > 0 {
				for _, discountPrincipal := range discountPrincipals {
					responses[index].PrincipalID = append(responses[index].PrincipalID, int(discountPrincipal.PrincipalID))
				}
			}
		}
	}
	// log.Info("SubTotalPrincipals Awal : ", SubTotalPrincipals)
	// log.Info("TotalProductPerDiscount Awal : ", TotalProductPerDiscount)
	// log.Info("Responses Awal : ", responses)

	slabRewards := map[string]int{}
	decreaseSlabRewards := map[string]int{}
	discountCriterias := map[string]model.DiscountCriteria{}
	for discountID, SubTotalPrincipal := range SubTotalPrincipals {
		if discountCriteria, err := service.DiscountRepository.FindDiscountCriteriaBySubTotal(discountID, SubTotalPrincipal); err == nil {
			slabReward := discountCriteria.SlabReward
			if discountCriteria.SlabRewardType == 2 {
				slabReward = math.Round((float64(slabReward) * float64(SubTotalPrincipal)) / 100)
			}
			slabRewards[discountID] = int(slabReward)
			decreaseSlabRewards[discountID] = slabRewards[discountID]
			discountCriterias[discountID] = discountCriteria
		}
	}
	// log.Info("slabRewards Awal : ", slabRewards)
	// log.Info("decreaseSlabRewards Awal : ", decreaseSlabRewards)
	// log.Info("discountCriterias Awal : ", discountCriterias)
	// log.Info("SubTotalPrincipals Awal : ", SubTotalPrincipals)
	// log.Info("TotalProductPerDiscount Awal : ", TotalProductPerDiscount)

	for index := range responses {
		if responses[index].DiscountID == "" {
			responses[index].SubTotalPrincipal = responses[index].SubTotal
		} else {
			responses[index].SubTotalPrincipal = SubTotalPrincipals[responses[index].DiscountID]
			responses[index].SlabDesc = discountCriterias[responses[index].DiscountID].SlabDesc
			responses[index].SlabReward = slabRewards[responses[index].DiscountID]

			rewardProduct := int(math.Round((float64(slabRewards[responses[index].DiscountID]) * float64(responses[index].SubTotal)) / float64(SubTotalPrincipals[responses[index].DiscountID])))
			TotalProductPerDiscount[responses[index].DiscountID]--
			if TotalProductPerDiscount[responses[index].DiscountID] <= 0 {
				rewardProduct = decreaseSlabRewards[responses[index].DiscountID]
			}
			decreaseSlabRewards[responses[index].DiscountID] -= rewardProduct
			responses[index].RewardProduct = rewardProduct
		}

		// log.Info("TotalProductPerDiscount[", responses[index].DiscountID, "]  : ", TotalProductPerDiscount[responses[index].DiscountID])
		// log.Info("decreaseSlabRewards[", responses[index].DiscountID, "] : ", decreaseSlabRewards[responses[index].DiscountID])
		// log.Info("Response ", index, " : ", responses[index])
	}


		log.Info("PENGHITUNGAN REWARD PRICE / REWARD PRODUCT")
			for promoID := range validatedPromoCriterias {
				log.Info("validatedPromoCriterias["+promoID+"] :", validatedPromoCriterias[promoID])

				var response entity.ConsultPromotionResponse
				response.PromotionID = promoID
				response.PromotionDesc = validatedPromoCriterias[promoID].PromoDesc
				response.SlabId = *validatedPromoCriterias[promoID].SlabID
				response.SlabDesc = validatedPromoCriterias[promoID].SlabDesc
				switch validatedPromoCriterias[promoID].SlabRewardType {
				case 3:
					response.SlabReward = math.Round((float64(subTotalValidatedPromoAdditionalCriteriaByProductGroups[promoID]) * validatedPromoCriterias[promoID].SlabReward) / 100.0)
				case 2:
					response.SlabReward = validatedPromoCriterias[promoID].SlabReward
				default:
					response.SlabReward = 0
				}

				slabReward := response.SlabReward
				for proID := range validatedPromoAdditionalCriteriaByProductGroups[promoID] {
					response.Products = append(response.Products, proID)

					if response.SlabReward == 0 {
						continue
					}

					log.Info("PENGHITUNGAN REWARD PRICE")
					// log.Info("validatedPromoAdditionalCriteriaByProductGroups["+promoID+"]["+strconv.Itoa(proID)+"].SubTotal :", validatedPromoAdditionalCriteriaByProductGroups[promoID][proID].SubTotal)
					rewardPrice := entity.ConsultPromotionRewardPriceResponse{}
					rewardPrice.ProID = proID
					rewardPrice.SubTotal = float64(validatedPromoAdditionalCriteriaByProductGroups[promoID][proID].SubTotal)

					reward := math.Round((rewardPrice.SubTotal * response.SlabReward) / float64(subTotalValidatedPromoAdditionalCriteriaByProductGroups[promoID]))
					slabReward -= reward
					if slabReward <= 0 {
						reward += slabReward
					}
					// log.Info("slabReward "+promoID+" :", slabReward)
					rewardPrice.Reward = reward
					rewardPrice.Total = rewardPrice.SubTotal - rewardPrice.Reward

					response.RewardPrice = append(response.RewardPrice, rewardPrice)
				}

				// log.Info("SlabRule", validatedPromoCriterias[promoID].SlabRule)
				// log.Info("SlabRuleTo", validatedPromoCriterias[promoID].SlabRuleTo)
				log.Info("PENGHITUNGAN REWARD PRODUCT")
				if response.SlabReward == 0 {
					rewards, _ := service.PromotionRepository.GetAllRewardProductFromStock(request, validatedPromoCriterias[promoID])

					multipliedValue := int64(1)
					if validatedPromoCriterias[promoID].IsMultiplied {
						multipliedValue = validatedPromoCriterias[promoID].SlabRule / int64(validatedPromoCriterias[promoID].SlabRuleTo)
					}
					totalQtyReward := validatedPromoCriterias[promoID].SlabReward * float64(multipliedValue)

					for _, reward := range rewards {
						var rewardProduct entity.ConsultPromotionRewardProductResponse

						qtyReward := totalQtyReward
						if totalQtyReward >= float64(reward.QtyStock) {
							qtyReward = float64(reward.QtyStock)
						}
						totalQtyReward -= qtyReward

						rewardProduct.ProID = int(reward.ProID)
						rewardProduct.Qty = qtyReward
						rewardProduct.UnitId = reward.UnitId
						rewardProduct.Uom = validatedPromoCriterias[promoID].SlabRewardUom

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


	return responses, nil
}
*/

func (service *discountServiceImpl) DetailGrp(DiscGrpId string) (data []entity.DetailDiscountGrp, err error) {
	discounts, err := service.DiscountRepository.FindDiscGrpId(DiscGrpId)
	if err != nil {
		return data, err
	}

	for _, row := range discounts {
		var vResp entity.DetailDiscountGrp
		structs.Automapper(row, &vResp)

		// vResp.CustID = ""
		// vResp.EffectiveFrom = row.EffectiveFrom.Format(constant.YYYY_MM_DD)
		// vResp.EffectiveTo = row.EffectiveTo.Format(constant.YYYY_MM_DD)
		// vResp.DiscountStatusDesc = vResp.GetDiscountStatusDesc()
		// vResp.PublishStatusDesc = vResp.GetPublishStatusDesc()

		data = append(data, vResp)
	}

	return data, err
}
