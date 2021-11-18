package cloudca

import (
	"fmt"

	cca "github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudcaSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: createSSHKey,
		Read:   readSSHKey,
		Delete: deleteSSHKey,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of environment where the SSH key should be created",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the SSH Key",
			},
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func createSSHKey(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	name := d.Get("name").(string)
	publicKey := d.Get("public_key").(string)

	sk := cloudca.SSHKey{
		Name:      name,
		PublicKey: publicKey,
	}
	newSk, err := ccaResources.SSHKeys.Create(sk)
	if err != nil {
		return fmt.Errorf("Error creating new SSH key %s", err)
	}
	d.SetId(newSk.ID)
	return readSSHKey(d, meta)
}

func readSSHKey(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}

	sk, err := ccaResources.SSHKeys.Get(d.Id())

	if err != nil {
		return handleNotFoundError("SSH key", false, err, d)
	}

	if err := d.Set("name", sk.Name); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	return nil
}

func deleteSSHKey(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}

	if _, err := ccaResources.SSHKeys.Delete(d.Id()); err != nil {
		return handleNotFoundError("SSH key", true, err, d)
	}

	return nil
}
