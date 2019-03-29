package cloudca

import (
	"fmt"

	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudcaVpnUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaVpnUserCreate,
		Read:   resourceCloudcaVpnUserRead,
		Delete: resourceCloudcaVpnUserDelete,

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the environment where the vpn should be created",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Username of the VPN user",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Password of the VPN user",
			},
		},
	}
}

func resourceCloudcaVpnUserCreate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))
	if rerr != nil {
		return rerr
	}

	remoteAccessVpnUser := cloudca.RemoteAccessVpnUser{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}
	_, err := ccaResources.RemoteAccessVpnUser.Create(remoteAccessVpnUser)
	if err != nil {
		return fmt.Errorf("Error adding VPN user: %s", err)
	}
	return resourceCloudcaVpnRead(d, meta)
}

func resourceCloudcaVpnUserRead(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))
	if rerr != nil {
		return rerr
	}

	// to set the Id, we need to list and loop through to find the Id we just created
	vpnUsers, err := ccaResources.RemoteAccessVpnUser.List()
	if err != nil {
		return fmt.Errorf("Error getting VPN user: %s", err)
	}
	found := false
	var vpnUser cloudca.RemoteAccessVpnUser
	for _, user := range vpnUsers {
		if user.Username == d.Get("username").(string) {
			vpnUser = user
			found = true
			break
		}
	}
	if !found {
		// If we can not find the user based on their 'username' then this
		// entity is "missing" (at least as far as terraform is concerned).
		d.SetId("")
		return nil
	}
	if err := d.Set("id", vpnUser.Id); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}
	if err := d.Set("username", vpnUser.Username); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}
	return nil
}

func resourceCloudcaVpnUserDelete(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))
	if rerr != nil {
		return rerr
	}
	remoteAccessVpnUser := cloudca.RemoteAccessVpnUser{
		Id:       d.Get("id").(string),
		Username: d.Get("username").(string),
	}
	if _, err := ccaResources.RemoteAccessVpnUser.Delete(remoteAccessVpnUser); err != nil {
		return handleNotFoundError("VPN UserDelete", true, err, d)
	}
	return nil
}
