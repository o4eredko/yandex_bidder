package account

import (
	"fmt"

	dbx "github.com/go-ozzo/ozzo-dbx"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
	sqlStore "gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/sql"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	repo struct {
		store *sqlStore.Store
	}
)

func New(store *sqlStore.Store) usecase.AccountRepo {
	return &repo{
		store: store,
	}
}

func (r *repo) Bids(account *entities.Account, strategy string) ([]*entities.Bid, error) {
	bids := make([]*entities.Bid, 0)
	rows, err := r.store.DB.
		Select("c.id", fmt.Sprintf("dbo.%s(c.id) as bid", strategy)).
		From("campaigns AS c").
		InnerJoin("accounts AS a", dbx.NewExp("a.id = c.account_id")).
		Where(dbx.HashExp{"a.id": account.ID}).
		Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		bid := new(entities.Bid)
		if err := rows.ScanStruct(bid); err != nil {
			return nil, err
		}
		if bid.Bid != nil {
			bids = append(bids, bid)
		}
	}

	return bids, nil
}
