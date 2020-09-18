package group

import (
	"github.com/rs/zerolog/log"
	"math"
	"time"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	useCase struct {
		concurrencyLimit int
		groupRepo        usecase.GroupRepo
		accountRepo      usecase.AccountRepo
		bidRepo          usecase.BidRepo
		jobRepo          usecase.JobRepo
	}

	UseCase interface {
		GetAll() ([]*entities.Group, error)
		Update(updateIn *domain.GroupUpdateIn) (*entities.Group, error)
		FixBids(groupID int) error
		Schedule(group *entities.Group) error
		ScheduleAll(stopOnErr bool) error
		Unschedule(group *entities.Group) error
		Start(groupID int) error
		Terminate(groupID int) error
	}
)

func New(
	concurrencyLimit int,
	groupRepo usecase.GroupRepo,
	accountRepo usecase.AccountRepo,
	bidRepo usecase.BidRepo,
	jobRepo usecase.JobRepo,
) UseCase {
	return &useCase{
		concurrencyLimit: concurrencyLimit,
		groupRepo:        groupRepo,
		accountRepo:      accountRepo,
		bidRepo:          bidRepo,
		jobRepo:          jobRepo,
	}
}

func (u *useCase) GetAll() ([]*entities.Group, error) {
	groups, err := u.groupRepo.GetAll()
	if err != nil {
		return nil, err
	}
	for _, group := range groups {
		group.IsScheduled = u.jobRepo.Scheduled(group.ID)
	}
	return groups, nil
}

func (u *useCase) Update(updateIn *domain.GroupUpdateIn) (*entities.Group, error) {
	group, err := u.groupRepo.GetByID(updateIn.ID)
	if err != nil {
		return nil, err
	}
	group.Start = updateIn.Start
	group.Interval = &updateIn.Interval
	group.Strategy = &updateIn.Strategy
	if err := u.groupRepo.Update(group); err != nil {
		return nil, err
	}
	if u.jobRepo.Scheduled(group.ID) {
		task := entities.NewTask(u.FixBids, group.ID)
		job := entities.NewJob(group.ID, group.Start, *group.Interval, task)
		if err := u.jobRepo.Update(job); err != nil {
			return nil, err
		}
	}
	return group, nil
}

func (u *useCase) FixBids(groupID int) error {
	now := time.Now().UTC()
	log.Info().Msgf("start FixBids for groupID=%d at %v", groupID, now)
	group, err := u.groupRepo.GetByID(groupID)
	if err != nil {
		return err
	}
	accounts, err := u.groupRepo.Accounts(group)
	if err != nil {
		return err
	}
	accountsWithBids := make([]*domain.AccountBids, 0, len(accounts))
	numOfWorkers := int(math.Min(float64(len(accounts)), float64(u.concurrencyLimit)))
	if numOfWorkers < 2 {
		accountsWithBids, err = u.calculateBids(accounts, *group.Strategy)
	} else {
		accountsWithBids, err = u.calculateWithWorkers(numOfWorkers, accounts, *group.Strategy)
	}
	if err != nil {
		return err
	}
	groupToUpdate := &domain.GroupToUpdateBids{
		Name:       group.Name,
		Accounts:   accountsWithBids,
		MaxRetries: group.MaxRetries(),
	}
	return u.bidRepo.Update(groupToUpdate)
}

func (u *useCase) Schedule(group *entities.Group) error {
	if scheduled := u.jobRepo.Scheduled(group.ID); scheduled {
		return domain.ErrJobAlreadyScheduled
	}
	if err := group.ValidateSchedule(); err != nil {
		return err
	}
	task := entities.NewTask(u.FixBids, group.ID)
	job := entities.NewJob(group.ID, group.Start, *group.Interval, task)
	group.ShouldSchedule = true
	if err := u.groupRepo.Update(group); err != nil {
		return err
	}
	return u.jobRepo.Add(job)
}

func (u *useCase) ScheduleAll(suppressErr bool) error {
	groups, err := u.GetAll()
	if err != nil {
		return err
	}
	for _, group := range groups {
		if group.ShouldSchedule {
			err := u.Schedule(group)
			if err != nil && !suppressErr {
				return err
			}
		}
	}
	return nil
}

func (u *useCase) Unschedule(group *entities.Group) error {
	if scheduled := u.jobRepo.Scheduled(group.ID); !scheduled {
		return domain.ErrJobNotFound
	}
	if err := u.jobRepo.Remove(group.ID); err != nil {
		return err
	}
	group.ShouldSchedule = false
	return u.groupRepo.Update(group)
}

func (u *useCase) Start(groupID int) error {
	group, err := u.groupRepo.GetByID(groupID)
	if err != nil {
		return err
	}
	return u.Schedule(group)
}

func (u *useCase) Terminate(groupID int) error {
	group, err := u.groupRepo.GetByID(groupID)
	if err != nil {
		return err
	}
	return u.Unschedule(group)
}
