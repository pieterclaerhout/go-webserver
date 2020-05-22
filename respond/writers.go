package respond

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Write write the response in the appropriate format based on the request header
func (resp Response) Write(w http.ResponseWriter, r *http.Request) error {

	if resp.StatusCode == http.StatusMovedPermanently {
		http.Redirect(w, r, resp.NewURL, resp.StatusCode)
		return nil
	}

	accepts := parseAccept(r.Header.Get("Accept"))
	for _, accept := range accepts {
		switch accept.SubType {
		case "html", "xhtml", "xhtml+xml":
			return resp.ToHTML(w)
		case "text":
			return resp.ToText(w)
		}
	}

	return resp.ToJSON(w)

}

// ToText returns the response as plain text
func (resp Response) ToText(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(resp.StatusCode)
	_, err := w.Write([]byte(fmt.Sprintf("%v", resp.Body)))
	return err
}

// ToHTML returns the response as HTML
func (resp Response) ToHTML(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(resp.StatusCode)
	w.Write([]byte("<pre>"))
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	enc.SetIndent("", "    ")
	if err := enc.Encode(resp.Body); err != nil {
		return err
	}
	w.Write([]byte("</pre>"))
	return nil
}

// ToJSON returns the response as JSON
func (resp Response) ToJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	return json.NewEncoder(w).Encode(resp.Body)
}
