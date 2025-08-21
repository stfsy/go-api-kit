package config

import (
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Env string `default:"production"`
}

type ServerConfig struct {
	MaxBodySize  int `default:"10485760" split_words:"true"`
	ReadTimeout  int `default:"10" split_words:"true"`  // seconds
	WriteTimeout int `default:"10" split_words:"true"`  // seconds
	IdleTimeout  int `default:"620" split_words:"true"` // seconds
}

type ContainerConfig struct {
	Port string `default:"8080"`
}

type Configuration struct {
	AppConfig
	ServerConfig
	ContainerConfig
}

var c *Configuration

func Load() {
	var a AppConfig
	envconfig.MustProcess("API_KIT", &a)

	var s ContainerConfig
	envconfig.MustProcess("", &s)

	c = &Configuration{
		AppConfig:       a,
		ContainerConfig: s,
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
