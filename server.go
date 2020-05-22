package webserver

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pieterclaerhout/go-log"
	"github.com/pieterclaerhout/go-webserver/middleware"
)

const port = ":8080"

// Server defines the webserver and it's router
type Server struct {
	router *chi.Mux // The router to use
}

// New returns a new app instance
func New() *Server {

	r := chi.NewRouter()

	r.Use(middleware.Health("/status"))
	r.Use(middleware.RobotsTxt())
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)

	return &Server{
		router: r,
	}

}

// RunWithApps starts the server given the apps and returns errors if any
func (s *Server) RunWithApps(apps ...App) error {

	for _, app := range apps {
		log.Info("Registering app:", app.Name())
		app.Register(s.router)
	}

	s.printRoutes()

	log.Info("Server for listening on", port)
	return http.ListenAndServe(port, s.router)

}
