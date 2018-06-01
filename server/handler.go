package server

import (
	"net/http"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/Peripli/service-manager/pkg/plugin"
	"github.com/Peripli/service-manager/rest"
	"github.com/gorilla/mux"
)

func newHttpHandler(filters []plugin.Filter, handler plugin.Handler) http.Handler {
	return httpHandler(chain(filters, handler))
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

	var body interface{}
	if request.Method == "PUT" || request.Method == "POST" {
		if err := rest.ReadJSONBody(request, &body); err != nil {
			return nil, err
		}
	}

	return &plugin.Request{
		Request:    request,
		PathParams: pathParams,
		Body:       body,
	}, nil
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
	logrus.Debugf("%d filters for endpoint %v", len(matches), endpoint)
	return matches
}

func matchPath(endpointPath string, pattern string) bool {
	if pattern == "" {
		return true
	}
	match, err := path.Match(pattern, endpointPath)
	if err != nil {
		logrus.Fatalf("Invalid endpoint path pattern %s: %v", endpointPath, err)
	}
	return match
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
