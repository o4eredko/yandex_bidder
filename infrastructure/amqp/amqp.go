package amqp

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"

	amqpRepo "gitlab.jooble.com/marketing_tech/yandex_bidder/adapter/repository/amqp"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/config"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/app"
)

type (
	consumer struct {
		store          amqpRepo.Store
		config         *config.Config
		handlerMapping map[string]func(amqp.Delivery)
	}
)

func New(app *app.App, config *config.Config) *consumer {
	handler := NewHandler(app.AMQPStore, config, app.BidUseCase)
	handlerMapping := map[string]func(amqp.Delivery){
		"updated": handler.UpdateBid,
	}
	for key := range handlerMapping {
		if _, ok := config.AMQP.Consumes[key]; !ok {
			panic("No config was found for queue " + key)
		}
	}
	return &consumer{
		store:          app.AMQPStore,
		config:         config,
		handlerMapping: handlerMapping,
	}
}

func (c *consumer) Start() {
	for key, handler := range c.handlerMapping {
		consumeConfig := c.config.AMQP.Consumes[key]
		if err := c.store.Subscribe(consumeConfig, handler); err != nil {
			log.Error().Msgf("Cannot subscribe to exchange: %s", consumeConfig.Exchange.Name)
		}
	}
	log.Info().Msg("AMQP server is running. To exit press CTRL+C")
}

func (c *consumer) Shutdown(ctx context.Context) {
	if err := c.store.Shutdown(); err != nil {
		panic(err)
	}
}
