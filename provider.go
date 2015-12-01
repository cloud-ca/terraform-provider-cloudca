package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/cloud-ca/terraform-cloudca/cloudca" 
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDCA_API_URL", "https://api.cloud.ca/v1"),
			},
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDCA_API_KEY", nil),
			},
		},
		ResourcesMap: mergeResourceMaps(
						cloudca.GetCloudCAResourceMap(),
		),
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIURL: d.Get("api_url").(string),
		APIKey: d.Get("api_key").(string),
	}

	return config.NewClient()
}

func mergeResourceMaps(resourceMaps ...map[string]*schema.Resource) map[string]*schema.Resource {
	mergedMap := map[string]*schema.Resource{}
	for _, resourceMap := range resourceMaps {
		for k, v := range resourceMap {
			mergedMap[k] = v
		}
	}
	return mergedMap
}