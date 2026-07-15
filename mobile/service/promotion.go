package service

import (
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/config/env"
	"mobile/repository"
	"strconv"
	"strings"

	// "context"
	// "errors"
	"math"

	// "mobile/pkg/constant"
	// "mobile/pkg/str"
	// "mobile/pkg/structs"

	"github.com/gofiber/fiber/v2/log"
)

type PromotionService interface {
	List(dataFilter entity.PromotionsQueryFilter, custId, parentCustId string) (data []entity.PromotionsResp, total int64, lastPage int, err error)
	ListMobile(dataFilter entity.PromotionMobileListQueryFilter, custId, parentCustId string) (data []entity.PromotionMobileListResponse, total int64, lastPage int, err error)
	DetailMobile(promoId string, custId, parentCustId string) (data entity.PromotionMobileDetailResponse, err error)
	ConsultPromotion(equest entity.ConsultPromotionBody) (responses []entity.ConsultPromotionResponse, err error)
	Conversion(conversionBody entity.CreateConversionBody, custID string, parentCustID string) (response entity.OrderConversionResponse, err error)
	OutletList(dataFilter entity.PromotionOutletListQueryFilter, otTypeID int64, custId, parentCustId string) (data []entity.PromotionOutletListResponse, total int64, lastPage int, err error)
}

type PromotionServiceImpl struct {
	Config              env.ConfigEnv
	PromotionRepository repository.PromotionRepository
	// MCustomerRepository repository.MCustomerRepository,
	Transaction repository.Dbtransaction
}

func NewPromotionService(
	promotionRepository repository.PromotionRepository,
	config env.ConfigEnv,
	transaction repository.Dbtransaction,
) *PromotionServiceImpl {
	return &PromotionServiceImpl{
		PromotionRepository: promotionRepository,
		Config:              config,
		Transaction:         transaction,
	}
}

func (service *PromotionServiceImpl) List(dataFilter entity.PromotionsQueryFilter, custId, parentCustId string) (response []entity.PromotionsResp, total int64, lastPage int, err error) {
	var (
		promo  entity.PromotionsResp
		iToStr string
	)
	for i := 1; i < 3; i++ {
		iToStr = strconv.Itoa(i)
		promo.Code = "0000" + iToStr
		promo.Title = iToStr + " Gratis 1 pcs produk A"

		promo.Tnc = []string{
			"Valid only during specified promo",
			"Applicate to eligle products",
			"Discount limited to one use per cust",
		}
		response = append(response, promo)
	}

	return response, total, lastPage, err
}

// ListMobile returns active promotions for mobile app
func (service *PromotionServiceImpl) ListMobile(dataFilter entity.PromotionMobileListQueryFilter, custId, parentCustId string) (response []entity.PromotionMobileListResponse, total int64, lastPage int, err error) {
	dataFilter.CustId = custId
	dataFilter.ParentCustId = parentCustId

	// Get active promotions
	promotions, total, lastPage, err := service.PromotionRepository.FindAllActivePromotions(dataFilter)
	if err != nil {
		return response, total, lastPage, err
	}

	if len(promotions) == 0 {
		return response, total, lastPage, nil
	}

	for _, promo := range promotions {
		var promoResp entity.PromotionMobileListResponse
		promoResp.PromoID = promo.PromoID
		promoResp.PromoDesc = promo.PromoDesc

		// Format dates
		if promo.EffectiveFrom != nil {
			promoResp.EffectiveFrom = promo.EffectiveFrom.Format("02-01-2006")
		}
		if promo.EffectiveTo != nil {
			promoResp.EffectiveTo = promo.EffectiveTo.Format("02-01-2006")
		}

		promoResp.IsMultiplied = promo.IsMultiplied
		promoResp.MaxInvoiceOutlet = promo.MaxInvoiceOutlet

		response = append(response, promoResp)
	}

	return response, total, lastPage, nil
}

