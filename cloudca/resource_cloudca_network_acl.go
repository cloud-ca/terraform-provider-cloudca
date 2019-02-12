package cloudca

import (
	"fmt"

	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudcaNetworkACL() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaNetworkACLCreate,
		Read:   resourceCloudcaNetworkACLRead,
		Delete: resourceCloudcaNetworkACLDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of environment where the network ACL should be created",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of network ACL",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Description of network ACL",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the VPC",
			},
		},
	}
}

func resourceCloudcaNetworkACLCreate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}

	aclToCreate := cloudca.NetworkAcl{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		VpcId:       d.Get("vpc_id").(string),
	}
	newACL, err := ccaResources.NetworkAcls.Create(aclToCreate)
	if err != nil {
		return fmt.Errorf("Error creating the new network ACL %s: %s", aclToCreate.Name, err)
	}
	d.SetId(newACL.Id)
	return resourceCloudcaNetworkACLRead(d, meta)
}

func resourceCloudcaNetworkACLRead(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	acl, aErr := ccaResources.NetworkAcls.Get(d.Id())
	if aErr != nil {
		if ccaError, ok := aErr.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				_ = fmt.Errorf("ACL %s not found", d.Id())
				d.SetId("")
				return nil
			}
		}
		return aErr
	}

	// Update the config
	if err := d.Set("name", acl.Name); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	if err := d.Set("description", acl.Description); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	if err := d.Set("vpc_id", acl.VpcId); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	return nil
}

func resourceCloudcaNetworkACLDelete(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	if _, err := ccaResources.NetworkAcls.Delete(d.Id()); err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				_ = fmt.Errorf("Network ACL %s not found", d.Id())
				d.SetId("")
				return nil
			}
		}
		return err
	}
	return nil
}
