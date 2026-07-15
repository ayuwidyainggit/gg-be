package service

import (
	"context"
	"errors"
	"fmt"
	"sales/entity"
	"sales/model"
	"sales/repository"
)

type HierarchyApprovalService interface {
	Store(request entity.CreateHierarchyApprovalBody) (err error)
	Companies(dataFilter entity.CompaniesQueryFilter) (datas []entity.CompaniesListResponse, total int64, lastPage int, err error)
	List(dataFilter entity.HierarcyApprovalQueryFilter) (data []entity.HierarchyApprovalListResp, total int64, lastPage int, err error)
	Detail(hierarchyApprovalID int64, custID string, parentCustID string) (response entity.HierarchyApprovalReadResp, err error)
	Delete(custID string, parentCustId string, hierarchyApprovalID int64, userID int64) (err error)
	Update(hierarchyApprovalID int64, request entity.UpdateHierarchyApprovalBody) (err error)
	Employees(dataFilter entity.EmployeeHierarchyQueryFilter) (datas []entity.EmployeeResp, total int64, lastPage int, err error)
	RequestApproval(request entity.RequestApprovalBody) (err error)
	RequestApprovalDetail(custID string, requestApprovalID int64) (detail entity.ApprovalDetail, err error)
}

func NewHierarchyApprovalService(HierarchyApprovalRepository repository.HierarchyApprovalRepository, orderRepository repository.OrderRepository, orderApprovalRequest repository.OrderApprovalRequestRepository, transaction repository.Dbtransaction) *HierarchyApprovalServiceImpl {
	return &HierarchyApprovalServiceImpl{
		HierarchyApprovalRepository: HierarchyApprovalRepository,
		OrderRepository:             orderRepository,
		OrderApprovalRequest:        orderApprovalRequest,
		Transaction:                 transaction,
	}
}

type HierarchyApprovalServiceImpl struct {
	HierarchyApprovalRepository repository.HierarchyApprovalRepository
	Transaction                 repository.Dbtransaction
	OrderRepository             repository.OrderRepository
	OrderApprovalRequest        repository.OrderApprovalRequestRepository
}

func (service *HierarchyApprovalServiceImpl) Employees(dataFilter entity.EmployeeHierarchyQueryFilter) (datas []entity.EmployeeResp, total int64, lastPage int, err error) {
	employeesModel, total, lastPage, err := service.HierarchyApprovalRepository.FindAllEmployeeByCustIdLookupMode(dataFilter)
	if err != nil {
		return
	}

	for _, empModel := range employeesModel {
		datas = append(datas, entity.EmployeeResp{
			CustID:       empModel.CustId,
			EmployeeId:   empModel.EmployeeId,
			EmployeeCode: empModel.EmployeeCode,
			EmployeeName: empModel.EmployeeName,
		})
	}
	return

}

func (service *HierarchyApprovalServiceImpl) Companies(dataFilter entity.CompaniesQueryFilter) (datas []entity.CompaniesListResponse, total int64, lastPage int, err error) {
	cust, err := service.HierarchyApprovalRepository.GetUser(dataFilter.CustId)
	if err != nil {
		return
	}
	var companies []model.SmcMCustomer

	if cust.CustId != cust.ParentCustID { // jika child (distributor) hanya ambil data distributor itu saja dan parentnya
		custID := dataFilter.CustId
		dataFilter.CustId = dataFilter.ParentCustId
		companies, total, lastPage, err = service.HierarchyApprovalRepository.FindCompanies(dataFilter)
		if err != nil {
			return
		}

		for _, company := range companies {
			var headOffice bool
			if company.CustId == company.ParentCustID {
				headOffice = true

				datas = append(datas, entity.CompaniesListResponse{
					HeadOffice:  headOffice,
					CompanyID:   company.CustId,
					CompanyCode: company.CompanyCode,
					CompanyName: company.CompanyName,
				})

			}

			if company.CustId == custID {
				datas = append(datas, entity.CompaniesListResponse{
					HeadOffice:  false,
					CompanyID:   company.CustId,
					CompanyCode: company.CompanyCode,
					CompanyName: company.CompanyName,
				})
			}
		}

		return
	}

	companies, total, lastPage, err = service.HierarchyApprovalRepository.FindCompanies(dataFilter)
	if err != nil {
		return
	}

	for _, company := range companies {
		var headOffice bool
		if company.CustId == company.ParentCustID {
			headOffice = true
		}
		datas = append(datas, entity.CompaniesListResponse{
			HeadOffice:  headOffice,
			CompanyID:   company.CustId,
			CompanyCode: company.CompanyCode,
			CompanyName: company.CompanyName,
		})
	}
	return

}

