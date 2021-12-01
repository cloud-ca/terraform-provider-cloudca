# cloudca_vpn_user

Associate a User with the Remote Access VPN in an environment.

## Example Usage

```hcl
resource "cloudca_vpn_user" "my_vpn_user" {
    environment_id = "4cad744d-bf1f-423d-887b-bbb34f4d1b5b"
    username       = "my_user"
    password       = "password"
}
```

## Argument Reference

The following arguments are supported:

- [environment_id](#environment_id) - (Required) ID of environment.
- [username](#username) - (Required) The username to create for VPN access.
- [password](#password) - (Required) The password of the created VPN user.

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are returned:

- [id](#id) - The VPN User ID.

## Import

VPN Users can be imported using the VPN User ID, e.g.

```bash
terraform import cloudca_vpn_user.my_vpn_user 56fd2565-edc9-444c-994d-9b7c46435d68
```
