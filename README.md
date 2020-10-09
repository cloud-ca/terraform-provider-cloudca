# terraform-provider-cloudca [![Build Status](https://github.com/cloud-ca/terraform-provider-cloudca/workflows/build/badge.svg)](https://github.com/cloud-ca/terraform-provider-cloudca/actions) [![GoDoc](https://godoc.org/github.com/cloud-ca/terraform-provider-cloudca?status.svg)](https://godoc.org/github.com/cloud-ca/terraform-provider-cloudca) [![Go Report Card](https://goreportcard.com/badge/github.com/cloud-ca/terraform-provider-cloudca)](https://goreportcard.com/report/github.com/cloud-ca/terraform-provider-cloudca) [![license](https://img.shields.io/github/license/cloud-ca/terraform-provider-cloudca.svg)](https://github.com/cloud-ca/terraform-provider-cloudca/blob/master/LICENSE)


Terraform provider for cloud.ca

Tested with Terraform version : 0.13.4

## Installation

### 0.13 and greater
1. Download the cloud.ca Terraform provider binary for your OS from the [releases page](https://github.com/cloud-ca/terraform-provider-cloudca/releases).
2. Create the folder ~/.terraform.d/plugins/cloud.ca/cloudca/cloudca
3. Under that, make a folder for the version, and underneath that, a folder for your download type (ex. linux\_amd64).
4. COpy the provider to that folder. The directory structure should look like `~/.terraform.d/plugins/cloud.ca/cloudca/cloudca/1.5.0/linux_amd64/terraform-provider-cloudca_v1.5.0`

### 0.12 and below
1. Download the cloud.ca Terraform provider binary for your OS from the [releases page](https://github.com/cloud-ca/terraform-provider-cloudca/releases).
2. Copy the provider to the plugin directory `~/.terraform.d/plugins`.

Alternate installation: [Terraform documentation](https://www.terraform.io/docs/plugins/basics.html)

## How to use

In your configuration file, define a variable that will hold your API key. This variable will have the value of the environment variable `TF_VAR_my_api_key`. Create a new "cloudca" provider with the `api_key`. Optionally, you can override the `api_url` field of the provider.

```hcl
variable "my_api_key" {}

provider "cloudca" {
    api_key = "${var.my_api_key}"
}
```

## Links

- [**Resources documentation**](https://github.com/cloud-ca/terraform-provider-cloudca/tree/master/doc)

## Build from source

Install [Go](https://golang.org/doc/install) (version 1.12 is required)

Download the provider source:

```Shell
go get github.com/cloud-ca/terraform-provider-cloudca
```

Compile the provider:

```Shell
cd $GOPATH/src/github.com/cloud-ca/terraform-provider-cloudca
make build
```

Copy the provider to the directory where terraform is located:

```Shell
sudo cp ./bin/$(go env GOOS)-$(go env GOARCH)/terraform-provider-cloudca_* $(dirname `which terraform`)
```

## Build for all OS/architectures

To build zip files containing the executables for each OS/architecture combination, use

```Shell
make build-all
```

## Prepare a Release

To prepare a new release, use one of the following:

```shell
make patch # e.g. move from v1.2.3 to v1.2.4
make minor # e.g. move from v1.2.3 to v1.3.0
make major # e.g. move from v1.2.3 to v2.0.0
```

or

```shell
make release version=x.y.z # where x, y and z are non-negative integers
```

also you can use `push=true` flag in all of the above to push the newly released tag to GitHub.

## License

This project is licensed under the terms of the MIT license.

```text
The MIT License (MIT)

Copyright (c) 2019 CloudOps

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
