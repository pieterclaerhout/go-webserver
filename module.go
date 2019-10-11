package webserver

import (
	"github.com/labstack/echo"
)

// Module defines a server module
type Module interface {
	Register(router *echo.Echo) // The function which registers the endpoints on the router
	Start()                     // Executed when the server starts
	Stop()                      // Executed when the server stops
}
