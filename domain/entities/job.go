package entities

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Job struct {
	ID              int
	start           *time.Time
	normalizedStart time.Time
	interval        int
	f               func()
}

func NewJob(id int, start *time.Time, interval int) *Job {
	return &Job{
		ID:       id,
		start:    start,
		interval: interval,
	}
}

func (j *Job) WithFunc(f func()) *Job {
	j.f = f
	return j
}

func (j *Job) Run() {
	if j.f == nil {
		log.Warn().Msgf("running empty job: %v", j)
		return
	}
	j.f()
}

func (j *Job) Interval() time.Duration {
	return time.Duration(j.interval) * time.Hour
}

func (j *Job) ScheduleAt() time.Time {
	if j.normalizedStart.IsZero() {
		j.normalizedStart = *j.start
	}
	now := time.Now().UTC().Add(time.Minute)
	interval := j.Interval()
	for j.normalizedStart.Before(now) {
		j.normalizedStart = j.normalizedStart.Add(interval)
	}
	return j.normalizedStart
}
