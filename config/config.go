package config

import (
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Env string `default:"production"`
}

type ServerConfig struct {
	Port string `default:"8080"`
}

type Configuration struct {
	AppConfig
	ServerConfig
}

var c *Configuration

func Load() {
	var s Configuration
	envconfig.MustProcess("API_KIT", &s)

	c = &s
}

func Get() Configuration {
	if c == nil {
		Load()
	}

	return *c
}

func IsProduction() bool {
	return Get().Env == "production"
}
