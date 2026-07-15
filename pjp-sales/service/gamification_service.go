package service

import (
	"context"
	"encoding/json"
	"sales/entity"
	"sales/model"
	"sales/pkg/str"
	"sales/pkg/structs"
	"sales/repository"
	"sort"
	"time"
)

type GamificationService interface {
	List(dataFilter entity.GamificationQueryFilter) (data []entity.GamificationListResponse, total int64, lastPage int, err error)
	Store(request entity.CreateGamificationBody) (err error)
	Detail(gamificationNo string, custID string, parentCustID string) (response entity.GamificationResponse, err error)
	Update(gamificationNo string, request entity.UpdateGamificationBody) (err error)
	Cancel(gamificationNo string, request entity.CancelGamificationBody) (err error)
	Stop(gamificationNo string, request entity.StopGamificationBody) (err error)

	GamificationStatusLookupList(entity.GeneralQueryFilter) (data []entity.GamificationStatusesLookupResponse, total int64, lastPage int, err error)
	EmployeeGroupLookupList(entity.GeneralQueryFilter) (data []entity.EmployeeGroupsLookupResponse, total int64, lastPage int, err error)
	SalesTeamLookupList(entity.GeneralQueryFilter) (data []entity.SalesTeamsLookupResponse, total int64, lastPage int, err error)
	OperationTypeLookupList(entity.GeneralQueryFilter) (data []entity.OperationTypesLookupResponse, total int64, lastPage int, err error)
	WarehouseLookupList(entity.GeneralQueryFilter) (data []entity.WarehousesLookupResponse, total int64, lastPage int, err error)
	SalesmanLookupList(entity.SalesmanQueryFilter) (data []entity.SalesmansLookupResponse, total int64, lastPage int, err error)
	ProductCategoryLookupList(entity.GeneralQueryFilter) (data []entity.ProductCategoriesLookupResponse, total int64, lastPage int, err error)
	ProductLineLookupList(entity.GeneralQueryFilter) (data []entity.ProductLinesLookupResponse, total int64, lastPage int, err error)
	BrandLookupList(entity.BrandQueryFilter) (data []entity.BrandsLookupResponse, total int64, lastPage int, err error)
	SubBrand1LookupList(entity.SubBrand1QueryFilter) (data []entity.SubBrands1LookupResponse, total int64, lastPage int, err error)
	SubBrand2LookupList(entity.GeneralQueryFilter) (data []entity.SubBrands2LookupResponse, total int64, lastPage int, err error)
	ProductLookupList(entity.ProductQueryFilter) (data []entity.ProductsLookupResponse, total int64, lastPage int, err error)
	OutletLocationLookupList(entity.GeneralQueryFilter) (data []entity.OutletLocationsLookupResponse, total int64, lastPage int, err error)
	OutletGroupLookupList(entity.GeneralQueryFilter) (data []entity.OutletGroupsLookupResponse, total int64, lastPage int, err error)
	OutletClassLookupList(entity.GeneralQueryFilter) (data []entity.OutletClassesLookupResponse, total int64, lastPage int, err error)
	MarketLookupList(entity.GeneralQueryFilter) (data []entity.MarketsLookupResponse, total int64, lastPage int, err error)
	DistrictLookupList(entity.GeneralQueryFilter) (data []entity.DistrictsLookupResponse, total int64, lastPage int, err error)
	IndustryLookupList(entity.GeneralQueryFilter) (data []entity.IndustriesLookupResponse, total int64, lastPage int, err error)
	OutletLookupList(entity.OutletQueryFilter) (data []entity.OutletsLookupResponse, total int64, lastPage int, err error)
	SourceLookupList(entity.GeneralQueryFilter) (data []entity.SourcesLookupResponse, total int64, lastPage int, err error)
	MeasurementLookupList(entity.MeasurementQueryFilter) (data []entity.MeasurementsLookupResponse, total int64, lastPage int, err error)
	AnnouncementLookupList(entity.GeneralQueryFilter) (data []entity.AnnouncementsLookupResponse, total int64, lastPage int, err error)
	SubAnnouncementLookupList(entity.SubAnnouncementQueryFilter) (data []entity.SubAnnouncementsLookupResponse, total int64, lastPage int, err error)
}

