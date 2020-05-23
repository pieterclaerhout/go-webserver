package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pieterclaerhout/go-log/versioninfo"
	"github.com/pieterclaerhout/go-webserver/v2/binder"
	"github.com/pieterclaerhout/go-webserver/v2/respond"
)

// SampleApp defines a sample web application
type SampleApp struct {
}

// Name returns the name of this app
func (a *SampleApp) Name() string {
	return "sample"
}

// Register registers the routes for this app
func (a *SampleApp) Register(r *chi.Mux) {

	r.Get("/json", a.handleJSON())
	r.Get("/html", a.handleHTML())
	r.Get("/text", a.handleText())
	r.Get("/auto", a.handleAuto())
	r.Get("/parse", a.handleParse())
	r.Post("/parse", a.handleParse())

	r.NotFound(a.handleNotFound())
	r.MethodNotAllowed(a.handleMethodNotAllowed())

}

func (a *SampleApp) handleJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.sampleResponse().ToJSON(w)
	}
}

func (a *SampleApp) handleHTML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.sampleResponse().ToHTML(w)
	}
}

func (a *SampleApp) handleText() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.sampleResponse().ToText(w)
	}
}

func (a *SampleApp) handleAuto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.sampleResponse().Write(w, r)
	}
}

func (a *SampleApp) handleParse() http.HandlerFunc {

	type request struct {
		Value string `json:"value" form:"value"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		if err := binder.Bind(r, &req); err != nil {
			respond.Error(err).Write(w, r)
			return
		}

		respond.OK(req).ToJSON(w)

	}

}

func (a *SampleApp) handleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond.NotFound("nothing here").Write(w, r)
	}
}

func (a *SampleApp) handleMethodNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond.MethodNotAllowed(r.Method+" is not supported").Write(w, r)
	}
}

func (a *SampleApp) sampleResponse() *respond.Response {

	type response struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Copyright   string `json:"copyright"`
		Version     string `json:"version"`
		Revision    string `json:"revision"`
		Branch      string `json:"branch"`
	}

	return respond.OK(
		response{
			Name:        versioninfo.ProjectName,
			Description: versioninfo.ProjectDescription,
			Copyright:   versioninfo.ProjectCopyright,
			Version:     versioninfo.Version,
			Revision:    versioninfo.Revision,
			Branch:      versioninfo.Branch,
		},
	)

}
