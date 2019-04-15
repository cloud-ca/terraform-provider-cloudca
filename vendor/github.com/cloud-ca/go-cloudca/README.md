# go-cloudca

[![GoDoc](https://godoc.org/github.com/cloud-ca/go-cloudca?status.svg)](https://godoc.org/github.com/cloud-ca/go-cloudca)
[![Build Status](https://circleci.com/gh/cloud-ca/go-cloudca.svg?style=svg)](https://circleci.com/gh/cloud-ca/go-cloudca)
[![license](https://img.shields.io/github/license/cloud-ca/go-cloudca.svg)](https://github.com/cloud-ca/go-cloudca/blob/master/LICENSE)

A cloud.ca client for the Go programming language

## How to use

Import

```go
import "github.com/cloud-ca/go-cloudca"

/* import the services you need */
import "github.com/cloud-ca/go-cloudca/services/cloudca"
```

Create a new CcaClient.

```go
ccaClient := cca.NewCcaClient("[your-api-key]")
```

Retrieve the list of environments

```go
environments, _ := ccaClient.Environments.List()
```

Get the ServiceResources object for a specific environment and service. Here, we assume that it is a cloudca service.

```go
resources, _ := ccaClient.GetResources("[service-code]", "[environment-name]")
ccaResources := resources.(cloudca.Resources)
```

Now with the cloudca ServiceResources object, we can execute operations on cloudca resources in the specified environment.

Retrieve the list of instances in the environment.

```go
instances, _ := ccaResources.Instances.List()
```

Get a specific volume in the environment.

```go
volume, _ := ccaResources.Volumes.Get("[some-volume-id]")
```

Create a new instance in the environment.

```go
createdInstance, _ := ccaResources.Instances.Create(cloudca.Instance{
    Name: "[new-instance-name]",
    TemplateId: "[some-template-id]",
    ComputeOfferingId:"[some-compute-offering-id]",
    NetworkId:"[some-network-id]",
})
```

## Handling Errors

When trying to get a volume with a bogus id, an error will be returned.

```go
// Get a volume with a bogus id
_, err := ccaResources.Volumes.Get("[some-volume-id]")
```

Two types of error can occur: an unexpected error (ex: unable to connect to server) or an API error (ex: service resource not found)
If an error has occured, then we first try to cast the error into a CcaErrorResponse. This object contains the HTTP status code returned by the server, an error code and a list of CcaError objects. If it's not a CcaErrorResponse, then the error was not returned by the API.

```go
if err != nil {
    if errorResponse, ok := err.(api.CcaErrorResponse); ok {
        if errorResponse.StatusCode == api.NOT_FOUND {
            fmt.Println("Volume was not found")
        } else {
            // Can get more details from the CcaErrors
            fmt.Println(errorResponse.Errors)
        }
    } else {
        // handle unexpected error
        panic("Unexpected error")
    }
}
```

## License

This project is licensed under the terms of the MIT license.

```txt
The MIT License (MIT)

Copyright (c) 2019 cloud.ca

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
