package cloudca

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func GetCloudCAResourceMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
			"cloudca_instance": resourceCloudcaInstance(),
		}
}