// DetailMobile returns promotion detail for mobile app
func (service *PromotionServiceImpl) DetailMobile(promoId string, custId, parentCustId string) (response entity.PromotionMobileDetailResponse, err error) {
	// Initialize AdditionalCriteria with empty slices to avoid null in JSON
	response.AdditionalCriteria = entity.PromotionMobileDetailAdditionalCriteria{
		SelectedOutletTypes:   []entity.PromotionMobileOutletType{},
		SelectedOutletGroups:  []entity.PromotionMobileOutletGroup{},
		SelectedOutletClasses: []entity.PromotionMobileOutletClass{},
	}

	// Get promotion detail
	promotion, err := service.PromotionRepository.FindPromotionDetailByPromoID(promoId, parentCustId)
	if err != nil {
		return response, err
	}

	response.PromoID = promotion.PromoID
	response.PromoDesc = promotion.PromoDesc
	response.PromoType = promotion.PromoType

	// Format dates
	if promotion.EffectiveFrom != nil {
		response.EffectiveFrom = promotion.EffectiveFrom.Format("02-01-2006")
	}
	if promotion.EffectiveTo != nil {
		response.EffectiveTo = promotion.EffectiveTo.Format("02-01-2006")
	}

	response.IsMultiplied = promotion.IsMultiplied
	response.MaxInvoiceOutlet = promotion.MaxInvoiceOutlet
	response.MaxPromoUsage = promotion.MaxPromoUsage
	response.MaxTotalRewardType = promotion.MaxTotalRewardType
	response.MaxTotalRewardValue = promotion.MaxTotalRewardValue

	// Get promoted products
	promotedProducts, err := service.PromotionRepository.FindPromotedProductsByPromoID(promoId, parentCustId)
	if err != nil {
		log.Error("DetailMobile, FindPromotedProductsByPromoID, err:", err.Error())
	} else {
		for _, pp := range promotedProducts {
			response.PromotedProduct = append(response.PromotedProduct, entity.PromotionMobileDetailPromotedProduct{
				ProID:       pp.ProID,
				ProCode:     pp.ProCode,
				ProName:     pp.ProName,
				Mandatory:   pp.Mandatory,
				MinBuyType:  pp.MinBuyType,
				MinBuyQty:   pp.MinBuyQty,
				MinBuyValue: pp.MinBuyValue,
				MinBuyUom:   pp.MinBuyUom,
			})
		}
	}

	// Get criteria
	criteria, err := service.PromotionRepository.FindPromotionCriteriaByPromoID(promoId)
	if err != nil {
		log.Error("DetailMobile, FindPromotionCriteriaByPromoID, err:", err.Error())
	} else {
		for _, c := range criteria {
			response.Criteria = append(response.Criteria, entity.PromotionMobileDetailCriteria{
				ProID:       c.ProID,
				ProCode:     c.ProCode,
				ProName:     c.ProName,
				CountPromo:  c.CountPromo,
				MinPurchase: c.MinPurchase,
				MaxPurchase: c.MaxPurchase,
				Uom:         c.Uom,
			})
		}
	}

	// Get promotion rewards (slabs or stratas)
	promotionRewards, err := service.PromotionRepository.FindPromotionRewardByPromoID(promoId, parentCustId)
	if err != nil {
		log.Error("DetailMobile, FindPromotionRewardByPromoID, err:", err.Error())
	}

	// Build promotion_reward array based on promo_type
	if strings.ToLower(promotion.PromoType) == "slab" {
		slabs, err := service.PromotionRepository.FindPromotionSlabByPromoID(promoId)
		if err != nil {
			log.Error("DetailMobile, FindPromotionSlabByPromoID, err:", err.Error())
		} else {
			for _, slab := range slabs {
				rewardItem := entity.PromotionMobileDetailRewardItem{
					SlabID:      &slab.SlabID,
					SlabName:    &slab.SlabName,
					Ordinal:     slab.Ordinal,
					RuleType:    slab.RuleType,
					RangeFrom:   slab.RangeFrom,
					RangeTo:     slab.RangeTo,
					RuleUom:     slab.RuleUom,
					RewardType:  slab.RewardType,
					RewardUom:   slab.RewardUom,
					RewardValue: slab.RewardValue,
				}

				// Add product_reward if reward_type is "product"
				if strings.ToLower(slab.RewardType) == "product" {
					for _, reward := range promotionRewards {
						rewardItem.ProductReward = append(rewardItem.ProductReward, entity.PromotionMobileDetailProductReward{
							RewardProID: reward.RewardProID,
							ProID:       reward.ProID,
							ProCode:     reward.ProCode,
							ProName:     reward.ProName,
						})
					}
				}

				response.PromotionReward = append(response.PromotionReward, rewardItem)
			}
		}
	} else if strings.ToLower(promotion.PromoType) == "strata" {
		stratas, err := service.PromotionRepository.FindPromotionStrataByPromoID(promoId)
		if err != nil {
			log.Error("DetailMobile, FindPromotionStrataByPromoID, err:", err.Error())
		} else {
			for _, strata := range stratas {
				rewardItem := entity.PromotionMobileDetailRewardItem{
					StrataID:    &strata.StrataID,
					StrataName:  &strata.StrataName,
					Ordinal:     strata.Ordinal,
					RuleType:    strata.RuleType,
					RangeFrom:   strata.RangeFrom,
					RangeTo:     strata.RangeTo,
					RuleUom:     strata.RuleUom,
					RewardType:  strata.RewardType,
					RewardUom:   strata.RewardUom,
					RewardValue: strata.RewardValue,
				}

				// Add product_reward if reward_type is "product"
				if strings.ToLower(strata.RewardType) == "product" {
					for _, reward := range promotionRewards {
						rewardItem.ProductReward = append(rewardItem.ProductReward, entity.PromotionMobileDetailProductReward{
							RewardProID: reward.RewardProID,
							ProID:       reward.ProID,
							ProCode:     reward.ProCode,
							ProName:     reward.ProName,
						})
					}
				}

				response.PromotionReward = append(response.PromotionReward, rewardItem)
			}
		}
	}

	// Get additional criteria
	outletTypes, err := service.PromotionRepository.FindOutletTypesByPromoID(promoId, parentCustId)
	if err != nil {
		log.Error("DetailMobile, FindOutletTypesByPromoID, err:", err.Error())
	} else {
		for _, ot := range outletTypes {
			response.AdditionalCriteria.SelectedOutletTypes = append(response.AdditionalCriteria.SelectedOutletTypes, entity.PromotionMobileOutletType{
				OutletTypeID:   ot.OutletTypeID,
				OutletTypeCode: ot.OutletTypeCode,
				OutletTypeName: ot.OutletTypeName,
			})
		}
	}

	outletGroups, err := service.PromotionRepository.FindOutletGroupsByPromoID(promoId, parentCustId)
	if err != nil {
		log.Error("DetailMobile, FindOutletGroupsByPromoID, err:", err.Error())
	} else {
		for _, og := range outletGroups {
			response.AdditionalCriteria.SelectedOutletGroups = append(response.AdditionalCriteria.SelectedOutletGroups, entity.PromotionMobileOutletGroup{
				OutletGroupID:   og.OutletGroupID,
				OutletGroupCode: og.OutletGroupCode,
				OutletGroupName: og.OutletGroupName,
			})
		}
	}

	outletClasses, err := service.PromotionRepository.FindOutletClassesByPromoID(promoId, parentCustId)
	if err != nil {
		log.Error("DetailMobile, FindOutletClassesByPromoID, err:", err.Error())
	} else {
		for _, oc := range outletClasses {
			response.AdditionalCriteria.SelectedOutletClasses = append(response.AdditionalCriteria.SelectedOutletClasses, entity.PromotionMobileOutletClass{
				OutletClassID:   oc.OutletClassID,
				OutletClassCode: oc.OutletClassCode,
				OutletClassName: oc.OutletClassName,
			})
		}
	}

	return response, nil
}

