package cloudca

import (
	"fmt"
	"log"

	cca "github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudcaVpn() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaVpnCreate,
		Read:   resourceCloudcaVpnRead,
		Delete: resourceCloudcaVpnDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the environment where the vpn should be created",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the VPC for the vpn",
			},
			"certificate": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Description: "Certificate to use when using IKEV2 vpn type",
			},
			"preshared_key": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Description: "Preshared key to use when using L2TP vpn type",
			},
			"public_ip": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Description: "Public IP address associated with the vpn",
			},
			"public_ip_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Description: "ID of the public IP address associated with the vpn",
			},
			"state": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Description: "State of the vpn",
			},
			"type": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
				Description: "Type of vpn connection",
			},
		},
	}
}

func resourceCloudcaVpnCreate(d *schema.ResourceData, meta interface{}) error {
	vpnIPPurpose := "SOURCE_NAT"
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))
	if rerr != nil {
		return rerr
	}

	var vpnPubIPID string
	pubIps, _ := ccaResources.PublicIps.List()
	for _, ip := range pubIps {
		if ip.VpcId == d.Get("vpc_id").(string) {
			for _, purpose := range ip.Purposes {
				if purpose == vpnIPPurpose {
					vpnPubIPID = ip.Id
					break
				}
			}
		}
		if vpnPubIPID != "" {
			break
		}
	}

	if vpnPubIPID == "" {
		return fmt.Errorf("Error enabling the VPN because no Source NAT IP was found for the VPC")
	}

	_, err := ccaResources.RemoteAccessVpn.Enable(vpnPubIPID)
	if err != nil {
		return fmt.Errorf("Error enabling the VPN: %s", err)
	}
	d.SetId(vpnPubIPID)
	return resourceCloudcaVpnRead(d, meta)
}

func resourceCloudcaVpnRead(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))
	if rerr != nil {
		return rerr
	}

	vpn, err := ccaResources.RemoteAccessVpn.Get(d.Id())
	if err != nil {
		return handleNotFoundError("VPN", false, err, d)
	}

	if vpn.State == "Disabled" {
		// If the VPN is disabled, it means the VPN is not active
		// so this entity is "missing" (at least as far as terraform is concerned).
		d.SetId("")
		return handleNotFoundError("VPN Disabled", false, err, d)
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
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				log.Printf("VPN with id=%s no longer exists", d.Id())
				d.SetId("")
				return nil
			}
			return handleNotFoundError("VPN Delete", true, err, d)
		}
		return handleNotFoundError("VPN Delete", true, err, d)
	}
	return nil
}
