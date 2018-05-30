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

// Package rest contains logic for building the Service Manager REST API
package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Peripli/service-manager/types"
)

// AllMethods matches all REST HTTP Methods
const AllMethods = "*"

// API is the primary interface for REST API registration
type API interface {
	// Controllers returns the controllers registered for the API
	Controllers() []Controller

	// RegisterControllers registers a set of controllers
	RegisterControllers(...Controller)
}

// Controller is an entity that wraps a set of HTTP Routes
type Controller interface {
	// Routes returns the set of routes for this controller
	Routes() []Route
}

// Route is a mapping between an Endpoint and a REST API Handler
type Route struct {
	// Endpoint is the combination of Path and HTTP Method for the specified route
	Endpoint Endpoint

	// Handler is the function that should handle incoming requests for this endpoint
	Handler http.Handler
}

// Endpoint is a combination of a Path and an HTTP Method
type Endpoint struct {
	Path, Method string
}

type OSBPlugin interface {
	Provisioner
	DeProvisioner
}

type Provisioner interface {
	Provision(http.ResponseWriter, *http.Request) error
}

type DeProvisioner interface {
	DeProvision(http.ResponseWriter, *http.Request) error
}

type OSBFilter struct {
	PluginHandler APIHandler
	Typee         string
}

func (f OSBFilter) Handler() APIHandler {
	return f.pluginHandler
}

func (f OSBFilter) Type() string {
	return f.Typee
}

func RegisterPlugin(plugin interface{}) {
	switch p := plugin.(type) {
	case Provisioner:
		// api.RegisterFilter()
		filter := OSBFilter{p.OnProvision, "inbound"}
		filters.Path("/v1/osb").Filter(filter)
	case DeProvisioner:
		filter := OSBFilter{p.OnDeprovision, "inbound"}
		filters.Path("/v1/osb").Filter(filter)

	}
}

type AuthFilter struct{}

func (f AuthFilter) Handler() APIHandler {
	return func(rw http.ResponseWriter, req *http.Request) error {
		fmt.Println("auth filter!")
		return types.NewErrorResponse(errors.New("asdf"), 401, "Unauthorized")
	}
}

func (f AuthFilter) Type() string {
	return "inbound"
}

type OrgFilter struct{}

func (f OrgFilter) Handler() APIHandler {
	return func(rw http.ResponseWriter, req *http.Request) error {
		fmt.Println("org filter!")
		return nil
	}
}

func (f OrgFilter) Type() string {
	return "outbound"
}

type Filters []Filter

func (filters Filters) Wrap(handler http.Handler) http.Handler {
	inboundFilters, outboundFilters := filters.separate()
	fmt.Printf("Filters separated! Inbound: %v, Outbound: %v\n", len(inboundFilters), len(outboundFilters))
	return APIHandler(func(rw http.ResponseWriter, req *http.Request) error {
		if err := filters.process(inboundFilters, rw, req); err != nil {
			return err
		}

		apiHandler, ok := handler.(APIHandler)
		if ok {
			if err := apiHandler(rw, req); err != nil {
				return err
			}
		} else {
			handler.ServeHTTP(rw, req)
		}

		return filters.process(outboundFilters, rw, req)
	})
}

func (filters Filters) separate() (Filters, Filters) {
	var inbound = make(Filters, 0)
	var outbound = make(Filters, 0)
	for _, filter := range filters {
		if filter.Type() == "inbound" {
			inbound = append(inbound, filter)
		} else {
			outbound = append(outbound, filter)
		}
	}

	return inbound, outbound
}

func (f Filters) process(filters Filters, rw http.ResponseWriter, req *http.Request) error {
	for _, filter := range filters {
		h := filter.Handler()
		if err := h(rw, req); err != nil {
			return err
		}
	}
	return nil
}

// APIHandler enriches http.HandlerFunc with an error response for further processing
type APIHandler func(http.ResponseWriter, *http.Request) error

func (ah APIHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if err := ah(rw, r); err != nil {
		HandleError(err, rw)
	}
}
