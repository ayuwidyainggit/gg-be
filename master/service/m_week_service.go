package service

import (
	"errors"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"strconv"
	"strings"
	"time"
)

type MWeekService interface {
	Detail(int, int, int, string, string) (entity.MWeekResponse, error)
	List(entity.MWeekQueryFilter, string, string) (data []entity.MWeekResponse, total int, lastPage int, err error)
	Store(entity.CreateMWeekBody) (entity.MWeekResponse, error)
	Update(int, int, int, entity.UpdateMWeekRequest) error
	Delete(string, int, int, int, int64, string) error
}

func NewMWeekService(mWeekRepository repository.MWeekRepository) *mWeekServiceImpl {
	return &mWeekServiceImpl{
		MWeekRepository: mWeekRepository,
	}
}

type mWeekServiceImpl struct {
	MWeekRepository repository.MWeekRepository
}

func (service *mWeekServiceImpl) Detail(PerYear int, PerId int, WeekId int, custId string, parentCustId string) (response entity.MWeekResponse, err error) {
	scopeCustId, generatedOnly, err := service.resolveReadScope(entity.MWeekQueryFilter{PerYear: strconv.Itoa(PerYear)}, custId, parentCustId)
	if err != nil {
		return response, err
	}

	mWeeks, err := service.MWeekRepository.FindOneByPerYearAndPerIdAndWeekIdAndCustId(PerYear, PerId, WeekId, scopeCustId, generatedOnly)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(mWeeks, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *mWeekServiceImpl) List(dataFilter entity.MWeekQueryFilter, custId string, parentCustId string) (data []entity.MWeekResponse, total int, lastPage int, err error) {
	dataFilter.ParentCustId = parentCustId
	scopeCustId, generatedOnly, err := service.resolveReadScope(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	mWeeks, total, lastPage, err := service.MWeekRepository.FindAllByCustId(dataFilter, scopeCustId, generatedOnly)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range mWeeks {
		var vResp entity.MWeekResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *mWeekServiceImpl) resolveReadScope(dataFilter entity.MWeekQueryFilter, custId string, parentCustId string) (string, bool, error) {
	ownerCustId := parentCustId
	if ownerCustId == "" {
		ownerCustId = custId
	}

	if custId != ownerCustId {
		hasLegacyRows, err := service.MWeekRepository.HasLegacyRows(custId, dataFilter.PerYear)
		if err != nil {
			return "", false, err
		}
		if hasLegacyRows {
			return custId, false, nil
		}
	}

	hasGeneratedRows, err := service.MWeekRepository.HasGeneratedCalendarRows(ownerCustId, dataFilter.PerYear)
	if err != nil {
		return "", false, err
	}
	if hasGeneratedRows {
		return ownerCustId, true, nil
	}

	return custId, false, nil
}

func (service *mWeekServiceImpl) Store(request entity.CreateMWeekBody) (response entity.MWeekResponse, err error) {
	// mWeeks, err := service.MWeekRepository.FindOneByPerYearAndPerIdAndWeekIdAndCustId(request.PerYear, request.PerId, request.WeekId, request.CustId)
	// if err == nil {
	// 	return response, errors.New("week_id: " + strconv.Itoa(mWeeks.WeekId) + " is already exists")
	// }

	var mWeekData model.MWeek
	if request.WeekStart != nil {
		if *request.WeekStart == "" {
			request.WeekStart = nil
		}
	}
	if request.WeekEnd != nil {
		if *request.WeekEnd == "" {
			request.WeekEnd = nil
		}
	}
	structs.Automapper(request, &mWeekData)
	perId, err := service.MWeekRepository.Store(mWeekData)
	if err != nil {
		return response, err
	}

	response.PerId = perId

	return response, err
}

func (service *mWeekServiceImpl) Update(perYear int, perId int, weekId int, request entity.UpdateMWeekRequest) (err error) {
	mWeeks, err := service.MWeekRepository.FindAllByGreaterThanWeekIdAndCustId(perYear, weekId, request.CustId)
	if err != nil {
		return err
	}

	if len(mWeeks) == 0 {
		return errors.New("data master week not found")
	}

	if mWeeks[0].IsClosed {
		return errors.New("the week has been closed")
	}

	log.Println("mWeeks[0].WeekStart:", *mWeeks[0].WeekStart)
	log.Println("request.WeekStart:", *request.WeekStart)
	weekStartStr := *mWeeks[0].WeekStart
	weekEndStr := *mWeeks[0].WeekEnd
	if weekStartStr[:10] != *request.WeekStart {
		return errors.New("week_start must be " + weekStartStr[:10])
	}

	weekStartDateReq, err := time.Parse("2006-01-02", *request.WeekStart)
	if err != nil {
		return err
	}

	weekEndDateReq, err := time.Parse("2006-01-02", *request.WeekEnd)
	if err != nil {
		log.Println("mWeekServiceImpl, Update, weekEndDateReq time.Parse, err:", err.Error())
		return err
	}

	if weekEndDateReq.Before(weekStartDateReq) {
		log.Println("mWeekServiceImpl, Update, weekEndDateReq.Before(weekStartDateReq), err:", err.Error())
		return errors.New("week_end must at least same with week_start")
	}

	weekEndDate, err := time.Parse("2006-01-02", weekEndStr[:10])
	if err != nil {
		log.Println("mWeekServiceImpl, Update, weekEndDate time.Parse, err:", err.Error())
		return err
	}

	weekEndSub := weekEndDateReq.Sub(weekEndDate)
	diffDay := int64(weekEndSub.Hours() / 24)
	log.Println("diffDay:", diffDay)

	distributors, err := service.MWeekRepository.FindDistributorByCustId(request.CustId)
	if err != nil {
		return errors.New("must have at least one distributor")
	}

	groupCustId := "'" + request.CustId + "',"
	for _, row := range distributors {
		groupCustId += "'" + row.CustId + "',"
	}
	groupCustId = strings.TrimSuffix(groupCustId, ",")
	log.Println("mWeekServiceImpl, groupCustId:", groupCustId)

	trx, err := service.MWeekRepository.TrxBegin()
	if err != nil {
		return err
	}

	for i, mWeek := range mWeeks {
		if i == 0 {

			// update m_week principal and distributor
			err = trx.UpdateTrx(perYear, perId, weekId, *request.WeekStart, *request.WeekEnd, groupCustId)
			if err != nil {
				trx.TrxRollback()
				return err
			}

			if diffDay == 0 {
				err = trx.DeleteMWorkDayByWeekIdExceptFirstDate(perYear, perId, weekId, *mWeek.WeekEnd, groupCustId)
				if err != nil {
					trx.TrxRollback()
					return err
				}
			} else if diffDay < 0 {
				isNeedUpdatePerId := false
				if len(mWeeks) > 1 {
					if mWeeks[1].PerId != mWeek.PerId {
						isNeedUpdatePerId = true
					}
				}
				err = trx.UpdateMWorkDayByWeekIdGreaterThanDate(perYear, perId, weekId, weekEndDateReq.Format("2006-01-02"), groupCustId, isNeedUpdatePerId)
				if err != nil {
					trx.TrxRollback()
					return err
				}
			} else if diffDay > 0 {
				err = trx.UpdateMWorkDayByWeekIdAndGreaterThanDate(perYear, perId, weekId, *mWeek.WeekEnd, *request.WeekEnd, groupCustId)
				if err != nil {
					trx.TrxRollback()
					return err
				}
			}
		} else {
			mWeekStartStr := *mWeeks[i].WeekStart
			mWeekEndStr := *mWeeks[i].WeekEnd
			mWeekStartTime, err := time.Parse("2006-01-02", mWeekStartStr[:10])
			if err != nil {
				trx.TrxRollback()
				return err
			}
			mWeekStartStrNew := mWeekStartTime.AddDate(0, 0, int(diffDay)).Format("2006-01-02")

			mWeekEndTime, err := time.Parse("2006-01-02", mWeekEndStr[:10])
			if err != nil {
				trx.TrxRollback()
				return err
			}
			mWeekEndStrNew := mWeekEndTime.AddDate(0, 0, int(diffDay)).Format("2006-01-02")

			err = trx.UpdateMWeekStartEndDate(perYear, mWeek.WeekId, mWeekStartStrNew, mWeekEndStrNew, groupCustId)
			if err != nil {
				trx.TrxRollback()
				return err
			}

			if diffDay < 0 {
				isNeedUpdatePerId := false
				if i < len(mWeeks)-1 {
					if mWeeks[i+1].PerId != mWeek.PerId {
						isNeedUpdatePerId = true
					}
				}
				err = trx.UpdateMWorkDayByWeekIdGreaterThanDate(mWeek.PerYear, mWeek.PerId, mWeek.WeekId, mWeekEndStrNew, groupCustId, isNeedUpdatePerId)
				if err != nil {
					trx.TrxRollback()
					return err
				}

				if i == (len(mWeeks) - 1) {
					err = trx.DeleteMWorkDayByGreaterThanDate(mWeek.PerYear, mWeekEndStrNew, groupCustId)
					if err != nil {
						trx.TrxRollback()
						return err
					}
				}
			} else if diffDay > 0 {
				isNeedUpdatePerId := false
				if mWeeks[i-1].PerId != mWeek.PerId {
					isNeedUpdatePerId = true
				}

				err = trx.UpdateMWorkDayByWeekIdLowerThanDate(mWeek.PerYear, mWeek.PerId, mWeek.WeekId, mWeekStartStrNew, groupCustId, isNeedUpdatePerId)
				if err != nil {
					trx.TrxRollback()
					return err
				}

				if i == (len(mWeeks) - 1) {
					for wd := 0; wd < int(diffDay); wd++ {
						workDateStr := mWeekEndTime.AddDate(0, 0, wd+1).Format("2006-01-02")
						mWorkDate := model.MWorkingDay{
							CustId:   request.CustId,
							PerYear:  mWeek.PerYear,
							PerId:    mWeek.PerId,
							WeekId:   mWeek.WeekId,
							WorkDate: &workDateStr,
						}
						err = trx.StoreMWorkDay(mWorkDate)
						if err != nil {
							trx.TrxRollback()
							return err
						}
					}
				}
			}

		}
	}

	trx.TrxCommit()

	return err
}

func (service *mWeekServiceImpl) Delete(custId string, PerYear int, PerId int, WeekId int, closedBy int64, closedByName string) (err error) {

	err = service.MWeekRepository.Delete(custId, PerYear, PerId, WeekId, closedBy, closedByName)
	if err != nil {
		return err
	}
	return err
}
