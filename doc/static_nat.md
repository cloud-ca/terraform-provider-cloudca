# cloudca_static_nat

Configures static NAT between a public and a private IP of an instance. Enabling static NAT is equivalent to forwarding every public port to every private port.

## Example Usage

```hcl
resource "cloudca_static_nat" "dev_static_nat" {
    environment_id = "4cad744d-bf1f-423d-887b-bbb34f4d1b5b"
    public_ip_id   = "10d523c1-907a-4f85-9181-9d62b16851c9"
    private_ip_id  = "c0d9824b-cb83-45ca-baca-e7e6c63a96a8"
}
```

## Argument Reference

The following arguments are supported:

- [environment_id](#environment_id) - (Required) ID of environment
- [public_ip_id](#public_ip_id) - (Required) The public IP to configure static NAT on. Cannot have any other purpose (e.g. load balancing, port forwarding)
- [private_ip_id](#private_ip_id) - (Required) A private IP of the instance to configure static NAT on. Must be in the same VPC as the public IP. Secondary IPs can be used here