func (service *PromotionServiceImpl) ConsultPromotion(request entity.ConsultPromotionBody) (responses []entity.ConsultPromotionResponse, err error) {
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
		validatedPromoAdditionalCriteriaByProductGroups := make(map[string]map[int]*entity.ConsultPromotionSubBody)
		subTotalValidatedPromoAdditionalCriteriaByProductGroups := make(map[string]int64)
		var validatedPromoList []string
		for promoID := range validatedPromoAdditionalCriteriaGroups {
			validatedPromoAdditionalCriteriaByProductGroups[promoID] = make(map[int]*entity.ConsultPromotionSubBody)
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
					log.Info("validatedPromoAdditionalCriteriaByProductGroups["+promoID+"]["+strconv.Itoa(req.ProID)+"] :", int64(req.SubTotal))
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
				}
				// log.Info("slabRuleValue SLAB ("+strconv.FormatInt(*promoCriterias[index].SlabID, 10)+") :", slabRuleValue)

				if promoCriterias[index].IsMultiplied || (slabRuleValue >= promoCriterias[index].SlabRuleFrom && slabRuleValue <= promoCriterias[index].SlabRuleTo) {
					validatedPromoCriterias[promoCriterias[index].PromoID] = promoCriterias[index]
					isPromoCriteriaValid = true
				}
			}
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

					qty := int64(0)
					var rewardProductConversion entity.CreateConversionBody
					rewardProductConversion.CustId = request.CustID
					rewardProductConversion.ProductId = reward.ProID

					switch validatedPromoCriterias[promoID].SlabRewardUom {
					case 1:
						rewardProductConversion.Qty1 = int64(qtyReward)
						rewardProductConversion.Qty2 = qty
						rewardProductConversion.Qty3 = qty
					case 2:
						rewardProductConversion.Qty1 = qty
						rewardProductConversion.Qty2 = int64(qtyReward)
						rewardProductConversion.Qty3 = qty
					default:
						rewardProductConversion.Qty1 = qty
						rewardProductConversion.Qty2 = qty
						rewardProductConversion.Qty3 = int64(qtyReward)
					}
					rewardProductConversionResut, _ := service.Conversion(rewardProductConversion, request.CustID, request.ParentCustID)

					rewardProduct.ProID = int(reward.ProID)
					rewardProduct.Qty1 = float64(rewardProductConversionResut.Qty1)
					rewardProduct.Qty2 = float64(rewardProductConversionResut.Qty2)
					rewardProduct.Qty3 = float64(rewardProductConversionResut.Qty3)
					// rewardProduct.UnitId = reward.UnitId
					// rewardProduct.Uom = validatedPromoCriterias[promoID].SlabRewardUom

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

	return responses, nil
}

