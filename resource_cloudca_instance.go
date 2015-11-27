package cloudca

import (
	"fmt"
	"log"
	"strings"

	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudcaInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaInstanceCreate,
		Read:   resourceCloudcaInstanceRead,
		Update: resourceCloudcaInstanceUpdate,
		Delete: resourceCloudcaInstanceDelete,

		Schema: map[string]*schema.Schema{
			"service_code": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"template": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"compute_offering": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"network": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"disk_offering": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"sshKeyName": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"publicKey": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"volumeToAttach": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"userData": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"purge": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceCloudcaInstanceCreate(d *schema.ResourceData, meta interface{}) error {

}

func resourceCloudcaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*gocca.CcaClient)
	ccaResources := ccaClient.GetResources(d.Get("service_code"), d.Get("environment_name")).(cloudca.Resources)

	// Get the virtual machine details
	instance, err := ccaResources.Instances.Get(d.Id())
	if err != nil {
		if ccaError, ok := err.(CcaErrorResponse); ok {
			if ccaError.statusCode == 404 {
				log.Printf("[DEBUG] Instance %s does no longer exist", d.Get("name").(string))
				d.SetId("")
				return nil
			}
		}
		return err
	}

	// Update the config
	d.Set("name", instance.Name)
	d.Set("template", instance.TemplateName)
	d.Set("compute_offering", instance.ComputeOfferingName)
	d.Set("network", instance.NetworkName)

	return nil
}

func resourceCloudcaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCloudcaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*gocca.CcaClient)
	ccaResources := ccaClient.GetResources(d.Get("service_code"), d.Get("environment_name")).(cloudca.Resources)

	log.Printf("[INFO] Destroying instance: %s", d.Get("name").(string))
	if _, err := cs.VirtualMachine.DestroyVirtualMachine(p); err != nil {
		if ccaError, ok := err.(CcaErrorResponse); ok {
			if ccaError.statusCode == 404 {
				log.Printf("[DEBUG] Instance %s does no longer exist", d.Get("name").(string))
				d.SetId("")
				return nil
			}
		}
		return err
	}

	return nil
}
