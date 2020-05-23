package binder

import (
	"encoding/json"
	"mime"
	"net/http"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/pieterclaerhout/go-log"
)

// Bind binds the request data to the object
//
// It checks the request content type on how the data should be parsed.
//
// Supports both HTTP form as well query string and JSON posts
func Bind(r *http.Request, val interface{}) error {

	// Parse as a JSON request
	if r.Method != http.MethodGet && isContentType(r, "application/json") {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&val); err != nil {
			return err
		}
		return nil
	}

	// Parse HTTP form and query strings
	if err := r.ParseForm(); err != nil {
		log.StackTrace(err)
		return err
	}
	return form.NewDecoder().Decode(&val, r.Form)

}

// isContentType checks if the request is the specific content type
func isContentType(r *http.Request, mimetype string) bool {
	contentType := r.Header.Get("Content-type")
	if contentType == "" {
		return mimetype == "application/octet-stream"
	}

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}
