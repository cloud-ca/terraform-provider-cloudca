# terraform-cloudca

Terraform provider for cloud.ca

# Installation
Download the provider:
```
$ go get github.com/cloud-ca/terraform-cloudca
```
Download and install the dependencies of the provider:
```
$ cd $GOPATH/src/github.com/cloud-ca/terraform-cloudca
$ godep restore
```
Compile the provider:
```
$ go build -o terraform-provider-cloudca
```
Copy it to the directory where terraform is located:
```
$ sudo cp terraform-provider-cloudca $(dirname `which terraform`)
```
# How to use

In your configuration file, define a variable that will hold your API key. This variable will have the value of the environment variable "TF_VAR_my_api_key". Create a new "cloudca" provider with the api_key. Optionally, you can override the api_url field of the provider.
```
var "my_api_key" {}

provider "cloudca" {
	api_key = "${var.my_api_key}"
}
```

# Links
- [**Resources documentation**](https://github.com/cloud-ca/terraform-cloudca/blob/master/cloudca/README.md)

# License

This project is licensed under the terms of the MIT license.
