package entities

import (
	"time"
)

type Job struct {
	ID              int
	start           *time.Time
	normalizedStart time.Time
	Interval        int
	Task            *Task
}

func NewJob(id int, start *time.Time, interval int, task *Task) *Job {
	return &Job{
		ID:              id,
		start:           start,
		normalizedStart: *start,
		Interval:        interval,
		Task:            task,
	}
}

func (j *Job) Start() time.Time {
	interval := time.Hour * time.Duration(j.Interval)
	now := time.Now().Add(time.Minute * 3)

	for j.normalizedStart.Before(now) {
		j.normalizedStart = j.normalizedStart.Add(interval)
	}

	return j.normalizedStart
}
