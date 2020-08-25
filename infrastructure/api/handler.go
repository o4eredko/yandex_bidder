package api

import (
	"github.com/kataras/iris/v12"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type Handler struct {
	groupUseCase    usecase.GroupUseCase
	strategyUseCase usecase.StrategyUseCase
	bidUseCase      usecase.BidUseCase
}

func NewHandler(
	groupUseCase usecase.GroupUseCase,
	strategyUseCase usecase.StrategyUseCase,
	bidUseCase usecase.BidUseCase,
) *Handler {
	return &Handler{
		groupUseCase:    groupUseCase,
		strategyUseCase: strategyUseCase,
		bidUseCase:      bidUseCase,
	}
}

func (h *Handler) Ping(c iris.Context) {
	c.JSON(iris.Map{"pong": true})
}

func (h *Handler) GetGroups() ([]*domain.Group, error) {
	return h.groupUseCase.GetAll()
}

func (h *Handler) UpdateGroup(id int, input *domain.GroupUpdateIn) (*domain.Group, error) {
	group, err := h.groupUseCase.Update(id, input)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (h *Handler) ChangeBid(id int) error {
	return h.bidUseCase.FixBids(id)
}

func (h *Handler) GetStrategies() ([]*domain.Strategy, error) {
	return h.strategyUseCase.GetStrategies()
}
