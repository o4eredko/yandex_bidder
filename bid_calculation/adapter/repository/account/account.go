package account

import (
	"fmt"

	"github.com/rs/zerolog/log"

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
		Select("c.id", fmt.Sprintf("dbo.%s(c.id) as bid", strategy), "bid as old_bid").
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
		if bid.Bid == nil {
			log.Warn().Msgf("bid for campaign_id=%v not found", bid.CampaignID)
			continue
		}

		if bid.OldBid == nil {
			log.Info().Msgf(
				"account=%s, campaign_id=%d, old_bid=nil, new_bid=%d",
				account.Name, bid.CampaignID, *bid.Bid,
			)
		} else {
			log.Info().Msgf(
				"account=%s, campaign_id=%d, old_bid=%d, new_bid=%d",
				account.Name, bid.CampaignID, *bid.OldBid, *bid.Bid,
			)
		}

		bids = append(bids, bid)
	}

	return bids, nil
}
