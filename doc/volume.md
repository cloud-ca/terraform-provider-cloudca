# cloudca_volume

Manages volumes. Modifying all fields with the exception of instance_id will result in destruction and recreation of the volume.

If the instance_id is updated, the volume will be detached from the previous instance and attached to the new instance.

**WARNING: Updating `size_in_gb` and/or `iops` of a volume will cause a REBOOT of the instance it's attached to.**

## Example Usage

```hcl
resource "cloudca_volume" "data_volume" {
    environment_id = "4cad744d-bf1f-423d-887b-bbb34f4d1b5b"
    name           = "Data Volume"
    disk_offering  = "20GB - 20 IOPS Min."
    instance_id    = "f932c530-5753-44ce-8aae-263672e1ae3f"
    size_in_gb     = "10"
}
```

## Argument Reference

The following arguments are supported:

- [environment_id](#environment_id) - (Required) ID of environment
- [name](#name) - (Required) The name of the volume to be created
- [disk_offering](#disk_offering) - (Required) The name or id of the disk offering to use for the volume
- [size_in_gb](#size_in_gb) - (Required) The size in GB of the volume.
- [iops](#iops) - (Optional) The number of IOPS of the volume. Only for disk offerings with custom iops.
- [instance_id](#instance_id) - The instance ID that the volume will be attached to. Note that changing the instance ID will _not_ result in the destruction of this volume

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are returned:

- [id](#id) - the volume ID

## Import

Volumes can be imported using the volume id, e.g.

```bash
terraform import cloudca_volume.data_volume b24f94f7-098f-458b-aeb3-b38992ae8d67
```
