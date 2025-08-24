package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testPayload struct {
	Name string `json:"name"`
}

func TestValidatingHandler_Success(t *testing.T) {
	payload := testPayload{Name: "test"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handlerCalled := false
	handler := func(w http.ResponseWriter, r *http.Request, p *testPayload) {
		handlerCalled = true
		if p == nil || p.Name != payload.Name {
			t.Errorf("expected name %s, got %v", payload.Name, p)
		}
	}

	ValidatingHandler[testPayload](handler)(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Result().StatusCode)
	}

	if !handlerCalled {
		t.Error("handler was not called on valid payload")
	}
}

func TestValidatingHandler_NoBodyMethods(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodHead, http.MethodOptions}
	for _, method := range methods {
		req := httptest.NewRequest(method, "/", nil)
		w := httptest.NewRecorder()
		handlerCalled := false
		handler := func(w http.ResponseWriter, r *http.Request, p *testPayload) {
			handlerCalled = true
			if p != nil {
				t.Errorf("expected nil payload, got %v", p)
			}
		}
		ValidatingHandler[testPayload](handler)(w, req)
		if !handlerCalled {
			t.Errorf("handler was not called for method %s", method)
		}
		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("expected status %d, got %d for method %s", http.StatusOK, w.Result().StatusCode, method)
		}
	}
}

func TestValidatingHandler_InvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"name":}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(invalidJSON))
	w := httptest.NewRecorder()

	handler := func(w http.ResponseWriter, r *http.Request, p *testPayload) {
		t.Error("handler should not be called on invalid JSON")
	}

	ValidatingHandler[testPayload](handler)(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Result().StatusCode)
	}
}

func TestValidatingHandler_DeleteWithNoBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	w := httptest.NewRecorder()
	handlerCalled := false
	handler := func(w http.ResponseWriter, r *http.Request, p *testPayload) {
		handlerCalled = true
		if p != nil {
			t.Errorf("expected nil payload, got %v", p)
		}
	}
	ValidatingHandler[testPayload](handler)(w, req)
	if !handlerCalled {
		t.Error("handler was not called for DELETE with no body")
	}
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d for DELETE", http.StatusOK, w.Result().StatusCode)
	}
}
