package configuration

import (
	"encoding/json"

	"github.com/cloud-ca/go-cloudca/api"
)

type Organization struct {
	Id           string        `json:"id,omitempty"`
	Name         string        `json:"name,omitempty"`
	EntryPoint   string        `json:"entryPoint,omitempty"`
	Users        []User        `json:"users"`
	Environments []Environment `json:"environments"`
	Roles        []Role        `json:"roles"`
}

type OrganizationService interface {
	Get(id string) (*Organization, error)
	List() ([]Organization, error)
	ListWithOptions(options map[string]string) ([]Organization, error)
}

type OrganizationApi struct {
	configurationService ConfigurationService
}

func NewOrganizationService(apiClient api.ApiClient) OrganizationService {
	return &OrganizationApi{
		configurationService: NewConfigurationService(apiClient, "organizations"),
	}
}

func parseOrganization(data []byte) *Organization {
	organization := Organization{}
	json.Unmarshal(data, &organization)
	return &organization
}

func parseOrganizationList(data []byte) []Organization {
	organizations := []Organization{}
	json.Unmarshal(data, &organizations)
	return organizations
}

func (organizationApi *OrganizationApi) Get(id string) (*Organization, error) {
	data, err := organizationApi.configurationService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseOrganization(data), nil
}

//List all organizations
func (organizationApi *OrganizationApi) List() ([]Organization, error) {
	return organizationApi.ListWithOptions(map[string]string{})
}

//List all organizations. Can use options to do sorting and paging.
func (organizationApi *OrganizationApi) ListWithOptions(options map[string]string) ([]Organization, error) {
	data, err := organizationApi.configurationService.List(options)
	if err != nil {
		return nil, err
	}
	return parseOrganizationList(data), nil
}
