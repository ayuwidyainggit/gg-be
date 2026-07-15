package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"finance/pkg/structs"
	"finance/repository"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type BankTransferService interface {
	Store(request entity.CreateBankTransferBody) (err error)
	Detail(BankTransferNo int, custID string, parentCustId string) (response entity.BankTransferResponse, err error)
	List(dataFilter entity.BankTransferQueryFilter) (data []entity.BankTransferResponse, total int64, lastPage int, err error)
	Delete(custId string, BankTransferNo int, userId int64) (err error)
	Update(BankTransferNo int, request entity.UpdateBankTransferBody) (err error)

	LookupBankTransfer(dataFilter entity.BankTransferQueryFilter) (data []entity.BankLookupBankTransfer, total int64, lastPage int, err error)
	LookupBankAccountBankTransfer(dataFilter entity.BankTransferQueryFilter) (data []entity.BankAccountLookupBankTransfer, total int64, lastPage int, err error)
}

type BankTransferServiceImpl struct {
	BankTransferRepository repository.BankTransferRepository
	Transaction            repository.Dbtransaction
}

func NewBankTransferService(repository repository.BankTransferRepository, transaction repository.Dbtransaction) *BankTransferServiceImpl {
	return &BankTransferServiceImpl{
		BankTransferRepository: repository,
		Transaction:            transaction,
	}
}

