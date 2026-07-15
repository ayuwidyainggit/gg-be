package service

import (
	"context"
	"errors"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/constant"
	"mobile/repository"

	"github.com/gofiber/fiber/v2/log"
)

type (
	PjpService interface {
		GetSalesmanDetail(empId int64, isDistributor bool) (response entity.PjpSalesmanResponse, err error)
		SubmitPjpDistributor(ctx context.Context, request entity.SubmitPjpDistributorRequest, custId string) error
		UpdatePJPDistributor(ctx context.Context, payload entity.UpdatePJPDistributorRequest) error
	}

	pjpServiceImpl struct {
		MSalesmanRepository      repository.MSalesmanRepository
		MWarehouseRepository     repository.MWarehouseRepository
		PjpDistributorRepository repository.PjpDistributorRepository
		PjpPrincipleRepository   repository.PjpPrincipalRepository
		Transaction              repository.Dbtransaction
	}
)

func NewPjpService(
	mSalesmanRepository repository.MSalesmanRepository,
	mWarehouseRepository repository.MWarehouseRepository,
	pjpDistributorRepository repository.PjpDistributorRepository,
	pjpPrincipalRepository repository.PjpPrincipalRepository,
	transaction repository.Dbtransaction,
) PjpService {
	return &pjpServiceImpl{
		MSalesmanRepository:      mSalesmanRepository,
		MWarehouseRepository:     mWarehouseRepository,
		PjpDistributorRepository: pjpDistributorRepository,
		PjpPrincipleRepository:   pjpPrincipalRepository,
		Transaction:              transaction,
	}
}

// GetSalesmanDetail retrieves PJP salesman detail information
func (service *pjpServiceImpl) GetSalesmanDetail(empId int64, isDistributor bool) (response entity.PjpSalesmanResponse, err error) {
	detail, err := service.MSalesmanRepository.FindPjpSalesmanDetail(empId)
	if err != nil {
		log.Error("PjpService, GetSalesmanDetail, FindPjpSalesmanDetail, err:", err.Error())
		return response, errors.New("record not found")
	}

	if detail.CustId == "" {
		return response, errors.New(constant.STATUS_DB_NOT_FOUND)
	}

	// get pjp info
	var pjp model.PermanentJourneyPlan
	if isDistributor {
		pjp, err = service.PjpDistributorRepository.GetPJPInfo(context.Background(), empId)
		if err != nil {
			log.Error("PjpService, GetSalesmanDetail, PjpDistributorRepository.GetPJPInfo, err:", err.Error())
		}
	} else {
		pjp, err = service.PjpPrincipleRepository.GetPJPInfo(context.Background(), empId)
		if err != nil {
			log.Error("PjpService, GetSalesmanDetail, PjpDistributorRepository.GetPJPInfo, err:", err.Error())
		}
	}

	var (
		wh       model.MWarehouse
		whCanvas model.MWarehouse
	)
	if detail.WHID != 0 {
		wh, err = service.MWarehouseRepository.FindOneByID(detail.WHID)
		if err != nil {
			log.Error("PjpService, GetSalesmanDetail, MWarehouseRepository.FindOneByID, err:", err.Error())
			return response, errors.New("record not found")
		}
	}

	if detail.WHIDCanvas != 0 {
		whCanvas, err = service.MWarehouseRepository.FindOneByID(detail.WHIDCanvas)
		if err != nil {
			log.Error("PjpService, GetSalesmanDetail, MWarehouseRepository.FindOneByID, err:", err.Error())
			return response, errors.New("record not found")
		}
	}

	response.CustId = detail.CustId
	response.CustName = detail.CustName
	response.DistributorId = detail.DistributorId // null if principal, not null if distributor
	response.DistributorCode = detail.DistributorCode
	response.DistributorName = detail.DistributorName
	response.EmpId = detail.EmpId
	response.SalesName = detail.SalesName
	response.SalesTeamId = detail.SalesTeamId
	response.SalesTeamCode = detail.SalesTeamCode
	response.SalesTeamName = detail.SalesTeamName

	if pjp.PjpCode > 0 {
		response.PJPCode = &pjp.PjpCode
	}

	if pjp.ApprovalStatus != nil {
		response.PJPStatus = pjp.ApprovalStatus
	}

	response.Data = map[string]any{
		"opr_type":        "",
		"wh_id":           0,
		"wh_code":         wh.WHCode,
		"wh_name":         wh.WHName,
		"wh_id_canvas":    0,
		"wh_code_canvas":  whCanvas.WHCode,
		"wh_name_canvas":  whCanvas.WHName,
		"opr_type_canvas": "",
	}

	if detail.IsTakingOrder {
		response.Data["opr_type"] = detail.OprType
		response.Data["wh_id"] = detail.WHID
	}

	if detail.IsActiveSalesmanCanvas {
		response.Data["opr_type_canvas"] = detail.OprTypeCanvas
		response.Data["wh_id_canvas"] = detail.WHIDCanvas
	}

	// for _, wh := range warehouses {
	// 	response.Data = append(response.Data, entity.PjpSalesmanWarehouseData{
	// 		OprType: wh.OprType,
	// 		WhId:    wh.WhId,
	// 		WhCode:  wh.WhCode,
	// 		WhName:  wh.WhName,
	// 	})
	// }

	return response, nil
}
