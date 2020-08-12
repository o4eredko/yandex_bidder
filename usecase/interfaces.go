package usecase

import "gitlab.jooble.com/marketing_tech/yandex_bidder/domain"

type (
	GroupRepo interface {
		GetAll() ([]*domain.Group, error)
		GetByID(int) (*domain.Group, error)
		GetStats(int) ([]*domain.Stats, error)
	}

	StrategyRepo interface {
		Get(string) (domain.BidHandler, error)
	}

	BidRepo interface {
		Update([]*domain.Stats, domain.BidHandler, int) error
	}
)
