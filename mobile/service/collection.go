package service

import (
	"context"
	"sort"
	"time"

	"mobile/model"
	"mobile/pkg/constant"
	"mobile/pkg/str"

	"mobile/entity"
	"mobile/pkg/config/env"

	"fmt"
	"mobile/pkg/structs"
	"mobile/repository"
)

type CollectionService interface {
	StoreCollection(request entity.CreateDepositBodyByCollection) (err error)
	StoreCollectionNoPayment(request entity.CreateCollectionNoPaymentRequest) (err error)
	List(dataFilter entity.CollectionQueryFilter) (data []entity.CollectionListResponse, total model.CollectionTotal, lastPage int, err error)
	CollectionList(dataFilter entity.CollectionQueryFilter) (data []entity.CollectionListV2Response, total model.CollectionTotal, lastPage int, err error)
	Detail(depositNo string, custID string) (response entity.DepositDetailResponse, err error)
	DetailInvoice(invoiceNo string, custID string) (response entity.DepositInvoiceDetailResponse, err error)

	ListMissedPayment(dataFilter entity.GeneralQueryFilter) (responses []entity.MissedPaymentReasonResp, err error)
	StoreCollectionList(req entity.CreateCollectionListBody) (err error)
}

func NewCollectionService(
	config env.ConfigEnv,
	collectionRepository repository.CollectionRepository,
	transaction repository.Dbtransaction,
	orderRepo repository.OrderRepository) *collectionServiceImpl {
	return &collectionServiceImpl{
		Config:               config,
		CollectionRepository: collectionRepository,
		Transaction:          transaction,
		OrderRepository:      orderRepo,
	}
}

type collectionServiceImpl struct {
	Config               env.ConfigEnv
	CollectionRepository repository.CollectionRepository
	Transaction          repository.Dbtransaction
	OrderRepository      repository.OrderRepository
}

