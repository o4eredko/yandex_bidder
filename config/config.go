package config

import (
	"fmt"

	"github.com/sherifabdlnaby/configuro"
)

type (
	Config struct {
		Api      apiConfig
		Cache    cacheConfig
		Logger   loggerConfig
	}

	apiConfig struct {
		Host string
		Port int
	}

	cacheConfig struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	loggerConfig struct {
		Level  string
		Format string
	}
)

func (c *apiConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *cacheConfig) URL() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
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
