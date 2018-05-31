/*
 * Copyright 2018 The Service Manager Authors
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

// Package server contains the logic of the Service Manager server
package server

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"fmt"

	"github.com/Peripli/service-manager/pkg/plugin"
	"github.com/Peripli/service-manager/rest"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server is the server to process incoming HTTP requests
type Server struct {
	Configuration *Config
	Router        *mux.Router
	api           rest.API
	Filters       []plugin.Filter
}

// New creates a new server with the provided REST API configuration and server configuration
// Returns the new server and an error if creation was not successful
func New(api rest.API, config *Config) (*Server, error) {
	router := mux.NewRouter().StrictSlash(true)

	return &Server{
		Configuration: config,
		Router:        router,
		api:           api,
	}, nil
}

func (server *Server) RegisterFilter(filter plugin.Filter) {
	server.Filters = append(server.Filters, filter)
}

// Run starts the server awaiting for incoming requests
func (server *Server) Run(ctx context.Context) {
	if err := registerControllers(server.Router, server.api.Controllers(), server.Filters); err != nil {
		logrus.Fatal(err)
	}

	handler := &http.Server{
		Handler:      server.Router,
		Addr:         server.Configuration.Address,
		WriteTimeout: server.Configuration.RequestTimeout,
		ReadTimeout:  server.Configuration.RequestTimeout,
	}
	startServer(ctx, handler, server.Configuration.ShutdownTimeout)
}

func registerRoutes(prefix string, fromRouter *mux.Router, toRouter *mux.Router) error {
	subRouter := toRouter.PathPrefix(prefix).Subrouter()
	return fromRouter.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {

		path, err := route.GetPathTemplate()
		if err != nil {
			return fmt.Errorf("register routes: %s", err)
		}
		r := subRouter.Handle(path, route.GetHandler())

		methods, err := route.GetMethods()
		if err != nil {
			return fmt.Errorf("register routes: %s", err)

		}
		if len(methods) > 0 {
			r.Methods(methods...)
		}

		logrus.Info("Registering route with methods: ", methods, " and path: ", path, " behind prefix ", prefix)
		return nil
	})
}

func registerControllers(router *mux.Router, controllers []rest.Controller, filters []plugin.Filter) error {
	for _, ctrl := range controllers {
		for _, route := range ctrl.Routes() {
			r := router.Handle(route.Endpoint.Path,
				newHttpHandler(matchFilters(&route.Endpoint, filters), route.Handler))
			if route.Endpoint.Method != rest.AllMethods {
				r.Methods(route.Endpoint.Method)
			}
		}
	}
	return nil
}

type httpHandler plugin.Handler

func (h httpHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if err := h.serve(res, req); err != nil {
		rest.HandleError(err, res)
	}
}

func (h httpHandler) serve(res http.ResponseWriter, req *http.Request) error {
	restReq, err := readRequest(req)
	if err != nil {
		return err
	}

	restRes, err := h(restReq)
	if err != nil {
		return err
	}

	// copy response headers
	for k, v := range restRes.Header {
		if k != "Content-Length" {
			res.Header()[k] = v
		}
	}

	code := restRes.StatusCode
	if code == 0 {
		code = http.StatusOK
	}
	if err := rest.SendJSON(res, code, restRes.Body); err != nil {
		return err
	}
	return nil
}

func readRequest(request *http.Request) (*plugin.Request, error) {
	pathParams := mux.Vars(request)

	queryParams := map[string]string{}
	for k, v := range request.URL.Query() {
		queryParams[k] = v[0]
	}

	var body interface{}
	if request.Method == "PUT" || request.Method == "POST" {
		if err := rest.ReadJSONBody(request, &body); err != nil {
			return nil, err
		}
	}

	return &plugin.Request{
		Request:     request,
		PathParams:  pathParams,
		QueryParams: queryParams,
		Body:        body,
	}, nil
}

func newHttpHandler(filters []plugin.Filter, handler plugin.Handler) http.Handler {
	return httpHandler(chain(filters, handler))
}

func chain(filters []plugin.Filter, handler plugin.Handler) plugin.Handler {
	if len(filters) == 0 {
		return handler
	}
	next := chain(filters[1:], handler)
	f := filters[0].Func
	return func(req *plugin.Request) (*plugin.Response, error) {
		return f(req, next)
	}
}

func matchFilters(endpoint *rest.Endpoint, filters []plugin.Filter) []plugin.Filter {
	matches := []plugin.Filter{}
	for _, filter := range filters {
		if matchPath(endpoint.Path, filter.PathPattern) &&
			matchMethod(endpoint.Method, filter.Methods) {
			matches = append(matches, filter)
		}
	}
	return matches
}

func matchPath(path string, pattern *regexp.Regexp) bool {
	return pattern == nil || pattern.MatchString(path)
}

func matchMethod(method string, methods []string) bool {
	if method == rest.AllMethods || methods == nil {
		return true
	}
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}

func startServer(ctx context.Context, server *http.Server, shutdownTimeout time.Duration) {
	go gracefulShutdown(ctx, server, shutdownTimeout)

	logrus.Debugf("Listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatal(err)
	}
}

func gracefulShutdown(ctx context.Context, server *http.Server, shutdownTimeout time.Duration) {
	<-ctx.Done()

	c, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	logrus.Debugf("Shutdown with timeout: %s", shutdownTimeout)

	if err := server.Shutdown(c); err != nil {
		logrus.Error("Error: ", err)
		if err := server.Close(); err != nil {
			logrus.Error("Error: ", err)
		}
	} else {
		logrus.Debug("Server stopped")
	}
}
