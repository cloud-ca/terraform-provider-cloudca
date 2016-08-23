package cloudca

import (
	"fmt"
	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
	"strings"
)

func resourceCloudcaNetworkAclRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaNetworkAclRuleCreate,
		Update: resourceCloudcaNetworkAclRuleUpdate,
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
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The rule number of network ACL",
			},
			"cidr": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The network ACL rule cidr",
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
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ICMP type. Can only be used with ICMP protocol.",
			},
			"icmp_code": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ICMP code. Can only be used with ICMP protocol.",
			},
			"start_port": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The start port. Can only be used with TCP/UDP protocol.",
			},
			"end_port": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The end port. Can only be used with TCP/UDP protocol.",
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

	var icmpType string
	if v, ok := d.GetOk("icmp_type"); ok {
		icmpType = strconv.Itoa(v.(int))
	}
	var icmpCode string
	if v, ok := d.GetOk("icmp_code"); ok {
		icmpCode = strconv.Itoa(v.(int))
	}
	var startPort string
	if v, ok := d.GetOk("start_port"); ok {
		startPort = strconv.Itoa(v.(int))
	}
	var endPort string
	if v, ok := d.GetOk("end_port"); ok {
		endPort = strconv.Itoa(v.(int))
	}

	aclRuleToCreate := cloudca.NetworkAclRule{
		RuleNumber:   strconv.Itoa(d.Get("rule_number").(int)),
		Cidr:         d.Get("cidr").(string),
		Action:       d.Get("action").(string),
		Protocol:     d.Get("protocol").(string),
		TrafficType:  d.Get("traffic_type").(string),
		IcmpType:     icmpType,
		IcmpCode:     icmpCode,
		StartPort:    startPort,
		EndPort:      endPort,
		NetworkAclId: d.Get("network_acl_id").(string),
	}
	newAclRule, err := ccaResources.NetworkAclRules.Create(aclRuleToCreate)
	if err != nil {
		return fmt.Errorf("Error creating the new network ACL rule %s: %s", aclRuleToCreate.RuleNumber, err)
	}
	d.SetId(newAclRule.Id)
	return resourceCloudcaNetworkAclRuleRead(d, meta)
}

func resourceCloudcaNetworkAclRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	var icmpType string
	if v, ok := d.GetOk("icmp_type"); ok {
		icmpType = strconv.Itoa(v.(int))
	}
	var icmpCode string
	if v, ok := d.GetOk("icmp_code"); ok {
		icmpCode = strconv.Itoa(v.(int))
	}
	var startPort string
	if v, ok := d.GetOk("start_port"); ok {
		startPort = strconv.Itoa(v.(int))
	}
	var endPort string
	if v, ok := d.GetOk("end_port"); ok {
		endPort = strconv.Itoa(v.(int))
	}

	aclRuleToUpdate := cloudca.NetworkAclRule{
		Id:          d.Id(),
		RuleNumber:  strconv.Itoa(d.Get("rule_number").(int)),
		Cidr:        d.Get("cidr").(string),
		Action:      d.Get("action").(string),
		TrafficType: d.Get("traffic_type").(string),
		IcmpType:    icmpType,
		IcmpCode:    icmpCode,
		StartPort:   startPort,
		EndPort:     endPort,
	}
	_, err := ccaResources.NetworkAclRules.Update(d.Id(), aclRuleToUpdate)
	if err != nil {
		return err
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

	var icmpType int
	if aclRule.IcmpType != "" {
		icmpType, _ = strconv.Atoi(aclRule.IcmpType)
	}
	var icmpCode int
	if aclRule.IcmpCode != "" {
		icmpCode, _ = strconv.Atoi(aclRule.IcmpCode)
	}
	var startPort int
	if aclRule.StartPort != "" {
		startPort, _ = strconv.Atoi(aclRule.StartPort)
	}
	var endPort int
	if aclRule.EndPort != "" {
		endPort, _ = strconv.Atoi(aclRule.EndPort)
	}

	d.Set("rule_number", aclRule.RuleNumber)
	d.Set("action", strings.ToLower(aclRule.Action))
	d.Set("protocol", strings.ToLower(aclRule.Protocol))
	d.Set("traffic_type", strings.ToLower(aclRule.TrafficType))
	d.Set("icmp_type", icmpType)
	d.Set("icmp_code", icmpCode)
	d.Set("start_port", startPort)
	d.Set("end_port", endPort)
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
