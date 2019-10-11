package webserver

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/pieterclaerhout/go-log"
	"github.com/pieterclaerhout/go-webserver/jobqueue"
)

// Server is an abstraction of a webserver
type Server struct {
	engine              *echo.Echo
	DefaultPort         string
	PrintRoutes         bool
	JobQueuePoolSize    int
	JobQueueConcurrency int
	modules             []Module
}

// New returns a new Server instacce
func New() *Server {

	runtime.GOMAXPROCS(runtime.NumCPU())
	log.PrintTimestamp = true

	return &Server{
		DefaultPort:         ":8080",
		modules:             []Module{},
		JobQueuePoolSize:    10,
		JobQueueConcurrency: 4,
	}

}

// Start starts the webserver on the indicated port
func (server *Server) Start() error {

	jobqueue.Default().Start(server.JobQueuePoolSize, server.JobQueueConcurrency)

	server.engine = echo.New()
	server.engine.HideBanner = true
	server.engine.Debug = log.DebugMode
	server.engine.HTTPErrorHandler = server.handleError

	for _, module := range server.modules {
		module.Register(server.engine)
		module.Start()
	}

	server.registerMiddlewares()

	if server.PrintRoutes {
		server.printRoutes()
	}

	port := server.port()
	return server.engine.Start(port)

}

// Stop stops the server and performs the shutdown action for each module
func (server *Server) Stop() error {

	jobqueue.Default().Stop()

	for _, module := range server.modules {
		module.Register(server.engine)
		module.Stop()
	}

	return server.engine.Close()

}

// Register registers the modules on the main router
func (server *Server) Register(modules ...Module) {
	server.modules = append(server.modules, modules...)
}

// printRoutes prints an overview with all routes
func (server *Server) printRoutes() {

	pkgPath := reflect.TypeOf(*server).PkgPath()
	for _, route := range server.engine.Routes() {
		if route.Name == "github.com/labstack/echo.(*Group).Use.func1" {
			continue
		}
		if route.Method != http.MethodPost && route.Method != http.MethodGet {
			continue
		}
		name := route.Name
		name = strings.ReplaceAll(name, pkgPath, "")
		name = strings.ReplaceAll(name, "-fm", "")
		log.Debug(fmt.Sprintf("%-4s %-30s %s", route.Method, route.Path, name[1:]))
	}

}

// port returns the port on which the server should listen
func (server *Server) port() string {

	port := os.Getenv("PORT")
	if port == "" {
		port = server.DefaultPort
	}

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return port

}

// registerMiddlewares registers the middleware which is going to be used
func (server *Server) registerMiddlewares() {

	server.engine.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	server.engine.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${status} ${uri}\n",
	}))

	server.engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

}
