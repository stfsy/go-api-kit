package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type HttpError struct {
	Title   string `json:"title"`
	Status  int    `json:"status"`
	Details any    `json:"details,omitempty"`
}

type ErrorDetails map[string]ErrorDetail

type ErrorDetail struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

type staticErrorResponse struct {
	title         string
	payload       []byte
	contentLength string
	status        int
}

func newStaticError(status int, title string) staticErrorResponse {
	b, _ := json.Marshal(HttpError{Title: title, Status: status})
	b = append(b, '\n')
	return staticErrorResponse{
		title:         title,
		payload:       b,
		contentLength: strconv.Itoa(len(b)),
		status:        status,
	}
}

var staticErrors = map[int]staticErrorResponse{
	http.StatusBadRequest:           newStaticError(http.StatusBadRequest, "Bad Request"),
	http.StatusUnauthorized:         newStaticError(http.StatusUnauthorized, "Unauthorized"),
	http.StatusForbidden:            newStaticError(http.StatusForbidden, "Forbidden"),
	http.StatusNotFound:             newStaticError(http.StatusNotFound, "Not Found"),
	http.StatusMethodNotAllowed:     newStaticError(http.StatusMethodNotAllowed, "Method Not Allowed"),
	http.StatusNotAcceptable:        newStaticError(http.StatusNotAcceptable, "Not Acceptable"),
	http.StatusRequestTimeout:       newStaticError(http.StatusRequestTimeout, "Request Timeout"),
	http.StatusConflict:             newStaticError(http.StatusConflict, "Conflict"),
	http.StatusGone:                 newStaticError(http.StatusGone, "Gone"),
	http.StatusLengthRequired:       newStaticError(http.StatusLengthRequired, "Length Required"),
	http.StatusPreconditionFailed:   newStaticError(http.StatusPreconditionFailed, "Precondition Failed"),
	413:                             newStaticError(413, "Payload Too Large"),
	414:                             newStaticError(414, "URI Too Long"),
	http.StatusUnsupportedMediaType: newStaticError(http.StatusUnsupportedMediaType, "Unsupported Media Type"),
	416:                             newStaticError(416, "Range Not Satisfiable"),
	http.StatusExpectationFailed:    newStaticError(http.StatusExpectationFailed, "Expectation Failed"),
	http.StatusUnprocessableEntity:  newStaticError(http.StatusUnprocessableEntity, "Unprocessable Entity"),
	http.StatusTooManyRequests:      newStaticError(http.StatusTooManyRequests, "Too Many Requests"),
	http.StatusInternalServerError:  newStaticError(http.StatusInternalServerError, "Internal Server Error"),
	http.StatusNotImplemented:        newStaticError(http.StatusNotImplemented, "Not Implemented"),
	http.StatusBadGateway:           newStaticError(http.StatusBadGateway, "Bad Gateway"),
	http.StatusServiceUnavailable:   newStaticError(http.StatusServiceUnavailable, "Service Unavailable"),
	http.StatusGatewayTimeout:       newStaticError(http.StatusGatewayTimeout, "Gateway Timeout"),
	http.StatusHTTPVersionNotSupported: newStaticError(http.StatusHTTPVersionNotSupported, "HTTP Version Not Supported"),
}

// sendError streams httpError as JSON directly to rw. If details are nil and the error
// is static, it serves the pre-marshaled payload with a precomputed Content-Length header.
func sendError(rw http.ResponseWriter, httpError HttpError) {
	rw.Header().Set(HeaderContentType, ContentTypeProblemJson)

	if httpError.Details == nil {
		if staticResp, ok := staticErrors[httpError.Status]; ok && staticResp.title == httpError.Title {
			rw.Header().Set("Content-Length", staticResp.contentLength)
			rw.WriteHeader(httpError.Status)
			_, _ = rw.Write(staticResp.payload)
			return
		}
	}

	b, err := json.Marshal(httpError)
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to encode error response as JSON %s", err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	b = append(b, '\n')
	rw.Header().Set("Content-Length", strconv.Itoa(len(b)))
	rw.WriteHeader(httpError.Status)
	_, _ = rw.Write(b)
}

