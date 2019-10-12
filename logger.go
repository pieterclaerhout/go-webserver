package webserver

import (
	"fmt"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pieterclaerhout/go-log"
)

// LoggerConfig defines the logger config
type LoggerConfig struct {
	Skipper middleware.Skipper // The skipper function to use
}

// Logger returns the default logger
func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(LoggerConfig{})
}

// LoggerWithConfig returns a Logger middleware using go-log
func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {

	if config.Skipper == nil {
		config.Skipper = middleware.DefaultLoggerConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			if config.Skipper(c) {
				return next(c)
			}

			if err = next(c); err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			n := res.Status

			logFunction := log.Info
			switch {
			case n >= 500:
				logFunction = log.Error
			case n >= 400:
				logFunction = log.Warn
			case n >= 300:
				logFunction = log.Warn
			}

			var buf strings.Builder
			buf.WriteString(req.Method)
			buf.WriteString(" ")
			buf.WriteString(fmt.Sprintf("%d", n))
			buf.WriteString(" ")
			buf.WriteString(req.RequestURI)
			buf.WriteString("\n")
			logFunction(buf.String())

			return nil

		}
	}

}
