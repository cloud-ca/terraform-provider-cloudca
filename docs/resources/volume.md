---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudca_volume Resource - terraform-provider-cloudca"
subcategory: ""
description: |-
  
---

# cloudca_volume (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **disk_offering** (String) The ID or name of the disk offering of the new volume
- **environment_id** (String) ID of environment where the volume should be created
- **instance_id** (String) The id of the instance to which the volume will be attached
- **name** (String) The name of the volume to be created

### Optional

- **id** (String) The ID of this resource.
- **iops** (Number) The number of iops of the volume
- **size_in_gb** (Number) The size of the volume in gigabytes

