package service

import (
	"finance/entity"
	"finance/pkg/structs"
	"finance/repository"
)

type ApListService interface {
	Detail(invNo string, ParentCustId string, custID string) (response entity.AccountPayableListDetailResponse, err error)
	List(dataFilter entity.AccountPayableListQueryFilter) (data []entity.AccountPayableListResponse, total int64, lastPage int, err error)
}

type ApListServiceImpl struct {
	Repository  repository.ApListRepository
	Transaction repository.Dbtransaction
}

func NewApListService(repository repository.ApListRepository, transaction repository.Dbtransaction) *ApListServiceImpl {
	return &ApListServiceImpl{
		Repository:  repository,
		Transaction: transaction,
	}
}

func (service *ApListServiceImpl) Detail(invNo string, ParentCustId string, custID string) (response entity.AccountPayableListDetailResponse, err error) {
	ap, err := service.Repository.FindByNo(invNo, ParentCustId, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ap, &response)
	if err != nil {
		return response, err
	}

	details, err := service.Repository.FindDetail(*response.InvNo, custID)
	if err != nil {
		return response, err
	}
	detailsData := make([]entity.AccountPayableListPaymentHistoryResponse, 0, len(details))
	// var amountInv = ap.InvAmount
	for _, detail := range details {
		if detail.PaymentMethod == nil {
			continue
		}

		var detailData entity.AccountPayableListPaymentHistoryResponse

		// 1: Cash; 2: Check; 3: Transfer; 4: Credit; 5:Return
		if *detail.PaymentMethod == 1 {
			detailData.PaymentMethodName = "Cash"
		}
		if *detail.PaymentMethod == 2 {
			detailData.PaymentMethodName = "Check"
		}
		if *detail.PaymentMethod == 3 {
			detailData.PaymentMethodName = "Transfer"
		}
		if *detail.PaymentMethod == 4 {
			detailData.PaymentMethodName = "Credit"
		}
		if *detail.PaymentMethod == 5 {
			detailData.PaymentMethodName = "Return"
		}

		if detail.PaymentDate != nil {
			paymentDate := detail.PaymentDate.Format("2006-01-02")
			detailData.PaymentDate = &paymentDate
		}

		if detail.PaymentBalance != nil {
			detailData.PaymentBalance = detail.PaymentBalance
		}

		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		// amountInv = amountInv - detailData.Amount
		// detailData.PaymentBalance = amountInv

		detailsData = append(detailsData, detailData)
	}

	if response.InvAmount != response.AmountPaid {
		response.InvStatus = "Outstanding"
	} else {
		response.InvStatus = "Paid"
		response.Aging = 0
	}

	invDate := ap.InvDate.Format("2006-01-02")
	response.InvDate = &invDate

	invDueDate := ap.InvDueDate.Format("2006-01-02")
	response.InvDueDate = &invDueDate

	CreatedAt := ap.CreatedAt.Format("2006-01-02")
	response.CreatedAt = &CreatedAt

	response.PaymentHistory = detailsData
	return response, nil
}

func (service *ApListServiceImpl) List(dataFilter entity.AccountPayableListQueryFilter) (data []entity.AccountPayableListResponse, total int64, lastPage int, err error) {
	aps, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range aps {
		var vResp entity.AccountPayableListResponse
		structs.Automapper(row, &vResp)
		if row.InvDate != nil {
			invDate := row.InvDate.Format("2006-01-02")
			vResp.InvDate = &invDate
		}
		if row.InvDueDate != nil {
			InvDueDate := row.InvDueDate.Format("2006-01-02")
			vResp.InvDueDate = &InvDueDate
		}

		if row.CreatedAt != nil {
			CreatedAt := row.CreatedAt.Format("2006-01-02")
			vResp.CreatedAt = &CreatedAt
		}

		if row.UpdatedAt != nil {
			UpdatedAt := row.UpdatedAt.Format("2006-01-02")
			vResp.UpdatedAt = &UpdatedAt
		}

		if row.InvAmount != row.AmountPaid {
			vResp.InvStatus = "Outstanding"
		} else {
			vResp.InvStatus = "Paid"
			vResp.Aging = 0
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err

}
