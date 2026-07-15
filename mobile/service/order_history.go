package service

import (
	"mobile/entity"
	"mobile/pkg/config/env"
	"mobile/pkg/structs"
	"mobile/repository"
)

type OrderHistoryService interface {
	Detail(RoNo string, custID string) (response entity.OrderHistoryDetailResponse, err error)
	List(dataFilter entity.OrderQueryFilter) (data []entity.OrderHistoryListResponse, total int64, lastPage int, err error)
}

func NewOrderHistoryService(
	config env.ConfigEnv,
	orderHistoryRepository repository.OrderHistoryRepository,
	transaction repository.Dbtransaction) *orderHistoryServiceImpl {
	return &orderHistoryServiceImpl{
		Config:                 config,
		orderHistoryRepository: orderHistoryRepository,
		Transaction:            transaction,
	}
}

type orderHistoryServiceImpl struct {
	Config                 env.ConfigEnv
	orderHistoryRepository repository.OrderHistoryRepository
	Transaction            repository.Dbtransaction
}

func (service *orderHistoryServiceImpl) Detail(RoNo string, custID string) (response entity.OrderHistoryDetailResponse, err error) {
	ro, err := service.orderHistoryRepository.FindByNo(RoNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ro, &response)
	if err != nil {
		return response, err
	}

	details, err := service.orderHistoryRepository.FindDetail(RoNo, custID)
	if err != nil {
		return response, err
	}
	for _, detail := range details {
		var detailData entity.OrderHistoryDetailResponse
		err = structs.Automapper(detail, &detailData)

		detailsProduct, err := service.orderHistoryRepository.FindDetailProductOrderHistory(detail.RoNo)
		if err != nil {
			return response, err
		}
		for _, detailRowProduct := range detailsProduct {
			var productDetail entity.OrderHistoryProductDetails
			err = structs.Automapper(detailRowProduct, &productDetail)
			if err != nil {
				return response, err
			}
			detailData.ProducDetails = append(detailData.ProducDetails, productDetail)
		}
		response.ProducDetails = append(response.ProducDetails, detailData.ProducDetails...)
	}
	if ro.RoDate != nil {
		roDate := ro.RoDate.Format("2006-01-02")
		response.RoDate = &roDate
	}
	if ro.ValDate != nil {
		ValDate := ro.ValDate.Format("2006-01-02")
		response.ValDate = &ValDate
	}
	if ro.DeliveryDate != nil {
		DelivDate := ro.DeliveryDate.Format("2006-01-02")
		response.DeliveryDate = &DelivDate
	}
	if ro.InvoiceDate != nil {
		invoiceDate := ro.InvoiceDate.Format("2006-01-02")
		response.InvoiceDate = &invoiceDate
	}

	statusName := response.GenerateDataStatusNameHistory()
	response.DataStatusName = statusName

	payTypeName := response.GeneratePayTypeNameHistory()
	response.PayTypeName = payTypeName

	return response, nil
}

func (service *orderHistoryServiceImpl) List(dataFilter entity.OrderQueryFilter) (data []entity.OrderHistoryListResponse, total int64, lastPage int, err error) {
	realOrders, total, lastPage, err := service.orderHistoryRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range realOrders {
		var vResp entity.OrderHistoryListResponse
		structs.Automapper(row, &vResp)

		details, err := service.orderHistoryRepository.FindDetailProductOrderHistory(row.RoNo)
		if err != nil {
			return nil, 0, 0, err
		}

		detailsPayment, err := service.orderHistoryRepository.FindDetailProductOrderHistoryPayment(row.RoNo)
		if err != nil {
			return nil, 0, 0, err
		}

		for _, detail := range detailsPayment {
			if detail.InvoiceNo != nil {
				isPayment := true
				vResp.IsPayment = &isPayment
			} else {
				isPayment := false
				vResp.IsPayment = &isPayment
			}
		}

		// if len(detailsPayment) > 0 {
		// 	isPayment := true
		// 	vResp.IsPayment = &isPayment
		// } else {
		// 	isPayment := false
		// 	vResp.IsPayment = &isPayment
		// }

		for _, detail := range details {
			var detailData entity.ProductImgRead
			err = structs.Automapper(detail, &detailData)
			if err != nil {
				return nil, 0, 0, err
			}
			vResp.ProductImg = append(vResp.ProductImg, detailData)
		}
		if row.RoDate != nil {
			roDate := row.RoDate.Format("2006-01-02")
			vResp.RoDate = &roDate
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
