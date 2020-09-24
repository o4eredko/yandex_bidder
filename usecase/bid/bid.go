package bid

import (
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	useCase struct {
		bidRepo usecase.BidRepo
	}

	UseCase interface {
		Update(bid *entities.Bid) error
	}
)

func New(bidRepo usecase.BidRepo) UseCase {
	return &useCase{
		bidRepo: bidRepo,
	}
}

func (u useCase) Update(bid *entities.Bid) error {
	return u.bidRepo.Update(bid)
}
