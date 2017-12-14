# terraform-provider-cloudca

Terraform provider for cloud.ca

Tested with Terraform version : 0.11.1

# Installation

1. Download the cloud.ca Terraform provider binary for your OS from the [releases page](https://github.com/cloud-ca/terraform-provider-cloudca/releases).
2. Copy the provider to the plugin directory `~/.terraform.d/plugins`.

Alternate installation: [Terraform documentation](https://www.terraform.io/docs/plugins/basics.html)

# How to use

In your configuration file, define a variable that will hold your API key. This variable will have the value of the environment variable "TF_VAR_my_api_key". Create a new "cloudca" provider with the api_key. Optionally, you can override the api_url field of the provider.
```hcl
variable "my_api_key" {}

provider "cloudca" {
	api_key = "${var.my_api_key}"
}
```

# Links
- [**Resources documentation**](https://github.com/cloud-ca/terraform-cloudca/blob/master/cloudca/README.md)

# Build from source
Install [Go](https://golang.org/doc/install) (version 1.8 is required)

Download the provider source:
```Shell
$ go get github.com/cloud-ca/terraform-provider-cloudca
```
Compile the provider:
```Shell
$ cd $GOPATH/src/github.com/cloud-ca/terraform-provider-cloudca
$ make init
$ make build
```
Copy the provider to the directory where terraform is located:
```Shell
$ sudo cp terraform-provider-cloudca $(dirname `which terraform`)
```
# Build for all OS/architectures
To build zip files containing the executables for each OS/architecture combination, use
```Shell
$ make build-all
```
# License

This project is licensed under the terms of the MIT license.
