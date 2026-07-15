package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sales/entity"
	"sales/model"
	"sales/pkg/constant"
	"sales/pkg/errmsg"
	"sales/pkg/str"
	"sales/pkg/structs"
	"sales/repository"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type PromotionService interface {
	Store(request entity.CreatePromotionBody) (err error)
	Detail(params entity.DetailPromotionParams) (response entity.Promotion, err error)
	List(dataFilter entity.PromotionQueryFilter) (data []entity.Promotion, total int64, lastPage int, err error)
	Update(promoID string, request entity.UpdatePromotionBody) (err error)
	Delete(params entity.DetailPromotionParams, deletedBy string) (err error)
	BulkUpdateStatus(equest entity.BulkUpdateStatusPromotionBody) (err error)
	ConsultPromotion(equest entity.ConsultPromotionBody) (responses []entity.ConsultPromotionResponse, err error)
	Conversion(conversionBody entity.CreateConversionBody, custID string, parentCustID string) (response entity.OrderConversionResponse, err error)
	StoreV2(request entity.CreatePromotionV2Body) (err error)
	DetailV2(params entity.DetailPromotionParams) (response entity.PromotionV2, err error)
	ListV2(dataFilter entity.PromotionV2QueryFilter) (data []entity.PromotionV2, total int64, lastPage int, err error)
	UpdateV2(promoID string, request entity.UpdatePromotionV2Body) (err error)
	UpdateV2Status(promoID string, promoStatus entity.PromotionV2Status, req entity.UpdateStatusPromotionV2Body) (err error)
	DuplicateV2(params entity.DetailPromotionParams) (newPromoID string, err error)
	ExistsPromo(custID, promoID string) (bool, error)
	DetailV2ForUpdate(params entity.DetailPromotionParams) (response entity.PromotionV2, err error)
	ConsultV2(req entity.ConsultPromoV2Req) (resp []entity.ConsultPromoResp, err error)
	PromoConversion(conversionBody entity.PromoConversionReq, custID string) (response entity.PromoConversionResp, err error)
	CloseExpiredPromotions() (err error)
}

func NewPromotionService(promotionRepository repository.PromotionRepository, promotionV2Repository repository.PromotionV2Repository, transaction repository.Dbtransaction) *promotionServiceImpl {
	return &promotionServiceImpl{
		PromotionRepository:   promotionRepository,
		PromotionV2Repository: promotionV2Repository,
		Transaction:           transaction,
	}
}

type promotionServiceImpl struct {
	PromotionRepository   repository.PromotionRepository
	PromotionV2Repository repository.PromotionV2Repository
	Transaction           repository.Dbtransaction
}

