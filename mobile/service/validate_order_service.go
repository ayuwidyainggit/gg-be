package service

import (
	"fmt"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/conversion"
	"mobile/pkg/structs"
	"mobile/repository"
	"time"
)

type ValidateOrderService interface {
	ValidateOrder(dataFilter entity.ValidateOrderBody) (data entity.ValidateResponse, total int64, lastPage int, err error)
	ValidateOrderDetail(dataFilter entity.ValidateOrderDetailBody) (data entity.ValidateDetailResponse, total int64, lastPage int, err error)
}

func NewValidateOrderService(validateOrderRepository repository.ValidateOrderRepository, transaction repository.Dbtransaction) *validateOrderServiceImpl {
	return &validateOrderServiceImpl{
		ValidateOrderRepository: validateOrderRepository,
		Transaction:             transaction,
	}
}

type validateOrderServiceImpl struct {
	ValidateOrderRepository repository.ValidateOrderRepository
	Transaction             repository.Dbtransaction
}

func (service *validateOrderServiceImpl) ValidateOrder(dataFilter entity.ValidateOrderBody) (data entity.ValidateResponse, total int64, lastPage int, err error) {

	var vRespvald entity.ValidateResponse

	vRespvald.IsSuccessValidate = false
	vRespvald.Validate1 = "Sufficient Stock"
	vRespvald.Validate2 = "Within Limit"
	vRespvald.Validate3 = "Allowed"
	vRespvald.Validate4 = "Allowed"

	dataFilter.ActiveProductOnly = "true"

	dataFilter.Date = time.Now().Format("2006-01-02")

	// fmt.Println(dataFilter.Date)

	for _, product := range dataFilter.ProStok {
		// fmt.Println("pro id", product.ProductId)
		dataFilter.ProID = append(dataFilter.ProID, product.ProductId)
	}

	stocks, total, lastPage, err := service.ValidateOrderRepository.StockReport(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	productsModel, err := service.ValidateOrderRepository.FindProductByListID(dataFilter.ProID)
	// if err != nil {
	// 	return err
	// }

	var productMap = model.MapProduct{}

	for _, productModel := range productsModel {
		productMap.SetProduct(productModel.ProductId, productModel)
	}

	//validate stock

	var cekValidate = true

	for _, row := range stocks {
		productModel, _ := productMap.GetByID(int64(row.ProID))
		// if err != nil {
		// 	return err
		// }

		// fmt.Print("row qty", int(row.Qty))

		cekValidate = true

		for _, product := range dataFilter.ProStok {
			dataFilter.ProID = append(dataFilter.ProID, product.ProductId)

			QtyUnit := &conversion.QtyUnit{
				Qty1:      int(product.Qty1),
				Qty2:      int(product.Qty2),
				Qty3:      int(product.Qty3),
				ConvUnit2: int(productModel.ConvUnit2),
				ConvUnit3: int(productModel.ConvUnit3),
			}

			totalQty, _ := QtyUnit.ToTotalQuantity()
			// fmt.Print(totalQty, "<", row.Qty, "/n")
			if totalQty > int(row.Qty) {
				cekValidate = false
				vRespvald.Validate1 = "Insufficient Stock"
			}
			// if err != nil {
			// 	return err
			// }
		}

		vRespvald.Validate1Success = cekValidate
	}

	invoices, total, lastPage, err := service.ValidateOrderRepository.FindAllArByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	var remainingLimit = 0.0
	var remainingObs = 0
	var dueDateCount = 0

	for _, row := range invoices {
		var vResp entity.ArListResponse
		structs.Automapper(row, &vResp)

		var dueDate time.Time
		if row.DueDate != nil {
			dueDate = *row.DueDate
			dueDateStr := dueDate.Format("2006-01-02")
			vResp.DueDate = &dueDateStr
		}

		invoicePaidAmount, errs := service.ValidateOrderRepository.CountInvoicePaidAmount(row.InvoiceNo, dataFilter.CustID)
		if errs != nil {
			return data, total, lastPage, errs
		}
		vResp.PaidAmount = invoicePaidAmount.PaidAmount
		vResp.RemainingAmount = vResp.InvoiceAmount - vResp.PaidAmount
		remainingLimit = remainingLimit + vResp.RemainingAmount
		dueDateStatus := int64(1)
		vResp.DueDateStatus = &dueDateStatus
		paymentDate := time.Now()

		if vResp.RemainingAmount > 0 {
			remainingObs++
		}

		if vResp.RemainingAmount > 0 && paymentDate.Format("2006-01-02") >= dueDate.Format("2006-01-02") {
			dueDateStatus = 2
			dueDateCount++
			cekValidate = false
		}
	}

	outlet, _ := service.ValidateOrderRepository.FindDetailOutletID(dataFilter.OutletID)
	// if err != nil {
	// 	return data, total, lastPage, err
	// }

	// fmt.Println(outlet)
	if *outlet.SalesInvLimitType == 1 {
		vRespvald.Validate3Success = true
		vRespvald.Validate3 = "Allowed (Unlimited)"
	} else {
		allowed := "Not Allowed"

		if dueDateCount >= *outlet.SalesInvLimit {
			vRespvald.Validate3Success = false
			allowed = "Not Allowed"

		} else {
			vRespvald.Validate3Success = true
			allowed = "Allowed"
		}
		vRespvald.Validate3 = allowed + " (" + fmt.Sprintf("%d", dueDateCount) + " of " + fmt.Sprintf("%d", *outlet.SalesInvLimit) + ")"
	}

	if *outlet.CreditLimitType == 1 {
		vRespvald.Validate2Success = true
		vRespvald.Validate2 = "Allowed (Unlimited)"
	} else {
		if remainingLimit+dataFilter.Total <= *outlet.CreditLimit {
			vRespvald.Validate2Success = true
		} else {
			// vRespvald.Validate2 = "Over Limit, (" + fmt.Sprintf("%.2f", remainingLimit+dataFilter.Total) + " of " + fmt.Sprintf("%.2f", *outlet.CreditLimit) + ")"
			vRespvald.Validate2 = "Over Limit"
		}
	}

	if *outlet.ObsType == 1 {
		vRespvald.Validate4Success = true
		vRespvald.Validate4 = "Allowed (Unlimited)"
	} else {
		allowed := "Not Allowed"
		if remainingObs <= *outlet.Obs {
			vRespvald.Validate4Success = true
			allowed = "Allowed"
		}
		vRespvald.Validate4 = allowed + " (" + fmt.Sprintf("%d", remainingObs) + " of " + fmt.Sprintf("%d", *outlet.Obs) + ")"
	}

	if vRespvald.Validate1Success && vRespvald.Validate2Success && vRespvald.Validate3Success && vRespvald.Validate4Success {
		vRespvald.IsSuccessValidate = true
	}

	return vRespvald, total, lastPage, err
}

func (service *validateOrderServiceImpl) ValidateOrderDetail(dataFilter entity.ValidateOrderDetailBody) (data entity.ValidateDetailResponse, total int64, lastPage int, err error) {

	var vRespvald entity.ValidateResponse

	vRespvald.IsSuccessValidate = false
	vRespvald.Validate3 = "Allowed"
	vRespvald.Validate4 = "Allowed"

	// dataFilter.ActiveProductOnly = "true"

	dataFilter.Date = time.Now().Format("2006-01-02")

	// fmt.Println(dataFilter.Date)

	invoices, total, lastPage, err := service.ValidateOrderRepository.FindAllArDetailByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	var remainingLimit = 0.0
	var remainingObs = 0
	var dueDateCount = 0

	var dueDateArray []entity.ArListResponse
	var outstandingArray []entity.ArListResponse
	var responseValidate entity.ValidateDetailResponse

	for _, row := range invoices {
		var vResp entity.ArListResponse
		structs.Automapper(row, &vResp)

		var dueDate time.Time
		if row.DueDate != nil {
			dueDate = *row.DueDate
			dueDateStr := dueDate.Format("2006-01-02")
			vResp.DueDate = &dueDateStr
		}

		invoicePaidAmount, errs := service.ValidateOrderRepository.CountInvoicePaidAmount(row.InvoiceNo, dataFilter.CustID)
		if errs != nil {
			return data, total, lastPage, errs
		}
		vResp.PaidAmount = invoicePaidAmount.PaidAmount
		vResp.RemainingAmount = vResp.InvoiceAmount - vResp.PaidAmount
		remainingLimit = remainingLimit + vResp.RemainingAmount
		dueDateStatus := int64(1)
		vResp.DueDateStatus = &dueDateStatus
		paymentDate := time.Now()

		if vResp.RemainingAmount > 0 && paymentDate.Format("2006-01-02") < dueDate.Format("2006-01-02") {
			remainingObs++
			difference := paymentDate.Sub(dueDate)
			aging := int64(difference.Hours() / 24)
			vResp.Aging = &aging
			vResp.InvoiceStatusName = "Outstanding"
			outstandingArray = append(dueDateArray, vResp) // Tambahkan vResp ke array
		}

		if vResp.RemainingAmount > 0 && paymentDate.Format("2006-01-02") >= dueDate.Format("2006-01-02") {
			dueDateStatus = 2
			dueDateCount++
			difference := paymentDate.Sub(dueDate)
			aging := int64(difference.Hours() / 24)
			vResp.Aging = &aging
			statusName := "Overdue"
			vResp.DueDateStatusName = &statusName
			dueDateArray = append(dueDateArray, vResp) // Tambahkan vResp ke array
		}
	}

	responseValidate.DuedateInvoive = dueDateArray
	responseValidate.OutstandingInvoice = outstandingArray

	return responseValidate, total, lastPage, err
}
