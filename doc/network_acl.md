# cloudca_network_acl

Create a network ACL.

## Example Usage

```hcl
resource "cloudca_network_acl" "my_acl" {
    environment_id = "4cad744d-bf1f-423d-887b-bbb34f4d1b5b"
    name           = "test-acl"
    description    = "This is a test acl"
    vpc_id         = "8b46e2d1-bbc4-4fad-b3bd-1b25fcba4cec"
}
```

## Argument Reference

The following arguments are supported:

- [environment_id](#environment_id) - (Required) ID of environment
- [name](#name) - (Required) Name of the network ACL
- [description](#description) - (Required) Description of the network ACL
- [vpc_id](#vpc_id) - (Required) ID of the VPC where the network ACL should be created

## Attribute Reference

The following attributes are returned:

- [id](#id) - ID of network ACL.
- [name](#name) - Name of network ACL.
