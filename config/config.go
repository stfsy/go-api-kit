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
	var a AppConfig
	envconfig.MustProcess("API_KIT", &a)

	var s ServerConfig
	envconfig.MustProcess("", &s)

	c = &Configuration{
		AppConfig:    a,
		ServerConfig: s,
	}
}

func Get() Configuration {
	if c == nil {
		Load()
	}

	return *c
}

func reset() {
	c = nil
}

func IsProduction() bool {
	return Get().Env == "production"
}
