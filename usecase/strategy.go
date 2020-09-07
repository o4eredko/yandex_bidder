package usecase

import (
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
)

type (
	strategyUseCase struct {
		strategyRepo StrategyRepo
	}

	StrategyUseCase interface {
		GetStrategies() ([]*entities.Strategy, error)
	}
)

func NewStrategyUseCase(strategyRepo StrategyRepo) StrategyUseCase {
	return &strategyUseCase{
		strategyRepo: strategyRepo,
	}
}

func (u *strategyUseCase) GetStrategies() ([]*entities.Strategy, error) {
	return u.strategyRepo.GetAll()
}
