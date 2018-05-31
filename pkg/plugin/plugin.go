package plugin

import (
	"encoding/json"
	"net/http"
	"regexp"
)

type JSON = interface{}
type Object = map[string]interface{}
type Array = []interface{}

type Request struct {
	*http.Request
	PathParams map[string]string
	Body       JSON
}

func (r *Request) String() string {
	return stringify(r)
}

type Response struct {
	// StatusCode is the HTTP status code
	StatusCode int

	Header http.Header

	// Body is the response body parsed as JSON
	Body JSON
}

func (r *Response) String() string {
	return stringify(r)
}

func stringify(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

type Handler func(*Request) (*Response, error)
type FilterFunc func(req *Request, next Handler) (*Response, error)

type RequestMatcher struct {
	// Methods match request method
	// if nil, matches any method
	// NOTE: This will work as long as each route handles a single method.
	// If a route handles multiple methods (e.g. *),
	// the filter could be called for methods which are not listed here.
	Methods []string

	// PathPattern matches endpoint path (as registered in mux)
	// if nil, matches any path
	PathPattern *regexp.Regexp
}

type Filter struct {
	RequestMatcher
	Func FilterFunc
}
