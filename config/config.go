package config

import (
	"github.com/kelseyhightower/envconfig"
)

// AppConfig holds application-level configuration.
type AppConfig struct {
	Env string `default:"production"`
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	MaxBodySize  int `default:"10485760" split_words:"true"`
	ReadTimeout  int `default:"10" split_words:"true"`  // seconds
	WriteTimeout int `default:"10" split_words:"true"`  // seconds
	IdleTimeout  int `default:"620" split_words:"true"` // seconds
}

// ContainerConfig holds container-specific configuration.
type ContainerConfig struct {
	Port string `default:"8080"`
}

// Configuration aggregates all config sections.
type Configuration struct {
	AppConfig
	ServerConfig
	ContainerConfig
}

var c *Configuration

// Load loads configuration from environment variables.
// Returns an error if loading fails.
func Load() error {
	var a AppConfig
	err := envconfig.Process("API_KIT", &a)
	if err != nil {
		return err
	}

	var s ServerConfig
	err = envconfig.Process("API_KIT", &s)
	if err != nil {
		return err
	}

	var ct ContainerConfig
	err = envconfig.Process("", &ct)
	if err != nil {
		return err
	}

	c = &Configuration{
		AppConfig:       a,
		ServerConfig:    s,
		ContainerConfig: ct,
	}
	return nil
}

// Get returns the loaded configuration, loading it if necessary.
// Panics if loading fails.
func Get() Configuration {
	if c == nil {
		err := Load()
		if err != nil {
			panic(err)
		}
	}
	return *c
}

// reset clears the loaded configuration (for testing).
func reset() {
	c = nil
}

// IsProduction returns true if the environment is production.
func IsProduction() bool {
	return Get().Env == "production"
}
