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

	"github.com/Peripli/service-manager/rest"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server is the server to process incoming HTTP requests
type Server struct {
	Configuration *Config
	Router        *mux.Router
}

// New creates a new server with the provided REST API configuration and server configuration
// Returns the new server and an error if creation was not successful
func New(api rest.API, config *Config) (*Server, error) {
	router := mux.NewRouter().StrictSlash(true)
	filters := make(map[string]rest.Filters)


	filters["/v1/osb"] = rest.Filters
	filters[".*"] = rest.Filters{rest.AuthFilter{}, rest.OrgFilter{}}

	if err := registerControllers(router, api.Controllers(), filters); err != nil {
		return nil, fmt.Errorf("new Config: %s", err)
	}

	return &Server{
		Configuration: config,
		Router:        router,
	}, nil
}

// Run starts the server awaiting for incoming requests
func (server *Server) Run(ctx context.Context) {
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

func registerControllers(router *mux.Router, controllers []rest.Controller, filters map[string]rest.Filters) error {
	for _, ctrl := range controllers {
		for _, route := range ctrl.Routes() {
			fromRouter, ok := route.Handler.(*mux.Router)
			if ok {
				if err := registerRoutes(route.Endpoint.Path, fromRouter, router); err != nil {
					return fmt.Errorf("register controllers: %s", err)
				}
			} else {
				var r *mux.Route
				matchedFilters, err := getAllMatchingFilters(route.Endpoint.Path, filters)
				if err != nil {
					return fmt.Errorf("problem occured while matching filters: %s", err)
				}
				if len(matchedFilters) > 0 {
					fmt.Printf("filters: %v\n", len(matchedFilters))
					filteredHandler := matchedFilters.Wrap(route.Handler)
					r = router.Handle(route.Endpoint.Path, filteredHandler)
				} else {
					r = router.Handle(route.Endpoint.Path, route.Handler)
				}

				if route.Endpoint.Method != rest.AllMethods {
					r.Methods(route.Endpoint.Method)
				}
			}
		}
	}
	return nil
}

func getAllMatchingFilters(path string, filters map[string]rest.Filters) (rest.Filters, error) {
	result := make(rest.Filters, 0)
	for filterPath, filtersPerRoute := range filters {
		match, err := regexp.MatchString(filterPath, path)
		if err != nil {
			return nil, err
		}
		if match {
			logrus.Infof("matching filter pattern %s to endpoint %s", filterPath, path)
			for _, filter := range filtersPerRoute {
				result = append(result, filter)
			}
		} else {
			logrus.Infof("No match for pattern %s to endpoint %s", filterPath, path)
		}
	}
	return result, nil
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
