package handlers

import (
	"testing"

	a "github.com/stretchr/testify/assert"
)

func TestCreateErrorDetail_ReturnsCorrectStructure(t *testing.T) {
	assert := a.New(t)

	ed := CreateErrorDetail("zip_code", "must match pattern")

	// top-level key
	assert.Len(ed, 1)
	inner, ok := ed["zip_code"]
	assert.True(ok, "expected key 'zip_code' to be present")

	// inner ErrorDetail contains the message field with expected value
	assert.Equal("must match pattern", inner.Message)
}

func TestCreateMustNotBeUndefinedDetail_Message(t *testing.T) {
	assert := a.New(t)

	ed := CreateMustNotBeUndefinedErrorDetail("email")
	inner, ok := ed["email"]
	assert.True(ok)
	assert.Equal("must not be undefined", inner.Message)
}

func TestCreateErrorDetail_AllowsEmptyKey(t *testing.T) {
	assert := a.New(t)

	ed := CreateErrorDetail("", "value")
	inner, ok := ed[""]
	assert.True(ok, "empty string key should be present in returned ErrorDetails")
	assert.Equal("value", inner.Message)
}

func TestCreateErrorDetail_ReturnsIndependentMaps(t *testing.T) {
	assert := a.New(t)

	a1 := CreateErrorDetail("k", "v1")
	a2 := CreateErrorDetail("k", "v2")

	// mutate first result
	ed1 := a1["k"]
	ed1.Message = "mutated"
	a1["k"] = ed1

	// ensure second result remains unchanged
	assert.Equal("mutated", a1["k"].Message)
	assert.Equal("v2", a2["k"].Message, "second call should not be affected by mutations to the first result")
}
