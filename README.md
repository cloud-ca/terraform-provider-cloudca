# terraform-cloudca

Terraform plugin for cloud.ca

# How to use

First step is to create a terraform configuration file.

In that file, define a variable that will hold your API key. This variable will have the value of the environment variable "TF_VAR_my_api_key". Create a new "cloudca" provider with the api_key. Optionally, you can override the api_url field of the provider.
```
var "my_api_key" {}

provider "cloudca" {
	api_key = "${var.my_api_key}"
}
```

Next step is to create a resource of that provider. 
Here, we are creating a new instance called "test-web-app" in the environment "dev" for the service "compute-east". 
```
resource "cloudca_instance" "web" {
	"service_code" = "compute-east"
	"environment_name" = "dev"
	"name" = "test-web-app"
	"template" = "CoreOS Stable"
	"compute_offering" = "1vCPU.1GB"
	"network" = "Web-Tier"
}
```

Alternatively, ids can be used instead of names.
```
resource "cloudca_instance" "web" {
  ...
	"network" = "db4c1e34-e1cd-4ca3-acfd-3b00042c49b7"
}
```

#License

This project is licensed under the terms of the MIT license.
