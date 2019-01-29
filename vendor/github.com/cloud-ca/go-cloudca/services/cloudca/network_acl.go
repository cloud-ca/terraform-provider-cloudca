package cloudca

import (
	"encoding/json"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

type NetworkAcl struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	VpcId       string `json:"vpcId,omitempty"`
}

type NetworkAclService interface {
	Get(id string) (*NetworkAcl, error)
	List() ([]NetworkAcl, error)
	ListByVpcId(vpcId string) ([]NetworkAcl, error)
	ListWithOptions(options map[string]string) ([]NetworkAcl, error)
	Create(networkAcl NetworkAcl) (*NetworkAcl, error)
	Delete(id string) (bool, error)
}

type NetworkAclApi struct {
	entityService services.EntityService
}

func NewNetworkAclService(apiClient api.ApiClient, serviceCode string, environmentName string) NetworkAclService {
	return &NetworkAclApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, NETWORK_ACL_ENTITY_TYPE),
	}
}

func parseNetworkAcl(data []byte) *NetworkAcl {
	networkAcl := NetworkAcl{}
	json.Unmarshal(data, &networkAcl)
	return &networkAcl
}

func parseNetworkAclList(data []byte) []NetworkAcl {
	networkAcls := []NetworkAcl{}
	json.Unmarshal(data, &networkAcls)
	return networkAcls
}

//Get network acl with the specified id for the current environment
func (networkAclApi *NetworkAclApi) Get(id string) (*NetworkAcl, error) {
	data, err := networkAclApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseNetworkAcl(data), nil
}

//List all network offerings for the current environment
func (networkAclApi *NetworkAclApi) List() ([]NetworkAcl, error) {
	return networkAclApi.ListWithOptions(map[string]string{})
}

//List all network offerings for the current environment
func (networkAclApi *NetworkAclApi) ListByVpcId(vpcId string) ([]NetworkAcl, error) {
	return networkAclApi.ListWithOptions(map[string]string{"vpc_id": vpcId})
}

//List all network offerings for the current environment. Can use options to do sorting and paging.
func (networkAclApi *NetworkAclApi) ListWithOptions(options map[string]string) ([]NetworkAcl, error) {
	data, err := networkAclApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseNetworkAclList(data), nil
}

func (networkAclApi *NetworkAclApi) Create(networkAcl NetworkAcl) (*NetworkAcl, error) {
	msg, err := json.Marshal(networkAcl)
	if err != nil {
		return nil, err
	}
	result, err := networkAclApi.entityService.Create(msg, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseNetworkAcl(result), nil
}

func (networkAclApi *NetworkAclApi) Delete(id string) (bool, error) {
	_, err := networkAclApi.entityService.Delete(id, []byte{}, map[string]string{})
	return err == nil, err
}
