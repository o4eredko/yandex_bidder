package strategy

import (
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	useCase struct {
		strategyRepo usecase.StrategyRepo
	}

	UseCase interface {
		GetAll() ([]entities.Strategy, error)
	}
)

func New(strategyRepo usecase.StrategyRepo) UseCase {
	return &useCase{
		strategyRepo: strategyRepo,
	}
}

func (u *useCase) GetAll() ([]entities.Strategy, error) {
	return u.strategyRepo.GetAll()
}
