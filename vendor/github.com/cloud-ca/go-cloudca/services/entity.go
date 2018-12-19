package services

import (
	"github.com/cloud-ca/go-cloudca/api"
)

//A generic service to access any entity
type EntityService interface {
	Get(id string, options map[string]string) ([]byte, error)
	List(options map[string]string) ([]byte, error)
	Execute(id string, operation string, body []byte, options map[string]string) ([]byte, error)
	Create(body []byte, options map[string]string) ([]byte, error)
	Update(id string, body []byte, options map[string]string) ([]byte, error)
	Delete(id string, body []byte, options map[string]string) ([]byte, error)
}

//Implementation of the EntityService
type EntityApi struct {
	apiClient       api.ApiClient
	taskService     TaskService
	serviceCode     string
	environmentName string
	entityType      string
}

func NewEntityService(apiClient api.ApiClient, serviceCode string, environmentName string, entityType string) EntityService {
	return &EntityApi{
		apiClient:       apiClient,
		taskService:     NewTaskService(apiClient),
		serviceCode:     serviceCode,
		environmentName: environmentName,
		entityType:      entityType,
	}
}

func (entityApi *EntityApi) buildEndpoint() string {
	return "/services/" + entityApi.serviceCode + "/" + entityApi.environmentName + "/" + entityApi.entityType
}

//Get an entity. Returns a []byte (of a json object) that should be unmarshalled to a specific entity
func (entityApi *EntityApi) Get(id string, options map[string]string) ([]byte, error) {
	request := api.CcaRequest{
		Method:   api.GET,
		Endpoint: entityApi.buildEndpoint() + "/" + id,
		Options:  options,
	}
	response, err := entityApi.apiClient.Do(request)
	if err != nil {
		return nil, err
	} else if response.IsError() {
		return nil, api.CcaErrorResponse(*response)
	}
	return response.Data, nil
}

//Get an entity list. Returns a []byte (of a json object) that should be unmarshalled to a specific entity
func (entityApi *EntityApi) List(options map[string]string) ([]byte, error) {
	request := api.CcaRequest{
		Method:   api.GET,
		Endpoint: entityApi.buildEndpoint(),
		Options:  options,
	}
	response, err := entityApi.apiClient.Do(request)
	if err != nil {
		return nil, err
	} else if response.IsError() {
		return nil, api.CcaErrorResponse(*response)
	}
	return response.Data, nil
}

//Execute a specific operation on an entity. Returns a []byte (of a json object) that should be unmarshalled to a specific entity
func (entityApi *EntityApi) Execute(id string, operation string, body []byte, options map[string]string) ([]byte, error) {
	optionsCopy := map[string]string{}
	for k, v := range options {
		optionsCopy[k] = v
	}
	optionsCopy["operation"] = operation
	endpoint := entityApi.buildEndpoint()
	if id != "" {
		endpoint = endpoint + "/" + id
	}
	request := api.CcaRequest{
		Method:   api.POST,
		Body:     body,
		Endpoint: endpoint,
		Options:  optionsCopy,
	}
	response, err := entityApi.apiClient.Do(request)
	if err != nil {
		return nil, err
	} else if response.IsError() {
		return nil, api.CcaErrorResponse(*response)
	}

	return entityApi.taskService.PollResponse(response, DEFAULT_POLLING_INTERVAL)
}

//Create a new entity described in the body parameter (json object). Returns a []byte (of a json object) that should be unmarshalled to a specific entity
func (entityApi *EntityApi) Create(body []byte, options map[string]string) ([]byte, error) {
	request := api.CcaRequest{
		Method:   api.POST,
		Body:     body,
		Endpoint: entityApi.buildEndpoint(),
		Options:  options,
	}
	response, err := entityApi.apiClient.Do(request)
	if err != nil {
		return nil, err
	} else if response.IsError() {
		return nil, api.CcaErrorResponse(*response)
	}
	return entityApi.taskService.PollResponse(response, DEFAULT_POLLING_INTERVAL)
}

//Update entity with specified id described in the body parameter (json object). Returns a []byte (of a json object) that should be unmarshalled to a specific entity
func (entityApi *EntityApi) Update(id string, body []byte, options map[string]string) ([]byte, error) {
	request := api.CcaRequest{
		Method:   api.PUT,
		Body:     body,
		Endpoint: entityApi.buildEndpoint() + "/" + id,
		Options:  options,
	}
	response, err := entityApi.apiClient.Do(request)
	if err != nil {
		return nil, err
	} else if response.IsError() {
		return nil, api.CcaErrorResponse(*response)
	}
	return entityApi.taskService.PollResponse(response, DEFAULT_POLLING_INTERVAL)
}

//Delete specified id described. A body (json object) can be provided if some fields must be sent to server. Returns a []byte (of a json object) that should be unmarshalled to a specific entity
func (entityApi EntityApi) Delete(id string, body []byte, options map[string]string) ([]byte, error) {
	request := api.CcaRequest{
		Method:   api.DELETE,
		Body:     body,
		Endpoint: entityApi.buildEndpoint() + "/" + id,
		Options:  options,
	}
	response, err := entityApi.apiClient.Do(request)
	if err != nil {
		return nil, err
	} else if response.IsError() {
		return nil, api.CcaErrorResponse(*response)
	}
	return entityApi.taskService.PollResponse(response, DEFAULT_POLLING_INTERVAL)
}
