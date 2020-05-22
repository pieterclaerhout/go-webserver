# go-webserver

[![Go Report Card](https://goreportcard.com/badge/github.com/pieterclaerhout/go-webserver)](https://goreportcard.com/report/github.com/pieterclaerhout/go-webserver)
[![Documentation](https://godoc.org/github.com/pieterclaerhout/go-webserver?status.svg)](http://godoc.org/github.com/pieterclaerhout/go-webserver)
[![license](https://img.shields.io/badge/license-Apache%20v2-orange.svg)](https://github.com/pieterclaerhout/go-webserver/raw/master/LICENSE)
[![GitHub version](https://badge.fury.io/gh/pieterclaerhout%2Fgo-webserver.svg)](https://badge.fury.io/gh/pieterclaerhout%2Fgo-webserver)
[![GitHub issues](https://img.shields.io/github/issues/pieterclaerhout/go-webserver.svg)](https://github.com/pieterclaerhout/go-webserver/issues)

This is the basic structure I use when creating a webserver.

It's based on the [Chi router](github.com/go-chi/chi) and it tries be as minimal as possible.

## Core Concepts

The core concept is that you compose your web app based on different App instances.

An app is a self-contained piece of functionality you expose via the webserver.

It needs to implement the following interface:

```go
type App interface {
	Name() string
	Register(r *chi.Mux)
}
```

## Getting started

To get started, you first need to create a struct implementing the `App` interface:

```go
package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pieterclaerhout/go-log/versioninfo"
	"github.com/pieterclaerhout/go-webserver/respond"
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
```

Then, you can run this app using the following boilerplate code:

```go
package main

import (
	"os"

	"github.com/pieterclaerhout/go-log"
	webserver "github.com/pieterclaerhout/go-webserver"
)

func main() {

	// Setup logging
	log.PrintColors = true
	log.PrintTimestamp = true
	log.DebugMode = (os.Getenv("DEBUG") == "1")

	// Run the app with the server
	err := webserver.New().RunWithApps(
		&SampleApp{},
	)
	log.CheckError(err)

}
```

You can easily create multiple apps and run them using the same server.

## Features

### Easy way to register dependencies

When you e.g. need a database connection in your web application, you can make define it in your `App` struct and re-use it through the complete web application.

### Defining handlers

To define handlers, it's best to return them as [http.HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc) instances:

```go
func (a *SampleApp) handleJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.sampleResponse().ToJSON(w)
	}
}
```

This makes it easy to define prerequisites which only need to be defined once:

```go
func (a *SampleApp) handleJSON() http.HandlerFunc {
    prerequisite := preparePrerequisite()
	return func(w http.ResponseWriter, r *http.Request) {
		// use prerequisite
	}
}
```

The advantage is that `preparePrerequisite` is executed only once.

The same for defining e.g. response types:

```go
func (a *SampleApp) handleJSON() http.HandlerFunc {
   
    type response struct {
        FieldA string
        FieldB string
    }

	return func(w http.ResponseWriter, r *http.Request) {
		resp := &response{}
	}
}
```

This keeps the package space decluttered and makes naming things a lot easier.

### Elegantly create responses

By using the `respond` package, you can create and send responses in a very easy and fluent way.

Creating a response can be done with one of the following functions:

* `respond.OK(body interface{}) *Response` -> HTTP 200 response
* `respond.Redirect(newURL string) *Response` -> HTTP 301 response
* `respond.NotFound(message string) *Response` -> HTTP 404 response
* `respond.MethodNotAllowed(message string) *Response` -> HTTP 405 response
* `respond.Error(err error) *Response` -> HTTP 500 response
* `responsd.ErrorWithCode(message string, statusCode int) *Response` -> HTTP response with custom error code

Once you created a response, you can decide how to output it:

* `(resp Response) ToText(w http.ResponseWriter)` -> Writes out the body as plain text
* `(resp Response) ToHTML(w http.ResponseWriter)` -> Writes out the body as HTML
* `(resp Response) ToJSON(w http.ResponseWriter)` -> Writes out the body as JSON

If you want to use auto-negotation for these 3 types, you can use the `Write` function:

* `(resp Response) Write(w http.ResponseWriter, r *http.Request)`

## Health endpoint

By default, each app has an endpoint `/status` which can be used for health checks.

The URL to which the health endpoint needs to be exposed can be customized by the `HealthEndpoint` property of the server instance.

## `robots.txt` support

By default, each app has an endpoint `/robots.txt` which outputs the following `robots.txt` file:

```
User-agent: *
Disallow: /
```

You can change it by altering the `RobotsTxt` property of the server.

## Error recovery

By default, error recovery is enabled and uses the `respond.Error` function to log the error.

## Logging

Logging is enabled as well and outputs in the following format:

```
2020-05-22 15:27:32.349 | INFO  |        383.724µs |        127.0.0.1 "GET /json HTTP/1.1" 200 136 "Paw/3.1.10 (Macintosh; OS X/10.15.4) GCDHTTPRequest"
2020-05-22 15:27:37.893 | INFO  |        146.911µs |        127.0.0.1 "GET /html HTTP/1.1" 200 184 "Paw/3.1.10 (Macintosh; OS X/10.15.4) GCDHTTPRequest"
2020-05-22 15:27:40.853 | INFO  |         93.863µs |        127.0.0.1 "GET /auto HTTP/1.1" 200 136 "Paw/3.1.10 (Macintosh; OS X/10.15.4) GCDHTTPRequest"
2020-05-22 15:27:44.466 | WARN  |         89.552µs |        127.0.0.1 "GET /invalid HTTP/1.1" 404 25 "Paw/3.1.10 (Macintosh; OS X/10.15.4) GCDHTTPRequest"
```