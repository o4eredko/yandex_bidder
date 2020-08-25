package usecase

import "gitlab.jooble.com/marketing_tech/yandex_bidder/domain"

type (
	bidUseCase struct {
		groupRepo   GroupRepo
		accountRepo AccountRepo
		bidRepo     BidRepo
	}

	BidUseCase interface {
		FixBids(groupID int) error
	}
)

func NewBidUseCase(groupRepo GroupRepo, accountRepo AccountRepo, bidRepo BidRepo) BidUseCase {
	return &bidUseCase{
		groupRepo:   groupRepo,
		accountRepo: accountRepo,
		bidRepo:     bidRepo,
	}
}

func (u *bidUseCase) FixBids(groupID int) error {
	group, err := u.groupRepo.GetByID(groupID)
	if err != nil {
		return err
	}

	accounts, err := u.groupRepo.Accounts(group)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		campaigns, err := u.accountRepo.Campaigns(account)
		if err != nil {
			return err
		}

		bids, err := u.bidRepo.Calculate(*group.Strategy, campaigns)
		if err != nil {
			return err
		}

		bidsOut := &domain.BidsOut{
			AccountName: account.Name,
			Bids:        bids,
			MaxRetries:  group.CalculateMaxRetries(),
		}

		if err := u.bidRepo.Update(bidsOut); err != nil {
			return err
		}
	}

	return nil
}
