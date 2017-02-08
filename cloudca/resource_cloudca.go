package cloudca

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/cloud-ca/go-cloudca"
	"github.com/cloud-ca/go-cloudca/api"
	"github.com/cloud-ca/go-cloudca/services/cloudca"
	"github.com/hashicorp/terraform/helper/schema"
)

func GetCloudCAResourceMap() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"cloudca_instance":             resourceCloudcaInstance(),
		"cloudca_environment":          resourceCloudcaEnvironment(),
		"cloudca_vpc":                  resourceCloudcaVpc(),
		"cloudca_tier":                 resourceCloudcaTier(),
		"cloudca_port_forwarding_rule": resourceCloudcaPortForwardingRule(),
		"cloudca_public_ip":            resourceCloudcaPublicIp(),
		"cloudca_volume":               resourceCloudcaVolume(),
		"cloudca_load_balancer_rule":   resourceCloudcaLoadBalancerRule(),
		"cloudca_network_acl":          resourceCloudcaNetworkAcl(),
		"cloudca_network_acl_rule":     resourceCloudcaNetworkAclRule(),
		"cloudca_static_nat":           resourceCloudcaStaticNat(),
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

func readIntFromString(valStr string) int {
	var valInt int
	if valStr != "" {
		valInt, _ = strconv.Atoi(valStr)
	}
	return valInt
}

// Provides a common, simple way to deal with 404s.
func handleNotFoundError(err error, d *schema.ResourceData) error {
	if ccaError, ok := err.(api.CcaErrorResponse); ok {
		if ccaError.StatusCode == 404 {
			fmt.Errorf("Entity (id=%s) not found", d.Id())
			d.SetId("")
			return nil
		}
	}
	return err
}

// Deals with all of the casting done to get a cloudca.Resources.
func getResources(d *schema.ResourceData, meta interface{}) cloudca.Resources {
	client := meta.(*cca.CcaClient)
	_resources, _ := client.GetResources(d.Get("service_code").(string), d.Get("environment_name").(string))
	return _resources.(cloudca.Resources)
}

// Deals with all of the casting done to get a cloudca.Resources.
func getResourcesForEnvironmentId(client *cca.CcaClient, environmentId string) (cloudca.Resources, error) {
	environment, err := client.Environments.Get(environmentId)
	if err != nil {
		return cloudca.Resources{}, err
	}
	resources, err := client.GetResources(environment.ServiceConnection.ServiceCode, environment.Name)
	if err != nil {
		return cloudca.Resources{}, err
	}
	return resources.(cloudca.Resources), nil
}
