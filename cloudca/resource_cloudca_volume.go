package cloudca

import (
	"fmt"
	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func resourceCloudcaVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaVolumeCreate,
		Read:   resourceCloudcaVolumeRead,
		Update: resourceCloudcaVolumeUpdate,
		Delete: resourceCloudcaVolumeDelete,

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
				Description: "Name of environment where port forwarding rule should be created",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the volume to be created",
			},
			"zone_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The zone into which the volume will be create",
			},
			"storage_tier": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The storage tier name",
			},
			"size_in_gb": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The size of the disk volume in gigabytes",
			},
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of the instance to which the volume should be attached",
			},
		},
	}
}

func resourceCloudcaVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	storageTier := d.Get("storage_tier").(string)
	size := d.Get("size_in_gb").(int)
	diskOfferingId, err := retrieveDiskOfferingId(&ccaResources, storageTier, size)
	if err != nil {
		return err
	}
	volumeToCreate := cloudca.Volume{
		Name:           d.Get("name").(string),
		DiskOfferingId: diskOfferingId,
	}

	if zoneName, ok := d.GetOk("zone_name"); ok {
		if zoneId, err := retrieveZoneId(&ccaResources, zoneName.(string)); err != nil {
			return err
		} else {
			volumeToCreate.ZoneId = zoneId
		}
	}

	newVolume, err := ccaResources.Volumes.Create(volumeToCreate)
	if err != nil {
		return err
	}
	if instanceId, ok := d.GetOk("instance_id"); ok {
		err = ccaResources.Volumes.AttachToInstance(newVolume, instanceId.(string))
		if err != nil {
			return err
		}
	}
	d.SetId(newVolume.Id)
	return resourceCloudcaVolumeRead(d, meta)
}

func resourceCloudcaVolumeRead(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	volume, err := ccaResources.Volumes.Get(d.Id())
	if err != nil {
		return handleVolumeNotFoundError(err, d)
	}
	d.Set("name", volume.Name)
	d.Set("zone_name", volume.ZoneName)
	d.Set("storage_tier", volume.StorageTier)
	d.Set("size_in_gb", volume.Size)
	d.Set("instance_id", volume.InstanceId)
	return nil
}

func resourceCloudcaVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

	d.Partial(true)
	if d.HasChange("instance_id") {
		oldInstanceId, newInstanceId := d.GetChange("instance_id")
		volume := &cloudca.Volume{
			Id: d.Id(),
		}
		if oldInstanceId != "" {
			err := ccaResources.Volumes.DetachFromInstance(volume)
			if err != nil {
				return fmt.Errorf("%s %s", oldInstanceId, newInstanceId)
			}
		}
		if newInstanceId != "" {
			err := ccaResources.Volumes.AttachToInstance(volume, newInstanceId.(string))
			if err != nil {
				return err
			}
		}
		d.SetPartial("instance_id")
	}
	d.Partial(false)
	return resourceCloudcaVolumeRead(d, meta)
}

func resourceCloudcaVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)
	if instanceId, ok := d.GetOk("instance_id"); ok && instanceId != "" {
		volume := &cloudca.Volume{
			Id: d.Id(),
		}
		err := ccaResources.Volumes.DetachFromInstance(volume)
		if err != nil {
			return err
		}
	}
	if err := ccaResources.Volumes.Delete(d.Id()); err != nil {
		return handleVolumeNotFoundError(err, d)
	}
	return nil
}

func retrieveDiskOfferingId(ccaResources *cloudca.Resources, storageTier string, size int) (id string, err error) {
	diskOfferings, err := ccaResources.DiskOfferings.List()
	if err != nil {
		return "", err
	}
	for _, diskOffering := range diskOfferings {
		if strings.EqualFold(diskOffering.StorageTier, storageTier) && diskOffering.GbSize == size {
			return diskOffering.Id, nil
		}
	}
	return "", fmt.Errorf("No valid disk offering's were found with storage tier: %s and size: %d", storageTier, size)
}

func retrieveZoneId(ccaResources *cloudca.Resources, zoneName string) (zoneId string, nerr error) {
	zones, err := ccaResources.Zones.List()
	if err != nil {
		return "", err
	}
	for _, zone := range zones {
		if strings.EqualFold(zone.Name, zoneName) {
			return zone.Id, nil
		}
	}
	return "", fmt.Errorf("Zone with name %s could not be found", zoneName)
}

func handleVolumeNotFoundError(err error, d *schema.ResourceData) error {
	if ccaError, ok := err.(api.CcaErrorResponse); ok {
		if ccaError.StatusCode == 404 {
			fmt.Errorf("Volume with id='%s' was not found", d.Id())
			d.SetId("")
			return nil
		}
	}
	return err
}
