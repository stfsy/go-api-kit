package utils

import (
	"strings"
	"testing"

	a "github.com/stretchr/testify/assert"
)

func TestGetSafeValue_AllowsPrintableASCII(t *testing.T) {
	assert := a.New(t)
	v, ok := GetSafeValue("abcDEF123_-.")
	assert.True(ok)
	assert.Equal("abcDEF123_-.", v)
}

func TestGetSafeValue_RejectsControl(t *testing.T) {
	assert := a.New(t)
	_, ok := GetSafeValue("bad\nvalue")
	assert.False(ok)
}

func TestGetSafeValue_RejectsNonASCII(t *testing.T) {
	assert := a.New(t)
	// include a non-ASCII byte 0x80
	_, ok := GetSafeValue("foo\x80bar")
	assert.False(ok)
}

func TestGetSafeValue_RejectsSpace(t *testing.T) {
	assert := a.New(t)
	v, ok := GetSafeValue("has space")
	assert.True(ok)
	assert.Equal("has space", v)
}

func TestGetSafeValue_RespectsSizeLimit(t *testing.T) {
	assert := a.New(t)
	okStr := strings.Repeat("a", 4095)
	v, ok := GetSafeValue(okStr)
	assert.True(ok)
	assert.Equal(okStr, v)

	tooLong := strings.Repeat("b", 4096)
	_, ok2 := GetSafeValue(tooLong)
	assert.False(ok2)
}
