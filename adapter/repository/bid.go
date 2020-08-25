package repository

import (
	"fmt"

	dbx "github.com/go-ozzo/ozzo-dbx"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	amqpstore "gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/amqp"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	bidRepo struct {
		db        dbx.Builder
		amqpStore amqpstore.Store
	}
)

func NewBidRepo(db dbx.Builder, amqpStore amqpstore.Store) usecase.BidRepo {
	return &bidRepo{
		db:        db,
		amqpStore: amqpStore,
	}
}

func (r *bidRepo) Calculate(strategy string, campaigns []*domain.Campaign) ([]*domain.Bid, error) {
	query := r.db.NewQuery(fmt.Sprintf("SELECT dbo.%s({:campaign_id}) AS bid", strategy))
	query.Prepare()
	defer query.Close()

	bids := make([]*domain.Bid, 0, len(campaigns))

	for _, campaign := range campaigns {
		bid := new(domain.Bid)
		query.Bind(dbx.Params{"campaign_id": campaign.ID})
		if err := query.One(bid); err != nil {
			return nil, err
		}
		bid.CampaignID = campaign.ID
		bids = append(bids, bid)
	}

	return bids, nil
}

func (r *bidRepo) Update(bids *domain.BidsOut) error {
	// fmt.Printf("Account: %s, MaxRetiries: %d\n", bids.AccountName, bids.MaxRetries)
	// for i, bid := range bids.Bids {
	// 	fmt.Printf("\t#%d => id: %d, bid: %d\n", i, bid.CampaignID, bid.Bid)
	// }
	r.amqpStore.Channel.Publish("", "change_bid", false, false, msg)

	return nil
}
