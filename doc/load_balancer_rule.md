# cloudca_load_balancer_rule

Manage load balancer rules. Modifying the ports or public IP will cause the rule to be recreated

## Example Usage

```hcl
resource "cloudca_load_balancer_rule" "lbr" {
    environment_id    = "4cad744d-bf1f-423d-887b-bbb34f4d1b5b"
    name              = "web_lb"
    network_id        = "7bb97867-8021-443b-b548-c15897e3816d"
    public_ip_id      = "5cd3a059-f15b-49f7-b7e1-254fef15968d"
    protocol          = "tcp"
    algorithm         = "leastconn"
    public_port       = 80
    private_port      = 80
    instance_ids      = ["071e2929-672e-45bc-a5b6-703d17c08367"]
    stickiness_method = "AppCookie"
    stickiness_params = {
        cookieName = "allo"
    }
}
```

## Argument Reference

The following arguments are supported:

- [environment_id](#environment_id) - (Required) ID of environment
- [name](#name) - (Required) Name of the load balancer rule
- [network_id](#network_id) - (Required) Id of the load balancing network to bind to
- [public_ip_id](#public_ip_id) - (Required) The id of the public IP to load balance on
- [protocol](#protocol) - (Required) The protocol to load balance
- [algorithm](#algorithm) - (Required) The algorithm to use for load balancing. Supports: "leastconn", "roundrobin" or "source"
- [instance_ids](#instance_ids) - (Optional) The list of instances to load balance
- [stickiness_method](#stickiness_method) - (Optional) The stickiness method to use. Supports : "LbCookie", "AppCookie" and "SourceBased"
- [stickiness_params](#stickiness_params) - (Optional) The additional parameters required for each stickiness method. See (TODO ADD LINK here) for more information

## Attribute Reference

In addition to the arguments listed above, the following computed attributes are returned:

- [id](#id) - the load balancer rule ID

## Import

Load balancer rules can be imported using the load balancer rule id, e.g.

```bash
terraform import cloudca_load_balancer_rule.lbr e798936b-b05d-4dbf-ade1-21f98c5fd0f0
```
