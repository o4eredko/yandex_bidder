package usecase

import "gitlab.jooble.com/marketing_tech/yandex_bidder/domain"

type (
	strategyUseCase struct {
		strategyRepo StrategyRepo
	}

	StrategyUseCase interface {
		GetStrategies() ([]*domain.Strategy, error)
	}
)

func NewStrategyUseCase(strategyRepo StrategyRepo) StrategyUseCase {
	return &strategyUseCase{
		strategyRepo: strategyRepo,
	}
}

func (u *strategyUseCase) GetStrategies() ([]*domain.Strategy, error) {
	return u.strategyRepo.GetAll()
}
