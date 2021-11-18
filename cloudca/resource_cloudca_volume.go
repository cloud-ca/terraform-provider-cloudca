package cloudca

import (
	"fmt"
	"log"
	"strings"

	cca "github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudcaVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudcaVolumeCreate,
		Read:   resourceCloudcaVolumeRead,
		Update: resourceCloudcaVolumeUpdate,
		Delete: resourceCloudcaVolumeDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of environment where the volume should be created",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the volume to be created",
			},
			"disk_offering": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID or name of the disk offering of the new volume",
				StateFunc: func(val interface{}) string {
					return strings.ToLower(val.(string))
				},
			},
			"size_in_gb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The size of the volume in gigabytes",
			},
			"iops": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of iops of the volume",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the instance to which the volume will be attached",
			},
		},
	}
}

func resourceCloudcaVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	diskOffering, err := retrieveDiskOffering(&ccaResources, d.Get("disk_offering").(string))
	if err != nil {
		return err
	}
	volumeToCreate := cloudca.Volume{
		Name:           d.Get("name").(string),
		DiskOfferingId: diskOffering.Id,
	}

	if val, ok := d.GetOk("size_in_gb"); ok {
		if !diskOffering.CustomSize {
			return fmt.Errorf("Disk offering %s doesn't allow custom size", diskOffering.Id)
		}
		volumeToCreate.GbSize = val.(int)
	}

	if val, ok := d.GetOk("iops"); ok {
		if !diskOffering.CustomIops {
			return fmt.Errorf("Disk offering %s doesn't allow custom IOPS", diskOffering.Id)
		}
		volumeToCreate.Iops = val.(int)
	}

	if zone, ok := d.GetOk("zone"); ok {
		if isID(zone.(string)) {
			volumeToCreate.ZoneId = zone.(string)
		} else {
			volumeToCreate.ZoneId, err = retrieveZoneID(&ccaResources, zone.(string))
			if err != nil {
				return err
			}
		}
	}

	if instanceID, ok := d.GetOk("instance_id"); ok {
		volumeToCreate.InstanceId = instanceID.(string)
	}

	newVolume, err := ccaResources.Volumes.Create(volumeToCreate)
	if err != nil {
		return err
	}
	d.SetId(newVolume.Id)
	return resourceCloudcaVolumeRead(d, meta)
}

func resourceCloudcaVolumeRead(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	volume, err := ccaResources.Volumes.Get(d.Id())
	if err != nil {
		return handleNotFoundError("Volume", false, err, d)
	}

	if err := d.Set("name", volume.Name); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	if err := setValueOrID(d, "disk_offering", strings.ToLower(volume.DiskOfferingName), volume.DiskOfferingId); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	if err := d.Set("size_in_gb", volume.GbSize); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	if err := d.Set("iops", volume.Iops); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	if err := d.Set("instance_id", volume.InstanceId); err != nil {
		return fmt.Errorf("Error reading Trigger: %s", err)
	}

	return nil
}

func resourceCloudcaVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	d.Partial(true)
	curVolume, err := ccaResources.Volumes.Get(d.Id())
	if err != nil {
		return err
	}
	if d.HasChange("instance_id") {
		oldInstanceID, newInstanceID := d.GetChange("instance_id")
		volume := &cloudca.Volume{
			Id: d.Id(),
		}
		if oldInstanceID != "" && curVolume.InstanceId != "" {
			err := ccaResources.Volumes.DetachFromInstance(volume)
			if err != nil {
				return err
			}
		}
		if newInstanceID != "" {
			err := ccaResources.Volumes.AttachToInstance(volume, newInstanceID.(string))
			if err != nil {
				return err
			}
		}
		d.SetPartial("instance_id")
	}
	if d.HasChange("size_in_gb") || d.HasChange("iops") {
		volumeToResize := cloudca.Volume{
			Id: d.Id(),
		}
		if val, ok := d.GetOk("size_in_gb"); ok {
			volumeToResize.GbSize = val.(int)
			if curVolume.GbSize > volumeToResize.GbSize {
				return fmt.Errorf("Cannot reduce size of a volume")
			}
		}
		if val, ok := d.GetOk("iops"); ok {
			volumeToResize.Iops = val.(int)
		}
		_ = ccaResources.Volumes.Resize(&volumeToResize)
		d.SetPartial("size_in_gb")
		d.SetPartial("iops")
	}
	d.Partial(false)
	return resourceCloudcaVolumeRead(d, meta)
}

func resourceCloudcaVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	ccaResources, rerr := getResourcesForEnvironmentID(meta.(*cca.CcaClient), d.Get("environment_id").(string))

	if rerr != nil {
		return rerr
	}
	if instanceID, ok := d.GetOk("instance_id"); ok && instanceID != "" {
		volume := &cloudca.Volume{
			Id: d.Id(),
		}
		err := ccaResources.Volumes.DetachFromInstance(volume)
		if err != nil {
			return err
		}
	}
	if err := ccaResources.Volumes.Delete(d.Id()); err != nil {
		return handleNotFoundError("Volume", true, err, d)
	}
	return nil
}

func retrieveZoneID(ccaResources *cloudca.Resources, zoneName string) (zoneID string, nerr error) {
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

func retrieveDiskOffering(ccaRes *cloudca.Resources, name string) (diskOffering *cloudca.DiskOffering, err error) {
	if isID(name) {
		return ccaRes.DiskOfferings.Get(name)
	}
	offerings, err := ccaRes.DiskOfferings.List()
	if err != nil {
		return nil, err
	}
	for _, offering := range offerings {
		if strings.EqualFold(offering.Name, name) {
			log.Printf("Found disk offering: %+v", offering)
			return &offering, nil
		}
	}
	return nil, fmt.Errorf("Disk offering with name %s not found", name)
}
