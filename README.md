<div align="center">

[![contributions - welcome](https://img.shields.io/badge/contributions-welcome-blue/green)](/CONTRIBUTING.md "Go to contributions doc")
[![GitHub License](https://img.shields.io/github/license/stfsy/go-api-kit.svg)](https://github.com/stfsy/go-api-kit/blob/master/LICENSE)
<br/>
[![Go Report Card](https://goreportcard.com/badge/github.com/stfsy/go-api-kit)](https://goreportcard.com/report/github.com/stfsy/go-api-kit)
[![Go](https://img.shields.io/github/go-mod/go-version/stfsy/go-api-kit
)](https://go.dev/ "Go to golang homepage")
</div>

<br/>

# go-api-kit
Kickstarts your API by providing out-of-the-box implementations for must-have modules and components for a successful API.

Provides a solid foundation for SaaS, Client/Server and API products. It provides out-of-the-box mitigations for the [10 OWASP risks for APIs](https://owasp.org/API-Security/editions/2023/en/0x11-t10/):


## 📦 Installation

```bash
go get https://github.com/stfsy/go-api-kit
```

## 🚀 Usage
### main.go
```golang
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/stfsy/go-api-kit/server"
)

var s *server.Server

func main() {
	startServerNonBlocking()
	stopServerAfterSignal()
}

func stopServerAfterSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	s.Stop()

	fmt.Println("Graceful shutdown complete.")
}

func startServerNonBlocking() {
	s = server.NewServer(&server.ServerConfig{
		MuxCallback: func(*http.ServeMux) {
			// add your endpoints and middlewares here
		},
		ListenCallback: func() {
			// do sth just after listen was called on the server instance and
			// just before the server starts serving requests
		},
		// port override is optional but can be used if you want to
		// define the port manually. If empty the value of env.PORT is used.
		PortOverride: "8080",
	})
	go func() {
		err := s.Start()
		if err != nil {
			panic(fmt.Errorf("unable to start server %w", err))
		}
	}()
}
```

## Configuration
This module will read the following environment variables.

### Env Vars
- `API_KIT_ENV`: default=production
- `API_KIT_MAX_BODY_SIZE`: default=10485760 (bytes) = 10 MB
- `API_KIT_READ_TIMEOUT`: default=10 (seconds)
- `API_KIT_WRITE_TIMEOUT`: default=10 (seconds)
- `API_KIT_IDLE_TIMEOUT`: default=620 (seconds)
### Standard Env Vars
- `PORT`: default=8080

## Middlewares
The module provides several ready-to use middlewares which are compatible with e.g. https://github.com/urfave/negroni.

### Access Log Middleware
Logs each incoming request to give insights about usage and response times.

```go
import (
	"net/http"
	"github.com/urfave/negroni"
	"github.com/stfsy/go-api-kit/server/middlewares"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	n := negroni.New()
	n.Use(middlewares.NewAccessLog())
	n.UseHandler(mux)
	http.ListenAndServe(":8080", n)
}
```
[Source](server/middlewares/access-log.go)

### Content Type Middleware
Validates the incoming content type, if the request method implies a state change (e.g. POST).

```go
import (
	"net/http"
	"github.com/urfave/negroni"
	"github.com/stfsy/go-api-kit/server/middlewares"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	n := negroni.New()
	n.Use(middlewares.NewRequireContentTypeMiddleware("application/json"))
	n.UseHandler(mux)
	http.ListenAndServe(":8080", n)
}
```
[Source](server/middlewares/require-content-type.go)

### Max Body Length Middleware
Limits the maximum allowed size of the request body.

```go
import (
	"net/http"
	"github.com/urfave/negroni"
	"github.com/stfsy/go-api-kit/server/middlewares"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	n := negroni.New()
	n.Use(middlewares.NewRequireMaxBodyLengthMiddleware())
	n.UseHandler(mux)
	http.ListenAndServe(":8080", n)
}
```
[Source](server/middlewares/require-body-length.go)

### Security Headers Middleware
Adds additional security headers to the response to prevent common attacks and protect users and their data.

```go
import (
	"net/http"
	"github.com/urfave/negroni"
	"github.com/stfsy/go-api-kit/server/middlewares"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	n := negroni.New()
	n.Use(middlewares.NewRespondWithSecurityHeadersMiddleware())
	n.UseHandler(mux)
	http.ListenAndServe(":8080", n)
}
```
[Source](server/middlewares/respond-with-security-headers.go)

### Upstream Cache Control Middleware
Instructs proxy servers between the client and the API to not cache responses.

```go
import (
	"net/http"
	"github.com/urfave/negroni"
	"github.com/stfsy/go-api-kit/server/middlewares"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	n := negroni.New()
	n.Use(middlewares.NewNoCacheHeadersMiddleware())
	n.UseHandler(mux)
	http.ListenAndServe(":8080", n)
}
```
[Source](server/middlewares/respond-with-no-cache-headers.go)


## Functions

### Response Sender Functions

These functions help you send plain text or JSON responses easily:

#### SendText
Sends a plain text response.

```go
import "github.com/stfsy/go-api-kit/server/handlers"

handlers.SendText(w, "Hello, world!")
```
[Source](server/handlers/response-sender.go)

#### SendJson
Sends a JSON response (sets Content-Type to application/json).

```go
import "github.com/stfsy/go-api-kit/server/handlers"

handlers.SendJson(w, []byte(`{"message":"ok"}`))
```
[Source](server/handlers/response-sender.go)

---

### Response Error Sender Functions

These functions send standardized error responses with the correct HTTP status code and a JSON body. Each function takes an `http.ResponseWriter` and an optional `details` map for additional error info.

#### Example usage:
```go
import "github.com/stfsy/go-api-kit/server/handlers"

handlers.SendBadRequest(w, map[string]string{"zip_code": "must match the required pattern"})
handlers.SendUnauthorized(w, map[string]string{"x-api-key": "must not be null"})
handlers.SendInternalServerError(w, nil)
```
[Source](server/handlers/response-error-sender.go)

#### Available error response functions:

| Status Code | Title                        | Function                   |
|-------------|------------------------------|----------------------------|
| 400         | Bad Request                  | SendBadRequest             |
| 401         | Unauthorized                 | SendUnauthorized           |
| 403         | Forbidden                    | SendForbidden              |
| 404         | Not Found                    | SendNotFound               |
| 405         | Method Not Allowed           | SendMethodNotAllowed       |
| 406         | Not Acceptable               | SendNotAcceptable          |
| 408         | Request Timeout              | SendRequestTimeout         |
| 409         | Conflict                     | SendConflict               |
| 410         | Gone                         | SendGone                   |
| 411         | Length Required              | SendLengthRequired         |
| 412         | Precondition Failed          | SendPreconditionFailed     |
| 413         | Payload Too Large            | SendPayloadTooLarge        |
| 414         | URI Too Long                 | SendURITooLong             |
| 415         | Unsupported Media Type       | SendUnsupportedMediaType   |
| 416         | Range Not Satisfiable        | SendRangeNotSatisfiable    |
| 417         | Expectation Failed           | SendExpectationFailed      |
| 422         | Unprocessable Entity         | SendUnprocessableEntity    |
| 429         | Too Many Requests            | SendTooManyRequests        |
| 500         | Internal Server Error        | SendInternalServerError    |
| 501         | Not Implemented              | SendNotImplemented         |
| 502         | Bad Gateway                  | SendBadGateway             |
| 503         | Service Unavailable          | SendServiceUnavailable     |
| 504         | Gateway Timeout              | SendGatewayTimeout         |
| 505         | HTTP Version Not Supported   | SendHTTPVersionNotSupported|


## 🧪 Running Tests
To run tests, run the following command

```bash
./test.sh
```

## 📄 License
[MIT](https://choosealicense.com/licenses/mit/)
