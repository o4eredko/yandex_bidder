package api

import (
	"net/http"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
)

func (h *Handler) Register(i *iris.Application) {
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodHead, http.MethodOptions,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	i.UseRouter(crs)

	api := i.Party("/api")
	v1 := api.Party("/v1")

	v1.Get("/ping", h.Ping)
	v1.ConfigureContainer(func(container *router.APIContainer) {
		container.Get("/groups", h.GetGroups)
		container.Put("/groups/{id:int}", h.UpdateGroup)
		container.Post("/groups/{id:int}", h.ChangeBid)
		container.Put("/groups/{id:int}/state", h.ToggleGroup)
		container.Put("/groups/{id:int}/bids", h.ChangeBid)

		container.Get("/strategies", h.GetAll)

		container.Get("/scheduler", h.SchedulerInfo)

		container.OnError(func(c *context.Context, err error) {
			if err == domain.ErrGroupNotFound || err == domain.ErrJobNotFound {
				c.StatusCode(http.StatusNotFound)
				c.JSON(map[string]string{"details": err.Error()})
			} else if err == domain.ErrJobAlreadyScheduled {
				c.StatusCode(http.StatusBadRequest)
				c.JSON(map[string]string{"details": err.Error()})
			} else {
				c.StatusCode(http.StatusInternalServerError)
				c.JSON(err)
			}
		})
	})
}
