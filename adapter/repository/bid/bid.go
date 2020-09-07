package bid

import (
	"encoding/json"
	"fmt"

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
	msg, err := json.MarshalIndent(bids, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(msg))

	return r.amqpStore.Publish(msg)
}
