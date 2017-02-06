package cloudca

import (
	"fmt"

	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudcaStaticNat() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaStaticNatCreate,
		Read:   resourceCloudcaStaticNatRead,
		Delete: resourceCloudcaStaticNatDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of environment where static NAT should be enabled",
			},
			"public_ip_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The public IP to enable static NAT on",
			},
			"private_ip_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private IP to enable static NAT on",
			},
		},
	}
}

func resourceCloudcaStaticNatCreate(d *schema.ResourceData, meta interface{}) error {
	ccaResources := getResourcesForEnvironmentId(d, meta)
	staticNatPublicIp := cloudca.PublicIp{
		Id:          d.Get("public_ip_id").(string),
		PrivateIpId: d.Get("private_ip_id").(string),
	}
	_, err := ccaResources.PublicIps.EnableStaticNat(staticNatPublicIp)
	if err != nil {
		return fmt.Errorf("Error enabling static NAT: %s", err)
	}
	d.SetId(staticNatPublicIp.Id)
	return resourceCloudcaStaticNatRead(d, meta)
}

func resourceCloudcaStaticNatRead(d *schema.ResourceData, meta interface{}) error {
	ccaResources := getResourcesForEnvironmentId(d, meta)
	publicIp, err := ccaResources.PublicIps.Get(d.Id())
	if err != nil {
		return handleNotFoundError(err, d)
	}
	if publicIp.PrivateIpId == "" {
		// If the private IP ID is missing, it means the public IP no longer has static NAT
		// enabled and so this entity is "missing" (at least as far as terraform is concerned).
		d.SetId("")
		return nil
	}
	d.Set("private_ip_id", publicIp.PrivateIpId)
	return nil
}

func resourceCloudcaStaticNatDelete(d *schema.ResourceData, meta interface{}) error {
	ccaResources := getResourcesForEnvironmentId(d, meta)
	_, err := ccaResources.PublicIps.DisableStaticNat(d.Id())
	return handleNotFoundError(err, d)
}