func (service *promotionServiceImpl) Store(request entity.CreatePromotionBody) (err error) {
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

	request.PromoStatusID = 1 // make it default status 'Draft'
	var promoModel model.Promotion
	err = structs.Automapper(request, &promoModel)
	if err != nil {
		return err
	}

	var promoStatusLogModel model.PromoStatusLog
	err = structs.Automapper(request, &promoStatusLogModel)
	if err != nil {
		return err
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.PromotionRepository.Store(txCtx, &promoModel)
		if err != nil {
			log.Error("err:", err.Error())
			return err
		}

		err = service.PromotionRepository.StoreStatusLog(txCtx, &promoStatusLogModel)
		if err != nil {
			log.Error("err:", err.Error())
			return err
		}

		isHaveRewardProducts := false
		for _, row := range request.PromoCriteria {
			var promoCriteriaModel model.PromoCriteria
			err := structs.Automapper(row, &promoCriteriaModel)
			if err != nil {
				return err
			}
			promoCriteriaModel.CustID = request.CustID
			promoCriteriaModel.PromoID = request.PromoID
			// log.Info("promoCriteriaModel:", structs.StructToJson(promoCriteriaModel))
			err = service.PromotionRepository.StorePromoCriteria(txCtx, &promoCriteriaModel)
			if err != nil {
				return err
			}

			if row.SlabRewardType == entity.PromoRewardTypeQuantity {
				isHaveRewardProducts = true
			}
		}

		if isHaveRewardProducts {
			if len(request.RewardProduct) < 1 {
				return errors.New("minimum have 1 reward product")
			}

			for _, row := range request.RewardProduct {
				var rewardProductModel model.PromoRewardProduct
				err := structs.Automapper(row, &rewardProductModel)
				if err != nil {
					return err
				}
				rewardProductModel.CustID = request.CustID
				rewardProductModel.PromoID = request.PromoID
				err = service.PromotionRepository.StorePromoRewardProduct(txCtx, &rewardProductModel)
				if err != nil {
					return err
				}
			}
		}

		for _, row := range request.PromoAdditionalCriteria {
			var promoAddCriteriaModel model.PromoAdditionalCriteria
			err := structs.Automapper(row, &promoAddCriteriaModel)
			if err != nil {
				return err
			}
			promoAddCriteriaModel.CustID = request.CustID
			promoAddCriteriaModel.PromoID = request.PromoID
			err = service.PromotionRepository.StorePromoAdditionalCriteria(txCtx, &promoAddCriteriaModel)
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

func (service *promotionServiceImpl) Detail(params entity.DetailPromotionParams) (response entity.Promotion, err error) {
	promo, err := service.PromotionRepository.FindByPromoID(params)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(promo, &response)
	if err != nil {
		return response, err
	}

	response.EffectiveFrom = promo.EffectiveFrom.Format(constant.YYYY_MM_DD)
	response.EffectiveTo = promo.EffectiveTo.Format(constant.YYYY_MM_DD)
	response.PromoTypeName = response.GetPromoTypeName()
	response.BudgetReferenceTypeName = response.GetPromoBudgetReferenceTypeName()
	response.MaxDiscountTypeName = constant.GetQtyAmountPercentDisplayName(response.MaxDiscountType)
	response.PromoStatusDesc = response.GetPromoStatusDesc()
	response.BudgetControlLevelName = constant.GetPromoScopeLevelName(response.BudgetControlLevel)
	response.ExecutionLevelName = constant.GetPromoScopeLevelName(response.ExecutionLevel)
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

func (service *promotionServiceImpl) List(dataFilter entity.PromotionQueryFilter) (data []entity.Promotion, total int64, lastPage int, err error) {
	promotions, total, lastPage, err := service.PromotionRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range promotions {
		var vResp entity.Promotion
		structs.Automapper(row, &vResp)

		vResp.EffectiveFrom = row.EffectiveFrom.Format(constant.YYYY_MM_DD)
		vResp.EffectiveTo = row.EffectiveTo.Format(constant.YYYY_MM_DD)
		vResp.PromoTypeName = vResp.GetPromoTypeName()
		vResp.BudgetReferenceTypeName = vResp.GetPromoBudgetReferenceTypeName()
		vResp.MaxDiscountTypeName = constant.GetQtyAmountPercentDisplayName(vResp.MaxDiscountType)
		vResp.PromoStatusDesc = vResp.GetPromoStatusDesc()
		vResp.BudgetControlLevelName = constant.GetPromoScopeLevelName(vResp.BudgetControlLevel)
		vResp.ExecutionLevelName = constant.GetPromoScopeLevelName(vResp.ExecutionLevel)
		vResp.MaxDiscountOutletUomName = constant.GetUomName(vResp.MaxDiscountOutletUom)

		// payTypeName := vResp.GeneratePayTypeName()
		// vResp.PayTypeName = payTypeName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *promotionServiceImpl) Update(promoID string, request entity.UpdatePromotionBody) (err error) {
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

	var promoModel model.Promotion
	err = structs.Automapper(request, &promoModel)
	if err != nil {
		return err
	}

	promoModel.CustID = ""
	promoModel.PromoStatusID = 0
	if request.PromoStatusID == 2 || request.PromoStatusID == 6 || request.PromoStatusID == 7 {
		promoModel.PromoStatusID = request.PromoStatusID
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.PromotionRepository.Update(txCtx, promoID, promoModel)
		if err != nil {
			return err
		}

		isHaveRewardProducts := false

		err := service.PromotionRepository.DeletePromoCriterias(txCtx, request.CustID, promoID)
		if err != nil {
			log.Error("DeletePromoCriterias, error:", err.Error())
		}

		err = service.PromotionRepository.DeletePromoAdditionalCriterias(txCtx, request.CustID, promoID)
		if err != nil {
			log.Error("DeletePromoRewardProducts, error:", err.Error())
		}

		for _, row := range request.PromoCriteria {

			var promoCritModel model.PromoCriteria
			err = structs.Automapper(row, &promoCritModel)
			if err != nil {
				return err
			}
			promoCritModel.CustID = request.CustID
			promoCritModel.PromoID = promoID
			promoCritModel.SlabID = nil
			err = service.PromotionRepository.StorePromoCriteria(txCtx, &promoCritModel)
			if err != nil {
				return err
			}

			if row.SlabRewardType == entity.PromoRewardTypeQuantity {
				isHaveRewardProducts = true
			}
		}

		err = service.PromotionRepository.DeletePromoRewardProducts(txCtx, request.CustID, promoID)
		if err != nil {
			log.Error("DeletePromoRewardProducts, error:", err.Error())
		}
		if isHaveRewardProducts {
			if len(request.RewardProduct) < 1 {
				return errors.New("minimum have 1 reward product")
			}

			for _, row := range request.RewardProduct {
				var rewardProductModel model.PromoRewardProduct
				err := structs.Automapper(row, &rewardProductModel)
				if err != nil {
					return err
				}
				rewardProductModel.CustID = request.CustID
				rewardProductModel.PromoID = promoID
				err = service.PromotionRepository.StorePromoRewardProduct(txCtx, &rewardProductModel)
				if err != nil {
					return err
				}
			}
		}

		for _, row := range request.PromoAdditionalCriteria {
			var promoCritAddModel model.PromoAdditionalCriteria
			err = structs.Automapper(row, &promoCritAddModel)
			if err != nil {
				return err
			}
			promoCritAddModel.CustID = request.CustID
			promoCritAddModel.PromoID = promoID
			promoCritAddModel.PromoAddCriteriaID = nil
			err = service.PromotionRepository.StorePromoAdditionalCriteria(txCtx, &promoCritAddModel)
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

func (service *promotionServiceImpl) Delete(params entity.DetailPromotionParams, deletedBy string) (err error) {
	c := context.Background()

	promo, err := service.PromotionRepository.FindByPromoID(params)
	if err != nil {
		return err
	}

	// validate if promo status is not draft or rejected
	if promo.PromoStatusID != 1 && promo.PromoStatusID != 4 {
		return errors.New("the promotion is not allow to be delete")
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.PromotionRepository.Delete(txCtx, params.CustID, params.PromoID)
		if err != nil {
			return err
		}

		err = service.PromotionRepository.DeletePromoCriterias(txCtx, params.CustID, params.PromoID)
		if err != nil {
			log.Error("DeletePromoCriterias, error:", err.Error())
		}

		err = service.PromotionRepository.DeletePromoAdditionalCriterias(txCtx, params.CustID, params.PromoID)
		if err != nil {
			log.Error("DeletePromoAdditionalCriterias, error:", err.Error())
		}

		err = service.PromotionRepository.DeletePromoRewardProducts(txCtx, params.CustID, params.PromoID)
		if err != nil {
			log.Error("DeletePromoRewardProducts, error:", err.Error())
		}
		return nil
	})

	return err
}

func (service *promotionServiceImpl) GetPromoCriterias(params entity.DetailPromotionParams, promoResponse *entity.Promotion) (err error) {
	promoCriterias, err := service.PromotionRepository.FindAllPromoCriteriasByPromoID(params)
	if err != nil {
		return err
	}

	for _, row := range promoCriterias {
		var promoCriteria entity.PromoCriteria
		err = structs.Automapper(row, &promoCriteria)
		if err != nil {
			return err
		}
		promoCriteria.CustID = ""
		promoCriteria.PromoID = ""
		promoCriteria.SlabRuleTypeName = constant.GetQtyAmountPercentDisplayName(promoCriteria.SlabRuleType)
		promoCriteria.SlabRewardTypeName = constant.GetQtyAmountPercentDisplayName(int(promoCriteria.SlabRewardType))
		promoCriteria.SlabRuleUomName = constant.GetUomName(promoCriteria.SlabRuleUom)
		promoCriteria.SlabRewardUomName = constant.GetUomName(promoCriteria.SlabRewardUom)

		promoResponse.PromoCriterias = append(promoResponse.PromoCriterias, promoCriteria)
	}

	return
}

func (service *promotionServiceImpl) GetPromoAdditionalCriterias(params entity.DetailPromotionParams, promoResponse *entity.Promotion) (err error) {
	promoAdditionalCriterias, err := service.PromotionRepository.FindAllPromoAdditionalCriteriasByPromoID(params)
	if err != nil {
		return err
	}

	for _, row := range promoAdditionalCriterias {
		var promoAddCriteria entity.PromoAdditionalCriteria
		err = structs.Automapper(row, &promoAddCriteria)
		if err != nil {
			return err
		}
		promoAddCriteria.CustID = ""
		promoAddCriteria.PromoID = ""
		promoAddCriteria.AttributeName = constant.GetPromoAttributeDisplayName(promoAddCriteria.Attribute)
		promoAddCriteria.ConditionName = constant.GetIncludeExcludeDisplayName(promoAddCriteria.Condition)
		promoAddCriteria.MinBuyTypeName = constant.GetQtyAmountPercentDisplayName(promoAddCriteria.MinBuyType)
		promoAddCriteria.MinBuyUomName = constant.GetUomName(promoAddCriteria.MinBuyUom)

		// Switch by Attribute
		switch promoAddCriteria.Attribute {

		case constant.AttrProduct:
			product, err := service.PromotionRepository.FindOneProductByProID(params.CustID, params.ParentCustId, promoAddCriteria.ReferenceID)
			if err == nil {
				promoAddCriteria.ReferenceCode = product.ReferenceCode
				promoAddCriteria.ReferenceName = product.ReferenceName
			}
		case constant.AttrOutletClass:
			outletClass, err := service.PromotionRepository.FindOneOutletClassByID(params.CustID, params.ParentCustId, promoAddCriteria.ReferenceID)
			if err == nil {
				promoAddCriteria.ReferenceCode = outletClass.ReferenceCode
				promoAddCriteria.ReferenceName = outletClass.ReferenceName
			}
		case constant.AttrOutletType:
			outletType, err := service.PromotionRepository.FindOneOutletTypeByID(params.CustID, params.ParentCustId, promoAddCriteria.ReferenceID)
			if err == nil {
				promoAddCriteria.ReferenceCode = outletType.ReferenceCode
				promoAddCriteria.ReferenceName = outletType.ReferenceName
			}
		case constant.AttrOutletGroup:
			outletGroup, err := service.PromotionRepository.FindOneOutletGroupByID(params.CustID, params.ParentCustId, promoAddCriteria.ReferenceID)
			if err == nil {
				promoAddCriteria.ReferenceCode = outletGroup.ReferenceCode
				promoAddCriteria.ReferenceName = outletGroup.ReferenceName
			}
		case constant.AttrSalesType:
			salesType, err := service.PromotionRepository.FindOneSalesTypeByID(params.CustID, params.ParentCustId, promoAddCriteria.ReferenceID)
			if err == nil {
				promoAddCriteria.ReferenceCode = salesType.ReferenceCode
				promoAddCriteria.ReferenceName = salesType.ReferenceName
			}
		case constant.AttrSalesTeam:
			salesTeam, err := service.PromotionRepository.FindOneSalesTeamByID(params.CustID, params.ParentCustId, promoAddCriteria.ReferenceID)
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

func (service *promotionServiceImpl) GetRewardProducts(params entity.DetailPromotionParams, promoResponse *entity.Promotion) (err error) {
	rewardProducts, err := service.PromotionRepository.FindAllRewardProductsByPromoID(params)
	if err != nil {
		return err
	}
	log.Info("rewardProducts:", structs.StructToJson(rewardProducts))
	for _, row := range rewardProducts {
		var rewardProduct entity.PromoRewardProduct
		err = structs.Automapper(row, &rewardProduct)
		if err != nil {
			return err
		}
		rewardProduct.CustID = ""
		rewardProduct.PromoID = ""
		productDetail, err := service.PromotionRepository.FindProductByProID(params.ParentCustId, row.ProID)
		if err == nil {
			rewardProduct.ProCode = productDetail.ProCode
			rewardProduct.ProName = productDetail.ProName
		}

		promoResponse.RewardProduct = append(promoResponse.RewardProduct, rewardProduct)
	}

	return
}

func (service *promotionServiceImpl) BulkUpdateStatus(request entity.BulkUpdateStatusPromotionBody) (err error) {
	c := context.Background()

	if request.PromoStatusID == 4 {
		if request.Remarks == "" {
			return errors.New("remarks is required")
		}
	} else {
		request.Remarks = ""
	}

	promos, err := service.PromotionRepository.FindAllByCustIdAndPromoID(request)
	if err != nil {
		return err
	}

	if len(promos) < 1 {
		return errors.New("promotion id not found")
	}

	promoTo := entity.Promotion{
		PromoStatusID: request.PromoStatusID,
	}
	promoStatusDescTo := promoTo.GetPromoStatusDesc()

	for _, r := range promos {
		promoFrom := entity.Promotion{
			PromoStatusID: r.PromoStatusID,
		}
		promoStatusDescFrom := promoFrom.GetPromoStatusDesc()

		errMsg := r.PromoID + ` with status '` + promoStatusDescFrom + `' is not allowed to update status to '` + promoStatusDescTo + `'`

		// validate from status draft to submitted
		if r.PromoStatusID == 1 && request.PromoStatusID != 2 {
			return errors.New(errMsg)
		}

		// validate from status submitted to approved or rejected
		if r.PromoStatusID == 2 && request.PromoStatusID != 3 && request.PromoStatusID != 4 {
			return errors.New(errMsg)
		}

		// validate from status rejected to submitted
		if r.PromoStatusID == 4 && request.PromoStatusID != 2 {
			return errors.New(errMsg)
		}

		// validate from status Approved to Active
		if r.PromoStatusID == 3 && request.PromoStatusID != 6 {
			return errors.New(errMsg)
		}

		// validate from status Active to Inactive
		if r.PromoStatusID == 6 && request.PromoStatusID != 7 {
			return errors.New(errMsg)
		}
	}

	// log.Info("promos:", structs.StructToJson(promos))

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.PromotionRepository.BulkUpdateStatus(txCtx, request)
		if err != nil {
			return err
		}

		for _, v := range request.PromoID {
			promoStatusLogModel := model.PromoStatusLog{
				CustID:        request.CustID,
				PromoID:       v,
				PromoStatusID: request.PromoStatusID,
				Remarks:       request.Remarks,
			}
			err = service.PromotionRepository.StoreStatusLog(txCtx, &promoStatusLogModel)
			if err != nil {
				log.Error("err:", err.Error())
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

func (service *promotionServiceImpl) ConsultPromotion(request entity.ConsultPromotionBody) (responses []entity.ConsultPromotionResponse, err error) {
	for index, detail := range request.Details {
		var detailConversion entity.CreateConversionBody
		detailConversion.CustId = request.CustID
		detailConversion.ProductId = detail.ProID
		detailConversion.Qty1 = int64(detail.Qty1)
		detailConversion.Qty2 = int64(detail.Qty2)
		detailConversion.Qty3 = int64(detail.Qty3)
		consultPromotionBodyDetail, _ := service.Conversion(detailConversion, request.CustID, request.ParentCustID)

		if consultPromotionBodyDetail.Qty1 != nil {
			request.Details[index].Qty1 = float64(*consultPromotionBodyDetail.Qty1)
		}
		if consultPromotionBodyDetail.Qty2 != nil {
			request.Details[index].Qty2 = float64(*consultPromotionBodyDetail.Qty2)
		}
		if consultPromotionBodyDetail.Qty3 != nil {
			request.Details[index].Qty3 = float64(*consultPromotionBodyDetail.Qty3)
		}
	}

	outlet, err := service.PromotionRepository.FindOutletByID(int64(request.OutletId), request.CustID, request.ParentCustID)
	if err != nil {
		log.Error(err.Error())
		return responses, fmt.Errorf("Outlet ID: %d not found", request.OutletId)
	}

	salesman, err := service.PromotionRepository.FindSalesmanByID(int64(request.SalesmanId), request.CustID, request.ParentCustID)
	if err != nil {
		log.Error(err.Error())
		return responses, fmt.Errorf("Salesman ID: %d not found", request.SalesmanId)
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
				case 3:
					response.SlabReward = math.Round((float64(subTotalValidatedPromoAdditionalCriteriaByProductGroups[promoCriteria.PromoID]) * validatedPromoCriterias[promoCriteria.PromoID].SlabReward) / 100.0)
				case 2:
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
						if rewardProductConversionResut.Qty1 != nil {
							rewardProduct.Qty1 = float64(*rewardProductConversionResut.Qty1)
						}
						if rewardProductConversionResut.Qty2 != nil {
							rewardProduct.Qty2 = float64(*rewardProductConversionResut.Qty2)
						}
						if rewardProductConversionResut.Qty3 != nil {
							rewardProduct.Qty3 = float64(*rewardProductConversionResut.Qty3)
						}
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

func (service *promotionServiceImpl) Conversion(conversionBody entity.CreateConversionBody, custID string, parentCustID string) (response entity.OrderConversionResponse, err error) {
	if conversionBody.ProductId == 0 {
		qty1 := conversionBody.Qty1
		qty2 := conversionBody.Qty2
		qty3 := conversionBody.Qty3
		response.Qty1 = &qty1
		response.Qty2 = &qty2
		response.Qty3 = &qty3
		response.TotalQty = qty1 + qty2 + qty3
		return response, nil
	}
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

	qty1Ptr := qty1
	qty2Ptr := qty2
	qty3Ptr := qty3
	response.Qty1 = &qty1Ptr
	response.Qty2 = &qty2Ptr
	response.Qty3 = &qty3Ptr

	response.TotalQty = (int64(product.ConvUnit2)*int64(product.ConvUnit3))*qty3 + (int64(product.ConvUnit2) * qty2) + qty1

	return response, err
}

func (service *promotionServiceImpl) StoreV2(request entity.CreatePromotionV2Body) (err error) {
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

	var promoModel model.PromotionV2
	err = structs.Automapper(request, &promoModel)
	if err != nil {
		return err
	}

	if err = applyPromotionV2ExtendedFields(&promoModel, request.BudgetID, request.ClaimDateFrom, request.ClaimDateTo, request.VatRate, request.WhtRate); err != nil {
		return err
	}

	promoModel.UpdatedAt = time.Now().UTC()

	// var promoStatusLogModel model.PromoStatusLog
	// err = structs.Automapper(request, &promoStatusLogModel)
	// if err != nil {
	// 	return err
	// }

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.PromotionV2Repository.Store(txCtx, &promoModel)
		if err != nil {
			log.Error("err:", err.Error())
			return err
		}

		if request.PromoType == entity.PromotionTypeSlab {
			var promoSlabs []model.PromotionV2Slabs
			for _, slab := range request.Slabs {
				var promoSlab model.PromotionV2Slabs
				err = structs.Automapper(slab, &promoSlab)
				if err != nil {
					return err
				}
				if promoSlab.PerScope != nil && *promoSlab.PerScope == "" {
					promoSlab.PerScope = nil
				}
				promoSlabs = append(promoSlabs, promoSlab)
			}

			err = service.PromotionV2Repository.StoreSlabs(txCtx, promoSlabs)
			if err != nil {
				log.Error("store slabs err:", err.Error())
				return err
			}
		}

		if request.PromoType == entity.PromotionTypeStrata {
			var promoStratas []model.PromotionV2Strata
			for _, strata := range request.Strata {
				var promoStrata model.PromotionV2Strata
				err = structs.Automapper(strata, &promoStrata)
				if err != nil {
					return err
				}
				promoStrata.ClaimRealizationPct = &request.ClaimRealizationPct
				promoStratas = append(promoStratas, promoStrata)
			}
			err = service.PromotionV2Repository.StoreStrata(txCtx, promoStratas)
			if err != nil {
				log.Error("store strata err:", err.Error())
				return err
			}
		}

		// promo product criterias
		var promoProductCriterias []model.PromotionProductCriteria
		for _, item := range request.ProductCriteria {
			var promoProductCriteria model.PromotionProductCriteria
			promoProductCriteria.CustID = request.ParentCustID // store under parent cust to align with queries
			promoProductCriteria.PromoID = request.PromoID
			promoProductCriteria.ProID = item.ProID
			promoProductCriteria.Mandatory = item.Mandatory
			if item.MinBuyType != nil {
				minBuyType := model.RuleType(*item.MinBuyType)
				promoProductCriteria.MinBuyType = &minBuyType
			}
			promoProductCriteria.MinBuyQty = item.MinBuyQty
			promoProductCriteria.MinBuyValue = item.MinBuyValue
			if item.MinBuyUom != nil {
				minBuyUom := model.UomType(*item.MinBuyUom)
				promoProductCriteria.MinBuyUom = &minBuyUom
			}
			promoProductCriterias = append(promoProductCriterias, promoProductCriteria)
		}
		// log.Info("promoProductCriterias:", structs.StructToJson(promoProductCriterias))
		err = service.PromotionV2Repository.StoreProductCriteria(txCtx, promoProductCriterias)
		if err != nil {
			return err
		}

		// promo reward products
		if len(request.RewardProducts) > 0 {
			var promoRewardProducts []model.PromotionRewardProduct
			for _, rewardProduct := range request.RewardProducts {
				var promoRewardProduct model.PromotionRewardProduct
				err = structs.Automapper(rewardProduct, &promoRewardProduct)
				promoRewardProduct.CustID = request.CustID
				promoRewardProduct.PromoID = request.PromoID
				promoRewardProducts = append(promoRewardProducts, promoRewardProduct)
			}
			err = service.PromotionV2Repository.StoreRewardProducts(txCtx, promoRewardProducts)
			if err != nil {
				log.Error("store reward products err:", err.Error())
				return err
			}
		}

		// Store Coverage Distributors
		if len(request.CoverageDistributors) > 0 {
			var coverageDistributors []model.PromotionCoverageDistributors
			for _, coverageDistributor := range request.CoverageDistributors {
				var promoCoverageDistributor model.PromotionCoverageDistributors
				promoCoverageDistributor.CustID = request.CustID
				promoCoverageDistributor.PromoID = request.PromoID
				promoCoverageDistributor.DistributorID = coverageDistributor.DistributorID
				coverageDistributors = append(coverageDistributors, promoCoverageDistributor)
			}
			err = service.PromotionV2Repository.StoreCoverageDistributors(txCtx, coverageDistributors)
			if err != nil {
				log.Error("store coverage distributors err:", err.Error())
				return err
			}
		}

		// Store Outlet Criteria
		// Create outlet criteria record
		var promoOutletCriteria model.PromotionOutletCriteria
		promoOutletCriteria.CustID = request.CustID
		promoOutletCriteria.PromoID = request.PromoID
		promoOutletCriteria.SelectionType = model.OutletSelType(request.OutletCriteria.SelectionType)

		// Store outlet criteria and get the ID
		outletCriteriaID, err := service.PromotionV2Repository.StoreOutletCriteria(txCtx, &promoOutletCriteria)
		if err != nil {
			log.Error("store outlet criteria err:", err.Error())
			return err
		}

		// Store outlet attributes based on selection type
		if request.OutletCriteria.SelectionType == "by_attribute" {
			// Store outlet type attribute
			var outletAttributeTypes []model.PromotionOutletAttributeType
			for _, outletTypeID := range request.OutletCriteria.OutletTypeIDs {
				var outletAttributeType model.PromotionOutletAttributeType
				outletAttributeType.CustID = request.CustID
				outletAttributeType.CriteriaID = outletCriteriaID
				outletAttributeType.OutletTypeID = outletTypeID
				outletAttributeTypes = append(outletAttributeTypes, outletAttributeType)
			}
			if len(outletAttributeTypes) > 0 {
				err = service.PromotionV2Repository.StoreOutletAttributeType(txCtx, outletAttributeTypes)
				if err != nil {
					log.Error("store outlet attribute type err:", err.Error())
					return err
				}
			}

			// Store outlet group attribute
			var outletAttributeGroups []model.PromotionOutletAttributeGroup
			for _, outletGroupID := range request.OutletCriteria.OutletGroupIDs {
				var outletAttributeGroup model.PromotionOutletAttributeGroup
				outletAttributeGroup.CustID = request.CustID
				outletAttributeGroup.CriteriaID = outletCriteriaID
				outletAttributeGroup.OutletGroupID = outletGroupID
				outletAttributeGroups = append(outletAttributeGroups, outletAttributeGroup)
			}
			if len(outletAttributeGroups) > 0 {
				err = service.PromotionV2Repository.StoreOutletAttributeGroup(txCtx, outletAttributeGroups)
				if err != nil {
					log.Error("store outlet attribute group err:", err.Error())
					return err
				}
			}

			// Store outlet class attribute
			var outletAttributeClasses []model.PromotionOutletAttributeClass
			for _, outletClassID := range request.OutletCriteria.OutletClassIDs {
				var outletAttributeClass model.PromotionOutletAttributeClass
				outletAttributeClass.CustID = request.CustID
				outletAttributeClass.CriteriaID = outletCriteriaID
				outletAttributeClass.OutletClassID = outletClassID
				outletAttributeClasses = append(outletAttributeClasses, outletAttributeClass)
			}
			if len(outletAttributeClasses) > 0 {
				err = service.PromotionV2Repository.StoreOutletAttributeClass(txCtx, outletAttributeClasses)
				if err != nil {
					log.Error("store outlet attribute class err:", err.Error())
					return err
				}
			}
		} else if request.OutletCriteria.SelectionType == "by_outlet" {
			// Store specific outlet IDs for by_outlet selection
			var outletsSelected []model.PromotionOutletsSelected
			for _, outletID := range request.OutletCriteria.OutletIDs {
				var outletSelected model.PromotionOutletsSelected
				outletSelected.CustID = request.CustID
				outletSelected.CriteriaID = outletCriteriaID
				outletSelected.OutletID = outletID
				outletsSelected = append(outletsSelected, outletSelected)
			}
			if len(outletsSelected) > 0 {
				err = service.PromotionV2Repository.StoreOutletsSelected(txCtx, outletsSelected)
				if err != nil {
					log.Error("store outlets selected err:", err.Error())
					return err
				}
			}

		}

		// Store sales team attributes
		var outletAttributeSalesTeams []model.PromotionOutletAttributeSalesTeam
		for _, salesTeamID := range request.OutletCriteria.SalesTeamIDs {
			var outletAttributeSalesTeam model.PromotionOutletAttributeSalesTeam
			outletAttributeSalesTeam.CustID = request.CustID
			outletAttributeSalesTeam.CriteriaID = outletCriteriaID
			outletAttributeSalesTeam.SalesTeamID = salesTeamID
			outletAttributeSalesTeams = append(outletAttributeSalesTeams, outletAttributeSalesTeam)
		}
		if len(outletAttributeSalesTeams) > 0 {
			err = service.PromotionV2Repository.StoreOutletAttributeSalesTeam(txCtx, outletAttributeSalesTeams)
			if err != nil {
				log.Error("store outlet attribute sales team err:", err.Error())
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

func (service *promotionServiceImpl) DetailV2(params entity.DetailPromotionParams) (response entity.PromotionV2, err error) {
	promo, err := service.PromotionV2Repository.FindByPromoID(params)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(promo, &response)
	if err != nil {
		return response, err
	}

	response.EffectiveFrom = promo.EffectiveFrom.Format(constant.YYYY_MM_DD)
	response.EffectiveTo = promo.EffectiveTo.Format(constant.YYYY_MM_DD)
	mapPromotionV2ExtendedFieldsToEntity(promo, &response)
	mapPromotionV2ResponseCustID(promo, &response)

	promoSlabs, err := service.PromotionV2Repository.FindPromoSlabsByPromoID(params)
	if err != nil {
		return response, err
	}
	for _, slab := range promoSlabs {
		var slabItem entity.PromoSlabItem
		err = structs.Automapper(slab, &slabItem)
		if err != nil {
			return response, err
		}
		slabItem.CustID = ""
		slabItem.PromoID = ""
		response.Slabs = append(response.Slabs, slabItem)
	}

	promoStratas, err := service.PromotionV2Repository.FindPromoStratasByPromoID(params)
	if err != nil {
		return response, err
	}
	for _, strata := range promoStratas {
		var strataItem entity.PromoStrataItem
		err = structs.Automapper(strata, &strataItem)
		if err != nil {
			return response, err
		}
		strataItem.CustID = ""
		strataItem.PromoID = ""
		response.Strata = append(response.Strata, strataItem)
	}

	promoProductCriterias, err := service.PromotionV2Repository.FindPromoProductCriteriasByPromoID(params)
	if err != nil {
		return response, err
	}
	for _, productCriteria := range promoProductCriterias {
		var productCriteriaItem entity.PromoProductCriteria
		err = structs.Automapper(productCriteria, &productCriteriaItem)
		if err != nil {
			return response, err
		}
		productCriteriaItem.CustID = ""
		productCriteriaItem.PromoID = ""
		response.ProductCriteria = append(response.ProductCriteria, productCriteriaItem)
	}

	rewardProducts, err := service.PromotionV2Repository.FindPromoRewardProductsByPromoID(params)
	if err != nil {
		return response, err
	}
	for _, rewardProduct := range rewardProducts {
		var rewardProductItem entity.PromotionRewardProduct
		err = structs.Automapper(rewardProduct, &rewardProductItem)
		if err != nil {
			return response, err
		}
		rewardProductItem.CustID = ""
		rewardProductItem.PromoID = ""
		response.RewardProducts = append(response.RewardProducts, rewardProductItem)
	}

	coverageDistributors, err := service.PromotionV2Repository.FindCoverageDistributorsByPromoID(params)
	if err != nil {
		return response, err
	}
	response.CoverageDistributors = make([]entity.PromoCoverageDistributor, 0, len(coverageDistributors))
	for _, coverageDistributor := range coverageDistributors {
		response.CoverageDistributors = append(response.CoverageDistributors, mapPromoCoverageDistributorFromModel(coverageDistributor))
	}

	outletCriteriaList, err := service.PromotionV2Repository.FindOutletCriteriaWithPreloads(params)
	if err != nil {
		return response, err
	}

	// Convert outlet criteria list to entity format
	if len(outletCriteriaList) > 0 {
		var outletCriteriaItem entity.PromoOutletCriteria
		err = structs.Automapper(outletCriteriaList[0], &outletCriteriaItem)
		if err != nil {
			return response, err
		}
		outletCriteriaItem.CustID = ""
		outletCriteriaItem.PromoID = ""
		outletCriteriaItem.SelectedOutlets = []entity.PromoOutletSelected{}
		for _, outlet := range outletCriteriaList[0].SelectedOutlets {
			outletCriteriaItem.SelectedOutlets = append(outletCriteriaItem.SelectedOutlets, mapPromoOutletSelectedFromModel(outlet))
		}
		outletCriteriaItem.SelectedOutletTypes = []entity.PromoOutletTypeSelected{}
		for _, outletType := range outletCriteriaList[0].AttributeTypes {
			var outletTypeSelected entity.PromoOutletTypeSelected
			err = structs.Automapper(outletType, &outletTypeSelected)
			if err != nil {
				return response, err
			}
			outletCriteriaItem.SelectedOutletTypes = append(outletCriteriaItem.SelectedOutletTypes, outletTypeSelected)
		}
		outletCriteriaItem.SelectedOutletGroups = []entity.PromoOutletGroupSelected{}
		for _, outletGroup := range outletCriteriaList[0].AttributeGroups {
			var outletGroupSelected entity.PromoOutletGroupSelected
			err = structs.Automapper(outletGroup, &outletGroupSelected)
			if err != nil {
				return response, err
			}
			outletCriteriaItem.SelectedOutletGroups = append(outletCriteriaItem.SelectedOutletGroups, outletGroupSelected)
		}
		outletCriteriaItem.SelectedOutletClasses = []entity.PromoOutletClassSelected{}
		for _, outletClass := range outletCriteriaList[0].AttributeClasses {
			var outletClassSelected entity.PromoOutletClassSelected
			err = structs.Automapper(outletClass, &outletClassSelected)
			if err != nil {
				return response, err
			}
			outletCriteriaItem.SelectedOutletClasses = append(outletCriteriaItem.SelectedOutletClasses, outletClassSelected)
		}
		outletCriteriaItem.SelectedSalesTeams = []entity.PromoOutletSalesTeamSelected{}
		for _, salesTeam := range outletCriteriaList[0].AttributeSalesTeams {
			outletCriteriaItem.SelectedSalesTeams = append(outletCriteriaItem.SelectedSalesTeams, mapPromoOutletSalesTeamFromModel(salesTeam))
		}
		response.OutletCriteria = outletCriteriaItem
	}

	return response, nil
}

func (service *promotionServiceImpl) ListV2(dataFilter entity.PromotionV2QueryFilter) (data []entity.PromotionV2, total int64, lastPage int, err error) {
	promotions, total, lastPage, err := service.PromotionV2Repository.FindAllByCustID(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range promotions {
		var vResp entity.PromotionV2
		structs.Automapper(row, &vResp)
		vResp.EffectiveFrom = row.EffectiveFrom.Format(constant.YYYY_MM_DD)
		vResp.EffectiveTo = row.EffectiveTo.Format(constant.YYYY_MM_DD)
		mapPromotionV2ExtendedFieldsToEntity(row, &vResp)
		mapPromotionV2ResponseCustID(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *promotionServiceImpl) UpdateV2(promoID string, request entity.UpdatePromotionV2Body) (err error) {
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

	var promoModel model.PromotionV2
	err = structs.Automapper(request, &promoModel)
	if err != nil {
		return err
	}

	if err = applyPromotionV2ExtendedFields(&promoModel, request.BudgetID, request.ClaimDateFrom, request.ClaimDateTo, request.VatRate, request.WhtRate); err != nil {
		return err
	}

	promoModel.DistributorCustID = request.CustID

	promoModel.CustID = request.ParentCustID
	promoModel.UpdatedAt = time.Now().UTC()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		// Update main promotion record
		err := service.PromotionV2Repository.Update(txCtx, promoID, &promoModel)
		if err != nil {
			log.Error("err:", err.Error())
			return err
		}

		// Delete existing related records
		err = service.PromotionV2Repository.DeleteSlabs(txCtx, request.ParentCustID, promoID)
		if err != nil {
			log.Error("DeleteSlabs, error:", err.Error())
		}

		err = service.PromotionV2Repository.DeleteStratas(txCtx, request.ParentCustID, promoID)
		if err != nil {
			log.Error("DeleteStratas, error:", err.Error())
		}

		// delete product criteria for both owner (parent) and distributor (legacy rows)
		custIDs := []string{request.ParentCustID}
		if request.CustID != "" && request.CustID != request.ParentCustID {
			custIDs = append(custIDs, request.CustID)
		}

		err = service.PromotionV2Repository.DeleteProductCriteria(txCtx, custIDs, promoID)
		if err != nil {
			log.Error("DeleteProductCriteria, error:", err.Error())
		}

		err = service.PromotionV2Repository.DeleteRewardProducts(txCtx, request.ParentCustID, promoID)
		if err != nil {
			log.Error("DeleteRewardProducts, error:", err.Error())
		}

		err = service.PromotionV2Repository.DeleteCoverageDistributors(txCtx, request.ParentCustID, promoID)
		if err != nil {
			log.Error("DeleteCoverageDistributors, error:", err.Error())
		}

		err = service.PromotionV2Repository.DeleteOutletCriteria(txCtx, request.ParentCustID, promoID)
		if err != nil {
			log.Error("DeleteOutletCriteria, error:", err.Error())
		}

		// Insert new related records
		if request.PromoType == entity.PromotionTypeSlab {
			var promoSlabs []model.PromotionV2Slabs
			for _, slab := range request.Slabs {
				var promoSlab model.PromotionV2Slabs
				err = structs.Automapper(slab, &promoSlab)
				if err != nil {
					return err
				}
				promoSlab.CustID = request.ParentCustID
				promoSlab.PromoID = promoID
				if promoSlab.PerScope != nil && *promoSlab.PerScope == "" {
					promoSlab.PerScope = nil
				}
				promoSlabs = append(promoSlabs, promoSlab)
			}

			err = service.PromotionV2Repository.StoreSlabs(txCtx, promoSlabs)
			if err != nil {
				log.Error("store slabs err:", err.Error())
				return err
			}
		}

		if request.PromoType == entity.PromotionTypeStrata {
			var promoStratas []model.PromotionV2Strata
			for _, strata := range request.Strata {
				var promoStrata model.PromotionV2Strata
				err = structs.Automapper(strata, &promoStrata)
				if err != nil {
					return err
				}
				promoStrata.CustID = request.ParentCustID
				promoStrata.PromoID = promoID
				promoStrata.ClaimRealizationPct = &request.ClaimRealizationPct
				if promoStrata.PerScope != nil && *promoStrata.PerScope == "" {
					promoStrata.PerScope = nil
				}
				promoStratas = append(promoStratas, promoStrata)
			}
			err = service.PromotionV2Repository.StoreStrata(txCtx, promoStratas)
			if err != nil {
				log.Error("store strata err:", err.Error())
				return err
			}
		}

		// promo product criterias
		var promoProductCriterias []model.PromotionProductCriteria
		for _, item := range request.ProductCriteria {
			var promoProductCriteria model.PromotionProductCriteria
			promoProductCriteria.CustID = request.ParentCustID
			promoProductCriteria.PromoID = promoID
			promoProductCriteria.ProID = item.ProID
			promoProductCriteria.Mandatory = item.Mandatory
			if item.MinBuyType != nil {
				minBuyType := model.RuleType(*item.MinBuyType)
				promoProductCriteria.MinBuyType = &minBuyType
			}
			promoProductCriteria.MinBuyQty = item.MinBuyQty
			promoProductCriteria.MinBuyValue = item.MinBuyValue
			if item.MinBuyUom != nil {
				minBuyUom := model.UomType(*item.MinBuyUom)
				promoProductCriteria.MinBuyUom = &minBuyUom
			}
			promoProductCriterias = append(promoProductCriterias, promoProductCriteria)
		}
		err = service.PromotionV2Repository.StoreProductCriteria(txCtx, promoProductCriterias)
		if err != nil {
			return err
		}

		// promo reward products
		if len(request.RewardProducts) > 0 {
			var promoRewardProducts []model.PromotionRewardProduct
			for _, rewardProduct := range request.RewardProducts {
				var promoRewardProduct model.PromotionRewardProduct
				err = structs.Automapper(rewardProduct, &promoRewardProduct)
				promoRewardProduct.CustID = request.ParentCustID
				promoRewardProduct.PromoID = promoID
				promoRewardProducts = append(promoRewardProducts, promoRewardProduct)
			}
			err = service.PromotionV2Repository.StoreRewardProducts(txCtx, promoRewardProducts)
			if err != nil {
				log.Error("store reward products err:", err.Error())
				return err
			}
		}

		// Store Coverage Distributors
		if len(request.CoverageDistributors) > 0 {
			var coverageDistributors []model.PromotionCoverageDistributors
			for _, coverageDistributor := range request.CoverageDistributors {
				var promoCoverageDistributor model.PromotionCoverageDistributors
				promoCoverageDistributor.CustID = request.ParentCustID
				promoCoverageDistributor.PromoID = promoID
				promoCoverageDistributor.DistributorID = coverageDistributor.DistributorID
				coverageDistributors = append(coverageDistributors, promoCoverageDistributor)
			}
			err = service.PromotionV2Repository.StoreCoverageDistributors(txCtx, coverageDistributors)
			if err != nil {
				log.Error("store coverage distributors err:", err.Error())
				return err
			}
		}

		// Store Outlet Criteria
		// Create outlet criteria record
		var promoOutletCriteria model.PromotionOutletCriteria
		promoOutletCriteria.CustID = request.ParentCustID
		promoOutletCriteria.PromoID = promoID
		promoOutletCriteria.SelectionType = model.OutletSelType(request.OutletCriteria.SelectionType)

		// Store outlet criteria and get the ID
		outletCriteriaID, err := service.PromotionV2Repository.StoreOutletCriteria(txCtx, &promoOutletCriteria)
		if err != nil {
			log.Error("store outlet criteria err:", err.Error())
			return err
		}

		// Store outlet attributes based on selection type
		if request.OutletCriteria.SelectionType == "by_attribute" {
			// Store outlet type attribute
			var outletAttributeTypes []model.PromotionOutletAttributeType
			for _, outletTypeID := range request.OutletCriteria.OutletTypeIDs {
				var outletAttributeType model.PromotionOutletAttributeType
				outletAttributeType.CustID = request.ParentCustID
				outletAttributeType.CriteriaID = outletCriteriaID
				outletAttributeType.OutletTypeID = outletTypeID
				outletAttributeTypes = append(outletAttributeTypes, outletAttributeType)
			}
			if len(outletAttributeTypes) > 0 {
				err = service.PromotionV2Repository.StoreOutletAttributeType(txCtx, outletAttributeTypes)
				if err != nil {
					log.Error("store outlet attribute type err:", err.Error())
					return err
				}
			}

			// Store outlet group attribute
			var outletAttributeGroups []model.PromotionOutletAttributeGroup
			for _, outletGroupID := range request.OutletCriteria.OutletGroupIDs {
				var outletAttributeGroup model.PromotionOutletAttributeGroup
				outletAttributeGroup.CustID = request.ParentCustID
				outletAttributeGroup.CriteriaID = outletCriteriaID
				outletAttributeGroup.OutletGroupID = outletGroupID
				outletAttributeGroups = append(outletAttributeGroups, outletAttributeGroup)
			}
			if len(outletAttributeGroups) > 0 {
				err = service.PromotionV2Repository.StoreOutletAttributeGroup(txCtx, outletAttributeGroups)
				if err != nil {
					log.Error("store outlet attribute group err:", err.Error())
					return err
				}
			}

			// Store outlet class attribute
			var outletAttributeClasses []model.PromotionOutletAttributeClass
			for _, outletClassID := range request.OutletCriteria.OutletClassIDs {
				var outletAttributeClass model.PromotionOutletAttributeClass
				outletAttributeClass.CustID = request.ParentCustID
				outletAttributeClass.CriteriaID = outletCriteriaID
				outletAttributeClass.OutletClassID = outletClassID
				outletAttributeClasses = append(outletAttributeClasses, outletAttributeClass)
			}
			if len(outletAttributeClasses) > 0 {
				err = service.PromotionV2Repository.StoreOutletAttributeClass(txCtx, outletAttributeClasses)
				if err != nil {
					log.Error("store outlet attribute class err:", err.Error())
					return err
				}
			}
		} else if request.OutletCriteria.SelectionType == "by_outlet" {
			// Store specific outlet IDs for by_outlet selection
			var outletsSelected []model.PromotionOutletsSelected
			for _, outletID := range request.OutletCriteria.OutletIDs {
				var outletSelected model.PromotionOutletsSelected
				outletSelected.CustID = request.ParentCustID
				outletSelected.CriteriaID = outletCriteriaID
				outletSelected.OutletID = outletID
				outletsSelected = append(outletsSelected, outletSelected)
			}
			if len(outletsSelected) > 0 {
				err = service.PromotionV2Repository.StoreOutletsSelected(txCtx, outletsSelected)
				if err != nil {
					log.Error("store outlets selected err:", err.Error())
					return err
				}
			}
		}

		// Store sales team attributes
		var outletAttributeSalesTeams []model.PromotionOutletAttributeSalesTeam
		for _, salesTeamID := range request.OutletCriteria.SalesTeamIDs {
			var outletAttributeSalesTeam model.PromotionOutletAttributeSalesTeam
			outletAttributeSalesTeam.CustID = request.ParentCustID
			outletAttributeSalesTeam.CriteriaID = outletCriteriaID
			outletAttributeSalesTeam.SalesTeamID = salesTeamID
			outletAttributeSalesTeams = append(outletAttributeSalesTeams, outletAttributeSalesTeam)
		}
		if len(outletAttributeSalesTeams) > 0 {
			err = service.PromotionV2Repository.StoreOutletAttributeSalesTeam(txCtx, outletAttributeSalesTeams)
			if err != nil {
				log.Error("store outlet attribute sales team err:", err.Error())
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

func (service *promotionServiceImpl) ExistsPromo(custID, promoID string) (bool, error) {
	exists, err := service.PromotionV2Repository.ExistsPromo(custID, promoID)
	return exists, err
}

func (service *promotionServiceImpl) UpdateV2Status(promoID string, promoStatus entity.PromotionV2Status, req entity.UpdateStatusPromotionV2Body) (err error) {
	c := context.Background()

	// Convert entity status to model status
	var modelStatus model.PromotionStatus
	switch promoStatus {
	case entity.PromoStatusDraft:
		modelStatus = model.PromoStatusDraft
	case entity.PromoStatusSubmit:
		modelStatus = model.PromoStatusSubmit
	case entity.PromoStatusApproved:
		modelStatus = model.PromoStatusApproved
	case entity.PromoStatusRejected:
		modelStatus = model.PromoStatusRejected
	case entity.PromoStatusInactive:
		modelStatus = model.PromoStatusInactive
	case entity.PromoStatusActive:
		modelStatus = model.PromoStatusActive
	case entity.PromoStatusClosed:
		modelStatus = model.PromoStatusClosed
	default:
		return errors.New("invalid promotion status")
	}

	if req.EffectiveTo != "" {
		effectiveTo, parseErr := str.DateStrToRfc3339String(req.EffectiveTo)
		if parseErr != nil {
			return parseErr
		}
		req.EffectiveTo = effectiveTo
	}

	var promoModel model.PromotionV2
	err = structs.Automapper(req, &promoModel)
	if err != nil {
		return err
	}

	promoModel.PromoStatus = modelStatus
	promoModel.UpdatedAt = time.Now().UTC()
	promoModel.CustID = req.ParentCustID
	promoModel.DistributorCustID = req.CustID

	// Update the promotion status
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.PromotionV2Repository.Update(txCtx, promoID, &promoModel)
		if err != nil {
			log.Error("Error updating promotion status:", err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
func (service *promotionServiceImpl) CloseExpiredPromotions() (err error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Error("CloseExpiredPromotions, LoadLocation, err:", err.Error())
		return err
	}

	now := time.Now().In(loc)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	affected, err := service.PromotionV2Repository.CloseExpiredPromotionStatuses(todayStart)
	if err != nil {
		log.Error("CloseExpiredPromotions, CloseExpiredPromotionStatuses, err:", err.Error())
		return err
	}
	if affected > 0 {
		log.Info("CloseExpiredPromotions closed ", affected, " promotion(s)")
	}

	return nil
}

func (service *promotionServiceImpl) DuplicateV2(params entity.DetailPromotionParams) (newPromoID string, err error) {
	c := context.Background()

	// Get the original promotion details
	originalPromo, err := service.PromotionV2Repository.FindByPromoID(params)
	if err != nil {
		log.Error("Error finding original promotion:", err.Error())
		return "", err
	}

	// Generate new promotion ID with auto-incrementing sequence number
	newPromoID, err = service.generateNextSequenceNumber(params.ParentCustId, originalPromo.PromoID)
	if err != nil {
		log.Error("Error generating sequence number:", err.Error())
		return "", err
	}

	// Create the duplicate promotion
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		// Create new promotion model with updated fields
		duplicatePromo := originalPromo
		duplicatePromo.PromoID = newPromoID
		duplicatePromo.PromoStatus = model.PromoStatusDraft // Set to draft as per requirements
		duplicatePromo.CreatedAt = time.Now().UTC()
		duplicatePromo.UpdatedAt = time.Now().UTC()
		duplicatePromo.CreatedBy = params.UserFullname
		duplicatePromo.UpdatedBy = params.UserFullname

		// Store the duplicate promotion
		err := service.PromotionV2Repository.Store(txCtx, &duplicatePromo)
		if err != nil {
			log.Error("Error storing duplicate promotion:", err.Error())
			return err
		}

		// Duplicate all related data
		// 1. Duplicate slabs
		if originalPromo.PromoType == model.PromotionTypeSlab {
			slabs, err := service.PromotionV2Repository.FindPromoSlabsByPromoID(params)
			if err != nil {
				log.Error("Error finding slabs:", err.Error())
				return err
			}
			if slabs != nil && len(slabs) > 0 {
				// Update slab data for new promotion
				for i := range slabs {
					slabs[i].PromoID = newPromoID
					slabs[i].CustID = params.ParentCustId
					slabs[i].ID = ""
				}
				err = service.PromotionV2Repository.StoreSlabs(txCtx, slabs)
				if err != nil {
					log.Error("Error storing duplicate slabs:", err.Error())
					return err
				}
			}
		}

		// 2. Duplicate strata
		if originalPromo.PromoType == model.PromotionTypeStrata {
			stratas, err := service.PromotionV2Repository.FindPromoStratasByPromoID(params)
			if err != nil {
				log.Error("Error finding stratas:", err.Error())
				return err
			}
			if stratas != nil && len(stratas) > 0 {
				// Update strata data for new promotion
				for i := range stratas {
					stratas[i].PromoID = newPromoID
					stratas[i].CustID = params.ParentCustId
					stratas[i].ID = ""
				}
				err = service.PromotionV2Repository.StoreStrata(txCtx, stratas)
				if err != nil {
					log.Error("Error storing duplicate stratas:", err.Error())
					return err
				}
			}
		}
		log.Info("pass insert strata")

		// 3. Duplicate product criteria
		productCriterias, err := service.PromotionV2Repository.FindPromoProductCriteriasByPromoID(params)
		if err != nil {
			log.Error("Error finding product criterias:", err.Error())
			return err
		}
		if productCriterias != nil && len(productCriterias) > 0 {
			// Update product criteria data for new promotion
			for i := range productCriterias {
				productCriterias[i].PromoID = newPromoID
				productCriterias[i].CustID = params.ParentCustId
				productCriterias[i].ID = ""
			}
			err = service.PromotionV2Repository.StoreProductCriteria(txCtx, productCriterias)
			if err != nil {
				log.Error("Error storing duplicate product criterias:", err.Error())
				return err
			}
		}

		// 4. Duplicate reward products
		rewardProducts, err := service.PromotionV2Repository.FindPromoRewardProductsByPromoID(params)
		if err != nil {
			log.Error("Error finding reward products:", err.Error())
			return err
		}
		if rewardProducts != nil && len(rewardProducts) > 0 {
			// Update reward products data for new promotion
			for i := range rewardProducts {
				rewardProducts[i].PromoID = newPromoID
				rewardProducts[i].CustID = params.ParentCustId
				rewardProducts[i].ID = ""
			}
			err = service.PromotionV2Repository.StoreRewardProducts(txCtx, rewardProducts)
			if err != nil {
				log.Error("Error storing duplicate reward products:", err.Error())
				return err
			}
		}

		// 5. Duplicate coverage distributors
		coverageDistributors, err := service.PromotionV2Repository.FindCoverageDistributorsByPromoID(params)
		if err != nil {
			log.Error("Error finding coverage distributors:", err.Error())
			return err
		}
		if coverageDistributors != nil && len(coverageDistributors) > 0 {
			// Update coverage distributors data for new promotion
			for i := range coverageDistributors {
				coverageDistributors[i].PromoID = newPromoID
				coverageDistributors[i].CustID = params.ParentCustId
				coverageDistributors[i].ID = ""
			}
			err = service.PromotionV2Repository.StoreCoverageDistributors(txCtx, coverageDistributors)
			if err != nil {
				log.Error("Error storing duplicate coverage distributors:", err.Error())
				return err
			}
		}

		// 6. Duplicate outlet criteria
		outletCriteriaList, err := service.PromotionV2Repository.FindOutletCriteriaWithPreloads(params)
		if err != nil {
			log.Error("Error finding outlet criteria:", err.Error())
			return err
		}
		if outletCriteriaList != nil && len(outletCriteriaList) > 0 {
			// Create new outlet criteria
			var newOutletCriteria model.PromotionOutletCriteria
			newOutletCriteria.CustID = params.ParentCustId
			newOutletCriteria.PromoID = newPromoID
			newOutletCriteria.SelectionType = outletCriteriaList[0].SelectionType
			newOutletCriteria.ID = ""

			// Store outlet criteria and get the ID
			outletCriteriaID, err := service.PromotionV2Repository.StoreOutletCriteria(txCtx, &newOutletCriteria)
			if err != nil {
				log.Error("Error storing duplicate outlet criteria:", err.Error())
				return err
			}

			// Duplicate selected outlets - check for nil slice
			if outletCriteriaList[0].SelectedOutlets != nil && len(outletCriteriaList[0].SelectedOutlets) > 0 {
				var outletsSelected []model.PromotionOutletsSelected
				for _, outlet := range outletCriteriaList[0].SelectedOutlets {
					var outletSelected model.PromotionOutletsSelected
					outletSelected.CustID = params.ParentCustId
					outletSelected.CriteriaID = outletCriteriaID
					outletSelected.OutletID = outlet.OutletID
					outletSelected.ID = ""
					outletsSelected = append(outletsSelected, outletSelected)
				}
				err = service.PromotionV2Repository.StoreOutletsSelected(txCtx, outletsSelected)
				if err != nil {
					log.Error("Error storing duplicate outlets selected:", err.Error())
					return err
				}
			}

			// Duplicate outlet attribute types - check for nil slice
			if outletCriteriaList[0].AttributeTypes != nil && len(outletCriteriaList[0].AttributeTypes) > 0 {
				var outletAttributeTypes []model.PromotionOutletAttributeType
				for _, outletType := range outletCriteriaList[0].AttributeTypes {
					var outletAttributeType model.PromotionOutletAttributeType
					outletAttributeType.CustID = params.ParentCustId
					outletAttributeType.CriteriaID = outletCriteriaID
					outletAttributeType.OutletTypeID = outletType.OutletTypeID
					outletAttributeType.ID = ""
					outletAttributeTypes = append(outletAttributeTypes, outletAttributeType)
				}
				err = service.PromotionV2Repository.StoreOutletAttributeType(txCtx, outletAttributeTypes)
				if err != nil {
					log.Error("Error storing duplicate outlet attribute types:", err.Error())
					return err
				}
			}

			// Duplicate outlet attribute groups - check for nil slice
			if outletCriteriaList[0].AttributeGroups != nil && len(outletCriteriaList[0].AttributeGroups) > 0 {
				var outletAttributeGroups []model.PromotionOutletAttributeGroup
				for _, outletGroup := range outletCriteriaList[0].AttributeGroups {
					var outletAttributeGroup model.PromotionOutletAttributeGroup
					outletAttributeGroup.CustID = params.ParentCustId
					outletAttributeGroup.CriteriaID = outletCriteriaID
					outletAttributeGroup.OutletGroupID = outletGroup.OutletGroupID
					outletAttributeGroup.ID = ""
					outletAttributeGroups = append(outletAttributeGroups, outletAttributeGroup)
				}
				err = service.PromotionV2Repository.StoreOutletAttributeGroup(txCtx, outletAttributeGroups)
				if err != nil {
					log.Error("Error storing duplicate outlet attribute groups:", err.Error())
					return err
				}
			}

			// Duplicate outlet attribute classes - check for nil slice
			if outletCriteriaList[0].AttributeClasses != nil && len(outletCriteriaList[0].AttributeClasses) > 0 {
				var outletAttributeClasses []model.PromotionOutletAttributeClass
				for _, outletClass := range outletCriteriaList[0].AttributeClasses {
					var outletAttributeClass model.PromotionOutletAttributeClass
					outletAttributeClass.CustID = params.ParentCustId
					outletAttributeClass.CriteriaID = outletCriteriaID
					outletAttributeClass.OutletClassID = outletClass.OutletClassID
					outletAttributeClass.ID = ""
					outletAttributeClasses = append(outletAttributeClasses, outletAttributeClass)
				}
				err = service.PromotionV2Repository.StoreOutletAttributeClass(txCtx, outletAttributeClasses)
				if err != nil {
					log.Error("Error storing duplicate outlet attribute classes:", err.Error())
					return err
				}
			}

			// Duplicate outlet attribute sales teams - check for nil slice
			if outletCriteriaList[0].AttributeSalesTeams != nil && len(outletCriteriaList[0].AttributeSalesTeams) > 0 {
				var outletAttributeSalesTeams []model.PromotionOutletAttributeSalesTeam
				for _, salesTeam := range outletCriteriaList[0].AttributeSalesTeams {
					var outletAttributeSalesTeam model.PromotionOutletAttributeSalesTeam
					outletAttributeSalesTeam.CustID = params.ParentCustId
					outletAttributeSalesTeam.CriteriaID = outletCriteriaID
					outletAttributeSalesTeam.SalesTeamID = salesTeam.SalesTeamID
					outletAttributeSalesTeam.ID = ""
					outletAttributeSalesTeams = append(outletAttributeSalesTeams, outletAttributeSalesTeam)
				}
				err = service.PromotionV2Repository.StoreOutletAttributeSalesTeam(txCtx, outletAttributeSalesTeams)
				if err != nil {
					log.Error("Error storing duplicate outlet attribute sales teams:", err.Error())
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return newPromoID, nil
}

// generateNextSequenceNumber generates the next available sequence number for promotion duplication
func (service *promotionServiceImpl) generateNextSequenceNumber(custID, basePromoID string) (string, error) {
	// Find all existing promotion IDs that start with the base name
	existingPromoIDs, err := service.PromotionV2Repository.FindPromoIDsByBaseName(custID, basePromoID)
	if err != nil {
		return "", err
	}

	// Create a map to track existing sequence numbers
	existingSequences := make(map[int]bool)

	// Parse existing promotion IDs to extract sequence numbers
	for _, promoID := range existingPromoIDs {
		if promoID == basePromoID {
			// This is the original promotion, skip it
			continue
		}

		// Check if this promotion ID follows the pattern: basePromoID-XXX
		if len(promoID) > len(basePromoID)+1 && promoID[:len(basePromoID)+1] == basePromoID+"-" {
			sequenceStr := promoID[len(basePromoID)+1:]
			// Try to parse the sequence number
			if sequence, parseErr := strconv.Atoi(sequenceStr); parseErr == nil {
				existingSequences[sequence] = true
			}
		}
	}

	// Find the next available sequence number starting from 1
	nextSequence := 1
	for existingSequences[nextSequence] {
		nextSequence++
	}

	// Format the sequence number with leading zeros (001, 002, 003, etc.)
	sequenceStr := fmt.Sprintf("%03d", nextSequence)
	return basePromoID + "-" + sequenceStr, nil
}

func (service *promotionServiceImpl) DetailV2ForUpdate(params entity.DetailPromotionParams) (response entity.PromotionV2, err error) {
	promo, err := service.PromotionV2Repository.FindByPromoID(params)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(promo, &response)
	if err != nil {
		return response, err
	}

	response.EffectiveFrom = promo.EffectiveFrom.Format(constant.YYYY_MM_DD)
	response.EffectiveTo = promo.EffectiveTo.Format(constant.YYYY_MM_DD)
	mapPromotionV2ExtendedFieldsToEntity(promo, &response)
	mapPromotionV2ResponseCustID(promo, &response)

	promoSlabs, err := service.PromotionV2Repository.FindPromoSlabsByPromoID(params)
	if err != nil {
		return response, err
	}
	for _, slab := range promoSlabs {
		var slabItem entity.PromoSlabItem
		err = structs.Automapper(slab, &slabItem)
		if err != nil {
			return response, err
		}
		slabItem.CustID = ""
		slabItem.PromoID = ""
		response.Slabs = append(response.Slabs, slabItem)
	}

	promoStratas, err := service.PromotionV2Repository.FindPromoStratasByPromoID(params)
	if err != nil {
		return response, err
	}
	for _, strata := range promoStratas {
		var strataItem entity.PromoStrataItem
		err = structs.Automapper(strata, &strataItem)
		if err != nil {
			return response, err
		}
		strataItem.CustID = ""
		strataItem.PromoID = ""
		response.Strata = append(response.Strata, strataItem)
	}

	promoProductCriterias, err := service.PromotionV2Repository.FindPromoProductCriteriasByPromoID(params)
	if err != nil {
		return response, err
	}
	for _, productCriteria := range promoProductCriterias {
		var productCriteriaItem entity.PromoProductCriteria
		err = structs.Automapper(productCriteria, &productCriteriaItem)
		if err != nil {
			return response, err
		}
		productCriteriaItem.CustID = ""
		productCriteriaItem.PromoID = ""
		response.ProductCriteria = append(response.ProductCriteria, productCriteriaItem)
	}

	rewardProducts, err := service.PromotionV2Repository.FindPromoRewardProductsByPromoID(params)
	if err != nil {
		return response, err
	}
	for _, rewardProduct := range rewardProducts {
		var rewardProductItem entity.PromotionRewardProduct
		err = structs.Automapper(rewardProduct, &rewardProductItem)
		if err != nil {
			return response, err
		}
		rewardProductItem.CustID = ""
		rewardProductItem.PromoID = ""
		response.RewardProducts = append(response.RewardProducts, rewardProductItem)
	}

	coverageDistributors, err := service.PromotionV2Repository.FindCoverageDistributorsByPromoID(params)
	if err != nil {
		return response, err
	}
	response.CoverageDistributors = make([]entity.PromoCoverageDistributor, 0, len(coverageDistributors))
	for _, coverageDistributor := range coverageDistributors {
		response.CoverageDistributors = append(response.CoverageDistributors, mapPromoCoverageDistributorFromModel(coverageDistributor))
	}

	outletCriteriaList, err := service.PromotionV2Repository.FindOutletCriteriaWithPreloads(params)
	if err != nil {
		return response, err
	}

	// Convert outlet criteria list to entity format
	if len(outletCriteriaList) > 0 {
		var outletCriteriaItem entity.PromoOutletCriteria
		err = structs.Automapper(outletCriteriaList[0], &outletCriteriaItem)
		if err != nil {
			return response, err
		}
		outletCriteriaItem.CustID = ""
		outletCriteriaItem.PromoID = ""
		outletCriteriaItem.SelectedOutlets = []entity.PromoOutletSelected{}
		for _, outlet := range outletCriteriaList[0].SelectedOutlets {
			outletCriteriaItem.SelectedOutlets = append(outletCriteriaItem.SelectedOutlets, mapPromoOutletSelectedFromModel(outlet))
		}
		outletCriteriaItem.SelectedOutletTypes = []entity.PromoOutletTypeSelected{}
		for _, outletType := range outletCriteriaList[0].AttributeTypes {
			var outletTypeSelected entity.PromoOutletTypeSelected
			err = structs.Automapper(outletType, &outletTypeSelected)
			if err != nil {
				return response, err
			}
			outletCriteriaItem.SelectedOutletTypes = append(outletCriteriaItem.SelectedOutletTypes, outletTypeSelected)
		}
		outletCriteriaItem.SelectedOutletGroups = []entity.PromoOutletGroupSelected{}
		for _, outletGroup := range outletCriteriaList[0].AttributeGroups {
			var outletGroupSelected entity.PromoOutletGroupSelected
			err = structs.Automapper(outletGroup, &outletGroupSelected)
			if err != nil {
				return response, err
			}
			outletCriteriaItem.SelectedOutletGroups = append(outletCriteriaItem.SelectedOutletGroups, outletGroupSelected)
		}
		outletCriteriaItem.SelectedOutletClasses = []entity.PromoOutletClassSelected{}
		for _, outletClass := range outletCriteriaList[0].AttributeClasses {
			var outletClassSelected entity.PromoOutletClassSelected
			err = structs.Automapper(outletClass, &outletClassSelected)
			if err != nil {
				return response, err
			}
			outletCriteriaItem.SelectedOutletClasses = append(outletCriteriaItem.SelectedOutletClasses, outletClassSelected)
		}
		outletCriteriaItem.SelectedSalesTeams = []entity.PromoOutletSalesTeamSelected{}
		for _, salesTeam := range outletCriteriaList[0].AttributeSalesTeams {
			outletCriteriaItem.SelectedSalesTeams = append(outletCriteriaItem.SelectedSalesTeams, mapPromoOutletSalesTeamFromModel(salesTeam))
		}
		response.OutletCriteria = outletCriteriaItem
	}

	return response, nil
}

func (service *promotionServiceImpl) ConsultV2(req entity.ConsultPromoV2Req) (responses []entity.ConsultPromoResp, err error) {
	responses = make([]entity.ConsultPromoResp, 0)

	requestedPromoIDs, err := entity.NormalizePromoIDList(req.PromoIDs)
	if err != nil {
		return nil, err
	}

	// Phase 1: Initial Quantity Conversion
	log.Info("Phase 1: Initial Quantity Conversion >>> ", structs.StructToJson(req.Details))
	for index := range req.Details {
		detail := &req.Details[index]
		var detailConversion entity.PromoConversionReq
		detailConversion.CustID = req.CustID
		detailConversion.ProductID = int64(detail.ProID)
		detailConversion.Qty1 = int64(detail.Qty1)
		detailConversion.Qty2 = int64(detail.Qty2)
		detailConversion.Qty3 = int64(detail.Qty3)

		conversionResult, err := service.PromoConversion(detailConversion, req.CustID)
		if err != nil && req.CustID != req.ParentCustID {
			conversionResult, err = service.PromoConversion(detailConversion, req.ParentCustID)
		}
		if err != nil {
			log.Error(fmt.Sprintf("Error converting quantities for product %d: %v", detail.ProID, err))
			continue
		}

		// Update detail with converted quantities
		detail.Qty1 = float64(conversionResult.Qty1)
		detail.Qty2 = float64(conversionResult.Qty2)
		detail.Qty3 = float64(conversionResult.Qty3)
		detail.Total = float64(conversionResult.TotalQty)
	}
	log.Info("Phase 1: After Quantity Conversion >>> ", structs.StructToJson(req.Details))

	// Phase 2: Outlet & Salesman Validation
	log.Info("Phase 2: Outlet & Salesman Validation")
	outlet, err := service.PromotionV2Repository.FindOutletByID(int64(req.OutletID), req.CustID)
	if err != nil {
		log.Error(err.Error())
		return responses, fmt.Errorf("Outlet ID: %d not found", req.OutletID)
	}
	// log.Info("Outlet >>> ", structs.StructToJson(outlet))

	salesman, err := service.PromotionV2Repository.FindSalesmanByID(int64(req.SalesmanID), req.CustID)
	if err != nil {
		log.Error(err.Error())
		return responses, fmt.Errorf("Salesman ID: %d not found", req.SalesmanID)
	}
	// log.Info("Salesman >>> ", structs.StructToJson(salesman))

	warehouse, err := service.PromotionV2Repository.FindWarehouseByID(int64(req.WhID), req.CustID)
	if err != nil {
		log.Error(err.Error())
		return responses, fmt.Errorf("Warehouse ID: %d not found", req.WhID)
	}
	log.Info("Warehouse >>> ", structs.StructToJson(warehouse))

	// Phase 3: Build Attribute Validation Criteria
	// In v2, we don't build a map like v1, but we use outlet criteria directly

	// Phase 4: Find & Validate Promotions by Outlet Criteria
	log.Info("Phase 4: Find & Validate Promotions by Outlet Criteria")
	promotions, err := service.PromotionV2Repository.FindActivePromotionsByOutletCriteria(req, outlet, salesman)
	if err != nil {
		log.Error("Error finding active promotions:", err.Error())
		return responses, err
	}

	if len(promotions) == 0 {
		return responses, nil
	}

	promoByID := make(map[string]model.PromotionV2, len(promotions))
	for _, promo := range promotions {
		promoByID[promo.PromoID] = promo
	}

	// Extract promo IDs
	var promoIDs []string
	for _, promo := range promotions {
		promoIDs = append(promoIDs, promo.PromoID)
	}

	// Phase 5 & 6: Find and validate product criteria
	log.Info("Phase 5 & 6: Find and validate product criteria >>> ", structs.StructToJson(promoIDs))
	productCriterias, err := service.PromotionV2Repository.FindProductCriteriasByPromoIDs(promoIDs, req.ParentCustID)
	if err != nil {
		log.Error("Error finding product criterias:", err.Error())
		return responses, err
	}

	// Validate product criteria against request details.
	// Option A: For "reward per order" (fixed_value/percentage per_order), reward is divided by the number of products
	// that match promotion product criteria AND are in the purchase (len(validatedPromoProductGroups[promoID])).
	// Ensure promotion product criteria includes all products that should share the reward (e.g. A,B,C in criteria → buy A,B,C → reward/3 each; buy A,B only → reward/2 each).
	validatedPromoProductGroups := make(map[string]map[int]*entity.ConPromoV2Det)
	subTotalValidatedPromoProductGroups := make(map[string]int64)
	var validatedPromoList []string

	// Group product criteria by promo ID
	productCriteriaByPromo := make(map[string][]model.PromotionProductCriteria)
	for _, pc := range productCriterias {
		productCriteriaByPromo[pc.PromoID] = append(productCriteriaByPromo[pc.PromoID], pc)
	}

	// Validate each promotion's product criteria
	for _, promo := range promotions {
		log.Info("PromoID >>> ", promo.PromoID)
		validatedPromoProductGroups[promo.PromoID] = make(map[int]*entity.ConPromoV2Det)
		subTotalValidatedPromoProductGroups[promo.PromoID] = 0
		isPromoProductCriteriaValid := false
		promoCriterias := productCriteriaByPromo[promo.PromoID]

		// If no product criteria, all products are eligible
		if len(promoCriterias) == 0 {
			log.Info("Validate each promotion's product criteria > if:")
			isPromoProductCriteriaValid = true
			for index, detail := range req.Details {
				validatedPromoProductGroups[promo.PromoID][detail.ProID] = &req.Details[index]
				subTotalValidatedPromoProductGroups[promo.PromoID] += int64(detail.GrossValue)
			}
		} else {
			log.Info("Validate each promotion's product criteria > else:")
			hasMandatoryCriteria := false
			for _, productCriteria := range promoCriterias {
				if productCriteria.Mandatory {
					hasMandatoryCriteria = true
					break
				}
			}

			// Check if all mandatory products are present
			for _, productCriteria := range promoCriterias {
				log.Info("mandatory productCriteria > ", structs.StructToJson(productCriteria))
				if !productCriteria.Mandatory {
					continue
				}

				productFound := false
				for index, detail := range req.Details {
					log.Info("Validate each promotion's product criteria > for > detail: ", structs.StructToJson(detail))
					if int64(detail.ProID) == productCriteria.ProID {
						log.Info("detail ProID > ", detail.ProID)
						// Calculate buy value
						buyValue := buyValueForCriteria(&detail, productCriteria)
						minBuyValue := minBuyValueForCriteria(productCriteria)

						log.Info("buyValue >= minBuyValue > ", buyValue, " >= ", minBuyValue)
						if buyValue >= minBuyValue {
							productFound = true
							isPromoProductCriteriaValid = true
							validatedPromoProductGroups[promo.PromoID][detail.ProID] = &req.Details[index]
							subTotalValidatedPromoProductGroups[promo.PromoID] += int64(detail.GrossValue)
						}
					}
				}
				log.Info("productFound > ", productFound)

				if !productFound {
					log.Info("Mandatory product not found or doesn't meet minimum")
					// Mandatory product not found or doesn't meet minimum
					isPromoProductCriteriaValid = false
					break
				}
			}

			// Non-mandatory products only count when there is no mandatory criteria,
			// or all mandatory products in the order already passed validation.
			log.Info("isPromoProductCriteriaValid > ", isPromoProductCriteriaValid)
			if !hasMandatoryCriteria || isPromoProductCriteriaValid {
				for _, productCriteria := range promoCriterias {
					log.Info("non-mandatory productCriteria > ", structs.StructToJson(productCriteria))
					if productCriteria.Mandatory {
						continue
					}

					for index, detail := range req.Details {
						if int64(detail.ProID) == productCriteria.ProID {
							buyValue := buyValueForCriteria(&detail, productCriteria)
							minBuyValue := minBuyValueForCriteria(productCriteria)

							log.Info("buyValue >= minBuyValue > ", buyValue, " >= ", minBuyValue)
							if buyValue >= minBuyValue {
								isPromoProductCriteriaValid = true
								validatedPromoProductGroups[promo.PromoID][detail.ProID] = &req.Details[index]
								subTotalValidatedPromoProductGroups[promo.PromoID] += int64(detail.GrossValue)
							}
						}
					}
				}
			}
		}

		if isPromoProductCriteriaValid && len(validatedPromoProductGroups[promo.PromoID]) > 0 {
			log.Info("validatedPromoList append - PromoID > ", promo.PromoID)
			validatedPromoList = append(validatedPromoList, promo.PromoID)
		} else {
			log.Info("validatedPromoList delete - PromoID > ", promo.PromoID)
			delete(validatedPromoProductGroups, promo.PromoID)
			delete(subTotalValidatedPromoProductGroups, promo.PromoID)
		}
	}

	log.Info("len(validatedPromoList) >>> ", len(validatedPromoList))
	if len(validatedPromoList) == 0 {
		return responses, nil
	}

	log.Info("validatedPromoList: ", structs.StructToJson(validatedPromoList))
	log.Info("validatedPromoProductGroups: ", structs.StructToJson(validatedPromoProductGroups))
	log.Info("subTotalValidatedPromoProductGroups: ", structs.StructToJson(subTotalValidatedPromoProductGroups))

	// Phase 7: Validate Slab Rules
	log.Info("Phase 7: Validate Slab Rules")
	slabs, err := service.PromotionV2Repository.FindSlabsByPromoIDs(validatedPromoList, req.ParentCustID)
	if err != nil {
		log.Error("Error finding slabs:", err.Error())
		return responses, err
	}

	log.Info("slabs: ", structs.StructToJson(slabs))

	validatedSlabs := make(map[string]model.PromotionV2Slabs)

	// Group slabs by promo ID and validate
	log.Info("Group slabs by promo ID and validate")
	for _, slab := range slabs {
		slabRuleValue := float64(subTotalValidatedPromoProductGroups[slab.PromoID])
		if slab.RuleType == model.RuleTypeQuantity {
			slabRuleValue = 0
			for _, detail := range validatedPromoProductGroups[slab.PromoID] {
				slabRuleValue += detailQtyByUom(detail, slab.RuleUom)
			}
		} else if slab.RuleType == model.RuleTypeValue {
			slabRuleValue = 0
			for _, detail := range validatedPromoProductGroups[slab.PromoID] {
				if detail == nil {
					continue
				}
				slabRuleValue += float64(detail.GrossValue)
			}
		}

		// Check if slab rule is valid
		rangeFrom := 0.0
		if slab.RangeFrom != nil {
			rangeFrom = *slab.RangeFrom
		}
		isMultiplied := false
		if promo, exists := promoByID[slab.PromoID]; exists && promo.SlabMultiplied != nil && *promo.SlabMultiplied {
			isMultiplied = true
		}

		log.Info(
			"Validating slab > promoID:", slab.PromoID,
			"| slabRuleValue:", slabRuleValue,
			"| range:", rangeFrom, "-", slab.RangeTo,
			"| isMultiplied:", isMultiplied,
		)
		if isMultiplied || (slabRuleValue >= rangeFrom && slabRuleValue <= slab.RangeTo) {
			log.Info("validatedSlabs append - PromoID > ", structs.StructToJson(slab))
			validatedSlabs[slab.PromoID] = slab
		}
	}

	log.Info("validatedSlabs: ", structs.StructToJson(validatedSlabs))

	// phase 7b: validate strata
	stratas, err := service.PromotionV2Repository.
		FindStratasByPromoIDs(validatedPromoList, req.ParentCustID)
	if err != nil {
		log.Error("Error finding strata:", err.Error())
		return responses, err
	}

	// log.Info("stratas: ", structs.StructToJson(stratas))
	validatedStrata := make(map[string]model.PromotionV2Strata)                    // single strata per promo (non-sequential)
	validatedStrataList := make(map[string][]model.PromotionV2Strata)              // all stratas per promo (sequential)
	validatedStrataListNonSequential := make(map[string][]model.PromotionV2Strata) // all stratas per promo (non-sequential, for response + percentage calc)

	for promoID, details := range validatedPromoProductGroups {
		promo, exists := promoByID[promoID]
		if !exists {
			log.Warn("Promotion not found in promoMap for promoID:", promoID)
			continue
		}

		// Filter stratas for this promo
		promoStratas := make([]model.PromotionV2Strata, 0)
		for _, s := range stratas {
			if s.PromoID == promoID {
				promoStratas = append(promoStratas, s)
			}
		}

		if len(promoStratas) == 0 {
			continue
		}

		// Determine if sequential calculation is enabled
		isSequential := promo.StrataSequential != nil && *promo.StrataSequential

		if isSequential {
			// Sequential: sort ASC by ordinal (1..5) for display and sequential calculation
			sort.SliceStable(promoStratas, func(i, j int) bool {
				return promoStratas[i].Ordinal < promoStratas[j].Ordinal
			})
			ruleValue := calculateStrataRuleValue(promoStratas[0], details, req.Details)
			// Find highest ordinal stratum that contains ruleValue (tier K we qualify for)
			highestOrdinal := 0
			for _, s := range promoStratas {
				if ruleValue >= s.RangeFrom && ruleValue <= s.RangeTo && s.Ordinal > highestOrdinal {
					highestOrdinal = s.Ordinal
				}
			}
			if highestOrdinal > 0 {
				// Store only strata 1..K (applicable tiers), not all 5
				applicableStrata := make([]model.PromotionV2Strata, 0, highestOrdinal)
				for _, s := range promoStratas {
					if s.Ordinal <= highestOrdinal {
						applicableStrata = append(applicableStrata, s)
					}
				}
				validatedStrataList[promoID] = applicableStrata
			}
		} else {
			// Non-sequential: sort ASC to get lowest ordinal first (lowest tier)
			sort.SliceStable(promoStratas, func(i, j int) bool {
				return promoStratas[i].Ordinal < promoStratas[j].Ordinal
			})
			ruleValue := calculateStrataRuleValue(promoStratas[0], details, req.Details)
			matchedOrdinal := resolveNonSequentialStrataOrdinal(promoStratas, ruleValue)
			for _, strata := range promoStratas {
				log.Info(
					"Validating strata > promoID:", promoID,
					"| ruleValue:", ruleValue,
					"| range:", strata.RangeFrom, "-", strata.RangeTo,
					"| ordinal:", strata.Ordinal,
					"| sequential:", isSequential,
				)
			}
			if matchedOrdinal > 0 {
				var matchedStrata model.PromotionV2Strata
				applicableStrata := make([]model.PromotionV2Strata, 0, matchedOrdinal)
				for _, s := range promoStratas {
					if s.Ordinal <= matchedOrdinal {
						applicableStrata = append(applicableStrata, s)
					}
					if s.Ordinal == matchedOrdinal {
						matchedStrata = s
					}
				}
				validatedStrata[promoID] = matchedStrata
				validatedStrataListNonSequential[promoID] = applicableStrata
			}
		}
	}

	log.Info("validatedStrata:", structs.StructToJson(validatedStrata))
	log.Info("validatedStrataList (sequential):", structs.StructToJson(validatedStrataList))

	// Combine promo IDs from validated slabs, single strata, and strata list (sequential)
	combinedPromoIDs := map[string]bool{}

	for promoID := range validatedSlabs {
		combinedPromoIDs[promoID] = true
	}
	for promoID := range validatedStrata {
		combinedPromoIDs[promoID] = true
	}
	for promoID := range validatedStrataList {
		combinedPromoIDs[promoID] = true
	}

	// Total gross value of entire order (all details) for response
	orderTotalGross := 0.0
	for _, d := range req.Details {
		orderTotalGross += float64(d.GrossValue)
	}

	// Phase 8: Calculate Rewards
	log.Info("Phase 8: Calculate Rewards")
	insufficientStockPromoIDs := make([]string, 0)
promoLoop:
	for promoID := range combinedPromoIDs {
		if len(requestedPromoIDs) > 0 && !promoIDInList(promoID, requestedPromoIDs) {
			continue
		}
		slab, hasSlab := validatedSlabs[promoID]
		strata, hasStrata := validatedStrata[promoID]
		strataList, hasStrataList := validatedStrataList[promoID]
		strataListNonSeq, hasStrataListNonSeq := validatedStrataListNonSequential[promoID]

		// For reward config: use single strata or first of strata list
		strataForCfg := strata
		hasStrataForCfg := hasStrata
		if hasStrataList && len(strataList) > 0 {
			strataForCfg = strataList[0]
			hasStrataForCfg = true
		}

		rewardCfg, ok := rewardConfigFromSlabStrata(slab, hasSlab, strataForCfg, hasStrataForCfg)
		if !ok {
			continue
		}

		rewardType := rewardCfg.rewardType
		ruleType := rewardCfg.ruleType
		rewardValuePtr := rewardCfg.rewardValuePtr
		rewardUom := rewardCfg.rewardUom
		perScope := rewardCfg.perScope

		log.Info("validatedSlabs > promoID: ", promoID, " > ", structs.StructToJson(slab))
		promo, exists := promoByID[promoID]
		if !exists {
			continue
		}

		isSequential := promo.StrataSequential != nil && *promo.StrataSequential
		useStrataSequentialPercentage := hasStrataList && isSequential && rewardType == model.RewardTypePercentage
		useStrataNonSequentialPercentage := hasStrata && !isSequential && rewardType == model.RewardTypePercentage && hasStrataListNonSeq && len(strataListNonSeq) > 0
		useStrataNonSequentialFixedValue := hasStrata && !isSequential && rewardType == model.RewardTypeFixedValue && hasStrataListNonSeq && len(strataListNonSeq) > 0
		useStrataNonSequentialProduct := hasStrata && !isSequential && rewardType == model.RewardTypeProduct && hasStrataListNonSeq && len(strataListNonSeq) > 0

		// log.Info("promo > ", structs.StructToJson(promo))
		var response entity.ConsultPromoResp
		response.PromoID = promoID
		response.PromoDesc = promo.PromoDesc
		if promo.PromoType == model.PromotionTypeSlab && perScope != nil {
			response.SlabPerScope = *perScope
		}
		if hasSlab {
			response.SlabID = slab.ID
			if slab.Description != nil {
				response.SlabDesc = *slab.Description
			}
		}

		// Set slab_* only for promo_type slab; leave empty for strata
		if promo.PromoType == model.PromotionTypeSlab {
			response.SlabRewardType = entity.RewardType(rewardType)
			response.SlabRuleType = entity.RuleType(rewardCfg.ruleType)
			if rewardUom != nil {
				response.SlabRewardUom = entity.UOMType(*rewardUom)
			}
			if rewardCfg.ruleUom != nil {
				response.SlabRuleUom = entity.UOMType(*rewardCfg.ruleUom)
			}
		}

		// total_gross_value: default = total order (all request details); for strata response overwritten below with total of eligible products only
		response.TotalGrossValue = orderTotalGross

		// Strata-only response: fill strata_* when we have sequential strata (prefer over slab)
		if useStrataSequentialPercentage {
			response.StrataID = make([]string, 0, len(strataList))
			response.StrataDesc = make([]string, 0, len(strataList))
			response.StrataReward = make([]float64, 0, len(strataList))
			for _, s := range strataList {
				response.StrataID = append(response.StrataID, s.ID)
				if s.Description != nil {
					response.StrataDesc = append(response.StrataDesc, *s.Description)
				} else {
					response.StrataDesc = append(response.StrataDesc, "")
				}
				rv := 0.0
				if s.RewardValue != nil {
					rv = *s.RewardValue
				}
				response.StrataReward = append(response.StrataReward, rv)
			}
			response.StrataRuleType = entity.RuleType(strataList[0].RuleType)
			if strataList[0].RuleUom != nil {
				response.StrataRuleUom = entity.UOMType(*strataList[0].RuleUom)
			}
			if strataList[0].RewardUom != nil {
				response.StrataRewardUom = entity.UOMType(*strataList[0].RewardUom)
			}
			response.StrataRewardType = entity.RewardType(strataList[0].RewardType)
			if strataList[0].PerScope != nil {
				response.StrataPerScope = *strataList[0].PerScope
			}
		}
		// Strata response for non-sequential percentage: fill strata_* from all stratas
		if useStrataNonSequentialPercentage {
			response.StrataID = make([]string, 0, len(strataListNonSeq))
			response.StrataDesc = make([]string, 0, len(strataListNonSeq))
			response.StrataReward = make([]float64, 0, len(strataListNonSeq))
			for _, s := range strataListNonSeq {
				response.StrataID = append(response.StrataID, s.ID)
				if s.Description != nil {
					response.StrataDesc = append(response.StrataDesc, *s.Description)
				} else {
					response.StrataDesc = append(response.StrataDesc, "")
				}
				rv := 0.0
				if s.RewardValue != nil {
					rv = *s.RewardValue
				}
				response.StrataReward = append(response.StrataReward, rv)
			}
			response.StrataRuleType = entity.RuleType(strataListNonSeq[0].RuleType)
			if strataListNonSeq[0].RuleUom != nil {
				response.StrataRuleUom = entity.UOMType(*strataListNonSeq[0].RuleUom)
			}
			if strataListNonSeq[0].RewardUom != nil {
				response.StrataRewardUom = entity.UOMType(*strataListNonSeq[0].RewardUom)
			}
			response.StrataRewardType = entity.RewardType(strataListNonSeq[0].RewardType)
			if strataListNonSeq[0].PerScope != nil {
				response.StrataPerScope = *strataListNonSeq[0].PerScope
			}
		}
		// Strata response for non-sequential fixed_value: fill strata_* from all stratas
		if useStrataNonSequentialFixedValue {
			response.StrataID = make([]string, 0, len(strataListNonSeq))
			response.StrataDesc = make([]string, 0, len(strataListNonSeq))
			response.StrataReward = make([]float64, 0, len(strataListNonSeq))
			for _, s := range strataListNonSeq {
				response.StrataID = append(response.StrataID, s.ID)
				if s.Description != nil {
					response.StrataDesc = append(response.StrataDesc, *s.Description)
				} else {
					response.StrataDesc = append(response.StrataDesc, "")
				}
				rv := 0.0
				if s.RewardValue != nil {
					rv = *s.RewardValue
				}
				response.StrataReward = append(response.StrataReward, rv)
			}
			response.StrataRuleType = entity.RuleType(strataListNonSeq[0].RuleType)
			if strataListNonSeq[0].RuleUom != nil {
				response.StrataRuleUom = entity.UOMType(*strataListNonSeq[0].RuleUom)
			}
			if strataListNonSeq[0].RewardUom != nil {
				response.StrataRewardUom = entity.UOMType(*strataListNonSeq[0].RewardUom)
			}
			response.StrataRewardType = entity.RewardType(strataListNonSeq[0].RewardType)
			if strataListNonSeq[0].PerScope != nil {
				response.StrataPerScope = *strataListNonSeq[0].PerScope
			}
		}
		// Strata response for non-sequential product: fill strata_* from all stratas
		if useStrataNonSequentialProduct {
			response.StrataID = make([]string, 0, len(strataListNonSeq))
			response.StrataDesc = make([]string, 0, len(strataListNonSeq))
			response.StrataReward = make([]float64, 0, len(strataListNonSeq))
			for _, s := range strataListNonSeq {
				response.StrataID = append(response.StrataID, s.ID)
				if s.Description != nil {
					response.StrataDesc = append(response.StrataDesc, *s.Description)
				} else {
					response.StrataDesc = append(response.StrataDesc, "")
				}
				rv := 0.0
				if s.RewardValue != nil {
					rv = *s.RewardValue
				}
				response.StrataReward = append(response.StrataReward, rv)
			}
			response.StrataRuleType = entity.RuleType(strataListNonSeq[0].RuleType)
			if strataListNonSeq[0].RuleUom != nil {
				response.StrataRuleUom = entity.UOMType(*strataListNonSeq[0].RuleUom)
			}
			if strataListNonSeq[0].RewardUom != nil {
				response.StrataRewardUom = entity.UOMType(*strataListNonSeq[0].RewardUom)
			}
			response.StrataRewardType = entity.RewardType(strataListNonSeq[0].RewardType)
			if strataListNonSeq[0].PerScope != nil {
				response.StrataPerScope = *strataListNonSeq[0].PerScope
			}
		}

		// Calculate reward (single value; not used for strata sequential percentage)
		rewardValue := 0.0
		if rewardValuePtr != nil {
			rewardValue = *rewardValuePtr
			log.Info("rewardValue > ", rewardValue)
		}

		// Calculate slab reward based on type (leave 0 for strata-only response)
		response.SlabReward = rewardValue
		if useStrataSequentialPercentage || useStrataNonSequentialPercentage || useStrataNonSequentialFixedValue || useStrataNonSequentialProduct {
			response.SlabReward = 0
		}

		// Calculate total gross value for this promo first
		totalGrossValueForPromo := 0.0
		for proID := range validatedPromoProductGroups[promoID] {
			totalGrossValueForPromo += float64(validatedPromoProductGroups[promoID][proID].GrossValue)
		}
		// Option A: For strata/slab per_order reward, total_gross_value = sum of eligible products (products that match promotion product criteria and are in the purchase)
		if useStrataSequentialPercentage || useStrataNonSequentialPercentage || useStrataNonSequentialFixedValue || useStrataNonSequentialProduct {
			response.TotalGrossValue = totalGrossValueForPromo
		}

		// Check if special multiplied calculation is needed
		isMultiplied := promo.SlabMultiplied != nil && *promo.SlabMultiplied
		isFixedValue := rewardType == model.RewardTypeFixedValue &&
			perScope != nil && (*perScope == string(model.PerScopeOrder) || *perScope == string(model.PerScopeProduct))
		useMultipliedCalculation := isMultiplied && isFixedValue && hasSlab && slab.RangeTo > 0

		var multipliedRewardPerProduct, multiplier, totalReward, totalQtyByUom float64 = 0, 0, 0, 0
		if useMultipliedCalculation {
			// Calculate multiplier: round down(total_gross_value / slab.range_to) for value rule; for quantity use total qty (by rule UOM) / range_to
			if ruleType == model.RuleTypeQuantity {
				// Quantity rule: total reward = slab_reward × multiplier (multiplier = floor(totalQtyByUom / rangeTo)); per_order = distribute totalReward equally by product count
				for _, d := range validatedPromoProductGroups[promoID] {
					totalQtyByUom += detailQtyByUom(d, rewardCfg.ruleUom)
				}
				if slab.RangeTo > 0 {
					multiplier = math.Floor(totalQtyByUom / slab.RangeTo)
				} else {
					multiplier = 1
				}
				totalReward = multiplier * rewardValue
				if *perScope == string(model.PerScopeOrder) {
					productCount := len(validatedPromoProductGroups[promoID])
					if productCount > 0 {
						multipliedRewardPerProduct = totalReward / float64(productCount)
					}
				}
				log.Infof(
					"[multiplied calculation quantity] promoID=%s totalQtyByUom=%.2f rangeTo=%.2f multiplier=%.0f totalReward=%.2f rewardPerProduct=%.2f",
					promoID, totalQtyByUom, slab.RangeTo, multiplier, totalReward, multipliedRewardPerProduct,
				)
			} else {
				multiplier = math.Floor(totalGrossValueForPromo / slab.RangeTo)
				totalReward = multiplier * rewardValue
				// Per-product reward: totalReward / number of products (for per_order value rule)
				if *perScope == string(model.PerScopeOrder) {
					productCount := len(validatedPromoProductGroups[promoID])
					if productCount > 0 {
						multipliedRewardPerProduct = totalReward / float64(productCount)
					}
					log.Infof(
						"[multiplied calculation value] promoID=%s totalGrossValue=%.2f rangeTo=%.2f multiplier=%.0f totalReward=%.2f productCount=%d rewardPerProduct=%.2f",
						promoID, totalGrossValueForPromo, slab.RangeTo, multiplier, totalReward, productCount, multipliedRewardPerProduct,
					)
				}
			}
		}

		log.Info("validatedPromoProductGroups > ", structs.StructToJson(validatedPromoProductGroups[promoID]))
		// Add products to response in request order (req.Details) so products_eligible and reward_percentage match expected order
		proIDs := make([]int, 0, len(validatedPromoProductGroups[promoID]))
		seen := make(map[int]bool)
		for _, detail := range req.Details {
			if validatedPromoProductGroups[promoID][detail.ProID] != nil && !seen[detail.ProID] {
				proIDs = append(proIDs, detail.ProID)
				seen[detail.ProID] = true
			}
		}
		for _, proID := range proIDs {
			log.Info("proID > ", proID)
			response.ProductsEligible = append(response.ProductsEligible, proID)

			log.Info("response.SlabReward > ", response.SlabReward)

			if useStrataSequentialPercentage {
				// Sequential percentage: promo1 = gross*P1%, promo2 = (gross-promo1)*P2%, ... net = gross - (promo1+...+promo5)
				rewardPercentage := entity.PromoRewardPercentage{}
				rewardPercentage.ProID = proID
				grossVal := float64(validatedPromoProductGroups[promoID][proID].GrossValue)
				rewardPercentage.GrossValue = grossVal
				remaining := grossVal
				for i, s := range strataList {
					pct := 0.0
					if s.RewardValue != nil {
						pct = *s.RewardValue
					}
					pr := math.Round(remaining * pct / 100.0)
					switch i {
					case 0:
						rewardPercentage.Promo1 = pr
					case 1:
						rewardPercentage.Promo2 = pr
					case 2:
						rewardPercentage.Promo3 = pr
					case 3:
						rewardPercentage.Promo4 = pr
					case 4:
						rewardPercentage.Promo5 = pr
					}
					remaining -= pr
				}
				rewardPercentage.NetValue = rewardPercentage.GrossValue - (rewardPercentage.Promo1 + rewardPercentage.Promo2 + rewardPercentage.Promo3 + rewardPercentage.Promo4 + rewardPercentage.Promo5)
				response.RewardPercentage = append(response.RewardPercentage, rewardPercentage)
				continue
			}

			if useStrataNonSequentialPercentage {
				// Non-sequential percentage: promo_i = gross_value * strata_reward[i]% for i=1..5, net = gross - (promo1+...+promo5)
				rewardPercentage := entity.PromoRewardPercentage{}
				rewardPercentage.ProID = proID
				grossVal := float64(validatedPromoProductGroups[promoID][proID].GrossValue)
				rewardPercentage.GrossValue = grossVal
				totalPromo := 0.0
				for i, s := range strataListNonSeq {
					pct := 0.0
					if s.RewardValue != nil {
						pct = *s.RewardValue
					}
					pr := math.Round(grossVal * pct / 100.0)
					totalPromo += pr
					switch i {
					case 0:
						rewardPercentage.Promo1 = pr
					case 1:
						rewardPercentage.Promo2 = pr
					case 2:
						rewardPercentage.Promo3 = pr
					case 3:
						rewardPercentage.Promo4 = pr
					case 4:
						rewardPercentage.Promo5 = pr
					}
				}
				rewardPercentage.NetValue = rewardPercentage.GrossValue - totalPromo
				response.RewardPercentage = append(response.RewardPercentage, rewardPercentage)
				continue
			}

			if useStrataNonSequentialFixedValue {
				// Non-sequential fixed_value: per_order = strata_reward[i]/productCount; per_product = qty * strata_reward[i]
				detail := validatedPromoProductGroups[promoID][proID]
				productCount := len(validatedPromoProductGroups[promoID])
				rv := entity.PromoRewardValue{}
				rv.ProID = proID
				rv.GrossValue = float64(detail.GrossValue)
				var promo1, promo2, promo3, promo4, promo5 float64
				if perScope != nil && *perScope == string(model.PerScopeOrder) {
					// per_order: reward per product = strata_reward[i] / productCount (same for all products)
					for i, s := range strataListNonSeq {
						strataReward := 0.0
						if s.RewardValue != nil {
							strataReward = *s.RewardValue
						}
						pr := math.Round(strataReward / float64(productCount))
						switch i {
						case 0:
							promo1 = pr
						case 1:
							promo2 = pr
						case 2:
							promo3 = pr
						case 3:
							promo4 = pr
						case 4:
							promo5 = pr
						}
					}
				} else {
					// per_product: promo_i = qty (by rule UOM) * strata_reward[i]
					qty := detailQtyByUom(detail, strataListNonSeq[0].RuleUom)
					for i, s := range strataListNonSeq {
						strataReward := 0.0
						if s.RewardValue != nil {
							strataReward = *s.RewardValue
						}
						pr := math.Round(qty * strataReward)
						switch i {
						case 0:
							promo1 = pr
						case 1:
							promo2 = pr
						case 2:
							promo3 = pr
						case 3:
							promo4 = pr
						case 4:
							promo5 = pr
						}
					}
				}
				rv.Promo1 = promo1
				rv.Promo2 = promo2
				rv.Promo3 = promo3
				rv.Promo4 = promo4
				rv.Promo5 = promo5
				rv.NetValue = rv.GrossValue - (promo1 + promo2 + promo3 + promo4 + promo5)
				response.RewardValue = append(response.RewardValue, rv)
				continue
			}

			if response.SlabReward == 0 {
				continue
			}

			// Calculate reward based on reward type
			switch rewardType {
			case model.RewardTypePercentage:
				log.Info("reward_type = percentage")
				rewardPercentage := entity.PromoRewardPercentage{}
				rewardPercentage.ProID = proID
				rewardPercentage.GrossValue = float64(validatedPromoProductGroups[promoID][proID].GrossValue)
				rewardValueCalculated := math.Round((rewardPercentage.GrossValue * rewardValue) / 100.0)

				log.Info("rewardPercentage.GrossValue > ", rewardPercentage.GrossValue)

				rewardPercentage.Promo1 = rewardValueCalculated
				rewardPercentage.NetValue = rewardPercentage.GrossValue - rewardValueCalculated

				response.RewardPercentage = append(response.RewardPercentage, rewardPercentage)
				log.Info("response > ", structs.StructToJson(response))

			case model.RewardTypeFixedValue:
				log.Info("reward_type = fixed_value")
				log.Info("validatedPromoProductGroups[promoID][proID] > ", structs.StructToJson(validatedPromoProductGroups[promoID][proID]))

				rewardValue := entity.PromoRewardValue{}
				rewardValue.ProID = proID
				rewardValue.GrossValue = float64(validatedPromoProductGroups[promoID][proID].GrossValue)

				log.Infof("rewardValue.GrossValue > %f", rewardValue.GrossValue)

				if useMultipliedCalculation {
					// Quantity rule + per_order: distribute totalReward equally (totalReward / productCount); works for smallest/middle/largest UOM
					if ruleType == model.RuleTypeQuantity && *perScope == string(model.PerScopeOrder) {
						rewardValue.Promo1 = multipliedRewardPerProduct
						log.Infof("Using multiplied calculation (quantity per_order equal split) > rewardValue.Promo1: %.2f", rewardValue.Promo1)
					} else if ruleType != model.RuleTypeQuantity && *perScope == string(model.PerScopeOrder) {
						// Value rule per_order: split total reward by product count
						log.Infof("if multipliedRewardPerProduct > %f", multipliedRewardPerProduct)
						rewardValue.Promo1 = multipliedRewardPerProduct
						log.Infof("Using multiplied calculation > rewardValue.Promo1: %f", rewardValue.Promo1)
					} else {
						// per_product: quantity rule = totalReward × qty (by rule UOM); value rule = totalReward × Qty3
						if ruleType == model.RuleTypeQuantity {
							qtyByUom := detailQtyByUom(validatedPromoProductGroups[promoID][proID], rewardCfg.ruleUom)
							rewardValue.Promo1 = totalReward * qtyByUom
							log.Infof("Using multiplied calculation (quantity per_product) > totalReward=%.0f qtyByUom=%.2f rewardValue.Promo1: %.2f", totalReward, qtyByUom, rewardValue.Promo1)
						} else {
							rewardValue.Promo1 = totalReward * float64(validatedPromoProductGroups[promoID][proID].Qty3)
							log.Infof("Using multiplied calculation (value per_product) > rewardValue.Promo1: %f", rewardValue.Promo1)
						}
					}
				} else {
					// Use existing perScopeRewardValue logic
					rewardValueCalculated := response.SlabReward
					log.Info("PerScopeOrder > ", len(validatedPromoProductGroups[promoID]))
					log.Info("PerScopeProduct > ruleType", structs.StructToJson(validatedPromoProductGroups[promoID][proID]), " > ", ruleType)
					rewardValue.Promo1 = perScopeRewardValue(
						rewardValueCalculated,
						perScope,
						ruleType,
						rewardCfg.ruleUom,
						validatedPromoProductGroups[promoID][proID],
						len(validatedPromoProductGroups[promoID]),
					)
				}
				rewardValue.NetValue = rewardValue.GrossValue - rewardValue.Promo1

				response.RewardValue = append(response.RewardValue, rewardValue)
			}

		}

		// Handle product rewards
		log.Info("slab.RewardType Product > ", rewardType)
		if rewardType == model.RewardTypeProduct {
			var rewardCtx model.RewardContext

			if hasStrataList && len(strataList) > 0 {
				rewardCtx = model.RewardContext{
					PromoID:   strataList[0].PromoID,
					RewardUom: strataList[0].RewardUom,
				}
			} else if hasStrata {
				rewardCtx = model.RewardContext{
					PromoID:   strata.PromoID,
					RewardUom: strata.RewardUom,
				}
			} else {
				rewardCtx = model.RewardContext{
					PromoID:   slab.PromoID,
					RewardUom: slab.RewardUom,
				}
			}

			log.Info("GetAllRewardProductFromStockV2 > promoID: ", promoID)
			rewards, _ := service.PromotionV2Repository.GetAllRewardProductFromStockV2(req, rewardCtx)

			log.Info("rewards: ", structs.StructToJson(rewards))

			multipliedValue := float64(1)
			if promo.SlabMultiplied != nil && *promo.SlabMultiplied {
				slabRuleValue := float64(subTotalValidatedPromoProductGroups[promoID])
				if slab.RuleType == model.RuleTypeQuantity {
					slabRuleValue = 0
					for _, detail := range validatedPromoProductGroups[promoID] {
						buyValue := 0.0
						if slab.RuleUom != nil {
							switch *slab.RuleUom {
							case model.UomTypeMiddle:
								buyValue = detail.Qty2
							case model.UomTypeSmallest:
								buyValue = detail.Total
							default:
								buyValue = detail.Qty3
							}
						}
						slabRuleValue += buyValue
					}
				}
				if slab.RangeTo > 0 {
					multipliedValue = math.Floor(slabRuleValue / slab.RangeTo)
				}
			}

			totalQtyReward := 0.0
			if useStrataNonSequentialProduct {
				for _, s := range strataListNonSeq {
					if s.RewardValue != nil {
						totalQtyReward += *s.RewardValue
					}
				}
			} else {
				log.Info("rewardValue > ", rewardValue)
				log.Info("multipliedValue > ", multipliedValue)
				if rewardValue > 0 {
					totalQtyReward = rewardValue * multipliedValue
				} else if multipliedValue > 0 {
					totalQtyReward = multipliedValue
				}
			}
			log.Info("totalQtyReward > ", totalQtyReward)

			if totalQtyReward > 0 {
				// Sum promo-available stock in reward UOM (warehouse stock minus order qty per reward product).
				var totalAvailableStock float64
				for i := range rewards {
					r := &rewards[i]
					stockInRewardUom, convErr := service.rewardProductStockInUom(r.QtyStock, int64(r.ConvUnit2), int64(r.ConvUnit3), rewardUom)
					if convErr != nil {
						log.Error("Error converting reward product for stock check:", convErr.Error())
						return nil, convErr
					}
					orderQty := orderQtyInRewardUom(req.Details, int(r.ProID), rewardUom)
					available := stockInRewardUom - orderQty
					if available < 0 {
						available = 0
					}
					totalAvailableStock += available
				}
				log.Info("totalAvailableStock (reward UOM, after order qty) > ", totalAvailableStock, " | totalQtyReward > ", totalQtyReward)

				if totalAvailableStock < totalQtyReward {
					log.Info("Insufficient stock for product reward, skipping promo from consult response")
					insufficientStockPromoIDs = append(insufficientStockPromoIDs, promoID)
					continue promoLoop
				}

				// Sufficient stock: allocate reward from stock
				remainingQtyReward := totalQtyReward
				for _, reward := range rewards {
					var rewardProduct entity.PromoRewardProductDet

					stockInRewardUom, convErr := service.rewardProductStockInUom(reward.QtyStock, int64(reward.ConvUnit2), int64(reward.ConvUnit3), rewardUom)
					if convErr != nil {
						log.Error("Error converting reward product:", convErr.Error())
						return nil, convErr
					}
					orderQty := orderQtyInRewardUom(req.Details, int(reward.ProID), rewardUom)
					availableStock := stockInRewardUom - orderQty
					if availableStock < 0 {
						availableStock = 0
					}
					log.Info("reward stock after order qty: pro_id=", reward.ProID, " available=", availableStock)

					qtyReward := remainingQtyReward
					if remainingQtyReward >= availableStock {
						qtyReward = availableStock
					}
					remainingQtyReward -= qtyReward

					if useStrataNonSequentialProduct {
						if len(response.RewardProduct) > 0 {
							if remainingQtyReward <= 0 {
								break
							}
							continue
						}

						var promoByOrdinal [5]float64
						var totalGross, sumQty1, sumQty2, sumQty3 float64
						for i, s := range strataListNonSeq {
							if i >= 5 || s.RewardValue == nil || *s.RewardValue <= 0 {
								continue
							}
							ordinalConversion := buildRewardProductConversionBody(req.CustID, reward.ProID, *s.RewardValue, rewardUom)
							ordinalConversionResult, convErr := service.ConversionWithPrice(ordinalConversion, req.DistributorID,
								req.OrderDate, req.CustID, req.ParentCustID)
							if convErr != nil {
								log.Error("Error converting reward product per ordinal:", convErr.Error())
								return responses, convErr
							}
							ordinalGross := math.Round(rewardGrossValueFromConversion(ordinalConversionResult))
							promoByOrdinal[i] = ordinalGross
							totalGross += ordinalGross
							sumQty1 += float64(ordinalConversionResult.Qty1)
							sumQty2 += float64(ordinalConversionResult.Qty2)
							sumQty3 += float64(ordinalConversionResult.Qty3)
						}

						rewardProduct.ProID = int(reward.ProID)
						rewardProduct.Qty1 = sumQty1
						rewardProduct.Qty2 = sumQty2
						rewardProduct.Qty3 = sumQty3
						rewardProduct.GrossValue = totalGross
						rewardProduct.Promo1 = promoByOrdinal[0]
						rewardProduct.Promo2 = promoByOrdinal[1]
						rewardProduct.Promo3 = promoByOrdinal[2]
						rewardProduct.Promo4 = promoByOrdinal[3]
						rewardProduct.Promo5 = promoByOrdinal[4]
					} else {
						rewardProductConversion := buildRewardProductConversionBody(req.CustID, reward.ProID, qtyReward, rewardUom)
						rewardProductConversionResult, convErr := service.ConversionWithPrice(rewardProductConversion, req.DistributorID,
							req.OrderDate, req.CustID, req.ParentCustID)
						if convErr != nil {
							log.Error("Error converting reward product:", convErr.Error())
							return responses, convErr
						}

						rewardGross := math.Round(rewardGrossValueFromConversion(rewardProductConversionResult))
						rewardProduct.ProID = int(reward.ProID)
						rewardProduct.Qty1 = float64(rewardProductConversionResult.Qty1)
						rewardProduct.Qty2 = float64(rewardProductConversionResult.Qty2)
						rewardProduct.Qty3 = float64(rewardProductConversionResult.Qty3)
						rewardProduct.GrossValue = rewardGross
						rewardProduct.Promo1 = rewardGross
					}

					if (rewardProduct.Qty1 + rewardProduct.Qty2 + rewardProduct.Qty3) > 0 {
						response.RewardProduct = append(response.RewardProduct, rewardProduct)
					}

					if remainingQtyReward <= 0 {
						break
					}
				}
			}
		}

		// Single append per promo: add response when any reward is present
		if consultPromoRespHasReward(response) {
			responses = append(responses, response)
		}
	}

	if err := resolveConsultPromoStockError(requestedPromoIDs, insufficientStockPromoIDs, responses); err != nil {
		return nil, err
	}
	return responses, nil
}

func (service *promotionServiceImpl) rewardProductStockInUom(qtyStock, convUnit2, convUnit3 int64, rewardUom *model.UomType) (float64, error) {
	conv, err := service.PromoConversionWithoutProductQuery(qtyStock, convUnit2, convUnit3)
	if err != nil {
		return 0, err
	}
	if rewardUom == nil {
		return float64(conv.Qty3), nil
	}
	switch *rewardUom {
	case model.UomTypeSmallest:
		return float64(conv.Qty1), nil
	case model.UomTypeMiddle:
		return float64(conv.Qty2), nil
	default:
		return float64(conv.Qty3), nil
	}
}

func orderQtyInRewardUom(details []entity.ConPromoV2Det, proID int, rewardUom *model.UomType) float64 {
	var total float64
	for i := range details {
		if details[i].ProID != proID {
			continue
		}
		total += detailQtyByUom(&details[i], rewardUom)
	}
	return total
}

func promoIDInList(promoID string, promoIDs []string) bool {
	for _, id := range promoIDs {
		if id == promoID {
			return true
		}
	}
	return false
}

func resolveConsultPromoStockError(requestedPromoIDs, insufficientStockPromoIDs []string, responses []entity.ConsultPromoResp) error {
	if len(insufficientStockPromoIDs) == 0 {
		return nil
	}
	if len(requestedPromoIDs) > 0 {
		for _, requestedID := range requestedPromoIDs {
			if promoIDInList(requestedID, insufficientStockPromoIDs) {
				return errors.New(errmsg.PromoInsufficientStockMessage(requestedID))
			}
		}
		return nil
	}
	if len(responses) > 0 {
		return nil
	}
	return errors.New(errmsg.PromoInsufficientStockMessage(insufficientStockPromoIDs[0]))
}

func consultPromoRespHasReward(response entity.ConsultPromoResp) bool {
	if len(response.RewardValue) > 0 || len(response.RewardPercentage) > 0 {
		return true
	}
	for _, product := range response.RewardProduct {
		if product.Qty1+product.Qty2+product.Qty3 > 0 {
			return true
		}
	}
	return false
}

func detailQtyByUom(detail *entity.ConPromoV2Det, uom *model.UomType) float64 {
	if detail == nil {
		return 0
	}
	if uom == nil {
		return detail.Qty3
	}
	switch *uom {
	case model.UomTypeMiddle:
		return detail.Qty2
	case model.UomTypeSmallest:
		return detail.Qty1
	default:
		return detail.Qty3
	}
}

func minBuyValueForCriteria(criteria model.PromotionProductCriteria) float64 {
	if criteria.MinBuyValue != nil {
		return *criteria.MinBuyValue
	}
	if criteria.MinBuyQty != nil {
		return *criteria.MinBuyQty
	}
	return 0
}

func buyValueForCriteria(detail *entity.ConPromoV2Det, criteria model.PromotionProductCriteria) float64 {
	if detail == nil {
		return 0
	}
	if criteria.MinBuyType != nil && *criteria.MinBuyType == model.RuleTypeQuantity {
		return detailQtyByUom(detail, criteria.MinBuyUom)
	}
	return float64(detail.GrossValue)
}

type rewardConfig struct {
	rewardType     model.RewardType
	ruleType       model.RuleType
	ruleUom        *model.UomType
	rewardValuePtr *float64
	rewardUom      *model.UomType
	perScope       *string
}

func rewardConfigFromSlabStrata(slab model.PromotionV2Slabs, hasSlab bool, strata model.PromotionV2Strata, hasStrata bool) (rewardConfig, bool) {
	if hasStrata {
		return rewardConfig{
			rewardType:     strata.RewardType,
			rewardValuePtr: strata.RewardValue,
			rewardUom:      strata.RewardUom,
			perScope:       strata.PerScope,
			ruleType:       strata.RuleType,
			ruleUom:        strata.RuleUom,
		}, true
	}
	if hasSlab {
		return rewardConfig{
			rewardType:     slab.RewardType,
			rewardValuePtr: slab.RewardValue,
			rewardUom:      slab.RewardUom,
			perScope:       slab.PerScope,
			ruleType:       slab.RuleType,
			ruleUom:        slab.RuleUom,
		}, true
	}
	return rewardConfig{}, false
}

func perScopeRewardValue(rewardValueCalculated float64, perScope *string, ruleType model.RuleType, ruleUom *model.UomType, detail *entity.ConPromoV2Det, groupSize int) float64 {
	if perScope == nil {
		return rewardValueCalculated
	}

	log.Info("perScopeRewardValue > perScope > ", *perScope)
	log.Info("perScopeRewardValue > detail > ", structs.StructToJson(detail))

	switch *perScope {
	case string(model.PerScopeOrder):
		if groupSize > 0 {
			return rewardValueCalculated / float64(groupSize)
		}
		return rewardValueCalculated
	case string(model.PerScopeProduct):
		if detail == nil {
			return rewardValueCalculated
		}
		switch ruleType {
		case model.RuleTypeValue:
			return rewardValueCalculated * float64(detail.Qty3)
		case model.RuleTypeQuantity:
			qty := detailQtyByUom(detail, ruleUom)
			return rewardValueCalculated * qty
		default:
			return rewardValueCalculated
		}
	default:
		return rewardValueCalculated
	}
}

func calculateStrataRuleValue(strata model.PromotionV2Strata, details map[int]*entity.ConPromoV2Det, orderDetails []entity.ConPromoV2Det) float64 {
	if strata.RuleType == model.RuleTypeValue {
		value := 0.0
		for _, d := range orderDetails {
			value += float64(d.GrossValue)
		}
		return value
	}

	value := 0.0
	for _, d := range details {
		qty := detailQtyByUom(d, strata.RuleUom)
		value += qty
	}
	return value
}

// resolveNonSequentialStrataOrdinal returns the highest strata ordinal that applies.
// Values inside a tier range use that tier; values above the top tier cap at the highest tier.
func resolveNonSequentialStrataOrdinal(strata []model.PromotionV2Strata, ruleValue float64) int {
	matchedOrdinal := 0
	for _, s := range strata {
		if ruleValue >= s.RangeFrom && ruleValue <= s.RangeTo && s.Ordinal > matchedOrdinal {
			matchedOrdinal = s.Ordinal
		}
	}
	if matchedOrdinal > 0 {
		return matchedOrdinal
	}

	topOrdinal := 0
	for _, s := range strata {
		if ruleValue >= s.RangeFrom && s.Ordinal > topOrdinal {
			topOrdinal = s.Ordinal
		}
	}
	return topOrdinal
}

func (service *promotionServiceImpl) PromoConversion(req entity.PromoConversionReq, custID string) (response entity.PromoConversionResp, err error) {
	// ponytail: import (data_source=3) may pass ProId=0 for blank ProCode rows.
	// passthrough qty without product lookup so the import row keeps flowing.
	if req.ProductID == 0 {
		response.Qty1 = req.Qty1
		response.Qty2 = req.Qty2
		response.Qty3 = req.Qty3
		response.TotalQty = req.Qty1 + req.Qty2 + req.Qty3
		return response, nil
	}
	product, err := service.PromotionRepository.FindProductByIDAndCustID(req.ProductID, custID)
	if err != nil {
		return response, err
	}

	qty1 := req.Qty1
	qty2 := req.Qty2
	qty3 := req.Qty3

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

	response.TotalQty = (int64(product.ConvUnit2)*int64(product.ConvUnit3))*qty3 + (int64(product.ConvUnit2) * qty2) + qty1

	response.Qty1 = response.TotalQty
	response.Qty2 = response.TotalQty / int64(product.ConvUnit2)
	response.Qty3 = qty3

	return response, err
}

func buildRewardProductConversionBody(custID string, productID int64, qtyReward float64, rewardUom *model.UomType) entity.CreateConversionBody {
	qty := int64(0)
	conversionBody := entity.CreateConversionBody{
		CustId:    custID,
		ProductId: productID,
	}
	if rewardUom == nil {
		return conversionBody
	}
	switch *rewardUom {
	case model.UomTypeSmallest:
		conversionBody.Qty1 = int64(qtyReward)
		conversionBody.Qty2 = qty
		conversionBody.Qty3 = qty
	case model.UomTypeMiddle:
		conversionBody.Qty1 = qty
		conversionBody.Qty2 = int64(qtyReward)
		conversionBody.Qty3 = qty
	default:
		conversionBody.Qty1 = qty
		conversionBody.Qty2 = qty
		conversionBody.Qty3 = int64(qtyReward)
	}
	return conversionBody
}

func rewardGrossValueFromConversion(conversionResult entity.ConversionWithPriceResp) float64 {
	return float64(conversionResult.Qty1)*conversionResult.SellPrice1 +
		float64(conversionResult.Qty2)*conversionResult.SellPrice2 +
		float64(conversionResult.Qty3)*conversionResult.SellPrice3
}

func (service *promotionServiceImpl) ConversionWithPrice(conversionBody entity.CreateConversionBody, distributorID int64, transDate, custID, parentCustID string) (response entity.ConversionWithPriceResp, err error) {
	product, err := service.PromotionRepository.FindProductAndPriceByID(conversionBody.ProductId, distributorID, transDate, custID, parentCustID)
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
	response.SellPrice1 = product.SellPrice1
	response.SellPrice2 = product.SellPrice2
	response.SellPrice3 = product.SellPrice3

	response.TotalQty = (int64(product.ConvUnit2)*int64(product.ConvUnit3))*qty3 + (int64(product.ConvUnit2) * qty2) + qty1

	return response, err
}

func (service *promotionServiceImpl) PromoConversionWithoutProductQuery(totalQty int64, convUnit2 int64, convUnit3 int64) (response entity.PromoConverResp, err error) {

	qty1 := totalQty
	qty2 := int64(0)
	qty3 := int64(0)

	rQty2 := qty1 / convUnit2
	if rQty2 > 0 {
		qty1 = qty1 % convUnit2
		qty2 += rQty2
	}

	rQty3 := qty2 / convUnit3
	if rQty3 > 0 {
		qty2 = qty2 % convUnit3
		qty3 += rQty3
	}

	response.TotalQty = (convUnit2*convUnit3)*qty3 + (convUnit2 * qty2) + qty1

	response.Qty1 = response.TotalQty
	response.Qty2 = response.TotalQty / convUnit2
	response.Qty3 = qty3

	return response, nil
}

func parseOptionalPromotionDate(dateStr string) (*time.Time, error) {
	trimmed := strings.TrimSpace(dateStr)
	if trimmed == "" {
		return nil, nil
	}

	layouts := []string{constant.DATE_FORMAT_DD_MM_YYYY, constant.YYYY_MM_DD}
	for _, layout := range layouts {
		t, err := time.Parse(layout, trimmed)
		if err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("invalid date format %q, use DD/MM/YYYY or YYYY-MM-DD", dateStr)
}

func applyPromotionV2ExtendedFields(promo *model.PromotionV2, budgetID, claimFrom, claimTo string, vatRate, whtRate *float64) error {
	if budgetID != "" {
		promo.BudgetID = &budgetID
	} else {
		promo.BudgetID = nil
	}

	claimFromDate, err := parseOptionalPromotionDate(claimFrom)
	if err != nil {
		return err
	}
	claimToDate, err := parseOptionalPromotionDate(claimTo)
	if err != nil {
		return err
	}
	if claimFromDate != nil && claimToDate != nil && claimToDate.Before(*claimFromDate) {
		return fmt.Errorf("claim_date_to must be on or after claim_date_from")
	}

	promo.ClaimDateFrom = claimFromDate
	promo.ClaimDateTo = claimToDate
	promo.VatRate = vatRate
	promo.WhtRate = whtRate

	return nil
}

func mapPromotionV2ExtendedFieldsToEntity(promo model.PromotionV2, response *entity.PromotionV2) {
	if promo.ExistingPromoID != nil {
		response.ExistingPromoID = *promo.ExistingPromoID
	} else {
		response.ExistingPromoID = ""
	}
	if promo.BudgetID != nil {
		response.BudgetID = *promo.BudgetID
	}
	if promo.ClaimDateFrom != nil {
		response.ClaimDateFrom = promo.ClaimDateFrom.Format(constant.DATE_FORMAT_DD_MM_YYYY)
	}
	if promo.ClaimDateTo != nil {
		response.ClaimDateTo = promo.ClaimDateTo.Format(constant.DATE_FORMAT_DD_MM_YYYY)
	}
	response.VatRate = promo.VatRate
	response.WhtRate = promo.WhtRate
}

func mapPromotionV2ResponseCustID(promo model.PromotionV2, response *entity.PromotionV2) {
	response.CustID = strings.TrimSpace(promo.DistributorCustID)
	if response.CustID == "" {
		response.CustID = strings.TrimSpace(promo.CustID)
	}
	response.DistributorCustID = ""
}

func mapPromoCoverageDistributorFromModel(row model.PromotionCoverageDistributors) entity.PromoCoverageDistributor {
	return entity.PromoCoverageDistributor{
		DistributorID:   row.DistributorID,
		DistributorCode: row.DistributorCode,
		DistributorName: row.DistributorName,
	}
}

func mapPromoOutletSalesTeamFromModel(row model.PromotionOutletAttributeSalesTeam) entity.PromoOutletSalesTeamSelected {
	return entity.PromoOutletSalesTeamSelected{
		SalesTeamID:   row.SalesTeamID,
		SalesTeamCode: row.SalesTeamCode,
		SalesTeamName: row.SalesTeamName,
	}
}

func mapPromoOutletSelectedFromModel(row model.PromotionOutletsSelected) entity.PromoOutletSelected {
	return entity.PromoOutletSelected{
		OutletID:        row.OutletID,
		OutletCode:      row.OutletCode,
		OutletName:      row.OutletName,
		DistributorCode: row.DistributorCode,
		DistributorName: row.DistributorName,
	}
}
