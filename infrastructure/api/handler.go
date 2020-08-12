package api

import (
	"github.com/kataras/iris/v12"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Ping(c iris.Context) {
	c.JSON(iris.Map{"pong": true})
}

func (h *Handler) GetGroups() {
}

func (h *Handler) UpdateGroup(id int) (*domain.Group, error) {
	return nil, nil
}

func (h *Handler) ChangeBid(id int) {

}
