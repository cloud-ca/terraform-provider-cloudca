package cloudca

import (
	"fmt"
	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func resourceCloudcaPortForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: createPortForwardingRules,
		Read: readPortForwardingRules,
		Delete: deletePortForwardingRule,

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
				Description: "Name of environment where port forwarding rule should be created",
			},
			"ip_address_id": &schema.Schema{
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The public IP to which these rules should be applied",
			},
			"forward": &schema.Schema{
				Type: schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance": &schema.Schema{
							Type: schema.TypeString,
							Required: false,
							ForceNew: true,
							Description: "Name or ID of the instance that this rule should be applied to. If a private IP is not specified, the instance's primary private IP will be selected.",
						},
						"ip_address_id": &schema.Schema{
							Type: schema.TypeString,
							Required: false,
							ForceNew: true,
							Description: "The ID of the private IP to bind to. Does not require an instance to be specified.",
						},
						"protocol": &schema.Schema{
							Type: schema.TypeString,
							Required: true,
							ForceNew: true,
							Description: "The protocol that this rule should use (eg. TCP, UDP)",
						},
						"private_port_range_start": &schema.Schema{
							Type: schema.TypeInt,
							Required: false,
							ForceNew: true,
							Description: "The start of the private port range for this rule",
						},
						"private_port_range_end": &schema.Schema{
							Type: schema.TypeInt,
							Required: false,
							ForceNew: true,
							Description: "The end of the private port range for this rule",
						},
						"public_port_range_start": &schema.Schema{
							Type: schema.TypeInt,
							Required: false,
							ForceNew: true,
							Description: "The start of the public port range for this rule",
						},
						"public_port_range_end": &schema.Schema{
							Type: schema.TypeInt,
							Required: false,
							ForceNew: true,
							Description: "The end of the public port range for this rule",
						},
						"private_port": &schema.Schema{
							Type: schema.TypeInt,
							Required: false,
							ForceNew: true,
							Description: "The private port for this rule",
						},
						"public_port": &schema.Schema{
							Type: schema.TypeInt,
							Required: false,
							ForceNew: true,
							Description: "The public port for this rule",
						},
					},
				},
			},
		},
	}
}

func arePortTypesMutuallyExcluded(d map[string]interface{}) bool {
	all_range_fields_present := d["private_port_range_start"] != nil && d["private_port_range_end"] != nil && d["public_port_range_start"] != nil && d["public_port_range_end"] != nil
	all_non_range_fields_present := d["private_port"] != nil && d["public_port"] != nil

	return all_range_fields_present != all_non_range_fields_present
}

func normalizePorts(d map[string]interface{}) {
	private_port := d["private_port"]
	public_port := d["public_port"]
	if private_port != nil && public_port != nil {
		d["private_port_range_start"] = private_port.(int)
		d["private_port_range_end"] = private_port.(int)

		d["public_port_range_start"] = public_port.(int)
		d["public_port_range_end"] = public_port.(int)
	}
}

func retrieveInstanceId(ccaResources *cloudca.Resources, name string) (string, error) {
	if isID(name) {
		return name, nil
	}

	instances, err := ccaResources.Instances.List()
	if err != nil {
		return "", err
	}

	for _, instance := range instances {
		if strings.EqualFold(instance.Name, name) {
			return instance.Id, nil
		}
	}

	return "", fmt.Errorf("Couldn't find any instance with ID %s", name)
}

func createPortForwardingRules(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("ip_address_id").(string))
	createdRules := resourceCloudcaPortForwardingRule().Schema["forward"].ZeroValue().(*schema.Set)
	forwards := d.Get("forward").(*schema.Set)

	for _, forward := range forwards.List() {
		if err := createPortForwardingRule(d, meta, forward.(map[string]interface{})); err != nil {
			return err
		}

		if forward.(map[string]interface{})["id"].(string) != "" {
			forwards.Add(forward)
		}
	}

	d.Set("forwards", createdRules)
	return readPortForwardingRules(d, meta)
}

func createPortForwardingRule(d *schema.ResourceData, meta interface{}, forward map[string]interface{}) error {
	if !arePortTypesMutuallyExcluded(forward) {
		return fmt.Errorf("Cannot mix port ranges and plain ports in port forwarding rule definition")
	}

	client := meta.(*cca.CcaClient)
	resources, err := client.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))

	if err != nil {
		return err
	}

	ccaResources := resources.(cloudca.Resources)
	instance_id, err := retrieveInstanceId(&ccaResources, forward["instance"].(string))
	if err != nil {
		return err
	}

	pfr := cloudca.PortForwardingRule{
		PublicIpId: d.Get("ip_address_id").(string),
		InstanceId: instance_id,
		Protocol: forward["protocol"].(string),
		PublicPortStart: forward["public_port_range_start"].(string),
		PublicPortEnd: forward["public_port_range_end"].(string),
		PrivateIpId: forward["ip_address_id"].(string),
		PrivatePortStart: forward["private_port_range_start"].(string),
		PrivatePortEnd: forward["private_port_range_end"].(string),
	}

	newPfr, err := ccaResources.PortForwardingRules.Create(pfr)
	if err != nil {
		return err
	}

	forward["id"] = newPfr.Id
	return nil
}

func readPortForwardingRules(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func deletePortForwardingRule(d *schema.ResourceData, meta interface{}) error {
	return nil
}
