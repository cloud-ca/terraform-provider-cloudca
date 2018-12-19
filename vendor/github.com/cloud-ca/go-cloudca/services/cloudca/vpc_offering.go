package cloudca

import (
	"encoding/json"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

type VpcOffering struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	State string `json:"state,omitempty"`
}

type VpcOfferingService interface {
	Get(id string) (*VpcOffering, error)
	List() ([]VpcOffering, error)
	ListWithOptions(options map[string]string) ([]VpcOffering, error)
}

type VpcOfferingApi struct {
	entityService services.EntityService
}

func NewVpcOfferingService(apiClient api.ApiClient, serviceCode string, environmentName string) VpcOfferingService {
	return &VpcOfferingApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, VPC_OFFERING_ENTITY_TYPE),
	}
}

func parseVpcOffering(data []byte) *VpcOffering {
	vpcOffering := VpcOffering{}
	json.Unmarshal(data, &vpcOffering)
	return &vpcOffering
}

func parseVpcOfferingList(data []byte) []VpcOffering {
	vpcOfferings := []VpcOffering{}
	json.Unmarshal(data, &vpcOfferings)
	return vpcOfferings
}

//Get disk offering with the specified id for the current environment
func (vpcOfferingApi *VpcOfferingApi) Get(id string) (*VpcOffering, error) {
	data, err := vpcOfferingApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseVpcOffering(data), nil
}

//List all disk offerings for the current environment
func (vpcOfferingApi *VpcOfferingApi) List() ([]VpcOffering, error) {
	return vpcOfferingApi.ListWithOptions(map[string]string{})
}

//List all disk offerings for the current environment. Can use options to do sorting and paging.
func (vpcOfferingApi *VpcOfferingApi) ListWithOptions(options map[string]string) ([]VpcOffering, error) {
	data, err := vpcOfferingApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseVpcOfferingList(data), nil
}
