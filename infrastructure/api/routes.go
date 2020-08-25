package api

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
)

func (h *Handler) Register(i *iris.Application) {
	api := i.Party("/api")
	v1 := api.Party("/v1")

	v1.Get("/ping", h.Ping)
	v1.ConfigureContainer(func(container *router.APIContainer) {
		container.Get("/groups", h.GetGroups)
		container.Put("/groups/{id:int}", h.UpdateGroup)
		container.Post("/groups/{id:int}", h.ChangeBid)

		container.Get("/strategies", h.GetStrategies)

		container.OnError(func(c *context.Context, err error) {
			if err == domain.ErrGroupNotFound {
				c.StatusCode(http.StatusNotFound)
				c.JSON(map[string]string{"details": err.Error()})
			} else {
				c.StatusCode(http.StatusInternalServerError)
				c.JSON(err)
			}
		})
	})
}
