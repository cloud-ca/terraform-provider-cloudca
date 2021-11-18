package cloudca

import (
	"fmt"

	cca "github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudcaStaticNAT() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaStaticNATCreate,
		Read:   resourceCloudcaStaticNATRead,
		Delete: resourceCloudcaStaticNATDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of environment where static NAT should be enabled",
			},
			"public_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The public IP to enable static NAT on",
			},
			"private_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private IP to enable static NAT on",
			},
		},
	}
}

func resourceCloudcaStaticNATCreate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	staticNATPublicIP := cloudca.PublicIp{
		Id:          d.Get("public_ip_id").(string),
		PrivateIpId: d.Get("private_ip_id").(string),
	}
	_, err := ccaResources.PublicIps.EnableStaticNat(staticNATPublicIP)
	if err != nil {
		return fmt.Errorf("Error enabling static NAT: %s", err)
	}
	d.SetId(staticNATPublicIP.Id)
	return resourceCloudcaStaticNATRead(d, meta)
}

func resourceCloudcaStaticNATRead(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	publicIP, err := ccaResources.PublicIps.Get(d.Id())
	if err != nil {
		return handleNotFoundError("Static NAT", false, err, d)
	}
	if publicIP.PrivateIpId == "" {
		// If the private IP ID is missing, it means the public IP no longer has static NAT
		// enabled and so this entity is "missing" (at least as far as terraform is concerned).
		d.SetId("")
		return nil
	}
	if err := d.Set("private_ip_id", publicIP.PrivateIpId); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}
	return nil
}

func resourceCloudcaStaticNATDelete(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	_, err := ccaResources.PublicIps.DisableStaticNat(d.Id())
	return handleNotFoundError("Static NAT", true, err, d)
}
