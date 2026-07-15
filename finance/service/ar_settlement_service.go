package service

import (
	"context"
	"finance/entity"
	"finance/pkg/structs"
	"finance/repository"
	"fmt"
	"log"
)

type ArSettlementService interface {
	Detail(depositNo string, custID string) (response entity.ArSettlementResponse, err error)
	List(dataFilter entity.ArSettlementQueryFilter) (data []entity.ArSettlementListResponse, total int64, lastPage int, err error)
	Approve(custId string, depositNo string, userId int64) (err error)
	BulkApprove(items []entity.BulkApproveArSettlementItem, userId int64) (err error)
	Reject(custId string, depositNo string, userId int64) (err error)
	BulkReject(items []entity.BulkApproveArSettlementItem, userId int64) (err error)

	CollectorLookupList(entity.GeneralQueryFilter) (data []entity.EmployeeLookupResponse, total int64, lastPage int, err error)
	DepositStatusLookupList(entity.GeneralQueryFilter) (data []entity.DepositStatusLookupResponse, total int64, lastPage int, err error)

	VerifyRejectData(depositNo string, custId string) (entity.RejectVerifyReport, error)
}

type arSettlementServiceImpl struct {
	Repository   repository.ArSettlementRepository
	DepositRepo  repository.DepositRepository
	Transaction  repository.Dbtransaction
}

func NewArSettlementService(repository repository.ArSettlementRepository, depositRepo repository.DepositRepository, transaction repository.Dbtransaction) *arSettlementServiceImpl {
	return &arSettlementServiceImpl{
		Repository:  repository,
		DepositRepo: depositRepo,
		Transaction: transaction,
	}
}

func (service *arSettlementServiceImpl) Detail(depositNo string, custID string) (response entity.ArSettlementResponse, err error) {
	// var settlement model.ArSettlementList
	// var branchSettlement model.ArBranchSettlementList
	var DetailsData []entity.ArSettlementPaymentResponse
	branchCode := "B"
	dNo := string([]rune(depositNo)[2])
	if dNo == branchCode {
		settlement, err := service.Repository.FindOneByBranchDepositNo(depositNo, custID)
		if err != nil {
			return response, err
		}

		if err = structs.Automapper(settlement, &response); err != nil {
			return response, err
		}

		if settlement.DepositDate != nil {
			depositDate := settlement.DepositDate.Format("2006-01-02")
			response.DepositDate = &depositDate
		}

		if settlement.CollectionDate != nil {
			collectionDate := settlement.CollectionDate.Format("2006-01-02")
			response.CollectionDate = &collectionDate
		}

		depositStatusName := response.GenerateDataDepositStatusName()
		response.DepositStatusName = &depositStatusName
	} else {
		settlement, err := service.Repository.FindOneByDepositNo(depositNo, custID)
		if err != nil {
			return response, err
		}

		if err = structs.Automapper(settlement, &response); err != nil {
			return response, err
		}

		if settlement.DepositDate != nil {
			depositDate := settlement.DepositDate.Format("2006-01-02")
			response.DepositDate = &depositDate
		}

		if settlement.CollectionDate != nil {
			collectionDate := settlement.CollectionDate.Format("2006-01-02")
			response.CollectionDate = &collectionDate
		}

		depositStatusName := response.GenerateDataDepositStatusName()
		response.DepositStatusName = &depositStatusName
	}

	if dNo == branchCode {
		Details, err := service.Repository.FindBranchDetail(depositNo, custID)
		if err != nil {
			return response, err
		}
		for _, detail := range Details {
			var detailData entity.ArSettlementPaymentResponse
			err = structs.Automapper(detail, &detailData)
			if err != nil {
				return response, err
			}

			if detail.InvoiceDate != nil {
				invoiceDate := detail.InvoiceDate.Format("2006-01-02")
				detailData.InvoiceDate = &invoiceDate
			}

			payTypeName := detailData.GenerateDataPayTypeName()
			detailData.PayTypeName = &payTypeName

			DetailsData = append(DetailsData, detailData)
		}
	} else {
		Details, err := service.Repository.FindDetail(depositNo, custID)
		if err != nil {
			return response, err
		}
		for _, detail := range Details {
			var detailData entity.ArSettlementPaymentResponse
			err = structs.Automapper(detail, &detailData)
			if err != nil {
				return response, err
			}

			if detail.InvoiceDate != nil {
				invoiceDate := detail.InvoiceDate.Format("2006-01-02")
				detailData.InvoiceDate = &invoiceDate
			}

			payTypeName := detailData.GenerateDataPayTypeName()
			detailData.PayTypeName = &payTypeName

			DetailsData = append(DetailsData, detailData)
		}

	}

	response.Details = DetailsData

	if dNo != branchCode {
		sumRem, sumErr := service.Repository.SumInvoiceRemainingForDepositSettlement(depositNo, custID)
		if sumErr != nil {
			return response, sumErr
		}
		response.RemainingAmount = &sumRem
	}

	expenses, err := service.DepositRepo.FindExpenseByDepositNo(depositNo, custID)
	if err != nil {
		return response, err
	}
	var expenseTotal float64
	expenseList := make([]entity.ArSettlementExpenseItem, 0, len(expenses))
	for _, e := range expenses {
		amt, _ := e.PaymentAmount.Float64()
		expenseTotal += amt
		expenseList = append(expenseList, entity.ArSettlementExpenseItem{
			DepositExpenseID: e.DepositExpenseID,
			DocNo:            e.DocNo,
			Balance:          e.Balance,
			PaymentAmount:    amt,
		})
	}
	response.ExpenseTotal = &expenseTotal
	response.Expense = expenseList

	return response, nil
}

