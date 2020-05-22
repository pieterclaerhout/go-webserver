package webserver

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pieterclaerhout/go-log"
)

func (s *Server) printRoutes() {

	var packageName = func(v interface{}) string {
		if v == nil {
			return ""
		}

		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Ptr {
			return val.Elem().Type().PkgPath()
		}
		return val.Type().PkgPath()
	}

	mainPackage := filepath.Dir(packageName(s))

	routes := []string{}

	chi.Walk(s.router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {

		route = strings.Replace(route, "/*/", "/", -1)

		handerName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		handerName = strings.ReplaceAll(handerName, mainPackage+"/", "")
		handerName = strings.TrimSuffix(handerName, ".func1")

		routes = append(routes, fmt.Sprintf("%-7s %-30s %s", method, route, handerName))

		return nil

	})

	sort.Strings(routes)

	for _, route := range routes {
		log.Debug(route)
	}

}
