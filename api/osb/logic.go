package osb

import (
	"net/http"

	//"github.com/pmorie/osb-broker-lib/pkg/broker"

	"fmt"

	"github.com/Peripli/service-manager/storage"
	"github.com/Peripli/service-manager/types"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	osbc "github.com/pmorie/go-open-service-broker-client/v2"
	"github.com/pmorie/osb-broker-lib/pkg/broker"
)

// BusinessLogic provides an implementation of the osb.BusinessLogic interface.
type BusinessLogic struct {
	createFunc    osbc.CreateFunc
	brokerStorage storage.Broker
}

var _ broker.Interface = &BusinessLogic{}

func NewBusinessLogic(createFunc osbc.CreateFunc, brokerStorage storage.Broker) *BusinessLogic {
	return &BusinessLogic{
		createFunc:    createFunc,
		brokerStorage: brokerStorage,
	}
}

func (b *BusinessLogic) GetCatalog(c *broker.RequestContext) (*broker.CatalogResponse, error) {
	client, err := b.osbClient(c.Request)
	if err != nil {
		return nil, err
	}
	response, err := client.GetCatalog()
	if err != nil {
		return nil, err
	}

	return &broker.CatalogResponse{
		CatalogResponse: *response,
	}, nil
}

func (b *BusinessLogic) Provision(request *osbc.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
	client, err := b.osbClient(c.Request)
	if err != nil {
		return nil, err
	}

	response, err := client.ProvisionInstance(request)
	if err != nil {
		return nil, err
	}

	return &broker.ProvisionResponse{
		ProvisionResponse: *response,
	}, nil
}

func (b *BusinessLogic) Deprovision(request *osbc.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
	client, err := b.osbClient(c.Request)
	if err != nil {
		return nil, err
	}
	response, err := client.DeprovisionInstance(request)
	if err != nil {
		return nil, err
	}

	return &broker.DeprovisionResponse{
		DeprovisionResponse: *response,
	}, nil
}

func (b *BusinessLogic) LastOperation(request *osbc.LastOperationRequest, c *broker.RequestContext) (*broker.LastOperationResponse, error) {
	client, err := b.osbClient(c.Request)
	if err != nil {
		return nil, err
	}
	response, err := client.PollLastOperation(request)
	if err != nil {
		return nil, err
	}

	return &broker.LastOperationResponse{
		LastOperationResponse: *response,
	}, nil
}

func (b *BusinessLogic) Bind(request *osbc.BindRequest, c *broker.RequestContext) (*broker.BindResponse, error) {
	client, err := b.osbClient(c.Request)
	if err != nil {
		return nil, err
	}

	response, err := client.Bind(request)
	if err != nil {
		return nil, err
	}

	return &broker.BindResponse{
		BindResponse: *response,
	}, nil

}

func (b *BusinessLogic) Unbind(request *osbc.UnbindRequest, c *broker.RequestContext) (*broker.UnbindResponse, error) {
	client, err := b.osbClient(c.Request)
	if err != nil {
		return nil, err
	}

	response, err := client.Unbind(request)
	if err != nil {
		return nil, err
	}

	return &broker.UnbindResponse{
		UnbindResponse: *response,
	}, nil
}

func (b *BusinessLogic) Update(request *osbc.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
	client, err := b.osbClient(c.Request)
	if err != nil {
		return nil, err
	}

	response, err := client.UpdateInstance(request)
	if err != nil {
		return nil, err
	}

	return &broker.UpdateInstanceResponse{
		UpdateInstanceResponse: *response,
	}, nil
}

func (b *BusinessLogic) ValidateBrokerAPIVersion(version string) error {
	expectedVersion := osbc.LatestAPIVersion().HeaderValue()
	if version != expectedVersion {
		return fmt.Errorf("error validating OSB Version: expected %s but was %s", expectedVersion, version)
	}
	return nil
}

func clientConfigForBroker(broker *types.Broker) *osbc.ClientConfiguration {
	config := osbc.DefaultClientConfiguration()
	config.Name = broker.Name
	config.URL = broker.URL
	config.AuthConfig = &osbc.AuthConfig{
		BasicAuthConfig: &osbc.BasicAuthConfig{
			Username: broker.User,
			Password: broker.Password,
		},
	}
	return config
}

func (b *BusinessLogic) osbClient(request *http.Request) (osbc.Client, error) {
	vars := mux.Vars(request)
	brokerID, ok := vars["brokerID"]
	if !ok {
		return nil, fmt.Errorf("Error creating OSB client: brokerID path parameter not found")
	}
	broker, err := b.brokerStorage.Find(nil, brokerID)
	if err != nil {
		return nil, fmt.Errorf("Error obtaining broker with id %s from storage: %s", brokerID, err)
	}
	config := clientConfigForBroker(broker)
	logrus.Debug("Building OSB client for broker with name: ", config.Name, " accesible at: ", config.URL)
	return b.createFunc(config)
}