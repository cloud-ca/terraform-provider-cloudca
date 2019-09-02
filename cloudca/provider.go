package cloudca

import (
	"os"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDCA_API_URL", "https://api.cloud.ca/v1"),
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDCA_API_KEY", nil),
			},
		},
		ResourcesMap: mergeResourceMaps(
			GetCloudCAResourceMap(),
		),
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	insecure, _ := strconv.ParseBool(os.Getenv("CLOUD_CA_INSECURE_CONNECTION"))
	config := Config{
		APIURL:   d.Get("api_url").(string),
		APIKey:   d.Get("api_key").(string),
		Insecure: insecure,
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