func (service *HierarchyApprovalServiceImpl) Store(request entity.CreateHierarchyApprovalBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		for _, setupFor := range request.SetupFor {

			setupHierarchyApproval, err := service.HierarchyApprovalRepository.FindBySetupFor(setupFor)
			if err == nil {
				return errors.New(fmt.Sprintf("company '%v' already setup with id %v", setupHierarchyApproval.CompanyName, setupHierarchyApproval.HierarchyApprovalID))
			}

			hierarchyApp := model.HierarchyApproval{
				SetupFor:              setupFor,
				HierarchyApprovalType: request.ApprovalType,
				CreatedBy:             &request.UserID,
			}

			err = service.HierarchyApprovalRepository.Store(txCtx, &hierarchyApp)
			if err != nil {
				return err
			}

			var employeeValidationMap = entity.TempEmployeeValidationMap{}

			for _, detail := range request.Details {
				hierarchyAppDetail := model.HierarchyApprovalDet{
					HierarchyApprovalID:        *hierarchyApp.HierarchyApprovalID,
					Level:                      detail.Level,
					MaxOverLimit:               detail.MaxOverLimit,
					HierarchyApprovalDetCustID: detail.CompanyID,
					IsActive:                   *detail.IsActive,
				}

				err := service.HierarchyApprovalRepository.StoreDetail(txCtx, &hierarchyAppDetail)
				if err != nil {
					return err
				}

				var hierarchyAppDetailEmps []*model.HierarchyApprovalDetEmp
				if len(detail.EmpIDs) == 2 {
					if detail.EmpIDs[0].EmpID == detail.EmpIDs[1].EmpID {
						return errors.New("employee id index 0 can't be same with employee index 1")
					}
				}

				for _, empID := range detail.EmpIDs {
					level, err := employeeValidationMap.GetByID(empID.EmpID)
					if err == nil && level != nil {
						return errors.New(fmt.Sprintf("employee id : %v already set on level %v", empID.EmpID, *level))
					}
					_, err = service.HierarchyApprovalRepository.FindOneByEmployeeIdAndCustId(empID.EmpID, detail.CompanyID)
					if err != nil {
						return err
					}

					hierarchyAppDetailEmps = append(hierarchyAppDetailEmps, &model.HierarchyApprovalDetEmp{
						HierarchyApprovalDetailID: *hierarchyAppDetail.HierarchyApprovalDetailID,
						EmpID:                     empID.EmpID,
						Seq:                       empID.Sequence,
					})
					employeeValidationMap.SetTempEmployeeValidationMap(empID.EmpID, detail.Level)
				}

				err = service.HierarchyApprovalRepository.StoreDetailEmp(txCtx, hierarchyAppDetailEmps)
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

	return
}

func (service *HierarchyApprovalServiceImpl) List(dataFilter entity.HierarcyApprovalQueryFilter) (data []entity.HierarchyApprovalListResp, total int64, lastPage int, err error) {
	hierarchies, total, lastPage, err := service.HierarchyApprovalRepository.FindAllByCustID(dataFilter)
	if err != nil {
		return
	}

	for _, hierarchy := range hierarchies {
		data = append(data, entity.HierarchyApprovalListResp{
			HierarchyApprovalID: hierarchy.HierarchyApprovalID,
			SetupFor:            hierarchy.SetupFor,
			SetupForName:        hierarchy.CompanyName,
			SetupForCode:        hierarchy.CompanyCode,
			ApprovalType:        hierarchy.HierarchyApprovalType,
			UpdatedBy:           hierarchy.UpdatedBy,
			UpdatedByname:       hierarchy.UpdatedByname,
			UpdatedAt:           hierarchy.UpdatedAt,
		})
	}

	return
}

func (service *HierarchyApprovalServiceImpl) Detail(hierarchyApprovalID int64, custID string, parentCustID string) (response entity.HierarchyApprovalReadResp, err error) {
	hierarchyApproval, err := service.HierarchyApprovalRepository.FindByCustID(hierarchyApprovalID, custID, parentCustID)
	if err != nil {
		return
	}

	hierarchyApprovalDetail, err := service.HierarchyApprovalRepository.FindDetail(hierarchyApprovalID)
	if err != nil {
		return
	}

	response = entity.HierarchyApprovalReadResp{
		HierarchyApprovalID: hierarchyApproval.HierarchyApprovalID,
		SetupFor:            hierarchyApproval.SetupFor,
		SetupForName:        hierarchyApproval.CompanyName,
		SetupForCode:        hierarchyApproval.CompanyCode,
		ApprovalType:        hierarchyApproval.HierarchyApprovalType,
	}

	for _, hierarchyApprovalDet := range hierarchyApprovalDetail {
		detailEntity := entity.HierarchyApprovalDetailResp{
			IsActive:     hierarchyApprovalDet.IsActive,
			Level:        hierarchyApprovalDet.Level,
			CompanyID:    hierarchyApprovalDet.HierarchyApprovalDetCustID,
			CompanyName:  hierarchyApprovalDet.HierarchyApprovalDetCustIDName,
			MaxOverLimit: hierarchyApprovalDet.MaxOverLimit,
		}

		hierarchyApprovalDetailEmps, err := service.HierarchyApprovalRepository.FindDetailEmp(hierarchyApprovalDet.HierarchyApprovalDetailID)
		if err != nil {
			return response, err
		}

		for _, hierarchyApprovalDetailEmp := range hierarchyApprovalDetailEmps {
			detailEmpEntity := entity.HierarchyApprovalDetailEmpResp{
				Sequence: hierarchyApprovalDetailEmp.Seq,
				EmpID:    hierarchyApprovalDetailEmp.EmpID,
				EmpName:  hierarchyApprovalDetailEmp.EmpName,
			}

			detailEntity.EmpIDs = append(detailEntity.EmpIDs, detailEmpEntity)
		}

		response.Details = append(response.Details, detailEntity)

	}

	return
}

func (service *HierarchyApprovalServiceImpl) Delete(custID string, parentCustId string, hierarchyApprovalID int64, userID int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.HierarchyApprovalRepository.Delete(txCtx, custID, parentCustId, hierarchyApprovalID, userID)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *HierarchyApprovalServiceImpl) Update(hierarchyApprovalID int64, request entity.UpdateHierarchyApprovalBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		hierarchyModel := model.HierarchyApproval{
			SetupFor:              request.SetupFor,
			HierarchyApprovalType: request.ApprovalType,
			UpdatedBy:             &request.UserID,
		}

		err = service.HierarchyApprovalRepository.Update(txCtx, hierarchyApprovalID, hierarchyModel)
		if err != nil {
			return err
		}

		hierarchyApprovalDetails, err := service.HierarchyApprovalRepository.FindDetail(hierarchyApprovalID)
		if err != nil {
			return err
		}

		var hierarchyApprovalDetReadMap = model.HierarchyApprovalDetReadMap{}

		for _, hierarchyApprovalDetail := range hierarchyApprovalDetails {
			hierarchyApprovalDetReadMap.Set(hierarchyApprovalDetail.Level, hierarchyApprovalDetail)
		}

		for _, detail := range request.Details {
			hierarchyApprovalDetailModel, err := hierarchyApprovalDetReadMap.GetByLevel(detail.Level)
			if err != nil { // jika tidak ditemukan, insert barus
				var employeeValidationMap = entity.TempEmployeeValidationMap{}

				hierarchyAppDetail := model.HierarchyApprovalDet{
					HierarchyApprovalID:        hierarchyApprovalID,
					Level:                      detail.Level,
					MaxOverLimit:               detail.MaxOverLimit,
					HierarchyApprovalDetCustID: detail.CompanyID,
					IsActive:                   *detail.IsActive,
				}

				err := service.HierarchyApprovalRepository.StoreDetail(txCtx, &hierarchyAppDetail)
				if err != nil {
					return err
				}

				var hierarchyAppDetailEmps []*model.HierarchyApprovalDetEmp
				if len(detail.EmpIDs) == 2 {
					if detail.EmpIDs[0].EmpID == detail.EmpIDs[1].EmpID {
						return errors.New("employee id index 0 can't be same with employee index 1")
					}
				}

				for _, empID := range detail.EmpIDs {
					level, err := employeeValidationMap.GetByID(empID.EmpID)
					if err == nil && level != nil {
						return errors.New(fmt.Sprintf("employee id : %v already set on level %v", empID.EmpID, *level))
					}
					_, err = service.HierarchyApprovalRepository.FindOneByEmployeeIdAndCustId(empID.EmpID, detail.CompanyID)
					if err != nil {
						return err
					}

					hierarchyAppDetailEmps = append(hierarchyAppDetailEmps, &model.HierarchyApprovalDetEmp{
						HierarchyApprovalDetailID: *hierarchyAppDetail.HierarchyApprovalDetailID,
						EmpID:                     empID.EmpID,
						Seq:                       empID.Sequence,
					})
					employeeValidationMap.SetTempEmployeeValidationMap(empID.EmpID, detail.Level)
				}

				err = service.HierarchyApprovalRepository.StoreDetailEmp(txCtx, hierarchyAppDetailEmps)
				if err != nil {
					return err
				}

				continue
			}

			hierarchyDetModel := model.HierarchyApprovalUpdate{
				HierarchyApprovalDetCustID: detail.CompanyID,
				IsActive:                   *detail.IsActive,
				MaxOverLimit:               detail.MaxOverLimit,
				UpdatedBy:                  &request.UserID,
			}

			err = service.HierarchyApprovalRepository.UpdateDetail(txCtx, hierarchyApprovalDetailModel.HierarchyApprovalDetailID, hierarchyDetModel)
			if err != nil {
				return err
			}

			hierarchyApprovalDetailEmps, err := service.HierarchyApprovalRepository.FindDetailEmp(hierarchyApprovalDetailModel.HierarchyApprovalDetailID)
			if err != nil {
				return err
			}

			var hierarchyApprovalDetEmpReadMap = model.HierarchyApprovalDetEmpReadMap{}

			for _, hierarchyApprovalDetailEmp := range hierarchyApprovalDetailEmps {
				hierarchyApprovalDetEmpReadMap.Set(hierarchyApprovalDetailEmp.Seq, hierarchyApprovalDetailEmp)
			}

			if len(detail.EmpIDs) == 1 { // hapus emp 2 jika hanya berisi 1
				service.HierarchyApprovalRepository.DeleteDetailEmp(txCtx, hierarchyApprovalDetailModel.HierarchyApprovalDetailID, 2)
			}

			for _, empID := range detail.EmpIDs {
				hierarchyApprovalDetaiEmplModel, err := hierarchyApprovalDetEmpReadMap.GetBySequence(empID.Sequence)
				if err != nil { // jika tidak ditemukan, berarti belum ada record. tambah record baru di DB

					var hierarchyAppDetailEmps []*model.HierarchyApprovalDetEmp
					hierarchyAppDetailEmps = append(hierarchyAppDetailEmps, &model.HierarchyApprovalDetEmp{
						HierarchyApprovalDetailID: hierarchyApprovalDetailModel.HierarchyApprovalDetailID,
						EmpID:                     empID.EmpID,
						Seq:                       empID.Sequence,
					})

					err = service.HierarchyApprovalRepository.StoreDetailEmp(txCtx, hierarchyAppDetailEmps)
					if err != nil {
						return err
					}

					continue
				}

				hierarchyApprovalDetEmpUpdate := model.HierarchyApprovalDetEmpUpdate{
					EmpID: empID.EmpID,
					Seq:   empID.Sequence,
				}

				err = service.HierarchyApprovalRepository.UpdateDetailEmp(txCtx, hierarchyApprovalDetaiEmplModel.HierarchyApprovalDetailEmpID, hierarchyApprovalDetEmpUpdate)
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
	return
}

func (service *HierarchyApprovalServiceImpl) RequestApproval(request entity.RequestApprovalBody) (err error) {

	ro, err := service.OrderRepository.FindByNo(request.RoNo, request.CustID)
	if err != nil {
		return err
	}

	if ro.DataStatus != nil {
		if *ro.DataStatus != entity.NEED_REVIEW {
			return errors.New("order must be 'NEED REVIEW' for request approval")
		}
	}

	orderApproval, err := service.OrderApprovalRequest.FindApprovalProcessedByRoNo(request.RoNo, request.CustID)
	if err == nil {
		return errors.New(fmt.Sprintf("ro no %v already request approval with ID = %v", request.RoNo, orderApproval.OrderApprovalRequestID))
	}

	setupHierarchyApproval, err := service.HierarchyApprovalRepository.FindBySetupForOnly(ro.CustID)
	if err != nil {
		return err
	}

	setupHierarchyApprovalDetails, err := service.HierarchyApprovalRepository.FindDetail(setupHierarchyApproval.HierarchyApprovalID)
	if err != nil {
		return err
	}

	if setupHierarchyApproval.HierarchyApprovalType == entity.APPROVAL_TYPE_DIRECT {
		var setupHierarchyApprovalDetIndex int

		for i := len(setupHierarchyApprovalDetails); 0 < i; i-- {
			setupHierarchyApprovalDetail := setupHierarchyApprovalDetails[i-1]
			if i > 1 {

				if *setupHierarchyApprovalDetail.MaxOverLimit > *ro.TotalFinal {

					setupHierarchyApprovalDetIndex = i - 1

					break
				}
			} else {
				setupHierarchyApprovalDetIndex = i - 1
			}
		}
		c := context.Background()

		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
			setupHierarchyApprovalDet := setupHierarchyApprovalDetails[setupHierarchyApprovalDetIndex]

			orderApprovalRequestModel := model.OrderApprovalRequest{
				CustID:    request.CustID,
				RoNo:      request.RoNo,
				CreatedBy: &request.UserID,
			}

			err := service.OrderApprovalRequest.Store(txCtx, &orderApprovalRequestModel)
			if err != nil {
				return err
			}

			hierarchyApprovalDetailEmps, err := service.HierarchyApprovalRepository.FindDetailEmp(setupHierarchyApprovalDet.HierarchyApprovalDetailID)
			if err != nil {
				return err
			}

			for _, hierarchyApprovalDetailEmp := range hierarchyApprovalDetailEmps {
				orderApprovalRequestDetModel := model.OrderApprovalRequestDetail{
					OrderApprovalRequestID: *orderApprovalRequestModel.OrderApprovalRequestID,
					EmpID:                  hierarchyApprovalDetailEmp.EmpID,
					Seq:                    hierarchyApprovalDetailEmp.Seq,
					Level:                  setupHierarchyApprovalDet.Level,
				}

				err := service.OrderApprovalRequest.StoreDetail(txCtx, &orderApprovalRequestDetModel)
				if err != nil {
					return err
				}
			}

			return nil
		})
		return
	} else {

		var setupHierarchyApprovalDetIndexs []int
		var approvalLimitFound bool
		for i := len(setupHierarchyApprovalDetails); 0 < i; i-- {
			setupHierarchyApprovalDetail := setupHierarchyApprovalDetails[i-1]
			if approvalLimitFound {
				setupHierarchyApprovalDetIndexs = append(setupHierarchyApprovalDetIndexs, i-1)
				continue
			}
			if i > 1 {
				if *setupHierarchyApprovalDetail.MaxOverLimit < *ro.TotalFinal {
					setupHierarchyApprovalDetIndexs = append(setupHierarchyApprovalDetIndexs, i-1)
					approvalLimitFound = true
				}
			} else {
				setupHierarchyApprovalDetIndexs = append(setupHierarchyApprovalDetIndexs, i-1)
			}
		}

		c := context.Background()

		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
			orderApprovalRequestModel := model.OrderApprovalRequest{
				CustID:    request.CustID,
				RoNo:      request.RoNo,
				CreatedBy: &request.UserID,
			}

			err := service.OrderApprovalRequest.Store(txCtx, &orderApprovalRequestModel)
			if err != nil {
				return err
			}

			for _, setupHierarchyApprovalDetIndex := range setupHierarchyApprovalDetIndexs {
				setupHierarchyApprovalDet := setupHierarchyApprovalDetails[setupHierarchyApprovalDetIndex]

				hierarchyApprovalDetailEmps, err := service.HierarchyApprovalRepository.FindDetailEmp(setupHierarchyApprovalDet.HierarchyApprovalDetailID)
				if err != nil {
					return err
				}

				for _, hierarchyApprovalDetailEmp := range hierarchyApprovalDetailEmps {
					orderApprovalRequestDetModel := model.OrderApprovalRequestDetail{
						OrderApprovalRequestID: *orderApprovalRequestModel.OrderApprovalRequestID,
						EmpID:                  hierarchyApprovalDetailEmp.EmpID,
						Seq:                    hierarchyApprovalDetailEmp.Seq,
						Level:                  setupHierarchyApprovalDet.Level,
					}

					err := service.OrderApprovalRequest.StoreDetail(txCtx, &orderApprovalRequestDetModel)
					if err != nil {
						return err
					}
				}
			}
			return nil
		})

	}

	return
}

