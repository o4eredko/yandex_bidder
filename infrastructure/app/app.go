package app

import (
	"github.com/rs/zerolog/log"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/adapter/repository"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/config"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/logger"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/amqp"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/sql"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	shutdowner interface {
		Shutdown() error
	}

	App struct {
		config          *config.Config
		cleanupTasks    []shutdowner
		GroupUseCase    usecase.GroupUseCase
		BidUseCase      usecase.BidUseCase
		StrategyUseCase usecase.StrategyUseCase
	}
)

func New(config *config.Config) (*App, error) {
	logger.ConfigureLogger(config.Logger.Level)

	app := &App{config: config}
	defer app.shutdownOnPanic()

	sqlStore := sql.New(config.Database.DSN())
	app.AddCleanupTask(sqlStore)
	amqpStore := amqp.New(config.AMQP.DSN())
	app.AddCleanupTask(sqlStore)

	groupRepo := repository.NewGroupRepo(sqlStore.DB)
	accountRepo := repository.NewAccountRepo(sqlStore.DB)
	strategyRepo := repository.NewStrategyRepo(sqlStore.DB)
	bidRepo := repository.NewBidRepo(sqlStore.DB, amqpStore)

	app.StrategyUseCase = usecase.NewStrategyUseCase(strategyRepo)
	app.GroupUseCase = usecase.NewGroupUseCase(groupRepo)
	app.BidUseCase = usecase.NewBidUseCase(groupRepo, accountRepo, bidRepo)

	return app, nil
}

func (a *App) AddCleanupTask(s shutdowner) {
	a.cleanupTasks = append(a.cleanupTasks, s)
}

func (a *App) Shutdown() {
	lastIndex := len(a.cleanupTasks) - 1

	for i := range a.cleanupTasks {
		if err := a.cleanupTasks[lastIndex-i].Shutdown(); err != nil {
			log.Error().Err(err)
		}
	}
}

func (a *App) shutdownOnPanic() {
	if r := recover(); r != nil {
		a.Shutdown()
		panic(r)
	}
}
