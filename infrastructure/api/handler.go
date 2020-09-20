package api

import (
	"github.com/kataras/iris/v12"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase/group"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase/scheduler"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase/strategy"
)

type Handler struct {
	groupUseCase     group.UseCase
	strategyUseCase  strategy.UseCase
	schedulerUseCase scheduler.UseCase
}

func NewHandler(
	groupUseCase group.UseCase,
	strategyUseCase strategy.UseCase,
	schedulerUseCase scheduler.UseCase,
) *Handler {
	return &Handler{
		groupUseCase:     groupUseCase,
		strategyUseCase:  strategyUseCase,
		schedulerUseCase: schedulerUseCase,
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

func (h *Handler) GetAll() ([]entities.Strategy, error) {
	return h.strategyUseCase.GetAll()
}

func (h *Handler) ToggleGroup(c iris.Context, id int, input *domain.GroupToggleIn) (err error) {
	if err = input.Validate(); err != nil {
		return err
	}
	switch input.Action {
	case "start":
		err = h.groupUseCase.Start(id)
	case "terminate":
		err = h.groupUseCase.Terminate(id)
	}
	if err == nil {
		c.JSON(iris.Map{"updated": true})
	}
	return
}

func (h *Handler) SchedulerInfo() *domain.SchedulerOut {
	return h.schedulerUseCase.ShowInfo()
}
