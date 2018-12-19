package cloudca

import (
	"encoding/json"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

const (
	VPC_RESTART_ROUTER_OPERATION = "restart"
)

type Vpc struct {
	Id            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	VpcOfferingId string `json:"vpcOfferingId,omitempty"`
	State         string `json:"state,omitempty"`
	Cidr          string `json:"cidr,omitempty"`
	ZoneId        string `json:"zoneId,omitempty"`
	ZoneName      string `json:"zoneName,omitempty"`
	NetworkDomain string `json:"networkDomain,omitempty"`
	SourceNatIp   string `json:"sourceNatIp,omitempty"`
	VpnStatus     string `json:"vpnStatus,omitempty"`
	Type          string `json:"type,omitempty"`
}

type VpcService interface {
	Get(id string) (*Vpc, error)
	List() ([]Vpc, error)
	ListWithOptions(options map[string]string) ([]Vpc, error)
	Create(vpc Vpc) (*Vpc, error)
	Update(vpc Vpc) (*Vpc, error)
	Destroy(id string) (bool, error)
	RestartRouter(id string) (bool, error)
}

type VpcApi struct {
	entityService services.EntityService
}

func NewVpcService(apiClient api.ApiClient, serviceCode string, environmentName string) VpcService {
	return &VpcApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, VPC_ENTITY_TYPE),
	}
}

func parseVpc(data []byte) *Vpc {
	vpc := Vpc{}
	json.Unmarshal(data, &vpc)
	return &vpc
}

func parseVpcList(data []byte) []Vpc {
	vpcs := []Vpc{}
	json.Unmarshal(data, &vpcs)
	return vpcs
}

//Get vpc with the specified id for the current environment
func (vpcApi *VpcApi) Get(id string) (*Vpc, error) {
	data, err := vpcApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseVpc(data), nil
}

//List all vpcs for the current environment
func (vpcApi *VpcApi) List() ([]Vpc, error) {
	return vpcApi.ListWithOptions(map[string]string{})
}

//List all vpcs for the current environment. Can use options to do sorting and paging.
func (vpcApi *VpcApi) ListWithOptions(options map[string]string) ([]Vpc, error) {
	data, err := vpcApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseVpcList(data), nil
}

//Create an vpc in the current environment
func (vpcApi *VpcApi) Create(vpc Vpc) (*Vpc, error) {
	send, merr := json.Marshal(vpc)
	if merr != nil {
		return nil, merr
	}
	body, err := vpcApi.entityService.Create(send, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseVpc(body), nil
}

//Create an vpc in the current environment
func (vpcApi *VpcApi) Update(vpc Vpc) (*Vpc, error) {
	send, merr := json.Marshal(vpc)
	if merr != nil {
		return nil, merr
	}
	body, err := vpcApi.entityService.Update(vpc.Id, send, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseVpc(body), nil
}

//Destroy a vpc with specified id in the current environment
func (vpcApi *VpcApi) Destroy(id string) (bool, error) {
	_, err := vpcApi.entityService.Delete(id, []byte{}, map[string]string{})
	return err == nil, err
}

//Restart the router of the vpc with the specified id exists in the current environment
func (vpcApi *VpcApi) RestartRouter(id string) (bool, error) {
	_, err := vpcApi.entityService.Execute(id, VPC_RESTART_ROUTER_OPERATION, []byte{}, map[string]string{})
	return err == nil, err
}
