package configuration

import (
	"github.com/cloud-ca/go-cloudca/api"
)

type ConfigurationService interface {
	Get(id string, options map[string]string) ([]byte, error)
	List(options map[string]string) ([]byte, error)
	Create(body []byte, options map[string]string) ([]byte, error)
	Update(id string, body []byte, options map[string]string) ([]byte, error)
	Delete(id string, body []byte, options map[string]string) ([]byte, error)
}

//Implementation of the ConfigurationService
type ConfigurationApi struct {
	apiClient         api.ApiClient
	configurationType string
}

func NewConfigurationService(apiClient api.ApiClient, configurationType string) ConfigurationService {
	return &ConfigurationApi{
		apiClient:         apiClient,
		configurationType: configurationType,
	}
}

func (configurationApi *ConfigurationApi) buildEndpoint() string {
	return "/" + configurationApi.configurationType
}

//Get. Returns a []byte (of a json object) that should be unmarshalled to a specific entity
func (configurationApi *ConfigurationApi) Get(id string, options map[string]string) ([]byte, error) {
	request := api.CcaRequest{
		Method:   api.GET,
		Endpoint: configurationApi.buildEndpoint() + "/" + id,
		Options:  options,
	}
	response, err := configurationApi.apiClient.Do(request)
	if err != nil {
		return nil, err
	} else if response.IsError() {
		return nil, api.CcaErrorResponse(*response)
	}
	return response.Data, nil
}

//Get list. Returns a []byte (of a json object) that should be unmarshalled to a specific entity
func (configurationApi *ConfigurationApi) List(options map[string]string) ([]byte, error) {
	request := api.CcaRequest{
		Method:   api.GET,
		Endpoint: configurationApi.buildEndpoint(),
		Options:  options,
	}
	response, err := configurationApi.apiClient.Do(request)
	if err != nil {
		return nil, err
	} else if response.IsError() {
		return nil, api.CcaErrorResponse(*response)
	}
	return response.Data, nil
}

//Create as described in the body parameter (json object). Returns a []byte (of a json object) that should be unmarshalled to a specific entity
func (configurationApi *ConfigurationApi) Create(body []byte, options map[string]string) ([]byte, error) {
	request := api.CcaRequest{
		Method:   api.POST,
		Body:     body,
		Endpoint: configurationApi.buildEndpoint(),
		Options:  options,
	}
	response, err := configurationApi.apiClient.Do(request)
	if err != nil {
		return nil, err
	} else if response.IsError() {
		return nil, api.CcaErrorResponse(*response)
	}
	return response.Data, nil
}

//Update specified id as described in the body parameter (json object). Returns a []byte (of a json object) that should be unmarshalled to a specific entity
func (configurationApi *ConfigurationApi) Update(id string, body []byte, options map[string]string) ([]byte, error) {
	request := api.CcaRequest{
		Method:   api.PUT,
		Body:     body,
		Endpoint: configurationApi.buildEndpoint() + "/" + id,
		Options:  options,
	}
	response, err := configurationApi.apiClient.Do(request)
	if err != nil {
		return nil, err
	} else if response.IsError() {
		return nil, api.CcaErrorResponse(*response)
	}
	return response.Data, nil
}

//Delete. A body (json object) can be provided if some fields must be sent to server. Returns a []byte (of a json object) that should be unmarshalled to a specific entity
func (configurationApi ConfigurationApi) Delete(id string, body []byte, options map[string]string) ([]byte, error) {
	request := api.CcaRequest{
		Method:   api.DELETE,
		Body:     body,
		Endpoint: configurationApi.buildEndpoint() + "/" + id,
		Options:  options,
	}
	response, err := configurationApi.apiClient.Do(request)
	if err != nil {
		return nil, err
	} else if response.IsError() {
		return nil, api.CcaErrorResponse(*response)
	}
	return response.Data, nil
}
