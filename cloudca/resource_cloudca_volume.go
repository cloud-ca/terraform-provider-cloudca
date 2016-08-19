package cloudca

import(
  "fmt"
  "github.com/cloud-ca/go-cloudca"
  "github.com/cloud-ca/go-cloudca/api"
  "github.com/cloud-ca/go-cloudca/services/cloudca"
  "github.com/hashicorp/terraform/helper/schema"
  "log"
  "strings"
)

func resourceCloudcaVolume() *schema.Resource {
  return &schema.Resource {
    Create : resourceCloudcaVolumeCreate,
    Read : resourceCloudcaVolumeRead,
    Update : resourceCloudcaVolumeUpdate,
    Delete : resourceCloudcaVolumeDelete,

    Schema: map[string]*schema.Schema {
      "service_code": &schema.Schema {
        Type:        schema.TypeString,
        Required:    true,
        ForceNew:    true,
        Description: "A cloudca service code",
      },
      "environment_name": &schema.Schema {
        Type:        schema.TypeString,
        Required:    true,
        ForceNew:    true,
        Description: "Name of environment where port forwarding rule should be created",
      },
      "name" : &schema.Schema {
        Type:        schema.TypeString,
        Required:    true,
        ForceNew:    true,
        Description: "The name of the volume to be created",
      },
      "zone_name": &schema.Schema {
        Type:        schema.TypeString,
        Optional:    true,
        ForceNew:    true,
        Description: "The zone into which the volume will be create",
      },
      "storage_tier" : &schema.Schema {
        Type:        schema.TypeString,
        Required:    true,
        ForceNew:    true,
        Description: "The storage tier name",
      },
      "size" : &schema.Schema {
        Type:        schema.TypeString,
        Required:    true,
        ForceNew:    true,
        Description: "The size of the disk volume in gigabytes",
      },
      "instance_id" : &schema.Schema {
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
  size := d.Get("size").(string)
  diskOfferingId, err := retrieveDiskOfferingId(&ccaResources,storageTier,size)
  if(err != nil) {
    return err
  }
  volumeToCreate := cloudca.Volume{
    Name: d.Get("name").(string),
    DiskOfferingId : diskOfferingId,
  }

  if instanceId, ok := d.GetOk("instance_id"); ok {
    volumeToCreate.InstanceId = instanceId.(string)
  }
  if zoneName, ok := d.GetOk("zone_name"); ok {
    volumeToCreate.ZoneName = zoneName.(string)
  }

  newVolume, err := ccaResources.Volumes.Create(volumeToCreate)
  if(err != nil) {
    return fmt.Errorf("Error creating the new volume %s", err)
  }
  d.SetId(newVolume.Id)
  return resourceCloudcaPublicIpRead(d, meta)
}

func resourceCloudcaVolumeRead(d *schema.ResourceData, meta interface{}) error {
  ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

  volume, err := ccaResources.Volumes.Get(d.Id())
  if(err != nil) {
    return handleVolumeNotFoundError(err, d)
  }
  d.Set("name", volume.Name)
  d.Set("zone_name", volume.ZoneName)
  d.Set("storage_tier", volume.StorageTier)
  d.Set("size", volume.Size)
  d.Set("instance_id", volume.InstanceId)
  return nil
}

func resourceCloudcaVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
  ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

  d.Partial(true)

  if d.HasChange("instanceId") {
    oldInstanceId,newInstanceId := d.GetChange("instance_id")
    if volume, err := ccaResources.Volumes.Get(d.Id()); err != nil {
      if nerr := handleVolumeNotFoundError(err, d); nerr != nil {
        return nerr
      }
      if (oldInstanceId == nil) {
        log.Printf("[DEBUG] Instance Id %s detected. Attaching volume to instance.", newInstanceId)
        if aerr := ccaResources.Volumes.AttachToInstance(volume, newInstanceId.(string)); aerr != nil {
          return aerr
        }

      } else {
        log.Printf("[DEBUG] Instance Id has changed from %s, to %s. Attempting attach volume to new instance", oldInstanceId, newInstanceId)
        if derr := ccaResources.Volumes.DetachFromInstance(volume); derr != nil {
          return derr
        }
        if aerr := ccaResources.Volumes.AttachToInstance(volume, newInstanceId.(string)); aerr != nil {
          return aerr
        }
      }
      d.SetPartial("instance_id")
    }
  }
  d.Partial(false)
  return nil
}

func resourceCloudcaVolumeDelete(d *schema.ResourceData, meta interface{}) error {
  ccaClient := meta.(*cca.CcaClient)
	resources, _ := ccaClient.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	ccaResources := resources.(cloudca.Resources)

  fmt.Println("[INFO] Deleting volume: %s", d.Get("name").(string))
  if derr := ccaResources.Volumes.Delete(d.Id()); derr != nil {
    return handleVolumeNotFoundError(derr, d)
  }
  return nil
}


func retrieveDiskOfferingId(ccaResources *cloudca.Resources, storageTier string, size string) (id string, err error) {
  diskOfferings, err := ccaResources.DiskOfferings.List()
  if(err != nil) {
    return "", err
  }
  for _,diskOffering := range diskOfferings {
    if(strings.EqualFold(diskOffering.StorageTier, storageTier) && strings.EqualFold(diskOffering.Name,size)) {
      return diskOffering.Id, nil
    }
  }
  return "", fmt.Errorf("No valid disk offering's were found with storage tier: %s and size: %s", storageTier, size)
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
