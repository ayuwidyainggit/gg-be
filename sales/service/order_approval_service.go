package service

import (
	"context"
	"sales/entity"
	"sales/model"
	"sales/pkg/structs"
	"sales/repository"
)

type OrderApprovalService interface {
	List(dataFilter entity.OrderApprovalQueryFilter) (data []entity.OrderApprovalListResponse, total int64, lastPage int, err error)
	UpdateStatusDetail(OrderApprovalRequestID int64, empID int64, status int) (err error)
}

func NewOrderApprovalService(orderApprovalRepo repository.OrderApprovalRepository, transaction repository.Dbtransaction) *OrderApprovalServiceImpl {
	return &OrderApprovalServiceImpl{
		Transaction:             transaction,
		OrderApprovalRepository: orderApprovalRepo,
	}
}

type OrderApprovalServiceImpl struct {
	Transaction             repository.Dbtransaction
	OrderApprovalRepository repository.OrderApprovalRepository
}

func (service *OrderApprovalServiceImpl) List(dataFilter entity.OrderApprovalQueryFilter) (datas []entity.OrderApprovalListResponse, total int64, lastPage int, err error) {
	var dataResults []model.OrderApprovalRead
	if dataFilter.Status == 0 { // jika perlu
		dataResults, total, lastPage, err = service.OrderApprovalRepository.FindNeedReview(dataFilter)
		if err != nil {
			return datas, 0, 0, err
		}

	} else if dataFilter.Status == 1 {
		dataResults, total, lastPage, err = service.OrderApprovalRepository.FindApproved(dataFilter)
		if err != nil {
			return datas, 0, 0, err
		}

	} else {
		dataResults, total, lastPage, err = service.OrderApprovalRepository.FindRejected(dataFilter)
		if err != nil {
			return datas, 0, 0, err
		}

	}

	for _, dataResult := range dataResults {
		var vResp entity.OrderApprovalListResponse

		structs.Automapper(dataResult, &vResp)
		if dataResult.RoDate != nil {
			roDate := dataResult.RoDate.Format("2006-01-02")
			vResp.RoDate = roDate
		}

		vResp.CreditLimitValue = dataResult.CreditLimitValue
		vResp.SalesInvLimitValue = dataResult.SalesInvLimitValue
		vResp.ObsLimitValue = dataResult.ObsLimitValue
		datas = append(datas, vResp)
	}

	return
}

func (service *OrderApprovalServiceImpl) UpdateStatusDetail(OrderApprovalRequestID int64, empID int64, status int) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.OrderApprovalRepository.UpdateStatusDetail(txCtx, OrderApprovalRequestID, empID, status)
		if err != nil {
			return err
		}

		_, err := service.OrderApprovalRepository.FindIfNeedsApproval(txCtx, OrderApprovalRequestID)
		if err != nil { // jika sudah tidak ada record, artinya sudah tidak ada yang perlu di review, maka set order approvalnya finished
			service.OrderApprovalRepository.UpdateOrderApprovalFinished(txCtx, OrderApprovalRequestID)

			if status == entity.ORDER_APPROVAL_APPROVED {
				orderApprovalRequest, err := service.OrderApprovalRepository.FindOrderApprovalByID(OrderApprovalRequestID)
				if err != nil {
					return err
				}
				service.OrderApprovalRepository.UpdateStatusOrder(txCtx, orderApprovalRequest.RoNo, entity.PROCESSED)
			}

			return nil
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
