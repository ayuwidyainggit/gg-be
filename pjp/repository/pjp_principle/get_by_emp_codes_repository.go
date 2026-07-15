package pjp_principle

import (
	"context"
	"errors"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

func (repo *pjpPrincipleRepository) GetPjpsByEmpCodes(ctx context.Context, tx *gorm.DB, empCode []string, currentCustomerId string) model.PjpPrinciple {
	var pjps model.PjpPrinciple

	err := tx.WithContext(ctx).
		Select("pjp_principles.permanent_journey_plans.*, cust.distributor_id").
		Where("pjp_principles.permanent_journey_plans.cust_id = ? AND pjp_principles.permanent_journey_plans.salesman_code IN ?", currentCustomerId, empCode).
		Joins("LEFT JOIN smc.m_customer cust ON pjp_principles.permanent_journey_plans.cust_id = cust.cust_id").
		Take(&pjps).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helper.ErrorPanic(errors.New("pjp principle not found"))
		}
		helper.ErrorPanic(errors.New("error getting pjp principle"))
	}

	return pjps
}
