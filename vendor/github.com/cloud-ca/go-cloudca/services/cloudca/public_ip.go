package cloudca

import (
	"encoding/json"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

const (
	PUBLIC_IP_ENABLE_STATIC_NAT_OPERATION  = "enableStaticNat"
	PUBLIC_IP_DISABLE_STATIC_NAT_OPERATION = "disableStaticNat"
)

type PublicIp struct {
	Id            string   `json:"id,omitempty"`
	IpAddress     string   `json:"ipaddress,omitempty"`
	State         string   `json:"state,omitempty"`
	ZoneId        string   `json:"zoneId,omitempty"`
	ZoneName      string   `json:"zoneName,omitempty"`
	NetworkId     string   `json:"networkId,omitempty"`
	NetworkName   string   `json:"networkName,omitempty"`
	VpcId         string   `json:"vpcId,omitempty"`
	VpcName       string   `json:"vpcName,omitempty"`
	PrivateIpId   string   `json:"privateIpId,omitempty"`
	InstanceNames []string `json:"instanceNames,omitempty"`
	InstanceId    string   `json:"instanceId,omitempty"`
	Purposes      []string `json:"purposes,omitempty"`
	Ports         []string `json:"ports,omitempty"`
}

type PublicIpService interface {
	Get(id string) (*PublicIp, error)
	List() ([]PublicIp, error)
	Acquire(publicIp PublicIp) (*PublicIp, error)
	Release(id string) (bool, error)
	EnableStaticNat(publicIp PublicIp) (bool, error)
	DisableStaticNat(id string) (bool, error)
}

type PublicIpApi struct {
	entityService services.EntityService
}

func NewPublicIpService(apiClient api.ApiClient, serviceCode string, environmentName string) PublicIpService {
	return &PublicIpApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, PUBLIC_IP_ENTITY_TYPE),
	}
}

func parsePublicIp(data []byte) *PublicIp {
	publicIp := PublicIp{}
	json.Unmarshal(data, &publicIp)
	return &publicIp
}

func parsePublicIpList(data []byte) []PublicIp {
	publicIps := []PublicIp{}
	json.Unmarshal(data, &publicIps)
	return publicIps
}

func (publicIpApi *PublicIpApi) Get(id string) (*PublicIp, error) {
	data, err := publicIpApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parsePublicIp(data), nil
}

func (publicIpApi *PublicIpApi) List() ([]PublicIp, error) {
	return publicIpApi.ListWithOptions(map[string]string{})
}

func (publicIpApi *PublicIpApi) ListWithOptions(options map[string]string) ([]PublicIp, error) {
	data, err := publicIpApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parsePublicIpList(data), nil
}

func (publicIpApi *PublicIpApi) Acquire(publicIp PublicIp) (*PublicIp, error) {
	send, merr := json.Marshal(publicIp)
	if merr != nil {
		return nil, merr
	}
	body, err := publicIpApi.entityService.Create(send, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parsePublicIp(body), nil
}

func (publicIpApi *PublicIpApi) Release(id string) (bool, error) {
	_, err := publicIpApi.entityService.Delete(id, []byte{}, map[string]string{})
	return err == nil, err
}

func (publicIpApi *PublicIpApi) EnableStaticNat(publicIp PublicIp) (bool, error) {
	send, merr := json.Marshal(publicIp)
	if merr != nil {
		return false, merr
	}
	_, err := publicIpApi.entityService.Execute(publicIp.Id, PUBLIC_IP_ENABLE_STATIC_NAT_OPERATION, send, map[string]string{})
	return err == nil, err
}

func (publicIpApi *PublicIpApi) DisableStaticNat(id string) (bool, error) {
	_, err := publicIpApi.entityService.Execute(id, PUBLIC_IP_DISABLE_STATIC_NAT_OPERATION, []byte{}, map[string]string{})
	return err == nil, err
}
