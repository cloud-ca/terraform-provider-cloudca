package cloudca

import (
	"fmt"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"regexp"
	"strings"
)

func GetCloudCAResourceMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"cloudca_instance":             resourceCloudcaInstance(),
		"cloudca_environment":          resourceCloudcaEnvironment(),
		"cloudca_vpc":                  resourceCloudcaVpc(),
		"cloudca_tier":                 resourceCloudcaTier(),
		"cloudca_port_forwarding_rule": resourceCloudcaPortForwardingRule(),
		"cloudca_publicip":             resourceCloudcaPublicIp(),
		"cloudca_networkacl":			resourceCloudcaNetworkAcl(),
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

func retrieveVpcId(ccaRes *cloudca.Resources, name string) (id string, err error) {
	if isID(name) {
		return name, nil
	}
	vpcs, err := ccaRes.Vpcs.List()
	if err != nil {
		return "", err
	}
	for _, vpc := range vpcs {
		if strings.EqualFold(vpc.Name, name) {
			log.Printf("Found vpc: %+v", vpc)
			return vpc.Id, nil
		}
	}
	return "", fmt.Errorf("Vpc with name %s not found", name)
}
