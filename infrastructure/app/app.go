package app

import (
	"gitlab.jooble.com/marketing_tech/yandex_bidder/config"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/cache"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/log"
)

type (
	shutdowner interface {
		Shutdown()
	}

	App struct {
		config       *config.Config
		cleanupTasks []shutdowner
	}
)

func New(config *config.Config) *App {
	log.ConfigureLogger(config.Logger.Level)

	app := &App{config: config}
	defer app.shutdownOnPanic()
	cacheStore := cache.New(config.Cache.URL())
	app.AddCleanupTask(cacheStore)

	return app
}

func (a *App) AddCleanupTask(s shutdowner) {
	a.cleanupTasks = append(a.cleanupTasks, s)
}

func (a *App) Shutdown() {
	lastIndex := len(a.cleanupTasks) - 1

	for i := range a.cleanupTasks {
		a.cleanupTasks[lastIndex-i].Shutdown()
	}
}

func (a *App) shutdownOnPanic() {
	if r := recover(); r != nil {
		a.Shutdown()
		panic(r)
	}
}
