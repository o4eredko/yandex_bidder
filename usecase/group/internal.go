package group

import (
	"sync"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
)

func (u *useCase) calculateBids(
	accounts []*entities.Account,
	strategy string,
) ([]*domain.AccountBids, error) {
	accountsWithBids := make([]*domain.AccountBids, 0, len(accounts))
	for _, account := range accounts {
		bids, err := u.accountRepo.Bids(account, strategy)
		if err != nil {
			return nil, err
		}
		accountWithBids := &domain.AccountBids{
			AccountName: account.Name,
			Bids:        bids,
		}
		accountsWithBids = append(accountsWithBids, accountWithBids)
	}

	return accountsWithBids, nil
}

func (u *useCase) bidsWorker(
	wg *sync.WaitGroup,
	strategy string,
	accounts <-chan *entities.Account,
	res chan<- interface{},
) {
	defer wg.Done()
	for account := range accounts {
		bids, err := u.accountRepo.Bids(account, strategy)
		if err != nil {
			res <- err
			continue
		}
		res <- &domain.AccountBids{
			AccountName: account.Name,
			Bids:        bids,
		}
	}
}

func (u *useCase) calculateWithWorkers(
	numOfWorkers int,
	accounts []*entities.Account,
	strategy string,
) ([]*domain.AccountBids, error) {
	accountsWithBids := make([]*domain.AccountBids, 0, len(accounts))
	wg := new(sync.WaitGroup)
	accountChan := make(chan *entities.Account, len(accounts))
	resChan := make(chan interface{}, numOfWorkers)

	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go u.bidsWorker(wg, strategy, accountChan, resChan)
	}
	go func() {
		wg.Wait()
		close(resChan)
	}()

	for _, acc := range accounts {
		accountChan <- acc
	}
	close(accountChan)

	for res := range resChan {
		switch res := res.(type) {
		case *domain.AccountBids:
			accountsWithBids = append(accountsWithBids, res)
		case error:
			return nil, res
		}
	}

	return accountsWithBids, nil
}