func NewGamificationService(gamificationRepository repository.GamificationRepository, transaction repository.Dbtransaction) *gamificationServiceImpl {
	return &gamificationServiceImpl{
		GamificationRepository: gamificationRepository,
		Transaction:            transaction,
	}
}

type gamificationServiceImpl struct {
	GamificationRepository repository.GamificationRepository
	Transaction            repository.Dbtransaction
}

func (service *gamificationServiceImpl) List(dataFilter entity.GamificationQueryFilter) (data []entity.GamificationListResponse, total int64, lastPage int, err error) {
	rtns, total, lastPage, err := service.GamificationRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range rtns {
		var vResp entity.GamificationListResponse
		structs.Automapper(row, &vResp)

		if row.GamificationDate != nil {
			InvDate := row.GamificationDate.Format("2006-01-02")
			vResp.GamificationDate = &InvDate
		}

		if row.StartDate != nil {
			InvDate := row.StartDate.Format("2006-01-02")
			vResp.StartDate = &InvDate
		}

		if row.EndDate != nil {
			InvDate := row.EndDate.Format("2006-01-02")
			vResp.EndDate = &InvDate
		}

		if row.FinishedDate != nil {
			InvDate := row.FinishedDate.Format("2006-01-02")
			vResp.FinishedDate = &InvDate
		}

		gamificationStatusName := vResp.GenerateGamificationStatusName()
		vResp.GamificationStatusName = &gamificationStatusName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) GamificationStatusLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.GamificationStatusesLookupResponse, total int64, lastPage int, err error) {
	GamificationStatuses, total, lastPage, err := service.GamificationRepository.FindAllGamificationStatusesLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, gamificationStatus := range GamificationStatuses {
		var vResp entity.GamificationStatusesLookupResponse
		structs.Automapper(gamificationStatus, &vResp)

		gamificationStatusName := vResp.GenerateDataGamificationStatusName()
		vResp.GamificationStatusName = &gamificationStatusName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) EmployeeGroupLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.EmployeeGroupsLookupResponse, total int64, lastPage int, err error) {
	EmployeeGroups, total, lastPage, err := service.GamificationRepository.FindAllEmployeeGroupByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range EmployeeGroups {
		var vResp entity.EmployeeGroupsLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) SalesTeamLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.SalesTeamsLookupResponse, total int64, lastPage int, err error) {
	SalesTeams, total, lastPage, err := service.GamificationRepository.FindAllSalesTeamByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range SalesTeams {
		var vResp entity.SalesTeamsLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) OperationTypeLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.OperationTypesLookupResponse, total int64, lastPage int, err error) {
	OperationTypes, total, lastPage, err := service.GamificationRepository.FindAllOperationTypeByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range OperationTypes {
		var vResp entity.OperationTypesLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) WarehouseLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.WarehousesLookupResponse, total int64, lastPage int, err error) {
	Warehouses, total, lastPage, err := service.GamificationRepository.FindAllWarehouseByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Warehouses {
		var vResp entity.WarehousesLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) SalesmanLookupList(dataFilter entity.SalesmanQueryFilter) (data []entity.SalesmansLookupResponse, total int64, lastPage int, err error) {
	Salesmans, total, lastPage, err := service.GamificationRepository.FindAllSalesmanByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Salesmans {
		var vResp entity.SalesmansLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) ProductCategoryLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.ProductCategoriesLookupResponse, total int64, lastPage int, err error) {
	ProductCategories, total, lastPage, err := service.GamificationRepository.FindAllProductCategoryByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range ProductCategories {
		var vResp entity.ProductCategoriesLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) ProductLineLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.ProductLinesLookupResponse, total int64, lastPage int, err error) {
	ProductLines, total, lastPage, err := service.GamificationRepository.FindAllProductLineByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range ProductLines {
		var vResp entity.ProductLinesLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) BrandLookupList(dataFilter entity.BrandQueryFilter) (data []entity.BrandsLookupResponse, total int64, lastPage int, err error) {
	Brands, total, lastPage, err := service.GamificationRepository.FindAllBrandByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Brands {
		var vResp entity.BrandsLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) SubBrand1LookupList(dataFilter entity.SubBrand1QueryFilter) (data []entity.SubBrands1LookupResponse, total int64, lastPage int, err error) {
	SubBrands1, total, lastPage, err := service.GamificationRepository.FindAllSubBrand1ByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range SubBrands1 {
		var vResp entity.SubBrands1LookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) SubBrand2LookupList(dataFilter entity.GeneralQueryFilter) (data []entity.SubBrands2LookupResponse, total int64, lastPage int, err error) {
	SubBrands2, total, lastPage, err := service.GamificationRepository.FindAllSubBrand2ByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range SubBrands2 {
		var vResp entity.SubBrands2LookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) ProductLookupList(dataFilter entity.ProductQueryFilter) (data []entity.ProductsLookupResponse, total int64, lastPage int, err error) {
	Products, total, lastPage, err := service.GamificationRepository.FindAllProductByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Products {
		var vResp entity.ProductsLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) OutletLocationLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.OutletLocationsLookupResponse, total int64, lastPage int, err error) {
	OutletLocations, total, lastPage, err := service.GamificationRepository.FindAllOutletLocationByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range OutletLocations {
		var vResp entity.OutletLocationsLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) OutletGroupLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.OutletGroupsLookupResponse, total int64, lastPage int, err error) {
	OutletGroups, total, lastPage, err := service.GamificationRepository.FindAllOutletGroupByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range OutletGroups {
		var vResp entity.OutletGroupsLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) OutletClassLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.OutletClassesLookupResponse, total int64, lastPage int, err error) {
	OutletClasses, total, lastPage, err := service.GamificationRepository.FindAllOutletClassByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range OutletClasses {
		var vResp entity.OutletClassesLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) MarketLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.MarketsLookupResponse, total int64, lastPage int, err error) {
	Markets, total, lastPage, err := service.GamificationRepository.FindAllMarketByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Markets {
		var vResp entity.MarketsLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) DistrictLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.DistrictsLookupResponse, total int64, lastPage int, err error) {
	Districts, total, lastPage, err := service.GamificationRepository.FindAllDistrictByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Districts {
		var vResp entity.DistrictsLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) IndustryLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.IndustriesLookupResponse, total int64, lastPage int, err error) {
	Industries, total, lastPage, err := service.GamificationRepository.FindAllIndustryByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Industries {
		var vResp entity.IndustriesLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) OutletLookupList(dataFilter entity.OutletQueryFilter) (data []entity.OutletsLookupResponse, total int64, lastPage int, err error) {
	Outlets, total, lastPage, err := service.GamificationRepository.FindAllOutletByCustIdLookupMode(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Outlets {
		var vResp entity.OutletsLookupResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) SourceLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.SourcesLookupResponse, total int64, lastPage int, err error) {
	var source entity.SourcesLookupResponse

	sources := source.GetDataSource()
	var keys []int
	for key := range sources {
		keys = append(keys, int(key))
	}

	sort.Ints(keys)

	for _, key := range keys {
		source.SourceId = int64(key)
		source.SourceName = sources[int64(key)]
		data = append(data, source)
	}

	total = int64(len(data))
	lastPage = 1

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) MeasurementLookupList(dataFilter entity.MeasurementQueryFilter) (data []entity.MeasurementsLookupResponse, total int64, lastPage int, err error) {
	var measurement entity.MeasurementsLookupResponse

	measurements := measurement.GetDataMeasurement()
	var keys []int
	for key := range measurements {
		keys = append(keys, int(key))
	}

	sort.Ints(keys)

	var measurementKeys []int
	if dataFilter.SourceId == 1 {
		measurementKeys = append(measurementKeys, 1, 2, 3)
	} else {
		measurementKeys = append(measurementKeys, 1, 2, 4)
	}

	for _, key := range measurementKeys {
		measurement.MeasurementId = int64(key)
		measurement.MeasurementName = measurements[int64(key)]
		data = append(data, measurement)
	}

	total = int64(len(data))
	lastPage = 1

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) AnnouncementLookupList(dataFilter entity.GeneralQueryFilter) (data []entity.AnnouncementsLookupResponse, total int64, lastPage int, err error) {
	var announcement entity.AnnouncementsLookupResponse

	announcements := announcement.GetDataAnnouncement()
	var keys []int
	for key := range announcements {
		keys = append(keys, int(key))
	}

	sort.Ints(keys)

	for _, key := range keys {
		announcement.AnnouncementId = int64(key)
		announcement.AnnouncementName = announcements[int64(key)]
		data = append(data, announcement)
	}

	total = int64(len(data))
	lastPage = 1

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) SubAnnouncementLookupList(dataFilter entity.SubAnnouncementQueryFilter) (data []entity.SubAnnouncementsLookupResponse, total int64, lastPage int, err error) {
	var subAnnouncement entity.SubAnnouncementsLookupResponse

	subAnnouncements := subAnnouncement.GetDataSubAnnouncement()
	var keys []int
	for key := range subAnnouncements {
		keys = append(keys, int(key))
	}

	sort.Ints(keys)

	var subAnnouncementKeys []int
	if dataFilter.AnnouncemnetId == 1 {
		subAnnouncementKeys = append(subAnnouncementKeys, 1, 2)
	} else {
		subAnnouncementKeys = append(subAnnouncementKeys, 3, 4, 5, 6)
	}

	for _, key := range subAnnouncementKeys {
		subAnnouncement.SubAnnouncementId = int64(key)
		subAnnouncement.SubAnnouncementName = subAnnouncements[int64(key)]
		data = append(data, subAnnouncement)
	}

	total = int64(len(data))
	lastPage = 1

	return data, total, lastPage, err
}

func (service *gamificationServiceImpl) Store(request entity.CreateGamificationBody) (err error) {
	c := context.Background()

	gamificationParticipants := request.Participants
	gamificationProducts := request.Products
	gamificationOutlets := request.Outlets
	gamificationAnnouncements := request.Announcements
	if err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		gamificationDate, err := str.DateStrToRfc3339String(time.Now().Format("2006-01-02"))
		if err != nil {
			return err
		}
		request.GamificationDate = gamificationDate

		startDate, err := str.DateStrToRfc3339String(request.StartDate)
		if err != nil {
			return err
		}
		request.StartDate = startDate

		endDate, err := str.DateStrToRfc3339String(request.EndDate)
		if err != nil {
			return err
		}
		request.EndDate = endDate

		var GamificationModel model.Gamification
		if err = structs.Automapper(request, &GamificationModel); err != nil {
			return err
		}
		GamificationModel.GamificationStatus = 1
		if err := service.GamificationRepository.Store(txCtx, &GamificationModel); err != nil {
			return err
		}

		for _, EmpID := range gamificationParticipants {
			var GamificationParticipantModel model.GamificationParticipant

			GamificationParticipantModel.CustID = request.CustID
			GamificationParticipantModel.GamificationNo = GamificationModel.GamificationNo
			GamificationParticipantModel.EmpID = EmpID
			GamificationParticipantModel.CreatedAt = time.Now()
			GamificationParticipantModel.CreatedBy = request.CreatedBy

			if err = service.GamificationRepository.StoreParticipant(txCtx, &GamificationParticipantModel); err != nil {
				return err
			}
		}

		for _, ProductID := range gamificationProducts {
			var GamificationProductModel model.GamificationProduct

			GamificationProductModel.CustID = request.CustID
			GamificationProductModel.GamificationNo = GamificationModel.GamificationNo
			GamificationProductModel.ProID = ProductID
			GamificationProductModel.CreatedAt = time.Now()
			GamificationProductModel.CreatedBy = request.CreatedBy

			if err = service.GamificationRepository.StoreProduct(txCtx, &GamificationProductModel); err != nil {
				return err
			}
		}

		for _, OutletID := range gamificationOutlets {
			var GamificationOutletModel model.GamificationOutlet

			GamificationOutletModel.CustID = request.CustID
			GamificationOutletModel.GamificationNo = GamificationModel.GamificationNo
			GamificationOutletModel.OutletID = OutletID
			GamificationOutletModel.CreatedAt = time.Now()
			GamificationOutletModel.CreatedBy = request.CreatedBy

			if err = service.GamificationRepository.StoreOutlet(txCtx, &GamificationOutletModel); err != nil {
				return err
			}
		}

		for _, AnnouncementDate := range gamificationAnnouncements {
			var GamificationAnnouncementModel model.GamificationAnnouncement

			AnnouncementDate, err = str.DateStrToRfc3339String(AnnouncementDate)
			if err != nil {
				return err
			}
			GamificationAnnouncementModel.CustID = request.CustID
			GamificationAnnouncementModel.GamificationNo = GamificationModel.GamificationNo

			announcementDate, err := time.Parse(time.RFC3339, AnnouncementDate)
			if err != nil {
				return err
			}
			GamificationAnnouncementModel.AnnouncementDate = announcementDate
			GamificationAnnouncementModel.CreatedAt = time.Now()
			GamificationAnnouncementModel.CreatedBy = request.CreatedBy

			if err = service.GamificationRepository.StoreAnnouncement(txCtx, &GamificationAnnouncementModel); err != nil {
				return err
			}
		}

		var GamificationFilterModel model.GamificationFilter
		GamificationFilterModel.GamificationNo = GamificationModel.GamificationNo

		salesTeams, err := json.Marshal(request.SalesTeams)
		if err != nil {
			return err
		}
		GamificationFilterModel.SalesTeam = string(salesTeams)

		operationTypes, err := json.Marshal(request.OperationTypes)
		if err != nil {
			return err
		}
		GamificationFilterModel.OperationType = string(operationTypes)

		warehouses, err := json.Marshal(request.Warehouses)
		if err != nil {
			return err
		}
		GamificationFilterModel.Warehouse = string(warehouses)

		productCategories, err := json.Marshal(request.ProductCategories)
		if err != nil {
			return err
		}
		GamificationFilterModel.ProductCategory = string(productCategories)

		productLines, err := json.Marshal(request.ProductLines)
		if err != nil {
			return err
		}
		GamificationFilterModel.ProductLine = string(productLines)

		brands, err := json.Marshal(request.Brands)
		if err != nil {
			return err
		}
		GamificationFilterModel.Brand = string(brands)

		subBrands1, err := json.Marshal(request.SubBrands1)
		if err != nil {
			return err
		}
		GamificationFilterModel.SubBrands1 = string(subBrands1)

		subBrands2, err := json.Marshal(request.SubBrands2)
		if err != nil {
			return err
		}
		GamificationFilterModel.SubBrands2 = string(subBrands2)

		outletLocations, err := json.Marshal(request.OutletLocations)
		if err != nil {
			return err
		}
		GamificationFilterModel.OutletLocation = string(outletLocations)

		outletGroups, err := json.Marshal(request.OutletGroups)
		if err != nil {
			return err
		}
		GamificationFilterModel.OutletGroup = string(outletGroups)

		outletClasses, err := json.Marshal(request.OutletClasses)
		if err != nil {
			return err
		}
		GamificationFilterModel.OutletClass = string(outletClasses)

		markets, err := json.Marshal(request.Markets)
		if err != nil {
			return err
		}
		GamificationFilterModel.Market = string(markets)

		districts, err := json.Marshal(request.Districts)
		if err != nil {
			return err
		}
		GamificationFilterModel.District = string(districts)

		industries, err := json.Marshal(request.Industries)
		if err != nil {
			return err
		}
		GamificationFilterModel.Industry = string(industries)

		if err = service.GamificationRepository.StoreFilter(txCtx, &GamificationFilterModel); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (service *gamificationServiceImpl) Detail(gamificationNo string, custID string, parentCustID string) (response entity.GamificationResponse, err error) {
	data, err := service.GamificationRepository.FindOneByGamificationNo(gamificationNo, custID, parentCustID)
	if err != nil {
		return response, err
	}

	if err = structs.Automapper(data, &response); err != nil {
		return response, err
	}

	Participants, err := service.GamificationRepository.FindGamificationParticipants(gamificationNo, custID, parentCustID)
	if err != nil {
		return response, err
	}

	var participantsData []entity.GamificationParticipantsResponse
	for _, participant := range Participants {
		var participantData entity.GamificationParticipantsResponse
		if err = structs.Automapper(participant, &participantData); err != nil {
			return response, err
		}

		participantsData = append(participantsData, participantData)
	}

	Products, err := service.GamificationRepository.FindGamificationProducts(gamificationNo, custID, parentCustID)
	if err != nil {
		return response, err
	}

	var productsData []entity.GamificationProductsResponse
	for _, product := range Products {
		var productData entity.GamificationProductsResponse
		if err = structs.Automapper(product, &productData); err != nil {
			return response, err
		}

		productsData = append(productsData, productData)
	}

	Outlets, err := service.GamificationRepository.FindGamificationOutlets(gamificationNo, custID, parentCustID)
	if err != nil {
		return response, err
	}

	var outletsData []entity.GamificationOutletsResponse
	for _, outlet := range Outlets {
		var outletData entity.GamificationOutletsResponse
		if err = structs.Automapper(outlet, &outletData); err != nil {
			return response, err
		}

		outletsData = append(outletsData, outletData)
	}

	Announcements, err := service.GamificationRepository.FindGamificationAnnouncements(gamificationNo, custID, parentCustID)
	if err != nil {
		return response, err
	}

	var announcementsData []entity.GamificationAnnouncementsResponse
	for _, announcement := range Announcements {
		var announcementData entity.GamificationAnnouncementsResponse
		if err = structs.Automapper(announcement, &announcementData); err != nil {
			return response, err
		}

		announcementData.AnnouncementDate = announcement.AnnouncementDate.Format("2006-01-02")
		announcementsData = append(announcementsData, announcementData)
	}

	if data.GamificationStatus == 2 || data.GamificationStatus == 3 || data.GamificationStatus == 4 {
		var Rankings []model.GamificationRankingRead
		switch data.MeasurementID {
		case 1:
			Rankings, err = service.GamificationRepository.FindGamificationRankingsByRevenue(data, custID, parentCustID)
			if err != nil {
				return response, err
			}
		case 2:
			Rankings, err = service.GamificationRepository.FindGamificationRankingsByQuantity(data, custID, parentCustID)
			if err != nil {
				return response, err
			}
		default:
			Rankings, err = service.GamificationRepository.FindGamificationRankingsByTotal(data, custID, parentCustID)
			if err != nil {
				return response, err
			}
		}

		var rankingsData []entity.GamificationRankingsResponse
		for index, ranking := range Rankings {
			var rankingData entity.GamificationRankingsResponse
			if err = structs.Automapper(ranking, &rankingData); err != nil {
				return response, err
			}
			rankingData.Ranking = int64(index) + 1

			rankingsData = append(rankingsData, rankingData)
		}

		response.Rankings = rankingsData
	}

	response.GamificationDate = data.GamificationDate.Format("2006-01-02")
	response.StartDate = data.StartDate.Format("2006-01-02")
	response.EndDate = data.EndDate.Format("2006-01-02")

	if data.FinishedDate != nil {
		finishedDate := data.FinishedDate.Format("2006-01-02")
		response.FinishedDate = &finishedDate
	}

	response.GamificationStatusName = response.GenerateGamificationStatusName()
	response.SourceName = response.GenerateGamificationSourceName()
	response.MeasurementName = response.GenerateGamificationMeasurementName()
	response.AnnouncementName = response.GenerateGamificationAnnouncementName()

	subAnnouncementName := response.GenerateGamificationSubAnnouncementName()
	response.SubAnnouncementName = &subAnnouncementName

	response.Participants = participantsData
	response.Products = productsData
	response.Outlets = outletsData
	response.Announcements = announcementsData

	gamificationFilter, err := service.GamificationRepository.FindGamificationFilter(gamificationNo, custID, parentCustID)
	if err != nil {
		return response, err
	}

	var salesTeams []int64
	if err := json.Unmarshal([]byte(gamificationFilter.SalesTeam), &salesTeams); err != nil {
		return response, err
	}
	response.SalesTeams = salesTeams

	var operationTypes []string
	if err := json.Unmarshal([]byte(gamificationFilter.OperationType), &operationTypes); err != nil {
		return response, err
	}
	response.OperationTypes = operationTypes

	var warehouses []int64
	if err := json.Unmarshal([]byte(gamificationFilter.Warehouse), &warehouses); err != nil {
		return response, err
	}
	response.Warehouses = warehouses

	var productCategories []int64
	if err := json.Unmarshal([]byte(gamificationFilter.ProductCategory), &productCategories); err != nil {
		return response, err
	}
	response.ProductCategories = productCategories

	var productLines []int64
	if err := json.Unmarshal([]byte(gamificationFilter.ProductLine), &productLines); err != nil {
		return response, err
	}
	response.ProductLines = productLines

	var brands []int64
	if err := json.Unmarshal([]byte(gamificationFilter.Brand), &brands); err != nil {
		return response, err
	}
	response.Brands = brands

	var subBrands1 []int64
	if err := json.Unmarshal([]byte(gamificationFilter.SubBrand1), &subBrands1); err != nil {
		return response, err
	}
	response.SubBrands1 = subBrands1

	var subBrands2 []int64
	if err := json.Unmarshal([]byte(gamificationFilter.SubBrand2), &subBrands2); err != nil {
		return response, err
	}
	response.SubBrands2 = subBrands2

	var outletLocations []int64
	if err := json.Unmarshal([]byte(gamificationFilter.OutletLocation), &outletLocations); err != nil {
		return response, err
	}
	response.OutletLocations = outletLocations

	var outletGroups []int64
	if err := json.Unmarshal([]byte(gamificationFilter.OutletGroup), &outletGroups); err != nil {
		return response, err
	}
	response.OutletGroups = outletGroups

	var outletClasses []int64
	if err := json.Unmarshal([]byte(gamificationFilter.OutletClass), &outletClasses); err != nil {
		return response, err
	}
	response.OutletClasses = outletClasses

	var markets []int64
	if err := json.Unmarshal([]byte(gamificationFilter.Market), &markets); err != nil {
		return response, err
	}
	response.Markets = markets

	var districts []int64
	if err := json.Unmarshal([]byte(gamificationFilter.District), &districts); err != nil {
		return response, err
	}
	response.Districts = districts

	var industries []int64
	if err := json.Unmarshal([]byte(gamificationFilter.Industry), &industries); err != nil {
		return response, err
	}
	response.Industries = industries

	return response, nil
}

func (service *gamificationServiceImpl) Update(gamificationNo string, request entity.UpdateGamificationBody) (err error) {
	c := context.Background()

	gamificationParticipants := request.Participants
	gamificationProducts := request.Products
	gamificationOutlets := request.Outlets
	gamificationAnnouncements := request.Announcements
	if err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		startDate, err := str.DateStrToRfc3339String(request.StartDate)
		if err != nil {
			return err
		}
		request.StartDate = startDate

		endDate, err := str.DateStrToRfc3339String(request.EndDate)
		if err != nil {
			return err
		}
		request.EndDate = endDate

		var GamificationModel model.Gamification
		if err = structs.Automapper(request, &GamificationModel); err != nil {
			return err
		}
		GamificationModel.CustID = ""
		GamificationModel.GamificationNo = gamificationNo
		GamificationModel.GamificationStatus = 1
		if err := service.GamificationRepository.Update(txCtx, &GamificationModel); err != nil {
			return err
		}

		if err := service.GamificationRepository.DeleteParticipant(txCtx, gamificationNo); err != nil {
			return err
		}
		for _, EmpID := range gamificationParticipants {
			var GamificationParticipantModel model.GamificationParticipant

			GamificationParticipantModel.CustID = request.CustID
			GamificationParticipantModel.GamificationNo = gamificationNo
			GamificationParticipantModel.EmpID = EmpID
			GamificationParticipantModel.CreatedAt = time.Now()
			GamificationParticipantModel.CreatedBy = request.UpdatedBy

			if err = service.GamificationRepository.StoreParticipant(txCtx, &GamificationParticipantModel); err != nil {
				return err
			}
		}

		if err := service.GamificationRepository.DeleteProduct(txCtx, gamificationNo); err != nil {
			return err
		}
		for _, ProductID := range gamificationProducts {
			var GamificationProductModel model.GamificationProduct

			GamificationProductModel.CustID = request.CustID
			GamificationProductModel.GamificationNo = gamificationNo
			GamificationProductModel.ProID = ProductID
			GamificationProductModel.CreatedAt = time.Now()
			GamificationProductModel.CreatedBy = request.UpdatedBy

			if err = service.GamificationRepository.StoreProduct(txCtx, &GamificationProductModel); err != nil {
				return err
			}
		}

		if err := service.GamificationRepository.DeleteOutlet(txCtx, gamificationNo); err != nil {
			return err
		}
		for _, OutletID := range gamificationOutlets {
			var GamificationOutletModel model.GamificationOutlet

			GamificationOutletModel.CustID = request.CustID
			GamificationOutletModel.GamificationNo = gamificationNo
			GamificationOutletModel.OutletID = OutletID
			GamificationOutletModel.CreatedAt = time.Now()
			GamificationOutletModel.CreatedBy = request.UpdatedBy

			if err = service.GamificationRepository.StoreOutlet(txCtx, &GamificationOutletModel); err != nil {
				return err
			}
		}

		if err := service.GamificationRepository.DeleteAnnouncement(txCtx, gamificationNo); err != nil {
			return err
		}
		for _, AnnouncementDate := range gamificationAnnouncements {
			var GamificationAnnouncementModel model.GamificationAnnouncement

			AnnouncementDate, err = str.DateStrToRfc3339String(AnnouncementDate)
			if err != nil {
				return err
			}
			GamificationAnnouncementModel.CustID = request.CustID
			GamificationAnnouncementModel.GamificationNo = gamificationNo

			announcementDate, err := time.Parse(time.RFC3339, AnnouncementDate)
			if err != nil {
				return err
			}
			GamificationAnnouncementModel.AnnouncementDate = announcementDate
			GamificationAnnouncementModel.CreatedAt = time.Now()
			GamificationAnnouncementModel.CreatedBy = request.UpdatedBy

			if err = service.GamificationRepository.StoreAnnouncement(txCtx, &GamificationAnnouncementModel); err != nil {
				return err
			}
		}

		if err := service.GamificationRepository.DeleteFilter(txCtx, gamificationNo); err != nil {
			return err
		}
		var GamificationFilterModel model.GamificationFilter
		GamificationFilterModel.GamificationNo = gamificationNo

		salesTeams, err := json.Marshal(request.SalesTeams)
		if err != nil {
			return err
		}
		GamificationFilterModel.SalesTeam = string(salesTeams)

		operationTypes, err := json.Marshal(request.OperationTypes)
		if err != nil {
			return err
		}
		GamificationFilterModel.OperationType = string(operationTypes)

		warehouses, err := json.Marshal(request.Warehouses)
		if err != nil {
			return err
		}
		GamificationFilterModel.Warehouse = string(warehouses)

		productCategories, err := json.Marshal(request.ProductCategories)
		if err != nil {
			return err
		}
		GamificationFilterModel.ProductCategory = string(productCategories)

		productLines, err := json.Marshal(request.ProductLines)
		if err != nil {
			return err
		}
		GamificationFilterModel.ProductLine = string(productLines)

		brands, err := json.Marshal(request.Brands)
		if err != nil {
			return err
		}
		GamificationFilterModel.Brand = string(brands)

		subBrands1, err := json.Marshal(request.SubBrands1)
		if err != nil {
			return err
		}
		GamificationFilterModel.SubBrands1 = string(subBrands1)

		subBrands2, err := json.Marshal(request.SubBrands2)
		if err != nil {
			return err
		}
		GamificationFilterModel.SubBrands2 = string(subBrands2)

		outletLocations, err := json.Marshal(request.OutletLocations)
		if err != nil {
			return err
		}
		GamificationFilterModel.OutletLocation = string(outletLocations)

		outletGroups, err := json.Marshal(request.OutletGroups)
		if err != nil {
			return err
		}
		GamificationFilterModel.OutletGroup = string(outletGroups)

		outletClasses, err := json.Marshal(request.OutletClasses)
		if err != nil {
			return err
		}
		GamificationFilterModel.OutletClass = string(outletClasses)

		markets, err := json.Marshal(request.Markets)
		if err != nil {
			return err
		}
		GamificationFilterModel.Market = string(markets)

		districts, err := json.Marshal(request.Districts)
		if err != nil {
			return err
		}
		GamificationFilterModel.District = string(districts)

		industries, err := json.Marshal(request.Industries)
		if err != nil {
			return err
		}
		GamificationFilterModel.Industry = string(industries)

		if err = service.GamificationRepository.StoreFilter(txCtx, &GamificationFilterModel); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (service *gamificationServiceImpl) Cancel(gamificationNo string, request entity.CancelGamificationBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		var Model model.Gamification
		err = structs.Automapper(request, &Model)
		if err != nil {
			return err
		}

		Model.CustID = ""
		Model.GamificationNo = gamificationNo

		// var reviewedBy = request.UpdatedBy
		// var reviewedAt = time.Now()

		Model.GamificationStatus = 5
		// Model.IsReviewed = true
		// Model.ReviewedBy = &reviewedBy
		// Model.ReviewedAt = &reviewedAt
		err = service.GamificationRepository.Update(txCtx, &Model)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *gamificationServiceImpl) Stop(gamificationNo string, request entity.StopGamificationBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		var Model model.Gamification
		err = structs.Automapper(request, &Model)
		if err != nil {
			return err
		}

		Model.CustID = ""
		Model.GamificationNo = gamificationNo
		finishedDate := time.Now()
		Model.FinishedDate = &finishedDate

		Model.GamificationStatus = 4
		err = service.GamificationRepository.Update(txCtx, &Model)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
