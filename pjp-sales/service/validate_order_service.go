package service

import (
	"fmt"
	"sales/entity"
	"sales/pkg/conversion"
	"sales/pkg/structs"
	"sales/repository"
	"strings"
	"time"
)

type ValidateOrderService interface {
	ValidateOrder(dataFilter entity.ValidateOrderBody) (data entity.ValidateResponse, total int64, lastPage int, err error)
	ValidateOrderWithoutStock(dataFilter entity.ValidateOrderBody) (data entity.ValidateResponse, total int64, lastPage int, err error)
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

func formatCurrency(amount float64) string {
	str := fmt.Sprintf("%.2f", amount)
	parts := strings.Split(str, ".")
	intPart := parts[0]
	//decPart := parts[1]

	// Format pemisah ribuan (pakai titik)
	var result []byte
	count := 0
	for i := len(intPart) - 1; i >= 0; i-- {
		result = append([]byte{intPart[i]}, result...)
		count++
		if count%3 == 0 && i != 0 {
			result = append([]byte{'.'}, result...)
		}
	}

	// Gabungkan dan ubah desimal menjadi koma
	return fmt.Sprintf("%s", string(result))
}

func (service *validateOrderServiceImpl) ValidateOrder(dataFilter entity.ValidateOrderBody) (data entity.ValidateResponse, total int64, lastPage int, err error) {
	return service.validateOrder(dataFilter, true)
}

func (service *validateOrderServiceImpl) ValidateOrderWithoutStock(dataFilter entity.ValidateOrderBody) (data entity.ValidateResponse, total int64, lastPage int, err error) {
	return service.validateOrder(dataFilter, false)
}

func (service *validateOrderServiceImpl) validateOrder(dataFilter entity.ValidateOrderBody, includeStockValidation bool) (data entity.ValidateResponse, total int64, lastPage int, err error) {

	var vRespvald entity.ValidateResponse

	vRespvald.IsSuccessValidate = false
	vRespvald.Validate1 = "Sufficient Stock"
	vRespvald.Validate2 = "Within Limit"
	vRespvald.Validate3 = "Allowed"
	vRespvald.Validate4 = "Allowed"

	dataFilter.ActiveProductOnly = "true"

	dataFilter.Date = time.Now().Format("2006-01-02")

	for _, product := range dataFilter.ProStok {
		// fmt.Println("pro id", product.ProductId)
		dataFilter.ProID = append(dataFilter.ProID, product.ProductId)
	}

	vRespvald.Validate1Success = true
	if includeStockValidation {
		whStocks, err := service.ValidateOrderRepository.GetWarehouseStockByProducts(dataFilter.CustID, dataFilter.WhID, dataFilter.ProID)
		if err != nil {
			return data, 0, 0, err
		}

		//validate stock
		for _, row := range whStocks {
			for _, product := range dataFilter.ProStok {
				if product.ProductId != row.ProID {
					continue
				}

				QtyUnit := &conversion.QtyUnit{
					Qty1:      int(product.Qty1),
					Qty2:      int(product.Qty2),
					Qty3:      int(product.Qty3),
					ConvUnit2: int(row.ConvUnit2),
					ConvUnit3: int(row.ConvUnit3),
				}
				totalQty, _ := QtyUnit.ToTotalQuantity()

				QtyChangeUnit := &conversion.QtyUnit{
					Qty1:      int(product.QtyChange1),
					Qty2:      int(product.QtyChange2),
					Qty3:      int(product.QtyChange3),
					ConvUnit2: int(row.ConvUnit2),
					ConvUnit3: int(row.ConvUnit3),
				}
				totalQtyChange, _ := QtyChangeUnit.ToTotalQuantity()

				// Formula: wh_stock + (oncust - qty_changes)
				availableStock := row.Qty + (float64(totalQty) - float64(totalQtyChange))

				if float64(totalQty) > availableStock {
					vRespvald.Validate1Success = false
					vRespvald.Validate1 = "Insufficient Stock"
				}
			}
		}
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
		}
	}

	outlet, err := service.ValidateOrderRepository.FindDetailOutletID(dataFilter.OutletID)
	if err != nil {
		return data, total, lastPage, err
	}

	// fmt.Println(outlet)

	if outlet.CreditLimitType != nil && (*outlet.CreditLimitType == 1 || *outlet.CreditLimitType == 0) {
		vRespvald.Validate2Success = true
		vRespvald.Validate2 = "Allowed (Unlimited)"
	} else {
		var creditLimit float64
		if outlet.CreditLimit != nil {
			creditLimit = *outlet.CreditLimit
		}
		remainingCredit := (remainingLimit + dataFilter.Total) - creditLimit

		if remainingCredit <= 0 {
			vRespvald.Validate2Success = true
		} else {
			//fmt.Println(fmt.Sprintf("############ %v - %v", remainingLimit, dataFilter.Total))

			vRespvald.Validate2value = remainingCredit
			vRespvald.Validate2 = "Over Limit (" + formatCurrency(remainingCredit) + ")"
			//vRespvald.Validate2 = "Over Limit"
		}
	}

	if outlet.SalesInvLimitType != nil && (*outlet.SalesInvLimitType == 1 || *outlet.SalesInvLimitType == 0) {
		vRespvald.Validate3Success = true
		vRespvald.Validate3 = "Allowed (Unlimited)"
	} else {
		allowed := "Not Allowed"
		var salesInvLimit int
		if outlet.SalesInvLimit != nil {
			salesInvLimit = *outlet.SalesInvLimit
		}

		if dueDateCount >= salesInvLimit {
			vRespvald.Validate3Success = false
			allowed = "Not Allowed"

		} else {
			vRespvald.Validate3Success = true
			allowed = "Allowed"
		}
		vRespvald.Validate3 = allowed + " (" + fmt.Sprintf("%d", dueDateCount) + " of " + fmt.Sprintf("%d", salesInvLimit) + ")"
		vRespvald.Validate3Value = dueDateCount
	}

	if outlet.ObsType != nil && (*outlet.ObsType == 1 || *outlet.ObsType == 0) {
		vRespvald.Validate4Success = true
		vRespvald.Validate4 = "Allowed (Unlimited)"
	} else {
		allowed := "Not Allowed"
		var obs int
		if outlet.Obs != nil {
			obs = *outlet.Obs
		}
		if remainingObs <= obs {
			vRespvald.Validate4Success = true
			allowed = "Allowed"
		}

		vRespvald.Validate4 = allowed + " (" + fmt.Sprintf("%d", remainingObs) + " of " + fmt.Sprintf("%d", obs) + ")"
		vRespvald.Validate4Value = remainingObs
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
		if vResp.RemainingAmount > 0 {
			remainingObs++
			difference := paymentDate.Sub(dueDate)
			aging := int64(difference.Hours() / 24)
			vResp.Aging = &aging
			vResp.InvoiceStatusName = "Outstanding"
			outstandingArray = append(outstandingArray, vResp) // Tambahkan vResp ke array
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
