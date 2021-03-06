package middleware

import (
	"net/http"

	"github.com/pieterclaerhout/go-log"
	"github.com/pieterclaerhout/go-webserver/v2/respond"
)

// Recovery returns the Recovery middleware
func Recovery(next http.Handler) http.Handler {
	log.Debug("Registering recovery")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				respond.Error(err.(error)).Write(w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
