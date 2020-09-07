package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/config"
)

type (
	store struct {
		Scheduler *gocron.Scheduler
	}
)

func New(config *config.Scheduler) *store {
	loc, err := time.LoadLocation(config.TimeZone)
	if err != nil {
		panic(err)
	}
	scheduler := gocron.NewScheduler(loc)
	scheduler.StartAsync()

	return &store{
		Scheduler: scheduler,
	}
}

func (s *store) Shutdown() error {
	s.Scheduler.Clear()
	s.Scheduler.Stop()

	return nil
}
