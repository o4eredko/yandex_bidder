package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/config"
)

type (
	store struct {
		Cron *cron.Cron
		Jobs map[int]cron.EntryID
	}
)

func New(config *config.Scheduler) *store {
	loc, err := time.LoadLocation(config.TimeZone)
	if err != nil {
		panic(err)
	}
	c := cron.New(cron.WithLocation(loc))
	c.Start()
	return &store{
		Cron: c,
	}
}

func (s *store) Shutdown() error {
	s.Cron.Stop()
	return nil
}