func (service *BankTransferServiceImpl) Store(request entity.CreateBankTransferBody) (err error) {
	c := context.Background()

	// Parse Transfer Date
	var trfDate time.Time
	if request.Document.TransferDate != "" {
		trfDate, err = time.Parse("2006-01-02", request.Document.TransferDate)
		if err != nil {
			return err
		}
	} else {
		trfDate = time.Now()
	}

	// Auto Generate DocNoBank
	datePart := trfDate.Format("060102") // YYMMDD
	prefix := "TF" + datePart
	lastDocNo, err := service.BankTransferRepository.GetLastDocNoBank(c, prefix)
	if err != nil && !strings.Contains(err.Error(), "record not found") && !strings.Contains(err.Error(), "no rows") {
		return err
	}

	seq := 1
	if lastDocNo != "" {
		// TFYYMMDDNNNN (12 chars)
		if len(lastDocNo) >= 12 {
			lastSeqStr := lastDocNo[8:]
			lastSeq, _ := strconv.Atoi(lastSeqStr)
			seq = lastSeq + 1
		}
	}
	newDocNo := fmt.Sprintf("%s%04d", prefix, seq)

	// Prepare Model
	outletID := request.Outlet.OutletID
	outletBankID := request.Outlet.OutletBankID
	accountNo := request.Bank.AccountNo
	accountName := request.Bank.AccountName

	bankTransferModel := model.BankTransfer{
		CustID:       request.CustID,
		DocNoBank:    newDocNo,
		OwnerID:      request.OwnerID,
		SalesmanID:   request.SalesmanID,
		SupplierID:   request.Supplier.SupID,
		OutletID:     &outletID,
		BankID:       request.Bank.BankID,
		AccountNo:    &accountNo,
		AccountName:  accountName,
		OutletBankID: &outletBankID,
		TransferDate: &trfDate,
		Amount:       request.Amount,
		StatusBank:   2, // Pending
	}

	if request.CreatedBy != nil {
		bankTransferModel.CreatedBy = int(*request.CreatedBy)
	}

	if bankTransferModel.OwnerID == 1 {
		bankTransferModel.SupplierID = nil
	} else {
		bankTransferModel.OutletID = nil
		bankTransferModel.SalesmanID = nil
		bankTransferModel.OutletBankID = nil
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.BankTransferRepository.Store(txCtx, &bankTransferModel)
		if err != nil {
			return err
		}

		// Store Files
		for _, file := range request.ProofOfPayment {
			fileModel := model.BankTransferFile{
				CustID:         request.CustID,
				BankTransferNo: newDocNo,
				FileName:       file.FileName,
				FileURL:        file.FileURL,
				FileKey:        file.FileKey,
				FileSize:       file.FileSize,
				MediaCategory:  file.MediaCategory,
				CreatedAt:      time.Now(),
			}
			if err := service.BankTransferRepository.StoreFile(txCtx, &fileModel); err != nil {
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

func (service *BankTransferServiceImpl) Detail(BankTransferNo int, custID string, parentCustId string) (response entity.BankTransferResponse, err error) {
	ctx := context.Background()
	BankTransfer, err := service.BankTransferRepository.FindByNo(BankTransferNo, custID, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(BankTransfer, &response)
	if err != nil {
		return response, err
	}

	response.AccountName = BankTransfer.AccountName

	TrfDate := BankTransfer.TransferDate.Format("2006-01-02")
	response.TransferDate = &TrfDate

	ownerName := entity.ConvStatusBankTransfer(entity.OwnerBank, response.OwnerID)
	response.OwnerName = ownerName

	statusText := entity.ConvStatusBankTransfer(entity.StatusBank, response.StatusBank)
	response.StatusBankText = &statusText

	// used_amount = sum(deposit_data.used_amount); remaining_amount = amount - used_amount
	response.DepositData = make([]entity.BankTransferDepositDataItem, 0)
	response.UsedAmount = 0
	depositRows, _ := service.BankTransferRepository.FindDepositDataByDocumentNo(ctx, BankTransfer.DocNoBank, custID)
	for _, row := range depositRows {
		depositDateStr := ""
		if row.DepositDate != nil {
			depositDateStr = row.DepositDate.Format("2006-01-02")
		}
		response.DepositData = append(response.DepositData, entity.BankTransferDepositDataItem{
			DepositNo:   row.DepositNo,
			DepositDate: depositDateStr,
			InvoiceNo:   row.InvoiceNo,
			UsedAmount:  row.UsedAmount,
		})
		response.UsedAmount += row.UsedAmount
	}
	response.RemainingAmount = response.Amount - response.UsedAmount

	response.ProofOfPayment = make([]entity.BankTransferFileBody, 0)
	files, err := service.BankTransferRepository.FindFilesByDocNo(ctx, BankTransfer.DocNoBank, custID)
	if err == nil && len(files) > 0 {
		for _, f := range files {
			response.ProofOfPayment = append(response.ProofOfPayment, entity.BankTransferFileBody{
				FileName:      f.FileName,
				FileURL:       f.FileURL,
				FileKey:       f.FileKey,
				FileSize:      f.FileSize,
				MediaCategory: f.MediaCategory,
			})
		}
	}

	return response, nil
}

func (service *BankTransferServiceImpl) List(dataFilter entity.BankTransferQueryFilter) (data []entity.BankTransferResponse, total int64, lastPage int, err error) {
	ctx := context.Background()
	BankTransfers, total, lastPage, err := service.BankTransferRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range BankTransfers {
		var vResp entity.BankTransferResponse
		structs.Automapper(row, &vResp)
		if row.TransferDate != nil {
			TransferDate := row.TransferDate.Format("2006-01-02")
			vResp.TransferDate = &TransferDate
		}

		ownerName := entity.ConvStatusBankTransfer(entity.OwnerBank, row.OwnerID)
		vResp.OwnerName = ownerName

		statusText := entity.ConvStatusBankTransfer(entity.StatusBank, row.StatusBank)
		vResp.StatusBankText = &statusText

		// used_amount = sum(deposit_data.used_amount); remaining_amount = amount - used_amount
		vResp.DepositData = make([]entity.BankTransferDepositDataItem, 0)
		vResp.UsedAmount = 0
		depositRows, _ := service.BankTransferRepository.FindDepositDataByDocumentNo(ctx, row.DocNoBank, dataFilter.CustId)
		for _, depositRow := range depositRows {
			depositDateStr := ""
			if depositRow.DepositDate != nil {
				depositDateStr = depositRow.DepositDate.Format("2006-01-02")
			}
			vResp.DepositData = append(vResp.DepositData, entity.BankTransferDepositDataItem{
				DepositNo:   depositRow.DepositNo,
				DepositDate: depositDateStr,
				InvoiceNo:   depositRow.InvoiceNo,
				UsedAmount:  depositRow.UsedAmount,
			})
			vResp.UsedAmount += depositRow.UsedAmount
		}
		vResp.RemainingAmount = row.Amount - vResp.UsedAmount

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *BankTransferServiceImpl) Delete(custId string, BankTransferNo int, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.BankTransferRepository.Delete(txCtx, custId, BankTransferNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *BankTransferServiceImpl) Update(BankTransferNo int, request entity.UpdateBankTransferBody) (err error) {
	c := context.Background()

	if request.TransferDate != nil {
		TrfDate, err := str.DateStrToRfc3339String(*request.TransferDate)
		if err != nil {
			return err
		}
		request.TransferDate = &TrfDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.BankTransfer
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}

	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.BankTransferRepository.Update(txCtx, BankTransferNo, request.CustID, Model)
		if err != nil {
			return err
		}

		// Handle ProofOfPayment Update
		if request.ProofOfPayment.Mode == "replace" {
			// Get existing data to find DocNoBank
			existingData, err := service.BankTransferRepository.FindOne(txCtx, BankTransferNo, request.CustID)
			if err != nil {
				return err
			}

			// Delete old files
			err = service.BankTransferRepository.DeleteFilesByDocNo(txCtx, existingData.DocNoBank, request.CustID)
			if err != nil {
				return err
			}

			// Insert new files
			for _, file := range request.ProofOfPayment.Files {
				fileModel := model.BankTransferFile{
					CustID:         request.CustID,
					BankTransferNo: existingData.DocNoBank,
					FileName:       file.FileName,
					FileURL:        file.FileURL,
					FileKey:        file.FileKey,
					FileSize:       file.FileSize,
					MediaCategory:  file.MediaCategory,
					CreatedAt:      time.Now(),
				}
				if err := service.BankTransferRepository.StoreFile(txCtx, &fileModel); err != nil {
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

func (service *BankTransferServiceImpl) LookupBankTransfer(dataFilter entity.BankTransferQueryFilter) (data []entity.BankLookupBankTransfer, total int64, lastPage int, err error) {
	BankTransfers, total, lastPage, err := service.BankTransferRepository.FindAllBankByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range BankTransfers {
		var vResp entity.BankLookupBankTransfer
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *BankTransferServiceImpl) LookupBankAccountBankTransfer(dataFilter entity.BankTransferQueryFilter) (data []entity.BankAccountLookupBankTransfer, total int64, lastPage int, err error) {
	BankTransfers, total, lastPage, err := service.BankTransferRepository.FindAllBankAccountByCustId(dataFilter, dataFilter.BankID)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range BankTransfers {
		var vResp entity.BankAccountLookupBankTransfer
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
