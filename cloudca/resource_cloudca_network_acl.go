package cloudca

import (
	"fmt"

	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudcaNetworkAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaNetworkAclCreate,
		Read:   resourceCloudcaNetworkAclRead,
		Delete: resourceCloudcaNetworkAclDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of environment where the network ACL should be created",
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
			"vpc_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the VPC",
			},
		},
	}
}

func resourceCloudcaNetworkAclCreate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentId(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}

	aclToCreate := cloudca.NetworkAcl{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		VpcId:       d.Get("vpc_id").(string),
	}
	newAcl, err := ccaResources.NetworkAcls.Create(aclToCreate)
	if err != nil {
		return fmt.Errorf("Error creating the new network ACL %s: %s", aclToCreate.Name, err)
	}
	d.SetId(newAcl.Id)
	return resourceCloudcaNetworkAclRead(d, meta)
}

func resourceCloudcaNetworkAclRead(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentId(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	acl, aErr := ccaResources.NetworkAcls.Get(d.Id())
	if aErr != nil {
		if ccaError, ok := aErr.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("ACL %s not found", d.Id())
				d.SetId("")
				return nil
			}
		}
		return aErr
	}

	// Update the config
	d.Set("name", acl.Name)
	d.Set("description", acl.Description)
	d.Set("vpc_id", acl.VpcId)

	return nil
}

func resourceCloudcaNetworkAclDelete(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentId(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
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
