package usecase

import "gitlab.jooble.com/marketing_tech/yandex_bidder/domain"

type (
	GroupRepo interface {
		GetAll() ([]*domain.Group, error)
		GetByID(id int) (*domain.Group, error)
		Accounts(group *domain.Group) ([]*domain.Account, error)
		Update(group *domain.Group) error
	}

	AccountRepo interface {
		Campaigns(account *domain.Account) ([]*domain.Campaign, error)
	}

	StrategyRepo interface {
		GetAll() ([]*domain.Strategy, error)
	}

	BidRepo interface {
		Update(bids *domain.BidsOut) error
		Calculate(strategy string, campaigns []*domain.Campaign) ([]*domain.Bid, error)
	}
)
