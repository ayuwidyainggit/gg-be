package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/constant"
	"mobile/pkg/structs"
	"mobile/repository"
	"strconv"
	"time"
)

type MOutletService interface {
	MOutletCheck(ctx context.Context, year int, cust_id string, status []string) (entity.MOutletCodeResponse, error)
	Store(context.Context, string, string, int64, entity.CreateMOutletBody) (entity.MOutletRespone, error)
	StoreFromList(context.Context, entity.ExtraCallOutlet) (entity.MOutletRespone, error)
	List(entity.MOutletQueryFilter, string, string, int64) (data []entity.MOutletListRespone, total int64, lastPage int, err error)
	ListOutletAdditionals(entity.MOutletQueryFilter, string, int64, string) (data []entity.MOutletListRespone, total int64, lastPage int, err error)
	Delete(string, int, int64) error
	Detail(int64, int64, string, string) (entity.OutletRespone, error)
	MobileOutletList(entity.MobileOutletListQueryFilter, string, int64) (data []entity.MobileOutletListResponse, total int64, lastPage int, err error)
	MobileOutletDetail(int64, string) (entity.MobileOutletDetailResponse, error)
	OutletPJPList(entity.OutletPJPListQuery) (data entity.OutletPJPListResponse, paging entity.Pagination, err error)
	GetRegionByDistributorID(ctx context.Context, distributorID int) (model.Region, error)
}

func NewMOutletService(
	mOutletRepository repository.MOutletRepository,
	pjpDistributorRepository repository.PjpDistributorRepository,
	pjpPrincipalRepository repository.PjpPrincipalRepository,
	invoicesRepository repository.InvoicesRepository) *mOutletServiceImpl {
	return &mOutletServiceImpl{
		MOutletRepository:        mOutletRepository,
		PjpDistributorRepository: pjpDistributorRepository,
		PjpPrincipalRepository:   pjpPrincipalRepository,
		InvoicesRepository:       invoicesRepository,
	}
}

type mOutletServiceImpl struct {
	MOutletRepository        repository.MOutletRepository
	PjpDistributorRepository repository.PjpDistributorRepository
	PjpPrincipalRepository   repository.PjpPrincipalRepository
	InvoicesRepository       repository.InvoicesRepository
}

