package cloudca

import (
	"fmt"
	"github.com/cloud-ca/go-cloudca" 
	"github.com/cloud-ca/go-cloudca/api"
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

			"ssh_keyname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"volume": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"user_data": &schema.Schema{
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
	ccaClient := meta.(*gocca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)
	
	// Retrieve the compute_offering ID
	computeOfferingId := "b72c010a-0cbb-49b5-9b19-84ea671d7b3f"

	// Retrieve the template ID
	templateId := "4cdeea05-ae6e-49e8-8385-00502d29e55c"

	// Retrieve the network ID
	networkId := "4485b91d-c772-414a-a3dd-d973aebe841b"

	//
	instanceToCreate := cloudca.Instance{Name: d.Get("name").(string),
		ComputeOfferingId: computeOfferingId,
		TemplateId:        templateId,
		NetworkId:         networkId,
	}

	newInstance, err := ccaResources.Instances.Create(instanceToCreate)
	if err != nil {
		return fmt.Errorf("Error creating the new instance %s: %s", instanceToCreate.Name, err)
	}

	d.SetId(newInstance.Id)
	d.SetConnInfo(map[string]string{
		"privateIp": newInstance.IpAddress,
		"password":  newInstance.Password,
	})

	return resourceCloudcaInstanceRead(d, meta)
}

func resourceCloudcaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*gocca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	// Get the virtual machine details
	instance, err := ccaResources.Instances.Get(d.Id())
	if err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Println("[DEBUG] Instance %s does no longer exist", d.Get("name").(string))
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
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	fmt.Println("[INFO] Destroying instance: %s", d.Get("name").(string))
	if _, err := ccaResources.Instances.Destroy(d.Id(), d.Get("purge").(bool)); err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Println("[DEBUG] Instance %s does no longer exist", d.Get("name").(string))
				d.SetId("")
				return nil
			}
		}
		return err
	}

	return nil
}
