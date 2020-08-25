package usecase

import "gitlab.jooble.com/marketing_tech/yandex_bidder/domain"

type (
	groupUseCase struct {
		groupRepo GroupRepo
	}

	GroupUseCase interface {
		GetAll() ([]*domain.Group, error)
		Update(id int, group *domain.GroupUpdateIn) (*domain.Group, error)
	}
)

func NewGroupUseCase(groupRepo GroupRepo) GroupUseCase {
	return &groupUseCase{
		groupRepo: groupRepo,
	}
}

func (u *groupUseCase) GetAll() ([]*domain.Group, error) {
	return u.groupRepo.GetAll()
}

func (u *groupUseCase) Update(id int, input *domain.GroupUpdateIn) (*domain.Group, error) {
	group, err := u.groupRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := input.Validate(); err != nil {
		return nil, err
	}

	group.Strategy = &input.StrategyName
	group.Start = input.ScheduleStart
	group.Interval = &input.ScheduleInterval

	if err := u.groupRepo.Update(group); err != nil {
		return nil, err
	}

	return group, err
}
