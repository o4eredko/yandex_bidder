package job

import (
	"strconv"
	"sync"

	"github.com/go-co-op/gocron"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	jobsClaim map[int]*gocron.Job

	repo struct {
		mu        sync.RWMutex
		scheduler *gocron.Scheduler
		jobs      jobsClaim
	}
)

func New(scheduler *gocron.Scheduler) usecase.JobRepo {
	return &repo{
		scheduler: scheduler,
		jobs:      make(jobsClaim),
	}
}

func (r *repo) Add(job *entities.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tag := []string{strconv.Itoa(job.ID)}
	start := job.Start()
	cronJob, err := r.scheduler.
		Every(uint64(job.Interval)).
		Minutes().
		SetTag(tag).
		StartAt(start).
		Do(job.Task.Func, job.Task.Args)
	if err != nil {
		return err
	}
	r.jobs[job.ID] = cronJob

	return nil
}

func (r *repo) Scheduled(jobID int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cronJob, scheduled := r.jobs[jobID]
	if !scheduled {
		scheduled = r.scheduler.Scheduled(cronJob)
	}

	return scheduled
}

func (r *repo) Remove(jobID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cronJob, ok := r.jobs[jobID]
	if !ok {
		return domain.ErrJobNotFound
	}

	r.scheduler.RemoveByReference(cronJob)
	delete(r.jobs, jobID)

	return nil
}
