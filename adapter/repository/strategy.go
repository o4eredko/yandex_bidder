package repository

import (
	"errors"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type strategyRepo struct {
	handlers map[string]domain.BidHandler
}

func NewStrategyRepo() usecase.StrategyRepo {
	return &strategyRepo{
		handlers: map[string]domain.BidHandler{
			"default": defaultHandler,
		},
	}
}

func (r *strategyRepo) Get(key string) (domain.BidHandler, error) {
	handler, ok := r.handlers[key]
	if !ok {
		return nil, errors.New("handler not found")
	}
	return handler, nil
}

func defaultHandler(a, b, c int) int {
	return 228
}
