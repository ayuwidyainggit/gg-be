package service

import (
	"context"
	"errors"
	"sales/entity"
	"sales/model"
	"sales/pkg/conversion"
	"sales/pkg/str"
	"sales/pkg/structs"
	"sales/repository"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type InvoiceService interface {
	Detail(RoNo string, custID string) (response entity.InvoiceResponse, err error)
	List(dataFilter entity.InvoiceQueryFilter) (data []entity.InvoiceListResponse, total int64, lastPage int, err error)
	Details(dataFilter entity.InvoiceQueryFilter) (data []entity.InvoiceListResponse, err error)
	BulkUpdate(custId string, request entity.BulkUpdateInvoiceBody) (err error)
	Print(custId string, invoiceNo string, userId int64) (err error)
}

func NewInvoiceService(invoiceRepository repository.InvoiceRepository, stockRepository repository.StockRepository, transaction repository.Dbtransaction) *invoiceServiceImpl {
	return &invoiceServiceImpl{
		InvoiceRepository: invoiceRepository,
		StockRepository:   stockRepository,
		Transaction:       transaction,
	}
}

type invoiceServiceImpl struct {
	InvoiceRepository repository.InvoiceRepository
	StockRepository   repository.StockRepository
	Transaction       repository.Dbtransaction
}

func (service *invoiceServiceImpl) Detail(RoNo string, custID string) (response entity.InvoiceResponse, err error) {
	invoice, err := service.InvoiceRepository.FindByNo(RoNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(invoice, &response)
	if err != nil {
		return response, err
	}

	details, err := service.InvoiceRepository.FindDetail(RoNo, custID)
	if err != nil {
		return response, err
	}

	finalTotals := calculateInvoiceFinalTotals(details)
	response.SubTotal = float64Ptr(finalTotals.Gross)
	response.PromoValue = float64Ptr(finalTotals.PromoTotal())
	response.DiscValue = float64Ptr(finalTotals.Discount)
	response.VatValue = float64Ptr(finalTotals.VAT)
	response.Total = float64Ptr(finalTotals.Net)

	for _, detail := range details {
		detailData, _, mapErr := mapInvoiceFinalDetailResponse(detail)
		if mapErr != nil {
			return response, mapErr
		}
		if detailData.ItemType == 1 {
			response.Details.Normal = append(response.Details.Normal, detailData)
		} else {
			response.Details.Promo = append(response.Details.Promo, detailData)
		}

	}

	if invoice.DeliveryDate != nil {
		delivDate := invoice.DeliveryDate.Format("2006-01-02")
		response.DeliveryDate = &delivDate
	}

	statusName := response.GenerateDataStatusName()
	response.DataStatusName = statusName

	payTypeName := response.GeneratePayTypeName()
	response.PayTypeName = payTypeName

	return response, nil
}

func (service *invoiceServiceImpl) List(dataFilter entity.InvoiceQueryFilter) (data []entity.InvoiceListResponse, total int64, lastPage int, err error) {
	invoice, total, lastPage, err := service.InvoiceRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range invoice {
		var vResp entity.InvoiceListResponse
		structs.Automapper(row, &vResp)

		if row.DeliveryDate != nil {
			delivDate := row.DeliveryDate.Format("2006-01-02")
			vResp.DeliveryDate = &delivDate
		}
		if row.DueDate != nil {
			dueDate := row.DueDate.Format("2006-01-02")
			vResp.DueDate = &dueDate
		}

		if row.InvoiceDate != nil {
			invoiceDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = invoiceDate
		}

		statusName := vResp.GenerateDataStatusName()
		vResp.DataStatusName = statusName

		payTypeName := vResp.GeneratePayTypeName()
		vResp.PayTypeName = payTypeName

		details, err := service.InvoiceRepository.FindDetail(row.OrderNo, row.CustID)
		if err != nil {
			return data, total, lastPage, err
		}

		finalTotals := calculateInvoiceFinalTotals(details)
		vResp.SubTotal = float64Ptr(finalTotals.Gross)
		vResp.PromoValue = float64Ptr(finalTotals.PromoTotal())
		vResp.DiscValue = float64Ptr(finalTotals.Discount)
		vResp.VatValue = float64Ptr(finalTotals.VAT)
		vResp.Total = float64Ptr(finalTotals.Net)

		for _, detail := range details {
			detailData, _, mapErr := mapInvoiceFinalDetailResponse(detail)
			if mapErr != nil {
				return data, total, lastPage, mapErr
			}
			vResp.TotalVolume += (detail.Volume1 * detail.Qty1Final) + (detail.Volume2 * detail.Qty2Final) + (detail.Volume3 * detail.Qty3Final)
			vResp.TotalWeight += (detail.Weight1 * detail.Qty1Final) + (detail.Weight2 * detail.Qty2Final) + (detail.Weight3 * detail.Qty3Final)
			vResp.Details = append(vResp.Details, detailData)
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *invoiceServiceImpl) Details(dataFilter entity.InvoiceQueryFilter) (data []entity.InvoiceListResponse, err error) {
	invoice, err := service.InvoiceRepository.FindAllByInvoiceNombersAndCustId(dataFilter)
	if err != nil {
		return data, err
	}

	for _, row := range invoice {
		var vResp entity.InvoiceListResponse
		structs.Automapper(row, &vResp)

		if row.DeliveryDate != nil {
			delivDate := row.DeliveryDate.Format("2006-01-02")
			vResp.DeliveryDate = &delivDate
		}
		if row.DueDate != nil {
			dueDate := row.DueDate.Format("2006-01-02")
			vResp.DueDate = &dueDate
		}

		statusName := vResp.GenerateDataStatusName()
		vResp.DataStatusName = statusName

		payTypeName := vResp.GeneratePayTypeName()
		vResp.PayTypeName = payTypeName

		details, err := service.InvoiceRepository.FindDetail(row.OrderNo, row.CustID)
		if err != nil {
			return data, err
		}

		finalTotals := calculateInvoiceFinalTotals(details)
		vResp.SubTotal = float64Ptr(finalTotals.Gross)
		vResp.PromoValue = float64Ptr(finalTotals.PromoTotal())
		vResp.DiscValue = float64Ptr(finalTotals.Discount)
		vResp.VatValue = float64Ptr(finalTotals.VAT)
		vResp.Total = float64Ptr(finalTotals.Net)

		for _, detail := range details {
			detailData, _, mapErr := mapInvoiceFinalDetailResponse(detail)
			if mapErr != nil {
				return data, mapErr
			}

			vResp.TotalVolume += (detail.Volume1 * detail.Qty1Final) + (detail.Volume2 * detail.Qty2Final) + (detail.Volume3 * detail.Qty3Final)
			vResp.TotalWeight += (detail.Weight1 * detail.Qty1Final) + (detail.Weight2 * detail.Qty2Final) + (detail.Weight3 * detail.Qty3Final)
			vResp.Details = append(vResp.Details, detailData)
		}

		data = append(data, vResp)
	}

	return data, err
}

func (service *invoiceServiceImpl) BulkUpdate(custId string, request entity.BulkUpdateInvoiceBody) (err error) {
	c := context.Background()
	dateNow := time.Now()
	invoiceDate, err := str.DateStrToRfc3339String(dateNow.Format("2006-01-02"))
	if err != nil {
		return err
	}

	for index := range request.Orders {
		dataStatus := int64(6)
		request.Orders[index].InvoiceDate = &invoiceDate
		request.Orders[index].DataStatus = &dataStatus

		if request.Orders[index].RoDate != nil {
			roDate, err := str.DateStrToRfc3339String(*request.Orders[index].RoDate)
			if err != nil {
				return err
			}
			request.Orders[index].RoDate = &roDate
		}

		if request.Orders[index].ValDate != nil {
			valDate, err := str.DateStrToRfc3339String(*request.Orders[index].ValDate)
			if err != nil {
				return err
			}
			request.Orders[index].ValDate = &valDate
		}

		if request.Orders[index].DueDate != nil {
			dueDate, err := str.DateStrToRfc3339String(*request.Orders[index].DueDate)
			if err != nil {
				return err
			}
			request.Orders[index].DueDate = &dueDate
		}

		if request.Orders[index].DeliveryDate != nil {
			deliveryDate, err := str.DateStrToRfc3339String(*request.Orders[index].DeliveryDate)
			if err != nil {
				return err
			}
			request.Orders[index].DeliveryDate = &deliveryDate
		}

		var invoiceModel model.Invoice
		err = structs.Automapper(request.Orders[index], &invoiceModel)
		if err != nil {
			return err
		}
		invoiceModel.CustID = ""

		const maxRetryAttempt = 3
		for retryAttempt := 1; retryAttempt <= maxRetryAttempt; retryAttempt++ {
			err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
				generatedInvoiceNo, err := service.InvoiceRepository.GenerateInvoiceNo(txCtx, custId, dateNow)
				if err != nil {
					return err
				}

				request.Orders[index].InvoiceNo = &generatedInvoiceNo
				invoiceModel.InvoiceNo = &generatedInvoiceNo

				invoiceExist, err := service.InvoiceRepository.FindByNo(request.Orders[index].RoNo, custId)
				if err != nil {
					return err
				}
				today := time.Now()

				if invoiceExist.PaymentType == entity.PAY_TYPE_CREDIT {
					nextDate := today.AddDate(0, 0, invoiceExist.TOP)
					nextDate = time.Date(nextDate.Year(), nextDate.Month(), nextDate.Day(), 0, 0, 0, 0, nextDate.Location())
					invoiceModel.DueDate = &nextDate
				} else {
					invoiceModel.DueDate = &today
				}

				details, err := service.InvoiceRepository.FindDetail(request.Orders[index].RoNo, custId)
				if err != nil {
					return err
				}

				finalTotals := calculateInvoiceFinalTotals(details)
				invoiceModel.SubTotal = nil
				invoiceModel.PromoValue = nil
				invoiceModel.DiscValue = nil
				invoiceModel.VatValue = nil
				invoiceModel.Total = nil
				invoiceModel.SubTotalFinal = float64Ptr(finalTotals.Gross)
				invoiceModel.PromoValueFinal = float64Ptr(finalTotals.PromoTotal())
				invoiceModel.DiscValueFinal = float64Ptr(finalTotals.Discount)
				invoiceModel.VatValueFinal = float64Ptr(finalTotals.VAT)
				invoiceModel.TotalFinal = float64Ptr(finalTotals.Net)

				err = service.InvoiceRepository.Update(txCtx, request.Orders[index].RoNo, custId, invoiceModel)
				if err != nil {
					return err
				}

				if invoiceExist.OutletID != nil {
					err = service.InvoiceRepository.UpdateOutletStatusFromPreDormantIfSet(
						txCtx, custId, *invoiceExist.OutletID, request.Orders[index].UpdatedBy,
					)
					if err != nil {
						return err
					}
				}

				invoiceSalesStockUpdateEntities := make([]*entity.InvoiceSalesStockUpdate, 0, len(details))
				for _, detail := range details {
					qtyUnit := &conversion.QtyUnit{
						Qty1:      int(detail.Qty1Final),
						Qty2:      int(detail.Qty2Final),
						Qty3:      int(detail.Qty3Final),
						ConvUnit2: int(detail.ConvUnit2),
						ConvUnit3: int(detail.ConvUnit3),
					}

					totalQtyFinal, err := qtyUnit.ToTotalQuantity()
					if err != nil {
						return err
					}

					salesOrderStockUpdateEntity := entity.InvoiceSalesStockUpdate{
						CustID:         detail.CustId,
						WhID:           *invoiceExist.WhId,
						ProID:          int64(detail.ProId),
						StockDate:      *invoiceExist.RoDate,
						TrCode:         request.Orders[index].RoNo[0:2],
						TrNo:           request.Orders[index].RoNo,
						QtyOrderBefore: float64(totalQtyFinal),
						UnitPrice:      detail.SellPrice1,
						RefDetId:       int64(*detail.OrderDetID),
					}
					invoiceSalesStockUpdateEntities = append(invoiceSalesStockUpdateEntities, &salesOrderStockUpdateEntity)
				}

				err = service.StockRepository.InvoiceSalesStockUpdates(txCtx, invoiceSalesStockUpdateEntities)
				if err != nil {
					return err
				}

				return nil
			})

			if err == nil {
				break
			}
			if !isUniqueViolation(err) || retryAttempt == maxRetryAttempt {
				return err
			}
		}
	}

	return nil
}

func mapInvoiceFinalDetailResponse(detail model.InvoiceDetRead) (entity.InvoiceDetResponse, invoiceFinalLineAmount, error) {
	var detailData entity.InvoiceDetResponse
	err := structs.Automapper(detail, &detailData)
	if err != nil {
		return detailData, invoiceFinalLineAmount{}, err
	}

	lineAmount := calculateInvoiceFinalLineAmount(detail)
	detailData.Qty1 = detail.Qty1Final
	detailData.Qty2 = detail.Qty2Final
	detailData.Qty3 = detail.Qty3Final
	detailData.SellPrice1 = detail.SellPriceFinal1
	detailData.SellPrice2 = detail.SellPriceFinal2
	detailData.SellPrice3 = detail.SellPriceFinal3
	detailData.Amount = lineAmount.Gross
	detailData.DiscValue = lineAmount.Discount
	detailData.VatValue = lineAmount.VAT
	detailData.NetValue = lineAmount.Net
	if detail.ExpDate != nil {
		detailData.ExpDate = detail.ExpDate.Format("2006-01-02")
	}

	return detailData, lineAmount, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return strings.Contains(err.Error(), "SQLSTATE 23505")
}

func (service *invoiceServiceImpl) Print(custId string, invoiceNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.InvoiceRepository.Print(txCtx, custId, invoiceNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
