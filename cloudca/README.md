#Resources
- [**cloudca_environment**](#cloudca_environment)
- [**cloudca_vpc**](#cloudca_vpc)
- [**cloudca_tier**](#cloudca_tier)
- [**cloudca_instance**](#cloudca_instance)
- [**cloudca_publicip**](#cloudca_publicip)
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
	admin_role_users = ["pat"]
	read_only_role_users = ["franz","bob"]
}
```
###Argument Reference
The following arguments are supported:
- service_code - (Required) Service code
- organization_code - (Required) Organization's entry point, i.e. <entry_point>.cloud.ca
- name - (Required) Name of environment to be created. Must be lower case, contain alphanumeric charaters, underscores or dashes
- description - (Required) Description for the environment
- admin_role_users - (Optional) List of users that will be given the Environment Admin role
- user_role_users - (Optional) List of users that will be given the User role
- read_only_role_users - (Optional) List of users that will be given the Read-only role

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
	vpc_id = "8b46e2d1-bbc4-4fad-b3bd-1b25fcba4cec"
	network_offering = "Standard Tier"
	network_acl = "default_allow"
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
resource "cloudca_tier" "my_instance" {
	service_code = "compute-east"
	environment_name = "dev"
	name = "test-tier"
	description = "This is a test tier"
	vpc_id = "8b46e2d1-bbc4-4fad-b3bd-1b25fcba4cec"
	network_offering = "Standard Tier"
	network_acl = "default_allow"
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

##cloudca_publicip
Acquires a public IP in a specific VPC. If you update any of the fields in the resource, then it will release this IP and recreate it.

###Example usage
```
resource "cloudca_publicip" "my_publicip" {
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
	instance_id = "5ec8564e-4793-4d2a-a4f3-218071c69c7e"
	private_ip_id = "30face92-f1cf-4064-aa7f-008ea09ef7f0"
	private_port_start = 8080
	private_port_end = 8080
	public_port_start = 80
	public_port_end = 80 
}
```

###Argument reference
- service_code - (Required)
- environment_name - (Required)
- instance_id - (Optional) If specified without a private_ip_id, applies the rule to the primary private IP address of this instance. **At least one of private_ip_id and instance_id must be provided**
- private_ip_id - (Optional) The private IP which should be used to create this rule. **At least one of private_ip_id and instance_id must be provided**
- private_port_start - (Required
- private_port_end - (Required)
- public_ip_id - (Required) The public IP which should be used to create this rule
- public_port_start - (Required)
- public_port_end - (Required)

###Attribute reference
- id - the rule ID