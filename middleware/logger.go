package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/pieterclaerhout/go-log"
)

// Logger returns the logging middleware
func Logger(next http.Handler) http.Handler {
	log.Debug("Registering logger")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ri := parseRequestInfo(r)

		m := httpsnoop.CaptureMetrics(next, w, r)

		ri.code = m.Code
		ri.size = m.Written
		ri.duration = m.Duration

		logF := log.Info
		if ri.code == http.StatusInternalServerError {
			logF = log.Error
		} else if ri.code == http.StatusMethodNotAllowed {
			logF = log.Error
		} else if ri.code == http.StatusNotFound {
			logF = log.Warn
		} else if ri.code == http.StatusUnauthorized {
			logF = log.Warn
		}

		logF(ri.String())

	})
}

// requestInfo contains the info about a request
type requestInfo struct {
	ts        time.Time
	proto     string
	method    string
	uri       string
	referer   string
	ipaddr    string
	code      int
	size      int64
	duration  time.Duration
	userAgent string
}

// String returns the request info as a string
func (ri requestInfo) String() string {
	return fmt.Sprintf(
		"%16v | %16s \"%s %s %s\" %d %d \"%s\"",
		ri.duration, ri.ipaddr, ri.method, ri.uri, ri.proto, ri.code, ri.size, ri.userAgent,
	)
}

func parseRequestInfo(r *http.Request) *requestInfo {
	return &requestInfo{
		ts:        time.Now(),
		proto:     r.Proto,
		method:    r.Method,
		uri:       r.URL.String(),
		referer:   r.Header.Get("Referer"),
		userAgent: r.Header.Get("User-Agent"),
		ipaddr:    requestGetRemoteAddress(r),
	}
}

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		return parts[0]
	}
	return hdrRealIP
}
