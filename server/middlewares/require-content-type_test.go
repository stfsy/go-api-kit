package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/urfave/negroni"
)

func TestContentType(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		name                string
		inputValue          string
		allowedContentTypes string
		want                int
		method              string
		body                string
	}{
		{
			name:                "should accept requests with a matching content type",
			inputValue:          "application/json; charset=UTF-8",
			allowedContentTypes: "application/json",
			want:                http.StatusOK,
			method:              http.MethodPost,
			body:                "{}",
		},
		{
			name:                "should accept requests with a matching content type no charset",
			inputValue:          "application/json",
			allowedContentTypes: "application/json",
			want:                http.StatusOK,
			method:              http.MethodPost,
			body:                "{}",
		},
		{
			name:                "should accept requests with a matching content-type with extra values",
			inputValue:          "application/json; foo=bar; charset=UTF-8; spam=eggs",
			allowedContentTypes: "application/json",
			want:                http.StatusOK,
			method:              http.MethodPost,
			body:                "{}",
		},
		{
			name:                "should accept requests with a matching content type when multiple content types are supported",
			inputValue:          "text/xml; charset=UTF-8",
			allowedContentTypes: "text/xml",
			want:                http.StatusOK,
			method:              http.MethodPost,
			body:                "<xml></xml>",
		},
		{
			name:                "should not accept requests with a mismatching content type",
			inputValue:          "text/plain; charset=latin-1",
			allowedContentTypes: "application/json",
			want:                http.StatusUnsupportedMediaType,
			method:              http.MethodPost,
			body:                "plain text",
		},
		{
			name:                "should not accept requests with a mismatching content type even if multiple content types are allowed",
			inputValue:          "text/plain; charset=Latin-1",
			allowedContentTypes: "text/xml",
			want:                http.StatusUnsupportedMediaType,
			method:              http.MethodPost,
			body:                "plain text",
		},
		// Additional cases for PATCH, PUT, DELETE, and empty body
		{
			name:                "PATCH with allowed content type",
			inputValue:          "application/json",
			allowedContentTypes: "application/json",
			want:                http.StatusOK,
			method:              http.MethodPatch,
			body:                "{}",
		},
		{
			name:                "PUT with disallowed content type",
			inputValue:          "text/plain",
			allowedContentTypes: "application/json",
			want:                http.StatusUnsupportedMediaType,
			method:              http.MethodPut,
			body:                "plain text",
		},
		{
			name:                "DELETE with body and allowed content type",
			inputValue:          "application/json",
			allowedContentTypes: "application/json",
			want:                http.StatusOK,
			method:              http.MethodDelete,
			body:                "{}",
		},
		{
			name:                "Content-Type with charset parameter",
			inputValue:          "application/json; charset=utf-8",
			allowedContentTypes: "application/json",
			want:                http.StatusOK,
			method:              http.MethodPost,
			body:                "{}",
		},
		{
			name:                "POST with empty body should skip enforcement",
			inputValue:          "application/json",
			allowedContentTypes: "application/json",
			want:                http.StatusOK,
			method:              http.MethodPost,
			body:                "",
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()

			r := negroni.New()
			r.Use(NewRequireContentTypeMiddleware(tt.allowedContentTypes))

			req := httptest.NewRequest(tt.method, "/", bytes.NewReader([]byte(tt.body)))
			req.Header.Set("Content-Type", tt.inputValue)

			r.ServeHTTP(recorder, req)
			res := recorder.Result()

			if res.StatusCode != tt.want {
				t.Errorf("response is incorrect, got %d, want %d", recorder.Code, tt.want)
			}
		})
	}
}
