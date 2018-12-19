package configuration

import (
	"encoding/json"
	"github.com/cloud-ca/go-cloudca/api"
)

type ServiceConnection struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	ServiceCode string `json:"serviceCode,omitempty"`
}

type ServiceConnectionService interface {
	Get(id string) (*ServiceConnection, error)
	List() ([]ServiceConnection, error)
	ListWithOptions(options map[string]string) ([]ServiceConnection, error)
}

type ServiceConnectionApi struct {
	configurationService ConfigurationService
}

func NewServiceConnectionService(apiClient api.ApiClient) ServiceConnectionService {
	return &ServiceConnectionApi{
		configurationService: NewConfigurationService(apiClient, "services/connections"),
	}
}

func parseServiceConnection(data []byte) *ServiceConnection {
	service_connection := ServiceConnection{}
	json.Unmarshal(data, &service_connection)
	return &service_connection
}

func parseServiceConnectionList(data []byte) []ServiceConnection {
	service_connections := []ServiceConnection{}
	json.Unmarshal(data, &service_connections)
	return service_connections
}

func (serviceConnectionApi *ServiceConnectionApi) Get(id string) (*ServiceConnection, error) {
	data, err := serviceConnectionApi.configurationService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseServiceConnection(data), nil
}

//List all service connections
func (serviceConnectionApi *ServiceConnectionApi) List() ([]ServiceConnection, error) {
	return serviceConnectionApi.ListWithOptions(map[string]string{})
}

//List all service connections. Can use options to do sorting and paging.
func (serviceConnectionApi *ServiceConnectionApi) ListWithOptions(options map[string]string) ([]ServiceConnection, error) {
	data, err := serviceConnectionApi.configurationService.List(options)
	if err != nil {
		return nil, err
	}
	return parseServiceConnectionList(data), nil
}
