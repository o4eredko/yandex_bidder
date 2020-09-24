package bid

import (
	dbx "github.com/go-ozzo/ozzo-dbx"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
	sqlStore "gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/sql"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	repo struct {
		sqlStore *sqlStore.Store
	}
)

func New(sqlStore *sqlStore.Store) usecase.BidRepo {
	return &repo{
		sqlStore: sqlStore,
	}
}

func (r *repo) Update(bid *entities.Bid) error {
	params := dbx.Params{"bid": bid.Bid}
	where := dbx.HashExp{"id": bid.CampaignID}
	query := r.sqlStore.DB.Update("campaigns", params, where)
	_, err := query.Execute()
	return err
}
