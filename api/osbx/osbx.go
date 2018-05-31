package osbx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/Peripli/service-manager/pkg/plugin"
	"github.com/Peripli/service-manager/rest"
)

func NewController() *controller {
	return &controller{os.Getenv("BROKER_URL")}
}

type controller struct {
	brokerURL string
}

const (
	baseURL            = "/v1/osb/{broker_id}"
	catalogURL         = baseURL + "/v2/catalog"
	serviceInstanceURL = baseURL + "/v2/service_instances/{instance_id}"
	serviceBindingURL  = baseURL + "/v2/service_instances/{instance_id}/service_bindings/{binding_id}"
)

func (c *controller) Routes() []rest.Route {
	// NOTE: for filters/plugins to work, create a separate route for each mothod
	return []rest.Route{
		{rest.Endpoint{"GET", catalogURL}, c.osbHandler},
		{rest.Endpoint{"PUT", serviceInstanceURL}, c.osbHandler},
		{rest.Endpoint{"DELETE", serviceInstanceURL}, c.osbHandler},
		{rest.Endpoint{"PUT", serviceBindingURL}, c.osbHandler},
		{rest.Endpoint{"DELETE", serviceBindingURL}, c.osbHandler},
	}
}

// osbHandler forwards request to actual broker
func (c *controller) osbHandler(req *plugin.Request) (*plugin.Response, error) {
	reqBody, err := json.Marshal(req.Body)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	osbURL := strings.TrimPrefix(req.URL.RequestURI(),
		"/v1/osb/"+req.PathParams["broker_id"])
	url := c.brokerURL + osbURL
	logrus.Debugf("Forwarding: %s %s", req.Method, url)

	request, err := http.NewRequest(req.Method, url, bytes.NewReader(reqBody))
	for k, v := range req.Header {
		if k != "Content-Length" && k != "Host" {
			request.Header[k] = v
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("Broker returned status %d", response.StatusCode)
	}
	var resBody interface{}
	err = json.NewDecoder(response.Body).Decode(&resBody)
	if err != nil {
		return nil, err
	}
	return &plugin.Response{
		Header:     response.Header,
		Body:       resBody,
		StatusCode: response.StatusCode,
	}, nil
}
