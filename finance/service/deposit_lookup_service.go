package service

import (
	"finance/entity"
	"finance/model"
	"finance/pkg/structs"
	"finance/repository"
	"math"
)

type DepositLookupService interface {
	LookupCollectionNo(dataFilter entity.GeneralQueryFilter) (data []entity.CollectionNoLookup, total int64, lastPage int, err error)
	LookupDepositNo(dataFilter entity.GeneralQueryFilter) (data []entity.DepositNoLookup, total int64, lastPage int, err error)
	LookupDepositStatus(dataFilter entity.GeneralQueryFilter) (data []entity.DepositStatusLookup, total int64, lastPage int, err error)

	ListInvoiceByCollection(dataFilter entity.GeneralQueryFilter) (data []entity.InvoiceByCollectionResponse, total int64, lastPage int, err error)

	ListBalancePaymentDepositByCustId(dataFilter entity.DepositLookupQueryFilter) (data map[string][]entity.DepositPaymentLookup, total int64, lastPage int, err error)
}

type DepositLookupServiceImpl struct {
	DepositLookupRepository repository.DepositLookupRepository
	Transaction             repository.Dbtransaction
}

func NewDepositLookupService(repository repository.DepositLookupRepository, transaction repository.Dbtransaction) *DepositLookupServiceImpl {
	return &DepositLookupServiceImpl{
		DepositLookupRepository: repository,
		Transaction:             transaction,
	}
}

func (service *DepositLookupServiceImpl) LookupCollectionNo(dataFilter entity.GeneralQueryFilter) (data []entity.CollectionNoLookup, total int64, lastPage int, err error) {
	DepositLookups, total, lastPage, err := service.DepositLookupRepository.FindAllCollectionNoByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range DepositLookups {
		var vResp entity.CollectionNoLookup
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *DepositLookupServiceImpl) LookupDepositNo(dataFilter entity.GeneralQueryFilter) (data []entity.DepositNoLookup, total int64, lastPage int, err error) {
	DepositLookups, total, lastPage, err := service.DepositLookupRepository.FindAllDepositNoByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range DepositLookups {
		var vResp entity.DepositNoLookup
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *DepositLookupServiceImpl) LookupDepositStatus(dataFilter entity.GeneralQueryFilter) (data []entity.DepositStatusLookup, total int64, lastPage int, err error) {
	DepositLookups, total, lastPage, err := service.DepositLookupRepository.FindAllDepositStatusByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range DepositLookups {
		var vResp entity.DepositStatusLookup
		structs.Automapper(row, &vResp)
		DepositStatusName := entity.ConvStatusBankTransfer(entity.StatusDeposit, row.DepositStatus)
		vResp.DepositStatusName = DepositStatusName
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *DepositLookupServiceImpl) ListInvoiceByCollection(dataFilter entity.GeneralQueryFilter) (data []entity.InvoiceByCollectionResponse, total int64, lastPage int, err error) {
	DepositLookups, total, lastPage, err := service.DepositLookupRepository.FindInvoiceByCollectionByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range DepositLookups {
		var vResp entity.InvoiceByCollectionResponse
		vResp.TotalInvoicePaymnet = new(float64) // Menggunakan pointer untuk float64
		*vResp.TotalInvoicePaymnet = 0.0         // Mengatur nilai
		structs.Automapper(row, &vResp)
		// if row.TransferDate != nil {
		// 	TransferDate := row.TransferDate.Format("2006-01-02")
		// 	vResp.TransferDate = &TransferDate
		// }

		// ownerName := entity.ConvStatusBankTransfer(entity.OwnerBank, row.OwnerID)
		// vResp.OwnerName = ownerName

		// statusText := entity.ConvStatusBankTransfer(entity.StatusBank, row.StatusBank)
		// vResp.StatusBankText = &statusText

		// vResp.UsedAmount = float64(0)
		// vResp.RemainingAmount = row.Amount - vResp.UsedAmount

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *DepositLookupServiceImpl) ListBalancePaymentDepositByCustId(dataFilter entity.DepositLookupQueryFilter) (mData map[string][]entity.DepositPaymentLookup, total int64, lastPage int, err error) {
	modes := []string{"check", "transfer", "cndn", "return"}

	if len(dataFilter.Mode) == 0 {
		dataFilter.Mode = modes
	}

	mData = make(map[string][]entity.DepositPaymentLookup, len(dataFilter.Mode))

	for _, mode := range dataFilter.Mode {
		var DepositLookups []model.DepositPaymentLookup
		var subTotal int64
		switch mode {
		case "check":
			DepositLookups, subTotal, _, err = service.DepositLookupRepository.FindChequeGiroBalance(dataFilter)
		case "transfer":
			DepositLookups, subTotal, _, err = service.DepositLookupRepository.FindBankTransferBalance(dataFilter)
		case "cndn":
			DepositLookups, subTotal, _, err = service.DepositLookupRepository.FindCNDNBalance(dataFilter)
		case "return":
			DepositLookups, subTotal, _, err = service.DepositLookupRepository.FindReturnBalance(dataFilter)
		default:
			continue
		}

		if err != nil {
			return
		}

		for _, row := range DepositLookups {
			var vResp entity.DepositPaymentLookup
			structs.Automapper(row, &vResp)

			mData[mode] = append(mData[mode], vResp)
			// data = append(data, vResp)
		}
		total += subTotal
	}

	lastPage = int(math.Ceil(float64(total) / float64(dataFilter.Limit)))

	return
}
