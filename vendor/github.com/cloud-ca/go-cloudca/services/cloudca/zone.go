package cloudca

import (
	"encoding/json"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

type Zone struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ZoneService interface {
	Get(string) (*Zone, error)
	List() ([]Zone, error)
	ListWithOptions(map[string]string) ([]Zone, error)
}

type ZoneApi struct {
	entityService services.EntityService
}

func NewZoneService(apiClient api.ApiClient, serviceCode string, environmentName string) ZoneService {
	return &ZoneApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, ZONE_ENTITY_TYPE),
	}
}

func parseZone(data []byte) *Zone {
	zone := Zone{}
	json.Unmarshal(data, &zone)
	return &zone
}

func parseZoneList(data []byte) []Zone {
	zones := []Zone{}
	json.Unmarshal(data, &zones)
	return zones
}

func (api *ZoneApi) Get(id string) (*Zone, error) {
	resp, err := api.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseZone(resp), nil
}

func (api *ZoneApi) List() ([]Zone, error) {
	return api.ListWithOptions(map[string]string{})
}

func (api *ZoneApi) ListWithOptions(options map[string]string) ([]Zone, error) {
	resp, err := api.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseZoneList(resp), nil
}
