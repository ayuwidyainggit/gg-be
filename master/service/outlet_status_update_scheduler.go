package service

import (
	"log"
	"strings"
	"time"
)

const (
	defaultOutletStatusUpdateTime = "01:00:00"
	outletStatusScheduleLocation  = "Asia/Jakarta"
)

func (service *outletServiceImpl) StartStatusUpdateScheduler(runAt string) {
	scheduleStr := strings.TrimSpace(runAt)
	if scheduleStr == "" {
		scheduleStr = defaultOutletStatusUpdateTime
	}
	scheduledTOD, err := time.Parse("15:04:05", scheduleStr)
	if err != nil {
		log.Printf("outlet status scheduler: invalid OUTLET_STATUS_UPDATE_TIME=%q, fallback to %s", scheduleStr, defaultOutletStatusUpdateTime)
		scheduledTOD, _ = time.Parse("15:04:05", defaultOutletStatusUpdateTime)
		scheduleStr = defaultOutletStatusUpdateTime
	}

	loc, err := time.LoadLocation(outletStatusScheduleLocation)
	if err != nil {
		log.Printf("outlet status scheduler: cannot load %s (%v), using local timezone", outletStatusScheduleLocation, err)
		loc = time.Local
	}

	go func() {
		for {
			now := time.Now().In(loc)
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), scheduledTOD.Hour(), scheduledTOD.Minute(), scheduledTOD.Second(), 0, loc)
			if !nextRun.After(time.Now()) {
				nextRun = nextRun.Add(24 * time.Hour)
			}

			wait := time.Until(nextRun)
			time.Sleep(wait)

			log.Printf("outlet status scheduler: start job at %s", time.Now().In(loc).Format(time.RFC3339))
			rows, err := service.UpdateStatuses()
			if err != nil {
				log.Printf("outlet status scheduler: error: %v", err)
				continue
			}
			log.Printf("outlet status scheduler: success, rows affected=%d", rows)
		}
	}()
}
