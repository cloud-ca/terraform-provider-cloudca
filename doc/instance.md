# cloudca_instance

Create and starts an instance.

## Example Usage

```hcl
resource "cloudca_instance" "my_instance" {
    environment_id         = "4cad744d-bf1f-423d-887b-bbb34f4d1b5b"
    name                   = "test-instance"
    network_id             = "672016ef-05ee-4e88-b68f-ac9cc462300b"
    template               = "Ubuntu 16.04.03 HVM"
    compute_offering       = "1vCPU.512MB"
    ssh_key_name           = "my_ssh_key"
    root_volume_size_in_gb = 100
    private_ip             = "10.2.1.124"
    dedicated_group_id      = "78fdce97-3a46-4b50-bca7-c70ef8449da8"
}
```

## Argument Reference

The following arguments are supported:

- [environment_id](#environment_id) - (Required) ID of environment
- [name](#name) - (Required) Name of instance
- [network_id](#network_id) - (Required) The ID of the network where the instance should be created
- [template](#template) - (Required) Name of template to use for the instance
- [compute_offering](#compute_offering) - (Required) Name of the compute offering to use for the instance
- [cpu_count](#cpu_count) - (Required) Number of CPUs the instance should be created with.
- [memory_in_mb](#memory_in_mb) - (Required) Amount of memory in MB the instance should be created with.
- [user_data](#user_data) - (Optional) User data to add to the instance
- [ssh_key_name](#ssh_key_name) - (Optional) Name of the SSH key pair to attach to the instance. Mutually exclusive with public_key.
- [public_key](#public_key) - (Optional) Public key to attach to the instance. Mutually exclusive with ssh_key_name.
- [root_volume_size_in_gb](#root_volume_size_in_gb) - (Optional) Size of the root volume of the instance. This only works for templates that allows root volume resize.
- [private_ip](#private_ip) - (Optional) Instance's private IPv4 address.
- [dedicated_group_id](#dedicated_group_id) - (Optional) Dedicated group id in which the instance will be created

## Attribute Reference

The following attributes are returned:

- [id](#id) - ID of instance.
- [private_ip_id](#private_ip_id) - ID of instance's private IP
- [private_ip](#private_ip) - Instance's private IP

## Import

Instances can be imported using the instance id, e.g.

```bash
terraform import cloudca_instance.my_instance c33dc4e3-0067-4c26-a588-53c9a936b9de
```
