package cloudca

import (
	"encoding/json"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

type NetworkOffering struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type NetworkOfferingService interface {
	Get(id string) (*NetworkOffering, error)
	List() ([]NetworkOffering, error)
	ListWithOptions(options map[string]string) ([]NetworkOffering, error)
}

type NetworkOfferingApi struct {
	entityService services.EntityService
}

func NewNetworkOfferingService(apiClient api.ApiClient, serviceCode string, environmentName string) NetworkOfferingService {
	return &NetworkOfferingApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, NETWORK_OFFERING_ENTITY_TYPE),
	}
}

func parseNetworkOffering(data []byte) *NetworkOffering {
	networkOffering := NetworkOffering{}
	json.Unmarshal(data, &networkOffering)
	return &networkOffering
}

func parseNetworkOfferingList(data []byte) []NetworkOffering {
	networkOfferings := []NetworkOffering{}
	json.Unmarshal(data, &networkOfferings)
	return networkOfferings
}

//Get network offering with the specified id for the current environment
func (networkOfferingApi *NetworkOfferingApi) Get(id string) (*NetworkOffering, error) {
	data, err := networkOfferingApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseNetworkOffering(data), nil
}

//List all network offerings for the current environment
func (networkOfferingApi *NetworkOfferingApi) List() ([]NetworkOffering, error) {
	return networkOfferingApi.ListWithOptions(map[string]string{})
}

//List all network offerings for the current environment. Can use options to do sorting and paging.
func (networkOfferingApi *NetworkOfferingApi) ListWithOptions(options map[string]string) ([]NetworkOffering, error) {
	data, err := networkOfferingApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseNetworkOfferingList(data), nil
}
