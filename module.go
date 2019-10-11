package webserver

import (
	"github.com/labstack/echo"
)

// Module defines a server module
type Module interface {
	Start()                     // Executed when the server starts
	Register(router *echo.Echo) // The function which registers the endpoints on the router
	Stop()                      // Executed when the server stops
}
