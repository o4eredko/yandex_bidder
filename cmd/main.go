package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/config"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/api"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/app"
)

func main() {
	c := config.NewConfig("config.toml")

	app, err := app.New(c)
	if err != nil {
		panic(err)
	}

	api := api.New(c, app)
	// Start server
	go api.Start()
	// Wait for interrupt signal to gracefully shutdown the server with
	// api timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	api.Shutdown(ctx)
}
