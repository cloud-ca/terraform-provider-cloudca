# cloudca_vpc

Create a vpc.

## Example Usage

```hcl
resource "cloudca_vpc" "my_vpc" {
    environment_id = "4cad744d-bf1f-423d-887b-bbb34f4d1b5b"
    name           = "test-vpc"
    description    = "This is a test vpc"
    vpc_offering   = "Default VPC offering"
}
```

## Argument Reference

The following arguments are supported:

- [environment_id](#environment_id) - (Required) ID of environment
- [name](#name) - (Required) Name of the VPC
- [description](#description) - (Required) Description of the VPC
- [vpc_offering](#vpc_offering) - (Required) The name of the VPC offering to use for the vpc
- [network_domain](#network_domain) - (Optional) A custom DNS suffix at the level of a network
- [zone](#zone) - (Optional) The zone name or ID where the VPC will be created

## Attribute Reference

The following attributes are returned:

- [id](#id) - ID of VPC.

## Import

VPCs can be imported using the VPC id, e.g.

```bash
terraform import cloudca_vpc.my_vpc 06dca131-8c68-4054-bd6b-9e47c5a099ea
```
