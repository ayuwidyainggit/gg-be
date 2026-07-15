package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"finance/pkg/structs"
	"finance/repository"
	"log"
	"time"
)

type ArService interface {
	// Store(request entity.CreateArBody) (err error)
	Detail(invoiceNo string, custID string, parentCustID string) (response entity.ArResponse, err error)
	List(dataFilter entity.InvoiceQueryFilter) (data []entity.ArListResponse, total int64, lastPage int, err error)
	// Delete(custId string, arNo string, userId int64) (err error)
	// Update(arNo string, request entity.UpdateArbody) (err error)
	StoreCollection(request entity.CreateCollectionBody) (err error)
	CollectionDetail(collectionNo string, custID string, parentCustId string) (response entity.CollectionResponse, err error)
	CollectionList(dataFilter entity.CollectionQueryFilter) (data []entity.CollectionListResponse, total int64, lastPage int, err error)
	DeleteCollection(custId string, collectionNo string, userId int64) (err error)
	UpdateCollection(collectionNo string, request entity.UpdateCollectionBody) (err error)
	ReplaceCollection(collectionNo string, request entity.UpdateCollectionBody) (err error)
	PrintCollection(custId string, collectionNo string, userId int64) (err error)
	EmployeeGroupLookupList(entity.GeneralQueryFilter) (data []entity.EmployeeGroupLookupResponse, total int64, lastPage int, err error)
	EmployeeLookupList(entity.EmployeeListQueryFilter) (data []entity.EmployeeLookupResponse, total int64, lastPage int, err error)
	InvoiceList(dataFilter entity.InvoiceQueryFilter) (data []entity.InvoiceListResponse, total int64, lastPage int, err error)
	OutletFilterLookupList(entity.GeneralQueryFilter) (data []entity.OutletLookupResponse, total int64, lastPage int, err error)
	CollectorLookupList(entity.GeneralQueryFilter) (data []entity.EmployeeLookupResponse, total int64, lastPage int, err error)
	OutletGroupFilterLookupList(entity.GeneralQueryFilter) (data []entity.OutletGroupLookupResponse, total int64, lastPage int, err error)
	SalesmanFilterLookupList(entity.GeneralQueryFilter) (data []entity.SalesmanLookupResponse, total int64, lastPage int, err error)
}

type arServiceImpl struct {
	Repository  repository.ArRepository
	Transaction repository.Dbtransaction
}

func NewArService(repository repository.ArRepository, transaction repository.Dbtransaction) *arServiceImpl {
	return &arServiceImpl{
		Repository:  repository,
		Transaction: transaction,
	}
}

