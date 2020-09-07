package usecase

import (
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
)

type (
	groupUseCase struct {
		groupRepo   GroupRepo
		accountRepo AccountRepo
		bidRepo     BidRepo
		jobRepo     JobRepo
	}

	GroupUseCase interface {
		GetAll() ([]*entities.Group, error)
		Update(updateIn *domain.GroupUpdateIn) (*entities.Group, error)
		FixBids(groupID int) error
		Schedule(group *entities.Group) error
		ScheduleAll(stopOnErr bool) error
		Start(groupID int) error
		Pause(groupID int) error
	}
)

func NewGroupUseCase(
	groupRepo GroupRepo,
	accountRepo AccountRepo,
	bidRepo BidRepo,
	jobRepo JobRepo,
) GroupUseCase {
	return &groupUseCase{
		groupRepo:   groupRepo,
		accountRepo: accountRepo,
		bidRepo:     bidRepo,
		jobRepo:     jobRepo,
	}
}

func (u *groupUseCase) GetAll() ([]*entities.Group, error) {
	groups, err := u.groupRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		group.IsScheduled = u.jobRepo.Scheduled(group.ID)
	}

	return groups, nil
}

func (u *groupUseCase) Update(updateIn *domain.GroupUpdateIn) (*entities.Group, error) {
	group, err := u.groupRepo.GetByID(updateIn.ID)
	if err != nil {
		return nil, err
	}

	if err := u.jobRepo.Remove(group.ID); err != nil {
		return nil, err
	}

	group.Start = updateIn.Start
	*group.Interval = updateIn.Interval
	*group.Strategy = updateIn.Strategy
	if err := u.groupRepo.Update(group); err != nil {
		return nil, err
	}

	if err := u.Schedule(group); err != nil {
		return nil, err
	}

	return group, err
}

func (u *groupUseCase) FixBids(groupID int) error {
	group, err := u.groupRepo.GetByID(groupID)
	if err != nil {
		return err
	}

	accounts, err := u.groupRepo.Accounts(group)
	if err != nil {
		return err
	}

	groupToUpdate := new(domain.GroupToUpdateBids)
	for _, account := range accounts {
		bids, err := u.accountRepo.Bids(account, *group.Strategy)
		if err != nil {
			return err
		}

		accountWithBids := &domain.AccountBids{
			AccountName: account.Name,
			Bids:        bids,
		}
		groupToUpdate.Accounts = append(groupToUpdate.Accounts, accountWithBids)
	}
	groupToUpdate.Name = group.Name
	groupToUpdate.MaxRetries = group.MaxRetries()

	if err := u.bidRepo.Update(groupToUpdate); err != nil {
		return err
	}

	return nil
}

func (u *groupUseCase) Schedule(group *entities.Group) error {
	if err := group.ValidateSchedule(); err != nil {
		return err
	}

	task := entities.NewTask(u.FixBids, group.ID)
	job := entities.NewJob(group.ID, group.Start, *group.Interval, task)

	return u.jobRepo.Add(job)
}

func (u *groupUseCase) ScheduleAll(suppressErr bool) error {
	groups, err := u.GetAll()
	if err != nil {
		return err
	}

	for _, group := range groups {
		err := u.Schedule(group)
		if err != nil && !suppressErr {
			return err
		}
	}

	return nil
}

func (u *groupUseCase) Start(groupID int) error {
	group, err := u.groupRepo.GetByID(groupID)
	if err != nil {
		return err
	}
	if scheduled := u.jobRepo.Scheduled(groupID); scheduled {
		return domain.ErrJobAlreadyScheduled
	}

	return u.Schedule(group)
}

func (u *groupUseCase) Pause(groupID int) error {
	if scheduled := u.jobRepo.Scheduled(groupID); !scheduled {
		return domain.ErrJobNotFound
	}

	return u.jobRepo.Remove(groupID)
}
