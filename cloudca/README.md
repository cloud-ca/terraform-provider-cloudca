#Resources
- [**cloudca_environment**](#cloudca_environment)
- [**cloudca_vpc**](#cloudca_vpc)
- [**cloudca_tier**](#cloudca_tier)
- [**cloudca_network_acl**](#cloudca_network_acl)
- [**cloudca_network_acl_rule**](#cloudca_network_acl_rule)
- [**cloudca_instance**](#cloudca_instance)
- [**cloudca_public_ip**](#cloudca_public_ip)
- [**cloudca_port_forwarding_rule**](#cloudca_port_forwarding_rule)

##cloudca_environment
Manages a cloud.ca environment

###Example usage
```
resource "cloudca_environment" "my_environment" {
	service_code = "compute-east"
	organization_code = "test"
	name = "production"
	description = "Environment for production workloads"
	admin_role = ["pat"]
	read_only_role = ["franz","bob"]
}
```
###Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- organization_code - (Required) Organization's entry point, i.e. \<entry_point\>.cloud.ca
- name - (Required) Name of environment to be created. Must be lower case, contain alphanumeric charaters, underscores or dashes
- description - (Required) Description for the environment
- admin_role - (Optional) List of users that will be given the Environment Admin role
- user_role - (Optional) List of users that will be given the User role
- read_only_role - (Optional) List of users that will be given the Read-only role

###Attribute Reference
- id - ID of the environment.
- name - Name of the environment.

##cloudca_vpc
Create a vpc.

###Example usage
```
resource "cloudca_vpc" "my_vpc" {
	service_code = "compute-east"
	environment_name = "dev"
	name = "test-vpc"
	description = "This is a test vpc"
	vpc_offering = "Default VPC offering"
}
```
###Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- environment_name - (Required) Name of environment
- name - (Required) Name of the VPC
- description - (Required) Description of the VPC
- vpc_offering - (Required) The name of the VPC offering to use for the vpc
- network_domain - (Optional) A custom DNS suffix at the level of a network

###Attribute Reference
- id - ID of VPC.

##cloudca_tier
Create a tier.

###Example usage
```
resource "cloudca_tier" "my_tier" {
	service_code = "compute-east"
	environment_name = "dev"
	name = "test-tier"
	description = "This is a test tier"
	vpc_id = "8b46e2d1-bbc4-4fad-b3bd-1b25fcba4cec"
	network_offering = "Standard Tier"
	network_acl_id = "7d428416-263d-47cd-9270-2cdbdf222f57"
}
```
###Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- environment_name - (Required) Name of environment
- name - (Required) Name of the tier
- description - (Required) Description of the tier
- vpc_id - (Required) The ID of the vpc where the tier should be created
- network_offering - (Required) The name of the network offering to use for the tier
- network_acl - (Required) The name of the network ACL to use for the tier

###Attribute Reference
- id - ID of tier.

##cloudca_network_acl
Create a network ACL.

###Example usage
```
resource "cloudca_network_acl" "my_acl" {
	service_code = "compute-east"
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
- vpc_id - (Required) ID of the vpc where the network ACL should be created

###Attribute Reference
- id - ID of network ACL.
- name - Name of network ACL.

##cloudca_network_acl_rule
Create a network ACL rule.

###Example usage
```
resource "cloudca_network_acl_rule" "my_acl" {
	service_code = "compute-east"
	environment_name = "dev"
	rule_number = 55
	action = "Allow"
	protocol = "TCP"
	start_port = 80
	end_port = 80
	traffic_type = "Ingress"
	network_acl_id = "c0731f8b-92f0-4fac-9cbd-245468955fdf"
}
```
###Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- environment_name - (Required) Name of environment
- network_acl_id - (Required) ID of the network ACL where the rule should be created
- rule_number - (Required) Rule number of the network ACL rule
- cidr - (Required) Cidr of the network ACL rule
- action - (Required) Action of the network ACL rule (i.e. Allow or Deny)
- protocol - (Required) Protocol of the network ACL rule (i.e. TCP, UDP, ICMP or All)
- traffic_type - (Required) TrafficType of the network ACL rule (i.e. Ingress or Egress)
- icmp_type - (Optional) The ICMP type. Can only be used with ICMP protocol
- icmp_code - (Optional) The ICMP code. Can only be used with ICMP protocol
- start_port - (Optional) The start port. Can only be used with TCP/UDP protocol
- end_port - (Optional) The end port. Can only be used with TCP/UDP protocol

###Attribute Reference
- id - ID of network ACL.
- name - Name of network ACL.

##cloudca_instance
Create and starts an instance.

###Example usage
```
resource "cloudca_instance" "my_instance" {
	service_code = "compute-east"
	environment_name = "dev"
	name = "test-instance"
	network_id = "672016ef-05ee-4e88-b68f-ac9cc462300b"
	template = "CentOS 6.7 base (64bit)"
	compute_offering = "1vCPU.512MB"
	ssh_key_name = "my_ssh_key"
}
```
###Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- environment_name - (Required) Name of environment
- name - (Required) Name of instance
- network_id - (Required) The ID of the network where the instance should be created
- template - (Required) Name of template to use for the instance
- compute_offering - (Required) Name of the compute offering to use for the instance
- user_date - (Optional) User data to add to the instance
- ssh_key_name - (Optional) Name of the SSH key pair to attach to the instance. Mutually exclusive with public_key.
- public_key - (Optional) Public key to attach to the instance. Mutually exclusive with ssh_key_name.
- purge - (Optional) If true, then it will purge the instance on destruction

###Attribute Reference
- id - ID of instance.
- private_ip_id - ID of instance's private IP

##cloudca_public_ip
Acquires a public IP in a specific VPC. If you update any of the fields in the resource, then it will release this IP and recreate it.

###Example usage
```
resource "cloudca_public_ip" "my_publicip" {
	service_code = "compute-east"
	environment_name = "dev"
	vpc_id = "8b46e2d1-bbc4-4fad-b3bd-1b25fcba4cec"
}
```
###Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- environment_name - (Required) Name of environment
- vpc_id - (Required) The ID of the vpc to acquire the public IP

###Attribute Reference
- id - The public IP ID.
- ip_address - The public IP address

##cloudca_port_forwarding_rule
Manages port forwarding rules. Modifying any field will result in destruction and recreation of the rule.

When adding a port forwarding rule to the default private IP of an instance, only the instance id is required. Alternatively, the private_ip_id can be used on its own (for example when targeting an instance secondary IP).

###Example usage
```
resource "cloudca_port_forwarding_rule" "web_pfr" {
	service_code = "compute-east"
	environment_name = "dev"

	public_ip_id = "319f508f-089b-482d-af17-0f3360520c69"
	public_port_start = 80
	private_ip_id = "30face92-f1cf-4064-aa7f-008ea09ef7f0"
	private_port_start = 8080
	protocol = "TCP"
}
```

###Argument reference
- service_code - (Required)
- environment_name - (Required)
- private_ip_id - (Required) The private IP which should be used to create this rule
- private_port_start - (Required)
- private_port_end - (Optional) If not specified, defaults to the private start port
- public_ip_id - (Required) The public IP which should be used to create this rule
- public_port_start - (Required)
- public_port_end - (Optional) If not specified, defaults to the public start port
- protocol - (Required) The protocol to be used for this rule - must be TCP or UDP

###Attribute reference
- id - the rule ID
- public_ip - the public IP address of this rule
- private_ip - the private IP address of this rule
- instance_id - the instance associated with the private IP address of this rule
