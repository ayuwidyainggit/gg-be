package service

import (
	"context"
	"sales/entity"
	"sales/model"
	"sales/pkg/constant"
	"sales/pkg/structs"
	"sales/repository"

	"github.com/gofiber/fiber/v2/log"
)

type PromoTemplateService interface {
	Store(request entity.CreatePromoTemplateBody) (err error)
	Detail(params entity.DetailPromoTemplateParams) (response entity.PromoTemplate, err error)
	List(dataFilter entity.PromoTemplateQueryFilter) (data []entity.PromoTemplate, total int64, lastPage int, err error)
	Update(promoID string, request entity.UpdatePromoTemplateBody) (err error)
	Delete(custId, promoID, deletedBy string) (err error)
}

func NewPromoTemplateService(promoTemplateRepository repository.PromoTemplateRepository, transaction repository.Dbtransaction) *promoTemplateServiceImpl {
	return &promoTemplateServiceImpl{
		PromoTemplateRepository: promoTemplateRepository,
		Transaction:             transaction,
	}
}

type promoTemplateServiceImpl struct {
	PromoTemplateRepository repository.PromoTemplateRepository
	Transaction             repository.Dbtransaction
}

func (service *promoTemplateServiceImpl) Store(request entity.CreatePromoTemplateBody) (err error) {
	c := context.Background()

	var promoModel model.PromoTemplate
	err = structs.Automapper(request, &promoModel)
	if err != nil {
		return err
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		promoTemplateID, err := service.PromoTemplateRepository.Store(txCtx, &promoModel)
		if err != nil {
			log.Error("err:", err.Error())
			return err
		}

		isHaveRewardProducts := false
		for _, row := range request.PromoTemplateCriteria {
			var promoCriteriaModel model.PromoTemplateCriteria
			err := structs.Automapper(row, &promoCriteriaModel)
			if err != nil {
				return err
			}
			promoCriteriaModel.CustID = request.CustID
			promoCriteriaModel.PromoTemplateID = promoTemplateID
			// log.Info("promoCriteriaModel:", structs.StructToJson(promoCriteriaModel))
			err = service.PromoTemplateRepository.StorePromoCriteria(txCtx, &promoCriteriaModel)
			if err != nil {
				return err
			}

			if row.SlabRewardType == 1 {
				isHaveRewardProducts = true
			}
		}

		if isHaveRewardProducts {
			for _, row := range request.PromoTemplateRewardProduct {
				var rewardProductModel model.PromoTemplateRewardProduct
				err := structs.Automapper(row, &rewardProductModel)
				if err != nil {
					return err
				}
				rewardProductModel.CustID = request.CustID
				rewardProductModel.PromoTemplateID = promoTemplateID
				err = service.PromoTemplateRepository.StorePromoRewardProduct(txCtx, &rewardProductModel)
				if err != nil {
					return err
				}
			}
		}

		for _, row := range request.PromoTemplateAdditionalCriteria {
			var promoAddCriteriaModel model.PromoTemplateAdditionalCriteria
			err := structs.Automapper(row, &promoAddCriteriaModel)
			if err != nil {
				return err
			}
			promoAddCriteriaModel.CustID = request.CustID
			promoAddCriteriaModel.PromoTemplateID = promoTemplateID
			err = service.PromoTemplateRepository.StorePromoAdditionalCriteria(txCtx, &promoAddCriteriaModel)
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

func (service *promoTemplateServiceImpl) Detail(params entity.DetailPromoTemplateParams) (response entity.PromoTemplate, err error) {
	promo, err := service.PromoTemplateRepository.FindByPromoTemplateID(params)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(promo, &response)
	if err != nil {
		return response, err
	}

	response.BudgetReferenceTypeName = response.GetPromoBudgetReferenceTypeName()
	response.MaxDiscountTypeName = constant.GetQtyAmountPercentDisplayName(response.MaxDiscountType)
	response.PromoStatusDesc = response.GetPromoStatusDesc()
	response.MaxDiscountOutletUomName = constant.GetUomName(response.MaxDiscountOutletUom)

	err = service.GetPromoCriterias(params, &response)
	if err != nil {
		return response, err
	}

	err = service.GetRewardProducts(params, &response)
	if err != nil {
		return response, err
	}

	err = service.GetPromoAdditionalCriterias(params, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (service *promoTemplateServiceImpl) List(dataFilter entity.PromoTemplateQueryFilter) (data []entity.PromoTemplate, total int64, lastPage int, err error) {
	promotions, total, lastPage, err := service.PromoTemplateRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range promotions {
		var vResp entity.PromoTemplate
		structs.Automapper(row, &vResp)

		// vResp.PromoTypeName = vResp.GetPromoTypeName()
		vResp.BudgetReferenceTypeName = vResp.GetPromoBudgetReferenceTypeName()
		vResp.MaxDiscountTypeName = constant.GetQtyAmountPercentDisplayName(vResp.MaxDiscountType)
		vResp.PromoStatusDesc = vResp.GetPromoStatusDesc()
		vResp.MaxDiscountOutletUomName = constant.GetUomName(vResp.MaxDiscountOutletUom)

		// payTypeName := vResp.GeneratePayTypeName()
		// vResp.PayTypeName = payTypeName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *promoTemplateServiceImpl) Update(promoID string, request entity.UpdatePromoTemplateBody) (err error) {
	c := context.Background()

	var promoModel model.PromoTemplate
	err = structs.Automapper(request, &promoModel)
	if err != nil {
		return err
	}

	promoModel.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.PromoTemplateRepository.Update(txCtx, promoID, promoModel)
		if err != nil {
			return err
		}

		isHaveRewardProducts := false

		err := service.PromoTemplateRepository.DeletePromoCriterias(txCtx, request.CustID, promoID)
		if err != nil {
			log.Error("DeletePromoCriteriasNotInIDs, error:", err.Error())
		}

		err = service.PromoTemplateRepository.DeletePromoAdditionalCriterias(txCtx, request.CustID, promoID)
		if err != nil {
			log.Error("DeletePromoCriteriasNotInIDs, error:", err.Error())
		}

		for _, row := range request.PromoTemplateCriteria {

			var promoCritModel model.PromoTemplateCriteria
			err = structs.Automapper(row, &promoCritModel)
			if err != nil {
				return err
			}
			promoCritModel.CustID = request.CustID
			promoCritModel.PromoTemplateID = promoID
			promoCritModel.PromoTemplateSlabID = nil
			err = service.PromoTemplateRepository.StorePromoCriteria(txCtx, &promoCritModel)
			if err != nil {
				return err
			}

			if row.SlabRewardType == 1 {
				isHaveRewardProducts = true
			}
		}

		err = service.PromoTemplateRepository.DeletePromoTemplateRewardProducts(txCtx, request.CustID, promoID)
		if err != nil {
			log.Error("DeletePromoTemplateRewardProducts, error:", err.Error())
		}
		if isHaveRewardProducts {
			for _, row := range request.PromoTemplateRewardProduct {
				var rewardProductModel model.PromoTemplateRewardProduct
				err := structs.Automapper(row, &rewardProductModel)
				if err != nil {
					return err
				}
				rewardProductModel.CustID = request.CustID
				rewardProductModel.PromoTemplateID = promoID
				err = service.PromoTemplateRepository.StorePromoRewardProduct(txCtx, &rewardProductModel)
				if err != nil {
					return err
				}
			}
		}

		for _, row := range request.PromoTemplateAdditionalCriteria {
			var promoCritAddModel model.PromoTemplateAdditionalCriteria
			err = structs.Automapper(row, &promoCritAddModel)
			if err != nil {
				return err
			}
			promoCritAddModel.CustID = request.CustID
			promoCritAddModel.PromoTemplateID = promoID
			promoCritAddModel.PromoTeamplateAddCriteriaID = nil
			err = service.PromoTemplateRepository.StorePromoAdditionalCriteria(txCtx, &promoCritAddModel)
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

func (service *promoTemplateServiceImpl) Delete(custId, promoID, deletedBy string) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.PromoTemplateRepository.Delete(txCtx, custId, promoID)
		if err != nil {
			return err
		}

		err = service.PromoTemplateRepository.DeletePromoCriterias(txCtx, custId, promoID)
		if err != nil {
			log.Error("DeletePromoCriterias, error:", err.Error())
		}

		err = service.PromoTemplateRepository.DeletePromoAdditionalCriterias(txCtx, custId, promoID)
		if err != nil {
			log.Error("DeletePromoAdditionalCriterias, error:", err.Error())
		}

		err = service.PromoTemplateRepository.DeletePromoTemplateRewardProducts(txCtx, custId, promoID)
		if err != nil {
			log.Error("DeletePromoTemplateRewardProducts, error:", err.Error())
		}

		return nil
	})

	return err
}

func (service *promoTemplateServiceImpl) GetPromoCriterias(params entity.DetailPromoTemplateParams, promoResponse *entity.PromoTemplate) (err error) {
	promoCriterias, err := service.PromoTemplateRepository.FindAllPromoCriteriasByPromoTemplateID(params)
	if err != nil {
		return err
	}

	for _, row := range promoCriterias {
		var promoCriteria entity.PromoTemplateCriteria
		err = structs.Automapper(row, &promoCriteria)
		if err != nil {
			return err
		}
		promoCriteria.CustID = ""
		promoCriteria.PromoTemplateID = ""
		promoCriteria.SlabRuleTypeName = constant.GetQtyAmountPercentDisplayName(int(promoCriteria.SlabRuleType))
		promoCriteria.SlabRewardTypeName = constant.GetQtyAmountPercentDisplayName(int(promoCriteria.SlabRewardType))
		promoCriteria.SlabRuleUomName = constant.GetUomName(promoCriteria.SlabRuleUom)
		promoCriteria.SlabRewardUomName = constant.GetUomName(int(promoCriteria.SlabRewardUom))

		promoResponse.PromoCriterias = append(promoResponse.PromoCriterias, promoCriteria)
	}

	return
}

func (service *promoTemplateServiceImpl) GetPromoAdditionalCriterias(params entity.DetailPromoTemplateParams, promoResponse *entity.PromoTemplate) (err error) {
	promoAdditionalCriterias, err := service.PromoTemplateRepository.FindAllPromoAdditionalCriteriasByPromoTemplateID(params)
	if err != nil {
		return err
	}

	for _, row := range promoAdditionalCriterias {
		var promoAddCriteria entity.PromoTemplateAdditionalCriteria
		err = structs.Automapper(row, &promoAddCriteria)
		if err != nil {
			return err
		}
		promoAddCriteria.CustID = ""
		promoAddCriteria.PromoTemplateID = ""
		promoAddCriteria.AttributeName = constant.GetPromoAttributeDisplayName(promoAddCriteria.Attribute)
		promoAddCriteria.ConditionName = constant.GetIncludeExcludeDisplayName(promoAddCriteria.Condition)
		promoAddCriteria.MinBuyTypeName = constant.GetQtyAmountPercentDisplayName(promoAddCriteria.MinBuyType)
		promoAddCriteria.MinBuyUomName = constant.GetUomName(promoAddCriteria.MinBuyUom)

		// Switch by Attribute
		switch promoAddCriteria.Attribute {

		case constant.AttrProduct:
			product, err := service.PromoTemplateRepository.FindOneProductByProID(params.CustID, params.ParentCustId, promoAddCriteria.ReferenceID)
			if err == nil {
				promoAddCriteria.ReferenceCode = product.ReferenceCode
				promoAddCriteria.ReferenceName = product.ReferenceName
			}
		case constant.AttrOutletClass:
			outletClass, err := service.PromoTemplateRepository.FindOneOutletClassByID(params.CustID, params.ParentCustId, promoAddCriteria.ReferenceID)
			if err == nil {
				promoAddCriteria.ReferenceCode = outletClass.ReferenceCode
				promoAddCriteria.ReferenceName = outletClass.ReferenceName
			}
		case constant.AttrOutletType:
			outletType, err := service.PromoTemplateRepository.FindOneOutletTypeByID(params.CustID, params.ParentCustId, promoAddCriteria.ReferenceID)
			if err == nil {
				promoAddCriteria.ReferenceCode = outletType.ReferenceCode
				promoAddCriteria.ReferenceName = outletType.ReferenceName
			}
		case constant.AttrOutletGroup:
			outletGroup, err := service.PromoTemplateRepository.FindOneOutletGroupByID(params.CustID, params.ParentCustId, promoAddCriteria.ReferenceID)
			if err == nil {
				promoAddCriteria.ReferenceCode = outletGroup.ReferenceCode
				promoAddCriteria.ReferenceName = outletGroup.ReferenceName
			}
		case constant.AttrSalesType:
			salesType, err := service.PromoTemplateRepository.FindOneSalesTypeByID(params.CustID, params.ParentCustId, promoAddCriteria.ReferenceID)
			if err == nil {
				promoAddCriteria.ReferenceCode = salesType.ReferenceCode
				promoAddCriteria.ReferenceName = salesType.ReferenceName
			}
		case constant.AttrSalesTeam:
			salesTeam, err := service.PromoTemplateRepository.FindOneSalesTeamByID(params.CustID, params.ParentCustId, promoAddCriteria.ReferenceID)
			if err == nil {
				promoAddCriteria.ReferenceCode = salesTeam.ReferenceCode
				promoAddCriteria.ReferenceName = salesTeam.ReferenceName
			}
		default:

		}

		promoResponse.PromoAdditionalCriterias = append(promoResponse.PromoAdditionalCriterias, promoAddCriteria)

	}

	return
}

func (service *promoTemplateServiceImpl) GetRewardProducts(params entity.DetailPromoTemplateParams, promoResponse *entity.PromoTemplate) (err error) {
	rewardProducts, err := service.PromoTemplateRepository.FindAllRewardProductsByPromoTemplateID(params)
	if err != nil {
		return err
	}
	log.Info("rewardProducts:", structs.StructToJson(rewardProducts))
	for _, row := range rewardProducts {
		var rewardProduct entity.PromoTemplateRewardProduct
		err = structs.Automapper(row, &rewardProduct)
		if err != nil {
			return err
		}
		rewardProduct.CustID = ""
		rewardProduct.PromoTemplateID = ""
		productDetail, err := service.PromoTemplateRepository.FindProductByProID(params.ParentCustId, row.ProID)
		if err == nil {
			rewardProduct.ProCode = productDetail.ProCode
			rewardProduct.ProName = productDetail.ProName
		}

		promoResponse.RewardProduct = append(promoResponse.RewardProduct, rewardProduct)
	}

	return
}
