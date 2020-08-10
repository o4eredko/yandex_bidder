package api

import (
	"github.com/kataras/iris/v12"
)

func (h *Handler) Register(i *iris.Application) {
	api := i.Party("/api")
	v1 := api.Party("/v1")

	v1.Get("/ping", h.Ping)
}
