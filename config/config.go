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
		Host      string
		Port      int
		Username  string
		Password  string
		Consumes  map[string]*ConsumeConfig
		Publishes map[string]*PublishConfig
	}

	ConsumeConfig struct {
		Name      string
		AutoAck   bool `config:"auto_ack"`
		Exclusive bool
		NoLocal   bool `config:"no_local"`
		NoWait    bool `config:"no_wait"`
		Args      map[string]interface{}
		Exchange  exchangeConfig
	}

	PublishConfig struct {
		Mandatory bool
		Immediate bool
		Exchange  exchangeConfig
	}

	exchangeConfig struct {
		Name       string
		Type       string
		Durable    bool
		AutoDelete bool `config:"auto_delete"`
		Internal   bool
		NoWait     bool   `config:"no_wait"`
		RoutingKey string `config:"routing_key"`
		Args       map[string]interface{}
		Queue      queueConfig
	}

	queueConfig struct {
		Name       string
		Durable    bool
		AutoDelete bool `config:"auto_delete"`
		Exclusive  bool
		Internal   bool
		NoWait     bool `config:"no_wait"`
		Args       map[string]interface{}
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
