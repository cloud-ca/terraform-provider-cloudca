# cloud.ca Provider

The cloud.ca provider is used to interact with the many resources supported by [cloud.ca](https://cloud.ca/). The provider needs to be configured with the proper credentials before it can be used. Optionally with a URL pointing to a running cloud.ca API.

In order to provide the required configuration options you need to supply the value for `api_key` field.

## Example Usage

```hcl
variable "my_api_key" {}

# Configure cloud.ca Provider
provider "cloudca" {
    api_key = "${var.my_api_key}"
}

# Create an Instance
resource "cloudca_instance" "instance" {
    # ...
}
```

## Argument Reference

The following arguments are supported:

- [api_key](#api_key) - (Required) This is the cloud.ca API key. It can also be sourced from the `CLOUDCA_API_KEY` environment variable.
- [api_url](#api_url) - (Optional) This is the cloud.ca API URL. It can also be sourced from the `CLOUDCA_API_URL` environment variable.

## Resources

- [**cloudca_environment**](environment.md)
- [**cloudca_instance**](instance.md)
- [**cloudca_load_balancer_rule**](load_balancer_rule.md)
- [**cloudca_network**](network.md)
- [**cloudca_network_acl**](network_acl.md)
- [**cloudca_network_acl_rule**](network_acl_rule.md)
- [**cloudca_port_forwarding_rule**](port_forwarding_rule.md)
- [**cloudca_public_ip**](public_ip.md)
- [**cloudca_static_nat**](static_nat.md)
- [**cloudca_ssh_key**](ssh_key.md)
- [**cloudca_volume**](volume.md)
- [**cloudca_vpc**](vpc.md)
