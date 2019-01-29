package cloudca

import (
	"encoding/json"

	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

type Template struct {
	Id               string   `json:"id,omitempty"`
	Name             string   `json:"name,omitempty"`
	Description      string   `json:"description,omitempty"`
	Size             int      `json:"size,omitempty"`
	IsPublic         bool     `json:"isPublic,omitempty"`
	IsReady          bool     `json:"isReady,omitempty"`
	SSHKeyEnabled    bool     `json:"sshKeyEnabled,omitempty"`
	Extractable      bool     `json:"extractable,omitempty"`
	Resizable        bool     `json:"resizable,omitempty"`
	OSType           string   `json:"osType,omitempty"`
	OSTypeId         string   `json:"osTypeId,omitempty"`
	Hypervisor       string   `json:"hypervisor,omitempty"`
	Format           string   `json:"format,omitempty"`
	ProjectId        string   `json:"projectId,omitempty"`
	AvailableInZones []string `json:"availableInZones,omitempty"`
}

type TemplateService interface {
	Get(id string) (*Template, error)
	List() ([]Template, error)
	ListWithOptions(options map[string]string) ([]Template, error)
}

type TemplateApi struct {
	entityService services.EntityService
}

func NewTemplateService(apiClient api.ApiClient, serviceCode string, environmentName string) TemplateService {
	return &TemplateApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, TEMPLATE_ENTITY_TYPE),
	}
}

func parseTemplate(data []byte) *Template {
	template := Template{}
	json.Unmarshal(data, &template)
	return &template
}

func parseTemplateList(data []byte) []Template {
	templates := []Template{}
	json.Unmarshal(data, &templates)
	return templates
}

//Get template with the specified id for the current environment
func (templateApi *TemplateApi) Get(id string) (*Template, error) {
	data, err := templateApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseTemplate(data), nil
}

//List all templates for the current environment
func (templateApi *TemplateApi) List() ([]Template, error) {
	return templateApi.ListWithOptions(map[string]string{})
}

//List all templates for the current environment. Can use options to do sorting and paging.
func (templateApi *TemplateApi) ListWithOptions(options map[string]string) ([]Template, error) {
	data, err := templateApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseTemplateList(data), nil
}
