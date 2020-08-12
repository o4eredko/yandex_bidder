package repository

import (
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	bidRepo struct {
		amqpStore AMQPStore
	}

	AMQPStore interface {
	}
)

func NewBidRepo(amqpStore AMQPStore) usecase.BidRepo {
	return &bidRepo{amqpStore: amqpStore}
}

func (b bidRepo) Update(stats []*domain.Stats, strategy domain.BidHandler, maxRetries int) error {
	result := make(map[string]*domain.BidsOut)
	for _, item := range stats {
		newBid := &domain.Bid{
			Bid:        strategy(item.Clicks, item.Cost, item.Impressions),
			CampaignId: item.CampaignId,
		}

		value, exists := result[item.AccountName]
		if !exists {
			result[item.AccountName] = &domain.BidsOut{
				Bids:       make([]*domain.Bid, 0),
				MaxRetries: maxRetries,
			}
			value = result[item.AccountName]
		}
		value.Bids = append(value.Bids, newBid)
	}
	return nil
}
