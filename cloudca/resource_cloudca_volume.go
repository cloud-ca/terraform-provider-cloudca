package cloudca

import(
  "fmt"
  "github.com/cloud-ca/go-cloudca"
  "github.com/cloud-ca/go-cloudca/services/cloudca"
  "github.com/hashicorp/terraform/helper/schema"
  "strings"
)

func resourceCloudcaVolume() *schema.Resource {
  return &schema.Resource {
    Create : resourceCloudcaVolumeCreate,
    Read : resourceCloudcaVolumeRead,

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
        Optional:    true;
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
        Type:        schema.TypeInt,
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

  zoneId, err := retrieveZoneId(&ccaResources ,d.Get("zone_name").(string))
  if (err != nil) {
    return "", err
  }

  size := d.Get("size").(string)
  storageTier := d.Get("storage_tier").(string)
  diskOfferingId := retrieveDiskOfferingId(&ccaResources,storageTier,size)

  volumeToCreate := cloudca.Voldumes {
    Name: d.Get("name").(string),
    ZoneId : zoneId,
    DiskOfferingId : diskOfferingId,
    InstanceId : d.Get("instance_id").(string)
  }

  newVolume, err := ccaResources.Volumes.CreateVolume(volumeToCreate)
  if(err != null) {
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
    return handleNotFoundError(err, d)
  }
  d.Set("name", volume.Name)
  d.Set("zone_name", volume.ZoneName)
  d.Set("storage_tier", volume.StorageTier)
  d.Set("size", volume.Size)
  d.Set("instance_id", volume.InstanceId)
}

func retrieveDiskOfferingId(ccaResources *cloudca.Resources, storageTier string, diskSize int) (id string, err error) {
  return "", nil
}

func retrieveZoneId(ccaResources *cloudca.Resources, name string) (id string, err error) {
  //worth checking if isId first?
  tiers, err := ccaResources.Tiers.List()
  if(err != nil) {
    return "", err
  }
  for _, tier := range tiers {
    if strings.EqualFold(tier.Name, name) {
      return tier.Id, nil
    }
  }
  return "", fmt.Errorf("Tier with name %s not found", name)
}

func handleNotFoundError(err error, d *schema.ResourceData) error {
  if ccaError, ok := err.(api.CcaErrorResponse); ok {
    if ccaError.StatusCode == 404 {
      log.Printf("Volume with id='%s' was not found", d.Id())
      d.setId("")
      return nil
    }
  }
}
