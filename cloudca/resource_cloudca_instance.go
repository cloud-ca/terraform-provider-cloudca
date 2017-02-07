package cloudca

import (
	"fmt"
	"log"
	"strings"

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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A cloudca service code",
			},
			"environment_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of environment where instance should be created",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of instance",
			},

			"template": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name or id of the template to use for this instance",
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},

			"compute_offering": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name or id of the compute offering to use for this instance",
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},

			"network_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the tier into which the new instance will be created",
			},

			"ssh_key_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSH key name to attach to the new instance. Note: Cannot be used with public key.",
			},

			"public_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Public key to attach to the new instance. Note: Cannot be used with SSH key name.",
			},

			"user_data": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional data passed to the new instance during its initialization",
			},

			"cpu_count": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The instances CPU count. If the compute offering is custom, this value is required",
			},

			"memory_in_mb": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The instance's memory in MB. If the compute offering is custom, this value is required",
			},
			"private_ip_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudcaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	computeOfferingId, cerr := retrieveComputeOfferingID(&ccaResources, d.Get("compute_offering").(string))

	if cerr != nil {
		return cerr
	}

	templateId, terr := retrieveTemplateID(&ccaResources, d.Get("template").(string))

	if terr != nil {
		return terr
	}

	instanceToCreate := cloudca.Instance{Name: d.Get("name").(string),
		ComputeOfferingId: computeOfferingId,
		TemplateId:        templateId,
		NetworkId:         d.Get("network_id").(string),
	}

	if sshKeyname, ok := d.GetOk("ssh_key_name"); ok {
		instanceToCreate.SSHKeyName = sshKeyname.(string)
	}
	if publicKey, ok := d.GetOk("public_key"); ok {
		instanceToCreate.PublicKey = publicKey.(string)
	}
	if userData, ok := d.GetOk("user_data"); ok {
		instanceToCreate.UserData = userData.(string)
	}

	hasCustomFields := false
	if cpuCount, ok := d.GetOk("cpu_count"); ok {
		instanceToCreate.CpuCount = cpuCount.(int)
		hasCustomFields = true
	}
	if memoryInMB, ok := d.GetOk("memory_in_mb"); ok {
		instanceToCreate.MemoryInMB = memoryInMB.(int)
		hasCustomFields = true
	}

	computeOffering, cerr := ccaResources.ComputeOfferings.Get(computeOfferingId)
	if cerr != nil {
		return cerr
	} else if !computeOffering.Custom && hasCustomFields {
		return fmt.Errorf("Cannot have a CPU count or memory in MB because \"%s\" isn't a custom compute offering", computeOffering.Name)
	}

	newInstance, err := ccaResources.Instances.Create(instanceToCreate)
	if err != nil {
		return fmt.Errorf("Error creating the new instance %s: %s", instanceToCreate.Name, err)
	}

	d.SetId(newInstance.Id)
	d.SetConnInfo(map[string]string{
		"host":     newInstance.IpAddress,
		"user":     newInstance.Username,
		"password": newInstance.Password,
	})

	return resourceCloudcaInstanceRead(d, meta)
}

func resourceCloudcaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	// Get the virtual machine details
	instance, err := ccaResources.Instances.Get(d.Id())
	if err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("Instance %s does no longer exist", d.Get("name").(string))
				d.SetId("")
				return nil
			}
		}
		return err
	}

	// Update the config
	d.Set("name", instance.Name)
	setValueOrID(d, "template", strings.ToLower(instance.TemplateName), instance.TemplateId)
	setValueOrID(d, "compute_offering", strings.ToLower(instance.ComputeOfferingName), instance.ComputeOfferingId)
	d.Set("network_id", instance.NetworkId)
	d.Set("private_ip_id", instance.IpAddressId)
	d.Set("private_ip", instance.IpAddress)

	return nil
}

func resourceCloudcaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	d.Partial(true)

	if d.HasChange("compute_offering") || d.HasChange("cpu_count") || d.HasChange("memory_in_mb") {
		newComputeOffering := d.Get("compute_offering").(string)
		log.Printf("[DEBUG] Compute offering has changed for %s, changing compute offering...", newComputeOffering)
		newComputeOfferingId, ferr := retrieveComputeOfferingID(&ccaResources, newComputeOffering)
		if ferr != nil {
			return ferr
		}
		instanceToUpdate := cloudca.Instance{Id: d.Id(),
			ComputeOfferingId: newComputeOfferingId,
		}

		hasCustomFields := false
		if cpuCount, ok := d.GetOk("cpu_count"); ok {
			instanceToUpdate.CpuCount = cpuCount.(int)
			hasCustomFields = true
		}
		if memoryInMB, ok := d.GetOk("memory_in_mb"); ok {
			instanceToUpdate.MemoryInMB = memoryInMB.(int)
			hasCustomFields = true
		}

		computeOffering, cerr := ccaResources.ComputeOfferings.Get(newComputeOfferingId)
		if cerr != nil {
			return cerr
		} else if !computeOffering.Custom && hasCustomFields {
			return fmt.Errorf("Cannot have a CPU count or memory in MB because \"%s\" isn't a custom compute offering", computeOffering.Name)
		}

		_, err := ccaResources.Instances.ChangeComputeOffering(instanceToUpdate)
		if err != nil {
			return err
		}
		d.SetPartial("compute_offering")
	}

	if d.HasChange("ssh_key_name") {
		sshKeyName := d.Get("ssh_key_name").(string)
		log.Printf("[DEBUG] SSH key name has changed for %s, associating new SSH key...", sshKeyName)
		_, err := ccaResources.Instances.AssociateSSHKey(d.Id(), sshKeyName)
		if err != nil {
			return err
		}
		d.SetPartial("ssh_key_name")
	}

	d.Partial(false)

	return nil
}

func resourceCloudcaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	fmt.Println("[INFO] Destroying instance: %s", d.Get("name").(string))
	if _, err := ccaResources.Instances.Destroy(d.Id(), true); err != nil {
		if ccaError, ok := err.(api.CcaErrorResponse); ok {
			if ccaError.StatusCode == 404 {
				fmt.Errorf("Instance %s does no longer exist", d.Get("name").(string))
				d.SetId("")
				return nil
			}
		}
		return err
	}

	return nil
}

func retrieveComputeOfferingID(ccaRes *cloudca.Resources, name string) (id string, err error) {
	if isID(name) {
		return name, nil
	}

	computeOfferings, err := ccaRes.ComputeOfferings.List()
	if err != nil {
		return "", err
	}
	for _, offering := range computeOfferings {

		if strings.EqualFold(offering.Name, name) {
			log.Printf("Found compute offering: %+v", offering)
			return offering.Id, nil
		}
	}

	return "", fmt.Errorf("Compute offering with name %s not found", name)
}

func retrieveTemplateID(ccaRes *cloudca.Resources, name string) (id string, err error) {
	if isID(name) {
		return name, nil
	}

	templates, err := ccaRes.Templates.List()
	if err != nil {
		return "", err
	}
	for _, template := range templates {

		if strings.EqualFold(template.Name, name) {
			log.Printf("Found template: %+v", template)
			return template.Id, nil
		}
	}

	return "", fmt.Errorf("Template with name %s not found", name)
}
