package config

import (
	"fmt"

	"github.com/sherifabdlnaby/configuro"
)

type (
	Config struct {
		Api      apiConfig
		Database databaseConfig `config:"database"`
		// Cache    cacheConfig
		Logger loggerConfig
	}

	apiConfig struct {
		Host string
		Port int
	}

	databaseConfig struct {
		Host     string
		Port     int
		Username string
		Password string
		DB       string `config:"database"`
	}

	// cacheConfig struct {
	// 	Host     string
	// 	Port     int
	// 	Username string
	// 	Password string
	// }

	loggerConfig struct {
		Level  string
		Format string
	}
)

func (c *apiConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *databaseConfig) DSN() string {
	return fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		c.Host, c.Username, c.Password, c.Port, c.DB)
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
