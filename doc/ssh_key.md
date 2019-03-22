# cloudca_ssh_key

Adds an SSH key to the environment so that it can be associated with instances.

## Example Usage

```hcl
resource "cloudca_ssh_key" "dev_ssh_key" {
    environment_id = "4cad744d-bf1f-423d-887b-bbb34f4d1b5b"
    name           = "my-ssh-key"
    public_key     = "my_public_key_data"
}
```

## Argument Reference

The following arguments are supported:

- [environment_id](#environment_id) - (Required) ID of environment
- [name](#name) - (Required) The name of the SSH key to add
- [public_key](#public_key) - (Required) The public key data

## Attribute Reference

Only the arguments listed above are returned.

## Import

SSH keys can be imported using the SSH key id, e.g.

```bash
terraform import cloudca_ssh_key.dev_ssh_key 919dd040-2b1e-4192-b25f-e3b8beca96e1
```
