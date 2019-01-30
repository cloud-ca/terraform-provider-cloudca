# cloudca_network

Create a network.

## Example Usage

```hcl
resource "cloudca_network" "my_network" {
    environment_id   = "4cad744d-bf1f-423d-887b-bbb34f4d1b5b"
    name             = "test-network"
    description      = "This is a test network"
    vpc_id           = "8b46e2d1-bbc4-4fad-b3bd-1b25fcba4cec"
    network_offering = "Standard Tier"
    network_acl      = "7d428416-263d-47cd-9270-2cdbdf222f57"
}
```

## Argument Reference

The following arguments are supported:

- [environment_id](#environment_id) - (Required) ID of environment
- [name](#name) - (Required) Name of the network
- [description](#description) - (Required) Description of the network
- [vpc_id](#vpc_id) - (Required) The ID of the vpc where the network should be created
- [network_offering](#network_offering) - (Required) The name of the network offering to use for the network
- [network_acl](#network_acl) - (Required) The id or name of the network ACL to use for the network

## Attribute Reference

The following attributes are returned:

- [id](#id) - ID of network.
- [cidr](#cidr) - Cidr of network