func (service *HierarchyApprovalServiceImpl) RequestApprovalDetail(custID string, requestApprovalID int64) (detail entity.ApprovalDetail, err error) {
	orderApproval, err := service.OrderApprovalRequest.FindOneByID(requestApprovalID, custID)
	if err != nil {
		return detail, err
	}

	orderApprovalDetails, err := service.OrderApprovalRequest.FindDetail(orderApproval.OrderApprovalRequestID)
	if err != nil {
		return detail, err
	}

	groupedDetails := make(map[int][]entity.RequestApprovalDetail)

	for _, orderApprovalDetail := range orderApprovalDetails {
		detail := entity.RequestApprovalDetail{
			EmployeeId:       orderApprovalDetail.EmployeeId,
			EmployeeCode:     orderApprovalDetail.EmployeeCode,
			EmployeeName:     orderApprovalDetail.EmployeeName,
			Sequence:         orderApprovalDetail.Seq,
			Status:           orderApprovalDetail.Status,
			ActDate:          orderApprovalDetail.ActDate,
			EmployeeImageURL: orderApprovalDetail.ImageURL,
		}

		groupedDetails[orderApprovalDetail.Level] = append(groupedDetails[orderApprovalDetail.Level], detail)
	}

	detail.FinishedAt = orderApproval.FinishedAt
	detail.RoNo = orderApproval.RoNo

	for level, groupedDetail := range groupedDetails {
		detail.Approvals = append(detail.Approvals, entity.GroupedApproval{
			Level:   level,
			Details: groupedDetail,
		})
	}
	return
}