func (service *collectionServiceImpl) List(dataFilter entity.CollectionQueryFilter) (data []entity.CollectionListResponse, total model.CollectionTotal, lastPage int, err error) {
	arpays, total, lastPage, err := service.CollectionRepository.FindAllByCustId(dataFilter)
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

func (service *collectionServiceImpl) StoreCollection(request entity.CreateDepositBodyByCollection) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.DepositDate != "" {
		DepositDate, err := str.DateStrToRfc3339String(request.DepositDate)
		if err != nil {
			return err
		}
		request.DepositDate = DepositDate
	}

	var depositModel model.Deposit
	err = structs.Automapper(request, &depositModel)
	if err != nil {
		return err
	}

	total, err := service.CollectionRepository.CountAllByCustId(request.CustID, request.DepositDate)
	if err != nil {
		return err
	}

	depositNo := entity.GenerateNumber("DP", total, depositModel.DepositDate)
	depositModel.DepositNo = depositNo

	fmt.Println(depositNo)

	if depositModel.CollectionNo != nil {
		depositModel.SalesmanID = nil
		depositModel.InvoiceDateFrom = nil
		depositModel.InvoiceDateTo = nil
		depositModel.DueDateTo = nil
		depositModel.DueDateTo = nil
	}

	depositModel.DepositStatus = 1

	var remainingAmount = 0.0

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		for _, detail := range request.Details {
			var depositDetails model.DepositDetail

			err := structs.Automapper(detail, &depositDetails)
			if err != nil {
				return err
			}

			depositDetails.DepositNo = depositNo
			remainingAmountByInv, errR := service.CollectionRepository.CountRemainingAmountByInvoice(txCtx, detail.InvoiceNo, request.CustID)
			if errR != nil {
				return errR
			}

			depositDetails.RemainingPayment = detail.InvoiceAmount - (remainingAmountByInv + depositDetails.TotalPayment)
			remainingAmount += depositDetails.RemainingPayment
			depositDetails.CustID = request.CustID
			_, err = service.CollectionRepository.StoreDetail(txCtx, &depositDetails)
			if err != nil {
				return err
			}

			indexGiro := -1
			for i, payment := range detail.Payment {
				if payment.PayType == 2 {
					indexGiro = i
					break
				}
			}

			for index, payment := range detail.Payment {
				var depositPayment model.DepositPayment

				err := structs.Automapper(payment, &depositPayment)
				if err != nil {
					return err
				}
				depositPayment.DepositNo = depositNo
				depositPayment.CustID = request.CustID

				if index == 0 {
					depositPayment.Discount = &detail.Discount
				}

				if indexGiro > -1 && indexGiro == index {
					depositPayment.Materai = &detail.Materai
				} else {
					zeroFloat := float64(0)
					depositPayment.Materai = &zeroFloat
				}

				_, err = service.CollectionRepository.StorePayment(txCtx, &depositPayment)
				if err != nil {
					return err
				}

				for _, image := range payment.Images {
					depositPaymentImage := model.DepositPaymentImage{
						DepositNo: depositNo,
						InvoiceNo: payment.InvoiceNo,
						ImageUrl:  image.ImageUrl,
					}

					_, err = service.CollectionRepository.StoreDepositPaymentImage(txCtx, &depositPaymentImage)
					if err != nil {
						return err
					}
				}
			}
		}

		depositModel.RemainingAmount = remainingAmount
		err := service.CollectionRepository.Store(txCtx, &depositModel)

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

func (service *collectionServiceImpl) StoreCollectionNoPayment(request entity.CreateCollectionNoPaymentRequest) error {
	c := context.Background()

	// Convert request to model
	if request.PaymentDate != "" {
		PaymentDate, err := str.DateStrToRfc3339String(request.PaymentDate)
		if err != nil {
			return err
		}
		request.PaymentDate = PaymentDate
	}
	var collectionNoPaymentModel model.CollectionNoPayment
	err := structs.Automapper(request, &collectionNoPaymentModel)
	if err != nil {
		return err
	}

	// Store in database
	err = service.CollectionRepository.StoreCollectionNoPayment(c, &collectionNoPaymentModel)
	if err != nil {
		return err
	}

	return nil
}

func (service *collectionServiceImpl) Detail(depositNo string, custID string) (response entity.DepositDetailResponse, err error) {
	Deposit, err := service.CollectionRepository.FindByNo(depositNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(Deposit, &response)
	if err != nil {
		return response, err
	}

	statusText := entity.ConvStatus(entity.StatusDeposit, response.DepositStatus)
	response.DepositStatusName = statusText

	details, err := service.CollectionRepository.FindDetailByNo(depositNo, custID)
	if err != nil {
		return response, err
	}

	response.Details = make([]entity.DepositDetail, len(details))
	for i, detail := range details {
		var detailData entity.DepositDetail
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		detailsPayment, err := service.CollectionRepository.FindDetailPaymentByNo(depositNo, detailData.InvoiceNo, custID)
		if err != nil {
			return response, err
		}

		detailData.Payment = make([]entity.DepositPayment, len(detailsPayment))
		for i, detail := range detailsPayment {
			var detailDataPayment entity.DepositPayment
			err = structs.Automapper(detail, &detailDataPayment)
			if err != nil {
				return response, err
			}
			// Fetch payment images for this payment
			paymentImages, err := service.CollectionRepository.FindPaymentImagesByNo(depositNo, detail.InvoiceNo)
			if err != nil {
				return response, err
			}

			// Map payment images to the DepositPayment
			detailDataPayment.Images = make([]entity.DepositPaymentImage, len(paymentImages))
			for j, img := range paymentImages {
				var imageData entity.DepositPaymentImage
				err = structs.Automapper(img, &imageData)
				if err != nil {
					return response, err
				}
				detailDataPayment.Images[j] = imageData
			}
			detailData.Payment[i] = detailDataPayment
		}

		response.Details[i] = detailData
	}

	if Deposit.DepositDate != nil {
		DepositDate := Deposit.DepositDate.Format("2006-01-02")
		response.DepositDate = &DepositDate
	}

	if Deposit.CollectionDate != nil {
		CollectionDate := Deposit.CollectionDate.Format("2006-01-02")
		response.CollectionDate = &CollectionDate
	}

	if Deposit.InvoiceDateFrom != nil {
		InvoiceDateFrom := Deposit.InvoiceDateFrom.Format("2006-01-02")
		response.InvoiceDateFrom = &InvoiceDateFrom
	}

	if Deposit.InvoiceDateTo != nil {
		InvoiceDateTo := Deposit.InvoiceDateTo.Format("2006-01-02")
		response.InvoiceDateTo = &InvoiceDateTo
	}

	if Deposit.DueDateFrom != nil {
		DueDateFrom := Deposit.DueDateFrom.Format("2006-01-02")
		response.DueDateFrom = &DueDateFrom
	}

	if Deposit.DueDateTo != nil {
		DueDateTo := Deposit.DueDateTo.Format("2006-01-02")
		response.DueDateTo = &DueDateTo
	}

	cashs, err := service.CollectionRepository.FindDetailPaymentInvoiceByNo(1, depositNo, custID)
	if err != nil {
		return response, err
	}

	response.Cash = make([]entity.DepositPaymentInvoice, len(cashs))

	for i, cash := range cashs {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(cash, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		detailData.PaymentBalance = Deposit.TotalPaymentBalance

		response.Cash[i] = detailData
	}

	cek, err := service.CollectionRepository.FindDetailPaymentInvoiceByNo(2, depositNo, custID)
	if err != nil {
		return response, err
	}

	response.Cek = make([]entity.DepositPaymentInvoice, len(cek))

	for i, ceks := range cek {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(ceks, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		detailData.PaymentBalance = Deposit.TotalPaymentBalance

		response.Cek[i] = detailData
	}

	transfers, err := service.CollectionRepository.FindDetailPaymentInvoiceByNo(3, depositNo, custID)
	if err != nil {
		return response, err
	}

	response.Trasfer = make([]entity.DepositPaymentInvoice, len(transfers))

	for i, transfer := range transfers {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(transfer, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		detailData.PaymentBalance = Deposit.TotalPaymentBalance

		response.Trasfer[i] = detailData
	}

	returns, err := service.CollectionRepository.FindDetailPaymentInvoiceByNo(4, depositNo, custID)
	if err != nil {
		return response, err
	}

	response.Return = make([]entity.DepositPaymentInvoice, len(returns))

	for i, returna := range returns {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(returna, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		detailData.PaymentBalance = Deposit.TotalPaymentBalance

		response.Return[i] = detailData
	}

	cndns, err := service.CollectionRepository.FindDetailPaymentInvoiceByNo(5, depositNo, custID)
	if err != nil {
		return response, err
	}

	response.CNDN = make([]entity.DepositPaymentInvoice, len(cndns))

	for i, cndn := range cndns {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(cndn, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		detailData.PaymentBalance = Deposit.TotalPaymentBalance

		response.CNDN[i] = detailData
	}

	// ownerName := entity.ConvStatus(entity.OwnerGiro, response.OwnerID)
	// response.OwnerName = ownerName

	// statusText := entity.ConvStatus(entity.StatusGiro, response.StatusCheque)
	// response.StatusChequeText = &statusText

	// response.UsedAmount = float64(0)
	// response.RemainingAmount = response.Amount - response.UsedAmount

	return response, nil
}

func (service *collectionServiceImpl) DetailInvoice(invoiceNo string, custID string) (response entity.DepositInvoiceDetailResponse, err error) {
	cashs, err := service.CollectionRepository.FindDetailPaymentByInvoice(1, invoiceNo, custID)
	if err != nil {
		return response, err
	}

	response.Cash = make([]entity.DepositPaymentInvoice, len(cashs))

	for i, cash := range cashs {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(cash, &detailData)
		if err != nil {
			return response, err
		}

		// Calculate TotalPayment based on the fields: PaymentAmount, Materai, and Discount
		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		// detailData.PaymentBalance = Deposit.TotalPaymentBalance

		paymentImages, err := service.CollectionRepository.FindPaymentImagesByNo(cash.DepositNo, cash.InvoiceNo)
		if err != nil {
			return response, err
		}

		detailData.Images = make([]entity.DepositPaymentImage, len(paymentImages))
		for j, img := range paymentImages {
			var imageData entity.DepositPaymentImage
			err = structs.Automapper(img, &imageData)
			if err != nil {
				return response, err
			}
			detailData.Images[j] = imageData
		}

		response.Cash[i] = detailData
	}

	cek, err := service.CollectionRepository.FindDetailPaymentByInvoice(2, invoiceNo, custID)
	if err != nil {
		return response, err
	}

	response.Cek = make([]entity.DepositPaymentInvoice, len(cek))

	for i, ceks := range cek {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(ceks, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		// detailData.PaymentBalance = Deposit.TotalPaymentBalance
		paymentImages, err := service.CollectionRepository.FindPaymentImagesByNo(ceks.DepositNo, ceks.InvoiceNo)
		if err != nil {
			return response, err
		}

		detailData.Images = make([]entity.DepositPaymentImage, len(paymentImages))
		for j, img := range paymentImages {
			var imageData entity.DepositPaymentImage
			err = structs.Automapper(img, &imageData)
			if err != nil {
				return response, err
			}
			detailData.Images[j] = imageData
		}

		response.Cek[i] = detailData
	}

	transfers, err := service.CollectionRepository.FindDetailPaymentByInvoice(3, invoiceNo, custID)
	if err != nil {
		return response, err
	}

	response.Trasfer = make([]entity.DepositPaymentInvoice, len(transfers))

	for i, transfer := range transfers {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(transfer, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		// detailData.PaymentBalance = Deposit.TotalPaymentBalance
		paymentImages, err := service.CollectionRepository.FindPaymentImagesByNo(transfer.DepositNo, transfer.InvoiceNo)
		if err != nil {
			return response, err
		}

		detailData.Images = make([]entity.DepositPaymentImage, len(paymentImages))
		for j, img := range paymentImages {
			var imageData entity.DepositPaymentImage
			err = structs.Automapper(img, &imageData)
			if err != nil {
				return response, err
			}
			detailData.Images[j] = imageData
		}

		response.Trasfer[i] = detailData
	}

	returns, err := service.CollectionRepository.FindDetailPaymentByInvoice(4, invoiceNo, custID)
	if err != nil {
		return response, err
	}

	response.Return = make([]entity.DepositPaymentInvoice, len(returns))

	for i, returna := range returns {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(returna, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		// detailData.PaymentBalance = Deposit.TotalPaymentBalance
		paymentImages, err := service.CollectionRepository.FindPaymentImagesByNo(returna.DepositNo, returna.InvoiceNo)
		if err != nil {
			return response, err
		}

		detailData.Images = make([]entity.DepositPaymentImage, len(paymentImages))
		for j, img := range paymentImages {
			var imageData entity.DepositPaymentImage
			err = structs.Automapper(img, &imageData)
			if err != nil {
				return response, err
			}
			detailData.Images[j] = imageData
		}

		response.Return[i] = detailData
	}

	cndns, err := service.CollectionRepository.FindDetailPaymentByInvoice(5, invoiceNo, custID)
	if err != nil {
		return response, err
	}

	response.CNDN = make([]entity.DepositPaymentInvoice, len(cndns))

	for i, cndn := range cndns {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(cndn, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		// detailData.PaymentBalance = Deposit.TotalPaymentBalance
		paymentImages, err := service.CollectionRepository.FindPaymentImagesByNo(cndn.DepositNo, cndn.InvoiceNo)
		if err != nil {
			return response, err
		}

		detailData.Images = make([]entity.DepositPaymentImage, len(paymentImages))
		for j, img := range paymentImages {
			var imageData entity.DepositPaymentImage
			err = structs.Automapper(img, &imageData)
			if err != nil {
				return response, err
			}
			detailData.Images[j] = imageData
		}

		response.CNDN[i] = detailData
	}

	// ownerName := entity.ConvStatus(entity.OwnerGiro, response.OwnerID)
	// response.OwnerName = ownerName

	// statusText := entity.ConvStatus(entity.StatusGiro, response.StatusCheque)
	// response.StatusChequeText = &statusText

	// response.UsedAmount = float64(0)
	// response.RemainingAmount = response.Amount - response.UsedAmount

	return response, nil
}

func (service *collectionServiceImpl) ListMissedPayment(dataFilter entity.GeneralQueryFilter) (responses []entity.MissedPaymentReasonResp, err error) {
	missedPayments, err := service.CollectionRepository.FindMissedPaymentReasons(dataFilter)
	if err != nil {
		return
	}

	responses = make([]entity.MissedPaymentReasonResp, 0)
	for _, missed_payment := range missedPayments {
		orderResp := entity.MissedPaymentReasonResp{
			MissedPaymentId:   missed_payment.MissedPaymentId,
			MissedPaymentName: missed_payment.MissedPaymentName,
			ImageUrl:          missed_payment.ImageUrl,
		}
		responses = append(responses, orderResp)
	}
	return
}

func (service *collectionServiceImpl) CollectionList(dataFilter entity.CollectionQueryFilter) (data []entity.CollectionListV2Response, total model.CollectionTotal, lastPage int, err error) {
	arpays, total, lastPage, err := service.CollectionRepository.GetCollectionList(context.Background(), dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range arpays {
		vResp := entity.CollectionListV2Response{
			CollectionNo:    row.CollectionNo,
			EmpID:           row.EmpID,
			InvoiceNo:       row.InvoiceNo,
			InvoiceAmount:   row.TotalAmount,
			PaidAmount:      row.PaidAmount,
			RemainingAmount: row.RemainingAmount,
			RoNo:            row.RoNo,
			OrderNo:         row.OrderNo,
		}
		if row.InvoiceDateFrom != nil {
			strInvoiceDate := row.InvoiceDateFrom.Format("2006-01-02")
			vResp.InvoiceDate = &strInvoiceDate
		}
		if row.DueDateFrom != nil {
			strDueDate := row.DueDateFrom.Format("2006-01-02")
			vResp.DueDate = &strDueDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *collectionServiceImpl) StoreCollectionList(req entity.CreateCollectionListBody) (err error) {
	c := context.Background()
	var collectionListModel model.CollectionModel
	err = structs.Automapper(req, &collectionListModel)
	if err != nil {
		return err
	}

	collectionNo, err := service.CollectionRepository.GetNewCollectionNo()
	if err != nil {
		return err
	}

	invoices, err := service.OrderRepository.GetInvoiceByNumbers(req.CustID, req.ParentCustID, req.EmpID, req.IsCollection)
	if err != nil {
		return err
	}

	if len(invoices) == 0 {
		return nil
	}

	now := time.Now()

	var details []model.CollectionDetail
	var totalAmount float64
	var totalRemaining float64
	var dueDates []time.Time
	var invoiceDates []time.Time
	for _, invoice := range invoices {
		totalAmount += invoice.InvoiceAmount
		totalRemaining += invoice.RemainingAmount
		if invoice.DueDate != nil {
			dueDates = append(dueDates, *invoice.DueDate)
		}
		if invoice.InvoiceDate != nil {
			invoiceDates = append(invoiceDates, *invoice.InvoiceDate)
		}
		detail := model.CollectionDetail{
			CustID:          req.CustID,
			CollectionNo:    collectionNo,
			InvoiceNo:       invoice.InvoiceNo,
			SalesmanID:      req.EmpID,
			InvoiceAmount:   invoice.InvoiceAmount,
			RemainingAmount: invoice.RemainingAmount,
			PaidAmount:      invoice.PaidAmount,
			CreatedBy:       &req.UserID,
			CreatedAt:       &now,
			Source:          constant.SourceMobile,
		}
		details = append(details, detail)
	}

	collectionListModel.CustID = req.CustID
	collectionListModel.CollectionNo = collectionNo
	collectionListModel.CollectionDate = &now
	collectionListModel.EmpID = &req.EmpID
	collectionListModel.TotalAmount = totalAmount
	collectionListModel.RemainingAmount = totalRemaining
	collectionListModel.CreatedBy = &req.UserID
	collectionListModel.UpdatedBy = &req.UserID
	collectionListModel.CreatedAt = now
	collectionListModel.UpdatedAt = now
	collectionListModel.Source = constant.SourceMobile

	if len(dueDates) > 0 {
		sort.Slice(dueDates, func(i, j int) bool { return dueDates[i].Before(dueDates[j]) })
		collectionListModel.DueDateFrom = &dueDates[0]
		collectionListModel.DueDateTo = &dueDates[len(dueDates)-1]
	}

	if len(invoiceDates) > 0 {
		sort.Slice(invoiceDates, func(i, j int) bool { return invoiceDates[i].Before(invoiceDates[j]) })
		collectionListModel.InvoiceDateFrom = &invoiceDates[0]
		collectionListModel.InvoiceDateTo = &invoiceDates[len(invoiceDates)-1]
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		if err := service.CollectionRepository.StoreCollection(txCtx, &collectionListModel); err != nil {
			return err
		}
		if err := service.CollectionRepository.StoreCollectionDetails(txCtx, details); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
