package cloudca

import (
	"fmt"
	"log"

	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudcaPublicIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaPublicIPCreate,
		Read:   resourceCloudcaPublicIPRead,
		Delete: resourceCloudcaPublicIPDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of environment where the public IP should be created",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the VPC",
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudcaPublicIPCreate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	vpcID := d.Get("vpc_id").(string)

	publicIPToCreate := cloudca.PublicIp{
		VpcId: vpcID,
	}
	newPublicIP, err := ccaResources.PublicIps.Acquire(publicIPToCreate)
	if err != nil {
		return fmt.Errorf("Error acquiring the new public ip %s", err)
	}
	d.SetId(newPublicIP.Id)
	return resourceCloudcaPublicIPRead(d, meta)
}

func resourceCloudcaPublicIPRead(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	publicIP, err := ccaResources.PublicIps.Get(d.Id())
	if err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				log.Printf("Public Ip with id='%s' was not found", d.Id())
				d.SetId("")
				return nil
			}
		}
		return err
	}
	_ = d.Set("vpc_id", publicIP.VpcId)
	_ = d.Set("ip_address", publicIP.IpAddress)
	return nil
}

func resourceCloudcaPublicIPDelete(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}

	if _, err := ccaResources.PublicIps.Release(d.Id()); err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				_ = fmt.Errorf("Public Ip %s not found", d.Id())
				d.SetId("")
				return nil
			}
		}
		return err
	}

	return nil
}
