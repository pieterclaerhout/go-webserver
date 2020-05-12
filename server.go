package webserver

import (
	"fmt"
	stdLog "log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"syscall"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/pieterclaerhout/go-log"
	"github.com/pieterclaerhout/go-webserver/jobqueue"
	wsmiddleware "github.com/pieterclaerhout/go-webserver/middleware"
	"github.com/pieterclaerhout/go-xray"
)

// Server is an abstraction of a webserver
type Server struct {
	engine              *echo.Echo
	DefaultPort         string
	PrintRoutes         bool
	UseJobQueue         bool
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
		UseJobQueue:         true,
		JobQueuePoolSize:    10,
		JobQueueConcurrency: 4,
	}

}

// setupShutdownHook handles the shutdown of the server
func (server *Server) setupShutdownHook() {

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		select {
		case <-signalChan:
			fmt.Println()
			server.Stop()
		}
	}()

}

// Start starts the webserver on the indicated port
func (server *Server) Start() error {
	return server.StartWithListener(nil)
}

// StartWithListener starts the webserver using the given listener
func (server *Server) StartWithListener(listener net.Listener) error {

	server.setupShutdownHook()

	if server.UseJobQueue {
		jobqueue.Default().Start(server.JobQueuePoolSize, server.JobQueueConcurrency)
	}

	server.engine = echo.New()
	if listener != nil {
		server.engine.Listener = listener
		server.engine.HidePort = true
	}
	server.engine.Logger = newLogger()
	server.engine.StdLogger = stdLog.New(server.engine.Logger.Output(), server.engine.Logger.Prefix()+": ", 0)
	server.engine.HideBanner = true
	server.engine.Debug = log.DebugMode
	server.engine.HTTPErrorHandler = server.handleError

	for _, module := range server.modules {
		log.Debug("Starting module:", xray.Name(module))
		module.Register(server.engine)
		module.Start()
	}

	server.registerMiddlewares()

	if server.PrintRoutes {
		server.printRoutes()
	}

	port := server.port()
	log.Debug("Starting webserver on port:", port)

	err := server.engine.Start(port)
	if err == http.ErrServerClosed {
		return nil
	}
	return err

}

// Stop stops the server and performs the shutdown action for each module
func (server *Server) Stop() {

	if server.UseJobQueue {
		log.Debug(("Stopping background task queue"))
		jobqueue.Default().Stop()
	}

	for _, module := range server.modules {
		log.Debug("Stopping module:", xray.Name(module))
		module.Stop()
	}

	log.Debug("Stopping webserver")
	server.engine.Close()

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

	middleware.DefaultLoggerConfig = middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			switch c.Request().RequestURI {
			case "/status", "/favicon.ico":
				return true
			default:
				return false
			}
		},
	}

	server.engine.Use(wsmiddleware.Logger())

	server.engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

}
