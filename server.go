package webserver

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pieterclaerhout/go-log"
	"github.com/pieterclaerhout/go-webserver/v2/middleware"
)

const port = ":8080"
const robotsTxt = "User-agent: *\nDisallow: /"
const healthEndpoint = "/status"

// Server defines the webserver and it's router
type Server struct {
	router         *chi.Mux // The router to use
	RobotsTxt      string   // The robots.txt file
	HealthEndpoint string   // The URL to which the health endpoint should be exposed
}

// New returns a new app instance
func New() *Server {

	r := chi.NewRouter()

	return &Server{
		router:         r,
		RobotsTxt:      robotsTxt,
		HealthEndpoint: healthEndpoint,
	}

}

// RunWithApps starts the server given the apps and returns errors if any
func (s *Server) RunWithApps(apps ...App) error {

	s.router.Use(middleware.Health(s.HealthEndpoint))
	s.router.Use(middleware.RobotsTxt(s.RobotsTxt))
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recovery)

	for _, app := range apps {
		log.Info("Registering app:", app.Name())
		app.Register(s.router)
	}

	s.printRoutes()

	log.Info("Server for listening on", port)
	return http.ListenAndServe(port, s.router)

}
