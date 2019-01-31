package cloudca

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	TCP  = "TCP"
	UDP  = "UDP"
	ICMP = "ICMP"
)

func resourceCloudcaNetworkACLRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaNetworkACLRuleCreate,
		Update: resourceCloudcaNetworkACLRuleUpdate,
		Read:   resourceCloudcaNetworkACLRuleRead,
		Delete: resourceCloudcaNetworkACLRuleDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of environment where the network ACL rule should be created",
			},
			"rule_number": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The rule number of network ACL",
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network ACL rule cidr",
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network ACL rule action (i.e. Allow or Deny)",
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The network ACL rule protocol (i.e. TCP, UDP, ICMP or All)",
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
			"traffic_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network ACL rule traffc type (i.e. Ingress or Egress)",
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
			"icmp_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ICMP type. Can only be used with ICMP protocol.",
			},
			"icmp_code": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ICMP code. Can only be used with ICMP protocol.",
			},
			"start_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The start port. Can only be used with TCP/UDP protocol.",
			},
			"end_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The end port. Can only be used with TCP/UDP protocol.",
			},
			"network_acl_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the network ACL of the network ACL rule",
			},
		},
	}
}

func resourceCloudcaNetworkACLRuleCreate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	aclRuleToCreate := cloudca.NetworkAclRule{
		RuleNumber:   strconv.Itoa(d.Get("rule_number").(int)),
		Cidr:         d.Get("cidr").(string),
		Action:       d.Get("action").(string),
		Protocol:     d.Get("protocol").(string),
		TrafficType:  d.Get("traffic_type").(string),
		NetworkAclId: d.Get("network_acl_id").(string),
	}
	fillPortFields(d, &aclRuleToCreate)
	fillIcmpFields(d, &aclRuleToCreate)
	if !(strings.EqualFold(TCP, aclRuleToCreate.Protocol) || strings.EqualFold(UDP, aclRuleToCreate.Protocol)) && (aclRuleToCreate.StartPort != "" || aclRuleToCreate.EndPort != "") {
		return fmt.Errorf("Cannot have ports if not TCP or UDP protocol")
	}
	if !strings.EqualFold(ICMP, aclRuleToCreate.Protocol) && (aclRuleToCreate.IcmpType != "" || aclRuleToCreate.IcmpCode != "") {
		return fmt.Errorf("Cannot have icmp fields if not ICMP protocol")
	}

	newACLRule, err := ccaResources.NetworkAclRules.Create(aclRuleToCreate)
	if err != nil {
		return fmt.Errorf("Error creating the new network ACL rule %s: %s", aclRuleToCreate.RuleNumber, err)
	}
	d.SetId(newACLRule.Id)
	return resourceCloudcaNetworkACLRuleRead(d, meta)
}

func resourceCloudcaNetworkACLRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	aclRuleToUpdate := cloudca.NetworkAclRule{
		Id:          d.Id(),
		RuleNumber:  strconv.Itoa(d.Get("rule_number").(int)),
		Cidr:        d.Get("cidr").(string),
		Action:      d.Get("action").(string),
		Protocol:    d.Get("protocol").(string),
		TrafficType: d.Get("traffic_type").(string),
	}
	fillPortFields(d, &aclRuleToUpdate)
	fillIcmpFields(d, &aclRuleToUpdate)

	_, err := ccaResources.NetworkAclRules.Update(d.Id(), aclRuleToUpdate)
	if err != nil {
		return err
	}
	return nil
}

func resourceCloudcaNetworkACLRuleRead(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
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
	d.Set("action", strings.ToLower(aclRule.Action))
	d.Set("protocol", strings.ToLower(aclRule.Protocol))
	d.Set("traffic_type", strings.ToLower(aclRule.TrafficType))
	d.Set("icmp_type", readIntFromString(aclRule.IcmpType))
	d.Set("icmp_code", readIntFromString(aclRule.IcmpCode))
	d.Set("start_port", readIntFromString(aclRule.StartPort))
	d.Set("end_port", readIntFromString(aclRule.EndPort))
	d.Set("network_acl_id", aclRule.NetworkAclId)

	return nil
}

func resourceCloudcaNetworkACLRuleDelete(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
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

func fillPortFields(d *schema.ResourceData, aclRule *cloudca.NetworkAclRule) {
	if v, ok := d.GetOk("start_port"); ok {
		aclRule.StartPort = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("end_port"); ok {
		aclRule.EndPort = strconv.Itoa(v.(int))
	}
}

func fillIcmpFields(d *schema.ResourceData, aclRule *cloudca.NetworkAclRule) {
	if v, ok := d.GetOk("icmp_type"); ok {
		aclRule.IcmpType = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("icmp_code"); ok {
		aclRule.IcmpCode = strconv.Itoa(v.(int))
	}
}