/*
	func (service *arServiceImpl) Store(request entity.CreateArBody) (err error) {
		c := context.Background()

		// parse time format YYYY-mm-dd to Rfc3339
		if request.ArDate != nil {
			arDate, err := str.DateStrToRfc3339String(*request.ArDate)
			if err != nil {
				return err
			}
			request.ArDate = &arDate
		}

		var Armodel model.Ar
		err = structs.Automapper(request, &Armodel)
		if err != nil {
			return err
		}
		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
			err := service.Repository.Store(txCtx, &Armodel)
			if err != nil {
				return err
			}

			for _, Detail := range request.Details {
				var detModel model.ArDet
				err = structs.Automapper(Detail, &detModel)
				if err != nil {
					return err
				}
				detModel.CustID = request.CustID
				detModel.ArNo = Armodel.ArNo
				err = service.Repository.StoreDetail(txCtx, &detModel)
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
*/
func (service *arServiceImpl) Detail(invoiceNo string, custID string, parentCustID string) (response entity.ArResponse, err error) {
	invoice, err := service.Repository.FindByInvoiceNo(invoiceNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(invoice, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.Repository.FindDetail(invoiceNo, custID, parentCustID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.ArPaymentResponse
	for _, detail := range Details {
		var detailData entity.ArPaymentResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		detailData.VisitDate = detail.VisitDate.Format("2006-01-02")

		if detail.VerifiedDate != nil {
			verifiedDate := detail.VerifiedDate.Format("2006-01-02")
			detailData.VerifiedDate = &verifiedDate
		}

		// hitung berdasarkan total payment terhadap no invoice untuk deposit detail
		detailData.CollectionStatus = 2
		if detail.TotalPayment > 0 {
			detailData.CollectionStatus = 1
		}

		// cek payment amount >= remaining amount (deposit_detail)
		detailData.PaymentOption = 2
		if detail.RemainingPayment <= 0 {
			detailData.PaymentOption = 1
		}

		detailData.CollectionStatusName = detailData.GenerateDataCollectionStatusName()
		detailData.PaymentOptionName = detailData.GenerateDataPaymentOptionName()
		detailData.PaymentMethodName = detailData.GenerateDataPaymentMethodName()
		detailData.VerificationStatusName = detailData.GenerateDataVerificationStatusName()

		DetailsData = append(DetailsData, detailData)
	}

	if invoice.InvoiceDate != nil {
		invoiceDate := invoice.InvoiceDate.Format("2006-01-02")
		response.InvoiceDate = &invoiceDate
	}

	var dueDate time.Time
	if invoice.DueDate != nil {
		dueDate = *invoice.DueDate
		dueDateStr := dueDate.Format("2006-01-02")
		response.DueDate = &dueDateStr
	}

	invoicePaidAmount, err := service.Repository.CountInvoicePaidAmount(invoiceNo, custID)
	if err != nil {
		return response, err
	}
	response.PaidAmount = invoicePaidAmount.PaidAmount
	response.RemainingAmount = response.InvoiceAmount - response.PaidAmount
	response.InvoiceStatus = int64(entity.InvoiceStatusOutstanding)
	dueDateStatus := int64(1)
	response.DueDateStatus = &dueDateStatus
	paymentDate := time.Now()

	if response.RemainingAmount <= 0 {
		response.InvoiceStatus = int64(entity.InvoiceStatusPaid)
		response.RemainingAmount = 0
		response.DueDateStatus = nil
	} else {
		if paymentDate.Format("2006-01-02") >= dueDate.Format("2006-01-02") {
			dueDateStatus = 2
			response.DueDateStatus = &dueDateStatus
		}

		difference := paymentDate.Sub(dueDate)
		aging := int64(difference.Hours() / 24)
		response.Aging = &aging
	}
	response.InvoiceStatusName = response.GenerateDataInvoiceStatusName()

	dueDateStatusName := response.GenerateDataDueDateStatusName()
	response.DueDateStatusName = &dueDateStatusName

	response.Details = DetailsData
	return response, nil
}
func (service *arServiceImpl) List(dataFilter entity.InvoiceQueryFilter) (data []entity.ArListResponse, total int64, lastPage int, err error) {
	invoices, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range invoices {
		var vResp entity.ArListResponse
		structs.Automapper(row, &vResp)
		if row.InvoiceDate != nil {
			invoiceDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = &invoiceDate
		}
		var dueDate time.Time
		if row.DueDate != nil {
			dueDate = *row.DueDate
			dueDateStr := dueDate.Format("2006-01-02")
			vResp.DueDate = &dueDateStr
		}

		invoicePaidAmount, errs := service.Repository.CountInvoicePaidAmount(row.InvoiceNo, dataFilter.CustId)
		if errs != nil {
			return data, total, lastPage, errs
		}
		vResp.PaidAmount = invoicePaidAmount.PaidAmount
		vResp.RemainingAmount = vResp.InvoiceAmount - vResp.PaidAmount
		vResp.InvoiceStatus = int64(entity.InvoiceStatusOutstanding)
		dueDateStatus := int64(1)
		vResp.DueDateStatus = &dueDateStatus
		paymentDate := time.Now()

		if vResp.RemainingAmount <= 0 {
			vResp.InvoiceStatus = int64(entity.InvoiceStatusPaid)
			vResp.RemainingAmount = 0
			vResp.DueDateStatus = nil
		} else {
			if paymentDate.Format("2006-01-02") >= dueDate.Format("2006-01-02") {
				dueDateStatus = 2
				vResp.DueDateStatus = &dueDateStatus
			}

			difference := paymentDate.Sub(dueDate)
			aging := int64(difference.Hours() / 24)
			vResp.Aging = &aging
		}
		vResp.InvoiceStatusName = vResp.GenerateDataInvoiceStatusName()

		dueDateStatusName := vResp.GenerateDataDueDateStatusName()
		vResp.DueDateStatusName = &dueDateStatusName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

/*
	func (service *arServiceImpl) Delete(custId string, arNo string, userId int64) (err error) {
		c := context.Background()
		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
			err = service.Repository.Delete(txCtx, custId, arNo, userId)
			if err != nil {
				return err
			}
			return nil
		})

		return err
	}

	func (service *arServiceImpl) Update(arNo string, request entity.UpdateArbody) (err error) {
		c := context.Background()

		if request.ArDate != nil {
			// parse time format YYYY-mm-dd to Rfc3339
			if *request.ArDate != "" {
				arDate, err := str.DateStrToRfc3339String(*request.ArDate)
				if err != nil {
					return err
				}
				request.ArDate = &arDate
			}

		}

		// End parse time format YYYY-mm-dd to Rfc339
		var Model model.Ar
		err = structs.Automapper(request, &Model)
		if err != nil {
			return err
		}
		if Model.IsPosted != nil {
			if *Model.IsPosted {
				now := time.Now()
				Model.PostedAt = &now
			}
		}
		Model.CustID = ""
		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
			err = service.Repository.Update(txCtx, arNo, Model)
			if err != nil {
				return err
			}
			DetailIds := []int64{}

			for _, detail := range request.Details {
				if detail.ArDetID != nil {
					DetailIds = append(DetailIds, *detail.ArDetID)
				}
			}
			if len(DetailIds) > 0 {
				err := service.Repository.DeleteDetailNotInIDs(txCtx, arNo, DetailIds)
				if err != nil {
					return err
				}
			}

			for _, detail := range request.Details {
				// parse time format YYYY-mm-dd to Rfc3339

				var arDetModel model.ArDet

				err = structs.Automapper(detail, &arDetModel)
				if err != nil {
					return err
				}
				arDetModel.CustID = request.CustID
				arDetModel.ArNo = arNo
				if detail.ArDetID == nil || *detail.ArDetID == 0 {
					arDetModel.ArDetID = nil
					err = service.Repository.StoreDetail(txCtx, &arDetModel)
					if err != nil {
						return err
					}
				} else {
					arDetModel.CustID = ""
					err = service.Repository.UpdateDetail(txCtx, &arDetModel)
					if err != nil {
						return err
					}

				}
			}
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	}
*/
func (service *arServiceImpl) StoreCollection(request entity.CreateCollectionBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.CollectionDate != nil {
		arDate, err := str.DateStrToRfc3339String(*request.CollectionDate)
		if err != nil {
			return err
		}
		request.CollectionDate = &arDate
	}

	var CollectionModel model.Collection
	err = structs.Automapper(request, &CollectionModel)
	if err != nil {
		return err
	}
	CollectionModel.UpdatedBy = CollectionModel.CreatedBy
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		totalRemainingAmount := 0.0

		log.Println("arService, StoreCollection, Collection Before StoreCollection:", CollectionModel)
		err := service.Repository.StoreCollection(txCtx, &CollectionModel)
		if err != nil {
			return err
		}

		for _, Detail := range request.Details {
			var detModel model.CollectionDet
			err = structs.Automapper(Detail, &detModel)
			if err != nil {
				return err
			}
			detModel.CustID = request.CustID
			detModel.CreatedBy = CollectionModel.CreatedBy
			detModel.CollectionNo = CollectionModel.CollectionNo

			detModel.RemainingAmount, err = service.calculateCollectionRemainingAmount(detModel.InvoiceNo, detModel.InvoiceAmount, detModel.PaidAmount, request.CustID)
			if err != nil {
				return err
			}
			totalRemainingAmount += detModel.RemainingAmount

			log.Println("arService, StoreCollection, Detail Before StoreCollectionDetail:", detModel)

			err = service.Repository.StoreCollectionDetail(txCtx, &detModel)
			if err != nil {
				return err
			}

		}

		if err := service.Repository.UpdateCollectionRemainingAmount(txCtx, CollectionModel.CollectionNo, request.CustID, totalRemainingAmount); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (service *arServiceImpl) CollectionDetail(collectionNo string, custID string, parentCustId string) (response entity.CollectionResponse, err error) {
	collectionDetail, err := service.Repository.FindCollectionByNo(collectionNo, custID, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(collectionDetail, &response)
	if err != nil {
		return response, err
	}

	Details, err := service.Repository.FindCollectionDetail(collectionNo, custID)
	if err != nil {
		return response, err
	}
	var DetailsData []entity.CollectionDetResponse
	for _, detail := range Details {
		var detailData entity.CollectionDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		if detail.InvoiceDate != nil {
			InvDate := detail.InvoiceDate.Format("2006-01-02")
			detailData.InvoiceDate = InvDate
		}

		if detail.DueDate != nil {
			DueDate := detail.DueDate.Format("2006-01-02")
			detailData.DueDate = DueDate
		}

		invoiceAmount, paidAmount, remainingAmount, invoicePayment, paidAmountByCollection := mapCollectionDetailAmounts(detail)

		invoicePaidAmount, err := service.Repository.CountInvoicePaidAmount(detail.InvoiceNo, custID)
		if err != nil {
			return response, err
		}

		invoicePayment = invoicePaidAmount.PaidAmount
		remainingAmount = invoiceAmount - invoicePayment
		if remainingAmount < 0 {
			remainingAmount = 0
		}

		detailData.InvoiceAmount = &invoiceAmount
		detailData.PaidAmount = &paidAmount
		detailData.RemainingAmount = &remainingAmount
		detailData.InvoicePayment = &invoicePayment
		detailData.PaidAmountByCollection = &paidAmountByCollection
		detailData.TotalInvoicePayment = &invoicePayment

		DetailsData = append(DetailsData, detailData)
	}

	if collectionDetail.CollectionDate != nil {
		collectionDate := collectionDetail.CollectionDate.Format("2006-01-02")
		response.CollectionDate = &collectionDate
	}

	if collectionDetail.InvoiceDateFrom != nil {
		invoiceDateFrom := collectionDetail.InvoiceDateFrom.Format("2006-01-02")
		response.InvoiceDateFrom = &invoiceDateFrom
	}

	if collectionDetail.InvoiceDateTo != nil {
		invoiceDateTo := collectionDetail.InvoiceDateTo.Format("2006-01-02")
		response.InvoiceDateTo = &invoiceDateTo
	}

	if collectionDetail.DueDateFrom != nil {
		dueDateFrom := collectionDetail.DueDateFrom.Format("2006-01-02")
		response.DueDateFrom = &dueDateFrom
	}

	if collectionDetail.DueDateTo != nil {
		dueDateTo := collectionDetail.DueDateTo.Format("2006-01-02")
		response.DueDateTo = &dueDateTo
	}

	totalInvoicePayment := calculateTotalInvoicePaymentFromResponse(DetailsData)
	totalRemainingAmount := calculateTotalRemainingAmountFromResponse(DetailsData)
	response.TotalInvoicePayment = &totalInvoicePayment
	response.RemainingAmount = &totalRemainingAmount
	response.Details = DetailsData
	return response, nil
}
func (service *arServiceImpl) CollectionList(dataFilter entity.CollectionQueryFilter) (data []entity.CollectionListResponse, total int64, lastPage int, err error) {
	arpays, total, lastPage, err := service.Repository.FindAllCollectionByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range arpays {
		var vResp entity.CollectionListResponse
		structs.Automapper(row, &vResp)
		if row.CollectionDate != nil {
			arPayDate := row.CollectionDate.Format("2006-01-02")
			vResp.CollectionDate = &arPayDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
func (service *arServiceImpl) DeleteCollection(custId string, collectionNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.DeleteCollection(txCtx, custId, collectionNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
func (service *arServiceImpl) UpdateCollection(collectionNo string, request entity.UpdateCollectionBody) (err error) {
	c := context.Background()

	if request.CollectionDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		if *request.CollectionDate != "" {
			collectionDate, err := str.DateStrToRfc3339String(*request.CollectionDate)
			if err != nil {
				return err
			}
			request.CollectionDate = &collectionDate
		}

	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.Collection
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = request.CustID
	Model.UpdatedAt = time.Now()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		totalRemainingAmount := 0.0

		err = service.Repository.UpdateCollection(txCtx, collectionNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details {
			if detail.CollectionDetID != nil {
				DetailIds = append(DetailIds, *detail.CollectionDetID)
			}
		}
		if len(DetailIds) > 0 {
			err := service.Repository.DeleteCollectionDetailNotInIDs(txCtx, collectionNo, DetailIds, request.CustID)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details {
			// parse time format YYYY-mm-dd to Rfc3339

			var collectionDetModel model.CollectionDet

			err = structs.Automapper(detail, &collectionDetModel)
			if err != nil {
				return err
			}
			collectionDetModel.CustID = request.CustID
			collectionDetModel.CollectionNo = collectionNo
			collectionDetModel.RemainingAmount, err = service.calculateCollectionRemainingAmount(collectionDetModel.InvoiceNo, collectionDetModel.InvoiceAmount, collectionDetModel.PaidAmount, request.CustID)
			if err != nil {
				return err
			}
			totalRemainingAmount += collectionDetModel.RemainingAmount
			if detail.CollectionDetID == nil || *detail.CollectionDetID == 0 {
				collectionDetModel.CollectionDetID = nil
				collectionDetModel.CreatedBy = Model.UpdatedBy
				err = service.Repository.StoreCollectionDetail(txCtx, &collectionDetModel)
				if err != nil {
					return err
				}
			} else {
				err = service.Repository.UpdateCollectionDetail(txCtx, &collectionDetModel)
				if err != nil {
					return err
				}

			}
		}

		if err := service.Repository.UpdateCollectionRemainingAmount(txCtx, collectionNo, request.CustID, totalRemainingAmount); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *arServiceImpl) ReplaceCollection(collectionNo string, request entity.UpdateCollectionBody) (err error) {
	c := context.Background()

	if request.CollectionDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		if *request.CollectionDate != "" {
			collectionDate, err := str.DateStrToRfc3339String(*request.CollectionDate)
			if err != nil {
				return err
			}
			request.CollectionDate = &collectionDate
		}

	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.Collection
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = request.CustID
	Model.UpdatedAt = time.Now()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		totalRemainingAmount := 0.0

		err = service.Repository.UpdateCollection(txCtx, collectionNo, Model)
		if err != nil {
			return err
		}

		// Delete ALL existing details for this collection
		err = service.Repository.DeleteAllCollectionDetails(txCtx, collectionNo, request.CustID)
		if err != nil {
			return err
		}

		for _, detail := range request.Details {
			var collectionDetModel model.CollectionDet

			err = structs.Automapper(detail, &collectionDetModel)
			if err != nil {
				return err
			}
			collectionDetModel.CustID = request.CustID
			collectionDetModel.CollectionNo = collectionNo

			// Always create new detail for replacement strategy
			collectionDetModel.CollectionDetID = nil
			collectionDetModel.CreatedBy = Model.UpdatedBy
			collectionDetModel.RemainingAmount, err = service.calculateCollectionRemainingAmount(collectionDetModel.InvoiceNo, collectionDetModel.InvoiceAmount, collectionDetModel.PaidAmount, request.CustID)
			if err != nil {
				return err
			}
			totalRemainingAmount += collectionDetModel.RemainingAmount

			err = service.Repository.StoreCollectionDetail(txCtx, &collectionDetModel)
			if err != nil {
				return err
			}
		}

		if err := service.Repository.UpdateCollectionRemainingAmount(txCtx, collectionNo, request.CustID, totalRemainingAmount); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *arServiceImpl) PrintCollection(custId string, collectionNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.PrintCollection(txCtx, custId, collectionNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *arServiceImpl) EmployeeGroupLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.EmployeeGroupLookupResponse, total int64, lastPage int, err error) {
	EmployeeGroups, total, lastPage, err := service.Repository.FindAllEmployeeGroupByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range EmployeeGroups {
		var vResp entity.EmployeeGroupLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *arServiceImpl) OutletGroupFilterLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.OutletGroupLookupResponse, total int64, lastPage int, err error) {
	OutletGroups, total, lastPage, err := service.Repository.FindAllOutletGroupFilterByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range OutletGroups {
		var vResp entity.OutletGroupLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *arServiceImpl) EmployeeLookupList(dataFilter entity.EmployeeListQueryFilter) (data []entity.EmployeeLookupResponse, total int64, lastPage int, err error) {
	Employees, total, lastPage, err := service.Repository.FindAllEmployeeByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Employees {
		var vResp entity.EmployeeLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *arServiceImpl) SalesmanFilterLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.SalesmanLookupResponse, total int64, lastPage int, err error) {
	Salesmans, total, lastPage, err := service.Repository.FindAllSalesmanFilterByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Salesmans {
		var vResp entity.SalesmanLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *arServiceImpl) InvoiceList(dataFilter entity.InvoiceQueryFilter) (data []entity.InvoiceListResponse, total int64, lastPage int, err error) {
	realOrders, total, lastPage, err := service.Repository.FindAllInvoiceByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range realOrders {
		var vResp entity.InvoiceListResponse
		structs.Automapper(row, &vResp)
		if row.InvoiceDate != nil {
			invoiceDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = &invoiceDate
		}

		if row.DueDate != nil {
			dueDate := row.DueDate.Format("2006-01-02")
			vResp.DueDate = &dueDate
		}

		// invoicePaidAmount, errs := service.Repository.CountInvoicePaidAmount(*row.InvoiceNo, dataFilter.CustId)
		// if errs != nil {
		// 	return data, total, lastPage, errs
		// }
		// vResp.PaidAmount = &invoicePaidAmount.PaidAmount
		// vResp.RemainingAmount = new(float64)                              // Allocate memory for RemainingAmount
		// *vResp.RemainingAmount = *vResp.InvoiceAmount - *vResp.PaidAmount // Dereference pointers for calculation

		// if vResp.RemainingAmount != nil && *vResp.RemainingAmount <= 0 {
		// 	*vResp.RemainingAmount = 0
		// }

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *arServiceImpl) OutletFilterLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.OutletLookupResponse, total int64, lastPage int, err error) {
	Outlets, total, lastPage, err := service.Repository.FindAllOutletFilterByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Outlets {
		var vResp entity.OutletLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *arServiceImpl) CollectorLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.EmployeeLookupResponse, total int64, lastPage int, err error) {
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

func mapCollectionDetailAmounts(detail model.CollectionDetList) (invoiceAmount float64, paidAmount float64, remainingAmount float64, invoicePayment float64, paidAmountByCollection float64) {
	invoiceAmount = safeFloat64(detail.InvoiceAmount)
	paidAmount = safeFloat64(detail.PaidAmount)
	remainingAmount = safeFloat64(detail.RemainingAmount)
	invoicePayment = safeFloat64(detail.TotalInvoicePayment)
	paidAmountByCollection = paidAmount
	return invoiceAmount, paidAmount, remainingAmount, invoicePayment, paidAmountByCollection
}

func calculateTotalInvoicePayment(details []model.CollectionDetList) float64 {
	totalInvoicePayment := 0.0
	for _, detail := range details {
		totalInvoicePayment += safeFloat64(detail.TotalInvoicePayment)
	}
	return totalInvoicePayment
}

func calculateTotalInvoicePaymentFromResponse(details []entity.CollectionDetResponse) float64 {
	totalInvoicePayment := 0.0
	for _, detail := range details {
		totalInvoicePayment += safeFloat64(detail.InvoicePayment)
	}
	return totalInvoicePayment
}

func calculateTotalRemainingAmountFromResponse(details []entity.CollectionDetResponse) float64 {
	totalRemainingAmount := 0.0
	for _, detail := range details {
		totalRemainingAmount += safeFloat64(detail.RemainingAmount)
	}
	return totalRemainingAmount
}

func safeFloat64(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func (service *arServiceImpl) calculateCollectionRemainingAmount(invoiceNo string, invoiceAmount float64, paidAmountByCollection float64, custID string) (float64, error) {
	invoicePaidAmount, err := service.Repository.CountInvoicePaidAmount(invoiceNo, custID)
	if err != nil {
		return 0, err
	}

	remainingAmount := invoiceAmount - (invoicePaidAmount.PaidAmount + paidAmountByCollection)
	if remainingAmount < 0 {
		return 0, nil
	}

	return remainingAmount, nil
}
