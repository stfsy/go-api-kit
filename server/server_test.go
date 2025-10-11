package server

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stfsy/go-api-kit/config"
	a "github.com/stretchr/testify/assert"
)

// startTestServer starts the server with a simple /test handler and returns a
// function to stop it. Tests should call defer stop().
func startTestServer(_t *testing.T) func() {
	srv := NewServer(&ServerConfig{
		MuxCallback: func(mux *http.ServeMux) {
			mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
				_, e := io.ReadAll(r.Body)
				if e != nil {
					http.Error(w, "error reading body", http.StatusRequestEntityTooLarge)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			})
		},
	})

	go func() {
		// Ignore Start errors in tests; Start logs errors via logger.
		_ = srv.Start()
	}()

	// Wait briefly for server to start. In production tests, prefer sync.
	time.Sleep(300 * time.Millisecond)

	return func() {
		// ensure graceful shutdown
		srv.Stop()
		time.Sleep(50 * time.Millisecond)
	}
}

func TestGET_ReturnsOK(t *testing.T) {
	stop := startTestServer(t)
	defer stop()

	resp, err := http.Get("http://localhost:8080/test")
	a := a.New(t)
	a.NoError(err)
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	a.Equal(http.StatusOK, resp.StatusCode)
	a.Equal("ok", string(b))
}

func TestSecurityHeaders_AreSet(t *testing.T) {
	stop := startTestServer(t)
	defer stop()

	resp, err := http.Get("http://localhost:8080/test")
	a := a.New(t)
	a.NoError(err)
	defer resp.Body.Close()

	// Check a few representative security headers
	a.Equal("nosniff", resp.Header.Get("X-Content-Type-Options"))
	a.Equal("DENY", resp.Header.Get("X-Frame-Options"))
}

func TestPOST_WithoutLengthOrTransferEncoding_Returns411(t *testing.T) {
	stop := startTestServer(t)
	defer stop()

	req, _ := http.NewRequest("POST", "http://localhost:8080/test", nil)
	resp, err := http.DefaultClient.Do(req)
	a := a.New(t)
	a.NoError(err)
	defer resp.Body.Close()
	a.Equal(http.StatusLengthRequired, resp.StatusCode)
}

func TestPOST_WithBodyButNoContentType_Returns415(t *testing.T) {
	stop := startTestServer(t)
	defer stop()

	req, _ := http.NewRequest("POST", "http://localhost:8080/test", strings.NewReader("x"))
	// Don't set Content-Type
	req.Header.Del("Content-Type")
	resp, err := http.DefaultClient.Do(req)
	a := a.New(t)
	a.NoError(err)
	defer resp.Body.Close()
	a.Equal(http.StatusUnsupportedMediaType, resp.StatusCode)
}

func TestCacheHeadersSet(t *testing.T) {
	stop := startTestServer(t)
	defer stop()

	resp, err := http.Get("http://localhost:8080/test")
	assert := a.New(t)
	assert.NoError(err)
	defer resp.Body.Close()

	// check a subset of no-cache headers
	assert.Equal("no-store, no-cache, must-revalidate, proxy-revalidate", resp.Header.Get("Cache-Control"))
	assert.Equal("0", resp.Header.Get("Expires"))
	assert.Equal("no-cache", resp.Header.Get("Pragma"))
}

func TestRequireHTTP11RejectsHTTP10(t *testing.T) {
	stop := startTestServer(t)
	defer stop()

	// Open raw TCP connection and send an HTTP/1.0 request line so the server sees ProtoMajor=1 ProtoMinor=0
	conn, err := net.Dial("tcp", "localhost:8080")
	assert := a.New(t)
	assert.NoError(err)
	defer conn.Close()

	req := "GET /test HTTP/1.0\r\nHost: localhost\r\n\r\n"
	_, err = conn.Write([]byte(req))
	assert.NoError(err)

	// Read status line from response
	r := bufio.NewReader(conn)
	statusLine, err := r.ReadString('\n')
	assert.NoError(err)
	// status line format: HTTP/1.1 400 Bad Request

	fields := strings.Split(statusLine, " ")
	// fields[1] should be the status code
	if len(fields) < 2 {
		t.Fatalf("unexpected status line: %q", statusLine)
	}
	assert.Equal("400", fields[1])
}
func TestMaxBodyLengthEnforced(t *testing.T) {
	stop := startTestServer(t)
	defer stop()

	cfg := config.Get()
	// create a body that is one byte larger than allowed
	big := bytes.Repeat([]byte("x"), cfg.MaxBodySize+1)
	req, _ := http.NewRequest("POST", "http://localhost:8080/test", bytes.NewReader(big))
	// set content-type so content-type middleware allows passing
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	assert := a.New(t)
	assert.NoError(err)
	assert.Equal(http.StatusRequestEntityTooLarge, resp.StatusCode)
	defer resp.Body.Close()
}
