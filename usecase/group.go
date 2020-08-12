package usecase

import "gitlab.jooble.com/marketing_tech/yandex_bidder/domain"

type (
	groupUseCase struct {
		groupRepo GroupRepo
	}

	GroupUseCase interface {
		GetAll() ([]*domain.Group, error)
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
