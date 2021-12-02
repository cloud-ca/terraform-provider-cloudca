# cloudca_public_ip

Acquires a public IP in a specific VPC. If you update any of the fields in the resource, then it will release this IP and recreate it.

## Example Usage

```hcl
resource "cloudca_public_ip" "my_publicip" {
    environment_id = "4cad744d-bf1f-423d-887b-bbb34f4d1b5b"
    vpc_id         = "8b46e2d1-bbc4-4fad-b3bd-1b25fcba4cec"
}
```

## Argument Reference

The following arguments are supported:

- [environment_id](#environment_id) - (Required) ID of environment
- [vpc_id](#vpc_id) - (Required) The ID of the VPC to acquire the public IP

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are returned:

- [id](#id) - The public IP ID.
- [ip_address](#ip_address) - The public IP address

## Import

Public IPs can be imported using the public IP id, e.g.

```bash
terraform import cloudca_public_ip.my_publicip 56fd2565-edc9-444c-994d-9b7c46435d68
```
