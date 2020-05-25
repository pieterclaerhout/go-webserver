package respond

import (
	"net/http"

	"github.com/pieterclaerhout/go-log"
)

// Response is used for a HTTP response
type Response struct {
	Body       interface{}
	NewURL     string `json:"-"`
	StatusCode int
}

// ErrorResponse is used for a HTTP error response
type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

// OK is used to send a HTTP 200 response
func OK(body interface{}) *Response {
	return &Response{
		Body:       body,
		StatusCode: http.StatusOK,
	}
}

// Redirect performs a server side redirect
func Redirect(newURL string) *Response {
	return &Response{
		NewURL:     newURL,
		StatusCode: http.StatusMovedPermanently,
	}
}

// NotFound is used to send a 404 error
func NotFound(message string) *Response {
	return ErrorWithCode(message, http.StatusNotFound)
}

// MethodNotAllowed is used to send a 405 error
func MethodNotAllowed(message string) *Response {
	return ErrorWithCode(message, http.StatusMethodNotAllowed)
}

// Error is used to send a HTTP response with a generic status code
func Error(err error) *Response {
	log.StackTrace(err)
	return ErrorWithCode(err.Error(), http.StatusInternalServerError)
}

// ErrorWithCode is used to send a HTTP response with a custom status code
func ErrorWithCode(message string, statusCode int) *Response {
	return &Response{
		Body: ErrorResponse{
			Error: message,
		},
		StatusCode: statusCode,
	}
}
