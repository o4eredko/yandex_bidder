package amqpRepo

import (
	"encoding/json"

	"github.com/streadway/amqp"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/config"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	Store interface {
		Publish(publishConfig *config.PublishConfig, message *amqp.Publishing) error
		Subscribe(consumeConfig *config.ConsumeConfig, handler func(amqp.Delivery)) error
		Shutdown() error
	}

	repo struct {
		amqpStore Store
		config    *config.AMQP
	}
)

func New(amqpStore Store, config *config.AMQP) usecase.AMQPRepo {
	return &repo{
		amqpStore: amqpStore,
		config:    config,
	}
}

func (r *repo) Update(bids *domain.GroupToUpdateBids) error {
	body, err := json.Marshal(bids)
	if err != nil {
		return err
	}
	msg := &amqp.Publishing{
		Body:        body,
		ContentType: "application/json",
	}
	return r.amqpStore.Publish(r.config.Publishes["change_bid"], msg)
}
