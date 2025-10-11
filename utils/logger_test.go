package utils

import (
	"log/slog"
	"testing"
)

func TestNewLogger_NonProduction(t *testing.T) {
	old := isProduction
	defer func() { isProduction = old }()
	isProduction = false

	l := NewLogger("testcomp")
	if l == nil {
		t.Fatalf("expected non-nil logger")
	}

	// ensure With doesn't panic and returns a slog.Logger
	l2 := l.With(slog.String("k", "v"))
	if l2 == nil {
		t.Fatalf("expected With to return non-nil logger")
	}
}

func TestNewLogger_Production(t *testing.T) {
	old := isProduction
	defer func() { isProduction = old }()
	isProduction = true

	l := NewLogger("testcomp")
	if l == nil {
		t.Fatalf("expected non-nil logger in production")
	}
	// basic smoke test: ensure logging a message doesn't panic
	l.Info("smoke", slog.String("k", "v"))
}
