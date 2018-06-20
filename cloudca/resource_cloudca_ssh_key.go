package cloudca

import (
	"fmt"
	"log"

	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudcaSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: createSSHKey,
		Read:   readSSHKey,
		Delete: deleteSSHKey,

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
	ccaResources, rerr := getResourcesForEnvironmentId(meta.(*cca.CcaClient), d.Get("environment_id").(string))

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
	ccaResources, rerr := getResourcesForEnvironmentId(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	sk, err := ccaResources.SSHKeys.Get(d.Id())
	if err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				log.Printf("SSH key with id='%s' was not found", d.Id())
				d.SetId("")
				return nil
			}
		}
		return err
	}
	d.Set("name", sk.Name)
	return nil
}

func deleteSSHKey(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentId(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}

	if _, err := ccaResources.SSHKeys.Delete(d.Id()); err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Printf("SSH key %s not found", d.Id())
				d.SetId("")
				return nil
			}
		}
		return err
	}

	return nil
}
