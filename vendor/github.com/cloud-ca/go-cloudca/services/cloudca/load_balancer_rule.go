package cloudca

import (
	"encoding/json"

	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services"
)

const (
	UPDATE_INSTANCES  = "updateInstances"
	UPDATE_STICKINESS = "updateStickiness"
)

type LoadBalancerRule struct {
	Id                         string            `json:"id,omitempty"`
	Name                       string            `json:"name,omitempty"`
	InstanceIds                []string          `json:"instanceIds,omitempty"`
	NetworkId                  string            `json:"networkId,omitempty"`
	PublicIp                   string            `json:"publicIp,omitempty"`
	PublicIpId                 string            `json:"publicIpId,omitempty"`
	PublicPort                 string            `json:"publicPort,omitempty"`
	PrivatePort                string            `json:"privatePort,omitempty"`
	Protocol                   string            `json:"protocol,omitempty"`
	Algorithm                  string            `json:"algorithm,omitempty"`
	StickinessMethod           string            `json:"stickinessMethod,omitempty"`
	StickinessPolicyParameters map[string]string `json:"stickinessPolicyParameters,omitempty"`
}

type LoadBalancerRuleService interface {
	Get(id string) (*LoadBalancerRule, error)
	List() ([]LoadBalancerRule, error)
	ListWithOptions(options map[string]string) ([]LoadBalancerRule, error)
	Create(lbr LoadBalancerRule) (*LoadBalancerRule, error)
	Delete(id string) error
	Update(lbr LoadBalancerRule) (*LoadBalancerRule, error)
	SetLoadBalancerRuleInstances(id string, instanceIds []string) error
	SetLoadBalancerRuleStickinessPolicy(id string, method string, stickinessPolicyParameters map[string]string) error
	RemoveLoadBalancerRuleStickinessPolicy(id string) error
}

type LoadBalancerRuleApi struct {
	entityService services.EntityService
}

func NewLoadBalancerRuleService(apiClient api.ApiClient, serviceCode string, environmentName string) LoadBalancerRuleService {
	return &LoadBalancerRuleApi{
		entityService: services.NewEntityService(apiClient, serviceCode, environmentName, LOAD_BALANCER_RULE_ENTITY_TYPE),
	}
}

func parseLoadBalancerRule(data []byte) *LoadBalancerRule {
	lbr := LoadBalancerRule{}
	json.Unmarshal(data, &lbr)
	return &lbr
}

func parseLoadBalancerRuleList(data []byte) []LoadBalancerRule {
	lbrs := []LoadBalancerRule{}
	json.Unmarshal(data, &lbrs)
	return lbrs
}

func (api *LoadBalancerRuleApi) Get(id string) (*LoadBalancerRule, error) {
	data, err := api.entityService.Get(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseLoadBalancerRule(data), nil
}

func (api *LoadBalancerRuleApi) ListWithOptions(options map[string]string) ([]LoadBalancerRule, error) {
	data, err := api.entityService.List(options)
	if err != nil {
		return nil, err
	}
	return parseLoadBalancerRuleList(data), nil
}

func (api *LoadBalancerRuleApi) List() ([]LoadBalancerRule, error) {
	return api.ListWithOptions(map[string]string{})
}

func (api *LoadBalancerRuleApi) Create(lbr LoadBalancerRule) (*LoadBalancerRule, error) {
	msg, err := json.Marshal(lbr)
	if err != nil {
		return nil, err
	}
	result, err := api.entityService.Create(msg, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseLoadBalancerRule(result), nil
}

func (api *LoadBalancerRuleApi) Update(lbr LoadBalancerRule) (*LoadBalancerRule, error) {
	msg, err := json.Marshal(lbr)
	if err != nil {
		return nil, err
	}
	result, err := api.entityService.Update(lbr.Id, msg, map[string]string{})
	if err != nil {
		return nil, err
	}
	return parseLoadBalancerRule(result), nil
}

func (api *LoadBalancerRuleApi) SetLoadBalancerRuleInstances(id string, instanceIds []string) error {
	lbr := LoadBalancerRule{
		Id:          id,
		InstanceIds: instanceIds,
	}
	msg, err := json.Marshal(lbr)
	if err != nil {
		return err
	}
	_, updateErr := api.entityService.Execute(id, UPDATE_INSTANCES, msg, map[string]string{})
	return updateErr
}

func (api *LoadBalancerRuleApi) SetLoadBalancerRuleStickinessPolicy(id string, method string, stickinessPolicyParameters map[string]string) error {
	lbr := LoadBalancerRule{
		Id:                         id,
		StickinessMethod:           method,
		StickinessPolicyParameters: stickinessPolicyParameters,
	}
	msg, err := json.Marshal(lbr)
	if err != nil {
		return err
	}
	_, updateErr := api.entityService.Execute(id, UPDATE_STICKINESS, msg, map[string]string{})
	return updateErr
}

func (api *LoadBalancerRuleApi) RemoveLoadBalancerRuleStickinessPolicy(id string) error {
	lbr := LoadBalancerRule{
		Id:               id,
		StickinessMethod: "none",
	}
	msg, err := json.Marshal(lbr)
	if err != nil {
		return err
	}
	_, updateErr := api.entityService.Execute(id, UPDATE_STICKINESS, msg, map[string]string{})
	return updateErr
}

func (api *LoadBalancerRuleApi) Delete(id string) error {
	_, err := api.entityService.Delete(id, []byte{}, map[string]string{})
	return err
}
