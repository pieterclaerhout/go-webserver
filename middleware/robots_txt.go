package middleware

import (
	"net/http"
	"strings"

	"github.com/pieterclaerhout/go-log"
	"github.com/pieterclaerhout/go-webserver/v2/respond"
)

const robotsTxtPath = "/robots.txt"

// RobotsTxt outputs the default robots.txt file
func RobotsTxt(robotsTxt string) func(http.Handler) http.Handler {
	log.Debug("Registering robots.txt")
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" && strings.EqualFold(r.URL.Path, robotsTxtPath) {
				respond.OK(robotsTxt).ToText(w)
				return
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
