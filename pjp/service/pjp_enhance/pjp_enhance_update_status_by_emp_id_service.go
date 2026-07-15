package pjpenhance

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
	"strconv"
)

func (service *pjpEnhanceService) UpdateStatusByEmpId(ctx context.Context, empId int, request request.UpdateStatusPjpEnhanceRequest, currentCustomerId string) {
	// Validasi struct menggunakan validator
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	// Mulai transaksi
	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data PJP berdasarkan empId dan customerId
	pjp := service.pjpRepository.GetPjpByEmpId(ctx, tx, empId, currentCustomerId)

	// Update field status (wajib)
	if request.Status != nil {
		pjp.Status = strconv.FormatBool(*request.Status)
	}

	// Update field opsional jika tidak nil
	if request.SalesmanName != nil {
		pjp.SalesmanName = *request.SalesmanName
	}
	if request.SalesmanCOde != nil {
		pjp.SalesmanCode = *request.SalesmanCOde
	}
	if request.WarehouseID != nil {
		pjp.WarehouseID = *request.WarehouseID
	}
	if request.WarehouseName != nil {
		pjp.WarehouseName = *request.WarehouseName
	}
	if request.OperationType != nil {
		pjp.OperationType = *request.OperationType
	}
	if request.TeamSalesMan != nil {
		pjp.TeamSalesMan = *request.TeamSalesMan
	}

	// Simpan perubahan
	service.pjpRepository.Update(ctx, tx, pjp)
}
