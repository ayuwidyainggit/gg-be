package service

import (
	"mobile/entity"
	"mobile/pkg/config/env"
	"mobile/repository"
)

type TakingOrderService interface {
	List(dataFilter entity.GeneralQueryFilter) (responses []entity.NoOrderReasonResp, err error)
}

type TakingOrderServiceImpl struct {
	Config                env.ConfigEnv
	TakingOrderRepository repository.MTakingOrderRepository
}

func NewTakingOrderService(
	config env.ConfigEnv, takingOrderRepo repository.MTakingOrderRepository,
) *TakingOrderServiceImpl {
	return &TakingOrderServiceImpl{
		Config:                config,
		TakingOrderRepository: takingOrderRepo,
	}
}

func (service *TakingOrderServiceImpl) List(dataFilter entity.GeneralQueryFilter) (responses []entity.NoOrderReasonResp, err error) {
	takingOrders, err := service.TakingOrderRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return
	}

	responses = make([]entity.NoOrderReasonResp, 0)
	for _, takingOrder := range takingOrders {
		orderResp := entity.NoOrderReasonResp{
			TakingOrderId: takingOrder.TakingOrderId,
			Reason:        takingOrder.TakingOrderName,
			ImageUrl:      takingOrder.ImageUrl,
		}
		responses = append(responses, orderResp)
	}
	return
}
