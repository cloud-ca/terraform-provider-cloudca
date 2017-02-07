# Resources
- [**cloudca_environment**](#cloudca_environment)
- [**cloudca_vpc**](#cloudca_vpc)
- [**cloudca_tier**](#cloudca_tier)
- [**cloudca_network_acl**](#cloudca_network_acl)
- [**cloudca_network_acl_rule**](#cloudca_network_acl_rule)
- [**cloudca_instance**](#cloudca_instance)
- [**cloudca_public_ip**](#cloudca_public_ip)
- [**cloudca_static_nat**](#cloudca_static_nat)
- [**cloudca_port_forwarding_rule**](#cloudca_port_forwarding_rule)
- [**cloudca_volume**](#cloudca_volume)
- [**cloudca_load_balancer_rule**](#cloudca_load_balancer_rule)

## cloudca_environment
Manages a cloud.ca environment

### Example usage
```hcl
resource "cloudca_environment" "my_environment" {
	service_code = "compute-qc"
	organization_code = "test"
	name = "production"
	description = "Environment for production workloads"
	admin_role = ["pat"]
	read_only_role = ["franz","bob"]
}
```
### Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- organization_code - (Required) Organization's entry point, i.e. \<entry_point\>.cloud.ca
- name - (Required) Name of environment to be created. Must be lower case, contain alphanumeric characters, underscores or dashes
- description - (Required) Description for the environment
- admin_role - (Optional) List of users that will be given the Environment Admin role
- user_role - (Optional) List of users that will be given the User role
- read_only_role - (Optional) List of users that will be given the Read-only role

### Attribute Reference
- id - ID of the environment.
- name - Name of the environment.

## cloudca_vpc
Create a vpc.

### Example usage
```hcl
resource "cloudca_vpc" "my_vpc" {
	service_code = "compute-qc"
	environment_name = "dev"
	name = "test-vpc"
	description = "This is a test vpc"
	vpc_offering = "Default VPC offering"
}
```
### Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- environment_name - (Required) Name of environment
- name - (Required) Name of the VPC
- description - (Required) Description of the VPC
- vpc_offering - (Required) The name of the VPC offering to use for the vpc
- network_domain - (Optional) A custom DNS suffix at the level of a network
- zone - (Optional) The zone name or ID where the VPC will be created

### Attribute Reference
- id - ID of VPC.

## cloudca_tier
Create a tier.

### Example usage
```hcl
resource "cloudca_tier" "my_tier" {
	service_code = "compute-qc"
	environment_name = "dev"
	name = "test-tier"
	description = "This is a test tier"
	vpc_id = "8b46e2d1-bbc4-4fad-b3bd-1b25fcba4cec"
	network_offering = "Standard Tier"
	network_acl_id = "7d428416-263d-47cd-9270-2cdbdf222f57"
}
```
### Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- environment_name - (Required) Name of environment
- name - (Required) Name of the tier
- description - (Required) Description of the tier
- vpc_id - (Required) The ID of the vpc where the tier should be created
- network_offering - (Required) The name of the network offering to use for the tier
- network_acl_id - (Required) The id of the network ACL to use for the tier

### Attribute Reference
- id - ID of tier.
- cidr - Cidr of tier

## cloudca_network_acl
Create a network ACL.

### Example usage
```hcl
resource "cloudca_network_acl" "my_acl" {
	service_code = "compute-qc"
	environment_name = "dev"
	name = "test-acl"
	description = "This is a test acl"
	vpc_id = "8b46e2d1-bbc4-4fad-b3bd-1b25fcba4cec"
}
```
###Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- environment_name - (Required) Name of environment
- name - (Required) Name of the network ACL
- description - (Required) Description of the network ACL
- vpc_id - (Required) ID of the VPC where the network ACL should be created

### Attribute Reference
- id - ID of network ACL.
- name - Name of network ACL.

## cloudca_network_acl_rule
Create a network ACL rule.

### Example usage
```hcl
resource "cloudca_network_acl_rule" "my_acl" {
	service_code = "compute-qc"
	environment_name = "dev"
	rule_number = 55
	cidr = "10.212.208.0/22"
	action = "Allow"
	protocol = "TCP"
	start_port = 80
	end_port = 80
	traffic_type = "Ingress"
	network_acl_id = "c0731f8b-92f0-4fac-9cbd-245468955fdf"
}
```
### Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- environment_name - (Required) Name of environment
- network_acl_id - (Required) ID of the network ACL where the rule should be created
- rule_number - (Required) Rule number of the network ACL rule
- cidr - (Required) CIDR of the network ACL rule
- action - (Required) Action of the network ACL rule (i.e. Allow or Deny)
- protocol - (Required) Protocol of the network ACL rule (i.e. TCP, UDP, ICMP or All)
- traffic_type - (Required) TrafficType of the network ACL rule (i.e. Ingress or Egress)
- icmp_type - (Optional) The ICMP type. Can only be used with ICMP protocol
- icmp_code - (Optional) The ICMP code. Can only be used with ICMP protocol
- start_port - (Optional) The start port. Can only be used with TCP/UDP protocol
- end_port - (Optional) The end port. Can only be used with TCP/UDP protocol

### Attribute Reference
- id - ID of network ACL.
- name - Name of network ACL.

## cloudca_instance
Create and starts an instance.

### Example usage
```hcl
resource "cloudca_instance" "my_instance" {
	service_code = "compute-qc"
	environment_name = "dev"
	name = "test-instance"
	network_id = "672016ef-05ee-4e88-b68f-ac9cc462300b"
	template = "CentOS 6.7 base (64bit)"
	compute_offering = "1vCPU.512MB"
	ssh_key_name = "my_ssh_key"
}
```
### Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- environment_name - (Required) Name of environment
- name - (Required) Name of instance
- network_id - (Required) The ID of the network where the instance should be created
- template - (Required) Name of template to use for the instance
- compute_offering - (Required) Name of the compute offering to use for the instance
- cpu_count - (Required if custom compute offering) Number of CPUs the instance should be created with. Can only be used for custom compute offerings.
- memory_in_mb - (Required if custom compute offering) Amount of memory in MB the instance should be created with. Can only be used for custom compute offerings.
- user_data - (Optional) User data to add to the instance
- ssh_key_name - (Optional) Name of the SSH key pair to attach to the instance. Mutually exclusive with public_key.
- public_key - (Optional) Public key to attach to the instance. Mutually exclusive with ssh_key_name.

### Attribute Reference
- id - ID of instance.
- private_ip_id - ID of instance's private IP
- private_ip - Instance's private IP

## cloudca_public_ip
Acquires a public IP in a specific VPC. If you update any of the fields in the resource, then it will release this IP and recreate it.

### Example usage
```hcl
resource "cloudca_public_ip" "my_publicip" {
	service_code = "compute-qc"
	environment_name = "dev"
	vpc_id = "8b46e2d1-bbc4-4fad-b3bd-1b25fcba4cec"
}
```
### Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- environment_name - (Required) Name of environment
- vpc_id - (Required) The ID of the VPC to acquire the public IP

### Attribute Reference
- id - The public IP ID.
- ip_address - The public IP address

## cloudca_static_nat
Configures static NAT between a public and a private IP of an instance. Enabling static NAT is equivalent to forwarding every public port to every private port.

### Example usage
```hcl
resource "cloudca_static_nat" "dev_static_nat" {
	service_code     = "compute-qc"
	environment_name = "dev"

	public_ip_id     = "10d523c1-907a-4f85-9181-9d62b16851c9"
	private_ip_id    = "c0d9824b-cb83-45ca-baca-e7e6c63a96a8"
}
```

### Argument reference
- service_code - (Required)
- environment_name - (Required)
- public_ip_id - (Required) The public IP to configure static NAT on. Cannot have any other purpose (e.g. load balancing, port forwarding)
- private_ip_id - (Required) A private IP of the instance to configure static NAT on. Must be in the same VPC as the public IP. Secondary IPs can be used here

## cloudca_port_forwarding_rule
Manages port forwarding rules. Modifying any field will result in destruction and recreation of the rule.

When adding a port forwarding rule to the default private IP of an instance, only the instance id is required. Alternatively, the private_ip_id can be used on its own (for example when targeting an instance secondary IP).

### Example usage
```hcl
resource "cloudca_port_forwarding_rule" "web_pfr" {
	service_code = "compute-qc"
	environment_name = "dev"

	public_ip_id = "319f508f-089b-482d-af17-0f3360520c69"
	public_port_start = 80
	private_ip_id = "30face92-f1cf-4064-aa7f-008ea09ef7f0"
	private_port_start = 8080
	protocol = "TCP"
}
```

### Argument reference
- service_code - (Required)
- environment_name - (Required)
- private_ip_id - (Required) The private IP which should be used to create this rule
- private_port_start - (Required)
- private_port_end - (Optional) If not specified, defaults to the private start port
- public_ip_id - (Required) The public IP which should be used to create this rule
- public_port_start - (Required)
- public_port_end - (Optional) If not specified, defaults to the public start port
- protocol - (Required) The protocol to be used for this rule - must be TCP or UDP

### Attribute reference
- id - the rule ID
- public_ip - the public IP address of this rule
- private_ip - the private IP address of this rule
- instance_id - the instance associated with the private IP address of this rule

## cloudca_volume
Manages volumes. Modifying all fields with the exception of instance_id will result in destruction and recreation of the volume.

If the instance_id is updated, the volume will be detached from the previous instance and attached to the new instance.

### Example usage
```hcl
resource "cloudca_volume" "data_volume" {
	service_code = "compute-qc"
	environment_name = "dev"

	name = "Data Volume"

	disk_offering = "High Performance SSD"
	size_in_gb = 20
	iops = 2000
	instance_id = "f932c530-5753-44ce-8aae-263672e1ae3f"
}
```

### Argument reference
- service_code - (Required)
- environment_name - (Required)
- name - (Required) The name of the volume to be created
- disk_offering - (Required) The name or id of the disk offering to use for the volume
- size_in_gb - (Optional) The size in GB of the volume. Only for disk offerings with custom size.
- iops - (Optional) The number of IOPS of the volume. Only for disk offerings with custom iops.
- instance_id - The instance ID that the volume will be attached to. Note that changing the instance ID will _not_ result in the destruction of this volume

### Attribute reference
- id - the volume ID


## cloudca_load_balancer_rule
Manage load balancer rules. Modifying the ports or public IP will cause the rule to be recreated

### Example usage
```hcl
resource "cloudca_load_balancer_rule" "lbr" {
   service_code = "compute-qc"
   environment_name = "dev"
   name="web_lb"
   network_id  = "7bb97867-8021-443b-b548-c15897e3816d"
   public_ip_id="5cd3a059-f15b-49f7-b7e1-254fef15968d"
   protocol="tcp"
   algorithm = "leastconn"
   public_port = 80
   private_port = 80
   instance_ids = ["071e2929-672e-45bc-a5b6-703d17c08367"]
   stickiness_method = "AppCookie"
   stickiness_params {
      cookieName = "allo"
   }
}
```

### Argument reference
- service_code - (Required)
- environment_name - (Required)
- name - (Required) Name of the load balancer rule
- network_id - (Required) Id of the load balancing tier to bind to
- public_ip_id - (Required) The id of the public IP to load balance on
- protocol - (Required) The protocol to load balance
- algorithm - (Required) The algorithm to use for load balancing. Supports: "leastconn", "roundrobin" or "source"
- instance_ids - (Optional) The list of instances to load balance
- stickiness_method - (Optional) The stickiness method to use. Supports : "LbCookie", "AppCookie" and "SourceBased"
- stickiness_params - (Optional) The additional parameters required for each stickiness method. See (TODO ADD LINK here) for more information

### Attribute reference
- id - the load balancer rule ID
