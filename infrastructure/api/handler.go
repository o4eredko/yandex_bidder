package api

import (
	"github.com/kataras/iris/v12"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type Handler struct {
	groupUseCase    usecase.GroupUseCase
	strategyUseCase usecase.StrategyUseCase
}

func NewHandler(
	groupUseCase usecase.GroupUseCase,
	strategyUseCase usecase.StrategyUseCase,
) *Handler {
	return &Handler{
		groupUseCase:    groupUseCase,
		strategyUseCase: strategyUseCase,
	}
}

func (h *Handler) Ping(c iris.Context) {
	c.JSON(iris.Map{"pong": true})
}

func (h *Handler) GetGroups() ([]*entities.Group, error) {
	return h.groupUseCase.GetAll()
}

func (h *Handler) UpdateGroup(id int, input *domain.GroupUpdateIn) (*entities.Group, error) {
	input.ID = id
	return h.groupUseCase.Update(input)
}

func (h *Handler) ChangeBid(id int) error {
	return h.groupUseCase.FixBids(id)
}

func (h *Handler) GetStrategies() ([]*entities.Strategy, error) {
	return h.strategyUseCase.GetStrategies()
}

func (h *Handler) ToggleGroup(id int, input *domain.GroupToggleIn) (err error) {
	if err = input.Validate(); err != nil {
		return
	}

	switch input.Action {
	case "start":
		err = h.groupUseCase.Start(id)
	case "pause":
		err = h.groupUseCase.Pause(id)
	}

	return
}
