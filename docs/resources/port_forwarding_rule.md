# cloudca_port_forwarding_rule

Manages port forwarding rules. Modifying any field will result in destruction and recreation of the rule.

When adding a port forwarding rule to the default private IP of an instance, only the instance id is required. Alternatively, the private_ip_id can be used on its own (for example when targeting an instance secondary IP).

## Example Usage

```hcl
resource "cloudca_port_forwarding_rule" "web_pfr" {
    environment_id     = "4cad744d-bf1f-423d-887b-bbb34f4d1b5b"
    public_ip_id       = "319f508f-089b-482d-af17-0f3360520c69"
    public_port_start  = 80
    private_ip_id      = "30face92-f1cf-4064-aa7f-008ea09ef7f0"
    private_port_start = 8080
    protocol           = "TCP"
}
```

## Argument Reference

The following arguments are supported:

- [environment_id](#environment_id) - (Required) ID of environment_id
- [private_ip_id](#private_ip_id) - (Required) The private IP which should be used to create this rule
- [private_port_start](#private_port_start) - (Required)
- [private_port_end](#private_port_end) - (Optional) If not specified, defaults to the private start port
- [public_ip_id](#public_ip_id) - (Required) The public IP which should be used to create this rule
- [public_port_start](#public_port_start) - (Required)
- [public_port_end](#public_port_end) - (Optional) If not specified, defaults to the public start port
- [protocol](#protocol) - (Required) The protocol to be used for this rule - must be TCP or UDP

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are returned:

- [id](#id) - the rule ID
- [public_ip](#public_ip) - the public IP address of this rule
- [private_ip](#private_ip) - the private IP address of this rule
- [instance_id](#instance_id) - the instance associated with the private IP address of this rule

## Import

Port forwarding rules can be imported using the port forwarding rules id, e.g.

```bash
terraform import cloudca_port_forwarding_rule.web_pfr 816bd39d-5379-45be-b7a1-6b2ea18cec62
```
