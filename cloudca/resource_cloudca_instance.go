package cloudca

import (
	"fmt"
	"log"
	"strings"

	cca "github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudcaInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaInstanceCreate,
		Read:   resourceCloudcaInstanceRead,
		Update: resourceCloudcaInstanceUpdate,
		Delete: resourceCloudcaInstanceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of environment where instance should be created",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of instance",
			},
			"template": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name or id of the template to use for this instance",
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
			"compute_offering": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name or id of the compute offering to use for this instance",
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the network into which the new instance will be created",
			},
			"ssh_key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "SSH key name to attach to the new instance. Note: Cannot be used with public key.",
			},
			"public_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Public key to attach to the new instance. Note: Cannot be used with SSH key name.",
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional data passed to the new instance during its initialization",
			},
			"cpu_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The instances CPU count. If the compute offering is custom, this value is required",
			},
			"memory_in_mb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The instance's memory in MB. If the compute offering is custom, this value is required",
			},
			"root_volume_size_in_gb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The size of the root volume in GB. This can only be set if the template allows choosing a custom root volume size.",
			},
			"private_ip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The IPv4 address of the instance. Must be within the network's CIDR and not collide with existing instances.",
			},
			"dedicated_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Id of the dedicated group into which the new instance will be created",
			},
		},
	}
}

func resourceCloudcaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}

	computeOfferingID, cerr := retrieveComputeOfferingID(&ccaResources, d.Get("compute_offering").(string))

	if cerr != nil {
		return cerr
	}

	templateID, terr := retrieveTemplateID(&ccaResources, d.Get("template").(string))

	if terr != nil {
		return terr
	}

	instanceToCreate := cloudca.Instance{Name: d.Get("name").(string),
		ComputeOfferingId: computeOfferingID,
		TemplateId:        templateID,
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
	if privateIP, ok := d.GetOk("private_ip"); ok {
		instanceToCreate.IpAddress = privateIP.(string)
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

	computeOffering, cerr := ccaResources.ComputeOfferings.Get(computeOfferingID)
	if cerr != nil {
		return cerr
	} else if !computeOffering.Custom && hasCustomFields {
		return fmt.Errorf("Cannot have a CPU count or memory in MB because \"%s\" isn't a custom compute offering", computeOffering.Name)
	}

	if rootVolumeSizeInGb, ok := d.GetOk("root_volume_size_in_gb"); ok {
		instanceToCreate.RootVolumeSizeInGb = rootVolumeSizeInGb.(int)
	}

	if dedicatedGroupID, ok := d.GetOk("dedicated_group_id"); ok {
		instanceToCreate.DedicatedGroupId = dedicatedGroupID.(string)
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
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	// Get the virtual machine details
	instance, err := ccaResources.Instances.Get(d.Id())
	if err != nil {
		return handleNotFoundError("Instance", false, err, d)
	}
	// Update the config
	if err := d.Set("name", instance.Name); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	if err := setValueOrID(d, "template", strings.ToLower(instance.TemplateName), instance.TemplateId); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	if err := setValueOrID(d, "compute_offering", strings.ToLower(instance.ComputeOfferingName), instance.ComputeOfferingId); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	if err := d.Set("network_id", instance.NetworkId); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	if err := d.Set("private_ip_id", instance.IpAddressId); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	if err := d.Set("private_ip", instance.IpAddress); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	dID, dIDErr := getDedicatedGroupID(ccaResources, instance)
	if dIDErr != nil {
		return dIDErr
	}

	if err := d.Set("dedicated_group_id", dID); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	return nil
}

func resourceCloudcaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	d.Partial(true)

	if d.HasChange("compute_offering") || d.HasChange("cpu_count") || d.HasChange("memory_in_mb") {
		newComputeOffering := d.Get("compute_offering").(string)
		log.Printf("[DEBUG] Compute offering has changed for %s, changing compute offering...", newComputeOffering)
		newComputeOfferingID, ferr := retrieveComputeOfferingID(&ccaResources, newComputeOffering)
		if ferr != nil {
			return ferr
		}
		instanceToUpdate := cloudca.Instance{Id: d.Id(),
			ComputeOfferingId: newComputeOfferingID,
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

		computeOffering, cerr := ccaResources.ComputeOfferings.Get(newComputeOfferingID)
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

	if d.HasChange("private_ip") {
		return fmt.Errorf("Cannot update the private IP of an instance")
	}

	d.Partial(false)

	return nil
}

func resourceCloudcaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	fmt.Printf("[INFO] Destroying instance: %s\n", d.Get("name").(string))
	if _, err := ccaResources.Instances.Destroy(d.Id(), true); err != nil {
		return handleNotFoundError("Instance", true, err, d)
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
			return template.ID, nil
		}
	}

	return "", fmt.Errorf("Template with name %s not found", name)
}

func getDedicatedGroupID(ccaRes cloudca.Resources, instance *cloudca.Instance) (string, error) {
	dedicatedGroups, err := ccaRes.AffinityGroups.ListWithOptions(map[string]string{
		"type": "ExplicitDedication",
	})
	if err != nil {
		return "", err
	}
	for _, dedicatedGroup := range dedicatedGroups {
		for _, affinityGroupID := range instance.AffinityGroupIds {
			if strings.EqualFold(dedicatedGroup.Id, affinityGroupID) {
				return dedicatedGroup.Id, nil
			}
		}
	}
	return "", nil
}
