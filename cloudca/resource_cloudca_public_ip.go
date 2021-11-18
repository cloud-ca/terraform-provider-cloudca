package cloudca

import (
	"fmt"

	cca "github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudcaPublicIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaPublicIPCreate,
		Read:   resourceCloudcaPublicIPRead,
		Delete: resourceCloudcaPublicIPDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
		return fmt.Errorf("Error acquiring the new public IP %s", err)
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
		return handleNotFoundError("Public IP", false, err, d)
	}

	if err := d.Set("vpc_id", publicIP.VpcId); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	if err := d.Set("ip_address", publicIP.IpAddress); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	return nil
}

func resourceCloudcaPublicIPDelete(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}

	if _, err := ccaResources.PublicIps.Release(d.Id()); err != nil {
		return handleNotFoundError("Public IP", true, err, d)
	}

	return nil
}
