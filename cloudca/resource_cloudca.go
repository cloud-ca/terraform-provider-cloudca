package cloudca

import (
	"regexp"
	"github.com/hashicorp/terraform/helper/schema"
)

func GetCloudCAResourceMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
			"cloudca_instance": resourceCloudcaInstance(),
			"cloudca_environment": resourceCloudcaEnvironment(),
			"cloudca_vpc": resourceCloudcaVpc(),
		}
}

func setValueOrID(d *schema.ResourceData, key string, value string, id string) {
	if isID(d.Get(key).(string)) {
		d.Set(key, id)
	} else {
		d.Set(key, value)
	}
}

func isID(id string) bool {
	re := regexp.MustCompile(`^([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$`)
	return re.MatchString(id)
}