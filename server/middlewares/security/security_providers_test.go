package security

import "testing"

func TestHeaderProvidersReturnExpectedPairs(t *testing.T) {
	cases := []struct {
		name      string
		p         HeaderKeyValueProvider
		wantName  string
		wantValue string
	}{
		{"x-content-type", NewXContentTypeOptions(), "X-Content-Type-Options", "nosniff"},
		{"referrer", NewReferrerPolicy(), "Referrer-Policy", "same-origin"},
		{"cors-resource", NewCrossOriginResourcePolicy(), "Cross-Origin-Resource-Policy", "same-site"},
		{"x-xss", NewXssProtection(), "X-XSS-Protection", "1; mode=block"},
		{"x-permitted", NewXPermittedCrossDomainOptions(), "X-Permitted-Cross-Domain-Policies", "none"},
		{"x-frame", NewXFrameOptions(), "X-Frame-Options", "DENY"},
		{"x-download", NewXDownloadOptions(), "X-Download-Options", "noopen"},
		{"sts", NewStrictTransportSecurityPolicy(), "Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload"},
		{"oac", NewOriginAgentClusterPolicy(), "Origin-Agent-Cluster", "?1"},
		{"coop", NewCrossOriginOpenerPolicy(), "Cross-Origin-Opener-Policy", "same-origin"},
		{"coep", NewCrossOriginEmbedderPolicy(), "Cross-Origin-Embedder-Policy", "require-corp"},
	}

	for _, c := range cases {
		kv := c.p.HeaderKeyValue
		if kv.Name != c.wantName {
			t.Fatalf("%s: Name = %q; want %q", c.name, kv.Name, c.wantName)
		}
		if kv.Value != c.wantValue {
			t.Fatalf("%s: Value = %q; want %q", c.name, kv.Value, c.wantValue)
		}
	}
}
