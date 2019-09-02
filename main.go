package main

import (
	"github.com/cloud-ca/terraform-provider-cloudca/cloudca"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cloudca.Provider,
	})
}
