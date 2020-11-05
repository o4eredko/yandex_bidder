package job

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
)

func (r *repo) schedule(j *entities.Job, s cron.Schedule) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.cron.Schedule(s, j)
	go r.cron.Entry(id).Job.Run()
	r.jobs[j.ID] = id
	delete(r.pending, j.ID)
}

func (r *repo) scheduleWithDelay(j *pendingJob, s cron.Schedule) {
	log.Info().Msgf("[%d] will be scheduled at: %v", j.ID, j.ScheduleAt())
	for {
		select {
		case <-j.ticker:
			log.Info().Msgf("[%d] scheduling at %v", j.ID, time.Now().UTC())
			r.schedule(j.Job, s)
			log.Info().Msgf("[%d] scheduled", j.ID)
			return
		case <-j.cancel:
			log.Info().Msgf("[%d] canceled", j.ID)
			return
		}
	}
}
