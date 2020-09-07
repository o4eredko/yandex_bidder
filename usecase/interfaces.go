package usecase

import (
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
)

type (
	GroupRepo interface {
		GetAll() ([]*entities.Group, error)
		GetByID(id int) (*entities.Group, error)
		Accounts(group *entities.Group) ([]*entities.Account, error)
		Update(group *entities.Group) error
	}

	AccountRepo interface {
		Bids(account *entities.Account, strategy string) ([]*entities.Bid, error)
	}

	StrategyRepo interface {
		GetAll() ([]*entities.Strategy, error)
	}

	BidRepo interface {
		Update(bids *domain.GroupToUpdateBids) error
	}

	JobRepo interface {
		Add(job *entities.Job) error
		Remove(jobID int) error
		Scheduled(jobID int) bool
	}
)
