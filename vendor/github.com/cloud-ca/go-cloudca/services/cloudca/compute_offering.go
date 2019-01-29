package cloudca

import (
	"encoding/json"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

type ComputeOffering struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	MemoryInMB int    `json:"memoryInMB,omitempty"`
	CpuCount   int    `json:"cpuCount,omitempty"`
	Custom     bool   `json:"custom,omitempty"`
}

type ComputeOfferingService interface {
	Get(id string) (*ComputeOffering, error)
	List() ([]ComputeOffering, error)
	ListWithOptions(options map[string]string) ([]ComputeOffering, error)
}

type ComputeOfferingApi struct {
	entityService services.EntityService
}

func NewComputeOfferingService(apiClient api.ApiClient, serviceCode string, environmentName string) ComputeOfferingService {
	return &ComputeOfferingApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, COMPUTE_OFFERING_ENTITY_TYPE),
	}
}

func parseComputeOffering(data []byte) *ComputeOffering {
	computeOffering := ComputeOffering{}
	json.Unmarshal(data, &computeOffering)
	return &computeOffering
}

func parseComputeOfferingList(data []byte) []ComputeOffering {
	computeOfferings := []ComputeOffering{}
	json.Unmarshal(data, &computeOfferings)
	return computeOfferings
}

//Get compute offering with the specified id for the current environment
func (computeOfferingApi *ComputeOfferingApi) Get(id string) (*ComputeOffering, error) {
	data, err := computeOfferingApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseComputeOffering(data), nil
}

//List all compute offerings for the current environment
func (computeOfferingApi *ComputeOfferingApi) List() ([]ComputeOffering, error) {
	return computeOfferingApi.ListWithOptions(map[string]string{})
}

//List all compute offerings for the current environment. Can use options to do sorting and paging.
func (computeOfferingApi *ComputeOfferingApi) ListWithOptions(options map[string]string) ([]ComputeOffering, error) {
	data, err := computeOfferingApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseComputeOfferingList(data), nil
}
