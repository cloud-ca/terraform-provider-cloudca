package cloudca

import (
	"encoding/json"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

type AffinityGroup struct {
	Id            string   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	Type          string   `json:"type,omitempty"`
	InstanceIds   []string `json:"instanceIds,omitempty"`
	InstanceNames []string `json:"instanceNames,omitempty"`
	ZoneIds       []string `json:"zoneIds,omitempty"`
}

type AffinityGroupService interface {
	Get(string) (*AffinityGroup, error)
	List() ([]AffinityGroup, error)
	ListWithOptions(map[string]string) ([]AffinityGroup, error)
}

type AffinityGroupApi struct {
	entityService services.EntityService
}

func NewAffinityGroupsService(apiClient api.ApiClient, serviceCode string, environmentName string) AffinityGroupService {
	return &AffinityGroupApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, AFFINITY_GROUP_ENTITY_TYPE),
	}
}

func parseAffinityGroup(data []byte) *AffinityGroup {
	affinityGroup := AffinityGroup{}
	json.Unmarshal(data, &affinityGroup)
	return &affinityGroup
}

func parseAffinityGroupList(data []byte) []AffinityGroup {
	affinityGroups := []AffinityGroup{}
	json.Unmarshal(data, &affinityGroups)
	return affinityGroups
}

func (api *AffinityGroupApi) Get(id string) (*AffinityGroup, error) {
	resp, err := api.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseAffinityGroup(resp), nil
}

func (api *AffinityGroupApi) List() ([]AffinityGroup, error) {
	return api.ListWithOptions(map[string]string{})
}

func (api *AffinityGroupApi) ListWithOptions(options map[string]string) ([]AffinityGroup, error) {
	resp, err := api.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseAffinityGroupList(resp), nil
}
