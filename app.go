package webserver

import (
	"github.com/go-chi/chi"
)

// App defines the app server by the webserver
type App interface {
	Name() string
	Register(r *chi.Mux)
}
