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
	"time"

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
func New(api *rest.API, config *Config) *Server {
	router := mux.NewRouter().StrictSlash(true)
	if err := registerControllers(router, api); err != nil {
		logrus.Fatal(err)
	}

	return &Server{
		Configuration: config,
		Router:        router,
	}
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

func registerControllers(router *mux.Router, api *rest.API) error {
	for _, ctrl := range api.Controllers {
		for _, route := range ctrl.Routes() {
			logrus.Debugf("Register endpoint: %s %s", route.Endpoint.Method, route.Endpoint.Path)
			r := router.Handle(route.Endpoint.Path,
				rest.NewHTTPHandler(rest.MatchFilters(&route.Endpoint, api.Filters), route.Handler))
			r.Methods(route.Endpoint.Method)
		}
	}
	return nil
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
