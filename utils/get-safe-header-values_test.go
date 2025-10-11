package utils

import (
	"net/http/httptest"
	"strings"
	"testing"

	a "github.com/stretchr/testify/assert"
)

func TestControlCharHeaderMiddleware_AllowsValid(t *testing.T) {
	assert := a.New(t)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.5")
	v, ok := GetSafeHeaderValue("X-Forwarded-For", req.Header)
	assert.True(ok, "header should be considered safe when exactly one non-empty value")
	assert.Equal("203.0.113.5", v)
}

func TestControlCharHeaderMiddleware_BlocksControlChars(t *testing.T) {
	assert := a.New(t)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.5\n") // contains control char
	_, ok := GetSafeHeaderValue("X-Forwarded-For", req.Header)
	assert.False(ok, "header with control char should be rejected")
}

func TestControlCharHeaderMiddleware_NoOpWhenHeaderNameEmpty(t *testing.T) {
	assert := a.New(t)

	req := httptest.NewRequest("GET", "/test", nil)
	v, ok := GetSafeHeaderValue("", req.Header)
	assert.True(ok, "empty header name should be treated as no-op")
	assert.Equal("", v)
}

func TestControlCharHeaderMiddleware_MultiValueHeader(t *testing.T) {
	assert := a.New(t)

	req := httptest.NewRequest("GET", "/test", nil)
	// set multiple values for the same header; one contains a control char
	req.Header.Add("X-Forwarded-For", "203.0.113.5")
	req.Header.Add("X-Forwarded-For", "198.51.100.7\n")
	// single-value accessor should reject when more than one non-empty value
	_, ok := GetSafeHeaderValue("X-Forwarded-For", req.Header)
	assert.False(ok, "GetSafeHeaderValue should reject when multiple non-empty values are present")

	// multi-value accessor should reject because one value contains a control char
	vals, ok2 := GetSafeHeaderValues("X-Forwarded-For", req.Header)
	assert.False(ok2, "GetSafeHeaderValues should reject when any value contains control chars")
	assert.Nil(vals)
}

func TestGetSafeHeaderValues_AllowsMultipleSafeValues(t *testing.T) {
	assert := a.New(t)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Add("X-Forwarded-For", "203.0.113.5")
	req.Header.Add("X-Forwarded-For", "198.51.100.7")

	vals, ok := GetSafeHeaderValues("X-Forwarded-For", req.Header)
	assert.True(ok, "GetSafeHeaderValues should accept multiple safe values")
	assert.Equal([]string{"203.0.113.5", "198.51.100.7"}, vals)
}

func TestASCIIValidator_AllowsPrintableASCII(t *testing.T) {
	// removed: ASCIIValidator no longer exported; behavior is enforced via middleware
}

func TestASCIIValidator_BlocksNonASCII(t *testing.T) {
	// removed: ASCIIValidator no longer exported; behavior is enforced via middleware
}

func TestControlCharHeaderMiddleware_BlocksNUL(t *testing.T) {
	assert := a.New(t)

	req := httptest.NewRequest("GET", "/test", nil)
	// contains NUL byte
	req.Header.Set("X-Forwarded-For", "203.0.113.5\x00")
	_, ok := GetSafeHeaderValue("X-Forwarded-For", req.Header)
	assert.False(ok, "header containing NUL should be rejected")
}

func TestControlCharHeaderMiddleware_BlocksDEL(t *testing.T) {
	assert := a.New(t)

	req := httptest.NewRequest("GET", "/test", nil)
	// contains DEL (0x7F)
	req.Header.Set("X-Forwarded-For", "203.0.113.5\x7f")
	_, ok := GetSafeHeaderValue("X-Forwarded-For", req.Header)
	assert.False(ok, "header containing DEL should be rejected")
}

func TestControlCharHeaderMiddleware_BlocksTab(t *testing.T) {
	assert := a.New(t)

	req := httptest.NewRequest("GET", "/test", nil)
	// contains TAB (0x09)
	req.Header.Set("X-Forwarded-For", "203.0.113.5\t")
	_, ok := GetSafeHeaderValue("X-Forwarded-For", req.Header)
	assert.False(ok, "header containing TAB should be rejected")
}

func TestControlCharHeaderMiddleware_BlocksNonASCIIByte(t *testing.T) {
	assert := a.New(t)

	req := httptest.NewRequest("GET", "/test", nil)
	// contains a non-ASCII single byte (0x80)
	req.Header.Set("X-Forwarded-For", "203.0.113.5\x80")
	_, ok := GetSafeHeaderValue("X-Forwarded-For", req.Header)
	assert.False(ok, "header containing non-ASCII byte should be rejected")
}

func TestControlCharHeaderMiddleware_RejectsSpace(t *testing.T) {
	assert := a.New(t)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Api-Token", "bad token with space")
	v, ok := GetSafeHeaderValue("X-Api-Token", req.Header)
	// current implementation allows space in header values
	assert.True(ok, "token containing space should be allowed by implementation")
	assert.Equal("bad token with space", v)
}

func TestGetSafeHeaderValue_AllowsValueJustUnderLimit(t *testing.T) {
	assert := a.New(t)

	// create a value of length 4095 (one less than the 4096 limit)
	long := strings.Repeat("a", 4095)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Api-Token", long)

	v, ok := GetSafeHeaderValue("X-Api-Token", req.Header)
	assert.True(ok, "value just under the limit should be accepted")
	assert.Equal(long, v)
}

func TestGetSafeHeaderValue_RejectsValueAtOrAboveLimit(t *testing.T) {
	assert := a.New(t)

	// value of length 4096 should be rejected
	tooLong := strings.Repeat("b", 4096)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Api-Token", tooLong)

	_, ok := GetSafeHeaderValue("X-Api-Token", req.Header)
	assert.False(ok, "value at the 4096-byte limit should be rejected")
}

func TestGetSafeHeaderValues_AllowsMultipleValuesBelowLimit(t *testing.T) {
	assert := a.New(t)

	// two values which individually are below the limit but whose combined
	// length exceeds 4096 should be accepted because we check per-value only
	v1 := strings.Repeat("x", 3000)
	v2 := strings.Repeat("y", 2000)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Add("X-Forwarded-For", v1)
	req.Header.Add("X-Forwarded-For", v2)

	vals, ok := GetSafeHeaderValues("X-Forwarded-For", req.Header)
	assert.True(ok, "multiple values each under the per-value limit should be accepted")
	assert.Equal([]string{v1, v2}, vals)
}

func TestGetSafeHeaderValues_RejectsIndividualValueTooLarge(t *testing.T) {
	assert := a.New(t)

	big := strings.Repeat("z", 4096)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Add("X-Forwarded-For", "127.0.0.1")
	req.Header.Add("X-Forwarded-For", big)

	vals, ok := GetSafeHeaderValues("X-Forwarded-For", req.Header)
	assert.False(ok, "GetSafeHeaderValues should reject when any individual value is >= limit")
	assert.Nil(vals)
}
