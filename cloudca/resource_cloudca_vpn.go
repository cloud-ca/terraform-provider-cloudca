package cloudca

import (
	"fmt"

	"github.com/cloud-ca/go-cloudca"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudcaVpn() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaVpnCreate,
		Read:   resourceCloudcaVpnRead,
		Delete: resourceCloudcaVpnDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the environment where the vpn should be created",
			},
			"certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Certificate to use when using IKEV2 vpn type",
			},
			"preshared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Preshared key to use when using L2TP vpn type",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Public IP address associated with the vpn",
			},
			"public_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the public IP address associated with the vpn",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "State of the vpn",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Type of vpn connection",
			},
		},
	}
}

func resourceCloudcaVpnCreate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))
	if rerr != nil {
		return rerr
	}
	_, err := ccaResources.RemoteAccessVpn.Enable(d.Get("public_ip_id").(string))
	if err != nil {
		return fmt.Errorf("Error enabling the VPN: %s", err)
	}
	d.SetId(d.Get("public_ip_id").(string))
	return resourceCloudcaVpnRead(d, meta)
}

func resourceCloudcaVpnRead(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))
	if rerr != nil {
		return rerr
	}

	vpn, err := ccaResources.RemoteAccessVpn.Get(d.Get("public_ip_id").(string))
	if err != nil {
		return handleNotFoundError("VPN", false, err, d)
	}

	if vpn.State == "Disabled" {
		// If the VPN is disabled, it means the VPN is not active
		// so this entity is "missing" (at least as far as terraform is concerned).
		d.SetId("")
		return nil
	}
	if err := d.Set("id", vpn.Id); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}
	if err := d.Set("state", vpn.State); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}
	if err := d.Set("certificate", vpn.Certificate); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}
	if err := d.Set("preshared_key", vpn.PresharedKey); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}
	if err := d.Set("public_ip", vpn.PublicIpAddress); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}
	if err := d.Set("public_ip_id", vpn.PublicIpAddressId); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}
	if err := d.Set("type", vpn.Type); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}
	return nil
}

func resourceCloudcaVpnDelete(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))
	if rerr != nil {
		return rerr
	}
	if _, err := ccaResources.RemoteAccessVpn.Disable(d.Id()); err != nil {
		return handleNotFoundError("VPN Delete", true, err, d)
	}
	return nil
}
