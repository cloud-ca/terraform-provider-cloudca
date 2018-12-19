package cloudca

import (
	"encoding/json"

	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

const (
	PFR_CREATE = "create"
	PFR_DELETE = "delete"
)

type PortForwardingRule struct {
	Id               string `json:"id,omitempty"`
	InstanceId       string `json:"instanceId,omitempty"`
	InstanceName     string `json:"instanceName,omitempty"`
	NetworkId        string `json:"networkId,omitempty"`
	PrivateIp        string `json:"privateIp,omitempty"`
	PrivateIpId      string `json:"privateIpId,omitempty"`
	PrivatePortStart string `json:"privatePortStart,omitempty"`
	PrivatePortEnd   string `json:"privatePortEnd,omitempty"`
	PublicIp         string `json:"ipAddress,omitempty"`
	PublicIpId       string `json:"ipAddressId,omitempty"`
	PublicPortStart  string `json:"publicPortStart,omitempty"`
	PublicPortEnd    string `json:"publicPortEnd,omitempty"`
	Protocol         string `json:"protocol,omitempty"`
	State            string `json:"state,omitempty"`
	VpcId            string `json:"vpcId,omitempty"`
}

type PortForwardingRuleService interface {
	Get(id string) (*PortForwardingRule, error)
	List() ([]PortForwardingRule, error)
	ListWithOptions(options map[string]string) ([]PortForwardingRule, error)
	Create(pfr PortForwardingRule) (*PortForwardingRule, error)
	Delete(id string) (bool, error)
}

type PortForwardingRuleApi struct {
	entityService services.EntityService
}

func NewPortForwardingRuleService(apiClient api.ApiClient, serviceCode string, environmentName string) PortForwardingRuleService {
	return &PortForwardingRuleApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, PORT_FORWARDING_RULE_ENTITY_TYPE),
	}
}

func parsePortForwardingRule(data []byte) *PortForwardingRule {
	pfr := PortForwardingRule{}
	json.Unmarshal(data, &pfr)
	return &pfr
}

func parsePortForwardingRuleList(data []byte) []PortForwardingRule {
	pfrs := []PortForwardingRule{}
	json.Unmarshal(data, &pfrs)
	return pfrs
}

func (api *PortForwardingRuleApi) Get(id string) (*PortForwardingRule, error) {
	data, err := api.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parsePortForwardingRule(data), nil
}

func (api *PortForwardingRuleApi) ListWithOptions(options map[string]string) ([]PortForwardingRule, error) {
	data, err := api.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parsePortForwardingRuleList(data), nil
}

func (api *PortForwardingRuleApi) List() ([]PortForwardingRule, error) {
	return api.ListWithOptions(map[string]string{})
}

func (api *PortForwardingRuleApi) Create(pfr PortForwardingRule) (*PortForwardingRule, error) {
	msg, err := json.Marshal(pfr)
	if err != nil {
		return nil, err
	}
	result, err := api.entityService.Create(msg, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parsePortForwardingRule(result), nil
}

func (api *PortForwardingRuleApi) Delete(id string) (bool, error) {
	_, err := api.entityService.Delete(id, []byte{}, map[string]string{})
	return err == nil, err
}
