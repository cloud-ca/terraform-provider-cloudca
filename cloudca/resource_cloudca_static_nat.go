package cloudca

import (
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
)

func resourceCloudcaStaticNat() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaStaticNatCreate,
		Read:   resourceCloudcaStaticNatRead,
		Delete: resourceCloudcaStaticNatDelete,

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
				Description: "Name of environment where tier should be created",
			},
			"organization_code": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Entry point of organization",
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
	resources := getResources(d, meta)
	staticNatPublicIp := cloudca.PublicIp{
		Id:          d.Get("public_ip_id").(string),
		PrivateIpId: d.Get("private_ip_id").(string),
	}
	_, err := resources.PublicIps.EnableStaticNat(staticNatPublicIp)
	if err != nil {
		return fmt.Errorf("Error enabling static NAT: %s", err)
	}
	d.SetId(staticNatPublicIp.Id)
	return resourceCloudcaStaticNatRead(d, meta)
}

func resourceCloudcaStaticNatRead(d *schema.ResourceData, meta interface{}) error {
	resources := getResources(d, meta)
	publicIp, err := resources.PublicIps.Get(d.Id())
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
	resources := getResources(d, meta)
	_, err := resources.PublicIps.DisableStaticNat(d.Id())
	return handleNotFoundError(err, d)
}
