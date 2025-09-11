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
```go
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/stfsy/go-api-kit/server"
	"github.com/urfave/negroni"
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
		// Register endpoints and custom middlewares
		MuxCallback: func(mux *http.ServeMux) {
			mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			})
		},
		// Customize the Negroni middleware stack
		MiddlewareCallback: func(n *negroni.Negroni) *negroni.Negroni {
			// n.Use(myCustomMiddleware)
			return n
		},
		// Called after the server starts listening, before serving requests
		ListenCallback: func() {
			fmt.Println("Server is now listening!")
		},
		// Manually set the port. If empty, uses env.PORT
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

#### `ServerConfig` struct

| Field                | Type                                   | Description                                                                                   |
|----------------------|----------------------------------------|-----------------------------------------------------------------------------------------------|
| `MuxCallback`        | `func(*http.ServeMux)`                 | Register endpoints and custom middlewares to the HTTP mux.                                    |
| `MiddlewareCallback` | `func(*negroni.Negroni) *negroni.Negroni` | Customize the Negroni middleware stack before the server starts.                              |
| `ListenCallback`     | `func()`                               | Called after the server starts listening, before serving requests.                            |
| `PortOverride`       | `string`                               | Manually set the port. If empty, uses the value of the `PORT` environment variable.           |

**Example:**
```go
server.NewServer(&server.ServerConfig{
	MuxCallback: func(mux *http.ServeMux) {
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})
	},
	MiddlewareCallback: func(n *negroni.Negroni) *negroni.Negroni {
		// Add custom middleware if needed
		return n
	},
	ListenCallback: func() {
		fmt.Println("Server is now listening!")
	},
	PortOverride: "8080",
})
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

#### SendStructAsJson
Sends a struct as a JSON response (sets Content-Type to application/json). The struct is marshaled to JSON automatically.

```go
import "github.com/stfsy/go-api-kit/server/handlers"

type MyResponse struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
}

handlers.SendStructAsJson(w, MyResponse{Status: 200, Title: "Success"})
```
[Source](server/handlers/response-sender.go)

---

### Response Error Sender Functions
These functions send standardized error responses with the correct HTTP status code and a JSON body. Each function takes an `http.ResponseWriter` and an optional `details` map for additional error info. The error response now uses the following structure:

```json
{
	"status": 400,
	"title": "Bad Request",
	"details": {
		"zip_code": {
			"validator": "required",
			"message": "must not be undefined",
		}
	}
}
```

#### Example usage:
```go
import "github.com/stfsy/go-api-kit/server/handlers"

// With details (field-level errors):
handlers.SendBadRequest(w, handlers.ErrorDetails{
	"zip_code": {
		"validator": "required",
		"message": "must not be undefined",
	}
})
handlers.SendUnauthorized(w, handlers.ErrorDetails{
	"x-api-key": {
		"message": "must not be null",
	},
})

// Without details (generic error):
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

---

### ValidatingHandler (Generic Request Validation)
Wraps your handler to automatically decode and validate JSON request bodies for POST, PUT, and PATCH methods. For other methods, the handler receives nil as the payload.

To enable JSON payload validation, add https://github.com/go-playground/validator compatible tags to your struct. 

#### Usage
```go
import "github.com/stfsy/go-api-kit/server/handlers"

type MyPayload struct {
    Name string `json:"name" validate:"required"`
}

handlers.ValidatingHandler[MyPayload](func(w http.ResponseWriter, r *http.Request, p *MyPayload) {
    // Use validated payload
    w.Write([]byte(p.Name))
})
```
[Source](server/handlers/validating_handler.go)

#### Validation Errors
Validation errors will be sent to the client automatically with status code `400`. The response will have content type `application/problem+json`. Here's an example:

```json
{
	"status": 400,
	"title": "Bad Request",
	"details": {
		"zip_code": {
			"validator": "required",
			"message": "must not be undefined",
		}
	}
}
```


## 🧪 Running Tests
To run tests, run the following command

```bash
./test.sh
```

## 📄 License
[MIT](https://choosealicense.com/licenses/mit/)