func (service *mOutletServiceImpl) Store(ctx context.Context, custID string, distributorCode string, distributorID int64, request entity.CreateMOutletBody) (response entity.MOutletRespone, err error) {
	y, _, _ := time.Now().Date()
	// get last 2 digit of year
	year := y % 100

	outletCodeConf, err := service.MOutletRepository.MOutletCheck(ctx, year, custID, []string{"Active"})
	if request.OutletCode == "" && err != nil && errors.Is(err, sql.ErrNoRows) {
		return response, errors.New("outlet code configuration not found for cust_id: " + custID + ", year: " + strconv.Itoa(year))
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return response, err
	}

	outlet, err := service.MOutletRepository.FindOneByOutletCodeAndCustId(request.OutletName, request.CustId, request.ParentCustId)
	if err == nil {
		return response, errors.New("outlet_name: " + *outlet.OutletName + " is already exists")
	}

	outletStatus, err := service.MOutletRepository.FindTopStatus(ctx, request.ParentCustId)
	if err != nil {
		return response, errors.New("error find top status")
	}

	// Get region info for outlet_principal_code generation (only if distributorID is provided)
	var region model.Region
	if distributorID > 0 {
		region, err = service.MOutletRepository.FindRegionByDistributorID(ctx, int(distributorID))
		if err != nil {
			return response, errors.New("distributor region not found")
		}
	}

	var (
		outletData         model.MOutlet
		timeNow            = time.Now().In(time.UTC)
		verificationStatus = constant.ApproveVerificationStatus
		sourceVal          = request.Source
	)

	err = structs.Automapper(request, &outletData)
	if err != nil {
		return response, err
	}

	outletData.VerificationStatus = &verificationStatus
	outletData.OutletStatus = outletStatus
	outletData.CreatedAt = &timeNow
	outletData.CreatedBy = &request.CreatedBy
	outletData.UpdatedBy = &request.UpdatedBy
	outletData.UpdatedAt = &timeNow
	outletData.FileUrl = &request.FileUrl
	outletData.Source = &sourceVal
	outletData.IsActive = true

	const (
		nonDiscountGroupID = 140
		nonOutletGroupID   = 114
		nonPriceGroupID    = 51
		nonOutletClassID   = 74
		nonOutletTypeID    = 100
	)

	outletData.NonDiscountGroupID = nonDiscountGroupID
	outletData.NonOutletGroupID = nonOutletGroupID
	outletData.NonPriceGroupID = nonPriceGroupID
	outletData.NonOutletClassID = nonOutletClassID
	outletData.NonOutletTypeID = nonOutletTypeID

	trx, err := service.MOutletRepository.TrxBegin()
	if err != nil {
		trx.TrxRollback()
		return response, err
	}

	// Generate outlet code within transaction
	// Check if outletCodeConf has a non-zero ID (meaning config exists)
	if outletCodeConf.Id != "" {
		conf := outletCodeConf
		// Parse last_sequence_no to int, increment by 1
		lastSeq, err := strconv.Atoi(conf.LastSequenceNo)
		if err != nil {
			trx.TrxRollback()
			return response, errors.New("invalid last_sequence_no format")
		}
		newSeq := lastSeq + 1
		// Format as 4 digit with leading zeros
		newSeqStr := fmt.Sprintf("%04d", newSeq)
		// Generate outlet code: serial_code + year_code (2 digit) + new sequence
		yearCodeStr := fmt.Sprintf("%02d", conf.YearCode)
		generatedCode := conf.SerialCode + yearCodeStr + newSeqStr

		// Check if the generated code already exists before using it
		outletCodeExists, err := service.MOutletRepository.CheckOutletByCode(ctx, generatedCode, request.CustId)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		if outletCodeExists {
			trx.TrxRollback()
			return response, fmt.Errorf("outlet_code: generated outlet code has already taken")
		}
		outletData.OutletCode = &generatedCode
		// Update last_sequence_no in m_outlet_code
		err = trx.UpdateOutletCodeSequence(conf.Id, newSeqStr, request.UpdatedBy)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}

		// Generate outlet_principal_code: region_code + distributor_code + DDMMYY + 3 digit sequential
		// Only generate if distributorID is provided
		if distributorID > 0 {
			now := time.Now()
			day := now.Format("02")
			month := now.Format("01")
			yy := now.Format("06")
			ddmmyy := day + month + yy
			// Use the same newSeq but format as 3 digits
			principalSeqStr := fmt.Sprintf("%03d", newSeq)
			principalCode := region.RegionCode + distributorCode + ddmmyy + principalSeqStr
			outletData.OutletPrincipalCode = &principalCode
		}
	} else if request.OutletCode != "" {
		// If no config found, use the outlet_code sent by user
		// Check if outlet code is already used
		outletCodeExists, err := service.MOutletRepository.CheckOutletByCode(ctx, request.OutletCode, request.CustId)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}

		if outletCodeExists {
			trx.TrxRollback()
			return response, fmt.Errorf("outlet_code: outlet code has already taken")
		}
		outletData.OutletCode = &request.OutletCode
	}

	err = trx.Store(&outletData)
	if err != nil {
		trx.TrxRollback()
		return response, err
	}

	response.OutletId = outletData.OutletId

	for _, detail := range request.Details.OutletContact {
		var outletContactData model.MOutletContact
		err = structs.Automapper(detail, &outletContactData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		outletContactData.CustID = request.CustId
		outletContactData.OutletID = int64(outletData.OutletId)
		err := trx.StoreDetailContact(&outletContactData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
	}

	trx.TrxCommit()

	return response, err
}

func (service *mOutletServiceImpl) StoreFromList(ctx context.Context, request entity.ExtraCallOutlet) (response entity.MOutletRespone, err error) {
	if request.RouteCode == "" || request.PJPCode == "" || request.PJPID == 0 {

		if request.IsDistributor {
			pjpD, err := service.PjpDistributorRepository.FindPJPByEmpIDAndCustID(ctx, request.EmpID, request.CustID)
			if err != nil {
				return response, err
			}

			if request.RouteCode == "" {
				request.RouteCode = pjpD.RouteCode
			}

			if request.PJPCode == "" {
				request.PJPCode = pjpD.PJPCode
			}

			if request.PJPID == 0 {
				request.PJPID = int(pjpD.PJPID)
			}

		} else {
			pjpD, err := service.PjpPrincipalRepository.FindPJPByEmpIDAndCustID(ctx, request.EmpID, request.CustID)
			if err != nil {
				return response, err
			}

			if request.RouteCode == "" {
				request.RouteCode = pjpD.RouteCode
			}

			if request.PJPCode == "" {
				request.PJPCode = pjpD.PJPCode
			}

			if request.PJPID == 0 {
				request.PJPID = int(pjpD.PJPID)
			}
		}
	}

	outlets, err := service.MOutletRepository.FindByIdOutletFromList(ctx, request.OutletIDs)
	if err != nil {
		return response, err
	}

	if len(outlets) == 0 {
		return response, errors.New("outlet not found")
	}

	trx, err := service.MOutletRepository.TrxBegin()
	if err != nil {
		return response, err
	}

	defer func() {
		if err != nil {
			trx.TrxRollback()
		}
	}()

	routeCode, err := strconv.Atoi(request.RouteCode)
	if err != nil {
		return response, fmt.Errorf("invalid route code")
	}

	route := model.Route{}
	if request.IsDistributor {
		route, err = service.PjpDistributorRepository.GetRouteByRouteCode(ctx, routeCode)
		if err != nil {
			return response, err
		}
	} else {
		route, err = service.PjpPrincipalRepository.GetRouteByRouteCode(ctx, routeCode)
		if err != nil {
			return response, err
		}
	}

	if int(route.PjpId) != request.PJPID {
		return response, errors.New("route not found")
	}

	pjpCode, err := strconv.Atoi(request.PJPCode)
	if err != nil {
		return response, fmt.Errorf("invalid pjp code")
	}

	for _, rows := range outlets {
		var (
			isAdditional         = true
			isCurrentYear        = true
			isExtraCall          = true
			timeNow              = time.Now().In(time.UTC)
			startWeek            = timeNow.Format(time.DateOnly)
			startOfYear          = time.Date(timeNow.Year(), 1, 1, 0, 0, 0, 0, timeNow.Location())
			daysSinceStartOfYear = int(timeNow.Sub(startOfYear).Hours() / 24)
			week                 = (daysSinceStartOfYear+int(startOfYear.Weekday()))/7 + 1
			yearNow              = timeNow.Year()
		)

		outletData := model.MOutletCreadFromList{
			RouteCode:       &routeCode,
			RouteName:       &route.RouteName,
			OutletId:        rows.OutletID,
			OutletCode:      rows.OutletCode,
			OutletName:      rows.OutletName,
			Longitude:       rows.Longitude,
			Latitude:        rows.Latitude,
			OutletStatus:    0,
			OutletAddress:   rows.Address,
			PjpId:           &request.PJPID,
			PjpCode:         &pjpCode,
			CustId:          request.CustID,
			Status:          "pending",
			CreatedAt:       &timeNow,
			UpdatedAt:       &timeNow,
			IsAdditional:    &isAdditional,
			StartWeek:       &startWeek,
			IsInCurrentYear: &isCurrentYear,
			Week:            &week,
			Year:            &yearNow,
			Date:            &timeNow,
			IsExtraCall:     &isExtraCall,
		}

		var (
			onlyDateVisitList   = time.Date(outletData.Date.Year(), outletData.Date.Month(), outletData.Date.Day(), 0, 0, 0, 0, outletData.Date.Location())
			isPlannedVisitList  = false
			outletDataVisitList = model.MOutletCreadFromListOutletVisitList{
				Week:        outletData.Week,
				Year:        outletData.Year,
				Date:        &onlyDateVisitList,
				Day:         outletData.Date.Format("Mon"),
				RouteCode:   outletData.RouteCode,
				OutletCode:  outletData.OutletCode,
				PjpId:       outletData.PjpId,
				PjpCode:     outletData.PjpCode,
				CreatedAt:   outletData.CreatedAt,
				OutletId:    outletData.OutletId,
				IsPlanned:   &isPlannedVisitList,
				IsExtraCall: &isExtraCall,
				Latitude:    outletData.Latitude,
				Longitude:   outletData.Longitude,
			}
		)

		if request.IsDistributor {
			err = trx.StoreFromList(ctx, &outletData)
			if err != nil {
				return response, err
			}

			err = trx.StoreFromListOutletVisitList(ctx, &outletDataVisitList)
			if err != nil {
				return response, err
			}
		} else {
			err = trx.StoreFromListPrinciple(ctx, &outletData)
			if err != nil {
				return response, err
			}

			err = trx.StoreFromListOutletVisitListPrinciple(ctx, &outletDataVisitList)
			if err != nil {
				return response, err
			}
		}
	}

	err = trx.TrxCommit()
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *mOutletServiceImpl) List(dataFilter entity.MOutletQueryFilter, custId, parentCustId string, empId int64) (data []entity.MOutletListRespone, total int64, lastPage int, err error) {
	outlets, total, lastPage, err := service.MOutletRepository.FindAllByCustId(dataFilter, custId, parentCustId, empId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(outlets) > 0 {
		// Get all outlet names first to avoid queries in loop
		outletNames := make([]string, len(outlets))
		for i, row := range outlets {
			outletNames[i] = *row.OutletName
		}

		// Get outlet counts for all names in a single query
		outletCounts, err := service.MOutletRepository.GetOutletCountsByNames(outletNames)
		if err != nil {
			return data, total, lastPage, err
		}

		// Convert slice to map for efficient lookup
		outletCountsMap := make(map[string]model.MOutletFromListSimilar)
		for _, count := range outletCounts {
			outletCountsMap[count.OutletName] = count
		}

		for _, row := range outlets {
			var vResp entity.MOutletListRespone
			structs.Automapper(row, &vResp)

			verificationStatusName := vResp.GenerateOutletVerificationStatusName()
			vResp.VerificationStatusName = &verificationStatusName

			switch row.OutletStatus {
			case 1:
				vResp.OutletStatusName = "Active"
			case 2:
				vResp.OutletStatusName = "Covered"
			case 3:
				vResp.OutletStatusName = "Non Active"
			case 4:
				vResp.OutletStatusName = "Closed"
			default:
				vResp.OutletStatusName = "Unknown"
			}

			log.Println("OutletService, List, vResp:", vResp)

			// Use the outlet counts from the map if available
			if outletCount, exists := outletCountsMap[*row.OutletName]; exists && outletCount.Count == 1 {
			}
			data = append(data, vResp)

		}
	}
	return data, total, lastPage, err
}

func (service *mOutletServiceImpl) ListOutletAdditionals(dataFilter entity.MOutletQueryFilter, custId string, empId int64, parentCustId string) (data []entity.MOutletListRespone, total int64, lastPage int, err error) {
	outlets, total, lastPage, err := service.MOutletRepository.FindAllByCustIdFromList(dataFilter, custId, empId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(outlets) > 0 {
		for _, row := range outlets {
			var vResp entity.MOutletListRespone
			structs.Automapper(row, &vResp)

			// ADD OUTLET FRROM LIST
			if row.OutletStatus == 0 {
				unPlanned := "Un Planned"
				vResp.VerificationStatusName = &unPlanned
			}

			// REGISTER OUTLET
			if row.OutletStatus == 1 {
				switch *row.VerificationStatus {
				case 1:
					approved := "Approved"
					vResp.VerificationStatusName = &approved
				case 2:
					inReview := "In Review"
					vResp.VerificationStatusName = &inReview
				case 3:
					nonActive := "Non Active"
					vResp.VerificationStatusName = &nonActive
				default:
					unknown := "Unknown"
					vResp.VerificationStatusName = &unknown
				}
			}

			log.Println("OutletService, List, vResp:", vResp)

			data = append(data, vResp)

		}
	}
	return data, total, lastPage, err
}

func (service *mOutletServiceImpl) Delete(custId string, outletId int, userId int64) (err error) {

	err = service.MOutletRepository.Delete(custId, outletId, userId)
	if err != nil {
		return err
	}
	return err
}

func (service *mOutletServiceImpl) Detail(outletID, salesmanID int64, custId string, parentCustId string) (response entity.OutletRespone, err error) {
	outlet, err := service.MOutletRepository.FindOneByOutletIdAndCustId(outletID, custId, parentCustId)
	if err != nil {
		return response, err
	}
	err = structs.Automapper(outlet, &response)
	if err != nil {
		return response, err
	}

	log.Println("outlet.ArStatus:", outlet.ArStatus)
	log.Println("AR_TYPE_NORMAL_ID:", entity.AR_TYPE_NORMAL_ID)
	log.Println("AR_TYPE_1X_GIRO_TOLAKAN_ID:", entity.AR_TYPE_1X_GIRO_TOLAKAN_ID)
	log.Println("AR_TYPE_2X_GIRO_TOLAKAN_ID:", entity.AR_TYPE_2X_GIRO_TOLAKAN_ID)
	log.Println("AR_TYPE_3X_GIRO_TOLAKAN_ID:", entity.AR_TYPE_3X_GIRO_TOLAKAN_ID)

	if *outlet.ArStatus == entity.AR_TYPE_NORMAL_ID {
		response.ArStatusName = entity.AR_TYPE_NORMAL
	} else if *outlet.ArStatus == entity.AR_TYPE_1X_GIRO_TOLAKAN_ID {
		response.ArStatusName = entity.AR_TYPE_1X_GIRO_TOLAKAN
	} else if *outlet.ArStatus == entity.AR_TYPE_2X_GIRO_TOLAKAN_ID {
		response.ArStatusName = entity.AR_TYPE_2X_GIRO_TOLAKAN
	} else if *outlet.ArStatus == entity.AR_TYPE_3X_GIRO_TOLAKAN_ID {
		response.ArStatusName = entity.AR_TYPE_3X_GIRO_TOLAKAN
	}

	if outlet.CloseDate != nil {
		closeDate := outlet.CloseDate.Format("2006-01-02")
		response.CloseDate = closeDate
	}

	outletContactDetail, err := service.MOutletRepository.GetDetailOutletContact(outletID, custId)
	if err != nil {
		return response, err
	}
	for _, contactDetail := range outletContactDetail {
		var outletContactResp entity.MOutletContactRead
		err = structs.Automapper(contactDetail, &outletContactResp)
		if err != nil {
			return response, err
		}
		response.Details.OutletContact = append(response.Details.OutletContact, outletContactResp)
	}

	outletPayment := []entity.MOutletPayment{
		{PaymentType: "", CreditLimitType: "", CreditLimit: 0, SalesInvLimitType: "", SalesInvLimit: 0},
	}
	for _, outletPaymentDetail := range outletPayment {
		var outletPaymentResp entity.MOutletPayment
		err = structs.Automapper(outletPaymentDetail, &outletPaymentResp)
		if err != nil {
			return response, err
		}
		response.Details.OutletPayment = append(response.Details.OutletPayment, outletPaymentResp)
	}

	// Fetch remaining outstanding amount for the outlet
	remainingOutstanding, err := service.InvoicesRepository.GetRemainingOutstandingByOutletID(context.Background(), outletID)
	if err != nil {
		return response, err
	}
	response.RemainingOutstandingAmount = remainingOutstanding

	// Ambil tanggal saat ini
	// currentDate := time.Now()

	// newDate := currentDate.AddDate(0, 0, response.Top)

	// formattedDate := newDate.Format("2006-01-02")

	// response.Duedate = formattedDate

	// verificationStatusName := response.GenerateOutletVerificationStatusName()
	// response.VerificationStatusName = &verificationStatusName

	// response.TaxInvoiceFormName = response.GenerateTaxInvoiceFormName()

	return response, err
}

func (service *mOutletServiceImpl) MobileOutletList(dataFilter entity.MobileOutletListQueryFilter, custId string, empID int64) (data []entity.MobileOutletListResponse, total int64, lastPage int, err error) {
	if dataFilter.Page <= 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 5
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "outlet_code:asc"
	}

	outlets, total, lastPage, err := service.MOutletRepository.FindMobileOutletList(dataFilter, custId, empID)
	if err != nil {
		return data, 0, 0, err
	}

	data = make([]entity.MobileOutletListResponse, 0)
	for _, outlet := range outlets {
		data = append(data, entity.MobileOutletListResponse{
			OutletId:     outlet.OutletId,
			OutletCode:   outlet.OutletCode,
			OutletName:   outlet.OutletName,
			OutletStatus: outlet.OutletStatus,
			Address1:     outlet.Address1,
			Latitude:     outlet.Latitude,
			Longitude:    outlet.Longitude,
			AvgSalesWeek: outlet.AvgSalesWeek,
		})
	}

	return data, total, lastPage, nil
}

func (service *mOutletServiceImpl) MOutletCheck(ctx context.Context, year int, cust_id string, status []string) (entity.MOutletCodeResponse, error) {
	var response entity.MOutletCodeResponse
	outlet, err := service.MOutletRepository.MOutletCheck(ctx, year, cust_id, status)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return response, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return response, nil
	}

	err = structs.Automapper(outlet, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (service *mOutletServiceImpl) MobileOutletDetail(outletId int64, custId string) (entity.MobileOutletDetailResponse, error) {
	var response entity.MobileOutletDetailResponse

	outlet, err := service.MOutletRepository.FindMobileOutletDetail(outletId, custId)
	if err != nil {
		return response, err
	}

	response.OutletId = outlet.OutletId
	response.OutletCode = outlet.OutletCode
	response.OutletName = outlet.OutletName
	response.Address1 = outlet.Address1
	response.PhoneNo = outlet.PhoneNo
	response.BuildingOwn = outlet.BuildingOwn
	response.FileUrl = outlet.FileUrl
	response.Longitude = outlet.Longitude
	response.Latitude = outlet.Latitude

	contact, err := service.MOutletRepository.FindMobileOutletContact(outletId, custId)
	if err != nil {
		return response, err
	}

	if contact != nil && contact.ContactName != "" {
		response.OtherContact = &entity.MobileOutletContactDetail{
			ContactName: contact.ContactName,
			JobTitle:    contact.JobTitle,
			PhoneNo:     contact.PhoneNo,
			WaNo:        contact.WaNo,
			Email:       contact.Email,
		}
	}

	return response, nil
}

func (service *mOutletServiceImpl) OutletPJPList(params entity.OutletPJPListQuery) (data entity.OutletPJPListResponse, paging entity.Pagination, err error) {

	if params.Page <= 0 {
		params.Page = 1
	}

	if params.Limit <= 0 || params.Limit > 999 {
		params.Limit = 5
	}

	if params.Sort == "" {
		params.Sort = "outlet_code"
	}

	if params.SortOrder == "" {
		params.SortOrder = "asc"
	}

	var (
		resp             entity.OutletPJPListResponse
		respOutlets      = make([]entity.OutletPJPResponse, 0)
		respDistributors = make([]entity.DistributorPJPResponse, 0)
	)

	if params.DistributorID > 0 {
		outlets, total, err := service.PjpDistributorRepository.FindAllOutletByDate(context.Background(), params)
		paging = entity.Pagination{
			TotalRecord: total,
			PageCurrent: params.Page,
			PageLimit:   params.Limit,
			PageTotal:   int(math.Ceil(float64(total) / float64(params.Limit))),
		}
		if err != nil {
			return resp, paging, err
		}

		for _, outlet := range outlets {
			resp.RouteName = outlet.RouteName
			resp.RouteCode = outlet.RouteCode
			resp.Week = outlet.Week
			resp.Year = outlet.Year
			resp.Date = outlet.Date
			if outlet.OutletID == 0 {
				continue
			}

			respOutlets = append(respOutlets, entity.OutletPJPResponse{
				OutletID:      outlet.OutletID,
				OutletCode:    outlet.OutletCode,
				OutletName:    outlet.OutletName,
				OutletStatus:  outlet.OutletStatus,
				OutletAddress: outlet.OutletAddress,
				Latitude:      outlet.Latitude,
				Longitude:     outlet.Longitude,
			})
		}

	} else {
		outlets, err := service.PjpPrincipalRepository.FindAllOutletByDate(context.Background(), params, "outlet")
		if err != nil {
			return resp, paging, err
		}

		for _, outlet := range outlets {
			resp.RouteName = outlet.RouteName
			resp.RouteCode = outlet.RouteCode
			resp.Week = outlet.Week
			resp.Year = outlet.Year
			resp.Date = outlet.Date
			if outlet.DestinationID == 0 {
				continue
			}
			respOutlets = append(respOutlets, entity.OutletPJPResponse{
				OutletID:      outlet.DestinationID,
				OutletCode:    outlet.DestinationCode,
				OutletName:    outlet.DestinationName,
				OutletStatus:  outlet.DestinationStatus,
				OutletAddress: outlet.DestinationAddress,
				Latitude:      outlet.Latitude,
				Longitude:     outlet.Longitude,
			})
		}

		distributors, err := service.PjpPrincipalRepository.FindAllOutletByDate(context.Background(), params, "distributor")
		if err != nil {
			return resp, paging, err
		}

		for _, distributor := range distributors {
			if len(outlets) == 0 {
				resp.RouteName = distributor.RouteName
				resp.RouteCode = distributor.RouteCode
				resp.Week = distributor.Week
				resp.Year = distributor.Year
				resp.Date = distributor.Date
			}
			if distributor.DestinationID == 0 {
				continue
			}

			respDistributors = append(respDistributors, entity.DistributorPJPResponse{
				DistributorID:      distributor.DestinationID,
				DistributorCode:    distributor.DestinationCode,
				DistributorName:    distributor.DestinationName,
				DistributorStatus:  distributor.DestinationStatus,
				DistributorAddress: distributor.DestinationAddress,
				Latitude:           distributor.Latitude,
				Longitude:          distributor.Longitude,
			})
		}

		totalCount, err := service.PjpPrincipalRepository.CountAllOutletByDate(context.Background(), params)
		if err != nil {
			return resp, paging, err
		}
		paging = entity.Pagination{
			TotalRecord: totalCount,
			PageCurrent: params.Page,
			PageLimit:   params.Limit,
			PageTotal:   int(math.Ceil(float64(totalCount) / float64(params.Limit))),
		}
	}

	resp.Outlets = respOutlets
	resp.Distributors = respDistributors
	return resp, paging, nil
}

func (service *mOutletServiceImpl) GetRegionByDistributorID(ctx context.Context, distributorID int) (model.Region, error) {
	region, err := service.MOutletRepository.FindRegionByDistributorID(ctx, distributorID)
	if err != nil {
		return region, err
	}

	return region, nil
}
