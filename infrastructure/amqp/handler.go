package amqp

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"

	amqpRepo "gitlab.jooble.com/marketing_tech/yandex_bidder/adapter/repository/amqp"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/config"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase/bid"
)

type (
	handler struct {
		store      amqpRepo.Store
		config     *config.Config
		bidUseCase bid.UseCase
	}

	Handler interface {
		UpdateBid(amqp.Delivery)
	}
)

func NewHandler(store amqpRepo.Store, config *config.Config, bidUseCase bid.UseCase) Handler {
	return &handler{
		store:      store,
		config:     config,
		bidUseCase: bidUseCase,
	}
}

func (h *handler) UpdateBid(message amqp.Delivery) {
	bid := new(entities.Bid)

	if err := json.Unmarshal(message.Body, bid); err != nil {
		log.Error().Msgf("Cannot map message: %s to json", message.Body)
		return
	}
	log.Info().Msgf("received bid: %v", bid)

	err := h.bidUseCase.Update(bid)
	if err != nil {
		log.Error().Msgf("Cannot update bid for campaign: %s, error: %v", bid.CampaignID, err)
		return
	}
	log.Info().Msgf("Updated bid: %v", bid)
}
