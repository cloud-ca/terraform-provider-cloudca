---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cloudca_load_balancer_rule Resource - terraform-provider-cloudca"
subcategory: ""
description: |-
  
---

# cloudca_load_balancer_rule (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **algorithm** (String) The algorithm used to load balance
- **environment_id** (String) ID of environment where load balancer rule should be created
- **name** (String) Name of the load balancer rule
- **network_id** (String) The network ID to bind to
- **private_port** (String) The port to which the traffic will be load balanced internally
- **protocol** (String) The protocol that this rule should use (eg. TCP, UDP)
- **public_ip_id** (String) ID of the public IP to which the rule should be applied
- **public_port** (String) The port on the public IP

### Optional

- **id** (String) The ID of this resource.
- **instance_ids** (Set of String) List of instance ids that will be load balanced
- **stickiness_method** (String) The stickiness method
- **stickiness_params** (Map of String) The stickiness policy parameters

### Read-Only

- **public_ip** (String) The public IP to which the rule should be applied

