package config

import (
	"testing"

	a "github.com/stretchr/testify/assert"
)

func TestDefaultPort(t *testing.T) {
	assert := a.New(t)

	c := Get()

	assert.Equal("8080", c.Port)
}

func TestReadPortEnvVar(t *testing.T) {
	assert := a.New(t)

	reset()
	t.Setenv("PORT", "8081")
	c := Get()

	assert.Equal("8081", c.Port)
}

func TestDefaultEnv(t *testing.T) {
	assert := a.New(t)

	c := Get()

	assert.Equal("production", c.Env)
}

func TestIsProductionByDefault(t *testing.T) {
	assert := a.New(t)

	_ = Get()

	assert.Equal(true, IsProduction())
}

func TestReadEnvEnvVar(t *testing.T) {
	assert := a.New(t)

	reset()
	t.Setenv("API_KIT_ENV", "development")
	c := Get()

	assert.Equal("development", c.Env)
}

func TestIsNotProduction(t *testing.T) {
	assert := a.New(t)

	reset()
	t.Setenv("API_KIT_ENV", "development")
	_ = Get()

	assert.Equal(false, IsProduction())
}
