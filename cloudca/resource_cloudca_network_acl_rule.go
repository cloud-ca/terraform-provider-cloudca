package cloudca

import (
	"fmt"
	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudcaNetworkAclRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaNetworkAclRuleCreate,
		Update: resoruceCloudcaNetworkAclRuleUpdate,
		Read:   resourceCloudcaNetworkAclRuleRead,
		Delete: resourceCloudcaNetworkAclRuleDelete,

		Schema: map[string]*schema.Schema{
			"service_code": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A cloudca service code",
			},
			"environment_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of environment where the network ACL rule should be created",
			},
			"rule_number": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The rule number of network ACL",
			},
			"action": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network ACL rule action (i.e. Allow or Deny)",
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
			"protocol": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The network ACL rule protocol (i.e. TCP, UDP, ICMP or ALL)",
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
			"traffic_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network ACL rule traffc type (i.e. Ingress or Egress)",
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
			"icmp_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ICMP type. Can only be used with ICMP protocol.",
			},
			"icmp_code": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ICMP code. Can only be used with ICMP protocol.",
			},
			"network_acl_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the network ACL of the network ACL rule",
			},
		},
	}
}

func resourceCloudcaNetworkAclRuleCreate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	aclRuleToCreate := cloudca.NetworkAclRule{
		RuleNumber:   d.Get("rule_number").(string),
		Action:       d.Get("action").(string),
		Protocol:     d.Get("protocol").(string),
		TrafficType:  d.Get("traffic_type").(string),
		IcmpType:     d.Get("icmp_type").(string),
		IcmpCode:     d.Get("icmp_code").(string),
		NetworkAclId: d.Get("network_acl_id").(string),
	}
	newAclRule, err := ccaResources.NetworkAclRules.Create(aclRuleToCreate)
	if err != nil {
		return fmt.Errorf("Error creating the new network ACL rule %s: %s", aclRuleToCreate.RuleNumber, err)
	}
	d.SetId(newAclRule.Id)
	return resourceCloudcaNetworkAclRead(d, meta)
}

func resourceCloudcaNetworkAclRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	if d.HasChange("rule_number") || d.HasChange("action") ||
		d.HasChange("traffic_type") || d.HasChange("icmp_type") ||
		d.HasChange("icmp_code") {
		aclRuleToUpdate := cloudca.NetworkAclRule{
			Id:          d.Id(),
			RuleNumber:  d.Get("rule_number").(string),
			Action:      d.Get("action").(string),
			TrafficType: d.Get("traffic_type").(string),
			IcmpType:    d.Get("icmp_type").(string),
			IcmpCode:    d.Get("icmp_code").(string),
		}
		_, err := ccaResources.NetworkAclRules.Update(d.Id(), aclRuleToUpdate)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceCloudcaNetworkAclRuleRead(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	aclRule, aErr := ccaResources.NetworkAclRules.Get(d.Id())
	if aErr != nil {
		if ccaError, ok := aErr.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("ACL rule %s not found", d.Id())
				d.SetId("")
				return nil
			}
		}
		return aErr
	}

	d.Set("rule_number", aclRule.RuleNumber)
	d.Set("action", aclRule.Action)
	d.Set("protocol", aclRule.Protocol)
	d.Set("traffic_type", aclRule.TrafficType)
	d.Set("network_acl_id", aclRule.NetworkAclId)

	return nil
}

func resourceCloudcaNetworkAclRuleDelete(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)
	if _, err := ccaResources.NetworkAclRules.Delete(d.Id()); err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("Network ACL rule %s not found", d.Id())
				d.SetId("")
				return nil
			}
		}
		return err
	}
	return nil
}
