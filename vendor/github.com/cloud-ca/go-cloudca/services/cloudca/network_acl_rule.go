package cloudca

import (
	"encoding/json"

	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

type NetworkAclRule struct {
	Id           string `json:"id,omitempty"`
	NetworkAclId string `json:"networkAclId,omitempty"`
	RuleNumber   string `json:"ruleNumber,omitempty"`
	Cidr         string `json:"cidr,omitempty"`
	Action       string `json:"action,omitempty"`
	Protocol     string `json:"protocol,omitempty"`
	StartPort    string `json:"startPort,omitempty"`
	EndPort      string `json:"endPort,omitempty"`
	IcmpType     string `json:"icmpType,omitempty"`
	IcmpCode     string `json:"icmpCode,omitempty"`
	TrafficType  string `json:"trafficType,omitempty"`
	State        string `json:"state,omitempty"`
}

type NetworkAclRuleService interface {
	Get(id string) (*NetworkAclRule, error)
	List() ([]NetworkAclRule, error)
	ListByNetworkAclId(networkAclId string) ([]NetworkAclRule, error)
	ListWithOptions(options map[string]string) ([]NetworkAclRule, error)
	Create(networkAclRule NetworkAclRule) (*NetworkAclRule, error)
	Update(id string, networkAclRule NetworkAclRule) (*NetworkAclRule, error)
	Delete(id string) (bool, error)
}

type NetworkAclRuleApi struct {
	entityService services.EntityService
}

func NewNetworkAclRuleService(apiClient api.ApiClient, serviceCode string, environmentName string) NetworkAclRuleService {
	return &NetworkAclRuleApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, NETWORK_ACL_RULE_ENTITY_TYPE),
	}
}

func parseNetworkAclRule(data []byte) *NetworkAclRule {
	networkAclRule := NetworkAclRule{}
	json.Unmarshal(data, &networkAclRule)
	return &networkAclRule
}

func parseNetworkAclRuleList(data []byte) []NetworkAclRule {
	aclRules := []NetworkAclRule{}
	json.Unmarshal(data, &aclRules)
	return aclRules
}

//Get network acl rule with the specified id for the current environment
func (networkAclRuleApi *NetworkAclRuleApi) Get(id string) (*NetworkAclRule, error) {
	data, err := networkAclRuleApi.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseNetworkAclRule(data), nil
}

func (networkAclRuleApi *NetworkAclRuleApi) ListWithOptions(options map[string]string) ([]NetworkAclRule, error) {
	data, err := networkAclRuleApi.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseNetworkAclRuleList(data), nil
}

func (networkAclRuleApi *NetworkAclRuleApi) List() ([]NetworkAclRule, error) {
	return networkAclRuleApi.ListWithOptions(map[string]string{})
}

//List all network acl rules for the NetworkAcl
func (networkAclRuleApi *NetworkAclRuleApi) ListByNetworkAclId(networkAclId string) ([]NetworkAclRule, error) {
	return networkAclRuleApi.ListWithOptions(map[string]string{"network_acl_id": networkAclId})
}

func (networkAclRuleApi *NetworkAclRuleApi) Create(networkAclRule NetworkAclRule) (*NetworkAclRule, error) {
	msg, err := json.Marshal(networkAclRule)
	if err != nil {
		return nil, err
	}
	result, err := networkAclRuleApi.entityService.Create(msg, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseNetworkAclRule(result), nil
}

func (networkAclRuleApi *NetworkAclRuleApi) Update(id string, networkAclRule NetworkAclRule) (*NetworkAclRule, error) {
	msg, err := json.Marshal(networkAclRule)
	if err != nil {
		return nil, err
	}
	result, err := networkAclRuleApi.entityService.Update(id, msg, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseNetworkAclRule(result), nil
}

func (networkAclRuleApi *NetworkAclRuleApi) Delete(id string) (bool, error) {
	_, err := networkAclRuleApi.entityService.Delete(id, []byte{}, map[string]string{})
	return err == nil, err
}
