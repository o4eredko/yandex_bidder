package job

import (
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	jobsClaim map[int]cron.EntryID

	pendingJob struct {
		*entities.Job
		ticker <-chan time.Time
		cancel chan struct{}
	}
	pendingJobs map[int]*pendingJob

	repo struct {
		cron    *cron.Cron
		mu      sync.RWMutex
		jobs    jobsClaim
		pending pendingJobs
	}
)

func New(cron *cron.Cron) usecase.JobRepo {
	return &repo{
		cron:    cron,
		jobs:    make(jobsClaim),
		pending: make(pendingJobs),
	}
}

func (r *repo) SchedulerInfo() *domain.SchedulerOut {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entries := r.cron.Entries()
	s := &domain.SchedulerOut{
		Count:   len(entries) + len(r.pending),
		Running: make([]*domain.RunningJobOut, 0, len(entries)),
		Pending: make([]*domain.PendingJobOut, 0, len(r.pending)),
	}
	for _, e := range entries {
		job := e.Job.(*entities.Job)
		s.Running = append(s.Running, &domain.RunningJobOut{
			ID:       job.ID,
			NextRun:  e.Next,
			PrevRun:  e.Prev,
			Interval: job.Interval().String(),
		})
	}
	for _, pendingJob := range r.pending {
		s.Pending = append(s.Pending, &domain.PendingJobOut{
			ID:         pendingJob.ID,
			ScheduleAt: pendingJob.ScheduleAt(),
		})
	}
	return s
}

func (r *repo) Add(job *entities.Job) error {
	if r.Scheduled(job.ID) {
		return domain.ErrJobAlreadyScheduled
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	at := job.ScheduleAt()
	delay := at.Sub(now)
	schedule := cron.Every(job.Interval())
	pending := &pendingJob{
		Job:    job,
		cancel: make(chan struct{}),
		ticker: time.After(delay),
	}
	go r.scheduleWithDelay(pending, schedule)
	r.pending[job.ID] = pending
	return nil
}

func (r *repo) Update(job *entities.Job) error {
	if err := r.Remove(job.ID); err != nil {
		return err
	}
	return r.Add(job)
}

func (r *repo) Scheduled(jobID int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.pending[jobID]; ok {
		return true
	}
	_, ok := r.jobs[jobID]
	return ok
}

func (r *repo) Remove(jobID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if pending, ok := r.pending[jobID]; ok {
		pending.cancel <- struct{}{}
		delete(r.pending, jobID)
	} else if id, ok := r.jobs[jobID]; ok {
		r.cron.Remove(id)
		delete(r.jobs, jobID)
	} else {
		return domain.ErrJobNotFound
	}
	return nil
}
