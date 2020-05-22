package middleware

import (
	"net/http"
	"strings"

	"github.com/pieterclaerhout/go-webserver/respond"
)

const robotsTxt = `User-agent: *
Disallow: /`

const robotsTxtPath = "/robots.txt"

// RobotsTxt outputs the default robots.txt file
func RobotsTxt() func(http.Handler) http.Handler {
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
