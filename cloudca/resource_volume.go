package cloudca

import(
  "fmt"
  "github.com/cloud-ca/go-cloudca"
  "github.com/cloud-ca/go-cloudca/api"
  "github.com/cloud-ca/go-cloudca/services/cloudca"
  "strconv"
)

func resourceVolume() *schema.Resource {
  return &schema.Resource {
    Create: createVolume,
    Read: readVolume,
    Update: updateVolume,
    Delete: deleteVolume,

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
        ForceNew:    true,
        Description: "The zone into which the volume will be create",
      },
      "storage_tier" : &schema.Schema {
        Type:        schema.TypeString,
        Required:    true,
        ForceNew:    true,
        Description: "The storage tier name"
      },
      "size" : &schema.Schema {
        Type:        schema.TypeInt,
        Required:    true,
        ForceNew:    true,
        Description: "The size of the disk volume in gigabytes"
      },
      "instance_id" : &schema.Schema {
        Type:        &schema.TypeString
        Description: "The id of the instance to which the volume should be attached."
      },
    }
  }
}
