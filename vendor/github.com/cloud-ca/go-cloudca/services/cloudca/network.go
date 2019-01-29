package cloudca

import (
	"encoding/json"

	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

type Service struct {
	Name         string                 `json:"name,omitempty"`
	Capabilities map[string]interface{} `json:"capabilities,omitempty"`
}

type Network struct {
	Id                string    `json:"id,omitempty"`
	Name              string    `json:"name,omitempty"`
	Description       string    `json:"description,omitempty"`
	VpcId             string    `json:"vpcId,omitempty"`
	NetworkOfferingId string    `json:"networkOfferingId,omitempty"`
	NetworkAclId      string    `json:"networkAclId,omitempty"`
	NetworkAclName    string    `json:"networkAclName,omitempty"`
	ZoneId            string    `json:"zoneid,omitempty"`
	ZoneName          string    `json:"zonename,omitempty"`
	Cidr              string    `json:"cidr,omitempty"`
	Type              string    `json:"type,omitempty"`
	State             string    `json:"state,omitempty"`
	Gateway           string    `json:"gateway,omitempty"`
	IsSystem          bool      `json:"issystem,omitempty"`
	Domain            string    `json:"domain,omitempty"`
	DomainId          string    `json:"domainid,omitempty"`
	Project           string    `json:"project,omitempty"`
	ProjectId         string    `json:"projectid,omitempty"`
	Services          []Service `json:"service,omitempty"`
}

type NetworkService interface {
	Get(id string) (*Network, error)
	List() ([]Network, error)
	ListOfVpc(vpcId string) ([]Network, error)
	ListWithOptions(options map[string]string) ([]Network, error)
	Create(network Network, options map[string]string) (*Network, error)
	Update(id string, network Network) (*Network, error)
	Delete(id string) (bool, error)
	ChangeAcl(id string, aclId string) (bool, error)
}

type NetworkApi struct {
	entityService services.EntityService
}

func NewNetworkService(apiClient api.ApiClient, serviceCode string, environmentName string) NetworkService {
	return &NetworkApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, NETWORK_ENTITY_TYPE),
	}
}

func parseNetwork(data []byte) *Network {
	network := Network{}
	json.Unmarshal(data, &network)
	return &network
}

func parseNetworkList(data []byte) []Network {
	networks := []Network{}
	json.Unmarshal(data, &networks)
	return networks
}

//Get network with the specified id for the current environment
func (networkApi *NetworkApi) Get(id string) (*Network, error) {
	data, err := networkApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseNetwork(data), nil
}

//List all networks for the current environment
func (networkApi *NetworkApi) List() ([]Network, error) {
	return networkApi.ListWithOptions(map[string]string{})
}

//List all networks of a vpc for the current environment
func (networkApi *NetworkApi) ListOfVpc(vpcId string) ([]Network, error) {
	return networkApi.ListWithOptions(map[string]string{
		vpcId: vpcId,
	})
}

//List all networks for the current environment. Can use options to do sorting and paging.
func (networkApi *NetworkApi) ListWithOptions(options map[string]string) ([]Network, error) {
	data, err := networkApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseNetworkList(data), nil
}

func (networkApi *NetworkApi) Create(network Network, options map[string]string) (*Network, error) {
	send, merr := json.Marshal(network)
	if merr != nil {
		return nil, merr
	}
	body, err := networkApi.entityService.Create(send, options)
	if err != nil {
		return nil, err
	}
	return parseNetwork(body), nil
}

func (networkApi *NetworkApi) Update(id string, network Network) (*Network, error) {
	send, merr := json.Marshal(network)
	if merr != nil {
		return nil, merr
	}
	body, err := networkApi.entityService.Update(id, send, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseNetwork(body), nil
}

func (networkApi *NetworkApi) Delete(id string) (bool, error) {
	_, err := networkApi.entityService.Delete(id, []byte{}, map[string]string{})
	return err == nil, err
}

func (networkApi *NetworkApi) ChangeAcl(id string, aclId string) (bool, error) {
	send, merr := json.Marshal(Network{
		NetworkAclId: aclId,
	})
	if merr != nil {
		return false, merr
	}
	_, err := networkApi.entityService.Execute(id, "replace", send, map[string]string{})
	return err == nil, err
}
