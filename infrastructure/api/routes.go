package api

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

func (h *Handler) Register(i *iris.Application) {
	api := i.Party("/api")
	v1 := api.Party("/v1")

	v1.Get("/ping", h.Ping)
	v1.ConfigureContainer(func(container *router.APIContainer) {
		container.Get("/groups", h.GetGroups)
		container.Put("/groups/{id:int}", h.UpdateGroup)
		container.Post("/groups/{id:int}", h.ChangeBid)
	})
}
