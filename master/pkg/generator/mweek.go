package generator

import (
	"time"
)

type Mweek struct {
	WeekStart   time.Time
	WeekIdStart int
	MPeriod     []MPeriod
}
type MPeriod struct {
	PerID     int
	WeekCount int
}

type MWeekCalculate struct {
	PerID     int
	WeekID    int
	WeekStart time.Time
	WeekEnd   time.Time
	MworkDays []MWorkDay
}

type MWorkDay struct {
	WorkDate time.Time
	IsWork   bool
}

func NewMweek(start time.Time, WeekIDid int) *Mweek {
	mweek := Mweek{
		WeekStart:   start,
		WeekIdStart: WeekIDid,
	}
	return &mweek
}

func (mw *Mweek) AddMPeriod(perID, weekCount int) {
	mperiod := MPeriod{
		PerID:     perID,
		WeekCount: weekCount,
	}
	mw.MPeriod = append(mw.MPeriod, mperiod)
}

func (mw *Mweek) Calculate() []MWeekCalculate {
	var Result []MWeekCalculate
	var tempWeekStart, tempWeekEnd time.Time
	var tempWeekID int = mw.WeekIdStart
	tempWeekStart = mw.WeekStart
	for _, MPeriod := range mw.MPeriod {
		for i := 0; i < MPeriod.WeekCount; i++ {
			tempWeekEnd = tempWeekStart.Add(time.Hour * 24 * 6)
			mweekcalc := MWeekCalculate{
				PerID:     MPeriod.PerID,
				WeekID:    tempWeekID,
				WeekEnd:   tempWeekEnd,
				WeekStart: tempWeekStart,
			}

			diffhour := tempWeekEnd.Sub(tempWeekStart).Hours()
			diffDay := diffhour / 24
			for i := 0; i <= int(diffDay); i++ {
				workdate := tempWeekStart.Add(time.Hour * 24 * time.Duration(i))
				isMonday := func(t time.Time) bool {
					return t.Weekday().String() != "Sunday"
				}
				mworkday := MWorkDay{
					WorkDate: workdate,
					IsWork:   isMonday(workdate),
				}
				mweekcalc.MworkDays = append(mweekcalc.MworkDays, mworkday)
			}

			Result = append(Result, mweekcalc)
			tempWeekID++
			tempWeekStart = tempWeekEnd.Add(time.Hour * 24)
		}

	}

	return Result
}
