package utils

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"scyllax-pjp/model"
	"scyllax-pjp/repository"
	"time"
)

func CheckWeek(ctx context.Context, routePopPermanentRepository repository.RoutePopPermanentRepository, routePopDailyRepository repository.RoutePopDailyRepository) {
	tn := time.Now().UTC()

	// Check if today is Saturday
	if tn.Weekday() == time.Saturday {
		_, week := tn.ISOWeek()
		fmt.Println("weeks", week)
		popPermanent, err := routePopPermanentRepository.FindByWeek(ctx, week)

		if err != nil {
			fmt.Println("week not found")
		}

		for _, pop := range popPermanent {
			routePopDailyPermanent := model.RoutePopDaily{
				RouteCode: pop.RouteCode,
				Week:      pop.Week + 1,
				Day:       pop.Day,
				Date:      pop.Date,
				PjpCode:   pop.PjpCode,
				PjpID:     pop.PjpID,
				Year:      pop.Year,
			}
			routePopDailyRepository.Insert(ctx, routePopDailyPermanent)
		}
	}
}

func CheckWeeklyJob(routePopPermanentRepository repository.RoutePopPermanentRepository, routePopDailyRepository repository.RoutePopDailyRepository) {
	local, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		local = time.UTC
	}

	s := gocron.NewScheduler(local)

	//testing
	//s.Every(10).Seconds().Do(func() {
	//	fmt.Println("Job is running in tes....")
	//	CheckWeek(context.Background(), routePopPermanentRepository, routePopDailyRepository)
	//})

	//production
	s.Every(1).Day().At("00:00").Do(func() {
		fmt.Println("Job is running....")
		CheckWeek(context.Background(), routePopPermanentRepository, routePopDailyRepository)
	})
	s.StartAsync()
}
