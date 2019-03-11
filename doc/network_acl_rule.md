# cloudca_network_acl_rule

Create a network ACL rule.

## Example Usage

```hcl
resource "cloudca_network_acl_rule" "my_acl" {
    environment_id = "4cad744d-bf1f-423d-887b-bbb34f4d1b5b"
    rule_number    = 55
    cidr           = "10.212.208.0/22"
    action         = "Allow"
    protocol       = "TCP"
    start_port     = 80
    end_port       = 80
    traffic_type   = "Ingress"
    network_acl_id = "c0731f8b-92f0-4fac-9cbd-245468955fdf"
}
```

## Argument Reference

The following arguments are supported:

- [environment_id](#environment_id) - (Required) ID of environment
- [network_acl_id](#network_acl_id) - (Required) ID of the network ACL where the rule should be created
- [rule_number](#rule_number) - (Required) Rule number of the network ACL rule
- [cidr](#cidr) - (Required) CIDR of the network ACL rule
- [action](#action) - (Required) Action of the network ACL rule (i.e. Allow or Deny)
- [protocol](#protocol) - (Required) Protocol of the network ACL rule (i.e. TCP, UDP, ICMP or All)
- [traffic_type](#traffic_type) - (Required) TrafficType of the network ACL rule (i.e. Ingress or Egress)
- [icmp_type](#icmp_type) - (Optional) The ICMP type. Can only be used with ICMP protocol
- [icmp_code](#icmp_code) - (Optional) The ICMP code. Can only be used with ICMP protocol
- [start_port](#start_port) - (Optional) The start port. Can only be used with TCP/UDP protocol
- [end_port](#end_port) - (Optional) The end port. Can only be used with TCP/UDP protocol

## Attribute Reference

The following attributes are returned:

- [id](#id) - ID of network ACL.
- [name](#name) - Name of network ACL.

## Import

Network ACL rules can be imported using the network ACL rule id, e.g.

```bash
terraform import cloudca_network_acl_rule.my_acl 24323470-336e-4244-be26-5b25a262bcce
```
