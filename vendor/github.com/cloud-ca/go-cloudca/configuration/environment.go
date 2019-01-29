package configuration

import (
	"encoding/json"

	"github.com/cloud-ca/go-cloudca/api"
)

const (
	ENVIRONMENT_CONFIGURATION_TYPE = "environments"
)

type Environment struct {
	Id                string            `json:"id,omitempty"`
	Name              string            `json:"name,omitempty"`
	Description       string            `json:"description,omitempty"`
	Organization      Organization      `json:"organization,omitempty"`
	ServiceConnection ServiceConnection `json:"serviceConnection,omitempty"`
	Users             []User            `json:"users"`
	Roles             []Role            `json:"roles"`
}

type EnvironmentService interface {
	Get(id string) (*Environment, error)
	List() ([]Environment, error)
	ListWithOptions(options map[string]string) ([]Environment, error)
	Create(environment Environment) (*Environment, error)
	Update(id string, environment Environment) (*Environment, error)
	Delete(id string) (bool, error)
}

type EnvironmentApi struct {
	configurationService ConfigurationService
}

func NewEnvironmentService(apiClient api.ApiClient) EnvironmentService {
	return &EnvironmentApi{
		configurationService: NewConfigurationService(apiClient, ENVIRONMENT_CONFIGURATION_TYPE),
	}
}

func parseEnvironment(data []byte) *Environment {
	environment := Environment{}
	json.Unmarshal(data, &environment)
	return &environment
}

func parseEnvironmentList(data []byte) []Environment {
	environments := []Environment{}
	json.Unmarshal(data, &environments)
	return environments
}

//Get environment with the specified id
func (environmentApi *EnvironmentApi) Get(id string) (*Environment, error) {
	data, err := environmentApi.configurationService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseEnvironment(data), nil
}

//List all environments
func (environmentApi *EnvironmentApi) List() ([]Environment, error) {
	return environmentApi.ListWithOptions(map[string]string{})
}

//List all instances for the current environment. Can use options to do sorting and paging.
func (environmentApi *EnvironmentApi) ListWithOptions(options map[string]string) ([]Environment, error) {
	data, err := environmentApi.configurationService.List(options)
	if err != nil {
		return nil, err
	}
	return parseEnvironmentList(data), nil
}

//Create environment
func (environmentApi *EnvironmentApi) Create(environment Environment) (*Environment, error) {
	send, merr := json.Marshal(environment)
	if merr != nil {
		return nil, merr
	}
	body, err := environmentApi.configurationService.Create(send, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseEnvironment(body), nil
}

func (environmentApi *EnvironmentApi) Update(id string, environment Environment) (*Environment, error) {
	send, merr := json.Marshal(environment)
	if merr != nil {
		return nil, merr
	}
	body, err := environmentApi.configurationService.Update(id, send, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseEnvironment(body), nil
}

func (environmentApi *EnvironmentApi) Delete(id string) (bool, error) {
	_, err := environmentApi.configurationService.Delete(id, []byte{}, map[string]string{})
	return err == nil, err
}
