package service

import (
	"errors"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/generator"
	"master/pkg/structs"
	"master/repository"
	"strconv"
	"strings"
	"time"
)

type MPeriodsService interface {
	Detail(int, int, string) (entity.MPeriodsResponse, error)
	List(entity.MPeriodsQueryFilter, string) (data []entity.MPeriodsResponse, total int, lastPage int, err error)
	Store(entity.CreateMPeriodsBody) (entity.MPeriodsResponse, error)
	Update(int, int, entity.UpdateMPeriodsRequest) error
	Delete(string, int, int, int64, string) error
	ListYear(entity.MPeriodsQueryFilter, string) (data []entity.MPeriodsListYear, err error)
}

func NewMPeriodsService(mPeriodsRepository repository.MPeriodsRepository) *mPeriodsServiceImpl {
	return &mPeriodsServiceImpl{
		MPeriodsRepository: mPeriodsRepository,
	}
}

type mPeriodsServiceImpl struct {
	MPeriodsRepository repository.MPeriodsRepository
}

func (service *mPeriodsServiceImpl) Detail(PerYear int, PerId int, custId string) (response entity.MPeriodsResponse, err error) {
	mPeriods, err := service.MPeriodsRepository.FindOneByPerYearAndPerIdAndCustId(PerYear, PerId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(mPeriods, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *mPeriodsServiceImpl) List(dataFilter entity.MPeriodsQueryFilter, custId string) (data []entity.MPeriodsResponse, total int, lastPage int, err error) {
	mPeriods, total, lastPage, err := service.MPeriodsRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range mPeriods {
		var vResp entity.MPeriodsResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *mPeriodsServiceImpl) Store(request entity.CreateMPeriodsBody) (response entity.MPeriodsResponse, err error) {
	log.Println("masuk service store")
	if request.CustId != request.ParentCustId {
		return response, errors.New("this endpoint only for principal")
	}
	hasGeneratedRows, err := service.MPeriodsRepository.HasGeneratedCalendarRows(request.CustId, request.PerYear)
	if err != nil {
		return response, err
	}
	if hasGeneratedRows {
		return response, errors.New("period mutation is not allowed because working day calendar data already exists for this year")
	}

	listConfig := "'start_date','max_periods','def_week','def_week_count','max_week_count'"
	mConfigs, err := service.MPeriodsRepository.FindConfigsByCustId(request.CustId, listConfig)
	if err != nil {
		return response, errors.New("configuration for create master period not found")
	}

	var startDate, maxPeriods, defWeek, defWeekCount, maxWeekCount string
	for _, row := range mConfigs {
		switch row.ConfigId {
		case "start_date":
			startDate = row.ConfigValue
		case "max_periods":
			maxPeriods = row.ConfigValue
		case "def_week":
			defWeek = row.ConfigValue
		case "def_week_count":
			defWeekCount = row.ConfigValue
		case "max_week_count":
			maxWeekCount = row.ConfigValue
		default:
		}
	}

	if startDate == "" || maxPeriods == "" || defWeek == "" || defWeekCount == "" || maxWeekCount == "" {
		return response, errors.New("missing configuration to create master period")
	}

	maxPeriodsInt, err := strconv.Atoi(maxPeriods)
	if err != nil {
		return response, err
	}

	if request.PerId > maxPeriodsInt {
		return response, errors.New("per_id must be " + maxPeriods + " or less")
	}

	_, err = service.MPeriodsRepository.FindOneByPerYearAndPerIdAndCustId(request.PerYear, request.PerId, request.CustId)
	if err == nil {
		errMsg := fmt.Sprintf("per_year: %d, per_id: %d already exists", request.PerYear, request.PerId)
		return response, errors.New(errMsg)
	}

	isFirst := true
	startNewWeek := 1
	mWeek, err := service.MPeriodsRepository.FindOneLastDataCustId(request.PerYear, request.CustId)
	if err == nil {
		lastWeek, err := time.Parse("2006-01-02T00:00:00Z", *mWeek.WeekEnd)
		if err != nil {
			log.Println("error disini")
			return response, err
		}
		log.Println("mWeek:", structs.StructToJson(mWeek))
		log.Println("lastWeek:", structs.StructToJson(lastWeek))
		nextDate := lastWeek.AddDate(0, 0, 1)
		startDate = nextDate.Format("2006-01-02")
		isFirst = false
		startNewWeek = mWeek.WeekId + 1
	} else {
		request.PerId = 1
		startDateFromReq := time.Date(request.PerYear, 1, 1, 0, 0, 0, 0, time.UTC)
		startDate = startDateFromReq.Format("2006-01-02")
	}

	log.Println("startDate:", startDate)

	layoutFormat := "2006-01-02"
	startDateTime, err := time.Parse(layoutFormat, startDate)
	if err != nil {
		return response, err
	}
	startYearStr := startDateTime.Format("2006")
	startYearInt, err := strconv.Atoi(startYearStr)
	if err != nil {
		return response, err
	}

	var lengthOfDefWeek = len([]rune(defWeek))
	firstDefWeekStr := defWeek[0:1]
	firstDefWeekInt, err := strconv.Atoi(firstDefWeekStr)
	if err != nil {
		return response, errors.New("failed convert default week config, error:" + err.Error())
	}
	if lengthOfDefWeek < 1 || firstDefWeekInt < 1 {
		return response, errors.New("config def_week must have at least one digit and greater than one")
	}

	mweekGenerator := generator.NewMweek(startDateTime, startNewWeek)

	distributors, err := service.MPeriodsRepository.FindDistributorByCustId(request.CustId)
	if err != nil {
		return response, errors.New("must have at least one distributor")
	}

	trx, err := service.MPeriodsRepository.TrxBegin()
	if err != nil {
		return response, err
	}

	defWeekCfgStr := defWeek
	sumOfWeeks := 0
	defWeekCountInt, err := strconv.Atoi(defWeekCount)
	if err != nil {
		return response, err
	}
	defMaxWeekCountInt, err := strconv.Atoi(maxWeekCount)
	if err != nil {
		return response, err
	}

	log.Println("startNewWeek:", startNewWeek)
	totalWeeks := startNewWeek + sumOfWeeks
	log.Println("maxPeriodsInt:", maxPeriodsInt)
	for perId := request.PerId; perId <= maxPeriodsInt; perId++ {
		log.Println("perId:", perId)
		log.Println("sumOfWeeks:", sumOfWeeks)
		log.Println("defWeekCountInt:", defWeekCountInt)
		log.Println("totalWeeks:", totalWeeks)
		log.Println("defMaxWeekCountInt:", defMaxWeekCountInt)

		if sumOfWeeks >= defWeekCountInt || totalWeeks > defMaxWeekCountInt {
			break
		}
		var mPeriod model.MPeriods
		timeNow := time.Now().In(time.UTC)
		structs.Automapper(request, &mPeriod)
		mPeriod.PerId = perId
		mPeriod.PerYear = startYearInt
		mPeriod.UpdatedBy = &request.UpdatedBy
		mPeriod.UpdatedAt = &timeNow
		defWeekCfgInt, err := strconv.Atoi(defWeekCfgStr[0:1])
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		defWeekCfgStr = defWeekCfgStr[1:]
		if defWeekCfgStr == "" {
			defWeekCfgStr = defWeek
		}
		mPeriod.WeekCount = &defWeekCfgInt
		isActivePeriod := false
		if isFirst && perId == 1 {
			isActivePeriod = true
			mPeriod.IsActive = &isActivePeriod
		}

		tempSumOfWeek := sumOfWeeks
		sumOfWeeks += defWeekCfgInt
		if sumOfWeeks > defWeekCountInt {
			defWeekCfgInt = defWeekCountInt - tempSumOfWeek
			mPeriod.WeekCount = &defWeekCfgInt
		}

		if !isFirst && request.WeekCount > 0 && perId == request.PerId {
			sumOfWeeks = mWeek.WeekId + request.WeekCount
			diffWeekCount := sumOfWeeks - defMaxWeekCountInt
			weekCount := request.WeekCount - diffWeekCount
			mPeriod.WeekCount = &weekCount
		}

		log.Println("before storePeriod, mPeriod:", structs.StructToJson(mPeriod))
		// store master period for principal
		err = trx.StorePeriod(mPeriod)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}

		// store master period for distributor
		for _, distributor := range distributors {
			var distPeriod model.MPeriods
			structs.Automapper(mPeriod, &distPeriod)
			distPeriod.CustId = distributor.CustId
			err = trx.StorePeriod(distPeriod)
			if err != nil {
				trx.TrxRollback()
				return response, err
			}
		}

		mweekGenerator.AddMPeriod(perId, *mPeriod.WeekCount)
	}
	results := mweekGenerator.Calculate()

	tempActiveWorkDay := false
	for i, result := range results {
		isActive := false
		if isFirst && i == 0 {
			isActive = true
		}
		weekStart := result.WeekStart.Format("2006-01-02")
		weekEnd := result.WeekEnd.Format("2006-01-02")
		mweekData := model.MWeek{
			CustId:    request.CustId,
			PerYear:   startYearInt,
			PerId:     result.PerID,
			WeekId:    result.WeekID,
			WeekStart: &weekStart,
			WeekEnd:   &weekEnd,
			IsActive:  &isActive,
		}
		// store master week for principal
		_, err := trx.StoreMweek(mweekData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}

		// store master week for distributor
		for _, distributor := range distributors {
			distWeek := model.MWeek{
				CustId:    distributor.CustId,
				PerYear:   startYearInt,
				PerId:     result.PerID,
				WeekId:    result.WeekID,
				WeekStart: &weekStart,
				WeekEnd:   &weekEnd,
				IsActive:  &isActive,
			}
			_, err := trx.StoreMweek(distWeek)
			if err != nil {
				trx.TrxRollback()
				return response, err
			}
		}

		for idxWd, mWorkDay := range result.MworkDays {
			isActiveWd := false
			if idxWd == 0 && !tempActiveWorkDay && isFirst {
				isActiveWd = true
				tempActiveWorkDay = true
			}
			workDate := mWorkDay.WorkDate.Format("2006-01-02")
			mworkdayModel := model.MWorkingDay{
				CustId:   request.CustId,
				PerYear:  startYearInt,
				PerId:    result.PerID,
				WeekId:   result.WeekID,
				WorkDate: &workDate,
				IsActive: &isActiveWd,
				IsWork:   &mWorkDay.IsWork,
			}

			// store master work day for principal
			_, err := trx.StoreMWorkingDay(mworkdayModel)
			if err != nil {
				trx.TrxRollback()
				return response, err
			}

			// store master work day for distributor
			for _, distributor := range distributors {
				distWorkDay := model.MWorkingDay{
					CustId:   distributor.CustId,
					PerYear:  startYearInt,
					PerId:    result.PerID,
					WeekId:   result.WeekID,
					WorkDate: &workDate,
					IsActive: &isActiveWd,
					IsWork:   &mWorkDay.IsWork,
				}
				_, err := trx.StoreMWorkingDay(distWorkDay)
				if err != nil {
					trx.TrxRollback()
					return response, err
				}
			}
		}
	}
	trx.TrxCommit()
	return response, err
}

func (service *mPeriodsServiceImpl) Update(perYear int, perId int, request entity.UpdateMPeriodsRequest) (err error) {
	if request.CustId != request.ParentCustId {
		return errors.New("this endpoint only for principal")
	}
	hasGeneratedRows, err := service.MPeriodsRepository.HasGeneratedCalendarRows(request.CustId, perYear)
	if err != nil {
		return err
	}
	if hasGeneratedRows {
		return errors.New("period mutation is not allowed because working day calendar data already exists for this year")
	}

	listConfig := "'max_week_count'"
	mConfigs, err := service.MPeriodsRepository.FindConfigsByCustId(request.CustId, listConfig)
	if err != nil {
		return errors.New("configuration for create master period not found")
	}

	var maxWeekCountStr string
	var maxWeekCountInt int
	for _, row := range mConfigs {
		switch row.ConfigId {
		case "max_week_count":
			maxWeekCountStr = row.ConfigValue
			maxWeekCountInt, err = strconv.Atoi(maxWeekCountStr)
			if err != nil {
				return err
			}
		default:
		}
	}

	distributors, err := service.MPeriodsRepository.FindDistributorByCustId(request.CustId)
	if err != nil {
		return errors.New("must have at least one distributor")
	}

	groupCustId := "'" + request.CustId + "',"
	for _, row := range distributors {
		groupCustId += "'" + row.CustId + "',"
	}
	groupCustId = strings.TrimSuffix(groupCustId, ",")
	// log.Println("mPeriodsServiceImpl, groupCustId:", groupCustId)

	trx, err := service.MPeriodsRepository.TrxBegin()
	if err != nil {
		return err
	}
	err = trx.Update(perYear, perId, request)
	if err != nil {
		trx.TrxRollback()
		return err
	}
	// update master period for distributor
	for _, distributor := range distributors {
		var distPeriodUpd entity.UpdateMPeriodsRequest
		structs.Automapper(request, &distPeriodUpd)
		distPeriodUpd.CustId = distributor.CustId
		err = trx.Update(perYear, perId, distPeriodUpd)
		if err != nil {
			trx.TrxRollback()
			return err
		}
	}

	sumTotalWeekCountExclReq, err := trx.FindTotalWeekCountExclByPerIdAndCustId(perYear, perId, request.CustId)
	if err != nil {
		trx.TrxRollback()
		return err
	}
	sumWeekCount := *sumTotalWeekCountExclReq.WeekCount + request.WeekCount

	if sumWeekCount > maxWeekCountInt {
		trx.TrxRollback()
		sumWeekCountStr := strconv.Itoa(sumWeekCount)
		errMsg := "total week count maximum is " + maxWeekCountStr + ", if including this request is " + sumWeekCountStr
		return errors.New(errMsg)
	}

	mPeriodsAll, err := trx.FindAllBymPeriodsCodeAndCustId(perYear, perId, request.CustId)
	if err != nil {
		trx.TrxRollback()
		return err
	}

	if mPeriodsAll[0].IsClosed {
		trx.TrxRollback()
		return errors.New("period is closed")
	}

	mWeeks, err := service.MPeriodsRepository.FindOneByWeekStartMinimum(perYear, perId, request.CustId)
	if err != nil {
		trx.TrxRollback()
		return err
	}
	err = trx.DeleteMWeekGreaterThanPeriod(perYear, perId, groupCustId)
	if err != nil {
		trx.TrxRollback()
		return err
	}
	err = trx.DeleteMWorkDayGreaterThanPeriod(perYear, perId, groupCustId)
	if err != nil {
		trx.TrxRollback()
		return err
	}
	//dateString := "2023-04-05T00:00:00Z"
	date, err := time.Parse(time.RFC3339, *mWeeks.WeekStart)
	if err != nil {
		trx.TrxRollback()
		return err
	}
	mweekGenerator := generator.NewMweek(date, mWeeks.WeekId)

	for _, mPeriod := range mPeriodsAll {
		mweekGenerator.AddMPeriod(mPeriod.PerId, *mPeriod.WeekCount)
	}
	results := mweekGenerator.Calculate()

	for _, result := range results {
		isActive := false
		weekStart := result.WeekStart.Format("2006-01-02")
		weekEnd := result.WeekEnd.Format("2006-01-02")
		mweekData := model.MWeek{
			CustId:    request.CustId,
			PerYear:   perYear,
			PerId:     result.PerID,
			WeekId:    result.WeekID,
			WeekStart: &weekStart,
			WeekEnd:   &weekEnd,
			IsActive:  &isActive,
		}
		_, err := trx.StoreMweek(mweekData)
		if err != nil {
			trx.TrxRollback()
			return err
		}
		// store master week for distributor
		for _, distributor := range distributors {
			distWeek := model.MWeek{
				CustId:    distributor.CustId,
				PerYear:   perYear,
				PerId:     result.PerID,
				WeekId:    result.WeekID,
				WeekStart: &weekStart,
				WeekEnd:   &weekEnd,
				IsActive:  &isActive,
			}
			_, err := trx.StoreMweek(distWeek)
			if err != nil {
				trx.TrxRollback()
				return err
			}
		}

		for _, mWorkDay := range result.MworkDays {
			workDate := mWorkDay.WorkDate.Format("2006-01-02")
			mworkdayModel := model.MWorkingDay{
				CustId:   request.CustId,
				PerYear:  perYear,
				PerId:    result.PerID,
				WeekId:   result.WeekID,
				WorkDate: &workDate,
				IsActive: &isActive,
				IsWork:   &mWorkDay.IsWork,
			}
			_, err := trx.StoreMWorkingDay(mworkdayModel)
			if err != nil {
				trx.TrxRollback()
				return err
			}
			// store master work day for distributor
			for _, distributor := range distributors {
				distWorkDay := model.MWorkingDay{
					CustId:   distributor.CustId,
					PerYear:  perYear,
					PerId:    result.PerID,
					WeekId:   result.WeekID,
					WorkDate: &workDate,
					IsActive: &isActive,
					IsWork:   &mWorkDay.IsWork,
				}
				_, err := trx.StoreMWorkingDay(distWorkDay)
				if err != nil {
					trx.TrxRollback()
					return err
				}
			}
		}
	}
	trx.TrxCommit()

	return err
}

func (service *mPeriodsServiceImpl) Delete(custId string, PerYear int, PerId int, closedBy int64, closedByName string) (err error) {
	hasGeneratedRows, err := service.MPeriodsRepository.HasGeneratedCalendarRows(custId, PerYear)
	if err != nil {
		return err
	}
	if hasGeneratedRows {
		return errors.New("period mutation is not allowed because working day calendar data already exists for this year")
	}

	err = service.MPeriodsRepository.Delete(custId, PerYear, PerId, closedBy, closedByName)
	if err != nil {
		return err
	}
	return err
}

func (service *mPeriodsServiceImpl) ListYear(dataFilter entity.MPeriodsQueryFilter, custId string) (data []entity.MPeriodsListYear, err error) {
	mPeriods, err := service.MPeriodsRepository.FindAllYearByCustId(dataFilter, custId)
	if err != nil {
		return data, err
	}

	for _, row := range mPeriods {
		var vResp entity.MPeriodsListYear
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, err
}
