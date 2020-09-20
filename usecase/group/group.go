package group

import (
	"math"
	"time"

	"github.com/rs/zerolog/log"

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
		ScheduleAll(stopOnErr bool) error
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
	if u.jobRepo.Scheduled(group.ID) {
		job := u.newJob(group)
		if err := u.jobRepo.Update(job); err != nil {
			return nil, err
		}
		group.IsScheduled = true
		group.ShouldSchedule = true
	}
	if err := u.groupRepo.Update(group); err != nil {
		return nil, err
	}
	return group, nil
}

func (u *useCase) FixBids(groupID int) error {
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
	if len(accountsWithBids) == 0 {
		log.Info().Msgf("[%d] bids for update not found, skipping", group.ID)
		return nil
	}
	groupToUpdate := &domain.GroupToUpdateBids{
		Name:       group.Name,
		Accounts:   accountsWithBids,
		MaxRetries: group.MaxRetries(),
	}
	return u.bidRepo.Update(groupToUpdate)
}

func (u *useCase) ScheduleAll(suppressErr bool) error {
	groups, err := u.GetAll()
	if err != nil {
		return err
	}
	for _, group := range groups {
		if group.ShouldSchedule {
			err := u.schedule(group)
			if err != nil && !suppressErr {
				return err
			}
		}
	}
	return nil
}

func (u *useCase) Start(groupID int) error {
	group, err := u.groupRepo.GetByID(groupID)
	if err != nil {
		return err
	}
	if err := u.schedule(group); err != nil {
		return err
	}
	group.IsScheduled = true
	group.ShouldSchedule = true
	return u.groupRepo.Update(group)
}

func (u *useCase) Terminate(groupID int) error {
	group, err := u.groupRepo.GetByID(groupID)
	if err != nil {
		return err
	}
	if err := u.jobRepo.Remove(group.ID); err != nil {
		return err
	}
	group.ShouldSchedule = false
	return u.groupRepo.Update(group)
}

func (u *useCase) schedule(group *entities.Group) error {
	if err := group.ValidateSchedule(); err != nil {
		return err
	}
	if u.jobRepo.Scheduled(group.ID) {
		return domain.ErrJobAlreadyScheduled
	}
	job := u.newJob(group)
	return u.jobRepo.Add(job)
}

func (u *useCase) newJob(group *entities.Group) *entities.Job {
	f := func() {
		start := time.Now().UTC()
		log.Info().Msgf("[%d] invoked at: %v", group.ID, start)
		if err := u.FixBids(group.ID); err != nil {
			log.Error().Msgf("[%d] failed: %v", group.ID, err)
		} else {
			log.Info().Msgf("[%d] took: %v", group.ID, time.Since(start))
		}
	}
	return entities.
		NewJob(group.ID, group.Start, *group.Interval).
		WithFunc(f)
}
