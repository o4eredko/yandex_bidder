package usecase

type (
	bidUseCase struct {
		groupRepo    GroupRepo
		strategyRepo StrategyRepo
		bidRepo      BidRepo
	}

	BidUseCase interface {
		FixBids(groupId int) error
	}
)

func NewBidUseCase(groupRepo GroupRepo, strategyRepo StrategyRepo, bidRepo BidRepo) BidUseCase {
	return &bidUseCase{
		groupRepo:    groupRepo,
		strategyRepo: strategyRepo,
		bidRepo:      bidRepo,
	}
}

func (u *bidUseCase) FixBids(groupId int) error {
	group, err := u.groupRepo.GetByID(groupId)
	if err != nil {
		return err
	}
	strategy, err := u.strategyRepo.Get(*group.Strategy)
	if err != nil {
		return err
	}
	stats, err := u.groupRepo.GetStats(groupId)
	if err != nil {
		return err
	}
	maxRetries := *group.Interval / 2
	err = u.bidRepo.Update(stats, strategy, maxRetries)
	if err != nil {
		return err
	}
	return nil
}