// 400 Bad Request
func SendBadRequest(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Bad Request",
		Status:  http.StatusBadRequest,
		Details: details,
	})
}

func SendValidationError(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Bad Request",
		Status:  http.StatusBadRequest,
		Details: details,
	})
}

// 401 Unauthorized
func SendUnauthorized(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Unauthorized",
		Status:  http.StatusUnauthorized,
		Details: details,
	})
}

// 403 Forbidden
func SendForbidden(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Forbidden",
		Status:  http.StatusForbidden,
		Details: details,
	})
}

// 404 Not Found
func SendNotFound(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Not Found",
		Status:  http.StatusNotFound,
		Details: details,
	})
}

// 405 Method Not Allowed
func SendMethodNotAllowed(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Method Not Allowed",
		Status:  http.StatusMethodNotAllowed,
		Details: details,
	})
}

// 406 Not Acceptable
func SendNotAcceptable(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Not Acceptable",
		Status:  http.StatusNotAcceptable,
		Details: details,
	})
}

// 408 Request Timeout
func SendRequestTimeout(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Request Timeout",
		Status:  http.StatusRequestTimeout,
		Details: details,
	})
}

// 409 Conflict
func SendConflict(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Conflict",
		Status:  http.StatusConflict,
		Details: details,
	})
}

// 410 Gone
func SendGone(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Gone",
		Status:  http.StatusGone,
		Details: details,
	})
}

// 411 Length Required
func SendLengthRequired(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Length Required",
		Status:  http.StatusLengthRequired,
		Details: details,
	})
}

// 412 Precondition Failed
func SendPreconditionFailed(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Precondition Failed",
		Status:  http.StatusPreconditionFailed,
		Details: details,
	})
}

// 413 Payload Too Large
func SendPayloadTooLarge(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Payload Too Large",
		Status:  413,
		Details: details,
	})
}

// 414 URI Too Long
func SendURITooLong(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "URI Too Long",
		Status:  414,
		Details: details,
	})
}

// 415 Unsupported Media Type
func SendUnsupportedMediaType(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Unsupported Media Type",
		Status:  http.StatusUnsupportedMediaType,
		Details: details,
	})
}

// 416 Range Not Satisfiable
func SendRangeNotSatisfiable(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Range Not Satisfiable",
		Status:  416,
		Details: details,
	})
}

// 417 Expectation Failed
func SendExpectationFailed(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Expectation Failed",
		Status:  http.StatusExpectationFailed,
		Details: details,
	})
}

// 422 Unprocessable Entity
func SendUnprocessableEntity(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Unprocessable Entity",
		Status:  http.StatusUnprocessableEntity,
		Details: details,
	})
}

// 429 Too Many Requests
func SendTooManyRequests(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Too Many Requests",
		Status:  http.StatusTooManyRequests,
		Details: details,
	})
}

// 500 Internal Server Error
func SendInternalServerError(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Internal Server Error",
		Status:  http.StatusInternalServerError,
		Details: details,
	})
}

// 501 Not Implemented
func SendNotImplemented(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Not Implemented",
		Status:  http.StatusNotImplemented,
		Details: details,
	})
}

// 502 Bad Gateway
func SendBadGateway(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Bad Gateway",
		Status:  http.StatusBadGateway,
		Details: details,
	})
}

// 503 Service Unavailable
func SendServiceUnavailable(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Service Unavailable",
		Status:  http.StatusServiceUnavailable,
		Details: details,
	})
}

// 504 Gateway Timeout
func SendGatewayTimeout(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "Gateway Timeout",
		Status:  http.StatusGatewayTimeout,
		Details: details,
	})
}

// 505 HTTP Version Not Supported
func SendHTTPVersionNotSupported(rw http.ResponseWriter, details ErrorDetails) {
	sendError(rw, HttpError{
		Title:   "HTTP Version Not Supported",
		Status:  http.StatusHTTPVersionNotSupported,
		Details: details,
	})
}
