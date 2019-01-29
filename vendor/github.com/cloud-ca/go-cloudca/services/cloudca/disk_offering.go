package cloudca

import (
	"encoding/json"

	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

type DiskOffering struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	GbSize     int    `json:"gbSize,omitempty"`
	MinIops    int    `json:"minIops,omitempty"`
	MaxIops    int    `json:"maxIops,omitempty"`
	CustomSize bool   `json:"customSize,omitempty"`
	CustomIops bool   `json:"customIops,omitempty"`
}

type DiskOfferingService interface {
	Get(id string) (*DiskOffering, error)
	List() ([]DiskOffering, error)
	ListWithOptions(options map[string]string) ([]DiskOffering, error)
}

type DiskOfferingApi struct {
	entityService services.EntityService
}

func NewDiskOfferingService(apiClient api.ApiClient, serviceCode string, environmentName string) DiskOfferingService {
	return &DiskOfferingApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, DISK_OFFERING_ENTITY_TYPE),
	}
}

func parseDiskOffering(data []byte) *DiskOffering {
	diskOffering := DiskOffering{}
	json.Unmarshal(data, &diskOffering)
	return &diskOffering
}

func parseDiskOfferingList(data []byte) []DiskOffering {
	diskOfferings := []DiskOffering{}
	json.Unmarshal(data, &diskOfferings)
	return diskOfferings
}

//Get disk offering with the specified id for the current environment
func (diskOfferingApi *DiskOfferingApi) Get(id string) (*DiskOffering, error) {
	data, err := diskOfferingApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseDiskOffering(data), nil
}

//List all disk offerings for the current environment
func (diskOfferingApi *DiskOfferingApi) List() ([]DiskOffering, error) {
	return diskOfferingApi.ListWithOptions(map[string]string{})
}

//List all disk offerings for the current environment. Can use options to do sorting and paging.
func (diskOfferingApi *DiskOfferingApi) ListWithOptions(options map[string]string) ([]DiskOffering, error) {
	data, err := diskOfferingApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseDiskOfferingList(data), nil
}
