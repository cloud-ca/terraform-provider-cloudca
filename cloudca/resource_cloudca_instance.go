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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "A cloudca service code",
			},
			"environment_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Name of environment where instance should be created",
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Name of instance",
			},

			"template": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Name or id of the template to use for this instance",
			},

			"compute_offering": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: "Name or id of the compute offering to use for this instance",
			},

			"network": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Name or id of the tier into which the new instance will be created",
			},

			"disk_offering": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Create and attach a volume with this disk offering (name or id) to the new instance",
			},

			"ssh_keyname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "SSH key name to attach to the new instance. Note: Cannot be used with public key.",
			},

			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Public key to attach to the new instance. Note: Cannot be used with SSH key name.",
			},

			"volume": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Name or id of volume to attach to this instance.",
			},

			"portsToForward": &schema.Schema{
				Type:     schema.TypeList,
				Elem: 	  &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Description: "List of port forwarding rules for this instance. Note: Might acquire a public IP if necessary",
			},

			"user_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Additional data passed to the new instance during its initialization",
			},

			"purge": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "If true, then it will purge the instance on destruction",
			},
		},
	}
}

func resourceCloudcaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*gocca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)
	
	computeOfferingId, _ := retrieveComputeOfferingID(&ccaResources, d.Get("compute_offering").(string))

	templateId, _ := retrieveTemplateID(&ccaResources, d.Get("template").(string))

	networkId, _ := retrieveNetworkID(&ccaResources, d.Get("network").(string))

	//
	instanceToCreate := cloudca.Instance{Name: d.Get("name").(string),
		ComputeOfferingId: computeOfferingId,
		TemplateId:        templateId,
		NetworkId:         networkId,
	}

	if sshKeyname, ok := d.GetOk("ssk_keyname"); ok {
		instanceToCreate.SSHKeyName = sshKeyname.(string)
	}
	if publicKey, ok := d.GetOk("public_key"); ok {
		instanceToCreate.PublicKey = publicKey.(string)
	}
	if volumeToAttach, ok := d.GetOk("volume"); ok {
		volumeToAttachId, _ := retrieveVolumeID(&ccaResources, volumeToAttach.(string))
		instanceToCreate.VolumeIdToAttach = volumeToAttachId
	}
	if portsToForward, ok := d.GetOk("ports_to_forward"); ok {
		instanceToCreate.PortsToForward = portsToForward.([]string)
	}
	if userData, ok := d.GetOk("user_data"); ok {
		instanceToCreate.UserData = userData.(string)
	}

	newInstance, err := ccaResources.Instances.Create(instanceToCreate)
	if err != nil {
		return fmt.Errorf("Error creating the new instance %s: %s", instanceToCreate.Name, err)
	}

	d.SetId(newInstance.Id)
	d.SetConnInfo(map[string]string{
		"host": newInstance.IpAddress,
		"user":  newInstance.Username,
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
				fmt.Errorf("[DEBUG] Instance %s does no longer exist", d.Get("name").(string))
				d.SetId("")
				return nil
			}
		}
		return err
	}

	// Update the config
	setValueOrID(d, "name", instance.Name, instance.Id)
	setValueOrID(d, "template", instance.TemplateName, instance.TemplateId)
	setValueOrID(d, "compute_offering", instance.ComputeOfferingName, instance.ComputeOfferingId)
	setValueOrID(d, "network", instance.NetworkName, instance.NetworkId)

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

	    if strings.EqualFold(offering.Name,name) {
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

	    if strings.EqualFold(template.Name,name) {
	    	log.Printf("Found template: %+v", template)
	    	return template.Id, nil
	    }
	}

	return "", fmt.Errorf("Template with name %s not found", name)
}

func retrieveNetworkID(ccaRes *cloudca.Resources, name string) (id string, err error) {
	if isID(name) {
		return name, nil
	}

	tiers, err := ccaRes.Tiers.List()
	if err != nil {
		return "", err
	}
	for _, tier := range tiers {

	    if strings.EqualFold(tier.Name,name) {
	    	log.Printf("Found tier: %+v", tier)
	    	return tier.Id, nil
	    }
	}

	return "", fmt.Errorf("Network with name %s not found", name)
}

func retrieveVolumeID(ccaRes *cloudca.Resources, name string) (id string, err error) {
	if isID(name) {
		return name, nil
	}

	volumes, err := ccaRes.Volumes.ListOfType(cloudca.VOLUME_TYPE_DATA)
	if err != nil {
		return "", err
	}
	for _, volume := range volumes {

		if strings.EqualFold(volume.Name,name) {
	    	log.Printf("Found volume: %+v", volume)
	    	return volume.Id, nil
	    }
	}

	return "", fmt.Errorf("Volume with name %s not found", name)
}

