package cca

import (
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/configuration"
	"github.com/cloud-ca/go-cloudca/services"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
)

const (
	DEFAULT_API_URL = "https://api.cloud.ca/v1/"
)

type CcaClient struct {
	apiClient          api.ApiClient
	Tasks              services.TaskService
	Environments       configuration.EnvironmentService
	Users              configuration.UserService
	ServiceConnections configuration.ServiceConnectionService
	Organizations      configuration.OrganizationService
}

//Create a CcaClient with the default URL
func NewCcaClient(apiKey string) *CcaClient {
	return NewCcaClientWithURL(DEFAULT_API_URL, apiKey)
}

//Create a CcaClient with a custom URL
func NewCcaClientWithURL(apiURL string, apiKey string) *CcaClient {
	apiClient := api.NewApiClient(apiURL, apiKey)
	return NewCcaClientWithApiClient(apiClient)
}

//Create a CcaClient with a custom URL that accepts insecure connections
func NewInsecureCcaClientWithURL(apiURL string, apiKey string) *CcaClient {
	apiClient := api.NewInsecureApiClient(apiURL, apiKey)
	return NewCcaClientWithApiClient(apiClient)
}

func NewCcaClientWithApiClient(apiClient api.ApiClient) *CcaClient {
	ccaClient := CcaClient{
		apiClient:          apiClient,
		Tasks:              services.NewTaskService(apiClient),
		Environments:       configuration.NewEnvironmentService(apiClient),
		Users:              configuration.NewUserService(apiClient),
		ServiceConnections: configuration.NewServiceConnectionService(apiClient),
		Organizations:      configuration.NewOrganizationService(apiClient),
	}
	return &ccaClient
}

//Get the Resources for a specific serviceCode and environmentName
//For now it assumes that the serviceCode belongs to a cloud.ca service type
func (c CcaClient) GetResources(serviceCode string, environmentName string) (services.ServiceResources, error) {
	//TODO: change to check service type of service code
	return cloudca.NewResources(c.apiClient, serviceCode, environmentName), nil
}

//Get the API url used to do he calls
func (c CcaClient) GetApiURL() string {
	return c.apiClient.GetApiURL()
}

//Get the API key used in the calls
func (c CcaClient) GetApiKey() string {
	return c.apiClient.GetApiKey()
}

//Get the API Client used by all the services
func (c CcaClient) GetApiClient() api.ApiClient {
	return c.apiClient
}
