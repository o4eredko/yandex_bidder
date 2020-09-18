package bid

import (
	"encoding/json"
	"github.com/rs/zerolog/log"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	amqpStore "gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/amqp"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	repo struct {
		amqpStore *amqpStore.Store
	}
)

func New(amqpStore *amqpStore.Store) usecase.BidRepo {
	return &repo{
		amqpStore: amqpStore,
	}
}

func (r *repo) Update(bids *domain.GroupToUpdateBids) error {
	msg, err := json.Marshal(bids)
	if err != nil {
		return err
	}
	log.Info().Msgf("bids for update: %v", string(msg))

	return r.amqpStore.Publish(msg)
}
