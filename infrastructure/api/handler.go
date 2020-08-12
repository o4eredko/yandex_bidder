package api

import (
	"github.com/kataras/iris/v12"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type Handler struct {
	groupUseCase usecase.GroupUseCase
}

func NewHandler(groupUseCase usecase.GroupUseCase) *Handler {
	return &Handler{
		groupUseCase: groupUseCase,
	}
}

func (h *Handler) Ping(c iris.Context) {
	c.JSON(iris.Map{"pong": true})
}

func (h *Handler) GetGroups() ([]*domain.Group, error) {
	return h.groupUseCase.GetAll()
}

func (h *Handler) UpdateGroup(id int) (*domain.Group, error) {
	return nil, nil
}

func (h *Handler) ChangeBid(id int) {

}