func (service *PromotionServiceImpl) Conversion(conversionBody entity.CreateConversionBody, custID string, parentCustID string) (response entity.OrderConversionResponse, err error) {
	product, err := service.PromotionRepository.FindProductByID(conversionBody.ProductId)
	if err != nil {
		return response, err
	}

	qty1 := conversionBody.Qty1
	qty2 := conversionBody.Qty2
	qty3 := conversionBody.Qty3

	rQty2 := qty1 / int64(product.ConvUnit2)
	if rQty2 > 0 {
		qty1 = qty1 % int64(product.ConvUnit2)
		qty2 += rQty2
	}

	rQty3 := qty2 / int64(product.ConvUnit3)
	if rQty3 > 0 {
		qty2 = qty2 % int64(product.ConvUnit3)
		qty3 += rQty3
	}

	response.Qty1 = qty1
	response.Qty2 = qty2
	response.Qty3 = qty3

	response.TotalQty = (int64(product.ConvUnit2)*int64(product.ConvUnit3))*qty3 + (int64(product.ConvUnit2) * qty2) + qty1

	return response, err
}

func (service *PromotionServiceImpl) OutletList(dataFilter entity.PromotionOutletListQueryFilter, otTypeID int64, custId, parentCustId string) (response []entity.PromotionOutletListResponse, total int64, lastPage int, err error) {
	outlets, total, lastPage, err := service.PromotionRepository.FindOutletsByTypeGroupClass(dataFilter, otTypeID, custId, parentCustId)
	if err != nil {
		return response, total, lastPage, err
	}

	response = make([]entity.PromotionOutletListResponse, 0, len(outlets))
	for _, outlet := range outlets {
		todayVisit := false
		if outlet.TodayVisit != nil {
			todayVisit = *outlet.TodayVisit
		}

		response = append(response, entity.PromotionOutletListResponse{
			OtTypeID:   outlet.OtTypeID,
			OutletCode: outlet.OutletCode,
			OutletName: outlet.OutletName,
			Address1:   outlet.Address1,
			TodayVisit: todayVisit,
		})
	}

	return response, total, lastPage, nil
}
