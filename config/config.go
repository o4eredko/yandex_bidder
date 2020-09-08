package config

import (
	"fmt"

	"github.com/sherifabdlnaby/configuro"
)

type (
	Config struct {
		App       *App
		Scheduler *Scheduler
		Database  *Database
		AMQP      *AMQP
		Api       *Api
		Logger    *Logger
	}

	App struct {
		ConcurrencyLimit int `config:"concurrency_limit"`
	}

	Scheduler struct {
		TimeZone                string `config:"time_zone"`
		SuppressErrorsOnStartup bool   `config:"suppress_errors_on_startup"`
	}

	Api struct {
		Host string
		Port int
	}

	Database struct {
		Host     string
		Port     int
		Username string
		Password string
		DB       string `config:"database"`
	}

	AMQP struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	Logger struct {
		Level  string
		Format string
	}
)

func (c *Api) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Database) DSN() string {
	return fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		c.Host, c.Username, c.Password, c.Port, c.DB)
}

func (c *AMQP) DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", c.Username, c.Password, c.Host, c.Port)
}

func NewConfig(configPath string) *Config {
	conf, err := configuro.NewConfig(configuro.WithLoadFromConfigFile(configPath, true))
	if err != nil {
		panic(err)
	}

	confStruct := &Config{}
	if err := conf.Load(confStruct); err != nil {
		panic(err)
	}

	return confStruct
}
