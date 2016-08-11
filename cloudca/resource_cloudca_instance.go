package cloudca

import (
	"fmt"
	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
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

			"network": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name or id of the tier into which the new instance will be created",
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
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

			"purge": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, then it will purge the instance on destruction",
			},
			"private_ip_id": &schema.Schema{
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

	networkId, nerr := retrieveNetworkID(&ccaResources, d.Get("network").(string))

	if nerr != nil {
		return nerr
	}

	instanceToCreate := cloudca.Instance{Name: d.Get("name").(string),
		ComputeOfferingId: computeOfferingId,
		TemplateId:        templateId,
		NetworkId:         networkId,
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
	setValueOrID(d, "network", strings.ToLower(instance.NetworkName), instance.NetworkId)
	d.Set("private_ip_id", instance.IpAddressId)

	return nil
}

func resourceCloudcaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	d.Partial(true)

	if d.HasChange("compute_offering") {
		newComputeOffering := d.Get("compute_offering").(string)
		log.Printf("[DEBUG] Compute offering has changed for %s, changing compute offering...", newComputeOffering)
		newComputeOfferingId, ferr := retrieveComputeOfferingID(&ccaResources, newComputeOffering)
		if ferr != nil {
			return ferr
		}
		_, err := ccaResources.Instances.ChangeComputeOffering(d.Id(), newComputeOfferingId)
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

func retrieveNetworkID(ccaRes *cloudca.Resources, name string) (id string, err error) {
	if isID(name) {
		return name, nil
	}

	tiers, err := ccaRes.Tiers.List()
	if err != nil {
		return "", err
	}
	for _, tier := range tiers {

		if strings.EqualFold(tier.Name, name) {
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

		if strings.EqualFold(volume.Name, name) {
			log.Printf("Found volume: %+v", volume)
			return volume.Id, nil
		}
	}

	return "", fmt.Errorf("Volume with name %s not found", name)
}
