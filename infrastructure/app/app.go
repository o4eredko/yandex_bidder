package app

import (
	"github.com/rs/zerolog/log"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/adapter/repository/account"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/adapter/repository/bid"
	groupRepo "gitlab.jooble.com/marketing_tech/yandex_bidder/adapter/repository/group"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/adapter/repository/job"
	strategyRepo "gitlab.jooble.com/marketing_tech/yandex_bidder/adapter/repository/strategy"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/config"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/logger"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/amqp"
	schedulerStore "gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/scheduler"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/sql"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase/group"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase/scheduler"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase/strategy"
)

type (
	shutdowner interface {
		Shutdown() error
	}

	App struct {
		config           *config.Config
		cleanupTasks     []shutdowner
		GroupUseCase     group.UseCase
		StrategyUseCase  strategy.UseCase
		SchedulerUseCase scheduler.UseCase
	}
)

func New(config *config.Config) (*App, error) {
	logger.ConfigureLogger(config.Logger.Level)

	app := &App{config: config}
	defer app.shutdownOnPanic()

	sqlStore := sql.New(config.Database)
	app.AddCleanupTask(sqlStore)
	amqpStore := amqp.New(config.AMQP)
	app.AddCleanupTask(amqpStore)
	schedulerStore := schedulerStore.New(config.Scheduler)
	app.AddCleanupTask(schedulerStore)

	groupRepo := groupRepo.New(sqlStore)
	accountRepo := account.New(sqlStore)
	strategyRepo := strategyRepo.New(sqlStore)
	bidRepo := bid.New(amqpStore)
	jobRepo := job.New(schedulerStore.Cron)

	app.StrategyUseCase = strategy.New(strategyRepo)
	app.GroupUseCase = group.New(config.App.ConcurrencyLimit, groupRepo, accountRepo, bidRepo, jobRepo)
	app.SchedulerUseCase = scheduler.New(jobRepo)

	if err := app.GroupUseCase.ScheduleAll(config.Scheduler.SuppressErrorsOnStartup); err != nil {
		return nil, err
	}

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