func (service *arSettlementServiceImpl) DetailOld(depositNo string, custID string) (response entity.ArSettlementResponse, err error) {
	settlement, err := service.Repository.FindOneByDepositNo(depositNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(settlement, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.Repository.FindDetail(depositNo, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.ArSettlementPaymentResponse
	for _, detail := range Details {
		var detailData entity.ArSettlementPaymentResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		if detail.InvoiceDate != nil {
			invoiceDate := detail.InvoiceDate.Format("2006-01-02")
			detailData.InvoiceDate = &invoiceDate
		}

		payTypeName := detailData.GenerateDataPayTypeName()
		detailData.PayTypeName = &payTypeName

		DetailsData = append(DetailsData, detailData)
	}
	log.Println("AR Settlement Service (Details) : ", Details)

	if settlement.DepositDate != nil {
		depositDate := settlement.DepositDate.Format("2006-01-02")
		response.DepositDate = &depositDate
	}

	if settlement.CollectionDate != nil {
		collectionDate := settlement.CollectionDate.Format("2006-01-02")
		response.CollectionDate = &collectionDate
	}

	depositStatusName := response.GenerateDataDepositStatusName()
	response.DepositStatusName = &depositStatusName

	response.Details = DetailsData
	return response, nil
}

func (service *arSettlementServiceImpl) List(dataFilter entity.ArSettlementQueryFilter) (data []entity.ArSettlementListResponse, total int64, lastPage int, err error) {
	// var settlements []model.ArSettlementList

	if dataFilter.CustId == dataFilter.ParentCustId {
		settlements, total, lastPage, err := service.Repository.FindAllByCustIdNew(dataFilter)
		if err != nil {
			return data, total, lastPage, err
		}

		for _, settlement := range settlements {
			var vResp entity.ArSettlementListResponse
			structs.Automapper(settlement, &vResp)
			if settlement.DepositDate != nil {
				depositDate := settlement.DepositDate.Format("2006-01-02")
				vResp.DepositDate = &depositDate
			}

			depositStatusName := vResp.GenerateDataDepositStatusName()
			vResp.DepositStatusName = &depositStatusName

			data = append(data, vResp)
		}

		return data, total, lastPage, err
	} else {
		settlements, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
		if err != nil {
			return data, total, lastPage, err
		}

		for _, settlement := range settlements {
			var vResp entity.ArSettlementListResponse
			structs.Automapper(settlement, &vResp)
			if settlement.DepositDate != nil {
				depositDate := settlement.DepositDate.Format("2006-01-02")
				vResp.DepositDate = &depositDate
			}

			if settlement.DepositNo != nil {
				sumRem, sumErr := service.Repository.SumInvoiceRemainingForDepositSettlement(*settlement.DepositNo, dataFilter.CustId)
				if sumErr != nil {
					return data, total, lastPage, sumErr
				}
				vResp.RemainingAmount = &sumRem
			}

			depositStatusName := vResp.GenerateDataDepositStatusName()
			vResp.DepositStatusName = &depositStatusName

			data = append(data, vResp)
		}

		return data, total, lastPage, err
	}
}

func (service *arSettlementServiceImpl) Approve(custId string, depositNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		branchCode := "B"
		dNo := string([]rune(depositNo)[2])
		fmt.Println(dNo)
		fmt.Println(branchCode)
		if dNo == branchCode {
			err = service.Repository.ApproveBranch(txCtx, custId, depositNo, userId)
			if err != nil {
				return err
			}

			// deposits, err := service.Repository.FindDetail(depositNo, custId)
			// if err != nil {
			// 	return err
			// }
			// var invoicesNoPaidOff []string
			// for _, deposit := range deposits {
			// 	depositStatus := []int{entity.DEPOSIT_STATUS_IN_REVIEW, entity.DEPOSIT_STATUS_IN_APPROVED}
			// 	depositDetails, err := service.Repository.FindBranchDetailByInvoice(txCtx, *deposit.InvoiceNo, depositStatus, custId)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	invoiceAmount := depositDetails[0].InvoiceAmount
			// 	var totalPayment float64
			// 	for _, depositDetail := range depositDetails {
			// 		totalPayment += depositDetail.TotalPayment
			// 	}

			// 	if totalPayment >= invoiceAmount {
			// 		invoicesNoPaidOff = append(invoicesNoPaidOff, *deposit.InvoiceNo)
			// 	}
			// }

			// if len(invoicesNoPaidOff) > 0 {
			// 	err := service.Repository.SetBranchInvoiceToPaidOff(txCtx, invoicesNoPaidOff, custId)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
		} else {
			err = service.Repository.Approve(txCtx, custId, depositNo, userId)
			if err != nil {
				return err
			}

			deposits, err := service.Repository.FindDetail(depositNo, custId)
			if err != nil {
				return err
			}
			var invoicesNoPaidOff []string
			for _, deposit := range deposits {
				depositStatus := []int{entity.DEPOSIT_STATUS_IN_REVIEW, entity.DEPOSIT_STATUS_IN_APPROVED}
				depositDetails, err := service.Repository.FindDetailByInvoice(txCtx, *deposit.InvoiceNo, depositStatus, custId)
				if err != nil {
					return err
				}

				invoiceAmount := depositDetails[0].InvoiceAmount
				var totalPayment float64
				for _, depositDetail := range depositDetails {
					totalPayment += depositDetail.TotalPayment
				}

				if totalPayment >= invoiceAmount {
					invoicesNoPaidOff = append(invoicesNoPaidOff, *deposit.InvoiceNo)
				}
			}

			if len(invoicesNoPaidOff) > 0 {
				err := service.Repository.SetInvoiceToPaidOff(txCtx, invoicesNoPaidOff, custId)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	return err
}

func (service *arSettlementServiceImpl) BulkApprove(items []entity.BulkApproveArSettlementItem, userId int64) (err error) {
	for _, item := range items {
		if err = service.Approve(item.CustId, item.DepositNo, userId); err != nil {
			return err
		}
	}
	return nil
}

func (service *arSettlementServiceImpl) Reject(custId string, depositNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		branchCode := "B"
		dNo := string([]rune(depositNo)[2])
		if dNo == branchCode {
			if err = service.Repository.RejectBranch(txCtx, custId, depositNo, userId); err != nil {
				return err
			}
		} else {
			if err = service.Repository.Reject(txCtx, custId, depositNo, userId); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (service *arSettlementServiceImpl) BulkReject(items []entity.BulkApproveArSettlementItem, userId int64) (err error) {
	for _, item := range items {
		if err = service.Reject(item.CustId, item.DepositNo, userId); err != nil {
			return err
		}
	}
	return nil
}

func (service *arSettlementServiceImpl) CollectorLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.EmployeeLookupResponse, total int64, lastPage int, err error) {
	Collectors, total, lastPage, err := service.Repository.FindAllCollectorByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Collectors {
		var vResp entity.EmployeeLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *arSettlementServiceImpl) DepositStatusLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.DepositStatusLookupResponse, total int64, lastPage int, err error) {
	DepositStatuses, total, lastPage, err := service.Repository.FindAllDepositStatusLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, depositStatus := range DepositStatuses {
		var vResp entity.DepositStatusLookupResponse
		structs.Automapper(depositStatus, &vResp)

		depositStatusName := vResp.GenerateDataDepositStatusName()
		vResp.DepositStatusName = &depositStatusName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *arSettlementServiceImpl) VerifyRejectData(depositNo string, custId string) (entity.RejectVerifyReport, error) {
	ctx := context.Background()
	return service.Repository.VerifyRejectData(ctx, depositNo, custId)
}
