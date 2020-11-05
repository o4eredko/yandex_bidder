package scheduler

import (
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	useCase struct {
		repo usecase.JobRepo
	}

	UseCase interface {
		ShowInfo() *domain.SchedulerOut
	}
)

func New(repo usecase.JobRepo) UseCase {
	return &useCase{
		repo: repo,
	}
}

func (u *useCase) ShowInfo() *domain.SchedulerOut {
	return u.repo.SchedulerInfo()
}
