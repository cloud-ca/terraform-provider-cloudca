package cloudca

import (
	"fmt"
	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceCloudcaNetworkAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaNetworkAclCreate,
		Read:   resourceCloudcaNetworkAclRead,
		Delete: resourceCloudcaNetworkAclDelete,

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
				Description: "Name of environment where the network ACL should be created",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of network ACL",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Description of network ACL",
			},
			"vpc": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name or id of the VPC",
			},
		},
	}
}

func resourceCloudcaNetworkAclCreate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	vpcId, verr := retrieveVpcId(&ccaResources, d.Get("vpc").(string))
	if verr != nil {
		return verr
	}

	aclToCreate := cloudca.NetworkAcl{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		VpcId:             vpcId,
	}
	options := map[string]string{}
	if orgId, ok := d.GetOk("organization_code"); ok {
		options["org_id"] = orgId.(string)
	}
	newAcl, err := ccaResources.NetworkAcls.Create(aclToCreate, options)
	if err != nil {
		return fmt.Errorf("Error creating the new network ACL %s: %s", aclToCreate.Name, err)
	}
	d.SetId(newAcl.Id)
	return resourceCloudcaNetworkAclRead(d, meta)
}

func resourceCloudcaNetworkAclRead(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	acl, aErr := ccaResources.NetworkAcls.Get(acl.Id)
	if aErr != nil {
		if ccaError, ok := aErr.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("ACL %s not found", acl.Id)
				d.SetId("")
				return nil
			}
		}
		return aErr
	}

	vpc, vErr := ccaResources.Vpcs.Get(acl.VpcId)
	if vErr != nil {
		if ccaError, ok := vErr.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("Vpc %s not found", acl.VpcId)
				d.SetId("")
				return nil
			}
		}
		return vErr
	}

	// Update the config
	d.Set("name", acl.Name)
	d.Set("description", acl.Description)
	setValueOrID(d, "vpc", vpc.Name, acl.VpcId)

	return nil
}

func resourceCloudcaNetworkAclDelete(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)
	if _, err := ccaResources.NetworkAcls.Delete(d.Id()); err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("Network ACL %s not found", d.Id())
				d.SetId("")
				return nil
			}
		}
		return err
	}
	return nil
}

