# cloudca_vpn

Associate a Remote Access VPN with a VPC in order to enable VPN connectivity to a VPC from a client workstation.

## Example Usage

```hcl
resource "cloudca_vpn" "my_vpn" {
    environment_id = "4cad744d-bf1f-423d-887b-bbb34f4d1b5b"
    vpc_id         = "8b46e2d1-bbc4-4fad-b3bd-1b25fcba4cec"
}
```

## Argument Reference

The following arguments are supported:

- [environment_id](#environment_id) - (Required) ID of environment.
- [vpc_id](#vpc_id) - (Required) The ID of the VPC to associate the VPN with.

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are returned:

- [id](#id) - The VPN ID.
- [certificate](#certificate) - The certificate associated with this VPN connection (will be empty if `preshared_key` is set).
- [preshared_key](#preshared_key) - The pre-shared key associated with this VPN connection (will be empty if `certificate` is set).
- [public_ip](#public_ip) - The public IP address associated with the VPN.
- [public_ip_id](#public_ip_id) - The ID of the public IP associated with the VPN.
- [state](#state) - The state of the VPN connection.
- [type](#type) - The type of VPN connection (`IPSEC` or `IKEV2`).

## Import

VPNs can be imported using the VPN ID, e.g.

```bash
terraform import cloudca_vpn.my_vpn 56fd2565-edc9-444c-994d-9b7c46435d68
```
